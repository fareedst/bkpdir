# Architecture and Specification Requirements Traceability

This document maps code components to their corresponding architecture requirements and specification requirements.

> **Note**: For testing requirements and architecture, see [Testing Requirements](testing.md).

## Code Quality and Linting Requirements

### Linting Standards
**Implementation**: `Makefile`, `.revive.toml`
**Specification Requirements**:
- **Linter**: Uses `revive` for Go code linting
  - Spec: "All code must pass revive linter checks"
  - Configuration: `.revive.toml` file in project root
  - Command: `make lint` runs revive linter
  - Rules: Standard Go best practices with custom configurations
- **Error Handling**: All errors must be properly handled
  - Spec: "No unhandled errors allowed"
  - All `fmt.Printf`, `fmt.Fprintf` return values must be checked
  - All file operations must handle errors appropriately
- **Code Style**: Consistent formatting and naming conventions
  - Spec: "Follow Go standard formatting and naming"
  - Function names must be descriptive and follow Go conventions
  - Variable names must be clear and meaningful
  - Comments must follow Go documentation standards

**Example Usage**:
```bash
# Run linter
make lint

# Build with linting
make build
```

### Resource Management Requirements
**Implementation**: `archive.go` - `ResourceManager`
**Specification Requirements**:
- **Resource Cleanup**: All temporary resources must be cleaned up
  - Spec: "No temporary files or directories should remain after operations"
  - Implementation: `ResourceManager` struct with automatic cleanup
  - Thread-safe: Uses mutex for concurrent access
  - Error-resilient: Continues cleanup even if individual operations fail
- **Atomic Operations**: Archive operations must be atomic
  - Spec: "Archive creation must be atomic to prevent corruption"
  - Implementation: Temporary files with atomic rename operations
  - Cleanup: Temporary files registered for automatic cleanup
- **Panic Recovery**: Operations must recover from panics
  - Spec: "Unexpected panics must not leave resources uncleaned"
  - Implementation: Defer functions with panic recovery
  - Logging: Panic information logged to stderr

**Example Usage**:
```go
// Create resource manager
rm := NewResourceManager()
defer rm.Cleanup()

// Register temporary resources
rm.AddTempFile("/tmp/archive.tmp")
rm.AddTempDir("/tmp/archive_work")
```

### Enhanced Error Handling Requirements
**Implementation**: `archive.go` - Enhanced error handling
**Specification Requirements**:
- **Structured Errors**: Use `ArchiveError` for consistent error handling
  - Spec: "All archive operations must return structured errors with status codes"
  - Fields:
    - `Message`: Human-readable error description
    - `StatusCode`: Numeric exit code for application
  - Implementation: Implements Go's `error` interface
  - Usage: Allows callers to extract both message and status code

**Example Usage**:
```go
// Create structured error
err := NewArchiveError("file not found", cfg.StatusFileNotFound)

// Check for ArchiveError type
if archiveErr, ok := err.(*ArchiveError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", archiveErr.Message)
    os.Exit(archiveErr.StatusCode)
}
```

### Output Formatting
**Implementation**: `formatter.go`
**Specification Requirements**:
- `NewOutputFormatter(cfg *Config) *OutputFormatter`
  - Spec: "Creates new output formatter with configuration"
  - Input: `cfg *Config` - Configuration containing format strings
  - Output: `*OutputFormatter` - New formatter instance
  - Behavior: Creates formatter with reference to configuration for format strings
  - Error Cases: None

- `FormatCreatedArchive(path string) string`
  - Spec: "Formats successful archive creation message using printf-style format"
  - Input: `path string` - Path to created archive
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatCreatedArchive` format string with path parameter
  - Error Cases: None (falls back to safe default on invalid format)

- `FormatIdenticalArchive(path string) string`
  - Spec: "Formats identical archive message using printf-style format"
  - Input: `path string` - Path to existing archive
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatIdenticalArchive` format string with path parameter
  - Error Cases: None (falls back to safe default on invalid format)

- `FormatListArchive(path, creationTime string) string`
  - Spec: "Formats archive listing entry using printf-style format"
  - Input:
    - `path string` - Path to archive file
    - `creationTime string` - Creation timestamp
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatListArchive` format string with path and time parameters
  - Error Cases: None (falls back to safe default on invalid format)

- `FormatConfigValue(name, value, source string) string`
  - Spec: "Formats configuration value display using printf-style format"
  - Input:
    - `name string` - Configuration parameter name
    - `value string` - Configuration value
    - `source string` - Source file or "default"
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatConfigValue` format string with name, value, and source parameters
  - Error Cases: None (falls back to safe default on invalid format)

- `FormatDryRunArchive(path string) string`
  - Spec: "Formats dry-run archive message using printf-style format"
  - Input: `path string` - Path that would be created
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatDryRunArchive` format string with path parameter
  - Error Cases: None (falls back to safe default on invalid format)

- `FormatError(message string) string`
  - Spec: "Formats error message using printf-style format"
  - Input: `message string` - Error message
  - Output: `string` - Formatted message
  - Behavior: Uses `cfg.FormatError` format string with message parameter
  - Error Cases: None (falls back to safe default on invalid format)

**Example Usage**:
```go
// Create output formatter
formatter := NewOutputFormatter(cfg)

// Format messages
createdMsg := formatter.FormatCreatedArchive("/path/to/archive")
errorMsg := formatter.FormatError("File not found")

// Print messages directly
formatter.PrintCreatedArchive("/path/to/archive")
formatter.PrintError("File not found")

