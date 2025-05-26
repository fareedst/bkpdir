package main

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

// CreateDirectorySnapshot creates a snapshot of the given directory
func CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
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

		// Check exclusion patterns
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
			hash, err := calculateFileHash(path)
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
func CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
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
		hash, err := calculateArchiveFileHash(file)
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
func CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool {
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
func IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error) {
	// Create snapshot of directory
	dirSnapshot, err := CreateDirectorySnapshot(dirPath, excludePatterns)
	if err != nil {
		return false, err
	}

	// Create snapshot of archive
	archiveSnapshot, err := CreateArchiveSnapshot(archivePath)
	if err != nil {
		return false, err
	}

	// Compare snapshots
	return CompareSnapshots(dirSnapshot, archiveSnapshot), nil
}

// FindMostRecentArchive finds the most recent archive in the archive directory
func FindMostRecentArchive(archiveDir string) (string, error) {
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return "", err
	}

	if len(archives) == 0 {
		return "", nil
	}

	// Find the most recent archive (archives are already sorted by creation time)
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

// Helper function to calculate SHA-256 hash of a file
func calculateFileHash(filePath string) (string, error) {
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

// Helper function to calculate SHA-256 hash of a file in a ZIP archive
func calculateArchiveFileHash(file *zip.File) (string, error) {
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

// GetDirectoryTreeSummary returns a summary of the directory tree for debugging
func GetDirectoryTreeSummary(dirPath string, excludePatterns []string) (string, error) {
	snapshot, err := CreateDirectorySnapshot(dirPath, excludePatterns)
	if err != nil {
		return "", err
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Directory: %s\n", dirPath))
	summary.WriteString(fmt.Sprintf("Files: %d\n", len(snapshot.Files)))

	for _, file := range snapshot.Files {
		if file.IsDir {
			summary.WriteString(fmt.Sprintf("  [DIR]  %s\n", file.RelativePath))
		} else {
			summary.WriteString(fmt.Sprintf("  [FILE] %s (%d bytes, %s)\n",
				file.RelativePath, file.Size, file.Hash[:8]))
		}
	}

	return summary.String(), nil
}

// GetArchiveTreeSummary returns a summary of the archive contents for debugging
func GetArchiveTreeSummary(archivePath string) (string, error) {
	snapshot, err := CreateArchiveSnapshot(archivePath)
	if err != nil {
		return "", err
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Archive: %s\n", archivePath))
	summary.WriteString(fmt.Sprintf("Files: %d\n", len(snapshot.Files)))

	for _, file := range snapshot.Files {
		summary.WriteString(fmt.Sprintf("  [FILE] %s (%d bytes, %s)\n",
			file.RelativePath, file.Size, file.Hash[:8]))
	}

	return summary.String(), nil
}
