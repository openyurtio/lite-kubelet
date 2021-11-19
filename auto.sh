#!/bin/bash

KUBE_BUILD_PLATFORMS=linux/amd64 make WHAT=cmd/kubelet GOFLAGS=-v && scp _output/local/bin/linux/amd64/kubelet testkubelet:~/download/

