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

// OUT-001: Delayed output management
// OutputMessage represents a message that can be displayed later
type OutputMessage struct {
	Content     string
	Destination string // "stdout" or "stderr"
	Type        string // "info", "error", "warning", etc.
}

// OutputCollector collects output messages for delayed display
type OutputCollector struct {
	messages []OutputMessage
}

// NewOutputCollector creates a new OutputCollector
func NewOutputCollector() *OutputCollector {
	return &OutputCollector{
		messages: make([]OutputMessage, 0),
	}
}

// AddStdout adds a stdout message to the collector
func (oc *OutputCollector) AddStdout(content, messageType string) {
	// OUT-001: Delayed output implementation
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stdout",
		Type:        messageType,
	})
}

// AddStderr adds a stderr message to the collector
func (oc *OutputCollector) AddStderr(content, messageType string) {
	// OUT-001: Delayed output implementation
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stderr",
		Type:        messageType,
	})
}

// GetMessages returns all collected messages
func (oc *OutputCollector) GetMessages() []OutputMessage {
	// OUT-001: Delayed output implementation
	return oc.messages
}

// FlushAll displays all collected messages and clears the collector
func (oc *OutputCollector) FlushAll() {
	// OUT-001: Delayed output implementation
	for _, msg := range oc.messages {
		if msg.Destination == "stderr" {
			fmt.Fprint(os.Stderr, msg.Content)
		} else {
			fmt.Print(msg.Content)
		}
	}
	oc.messages = make([]OutputMessage, 0)
}

// FlushStdout displays only stdout messages and removes them from the collector
func (oc *OutputCollector) FlushStdout() {
	// OUT-001: Delayed output implementation
	remaining := make([]OutputMessage, 0)
	for _, msg := range oc.messages {
		if msg.Destination == "stdout" {
			fmt.Print(msg.Content)
		} else {
			remaining = append(remaining, msg)
		}
	}
	oc.messages = remaining
}

// FlushStderr displays only stderr messages and removes them from the collector
func (oc *OutputCollector) FlushStderr() {
	// OUT-001: Delayed output implementation
	remaining := make([]OutputMessage, 0)
	for _, msg := range oc.messages {
		if msg.Destination == "stderr" {
			fmt.Fprint(os.Stderr, msg.Content)
		} else {
			remaining = append(remaining, msg)
		}
	}
	oc.messages = remaining
}

// Clear removes all collected messages without displaying them
func (oc *OutputCollector) Clear() {
	// OUT-001: Delayed output implementation
	oc.messages = make([]OutputMessage, 0)
}

// CFG-003: Output formatting interface
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// REFACTOR-001: Output formatting interface contracts defined
// REFACTOR-001: Config dependency interface required for extraction
// OutputFormatter provides methods for formatting and printing output for BkpDir operations.
// It supports both printf-style and template-based formatting, with optional delayed output.
type OutputFormatter struct {
	cfg       *Config
	collector *OutputCollector // OUT-001: Optional output collector for delayed display
}

// CFG-003: Output formatter constructor
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// NewOutputFormatter creates a new OutputFormatter with the given configuration.
// It initializes the formatter with the provided config for use in formatting operations.
func NewOutputFormatter(cfg *Config) *OutputFormatter {
	return &OutputFormatter{cfg: cfg, collector: nil}
}

// OUT-001: Output formatter constructor with collector
// NewOutputFormatterWithCollector creates a new OutputFormatter with delayed output support.
// Messages will be collected in the provided OutputCollector instead of being printed immediately.
func NewOutputFormatterWithCollector(cfg *Config, collector *OutputCollector) *OutputFormatter {
	return &OutputFormatter{cfg: cfg, collector: collector}
}

// OUT-001: Check if formatter is in delayed output mode
// IsDelayedMode returns true if the formatter is collecting output instead of printing immediately.
func (f *OutputFormatter) IsDelayedMode() bool {
	return f.collector != nil
}

// OUT-001: Get the output collector
// GetCollector returns the OutputCollector if in delayed mode, nil otherwise.
func (f *OutputFormatter) GetCollector() *OutputCollector {
	return f.collector
}

