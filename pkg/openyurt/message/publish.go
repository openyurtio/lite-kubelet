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
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

type ErrorType int

const (
	ErrorTypeStatusError ErrorType = iota
	ErrorTypeStringError
	ErrorTypeNil
)

const (
	ObjectTypeEvent = "event"
	ObjectTypeLease = "lease"
	ObjectTypeNode  = "node"
	ObjectTypePod   = "pod"
	ObjectTypeStart = "start"
)

const (
	OperateTypeCreate = "create"
	OperateTypeUpdate = "update"
	OperateTypePatch  = "patch"
	OperateTypeDelete = "delete"
	OperateTypeGet    = "get"
	OperateTypeStart  = "start"
)

type PublishData struct {
	// object type, event, pod, node, lease
	ObjectType string `json:"object_type,omitempty"`
	// operate type create, update, get , patch, delete
	OperateType string `json:"operate_type,omitempty"`

	ObjectName        string          `json:"object_name,omitempty"`
	ObjectNS          string          `json:"object_ns,omitempty"`
	Object            interface{}     `json:"object,omitempty"`
	Options           interface{}     `json:"options,omitempty"`
	Identity          string          `json:"identity,omitempty"`
	PatchType         types.PatchType `json:"patch_type,omitempty"`
	PatchData         []byte          `json:"patch_data,omitempty"`
	PatchSubResources []string        `json:"patch_sub_resources,omitempty"`
	NodeName          string          `json:"node_name,omitempty"`
	NeedAck           bool            `json:"need_ack,omitempty"`
}

func (p *PublishData) String() string {
	return fmt.Sprintf("%s %s [%s/%s] identity [%s]", p.OperateType, p.ObjectType,
		p.ObjectNS, p.ObjectName, p.Identity)
}

var _ fmt.Stringer = &PublishData{}

func newPublishData(
	objectType, operateType string,
	needAck bool,
	nodename string,
	obj metav1.Object,
	options interface{},
	pathType types.PatchType,
	patchData []byte,
	subresources []string) *PublishData {
	var name, ns string

	if obj != nil {
		name = obj.GetName()
		ns = obj.GetNamespace()
	}
	data := &PublishData{
		ObjectType:        objectType,
		OperateType:       operateType,
		ObjectName:        name,
		ObjectNS:          ns,
		NodeName:          nodename,
		Object:            obj,
		Options:           options,
		Identity:          string(uuid.NewUUID()),
		PatchType:         pathType,
		PatchData:         patchData,
		PatchSubResources: subresources,
		NeedAck:           needAck,
	}
	return data
}

// return value: errInfo , err
func (data *AckData) UnmarshalAckData(k8sobj interface{}) (error, error) {
	if k8sobj != nil {
		objData, err := yaml.Marshal(data.Object)
		if err != nil {
			return nil, fmt.Errorf("marshal error %v", err)
		}
		if err := yaml.Unmarshal(objData, k8sobj); err != nil {
			return nil, err
		}
	}

	if data.ErrorType == ErrorTypeNil {
		return nil, nil
	}

	edata, err := yaml.Marshal(data.Error)
	if err != nil {
		return nil, err
	}

	switch data.ErrorType {
	case ErrorTypeStringError:
		errInfo := fmt.Errorf("%v", edata)
		return errInfo, nil
	case ErrorTypeStatusError:
		s := &errors.StatusError{}
		if err := yaml.Unmarshal(edata, s); err != nil {
			return nil, err
		}
		return s, nil
	default:
		klog.Errorf("Error AckData errorType %v", data.ErrorType)
		return nil, fmt.Errorf("error AckData errortype %v", data.ErrorType)
	}
}

func PublishGetData(objectType string, needack bool, nodename string, object metav1.Object, opt metav1.GetOptions) *PublishData {
	return newPublishData(objectType, OperateTypeGet, needack, nodename, object, opt, "", nil, nil)
}

func PublishCreateData(objectType string, needack bool, nodename string, object metav1.Object, options metav1.CreateOptions) *PublishData {
	return newPublishData(objectType, OperateTypeCreate, needack, nodename, object, options, "", nil, nil)
}

func PublishDeleteData(objectType string, needack bool, nodename string, object metav1.Object, options metav1.DeleteOptions) *PublishData {
	return newPublishData(objectType, OperateTypeDelete, needack, nodename, object, options, "", nil, nil)
}

func PublishPatchData(objectType string, needack bool, nodename string, object metav1.Object, patchType types.PatchType, patchData []byte, options metav1.PatchOptions, subresources ...string) *PublishData {
	return newPublishData(objectType, OperateTypePatch, needack, nodename, object, options, patchType, patchData, subresources)
}

func PublishUpdateData(objectType string, needack bool, nodename string, object metav1.Object, options metav1.UpdateOptions) *PublishData {
	return newPublishData(objectType, OperateTypeUpdate, needack, nodename, object, options, "", nil, nil)
}

func PublishStartData(nodename string) *PublishData {
	return newPublishData(ObjectTypeStart, OperateTypeStart, true, nodename, nil, nil, "", nil, nil)
}

