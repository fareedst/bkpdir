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

// REFACTOR-001: Archive management interface contracts defined
// REFACTOR-001: Multiple dependency interfaces required for extraction
// REFACTOR-005: Structure optimization - Standardized naming and interface preparation

// REFACTOR-005: Extraction preparation - Standardized naming conventions
// ArchiveConfig holds configuration for generating archive names.
type ArchiveConfig struct {
	Prefix             string
	Timestamp          string
	GitBranch          string
	GitHash            string
	GitIsClean         bool
	ShowGitDirtyStatus bool
	Note               string
	IsGit              bool
	IsIncremental      bool
	BaseName           string
}

// REFACTOR-005: Structure optimization - Consistent naming pattern
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

// REFACTOR-005: Structure optimization - Interface-ready configuration
// ArchiveVerificationOptions holds configuration for archive verification
type ArchiveVerificationOptions struct {
	Path   string
	Config ArchiveConfigInterface
}

// REFACTOR-005: Structure optimization - Interface-based configuration abstraction
// ArchiveConfigInterface abstracts configuration dependencies for archive operations
type ArchiveConfigInterface interface {
	GetArchiveDirPath() string
	GetUseCurrentDirName() bool
	GetExcludePatterns() []string
	GetIncludeGitInfo() bool
	GetShowGitDirtyStatus() bool
	GetVerification() *VerificationConfig
	GetStatusCodes() map[string]int
	GetStatusDirectoryNotFound() int
	GetStatusDiskFull() int
	GetStatusConfigError() int
}

// REFACTOR-005: Structure optimization - Interface-ready configuration
// ArchiveCreationOptions holds configuration for archive creation
type ArchiveCreationOptions struct {
	Context     context.Context
	CWD         string
	Path        string
	Files       []string
	Config      ArchiveConfigInterface
	Verify      bool
	ResourceMgr *ResourceManager
}

// REFACTOR-005: Structure optimization - Interface-based formatter abstraction
// ArchiveFormatterInterface abstracts formatter dependencies for archive operations
type ArchiveFormatterInterface interface {
	PrintDryRunFilesHeader()
	PrintDryRunFileEntry(file string)
	PrintDryRunArchive(path string)
	PrintIncrementalCreated(path string)
}

// REFACTOR-005: Structure optimization - Interface wrapper for Config backward compatibility
// ConfigToArchiveConfigAdapter adapts Config to ArchiveConfigInterface
type ConfigToArchiveConfigAdapter struct {
	cfg *Config
}

func (a *ConfigToArchiveConfigAdapter) GetArchiveDirPath() string {
	return a.cfg.ArchiveDirPath
}

func (a *ConfigToArchiveConfigAdapter) GetUseCurrentDirName() bool {
	return a.cfg.UseCurrentDirName
}

func (a *ConfigToArchiveConfigAdapter) GetExcludePatterns() []string {
	return a.cfg.ExcludePatterns
}

func (a *ConfigToArchiveConfigAdapter) GetIncludeGitInfo() bool {
	return a.cfg.IncludeGitInfo
}

func (a *ConfigToArchiveConfigAdapter) GetShowGitDirtyStatus() bool {
	return a.cfg.ShowGitDirtyStatus
}

func (a *ConfigToArchiveConfigAdapter) GetVerification() *VerificationConfig {
	return a.cfg.Verification
}

func (a *ConfigToArchiveConfigAdapter) GetStatusCodes() map[string]int {
	return a.cfg.GetStatusCodes()
}

func (a *ConfigToArchiveConfigAdapter) GetStatusDirectoryNotFound() int {
	return a.cfg.StatusDirectoryNotFound
}

func (a *ConfigToArchiveConfigAdapter) GetStatusDiskFull() int {
	return a.cfg.StatusDiskFull
}

func (a *ConfigToArchiveConfigAdapter) GetStatusConfigError() int {
	return a.cfg.StatusConfigError
}

// REFACTOR-005: Structure optimization - Interface wrapper for OutputFormatter backward compatibility
// OutputFormatterToArchiveFormatterAdapter adapts OutputFormatter to ArchiveFormatterInterface
type OutputFormatterToArchiveFormatterAdapter struct {
	formatter *OutputFormatter
}

