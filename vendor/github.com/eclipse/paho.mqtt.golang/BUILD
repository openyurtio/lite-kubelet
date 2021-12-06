load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "components.go",
        "filestore.go",
        "memstore.go",
        "message.go",
        "messageids.go",
        "net.go",
        "netconn.go",
        "oops.go",
        "options.go",
        "options_reader.go",
        "ping.go",
        "router.go",
        "store.go",
        "token.go",
        "topic.go",
        "trace.go",
        "websocket.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/github.com/eclipse/paho.mqtt.golang",
    importpath = "github.com/eclipse/paho.mqtt.golang",
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/eclipse/paho.mqtt.golang/packets:go_default_library",
        "//vendor/github.com/gorilla/websocket:go_default_library",
        "//vendor/golang.org/x/net/proxy:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [
        ":package-srcs",
        "//vendor/github.com/eclipse/paho.mqtt.golang/packets:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
