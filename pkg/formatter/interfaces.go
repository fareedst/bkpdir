// Interfaces for the formatter package providing clean abstractions for output formatting.
// These interfaces enable dependency injection and testing while maintaining
// backward compatibility with the original formatter functionality.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

// ⭐ EXTRACT-003: Core interfaces - 🔧 Configuration provider abstraction
// ConfigProvider abstracts configuration access for formatter components
type ConfigProvider interface {
	GetFormatString(formatType string) string
	GetTemplateString(templateType string) string
	GetPattern(patternType string) string
	GetErrorFormat(errorType string) string
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Output destination abstraction
// OutputDestination abstracts output handling for formatter components
type OutputDestination interface {
	Print(message string)
	PrintError(message string)
	IsDelayedMode() bool
	SetCollector(collector *OutputCollector)
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Pattern extraction interface
// PatternExtractor defines contract for regex-based data extraction
type PatternExtractor interface {
	ExtractArchiveFilenameData(filename string) map[string]string
	ExtractBackupFilenameData(filename string) map[string]string
	ExtractPatternData(pattern, text string) map[string]string
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Printf formatter interface
// Formatter provides printf-style formatting operations
type Formatter interface {
	FormatCreatedArchive(path string) string
	FormatIdenticalArchive(path string) string
	FormatListArchive(path, creationTime string) string
	FormatConfigValue(name, value, source string) string
	FormatError(message string) string
	// Extended formatting operations
	FormatDryRunArchive(path string) string
	FormatCreatedBackup(path string) string
	FormatIdenticalBackup(path string) string
	FormatListBackup(path, creationTime string) string
	FormatDryRunBackup(path string) string
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Template formatter interface
// TemplateFormatter provides template-based formatting operations
type TemplateFormatter interface {
	FormatWithTemplate(input, pattern, tmplStr string) (string, error)
	FormatWithPlaceholders(format string, data map[string]string) string
	TemplateCreatedArchive(data map[string]string) string
	TemplateIdenticalArchive(data map[string]string) string
	// Extended template operations
	TemplateListArchive(data map[string]string) string
	TemplateConfigValue(data map[string]string) string
	TemplateDryRunArchive(data map[string]string) string
	TemplateError(data map[string]string) string
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Error formatter interface
// ErrorFormatter provides specialized error formatting
type ErrorFormatter interface {
	FormatDiskFullError(err error) string
	FormatPermissionError(err error) string
	FormatDirectoryNotFound(err error) string
	FormatFileNotFound(err error) string
	FormatInvalidDirectory(err error) string
	FormatInvalidFile(err error) string
	// Template-based error formatting
	TemplateDiskFullError(err error) string
	TemplatePermissionError(err error) string
	TemplateDirectoryNotFound(err error) string
	TemplateFileNotFound(err error) string
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Print operations interface
// PrintFormatter provides print operations with output destination support
type PrintFormatter interface {
	PrintCreatedArchive(path string)
	PrintIdenticalArchive(path string)
	PrintListArchive(path, creationTime string)
	PrintConfigValue(name, value, source string)
	PrintError(message string)
	// Extended print operations
	PrintDryRunArchive(path string)
	PrintCreatedBackup(path string)
	PrintIdenticalBackup(path string)
	PrintListBackup(path, creationTime string)
	PrintDryRunBackup(path string)
}

// ⭐ EXTRACT-003: Core interfaces - 🔧 Comprehensive formatter interface
// OutputFormatterInterface combines all formatting capabilities
type OutputFormatterInterface interface {
	Formatter
	TemplateFormatter
	ErrorFormatter
	PrintFormatter
	PatternExtractor

	// Delayed output support
	IsDelayedMode() bool
	GetCollector() *OutputCollector
	SetCollector(collector *OutputCollector)
}
