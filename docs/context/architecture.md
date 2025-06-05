# BkpDir Architecture

## Overview

BkpDir is a directory archiving and file backup tool that creates ZIP-based archives and individual file backups with Git integration, file exclusion patterns, and comprehensive verification capabilities. The architecture incorporates advanced features including structured error handling, resource management, template-based formatting, context-aware operations, enhanced configuration management with environment variable support, and dual-mode operation for both directory archiving and individual file backup.

## Core Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                            │
├─────────────────────────────────────────────────────────────┤
│  Command Handlers  │  Flag Processing  │  Output Formatting │
│  Context Support   │  Error Handling   │  Resource Cleanup  │
│  Dual-Mode Ops     │  Backward Compat  │  File Operations   │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Layer                      │
├─────────────────────────────────────────────────────────────┤
│  Config Discovery  │  Environment Vars │  Status Codes      │
│  Template Configs  │  Format Strings   │  Regex Patterns    │
│  Dual-Mode Config  │  File Backup Cfg  │  Archive Config    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                      Core Services                          │
├─────────────────────────────────────────────────────────────┤
│  Archive Service   │  Git Service      │  Verification      │
│  Resource Manager  │  Template Engine  │  Error Handler     │
│  Context Manager   │  File Operations  │  Backup Service    │
│  Comparison Svc    │  Formatter Svc    │  Config Service    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                     Storage Layer                           │
├─────────────────────────────────────────────────────────────┤
│  File System        │  ZIP Archives    │  Checksums         │
│  File Backups       │  Metadata        │  Verification      │
│  Directory Trees    │  Atomic Ops      │  Resource Cleanup  │
└─────────────────────────────────────────────────────────────┘
```

## Data Models

### Core Configuration Object

```go
type Config struct {
    // Directory archiving settings
    ArchiveDirPath    string   `yaml:"archive_dir_path"`
    UseCurrentDirName bool     `yaml:"use_current_dir_name"`
    ExcludePatterns   []string `yaml:"exclude_patterns"`
    IncludeGitInfo    bool     `yaml:"include_git_info"`
    
    // File backup settings
    BackupDirPath             string `yaml:"backup_dir_path"`
    UseCurrentDirNameForFiles bool   `yaml:"use_current_dir_name_for_files"`
    
    // Verification settings
    Verification *VerificationConfig `yaml:"verification"`
    
    // Status codes for directory operations
    StatusCreatedArchive                        int `yaml:"status_created_archive"`
    StatusFailedToCreateArchiveDirectory        int `yaml:"status_failed_to_create_archive_directory"`
    StatusDirectoryIsIdenticalToExistingArchive int `yaml:"status_directory_is_identical_to_existing_archive"`
    StatusDirectoryNotFound                     int `yaml:"status_directory_not_found"`
    StatusInvalidDirectoryType                  int `yaml:"status_invalid_directory_type"`
    StatusPermissionDenied                      int `yaml:"status_permission_denied"`
    StatusDiskFull                              int `yaml:"status_disk_full"`
    StatusConfigError                           int `yaml:"status_config_error"`
    
    // Status codes for file operations
    StatusCreatedBackup                   int `yaml:"status_created_backup"`
    StatusFailedToCreateBackupDirectory   int `yaml:"status_failed_to_create_backup_directory"`
    StatusFileIsIdenticalToExistingBackup int `yaml:"status_file_is_identical_to_existing_backup"`
    StatusFileNotFound                    int `yaml:"status_file_not_found"`
    StatusInvalidFileType                 int `yaml:"status_invalid_file_type"`
    
    // Printf-style format strings for directory operations
    FormatCreatedArchive   string `yaml:"format_created_archive"`
    FormatIdenticalArchive string `yaml:"format_identical_archive"`
    FormatListArchive      string `yaml:"format_list_archive"`
    FormatConfigValue      string `yaml:"format_config_value"`
    FormatDryRunArchive    string `yaml:"format_dry_run_archive"`
    FormatError            string `yaml:"format_error"`
    
    // Printf-style format strings for file operations
    FormatCreatedBackup   string `yaml:"format_created_backup"`
    FormatIdenticalBackup string `yaml:"format_identical_backup"`
    FormatListBackup      string `yaml:"format_list_backup"`
    FormatDryRunBackup    string `yaml:"format_dry_run_backup"`
    
    // Template-based format strings for directory operations
    TemplateCreatedArchive   string `yaml:"template_created_archive"`
    TemplateIdenticalArchive string `yaml:"template_identical_archive"`
    TemplateListArchive      string `yaml:"template_list_archive"`
    TemplateConfigValue      string `yaml:"template_config_value"`
    TemplateDryRunArchive    string `yaml:"template_dry_run_archive"`
    TemplateError            string `yaml:"template_error"`
    
    // Template-based format strings for file operations
    TemplateCreatedBackup   string `yaml:"template_created_backup"`
    TemplateIdenticalBackup string `yaml:"template_identical_backup"`
    TemplateListBackup      string `yaml:"template_list_backup"`
    TemplateDryRunBackup    string `yaml:"template_dry_run_backup"`
    
    // Regex patterns for data extraction
    PatternArchiveFilename string `yaml:"pattern_archive_filename"`
    PatternBackupFilename  string `yaml:"pattern_backup_filename"`
    PatternConfigLine      string `yaml:"pattern_config_line"`
    PatternTimestamp       string `yaml:"pattern_timestamp"`
}

type VerificationConfig struct {
    VerifyOnCreate    bool   `yaml:"verify_on_create"`
    ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}
```

### Enhanced Data Objects

```go
type ConfigValue struct {
    Name   string
    Value  string
    Source string
}

type ArchiveError struct {
    Message    string
    StatusCode int
    Operation  string
    Path       string
    Err        error
}

type BackupError struct {
    Message    string
    StatusCode int
    Operation  string
    Path       string
    Err        error
}

type ResourceManager struct {
    resources []Resource
    mutex     sync.RWMutex
}

type TemplateFormatter struct {
    config *Config
}

type OutputFormatter struct {
    cfg *Config
}

type TemplateFormatter struct {
    config *Config
}

type Backup struct {
    Name         string
    Path         string
    CreationTime time.Time
    SourceFile   string
    Note         string
}

type Archive struct {
    Name               string
    Path               string
    CreationTime       time.Time
    IsIncremental      bool
    GitBranch          string
    GitHash            string
    Note               string
    BaseArchive        string
    VerificationStatus *VerificationStatus
}

type BackupInfo struct {
    Name         string
    Path         string
    CreationTime time.Time
    Size         int64
}

