go_library(
    name = 'goget',
    srcs = ['goget.go'],
    deps = [
        '//third_party/go:logging',
    ],
    visibility = ['//tools/please_go_tool/...'],
)

go_test(
    name = 'goget_test',
    srcs = ['goget_test.go'],
    deps = [
        ':goget',
        '//third_party/go:testify',
    ],
)
