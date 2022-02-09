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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
	"k8s.io/kubernetes/pkg/openyurt/manifest"
	"sigs.k8s.io/yaml"
)

func subscribeHandlers(nodeName, rootTopic string) map[string]mqtt.MessageHandler {
	handlers := make(map[string]mqtt.MessageHandler)
	handlers[GetPublishTopic(rootTopic, nodeName)] = func(client mqtt.Client, mes mqtt.Message) {
		subcribeData, err := UnmarshalPayloadToSubscribeData(mes.Payload())
		if err != nil {
			klog.Errorf("Unmarshal message payload[topic %s] to subscribedata error %v.", mes.Topic(), err)
			return
		}
		if err := dealWithSubcribeDate(subcribeData); err != nil {
			klog.Errorf("dealWithSubcribeData %s error %v.", subcribeData, err)
			return
		}
	}
	return handlers
}

func dealWithSubcribeDate(subcribeData *SubscribeData) error {
	switch subcribeData.DataType {
	case SubscribeDataTypeAck:
		d := &AckData{}
		dateBytes, err := json.Marshal(subcribeData.Data)
		if err != nil {
			klog.Errorf("SubcribeData.Data Marshal byte error %v", err)
			return err
		}
		if err := json.Unmarshal(dateBytes, d); err != nil {
			klog.Errorf("SubcribeData.Data UnMarshal AckData object error %v", err)
			return err
		}
		GetDefaultTimeoutCache().Set(d)
		klog.V(2).Infof("Subscribe ack payload %s successful", d)
	case SubscribeDataTypeNode:
		d := &corev1.Node{}
		dateBytes, err := yaml.Marshal(subcribeData.Data)
		if err != nil {
			klog.Errorf("SubcribeData.Data Marshal byte error %v", err)
			return err
		}
		if err := yaml.Unmarshal(dateBytes, d); err != nil {
			klog.Errorf("SubcribeData.Data UnMarshal Node object error %v", err)
			return err
		}
		if err := saveToObjectFile(dateBytes, d, manifest.GetNodesManifestPath()); err != nil {
			klog.Errorf("Save node object to file error %v", err)
			return err
		}
		klog.V(2).Infof("Subscribe node payload %s successful", d)
	case SubscribeDataTypePod:
		d := &corev1.Pod{}
		dateBytes, err := yaml.Marshal(subcribeData.Data)
		if err != nil {
			klog.Errorf("SubcribeData.Data Marshal byte error %v", err)
			return err
		}
		if err := yaml.Unmarshal(dateBytes, d); err != nil {
			klog.Errorf("SubcribeData.Data UnMarshal Node object error %v", err)
			return err
		}
		if err := saveToObjectFile(dateBytes, d, manifest.GetPodsManifestPath()); err != nil {
			klog.Errorf("Save pod object to file error %v", err)
			return err
		}
		klog.V(2).Infof("Subscribe pod payload %s successful", d)
	default:
		return fmt.Errorf("wrong subscribedata type %s", subcribeData.DataType)
	}
	return nil
}

func saveToObjectFile(payload []byte, obj metav1.Object, objectManifestPath string) error {

	name, err := fileCache.CreateFileNameByNamespacedObject(obj)
	if err != nil {
		return fmt.Errorf("get object filename error %v", err)
	}
	filePath := filepath.Join(objectManifestPath, name)

	if obj.GetDeletionTimestamp() != nil {
		klog.Warningf("Find object[%v/%v] deletionTimestamp is not nil , need to delete localcachefile %s", obj.GetNamespace(), obj.GetName(), filePath)
		err = os.RemoveAll(filePath)
		if err != nil {
			klog.Errorf("Remove cache file %s error %v", filePath, err)
			return err
		}
		klog.Warningf("Delete localcache file %s succefully", filePath)
		return nil
	}

	// must use CREATE AND TRUNC
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("openfile %s error %v", filePath, err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	if _, err := write.Write(payload); err != nil {
		return fmt.Errorf("write payload error %v", err)
	}
	return write.Flush()
}