// Use in archive operations
if archiveCreated {
    formatter.PrintCreatedArchive(archivePath)
} else {
    formatter.PrintIdenticalArchive(existingArchivePath)
}
```

## Data Objects

### Config
**Implementation**: `config.go`
**Specification Requirements**:
- Configuration stored in YAML file `.bkpdir.yml` at root directory
- Default values used if file not present
- Fields:
  - `ArchiveDirPath`: `Config.ArchiveDirPath`
    - Spec: "Specifies where archives are stored"
    - Default: "../.bkpdir" relative to current directory
    - YAML key: "archive_dir_path"
  - `UseCurrentDirName`: `Config.UseCurrentDirName`
    - Spec: "Controls whether to include current directory name in archive path"
    - Default: true
    - YAML key: "use_current_dir_name"
  - `ExcludePatterns`: `Config.ExcludePatterns`
    - Spec: "List of patterns to exclude from archiving"
    - Default: [".git/", "vendor/"]
    - YAML key: "exclude_patterns"
  - `Verification`: `Config.Verification`
    - Spec: "Controls archive verification behavior"
  - `StatusCreatedArchive`: `Config.StatusCreatedArchive`
    - Spec: "Exit code when a new archive is successfully created"
    - Default: 0
    - YAML key: "status_created_archive"
  - `StatusFailedToCreateArchiveDirectory`: `Config.StatusFailedToCreateArchiveDirectory`
    - Spec: "Exit code when archive directory creation fails"
    - Default: 31
    - YAML key: "status_failed_to_create_archive_directory"
  - `StatusDirectoryIsIdenticalToExistingArchive`: `Config.StatusDirectoryIsIdenticalToExistingArchive`
    - Spec: "Exit code when directory is identical to most recent archive"
    - Default: 0
    - YAML key: "status_directory_is_identical_to_existing_archive"
  - `StatusDirectoryNotFound`: `Config.StatusDirectoryNotFound`
    - Spec: "Exit code when source directory does not exist"
    - Default: 20
    - YAML key: "status_directory_not_found"
  - `StatusInvalidDirectoryType`: `Config.StatusInvalidDirectoryType`
    - Spec: "Exit code when source path is not a directory"
    - Default: 21
    - YAML key: "status_invalid_directory_type"
  - `StatusPermissionDenied`: `Config.StatusPermissionDenied`
    - Spec: "Exit code when directory access is denied"
    - Default: 22
    - YAML key: "status_permission_denied"
  - `StatusDiskFull`: `Config.StatusDiskFull`
    - Spec: "Exit code when disk space is insufficient"
    - Default: 30
    - YAML key: "status_disk_full"
  - `StatusConfigError`: `Config.StatusConfigError`
    - Spec: "Exit code when configuration is invalid"
    - Default: 10
    - YAML key: "status_config_error"
  - `FormatCreatedArchive`: `Config.FormatCreatedArchive`
    - Spec: "Printf-style format string for successful archive creation messages"
    - Default: "Created archive: %s\n"
    - YAML key: "format_created_archive"
    - Supports: ANSI color codes and text formatting
  - `FormatIdenticalArchive`: `Config.FormatIdenticalArchive`
    - Spec: "Printf-style format string for identical directory messages"
    - Default: "Directory is identical to existing archive: %s\n"
    - YAML key: "format_identical_archive"
    - Supports: ANSI color codes and text formatting
  - `FormatListArchive`: `Config.FormatListArchive`
    - Spec: "Printf-style format string for archive listing entries"
    - Default: "%s (created: %s)\n"
    - YAML key: "format_list_archive"
    - Supports: ANSI color codes and text formatting
  - `FormatConfigValue`: `Config.FormatConfigValue`
    - Spec: "Printf-style format string for configuration value display"
    - Default: "%s: %s (source: %s)\n"
    - YAML key: "format_config_value"
    - Supports: ANSI color codes and text formatting
  - `FormatDryRunArchive`: `Config.FormatDryRunArchive`
    - Spec: "Printf-style format string for dry-run archive messages"
    - Default: "Would create archive: %s\n"
    - YAML key: "format_dry_run_archive"
    - Supports: ANSI color codes and text formatting
  - `FormatError`: `Config.FormatError`
    - Spec: "Printf-style format string for error messages"
    - Default: "Error: %s\n"
    - YAML key: "format_error"
    - Supports: ANSI color codes and text formatting
  - `TemplateCreatedArchive`: `Config.TemplateCreatedArchive`
    - Spec: "Template string for successful archive creation messages with named placeholders"
    - Default: "Created archive: %{path}\n"
    - YAML key: "template_created_archive"
    - Supports: Go text/template syntax and %{name} placeholders
  - `TemplateIdenticalArchive`: `Config.TemplateIdenticalArchive`
    - Spec: "Template string for identical directory messages with named placeholders"
    - Default: "Directory is identical to existing archive: %{path}\n"
    - YAML key: "template_identical_archive"
    - Supports: Go text/template syntax and %{name} placeholders
  - `TemplateListArchive`: `Config.TemplateListArchive`
    - Spec: "Template string for archive listing entries with named placeholders"
    - Default: "%{path} (created: %{creation_time})\n"
    - YAML key: "template_list_archive"
    - Supports: Go text/template syntax and %{name} placeholders
  - `TemplateConfigValue`: `Config.TemplateConfigValue`
    - Spec: "Template string for configuration value display with named placeholders"
    - Default: "%{name}: %{value} (source: %{source})\n"
    - YAML key: "template_config_value"
    - Supports: Go text/template syntax and %{name} placeholders
  - `TemplateDryRunArchive`: `Config.TemplateDryRunArchive`
    - Spec: "Template string for dry-run archive messages with named placeholders"
    - Default: "Would create archive: %{path}\n"
    - YAML key: "template_dry_run_archive"
    - Supports: Go text/template syntax and %{name} placeholders
  - `TemplateError`: `Config.TemplateError`
    - Spec: "Template string for error messages with named placeholders"
    - Default: "Error: %{message}\n"
    - YAML key: "template_error"
    - Supports: Go text/template syntax and %{name} placeholders
  - `PatternArchiveFilename`: `Config.PatternArchiveFilename`
    - Spec: "Named regex pattern for parsing archive filenames"
    - Default: "(?P<prefix>[^-]*)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\\.zip"
    - YAML key: "pattern_archive_filename"
    - Supports: Named capture groups for data extraction
  - `PatternConfigLine`: `Config.PatternConfigLine`
    - Spec: "Named regex pattern for parsing configuration display lines"
    - Default: "(?P<name>[^:]+):\\s*(?P<value>[^(]+)\\s*\\(source:\\s*(?P<source>[^)]+)\\)"
    - YAML key: "pattern_config_line"
    - Supports: Named capture groups for data extraction
  - `PatternTimestamp`: `Config.PatternTimestamp`
    - Spec: "Named regex pattern for parsing timestamps"
    - Default: "(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})\\s+(?P<hour>\\d{2}):(?P<minute>\\d{2}):(?P<second>\\d{2})"
    - YAML key: "pattern_timestamp"
    - Supports: Named capture groups for data extraction

**Example Usage**:
```go
// Load configuration from YAML or use defaults
cfg, err := LoadConfig(".")
if err != nil {
    log.Fatal(err)
}

// Access configuration values
archivePath := cfg.ArchiveDirPath
excludePatterns := cfg.ExcludePatterns