type BackupOptions struct {
    Context   context.Context
    Config    *Config
    Formatter *OutputFormatter
    FilePath  string
    Note      string
    DryRun    bool
}
```

## Service Architecture

### Archive Service

**Responsibilities:**
- Directory scanning and file collection with context support
- ZIP archive creation with compression and cancellation
- Incremental archive logic with base archive tracking
- File exclusion pattern matching with doublestar patterns
- Archive naming with timestamps, Git info, and notes
- Context-aware operations with cancellation support

**Key Components:**
- `ArchiveCreator`: Handles full and incremental archives with context
- `FileScanner`: Traverses directories with exclusion patterns and cancellation checks
- `CompressionEngine`: ZIP creation with optimal compression and context support
- `NamingService`: Generates archive names with metadata and Git integration

**Context-Aware Functions:**
- `CreateArchiveWithContext()`: Archive creation with cancellation support
- `CreateFullArchiveWithContext()`: Full archive with context and resource cleanup
- `createZipArchiveWithContext()`: ZIP creation with periodic cancellation checks

### File Backup Service

**Responsibilities:**
- Single file backup creation with atomic operations
- File comparison for identical backup detection
- Directory structure preservation in backup paths
- Context-aware backup operations with cancellation
- Enhanced error handling and resource cleanup

**Key Components:**
- `BackupCreator`: Handles individual file backup creation
- `FileComparator`: Compares files for identical backup detection
- `BackupNaming`: Generates backup names with timestamps and notes
- `BackupListing`: Lists and sorts file backups by creation time

**Context-Aware Functions:**
- `CreateFileBackupWithContext()`: File backup with cancellation support
- `CopyFileWithContext()`: File copying with periodic cancellation checks
- `CheckForIdenticalFileBackup()`: Backup comparison with existing files

### Git Integration Service

**Responsibilities:**
- Git repository detection and validation
- Branch and commit hash extraction
- Git status checking for clean state
- Integration with archive and backup naming

**Key Components:**
- `GitDetector`: Repository discovery and validation
- `BranchExtractor`: Current branch identification
- `HashExtractor`: Commit hash retrieval
- `StatusChecker`: Working directory state analysis

### Resource Management Service

**Responsibilities:**
- Automatic cleanup of temporary files and directories
- Thread-safe resource tracking with mutex protection
- Error-resilient cleanup operations that continue on failures
- Context-aware resource management with cancellation support
- Panic recovery mechanisms for robust cleanup

**Key Components:**
- `ResourceTracker`: Maintains thread-safe resource registry
- `CleanupManager`: Handles automatic cleanup with error resilience
- `TempFileManager`: Temporary file lifecycle management
- `AtomicOperations`: Safe file operations with cleanup tracking

**Resource Types:**
- `TempFile`: Temporary file resource with cleanup
- `TempDir`: Temporary directory resource with recursive cleanup
- `Resource`: Interface for all cleanable resources

### Template Formatting Service

**Responsibilities:**
- Advanced template-based output formatting with named placeholders
- Regex pattern extraction for rich data formatting from filenames and paths
- Support for both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
- ANSI color support for enhanced readability and text highlighting
- Fallback mechanisms for invalid templates with graceful degradation
- Operation context integration for enhanced error messages

**Key Components:**
- `TemplateFormatter`: Advanced template processing with regex extraction and data binding
- `OutputFormatter`: Printf-style and template-based formatting with dual-mode support
- `PlaceholderResolver`: Named parameter substitution with %{name} syntax
- `PatternExtractor`: Regex-based data extraction with named groups for archives and backups
- `ColorFormatter`: ANSI color code handling and text highlighting support
- `ContextManager`: Operation context tracking for enhanced error display

**Template Methods for Directory Operations:**
- `TemplateCreatedArchive()`: Format archive creation with extracted metadata
- `TemplateIdenticalArchive()`: Format identical directory with archive information
- `TemplateListArchive()`: Format archive listing with rich data extraction
- `TemplateDryRunArchive()`: Format dry-run with planned archive information

**Template Methods for File Operations:**
- `TemplateCreatedBackup()`: Format backup creation with extracted filename data
- `TemplateIdenticalBackup()`: Format identical file with backup information
- `TemplateListBackup()`: Format backup listing with rich data extraction
- `TemplateDryRunBackup()`: Format dry-run with planned backup information

**Common Template Methods:**
- `TemplateConfigValue()`: Format configuration display with conditional formatting
- `TemplateError()`: Format error messages with operation context
- `FormatWithTemplate()`: Apply Go text/template to regex-extracted data
- `FormatWithPlaceholders()`: Replace %{name} placeholders with data values

**Print Methods:**
- Direct printing methods for all template operations: `PrintTemplateCreatedArchive()`, etc.
- Stdout and stderr routing based on message type
- Consistent formatting across all output types

### Output Formatting Service

**Responsibilities:**
- Centralized printf-style formatting for all application output
- Consistent message formatting across directory and file operations
- ANSI color code support for enhanced readability
- Configurable format strings with fallback defaults
- Integration with template formatting for dual-mode support

**Key Components:**
- `OutputFormatter`: Printf-style formatting with configuration integration
- `MessageRouter`: Stdout/stderr routing based on message type
- `FormatValidator`: Format string validation and error handling
- `ColorManager`: ANSI color code management and terminal detection

**Format Methods for Directory Operations:**
- `FormatCreatedArchive()`: Format archive creation messages
- `FormatIdenticalArchive()`: Format identical directory messages
- `FormatListArchive()`: Format archive listing entries
- `FormatDryRunArchive()`: Format dry-run archive messages

**Format Methods for File Operations:**
- `FormatCreatedBackup()`: Format backup creation messages
- `FormatIdenticalBackup()`: Format identical file messages
- `FormatListBackup()`: Format backup listing entries
- `FormatDryRunBackup()`: Format dry-run backup messages

**Common Format Methods:**
- `FormatConfigValue()`: Format configuration value display
- `FormatError()`: Format error messages
- Print methods for all operations: `PrintCreatedArchive()`, `PrintError()`, etc.

### Error Handling Service

**Responsibilities:**
- Structured error creation with status codes and context
- Enhanced error detection for disk space and permissions
- Status code management with configurable exit codes
- Error context preservation with operation and path information
- Panic recovery mechanisms for robust error handling

**Key Components:**
- `ErrorBuilder`: Structured error construction with context
- `StatusManager`: Exit code configuration and management
- `ContextPreserver`: Error context tracking with operation details
- `PanicHandler`: Recovery mechanisms with cleanup
- `ErrorDetector`: Enhanced error pattern detection

**Error Detection Functions:**
- `IsDiskFullError()`: Comprehensive disk space error detection
- `IsPermissionError()`: Permission and access error detection
- `IsDirectoryNotFoundError()`: Directory existence error detection

## Configuration Architecture

### Configuration Discovery

```
1. Command line flags (highest priority)
2. Environment variables (BKPDIR_CONFIG for search paths)
3. Configuration files in search path order:
   - ./.bkpdir.yml (current directory)
   - ~/.bkpdir.yml (home directory)
   - Custom paths from BKPDIR_CONFIG
4. Default values (lowest priority)
```

### Configuration Sources

**Environment Variable Support:**
- `BKPDIR_CONFIG`: Colon-separated list of configuration file paths
- Home directory expansion support with `~` prefix
- Configurable search path with precedence rules
- Earlier files in search path take precedence over later files

**File-based Configuration:**
- YAML format with comprehensive validation
- Schema-based structure with all documented fields
- Default value inheritance and merging
- Status code configuration with defaults
- Format string configuration with printf-style and template support
- Regex pattern configuration for data extraction

**Configuration Merging:**
- Multiple configuration files processed in search path order
- Earlier files override later files for non-default values
- Default values preserved when not explicitly overridden
- Status codes, format strings, and templates all configurable

### ⭐ Layered Configuration Inheritance Architecture (CFG-005)

**🔻 CFG-005: Layered Configuration Inheritance Implementation**

The configuration system supports layered inheritance enabling configuration files to inherit from other configuration files with flexible merge strategies. This provides hierarchical configuration management for complex project structures.

#### Inheritance Configuration Structure

```go
// Configuration inheritance metadata
type ConfigInheritance struct {
    Inherit []string `yaml:"inherit"` // List of parent configuration files
}

// Inheritable configuration wrapper
type InheritableConfig struct {
    *Config `yaml:",inline"`         // Inline main configuration
    ConfigInheritance `yaml:",inline"` // Inline inheritance metadata
}

// Inheritance chain tracking
type InheritanceChain struct {
    Files   []string           // Configuration files in dependency order
    Visited map[string]bool    // Circular dependency prevention
    Sources map[string]string  // Source tracking for debugging
}
```

#### Merge Strategy Architecture

```go
// Merge strategy processor
type MergeStrategyProcessor struct {
    strategies map[string]MergeStrategy
}

// Merge strategy interface
type MergeStrategy interface {
    Merge(dest, src reflect.Value) error
    GetPrefix() string
    GetDescription() string
}

// Strategy implementations
type StandardOverrideStrategy struct{}  // No prefix - replace parent values
type ArrayMergeStrategy struct{}        // + prefix - append to parent arrays
type ArrayPrependStrategy struct{}      // ^ prefix - prepend to parent arrays
type ArrayReplaceStrategy struct{}      // ! prefix - replace parent arrays
type DefaultValueStrategy struct{}      // = prefix - use if parent not set
```

#### Configuration Loading Engine Enhancement

```go
// Enhanced configuration loader with inheritance support
type InheritanceConfigLoader struct {
    baseLoader     ConfigLoader
    pathResolver   PathResolver
    chainBuilder   InheritanceChainBuilder
    strategyEngine MergeStrategyProcessor
    circularDetector CircularDependencyDetector
}

// Inheritance chain processing
func (l *InheritanceConfigLoader) LoadConfigWithInheritance(root string) (*Config, error) {
    // 1. Build inheritance dependency graph
    chain, err := l.buildInheritanceChain(root)
    if err != nil {
        return nil, err
    }
    
    // 2. Load configurations in dependency order
    configs, err := l.loadConfigurationChain(chain)
    if err != nil {
        return nil, err
    }
    
    // 3. Apply merge strategies and resolve final configuration
    return l.mergeConfigurationChain(configs)
}
```

#### Merge Strategy Processing

```go
// Key prefix analysis for merge strategy determination
type PrefixedKeyProcessor struct {
    prefixMap map[string]MergeStrategy
}

func (p *PrefixedKeyProcessor) ProcessPrefixedKeys(config map[string]interface{}) ProcessedConfig {
    processed := ProcessedConfig{
        StandardKeys: make(map[string]interface{}),
        MergeOps:     make([]MergeOperation, 0),
    }
    
    for key, value := range config {
        strategy, cleanKey := p.extractMergeStrategy(key)
        processed.MergeOps = append(processed.MergeOps, MergeOperation{
            Key:      cleanKey,
            Value:    value,
            Strategy: strategy,
        })
    }
    
    return processed
}
```

#### Circular Dependency Detection

```go
// Depth-first search for circular dependency detection
type CircularDependencyDetector struct {
    visited map[string]bool
    inStack map[string]bool
    path    []string
}

func (d *CircularDependencyDetector) DetectCycle(startFile string, resolver PathResolver) error {
    d.visited = make(map[string]bool)
    d.inStack = make(map[string]bool)
    d.path = make([]string, 0)
    
    return d.dfsDetectCycle(startFile, resolver)
}
```

#### Configuration Source Tracking

```go
// Enhanced source tracking with inheritance information
type InheritanceSourceTracker struct {
    configSources map[string]ConfigSource
    valueOrigins  map[string]ValueOrigin
    chainMetadata InheritanceChainMetadata
}

type ValueOrigin struct {
    SourceFile    string    // Configuration file containing the value
    SourceType    string    // "inherit", "override", "merge", "prepend", "replace", "default"
    MergeStrategy string    // Merge strategy applied
    LineNumber    int       // Line number in source file
    ChainDepth    int       // Depth in inheritance chain
}
```

#### Integration with Existing Configuration System

```go
// Backward compatibility layer
type ConfigAdapterWithInheritance struct {
    inheritanceLoader *InheritanceConfigLoader
    legacyLoader      *DefaultConfigLoader
    inheritanceEnabled bool
}

func (a *ConfigAdapterWithInheritance) LoadConfig(root string) (*Config, error) {
    // Check if inheritance is being used
    if a.hasInheritanceDeclarations(root) {
        return a.inheritanceLoader.LoadConfigWithInheritance(root)
    }
    
    // Fall back to legacy loading for non-inheritance configurations
    return a.legacyLoader.LoadConfig(root)
}
```

#### Performance Optimization

```go
// Inheritance caching for performance
type InheritanceCache struct {
    chainCache    map[string]*InheritanceChain    // Cached dependency chains
    configCache   map[string]*Config              // Cached resolved configurations
    mTimeTracker  map[string]time.Time            // File modification time tracking
    cacheEnabled  bool
    maxCacheSize  int
}

