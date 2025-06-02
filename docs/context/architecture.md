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