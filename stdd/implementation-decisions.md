# Implementation Decisions

**STDD Methodology Version**: 1.0.2

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
    IncludeGitInfo    bool     `yaml:"include_git_info"`      // Legacy - use Git.IncludeInfo
    ShowGitDirtyStatus bool    `yaml:"show_git_dirty_status"` // Legacy - use Git.ShowDirtyStatus
    
    // GIT-005: Git configuration for repository detection and information extraction
    Git *GitConfig `yaml:"git,omitempty"`
    
    // File backup settings
    BackupDirPath             string `yaml:"backup_dir_path"`
    UseCurrentDirNameForFiles bool   `yaml:"use_current_dir_name_for_files"`
    
    
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

### Configuration Display Flattening [IMPL:CONFIG_DISPLAY_FLATTENING] [ARCH:CFG_006] [REQ:CFG_006]

**Decision**: Configuration keys are displayed as top-level keys in `bkpdir config` output, even when they are nested in the internal struct.

**Rationale:**
- The reflection-based configuration discovery system (`GetAllConfigFields`) uses YAML tag names for display
- Nested struct fields (like `Git.IncludeBranch`) are displayed using their YAML tag names (e.g., `include_branch`) as top-level keys
- This provides a flat, consistent view of all configuration options regardless of internal struct organization
- Users see configuration keys as they appear in YAML files, not as nested struct paths

**Implementation Details:**
- Internal struct has `Git *GitConfig` with nested fields
- YAML files can use either:
  - Top-level keys: `include_git_info`, `include_branch`, `show_git_dirty_status` (legacy and current format)
  - Nested format: `git: { include_info: true, include_branch: true }` (also supported but flattened in display)
- The `bkpdir config` command uses `field.YAMLName` for display, which results in top-level keys
- Reflection system (`reflectConfigFields`) recursively processes nested structs and uses YAML tag names from each field
- Configuration merging (`mergeGitSettings`) maintains backward compatibility with both formats

**Configuration File Format:**
```yaml
# Top-level keys (actual format used in practice)
include_git_info: true
show_git_dirty_status: true
include_branch: true
include_hash: true

# Nested format (also supported internally, but displayed as top-level)
# git:
#   include_info: true
#   include_branch: true
```

**Code Markers**: `GetAllConfigFields`, `reflectConfigFields`, `field.YAMLName`, `config.go` line 3424

**Cross-References**: [ARCH:CFG_006], [REQ:CFG_006], [IMPL:CFG_006]

### Default Values
- ArchiveDirPath: "../.bkpdir"
- UseCurrentDirName: true
- ExcludePatterns: [".git/", "vendor/"]
- IncludeGitInfo: false
- ShowGitDirtyStatus: true
- Git: DefaultGitConfig() (includes nested GitConfig with defaults)
- BackupDirPath: "../.bkpdir"
- UseCurrentDirNameForFiles: true
- Status codes: Various defaults (0 for success, non-zero for errors)
- Format strings: Default printf-style formats
- Template strings: Default template formats with placeholders

## 2. ZIP Archive Format Implementation [IMPL:ZIP_FORMAT] [ARCH:ARCHIVE_FORMAT]

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

## 9. Testing Implementation

## EXTRACT-009 Test Utilities Implementation [IMPL:EXTRACT-009_TESTUTIL] [ARCH:EXTRACT-009_TESTUTIL] [REQ:EXTRACT-009]

### Summary
Implement `pkg/testutil` with the following components:

- `interfaces.go` — small interfaces for providers (TestUtilProvider) to avoid tight coupling
- `assertions.go` — simple test assertion helpers (`AssertEqualString`, `AssertBool`, etc.)
- `fshelpers.go` — temp directory and file helpers, `CreateTempDir`, `WriteFile`, `CreateTestDirectoryStructure`
- `archivehelpers.go` — create zip archives for tests (`CreateTestZipArchive`)
- `githelpers.go` — create lightweight git repos for testing (`InitTestGitRepo`)
- `cmdhelpers.go` — helpers to execute Cobra commands and capture output (`ExecuteCommand`) with cleanup
- `provider.go` — default provider that composes utilities for convenience
- `doc.go` + godoc comments and examples
- `testutil_test.go` — unit tests for all exported helpers

### Key Implementation Notes
- Keep `pkg/testutil` free of project-specific domain types — accept simple primitives or small interfaces.
- Minimize external dependencies; prefer stdlib for file and zip operations.
- Provide migration adapters so existing tests can switch incrementally to `pkg/testutil`.
- Add `[IMPL:EXTRACT-009_TESTUTIL]` tokens in code comments for exported functions and tests that rely on this package.

### Validation
- Add unit tests exercising assertion helpers, temp dir lifecycle, archive creation, and command execution patterns.
- Migrate a small set of tests (e.g., `pkg/config`, `pkg/formatter`, `pkg/fileops`) to demonstrate usage and validate no regressions.


 [IMPL:TESTING] [ARCH:TESTING_STRATEGY] [REQ:*]

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

**Cross-References**: [ARCH:FILE_STATISTICS], [REQ:OUTPUT_FORMATTING], [REQ:OUT_002]

## 17. File Statistics Template Processing Fix [IMPL:FILE_STATISTICS_TEMPLATE_FIX] [ARCH:FILE_STATISTICS] [REQ:OUT_002] [REQ:OUTPUT_FORMATTING]

### Decision: Simplified template processing with direct placeholder replacement and fixed AIFormatterAdapter delegation
**Rationale:**
- Fixed critical bug where `AIFormatterAdapter` was bypassing template processing
- Simplified template processing to use direct placeholder replacement (same approach as `formatListArchiveSimple`)
- More reliable than complex `formatTemplate` function which had issues with Go's `text/template` internal `fmt` usage
- Ensures consistent behavior between `OutputFormatter` and `AIFormatterAdapter`
- Final safety check now correctly checks data map first before using defaults

### Implementation Approach:
- **Simplified Template Processing**: Replaced complex `formatTemplate` call with direct placeholder replacement loop in `FormatCreatedArchiveWithStats` and `FormatIncrementalCreatedWithStats`
  - Direct `strings.ReplaceAll` loop for all placeholders from data map
  - Final safety check that checks data map first before using defaults
  - Avoids issues with Go's `text/template` internal `fmt` usage
- **Fixed AIFormatterAdapter Delegation**: Updated `FormatCreatedArchiveWithStats` and `FormatIncrementalCreatedWithStats` in `AIFormatterAdapter` to correctly delegate to `OutputFormatter` implementation
  - Previously was calling basic `FormatCreatedArchive` and `FormatIncrementalCreated`, bypassing template processing
  - Now creates temporary `OutputFormatter` instance to use its template processing logic
  - Ensures consistent behavior regardless of which formatter implementation is used

### Template Processing Logic:
```go
// Simplified approach (same as formatListArchiveSimple)
result := templateStr
for key, value := range data {
    placeholder := fmt.Sprintf("#{%s}", key)
    result = strings.ReplaceAll(result, placeholder, value)
}

// Final safety check: check data map first before using defaults
if strings.Contains(result, "#{") {
    finalReplacements := map[string]string{
        "#{size_human}": func() string {
            if sh, ok := data["size_human"]; ok {
                return sh
            }
            return "unknown"
        }(),
        // ... similar for other placeholders
    }
    for placeholder, value := range finalReplacements {
        result = strings.ReplaceAll(result, placeholder, value)
    }
}
```

### AIFormatterAdapter Fix:
```go
func (fa *AIFormatterAdapter) FormatCreatedArchiveWithStats(path string) string {
    // Delegate to the OutputFormatter implementation which has the correct template processing
    outputFormatter := &OutputFormatter{cfg: fa.config}
    return outputFormatter.FormatCreatedArchiveWithStats(path)
}
```

**Code Structure:**
- `formatter.go` lines ~1388-1467: `FormatCreatedArchiveWithStats()` - Simplified template processing with direct placeholder replacement
- `formatter.go` lines ~1468-1547: `FormatIncrementalCreatedWithStats()` - Simplified template processing with direct placeholder replacement
- `ai_formatter_adapter.go` lines ~977-989: `FormatCreatedArchiveWithStats()` and `FormatIncrementalCreatedWithStats()` - Fixed to delegate to `OutputFormatter`

**Code Markers**:
- `formatter.go` line ~1417-1423: Direct placeholder replacement loop in `FormatCreatedArchiveWithStats`
- `formatter.go` line ~1427-1478: Final safety check with data map lookup first
- `formatter.go` line ~1496-1503: Direct placeholder replacement loop in `FormatIncrementalCreatedWithStats`
- `ai_formatter_adapter.go` line ~977-982: Fixed `FormatCreatedArchiveWithStats` delegation
- `ai_formatter_adapter.go` line ~984-989: Fixed `FormatIncrementalCreatedWithStats` delegation

**Cross-References**: [ARCH:FILE_STATISTICS], [REQ:OUT_002], [REQ:OUTPUT_FORMATTING], [IMPL:FILE_STATISTICS]

