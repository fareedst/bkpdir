# Implementation Decisions

**STDD Methodology Version**: 1.0.1

## Overview
This document captures detailed implementation decisions for this project, including specific APIs, data structures, and algorithms. All decisions are cross-referenced with architecture decisions using `[ARCH:*]` tokens and requirements using `[REQ:*]` tokens for traceability.

## Template Structure

When documenting implementation decisions, use this format:

```markdown
## N. Implementation Title [IMPL:IDENTIFIER] [ARCH:RELATED_ARCHITECTURE] [REQ:RELATED_REQUIREMENT]

### Decision: Brief description of the implementation decision
**Rationale:**
- Why this implementation approach was chosen
- What problems it solves
- How it fulfills the architecture decision

### Implementation Approach:
- Specific technical details
- Code structure or patterns
- API design decisions

**Code Markers**: Specific code locations, function names, or patterns to look for

**Cross-References**: [ARCH:RELATED_ARCHITECTURE], [REQ:RELATED_REQUIREMENT]
```

## Notes

- All implementation decisions MUST be recorded here IMMEDIATELY when made
- Each decision MUST include `[IMPL:*]` token and cross-reference both `[ARCH:*]` and `[REQ:*]` tokens
- Implementation decisions are dependent on both architecture decisions and requirements
- DO NOT defer implementation documentation - record decisions as they are made

---


## 1. Configuration Structure [IMPL:CONFIG_STRUCT] [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]

### Config Type
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
    
    // Printf-style format strings
    FormatCreatedArchive   string `yaml:"format_created_archive"`
    FormatIdenticalArchive string `yaml:"format_identical_archive"`
    FormatListArchive      string `yaml:"format_list_archive"`
    FormatCreatedBackup    string `yaml:"format_created_backup"`
    FormatIdenticalBackup  string `yaml:"format_identical_backup"`
    FormatListBackup       string `yaml:"format_list_backup"`
    FormatError            string `yaml:"format_error"`
    
    // Template-based format strings
    TemplateCreatedArchive   string `yaml:"template_created_archive"`
    TemplateIdenticalArchive string `yaml:"template_identical_archive"`
    TemplateListArchive      string `yaml:"template_list_archive"`
    TemplateCreatedBackup    string `yaml:"template_created_backup"`
    TemplateIdenticalBackup  string `yaml:"template_identical_backup"`
    TemplateListBackup       string `yaml:"template_list_backup"`
    TemplateError            string `yaml:"template_error"`
    
    // Regex patterns for data extraction
    PatternArchiveFilename string `yaml:"pattern_archive_filename"`
    PatternBackupFilename  string `yaml:"pattern_backup_filename"`
    PatternConfigLine      string `yaml:"pattern_config_line"`
    PatternTimestamp       string `yaml:"pattern_timestamp"`
}
```

### Default Values
- ArchiveDirPath: "../.bkpdir"
- UseCurrentDirName: true
- ExcludePatterns: [".git/", "vendor/"]
- IncludeGitInfo: false
- BackupDirPath: "../.bkpdir"
- UseCurrentDirNameForFiles: true
- Status codes: Various defaults (0 for success, non-zero for errors)
- Format strings: Default printf-style formats
- Template strings: Default template formats with placeholders

## 2. ZIP Archive Format Implementation [IMPL:ZIP_FORMAT] [ARCH:ARCHIVE_FORMAT] [REQ:ARCHIVE_VERIFICATION]

### Decision: Use ZIP format for all archive operations
**Rationale:**
- Cross-platform compatibility
- Built-in compression
- Wide tool support
- Standard library support in Go

**Implementation Approach:**
- Uses Go's `archive/zip` package
- Atomic file operations with temporary files
- Compression level configurable
- Supports incremental archives

**Code Markers**: `archive/zip` imports, `*.zip` file extensions

## 3. Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING] [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]

### Decision: Support both printf-style and template-based output formatting
**Rationale:**
- Printf for simple, backward-compatible formatting
- Templates for rich data extraction and advanced formatting
- Graceful fallback from template to printf
- Supports ANSI colors and structured output

**Implementation Approach:**
- `OutputFormatter`: Printf-style formatting methods
- `TemplateFormatter`: Advanced template processing
- Regex extraction for rich data formatting
- Support for both Go text/template and placeholder syntax

**Code Markers**: `FormatXXX()` and `TemplateXXX()` functions, regex patterns

## 4. Structured Error Handling [IMPL:STRUCTURED_ERRORS] [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING]

### Decision: Use structured error types with status codes and operation context
**Rationale:**
- Consistent error handling across operations
- Machine-readable status codes for scripting
- Enhanced debugging with operation context
- Supports template-based error formatting

### Error Types
```go
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
```

**Error Wrapping:**
```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

**Error Reporting:**
- Error logging to stderr
- Status code propagation
- Operation context in error messages
- Template formatting support

**Code Markers**: `ArchiveError`, `BackupError` types, `NewXXXError()` functions

## 5. Git Command-line Integration [IMPL:GIT_CLI] [ARCH:GIT_INTEGRATION] [REQ:GIT_INTEGRATION]

### Decision: Use Git command-line interface for repository information
**Rationale:**
- Simplicity over Git library dependencies
- Relies on user's Git installation
- Consistent with user's Git configuration
- Lightweight implementation

**Implementation Approach:**
- Uses `exec.Command("git", ...)` for Git operations
- Repository detection: `git rev-parse --is-inside-work-tree`
- Branch extraction: `git rev-parse --abbrev-ref HEAD`
- Hash extraction: `git rev-parse --short HEAD`
- Status detection: `git status --porcelain`

**Code Markers**: `exec.Command("git", ...)` calls, Git command parsing

## 6. Resource Management with Cleanup [IMPL:RESOURCE_MANAGER] [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]

### Decision: Implement ResourceManager for automatic cleanup with panic recovery
**Rationale:**
- Prevents resource leaks on operation failure
- Handles panic scenarios gracefully
- Thread-safe for concurrent operations
- Simplifies error handling in operations

### ResourceManager Structure
```go
type ResourceManager struct {
    tempFiles []*TempFile
    tempDirs  []*TempDir
    mutex     sync.RWMutex
}

type TempFile struct {
    Path string
}

type TempDir struct {
    Path string
}
```

**Methods:**
- `NewResourceManager()`: Create new resource manager
- `AddTempFile(path string)`: Register temporary file for cleanup
- `AddTempDir(path string)`: Register temporary directory for cleanup
- `RemoveResource(resource Resource)`: Remove resource from tracking
- `Cleanup()`: Clean up all registered resources
- `CleanupWithPanicRecovery()`: Clean up with panic recovery

**Code Markers**: `ResourceManager` type, `CleanupWithPanicRecovery()` calls

## 7. Context-Aware Operations [IMPL:CONTEXT_OPS] [ARCH:CONTEXT_SUPPORT] [REQ:CONTEXT_SUPPORT]

### Decision: Support context cancellation for long-running operations
**Rationale:**
- Enables operation timeouts
- Supports graceful cancellation
- Better user experience for long operations
- Standard Go pattern for cancellable operations

**Implementation Approach:**
- All long-running operations accept `context.Context`
- Periodic cancellation checks during operations
- Context-aware file operations
- Timeout support via `context.WithTimeout()`

**Functions:**
- `CreateArchiveWithContext(ctx context.Context, ...)`
- `CreateFileBackupWithContext(ctx context.Context, ...)`
- `CopyFileWithContext(ctx context.Context, src, dst string)`
- `CompareFilesWithContext(ctx context.Context, file1, file2 string)`

**Code Markers**: `context.Context` parameters, `checkContextCancellation()` calls

## 8. Atomic File Operations [IMPL:ATOMIC_OPS] [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]

### Decision: Use temporary files with atomic rename for all file creation
**Rationale:**
- Prevents corruption on operation failure
- Ensures consistency during concurrent access
- Standard pattern for safe file operations
- Integrates with resource cleanup

**Implementation Approach:**
- Create temporary files with `.tmp` extension
- Write content to temporary file
- Atomic rename to final location
- Register temporary files with ResourceManager

**Code Markers**: `.tmp` file extensions, `os.Rename()` calls

## 9. Testing Implementation [IMPL:TESTING] [ARCH:TESTING_STRATEGY] [REQ:*]

**Note**: This implementation realizes the validation criteria specified in `requirements.md` and follows the testing strategy defined in `architecture-decisions.md`. Each test validates specific satisfaction criteria from requirements.

