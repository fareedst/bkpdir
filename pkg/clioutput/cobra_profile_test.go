// [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]
package clioutput

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestApplyToRoot_HelpDecoratedTTY_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := &cobra.Command{Use: "testcli", Short: "Short description"}
	root.AddCommand(&cobra.Command{Use: "sub", Short: "Subcommand"})
	var stdout, stderr bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&stdout, &stderr), WithEnv(map[string]string{}), WithTTYFunc(func(w io.Writer) bool { return true }))
	ApplyToRoot(root, p)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	help := stderr.String()
	if !strings.Contains(help, "\033") && !strings.Contains(help, "\x1b") {
		t.Fatalf("expected decorated help on TTY stderr, got: %q", help)
	}
}

func TestApplyToRoot_HelpPlainNoESC_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := &cobra.Command{Use: "testcli", Short: "Short description"}
	var stdout, stderr bytes.Buffer
	p := ResolvePolicy(PeeledFlags{Plain: true}, WithWriters(&stdout, &stderr), WithEnv(map[string]string{}), WithTTYFunc(func(w io.Writer) bool { return true }))
	ApplyToRoot(root, p)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	help := stderr.String()
	if strings.Contains(help, "\033") || strings.Contains(help, "\x1b") {
		t.Fatalf("plain help must not contain ESC: %q", help)
	}
}

func TestApplyToRoot_HelpNonTTYPlain_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := &cobra.Command{Use: "testcli", Short: "Short description"}
	var stdout, stderr bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&stdout, &stderr), WithEnv(map[string]string{}), WithTTYFunc(func(io.Writer) bool { return false }))
	ApplyToRoot(root, p)
	root.SetArgs([]string{"--help"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	help := stderr.String()
	if strings.Contains(help, "\033") || strings.Contains(help, "\x1b") {
		t.Fatalf("non-TTY help must be plain: %q", help)
	}
}

func TestCaptureHelp_ReapplySameTree_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	root := &cobra.Command{Use: "cap", Short: "Cap test"}
	var stdout, stderr bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&stdout, &stderr), WithEnv(map[string]string{}), WithTTYFunc(func(io.Writer) bool { return true }))
	stdout1, stderr1, err := CaptureHelp(root, p, []string{"--help"})
	if err != nil {
		t.Fatal(err)
	}
	stdout2, stderr2, err := CaptureHelp(root, p, []string{"--help"})
	if err != nil {
		t.Fatal(err)
	}
	if stdout1 != stdout2 || stderr1 != stderr2 {
		t.Fatalf("CaptureHelp should be stable on reapply\nfirst=%q\nsecond=%q", stderr1, stderr2)
	}
}
