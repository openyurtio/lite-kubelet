package client

import (
	"path/filepath"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
)

type AckSubTopic struct {
}

func (l *AckSubTopic) GetResourceTopic(nodename string) (string, mqtt.MessageHandler) {
	// ack tpoic : /lite/ack/
	return filepath.Join(MqttCloudAckRootTopic, nodename),
		func(client mqtt.Client, message mqtt.Message) {
			klog.V(4).Infof("Get ack payload:\n%s", string(message.Payload()))
			ack, err := UnmarshalPayloadToPublishAckData(message.Payload())
			if err != nil {
				klog.Errorf("Unmarshal ack payload data error %v", err)
			}
			klog.V(4).Infof("Get ack payload[%s] successfull", ack.Identity)
			GetDefaultTimeoutCache().Set(ack)
		}
}

var _ SubTopicor = &AckSubTopic{}
