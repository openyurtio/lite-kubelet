package v1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// FakeEvents implements EventInterface
type FakeEvents struct {
	MQTTClient client.KubeletOperatorInterface
	fakecorev1.FakeEvents
	ns string
}

// Create takes the representation of a event and creates it.  Returns the server's representation of the event, and an error, if there is any.
func (c *FakeEvents) Create(ctx context.Context, event *corev1.Event, opts v1.CreateOptions) (result *corev1.Event, err error) {
	panic("need to implement: events create")
}

// Patch applies the patch and returns the patched event.
func (c *FakeEvents) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *corev1.Event, err error) {
	panic("need to implement: event patch")
}
