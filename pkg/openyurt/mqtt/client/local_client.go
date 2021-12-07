package client

import (
	"fmt"
	"time"

	"k8s.io/klog/v2"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"
	"sigs.k8s.io/yaml"
)

type Topicor interface {
	GetPreTopic() string
}

// KubeletOperatorInterface is interface of mqttclient
type KubeletOperatorInterface interface {
	PodsGetter
	NodesGetter
	EventsGetter
	LeasesGetter
	MessageSendor
}

type MessageSendor interface {
	Send(topic string, qos byte, retained bool, obj interface{}, timeout time.Duration) error
}

type LocalClient struct {
	send   mqtt.Client
	nodes  cache.Indexer
	leases cache.Indexer
	events cache.Indexer
	//pods   cache.Indexer
}

func (l *LocalClient) Send(topic string, qos byte, retained bool, obj interface{}, timeout time.Duration) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("Marshal object error %v", err)
	}
	token := l.send.Publish(topic, qos, retained, data)
	out := token.WaitTimeout(timeout)
	if !out {
		return fmt.Errorf("Publish data timeout")
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("Publish data error %v", err)
	}
	return nil
}

func (l *LocalClient) Pods(namespace string) PodInstance {
	return newPods(namespace, l)
}

func (l *LocalClient) Nodes() NodeInstance {
	return newNodes(l.nodes, l)
}

func (l *LocalClient) Events(namespace string) EventInstance {
	return newEvents(namespace, l.events, l)
}

func (l *LocalClient) Leases(namespace string) LeaseInstance {
	return newLeases(namespace, l.leases, l)
}

func (l *LocalClient) GetNodesIndexer() cache.Indexer {
	return l.nodes
}

func NewLocalClient(broker string, port int, clientid, username, passwd string) (*LocalClient, error) {

	klog.V(4).Infof("create mqtt client  broker[%v] port[%v] clientid[%v] username[%v], passwd[%v]",
		broker, port, clientid, username, passwd)

	if len(broker) == 0 || port == 0 || len(clientid) == 0 || len(username) == 0 || len(passwd) == 0 {
		return nil, fmt.Errorf("now broker[%v] port[%v] clientid[%v] username[%v], passwd[%v], some of them is nil",
			broker, port, clientid, username, passwd)
	}

	nodeIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileNodeDeps(manifest.GetNodesManifestPath()), false, nil)
	leaseIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileLeaseDeps(manifest.GetLeasesManifestPath()), false, nil)
	//podIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFilePodDeps(manifest.GetPodsManifestPath(mqttManifestDir)), false, nil)
	eventIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileEventDeps(manifest.GetEventsManifestPath()), false, nil)
	c := NewMqttClient(broker, port, clientid, username, passwd)

	l := &LocalClient{
		send:   c,
		nodes:  nodeIndexer,
		leases: leaseIndexer,
		//pods:   podIndexer,
		events: eventIndexer,
	}
	return l, nil
}

var _ KubeletOperatorInterface = &LocalClient{}