### Unit Test Structure
```go
func TestExampleFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    InputType
        expected OutputType
    }{
        {
            name:     "test case 1",
            input:    inputValue,
            expected: expectedValue,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := functionUnderTest(tt.input)
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Integration Test Structure
```go
func TestIntegrationScenario(t *testing.T) {
    // Setup
    // Execute
    // Verify
}
```

### Test Infrastructure

**Test Utilities** (`testutil/` package):
- `CreateTestDirectory()`: Creates temporary test directories with files
- `CreateTestFile()`: Creates temporary test files with specified content
- `CreateTestGitRepo()`: Creates temporary Git repositories for testing
- `CreateTestArchive()`: Creates test ZIP archives
- `CreateTestBackup()`: Creates test file backups
- `AssertNoTempFiles()`: Verifies no temporary files remain
- `AssertArchiveContents()`: Validates archive contents
- `AssertBackupContents()`: Validates backup contents
- `MockFileSystem()`: Provides file system mocking capabilities
- `CaptureOutput()`: Captures stdout/stderr for testing
- `SetupTestConfig()`: Creates test configuration files
- `SetupTestConfigForModification()`: Creates test configuration files for modification testing
- `SetupTestContext()`: Creates test contexts with timeouts
- `CleanupTestResources()`: Ensures test cleanup
- `VerifyConfigFileContents()`: Validates configuration file contents and structure
- `CreateTestConfigFile()`: Creates test configuration files with specific values
- `AssertConfigValue()`: Validates configuration values in files
- `AssertConfigFileStructure()`: Validates YAML structure and formatting
- `VerifyResourceCleanup()`: Verifies all resources are cleaned up
- `CreateLargeTestFile()`: Creates large files for performance testing
- `SimulateDiskFull()`: Simulates disk full conditions
- `SimulatePermissionDenied()`: Simulates permission denied conditions

**Testing Frameworks**:
- Archive corruption testing framework
- Disk space simulation framework
- Permission testing framework
- Context cancellation testing helpers
- Error injection framework

**Test Categories**:
- Unit tests: Individual function testing with table-driven tests
- Integration tests: End-to-end workflow testing
- Performance tests: Benchmark testing for various scenarios
- Stress tests: Resource exhaustion and high concurrency testing

**Test Execution**:
- Local development: `make test`, `go test -cover ./...`
- Build tags: `-tags=integration`, `-tags=performance`, `-tags=stress`, `-tags=context`
- Race detection: `go test -race ./...`
- Coverage requirements: Minimum 90% overall, 100% for critical functions

**Quality Gates**:
- Code coverage: Minimum 90% overall, 100% for critical functions
- Performance benchmarks: Must not regress
- Memory usage: Must stay within reasonable bounds, no leaks
- Concurrent safety: No race conditions
- Resource cleanup: 100% cleanup verification
- Platform compatibility: Tests pass on Linux and macOS

## 10. Code Style and Conventions [IMPL:CODE_STYLE]

### Naming
- Use descriptive names
- Follow Go naming conventions
- Exported types/functions: PascalCase
- Unexported: camelCase

### Documentation
- Package-level documentation
- Exported function documentation
- Inline comments for complex logic
- Examples in test files

### Formatting
- Use `go fmt` for code formatting
- Use `revive` linter for code quality
- Follow Go standard formatting

## 11. Package Extraction Implementation [IMPL:PACKAGE_EXTRACTION] [ARCH:PACKAGE_EXTRACTION] [REQ:MAINTAINABILITY]

### Decision: Interface-driven extraction with layered dependency order
**Rationale:**
- Enables clean component boundaries
- Prevents circular dependencies
- Maintains backward compatibility
- Supports gradual migration

### Extraction Order
**Phase 1: Foundation Components (No Dependencies)**
1. `pkg/config` - Configuration management (1,322 lines)
2. `pkg/git` - Git integration utilities (255 lines)
3. `pkg/fileops` - File operations (1,173 lines)

**Phase 2: Interface-Dependent Components**
4. `pkg/formatter` - Output formatting (1,056 lines)
5. `pkg/errors` - Error handling (918 lines)
6. `pkg/resources` - Resource management (486 lines)

**Phase 3: Framework Layer**
7. `pkg/cli` - CLI command framework (722 lines)

**Phase 4: Pattern Layer**
8. `pkg/processing` - Data processing patterns (1,705 lines)

### Interface Contracts
**Before Extraction:**
- Define interfaces for all component dependencies
- Example: `ConfigProvider`, `ConfigMerger` for config package
- Example: `ArchiveConfigInterface`, `ArchiveFormatterInterface` for archive operations

**After Extraction:**
- Components expose clean interfaces
- No internal package dependencies
- External dependencies only (standard library, third-party)

### Backward Compatibility Strategy
**Legacy Type Aliases:**
```go
type FileInfo = fileops.FileInfo
type DirectorySnapshot = fileops.DirectorySnapshot
type PatternMatcher = fileops.PatternMatcher
```

**Wrapper Functions:**
```go
func CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
    return fileops.CreateDirectorySnapshot(rootPath, excludePatterns)
}
```

**Code Markers**: `EXTRACT-*` tokens, interface definitions, wrapper functions

## 12. CLI Framework Implementation [IMPL:CLI_FRAMEWORK] [ARCH:CLI_FRAMEWORK] [REQ:USABILITY]

### Decision: Builder pattern with manager interfaces for CLI construction
**Rationale:**
- Standardizes command construction
- Enables testability through interfaces
- Provides consistent flag management
- Supports context-aware execution

### Core Interfaces
**CommandBuilder:**
```go
type CommandBuilder interface {
    NewCommand(name, short, long string) *cobra.Command
    WithHandler(cmd *cobra.Command, handler func(*cobra.Command, []string) error) *cobra.Command
    WithFlags(cmd *cobra.Command, flags []string) *cobra.Command
    WithSubcommands(parent *cobra.Command, children ...*cobra.Command) *cobra.Command
}
```

**FlagManager:**
```go
type FlagManager interface {
    AddGlobalFlags(cmd *cobra.Command) error
    AddDryRunFlag(cmd *cobra.Command, target *bool) error
    AddNoteFlag(cmd *cobra.Command, target *string) error
    AddConfigFlag(cmd *cobra.Command, target *bool) error
}
```

**ContextManager:**
```go
type ContextManager interface {
    Create(parent context.Context) (context.Context, context.CancelFunc)
    WithTimeout(parent context.Context, timeout string) (context.Context, context.CancelFunc)
    HandleSignals(cancel context.CancelFunc)
}
```

**DryRunManager:**
```go
type DryRunManager interface {
    Execute(ctx CommandContext, op DryRunOperation) error
    Log(ctx CommandContext, message string)
}
```

### Implementation Structures
**DefaultCommandBuilder:**
- Wraps Cobra command creation
- Provides fluent API for command setup
- Integrates with FlagManager

**DefaultFlagManager:**
- Standard flag registration methods
- FlagSet structure for batch flag addition
- Consistent flag naming conventions

**CLIApp:**
- Complete CLI application container
- Manages all managers and builders
- Provides unified API for application setup

**Code Markers**: `pkg/cli/` package, builder methods, manager interfaces

## 13. File Operations Implementation [IMPL:FILE_OPERATIONS] [ARCH:FILE_OPERATIONS] [REQ:RELIABILITY]

### Decision: Atomic operations with comprehensive validation and pattern-based exclusion
**Rationale:**
- Prevents data corruption
- Enhances security
- Supports complex traversal scenarios
- Enables efficient file comparison

### Atomic Operations
**AtomicWriter:**
```go
type AtomicWriter struct {
    filename string
    tempFile *os.File
    perm     os.FileMode
}

func NewAtomicWriter(filename string, perm os.FileMode) (*AtomicWriter, error)
func (aw *AtomicWriter) Write(p []byte) (n int, err error)
func (aw *AtomicWriter) Commit() error
func (aw *AtomicWriter) Rollback() error
```

**Atomic Functions:**
- `AtomicWriteFile(filename string, data []byte, perm os.FileMode) error`
- `AtomicCopy(src, dst string) error`

### Path Validation
**ValidationOptions:**
```go
type ValidationOptions struct {
    MustExist      bool
    MustBeReadable bool
    MustBeWritable bool
    MustBeFile     bool
    MustBeDir      bool
    CheckSecurity  bool
}
```

**Validation Functions:**
- `ValidatePath(path string, options ValidationOptions) error`
- `ValidateExistence(path string) error`
- `ValidateReadable(path string) error`
- `ValidateWritable(path string) error`
- `IsSecurePath(path string) bool`

### Directory Traversal
**TraversalOptions:**
```go
type TraversalOptions struct {
    Exclusions     []string
    FollowSymlinks bool
    MaxDepth       int
    IncludeHidden  bool
    SortOrder      SortOrder
}
```

**Traversal Functions:**
- `Walk(root string, walkFn filepath.WalkFunc) error`
- `WalkWithExclusions(root string, exclusions []string, walkFn filepath.WalkFunc) error`
- `WalkWithOptions(root string, options TraversalOptions, walkFn filepath.WalkFunc) error`

### Pattern Exclusion
**PatternMatcher:**
```go
type PatternMatcher struct {
    patterns []string
}