// OUT-001: Set or remove the output collector
// SetCollector sets the output collector for delayed output, or removes it if nil.
func (f *OutputFormatter) SetCollector(collector *OutputCollector) {
	f.collector = collector
}

// CFG-003: Printf-style archive creation formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatCreatedArchive formats a message for a created archive.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatCreatedArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedArchive, path)
}

// CFG-003: Printf-style identical archive formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalArchive formats a message for an identical archive.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatIdenticalArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalArchive, path)
}

// CFG-003: Printf-style archive listing formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListArchive formats a message for listing an archive.
// It uses the configured format string to create the output message with path and creation time.
func (f *OutputFormatter) FormatListArchive(path, creationTime string) string {
	return fmt.Sprintf(f.cfg.FormatListArchive, path, creationTime)
}

// CFG-003: Printf-style configuration value formatting
// IMMUTABLE-REF: Output Formatting Requirements, Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// FormatConfigValue formats a configuration value for display.
// It uses the configured format string to create the output message with name, value, and source.
func (f *OutputFormatter) FormatConfigValue(name, value, source string) string {
	return fmt.Sprintf(f.cfg.FormatConfigValue, name, value, source)
}

// CFG-003: Printf-style dry run formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatDryRunArchive formats a message for a dry-run archive operation.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatDryRunArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunArchive, path)
}

// CFG-003: Printf-style error formatting
// IMMUTABLE-REF: Output Formatting Requirements, Error Handling Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatError formats an error message.
// It uses the configured format string to create the error output message.
func (f *OutputFormatter) FormatError(message string) string {
	return fmt.Sprintf(f.cfg.FormatError, message)
}

// CFG-003: Output printing for archives
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintCreatedArchive prints a created archive message to stdout.
// It formats the message using FormatCreatedArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintCreatedArchive(path string) {
	message := f.FormatCreatedArchive(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Output printing for identical archives
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintIdenticalArchive prints an identical archive message to stdout.
// It formats the message using FormatIdenticalArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintIdenticalArchive(path string) {
	message := f.FormatIdenticalArchive(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Output printing for archive listings
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintListArchive prints a list archive message to stdout.
// It formats the message using FormatListArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintListArchive(path, creationTime string) {
	message := f.FormatListArchive(path, creationTime)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Output printing for configuration values
// IMMUTABLE-REF: Output Formatting Requirements, Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintConfigValue prints a config value message to stdout.
// It formats the message using FormatConfigValue and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintConfigValue(name, value, source string) {
	message := f.FormatConfigValue(name, value, source)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Output printing for dry run operations
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintDryRunArchive prints a dry-run archive message to stdout.
// It formats the message using FormatDryRunArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintDryRunArchive(path string) {
	message := f.FormatDryRunArchive(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Error output printing
// IMMUTABLE-REF: Output Formatting Requirements, Error Handling Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintError prints an error message to stderr.
// It formats the message using FormatError and writes it to stderr.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintError(message string) {
	errorMessage := f.FormatError(message)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(errorMessage, "error")
	} else {
		fmt.Fprint(os.Stderr, errorMessage)
	}
}

// CFG-003: Printf-style backup creation formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintCreatedBackup prints a created backup message to stdout.
// It formats the message using FormatCreatedBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintCreatedBackup(path string) {
	message := f.FormatCreatedBackup(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Printf-style identical backup formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintIdenticalBackup prints an identical backup message to stdout.
// It formats the message using FormatIdenticalBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintIdenticalBackup(path string) {
	message := f.FormatIdenticalBackup(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Printf-style backup listing formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintListBackup prints a list backup message to stdout.
// It formats the message using FormatListBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintListBackup(path, creationTime string) {
	message := f.FormatListBackup(path, creationTime)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Printf-style backup dry run formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: Enhanced with delayed output support
// PrintDryRunBackup prints a dry-run backup message to stdout.
// It formats the message using FormatDryRunBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintDryRunBackup(path string) {
	message := f.FormatDryRunBackup(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: Regex pattern data extraction for archives
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// ExtractArchiveFilenameData extracts data from an archive filename using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractArchiveFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternArchiveFilename, filename)
}

// CFG-003: Regex pattern data extraction for backups
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// ExtractBackupFilenameData extracts data from a backup filename using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractBackupFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternBackupFilename, filename)
}

// CFG-003: Regex pattern data extraction for config lines
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// ExtractConfigLineData extracts data from a config line using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractConfigLineData(line string) map[string]string {
	return f.extractPatternData(f.cfg.PatternConfigLine, line)
}

// CFG-003: Regex pattern data extraction for timestamps
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// ExtractTimestampData extracts data from a timestamp using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractTimestampData(timestamp string) map[string]string {
	return f.extractPatternData(f.cfg.PatternTimestamp, timestamp)
}

