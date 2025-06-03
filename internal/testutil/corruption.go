// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically archive corruption testing utilities.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-A: Archive corruption testing framework - 📝
// DECISION-REF: DEC-001 (ZIP format)
// IMPLEMENTATION-NOTES: Provides systematic corruption for verification testing

package testutil

import (
	"archive/zip"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
)

// CorruptionType represents different types of archive corruption
type CorruptionType int

const (
	// CRC errors - corrupt file checksums
	CorruptionCRC CorruptionType = iota
	// Header corruption - corrupt ZIP file headers
	CorruptionHeader
	// File truncation - cut off end of archive
	CorruptionTruncate
	// Central directory corruption - corrupt ZIP central directory
	CorruptionCentralDir
	// Local header corruption - corrupt individual file headers
	CorruptionLocalHeader
	// Data corruption - corrupt actual file data
	CorruptionData
	// Signature corruption - corrupt ZIP file signatures
	CorruptionSignature
	// Comment corruption - corrupt archive comments
	CorruptionComment
)

// CorruptionConfig configures how corruption should be applied
type CorruptionConfig struct {
	Type           CorruptionType
	Seed           int64   // For reproducible corruption
	TargetFile     string  // Specific file to corrupt (empty for archive-level)
	CorruptionSize int     // Number of bytes to corrupt
	Offset         int     // Byte offset for corruption (-1 for random)
	Severity       float32 // 0.0-1.0, how severe the corruption should be
}

// CorruptionResult contains information about applied corruption
type CorruptionResult struct {
	Type           CorruptionType
	AppliedAt      []int  // Byte offsets where corruption was applied
	OriginalBytes  []byte // Original bytes for potential recovery
	CorruptedBytes []byte // The corrupted bytes that were written
	Description    string // Human-readable description of corruption
	Recoverable    bool   // Whether this corruption can be recovered from
}

// ArchiveCorruptor provides utilities for controlled ZIP corruption
type ArchiveCorruptor struct {
	originalPath string
	backupPath   string
	config       CorruptionConfig
}

// NewArchiveCorruptor creates a new archive corruptor for the given archive
func NewArchiveCorruptor(archivePath string, config CorruptionConfig) (*ArchiveCorruptor, error) {
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("archive file does not exist: %s", archivePath)
	}

	backupPath := archivePath + ".corruption-backup"

	return &ArchiveCorruptor{
		originalPath: archivePath,
		backupPath:   backupPath,
		config:       config,
	}, nil
}

