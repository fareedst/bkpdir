// This file is part of bkpdir
//
// Package main provides output formatting for BkpDir CLI and tests.
// It handles printf-style and template-based formatting of output messages.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// OutputFormatter provides methods for formatting and printing output for BkpDir operations.
type OutputFormatter struct {
	cfg *Config
}

// NewOutputFormatter creates a new OutputFormatter with the given config.
func NewOutputFormatter(cfg *Config) *OutputFormatter {
	return &OutputFormatter{cfg: cfg}
}

// FormatCreatedArchive formats a message for a created archive.
func (f *OutputFormatter) FormatCreatedArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedArchive, path)
}

// FormatIdenticalArchive formats a message for an identical archive.
func (f *OutputFormatter) FormatIdenticalArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalArchive, path)
}

// FormatListArchive formats a message for listing an archive.
func (f *OutputFormatter) FormatListArchive(path, creationTime string) string {
	return fmt.Sprintf(f.cfg.FormatListArchive, path, creationTime)
}

// FormatConfigValue formats a configuration value for display.
func (f *OutputFormatter) FormatConfigValue(name, value, source string) string {
	return fmt.Sprintf(f.cfg.FormatConfigValue, name, value, source)
}

// FormatDryRunArchive formats a message for a dry-run archive operation.
func (f *OutputFormatter) FormatDryRunArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunArchive, path)
}

// FormatError formats an error message.
func (f *OutputFormatter) FormatError(message string) string {
	return fmt.Sprintf(f.cfg.FormatError, message)
}

// PrintCreatedArchive prints a created archive message to stdout.
func (f *OutputFormatter) PrintCreatedArchive(path string) {
	fmt.Print(f.FormatCreatedArchive(path))
}

// PrintIdenticalArchive prints an identical archive message to stdout.
func (f *OutputFormatter) PrintIdenticalArchive(path string) {
	fmt.Print(f.FormatIdenticalArchive(path))
}

// PrintListArchive prints a list archive message to stdout.
func (f *OutputFormatter) PrintListArchive(path, creationTime string) {
	fmt.Print(f.FormatListArchive(path, creationTime))
}

// PrintConfigValue prints a config value message to stdout.
func (f *OutputFormatter) PrintConfigValue(name, value, source string) {
	fmt.Print(f.FormatConfigValue(name, value, source))
}

// PrintDryRunArchive prints a dry-run archive message to stdout.
func (f *OutputFormatter) PrintDryRunArchive(path string) {
	fmt.Print(f.FormatDryRunArchive(path))
}

// PrintError prints an error message to stderr.
func (f *OutputFormatter) PrintError(message string) {
	fmt.Fprint(os.Stderr, f.FormatError(message))
}

// PrintCreatedBackup prints a created backup message to stdout.
func (f *OutputFormatter) PrintCreatedBackup(path string) {
	fmt.Print(f.FormatCreatedBackup(path))
}

// PrintIdenticalBackup prints an identical backup message to stdout.
func (f *OutputFormatter) PrintIdenticalBackup(path string) {
	fmt.Print(f.FormatIdenticalBackup(path))
}

// PrintListBackup prints a list backup message to stdout.
func (f *OutputFormatter) PrintListBackup(path, creationTime string) {
	fmt.Print(f.FormatListBackup(path, creationTime))
}

// PrintDryRunBackup prints a dry-run backup message to stdout.
func (f *OutputFormatter) PrintDryRunBackup(path string) {
	fmt.Print(f.FormatDryRunBackup(path))
}

// ExtractArchiveFilenameData extracts data from an archive filename using a regex pattern.
func (f *OutputFormatter) ExtractArchiveFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternArchiveFilename, filename)
}

// ExtractBackupFilenameData extracts data from a backup filename using a regex pattern.
func (f *OutputFormatter) ExtractBackupFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternBackupFilename, filename)
}

// ExtractConfigLineData extracts data from a config line using a regex pattern.
func (f *OutputFormatter) ExtractConfigLineData(line string) map[string]string {
	return f.extractPatternData(f.cfg.PatternConfigLine, line)
}

// ExtractTimestampData extracts data from a timestamp string using a regex pattern.
func (f *OutputFormatter) ExtractTimestampData(timestamp string) map[string]string {
	return f.extractPatternData(f.cfg.PatternTimestamp, timestamp)
}

// FormatArchiveWithExtraction formats an archive path using extracted data and a template.
func (f *OutputFormatter) FormatArchiveWithExtraction(archivePath string) string {
	// Extract data from archive filename
	filename := getFilenameFromPath(archivePath)
	data := f.ExtractArchiveFilenameData(filename)
	data["path"] = archivePath

	// Use template formatting if we have extracted data
	if len(data) > 1 { // More than just "path"
		return f.FormatCreatedArchiveTemplate(data)
	}

	// Fall back to printf-style formatting
	return f.FormatCreatedArchive(archivePath)
}

