// This file is part of bkpdir
//
// Package main provides file backup functionality for BkpDir.
// It handles individual file backup creation, comparison, and management.
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// REFACTOR-001: Backup management interface contracts defined
// REFACTOR-001: Multiple dependency interfaces required for extraction
// REFACTOR-004: Error handling now standardized in errors.go
// REFACTOR-005: Structure optimization - Interface-based backup operations

// REFACTOR-005: Structure optimization - Interface-based configuration abstraction
// BackupConfigInterface abstracts configuration dependencies for backup operations
type BackupConfigInterface interface {
	GetBackupDirPath() string
	GetUseCurrentDirNameForFiles() bool
	GetStatusFileNotFound() int
	GetStatusPermissionDenied() int
	GetStatusInvalidFileType() int
	GetStatusDiskFull() int
	GetStatusFileIsIdenticalToExistingBackup() int
}

// REFACTOR-005: Structure optimization - Interface-based formatter abstraction
// BackupFormatterInterface abstracts formatter dependencies for backup operations
type BackupFormatterInterface interface {
	PrintDryRunBackup(path string)
	PrintIdenticalBackup(path string)
	PrintBackupCreated(path string)
	PrintNoBackupsFound(filename, backupDir string)
}

// REFACTOR-005: Structure optimization - Interface wrapper for Config backward compatibility
// ConfigToBackupConfigAdapter adapts Config to BackupConfigInterface
type ConfigToBackupConfigAdapter struct {
	cfg *Config
}

func (a *ConfigToBackupConfigAdapter) GetBackupDirPath() string {
	return a.cfg.BackupDirPath
}

func (a *ConfigToBackupConfigAdapter) GetUseCurrentDirNameForFiles() bool {
	return a.cfg.UseCurrentDirNameForFiles
}

func (a *ConfigToBackupConfigAdapter) GetStatusFileNotFound() int {
	return a.cfg.StatusFileNotFound
}

func (a *ConfigToBackupConfigAdapter) GetStatusPermissionDenied() int {
	return a.cfg.StatusPermissionDenied
}

func (a *ConfigToBackupConfigAdapter) GetStatusInvalidFileType() int {
	return a.cfg.StatusInvalidFileType
}

func (a *ConfigToBackupConfigAdapter) GetStatusDiskFull() int {
	return a.cfg.StatusDiskFull
}

func (a *ConfigToBackupConfigAdapter) GetStatusFileIsIdenticalToExistingBackup() int {
	return a.cfg.StatusFileIsIdenticalToExistingBackup
}

// REFACTOR-005: Structure optimization - Interface wrapper for OutputFormatter backward compatibility
// OutputFormatterToBackupFormatterAdapter adapts OutputFormatter to BackupFormatterInterface
type OutputFormatterToBackupFormatterAdapter struct {
	formatter *OutputFormatter
}

func (a *OutputFormatterToBackupFormatterAdapter) PrintDryRunBackup(path string) {
	a.formatter.PrintDryRunBackup(path)
}

func (a *OutputFormatterToBackupFormatterAdapter) PrintIdenticalBackup(path string) {
	a.formatter.PrintIdenticalBackup(path)
}

func (a *OutputFormatterToBackupFormatterAdapter) PrintBackupCreated(path string) {
	a.formatter.PrintBackupCreated(path)
}

func (a *OutputFormatterToBackupFormatterAdapter) PrintNoBackupsFound(filename, backupDir string) {
	a.formatter.PrintNoBackupsFound(filename, backupDir)
}

// BackupInfo represents information about a file backup
type BackupInfo struct {
	Name         string
	Path         string
	CreationTime time.Time
	Size         int64
}

// Backup represents a single file backup with metadata
type Backup struct {
	Name         string
	Path         string
	CreationTime time.Time
	SourceFile   string
	Note         string
}

// BackupOptions holds parameters for backup creation functions
type BackupOptions struct {
	Context   context.Context
	Config    *Config
	Formatter *OutputFormatter
	FilePath  string
	Note      string
	DryRun    bool
}

// FILE-002: File backup creation implementation
// IMMUTABLE-REF: Commands - Create File Backup, File Backup Operations
// TEST-REF: TestCreateFileBackup
// DECISION-REF: DEC-002
// REFACTOR-004: Error handling now uses standardized patterns from errors.go
// CreateFileBackup creates a backup of a single file
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
	opts := BackupOptions{
		Context:   context.Background(),
		Config:    cfg,
		Formatter: nil,
		FilePath:  filePath,
		Note:      note,
		DryRun:    dryRun,
	}
	return createFileBackupInternal(opts)
}

