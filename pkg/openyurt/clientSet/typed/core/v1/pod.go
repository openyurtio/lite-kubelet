package v1

import (
	"context"

	"k8s.io/apimachinery/pkg/types"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// FakePods implements PodInterface
type FakePods struct {
	MQTTClient client.KubeletOperatorInterface
	fakecorev1.FakePods
	ns string
}

// Create takes the representation of a pod and creates it.  Returns the server's representation of the pod, and an error, if there is any.
func (c *FakePods) Create(ctx context.Context, pod *corev1.Pod, opts v1.CreateOptions) (result *corev1.Pod, err error) {
	panic("need to implement: pod create")
}

// Delete takes name of the pod and deletes it. Returns an error if one occurs.
func (c *FakePods) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	panic("need to implement:pod delete")
}

// Patch applies the patch and returns the patched pod.
func (c *FakePods) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Pod, err error) {
	panic("need to implement:pod patch")
}
