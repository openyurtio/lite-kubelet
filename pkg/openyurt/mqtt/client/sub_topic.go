package client

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"k8s.io/klog/v2"
)

var allSubTopicFuncs map[string]mqtt.MessageHandler

type SubTopicor interface {
	GetResourceTopic(nodename string) (string, mqtt.MessageHandler)
}

func init() {
	allSubTopicFuncs = make(map[string]mqtt.MessageHandler)
}

func GetAllTopicFuncs() map[string]mqtt.MessageHandler {
	return allSubTopicFuncs
}

func RegisterSubtopicor(nodename string, s SubTopicor) {
	if t, f := s.GetResourceTopic(nodename); len(t) != 0 {
		if _, exist := allSubTopicFuncs[t]; exist {
			klog.Fatalf("Has seem subtopic %s", t)
		}
		allSubTopicFuncs[t] = f
	}
}
