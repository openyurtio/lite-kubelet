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
package oyProber

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/prober"
	"k8s.io/kubernetes/pkg/kubelet/prober/results"
	"k8s.io/kubernetes/pkg/kubelet/status"
)

type manager struct {
}

func (m *manager) AddPod(pod *corev1.Pod) {
	klog.V(4).Infof("prober AddPod do nothing")
	return
}

func (m *manager) RemovePod(pod *corev1.Pod) {
	klog.V(4).Infof("prober RemovePod do nothing")
	return
}

func (m *manager) CleanupPods(desiredPods map[types.UID]sets.Empty) {
	klog.V(4).Infof("prober CleanupPods do nothing")
	return
}

func (m *manager) UpdatePodStatus(podUID types.UID, podStatus *corev1.PodStatus) {
	for i, c := range podStatus.ContainerStatuses {
		var started bool
		if c.State.Running == nil {
			started = false
		} else {
			// The check whether there is a probe which hasn't run yet.
			// _, exists := m.getWorker(podUID, c.Name, startup)
			started = true //!exists
		}
		podStatus.ContainerStatuses[i].Started = &started

		if started {
			var ready bool
			if c.State.Running == nil {
				ready = false
			} else {
				// The check whether there is a probe which hasn't run yet.
				//_, exists := m.getWorker(podUID, c.Name, readiness)
				ready = true // !exists
			}
			podStatus.ContainerStatuses[i].Ready = ready
		}
	}
	// init containers are ready if they have exited with success or if a readiness probe has
	// succeeded.
	for i, c := range podStatus.InitContainerStatuses {
		var ready bool
		if c.State.Terminated != nil && c.State.Terminated.ExitCode == 0 {
			ready = true
		}
		podStatus.InitContainerStatuses[i].Ready = ready
	}
}

func (m *manager) Start() {
	klog.V(4).Infof("prober start do nothing")
}

// NewManager creates a Manager for pod probing.
func NewManager(
	statusManager status.Manager,
	livenessManager results.Manager,
	startupManager results.Manager,
	runner kubecontainer.CommandRunner,
	recorder record.EventRecorder) prober.Manager {

	return &manager{}
}

var _ prober.Manager = &manager{}
