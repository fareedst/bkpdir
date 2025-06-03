// This file is part of bkpdir
//
// Package main provides adapter implementations for component extraction.
// These adapters bridge existing structures to new interfaces while maintaining
// full backward compatibility.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// 🔶 REFACTOR-005: Structure optimization - Comprehensive adapter implementations - 📝
// This file contains all adapter implementations to bridge existing code to new interfaces

// ========================================
// CONFIGURATION ADAPTERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Universal configuration provider adapter - 🔧
// ConfigProviderAdapter adapts Config struct to ConfigProviderInterface
type ConfigProviderAdapter struct {
	cfg *Config
}

// NewConfigProviderAdapter creates a new ConfigProviderAdapter
func NewConfigProviderAdapter(cfg *Config) *ConfigProviderAdapter {
	return &ConfigProviderAdapter{cfg: cfg}
}

func (a *ConfigProviderAdapter) GetArchiveConfig() ArchiveConfigInterface {
	return &ArchiveConfigAdapter{cfg: a.cfg}
}

func (a *ConfigProviderAdapter) GetBackupConfig() BackupConfigInterface {
	return &BackupConfigAdapter{cfg: a.cfg}
}

func (a *ConfigProviderAdapter) GetFormatterConfig() FormatterConfigInterface {
	return &FormatterConfigAdapter{cfg: a.cfg}
}

func (a *ConfigProviderAdapter) GetErrorConfig() ErrorConfigInterface {
	return &ErrorConfigAdapter{cfg: a.cfg}
}

func (a *ConfigProviderAdapter) GetGitConfig() GitConfigInterface {
	return &GitConfigAdapter{cfg: a.cfg}
}

// 🔶 REFACTOR-005: Structure optimization - Formatter configuration adapter - 🔧
// FormatterConfigAdapter adapts Config struct to FormatterConfigInterface
type FormatterConfigAdapter struct {
	cfg *Config
}

func (a *FormatterConfigAdapter) GetFormatCreatedArchive() string { return a.cfg.FormatCreatedArchive }
func (a *FormatterConfigAdapter) GetFormatIdenticalArchive() string {
	return a.cfg.FormatIdenticalArchive
}
func (a *FormatterConfigAdapter) GetFormatListArchive() string   { return a.cfg.FormatListArchive }
func (a *FormatterConfigAdapter) GetFormatConfigValue() string   { return a.cfg.FormatConfigValue }
func (a *FormatterConfigAdapter) GetFormatDryRunArchive() string { return a.cfg.FormatDryRunArchive }
func (a *FormatterConfigAdapter) GetFormatError() string         { return a.cfg.FormatError }
func (a *FormatterConfigAdapter) GetFormatCreatedBackup() string { return a.cfg.FormatCreatedBackup }
func (a *FormatterConfigAdapter) GetFormatIdenticalBackup() string {
	return a.cfg.FormatIdenticalBackup
}
func (a *FormatterConfigAdapter) GetFormatListBackup() string   { return a.cfg.FormatListBackup }
func (a *FormatterConfigAdapter) GetFormatDryRunBackup() string { return a.cfg.FormatDryRunBackup }

func (a *FormatterConfigAdapter) GetFormatNoArchivesFound() string {
	return a.cfg.FormatNoArchivesFound
}
func (a *FormatterConfigAdapter) GetFormatVerificationFailed() string {
	return a.cfg.FormatVerificationFailed
}
func (a *FormatterConfigAdapter) GetFormatVerificationSuccess() string {
	return a.cfg.FormatVerificationSuccess
}
func (a *FormatterConfigAdapter) GetFormatVerificationWarning() string {
	return a.cfg.FormatVerificationWarning
}
func (a *FormatterConfigAdapter) GetFormatConfigurationUpdated() string {
	return a.cfg.FormatConfigurationUpdated
}
func (a *FormatterConfigAdapter) GetFormatConfigFilePath() string { return a.cfg.FormatConfigFilePath }
func (a *FormatterConfigAdapter) GetFormatDryRunFilesHeader() string {
	return a.cfg.FormatDryRunFilesHeader
}
func (a *FormatterConfigAdapter) GetFormatDryRunFileEntry() string {
	return a.cfg.FormatDryRunFileEntry
}
func (a *FormatterConfigAdapter) GetFormatNoFilesModified() string {
	return a.cfg.FormatNoFilesModified
}
func (a *FormatterConfigAdapter) GetFormatIncrementalCreated() string {
	return a.cfg.FormatIncrementalCreated
}
func (a *FormatterConfigAdapter) GetFormatNoBackupsFound() string { return a.cfg.FormatNoBackupsFound }
func (a *FormatterConfigAdapter) GetFormatBackupWouldCreate() string {
	return a.cfg.FormatBackupWouldCreate
}
func (a *FormatterConfigAdapter) GetFormatBackupIdentical() string {
	return a.cfg.FormatBackupIdentical
}
func (a *FormatterConfigAdapter) GetFormatBackupCreated() string { return a.cfg.FormatBackupCreated }

