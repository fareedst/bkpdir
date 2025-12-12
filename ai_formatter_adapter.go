// AI-First Formatter Adapter
// Bridges the main application to use the new AI-first formatter from pkg/formatter
// while maintaining backward compatibility with existing interfaces.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:core-functionality]
package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"strings"
)

// [CRITICAL] FMT-001: AI-first formatter config adapter - [ACTION:configure-modify]
type AIFormatterConfigAdapter struct {
	config *Config
}

// [CRITICAL] FMT-001: AI-first formatter config adapter constructor - [ACTION:configure-modify]
func NewAIFormatterConfigAdapter(config *Config) *AIFormatterConfigAdapter {
	return &AIFormatterConfigAdapter{config: config}
}

// [CRITICAL] FMT-001: Format string access - [ACTION:configure-modify]
func (a *AIFormatterConfigAdapter) GetFormatString(formatType formatter.FormatType) (string, error) {
	switch formatType {
	case formatter.FormatTypeCreated:
		return a.config.FormatCreatedArchive, nil
	case formatter.FormatTypeIdentical:
		return a.config.FormatIdenticalArchive, nil
	case formatter.FormatTypeList:
		return a.config.FormatListArchive, nil
	case formatter.FormatTypeDryRun:
		return a.config.FormatDryRunArchive, nil
	case formatter.FormatTypeError:
		return a.config.FormatError, nil
	default:
		return "", fmt.Errorf("unsupported format type: %s", formatType)
	}
}

// [CRITICAL] FMT-001: Template string access - [ACTION:configure-modify]
func (a *AIFormatterConfigAdapter) GetTemplateString(templateType formatter.FormatType) (string, error) {
	switch templateType {
	case formatter.FormatTypeCreated:
		return a.config.TemplateCreatedArchive, nil
	case formatter.FormatTypeIdentical:
		return a.config.TemplateIdenticalArchive, nil
	case formatter.FormatTypeList:
		return a.config.TemplateListArchive, nil
	case formatter.FormatTypeDryRun:
		return a.config.TemplateDryRunArchive, nil
	case formatter.FormatTypeError:
		return a.config.TemplateError, nil
	case formatter.FormatTypeConfig:
		return a.config.TemplateConfigValue, nil
	default:
		return "", fmt.Errorf("unsupported template type: %s", templateType)
	}
}

// [CRITICAL] FMT-001: Error format access - [ACTION:configure-modify]
func (a *AIFormatterConfigAdapter) GetErrorFormat(errorType formatter.ErrorType) (string, error) {
	switch errorType {
	case formatter.ErrorTypeDiskFull:
		return a.config.FormatDiskFullError, nil
	case formatter.ErrorTypePermission:
		return a.config.FormatPermissionError, nil
	case formatter.ErrorTypeDirectoryNotFound:
		return a.config.FormatDirectoryNotFound, nil
	case formatter.ErrorTypeFileNotFound:
		return a.config.FormatFileNotFound, nil
	case formatter.ErrorTypeInvalidDirectory:
		return a.config.FormatInvalidDirectory, nil
	case formatter.ErrorTypeInvalidFile:
		return a.config.FormatInvalidFile, nil
	case formatter.ErrorTypeGeneric:
		return a.config.FormatError, nil
	default:
		return "", fmt.Errorf("unsupported error type: %s", errorType)
	}
}

// [CRITICAL] FMT-001: Pattern access - [ACTION:discovery]
func (a *AIFormatterConfigAdapter) GetPattern(patternType formatter.PatternType) (string, error) {
	switch patternType {
	case formatter.PatternTypeArchiveFilename:
		return a.config.PatternArchiveFilename, nil
	case formatter.PatternTypeBackupFilename:
		return a.config.PatternBackupFilename, nil
	case formatter.PatternTypeConfigLine:
		return a.config.PatternConfigLine, nil
	case formatter.PatternTypeTimestamp:
		return a.config.PatternTimestamp, nil
	default:
		return "", fmt.Errorf("unsupported pattern type: %s", patternType)
	}
}