// CFG-003: Template-based archive formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatArchiveWithExtraction formats an archive message using template-based formatting.
// It extracts data from the archive filename and applies the configured template.
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

// CFG-003: Template-based archive listing formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListArchiveWithExtraction formats a list archive message using template-based formatting.
// It extracts data from the archive filename and applies the configured template.
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

// CFG-003: Printf-style backup creation formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatCreatedBackup formats a message for a created backup.
func (f *OutputFormatter) FormatCreatedBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedBackup, path)
}

// CFG-003: Printf-style identical backup formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalBackup formats a message for an identical backup.
func (f *OutputFormatter) FormatIdenticalBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalBackup, path)
}

// CFG-003: Printf-style backup listing formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListBackup formats a message for listing a backup.
func (f *OutputFormatter) FormatListBackup(path, creationTime string) string {
	return fmt.Sprintf(f.cfg.FormatListBackup, path, creationTime)
}

// CFG-003: Printf-style backup dry run formatting
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatDryRunBackup formats a message for a dry-run backup operation.
func (f *OutputFormatter) FormatDryRunBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunBackup, path)
}

// CFG-003: Template-based archive creation formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatCreatedArchiveTemplate formats a created archive message using a template.
func (f *OutputFormatter) FormatCreatedArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateCreatedArchive, data)
}

// CFG-003: Template-based identical archive formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalArchiveTemplate formats an identical archive message using a template.
func (f *OutputFormatter) FormatIdenticalArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIdenticalArchive, data)
}

// CFG-003: Template-based archive listing formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
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

// CFG-003: Template-based formatting interface
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateFormatter provides methods for template-based output formatting.
// It supports both pattern-based and placeholder-based template formatting.
type TemplateFormatter struct {
	config *Config
}

// CFG-003: Template formatter constructor
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// NewTemplateFormatter creates a new TemplateFormatter with the given configuration.
// It initializes the formatter with the provided config for use in template operations.
func NewTemplateFormatter(cfg *Config) *TemplateFormatter {
	return &TemplateFormatter{config: cfg}
}

// CFG-003: Pattern-based template formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatWithTemplate formats input using a pattern and template string.
// It extracts data using the pattern and applies the template to the extracted data.
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

// CFG-003: Placeholder-based template formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatWithPlaceholders formats a string using placeholder-based template formatting.
// It replaces placeholders in the format string with values from the data map.
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

// CFG-003: Template-based archive creation formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateCreatedArchive formats a created archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedArchive, data)
}

// CFG-003: Template-based identical archive formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateIdenticalArchive formats an identical archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalArchive, data)
}

// CFG-003: Template-based archive listing formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateListArchive formats a list archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateListArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListArchive, data)
}

// CFG-003: Template-based configuration value formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateConfigValue formats a configuration value message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateConfigValue(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateConfigValue, data)
}

// CFG-003: Template-based dry run formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateDryRunArchive formats a dry-run archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunArchive, data)
}

// CFG-003: Template-based error formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateError formats an error message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateError(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateError, data)
}