// Cache invalidation based on file modification times
func (c *InheritanceCache) IsCacheValid(filePath string) bool {
    if !c.cacheEnabled {
        return false
    }
    
    cachedTime, exists := c.mTimeTracker[filePath]
    if !exists {
        return false
    }
    
    currentStat, err := os.Stat(filePath)
    if err != nil {
        return false
    }
    
    return currentStat.ModTime().Equal(cachedTime)
}
```

#### Error Handling and Debugging

```go
// Inheritance-specific error types
type InheritanceError struct {
    Type        string    // "circular", "missing_file", "invalid_syntax", "merge_conflict"
    Message     string    // Human-readable error message
    ChainPath   []string  // Files involved in inheritance chain
    SourceFile  string    // File where error originated
    LineNumber  int       // Line number of problematic declaration
}

// Inheritance chain visualization for debugging
type InheritanceChainVisualizer struct {
    chain    *InheritanceChain
    formatter OutputFormatter
}

func (v *InheritanceChainVisualizer) GenerateChainDiagram() string {
    // Generate ASCII art diagram of inheritance relationships
    // Include merge strategies and value origins
    // Highlight conflicts and resolution order
}
```

#### Configuration Example Processing

```yaml
# ~/.bkpdir.yml (base configuration)
archive_dir_path: "~/Archives"
backup_dir_path: "~/Backups"
exclude_patterns:
  - "*.tmp"
  - "*.log"
  - ".DS_Store"
include_git_info: true

# /project/.bkpdir.yml (inheriting configuration)
inherit: "~/.bkpdir.yml"

# Standard override
archive_dir_path: "./project-archives"
backup_dir_path: "./project-backups"

# Array merge strategy (append)
+exclude_patterns:
  - "node_modules/"
  - "dist/"
  - "*.secret"

# Array prepend strategy (high priority)
^exclude_patterns:
  - "*.private"

# Default strategy (use if not set)
=verification:
  verify_on_create: true
```

**Resolved Configuration Result:**
```yaml
archive_dir_path: "./project-archives"  # Overridden by child
backup_dir_path: "./project-backups"    # Overridden by child
exclude_patterns:                       # Merged: prepend + base + append
  - "*.private"                         # ^ prepend strategy
  - "*.tmp"                             # Base configuration
  - "*.log"                             # Base configuration
  - ".DS_Store"                         # Base configuration
  - "node_modules/"                     # + merge strategy
  - "dist/"                             # + merge strategy
  - "*.secret"                          # + merge strategy
include_git_info: true                  # Inherited from base
verification:                           # = default strategy applied
  verify_on_create: true
```

This architecture provides comprehensive inheritance capabilities while maintaining backward compatibility, performance, and robust error handling for complex configuration scenarios.

## Archive Format Architecture

### Archive Naming Convention

```
Format: [prefix-]timestamp[=branch=hash][=note].zip

Examples:
- backup-2024-03-15-14-30.zip
- proj-2024-03-15-14-30=main=a1b2c3d.zip
- data-2024-03-15-14-30=feature=x9y8z7w=release.zip

Incremental Format: basename_update=timestamp[=branch=hash][=note].zip
- proj-2024-03-15-14-30_update=2024-03-15-16-45=main=def456.zip
```

### File Backup Naming Convention

```
Format: filename-timestamp[=note]

Examples:
- document.txt-2024-03-15-14-30
- config.yml-2024-03-15-14-30=before-changes
- script.sh-2024-03-15-14-30=working-version
```

### Archive Structure

```
archive.zip
├── metadata.json          # Archive metadata
├── checksums.sha256      # File checksums
├── git-info.json         # Git repository info (if enabled)
└── content/              # Archived directory content
    ├── file1.txt
    ├── subdir/
    │   └── file2.txt
    └── ...
```

### Metadata Format

```json
{
    "created_at": "2024-03-15T14:30:22Z",
    "source_path": "/path/to/source",
    "archive_type": "full|incremental",
    "git_info": {
        "branch": "main",
        "commit": "a1b2c3d4e5f6",
        "is_clean": true
    },
    "excluded_patterns": [".git/", "vendor/"],
    "file_count": 1234,
    "total_size": 5678901
}
```

## Output Formatting Architecture

### Printf-Style Formatting

```go
// Format string configuration for directory operations
FormatStrings: map[string]string{
    "format_created_archive": "Created archive: %s\n",
    "format_identical_archive": "Directory is identical to existing archive: %s\n",
    "format_list_archive": "%s (created: %s)\n",
    "format_config_value": "%s: %s (source: %s)\n",
    "format_dry_run_archive": "Would create archive: %s\n",
    "format_error": "Error: %s\n",
}

// Format string configuration for file operations
FileFormatStrings: map[string]string{
    "format_created_backup": "Created backup: %s\n",
    "format_identical_backup": "File is identical to existing backup: %s\n",
    "format_list_backup": "%s (created: %s)\n",
    "format_dry_run_backup": "Would create backup: %s\n",
}
```

### Template-Based Formatting

```go
// Template configuration for directory operations
Templates: map[string]string{
    "template_created_archive": "Created archive: %{path}\n",
    "template_identical_archive": "Directory is identical to existing archive: %{path}\n",
    "template_list_archive": "%{path} (created: %{creation_time})\n",
    "template_config_value": "%{name}: %{value} (source: %{source})\n",
    "template_dry_run_archive": "Would create archive: %{path}\n",
    "template_error": "Error: %{message}\n",
}

// Template configuration for file operations
FileTemplates: map[string]string{
    "template_created_backup": "Created backup: %{path}\n",
    "template_identical_backup": "File is identical to existing backup: %{path}\n",
    "template_list_backup": "%{path} (created: %{creation_time})\n",
    "template_dry_run_backup": "Would create backup: %{path}\n",
}
```

### Regex Pattern Integration

```go
// Named regex patterns for data extraction
Patterns: map[string]string{
    "pattern_archive_filename": `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`,
    "pattern_backup_filename": `(?P<filename>[^/]+)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<note>.+))?`,
    "pattern_config_line": `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`,
    "pattern_timestamp": `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`,
}
```

## Error Handling Architecture

### Error Categories

1. **Configuration Errors** (Status: 10)
   - Invalid configuration files
   - Missing required settings
   - Environment variable issues

2. **File System Errors** (Status: 20-22)
   - Directory/file not found (20)
   - Invalid directory/file type (21)
   - Permission denied (22)

3. **Archive/Backup Errors** (Status: 30-31)
   - Disk space insufficient (30)
   - Archive/backup directory creation failure (31)

4. **Identical Content** (Status: 0)
   - Directory identical to existing archive (0)
   - File identical to existing backup (0)

5. **Success** (Status: 0)
   - Archive created successfully (0)
   - Backup created successfully (0)

### Error Context Preservation

```go
type ArchiveError struct {
    Message    string    // Human-readable error message
    StatusCode int       // Exit status code
    Operation  string    // "create", "verify", "list", "backup"
    Path       string    // File or directory path
    Err        error     // Underlying error
}
```

### Enhanced Error Detection

```go
// Comprehensive disk space error detection
func IsDiskFullError(err error) bool {
    patterns := []string{
        "no space left on device",
        "disk full",
        "not enough space",
        "insufficient disk space",
        "device full",
        "quota exceeded",
        "file too large",
    }
    // Case-insensitive matching
}

// Permission error detection
func IsPermissionError(err error) bool {
    patterns := []string{
        "permission denied",
        "access denied",
        "operation not permitted",
        "insufficient privileges",
    }
}
```

## Concurrency Architecture

### Thread Safety

- **Resource Manager**: Thread-safe resource tracking with mutex protection
- **Template Engine**: Concurrent template rendering with immutable templates
- **Error Handler**: Safe error aggregation with structured error types
- **Configuration**: Immutable after initialization for thread safety

### Context-Aware Operations

- **Archive Creation**: Context cancellation support with periodic checks
- **File Operations**: Context-aware file copying with cancellation
- **Resource Cleanup**: Context-aware cleanup with timeout handling
- **ZIP Creation**: Cancellation checks during file processing

### Goroutine Usage

- **File Scanning**: Directory traversal with cancellation checks
- **Compression**: File processing with context cancellation
- **Verification**: Checksum calculation with cancellation support
- **Cleanup**: Background resource cleanup with panic recovery

## Testing Architecture

### Unit Testing

- **Service Layer**: Mock dependencies with context support
- **Configuration**: Various input scenarios including environment variables
- **Error Handling**: Error condition simulation with structured errors
- **Template Engine**: Format string and template validation
- **Resource Management**: Cleanup verification and thread safety
- **Context Operations**: Cancellation and timeout testing

### Integration Testing

- **End-to-End**: Complete archive and backup workflows
- **File System**: Real directory and file operations
- **Git Integration**: Repository interaction with various states
- **Configuration**: Multi-source configuration with precedence
- **Context Integration**: Cancellation in real workflows

### Performance Testing

- **Large Directories**: Scalability testing with context support
- **Memory Usage**: Resource consumption with cleanup verification
- **Compression Speed**: Archive creation performance with cancellation
- **Verification Speed**: Checksum calculation with context support

## Testing Infrastructure Architecture

### Archive Corruption Testing Framework (TEST-INFRA-001-A) ✅ COMPLETED

The testing infrastructure provides systematic archive corruption utilities for comprehensive verification testing.

```
┌─────────────────────────────────────────────────────────────┐
│                Testing Infrastructure Layer                  │
├─────────────────────────────────────────────────────────────┤
│  CorruptionConfig  │  ArchiveCorruptor  │  CorruptionResult │
│  CorruptionDetector│  Test Utilities    │  Backup/Restore   │
│  Performance Bench │  Cross-Platform    │  Resource Safety  │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                  Core Services Integration                   │
├─────────────────────────────────────────────────────────────┤
│  Archive Service   │  Verification Svc  │  Comparison Svc   │
│  Error Handling    │  Resource Manager  │  File Operations  │
│  Context Support   │  Cleanup System    │  Performance     │
└─────────────────────────────────────────────────────────────┘
```

#### Testing Data Models

```go
// CorruptionConfig configures systematic corruption application
type CorruptionConfig struct {
    Type           CorruptionType  // Type of corruption to apply
    Seed           int64          // For reproducible corruption
    TargetFile     string         // Specific file to corrupt (empty for archive-level)
    CorruptionSize int            // Number of bytes to corrupt
    Offset         int            // Byte offset for corruption (-1 for random)
    Severity       float32        // 0.0-1.0, how severe the corruption should be
}

