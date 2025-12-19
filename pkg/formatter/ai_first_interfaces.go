// AI-First Formatter Interfaces
// Provides clear component separation and interface standardization for optimal AI assistant comprehension.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:core-functionality]
// [REQ:OUTPUT_FORMATTING] Output formatting requirement
// [ARCH:OUTPUT_FORMATTING] Output formatting architecture
// [IMPL:DUAL_FORMATTING] Dual-mode formatting implementation
package formatter

import (
	"fmt"
	"time"
)

// [CRITICAL] FMT-001: Format type enumeration - [ACTION:core-functionality]
type FormatType string

const (
	FormatTypeCreated   FormatType = "created"
	FormatTypeIdentical FormatType = "identical"
	FormatTypeList      FormatType = "list"
	FormatTypeDryRun    FormatType = "dry_run"
	FormatTypeError     FormatType = "error"
	FormatTypeConfig    FormatType = "config"
)

// [CRITICAL] FMT-001: Error type enumeration - [ACTION:core-functionality]
type ErrorType string

const (
	ErrorTypeDiskFull          ErrorType = "disk_full"
	ErrorTypePermission        ErrorType = "permission"
	ErrorTypeDirectoryNotFound ErrorType = "directory_not_found"
	ErrorTypeFileNotFound      ErrorType = "file_not_found"
	ErrorTypeInvalidDirectory  ErrorType = "invalid_directory"
	ErrorTypeInvalidFile       ErrorType = "invalid_file"
	ErrorTypeGeneric           ErrorType = "generic"
)

// [CRITICAL] FMT-001: Pattern type enumeration - [ACTION:discovery]
type PatternType string

const (
	PatternTypeArchiveFilename PatternType = "archive_filename"
	PatternTypeBackupFilename  PatternType = "backup_filename"
	PatternTypeConfigLine      PatternType = "config_line"
	PatternTypeTimestamp       PatternType = "timestamp"
)

// [CRITICAL] FMT-001: AI-first output destination enumeration - [ACTION:format-processing]
type AIOutputDestination string

const (
	AIOutputDestinationStdout AIOutputDestination = "stdout"
	AIOutputDestinationStderr AIOutputDestination = "stderr"
)

// [CRITICAL] FMT-001: AI-first message type enumeration - [ACTION:core-functionality]
type AIMessageType string

const (
	AIMessageTypeInfo    AIMessageType = "info"
	AIMessageTypeError   AIMessageType = "error"
	AIMessageTypeWarning AIMessageType = "warning"
	AIMessageTypeConfig  AIMessageType = "config"
)

// [CRITICAL] FMT-001: AI-first structured data types - [ACTION:core-functionality]
type AIArchiveData struct {
	Prefix string `json:"prefix"`
	Year   string `json:"year"`
	Month  string `json:"month"`
	Day    string `json:"day"`
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
	Branch string `json:"branch"`
	Hash   string `json:"hash"`
	Note   string `json:"note"`
}