## 18. Directory Comparison Implementation [IMPL:DIRECTORY_COMPARISON] [ARCH:DIRECTORY_COMPARISON]

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
    Path     string   // Full field path (e.g., "Git.IncludeInfo")
    Name     string   // Field name
    Type     string   // Type string representation
    Category string   // Field category (archive, backup, formatting, etc.)
    Value    interface{} // Current field value
}
```

**Recursive Field Traversal:**
- Handles nested structs (e.g., GitConfig)
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

### Legacy Registry Migration Implementation [IMPL:TOKEN_SYSTEM] [ARCH:TOKEN_SYSTEM] [REQ:DOC_016] [REQ:GOV_REGISTRY_COMPLETENESS]

**Status**: ✅ Completed

**Decision**: Implement a deterministic migration workflow that relocates every legacy token entry from `project-tokens.yaml` into `stdd/semantic-tokens.md`, generating STDD-compliant cross-links during the transfer.

**Implementation Approach:**
1. **Structured Extraction**
   - Use `.venv/bin/python` helpers (PyYAML) to parse `project-tokens.yaml` and emit normalized JSON for each group (`actions`, `features`, `semantic_tokens`).
   - Capture metadata (description, status, source paths, subtasks) for downstream documentation.
2. **Cross-Link Synthesis**
   - For each legacy token, resolve canonical `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`, and test references using prefix heuristics (e.g., `DIR-001 → [REQ:IMMUTABLE_DIRECTORY_OPERATIONS]`).
   - Flag gaps by opening subtasks in `stdd/tasks.md` whenever a requirement/decision/test link is missing.
3. **Registry Authoring**
   - Generate Markdown tables per token category with columns: Description, Linked Requirement(s), Linked Architecture, Linked Implementation, Linked Tests/Code, Source Document, Status.
   - Embed migration status badges so governance tooling can verify coverage.
4. **YAML Decommissioning**
   - After each category is represented in markdown, replace the YAML entries with a pointer stub until the file can be removed in a final cleanup change. `project-tokens.yaml` now contains a canonical pointer to `stdd/semantic-tokens.md`.
5. **Validation Hooks**
   - Run `scripts/validate-token-traceability.sh`, `scripts/token-coverage-analysis.sh`, and targeted `rg "\[TOKEN\]"` sweeps to prove each migrated token exists in requirements, architecture decisions, implementation notes, code, and tests.
   - Record validation results inside the migration tables for auditability.

**Migration Outcome:**
- All legacy tokens were migrated into `stdd/semantic-tokens.md` and validated via `scripts/token-coverage-analysis.sh` and `scripts/validate-token-traceability.sh`.
- Validation report: `docs/validation-reports/token-coverage-analysis.md` (100% source/test/overall coverage).
- `project-tokens.yaml` replaced with a pointer stub to the canonical STDD registry.
- Created implementation marker: `[IMPL:TOKEN_MIGRATION_COMPLETE]` to indicate the migration is finished and audited.

**Pseudo-Code:**
```bash
# Extraction artifacts archived in docs/governance and migration logs
.venv/bin/python scripts/token_inventory.py > /tmp/tokens.json
.venv/bin/python scripts/render_token_tables.py --group features >> stdd/semantic-tokens.md
scripts/validate-token-traceability.sh
scripts/token-coverage-analysis.sh
```

**Code Markers**: `project-tokens.yaml` migration scripts, `stdd/semantic-tokens.md` tables annotated with `[REQ:DOC_016]`, validation logs appended to docs/governance summaries.

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

## 23. List Command Format Safety Fix [IMPL:LIST_FORMAT_SAFETY] [ARCH:OUTPUT_FORMATTING] [REQ:OUT_002]

### Decision: Guard printf usage and ensure template-style list formatting is used safely
**Rationale:**
- Prevent `fmt` from interpreting template placeholders or appending EXTRA args when format strings contain no printf verbs (this produced `%!(EXTRA ...)` in user output).
- Ensure AI formatter and core formatter consistently handle template placeholders (`#{...}`) and printf formats.

### Implementation Approach:
- Added guarded printf usage in `FormatListArchive` and `FormatListBackup` to call `fmt.Sprintf` only when the format string contains printf verbs (`%`). Otherwise, return the template-format string as-is or perform placeholder replacement.
- Updated AI core formatter (`pkg/formatter/ai_core_formatter.go`) `FormatWithContext` list handling to detect template placeholders (`#{`) and use `FormatWithPlaceholders` with gathered file stats when present; otherwise use `fmt.Sprintf` only when printf verbs exist.
- Ensured `AIFormatterAdapter` and `AIFormatterAdapter.FormatListArchiveWithExtraction` delegate to placeholder-based formatting and gather file stats for `#{size_human}` and related placeholders.

**Code Markers:** `formatter.go` (FormatListArchive, FormatListBackup), `pkg/formatter/ai_core_formatter.go` (FormatWithContext list handling), `ai_formatter_adapter.go` (FormatListArchiveWithExtraction, FormatListBackupWithExtraction)

**Cross-References:** [ARCH:OUTPUT_FORMATTING], [REQ:OUT_002], [IMPL:FILE_STATISTICS]

## 23. Configuration File Precedence Fix [IMPL:CFG_PRECEDENCE_FIX] [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]

### Decision: Fix configuration merge logic to respect earlier file precedence for sequential file processing
**Rationale:**
- Bug fix: Configuration files were not respecting the requirement that "earlier files take precedence over later files" when processing sequential config files
- The `mergeConfigs` and `mergeBasicSettings` functions were unconditionally overriding values without checking if they were already set by earlier files
- This caused local config files (e.g., `./.bkpdir.yml`) to be overridden by home directory config files (e.g., `~/.bkpdir.yml`), violating the precedence specification
- Special case: When an earlier file explicitly sets a value to the default (e.g., `include_git_info: false`), we need to detect that it was explicitly set, not just compare against defaults
- The fix ensures that when processing sequential files (not inheritance chains), values set by earlier files are preserved, even if they equal the default

### Implementation Approach:
- Modified `LoadConfig` fallback path to use `fileProcessed` flag instead of checking if cfg equals defaults (which fails when first file sets values to defaults)
- Modified `applyMergeStrategies` to accept `initialDefaultCfg` parameter to track initial defaults before any files were processed
- Modified `mergeConfigs` to accept `inheritContext`, `defaultCfg`, `rawSrcMap`, `initialDefaultCfg`, and `dstBeforeMerge` parameters
- Added `dstBeforeMerge` tracking: Save state of `dst` before merge to detect if earlier files modified values
- Updated `mergeBasicSettings`, `mergeFileBackupSettings`, and `mergeGitSettings` to check if `dstBeforeMerge` equals `initialDefaultCfg` to detect if earlier files modified values
- Use `rawSrcMap` to check if fields were explicitly set in source files (not just present in unmarshaled Config)
- Updated `applyReplace` to respect earlier file precedence even with explicit `!` (replace) prefix

**Key Changes:**
1. `LoadConfig` fallback: Uses `fileProcessed` flag and `initialDefaultCfg` to track initial state
2. `applyMergeStrategies(dst, src *Config, inheritContext bool, rawSrcMap map[string]interface{}, initialDefaultCfg *Config)` - Added `initialDefaultCfg` parameter
3. `mergeConfigs(dst, src *Config, inheritContext bool, defaultCfg *Config, rawSrcMap map[string]interface{}, initialDefaultCfg *Config, dstBeforeMerge *Config)` - Added context and state tracking parameters
4. `mergeBasicSettings` - Added precedence check: `dstWasNotModified := dstBeforeMerge != nil && initialDefaultCfg != nil && dstBeforeMerge.IncludeGitInfo == initialDefaultCfg.IncludeGitInfo`
5. `applyReplace` - Added precedence check for sequential file processing
6. All merge functions updated to accept and use new parameters

**Behavior:**
- **Inheritance chains** (`inheritContext=true`): Child configs override parent configs (normal inheritance behavior)
- **Sequential files** (`inheritContext=false`): Earlier files take precedence - if `dstBeforeMerge` equals `initialDefaultCfg`, no earlier file modified the field, so allow setting; otherwise preserve earlier file's value
- **Explicit field setting**: Use `rawSrcMap` to check if a field was explicitly set in the source file (not just present due to defaults)
- **Default value handling**: When first file sets `include_git_info: false` (default), `dstBeforeMerge` for second file will have `false` (from first file), which equals `initialDefaultCfg.IncludeGitInfo` (also `false`), so we correctly detect that no earlier file modified it from initial default

**Code Markers**: `mergeConfigs`, `mergeBasicSettings`, `applyOverride`, `applyReplace`, `inheritContext`, `fileProcessed`, `initialDefaultCfg`, `dstBeforeMerge`, `rawSrcMap`, `// CFG-001: Earlier files take precedence`

**Test Case:**
- Local config: `/tmp/myapp/.bkpdir.yml` with `archive_dir_path: "./archives"` and `include_git_info: false`
- Home config: `/Users/fareed/.bkpdir.yml` with `archive_dir_path: "/Users/fareed/.bkpdir"` and `include_git_info: true`
- Expected: Local config values (`"./archives"` and `false`) are preserved
- Before fix: Home config values (`"/Users/fareed/.bkpdir"` and `true`) were incorrectly used
- After fix: Local config values (`"./archives"` and `false`) are correctly preserved

**Cross-References**: [ARCH:CONFIG_SYSTEM], [REQ:CONFIGURATION], [IMPL:CONFIG_STRUCT]

## 24. Field-Level Merge Behavior Registry [IMPL:CFG_MERGE_BEHAVIOR_REGISTRY] [ARCH:CFG_005] [ARCH:CFG_001] [REQ:CFG_005] [REQ:CONFIGURATION]

### Decision: Implement field-level merge behavior registry to resolve conflict between CFG-001 and CFG-005
**Rationale:**
- Conflict resolution: CFG-005 requires array fields to merge/accumulate by default, while CFG-001 requires earlier files to take precedence
- Solution: Field-level registry specifies merge behavior per field, allowing different fields to have different precedence rules
- Enables fine-grained control: Some fields (like `exclude_patterns`) should always merge, while others (like `archive_dir_path`) should respect earlier file precedence
- Maintains backward compatibility: Default behavior preserved, explicit prefixes still work

### Implementation Approach:
- Created `FieldMergeBehavior` type with two behaviors: `MergeBehaviorAccumulate` and `MergeBehaviorPrecedence`
- Created `fieldMergeBehaviors` registry map specifying behavior per field
- Updated `getFieldMergeBehavior()` helper function to retrieve field-specific behavior
- Updated `applyMergeStrategies()` to use registry for default strategy selection
- Updated `applyReplace()` to check registry and respect field-specific behavior
- Updated explicit prefix detection to work for all fields (not just `exclude_patterns`)

**Key Changes:**
1. `FieldMergeBehavior` type and constants (`MergeBehaviorAccumulate`, `MergeBehaviorPrecedence`)
2. `fieldMergeBehaviors` registry map with field-specific behaviors
3. `getFieldMergeBehavior(fieldName string) FieldMergeBehavior` helper function
4. `applyMergeStrategies()`: Uses registry to determine default strategy for accumulate fields
5. `applyReplace()`: Checks registry to determine if precedence check should apply
6. Explicit prefix detection: Now works for all fields via `hasExplicitPrefix` map

**Field Behavior Matrix:**
| Field | Behavior | Default Strategy | Sequential Files | Inheritance Chains |
|-------|----------|------------------|------------------|-------------------|
| `exclude_patterns` | Accumulate | Merge | Merge | Merge |
| `archive_dir_path` | Precedence | Override | Earlier wins | Child overrides |
| `include_git_info` | Precedence | Override | Earlier wins | Child overrides |
| `!exclude_patterns` | Accumulate* | Replace | Earlier wins** | Replace |

*Note: Even with `!` prefix, `exclude_patterns` is still marked as accumulate, but explicit prefix overrides default merge strategy.  
**Note: For sequential files, earlier file precedence applies even with `!` prefix if field was set by earlier file (CFG-001).

**Behavior:**
- **MergeBehaviorAccumulate fields** (e.g., `exclude_patterns`):
  - Default (no prefix): Always merge/accumulate (CFG-005)
  - Explicit `!` prefix: Replace, but still respect earlier file precedence if field was set by earlier file (CFG-001)
  - Explicit `+` prefix: Explicit merge (same as default for this field)
  - Explicit `^` prefix: Prepend to existing values
  - Explicit `=` prefix: Use only if not set by earlier file

