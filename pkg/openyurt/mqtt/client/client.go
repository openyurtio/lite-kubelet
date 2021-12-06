package client

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	klog.V(4).Infof("Publish Message: %s from topic: %s", string(msg.Payload()), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	klog.V(4).Infof("Connected mqtt broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	klog.V(4).Infof("Connect lost:%v", err)
}

func NewMqttClient(broker string, port int, clientid, username, passwd string) mqtt.Client {

	opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("wss://%s:%d", broker, port))
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
		fmt.Printf("connect error %v\n", token.Error())
		panic(token.Error())
	} else {
		fmt.Printf("connectd success")
	}
}
