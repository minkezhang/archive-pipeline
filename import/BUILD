load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "contact",
    srcs = ["contact.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/contact",
    deps = [
        "//api:data_go_proto",
    ],
)

go_library(
    name = "record",
    srcs = ["record.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/record",
    deps = [
        ":contact",
        "//api:data_go_proto",
    ],
)

go_library(
    name = "interfaces",
    srcs = ["interfaces.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/interfaces",
    deps = [
        ":record",
    ],
)

go_binary(
    name = "main",
    srcs = ["main.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/main",
    data = [
        "//data:data",
    ],
    deps = [
        ":interfaces",
        ":record",
        "//import/googlevoice:importer",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
