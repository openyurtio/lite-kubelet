#!/bin/bash

set -e

BINARY=$0
SUB_CMD=$1


function console() {

> test.log

./kubelet --container-log-max-files 10 --container-log-max-size=100Mi --max-pods 64 --pod-max-pids 16384 --pod-manifest-path=/etc/kubernetes/manifests --v 4 --hostname-override=cn-beijing.172.23.142.26 --cgroup-driver=systemd --node-labels=type=lite-kubelet --mqtt-broker="post-cn-zvp2gmir011.mqtt.aliyuncs.com" --mqtt-broker-port 8883 --mqtt-clientid="GID_NODE@@@lite-26" --mqtt-username="Signature|access-key|post-cn-zvp2gmir011" --mqtt-password="paasswd="  --pod-infra-container-image=registry-vpc.cn-beijing.aliyuncs.com/acs/pause:3.5  2>&1  | tee test.log


}

function core() {
./lite-core --mqtt-broker="post-cn-zvp2gmir011.mqtt.aliyuncs.com" --mqtt-broker-port 8883 --mqtt-clientid="GID_NODE@@@lite-core" --mqtt-username="Signature|access key|post-cn-zvp2gmir011" --mqtt-password="passwd=" --kubeconfig=/root/.kube/config --v 4
}

function log() {

        echo "nil"

}

function help() {
        echo """
$BINARY console
$BINARY log
$BINARY core
"""
}

case $SUB_CMD in
    "console")
        console
        ;;
    "log")
        log
        ;;
    "core")
        core 
        ;;

    *)
        echo "wrong cmd , help info:"
        help
esac


# /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --max-pods 64 --pod-max-pids 16384 --pod-manifest-path=/etc/kubernetes/manifests --feature-gates=IPv6DualStack=true --network-plugin=cni --cni-conf-dir=/etc/cni/net.d --cni-bin-dir=/opt/cni/bin --dynamic-config-dir=/etc/kubernetes/kubelet-config --v=3 --enable-controller-attach-detach=true --cluster-dns=192.168.0.10 --pod-infra-container-image=registry-vpc.cn-beijing.aliyuncs.com/acs/pause:3.5 --enable-load-reader --cluster-domain=cluster.local --cloud-provider=external --hostname-override=cn-beijing.172.23.142.26 --provider-id=cn-beijing.i-2zehxjofqxse3cpv3fm5 --authorization-mode=Webhook --authentication-token-webhook=true --anonymous-auth=false --client-ca-file=/etc/kubernetes/pki/ca.crt --cgroup-driver=systemd --tls-cipher-suites=TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256 --tls-cert-file=/var/lib/kubelet/pki/kubelet.crt --tls-private-key-file=/var/lib/kubelet/pki/kubelet.key --rotate-certificates=true --cert-dir=/var/lib/kubelet/pki --node-labels=alibabacloud.com/nodepool-id=npd9fdb5079fc849e5837c357d8ea9082b,ack.aliyun.com=c518956bb66714443a8a2d05c9f148a2a --eviction-hard=imagefs.available<15%,memory.available<300Mi,nodefs.available<10%,nodefs.inodesFree<5% --system-reserved=cpu=50m,memory=849Mi --kube-reserved=cpu=50m,memory=849Mi --kube-reserved=pid=1000 --system-reserved=pid=1000
