package client

import (
	"context"

	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type LeasesGetter interface {
	Leases(namespace string) LeaseInstance
}

type LeaseInstance interface {
	Get(ctx context.Context, name string, options v1.GetOptions) (result *coordinationv1.Lease, err error)
	Create(ctx context.Context, lease *coordinationv1.Lease, opts v1.CreateOptions) (result *coordinationv1.Lease, err error)
	Update(ctx context.Context, lease *coordinationv1.Lease, opts v1.UpdateOptions) (result *coordinationv1.Lease, err error)
}

type leases struct {
	namespace string
}

func (l *leases) Get(ctx context.Context, name string, options v1.GetOptions) (result *coordinationv1.Lease, err error) {
	// GET 方法返回空即可
	klog.Warningf("implement me: get lease [%s][%s]", l.namespace, name)
	return nil, apierrors.NewNotFound(corev1.Resource("leases"), name)
}

func (l *leases) Create(ctx context.Context, lease *coordinationv1.Lease, opts v1.CreateOptions) (result *coordinationv1.Lease, err error) {
	klog.Warningf("implement me: create lease ")
	return lease, nil
}

func (l *leases) Update(ctx context.Context, lease *coordinationv1.Lease, opts v1.UpdateOptions) (result *coordinationv1.Lease, err error) {
	klog.Warningf("implement me: update lease ")
	return lease, nil
}

func newLeases(namespace string) *leases {
	return &leases{
		namespace: namespace,
	}
}

var _ LeaseInstance = &leases{}
