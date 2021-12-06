package(default_visibility = ["//visibility:public"])

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
        "//cmd/clicheck:all-srcs",
        "//cmd/dependencycheck:all-srcs",
        "//cmd/gendocs:all-srcs",
        "//cmd/genman:all-srcs",
        "//cmd/genswaggertypedocs:all-srcs",
        "//cmd/genutils:all-srcs",
        "//cmd/genyaml:all-srcs",
        "//cmd/importverifier:all-srcs",
        "//cmd/kubectl:all-srcs",
        "//cmd/kubectl-convert:all-srcs",
        "//cmd/kubelet:all-srcs",
        "//cmd/linkcheck:all-srcs",
        "//cmd/preferredimports:all-srcs",
        "//cmd/verifydependencies:all-srcs",
    ],
    tags = ["automanaged"],
)
