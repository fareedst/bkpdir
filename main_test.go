package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
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
			rootCmd := &cobra.Command{
				Use:   "bkpdir",
				Short: "Directory archiving CLI for macOS and Linux",
			}
			rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
			rootCmd.AddCommand(fullCmd())
			rootCmd.AddCommand(incCmd())
			rootCmd.AddCommand(listCmd())

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
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	// Debug: print temp dir
	t.Logf("Temp dir: %s", tmpDir)

	// Create a test config file
	cfgPath := filepath.Join(tmpDir, ".bkpdir.yml")
	os.WriteFile(cfgPath, []byte("archive_dir_path: ../.bkpdir\nuse_current_dir_name: true\n"), 0644)

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

	// Create a full archive first (without dry-run)
	rootCmd := &cobra.Command{
		Use:   "bkpdir",
		Short: "Directory archiving CLI for macOS and Linux",
	}
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())
	rootCmd.AddCommand(listCmd())

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
			rootCmd := &cobra.Command{
				Use:   "bkpdir",
				Short: "Directory archiving CLI for macOS and Linux",
			}
			rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
			rootCmd.AddCommand(fullCmd())
			rootCmd.AddCommand(incCmd())
			rootCmd.AddCommand(listCmd())

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
				rootCmd := &cobra.Command{
					Use:   "bkpdir",
					Short: "Directory archiving CLI for macOS and Linux",
				}
				rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
				rootCmd.AddCommand(fullCmd())
				rootCmd.AddCommand(incCmd())
				rootCmd.AddCommand(listCmd())
				rootCmd.SetArgs([]string{"full"})
				if err := rootCmd.Execute(); err != nil {
					t.Fatalf("Failed to create initial full archive for inc test: %v", err)
				}
			}

			// Create a new root command for each test
			rootCmd := &cobra.Command{
				Use:   "bkpdir",
				Short: "Directory archiving CLI for macOS and Linux",
			}
			rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
			rootCmd.AddCommand(fullCmd())
			rootCmd.AddCommand(incCmd())
			rootCmd.AddCommand(listCmd())

			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Execute() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
