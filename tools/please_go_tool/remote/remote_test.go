package remote

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueDeps(t *testing.T) {
	jps := jsonPackages{
		&jsonPackage{ImportPath: "core", Deps: []string{"net", "net/http", "sync", "archive/tar"}},
		&jsonPackage{ImportPath: "jarcat", Deps: []string{"archive/tar", "archive/zip"}},
	}
	expected := []string{"archive/tar", "archive/zip", "core", "jarcat", "net", "net/http", "sync"}
	assert.Equal(t, expected, jps.UniqueDeps())
}

func TestToMap(t *testing.T) {
	jps := jsonPackages{
		&jsonPackage{ImportPath: "github.com/thought-machine/please/tools"},
		&jsonPackage{ImportPath: "github.com/thought-machine/please/src"},
	}
	expected := map[string]*jsonPackage{
		"github.com/thought-machine/please/tools": &jsonPackage{ImportPath: "github.com/thought-machine/please/tools"},
		"github.com/thought-machine/please/src":   &jsonPackage{ImportPath: "github.com/thought-machine/please/src"},
	}
	assert.Equal(t, expected, jps.ToMap())
}

func TestToGitMap(t *testing.T) {
	jps := jsonPackages{}
	assert.NoError(t, jps.FromJSON([]byte(samplePackages)))
	jps[0].GitURL = "github.com/thought-machine/please"
	jps[1].GitURL = "github.com/thought-machine/please"
	jps[2].GitURL = "github.com/thought-machine/pleasings"
	expected := map[string]*jsonPackage{
		"github.com/thought-machine/please":    jps[1],
		"github.com/thought-machine/pleasings": jps[2],
	}
	assert.Equal(t, expected, jps.ToGitMap([]string{"cli"}))
}

func TestToShortFormatString(t *testing.T) {
	jps := jsonPackages{}
	assert.NoError(t, jps.FromJSON([]byte(samplePackages)))
	expected := "core|src/core|core/core.a|build_env.go,build_label.go,build_target.go,cache.go,config.go,file_label.go,glob.go,graph.go,lock.go,package.go,state.go,test_results.go,utils.go,version.go|cli\n"
	assert.Equal(t, expected, jps[0].ToShortFormatString(jps.ToMap()))
}

func TestToShortFormatStringCgo(t *testing.T) {
	jps := jsonPackages{}
	assert.NoError(t, jps.FromJSON([]byte(samplePackages)))
	expected := "parse|src/parse|parse/parse.a|builtin_rules.go,parse_step.go,suggest.go|interpreter.go|interpreter.c|interpreter.h|--std=c99^-Werror|-ldl|core\n"
	assert.Equal(t, expected, jps[2].ToShortFormatString(jps.ToMap()))
}

func TestRepoNameToRuleName(t *testing.T) {
	assert.Equal(t, "please", repoNameToRuleName("github.com/thought-machine/please"))
	assert.Equal(t, "please", repoNameToRuleName("github.com/thought-machine/please.git"))
}

