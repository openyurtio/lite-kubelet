package v1

import (
	v1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	fakecoordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

type FakeCoordinationV1 struct {
	LocalClient client.KubeletOperatorInterface
	fakecoordinationv1.FakeCoordinationV1
}

func (f *FakeCoordinationV1) Leases(namespace string) v1.LeaseInterface {
	return &FakeLeases{
		LocalClient: f.LocalClient,
		FakeLeases:  fakecoordinationv1.FakeLeases{Fake: &f.FakeCoordinationV1},
		ns:          namespace,
	}
}