// CreateBackup creates a backup of the original archive
func (ac *ArchiveCorruptor) CreateBackup() error {
	input, err := os.Open(ac.originalPath)
	if err != nil {
		return fmt.Errorf("failed to open original archive: %w", err)
	}
	defer input.Close()

	output, err := os.Create(ac.backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

// RestoreFromBackup restores the archive from backup
func (ac *ArchiveCorruptor) RestoreFromBackup() error {
	if _, err := os.Stat(ac.backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", ac.backupPath)
	}

	return os.Rename(ac.backupPath, ac.originalPath)
}

// Cleanup removes the backup file
func (ac *ArchiveCorruptor) Cleanup() error {
	if _, err := os.Stat(ac.backupPath); os.IsNotExist(err) {
		return nil // Nothing to clean up
	}
	return os.Remove(ac.backupPath)
}

// ApplyCorruption applies the configured corruption to the archive
func (ac *ArchiveCorruptor) ApplyCorruption() (*CorruptionResult, error) {
	// Create backup first
	if err := ac.CreateBackup(); err != nil {
		return nil, fmt.Errorf("failed to create backup: %w", err)
	}

	switch ac.config.Type {
	case CorruptionCRC:
		return ac.corruptCRC()
	case CorruptionHeader:
		return ac.corruptHeader()
	case CorruptionTruncate:
		return ac.corruptTruncate()
	case CorruptionCentralDir:
		return ac.corruptCentralDirectory()
	case CorruptionLocalHeader:
		return ac.corruptLocalHeader()
	case CorruptionData:
		return ac.corruptData()
	case CorruptionSignature:
		return ac.corruptSignature()
	case CorruptionComment:
		return ac.corruptComment()
	default:
		return nil, fmt.Errorf("unsupported corruption type: %d", ac.config.Type)
	}
}

// findLocalHeaderOffsets finds all local header offsets in a deterministic order
func (ac *ArchiveCorruptor) findLocalHeaderOffsets(data []byte) []int {
	var offsets []int
	for i := 0; i < len(data)-4; i++ {
		// Look for local file header signature (0x04034b50)
		if binary.LittleEndian.Uint32(data[i:i+4]) == 0x04034b50 {
			offsets = append(offsets, i)
		}
	}
	// Sort to ensure deterministic order
	sort.Ints(offsets)
	return offsets
}

// corruptCRC corrupts file CRC checksums
func (ac *ArchiveCorruptor) corruptCRC() (*CorruptionResult, error) {
	// Read entire file into memory for manipulation
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionCRC,
		Description: "Corrupted CRC checksums in ZIP headers",
		Recoverable: true, // CRC errors are often recoverable
	}

	// Find all local header offsets in deterministic order
	headerOffsets := ac.findLocalHeaderOffsets(data)

	corruptedCount := 0
	for _, headerOffset := range headerOffsets {
		crcOffset := headerOffset + 14 // CRC is at offset 14 from header start
		if crcOffset+4 <= len(data) {
			// Store original bytes
			originalCRC := make([]byte, 4)
			copy(originalCRC, data[crcOffset:crcOffset+4])
			result.OriginalBytes = append(result.OriginalBytes, originalCRC...)

			// Generate corrupted CRC using seed for reproducibility
			corruptedCRC := make([]byte, 4)
			if ac.config.Seed != 0 {
				// Use seed + offset for deterministic but different corruption per location
				seedValue := uint32(ac.config.Seed + int64(headerOffset))
				binary.LittleEndian.PutUint32(corruptedCRC, seedValue)
			} else {
				rand.Read(corruptedCRC)
			}

			// Apply corruption
			copy(data[crcOffset:crcOffset+4], corruptedCRC)
			result.AppliedAt = append(result.AppliedAt, crcOffset)
			result.CorruptedBytes = append(result.CorruptedBytes, corruptedCRC...)
			corruptedCount++

			if ac.config.CorruptionSize > 0 && corruptedCount >= ac.config.CorruptionSize {
				break
			}
		}
	}

	if corruptedCount == 0 {
		return nil, fmt.Errorf("no CRC locations found to corrupt")
	}

	result.Description = fmt.Sprintf("Corrupted %d CRC checksums", corruptedCount)

	// Write corrupted data back to file
	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptHeader corrupts ZIP file headers - more aggressively to ensure unreadability
func (ac *ArchiveCorruptor) corruptHeader() (*CorruptionResult, error) {
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionHeader,
		Description: "Corrupted ZIP file headers",
		Recoverable: false, // Header corruption is usually fatal
	}

	// Find first local file header and corrupt multiple fields to ensure unreadability
	headerOffsets := ac.findLocalHeaderOffsets(data)
	if len(headerOffsets) == 0 {
		return nil, fmt.Errorf("no headers found to corrupt")
	}

	headerOffset := headerOffsets[0]

	// Corrupt signature (first 4 bytes)
	if headerOffset+4 <= len(data) {
		result.OriginalBytes = make([]byte, 4)
		copy(result.OriginalBytes, data[headerOffset:headerOffset+4])

		// Use recognizably corrupted signature
		corruptedSig := []byte{0xDE, 0xAD, 0xBE, 0xEF}
		if ac.config.Seed != 0 {
			binary.LittleEndian.PutUint32(corruptedSig, uint32(ac.config.Seed))
		}

		copy(data[headerOffset:headerOffset+4], corruptedSig)
		result.AppliedAt = append(result.AppliedAt, headerOffset)
		result.CorruptedBytes = append(result.CorruptedBytes, corruptedSig...)
	}

	// Also corrupt version field for good measure
	versionOffset := headerOffset + 4
	if versionOffset+2 <= len(data) {
		originalVersion := make([]byte, 2)
		copy(originalVersion, data[versionOffset:versionOffset+2])
		result.OriginalBytes = append(result.OriginalBytes, originalVersion...)

		corruptedVersion := []byte{0xFF, 0xFF}
		copy(data[versionOffset:versionOffset+2], corruptedVersion)
		result.AppliedAt = append(result.AppliedAt, versionOffset)
		result.CorruptedBytes = append(result.CorruptedBytes, corruptedVersion...)
	}

	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptTruncate truncates the archive file
func (ac *ArchiveCorruptor) corruptTruncate() (*CorruptionResult, error) {
	info, err := os.Stat(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat archive: %w", err)
	}

	originalSize := info.Size()
	truncateSize := originalSize

	if ac.config.CorruptionSize > 0 {
		truncateSize = originalSize - int64(ac.config.CorruptionSize)
	} else {
		// Default: truncate 10% of the file
		truncateSize = int64(float64(originalSize) * (1.0 - float64(ac.config.Severity)))
	}

	if truncateSize < 0 {
		truncateSize = 0
	}

	err = os.Truncate(ac.originalPath, truncateSize)
	if err != nil {
		return nil, fmt.Errorf("failed to truncate archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionTruncate,
		AppliedAt:   []int{int(truncateSize)},
		Description: fmt.Sprintf("Truncated archive from %d to %d bytes", originalSize, truncateSize),
		Recoverable: false, // Truncation loses data permanently
	}

	return result, nil
}

// corruptCentralDirectory corrupts the ZIP central directory
func (ac *ArchiveCorruptor) corruptCentralDirectory() (*CorruptionResult, error) {
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionCentralDir,
		Description: "Corrupted ZIP central directory",
		Recoverable: false, // Central directory corruption is usually fatal
	}

	// Find End of Central Directory signature (0x06054b50)
	for i := len(data) - 4; i >= 0; i-- {
		if binary.LittleEndian.Uint32(data[i:i+4]) == 0x06054b50 {
			// Found EOCD, now corrupt central directory offset
			if i+16 < len(data) {
				cdOffset := i + 16 // Central directory offset is at this position

				// Store original bytes
				result.OriginalBytes = make([]byte, 4)
				copy(result.OriginalBytes, data[cdOffset:cdOffset+4])

				// Corrupt the offset
				corruptedOffset := make([]byte, 4)
				if ac.config.Seed != 0 {
					binary.LittleEndian.PutUint32(corruptedOffset, uint32(ac.config.Seed))
				} else {
					binary.LittleEndian.PutUint32(corruptedOffset, 0xFFFFFFFF)
				}

				copy(data[cdOffset:cdOffset+4], corruptedOffset)
				result.AppliedAt = append(result.AppliedAt, cdOffset)
				result.CorruptedBytes = corruptedOffset
				break
			}
		}
	}

	if len(result.AppliedAt) == 0 {
		return nil, fmt.Errorf("no central directory found to corrupt")
	}

	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptLocalHeader corrupts local file headers
func (ac *ArchiveCorruptor) corruptLocalHeader() (*CorruptionResult, error) {
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionLocalHeader,
		Description: "Corrupted local file headers",
		Recoverable: true, // Sometimes recoverable with repair tools
	}

	// Find local file headers and corrupt version field
	headerOffsets := ac.findLocalHeaderOffsets(data)

	corruptedCount := 0
	for _, headerOffset := range headerOffsets {
		versionOffset := headerOffset + 4 // Version field is right after signature

		if versionOffset+2 <= len(data) {
			// Store original bytes
			original := make([]byte, 2)
			copy(original, data[versionOffset:versionOffset+2])
			result.OriginalBytes = append(result.OriginalBytes, original...)

			// Corrupt version field
			corruptedVersion := []byte{0xFF, 0xFF}
			if ac.config.Seed != 0 {
				binary.LittleEndian.PutUint16(corruptedVersion, uint16(ac.config.Seed+int64(headerOffset)))
			}

			copy(data[versionOffset:versionOffset+2], corruptedVersion)
			result.AppliedAt = append(result.AppliedAt, versionOffset)
			result.CorruptedBytes = append(result.CorruptedBytes, corruptedVersion...)
			corruptedCount++

			if ac.config.CorruptionSize > 0 && corruptedCount >= ac.config.CorruptionSize {
				break
			}
		}
	}

	if corruptedCount == 0 {
		return nil, fmt.Errorf("no local headers found to corrupt")
	}

	result.Description = fmt.Sprintf("Corrupted %d local file headers", corruptedCount)
	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptData corrupts actual file data within the archive
func (ac *ArchiveCorruptor) corruptData() (*CorruptionResult, error) {
	// Open archive to find file data locations
	reader, err := zip.OpenReader(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionData,
		Description: "Corrupted file data within archive",
		Recoverable: true, // Data corruption might be detectable/recoverable
	}

	// Find file data and corrupt it
	if len(reader.File) > 0 {
		file := reader.File[0] // Corrupt first file's data

		// Calculate approximate data location
		// This is a simplified approach - in reality, finding exact data location
		// in ZIP requires parsing headers carefully
		headerSize := 30 + len(file.Name) // Approximate local header size
		dataStart := headerSize
		dataSize := int(file.CompressedSize64)

		if dataStart+dataSize < len(data) {
			corruptionSize := ac.config.CorruptionSize
			if corruptionSize == 0 {
				corruptionSize = min(dataSize/10, 100) // Corrupt up to 10% or 100 bytes
			}

			// Store original bytes
			result.OriginalBytes = make([]byte, corruptionSize)
			copy(result.OriginalBytes, data[dataStart:dataStart+corruptionSize])

			// Generate corruption
			corruptedBytes := make([]byte, corruptionSize)
			if ac.config.Seed != 0 {
				// Deterministic corruption based on seed
				for i := range corruptedBytes {
					corruptedBytes[i] = byte((ac.config.Seed + int64(i)) % 256)
				}
			} else {
				rand.Read(corruptedBytes)
			}

			copy(data[dataStart:dataStart+corruptionSize], corruptedBytes)
			result.AppliedAt = append(result.AppliedAt, dataStart)
			result.CorruptedBytes = corruptedBytes

			result.Description = fmt.Sprintf("Corrupted %d bytes of file data", corruptionSize)
		}
	}

	if len(result.AppliedAt) == 0 {
		return nil, fmt.Errorf("no file data found to corrupt")
	}

	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptSignature corrupts ZIP file signatures more thoroughly
func (ac *ArchiveCorruptor) corruptSignature() (*CorruptionResult, error) {
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionSignature,
		Description: "Corrupted ZIP file signatures",
		Recoverable: false, // Signature corruption usually makes files unreadable
	}

	// Corrupt the very first bytes completely (should be local file header signature)
	if len(data) >= 4 {
		result.OriginalBytes = make([]byte, 4)
		copy(result.OriginalBytes, data[0:4])

		// Use completely invalid signature that will definitely break ZIP reading
		corruptedSig := []byte{0x00, 0x00, 0x00, 0x00}
		if ac.config.Seed != 0 {
			binary.LittleEndian.PutUint32(corruptedSig, uint32(ac.config.Seed))
		}

		copy(data[0:4], corruptedSig)
		result.AppliedAt = append(result.AppliedAt, 0)
		result.CorruptedBytes = corruptedSig
	}

	if len(result.AppliedAt) == 0 {
		return nil, fmt.Errorf("archive too small to corrupt signature")
	}

	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// corruptComment corrupts ZIP archive comments
func (ac *ArchiveCorruptor) corruptComment() (*CorruptionResult, error) {
	data, err := os.ReadFile(ac.originalPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive: %w", err)
	}

	result := &CorruptionResult{
		Type:        CorruptionComment,
		Description: "Corrupted ZIP archive comment",
		Recoverable: true, // Comment corruption is usually non-fatal
	}

	// Find End of Central Directory and add/corrupt comment
	for i := len(data) - 22; i >= 0; i-- { // EOCD is at least 22 bytes
		if binary.LittleEndian.Uint32(data[i:i+4]) == 0x06054b50 {
			commentLenOffset := i + 20
			if commentLenOffset+2 <= len(data) {
				// Set comment length to indicate corrupted comment
				commentLen := uint16(10) // 10 bytes of corrupted comment
				binary.LittleEndian.PutUint16(data[commentLenOffset:commentLenOffset+2], commentLen)

				// Append corrupted comment data
				corruptedComment := []byte("CORRUPTED!")
				if ac.config.Seed != 0 {
					corruptedComment = []byte(fmt.Sprintf("SEED:%d", ac.config.Seed))
				}

				// Extend the file with the comment
				data = append(data, corruptedComment...)
				result.AppliedAt = append(result.AppliedAt, len(data)-len(corruptedComment))
				result.CorruptedBytes = corruptedComment
				break
			}
		}
	}

	if len(result.AppliedAt) == 0 {
		return nil, fmt.Errorf("no suitable location found for comment corruption")
	}

	return result, os.WriteFile(ac.originalPath, data, 0644)
}

// Helper function for min since it's not available in older Go versions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CorruptionDetector provides utilities for detecting archive corruption
type CorruptionDetector struct {
	archivePath string
}

// NewCorruptionDetector creates a new corruption detector
func NewCorruptionDetector(archivePath string) *CorruptionDetector {
	return &CorruptionDetector{archivePath: archivePath}
}

// DetectCorruption attempts to detect various types of corruption
func (cd *CorruptionDetector) DetectCorruption() ([]CorruptionType, error) {
	var detectedTypes []CorruptionType

	// Check for basic corruption first
	if isTruncated(cd.archivePath) {
		detectedTypes = append(detectedTypes, CorruptionTruncate)
	}

	if isSignatureCorruption(cd.archivePath) {
		detectedTypes = append(detectedTypes, CorruptionSignature)
	}

	// Try to open as ZIP archive
	reader, err := zip.OpenReader(cd.archivePath)
	if err != nil {
		// If we can't open it but haven't detected signature/truncation, it's likely header corruption
		if len(detectedTypes) == 0 {
			detectedTypes = append(detectedTypes, CorruptionHeader)
		}
		return detectedTypes, nil
	}
	defer reader.Close()

	// Archive opened successfully, check for other types of corruption
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			// Could be various types of corruption
			detectedTypes = append(detectedTypes, CorruptionCRC)
			continue
		}

		// Try to read the file data
		_, err = io.Copy(io.Discard, rc)
		rc.Close()

		if err != nil {
			detectedTypes = append(detectedTypes, CorruptionData)
		}
	}

	return detectedTypes, nil
}

// isSignatureCorruption checks if the file has signature corruption
func isSignatureCorruption(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil || len(data) < 4 {
		return false
	}

	// Check if first 4 bytes are not a valid ZIP signature
	sig := binary.LittleEndian.Uint32(data[0:4])
	validSignatures := []uint32{
		0x04034b50, // Local file header
		0x08074b50, // Data descriptor
		0x02014b50, // Central directory header
		0x06054b50, // End of central directory
	}

	for _, validSig := range validSignatures {
		if sig == validSig {
			return false
		}
	}

	return true
}

// isHeaderCorruption checks for header corruption beyond signature
func isHeaderCorruption(path string) bool {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return true // If we can't open it, assume header corruption
	}
	defer reader.Close()

	// If we can open it but can't read files, might be header corruption
	if len(reader.File) == 0 {
		return true
	}

	return false
}

// isTruncated checks if the file appears to be truncated
func isTruncated(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	if len(data) < 22 { // Minimum size for ZIP file with EOCD
		return true
	}

	// Look for End of Central Directory signature at the end
	found := false
	for i := len(data) - 22; i >= max(0, len(data)-65557); i-- { // Max comment size is 65535
		if len(data[i:]) >= 4 && binary.LittleEndian.Uint32(data[i:i+4]) == 0x06054b50 {
			found = true
			break
		}
	}

	return !found
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// CreateTestArchive creates a test archive with specified files for corruption testing
func CreateTestArchive(archivePath string, files map[string]string) error {
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for filePath, content := range files {
		fileWriter, err := zipWriter.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s in archive: %w", filePath, err)
		}

		_, err = fileWriter.Write([]byte(content))
		if err != nil {
			return fmt.Errorf("failed to write content to file %s: %w", filePath, err)
		}
	}

	return nil
}

// CreateCorruptedTestArchive creates a test archive and immediately applies corruption
func CreateCorruptedTestArchive(archivePath string, files map[string]string, config CorruptionConfig) (*CorruptionResult, error) {
	// First create the archive
	if err := CreateTestArchive(archivePath, files); err != nil {
		return nil, fmt.Errorf("failed to create test archive: %w", err)
	}

	// Then corrupt it
	corruptor, err := NewArchiveCorruptor(archivePath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create corruptor: %w", err)
	}
	defer corruptor.Cleanup()

	return corruptor.ApplyCorruption()
}