// FormatListArchiveWithExtraction formats a list archive message using extracted data from the filename.
func (f *OutputFormatter) FormatListArchiveWithExtraction(archivePath, creationTime string) string {
	// Extract data from archive filename
	filename := getFilenameFromPath(archivePath)
	data := f.ExtractArchiveFilenameData(filename)
	data["path"] = archivePath
	data["creation_time"] = creationTime

	// Use template formatting if we have extracted data
	if len(data) > 2 { // More than just "path" and "creation_time"
		return f.FormatListArchiveTemplate(data)
	}

	// Fall back to printf-style formatting
	return f.FormatListArchive(archivePath, creationTime)
}

func getFilenameFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// FormatCreatedBackup formats a message for a created backup.
func (f *OutputFormatter) FormatCreatedBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedBackup, path)
}

// FormatIdenticalBackup formats a message for an identical backup.
func (f *OutputFormatter) FormatIdenticalBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalBackup, path)
}

// FormatListBackup formats a message for listing a backup.
func (f *OutputFormatter) FormatListBackup(path, creationTime string) string {
	return fmt.Sprintf(f.cfg.FormatListBackup, path, creationTime)
}

// FormatDryRunBackup formats a message for a dry-run backup operation.
func (f *OutputFormatter) FormatDryRunBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunBackup, path)
}

// FormatCreatedArchiveTemplate formats a created archive message using a template.
func (f *OutputFormatter) FormatCreatedArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateCreatedArchive, data)
}

// FormatIdenticalArchiveTemplate formats an identical archive message using a template.
func (f *OutputFormatter) FormatIdenticalArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIdenticalArchive, data)
}

// FormatListArchiveTemplate formats a list archive message using a template.
func (f *OutputFormatter) FormatListArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateListArchive, data)
}

// FormatConfigValueTemplate formats a config value message using a template.
func (f *OutputFormatter) FormatConfigValueTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateConfigValue, data)
}

// FormatDryRunArchiveTemplate formats a dry-run archive message using a template.
func (f *OutputFormatter) FormatDryRunArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateDryRunArchive, data)
}

// FormatErrorTemplate formats an error message using a template.
func (f *OutputFormatter) FormatErrorTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateError, data)
}

// FormatCreatedBackupTemplate formats a created backup message using a template.
func (f *OutputFormatter) FormatCreatedBackupTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateCreatedBackup, data)
}

// FormatIdenticalBackupTemplate formats an identical backup message using a template.
func (f *OutputFormatter) FormatIdenticalBackupTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIdenticalBackup, data)
}

// FormatListBackupTemplate formats a list backup message using a template.
func (f *OutputFormatter) FormatListBackupTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateListBackup, data)
}

// FormatDryRunBackupTemplate formats a dry-run backup message using a template.
func (f *OutputFormatter) FormatDryRunBackupTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateDryRunBackup, data)
}

// Template formatting helper
func (f *OutputFormatter) formatTemplate(templateStr string, data map[string]string) string {
	// First handle %{name} style placeholders
	result := templateStr
	for key, value := range data {
		placeholder := fmt.Sprintf("%%{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Then handle Go text/template style {{.name}} placeholders
	tmpl, err := template.New("format").Parse(result)
	if err != nil {
		// Fall back to simple replacement if template parsing fails
		return result
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		// Fall back to simple replacement if template execution fails
		return result
	}

	return buf.String()
}

// Regex pattern extraction methods
func (f *OutputFormatter) extractPatternData(pattern, text string) map[string]string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return make(map[string]string)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}

// FormatBackupWithExtraction formats a backup path using extracted data and a template.
func (f *OutputFormatter) FormatBackupWithExtraction(backupPath string) string {
	// Extract data from backup filename
	filename := getFilenameFromPath(backupPath)
	data := f.ExtractBackupFilenameData(filename)
	data["path"] = backupPath

	// Use template formatting if we have extracted data
	if len(data) > 1 { // More than just "path"
		return f.FormatCreatedBackupTemplate(data)
	}

	// Fall back to printf-style formatting
	return f.FormatCreatedBackup(backupPath)
}

// FormatListBackupWithExtraction formats a list backup message using extracted data from the filename.
func (f *OutputFormatter) FormatListBackupWithExtraction(backupPath, creationTime string) string {
	// Extract data from backup filename
	filename := getFilenameFromPath(backupPath)
	data := f.ExtractBackupFilenameData(filename)
	data["path"] = backupPath
	data["creation_time"] = creationTime

	// Use template formatting if we have extracted data
	if len(data) > 2 { // More than just "path" and "creation_time"
		return f.FormatListBackupTemplate(data)
	}

	// Fall back to printf-style formatting
	return f.FormatListBackup(backupPath, creationTime)
}

// TemplateFormatter provides advanced template processing with regex extraction
type TemplateFormatter struct {
	config *Config
}

// NewTemplateFormatter creates a new template formatter
func NewTemplateFormatter(cfg *Config) *TemplateFormatter {
	return &TemplateFormatter{config: cfg}
}

// FormatWithTemplate applies template to regex groups extracted from input
func (tf *TemplateFormatter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
	// Extract data using regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}

	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return tmplStr, nil // Return template as-is if no matches
	}

	// Build data map from named groups
	data := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			data[name] = matches[i]
		}
	}

	// Apply template formatting
	return tf.FormatWithPlaceholders(tmplStr, data), nil
}

