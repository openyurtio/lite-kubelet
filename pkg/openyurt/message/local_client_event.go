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
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type EventsGetter interface {
	Events(namespace string) EventInstance
}

type EventInstance interface {
	CreateWithEventNamespace(event *corev1.Event) (*corev1.Event, error)
	UpdateWithEventNamespace(event *corev1.Event) (*corev1.Event, error)
	PatchWithEventNamespace(event *corev1.Event, data []byte) (*corev1.Event, error)
}

type events struct {
	nodename  string
	namespace string
	index     cache.Indexer
	client    MessageSendor
}

func (e *events) CreateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {

	if err := e.client.Send(PublishCreateData(ObjectTypeEvent, false, e.nodename, event, metav1.CreateOptions{})); err != nil {
		klog.Errorf("Publish create event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create event data error %v", err))
	}
	klog.V(4).Infof("Create event[%s][%s] successfully", event.GetNamespace(), event.GetName())
	return event, nil
}

func (e *events) UpdateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {
	if err := e.client.Send(PublishUpdateData(ObjectTypeEvent, false, e.nodename, event, metav1.UpdateOptions{})); err != nil {
		klog.Errorf("Publish update event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish update event data error %v", err))
	}
	// do not dealwith ack
	klog.V(4).Infof("Update event[%s][%s] successfully", event.GetNamespace(), event.GetName())
	return event, nil
}

func (e *events) PatchWithEventNamespace(event *corev1.Event, data []byte) (*corev1.Event, error) {

	if err := e.client.Send(PublishPatchData(ObjectTypeEvent, false, e.nodename, event, types.StrategicMergePatchType, data, metav1.PatchOptions{})); err != nil {
		klog.Errorf("Publish patch event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish patch event data error %v", err))
	}
	// do not dealwith ack
	klog.V(4).Infof("Patch event[%s][%s] successfully", event.GetNamespace(), event.GetName())
	return event, nil
}

func newEvents(nodename, namespace string, index cache.Indexer, c MessageSendor) *events {
	return &events{
		nodename:  nodename,
		namespace: namespace,
		index:     index,
		client:    c,
	}
}

var _ EventInstance = &events{}
