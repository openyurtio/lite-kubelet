package client

import (
	"path/filepath"

	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	coordinationv1 "k8s.io/api/coordination/v1"
	"k8s.io/klog/v2"
)

type LeaseSubTopic struct {
}

func (l *LeaseSubTopic) GetResourceTopic(nodename string) (string, mqtt.MessageHandler) {
	// lease tpoic : /lite/cloud/leases/{nodename}
	return filepath.Join(MqttCloudPublishRootTopic, "leases", nodename),
		func(client mqtt.Client, message mqtt.Message) {
			if err := saveMessageToObjectFile(message, &coordinationv1.Lease{}, manifest.GetLeasesManifestPath()); err != nil {
				klog.Errorf("Save message[topic %s] payload to lease manifest error %v", message.Topic(), err)
			}
		}
}

var _ SubTopicor = &LeaseSubTopic{}
