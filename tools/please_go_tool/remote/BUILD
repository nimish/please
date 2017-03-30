go_library(
    name = 'remote',
    srcs = ['remote.go'],
    deps = [
        '//third_party/go:logging',
    ],
    visibility = ['//tools/please_go_tool/...'],
)

go_test(
    name = 'remote_test',
    srcs = ['remote_test.go'],
    deps = [
        ':remote',
        '//third_party/go:testify',
    ],
)