func (a *FormatterConfigAdapter) GetFormatDiskFullError() string { return a.cfg.FormatDiskFullError }
func (a *FormatterConfigAdapter) GetFormatPermissionError() string {
	return a.cfg.FormatPermissionError
}
func (a *FormatterConfigAdapter) GetFormatDirectoryNotFound() string {
	return a.cfg.FormatDirectoryNotFound
}
func (a *FormatterConfigAdapter) GetFormatFileNotFound() string { return a.cfg.FormatFileNotFound }
func (a *FormatterConfigAdapter) GetFormatInvalidDirectory() string {
	return a.cfg.FormatInvalidDirectory
}
func (a *FormatterConfigAdapter) GetFormatInvalidFile() string { return a.cfg.FormatInvalidFile }
func (a *FormatterConfigAdapter) GetFormatFailedWriteTemp() string {
	return a.cfg.FormatFailedWriteTemp
}
func (a *FormatterConfigAdapter) GetFormatFailedFinalizeFile() string {
	return a.cfg.FormatFailedFinalizeFile
}
func (a *FormatterConfigAdapter) GetFormatFailedCreateDirDisk() string {
	return a.cfg.FormatFailedCreateDirDisk
}
func (a *FormatterConfigAdapter) GetFormatFailedCreateDirPerm() string {
	return a.cfg.FormatFailedCreateDir
}
func (a *FormatterConfigAdapter) GetFormatFailedRenamePerm() string {
	return a.cfg.FormatFailedAccessDir
}
func (a *FormatterConfigAdapter) GetFormatFailedRenameFile() string {
	return a.cfg.FormatFailedAccessFile
}

func (a *FormatterConfigAdapter) GetTemplateCreatedArchive() string {
	return a.cfg.TemplateCreatedArchive
}
func (a *FormatterConfigAdapter) GetTemplateIdenticalArchive() string {
	return a.cfg.TemplateIdenticalArchive
}
func (a *FormatterConfigAdapter) GetTemplateListArchive() string { return a.cfg.TemplateListArchive }
func (a *FormatterConfigAdapter) GetTemplateConfigValue() string { return a.cfg.TemplateConfigValue }
func (a *FormatterConfigAdapter) GetTemplateDryRunArchive() string {
	return a.cfg.TemplateDryRunArchive
}
func (a *FormatterConfigAdapter) GetTemplateError() string { return a.cfg.TemplateError }
func (a *FormatterConfigAdapter) GetTemplateCreatedBackup() string {
	return a.cfg.TemplateCreatedBackup
}
func (a *FormatterConfigAdapter) GetTemplateIdenticalBackup() string {
	return a.cfg.TemplateIdenticalBackup
}
func (a *FormatterConfigAdapter) GetTemplateListBackup() string   { return a.cfg.TemplateListBackup }
func (a *FormatterConfigAdapter) GetTemplateDryRunBackup() string { return a.cfg.TemplateDryRunBackup }

func (a *FormatterConfigAdapter) GetArchiveFilenamePattern() string {
	return a.cfg.PatternArchiveFilename
}
func (a *FormatterConfigAdapter) GetBackupFilenamePattern() string {
	return a.cfg.PatternBackupFilename
}
func (a *FormatterConfigAdapter) GetConfigLinePattern() string { return a.cfg.PatternConfigLine }
func (a *FormatterConfigAdapter) GetTimestampPattern() string  { return a.cfg.PatternTimestamp }

// 🔶 REFACTOR-005: Structure optimization - Git configuration adapter - 🔧
// GitConfigAdapter adapts Config struct to GitConfigInterface
type GitConfigAdapter struct {
	cfg *Config
}

func (a *GitConfigAdapter) GetIncludeGitInfo() bool     { return a.cfg.IncludeGitInfo }
func (a *GitConfigAdapter) GetShowGitDirtyStatus() bool { return a.cfg.ShowGitDirtyStatus }

// 🔶 REFACTOR-005: Structure optimization - Error configuration adapter - 🔧
// ErrorConfigAdapter adapts Config struct to ErrorConfigInterface
type ErrorConfigAdapter struct {
	cfg *Config
}

