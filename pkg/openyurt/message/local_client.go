/*
Copyright 2022 The OpenYurt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package message

import (
	"fmt"
	"path/filepath"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
	"k8s.io/kubernetes/pkg/openyurt/manifest"
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
	Send(data *PublishData) error
}

type LocalClient struct {
	nodename  string
	rootTopic string
	send      mqtt.Client
	nodes     cache.Indexer
	leases    cache.Indexer
	pods      cache.Indexer
	events    cache.Indexer
}

func (l *LocalClient) Send(obj *PublishData) error {
	if err := Mqtt3Send(l.send, filepath.Join(l.rootTopic, "edge"), obj); err != nil {
		klog.Errorf("Publish %s error", obj)
	}
	klog.V(4).Infof("Publish %s successful", obj)
	return nil
}

func (l *LocalClient) Pods(namespace string) PodInstance {
	return newPods(l.nodename, namespace, l.pods, l)
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

func NewLocalClient(nodename, broker string, port int,
	accessKey, secretKey, instance, group, rootTopic string) (*LocalClient, error) {

	nodeIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileNodeDeps(manifest.GetNodesManifestPath()), false, nil)

	clientID := fmt.Sprintf("%s@@@%s", group, nodename)
	username := fmt.Sprintf("Signature|%s|%s", accessKey, instance)
	passwd := GetSignature(clientID, secretKey)

	//c := NewMqttClient(broker, port, clientID, username, passwd)
	subHandlers := subscribeHandlers(nodename, rootTopic)

	c := newMqtt3Client(broker, port,
		clientID, username, passwd,
		true, false,
		func(client mqtt.Client) {
			for t, f := range subHandlers {
				token := client.Subscribe(t, 1, f)
				token.Wait()
				if err := token.Error(); err != nil {
					klog.Fatalf("Subscribe topic %s error %v", t, err)
				}
				klog.Infof("Subscribe topic %s successfully", t)
			}
		}, func(client mqtt.Client, err error) {
			klog.Warningf("mqtt client connect lost:%v", err)
		})

	l := &LocalClient{
		nodename:  nodename,
		rootTopic: rootTopic,
		send:      c,
		nodes:     nodeIndexer,
		leases:    nil,
		events:    nil,
		pods:      nil,
	}
	return l, nil
}

var _ KubeletOperatorInterface = &LocalClient{}