// FILE-002: Core backup logic implementation
// IMMUTABLE-REF: File Backup Operations, Atomic Operations
// TEST-REF: TestCreateFileBackup
// DECISION-REF: DEC-002
// createFileBackupInternal handles the core backup logic
func createFileBackupInternal(opts BackupOptions) error {
	// Validate file
	if err := validateFileForBackup(opts.FilePath, opts.Config); err != nil {
		return err
	}

	// Generate backup path
	backupPath, err := generateBackupPath(opts.Config, opts.FilePath, opts.Note)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return handleDryRunBackup(opts.Formatter, backupPath)
	}

	return performBackupOperation(opts, backupPath)
}

// FILE-001: Backup path generation implementation
// IMMUTABLE-REF: File Backup Naming Convention
// TEST-REF: TestGenerateBackupName
// DECISION-REF: DEC-002
// generateBackupPath creates the full backup path including note
func generateBackupPath(cfg *Config, filePath, note string) (string, error) {
	backupPath, err := determineBackupPath(cfg, filePath)
	if err != nil {
		return "", err
	}

	if note != "" {
		backupPath += "=" + note
	}

	return backupPath, nil
}

// CFG-003: Dry run output formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestCreateFileBackup
// DECISION-REF: DEC-003
// handleDryRunBackup handles dry run mode output
func handleDryRunBackup(formatter *OutputFormatter, backupPath string) error {
	if formatter != nil {
		formatter.PrintDryRunBackup(backupPath)
	} else {
		// Fallback for legacy code paths
		cfg := DefaultConfig()
		fallbackFormatter := NewOutputFormatter(cfg)
		fallbackFormatter.PrintBackupWouldCreate(backupPath)
	}
	return nil
}

// FILE-002: Backup operation coordination
// IMMUTABLE-REF: File Backup Operations, Atomic Operations
// TEST-REF: TestCreateFileBackup
// DECISION-REF: DEC-002
// performBackupOperation performs the actual backup operation
func performBackupOperation(opts BackupOptions, backupPath string) error {
	backupDir := filepath.Dir(backupPath)
	baseFilename := filepath.Base(opts.FilePath)

	// Check for identical backup
	if err := checkAndHandleIdenticalBackup(opts, backupDir, baseFilename); err != nil {
		return err
	}

	// Create backup directory
	if err := SafeMkdirAll(backupDir, 0755, opts.Config); err != nil {
		return err
	}

	return executeBackupWithCleanup(opts, backupPath)
}

// FILE-003: Identical file backup detection
// IMMUTABLE-REF: File Backup Operations, Identical File Detection
// TEST-REF: TestCheckForIdenticalFileBackup
// DECISION-REF: DEC-002
// checkAndHandleIdenticalBackup checks if file is identical to existing backup
func checkAndHandleIdenticalBackup(opts BackupOptions, backupDir, baseFilename string) error {
	identical, existingBackup, err := CheckForIdenticalFileBackup(opts.FilePath, backupDir, baseFilename)
	if err == nil && identical {
		if opts.Formatter != nil {
			opts.Formatter.PrintIdenticalBackup(existingBackup)
		} else {
			// Fallback for legacy code paths
			cfg := DefaultConfig()
			fallbackFormatter := NewOutputFormatter(cfg)
			fallbackFormatter.PrintBackupIdentical(existingBackup)
		}
		os.Exit(opts.Config.StatusFileIsIdenticalToExistingBackup)
	}
	return nil
}

