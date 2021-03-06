load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "metrics_cached.go",
        "metrics_du.go",
        "metrics_errors.go",
        "metrics_nil.go",
        "metrics_statfs.go",
        "noop_expandable_plugin.go",
        "plugins.go",
        "volume.go",
        "volume_linux.go",
        "volume_unsupported.go",
    ],
    importpath = "k8s.io/kubernetes/pkg/volume",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/proxy/util:go_default_library",
        "//pkg/volume/util/fs:go_default_library",
        "//pkg/volume/util/hostutil:go_default_library",
        "//pkg/volume/util/recyclerclient:go_default_library",
        "//pkg/volume/util/subpath:go_default_library",
        "//staging/src/k8s.io/api/authentication/v1:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/api/resource:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/types:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/errors:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/sets:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/util/validation:go_default_library",
        "//staging/src/k8s.io/client-go/informers:go_default_library",
        "//staging/src/k8s.io/client-go/kubernetes:go_default_library",
        "//staging/src/k8s.io/client-go/listers/storage/v1:go_default_library",
        "//staging/src/k8s.io/client-go/tools/cache:go_default_library",
        "//staging/src/k8s.io/client-go/tools/record:go_default_library",
        "//staging/src/k8s.io/cloud-provider:go_default_library",
        "//staging/src/k8s.io/mount-utils:go_default_library",
        "//vendor/k8s.io/klog/v2:go_default_library",
        "//vendor/k8s.io/utils/exec:go_default_library",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "//pkg/features:go_default_library",
            "//staging/src/k8s.io/apiserver/pkg/util/feature:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//pkg/features:go_default_library",
            "//staging/src/k8s.io/apiserver/pkg/util/feature:go_default_library",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "go_default_test",
    srcs = [
        "metrics_du_test.go",
        "metrics_nil_test.go",
        "metrics_statfs_test.go",
        "plugins_test.go",
        "volume_linux_test.go",
    ],
    embed = [":go_default_library"],
    deps = [
        "//pkg/volume/testing:go_default_library",
        "//staging/src/k8s.io/api/core/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/types:go_default_library",
        "//staging/src/k8s.io/client-go/util/testing:go_default_library",
    ] + select({
        "@io_bazel_rules_go//go/platform:android": [
            "//pkg/features:go_default_library",
            "//staging/src/k8s.io/apiserver/pkg/util/feature:go_default_library",
            "//staging/src/k8s.io/component-base/featuregate/testing:go_default_library",
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//pkg/features:go_default_library",
            "//staging/src/k8s.io/apiserver/pkg/util/feature:go_default_library",
            "//staging/src/k8s.io/component-base/featuregate/testing:go_default_library",
            "//vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "//conditions:default": [],
    }),
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
        "//pkg/volume/configmap:all-srcs",
        "//pkg/volume/csimigration:all-srcs",
        "//pkg/volume/emptydir:all-srcs",
        "//pkg/volume/hostpath:all-srcs",
        "//pkg/volume/secret:all-srcs",
        "//pkg/volume/testing:all-srcs",
        "//pkg/volume/util:all-srcs",
        "//pkg/volume/validation:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
