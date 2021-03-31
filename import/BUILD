load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "googlevoice",
    srcs = ["googlevoice.go"],
    importpath = "github.com/minkezhang/archive-pipeline/import/googlevoice",
    deps = [
        "//api:constants_go_proto",
        "//api:data_go_proto",
	"@org_golang_google_protobuf//types/known/timestamppb:go_default_library",
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
        ":googlevoice",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
