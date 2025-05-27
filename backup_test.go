// This file is part of bkpdir

package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateBackupName(t *testing.T) {
	t.Run("basic backup naming without notes", func(t *testing.T) {
		testBasicBackupNaming(t)
	})

	t.Run("backup naming with notes", func(t *testing.T) {
		testBackupNamingWithNotes(t)
	})

	t.Run("various source file paths and extensions", func(t *testing.T) {
		testVariousFilePaths(t)
	})

	t.Run("timestamp format validation", func(t *testing.T) {
		testTimestampFormat(t)
	})

	t.Run("special characters in filenames", func(t *testing.T) {
		testSpecialCharacters(t)
	})

	t.Run("long filenames", func(t *testing.T) {
		testLongFilenames(t)
	})

	t.Run("files with multiple extensions", func(t *testing.T) {
		testMultipleExtensions(t)
	})
}

// Helper functions for TestGenerateBackupName to reduce complexity

func testBasicBackupNaming(t *testing.T) {
	filePath := "/path/to/test.txt"
	timestamp := "20240320-153045"

	name := GenerateBackupName(filePath, timestamp, "")

	expected := "test.txt-20240320-153045"
	if name != expected {
		t.Errorf("GenerateBackupName() = %q, want %q", name, expected)
	}
}

func testBackupNamingWithNotes(t *testing.T) {
	filePath := "/path/to/test.txt"
	note := "before-changes"
	timestamp := "20240320-153045"

	name := GenerateBackupName(filePath, timestamp, note)

	expected := "test.txt-20240320-153045=before-changes"
	if name != expected {
		t.Errorf("GenerateBackupName() = %q, want %q", name, expected)
	}
}

func testVariousFilePaths(t *testing.T) {
	timestamp := "20240320-153045"

	testCases := []struct {
		filePath string
		expected string
	}{
		{"/simple.txt", "simple.txt-20240320-153045"},
		{"/path/to/file.go", "file.go-20240320-153045"},
		{"/path/to/file.tar.gz", "file.tar.gz-20240320-153045"},
		{"relative/path/file.json", "file.json-20240320-153045"},
		{"/no-extension", "no-extension-20240320-153045"},
	}

	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			name := GenerateBackupName(tc.filePath, timestamp, "")
			if name != tc.expected {
				t.Errorf("GenerateBackupName(%q) = %q, want %q", tc.filePath, name, tc.expected)
			}
		})
	}
}

func testTimestampFormat(t *testing.T) {
	filePath := "/test.txt"
	timestamp := "20241231-235959"

	name := GenerateBackupName(filePath, timestamp, "")

	expected := "test.txt-20241231-235959"
	if name != expected {
		t.Errorf("GenerateBackupName() = %q, want %q", name, expected)
	}
}

func testSpecialCharacters(t *testing.T) {
	timestamp := "20240320-153045"

	testCases := []struct {
		filePath string
		expected string
	}{
		{"/file with spaces.txt", "file with spaces.txt-20240320-153045"},
		{"/file-with-dashes.txt", "file-with-dashes.txt-20240320-153045"},
		{"/file_with_underscores.txt", "file_with_underscores.txt-20240320-153045"},
	}

	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			name := GenerateBackupName(tc.filePath, timestamp, "")
			if name != tc.expected {
				t.Errorf("GenerateBackupName(%q) = %q, want %q", tc.filePath, name, tc.expected)
			}
		})
	}
}

func testLongFilenames(t *testing.T) {
	longName := strings.Repeat("a", 200) + ".txt"
	filePath := "/" + longName
	timestamp := "20240320-153045"

	name := GenerateBackupName(filePath, timestamp, "")

	expected := longName + "-20240320-153045"
	if name != expected {
		t.Errorf("GenerateBackupName() with long filename failed")
	}
}