func (a *OutputFormatterToArchiveFormatterAdapter) PrintDryRunFilesHeader() {
	a.formatter.PrintDryRunFilesHeader()
}

func (a *OutputFormatterToArchiveFormatterAdapter) PrintDryRunFileEntry(file string) {
	a.formatter.PrintDryRunFileEntry(file)
}

func (a *OutputFormatterToArchiveFormatterAdapter) PrintDryRunArchive(path string) {
	a.formatter.PrintDryRunArchive(path)
}

func (a *OutputFormatterToArchiveFormatterAdapter) PrintIncrementalCreated(path string) {
	a.formatter.PrintIncrementalCreated(path)
}

// REFACTOR-005: Extraction preparation - Interface-based archive name generation
// GenerateArchiveNameWithInterface creates an archive name using interface abstractions
func GenerateArchiveNameWithInterface(cfg ArchiveConfig) string {
	if cfg.IsIncremental && cfg.BaseName != "" {
		return generateIncrementalArchiveName(cfg)
	}
	return generateFullArchiveNameFromConfig(cfg)
}

// ARCH-001: Archive naming convention implementation
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
// REFACTOR-005: Structure optimization - Standardized naming function
// GenerateArchiveName creates an archive name according to the spec.
// It handles both full and incremental archive naming based on the provided configuration.
func GenerateArchiveName(cfg ArchiveConfig) string {
	if cfg.IsIncremental && cfg.BaseName != "" {
		return generateIncrementalArchiveName(cfg)
	}
	return generateFullArchiveNameFromConfig(cfg)
}

// ARCH-003: Incremental archive naming implementation
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
// REFACTOR-005: Structure optimization - Consistent internal naming
// generateIncrementalArchiveName generates name for incremental archives
func generateIncrementalArchiveName(cfg ArchiveConfig) string {
	baseName := strings.TrimSuffix(cfg.BaseName, ".zip")
	name := baseName + "_update=" + cfg.Timestamp
	if cfg.IsGit && cfg.GitBranch != "" && cfg.GitHash != "" {
		name += "=" + cfg.GitBranch + "=" + cfg.GitHash
		if !cfg.GitIsClean && cfg.ShowGitDirtyStatus {
			name += "-dirty"
		}
	}
	if cfg.Note != "" {
		name += "=" + cfg.Note
	}
	return name + ".zip"
}

// ARCH-001: Full archive naming implementation
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
// REFACTOR-005: Structure optimization - Consistent internal naming
// generateFullArchiveNameFromConfig generates name for full archives from config
func generateFullArchiveNameFromConfig(cfg ArchiveConfig) string {
	var name string
	if cfg.Prefix != "" {
		name = cfg.Prefix + "-" + cfg.Timestamp
	} else {
		name = cfg.Timestamp
	}

	if cfg.IsGit && cfg.GitBranch != "" && cfg.GitHash != "" {
		name += "=" + cfg.GitBranch + "=" + cfg.GitHash
		if !cfg.GitIsClean && cfg.ShowGitDirtyStatus {
			name += "-dirty"
		}
	}

	if cfg.Note != "" {
		name += "=" + cfg.Note
	}

	return name + ".zip"
}

// ARCH-001: Archive naming with Git integration
// GIT-001: Git information extraction for naming
// GIT-003: Git status detection for naming
// IMMUTABLE-REF: Archive Naming Convention, Git Integration Requirements
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
// REFACTOR-005: Structure optimization - Standardized Git integration function
// GenerateFullArchiveName creates a full archive name with optional Git integration and note.
// It uses the current directory name as prefix and includes Git branch/hash if available.
func GenerateFullArchiveName(cfg *Config, cwd string, note string) (string, error) {
	timestamp := time.Now().Format("2006-01-02-15-04")
	prefix := filepath.Base(cwd)

	archiveConfig := ArchiveConfig{
		Prefix:             prefix,
		Timestamp:          timestamp,
		Note:               note,
		ShowGitDirtyStatus: cfg.ShowGitDirtyStatus,
	}

	if cfg.IncludeGitInfo {
		if IsGitRepository(cwd) {
			branch, hash, isClean := GetGitInfoWithStatus(cwd)
			archiveConfig.IsGit = true
			archiveConfig.GitBranch = branch
			archiveConfig.GitHash = hash
			archiveConfig.GitIsClean = isClean
		}
	}

	return GenerateArchiveName(archiveConfig), nil
}