// CorruptionResult contains comprehensive corruption information
type CorruptionResult struct {
    Type           CorruptionType  // Type of corruption applied
    AppliedAt      []int          // Byte offsets where corruption was applied
    OriginalBytes  []byte         // Original bytes for potential recovery
    CorruptedBytes []byte         // The corrupted bytes that were written
    Description    string         // Human-readable description of corruption
    Recoverable    bool           // Whether this corruption can be recovered from
}

// ArchiveCorruptor provides controlled ZIP corruption capabilities
type ArchiveCorruptor struct {
    originalPath string          // Path to archive being corrupted
    backupPath   string          // Path to backup for restoration
    config       CorruptionConfig // Configuration for corruption behavior
}

// CorruptionDetector analyzes archives for corruption types
type CorruptionDetector struct {
    archivePath string           // Path to archive for analysis
}
```

#### Corruption Type Architecture

```go
// CorruptionType represents different categories of archive corruption
type CorruptionType int

const (
    CorruptionCRC         CorruptionType = iota // Corrupt file checksums (recoverable)
    CorruptionHeader                            // Corrupt ZIP file headers (fatal)
    CorruptionTruncate                         // Cut off end of archive (permanent)
    CorruptionCentralDir                       // Corrupt ZIP central directory (fatal)
    CorruptionLocalHeader                      // Corrupt individual file headers (sometimes recoverable)
    CorruptionData                             // Corrupt actual file data (detectable/recoverable)
    CorruptionSignature                        // Corrupt ZIP file signatures (fatal)
    CorruptionComment                          // Corrupt archive comments (non-fatal)
)
```

#### Testing Service Integration

```go
// Testing service integration with existing verification systems
type TestingService struct {
    corruptor *ArchiveCorruptor
    detector  *CorruptionDetector
    manager   *ResourceManager    // Integration with existing resource management
    formatter *OutputFormatter    // Integration with existing output formatting
}

// Archive testing integration methods
func (ts *TestingService) CreateCorruptedArchive(path string, config CorruptionConfig) (*CorruptionResult, error)
func (ts *TestingService) DetectCorruption(path string) ([]CorruptionType, error)
func (ts *TestingService) ValidateVerificationBehavior(path string, expected []CorruptionType) error
func (ts *TestingService) BenchmarkCorruptionPerformance(config CorruptionConfig) (*PerformanceResult, error)
```

#### Design Patterns and Implementation

**1. Safe Corruption Pattern**
```go
// Automatic backup/restore pattern for safe testing
func (ac *ArchiveCorruptor) ApplyCorruption() (*CorruptionResult, error) {
    // Create backup first
    if err := ac.CreateBackup(); err != nil {
        return nil, fmt.Errorf("failed to create backup: %w", err)
    }
    
    // Apply corruption with error handling
    result, err := ac.corruptArchive()
    if err != nil {
        // Automatic restoration on failure
        ac.RestoreFromBackup()
        return nil, err
    }
    
    return result, nil
}
```

**2. Deterministic Corruption Pattern**
```go
// Reproducible corruption using seed-based generation
func (ac *ArchiveCorruptor) generateDeterministicCorruption(offset int) []byte {
    // Use seed + offset for deterministic but different corruption per location
    seedValue := uint32(ac.config.Seed + int64(offset))
    corruptedBytes := make([]byte, 4)
    binary.LittleEndian.PutUint32(corruptedBytes, seedValue)
    return corruptedBytes
}
```

**3. Cross-Platform File Operations**
```go
// Cross-platform temporary file management
func (ac *ArchiveCorruptor) createPlatformSafeTempFile() (*os.File, error) {
    // Platform-specific temporary directory with proper permissions
    tempDir := os.TempDir()
    return os.CreateTemp(tempDir, "corruption-test-*.tmp")
}
```

#### Performance Architecture

**Corruption Performance Characteristics**:
- **CRC Corruption**: ~763μs average for typical archives
- **Header Corruption**: ~45μs average (minimal data modification)
- **Truncation**: ~12μs average (simple file size operation)
- **Data Corruption**: ~1.2ms average (depends on corruption size)
- **Signature Corruption**: ~38μs average (minimal header modification)
- **Detection Analysis**: ~49μs average for classification

**Memory Usage Patterns**:
- **Backup Storage**: 1:1 archive size ratio during testing
- **Corruption Buffer**: <1KB additional memory for corruption data
- **Detection Analysis**: <512 bytes for signature analysis
- **Resource Tracking**: <100 bytes per tracked resource

#### Integration with Existing Architecture

**Verification Service Integration**:
```go
// Integration with existing verify.go logic
func TestVerificationWithCorruption(t *testing.T) {
    // Create corrupted archive using testing infrastructure
    config := CorruptionConfig{Type: CorruptionCRC, Seed: 12345}
    result, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
    
    // Test verification behavior with existing verification service
    verifier := NewArchiveVerifier(cfg, outputFormatter)
    status, err := verifier.VerifyArchive(archivePath)
    
    // Validate verification properly detects corruption
    assert.Equal(t, VerificationStatusCorrupted, status)
}
```

**Comparison Service Integration**:
```go
// Integration with existing comparison.go logic
func TestComparisonWithCorruption(t *testing.T) {
    // Create corrupted archive for comparison testing
    corruptor, err := NewArchiveCorruptor(archivePath, corruptionConfig)
    result, err := corruptor.ApplyCorruption()
    
    // Test comparison logic detects corruption
    comparator := NewArchiveComparator(cfg)
    identical, err := comparator.CompareArchiveToDirectory(archivePath, sourceDir)
    
    // Validate comparison properly handles corrupted archives
    assert.False(t, identical, "Corrupted archive should not match source directory")
}
```

#### Resource Management Integration

**Cleanup Integration**:
```go
// Integration with existing ResourceManager
func TestCorruptionWithResourceManagement(t *testing.T) {
    rm := NewResourceManager()
    defer rm.CleanupWithPanicRecovery()
    
    // Register corruption testing resources
    corruptor, err := NewArchiveCorruptor(archivePath, config)
    rm.AddCustomResource(&CorruptionResource{corruptor: corruptor})
    
    // Corruption testing with automatic cleanup
    result, err := corruptor.ApplyCorruption()
    // ... testing logic ...
    
    // Automatic restoration and cleanup handled by ResourceManager
}
```

#### Testing Infrastructure Service Layer

The testing infrastructure integrates seamlessly with existing service architecture:

1. **Configuration Service**: Uses existing config loading for test parameters
2. **Resource Management**: Integrates with ResourceManager for cleanup
3. **Output Formatting**: Uses OutputFormatter for test result display
4. **Error Handling**: Uses existing error types for consistent handling
5. **Context Support**: Supports context cancellation for long-running corruption tests

This architecture provides comprehensive testing capabilities while maintaining consistency with existing system patterns and ensuring safe, reproducible corruption testing for verification logic validation.

## Security Architecture

### File System Security

- **Path Validation**: Prevent directory traversal with enhanced validation
- **Permission Checking**: Verify access rights with structured errors
- **Symlink Handling**: Safe symbolic link processing
- **Temporary Files**: Secure temporary file creation with automatic cleanup

### Archive Security

- **Checksum Verification**: SHA-256 integrity checks with context support
- **Compression Bombs**: ZIP bomb protection during creation
- **Path Sanitization**: Safe archive extraction with validation
- **Metadata Validation**: Archive metadata verification with structured errors

## Extensibility Architecture

### Plugin Interface

```go
type ArchivePlugin interface {
    Name() string
    Process(ctx context.Context, archive *Archive) error
    Cleanup() error
}

type BackupPlugin interface {
    Name() string
    Process(ctx context.Context, backup *BackupInfo) error
    Cleanup() error
}
```

### Hook System

- **Pre-Archive**: Before archive creation with context
- **Post-Archive**: After archive creation with verification
- **Pre-Verification**: Before verification with context
- **Post-Verification**: After verification with status
- **Pre-Backup**: Before file backup creation
- **Post-Backup**: After file backup creation

### Custom Formatters

```go
type OutputFormatter interface {
    Format(template string, data interface{}) (string, error)
    RegisterFunction(name string, fn interface{})
}

