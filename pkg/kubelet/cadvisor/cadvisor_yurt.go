package cadvisor

import (
	"errors"

	"github.com/google/cadvisor/container/docker"
	"github.com/google/cadvisor/events"
	"github.com/google/cadvisor/fs"
	cadvisorapi "github.com/google/cadvisor/info/v1"
	info "github.com/google/cadvisor/info/v1"
	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
	"github.com/google/cadvisor/machine"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/version"
	"k8s.io/klog/v2"
)

type cadvisorClientYurt struct {
	sysFs  sysfs.SysFs
	fsInfo fs.FsInfo
}

// NewYurt creates a new cAdvisor Interface for linux systems.
func NewYurt() (Interface, error) {
	sysFs := sysfs.NewRealSysFs()
	fsInfo, err := fs.NewFsInfo(fs.Context{})
	if err != nil {
		return nil, err
	}

	return &cadvisorClientYurt{
		sysFs:  sysFs,
		fsInfo: fsInfo,
	}, nil
}

var errUnsupported = errors.New("cAdvisor yurt is unsupported in this build")

func (c cadvisorClientYurt) Start() error {
	klog.V(4).Infof("cadvisorClientYurt Started, by do nothing")
	return nil
}

func (c cadvisorClientYurt) DockerContainer(name string, req *cadvisorapi.ContainerInfoRequest) (cadvisorapi.ContainerInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: DockerContainer")
	return cadvisorapi.ContainerInfo{}, errUnsupported
}

func (c cadvisorClientYurt) ContainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (*cadvisorapi.ContainerInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: ContainerInfo")
	return nil, errUnsupported
}

func (c cadvisorClientYurt) ContainerInfoV2(name string, options cadvisorapiv2.RequestOptions) (map[string]cadvisorapiv2.ContainerInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: ContainerInfoV2")
	return nil, errUnsupported
}

func (c cadvisorClientYurt) GetRequestedContainersInfo(containerName string, options cadvisorapiv2.RequestOptions) (map[string]*cadvisorapi.ContainerInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: GetRequestedContainersInfo")
	return nil, errUnsupported
}

func (c cadvisorClientYurt) SubcontainerInfo(name string, req *cadvisorapi.ContainerInfoRequest) (map[string]*cadvisorapi.ContainerInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: SubcontainerInfo")
	return nil, errUnsupported
}

func (c cadvisorClientYurt) MachineInfo() (*cadvisorapi.MachineInfo, error) {
	klog.V(4).Infof("cadvisorClientYurt: MachineInfo")
	return machine.Info(c.sysFs, c.fsInfo, true)
}

func (c cadvisorClientYurt) VersionInfo() (*cadvisorapi.VersionInfo, error) {
	klog.V(4).Infof("cadvisorClientYurt: VersionInfo")
	kernelVersion := machine.KernelVersion()
	osVersion := machine.ContainerOsVersion()
	dockerVersion, err := docker.VersionString()
	if err != nil {
		return nil, err
	}
	dockerAPIVersion, err := docker.APIVersionString()
	if err != nil {
		return nil, err
	}

	return &info.VersionInfo{
		KernelVersion:      kernelVersion,
		ContainerOsVersion: osVersion,
		DockerVersion:      dockerVersion,
		DockerAPIVersion:   dockerAPIVersion,
		CadvisorVersion:    version.Info["version"],
		CadvisorRevision:   version.Info["revision"],
	}, nil
}

func (c cadvisorClientYurt) ImagesFsInfo() (cadvisorapiv2.FsInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: ImagesFsInfo")
	return cadvisorapiv2.FsInfo{}, errUnsupported
}

func (c cadvisorClientYurt) RootFsInfo() (cadvisorapiv2.FsInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: RootFsInfo")
	return cadvisorapiv2.FsInfo{}, errUnsupported
}

func (c cadvisorClientYurt) WatchEvents(request *events.Request) (*events.EventChannel, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: WatchEvents")
	return nil, errUnsupported
}

func (c cadvisorClientYurt) GetDirFsInfo(path string) (cadvisorapiv2.FsInfo, error) {
	klog.V(4).Infof("implement me: cadvisorClientYurt: GetDirFsInfo")
	return cadvisorapiv2.FsInfo{}, nil
}

var _ Interface = &cadvisorClientYurt{}
