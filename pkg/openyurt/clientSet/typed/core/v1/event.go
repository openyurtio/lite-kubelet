package v1

import (
	corev1 "k8s.io/api/core/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// FakeEvents implements EventInterface
type FakeEvents struct {
	LocalClient client.KubeletOperatorInterface
	fakecorev1.FakeEvents
	ns string
}

func (c *FakeEvents) CreateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {
	return c.LocalClient.Events(c.ns).CreateWithEventNamespace(event)
}

// UpdateWithEventNamespace is the same as a Update, except that it sends the request to the event.Namespace.
func (c *FakeEvents) UpdateWithEventNamespace(event *corev1.Event) (*corev1.Event, error) {
	return c.LocalClient.Events(c.ns).UpdateWithEventNamespace(event)

}
func (c *FakeEvents) PatchWithEventNamespace(event *corev1.Event, data []byte) (*corev1.Event, error) {
	return c.LocalClient.Events(c.ns).PatchWithEventNamespace(event, data)
}
