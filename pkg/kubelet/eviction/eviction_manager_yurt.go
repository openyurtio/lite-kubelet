package eviction

import (
	"time"

	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
)

type yurtManagerImpl struct {
}

// NewManager returns a configured Manager and an associated admission handler to enforce eviction configuration.
func NewManagerYurt() (Manager, lifecycle.PodAdmitHandler) {
	manage := &yurtManagerImpl{}
	return manage, manage
}

func (y *yurtManagerImpl) Admit(attrs *lifecycle.PodAdmitAttributes) lifecycle.PodAdmitResult {
	klog.V(4).Infof("implement me: yurtManagerImpl Admit")
	return lifecycle.PodAdmitResult{Admit: true}
}

func (y *yurtManagerImpl) Start(diskInfoProvider DiskInfoProvider, podFunc ActivePodsFunc, podCleanedUpFunc PodCleanedUpFunc, monitoringInterval time.Duration) {
	klog.V(4).Infof("implement me: yurtManagerImpl Start")
	return
}

func (y *yurtManagerImpl) IsUnderMemoryPressure() bool {
	klog.V(4).Infof("implement me: yurtManagerImpl IsUnderMemoryPressure")
	return true
}

func (y *yurtManagerImpl) IsUnderDiskPressure() bool {
	klog.V(4).Infof("implement me: yurtManagerImpl IsUnderDiskPressure")
	return true
}

func (y *yurtManagerImpl) IsUnderPIDPressure() bool {
	klog.V(4).Infof("implement me: yurtManagerImpl IsUnderPIDPressure")
	return true
}

// ensure it implements the required interface
var _ Manager = &yurtManagerImpl{}
var _ lifecycle.PodAdmitHandler = &yurtManagerImpl{}
