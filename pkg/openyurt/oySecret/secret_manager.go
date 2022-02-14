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
package oySecret

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/secret"
)

type SecretManager struct {
	index cache.Indexer
}

func NewSecretManager(index cache.Indexer) secret.Manager {
	return &SecretManager{
		index: index,
	}
}

func (s *SecretManager) GetSecret(namespace, name string) (*corev1.Secret, error) {
	obj, exists, err := s.index.Get(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	})
	if !exists {
		klog.Errorf("Can not get secret[%s/%s]", namespace, name)
		return nil, fmt.Errorf("can not get secret")
	}

	finnal, ok := obj.(*corev1.Secret)
	if !ok {
		klog.Errorf("Cache obj convert to *corev1.Secret error", err)
		return nil, fmt.Errorf("cache obj convert to *corev1.Node error")
	}
	return finnal, nil
}

func (s *SecretManager) RegisterPod(pod *corev1.Pod) {
	klog.V(4).Infof("implement me RegisterPod")
	return
}

func (s *SecretManager) UnregisterPod(pod *corev1.Pod) {
	klog.V(4).Infof("implement me UnRegisterPod")
	return
}

var _ secret.Manager = &SecretManager{}