// FILE-002: Atomic backup execution with cleanup
// IMMUTABLE-REF: File Backup Operations, Atomic Operations, Resource Cleanup
// TEST-REF: TestCreateFileBackupWithCleanup
// DECISION-REF: DEC-002
// executeBackupWithCleanup performs backup with resource cleanup
func executeBackupWithCleanup(opts BackupOptions, backupPath string) error {
	// Create resource manager for cleanup
	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Create temporary file for atomic operation
	tempFile := backupPath + ".tmp"
	rm.AddTempFile(tempFile)

	// Copy file to backup location
	if err := copyFile(opts.FilePath, tempFile); err != nil {
		return NewArchiveErrorWithCause("Failed to create backup", opts.Config.StatusDiskFull, err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, backupPath); err != nil {
		return NewArchiveErrorWithCause("Failed to finalize backup", opts.Config.StatusDiskFull, err)
	}

	// Remove from cleanup list since operation succeeded
	rm.RemoveResource(&TempFile{Path: tempFile})

	// Create formatter for output (fallback since this function doesn't have direct access to opts.Formatter)
	formatter := NewOutputFormatter(opts.Config)
	formatter.PrintBackupCreated(backupPath)
	return nil
}

// CreateFileBackupEnhanced creates a backup with enhanced error handling and formatting
func CreateFileBackupEnhanced(opts BackupOptions) error {
	return createFileBackupInternal(opts)
}

// CheckForIdenticalFileBackup checks if the file is identical to the most recent backup
func CheckForIdenticalFileBackup(filePath, backupDir, baseFilename string) (bool, string, error) {
	// Find most recent backup for this file
	backups, err := ListFileBackups(backupDir, baseFilename)
	if err != nil || len(backups) == 0 {
		return false, "", err
	}

	// Get the most recent backup
	mostRecent := backups[0]

	// Compare file sizes first (quick check)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, "", err
	}

	backupInfo, err := os.Stat(mostRecent.Path)
	if err != nil {
		return false, "", err
	}

	if fileInfo.Size() != backupInfo.Size() {
		return false, "", nil
	}

	// Compare file contents
	identical, err := compareFiles(filePath, mostRecent.Path)
	if err != nil {
		return false, "", err
	}

	return identical, mostRecent.Path, nil
}

// ListFileBackups lists all backups for a specific file
func ListFileBackups(backupDir, baseFilename string) ([]BackupInfo, error) {
	var backups []BackupInfo

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return backups, nil
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	prefix := baseFilename + "-"
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backupPath := filepath.Join(backupDir, name)
		backup := BackupInfo{
			Name:         name,
			Path:         backupPath,
			CreationTime: info.ModTime(),
			Size:         info.Size(),
		}

		backups = append(backups, backup)
	}

	// Sort by creation time (most recent first)
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].CreationTime.Before(backups[j].CreationTime) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}

	return backups, nil
}

// ListFileBackupsEnhanced lists backups with enhanced formatting
func ListFileBackupsEnhanced(cfg *Config, formatter *OutputFormatter, filePath string) error {
	baseFilename := filepath.Base(filePath)

	backupDir := cfg.BackupDirPath
	if cfg.UseCurrentDirNameForFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
		}

		// Maintain directory structure in backup path
		relPath, err := filepath.Rel(cwd, filePath)
		if err != nil {
			// If we can't get relative path, use absolute path structure
			relPath = strings.TrimPrefix(filePath, "/")
		}

		backupDir = filepath.Join(backupDir, filepath.Dir(relPath))
	}

	backups, err := ListFileBackups(backupDir, baseFilename)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to list backups", 1, err)
	}

	if len(backups) == 0 {
		formatter.PrintNoBackupsFound(baseFilename, backupDir)
		return nil
	}

	for _, backup := range backups {
		creationTime := backup.CreationTime.Format("2006-01-02 15:04:05")
		output := formatter.FormatListBackupWithExtraction(backup.Path, creationTime)
		fmt.Print(output)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// compareFiles compares two files byte by byte with simplified logic
func compareFiles(file1, file2 string) (bool, error) {
	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	return compareFileContents(f1, f2)
}

// compareFileContents compares the contents of two open files
func compareFileContents(f1, f2 *os.File) (bool, error) {
	buf1 := make([]byte, 4096)
	buf2 := make([]byte, 4096)

	for {
		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		if n1 != n2 {
			return false, nil
		}

		// Handle end of file cases
		if err1 != nil || err2 != nil {
			return handleFileComparisonErrors(err1, err2)
		}

		// Compare buffer contents
		if !compareBuffers(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
	}
}

// handleFileComparisonErrors handles EOF and other errors during file comparison
func handleFileComparisonErrors(err1, err2 error) (bool, error) {
	if err1 == io.EOF && err2 == io.EOF {
		return true, nil
	}
	if err1 != nil {
		return false, err1
	}
	return false, err2
}

// compareBuffers compares two byte slices
func compareBuffers(buf1, buf2 []byte) bool {
	for i := 0; i < len(buf1); i++ {
		if buf1[i] != buf2[i] {
			return false
		}
	}
	return true
}

// validateFileForBackup checks if the file exists and is a regular file
func validateFileForBackup(filePath string, cfg *Config) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewArchiveError("File not found", cfg.StatusFileNotFound)
		}
		return NewArchiveErrorWithCause("Failed to access file", cfg.StatusPermissionDenied, err)
	}

	if !fileInfo.Mode().IsRegular() {
		return NewArchiveError("Path is not a regular file", cfg.StatusInvalidFileType)
	}

	return nil
}