func NewPatternMatcher(patterns []string) *PatternMatcher
func (pm *PatternMatcher) Match(path string) bool
func (pm *PatternMatcher) AddPattern(pattern string)
```

**Code Markers**: `pkg/fileops/` package, atomic operations, validation functions

## 14. Processing Patterns Implementation [IMPL:PROCESSING_PATTERNS] [ARCH:PROCESSING_PATTERNS] [REQ:PERFORMANCE]

### Decision: Pipeline-based processing with concurrent execution support
**Rationale:**
- Enables reusable processing patterns
- Supports high-performance concurrent operations
- Standardizes naming conventions
- Integrates verification capabilities

### Pipeline Pattern
**Pipeline Interface:**
```go
type PipelineInterface interface {
    AddStage(stage Stage) PipelineInterface
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    GetProgress() Progress
}
```

**Stage Interface:**
```go
type Stage interface {
    Process(ctx context.Context, input interface{}) (interface{}, error)
    Name() string
}
```

### Concurrent Processing
**ConcurrentProcessor:**
```go
type ConcurrentProcessor struct {
    workers    int
    queue      chan WorkItem
    results    chan Result
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    processed  int64
    failed     int64
}

func NewConcurrentProcessor(workers int) *ConcurrentProcessor
func (cp *ConcurrentProcessor) Submit(item WorkItem) error
func (cp *ConcurrentProcessor) Process(ctx context.Context) error
func (cp *ConcurrentProcessor) GetStats() Stats
```

### Naming Conventions
**NamingProvider:**
```go
type NamingProviderInterface interface {
    GenerateName(pattern string, metadata map[string]string) (string, error)
    ParseName(name string, pattern string) (map[string]string, error)
}
```

**Naming Patterns:**
- ISO 8601 timestamp: `2006-01-02T150405`
- Git integration: `{branch}-{hash}-{dirty}`
- Archive pattern: `{prefix}-{timestamp}-{note}.zip`
- Backup pattern: `{name}-{timestamp}-{note}.bak`

### Verification Integration
**VerificationManager:**
```go
type VerificationManager struct {
    providers map[string]VerificationProvider
    defaultAlg string
}

func NewVerificationManager() *VerificationManager
func (vm *VerificationManager) AddProvider(name string, provider VerificationProvider)
func (vm *VerificationManager) Calculate(path string, algorithm string) (string, error)
func (vm *VerificationManager) Verify(path string, checksum string, algorithm string) (bool, error)
```

**Code Markers**: `pkg/processing/` package, pipeline stages, worker pools

## 15. Auto-Detection Implementation [IMPL:AUTO_DETECTION] [ARCH:AUTO_DETECTION] [REQ:USABILITY]

### Decision: Path type detection with command routing logic
**Rationale:**
- Improves user experience
- Reduces command verbosity
- Maintains backward compatibility
- Supports intuitive CLI usage

### Path Type Detection
**Detection Functions:**
```go
func isFile(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return info.Mode().IsRegular()
}

func isDirectory(path string) bool {
    info, err := os.Stat(path)
    if err != nil {
        return false
    }
    return info.IsDir()
}

func validatePath(path string) error {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("path does not exist: %s", path)
        }
        if os.IsPermission(err) {
            return fmt.Errorf("permission denied accessing path: %s", path)
        }
        return fmt.Errorf("error accessing path %s: %v", path, err)
    }
    return nil
}
```

### Command Routing
**Routing Logic:**
```go
func handleAutoDetectedCommand(args []string) {
    path := args[0]
    
    if err := validatePath(path); err != nil {
        // Handle error
    }
    
    if isFile(path) {
        handleAutoDetectedFileBackup(args)
    } else if isDirectory(path) {
        handleAutoDetectedDirectoryArchive(args)
    } else {
        // Handle unsupported types
    }
}
```

**Known Commands Bypass:**
- Explicit commands: `create`, `config`, `list`, `verify`, `backup`, `version`
- Global flags: `--config`, `--dry-run`, `--list`
- Help flags: `--help`, `-h`, `--version`, `-v`

**Integration with Cobra:**
```go
func executeWithAutoDetection(rootCmd *cobra.Command) error {
    args := os.Args[1:]
    
    // Check for known commands/flags
    // If not known, trigger auto-detection
    // Otherwise execute normally
}
```

**Code Markers**: `main.go` path detection functions, `executeWithAutoDetection()`, `handleAutoDetectedCommand()`

## 16. File Statistics Implementation [IMPL:FILE_STATISTICS] [ARCH:FILE_STATISTICS] [REQ:OUTPUT_FORMATTING]

### Decision: FileStatInfo structure with human-readable formatting
**Rationale:**
- Provides rich file information
- Enhances output display
- Supports template formatting
- Improves user experience

### FileStatInfo Structure
```go
type FileStatInfo struct {
    Path      string      // Full file path
    Name      string      // File name only
    Size      int64       // File size in bytes
    SizeHuman string      // Human-readable size (1.2MB, 455KB, etc.)
    MTime     time.Time   // Modification time
    MTimeUnix int64       // Modification time as unix timestamp
    Mode      os.FileMode // File permissions/mode
    Type      string      // File type (regular, directory, symlink)
}
```

### Gathering Function
```go
func GatherFileStatInfo(path string) (*FileStatInfo, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
    }
    
    return &FileStatInfo{
        Path:      path,
        Name:      filepath.Base(path),
        Size:      info.Size(),
        SizeHuman: formatHumanSize(info.Size()),
        MTime:     info.ModTime(),
        MTimeUnix: info.ModTime().Unix(),
        Mode:      info.Mode(),
        Type:      getFileType(info),
    }, nil
}
```

### Human-Readable Size Formatting
```go
func formatHumanSize(size int64) string {
    const (
        KB = 1024
        MB = KB * 1024
        GB = MB * 1024
        TB = GB * 1024
    )
    
    switch {
    case size >= TB:
        return fmt.Sprintf("%.1fTB", float64(size)/TB)
    case size >= GB:
        return fmt.Sprintf("%.1fGB", float64(size)/GB)
    case size >= MB:
        return fmt.Sprintf("%.1fMB", float64(size)/MB)
    case size >= KB:
        return fmt.Sprintf("%.1fKB", float64(size)/KB)
    default:
        return fmt.Sprintf("%dB", size)
    }
}
```

### File Type Detection
```go
func getFileType(info os.FileInfo) string {
    mode := info.Mode()
    switch {
    case mode.IsRegular():
        return "regular"
    case mode.IsDir():
        return "directory"
    case mode&os.ModeSymlink != 0:
        return "symlink"
    case mode&os.ModeDevice != 0:
        return "device"
    case mode&os.ModeNamedPipe != 0:
        return "pipe"
    case mode&os.ModeSocket != 0:
        return "socket"
    default:
        return "other"
    }
}
```

**Code Markers**: `file_stats.go`, `FileStatInfo` struct, `GatherFileStatInfo()` function

## 17. Directory Comparison Implementation [IMPL:DIRECTORY_COMPARISON] [ARCH:DIRECTORY_COMPARISON] [REQ:ARCHIVE_VERIFICATION]

### Decision: Snapshot-based comparison with hash-based content verification
**Rationale:**
- Enables efficient comparison
- Supports archive-to-directory comparison
- Provides accurate change detection
- Handles file renames and moves

### Snapshot Structures
**DirectorySnapshot:**
```go
type DirectorySnapshot struct {
    Files map[string]FileInfo
    Root  string
}

