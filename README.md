# lite-kubelet 

lite-kubelet is a lightweight version of Kubelet. It mainly changes the communication mechanism between Kubelet and APIServer. lite-kubelet uses MQTT to communicate with APIServer.

We also trimmed some of the code in lite-kubelet, but retained the core kubelet logic to make it less memory and CPU intensive.

lite-kubelet code is forked from [kubernetes].

----

## To start developing K8s

##### You have a working [Go environment].

```
mkdir -p $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io
git clone https://github.com/openyurtio/lite-kubelet.git
mv -vf lite-kubelet  kubernetes
# for linux, amd64
KUBE_BUILD_PLATFORMS=linux/amd64 make WHAT=cmd/kubelet GOFLAGS=-v
```

## What changes did we make based on [Kubelet]

  We try to keep the core logic of Kubelet, some of the code will be trimmed, some of the interface will be re-implemented, the re-implemented code will be put in the `pkg/openyurt` directory in the future.

- Baseon kubernetes v1.20.4 code 

- Delete KubeletConfigController DynamicKubeletConfig

  Lite-kubelet does not know the address of apiserver, he only needs to know the information about the MQTT broker
  
  cmd/kubelet/app/server.go 
  
  ```
  var kubeletConfigController *dynamickubeletconfig.Controller
  ``` 
 
- Delete internal plugin tree: 

  configmap, secret,cephfs, csi, downwardapi, fc, flocker, git_repo, glusterfs, iscsi, nfs, portworx, projected, quobyte, rbd, scaleio, storageos and so on. 
  Only keep hostpath, emptydir.  
  configmap and secret will be supported in the future.
  
  /cmd/kubelet/app/plugins.go
  ```
  func ProbeVolumePlugins(featureGate featuregate.FeatureGate) ([]volume.VolumePlugin, error) {
  ``` 
  
- Disable start kubelet server: 10250, 10255 port
  cmd/kubelet/app/server.go
  ```
    func startKubelet(k kubelet.Bootstrap, podCfg *config.PodConfig, kubeCfg *kubeletconfiginternal.KubeletConfiguration, kubeDeps *kubelet.Dependencies, enableCAdvisorJSONEndpoints, enableServer bool) {

  ```
  
- Disable crio server

- Disable eviction_manager 
  Reimplements the Eviction.Manager interface, but  do nothing
  
  pkg/kubelet/eviction/eviction_manager_yurt.go
    
  It will be placed in `pkg/openyurt` dir in the future   
   
  pkg/kubelet/kubelet.go
  ```
  evictionManager, evictionAdmitHandler := eviction.NewManagerYurt()
  ```
- Disable cadvisor
  Reimplements the CAdvisorInterface , but do nothing
  
  pkg/kubelet/cadvisor/cadvisor_yurt.go
  
  It will be placed in `pkg/openyurt` dir in the future   
  
- Use containerd as default runtime

- Disable pod probe manager(livenessManger and startupManager)

- Delete runtimeClassManager
  
- Add mqtt source file 
 
 
## TODO list

1 替换构建的二进制名字
2 删除无关的代码

[kubernetes]: https://github.com/kubernetes/kubernetes
[Kubelet]: https://github.com/kubernetes/kubernetes/tree/master/cmd/kubelet