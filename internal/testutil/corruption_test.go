// This file is part of bkpdir
//
// Package testutil provides comprehensive tests for archive corruption utilities.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// TEST-INFRA-001-A: Archive corruption testing framework tests
// DECISION-REF: DEC-001 (ZIP format)
// IMPLEMENTATION-NOTES: Comprehensive test coverage for corruption utilities

package testutil

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// Test data for archive corruption testing
var testFiles = map[string]string{
	"file1.txt":        "This is test file 1 content",
	"file2.txt":        "This is test file 2 content with more data",
	"subdir/file3.txt": "This is test file 3 in a subdirectory",
	"large.txt":        "Large file content: " + generateLargeContent(1000),
}

// generateLargeContent creates large test content
func generateLargeContent(size int) string {
	content := make([]byte, size)
	for i := range content {
		content[i] = byte('A' + (i % 26))
	}
	return string(content)
}

// TestCorruptionType tests the corruption type enumeration
func TestCorruptionType(t *testing.T) {
	tests := []struct {
		name string
		ct   CorruptionType
		want int
	}{
		{"CRC corruption", CorruptionCRC, 0},
		{"Header corruption", CorruptionHeader, 1},
		{"Truncate corruption", CorruptionTruncate, 2},
		{"Central directory corruption", CorruptionCentralDir, 3},
		{"Local header corruption", CorruptionLocalHeader, 4},
		{"Data corruption", CorruptionData, 5},
		{"Signature corruption", CorruptionSignature, 6},
		{"Comment corruption", CorruptionComment, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.ct) != tt.want {
				t.Errorf("CorruptionType value = %d, want %d", int(tt.ct), tt.want)
			}
		})
	}
}

// TestCreateTestArchive tests the test archive creation utility
func TestCreateTestArchive(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("CreateTestArchive failed: %v", err)
	}

	// Verify archive was created and contains expected files
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Fatalf("Failed to open created archive: %v", err)
	}
	defer reader.Close()

	if len(reader.File) != len(testFiles) {
		t.Errorf("Archive contains %d files, expected %d", len(reader.File), len(testFiles))
	}

	// Verify each file exists and has correct content
	for _, file := range reader.File {
		expectedContent, exists := testFiles[file.Name]
		if !exists {
			t.Errorf("Unexpected file in archive: %s", file.Name)
			continue
		}

		rc, err := file.Open()
		if err != nil {
			t.Errorf("Failed to open file %s: %v", file.Name, err)
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file.Name, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("File %s content mismatch: got %q, want %q", file.Name, string(content), expectedContent)
		}
	}
}

// TestNewArchiveCorruptor tests archive corruptor creation
func TestNewArchiveCorruptor(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Test with non-existent archive
	t.Run("nonexistent archive", func(t *testing.T) {
		config := CorruptionConfig{Type: CorruptionCRC}
		_, err := NewArchiveCorruptor("/nonexistent/archive.zip", config)
		if err == nil {
			t.Error("Expected error for nonexistent archive")
		}
	})

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Test with existing archive
	t.Run("existing archive", func(t *testing.T) {
		config := CorruptionConfig{Type: CorruptionCRC}
		corruptor, err := NewArchiveCorruptor(archivePath, config)
		if err != nil {
			t.Errorf("NewArchiveCorruptor failed: %v", err)
		}
		if corruptor == nil {
			t.Error("Expected non-nil corruptor")
		}
		if corruptor.originalPath != archivePath {
			t.Errorf("Original path = %s, want %s", corruptor.originalPath, archivePath)
		}
	})
}

// TestArchiveCorruptor_Backup tests backup functionality
func TestArchiveCorruptor_Backup(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{Type: CorruptionCRC}
	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Test backup creation
	err = corruptor.CreateBackup()
	if err != nil {
		t.Errorf("CreateBackup failed: %v", err)
	}

	// Verify backup file exists
	if _, err := os.Stat(corruptor.backupPath); os.IsNotExist(err) {
		t.Error("Backup file was not created")
	}

	// Verify backup content matches original
	originalData, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("Failed to read original archive: %v", err)
	}

	backupData, err := os.ReadFile(corruptor.backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup archive: %v", err)
	}

	if !bytes.Equal(originalData, backupData) {
		t.Error("Backup content does not match original")
	}
}

