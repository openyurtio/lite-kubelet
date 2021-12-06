load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "connack.go",
        "connect.go",
        "disconnect.go",
        "packets.go",
        "pingreq.go",
        "pingresp.go",
        "puback.go",
        "pubcomp.go",
        "publish.go",
        "pubrec.go",
        "pubrel.go",
        "suback.go",
        "subscribe.go",
        "unsuback.go",
        "unsubscribe.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/eclipse/paho.mqtt.golang/packets",
    importpath = "github.com/eclipse/paho.mqtt.golang/packets",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
