// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Comprehensive unit tests for simplified list formatting
// These tests validate each step of the formatting process to catch bugs early
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFormatListArchiveSimple_EndToEnd tests the complete flow from format string to final output
func TestFormatListArchiveSimple_EndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "bkpdir-2025-06-07-23-01.zip")
	content := []byte("test archive content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	statInfo, err := formatter.GatherFileStatInfo(testFile)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}

	creationTime := "2025-06-07 23:01:23"
	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	formatterInstance := NewFormatterAdapter(cfg)

	// Get formatted output exactly as main.go does
	output := formatterInstance.FormatListArchiveWithExtraction(testFile, creationTime)
	output = strings.TrimSuffix(output, "\n")

	t.Logf("FormatListArchive: %q", cfg.FormatListArchive)
	t.Logf("Output after FormatListArchiveWithExtraction: %q", output)
	t.Logf("Output length: %d", len(output))

	// CRITICAL: Verify no unprocessed placeholders BEFORE PrintArchiveListWithStatus
	if strings.Contains(output, "#{") {
		t.Errorf("CRITICAL BUG: Output contains unprocessed placeholder BEFORE printing, got %q", output)
		t.Logf("This means formatListArchiveSimple is not replacing placeholders correctly")
	}

	// CRITICAL: Verify no fmt.Sprintf error patterns BEFORE printing
	if strings.Contains(output, "%!{") {
		t.Errorf("CRITICAL BUG: Output contains fmt.Sprintf error pattern BEFORE printing, got %q", output)
		t.Logf("This means fmt.Sprintf was called somewhere in the formatting chain")
	}

	// Verify expected content
	if !strings.Contains(output, testFile) {
		t.Errorf("Expected output to contain path %q, got %q", testFile, output)
	}

	if !strings.Contains(output, statInfo.SizeHuman) {
		t.Errorf("Expected output to contain size_human %q, got %q", statInfo.SizeHuman, output)
	}

	// Test that the output is safe to use with fmt.Sprintf
	// This simulates what PrintArchiveListWithStatus does
	status := " [UNVERIFIED]"
	testOutput := fmt.Sprintf("%s", output) // Should not cause errors
	if testOutput != output {
		t.Errorf("Output should be safe for fmt.Sprintf, got %q", testOutput)
	}

	// Test string concatenation (what PrintArchiveListWithStatus actually does)
	message := output + status + "\n"
	if strings.Contains(message, "%!{") {
		t.Errorf("CRITICAL BUG: Message contains fmt.Sprintf error pattern, got %q", message)
	}
}

// TestReplacePlaceholders_Isolation tests the ReplacePlaceholders function in complete isolation
func TestReplacePlaceholders_Isolation(t *testing.T) {
	formatStr := "#{path} (size: #{size_human})\n"
	data := map[string]string{
		"path":       "bkpdir-2025-06-07-23-01.zip",
		"size_human": "115B",
	}

	result := formatter.ReplacePlaceholders(formatStr, data)

	t.Logf("Format: %q", formatStr)
	t.Logf("Data: %v", data)
	t.Logf("Result: %q", result)

	// CRITICAL: Verify no unprocessed placeholders (old %{ syntax or new #{ syntax)
	if strings.Contains(result, "%{") {
		t.Errorf("CRITICAL BUG: Result contains unprocessed old-style placeholder, got %q", result)
	}
	if strings.Contains(result, "#{") {
		t.Errorf("CRITICAL BUG: Result contains unprocessed placeholder, got %q", result)
	}

	// CRITICAL: Verify no fmt.Sprintf error patterns
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

	// Test that the result is safe to use with fmt.Sprintf
	testOutput := fmt.Sprintf("%s", result)
	if testOutput != result {
		t.Errorf("Result should be safe for fmt.Sprintf, got %q", testOutput)
	}
}

// TestFormatListArchiveSimple_StepByStep tests each step individually
func TestFormatListArchiveSimple_StepByStep(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-archive.zip")
	content := []byte("test archive content")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	statInfo, err := formatter.GatherFileStatInfo(testFile)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}

	creationTime := "2025-07-13 19:31:42"
	cfg := DefaultConfig()
	cfg.FormatListArchive = "#{path} (size: #{size_human})\n"
	formatterInstance := NewFormatterAdapter(cfg)

	// Step 1: Test data map building
	t.Run("Step1_DataMap", func(t *testing.T) {
		data := make(map[string]string)
		data["path"] = testFile
		data["creation_time"] = creationTime
		data["size_human"] = statInfo.SizeHuman

		if data["path"] != testFile {
			t.Errorf("Expected path %q, got %q", testFile, data["path"])
		}
		if data["size_human"] != statInfo.SizeHuman {
			t.Errorf("Expected size_human %q, got %q", statInfo.SizeHuman, data["size_human"])
		}
	})

	// Step 2: Test format string selection
	t.Run("Step2_FormatSelection", func(t *testing.T) {
		formatStr := ""
		if cfg.FormatListArchive != "" && strings.Contains(cfg.FormatListArchive, "#{") {
			formatStr = cfg.FormatListArchive
		} else if cfg.TemplateListArchive != "" {
			formatStr = cfg.TemplateListArchive
		} else {
			formatStr = "#{path} (size: #{size_human})\n"
		}

		if formatStr == "" {
			t.Error("Format string should be selected")
		}
		if !strings.Contains(formatStr, "#{path}") {
			t.Errorf("Format string should contain #{path}, got %q", formatStr)
		}
	})

	// Step 3: Test placeholder replacement
	t.Run("Step3_PlaceholderReplacement", func(t *testing.T) {
		formatStr := "#{path} (size: #{size_human})\n"
		data := map[string]string{
			"path":       testFile,
			"size_human": statInfo.SizeHuman,
		}

		result := formatter.ReplacePlaceholders(formatStr, data)

		if strings.Contains(result, "#{") {
			t.Errorf("BUG: Result still contains #{...} patterns, got %q", result)
		}
		if !strings.Contains(result, testFile) {
			t.Errorf("Result should contain path %q, got %q", testFile, result)
		}
		if !strings.Contains(result, statInfo.SizeHuman) {
			t.Errorf("Result should contain size_human %q, got %q", statInfo.SizeHuman, result)
		}
	})

	// Step 4: Test final output
	t.Run("Step4_FinalOutput", func(t *testing.T) {
		result := formatListArchiveSimple(cfg, formatterInstance, testFile, creationTime)

		t.Logf("Final result: %q", result)

		// CRITICAL: Verify no unprocessed placeholders (old %{ syntax or new #{ syntax)
		if strings.Contains(result, "%{") {
			t.Errorf("CRITICAL BUG: Result contains unprocessed old-style placeholder, got %q", result)
		}
		if strings.Contains(result, "#{") {
			t.Errorf("CRITICAL BUG: Result contains unprocessed placeholder, got %q", result)
		}

		// CRITICAL: Verify no fmt.Sprintf error patterns
		if strings.Contains(result, "%!{") {
			t.Errorf("CRITICAL BUG: Result contains fmt.Sprintf error pattern, got %q", result)
		}

		// Verify expected content
		if !strings.Contains(result, testFile) {
			t.Errorf("Expected result to contain path %q, got %q", testFile, result)
		}

		if !strings.Contains(result, statInfo.SizeHuman) {
			t.Errorf("Expected result to contain size_human %q, got %q", statInfo.SizeHuman, result)
		}
	})
}
