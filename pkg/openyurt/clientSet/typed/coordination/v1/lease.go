package v1

import (
	"context"

	coordinationv1 "k8s.io/api/coordination/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakecoordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// FakeLeases implements LeasesInterface
type FakeLeases struct {
	LocalClient client.KubeletOperatorInterface
	fakecoordinationv1.FakeLeases
	ns string
}

// Get takes name of the lease, and returns the corresponding lease object, and an error if there is any.
func (c *FakeLeases) Get(ctx context.Context, name string, options v1.GetOptions) (result *coordinationv1.Lease, err error) {
	return c.LocalClient.Leases(c.ns).Get(ctx, name, options)
}

// Create takes the representation of a lease and creates it.  Returns the server's representation of the lease, and an error, if there is any.
func (c *FakeLeases) Create(ctx context.Context, lease *coordinationv1.Lease, opts v1.CreateOptions) (result *coordinationv1.Lease, err error) {
	return c.LocalClient.Leases(c.ns).Create(ctx, lease, opts)
}

// Update takes the representation of a lease and updates it. Returns the server's representation of the lease, and an error, if there is any.
func (c *FakeLeases) Update(ctx context.Context, lease *coordinationv1.Lease, opts v1.UpdateOptions) (result *coordinationv1.Lease, err error) {
	panic("need to implement")
}
