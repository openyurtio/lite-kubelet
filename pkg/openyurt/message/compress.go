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
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"k8s.io/klog/v2"
)

type MessageThin interface {
	Compress(data []byte) ([]byte, error)
	UnCompress(data []byte) ([]byte, error)
}

//Gzip Compress/EnCompress
type Gzip struct {
}

func (g *Gzip) Compress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		klog.Warning("Gzip Compress fail:", err)
		return nil, err
	}
	// We should close the writer immediately instead of using defer.
	if err = writer.Close(); err != nil {
		klog.Warning("Close the Gzip Object fail:", err)
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (g *Gzip) UnCompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = reader.Close(); err != nil {
			klog.Errorf("Gzip EnCompress fail:", err.Error())
		}
	}()
	return ioutil.ReadAll(reader)
}

var _ MessageThin = &Gzip{}
