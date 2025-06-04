// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// DefaultFileSystemTestHelper provides standard file system testing functionality.
type DefaultFileSystemTestHelper struct{}

// NewFileSystemTestHelper creates a new file system test helper.
func NewFileSystemTestHelper() FileSystemTestHelper {
	return &DefaultFileSystemTestHelper{}
}

// CreateTempDir creates a temporary directory and returns its path.
// This enhances the standard t.TempDir() with custom prefix support.
//
// ⭐ EXTRACT-009: Temporary directory creation - 🔧
func (h *DefaultFileSystemTestHelper) CreateTempDir(t *testing.T, prefix string) string {
	t.Helper()

	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory with prefix %q: %v", prefix, err)
	}

	// Register cleanup with the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

// CreateTempFile creates a temporary file with content and returns its path.
//
// ⭐ EXTRACT-009: Temporary file creation - 🔧
func (h *DefaultFileSystemTestHelper) CreateTempFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()

	filePath := filepath.Join(dir, name)

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatalf("Failed to create parent directories for %q: %v", filePath, err)
	}

	// Create the file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("Failed to create temp file %q: %v", filePath, err)
	}

	return filePath
}

// CreateDirectory creates a directory structure from a map of paths to content.
// Extracted from comparison_test.go createTestDirectory function.
//
// ⭐ EXTRACT-009: Directory structure creation - 🔧
func (h *DefaultFileSystemTestHelper) CreateDirectory(t *testing.T, root string, files map[string]string) error {
	t.Helper()

	for filePath, content := range files {
		fullPath := filepath.Join(root, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directories for %q: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %q: %w", fullPath, err)
		}
	}
	return nil
}

// CreateZipArchive creates a ZIP archive from a map of paths to content.
// Extracted from comparison_test.go createTestZipArchive function.
//
// ⭐ EXTRACT-009: ZIP archive creation - 🔧
func (h *DefaultFileSystemTestHelper) CreateZipArchive(t *testing.T, archivePath string, files map[string]string) error {
	t.Helper()

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(archivePath), 0755); err != nil {
		return fmt.Errorf("failed to create directories for archive %q: %w", archivePath, err)
	}

	zipFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file %q: %w", archivePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for filePath, content := range files {
		fileWriter, err := zipWriter.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %q in ZIP: %w", filePath, err)
		}
		_, err = fileWriter.Write([]byte(content))
		if err != nil {
			return fmt.Errorf("failed to write content to file %q in ZIP: %w", filePath, err)
		}
	}

	return nil
}

// CreateTestFiles creates multiple test files in a directory from a map.
// This is a convenience function that combines directory and file creation.
//
// ⭐ EXTRACT-009: Multiple test file creation - 🔧
func (h *DefaultFileSystemTestHelper) CreateTestFiles(t *testing.T, baseDir string, files map[string]string) {
	t.Helper()

	for filePath, content := range files {
		fullPath := filepath.Join(baseDir, filePath)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directories for %q: %v", fullPath, err)
		}

		// Create the file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %q: %v", fullPath, err)
		}
	}
}

// Package-level convenience functions

var defaultFileSystemHelper = NewFileSystemTestHelper()

// CreateTempDir is a package-level convenience function for creating temporary directories.
//
// ⭐ EXTRACT-009: Package-level temp directory - 🔧
func CreateTempDir(t *testing.T, prefix string) string {
	return defaultFileSystemHelper.CreateTempDir(t, prefix)
}

// CreateTempFile is a package-level convenience function for creating temporary files.
//
// ⭐ EXTRACT-009: Package-level temp file - 🔧
func CreateTempFile(t *testing.T, dir, name string, content []byte) string {
	return defaultFileSystemHelper.CreateTempFile(t, dir, name, content)
}

// CreateDirectory is a package-level convenience function for creating directory structures.
//
// ⭐ EXTRACT-009: Package-level directory creation - 🔧
func CreateDirectory(t *testing.T, root string, files map[string]string) error {
	return defaultFileSystemHelper.CreateDirectory(t, root, files)
}

// CreateZipArchive is a package-level convenience function for creating ZIP archives.
//
// ⭐ EXTRACT-009: Package-level ZIP creation - 🔧
func CreateZipArchive(t *testing.T, archivePath string, files map[string]string) error {
	return defaultFileSystemHelper.CreateZipArchive(t, archivePath, files)
}

// CreateTestFiles is a package-level convenience function for creating multiple test files.
//
// ⭐ EXTRACT-009: Package-level test files - 🔧
func CreateTestFiles(t *testing.T, baseDir string, files map[string]string) {
	defaultFileSystemHelper.CreateTestFiles(t, baseDir, files)
}

// WithTempDir executes a function with a temporary directory and cleans up automatically.
// This is a convenience function for tests that need a temporary directory.
//
// ⭐ EXTRACT-009: Temporary directory with cleanup - 🔧
func WithTempDir(t *testing.T, prefix string, fn func(dir string)) {
	t.Helper()

	tempDir := CreateTempDir(t, prefix)
	fn(tempDir)
	// Cleanup is automatically handled by t.Cleanup in CreateTempDir
}

// WithTestFiles executes a function with test files created in a temporary directory.
// This combines temporary directory creation with file setup.
//
// ⭐ EXTRACT-009: Test files with cleanup - 🔧
func WithTestFiles(t *testing.T, files map[string]string, fn func(dir string)) {
	t.Helper()

	WithTempDir(t, "testutil-files-", func(dir string) {
		CreateTestFiles(t, dir, files)
		fn(dir)
	})
}

// WithTestArchive executes a function with a test ZIP archive and cleans up automatically.
//
// ⭐ EXTRACT-009: Test archive with cleanup - 🔧
func WithTestArchive(t *testing.T, files map[string]string, fn func(archivePath string)) {
	t.Helper()

	WithTempDir(t, "testutil-archive-", func(dir string) {
		archivePath := filepath.Join(dir, "test.zip")

		if err := CreateZipArchive(t, archivePath, files); err != nil {
			t.Fatalf("Failed to create test archive: %v", err)
		}

		fn(archivePath)
	})
}