// TestArchiveCorruptor_RestoreFromBackup tests restore functionality
func TestArchiveCorruptor_RestoreFromBackup(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Read original data
	originalData, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("Failed to read original archive: %v", err)
	}

	config := CorruptionConfig{Type: CorruptionCRC}
	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Create backup and corrupt archive
	err = corruptor.CreateBackup()
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Corrupt the archive by writing random data
	err = os.WriteFile(archivePath, []byte("corrupted data"), 0644)
	if err != nil {
		t.Fatalf("Failed to corrupt archive: %v", err)
	}

	// Restore from backup
	err = corruptor.RestoreFromBackup()
	if err != nil {
		t.Errorf("RestoreFromBackup failed: %v", err)
	}

	// Verify restored content matches original
	restoredData, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("Failed to read restored archive: %v", err)
	}

	if !bytes.Equal(originalData, restoredData) {
		t.Error("Restored content does not match original")
	}

	// Test restore without backup
	t.Run("no backup file", func(t *testing.T) {
		config2 := CorruptionConfig{Type: CorruptionCRC}
		corruptor2, err := NewArchiveCorruptor(archivePath, config2)
		if err != nil {
			t.Fatalf("NewArchiveCorruptor failed: %v", err)
		}

		err = corruptor2.RestoreFromBackup()
		if err == nil {
			t.Error("Expected error when restoring without backup")
		}
	})
}

// TestCorruptionCRC tests CRC corruption
func TestCorruptionCRC(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type: CorruptionCRC,
		Seed: 12345, // For reproducible corruption
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply CRC corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionCRC {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionCRC)
	}

	if !result.Recoverable {
		t.Error("CRC corruption should be marked as recoverable")
	}

	if len(result.AppliedAt) == 0 {
		t.Error("No corruption locations recorded")
	}

	if len(result.OriginalBytes) == 0 {
		t.Error("No original bytes recorded")
	}

	// Verify archive still opens but has CRC errors
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Errorf("Archive should still open after CRC corruption: %v", err)
	} else {
		reader.Close()
	}
}

// TestCorruptionHeader tests header corruption
func TestCorruptionHeader(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type: CorruptionHeader,
		Seed: 54321,
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply header corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionHeader {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionHeader)
	}

	if result.Recoverable {
		t.Error("Header corruption should be marked as non-recoverable")
	}

	// Try to open archive - should either fail or have major issues
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		// Expected - archive should not open
		return
	}
	defer reader.Close()

	// If it does open (Go's zip reader is resilient), it should have problems
	if len(reader.File) == 0 {
		// No files readable - this is also corruption success
		return
	}

	// Try reading files - should have issues
	hasErrors := false
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			hasErrors = true
			break
		}
		_, err = io.ReadAll(rc)
		rc.Close()
		if err != nil {
			hasErrors = true
			break
		}
	}

	if !hasErrors {
		t.Error("Header corruption should cause archive reading issues")
	}
}

// TestCorruptionTruncate tests truncation corruption
func TestCorruptionTruncate(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Get original size
	originalInfo, err := os.Stat(archivePath)
	if err != nil {
		t.Fatalf("Failed to stat original archive: %v", err)
	}
	originalSize := originalInfo.Size()

	config := CorruptionConfig{
		Type:           CorruptionTruncate,
		CorruptionSize: 100, // Truncate 100 bytes
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply truncation corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionTruncate {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionTruncate)
	}

	if result.Recoverable {
		t.Error("Truncation corruption should be marked as non-recoverable")
	}

	// Verify file is actually truncated
	truncatedInfo, err := os.Stat(archivePath)
	if err != nil {
		t.Fatalf("Failed to stat truncated archive: %v", err)
	}
	truncatedSize := truncatedInfo.Size()

	expectedSize := originalSize - 100
	if truncatedSize != expectedSize {
		t.Errorf("Truncated size = %d, want %d", truncatedSize, expectedSize)
	}
}