// Access status code configuration
createdArchiveStatus := cfg.StatusCreatedArchive
identicalDirectoryStatus := cfg.StatusDirectoryIsIdenticalToExistingArchive
directoryNotFoundStatus := cfg.StatusDirectoryNotFound

// Access format string configuration
createdFormat := cfg.FormatCreatedArchive
errorFormat := cfg.FormatError
```

### ConfigValue
**Implementation**: `config.go`
**Specification Requirements**:
- Fields:
  - `Name`: `ConfigValue.Name`
    - Spec: "Configuration parameter name"
  - `Value`: `ConfigValue.Value`
    - Spec: "Computed configuration value including defaults"
  - `Source`: `ConfigValue.Source`
    - Spec: "Source file path or 'default' for default values"

**Example Usage**:
```go
// Create configuration value entry
configValue := &ConfigValue{
    Name:   "archive_dir_path",
    Value:  "../.bkpdir",
    Source: "~/.bkpdir.yml",
}
```

### VerificationConfig
**Implementation**: `config.go`
**Specification Requirements**:
- Fields:
  - `VerifyOnCreate`: `VerificationConfig.VerifyOnCreate`
    - Spec: "Automatically verify archives after creation"
    - Default: false
    - YAML key: "verify_on_create"
  - `ChecksumAlgorithm`: `VerificationConfig.ChecksumAlgorithm`
    - Spec: "Algorithm used for checksums"
    - Default: "sha256"
    - YAML key: "checksum_algorithm"

**Example Usage**:
```go
// Create verification config with custom settings
verifCfg := &VerificationConfig{
    VerifyOnCreate: true,
    ChecksumAlgorithm: "sha256",
}

// Use in main config
cfg.Verification = verifCfg
```

### Archive
**Implementation**: `archive.go`
**Specification Requirements**:
- Fields:
  - `Name`: `Archive.Name`
    - Spec: "Archive naming format for Git/non-Git repositories"
  - `Path`: `Archive.Path`
    - Spec: "Full path to archive file"
  - `CreationTime`: `Archive.CreationTime`
    - Spec: "When the archive was created"
  - `IsIncremental`: `Archive.IsIncremental`
    - Spec: "Whether this is an incremental archive"
  - `GitBranch`: `Archive.GitBranch`
    - Spec: "Current git branch (if in git repo)"
  - `GitHash`: `Archive.GitHash`
    - Spec: "Current git commit hash (if in git repo)"
  - `Note`: `Archive.Note`
    - Spec: "Optional note for the archive"
  - `BaseArchive`: `Archive.BaseArchive`
    - Spec: "Name of base archive for incremental backups"
  - `VerificationStatus`: `Archive.VerificationStatus`
    - Spec: "Result of last verification"

**Example Usage**:
```go
// Create a new archive object
archive := &Archive{
    Name: "2024-03-20-15-30=main=abc123=backup.zip",
    Path: "/path/to/archive.zip",
    CreationTime: time.Now(),
    IsIncremental: false,
    GitBranch: "main",
    GitHash: "abc123",
    Note: "backup",
}
```

### VerificationStatus
**Implementation**: `verify.go`
**Specification Requirements**:
- Fields:
  - `VerifiedAt`: `VerificationStatus.VerifiedAt`
    - Spec: "When the archive was last verified"
  - `IsVerified`: `VerificationStatus.IsVerified`
    - Spec: "Whether the archive passed verification"
  - `HasChecksums`: `VerificationStatus.HasChecksums`
    - Spec: "Whether checksums were verified"
  - `Errors`: `VerificationStatus.Errors`
    - Spec: "List of verification errors"

### ArchiveError
**Implementation**: `archive.go`
**Specification Requirements**:
- **Structured Error Handling**: Provides consistent error reporting
  - Spec: "All archive operations return structured errors with status codes"
  - Fields:
    - `Message`: Human-readable error description
    - `StatusCode`: Numeric exit code for application
  - Implementation: Implements Go's `error` interface
  - Usage: Allows callers to extract both message and status code

**Example Usage**:
```go
// Create structured error
err := NewArchiveError("directory not found", cfg.StatusDirectoryNotFound)

// Check for ArchiveError type
if archiveErr, ok := err.(*ArchiveError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", archiveErr.Message)
    os.Exit(archiveErr.StatusCode)
}
```

### ResourceManager
**Implementation**: `archive.go`
**Specification Requirements**:
- **Resource Tracking**: Thread-safe tracking of temporary resources
  - Spec: "Track all temporary files and directories for cleanup"
  - Fields:
    - `tempFiles`: List of temporary files to clean up
    - `tempDirs`: List of temporary directories to clean up
    - `mutex`: Mutex for thread-safe access
  - Methods:
    - `AddTempFile()`: Register temporary file for cleanup
    - `AddTempDir()`: Register temporary directory for cleanup
    - `Cleanup()`: Remove all registered resources

**Example Usage**:
```go
// Create and use resource manager
rm := NewResourceManager()
defer rm.Cleanup()

rm.AddTempFile("/tmp/archive.tmp")
rm.AddTempDir("/tmp/archive_work")
```

### TemplateFormatter
**Implementation**: `formatter.go`
**Specification Requirements**:
- **Template-Based Formatting**: Provides centralized template-based formatting using named placeholders
  - Spec: "All standard output must support template-based formatting with named data extraction"
  - Fields:
    - `config`: *Config - Reference to configuration for template strings and regex patterns
  - Methods:
    - `NewTemplateFormatter(cfg *Config) *TemplateFormatter`: Creates new template formatter
    - `FormatWithTemplate(input, pattern, tmplStr string) (string, error)`: Applies text/template to named regex groups
    - `FormatWithPlaceholders(format string, data map[string]string) string`: Replaces %{name} placeholders
    - `TemplateCreatedArchive(path string) string`: Formats archive creation using template
    - `TemplateIdenticalArchive(path string) string`: Formats identical directory message using template
    - `TemplateListArchive(path, creationTime string) string`: Formats archive listing using template
    - `TemplateConfigValue(name, value, source string) string`: Formats configuration display using template
    - `TemplateDryRunArchive(path string) string`: Formats dry-run message using template
    - `TemplateError(message, operation string) string`: Formats error message using template

**Example Usage**:
```go
// Create template formatter with configuration
templateFormatter := NewTemplateFormatter(cfg)

