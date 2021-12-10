package client

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

type Identityor interface {
	GetIdentity() string
}
type PublishData struct {
	ObjectName        string          `json:"object_name,omitempty"`
	Object            interface{}     `json:"object,omitempty"`
	Options           interface{}     `json:"options,omitempty"`
	Identity          string          `json:"identity,omitempty"`
	PatchType         types.PatchType `json:"patch_type,omitempty"`
	PatchData         []byte          `json:"patch_data,omitempty"`
	PatchSubResources []string        `json:"patch_sub_resources,omitempty"`
	NodeName          string          `json:"node_name,omitempty"`
}

func (p *PublishData) GetIdentity() string {
	return p.Identity
}

type PublishAckData struct {
	Object    interface{} `json:"object,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Identity  string      `json:"identity,omitempty"`
	ErrorType ErrorType   `json:"error_type,omitempty"`
}

func (p *PublishAckData) GetIdentity() string {
	return p.Identity
}

var _ Identityor = &PublishData{}
var _ Identityor = &PublishAckData{}

func newPublishData(nodename, objectName string, obj metav1.Object, options interface{}, pathType types.PatchType, patchData []byte, subresources []string) *PublishData {
	return &PublishData{
		ObjectName:        objectName,
		NodeName:          nodename,
		Object:            obj,
		Options:           options,
		Identity:          string(uuid.NewUUID()),
		PatchType:         pathType,
		PatchData:         patchData,
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

func PublishCreateData(nodename string, object metav1.Object, options metav1.CreateOptions) *PublishData {
	return newPublishData(nodename, object.GetName(), object, options, "", nil, nil)
}

func PublishDeleteData(nodename, name string, options metav1.DeleteOptions) *PublishData {
	return newPublishData(nodename, name, nil, options, "", nil, nil)
}

func PublishPatchData(nodename, name string, patchType types.PatchType, patchData []byte, options metav1.PatchOptions, subresources ...string) *PublishData {
	return newPublishData(nodename, name, nil, options, patchType, patchData, subresources)
}

func PublishUpdateData(nodename string, object metav1.Object, options metav1.UpdateOptions) *PublishData {
	return newPublishData(nodename, object.GetName(), object, options, "", nil, nil)
}
