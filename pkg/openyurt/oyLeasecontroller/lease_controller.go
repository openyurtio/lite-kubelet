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
package oyLeasecontroller

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/component-helpers/apimachinery/lease"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/message"
)

type controller struct {
	client    mqtt.Client
	rootTopic string
	nodeName  string
}

func (c *controller) Run(stopCh <-chan struct{}) {
	ticker := time.NewTicker(time.Minute)
	defer func() {
		ticker.Stop()
		klog.Errorf("Heartbeat controller stoped ...")
	}()
	pbData := message.PublishOnlineData(c.nodeName)

	if err := message.Mqtt3Send(c.client, message.GetDataTopic(c.rootTopic),
		pbData); err != nil {
		klog.Errorf("First send heartbeat online data %s error %v", pbData, err)
		// no not return
	} else {
		klog.V(4).Infof("First send online heartbeat data %s successful", pbData)
	}

	for {
		// Wait for next probe tick.
		select {
		case <-ticker.C:
			if err := message.Mqtt3Send(c.client, message.GetDataTopic(c.rootTopic),
				pbData); err != nil {
				klog.Errorf("Send heartbeat online data %s error %v", pbData, err)
				// break this select, continue next for loop
				break
			}
			klog.V(4).Infof("Send online heartbeat data %s successful", pbData)
		case <-stopCh:
			return
		}
	}

}

// NewController constructs and returns a controller
func NewController(mqttClient mqtt.Client, rootTopic, nodeName string) lease.Controller {
	return &controller{
		client:    mqttClient,
		rootTopic: rootTopic,
		nodeName:  nodeName,
	}
}