// Format messages using templates
createdMsg := templateFormatter.TemplateCreatedArchive("/path/to/archive")
errorMsg := templateFormatter.TemplateError("Directory not found", "archive_creation")

// Use template formatting with regex extraction
archivePath := "project-2024-03-21-15-30=main=abc123=important.zip"
formatted, err := templateFormatter.FormatWithTemplate(
    archivePath,
    cfg.PatternArchiveFilename,
    "Archive of {{.prefix}} from {{.year}}-{{.month}}-{{.day}} ({{.branch}}@{{.hash}}) with note: {{.note}}",
)

// Use placeholder formatting
data := map[string]string{
    "prefix": "myproject",
    "year": "2024",
    "month": "03",
    "day": "21",
    "branch": "main",
    "hash": "abc123",
}
result := templateFormatter.FormatWithPlaceholders(
    "Directory %{prefix} archived on %{year}-%{month}-%{day} from %{branch}@%{hash}",
    data,
)
```

## Core Functions

### Configuration Management
**Implementation**: `config.go`
**Specification Requirements**:
- `DefaultConfig()`: `DefaultConfig()`
  - Spec: "Creates default configuration with specified defaults"
  - Input: None
  - Output: `*Config` - Returns a new Config with default values
  - Behavior: Creates a new Config struct with all fields set to their default values
  - Default Values:
    - ArchiveDirPath: "../.bkpdir"
    - UseCurrentDirName: true
    - ExcludePatterns: [".git/", "vendor/"]
    - Verification.VerifyOnCreate: false
    - Verification.ChecksumAlgorithm: "sha256"
    - StatusCreatedArchive: 0
    - StatusFailedToCreateArchiveDirectory: 31
    - StatusDirectoryIsIdenticalToExistingArchive: 0
    - StatusDirectoryNotFound: 20
    - StatusInvalidDirectoryType: 21
    - StatusPermissionDenied: 22
    - StatusDiskFull: 30
    - StatusConfigError: 10

- `LoadConfig()`: `LoadConfig(root string) (*Config, error)`
  - Spec: "Loads config from YAML or uses defaults"
  - Input: `root string` - Path to root directory containing .bkpdir.yml
  - Output: `(*Config, error)` - Returns config and any error encountered
  - Behavior:
    - Attempts to read .bkpdir.yml from root directory
    - If file exists, merges with default values
    - If file doesn't exist, returns default config
    - Returns error if file exists but is invalid YAML
  - Error Cases:
    - Invalid YAML format
    - Invalid configuration values
    - File system errors

- `GetConfigSearchPath()`: `GetConfigSearchPath() []string`
  - Spec: "Returns list of configuration file paths to search"
  - Input: None
  - Output: `[]string` - List of configuration file paths
  - Behavior:
    - Reads `BKPDIR_CONFIG` environment variable for configuration search path
    - Returns default path if environment variable not set: `./.bkpdir.yml:~/.bkpdir.yml`
    - Handles colon-separated path list parsing
    - Supports home directory expansion

- `DisplayConfig()`: `DisplayConfig() error`
  - Spec: "Displays computed configuration values and exits"
  - Input: None
  - Output: `error` - Any error encountered
  - Behavior:
    - Processes configuration files from `BKPDIR_CONFIG` environment variable
    - If `BKPDIR_CONFIG` not set, uses hard-coded default search path: `./.bkpdir.yml:~/.bkpdir.yml`
    - Shows each configuration value with name, computed value, and source file
    - Displays format: `name: value (source: source_file)`
    - Default values show source as "default"
    - Application exits after displaying values
  - Error Cases:
    - Configuration file read errors
    - Invalid YAML format
    - Environment variable parsing errors

**Example Usage**:
```go
// Create default configuration
cfg := DefaultConfig()

// Load configuration from YAML file
cfg, err := LoadConfig(".")
if err != nil {
    log.Fatal(err)
}

// Display configuration values and exit
err = DisplayConfig()
if err != nil {
    log.Fatal(err)
}
```

### Git Integration
**Implementation**: `git.go`
**Specification Requirements**:
- `IsGitRepository(dir string) bool`
  - Spec: "Checks if directory is a git repo"
  - Input: `dir string` - Path to directory to check
  - Output: `bool` - True if directory is a git repository
  - Behavior: Checks for .git directory or git configuration
  - Error Cases: None (returns false for any error)

- `GetGitBranch(dir string) string`
  - Spec: "Gets current branch name"
  - Input: `dir string` - Path to git repository
  - Output: `string` - Current branch name or empty string
  - Behavior: Executes git branch command to get current branch
  - Error Cases: Returns empty string if not a git repo or error occurs

- `GetGitShortHash(dir string) string`
  - Spec: "Gets current commit hash"
  - Input: `dir string` - Path to git repository
  - Output: `string` - Short commit hash or empty string
  - Behavior: Executes git rev-parse to get abbreviated commit hash
  - Error Cases: Returns empty string if not a git repo or error occurs

**Example Usage**:
```go
// Check if directory is a git repository
if IsGitRepository(".") {
    // Get git information
    branch := GetGitBranch(".")
    hash := GetGitShortHash(".")
    fmt.Printf("Git: %s@%s\n", branch, hash)
}
```

### File System Operations
**Implementation**: `archive.go`
**Specification Requirements**:
- `ShouldExcludeFile(filePath string, patterns []string) bool`
  - Spec: "Checks if file matches exclude patterns using doublestar glob matching"
  - Input:
    - `filePath string` - Path to file to check
    - `patterns []string` - List of glob patterns to match against
  - Output: `bool` - True if file should be excluded
  - Behavior: Matches file path against each pattern using doublestar glob matching
  - Error Cases: None (returns false for any error)

- `CopyFileWithContext(ctx context.Context, src, dst string) error`
  - Spec: "Context-aware file copying with cancellation support"
  - Input:
    - `ctx context.Context` - Context for cancellation
    - `src string` - Path to source file
    - `dst string` - Path to destination file
  - Output: `error` - Any error encountered
  - Behavior:
    - Creates destination directory if needed
    - Copies file with all permissions preserved
    - Preserves original file's modification time
    - Checks for cancellation at multiple points during operation
    - Returns context.Canceled if operation is cancelled
  - Error Cases:
    - Source file not found
    - Permission denied
    - Disk full
    - Context cancellation
    - Other file system errors

- `CompareDirectories(dir1, dir2 string) (bool, error)`
  - Spec: "Performs comparison of two directory trees"
  - Input:
    - `dir1 string` - Path to first directory
    - `dir2 string` - Path to second directory
  - Output: `(bool, error)` - True if directories are identical, error if comparison fails
  - Behavior:
    - Walks both directory trees
    - Compares file structure and content
    - Excludes files matching configured patterns
    - Returns false immediately if any difference found
  - Error Cases:
    - Directory not found
    - Permission denied
    - File system errors

**Example Usage**:
```go
// Check if file should be excluded
patterns := []string{".git/", "vendor/", "*.tmp"}
if ShouldExcludeFile("path/to/file.tmp", patterns) {
    fmt.Println("File should be excluded")
}

