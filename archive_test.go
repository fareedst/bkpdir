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

// Package main provides tests for archive creation and management in BkpDir.
// It verifies archive creation, naming, and verification behavior.
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateArchiveName tests archive naming functionality
func TestGenerateArchiveName(t *testing.T) {
	// ARCH-001: Archive naming convention validation
	// TEST-REF: Feature tracking matrix ARCH-001
	// IMMUTABLE-REF: Archive Naming Convention
	tests := []struct {
		cfg    ArchiveNameConfig
		expect string
		name   string
	}{
		{
			cfg: ArchiveNameConfig{
				Prefix:        "HOME",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "main",
				GitHash:       "abc123",
				GitIsClean:    true,
				Note:          "",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "HOME-2024-06-01-12-00=main=abc123.zip",
			name:   "full with git",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "note",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=note.zip",
			name:   "full with note only",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "dev",
				GitHash:       "def456",
				GitIsClean:    true,
				Note:          "backup",
				IsGit:         true,
				IsIncremental: true,
				BaseName:      "HOME-2024-06-01-12-00=main=abc123",
			},
			expect: "HOME-2024-06-01-12-00=main=abc123_update=2024-06-01-12-00=dev=def456=backup.zip",
			name:   "incremental with git and note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "",
				IsGit:         false,
				IsIncremental: true,
				BaseName:      "ARCHIVE-2024-06-01-12-00",
			},
			expect: "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00.zip",
			name:   "incremental no git, no note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00.zip",
			name:   "full no git, no note",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00.zip",
			name:   "full git, no branch/hash",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "test note",
				IsGit:         false,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=test note.zip",
			name:   "full with note containing spaces",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "",
				GitHash:       "",
				GitIsClean:    true,
				Note:          "test-note",
				IsGit:         false,
				IsIncremental: true,
				BaseName:      "ARCHIVE-2024-06-01-12-00",
			},
			expect: "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00=test-note.zip",
			name:   "incremental with note only",
		},
		{
			cfg: ArchiveNameConfig{
				Prefix:        "",
				Timestamp:     "2024-06-01-12-00",
				GitBranch:     "main",
				GitHash:       "abc123",
				GitIsClean:    true,
				Note:          "test note",
				IsGit:         true,
				IsIncremental: false,
				BaseName:      "",
			},
			expect: "2024-06-01-12-00=main=abc123=test note.zip",
			name:   "full with git and note with spaces",
		},
	}
	for _, tt := range tests {
		got := GenerateArchiveName(tt.cfg)
		if got != tt.expect {
			t.Errorf("%s: GenerateArchiveName() = %q, want %q", tt.name, got, tt.expect)
		}
	}
}

// TestCreateFullArchive tests the full archive creation functionality for ARCH-002 feature
func TestCreateFullArchive(t *testing.T) {
	// ARCH-002: Full archive creation validation
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
	// ARCH-003: Incremental archive creation validation
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
