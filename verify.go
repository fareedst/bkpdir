// This file is part of bkpdir
//
// Package main provides archive verification for BkpDir.
// It handles integrity checking and checksum verification of archives.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// VerificationStatus represents the result of an archive verification
type VerificationStatus struct {
	VerifiedAt   time.Time `json:"verified_at"`
	IsVerified   bool      `json:"is_verified"`
	HasChecksums bool      `json:"has_checksums"`
	Errors       []string  `json:"errors,omitempty"`
}

// VerifyArchive verifies the integrity of an archive
func VerifyArchive(archivePath string) (*VerificationStatus, error) {
	// Archive verification implementation
	// DECISION-REF: DEC-001
	status := &VerificationStatus{
		VerifiedAt: time.Now(),
		IsVerified: true,
	}

	// Open the archive
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		status.IsVerified = false
		status.Errors = append(status.Errors, fmt.Sprintf("Failed to open archive: %v", err))
		return status, nil
	}
	defer reader.Close()

	// Check each file in the archive
	for _, file := range reader.File {
		if err := verifyFile(file); err != nil {
			status.IsVerified = false
			status.Errors = append(status.Errors, err.Error())
		}
	}

	return status, nil
}

// verifyFile verifies a single file in the archive
func verifyFile(file *zip.File) error {
	// Individual file verification
	// DECISION-REF: DEC-001
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", file.Name, err)
	}
	defer rc.Close()

	// Read a small portion to verify the file can be read
	buf := make([]byte, 1024)
	_, err = rc.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file %s: %v", file.Name, err)
	}

	return nil
}

// GenerateChecksums generates checksums for files in the map
func GenerateChecksums(fileMap map[string]string, _ string) (map[string]string, error) {
	// Checksum generation for verification
	// DECISION-REF: DEC-001
	checksums := make(map[string]string)
	for relPath, absPath := range fileMap {
		checksum, err := calculateFileChecksum(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate checksum for %s: %w", relPath, err)
		}
		checksums[relPath] = checksum
	}
	return checksums, nil
}

// calculateFileChecksum calculates the SHA-256 checksum of a file
func calculateFileChecksum(filePath string) (string, error) {
	// SHA-256 checksum calculation
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// StoreChecksums stores checksums in the archive
func StoreChecksums(archive *Archive, checksums map[string]string) error {
	// Checksum storage in archive
	// DECISION-REF: DEC-001, DEC-008
	// Create a temporary file for checksums
	tmpFile, err := createChecksumsTempFile(checksums)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Create a new archive with the checksums file
	newPath := archive.Path + ".new"
	if err := createNewArchiveWithChecksums(archive.Path, newPath, tmpFile.Name()); err != nil {
		return err
	}
	defer os.Remove(newPath)

	// Replace the original archive with the new one
	if err := os.Rename(newPath, archive.Path); err != nil {
		return fmt.Errorf("failed to replace original archive: %w", err)
	}

	return nil
}

// createChecksumsTempFile creates a temporary file containing the checksums
func createChecksumsTempFile(checksums map[string]string) (*os.File, error) {
	// Temporary checksum file creation
	// DECISION-REF: DEC-008
	tmpFile, err := os.CreateTemp("", "bkpdir-checksums-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	encoder := json.NewEncoder(tmpFile)
	if err := encoder.Encode(checksums); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to encode checksums: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	return tmpFile, nil
}

// createNewArchiveWithChecksums creates a new archive with the checksums file
func createNewArchiveWithChecksums(archivePath, newPath, checksumsPath string) error {
	// Archive reconstruction with checksums
	// DECISION-REF: DEC-001, DEC-008
	// Open the original archive
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	// Create the new archive
	writer, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("failed to create new archive: %w", err)
	}
	defer writer.Close()

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// Copy all files from the original archive
	if err := copyArchiveFiles(reader, zipWriter); err != nil {
		return err
	}

	// Add the checksums file
	return addChecksumsFile(zipWriter, checksumsPath)
}

// copyArchiveFiles copies all files from the original archive to the new one
func copyArchiveFiles(reader *zip.ReadCloser, writer *zip.Writer) error {
	// Archive file copying during reconstruction
	for _, file := range reader.File {
		if err := copyArchiveFile(file, writer); err != nil {
			return err
		}
	}
	return nil
}

// copyArchiveFile copies a single file from the original archive to the new one
func copyArchiveFile(file *zip.File, writer *zip.Writer) error {
	// Individual file copying in archive
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer rc.Close()

	header, err := zip.FileInfoHeader(file.FileInfo())
	if err != nil {
		return fmt.Errorf("failed to create header: %w", err)
	}
	header.Name = file.Name

	w, err := writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create file in new archive: %w", err)
	}

	if _, err := io.Copy(w, rc); err != nil {
		return fmt.Errorf("failed to copy file to new archive: %w", err)
	}

	return nil
}

