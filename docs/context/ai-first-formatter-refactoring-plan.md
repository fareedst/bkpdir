# [CRITICAL] AI-First Formatter Refactoring Plan

**Implementation Token**: `// [CRITICAL] FMT-001: AI-first formatter refactoring`
**Purpose**: Establish AI-first data structures for formatter component that enable optimal AI assistant comprehension, testing, and maintenance.

## [PROBLEM] Current Formatter Issues

### [ANALYSIS] Mixed Concerns and Tight Coupling
**Current State**: `formatter.go` (1714 lines) with mixed responsibilities:
- **Formatting logic** mixed with **printing operations**
- **Pattern extraction** embedded in multiple components
- **Template processing** tightly coupled to configuration
- **Error formatting** scattered across multiple interfaces

### [IMPACT] AI Assistant Navigation Difficulties
1. **Complex component boundaries** - AI cannot easily understand responsibilities
2. **Inconsistent interface patterns** - Multiple overlapping interface definitions
3. **Tight configuration coupling** - Direct `*Config` struct dependencies
4. **Mixed abstraction levels** - Some components abstracted, others not

## [SOLUTION] AI-First Data Structure Design

### [ARCHITECTURE] Clear Component Separation

#### **1. Core Formatting Engine**
```go
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
}

// [CRITICAL] FMT-001: Format type enumeration - [ACTION:core-functionality]
type FormatType string

const (
    FormatTypeCreated    FormatType = "created"
    FormatTypeIdentical  FormatType = "identical"
    FormatTypeList       FormatType = "list"
    FormatTypeDryRun     FormatType = "dry_run"
    FormatTypeError      FormatType = "error"
)

// [CRITICAL] FMT-001: Error type enumeration - [ACTION:core-functionality]
type ErrorType string

const (
    ErrorTypeDiskFull        ErrorType = "disk_full"
    ErrorTypePermission      ErrorType = "permission"
    ErrorTypeDirectoryNotFound ErrorType = "directory_not_found"
    ErrorTypeFileNotFound    ErrorType = "file_not_found"
    ErrorTypeInvalidDirectory ErrorType = "invalid_directory"
    ErrorTypeInvalidFile     ErrorType = "invalid_file"
)
```

#### **2. Pattern Extraction Engine**
```go
// [CRITICAL] FMT-001: Pattern extraction engine - [ACTION:discovery]
type PatternExtractor interface {
    // Core extraction operations
    ExtractArchiveData(filename string) (ArchiveData, error)
    ExtractBackupData(filename string) (BackupData, error)
    ExtractConfigData(line string) (ConfigData, error)
    ExtractTimestampData(timestamp string) (TimestampData, error)
    
    // Generic pattern extraction
    ExtractPattern(pattern, text string) (map[string]string, error)
}

// [CRITICAL] FMT-001: Structured data types - [ACTION:core-functionality]
type ArchiveData struct {
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

type BackupData struct {
    Filename  string `json:"filename"`
    Year      string `json:"year"`
    Month     string `json:"month"`
    Day       string `json:"day"`
    Hour      string `json:"hour"`
    Minute    string `json:"minute"`
    Note      string `json:"note"`
}
```

#### **3. Output Management Engine**
```go
// [CRITICAL] FMT-001: Output management engine - [ACTION:format-processing]
type OutputManager interface {
    // Direct output operations
    Print(message string) error
    PrintError(message string) error
    
    // Delayed output operations
    Collect(message OutputMessage) error
    Flush() error
    FlushStdout() error
    FlushStderr() error
    Clear() error
    
    // Output state management
    IsDelayedMode() bool
    SetDelayedMode(enabled bool) error
    GetCollectedMessages() []OutputMessage
}

// [CRITICAL] FMT-001: Output message structure - [ACTION:core-functionality]
type OutputMessage struct {
    Content     string            `json:"content"`
    Destination OutputDestination  `json:"destination"`
    Type        MessageType       `json:"type"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    Timestamp   time.Time         `json:"timestamp"`
}

type OutputDestination string

const (
    OutputDestinationStdout OutputDestination = "stdout"
    OutputDestinationStderr OutputDestination = "stderr"
)

type MessageType string

const (
    MessageTypeInfo    MessageType = "info"
    MessageTypeError   MessageType = "error"
    MessageTypeWarning MessageType = "warning"
    MessageTypeConfig  MessageType = "config"
)
```

#### **4. Configuration Abstraction Layer**
```go
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

type PatternType string