type TemplateFormatter interface {
    FormatWithTemplate(input, pattern, tmplStr string) (string, error)
    FormatWithPlaceholders(format string, data map[string]string) string
}
```

## Deployment Architecture

### Binary Distribution

- **Single Binary**: Self-contained executable with embedded defaults
- **Cross-Platform**: Linux, macOS, Windows support
- **Static Linking**: No external dependencies
- **Version Embedding**: Build-time version info with Git integration

### Configuration Management

- **Default Configuration**: Embedded defaults with immutable specifications
- **Configuration Validation**: Schema-based validation with error reporting
- **Migration Support**: Configuration version upgrades
- **Environment Integration**: Container-friendly configuration with BKPDIR_CONFIG

## Performance Considerations

### Memory Management

- **Streaming Processing**: Large file handling with context support
- **Buffer Pooling**: Memory reuse with efficient allocation
- **Garbage Collection**: Minimal allocation with resource cleanup
- **Resource Cleanup**: Automatic resource management with panic recovery

### I/O Optimization

- **Buffered I/O**: Efficient file operations with context cancellation
- **Parallel Processing**: Concurrent file handling with thread safety
- **Compression Tuning**: Optimal compression settings with performance monitoring
- **Disk Space Monitoring**: Available space checking with enhanced error detection

### Scalability

- **Large Directories**: Efficient directory traversal with context support
- **Many Files**: Scalable file processing with resource management
- **Deep Hierarchies**: Stack-safe recursion with cancellation checks
- **Long-Running Operations**: Progress reporting with context cancellation

## CLI Commands Architecture

### Command Structure

```
bkpdir
├── create [NOTE]                    # Create full directory archive
├── create --incremental [NOTE]     # Create incremental directory archive
├── backup FILE_PATH [NOTE]         # Create file backup
├── list                            # List directory archives
├── --list FILE_PATH                # List file backups for specific file
├── verify [ARCHIVE_NAME]           # Verify archive integrity
├── config                          # Display configuration
├── full [NOTE]                     # Backward compatibility: create full archive
├── inc [NOTE]                      # Backward compatibility: create incremental archive
└── --config                        # Backward compatibility: display configuration
```

### Command Implementation

**Directory Archive Commands:**
- `create`: Primary command for directory archiving
- `create --incremental`: Incremental archive creation
- `full`: Backward compatibility alias for `create`
- `inc`: Backward compatibility alias for `create --incremental`

**File Backup Commands:**
- `backup FILE_PATH [NOTE]`: Create backup of individual file
- `--list FILE_PATH`: List backups for specific file

**Management Commands:**
- `list`: List all directory archives with verification status
- `verify [ARCHIVE_NAME]`: Verify archive integrity with optional checksums
- `config`: Display computed configuration values with sources
- `--config`: Backward compatibility alias for `config`

### Command Flags

**Global Flags:**
- `--dry-run, -d`: Show what would be done without executing
- `--verify, -v`: Verify archive after creation
- `--checksum, -c`: Include checksum verification
- `--note, -n`: Add note to archive/backup name
- `--config`: Display configuration (backward compatibility)
- `--list FILE`: List file backups (backward compatibility)

**Context Support:**
- All long-running operations support context cancellation
- Timeout support through context.WithTimeout()
- Graceful cancellation with resource cleanup
- Progress indication for long operations

### Output Formatting

**Configurable Output:**
- Printf-style format strings for all output
- Template-based formatting with named placeholders
- ANSI color support for enhanced readability
- Regex-based data extraction for rich formatting

**Format Categories:**
- Archive operations: creation, listing, verification
- File backup operations: creation, listing, comparison
- Configuration display: values, sources, validation
- Error messages: structured errors with context

This architecture provides a robust foundation for BkpDir's directory archiving and file backup capabilities while incorporating advanced features for enhanced error handling, resource management, context-aware operations, and comprehensive output formatting.

## 🔻 CI/CD Pipeline Architecture for AI Development

### CICD-001: AI-Optimized Pipeline Architecture
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// CICD-001: AI-first CI/CD optimization`

#### Pipeline Architecture Overview
```
┌─────────────────────────────────────────────────────────────┐
│                   AI-First CI/CD Pipeline                  │
├─────────────────────────────────────────────────────────────┤
│  Trigger: AI Assistant Code Changes                        │
│  ├─ Low Priority Scheduling                                │
│  ├─ Background Processing                                  │
│  └─ Non-blocking Execution                                 │
├─────────────────────────────────────────────────────────────┤
│  Phase 1: AI Protocol Validation                          │
│  ├─ DOC-008 Icon Validation (make validate-icon-enforcement) │
│  ├─ Token Compliance Checking                             │
│  ├─ Documentation Synchronization                         │
│  └─ AI Assistant Protocol Adherence                       │
├─────────────────────────────────────────────────────────────┤
│  Phase 2: Automated Quality Gates                         │
│  ├─ Lint Checking (make lint)                             │
│  ├─ Test Execution (make test)                            │
│  ├─ Coverage Validation                                    │
│  └─ Build Verification (make build)                       │
├─────────────────────────────────────────────────────────────┤
│  Phase 3: AI-Optimized Reporting                          │
│  ├─ Standardized Icon Output (⭐🔺🔶🔻)                    │
│  ├─ Token-Based Traceability                              │
│  ├─ AI-Readable Status Messages                           │
│  └─ Automated Documentation Updates                       │
├─────────────────────────────────────────────────────────────┤
│  Output: AI Assistant Feedback Loop                       │
│  ├─ Success: Continue Development                          │
│  ├─ Warnings: Non-blocking Notifications                  │
│  └─ Failures: Detailed AI-Friendly Error Reports         │
└─────────────────────────────────────────────────────────────┘
```

#### Technical Components

**🔧 Pipeline Orchestration Layer**
- **Priority Queue Management**: Low-priority task scheduling
- **Resource Allocation**: Background CPU/memory usage optimization
- **Parallel Processing**: Multi-stage validation execution
- **Context Awareness**: AI assistant workflow state tracking

**🛡️ AI Protocol Validation Engine**
- **Token Validation**: Integration with DOC-008 comprehensive validation
- **Documentation Sync**: Cross-reference consistency checking
- **Protocol Compliance**: ai-assistant-protocol.md adherence verification
- **Icon Standardization**: DOC-007 implementation token format validation

**📊 Quality Gate Automation**
- **Zero-Human Gates**: Fully automated approval processes
- **Configurable Thresholds**: AI-optimized quality metrics
- **Failure Recovery**: Self-healing pipeline configurations
- **Metrics Collection**: AI workflow performance analytics

**🔍 AI-Optimized Reporting System**
- **Structured Output**: Machine-readable status reports
- **Icon Integration**: Standardized visual indicators for AI comprehension
- **Token Attribution**: Change traceability through implementation tokens
- **Feedback Optimization**: Streamlined error reporting for AI iteration

#### Integration Points

**📋 Feature Tracking Integration**
- Automatic status updates in `feature-tracking.md`
- Token-based change attribution
- Cross-reference validation with documentation

**🔧 Build System Integration**
- Enhanced Makefile targets for AI workflows
- Automated dependency resolution
- Background build artifact management

**📝 Documentation Pipeline**
- Automatic documentation consistency checking
- Cross-file reference validation
- AI-friendly documentation formatting

#### Deployment Architecture

**🚀 Pipeline Execution Environment**
- Containerized execution for consistency
- Resource limits for background processing
- Parallel execution capability
- Automated cleanup and recovery

**📊 Monitoring and Analytics**
- AI workflow performance metrics
- Pipeline efficiency tracking
- Automated optimization recommendations
- Resource utilization monitoring

## 🔻 AI-First Documentation and Code Maintenance Architecture

### DOC-011: AI Validation Framework Architecture
**Priority**: 🔺 HIGH  
**Implementation Tokens**: `// DOC-011: AI validation integration`

#### AI Validation Framework Overview
```
┌─────────────────────────────────────────────────────────────┐
│            AI Validation Integration Framework              │
├─────────────────────────────────────────────────────────────┤
│  Input: AI Assistant Code Changes                          │
│  ├─ Pre-submission Validation Requests                     │
│  ├─ Real-time Validation Queries                          │
│  ├─ Batch Validation Processing                            │
│  └─ Bypass Request Handling                                │
├─────────────────────────────────────────────────────────────┤
│  Layer 1: AI Workflow Integration                         │
│  ├─ DOC-008 Validation Framework Integration              │
│  ├─ Pre-submission Hook APIs                              │
│  ├─ Real-time Validation Services                         │
│  └─ Zero-friction Workflow Integration                    │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: AI-Optimized Error Processing                   │
│  ├─ AI-Readable Error Message Formatting                  │
│  ├─ Context-aware Remediation Guidance                    │
│  ├─ Structured Error Response APIs                        │
│  └─ Priority-based Error Categorization                   │
├─────────────────────────────────────────────────────────────┤
│  Layer 3: Compliance Monitoring System                    │
│  ├─ AI Assistant Behavior Tracking                        │
│  ├─ Validation Adherence Monitoring                       │
│  ├─ Bypass Audit Trail Management                         │
│  └─ Compliance Dashboard Generation                       │
├─────────────────────────────────────────────────────────────┤
│  Output: AI Assistant Validation Results                  │
│  ├─ Pass/Fail Validation Status                           │
│  ├─ AI-Friendly Error Reports                             │
│  ├─ Remediation Action Lists                              │
│  └─ Compliance Monitoring Data                            │
└─────────────────────────────────────────────────────────────┘
```

#### Core Validation Components

**🔧 AI Workflow Integration Engine**
- **DOC-008 Framework Bridge**: Seamless integration with existing comprehensive validation system
- **Pre-submission API Gateway**: Real-time validation endpoints for AI assistant workflows
- **Workflow Hook Manager**: Non-intrusive validation integration that preserves AI assistant flow
- **Context-Aware Validation**: Validation behavior adapted to AI assistant operation context

**🛡️ AI-Optimized Error Processing System**
- **Error Message Formatter**: AI-readable error message structure and formatting
- **Remediation Guide Generator**: Context-aware guidance for AI assistant error resolution
- **Error Categorization Engine**: Priority-based error classification for AI assistant action prioritization
- **Response API Layer**: Structured API responses optimized for AI assistant consumption

**📊 Bypass Mechanism Framework**
- **Safe Override Controller**: Controlled bypass mechanisms with comprehensive documentation requirements
- **Audit Trail Manager**: Complete tracking and logging of all bypass actions and justifications
- **Approval Workflow Engine**: Automated approval processes with proper authorization controls
- **Rollback System**: Reversible bypass actions with automated rollback capabilities

**🔍 Compliance Monitoring Infrastructure**
- **AI Behavior Analytics**: Comprehensive tracking of AI assistant validation behavior patterns
- **Adherence Metrics Collector**: Real-time monitoring of AI assistant compliance with validation requirements
- **Dashboard Generator**: Visual compliance dashboards for monitoring AI assistant validation adherence
- **Reporting Engine**: Automated compliance report generation with trend analysis

