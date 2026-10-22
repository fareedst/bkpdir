//go:build !windows

// [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellCLIOutputProfile_PipedHelpNoANSI_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := mustRepoRoot(t)
	path := filepath.Join(root, "scripts/run-mutation-pilot.sh")
	cmd := exec.Command("bash", path, "--help")
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "NO_COLOR=1")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	out := stderr.String()
	if strings.Contains(out, "\033") || strings.Contains(out, "\x1b") {
		t.Fatalf("piped/plain shell help must not contain ESC: %q", out)
	}
}