// CFG-003: Template-based backup creation formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateCreatedBackup formats a created backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateCreatedBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedBackup, data)
}

// CFG-003: Template-based identical backup formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateIdenticalBackup formats an identical backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateIdenticalBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalBackup, data)
}

// CFG-003: Template-based backup listing formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateListBackup formats a list backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateListBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListBackup, data)
}

// CFG-003: Template-based backup dry run formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateDryRunBackup formats a dry-run backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateDryRunBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunBackup, data)
}

// CFG-003: Template-based archive creation printing
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// PrintTemplateCreatedArchive prints a created archive message using template formatting.
// It extracts data from the archive filename and prints the formatted message to stdout.
func (tf *TemplateFormatter) PrintTemplateCreatedArchive(path string) {
	// Extract data from archive filename
	filename := getFilenameFromPath(path)
	data := tf.extractArchiveData(filename)
	data["path"] = path
	fmt.Print(tf.TemplateCreatedArchive(data))
}

// CFG-003: Template-based backup creation printing
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// PrintTemplateCreatedBackup prints a created backup message using template formatting.
// It extracts data from the backup filename and prints the formatted message to stdout.
func (tf *TemplateFormatter) PrintTemplateCreatedBackup(path string) {
	// Extract data from backup filename
	filename := getFilenameFromPath(path)
	data := tf.extractBackupData(filename)
	data["path"] = path
	fmt.Print(tf.TemplateCreatedBackup(data))
}

// CFG-003: Template-based backup listing printing
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// PrintTemplateListBackup prints a list backup message using template formatting.
// It extracts data from the backup filename and prints the formatted message to stdout.
func (tf *TemplateFormatter) PrintTemplateListBackup(path, creationTime string) {
	// Extract data from backup filename
	filename := getFilenameFromPath(path)
	data := tf.extractBackupData(filename)
	data["path"] = path
	data["creation_time"] = creationTime
	fmt.Print(tf.TemplateListBackup(data))
}

// CFG-003: Template-based error printing
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// PrintTemplateError prints an error message using template formatting.
// It formats the message with the provided operation and prints it to stderr.
func (tf *TemplateFormatter) PrintTemplateError(message, operation string) {
	data := map[string]string{
		"message":   message,
		"operation": operation,
	}
	fmt.Print(tf.TemplateError(data))
}

// CFG-003: Archive data extraction
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// extractArchiveData extracts data from an archive filename using regex patterns.
// It returns a map of named capture groups from the configured patterns.
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

// CFG-003: Backup data extraction
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// extractBackupData extracts data from a backup filename using regex patterns.
// It returns a map of named capture groups from the configured patterns.
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

// CFG-003: Template-based backup formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatBackupWithExtraction formats a backup message using template-based formatting.
// It extracts data from the backup filename and applies the configured template.
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

// CFG-003: Template-based backup listing formatting
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListBackupWithExtraction formats a list backup message using template-based formatting.
// It extracts data from the backup filename and applies the configured template.
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

// CFG-004: Incremental created formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatIncrementalCreated(path string) string {
	return fmt.Sprintf(f.cfg.FormatIncrementalCreated, path)
}

// Printf-style formatting methods for archive operations
// CFG-004: No archives found formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoArchivesFound(archiveDir string) string {
	return fmt.Sprintf(f.cfg.FormatNoArchivesFound, archiveDir)
}

// CFG-004: Verification failed formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationFailed(archiveName string, err error) string {
	return fmt.Sprintf(f.cfg.FormatVerificationFailed, archiveName, err)
}

// CFG-004: Verification success formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationSuccess(archiveName string) string {
	return fmt.Sprintf(f.cfg.FormatVerificationSuccess, archiveName)
}

// CFG-004: Verification warning formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationWarning(archiveName string, err error) string {
	return fmt.Sprintf(f.cfg.FormatVerificationWarning, archiveName, err)
}

// CFG-004: Configuration updated formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigurationUpdated(key string, value interface{}) string {
	return fmt.Sprintf(f.cfg.FormatConfigurationUpdated, key, value)
}