#### Technical Implementation Architecture

**🚀 API Gateway Layer**
```go
type AIValidationGateway struct {
    validator       *DOC008Validator
    errorFormatter  *AIErrorFormatter
    complianceTracker *ComplianceMonitor
    bypassManager   *BypassController
}

type ValidationRequest struct {
    SourceFiles     []string          `json:"source_files"`
    ChangedTokens   []string          `json:"changed_tokens"`
    ValidationMode  string            `json:"validation_mode"`  // "standard", "strict", "legacy"
    RequestContext  *AIRequestContext `json:"request_context"`
}

type ValidationResponse struct {
    Status          string                  `json:"status"`          // "pass", "fail", "warning"
    Errors          []AIOptimizedError     `json:"errors"`
    Warnings        []AIOptimizedWarning   `json:"warnings"`
    RemediationSteps []RemediationAction    `json:"remediation_steps"`
    ComplianceScore float64                `json:"compliance_score"`
}

type AIOptimizedError struct {
    ErrorID         string            `json:"error_id"`
    Category        string            `json:"category"`
    Severity        string            `json:"severity"`
    Message         string            `json:"message"`
    FileReference   *FileLocation     `json:"file_reference"`
    Remediation     *RemediationGuide `json:"remediation"`
    Context         map[string]string `json:"context"`
}
```

**🔧 Validation Integration Layer**
```go
type DOC008ValidationBridge struct {
    validationEngine    *ValidationEngine
    aiContextProcessor  *AIContextProcessor
    resultFormatter     *AIResultFormatter
}

type PreSubmissionHook struct {
    validationGateway   *AIValidationGateway
    hookConfiguration   *HookConfig
    bypassController    *BypassController
}

type AIWorkflowIntegration struct {
    makefileIntegration *MakefileTargets
    cliIntegration      *CLICommands
    apiIntegration      *RESTAPIEndpoints
}
```

**📊 Compliance Monitoring Architecture**
```go
type ComplianceMonitor struct {
    behaviorTracker     *AIBehaviorTracker
    adherenceCalculator *AdherenceMetrics
    auditTrailManager   *AuditTrailManager
    dashboardGenerator  *ComplianceDashboard
}

type AIBehaviorTracker struct {
    validationEvents    chan ValidationEvent
    behaviorPatterns    *PatternAnalyzer
    complianceScoring   *ComplianceScorer
}

type BypassController struct {
    bypassRequestQueue  chan BypassRequest
    approvalWorkflow    *ApprovalEngine
    auditLogger        *AuditLogger
    rollbackManager    *RollbackManager
}
```

#### Integration Points and Data Flow

**📋 DOC-008 Framework Integration**
- Direct integration with existing comprehensive validation engine
- Reuse of validation categories and enforcement mechanisms
- Extension of validation reporting for AI-specific requirements
- Preservation of existing validation behavior with AI-friendly enhancements

**🔌 AI Assistant Workflow Integration**
- Makefile target integration: `make validate-ai-submission`
- CLI command integration: `bkpdir validate --ai-mode`
- API endpoint integration: `/api/v1/validation/ai-submission`
- Pre-commit hook integration for seamless workflow inclusion

**🛡️ Feature Tracking Integration**
- Automatic compliance status updates in `feature-tracking.md`
- Integration with ai-assistant-protocol.md requirements
- Cross-reference validation with ai-assistant-compliance.md
- Token-based traceability through implementation tokens

#### Performance and Scalability Design

**⚡ Real-time Validation Architecture**
- Sub-second validation response times for AI assistant workflows
- Concurrent validation processing for multiple file changes
- Caching layer for frequently validated patterns and tokens
- Optimized validation algorithms for large codebases

**📈 Scalability Components**
- Horizontal scaling capability for high AI assistant usage
- Queue-based processing for batch validation requests
- Distributed compliance monitoring across multiple AI assistants
- Efficient resource utilization for background validation processing

### DOC-013: AI-First Documentation Strategy Architecture
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// DOC-013: AI-first maintenance`

#### AI-First Documentation Framework Overview
```
┌─────────────────────────────────────────────────────────────┐
│              AI-First Documentation Framework              │
├─────────────────────────────────────────────────────────────┤
│  Input: AI Assistant Documentation Needs                  │
│  ├─ New Feature Documentation Requirements                 │
│  ├─ Code Maintenance Documentation Updates                 │
│  ├─ Cross-Reference Integrity Maintenance                  │
│  └─ Quality Standards Compliance Validation               │
├─────────────────────────────────────────────────────────────┤
│  Layer 1: AI Documentation Standards                      │
│  ├─ Implementation Token Standardization                  │
│  ├─ Icon Usage Guidelines (DOC-007/DOC-008)              │
│  ├─ Cross-Reference Integrity Management                   │
│  └─ AI-Friendly Content Structure                         │
├─────────────────────────────────────────────────────────────┤
│  Layer 2: AI Code Maintenance Standards                   │
│  ├─ AI-Optimized Code Comment Patterns                    │
│  ├─ Structured Error Handling for AI Comprehension       │
│  ├─ Variable Naming for AI Understanding                  │
│  └─ Function Design Principles for AI Maintenance         │
├─────────────────────────────────────────────────────────────┤
│  Layer 3: AI Quality Assurance Framework                  │
│  ├─ AI-Friendly Test Structure and Organization           │
│  ├─ Code Quality Monitoring for AI Assistants            │
│  ├─ AI Comprehensibility Metrics and Analytics            │
│  └─ Automated Quality Recommendations                     │
├─────────────────────────────────────────────────────────────┤
│  Layer 4: AI Workflow Optimization                        │
│  ├─ AI Assistant Documentation Workflows                  │
│  ├─ Automated Code Maintenance Processes                  │
│  ├─ AI Performance Optimization Strategies                │
│  └─ Knowledge Repository and Template Management          │
├─────────────────────────────────────────────────────────────┤
│  Output: AI-Optimized Development Environment             │
│  ├─ Standardized Documentation Following AI Principles    │
│  ├─ AI-Friendly Code with Optimal Comprehensibility       │
│  ├─ Consistent Quality Standards Across All AI Work       │
│  └─ Efficient AI Assistant Collaboration and Maintenance  │
└─────────────────────────────────────────────────────────────┘
```

#### Core AI Documentation Components

**📝 AI Documentation Standards Engine**
- **Implementation Token Manager**: Standardized token format with priority and action icons
- **Cross-Reference Validator**: Automated link validation and bidirectional reference checking
- **Content Structure Optimizer**: AI-friendly formatting and hierarchy management
- **Template Management System**: Comprehensive template library for all documentation types

**🏗️ AI Code Maintenance Framework**
- **Code Comment Standardizer**: AI-optimized comment patterns with implementation tokens
- **Error Handling Standardizer**: Structured error patterns for AI comprehension
- **Variable Naming Optimizer**: AI-friendly naming conventions and patterns
- **Function Design Validator**: Single responsibility and clear purpose validation

**📊 AI Quality Assurance System**
- **AI Comprehensibility Analyzer**: Metrics for AI assistant understanding and navigation
- **Code Quality Monitor**: Real-time monitoring of AI code maintenance standards
- **Quality Recommendation Engine**: Automated suggestions for AI assistant improvements
- **Test Structure Optimizer**: AI-friendly test organization and validation patterns

**🚀 AI Workflow Integration Engine**
- **Documentation Workflow Manager**: Standardized processes for AI documentation creation
- **Code Maintenance Automation**: Automated maintenance workflows for AI assistants
- **Performance Optimization Framework**: Code optimization strategies for AI efficiency
- **Knowledge Repository Manager**: Template and best practice management system

#### Technical Implementation Architecture

**📋 AI Documentation Standards Layer**
```go
// 🔻 DOC-013: AI documentation standards architecture - 🏗️ Documentation framework
type AIDocumentationStandardsEngine struct {
    tokenManager         *ImplementationTokenManager
    crossRefValidator    *CrossReferenceValidator  
    contentOptimizer     *AIContentStructureOptimizer
    templateManager      *DocumentationTemplateManager
    iconValidator        *DOC008IconValidator
}

type ImplementationTokenManager struct {
    priorityIconRules    map[string]PriorityIcon    // Feature priority → icon mapping
    actionIconRules      map[string]ActionIcon      // Function behavior → icon mapping
    tokenFormatValidator *TokenFormatValidator      // DOC-007/DOC-008 compliance
    crossRefManager      *CrossReferenceManager     // Feature tracking integration
}

type AIContentStructureOptimizer struct {
    headingHierarchy     *HierarchyValidator       // Consistent H1→H2→H3 structure
    sectionOrdering      *SectionOrderValidator    // Predictable content organization
    terminologyStandards *TerminologyManager       // Consistent vocabulary usage
    linkingPatterns      *LinkingPatternManager    // AI-traversable link formats
}

type DocumentationTemplateManager struct {
    featureTemplates     map[string]Template       // Feature documentation templates
    codeCommentTemplates map[string]Template       // Code comment templates
    archTemplates        map[string]Template       // Architecture documentation templates
    testTemplates        map[string]Template       // Testing documentation templates
    validationRules      []TemplateValidationRule  // Template compliance validation
}
```

