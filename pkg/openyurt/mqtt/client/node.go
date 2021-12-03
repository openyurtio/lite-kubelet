package client

import (
	"context"

	"k8s.io/klog/v2"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type NodesGetter interface {
	Nodes() NodeInstance
}

type NodeInstance interface {
	Create(ctx context.Context, node *corev1.Node, opts v1.CreateOptions) (result *corev1.Node, err error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Node, err error)
}

type nodes struct {
}

func (n *nodes) Create(ctx context.Context, node *corev1.Node, opts v1.CreateOptions) (result *corev1.Node, err error) {
	klog.Warningf("implement me create node ")
	return node, nil
}

func (n *nodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Node, err error) {
	klog.Warningf("implement me patch node")
	return nil, nil
}

func newNodes() *nodes {
	return &nodes{}
}

var _ NodeInstance = &nodes{}