- **MergeBehaviorPrecedence fields** (e.g., `archive_dir_path`, `include_git_info`):
  - Default (no prefix): Earlier files take precedence (CFG-001)
  - Explicit prefixes: Work as normal, but precedence still applies for sequential files

**Code Markers**: `FieldMergeBehavior`, `fieldMergeBehaviors`, `getFieldMergeBehavior`, `applyMergeStrategies`, `applyReplace`, `// CFG-001 + CFG-005: Field merge behavior configuration`

**Test Cases:**
- `exclude_patterns` without prefix: Merges from all files (CFG-005)
- `exclude_patterns` with `!` prefix in later file: Respects earlier file if it was set (CFG-001)
- `archive_dir_path` in later file: Earlier file value preserved (CFG-001)
- `include_git_info` in later file: Earlier file value preserved (CFG-001)
- `+archive_dir_path` in later file: Earlier file value preserved (CFG-001) - Fixed in [IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX]
- `^archive_dir_path` in later file: Earlier file value preserved (CFG-001) - Fixed in [IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX]

**Cross-References**: [ARCH:CFG_005], [ARCH:CFG_001], [REQ:CFG_005], [REQ:CONFIGURATION], [IMPL:CFG_PRECEDENCE_FIX]

## 25. Mixed-Mode Merge Strategy Fix [IMPL:CFG_MIXED_MODE_MERGE_FIX] [ARCH:CFG_005] [ARCH:CFG_001] [REQ:CFG_005] [REQ:CONFIGURATION]

### Decision: Fix mixed-mode merge to ensure CFG-005 merge behavior for accumulate fields in all contexts
**Rationale:**
- Bug fix: `exclude_patterns` and other `MergeBehaviorAccumulate` fields were not consistently merging with defaults per CFG-005 requirement
- The logic was incorrectly keeping "override" strategy for first file in sequential processing, causing defaults to be replaced instead of merged
- Test expectations were misaligned with CFG-005 requirement - updated to reflect correct merge behavior
- Ensures consistent behavior: accumulate fields always merge (unless explicit `!` prefix), regardless of context

### Implementation Approach:
- Simplified merge strategy logic in `applyMergeStrategies` to always change "override" to "merge" for `MergeBehaviorAccumulate` fields (unless explicit prefix)
- Updated `applyOverride` to check if earlier files were processed and preserve their state (including defaults)
- Updated `applyMergeOperation` to pass `originalDstValue` (before `mergeConfigs` modifications) for accurate precedence checking
- Updated `applyOverride` to reset values back to original when override is prevented
- Updated test expectations in `TestLoadConfigMultipleFiles` to match CFG-005 behavior (merge with defaults)

### Key Changes:

#### 1. Simplified Merge Strategy Logic (`applyMergeStrategies`)
**Change**: For `MergeBehaviorAccumulate` fields with no explicit prefix, always change "override" to "merge" strategy. Removed complex conditional logic that was preventing merge in certain contexts.

**Before**: Complex logic checking `isFirstFileInSequential`, `isTrueInheritanceChain`, etc., sometimes keeping "override" strategy
**After**: Simple logic - always merge for accumulate fields unless explicit `!` prefix is used

**Code Location**: `config.go`, `applyMergeStrategies` function, lines ~1963-1985

#### 2. Earlier File Precedence Enhancement (`applyOverride`)
**Change**: Enhanced precedence checking to preserve earlier file state even when values equal defaults. Added check for `earlierFilesProcessed` to ensure later files cannot override earlier file's state.

**Before**: Only checked if field was explicitly set by earlier file or if value differed from default
**After**: Also checks if any earlier files were processed, preserving their entire state (including defaults)

**Code Location**: `config.go`, `applyOverride` function, lines ~2460-2485

#### 3. Original Value Tracking (`applyMergeOperation`)
**Change**: Added `originalDstValue` parameter to track values before `mergeConfigs` modifications. This ensures precedence checks use correct original values, not values modified by `mergeConfigs`.

**Before**: Used `dstValue` which could be modified by `mergeConfigs` before precedence check
**After**: Uses `originalDstValue` from before `mergeConfigs` modifications

**Code Location**: `config.go`, `applyMergeOperation` function, lines ~2422-2453

#### 4. Value Reset on Override Prevention (`applyOverride`)
**Change**: When override is prevented due to earlier file precedence, reset the value back to original `dstValue` to undo any modifications made by `mergeConfigs`.

**Before**: Just returned `nil` without resetting value, leaving `mergeConfigs` modifications in place
**After**: Calls `setConfigField(result, key, dstValue)` to reset to original value

**Code Location**: `config.go`, `applyOverride` function, line ~2483

#### 5. Test Expectation Update (`TestLoadConfigMultipleFiles`)
**Change**: Updated test expectation to match CFG-005 behavior - `exclude_patterns` should merge with defaults, not replace them.

**Before**: Expected `[local1, local2]` (replacement behavior)
**After**: Expected `[.git/ vendor/ local1, local2]` (merge behavior per CFG-005)

**Code Location**: `config_test.go`, `TestLoadConfigMultipleFiles` function, line ~240

### Behavior:

- **Single file scenario**: `exclude_patterns` merges with defaults → `[.git/ vendor/ user_patterns]` (CFG-005)
- **Multiple files, first file**: Merges with defaults → `[.git/ vendor/ first_file_patterns]` (CFG-005)
- **Multiple files, subsequent files**: Merges with accumulated result → `[.git/ vendor/ first_file_patterns second_file_patterns]` (CFG-005)
- **Explicit `!` prefix**: Replaces, but still respects earlier file precedence if field was set by earlier file (CFG-001)
- **Earlier file precedence**: For `MergeBehaviorPrecedence` fields, earlier files take precedence even if values equal defaults

### Test Cases:
- `TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE`: Validates merge with defaults for single file
- `TestLoadConfigMultipleFiles`: Validates merge with defaults and earlier file precedence
- `TestLoadConfigMultipleFiles/default_search_path_processes_multiple_files`: Validates earlier file precedence for non-accumulate fields

### Validation:
- ✅ `exclude_patterns` merges with defaults in all contexts (CFG-005)
- ✅ Earlier file precedence preserved for `MergeBehaviorPrecedence` fields (CFG-001)
- ✅ Explicit `!` prefix still respects earlier file precedence (CFG-001)
- ✅ All tests pass with updated expectations

**Code Markers**: `applyMergeStrategies`, `applyOverride`, `applyMergeOperation`, `applyReplace`, `originalDstValue`, `earlierFilesProcessed`, `// CFG-005: Array fields default to merge`, `// CFG-001: Earlier files take precedence`

**Cross-References**: [ARCH:CFG_005], [ARCH:CFG_001], [REQ:CFG_005], [REQ:CONFIGURATION], [IMPL:CFG_PRECEDENCE_FIX], [IMPL:EXCLUDE_MERGE_FIX]

---

## 26. Quoted YAML Key Prefix Support [IMPL:CFG_QUOTED_KEY_PREFIX] [ARCH:CFG_005] [REQ:CFG_005] [REQ:CONFIGURATION]

**Date**: 2025-12-11  
**Status**: Implemented  
**Priority**: P1 (Important)

### Problem:
When YAML keys with merge strategy prefixes (`+`, `^`, `!`, `=`) are quoted (e.g., `"+exclude_patterns"`), the YAML parser preserves the quotes as part of the key string. This caused the prefix detection logic to fail because it was checking the first character of the key, which was `"` instead of `+`.

**Symptoms**:
- Merge strategy prefixes in quoted YAML keys were not recognized
- Tests using quoted keys (e.g., `"+exclude_patterns":`) failed because the strategy was extracted as "override" instead of "merge"
- Inheritance chain merge strategies with quoted keys did not work correctly

### Solution:
Updated three key functions to handle quoted YAML keys:

1. **`extractStrategy` function**: Strip quotes before checking for merge strategy prefixes
2. **`hasExplicitPrefix` detection**: Strip quotes before checking for prefix characters
3. **`processKeys` function**: Prioritize keys with explicit prefixes when duplicate keys exist

### Implementation Details:

#### 1. Quote Handling in `extractStrategy`
**Change**: Added logic to detect and strip quotes (both double `"` and single `'`) from YAML keys before checking for merge strategy prefixes.

**Before**: 
```go
switch key[0] {
case '+':
    return "merge", key[1:]
```

**After**:
```go
// Handle quoted YAML keys
cleanKey := key
if len(key) >= 2 && ((key[0] == '"' && key[len(key)-1] == '"') || (key[0] == '\'' && key[len(key)-1] == '\'')) {
    cleanKey = key[1 : len(key)-1]
}
switch cleanKey[0] {
case '+':
    return "merge", cleanKey[1:]
```

**Code Location**: `config.go`, `extractStrategy` function, lines ~2394-2411

#### 2. Quote Handling in `hasExplicitPrefix` Detection
**Change**: Updated the prefix detection logic to strip quotes before checking for prefix characters.

**Before**:
```go
if len(origKey) > 0 && (origKey[0] == '!' || origKey[0] == '+' || origKey[0] == '^' || origKey[0] == '=') {
```

**After**:
```go
// Handle quoted YAML keys
cleanKey := origKey
if len(origKey) >= 2 && ((origKey[0] == '"' && origKey[len(origKey)-1] == '"') || (origKey[0] == '\'' && origKey[len(origKey)-1] == '\'')) {
    cleanKey = origKey[1 : len(origKey)-1]
}
if len(cleanKey) > 0 && (cleanKey[0] == '!' || cleanKey[0] == '+' || cleanKey[0] == '^' || cleanKey[0] == '=') {
```

**Code Location**: `config.go`, `applyMergeStrategies` function, lines ~1936-1946

#### 3. Duplicate Key Prioritization in `processKeys`
**Change**: Added logic to identify keys with explicit prefixes and prioritize them when duplicate keys exist (one with prefix, one without).

**Implementation**:
- First pass: Identify all keys with explicit prefixes and map them to their base keys
- Second pass: When processing keys, skip duplicate keys without prefixes if an explicit prefix key exists for the same base key
- This ensures that `"+exclude_patterns"` takes precedence over `"exclude_patterns"` if both exist

**Code Location**: `config.go`, `processKeys` function, lines ~2383-2465

### Behavior:

- **Quoted keys with prefixes**: `"+exclude_patterns"` is correctly recognized as merge strategy
- **Unquoted keys with prefixes**: `+exclude_patterns` continues to work as before
- **Duplicate keys**: Keys with explicit prefixes take precedence over plain keys
- **Both quote types**: Supports both double quotes (`"`) and single quotes (`'`)