func TestToBuildRule(t *testing.T) {
	jps := jsonPackages{}
	assert.NoError(t, jps.FromJSON([]byte(samplePackages)))
	jps[0].GitURL = "github.com/thought-machine/please"
	jps[1].GitURL = "github.com/thought-machine/please"
	jps[2].GitURL = "github.com/thought-machine/pleasings"
	jps[0].Revision = "cf4e57e3bc210d18d3e6caedb7db6b57655e2be8"
	jps[1].Revision = "cf4e57e3bc210d18d3e6caedb7db6b57655e2be8"
	jps[2].Revision = "b916153623b843b3b4a34854bfd5ecb4577c083f"
	jps[2].Imports = append(jps[2].Imports, "github.com/thought-machine/please")
	originalPackages := []string{"github.com/thought-machine/please", "github.com/thought-machine/pleasings/grm"}
	expected1 := `go_remote_library(
    name = 'please',
    get = 'github.com/thought-machine/please',
    revision = 'cf4e57e3bc210d18d3e6caedb7db6b57655e2be8',
)
`
	expected2 := `go_remote_library(
    name = 'pleasings',
    get = 'github.com/thought-machine/pleasings',
    packages = [
        'grm',
    ],
    revision = 'b916153623b843b3b4a34854bfd5ecb4577c083f',
    deps = [
        ':please',
    ],
)
`
	expected3 := `go_remote_library(
    name = 'please',
    get = 'github.com/thought-machine/please',
    revision = 'cf4e57e3bc210d18d3e6caedb7db6b57655e2be8',
    hashes = ['lKBCbo0yA9pUaMzwxiT5PLN2AeI'],
)
`
	assert.Equal(t, expected1, jps[0].ToBuildRule(false, jps.ToGitMap(nil), nil, originalPackages))
	assert.Equal(t, expected2, jps[2].ToBuildRule(false, jps.ToGitMap(nil), nil, originalPackages))
	assert.Equal(t, expected3, jps[0].ToBuildRule(true, jps.ToGitMap(nil), nil, originalPackages))
}

func TestJsonPackagesSort(t *testing.T) {
	jps := jsonPackages{
		&jsonPackage{
			ImportPath: "gopkg.in/gcfg.v1",
			Deps: []string{
				"gopkg.in/gcfg.v1/scanner",
				"gopkg.in/gcfg.v1/token",
				"gopkg.in/gcfg.v1/types",
				"gopkg.in/warnings.v0",
			},
		},
		&jsonPackage{
			ImportPath: "gopkg.in/gcfg.v1/scanner",
			Deps: []string{
				"gopkg.in/gcfg.v1/token",
			},
		},
		&jsonPackage{
			ImportPath: "gopkg.in/gcfg.v1/token",
		},
		&jsonPackage{
			ImportPath: "gopkg.in/gcfg.v1/types",
		},
		&jsonPackage{
			ImportPath: "gopkg.in/warnings.v0",
		},
	}
	assert.False(t, jps.Less(0, 1))
	assert.True(t, jps.Less(1, 0))
	sort.Sort(jps)
	assert.Equal(t, "gopkg.in/gcfg.v1/token", jps[0].ImportPath)
	assert.Equal(t, "gopkg.in/gcfg.v1/scanner", jps[1].ImportPath)
	assert.Equal(t, "gopkg.in/gcfg.v1/types", jps[2].ImportPath)
	assert.Equal(t, "gopkg.in/warnings.v0", jps[3].ImportPath)
	assert.Equal(t, "gopkg.in/gcfg.v1", jps[4].ImportPath)
}