// Copy with context and timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err = CopyFileWithContext(ctx, "source.txt", "backup.txt")
if err == context.Canceled {
    log.Println("Copy operation was cancelled")
}

// Compare two directories
identical, err := CompareDirectories("source/", "backup/")
if err != nil {
    log.Fatal(err)
}
if identical {
    fmt.Println("Directories are identical")
}
```

### Enhanced Error Detection
**Implementation**: `archive.go`
**Specification Requirements**:
- `isDiskFullError(err error) bool`
  - Spec: "Enhanced disk space error detection"
  - Input: `err error` - Error to check
  - Output: `bool` - True if error indicates disk space issues
  - Behavior:
    - Checks error message for multiple disk space indicators
    - Indicators: "no space left", "disk full", "not enough space", "insufficient disk space", "device full", "quota exceeded", "file too large"
    - Case-insensitive matching
  - Error Cases: None (returns false for nil or unrelated errors)

**Example Usage**:
```go
// Check for disk space errors
if err := someFileOperation(); err != nil {
    if isDiskFullError(err) {
        return NewArchiveError("Disk full", cfg.StatusDiskFull)
    }
    return err
}
```

### Archive Management
**Implementation**: `archive.go`
**Specification Requirements**:
- `GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, isGit, isIncremental bool, baseName string) string`
  - Spec: "Generates archive filename according to specified format"
  - Input:
    - `prefix string` - Optional prefix for archive name
    - `timestamp string` - Creation timestamp
    - `gitBranch string` - Current git branch
    - `gitHash string` - Current git commit hash
    - `note string` - Optional note for archive
    - `isGit bool` - Whether in git repository
    - `isIncremental bool` - Whether this is an incremental archive
    - `baseName string` - Name of base archive for incremental backups
  - Output: `string` - Generated archive filename
  - Behavior: Constructs filename according to format:
    - Full: `[prefix-]timestamp[=branch=hash][=note].zip`
    - Incremental: `baseName_update=timestamp[=branch=hash][=note].zip`
  - Error Cases: None

- `ListArchives(archiveDir string) ([]Archive, error)`
  - Spec: "Gets all archives in directory"
  - Input: `archiveDir string` - Path to archive directory
  - Output: `([]Archive, error)` - List of archives and any error
  - Behavior:
    - Scans directory for .zip files
    - Creates Archive objects for each file
    - Includes verification status if available
    - Sorts archives by creation time (most recent first)
  - Error Cases:
    - Directory not found
    - Permission denied
    - Invalid archive files

- `CreateFullArchive(cfg *Config, note string, dryRun bool) error`
  - Spec: "Creates full archive with specified behavior"
  - Input:
    - `cfg *Config` - Configuration to use
    - `note string` - Optional note for archive
    - `dryRun bool` - Whether to simulate creation
  - Output: `error` - Any error encountered
  - Behavior:
    - Walks directory collecting files
    - Excludes files matching patterns
    - Creates ZIP archive
    - Optionally verifies after creation
    - Uses configured status codes for different exit conditions
    - Includes panic recovery for unexpected errors
  - Error Cases:
    - Invalid configuration (exits with `cfg.StatusConfigError`)
    - Directory not found (exits with `cfg.StatusDirectoryNotFound`)
    - Invalid directory type (exits with `cfg.StatusInvalidDirectoryType`)
    - Permission denied (exits with `cfg.StatusPermissionDenied`)
    - Disk full (exits with `cfg.StatusDiskFull`)
    - Failed to create archive directory (exits with `cfg.StatusFailedToCreateArchiveDirectory`)
    - Directory identical to existing archive (exits with `cfg.StatusDirectoryIsIdenticalToExistingArchive`)
    - Successful archive creation (exits with `cfg.StatusCreatedArchive`)

- `CreateFullArchiveWithCleanup(cfg *Config, note string, dryRun bool) error`
  - Spec: "Creates full archive with automatic resource cleanup"
  - Input: Same as CreateFullArchive
  - Output: `error` - Any error encountered
  - Behavior:
    - All CreateFullArchive functionality plus:
    - Automatic resource cleanup via ResourceManager
    - Atomic operations using temporary files
    - Cleanup on errors or panics
    - No temporary files left behind
  - Error Cases: Same as CreateFullArchive

- `CreateFullArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool) error`
  - Spec: "Context-aware full archive creation"
  - Input:
    - `ctx context.Context` - Context for cancellation
    - Other inputs same as CreateFullArchive
  - Output: `error` - Any error encountered
  - Behavior:
    - All CreateFullArchive functionality plus:
    - Context cancellation support
    - Cancellation checks at multiple points
    - Returns appropriate error on cancellation
  - Error Cases: Same as CreateFullArchive plus context cancellation

- `CreateFullArchiveWithContextAndCleanup(ctx context.Context, cfg *Config, note string, dryRun bool) error`
  - Spec: "Context-aware full archive creation with resource cleanup"
  - Input:
    - `ctx context.Context` - Context for cancellation
    - Other inputs same as CreateFullArchive
  - Output: `error` - Any error encountered
  - Behavior:
    - Combines all features of CreateFullArchiveWithCleanup and CreateFullArchiveWithContext
    - Context cancellation support with automatic cleanup
    - Atomic operations with cleanup on cancellation
    - Most robust archive creation function
  - Error Cases: Same as CreateFullArchive plus context cancellation

- `CreateIncrementalArchive(cfg *Config, note string, dryRun bool) error`
  - Spec: "Creates incremental archive based on last full archive"
  - Input:
    - `cfg *Config` - Configuration to use
    - `note string` - Optional note for archive
    - `dryRun bool` - Whether to simulate creation
  - Output: `error` - Any error encountered
  - Behavior:
    - Finds most recent full archive
    - Collects files modified since full archive
    - Creates incremental ZIP archive
    - Optionally verifies after creation
    - Uses configured status codes for different exit conditions
  - Error Cases:
    - No full archive found
    - Invalid configuration
    - File system errors
    - Archive creation errors
    - Verification failures

**Example Usage**:
```go
// Generate archive name
name := GenerateArchiveName("prefix", "2024-03-20-15-30", "main", "abc123", "backup", true, false, "")

