// Package fileops provides file operations and utilities for CLI applications.
//
// This package extracts common file operation patterns including comparison,
// exclusion, validation, and directory traversal functionality.
package fileops

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ⭐ EXTRACT-006: File comparison interface and implementation extracted - 🔧

// FileInfo represents information about a file for comparison
type FileInfo struct {
	RelativePath string
	Size         int64
	ModTime      time.Time
	IsDir        bool
	Hash         string // SHA-256 hash for content comparison
}

// DirectorySnapshot represents a snapshot of a directory's contents
type DirectorySnapshot struct {
	Files []FileInfo
}

// Comparer defines the interface for file and directory comparison operations
type Comparer interface {
	CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error)
	CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error)
	CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool
	IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error)
}

// DefaultComparer implements the Comparer interface with standard comparison logic
type DefaultComparer struct{}

// NewComparer creates a new DefaultComparer instance
func NewComparer() Comparer {
	// ⭐ EXTRACT-006: Comparer interface implementation - 🔧
	return &DefaultComparer{}
}

// CreateDirectorySnapshot creates a snapshot of the given directory
func (c *DefaultComparer) CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Directory snapshot creation extracted - 🔧
	var files []FileInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Check exclusion patterns using the exclusion system
		if ShouldExcludeFile(relPath, excludePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		fileInfo := FileInfo{
			RelativePath: relPath,
			Size:         info.Size(),
			ModTime:      info.ModTime(),
			IsDir:        info.IsDir(),
		}

		// Calculate hash for regular files
		if !info.IsDir() {
			hash, err := c.calculateFileHash(path)
			if err != nil {
				return err
			}
			fileInfo.Hash = hash
		}

		files = append(files, fileInfo)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files for consistent comparison
	sort.Slice(files, func(i, j int) bool {
		return files[i].RelativePath < files[j].RelativePath
	})

	return &DirectorySnapshot{Files: files}, nil
}

// CreateArchiveSnapshot creates a snapshot from a ZIP archive
func (c *DefaultComparer) CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Archive snapshot creation extracted - 🔧
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var files []FileInfo

	for _, file := range reader.File {
		// Skip directories in ZIP files (they're implicit)
		if strings.HasSuffix(file.Name, "/") {
			continue
		}

		fileInfo := FileInfo{
			RelativePath: file.Name,
			Size:         int64(file.UncompressedSize64),
			ModTime:      file.Modified,
			IsDir:        false,
		}

		// Calculate hash of file content in archive
		hash, err := c.calculateArchiveFileHash(file)
		if err != nil {
			return nil, err
		}
		fileInfo.Hash = hash

		files = append(files, fileInfo)
	}

	// Sort files for consistent comparison
	sort.Slice(files, func(i, j int) bool {
		return files[i].RelativePath < files[j].RelativePath
	})

	return &DirectorySnapshot{Files: files}, nil
}

// CompareSnapshots compares two directory snapshots
func (c *DefaultComparer) CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool {
	// ⭐ EXTRACT-006: Snapshot comparison implementation extracted - 🔧
	if len(snapshot1.Files) != len(snapshot2.Files) {
		return false
	}

	for i, file1 := range snapshot1.Files {
		file2 := snapshot2.Files[i]

		// Compare basic properties
		if file1.RelativePath != file2.RelativePath ||
			file1.Size != file2.Size ||
			file1.IsDir != file2.IsDir {
			return false
		}

		// Compare content hash for regular files
		if !file1.IsDir && file1.Hash != file2.Hash {
			return false
		}
	}

	return true
}

// IsDirectoryIdenticalToArchive checks if a directory is identical to an archive
func (c *DefaultComparer) IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error) {
	// ⭐ EXTRACT-006: Directory-to-archive comparison extracted - 🔍
	// Create snapshot of directory
	dirSnapshot, err := c.CreateDirectorySnapshot(dirPath, excludePatterns)
	if err != nil {
		return false, err
	}

	// Create snapshot of archive
	archiveSnapshot, err := c.CreateArchiveSnapshot(archivePath)
	if err != nil {
		return false, err
	}

	// Compare snapshots
	return c.CompareSnapshots(dirSnapshot, archiveSnapshot), nil
}

// calculateFileHash calculates SHA-256 hash of a file
func (c *DefaultComparer) calculateFileHash(filePath string) (string, error) {
	// ⭐ EXTRACT-006: File hash calculation extracted - 🔧
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// calculateArchiveFileHash calculates SHA-256 hash of a file in an archive
func (c *DefaultComparer) calculateArchiveFileHash(file *zip.File) (string, error) {
	// ⭐ EXTRACT-006: Archive file hash calculation extracted - 🔧
	reader, err := file.Open()
	if err != nil {
		return "", err
	}
	defer reader.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Convenience functions for backward compatibility

// CreateDirectorySnapshot creates a snapshot using the default comparer
func CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Backward compatibility function - 🔧
	comparer := NewComparer()
	return comparer.CreateDirectorySnapshot(rootPath, excludePatterns)
}

// CreateArchiveSnapshot creates a snapshot using the default comparer
func CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Backward compatibility function - 🔧
	comparer := NewComparer()
	return comparer.CreateArchiveSnapshot(archivePath)
}

// CompareSnapshots compares snapshots using the default comparer
func CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool {
	// ⭐ EXTRACT-006: Backward compatibility function - 🔧
	comparer := NewComparer()
	return comparer.CompareSnapshots(snapshot1, snapshot2)
}

// IsDirectoryIdenticalToArchive checks identity using the default comparer
func IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error) {
	// ⭐ EXTRACT-006: Backward compatibility function - 🔧
	comparer := NewComparer()
	return comparer.IsDirectoryIdenticalToArchive(dirPath, archivePath, excludePatterns)
}