### Test Cases:

- `TestMergeStrategiesInInheritanceChains/merge_strategy_(+)_in_inheritance_chain`: Validates merge strategy with quoted `"+exclude_patterns"` key
- `TestMergeStrategiesInSequentialFiles/merge_strategy_(+)_appends_in_sequential_files`: Validates merge strategy with quoted keys in sequential files
- `TestMergeStrategiesInSequentialFiles/prepend_strategy_(^)_prepends_in_sequential_files`: Validates prepend strategy with quoted keys

### Validation:

- ✅ Quoted YAML keys with merge strategy prefixes are correctly recognized
- ✅ Merge strategies work correctly with quoted keys in inheritance chains
- ✅ Merge strategies work correctly with quoted keys in sequential files
- ✅ Duplicate keys are handled correctly (explicit prefix keys take precedence)
- ✅ All tests pass with quoted key syntax

**Code Markers**: `extractStrategy`, `hasExplicitPrefix`, `processKeys`, `quoted YAML keys`, `cleanKey`, `explicitPrefixKeys`

**Cross-References**: [ARCH:CFG_005], [REQ:CFG_005], [REQ:CONFIGURATION], [IMPL:CFG_MIXED_MODE_MERGE_FIX]

## 26. Merge and Prepend Strategy Precedence Fix [IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX] [ARCH:CFG_001] [ARCH:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION]

### Decision: Fix `applyMerge()` and `applyPrepend()` to respect earlier file precedence for scalar/precedence fields
**Rationale:**
- Bug fix: `+` and `^` prefixes on scalar/precedence fields (e.g., `+archive_dir_path`, `^archive_dir_path`) were bypassing precedence checks
- Per CFG-001 and implementation decision [IMPL:CFG_MERGE_BEHAVIOR_REGISTRY], explicit prefixes should "work as normal, but precedence still applies for sequential files"
- `applyMerge()` and `applyPrepend()` were designed for arrays and fell back to `setConfigField()` for scalars without checking precedence
- This violated CFG-001 requirement that earlier files take precedence over later files

### Implementation Approach:
- Updated `applyMerge()` and `applyPrepend()` function signatures to accept `inheritContext`, `defaultCfg`, and `explicitlySetFields` parameters
- Added precedence checking logic for scalar/precedence fields before falling back to `setConfigField()`
- Updated `applyMergeOperation()` to pass new parameters to `applyMerge()` and `applyPrepend()`
- Updated tests to expect correct behavior (earlier file precedence)

### Key Changes:

#### 1. Updated `applyMerge()` Function Signature
**Change**: Added parameters for precedence checking: `inheritContext bool`, `defaultCfg *Config`, `explicitlySetFields map[string]bool`

**Before**: `func applyMerge(result *Config, key string, value interface{}, dstValue interface{}) error`
**After**: `func applyMerge(result *Config, key string, value interface{}, dstValue interface{}, inheritContext bool, defaultCfg *Config, explicitlySetFields map[string]bool) error`

**Code Location**: `config.go`, `applyMerge` function, line ~2563

#### 2. Added Precedence Check in `applyMerge()` for Scalar Fields
**Change**: When value is not an array, check if field is `MergeBehaviorPrecedence` and respect earlier file precedence before setting value.

**Before**: Fell back to `setConfigField(result, key, value)` without precedence check
**After**: Checks precedence similar to `applyOverride()` - if earlier file set the field, preserve its value

**Code Location**: `config.go`, `applyMerge` function, lines ~2598-2615

#### 3. Updated `applyPrepend()` Function Signature
**Change**: Added parameters for precedence checking: `inheritContext bool`, `defaultCfg *Config`, `explicitlySetFields map[string]bool`

**Before**: `func applyPrepend(result *Config, key string, value interface{}, dstValue interface{}) error`
**After**: `func applyPrepend(result *Config, key string, value interface{}, dstValue interface{}, inheritContext bool, defaultCfg *Config, explicitlySetFields map[string]bool) error`

**Code Location**: `config.go`, `applyPrepend` function, line ~2733

#### 4. Added Precedence Check in `applyPrepend()` for Scalar Fields
**Change**: When value is not an array, check if field is `MergeBehaviorPrecedence` and respect earlier file precedence before setting value.

**Before**: Fell back to `setConfigField(result, key, value)` without precedence check
**After**: Checks precedence similar to `applyOverride()` - if earlier file set the field, preserve its value

**Code Location**: `config.go`, `applyPrepend` function, lines ~2746-2763

#### 5. Updated `applyMergeOperation()` Calls
**Change**: Updated calls to `applyMerge()` and `applyPrepend()` to pass new parameters.

**Before**: 
```go
case "merge":
    return applyMerge(result, key, operation.value, dstValue)
case "prepend":
    return applyPrepend(result, key, operation.value, dstValue)
```

**After**:
```go
case "merge":
    return applyMerge(result, key, operation.value, dstValue, inheritContext, defaultCfg, explicitlySetFields)
case "prepend":
    return applyPrepend(result, key, operation.value, dstValue, inheritContext, defaultCfg, explicitlySetFields)
```

**Code Location**: `config.go`, `applyMergeOperation` function, lines ~2514-2516

#### 6. Updated Test Expectations
**Change**: Updated `TestAllStrategiesWithPrecedenceFields` tests to expect earlier file precedence for `+` and `^` prefixes.

**Before**: Tests accepted either file1 or file2 value (masking the bug)
**After**: Tests expect file1's value to be preserved (correct CFG-001 behavior)

**Code Location**: `config_test.go`, `TestAllStrategiesWithPrecedenceFields` function, lines ~5617-5665

### Behavior:

**Before Fix:**
- `+archive_dir_path` in file2: Overrode file1's value ❌
- `^archive_dir_path` in file2: Overrode file1's value ❌
- `!archive_dir_path` in file2: Correctly respected file1's value ✅
- `archive_dir_path` (no prefix) in file2: Correctly respected file1's value ✅

**After Fix:**
- `+archive_dir_path` in file2: Respects file1's value ✅
- `^archive_dir_path` in file2: Respects file1's value ✅
- `!archive_dir_path` in file2: Respects file1's value ✅
- `archive_dir_path` (no prefix) in file2: Respects file1's value ✅

### Test Coverage:
- `TestAllStrategiesWithPrecedenceFields/merge_strategy (+)` - Validates `+` prefix respects precedence
- `TestAllStrategiesWithPrecedenceFields/prepend_strategy (^)` - Validates `^` prefix respects precedence
- `TestAllStrategiesWithPrecedenceFields/replace_strategy (!)` - Validates `!` prefix respects precedence (already working)
- `TestAllStrategiesWithPrecedenceFields/no prefix` - Validates no prefix respects precedence (already working)

### Code Markers:
- `applyMerge`, `applyPrepend`, `applyMergeOperation`
- `// CFG-001 + CFG-005: For scalar/precedence fields with +/^ prefix, respect earlier file precedence`
- `getFieldMergeBehavior`, `MergeBehaviorPrecedence`

### Cross-References:
- [ARCH:CFG_001] - Configuration Discovery Architecture
- [ARCH:CFG_005] - Layered Configuration Inheritance Architecture
- [REQ:CFG_001] - Configuration Discovery (earlier file precedence)
- [REQ:CONFIGURATION] - Configuration Management
- [IMPL:CFG_MERGE_BEHAVIOR_REGISTRY] - Field-Level Merge Behavior Registry
- [IMPL:CFG_MIXED_MODE_MERGE_FIX] - Mixed-Mode Merge Strategy Fix

## 27. Priority 1 Configuration Merge Tests Implementation [IMPL:TEST_CFG_005_P1] [ARCH:CFG_005] [REQ:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION]

**Date**: 2025-12-11
**Status**: Implemented
**Priority**: P1 (Important)

### Decision: Implement comprehensive Priority 1 tests for configuration merge edge cases and inheritance scenarios
**Rationale:**
- Validates important edge cases for configuration inheritance and merging
- Ensures relative path resolution, home directory expansion, and error handling work correctly
- Tests cover multiple inheritance sources, deep chains, type mismatches, and null value handling
- Provides comprehensive coverage beyond Priority 0 critical tests

### Implementation Approach:
- **TestMultipleInheritanceSources_REQ_CFG_005**: Tests child inheriting from multiple parent files (parent1.yml, parent2.yml)
  - Verifies precedence fields (child overrides parents in inheritance chains)
  - Verifies accumulate fields merge from all sources including defaults
- **TestRelativePathInheritance_REQ_CFG_005**: Tests relative path resolution (../base.yml, ./sibling.yml)
  - Verifies paths resolve correctly relative to config file location
  - Tests both parent directory (../) and sibling (./) path resolution
- **TestHomeDirectoryExpansion_REQ_CFG_005**: Tests home directory expansion in inherit paths (~/.bkpdir-base.yml)
  - Verifies ~ expansion works in inheritance chain resolution
  - Tests loading base config from user's home directory
- **TestMissingInheritanceFile_REQ_CFG_005**: Tests error handling when inheritance file is missing
  - Verifies graceful error handling or skip behavior
  - Ensures child config still loads when parent is missing
- **TestInvalidYAMLHandling_REQ_CFG_005**: Tests invalid YAML in one file
  - Verifies first file loads successfully, second file skipped with error logged
  - Ensures no data loss when one file has errors
- **TestDeepInheritanceChain_REQ_CFG_005**: Tests very long inheritance chain (12 files)
  - Currently skipped due to known implementation bug where only last level is processed
  - Documents expected behavior: all levels should have patterns merged
- **TestTypeMismatchHandling_REQ_CFG_005**: Tests type mismatches (array vs string)
  - Verifies graceful handling when field types don't match
  - Tests first file array, second file string scenario
- **TestNullValueHandling_REQ_CFG_005**: Tests nil/null value handling
  - Verifies null values treated as not set
  - Tests that first file's values preserved when second file has null
- **TestWhitespaceStringHandling_REQ_CFG_005**: Tests whitespace-only string handling
  - Verifies whitespace strings handled correctly (treated as empty or preserved)
  - Tests precedence behavior with whitespace values

### Test Expectations:
- All tests account for default merging behavior (defaults included in accumulate fields)
- Inheritance chains: child overrides parent for precedence fields
- Sequential files: first file merges with defaults, subsequent files respect precedence
- Tests use semantic token references `[REQ:CFG_005]` in test names and comments

