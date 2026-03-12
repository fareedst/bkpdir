// [IMPL-LIST_FORMAT_SAFETY] [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] [REQ-OUT_002]
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFormatListArchiveSimple_Comprehensive tests the simplified list formatting implementation
// with comprehensive test cases covering all scenarios
func TestFormatListArchiveSimple_Comprehensive(t *testing.T) {
	// Setup: Create a test file with known content
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-archive.zip")
	content := []byte("test archive content for size calculation - this is a longer string to make the file size more visible and testable")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Get file stats for validation
	statInfo, err := formatter.GatherFileStatInfo(testFile)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}

	creationTime := "2025-07-13 19:31:42"

	tests := []struct {
		name                string
		formatListArchive   string
		templateListArchive string
		expectedContains    []string
		expectedNotContains []string
		description         string
	}{
		{
			name:                "FormatListArchive_with_template_placeholders",
			formatListArchive:   "#{path} (size: #{size_human})\n",
			templateListArchive: "",
			expectedContains:    []string{testFile, "size:", statInfo.SizeHuman},
			expectedNotContains: []string{"#{path}", "#{size_human}", "%!{"},
			description:         "FormatListArchive with template placeholders should work",
		},
		{
			name:                "FormatListArchive_template_only",
			formatListArchive:   "#{path} (size: #{size_human})\n",
			templateListArchive: "",
			expectedContains:    []string{testFile, "size:", statInfo.SizeHuman},
			expectedNotContains: []string{"#{path}", "#{size_human}", "%!{"},
			description:         "FormatListArchive with template placeholders only should work",
		},
		{
			name:                "FormatListArchive_with_creation_time",
			formatListArchive:   "#{path} (created: #{creation_time})\n",
			templateListArchive: "",
			expectedContains:    []string{testFile, creationTime},
			expectedNotContains: []string{"#{path}", "#{creation_time}", "%!{"},
			description:         "FormatListArchive with creation_time placeholder should work",
		},
		{
			name:                "TemplateListArchive_fallback",
			formatListArchive:   "",
			templateListArchive: "#{path} (size: #{size_human})\n",
			expectedContains:    []string{testFile, "size:", statInfo.SizeHuman},
			expectedNotContains: []string{"#{path}", "#{size_human}", "%!{"},
			description:         "TemplateListArchive should be used when FormatListArchive is empty",
		},
		{
			name:                "Default_template_fallback",
			formatListArchive:   "",
			templateListArchive: "",
			expectedContains:    []string{testFile, "size:"},
			expectedNotContains: []string{"#{path}", "#{size_human}", "%!{"},
			description:         "Default template format should be used when both are empty",
		},
		{
			name:                "FormatListArchive_without_template_falls_back",
			formatListArchive:   "plain text without placeholders\n",
			templateListArchive: "#{path} (size: #{size_human})\n",
			expectedContains:    []string{testFile, "size:"},
			expectedNotContains: []string{"#{path}", "#{size_human}", "%!{"},
			description:         "When FormatListArchive doesn't contain #{, should fall back to TemplateListArchive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := DefaultConfig()
			cfg.FormatListArchive = tt.formatListArchive
			cfg.TemplateListArchive = tt.templateListArchive
			formatterInstance := NewFormatterAdapter(cfg)

			result := formatListArchiveSimple(cfg, formatterInstance, testFile, creationTime)

			t.Logf("Test: %s", tt.description)
			t.Logf("Config FormatListArchive: %q", cfg.FormatListArchive)
			t.Logf("Config TemplateListArchive: %q", cfg.TemplateListArchive)
			t.Logf("Test FormatListArchive: %q", tt.formatListArchive)
			t.Logf("Test TemplateListArchive: %q", tt.templateListArchive)
			t.Logf("Result: %q", result)

			// Verify expected content
			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected result to contain %q, got %q", expected, result)
				}
			}

			// Verify unexpected content (errors, unprocessed placeholders)
			for _, notExpected := range tt.expectedNotContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("BUG: Result contains unexpected content %q, got %q", notExpected, result)
				}
			}

			// CRITICAL: Verify no Go template errors
			if strings.Contains(result, "%!{") {
				t.Errorf("CRITICAL BUG: Result contains Go template error format, got %q", result)
			}

			// CRITICAL: Verify no unprocessed template placeholders
			if strings.Contains(result, "#{") {
				t.Errorf("CRITICAL BUG: Result contains unprocessed template placeholder, got %q", result)
			}
		})
	}
}

// TestFormatListArchiveSimple_FileStats tests file statistics gathering and inclusion
func TestFormatListArchiveSimple_FileStats(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-archive.zip")
	content := []byte("test content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	statInfo, err := formatter.GatherFileStatInfo(testFile)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}

	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human}, bytes: #{size})\n"
	formatterInstance := NewFormatterAdapter(cfg)

	result := formatListArchiveSimple(cfg, formatterInstance, testFile, "2025-07-13 19:31:42")

	t.Logf("File size: %d bytes", statInfo.Size)
	t.Logf("File size_human: %s", statInfo.SizeHuman)
	t.Logf("Result: %q", result)

	// Verify size_human is included
	if !strings.Contains(result, statInfo.SizeHuman) {
		t.Errorf("Expected result to contain size_human %q, got %q", statInfo.SizeHuman, result)
	}

	// Verify size is included
	if !strings.Contains(result, fmt.Sprintf("%d", statInfo.Size)) {
		t.Errorf("Expected result to contain size %d, got %q", statInfo.Size, result)
	}
}

// TestFormatListArchiveSimple_MissingFile tests behavior when file doesn't exist
func TestFormatListArchiveSimple_MissingFile(t *testing.T) {
	nonExistentFile := "/nonexistent/path/to/file.zip"
	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	formatterInstance := NewFormatterAdapter(cfg)

	result := formatListArchiveSimple(cfg, formatterInstance, nonExistentFile, "2025-07-13 19:31:42")

	t.Logf("Result for non-existent file: %q", result)

	// Should still format correctly with "unknown" for size_human
	if !strings.Contains(result, nonExistentFile) {
		t.Errorf("Expected result to contain file path, got %q", result)
	}

	// Should contain "unknown" for size_human when file doesn't exist
	if !strings.Contains(result, "unknown") {
		t.Errorf("Expected result to contain 'unknown' for missing file stats, got %q", result)
	}

	// Should not contain unprocessed placeholders (old %{ syntax or new #{ syntax)
	if strings.Contains(result, "%{") {
		t.Errorf("BUG: Result contains unprocessed old-style template placeholder, got %q", result)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed template placeholder, got %q", result)
	}
}