const samplePackages = `
{
	"Dir": "/home/pebers/git/please/src/core",
	"ImportPath": "core",
	"Name": "core",
	"Doc": "Tests on specific functions in build_target.go",
	"Target": "/home/pebers/git/please/pkg/linux_amd64/core.a",
	"Stale": true,
	"StaleReason": "cannot stat install target",
	"Root": "/home/pebers/git/please",
	"GoFiles": [
		"build_env.go",
		"build_label.go",
		"build_target.go",
		"cache.go",
		"config.go",
		"file_label.go",
		"glob.go",
		"graph.go",
		"lock.go",
		"package.go",
		"state.go",
		"test_results.go",
		"utils.go",
		"version.go"
	],
	"Imports": [
		"bytes",
		"cli",
		"context",
		"crypto/sha1",
		"encoding/base64",
		"encoding/gob",
		"fmt",
		"github.com/Workiva/go-datastructures/queue",
		"github.com/coreos/go-semver/semver",
		"gopkg.in/gcfg.v1",
		"gopkg.in/op/go-logging.v1",
		"io",
		"io/ioutil",
		"os",
		"os/exec",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"time"
	],
	"Deps": [
		"bufio",
		"bytes",
		"cli",
		"container/list",
		"context",
		"crypto",
		"crypto/sha1",
		"encoding",
		"encoding/base64",
		"encoding/binary",
		"encoding/gob",
		"errors",
		"fmt",
		"github.com/Workiva/go-datastructures/queue",
		"github.com/coreos/go-semver/semver",
		"github.com/dustin/go-humanize",
		"github.com/jessevdk/go-flags",
		"golang.org/x/crypto/ssh/terminal",
		"gopkg.in/gcfg.v1",
		"gopkg.in/gcfg.v1/scanner",
		"gopkg.in/gcfg.v1/token",
		"gopkg.in/gcfg.v1/types",
		"gopkg.in/op/go-logging.v1",
		"gopkg.in/warnings.v0",
		"hash",
		"internal/nettrace",
		"internal/race",
		"internal/singleflight",
		"io",
		"io/ioutil",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/rand",
		"net",
		"net/url",
		"os",
		"os/exec",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe"
	],
	"TestGoFiles": [
		"build_env_test.go",
		"build_label_test.go",
		"build_target_test.go",
		"config_test.go",
		"glob_test.go",
		"graph_test.go",
		"label_parse_test.go",
		"lock_test.go",
		"package_test.go",
		"state_test.go",
		"test_results_test.go",
		"utils_test.go"
	],
	"TestImports": [
		"cli",
		"context",
		"crypto/sha1",
		"encoding/base64",
		"fmt",
		"github.com/stretchr/testify/assert",
		"io/ioutil",
		"os",
		"strings",
		"syscall",
		"testing",
		"time"
	]
}
{
	"Dir": "/home/pebers/git/please/src/cli",
	"ImportPath": "cli",
	"Name": "cli",
	"Doc": "Package cli contains helper functions related to flag parsing and logging.",
	"Target": "/home/pebers/git/please/pkg/linux_amd64/cli.a",
	"Stale": true,
	"StaleReason": "cannot stat install target",
	"Root": "/home/pebers/git/please",
	"GoFiles": [
		"flags.go",
		"logging.go",
		"progress.go",
		"replacements.go",
		"window.go"
	],
	"Imports": [
		"bytes",
		"container/list",
		"fmt",
		"github.com/coreos/go-semver/semver",
		"github.com/dustin/go-humanize",
		"github.com/jessevdk/go-flags",
		"golang.org/x/crypto/ssh/terminal",
		"gopkg.in/op/go-logging.v1",
		"io",
		"net/url",
		"os",
		"path",
		"reflect",
		"regexp",
		"runtime",
		"strconv",
		"strings",
		"sync",
		"syscall",
		"time",
		"unsafe"
	],
	"Deps": [
		"bufio",
		"bytes",
		"container/list",
		"context",
		"encoding/binary",
		"errors",
		"fmt",
		"github.com/coreos/go-semver/semver",
		"github.com/dustin/go-humanize",
		"github.com/jessevdk/go-flags",
		"golang.org/x/crypto/ssh/terminal",
		"gopkg.in/op/go-logging.v1",
		"internal/nettrace",
		"internal/race",
		"internal/singleflight",
		"io",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/rand",
		"net",
		"net/url",
		"os",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe"
	],
	"TestGoFiles": [
		"flags_test.go",
		"logging_test.go"
	],
	"TestImports": [
		"github.com/stretchr/testify/assert",
		"strings",
		"testing",
		"time"
	]
}
{
	"Dir": "/home/pebers/git/please/src/parse",
	"ImportPath": "parse",
	"Name": "parse",
	"Doc": "Package responsible for parsing build files and constructing build targets \u0026 the graph.",
	"Target": "/home/pebers/git/please/pkg/linux_amd64/parse.a",
	"Stale": true,
	"StaleReason": "cannot stat install target",
	"Root": "/home/pebers/git/please",
	"GoFiles": [
		"builtin_rules.go",
		"parse_step.go",
		"suggest.go"
	],
	"CgoFiles": [
		"interpreter.go"
	],
	"CFiles": [
		"interpreter.c"
	],
	"HFiles": [
		"interpreter.h"
	],
	"CgoCFLAGS": [
		"--std=c99",
		"-Werror"
	],
	"CgoLDFLAGS": [
		"-ldl"
	],
	"Imports": [
		"C",
		"bytes",
		"compress/gzip",
		"core",
		"crypto/sha1",
		"fmt",
		"github.com/kardianos/osext",
		"gopkg.in/op/go-logging.v1",
		"io",
		"io/ioutil",
		"os",
		"path",
		"path/filepath",
		"runtime",
		"sort",
		"strings",
		"sync",
		"time",
		"unsafe",
		"update",
		"utils"
	],
	"Deps": [
		"archive/tar",
		"bufio",
		"bytes",
		"cli",
		"compress/bzip2",
		"compress/flate",
		"compress/gzip",
		"container/list",
		"context",
		"core",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/des",
		"crypto/dsa",
		"crypto/ecdsa",
		"crypto/elliptic",
		"crypto/hmac",
		"crypto/internal/cipherhw",
		"crypto/md5",
		"crypto/rand",
		"crypto/rc4",
		"crypto/rsa",
		"crypto/sha1",
		"crypto/sha256",
		"crypto/sha512",
		"crypto/subtle",
		"crypto/tls",
		"crypto/x509",
		"crypto/x509/pkix",
		"encoding",
		"encoding/asn1",
		"encoding/base64",
		"encoding/binary",
		"encoding/gob",
		"encoding/hex",
		"encoding/pem",
		"errors",
		"fmt",
		"github.com/Workiva/go-datastructures/queue",
		"github.com/coreos/go-semver/semver",
		"github.com/dustin/go-humanize",
		"github.com/jessevdk/go-flags",
		"github.com/kardianos/osext",
		"github.com/texttheater/golang-levenshtein/levenshtein",
		"golang.org/x/crypto/ssh/terminal",
		"gopkg.in/gcfg.v1",
		"gopkg.in/gcfg.v1/scanner",
		"gopkg.in/gcfg.v1/token",
		"gopkg.in/gcfg.v1/types",
		"gopkg.in/op/go-logging.v1",
		"gopkg.in/warnings.v0",
		"hash",
		"hash/crc32",
		"internal/nettrace",
		"internal/race",
		"internal/singleflight",
		"internal/syscall/unix",
		"io",
		"io/ioutil",
		"log",
		"log/syslog",
		"math",
		"math/big",
		"math/rand",
		"mime",
		"mime/multipart",
		"mime/quotedprintable",
		"net",
		"net/http",
		"net/http/httptrace",
		"net/http/internal",
		"net/textproto",
		"net/url",
		"os",
		"os/exec",
		"os/signal",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/internal/atomic",
		"runtime/internal/sys",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe",
		"update",
		"utils",
		"vendor/golang_org/x/crypto/chacha20poly1305",
		"vendor/golang_org/x/crypto/chacha20poly1305/internal/chacha20",
		"vendor/golang_org/x/crypto/curve25519",
		"vendor/golang_org/x/crypto/poly1305",
		"vendor/golang_org/x/net/http2/hpack",
		"vendor/golang_org/x/net/idna",
		"vendor/golang_org/x/net/lex/httplex",
		"vendor/golang_org/x/text/transform",
		"vendor/golang_org/x/text/unicode/norm",
		"vendor/golang_org/x/text/width"
	],
	"TestGoFiles": [
		"interpreter_test.go",
		"parse_step_test.go",
		"suggest_test.go"
	],
	"TestImports": [
		"core",
		"github.com/stretchr/testify/assert",
		"os",
		"strings",
		"testing",
		"unsafe"
	]
}

`
