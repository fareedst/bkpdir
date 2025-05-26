package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Archive struct {
	Name               string
	Path               string
	CreationTime       time.Time
	IsIncremental      bool
	GitBranch          string
	GitHash            string
	Note               string
	BaseArchive        string // for incremental
	VerificationStatus *VerificationStatus
}

// GenerateArchiveName creates an archive name according to the spec.
func GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, isGit, isIncremental bool, baseName string) string {
	name := ""
	if prefix != "" {
		name += prefix + "-"
	}
	name += timestamp
	if isGit && gitBranch != "" && gitHash != "" {
		name += "=" + gitBranch + "=" + gitHash
	}
	if isIncremental && baseName != "" {
		name = strings.TrimSuffix(baseName, ".zip") + "_update=" + timestamp
		if isGit && gitBranch != "" && gitHash != "" {
			name += "=" + gitBranch + "=" + gitHash
		}
	}
	if note != "" {
		name += "=" + note
	}
	name += ".zip"
	return name
}

// ListArchives lists all archives in the archive directory for the current source.
func ListArchives(archiveDir string) ([]Archive, error) {
	// Create archive directory if it doesn't exist //
	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create archive directory: %w", err)
	}

	var archives []Archive
	dirEntries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".zip") {
			continue
		}

		// Get file info to extract creation time
		archivePath := filepath.Join(archiveDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			// If we can't get file info, use zero time
			fileInfo = nil
		}

		archive := Archive{
			Name:          entry.Name(),
			Path:          archivePath,
			IsIncremental: strings.Contains(entry.Name(), "_update="),
		}

		// Set creation time from file modification time
		if fileInfo != nil {
			archive.CreationTime = fileInfo.ModTime()
		}

		// Load verification status if available
		status, err := LoadVerificationStatus(&archive)
		if err == nil && status != nil {
			archive.VerificationStatus = status
		}

		archives = append(archives, archive)
	}
	return archives, nil
}

// CreateFullArchive creates a full archive
func CreateFullArchive(cfg *Config, note string, dryRun bool, verify bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	// Validate directory
	if err := ValidateDirectoryPath(cwd, cfg); err != nil {
		return err
	}

	// Create resource manager for cleanup
	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Determine archive directory
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		if err := SafeMkdirAll(archiveDir, 0755, cfg); err != nil {
			return err
		}
	}

	// Walk directory and collect files
	var files []string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if ShouldExcludeFile(rel, cfg.ExcludePatterns) {
			return nil
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return NewArchiveErrorWithCause("Failed to scan directory", 1, err)
	}

	// Git info
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit && cfg.IncludeGitInfo {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	// Archive name
	timestamp := time.Now().Format("2006-01-02-15-04")
	prefix := ""
	if cfg.UseCurrentDirName {
		prefix = filepath.Base(cwd)
	}
	archiveName := GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note, isGit && cfg.IncludeGitInfo, false, "")
	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		formatter := NewOutputFormatter(cfg)
		fmt.Println("[Dry Run] Files to include:")
		for _, f := range files {
			fmt.Println("  ", f)
		}
		formatter.PrintDryRunArchive(archivePath)
		return nil
	}

	// Create temporary file for atomic operation
	tempPath := archivePath + ".tmp"
	rm.AddTempFile(tempPath)

	// Create zip archive
	f, err := os.Create(tempPath)
	if err != nil {
		if IsDiskFullError(err) {
			return NewArchiveError("Insufficient disk space to create archive", cfg.StatusDiskFull)
		}
		return NewArchiveErrorWithCause("Failed to create archive file", 1, err)
	}
	defer f.Close()
	zipw := zip.NewWriter(f)
	defer zipw.Close()

	for _, rel := range files {
		abs := filepath.Join(cwd, rel)
		info, err := os.Lstat(abs)
		if err != nil {
			return NewArchiveErrorWithCause("Failed to read file info", 1, err)
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return NewArchiveErrorWithCause("Failed to create file header", 1, err)
		}
		hdr.Name = rel
		hdr.Method = zip.Deflate
		w, err := zipw.CreateHeader(hdr)
		if err != nil {
			if IsDiskFullError(err) {
				return NewArchiveError("Insufficient disk space to create archive", cfg.StatusDiskFull)
			}
			return NewArchiveErrorWithCause("Failed to create file in archive", 1, err)
		}
		if !info.IsDir() {
			rf, err := os.Open(abs)
			if err != nil {
				if IsPermissionError(err) {
					return NewArchiveError("Permission denied reading file", cfg.StatusPermissionDenied)
				}
				return NewArchiveErrorWithCause("Failed to open file", 1, err)
			}
			_, err = io.Copy(w, rf)
			rf.Close()
			if err != nil {
				if IsDiskFullError(err) {
					return NewArchiveError("Insufficient disk space to create archive", cfg.StatusDiskFull)
				}
				return NewArchiveErrorWithCause("Failed to write file to archive", 1, err)
			}
		}
	}

	// Close the zip writer to ensure all data is written
	zipw.Close()
	f.Close()

	// Atomic rename
	if err := os.Rename(tempPath, archivePath); err != nil {
		return NewArchiveErrorWithCause("Failed to finalize archive", 1, err)
	}

	// Remove temp file from cleanup list since operation succeeded
	rm.RemoveResource(&TempFile{Path: tempPath})

	// Create archive object for verification
	archive := &Archive{
		Name:          archiveName,
		Path:          archivePath,
		CreationTime:  time.Now(),
		IsIncremental: false,
		GitBranch:     gitBranch,
		GitHash:       gitHash,
		Note:          note,
	}

	// Generate and store checksums
	// Prepare absolute paths for checksum calculation
	var absFiles []string
	for _, rel := range files {
		absFiles = append(absFiles, filepath.Join(cwd, rel))
	}
	checksums, err := GenerateChecksums(absFiles, cfg.Verification.ChecksumAlgorithm)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to generate checksums", 1, err)
	}

	if err := StoreChecksums(archive, checksums); err != nil {
		return NewArchiveErrorWithCause("Failed to store checksums", 1, err)
	}

	// Verify the archive if requested
	if verify || cfg.Verification.VerifyOnCreate {
		status, err := VerifyArchive(archivePath)
		if err != nil {
			return NewArchiveErrorWithCause("Verification failed", 1, err)
		}

		if !status.IsVerified {
			return NewArchiveErrorWithContext("Archive verification failed", 1, "verify", archivePath, fmt.Errorf("errors: %v", status.Errors))
		}

		// Store verification status
		if err := StoreVerificationStatus(archive, status); err != nil {
			// Don't fail if we can't store status, just warn
			fmt.Printf("Warning: Could not store verification status: %v\n", err)
		}

		fmt.Println("Archive verified successfully")
	}

	// Use formatter for output
	formatter := NewOutputFormatter(cfg)
	formatter.PrintCreatedArchive(archivePath)
	return nil
}