### Code Markers:
- `config_test.go`: TestMultipleInheritanceSources_REQ_CFG_005, TestRelativePathInheritance_REQ_CFG_005, TestHomeDirectoryExpansion_REQ_CFG_005, TestMissingInheritanceFile_REQ_CFG_005, TestInvalidYAMLHandling_REQ_CFG_005, TestDeepInheritanceChain_REQ_CFG_005, TestTypeMismatchHandling_REQ_CFG_005, TestNullValueHandling_REQ_CFG_005, TestWhitespaceStringHandling_REQ_CFG_005

**Cross-References**: [ARCH:CFG_005], [REQ:CFG_005], [REQ:CFG_001], [REQ:CONFIGURATION], [IMPL:CFG_INHERITANCE_PATH_RESOLUTION]

## 28. Inheritance Path Resolution Fix [IMPL:CFG_INHERITANCE_PATH_RESOLUTION] [ARCH:CFG_005] [REQ:CFG_005] [REQ:CONFIGURATION]

**Date**: 2025-12-11
**Status**: Implemented
**Priority**: P1 (Important)

### Decision: Fix inheritance chain path resolution to correctly handle relative paths and home directory expansion
**Rationale:**
- Bug fix: `buildChainRecursive` was passing `resolvedPath` (file path) instead of `filepath.Dir(resolvedPath)` (directory) as basePath
- Bug fix: `resolvePath` was not expanding `~` before resolving relative paths
- Bug fix: `resolvePath` was not handling directory paths correctly (was calling `filepath.Dir()` on directories)
- These bugs prevented inheritance chains from resolving parent files correctly, especially with relative paths and home directory expansion

### Implementation Approach:
- **Fixed buildChainRecursive basePath**: Changed line 2289 from `resolvedPath` to `filepath.Dir(resolvedPath)` when passing basePath to recursive calls
  - Ensures parent file paths are resolved relative to the directory containing the child file
- **Added ExpandPath method**: Added `ExpandPath` method to `defaultPathResolver` to handle `~` expansion
  - Expands `~/` to user's home directory before path resolution
  - Handles environment variable expansion
- **Fixed resolvePath directory handling**: Updated `resolvePath` to:
  - Call `ExpandPath` first to expand `~` before checking if path is absolute
  - Check if `basePath` is a directory or file using `os.Stat()`
  - Use `basePath` directly if it's a directory, or `filepath.Dir(basePath)` if it's a file
  - Handle case where `basePath` doesn't exist (assume it's a directory from `filepath.Dir()`)

### Code Changes:
- `config.go` line 2289: `buildChainRecursive(parentPath, filepath.Dir(resolvedPath), pathResolver, chain)`
- `config.go` line 2212-2240: Updated `resolvePath` to expand `~` and handle directories correctly
- `config.go` line 2241-2253: Added `ExpandPath` method to `defaultPathResolver`

### Test Coverage:
- TestConfigInheritance: Now passes with relative path `base.yml` resolving correctly
- TestRelativePathInheritance_REQ_CFG_005: Tests `../base/base.yml` and `./sibling.yml` resolution
- TestHomeDirectoryExpansion_REQ_CFG_005: Tests `~/.bkpdir-base.yml` expansion in inherit paths

**Cross-References**: [ARCH:CFG_005], [REQ:CFG_005], [REQ:CONFIGURATION], [IMPL:TEST_CFG_005_P1]

## 42. Mixed Sequential and Inheritance File Processing [IMPL:CFG_MIXED_SEQUENTIAL_INHERITANCE] [ARCH:CFG_005] [REQ:CFG_005] [REQ:CFG_001]

### Decision: Implement field tracking to preserve sequential file values when processing inheritance chains that come after sequential files
**Rationale:**
- Requirement: When a sequential file (no inheritance) is processed first, followed by an inheritance chain, the sequential file's precedence fields should be preserved [REQ:CFG_001]
- Requirement: Within inheritance chains, child files should override parent files (normal inheritance behavior) [REQ:CFG_005]
- Problem: Without field tracking, inheritance chain files would override sequential file values, violating CFG-001 precedence rules
- Solution: Track fields from sequential files (single-file chains) in `explicitlySetFields`, but do NOT track fields from true inheritance chain files, allowing child files to override parent files while preserving sequential file values

### Implementation Approach:
- **Field Tracking Strategy**:
  - Track fields in `explicitlySetFields` only from sequential files (single-file chains, `len(chain.files) == 1`)
  - Do NOT track fields from true inheritance chain files (`len(chain.files) > 1`) to allow child files to override parent files
  - When processing inheritance chain files, check `explicitlySetFields` to preserve sequential file values for precedence fields
- **Merge Function Updates**:
  - Updated `mergeBasicSettings`, `mergeGitSettings`, `mergeFileBackupSettings` to check `explicitlySetFields` even when `inheritContext=true` (inheritance chains)
  - Updated `applyOverride` to check `explicitlySetFields` when `inheritContext=true` to preserve sequential file values
  - Within inheritance chains, fields can still override each other (normal inheritance), but cannot override sequential file values
- **Single-File Chain Detection**:
  - Added `isSingleFileChain` flag to detect when `buildChain` returns a single-file chain (file without inheritance)
  - Single-file chains are processed through the inheritance chain path but have their fields tracked
  - This handles the case where `buildChain` succeeds for files without `inherit` field (returns single-file chain)

### Code Changes:
- `config.go` line 1789: Added `isSingleFileChain := len(chain.files) == 1` detection
- `config.go` line 1833-1843: Track fields from single-file chains in `explicitlySetFields`
- `config.go` line 675-692: Updated `mergeBasicSettings` to check `explicitlySetFields` when `inheritContext=true`
- `config.go` line 728-747: Updated `mergeBasicSettings` for `include_git_info` to check `explicitlySetFields`
- `config.go` line 806-837: Updated `mergeGitSettings` to check `explicitlySetFields` when `inheritContext=true`
- `config.go` line 928-940: Updated `mergeFileBackupSettings` to check `explicitlySetFields` when `inheritContext=true`
- `config.go` line 2649-2681: Updated `applyOverride` to check `explicitlySetFields` when `inheritContext=true`

### Test Coverage:
- `TestMixedSequentialAndInheritance`: Verifies that sequential file precedence fields are preserved when processing inheritance chains
  - Sequential file sets `archive_dir_path`, `include_git_info`, `skip_broken_symlinks` (precedence fields)
  - Inheritance chain (base.yml + inherited.yml) tries to override these fields
  - Expected: Sequential file values are preserved, accumulate fields (`exclude_patterns`) merge correctly

### Behavior:
- **Sequential files processed first**: Fields are tracked in `explicitlySetFields`
- **Inheritance chains processed after sequential files**: 
  - Precedence fields from sequential files are preserved (checked via `explicitlySetFields`)
  - Accumulate fields from sequential files and inheritance chain merge correctly
- **Within inheritance chains**: Child files can override parent files (normal inheritance behavior)
- **Single-file chains**: Fields are tracked to preserve values when inheritance chains are processed later

**Code Markers**: `isSingleFileChain`, `explicitlySetFields`, `mergeBasicSettings`, `mergeGitSettings`, `mergeFileBackupSettings`, `applyOverride`, `// CFG-001: Track fields explicitly set in this file for later precedence checks`

**Cross-References**: [ARCH:CFG_005], [REQ:CFG_005], [REQ:CFG_001], [REQ:CONFIGURATION]

## 34. Unicode and Special Character Handling Tests [IMPL:TEST_UNICODE_HANDLING] [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]

### Decision: Comprehensive test coverage for Unicode and special character handling in configuration
**Rationale:**
- Ensures configuration system correctly handles Unicode characters in paths and patterns
- Validates that special characters in config file paths don't break loading
- Critical for internationalization and real-world usage scenarios
- Part of comprehensive configuration merge test plan

### Implementation Approach:
- **TestUnicodeHandling**: Tests Unicode characters (émojis, 特殊字符) in configuration values
  - Tests Unicode in `archive_dir_path` field
  - Tests Unicode in `exclude_patterns` array field
  - Verifies all characters are preserved correctly after loading
  - Validates merge behavior with Unicode patterns (CFG-005: array fields default to merge)
  
- **TestSpecialCharactersInPaths**: Tests special characters in config file paths
  - Tests config file paths with spaces: `/path/with spaces/.bkpdir.yml`
  - Tests config file paths with Unicode: `/path/with-特殊/.bkpdir.yml`
  - Tests config file paths with both spaces and Unicode
  - Verifies LoadConfig correctly loads files regardless of path characters

### Test Structure:
```go
// TestUnicodeHandling tests Unicode and special characters in configuration values
// [REQ:CONFIGURATION] [REQ:CFG_005] Validates Unicode character preservation
func TestUnicodeHandling(t *testing.T) {
    // Creates config with Unicode in archive_dir_path and exclude_patterns
    // Verifies all Unicode characters preserved correctly
    // Validates merge behavior with defaults (CFG-005)
}

// TestSpecialCharactersInPaths tests special characters in config file paths
// [REQ:CONFIGURATION] [REQ:CFG_001] Validates config file loading with special paths
func TestSpecialCharactersInPaths(t *testing.T) {
    // Tests spaces in path
    // Tests Unicode in path
    // Tests both spaces and Unicode in path
    // Verifies LoadConfig works correctly
}
```

### Test Coverage:
- Unicode characters in `archive_dir_path`: `/path/with/émojis/🚀/and/特殊字符`
- Unicode patterns in `exclude_patterns`: `["*.文件", "测试/*", "unicode-文件.log", "émoji-🚀.tmp"]`
- Config file paths with spaces: `path with spaces/.bkpdir.yml`
- Config file paths with Unicode: `path-with-特殊/.bkpdir.yml`
- Config file paths with both: `path with 特殊 chars/.bkpdir.yml`
- Merge behavior validation: Unicode patterns merge correctly with defaults per CFG-005

### Behavior:
- **Unicode in values**: All Unicode characters preserved exactly as specified
- **Unicode in patterns**: Unicode patterns work correctly in exclude_patterns
- **Special chars in paths**: Config files load correctly regardless of path characters
- **Merge behavior**: Unicode patterns merge with defaults per CFG-005 (array fields default to merge)

**Code Markers**: `config_test.go` line ~6470-6600, `TestUnicodeHandling`, `TestSpecialCharactersInPaths`, `createTestConfigFileWithData`

**Cross-References**: [ARCH:TESTING_STRATEGY], [REQ:CONFIGURATION], [REQ:CFG_005], [REQ:CFG_001]

## 35. Empty String Handling Tests [IMPL:TEST_EMPTY_STRING_HANDLING] [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005]

### Decision: Comprehensive test coverage for empty string handling in configuration merging
**Rationale:**
- Ensures configuration system correctly handles empty strings in merging scenarios
- Validates whether empty strings are treated as zero values or explicit empty values
- Critical for understanding precedence behavior when empty strings are explicitly set
- Part of comprehensive configuration merge test plan