// TestCorruptionCentralDirectory tests central directory corruption
func TestCorruptionCentralDirectory(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type: CorruptionCentralDir,
		Seed: 98765,
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply central directory corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionCentralDir {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionCentralDir)
	}

	if result.Recoverable {
		t.Error("Central directory corruption should be marked as non-recoverable")
	}

	// Archive might still open but should have issues accessing files
	reader, err := zip.OpenReader(archivePath)
	if err == nil {
		// If it opens, try to access files
		for _, file := range reader.File {
			rc, err := file.Open()
			if err == nil {
				rc.Close()
			}
		}
		reader.Close()
	}
}

// TestCorruptionData tests file data corruption
func TestCorruptionData(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type:           CorruptionData,
		Seed:           11111,
		CorruptionSize: 50, // Corrupt 50 bytes of data
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply data corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionData {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionData)
	}

	if !result.Recoverable {
		t.Error("Data corruption should be marked as recoverable")
	}

	if len(result.CorruptedBytes) != 50 {
		t.Errorf("Corrupted bytes length = %d, want 50", len(result.CorruptedBytes))
	}
}

// TestCorruptionSignature tests signature corruption
func TestCorruptionSignature(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type: CorruptionSignature,
		Seed: 77777,
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply signature corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionSignature {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionSignature)
	}

	if result.Recoverable {
		t.Error("Signature corruption should be marked as non-recoverable")
	}

	// Try to open archive - should either fail or have major issues
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		// Expected - archive should not open
		return
	}
	defer reader.Close()

	// If it does open despite signature corruption, it should have problems
	if len(reader.File) == 0 {
		// No files readable - this is corruption success
		return
	}

	// Try reading files - should have issues
	hasErrors := false
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			hasErrors = true
			break
		}
		_, err = io.ReadAll(rc)
		rc.Close()
		if err != nil {
			hasErrors = true
			break
		}
	}

	if !hasErrors {
		t.Error("Signature corruption should cause archive reading issues")
	}
}

// TestCorruptionComment tests comment corruption
func TestCorruptionComment(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	config := CorruptionConfig{
		Type: CorruptionComment,
		Seed: 88888,
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Apply comment corruption
	result, err := corruptor.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption failed: %v", err)
	}

	// Verify corruption result
	if result.Type != CorruptionComment {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionComment)
	}

	if !result.Recoverable {
		t.Error("Comment corruption should be marked as recoverable")
	}

	// Archive should still open after comment corruption
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		t.Errorf("Archive should open after comment corruption: %v", err)
	} else {
		reader.Close()
	}
}

