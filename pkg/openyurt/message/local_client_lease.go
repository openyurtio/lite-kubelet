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
	"time"

	coordinationv1 "k8s.io/api/coordination/v1"
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

func (l *leases) Get(ctx context.Context, name string, options metav1.GetOptions) (*coordinationv1.Lease, error) {
	data := PublishGetData(ObjectTypeLease, true, l.nodename, &coordinationv1.Lease{ObjectMeta: metav1.ObjectMeta{
		Name:      name,
		Namespace: l.namespace,
	}}, options)

	if err := l.client.Send(data); err != nil {
		klog.Errorf("Publish get lease[%s][%s] data error %v", l.namespace, name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish get lease data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().PopWait(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout  when get lease", data.Identity)
		return nil, errors.NewTimeoutError("lease", 5)
	}
	nl := &coordinationv1.Lease{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Get lease[%s][%s]  errorInfo %v", l.namespace, name, errInfo)
	return nl, errInfo
}

func (l *leases) Create(ctx context.Context, lease *coordinationv1.Lease, opts metav1.CreateOptions) (*coordinationv1.Lease, error) {

	data := PublishCreateData(ObjectTypeLease, true, l.nodename, lease, opts)

	if err := l.client.Send(data); err != nil {
		klog.Errorf("Publish create lease[%s][%s] data error %v", lease.Namespace, lease.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create lease data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().PopWait(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout  when create lease", data.Identity)
		return lease, errors.NewTimeoutError("lease", 5)
	}
	nl := &coordinationv1.Lease{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return lease, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Create lease[%s][%s] errorInfo %v", lease.GetNamespace(), lease.GetName(), errInfo)
	return nl, errInfo
}

func (l *leases) Update(ctx context.Context, lease *coordinationv1.Lease, opts metav1.UpdateOptions) (result *coordinationv1.Lease, err error) {
	data := PublishUpdateData(ObjectTypeLease, true, l.nodename, lease, opts)

	if err := l.client.Send(data); err != nil {
		klog.Errorf("Publish update lease[%s][%s] data error %v", lease.Namespace, lease.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish update lease data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().PopWait(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeoutCache timeout  when update lease", data.Identity)
		return lease, errors.NewTimeoutError("lease", 5)
	}
	nl := &coordinationv1.Lease{}
	errInfo, err := ackdata.UnmarshalAckData(nl)
	if err != nil {
		klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return lease, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Update lease[%s][%s] errorInfo %v", lease.GetNamespace(), lease.GetName(), errInfo)
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
