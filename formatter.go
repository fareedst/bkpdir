package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

type OutputFormatter struct {
	cfg *Config
}

func NewOutputFormatter(cfg *Config) *OutputFormatter {
	return &OutputFormatter{cfg: cfg}
}

// Printf-style formatting methods
func (f *OutputFormatter) FormatCreatedArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatCreatedArchive, path)
}

func (f *OutputFormatter) FormatIdenticalArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatIdenticalArchive, path)
}

func (f *OutputFormatter) FormatListArchive(path, creationTime string) string {
	return fmt.Sprintf(f.cfg.FormatListArchive, path, creationTime)
}

func (f *OutputFormatter) FormatConfigValue(name, value, source string) string {
	return fmt.Sprintf(f.cfg.FormatConfigValue, name, value, source)
}

func (f *OutputFormatter) FormatDryRunArchive(path string) string {
	return fmt.Sprintf(f.cfg.FormatDryRunArchive, path)
}

func (f *OutputFormatter) FormatError(message string) string {
	return fmt.Sprintf(f.cfg.FormatError, message)
}

// Template-based formatting methods
func (f *OutputFormatter) FormatCreatedArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateCreatedArchive, data)
}

func (f *OutputFormatter) FormatIdenticalArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateIdenticalArchive, data)
}

func (f *OutputFormatter) FormatListArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateListArchive, data)
}

func (f *OutputFormatter) FormatConfigValueTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateConfigValue, data)
}

func (f *OutputFormatter) FormatDryRunArchiveTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateDryRunArchive, data)
}

func (f *OutputFormatter) FormatErrorTemplate(data map[string]string) string {
	return f.formatTemplate(f.cfg.TemplateError, data)
}

// Print methods that output directly to stdout/stderr
func (f *OutputFormatter) PrintCreatedArchive(path string) {
	fmt.Print(f.FormatCreatedArchive(path))
}

func (f *OutputFormatter) PrintIdenticalArchive(path string) {
	fmt.Print(f.FormatIdenticalArchive(path))
}

func (f *OutputFormatter) PrintListArchive(path, creationTime string) {
	fmt.Print(f.FormatListArchive(path, creationTime))
}

func (f *OutputFormatter) PrintConfigValue(name, value, source string) {
	fmt.Print(f.FormatConfigValue(name, value, source))
}

func (f *OutputFormatter) PrintDryRunArchive(path string) {
	fmt.Print(f.FormatDryRunArchive(path))
}

func (f *OutputFormatter) PrintError(message string) {
	fmt.Fprint(os.Stderr, f.FormatError(message))
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
func (f *OutputFormatter) ExtractArchiveFilenameData(filename string) map[string]string {
	return f.extractPatternData(f.cfg.PatternArchiveFilename, filename)
}

func (f *OutputFormatter) ExtractConfigLineData(line string) map[string]string {
	return f.extractPatternData(f.cfg.PatternConfigLine, line)
}

func (f *OutputFormatter) ExtractTimestampData(timestamp string) map[string]string {
	return f.extractPatternData(f.cfg.PatternTimestamp, timestamp)
}

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

// Enhanced formatting methods that use regex extraction
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
