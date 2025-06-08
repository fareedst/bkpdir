// Backward compatibility adapter for the extracted formatter package.
// This adapter bridges the original OutputFormatter with the new pkg/formatter
// components while maintaining all existing functionality and API compatibility.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
)

// ⭐ EXTRACT-003: Backward compatibility adapter - 🔧 Configuration provider implementation
// FormatterConfigProvider adapts the Config struct to the formatter.ConfigProvider interface
type FormatterConfigProvider struct {
	config *Config
}

// NewFormatterConfigProvider creates a new FormatterConfigProvider
func NewFormatterConfigProvider(config *Config) *FormatterConfigProvider {
	return &FormatterConfigProvider{config: config}
}

// GetFormatString returns format strings from the configuration
func (fcp *FormatterConfigProvider) GetFormatString(formatType string) string {
	switch formatType {
	case "created_archive":
		result := fcp.config.FormatCreatedArchive
		fmt.Fprintf(os.Stderr, "DEBUG: GetFormatString(created_archive) = %q\n", result)
		return result
	case "identical_archive":
		return fcp.config.FormatIdenticalArchive
	case "list_archive":
		return fcp.config.FormatListArchive
	case "config_value":
		return fcp.config.FormatConfigValue
	case "dry_run_archive":
		return fcp.config.FormatDryRunArchive
	case "created_backup":
		return fcp.config.FormatCreatedBackup
	case "identical_backup":
		return fcp.config.FormatIdenticalBackup
	case "list_backup":
		return fcp.config.FormatListBackup
	case "dry_run_backup":
		return fcp.config.FormatDryRunBackup
	default:
		return ""
	}
}

// GetTemplateString returns template strings from the configuration
func (fcp *FormatterConfigProvider) GetTemplateString(templateType string) string {
	switch templateType {
	case "created_archive":
		return fcp.config.TemplateCreatedArchive
	case "identical_archive":
		return fcp.config.TemplateIdenticalArchive
	case "list_archive":
		return fcp.config.TemplateListArchive
	case "config_value":
		return fcp.config.TemplateConfigValue
	case "dry_run_archive":
		return fcp.config.TemplateDryRunArchive
	case "error":
		return fcp.config.TemplateError
	case "created_backup":
		return fcp.config.TemplateCreatedBackup
	case "identical_backup":
		return fcp.config.TemplateIdenticalBackup
	case "list_backup":
		return fcp.config.TemplateListBackup
	case "dry_run_backup":
		return fcp.config.TemplateDryRunBackup
	case "error_disk_full":
		return fcp.config.TemplateDiskFullError
	case "error_permission":
		return fcp.config.TemplatePermissionError
	case "error_directory_not_found":
		return fcp.config.TemplateDirectoryNotFound
	case "error_file_not_found":
		return fcp.config.TemplateFileNotFound
	default:
		return ""
	}
}

// GetPattern returns regex patterns from the configuration
func (fcp *FormatterConfigProvider) GetPattern(patternType string) string {
	switch patternType {
	case "archive_filename":
		return fcp.config.PatternArchiveFilename
	case "backup_filename":
		return fcp.config.PatternBackupFilename
	case "config_line":
		return fcp.config.PatternConfigLine
	case "timestamp":
		return fcp.config.PatternTimestamp
	default:
		return ""
	}
}

// GetErrorFormat returns error format strings from the configuration
func (fcp *FormatterConfigProvider) GetErrorFormat(errorType string) string {
	switch errorType {
	case "generic":
		return fcp.config.FormatError
	case "disk_full":
		return fcp.config.FormatDiskFullError
	case "permission":
		return fcp.config.FormatPermissionError
	case "directory_not_found":
		return fcp.config.FormatDirectoryNotFound
	case "file_not_found":
		return fcp.config.FormatFileNotFound
	case "invalid_directory":
		return fcp.config.FormatInvalidDirectory
	case "invalid_file":
		return fcp.config.FormatInvalidFile
	default:
		return ""
	}
}

// ⭐ OUT-002: Enhanced output with file statistics - Enhanced configuration support
// GetDetailedFormatString returns detailed format strings with file statistics
func (fcp *FormatterConfigProvider) GetDetailedFormatString(formatType string) string {
	switch formatType {
	case "created_archive":
		return fcp.config.FormatCreatedArchiveDetailed
	case "incremental_created":
		return fcp.config.FormatIncrementalCreatedDetailed
	default:
		return ""
	}
}

