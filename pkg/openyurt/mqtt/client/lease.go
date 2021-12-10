package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	coordinationv1 "k8s.io/api/coordination/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type LeasesGetter interface {
	Leases(namespace string) LeaseInstance
}

type LeaseInstance interface {
	PublishTopicor
	Get(ctx context.Context, name string, options metav1.GetOptions) (result *coordinationv1.Lease, err error)
	Create(ctx context.Context, lease *coordinationv1.Lease, opts metav1.CreateOptions) (result *coordinationv1.Lease, err error)
	Update(ctx context.Context, lease *coordinationv1.Lease, opts metav1.UpdateOptions) (result *coordinationv1.Lease, err error)
}

type leases struct {
	nodename  string
	namespace string
	index     cache.Indexer
	client    MessageSendor
}

func (l *leases) GetPublishDeleteTopic(name string) string {
	return filepath.Join(l.GetPublishPreTopic(), name, "delete")
}

func (l *leases) GetPublishCreateTopic(name string) string {
	return filepath.Join(l.GetPublishPreTopic(), name, "create")
}

func (l *leases) GetPublishUpdateTopic(name string) string {
	return filepath.Join(l.GetPublishPreTopic(), name, "update")
}

func (l *leases) GetPublishPatchTopic(name string) string {
	return filepath.Join(l.GetPublishPreTopic(), name, "patch")
}

func (l *leases) GetPublishPreTopic() string {
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

func (l *leases) Create(ctx context.Context, lease *coordinationv1.Lease, opts metav1.CreateOptions) (*coordinationv1.Lease, error) {

	createTopic := l.GetPublishCreateTopic(lease.GetName())
	data := PublishCreateData(l.nodename, lease, opts)

	if err := l.client.Send(createTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish create lease[%s][%s] data error %v", lease.Namespace, lease.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("Publish create lease data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[Indentify %s] from timeoutcache timeout  when create lease", data.Identity)
		return lease, errors.NewTimeoutError("lease", 5)
	}
	nl := &coordinationv1.Lease{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return lease, errors.NewInternalError(err)
	}

	klog.Infof("###### Create lease[%s][%s] by topic[%s]: finnal errorinfo %v", lease.GetNamespace(), lease.GetName(), createTopic, errInfo)
	return nl, errInfo
}

func (l *leases) Update(ctx context.Context, lease *coordinationv1.Lease, opts metav1.UpdateOptions) (result *coordinationv1.Lease, err error) {
	updateTopic := l.GetPublishUpdateTopic(lease.GetName())
	data := PublishUpdateData(l.nodename, lease, opts)

	if err := l.client.Send(updateTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish update lease[%s][%s] data error %v", lease.Namespace, lease.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("Publish update lease data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[Indentify %s] from timeoutcache timeout  when update lease", data.Identity)
		return lease, errors.NewTimeoutError("lease", 5)
	}
	nl := &coordinationv1.Lease{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return lease, errors.NewInternalError(err)
	}

	klog.Infof("###### Update lease[%s][%s] by topic[%s]: finnal errorinfo %v", lease.GetNamespace(), lease.GetName(), updateTopic, errInfo)
	return nl, errInfo
}

func newLeases(nodename, namespace string, index cache.Indexer, c MessageSendor) *leases {
	return &leases{
		nodename:  nodename,
		namespace: namespace,
		index:     index,
		client:    c,
	}
}

var _ LeaseInstance = &leases{}