const (
    PatternTypeArchiveFilename PatternType = "archive_filename"
    PatternTypeBackupFilename  PatternType = "backup_filename"
    PatternTypeConfigLine      PatternType = "config_line"
    PatternTypeTimestamp       PatternType = "timestamp"
)
```

### [INTEGRATION] AI-First Interface Composition

#### **Primary Formatter Interface**
```go
// [CRITICAL] FMT-001: Primary formatter interface - [ACTION:core-functionality]
type AIFirstFormatter interface {
    // Core formatting capabilities
    CoreFormatter
    PatternExtractor
    OutputManager
    
    // Configuration management
    SetConfig(config FormatterConfig) error
    GetConfig() FormatterConfig
    
    // AI-friendly operations
    FormatWithContext(ctx FormatContext) (string, error)
    ExtractWithContext(ctx ExtractContext) (interface{}, error)
    PrintWithContext(ctx PrintContext) error
}

// [CRITICAL] FMT-001: Context structures for AI comprehension - [ACTION:core-functionality]
type FormatContext struct {
    FormatType    FormatType                `json:"format_type"`
    Data          map[string]interface{}    `json:"data"`
    Options       FormatOptions             `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}

type ExtractContext struct {
    PatternType   PatternType               `json:"pattern_type"`
    Input         string                    `json:"input"`
    Options       ExtractOptions            `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}

type PrintContext struct {
    Message       string                    `json:"message"`
    Destination   OutputDestination         `json:"destination"`
    Type          MessageType               `json:"type"`
    Options       PrintOptions              `json:"options"`
    Metadata      map[string]string         `json:"metadata"`
}
```

## [IMPLEMENTATION] AI-First Refactoring Strategy

### [PHASE_1] Core Component Extraction (Week 1)
**Priority**: [CRITICAL] - Foundation for all other refactoring

1. **Extract CoreFormatter** - Pure formatting logic
2. **Extract PatternExtractor** - Pattern extraction logic
3. **Extract OutputManager** - Output handling logic
4. **Create FormatterConfig** - Configuration abstraction

### [PHASE_2] Interface Standardization (Week 2)
**Priority**: [HIGH] - Establish AI-friendly contracts

1. **Define AIFirstFormatter interface** - Primary composition interface
2. **Implement context structures** - AI-friendly operation contexts
3. **Create validation framework** - Interface compliance checking
4. **Establish error handling patterns** - Consistent error structures

### [PHASE_3] AI Assistant Integration (Week 3)
**Priority**: [HIGH] - Enable optimal AI comprehension

1. **Update AI assistant protocols** - Use new interface patterns
2. **Create AI-friendly documentation** - Clear component boundaries
3. **Implement AI navigation helpers** - Easy component discovery
4. **Establish testing patterns** - AI-friendly test structures

### [PHASE_4] Backward Compatibility (Week 4)
**Priority**: [MEDIUM] - Maintain existing functionality

1. **Create adapter layer** - Bridge old and new interfaces
2. **Implement gradual migration** - Incremental adoption
3. **Update existing tests** - Maintain test coverage
4. **Document migration guide** - Clear upgrade path

## [VALIDATION] AI-First Quality Metrics

### [COMPREHENSION] AI Assistant Understanding
- **Component Boundary Clarity**: Each component has single, clear responsibility
- **Interface Consistency**: All interfaces follow same patterns
- **Context Structure**: All operations use AI-friendly context structures
- **Error Handling**: Consistent error patterns across all components

### [MAINTAINABILITY] Code Quality Metrics
- **Interface Compliance**: 100% of components implement required interfaces
- **Context Usage**: 100% of operations use context structures
- **Error Consistency**: 100% of errors follow structured patterns
- **Documentation Coverage**: 100% of interfaces have AI-friendly documentation

### [TESTABILITY] Testing Framework
- **Unit Test Coverage**: >95% for all new components
- **Interface Testing**: 100% of interfaces have contract tests
- **Context Testing**: All context structures have validation tests
- **Integration Testing**: End-to-end formatter workflow tests

## [SUCCESS] Expected Outcomes

### [IMMEDIATE] AI Assistant Benefits
1. **Clear Component Boundaries** - AI can easily understand responsibilities
2. **Consistent Interface Patterns** - Predictable API design
3. **Context-Driven Operations** - AI-friendly operation contexts
4. **Structured Error Handling** - Consistent error patterns

### [LONG_TERM] System Benefits
1. **Improved Maintainability** - Clear separation of concerns
2. **Enhanced Testability** - Focused component testing
3. **Better Extensibility** - Easy to add new formatting capabilities
4. **AI-First Design** - Optimized for AI assistant comprehension

---

**Implementation Status**: [ACTION:format-processing] **PLANNING** - Ready for immediate implementation
**Priority**: [CRITICAL] - Foundation for AI-first development
**Dependencies**: None - Can be implemented independently
**Timeline**: 4 weeks for complete implementation 