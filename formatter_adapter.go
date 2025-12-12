// Backward compatibility adapter for the extracted formatter package.
// This adapter bridges the original OutputFormatter with the new pkg/formatter
// components while maintaining all existing functionality and API compatibility.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// FORMATTER-COMPATIBILITY-001: Formatter compatibility specification - Backward compatibility formatter adapter [ACTION:format-processing]
// Source: formatter_adapter.go - FORMATTER-COMPATIBILITY-001
// Impact: Core functionality requirement for formatter compatibility

// SERVICE-FORMATTER-001: Formatter service architecture decision - Formatter service implementation [ACTION:core-functionality]
// Source: formatter_adapter.go - SERVICE-FORMATTER-001
// Impact: Formatter service implementation decision
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"strings"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: GetFormatString(created_archive) = %q\n", result)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
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

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]

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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]

// Extended formatting methods
func (fa *FormatterAdapter) FormatNoArchivesFound(archiveDir string) string {
	// Use config directly since these aren't in the extracted formatter yet
	return fmt.Sprintf(fa.config.FormatNoArchivesFound, archiveDir)
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]

// PrintArchiveListWithStatus prints archive list with status
// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] CRITICAL: output may contain #{...} patterns if placeholder replacement failed
// Use os.Stdout.WriteString instead of fmt.Print to avoid fmt misinterpreting #{...} as format verbs
func (fa *FormatterAdapter) PrintArchiveListWithStatus(output, status string) {
	// CRITICAL: Use os.Stdout.WriteString, NOT fmt.Print, to avoid fmt misinterpreting #{...} patterns
	// fmt.Print internally uses fmt.Sprintf("%v", ...) which can misinterpret #{...} as format verbs
	message := output + status + "\n"
	if fa.formatter.IsDelayedMode() {
		fa.formatter.GetCollector().AddStdout(message, "info")
	} else {
		os.Stdout.WriteString(message)
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// FormatListArchiveWithExtraction formats archive listing with data extraction
// FormatListArchiveWithExtraction formats archive listing with a simplified, testable implementation.
// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] This implementation:
// 1. Always gathers file statistics (needed for template placeholders)
// 2. Uses FormatListArchive if it contains template placeholders, otherwise TemplateListArchive
// 3. Processes all placeholders in a single, clear code path
// 4. Returns formatted string ready for output
func (fa *FormatterAdapter) FormatListArchiveWithExtraction(archivePath, creationTime string) string {
	// Use the simplified implementation
	return formatListArchiveSimple(fa.config, fa, archivePath, creationTime)
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

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Gather file statistics to populate size_human and other stat fields
	// Per requirements: list output MUST support template-style placeholders for file attributes
	statInfo, err := formatter.GatherFileStatInfo(backupPath)
	if err == nil {
		// Add file statistics to data map
		data["size"] = fmt.Sprintf("%d", statInfo.Size)
		data["size_human"] = statInfo.SizeHuman
		data["mtime"] = statInfo.MTime.Format("2006-01-02 15:04:05")
		data["mtime_unix"] = fmt.Sprintf("%d", statInfo.MTimeUnix)
		data["mode"] = statInfo.Mode.String()
		data["type"] = statInfo.Type
		data["name"] = statInfo.Name
	} else {
		// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Provide default values for file attributes when stats can't be gathered
		// This ensures placeholders like #{size_human} are always available in the data map
		if _, exists := data["size"]; !exists {
			data["size"] = "0"
		}
		if _, exists := data["size_human"]; !exists {
			data["size_human"] = "unknown"
		}
		if _, exists := data["mtime"]; !exists {
			data["mtime"] = creationTime // Use creation_time as fallback
		}
		if _, exists := data["mtime_unix"]; !exists {
			data["mtime_unix"] = "0"
		}
		if _, exists := data["mode"]; !exists {
			data["mode"] = "unknown"
		}
		if _, exists := data["type"]; !exists {
			data["type"] = "unknown"
		}
		if _, exists := data["name"]; !exists {
			// Extract filename from path
			parts := strings.Split(backupPath, "/")
			if len(parts) > 0 {
				data["name"] = parts[len(parts)-1]
			} else {
				data["name"] = backupPath
			}
		}
	}

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Per requirements: list output MUST use template-style format strings
	// Always prefer TemplateListBackup for list output to support file attributes
	templateStr := fa.config.TemplateListBackup
	if templateStr == "" {
		// Use default template format that includes size_human as per requirements
		templateStr = "#{path} (size: #{size_human})\n"
	}
	return fa.formatter.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
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
