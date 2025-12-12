// This file is part of bkpdir
//
// Package main provides output formatting for BkpDir CLI and tests.
// It handles printf-style and template-based formatting of output messages.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// SERVICE-TEMPLATE-001: Template Formatting Service Architecture Decision - Template formatting service implementation [ACTION:format-processing]
// Source: docs/context/architecture.md - Template Formatting Service section
// Impact: Template formatting service architecture decision

// SERVICE-OUTPUT-001: Output Formatting Service Architecture Decision - Output formatting service implementation [ACTION:format-processing]
// Source: docs/context/architecture.md - Output Formatting Service section
// Impact: Output formatting service architecture decision

// OUTPUT-FORMATTING-001: Output formatting specification - Output formatting and display [ACTION:format-processing]
// Source: formatter.go - OUTPUT-FORMATTING-001
// Impact: Core functionality requirement for output formatting

// SERVICE-FORMAT-001: Format service architecture decision - Format service implementation [ACTION:core-functionality]
// Source: formatter.go - SERVICE-FORMAT-001
// Impact: Format service implementation decision

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Component boundaries identified: OutputCollector, PrintfFormatter, TemplateFormatter, PatternExtractor, ErrorFormatter
// Ready for EXTRACT-003 (Output Formatting System) with config interface abstraction
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// These interfaces define contracts for clean component extraction

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Abstracts configuration dependency for formatter components
type FormatProvider interface {
	GetFormatString(formatType string) string
	GetTemplateString(templateType string) string
	GetPattern(patternType string) string
	GetErrorFormat(errorType string) string
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Abstracts output handling for formatter components
type OutputDestination interface {
	Print(message string)
	PrintError(message string)
	IsDelayedMode() bool
	SetCollector(collector *OutputCollector)
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Defines contract for regex-based data extraction
type PatternExtractor interface {
	ExtractArchiveFilenameData(filename string) map[string]string
	ExtractBackupFilenameData(filename string) map[string]string
	ExtractPatternData(pattern, text string) map[string]string
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Primary formatter interface for printf-style formatting
type FormatterInterface interface {
	FormatCreatedArchive(path string) string
	FormatIdenticalArchive(path string) string
	FormatListArchive(path, creationTime string) string
	FormatConfigValue(name, value, source string) string
	FormatError(message string) string
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Interface for template-based formatting operations
type TemplateFormatterInterface interface {
	FormatWithTemplate(input, pattern, tmplStr string) (string, error)
	FormatWithPlaceholders(format string, data map[string]string) string
	TemplateCreatedArchive(data map[string]string) string
	TemplateIdenticalArchive(data map[string]string) string
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// OutputMessage represents a message that can be displayed later
type OutputMessage struct {
	Content     string
	Destination string // "stdout" or "stderr"
	Type        string // "info", "error", "warning", etc.
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Output collector ready for immediate extraction
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
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stdout",
		Type:        messageType,
	})
}

// AddStderr adds a stderr message to the collector
func (oc *OutputCollector) AddStderr(content, messageType string) {
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stderr",
		Type:        messageType,
	})
}

// GetMessages returns all collected messages
func (oc *OutputCollector) GetMessages() []OutputMessage {
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
	return oc.messages
}

// FlushAll displays all collected messages and clears the collector
func (oc *OutputCollector) FlushAll() {
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
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
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
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
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
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
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
	oc.messages = make([]OutputMessage, 0)
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Printf Formatter Component (Lines 120-610) - [NOTE] [DECISION:maintenance]
// Configuration dependency requires interface abstraction for extraction
// OutputFormatter provides methods for formatting and printing output for BkpDir operations.
// It supports both printf-style and template-based formatting, with optional delayed output.
type OutputFormatter struct {
	cfg *Config
	// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
	collector *OutputCollector
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// IsDelayedMode returns true if the formatter is collecting output instead of printing immediately.
func (f *OutputFormatter) IsDelayedMode() bool {
	return f.collector != nil
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// GetCollector returns the OutputCollector if in delayed mode, nil otherwise.
func (f *OutputFormatter) GetCollector() *OutputCollector {
	return f.collector
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// SetCollector sets the output collector for delayed output, or removes it if nil.
func (f *OutputFormatter) SetCollector(collector *OutputCollector) {
	f.collector = collector
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Core Printf Formatters (Lines 166-226) - [NOTE] [DECISION:maintenance]
// Direct config dependency - format string access needs interface abstraction
// FormatCreatedArchive formats a message for a created archive.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatCreatedArchive(path string) string {
	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	result := fmt.Sprintf(f.cfg.FormatCreatedArchive, path)
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: FormatCreatedArchive called with path: %s\n", path)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: Format string: %q\n", f.cfg.FormatCreatedArchive)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: Result: %q\n", result)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	return result
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalArchive formats a message for an identical archive.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatIdenticalArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalArchive, path)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListArchive formats a message for listing an archive.
// It uses the configured format string to create the output message with path and creation time.
// FormatListArchive formats a list archive message.
// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Supports both printf-style (%s) and template-style (#{name}) placeholders.
// If template placeholders are detected, gathers file statistics and uses template formatting.
// Otherwise, uses printf formatting for backward compatibility.
func (f *OutputFormatter) FormatListArchive(path, creationTime string) string {
	formatStr := f.cfg.FormatListArchive

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Check if format string contains template placeholders
	hasTemplatePlaceholders := strings.Contains(formatStr, "#{")

	if hasTemplatePlaceholders {
		// Use template formatting with file statistics
		data := make(map[string]string)
		data["path"] = path
		data["creation_time"] = creationTime

		// Gather file statistics to populate size_human and other stat fields
		statInfo, err := GatherFileStatInfo(path)
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
			// Provide default values for file attributes when stats can't be gathered
			data["size"] = "0"
			data["size_human"] = "unknown"
			data["mtime"] = creationTime
			data["mtime_unix"] = "0"
			data["mode"] = "unknown"
			data["type"] = "unknown"
			// Extract filename from path
			parts := strings.Split(path, "/")
			if len(parts) > 0 {
				data["name"] = parts[len(parts)-1]
			} else {
				data["name"] = path
			}
		}

		// Use template formatter to process placeholders
		return f.formatTemplate(formatStr, data)
	}

	// Use printf formatting for backward compatibility
	return fmt.Sprintf(formatStr, path, creationTime)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements, Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// FormatConfigValue formats a configuration value for display.
// It uses the configured format string to create the output message with name, value, and source.
func (f *OutputFormatter) FormatConfigValue(name, value, source string) string {
	return fmt.Sprintf(f.cfg.FormatConfigValue, name, value, source)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatDryRunArchive formats a message for a dry-run archive operation.
// It uses the configured format string to create the output message.
func (f *OutputFormatter) FormatDryRunArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunArchive, path)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements, Error Handling Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatError formats an error message.
// It uses the configured format string to create the error output message.
func (f *OutputFormatter) FormatError(message string) string {
	return fmt.Sprintf(f.cfg.FormatError, message)
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Print Output Methods (Lines 228-405) - [NOTE] [DECISION:maintenance]
// Format + Print with optional delayed output via collector
// PrintCreatedArchive prints a message for a created archive.
// Uses delayed output if collector is set, otherwise prints immediately.
func (f *OutputFormatter) PrintCreatedArchive(path string) {
	message := f.FormatCreatedArchive(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintIdenticalArchive prints an identical archive message to stdout.
// It formats the message using FormatIdenticalArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintIdenticalArchive(path string) {
	message := f.FormatIdenticalArchive(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintListArchive prints a list archive message to stdout.
// It formats the message using FormatListArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintListArchive(path, creationTime string) {
	message := f.FormatListArchive(path, creationTime)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements, Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintConfigValue prints a config value message to stdout.
// It formats the message using FormatConfigValue and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintConfigValue(name, value, source string) {
	message := f.FormatConfigValue(name, value, source)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintDryRunArchive prints a dry-run archive message to stdout.
// It formats the message using FormatDryRunArchive and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintDryRunArchive(path string) {
	message := f.FormatDryRunArchive(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements, Error Handling Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintError prints an error message to stderr.
// It formats the message using FormatError and writes it to stderr.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintError(message string) {
	errorMessage := f.FormatError(message)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(errorMessage, "error")
	} else {
		fmt.Fprint(os.Stderr, errorMessage)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintCreatedBackup prints a created backup message to stdout.
// It formats the message using FormatCreatedBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintCreatedBackup(path string) {
	message := f.FormatCreatedBackup(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintIdenticalBackup prints an identical backup message to stdout.
// It formats the message using FormatIdenticalBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintIdenticalBackup(path string) {
	message := f.FormatIdenticalBackup(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintListBackup prints a list backup message to stdout.
// It formats the message using FormatListBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintListBackup(path, creationTime string) {
	message := f.FormatListBackup(path, creationTime)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
// PrintDryRunBackup prints a dry-run backup message to stdout.
// It formats the message using FormatDryRunBackup and writes it to stdout.
// If in delayed mode, the message is collected instead of printed immediately.
func (f *OutputFormatter) PrintDryRunBackup(path string) {
	message := f.FormatDryRunBackup(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Pattern Extraction Methods (Lines 406-482) - [NOTE] [DECISION:maintenance]
// Regex-based data extraction - shared functionality
// ExtractArchiveFilenameData extracts data from archive filename patterns.
func (f *OutputFormatter) ExtractArchiveFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternArchiveFilename, filename)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// ExtractBackupFilenameData extracts data from a backup filename using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractBackupFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternBackupFilename, filename)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-003
// ExtractConfigLineData extracts data from a config line using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractConfigLineData(line string) map[string]string {
	return f.extractPatternData(f.cfg.PatternConfigLine, line)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// ExtractTimestampData extracts data from a timestamp using a regex pattern.
// It returns a map of named capture groups from the configured pattern.
func (f *OutputFormatter) ExtractTimestampData(timestamp string) map[string]string {
	return f.extractPatternData(f.cfg.PatternTimestamp, timestamp)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

	// Gather file statistics to populate size_human and other stat fields
	statInfo, err := GatherFileStatInfo(archivePath)
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
		// If stats can't be gathered, provide default values for common placeholders
		// to prevent template placeholders from appearing literally in output
		if _, exists := data["size_human"]; !exists {
			data["size_human"] = "unknown"
		}
		if _, exists := data["size"]; !exists {
			data["size"] = "0"
		}
		if _, exists := data["mtime"]; !exists {
			data["mtime"] = creationTime // Use creation_time as fallback
		}
	}

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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatCreatedBackup formats a message for a created backup.
func (f *OutputFormatter) FormatCreatedBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedBackup, path)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalBackup formats a message for an identical backup.
func (f *OutputFormatter) FormatIdenticalBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalBackup, path)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatListBackup formats a message for listing a backup.
// FormatListBackup formats a list backup message.
// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Supports both printf-style (%s) and template-style (#{name}) placeholders.
// If template placeholders are detected, gathers file statistics and uses template formatting.
// Otherwise, uses printf formatting for backward compatibility.
func (f *OutputFormatter) FormatListBackup(path, creationTime string) string {
	formatStr := f.cfg.FormatListBackup

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Check if format string contains template placeholders
	hasTemplatePlaceholders := strings.Contains(formatStr, "#{")

	if hasTemplatePlaceholders {
		// Use template formatting with file statistics
		data := make(map[string]string)
		data["path"] = path
		data["creation_time"] = creationTime

		// Gather file statistics to populate size_human and other stat fields
		statInfo, err := GatherFileStatInfo(path)
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
			// Provide default values for file attributes when stats can't be gathered
			data["size"] = "0"
			data["size_human"] = "unknown"
			data["mtime"] = creationTime
			data["mtime_unix"] = "0"
			data["mode"] = "unknown"
			data["type"] = "unknown"
			// Extract filename from path
			parts := strings.Split(path, "/")
			if len(parts) > 0 {
				data["name"] = parts[len(parts)-1]
			} else {
				data["name"] = path
			}
		}

		// Use template formatter to process placeholders
		return f.formatTemplate(formatStr, data)
	}

	// Use printf formatting for backward compatibility
	return fmt.Sprintf(formatStr, path, creationTime)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Output Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatDryRunBackup formats a message for a dry-run backup operation.
func (f *OutputFormatter) FormatDryRunBackup(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunBackup, path)
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Template Integration Methods (Lines 532-609) - [NOTE] [DECISION:maintenance]
// Bridge between printf and template systems
// FormatCreatedArchiveTemplate formats using template with extracted data.
func (f *OutputFormatter) FormatCreatedArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateCreatedArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// FormatIdenticalArchiveTemplate formats an identical archive message using a template.
func (f *OutputFormatter) FormatIdenticalArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIdenticalArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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
	// First handle #{name} style placeholders - replace ALL placeholders from data map first
	result := templateStr
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Replace known placeholders that might be missing with defaults
	// This handles cases where placeholders weren't in the data map (e.g., missing file stats)
	// CRITICAL: This must happen AFTER the data map replacement to catch any that weren't replaced
	// Only replace known placeholders, leave unknown ones as-is
	knownPlaceholders := map[string]string{
		"#{size_human}": "unknown",
		"#{size}":       "0",
	}
	// Add mtime default if creation_time is available
	if creationTime, ok := data["creation_time"]; ok {
		knownPlaceholders["#{mtime}"] = creationTime
	} else {
		knownPlaceholders["#{mtime}"] = "unknown"
	}

	// Replace only known placeholders that are still present (weren't replaced by data map)
	for placeholder, defaultValue := range knownPlaceholders {
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, defaultValue)
		}
	}

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Handle printf-style placeholders (%s, %d, etc.) after template placeholders are replaced
	// This allows mixed format strings like "%s (size: #{size_human})\n"
	// We need to replace printf placeholders with values from the data map
	// Common mappings: %s -> path or name, %d -> size, etc.
	if strings.Contains(result, "%s") {
		// Replace %s with path if available, otherwise name
		if path, ok := data["path"]; ok {
			result = strings.Replace(result, "%s", path, 1) // Replace first occurrence
		} else if name, ok := data["name"]; ok {
			result = strings.Replace(result, "%s", name, 1)
		}
		// If there are multiple %s, replace them with available values
		// Second %s could be creation_time
		if strings.Contains(result, "%s") {
			if creationTime, ok := data["creation_time"]; ok {
				result = strings.Replace(result, "%s", creationTime, 1)
			}
		}
	}

	// CRITICAL: Final safety check - replace any remaining #{...} patterns
	// This prevents fmt.Sprintf from misinterpreting them as format verbs
	// We MUST ensure no #{...} patterns remain before returning
	if strings.Contains(result, "#{") {
		// Replace any remaining known placeholders with safe defaults
		finalReplacements := map[string]string{
			"#{path}": func() string {
				if p, ok := data["path"]; ok {
					return p
				}
				return "unknown"
			}(),
			"#{name}": func() string {
				if n, ok := data["name"]; ok {
					return n
				}
				return "unknown"
			}(),
			"#{creation_time}": func() string {
				if ct, ok := data["creation_time"]; ok {
					return ct
				}
				return "unknown"
			}(),
			"#{size_human}": "unknown",
			"#{size}":       "0",
			"#{mtime}": func() string {
				if ct, ok := data["creation_time"]; ok {
					return ct
				}
				return "unknown"
			}(),
			"#{mtime_unix}": "0",
			"#{mode}":       "unknown",
			"#{type}":       "unknown",
		}
		for placeholder, value := range finalReplacements {
			if strings.Contains(result, placeholder) {
				result = strings.ReplaceAll(result, placeholder, value)
			}
		}
	}

	// CRITICAL: Do NOT call tmpl.Execute if result still contains #{...} patterns
	// Go's text/template uses fmt internally which will misinterpret #{...} as format verbs
	// Only parse as Go template if there are {{.}} patterns AND no #{...} patterns remain
	hasGoTemplatePatterns := strings.Contains(result, "{{") && strings.Contains(result, "}}")
	hasPlaceholderPatterns := strings.Contains(result, "#{")

	if hasGoTemplatePatterns && !hasPlaceholderPatterns {
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

	return result
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Pattern Extraction Methods (Lines 406-482) - [NOTE] [DECISION:maintenance]
// Regex-based data extraction - shared functionality
// ExtractArchiveFilenameData extracts data from archive filename patterns.
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

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Template Formatter Component (Lines 637-928) - [NOTE] [DECISION:maintenance]
// Configuration dependency requires interface abstraction for extraction
// TemplateFormatter provides methods for template-based output formatting.
// It supports both pattern-based and placeholder-based template formatting.
type TemplateFormatter struct {
	config *Config
}

// NewTemplateFormatter creates a new TemplateFormatter with the given configuration.
// It initializes the formatter with the provided config for use in template operations.
func NewTemplateFormatter(cfg *Config) *TemplateFormatter {
	return &TemplateFormatter{config: cfg}
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Template Engine Core (Lines 657-717) - [NOTE] [DECISION:maintenance]
// Self-contained template processing with pattern extraction
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

// FormatWithPlaceholders formats a string using placeholder-based template formatting.
// It replaces placeholders in the format string with values from the data map.
func (tf *TemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
	result := format

	// Handle #{name} style placeholders - replace ALL placeholders from data map first
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] Replace known placeholders that might be missing with defaults
	// This handles cases where placeholders weren't in the data map (e.g., missing file stats)
	// CRITICAL: This must happen AFTER the data map replacement to catch any that weren't replaced
	// Only replace known placeholders, leave unknown ones as-is
	knownPlaceholders := map[string]string{
		"#{size_human}": "unknown",
		"#{size}":       "0",
	}
	// Add mtime default if creation_time is available
	if creationTime, ok := data["creation_time"]; ok {
		knownPlaceholders["#{mtime}"] = creationTime
	} else {
		knownPlaceholders["#{mtime}"] = "unknown"
	}

	// Replace only known placeholders that are still present (weren't replaced by data map)
	for placeholder, defaultValue := range knownPlaceholders {
		if strings.Contains(result, placeholder) {
			result = strings.ReplaceAll(result, placeholder, defaultValue)
		}
	}

	// CRITICAL: Final safety check - replace any remaining #{...} patterns
	// This prevents fmt.Sprintf from misinterpreting them as format verbs
	// We MUST ensure no #{...} patterns remain before returning
	if strings.Contains(result, "#{") {
		// Replace any remaining known placeholders with safe defaults
		finalReplacements := map[string]string{
			"#{path}": func() string {
				if p, ok := data["path"]; ok {
					return p
				}
				return "unknown"
			}(),
			"#{name}": func() string {
				if n, ok := data["name"]; ok {
					return n
				}
				return "unknown"
			}(),
			"#{creation_time}": func() string {
				if ct, ok := data["creation_time"]; ok {
					return ct
				}
				return "unknown"
			}(),
			"#{size_human}": "unknown",
			"#{size}":       "0",
			"#{mtime}": func() string {
				if ct, ok := data["creation_time"]; ok {
					return ct
				}
				return "unknown"
			}(),
			"#{mtime_unix}": "0",
			"#{mode}":       "unknown",
			"#{type}":       "unknown",
		}
		for placeholder, value := range finalReplacements {
			if strings.Contains(result, placeholder) {
				result = strings.ReplaceAll(result, placeholder, value)
			}
		}
	}

	// CRITICAL: Do NOT call tmpl.Execute if result still contains #{...} patterns
	// Go's text/template uses fmt internally which will misinterpret #{...} as format verbs
	// Only parse as Go template if there are {{.}} patterns AND no #{...} patterns remain
	hasGoTemplatePatterns := strings.Contains(result, "{{") && strings.Contains(result, "}}")
	hasPlaceholderPatterns := strings.Contains(result, "#{")

	if hasGoTemplatePatterns && !hasPlaceholderPatterns {
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

	return result
}

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Template Method Series (Lines 718-817) - [NOTE] [DECISION:maintenance]
// Direct config template dependency - needs interface abstraction
// TemplateCreatedArchive formats a created archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateIdenticalArchive formats an identical archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateListArchive formats a list archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateListArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateConfigValue formats a configuration value message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateConfigValue(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateConfigValue, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateDryRunArchive formats a dry-run archive message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunArchive, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateError formats an error message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateError(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateError, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateCreatedBackup formats a created backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateCreatedBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateCreatedBackup, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateIdenticalBackup formats an identical backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateIdenticalBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateIdenticalBackup, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateListBackup formats a list backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateListBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateListBackup, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// TemplateDryRunBackup formats a dry-run backup message using a template.
// It applies the configured template to the provided data map.
func (tf *TemplateFormatter) TemplateDryRunBackup(data map[string]string) string {
	return tf.FormatWithPlaceholders(tf.config.TemplateDryRunBackup, data)
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
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

// REFACTOR-002: See architecture.md - Formatter Decomposition [DECISION:maintenance]
// Extended Printf Formatters (Lines 929-1084) - [CHECK] [DECISION:maintenance]
// Complex formatting requiring data extraction - extends core printf functionality
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

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatIncrementalCreated(path string) string {
	return fmt.Sprintf(f.cfg.FormatIncrementalCreated, path)
}

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
// FormatCreatedArchiveWithStats formats a created archive message with file statistics using named replacements
func (f *OutputFormatter) FormatCreatedArchiveWithStats(path string) string {
	// Gather file statistics
	statInfo, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic format if stat gathering fails
		return f.FormatCreatedArchive(path)
	}

	// Create data map for named replacements
	data := map[string]string{
		"path":       statInfo.Path,
		"name":       statInfo.Name,
		"size":       fmt.Sprintf("%d", statInfo.Size),
		"size_human": statInfo.SizeHuman,
		"mtime":      statInfo.MTime.Format("2006-01-02 15:04:05"),
		"mtime_unix": fmt.Sprintf("%d", statInfo.MTimeUnix),
		"mode":       statInfo.Mode.String(),
		"type":       statInfo.Type,
	}

	// Use detailed template string with named replacements
	result := f.formatTemplate(f.cfg.TemplateCreatedArchiveDetailed, data)
	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	if strings.Contains(result, "#{") {
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: Template not processed correctly: %q\n", result)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: Template input was: %q\n", f.cfg.TemplateCreatedArchiveDetailed)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: Data was: %+v\n", data)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	}
	return result
}

// FormatIncrementalCreatedWithStats formats an incremental archive message with file statistics using named replacements
func (f *OutputFormatter) FormatIncrementalCreatedWithStats(path string) string {
	// Gather file statistics
	statInfo, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic format if stat gathering fails
		return f.FormatIncrementalCreated(path)
	}

	// Create data map for named replacements
	data := map[string]string{
		"path":       statInfo.Path,
		"name":       statInfo.Name,
		"size":       fmt.Sprintf("%d", statInfo.Size),
		"size_human": statInfo.SizeHuman,
		"mtime":      statInfo.MTime.Format("2006-01-02 15:04:05"),
		"mtime_unix": fmt.Sprintf("%d", statInfo.MTimeUnix),
		"mode":       statInfo.Mode.String(),
		"type":       statInfo.Type,
	}

	// Use detailed format string with named replacements
	return f.formatTemplate(f.cfg.FormatIncrementalCreatedDetailed, data)
}

// PrintCreatedArchiveWithStats prints a created archive message with file statistics
func (f *OutputFormatter) PrintCreatedArchiveWithStats(path string) {
	message := f.FormatCreatedArchiveWithStats(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// PrintIncrementalCreatedWithStats prints an incremental archive message with file statistics
func (f *OutputFormatter) PrintIncrementalCreatedWithStats(path string) {
	message := f.FormatIncrementalCreatedWithStats(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
// TemplateCreatedArchiveWithStats formats a created archive message using template with file statistics
func (f *OutputFormatter) TemplateCreatedArchiveWithStats(path string) string {
	// Gather file statistics
	statInfo, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat gathering fails
		data := map[string]string{"path": path}
		return f.formatTemplate(f.cfg.TemplateCreatedArchive, data)
	}

	// Create data map for template processing
	data := map[string]string{
		"path":       statInfo.Path,
		"name":       statInfo.Name,
		"size":       fmt.Sprintf("%d", statInfo.Size),
		"size_human": statInfo.SizeHuman,
		"mtime":      statInfo.MTime.Format("2006-01-02 15:04:05"),
		"mtime_unix": fmt.Sprintf("%d", statInfo.MTimeUnix),
		"mode":       statInfo.Mode.String(),
		"type":       statInfo.Type,
	}

	// Use detailed template with named replacements
	return f.formatTemplate(f.cfg.TemplateCreatedArchiveDetailed, data)
}

// TemplateIncrementalCreatedWithStats formats an incremental archive message using template with file statistics
func (f *OutputFormatter) TemplateIncrementalCreatedWithStats(path string) string {
	// Gather file statistics
	statInfo, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat gathering fails
		data := map[string]string{"path": path}
		return f.formatTemplate(f.cfg.TemplateIncrementalCreated, data)
	}

	// Create data map for template processing
	data := map[string]string{
		"path":       statInfo.Path,
		"name":       statInfo.Name,
		"size":       fmt.Sprintf("%d", statInfo.Size),
		"size_human": statInfo.SizeHuman,
		"mtime":      statInfo.MTime.Format("2006-01-02 15:04:05"),
		"mtime_unix": fmt.Sprintf("%d", statInfo.MTimeUnix),
		"mode":       statInfo.Mode.String(),
		"type":       statInfo.Type,
	}

	// Use detailed template string with named replacements
	return f.formatTemplate(f.cfg.TemplateIncrementalCreatedDetailed, data)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoArchivesFound(archiveDir string) string {
	return fmt.Sprintf(f.cfg.FormatNoArchivesFound, archiveDir)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigurationUpdated(key string, value interface{}) string {
	return fmt.Sprintf(f.cfg.FormatConfigurationUpdated, key, value)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatConfigFilePath(path string) string {
	return fmt.Sprintf(f.cfg.FormatConfigFilePath, path)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFilesHeader() string {
	return f.cfg.FormatDryRunFilesHeader
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatDryRunFileEntry(file string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunFileEntry, file)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) FormatNoFilesModified() string {
	return f.cfg.FormatNoFilesModified
}

// Printf-style formatting methods for archive operations
// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// IMMUTABLE-REF: String externalization requirements
// TEST-REF: TestStringExternalization
// DECISION-REF: DEC-009
func (f *OutputFormatter) PrintNoArchivesFound(archiveDir string) {
	message := f.FormatNoArchivesFound(archiveDir)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintConfigurationUpdated(key string, value interface{}) {
	message := f.FormatConfigurationUpdated(key, value)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintConfigFilePath(path string) {
	message := f.FormatConfigFilePath(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "config")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintDryRunFilesHeader() {
	message := f.FormatDryRunFilesHeader()
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintDryRunFileEntry(file string) {
	message := f.FormatDryRunFileEntry(file)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintNoFilesModified() {
	message := f.FormatNoFilesModified()
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintIncrementalCreated(path string) {
	message := f.FormatIncrementalCreated(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
func (f *OutputFormatter) FormatNoBackupsFound(filename, backupDir string) string {
	return fmt.Sprintf(f.cfg.FormatNoBackupsFound, filename, backupDir)
}

func (f *OutputFormatter) FormatBackupWouldCreate(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupWouldCreate, path)
}

func (f *OutputFormatter) FormatBackupIdentical(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupIdentical, path)
}

func (f *OutputFormatter) FormatBackupCreated(path string) string {
	return fmt.Sprintf(f.cfg.FormatBackupCreated, path)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintNoBackupsFound(filename, backupDir string) {
	message := f.FormatNoBackupsFound(filename, backupDir)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintBackupWouldCreate(path string) {
	message := f.FormatBackupWouldCreate(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "dry-run")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintBackupIdentical(path string) {
	message := f.FormatBackupIdentical(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintBackupCreated(path string) {
	message := f.FormatBackupCreated(path)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		fmt.Print(message)
	}
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
func (f *OutputFormatter) FormatDiskFullError(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatDiskFullError, err)
}

func (f *OutputFormatter) FormatPermissionError(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatPermissionError, err)
}

func (f *OutputFormatter) FormatDirectoryNotFound(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatDirectoryNotFound, err)
}

func (f *OutputFormatter) FormatFileNotFound(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFileNotFound, err)
}

func (f *OutputFormatter) FormatInvalidDirectory(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatInvalidDirectory, err)
}

func (f *OutputFormatter) FormatInvalidFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatInvalidFile, err)
}

func (f *OutputFormatter) FormatFailedWriteTemp(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedWriteTemp, err)
}

func (f *OutputFormatter) FormatFailedFinalizeFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedFinalizeFile, err)
}

func (f *OutputFormatter) FormatFailedCreateDirDisk(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedCreateDirDisk, err)
}

func (f *OutputFormatter) FormatFailedCreateDir(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedCreateDir, err)
}

func (f *OutputFormatter) FormatFailedAccessDir(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedAccessDir, err)
}

func (f *OutputFormatter) FormatFailedAccessFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	return fmt.Sprintf(f.cfg.FormatFailedAccessFile, err)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
func (f *OutputFormatter) TemplateDiskFullError(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateDiskFullError, data)
}

func (f *OutputFormatter) TemplatePermissionError(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplatePermissionError, data)
}

func (f *OutputFormatter) TemplateDirectoryNotFound(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateDirectoryNotFound, data)
}

func (f *OutputFormatter) TemplateFileNotFound(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFileNotFound, data)
}

func (f *OutputFormatter) TemplateInvalidDirectory(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateInvalidDirectory, data)
}

func (f *OutputFormatter) TemplateInvalidFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateInvalidFile, data)
}

func (f *OutputFormatter) TemplateFailedWriteTemp(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedWriteTemp, data)
}

func (f *OutputFormatter) TemplateFailedFinalizeFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedFinalizeFile, data)
}

func (f *OutputFormatter) TemplateFailedCreateDirDisk(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedCreateDirDisk, data)
}

func (f *OutputFormatter) TemplateFailedCreateDir(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedCreateDir, data)
}

func (f *OutputFormatter) TemplateFailedAccessDir(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedAccessDir, data)
}