// REFACTOR-005: Structure optimization - Backward compatibility wrapper
// generateFullArchiveName maintains backward compatibility while using new structure
func generateFullArchiveName(cfg *Config, cwd string, note string) (string, error) {
	return GenerateFullArchiveName(cfg, cwd, note)
}

// ARCH-002: Archive listing implementation
// IMMUTABLE-REF: Commands - List Archives
// TEST-REF: TestListArchives
// DECISION-REF: DEC-001
// ListArchives lists all archives in the archive directory for the current source.
// It returns a slice of Archive structs containing metadata for each archive found.
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

// ARCH-002: Archive metadata extraction
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestListArchives
// DECISION-REF: DEC-001
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

// ARCH-002: Archive creation with context
// IMMUTABLE-REF: Commands - Create Archive
// TEST-REF: TestCreateFullArchive
// DECISION-REF: DEC-001
// CreateArchiveWithContext creates a new archive with the given configuration and note.
// It supports both dry-run mode and verification of the created archive.
func CreateArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error {
	return CreateFullArchiveWithContext(ctx, cfg, note, dryRun, verify)
}

// ARCH-002: File collection for archiving
// IMMUTABLE-REF: Directory Operations, File Exclusion Requirements
// TEST-REF: TestCreateFullArchive
// DECISION-REF: DEC-001
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

	// REFACTOR-005: Structure optimization - Use interface adapter for reduced coupling
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: cfg}

	archiveDir, err := prepareArchiveDirectoryWithInterface(archiveConfig, cwd, dryRun)
	if err != nil {
		return err
	}

	files, err := collectFilesToArchiveWithInterface(ctx, cwd, archiveConfig.GetExcludePatterns())
	if err != nil {
		return NewArchiveErrorWithCause("Failed to scan directory", 1, err)
	}

	archiveName, err := generateFullArchiveNameWithInterface(archiveConfig, cwd, note)
	if err != nil {
		return err
	}

	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		printDryRunInfoWithInterface(files, archivePath, archiveConfig)
		return nil
	}

	return createAndVerifyArchive(ArchiveCreationOptions{
		Context:     ctx,
		CWD:         cwd,
		Path:        archivePath,
		Files:       files,
		Config:      archiveConfig,
		Verify:      verify,
		ResourceMgr: rm,
	})
}

// REFACTOR-005: Structure optimization - Interface-based directory preparation
// prepareArchiveDirectoryWithInterface prepares the archive directory using interface abstractions
func prepareArchiveDirectoryWithInterface(cfg ArchiveConfigInterface, cwd string, dryRun bool) (string, error) {
	archiveDir := cfg.GetArchiveDirPath()
	if cfg.GetUseCurrentDirName() {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		// Use the interface's config for SafeMkdirAll
		if concreteCfg, ok := cfg.(*ConfigToArchiveConfigAdapter); ok {
			if err := SafeMkdirAll(archiveDir, 0755, concreteCfg.cfg); err != nil {
				return "", err
			}
		} else {
			// Fallback: create directory without config-specific error handling
			if err := os.MkdirAll(archiveDir, 0755); err != nil {
				return "", NewArchiveError(fmt.Sprintf("Failed to create directory: %s", archiveDir), 1)
			}
		}
	}
	return archiveDir, nil
}

