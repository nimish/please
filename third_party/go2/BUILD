go_remote_library(
    name = 'atime',
    get = 'github.com/djherbis/atime',
    revision = '89517e96e10b93292169a79fd4523807bdc5d5fa',
)

go_remote_library(
    name = 'client_golang',
    get = 'github.com/prometheus/client_golang',
    packages = [
        'prometheus',
    ],
    revision = '738ed6c0b9b30ec936d4b74c1600ae76324efff5',
    deps = [
        ':client_model',
        ':common',
        ':perks',
        ':procfs',
        ':protobuf',
    ],
)

go_remote_library(
    name = 'client_model',
    get = 'github.com/prometheus/client_model',
    packages = [
        'go',
    ],
    revision = '6f3806018612930941127f2a7c6c453ba2c527d2',
    deps = [
        ':protobuf',
    ],
)

go_remote_library(
    name = 'common',
    get = 'github.com/prometheus/common',
    packages = [
        'expfmt',
        'internal/bitbucket.org/ww/goautoneg',
        'model',
    ],
    revision = '49fee292b27bfff7f354ee0f64e1bc4850462edf',
    deps = [
        ':client_model',
        ':golang_protobuf_extensions',
        ':protobuf',
    ],
)

go_remote_library(
    name = 'concurrent-map',
    get = 'github.com/streamrail/concurrent-map',
    revision = '8bf1e9bacbf65b10c81d0f4314cf2b1ebef728b5',
)

go_remote_library(
    name = 'context',
    get = 'github.com/gorilla/context',
    revision = '08b5f424b9271eedf6f9f0ce86cb9396ed337a42',
)

go_remote_library(
    name = 'crypto',
    get = 'golang.org/x/crypto',
    packages = [
        'ssh/terminal',
    ],
    repo = 'go.googlesource.com/crypto',
    revision = '459e26527287adbc2adcc5d0d49abff9a5f315a7',
)

go_remote_library(
    name = 'dns',
    get = 'github.com/miekg/dns',
    revision = '767422ac12884e2baed0afd7303cf06cff90fef6',
)

go_remote_library(
    name = 'fsnotify',
    get = 'github.com/fsnotify/fsnotify',
    revision = '4da3e2cfbabc9f751898f250b49f2439785783a1',
    deps = [
        ':sys',
    ],
)

go_remote_library(
    name = 'gcfg.v1',
    get = 'gopkg.in/gcfg.v1',
    revision = '27e4946190b4a327b539185f2b5b1f7c84730728',
    deps = [
        ':warnings.v0',
    ],
)

go_remote_library(
    name = 'go-bindata',
    get = 'github.com/jteeuwen/go-bindata',
    packages = [
        '',
        'go-bindata',
    ],
    revision = 'a0ff2567cfb70903282db057e799fd826784d41d',
)

go_remote_library(
    name = 'go-datastructures',
    get = 'github.com/Workiva/go-datastructures',
    packages = [
        'queue',
    ],
    revision = 'dd1dc71ac92298310606ac77d9519ddf1c6e2500',
)

go_remote_library(
    name = 'go-flags',
    get = 'github.com/jessevdk/go-flags',
    revision = '460c7bb0abd6e927f2767cadc91aa6ef776a98b4',
)

go_remote_library(
    name = 'go-humanize',
    get = 'github.com/dustin/go-humanize',
    revision = '259d2a102b871d17f30e3cd9881a642961a1e486',
)

go_remote_library(
    name = 'go-isatty',
    get = 'github.com/mattn/go-isatty',
    revision = 'fc9e8d8ef48496124e79ae0df75490096eccf6fe',
)

go_remote_library(
    name = 'go-logging.v1',
    get = 'gopkg.in/op/go-logging.v1',
    revision = 'b2cb9fa56473e98db8caba80237377e83fe44db5',
)

go_remote_library(
    name = 'go-metrics',
    get = 'github.com/armon/go-metrics',
    revision = '93f237eba9b0602f3e73710416558854a81d9337',
)

go_remote_library(
    name = 'go-msgpack',
    get = 'github.com/hashicorp/go-msgpack',
    packages = [
        'codec',
    ],
    revision = 'fa3f63826f7c23912c15263591e65d54d080b458',
)

go_remote_library(
    name = 'go-multierror',
    get = 'github.com/hashicorp/go-multierror',
    packages = [
        '',
        'vendor/github.com/hashicorp/errwrap',
    ],
    revision = 'ed905158d87462226a13fe39ddf685ea65f1c11f',
)

