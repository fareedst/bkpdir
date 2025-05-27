// This file is part of bkpdir
//
// Package main provides archive creation and management for BkpDir.
// It handles full and incremental archive creation, naming, and verification.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArchiveNameConfig holds configuration for generating archive names.
type ArchiveNameConfig struct {
	Prefix        string
	Timestamp     string
	GitBranch     string
	GitHash       string
	Note          string
	IsGit         bool
	IsIncremental bool
	BaseName      string
}

// Archive represents a directory archive with metadata including name, path,
// creation time, Git information, and verification status. It supports both
// full and incremental archives.
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

// ArchiveVerificationConfig holds configuration for archive verification
type ArchiveVerificationConfig struct {
	Path   string
	Config *Config
}

// ArchiveCreationConfig holds configuration for archive creation
type ArchiveCreationConfig struct {
	Context     context.Context
	CWD         string
	Path        string
	Files       []string
	Config      *Config
	Verify      bool
	ResourceMgr *ResourceManager
}

// GenerateArchiveName creates an archive name according to the spec.
func GenerateArchiveName(cfg ArchiveNameConfig) string {
	if cfg.IsIncremental && cfg.BaseName != "" {
		return generateIncrementalArchiveName(cfg)
	}
	return generateFullArchiveNameFromConfig(cfg)
}

// generateIncrementalArchiveName generates name for incremental archives
func generateIncrementalArchiveName(cfg ArchiveNameConfig) string {
	baseName := strings.TrimSuffix(cfg.BaseName, ".zip")
	name := baseName + "_update=" + cfg.Timestamp
	if cfg.IsGit && cfg.GitBranch != "" && cfg.GitHash != "" {
		name += "=" + cfg.GitBranch + "=" + cfg.GitHash
	}
	if cfg.Note != "" {
		name += "=" + cfg.Note
	}
	return name + ".zip"
}

// generateFullArchiveNameFromConfig generates name for full archives from config
func generateFullArchiveNameFromConfig(cfg ArchiveNameConfig) string {
	var name string
	if cfg.Prefix != "" {
		name = cfg.Prefix + "-" + cfg.Timestamp
	} else {
		name = cfg.Timestamp
	}

	if cfg.IsGit && cfg.GitBranch != "" && cfg.GitHash != "" {
		name += "=" + cfg.GitBranch + "=" + cfg.GitHash
	}

	if cfg.Note != "" {
		name += "=" + cfg.Note
	}

	return name + ".zip"
}

// generateFullArchiveName generates the name for a full archive.
func generateFullArchiveName(cfg *Config, cwd, note string) (string, error) {
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit && cfg.IncludeGitInfo {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	timestamp := time.Now().Format("2006-01-02-15-04")
	prefix := ""
	if cfg.UseCurrentDirName {
		prefix = filepath.Base(cwd)
	}

	nameCfg := ArchiveNameConfig{
		Prefix:        prefix,
		Timestamp:     timestamp,
		GitBranch:     gitBranch,
		GitHash:       gitHash,
		Note:          note,
		IsGit:         isGit && cfg.IncludeGitInfo,
		IsIncremental: false,
	}

	return GenerateArchiveName(nameCfg), nil
}

// ListArchives lists all archives in the archive directory for the current source.
func ListArchives(archiveDir string) ([]Archive, error) {
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

		archive, err := createArchiveFromEntry(archiveDir, entry)
		if err != nil {
			continue // Skip entries we can't process
		}
		archives = append(archives, archive)
	}
	return archives, nil
}

// createArchiveFromEntry creates an Archive from a directory entry.
func createArchiveFromEntry(archiveDir string, entry os.DirEntry) (Archive, error) {
	archivePath := filepath.Join(archiveDir, entry.Name())
	fileInfo, err := entry.Info()
	if err != nil {
		return Archive{}, err
	}

	archive := Archive{
		Name:          entry.Name(),
		Path:          archivePath,
		IsIncremental: strings.Contains(entry.Name(), "_update="),
		CreationTime:  fileInfo.ModTime(),
	}

	// Load verification status if available
	status, err := LoadVerificationStatus(&archive)
	if err == nil && status != nil {
		archive.VerificationStatus = status
	}

	return archive, nil
}

// CreateArchiveWithContext creates an archive with context support for cancellation
func CreateArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error {
	return CreateFullArchiveWithContext(ctx, cfg, note, dryRun, verify)
}

// collectFilesToArchive walks the directory and collects files to archive
func collectFilesToArchive(ctx context.Context, cwd string, excludePatterns []string) ([]string, error) {
	var files []string
	err := filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err := checkContextCancellation(ctx); err != nil {
			return err
		}

		if err != nil {
			return err
		}

		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}

		if rel == "." || info.IsDir() || ShouldExcludeFile(rel, excludePatterns) {
			return nil
		}

		files = append(files, rel)
		return nil
	})
	return files, err
}