// addChecksumsFile adds the checksums file to the archive
func addChecksumsFile(writer *zip.Writer, checksumsPath string) error {
	// ARCH-002: Checksum file addition to archive
	// DECISION-REF: DEC-001
	checksumFile, err := os.Open(checksumsPath)
	if err != nil {
		return fmt.Errorf("failed to open checksums file: %w", err)
	}
	defer checksumFile.Close()

	header := &zip.FileHeader{
		Name:   ".checksums",
		Method: zip.Deflate,
	}
	w, err := writer.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create checksums file in archive: %w", err)
	}

	if _, err := io.Copy(w, checksumFile); err != nil {
		return fmt.Errorf("failed to write checksums to archive: %w", err)
	}

	return nil
}

// ReadChecksums reads checksums from an archive
func ReadChecksums(archive *Archive) (map[string]string, error) {
	// ARCH-002: Checksum reading from archive
	// DECISION-REF: DEC-001
	reader, err := zip.OpenReader(archive.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	checksumFile, err := findChecksumsFile(reader)
	if err != nil {
		return nil, err
	}

	return readChecksumsFromFile(checksumFile)
}

// findChecksumsFile finds the checksums file in the archive
func findChecksumsFile(reader *zip.ReadCloser) (*zip.File, error) {
	// Checksum file location in archive
	for _, file := range reader.File {
		if file.Name == ".checksums" {
			return file, nil
		}
	}
	return nil, fmt.Errorf("checksums file not found in archive")
}

// readChecksumsFromFile reads checksums from a file in the archive
func readChecksumsFromFile(file *zip.File) (map[string]string, error) {
	// Checksum data extraction from file
	rc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open checksums file: %w", err)
	}
	defer rc.Close()

	var checksums map[string]string
	decoder := json.NewDecoder(rc)
	if err := decoder.Decode(&checksums); err != nil {
		return nil, fmt.Errorf("failed to decode checksums: %w", err)
	}

	return checksums, nil
}

// VerifyChecksums verifies file checksums against stored values
func VerifyChecksums(archivePath string) (*VerificationStatus, error) {
	// ARCH-002: Complete checksum verification process
	// DECISION-REF: DEC-001
	status := &VerificationStatus{
		VerifiedAt: time.Now(),
		IsVerified: true,
	}

	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return handleVerificationError(status, "Failed to open archive: %v", err)
	}
	defer reader.Close()

	checksumFile, err := findChecksumsFile(reader)
	if err != nil {
		return handleVerificationError(status, "Checksums file not found in archive")
	}

	storedChecksums, err := readChecksumsFromFile(checksumFile)
	if err != nil {
		return handleVerificationError(status, "Failed to read checksums: %v", err)
	}

	if err := verifyArchiveChecksums(reader, storedChecksums, status); err != nil {
		return handleVerificationError(status, err.Error())
	}

	status.HasChecksums = true
	return status, nil
}

// handleVerificationError handles verification errors
func handleVerificationError(
	status *VerificationStatus,
	format string,
	args ...interface{},
) (*VerificationStatus, error) {
	// CFG-002: Verification error handling
	// DECISION-REF: DEC-004
	status.IsVerified = false
	status.Errors = append(status.Errors, fmt.Sprintf(format, args...))
	return status, nil
}

// verifyArchiveChecksums verifies checksums for all files in the archive
func verifyArchiveChecksums(
	reader *zip.ReadCloser,
	storedChecksums map[string]string,
	status *VerificationStatus,
) error {
	// Archive-wide checksum verification
	for _, file := range reader.File {
		if file.Name == ".checksums" {
			continue
		}

		if err := verifyFileChecksum(file, storedChecksums, status); err != nil {
			return err
		}
	}
	return nil
}

// verifyFileChecksum verifies the checksum of a single file
func verifyFileChecksum(file *zip.File, storedChecksums map[string]string, _ *VerificationStatus) error {
	// Individual file checksum verification
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", file.Name, err)
	}
	defer rc.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, rc); err != nil {
		return fmt.Errorf("failed to calculate checksum for %s: %v", file.Name, err)
	}

	calculatedChecksum := hex.EncodeToString(hash.Sum(nil))
	storedChecksum, exists := storedChecksums[file.Name]

	if !exists {
		return fmt.Errorf("no stored checksum for %s", file.Name)
	}
	if calculatedChecksum != storedChecksum {
		return fmt.Errorf("checksum mismatch for %s", file.Name)
	}

	return nil
}

// StoreVerificationStatus stores verification status in a metadata file
func StoreVerificationStatus(archive *Archive, status *VerificationStatus) error {
	// ARCH-002: Verification status persistence
	// DECISION-REF: DEC-008
	metadataDir := filepath.Join(filepath.Dir(archive.Path), ".metadata")
	if err := os.MkdirAll(metadataDir, 0o755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	metadataPath := filepath.Join(metadataDir, archive.Name+".json")
	file, err := os.Create(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(status); err != nil {
		return fmt.Errorf("failed to encode verification status: %w", err)
	}

	return nil
}

// LoadVerificationStatus loads verification status from a metadata file
func LoadVerificationStatus(archive *Archive) (*VerificationStatus, error) {
	// ARCH-002: Verification status loading
	metadataDir := filepath.Join(filepath.Dir(archive.Path), ".metadata")
	metadataPath := filepath.Join(metadataDir, archive.Name+".json")

	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return nil, nil
	}

	file, err := os.Open(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer file.Close()

	var status VerificationStatus
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode verification status: %w", err)
	}

	return &status, nil
}
