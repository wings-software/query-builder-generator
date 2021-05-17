load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"], exclude = ["*_test.go"]),
    importpath = "github.com/query-builder-generator/src/compiler",
    visibility = ["//visibility:public"],
    deps = [
        "//src/dom:go_default_library",
    ]
)

go_test(
    name = "go_default_test",
    srcs = glob(["*_test.go"]),
    embed = [":go_default_library"],
    deps = [
        "//src/dom:go_default_library",
        "@com_github_stretchr_testify//assert:go_default_library",
    ]
)