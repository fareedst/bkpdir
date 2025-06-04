// This file is part of bkpdir

// Copyright 2024 BkpDir Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main provides comprehensive tests for archive creation, verification, and management.
// It tests both full and incremental archive operations with various configurations.
package main

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Test archive name generation
func TestGenerateArchiveName(t *testing.T) {
	tests := []struct {
		name     string
		config   ArchiveConfig
		expected string
	}{
		{
			name: "basic archive name",
			config: ArchiveConfig{
				Prefix:    "test",
				Timestamp: "2024-01-01-12-00",
			},
			expected: "test-2024-01-01-12-00.zip",
		},
		{
			name: "archive with git info",
			config: ArchiveConfig{
				Prefix:    "test",
				Timestamp: "2024-01-01-12-00",
				IsGit:     true,
				GitBranch: "main",
				GitHash:   "abc123",
			},
			expected: "test-2024-01-01-12-00=main=abc123.zip",
		},
		{
			name: "archive with git info and note",
			config: ArchiveConfig{
				Prefix:    "test",
				Timestamp: "2024-01-01-12-00",
				IsGit:     true,
				GitBranch: "main",
				GitHash:   "abc123",
				Note:      "backup",
			},
			expected: "test-2024-01-01-12-00=main=abc123=backup.zip",
		},
		{
			name: "archive with dirty git status",
			config: ArchiveConfig{
				Prefix:             "test",
				Timestamp:          "2024-01-01-12-00",
				IsGit:              true,
				GitBranch:          "main",
				GitHash:            "abc123",
				GitIsClean:         false,
				ShowGitDirtyStatus: true,
			},
			expected: "test-2024-01-01-12-00=main=abc123-dirty.zip",
		},
		{
			name: "archive with note only",
			config: ArchiveConfig{
				Prefix:    "test",
				Timestamp: "2024-01-01-12-00",
				Note:      "backup",
			},
			expected: "test-2024-01-01-12-00=backup.zip",
		},
		{
			name: "archive without prefix",
			config: ArchiveConfig{
				Timestamp: "2024-01-01-12-00",
			},
			expected: "2024-01-01-12-00.zip",
		},
		{
			name: "incremental archive",
			config: ArchiveConfig{
				Timestamp:     "2024-01-01-12-00",
				IsIncremental: true,
				BaseName:      "test-2024-01-01-10-00.zip",
			},
			expected: "test-2024-01-01-10-00_update=2024-01-01-12-00.zip",
		},
		{
			name: "incremental archive with git info",
			config: ArchiveConfig{
				Timestamp:     "2024-01-01-12-00",
				IsIncremental: true,
				BaseName:      "test-2024-01-01-10-00.zip",
				IsGit:         true,
				GitBranch:     "main",
				GitHash:       "abc123",
			},
			expected: "test-2024-01-01-10-00_update=2024-01-01-12-00=main=abc123.zip",
		},
		{
			name: "incremental archive with note",
			config: ArchiveConfig{
				Timestamp:     "2024-01-01-12-00",
				IsIncremental: true,
				BaseName:      "test-2024-01-01-10-00.zip",
				Note:          "update",
			},
			expected: "test-2024-01-01-10-00_update=2024-01-01-12-00=update.zip",
		},
		{
			name: "incremental archive with git dirty status",
			config: ArchiveConfig{
				Timestamp:          "2024-01-01-12-00",
				IsIncremental:      true,
				BaseName:           "test-2024-01-01-10-00.zip",
				IsGit:              true,
				GitBranch:          "main",
				GitHash:            "abc123",
				GitIsClean:         false,
				ShowGitDirtyStatus: true,
			},
			expected: "test-2024-01-01-10-00_update=2024-01-01-12-00=main=abc123-dirty.zip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateArchiveName(tt.config)
			if result != tt.expected {
				t.Errorf("GenerateArchiveName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCreateFullArchive tests the full archive creation functionality for ARCH-002 feature
func TestCreateFullArchive(t *testing.T) {
	// ⭐ ARCH-002: Full archive creation validation - 🔧
	// TEST-REF: Feature tracking matrix ARCH-002
	// IMMUTABLE-REF: Commands - Create Archive
	tempDir := t.TempDir()

	sourceDir := filepath.Join(tempDir, "source")
	archiveDir := filepath.Join(tempDir, "archives")

	// Create source directory with test files
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test files
	testFiles := map[string]string{
		"file1.txt":          "content1",
		"file2.txt":          "content2",
		"subdir/file3.txt":   "content3",
		"subdir/subfile.txt": "subcontent",
	}

	for filePath, content := range testFiles {
		fullPath := filepath.Join(sourceDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Save current directory and change to source directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	if err := os.Chdir(sourceDir); err != nil {
		t.Fatal(err)
	}

	// Create a test config with our test archive directory
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false // Don't create subdirectory

	// Test archive creation with valid config
	err = CreateFullArchive(cfg, "test-note", false, false)
	if err != nil {
		t.Fatalf("CreateFullArchive failed: %v", err)
	}

	// Verify archive was created
	archives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(archives) == 0 {
		t.Error("No archives were created")
	}

	// Test with exclusion patterns
	t.Run("WithExclusion", func(t *testing.T) {
		// Create config with exclusion patterns
		cfgExclude := DefaultConfig()
		cfgExclude.ArchiveDirPath = archiveDir
		cfgExclude.UseCurrentDirName = false
		cfgExclude.ExcludePatterns = []string{"subdir/*.txt"}

		// Create archive with exclusion
		err = CreateFullArchive(cfgExclude, "exclude-test", false, false)
		if err != nil {
			t.Fatalf("CreateFullArchive with exclusion failed: %v", err)
		}

		// Verify more archives were created
		newArchives, err := ListArchives(archiveDir)
		if err != nil {
			t.Fatal(err)
		}

		if len(newArchives) <= len(archives) {
			t.Error("Expected more archives after exclusion test")
		}
	})
}

// TestCreateIncremental tests the incremental archive creation functionality for ARCH-003 feature
func TestCreateIncremental(t *testing.T) {
	// ⭐ ARCH-003: Incremental archive creation validation - 🔧
	// TEST-REF: Feature tracking matrix ARCH-003
	// IMMUTABLE-REF: Commands - Create Incremental Archive
	tempDir := t.TempDir()

	sourceDir := filepath.Join(tempDir, "source")
	archiveDir := filepath.Join(tempDir, "archives")

	// Create directories
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create initial files
	initialFiles := map[string]string{
		"file1.txt": "content1",
		"file2.txt": "content2",
	}

	for filePath, content := range initialFiles {
		fullPath := filepath.Join(sourceDir, filePath)
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Save current directory and change to source directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	if err := os.Chdir(sourceDir); err != nil {
		t.Fatal(err)
	}

	// Create initial full archive
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false

	err = CreateFullArchive(cfg, "base-archive", false, false)
	if err != nil {
		t.Fatalf("CreateFullArchive failed: %v", err)
	}

	// Add new files for incremental
	newFiles := map[string]string{
		"file3.txt":        "content3",
		"subdir/file4.txt": "content4",
	}

	for filePath, content := range newFiles {
		fullPath := filepath.Join(sourceDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Get initial archive count
	initialArchives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatal(err)
	}

	// Create incremental archive
	err = CreateIncrementalArchive(cfg, "incremental-test", false, false)
	if err != nil {
		t.Fatalf("CreateIncrementalArchive failed: %v", err)
	}

	// Verify more archives exist
	finalArchives, err := ListArchives(archiveDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(finalArchives) <= len(initialArchives) {
		t.Error("Expected more archives after incremental creation")
	}

	// Test when no changes are needed
	t.Run("NoChangesNeeded", func(t *testing.T) {
		// Reset directory to match base archive state by removing new files
		for filePath := range newFiles {
			fullPath := filepath.Join(sourceDir, filePath)
			os.Remove(fullPath)
		}
		// Remove empty subdirectory
		os.Remove(filepath.Join(sourceDir, "subdir"))

		// This should detect that no incremental is needed
		err = CreateIncrementalArchive(cfg, "no-change-test", false, false)
		// The function should handle this case gracefully
		if err != nil {
			t.Logf("CreateIncrementalArchive with no changes: %v", err)
		}
	})
}

func TestSkipBrokenSymlinks(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bkpdir_symlink_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a regular file
	regularFile := filepath.Join(tempDir, "regular.txt")
	if err := os.WriteFile(regularFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create regular file: %v", err)
	}

	// Create a broken symlink
	brokenSymlink := filepath.Join(tempDir, "broken_link")
	if err := os.Symlink("/nonexistent/path", brokenSymlink); err != nil {
		t.Fatalf("Failed to create broken symlink: %v", err)
	}

	// Test with SkipBrokenSymlinks = true
	cfg := &Config{
		SkipBrokenSymlinks: true,
	}
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: cfg}

	// Create a zip writer for testing
	var buf bytes.Buffer
	zipw := zip.NewWriter(&buf)

	// Test adding the broken symlink - should not fail
	err = addFileToZipWithConfig(tempDir, "broken_link", zipw, archiveConfig)
	if err != nil {
		t.Errorf("Expected no error when skipping broken symlinks, got: %v", err)
	}

	// Test adding the regular file - should work
	err = addFileToZipWithConfig(tempDir, "regular.txt", zipw, archiveConfig)
	if err != nil {
		t.Errorf("Failed to add regular file: %v", err)
	}

	zipw.Close()

	// Test with SkipBrokenSymlinks = false
	cfg.SkipBrokenSymlinks = false
	var buf2 bytes.Buffer
	zipw2 := zip.NewWriter(&buf2)

	// Test adding the broken symlink - should fail
	err = addFileToZipWithConfig(tempDir, "broken_link", zipw2, archiveConfig)
	if err == nil {
		t.Error("Expected error when not skipping broken symlinks, got nil")
	}

	zipw2.Close()
}
