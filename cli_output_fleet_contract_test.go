//go:build !windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFleetCLIOutputContract_CobraTools_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		pkg       string
		contains  string
	}{
		{"token-suggester", "./cmd/token-suggester", "token-suggester"},
		{"ai-validation", "./cmd/ai-validation", "ai-validation"},
		{"realtime-validator", "./cmd/realtime-validator", "realtime-validator"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			exe := buildGoPkgBinary(t, tc.pkg, tc.name)
			assertHelpDefaultStderr(t, exe, []string{"--help"}, tc.contains)
		})
	}
}

func TestFleetCLIOutputContract_SpecctlUsage_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	exe := buildGoPkgBinary(t, "./cmd/specctl", "specctl")
	cmd := exec.Command(exe, "not-a-subcommand")
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err == nil {
		t.Fatal("expected non-zero exit for unknown subcommand")
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout should be empty on usage error, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "usage:") {
		t.Errorf("expected usage on stderr, got %q", stderr.String())
	}
}

func buildGoPkgBinary(t *testing.T, pkg, name string) string {
	t.Helper()
	root := repoRoot(t)
	dir := t.TempDir()
	exe := filepath.Join(dir, name)
	cmd := exec.Command("go", "build", "-o", exe, pkg)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s: %v\n%s", pkg, err, out)
	}
	return exe
}
