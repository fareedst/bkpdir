package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

// [IMPL:LIST_FORMAT_SAFETY] TestFormatListArchive_IMPL_LIST_FORMAT_SAFETY verifies template-style placeholder formatting works and is referenced by the semantic token
func TestFormatListArchive_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	// create temp file to have a real size
	f, err := ioutil.TempFile("", "bkpdir-test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if _, err := f.WriteString("hello world"); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	cfg := DefaultConfig()
	// Use template-style placeholders
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	formatter := &OutputFormatter{cfg: cfg}

	creationTime := "2025-12-18 14:02:04"
	out := formatter.FormatListArchiveWithExtraction(f.Name(), creationTime)

	if out == "" {
		t.Fatalf("expected non-empty output")
	}
	if !contains(out, f.Name()) {
		t.Fatalf("expected output to contain path %q; got %q", f.Name(), out)
	}
	// size_human should be replaced, not left as placeholder
	if contains(out, "#{size_human}") {
		t.Fatalf("placeholder #{size_human} was not replaced: %q", out)
	}
}

// [IMPL:LIST_FORMAT_SAFETY] TestPrintfStyleListArchive_IMPL_LIST_FORMAT_SAFETY verifies printf-style formats still work and do not trigger EXTRA output
func TestPrintfStyleListArchive_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FormatListArchive = "%s (created: %s)\n"
	formatter := &OutputFormatter{cfg: cfg}

	path := "some/path/file.txt"
	creationTime := "2025-12-18 14:02:04"
	out := formatter.FormatListArchiveWithExtraction(path, creationTime)
	expected := fmt.Sprintf(cfg.FormatListArchive, path, creationTime)
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (indexOf(s, sub) >= 0)))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
