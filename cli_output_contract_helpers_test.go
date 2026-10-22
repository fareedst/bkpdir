//go:build !windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var cachedModuleRoot string

// cacheModuleRootForTests records the module root before tests that os.Chdir without restore.
func cacheModuleRootForTests() {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	if root, ok := moduleRootFromDir(wd); ok {
		cachedModuleRoot = root
	}
}

func moduleRootFromDir(start string) (string, bool) {
	for d := start; ; {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d, true
		}
		parent := filepath.Dir(d)
		if parent == d {
			return "", false
		}
		d = parent
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	if cachedModuleRoot != "" {
		return cachedModuleRoot
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if root, ok := moduleRootFromDir(wd); ok {
		return root
	}
	t.Fatal("go.mod not found")
	return ""
}

func assertHelpDefaultStderr(t *testing.T, exe string, helpArgs []string, stderrContains string) {
	t.Helper()
	cmd := exec.Command(exe, helpArgs...)
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v: run failed: %v stderr=%q", exe, helpArgs, err, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("%s %v: stdout should be empty on default help, got %q", exe, helpArgs, stdout.String())
	}
	if stderrContains != "" && !strings.Contains(stderr.String(), stderrContains) {
		t.Errorf("%s %v: expected stderr to contain %q, got %q", exe, helpArgs, stderrContains, stderr.String())
	}
}