### Implementation Approach:
- **TestEmptyStringHandling**: Tests empty string handling in various merge scenarios
  - First file empty string, second file value: Tests if empty string is treated as zero value (allows override) or explicit value (preserves empty)
  - First file value, second file empty string: Tests CFG-001 precedence (first file value preserved)
  - Both files empty string: Tests that empty string is preserved when both files have it

### Test Structure:
```go
// TestEmptyStringHandling tests empty string handling in configuration merging
// [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005] Validates empty string treatment
func TestEmptyStringHandling(t *testing.T) {
    // Test: first file empty, second file value
    // Test: first file value, second file empty
    // Test: both files empty
}
```

### Test Coverage:
- Empty string in first file, value in second file: Verifies behavior (zero value vs explicit empty)
- Value in first file, empty string in second file: Verifies CFG-001 precedence (first file preserved)
- Empty string in both files: Verifies empty string is preserved
- Merge behavior validation: Empty strings handled correctly per CFG-001 and CFG-005

### Behavior:
- **Empty string as zero value**: If treated as zero value, later files can override
- **Empty string as explicit value**: If treated as explicit, CFG-001 precedence applies (first file preserved)
- **Precedence with empty strings**: First file's value (even if empty) takes precedence per CFG-001
- **Both files empty**: Result is empty string

**Code Markers**: `config_test.go` line ~6643-6750, `TestEmptyStringHandling`, `createTestConfigFileWithData`, `isZeroValue`

**Cross-References**: [ARCH:TESTING_STRATEGY], [REQ:CONFIGURATION], [REQ:CFG_001], [REQ:CFG_005]

## 36. Prepend Strategy Ordering Tests [IMPL:TEST_PREPEND_ORDERING] [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]

### Decision: Comprehensive test coverage for prepend strategy ordering in configuration merging
**Rationale:**
- Ensures prepend strategy (`^` prefix) correctly maintains order (new values before existing values)
- Validates prepend behavior in both inheritance chains and sequential file processing
- Critical for understanding merge strategy ordering semantics
- Part of comprehensive configuration merge test plan

### Implementation Approach:
- **TestPrependStrategyOrdering**: Tests prepend strategy ordering in various scenarios
  - Inheritance chain: Parent has patterns, child prepends - verifies child patterns come before parent patterns
  - Sequential files: First file has patterns, second file prepends - verifies second file patterns come before first file patterns
  - With defaults: Prepend with only defaults - verifies prepended values come before defaults
  - Multiple prepends: Multiple prepend operations - verifies CFG-001 precedence (earlier files take precedence)

### Test Structure:
```go
// TestPrependStrategyOrdering tests prepend strategy ordering in configuration merging
// [REQ:CONFIGURATION] [REQ:CFG_005] Validates prepend strategy maintains correct order
func TestPrependStrategyOrdering(t *testing.T) {
    // Test: prepend in inheritance chain
    // Test: prepend in sequential files
    // Test: prepend with defaults
    // Test: multiple prepend operations
}
```

### Test Coverage:
- Inheritance chain prepend: Parent `["parent1", "parent2"]`, child `^["child1", "child2"]` → `["child1", "child2", "parent1", "parent2"]`
- Sequential files prepend: First file patterns, second file prepends → new values before existing
- Prepend with defaults: Prepend patterns come before default patterns
- Multiple prepend operations: CFG-001 precedence applies (earlier prepends preserved)

### Behavior:
- **Prepend ordering**: New values (source) are placed before existing values (destination)
- **Inheritance chains**: Child prepended values come before parent values
- **Sequential files**: Later file prepended values come before earlier file values (but CFG-001 precedence may apply)
- **With defaults**: Prepended values come before default values
- **Multiple prepends**: CFG-001 precedence ensures earlier file prepends are preserved

**Code Markers**: `config_test.go` line ~6753-6850, `TestPrependStrategyOrdering`, `createTestConfigFileWithData`, `ArrayPrependStrategy`, `applyPrepend`

**Cross-References**: [ARCH:TESTING_STRATEGY], [REQ:CONFIGURATION], [REQ:CFG_005], [REQ:CFG_001]

## 37. Default Strategy Edge Cases Tests [IMPL:TEST_DEFAULT_STRATEGY_EDGES] [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]

### Decision: Comprehensive test coverage for default strategy (`=` prefix) edge cases
**Rationale:**
- Ensures default strategy only applies when destination is zero value
- Validates that default strategy does NOT apply when destination is non-zero or equals default
- Critical for understanding default strategy semantics
- Part of comprehensive configuration merge test plan

### Implementation Approach:
- **TestDefaultStrategyEdgeCases**: Tests default strategy in various edge case scenarios
  - String field zero value: Empty string → default strategy should apply
  - String field non-zero: Non-empty string → default strategy should NOT apply
  - String field equals default: Value equals default but not zero → default strategy should NOT apply
  - Array field zero value: Empty slice → default strategy should apply
  - Array field non-zero: Non-empty slice → default strategy should NOT apply
  - Bool field zero value: false → default strategy should apply
  - Bool field non-zero: true → default strategy should NOT apply

### Test Structure:
```go
// TestDefaultStrategyEdgeCases tests default strategy edge cases
// [REQ:CONFIGURATION] [REQ:CFG_005] Validates default strategy only applies when destination is zero value
func TestDefaultStrategyEdgeCases(t *testing.T) {
    // Test: default strategy when field is zero value
    // Test: default strategy when field is non-zero
    // Test: default strategy when field equals default
    // Test: default strategy with array field (zero and non-zero)
    // Test: default strategy with bool field (zero and non-zero)
}
```

### Test Coverage:
- String field zero value: `=archive_dir_path: "/custom/path"` when destination is `""` → applies
- String field non-zero: `=archive_dir_path: "/custom/path"` when destination is `"/existing/path"` → does NOT apply
- String field equals default: `=archive_dir_path: "/custom/path"` when destination is `"../.bkpdir"` → does NOT apply (not zero, even if equals default)
- Array field zero value: `=exclude_patterns: ["custom1"]` when destination is `[]` → **Important**: Empty array in first file merges with defaults first (CFG-005), so by the time default strategy is evaluated, destination is `[".git/", "vendor/"]` (non-zero), and default strategy does NOT apply. Result: defaults preserved, not custom values.
- Array field non-zero: `=exclude_patterns: ["custom1"]` when destination is `["existing1"]` → does NOT apply
- Bool field zero value: `=include_git_info: true` when destination is `false` → applies
- Bool field non-zero: `=include_git_info: false` when destination is `true` → does NOT apply

### Behavior:
- **Zero value check**: Default strategy uses `isZeroValue()` to check if destination is zero value
- **Zero value types**: Empty string `""`, empty slice `[]`, `false` for bool, `0` for int, `nil`
- **Non-zero values**: Any non-zero value prevents default strategy from applying
- **Equals default but not zero**: Even if value equals default (e.g., `"../.bkpdir"`), if it's not zero value, default strategy does NOT apply
- **Array field special case**: When first file has empty array `[]`, it merges with defaults first (CFG-005), so by the time default strategy is evaluated, destination is non-zero (contains defaults), and default strategy does NOT apply
- **Core principle**: Default strategy only applies when destination is zero value at the time of evaluation, regardless of whether it equals the default

**Code Markers**: `config_test.go` line ~6937-7155, `TestDefaultStrategyEdgeCases`, `createTestConfigFileWithData`, `applyDefault`, `isZeroValue`, `DefaultValueStrategy`

**Cross-References**: [ARCH:TESTING_STRATEGY], [REQ:CONFIGURATION], [REQ:CFG_005]

## N. List Command Limit Implementation [IMPL:LIST_LIMIT] [ARCH:LIST_LIMIT] [REQ:LIST_LIMIT]

### Decision: Add --limit flag to list commands with default value of 10
**Rationale:**
- Simple command-line flag provides immediate control
- Default of 10 prevents overwhelming output
- Applies to both directory archives and file backups consistently
- Maintains backward compatibility (existing behavior when limit not specified)

### Implementation Approach:
- Add `--limit` flag (short: `-n`) to `listCmd()` in `main.go`
- Add `--limit` flag to `--list` flag handling for file backups
- Default value: 10 (show newest 10 files)
- Special handling: `--limit 0` or `--all` flag shows all files
- Apply limit after sorting in `ListArchivesEnhanced()` and `ListFileBackupsEnhanced()`
- Limit logic: `if limit > 0 && len(archives) > limit { archives = archives[:limit] }`

**Code Structure:**
- `main.go`: Global variable `listLimit` (default 10) added to persistent flags
- `listCmd()`: Command inherits `--limit` flag from persistent flags
- `handleListCommand()`: Passes `listLimit` to `ListArchivesEnhanced()`
- `ListArchivesEnhanced()`: Accepts limit parameter, applies after sorting (line ~1157-1160)
- `handleListFileBackupsCommand()`: Passes `listLimit` to `ListFileBackupsEnhanced()`
- `ListFileBackupsEnhanced()`: Accepts limit parameter, applies after sorting (backup.go line ~392-396)
- `HandleListArchives()`: CommandHandler applies default limit of 10
- `HandleListFileBackups()`: CommandHandler applies default limit of 10

**Code Markers**: 
- `main.go` line ~41: `listLimit` global variable declaration
- `main.go` line ~387: Persistent flag `--limit` definition
- `main.go` line ~1069-1080: `listCmd()` function
- `main.go` line ~780-807: `handleListCommand()` function
- `main.go` line ~1128-1186: `ListArchivesEnhanced()` function with limit parameter
- `main.go` line ~1188-1219: `handleListFileBackupsCommand()` function
- `backup.go` line ~372-419: `ListFileBackupsEnhanced()` function with limit parameter
- `main_test.go` line ~1301-1610: Comprehensive test suite `TestListArchivesEnhanced_WithLimit`
- `backup_test.go` line ~964-1210: Comprehensive test suite `TestListFileBackupsEnhanced_WithLimit`

**Cross-References**: [ARCH:LIST_LIMIT], [REQ:LIST_LIMIT], [REQ:USABILITY]

## N+1. Diff Command Implementation [IMPL:DIFF_COMMAND] [ARCH:DIFF_COMMAND] [REQ:DIFF_COMMAND]

### Decision: Implement diff command with archive state reconstruction and comparison engine
**Rationale:**
- Reconstructs effective state by applying incremental on top of full archive
- Provides accurate change detection by comparing against reconstructed state
- Reuses comparison logic for duplicate prevention consistency
- Integrates with existing CLI framework and output formatting system
- Supports context cancellation for long-running operations
- Provides user-friendly error messages for edge cases