func (a *ErrorConfigAdapter) GetStatusCodes() map[string]int {
	return map[string]int{
		"created_archive":                      a.cfg.StatusCreatedArchive,
		"failed_to_create_archive_directory":   a.cfg.StatusFailedToCreateArchiveDirectory,
		"directory_is_identical_to_existing":   a.cfg.StatusDirectoryIsIdenticalToExistingArchive,
		"directory_not_found":                  a.cfg.StatusDirectoryNotFound,
		"invalid_directory_type":               a.cfg.StatusInvalidDirectoryType,
		"permission_denied":                    a.cfg.StatusPermissionDenied,
		"disk_full":                            a.cfg.StatusDiskFull,
		"config_error":                         a.cfg.StatusConfigError,
		"created_backup":                       a.cfg.StatusCreatedBackup,
		"failed_to_create_backup_directory":    a.cfg.StatusFailedToCreateBackupDirectory,
		"file_is_identical_to_existing_backup": a.cfg.StatusFileIsIdenticalToExistingBackup,
		"file_not_found":                       a.cfg.StatusFileNotFound,
		"invalid_file_type":                    a.cfg.StatusInvalidFileType,
	}
}

func (a *ErrorConfigAdapter) GetErrorFormatStrings() map[string]string {
	return map[string]string{
		"disk_full":              a.cfg.FormatDiskFullError,
		"permission":             a.cfg.FormatPermissionError,
		"directory_not_found":    a.cfg.FormatDirectoryNotFound,
		"file_not_found":         a.cfg.FormatFileNotFound,
		"invalid_directory":      a.cfg.FormatInvalidDirectory,
		"invalid_file":           a.cfg.FormatInvalidFile,
		"failed_write_temp":      a.cfg.FormatFailedWriteTemp,
		"failed_finalize":        a.cfg.FormatFailedFinalizeFile,
		"failed_create_dir_disk": a.cfg.FormatFailedCreateDirDisk,
		"failed_create_dir_perm": a.cfg.FormatFailedCreateDir,
		"failed_rename_perm":     a.cfg.FormatFailedAccessDir,
		"failed_rename_file":     a.cfg.FormatFailedAccessFile,
	}
}

func (a *ErrorConfigAdapter) GetDirectoryPermissions() os.FileMode { return 0755 } // Default directory permissions
func (a *ErrorConfigAdapter) GetFilePermissions() os.FileMode      { return 0644 } // Default file permissions
func (a *ErrorConfigAdapter) GetStatusFileNotFound() int           { return a.cfg.StatusFileNotFound }
func (a *ErrorConfigAdapter) GetStatusPermissionDenied() int       { return a.cfg.StatusPermissionDenied }
func (a *ErrorConfigAdapter) GetStatusInvalidFileType() int        { return a.cfg.StatusInvalidFileType }
func (a *ErrorConfigAdapter) GetStatusDiskFull() int               { return a.cfg.StatusDiskFull }
func (a *ErrorConfigAdapter) GetStatusDirectoryNotFound() int      { return a.cfg.StatusDirectoryNotFound }
func (a *ErrorConfigAdapter) GetStatusConfigError() int            { return a.cfg.StatusConfigError }

// 🔶 REFACTOR-005: Structure optimization - Shortened adapter names - 🔧
// ArchiveConfigAdapter adapts Config struct to ArchiveConfigInterface (renamed from ConfigToArchiveConfigAdapter)
type ArchiveConfigAdapter struct {
	cfg *Config
}

func (a *ArchiveConfigAdapter) GetArchiveDirPath() string            { return a.cfg.ArchiveDirPath }
func (a *ArchiveConfigAdapter) GetUseCurrentDirName() bool           { return a.cfg.UseCurrentDirName }
func (a *ArchiveConfigAdapter) GetExcludePatterns() []string         { return a.cfg.ExcludePatterns }
func (a *ArchiveConfigAdapter) GetIncludeGitInfo() bool              { return a.cfg.IncludeGitInfo }
func (a *ArchiveConfigAdapter) GetShowGitDirtyStatus() bool          { return a.cfg.ShowGitDirtyStatus }
func (a *ArchiveConfigAdapter) GetVerification() *VerificationConfig { return a.cfg.Verification }
func (a *ArchiveConfigAdapter) GetStatusCodes() map[string]int {
	return map[string]int{
		"created_archive":                    a.cfg.StatusCreatedArchive,
		"failed_to_create_archive_directory": a.cfg.StatusFailedToCreateArchiveDirectory,
		"directory_is_identical_to_existing": a.cfg.StatusDirectoryIsIdenticalToExistingArchive,
		"directory_not_found":                a.cfg.StatusDirectoryNotFound,
		"invalid_directory_type":             a.cfg.StatusInvalidDirectoryType,
		"permission_denied":                  a.cfg.StatusPermissionDenied,
		"disk_full":                          a.cfg.StatusDiskFull,
		"config_error":                       a.cfg.StatusConfigError,
	}
}
func (a *ArchiveConfigAdapter) GetStatusDirectoryNotFound() int { return a.cfg.StatusDirectoryNotFound }
func (a *ArchiveConfigAdapter) GetStatusDiskFull() int          { return a.cfg.StatusDiskFull }
func (a *ArchiveConfigAdapter) GetStatusConfigError() int       { return a.cfg.StatusConfigError }

