package client

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/mqtt/manifest"
)

type PodSubTopic struct {
}

func (l *PodSubTopic) GetResourceTopic(nodename string) (string, mqtt.MessageHandler) {
	// lease tpoic : /lite/cloud/pods/{nodename}
	return filepath.Join(MqttCloudPublishRootTopic, "pods", nodename),
		func(client mqtt.Client, message mqtt.Message) {
			if err := saveMessageToObjectFile(message, &corev1.Pod{}, manifest.GetPodsManifestPath()); err != nil {
				klog.Errorf("Save message[topic %s] payload to lease manifest error %v", message.Topic(), err)
			}
		}
}

var _ SubTopicor = &PodSubTopic{}
