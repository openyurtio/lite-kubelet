#!/bin/bash

set -e

BINARY=$0
SUB_CMD=$1


function console() {

> test.log


./kubelet --container-runtime=remote --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock --container-log-max-files 10 --container-log-max-size=100Mi --max-pods 10 --pod-max-pids 16384 --pod-manifest-path=/etc/kubernetes/manifests --v 4 --hostname-override={nodename} --cgroup-driver=systemd  --mqtt-access-key={access key} --mqtt-broker="{broker}" --mqtt-broker-port 8883 --mqtt-group="{group id}" --mqtt-instance="{mqtt instance}" --mqtt-root-topic="{topic}"  --mqtt-secret-key="{secret key}" 2>&1  | tee test.log

# docker
#./kubelet --container-log-max-files 10 --container-log-max-size=100Mi --max-pods 10 --pod-max-pids 16384 --pod-manifest-path=/etc/kubernetes/manifests --v 4 --hostname-override={nodename} --cgroup-driver=systemd  --mqtt-access-key={access key} --mqtt-broker="{mqtt broker}" --mqtt-broker-port 8883 --mqtt-group="{mqtt group}" --mqtt-instance="{mqtt instance}" --mqtt-root-topic="{mqtt topic}"  --mqtt-secret-key="{mqtt secret}" --pod-infra-container-image=registry-vpc.cn-zhangjiakou/acs/pause:3.5  2>&1  | tee test.log

}

function log() {

        echo "nil"

}

function help() {
        echo """
$BINARY console
$BINARY log
"""
}

case $SUB_CMD in
    "console")
        console
        ;;
    "log")
        log
        ;;
    *)
        echo "wrong cmd , help info:"
        help
esac