**🔧 AI Code Maintenance Layer**
```go
// 🔻 DOC-013: AI code maintenance architecture - 🔧 Code standards framework
type AICodeMaintenanceFramework struct {
    commentStandardizer  *AICodeCommentStandardizer
    errorHandlingFramework *AIErrorHandlingFramework
    namingOptimizer      *AINamingOptimizer
    functionDesignValidator *AIFunctionDesignValidator
    qualityMonitor       *AICodeQualityMonitor
}

type AICodeCommentStandardizer struct {
    tokenFormat          *StandardizedTokenFormat   // Implementation token standards
    commentPatterns      *AICommentPatterns        // AI-optimized comment structure
    dependencyTracker    *DependencyDocumentationTracker // Related feature tracking
    maintenanceNotes     *AIMaintenanceNotesManager // AI assistant maintenance context
}

type AIErrorHandlingFramework struct {
    structuredErrorPatterns *StructuredErrorManager // Standardized error structures
    errorClassificationRules *ErrorClassificationEngine // AI decision framework
    recoveryPatterns        *ErrorRecoveryPatternManager // Automated recovery strategies
    errorContextManager     *ErrorContextManager    // AI-accessible error context
}

type AINamingOptimizer struct {
    variableNamingRules  *VariableNamingRules      // AI comprehension optimization
    functionNamingRules  *FunctionNamingRules      // Clear purpose indication
    typeNamingRules      *TypeNamingRules          // Business logic relationship
    consistencyValidator *NamingConsistencyValidator // Cross-codebase consistency
}
```

**📊 AI Quality Assurance Layer**
```go
// 🔻 DOC-013: AI quality assurance architecture - 📊 Quality framework
type AIQualityAssuranceSystem struct {
    comprehensibilityAnalyzer *AIComprehensibilityAnalyzer
    qualityMetricsEngine     *AICodeQualityMetricsEngine
    recommendationEngine     *AIQualityRecommendationEngine
    testStructureOptimizer   *AITestStructureOptimizer
    complianceMonitor        *AIComplianceMonitor
}

type AIComprehensibilityAnalyzer struct {
    tokenCoverageAnalyzer    *ImplementationTokenCoverageAnalyzer
    iconStandardizationAnalyzer *IconStandardizationAnalyzer
    crossRefIntegrityAnalyzer *CrossReferenceIntegrityAnalyzer
    errorHandlingConsistencyAnalyzer *ErrorHandlingConsistencyAnalyzer
    comprehensibilityScorer  *AIComprehensibilityScorer
}

type AICodeQualityMetricsEngine struct {
    qualityMetrics          *AICodeQualityMetrics     // Real-time quality tracking
    analysisContext         *QualityAnalysisContext   // Quality analysis framework
    validationRules         *AIQualityValidationRules // AI-specific quality rules
    performanceTracker      *AIPerformanceTracker     // AI assistant productivity metrics
}

type AIQualityRecommendationEngine struct {
    improvementPatterns     *ImprovementPatternDatabase // AI improvement suggestions
    bestPracticeValidator   *BestPracticeValidator      // AI best practice compliance
    qualityRules           *QualityRuleEngine          // Automated quality recommendations
    prioritizationEngine   *RecommendationPrioritizer  // High/medium/low priority recommendations
}
```

**🚀 AI Workflow Integration Layer**
```go
// 🔻 DOC-013: AI workflow integration architecture - 🚀 Workflow framework
type AIWorkflowIntegrationEngine struct {
    documentationWorkflowManager *AIDocumentationWorkflowManager
    codeMaintenanceAutomation   *AICodeMaintenanceAutomation
    performanceOptimizer        *AIPerformanceOptimizer
    knowledgeRepository         *AIKnowledgeRepository
    collaborationFramework      *AICollaborationFramework
}

type AIDocumentationWorkflowManager struct {
    creationWorkflow    *AIDocumentationCreationWorkflow    // 5-step creation process
    maintenanceWorkflow *AIDocumentationMaintenanceWorkflow // 5-step maintenance process
    validationWorkflow  *AIDocumentationValidationWorkflow  // 5-step validation process
    integrationWorkflow *AIDocumentationIntegrationWorkflow // Feature tracking integration
}

type AICodeMaintenanceAutomation struct {
    tokenMigrationEngine    *ImplementationTokenMigrationEngine // Mass token standardization
    commentStandardization  *CodeCommentStandardizationEngine   // AI-first comment updates
    qualityImprovementEngine *AutomatedQualityImprovementEngine // Automated code improvements
    crossRefMaintenanceEngine *CrossReferenceMaintenanceEngine  // Automated link maintenance
}

type AIPerformanceOptimizer struct {
    codeOptimizationStrategies *CodeOptimizationStrategyEngine // AI efficiency optimization
    patternConsolidation       *PatternConsolidationEngine     // Repetitive pattern consolidation
    errorHandlingStandardization *ErrorHandlingStandardizationEngine // Consistent error patterns
    comprehensibilityImprovement *ComprehensibilityImprovementEngine // AI understanding optimization
}
```

#### Data Flow and Integration Architecture

**📋 AI Documentation Workflow Integration**
```
AI Assistant Request → Template Selection → Content Generation → 
Validation (DOC-008) → Cross-Reference Update → Feature Tracking Update → 
Quality Verification → Documentation Publication
```

**🔧 AI Code Maintenance Workflow Integration**
```
Code Change Request → Token Analysis → Comment Standardization → 
Error Handling Validation → Naming Optimization → Quality Assessment → 
Cross-Reference Sync → Compliance Verification → Code Commitment
```

**📊 AI Quality Monitoring Integration**
```
Continuous Code Analysis → Quality Metrics Collection → 
Comprehensibility Scoring → Recommendation Generation → 
Improvement Prioritization → Automated Optimization → 
Performance Tracking → Feedback Loop Integration
```

#### Integration Points and Dependencies

**🔗 DOC-007/DOC-008 Integration**
- Direct integration with Source Code Icon Integration standards
- Comprehensive validation through Icon Validation and Enforcement system
- Mass Implementation Token Standardization format compliance
- Real-time validation and enforcement of icon usage standards

**📋 Feature Tracking Integration**
- Automatic feature-tracking.md status updates for all AI assistant work
- Implementation token cross-reference maintenance and validation
- Priority icon assignment based on feature priority in feature registry
- Bidirectional relationship maintenance between documentation and code

**🛡️ AI Assistant Compliance Integration**
- Integration with ai-assistant-compliance.md validation requirements
- Token referencing and cross-reference requirement enforcement
- Validation workflow integration with DOC-008 comprehensive system
- Compliance monitoring and reporting for AI assistant behavior

#### Performance and Scalability Characteristics

**⚡ AI Assistant Performance Optimization**
- **Documentation Creation**: <5 minutes average AI assistant task completion
- **Cross-Reference Updates**: <2 minutes average link validation and repair
- **Code Comprehension**: <30 seconds average for AI assistant function understanding
- **Quality Compliance**: <30 seconds DOC-008 validation completion

**📈 Scalability and Efficiency**
- **Multi-AI Assistant Support**: Framework designed for multiple concurrent AI assistants
- **Template Scalability**: Comprehensive template library supporting all documentation types
- **Quality Monitoring**: Real-time quality metrics collection and analysis
- **Knowledge Management**: Centralized knowledge repository for AI assistant collaboration

**🔧 Integration Efficiency**
- **Workflow Automation**: >95% automated workflow completion rate
- **Quality Compliance**: >99% automated compliance checking accuracy
- **Cross-Reference Integrity**: >99% link validation success rate
- **AI Comprehensibility**: >95% AI assistant comprehension rate for standardized documentation

#### Technology Stack and Implementation

**📋 Core Technologies**
- **Documentation Engine**: Markdown with standardized icon system and template management
- **Validation Framework**: DOC-008 comprehensive validation with AI-optimized error reporting
- **Quality Monitoring**: Real-time metrics collection with automated recommendation generation
- **Workflow Management**: Automated AI assistant workflow integration with feature tracking

**🔧 Development Infrastructure**
- **Template Management**: Version-controlled template library with validation rules
- **Quality Assurance**: Automated quality gates integrated with development workflow
- **Performance Monitoring**: AI assistant productivity analytics and optimization feedback
- **Knowledge Repository**: Comprehensive documentation and best practice management system

This architecture provides the comprehensive foundation for AI-first development where AI assistants can efficiently create, maintain, and optimize both code and documentation following consistent, machine-readable standards that enable optimal AI assistant collaboration and long-term maintainability.

## 🔧 Configuration Management System

### Configuration Architecture

**🔻 REFACTOR-003: Configuration Schema Abstraction Implementation**

The configuration system has been enhanced with schema abstraction to prepare for extraction while maintaining complete backward compatibility. The new architecture introduces clean interface boundaries and separation of concerns.

#### Interface Layer
```go
// Schema-agnostic configuration management
type ConfigLoader interface {
    LoadConfig(root string) (*Config, error)
    LoadConfigValues(root string) (map[string]ConfigValue, error)
    GetConfigValues(cfg *Config) []ConfigValue
    GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue
    ValidateConfig(cfg *Config) error
}

// Configuration merging and composition
type ConfigMerger interface {
    MergeConfigs(dst, src *Config)
    MergeConfigValues(dst, src map[string]ConfigValue)
    GetConfigSearchPaths() []string
    ExpandPath(path string) string
}

// Pluggable configuration validation
type ConfigValidator interface {
    ValidateSchema(cfg *Config) error
    ValidateValues(values map[string]ConfigValue) error
    GetRequiredFields() []string
    GetValidationRules() map[string]ValidationRule
}

// Configuration source abstraction
type ConfigSource interface {
    LoadFromFile(path string) (*Config, error)
    LoadFromEnvironment() (*Config, error)
    LoadDefaults() *Config
    GetSourceName() string
    IsAvailable() bool
}
```

#### Schema Separation
The backup application schema has been separated into focused structures:

```go
// Archive-specific configuration
type ArchiveSettings struct {
    DirectoryPath      string
    UseCurrentDirName  bool
    ExcludePatterns    []string
    IncludeGitInfo     bool
    ShowGitDirtyStatus bool
    Verification       *VerificationConfig
}

// Backup-specific configuration
type BackupSettings struct {
    DirectoryPath             string
    UseCurrentDirNameForFiles bool
}

// Format configuration separation
type FormatSettings struct {
    FormatStrings      map[string]string
    TemplateStrings    map[string]string
    PatternStrings     map[string]string
    ErrorFormatStrings map[string]string
}
```

