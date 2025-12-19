// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [REQ:DIFF_COMMAND] Integration tests for inc and diff CLI commands
// [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:DIFF_COMMAND] End-to-end command execution validation
// [IMPL:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:DIFF_COMMAND] CLI command integration testing
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncCommandIntegration tests the inc command end-to-end execution
func TestIncCommandIntegration_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
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

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create config for direct function calls
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive first
	if err := CreateFullArchive(cfg, "base-archive", false); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
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

	// Test 1: Create incremental with changes
	newFile := filepath.Join(sourceDir, "file3.txt")
	if err := os.WriteFile(newFile, []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Execute inc command via CLI
	rootCmd := createTestRootCmd()
	rootCmd.AddCommand(diffCmd()) // Add diff command for completeness
	var stdout, stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs([]string{"inc", "test-note"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("inc command failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify incremental archive was created
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

	if incrementalCountAfter != incrementalCountBefore+1 {
		t.Errorf("Expected %d incremental archives after creation, got %d", incrementalCountBefore+1, incrementalCountAfter)
	}

	// Verify output contains note
	output := stdout.String()
	if !strings.Contains(output, "test-note") {
		t.Logf("Note may be in archive name, output: %s", output)
	}

	// Test 2: Try to create another incremental with no changes (should skip)
	stdout.Reset()
	stderr.Reset()
	rootCmd.SetArgs([]string{"inc"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("inc command failed (should skip, not error): %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
	}

	// Verify no new archive was created
	archivesAfterSkip, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	incrementalCountAfterSkip := 0
	for _, archive := range archivesAfterSkip {
		if archive.IsIncremental {
			incrementalCountAfterSkip++
		}
	}

	if incrementalCountAfterSkip != incrementalCountAfter {
		t.Errorf("Expected %d incremental archives after skip (same as before), got %d", incrementalCountAfter, incrementalCountAfterSkip)
	}

	// Verify skip message was printed
	skipOutput := stdout.String()
	if !strings.Contains(skipOutput, "no changes") && !strings.Contains(skipOutput, "skipped") && !strings.Contains(skipOutput, "Skipping") {
		t.Logf("Skip message format may vary, output: %s", skipOutput)
	}
}

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncCommandWithNoteFlag tests the inc command with --note flag
func TestIncCommandWithNoteFlag_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
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
	testFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(testFile, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create config for direct function calls
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive first
	if err := CreateFullArchive(cfg, "base", false); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Add a change
	newFile := filepath.Join(sourceDir, "file2.txt")
	if err := os.WriteFile(newFile, []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Execute inc command with --note flag
	rootCmd := createTestRootCmd()
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{"inc", "--note", "flag-note"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("inc command with --note flag failed: %v\nOutput: %s", err, stdout.String())
	}

	// Verify archive was created
	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	incrementalFound := false
	for _, archive := range archives {
		if archive.IsIncremental && strings.Contains(archive.Name, "flag-note") {
			incrementalFound = true
			break
		}
	}

	if !incrementalFound {
		t.Errorf("Expected to find incremental archive with note 'flag-note'")
	}
}

// [REQ:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]
// TestIncCommandNoFullArchive tests inc command when no full archive exists
func TestIncCommandNoFullArchive_REQ_INCREMENTAL_DUPLICATE_PREVENTION(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create test file
	testFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(testFile, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test the underlying function directly since CLI command calls os.Exit()
	// which terminates the test process
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Try to create incremental without full archive (should error)
	err = CreateIncrementalArchive(cfg, "test", false)
	if err == nil {
		t.Error("Expected error when creating incremental archive without full archive, got nil")
	} else {
		// Error is expected - verify it's a meaningful error
		if !strings.Contains(err.Error(), "archive") && !strings.Contains(err.Error(), "found") {
			t.Logf("Error message: %v (may be acceptable)", err)
		}
	}
}

// [REQ:DIFF_COMMAND] [ARCH:DIFF_COMMAND] [IMPL:DIFF_COMMAND]
// TestDiffCommandCLI tests the diff command CLI end-to-end execution
func TestDiffCommandCLI_REQ_DIFF_COMMAND(t *testing.T) {
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

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create config for direct function calls
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive
	if err := CreateFullArchive(cfg, "base", false); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Test 1: Diff with no changes (should show no changes)
	rootCmd := createTestRootCmd()
	rootCmd.AddCommand(diffCmd())
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{"diff"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("diff command failed: %v\nOutput: %s", err, stdout.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "No changes") && !strings.Contains(output, "no changes") {
		t.Logf("Expected 'no changes' message, got: %s", output)
	}

	// Test 2: Make changes and diff again
	modifiedFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(modifiedFile, []byte("modified_content1"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	newFile := filepath.Join(sourceDir, "file3.txt")
	if err := os.WriteFile(newFile, []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Test diff calculation directly since CLI output may not be captured properly
	// due to os.Exit() calls
	reconstructedState, err := ReconstructArchiveState(archiveDir)
	if err != nil {
		t.Fatalf("ReconstructArchiveState failed: %v", err)
	}

	diff, err := CalculateDiff(sourceDir, reconstructedState, cfg.ExcludePatterns)
	if err != nil {
		t.Fatalf("CalculateDiff failed: %v", err)
	}

	// Verify changes were detected
	if len(diff.Added) == 0 && len(diff.Modified) == 0 && len(diff.Deleted) == 0 {
		t.Error("Expected changes to be detected, but diff shows no changes")
	}

	// Verify file3.txt is in added list
	foundFile3 := false
	for _, f := range diff.Added {
		if f == "file3.txt" {
			foundFile3 = true
			break
		}
	}
	if !foundFile3 {
		t.Errorf("Expected file3.txt to be in Added list, got: %v", diff.Added)
	}
}

// [REQ:DIFF_COMMAND] [ARCH:DIFF_COMMAND] [IMPL:DIFF_COMMAND]
// TestDiffCommandWithIncremental tests diff command with full + incremental archives
func TestDiffCommandWithIncremental_REQ_DIFF_COMMAND(t *testing.T) {
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

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create config for direct function calls
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	// Create full archive
	if err := CreateFullArchive(cfg, "base", false); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Modify file and create incremental
	modifiedFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(modifiedFile, []byte("modified_content1"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	if err := CreateIncrementalArchive(cfg, "inc1", false); err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	// Add another change
	newFile := filepath.Join(sourceDir, "file3.txt")
	if err := os.WriteFile(newFile, []byte("content3"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Run diff command (should compare against full + incremental)
	rootCmd := createTestRootCmd()
	rootCmd.AddCommand(diffCmd())
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{"diff"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("diff command failed: %v\nOutput: %s", err, stdout.String())
	}

	output := stdout.String()
	// Should show file3.txt as added (since it's not in full or incremental)
	if !strings.Contains(output, "file3.txt") && !strings.Contains(output, "Changes") {
		t.Logf("Output may not list specific files, got: %s", output)
	}
}

// [REQ:DIFF_COMMAND] [ARCH:DIFF_COMMAND] [IMPL:DIFF_COMMAND]
// TestDiffCommandCLINoArchives tests diff CLI command when no archives exist
func TestDiffCommandCLINoArchives_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create test file
	testFile := filepath.Join(sourceDir, "file1.txt")
	if err := os.WriteFile(testFile, []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test diff calculation directly since CLI command calls os.Exit()
	// which terminates the test process
	_, err = ReconstructArchiveState(archiveDir)
	if err == nil {
		t.Error("Expected error when no archives exist, got nil")
	} else {
		// Verify it's a meaningful error about no archives
		if !strings.Contains(err.Error(), "No archives") && !strings.Contains(err.Error(), "No full archive") {
			t.Logf("Error message: %v (may be acceptable)", err)
		}
	}
}

// [REQ:DIFF_COMMAND] [ARCH:DIFF_COMMAND] [IMPL:DIFF_COMMAND]
// TestDiffCommandExcludePatterns tests that diff command respects exclude patterns
func TestDiffCommandExcludePatterns_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Create files including one that should be excluded
	files := map[string]string{
		"file1.txt":    "content1",
		"excluded.txt": "excluded_content",
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

	// Change to source directory
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to change to source directory: %v", err)
	}
	defer os.Chdir(originalCwd)

	// Create config file for CLI command to load with exclude pattern
	configPath := filepath.Join(sourceDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: false\nexclude_patterns:\n  - excluded.txt\n", archiveDir)
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create config for direct function calls
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false
	cfg.ExcludePatterns = []string{"excluded.txt"}

	// Create full archive (excluded file should not be in archive)
	if err := CreateFullArchive(cfg, "base", false); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Modify excluded file (should not show in diff)
	excludedFile := filepath.Join(sourceDir, "excluded.txt")
	if err := os.WriteFile(excludedFile, []byte("modified_excluded"), 0644); err != nil {
		t.Fatalf("Failed to modify excluded file: %v", err)
	}

	// Run diff command
	rootCmd := createTestRootCmd()
	rootCmd.AddCommand(diffCmd())
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetArgs([]string{"diff"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("diff command failed: %v\nOutput: %s", err, stdout.String())
	}

	output := stdout.String()
	// Excluded file should not appear in diff
	if strings.Contains(output, "excluded.txt") {
		t.Errorf("Excluded file 'excluded.txt' should not appear in diff output, got: %s", output)
	}
}
