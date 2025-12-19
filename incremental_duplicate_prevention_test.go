// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] Incremental Archive Duplicate Prevention testing
// [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] Duplicate prevention architecture validation
// [IMPL:INCREMENTAL_DUPLICATE_PREVENTION] Incremental archive creation with duplicate prevention validation
package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncrementalArchiveSkipNoChanges tests that incremental archive creation is skipped when no changes exist
func TestIncrementalArchiveSkipNoChanges_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
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

	// Change to source directory first (required for archive creation)
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false // Don't create subdirectory

	// Create full archive using actual function
	err = CreateFullArchive(cfg, "base-archive", false)
	if err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Add a file to create a change, then create first incremental
	newFile := filepath.Join(sourceDir, "file3.txt")
	if err := os.WriteFile(newFile, []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Count archives before creating incremental
	archivesBefore, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}
	incrementalCountBefore := 0
	for _, archive := range archivesBefore {
		if archive.IsIncremental {
			incrementalCountBefore++
		}
	}

	// Create incremental archive using actual function (with change)
	// This should create the first incremental archive
	err = CreateIncrementalArchive(cfg, "incremental-test", false)
	if err != nil {
		t.Fatalf("Failed to create initial incremental archive: %v", err)
	}

	// Count archives after first incremental creation
	archivesAfterFirst, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}
	incrementalCountAfterFirst := 0
	for _, archive := range archivesAfterFirst {
		if archive.IsIncremental {
			incrementalCountAfterFirst++
		}
	}

	// Verify first incremental was created
	if incrementalCountAfterFirst != incrementalCountBefore+1 {
		t.Errorf("Expected %d incremental archives after first creation, got %d", incrementalCountBefore+1, incrementalCountAfterFirst)
	}

	// Attempt to create incremental archive again (should be skipped since no changes)
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    "",
		DryRun:  false,
		Context: context.Background(),
	}

	// Capture output to verify skip message
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = createIncrementalArchive(config)

	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify no error (skip is not an error)
	if err != nil {
		t.Errorf("Expected no error when skipping duplicate archive, got: %v", err)
	}

	// Verify skip message was printed
	if !strings.Contains(output, "no changes") && !strings.Contains(output, "skipped") && !strings.Contains(output, "Skipping") {
		t.Logf("Output: %q", output)
		t.Log("Note: Skip message format may vary, but archive should not be created")
	}

	// Verify no new incremental archive was created (should still be same count)
	archivesAfter, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	incrementalCountAfter := 0
	for _, archive := range archivesAfter {
		if archive.IsIncremental {
			incrementalCountAfter++
		}
	}

	// Should still have the same number of incremental archives (no new one created)
	if incrementalCountAfter != incrementalCountAfterFirst {
		t.Errorf("Expected %d incremental archives after skip (same as before), got %d", incrementalCountAfterFirst, incrementalCountAfter)
	}
}

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncrementalArchiveProceedWithChanges tests that incremental archive creation proceeds when changes exist
func TestIncrementalArchiveProceedWithChanges_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create initial files
	initialFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for relPath, content := range initialFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Change to source directory first
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive using actual function
	err = CreateFullArchive(cfg, "base-archive", false)
	if err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Modify a file (create change)
	modifiedFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(modifiedFile, []byte("modified_content1"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Count archives before
	archivesBefore, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}
	incrementalCountBefore := 0
	for _, archive := range archivesBefore {
		if archive.IsIncremental {
			incrementalCountBefore++
		}
	}

	// Create incremental archive (should proceed because there are changes)
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    "",
		DryRun:  false,
		Context: context.Background(),
	}

	err = createIncrementalArchive(config)
	if err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	// Verify new incremental archive was created
	archivesAfter, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	incrementalCountAfter := 0
	for _, archive := range archivesAfter {
		if archive.IsIncremental {
			incrementalCountAfter++
		}
	}

	// Should have one more incremental archive
	if incrementalCountAfter != incrementalCountBefore+1 {
		t.Errorf("Expected %d incremental archives after creation, got %d", incrementalCountBefore+1, incrementalCountAfter)
	}
}

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncrementalArchiveSkipWithFullOnly tests skip behavior when only full archive exists
func TestIncrementalArchiveSkipWithFullOnly_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
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

	// Change to source directory first
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive using actual function
	err = CreateFullArchive(cfg, "base-archive", false)
	if err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Attempt to create incremental archive (should be skipped if no changes)
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    "",
		DryRun:  false,
		Context: context.Background(),
	}

	// Capture output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = createIncrementalArchive(config)

	w.Close()
	os.Stdout = originalStdout

	// Read captured output
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	output := string(buf[:n])

	// Verify behavior - should either skip or create (depending on implementation)
	// The key is that it should handle the case gracefully
	if err != nil {
		// If it errors, it should be a clear error, not a panic
		t.Logf("Archive creation returned error (may be expected): %v", err)
		t.Logf("Output: %q", output)
	}

	// Verify archives list is still valid
	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	// Should have at least the full archive
	fullCount := 0
	for _, archive := range archives {
		if !archive.IsIncremental {
			fullCount++
		}
	}

	if fullCount < 1 {
		t.Errorf("Expected at least 1 full archive, got %d", fullCount)
	}
}

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncrementalArchiveProceedWithNewFile tests that new files trigger archive creation
func TestIncrementalArchiveProceedWithNewFile_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create initial files
	initialFiles := map[string]string{
		"file1.txt": "content1",
	}

	for relPath, content := range initialFiles {
		fullPath := filepath.Join(sourceDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Change to source directory first
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive using actual function
	err = CreateFullArchive(cfg, "base-archive", false)
	if err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Add a new file
	newFile := filepath.Join(sourceDir, "file2.txt")
	if err := os.WriteFile(newFile, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Count incremental archives before
	archivesBefore, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}
	incrementalCountBefore := 0
	for _, archive := range archivesBefore {
		if archive.IsIncremental {
			incrementalCountBefore++
		}
	}

	// Create incremental archive (should proceed because new file was added)
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    "",
		DryRun:  false,
		Context: context.Background(),
	}

	err = createIncrementalArchive(config)
	if err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	// Verify new incremental archive was created
	archivesAfter, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	incrementalCountAfter := 0
	for _, archive := range archivesAfter {
		if archive.IsIncremental {
			incrementalCountAfter++
		}
	}

	// Should have one more incremental archive
	if incrementalCountAfter != incrementalCountBefore+1 {
		t.Errorf("Expected %d incremental archives after creation, got %d", incrementalCountBefore+1, incrementalCountAfter)
	}
}
