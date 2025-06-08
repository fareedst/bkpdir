package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// ⭐ OUT-002: File stat information gathering - 🧪
// TestGatherFileStatInfo tests the file statistics gathering functionality
func TestGatherFileStatInfo(t *testing.T) {
	// Create a temporary file for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "testfile.txt")
	testContent := "Hello, World! This is test content for file stats."

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test gathering file stats
	stats, err := GatherFileStatInfo(testFile)
	if err != nil {
		t.Fatalf("GatherFileStatInfo failed: %v", err)
	}

	// Verify the basic fields
	if stats.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, stats.Path)
	}

	if stats.Name != "testfile.txt" {
		t.Errorf("Expected name testfile.txt, got %s", stats.Name)
	}

	if stats.Size != int64(len(testContent)) {
		t.Errorf("Expected size %d, got %d", len(testContent), stats.Size)
	}

	if stats.Type != "regular" {
		t.Errorf("Expected type regular, got %s", stats.Type)
	}

	// Verify the modification time is recent (within last minute)
	now := time.Now()
	if stats.MTime.After(now) || stats.MTime.Before(now.Add(-time.Minute)) {
		t.Errorf("Modification time %v is not recent", stats.MTime)
	}

	if stats.MTimeUnix != stats.MTime.Unix() {
		t.Errorf("Unix timestamp mismatch: %d vs %d", stats.MTimeUnix, stats.MTime.Unix())
	}

	// Verify human-readable size
	expectedSize := formatHumanSize(int64(len(testContent)))
	if stats.SizeHuman != expectedSize {
		t.Errorf("Expected human size %s, got %s", expectedSize, stats.SizeHuman)
	}
}

// ⭐ OUT-002: File stat information gathering - 🧪
// TestGatherFileStatInfoDirectory tests file stats for directories
func TestGatherFileStatInfoDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	stats, err := GatherFileStatInfo(tmpDir)
	if err != nil {
		t.Fatalf("GatherFileStatInfo failed for directory: %v", err)
	}

	if stats.Type != "directory" {
		t.Errorf("Expected type directory, got %s", stats.Type)
	}

	if stats.Name != filepath.Base(tmpDir) {
		t.Errorf("Expected name %s, got %s", filepath.Base(tmpDir), stats.Name)
	}
}

// ⭐ OUT-002: File stat information gathering - 🧪
// TestGatherFileStatInfoNonexistent tests error handling for nonexistent files
func TestGatherFileStatInfoNonexistent(t *testing.T) {
	nonexistentFile := "/path/that/does/not/exist"

	_, err := GatherFileStatInfo(nonexistentFile)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// ⭐ OUT-002: File stat information gathering - 🧪
// TestFormatHumanSize tests the human-readable size formatting
func TestFormatHumanSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0B"},
		{100, "100B"},
		{1023, "1023B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
		{1572864, "1.5MB"},
		{1073741824, "1.0GB"},
		{1610612736, "1.5GB"},
		{1099511627776, "1.0TB"},
		{1649267441664, "1.5TB"},
	}

	for _, test := range tests {
		result := formatHumanSize(test.size)
		if result != test.expected {
			t.Errorf("formatHumanSize(%d) = %s, expected %s", test.size, result, test.expected)
		}
	}
}

// ⭐ OUT-002: File stat information gathering - 🧪
// TestGetFileType tests file type detection
func TestGetFileType(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "testfile.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test regular file
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	fileType := getFileType(info)
	if fileType != "regular" {
		t.Errorf("Expected type regular, got %s", fileType)
	}

	// Test directory
	info, err = os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("Failed to stat directory: %v", err)
	}

	fileType = getFileType(info)
	if fileType != "directory" {
		t.Errorf("Expected type directory, got %s", fileType)
	}
}
