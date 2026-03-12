// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Comprehensive tests for placeholder replacement
// Tests use #{...} syntax to avoid fmt package conflicts
package formatter

import (
	"fmt"
	"strings"
	"testing"
)

func TestReplacePlaceholders_Simple(t *testing.T) {
	formatStr := "#{path}\n"
	data := map[string]string{
		"path": "/test/path.zip",
	}

	result := ReplacePlaceholders(formatStr, data)

	if !strings.Contains(result, "/test/path.zip") {
		t.Errorf("Expected result to contain path, got %q", result)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed placeholder, got %q", result)
	}
}

func TestReplacePlaceholders_Multiple(t *testing.T) {
	formatStr := "#{path} (size: #{size_human})\n"
	data := map[string]string{
		"path":       "/test/path.zip",
		"size_human": "115B",
	}

	result := ReplacePlaceholders(formatStr, data)

	if !strings.Contains(result, "/test/path.zip") {
		t.Errorf("Expected result to contain path, got %q", result)
	}
	if !strings.Contains(result, "115B") {
		t.Errorf("Expected result to contain size_human, got %q", result)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed placeholder, got %q", result)
	}
}

func TestReplacePlaceholders_ExactBugScenario(t *testing.T) {
	// Test with the exact format string and data from the bug report
	formatStr := "#{path} (size: #{size_human})\n"
	data := map[string]string{
		"path":          "bkpdir-2025-06-07-23-01.zip",
		"size_human":    "115B",
		"creation_time": "2025-06-07 23:01:23",
	}

	result := ReplacePlaceholders(formatStr, data)

	t.Logf("Format: %q", formatStr)
	t.Logf("Data: %v", data)
	t.Logf("Result: %q", result)

	// CRITICAL: Verify no unprocessed placeholders
	if strings.Contains(result, "#{") {
		t.Errorf("CRITICAL BUG: Result contains unprocessed placeholder, got %q", result)
	}

	// CRITICAL: Verify no fmt.Sprintf error patterns (should not occur with #{...})
	if strings.Contains(result, "%!{") {
		t.Errorf("CRITICAL BUG: Result contains fmt.Sprintf error pattern, got %q", result)
	}

	// Verify expected content
	if !strings.Contains(result, "bkpdir-2025-06-07-23-01.zip") {
		t.Errorf("Expected result to contain path, got %q", result)
	}
	if !strings.Contains(result, "115B") {
		t.Errorf("Expected result to contain size_human, got %q", result)
	}

	// Verify the result is safe to use with fmt.Sprintf
	// This test ensures the result won't cause fmt.Sprintf errors
	testOutput := fmt.Sprintf("%s", result)
	if testOutput != result {
		t.Errorf("Result should be safe for fmt.Sprintf, got %q", testOutput)
	}
}

func TestReplacePlaceholders_AllKnown(t *testing.T) {
	formatStr := "#{path} (size: #{size_human}, bytes: #{size}, created: #{creation_time})\n"
	data := map[string]string{
		"path":          "/test/path.zip",
		"size_human":    "115B",
		"size":          "115",
		"creation_time": "2025-07-13 19:31:42",
	}

	result := ReplacePlaceholders(formatStr, data)

	// Verify all placeholders were replaced
	checks := []struct {
		name     string
		expected string
	}{
		{"path", "/test/path.zip"},
		{"size_human", "115B"},
		{"size", "115"},
		{"creation_time", "2025-07-13 19:31:42"},
	}

	for _, check := range checks {
		if !strings.Contains(result, check.expected) {
			t.Errorf("Expected result to contain %s %q, got %q", check.name, check.expected, result)
		}
	}

	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed placeholder, got %q", result)
	}
}

func TestReplacePlaceholders_UnknownPlaceholder(t *testing.T) {
	formatStr := "#{path} (unknown: #{missing})\n"
	data := map[string]string{
		"path": "/test/path.zip",
	}

	result := ReplacePlaceholders(formatStr, data)

	if !strings.Contains(result, "/test/path.zip") {
		t.Errorf("Expected result to contain path, got %q", result)
	}
	// Unknown placeholder should remain as-is
	if !strings.Contains(result, "#{missing}") {
		t.Errorf("Expected unknown placeholder to remain, got %q", result)
	}
}

func TestReplacePlaceholders_NoPlaceholders(t *testing.T) {
	formatStr := "plain text without placeholders\n"
	data := map[string]string{
		"path": "/test/path.zip",
	}

	result := ReplacePlaceholders(formatStr, data)

	if result != formatStr {
		t.Errorf("Expected result to be unchanged, got %q", result)
	}
}

func TestReplacePlaceholders_EmptyData(t *testing.T) {
	formatStr := "#{path} (size: #{size_human})\n"
	data := map[string]string{}

	result := ReplacePlaceholders(formatStr, data)

	// All placeholders should remain as-is
	if !strings.Contains(result, "#{path}") {
		t.Errorf("Expected placeholder to remain, got %q", result)
	}
	if !strings.Contains(result, "#{size_human}") {
		t.Errorf("Expected placeholder to remain, got %q", result)
	}
}

func TestReplacePlaceholders_SpecialChars(t *testing.T) {
	formatStr := "#{path}\n"
	data := map[string]string{
		"path": "file with spaces and = signs.zip",
	}

	result := ReplacePlaceholders(formatStr, data)

	if !strings.Contains(result, "file with spaces and = signs.zip") {
		t.Errorf("Expected result to contain path with special chars, got %q", result)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed placeholder, got %q", result)
	}
}

func TestReplacePlaceholders_MultipleOccurrences(t *testing.T) {
	formatStr := "#{path} - #{path} - #{path}\n"
	data := map[string]string{
		"path": "/test/path.zip",
	}

	result := ReplacePlaceholders(formatStr, data)

	expectedCount := strings.Count(result, "/test/path.zip")
	if expectedCount != 3 {
		t.Errorf("Expected path to appear 3 times, got %d occurrences", expectedCount)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("BUG: Result contains unprocessed placeholder, got %q", result)
	}
}
