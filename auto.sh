#!/bin/bash

#set -e

BINARY=$0
SUB_CMD=$1

LINUX_KUBELET_BIN=_output/local/bin/linux/amd64/kubelet

function buildbin() {

	echo "完整版kubelet 大小:119 M"
	echo "原来的kubelet 大小:"
	du -m ${LINUX_KUBELET_BIN}
	rm -rf ${LINUX_KUBELET_BIN}

	echo "开始构建 ..."
	KUBE_BUILD_PLATFORMS=linux/amd64 make WHAT=cmd/kubelet GOFLAGS=-v 

	echo "现在的kubelet 大小:"
	date 2>&1 | >> kubelet.txt
	du -m ${LINUX_KUBELET_BIN}

}

function scpbin() {
	scp ${LINUX_KUBELET_BIN} testkubelet:~/download/
}

function scpbindocker() {
	scp ${LINUX_KUBELET_BIN} testkubeletdocker:~/download/
}

function camera() {

	echo "开始构建arm 版本"
	KUBE_BUILD_PLATFORMS=linux/arm make WHAT=cmd/kubelet GOFLAGS=-v
    cp _output/local/bin/linux/arm/kubelet  ~/Downloads/lite-kubelet/
}

function help() {
        echo """
$BINARY build 
	单纯构建linux 环境的kubelet
$BINARY scp 
	将linux 环境的kubelet scp 到目标机器
$BINARY scpdocker
	将linux 环境的kubelet scp 到目标docker机器上

$BINARY camera 
    即构建arm 版本，cp ~/Download 

$BINARY all 
    即构建，又scp
$BINARY alldocker
    即构建，又scp 到docker

        """
        exit 0
}

case $SUB_CMD in
    "build")
        buildbin 
        ;;
    "scp")
        scpbin 
        ;;
    "scpdocker")
        scpbindocker 
        ;;
    "all")
        buildbin
	    scpbin	
        ;;
    "alldocker")
        buildbin
	    scpbindocker	
        ;;
    "camera")
        camera 
        ;;
    *)
        echo "wrong cmd , help info:"
        help
esac
