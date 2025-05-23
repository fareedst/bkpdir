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
		rc, err := file.Open()
		if err != nil {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to open file %s: %v", file.Name, err))
			continue
		}

		// Read a small portion to verify the file can be read
		buf := make([]byte, 1024)
		_, err = rc.Read(buf)
		if err != nil && err != io.EOF {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to read file %s: %v", file.Name, err))
		}
		rc.Close()
	}

	return status, nil
}

// GenerateChecksums generates checksums for a list of files
func GenerateChecksums(files []string, algorithm string) (map[string]string, error) {
	checksums := make(map[string]string)

	for _, file := range files {
		info, err := os.Lstat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", file, err)
		}

		// Skip directories
		if info.IsDir() {
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", file, err)
		}
		defer f.Close()

		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			return nil, fmt.Errorf("failed to calculate checksum for %s: %w", file, err)
		}

		// Use the base name of the file as the key
		checksums[filepath.Base(file)] = hex.EncodeToString(hash.Sum(nil))
	}

	return checksums, nil
}

// StoreChecksums stores checksums in the archive
func StoreChecksums(archive *Archive, checksums map[string]string) error {
	// Create a temporary file for checksums
	tmpFile, err := os.CreateTemp("", "bkpdir-checksums-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write checksums to the temporary file
	encoder := json.NewEncoder(tmpFile)
	if err := encoder.Encode(checksums); err != nil {
		return fmt.Errorf("failed to encode checksums: %w", err)
	}
	tmpFile.Close()

	// Open the archive for writing
	reader, err := zip.OpenReader(archive.Path)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	// Create a new archive with the checksums file
	newPath := archive.Path + ".new"
	writer, err := os.Create(newPath)
	if err != nil {
		return fmt.Errorf("failed to create new archive: %w", err)
	}
	defer os.Remove(newPath)

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	// Copy all files from the original archive
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		header, err := zip.FileInfoHeader(file.FileInfo())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create header: %w", err)
		}
		header.Name = file.Name

		w, err := zipWriter.CreateHeader(header)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file in new archive: %w", err)
		}

		_, err = io.Copy(w, rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file to new archive: %w", err)
		}
	}

	// Add the checksums file
	checksumFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open checksums file: %w", err)
	}
	defer checksumFile.Close()

	header := &zip.FileHeader{
		Name:   ".checksums",
		Method: zip.Deflate,
	}
	w, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create checksums file in archive: %w", err)
	}

	_, err = io.Copy(w, checksumFile)
	if err != nil {
		return fmt.Errorf("failed to write checksums to archive: %w", err)
	}

	zipWriter.Close()
	writer.Close()

	// Replace the original archive with the new one
	if err := os.Rename(newPath, archive.Path); err != nil {
		return fmt.Errorf("failed to replace original archive: %w", err)
	}

	return nil
}

// ReadChecksums reads checksums from an archive
func ReadChecksums(archive *Archive) (map[string]string, error) {
	reader, err := zip.OpenReader(archive.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	// Find the checksums file
	var checksumFile *zip.File
	for _, file := range reader.File {
		if file.Name == ".checksums" {
			checksumFile = file
			break
		}
	}

	if checksumFile == nil {
		return nil, fmt.Errorf("checksums file not found in archive")
	}

	// Read the checksums file
	rc, err := checksumFile.Open()
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

	// Find the checksums file
	var checksumFile *zip.File
	for _, file := range reader.File {
		if file.Name == ".checksums" {
			checksumFile = file
			break
		}
	}

	if checksumFile == nil {
		status.IsVerified = false
		status.Errors = append(status.Errors, "Checksums file not found in archive")
		return status, nil
	}

	// Read the checksums
	rc, err := checksumFile.Open()
	if err != nil {
		status.IsVerified = false
		status.Errors = append(status.Errors, fmt.Sprintf("Failed to open checksums file: %v", err))
		return status, nil
	}
	defer rc.Close()

	var storedChecksums map[string]string
	decoder := json.NewDecoder(rc)
	if err := decoder.Decode(&storedChecksums); err != nil {
		status.IsVerified = false
		status.Errors = append(status.Errors, fmt.Sprintf("Failed to decode checksums: %v", err))
		return status, nil
	}

	// Verify each file's checksum
	for _, file := range reader.File {
		if file.Name == ".checksums" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to open file %s: %v", file.Name, err))
			continue
		}

		hash := sha256.New()
		if _, err := io.Copy(hash, rc); err != nil {
			rc.Close()
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("Failed to calculate checksum for %s: %v", file.Name, err))
			continue
		}
		rc.Close()

		calculatedChecksum := hex.EncodeToString(hash.Sum(nil))
		storedChecksum, exists := storedChecksums[file.Name]

		if !exists {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("No stored checksum for %s", file.Name))
		} else if calculatedChecksum != storedChecksum {
			status.IsVerified = false
			status.Errors = append(status.Errors, fmt.Sprintf("Checksum mismatch for %s", file.Name))
		}
	}

	status.HasChecksums = true
	return status, nil
}

// StoreVerificationStatus stores verification status in a metadata file
func StoreVerificationStatus(archive *Archive, status *VerificationStatus) error {
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
