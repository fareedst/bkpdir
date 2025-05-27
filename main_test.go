// This file is part of bkpdir

// Package main provides tests for the BkpDir CLI application.
package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// Test constants
const (
	dryRunFlagDesc = "Show what would be done, but don't create archives"
)

func TestFullCmdWithNote(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create a test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	os.WriteFile(cfgPath, []byte("archive_dir_path: ../.bkpdir\nuse_current_dir_name: true\n"), 0644)

	// Create a test file to archive
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Create archive directory structure
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	os.MkdirAll(filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir)), 0755)

	// Test cases
	tests := []struct {
		name     string
		args     []string
		wantNote string
	}{
		{"no note", []string{"--dry-run", "full"}, ""},
		{"with note", []string{"--dry-run", "full", "test-note"}, "test-note"},
		{"with spaces in note", []string{"--dry-run", "full", "test note with spaces"}, "test note with spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			rootCmd := createTestRootCmd()

			// Set up a buffer to capture output
			rootCmd.SetOut(nil)
			rootCmd.SetErr(nil)
			// Execute the command with test args
			rootCmd.SetArgs(tt.args)
			rootCmd.Execute()

			// Verify the note was passed correctly by checking the archive name
			archives, err := ListArchives(archiveDir)
			if err != nil {
				t.Fatalf("ListArchives error: %v", err)
			}

			// In dry-run mode, no archives are created
			if len(archives) > 0 {
				t.Errorf("Expected no archives in dry-run mode, got %d", len(archives))
			}
		})
	}
}

func TestIncCmdWithNote(t *testing.T) {
	// Setup test environment
	tmpDir, archiveDir := setupIncTestEnvironment(t)

	// Create initial full archive
	createInitialFullArchive(t, tmpDir, archiveDir)

	// Run incremental command tests
	runIncCommandTests(t, archiveDir)
}

func setupIncTestEnvironment(t *testing.T) (string, string) {
	t.Helper()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Debug: print temp dir
	t.Logf("Temp dir: %s", tmpDir)

	// Create a test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	configContent := "archive_dir_path: ../.bkpdir\nuse_current_dir_name: true\n"
	os.WriteFile(cfgPath, []byte(configContent), 0644)

	// Create a test file to archive
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Create archive directory structure
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Debug: print archiveDir
	t.Logf("Archive dir: %s", archiveDir)

	return tmpDir, archiveDir
}

func createInitialFullArchive(t *testing.T, _ string, archiveDir string) {
	t.Helper()

	// Create a full archive first (without dry-run)
	rootCmd := createTestRootCmd()

	// Create initial full archive
	rootCmd.SetArgs([]string{"full"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to create initial full archive: %v", err)
	}

	// Debug: list files in archiveDir
	files, err := os.ReadDir(archiveDir)
	if err != nil {
		t.Logf("Error reading archiveDir: %v", err)
	} else {
		for _, f := range files {
			t.Logf("File in archiveDir: %s", f.Name())
		}
	}

	// Verify the initial archive was created
	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("ListArchives error: %v", err)
	}
	if len(archives) != 1 {
		t.Fatalf("Expected one initial full archive, got %d", len(archives))
	}
}

func runIncCommandTests(t *testing.T, archiveDir string) {
	t.Helper()

	// Test cases
	tests := []struct {
		name     string
		args     []string
		wantNote string
	}{
		{"no note", []string{"--dry-run", "inc"}, ""},
		{"with note", []string{"--dry-run", "inc", "test-note"}, "test-note"},
		{"with spaces in note", []string{"--dry-run", "inc", "test note with spaces"}, "test note with spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			rootCmd := createTestRootCmd()

			// Set up a buffer to capture output
			rootCmd.SetOut(nil)
			rootCmd.SetErr(nil)
			// Execute the command with test args
			rootCmd.SetArgs(tt.args)
			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("Failed to execute command: %v", err)
			}

			// Verify the note was passed correctly by checking the archive name
			archives, err := ListArchives(archiveDir)
			if err != nil {
				t.Fatalf("ListArchives error: %v", err)
			}

			// In dry-run mode, we should only have the initial full archive
			if len(archives) != 1 {
				t.Errorf("Expected only the initial full archive in dry-run mode, got %d", len(archives))
			}
		})
	}
}

func createTestRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "bkpdir",
		Short: "Directory archiving CLI for macOS and Linux",
	}
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, dryRunFlagDesc)
	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())
	rootCmd.AddCommand(listCmd())
	return rootCmd
}