type AIBackupData struct {
	Filename string `json:"filename"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	Day      string `json:"day"`
	Hour     string `json:"hour"`
	Minute   string `json:"minute"`
	Note     string `json:"note"`
}

type AIConfigData struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

type AITimestampData struct {
	Year   string `json:"year"`
	Month  string `json:"month"`
	Day    string `json:"day"`
	Hour   string `json:"hour"`
	Minute string `json:"minute"`
	Second string `json:"second"`
}

// [CRITICAL] FMT-001: AI-first output message structure - [ACTION:core-functionality]
type AIOutputMessage struct {
	Content     string              `json:"content"`
	Destination AIOutputDestination `json:"destination"`
	Type        AIMessageType       `json:"type"`
	Metadata    map[string]string   `json:"metadata,omitempty"`
	Timestamp   time.Time           `json:"timestamp"`
}

// [CRITICAL] FMT-001: Context structures for AI comprehension - [ACTION:core-functionality]
type FormatContext struct {
	FormatType FormatType             `json:"format_type"`
	Data       map[string]interface{} `json:"data"`
	Options    FormatOptions          `json:"options"`
	Metadata   map[string]string      `json:"metadata"`
}

type ExtractContext struct {
	PatternType PatternType       `json:"pattern_type"`
	Input       string            `json:"input"`
	Options     ExtractOptions    `json:"options"`
	Metadata    map[string]string `json:"metadata"`
}

type PrintContext struct {
	Message     string              `json:"message"`
	Destination AIOutputDestination `json:"destination"`
	Type        AIMessageType       `json:"type"`
	Options     PrintOptions        `json:"options"`
	Metadata    map[string]string   `json:"metadata"`
}

// [CRITICAL] FMT-001: Options structures for AI comprehension - [ACTION:core-functionality]
type FormatOptions struct {
	UseTemplate  bool              `json:"use_template"`
	Template     string            `json:"template,omitempty"`
	Placeholders map[string]string `json:"placeholders,omitempty"`
}

type ExtractOptions struct {
	ValidatePattern  bool `json:"validate_pattern"`
	ReturnStructured bool `json:"return_structured"`
}

type PrintOptions struct {
	Delayed bool `json:"delayed"`
	Flush   bool `json:"flush"`
}

// [CRITICAL] FMT-001: Core formatting engine - [ACTION:core-functionality]
type CoreFormatter interface {
	// Pure formatting operations (no side effects)
	FormatArchive(path string, formatType FormatType) (string, error)
	FormatBackup(path string, formatType FormatType) (string, error)
	FormatConfig(name, value, source string) (string, error)
	FormatError(err error, errorType ErrorType) (string, error)

	// Template-based formatting
	FormatWithTemplate(template string, data map[string]interface{}) (string, error)
	FormatWithPlaceholders(format string, data map[string]string) (string, error)

	// AI-friendly operations
	FormatWithContext(ctx FormatContext) (string, error)
}

// [CRITICAL] FMT-001: AI-first pattern extraction engine - [ACTION:discovery]
type AIPatternExtractor interface {
	// Core extraction operations
	ExtractArchiveData(filename string) (AIArchiveData, error)
	ExtractBackupData(filename string) (AIBackupData, error)
	ExtractConfigData(line string) (AIConfigData, error)
	ExtractTimestampData(timestamp string) (AITimestampData, error)

	// Generic pattern extraction
	ExtractPattern(pattern, text string) (map[string]string, error)

	// AI-friendly operations
	ExtractWithContext(ctx ExtractContext) (interface{}, error)
}

// [CRITICAL] FMT-001: AI-first output management engine - [ACTION:format-processing]
type AIOutputManager interface {
	// Direct output operations
	Print(message string) error
	PrintError(message string) error

	// Delayed output operations
	Collect(message AIOutputMessage) error
	Flush() error
	FlushStdout() error
	FlushStderr() error
	Clear() error

	// Output state management
	IsDelayedMode() bool
	SetDelayedMode(enabled bool) error
	GetCollectedMessages() []AIOutputMessage

	// AI-friendly operations
	PrintWithContext(ctx PrintContext) error
}

// [CRITICAL] FMT-001: Configuration abstraction layer - [ACTION:configure-modify]
type FormatterConfig interface {
	// Format string access
	GetFormatString(formatType FormatType) (string, error)
	GetTemplateString(templateType FormatType) (string, error)
	GetErrorFormat(errorType ErrorType) (string, error)

	// Pattern access
	GetPattern(patternType PatternType) (string, error)

	// Configuration validation
	Validate() error
	GetValidationErrors() []ConfigError
}

// [CRITICAL] FMT-001: Configuration error structure - [ACTION:core-functionality]
type ConfigError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("config error [%s]: %s - %s", ce.Code, ce.Field, ce.Message)
}

// [CRITICAL] FMT-001: Primary formatter interface - [ACTION:core-functionality]
type AIFirstFormatter interface {
	// Core formatting capabilities
	CoreFormatter
	AIPatternExtractor
	AIOutputManager

	// Configuration management
	SetConfig(config FormatterConfig) error
	GetConfig() FormatterConfig

	// AI-friendly operations
	FormatWithContext(ctx FormatContext) (string, error)
	ExtractWithContext(ctx ExtractContext) (interface{}, error)
	PrintWithContext(ctx PrintContext) error
}