// CFG-004: Config file path formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigFilePath(path string) string {
	return fmt.Sprintf(f.cfg.FormatConfigFilePath, path)
}

// CFG-004: Dry run files header formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFilesHeader() string {
	return f.cfg.FormatDryRunFilesHeader
}

// CFG-004: Dry run file entry formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFileEntry(file string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunFileEntry, file)
}

// CFG-004: No files modified formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoFilesModified() string {
	return f.cfg.FormatNoFilesModified
}

// Printf-style formatting methods for backup operations
// CFG-004: No backups found formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoBackupsFound(filename, backupDir string) string {
	return fmt.Sprintf(f.cfg.FormatNoBackupsFound, filename, backupDir)
}

// CFG-004: Backup would create formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupWouldCreate(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupWouldCreate, path)
}

// CFG-004: Backup identical formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupIdentical(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupIdentical, path)
}

// CFG-004: Backup created formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupCreated(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupCreated, path)
}

// Template-based formatting methods for archive operations
// CFG-004: No archives found template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoArchivesFoundTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateNoArchivesFound, data)
}

// CFG-004: Verification failed template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationFailedTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateVerificationFailed, data)
}

// CFG-004: Verification success template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationSuccessTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateVerificationSuccess, data)
}

// CFG-004: Verification warning template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatVerificationWarningTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateVerificationWarning, data)
}

// CFG-004: Configuration updated template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigurationUpdatedTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateConfigurationUpdated, data)
}

// CFG-004: Config file path template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigFilePathTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateConfigFilePath, data)
}

// CFG-004: Dry run files header template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFilesHeaderTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateDryRunFilesHeader, data)
}

// CFG-004: Dry run file entry template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFileEntryTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateDryRunFileEntry, data)
}

// CFG-004: No files modified template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoFilesModifiedTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateNoFilesModified, data)
}

// CFG-004: Incremental created template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatIncrementalCreatedTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIncrementalCreated, data)
}

// Template-based formatting methods for backup operations
// CFG-004: No backups found template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoBackupsFoundTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateNoBackupsFound, data)
}

// CFG-004: Backup would create template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupWouldCreateTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateBackupWouldCreate, data)
}

// CFG-004: Backup identical template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupIdenticalTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateBackupIdentical, data)
}

// CFG-004: Backup created template formatting
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatBackupCreatedTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateBackupCreated, data)
}

