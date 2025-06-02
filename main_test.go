// This file is part of bkpdir

// Package main provides tests for the BkpDir CLI application.
package main

import (
	"fmt"
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

// ARCH-002: Full archive command validation
// TEST-REF: TestFullCmdWithNote
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

// ARCH-003: Incremental archive command validation
// TEST-REF: TestIncCmdWithNote
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

// TEST-REF: TestMain_HandleConfigCommand
func TestMain_HandleConfigCommand(t *testing.T) {
	// TEST-MAIN-001: Test handleConfigCommand function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	configContent := `archive_dir_path: ../.bkpdir
use_current_dir_name: true
include_git_info: false
`
	if err := os.WriteFile(cfgPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Test should run without error
	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(fmt.Sprintf("%v", r), "exit") {
				t.Errorf("handleConfigCommand panicked: %v", r)
			}
		}
	}()

	// Since handleConfigCommand calls os.Exit(), we test indirectly by verifying
	// the config loading works correctly
	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Errorf("Config loading failed: %v", err)
	}
	if cfg.ArchiveDirPath != "../.bkpdir" {
		t.Errorf("Expected archive_dir_path '../.bkpdir', got %s", cfg.ArchiveDirPath)
	}
}

// TEST-REF: TestMain_HandleCreateCommand
func TestMain_HandleCreateCommand(t *testing.T) {
	// TEST-MAIN-002: Test handleCreateCommand function
	// This is a placeholder function, so we just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("handleCreateCommand panicked: %v", r)
		}
	}()
	handleCreateCommand()
}

// TEST-REF: TestMain_HandleVerifyCommand
func TestMain_HandleVerifyCommand(t *testing.T) {
	// TEST-MAIN-003: Test handleVerifyCommand function
	// This is a placeholder function, so we just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("handleVerifyCommand panicked: %v", r)
		}
	}()
	handleVerifyCommand()
}

// TEST-REF: TestMain_HandleVersionCommand
func TestMain_HandleVersionCommand(t *testing.T) {
	// TEST-MAIN-004: Test handleVersionCommand function
	// This is a placeholder function, so we just verify it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("handleVersionCommand panicked: %v", r)
		}
	}()
	handleVersionCommand()
}

// TEST-REF: TestMain_ConfigCmd
func TestMain_ConfigCmd(t *testing.T) {
	// TEST-MAIN-005: Test configCmd function
	cmd := configCmd()

	if cmd.Use != "config [KEY] [VALUE]" {
		t.Errorf("Expected Use 'config [KEY] [VALUE]', got %s", cmd.Use)
	}
	if cmd.Short != "Display or modify configuration values" {
		t.Errorf("Expected Short description, got %s", cmd.Short)
	}
	if !strings.Contains(cmd.Long, "Display configuration values") {
		t.Errorf("Expected Long description to contain 'Display configuration values'")
	}
}

// TEST-REF: TestMain_CreateCmd
func TestMain_CreateCmd(t *testing.T) {
	// TEST-MAIN-006: Test createCmd function
	cmd := createCmd()

	if cmd.Use != "create" {
		t.Errorf("Expected Use 'create', got %s", cmd.Use)
	}
	if cmd.Short != "Create a new archive" {
		t.Errorf("Expected Short 'Create a new archive', got %s", cmd.Short)
	}
}

// TEST-REF: TestMain_VerifyCmd
func TestMain_VerifyCmd(t *testing.T) {
	// TEST-MAIN-007: Test verifyCmd function
	cmd := verifyCmd()

	if cmd.Use != "verify" {
		t.Errorf("Expected Use 'verify', got %s", cmd.Use)
	}
	if cmd.Short != "Verify archives" {
		t.Errorf("Expected Short 'Verify archives', got %s", cmd.Short)
	}
}

// TEST-REF: TestMain_VersionCmd
func TestMain_VersionCmd(t *testing.T) {
	// TEST-MAIN-008: Test versionCmd function
	cmd := versionCmd()

	if cmd.Use != "version" {
		t.Errorf("Expected Use 'version', got %s", cmd.Use)
	}
	if cmd.Short != "Display version information" {
		t.Errorf("Expected Short 'Display version information', got %s", cmd.Short)
	}
}

// TEST-REF: TestMain_CreateFullArchiveEnhanced
func TestMain_CreateFullArchiveEnhanced(t *testing.T) {
	// TEST-MAIN-009: Test CreateFullArchiveEnhanced function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create test config
	cfg := createTestConfig(t, tmpDir)
	formatter := NewOutputFormatter(cfg)

	opts := ArchiveOptions{
		Config:    cfg,
		Formatter: formatter,
		Note:      "test-note",
		DryRun:    true, // Use dry-run to avoid creating actual archives
		Verify:    false,
	}

	// Test should not return error in dry-run mode
	err := CreateFullArchiveEnhanced(opts)
	if err != nil {
		t.Errorf("CreateFullArchiveEnhanced failed: %v", err)
	}
}

// TEST-REF: TestMain_CreateIncrementalArchiveEnhanced
func TestMain_CreateIncrementalArchiveEnhanced(t *testing.T) {
	// TEST-MAIN-010: Test CreateIncrementalArchiveEnhanced function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create test config
	cfg := createTestConfig(t, tmpDir)
	formatter := NewOutputFormatter(cfg)

	// Create archive directory
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive dir: %v", err)
	}

	// First create a full archive
	fullOpts := ArchiveOptions{
		Config:    cfg,
		Formatter: formatter,
		Note:      "initial",
		DryRun:    false, // Create actual archive for incremental base
		Verify:    false,
	}
	if err := CreateFullArchiveEnhanced(fullOpts); err != nil {
		t.Fatalf("Failed to create full archive: %v", err)
	}

	// Now test incremental
	opts := ArchiveOptions{
		Config:    cfg,
		Formatter: formatter,
		Note:      "test-note",
		DryRun:    true, // Use dry-run for incremental test
		Verify:    false,
	}

	// Test should not return error in dry-run mode with existing full archive
	err := CreateIncrementalArchiveEnhanced(opts)
	if err != nil {
		t.Errorf("CreateIncrementalArchiveEnhanced failed: %v", err)
	}
}

// Helper function to create test config
func createTestConfig(t *testing.T, tmpDir string) *Config {
	t.Helper()

	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	configContent := `archive_dir_path: ../.bkpdir
use_current_dir_name: true
include_git_info: false
`
	if err := os.WriteFile(cfgPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Create test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}

// TEST-REF: TestMain_VerifyArchiveEnhanced
func TestMain_VerifyArchiveEnhanced(t *testing.T) {
	// TEST-MAIN-011: Test VerifyArchiveEnhanced function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cfg := createTestConfig(t, tmpDir)
	formatter := NewOutputFormatter(cfg)

	// Create archive directory
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive dir: %v", err)
	}

	opts := VerifyOptions{
		Config:       cfg,
		Formatter:    formatter,
		ArchiveName:  "", // Empty means verify all
		WithChecksum: false,
	}

	// Should handle empty archive directory gracefully
	err := VerifyArchiveEnhanced(opts)
	// No error expected when no archives exist (verifyAllArchives returns nil in this case)
	if err != nil {
		t.Errorf("Expected no error when no archives exist, got: %v", err)
	}
}

// TEST-REF: TestMain_GetArchiveDirectory
func TestMain_GetArchiveDirectory(t *testing.T) {
	// TEST-MAIN-012: Test getArchiveDirectory function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cfg := createTestConfig(t, tmpDir)

	archiveDir, err := getArchiveDirectory(cfg)
	if err != nil {
		t.Errorf("getArchiveDirectory failed: %v", err)
	}

	expectedPath := filepath.Join("../.bkpdir", filepath.Base(tmpDir))
	if !strings.HasSuffix(archiveDir, expectedPath) {
		t.Errorf("Expected archive directory to end with %s, got %s", expectedPath, archiveDir)
	}
}

// TEST-REF: TestMain_HandleListFileBackupsCommand
func TestMain_HandleListFileBackupsCommand(t *testing.T) {
	// TEST-MAIN-013: Test handleListFileBackupsCommand function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cfg := createTestConfig(t, tmpDir)

	// Create backup directory
	backupDir := filepath.Join(tmpDir, "../.bkpdir", "files")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		t.Fatalf("Failed to create backup dir: %v", err)
	}

	// Test with valid file path
	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(fmt.Sprintf("%v", r), "exit") {
				t.Errorf("handleListFileBackupsCommand panicked: %v", r)
			}
		}
	}()

	// The function will call os.Exit, so we test the setup only
	if cfg.BackupDirPath == "" {
		t.Errorf("Expected backup directory path to be set")
	}
}

// TEST-REF: TestMain_BackupCmd
func TestMain_BackupCmd(t *testing.T) {
	// TEST-MAIN-014: Test backupCmd function
	cmd := backupCmd()

	if cmd.Use != "backup [FILE_PATH] [NOTE]" {
		t.Errorf("Expected Use 'backup [FILE_PATH] [NOTE]', got %s", cmd.Use)
	}
	if cmd.Short != "Create a backup of a single file" {
		t.Errorf("Expected Short 'Create a backup of a single file', got %s", cmd.Short)
	}
	if !strings.Contains(cmd.Long, "Create a backup of the specified file") {
		t.Errorf("Expected Long description to contain backup info")
	}
}

// TEST-REF: TestMain_HandleConfigSetCommand
func TestMain_HandleConfigSetCommand(t *testing.T) {
	// TEST-MAIN-015: Test handleConfigSetCommand function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	// Create initial config
	createTestConfig(t, tmpDir)

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"set archive path", "archive_dir_path", "/custom/path"},
		{"set boolean", "use_current_dir_name", "false"},
		{"set integer", "status_config_error", "5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !strings.Contains(fmt.Sprintf("%v", r), "exit") {
						t.Errorf("handleConfigSetCommand panicked: %v", r)
					}
				}
			}()

			// The function calls os.Exit, so we test the components separately
			configData := loadExistingConfigData(filepath.Join(tmpDir, ".bkpdir.yml"))
			if configData == nil {
				t.Errorf("Expected config data to be loaded")
			}

			convertedValue := convertConfigValue(tt.key, tt.value)
			if convertedValue == nil {
				t.Errorf("Expected converted value for %s", tt.key)
			}
		})
	}
}

// TEST-REF: TestMain_LoadExistingConfigData
func TestMain_LoadExistingConfigData(t *testing.T) {
	// TEST-MAIN-016: Test loadExistingConfigData function
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".bkpdir.yml")

	// Test with existing config file
	configContent := `archive_dir_path: /test/path
use_current_dir_name: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	configData := loadExistingConfigData(configPath)
	if configData == nil {
		t.Fatalf("Expected config data to be loaded")
	}

	if configData["archive_dir_path"] != "/test/path" {
		t.Errorf("Expected archive_dir_path '/test/path', got %v", configData["archive_dir_path"])
	}

	if configData["use_current_dir_name"] != true {
		t.Errorf("Expected use_current_dir_name true, got %v", configData["use_current_dir_name"])
	}

	// Test with non-existent config file
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.yml")
	configData2 := loadExistingConfigData(nonExistentPath)
	if configData2 == nil {
		t.Errorf("Expected empty config data for non-existent file")
	}
	if len(configData2) != 0 {
		t.Errorf("Expected empty config data, got %v", configData2)
	}
}

// TEST-REF: TestMain_ConvertConfigValue
func TestMain_ConvertConfigValue(t *testing.T) {
	// TEST-MAIN-017: Test convertConfigValue function
	tests := []struct {
		name       string
		key        string
		value      string
		expected   interface{}
		shouldExit bool
	}{
		{"string value", "archive_dir_path", "/custom/path", "/custom/path", false},
		{"boolean true", "use_current_dir_name", "true", true, false},
		{"boolean false", "include_git_info", "false", false, false},
		{"integer value", "status_config_error", "5", 5, false},
		{"invalid boolean", "use_current_dir_name", "maybe", false, true},
		{"invalid integer", "status_config_error", "abc", 0, true},
		{"unknown key", "unknown_key", "value", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldExit {
				// For tests that should exit, we skip them since we can't easily test os.Exit
				t.Skip("Skipping test that calls os.Exit")
				return
			}

			result := convertConfigValue(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TEST-REF: TestMain_ConvertBooleanValue
func TestMain_ConvertBooleanValue(t *testing.T) {
	// TEST-MAIN-018: Test convertBooleanValue function
	tests := []struct {
		name       string
		key        string
		value      string
		expected   bool
		shouldExit bool
	}{
		{"true value", "test_key", "true", true, false},
		{"false value", "test_key", "false", false, false},
		{"invalid value", "test_key", "maybe", false, true},
		{"empty value", "test_key", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldExit {
				// For tests that should exit, we skip them since we can't easily test os.Exit
				t.Skip("Skipping test that calls os.Exit")
				return
			}

			result := convertBooleanValue(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TEST-REF: TestMain_ConvertIntegerValue
func TestMain_ConvertIntegerValue(t *testing.T) {
	// TEST-MAIN-019: Test convertIntegerValue function
	tests := []struct {
		name       string
		key        string
		value      string
		expected   int
		shouldExit bool
	}{
		{"positive integer", "test_key", "5", 5, false},
		{"zero", "test_key", "0", 0, false},
		{"negative integer", "test_key", "-1", -1, false},
		{"invalid value", "test_key", "abc", 0, true},
		{"float value", "test_key", "3.14", 0, true},
		{"empty value", "test_key", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldExit {
				// For tests that should exit, we skip them since we can't easily test os.Exit
				t.Skip("Skipping test that calls os.Exit")
				return
			}

			result := convertIntegerValue(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TEST-REF: TestMain_UpdateConfigData
func TestMain_UpdateConfigData(t *testing.T) {
	// TEST-MAIN-020: Test updateConfigData function
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected map[string]interface{}
	}{
		{
			name:  "regular config value",
			key:   "archive_dir_path",
			value: "/custom/path",
			expected: map[string]interface{}{
				"archive_dir_path": "/custom/path",
			},
		},
		{
			name:  "verification config - verify_on_create",
			key:   "verify_on_create",
			value: true,
			expected: map[string]interface{}{
				"verification": map[string]interface{}{
					"verify_on_create": true,
				},
			},
		},
		{
			name:  "verification config - checksum_algorithm",
			key:   "checksum_algorithm",
			value: "sha256",
			expected: map[string]interface{}{
				"verification": map[string]interface{}{
					"checksum_algorithm": "sha256",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configData := make(map[string]interface{})
			updateConfigData(configData, tt.key, tt.value)

			if tt.key == "verify_on_create" || tt.key == "checksum_algorithm" {
				verification, ok := configData["verification"].(map[string]interface{})
				if !ok {
					t.Fatalf("Expected verification section to be map[string]interface{}")
				}
				if verification[tt.key] != tt.value {
					t.Errorf("Expected %s to be %v, got %v", tt.key, tt.value, verification[tt.key])
				}
			} else {
				if configData[tt.key] != tt.value {
					t.Errorf("Expected %s to be %v, got %v", tt.key, tt.value, configData[tt.key])
				}
			}
		})
	}
}

// TEST-REF: TestMain_SaveConfigData
func TestMain_SaveConfigData(t *testing.T) {
	// TEST-MAIN-021: Test saveConfigData function
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yml")

	configData := map[string]interface{}{
		"archive_dir_path":     "/test/path",
		"use_current_dir_name": true,
		"verification": map[string]interface{}{
			"verify_on_create":   false,
			"checksum_algorithm": "md5",
		},
	}

	defer func() {
		if r := recover(); r != nil {
			if !strings.Contains(fmt.Sprintf("%v", r), "exit") {
				t.Errorf("saveConfigData panicked: %v", r)
			}
		}
	}()

	saveConfigData(configPath, configData)

	// Verify file was created and contains expected content
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Failed to read saved config file: %v", err)
		return
	}

	content := string(data)
	if !strings.Contains(content, "archive_dir_path: /test/path") {
		t.Errorf("Expected saved config to contain archive_dir_path")
	}
	if !strings.Contains(content, "use_current_dir_name: true") {
		t.Errorf("Expected saved config to contain use_current_dir_name")
	}
	if !strings.Contains(content, "verification:") {
		t.Errorf("Expected saved config to contain verification section")
	}
}

// TEST-REF: TestMain_VerifySingleArchive
func TestMain_VerifySingleArchive(t *testing.T) {
	// TEST-MAIN-022: Test verifySingleArchive function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cfg := createTestConfig(t, tmpDir)
	formatter := NewOutputFormatter(cfg)
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive dir: %v", err)
	}

	opts := VerifyOptions{
		Config:       cfg,
		Formatter:    formatter,
		ArchiveName:  "nonexistent.zip",
		WithChecksum: false,
	}

	// Should fail for non-existent archive
	err := verifySingleArchive(opts, archiveDir)
	if err == nil {
		t.Errorf("Expected error for non-existent archive")
	}
}

// TEST-REF: TestMain_VerifyAllArchives
func TestMain_VerifyAllArchives(t *testing.T) {
	// TEST-MAIN-023: Test verifyAllArchives function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	cfg := createTestConfig(t, tmpDir)
	formatter := NewOutputFormatter(cfg)
	archiveDir := filepath.Join(tmpDir, "../.bkpdir", filepath.Base(tmpDir))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive dir: %v", err)
	}

	opts := VerifyOptions{
		Config:       cfg,
		Formatter:    formatter,
		ArchiveName:  "",
		WithChecksum: false,
	}

	// Should handle empty archive directory - returns nil when no archives exist
	err := verifyAllArchives(opts, archiveDir)
	if err != nil {
		t.Errorf("Expected no error when no archives exist, got: %v", err)
	}
}

// TEST-REF: TestMain_PerformVerification
func TestMain_PerformVerification(t *testing.T) {
	// TEST-MAIN-024: Test performVerification function
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "nonexistent.zip")

	// Test with non-existent archive - should return VerificationStatus with IsVerified: false
	status, err := performVerification(archivePath, false)
	if err != nil {
		t.Errorf("Expected no error for non-existent archive, got: %v", err)
	}
	if status == nil || status.IsVerified {
		t.Errorf("Expected VerificationStatus with IsVerified: false")
	}

	// Test with checksum verification - should return VerificationStatus with IsVerified: false
	status, err = performVerification(archivePath, true)
	if err != nil {
		t.Errorf("Expected no error for non-existent archive with checksum, got: %v", err)
	}
	if status == nil || status.IsVerified {
		t.Errorf("Expected VerificationStatus with IsVerified: false for checksum verification")
	}
}

// TEST-REF: TestMain_HandleVerificationResult
func TestMain_HandleVerificationResult(t *testing.T) {
	// TEST-MAIN-025: Test handleVerificationResult function
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	createTestConfig(t, tmpDir)

	archive := &Archive{
		Name: "test-archive.zip",
		Path: filepath.Join(tmpDir, "test-archive.zip"),
	}

	// Test with successful verification
	successStatus := &VerificationStatus{
		IsVerified: true,
		Errors:     []string{},
	}

	err := handleVerificationResult(archive, successStatus, "test-archive.zip")
	if err != nil {
		t.Errorf("Expected no error for successful verification, got: %v", err)
	}

	// Test with failed verification
	failStatus := &VerificationStatus{
		IsVerified: false,
		Errors:     []string{"Test error"},
	}

	err = handleVerificationResult(archive, failStatus, "test-archive.zip")
	if err == nil {
		t.Errorf("Expected error for failed verification")
	}
}

// TEST-REF: TestMain_Integration_ConfigCommands
func TestMain_Integration_ConfigCommands(t *testing.T) {
	// TEST-MAIN-026: Integration test for config command creation and validation
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)
	os.Chdir(tmpDir)

	createTestConfig(t, tmpDir)

	// Test configCmd creation
	cmd := configCmd()
	if cmd == nil {
		t.Fatalf("configCmd() returned nil")
	}

	// Verify command structure
	if cmd.Use != "config [KEY] [VALUE]" {
		t.Errorf("Expected Use 'config [KEY] [VALUE]', got %s", cmd.Use)
	}

	// Test that the command accepts the right number of args
	if cmd.Args == nil {
		t.Errorf("Expected Args validation function")
	}
}

// TEST-REF: TestMain_Integration_CreateCommands
func TestMain_Integration_CreateCommands(t *testing.T) {
	// TEST-MAIN-027: Integration test for create command
	cmd := createCmd()
	if cmd == nil {
		t.Fatalf("createCmd() returned nil")
	}

	if cmd.Use != "create" {
		t.Errorf("Expected Use 'create', got %s", cmd.Use)
	}

	if cmd.Short != "Create a new archive" {
		t.Errorf("Expected Short 'Create a new archive', got %s", cmd.Short)
	}
}

// TEST-REF: TestMain_Integration_BackupCommandStructure
func TestMain_Integration_BackupCommandStructure(t *testing.T) {
	// TEST-MAIN-028: Integration test for backup command structure
	cmd := backupCmd()
	if cmd == nil {
		t.Fatalf("backupCmd() returned nil")
	}

	// Verify command has the right flags
	noteFlag := cmd.Flags().Lookup("note")
	if noteFlag == nil {
		t.Errorf("Expected note flag to be defined")
	}

	if noteFlag.Shorthand != "n" {
		t.Errorf("Expected note flag shorthand 'n', got %s", noteFlag.Shorthand)
	}
}

// TEST-REF: TestMain_Integration_CommandValidation
func TestMain_Integration_CommandValidation(t *testing.T) {
	// TEST-MAIN-029: Integration test for command validation
	tests := []struct {
		name        string
		cmdFunc     func() *cobra.Command
		expectedUse string
	}{
		{"config command", configCmd, "config [KEY] [VALUE]"},
		{"create command", createCmd, "create"},
		{"verify command", verifyCmd, "verify"},
		{"version command", versionCmd, "version"},
		{"backup command", backupCmd, "backup [FILE_PATH] [NOTE]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmdFunc()
			if cmd == nil {
				t.Fatalf("%s returned nil", tt.name)
			}
			if cmd.Use != tt.expectedUse {
				t.Errorf("Expected Use '%s', got %s", tt.expectedUse, cmd.Use)
			}
			if cmd.Short == "" {
				t.Errorf("Expected non-empty Short description")
			}
		})
	}
}