// 🔶 REFACTOR-005: Structure optimization - Shortened adapter names - 🔧
// BackupConfigAdapter adapts Config struct to BackupConfigInterface (renamed from ConfigToBackupConfigAdapter)
type BackupConfigAdapter struct {
	cfg *Config
}

func (a *BackupConfigAdapter) GetBackupDirPath() string { return a.cfg.BackupDirPath }
func (a *BackupConfigAdapter) GetUseCurrentDirNameForFiles() bool {
	return a.cfg.UseCurrentDirNameForFiles
}
func (a *BackupConfigAdapter) GetStatusFileNotFound() int     { return a.cfg.StatusFileNotFound }
func (a *BackupConfigAdapter) GetStatusPermissionDenied() int { return a.cfg.StatusPermissionDenied }
func (a *BackupConfigAdapter) GetStatusInvalidFileType() int  { return a.cfg.StatusInvalidFileType }
func (a *BackupConfigAdapter) GetStatusDiskFull() int         { return a.cfg.StatusDiskFull }
func (a *BackupConfigAdapter) GetStatusFileIsIdenticalToExistingBackup() int {
	return a.cfg.StatusFileIsIdenticalToExistingBackup
}

// ========================================
// FORMATTER ADAPTERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Output formatter adapter - 🔧
// OutputFormatterAdapter adapts OutputFormatter struct to OutputFormatterInterface
type OutputFormatterAdapter struct {
	formatter *OutputFormatter
}

// NewOutputFormatterAdapter creates a new OutputFormatterAdapter
func NewOutputFormatterAdapter(formatter *OutputFormatter) *OutputFormatterAdapter {
	return &OutputFormatterAdapter{formatter: formatter}
}

// Core formatting methods
func (a *OutputFormatterAdapter) FormatCreatedArchive(path string) string {
	return a.formatter.FormatCreatedArchive(path)
}
func (a *OutputFormatterAdapter) FormatIdenticalArchive(path string) string {
	return a.formatter.FormatIdenticalArchive(path)
}
func (a *OutputFormatterAdapter) FormatListArchive(path, creationTime string) string {
	return a.formatter.FormatListArchive(path, creationTime)
}
func (a *OutputFormatterAdapter) FormatConfigValue(name, value, source string) string {
	return a.formatter.FormatConfigValue(name, value, source)
}
func (a *OutputFormatterAdapter) FormatDryRunArchive(path string) string {
	return a.formatter.FormatDryRunArchive(path)
}
func (a *OutputFormatterAdapter) FormatError(message string) string {
	return a.formatter.FormatError(message)
}

// Core printing methods
func (a *OutputFormatterAdapter) PrintCreatedArchive(path string) {
	a.formatter.PrintCreatedArchive(path)
}
func (a *OutputFormatterAdapter) PrintIdenticalArchive(path string) {
	a.formatter.PrintIdenticalArchive(path)
}
func (a *OutputFormatterAdapter) PrintListArchive(path, creationTime string) {
	a.formatter.PrintListArchive(path, creationTime)
}
func (a *OutputFormatterAdapter) PrintConfigValue(name, value, source string) {
	a.formatter.PrintConfigValue(name, value, source)
}
func (a *OutputFormatterAdapter) PrintDryRunArchive(path string) {
	a.formatter.PrintDryRunArchive(path)
}
func (a *OutputFormatterAdapter) PrintError(message string) { a.formatter.PrintError(message) }

// Backup formatting methods
func (a *OutputFormatterAdapter) FormatCreatedBackup(path string) string {
	return a.formatter.FormatCreatedBackup(path)
}
func (a *OutputFormatterAdapter) FormatIdenticalBackup(path string) string {
	return a.formatter.FormatIdenticalBackup(path)
}
func (a *OutputFormatterAdapter) FormatListBackup(path, creationTime string) string {
	return a.formatter.FormatListBackup(path, creationTime)
}
func (a *OutputFormatterAdapter) FormatDryRunBackup(path string) string {
	return a.formatter.FormatDryRunBackup(path)
}

// Backup printing methods
func (a *OutputFormatterAdapter) PrintCreatedBackup(path string) {
	a.formatter.PrintCreatedBackup(path)
}
func (a *OutputFormatterAdapter) PrintIdenticalBackup(path string) {
	a.formatter.PrintIdenticalBackup(path)
}
func (a *OutputFormatterAdapter) PrintListBackup(path, creationTime string) {
	a.formatter.PrintListBackup(path, creationTime)
}
func (a *OutputFormatterAdapter) PrintDryRunBackup(path string) { a.formatter.PrintDryRunBackup(path) }

