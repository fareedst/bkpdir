// This file is part of bkpdir
//
// Package main provides comprehensive interface definitions for component extraction.
// These interfaces enable clean separation and extraction of components while maintaining
// backward compatibility.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"os"
)

// 🔶 REFACTOR-005: Structure optimization - Comprehensive interface definitions for extraction - 📝
// This file contains all interface definitions required for clean component extraction

// ========================================
// CONFIGURATION INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Universal configuration provider interface - 🔧
// ConfigProviderInterface provides centralized access to all configuration interfaces
type ConfigProviderInterface interface {
	GetArchiveConfig() ArchiveConfigInterface
	GetBackupConfig() BackupConfigInterface
	GetFormatterConfig() FormatterConfigInterface
	GetErrorConfig() ErrorConfigInterface
	GetGitConfig() GitConfigInterface
}

// 🔶 REFACTOR-005: Structure optimization - Comprehensive formatter configuration interface - 🔧
// FormatterConfigInterface abstracts all configuration dependencies for formatter components
type FormatterConfigInterface interface {
	// Printf-style format strings
	GetFormatCreatedArchive() string
	GetFormatIdenticalArchive() string
	GetFormatListArchive() string
	GetFormatConfigValue() string
	GetFormatDryRunArchive() string
	GetFormatError() string
	GetFormatCreatedBackup() string
	GetFormatIdenticalBackup() string
	GetFormatListBackup() string
	GetFormatDryRunBackup() string

	// Extended format strings
	GetFormatNoArchivesFound() string
	GetFormatVerificationFailed() string
	GetFormatVerificationSuccess() string
	GetFormatVerificationWarning() string
	GetFormatConfigurationUpdated() string
	GetFormatConfigFilePath() string
	GetFormatDryRunFilesHeader() string
	GetFormatDryRunFileEntry() string
	GetFormatNoFilesModified() string
	GetFormatIncrementalCreated() string
	GetFormatNoBackupsFound() string
	GetFormatBackupWouldCreate() string
	GetFormatBackupIdentical() string
	GetFormatBackupCreated() string

	// Error format strings
	GetFormatDiskFullError() string
	GetFormatPermissionError() string
	GetFormatDirectoryNotFound() string
	GetFormatFileNotFound() string
	GetFormatInvalidDirectory() string
	GetFormatInvalidFile() string
	GetFormatFailedWriteTemp() string
	GetFormatFailedFinalizeFile() string
	GetFormatFailedCreateDirDisk() string
	GetFormatFailedCreateDirPerm() string
	GetFormatFailedRenamePerm() string
	GetFormatFailedRenameFile() string

	// Template strings
	GetTemplateCreatedArchive() string
	GetTemplateIdenticalArchive() string
	GetTemplateListArchive() string
	GetTemplateConfigValue() string
	GetTemplateDryRunArchive() string
	GetTemplateError() string
	GetTemplateCreatedBackup() string
	GetTemplateIdenticalBackup() string
	GetTemplateListBackup() string
	GetTemplateDryRunBackup() string

	// Pattern strings for regex extraction
	GetArchiveFilenamePattern() string
	GetBackupFilenamePattern() string
	GetConfigLinePattern() string
	GetTimestampPattern() string
}

// 🔶 REFACTOR-005: Structure optimization - Git configuration interface - 🔧
// GitConfigInterface abstracts Git-related configuration
type GitConfigInterface interface {
	GetIncludeGitInfo() bool
	GetShowGitDirtyStatus() bool
}

// 🔶 REFACTOR-005: Structure optimization - Renamed for consistency - 🔧
// ErrorConfigInterface abstracts error handling configuration (renamed from ErrorConfig)
type ErrorConfigInterface interface {
	GetStatusCodes() map[string]int
	GetErrorFormatStrings() map[string]string
	GetDirectoryPermissions() os.FileMode
	GetFilePermissions() os.FileMode
	GetStatusFileNotFound() int
	GetStatusPermissionDenied() int
	GetStatusInvalidFileType() int
	GetStatusDiskFull() int
	GetStatusDirectoryNotFound() int
	GetStatusConfigError() int
}

