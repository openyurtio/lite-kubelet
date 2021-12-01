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

function help() {
        echo """
$BINARY build 
	单纯构建linux 环境的kubelet
$BINARY scp 
	将linux 环境的kubelet scp 到目标机器
$BINARY all 
    即构建，又scp
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
    "all")
        buildbin
	    scpbin	
        ;;

    *)
        echo "wrong cmd , help info:"
        help
esac