type FileInfo struct {
    Path     string
    Size     int64
    ModTime  time.Time
    Hash     string
    IsDir    bool
}
```

### Snapshot Creation
**From Directory:**
```go
func CreateDirectorySnapshot(root string, exclusions []string) (*DirectorySnapshot, error) {
    snapshot := &DirectorySnapshot{
        Files: make(map[string]FileInfo),
        Root:  root,
    }
    
    err := WalkWithExclusions(root, exclusions, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        hash := calculateFileHash(path)
        snapshot.Files[path] = FileInfo{
            Path:    path,
            Size:    info.Size(),
            ModTime: info.ModTime(),
            Hash:    hash,
            IsDir:   info.IsDir(),
        }
        return nil
    })
    
    return snapshot, err
}
```

**From Archive:**
```go
func CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
    reader, err := zip.OpenReader(archivePath)
    if err != nil {
        return nil, err
    }
    defer reader.Close()
    
    snapshot := &DirectorySnapshot{
        Files: make(map[string]FileInfo),
        Root:  archivePath,
    }
    
    for _, file := range reader.File {
        hash := calculateArchiveFileHash(file)
        snapshot.Files[file.Name] = FileInfo{
            Path:    file.Name,
            Size:    int64(file.UncompressedSize64),
            ModTime: file.ModTime(),
            Hash:    hash,
            IsDir:   file.FileInfo().IsDir(),
        }
    }
    
    return snapshot, nil
}
```

### Comparison Algorithm
```go
func CompareSnapshots(snap1, snap2 *DirectorySnapshot) bool {
    if len(snap1.Files) != len(snap2.Files) {
        return false
    }
    
    for path, info1 := range snap1.Files {
        info2, exists := snap2.Files[path]
        if !exists {
            return false
        }
        
        if info1.Hash != info2.Hash || info1.Size != info2.Size {
            return false
        }
    }
    
    return true
}
```

**Code Markers**: `pkg/fileops/comparison.go`, snapshot structures, comparison functions

## 18. Exclusion Patterns Implementation [IMPL:EXCLUSION_PATTERNS] [ARCH:EXCLUSION_PATTERNS] [REQ:CONFIGURATION]

### Decision: Doublestar glob pattern matching with compiled matcher
**Rationale:**
- Supports flexible file filtering
- Enables efficient pattern matching
- Handles recursive patterns (`**`)
- Integrates with configuration system

### PatternMatcher Implementation
```go
type PatternMatcher struct {
    patterns []string
    compiled []glob.Glob
}

func NewPatternMatcher(patterns []string) *PatternMatcher {
    pm := &PatternMatcher{
        patterns: patterns,
        compiled: make([]glob.Glob, 0, len(patterns)),
    }
    
    for _, pattern := range patterns {
        if g, err := glob.Compile(pattern); err == nil {
            pm.compiled = append(pm.compiled, g)
        }
    }
    
    return pm
}

func (pm *PatternMatcher) Match(path string) bool {
    for _, g := range pm.compiled {
        if g.Match(path) {
            return true
        }
    }
    return false
}
```

### Convenience Functions
```go
func ShouldExcludeFile(path string, exclusions []string) bool {
    if len(exclusions) == 0 {
        return false
    }
    
    matcher := NewPatternMatcher(exclusions)
    return matcher.Match(path)
}

func FilterPaths(paths []string, exclusions []string) []string {
    matcher := NewPatternMatcher(exclusions)
    var filtered []string
    
    for _, path := range paths {
        if !matcher.Match(path) {
            filtered = append(filtered, path)
        }
    }
    
    return filtered
}
```

### Pattern Support
- **Doublestar**: `**` for recursive matching
- **Glob patterns**: `*.tmp`, `*.log`, `build/`
- **Directory patterns**: `.git/`, `node_modules/`
- **Multiple patterns**: Array of patterns for complex exclusions

**Code Markers**: `pkg/fileops/exclusion.go`, `PatternMatcher` struct, doublestar imports

## 19. Exclude Patterns Merge Testing [IMPL:TEST_EXCLUDE_MERGE] [ARCH:TEST_EXCLUDE_MERGE] [REQ:TEST_EXCLUDE_MERGE] [REQ:CONFIGURATION] [REQ:CFG_006]

### Decision: Comprehensive test scenario for exclude patterns merging and source tracking
**Rationale:**
- Validates merge behavior fix for exclude patterns
- Ensures source tracking shows correct attribution
- Provides regression protection
- Enables debugging through test scenarios

### Test Implementation Structure
```go
// Test scenario setup
func TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE(t *testing.T) {
    // 1. Create temp directory with config files
    // 2. Load config and verify merge
    // 3. Verify config command output
    // 4. Validate source attribution
}
```

### Test Scenarios
1. **Basic Merge Test**: Default + local config merge
2. **Source Attribution Test**: Config command shows correct source
3. **Deduplication Test**: Duplicate patterns handled correctly
4. **Order Preservation Test**: Defaults first, then additions

### Test Validation Points
- Merged patterns contain all patterns from all sources
- Config command output shows config file path (not "default")
- Inheritance chain shows all contributing files
- Patterns are deduplicated
- Order is preserved

**Code Markers**: `config_test.go`, `TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE`

**Cross-References**: [REQ:TEST_EXCLUDE_MERGE], [ARCH:TEST_EXCLUDE_MERGE], [REQ:CONFIGURATION], [REQ:CFG_006]

## 20. Array Field Default Merge Strategy Implementation [IMPL:EXCLUDE_MERGE_FIX] [ARCH:EXCLUDE_MERGE_FIX] [REQ:CFG_005] [REQ:CONFIGURATION]

### Decision: Implement array field default merge strategy to satisfy CFG-005 requirement
**Rationale:**
- Implements CFG-005 requirement that array fields default to merge (accumulate) strategy in all contexts (inheritance chains and sequential file processing)
- Fixes implementation bug where array fields were using "override" instead of "merge" by default
- Ensures default patterns are preserved when local config adds patterns
- Provides proper deduplication during merge
- Uses current state for merge operations to ensure proper accumulation
- Handles YAML unmarshaling type conversions gracefully
- Gracefully skips unknown config fields instead of aborting merge operations

### Implementation Changes

#### 1. Merge Strategy Detection Fix (`extractStrategy`)
**Location**: `config.go:1760-1783`
**Change**: Array fields (like `exclude_patterns`) now default to "merge" strategy instead of "override" in all contexts. This applies to both inheritance chains and sequential file processing, ensuring consistent behavior per CFG-005.

#### 2. Merge State Management Fix (`applyMergeStrategies`)
**Location**: `config.go:1700-1875`
**Changes**:
- Initialize `result` as a deep copy of `dst` instead of `DefaultConfig()` to preserve accumulated values
- Save `dst.ExcludePatterns` before calling `mergeConfigs` and restore after to prevent `mergeConfigs` from interfering with merge strategy processing
- Create temporary copy of `src` with `ExcludePatterns` set to `nil` before calling `mergeConfigs` to prevent double-processing
- Use `resultMap` (current state) instead of `dstMap` (original state) for merge operations
- Update `resultMap` after each merge operation to reflect changes
- Skip unknown config fields (like `inherit`, `verification`) gracefully instead of aborting merge

#### 3. Unknown Field Handling (`isKnownConfigField` and `applyMergeStrategies`)
**Location**: `config.go:1840-1844, 1861-1867, 2345-2358`
**New Function**: `isKnownConfigField` helper function checks if a field name is a valid config field before processing.
```go
func isKnownConfigField(key string) bool {
    knownFields := map[string]bool{
        "archive_dir_path":     true,
        "use_current_dir_name": true,
        "exclude_patterns":     true,
        "include_git_info":     true,
        "backup_dir_path":      true,
        "skip_broken_symlinks": true,
        "status_created_archive": true,
        "status_disk_full":      true,
    }
    return knownFields[key]
}
```
**Changes**:
- Skip `inherit` metadata field (used for inheritance processing, not config)
- Skip unknown config fields using `isKnownConfigField` check before processing
- Fallback error handler skips unknown fields if they slip through (defensive programming)
- Prevents merge operations from aborting when encountering metadata or unknown fields

#### 4. YAML Type Conversion Fix (`applyMerge`)
**Location**: `config.go:2115-2256`
**Changes**:
- Handle `[]interface{}` from YAML unmarshaling by converting to `[]string` before processing
- Convert both source (`value`) and destination (`dstValue`) from `[]interface{}` to `[]string` if needed
- Ensures merge logic correctly processes string slices regardless of YAML unmarshaling type
- Preserves existing merge logic with deduplication and order preservation

#### 5. Field Setting Type Conversion Fix (`setConfigField`)
**Location**: `config.go:2285-2343`
**Changes**:
- Added `else if` block to handle `[]interface{}` values for `exclude_patterns`
- Converts `[]interface{}` to `[]string` before assignment
- Ensures `applyReplace` and other merge operations work correctly with YAML-unmarshaled values

#### 6. Array Field Default Merge in All Contexts
**Location**: `config.go:1770-1783`
**Change**: Removed `inheritContext` condition that was preventing array fields from defaulting to merge in sequential file processing. Array fields now default to merge in all contexts (inheritance chains and sequential files), per CFG-005 requirement. Explicit prefixes (`!`, `+`, `^`, `=`) are still respected.

### Fix Summary
1. **Strategy Detection**: Array fields default to "merge" instead of "override" in all contexts
2. **State Management**: Uses current state (`resultMap`) instead of original state (`dstMap`)
3. **Pattern Preservation**: Explicitly copies default patterns before merge
4. **Deduplication**: Proper deduplication during merge operations
5. **Nil Handling**: Handles nil `dstValue` by getting current value from result
6. **YAML Type Conversion**: Handles `[]interface{}` from YAML unmarshaling in both `applyMerge` and `setConfigField`
7. **Unknown Field Handling**: Gracefully skips unknown config fields instead of aborting merge
8. **Metadata Field Handling**: Skips `inherit` and other metadata fields that aren't part of Config struct

### Test Updates
**Location**: `config_test.go`
**Changes**:
- Updated `TestLoadConfigMultipleFiles/default_search_path_processes_multiple_files` to use `!exclude_patterns` prefix where override behavior is expected
- Tests now explicitly use merge strategy prefixes to align with CFG-005 requirement that array fields default to merge

**Code Markers**: `config.go`, `extractStrategy`, `applyMergeStrategies`, `applyMerge`, `setConfigField`, `isKnownConfigField`

**Cross-References**: [REQ:CFG_005], [ARCH:EXCLUDE_MERGE_FIX], [REQ:CONFIGURATION], [REQ:TEST_EXCLUDE_MERGE]

## 20. Verification Implementation [IMPL:VERIFICATION] [ARCH:VERIFICATION] [REQ:ARCHIVE_VERIFICATION]

### Decision: Multi-algorithm checksum verification with status tracking
**Rationale:**
- Ensures archive integrity
- Supports multiple checksum algorithms
- Provides detailed verification status
- Enables corruption detection

### VerificationStatus Structure
```go
type VerificationStatus struct {
    VerifiedAt   time.Time `json:"verified_at"`
    IsVerified   bool      `json:"is_verified"`
    HasChecksums bool      `json:"has_checksums"`
    Errors       []string  `json:"errors,omitempty"`
}
```

### Archive Verification
```go
func VerifyArchive(archivePath string) (*VerificationStatus, error) {
    status := &VerificationStatus{
        VerifiedAt: time.Now(),
        IsVerified: true,
    }
    
    reader, err := zip.OpenReader(archivePath)
    if err != nil {
        status.IsVerified = false
        status.Errors = append(status.Errors, fmt.Sprintf("Failed to open archive: %v", err))
        return status, nil
    }
    defer reader.Close()
    
    for _, file := range reader.File {
        if err := verifyFile(file); err != nil {
            status.IsVerified = false
            status.Errors = append(status.Errors, err.Error())
        }
    }
    
    return status, nil
}
```

### File Verification
```go
func verifyFile(file *zip.File) error {
    rc, err := file.Open()
    if err != nil {
        return fmt.Errorf("failed to open file %s: %v", file.Name, err)
    }
    defer rc.Close()
    
    buf := make([]byte, 1024)
    _, err = rc.Read(buf)
    if err != nil && err != io.EOF {
        return fmt.Errorf("failed to read file %s: %v", file.Name, err)
    }
    
    return nil
}
```

### Checksum Generation
```go
func GenerateChecksums(fileMap map[string]string, algorithm string) (map[string]string, error) {
    checksums := make(map[string]string)
    
    for relPath, absPath := range fileMap {
        checksum, err := calculateFileChecksum(absPath, algorithm)
        if err != nil {
            return nil, fmt.Errorf("failed to calculate checksum for %s: %w", relPath, err)
        }
        checksums[relPath] = checksum
    }
    
    return checksums, nil
}

