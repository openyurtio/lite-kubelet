package clientSet

import (
	clientset "k8s.io/client-go/kubernetes"
	fakekube "k8s.io/client-go/kubernetes/fake"
	coordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	fakecoordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	fakecorev1 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	openyurtcoordinationv1 "k8s.io/kubernetes/pkg/openyurt/clientSet/typed/coordination/v1"
	openyurtcorev1 "k8s.io/kubernetes/pkg/openyurt/clientSet/typed/core/v1"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/client"
)

// Clientset implement clientset.Interface clientset "k8s.io/client-go/kubernetes"
type Clientset struct {
	*fakekube.Clientset
	MQTTClient client.KubeletOperatorInterface
}

// NewSimpleClientset is new clientset.Interface by mqtt
func NewSimpleClientset(metaClient client.KubeletOperatorInterface) clientset.Interface {
	return &Clientset{
		Clientset:  fakekube.NewSimpleClientset(),
		MQTTClient: metaClient}
}

func (c *Clientset) CoreV1() corev1.CoreV1Interface {
	return &openyurtcorev1.FakeCoreV1{FakeCoreV1: fakecorev1.FakeCoreV1{Fake: &c.Fake}, MQTTClient: c.MQTTClient}
}

func (c *Clientset) CoordinationV1() coordinationv1.CoordinationV1Interface {
	return &openyurtcoordinationv1.FakeCoordinationV1{FakeCoordinationV1: fakecoordinationv1.FakeCoordinationV1{Fake: &c.Fake}, MQTTClient: c.MQTTClient}
}
