package client

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"
)

type NodeSubTopic struct {
}

func (l *NodeSubTopic) GetResourceTopic(nodename string) (string, mqtt.MessageHandler) {
	// lease tpoic : /lite/cloud/nodes/{nodename}
	return filepath.Join(MqttCloudPublishRootTopic, "nodes", nodename),
		func(client mqtt.Client, message mqtt.Message) {
			if err := saveMessageToObjectFile(message, &corev1.Node{}, manifest.GetNodesManifestPath()); err != nil {
				klog.Errorf("Save message[topic %s] payload to nodes manifest error %v", message.Topic(), err)
			}
		}
}

var _ SubTopicor = &NodeSubTopic{}
