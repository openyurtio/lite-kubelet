load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["clientset.go"],
    importpath = "k8s.io/kubernetes/pkg/openyurt/clientSet",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/openyurt/clientSet/typed/coordination/v1:go_default_library",
        "//pkg/openyurt/clientSet/typed/core/v1:go_default_library",
        "//pkg/openyurt/mqtt/client:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/fake:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/typed/coordination/v1:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/typed/coordination/v1/fake:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes/typed/core/v1/fake:go_default_library",
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
        "//pkg/openyurt/clientSet/typed/coordination/v1:all-srcs",
        "//pkg/openyurt/clientSet/typed/core/v1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