// TestCorruptionReproducibility tests that corruption is reproducible with same seed
func TestCorruptionReproducibility(t *testing.T) {
	tempDir := t.TempDir()

	// Create identical archives by copying bytes rather than recreating
	archive1Path := filepath.Join(tempDir, "test1.zip")
	archive2Path := filepath.Join(tempDir, "test2.zip")

	// Create first archive
	err := CreateTestArchive(archive1Path, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive 1: %v", err)
	}

	// Copy the exact same bytes to create identical archive
	data, err := os.ReadFile(archive1Path)
	if err != nil {
		t.Fatalf("Failed to read first archive: %v", err)
	}

	err = os.WriteFile(archive2Path, data, 0644)
	if err != nil {
		t.Fatalf("Failed to create identical archive: %v", err)
	}

	// Verify archives are identical before corruption
	data1, err := os.ReadFile(archive1Path)
	if err != nil {
		t.Fatalf("Failed to read archive 1: %v", err)
	}
	data2, err := os.ReadFile(archive2Path)
	if err != nil {
		t.Fatalf("Failed to read archive 2: %v", err)
	}

	if !bytes.Equal(data1, data2) {
		t.Fatal("Archives are not identical before corruption")
	}

	// Apply same corruption with same seed to both
	seed := int64(99999)
	config := CorruptionConfig{
		Type: CorruptionCRC,
		Seed: seed,
	}

	// Corrupt first archive
	corruptor1, err := NewArchiveCorruptor(archive1Path, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor 1 failed: %v", err)
	}
	defer corruptor1.Cleanup()

	result1, err := corruptor1.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption 1 failed: %v", err)
	}

	// Corrupt second archive
	corruptor2, err := NewArchiveCorruptor(archive2Path, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor 2 failed: %v", err)
	}
	defer corruptor2.Cleanup()

	result2, err := corruptor2.ApplyCorruption()
	if err != nil {
		t.Fatalf("ApplyCorruption 2 failed: %v", err)
	}

	// Verify both corruptions are identical
	if len(result1.AppliedAt) != len(result2.AppliedAt) {
		t.Errorf("Different number of corruption locations: %d vs %d", len(result1.AppliedAt), len(result2.AppliedAt))
	}

	for i, offset1 := range result1.AppliedAt {
		if i < len(result2.AppliedAt) {
			offset2 := result2.AppliedAt[i]
			if offset1 != offset2 {
				t.Errorf("Corruption offset %d differs: %d vs %d", i, offset1, offset2)
			}
		}
	}

	if !bytes.Equal(result1.CorruptedBytes, result2.CorruptedBytes) {
		t.Error("Corrupted bytes are different between identical operations")
	}

	// Verify the actual corrupted files are identical
	corruptedData1, err := os.ReadFile(archive1Path)
	if err != nil {
		t.Fatalf("Failed to read corrupted archive 1: %v", err)
	}

	corruptedData2, err := os.ReadFile(archive2Path)
	if err != nil {
		t.Fatalf("Failed to read corrupted archive 2: %v", err)
	}

	if !bytes.Equal(corruptedData1, corruptedData2) {
		t.Error("Corrupted archives are not identical despite same seed")
	}
}

// TestCorruptionDetector tests the corruption detection functionality
func TestCorruptionDetector(t *testing.T) {
	tempDir := t.TempDir()

	// Test various types of corruption detection
	testCases := []struct {
		name              string
		corruptionType    CorruptionType
		expectedDetection []CorruptionType
	}{
		{
			name:              "signature corruption",
			corruptionType:    CorruptionSignature,
			expectedDetection: []CorruptionType{CorruptionSignature},
		},
		{
			name:              "header corruption",
			corruptionType:    CorruptionHeader,
			expectedDetection: []CorruptionType{CorruptionHeader, CorruptionSignature}, // May detect either
		},
		{
			name:              "truncation corruption",
			corruptionType:    CorruptionTruncate,
			expectedDetection: []CorruptionType{CorruptionTruncate},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			archivePath := filepath.Join(tempDir, tc.name+".zip")

			// Create and corrupt archive
			config := CorruptionConfig{
				Type:           tc.corruptionType,
				Seed:           12345,
				CorruptionSize: 100,
				Severity:       0.1,
			}

			_, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
			if err != nil {
				t.Fatalf("Failed to create corrupted archive: %v", err)
			}

			// Test detection
			detector := NewCorruptionDetector(archivePath)
			detected, err := detector.DetectCorruption()
			if err != nil {
				t.Errorf("DetectCorruption failed: %v", err)
			}

			// Verify at least one expected corruption type was detected
			found := false
			for _, expectedType := range tc.expectedDetection {
				for _, detectedType := range detected {
					if detectedType == expectedType {
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			if !found {
				t.Errorf("Expected corruption type(s) %v not detected, got %v", tc.expectedDetection, detected)
			}
		})
	}
}

// TestCreateCorruptedTestArchive tests the convenience function
func TestCreateCorruptedTestArchive(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "corrupted.zip")

	config := CorruptionConfig{
		Type: CorruptionCRC,
		Seed: 55555,
	}

	result, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
	if err != nil {
		t.Fatalf("CreateCorruptedTestArchive failed: %v", err)
	}

	// Verify archive was created and corrupted
	if result.Type != CorruptionCRC {
		t.Errorf("Corruption type = %v, want %v", result.Type, CorruptionCRC)
	}

	if len(result.AppliedAt) == 0 {
		t.Error("No corruption was applied")
	}

	// Verify archive exists
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("Corrupted archive was not created")
	}
}

