// This file is part of bkpdir
// 🔺 LINT-001: Lint compliance - 🔧

package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Helper function to create a test ZIP archive
func createTestZipArchive(archivePath string, files map[string]string) error {
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for filePath, content := range files {
		fileWriter, err := zipWriter.Create(filePath)
		if err != nil {
			return err
		}
		_, err = fileWriter.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to create a test directory structure
func createTestDirectory(rootPath string, files map[string]string) error {
	for filePath, content := range files {
		fullPath := filepath.Join(rootPath, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// TestCreateArchiveSnapshot tests the CreateArchiveSnapshot function - 0% coverage
func TestCreateArchiveSnapshot(t *testing.T) {
	tests := []struct {
		name      string
		files     map[string]string
		expectErr bool
	}{
		{
			name: "successful archive snapshot",
			files: map[string]string{
				"file1.txt":        "content1",
				"file2.txt":        "content2",
				"subdir/file3.txt": "content3",
			},
			expectErr: false,
		},
		{
			name:      "empty archive",
			files:     map[string]string{},
			expectErr: false,
		},
		{
			name: "archive with single file",
			files: map[string]string{
				"single.txt": "single content",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary archive
			tempDir := t.TempDir()
			archivePath := filepath.Join(tempDir, "test.zip")

			// Create test ZIP archive
			err := createTestZipArchive(archivePath, tt.files)
			if err != nil {
				t.Fatalf("Failed to create test archive: %v", err)
			}

			// Test CreateArchiveSnapshot
			snapshot, err := CreateArchiveSnapshot(archivePath)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if snapshot == nil {
				t.Errorf("Expected snapshot, got nil")
				return
			}

			// Verify file count
			expectedCount := len(tt.files)
			if len(snapshot.Files) != expectedCount {
				t.Errorf("Expected %d files, got %d", expectedCount, len(snapshot.Files))
			}

			// Verify files are sorted
			for i := 1; i < len(snapshot.Files); i++ {
				if snapshot.Files[i-1].RelativePath >= snapshot.Files[i].RelativePath {
					t.Errorf("Files not sorted: %s >= %s",
						snapshot.Files[i-1].RelativePath,
						snapshot.Files[i].RelativePath)
				}
			}

			// Verify file properties
			for _, file := range snapshot.Files {
				if file.IsDir {
					t.Errorf("Expected no directories in archive snapshot, found: %s", file.RelativePath)
				}
				if file.Hash == "" {
					t.Errorf("Expected hash for file %s, got empty", file.RelativePath)
				}
				if file.Size <= 0 && tt.files[file.RelativePath] != "" {
					t.Errorf("Expected positive size for file %s, got %d", file.RelativePath, file.Size)
				}
			}
		})
	}
}

// TestCreateArchiveSnapshot_ErrorCases tests error scenarios
func TestCreateArchiveSnapshot_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		archivePath string
		expectErr   bool
	}{
		{
			name:        "nonexistent archive",
			archivePath: "/nonexistent/archive.zip",
			expectErr:   true,
		},
		{
			name:        "invalid archive path",
			archivePath: "",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateArchiveSnapshot(tt.archivePath)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestCalculateArchiveFileHash tests the calculateArchiveFileHash function - 0% coverage
func TestCalculateArchiveFileHash(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive with known content
	files := map[string]string{
		"test.txt": "test content",
	}

	err := createTestZipArchive(archivePath, files)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Open the archive and get the file
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("Failed to open archive: %v", err)
	}
	defer reader.Close()

	if len(reader.File) == 0 {
		t.Fatalf("No files in archive")
	}

	file := reader.File[0]
	hash, err := calculateArchiveFileHash(file)
	if err != nil {
		t.Errorf("calculateArchiveFileHash failed: %v", err)
	}

	if hash == "" {
		t.Errorf("Expected non-empty hash")
	}

	// Verify hash length (SHA-256 produces 64-character hex string)
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}

	// Test consistency - same content should produce same hash
	hash2, err := calculateArchiveFileHash(file)
	if err != nil {
		t.Errorf("Second calculateArchiveFileHash failed: %v", err)
	}

	if hash != hash2 {
		t.Errorf("Hash inconsistency: %s != %s", hash, hash2)
	}
}

// TestIsDirectoryIdenticalToArchive tests the IsDirectoryIdenticalToArchive function - 0% coverage
func TestIsDirectoryIdenticalToArchive(t *testing.T) {
	tests := []struct {
		name            string
		dirFiles        map[string]string
		archiveFiles    map[string]string
		excludePattern  []string
		expectIdentical bool
		expectErr       bool
	}{
		{
			name: "identical directory and archive",
			dirFiles: map[string]string{
				"file1.txt": "content1",
				"file2.txt": "content2",
			},
			archiveFiles: map[string]string{
				"file1.txt": "content1",
				"file2.txt": "content2",
			},
			expectIdentical: true,
			expectErr:       false,
		},
		{
			name: "different content",
			dirFiles: map[string]string{
				"file1.txt": "content1",
			},
			archiveFiles: map[string]string{
				"file1.txt": "different content",
			},
			expectIdentical: false,
			expectErr:       false,
		},
		{
			name: "different files",
			dirFiles: map[string]string{
				"file1.txt": "content1",
			},
			archiveFiles: map[string]string{
				"file2.txt": "content1",
			},
			expectIdentical: false,
			expectErr:       false,
		},
		{
			name: "with exclusion patterns",
			dirFiles: map[string]string{
				"file1.txt":  "content1",
				"ignore.tmp": "ignore this",
			},
			archiveFiles: map[string]string{
				"file1.txt": "content1",
			},
			excludePattern:  []string{"*.tmp"},
			expectIdentical: true,
			expectErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			// Create test directory
			dirPath := filepath.Join(tempDir, "testdir")
			err := os.Mkdir(dirPath, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			err = createTestDirectory(dirPath, tt.dirFiles)
			if err != nil {
				t.Fatalf("Failed to create test directory structure: %v", err)
			}

			// Create test archive
			archivePath := filepath.Join(tempDir, "test.zip")
			err = createTestZipArchive(archivePath, tt.archiveFiles)
			if err != nil {
				t.Fatalf("Failed to create test archive: %v", err)
			}

			// Test the function
			identical, err := IsDirectoryIdenticalToArchive(dirPath, archivePath, tt.excludePattern)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if identical != tt.expectIdentical {
				t.Errorf("Expected identical=%v, got identical=%v", tt.expectIdentical, identical)
			}
		})
	}
}

// TestIsDirectoryIdenticalToArchive_ErrorCases tests error scenarios
func TestIsDirectoryIdenticalToArchive_ErrorCases(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		dirPath     string
		archivePath string
		expectErr   bool
	}{
		{
			name:        "nonexistent directory",
			dirPath:     "/nonexistent/dir",
			archivePath: filepath.Join(tempDir, "test.zip"),
			expectErr:   true,
		},
		{
			name:        "nonexistent archive",
			dirPath:     tempDir,
			archivePath: "/nonexistent/archive.zip",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a valid archive for the directory test case
			if tt.name == "nonexistent directory" {
				err := createTestZipArchive(tt.archivePath, map[string]string{"test.txt": "content"})
				if err != nil {
					t.Fatalf("Failed to create test archive: %v", err)
				}
			}

			_, err := IsDirectoryIdenticalToArchive(tt.dirPath, tt.archivePath, []string{})

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestFindMostRecentArchive tests the FindMostRecentArchive function - 0% coverage
func TestFindMostRecentArchive(t *testing.T) {
	tempDir := t.TempDir()
	archiveDir := filepath.Join(tempDir, "archives")
	err := os.MkdirAll(archiveDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	// Test empty directory
	t.Run("empty directory", func(t *testing.T) {
		mostRecent, err := FindMostRecentArchive(archiveDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if mostRecent != "" {
			t.Errorf("Expected empty string, got %s", mostRecent)
		}
	})

	// Create test archives with different timestamps
	archives := []struct {
		name    string
		content map[string]string
		delay   time.Duration
	}{
		{
			name:    "archive1.zip",
			content: map[string]string{"file1.txt": "content1"},
			delay:   0,
		},
		{
			name:    "archive2.zip",
			content: map[string]string{"file2.txt": "content2"},
			delay:   100 * time.Millisecond,
		},
		{
			name:    "incremental_update=something.zip", // Should be ignored
			content: map[string]string{"file3.txt": "content3"},
			delay:   200 * time.Millisecond,
		},
		{
			name:    "archive3.zip",
			content: map[string]string{"file4.txt": "content4"},
			delay:   300 * time.Millisecond,
		},
	}

	// Create archives with proper timing
	for _, archive := range archives {
		time.Sleep(archive.delay)
		archivePath := filepath.Join(archiveDir, archive.name)
		err := createTestZipArchive(archivePath, archive.content)
		if err != nil {
			t.Fatalf("Failed to create archive %s: %v", archive.name, err)
		}
	}

	// Test finding most recent archive
	t.Run("find most recent", func(t *testing.T) {
		mostRecent, err := FindMostRecentArchive(archiveDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		expectedName := "archive3.zip"
		if !strings.HasSuffix(mostRecent, expectedName) {
			t.Errorf("Expected most recent archive to be %s, got %s", expectedName, mostRecent)
		}
	})

	// Test directory with only incremental archives
	tempDir2 := t.TempDir()
	incrementalDir := filepath.Join(tempDir2, "incremental")
	err = os.MkdirAll(incrementalDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create incremental directory: %v", err)
	}

	incrementalPath := filepath.Join(incrementalDir, "inc_update=test.zip")
	err = createTestZipArchive(incrementalPath, map[string]string{"inc.txt": "content"})
	if err != nil {
		t.Fatalf("Failed to create incremental archive: %v", err)
	}

	t.Run("only incremental archives", func(t *testing.T) {
		mostRecent, err := FindMostRecentArchive(incrementalDir)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if mostRecent != "" {
			t.Errorf("Expected empty string for incremental-only directory, got %s", mostRecent)
		}
	})
}

// TestFindMostRecentArchive_ErrorCases tests error scenarios
func TestFindMostRecentArchive_ErrorCases(t *testing.T) {
	tests := []struct {
		name       string
		archiveDir string
		expectErr  bool
	}{
		{
			name:       "nonexistent directory",
			archiveDir: "/nonexistent/path",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FindMostRecentArchive(tt.archiveDir)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestCheckForIdenticalArchive tests the CheckForIdenticalArchive function - 0% coverage
func TestCheckForIdenticalArchive(t *testing.T) {
	tempDir := t.TempDir()

	// Test with no archives
	t.Run("no archives", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir1")
		archiveDir := filepath.Join(tempDir, "archives1")

		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err = os.MkdirAll(archiveDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create archive directory: %v", err)
		}

		identical, archivePath, err := CheckForIdenticalArchive(dirPath, archiveDir, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if identical {
			t.Errorf("Expected not identical with no archives")
		}
		if archivePath != "" {
			t.Errorf("Expected empty archive path, got %s", archivePath)
		}
	})

	// Test with identical directory and archive
	t.Run("identical directory and archive", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir2")
		archiveDir := filepath.Join(tempDir, "archives2")

		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err = os.MkdirAll(archiveDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create archive directory: %v", err)
		}

		// Create test files
		files := map[string]string{
			"file1.txt": "content1",
			"file2.txt": "content2",
		}

		err = createTestDirectory(dirPath, files)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create identical archive
		archivePath := filepath.Join(archiveDir, "test.zip")
		err = createTestZipArchive(archivePath, files)
		if err != nil {
			t.Fatalf("Failed to create test archive: %v", err)
		}

		identical, returnedPath, err := CheckForIdenticalArchive(dirPath, archiveDir, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !identical {
			t.Errorf("Expected identical directory and archive")
		}
		if returnedPath != archivePath {
			t.Errorf("Expected archive path %s, got %s", archivePath, returnedPath)
		}
	})

	// Test with different directory and archive
	t.Run("different directory and archive", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "testdir3")
		archiveDir := filepath.Join(tempDir, "archives3")

		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err = os.MkdirAll(archiveDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create archive directory: %v", err)
		}

		// Create different files
		dirFiles := map[string]string{
			"file1.txt": "content1",
		}

		archiveFiles := map[string]string{
			"file2.txt": "content2",
		}

		err = createTestDirectory(dirPath, dirFiles)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		archivePath := filepath.Join(archiveDir, "test.zip")
		err = createTestZipArchive(archivePath, archiveFiles)
		if err != nil {
			t.Fatalf("Failed to create test archive: %v", err)
		}

		identical, returnedPath, err := CheckForIdenticalArchive(dirPath, archiveDir, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if identical {
			t.Errorf("Expected different directory and archive")
		}
		if returnedPath != archivePath {
			t.Errorf("Expected archive path %s, got %s", archivePath, returnedPath)
		}
	})
}

// TestGetDirectoryTreeSummary tests the GetDirectoryTreeSummary function - 0% coverage
func TestGetDirectoryTreeSummary(t *testing.T) {
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Test empty directory
	t.Run("empty directory", func(t *testing.T) {
		summary, err := GetDirectoryTreeSummary(testDir, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !strings.Contains(summary, testDir) {
			t.Errorf("Expected summary to contain directory path")
		}
		if !strings.Contains(summary, "Files: 0") {
			t.Errorf("Expected summary to show 0 files")
		}
	})

	// Create test files
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	err = createTestDirectory(testDir, files)
	if err != nil {
		t.Fatalf("Failed to create test directory structure: %v", err)
	}

	t.Run("directory with files", func(t *testing.T) {
		summary, err := GetDirectoryTreeSummary(testDir, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify summary contains expected information
		if !strings.Contains(summary, testDir) {
			t.Errorf("Expected summary to contain directory path")
		}

		// Should have subdir + 3 files = 4 entries
		if !strings.Contains(summary, "Files: 4") {
			t.Errorf("Expected summary to show 4 files, got: %s", summary)
		}

		// Verify file entries
		if !strings.Contains(summary, "[FILE] file1.txt") {
			t.Errorf("Expected summary to contain file1.txt")
		}
		if !strings.Contains(summary, "[DIR]  subdir") {
			t.Errorf("Expected summary to contain subdir")
		}
		if !strings.Contains(summary, "[FILE] subdir/file3.txt") {
			t.Errorf("Expected summary to contain subdir/file3.txt")
		}
	})

	// Test with exclusion patterns
	t.Run("with exclusion patterns", func(t *testing.T) {
		// Add a file to exclude
		excludeFile := filepath.Join(testDir, "exclude.tmp")
		err := os.WriteFile(excludeFile, []byte("exclude this"), 0644)
		if err != nil {
			t.Fatalf("Failed to create exclude file: %v", err)
		}

		summary, err := GetDirectoryTreeSummary(testDir, []string{"*.tmp"})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Should still show 4 files (subdir + 3 original files, excluding the .tmp)
		if !strings.Contains(summary, "Files: 4") {
			t.Errorf("Expected summary to show 4 files after exclusion, got: %s", summary)
		}

		// Should not contain the excluded file
		if strings.Contains(summary, "exclude.tmp") {
			t.Errorf("Expected summary to exclude .tmp file")
		}
	})
}

// TestGetDirectoryTreeSummary_ErrorCases tests error scenarios
func TestGetDirectoryTreeSummary_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		dirPath   string
		expectErr bool
	}{
		{
			name:      "nonexistent directory",
			dirPath:   "/nonexistent/directory",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDirectoryTreeSummary(tt.dirPath, []string{})

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestGetArchiveTreeSummary tests the GetArchiveTreeSummary function - 0% coverage
func TestGetArchiveTreeSummary(t *testing.T) {
	tempDir := t.TempDir()

	// Test empty archive
	t.Run("empty archive", func(t *testing.T) {
		archivePath := filepath.Join(tempDir, "empty.zip")
		err := createTestZipArchive(archivePath, map[string]string{})
		if err != nil {
			t.Fatalf("Failed to create empty archive: %v", err)
		}

		summary, err := GetArchiveTreeSummary(archivePath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !strings.Contains(summary, archivePath) {
			t.Errorf("Expected summary to contain archive path")
		}
		if !strings.Contains(summary, "Files: 0") {
			t.Errorf("Expected summary to show 0 files")
		}
	})

	// Test archive with files
	t.Run("archive with files", func(t *testing.T) {
		files := map[string]string{
			"file1.txt":        "content1",
			"file2.txt":        "content2",
			"subdir/file3.txt": "content3",
		}

		archivePath := filepath.Join(tempDir, "test.zip")
		err := createTestZipArchive(archivePath, files)
		if err != nil {
			t.Fatalf("Failed to create test archive: %v", err)
		}

		summary, err := GetArchiveTreeSummary(archivePath)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify summary contains expected information
		if !strings.Contains(summary, archivePath) {
			t.Errorf("Expected summary to contain archive path")
		}
		if !strings.Contains(summary, "Files: 3") {
			t.Errorf("Expected summary to show 3 files, got: %s", summary)
		}

		// Verify file entries (archives don't show directories separately)
		if !strings.Contains(summary, "[FILE] file1.txt") {
			t.Errorf("Expected summary to contain file1.txt")
		}
		if !strings.Contains(summary, "[FILE] subdir/file3.txt") {
			t.Errorf("Expected summary to contain subdir/file3.txt")
		}

		// Should show hash snippets
		lines := strings.Split(summary, "\n")
		fileLineFound := false
		for _, line := range lines {
			if strings.Contains(line, "[FILE]") && strings.Contains(line, "bytes") {
				fileLineFound = true
				// Should contain hash snippet (first 8 characters)
				if !strings.Contains(line, "(") || !strings.Contains(line, ")") {
					t.Errorf("Expected file line to contain hash snippet: %s", line)
				}
				break
			}
		}
		if !fileLineFound {
			t.Errorf("Expected to find at least one file line with hash")
		}
	})
}

// TestGetArchiveTreeSummary_ErrorCases tests error scenarios
func TestGetArchiveTreeSummary_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		archivePath string
		expectErr   bool
	}{
		{
			name:        "nonexistent archive",
			archivePath: "/nonexistent/archive.zip",
			expectErr:   true,
		},
		{
			name:        "invalid archive path",
			archivePath: "",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetArchiveTreeSummary(tt.archivePath)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Benchmark tests for performance baselines
func BenchmarkCreateArchiveSnapshot(b *testing.B) {
	tempDir := b.TempDir()

	// Create a test archive
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	archivePath := filepath.Join(tempDir, "benchmark.zip")
	err := createTestZipArchive(archivePath, files)
	if err != nil {
		b.Fatalf("Failed to create test archive: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CreateArchiveSnapshot(archivePath)
		if err != nil {
			b.Fatalf("CreateArchiveSnapshot failed: %v", err)
		}
	}
}

func BenchmarkGetDirectoryTreeSummary(b *testing.B) {
	tempDir := b.TempDir()
	testDir := filepath.Join(tempDir, "benchmark")

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		b.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test files
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	err = createTestDirectory(testDir, files)
	if err != nil {
		b.Fatalf("Failed to create test directory structure: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetDirectoryTreeSummary(testDir, []string{})
		if err != nil {
			b.Fatalf("GetDirectoryTreeSummary failed: %v", err)
		}
	}
}

// TestCompareFiles tests the file comparison functionality for FILE-003 feature
func TestCompareFiles(t *testing.T) {
	// ⭐ FILE-003: File comparison validation - 🔍
	// TEST-REF: Feature tracking matrix FILE-003
	// IMMUTABLE-REF: Identical File Detection

	t.Run("IdenticalSnapshots", func(t *testing.T) {
		// Create temporary directories for testing
		tmpDir1 := t.TempDir()
		tmpDir2 := t.TempDir()

		// Create identical files in both directories
		content := "test content"
		file1 := filepath.Join(tmpDir1, "test.txt")
		file2 := filepath.Join(tmpDir2, "test.txt")

		if err := os.WriteFile(file1, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file2, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		snapshot1, err := CreateDirectorySnapshot(tmpDir1, []string{})
		if err != nil {
			t.Fatal(err)
		}

		snapshot2, err := CreateDirectorySnapshot(tmpDir2, []string{})
		if err != nil {
			t.Fatal(err)
		}

		if !CompareSnapshots(snapshot1, snapshot2) {
			t.Error("Expected identical snapshots to compare as equal")
		}
	})

	t.Run("DifferentSnapshots", func(t *testing.T) {
		// Create temporary directories for testing
		tmpDir1 := t.TempDir()
		tmpDir2 := t.TempDir()

		// Create different files
		file1 := filepath.Join(tmpDir1, "different1.txt")
		file2 := filepath.Join(tmpDir2, "different2.txt")

		if err := os.WriteFile(file1, []byte("content1"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(file2, []byte("content2"), 0644); err != nil {
			t.Fatal(err)
		}

		snapshot1, err := CreateDirectorySnapshot(tmpDir1, []string{})
		if err != nil {
			t.Fatal(err)
		}

		snapshot2, err := CreateDirectorySnapshot(tmpDir2, []string{})
		if err != nil {
			t.Fatal(err)
		}

		if CompareSnapshots(snapshot1, snapshot2) {
			t.Error("Expected different snapshots to compare as not equal")
		}
	})

	t.Run("WithExclusionPatterns", func(t *testing.T) {
		// Create temporary directory for testing
		tmpDir := t.TempDir()

		// Create files, some to be excluded
		testFile := filepath.Join(tmpDir, "include.txt")
		excludeFile := filepath.Join(tmpDir, "exclude.tmp")

		if err := os.WriteFile(testFile, []byte("include"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(excludeFile, []byte("exclude"), 0644); err != nil {
			t.Fatal(err)
		}

		patterns := []string{"*.tmp"}
		snapshot, err := CreateDirectorySnapshot(tmpDir, patterns)
		if err != nil {
			t.Fatal(err)
		}

		// Should only include .txt file, not .tmp file
		if len(snapshot.Files) != 1 {
			t.Errorf("Expected 1 file after exclusion, got %d", len(snapshot.Files))
		}
		if snapshot.Files[0].RelativePath != "include.txt" {
			t.Errorf("Expected include.txt, got %s", snapshot.Files[0].RelativePath)
		}
	})
}

// TestCreateDirectorySnapshot tests directory snapshot creation
func TestCreateDirectorySnapshot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with known content
	files := map[string]string{
		"file1.txt":        "content1",
		"file2.txt":        "content2",
		"subdir/file3.txt": "content3",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(tmpDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	snapshot, err := CreateDirectorySnapshot(tmpDir, []string{})
	if err != nil {
		t.Fatal(err)
	}

	// Should have created directory and 3 files (4 total entries)
	expectedFiles := 4 // subdir + 3 files
	if len(snapshot.Files) != expectedFiles {
		t.Errorf("Expected %d files in snapshot, got %d", expectedFiles, len(snapshot.Files))
	}

	// Verify files are sorted
	for i := 1; i < len(snapshot.Files); i++ {
		if snapshot.Files[i-1].RelativePath >= snapshot.Files[i].RelativePath {
			t.Error("Files not sorted in snapshot")
		}
	}
}

// TestCompareSnapshots tests snapshot comparison functionality
func TestCompareSnapshots(t *testing.T) {
	// Create test snapshots
	snapshot1 := &DirectorySnapshot{
		Files: []FileInfo{
			{RelativePath: "file1.txt", Size: 100, Hash: "hash1", IsDir: false},
			{RelativePath: "file2.txt", Size: 200, Hash: "hash2", IsDir: false},
		},
	}

	snapshot2 := &DirectorySnapshot{
		Files: []FileInfo{
			{RelativePath: "file1.txt", Size: 100, Hash: "hash1", IsDir: false},
			{RelativePath: "file2.txt", Size: 200, Hash: "hash2", IsDir: false},
		},
	}

	snapshot3 := &DirectorySnapshot{
		Files: []FileInfo{
			{RelativePath: "file1.txt", Size: 100, Hash: "hash1", IsDir: false},
			{RelativePath: "file2.txt", Size: 200, Hash: "different_hash", IsDir: false},
		},
	}

	// Test identical snapshots
	if !CompareSnapshots(snapshot1, snapshot2) {
		t.Error("Expected identical snapshots to be equal")
	}

	// Test different snapshots
	if CompareSnapshots(snapshot1, snapshot3) {
		t.Error("Expected different snapshots to be not equal")
	}

	// Test different lengths
	snapshot4 := &DirectorySnapshot{
		Files: []FileInfo{
			{RelativePath: "file1.txt", Size: 100, Hash: "hash1", IsDir: false},
		},
	}

	if CompareSnapshots(snapshot1, snapshot4) {
		t.Error("Expected snapshots with different lengths to be not equal")
	}
}