// List all archives
archives, err := ListArchives("/path/to/archives")
if err != nil {
    log.Fatal(err)
}

// Create full archive (basic)
err = CreateFullArchive(cfg, "backup note", false)

// Create full archive with cleanup (recommended for production)
err = CreateFullArchiveWithCleanup(cfg, "backup note", false)

// Create full archive with context and cleanup (most robust)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
defer cancel()
err = CreateFullArchiveWithContextAndCleanup(ctx, cfg, "backup note", false)

// Create incremental archive
err = CreateIncrementalArchive(cfg, "incremental update", false)
if err != nil {
    log.Fatal(err)
}
```

### Archive Verification
**Implementation**: `verify.go`
**Specification Requirements**:
- `VerifyArchive(archivePath string) (*VerificationStatus, error)`
  - Spec: "Verifies archive structure and integrity"
  - Input: `archivePath string` - Path to archive file
  - Output: `(*VerificationStatus, error)` - Verification status and any error
  - Behavior:
    - Validates ZIP file structure
    - Checks for corruption
    - Verifies file headers
  - Error Cases:
    - File not found
    - Invalid ZIP format
    - Corrupted archive
    - Permission denied

- `VerifyChecksums(archivePath string) (*VerificationStatus, error)`
  - Spec: "Verifies file checksums against stored values"
  - Input: `archivePath string` - Path to archive file
  - Output: `(*VerificationStatus, error)` - Verification status and any error
  - Behavior:
    - Extracts files to temporary directory
    - Generates checksums for extracted files
    - Compares with stored checksums
  - Error Cases:
    - No checksums found
    - Checksum mismatch
    - Extraction errors
    - File system errors

- `GenerateChecksums(files []string, algorithm string) (map[string]string, error)`
  - Spec: "Generates checksums for files using specified algorithm"
  - Input:
    - `files []string` - List of file paths
    - `algorithm string` - Checksum algorithm to use
  - Output: `(map[string]string, error)` - Map of filename to checksum
  - Behavior:
    - Reads each file
    - Generates checksum using specified algorithm
    - Uses base filename as key
  - Error Cases:
    - Invalid algorithm
    - File read errors
    - Permission denied

- `StoreChecksums(archive *Archive, checksums map[string]string) error`
  - Spec: "Stores checksums in archive"
  - Input:
    - `archive *Archive` - Archive to store checksums in
    - `checksums map[string]string` - Map of filename to checksum
  - Output: `error` - Any error encountered
  - Behavior:
    - Creates .checksums file in archive
    - Stores checksums as JSON
  - Error Cases:
    - Invalid archive
    - Write errors
    - JSON encoding errors

- `ReadChecksums(archive *Archive) (map[string]string, error)`
  - Spec: "Reads checksums from archive"
  - Input: `archive *Archive` - Archive to read checksums from
  - Output: `(map[string]string, error)` - Map of filename to checksum
  - Behavior:
    - Reads .checksums file from archive
    - Parses JSON into map
  - Error Cases:
    - No checksums file
    - Invalid JSON
    - Read errors

**Example Usage**:
```go
// Verify archive structure
status, err := VerifyArchive("/path/to/archive.zip")
if err != nil {
    log.Fatal(err)
}

// Generate and store checksums
files := []string{"file1.txt", "file2.txt"}
checksums, err := GenerateChecksums(files, "sha256")
if err != nil {
    log.Fatal(err)
}

// Store checksums in archive
err = StoreChecksums(archive, checksums)
if err != nil {
    log.Fatal(err)
}

// Read and verify checksums
storedChecksums, err := ReadChecksums(archive)
if err != nil {
    log.Fatal(err)
}

