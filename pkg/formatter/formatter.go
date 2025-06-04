// Main output formatter implementation combining all formatter components.
// Provides printf-style formatting operations with optional delayed output
// collection, pattern extraction, and template support.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

import (
	"fmt"
	"os"
)

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Main formatter combining all components
// DefaultOutputFormatter provides comprehensive output formatting functionality
type DefaultOutputFormatter struct {
	configProvider    ConfigProvider
	templateFormatter TemplateFormatter
	patternExtractor  PatternExtractor
	collector         *OutputCollector
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Constructor
// NewDefaultOutputFormatter creates a new DefaultOutputFormatter
func NewDefaultOutputFormatter(configProvider ConfigProvider) *DefaultOutputFormatter {
	return &DefaultOutputFormatter{
		configProvider:    configProvider,
		templateFormatter: NewDefaultTemplateFormatter(configProvider),
		patternExtractor:  NewDefaultPatternExtractor(configProvider),
		collector:         nil,
	}
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Constructor with collector
// NewDefaultOutputFormatterWithCollector creates a formatter with delayed output support
func NewDefaultOutputFormatterWithCollector(configProvider ConfigProvider, collector *OutputCollector) *DefaultOutputFormatter {
	return &DefaultOutputFormatter{
		configProvider:    configProvider,
		templateFormatter: NewDefaultTemplateFormatter(configProvider),
		patternExtractor:  NewDefaultPatternExtractor(configProvider),
		collector:         collector,
	}
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔍 Delayed output mode check
// IsDelayedMode returns true if the formatter is collecting output instead of printing immediately
func (f *DefaultOutputFormatter) IsDelayedMode() bool {
	return f.collector != nil
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔍 Get collector
// GetCollector returns the OutputCollector if in delayed mode, nil otherwise
func (f *DefaultOutputFormatter) GetCollector() *OutputCollector {
	return f.collector
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Set collector
// SetCollector sets the OutputCollector for delayed output mode
func (f *DefaultOutputFormatter) SetCollector(collector *OutputCollector) {
	f.collector = collector
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Printf-style formatting operations

// FormatCreatedArchive formats a created archive message using printf-style formatting
func (f *DefaultOutputFormatter) FormatCreatedArchive(path string) string {
	formatStr := f.configProvider.GetFormatString("created_archive")
	if formatStr == "" {
		formatStr = "Created archive: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// FormatIdenticalArchive formats an identical archive message using printf-style formatting
func (f *DefaultOutputFormatter) FormatIdenticalArchive(path string) string {
	formatStr := f.configProvider.GetFormatString("identical_archive")
	if formatStr == "" {
		formatStr = "Identical archive exists: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// FormatListArchive formats a list archive message using printf-style formatting
func (f *DefaultOutputFormatter) FormatListArchive(path, creationTime string) string {
	formatStr := f.configProvider.GetFormatString("list_archive")
	if formatStr == "" {
		formatStr = "%s (created: %s)\n"
	}
	return fmt.Sprintf(formatStr, path, creationTime)
}

// FormatConfigValue formats a configuration value message using printf-style formatting
func (f *DefaultOutputFormatter) FormatConfigValue(name, value, source string) string {
	formatStr := f.configProvider.GetFormatString("config_value")
	if formatStr == "" {
		formatStr = "%s=%s (source: %s)\n"
	}
	return fmt.Sprintf(formatStr, name, value, source)
}

// FormatError formats an error message using printf-style formatting
func (f *DefaultOutputFormatter) FormatError(message string) string {
	formatStr := f.configProvider.GetErrorFormat("generic")
	if formatStr == "" {
		formatStr = "Error: %s\n"
	}
	return fmt.Sprintf(formatStr, message)
}

// FormatDryRunArchive formats a dry-run archive message using printf-style formatting
func (f *DefaultOutputFormatter) FormatDryRunArchive(path string) string {
	formatStr := f.configProvider.GetFormatString("dry_run_archive")
	if formatStr == "" {
		formatStr = "Would create archive: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// FormatCreatedBackup formats a created backup message using printf-style formatting
func (f *DefaultOutputFormatter) FormatCreatedBackup(path string) string {
	formatStr := f.configProvider.GetFormatString("created_backup")
	if formatStr == "" {
		formatStr = "Created backup: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// FormatIdenticalBackup formats an identical backup message using printf-style formatting
func (f *DefaultOutputFormatter) FormatIdenticalBackup(path string) string {
	formatStr := f.configProvider.GetFormatString("identical_backup")
	if formatStr == "" {
		formatStr = "Identical backup exists: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// FormatListBackup formats a list backup message using printf-style formatting
func (f *DefaultOutputFormatter) FormatListBackup(path, creationTime string) string {
	formatStr := f.configProvider.GetFormatString("list_backup")
	if formatStr == "" {
		formatStr = "%s (created: %s)\n"
	}
	return fmt.Sprintf(formatStr, path, creationTime)
}

// FormatDryRunBackup formats a dry-run backup message using printf-style formatting
func (f *DefaultOutputFormatter) FormatDryRunBackup(path string) string {
	formatStr := f.configProvider.GetFormatString("dry_run_backup")
	if formatStr == "" {
		formatStr = "Would create backup: %s\n"
	}
	return fmt.Sprintf(formatStr, path)
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 📝 Print operations

// PrintCreatedArchive prints a created archive message
func (f *DefaultOutputFormatter) PrintCreatedArchive(path string) {
	message := f.FormatCreatedArchive(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintIdenticalArchive prints an identical archive message
func (f *DefaultOutputFormatter) PrintIdenticalArchive(path string) {
	message := f.FormatIdenticalArchive(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintListArchive prints a list archive message
func (f *DefaultOutputFormatter) PrintListArchive(path, creationTime string) {
	message := f.FormatListArchive(path, creationTime)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintConfigValue prints a configuration value message
func (f *DefaultOutputFormatter) PrintConfigValue(name, value, source string) {
	message := f.FormatConfigValue(name, value, source)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// PrintError prints an error message
func (f *DefaultOutputFormatter) PrintError(message string) {
	formattedMessage := f.FormatError(message)
	if f.IsDelayedMode() {
		f.collector.AddStderr(formattedMessage, "error")
	} else {
		fmt.Fprint(os.Stderr, formattedMessage)
	}
}

// PrintDryRunArchive prints a dry-run archive message
func (f *DefaultOutputFormatter) PrintDryRunArchive(path string) {
	message := f.FormatDryRunArchive(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintCreatedBackup prints a created backup message
func (f *DefaultOutputFormatter) PrintCreatedBackup(path string) {
	message := f.FormatCreatedBackup(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintIdenticalBackup prints an identical backup message
func (f *DefaultOutputFormatter) PrintIdenticalBackup(path string) {
	message := f.FormatIdenticalBackup(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintListBackup prints a list backup message
func (f *DefaultOutputFormatter) PrintListBackup(path, creationTime string) {
	message := f.FormatListBackup(path, creationTime)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintDryRunBackup prints a dry-run backup message
func (f *DefaultOutputFormatter) PrintDryRunBackup(path string) {
	message := f.FormatDryRunBackup(path)
	if f.IsDelayedMode() {
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Delegate template operations to TemplateFormatter

// FormatWithTemplate delegates to the template formatter
func (f *DefaultOutputFormatter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
	return f.templateFormatter.FormatWithTemplate(input, pattern, tmplStr)
}

// FormatWithPlaceholders delegates to the template formatter
func (f *DefaultOutputFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
	return f.templateFormatter.FormatWithPlaceholders(format, data)
}

// TemplateCreatedArchive delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateCreatedArchive(data map[string]string) string {
	return f.templateFormatter.TemplateCreatedArchive(data)
}

// TemplateIdenticalArchive delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return f.templateFormatter.TemplateIdenticalArchive(data)
}

// TemplateListArchive delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateListArchive(data map[string]string) string {
	return f.templateFormatter.TemplateListArchive(data)
}

// TemplateConfigValue delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateConfigValue(data map[string]string) string {
	return f.templateFormatter.TemplateConfigValue(data)
}

// TemplateDryRunArchive delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateDryRunArchive(data map[string]string) string {
	return f.templateFormatter.TemplateDryRunArchive(data)
}

// TemplateError delegates to the template formatter
func (f *DefaultOutputFormatter) TemplateError(data map[string]string) string {
	return f.templateFormatter.TemplateError(data)
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Delegate pattern extraction to PatternExtractor

// ExtractArchiveFilenameData delegates to the pattern extractor
func (f *DefaultOutputFormatter) ExtractArchiveFilenameData(filename string) map[string]string {
	return f.patternExtractor.ExtractArchiveFilenameData(filename)
}

// ExtractBackupFilenameData delegates to the pattern extractor
func (f *DefaultOutputFormatter) ExtractBackupFilenameData(filename string) map[string]string {
	return f.patternExtractor.ExtractBackupFilenameData(filename)
}

// ExtractPatternData delegates to the pattern extractor
func (f *DefaultOutputFormatter) ExtractPatternData(pattern, text string) map[string]string {
	return f.patternExtractor.ExtractPatternData(pattern, text)
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Error formatting operations

// FormatDiskFullError formats a disk full error message
func (f *DefaultOutputFormatter) FormatDiskFullError(err error) string {
	formatStr := f.configProvider.GetErrorFormat("disk_full")
	if formatStr == "" {
		formatStr = "Error: Disk full - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// FormatPermissionError formats a permission error message
func (f *DefaultOutputFormatter) FormatPermissionError(err error) string {
	formatStr := f.configProvider.GetErrorFormat("permission")
	if formatStr == "" {
		formatStr = "Error: Permission denied - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// FormatDirectoryNotFound formats a directory not found error message
func (f *DefaultOutputFormatter) FormatDirectoryNotFound(err error) string {
	formatStr := f.configProvider.GetErrorFormat("directory_not_found")
	if formatStr == "" {
		formatStr = "Error: Directory not found - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// FormatFileNotFound formats a file not found error message
func (f *DefaultOutputFormatter) FormatFileNotFound(err error) string {
	formatStr := f.configProvider.GetErrorFormat("file_not_found")
	if formatStr == "" {
		formatStr = "Error: File not found - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// FormatInvalidDirectory formats an invalid directory error message
func (f *DefaultOutputFormatter) FormatInvalidDirectory(err error) string {
	formatStr := f.configProvider.GetErrorFormat("invalid_directory")
	if formatStr == "" {
		formatStr = "Error: Invalid directory - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// FormatInvalidFile formats an invalid file error message
func (f *DefaultOutputFormatter) FormatInvalidFile(err error) string {
	formatStr := f.configProvider.GetErrorFormat("invalid_file")
	if formatStr == "" {
		formatStr = "Error: Invalid file - %s\n"
	}
	return fmt.Sprintf(formatStr, err.Error())
}

// ⭐ EXTRACT-003: OutputFormatter implementation - 🔧 Template error formatting

// TemplateDiskFullError formats a disk full error using template
func (f *DefaultOutputFormatter) TemplateDiskFullError(err error) string {
	data := map[string]string{"error": err.Error(), "type": "disk_full"}
	templateStr := f.configProvider.GetTemplateString("error_disk_full")
	if templateStr == "" {
		templateStr = "Error: Disk full - %{error}"
	}
	return f.FormatWithPlaceholders(templateStr, data)
}

// TemplatePermissionError formats a permission error using template
func (f *DefaultOutputFormatter) TemplatePermissionError(err error) string {
	data := map[string]string{"error": err.Error(), "type": "permission"}
	templateStr := f.configProvider.GetTemplateString("error_permission")
	if templateStr == "" {
		templateStr = "Error: Permission denied - %{error}"
	}
	return f.FormatWithPlaceholders(templateStr, data)
}

// TemplateDirectoryNotFound formats a directory not found error using template
func (f *DefaultOutputFormatter) TemplateDirectoryNotFound(err error) string {
	data := map[string]string{"error": err.Error(), "type": "directory_not_found"}
	templateStr := f.configProvider.GetTemplateString("error_directory_not_found")
	if templateStr == "" {
		templateStr = "Error: Directory not found - %{error}"
	}
	return f.FormatWithPlaceholders(templateStr, data)
}

// TemplateFileNotFound formats a file not found error using template
func (f *DefaultOutputFormatter) TemplateFileNotFound(err error) string {
	data := map[string]string{"error": err.Error(), "type": "file_not_found"}
	templateStr := f.configProvider.GetTemplateString("error_file_not_found")
	if templateStr == "" {
		templateStr = "Error: File not found - %{error}"
	}
	return f.FormatWithPlaceholders(templateStr, data)
}