func calculateFileChecksum(filePath string, algorithm string) (string, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", err
    }
    
    var hash []byte
    switch algorithm {
    case "sha256":
        h := sha256.Sum256(data)
        hash = h[:]
    case "sha512":
        h := sha512.Sum512(data)
        hash = h[:]
    case "md5":
        h := md5.Sum(data)
        hash = h[:]
    default:
        return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
    }
    
    return hex.EncodeToString(hash), nil
}
```

**Code Markers**: `verify.go`, `VerificationStatus` struct, checksum functions

## 20. Configuration Reflection Implementation [IMPL:CFG_006] [ARCH:CFG_006] [REQ:CFG_006]

### Decision: Reflection-based field discovery with caching and lazy evaluation
**Rationale:**
- Automatic field discovery eliminates manual maintenance overhead
- Performance optimization ensures responsive configuration inspection
- Complete source tracking enables debugging and troubleshooting
- Backward compatibility preserved through wrapper functions

### Field Discovery Implementation

**GetAllConfigFields() Function:**
```go
// [HIGH] CFG-006: Complete config visibility - Automatic field discovery
func GetAllConfigFields(cfg *Config) []configFieldInfo {
    // Check cache first
    if cached := getCachedFields(); cached != nil {
        return cached
    }
    
    // Perform reflection-based discovery
    fields := reflectConfigFields(reflect.TypeOf(*cfg), "", nil)
    
    // Cache results
    cacheFields(fields)
    
    return fields
}
```

**Field Metadata Structure:**
```go
type configFieldInfo struct {
    Path     string   // Full field path (e.g., "Verification.VerifyOnCreate")
    Name     string   // Field name
    Type     string   // Type string representation
    Category string   // Field category (archive, backup, formatting, etc.)
    Value    interface{} // Current field value
}
```

**Recursive Field Traversal:**
- Handles nested structs (e.g., VerificationConfig)
- Supports embedded fields and anonymous structs
- Prevents circular references with visited type tracking
- Maximum depth limiting for safety
- Type-aware value extraction

**Field Categorization:**
- Automatic categorization based on field name patterns
- Categories: archive, backup, formatting, verification, status, template, pattern
- Supports custom category assignment

### Source Tracking Implementation

**GetAllConfigValuesWithSources() Function:**
```go
// [HIGH] CFG-006: Complete config visibility - Enhanced source tracking
func GetAllConfigValuesWithSources(cfg *Config, root string) []ConfigValueWithMetadata {
    fields := GetAllConfigFields(cfg)
    values := make([]ConfigValueWithMetadata, 0, len(fields))
    
    for _, field := range fields {
        value := getFieldValue(cfg, field.Path)
        source := determineSource(cfg, field.Path)
        
        values = append(values, ConfigValueWithMetadata{
            Field:  field,
            Value:  value,
            Source: source,
            // ... additional metadata
        })
    }
    
    return values
}
```

**ConfigValueWithMetadata Structure:**
```go
type ConfigValueWithMetadata struct {
    Field         configFieldInfo
    Value         interface{}
    Source        string              // Source type (default, file, environment, inheritance)
    SourceFile    string              // Source file path (if applicable)
    InheritanceChain []string         // Complete inheritance chain
    MergeStrategy string              // Merge strategy used (override, append, etc.)
    IsOverride    bool                // Whether value was overridden
}
```

**Source Determination:**
- Integration with CFG-005 inheritance system
- Environment variable detection
- Default value identification
- Inheritance chain tracking
- Merge strategy attribution

### Display Formatting Implementation

**Format Functions:**
```go
// [HIGH] CFG-006: Complete config visibility - Type-aware formatting
func formatFieldValue(value interface{}, fieldType reflect.Kind) string {
    switch fieldType {
    case reflect.String:
        return formatStringValue(value)
    case reflect.Slice:
        return formatSliceValue(value)
    case reflect.Bool:
        return formatBoolValue(value)
    // ... other types
    }
}
```

**Output Formats:**
- **Table Format**: Tabular display with columns (Field, Value, Source, Category)
- **Tree Format**: Hierarchical display showing inheritance relationships
- **JSON Format**: Structured JSON output for programmatic use

**Zero Value Detection:**
- Detects zero values for all types
- Displays appropriate zero value representation
- Handles nil pointers and empty collections

### Performance Optimization Implementation

**ConfigFieldCache Structure:**
```go
// [MEDIUM] CFG-006: Performance optimization - Reflection result caching
type configFieldCache struct {
    mu         sync.RWMutex
    fields     []configFieldInfo
    structHash uint64        // Hash of Config struct to detect changes
    lastUpdate time.Time
    valid      bool
}