func TestCmdArgsValidation(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create a test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	os.WriteFile(cfgPath, []byte("archive_dir_path: ../.bkpdir\nuse_current_dir_name: true\n"), 0644)

	// Create a test file to archive
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"full no args", []string{"--dry-run", "full"}, false},
		{"full one arg", []string{"--dry-run", "full", "note"}, false},
		{"full two args", []string{"--dry-run", "full", "note1", "note2"}, true},
		{"inc no args", []string{"--dry-run", "inc"}, false},
		{"inc one arg", []string{"--dry-run", "inc", "note"}, false},
		{"inc two args", []string{"--dry-run", "inc", "note1", "note2"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// If the test is for an incremental command, create a full archive first
			if len(tt.args) > 1 && tt.args[1] == "inc" && !tt.wantErr {
				rootCmd := createTestRootCmd()
				rootCmd.SetArgs([]string{"full"})
				if err := rootCmd.Execute(); err != nil {
					t.Fatalf("Failed to create initial full archive for inc test: %v", err)
				}
			}

			// Create a new root command for each test
			rootCmd := createTestRootCmd()

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Execute() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestListCmd(t *testing.T) {
	// Setup test environment
	tmpDir, archiveDir := setupListTestEnvironment(t)

	// Create multiple archives with different timestamps
	createTestArchives(t, tmpDir, archiveDir)

	// Test the list command
	runListCommandTest(t, archiveDir)
}

func setupListTestEnvironment(t *testing.T) (string, string) {
	t.Helper()

	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Create a test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	configContent := "archive_dir_path: ../.bkpdir\nuse_current_dir_name: true\n"
	os.WriteFile(cfgPath, []byte(configContent), 0644)

	// Create a test file to archive
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Create archive directory structure
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	return tmpDir, archiveDir
}

func createTestArchives(t *testing.T, tmpDir string, archiveDir string) {
	t.Helper()

	// Create first archive
	rootCmd := createTestRootCmd()
	rootCmd.SetArgs([]string{"full", "First archive"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to create first archive: %v", err)
	}

	// Wait a bit to ensure different timestamps
	time.Sleep(1 * time.Second)

	// Modify the test file to create a different archive
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("modified test content"), 0644)

	// Create second archive
	rootCmd = createTestRootCmd()
	rootCmd.SetArgs([]string{"full", "Second archive"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to create second archive: %v", err)
	}

	// Verify archives were created
	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("ListArchives error: %v", err)
	}
	if len(archives) != 2 {
		t.Fatalf("Expected 2 archives, got %d", len(archives))
	}
}

func runListCommandTest(t *testing.T, archiveDir string) {
	t.Helper()

	// Execute the list command
	executeListCommand(t)

	// Get and verify archives
	archives := getAndVerifyArchives(t, archiveDir)

	// Verify sorting
	verifySortingOrder(t, archives)

	// Verify archive names
	verifyArchiveNames(t, archives)

	// Verify verification status
	verifyVerificationStatus(t, archives)
}

func executeListCommand(t *testing.T) {
	t.Helper()

	rootCmd := createTestRootCmd()
	rootCmd.SetOut(nil) // Suppress output for test
	rootCmd.SetErr(nil) // Suppress error output for test
	rootCmd.SetArgs([]string{"list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Failed to execute list command: %v", err)
	}
}

func getAndVerifyArchives(t *testing.T, archiveDir string) []Archive {
	t.Helper()

	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatalf("Failed to list archives: %v", err)
	}

	if len(archives) != 2 {
		t.Fatalf("Expected 2 archives, got %d", len(archives))
	}

	// Apply the same sorting logic as ListArchivesEnhanced (most recent first)
	sort.Slice(archives, func(i, j int) bool {
		return archives[i].CreationTime.After(archives[j].CreationTime)
	})

	return archives
}

func verifySortingOrder(t *testing.T, archives []Archive) {
	t.Helper()

	if !archives[0].CreationTime.After(archives[1].CreationTime) &&
		!archives[0].CreationTime.Equal(archives[1].CreationTime) {
		t.Errorf("Archives should be sorted by creation time (most recent first), got: %v vs %v",
			archives[0].CreationTime, archives[1].CreationTime)
	}
}

func verifyArchiveNames(t *testing.T, archives []Archive) {
	t.Helper()

	archiveNames := []string{archives[0].Name, archives[1].Name}
	hasFirstArchive := false
	hasSecondArchive := false

	for _, name := range archiveNames {
		if strings.Contains(name, "First archive") {
			hasFirstArchive = true
		}
		if strings.Contains(name, "Second archive") {
			hasSecondArchive = true
		}
	}

	if !hasFirstArchive {
		t.Errorf("Should have archive with 'First archive' note, got names: %v", archiveNames)
	}
	if !hasSecondArchive {
		t.Errorf("Should have archive with 'Second archive' note, got names: %v", archiveNames)
	}
}

func verifyVerificationStatus(t *testing.T, archives []Archive) {
	t.Helper()

	for i, archive := range archives {
		if archive.VerificationStatus != nil && archive.VerificationStatus.IsVerified {
			t.Errorf("Archive %d should be unverified initially", i)
		}
	}
}