// ========================================
// FORMATTER INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Complete output formatter interface - 🔧
// OutputFormatterInterface abstracts all formatter dependencies for clean extraction
type OutputFormatterInterface interface {
	// Core formatting methods
	FormatCreatedArchive(path string) string
	FormatIdenticalArchive(path string) string
	FormatListArchive(path, creationTime string) string
	FormatConfigValue(name, value, source string) string
	FormatDryRunArchive(path string) string
	FormatError(message string) string

	// Core printing methods
	PrintCreatedArchive(path string)
	PrintIdenticalArchive(path string)
	PrintListArchive(path, creationTime string)
	PrintConfigValue(name, value, source string)
	PrintDryRunArchive(path string)
	PrintError(message string)

	// Backup formatting methods
	FormatCreatedBackup(path string) string
	FormatIdenticalBackup(path string) string
	FormatListBackup(path, creationTime string) string
	FormatDryRunBackup(path string) string

	// Backup printing methods
	PrintCreatedBackup(path string)
	PrintIdenticalBackup(path string)
	PrintListBackup(path, creationTime string)
	PrintDryRunBackup(path string)

	// Extended formatting methods
	FormatNoArchivesFound(archiveDir string) string
	FormatVerificationFailed(archiveName string, err error) string
	FormatVerificationSuccess(archiveName string) string
	FormatVerificationWarning(archiveName string, err error) string
	FormatConfigurationUpdated(key string, value interface{}) string
	FormatConfigFilePath(path string) string
	FormatDryRunFilesHeader() string
	FormatDryRunFileEntry(file string) string
	FormatNoFilesModified() string
	FormatIncrementalCreated(path string) string

	// Extended printing methods
	PrintNoArchivesFound(archiveDir string)
	PrintVerificationFailed(archiveName string, err error)
	PrintVerificationSuccess(archiveName string)
	PrintVerificationWarning(archiveName string, err error)
	PrintConfigurationUpdated(key string, value interface{})
	PrintConfigFilePath(path string)
	PrintDryRunFilesHeader()
	PrintDryRunFileEntry(file string)
	PrintNoFilesModified()
	PrintIncrementalCreated(path string)

	// Backup extended methods
	FormatNoBackupsFound(filename, backupDir string) string
	FormatBackupWouldCreate(path string) string
	FormatBackupIdentical(path string) string
	FormatBackupCreated(path string) string
	PrintNoBackupsFound(filename, backupDir string)
	PrintBackupWouldCreate(path string)
	PrintBackupIdentical(path string)
	PrintBackupCreated(path string)

	// Error formatting methods
	FormatDiskFullError(err error) string
	FormatPermissionError(err error) string
	FormatDirectoryNotFound(err error) string
	FormatFileNotFound(err error) string
	FormatInvalidDirectory(err error) string
	FormatInvalidFile(err error) string
	FormatFailedWriteTemp(err error) string
	FormatFailedFinalizeFile(err error) string
	FormatFailedCreateDirDisk(err error) string
	FormatFailedCreateDirPerm(err error) string
	FormatFailedRenamePerm(err error) string
	FormatFailedRenameFile(err error) string

	// Error printing methods
	PrintDiskFullError(err error)
	PrintPermissionError(err error)
	PrintDirectoryNotFound(err error)
	PrintFileNotFound(err error)
	PrintInvalidDirectory(err error)
	PrintInvalidFile(err error)
	PrintFailedWriteTemp(err error)
	PrintFailedFinalizeFile(err error)
	PrintFailedCreateDirDisk(err error)
	PrintFailedCreateDirPerm(err error)
	PrintFailedRenamePerm(err error)
	PrintFailedRenameFile(err error)

	// Template methods
	TemplateCreatedArchive(data map[string]string) string
	TemplateIdenticalArchive(data map[string]string) string
	TemplateListArchive(data map[string]string) string
	TemplateConfigValue(data map[string]string) string
	TemplateDryRunArchive(data map[string]string) string
	TemplateError(data map[string]string) string

	// Template printing methods
	PrintTemplateCreatedArchive(path string)
	PrintTemplateCreatedBackup(path string)
	PrintTemplateListBackup(path, creationTime string)
	PrintTemplateError(message, operation string)

	// Pattern extraction methods
	ExtractArchiveFilenameData(filename string) map[string]string
	ExtractBackupFilenameData(filename string) map[string]string
	ExtractConfigLineData(line string) map[string]string
	ExtractTimestampData(timestamp string) map[string]string

	// Output control methods
	IsDelayedMode() bool
	SetCollector(collector *OutputCollector)
	GetCollector() *OutputCollector
}