var globalFieldCache = &configFieldCache{}
```

**Caching Strategy:**
- Thread-safe caching with sync.RWMutex
- Schema hash-based cache invalidation
- Automatic cache refresh on struct changes
- Singleton pattern for optimal memory usage

**Lazy Source Evaluation:**
```go
// [MEDIUM] CFG-006: Performance optimization - Lazy source evaluation
func GetConfigValuesWithSourcesFiltered(cfg *Config, root string, filter *ConfigFilter) []ConfigValueWithMetadata {
    fields := GetAllConfigFields(cfg)
    
    // Pre-filter fields before source evaluation
    filteredFields := applyFilter(fields, filter)
    
    // Only resolve sources for filtered fields
    values := make([]ConfigValueWithMetadata, 0, len(filteredFields))
    for _, field := range filteredFields {
        source := determineSource(cfg, field.Path) // Lazy evaluation
        values = append(values, ConfigValueWithMetadata{
            Field:  field,
            Source: source,
        })
    }
    
    return values
}
```

**Incremental Resolution:**
```go
// [MEDIUM] CFG-006: Performance optimization - Incremental resolution
func GetConfigFieldValue(cfg *Config, fieldPath string) (ConfigValueWithMetadata, error) {
    // Direct field access without full enumeration
    field := findFieldByPath(fieldPath)
    value := getFieldValue(cfg, fieldPath)
    source := determineSource(cfg, fieldPath)
    
    return ConfigValueWithMetadata{
        Field:  field,
        Value:  value,
        Source: source,
    }, nil
}
```

**Pattern-Based Queries:**
```go
// [MEDIUM] CFG-006: Performance optimization - Pattern-based queries
func GetConfigFieldByPattern(cfg *Config, pattern string) ([]configFieldInfo, error) {
    fields := GetAllConfigFields(cfg)
    matched := make([]configFieldInfo, 0)
    
    for _, field := range fields {
        if matchPattern(field.Path, pattern) {
            matched = append(matched, field)
        }
    }
    
    return matched, nil
}
```

### Testing Implementation

**Test Functions:**
- `TestConfigReflection()`: Main functionality tests
- `TestAdvancedFieldDiscovery()`: Edge case testing
- `TestSourceAttributionAccuracy()`: Source tracking validation
- `TestDisplayFormatting()`: Output format testing
- `TestFilteringFunctionality()`: Filtering validation
- `TestPerformanceOptimization()`: Performance benchmarks
- `BenchmarkConfigReflectionOperations()`: Performance benchmarks

**Test Coverage:**
- >95% coverage for all CFG-006 functionality
- Edge cases: anonymous fields, circular refs, error recovery
- Integration testing with CFG-005
- Performance validation with benchmarks

**Code Markers**: `config.go` (lines ~1653-1932), `config_test.go`, `config_bench_test.go`, `// [HIGH] CFG-006: Complete config visibility`, `// [MEDIUM] CFG-006: Performance optimization`

## 21. Data Models [IMPL:DATA_MODELS] [ARCH:SYSTEM_COMPONENTS] [REQ:*]

### Decision: Comprehensive data structures for archive and backup operations
**Rationale:**
- Provides structured data representation
- Enables type safety
- Supports serialization and persistence
- Facilitates testing and validation

### Archive and Backup Data Structures

**Archive Structure:**
```go
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
```

**Backup Structure:**
```go
type Backup struct {
    Name         string
    Path         string
    CreationTime time.Time
    SourceFile   string
    Note         string
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

**Formatter Structures:**
```go
type OutputFormatter struct {
    cfg *Config
}

type TemplateFormatter struct {
    config *Config
}
```

**Code Markers**: Archive and backup struct definitions, formatter types

## 11. Documentation Enhancement Framework [IMPL:DOC_ENHANCEMENT] [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_001]

### Decision: Implement comprehensive documentation enhancement system with semantic linking, synchronization, and traceability
**Rationale:**
- Cross-document repetition identified as highly valuable for LLM consumption
- Need automated systems to maintain consistency during changes
- Enhanced traceability prevents functionality loss during evolution
- Semantic linking improves understanding and change impact analysis

**Implementation Approach:**
- Semantic token system for cross-referencing
- Automated validation and consistency checking
- Behavioral contracts for feature stability
- Dependency mapping for change impact analysis

**Code Markers**: `// DOC-XXX: Documentation enhancement` tokens

## 12. Semantic Cross-Referencing Strategy [IMPL:SEMANTIC_CROSS_REF] [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_001]

### Decision: Use bi-directional linking with feature reference format including all document layers
**Rationale:**
- LLMs benefit from rich cross-references to understand relationships
- Bi-directional links prevent orphaned references
- Standardized format ensures consistency across documents
- Forward/backward/sibling links provide comprehensive navigation

**Implementation Approach:**
- Feature reference blocks with links to all documentation layers
- Automated link validation
- Cross-reference consistency checking
- Semantic token-based linking

**Code Markers**: LinkingEngine implementation, automated link validation

## 13. Enhanced Traceability Design [IMPL:TRACEABILITY] [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_003]

### Decision: Implement feature fingerprints with behavioral contracts and dependency mapping
**Rationale:**
- Feature fingerprints ensure stable identity through changes
- Behavioral contracts define what cannot change without version bump
- Dependency mapping shows change impact chains
- Automated regression prevention protects against functionality loss

**Implementation Approach:**
- Feature fingerprint generation and validation
- Behavioral contract definition and enforcement
- Dependency graph construction and analysis
- Automated change impact assessment

**Code Markers**: TraceabilitySystem, behavioral contract validation, dependency analysis

## 14. Token System Implementation [IMPL:TOKEN_SYSTEM] [ARCH:TOKEN_SYSTEM] [REQ:DOC_016]

### Decision: Implement comprehensive semantic token system for cross-layer traceability
**Rationale:**
- Ensures complete traceability from requirements through architecture to implementation
- Enables AI assistant navigation and feature discovery
- Supports automated validation of token consistency
- Maintains long-term code maintainability

**Implementation Approach:**
- Central token registry in `semantic-tokens.md` with cross-references
- Token format: `[TYPE:IDENTIFIER]` with UPPER_SNAKE_CASE identifiers
- Cross-layer token references: Implementation tokens reference `[ARCH:*]` and `[REQ:*]` tokens
- Automated validation scripts for token format, consistency, and traceability
- Token-based code navigation for AI assistants

**Token Usage in Code:**
```go
// [REQ:FILE_BACKUP] Create backup of single file with comparison
// [IMPL:ATOMIC_OPS] [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
    // ...
}
```

**Token Usage in Tests:**
```go
// Test validates [REQ:FILE_BACKUP] is met
func TestCreateFileBackup_REQ_FILE_BACKUP(t *testing.T) {
    // ...
}
```

**Validation Implementation:**
- Token format validation: Check `[TYPE:IDENTIFIER]` pattern compliance
- Cross-layer consistency: Verify tokens exist in all required layers
- Traceability validation: Ensure all requirements have implementation tokens
- Missing token detection: Identify features without proper token coverage

**Code Markers**: Token validation scripts, token registry maintenance, cross-reference validation

## 15. Pre-Extraction Refactoring Strategy [IMPL:REFACTOR_PREP] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

### Decision: Implement comprehensive pre-extraction refactoring to ensure clean component boundaries
**Rationale:**
- Prevents circular dependencies in extracted packages
- Ensures clean interfaces and minimal coupling
- Preserves backward compatibility during extraction
- Enables reliable, maintainable extracted components

**Implementation Approach:**
- Interface-first design before extraction
- Large file decomposition for better boundaries
- Dependency analysis and cleanup
- Comprehensive testing before extraction

**Code Markers**: `// REFACTOR-XXX: Preparation for extraction` tokens

## 16. Interface-First Extraction Approach [IMPL:INTERFACE_FIRST] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

### Decision: Define all component interfaces before extraction begins
**Rationale:**
- Prevents tight coupling between extracted packages
- Enables independent testing and development of components
- Provides clear contracts for component interaction
- Facilitates future evolution of implementations

**Implementation Approach:**
- Define interfaces before concrete implementations
- Use Go interfaces for abstraction
- Independent package testing
- Versioned interface evolution

**Code Markers**: Interface definitions with `// REFACTOR-001: Interface standardization`

## 17. Large File Decomposition Strategy [IMPL:LARGE_FILE_DECOMP] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

### Decision: Decompose large files (>1000 lines) before extraction for better package boundaries
**Rationale:**
- 1675-line formatter.go contains multiple logical components
- Clean separation enables focused extracted packages
- Reduces complexity and improves maintainability
- Enables independent evolution of formatter components

**Implementation Approach:**
- Identify logical component boundaries
- Extract components into separate files/packages
- Maintain interface compatibility
- Comprehensive testing after decomposition

**Code Markers**: `// REFACTOR-002: Component boundary` markings for logical separations

## 18. Feature Implementation Notes

### GIT-006: Configurable Git Dirty Status [IMPL:GIT_DIRTY_CONFIG] [ARCH:GIT_INTEGRATION] [REQ:GIT_INTEGRATION]

**Status**: ✅ COMPLETED  
**Priority**: High

**Description**: Make the Git dirty status indicator ('-dirty' suffix) configurable through the configuration system, allowing users to enable/disable this feature.

**Implementation Details:**
- Added `show_git_dirty_status` boolean option to Config struct
- Default to true for backward compatibility
- Updated Git status detection to respect configuration
- Modified archive naming functions to check the option before adding "-dirty" suffix
- Added proper merging of the new option in `mergeBasicSettings`

**Files Modified:**
- `config.go` - Added `show_git_dirty_status` option to Config struct
- `archive.go` - Updated archive naming to respect configuration

**Code Markers**: `show_git_dirty_status` configuration option, conditional dirty status in archive naming

### CFG-004: Eliminate Hardcoded Strings and Enhance Template Formatting [IMPL:CONFIGURABLE_STRINGS] [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]

**Status**: ✅ COMPLETED  
**Priority**: High

