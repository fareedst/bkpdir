//go:build !windows

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellCLIOutputContract_MigratedScripts_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := mustRepoRoot(t)
	scripts := []struct {
		name    string
		rel     string
		args    []string
		contain string
	}{
		{"run-mutation-pilot", "scripts/run-mutation-pilot.sh", []string{"--help"}, "run-mutation-pilot"},
		{"run-pseudocode-layer-b-batch", "scripts/run-pseudocode-layer-b-batch.sh", []string{"--help"}, "run-pseudocode-layer-b-batch"},
		{"run-layer-b-ci-enforcement", "scripts/run-layer-b-ci-enforcement.sh", []string{"--help"}, "run-layer-b-ci-enforcement"},
		{"token-navigate", "scripts/token-navigate.sh", []string{"help"}, "Token Navigation"},
		{"validate-test-tokens", "scripts/validate-test-tokens.sh", []string{"--help"}, "Usage:"},
		{"generate-missing-tokens", "scripts/generate-missing-tokens.sh", []string{"--help"}, "Usage:"},
		{"enhance-cross-reference-validation", "scripts/enhance-cross-reference-validation.sh", []string{"--help"}, "Usage:"},
	}
	for _, sc := range scripts {
		sc := sc
		t.Run(sc.name, func(t *testing.T) {
			path := filepath.Join(root, sc.rel)
			cmd := exec.Command("bash", append([]string{path}, sc.args...)...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "NO_COLOR=1")
			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				t.Fatalf("run %s: %v stderr=%q", sc.name, err, stderr.String())
			}
			if stdout.Len() != 0 {
				t.Errorf("%s: stdout should be empty on help, got %q", sc.name, stdout.String())
			}
			if !strings.Contains(stderr.String(), sc.contain) {
				t.Errorf("%s: expected stderr containing %q, got %q", sc.name, sc.contain, stderr.String())
			}
		})
	}
}

func TestShellCLIOutputContract_StdoutHelpEnv_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := mustRepoRoot(t)
	path := filepath.Join(root, "scripts/run-mutation-pilot.sh")
	cmd := exec.Command("bash", path, "--help")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "TIED_CLI_HELP_STDOUT=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr should be empty, got %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), "run-mutation-pilot") {
		t.Errorf("expected help on stdout, got %q", stdout.String())
	}
}

func mustRepoRoot(t *testing.T) string {
	t.Helper()
	return repoRoot(t)
}
