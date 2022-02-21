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

[kubernetes]: https://github.com/kubernetes/kubernetes