// [CRITICAL] FMT-001: Configuration validation - [ACTION:configure-modify]
func (a *AIFormatterConfigAdapter) Validate() error {
	// Basic validation - ensure required fields are present
	if a.config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	return nil
}

// [CRITICAL] FMT-001: Configuration error access - [ACTION:configure-modify]
func (a *AIFormatterConfigAdapter) GetValidationErrors() []formatter.ConfigError {
	return []formatter.ConfigError{} // No validation errors for now
}

// [CRITICAL] FMT-001: AI-first formatter adapter - [ACTION:core-functionality]
type AIFormatterAdapter struct {
	aiFormatter formatter.AIFirstFormatter
	config      *Config
}

// [CRITICAL] FMT-001: AI-first formatter adapter constructor - [ACTION:core-functionality]
func NewAIFormatterAdapter(config *Config) *AIFormatterAdapter {
	configAdapter := NewAIFormatterConfigAdapter(config)
	aiFormatter := formatter.NewAIFirstFormatter(configAdapter)

	return &AIFormatterAdapter{
		aiFormatter: aiFormatter,
		config:      config,
	}
}

// [CRITICAL] FMT-001: AI-first formatter adapter with collector - [ACTION:core-functionality]
func NewAIFormatterAdapterWithCollector(config *Config, collector *formatter.OutputCollector) *AIFormatterAdapter {
	configAdapter := NewAIFormatterConfigAdapter(config)
	aiFormatter := formatter.NewAIFirstFormatterWithCollector(configAdapter, collector)

	return &AIFormatterAdapter{
		aiFormatter: aiFormatter,
		config:      config,
	}
}

// [CRITICAL] FMT-001: Delayed mode management - [ACTION:format-processing]
func (fa *AIFormatterAdapter) IsDelayedMode() bool {
	return fa.aiFormatter.IsDelayedMode()
}

func (fa *AIFormatterAdapter) GetCollector() *formatter.OutputCollector {
	// The AI-first formatter doesn't expose collector directly
	// We'll need to handle this differently
	return nil
}

func (fa *AIFormatterAdapter) SetCollector(collector *formatter.OutputCollector) {
	// Create a new formatter with the collector
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	fa.aiFormatter = formatter.NewAIFirstFormatterWithCollector(configAdapter, collector)
}

