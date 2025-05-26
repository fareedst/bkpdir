package main

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestVerifyArchive(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "bkpdir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
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
		t.Fatalf("Failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Add the test file to the archive
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileWriter, err := zipWriter.Create("test.txt")
	if err != nil {
		t.Fatalf("Failed to create file in zip: %v", err)
	}

	if _, err := fileWriter.Write([]byte("test content")); err != nil {
		t.Fatalf("Failed to write to zip: %v", err)
	}

	// Close the zip writer to ensure all data is written
	zipWriter.Close()
	zipFile.Close()

	// Generate checksums before verification
	fileMap := map[string]string{
		"test.txt": testFile,
	}
	checksums, err := GenerateChecksums(fileMap, "sha256")
	if err != nil {
		t.Fatalf("GenerateChecksums failed: %v", err)
	}

	// Store checksums in the archive
	if err := StoreChecksums(archive, checksums); err != nil {
		t.Fatalf("StoreChecksums failed: %v", err)
	}

	// Test verification
	status, err := VerifyArchive(archivePath)
	if err != nil {
		t.Fatalf("VerifyArchive failed: %v", err)
	}

	if !status.IsVerified {
		t.Errorf("Archive verification failed: %v", status.Errors)
	}

	// Test storing and loading verification status
	if err := StoreVerificationStatus(archive, status); err != nil {
		t.Fatalf("StoreVerificationStatus failed: %v", err)
	}

	loadedStatus, err := LoadVerificationStatus(archive)
	if err != nil {
		t.Fatalf("LoadVerificationStatus failed: %v", err)
	}

	if loadedStatus == nil {
		t.Fatalf("Loaded status is nil")
	}

	if !loadedStatus.IsVerified {
		t.Errorf("Loaded status indicates verification failed: %v", loadedStatus.Errors)
	}

	// Test checksum verification
	checksumStatus, err := VerifyChecksums(archivePath)
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