### Implementation Approach:
- **Archive State Reconstruction**: Function to reconstruct effective state from full + most recent incremental archive
  - Load most recent full archive snapshot
  - Load most recent incremental archive snapshot (if exists)
  - Apply incremental changes on top of full archive state
  - Return reconstructed directory snapshot
  - Handle case where no archives exist with user-friendly message
- **Directory Comparison**: Compare current directory against reconstructed state
  - Create snapshot of current directory
  - Compare snapshots using existing comparison infrastructure
  - Identify added, modified, and deleted files
  - Respect exclude patterns from configuration
- **Change Reporting**: Format and display changes
  - Use configurable format strings for output
  - Display added files, modified files, deleted files
  - Handle case where no changes exist
- **CLI Command Integration**: Add `diff` command to CLI framework
  - Add `diffCmd()` function in `main.go`
  - Support context cancellation with periodic checks
  - Add "diff" to known commands list in auto-detection logic
  - Respect exclude patterns from configuration
  - Improved error handling for "no archives" case

### Implementation Details:
- **Archive Selection Strategy**: All archive selection functions use name-based sorting instead of file modification times
  - Archive names include timestamps in ISO 8601 format (`YYYY-MM-DDTHHmmss`) which are alphabetically sortable
  - Functions filter appropriate archives (full vs incremental), sort alphabetically by name, and return the last archive (most recent when sorted ascending)
  - More reliable than file system modification times, consistent with archive naming pattern, simpler implementation (no file system stat calls needed)
  - Affected functions: `FindMostRecentArchive()`, `findLatestFullArchive()`, `findLatestIncrementalArchive()`
- **Context Support**: Command handler checks for context cancellation before diff calculation and before printing results
- **Error Handling**: Improved handling of "no archives found" case - displays user-friendly message using `PrintNoArchivesFound()` instead of generic error
- **Auto-Detection Integration**: Added "diff" to `knownCommands` list in `executeWithAutoDetection()` to prevent path auto-detection from intercepting the command
- **Testing**: Comprehensive test coverage including:
  - Unit tests for archive state reconstruction (full only, full + incremental)
  - Unit tests for diff calculation (added, modified, deleted, no changes)
  - Integration test for end-to-end diff command
  - Edge case test for no archives scenario
  - Format and print method tests

**Code Structure:**
- `comparison.go`: `ReconstructArchiveState()` function to reconstruct state from full + incremental
- `comparison.go`: `CalculateDiff()` function to compare current directory against reconstructed state
- `comparison.go`: `DiffResult` struct to hold change information (added, modified, deleted files)
- `comparison.go`: `findLatestIncrementalArchive()` helper function - sorts by name (timestamps are alphabetically sortable)
- `comparison.go`: `FindMostRecentArchive()` function - sorts by name instead of modification time
- `archive.go`: `findLatestFullArchive()` function - sorts by name instead of modification time
- `main.go`: `diffCmd()` function for CLI command with context support
- `main.go`: Added "diff" to `knownCommands` list in `executeWithAutoDetection()` (line ~250)
- `formatter.go`: `FormatDiffResult()` and `PrintDiffResult()` methods for diff output
- `config.go`: Format string configuration for diff output (`FormatDiffNoChanges`, `FormatDiffChanges`, `FormatDiffAdded`, `FormatDiffModified`, `FormatDiffDeleted`)
- `diff_command_test.go`: Comprehensive test suite with semantic token references

**Archive Selection Strategy:**
- **Name-Based Sorting**: All archive selection functions (`FindMostRecentArchive`, `findLatestFullArchive`, `findLatestIncrementalArchive`) use name-based sorting instead of file modification times
- **Rationale**: Archive names include timestamps in ISO 8601 format (`YYYY-MM-DDTHHmmss`) which are alphabetically sortable, making name-based sorting more reliable and consistent with archive naming conventions
- **Implementation**: Filter appropriate archives (full vs incremental), sort alphabetically by name, return the last archive (most recent when sorted ascending)
- **Benefits**: More reliable than file system modification times, consistent with archive naming pattern, simpler implementation (no file system stat calls needed)

**Code Markers**:
- `comparison.go` line ~68-101: `FindMostRecentArchive(archiveDir string) (string, error)` - Finds most recent full archive by sorting names
- `comparison.go` line ~177-226: `findLatestIncrementalArchive(archiveDir string, baseFullArchive *Archive) (*Archive, error)` - Finds most recent incremental archive by sorting names
- `comparison.go` line ~231-330: `ReconstructArchiveState(archiveDir string) (*DirectorySnapshot, error)` - Reconstructs state from full + incremental
- `comparison.go` line ~333-409: `CalculateDiff(cwd string, reconstructedState *DirectorySnapshot, excludePatterns []string) (*DiffResult, error)` - Calculates differences
- `comparison.go` line ~169-174: `type DiffResult struct { Added []string; Modified []string; Deleted []string }` - Holds change information
- `archive.go` line ~869-904: `findLatestFullArchive(archiveDir string) (*Archive, error)` - Finds most recent full archive by sorting names
- `main.go` line ~1090-1143: `diffCmd() *cobra.Command` - CLI command implementation with context support
- `main.go` line ~250: Added "diff" to `knownCommands` list
- `formatter.go` line ~1752-1771: `FormatDiffResult(diff *DiffResult) string` - Format diff output
- `formatter.go` line ~1775-1786: `PrintDiffResult(diff *DiffResult)` - Print diff output
- `diff_command_test.go`: Comprehensive test suite with 8 test functions

**Cross-References**: [ARCH:DIFF_COMMAND], [REQ:DIFF_COMMAND], [REQ:INCREMENTAL_DUPLICATE_PREVENTION], [ARCH:DIRECTORY_COMPARISON], [REQ:OUTPUT_FORMATTING], [REQ:CONTEXT_SUPPORT], [ARCH:CLI_COMMANDS]

## N+2. Incremental Archive Duplicate Prevention Implementation [IMPL:INCREMENTAL_DUPLICATE_PREVENTION] [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [REQ:INCREMENTAL_DUPLICATE_PREVENTION]

### Decision: Reuse diff command analysis to prevent duplicate incremental archives
**Rationale:**
- Ensures consistency by using same comparison logic as diff command
- Prevents code duplication
- Compares against reconstructed state (full + most recent incremental) for accuracy
- Provides clear user feedback when archive creation is skipped
- Archive selection uses name-based sorting (consistent with diff command and archive naming conventions)

### Implementation Approach:
- **Reuse Diff Logic**: Leverage `CalculateDiff()` function from diff command implementation
- **State Reconstruction**: Use `ReconstructArchiveState()` to get effective state
- **Archive Selection**: Use `findLatestFullArchive()` which uses name-based sorting (consistent with diff command)
- **Change Detection**: Use diff result to determine if changes exist
- **Skip Creation**: Skip incremental archive creation if no changes detected
- **User Feedback**: Display appropriate message using format strings

**Code Structure:**
- `archive.go`: Modify `createIncrementalArchive()` to use diff analysis before creating archive
- `archive.go`: Add check for changes using `CalculateDiff()` function
- `archive.go`: Skip archive creation if `diffResult` shows no changes
- `formatter.go`: Add format method for skip message (`FormatIncrementalSkippedNoChanges()`, `PrintIncrementalSkippedNoChanges()`)
- `config.go`: Add format string configuration for skip message (`FormatIncrementalSkippedNoChanges`)

**Code Markers**:
- `archive.go` line ~695-738: `createIncrementalArchive()` - Integrated diff analysis before archive creation
  - Line ~697: Calls `ReconstructArchiveState()` to get effective state (full + most recent incremental)
  - Line ~704: Calls `findLatestFullArchive()` for fallback path (uses name-based sorting)
  - Line ~744: Calls `findLatestFullArchive()` to get base full archive for naming (uses name-based sorting)
  - Line ~720: Calls `CalculateDiff()` to detect changes, passing exclude patterns
  - Line ~726-729: Conditional skip when no changes detected, displays skip message
  - Line ~700-717: Fallback to old behavior when reconstruction fails (backward compatibility)
- `archive.go` line ~869-904: `findLatestFullArchive()` - Finds most recent full archive by sorting names (used by incremental archive creation)
- `formatter.go` line ~1785-1799: `FormatIncrementalSkippedNoChanges()` and `PrintIncrementalSkippedNoChanges()` - Format and print skip message
- `config.go` line ~153: `FormatIncrementalSkippedNoChanges` - Format string configuration field
- `config.go` line ~336: Default format string: "No changes detected since last incremental archive. Skipping archive creation.\n"
- `incremental_duplicate_prevention_test.go`: Comprehensive test suite with 4 test functions covering all scenarios

**Implementation Status**: ✅ Complete
- All code changes implemented and tested
- All tests passing (4 test functions, all scenarios covered)
- Exclude patterns respected via `archiveConfig.GetExcludePatterns()` passed to `CalculateDiff()`
- Edge cases handled: no incremental (compares against full), no archives (falls back to old behavior)
- Consistent with diff command behavior (reuses same comparison logic)

**Cross-References**: [ARCH:INCREMENTAL_DUPLICATE_PREVENTION], [REQ:INCREMENTAL_DUPLICATE_PREVENTION], [REQ:DIFF_COMMAND], [IMPL:DIFF_COMMAND], [REQ:OUTPUT_FORMATTING]

## EXTRACT-008 Documentation Migration & Preservation [IMPL:EXTRACT_008_DOC_MIGRATION] [ARCH:EXTRACT_008_INTERDEP] [REQ:EXTRACT_008_INTERDEP_MAPPING]

### Decision: Preserve working-plan content by migrating it directly into STDD core documents and decommission the working-plan artifact.

**Implementation Steps:**
1. Record the requirement token `[REQ:EXTRACT_008_INTERDEP_MAPPING]` in `stdd/requirements.md` (done).
2. Record the architecture decision `[ARCH:EXTRACT_008_INTERDEP]` in `stdd/architecture-decisions.md` (done).
3. Record this implementation decision `[IMPL:EXTRACT_008_DOC_MIGRATION]` here to document the migration approach and rationale.
4. Create the canonical deliverable `docs/package-interdependency-mapping.md` describing package relationships, interface contracts, and example integrations (outcome document—authoring to be performed as follow-up task).
5. Remove any supplemental preservation files (they are not canonical) and remove or decommission `working-plan-extract-008.md` once migration is complete.

**Rationale:**
STDD requires authoritative records to live in the requirements, architecture, and implementation documents for traceability. Working-plan or supplemental files create duplication and risk drift; migrating content into the canonical STDD files preserves traceability.

