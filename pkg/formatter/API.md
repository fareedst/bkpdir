# AI-First Formatter API Documentation

This document provides comprehensive API documentation for the AI-first formatter interfaces, designed for optimal AI assistant comprehension and usage.

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Data Types](#data-types)
3. [Context Structures](#context-structures)
4. [Implementation Details](#implementation-details)
5. [Usage Patterns](#usage-patterns)
6. [Error Handling](#error-handling)
7. [Testing Guidelines](#testing-guidelines)

## Core Interfaces

### AIFirstFormatter

**Purpose**: Primary interface that combines all formatting capabilities with AI-friendly operations.

**Interface Definition**:
```go
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
```

**Key Methods**:

#### `SetConfig(config FormatterConfig) error`
- **Purpose**: Updates the formatter's configuration
- **Parameters**: `config` - New configuration to apply
- **Returns**: Error if configuration is invalid or nil
- **Usage**: Call when configuration needs to be updated dynamically

#### `GetConfig() FormatterConfig`
- **Purpose**: Retrieves the current configuration
- **Returns**: Current configuration instance
- **Usage**: Access current configuration for inspection or validation

#### `FormatWithContext(ctx FormatContext) (string, error)`
- **Purpose**: Formats output using AI-friendly context structure
- **Parameters**: `ctx` - Format context with type, data, and options
- **Returns**: Formatted string and error
- **Usage**: Primary formatting method for AI assistants

#### `ExtractWithContext(ctx ExtractContext) (interface{}, error)`
- **Purpose**: Extracts structured data using AI-friendly context
- **Parameters**: `ctx` - Extract context with pattern type and input
- **Returns**: Structured data (interface{}) and error
- **Usage**: Extract data from text using configured patterns

#### `PrintWithContext(ctx PrintContext) error`
- **Purpose**: Prints output using AI-friendly context structure
- **Parameters**: `ctx` - Print context with message and options
- **Returns**: Error if printing fails
- **Usage**: Output messages with rich context and options

### CoreFormatter

**Purpose**: Provides pure formatting operations without side effects.

**Interface Definition**:
```go
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
```

**Key Methods**:

#### `FormatArchive(path string, formatType FormatType) (string, error)`
- **Purpose**: Formats archive-related messages
- **Parameters**: 
  - `path` - Archive file path
  - `formatType` - Type of format to apply
- **Returns**: Formatted string and error
- **Format Types**: `FormatTypeCreated`, `FormatTypeIdentical`, `FormatTypeList`, `FormatTypeDryRun`

#### `FormatBackup(path string, formatType FormatType) (string, error)`
- **Purpose**: Formats backup-related messages
- **Parameters**: 
  - `path` - Backup file path
  - `formatType` - Type of format to apply
- **Returns**: Formatted string and error
- **Format Types**: `FormatTypeCreated`, `FormatTypeIdentical`, `FormatTypeList`, `FormatTypeDryRun`

#### `FormatConfig(name, value, source string) (string, error)`
- **Purpose**: Formats configuration value display
- **Parameters**: 
  - `name` - Configuration key name
  - `value` - Configuration value
  - `source` - Configuration source
- **Returns**: Formatted string and error

#### `FormatError(err error, errorType ErrorType) (string, error)`
- **Purpose**: Formats error messages with specific error types
- **Parameters**: 
  - `err` - Error to format
  - `errorType` - Type of error for specialized formatting
- **Returns**: Formatted error string and error
- **Error Types**: `ErrorTypeDiskFull`, `ErrorTypePermission`, `ErrorTypeGeneric`, etc.

#### `FormatWithTemplate(template string, data map[string]interface{}) (string, error)`
- **Purpose**: Formats using Go text/template
- **Parameters**: 
  - `template` - Go template string
  - `data` - Data to inject into template
- **Returns**: Formatted string and error

#### `FormatWithPlaceholders(format string, data map[string]string) (string, error)`
- **Purpose**: Formats using placeholder substitution
- **Parameters**: 
  - `format` - Format string with placeholders like `%{key}`
  - `data` - Data to substitute
- **Returns**: Formatted string and error

### AIPatternExtractor

**Purpose**: Provides structured data extraction from text patterns.

**Interface Definition**:
```go
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
```

**Key Methods**:

#### `ExtractArchiveData(filename string) (AIArchiveData, error)`
- **Purpose**: Extracts structured data from archive filenames
- **Parameters**: `filename` - Archive filename to parse
- **Returns**: `AIArchiveData` structure and error
- **Data Fields**: `Prefix`, `Year`, `Month`, `Day`, `Hour`, `Minute`, `Branch`, `Hash`, `Note`

#### `ExtractBackupData(filename string) (AIBackupData, error)`
- **Purpose**: Extracts structured data from backup filenames
- **Parameters**: `filename` - Backup filename to parse
- **Returns**: `AIBackupData` structure and error
- **Data Fields**: `Filename`, `Year`, `Month`, `Day`, `Hour`, `Minute`, `Note`

#### `ExtractConfigData(line string) (AIConfigData, error)`
- **Purpose**: Extracts structured data from configuration lines
- **Parameters**: `line` - Configuration line to parse
- **Returns**: `AIConfigData` structure and error
- **Data Fields**: `Name`, `Value`, `Source`

#### `ExtractTimestampData(timestamp string) (AITimestampData, error)`
- **Purpose**: Extracts structured data from timestamp strings
- **Parameters**: `timestamp` - Timestamp string to parse
- **Returns**: `AITimestampData` structure and error
- **Data Fields**: `Year`, `Month`, `Day`, `Hour`, `Minute`, `Second`

#### `ExtractPattern(pattern, text string) (map[string]string, error)`
- **Purpose**: Generic pattern extraction using regex
- **Parameters**: 
  - `pattern` - Regex pattern with named groups
  - `text` - Text to extract from
- **Returns**: Map of named groups and error

### AIOutputManager

**Purpose**: Provides output handling with delayed output support.

**Interface Definition**:
```go
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
```

**Key Methods**:

#### `Print(message string) error`
- **Purpose**: Prints message to stdout
- **Parameters**: `message` - Message to print
- **Returns**: Error if printing fails
- **Behavior**: Prints immediately or collects based on mode

#### `PrintError(message string) error`
- **Purpose**: Prints error message to stderr
- **Parameters**: `message` - Error message to print
- **Returns**: Error if printing fails
- **Behavior**: Prints immediately or collects based on mode

#### `Collect(message AIOutputMessage) error`
- **Purpose**: Collects message for delayed output
- **Parameters**: `message` - Message with metadata
- **Returns**: Error if collection fails
- **Usage**: Use in delayed mode to batch output

#### `Flush() error`
- **Purpose**: Flushes all collected messages
- **Returns**: Error if flush fails
- **Usage**: Display all collected messages at once

#### `IsDelayedMode() bool`
- **Purpose**: Checks if output is being collected
- **Returns**: True if in delayed mode
- **Usage**: Determine current output behavior

#### `SetDelayedMode(enabled bool) error`
- **Purpose**: Enables or disables delayed output mode
- **Parameters**: `enabled` - Whether to enable delayed mode
- **Returns**: Error if mode change fails
- **Usage**: Switch between immediate and delayed output

## Data Types

### FormatType

**Purpose**: Strongly typed enumeration for format types.

```go
type FormatType string

const (
    FormatTypeCreated    FormatType = "created"
    FormatTypeIdentical  FormatType = "identical"
    FormatTypeList       FormatType = "list"
    FormatTypeDryRun     FormatType = "dry_run"
    FormatTypeError      FormatType = "error"
    FormatTypeConfig     FormatType = "config"
)
```

**Usage**: Use these constants when calling formatting methods to ensure type safety.

### ErrorType

**Purpose**: Strongly typed enumeration for error types.

```go
type ErrorType string

const (
    ErrorTypeDiskFull        ErrorType = "disk_full"
    ErrorTypePermission      ErrorType = "permission"
    ErrorTypeDirectoryNotFound ErrorType = "directory_not_found"
    ErrorTypeFileNotFound    ErrorType = "file_not_found"
    ErrorTypeInvalidDirectory ErrorType = "invalid_directory"
    ErrorTypeInvalidFile     ErrorType = "invalid_file"
    ErrorTypeGeneric         ErrorType = "generic"
)
```

**Usage**: Use these constants when formatting errors for specialized error handling.

### Structured Data Types

#### AIArchiveData

**Purpose**: Structured data extracted from archive filenames.

```go
type AIArchiveData struct {
    Prefix    string `json:"prefix"`
    Year      string `json:"year"`
    Month     string `json:"month"`
    Day       string `json:"day"`
    Hour      string `json:"hour"`
    Minute    string `json:"minute"`
    Branch    string `json:"branch"`
    Hash      string `json:"hash"`
    Note      string `json:"note"`
}
```

**Example**: `"myproject-2024-01-15-14-30-main-abc123-release.tar.gz"` extracts to:
```json
{
    "prefix": "myproject",
    "year": "2024",
    "month": "01",
    "day": "15",
    "hour": "14",
    "minute": "30",
    "branch": "main",
    "hash": "abc123",
    "note": "release"
}
```

#### AIBackupData

**Purpose**: Structured data extracted from backup filenames.

```go
type AIBackupData struct {
    Filename  string `json:"filename"`
    Year      string `json:"year"`
    Month     string `json:"month"`
    Day       string `json:"day"`
    Hour      string `json:"hour"`
    Minute    string `json:"minute"`
    Note      string `json:"note"`
}
```

**Example**: `"config.yaml-2024-01-15-14-30-before-update.bak"` extracts to:
```json
{
    "filename": "config.yaml",
    "year": "2024",
    "month": "01",
    "day": "15",
    "hour": "14",
    "minute": "30",
    "note": "before-update"
}
```

## Context Structures

### FormatContext

**Purpose**: AI-friendly context for formatting operations.

```go
type FormatContext struct {
    FormatType    FormatType                `json:"format_type"`
    Data          map[string]interface{}    `json:"data"`
    Options       FormatOptions             `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

**Fields**:
- `FormatType`: Type of format to apply
- `Data`: Map of data to format (e.g., `{"path": "/file.txt"}`)
- `Options`: Formatting options (template usage, placeholders)
- `Metadata`: Additional context for AI assistants

**Example Usage**:
```go
ctx := FormatContext{
    FormatType: FormatTypeCreated,
    Data: map[string]interface{}{
        "path": "/backups/archive.tar.gz",
        "size": "1024",
    },
    Options: FormatOptions{
        UseTemplate: true,
        Template:    "Created {{.path}} ({{.size}} bytes)",
    },
    Metadata: map[string]string{
        "operation": "backup",
        "user":      "admin",
    },
}
```

### ExtractContext

**Purpose**: AI-friendly context for data extraction operations.

```go
type ExtractContext struct {
    PatternType   PatternType               `json:"pattern_type"`
    Input         string                    `json:"input"`
    Options       ExtractOptions            `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

**Fields**:
- `PatternType`: Type of pattern to use for extraction
- `Input`: Text to extract data from
- `Options`: Extraction options (validation, structured return)
- `Metadata`: Additional context for AI assistants

**Example Usage**:
```go
ctx := ExtractContext{
    PatternType: PatternTypeArchiveFilename,
    Input:       "myproject-2024-01-15-14-30-main-abc123-release.tar.gz",
    Options: ExtractOptions{
        ValidatePattern: true,
        ReturnStructured: true,
    },
    Metadata: map[string]string{
        "source": "user_input",
        "validation": "strict",
    },
}
```

### PrintContext

**Purpose**: AI-friendly context for output operations.

```go
type PrintContext struct {
    Message       string                    `json:"message"`
    Destination   AIOutputDestination       `json:"destination"`
    Type          AIMessageType             `json:"type"`
    Options       PrintOptions              `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

**Fields**:
- `Message`: Message content to print
- `Destination`: Output destination (stdout/stderr)
- `Type`: Message type (info/error/warning/config)
- `Options`: Print options (delayed, flush)
- `Metadata`: Additional context for AI assistants

**Example Usage**:
```go
ctx := PrintContext{
    Message:     "Archive created successfully",
    Destination: AIOutputDestinationStdout,
    Type:        AIMessageTypeInfo,
    Options: PrintOptions{
        Delayed: true,
        Flush:   false,
    },
    Metadata: map[string]string{
        "operation": "backup",
        "priority":  "normal",
    },
}
```

## Implementation Details

### Constructor Functions

#### `NewAIFirstFormatter(config FormatterConfig) *AIFirstFormatterImpl`
- **Purpose**: Creates a new AI-first formatter instance
- **Parameters**: `config` - Configuration provider
- **Returns**: Configured formatter instance
- **Usage**: Primary constructor for AI-first formatter

#### `NewAIFirstFormatterWithCollector(config FormatterConfig, collector *OutputCollector) *AIFirstFormatterImpl`
- **Purpose**: Creates formatter with delayed output support
- **Parameters**: 
  - `config` - Configuration provider
  - `collector` - Output collector for delayed mode
- **Returns**: Configured formatter with collector
- **Usage**: When delayed output is needed

### Component Constructors

#### `NewAICoreFormatter(config FormatterConfig) *AICoreFormatter`
- **Purpose**: Creates core formatter component
- **Parameters**: `config` - Configuration provider
- **Returns**: Core formatter instance
- **Usage**: For pure formatting operations

#### `NewAIPatternExtractor(config FormatterConfig) *AIPatternExtractorImpl`
- **Purpose**: Creates pattern extractor component
- **Parameters**: `config` - Configuration provider
- **Returns**: Pattern extractor instance
- **Usage**: For structured data extraction

#### `NewAIOutputManager() *AIOutputManagerImpl`
- **Purpose**: Creates output manager component
- **Parameters**: None
- **Returns**: Output manager instance
- **Usage**: For output handling

#### `NewAIOutputManagerWithCollector(collector *OutputCollector) *AIOutputManagerImpl`
- **Purpose**: Creates output manager with delayed output support
- **Parameters**: `collector` - Output collector
- **Returns**: Output manager with collector
- **Usage**: When delayed output is needed

## Usage Patterns

### Basic Formatting Pattern

```go
// 1. Create formatter
config := NewMockFormatterConfig()
formatter := NewAIFirstFormatter(config)

// 2. Create format context
ctx := FormatContext{
    FormatType: FormatTypeCreated,
    Data: map[string]interface{}{
        "path": "/backups/archive.tar.gz",
    },
    Options:  FormatOptions{},
    Metadata: make(map[string]string),
}

// 3. Format with context
result, err := formatter.FormatWithContext(ctx)
if err != nil {
    log.Fatal(err)
}

// 4. Print result
printCtx := PrintContext{
    Message:     result,
    Destination: AIOutputDestinationStdout,
    Type:        AIMessageTypeInfo,
    Options:     PrintOptions{},
    Metadata:    make(map[string]string),
}

err = formatter.PrintWithContext(printCtx)
if err != nil {
    log.Fatal(err)
}
```

### Pattern Extraction Pattern

```go
// 1. Create extract context
ctx := ExtractContext{
    PatternType: PatternTypeArchiveFilename,
    Input:       "myproject-2024-01-15-14-30-main-abc123-release.tar.gz",
    Options:     ExtractOptions{},
    Metadata:    make(map[string]string),
}

// 2. Extract data
result, err := formatter.ExtractWithContext(ctx)
if err != nil {
    log.Fatal(err)
}

// 3. Type assert and use structured data
if archiveData, ok := result.(AIArchiveData); ok {
    fmt.Printf("Archive: %s, Branch: %s, Hash: %s\n", 
        archiveData.Prefix, archiveData.Branch, archiveData.Hash)
}
```

### Delayed Output Pattern

```go
// 1. Create formatter with collector
collector := NewOutputCollector()
formatter := NewAIFirstFormatterWithCollector(config, collector)

// 2. Collect multiple messages
for _, item := range items {
    ctx := PrintContext{
        Message:     fmt.Sprintf("Processing %s", item),
        Destination: AIOutputDestinationStdout,
        Type:        AIMessageTypeInfo,
        Options: PrintOptions{
            Delayed: true,
            Flush:   false,
        },
        Metadata: make(map[string]string),
    }
    
    err := formatter.PrintWithContext(ctx)
    if err != nil {
        log.Fatal(err)
    }
}

// 3. Flush all collected messages
err := formatter.Flush()
if err != nil {
    log.Fatal(err)
}
```

### Error Handling Pattern

```go
// 1. Create error context
ctx := FormatContext{
    FormatType: FormatTypeError,
    Data: map[string]interface{}{
        "error":     err,
        "errorType": ErrorTypeDiskFull,
    },
    Options:  FormatOptions{},
    Metadata: make(map[string]string),
}

// 2. Format error
result, err := formatter.FormatWithContext(ctx)
if err != nil {
    log.Fatal(err)
}

// 3. Print error to stderr
printCtx := PrintContext{
    Message:     result,
    Destination: AIOutputDestinationStderr,
    Type:        AIMessageTypeError,
    Options:     PrintOptions{},
    Metadata:    make(map[string]string),
}

err = formatter.PrintWithContext(printCtx)
if err != nil {
    log.Fatal(err)
}
```

## Error Handling

### Error Types

1. **Configuration Errors**: Invalid configuration provided
2. **Format Errors**: Invalid format strings or templates
3. **Pattern Errors**: Invalid regex patterns or no matches
4. **Output Errors**: Output destination unavailable
5. **Context Errors**: Invalid context data or missing required fields

### Error Handling Best Practices

1. **Always Check Errors**: Check all error returns from interface methods
2. **Provide Context**: Include relevant context in error messages
3. **Use Structured Errors**: Prefer `ConfigError` for configuration issues
4. **Graceful Degradation**: Provide fallback behavior when possible
5. **Log Errors**: Log errors with sufficient context for debugging

### Error Recovery Patterns

```go
// Pattern 1: Retry with fallback
result, err := formatter.FormatWithContext(ctx)
if err != nil {
    // Try with simpler context
    fallbackCtx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": ctx.Data["path"],
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    result, err = formatter.FormatWithContext(fallbackCtx)
    if err != nil {
        log.Fatal(err)
    }
}

// Pattern 2: Use default formatting
result, err := formatter.FormatArchive(path, FormatTypeCreated)
if err != nil {
    // Use simple string formatting
    result = fmt.Sprintf("Created: %s", path)
}
```

## Testing Guidelines

### Mock Configuration

```go
type MockFormatterConfig struct {
    formatStrings   map[FormatType]string
    templateStrings map[FormatType]string
    errorFormats    map[ErrorType]string
    patterns        map[PatternType]string
}

func NewMockFormatterConfig() *MockFormatterConfig {
    return &MockFormatterConfig{
        formatStrings:   make(map[FormatType]string),
        templateStrings: make(map[FormatType]string),
        errorFormats:    make(map[ErrorType]string),
        patterns:        make(map[PatternType]string),
    }
}
```

### Test Structure

```go
func TestAIFirstFormatter(t *testing.T) {
    // 1. Setup
    config := NewMockFormatterConfig()
    formatter := NewAIFirstFormatter(config)
    
    // 2. Configure test data
    config.formatStrings[FormatTypeCreated] = "Created: %s"
    
    // 3. Create test context
    ctx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": "/test/file.txt",
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    // 4. Execute test
    result, err := formatter.FormatWithContext(ctx)
    
    // 5. Assert results
    if err != nil {
        t.Errorf("FormatWithContext failed: %v", err)
    }
    
    expected := "Created: /test/file.txt"
    if result != expected {
        t.Errorf("expected '%s', got '%s'", expected, result)
    }
}
```

### Integration Testing

```go
func TestAIFirstFormatterIntegration(t *testing.T) {
    // 1. Setup complete workflow
    config := NewMockFormatterConfig()
    collector := NewOutputCollector()
    formatter := NewAIFirstFormatterWithCollector(config, collector)
    
    // 2. Test format -> extract -> print workflow
    formatCtx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": "/test/archive.tar.gz",
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    result, err := formatter.FormatWithContext(formatCtx)
    if err != nil {
        t.Fatal(err)
    }
    
    printCtx := PrintContext{
        Message:     result,
        Destination: AIOutputDestinationStdout,
        Type:        AIMessageTypeInfo,
        Options: PrintOptions{
            Delayed: true,
            Flush:   false,
        },
        Metadata: make(map[string]string),
    }
    
    err = formatter.PrintWithContext(printCtx)
    if err != nil {
        t.Fatal(err)
    }
    
    // 3. Verify results
    messages := formatter.GetCollectedMessages()
    if len(messages) != 1 {
        t.Errorf("expected 1 message, got %d", len(messages))
    }
    
    if messages[0].Content != result {
        t.Errorf("collected message does not match formatted result")
    }
}
```

### Performance Testing

```go
func BenchmarkAIFirstFormatter(b *testing.B) {
    config := NewMockFormatterConfig()
    formatter := NewAIFirstFormatter(config)
    
    ctx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": "/test/archive.tar.gz",
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := formatter.FormatWithContext(ctx)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Best Practices for AI Assistants

1. **Use Context Structures**: Always use context structures for complex operations
2. **Handle All Errors**: Check error returns from all interface methods
3. **Validate Inputs**: Use the validation methods provided by interfaces
4. **Use Structured Data**: Prefer structured data types over raw strings
5. **Document Operations**: Use clear, descriptive operation names
6. **Test Thoroughly**: Write comprehensive tests for all interface implementations
7. **Follow Patterns**: Use the established usage patterns for consistency
8. **Provide Context**: Include relevant metadata in context structures
9. **Use Type Safety**: Leverage strongly typed enums and structures
10. **Graceful Degradation**: Implement fallback behavior when possible 