go_remote_library(
    name = 'go-semver',
    get = 'github.com/coreos/go-semver',
    packages = [
        'semver',
    ],
    revision = '5e3acbb5668c4c3deb4842615c4098eb61fb6b1e',
)

go_remote_library(
    name = 'go-sockaddr',
    get = 'github.com/hashicorp/go-sockaddr',
    revision = 'acd314c5781ea706c710d9ea70069fd2e110d61d',
)

go_remote_library(
    name = 'golang-levenshtein',
    get = 'github.com/texttheater/golang-levenshtein',
    packages = [
        'levenshtein',
    ],
    revision = 'a809b0abc4bdd4591a275c2dcf0982b0581da4e6',
)

go_remote_library(
    name = 'golang_protobuf_extensions',
    get = 'github.com/matttproud/golang_protobuf_extensions',
    packages = [
        'pbutil',
    ],
    revision = 'c12348ce28de40eed0136aa2b644d0ee0650e56c',
    deps = [
        ':protobuf',
    ],
)

go_remote_library(
    name = 'grpc-go',
    get = 'github.com/grpc/grpc-go',
    revision = '7fc29d0caa400efc53212b222356f99efe87d7d5',
    deps = [
        ':net',
        ':protobuf',
    ],
)

go_remote_library(
    name = 'memberlist',
    get = 'github.com/hashicorp/memberlist',
    revision = '6cc6075ba9fba1915fa0416f00d2b4efa9dc2262',
    deps = [
        ':dns',
        ':go-metrics',
        ':go-msgpack',
        ':go-multierror',
        ':go-sockaddr',
        ':seed',
    ],
)

go_remote_library(
    name = 'mux',
    get = 'github.com/gorilla/mux',
    revision = '599cba5e7b6137d46ddf58fb1765f5d928e69604',
)

go_remote_library(
    name = 'net',
    get = 'golang.org/x/net',
    packages = [
        'context',
    ],
    repo = 'go.googlesource.com/net',
    revision = 'ffcf1bedda3b04ebb15a168a59800a73d6dc0f4d',
    deps = [
        ':text',
    ],
)

go_remote_library(
    name = 'osext',
    get = 'github.com/kardianos/osext',
    revision = '9d302b58e975387d0b4d9be876622c86cefe64be',
)

go_remote_library(
    name = 'perks',
    get = 'github.com/beorn7/perks',
    packages = [
        'quantile',
    ],
    revision = '4c0e84591b9aa9e6dcfdf3e020114cd81f89d5f9',
)

go_remote_library(
    name = 'procfs',
    get = 'github.com/prometheus/procfs',
    revision = 'a1dba9ce8baed984a2495b658c82687f8157b98f',
    strip = ['fixtures'],
)

go_remote_library(
    name = 'prompter',
    get = 'github.com/Songmu/prompter',
    revision = 'f49666b0047d12850875d771340e1d862d9e7a0c',
    deps = [
        ':crypto',
        ':go-isatty',
    ],
)

go_remote_library(
    name = 'protobuf',
    get = 'github.com/golang/protobuf',
    packages = [
        'protoc-gen-go',
    ],
    revision = 'c9c7427a2a70d2eb3bafa0ab2dc163e45f143317',
)

go_remote_library(
    name = 'seed',
    get = 'github.com/sean-/seed',
    revision = 'e2103e2c35297fb7e17febb81e49b312087a2372',
)

go_remote_library(
    name = 'shlex',
    get = 'github.com/google/shlex',
    revision = '6f45313302b9c56850fc17f99e40caebce98c716',
)

go_remote_library(
    name = 'sys',
    get = 'golang.org/x/sys',
    packages = [
        'unix',
    ],
    repo = 'go.googlesource.com/sys',
    revision = '9a7256cb28ed514b4e1e5f68959914c4c28a92e0',
)

go_remote_library(
    name = 'testify',
    get = 'github.com/stretchr/testify',
    revision = '4d4bfba8f1d1027c4fdbe371823030df51419987',
)

go_remote_library(
    name = 'text',
    get = 'golang.org/x/text',
    repo = 'go.googlesource.com/text',
    revision = '3aa6d803466d99ea1bdf4af99e6db53299c02bd7',
)

go_remote_library(
    name = 'tools',
    get = 'golang.org/x/tools',
    packages = [
        'cmd/stringer',
    ],
    repo = 'go.googlesource.com/tools',
    revision = 'd63e2b22b05a9682de336cd4802bba367ed429e7',
)

go_remote_library(
    name = 'warnings.v0',
    get = 'gopkg.in/warnings.v0',
    revision = '8a331561fe74dadba6edfc59f3be66c22c3b065d',
)