// checkContextCancellation checks if the context has been cancelled.
func checkContextCancellation(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// CreateFullArchiveWithContext creates a full archive with context support
func CreateFullArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	if err := checkContextCancellation(ctx); err != nil {
		return err
	}

	if err := ValidateDirectoryPath(cwd, cfg); err != nil {
		return err
	}

	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	archiveDir, err := prepareArchiveDirectory(cfg, cwd, dryRun)
	if err != nil {
		return err
	}

	files, err := collectFilesToArchive(ctx, cwd, cfg.ExcludePatterns)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to scan directory", 1, err)
	}

	archiveName, err := generateFullArchiveName(cfg, cwd, note)
	if err != nil {
		return err
	}

	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		printDryRunInfo(files, archivePath, cfg)
		return nil
	}

	return createAndVerifyArchive(ArchiveCreationConfig{
		Context:     ctx,
		CWD:         cwd,
		Path:        archivePath,
		Files:       files,
		Config:      cfg,
		Verify:      verify,
		ResourceMgr: rm,
	})
}

// prepareArchiveDirectory prepares the archive directory.
func prepareArchiveDirectory(cfg *Config, cwd string, dryRun bool) (string, error) {
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		if err := SafeMkdirAll(archiveDir, 0755, cfg); err != nil {
			return "", err
		}
	}
	return archiveDir, nil
}

// printDryRunInfo prints information about what would be archived.
func printDryRunInfo(files []string, archivePath string, cfg *Config) {
	fmt.Println("[Dry Run] Files to include:")
	for _, f := range files {
		fmt.Println("  ", f)
	}
	formatter := NewOutputFormatter(cfg)
	formatter.PrintDryRunArchive(archivePath)
}

// createAndVerifyArchive creates and verifies an archive.
func createAndVerifyArchive(cfg ArchiveCreationConfig) error {
	tempFile := cfg.Path + ".tmp"
	cfg.ResourceMgr.AddTempFile(tempFile)

	if err := createZipArchiveWithContext(cfg.Context, cfg.CWD, tempFile, cfg.Files); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to create archive",
			cfg.Config.StatusDiskFull,
			err,
		)
	}

	if err := os.Rename(tempFile, cfg.Path); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to finalize archive",
			cfg.Config.StatusDiskFull,
			err,
		)
	}

	cfg.ResourceMgr.RemoveResource(&TempFile{Path: tempFile})

	if cfg.Verify {
		verifyCfg := ArchiveVerificationConfig{
			Path:   cfg.Path,
			Config: cfg.Config,
		}
		return verifyArchive(verifyCfg)
	}

	return nil
}

// verifyArchive verifies an archive.
func verifyArchive(cfg ArchiveVerificationConfig) error {
	status, err := VerifyArchive(cfg.Path)
	if err != nil {
		return NewArchiveErrorWithCause(
			"Archive verification failed",
			cfg.Config.StatusConfigError,
			err,
		)
	}
	if !status.IsVerified {
		return NewArchiveErrorWithCause(
			"Archive verification failed",
			cfg.Config.StatusConfigError,
			fmt.Errorf("verification errors: %v", status.Errors),
		)
	}
	return nil
}

// CreateFullArchive creates a full archive without context (backward compatibility)
func CreateFullArchive(cfg *Config, note string, dryRun bool, verify bool) error {
	return CreateFullArchiveWithContext(context.Background(), cfg, note, dryRun, verify)
}

// CreateFullArchiveWithCleanup creates a full archive with enhanced resource cleanup
func CreateFullArchiveWithCleanup(cfg *Config, note string, dryRun bool, verify bool) error {
	return CreateFullArchiveWithContext(context.Background(), cfg, note, dryRun, verify)
}

// IncrementalArchiveConfig holds configuration for creating incremental archives
type IncrementalArchiveConfig struct {
	Config  *Config
	Note    string
	DryRun  bool
	Verify  bool
	Context context.Context
}

// CreateIncrementalArchive creates an incremental archive without context (backward compatibility)
func CreateIncrementalArchive(cfg *Config, note string, dryRun bool, verify bool) error {
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    note,
		DryRun:  dryRun,
		Verify:  verify,
		Context: context.Background(),
	}
	return createIncrementalArchive(config)
}

