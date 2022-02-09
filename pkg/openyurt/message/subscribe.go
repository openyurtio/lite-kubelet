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

	"sigs.k8s.io/yaml"
)

const (
	SubscribeDataTypeAck  = "ack"
	SubscribeDataTypeNode = "node"
	SubscribeDataTypePod  = "pod"
)

type SubscribeData struct {
	DataType string      `json:"data_type,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type AckData struct {
	Object    interface{} `json:"object,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Identity  string      `json:"identity,omitempty"`
	ErrorType ErrorType   `json:"error_type,omitempty"`
}

func (p *AckData) String() string {
	return fmt.Sprintf("identity %s", p.Identity)
}

func (s *SubscribeData) String() string {
	return fmt.Sprintf("datetype %s", s.DataType)
}

func UnmarshalPayloadToSubscribeData(payload []byte) (*SubscribeData, error) {
	d := &SubscribeData{}
	if err := yaml.Unmarshal(payload, d); err != nil {
		return nil, err
	}
	return d, nil
}

var _ fmt.Stringer = &SubscribeData{}
var _ fmt.Stringer = &AckData{}
