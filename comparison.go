// This file is part of bkpdir
//
// Package main provides directory comparison functionality for the BkpDir application.
// It handles comparing directories to detect changes and identical states.
package main

import (
	"fmt"
	"os"
	"time"

	"bkpdir/pkg/fileops"
)

// ⭐ EXTRACT-006: Updated to use extracted fileops package - 🔧
// 🔶 REFACTOR-001: File comparison interface contracts defined - 📝
// 🔶 REFACTOR-001: Dependency analysis - clean boundary confirmed - 📝

// Legacy type aliases for backward compatibility
type (
	FileInfo          = fileops.FileInfo
	DirectorySnapshot = fileops.DirectorySnapshot
)

// Legacy function wrappers for backward compatibility

// CreateDirectorySnapshot creates a snapshot of the given directory using the extracted package
func CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔧
	return fileops.CreateDirectorySnapshot(rootPath, excludePatterns)
}

// CreateArchiveSnapshot creates a snapshot from a ZIP archive using the extracted package
func CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔧
	return fileops.CreateArchiveSnapshot(archivePath)
}

// CompareSnapshots compares two directory snapshots using the extracted package
func CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔧
	return fileops.CompareSnapshots(snapshot1, snapshot2)
}

// IsDirectoryIdenticalToArchive checks if a directory is identical to an archive using the extracted package
func IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error) {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔧
	return fileops.IsDirectoryIdenticalToArchive(dirPath, archivePath, excludePatterns)
}

// FindMostRecentArchive finds the most recent archive in the archive directory
func FindMostRecentArchive(archiveDir string) (string, error) {
	// ⭐ ARCH-003: Most recent archive detection for incremental comparison - 🔍
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return "", err
	}

	if len(archives) == 0 {
		return "", nil
	}

	// Find the most recent full archive (skip incremental archives)
	var mostRecent Archive
	var mostRecentTime time.Time

	for _, archive := range archives {
		// Skip incremental archives for this comparison
		if archive.IsIncremental {
			continue
		}

		// Get file modification time
		info, err := os.Stat(archive.Path)
		if err != nil {
			continue
		}

		if info.ModTime().After(mostRecentTime) {
			mostRecentTime = info.ModTime()
			mostRecent = archive
		}
	}

	if mostRecent.Path == "" {
		return "", nil
	}

	return mostRecent.Path, nil
}

// CheckForIdenticalArchive checks if the directory is identical to the most recent archive
func CheckForIdenticalArchive(dirPath, archiveDir string, excludePatterns []string) (bool, string, error) {
	// ⭐ ARCH-003: Identical archive detection for incremental backup - 🔍
	// Find most recent archive
	mostRecentArchive, err := FindMostRecentArchive(archiveDir)
	if err != nil {
		return false, "", err
	}

	if mostRecentArchive == "" {
		// No archives exist
		return false, "", nil
	}

	// Check if directory is identical to most recent archive
	identical, err := IsDirectoryIdenticalToArchive(dirPath, mostRecentArchive, excludePatterns)
	if err != nil {
		return false, "", err
	}

	return identical, mostRecentArchive, nil
}

// GetDirectoryTreeSummary returns a summary of directory structure and content
func GetDirectoryTreeSummary(dirPath string, excludePatterns []string) (string, error) {
	// ⭐ EXTRACT-006: Using extracted directory traversal - 🔍
	snapshot, err := CreateDirectorySnapshot(dirPath, excludePatterns)
	if err != nil {
		return "", err
	}

	summary := fmt.Sprintf("Directory: %s\nFiles: %d\n", dirPath, len(snapshot.Files))
	for _, file := range snapshot.Files {
		if file.IsDir {
			summary += fmt.Sprintf("  [DIR]  %s\n", file.RelativePath)
		} else {
			summary += fmt.Sprintf("  [FILE] %s (%d bytes)\n", file.RelativePath, file.Size)
		}
	}

	return summary, nil
}

// GetArchiveTreeSummary returns a summary of archive structure and content
func GetArchiveTreeSummary(archivePath string) (string, error) {
	// ⭐ EXTRACT-006: Using extracted archive comparison - 🔍
	snapshot, err := CreateArchiveSnapshot(archivePath)
	if err != nil {
		return "", err
	}

	summary := fmt.Sprintf("Archive: %s\nFiles: %d\n", archivePath, len(snapshot.Files))
	for _, file := range snapshot.Files {
		summary += fmt.Sprintf("  [FILE] %s (%d bytes)\n", file.RelativePath, file.Size)
	}

	return summary, nil
}
