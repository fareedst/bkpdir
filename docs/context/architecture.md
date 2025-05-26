# BkpDir Architecture

## Overview

BkpDir is a directory archiving tool that creates ZIP-based archives with Git integration, file exclusion patterns, and comprehensive verification capabilities. The architecture incorporates advanced features from BkpFile including structured error handling, resource management, template-based formatting, and enhanced configuration management.

## Core Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                            │
├─────────────────────────────────────────────────────────────┤
│  Command Handlers  │  Flag Processing  │  Output Formatting │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Layer                      │
├─────────────────────────────────────────────────────────────┤
│  Config Discovery  │  Environment Vars │  Status Codes     │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                     Core Services                          │
├─────────────────────────────────────────────────────────────┤
│  Archive Service   │  Git Service      │  Verification     │
│  Resource Manager  │  Template Engine  │  Error Handler    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Storage Layer                           │
├─────────────────────────────────────────────────────────────┤
│  File System      │  ZIP Archives     │  Checksums        │
└─────────────────────────────────────────────────────────────┘
```

## Data Models

### Core Configuration Object

```go
type Config struct {
    // Archive settings
    ArchiveDir      string   `yaml:"archive_dir"`
    Prefix          string   `yaml:"prefix"`
    ExcludePatterns []string `yaml:"exclude_patterns"`
    
    // Git integration
    IncludeGitInfo  bool     `yaml:"include_git_info"`
    
    // Output formatting (from BkpFile)
    FormatStrings   map[string]string `yaml:"format_strings"`
    StatusCodes     map[string]int    `yaml:"status_codes"`
    
    // Enhanced features
    UseTemplates    bool     `yaml:"use_templates"`
    ColorOutput     bool     `yaml:"color_output"`
}
```

### Enhanced Data Objects (from BkpFile)

```go
type ConfigValue struct {
    Name   string
    Value  interface{}
    Source string
}

type ArchiveError struct {
    Operation  string
    Path       string
    Err        error
    StatusCode int
}

type ResourceManager struct {
    resources []Resource
    mutex     sync.RWMutex
}

type TemplateFormatter struct {
    templates map[string]*template.Template
    patterns  map[string]*regexp.Regexp
}
```

## Service Architecture

### Archive Service

**Responsibilities:**
- Directory scanning and file collection
- ZIP archive creation with compression
- Incremental archive logic
- File exclusion pattern matching
- Archive naming with timestamps and Git info

**Key Components:**
- `ArchiveCreator`: Handles full and incremental archives
- `FileScanner`: Traverses directories with exclusion patterns
- `CompressionEngine`: ZIP creation with optimal compression
- `NamingService`: Generates archive names with metadata

### Git Integration Service

**Responsibilities:**
- Git repository detection
- Branch and commit hash extraction
- Git status checking for clean state
- Integration with archive naming

**Key Components:**
- `GitDetector`: Repository discovery
- `BranchExtractor`: Current branch identification
- `HashExtractor`: Commit hash retrieval
- `StatusChecker`: Working directory state

### Resource Management Service (from BkpFile)

**Responsibilities:**
- Automatic cleanup of temporary files
- Thread-safe resource tracking
- Error-resilient cleanup operations
- Context-aware resource management

**Key Components:**
- `ResourceTracker`: Maintains resource registry
- `CleanupManager`: Handles automatic cleanup
- `TempFileManager`: Temporary file lifecycle
- `AtomicOperations`: Safe file operations

### Template Engine (from BkpFile)

**Responsibilities:**
- Template-based output formatting
- Named placeholder resolution
- Regex pattern extraction
- ANSI color support

**Key Components:**
- `TemplateParser`: Go text/template integration
- `PlaceholderResolver`: Named parameter substitution
- `PatternExtractor`: Regex-based data extraction
- `ColorFormatter`: ANSI color code handling

### Error Handling Service (from BkpFile)

**Responsibilities:**
- Structured error creation
- Status code management
- Error context preservation
- Panic recovery

**Key Components:**
- `ErrorBuilder`: Structured error construction
- `StatusManager`: Exit code configuration
- `ContextPreserver`: Error context tracking
- `PanicHandler`: Recovery mechanisms

## Configuration Architecture

### Configuration Discovery

```
1. Command line flags (highest priority)
2. Environment variables (BKPDIR_CONFIG)
3. .bkpdir.yml in current directory
4. .bkpdir.yml in home directory
5. Default values (lowest priority)
```

### Configuration Sources

**File-based Configuration:**
- YAML format with validation
- Schema-based structure
- Default value inheritance
- Environment variable expansion

**Environment Variables:**
- `BKPDIR_CONFIG`: Configuration file path
- `BKPDIR_ARCHIVE_DIR`: Archive directory override
- `BKPDIR_PREFIX`: Archive prefix override

**Command Line Flags:**
- `--config`: Configuration file path
- `--archive-dir`: Archive directory
- `--prefix`: Archive prefix
- `--exclude`: Additional exclusion patterns

## Archive Format Architecture

### Archive Naming Convention

```
Format: [prefix-]timestamp[=branch=hash][=note].zip

