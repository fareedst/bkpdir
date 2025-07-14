# Formatter Package

The formatter package provides comprehensive output formatting functionality for BkpDir operations. It supports both printf-style and template-based formatting with optional delayed output collection.

## Architecture Overview

### AI-First Design Principles

The formatter package follows AI-first design principles to ensure optimal AI assistant comprehension and maintenance:

1. **Clear Component Separation**: Each component has a single, well-defined responsibility
2. **Interface Standardization**: Consistent interface patterns across all components
3. **Structured Data Types**: JSON-serializable data structures for easy AI comprehension
4. **Context-Aware Operations**: AI-friendly context structures for complex operations
5. **Comprehensive Error Handling**: Consistent error patterns with detailed context

### Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AIFirstFormatter                        │
│  (Primary Interface - Composes All Capabilities)          │
├─────────────────────────────────────────────────────────────┤
│  CoreFormatter  │  AIPatternExtractor  │  AIOutputManager │
│  (Pure Logic)   │  (Data Extraction)   │  (I/O Handling) │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    ┌─────────────────┐
                    │ FormatterConfig │
                    │ (Configuration) │
                    └─────────────────┘
```

## Core Interfaces

### AIFirstFormatter

The primary interface that combines all formatting capabilities:

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

**Key Features:**
- **Composition Pattern**: Combines three specialized interfaces
- **Configuration Management**: Dynamic configuration updates
- **Context-Aware Operations**: AI-friendly operation contexts
- **Error Handling**: Comprehensive error context and validation

### CoreFormatter

Provides pure formatting operations without side effects:

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

**Key Features:**
- **Pure Functions**: No side effects, deterministic output
- **Type Safety**: Strongly typed format and error types
- **Template Support**: Both Go templates and placeholder substitution
- **Context Awareness**: AI-friendly context structures

### AIPatternExtractor

Provides structured data extraction from text patterns:

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

**Key Features:**
- **Structured Data**: JSON-serializable data structures
- **Type Safety**: Strongly typed return values
- **Pattern Flexibility**: Configurable regex patterns
- **Error Context**: Detailed error information

### AIOutputManager

Provides output handling with delayed output support:

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

**Key Features:**
- **Delayed Output**: Collect and display messages later
- **Destination Control**: Separate stdout/stderr handling
- **State Management**: Dynamic mode switching
- **Message Metadata**: Rich message context

## Data Types

### Format Types

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

### Error Types

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

### Structured Data Types

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

## Context Structures

### FormatContext

AI-friendly context for formatting operations:

```go
type FormatContext struct {
    FormatType    FormatType                `json:"format_type"`
    Data          map[string]interface{}    `json:"data"`
    Options       FormatOptions             `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

### ExtractContext

AI-friendly context for data extraction:

```go
type ExtractContext struct {
    PatternType   PatternType               `json:"pattern_type"`
    Input         string                    `json:"input"`
    Options       ExtractOptions            `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

### PrintContext

AI-friendly context for output operations:

```go
type PrintContext struct {
    Message       string                    `json:"message"`
    Destination   AIOutputDestination       `json:"destination"`
    Type          AIMessageType             `json:"type"`
    Options       PrintOptions              `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

## Usage Examples

### Basic Usage

```go
// Create formatter with configuration
config := NewMockFormatterConfig()
formatter := NewAIFirstFormatter(config)

// Format with context
ctx := FormatContext{
    FormatType: FormatTypeCreated,
    Data: map[string]interface{}{
        "path": "/test/archive.tar.gz",
    },
    Options:  FormatOptions{},
    Metadata: make(map[string]string),
}

result, err := formatter.FormatWithContext(ctx)
if err != nil {
    log.Fatal(err)
}
```

### Pattern Extraction

```go
// Extract archive data
ctx := ExtractContext{
    PatternType: PatternTypeArchiveFilename,
    Input:       "test-2024-01-15-14-30-main-abc123-note.tar.gz",
    Options:     ExtractOptions{},
    Metadata:    make(map[string]string),
}

result, err := formatter.ExtractWithContext(ctx)
if err != nil {
    log.Fatal(err)
}

if archiveData, ok := result.(AIArchiveData); ok {
    fmt.Printf("Archive: %s, Year: %s\n", archiveData.Prefix, archiveData.Year)
}
```

### Delayed Output

```go
// Create formatter with collector
collector := NewOutputCollector()
formatter := NewAIFirstFormatterWithCollector(config, collector)

// Print with context
printCtx := PrintContext{
    Message:     "Processing archive...",
    Destination: AIOutputDestinationStdout,
    Type:        AIMessageTypeInfo,
    Options: PrintOptions{
        Delayed: true,
        Flush:   false,
    },
    Metadata: make(map[string]string),
}

err := formatter.PrintWithContext(printCtx)
if err != nil {
    log.Fatal(err)
}

// Flush all collected messages
err = formatter.Flush()
if err != nil {
    log.Fatal(err)
}
```

## Configuration

### FormatterConfig Interface

```go
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
```

### Configuration Error Handling

```go
type ConfigError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Code    string `json:"code"`
}
```

## Testing

### Mock Configuration

```go
type MockFormatterConfig struct {
    formatStrings   map[FormatType]string
    templateStrings map[FormatType]string
    errorFormats    map[ErrorType]string
    patterns        map[PatternType]string
}
```

### Test Examples

```go
func TestAIFirstFormatter(t *testing.T) {
    config := NewMockFormatterConfig()
    formatter := NewAIFirstFormatter(config)
    
    // Test format with context
    ctx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": "/test/archive.tar.gz",
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    result, err := formatter.FormatWithContext(ctx)
    if err != nil {
        t.Errorf("FormatWithContext failed: %v", err)
    }
    
    if result == "" {
        t.Error("FormatWithContext returned empty result")
    }
}
```

## Migration Guide

### From Legacy Formatter

1. **Replace Direct Usage**:
   ```go
   // Old
   formatter := NewOutputFormatter(config)
   result := formatter.FormatCreatedArchive(path)
   
   // New
   formatter := NewAIFirstFormatter(config)
   ctx := FormatContext{
       FormatType: FormatTypeCreated,
       Data: map[string]interface{}{"path": path},
       Options: FormatOptions{},
       Metadata: make(map[string]string),
   }
   result, err := formatter.FormatWithContext(ctx)
   ```

2. **Update Error Handling**:
   ```go
   // Old
   result := formatter.FormatError(err)
   
   // New
   result, err := formatter.FormatError(err, ErrorTypeGeneric)
   ```

3. **Migrate Pattern Extraction**:
   ```go
   // Old
   data := formatter.ExtractArchiveFilenameData(filename)
   
   // New
   data, err := formatter.ExtractArchiveData(filename)
   ```

## Best Practices

### For AI Assistants

1. **Use Context Structures**: Always use context structures for complex operations
2. **Handle Errors**: Check all error returns and provide meaningful context
3. **Validate Inputs**: Use the validation methods provided by interfaces
4. **Use Structured Data**: Prefer structured data types over raw strings
5. **Document Operations**: Use clear, descriptive operation names

### For Developers

1. **Interface Composition**: Compose interfaces rather than creating monolithic types
2. **Error Context**: Provide detailed error information with context
3. **Type Safety**: Use strongly typed enums and structures
4. **Testing**: Write comprehensive tests for all interface implementations
5. **Documentation**: Keep documentation updated with interface changes

## Performance Considerations

1. **Lazy Initialization**: Components are created on-demand
2. **Memory Efficiency**: Structured data types are optimized for JSON serialization
3. **Error Handling**: Errors are handled efficiently without excessive allocations
4. **Context Reuse**: Context structures can be reused for multiple operations

## Future Enhancements

1. **Async Operations**: Support for asynchronous formatting operations
2. **Streaming Output**: Real-time output streaming capabilities
3. **Custom Patterns**: User-defined pattern extraction rules
4. **Performance Metrics**: Built-in performance monitoring
5. **Plugin System**: Extensible formatter plugin architecture 