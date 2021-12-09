package client

import (
	"context"
	"fmt"
	"path/filepath"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
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
	klog.Warningf("implement me create node ")
	return node, nil
}

func (n *nodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Node, err error) {
	klog.Warningf("implement me patch node")
	return nil, nil
}

func (n *nodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *corev1.Node, err error) {

	klog.Infof("###### Prepare to Get Node %s from cache ", name)
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

	klog.Infof("###### Get Node %s from cache succefully", name)
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
