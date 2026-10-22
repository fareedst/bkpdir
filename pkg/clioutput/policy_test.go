package clioutput

import (
	"bytes"
	"io"
	"testing"
)

func TestResolvePolicy_DefaultHelpToStderr_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	var out, err bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&out, &err), WithTTYFunc(func(io.Writer) bool { return true }))
	if p.HelpDest != HelpToStderr {
		t.Fatalf("HelpDest = %v, want stderr", p.HelpDest)
	}
	if p.HelpWriter() != &err {
		t.Fatal("expected help writer stderr buffer")
	}
}

func TestResolvePolicy_StdoutOverride_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	var out, err bytes.Buffer
	p := ResolvePolicy(PeeledFlags{Stdout: true}, WithWriters(&out, &err), WithTTYFunc(func(io.Writer) bool { return true }))
	if p.HelpDest != HelpToStdout {
		t.Fatalf("HelpDest = %v, want stdout", p.HelpDest)
	}
	if p.HelpWriter() != &out {
		t.Fatal("expected help writer stdout buffer")
	}
}

func TestResolvePolicy_PlainFromEnv_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	p := ResolvePolicy(PeeledFlags{}, WithEnv(map[string]string{"TIED_CLI_PLAIN": "1"}), WithTTYFunc(func(io.Writer) bool { return true }))
	if !p.Plain {
		t.Fatal("expected plain from env")
	}
}

func TestPeelTextureArgs_RemovesGlobals_REQ_BKPDIR_CLI_OUTPUT_STANDARD(t *testing.T) {
	rem, flags := PeelTextureArgs([]string{"--plain", "create", "--stdout", "-p"})
	if !flags.Plain || !flags.Stdout {
		t.Fatalf("flags = %+v", flags)
	}
	want := []string{"create"}
	if len(rem) != len(want) || rem[0] != want[0] {
		t.Fatalf("remaining = %v, want %v", rem, want)
	}
}