// Extended formatting methods
func (a *OutputFormatterAdapter) FormatNoArchivesFound(archiveDir string) string {
	return a.formatter.FormatNoArchivesFound(archiveDir)
}
func (a *OutputFormatterAdapter) FormatVerificationFailed(archiveName string, err error) string {
	return a.formatter.FormatVerificationFailed(archiveName, err)
}
func (a *OutputFormatterAdapter) FormatVerificationSuccess(archiveName string) string {
	return a.formatter.FormatVerificationSuccess(archiveName)
}
func (a *OutputFormatterAdapter) FormatVerificationWarning(archiveName string, err error) string {
	return a.formatter.FormatVerificationWarning(archiveName, err)
}
func (a *OutputFormatterAdapter) FormatConfigurationUpdated(key string, value interface{}) string {
	return a.formatter.FormatConfigurationUpdated(key, value)
}
func (a *OutputFormatterAdapter) FormatConfigFilePath(path string) string {
	return a.formatter.FormatConfigFilePath(path)
}
func (a *OutputFormatterAdapter) FormatDryRunFilesHeader() string {
	return a.formatter.FormatDryRunFilesHeader()
}
func (a *OutputFormatterAdapter) FormatDryRunFileEntry(file string) string {
	return a.formatter.FormatDryRunFileEntry(file)
}
func (a *OutputFormatterAdapter) FormatNoFilesModified() string {
	return a.formatter.FormatNoFilesModified()
}
func (a *OutputFormatterAdapter) FormatIncrementalCreated(path string) string {
	return a.formatter.FormatIncrementalCreated(path)
}

// Extended printing methods
func (a *OutputFormatterAdapter) PrintNoArchivesFound(archiveDir string) {
	a.formatter.PrintNoArchivesFound(archiveDir)
}
func (a *OutputFormatterAdapter) PrintVerificationFailed(archiveName string, err error) {
	a.formatter.PrintVerificationFailed(archiveName, err)
}
func (a *OutputFormatterAdapter) PrintVerificationSuccess(archiveName string) {
	a.formatter.PrintVerificationSuccess(archiveName)
}
func (a *OutputFormatterAdapter) PrintVerificationWarning(archiveName string, err error) {
	a.formatter.PrintVerificationWarning(archiveName, err)
}
func (a *OutputFormatterAdapter) PrintConfigurationUpdated(key string, value interface{}) {
	a.formatter.PrintConfigurationUpdated(key, value)
}
func (a *OutputFormatterAdapter) PrintConfigFilePath(path string) {
	a.formatter.PrintConfigFilePath(path)
}
func (a *OutputFormatterAdapter) PrintDryRunFilesHeader() { a.formatter.PrintDryRunFilesHeader() }
func (a *OutputFormatterAdapter) PrintDryRunFileEntry(file string) {
	a.formatter.PrintDryRunFileEntry(file)
}
func (a *OutputFormatterAdapter) PrintNoFilesModified() { a.formatter.PrintNoFilesModified() }
func (a *OutputFormatterAdapter) PrintIncrementalCreated(path string) {
	a.formatter.PrintIncrementalCreated(path)
}

// Backup extended methods
func (a *OutputFormatterAdapter) FormatNoBackupsFound(filename, backupDir string) string {
	return a.formatter.FormatNoBackupsFound(filename, backupDir)
}
func (a *OutputFormatterAdapter) FormatBackupWouldCreate(path string) string {
	return a.formatter.FormatBackupWouldCreate(path)
}
func (a *OutputFormatterAdapter) FormatBackupIdentical(path string) string {
	return a.formatter.FormatBackupIdentical(path)
}
func (a *OutputFormatterAdapter) FormatBackupCreated(path string) string {
	return a.formatter.FormatBackupCreated(path)
}
func (a *OutputFormatterAdapter) PrintNoBackupsFound(filename, backupDir string) {
	a.formatter.PrintNoBackupsFound(filename, backupDir)
}
func (a *OutputFormatterAdapter) PrintBackupWouldCreate(path string) {
	a.formatter.PrintBackupWouldCreate(path)
}
func (a *OutputFormatterAdapter) PrintBackupIdentical(path string) {
	a.formatter.PrintBackupIdentical(path)
}
func (a *OutputFormatterAdapter) PrintBackupCreated(path string) {
	a.formatter.PrintBackupCreated(path)
}