// determineBackupPath determines the backup directory and filename
func determineBackupPath(cfg *Config, filePath string) (string, error) {
	backupDir := cfg.BackupDirPath
	if cfg.UseCurrentDirNameForFiles {
		cwd, err := os.Getwd()
		if err != nil {
			return "", NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
		}

		// Maintain directory structure in backup path
		relPath, err := filepath.Rel(cwd, filePath)
		if err != nil {
			// If we can't get relative path, use absolute path structure
			relPath = strings.TrimPrefix(filePath, "/")
		}

		backupDir = filepath.Join(backupDir, filepath.Dir(relPath))
	}

	// Generate backup filename
	baseFilename := filepath.Base(filePath)
	timestamp := time.Now().Format("2006-01-02-15-04")
	backupFilename := fmt.Sprintf("%s-%s", baseFilename, timestamp)

	return filepath.Join(backupDir, backupFilename), nil
}

// CreateFileBackupWithContext creates a backup with context support for cancellation
func CreateFileBackupWithContext(ctx context.Context, cfg *Config, filePath string, note string, dryRun bool) error {
	opts := BackupOptions{
		Context:   ctx,
		Config:    cfg,
		Formatter: nil,
		FilePath:  filePath,
		Note:      note,
		DryRun:    dryRun,
	}
	return createFileBackupWithContextInternal(opts)
}

// createFileBackupWithContextInternal handles context-aware backup creation
func createFileBackupWithContextInternal(opts BackupOptions) error {
	// Check for cancellation
	if err := opts.Context.Err(); err != nil {
		return err
	}

	// Validate file
	if err := validateFileForBackup(opts.FilePath, opts.Config); err != nil {
		return err
	}

	// Generate backup path
	backupPath, err := generateBackupPath(opts.Config, opts.FilePath, opts.Note)
	if err != nil {
		return err
	}

	if opts.DryRun {
		return handleDryRunBackup(opts.Formatter, backupPath)
	}

	return performContextAwareBackup(opts, backupPath)
}

// performContextAwareBackup performs backup with context cancellation support
func performContextAwareBackup(opts BackupOptions, backupPath string) error {
	backupDir := filepath.Dir(backupPath)
	baseFilename := filepath.Base(opts.FilePath)

	// Check for identical backup
	if err := checkAndHandleIdenticalBackup(opts, backupDir, baseFilename); err != nil {
		return err
	}

	// Create backup directory
	if err := SafeMkdirAll(backupDir, 0755, opts.Config); err != nil {
		return err
	}

	return executeContextAwareBackup(opts, backupPath)
}

