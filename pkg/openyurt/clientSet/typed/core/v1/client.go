package v1

import (
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	"k8s.io/kubernetes/pkg/openyurt/message"
)

type FakeCoreV1 struct {
	LocalClient message.KubeletOperatorInterface
	fakecorev1.FakeCoreV1
}

func (f *FakeCoreV1) Pods(namespace string) corev1.PodInterface {
	return &FakePods{
		LocalClient: f.LocalClient,
		FakePods:    fakecorev1.FakePods{Fake: &f.FakeCoreV1},
		ns:          namespace,
	}
}

func (f *FakeCoreV1) Nodes() corev1.NodeInterface {
	return &FakeNodes{
		LocalClient: f.LocalClient,
		FakeNodes:   fakecorev1.FakeNodes{Fake: &f.FakeCoreV1},
	}
}

func (f *FakeCoreV1) Events(namespace string) corev1.EventInterface {
	return &FakeEvents{
		LocalClient: f.LocalClient,
		FakeEvents:  fakecorev1.FakeEvents{Fake: &f.FakeCoreV1},
		ns:          namespace,
	}
}