// GetDetailedTemplateString returns detailed template strings with file statistics
func (fcp *FormatterConfigProvider) GetDetailedTemplateString(templateType string) string {
	switch templateType {
	case "created_archive":
		return fcp.config.TemplateCreatedArchiveDetailed
	case "incremental_created":
		return fcp.config.TemplateIncrementalCreatedDetailed
	default:
		return ""
	}
}

// ⭐ EXTRACT-003: Backward compatibility adapter - 🔧 OutputFormatter adapter
// FormatterAdapter wraps the extracted formatter to maintain backward compatibility
type FormatterAdapter struct {
	formatter formatter.OutputFormatterInterface
	config    *Config
}

// NewFormatterAdapter creates a new FormatterAdapter
func NewFormatterAdapter(config *Config) *FormatterAdapter {
	configProvider := NewFormatterConfigProvider(config)
	return &FormatterAdapter{
		formatter: formatter.NewDefaultOutputFormatter(configProvider),
		config:    config,
	}
}

// NewFormatterAdapterWithCollector creates an adapter with delayed output support
func NewFormatterAdapterWithCollector(config *Config, collector *formatter.OutputCollector) *FormatterAdapter {
	configProvider := NewFormatterConfigProvider(config)
	return &FormatterAdapter{
		formatter: formatter.NewDefaultOutputFormatterWithCollector(configProvider, collector),
		config:    config,
	}
}

// ⭐ EXTRACT-003: Backward compatibility adapter - 📝 Delegate all methods to extracted formatter

// IsDelayedMode delegates to the extracted formatter
func (fa *FormatterAdapter) IsDelayedMode() bool {
	return fa.formatter.IsDelayedMode()
}

// GetCollector delegates to the extracted formatter
func (fa *FormatterAdapter) GetCollector() *formatter.OutputCollector {
	return fa.formatter.GetCollector()
}

// SetCollector delegates to the extracted formatter
func (fa *FormatterAdapter) SetCollector(collector *formatter.OutputCollector) {
	fa.formatter.SetCollector(collector)
}

// Printf-style formatting operations
func (fa *FormatterAdapter) FormatCreatedArchive(path string) string {
	return fa.formatter.FormatCreatedArchive(path)
}

func (fa *FormatterAdapter) FormatIdenticalArchive(path string) string {
	return fa.formatter.FormatIdenticalArchive(path)
}

func (fa *FormatterAdapter) FormatListArchive(path, creationTime string) string {
	return fa.formatter.FormatListArchive(path, creationTime)
}

func (fa *FormatterAdapter) FormatConfigValue(name, value, source string) string {
	return fa.formatter.FormatConfigValue(name, value, source)
}

func (fa *FormatterAdapter) FormatError(message string) string {
	return fa.formatter.FormatError(message)
}

func (fa *FormatterAdapter) FormatDryRunArchive(path string) string {
	return fa.formatter.FormatDryRunArchive(path)
}

func (fa *FormatterAdapter) FormatCreatedBackup(path string) string {
	return fa.formatter.FormatCreatedBackup(path)
}

func (fa *FormatterAdapter) FormatIdenticalBackup(path string) string {
	return fa.formatter.FormatIdenticalBackup(path)
}

func (fa *FormatterAdapter) FormatListBackup(path, creationTime string) string {
	return fa.formatter.FormatListBackup(path, creationTime)
}

func (fa *FormatterAdapter) FormatDryRunBackup(path string) string {
	return fa.formatter.FormatDryRunBackup(path)
}

// Print operations
func (fa *FormatterAdapter) PrintCreatedArchive(path string) {
	fa.formatter.PrintCreatedArchive(path)
}

func (fa *FormatterAdapter) PrintIdenticalArchive(path string) {
	fa.formatter.PrintIdenticalArchive(path)
}

func (fa *FormatterAdapter) PrintListArchive(path, creationTime string) {
	fa.formatter.PrintListArchive(path, creationTime)
}

func (fa *FormatterAdapter) PrintConfigValue(name, value, source string) {
	fa.formatter.PrintConfigValue(name, value, source)
}

func (fa *FormatterAdapter) PrintError(message string) {
	fa.formatter.PrintError(message)
}

func (fa *FormatterAdapter) PrintDryRunArchive(path string) {
	fa.formatter.PrintDryRunArchive(path)
}

func (fa *FormatterAdapter) PrintCreatedBackup(path string) {
	fa.formatter.PrintCreatedBackup(path)
}

func (fa *FormatterAdapter) PrintIdenticalBackup(path string) {
	fa.formatter.PrintIdenticalBackup(path)
}

func (fa *FormatterAdapter) PrintListBackup(path, creationTime string) {
	fa.formatter.PrintListBackup(path, creationTime)
}