// ========================================
// RESOURCE MANAGEMENT INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Resource management factory interface - 🔧
// ResourceManagerFactoryInterface provides centralized resource manager creation
type ResourceManagerFactoryInterface interface {
	CreateResourceManager() ResourceManagerInterface
	CreateResourceManagerWithContext(ctx context.Context) ResourceManagerInterface
}

// ========================================
// GIT PROVIDER INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Renamed for consistency - 🔧
// GitProviderInterface provides Git operations (renamed for consistency)
type GitProviderInterface interface {
	IsGitRepository(dir string) bool
	GetCurrentBranch(dir string) (string, error)
	GetCurrentCommitHash(dir string) (string, error)
	IsWorkingTreeClean(dir string) (bool, error)
	GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool)
	ExecuteGitCommand(dir string, args ...string) (string, error)
}

// ========================================
// FILE OPERATIONS INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - File operations interface - 🔧
// FileOperationsInterface provides file and directory operations
type FileOperationsInterface interface {
	// Directory operations
	SafeMkdirAll(path string, perm os.FileMode) error
	ValidateDirectoryPath(path string) error

	// File operations
	CopyFile(src, dst string) error
	CopyFileWithContext(ctx context.Context, src, dst string) error
	ValidateFilePath(path string) error
	CompareFiles(file1, file2 string) (bool, error)

	// Atomic operations
	AtomicWriteFile(path string, data []byte, rm ResourceManagerInterface) error
	AtomicCopyFile(src, dst string, rm ResourceManagerInterface) error
}

// 🔶 REFACTOR-005: Structure optimization - Exclusion management interface - 🔧
// ExclusionManagerInterface provides file exclusion operations
type ExclusionManagerInterface interface {
	ShouldExclude(path string, excludePatterns []string) bool
	FilterFiles(files []string, excludePatterns []string) []string
	LoadExclusionPatterns() ([]string, error)
	AddExclusionPattern(pattern string)
	RemoveExclusionPattern(pattern string)
	ValidatePattern(pattern string) error
}

// ========================================
// EXTERNAL DEPENDENCY INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - YAML processor interface for dependency injection - 🔧
// YAMLProcessorInterface abstracts YAML operations for dependency injection
type YAMLProcessorInterface interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

// 🔶 REFACTOR-005: Structure optimization - Filesystem interface for dependency injection - 🔧
// FilesystemInterface abstracts filesystem operations for dependency injection
type FilesystemInterface interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	MkdirAll(path string, perm os.FileMode) error
}

// 🔶 REFACTOR-005: Structure optimization - Command executor interface for dependency injection - 🔧
// CommandExecutorInterface abstracts command execution for dependency injection
type CommandExecutorInterface interface {
	Execute(dir string, name string, args ...string) ([]byte, error)
	ExecuteWithContext(ctx context.Context, dir string, name string, args ...string) ([]byte, error)
}

// ========================================
// SERVICE PROVIDER INTERFACES
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Universal service provider interface - 🔧
// ServiceProviderInterface aggregates all component interfaces for easy dependency injection
type ServiceProviderInterface interface {
	GetConfigProvider() ConfigProviderInterface
	GetFormatterProvider() OutputFormatterInterface
	GetResourceManagerFactory() ResourceManagerFactoryInterface
	GetGitProvider() GitProviderInterface
	GetFileOperations() FileOperationsInterface
	GetExclusionManager() ExclusionManagerInterface
	GetYAMLProcessor() YAMLProcessorInterface
	GetFilesystem() FilesystemInterface
	GetCommandExecutor() CommandExecutorInterface
}
