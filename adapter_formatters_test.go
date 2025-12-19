package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

// Tests for adapter paths (FormatterAdapter and AIFormatterAdapter)

// [IMPL:LIST_FORMAT_SAFETY] TestFormatterAdapter_PrintfStyle_IMPL_LIST_FORMAT_SAFETY ensures adapter fallback behavior for printf-style formats
func TestFormatterAdapter_PrintfStyle_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	// The simplified FormatterAdapter implementation does NOT interpret printf-style formats (%s)
	// and will fall back to template formatting (default or TemplateListArchive). Ensure safe fallback.
	cfg := DefaultConfig()
	cfg.FormatListArchive = "%s (created: %s)\n" // printf-style (not supported by simplified adapter)
	fa := NewFormatterAdapter(cfg)

	path := "some/path/file.txt"
	creationTime := "2025-12-18 14:02:04"
	out := fa.FormatListArchiveWithExtraction(path, creationTime)
	// Expect fallback to default template which uses size_human (unknown for non-existent path)
	expected := fmt.Sprintf("%s (size: %s)\n", path, "unknown")
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

// [IMPL:LIST_FORMAT_SAFETY] TestFormatterAdapter_TemplateStyle_IMPL_LIST_FORMAT_SAFETY verifies adapter replaces template placeholders
func TestFormatterAdapter_TemplateStyle_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	// create temp file to have real stats
	f, err := ioutil.TempFile("", "bkpdir-adapter-test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.WriteString("hello adapter"); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	fa := NewFormatterAdapter(cfg)

	out := fa.FormatListArchiveWithExtraction(f.Name(), "2025-12-18 14:02:04")
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
	if strings.Contains(out, "#{size_human}") {
		t.Fatalf("placeholder #{size_human} was not replaced: %q", out)
	}
}

// [IMPL:LIST_FORMAT_SAFETY] TestAIFormatterAdapter_TemplateStyle_IMPL_LIST_FORMAT_SAFETY verifies AI adapter replaces template placeholders
func TestAIFormatterAdapter_TemplateStyle_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	// create temp file to have real stats
	f, err := ioutil.TempFile("", "bkpdir-ai-adapter-test-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.WriteString("hello ai adapter"); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	ai := NewAIFormatterAdapter(cfg)

	out := ai.FormatListArchiveWithExtraction(f.Name(), "2025-12-18 14:02:04")
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
	if strings.Contains(out, "#{size_human}") {
		t.Fatalf("placeholder #{size_human} was not replaced by AI adapter: %q", out)
	}
}

// [IMPL:LIST_FORMAT_SAFETY] TestAIFormatterAdapter_TemplateFallback_IMPL_LIST_FORMAT_SAFETY verifies AI adapter falls back to TemplateListArchive
func TestAIFormatterAdapter_TemplateFallback_IMPL_LIST_FORMAT_SAFETY(t *testing.T) {
	// create temp file to have real stats
	f, err := ioutil.TempFile("", "bkpdir-ai-adapter-test-2-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.WriteString("hello ai adapter 2"); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	cfg := DefaultConfig()
	cfg.FormatListArchive = "" // not set
	cfg.TemplateListArchive = "#{path} (size: #{size_human})\n"
	ai := NewAIFormatterAdapter(cfg)

	out := ai.FormatListArchiveWithExtraction(f.Name(), "2025-12-18 14:02:04")
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
	if strings.Contains(out, "#{size_human}") {
		t.Fatalf("placeholder #{size_human} was not replaced by AI adapter fallback: %q", out)
	}
}