// Print methods for CFG-004 format strings
// CFG-004: Print methods for archive operations
// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintNoArchivesFound(archiveDir string) {
	message := f.FormatNoArchivesFound(archiveDir)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintVerificationFailed(archiveName string, err error) {
	message := f.FormatVerificationFailed(archiveName, err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "error")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintVerificationSuccess(archiveName string) {
	message := f.FormatVerificationSuccess(archiveName)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintVerificationWarning(archiveName string, err error) {
	message := f.FormatVerificationWarning(archiveName, err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "warning")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintConfigurationUpdated(key string, value interface{}) {
	message := f.FormatConfigurationUpdated(key, value)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintConfigFilePath(path string) {
	message := f.FormatConfigFilePath(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintDryRunFilesHeader() {
	message := f.FormatDryRunFilesHeader()
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintDryRunFileEntry(file string) {
	message := f.FormatDryRunFileEntry(file)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintNoFilesModified() {
	message := f.FormatNoFilesModified()
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintIncrementalCreated(path string) {
	message := f.FormatIncrementalCreated(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-004: Print methods for backup operations
// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintNoBackupsFound(filename, backupDir string) {
	message := f.FormatNoBackupsFound(filename, backupDir)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintBackupWouldCreate(path string) {
	message := f.FormatBackupWouldCreate(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintBackupIdentical(path string) {
	message := f.FormatBackupIdentical(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintBackupCreated(path string) {
	message := f.FormatBackupCreated(path)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-004: Error message formatting methods
func (f *OutputFormatter) FormatDiskFullError(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatDiskFullError, err)
}

func (f *OutputFormatter) FormatPermissionError(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatPermissionError, err)
}

func (f *OutputFormatter) FormatDirectoryNotFound(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatDirectoryNotFound, err)
}

func (f *OutputFormatter) FormatFileNotFound(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFileNotFound, err)
}

func (f *OutputFormatter) FormatInvalidDirectory(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatInvalidDirectory, err)
}

func (f *OutputFormatter) FormatInvalidFile(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatInvalidFile, err)
}

func (f *OutputFormatter) FormatFailedWriteTemp(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedWriteTemp, err)
}

func (f *OutputFormatter) FormatFailedFinalizeFile(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedFinalizeFile, err)
}

func (f *OutputFormatter) FormatFailedCreateDirDisk(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedCreateDirDisk, err)
}

func (f *OutputFormatter) FormatFailedCreateDir(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedCreateDir, err)
}

func (f *OutputFormatter) FormatFailedAccessDir(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedAccessDir, err)
}

func (f *OutputFormatter) FormatFailedAccessFile(err error) string {
	// CFG-004: Implementation token
	return fmt.Sprintf(f.cfg.FormatFailedAccessFile, err)
}

// CFG-004: Template-based error message formatting methods
func (f *OutputFormatter) TemplateDiskFullError(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateDiskFullError, data)
}

func (f *OutputFormatter) TemplatePermissionError(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplatePermissionError, data)
}

func (f *OutputFormatter) TemplateDirectoryNotFound(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateDirectoryNotFound, data)
}

func (f *OutputFormatter) TemplateFileNotFound(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFileNotFound, data)
}

func (f *OutputFormatter) TemplateInvalidDirectory(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateInvalidDirectory, data)
}

func (f *OutputFormatter) TemplateInvalidFile(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateInvalidFile, data)
}

func (f *OutputFormatter) TemplateFailedWriteTemp(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedWriteTemp, data)
}

func (f *OutputFormatter) TemplateFailedFinalizeFile(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedFinalizeFile, data)
}

func (f *OutputFormatter) TemplateFailedCreateDirDisk(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedCreateDirDisk, data)
}

func (f *OutputFormatter) TemplateFailedCreateDir(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedCreateDir, data)
}

func (f *OutputFormatter) TemplateFailedAccessDir(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedAccessDir, data)
}

func (f *OutputFormatter) TemplateFailedAccessFile(err error) string {
	// CFG-004: Implementation token
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedAccessFile, data)
}

// CFG-004: Print methods for error messages
// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintDiskFullError(err error) {
	// CFG-004: Implementation token
	message := f.FormatDiskFullError(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintPermissionError(err error) {
	// CFG-004: Implementation token
	message := f.FormatPermissionError(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintDirectoryNotFound(err error) {
	// CFG-004: Implementation token
	message := f.FormatDirectoryNotFound(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFileNotFound(err error) {
	// CFG-004: Implementation token
	message := f.FormatFileNotFound(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintInvalidDirectory(err error) {
	// CFG-004: Implementation token
	message := f.FormatInvalidDirectory(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintInvalidFile(err error) {
	// CFG-004: Implementation token
	message := f.FormatInvalidFile(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedWriteTemp(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedWriteTemp(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedFinalizeFile(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedFinalizeFile(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedCreateDirDisk(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedCreateDirDisk(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedCreateDir(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedCreateDir(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedAccessDir(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedAccessDir(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintFailedAccessFile(err error) {
	// CFG-004: Implementation token
	message := f.FormatFailedAccessFile(err)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// CFG-004: Print method for verification error details
// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintVerificationErrorDetail(errMsg string) {
	message := fmt.Sprintf("  - %s\n", errMsg)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "error")
	} else {
		fmt.Print(message)
	}
}

// CFG-004: Print method for archive list with status
// OUT-001: Enhanced with delayed output support
func (f *OutputFormatter) PrintArchiveListWithStatus(output, status string) {
	message := fmt.Sprintf("%s%s\n", output, status)
	if f.collector != nil {
		// OUT-001: Delayed output implementation
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}