**Description**: Implement comprehensive string externalization system allowing all user-facing strings to be loaded from configuration files rather than hardcoded, with enhanced template formatting using named elements in data structures.

**Implementation Details:**
- Added 12 new error message format strings (FormatDiskFullError, FormatPermissionError, etc.)
- Added corresponding template versions with named data structure support
- Updated HandleArchiveError to use configurable error messages
- Maintained full backward compatibility with existing format strings
- All error messages now use OutputFormatter methods instead of hardcoded strings

**Files Modified:**
- `config.go` - Extended Config struct with comprehensive format strings
- `formatter.go` - Added OutputFormatter methods for all new format strings
- `errors.go` - Updated error handling to use configurable messages

**Code Markers**: Configurable format strings, template-based error formatting with named elements

### OUT-001: Delayed Output Management [IMPL:DELAYED_OUTPUT] [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]

**Status**: ✅ COMPLETED  
**Priority**: Medium

**Description**: Implement delayed output functionality by returning output messages to calling functions instead of direct stdout/stderr printing, enabling better control over when and how output is displayed.

**Implementation Details:**
- Added OutputMessage struct with Content, Destination, and Type fields
- Added OutputCollector with methods: AddStdout, AddStderr, FlushAll, FlushStdout, FlushStderr, Clear
- Enhanced OutputFormatter with optional collector field and delayed mode support
- Added NewOutputFormatterWithCollector constructor for delayed output mode
- Updated all Print methods to check for collector and route messages accordingly
- Maintained full backward compatibility - existing code works unchanged

**Files Modified:**
- `formatter.go` - Added OutputCollector and OutputMessage types, enhanced OutputFormatter with delayed output support

**Code Markers**: `// OUT-001: Delayed output`, OutputCollector, OutputMessage types

### TEST-001: Comprehensive formatter.go Test Coverage [IMPL:TEST_COVERAGE] [ARCH:TESTING_STRATEGY] [REQ:*]

**Status**: ✅ COMPLETED  
**Priority**: Medium

**Description**: Test all 0% coverage functions in formatter.go to improve code coverage and ensure reliability of OutputCollector, template methods, and error formatting functionality.

**Implementation Details:**
- Added TestOutputCollector with tests for NewOutputCollector, AddStdout, AddStderr, GetMessages, Clear
- Added TestDelayedOutputMode testing OutputFormatter with collector integration
- Added TestTemplateFormattingMethods for all template functions with 0% coverage
- Added TestErrorFormattingMethods for error formatting functions
- Added TestPrintMethods testing Print methods in delayed mode
- Added TestTemplateFormatterAdvanced for TemplateFormatter functionality
- Added TestOutputCollectorFlushMethods for comprehensive testing of FlushAll, FlushStdout, FlushStderr
- Eliminated ALL 0% coverage functions in formatter.go (152 total functions tracked)
- Achieved 116 functions at 100% coverage, 36 functions with partial coverage (75%+ typical)

**Files Modified:**
- `formatter_test.go` - Added comprehensive test suites for previously untested functions

**Code Markers**: Comprehensive test coverage for OutputCollector, template methods, and error formatting

## 19. Extraction Principles and Design Decisions [IMPL:EXTRACTION_PRINCIPLES] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

### Maintain Backward Compatibility [IMPL:BACKWARD_COMPAT] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Decision**: Extract components without breaking existing backup application  
**Rationale**: Allows gradual migration and maintains stability of existing functionality  
**Implementation**: Use Go modules and interfaces to isolate extracted code  
**Design Note**: Original application continues to work while providing reusable components

### Interface-Driven Design [IMPL:INTERFACE_DRIVEN] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Decision**: Create clear interfaces for all extracted components  
**Rationale**: Enables flexibility, testing, and future evolution of implementations  
**Implementation**: Define contracts before extracting concrete implementations  
**Design Note**: Prevents tight coupling between extracted packages

### Zero-Breaking-Change Extraction [IMPL:ZERO_BREAKING] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Decision**: Extraction must not change existing application behavior  
**Rationale**: Maintains confidence in extraction process and preserves existing functionality  
**Implementation**: Comprehensive test coverage and behavioral verification  
**Design Note**: Existing tests must continue to pass without modification

### Layered Extraction Approach [IMPL:LAYERED_EXTRACTION] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Decision**: Extract in dependency order - core utilities first, then higher-level components  
**Rationale**: Prevents circular dependencies and ensures stable foundation  
**Implementation**: Infrastructure (config, errors) → Utilities (formatter, git) → Framework (cli) → Patterns (processing)  
**Design Note**: Each layer builds on previous layers without circular references

## 20. Extraction Challenges and Solutions [IMPL:EXTRACTION_CHALLENGES] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

### Large File Decomposition [IMPL:LARGE_FILE_CHALLENGE] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Challenge**: `formatter.go` (1675 lines) is large and complex  
**Solution**: Break into multiple focused packages while maintaining interface compatibility  
**Approach**: Extract template engine, printf formatter, output collector, and ANSI support separately  
**Design Note**: Size indicates rich functionality perfect for reuse, but needs careful decomposition

### Configuration Schema Flexibility [IMPL:CONFIG_SCHEMA_FLEX] [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]

**Challenge**: Current config is backup-specific but needs to be generic  
**Solution**: Extract configuration loading/merging logic with pluggable schema validation  
**Approach**: Create ConfigLoader interface that can handle different struct types  
**Design Note**: Preserve powerful discovery and merging while enabling different schemas

### Dependency Management [IMPL:DEPENDENCY_MGMT] [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]

**Challenge**: Extracted packages will have interdependencies  
**Solution**: Design clear dependency hierarchy and use Go modules for versioning  
**Approach**: Core packages (config, errors) have no internal dependencies, higher-level packages compose core packages  
**Design Note**: Clear layering prevents circular dependencies and simplifies usage

### Testing Complexity [IMPL:TESTING_COMPLEXITY] [ARCH:TESTING_STRATEGY] [REQ:*]

**Challenge**: Extracted components need comprehensive testing without duplicating existing tests  
**Solution**: Create package-specific tests while ensuring original integration tests still pass  
**Approach**: Extract test utilities first, then create focused tests for each package  
**Design Note**: Comprehensive testing ensures extracted components are production-ready

## 21. Configuration Output Grouping Implementation [IMPL:CONFIG_OUTPUT_GROUPING] [ARCH:CONFIG_OUTPUT_GROUPING] [REQ:CONFIG_OUTPUT_GROUPING]

### Decision: Implement grouping via metadata extension and display refactoring
**Rationale:**
- Extends existing metadata system without breaking changes
- Reuses existing display logic where possible
- Provides flexible ranking system

### Implementation Approach:
1.  **Metadata Extension**:
    -   Added `Importance` field to `configFieldInfo` struct
    -   Added `CategoryPriority` map for category sorting
    -   Added `getFieldImportance` helper for assigning importance

2.  **Display Logic**:
    -   `displayConfigGrouped`: New function for grouped output
    -   `displayConfigFlat`: Renamed from `displayConfigTable`
    -   `handleEnhancedConfigCommand`: Updated to switch between formats

3.  **Flags**:
    -   Added `--flat` flag to `configCmd`

**Code Markers**: `displayConfigGrouped`, `Importance`, `CategoryPriority`

**Cross-References**: [ARCH:CONFIG_OUTPUT_GROUPING], [REQ:CONFIG_OUTPUT_GROUPING]

## 22. User-Customizable Format Strings Implementation [IMPL:CUSTOMIZABLE_FORMAT_STRINGS] [ARCH:CUSTOMIZABLE_FORMAT_STRINGS] [REQ:CUSTOMIZABLE_FORMAT_STRINGS]

### Decision: Implement format string validation and comprehensive documentation for user customization. FormatListArchive and FormatListBackup support both printf-style and template-style placeholders.

**Rationale:**
- Existing infrastructure already supports customization (YAML tags, config loading)
- Main work is validation and documentation to prevent user mistakes
- Validation provides helpful guidance for correct placeholder usage
- Documentation enables users to discover and customize format strings

### Implementation Components:

#### 1. Format String Validation

**Location**: `config.go` - Functions: `ValidateFormatString()`, `validateAllFormatStrings()`, `getExpectedPlaceholders()`, `extractPlaceholders()`

**Key Decisions:**
- Validation runs automatically on configuration load (non-fatal warnings)
- Supports both printf-style (`%s`, `%d`) and template-style (`#{path}`, `#{name}`) placeholders
- Special placeholders like `#{size_human}` are validated per field
- Warnings printed to stderr, configuration loading continues
- Expected placeholders defined per field in `getExpectedPlaceholders()`
- **CRITICAL DECISION**: Template-style placeholders (`#{name}`) are replaced using simple string replacement, NOT Go text/template
- **CRITICAL DECISION**: Go text/template (`tmpl.Execute`) is ONLY used for `{{.name}}` style templates, and ONLY after all `#{...}` patterns are replaced
- **CRITICAL DECISION**: If any `#{...}` patterns remain after replacement, Go text/template processing is skipped to prevent fmt package errors

