package client

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"k8s.io/kubernetes/pkg/openyurt/fileCache"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	coordinationv1 "k8s.io/api/coordination/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"
	"sigs.k8s.io/yaml"
)

// Edge publish 的topic
const MqttEdgePublishRootTopic = "/lite/edge/"
const MqttCloudPublishRootTopic = "/lite/cloud/"

var mqttclient_once sync.Once
var mqttclient mqtt.Client

var subscribeMap = make(map[string]mqtt.MessageHandler)

func SubCloudLeasesOperator(client mqtt.Client, message mqtt.Message) {
	obj := &coordinationv1.Lease{}
	if err := saveMessageToObjectFile(message, obj, manifest.GetLeasesManifestPath()); err != nil {
		klog.Errorf("Save message[topic %s] payload to lease manifest error %v", message.Topic(), err)
		return
	}
	return
}

func saveMessageToObjectFile(message mqtt.Message, obj interface{}, objectManifestPath string) error {

	if err := yaml.Unmarshal(message.Payload(), obj); err != nil {
		return fmt.Errorf("unmarshal mqtt message error %v", err)
	}
	name, err := fileCache.CreateFileNameByNamespacedObject(obj)
	if err != nil {
		return fmt.Errorf("get object filename error %v", err)
	}
	filePath := filepath.Join(objectManifestPath, name)

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

/// sets the MessageHandler that will be called when a message
// is received that does not match any known subscriptions.
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	klog.V(4).Infof("Receive Message: %s from topic: %s, this message is received but does not match any known subscriptions", string(msg.Payload()), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	klog.V(4).Infof("Connected mqtt broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	klog.V(4).Infof("Connect lost:%v", err)
}

func NewMqttClient(broker string, port int, clientid, username, passwd string) mqtt.Client {

	mqttclient_once.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(fmt.Sprintf("ssl://%s:%d", broker, port))
		// 设置为false 时候， 代表持久性回话，客户端重连时候， 服务端会记录这个session , 这样之前没收到的消息可以重新发送
		opts.SetCleanSession(false)

		opts.SetClientID(clientid)
		opts.SetUsername(username)
		opts.SetPassword(passwd)

		opts.SetDefaultPublishHandler(messagePubHandler)
		opts.SetAutoReconnect(true)
		opts.SetConnectRetry(true)
		opts.SetOnConnectHandler(connectHandler)
		opts.SetConnectionLostHandler(connectLostHandler)
		opts.OnConnectionLost = connectLostHandler
		opts.SetKeepAlive(10 * time.Second)
		opts.SetConnectTimeout(20 * time.Second)
		opts.SetConnectRetryInterval(20 * time.Second)
		opts.SetPingTimeout(20 * time.Second)
		opts.SetWriteTimeout(10 * time.Second)

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			klog.Fatalf("Connect mqtt broker error %v", token.Error())
		}

		/*
			for t, f := range subscribeMap {
				klog.V(4).Infof("Prepare subscribe topic %s", t)
				token := client.Subscribe(t, 1, f)
				token.Wait()
				if err := token.Error(); err != nil {
					klog.Fatalf("Subscribe topic %s error %v", t, err)
				}
				klog.V(4).Infof("Subscribe topic %s successfully", t)
			}
		*/
		mqttclient = client
	})

	return mqttclient
}
