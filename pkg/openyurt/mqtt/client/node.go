package client

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
)

type NodesGetter interface {
	Nodes() NodeInstance
}

type NodeInstance interface {
	PublishTopicor
	Create(ctx context.Context, node *corev1.Node, opts v1.CreateOptions) (result *corev1.Node, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Node, err error)
	Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Node, err error)
}

type nodes struct {
	nodename string
	index    cache.Indexer
	client   MessageSendor
}

func (n *nodes) GetPublishGetTopic(name string) string {
	return filepath.Join(n.GetPublishPreTopic(), name, "get")
}
func (n *nodes) GetPublishDeleteTopic(name string) string {
	return filepath.Join(n.GetPublishPreTopic(), name, "delete")
}

func (n *nodes) GetPublishCreateTopic(name string) string {
	return filepath.Join(n.GetPublishPreTopic(), name, "create")
}

func (n *nodes) GetPublishUpdateTopic(name string) string {
	return filepath.Join(n.GetPublishPreTopic(), name, "update")
}

func (n *nodes) GetPublishPatchTopic(name string) string {
	return filepath.Join(n.GetPublishPreTopic(), name, "patch")
}

func (n *nodes) GetPublishPreTopic() string {
	return filepath.Join(MqttEdgePublishRootTopic, "nodes")
}

func (n *nodes) Create(ctx context.Context, node *corev1.Node, opts v1.CreateOptions) (result *corev1.Node, err error) {

	createTopic := n.GetPublishCreateTopic(node.GetName())
	data := PublishCreateData(true, n.nodename, node, opts)

	if err := n.client.Send(createTopic, 1, false, data, time.Second*5); err != nil {
		klog.Errorf("Publish create node[%s] data error %v", node.Name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish create node data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeout cache timeout, when create node", data.Identity)
		return node, errors.NewTimeoutError("node", 5)
	}
	nl := &corev1.Node{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("publish ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return node, errors.NewInternalError(err)
	}

	klog.V(4).Infof("[%s] Create node [%s] by topic[%s]: errorInfo %v", ackdata.Identity, node.GetName(), createTopic, errInfo)
	return nl, errInfo
}

func (n *nodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Node, err error) {
	patchTopic := n.GetPublishPatchTopic(name)
	patchData := PublishPatchData(true, n.nodename, name, "", nil, pt, data, opts, subresources...)

	if err := n.client.Send(patchTopic, 1, false, patchData, time.Second*5); err != nil {
		klog.Errorf("Publish patch node[%s] data error %v", name, err)
		return nil, apierrors.NewInternalError(fmt.Errorf("publish patch node data error %v", err))
	}
	ackdata, ok := GetDefaultTimeoutCache().Pop(patchData.Identity, time.Second*5)
	if !ok {
		klog.Errorf("Get ack data[%s] from timeout cache timeout, when patch node", patchData.Identity)
		return nil, errors.NewTimeoutError("node", 5)
	}
	nl := &corev1.Node{}
	errInfo, err := ackdata.UnmarshalPublishAckData(nl)
	if err != nil {
		klog.Errorf("publish ack data unmarshal error %v,data:\n%v", err, *ackdata)
		return nil, errors.NewInternalError(err)
	}

	klog.V(4).Infof("Patch node [%s] by topic[%s]: errorInfo %v", name, patchTopic, errInfo)
	return nl, errInfo
}

func (n *nodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Node, err error) {

	klog.Infof("Prepare to Get Node %s from cache ", name)
	obj, exists, err := n.index.GetByKey(name)
	if err != nil {
		klog.Errorf("Cache index get node %s error %v", name, err)
		return nil, err
	}
	if !exists {
		klog.Errorf("Can not get node %++v", name)
		return nil, apierrors.NewNotFound(corev1.Resource("nodes"), name)
	}

	finnal, ok := obj.(*corev1.Node)
	if !ok {
		klog.Errorf("Cache obj convert to *corev1.Node error", err)
		return nil, apierrors.NewInternalError(fmt.Errorf("cache obj convert to *corev1.Node error"))
	}

	klog.V(4).Infof("Get Node %s from cache succefully", name)
	return finnal, nil
}

func newNodes(nodename string, index cache.Indexer, c MessageSendor) *nodes {
	return &nodes{
		nodename: nodename,
		index:    index,
		client:   c,
	}
}

var _ NodeInstance = &nodes{}
