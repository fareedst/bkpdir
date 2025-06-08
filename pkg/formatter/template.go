// Template-based formatting functionality for the formatter package.
// Provides pattern-based and placeholder-based template formatting with
// support for both %{name} style and Go text/template {{.name}} placeholders.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Template formatting implementation
// DefaultTemplateFormatter provides template-based formatting functionality
type DefaultTemplateFormatter struct {
	configProvider ConfigProvider
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Constructor
// NewDefaultTemplateFormatter creates a new DefaultTemplateFormatter
func NewDefaultTemplateFormatter(configProvider ConfigProvider) *DefaultTemplateFormatter {
	return &DefaultTemplateFormatter{
		configProvider: configProvider,
	}
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Pattern and template processing
// FormatWithTemplate formats input using a pattern and template string
func (tf *DefaultTemplateFormatter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
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

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Placeholder-based formatting
// FormatWithPlaceholders formats a string using placeholder-based template formatting
func (tf *DefaultTemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
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

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Archive template operations
// TemplateCreatedArchive formats a created archive message using a template
func (tf *DefaultTemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("created_archive")
	if templateStr == "" {
		templateStr = "Created archive: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Identical archive template
// TemplateIdenticalArchive formats an identical archive message using a template
func (tf *DefaultTemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("identical_archive")
	if templateStr == "" {
		templateStr = "Identical archive: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 List archive template
// TemplateListArchive formats a list archive message using a template
func (tf *DefaultTemplateFormatter) TemplateListArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("list_archive")
	if templateStr == "" {
		templateStr = "%{path} (created: %{creation_time})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Configuration value template
// TemplateConfigValue formats a configuration value message using a template
func (tf *DefaultTemplateFormatter) TemplateConfigValue(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("config_value")
	if templateStr == "" {
		templateStr = "%{name}=%{value} (source: %{source})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Dry run archive template
// TemplateDryRunArchive formats a dry-run archive message using a template
func (tf *DefaultTemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("dry_run_archive")
	if templateStr == "" {
		templateStr = "Would create archive: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Error template
// TemplateError formats an error message using a template
func (tf *DefaultTemplateFormatter) TemplateError(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("error")
	if templateStr == "" {
		templateStr = "Error: %{message}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Backup template operations
// TemplateCreatedBackup formats a created backup message using a template
func (tf *DefaultTemplateFormatter) TemplateCreatedBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("created_backup")
	if templateStr == "" {
		templateStr = "Created backup: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Identical backup template
// TemplateIdenticalBackup formats an identical backup message using a template
func (tf *DefaultTemplateFormatter) TemplateIdenticalBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("identical_backup")
	if templateStr == "" {
		templateStr = "Identical backup: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 List backup template
// TemplateListBackup formats a list backup message using a template
func (tf *DefaultTemplateFormatter) TemplateListBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("list_backup")
	if templateStr == "" {
		templateStr = "%{path} (created: %{creation_time})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Dry run backup template
// TemplateDryRunBackup formats a dry-run backup message using a template
func (tf *DefaultTemplateFormatter) TemplateDryRunBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("dry_run_backup")
	if templateStr == "" {
		templateStr = "Would create backup: %{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple template formatter without config
// SimpleTemplateFormatter provides template formatting without configuration dependency
type SimpleTemplateFormatter struct{}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple constructor
// NewSimpleTemplateFormatter creates a SimpleTemplateFormatter
func NewSimpleTemplateFormatter() *SimpleTemplateFormatter {
	return &SimpleTemplateFormatter{}
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple pattern and template processing
// FormatWithTemplate formats input using a pattern and template string
func (stf *SimpleTemplateFormatter) FormatWithTemplate(input, pattern, tmplStr string) (string, error) {
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
	return stf.FormatWithPlaceholders(tmplStr, data), nil
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple placeholder formatting
// FormatWithPlaceholders formats a string using placeholder-based template formatting
func (stf *SimpleTemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
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

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple template operations with defaults
// TemplateCreatedArchive formats a created archive message using default template
func (stf *SimpleTemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Created archive: %{path}", data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple identical archive template
// TemplateIdenticalArchive formats an identical archive message using default template
func (stf *SimpleTemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Identical archive: %{path}", data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple list archive template
// TemplateListArchive formats a list archive message using default template
func (stf *SimpleTemplateFormatter) TemplateListArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("%{path} (created: %{creation_time})", data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple config value template
// TemplateConfigValue formats a configuration value message using default template
func (stf *SimpleTemplateFormatter) TemplateConfigValue(data map[string]string) string {
	return stf.FormatWithPlaceholders("%{name}=%{value} (source: %{source})", data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple dry run template
// TemplateDryRunArchive formats a dry-run archive message using default template
func (stf *SimpleTemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Would create archive: %{path}", data)
}

// ⭐ EXTRACT-003: TemplateFormatter component - 🔧 Simple error template
// TemplateError formats an error message using default template
func (stf *SimpleTemplateFormatter) TemplateError(data map[string]string) string {
	return stf.FormatWithPlaceholders("Error: %{message}", data)
}

// ⭐ OUT-002: Enhanced output with file statistics - Enhanced template methods
// Enhanced template operations using file statistics for detailed output

// TemplateCreatedArchiveWithStats formats a created archive message using templates with file statistics (DefaultTemplateFormatter)
func (tf *DefaultTemplateFormatter) TemplateCreatedArchiveWithStats(path string) string {
	stats, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat fails
		data := map[string]string{"path": path, "name": filepath.Base(path)}
		return tf.TemplateCreatedArchive(data)
	}

	templateStr := tf.configProvider.GetDetailedTemplateString("created_archive")
	if templateStr == "" {
		templateStr = "Created archive: {path} ({size_human}, {mtime})"
	}

	data := buildStatsTemplateData(stats)
	return formatTemplate(templateStr, data)
}

// TemplateIncrementalCreatedWithStats formats an incremental created message using templates with file statistics (DefaultTemplateFormatter)
func (tf *DefaultTemplateFormatter) TemplateIncrementalCreatedWithStats(path string) string {
	stats, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat fails
		data := map[string]string{"path": path, "name": filepath.Base(path)}
		return tf.TemplateCreatedArchive(data) // Use same as full archive for now
	}

	templateStr := tf.configProvider.GetDetailedTemplateString("incremental_created")
	if templateStr == "" {
		templateStr = "Created incremental archive: {path} ({size_human}, {mtime})"
	}

	data := buildStatsTemplateData(stats)
	return formatTemplate(templateStr, data)
}

// TemplateCreatedArchiveWithStats formats a created archive message using templates with file statistics (SimpleTemplateFormatter)
func (stf *SimpleTemplateFormatter) TemplateCreatedArchiveWithStats(path string) string {
	stats, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat fails
		data := map[string]string{"path": path, "name": filepath.Base(path)}
		return stf.TemplateCreatedArchive(data)
	}

	data := buildStatsTemplateData(stats)
	return formatTemplate("Created archive: {path} ({size_human}, {mtime})", data)
}

// TemplateIncrementalCreatedWithStats formats an incremental created message using templates with file statistics (SimpleTemplateFormatter)
func (stf *SimpleTemplateFormatter) TemplateIncrementalCreatedWithStats(path string) string {
	stats, err := GatherFileStatInfo(path)
	if err != nil {
		// Fallback to basic template if stat fails
		data := map[string]string{"path": path, "name": filepath.Base(path)}
		return stf.TemplateCreatedArchive(data) // Use same as full archive for now
	}

	data := buildStatsTemplateData(stats)
	return formatTemplate("Created incremental archive: {path} ({size_human}, {mtime})", data)
}

// ⭐ OUT-002: Enhanced output with file statistics - Helper functions for template formatting

// buildStatsTemplateData builds template data from file statistics
func buildStatsTemplateData(stats *FileStatInfo) map[string]string {
	return map[string]string{
		"path":       stats.Path,
		"name":       stats.Name,
		"size":       fmt.Sprintf("%d", stats.Size),
		"size_human": stats.SizeHuman,
		"mtime":      stats.MTime.Format("2006-01-02 15:04:05"),
		"mtime_unix": fmt.Sprintf("%d", stats.MTimeUnix),
		"mode":       stats.Mode.String(),
		"type":       stats.Type,
	}
}

// formatTemplate processes template strings with {placeholder} format
func formatTemplate(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}
