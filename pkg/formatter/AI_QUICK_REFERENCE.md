# AI-First Formatter Quick Reference

**For AI Assistants** - Quick reference for understanding and using the AI-first formatter interfaces.

## 🎯 Core Concept

The AI-first formatter uses **context structures** and **composition patterns** to provide clear, AI-friendly interfaces.

## 📋 Quick Interface Reference

### Primary Interface
```go
type AIFirstFormatter interface {
    CoreFormatter
    AIPatternExtractor  
    AIOutputManager
    
    SetConfig(config FormatterConfig) error
    GetConfig() FormatterConfig
    
    FormatWithContext(ctx FormatContext) (string, error)
    ExtractWithContext(ctx ExtractContext) (interface{}, error)
    PrintWithContext(ctx PrintContext) error
}
```

### Key Context Structures
```go
// For formatting operations
type FormatContext struct {
    FormatType    FormatType                `json:"format_type"`
    Data          map[string]interface{}    `json:"data"`
    Options       FormatOptions             `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}

// For data extraction
type ExtractContext struct {
    PatternType   PatternType               `json:"pattern_type"`
    Input         string                    `json:"input"`
    Options       ExtractOptions            `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}

// For output operations
type PrintContext struct {
    Message       string                    `json:"message"`
    Destination   AIOutputDestination       `json:"destination"`
    Type          AIMessageType             `json:"type"`
    Options       PrintOptions              `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

## 🚀 Common Usage Patterns

### 1. Basic Formatting
```go
// Create formatter
config := NewMockFormatterConfig()
formatter := NewAIFirstFormatter(config)

// Format with context
ctx := FormatContext{
    FormatType: FormatTypeCreated,
    Data: map[string]interface{}{
        "path": "/backups/archive.tar.gz",
    },
    Options:  FormatOptions{},
    Metadata: make(map[string]string),
}

result, err := formatter.FormatWithContext(ctx)
```

### 2. Pattern Extraction
```go
// Extract structured data
ctx := ExtractContext{
    PatternType: PatternTypeArchiveFilename,
    Input:       "myproject-2024-01-15-14-30-main-abc123-release.tar.gz",
    Options:     ExtractOptions{},
    Metadata:    make(map[string]string),
}

result, err := formatter.ExtractWithContext(ctx)
if archiveData, ok := result.(AIArchiveData); ok {
    // Use structured data
    fmt.Printf("Archive: %s, Branch: %s\n", archiveData.Prefix, archiveData.Branch)
}
```

### 3. Delayed Output
```go
// Create with collector
collector := NewOutputCollector()
formatter := NewAIFirstFormatterWithCollector(config, collector)

// Collect messages
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

// Flush all collected messages
err = formatter.Flush()
```

## 📊 Data Types Quick Reference

### Format Types
```go
FormatTypeCreated    // "created"
FormatTypeIdentical  // "identical" 
FormatTypeList       // "list"
FormatTypeDryRun     // "dry_run"
FormatTypeError      // "error"
FormatTypeConfig     // "config"
```

### Error Types
```go
ErrorTypeDiskFull        // "disk_full"
ErrorTypePermission      // "permission"
ErrorTypeDirectoryNotFound // "directory_not_found"
ErrorTypeFileNotFound    // "file_not_found"
ErrorTypeInvalidDirectory // "invalid_directory"
ErrorTypeInvalidFile     // "invalid_file"
ErrorTypeGeneric         // "generic"
```

### Pattern Types
```go
PatternTypeArchiveFilename // "archive_filename"
PatternTypeBackupFilename  // "backup_filename"
PatternTypeConfigLine      // "config_line"
PatternTypeTimestamp       // "timestamp"
```

### Output Destinations
```go
AIOutputDestinationStdout // "stdout"
AIOutputDestinationStderr // "stderr"
```

### Message Types
```go
AIMessageTypeInfo    // "info"
AIMessageTypeError   // "error"
AIMessageTypeWarning // "warning"
AIMessageTypeConfig  // "config"
```

## 🔧 Constructor Functions

```go
// Basic formatter
formatter := NewAIFirstFormatter(config)

// With delayed output
collector := NewOutputCollector()
formatter := NewAIFirstFormatterWithCollector(config, collector)

// Individual components
coreFormatter := NewAICoreFormatter(config)
patternExtractor := NewAIPatternExtractor(config)
outputManager := NewAIOutputManager()
outputManagerWithCollector := NewAIOutputManagerWithCollector(collector)
```

## 📝 Structured Data Types

### AIArchiveData
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

### AIBackupData
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

## ⚠️ Error Handling Patterns

### Always Check Errors
```go
result, err := formatter.FormatWithContext(ctx)
if err != nil {
    // Handle error appropriately
    log.Printf("Formatting failed: %v", err)
    return err
}
```

### Type Assertion for Structured Data
```go
result, err := formatter.ExtractWithContext(ctx)
if err != nil {
    return err
}

if archiveData, ok := result.(AIArchiveData); ok {
    // Use structured data
    fmt.Printf("Archive: %s\n", archiveData.Prefix)
} else {
    // Handle unexpected type
    return fmt.Errorf("expected AIArchiveData, got %T", result)
}
```

## 🧪 Testing Quick Reference

### Mock Configuration
```go
config := NewMockFormatterConfig()
config.formatStrings[FormatTypeCreated] = "Created: %s"
config.errorFormats[ErrorTypeGeneric] = "Error: %s"
```

### Test Structure
```go
func TestExample(t *testing.T) {
    // Setup
    config := NewMockFormatterConfig()
    formatter := NewAIFirstFormatter(config)
    
    // Configure test data
    config.formatStrings[FormatTypeCreated] = "Created: %s"
    
    // Create context
    ctx := FormatContext{
        FormatType: FormatTypeCreated,
        Data: map[string]interface{}{
            "path": "/test/file.txt",
        },
        Options:  FormatOptions{},
        Metadata: make(map[string]string),
    }
    
    // Execute
    result, err := formatter.FormatWithContext(ctx)
    
    // Assert
    if err != nil {
        t.Errorf("FormatWithContext failed: %v", err)
    }
    
    expected := "Created: /test/file.txt"
    if result != expected {
        t.Errorf("expected '%s', got '%s'", expected, result)
    }
}
```

## 🎯 AI Assistant Best Practices

### 1. Use Context Structures
- **Always** use context structures for complex operations
- Provide rich metadata for better AI comprehension
- Use structured data types over raw strings

### 2. Handle Errors Properly
- **Always** check error returns
- Provide meaningful error context
- Use type assertions for structured data

### 3. Follow Composition Pattern
- Use the primary `AIFirstFormatter` interface
- Leverage individual components when needed
- Maintain clear separation of concerns

### 4. Use Type Safety
- Use strongly typed enums (`FormatType`, `ErrorType`, etc.)
- Prefer structured data types over maps
- Validate data types before use

### 5. Provide Rich Context
- Include relevant metadata in context structures
- Use descriptive operation names
- Document complex operations

## 🔍 Common Patterns

### Format -> Extract -> Print Workflow
```go
// 1. Format message
formatCtx := FormatContext{
    FormatType: FormatTypeCreated,
    Data: map[string]interface{}{
        "path": "/backups/archive.tar.gz",
    },
    Options:  FormatOptions{},
    Metadata: make(map[string]string),
}

result, err := formatter.FormatWithContext(formatCtx)
if err != nil {
    return err
}

// 2. Extract data from filename
extractCtx := ExtractContext{
    PatternType: PatternTypeArchiveFilename,
    Input:       "myproject-2024-01-15-14-30-main-abc123-release.tar.gz",
    Options:     ExtractOptions{},
    Metadata:    make(map[string]string),
}

extractedData, err := formatter.ExtractWithContext(extractCtx)
if err != nil {
    return err
}

// 3. Print with context
printCtx := PrintContext{
    Message:     result,
    Destination: AIOutputDestinationStdout,
    Type:        AIMessageTypeInfo,
    Options:     PrintOptions{},
    Metadata:    make(map[string]string),
}

err = formatter.PrintWithContext(printCtx)
if err != nil {
    return err
}
```

### Error Handling with Context
```go
// Format error with context
errorCtx := FormatContext{
    FormatType: FormatTypeError,
    Data: map[string]interface{}{
        "error":     err,
        "errorType": ErrorTypeDiskFull,
    },
    Options:  FormatOptions{},
    Metadata: map[string]string{
        "operation": "backup",
        "user":      "admin",
    },
}

errorResult, err := formatter.FormatWithContext(errorCtx)
if err != nil {
    return err
}

// Print error to stderr
errorPrintCtx := PrintContext{
    Message:     errorResult,
    Destination: AIOutputDestinationStderr,
    Type:        AIMessageTypeError,
    Options:     PrintOptions{},
    Metadata:    make(map[string]string),
}

err = formatter.PrintWithContext(errorPrintCtx)
```

## 📚 Related Documentation

- **Full API Documentation**: `API.md`
- **Package Documentation**: `README.md`
- **Implementation Files**: 
  - `ai_first_interfaces.go` - Interface definitions
  - `ai_first_formatter.go` - Main implementation
  - `ai_core_formatter.go` - Core formatting
  - `ai_pattern_extractor.go` - Pattern extraction
  - `ai_output_manager.go` - Output management
  - `ai_first_formatter_test.go` - Comprehensive tests

## 🎯 Key Takeaways for AI Assistants

1. **Context is King**: Always use context structures for operations
2. **Error Handling**: Check all error returns and provide context
3. **Type Safety**: Use strongly typed enums and structures
4. **Composition**: Leverage the composition pattern for clean interfaces
5. **Structured Data**: Prefer structured data types over raw strings
6. **Testing**: Write comprehensive tests for all implementations
7. **Documentation**: Keep documentation updated with interface changes
8. **Performance**: Consider lazy initialization and memory efficiency
9. **Extensibility**: Design for future enhancements and plugin systems
10. **AI-Friendly**: Always consider AI assistant comprehension and usage 