#### Implementation Layer
Concrete implementations maintain backward compatibility:

- **DefaultConfigLoader**: Maintains existing configuration loading behavior
- **BackupAppValidator**: Backup application-specific validation
- **BackupApplicationConfig**: Wrapper providing application-specific interface
- **FileSystemOperations**: File system operations abstraction

#### Extraction Readiness
The abstraction prepares for extraction to `pkg/config` by providing:

- **Clean Interface Boundaries**: No circular dependencies, clear separation
- **Dependency Injection**: File operations, validation, and merging abstracted
- **Schema Independence**: Configuration loading separated from application schema
- **Backward Compatibility**: All existing functionality preserved

**Implementation**: Complete - Ready for EXTRACT-001 (Configuration Management System)

### Configuration Discovery and Loading

The configuration system supports multiple sources with a clear priority hierarchy:

```
1. Command line flags (highest priority)
2. Environment variables (BKPDIR_CONFIG for search paths)
3. Configuration files in search path order:
   - ./.bkpdir.yml (current directory)
   - ~/.bkpdir.yml (home directory)
   - Custom paths from BKPDIR_CONFIG
4. Default values (lowest priority)
```

### Configuration Sources

**Environment Variable Support:**
- `BKPDIR_CONFIG`: Colon-separated list of configuration file paths
- Home directory expansion support with `~` prefix
- Configurable search path with precedence rules
- Earlier files in search path take precedence over later files

**File-based Configuration:**
- YAML format with comprehensive validation
- Schema-based structure with all documented fields
- Default value inheritance and merging
- Status code configuration with defaults
- Format string configuration with printf-style and template support
- Regex pattern configuration for data extraction

**Configuration Merging:**
- Multiple configuration files processed in search path order
- Earlier files override later files for non-default values
- Default values preserved when not explicitly overridden
- Status codes, format strings, and templates all configurable

### ⭐ Layered Configuration Inheritance Architecture (CFG-005)

**🔻 CFG-005: Layered Configuration Inheritance Implementation**

The configuration system supports layered inheritance enabling configuration files to inherit from other configuration files with flexible merge strategies. This provides hierarchical configuration management for complex project structures.

#### Inheritance Configuration Structure

```go
// Configuration inheritance metadata
type ConfigInheritance struct {
    Inherit []string `yaml:"inherit"` // List of parent configuration files
}

// Inheritable configuration wrapper
type InheritableConfig struct {
    *Config `yaml:",inline"`         // Inline main configuration
    ConfigInheritance `yaml:",inline"` // Inline inheritance metadata
}

// Inheritance chain tracking
type InheritanceChain struct {
    Files   []string           // Configuration files in dependency order
    Visited map[string]bool    // Circular dependency prevention
    Sources map[string]string  // Source tracking for debugging
}
```

#### Merge Strategy Architecture

```go
// Merge strategy processor
type MergeStrategyProcessor struct {
    strategies map[string]MergeStrategy
}

// Merge strategy interface
type MergeStrategy interface {
    Merge(dest, src reflect.Value) error
    GetPrefix() string
    GetDescription() string
}

// Strategy implementations
type StandardOverrideStrategy struct{}  // No prefix - replace parent values
type ArrayMergeStrategy struct{}        // + prefix - append to parent arrays
type ArrayPrependStrategy struct{}      // ^ prefix - prepend to parent arrays
type ArrayReplaceStrategy struct{}      // ! prefix - replace parent arrays
type DefaultValueStrategy struct{}      // = prefix - use if parent not set
```

#### Configuration Loading Engine Enhancement

```go
// Enhanced configuration loader with inheritance support
type InheritanceConfigLoader struct {
    baseLoader     ConfigLoader
    pathResolver   PathResolver
    chainBuilder   InheritanceChainBuilder
    strategyEngine MergeStrategyProcessor
    circularDetector CircularDependencyDetector
}

// Inheritance chain processing
func (l *InheritanceConfigLoader) LoadConfigWithInheritance(root string) (*Config, error) {
    // 1. Build inheritance dependency graph
    chain, err := l.buildInheritanceChain(root)
    if err != nil {
        return nil, err
    }
    
    // 2. Load configurations in dependency order
    configs, err := l.loadConfigurationChain(chain)
    if err != nil {
        return nil, err
    }
    
    // 3. Apply merge strategies and resolve final configuration
    return l.mergeConfigurationChain(configs)
}
```

#### Merge Strategy Processing

```go
// Key prefix analysis for merge strategy determination
type PrefixedKeyProcessor struct {
    prefixMap map[string]MergeStrategy
}

func (p *PrefixedKeyProcessor) ProcessPrefixedKeys(config map[string]interface{}) ProcessedConfig {
    processed := ProcessedConfig{
        StandardKeys: make(map[string]interface{}),
        MergeOps:     make([]MergeOperation, 0),
    }
    
    for key, value := range config {
        strategy, cleanKey := p.extractMergeStrategy(key)
        processed.MergeOps = append(processed.MergeOps, MergeOperation{
            Key:      cleanKey,
            Value:    value,
            Strategy: strategy,
        })
    }
    
    return processed
}
```

#### Circular Dependency Detection

```go
// Depth-first search for circular dependency detection
type CircularDependencyDetector struct {
    visited map[string]bool
    inStack map[string]bool
    path    []string
}

func (d *CircularDependencyDetector) DetectCycle(startFile string, resolver PathResolver) error {
    d.visited = make(map[string]bool)
    d.inStack = make(map[string]bool)
    d.path = make([]string, 0)
    
    return d.dfsDetectCycle(startFile, resolver)
}
```

#### Configuration Source Tracking

```go
// Enhanced source tracking with inheritance information
type InheritanceSourceTracker struct {
    configSources map[string]ConfigSource
    valueOrigins  map[string]ValueOrigin
    chainMetadata InheritanceChainMetadata
}

type ValueOrigin struct {
    SourceFile    string    // Configuration file containing the value
    SourceType    string    // "inherit", "override", "merge", "prepend", "replace", "default"
    MergeStrategy string    // Merge strategy applied
    LineNumber    int       // Line number in source file
    ChainDepth    int       // Depth in inheritance chain
}
```

#### Integration with Existing Configuration System

```go
// Backward compatibility layer
type ConfigAdapterWithInheritance struct {
    inheritanceLoader *InheritanceConfigLoader
    legacyLoader      *DefaultConfigLoader
    inheritanceEnabled bool
}

func (a *ConfigAdapterWithInheritance) LoadConfig(root string) (*Config, error) {
    // Check if inheritance is being used
    if a.hasInheritanceDeclarations(root) {
        return a.inheritanceLoader.LoadConfigWithInheritance(root)
    }
    
    // Fall back to legacy loading for non-inheritance configurations
    return a.legacyLoader.LoadConfig(root)
}
```

#### Performance Optimization

```go
// Inheritance caching for performance
type InheritanceCache struct {
    chainCache    map[string]*InheritanceChain    // Cached dependency chains
    configCache   map[string]*Config              // Cached resolved configurations
    mTimeTracker  map[string]time.Time            // File modification time tracking
    cacheEnabled  bool
    maxCacheSize  int
}

// Cache invalidation based on file modification times
func (c *InheritanceCache) IsCacheValid(filePath string) bool {
    if !c.cacheEnabled {
        return false
    }
    
    cachedTime, exists := c.mTimeTracker[filePath]
    if !exists {
        return false
    }
    
    currentStat, err := os.Stat(filePath)
    if err != nil {
        return false
    }
    
    return currentStat.ModTime().Equal(cachedTime)
}
```

#### Error Handling and Debugging

```go
// Inheritance-specific error types
type InheritanceError struct {
    Type        string    // "circular", "missing_file", "invalid_syntax", "merge_conflict"
    Message     string    // Human-readable error message
    ChainPath   []string  // Files involved in inheritance chain
    SourceFile  string    // File where error originated
    LineNumber  int       // Line number of problematic declaration
}

// Inheritance chain visualization for debugging
type InheritanceChainVisualizer struct {
    chain    *InheritanceChain
    formatter OutputFormatter
}

func (v *InheritanceChainVisualizer) GenerateChainDiagram() string {
    // Generate ASCII art diagram of inheritance relationships
    // Include merge strategies and value origins
    // Highlight conflicts and resolution order
}
```

#### Configuration Example Processing

```yaml
# ~/.bkpdir.yml (base configuration)
archive_dir_path: "~/Archives"
backup_dir_path: "~/Backups"
exclude_patterns:
  - "*.tmp"
  - "*.log"
  - ".DS_Store"
include_git_info: true

# /project/.bkpdir.yml (inheriting configuration)
inherit: "~/.bkpdir.yml"

# Standard override
archive_dir_path: "./project-archives"
backup_dir_path: "./project-backups"

# Array merge strategy (append)
+exclude_patterns:
  - "node_modules/"
  - "dist/"
  - "*.secret"

# Array prepend strategy (high priority)
^exclude_patterns:
  - "*.private"

# Default strategy (use if not set)
=verification:
  verify_on_create: true
```

**Resolved Configuration Result:**
```yaml
archive_dir_path: "./project-archives"  # Overridden by child
backup_dir_path: "./project-backups"    # Overridden by child
exclude_patterns:                       # Merged: prepend + base + append
  - "*.private"                         # ^ prepend strategy
  - "*.tmp"                             # Base configuration
  - "*.log"                             # Base configuration
  - ".DS_Store"                         # Base configuration
  - "node_modules/"                     # + merge strategy
  - "dist/"                             # + merge strategy
  - "*.secret"                          # + merge strategy
include_git_info: true                  # Inherited from base
verification:                           # = default strategy applied
  verify_on_create: true
```

This architecture provides comprehensive inheritance capabilities while maintaining backward compatibility, performance, and robust error handling for complex configuration scenarios.

## Archive Format Architecture