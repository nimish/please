// package main implements various tools that Please finds useful for Go.
//
// Firstly, it implements a code templater for Go tests.
// This is essentially equivalent to what 'go test' does but it lifts some restrictions
// on file organisation and allows the code to be instrumented for coverage as separate
// build targets rather than having to repeat it for every test.
package main

import (
	"fmt"

	"gopkg.in/op/go-logging.v1"

	"cli"
	"tools/please_go_tool/remote"
	"tools/please_go_tool/testmain"
)

var log = logging.MustGetLogger("plz_go_test")

var opts = struct {
	Usage     string
	Verbosity int    `short:"v" long:"verbose" default:"2" description:"Verbosity of output (higher number = more output, default 2 -> notice, warnings and errors only)"`
	Go        string `short:"g" long:"go" default:"go" description:"Go binary to run"`

	TestMain struct {
		Dir     string   `short:"d" long:"dir" description:"Directory to search for Go package files for coverage"`
		Exclude []string `short:"x" long:"exclude" default:"third_party/go" description:"Directories to exclude from search"`
		Output  string   `short:"o" long:"output" description:"Output filename" required:"true"`
		Package string   `short:"p" long:"package" description:"Package containing this test" env:"PKG"`
		Args    struct {
			Sources []string `positional-arg-name:"sources" description:"Test source files" required:"true"`
		} `positional-args:"true" required:"true"`
	} `command:"testmain" description:"Templates a test main."`

	Remote struct {
		ShortFormat bool `short:"s" long:"short_format" description:"Prints a shorter format that is used for deriving individual generated rules."`
		Args        struct {
			Packages []string `positional-arg-name:"packages" description:"Packages to fetch" required:"true"`
		} `positional-args:"true" required:"true"`
	} `command:"remote" description:"Gets and prints some remote libraries."`
}{
	Usage: `
please_go_tool implements various tools that Please finds useful for building Go code.

Firstly, it implements a code templater for Go tests.
This is essentially equivalent to what 'go test' does but it lifts some restrictions
on file organisation and allows the code to be instrumented for coverage as separate
build targets rather than having to repeat it for every test.

It also implements a dependency fetcher, akin to a Pleaseish version of go get, which
generates a bunch of build rules for separate dependencies (and at runtime is re-invoked to
find source files etc).
`,
}

func main() {
	parser := cli.ParseFlagsOrDie("please_go_tool", "7.9.0", &opts)
	cli.InitLogging(opts.Verbosity)

	if parser.Active.Name == "testmain" {
		coverVars, err := testmain.FindCoverVars(opts.TestMain.Dir, opts.TestMain.Exclude, opts.TestMain.Args.Sources)
		if err != nil {
			log.Fatalf("Error scanning for coverage: %s", err)
		}
		if err = testmain.WriteTestMain(opts.TestMain.Package, testmain.IsVersion18(opts.Go), opts.TestMain.Args.Sources, opts.TestMain.Output, coverVars); err != nil {
			log.Fatalf("Error writing test main: %s", err)
		}
	} else if parser.Active.Name == "remote" {
		s, err := remote.FetchLibraries(opts.Go, opts.Remote.ShortFormat, opts.Remote.Args.Packages...)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		fmt.Print(s)
	}
}
