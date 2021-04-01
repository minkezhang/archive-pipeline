load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

package(default_visibility=["//visibility:public"])

go_library(
    name = "importer",
    srcs = ["importer.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/googlevoice/importer",
    deps = [
        ":transaction",
        "//import:contact",
        "//import:record",
    ],
)

go_library(
    name = "transaction",
    srcs = ["transaction.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/googlevoice/transaction",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
        "@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
    ],
)