func testMultipleExtensions(t *testing.T) {
	timestamp := "20240320-153045"

	testCases := []struct {
		filePath string
		expected string
	}{
		{"/file.tar.gz", "file.tar.gz-20240320-153045"},
		{"/file.backup.old", "file.backup.old-20240320-153045"},
		{"/file.min.js", "file.min.js-20240320-153045"},
	}

	for _, tc := range testCases {
		t.Run(tc.filePath, func(t *testing.T) {
			name := GenerateBackupName(tc.filePath, timestamp, "")
			if name != tc.expected {
				t.Errorf("GenerateBackupName(%q) = %q, want %q", tc.filePath, name, tc.expected)
			}
		})
	}
}

func TestListFileBackups(t *testing.T) {
	t.Run("empty backup directory", func(t *testing.T) {
		testEmptyBackupDirectory(t)
	})

	t.Run("multiple backups with different timestamps", func(t *testing.T) {
		testMultipleBackupsWithTimestamps(t)
	})

	t.Run("backups with notes", func(t *testing.T) {
		testBackupsWithNotes(t)
	})

	t.Run("backups in nested directory structures", func(t *testing.T) {
		testBackupsInNestedDirectories(t)
	})

	t.Run("sorting by creation time", func(t *testing.T) {
		testBackupSortingByCreationTime(t)
	})

	t.Run("files with special characters", func(t *testing.T) {
		testBackupsWithSpecialCharacters(t)
	})

	t.Run("large number of backups", func(t *testing.T) {
		testLargeNumberOfBackups(t)
	})

	t.Run("permission denied scenarios", func(t *testing.T) {
		testBackupPermissionDenied(t)
	})
}

// Helper functions for TestListFileBackups to reduce complexity

func testEmptyBackupDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "test.txt"
	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("ListFileBackups() returned %d backups, want 0", len(backups))
	}
}

func testMultipleBackupsWithTimestamps(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "test.txt"

	// Create test backup files with different timestamps
	backupFiles := []string{
		"test.txt-20240320-100000",
		"test.txt-20240320-120000",
		"test.txt-20240320-140000",
	}

	for _, backupFile := range backupFiles {
		backupPath := filepath.Join(backupDir, backupFile)
		os.WriteFile(backupPath, []byte("test content"), 0644)
	}

	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 3 {
		t.Errorf("ListFileBackups() returned %d backups, want 3", len(backups))
	}
}

func testBackupsWithNotes(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "test.txt"

	// Create test backup files with notes
	backupFiles := []string{
		"test.txt-20240320-100000-before-changes",
		"test.txt-20240320-120000-after-fix",
	}

	for _, backupFile := range backupFiles {
		backupPath := filepath.Join(backupDir, backupFile)
		os.WriteFile(backupPath, []byte("test content"), 0644)
	}

	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 2 {
		t.Errorf("ListFileBackups() returned %d backups, want 2", len(backups))
	}
}

func testBackupsInNestedDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")

	// Create nested directory structure
	nestedDir := filepath.Join(backupDir, "subdir")
	os.MkdirAll(nestedDir, 0755)

	filePath := "test.txt"

	// Create backup file in nested directory
	backupFile := "test.txt-20240320-100000"
	backupPath := filepath.Join(nestedDir, backupFile)
	os.WriteFile(backupPath, []byte("test content"), 0644)

	// ListFileBackups only searches in the specified directory, not subdirectories
	backups, err := ListFileBackups(nestedDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("ListFileBackups() returned %d backups, want 1", len(backups))
	}
}

func testBackupSortingByCreationTime(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "test.txt"

	// Create backup files with different modification times
	backupFiles := []struct {
		name  string
		delay time.Duration
	}{
		{"test.txt-20240320-100000", 0},                     // Created first (oldest)
		{"test.txt-20240320-120000", 10 * time.Millisecond}, // Created second
		{"test.txt-20240320-140000", 20 * time.Millisecond}, // Created last (newest)
	}

	for _, backupFile := range backupFiles {
		time.Sleep(backupFile.delay)
		backupPath := filepath.Join(backupDir, backupFile.name)
		os.WriteFile(backupPath, []byte("test content"), 0644)
	}

	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 3 {
		t.Errorf("ListFileBackups() returned %d backups, want 3", len(backups))
	}

	// Verify sorting (most recent first by modification time)
	expectedOrder := []string{
		"test.txt-20240320-140000", // Created last (newest modification time)
		"test.txt-20240320-120000", // Created second
		"test.txt-20240320-100000", // Created first (oldest modification time)
	}

	for i, backup := range backups {
		if backup.Name != expectedOrder[i] {
			t.Errorf("Backup[%d].Name = %q, want %q", i, backup.Name, expectedOrder[i])
		}
	}
}

func testBackupsWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "file with spaces.txt"

	// Create backup file
	backupFile := "file with spaces.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte("test content"), 0644)

	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("ListFileBackups() returned %d backups, want 1", len(backups))
	}
}

func testLargeNumberOfBackups(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	filePath := "test.txt"

	// Create 100 backup files
	for i := 0; i < 100; i++ {
		timestamp := time.Date(2024, 3, 20, 10, i, 0, 0, time.UTC).Format("20060102-150405")
		backupFile := GenerateBackupName(filePath, timestamp, "")
		backupPath := filepath.Join(backupDir, backupFile)
		os.WriteFile(backupPath, []byte("test content"), 0644)
	}

	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 100 {
		t.Errorf("ListFileBackups() returned %d backups, want 100", len(backups))
	}
}

func testBackupPermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create a subdirectory with no read permissions
	restrictedDir := filepath.Join(backupDir, "restricted")
	os.MkdirAll(restrictedDir, 0000)
	defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup

	filePath := "test.txt"

	// This should not fail even if some subdirectories are inaccessible
	backups, err := ListFileBackups(backupDir, filePath)

	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("ListFileBackups() returned %d backups, want 0", len(backups))
	}
}

func TestCreateFileBackup(t *testing.T) {
	t.Run("successful backup creation", func(t *testing.T) {
		testSuccessfulBackupCreation(t)
	})

	t.Run("file identical to existing backup", func(t *testing.T) {
		testIdenticalFileBackup(t)
	})

	t.Run("file not found", func(t *testing.T) {
		testFileNotFound(t)
	})

	t.Run("source path is directory", func(t *testing.T) {
		testSourcePathIsDirectory(t)
	})

	t.Run("backup directory creation", func(t *testing.T) {
		testBackupDirectoryCreation(t)
	})
}

// Helper functions for TestCreateFileBackup to reduce complexity

func testSuccessfulBackupCreation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "test content for backup"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create config
	cfg := DefaultConfig()
	cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
	cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

	// Create backup
	err := CreateFileBackup(cfg, sourceFile, "test-note", false)

	if err != nil {
		t.Fatalf("CreateFileBackup() error = %v", err)
	}

	// Verify backup was created
	backups, err := ListFileBackups(cfg.BackupDirPath, filepath.Base(sourceFile))
	if err != nil {
		t.Fatalf("ListFileBackups() error = %v", err)
	}
	if len(backups) != 1 {
		t.Errorf("Expected 1 backup, got %d", len(backups))
	}

	// Verify backup content
	backupPath := filepath.Join(cfg.BackupDirPath, backups[0].Name)
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	if string(backupContent) != sourceContent {
		t.Errorf("Backup content = %q, want %q", string(backupContent), sourceContent)
	}
}

func testIdenticalFileBackup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "identical content"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create config
	cfg := DefaultConfig()
	cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
	cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

	// Create first backup
	err := CreateFileBackup(cfg, sourceFile, "", false)
	if err != nil {
		t.Fatalf("First CreateFileBackup() error = %v", err)
	}

	// Verify first backup was created
	backups, listErr := ListFileBackups(cfg.BackupDirPath, filepath.Base(sourceFile))
	if listErr != nil {
		t.Fatalf("ListFileBackups() error = %v", listErr)
	}
	if len(backups) != 1 {
		t.Errorf("Expected 1 backup after first creation, got %d", len(backups))
	}

	// Note: The second backup attempt would call os.Exit() when it detects identical content
	// This is the expected behavior of the application, so we don't test it here
	// as it would terminate the test process
}

func testFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config
	cfg := DefaultConfig()
	cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
	cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

	// Try to backup non-existent file
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
	err := CreateFileBackup(cfg, nonExistentFile, "", false)

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func testSourcePathIsDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config
	cfg := DefaultConfig()
	cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
	cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

	// Try to backup a directory
	sourceDir := filepath.Join(tmpDir, "sourcedir")
	os.MkdirAll(sourceDir, 0755)

	err := CreateFileBackup(cfg, sourceDir, "", false)

	if err == nil {
		t.Error("Expected error when trying to backup a directory")
	}
}

func testBackupDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	os.WriteFile(sourceFile, []byte("test content"), 0644)

	// Create config with non-existent backup directory
	cfg := DefaultConfig()
	cfg.BackupDirPath = filepath.Join(tmpDir, "nonexistent", "backups")
	cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

	// Create backup (should create directory structure)
	err := CreateFileBackup(cfg, sourceFile, "", false)

	if err != nil {
		t.Fatalf("CreateFileBackup() error = %v", err)
	}

	// Verify backup directory was created
	if _, err := os.Stat(cfg.BackupDirPath); os.IsNotExist(err) {
		t.Error("Backup directory was not created")
	}
}

func TestCreateFileBackupWithContext(t *testing.T) {
	t.Run("successful backup with context", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source file
		sourceFile := filepath.Join(tmpDir, "source.txt")
		os.WriteFile(sourceFile, []byte("test content"), 0644)

		// Create config
		cfg := DefaultConfig()
		cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
		cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

		// Create backup with context
		ctx := context.Background()
		err := CreateFileBackupWithContext(ctx, cfg, sourceFile, "", false)

		if err != nil {
			t.Fatalf("CreateFileBackupWithContext() error = %v", err)
		}

		// Verify backup was created
		backups, err := ListFileBackups(cfg.BackupDirPath, filepath.Base(sourceFile))
		if err != nil {
			t.Fatalf("ListFileBackups() error = %v", err)
		}
		if len(backups) != 1 {
			t.Errorf("Expected 1 backup, got %d", len(backups))
		}
	})

	t.Run("backup cancelled via context", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source file
		sourceFile := filepath.Join(tmpDir, "source.txt")
		os.WriteFile(sourceFile, []byte("test content"), 0644)

		// Create config
		cfg := DefaultConfig()
		cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
		cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

		// Create cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := CreateFileBackupWithContext(ctx, cfg, sourceFile, "", false)

		if err == nil {
			t.Error("Expected error for cancelled context")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("backup with timeout", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source file
		sourceFile := filepath.Join(tmpDir, "source.txt")
		os.WriteFile(sourceFile, []byte("test content"), 0644)

		// Create config
		cfg := DefaultConfig()
		cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
		cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

		// Create context with very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		// Wait for timeout
		time.Sleep(1 * time.Millisecond)

		err := CreateFileBackupWithContext(ctx, cfg, sourceFile, "", false)

		if err == nil {
			t.Error("Expected timeout error")
		}
		if err != context.DeadlineExceeded {
			t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
		}
	})

	t.Run("already cancelled context handling", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create source file
		sourceFile := filepath.Join(tmpDir, "source.txt")
		os.WriteFile(sourceFile, []byte("test content"), 0644)

		// Create config
		cfg := DefaultConfig()
		cfg.BackupDirPath = filepath.Join(tmpDir, "backups")
		cfg.UseCurrentDirNameForFiles = false // Disable complex directory structure

		// Create already cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := CreateFileBackupWithContext(ctx, cfg, sourceFile, "", false)

		if err == nil {
			t.Error("Expected error for already cancelled context")
		}
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
	})
}

