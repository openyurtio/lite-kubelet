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
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/openyurt/fileCache"
	"sigs.k8s.io/yaml"
)

func subscribeHandlers(nodeName, rootTopic string) map[string]mqtt.MessageHandler {
	handlers := make(map[string]mqtt.MessageHandler)
	handlers[GetCtlTopic(rootTopic, nodeName)] = func(client mqtt.Client, mes mqtt.Message) {
		subcribeData, err := UnmarshalPayloadToSubscribeData(mes.Payload())
		if err != nil {
			klog.Errorf("Unmarshal message payload[topic %s] to subscribedata error %v.", mes.Topic(), err)
			return
		}
		if err := dealWithSubcribeDate(client, rootTopic, nodeName, subcribeData); err != nil {
			klog.Errorf("dealWithSubcribeData %s error %v.", subcribeData, err)
			return
		}
	}
	return handlers
}

func dealWithSubcribeDate(client mqtt.Client, rootTopic, nodeName string, subcribeData *SubscribeData) error {
	switch subcribeData.DataType {
	case SubscribeDataTypeAck:
		d := subcribeData.AckData
		GetDefaultTimeoutCache().Set(d)
		klog.V(2).Infof("Subscribe ack payload %s successful", d)
	case SubscribeDataTypeNode:
		if err := saveNodeToObjectFile(subcribeData.NodeData); err != nil {
			klog.Errorf("Save node object to file error %v", err)
			return err
		}
		klog.V(4).Infof("Subscribe node payload %s to localcache successful", subcribeData.NodeData.Name)
	case SubscribeDataTypePod:
		for i, _ := range subcribeData.SecretsData {
			s := subcribeData.SecretsData[i]
			if s == nil {
				continue
			}
			if err := saveSecretToObjectFile(s); err != nil {
				klog.Error("save secret to objectfile error %v", err)
				continue
			}
		}

		if err := savePodToObjectFile(subcribeData.PodData); err != nil {
			klog.Errorf("Save pod object to file error %v", err)
			return err
		}

		klog.V(2).Infof("Subscribe pod payload [%s/%s] to localcache successful", subcribeData.PodData.GetNamespace(), subcribeData.PodData.GetName())
		/*
			case SubscribeDataTypeRequestHeartBeat:
				if err := Mqtt3Send(client, GetDataTopic(rootTopic), PublishOnlineData(nodeName)); err != nil {
					klog.Errorf("Send heartbeat online data error %v", err)
					return err
				}
				klog.V(2).Infof("Subscribe request heartbeat payload successful, and send online heartbeat data successful")
		*/
	default:
		return fmt.Errorf("wrong subscribedata type %s", subcribeData.DataType)
	}
	return nil
}

func saveNodeToObjectFile(s *corev1.Node) error {

	dateBytes, err := yaml.Marshal(s)
	if err != nil {
		klog.Errorf("Marshal node[%s/%s] to bytes error %v",
			s.GetNamespace(), s.GetName(), err)
		return err
	}
	path, err := fileCache.NewDefaultFileNodeDeps().GetFullFileName(s)
	if err != nil {
		klog.Errorf("Get node object %++v full file name error %v", *s, err)
		return err
	}
	if err := saveToObjectFile(dateBytes, path); err != nil {
		klog.Errorf("Save node object[%s] to file error %v", s.GetName(), err)
		return err
	}
	return nil
}

func savePodToObjectFile(s *corev1.Pod) error {

	dateBytes, err := yaml.Marshal(s)
	if err != nil {
		klog.Errorf("Marshal pod[%s/%s] to bytes error %v",
			s.GetNamespace(), s.GetName(), err)
		return err
	}
	path, err := fileCache.NewDefaultFilePodDeps().GetFullFileName(s)
	if err != nil {
		klog.Errorf("Get pod object %++v full file name error %v", *s, err)
		return err
	}
	if err := saveToObjectFile(dateBytes, path); err != nil {
		klog.Errorf("Save pod object[%s/%s] to file error %v", s.GetNamespace(), s.GetName(), err)
		return err
	}
	return nil
}

func saveSecretToObjectFile(s *corev1.Secret) error {

	dateBytes, err := yaml.Marshal(s)
	if err != nil {
		klog.Errorf("Marshal secret[%s/%s] to bytes error %v",
			s.GetNamespace(), s.GetName(), err)
		return err
	}
	path, err := fileCache.NewDefaultFileSecretDeps().GetFullFileName(s)
	if err != nil {
		klog.Errorf("Get secret object %++v full file name error %v", *s, err)
		return err
	}
	if err := saveToObjectFile(dateBytes, path); err != nil {
		klog.Errorf("Save secrect object[%s/%s] to file error %v", s.GetNamespace(), s.GetName(), err)
		return err
	}
	return nil
}

func saveToObjectFile(payload []byte, filePath string) error {

	/*
		name, err := fileCache.CreateFileNameByObject(obj)
		if err != nil {
			return fmt.Errorf("get object filename error %v", err)
		}
		filePath := filepath.Join(objectManifestPath, name)
	*/

	/*
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
	*/

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