// Error formatting methods
func (a *OutputFormatterAdapter) FormatDiskFullError(err error) string {
	return a.formatter.FormatDiskFullError(err)
}
func (a *OutputFormatterAdapter) FormatPermissionError(err error) string {
	return a.formatter.FormatPermissionError(err)
}
func (a *OutputFormatterAdapter) FormatDirectoryNotFound(err error) string {
	return a.formatter.FormatDirectoryNotFound(err)
}
func (a *OutputFormatterAdapter) FormatFileNotFound(err error) string {
	return a.formatter.FormatFileNotFound(err)
}
func (a *OutputFormatterAdapter) FormatInvalidDirectory(err error) string {
	return a.formatter.FormatInvalidDirectory(err)
}
func (a *OutputFormatterAdapter) FormatInvalidFile(err error) string {
	return a.formatter.FormatInvalidFile(err)
}
func (a *OutputFormatterAdapter) FormatFailedWriteTemp(err error) string {
	return a.formatter.FormatFailedWriteTemp(err)
}
func (a *OutputFormatterAdapter) FormatFailedFinalizeFile(err error) string {
	return a.formatter.FormatFailedFinalizeFile(err)
}
func (a *OutputFormatterAdapter) FormatFailedCreateDirDisk(err error) string {
	return a.formatter.FormatFailedCreateDirDisk(err)
}
func (a *OutputFormatterAdapter) FormatFailedCreateDirPerm(err error) string {
	return a.formatter.FormatFailedCreateDir(err)
}
func (a *OutputFormatterAdapter) FormatFailedRenamePerm(err error) string {
	return a.formatter.FormatFailedAccessDir(err)
}
func (a *OutputFormatterAdapter) FormatFailedRenameFile(err error) string {
	return a.formatter.FormatFailedAccessFile(err)
}

// Error printing methods
func (a *OutputFormatterAdapter) PrintDiskFullError(err error) { a.formatter.PrintDiskFullError(err) }
func (a *OutputFormatterAdapter) PrintPermissionError(err error) {
	a.formatter.PrintPermissionError(err)
}
func (a *OutputFormatterAdapter) PrintDirectoryNotFound(err error) {
	a.formatter.PrintDirectoryNotFound(err)
}
func (a *OutputFormatterAdapter) PrintFileNotFound(err error) { a.formatter.PrintFileNotFound(err) }
func (a *OutputFormatterAdapter) PrintInvalidDirectory(err error) {
	a.formatter.PrintInvalidDirectory(err)
}
func (a *OutputFormatterAdapter) PrintInvalidFile(err error) { a.formatter.PrintInvalidFile(err) }
func (a *OutputFormatterAdapter) PrintFailedWriteTemp(err error) {
	a.formatter.PrintFailedWriteTemp(err)
}
func (a *OutputFormatterAdapter) PrintFailedFinalizeFile(err error) {
	a.formatter.PrintFailedFinalizeFile(err)
}
func (a *OutputFormatterAdapter) PrintFailedCreateDirDisk(err error) {
	a.formatter.PrintFailedCreateDirDisk(err)
}
func (a *OutputFormatterAdapter) PrintFailedCreateDirPerm(err error) {
	a.formatter.PrintFailedCreateDir(err)
}
func (a *OutputFormatterAdapter) PrintFailedRenamePerm(err error) {
	a.formatter.PrintFailedAccessDir(err)
}
func (a *OutputFormatterAdapter) PrintFailedRenameFile(err error) {
	a.formatter.PrintFailedAccessFile(err)
}

// Template methods - delegate to TemplateFormatter
func (a *OutputFormatterAdapter) TemplateCreatedArchive(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateCreatedArchive(data)
}
func (a *OutputFormatterAdapter) TemplateIdenticalArchive(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateIdenticalArchive(data)
}
func (a *OutputFormatterAdapter) TemplateListArchive(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateListArchive(data)
}
func (a *OutputFormatterAdapter) TemplateConfigValue(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateConfigValue(data)
}
func (a *OutputFormatterAdapter) TemplateDryRunArchive(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateDryRunArchive(data)
}
func (a *OutputFormatterAdapter) TemplateError(data map[string]string) string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.TemplateError(data)
}

// Template printing methods - delegate to TemplateFormatter
func (a *OutputFormatterAdapter) PrintTemplateCreatedArchive(path string) {
	tf := NewTemplateFormatter(a.formatter.cfg)
	tf.PrintTemplateCreatedArchive(path)
}
func (a *OutputFormatterAdapter) PrintTemplateCreatedBackup(path string) {
	tf := NewTemplateFormatter(a.formatter.cfg)
	tf.PrintTemplateCreatedBackup(path)
}
func (a *OutputFormatterAdapter) PrintTemplateListBackup(path, creationTime string) {
	tf := NewTemplateFormatter(a.formatter.cfg)
	tf.PrintTemplateListBackup(path, creationTime)
}
func (a *OutputFormatterAdapter) PrintTemplateError(message, operation string) {
	tf := NewTemplateFormatter(a.formatter.cfg)
	tf.PrintTemplateError(message, operation)
}

// Pattern extraction methods - delegate to TemplateFormatter
func (a *OutputFormatterAdapter) ExtractArchiveFilenameData(filename string) map[string]string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.extractArchiveData(filename)
}
func (a *OutputFormatterAdapter) ExtractBackupFilenameData(filename string) map[string]string {
	tf := NewTemplateFormatter(a.formatter.cfg)
	return tf.extractBackupData(filename)
}
func (a *OutputFormatterAdapter) ExtractConfigLineData(line string) map[string]string {
	return a.formatter.ExtractConfigLineData(line)
}
func (a *OutputFormatterAdapter) ExtractTimestampData(timestamp string) map[string]string {
	return a.formatter.ExtractTimestampData(timestamp)
}