func (f *OutputFormatter) TemplateFailedAccessFile(err error) string {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	data := map[string]string{
		"error": err.Error(),
	}
	return f.formatTemplate(f.cfg.TemplateFailedAccessFile, data)
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintDiskFullError(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatDiskFullError(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintPermissionError(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatPermissionError(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintDirectoryNotFound(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatDirectoryNotFound(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFileNotFound(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFileNotFound(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintInvalidDirectory(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatInvalidDirectory(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintInvalidFile(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatInvalidFile(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedWriteTemp(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedWriteTemp(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedFinalizeFile(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedFinalizeFile(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedCreateDirDisk(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedCreateDirDisk(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedCreateDir(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedCreateDir(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedAccessDir(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedAccessDir(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintFailedAccessFile(err error) {
	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	message := f.FormatFailedAccessFile(err)
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStderr(message, "error")
	} else {
		fmt.Fprint(os.Stderr, message)
	}
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
func (f *OutputFormatter) PrintArchiveListWithStatus(output, status string) {
	// [REQ:CUSTOMIZABLE_FORMAT_STRINGS] CRITICAL: Use os.Stdout.WriteString, NOT fmt.Print
	// to avoid fmt misinterpreting #{...} patterns as format verbs
	// fmt.Print internally uses fmt.Sprintf("%v", ...) which can misinterpret #{...} as format verbs
	message := output + status + "\n"
	if f.collector != nil {
		// OUT-001: See specification.md - Delayed Output [DECISION:maintenance]
		f.collector.AddStdout(message, "info")
	} else {
		os.Stdout.WriteString(message)
	}
}
