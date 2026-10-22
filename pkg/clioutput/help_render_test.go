// [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]
package clioutput

import (
	"os"
	"strings"
	"testing"
)

func loadHelpFixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestRenderHelp_Plain_NoANSI_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	in := loadHelpFixture(t, "help_minimal_root.txt")
	out := RenderHelp(in, RenderOptions{Plain: true})
	if strings.Contains(out, "\033") || strings.Contains(out, "\x1b") {
		t.Fatalf("plain help must not contain ESC: %q", out)
	}
	if strings.Contains(out, "▶") {
		t.Fatal("plain help must not contain section arrow")
	}
	if !strings.Contains(out, strings.Repeat("-", 56)) {
		t.Fatal("plain help expected 56-char ASCII rule")
	}
}

func TestRenderHelp_Decorated_SectionAndRule_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	in := loadHelpFixture(t, "help_minimal_root.txt")
	out := RenderHelp(in, RenderOptions{Plain: false})
	if !strings.Contains(out, "\033[1m") {
		t.Fatal("decorated help expected bold sequences")
	}
	if !strings.Contains(out, "\033[36m") {
		t.Fatal("decorated help expected cyan section treatment")
	}
	if !strings.Contains(out, "▶") {
		t.Fatal("decorated help expected section arrow")
	}
	if !strings.Contains(out, "─") {
		t.Fatal("decorated help expected unicode rule character")
	}
}

func TestRenderHelp_Decorated_CommandRow_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	in := loadHelpFixture(t, "help_minimal_root.txt")
	out := RenderHelp(in, RenderOptions{Plain: false})
	if !strings.Contains(out, "backup") || !strings.Contains(out, "Backup a single file") {
		t.Fatalf("command row content missing: %q", out)
	}
	// Bold wraps command name in decorated mode.
	if !strings.Contains(out, "\033[1mbackup\033[0m") {
		t.Fatalf("expected bold command name, got: %q", out)
	}
}

func TestRenderHelp_EmptyInput_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	if got := RenderHelp("", RenderOptions{Plain: true}); got != "" {
		t.Fatalf("empty plain: %q", got)
	}
	if got := RenderHelp("", RenderOptions{Plain: false}); got != "" {
		t.Fatalf("empty decorated: %q", got)
	}
}

func TestRenderHelp_PreservesUnknownLine_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	in := "Custom note line\nUsage:\n  cmd\n"
	out := RenderHelp(in, RenderOptions{Plain: true})
	if !strings.Contains(out, "Custom note line") {
		t.Fatalf("unknown line should be preserved: %q", out)
	}
}
