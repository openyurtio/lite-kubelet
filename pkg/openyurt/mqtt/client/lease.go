package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type LeasesGetter interface {
	Leases(namespace string) LeaseInstance
}

type LeaseInstance interface {
	Topicor
	Get(ctx context.Context, name string, options metav1.GetOptions) (result *coordinationv1.Lease, err error)
	Create(ctx context.Context, lease *coordinationv1.Lease, opts metav1.CreateOptions) (result *coordinationv1.Lease, err error)
	Update(ctx context.Context, lease *coordinationv1.Lease, opts metav1.UpdateOptions) (result *coordinationv1.Lease, err error)
}

type leases struct {
	namespace string
	index     cache.Indexer
	client    MessageSendor
}

func (l *leases) GetPreTopic() string {
	return filepath.Join(MqttEdgePublishRootTopic, "leases", l.namespace)
}

func (l *leases) Get(ctx context.Context, name string, options metav1.GetOptions) (*coordinationv1.Lease, error) {
	// GET 方法返回空即可
	t := &coordinationv1.Lease{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: l.namespace,
			Name:      name,
		},
	}
	obj, exists, err := l.index.Get(t)
	if err != nil {
		klog.Errorf("Cache index get lease %++v error %v", *t, err)
		return nil, err
	}
	if !exists {
		klog.Errorf("Can not get lease %++v", *t)
		return nil, apierrors.NewNotFound(corev1.Resource("leases"), name)
	}

	finnal, ok := obj.(*coordinationv1.Lease)
	if !ok {
		klog.Errorf("Cache obj convert to *coordinationv1.Lease error", *t, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("cache obj convert to *coordinationv1.Lease error"))
	}
	klog.Infof("###### Get lease [%s][%s] from cache succefully", l.namespace, name)
	return finnal, nil
}

func (l *leases) Create(ctx context.Context, lease *coordinationv1.Lease, opts metav1.CreateOptions) (result *coordinationv1.Lease, err error) {
	createTopic := filepath.Join(l.GetPreTopic(), lease.GetName())
	klog.V(4).Infof("Lease create topic %s", createTopic)

	if err := l.client.Send(createTopic, 1, false, lease, time.Second*5); err != nil {
		klog.Errorf("Publish lease[%s][%s] data error %v", lease.Namespace, lease.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("Publish Lease data error %v", err))
	}

	klog.Infof("###### Create lease[%s][%s] successfully", lease.GetNamespace(), lease.GetName())
	return lease, nil
}

func (l *leases) Update(ctx context.Context, lease *coordinationv1.Lease, opts metav1.UpdateOptions) (result *coordinationv1.Lease, err error) {
	klog.Warningf("implement me: update lease ")
	return lease, nil
}

func newLeases(namespace string, index cache.Indexer, c MessageSendor) *leases {
	return &leases{
		namespace: namespace,
		index:     index,
		client:    c,
	}
}

var _ LeaseInstance = &leases{}