**Validation:**
- Confirm tokens appear in `stdd/semantic-tokens.md` and trace to entries in requirements, architecture, and implementation documents.
- Confirm `docs/package-interdependency-mapping.md` exists and links back to these tokens.

## 19. Semantic Token Coverage Audit [IMPL:TOKEN_COVERAGE_AUDIT] [ARCH:TOKEN_SYSTEM] [REQ:DOC_016]

### Decision: Verify that each small module and its tests carry semantic tokens linking them to their requirements, architecture decisions, and implementation anchors.
**Rationale:**
- Demonstrates that the AI-first toolchain can rely on tokens to trace code/tests back to intent.
- Highlights gaps (e.g., `formatter.go`/`formatter_test.go`) and fixes them before broader refactors.
- Provides a repeatable checklist for future module audits.

### Implementation Approach:
- **Archive module (`archive.go`, `archive_test.go`)**: Documented `[REQ:FILE_BACKUP]`, `[ARCH:ARCHIVE_FORMAT]`, `[ARCH:PROCESSING_PATTERNS]`, `[ARCH:CONTEXT_SUPPORT]`, `[IMPL:ZIP_FORMAT]`, `[IMPL:PROCESSING_PATTERNS]`, `[IMPL:CONTEXT_OPS]`, `[IMPL:INCREMENTAL_DUPLICATE_PREVENTION]` with matching test references (e.g., `TestCreateFullArchive_REQ_FILE_BACKUP`).
- **Backup module (`backup.go`, `backup_test.go`)**: Annotated `[REQ:FILE_BACKUP]`, `[ARCH:RESOURCE_MANAGEMENT]`, `[IMPL:ATOMIC_OPS]`, additional `[REQ:LIST_LIMIT]` tokens for list helpers so the requirement coverage remains visible.
- **Errors module (`errors.go`, `errors_test.go`)**: Structured errors are marked with `[REQ:ERROR_HANDLING]`, `[ARCH:ERROR_HANDLING]`, `[IMPL:STRUCTURED_ERRORS]`; tests mirror these tags.
- **Exclusion module (`exclude.go`, `exclude_test.go`)**: Pattern helpers keep `[REQ:CONFIGURATION]`, `[ARCH:EXCLUSION_PATTERNS]`, `[ARCH:PACKAGE_EXTRACTION]`, `[IMPL:EXCLUSION_PATTERNS]` references across adapters and wrappers.
- **Comparison module (`comparison.go`, `comparison_test.go`)**: Snapshot utilities reference `[ARCH:DIRECTORY_COMPARISON]`, `[ARCH:DIFF_COMMAND]`, `[IMPL:DIRECTORY_COMPARISON]`, `[IMPL:DIFF_COMMAND]`, `[REQ:DIFF_COMMAND]` for diffing features.
- **Formatter module (`formatter.go`, `formatter_test.go`)**: Added `[REQ:OUTPUT_FORMATTING]`, `[ARCH:OUTPUT_FORMATTING]`, `[IMPL:DUAL_FORMATTING]` anchors so the formatting requirement is traceable from both implementation and tests; helper functions continue to reference `[REQ:CUSTOMIZABLE_FORMAT_STRINGS]`, `[IMPL:DIFF_COMMAND]`, and `[IMPL:INCREMENTAL_DUPLICATE_PREVENTION]` for contextual behavior.
- **Git module (`git.go`, `git_test.go`)**: Maintains `[REQ:GIT_INTEGRATION]`, `[ARCH:GIT_INTEGRATION]`, `[IMPL:GIT_CLI]`, and tests now restate them for repository detection and naming flows.

### Documentation:
- Capture this audit in `stdd/semantic-tokens.md` by registering `[IMPL:TOKEN_COVERAGE_AUDIT]` with the contexts above.
- Point future audits to this decision to reuse the pattern and ensure coverage plans explicitly mention `REQ:DOC_016` and `ARCH:TOKEN_SYSTEM`.

**Code Markers**: `archive.go`, `backup.go`, `errors.go`, `comparison.go`, `exclude.go`, `formatter.go`, `git.go` and their `_test.go` files that now consistently mention the appropriate tokens.

## 50. Configuration Hierarchy Preservation Fix [IMPL:CFG_HIERARCHY_PRESERVATION] [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION] [REQ:CFG_001]

### Decision: Fix configuration hierarchy to preserve values from earlier files when later files don't set those fields
**Rationale:**
- Bug fix: When a local config file exists but doesn't set a field (e.g., `archive_dir_path`), the system was falling back to compiled defaults instead of preserving values from home directory config files
- The expected behavior is that compiled defaults should only apply if the value is not set anywhere in the configuration hierarchy
- If a home config file sets `archive_dir_path = /Users/fareed/.bkpdir`, and a local config file exists but doesn't set `archive_dir_path`, the home config value should be preserved, not replaced with the compiled default
- This ensures proper configuration hierarchy: home config → local config → compiled defaults, with each level only applying when the previous level doesn't set a value

### Implementation Approach:
- **Updated `mergeBasicSettings`** (lines 682-689 in `config.go`): For sequential files (`inheritContext=false`), only override a field if:
  - The earlier file didn't explicitly set it (`!explicitlySetByEarlier`)
  - The destination value equals the default (`!dstDiffersFromDefault`)
  - The source file explicitly sets it (`explicitlySetInSrc`)
  - This ensures that if an earlier file set a field (even if not explicitly tracked), we preserve it, but if dstValue equals default (meaning earlier file didn't set it), allow later file to set it

- **Updated `applyOverride`** (lines 2647-2652 in `config.go`): Changed the preservation logic to only preserve if:
  - The field was explicitly set by an earlier file (`wasSetByEarlierFile`), OR
  - The destination value differs from the default (`dstDiffersFromDefault`)
  - Removed the `earlierFilesProcessed` check that was preserving fields when ANY field was set by earlier files, not just the specific field
  - This ensures that if an earlier file didn't set a field (dstValue equals default), later files can set it

**Key Changes:**
1. `mergeBasicSettings` - Updated sequential file logic to check `dstDiffersFromDefault` in addition to `explicitlySetByEarlier`
2. `applyOverride` - Simplified preservation logic to only check if the specific field was set or differs from default, not if any field was set

**Behavior:**
- **Sequential files** (`inheritContext=false`): 
  - If earlier file set a field (even if equals default), preserve it
  - If earlier file didn't set a field (dstValue equals default), allow later file to set it
  - This ensures proper hierarchy: home config values are preserved when local config doesn't set them
- **Inheritance chains** (`inheritContext=true`): Behavior unchanged - child configs override parent configs

**Test Case:**
- Home config (`~/.bkpdir.yml`): `archive_dir_path: /Users/fareed/.bkpdir`
- Local config (`./.bkpdir.yml`): `archive_dir_path` is not set (commented out or missing)
- Expected: Home config value (`/Users/fareed/.bkpdir`) is preserved
- Before fix: Compiled default (`../.bkpdir`) was incorrectly used
- After fix: Home config value (`/Users/fareed/.bkpdir`) is correctly preserved

**Code Markers**: `mergeBasicSettings`, `applyOverride`, `dstDiffersFromDefault`, `explicitlySetByEarlier`, `explicitlySetInSrc`, `// CFG-001: Earlier files take precedence`

**Cross-References**: [ARCH:CONFIG_SYSTEM], [REQ:CONFIGURATION], [REQ:CFG_001], [IMPL:CFG_PRECEDENCE_FIX]


## STDD Visualization Data Pipeline [IMPL:STDD_VIS_DATA_PIPELINE] [ARCH:STDD_VIS_FLOW] [REQ:STDD_VIS]

### Decision: Build a repeatable data pipeline that derives a token graph from STDD documents and representative code/tests to feed visualization artifacts.

**Rationale:**
- Keeps visuals synchronized with the canonical token registry.
- Enables regeneration when tokens or docs change.
- Provides a small, reviewable JSON/CSV snapshot that reviewers can validate.

**Implementation Approach (planning, no code yet):**
- Input sources: `stdd/semantic-tokens.md`, `stdd/requirements.md`, `stdd/architecture-decisions.md`, `stdd/implementation-decisions.md`, and sampled code/test anchors (grep for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`).
- Normalize nodes by type (REQ/ARCH/IMPL/TEST/CODE) and emit edges for cross-references found in docs and inline token comments; capture token metadata (type, status, sample refs) to drive token prominence in visuals.
- Export artifacts: `docs/data/stdd-trace.json` (graph) and `docs/data/stdd-samples.csv` (sample chains) with token-first styling hints (e.g., color map by token type, prominence weights).
- Validation hooks: schema check for nodes/edges, count consistency vs token registry, and spot-check two sample chains (e.g., `[REQ:CFG_005]` and `[REQ:LIST_LIMIT]`), confirming tokens appear at every hop.
- Module boundary: `DataExtraction` outputs normalized graph once; downstream modules consume it without re-parsing docs, preserving token-centric styling data.

**Cross-References**: [ARCH:STDD_VIS_FLOW], [REQ:STDD_VIS], [REQ:MODULE_VALIDATION]


## STDD Visualization Asset Plan [IMPL:STDD_VIS_ASSETS] [ARCH:STDD_VIS_FLOW] [REQ:STDD_VIS]

### Decision: Produce two complementary artifacts (layered flow + timeline/animation) sourced from the data pipeline, with embedded token callouts and legends.

**Rationale:**
- Combines static clarity (layered flow) with process storytelling (timeline/animation).
- Ensures artifacts stay anchored to real tokens and code/test references.

**Implementation Approach (planning, no code yet):**
- **Layered Flow Diagram**: Use the token graph to render a swimlane/Sankey-style view: lanes for REQ/ARCH/IMPL/TEST/CODE; tokens are visually dominant (size/color/labels); callouts for example chains; legend explaining token roles and color map.
- **Timeline/Animation**: Stepwise frames showing a requirement moving through architecture, implementation, tests, and code; token badges stay visible per step so the token remains the protagonist. Output as lightweight GIF/MP4 or animated SVG.
- **Asset locations**: `docs/images/stdd-trace-flow.svg` (or .png) and `docs/images/stdd-trace-timeline.gif/mp4`; link from STDD docs with captions describing token chains.
- **Validation**: Visual review checklist centered on token prominence (label clarity, color map alignment, edge correctness), alignment with two sample chains from the data snapshot, and confirmation that assets load in docs.
- Module boundaries: `Visualization` (static) and `AnimationTimeline` (dynamic) consume the same data snapshot; `Integration` links assets into docs and records validation notes with token focus.

**Cross-References**: [ARCH:STDD_VIS_FLOW], [REQ:STDD_VIS], [REQ:MODULE_VALIDATION]
