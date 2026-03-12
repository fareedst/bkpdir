// This file is part of bkpdir
// [REQ-OUTPUT_FORMATTING] Validates output formatting features via tests
// [ARCH-OUTPUT_FORMATTING] Exercise output formatting architecture decisions
// [IMPL-DUAL_FORMATTING] Targets printf/template formatter implementation

package main

import (
	"bkpdir/pkg/formatter"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplateFormatter tests the template formatting functionality
func TestTemplateFormatter(t *testing.T) {
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("ExtractionWithFormatting", func(t *testing.T) {
		archivePath := "/archives/HOME-2024-06-01-12-00=main=abc123=test.zip"
		result := fa.FormatListArchiveWithExtraction(archivePath, "2024-06-01 12:00:00")
		if !strings.Contains(result, "HOME-2024-06-01-12-00=main=abc123=test.zip") {
			t.Errorf("Expected result to contain archive filename, got %q", result)
		}
	})

	t.Run("CustomTemplateHandling", func(t *testing.T) {
		// Test custom template format with placeholders
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "feature-branch",
			"note":   "custom note",
		}

		// Test placeholder replacement
		format := "Archive: #{path} on #{branch} - #{note}"
		result := fa.FormatWithPlaceholders(format, data)
		expected := "Archive: /test/archive.zip on feature-branch - custom note"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test with missing placeholders (should remain as-is)
		format = "Archive: #{path} unknown: #{missing}"
		result = fa.FormatWithPlaceholders(format, data)
		expected = "Archive: /test/archive.zip unknown: #{missing}"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Test that list output displays size correctly with #{size_human} placeholder
	t.Run("ListOutputWithSizeHuman", func(t *testing.T) {
		// Create a temporary file to get real file stats
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-archive.zip")

		// Create a test file with some content
		content := []byte("test archive content for size calculation")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Get file stats to verify size_human is populated
		statInfo, err := formatter.GatherFileStatInfo(testFile)
		if err != nil {
			t.Fatalf("Failed to gather file stats: %v", err)
		}

		// Verify size_human is not empty
		if statInfo.SizeHuman == "" {
			t.Fatalf("Expected size_human to be populated, got empty string")
		}

		// Test FormatListArchiveWithExtraction with real file
		creationTime := "2025-07-13 19:31:42"
		result := fa.FormatListArchiveWithExtraction(testFile, creationTime)

		// The result should contain the file path
		if !strings.Contains(result, testFile) {
			t.Errorf("Expected result to contain file path %q, got %q", testFile, result)
		}

		// CRITICAL: The result should NOT contain the literal #{size_human} placeholder
		if strings.Contains(result, "#{size_human}") {
			t.Errorf("BUG: Result contains literal #{size_human} placeholder, got %q", result)
		}

		// CRITICAL: The result should NOT contain Go template error format
		if strings.Contains(result, "%!{") {
			t.Errorf("BUG: Result contains Go template error format %%!{...}, got %q", result)
		}

		// The result should contain the actual size_human value or "unknown"
		if !strings.Contains(result, statInfo.SizeHuman) && !strings.Contains(result, "unknown") {
			t.Errorf("Expected result to contain size_human value %q or 'unknown', got %q", statInfo.SizeHuman, result)
		}

		// The result should contain "size:" from the template format
		if !strings.Contains(result, "size:") {
			t.Errorf("Expected result to contain 'size:' from template, got %q", result)
		}

		// Print the result for debugging
		t.Logf("File size_human: %q", statInfo.SizeHuman)
		t.Logf("Result: %q", result)
	})

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Comprehensive test for mixed placeholder format strings
	// This test validates that FormatListArchive correctly processes format strings with both
	// printf-style (%s) and template-style (#{size_human}) placeholders
	t.Run("FormatListArchive_MixedPlaceholders", func(t *testing.T) {
		// Create a temporary file to get real file stats
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-archive.zip")

		// Create a test file with known content to verify size calculation
		content := []byte("test archive content for size calculation - this is a longer string to make the file size more visible and testable")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Get file stats to verify size_human is populated
		statInfo, err := formatter.GatherFileStatInfo(testFile)
		if err != nil {
			t.Fatalf("Failed to gather file stats: %v", err)
		}

		cfg := DefaultConfig()
		cfg.FormatListArchive = "%s (size: #{size_human})\n"

		formatterAdapter := NewFormatterAdapter(cfg)
		creationTime := "2025-07-13 19:31:42"

		result := formatterAdapter.FormatListArchive(testFile, creationTime)

		t.Logf("Format string: %%s (size: #{size_human})\\n")
		t.Logf("Result: %q", result)

		// CRITICAL: The result should NOT contain the literal #{size_human} placeholder
		if strings.Contains(result, "#{size_human}") {
			t.Errorf("BUG: Result contains literal #{size_human} placeholder, got %q", result)
		}

		// CRITICAL: The result should NOT contain Go template error format
		if strings.Contains(result, "%!{") {
			t.Errorf("BUG: Result contains Go template error format %%!{...}, got %q", result)
		}

		// CRITICAL: The result should NOT contain any #{...} patterns
		if strings.Contains(result, "#{") {
			t.Errorf("BUG: Result contains unprocessed template placeholder pattern #{...}, got %q", result)
		}

		// The result should contain the file path
		if !strings.Contains(result, testFile) {
			t.Errorf("Expected result to contain file path %q, got %q", testFile, result)
		}

		// The result should contain the actual size_human value or "unknown"
		if !strings.Contains(result, statInfo.SizeHuman) && !strings.Contains(result, "unknown") {
			t.Errorf("Expected result to contain size_human value %q or 'unknown', got %q", statInfo.SizeHuman, result)
		}

		// The result should contain "size:" from the template format
		if !strings.Contains(result, "size:") {
			t.Errorf("Expected result to contain 'size:' from template, got %q", result)
		}
	})

	// Test FormatListArchiveWithExtraction with template placeholders
	t.Run("FormatListArchiveWithExtraction_TemplatePlaceholders", func(t *testing.T) {
		// Create a temporary file to get real file stats
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test-archive.zip")

		// Create a test file with known content to verify size calculation
		content := []byte("test archive content for size calculation - this is a longer string to make the file size more visible and testable")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Get file stats to verify size_human is populated
		statInfo, err := formatter.GatherFileStatInfo(testFile)
		if err != nil {
			t.Fatalf("Failed to gather file stats: %v", err)
		}

		cfg := DefaultConfig()
		cfg.FormatListArchive = "#{path} (size: #{size_human})\n"

		formatterAdapter := NewFormatterAdapter(cfg)
		creationTime := "2025-07-13 19:31:42"

		result := formatterAdapter.FormatListArchiveWithExtraction(testFile, creationTime)

		t.Logf("Format string: #{path} (size: #{size_human})\\n")
		t.Logf("Result: %q", result)

		// CRITICAL: The result should NOT contain the literal #{size_human} placeholder
		if strings.Contains(result, "#{size_human}") {
			t.Errorf("BUG: Result contains literal #{size_human} placeholder, got %q", result)
		}

		// CRITICAL: The result should NOT contain Go template error format
		if strings.Contains(result, "%!{") {
			t.Errorf("BUG: Result contains Go template error format %%!{...}, got %q", result)
		}

		// CRITICAL: The result should NOT contain any #{...} patterns
		if strings.Contains(result, "#{") {
			t.Errorf("BUG: Result contains unprocessed template placeholder pattern #{...}, got %q", result)
		}

		// The result should contain the file path
		if !strings.Contains(result, testFile) {
			t.Errorf("Expected result to contain file path %q, got %q", testFile, result)
		}

		// The result should contain the actual size_human value or "unknown"
		if !strings.Contains(result, statInfo.SizeHuman) && !strings.Contains(result, "unknown") {
			t.Errorf("Expected result to contain size_human value %q or 'unknown', got %q", statInfo.SizeHuman, result)
		}

		// The result should contain "size:" from the template format
		if !strings.Contains(result, "size:") {
			t.Errorf("Expected result to contain 'size:' from template, got %q", result)
		}
	})
}
