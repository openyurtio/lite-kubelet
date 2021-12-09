package client

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

type PublishData struct {
	Object            interface{}     `json:"object,omitempty"`
	Options           interface{}     `json:"options,omitempty"`
	Identity          string          `json:"identity,omitempty"`
	PatchType         types.PatchType `json:"patch_type,omitempty"`
	PatchSubResources []string        `json:"patch_sub_resources,omitempty"`
	NodeName          string          `json:"node_name,omitempty"`
}

type PublishAckData struct {
	Object    interface{} `json:"object,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Identity  string      `json:"identity,omitempty"`
	ErrorType ErrorType   `json:"error_type,omitempty"`
}

func newPublishData(nodename string, obj runtime.Object, options interface{}, pathType types.PatchType, subresources []string) *PublishData {
	return &PublishData{
		NodeName:          nodename,
		Object:            obj,
		Options:           options,
		Identity:          string(uuid.NewUUID()),
		PatchType:         pathType,
		PatchSubResources: subresources,
	}
}

func UnmarshalPayloadToPublishAckData(payload []byte) (*PublishAckData, error) {
	d := &PublishAckData{}
	if err := yaml.Unmarshal(payload, d); err != nil {
		return nil, err
	}
	return d, nil
}

// return value: errInfo , err
func (data *PublishAckData) UnmarshalPublishAckData(k8sobj interface{}) (error, error) {
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
		klog.Errorf("Error publishAckData errorType %v", data.ErrorType)
		return nil, fmt.Errorf("error publishAckData errortype %v", data.ErrorType)
	}
}

func PublishCreateData(nodename string, object runtime.Object, options metav1.CreateOptions) *PublishData {
	return newPublishData(nodename, object, options, "", nil)
}

func PublishDeleteData(nodename, name string, options metav1.DeleteOptions) *PublishData {
	panic("implement me")
}

func PublishPatchData(nodename, name string, patchType types.PatchType, bytes []byte, options metav1.PatchOptions, s2 ...string) *PublishData {
	panic("implement me")
}

func PublishUpdateData(nodename string, object runtime.Object, options metav1.UpdateOptions) *PublishData {
	panic("implement me")
}