// executeContextAwareBackup performs the backup with context support
func executeContextAwareBackup(opts BackupOptions, backupPath string) error {
	// Create resource manager for cleanup
	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Create temporary file for atomic operation
	tempFile := backupPath + ".tmp"
	rm.AddTempFile(tempFile)

	// Copy file to backup location with context
	if err := CopyFileWithContext(opts.Context, opts.FilePath, tempFile); err != nil {
		return NewArchiveErrorWithCause("Failed to create backup", opts.Config.StatusDiskFull, err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, backupPath); err != nil {
		return NewArchiveErrorWithCause("Failed to finalize backup", opts.Config.StatusDiskFull, err)
	}

	// Remove from cleanup list since operation succeeded
	rm.RemoveResource(&TempFile{Path: tempFile})

	// Create formatter for output (fallback since this function doesn't have direct access to opts.Formatter)
	formatter := NewOutputFormatter(opts.Config)
	formatter.PrintBackupCreated(backupPath)
	return nil
}

// CreateFileBackupWithContextAndCleanup creates a backup with context and enhanced cleanup
func CreateFileBackupWithContextAndCleanup(ctx context.Context, cfg *Config, filePath string,
	note string, dryRun bool) error {
	return CreateFileBackupWithContext(ctx, cfg, filePath, note, dryRun)
}

// CopyFileWithContext copies a file with context support for cancellation
func CopyFileWithContext(ctx context.Context, src, dst string) error {
	// Check for cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	return copyWithContextChecks(ctx, sourceFile, destFile)
}

// copyWithContextChecks performs the copy operation with periodic context checks
func copyWithContextChecks(ctx context.Context, src, dst *os.File) error {
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		// Check for cancellation
		if err := ctx.Err(); err != nil {
			return err
		}

		n, err := src.Read(buf)
		if n > 0 {
			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// Copy file permissions
	sourceInfo, err := src.Stat()
	if err != nil {
		return err
	}

	return dst.Chmod(sourceInfo.Mode())
}

// CopyFile creates exact copy with permissions preserved
func CopyFile(src, dst string) error {
	return CopyFileWithContext(context.Background(), src, dst)
}

// GenerateBackupName generates backup filename with timestamp and note
func GenerateBackupName(sourcePath, timestamp, note string) string {
	baseFilename := filepath.Base(sourcePath)
	backupName := fmt.Sprintf("%s-%s", baseFilename, timestamp)

	if note != "" {
		backupName += "=" + note
	}

	return backupName
}

// CreateFileBackupWithCleanup creates a backup with enhanced resource cleanup
func CreateFileBackupWithCleanup(cfg *Config, filePath string, note string, dryRun bool) error {
	return CreateFileBackupWithContext(context.Background(), cfg, filePath, note, dryRun)
}

// CompareFilesWithContext compares two files with context support for cancellation
func CompareFilesWithContext(ctx context.Context, file1, file2 string) (bool, error) {
	// Check for cancellation
	if err := ctx.Err(); err != nil {
		return false, err
	}

	f1, err := os.Open(file1)
	if err != nil {
		return false, err
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		return false, err
	}
	defer f2.Close()

	return compareFileContentsWithContext(ctx, f1, f2)
}

// compareFileContentsWithContext compares the contents of two open files with context
func compareFileContentsWithContext(ctx context.Context, f1, f2 *os.File) (bool, error) {
	buf1 := make([]byte, 4096)
	buf2 := make([]byte, 4096)

	for {
		// Check for cancellation
		if err := ctx.Err(); err != nil {
			return false, err
		}

		n1, err1 := f1.Read(buf1)
		n2, err2 := f2.Read(buf2)

		if n1 != n2 {
			return false, nil
		}

		// Handle end of file cases
		if err1 != nil || err2 != nil {
			return handleFileComparisonErrors(err1, err2)
		}

		// Compare buffer contents
		if !compareBuffers(buf1[:n1], buf2[:n2]) {
			return false, nil
		}
	}
}

// ListFileBackupsWithContext lists all backups for a specific file with context support
func ListFileBackupsWithContext(ctx context.Context, backupDir, sourceFile string) ([]*Backup, error) {
	// Check for cancellation
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var backups []*Backup
	baseFilename := filepath.Base(sourceFile)

	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return backups, nil
	}

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	backups, err = processBackupEntries(ctx, entries, backupDir, baseFilename, sourceFile)
	if err != nil {
		return nil, err
	}

	// Sort by creation time (most recent first)
	sortBackupsByCreationTime(backups)

	return backups, nil
}

// processBackupEntries processes directory entries to create backup objects
func processBackupEntries(
	ctx context.Context,
	entries []os.DirEntry,
	backupDir, baseFilename, sourceFile string,
) ([]*Backup, error) {
	var backups []*Backup
	prefix := baseFilename + "-"

	for _, entry := range entries {
		// Check for cancellation
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}

		backup, err := createBackupFromEntry(entry, backupDir, sourceFile)
		if err != nil {
			continue
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// createBackupFromEntry creates a Backup object from a directory entry
func createBackupFromEntry(entry os.DirEntry, backupDir, sourceFile string) (*Backup, error) {
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	name := entry.Name()
	backupPath := filepath.Join(backupDir, name)

	// Extract note from filename if present
	note := ""
	if idx := strings.Index(name, "="); idx != -1 {
		note = name[idx+1:]
	}

	backup := &Backup{
		Name:         name,
		Path:         backupPath,
		CreationTime: info.ModTime(),
		SourceFile:   sourceFile,
		Note:         note,
	}

	return backup, nil
}

// sortBackupsByCreationTime sorts backups by creation time (most recent first)
func sortBackupsByCreationTime(backups []*Backup) {
	for i := 0; i < len(backups)-1; i++ {
		for j := i + 1; j < len(backups); j++ {
			if backups[i].CreationTime.Before(backups[j].CreationTime) {
				backups[i], backups[j] = backups[j], backups[i]
			}
		}
	}
}

// GetMostRecentBackup finds the most recent backup for a given file
func GetMostRecentBackup(backupDir, sourceFile string) (*Backup, error) {
	backups, err := ListFileBackupsWithContext(context.Background(), backupDir, sourceFile)
	if err != nil {
		return nil, err
	}

	if len(backups) == 0 {
		return nil, nil
	}

	return backups[0], nil // Already sorted with most recent first
}

// ValidateFileForBackup validates that a file can be backed up (public version)
func ValidateFileForBackup(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filePath)
		}
		return fmt.Errorf("failed to access file: %v", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return fmt.Errorf("path is not a regular file: %s", filePath)
	}

	return nil
}