func (fa *FormatterAdapter) PrintDryRunBackup(path string) {
	fa.formatter.PrintDryRunBackup(path)
}

// Template operations
func (fa *FormatterAdapter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
	return fa.formatter.FormatWithTemplate(input, pattern, tmplStr)
}

func (fa *FormatterAdapter) FormatWithPlaceholders(format string, data map[string]string) string {
	return fa.formatter.FormatWithPlaceholders(format, data)
}

func (fa *FormatterAdapter) TemplateCreatedArchive(data map[string]string) string {
	return fa.formatter.TemplateCreatedArchive(data)
}

func (fa *FormatterAdapter) TemplateIdenticalArchive(data map[string]string) string {
	return fa.formatter.TemplateIdenticalArchive(data)
}

func (fa *FormatterAdapter) TemplateListArchive(data map[string]string) string {
	return fa.formatter.TemplateListArchive(data)
}

func (fa *FormatterAdapter) TemplateConfigValue(data map[string]string) string {
	return fa.formatter.TemplateConfigValue(data)
}

func (fa *FormatterAdapter) TemplateDryRunArchive(data map[string]string) string {
	return fa.formatter.TemplateDryRunArchive(data)
}

func (fa *FormatterAdapter) TemplateError(data map[string]string) string {
	return fa.formatter.TemplateError(data)
}

// Pattern extraction operations
func (fa *FormatterAdapter) ExtractArchiveFilenameData(filename string) map[string]string {
	return fa.formatter.ExtractArchiveFilenameData(filename)
}

func (fa *FormatterAdapter) ExtractBackupFilenameData(filename string) map[string]string {
	return fa.formatter.ExtractBackupFilenameData(filename)
}

func (fa *FormatterAdapter) ExtractPatternData(pattern, text string) map[string]string {
	return fa.formatter.ExtractPatternData(pattern, text)
}

// Error formatting operations
func (fa *FormatterAdapter) FormatDiskFullError(err error) string {
	return fa.formatter.FormatDiskFullError(err)
}

func (fa *FormatterAdapter) FormatPermissionError(err error) string {
	return fa.formatter.FormatPermissionError(err)
}

func (fa *FormatterAdapter) FormatDirectoryNotFound(err error) string {
	return fa.formatter.FormatDirectoryNotFound(err)
}

func (fa *FormatterAdapter) FormatFileNotFound(err error) string {
	return fa.formatter.FormatFileNotFound(err)
}

func (fa *FormatterAdapter) FormatInvalidDirectory(err error) string {
	return fa.formatter.FormatInvalidDirectory(err)
}

func (fa *FormatterAdapter) FormatInvalidFile(err error) string {
	return fa.formatter.FormatInvalidFile(err)
}

// Template error formatting operations
func (fa *FormatterAdapter) TemplateDiskFullError(err error) string {
	return fa.formatter.TemplateDiskFullError(err)
}

func (fa *FormatterAdapter) TemplatePermissionError(err error) string {
	return fa.formatter.TemplatePermissionError(err)
}

func (fa *FormatterAdapter) TemplateDirectoryNotFound(err error) string {
	return fa.formatter.TemplateDirectoryNotFound(err)
}

func (fa *FormatterAdapter) TemplateFileNotFound(err error) string {
	return fa.formatter.TemplateFileNotFound(err)
}

// ⭐ EXTRACT-003: Extended formatting methods - 📝 Additional methods from original OutputFormatter

// Extended formatting methods
func (fa *FormatterAdapter) FormatNoArchivesFound(archiveDir string) string {
	// Use config directly since these aren't in the extracted formatter yet
	return fmt.Sprintf(fa.config.FormatNoArchivesFound, archiveDir)
}

func (fa *FormatterAdapter) FormatVerificationFailed(archiveName string, err error) string {
	return fmt.Sprintf(fa.config.FormatVerificationFailed, archiveName, err.Error())
}

func (fa *FormatterAdapter) FormatVerificationSuccess(archiveName string) string {
	return fmt.Sprintf(fa.config.FormatVerificationSuccess, archiveName)
}

func (fa *FormatterAdapter) FormatVerificationWarning(archiveName string, err error) string {
	return fmt.Sprintf(fa.config.FormatVerificationWarning, archiveName, err.Error())
}

func (fa *FormatterAdapter) FormatConfigurationUpdated(key string, value interface{}) string {
	return fmt.Sprintf(fa.config.FormatConfigurationUpdated, key, value)
}

func (fa *FormatterAdapter) FormatConfigFilePath(path string) string {
	return fmt.Sprintf(fa.config.FormatConfigFilePath, path)
}

