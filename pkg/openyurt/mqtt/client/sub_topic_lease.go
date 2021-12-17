package client

/*
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

*/
