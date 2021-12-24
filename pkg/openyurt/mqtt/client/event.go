package client

import (
	"fmt"
	"path/filepath"
	"time"

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
	PublishTopicor
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

func (e *events) GetPublishGetTopic(name string) string {
	return filepath.Join(e.GetPublishPreTopic(), name, "get")
}

func (e *events) GetPublishDeleteTopic(name string) string {
	return filepath.Join(e.GetPublishPreTopic(), name, "delete")
}

func (e *events) GetPublishCreateTopic(name string) string {
	return filepath.Join(e.GetPublishPreTopic(), name, "create")
}

func (e *events) GetPublishUpdateTopic(name string) string {
	return filepath.Join(e.GetPublishPreTopic(), name, "update")
}

func (e *events) GetPublishPatchTopic(name string) string {
	return filepath.Join(e.GetPublishPreTopic(), name, "patch")
}

func (e *events) GetPublishPreTopic() string {
	if len(e.namespace) == 0 {
		e.namespace = "default"
	}
	return filepath.Join(MqttEdgePublishRootTopic, "events", e.namespace)
}

func (e *events) CreateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {
	createTopic := e.GetPublishCreateTopic(event.GetName())
	data := PublishCreateData(false, e.nodename, event, metav1.CreateOptions{})

	if err := e.client.Send(createTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish create event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create event data error %v", err))
	}
	// Do not deal with ack
	/*
		ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*2)
		if !ok {
			klog.Errorf("Get ack data[%s] from timeoutCache timeout  when create event", data.Identity)
			return event, errors.NewTimeoutError("lease", 2)
		}
		nl := &corev1.Event{}
		errInfo, err := ackdata.UnmarshalPublishAckData(nl)
		if err != nil {
			klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
			return event, errors.NewInternalError(err)
		}

		klog.Infof("Create event[%s][%s] by topic[%s]: errorInfo %v",
			event.GetNamespace(), event.GetName(), createTopic, errInfo)
		return nl, errInfo
	*/
	klog.V(4).Infof("Create event[%s][%s] by topic[%s] successfully", event.GetNamespace(), event.GetName(), createTopic)
	return event, nil
}

func (e *events) UpdateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {
	updateTopic := e.GetPublishUpdateTopic(event.GetName())
	data := PublishUpdateData(false, e.nodename, event, metav1.UpdateOptions{})

	if err := e.client.Send(updateTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish update event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish update event data error %v", err))
	}
	// do not dealwith ack
	klog.V(4).Infof("Update event[%s][%s] by topic[%s] successfully", event.GetNamespace(), event.GetName(), updateTopic)
	return event, nil
}

func (e *events) PatchWithEventNamespace(event *corev1.Event, data []byte) (*corev1.Event, error) {
	patchTopic := e.GetPublishPatchTopic(event.GetName())
	pathData := PublishPatchData(false, e.nodename, event.GetName(), event.GetNamespace(),
		event, types.StrategicMergePatchType, data, metav1.PatchOptions{})

	if err := e.client.Send(patchTopic, 1, false, pathData, time.Second*5); err != nil {
		klog.Errorf("Publish patch event[%s][%s] data error %v", event.Namespace, event.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish patch event data error %v", err))
	}
	// do not dealwith ack
	klog.V(4).Infof("Patch event[%s][%s] by topic[%s] successfully", event.GetNamespace(), event.GetName(), patchTopic)
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
