// Package remote implements a helper for 'go get' to facilitate breaking it up
// into separate plz targets.
package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"sort"
	"strings"
	"sync"

	"gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("remote")

const template = `go_remote_library(
    name = '%s',
    get = '%s',
    revision = '%s',
    deps = [
        '%s',
    ],
)
`
const noDepsTemplate = `go_remote_library(
    name = '%s',
    get = '%s',
    revision = '%s',
)
`

// FetchLibraries uses 'go get' to fetch a series of libraries, and generate either a sequence of
// build rules describing them or a pithy description of them which can be parsed back
// into BUILD rules later. The BUILD rules generated by the former re-invoke this using the latter
// format to determine what exactly to build and how.
func FetchLibraries(gotool string, shortFormat bool, packages ...string) (string, error) {
	if out, err := goCommand(gotool, "get", "-d", packages...); err != nil {
		return "", fmt.Errorf("%s: %s", err, string(out))
	}
	packageData, err := goList(gotool, packages...)
	if err != nil {
		return "", err
	}
	// This gives us all their dependencies. go get might have fetched some others that we
	// don't know about, so we ask go list to re-describe them all to work out which are
	// system or not.
	packageData, err = goList(gotool, packageData.UniqueDeps()...)
	if err != nil {
		return "", err
	}
	// Now build up the response.
	if shortFormat {
		// This orders the dependencies such that dependors come after dependees, which
		// is important when we generate the build rules later.
		sort.Sort(packageData)
		var buf bytes.Buffer
		m := packageData.ToMap()
		for _, pkg := range packageData {
			if !pkg.Standard {
				buf.WriteString(pkg.ToShortFormatString(m))
			}
		}
		return buf.String(), nil
	}
	if err := packageData.AnnotateGitURLs(); err != nil {
		return "", err
	}
	m := packageData.ToGitMap()
	out := []string{}
	for _, pkg := range m {
		out = append(out, pkg.ToBuildRule(m))
	}
	sort.Strings(out)
	return strings.Join(out, "\n"), nil
}

// goCommand runs a Go command and returns its output.
func goCommand(gotool string, command, flag string, packages ...string) ([]byte, error) {
	if !strings.HasPrefix(gotool, "/") {
		path, err := exec.LookPath(gotool)
		if err != nil {
			return nil, err
		}
		gotool = path
	}
	log.Debug("Running %s %s %s %s...", gotool, command, flag, strings.Join(packages, " "))
	args := append([]string{command, flag}, packages...)
	cmd := exec.Command(gotool, args...)
	return cmd.Output()
}

// goList runs "go list -json" on the given packages and parses it to a struct.
func goList(gotool string, packages ...string) (jsonPackages, error) {
	out, err := goCommand(gotool, "list", "-json", packages...)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(out))
	}
	packageData := jsonPackages{}
	return packageData, packageData.FromJSON(out)
}

// A jsonPackage is a minimal copy of go list's builtin struct definition.
// Note that we don't support every possible feature here, only those that map to Please.
type jsonPackage struct {
	Dir          string // N.B. absolute path.
	Root         string
	ImportPath   string
	Target       string
	Standard     bool
	GoFiles      []string
	CgoFiles     []string
	CFiles       []string
	HFiles       []string
	CgoCFLAGS    []string
	CgoLDFLAGS   []string
	CgoPkgConfig []string // TODO(pebers): add support for this to cgo_library
	Imports      []string
	Deps         []string

	// GitURL is not in the upstream structure. We annotate it ourselves later.
	GitURL      string
	Revision    string
	RepoImports map[string]bool
}

// ToShortFormatString returns a short delimited string format that Please will re-parse later
// to create a build rule from.
func (jp *jsonPackage) ToShortFormatString(packages map[string]*jsonPackage) string {
	comma := func(s []string) string { return strings.Join(s, ",") }
	caret := func(s []string) string { return strings.Join(s, "^") }

	name := strings.Replace(strings.Replace(jp.ImportPath, "/", "_", -1), ".", "_", -1)
	dir := jp.trimRoot(jp.Dir)
	gofiles := comma(jp.GoFiles)
	deps := comma(jp.deps(packages))
	if len(jp.CgoFiles) == 0 {
		return fmt.Sprintf("%s|%s|%s|%s\n", name, dir, gofiles, deps)
	}
	// Cgo packages need quite a bit more information.
	cgofiles := comma(jp.CgoFiles)
	cfiles := comma(jp.CFiles)
	hfiles := comma(jp.HFiles)
	cflags := caret(jp.CgoCFLAGS)
	ldflags := caret(jp.CgoLDFLAGS)
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s\n", name, dir, gofiles, cgofiles, cfiles, hfiles, cflags, ldflags, deps)
}

// ToBuildRule returns a build rule representation suitable for copying into a BUILD file.
func (jp *jsonPackage) ToBuildRule(packages map[string]*jsonPackage) string {
	name := repoNameToRuleName(jp.GitURL)
	deps := jp.repoDeps(packages)
	if len(deps) == 0 {
		return fmt.Sprintf(noDepsTemplate, name, jp.GitURL, jp.Revision)
	}
	return fmt.Sprintf(template, name, jp.GitURL, jp.Revision, strings.Join(deps, "',\n        '"))
}

// trimRoot strips the root from the given string.
func (jp *jsonPackage) trimRoot(s string) string {
	return strings.TrimLeft(strings.TrimPrefix(s, jp.Root), "/")
}

