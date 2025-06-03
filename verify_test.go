// This file is part of bkpdir

// Package main provides tests for archive verification functionality.
// It verifies integrity checking and checksum validation.
package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

// TestArchiveData holds test archive setup data
type TestArchiveData struct {
	TestFile    string
	ArchivePath string
	Archive     *Archive
}

func setupTestArchive(t *testing.T) (*TestArchiveData, func()) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bkpdir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a test archive
	archivePath := filepath.Join(tempDir, "test.zip")
	archive := &Archive{
		Name: "test.zip",
		Path: archivePath,
	}

	// Create a simple zip archive
	zipFile, err := os.Create(archivePath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create zip file: %v", err)
	}

	// Add the test file to the archive
	zipWriter := zip.NewWriter(zipFile)
	fileWriter, err := zipWriter.Create("test.txt")
	if err != nil {
		zipFile.Close()
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create file in zip: %v", err)
	}

	if _, err := fileWriter.Write([]byte("test content")); err != nil {
		zipWriter.Close()
		zipFile.Close()
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to write to zip: %v", err)
	}

	// Close the zip writer to ensure all data is written
	zipWriter.Close()
	zipFile.Close()

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	data := &TestArchiveData{
		TestFile:    testFile,
		ArchivePath: archivePath,
		Archive:     archive,
	}

	return data, cleanup
}

// ⭐ ARCH-002: Archive verification functionality validation - 🔧
// TEST-REF: TestVerifyArchive
func TestVerifyArchive(t *testing.T) {
	data, cleanup := setupTestArchive(t)
	defer cleanup()

	// Test basic verification
	testBasicVerification(t, data)

	// Test checksum verification
	testChecksumVerification(t, data)

	// Test status storage and loading
	testStatusStorageAndLoading(t, data)
}

func testBasicVerification(t *testing.T, data *TestArchiveData) {
	status, err := VerifyArchive(data.ArchivePath)
	if err != nil {
		t.Fatalf("VerifyArchive failed: %v", err)
	}

	if !status.IsVerified {
		t.Errorf("Archive verification failed: %v", status.Errors)
	}
}

func testChecksumVerification(t *testing.T, data *TestArchiveData) {
	// Generate checksums before verification
	fileMap := map[string]string{
		"test.txt": data.TestFile,
	}
	checksums, err := GenerateChecksums(fileMap, "sha256")
	if err != nil {
		t.Fatalf("GenerateChecksums failed: %v", err)
	}

	// Store checksums in the archive
	if err := StoreChecksums(data.Archive, checksums); err != nil {
		t.Fatalf("StoreChecksums failed: %v", err)
	}

	// Test checksum verification
	checksumStatus, err := VerifyChecksums(data.ArchivePath)
	if err != nil {
		t.Fatalf("VerifyChecksums failed: %v", err)
	}

	if !checksumStatus.IsVerified {
		t.Errorf("Checksum verification failed: %v", checksumStatus.Errors)
	}

	if !checksumStatus.HasChecksums {
		t.Errorf("HasChecksums should be true")
	}
}

func testStatusStorageAndLoading(t *testing.T, data *TestArchiveData) {
	// First verify the archive to get a status
	status, err := VerifyArchive(data.ArchivePath)
	if err != nil {
		t.Fatalf("VerifyArchive failed: %v", err)
	}

	// Test storing verification status
	if err := StoreVerificationStatus(data.Archive, status); err != nil {
		t.Fatalf("StoreVerificationStatus failed: %v", err)
	}

	// Test loading verification status
	loadedStatus, err := LoadVerificationStatus(data.Archive)
	if err != nil {
		t.Fatalf("LoadVerificationStatus failed: %v", err)
	}

	if loadedStatus == nil {
		t.Fatalf("Loaded status is nil")
	}

	if !loadedStatus.IsVerified {
		t.Errorf("Loaded status indicates verification failed: %v", loadedStatus.Errors)
	}
}

func TestVerifyCorruptedArchive(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bkpdir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a corrupted zip file
	archivePath := filepath.Join(tempDir, "corrupted.zip")
	if err := os.WriteFile(archivePath, []byte("not a zip file"), 0644); err != nil {
		t.Fatalf("Failed to create corrupted file: %v", err)
	}

	// Test verification
	status, err := VerifyArchive(archivePath)
	if err != nil {
		t.Fatalf("VerifyArchive failed: %v", err)
	}

	if status.IsVerified {
		t.Errorf("Corrupted archive should not be verified")
	}

	if len(status.Errors) == 0 {
		t.Errorf("Expected errors for corrupted archive")
	}
}
