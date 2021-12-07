package client

import (
	"context"
	"path/filepath"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
)

type EventsGetter interface {
	Events(namespace string) EventInstance
}

type EventInstance interface {
	Topicor
	Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error)
}

type events struct {
	namespace string
	index     cache.Indexer
	client    MessageSendor
}

func (e *events) GetPreTopic() string {
	return filepath.Join(MqttEdgePublishRootTopic, "events", e.namespace)
}

func (e *events) Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error) {
	klog.Warningf("implement me: create event %++v", *event)
	return event, nil
}

func (e *events) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error) {
	klog.Warningf("implement me: patch event %++v", string(data))
	return result, nil
}

func newEvents(namespace string, index cache.Indexer, c MessageSendor) *events {
	return &events{
		namespace: namespace,
		index:     index,
		client:    c,
	}
}

var _ EventInstance = &events{}