// CreateIncrementalArchive creates an incremental archive
func CreateIncrementalArchive(cfg *Config, note string, dryRun bool, verify bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Determine archive directory
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			return err
		}
	}

	// Find latest full archive
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return err
	}
	if len(archives) == 0 {
		return fmt.Errorf("no archives found in %s", archiveDir)
	}

	// Find the most recent full archive
	var latestFullArchive *Archive
	for i := len(archives) - 1; i >= 0; i-- {
		if !archives[i].IsIncremental {
			latestFullArchive = &archives[i]
			break
		}
	}
	if latestFullArchive == nil {
		return fmt.Errorf("no full archive found in %s", archiveDir)
	}

	// Get the modification time of the latest full archive
	latestFullInfo, err := os.Stat(latestFullArchive.Path)
	if err != nil {
		return err
	}
	latestFullTime := latestFullInfo.ModTime()

	// Walk directory and collect modified files
	var modifiedFiles []string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if ShouldExcludeFile(rel, cfg.ExcludePatterns) {
			return nil
		}
		if info.ModTime().After(latestFullTime) {
			modifiedFiles = append(modifiedFiles, rel)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(modifiedFiles) == 0 {
		fmt.Println("No files modified since last full archive")
		return nil
	}

	// Git info
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	// Archive name
	timestamp := time.Now().Format("2006-01-02-15-04")
	archiveName := GenerateArchiveName("", timestamp, gitBranch, gitHash, note, isGit, true, latestFullArchive.Name)
	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		fmt.Println("[Dry Run] Files to include:")
		for _, f := range modifiedFiles {
			fmt.Println("  ", f)
		}
		fmt.Println("[Dry Run] Archive would be:", archivePath)
		return nil
	}

	// Create zip archive
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	zipw := zip.NewWriter(f)
	defer zipw.Close()

	for _, rel := range modifiedFiles {
		abs := filepath.Join(cwd, rel)
		info, err := os.Lstat(abs)
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = rel
		hdr.Method = zip.Deflate
		w, err := zipw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rf, err := os.Open(abs)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, rf)
			rf.Close()
			if err != nil {
				return err
			}
		}
	}

	// Close the zip writer to ensure all data is written
	zipw.Close()
	f.Close()

	// Create archive object for verification
	archive := &Archive{
		Name:          archiveName,
		Path:          archivePath,
		CreationTime:  time.Now(),
		IsIncremental: true,
		GitBranch:     gitBranch,
		GitHash:       gitHash,
		Note:          note,
		BaseArchive:   latestFullArchive.Name,
	}

	// Generate and store checksums
	if !dryRun {
		checksums, err := GenerateChecksums(modifiedFiles, cfg.Verification.ChecksumAlgorithm)
		if err != nil {
			return fmt.Errorf("failed to generate checksums: %w", err)
		}

		if err := StoreChecksums(archive, checksums); err != nil {
			return fmt.Errorf("failed to store checksums: %w", err)
		}
	}

	// Verify the archive if requested
	if verify || cfg.Verification.VerifyOnCreate {
		status, err := VerifyArchive(archivePath)
		if err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}

		if !status.IsVerified {
			return fmt.Errorf("archive verification failed: %v", status.Errors)
		}

		// Store verification status
		if err := StoreVerificationStatus(archive, status); err != nil {
			return fmt.Errorf("failed to store verification status: %w", err)
		}

		fmt.Println("Archive verified successfully")
	}

	fmt.Println("Created incremental archive:", archivePath)
	return nil
}