// Output control methods
func (a *OutputFormatterAdapter) IsDelayedMode() bool { return a.formatter.IsDelayedMode() }
func (a *OutputFormatterAdapter) SetCollector(collector *OutputCollector) {
	a.formatter.SetCollector(collector)
}
func (a *OutputFormatterAdapter) GetCollector() *OutputCollector { return a.formatter.GetCollector() }

// ========================================
// RESOURCE MANAGEMENT ADAPTERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Resource manager factory implementation - 🔧
// DefaultResourceManagerFactory implements ResourceManagerFactoryInterface
type DefaultResourceManagerFactory struct{}

func (f *DefaultResourceManagerFactory) CreateResourceManager() ResourceManagerInterface {
	return NewResourceManager()
}

func (f *DefaultResourceManagerFactory) CreateResourceManagerWithContext(ctx context.Context) ResourceManagerInterface {
	// Create resource manager with context - context is stored separately
	return NewResourceManager()
}

// ========================================
// GIT PROVIDER ADAPTERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Git provider adapter - 🔧
// DefaultGitProvider implements GitProviderInterface using existing Git functions
type DefaultGitProvider struct{}

func (g *DefaultGitProvider) IsGitRepository(dir string) bool {
	return IsGitRepository(dir)
}

func (g *DefaultGitProvider) GetCurrentBranch(dir string) (string, error) {
	branch := GetGitBranch(dir)
	if branch == "" {
		return "", &GitError{Operation: "get-branch", Err: fmt.Errorf("not in a git repository")}
	}
	return branch, nil
}

func (g *DefaultGitProvider) GetCurrentCommitHash(dir string) (string, error) {
	// Use GetGitInfo to get the hash
	_, hash := GetGitInfo(dir)
	if hash == "" {
		return "", &GitError{Operation: "get-hash", Err: fmt.Errorf("not in a git repository")}
	}
	return hash, nil
}

func (g *DefaultGitProvider) IsWorkingTreeClean(dir string) (bool, error) {
	return IsGitWorkingDirectoryClean(dir), nil
}

func (g *DefaultGitProvider) GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	return GetGitInfoWithStatus(dir)
}

func (g *DefaultGitProvider) ExecuteGitCommand(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", &GitError{Operation: "execute", Err: err}
	}
	return strings.TrimSpace(string(out)), nil
}

// ========================================
// EXTERNAL DEPENDENCY ADAPTERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - YAML processor implementation - 🔧
// DefaultYAMLProcessor implements YAMLProcessorInterface using gopkg.in/yaml.v3
type DefaultYAMLProcessor struct{}

func (y *DefaultYAMLProcessor) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *DefaultYAMLProcessor) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

// 🔶 REFACTOR-005: Structure optimization - Filesystem implementation - 🔧
// DefaultFilesystem implements FilesystemInterface using standard library
type DefaultFilesystem struct{}

