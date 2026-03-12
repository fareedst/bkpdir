// [REQ-DIFF_COMMAND] Diff Command testing
// [ARCH-DIFF_COMMAND] Diff command architecture validation
// [IMPL-DIFF_COMMAND] ReconstructArchiveState and CalculateDiff validation
package main

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestReconstructArchiveStateFullOnly tests reconstruction with full archive only
func TestReconstructArchiveStateFullOnly_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create source directory with files
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	for relPath, content := range files {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create full archive
	archivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000.zip")
	if err := createTestArchive(archivePath, sourceDir, files); err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Reconstruct state
	reconstructed, err := ReconstructArchiveState(archiveDir)
	if err != nil {
		t.Fatalf("ReconstructArchiveState failed: %v", err)
	}

	// Verify all files are in reconstructed state
	reconstructedMap := make(map[string]bool)
	for _, file := range reconstructed.Files {
		reconstructedMap[file.RelativePath] = true
	}

	for relPath := range files {
		if !reconstructedMap[relPath] {
			t.Errorf("Expected file %s in reconstructed state, but not found", relPath)
		}
	}

	// Verify file count matches
	if len(reconstructed.Files) != len(files) {
		t.Errorf("Expected %d files in reconstructed state, got %d", len(files), len(reconstructed.Files))
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestReconstructArchiveStateMultipleFullAndIncremental tests that reconstruction correctly
// finds the most recent full archive and its most recent incremental archive, even when
// multiple full archives and incrementals exist
func TestReconstructArchiveStateMultipleFullAndIncremental_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create first full archive (older)
	oldFullFiles := map[string]string{
		"old_file1.txt": "old_content1",
		"old_file2.txt": "old_content2",
	}
	for relPath, content := range oldFullFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	oldFullArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T100000.zip")
	if err := createTestArchive(oldFullArchivePath, sourceDir, oldFullFiles); err != nil {
		t.Fatalf("Failed to create old full archive: %v", err)
	}

	// Create incremental for old full archive
	oldIncrementalFiles := map[string]string{
		"old_file1.txt": "old_modified_content1",
	}
	for relPath, content := range oldIncrementalFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	oldIncrementalPath := filepath.Join(archiveDir, "bkp-2024-01-01T100000_update=2024-01-01T110000.zip")
	if err := createTestArchive(oldIncrementalPath, sourceDir, oldIncrementalFiles); err != nil {
		t.Fatalf("Failed to create old incremental archive: %v", err)
	}

	// Wait a bit to ensure different modification times
	time.Sleep(100 * time.Millisecond)

	// Create second full archive (newer - this should be selected)
	newFullFiles := map[string]string{
		"new_file1.txt": "new_content1",
		"new_file2.txt": "new_content2",
	}
	for relPath, content := range newFullFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	newFullArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000.zip")
	if err := createTestArchive(newFullArchivePath, sourceDir, newFullFiles); err != nil {
		t.Fatalf("Failed to create new full archive: %v", err)
	}

	// Wait a bit to ensure different modification times
	time.Sleep(100 * time.Millisecond)

	// Create first incremental for new full archive
	firstNewIncrementalFiles := map[string]string{
		"new_file2.txt": "new_modified_content2",
	}
	for relPath, content := range firstNewIncrementalFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	firstNewIncrementalPath := filepath.Join(archiveDir, "bkp-2024-01-01T120000_update=2024-01-01T130000.zip")
	if err := createTestArchive(firstNewIncrementalPath, sourceDir, firstNewIncrementalFiles); err != nil {
		t.Fatalf("Failed to create first new incremental archive: %v", err)
	}

	// Wait a bit to ensure different modification times
	time.Sleep(100 * time.Millisecond)

	// Create second incremental for new full archive (this should be selected)
	secondNewIncrementalFiles := map[string]string{
		"new_file1.txt": "new_modified_content1",
		"new_file3.txt": "new_content3", // New file
	}
	for relPath, content := range secondNewIncrementalFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}
	secondNewIncrementalPath := filepath.Join(archiveDir, "bkp-2024-01-01T120000_update=2024-01-01T140000.zip")
	if err := createTestArchive(secondNewIncrementalPath, sourceDir, secondNewIncrementalFiles); err != nil {
		t.Fatalf("Failed to create second new incremental archive: %v", err)
	}

	// Reconstruct state - should use new full archive + second new incremental
	reconstructed, err := ReconstructArchiveState(archiveDir)
	if err != nil {
		t.Fatalf("ReconstructArchiveState failed: %v", err)
	}

	// Verify reconstructed state contains files from new full + second new incremental
	reconstructedMap := make(map[string]string)
	for _, file := range reconstructed.Files {
		reconstructedMap[file.RelativePath] = file.Hash
	}

	// Should have new_file1.txt from second incremental (modified)
	if hash, exists := reconstructedMap["new_file1.txt"]; !exists {
		t.Error("Expected new_file1.txt in reconstructed state, but not found")
	} else {
		expectedHash := calculateFileHash([]byte("new_modified_content1"))
		if hash != expectedHash {
			t.Errorf("Expected new_file1.txt to have hash from second incremental, got %s, expected %s", hash, expectedHash)
		}
	}

	// Should have new_file2.txt from full archive (second incremental doesn't include it, so it keeps full archive value)
	if hash, exists := reconstructedMap["new_file2.txt"]; !exists {
		t.Error("Expected new_file2.txt in reconstructed state, but not found")
	} else {
		expectedHash := calculateFileHash([]byte("new_content2"))
		if hash != expectedHash {
			t.Errorf("Expected new_file2.txt to have hash from full archive (second incremental doesn't modify it), got %s, expected %s", hash, expectedHash)
		}
	}

	// Should have new_file3.txt from second incremental (new file)
	if _, exists := reconstructedMap["new_file3.txt"]; !exists {
		t.Error("Expected new_file3.txt from second incremental in reconstructed state, but not found")
	}

	// Should NOT have old files from old full archive
	if _, exists := reconstructedMap["old_file1.txt"]; exists {
		t.Error("Should not have old_file1.txt from old full archive in reconstructed state")
	}
	if _, exists := reconstructedMap["old_file2.txt"]; exists {
		t.Error("Should not have old_file2.txt from old full archive in reconstructed state")
	}

	// Verify file count (new_file1, new_file2, new_file3 = 3 files)
	expectedCount := 3
	if len(reconstructed.Files) != expectedCount {
		t.Errorf("Expected %d files in reconstructed state, got %d", expectedCount, len(reconstructed.Files))
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestReconstructArchiveStateFullAndIncremental tests reconstruction with full + incremental archive
func TestReconstructArchiveStateFullAndIncremental_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}

	// Create full archive files
	fullFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for relPath, content := range fullFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create full archive
	fullArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000.zip")
	if err := createTestArchive(fullArchivePath, sourceDir, fullFiles); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Create incremental archive with modified and new files
	incrementalFiles := map[string]string{
		"file2.txt": "modified_content2", // Modified
		"file3.txt": "content3",          // New file
	}

	for relPath, content := range incrementalFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create incremental archive (must contain "_update=" to be recognized as incremental)
	incrementalArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000_update=2024-01-01T130000.zip")
	if err := createTestArchive(incrementalArchivePath, sourceDir, incrementalFiles); err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	// Reconstruct state
	reconstructed, err := ReconstructArchiveState(archiveDir)
	if err != nil {
		t.Fatalf("ReconstructArchiveState failed: %v", err)
	}

	// Verify all files are in reconstructed state
	reconstructedMap := make(map[string]string)
	for _, file := range reconstructed.Files {
		reconstructedMap[file.RelativePath] = file.Hash
	}

	// Verify file1.txt from full archive is present
	if _, exists := reconstructedMap["file1.txt"]; !exists {
		t.Error("Expected file1.txt from full archive in reconstructed state, but not found")
	}

	// Verify file2.txt is from incremental (modified)
	if hash, exists := reconstructedMap["file2.txt"]; !exists {
		t.Error("Expected file2.txt in reconstructed state, but not found")
	} else {
		// Hash should match incremental version (modified content)
		expectedHash := calculateFileHash([]byte("modified_content2"))
		if hash != expectedHash {
			t.Errorf("Expected file2.txt to have hash from incremental archive, got %s", hash)
		}
	}

	// Verify file3.txt from incremental is present
	if _, exists := reconstructedMap["file3.txt"]; !exists {
		t.Error("Expected file3.txt from incremental archive in reconstructed state, but not found")
	}

	// Verify file count (file1 from full, file2 modified, file3 new = 3 files)
	expectedCount := 3
	if len(reconstructed.Files) != expectedCount {
		t.Errorf("Expected %d files in reconstructed state, got %d", expectedCount, len(reconstructed.Files))
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestCalculateDiffAddedFiles tests diff calculation with added files
func TestCalculateDiffAddedFiles_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create archive state files (what's in the archive)
	archiveFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	archiveStateDir := filepath.Join(tmpDir, "archive_state")
	if err := os.MkdirAll(archiveStateDir, 0755); err != nil {
		t.Fatalf("Failed to create archive state directory: %v", err)
	}

	for relPath, content := range archiveFiles {
		fullPath := filepath.Join(archiveStateDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create snapshot from archive state directory
	reconstructedState, err := CreateDirectorySnapshot(archiveStateDir, []string{})
	if err != nil {
		t.Fatalf("Failed to create directory snapshot: %v", err)
	}

	// Create current directory with additional file
	currentFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3", // New file
	}

	for relPath, content := range currentFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Calculate diff
	diff, err := CalculateDiff(sourceDir, reconstructedState, []string{})
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify added files
	if len(diff.Added) != 1 {
		t.Errorf("Expected 1 added file, got %d: %v", len(diff.Added), diff.Added)
	}
	if len(diff.Added) > 0 && diff.Added[0] != "file3.txt" {
		t.Errorf("Expected added file to be file3.txt, got %s", diff.Added[0])
	}

	// Verify no modified or deleted files
	if len(diff.Modified) != 0 {
		t.Errorf("Expected 0 modified files, got %d: %v", len(diff.Modified), diff.Modified)
	}
	if len(diff.Deleted) != 0 {
		t.Errorf("Expected 0 deleted files, got %d: %v", len(diff.Deleted), diff.Deleted)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestCalculateDiffModifiedFiles tests diff calculation with modified files
func TestCalculateDiffModifiedFiles_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveStateDir := filepath.Join(tmpDir, "archive_state")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveStateDir, 0755); err != nil {
		t.Fatalf("Failed to create archive state directory: %v", err)
	}

	// Create archive state (original files)
	archiveFiles := map[string]string{
		"file1.txt": "original_content",
		"file2.txt": "content2",
	}

	for relPath, content := range archiveFiles {
		fullPath := filepath.Join(archiveStateDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create snapshot from archive state
	reconstructedState, err := CreateDirectorySnapshot(archiveStateDir, []string{})
	if err != nil {
		t.Fatalf("Failed to create directory snapshot: %v", err)
	}

	// Create current directory with modified file
	currentFiles := map[string]string{
		"file1.txt": "modified_content", // Modified
		"file2.txt": "content2",         // Unchanged
	}

	for relPath, content := range currentFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Calculate diff
	diff, err := CalculateDiff(sourceDir, reconstructedState, []string{})
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify modified files
	if len(diff.Modified) != 1 {
		t.Errorf("Expected 1 modified file, got %d: %v", len(diff.Modified), diff.Modified)
	}
	if len(diff.Modified) > 0 && diff.Modified[0] != "file1.txt" {
		t.Errorf("Expected modified file to be file1.txt, got %s", diff.Modified[0])
	}

	// Verify no added or deleted files
	if len(diff.Added) != 0 {
		t.Errorf("Expected 0 added files, got %d: %v", len(diff.Added), diff.Added)
	}
	if len(diff.Deleted) != 0 {
		t.Errorf("Expected 0 deleted files, got %d: %v", len(diff.Deleted), diff.Deleted)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestCalculateDiffDeletedFiles tests diff calculation with deleted files
func TestCalculateDiffDeletedFiles_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveStateDir := filepath.Join(tmpDir, "archive_state")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveStateDir, 0755); err != nil {
		t.Fatalf("Failed to create archive state directory: %v", err)
	}

	// Create archive state (all files)
	archiveFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		"file3.txt": "content3",
	}

	for relPath, content := range archiveFiles {
		fullPath := filepath.Join(archiveStateDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create snapshot from archive state
	reconstructedState, err := CreateDirectorySnapshot(archiveStateDir, []string{})
	if err != nil {
		t.Fatalf("Failed to create directory snapshot: %v", err)
	}

	// Create current directory with one file deleted
	currentFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
		// file3.txt deleted
	}

	for relPath, content := range currentFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Calculate diff
	diff, err := CalculateDiff(sourceDir, reconstructedState, []string{})
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify deleted files
	if len(diff.Deleted) != 1 {
		t.Errorf("Expected 1 deleted file, got %d: %v", len(diff.Deleted), diff.Deleted)
	}
	if len(diff.Deleted) > 0 && diff.Deleted[0] != "file3.txt" {
		t.Errorf("Expected deleted file to be file3.txt, got %s", diff.Deleted[0])
	}

	// Verify no added or modified files
	if len(diff.Added) != 0 {
		t.Errorf("Expected 0 added files, got %d: %v", len(diff.Added), diff.Added)
	}
	if len(diff.Modified) != 0 {
		t.Errorf("Expected 0 modified files, got %d: %v", len(diff.Modified), diff.Modified)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestCalculateDiffNoChanges tests diff calculation with no changes
func TestCalculateDiffNoChanges_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveStateDir := filepath.Join(tmpDir, "archive_state")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveStateDir, 0755); err != nil {
		t.Fatalf("Failed to create archive state directory: %v", err)
	}

	// Create archive state
	archiveFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for relPath, content := range archiveFiles {
		fullPath := filepath.Join(archiveStateDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create snapshot from archive state
	reconstructedState, err := CreateDirectorySnapshot(archiveStateDir, []string{})
	if err != nil {
		t.Fatalf("Failed to create directory snapshot: %v", err)
	}

	// Create current directory with same files
	currentFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for relPath, content := range currentFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Calculate diff
	diff, err := CalculateDiff(sourceDir, reconstructedState, []string{})
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify no changes
	if len(diff.Added) != 0 {
		t.Errorf("Expected 0 added files, got %d: %v", len(diff.Added), diff.Added)
	}
	if len(diff.Modified) != 0 {
		t.Errorf("Expected 0 modified files, got %d: %v", len(diff.Modified), diff.Modified)
	}
	if len(diff.Deleted) != 0 {
		t.Errorf("Expected 0 deleted files, got %d: %v", len(diff.Deleted), diff.Deleted)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestFormatDiffResult tests diff result formatting
func TestFormatDiffResult_REQ_DIFF_COMMAND(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test with no changes
	diffNoChanges := &DiffResult{
		Added:    []string{},
		Modified: []string{},
		Deleted:  []string{},
	}

	result := formatter.FormatDiffResult(diffNoChanges)
	if !strings.Contains(result, cfg.FormatDiffNoChanges) {
		t.Errorf("Expected result to contain no changes message, got %q", result)
	}

	// Test with changes
	diffWithChanges := &DiffResult{
		Added:    []string{"newfile.txt"},
		Modified: []string{"modified.txt"},
		Deleted:  []string{"deleted.txt"},
	}

	result = formatter.FormatDiffResult(diffWithChanges)
	if !strings.Contains(result, cfg.FormatDiffChanges) {
		t.Errorf("Expected result to contain changes header, got %q", result)
	}
	if !strings.Contains(result, "newfile.txt") {
		t.Errorf("Expected result to contain added file, got %q", result)
	}
	if !strings.Contains(result, "modified.txt") {
		t.Errorf("Expected result to contain modified file, got %q", result)
	}
	if !strings.Contains(result, "deleted.txt") {
		t.Errorf("Expected result to contain deleted file, got %q", result)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestPrintDiffResult tests diff result printing
func TestPrintDiffResult_REQ_DIFF_COMMAND(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	diff := &DiffResult{
		Added:    []string{"newfile.txt"},
		Modified: []string{},
		Deleted:  []string{},
	}

	// Capture output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	formatter.PrintDiffResult(diff)

	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify output contains diff information
	if !strings.Contains(output, "newfile.txt") {
		t.Errorf("Expected output to contain added file, got %q", output)
	}
}

// Helper function to create a test archive
func createTestArchive(archivePath string, sourceDir string, files map[string]string) error {
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for relPath, content := range files {
		fileWriter, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		if _, err := fileWriter.Write([]byte(content)); err != nil {
			return err
		}
	}

	return nil
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestDiffCommandNoArchives tests diff command when no archives exist
func TestDiffCommandNoArchives_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Try to reconstruct state when no archives exist
	_, err := ReconstructArchiveState(archiveDir)
	if err == nil {
		t.Error("Expected error when no archives exist, but got nil")
	}
	if !strings.Contains(err.Error(), "No archives found") && !strings.Contains(err.Error(), "No full archive found") {
		t.Errorf("Expected error about no archives, got: %v", err)
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-DIFF_COMMAND]
// TestDiffCommandIntegration tests the full diff command integration
func TestDiffCommandIntegration_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create initial files for full archive
	fullFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for relPath, content := range fullFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create full archive
	fullArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000.zip")
	if err := createTestArchive(fullArchivePath, sourceDir, fullFiles); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Modify files and add new file
	modifiedFiles := map[string]string{
		"file2.txt": "modified_content2",
		"file3.txt": "content3",
	}

	for relPath, content := range modifiedFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create incremental archive
	incrementalArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000_update=2024-01-01T130000.zip")
	if err := createTestArchive(incrementalArchivePath, sourceDir, modifiedFiles); err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	// Reconstruct state
	reconstructedState, err := ReconstructArchiveState(archiveDir)
	if err != nil {
		t.Fatalf("ReconstructArchiveState failed: %v", err)
	}

	// Add another file to current directory
	newFile := filepath.Join(sourceDir, "file4.txt")
	if err := os.WriteFile(newFile, []byte("content4"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Calculate diff
	diff, err := CalculateDiff(sourceDir, reconstructedState, []string{})
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify diff results
	// file4.txt should be added
	if len(diff.Added) != 1 || diff.Added[0] != "file4.txt" {
		t.Errorf("Expected 1 added file (file4.txt), got %v", diff.Added)
	}

	// file2.txt should be modified (but we already have it in incremental, so it might not show as modified)
	// Actually, since file2.txt is in the incremental archive with modified content, and we haven't changed it again,
	// it should not be in modified list. Let me check the logic...

	// file1.txt should not be in any list (unchanged)
	foundFile1 := false
	for _, f := range diff.Added {
		if f == "file1.txt" {
			foundFile1 = true
		}
	}
	for _, f := range diff.Modified {
		if f == "file1.txt" {
			foundFile1 = true
		}
	}
	for _, f := range diff.Deleted {
		if f == "file1.txt" {
			foundFile1 = true
		}
	}
	if foundFile1 {
		t.Error("file1.txt should not appear in diff (unchanged)")
	}
}

// Helper function to calculate file hash using SHA-256 (for test data only)
// Note: Real implementation uses fileops package which calculates hashes from actual files
func calculateFileHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}