**Integration Points:**
- Integrated into `LoadConfig()` and `LoadConfigWithInheritance()` functions
- Validation occurs after configuration merging and before returning config
- See code comments in `config.go` for detailed implementation

#### 2. Documentation Structure

**Files Created:**
- `docs/format-strings-reference.md`: Comprehensive reference listing all format strings, organized by category, with placeholders, defaults, and examples
- `example-custom-formats.yml`: Brief example configuration demonstrating common customizations

**Documentation Features:**
- Lists all available format strings with YAML field names
- Documents supported placeholders for each format string
- Explains printf-style vs template-style formats
- Provides examples of common customizations (emoji, internationalization)

#### 3. Testing Implementation

**Test Files:**
- `config_test.go`: Format string validation tests (`TestFormatStringValidation_REQ_CUSTOMIZABLE_FORMAT_STRINGS`, `TestCustomFormatStringsLoad_REQ_CUSTOMIZABLE_FORMAT_STRINGS`, `TestDefaultFormatStrings_REQ_CUSTOMIZABLE_FORMAT_STRINGS`, `TestFormatStringValidationWarnings_REQ_CUSTOMIZABLE_FORMAT_STRINGS`)

**Test Coverage:**
- Validates correct placeholder detection for printf and template styles
- Verifies custom format strings load correctly from configuration
- Ensures default format strings work when not specified
- Confirms validation warnings are printed to stderr

#### 4. Simplified List Formatting Implementation

**Location**: `formatter_adapter_simple.go` - Function: `formatListArchiveSimple()`

**Key Decisions:**
- **SIMPLIFIED DESIGN**: Single, clear code path for list formatting
- **TEMPLATE-ONLY PLACEHOLDERS**: Only supports `#{name}` style placeholders. Printf-style (`%s`, `%d`) is NOT supported.
- **ALWAYS GATHER STATS**: File statistics are always gathered (needed for template placeholders)
- **SIMPLE STRING REPLACEMENT**: Uses `strings.ReplaceAll()` for placeholder replacement - no complex template engine
- **PRIORITY-BASED FORMAT SELECTION**:
  1. Use `FormatListArchive` if it contains template placeholders (`#{`)
  2. If `FormatListArchive` is empty or doesn't contain `#{`, use `TemplateListArchive`
  3. If both are empty, use default template format `#{path} (size: #{size_human})\n`
- **COMPREHENSIVE TESTING**: Detailed test suite validates all scenarios including:
  - Template placeholders replacement
  - Multiple placeholders
  - Fallback to TemplateListArchive
  - Fallback to default format
  - Missing files (graceful degradation with "unknown" defaults)
  - Unknown placeholders (left as-is)

**Rationale:**
- **Simplicity**: Single code path with simple string replacement is easier to understand, test, and maintain
- **Reliability**: Always gathering stats ensures template placeholders always work
- **Testability**: Clear function boundaries enable comprehensive unit testing
- **Maintainability**: No complex template engine means fewer bugs and easier updates
- **Clarity**: Only `#{name}` placeholders eliminates confusion about printf vs template formats

**Implementation Details:**
- Function signature: `formatListArchiveSimple(cfg *Config, formatterInstance *FormatterAdapter, archivePath, creationTime string) string`
- Always populates data map with: path, creation_time, name, size, size_human, mtime, mode, type
- Provides defaults ("unknown", "0") when file stats can't be gathered
- Uses simple `strings.ReplaceAll()` for placeholder replacement - no template engine needed
- Unknown placeholders (not in data map) are left as-is in the output

#### 5. Code Structure

**Files Modified:**
- `config.go`: Added validation functions and integrated into load process. Updated default format strings and validation for `FormatListArchive` and `FormatListBackup`.
- `formatter.go`: Enhanced template formatter to handle missing `#{size_human}` with fallback. Updated `FormatListArchive` and `FormatListBackup` to support template placeholders.
- `pkg/formatter/formatter.go`: Updated `FormatListArchive` and `FormatListBackup` to support template placeholders with automatic file statistics gathering.
- `formatter_adapter.go`: Enhanced to gather file statistics for template formatting
- `pkg/formatter/template.go`: Added fallback handling for missing placeholders

**Files Created:**
- `docs/format-strings-reference.md`: Comprehensive reference documentation
- `example-custom-formats.yml`: Brief example configuration

**No Changes Required:**
- `config.go` struct: Already has YAML tags
- `LoadConfig()`: Already loads format strings (validation added)
- `mergeFormatStrings()`: Already merges format strings
- `DefaultConfig()`: Already provides defaults

### Validation Strategy:

**Non-Fatal Warnings:**
- Validation warnings printed to stderr
- Configuration loading continues
- Users informed of potential issues but not blocked

**Error Message Format:**
- `Warning: Field 'FieldName': unexpected placeholder 'placeholder'. Expected one of: [list]`
- Clear indication of field, unexpected placeholder, and expected alternatives

### Backward Compatibility:

- All format strings have defaults in `DefaultConfig()`
- Users who don't specify format strings get default behavior
- Existing configurations continue to work unchanged
- Validation is non-fatal (warnings only)
- No breaking changes to configuration schema

### Performance Considerations:

- Validation runs once at configuration load time
- Regex compilation cached for efficiency

#### 6. Placeholder Syntax Migration: %{...} to #{...}

**Location**: All format string processing code - `pkg/formatter/ai_core_formatter.go`, `pkg/formatter/placeholder_replace.go`, `pkg/formatter/template.go`, `ai_formatter_adapter.go`, `formatter_adapter_simple.go`, `config.go`, and all test files

**Key Decisions:**
- **MIGRATION COMPLETE**: All placeholder syntax migrated from `%{...}` to `#{...}` to avoid conflicts with Go's `fmt` package
- **ROOT CAUSE**: Go's `fmt` package interprets `%{` as the start of a format verb, causing `%!(EXTRA ...)` errors when placeholders weren't replaced
- **SOLUTION**: Changed all placeholder syntax from `%{key}` to `#{key}` throughout the codebase
- **IMPACT**: Eliminates fmt package misinterpretation errors, simplifies code, improves maintainability
- **BACKWARD COMPATIBILITY**: Breaking change - users with old `%{...}` format strings must update to `#{...}` syntax

**Files Updated:**
- **Core Implementation**: `pkg/formatter/ai_core_formatter.go` - Updated `FormatWithPlaceholders` to use `#{}` syntax
- **Placeholder Replacement**: `pkg/formatter/placeholder_replace.go` - Already using `#{}` syntax (no change needed)
- **Template Processing**: `pkg/formatter/template.go` - Updated all placeholder generation and replacement to use `#{}`
- **Adapters**: `ai_formatter_adapter.go`, `formatter_adapter.go` - Updated fallback replacement logic to use `#{}`
- **Configuration**: `config.go` - Updated all default format strings and validation to use `#{}`
- **Tests**: All test files updated to use `#{}` syntax in format strings and assertions
- **Documentation**: All documentation files updated to reflect `#{}` syntax

**Rationale:**
- **Reliability**: Eliminates fmt package conflicts that caused `%!(EXTRA ...)` errors
- **Clarity**: `#{}` syntax is clearly distinct from Go's fmt format verbs
- **Simplicity**: No need for complex escaping or workarounds to prevent fmt misinterpretation
- **Consistency**: Single placeholder syntax throughout the codebase

**Implementation Details:**
- All `fmt.Sprintf("%{%s}", key)` changed to `fmt.Sprintf("#{%s}", key)`
- All `strings.Contains(formatStr, "%{")` changed to `strings.Contains(formatStr, "#{")`
- All default format strings updated (e.g., `"#{path} (size: #{size_human})\n"`)
- All validation logic updated to recognize `#{}` placeholders
- All test assertions updated to check for `#{}` instead of `%{}`
- All documentation examples updated to use `#{}` syntax

**Migration Status**: ✅ Complete - All code, tests, and documentation updated
- No runtime overhead for format string usage

**Code Markers**: `ValidateFormatString`, `validateAllFormatStrings`, `getExpectedPlaceholders`, `extractPlaceholders`, `// [REQ:CUSTOMIZABLE_FORMAT_STRINGS]`

**Cross-References**: [ARCH:CUSTOMIZABLE_FORMAT_STRINGS], [REQ:CUSTOMIZABLE_FORMAT_STRINGS], [ARCH:CONFIG_SYSTEM], [ARCH:OUTPUT_FORMATTING], [IMPL:CONFIG_STRUCT], [IMPL:DUAL_FORMATTING]