// TestUnsupportedCorruptionType tests error handling for unsupported corruption types
func TestUnsupportedCorruptionType(t *testing.T) {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test.zip")

	// Create test archive
	err := CreateTestArchive(archivePath, testFiles)
	if err != nil {
		t.Fatalf("Failed to create test archive: %v", err)
	}

	// Use invalid corruption type
	config := CorruptionConfig{
		Type: CorruptionType(999), // Invalid type
	}

	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		t.Fatalf("NewArchiveCorruptor failed: %v", err)
	}
	defer corruptor.Cleanup()

	// Should fail with unsupported type
	_, err = corruptor.ApplyCorruption()
	if err == nil {
		t.Error("Expected error for unsupported corruption type")
	}
}

// TestEdgeCases tests various edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("empty archive", func(t *testing.T) {
		tempDir := t.TempDir()
		archivePath := filepath.Join(tempDir, "empty.zip")

		// Create empty archive
		err := CreateTestArchive(archivePath, map[string]string{})
		if err != nil {
			t.Fatalf("Failed to create empty archive: %v", err)
		}

		config := CorruptionConfig{Type: CorruptionCRC}
		corruptor, err := NewArchiveCorruptor(archivePath, config)
		if err != nil {
			t.Fatalf("NewArchiveCorruptor failed: %v", err)
		}
		defer corruptor.Cleanup()

		// Should handle empty archive gracefully
		_, err = corruptor.ApplyCorruption()
		// May or may not succeed depending on what's in an "empty" ZIP
		// Just ensure it doesn't panic
	})

	t.Run("very small archive", func(t *testing.T) {
		tempDir := t.TempDir()
		archivePath := filepath.Join(tempDir, "tiny.zip")

		// Create tiny archive
		err := CreateTestArchive(archivePath, map[string]string{"tiny.txt": "x"})
		if err != nil {
			t.Fatalf("Failed to create tiny archive: %v", err)
		}

		config := CorruptionConfig{
			Type:     CorruptionSignature,
			Severity: 1.0, // Maximum corruption
		}

		corruptor, err := NewArchiveCorruptor(archivePath, config)
		if err != nil {
			t.Fatalf("NewArchiveCorruptor failed: %v", err)
		}
		defer corruptor.Cleanup()

		// Should handle small archive
		result, err := corruptor.ApplyCorruption()
		if err != nil {
			t.Errorf("ApplyCorruption failed on small archive: %v", err)
		} else if result == nil {
			t.Error("Expected non-nil result")
		}
	})
}

// Benchmark tests for performance
func BenchmarkCorruptionCRC(b *testing.B) {
	tempDir := b.TempDir()

	// Create a larger test archive for benchmarking
	largeFiles := map[string]string{}
	for i := 0; i < 10; i++ {
		largeFiles[fmt.Sprintf("file%d.txt", i)] = generateLargeContent(10000)
	}

	config := CorruptionConfig{
		Type: CorruptionCRC,
		Seed: 12345,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		archivePath := filepath.Join(tempDir, fmt.Sprintf("bench%d.zip", i))

		// Create and corrupt archive
		_, err := CreateCorruptedTestArchive(archivePath, largeFiles, config)
		if err != nil {
			b.Fatalf("CreateCorruptedTestArchive failed: %v", err)
		}

		// Clean up
		os.Remove(archivePath)
	}
}

func BenchmarkCorruptionDetection(b *testing.B) {
	tempDir := b.TempDir()
	archivePath := filepath.Join(tempDir, "bench.zip")

	// Create corrupted archive once
	config := CorruptionConfig{
		Type: CorruptionSignature,
		Seed: 12345,
	}

	_, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
	if err != nil {
		b.Fatalf("CreateCorruptedTestArchive failed: %v", err)
	}

	detector := NewCorruptionDetector(archivePath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := detector.DetectCorruption()
		if err != nil {
			b.Fatalf("DetectCorruption failed: %v", err)
		}
	}
}