func (fa *FormatterAdapter) FormatDryRunFilesHeader() string {
	return fa.config.FormatDryRunFilesHeader
}

func (fa *FormatterAdapter) FormatDryRunFileEntry(file string) string {
	return fmt.Sprintf(fa.config.FormatDryRunFileEntry, file)
}

func (fa *FormatterAdapter) FormatNoFilesModified() string {
	return fa.config.FormatNoFilesModified
}

func (fa *FormatterAdapter) FormatIncrementalCreated(path string) string {
	return fmt.Sprintf(fa.config.FormatIncrementalCreated, path)
}

func (fa *FormatterAdapter) FormatNoBackupsFound(filename, backupDir string) string {
	return fmt.Sprintf(fa.config.FormatNoBackupsFound, filename, backupDir)
}

func (fa *FormatterAdapter) FormatBackupWouldCreate(path string) string {
	return fmt.Sprintf(fa.config.FormatBackupWouldCreate, path)
}

func (fa *FormatterAdapter) FormatBackupIdentical(path string) string {
	return fmt.Sprintf(fa.config.FormatBackupIdentical, path)
}

func (fa *FormatterAdapter) FormatBackupCreated(path string) string {
	return fmt.Sprintf(fa.config.FormatBackupCreated, path)
}

// Extended print methods
func (fa *FormatterAdapter) PrintNoArchivesFound(archiveDir string) {
	message := fa.FormatNoArchivesFound(archiveDir)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintVerificationFailed(archiveName string, err error) {
	message := fa.FormatVerificationFailed(archiveName, err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

func (fa *FormatterAdapter) PrintVerificationSuccess(archiveName string) {
	message := fa.FormatVerificationSuccess(archiveName)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintVerificationWarning(archiveName string, err error) {
	message := fa.FormatVerificationWarning(archiveName, err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "warning")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

func (fa *FormatterAdapter) PrintConfigurationUpdated(key string, value interface{}) {
	message := fa.FormatConfigurationUpdated(key, value)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintConfigFilePath(path string) {
	message := fa.FormatConfigFilePath(path)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintDryRunFilesHeader() {
	message := fa.FormatDryRunFilesHeader()
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintDryRunFileEntry(file string) {
	message := fa.FormatDryRunFileEntry(file)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintNoFilesModified() {
	message := fa.FormatNoFilesModified()
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintIncrementalCreated(path string) {
	message := fa.FormatIncrementalCreated(path)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintNoBackupsFound(filename, backupDir string) {
	message := fa.FormatNoBackupsFound(filename, backupDir)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintBackupWouldCreate(path string) {
	message := fa.FormatBackupWouldCreate(path)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintBackupIdentical(path string) {
	message := fa.FormatBackupIdentical(path)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

func (fa *FormatterAdapter) PrintBackupCreated(path string) {
	message := fa.FormatBackupCreated(path)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// ⭐ EXTRACT-003: Backward compatibility adapter - 🔧 Constructor replacement
// NewOutputFormatter creates a FormatterAdapter instead of the original OutputFormatter
func NewOutputFormatter(cfg *Config) *FormatterAdapter {
	return NewFormatterAdapter(cfg)
}

// NewOutputFormatterWithCollector creates a FormatterAdapter with collector
func NewOutputFormatterWithCollector(cfg *Config, collector *formatter.OutputCollector) *FormatterAdapter {
	return NewFormatterAdapterWithCollector(cfg, collector)
}

// ⭐ EXTRACT-003: FormatterAdapter - 📝 Additional print methods for compatibility
// PrintVerificationErrorDetail prints verification error details
func (fa *FormatterAdapter) PrintVerificationErrorDetail(errMsg string) {
	message := fmt.Sprintf("  - %s\n", errMsg)
	if fa.formatter.IsDelayedMode() {
		fa.formatter.GetCollector().AddStdout(message, "error")
	} else {
		fmt.Print(message)
	}
}

// PrintArchiveListWithStatus prints archive list with status
func (fa *FormatterAdapter) PrintArchiveListWithStatus(output, status string) {
	message := fmt.Sprintf("%s%s\n", output, status)
	if fa.formatter.IsDelayedMode() {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// ⭐ EXTRACT-003: FormatterAdapter - 📝 Extended formatting methods with extraction
// FormatListArchiveWithExtraction formats archive listing with data extraction
func (fa *FormatterAdapter) FormatListArchiveWithExtraction(archivePath, creationTime string) string {
	// Extract data from archive filename and format with template
	data := fa.formatter.ExtractArchiveFilenameData(archivePath)
	if data == nil {
		data = make(map[string]string)
	}
	data["path"] = archivePath
	data["creation_time"] = creationTime

	// Try template formatting first, fall back to printf formatting
	if templateStr := fa.config.TemplateListArchive; templateStr != "" {
		return fa.formatter.FormatWithPlaceholders(templateStr, data)
	}

	// Fall back to standard formatting
	return fa.formatter.FormatListArchive(archivePath, creationTime)
}

// FormatListBackupWithExtraction formats backup listing with data extraction
func (fa *FormatterAdapter) FormatListBackupWithExtraction(backupPath, creationTime string) string {
	// Extract data from backup filename and format with template
	data := fa.formatter.ExtractBackupFilenameData(backupPath)
	if data == nil {
		data = make(map[string]string)
	}
	data["path"] = backupPath
	data["creation_time"] = creationTime

	// Try template formatting first, fall back to printf formatting
	if templateStr := fa.config.TemplateListBackup; templateStr != "" {
		return fa.formatter.FormatWithPlaceholders(templateStr, data)
	}

	// Fall back to standard formatting
	return fa.formatter.FormatListBackup(backupPath, creationTime)
}

// ⭐ EXTRACT-003: FormatterAdapter - 📝 Error print methods for compatibility
// PrintDiskFullError prints disk full error message
func (fa *FormatterAdapter) PrintDiskFullError(err error) {
	message := fa.FormatDiskFullError(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintPermissionError prints permission error message
func (fa *FormatterAdapter) PrintPermissionError(err error) {
	message := fa.FormatPermissionError(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintDirectoryNotFound prints directory not found error message
func (fa *FormatterAdapter) PrintDirectoryNotFound(err error) {
	message := fa.FormatDirectoryNotFound(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFileNotFound prints file not found error message
func (fa *FormatterAdapter) PrintFileNotFound(err error) {
	message := fa.FormatFileNotFound(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintInvalidDirectory prints invalid directory error message
func (fa *FormatterAdapter) PrintInvalidDirectory(err error) {
	message := fa.FormatInvalidDirectory(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintInvalidFile prints invalid file error message
func (fa *FormatterAdapter) PrintInvalidFile(err error) {
	message := fa.FormatInvalidFile(err)
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedWriteTemp prints failed write temp error message
func (fa *FormatterAdapter) PrintFailedWriteTemp(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedFinalizeFile prints failed finalize file error message
func (fa *FormatterAdapter) PrintFailedFinalizeFile(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedCreateDirDisk prints failed create directory disk error message
func (fa *FormatterAdapter) PrintFailedCreateDirDisk(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedCreateDir prints failed create directory error message
func (fa *FormatterAdapter) PrintFailedCreateDir(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedAccessDir prints failed access directory error message
func (fa *FormatterAdapter) PrintFailedAccessDir(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// PrintFailedAccessFile prints failed access file error message
func (fa *FormatterAdapter) PrintFailedAccessFile(err error) {
	message := fa.FormatError(err.Error())
	if fa.formatter.GetCollector() != nil {
		fa.formatter.GetCollector().AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// ⭐ OUT-002: Enhanced output with file statistics - Enhanced formatting methods
// Enhanced methods using file statistics for detailed output

// FormatCreatedArchiveWithStats delegates to the enhanced formatter method
func (fa *FormatterAdapter) FormatCreatedArchiveWithStats(path string) string {
	return fa.formatter.FormatCreatedArchiveWithStats(path)
}

// FormatIncrementalCreatedWithStats delegates to the enhanced formatter method
func (fa *FormatterAdapter) FormatIncrementalCreatedWithStats(path string) string {
	return fa.formatter.FormatIncrementalCreatedWithStats(path)
}

// TemplateCreatedArchiveWithStats delegates to the enhanced template method
func (fa *FormatterAdapter) TemplateCreatedArchiveWithStats(path string) string {
	return fa.formatter.TemplateCreatedArchiveWithStats(path)
}

// TemplateIncrementalCreatedWithStats delegates to the enhanced template method
func (fa *FormatterAdapter) TemplateIncrementalCreatedWithStats(path string) string {
	return fa.formatter.TemplateIncrementalCreatedWithStats(path)
}

// PrintCreatedArchiveWithStats delegates to the enhanced print method
func (fa *FormatterAdapter) PrintCreatedArchiveWithStats(path string) {
	fa.formatter.PrintCreatedArchiveWithStats(path)
}

// PrintIncrementalCreatedWithStats delegates to the enhanced print method
func (fa *FormatterAdapter) PrintIncrementalCreatedWithStats(path string) {
	fa.formatter.PrintIncrementalCreatedWithStats(path)
}
