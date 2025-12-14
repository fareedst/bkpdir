// [REQ:OUT_002] Enhanced Command Output with File Statistics testing
// [ARCH:FILE_STATISTICS] File statistics output formatting validation
// [IMPL:FILE_STATISTICS] FormatCreatedArchiveWithStats and FormatIncrementalCreatedWithStats validation
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestFormatCreatedArchiveWithStats tests format string processing with named replacements for full archives
func TestFormatCreatedArchiveWithStats_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-archive.zip")
	
	// Create a test archive file with known content
	content := []byte("test archive content for statistics")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test with default template format string
	result := formatter.FormatCreatedArchiveWithStats(testArchive)
	
	// Verify result contains file path
	if !strings.Contains(result, testArchive) {
		t.Errorf("Expected result to contain archive path %q, got %q", testArchive, result)
	}

	// Verify result contains size information (either size_human or size)
	statInfo, err := GatherFileStatInfo(testArchive)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}
	
	if !strings.Contains(result, statInfo.SizeHuman) && !strings.Contains(result, "unknown") {
		t.Errorf("Expected result to contain size_human %q or 'unknown', got %q", statInfo.SizeHuman, result)
	}

	// Verify no unprocessed placeholders remain
	if strings.Contains(result, "#{") {
		t.Errorf("Result contains unprocessed placeholders, got %q", result)
	}

	// Test with custom template format string
	cfg.TemplateCreatedArchiveDetailed = "Archive: #{path} (size: #{size_human}, type: #{type}, modified: #{mtime})\n"
	formatter = NewOutputFormatter(cfg)
	result = formatter.FormatCreatedArchiveWithStats(testArchive)

	// Verify all placeholders are replaced
	if strings.Contains(result, "#{path}") || strings.Contains(result, "#{size_human}") || 
		strings.Contains(result, "#{type}") || strings.Contains(result, "#{mtime}") {
		t.Errorf("Result contains unprocessed placeholders, got %q", result)
	}

	// Verify custom format content
	if !strings.Contains(result, "Archive:") {
		t.Errorf("Expected result to contain 'Archive:', got %q", result)
	}
	if !strings.Contains(result, "size:") {
		t.Errorf("Expected result to contain 'size:', got %q", result)
	}
	if !strings.Contains(result, statInfo.Type) {
		t.Errorf("Expected result to contain file type %q, got %q", statInfo.Type, result)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestFormatIncrementalCreatedWithStats tests format string processing with named replacements for incremental archives
func TestFormatIncrementalCreatedWithStats_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-incremental.zip")
	
	// Create a test archive file with known content
	content := []byte("test incremental archive content")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test with default template format string
	result := formatter.FormatIncrementalCreatedWithStats(testArchive)
	
	// Verify result contains file path
	if !strings.Contains(result, testArchive) {
		t.Errorf("Expected result to contain archive path %q, got %q", testArchive, result)
	}

	// Verify result contains size information
	statInfo, err := GatherFileStatInfo(testArchive)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}
	
	if !strings.Contains(result, statInfo.SizeHuman) && !strings.Contains(result, "unknown") {
		t.Errorf("Expected result to contain size_human %q or 'unknown', got %q", statInfo.SizeHuman, result)
	}

	// Verify no unprocessed placeholders remain
	if strings.Contains(result, "#{") {
		t.Errorf("Result contains unprocessed placeholders, got %q", result)
	}

	// Test with custom template format string
	cfg.TemplateIncrementalCreatedDetailed = "Incremental: #{path} (size: #{size_human}, modified: #{mtime})\n"
	formatter = NewOutputFormatter(cfg)
	result = formatter.FormatIncrementalCreatedWithStats(testArchive)

	// Verify all placeholders are replaced
	if strings.Contains(result, "#{path}") || strings.Contains(result, "#{size_human}") || 
		strings.Contains(result, "#{mtime}") {
		t.Errorf("Result contains unprocessed placeholders, got %q", result)
	}

	// Verify custom format content
	if !strings.Contains(result, "Incremental:") {
		t.Errorf("Expected result to contain 'Incremental:', got %q", result)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestFormatCreatedArchiveWithStatsAllPlaceholders tests all available placeholders
func TestFormatCreatedArchiveWithStatsAllPlaceholders_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-all-placeholders.zip")
	
	content := []byte("test content")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	// Use a template that includes all possible placeholders
	cfg.TemplateCreatedArchiveDetailed = "Path: #{path}, Name: #{name}, Size: #{size}, SizeHuman: #{size_human}, MTime: #{mtime}, MTimeUnix: #{mtime_unix}, Mode: #{mode}, Type: #{type}\n"
	formatter := NewOutputFormatter(cfg)

	result := formatter.FormatCreatedArchiveWithStats(testArchive)
	statInfo, err := GatherFileStatInfo(testArchive)
	if err != nil {
		t.Fatalf("Failed to gather file stats: %v", err)
	}

	// Verify all placeholders are replaced with actual values
	checks := []struct {
		placeholder string
		expected    string
		description string
	}{
		{"#{path}", statInfo.Path, "path"},
		{"#{name}", statInfo.Name, "name"},
		{"#{size}", "", "size (numeric)"},
		{"#{size_human}", statInfo.SizeHuman, "size_human"},
		{"#{mtime}", "", "mtime (formatted)"},
		{"#{mtime_unix}", "", "mtime_unix (numeric)"},
		{"#{mode}", "", "mode"},
		{"#{type}", statInfo.Type, "type"},
	}

	for _, check := range checks {
		if strings.Contains(result, check.placeholder) {
			t.Errorf("Placeholder %s (%s) was not replaced in result: %q", check.placeholder, check.description, result)
		}
	}

	// Verify actual values are present
	if !strings.Contains(result, statInfo.Path) {
		t.Errorf("Expected result to contain path %q, got %q", statInfo.Path, result)
	}
	if !strings.Contains(result, statInfo.Name) {
		t.Errorf("Expected result to contain name %q, got %q", statInfo.Name, result)
	}
	if !strings.Contains(result, statInfo.SizeHuman) {
		t.Errorf("Expected result to contain size_human %q, got %q", statInfo.SizeHuman, result)
	}
	if !strings.Contains(result, statInfo.Type) {
		t.Errorf("Expected result to contain type %q, got %q", statInfo.Type, result)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestFormatCreatedArchiveWithStatsErrorHandling tests error handling when stat gathering fails
func TestFormatCreatedArchiveWithStatsErrorHandling_REQ_OUT_002(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test with nonexistent file - should fallback to basic format
	nonexistentPath := "/path/that/does/not/exist/archive.zip"
	result := formatter.FormatCreatedArchiveWithStats(nonexistentPath)

	// Should fallback to basic FormatCreatedArchive (which doesn't require stats)
	// Result should still contain the path (even if file doesn't exist)
	if !strings.Contains(result, nonexistentPath) {
		t.Errorf("Expected result to contain path even for nonexistent file, got %q", result)
	}

	// Should not contain unprocessed placeholders
	if strings.Contains(result, "#{size_human}") {
		t.Errorf("Result should not contain unprocessed placeholders for nonexistent file, got %q", result)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestPrintCreatedArchiveWithStats tests the print method for full archives
func TestPrintCreatedArchiveWithStats_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-print.zip")
	
	content := []byte("test content")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Capture output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	formatter.PrintCreatedArchiveWithStats(testArchive)

	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify output contains archive path
	if !strings.Contains(output, testArchive) {
		t.Errorf("Expected output to contain archive path %q, got %q", testArchive, output)
	}

	// Verify no unprocessed placeholders
	if strings.Contains(output, "#{") {
		t.Errorf("Output contains unprocessed placeholders, got %q", output)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestPrintIncrementalCreatedWithStats tests the print method for incremental archives
func TestPrintIncrementalCreatedWithStats_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-incremental-print.zip")
	
	content := []byte("test incremental content")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Capture output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	formatter.PrintIncrementalCreatedWithStats(testArchive)

	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify output contains archive path
	if !strings.Contains(output, testArchive) {
		t.Errorf("Expected output to contain archive path %q, got %q", testArchive, output)
	}

	// Verify no unprocessed placeholders
	if strings.Contains(output, "#{") {
		t.Errorf("Output contains unprocessed placeholders, got %q", output)
	}
}

// [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]
// TestBackwardCompatibilityFormatStrings tests that existing format strings still work
func TestBackwardCompatibilityFormatStrings_REQ_OUT_002(t *testing.T) {
	tmpDir := t.TempDir()
	testArchive := filepath.Join(tmpDir, "test-backward-compat.zip")
	
	content := []byte("test content")
	if err := os.WriteFile(testArchive, content, 0644); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test that basic format methods still work (backward compatibility)
	basicResult := formatter.FormatCreatedArchive(testArchive)
	if basicResult == "" {
		t.Error("FormatCreatedArchive returned empty string")
	}
	if !strings.Contains(basicResult, testArchive) {
		t.Errorf("Expected FormatCreatedArchive to contain path, got %q", basicResult)
	}

	// Test that stats methods work alongside basic methods
	statsResult := formatter.FormatCreatedArchiveWithStats(testArchive)
	if statsResult == "" {
		t.Error("FormatCreatedArchiveWithStats returned empty string")
	}
	if !strings.Contains(statsResult, testArchive) {
		t.Errorf("Expected FormatCreatedArchiveWithStats to contain path, got %q", statsResult)
	}

	// Both should produce valid output
	if basicResult == statsResult {
		// This is acceptable - stats method may fallback to basic if template is same
		t.Logf("Both methods produced same output (acceptable): %q", basicResult)
	}
}

