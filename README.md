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
  
  Kubelet gets the pod declaration file in three ways: static pod Path, manifest-url, and apiserver.
  
  We removed the logic to get pod declaration files from manifest-URL and Apiserver. Added new logic to get pod declaration file for mqttfile.
  
  To get the pod declaration file via mqtt, the pod contents are first retrieved from the MQTT, cached locally, and then created in a staticPod-like fashion 
  
 pkg/kubelet/kubelet.go L278
 ```
 	updates := cfg.Channel(kubetypes.MqttFileSource)
 	send := func(cache cache.Indexer) {
 		pods := make([]*v1.Pod, 0, 10)
 		for _, o := range cache.List() {
 			if p, ok := o.(*v1.Pod); ok {
 				klog.V(4).Infof("Get Pod [%s][%s] from local mqtt cache", p.GetNamespace(), p.GetName())
 				pods = append(pods, p)
 			}
 		}
 		updates <- kubetypes.PodUpdate{Pods: pods, Op: kubetypes.SET, Source: kubetypes.MqttFileSource}
 	}
 	fileCache.NewFileObiectIndexer(fileCache.NewDefaultFilePodDeps(), false, send)
 
 ```

- The `clientset.Interface` Interface is reimplemented
  We redefined the `Clientset` structure, reimplemented the [clientset.Interface] Interface through [fakekube.Clientset], and reimplemented some methods.
  
  The method of reimplementation is as follows:
  ```
  CoreV1().Pods().Get()  
  CoreV1().Pods().Create()  
  CoreV1().Pods().Delete()  
  CoreV1().Pods().Patch()  
  ```
  
  ```
  CoreV1().Nodes().Create()
  CoreV1().Nodes().Patch()
  CoreV1().Nodes().Delete()
  ```
  
  ```
  CoreV1().Events().CreateWithEventNamespace()
  CoreV1().Events().UpdateWithEventNamespace()  
  CoreV1().Events().PatchWithEventNamespace()
  ```
  
  ```
  CoordinationV1().Leases().Get()
  CoordinationV1().Leases().Create()
  CoordinationV1().Leases().Update()
  ```
  
   After tailoring kubelet, we found that kubelet only needs to interact with `Node/Pod/Events/Lease` resource objects. So we reimplemented the methods of these resource objects from the interface to call apiserver to interact through MQTT. Then kole-Controller, the cloud controller, requests apiserver on behalf of it
  
     
  ref: pkg/openyurt/clientSet/clientset.go
  

## TODO list

- Change the built binary name kubelet to lite-kubelet
- Delete irrelevant code

[kubernetes]: https://github.com/kubernetes/kubernetes
[Kubelet]: https://github.com/kubernetes/kubernetes/tree/master/cmd/kubelet
[clientset.Interface]: https://github.com/kubernetes/client-go/blob/cc43a708a08eb9ff6a436f0cb00c5ee05121d2cd/kubernetes/clientset.go#L75
[fakekube.Clientset]: https://github.com/kubernetes/client-go/blob/cc43a708a08eb9ff6a436f0cb00c5ee05121d2cd/kubernetes/fake/clientset_generated.go#L151