// Verify checksums
status, err = VerifyChecksums("/path/to/archive.zip")
if err != nil {
    log.Fatal(err)
}
```

## Main Application Structure

### CLI Interface
**Implementation**: `main.go`
**Specification Requirements**:
- Uses `cobra` for command-line interface
- Global flags:
  - `--dry-run`: Implemented in `main.go`
    - Spec: "Show what would be done without creating archives"
    - Shows paths relative to current directory
  - `--config`: Implemented in `main.go`
    - Spec: "Display computed configuration values and exit"
    - Usage: `bkpdir --config`
- Commands:
  - `full [NOTE]`: `fullCmd()`
    - Spec: "Create full archive with optional note"
    - Usage: `bkpdir full [NOTE]`
    - Output: Shows archive path (relative to current directory) and creation time
    - When a new archive is created: Uses `FormatCreatedArchive` or `TemplateCreatedArchive` configuration
    - When directory is identical to existing archive: Uses `FormatIdenticalArchive` or `TemplateIdenticalArchive` configuration
    - Uses configured status codes for application exit
  - `inc [NOTE]`: `incCmd()`
    - Spec: "Create incremental archive with optional note"
    - Usage: `bkpdir inc [NOTE]`
    - Output: Shows archive path (relative to current directory) and creation time
    - Uses configured status codes for application exit
  - `list`: `listCmd()`
    - Spec: "List all archives"
    - Usage: `bkpdir list`
    - Shows archive names with verification status if available
    - Output formatting uses `FormatListArchive` or `TemplateListArchive` configuration
    - Shows [VERIFIED], [FAILED], or [UNVERIFIED] status indicators
  - `verify [ARCHIVE_NAME]`: `verifyCmd()`
    - Spec: "Verify archive integrity"
    - Usage: `bkpdir verify [ARCHIVE_NAME]`
    - Flags:
      - `--checksum`: Include checksum verification
    - Performs archive structure and integrity verification
    - With --checksum flag: verifies file contents against stored checksums
    - Stores verification results for display in list command
    - Uses configured status codes for application exit

### Workflow Implementation
**Implementation**: `archive.go` and `verify.go`
**Specification Requirements**:
- Full archive workflow: `CreateFullArchive()` and enhanced variants
  - Spec: "Compresses current directory tree into ZIP archive"
  - Steps:
    1. Load config
    2. Initialize output formatter and template formatter with configuration
    3. Validate source directory exists and is accessible
    4. Convert directory path to relative path if needed
    5. Compare directory with most recent archive
       - If identical, use formatter to display existing archive message and exit with `cfg.StatusDirectoryIsIdenticalToExistingArchive`
       - If different, proceed with archive creation
    6. Generate archive name using directory name
    7. Create archive directory structure
    8. Walk directory collecting files (excluding patterns, skipping directories)
    9. Get git info if applicable
    10. Create ZIP archive (or simulate in dry-run) with atomic operations
    11. Use formatter to display success message
    12. Clean up temporary resources
    13. If config.Verification.VerifyOnCreate is true, verify archive
    14. Exit with `cfg.StatusCreatedArchive` on successful archive creation
  - Error handling uses configured status codes and formatted error messages:
    - Directory not found: exit with `cfg.StatusDirectoryNotFound`
    - Invalid directory type: exit with `cfg.StatusInvalidDirectoryType`
    - Permission denied: exit with `cfg.StatusPermissionDenied`
    - Disk full: exit with `cfg.StatusDiskFull`
    - Failed to create archive directory: exit with `cfg.StatusFailedToCreateArchiveDirectory`
    - Configuration error: exit with `cfg.StatusConfigError`
  - Enhanced error handling:
    - Panic recovery with proper logging
    - Context cancellation support
    - Automatic resource cleanup on all error paths
    - All error messages use configurable format strings or templates

- Incremental archive workflow: `CreateIncrementalArchive()`
  - Spec: "Creates update archive containing only changed files"
  - Steps:
    1. Load config
    2. Initialize output formatter and template formatter with configuration
    3. Determine archive directory path
    4. Find most recent full archive
    5. Walk directory and collect files modified since the full archive (skipping directories)
    6. Get git info if applicable
    7. Generate incremental archive name
    8. Create zip archive (or simulate in dry-run)
    9. Use formatter to display success message
    10. Clean up temporary resources
    11. If config.Verification.VerifyOnCreate is true, verify archive

- Archive listing workflow: `ListArchives()`
  - Spec: "Displays all archives with verification status"
  - Steps:
    1. Load config
    2. Initialize output formatter and template formatter with configuration
    3. Get archive directory path
    4. List all .zip files
    5. Extract archive information and notes using regex patterns if template formatting is enabled
    6. Sort archives by creation time
    7. Use formatter to display archive names with verification status if available
    8. Shows [VERIFIED], [FAILED], or [UNVERIFIED] status indicators

- Archive verification workflow: `VerifyArchive()` and `VerifyChecksums()`
  - Spec: "Verifies archive integrity and optionally checksums"
  - Steps:
    1. Load config
    2. Initialize output formatter and template formatter with configuration
    3. Get archive path
    4. Verify archive structure and integrity
    5. If --checksum flag is set, verify file checksums
    6. Store verification status
    7. Use formatter to display verification results

- Configuration display workflow: `DisplayConfig()`
  - Spec: "Displays computed configuration values and exits"
  - Steps:
    1. Read `BKPDIR_CONFIG` environment variable or use default search path
    2. Process configuration files in order with precedence rules
    3. Merge configuration values with defaults
    4. Track source file for each configuration value
    5. Initialize output formatter and template formatter with configuration
    6. Use formatter to display each configuration value with configurable format
    7. Exit application after display

### Utility Functions
**Implementation**: Various files
**Specification Requirements**:
- File exclusion: `ShouldExcludeFile()` in `archive.go`
  - Spec: "Uses doublestar glob pattern matching"
- Archive naming: `GenerateArchiveName()` in `archive.go`
  - Spec: "Follows specified naming format"
- Checksum generation: `GenerateChecksums()` in `verify.go`
  - Spec: "Uses SHA-256 by default"
- Resource management: `ResourceManager` in `archive.go`
  - Spec: "Tracks and cleans up temporary resources automatically"
- Error handling: `ArchiveError` and related functions in `archive.go`
  - Spec: "Provides structured error handling with status codes"
- Directory comparison: `CompareDirectories()` in `archive.go`
  - Spec: "Performs tree comparison to detect changes"
  - Input: Source directory path and most recent archive path
  - Output: Boolean indicating if directory is identical to archive
  - Behavior: Compares directory tree with archive contents to detect changes

## Build and Development Requirements

### Build System
**Implementation**: `Makefile`
**Specification Requirements**:
- **Linting**: `make lint` command
  - Spec: "Run revive linter on all Go code"
  - Must pass before code can be committed
  - Uses `.revive.toml` configuration file
- **Testing**: `make test` command
  - Spec: "Run all unit tests with verbose output"
  - Must pass before code can be committed
  - Includes resource cleanup tests
- **Building**: `make build` command
  - Spec: "Build the application binary"
  - Depends on linting and testing passing
- **Cleaning**: `make clean` command
  - Spec: "Remove build artifacts"

### Development Workflow
**Specification Requirements**:
- **Code Quality**: All code must pass linting before commit
- **Testing**: All tests must pass before commit
- **Error Handling**: All errors must be properly handled
- **Resource Management**: All temporary resources must be cleaned up
- **Documentation**: All public functions must be documented
- **Backward Compatibility**: New features must not break existing functionality

## Testing Architecture

### Unit Tests
**Implementation**: `*_test.go` files
**Specification Requirements**:
- `TestGenerateArchiveName`: Tests archive naming
  - Spec: "Validates naming format for Git/non-Git repositories"
- `TestShouldExcludeFile`: Tests file exclusion
  - Spec: "Validates pattern matching behavior"
- `TestDefaultConfig` and `TestLoadConfig`: Tests configuration
  - Spec: "Validates configuration loading and defaults"
  - Tests configuration discovery with `BKPDIR_CONFIG` environment variable
  - Tests multiple configuration file precedence
  - Tests home directory expansion in paths
  - Tests status code configuration loading and defaults
  - Validates all status code fields are properly loaded from YAML
  - Tests status code configuration with custom values
- `TestGetConfigSearchPath`: Tests configuration path discovery
  - Validates environment variable parsing
  - Tests hard-coded default path when environment variable not set
  - Tests colon-separated path list parsing
  - Tests home directory expansion
- `TestDisplayConfig`: Tests configuration value display
  - Validates configuration value computation and source tracking
  - Tests environment variable processing for configuration paths
  - Tests hard-coded default path handling
  - Tests default value handling and source attribution
  - Tests output format with name, value, and source
  - Tests display of status code configuration values
- `TestVerifyArchive`: Tests archive verification
  - Spec: "Validates archive structure verification"
  - Tests verification status storage and loading
  - Tests checksum verification using base filename as key
- `TestVerifyChecksums`: Tests checksum verification
  - Spec: "Validates file content verification"
- `TestGenerateChecksums`: Tests checksum generation
  - Spec: "Validates checksum generation and storage"
  - Ensures directories are skipped
- `TestVerifyCorruptedArchive`: Tests corruption detection
  - Spec: "Validates error handling for corrupted archives"
- `TestResourceManager`: Tests resource management functionality
  - Validates thread-safe resource tracking
  - Tests automatic cleanup mechanisms
  - Test cases:
    - Basic resource registration and cleanup
    - Thread-safe concurrent access
    - Cleanup of both files and directories
    - Error-resilient cleanup (continues on individual failures)
    - Cleanup warnings logged to stderr
- `TestArchiveError`: Tests structured error handling
  - Validates ArchiveError functionality
  - Test cases:
    - Error creation with message and status code
    - Error interface implementation
    - Status code extraction
    - Error message formatting
- `TestIsDiskFullError`: Tests enhanced disk space detection
  - Validates disk full error detection
  - Test cases:
    - Various disk full error messages
    - Case-insensitive matching
    - Multiple disk space indicators
    - Non-disk-full errors (should return false)
    - Nil error handling
- `TestCompareDirectories`: Tests directory comparison
  - Validates directory tree comparison
  - Test cases:
    - Identical directories
    - Different directories
    - Directories with different file sizes
    - Empty directories
    - Large directories
    - Directories with special characters
- `TestConfigurationDiscovery`: Tests configuration file discovery
  - Tests multiple configuration files with different precedence
  - Tests environment variable override behavior
  - Tests hard-coded default path behavior
  - Tests missing configuration files handling
  - Tests invalid configuration file handling
  - Tests configuration merging with defaults
  - Tests status code configuration precedence and merging
- `TestStatusCodeConfiguration`: Tests status code configuration
  - Validates status code loading from YAML
  - Tests status code defaults
  - Tests status code precedence with multiple configuration files
  - Tests invalid status code values handling

### Output Formatting Tests
**Implementation**: `*_test.go` files
**Specification Requirements**:
- `TestOutputFormatter`: Tests output formatter functionality
  - Validates printf-style formatting with various inputs
  - Tests ANSI color code support
  - Test cases:
    - Basic format string application
    - Format strings with ANSI color codes
    - Format strings with special characters
    - Invalid format strings (should fall back to safe defaults)
    - Empty format strings
    - Format strings with multiple parameters
- `TestFormatCreatedArchive`: Tests archive creation message formatting
  - Validates format string application for successful archive messages
- `TestFormatIdenticalArchive`: Tests identical directory message formatting
  - Validates format string application for identical directory messages
- `TestFormatListArchive`: Tests archive listing entry formatting
  - Validates format string application for archive list entries
- `TestFormatConfigValue`: Tests configuration value display formatting
  - Validates format string application for configuration display
- `TestFormatDryRunArchive`: Tests dry-run message formatting
  - Validates format string application for dry-run messages
- `TestFormatError`: Tests error message formatting
  - Validates format string application for error messages
- `TestTemplateFormatter`: Tests template formatter functionality
  - Validates template-based formatting with named placeholders
  - Tests regex pattern extraction and template application
  - Test cases:
    - Basic template string application
    - Template strings with named placeholders
    - Regex pattern extraction with named groups
    - Template strings with Go text/template syntax
    - Invalid template strings (should fall back to safe defaults)
    - Template strings with conditional formatting

### Integration Tests
**Implementation**: `*_test.go` files
**Specification Requirements**:
- `TestFullCmdWithNote`: Tests full archive command
  - Spec: "Validates full archive creation with notes"
  - Tests status code exit behavior for different conditions
  - Tests panic recovery and error handling
  - Test cases:
    - Successful archive creation (should exit with `cfg.StatusCreatedArchive`)
    - Directory identical to existing archive (should exit with `cfg.StatusDirectoryIsIdenticalToExistingArchive`)
    - Directory not found (should exit with `cfg.StatusDirectoryNotFound`)
    - Invalid directory type (should exit with `cfg.StatusInvalidDirectoryType`)
    - Permission denied (should exit with `cfg.StatusPermissionDenied`)
    - Disk full scenarios (should exit with `cfg.StatusDiskFull`)
    - Archive directory creation failure (should exit with `cfg.StatusFailedToCreateArchiveDirectory`)
    - Configuration errors (should exit with `cfg.StatusConfigError`)
    - Panic recovery scenarios
- `TestIncCmdWithNote`: Tests incremental archive command
  - Spec: "Validates incremental archive creation with notes"
- `TestCmdArgsValidation`: Tests command-line arguments
  - Spec: "Validates command-line interface"
- `TestVerifyCmdWithChecksums`: Tests verify command
  - Spec: "Validates archive verification workflow"
- `TestVerifyOnCreate`: Tests automatic verification
  - Spec: "Validates automatic verification after creation"
- `TestCreateArchiveWithCleanup`: Tests archive creation with resource cleanup
  - Validates automatic resource cleanup functionality
  - Tests atomic operations with temporary files
  - Test cases:
    - Successful archive with cleanup verification
    - Archive failure with cleanup verification
    - No temporary files left after operations
    - Atomic file operations
- `TestCreateArchiveWithContext`: Tests context-aware archive creation
  - Validates archive creation with cancellation support
  - Test cases:
    - Successful archive with context
    - Archive cancelled via context
    - Archive with timeout
    - Context cancellation at various stages
- `TestCreateArchiveWithContextAndCleanup`: Tests context-aware archive with cleanup
  - Validates most robust archive creation functionality
  - Combines context support with resource cleanup
  - Test cases:
    - Successful archive with context and cleanup
    - Cancelled archive with proper cleanup
    - Timeout scenarios with cleanup verification
    - No resource leaks on cancellation

### Test Approach
**Specification Requirements**:
- Uses temporary directories for testing
- Simulates user environment with test files and directories
- Tests both dry-run and actual execution modes
- Verifies correct behavior for edge cases and error conditions
- Tests verification with both valid and corrupted archives
- Ensures that checksum keys match archive file entries (base filename)
- Ensures directories are not included in checksums and that corrupted archives are detected
- **Resource cleanup verification in all test scenarios**
- **Context cancellation and timeout handling testing**
- **Performance benchmarks for critical operations**
- **All code must pass linting before commit** 