// paths updates a sequence of paths with the given prefix.
func (jp *jsonPackage) paths(dir string, ps []string) []string {
	for i, p := range ps {
		ps[i] = path.Join(dir, p)
	}
	return ps
}

// deps returns the non-system dependencies for this package.
func (jp *jsonPackage) deps(packages map[string]*jsonPackage) []string {
	ret := make([]string, 0, len(jp.Imports))
	for _, imp := range jp.Imports {
		if pkg, present := packages[imp]; present && !pkg.Standard {
			ret = append(ret, strings.Replace(imp, "/", "_", -1))
		}
	}
	return ret
}

// repoDeps returns the non-system dependencies that we have build rules for.
func (jp *jsonPackage) repoDeps(packages map[string]*jsonPackage) []string {
	ret := make([]string, 0, len(jp.RepoImports))
	for dep := range jp.RepoImports {
		if pkg, present := packages[dep]; present && !pkg.Standard {
			if dep != jp.GitURL {
				ret = append(ret, ":"+repoNameToRuleName(dep))
			}
		}
	}
	sort.Strings(ret)
	return ret
}

// repoNameToRuleName converts a git repo name into a name suitable for a build rule.
func repoNameToRuleName(repoName string) string {
	return strings.TrimSuffix(repoName[strings.LastIndex(repoName, "/")+1:], ".git")
}

// FindGitURL finds the upstream Git URL of this package.
func (jp *jsonPackage) AnnotateGitURL() error {
	log.Debug("Running git config --get remote.origin.url in %s...", jp.Dir)
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = jp.Dir
	out, err := cmd.Output()
	if err != nil {
		// We need a bit of verbosity here so we don't just get 'exit status 1'
		return fmt.Errorf("%s in %s: %s", err, jp.Dir, string(out))
	}
	// Strip https:// prefix for more natural Go paths. We can assume it again later.
	jp.GitURL = strings.TrimSpace(strings.TrimPrefix(string(out), "https://"))
	log.Debug("Running %s in %s...", "git log -n 1 --pretty=format:'%H'", jp.Dir)
	cmd = exec.Command("git", "log", "-n", "1", "--pretty=format:'%H'")
	cmd.Dir = jp.Dir
	out, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s: %s", err, jp.Dir, string(out))
	}
	jp.Revision = strings.Trim(string(out), "'")
	return nil
}

// HasDep returns true if this package has a dependency on the given package.
func (jp *jsonPackage) HasDep(dep string) bool {
	for _, d := range jp.Deps {
		if d == dep {
			return true
		}
	}
	return false
}

type jsonPackages []*jsonPackage

// UniqueDeps returns the unique set of deps from a set of packages, including the packages themselves..
func (jps jsonPackages) UniqueDeps() []string {
	m := map[string]struct{}{}
	for _, jp := range jps {
		for _, dep := range jp.Deps {
			m[dep] = struct{}{}
		}
		m[jp.ImportPath] = struct{}{}
	}
	ret := make([]string, 0, len(m))
	for pkg := range m {
		ret = append(ret, pkg)
	}
	sort.Strings(ret)
	return ret
}

// ToMap converts the jsonPackages slice to a map of package import path -> package.
func (jps jsonPackages) ToMap() map[string]*jsonPackage {
	m := make(map[string]*jsonPackage, len(jps))
	for _, jp := range jps {
		m[jp.ImportPath] = jp
	}
	return m
}

// ToGitMap converts the jsonPackages slice to a map of git repo -> representative package.
func (jps jsonPackages) ToGitMap() map[string]*jsonPackage {
	byImport := jps.ToMap()
	m := map[string]*jsonPackage{}
	for _, jp := range jps {
		if jp.GitURL == "" {
			continue
		}
		p, present := m[jp.GitURL]
		if !present {
			m[jp.GitURL] = jp
			p = jp
			p.RepoImports = map[string]bool{}
		}
		for _, imp := range jp.Imports {
			if existing, present := byImport[imp]; present {
				p.RepoImports[existing.GitURL] = true
			}
		}
	}
	return m
}

// AnnotateGitURLs attempts to find the Git URL for each package.
func (jps jsonPackages) AnnotateGitURLs() error {
	var err error
	var wg sync.WaitGroup
	wg.Add(len(jps))
	for i, jp := range jps {
		go func(i int, jp *jsonPackage) {
			if !jp.Standard {
				if e := jp.AnnotateGitURL(); e != nil {
					err = e
				}
			}
			wg.Done()
		}(i, jp)
	}
	wg.Wait()
	return err
}

// FromJSON loads this set of packages from JSON.
// This is not as easy as you'd think since the output for multiple packages is not valid JSON -
// it's a sequence of top-level JSON objects one after another.
func (jps *jsonPackages) FromJSON(data []byte) error {
	d := append([]byte{'['}, bytes.Replace(data, []byte("}\n{"), []byte("},{"), -1)...)
	d = append(d, ']')
	return json.Unmarshal(d, jps)
}

func (jps jsonPackages) Len() int { return len(jps) }
func (jps jsonPackages) Less(i, j int) bool {
	return !jps[i].HasDep(jps[j].ImportPath) && (jps[j].HasDep(jps[i].ImportPath) || jps[i].ImportPath < jps[j].ImportPath)
}
func (jps jsonPackages) Swap(i, j int) {
	jps[i], jps[j] = jps[j], jps[i]
}