func (f *DefaultFilesystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (f *DefaultFilesystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (f *DefaultFilesystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (f *DefaultFilesystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (f *DefaultFilesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// 🔶 REFACTOR-005: Structure optimization - Command executor implementation - 🔧
// DefaultCommandExecutor implements CommandExecutorInterface using os/exec
type DefaultCommandExecutor struct{}

func (e *DefaultCommandExecutor) Execute(dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Output()
}

func (e *DefaultCommandExecutor) ExecuteWithContext(ctx context.Context, dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Output()
}

// ========================================
// SERVICE PROVIDER IMPLEMENTATION
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Universal service provider implementation - 🔧
// DefaultServiceProvider implements ServiceProviderInterface with default implementations
type DefaultServiceProvider struct {
	config           *Config
	formatter        *OutputFormatter
	resourceFactory  ResourceManagerFactoryInterface
	gitProvider      GitProviderInterface
	fileOperations   FileOperationsInterface
	exclusionManager ExclusionManagerInterface
	yamlProcessor    YAMLProcessorInterface
	filesystem       FilesystemInterface
	commandExecutor  CommandExecutorInterface
}

// NewDefaultServiceProvider creates a comprehensive service provider with all components
func NewDefaultServiceProvider(config *Config, formatter *OutputFormatter) *DefaultServiceProvider {
	return &DefaultServiceProvider{
		config:           config,
		formatter:        formatter,
		resourceFactory:  &DefaultResourceManagerFactory{},
		gitProvider:      &DefaultGitProvider{},
		fileOperations:   &DefaultFileOperations{},
		exclusionManager: NewDefaultExclusionManager(config.ExcludePatterns),
		yamlProcessor:    &DefaultYAMLProcessor{},
		filesystem:       &DefaultFilesystem{},
		commandExecutor:  &DefaultCommandExecutor{},
	}
}

func (s *DefaultServiceProvider) GetConfigProvider() ConfigProviderInterface {
	return NewConfigProviderAdapter(s.config)
}

func (s *DefaultServiceProvider) GetFormatterProvider() OutputFormatterInterface {
	return NewOutputFormatterAdapter(s.formatter)
}

func (s *DefaultServiceProvider) GetResourceManagerFactory() ResourceManagerFactoryInterface {
	return s.resourceFactory
}

func (s *DefaultServiceProvider) GetGitProvider() GitProviderInterface {
	return s.gitProvider
}

func (s *DefaultServiceProvider) GetFileOperations() FileOperationsInterface {
	return s.fileOperations
}

func (s *DefaultServiceProvider) GetExclusionManager() ExclusionManagerInterface {
	return s.exclusionManager
}

func (s *DefaultServiceProvider) GetYAMLProcessor() YAMLProcessorInterface {
	return s.yamlProcessor
}

func (s *DefaultServiceProvider) GetFilesystem() FilesystemInterface {
	return s.filesystem
}

func (s *DefaultServiceProvider) GetCommandExecutor() CommandExecutorInterface {
	return s.commandExecutor
}

// ========================================
// MISSING ADAPTER IMPLEMENTATIONS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - File operations adapter - 🔧
// DefaultFileOperations provides default implementation of FileOperationsInterface
type DefaultFileOperations struct{}

func (f *DefaultFileOperations) SafeMkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *DefaultFileOperations) ValidateDirectoryPath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("cannot access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", path)
	}
	return nil
}

func (f *DefaultFileOperations) CopyFile(src, dst string) error {
	return f.CopyFileWithContext(context.Background(), src, dst)
}

func (f *DefaultFileOperations) CopyFileWithContext(ctx context.Context, src, dst string) error {
	// Implementation using standard library functions
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Use context-aware copy if available
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err = destFile.ReadFrom(sourceFile)
		return err
	}
}

func (f *DefaultFileOperations) ValidateFilePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}
	return nil
}

func (f *DefaultFileOperations) CompareFiles(file1, file2 string) (bool, error) {
	// Simple implementation using file content comparison
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", file1, err)
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", file2, err)
	}

	return string(content1) == string(content2), nil
}

func (f *DefaultFileOperations) AtomicWriteFile(path string, data []byte, rm ResourceManagerInterface) error {
	// Create temporary file in same directory as target
	tmpFile := path + ".tmp"

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Add cleanup to resource manager
	if rm != nil {
		rm.AddTempFile(tmpFile)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, path); err != nil {
		os.Remove(tmpFile) // Cleanup on failure
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

func (f *DefaultFileOperations) AtomicCopyFile(src, dst string, rm ResourceManagerInterface) error {
	// Create temporary file in same directory as destination
	tmpFile := dst + ".tmp"

	if err := f.CopyFile(src, tmpFile); err != nil {
		return fmt.Errorf("failed to copy to temporary file: %w", err)
	}

	// Add cleanup to resource manager
	if rm != nil {
		rm.AddTempFile(tmpFile)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, dst); err != nil {
		os.Remove(tmpFile) // Cleanup on failure
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}

// 🔶 REFACTOR-005: Structure optimization - Exclusion manager adapter - 🔧
// DefaultExclusionManager provides default implementation of ExclusionManagerInterface
type DefaultExclusionManager struct {
	patterns []string
}

func NewDefaultExclusionManager(patterns []string) *DefaultExclusionManager {
	return &DefaultExclusionManager{patterns: patterns}
}

func (e *DefaultExclusionManager) ShouldExclude(path string, excludePatterns []string) bool {
	// Use provided patterns or fall back to configured patterns
	patterns := excludePatterns
	if len(patterns) == 0 {
		patterns = e.patterns
	}

	for _, pattern := range patterns {
		// Simple pattern matching - in production would use filepath.Match or similar
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func (e *DefaultExclusionManager) FilterFiles(files []string, excludePatterns []string) []string {
	var filtered []string
	for _, file := range files {
		if !e.ShouldExclude(file, excludePatterns) {
			filtered = append(filtered, file)
		}
	}
	return filtered
}

func (e *DefaultExclusionManager) LoadExclusionPatterns() ([]string, error) {
	// Return current patterns - in production might load from file
	return e.patterns, nil
}

func (e *DefaultExclusionManager) AddExclusionPattern(pattern string) {
	e.patterns = append(e.patterns, pattern)
}

func (e *DefaultExclusionManager) RemoveExclusionPattern(pattern string) {
	for i, p := range e.patterns {
		if p == pattern {
			e.patterns = append(e.patterns[:i], e.patterns[i+1:]...)
			return
		}
	}
}

func (e *DefaultExclusionManager) ValidatePattern(pattern string) error {
	// Basic validation - in production would use filepath.Match
	if pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	return nil
}
