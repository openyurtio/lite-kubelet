package client

// KubeletOperatorInterface is interface of mqttclient
type KubeletOperatorInterface interface {
	PodsGetter
	NodesGetter
	EventsGetter
	LeasesGetter
}

type MQTTClient struct {
}

func (M MQTTClient) Pods(namespace string) PodInstance {
	return newPods(namespace)
}

func (M MQTTClient) Nodes() NodeInstance {
	return newNodes()
}

func (M MQTTClient) Events(namespace string) EventInstance {
	return newEvents(namespace)
}

func (M MQTTClient) Leases(namespace string) LeaseInstance {
	return newLeases(namespace)
}

var _ KubeletOperatorInterface = &MQTTClient{}
