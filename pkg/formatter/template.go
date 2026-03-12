// Template-based formatting functionality for the formatter package.
// Provides pattern-based and placeholder-based template formatting with
// support for both #{name} style and Go text/template {{.name}} placeholders.
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// DefaultTemplateFormatter provides template-based formatting functionality
type DefaultTemplateFormatter struct {
	configProvider ConfigProvider
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// NewDefaultTemplateFormatter creates a new DefaultTemplateFormatter
func NewDefaultTemplateFormatter(configProvider ConfigProvider) *DefaultTemplateFormatter {
	return &DefaultTemplateFormatter{
		configProvider: configProvider,
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// FormatWithPlaceholders formats a string using placeholder-based template formatting
func (tf *DefaultTemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
	result := format

	// Handle #{name} style placeholders - replace ALL placeholders from data map first
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Replace known placeholders that might be missing with defaults
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

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Handle printf-style placeholders (%s, %d, etc.) after template placeholders are replaced
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateCreatedArchive formats a created archive message using a template
func (tf *DefaultTemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("created_archive")
	if templateStr == "" {
		templateStr = "Created archive: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateIdenticalArchive formats an identical archive message using a template
func (tf *DefaultTemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("identical_archive")
	if templateStr == "" {
		templateStr = "Identical archive: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateListArchive formats a list archive message using a template
func (tf *DefaultTemplateFormatter) TemplateListArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("list_archive")
	if templateStr == "" {
		templateStr = "#{path} (created: #{creation_time})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateConfigValue formats a configuration value message using a template
func (tf *DefaultTemplateFormatter) TemplateConfigValue(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("config_value")
	if templateStr == "" {
		templateStr = "#{name}=#{value} (source: #{source})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateDryRunArchive formats a dry-run archive message using a template
func (tf *DefaultTemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("dry_run_archive")
	if templateStr == "" {
		templateStr = "Would create archive: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateError formats an error message using a template
func (tf *DefaultTemplateFormatter) TemplateError(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("error")
	if templateStr == "" {
		templateStr = "Error: #{message}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateCreatedBackup formats a created backup message using a template
func (tf *DefaultTemplateFormatter) TemplateCreatedBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("created_backup")
	if templateStr == "" {
		templateStr = "Created backup: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateIdenticalBackup formats an identical backup message using a template
func (tf *DefaultTemplateFormatter) TemplateIdenticalBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("identical_backup")
	if templateStr == "" {
		templateStr = "Identical backup: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateListBackup formats a list backup message using a template
func (tf *DefaultTemplateFormatter) TemplateListBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("list_backup")
	if templateStr == "" {
		templateStr = "#{path} (created: #{creation_time})"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateDryRunBackup formats a dry-run backup message using a template
func (tf *DefaultTemplateFormatter) TemplateDryRunBackup(data map[string]string) string {
	templateStr := tf.configProvider.GetTemplateString("dry_run_backup")
	if templateStr == "" {
		templateStr = "Would create backup: #{path}"
	}
	return tf.FormatWithPlaceholders(templateStr, data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// SimpleTemplateFormatter provides template formatting without configuration dependency
type SimpleTemplateFormatter struct{}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// NewSimpleTemplateFormatter creates a SimpleTemplateFormatter
func NewSimpleTemplateFormatter() *SimpleTemplateFormatter {
	return &SimpleTemplateFormatter{}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// FormatWithPlaceholders formats a string using placeholder-based template formatting
func (stf *SimpleTemplateFormatter) FormatWithPlaceholders(format string, data map[string]string) string {
	result := format

	// Handle #{name} style placeholders - replace ALL placeholders from data map first
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Replace known placeholders that might be missing with defaults
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

	// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Handle printf-style placeholders (%s, %d, etc.) after template placeholders are replaced
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateCreatedArchive formats a created archive message using default template
func (stf *SimpleTemplateFormatter) TemplateCreatedArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Created archive: #{path}", data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateIdenticalArchive formats an identical archive message using default template
func (stf *SimpleTemplateFormatter) TemplateIdenticalArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Identical archive: #{path}", data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateListArchive formats a list archive message using default template
func (stf *SimpleTemplateFormatter) TemplateListArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("#{path} (created: #{creation_time})", data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateConfigValue formats a configuration value message using default template
func (stf *SimpleTemplateFormatter) TemplateConfigValue(data map[string]string) string {
	return stf.FormatWithPlaceholders("#{name}=#{value} (source: #{source})", data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateDryRunArchive formats a dry-run archive message using default template
func (stf *SimpleTemplateFormatter) TemplateDryRunArchive(data map[string]string) string {
	return stf.FormatWithPlaceholders("Would create archive: #{path}", data)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TemplateError formats an error message using default template
func (stf *SimpleTemplateFormatter) TemplateError(data map[string]string) string {
	return stf.FormatWithPlaceholders("Error: #{message}", data)
}

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
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

// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]

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