// REFACTOR-005: Structure optimization - Interface-based file collection
// collectFilesToArchiveWithInterface walks the directory and collects files to archive using interface abstractions
func collectFilesToArchiveWithInterface(ctx context.Context, cwd string, excludePatterns []string) ([]string, error) {
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

// REFACTOR-005: Structure optimization - Interface-based archive name generation
// generateFullArchiveNameWithInterface creates a full archive name using interface abstractions
func generateFullArchiveNameWithInterface(cfg ArchiveConfigInterface, cwd string, note string) (string, error) {
	timestamp := time.Now().Format("2006-01-02-15-04")
	prefix := filepath.Base(cwd)

	archiveConfig := ArchiveConfig{
		Prefix:             prefix,
		Timestamp:          timestamp,
		Note:               note,
		ShowGitDirtyStatus: cfg.GetShowGitDirtyStatus(),
	}

	if cfg.GetIncludeGitInfo() {
		if IsGitRepository(cwd) {
			branch, hash, isClean := GetGitInfoWithStatus(cwd)
			archiveConfig.IsGit = true
			archiveConfig.GitBranch = branch
			archiveConfig.GitHash = hash
			archiveConfig.GitIsClean = isClean
		}
	}

	return GenerateArchiveNameWithInterface(archiveConfig), nil
}

// REFACTOR-005: Structure optimization - Interface-based dry run printing
// printDryRunInfoWithInterface prints information about what would be archived using interface abstractions
func printDryRunInfoWithInterface(files []string, archivePath string, cfg ArchiveConfigInterface) {
	// Use the adapter to get the original config for OutputFormatter
	if concreteCfg, ok := cfg.(*ConfigToArchiveConfigAdapter); ok {
		formatter := NewOutputFormatter(concreteCfg.cfg)
		archiveFormatter := &OutputFormatterToArchiveFormatterAdapter{formatter: formatter}

		archiveFormatter.PrintDryRunFilesHeader()
		for _, f := range files {
			archiveFormatter.PrintDryRunFileEntry(f)
		}
		archiveFormatter.PrintDryRunArchive(archivePath)
	}
}

// REFACTOR-005: Structure optimization - Interface-based verification
// verifyArchiveWithInterface verifies an archive using interface abstractions
func verifyArchiveWithInterface(cfg ArchiveVerificationOptions) error {
	status, err := VerifyArchive(cfg.Path)
	if err != nil {
		return NewArchiveErrorWithCause(
			"Archive verification failed",
			cfg.Config.GetStatusConfigError(),
			err,
		)
	}
	if !status.IsVerified {
		return NewArchiveErrorWithCause(
			"Archive verification failed",
			cfg.Config.GetStatusConfigError(),
			fmt.Errorf("verification errors: %v", status.Errors),
		)
	}
	return nil
}

// prepareArchiveDirectory prepares the archive directory (backward compatibility).
func prepareArchiveDirectory(cfg *Config, cwd string, dryRun bool) (string, error) {
	// REFACTOR-005: Extraction preparation - Backward compatibility wrapper
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: cfg}
	return prepareArchiveDirectoryWithInterface(archiveConfig, cwd, dryRun)
}

// printDryRunInfo prints information about what would be archived (backward compatibility).
func printDryRunInfo(files []string, archivePath string, cfg *Config) {
	// REFACTOR-005: Extraction preparation - Backward compatibility wrapper
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: cfg}
	printDryRunInfoWithInterface(files, archivePath, archiveConfig)
}

// createAndVerifyArchive creates and verifies an archive.
func createAndVerifyArchive(cfg ArchiveCreationOptions) error {
	tempFile := cfg.Path + ".tmp"
	cfg.ResourceMgr.AddTempFile(tempFile)

	if err := createZipArchiveWithContext(cfg.Context, cfg.CWD, tempFile, cfg.Files); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to create archive",
			cfg.Config.GetStatusDiskFull(),
			err,
		)
	}

	if err := os.Rename(tempFile, cfg.Path); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to finalize archive",
			cfg.Config.GetStatusDiskFull(),
			err,
		)
	}

	cfg.ResourceMgr.RemoveResource(&TempFile{Path: tempFile})

	if cfg.Verify {
		verifyCfg := ArchiveVerificationOptions{
			Path:   cfg.Path,
			Config: cfg.Config,
		}
		return verifyArchiveWithInterface(verifyCfg)
	}

	return nil
}