Examples:
- backup-20240315-143022.zip
- proj-20240315-143022=main=a1b2c3d.zip
- data-20240315-143022=feature=x9y8z7w=release.zip
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

### Printf-Style Formatting (from BkpFile)

```go
// Format string configuration
FormatStrings: map[string]string{
    "archive_created": "Archive created: %{path} (%{size} bytes)",
    "files_processed": "Processed %{count} files",
    "error_occurred": "\033[31mError: %{message}\033[0m",
}
```

### Template-Based Formatting (from BkpFile)

```go
// Template configuration
Templates: map[string]string{
    "archive_summary": `
Archive Summary:
  Path: {{.Path}}
  Size: {{.Size | humanize}}
  Files: {{.FileCount}}
  {{if .GitInfo}}Git: {{.GitInfo.Branch}}@{{.GitInfo.Commit}}{{end}}
`,
}
```

## Error Handling Architecture

### Error Categories

1. **Configuration Errors** (Status: 1)
   - Invalid configuration files
   - Missing required settings
   - Environment variable issues

2. **File System Errors** (Status: 2)
   - Permission denied
   - Disk space insufficient
   - Path not found

3. **Archive Errors** (Status: 3)
   - Compression failures
   - Corruption detected
   - Verification failures

4. **Git Errors** (Status: 4)
   - Repository not found
   - Git command failures
   - Dirty working directory

### Error Context Preservation

```go
type ArchiveError struct {
    Operation  string    // "create", "verify", "list"
    Path       string    // File or directory path
    Err        error     // Underlying error
    StatusCode int       // Exit status code
    Context    map[string]interface{} // Additional context
}
```

## Concurrency Architecture

### Thread Safety

- **Resource Manager**: Thread-safe resource tracking
- **Template Engine**: Concurrent template rendering
- **Error Handler**: Safe error aggregation
- **Configuration**: Immutable after initialization

### Goroutine Usage

- **File Scanning**: Parallel directory traversal
- **Compression**: Concurrent file processing
- **Verification**: Parallel checksum calculation
- **Cleanup**: Background resource cleanup

## Testing Architecture

### Unit Testing

- **Service Layer**: Mock dependencies
- **Configuration**: Various input scenarios
- **Error Handling**: Error condition simulation
- **Template Engine**: Format string validation

### Integration Testing

- **End-to-End**: Complete archive workflows
- **File System**: Real directory operations
- **Git Integration**: Repository interaction
- **Configuration**: Multi-source configuration

### Performance Testing

- **Large Directories**: Scalability testing
- **Memory Usage**: Resource consumption
- **Compression Speed**: Archive creation performance
- **Verification Speed**: Checksum calculation

## Security Architecture

### File System Security

- **Path Validation**: Prevent directory traversal
- **Permission Checking**: Verify access rights
- **Symlink Handling**: Safe symbolic link processing
- **Temporary Files**: Secure temporary file creation

### Archive Security

- **Checksum Verification**: SHA-256 integrity checks
- **Compression Bombs**: ZIP bomb protection
- **Path Sanitization**: Safe archive extraction
- **Metadata Validation**: Archive metadata verification

## Extensibility Architecture

### Plugin Interface

```go
type ArchivePlugin interface {
    Name() string
    Process(ctx context.Context, archive *Archive) error
    Cleanup() error
}
```

### Hook System

- **Pre-Archive**: Before archive creation
- **Post-Archive**: After archive creation
- **Pre-Verification**: Before verification
- **Post-Verification**: After verification

### Custom Formatters

```go
type OutputFormatter interface {
    Format(template string, data interface{}) (string, error)
    RegisterFunction(name string, fn interface{})
}
```

## Deployment Architecture

### Binary Distribution

- **Single Binary**: Self-contained executable
- **Cross-Platform**: Linux, macOS, Windows
- **Static Linking**: No external dependencies
- **Version Embedding**: Build-time version info

### Configuration Management

- **Default Configuration**: Embedded defaults
- **Configuration Validation**: Schema-based validation
- **Migration Support**: Configuration version upgrades
- **Environment Integration**: Container-friendly configuration

## Performance Considerations

### Memory Management

- **Streaming Processing**: Large file handling
- **Buffer Pooling**: Memory reuse
- **Garbage Collection**: Minimal allocation
- **Resource Cleanup**: Automatic resource management

### I/O Optimization

- **Buffered I/O**: Efficient file operations
- **Parallel Processing**: Concurrent file handling
- **Compression Tuning**: Optimal compression settings
- **Disk Space Monitoring**: Available space checking

### Scalability

- **Large Directories**: Efficient directory traversal
- **Many Files**: Scalable file processing
- **Deep Hierarchies**: Stack-safe recursion
- **Long-Running Operations**: Progress reporting

This architecture provides a robust foundation for BkpDir's directory archiving capabilities while incorporating advanced features from BkpFile for enhanced error handling, resource management, and output formatting. 