// createIncrementalArchive is the core implementation for incremental archive creation
func createIncrementalArchive(config IncrementalArchiveConfig) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	archiveDir, err := prepareArchiveDirectory(config.Config, cwd, config.DryRun)
	if err != nil {
		return err
	}

	latestFullArchive, err := findLatestFullArchive(archiveDir)
	if err != nil {
		return err
	}

	modifiedFiles, err := collectModifiedFiles(cwd, latestFullArchive, config.Config.ExcludePatterns)
	if err != nil {
		return err
	}

	if len(modifiedFiles) == 0 {
		fmt.Println("No files modified since last full archive")
		return nil
	}

	archivePath, err := prepareIncrementalArchive(
		cwd, latestFullArchive, config.Config, config.Note)
	if err != nil {
		return err
	}

	if config.DryRun {
		printDryRunInfo(modifiedFiles, archivePath, config.Config)
		return nil
	}

	return createAndVerifyIncrementalArchive(ArchiveCreationConfig{
		Context: config.Context,
		CWD:     cwd,
		Path:    archivePath,
		Files:   modifiedFiles,
		Config:  config.Config,
		Verify:  config.Verify,
	})
}

// collectModifiedFiles collects files modified since the last full archive
func collectModifiedFiles(cwd string, latestFullArchive *Archive, excludePatterns []string) ([]string, error) {
	latestFullInfo, err := os.Stat(latestFullArchive.Path)
	if err != nil {
		return nil, err
	}
	latestFullTime := latestFullInfo.ModTime()

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
		if ShouldExcludeFile(rel, excludePatterns) {
			return nil
		}
		if info.ModTime().After(latestFullTime) {
			modifiedFiles = append(modifiedFiles, rel)
		}
		return nil
	})
	return modifiedFiles, err
}

// prepareIncrementalArchive prepares the archive name and path
func prepareIncrementalArchive(
	cwd string, latestFullArchive *Archive, cfg *Config, note string) (string, error) {
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit && cfg.IncludeGitInfo {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	timestamp := time.Now().Format("2006-01-02-15-04")
	nameCfg := ArchiveNameConfig{
		Prefix:        "",
		Timestamp:     timestamp,
		GitBranch:     gitBranch,
		GitHash:       gitHash,
		Note:          note,
		IsGit:         isGit && cfg.IncludeGitInfo,
		IsIncremental: true,
		BaseName:      latestFullArchive.Name,
	}
	archiveName := GenerateArchiveName(nameCfg)
	archivePath := filepath.Join(cfg.ArchiveDirPath, archiveName)
	return archivePath, nil
}

// createAndVerifyIncrementalArchive creates and verifies an incremental archive
func createAndVerifyIncrementalArchive(cfg ArchiveCreationConfig) error {
	if err := createZipArchiveWithContext(cfg.Context, cfg.CWD, cfg.Path, cfg.Files); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to create archive",
			cfg.Config.StatusDiskFull,
			err,
		)
	}

	if cfg.Verify || cfg.Config.Verification.VerifyOnCreate {
		verifyCfg := ArchiveVerificationConfig{
			Path:   cfg.Path,
			Config: cfg.Config,
		}
		if err := verifyArchive(verifyCfg); err != nil {
			return err
		}
	}

	fmt.Println("Created incremental archive:", cfg.Path)
	return nil
}

// CreateIncrementalArchiveWithContext creates an incremental archive with context support
func CreateIncrementalArchiveWithContext(
	ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error {
	config := IncrementalArchiveConfig{
		Config:  cfg,
		Note:    note,
		DryRun:  dryRun,
		Verify:  verify,
		Context: ctx,
	}
	return createIncrementalArchive(config)
}

// createZipArchiveWithContext creates a ZIP archive with context cancellation support
func createZipArchiveWithContext(ctx context.Context, sourceDir, archivePath string, files []string) error {
	if err := checkContextCancellation(ctx); err != nil {
		return err
	}

	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zipw := zip.NewWriter(f)
	defer zipw.Close()

	return addFilesToZip(ctx, sourceDir, files, zipw)
}

// addFilesToZip adds files to a zip archive
func addFilesToZip(ctx context.Context, sourceDir string, files []string, zipw *zip.Writer) error {
	for _, rel := range files {
		if err := checkContextCancellation(ctx); err != nil {
			return err
		}

		if err := addFileToZip(sourceDir, rel, zipw); err != nil {
			return err
		}
	}
	return nil
}

// addFileToZip adds a single file to a zip archive
func addFileToZip(sourceDir, rel string, zipw *zip.Writer) error {
	abs := filepath.Join(sourceDir, rel)
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

	return nil
}

// findLatestFullArchive finds the most recent full archive in the archive directory.
func findLatestFullArchive(archiveDir string) (*Archive, error) {
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return nil, err
	}
	if len(archives) == 0 {
		return nil, NewArchiveError("No archives found", 1)
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
		return nil, NewArchiveError("No full archive found", 1)
	}

	// Get the modification time of the latest full archive
	latestFullInfo, err := os.Stat(latestFullArchive.Path)
	if err != nil {
		return nil, NewArchiveErrorWithCause(
			"Failed to stat latest full archive",
			1,
			err,
		)
	}
	latestFullArchive.CreationTime = latestFullInfo.ModTime()

	return latestFullArchive, nil
}
