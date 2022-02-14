/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package message

import (
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
)

type PodsGetter interface {
	Pods(namespace string) PodInstance
}

type PodInstance interface {
	Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (result *corev1.Pod, err error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1.Pod, err error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (result *corev1.Pod, err error)
}

type pods struct {
	nodename  string
	namespace string
	client    MessageSendor
	index     cache.Indexer
}

func (p *pods) deleteLocalCache(name string) error {
	podTmp := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: p.namespace,
		},
	}

	filePath, err := fileCache.NewDefaultFilePodDeps().GetFullFileName(podTmp)
	if err != nil {
		return fmt.Errorf("get object filename error %v", err)
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		klog.Errorf("Remove cache file %s error %v", filePath, err)
		// no need return
	} else {
		klog.Infof("Can not find pod[%s/%s] from cloud, so delete localcache file %s succefully",
			podTmp.GetNamespace(), podTmp.GetName(), filePath)
	}
	return nil
}

func (p *pods) Get(ctx context.Context, name string, options metav1.GetOptions) (result *corev1.Pod, err error) {
	podTmp := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: p.namespace,
		},
	}
	data := PublishGetData(ObjectTypePod, true, p.nodename, podTmp, options)

	if err := p.client.Send(data); err != nil {
		klog.Errorf("Publish get pod[%s][%s] data error %v", p.namespace, name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish get pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout  when get pod", data.Identity)
		return nil, errors.NewTimeoutError("pod", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	if errInfo != nil && apierrors.IsNotFound(errInfo) {
		if err := p.deleteLocalCache(name); err != nil {
			klog.Error("Delete local pod cache error %v", err)
			// no need return
		}
	}

	klog.V(4).Infof("Get pod[%s][%s] errorInfo %v", p.namespace, name, errInfo)
	return nl, errInfo
}
func (p *pods) Create(ctx context.Context, pod *corev1.Pod, opts metav1.CreateOptions) (result *corev1.Pod, err error) {

	data := PublishCreateData(ObjectTypePod, true, p.nodename, pod, opts)

	if err := p.client.Send(data); err != nil {
		klog.Errorf("Publish create pod[%s][%s] data error %v", pod.Namespace, pod.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout when create pod", data.Identity)
		return pod, errors.NewTimeoutError("lease", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return pod, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Create pod[%s][%s]  errorInfo %v",
		pod.GetNamespace(), pod.GetName(), errInfo)
	return nl, errInfo
}

func (p *pods) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	// first delete local manifiest

	podTmp := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: p.namespace,
			Name:      name,
		},
	}

	if err := p.deleteLocalCache(name); err != nil {
		klog.Error("Delete local pod cache error %v", err)
		// no need return
	}

	data := PublishDeleteData(ObjectTypePod, true, p.nodename, podTmp, opts)

	if err := p.client.Send(data); err != nil {
		klog.Errorf("Publish delete pod[%s][%s] data error %v", p.namespace, name, err)
		return apierrors.NewInternalError(fmt.Errorf("publish delete pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout when delete pod", data.Identity)
		return errors.NewTimeoutError("pods", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return errors.NewInternalError(err)
	}

	klog.V(4).Infof("Delete pod[%s][%s]  errorInfo %v",
		p.namespace, name, errInfo)
	return errInfo
}

func (p *pods) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1.Pod, err error) {
	patchData := PublishPatchData(ObjectTypePod, true, p.nodename, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: p.namespace,
		},
	}, pt, data, opts, subresources...)

	if err := p.client.Send(patchData); err != nil {
		klog.Errorf("Publish patch pod[%s] data error %v", name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish patch pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(patchData.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeout cache timeout, when patch pod", patchData.Identity)
		return nil, errors.NewTimeoutError("pod", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("publish ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Patch pod [%s] errorInfo %v", name, errInfo)
	return nl, errInfo
}

func newPods(nodename, namespace string, index cache.Indexer, c MessageSendor) *pods {
	return &pods{
		nodename:  nodename,
		namespace: namespace,
		client:    c,
		index:     index,
	}
}

var _ PodInstance = &pods{}
