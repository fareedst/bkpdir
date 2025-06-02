// This file is part of bkpdir
// LINT-001: Lint compliance

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestCompareFiles tests the file comparison functionality for FILE-003 feature
func TestCompareFiles(t *testing.T) {
	// FILE-003: File comparison validation
	// TEST-REF: Feature tracking matrix FILE-003
	// IMMUTABLE-REF: Identical File Detection

	t.Run("IdenticalSnapshots", func(t *testing.T) {
		// Create temporary directories for testing
		tmpDir1, err := ioutil.TempDir("", "compare_test1_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir1)

		tmpDir2, err := ioutil.TempDir("", "compare_test2_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir2)

		// Create identical files in both directories
		content := "test content"
		file1 := filepath.Join(tmpDir1, "test.txt")
		file2 := filepath.Join(tmpDir2, "test.txt")

		if err := ioutil.WriteFile(file1, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(file2, []byte(content), 0644); err != nil {
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
		tmpDir1, err := ioutil.TempDir("", "compare_test1_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir1)

		tmpDir2, err := ioutil.TempDir("", "compare_test2_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir2)

		// Create different files
		file1 := filepath.Join(tmpDir1, "different1.txt")
		file2 := filepath.Join(tmpDir2, "different2.txt")

		if err := ioutil.WriteFile(file1, []byte("content1"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(file2, []byte("content2"), 0644); err != nil {
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
		tmpDir, err := ioutil.TempDir("", "exclude_test_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		// Create files, some to be excluded
		testFile := filepath.Join(tmpDir, "include.txt")
		excludeFile := filepath.Join(tmpDir, "exclude.tmp")

		if err := ioutil.WriteFile(testFile, []byte("include"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(excludeFile, []byte("exclude"), 0644); err != nil {
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
	tmpDir, err := ioutil.TempDir("", "snapshot_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
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