// FormatWithPlaceholders replaces %{name} placeholders with data values
func (tf *TemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
	result := format

	// Handle %{name} style placeholders
	for key, value := range data {
		placeholder := fmt.Sprintf("%%{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Handle Go text/template style {{.name}} placeholders
	tmpl, err := template.New("format").Parse(result)
	if err != nil {
		// Fall back to simple replacement if template parsing fails
		return result
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		// Fall back to simple replacement if template execution fails
		return result
	}

	return buf.String()
}

// TemplateCreatedArchive formats a created archive message using template-based formatting.
func (tf *TemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedArchive, data)
}

// TemplateIdenticalArchive formats an identical archive message using template-based formatting.
func (tf *TemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalArchive, data)
}

// TemplateListArchive formats a list archive message using template-based formatting.
func (tf *TemplateFormatter) TemplateListArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListArchive, data)
}

// TemplateConfigValue formats a configuration value message using template-based formatting.
func (tf *TemplateFormatter) TemplateConfigValue(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateConfigValue, data)
}

// TemplateDryRunArchive formats a dry-run archive message using template-based formatting.
func (tf *TemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunArchive, data)
}

// TemplateError formats an error message using template-based formatting.
func (tf *TemplateFormatter) TemplateError(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateError, data)
}

// TemplateCreatedBackup formats a created backup message using template-based formatting.
func (tf *TemplateFormatter) TemplateCreatedBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedBackup, data)
}

// TemplateIdenticalBackup formats an identical backup message using template-based formatting.
func (tf *TemplateFormatter) TemplateIdenticalBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalBackup, data)
}

// TemplateListBackup formats a list backup message using template-based formatting.
func (tf *TemplateFormatter) TemplateListBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListBackup, data)
}

// TemplateDryRunBackup formats a dry-run backup message using template-based formatting.
func (tf *TemplateFormatter) TemplateDryRunBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunBackup, data)
}

// Print methods for template formatting

// PrintTemplateCreatedArchive prints a created archive message using template formatting.
func (tf *TemplateFormatter) PrintTemplateCreatedArchive(path string) {
	// Extract data from archive filename
	filename := getFilenameFromPath(path)
	data := tf.extractArchiveData(filename)
	data["path"] = path
	fmt.Print(tf.TemplateCreatedArchive(data))
}

// PrintTemplateCreatedBackup prints a created backup message using template formatting.
func (tf *TemplateFormatter) PrintTemplateCreatedBackup(path string) {
	// Extract data from backup filename
	filename := getFilenameFromPath(path)
	data := tf.extractBackupData(filename)
	data["path"] = path
	fmt.Print(tf.TemplateCreatedBackup(data))
}

// PrintTemplateListBackup prints a list backup message using template formatting.
func (tf *TemplateFormatter) PrintTemplateListBackup(path, creationTime string) {
	// Extract data from backup filename
	filename := getFilenameFromPath(path)
	data := tf.extractBackupData(filename)
	data["path"] = path
	data["creation_time"] = creationTime
	fmt.Print(tf.TemplateListBackup(data))
}

// PrintTemplateError prints an error message using template formatting.
func (tf *TemplateFormatter) PrintTemplateError(message, operation string) {
	data := map[string]string{
		"message":   message,
		"operation": operation,
	}
	fmt.Print(tf.TemplateError(data))
}

// Helper methods for data extraction

// extractArchiveData extracts data from archive filename using regex pattern
func (tf *TemplateFormatter) extractArchiveData(filename string) map[string]string {
	re, err := regexp.Compile(tf.config.PatternArchiveFilename)
	if err != nil {
		return make(map[string]string)
	}

	matches := re.FindStringSubmatch(filename)
	if matches == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}

// extractBackupData extracts data from backup filename using regex pattern
func (tf *TemplateFormatter) extractBackupData(filename string) map[string]string {
	re, err := regexp.Compile(tf.config.PatternBackupFilename)
	if err != nil {
		return make(map[string]string)
	}

	matches := re.FindStringSubmatch(filename)
	if matches == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}

// Enhanced formatting methods that use regex extraction for file backups
