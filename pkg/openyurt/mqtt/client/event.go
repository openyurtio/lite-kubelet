package client

import (
	"context"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type EventsGetter interface {
	Events(namespace string) EventInstance
}

type EventInstance interface {
	Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error)
}

type events struct {
	namespace string
}

func (e *events) Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error) {
	klog.Warningf("implement me: create event %++v", *event)
	return event, nil
}

func (e *events) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error) {
	klog.Warningf("implement me: patch event %++v", string(data))
	return result, nil
}

func newEvents(namespace string) *events {
	return &events{
		namespace: namespace,
	}
}

var _ EventInstance = &events{}
