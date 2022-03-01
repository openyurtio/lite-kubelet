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
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/component-helpers/apimachinery/lease"
)

type controller struct {
	client    mqtt.Client
	rootTopic string
	nodeName  string
}

// Run
// Normally kubelet updates the lease to Apiserver every 10 seconds to maintain the node ready.
// If the lease heartbeat is sent through MQTT protocol, and the number of hosts is very large, too many heartbeats will be sent every day.
// Due to the limited ability of kole-controller to consume messages, This may cause message accumulation and the host lease status cannot be updated quickly.
// Too much news also brings more economic costs.

// The Run method regularly sends heartbeat information to the kole-controller through MQTT
// and reduces the heartbeat sending time interval.
// We believe that pod eviction is not required in the lightweight scenario, so it is acceptable to delay the update of host state.
// Considering the cost and the consumption speed of messages, We temporarily set the heartbeat interval to 5 minutes.
// But the heartbeat information is not a lease object.

//The kole-controller updates the lease object every 10 seconds instead of the original kubelet lease logic after receiving the heartbeat message .
// When the kole-controller does not receive the heartbeat message from lite-kubelet within 1.5 x 5 m,
// it considers that the lite-node is offline and stops updating the lease information.

// After some time, kube-controller-manger determines whether the node is readdy based on the lease update time, This is the original logic of kube-controller-manager

// TODO
// The MQTT3 protocol supports testamentary mode (Will Message),
// in which MTQT Brokers send will messages indicating that certain lite-node are offline when they are unexpectedly disconnected, offline, or disconnected from the MQTT broker.
// In the future, if kole-controller supports the ability to subscribe to wills, lite-kubelet will not need to send heartbeat packets regularly at all,
// further reducing the number of messages to be sent.
func (c *controller) Run(stopCh <-chan struct{}) {

	/*
	ticker := time.NewTicker(time.Minute*5)
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
	 */
}

// NewController constructs and returns a controller
func NewController(mqttClient mqtt.Client, rootTopic, nodeName string) lease.Controller {
	return &controller{
		client:    mqttClient,
		rootTopic: rootTopic,
		nodeName:  nodeName,
	}
}