func TestCheckForIdenticalFileBackup(t *testing.T) {
	t.Run("file identical to most recent backup", func(t *testing.T) {
		testIdenticalToRecentBackup(t)
	})

	t.Run("file different from existing backups", func(t *testing.T) {
		testDifferentFromExistingBackups(t)
	})

	t.Run("no existing backups", func(t *testing.T) {
		testNoExistingBackups(t)
	})

	t.Run("multiple backups with different content", func(t *testing.T) {
		testMultipleBackupsWithDifferentContent(t)
	})

	t.Run("large file comparison", func(t *testing.T) {
		testLargeFileComparison(t)
	})

	t.Run("binary file comparison", func(t *testing.T) {
		testBinaryFileComparison(t)
	})

	t.Run("files with special characters", func(t *testing.T) {
		testFilesWithSpecialCharacters(t)
	})

	t.Run("permission denied during comparison", func(t *testing.T) {
		testPermissionDeniedDuringComparison(t)
	})
}

// Helper functions for TestCheckForIdenticalFileBackup to reduce complexity

func testIdenticalToRecentBackup(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "identical content"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create backup file
	backupFile := "source.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte(sourceContent), 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if !isIdentical {
		t.Error("Expected files to be identical")
	}
}

func testDifferentFromExistingBackups(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "new content"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create backup file with different content
	backupFile := "source.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte("old content"), 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if isIdentical {
		t.Error("Expected files to be different")
	}
}

func testNoExistingBackups(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	os.WriteFile(sourceFile, []byte("content"), 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if isIdentical {
		t.Error("Expected no identical backup when no backups exist")
	}
}

func testMultipleBackupsWithDifferentContent(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	sourceContent := "current content"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create multiple backup files with different content
	backupFiles := []struct {
		name    string
		content string
	}{
		{"source.txt-20240320-100000", "old content 1"},
		{"source.txt-20240320-110000", "old content 2"},
		{"source.txt-20240320-120000", sourceContent}, // This one matches
	}

	for _, backup := range backupFiles {
		backupPath := filepath.Join(backupDir, backup.name)
		os.WriteFile(backupPath, []byte(backup.content), 0644)
	}

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if !isIdentical {
		t.Error("Expected to find identical backup among multiple backups")
	}
}

func testLargeFileComparison(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create large source file (1MB)
	sourceFile := filepath.Join(tmpDir, "large.txt")
	largeContent := strings.Repeat("a", 1024*1024)
	os.WriteFile(sourceFile, []byte(largeContent), 0644)

	// Create identical backup file
	backupFile := "large.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte(largeContent), 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if !isIdentical {
		t.Error("Expected large files to be identical")
	}
}

func testBinaryFileComparison(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create binary source file
	sourceFile := filepath.Join(tmpDir, "binary.dat")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	os.WriteFile(sourceFile, binaryContent, 0644)

	// Create identical backup file
	backupFile := "binary.dat-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, binaryContent, 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if !isIdentical {
		t.Error("Expected binary files to be identical")
	}
}

func testFilesWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file with special characters in name
	sourceFile := filepath.Join(tmpDir, "file with spaces.txt")
	sourceContent := "content with special chars: éñ中文"
	os.WriteFile(sourceFile, []byte(sourceContent), 0644)

	// Create backup file
	backupFile := "file with spaces.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte(sourceContent), 0644)

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	if err != nil {
		t.Fatalf("CheckForIdenticalFileBackup() error = %v", err)
	}
	if !isIdentical {
		t.Error("Expected files with special characters to be identical")
	}
}

func testPermissionDeniedDuringComparison(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()
	backupDir := filepath.Join(tmpDir, "backups")
	os.MkdirAll(backupDir, 0755)

	// Create source file
	sourceFile := filepath.Join(tmpDir, "source.txt")
	os.WriteFile(sourceFile, []byte("content"), 0644)

	// Create backup file with no read permissions
	backupFile := "source.txt-20240320-100000"
	backupPath := filepath.Join(backupDir, backupFile)
	os.WriteFile(backupPath, []byte("content"), 0000)
	defer os.Chmod(backupPath, 0644) // Restore permissions for cleanup

	isIdentical, _, err := CheckForIdenticalFileBackup(sourceFile, backupDir, filepath.Base(sourceFile))

	// Should handle permission errors gracefully
	if err != nil {
		t.Logf("Permission error handled: %v", err)
	}
	if isIdentical {
		t.Error("Should not report identical when backup file is not readable")
	}
}
