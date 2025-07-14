// AI-First Core Formatter Implementation
// Provides pure formatting operations without side effects for optimal AI assistant comprehension.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:core-functionality]
package formatter

import (
	"fmt"
	"strings"
	"text/template"
)

// [CRITICAL] FMT-001: AI-first core formatter implementation - [ACTION:core-functionality]
type AICoreFormatter struct {
	config FormatterConfig
}

// [CRITICAL] FMT-001: AI-first core formatter constructor - [ACTION:core-functionality]
func NewAICoreFormatter(config FormatterConfig) *AICoreFormatter {
	return &AICoreFormatter{
		config: config,
	}
}

// [CRITICAL] FMT-001: Archive formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatArchive(path string, formatType FormatType) (string, error) {
	formatStr, err := f.config.GetFormatString(formatType)
	if err != nil {
		return "", fmt.Errorf("failed to get format string for type %s: %w", formatType, err)
	}

	if formatStr == "" {
		// Provide default format strings
		switch formatType {
		case FormatTypeCreated:
			formatStr = "Created archive: %s\n"
		case FormatTypeIdentical:
			formatStr = "Identical archive exists: %s\n"
		case FormatTypeList:
			formatStr = "%s (created: %s)\n"
		case FormatTypeDryRun:
			formatStr = "Would create archive: %s\n"
		default:
			return "", fmt.Errorf("unsupported format type: %s", formatType)
		}
	}

	return fmt.Sprintf(formatStr, path), nil
}

// [CRITICAL] FMT-001: Backup formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatBackup(path string, formatType FormatType) (string, error) {
	formatStr, err := f.config.GetFormatString(formatType)
	if err != nil {
		return "", fmt.Errorf("failed to get format string for type %s: %w", formatType, err)
	}

	if formatStr == "" {
		// Provide default format strings
		switch formatType {
		case FormatTypeCreated:
			formatStr = "Created backup: %s\n"
		case FormatTypeIdentical:
			formatStr = "Identical backup exists: %s\n"
		case FormatTypeList:
			formatStr = "%s (created: %s)\n"
		case FormatTypeDryRun:
			formatStr = "Would create backup: %s\n"
		default:
			return "", fmt.Errorf("unsupported format type: %s", formatType)
		}
	}

	return fmt.Sprintf(formatStr, path), nil
}

// [CRITICAL] FMT-001: Config formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatConfig(name, value, source string) (string, error) {
	formatStr, err := f.config.GetFormatString(FormatTypeConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get config format string: %w", err)
	}

	if formatStr == "" {
		formatStr = "%s=%s (source: %s)\n"
	}

	return fmt.Sprintf(formatStr, name, value, source), nil
}

// [CRITICAL] FMT-001: Error formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatError(err error, errorType ErrorType) (string, error) {
	formatStr, err2 := f.config.GetErrorFormat(errorType)
	if err2 != nil {
		return "", fmt.Errorf("failed to get error format for type %s: %w", errorType, err2)
	}

	if formatStr == "" {
		// Provide default error format strings
		switch errorType {
		case ErrorTypeDiskFull:
			formatStr = "Error: Disk full - %s\n"
		case ErrorTypePermission:
			formatStr = "Error: Permission denied - %s\n"
		case ErrorTypeDirectoryNotFound:
			formatStr = "Error: Directory not found - %s\n"
		case ErrorTypeFileNotFound:
			formatStr = "Error: File not found - %s\n"
		case ErrorTypeInvalidDirectory:
			formatStr = "Error: Invalid directory - %s\n"
		case ErrorTypeInvalidFile:
			formatStr = "Error: Invalid file - %s\n"
		case ErrorTypeGeneric:
			formatStr = "Error: %s\n"
		default:
			return "", fmt.Errorf("unsupported error type: %s", errorType)
		}
	}

	return fmt.Sprintf(formatStr, err.Error()), nil
}

// [CRITICAL] FMT-001: Template-based formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatWithTemplate(templateStr string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("format").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// [CRITICAL] FMT-001: Placeholder-based formatting operations - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatWithPlaceholders(format string, data map[string]string) (string, error) {
	result := format

	for key, value := range data {
		placeholder := "%{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result, nil
}

// [CRITICAL] FMT-001: AI-friendly context formatting - [ACTION:core-functionality]
func (f *AICoreFormatter) FormatWithContext(ctx FormatContext) (string, error) {
	switch ctx.FormatType {
	case FormatTypeCreated, FormatTypeIdentical, FormatTypeList, FormatTypeDryRun:
		if path, ok := ctx.Data["path"].(string); ok {
			return f.FormatArchive(path, ctx.FormatType)
		}
		return "", fmt.Errorf("missing path in format context")
	case FormatTypeError:
		if err, ok := ctx.Data["error"].(error); ok {
			if errorType, ok := ctx.Data["errorType"].(ErrorType); ok {
				return f.FormatError(err, errorType)
			}
		}
		return "", fmt.Errorf("missing error or errorType in format context")
	case FormatTypeConfig:
		if name, ok := ctx.Data["name"].(string); ok {
			if value, ok := ctx.Data["value"].(string); ok {
				if source, ok := ctx.Data["source"].(string); ok {
					return f.FormatConfig(name, value, source)
				}
			}
		}
		return "", fmt.Errorf("missing name, value, or source in format context")
	default:
		return "", fmt.Errorf("unsupported format type: %s", ctx.FormatType)
	}
}
