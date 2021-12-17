package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type PodsGetter interface {
	Pods(namespace string) PodInstance
}

type PodInstance interface {
	PublishTopicor
	Create(ctx context.Context, pod *corev1.Pod, opts v1.CreateOptions) (result *corev1.Pod, err error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Pod, err error)
	Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Pod, err error)
}

type pods struct {
	nodename  string
	namespace string
	client    MessageSendor
	index     cache.Indexer
}

func (p *pods) GetPublishGetTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "get")
}

func (p *pods) GetPublishDeleteTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "delete")
}

func (p *pods) GetPublishCreateTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "create")
}

func (p *pods) GetPublishUpdateTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "update")
}

func (p *pods) GetPublishPatchTopic(name string) string {
	return filepath.Join(p.GetPublishPreTopic(), name, "patch")
}

func (p *pods) GetPublishPreTopic() string {
	return filepath.Join(MqttEdgePublishRootTopic, "pods", p.namespace)
}

func (p *pods) Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Pod, err error) {
	/*
		// GET 方法返回空即可
		t := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				Namespace: p.namespace,
				Name:      name,
			},
		}
		obj, exists, err := p.index.Get(t)
		if err != nil {
			klog.Errorf("Cache index get pod %++v error %v", *t, err)
			return nil, err
		}
		if !exists {
			klog.Errorf("Can not get pod %++v", *t)
			return nil, apierrors.NewNotFound(corev1.Resource("pods"), name)
		}

		finnal, ok := obj.(*corev1.Pod)
		if !ok {
			klog.Errorf("Cache obj convert to *corev1.Pod error", *t, err)
			return nil, apierrors.NewInternalError(fmt.Errorf("cache obj convert to *corev1.Pod error"))
		}
		klog.Infof("###### Get pod [%s][%s] from cache succefully", p.namespace, name)
		return finnal, nil
	*/

	getTopic := p.GetPublishGetTopic(name)
	data := PublishGetData(p.nodename, p.namespace, name, options)

	if err := p.client.Send(getTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish get pod[%s][%s] data error %v", p.namespace, name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish get pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout  when get pod", data.Identity)
		return nil, errors.NewTimeoutError("pod", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	klog.Infof("###### get pod[%s][%s] by topic[%s]: errorInfo %v", p.namespace, name, getTopic, errInfo)
	return nl, errInfo
}
func (p *pods) Create(ctx context.Context, pod *corev1.Pod, opts v1.CreateOptions) (result *corev1.Pod, err error) {

	createTopic := p.GetPublishCreateTopic(pod.GetName())
	data := PublishCreateData(p.nodename, pod, opts)

	if err := p.client.Send(createTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish create pod[%s][%s] data error %v", pod.Namespace, pod.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout when create pod", data.Identity)
		return pod, errors.NewTimeoutError("lease", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return pod, errors.NewInternalError(err)
	}

	klog.Infof("###### Create pod[%s][%s] by topic[%s]: errorInfo %v",
		pod.GetNamespace(), pod.GetName(), createTopic, errInfo)
	return nl, errInfo
}

func (p *pods) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {

	deleteTopic := p.GetPublishDeleteTopic(name)
	data := PublishDeleteData(p.nodename, name, p.namespace, opts)

	if err := p.client.Send(deleteTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish delete pod[%s][%s] data error %v", p.namespace, name, err)
		return apierrors.NewInternalError(fmt.Errorf("publish delete pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout when delete pod", data.Identity)
		return errors.NewTimeoutError("pods", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return errors.NewInternalError(err)
	}

	klog.Infof("###### delete pod[%s][%s] by topic[%s]: errorInfo %v",
		p.namespace, name, deleteTopic, errInfo)
	return errInfo
}

func (p *pods) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Pod, err error) {
	patchTopic := p.GetPublishPatchTopic(name)
	patchData := PublishPatchData(p.nodename, name, p.namespace,
		nil, pt, data, opts, subresources...)

	if err := p.client.Send(patchTopic, 1, false, patchData, time.Second*5); err != nil {
		klog.Errorf("Publish patch pod[%s] data error %v", name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish patch pod data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(patchData.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeout cache timeout, when patch pod", patchData.Identity)
		return nil, errors.NewTimeoutError("pod", 5)
	}
	nl := &corev1.Pod{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("publish ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	klog.Infof("###### Patch pod [%s] by topic[%s]: errorInfo %v", name, patchTopic, errInfo)
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
