//go:build !windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIOutputContract_RootHelpDefaultStderr_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	exe := buildBkpdirTestBinary(t)
	cmd := exec.Command(exe, "--help")
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("run: %v stderr=%q", err, stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should be empty on default help, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "bkpdir") {
		t.Errorf("expected help on stderr, got %q", stderr.String())
	}
}

func TestCLIOutputContract_StdoutHelpOverride_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	exe := buildBkpdirTestBinary(t)
	cmd := exec.Command(exe, "--stdout", "--help")
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr should be empty on successful stdout help, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "bkpdir") {
		t.Errorf("expected help on stdout, got %q", stdout.String())
	}
}

func buildBkpdirTestBinary(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	dir := t.TempDir()
	exe := filepath.Join(dir, "bkpdir")
	cmd := exec.Command("go", "build", "-o", exe, ".")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build bkpdir: %v\n%s", err, out)
	}
	return exe
}
