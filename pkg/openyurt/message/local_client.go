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
	"os"
	"path/filepath"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
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
	if err := Mqtt3Send(l.send, GetDataTopic(l.rootTopic), obj); err != nil {
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

	nodeIndexer := fileCache.NewFileObiectIndexer(fileCache.NewDefaultFileNodeDeps(), false, nil)

	clientID := fmt.Sprintf("%s@@@%s", group, nodename)
	username := fmt.Sprintf("Signature|%s|%s", accessKey, instance)
	passwd := GetSignature(clientID, secretKey)

	//c := NewMqttClient(broker, port, clientID, username, passwd)
	subHandlers := subscribeHandlers(nodename, rootTopic)

	cleanLocalCache := func(client mqtt.Client, nodename, rootTopic string) {
		// try 3 times
		klog.Infof("Prepare to clean local useless caches")
		defer klog.Infof("Clean local useless caches end ...")

		for i := 0; i < 3; i++ {
			data := PublishStartData(nodename)
			if err := Mqtt3Send(client, GetDataTopic(rootTopic), data); err != nil {
				klog.Errorf("Send start data error when mqtt client connected. %v", err)
				continue
			}
			ackdata, ok := GetDefaultTimeoutCache().Pop(data.Identity, time.Second*5)
			if !ok {
				klog.Errorf("Get ack data[%s] from timeoutCache timeout  when get start data", data.Identity)
				continue
			}

			startdata := &AckDataStartObject{}
			_, err := ackdata.UnmarshalAckData(startdata)
			if err != nil {
				klog.Errorf("ack data unmarshal error %v,data:\n%v", err, *ackdata)
				continue
			}

			cloudMap := make(map[string]struct{})
			for i, _ := range startdata.SecretList {
				s := startdata.SecretList[i]
				if err := saveSecretToObjectFile(s); err != nil {
					klog.Error("Save secret[%s/%s] to object file error %v", s.GetNamespace(), s.GetName(), err)
					continue
				}
			}
			for i, _ := range startdata.PodsList {
				tmpPod := startdata.PodsList[i]
				filename, err := fileCache.CreateFileNameByObject(tmpPod)
				if err != nil {
					klog.Errorf("create filename by pod[%s/%s] error %v", tmpPod.GetNamespace(), tmpPod.GetName(), err)
					continue
				}
				cloudMap[filepath.Base(filename)] = struct{}{}

				if err := savePodToObjectFile(tmpPod); err != nil {
					klog.Error("Save pod[%s/%s] to object file error %v", tmpPod.GetNamespace(), tmpPod.GetName(), err)
					continue
				}
			}

			podDeps := fileCache.NewDefaultFilePodDeps()
			for _, f := range podDeps.GetAllFiles() {
				if _, ok := cloudMap[filepath.Base(f)]; !ok {
					klog.Warningf("The pod corresponding to localcache file %s does not exist in the cloud, so need to delete the localcache file", f)

					_, err := os.Stat(f)
					if err != nil {
						klog.Errorf("Can't get stat for %q: %v", f, err)
						continue
					}
					err = os.RemoveAll(f)
					if err != nil {
						klog.Errorf("Remove cache file %s error %v", f, err)
					} else {
						klog.Infof("Delete localcache file %s succefully", f)
					}
				} else {
					klog.Infof("The pod corresponding to localcache file %s exist in the cloud, so keep the localcache file, do nothing", f)
				}
			}
			return
		}
	}

	c := newMqtt3Client(broker, port,
		clientID, username, passwd,
		true, false,
		func(client mqtt.Client) {
			// subscribe all topic
			for t, f := range subHandlers {
				token := client.Subscribe(t, 1, f)
				token.Wait()
				if err := token.Error(); err != nil {
					klog.Fatalf("Subscribe topic %s error %v", t, err)
				}
				klog.Infof("Subscribe topic %s successfully", t)
			}
			klog.Infof("Subscribe all topic successfully")
			// publish register info
			cleanLocalCache(client, nodename, rootTopic)

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