// [CRITICAL] FMT-001: Format operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatCreatedArchive(path string) string {
	result, err := fa.aiFormatter.FormatArchive(path, formatter.FormatTypeCreated)
	if err != nil {
		return fmt.Sprintf("Created archive: %s\n", path)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatIdenticalArchive(path string) string {
	result, err := fa.aiFormatter.FormatArchive(path, formatter.FormatTypeIdentical)
	if err != nil {
		return fmt.Sprintf("Identical archive exists: %s\n", path)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatListArchive(path, creationTime string) string {
	ctx := formatter.FormatContext{
		FormatType: formatter.FormatTypeList,
		Data: map[string]interface{}{
			"path":         path,
			"creationTime": creationTime,
		},
	}
	result, err := fa.aiFormatter.FormatWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("%s (created: %s)\n", path, creationTime)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatConfigValue(name, value, source string) string {
	result, err := fa.aiFormatter.FormatConfig(name, value, source)
	if err != nil {
		return fmt.Sprintf("%s = %s (from %s)\n", name, value, source)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatError(message string) string {
	result, err := fa.aiFormatter.FormatError(fmt.Errorf(message), formatter.ErrorTypeGeneric)
	if err != nil {
		return fmt.Sprintf("Error: %s\n", message)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatDryRunArchive(path string) string {
	result, err := fa.aiFormatter.FormatArchive(path, formatter.FormatTypeDryRun)
	if err != nil {
		return fmt.Sprintf("Would create archive: %s\n", path)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatCreatedBackup(path string) string {
	result, err := fa.aiFormatter.FormatBackup(path, formatter.FormatTypeCreated)
	if err != nil {
		return fmt.Sprintf("Created backup: %s\n", path)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatIdenticalBackup(path string) string {
	result, err := fa.aiFormatter.FormatBackup(path, formatter.FormatTypeIdentical)
	if err != nil {
		return fmt.Sprintf("Identical backup exists: %s\n", path)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatListBackup(path, creationTime string) string {
	ctx := formatter.FormatContext{
		FormatType: formatter.FormatTypeList,
		Data: map[string]interface{}{
			"path":         path,
			"creationTime": creationTime,
		},
	}
	result, err := fa.aiFormatter.FormatWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("%s (created: %s)\n", path, creationTime)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatDryRunBackup(path string) string {
	result, err := fa.aiFormatter.FormatBackup(path, formatter.FormatTypeDryRun)
	if err != nil {
		return fmt.Sprintf("Would create backup: %s\n", path)
	}
	return result
}

// [CRITICAL] FMT-001: Print operations - [ACTION:format-processing]
func (fa *AIFormatterAdapter) PrintCreatedArchive(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatCreatedArchive(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintIdenticalArchive(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatIdenticalArchive(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintListArchive(path, creationTime string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatListArchive(path, creationTime),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintConfigValue(name, value, source string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatConfigValue(name, value, source),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeConfig,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintError(message string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(message),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintDryRunArchive(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDryRunArchive(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintCreatedBackup(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatCreatedBackup(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintIdenticalBackup(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatIdenticalBackup(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintListBackup(path, creationTime string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatListBackup(path, creationTime),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintDryRunBackup(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDryRunBackup(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

// [CRITICAL] FMT-001: Template operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
	return fa.aiFormatter.FormatWithTemplate(tmplStr, map[string]interface{}{
		"input":   input,
		"pattern": pattern,
	})
}

func (fa *AIFormatterAdapter) FormatWithPlaceholders(format string, data map[string]string) string {
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatWithPlaceholders called with format=%q, data=%+v\n", format, data)
	}
	result, err := fa.aiFormatter.FormatWithPlaceholders(format, data)
	if err != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatWithPlaceholders error, using fallback: %v\n", err)
		}
		// Fallback to simple string replacement
		result := format
		for key, value := range data {
			placeholder := "#{" + key + "}"
			if strings.Contains(result, placeholder) {
				if debug {
					fmt.Fprintf(os.Stderr, "DEBUG: Replacing %q with %q\n", placeholder, value)
				}
				result = strings.ReplaceAll(result, placeholder, value)
			}
		}
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatWithPlaceholders fallback result=%q\n", result)
		}
		return result
	}
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatWithPlaceholders result=%q\n", result)
	}
	return result
}

func (fa *AIFormatterAdapter) TemplateCreatedArchive(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeCreated)
	if err != nil || template == "" {
		return fmt.Sprintf("Created archive: %s\n", data["path"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

func (fa *AIFormatterAdapter) TemplateIdenticalArchive(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeIdentical)
	if err != nil || template == "" {
		return fmt.Sprintf("Identical archive exists: %s\n", data["path"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

func (fa *AIFormatterAdapter) TemplateListArchive(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeList)
	if err != nil || template == "" {
		return fmt.Sprintf("%s (created: %s)\n", data["path"], data["creation_time"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

func (fa *AIFormatterAdapter) TemplateConfigValue(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeConfig)
	if err != nil || template == "" {
		return fmt.Sprintf("%s = %s (from %s)\n", data["name"], data["value"], data["source"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

func (fa *AIFormatterAdapter) TemplateDryRunArchive(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeDryRun)
	if err != nil || template == "" {
		return fmt.Sprintf("Would create archive: %s\n", data["path"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

func (fa *AIFormatterAdapter) TemplateError(data map[string]string) string {
	configAdapter := NewAIFormatterConfigAdapter(fa.config)
	template, err := configAdapter.GetTemplateString(formatter.FormatTypeError)
	if err != nil || template == "" {
		return fmt.Sprintf("Error: %s\n", data["message"])
	}
	return fa.FormatWithPlaceholders(template, data)
}

// [CRITICAL] FMT-001: Pattern extraction operations - [ACTION:discovery]
func (fa *AIFormatterAdapter) ExtractArchiveFilenameData(filename string) map[string]string {
	data, err := fa.aiFormatter.ExtractArchiveData(filename)
	if err != nil {
		return make(map[string]string)
	}

	// Convert AIArchiveData to map[string]string
	result := make(map[string]string)
	result["prefix"] = data.Prefix
	result["year"] = data.Year
	result["month"] = data.Month
	result["day"] = data.Day
	result["hour"] = data.Hour
	result["minute"] = data.Minute
	result["branch"] = data.Branch
	result["hash"] = data.Hash
	result["note"] = data.Note

	return result
}

func (fa *AIFormatterAdapter) ExtractBackupFilenameData(filename string) map[string]string {
	data, err := fa.aiFormatter.ExtractBackupData(filename)
	if err != nil {
		return make(map[string]string)
	}

	// Convert AIBackupData to map[string]string
	result := make(map[string]string)
	result["filename"] = data.Filename
	result["year"] = data.Year
	result["month"] = data.Month
	result["day"] = data.Day
	result["hour"] = data.Hour
	result["minute"] = data.Minute
	result["note"] = data.Note

	return result
}

func (fa *AIFormatterAdapter) ExtractPatternData(pattern, text string) map[string]string {
	result, err := fa.aiFormatter.ExtractPattern(pattern, text)
	if err != nil {
		return make(map[string]string)
	}
	return result
}

// [CRITICAL] FMT-001: Error formatting operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatDiskFullError(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypeDiskFull)
	if err2 != nil {
		return fmt.Sprintf("Error: Disk full - %s\n", err.Error())
	}
	return result
}

func (fa *AIFormatterAdapter) FormatPermissionError(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypePermission)
	if err2 != nil {
		return fmt.Sprintf("Error: Permission denied - %s\n", err.Error())
	}
	return result
}

func (fa *AIFormatterAdapter) FormatDirectoryNotFound(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypeDirectoryNotFound)
	if err2 != nil {
		return fmt.Sprintf("Error: Directory not found - %s\n", err.Error())
	}
	return result
}

func (fa *AIFormatterAdapter) FormatFileNotFound(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypeFileNotFound)
	if err2 != nil {
		return fmt.Sprintf("Error: File not found - %s\n", err.Error())
	}
	return result
}

func (fa *AIFormatterAdapter) FormatInvalidDirectory(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypeInvalidDirectory)
	if err2 != nil {
		return fmt.Sprintf("Error: Invalid directory - %s\n", err.Error())
	}
	return result
}

func (fa *AIFormatterAdapter) FormatInvalidFile(err error) string {
	result, err2 := fa.aiFormatter.FormatError(err, formatter.ErrorTypeInvalidFile)
	if err2 != nil {
		return fmt.Sprintf("Error: Invalid file - %s\n", err.Error())
	}
	return result
}

// [CRITICAL] FMT-001: Template error operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) TemplateDiskFullError(err error) string {
	data := map[string]string{"message": err.Error()}
	return fa.TemplateError(data)
}

func (fa *AIFormatterAdapter) TemplatePermissionError(err error) string {
	data := map[string]string{"message": err.Error()}
	return fa.TemplateError(data)
}

func (fa *AIFormatterAdapter) TemplateDirectoryNotFound(err error) string {
	data := map[string]string{"message": err.Error()}
	return fa.TemplateError(data)
}

func (fa *AIFormatterAdapter) TemplateFileNotFound(err error) string {
	data := map[string]string{"message": err.Error()}
	return fa.TemplateError(data)
}

// [CRITICAL] FMT-001: Additional format operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatNoArchivesFound(archiveDir string) string {
	return fmt.Sprintf("No archives found in: %s\n", archiveDir)
}

func (fa *AIFormatterAdapter) FormatConfigurationUpdated(key string, value interface{}) string {
	return fmt.Sprintf("Configuration updated: %s = %v\n", key, value)
}

func (fa *AIFormatterAdapter) FormatConfigFilePath(path string) string {
	return fmt.Sprintf("Configuration file: %s\n", path)
}

func (fa *AIFormatterAdapter) FormatDryRunFilesHeader() string {
	return "Files that would be processed:\n"
}

func (fa *AIFormatterAdapter) FormatDryRunFileEntry(file string) string {
	return fmt.Sprintf("  - %s\n", file)
}

func (fa *AIFormatterAdapter) FormatNoFilesModified() string {
	return "No files were modified.\n"
}

func (fa *AIFormatterAdapter) FormatIncrementalCreated(path string) string {
	return fmt.Sprintf("Created incremental archive: %s\n", path)
}

func (fa *AIFormatterAdapter) FormatNoBackupsFound(filename, backupDir string) string {
	return fmt.Sprintf("No backups found for %s in: %s\n", filename, backupDir)
}

func (fa *AIFormatterAdapter) FormatBackupWouldCreate(path string) string {
	return fmt.Sprintf("Would create backup: %s\n", path)
}

func (fa *AIFormatterAdapter) FormatBackupIdentical(path string) string {
	return fmt.Sprintf("Identical backup exists: %s\n", path)
}

func (fa *AIFormatterAdapter) FormatBackupCreated(path string) string {
	return fmt.Sprintf("Created backup: %s\n", path)
}

// [CRITICAL] FMT-001: Additional print operations - [ACTION:format-processing]
func (fa *AIFormatterAdapter) PrintNoArchivesFound(archiveDir string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatNoArchivesFound(archiveDir),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintConfigurationUpdated(key string, value interface{}) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatConfigurationUpdated(key, value),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeConfig,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintConfigFilePath(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatConfigFilePath(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeConfig,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintDryRunFilesHeader() {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDryRunFilesHeader(),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintDryRunFileEntry(file string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDryRunFileEntry(file),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintNoFilesModified() {
	ctx := formatter.PrintContext{
		Message:     fa.FormatNoFilesModified(),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintIncrementalCreated(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatIncrementalCreated(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintNoBackupsFound(filename, backupDir string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatNoBackupsFound(filename, backupDir),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintBackupWouldCreate(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatBackupWouldCreate(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintBackupIdentical(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatBackupIdentical(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintBackupCreated(path string) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatBackupCreated(path),
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

// [CRITICAL] FMT-001: New output formatter functions - [ACTION:core-functionality]
func NewOutputFormatter(cfg *Config) *AIFormatterAdapter {
	return NewAIFormatterAdapter(cfg)
}

func NewOutputFormatterWithCollector(cfg *Config, collector *formatter.OutputCollector) *AIFormatterAdapter {
	return NewAIFormatterAdapterWithCollector(cfg, collector)
}

// [CRITICAL] FMT-001: Additional error print operations - [ACTION:format-processing]
func (fa *AIFormatterAdapter) PrintDiskFullError(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDiskFullError(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintPermissionError(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatPermissionError(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintDirectoryNotFound(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatDirectoryNotFound(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFileNotFound(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatFileNotFound(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintInvalidDirectory(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatInvalidDirectory(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintInvalidFile(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatInvalidFile(err),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedWriteTemp(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedFinalizeFile(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedCreateDirDisk(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedCreateDir(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedAccessDir(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintFailedAccessFile(err error) {
	ctx := formatter.PrintContext{
		Message:     fa.FormatError(err.Error()),
		Destination: formatter.AIOutputDestinationStderr,
		Type:        formatter.AIMessageTypeError,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

func (fa *AIFormatterAdapter) PrintArchiveListWithStatus(output, status string) {
	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] CRITICAL: Use string concatenation, NOT fmt.Sprintf
	// to avoid fmt misinterpreting #{...} patterns as format verbs
	ctx := formatter.PrintContext{
		Message:     output + status + "\n",
		Destination: formatter.AIOutputDestinationStdout,
		Type:        formatter.AIMessageTypeInfo,
	}
	fa.aiFormatter.PrintWithContext(ctx)
}

// [CRITICAL] FMT-001: Enhanced format operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatListArchiveWithExtraction(archivePath, creationTime string) string {
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction called with path=%q, creationTime=%q\n", archivePath, creationTime)
	}
	// Extract data from archive filename and format with template
	data := fa.ExtractArchiveFilenameData(archivePath)
	if data == nil {
		data = make(map[string]string)
	}
	data["path"] = archivePath
	data["creation_time"] = creationTime

	// Gather file statistics for placeholders like #{size_human}
	statInfo, err := formatter.GatherFileStatInfo(archivePath)
	if err == nil && statInfo != nil {
		data["size"] = fmt.Sprintf("%d", statInfo.Size)
		data["size_human"] = statInfo.SizeHuman
		data["mtime"] = statInfo.MTime.Format("2006-01-02 15:04:05")
		data["mtime_unix"] = fmt.Sprintf("%d", statInfo.MTimeUnix)
		data["mode"] = statInfo.Mode.String()
		data["type"] = statInfo.Type
		if statInfo.Name != "" {
			data["name"] = statInfo.Name
		}
	} else {
		// Provide defaults for missing stats
		data["size"] = "0"
		data["size_human"] = "unknown"
		data["mtime"] = creationTime
		data["mtime_unix"] = "0"
		data["mode"] = "unknown"
		data["type"] = "unknown"
		// Extract filename from path
		parts := strings.Split(archivePath, "/")
		if len(parts) > 0 {
			data["name"] = parts[len(parts)-1]
		} else {
			data["name"] = archivePath
		}
	}

	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction data map: %+v\n", data)
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction FormatListArchive=%q\n", fa.config.FormatListArchive)
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction TemplateListArchive=%q\n", fa.config.TemplateListArchive)
	}

	// Priority: FormatListArchive (if contains #{) > TemplateListArchive > default
	var formatStr string
	if fa.config.FormatListArchive != "" && strings.Contains(fa.config.FormatListArchive, "#{") {
		formatStr = fa.config.FormatListArchive
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction using FormatListArchive=%q\n", formatStr)
		}
	} else if fa.config.TemplateListArchive != "" {
		formatStr = fa.config.TemplateListArchive
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction using TemplateListArchive=%q\n", formatStr)
		}
	} else {
		// Default template format when both are empty
		formatStr = "#{path} (size: #{size_human})\n"
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction using default format=%q\n", formatStr)
		}
	}

	// Use FormatWithPlaceholders to replace all placeholders
	result := fa.FormatWithPlaceholders(formatStr, data)
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: AIFormatterAdapter.FormatListArchiveWithExtraction final result=%q\n", result)
	}
	return result
}

func (fa *AIFormatterAdapter) FormatListBackupWithExtraction(backupPath, creationTime string) string {
	// Extract data from backup filename and format with template
	data := fa.ExtractBackupFilenameData(backupPath)
	if data == nil {
		data = make(map[string]string)
	}
	data["path"] = backupPath
	data["creation_time"] = creationTime

	// Try template formatting first, fall back to printf formatting
	if templateStr := fa.config.TemplateListBackup; templateStr != "" {
		return fa.FormatWithPlaceholders(templateStr, data)
	}

	// Fall back to standard formatting
	return fa.FormatListBackup(backupPath, creationTime)
}

// [CRITICAL] FMT-001: Stats format operations - [ACTION:core-functionality]
func (fa *AIFormatterAdapter) FormatCreatedArchiveWithStats(path string) string {
	return fa.FormatCreatedArchive(path)
}

func (fa *AIFormatterAdapter) FormatIncrementalCreatedWithStats(path string) string {
	return fa.FormatIncrementalCreated(path)
}

func (fa *AIFormatterAdapter) TemplateCreatedArchiveWithStats(path string) string {
	return fa.TemplateCreatedArchive(map[string]string{"path": path})
}

func (fa *AIFormatterAdapter) TemplateIncrementalCreatedWithStats(path string) string {
	return fa.TemplateCreatedArchive(map[string]string{"path": path})
}

func (fa *AIFormatterAdapter) PrintCreatedArchiveWithStats(path string) {
	fa.PrintCreatedArchive(path)
}

func (fa *AIFormatterAdapter) PrintIncrementalCreatedWithStats(path string) {
	fa.PrintIncrementalCreated(path)
}
