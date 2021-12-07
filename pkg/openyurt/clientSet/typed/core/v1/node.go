package v1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// FakeNodes implements NodeInterface
type FakeNodes struct {
	LocalClient client.KubeletOperatorInterface
	fakecorev1.FakeNodes
}

// Create takes the representation of a node and creates it.  Returns the server's representation of the node, and an error, if there is any.
func (c *FakeNodes) Create(ctx context.Context, node *corev1.Node, opts metav1.CreateOptions) (result *corev1.Node, err error) {
	panic("need to implement: node create")
}

// Patch applies the patch and returns the patched node.
func (c *FakeNodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *corev1.Node, err error) {
	panic("need to implement: node patch")
}

func (c *FakeNodes) Get(ctx context.Context, name string, opts metav1.GetOptions) (*corev1.Node, error) {
	return c.LocalClient.Nodes().Get(ctx, name, opts)
}
