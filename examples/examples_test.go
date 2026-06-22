package examples_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var examples = []struct {
	name  string
	pkg   string
	deps  []string
	avoid []string
}{
	{
		name:  "minimal",
		pkg:   "./examples/minimal",
		avoid: []string{"go.uber.org/fx"},
	},
	{
		name: "fx",
		pkg:  "./examples/fx",
		deps: []string{"go.uber.org/fx"},
	},
}

func TestExamplesImportBoundary(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go binary not available")
	}

	root := repoRoot(t)

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			deps := listDeps(t, goBin, root, example.pkg)
			for _, want := range example.deps {
				if !hasDependency(deps, want) {
					t.Fatalf("%s example does not depend on %s:\n%s", example.name, want, deps)
				}
			}
			for _, avoid := range example.avoid {
				if hasDependency(deps, avoid) {
					t.Fatalf("%s example depends on %s:\n%s", example.name, avoid, deps)
				}
			}
		})
	}
}

func TestExamplesOptimizedBuildSizes(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go binary not available")
	}

	root := repoRoot(t)
	binDir := t.TempDir()

	var rows []exampleSize
	for _, example := range examples {
		deps := listDeps(t, goBin, root, example.pkg)
		rows = append(rows, buildOptimizedExample(t, goBin, root, binDir, example.name, example.pkg, hasDependency(deps, "go.uber.org/fx")))
	}

	t.Logf("\n%s", formatExampleSizes(rows))
}

func repoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	return filepath.Dir(wd)
}

type exampleSize struct {
	name    string
	pkg     string
	binary  string
	bytes   int64
	hasFx   bool
	goos    string
	goarch  string
	cgo     string
	ldflags string
}

func listDeps(t *testing.T, goBin, root, pkg string) string {
	t.Helper()

	cmd := exec.Command(goBin, "list", "-deps", pkg)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off")

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps %s failed: %v\n%s", pkg, err, out)
	}

	return string(out)
}

func buildOptimizedExample(t *testing.T, goBin, root, binDir, name, pkg string, hasFx bool) exampleSize {
	t.Helper()

	ldflags := "-s -w -buildid="
	out := filepath.Join(binDir, name)
	if runtime.GOOS == "windows" {
		out += ".exe"
	}

	cmd := exec.Command(
		goBin,
		"build",
		"-trimpath",
		"-buildvcs=false",
		"-ldflags", ldflags,
		"-o", out,
		pkg,
	)
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=0",
		"GOWORK=off",
	)

	buildOut, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("optimized build %s failed: %v\n%s", pkg, err, buildOut)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat optimized build %s failed: %v", pkg, err)
	}

	return exampleSize{
		name:    name,
		pkg:     pkg,
		binary:  out,
		bytes:   info.Size(),
		hasFx:   hasFx,
		goos:    runtime.GOOS,
		goarch:  runtime.GOARCH,
		cgo:     "0",
		ldflags: ldflags,
	}
}

func formatExampleSizes(rows []exampleSize) string {
	var b strings.Builder

	b.WriteString("Optimized example binary sizes\n")
	b.WriteString("Flags: CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags='-s -w -buildid='\n\n")
	b.WriteString("| Example | Package | GOOS/GOARCH | Fx dep | Bytes | KiB |\n")
	b.WriteString("| --- | --- | --- | --- | ---: | ---: |\n")
	for _, row := range rows {
		fmt.Fprintf(
			&b,
			"| %s | `%s` | %s/%s | %s | %d | %.1f |\n",
			row.name,
			row.pkg,
			row.goos,
			row.goarch,
			yesNo(row.hasFx),
			row.bytes,
			float64(row.bytes)/1024,
		)
	}

	return b.String()
}

func yesNo(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func hasDependency(deps string, path string) bool {
	for _, line := range strings.Split(deps, "\n") {
		if line == path || strings.HasPrefix(line, path+"/") {
			return true
		}
	}

	return false
}