// verifyArchive verifies an archive (backward compatibility).
func verifyArchive(cfg ArchiveVerificationOptions) error {
	// REFACTOR-005: Extraction preparation - Backward compatibility wrapper
	return verifyArchiveWithInterface(cfg)
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

	// REFACTOR-005: Structure optimization - Use interface adapter for reduced coupling
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: config.Config}

	archiveDir, err := prepareArchiveDirectoryWithInterface(archiveConfig, cwd, config.DryRun)
	if err != nil {
		return err
	}

	latestFullArchive, err := findLatestFullArchive(archiveDir)
	if err != nil {
		return err
	}

	modifiedFiles, err := collectModifiedFiles(cwd, latestFullArchive, archiveConfig.GetExcludePatterns())
	if err != nil {
		return err
	}

	if len(modifiedFiles) == 0 {
		// Use the adapter to get the original config for OutputFormatter
		formatter := NewOutputFormatter(config.Config)
		formatter.PrintNoFilesModified()
		return nil
	}

	archivePath, err := prepareIncrementalArchiveWithInterface(
		cwd, latestFullArchive, archiveConfig, config.Note)
	if err != nil {
		return err
	}

	if config.DryRun {
		printDryRunInfoWithInterface(modifiedFiles, archivePath, archiveConfig)
		return nil
	}

	return createAndVerifyIncrementalArchive(ArchiveCreationOptions{
		Context: config.Context,
		CWD:     cwd,
		Path:    archivePath,
		Files:   modifiedFiles,
		Config:  archiveConfig,
		Verify:  config.Verify,
	})
}

// REFACTOR-005: Structure optimization - Interface-based incremental archive preparation
// prepareIncrementalArchiveWithInterface prepares the archive name and path using interface abstractions
func prepareIncrementalArchiveWithInterface(
	cwd string, latestFullArchive *Archive, cfg ArchiveConfigInterface, note string) (string, error) {
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash, gitIsClean := "", "", false
	if isGit && cfg.GetIncludeGitInfo() {
		gitBranch, gitHash, gitIsClean = GetGitInfoWithStatus(cwd)
	}

	timestamp := time.Now().Format("2006-01-02-15-04")
	nameCfg := ArchiveConfig{
		Prefix:             "",
		Timestamp:          timestamp,
		GitBranch:          gitBranch,
		GitHash:            gitHash,
		GitIsClean:         gitIsClean,
		ShowGitDirtyStatus: cfg.GetShowGitDirtyStatus(),
		Note:               note,
		IsGit:              isGit && cfg.GetIncludeGitInfo(),
		IsIncremental:      true,
		BaseName:           latestFullArchive.Name,
	}
	archiveName := GenerateArchiveNameWithInterface(nameCfg)
	archivePath := filepath.Join(cfg.GetArchiveDirPath(), archiveName)
	return archivePath, nil
}

// createAndVerifyIncrementalArchive creates and verifies an incremental archive
func createAndVerifyIncrementalArchive(cfg ArchiveCreationOptions) error {
	if err := createZipArchiveWithContext(cfg.Context, cfg.CWD, cfg.Path, cfg.Files); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to create archive",
			cfg.Config.GetStatusDiskFull(),
			err,
		)
	}

	verificationConfig := cfg.Config.GetVerification()
	if cfg.Verify || verificationConfig.VerifyOnCreate {
		verifyCfg := ArchiveVerificationOptions{
			Path:   cfg.Path,
			Config: cfg.Config,
		}
		if err := verifyArchiveWithInterface(verifyCfg); err != nil {
			return err
		}
	}

	// Use the adapter to get the original config for OutputFormatter
	if concreteCfg, ok := cfg.Config.(*ConfigToArchiveConfigAdapter); ok {
		formatter := NewOutputFormatter(concreteCfg.cfg)
		formatter.PrintIncrementalCreated(cfg.Path)
	}
	return nil
}

// prepareIncrementalArchive prepares the archive name and path (backward compatibility)
func prepareIncrementalArchive(
	cwd string, latestFullArchive *Archive, cfg *Config, note string) (string, error) {
	// REFACTOR-005: Extraction preparation - Backward compatibility wrapper
	archiveConfig := &ConfigToArchiveConfigAdapter{cfg: cfg}
	return prepareIncrementalArchiveWithInterface(cwd, latestFullArchive, archiveConfig, note)
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
