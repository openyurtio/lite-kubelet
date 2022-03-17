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
	"time"

	"sigs.k8s.io/yaml"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
)

type WillOptions struct {
	topic    string
	payload  string
	qos      byte
	retained bool
}

var defaultConnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	klog.V(5).Infof("Connected mqtt broker ...")
}

var defaultConnectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	klog.V(5).Infof("Connect lost:%v", err)
}

func newSessionMqtt3Client(broker string, port int, clientid, username, passwd string) mqtt.Client {
	return newMqtt3Client(broker, port, clientid, username, passwd, false, true, nil, defaultConnectHandler, defaultConnectLostHandler)
}

func newMqtt3Client(
	broker string,
	port int,
	clientid,
	username,
	passwd string,
	cleanSession bool,
	order bool,
	will *WillOptions,
	connectHandler mqtt.OnConnectHandler,
	connectLostHandler mqtt.ConnectionLostHandler) mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
	opts.SetCleanSession(cleanSession)

	if will != nil {
		opts.SetWill(will.topic, will.payload, will.qos, will.retained)
	}

	opts.SetClientID(clientid)
	opts.SetUsername(username)
	opts.SetPassword(passwd)
	opts.SetOrderMatters(order)
	// 设置重新使用resumesub
	opts.SetResumeSubs(true)
	// Do not set default publishHandler
	// opts.SetDefaultPublishHandler(messagePubHandler)
	//opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	if connectHandler != nil {
		opts.SetOnConnectHandler(connectHandler)
	}
	if connectLostHandler != nil {
		opts.SetConnectionLostHandler(connectLostHandler)
	}
	opts.SetKeepAlive(5 * time.Second)
	//opts.SetConnectTimeout(30 * time.Second)
	//opts.SetConnectRetryInterval(10 * time.Second)
	//opts.SetMaxReconnectInterval(10 * time.Minute)
	//opts.SetPingTimeout(10 * time.Second)
	//opts.SetWriteTimeout(10 * time.Second)
	//opts.SetReconnectingHandler()

	client := mqtt.NewClient(opts)

	for {
		klog.V(5).Infof("%s prepare to connect mqtt ...", clientid)

		token := client.Connect()
		// done
		if token.WaitTimeout(time.Second * 5) {
			if token.Error() != nil {
				klog.Errorf("Client %s connect mqtt broker error %v", clientid, token.Error())
				time.Sleep(time.Second)
				continue
			} else {
				klog.V(5).Infof("Client %s connect mqtt success...", clientid)
				break
			}
		} else {
			// timeout
			klog.Errorf("Client %s connect mqtt broker timeout, continue connect...", clientid)
			time.Sleep(time.Second)
			continue
		}
	}
	klog.Infof("Client %s connect mqtt broker successful", clientid)
	return client
}

func Mqtt3Send(c mqtt.Client, topic string, object *PublishData, thin MessageThin) error {
	opts := c.OptionsReader()
	clientID := opts.ClientID()

	data, err := yaml.Marshal(object)
	if err != nil {
		klog.Errorf("Mqtt3Send topic %s marshal error %v", topic, err)
		return err
	}

	thinData, err := thin.Compress(data)
	if err != nil {
		klog.Errorf("Mqtt3Send topic %s data ,compress error %v", topic, err)
		return err
	}

	token := c.Publish(topic, 1, false, thinData)
	if token.WaitTimeout(time.Second * 3) {
		if err := token.Error(); err != nil {
			klog.Errorf("%s publish topic[%s] data error %v", clientID, topic, err)
			return err
		}
	} else {
		// timeout
		klog.Errorf("%s publish topic[%s] data timeout", clientID, topic)
		return fmt.Errorf("%s publish topic[%s] data timeout", clientID, topic)
	}

	return nil
}
