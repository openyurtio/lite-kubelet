package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"
	"sigs.k8s.io/yaml"
)

// KubeletOperatorInterface is interface of mqttclient
type KubeletOperatorInterface interface {
	PodsGetter
	NodesGetter
	EventsGetter
	LeasesGetter
	MessageSendor
}

type MessageSendor interface {
	Send(topic string, qos byte, retained bool, obj Identityor, timeout time.Duration) error
}

type LocalClient struct {
	nodename string
	send     mqtt.Client
	nodes    cache.Indexer
	leases   cache.Indexer
	events   cache.Indexer
}

func (l *LocalClient) SubscribeTopics(nodename string) {

	RegisterSubtopicor(nodename, &LeaseSubTopic{})
	RegisterSubtopicor(nodename, &NodeSubTopic{})
	RegisterSubtopicor(nodename, &AckSubTopic{})

	for t, f := range GetAllTopicFuncs() {
		klog.V(4).Infof("Prepare subscribe topic %s", t)
		token := l.send.Subscribe(t, 1, f)
		token.Wait()
		if err := token.Error(); err != nil {
			klog.Fatalf("Subscribe topic %s error %v", t, err)
		}
		klog.V(4).Infof("Subscribe topic %s successfully", t)
	}

}

func (l *LocalClient) Send(topic string, qos byte, retained bool, obj Identityor, timeout time.Duration) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return fmt.Errorf("Marshal object error %v", err)
	}
	klog.V(4).Infof("###### [%s] Prepare to send to topic %s, data:\n%s", obj.GetIdentity(), topic, string(data))
	token := l.send.Publish(topic, qos, retained, data)
	out := token.WaitTimeout(timeout)
	if !out {
		return fmt.Errorf("[%s] Publish data timeout", obj.GetIdentity())
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("[%s] Publish data error %v", obj.GetIdentity(), err)
	}
	klog.V(4).Infof("###### [%s] Send to topic %s, data successful", obj.GetIdentity(), topic)
	return nil
}

func (l *LocalClient) Pods(namespace string) PodInstance {
	return newPods(l.nodename, namespace, l)
}

func (l *LocalClient) Nodes() NodeInstance {
	return newNodes(l.nodename, l.nodes, l)
}

func (l *LocalClient) Events(namespace string) EventInstance {
	return newEvents(l.nodename, namespace, l.events, l)
}

func (l *LocalClient) Leases(namespace string) LeaseInstance {
	return newLeases(l.nodename, namespace, l.leases, l)
}

func (l *LocalClient) GetNodesIndexer() cache.Indexer {
	return l.nodes
}

func NewLocalClient(nodename, broker string, port int, clientid, username, passwd string) (*LocalClient, error) {

	klog.V(4).Infof("create mqtt client  broker[%v] port[%v] clientid[%v] username[%v], passwd[%v]",
		broker, port, clientid, username, passwd)

	if len(broker) == 0 || port == 0 || len(clientid) == 0 || len(username) == 0 || len(passwd) == 0 {
		return nil, fmt.Errorf("now broker[%v] port[%v] clientid[%v] username[%v], passwd[%v], some of them is nil",
			broker, port, clientid, username, passwd)
	}

	nodeIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileNodeDeps(manifest.GetNodesManifestPath()), false, nil)
	leaseIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileLeaseDeps(manifest.GetLeasesManifestPath()), false, nil)
	eventIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileEventDeps(manifest.GetEventsManifestPath()), false, nil)
	c := NewMqttClient(broker, port, clientid, username, passwd)

	l := &LocalClient{
		nodename: nodename,
		send:     c,
		nodes:    nodeIndexer,
		leases:   leaseIndexer,
		events:   eventIndexer,
	}
	return l, nil
}

func saveMessageToObjectFile(message mqtt.Message, obj metav1.Object, objectManifestPath string) error {

	if err := yaml.Unmarshal(message.Payload(), obj); err != nil {
		return fmt.Errorf("unmarshal mqtt message error %v", err)
	}

	name, err := fileCache.CreateFileNameByNamespacedObject(obj)
	if err != nil {
		return fmt.Errorf("get object filename error %v", err)
	}
	filePath := filepath.Join(objectManifestPath, name)

	klog.V(4).Infof("Prepare to save object[topic %s] payload to localcache %s", message.Topic(), filePath)

	if obj.GetDeletionTimestamp() != nil {
		klog.Warningf("Find object[%v/%v] deletionTimestamp is not nil , need to delete localcachefile %s", obj.GetNamespace(), obj.GetName(), filePath)
		err = os.RemoveAll(filePath)
		if err != nil {
			klog.Errorf("Remove cache file %s error %v", filePath, err)
			return err
		}
		klog.Warningf("Delete localcache file %s succefully", filePath)
	}

	// must use CREATE AND TRUNC
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("openfile %s error %v", filePath, err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	if _, err := write.Write(message.Payload()); err != nil {
		return fmt.Errorf("write payload error %v", err)
	}
	return write.Flush()
}

var _ KubeletOperatorInterface = &LocalClient{}
