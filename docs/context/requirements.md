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

# Run all tests
make test

# Clean build artifacts
make clean
```

### Build System Integration
**Implementation**: `Makefile`, CI/CD pipeline
**Specification Requirements**:
- **Quality Gates**: Build system enforces quality standards before compilation
  - Spec: "All code must pass linting and testing before successful build"
  - Implementation: `make build` depends on `make lint` and `make test` targets
  - Behavior: Build fails if linting or testing fails
  - Error Cases: Non-zero exit codes from linting or testing prevent build
- **Dependency Management**: Proper ordering of build steps
  - Spec: "Build system must enforce proper dependency order"
  - Implementation: Makefile targets with proper dependencies
  - Behavior: Linting runs before testing, testing runs before building
  - Error Cases: Dependency failures prevent subsequent steps
- **Artifact Management**: Clean build artifacts and temporary files
  - Spec: "Build system must provide clean target for artifact removal"
  - Implementation: `make clean` target removes all build artifacts
  - Behavior: Removes compiled binaries and temporary files
  - Error Cases: None (clean continues even if some files cannot be removed)
- **Continuous Integration**: Automated quality checks in CI/CD
  - Spec: "CI/CD pipeline must run all quality checks automatically"
  - Implementation: CI configuration runs make targets in proper order
  - Behavior: Automated testing and linting on code changes
  - Error Cases: CI fails if any quality check fails

**Example Usage**:
```bash
# Full build with all quality checks
make build

# Individual quality checks
make lint
make test

# Clean all artifacts
make clean

# CI/CD pipeline commands
make ci-lint
make ci-test
make ci-build
```

### Resource Management Requirements
**Implementation**: `archive.go`, `backup.go` - `ResourceManager`
**Specification Requirements**:
- **Resource Cleanup**: All temporary resources must be cleaned up
  - Spec: "No temporary files or directories should remain after operations"
  - Implementation: `ResourceManager` struct with automatic cleanup
  - Thread-safe: Uses mutex for concurrent access
  - Error-resilient: Continues cleanup even if individual operations fail
- **Atomic Operations**: Archive and backup operations must be atomic
  - Spec: "Archive and backup creation must be atomic to prevent corruption"
  - Implementation: Temporary files with atomic rename operations
  - Cleanup: Temporary files registered for automatic cleanup
- **Panic Recovery**: Operations must recover from panics
  - Spec: "Unexpected panics must not leave resources uncleaned"
  - Implementation: Defer functions with panic recovery
  - Logging: Panic information logged to stderr
- **Thread Safety**: Resource manager must be thread-safe
  - Spec: "Support concurrent operations with proper resource tracking"
  - Implementation: Mutex-protected resource lists
  - Methods: `AddTempFile()`, `AddTempDir()`, `RemoveResource()`, `CleanupWithPanicRecovery()`

**Example Usage**:
```go
// Create resource manager
rm := NewResourceManager()
defer rm.CleanupWithPanicRecovery()

// Register temporary resources
rm.AddTempFile("/tmp/archive.tmp")
rm.AddTempDir("/tmp/archive_work")

// Remove resource from tracking (successful operation)
rm.RemoveResource(&TempFile{Path: "/tmp/archive.tmp"})
```

### Enhanced Error Handling Requirements
**Implementation**: `archive.go`, `backup.go`, `errors.go` - Enhanced error handling
**Specification Requirements**:
- **Structured Errors**: Use `ArchiveError` and `BackupError` for consistent error handling
  - Spec: "All operations must return structured errors with status codes"
  - Fields:
    - `Message`: Human-readable error description
    - `StatusCode`: Numeric exit code for application
    - `Operation`: Context about what operation failed
    - `Path`: File or directory path involved in the error
    - `Err`: Underlying error for debugging
  - Implementation: Implements Go's `error` interface
  - Usage: Allows callers to extract both message and status code
- **Enhanced Disk Space Detection**: Comprehensive disk space error detection
  - Spec: "Detect various disk full conditions with case-insensitive matching"
  - Function: `isDiskFullError(err error) bool`
  - Detects: "no space left", "disk full", "not enough space", "insufficient disk space", "device full", "quota exceeded", "file too large"
  - Implementation: Case-insensitive string matching on error messages
- **Operation Context Support**: Enhanced error messages with operation context
  - Spec: "Error messages must include operation context for better debugging"
  - Implementation: Error types include operation field for context
  - Usage: Template formatting can use operation context for rich error display
  - Examples: "archive creation", "backup creation", "file listing", "directory scanning"

**Example Usage**:
```go
// Create structured archive error with operation context
err := NewArchiveError("file not found", cfg.StatusFileNotFound, "archive creation", "/path/to/source")

// Create structured backup error with operation context
err := NewBackupError("permission denied", cfg.StatusPermissionDenied, "backup creation", "/path/to/file")

// Check for ArchiveError type
if archiveErr, ok := err.(*ArchiveError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", archiveErr.Message)
    os.Exit(archiveErr.StatusCode)
}

// Check for BackupError type
if backupErr, ok := err.(*BackupError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", backupErr.Message)
    os.Exit(backupErr.StatusCode)
}

// Check for disk space errors
if err := someFileOperation(); err != nil {
    if isDiskFullError(err) {
        return NewArchiveError("Disk full", cfg.StatusDiskFull, "archive creation", sourcePath)
    }
    return err
}

// Use with template formatting for rich error display
templateFormatter := NewTemplateFormatter(cfg)
if backupErr, ok := err.(*BackupError); ok {
    templateFormatter.PrintTemplateError(backupErr.Message, backupErr.Operation)
}
```

### Context Support Requirements
**Implementation**: `archive.go`, `backup.go` - Context-aware operations
**Specification Requirements**:
- **Context-Aware Operations**: All long-running operations must support context cancellation
  - Spec: "Support operation cancellation and timeouts"
  - Functions: `CreateArchiveWithContext()`, `CreateFileBackupWithContext()`, `CopyFileWithContext()`
  - Behavior: Check for cancellation at multiple operation points
  - Error Handling: Return appropriate error on cancellation
- **Timeout Support**: Operations must support configurable timeouts
  - Spec: "Prevent operations from hanging indefinitely"
  - Implementation: Context with timeout support
  - Usage: `context.WithTimeout()` for operation limits
- **Enhanced File Operations**: Context-aware file backup operations
  - Spec: "All file backup operations must support context cancellation and timeouts"
  - Functions: 
    - `CreateFileBackupWithContext(ctx context.Context, cfg *Config, filePath, note string, dryRun bool) error`
    - `CopyFileWithContext(ctx context.Context, src, dst string) error`
    - `CompareFilesWithContext(ctx context.Context, file1, file2 string) (bool, error)`
    - `ListFileBackupsWithContext(ctx context.Context, backupDir, sourceFile string) ([]*Backup, error)`
  - Behavior: Check context cancellation during file I/O operations
  - Error Handling: Return context.Canceled or context.DeadlineExceeded appropriately

**Example Usage**:
```go
// Create context with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use context-aware archive operations
err := CreateArchiveWithContext(ctx, cfg, note, dryRun, verify)
if err == context.Canceled {
    log.Println("Archive operation was cancelled")
} else if err == context.DeadlineExceeded {
    log.Println("Archive operation timed out")
}

// File backup with context
err = CreateFileBackupWithContext(ctx, cfg, filePath, note, dryRun)
if err == context.Canceled {
    log.Println("File backup was cancelled")
} else if err == context.DeadlineExceeded {
    log.Println("File backup timed out")
}

// File operations with context
identical, err := CompareFilesWithContext(ctx, "file1.txt", "file2.txt")
if err == context.Canceled {
    log.Println("File comparison was cancelled")
}

// List backups with context
backups, err := ListFileBackupsWithContext(ctx, cfg.BackupDirPath, "document.txt")
if err == context.Canceled {
    log.Println("Backup listing was cancelled")
}
```

### File Backup Requirements
**Implementation**: `backup.go`
**Specification Requirements**:
- **Individual File Backup**: Create backups of single files with comparison
  - Spec: "Support backing up individual files with timestamp and note"
  - Functions: `CreateFileBackup()`, `CreateFileBackupWithContext()`
  - Features: Atomic operations, context support, identical file detection
  - Naming: `filename-timestamp[=note]` format
- **File Comparison**: Compare files to detect identical backups
  - Spec: "Detect when file is identical to most recent backup"
  - Function: `CheckForIdenticalFileBackup()`, `CompareFiles()`, `CompareFilesWithContext()`
  - Implementation: Byte-by-byte comparison with existing backups
  - Status: Return appropriate status code for identical files
  - Context support: Cancellable file comparison operations
- **Backup Listing**: List backups for specific files
  - Spec: "List all backups for a given file path sorted by creation time"
  - Function: `ListFileBackups()`, `ListFileBackupsWithContext()`
  - Sorting: Most recent first
  - Display: Configurable format strings and template formatting
  - Context support: Cancellable backup listing operations
- **Enhanced File Operations**: Complete file system operations for backup management
  - Spec: "Provide comprehensive file operations with error handling and context support"
  - Functions:
    - `CopyFile(src, dst string) error`: Create exact copy with permissions preserved
    - `CopyFileWithContext(ctx context.Context, src, dst string) error`: Context-aware file copying
    - `GenerateBackupName(sourceFile, note string) string`: Generate backup filename
    - `GetMostRecentBackup(backupDir, sourceFile string) (*Backup, error)`: Find latest backup
    - `ValidateFileForBackup(filePath string) error`: Validate file can be backed up
  - Features: Atomic operations, permission preservation, modification time handling
  - Error handling: Structured errors with operation context

**Example Usage**:
```go
// Create file backup
err := CreateFileBackup(cfg, "document.txt", "before-changes", false)

// Create file backup with context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
err = CreateFileBackupWithContext(ctx, cfg, "script.sh", "working-version", false)

// List file backups
backups, err := ListFileBackups(cfg.BackupDirPath, "document.txt")
for _, backup := range backups {
    fmt.Printf("%s (created: %s)\n", backup.Name, backup.CreationTime)
}

// List file backups with context
backups, err = ListFileBackupsWithContext(ctx, cfg.BackupDirPath, "document.txt")
if err == context.Canceled {
    log.Println("Backup listing was cancelled")
}

// Check for identical backup
isIdentical, err := CheckForIdenticalFileBackup("document.txt", cfg)
if isIdentical {
    fmt.Println("File is identical to existing backup")
}

// Compare files with context
identical, err := CompareFilesWithContext(ctx, "file1.txt", "file2.txt")
if err == context.Canceled {
    log.Println("File comparison was cancelled")
}

// Copy file with context
err = CopyFileWithContext(ctx, "source.txt", "destination.txt")
if err == context.Canceled {
    log.Println("File copy was cancelled")
}

// Generate backup name
backupName := GenerateBackupName("document.txt", "important-version")
// Returns: "document.txt-2024-03-21-15-30=important-version"

// Get most recent backup
mostRecent, err := GetMostRecentBackup(cfg.BackupDirPath, "document.txt")
if mostRecent != nil {
    fmt.Printf("Most recent backup: %s\n", mostRecent.Path)
}

// Validate file for backup
err = ValidateFileForBackup("document.txt")
if err != nil {
    log.Printf("File validation failed: %v", err)
}
```

### Output Formatting Requirements
**Implementation**: `formatter.go`
**Specification Requirements**:
- **Printf-Style Formatting**: Centralized printf-style formatting for all application output
  - Spec: "All standard output must use configurable printf-style format specifications"
  - Methods for directory operations: `FormatCreatedArchive()`, `FormatIdenticalArchive()`, `FormatListArchive()`, `FormatDryRunArchive()`
  - Methods for file operations: `FormatCreatedBackup()`, `FormatIdenticalBackup()`, `FormatListBackup()`, `FormatDryRunBackup()`
  - Print methods: `PrintCreatedArchive()`, `PrintIdenticalArchive()`, etc.
  - Error handling: `FormatError()`, `PrintError()`

- **Template-Based Formatting**: Advanced template formatting with regex extraction
  - Spec: "Support both Go text/template syntax and placeholder syntax with data extraction"
  - Template methods: `FormatCreatedArchiveTemplate()`, `FormatIdenticalArchiveTemplate()`, etc.
  - Regex extraction: `ExtractArchiveFilenameData()`, `ExtractBackupFilenameData()`
  - Enhanced formatting: `FormatArchiveWithExtraction()`, `FormatListArchiveWithExtraction()`

**Example Usage**:
```go
// Create output formatter
formatter := NewOutputFormatter(cfg)

// Printf-style formatting for archives
createdMsg := formatter.FormatCreatedArchive("/path/to/archive")
formatter.PrintCreatedArchive("/path/to/archive")

// Printf-style formatting for file backups
backupMsg := formatter.FormatCreatedBackup("/path/to/backup")
formatter.PrintCreatedBackup("/path/to/backup")

// Template-based formatting with extraction
data := formatter.ExtractArchiveFilenameData("project-2024-03-21-15-30=main=abc123=note.zip")
templateMsg := formatter.FormatCreatedArchiveTemplate(data)

// Template-based formatting for file backups
backupData := formatter.ExtractBackupFilenameData("document.txt-2024-03-21-15-30=important")
backupTemplateMsg := formatter.FormatCreatedBackupTemplate(backupData)
```

### Template Formatting Requirements
**Implementation**: `template_formatter.go` - `TemplateFormatter`
**Specification Requirements**:
- **Advanced Template Processing**: Support for both Go text/template and placeholder syntax
  - Spec: "Provide rich template formatting with named placeholders and regex extraction"
  - Methods:
    - `NewTemplateFormatter(cfg *Config) *TemplateFormatter`: Create formatter with configuration
    - `FormatWithTemplate(input, pattern, tmplStr string) (string, error)`: Apply Go text/template to named regex groups
    - `FormatWithPlaceholders(format string, data map[string]string) string`: Replace %{name} placeholders
  - Template syntax: Supports both `{{.name}}` (Go text/template) and `%{name}` (placeholder) syntax
  - Regex integration: Extract named groups for template data
  - Error handling: Graceful degradation on template errors

- **Template Methods for Directory Operations**: Rich template formatting for archive operations
  - Spec: "Provide template-based formatting for all directory archive operations"
  - Methods:
    - `TemplateCreatedArchive(path string) string`: Format archive creation using template with data extraction
    - `TemplateIdenticalArchive(path string) string`: Format identical directory message using template
    - `TemplateListArchive(path, creationTime string) string`: Format archive listing using template
    - `TemplateDryRunArchive(path string) string`: Format dry-run message using template
  - Data extraction: Uses `cfg.PatternArchiveFilename` for extracting archive information
  - Fallback: Falls back to placeholder formatting on template errors

- **Template Methods for File Operations**: Rich template formatting for file backup operations
  - Spec: "Provide template-based formatting for all file backup operations"
  - Methods:
    - `TemplateCreatedBackup(path string) string`: Format backup creation using template with data extraction
    - `TemplateIdenticalBackup(path string) string`: Format identical file message using template
    - `TemplateListBackup(path, creationTime string) string`: Format backup listing using template
    - `TemplateDryRunBackup(path string) string`: Format dry-run message using template
  - Data extraction: Uses `cfg.PatternBackupFilename` for extracting backup information
  - Fallback: Falls back to placeholder formatting on template errors

- **Template Methods for Common Operations**: Template formatting for configuration and errors
  - Spec: "Provide template-based formatting for configuration display and error messages"
  - Methods:
    - `TemplateConfigValue(name, value, source string) string`: Format configuration display using template
    - `TemplateError(message, operation string) string`: Format error message using template with operation context
  - Enhanced features: Supports conditional formatting based on source type and operation context
  - Error context: Operation context enhances error information (e.g., "archive creation", "backup creation", "file listing")

- **Template Print Methods**: Direct printing with template formatting
  - Spec: "Provide direct printing methods for template-formatted output"
  - Methods for directory operations:
    - `PrintTemplateCreatedArchive(path string)`: Print template-formatted archive creation to stdout
    - `PrintTemplateIdenticalArchive(path string)`: Print template-formatted identical directory to stdout
    - `PrintTemplateListArchive(path, creationTime string)`: Print template-formatted archive listing to stdout
    - `PrintTemplateDryRunArchive(path string)`: Print template-formatted dry-run to stdout
  - Methods for file operations:
    - `PrintTemplateCreatedBackup(path string)`: Print template-formatted backup creation to stdout
    - `PrintTemplateIdenticalBackup(path string)`: Print template-formatted identical file to stdout
    - `PrintTemplateListBackup(path, creationTime string)`: Print template-formatted backup listing to stdout
    - `PrintTemplateDryRunBackup(path string)`: Print template-formatted dry-run to stdout
  - Common methods:
    - `PrintTemplateConfigValue(name, value, source string)`: Print template-formatted configuration to stdout
    - `PrintTemplateError(message, operation string)`: Print template-formatted error to stderr

- **Regex-Based Data Extraction**: Extract structured data from filenames and paths
  - Spec: "Use named regex patterns to extract rich data for template formatting"
  - Patterns: Archive filenames, backup filenames, timestamps, configuration lines
  - Named groups: Extract prefix, year, month, day, hour, minute, branch, hash, note, filename, etc.
  - Integration: Seamless integration with template formatting
  - Error handling: Invalid regex patterns handled gracefully

**Example Usage**:
```go
// Create template formatter
templateFormatter := NewTemplateFormatter(cfg)

// Template formatting with regex extraction for archives
archivePath := "project-2024-03-21-15-30=main=abc123d=important.zip"
formatted, err := templateFormatter.FormatWithTemplate(
    archivePath,
    cfg.PatternArchiveFilename,
    "Archive of {{.prefix}} from {{.year}}-{{.month}}-{{.day}} ({{.branch}}/{{.hash}}) with note: {{.note}}",
)

// Template formatting with regex extraction for backups
backupPath := "file.txt-2024-03-21-15-30=important"
formatted, err := templateFormatter.FormatWithTemplate(
    backupPath,
    cfg.PatternBackupFilename,
    "Backup of {{.filename}} from {{.year}}-{{.month}}-{{.day}} with note: {{.note}}",
)

// Placeholder formatting
data := map[string]string{
    "filename": "document.txt",
    "year": "2024",
    "month": "03",
    "day": "21",
}
result := templateFormatter.FormatWithPlaceholders(
    "File %{filename} backed up on %{year}-%{month}-%{day}",
    data,
)

// Direct template formatting methods
createdMsg := templateFormatter.TemplateCreatedArchive("/path/to/archive-2024-03-21-15-30=main=abc123d.zip")
backupMsg := templateFormatter.TemplateCreatedBackup("/path/to/backup-2024-03-21-15-30=note")
errorMsg := templateFormatter.TemplateError("File not found", "backup creation")

// Print methods for direct output
templateFormatter.PrintTemplateCreatedArchive("/path/to/archive")
templateFormatter.PrintTemplateCreatedBackup("/path/to/backup")
templateFormatter.PrintTemplateError("Operation failed", "archive creation")
```

## Data Objects

### Config
**Implementation**: `config.go`
**Specification Requirements**:
- Configuration stored in YAML files with configurable discovery
- Environment variable support: `BKPDIR_CONFIG` for search paths
- Default search path: `./.bkpdir.yml:~/.bkpdir.yml`
- Home directory expansion support
- Configuration precedence: Earlier files override later files

**Directory Archiving Fields**:
- `ArchiveDirPath`: Archive storage location (default: "../.bkpdir")
- `UseCurrentDirName`: Include directory name in path (default: true)
- `ExcludePatterns`: Glob patterns to exclude (default: [".git/", "vendor/"])
- `IncludeGitInfo`: Include Git branch/hash in names (default: false)
- `Verification`: Archive verification settings

**File Backup Fields**:
- `BackupDirPath`: File backup storage location (default: "../.bkpdir")
- `UseCurrentDirNameForFiles`: Include directory structure in backup path (default: true)

**Status Codes for Directory Operations**:
- `StatusCreatedArchive`: New archive created (default: 0)
- `StatusFailedToCreateArchiveDirectory`: Archive directory creation failed (default: 31)
- `StatusDirectoryIsIdenticalToExistingArchive`: Directory unchanged (default: 0)
- `StatusDirectoryNotFound`: Source directory missing (default: 20)
- `StatusInvalidDirectoryType`: Source not a directory (default: 21)
- `StatusPermissionDenied`: Access denied (default: 22)
- `StatusDiskFull`: Insufficient disk space (default: 30)
- `StatusConfigError`: Invalid configuration (default: 10)

**Status Codes for File Operations**:
- `StatusCreatedBackup`: New backup created (default: 0)
- `StatusFailedToCreateBackupDirectory`: Backup directory creation failed (default: 31)
- `StatusFileIsIdenticalToExistingBackup`: File unchanged (default: 0)
- `StatusFileNotFound`: Source file missing (default: 20)
- `StatusInvalidFileType`: Source not a regular file (default: 21)

**Printf-Style Format Strings for Directory Operations**:
- `FormatCreatedArchive`: Archive creation message (default: "Created archive: %s\n")
- `FormatIdenticalArchive`: Identical directory message (default: "Directory is identical to existing archive: %s\n")
- `FormatListArchive`: Archive listing format (default: "%s (created: %s)\n")
- `FormatDryRunArchive`: Dry-run message (default: "Would create archive: %s\n")

**Printf-Style Format Strings for File Operations**:
- `FormatCreatedBackup`: Backup creation message (default: "Created backup: %s\n")
- `FormatIdenticalBackup`: Identical file message (default: "File is identical to existing backup: %s\n")
- `FormatListBackup`: Backup listing format (default: "%s (created: %s)\n")
- `FormatDryRunBackup`: Dry-run message (default: "Would create backup: %s\n")

**Common Format Strings**:
- `FormatConfigValue`: Configuration display (default: "%s: %s (source: %s)\n")
- `FormatError`: Error messages (default: "Error: %s\n")

**Template-Based Format Strings for Directory Operations**:
- `TemplateCreatedArchive`: Archive creation template (default: "Created archive: %{path}\n")
- `TemplateIdenticalArchive`: Identical directory template (default: "Directory is identical to existing archive: %{path}\n")
- `TemplateListArchive`: Archive listing template (default: "%{path} (created: %{creation_time})\n")
- `TemplateDryRunArchive`: Dry-run template (default: "Would create archive: %{path}\n")

**Template-Based Format Strings for File Operations**:
- `TemplateCreatedBackup`: Backup creation template (default: "Created backup: %{path}\n")
- `TemplateIdenticalBackup`: Identical file template (default: "File is identical to existing backup: %{path}\n")
- `TemplateListBackup`: Backup listing template (default: "%{path} (created: %{creation_time})\n")
- `TemplateDryRunBackup`: Dry-run template (default: "Would create backup: %{path}\n")

**Common Template Strings**:
- `TemplateConfigValue`: Configuration display template (default: "%{name}: %{value} (source: %{source})\n")
- `TemplateError`: Error message template (default: "Error: %{message}\n")

**Regex Patterns for Data Extraction**:
- `PatternArchiveFilename`: Parse archive filenames with named groups
- `PatternBackupFilename`: Parse backup filenames with named groups
- `PatternConfigLine`: Parse configuration display lines
- `PatternTimestamp`: Parse timestamp strings

**Example Usage**:
```go
// Load configuration with environment variable support
cfg, err := LoadConfig(".")
if err != nil {
    log.Fatal(err)
}

// Access directory archiving configuration
archivePath := cfg.ArchiveDirPath
excludePatterns := cfg.ExcludePatterns
includeGit := cfg.IncludeGitInfo

// Access file backup configuration
backupPath := cfg.BackupDirPath
useCurrentDir := cfg.UseCurrentDirNameForFiles

// Access status codes
createdStatus := cfg.StatusCreatedArchive
identicalStatus := cfg.StatusDirectoryIsIdenticalToExistingArchive
createdBackupStatus := cfg.StatusCreatedBackup
identicalBackupStatus := cfg.StatusFileIsIdenticalToExistingBackup

// Access format strings
createdFormat := cfg.FormatCreatedArchive
createdTemplate := cfg.TemplateCreatedArchive
createdBackupFormat := cfg.FormatCreatedBackup
createdBackupTemplate := cfg.TemplateCreatedBackup

// Access regex patterns
archivePattern := cfg.PatternArchiveFilename
backupPattern := cfg.PatternBackupFilename
```

### ConfigValue
**Implementation**: `config.go`
**Specification Requirements**:
- Fields:
  - `Name`: Configuration parameter name
  - `Value`: Computed configuration value including defaults
  - `Source`: Source file path or "default" for default values

**Example Usage**:
```go
// Create configuration value entry
configValue := &ConfigValue{
    Name:   "archive_dir_path",
    Value:  "../.bkpdir",
    Source: "~/.bkpdir.yml",
}
```

### Backup
**Implementation**: `backup.go`
**Specification Requirements**:
- **File Backup Representation**: Represents a single file backup with metadata
  - Spec: "Provide structured representation of file backups with creation time and source information"
  - Fields:
    - `Name`: Backup filename (e.g., "file.txt-2024-03-21-15-30=note")
    - `Path`: Full path to backup file
    - `CreationTime`: Time when backup was created
    - `SourceFile`: Path to original source file
    - `Note`: Optional note extracted from backup filename
  - Usage: Used for backup listing, comparison, and metadata operations
  - Integration: Works with template formatting for rich display

**Example Usage**:
```go
// Create a new backup object
backup := &Backup{
    Name: "file.txt-2024-03-20-15-30=important_backup",
    Path: "/path/to/backup/file.txt-2024-03-20-15-30=important_backup",
    CreationTime: time.Now(),
    SourceFile: "/path/to/original/file.txt",
    Note: "important_backup",
}

// Use in backup listing
backups, err := ListFileBackups(cfg.BackupDirPath, "file.txt")
for _, backup := range backups {
    fmt.Printf("Backup: %s (created: %s)\n", backup.Name, backup.CreationTime.Format("2006-01-02 15:04:05"))
    if backup.Note != "" {
        fmt.Printf("  Note: %s\n", backup.Note)
    }
}

// Use with template formatting
templateFormatter := NewTemplateFormatter(cfg)
for _, backup := range backups {
    templateFormatter.PrintTemplateListBackup(backup.Path, backup.CreationTime.Format("2006-01-02 15:04:05"))
}
```

### BackupError
**Implementation**: `backup.go`
**Specification Requirements**:
- **Structured Error Handling**: Provides consistent error reporting for file backup operations
  - Spec: "All backup operations return structured errors with status codes and operation context"
  - Fields:
    - `Message`: Human-readable error description
    - `StatusCode`: Numeric exit code for application
    - `Operation`: Context about what operation failed (e.g., "backup creation", "file listing")
    - `Path`: File path involved in the error
    - `Err`: Underlying error for debugging
  - Implementation: Implements Go's `error` interface
  - Usage: Allows callers to extract both message and status code
  - Integration: Works with template formatting for rich error display

**Example Usage**:
```go
// Create structured backup error with operation context
err := NewBackupError("file not found", cfg.StatusFileNotFound, "backup creation", "/path/to/file")

// Check for BackupError type
if backupErr, ok := err.(*BackupError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", backupErr.Message)
    
    // Use with template formatting for rich error display
    templateFormatter := NewTemplateFormatter(cfg)
    templateFormatter.PrintTemplateError(backupErr.Message, backupErr.Operation)
    
    os.Exit(backupErr.StatusCode)
}

// Create error with underlying cause
underlyingErr := errors.New("permission denied")
err = NewBackupErrorWithCause("cannot create backup", cfg.StatusPermissionDenied, "backup creation", "/path/to/file", underlyingErr)
```

### ArchiveError
**Implementation**: `archive.go`
**Specification Requirements**:
- **Structured Error Handling**: Provides consistent error reporting for archive operations
  - Spec: "All archive operations return structured errors with status codes and operation context"
  - Fields:
    - `Message`: Human-readable error description
    - `StatusCode`: Numeric exit code for application
    - `Operation`: Context about what operation failed (e.g., "archive creation", "directory scanning")
    - `Path`: Directory path involved in the error
    - `Err`: Underlying error for debugging
  - Implementation: Implements Go's `error` interface
  - Usage: Allows callers to extract both message and status code
  - Integration: Works with template formatting for rich error display

**Example Usage**:
```go
// Create structured archive error with operation context
err := NewArchiveError("directory not found", cfg.StatusDirectoryNotFound, "archive creation", "/path/to/directory")

// Check for ArchiveError type
if archiveErr, ok := err.(*ArchiveError); ok {
    fmt.Fprintf(os.Stderr, "Error: %s\n", archiveErr.Message)
    
    // Use with template formatting for rich error display
    templateFormatter := NewTemplateFormatter(cfg)
    templateFormatter.PrintTemplateError(archiveErr.Message, archiveErr.Operation)
    
    os.Exit(archiveErr.StatusCode)
}

// Create error with underlying cause
underlyingErr := errors.New("permission denied")
err = NewArchiveErrorWithCause("cannot access directory", cfg.StatusPermissionDenied, "archive creation", "/path/to/directory", underlyingErr)
```

### VerificationConfig
**Implementation**: `config.go`
**Specification Requirements**:
- Fields:
  - `VerifyOnCreate`: Automatically verify archives after creation (default: false)
  - `ChecksumAlgorithm`: Algorithm used for checksums (default: "sha256")

### Archive
**Implementation**: `archive.go`
**Specification Requirements**:
- Fields:
  - `Name`: Archive filename with timestamp and optional Git info
  - `Path`: Full path to archive file
  - `CreationTime`: Archive creation timestamp
  - `IsIncremental`: Whether this is an incremental archive
  - `GitBranch`: Git branch name (if in Git repository)
  - `GitHash`: Git commit hash (if in Git repository)
  - `Note`: Optional user-provided note
  - `BaseArchive`: Base archive name for incremental archives
  - `VerificationStatus`: Result of last verification

### BackupInfo
**Implementation**: `backup.go`
**Specification Requirements**:
- Fields:
  - `Name`: Backup filename with timestamp and optional note
  - `Path`: Full path to backup file
  - `CreationTime`: Backup creation timestamp
  - `Size`: File size in bytes

### VerificationStatus
**Implementation**: `archive.go`
**Specification Requirements**:
- Fields:
  - `IsVerified`: Whether verification passed
  - `ChecksumVerified`: Whether checksum verification was performed
  - `Errors`: List of verification errors
  - `VerifiedAt`: When verification was performed

### ResourceManager
**Implementation**: `archive.go`, `backup.go`
**Specification Requirements**:
- **Thread-Safe Resource Tracking**: Mutex-protected resource management
  - Fields:
    - `tempFiles`: List of temporary files to clean up
    - `tempDirs`: List of temporary directories to clean up
    - `mutex`: Mutex for thread-safe access
  - Methods:
    - `NewResourceManager()`: Create new resource manager
    - `AddTempFile(path string)`: Register temporary file for cleanup
    - `AddTempDir(path string)`: Register temporary directory for cleanup
    - `RemoveResource(resource Resource)`: Remove resource from tracking
    - `Cleanup()`: Clean up all registered resources
    - `CleanupWithPanicRecovery()`: Clean up with panic recovery

**Example Usage**:
```go
// Create and use resource manager
rm := NewResourceManager()
defer rm.CleanupWithPanicRecovery()

// Register resources
rm.AddTempFile("/tmp/archive.tmp")
rm.AddTempDir("/tmp/work")

// Remove successful resource
rm.RemoveResource(&TempFile{Path: "/tmp/archive.tmp"})
```

### TemplateFormatter
**Implementation**: `formatter.go`
**Specification Requirements**:
- **Advanced Template Processing**: Template-based formatting with regex extraction
  - Fields:
    - `config`: Reference to configuration for template strings and patterns
  - Methods:
    - `NewTemplateFormatter(cfg *Config) *TemplateFormatter`: Create formatter
    - `FormatWithTemplate(input, pattern, tmplStr string) (string, error)`: Apply template to regex groups
    - `FormatWithPlaceholders(format string, data map[string]string) string`: Replace placeholders
    - Template methods for all operations: `TemplateCreatedArchive()`, `TemplateCreatedBackup()`, etc.

## Core Functions

### Configuration Management
**Implementation**: `config.go`
**Specification Requirements**:

- `DefaultConfig() *Config`
  - Spec: "Creates default configuration with all specified defaults"
  - Returns: New Config with all fields set to default values
  - Behavior: Initializes all directory archiving and file backup settings

- `LoadConfig(root string) (*Config, error)`
  - Spec: "Loads configuration with environment variable support and precedence"
  - Input: Root directory path
  - Behavior:
    - Reads `BKPDIR_CONFIG` environment variable for search paths
    - Falls back to default search path if not set
    - Processes files in order with precedence rules
    - Supports home directory expansion
    - Merges with default values
  - Error Cases: Invalid YAML, file system errors

- `getConfigSearchPaths() []string`
  - Spec: "Returns configuration file search paths from environment or defaults"
  - Behavior: Checks `BKPDIR_CONFIG`, returns colon-separated paths or defaults
  - Default: `./.bkpdir.yml:~/.bkpdir.yml`

- `expandPath(path string) string`
  - Spec: "Expands home directory (~) in configuration paths"
  - Input: Path that may contain ~ prefix
  - Output: Expanded absolute path

- `mergeConfigs(dst, src *Config)`
  - Spec: "Merges non-default values from source into destination config"
  - Behavior: Only overwrites fields that differ from defaults

- `GetConfigValues(cfg *Config) []ConfigValue`
  - Spec: "Returns all configuration values with their sources for display"
  - Output: List of configuration values showing name, value, and source

**Example Usage**:
```go
// Load configuration with environment support
cfg, err := LoadConfig(".")
if err != nil {
    log.Fatal(err)
}

// Get configuration values for display
values := GetConfigValues(cfg)
for _, cv := range values {
    fmt.Printf("%s: %s (source: %s)\n", cv.Name, cv.Value, cv.Source)
}
```

### File System Operations
**Implementation**: `archive.go`, `backup.go`
**Specification Requirements**:

- `CopyFile(src, dst string) error`
  - Spec: "Creates exact copy with permissions preserved"
  - Behavior: Copies file content and preserves permissions
  - Error Cases: Source not found, permission denied, disk full

- `CopyFileWithContext(ctx context.Context, src, dst string) error`
  - Spec: "Context-aware file copying with cancellation support"
  - Behavior: Same as CopyFile with cancellation checks
  - Error Cases: All CopyFile errors plus context cancellation

- `compareFiles(file1, file2 string) (bool, error)`
  - Spec: "Performs byte-by-byte comparison of files"
  - Behavior: Compares file sizes first, then content
  - Returns: True if files are identical

- `SafeMkdirAll(path string, perm os.FileMode, cfg *Config) error`
  - Spec: "Creates directory with enhanced error detection"
  - Behavior: Creates directory structure with proper error handling
  - Error Cases: Permission denied, disk full (detected via isDiskFullError)

**Example Usage**:
```go
// Copy file with context
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err := CopyFileWithContext(ctx, "source.txt", "backup.txt")

// Compare files
identical, err := compareFiles("file1.txt", "file2.txt")
if err != nil {
    log.Fatal(err)
}
```

### Enhanced Error Detection
**Implementation**: `archive.go`, `backup.go`
**Specification Requirements**:

- `isDiskFullError(err error) bool`
  - Spec: "Enhanced disk space error detection with comprehensive pattern matching"
  - Input: Error to check
  - Behavior: Case-insensitive matching for multiple disk space indicators
  - Patterns: "no space left", "disk full", "not enough space", "insufficient disk space", "device full", "quota exceeded", "file too large"
  - Returns: True if error indicates disk space issues

**Example Usage**:
```go
// Enhanced error detection
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

- `CreateFullArchive(cfg *Config, note string, dryRun bool, verify bool) error`
  - Spec: "Creates complete directory archive with comparison"
  - Behavior: Compares with existing archives, creates new archive if different
  - Features: Git integration, exclusion patterns, verification

- `CreateFullArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error`
  - Spec: "Context-aware archive creation with cancellation support"
  - Behavior: Same as CreateFullArchive with context cancellation

- `CreateIncrementalArchive(cfg *Config, note string, dryRun bool, verify bool) error`
  - Spec: "Creates incremental archive with changed files only"
  - Behavior: Requires existing full archive, includes only modified files

- `ListArchives(archiveDir string) ([]Archive, error)`
  - Spec: "Lists all archives with metadata extraction"
  - Behavior: Scans directory, extracts Git info and notes from filenames

- `CheckForIdenticalArchive(sourceDir, archiveDir string, excludePatterns []string) (bool, string, error)`
  - Spec: "Compares directory with most recent archive"
  - Behavior: Directory tree comparison with exclusion patterns

**Example Usage**:
```go
// Create archive with context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
err := CreateFullArchiveWithContext(ctx, cfg, "backup note", false, true)

// List archives
archives, err := ListArchives("/path/to/archives")
if err != nil {
    log.Fatal(err)
}
```

### File Backup Management
**Implementation**: `backup.go`
**Specification Requirements**:

- `CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error`
  - Spec: "Creates backup of single file with comparison"
  - Behavior: Compares with existing backups, creates new backup if different

- `CreateFileBackupEnhanced(ctx context.Context, cfg *Config, formatter *OutputFormatter, filePath string, note string, dryRun bool) error`
  - Spec: "Enhanced backup creation with context and formatting"
  - Behavior: Context support, enhanced formatting, resource cleanup

- `CreateFileBackupWithContext(ctx context.Context, cfg *Config, filePath string, note string, dryRun bool) error`
  - Spec: "Context-aware file backup creation"
  - Behavior: Same as CreateFileBackup with cancellation support

- `CreateFileBackupWithCleanup(cfg *Config, filePath string, note string, dryRun bool) error`
  - Spec: "File backup with automatic resource cleanup"
  - Behavior: Uses ResourceManager for atomic operations and cleanup

- `CreateFileBackupWithContextAndCleanup(ctx context.Context, cfg *Config, filePath string, note string, dryRun bool) error`
  - Spec: "Most robust backup creation with context and cleanup"
  - Behavior: Combines context support with resource management

- `CheckForIdenticalFileBackup(filePath, backupDir, baseFilename string) (bool, string, error)`
  - Spec: "Compares file with most recent backup using byte comparison"
  - Behavior: Size check first, then byte-by-byte comparison

- `ListFileBackups(backupDir, baseFilename string) ([]BackupInfo, error)`
  - Spec: "Lists all backups for specific file"
  - Behavior: Scans directory, extracts notes, sorts by creation time

- `ListFileBackupsEnhanced(cfg *Config, formatter *OutputFormatter, filePath string) error`
  - Spec: "Enhanced backup listing with formatting"
  - Behavior: Uses enhanced formatting with template support

**Example Usage**:
```go
// Create backup with full features
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
err := CreateFileBackupWithContextAndCleanup(ctx, cfg, "/path/to/file.txt", "important", false)

// List backups
backups, err := ListFileBackups("/backup/dir", "file.txt")
if err != nil {
    log.Fatal(err)
}
```

### Utility Functions
**Implementation**: Various files
**Specification Requirements**:

- `GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string) string`
  - Spec: "Generates archive filename with Git info and note"
  - Format: `prefix-YYYY-MM-DD-hh-mm[=branch=hash][=note].zip`

- `GenerateBackupName(sourcePath, timestamp, note string) string`
  - Spec: "Generates backup filename with timestamp and note"
  - Format: `filename-YYYY-MM-DD-hh-mm[=note]`

- `ValidateDirectoryPath(path string, cfg *Config) error`
  - Spec: "Validates directory exists and is accessible"
  - Error Cases: Not found, not directory, permission denied

- `ValidateFilePath(path string, cfg *Config) error`
  - Spec: "Validates file exists and is regular file"
  - Error Cases: Not found, not regular file, permission denied

**Example Usage**:
```go
// Generate names
archiveName := GenerateArchiveName("project", "2024-03-21-15-30", "main", "abc123", "backup")
backupName := GenerateBackupName("/path/to/file.txt", "2024-03-21-15-30", "important")

// Validate paths
err := ValidateDirectoryPath("/path/to/dir", cfg)
err = ValidateFilePath("/path/to/file.txt", cfg)
```

## Main Application Structure

### CLI Interface
**Implementation**: `main.go`
**Specification Requirements**:
- Uses `cobra` for command-line interface
- Supports both directory archiving and file backup operations
- Global flags:
  - `--dry-run`: Show what would be done without creating archives/backups
  - `--list`: List backups for specified file
  - `--config`: Display computed configuration values and exit
- Commands:
  - `full [NOTE]`: Create full directory archive
  - `inc [NOTE]`: Create incremental directory archive
  - `backup [FILE_PATH] [NOTE]`: Create file backup
  - `list`: List directory archives
  - `verify [ARCHIVE_NAME]`: Verify archive integrity
  - `version`: Show version information

### Enhanced Workflow Implementation
**Implementation**: `archive.go`, `backup.go`
**Specification Requirements**:

- **Directory Archive Workflow**: Enhanced with context and cleanup
  - Steps: Load config, validate directory, check for identical archive, create with resource management
  - Features: Git integration, exclusion patterns, verification

- **File Backup Workflow**: Enhanced with context and cleanup
  - Steps: Load config, validate file, check for identical backup, create with atomic operations
  - Features: Directory structure preservation, byte comparison, resource management

- **Configuration Display Workflow**: Enhanced with environment variable support
  - Steps: Process BKPDIR_CONFIG, merge configurations, display with sources

### Build and Development Requirements

### Build System
**Implementation**: `Makefile`
**Specification Requirements**:
- **Enhanced Build System**: Comprehensive build pipeline with quality gates
  - `make lint`: Run revive linter with `.revive.toml` configuration
  - `make test`: Run all tests including resource cleanup verification
  - `make build`: Build application (depends on lint and test)
  - `make clean`: Remove build artifacts
  - Quality gates: Linting and testing must pass before build

### Development Workflow
**Specification Requirements**:
- **Code Quality**: All code must pass linting and testing
- **Resource Management**: All temporary resources must be cleaned up
- **Context Support**: Long-running operations must support cancellation
- **Error Handling**: Enhanced error detection and structured error types
- **Template Support**: Rich formatting with regex extraction
- **Thread Safety**: Concurrent operations must be properly synchronized

**Example Build Usage**:
```bash
# Full build pipeline
make clean
make lint
make test
make build

# Development workflow
make lint && make test && make build
```

## Git Integration Requirements
- **Repository Detection**: Automatic Git repository detection using `git rev-parse --is-inside-work-tree`
- **Info Extraction**: Git branch and commit hash extraction using `git rev-parse` commands
- **Archive Naming**: Git information in archive names when enabled via configuration
- **Error Handling**: Graceful handling of non-Git directories and Git command failures
- **Configuration**: Git integration toggle via `include_git_info` setting
These Git integration requirements are mandatory and must be preserved

## Archive Verification Requirements

### Testing Infrastructure Requirements

#### Archive Corruption Testing Framework (TEST-INFRA-001-A) ✅ COMPLETED
**Implementation**: `internal/testutil/corruption.go`, `internal/testutil/corruption_test.go`
**Specification Requirements**:
- **Controlled ZIP Corruption Utilities**: Systematic corruption for verification testing
  - Spec: "Provide systematic corruption capabilities for testing archive verification logic"
  - Implementation: `ArchiveCorruptor` class with 8 distinct corruption types
  - Safety: Backup/restore functionality ensures safe testing without data loss
  - Configuration: `CorruptionConfig` struct allows precise control over corruption parameters
- **Corruption Type Enumeration**: Complete coverage of ZIP corruption scenarios
  - Spec: "Support all major ZIP corruption types for comprehensive verification testing"
  - Types: CRC errors, header corruption, truncation, invalid central directory, local header corruption, data corruption, signature corruption, comment corruption
  - Classification: Categorizes corruption as recoverable vs fatal for appropriate test expectations
  - Detection: Automatic corruption type identification through `CorruptionDetector`
- **Deterministic Corruption Patterns**: Reproducible corruption for consistent test results
  - Spec: "Corruption must be reproducible across test runs for reliable CI/CD testing"
  - Implementation: Seed-based corruption generation with offset variation
  - Reproducibility: Identical archives with identical seeds produce identical corruption
  - Verification: Test suite validates reproducibility across multiple test runs
- **Archive Repair Detection**: Test recovery behavior from various corruption types
  - Spec: "Enable testing of archive repair and recovery mechanisms"
  - Implementation: `CorruptionDetector` for automatic corruption classification
  - Recovery Testing: Framework tracks original bytes for potential recovery scenarios
  - Integration: Designed for use with `verify.go` and `comparison.go` testing

**Example Usage**:
```go
// Create systematic corruption for verification testing
config := CorruptionConfig{
    Type:           CorruptionCRC,
    Seed:           12345,
    CorruptionSize: 4,
    Severity:       0.5,
}

corruptor, err := NewArchiveCorruptor(archivePath, config)
if err != nil {
    return fmt.Errorf("failed to create corruptor: %w", err)
}
defer corruptor.Cleanup()

// Apply controlled corruption
result, err := corruptor.ApplyCorruption()
if err != nil {
    return fmt.Errorf("corruption failed: %w", err)
}

// Test verification behavior with corrupted archive
detector := NewCorruptionDetector(archivePath)
detected, err := detector.DetectCorruption()
if err != nil {
    return fmt.Errorf("detection failed: %w", err)
}

// Validate detection matches applied corruption
expectedTypes := []CorruptionType{CorruptionCRC}
if !containsAllTypes(detected, expectedTypes) {
    return fmt.Errorf("detection mismatch: got %v, want %v", detected, expectedTypes)
}
```

**Implementation Areas**:
- Testing infrastructure for `verify.go` archive verification logic
- Complex scenario testing for `comparison.go` archive comparison
- Error path testing for archive operations that are difficult to reproduce
- Performance baseline establishment for corruption detection algorithms
- Cross-platform testing utilities for ZIP archive manipulation

**Design Decisions**:
- **Comprehensive Corruption Types**: Implemented 8 different corruption types covering all major ZIP structure components
  - Rationale: Enables thorough testing of verification logic against real-world corruption scenarios
  - Implementation: Each corruption type targets specific ZIP format structures (headers, data, metadata)
  - Testing: Systematic validation ensures each corruption type behaves as expected
- **Deterministic Reproduction**: Used seed-based corruption generation for consistent test results
  - Rationale: CI/CD testing requires reproducible results for reliable quality gates
  - Implementation: Seed + offset combination generates different but deterministic corruption per location
  - Validation: Test suite verifies identical corruption across multiple runs with same seed
- **Safe Corruption Testing**: Implemented backup/restore functionality for safe testing
  - Rationale: Prevents accidental data loss during corruption testing
  - Implementation: Automatic backup creation before corruption, cleanup after testing
  - Safety: Panic recovery ensures cleanup occurs even during test failures
- **Go ZIP Reader Compatibility**: Designed corruption testing around Go's ZIP reader behavior
  - Rationale: Go's `zip.OpenReader` is surprisingly resilient to corruption
  - Implementation: Tests verify corruption effects rather than complete reading failure
  - Adaptation: Adjusted test expectations to match Go's actual ZIP handling behavior

**Performance Characteristics**:
- **CRC Corruption**: ~763μs average for typical archives
- **Corruption Detection**: ~49μs average for classification
- **Memory Usage**: Minimal additional memory overhead during testing
- **Cleanup Performance**: Automatic resource cleanup with <1ms overhead

**Cross-Platform Compatibility**:
- **Unix Systems**: Full functionality with proper file permission handling
- **Windows Systems**: Compatible file operations with appropriate path handling
- **Temporary Files**: Cross-platform temporary directory usage with proper cleanup
- **File Permissions**: Respects platform-specific file permission models

## 🔻 CI/CD Requirements for AI Development

### CICD-001: AI-First Development Optimization Requirements
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// CICD-001: AI-first CI/CD optimization`

#### Core Requirements
- **R-CICD-001-1**: All CI/CD tasks MUST be configured with low priority scheduling
- **R-CICD-001-2**: Pipeline configurations MUST assume zero human developer intervention
- **R-CICD-001-3**: All documentation and code MUST be maintained for AI assistant comprehension
- **R-CICD-001-4**: CI/CD outputs MUST use standardized icon system from DOC-007/DOC-008
- **R-CICD-001-5**: Pipeline MUST integrate with implementation token validation system

#### AI Assistant Optimization Requirements
- **R-CICD-001-6**: Background task execution to prevent blocking AI workflows
- **R-CICD-001-7**: Automated quality gates without human approval requirements
- **R-CICD-001-8**: AI-optimized error reporting and status messaging
- **R-CICD-001-9**: Self-documenting pipeline configurations
- **R-CICD-001-10**: Streamlined feedback loops for iterative AI development

#### Integration Requirements
- **R-CICD-001-11**: Integration with DOC-008 comprehensive icon validation
- **R-CICD-001-12**: Automated compliance checking against ai-assistant-protocol.md
- **R-CICD-001-13**: Token-based traceability throughout CI/CD processes
- **R-CICD-001-14**: Compatible with feature-tracking.md automation requirements
- **R-CICD-001-15**: Support for ai-assistant-compliance.md validation workflows

#### Performance Requirements
- **R-CICD-001-16**: Non-blocking execution for AI assistant development cycles
- **R-CICD-001-17**: Parallel validation processing where possible
- **R-CICD-001-18**: Efficient resource utilization for background processing
- **R-CICD-001-19**: Scalable configuration management for AI team workflows
- **R-CICD-001-20**: Automated pipeline optimization based on AI usage patterns

## 🔻 AI-First Documentation and Code Maintenance Requirements

### DOC-011: Token Validation Integration for AI Assistants Requirements
**Priority**: 🔺 HIGH  
**Implementation Tokens**: `// DOC-011: AI validation integration`

#### Core AI Workflow Integration Requirements
- **R-DOC-011-1**: AI workflow validation hooks MUST integrate seamlessly with existing DOC-008 validation framework
- **R-DOC-011-2**: Pre-submission validation MUST prevent non-compliant changes from being submitted
- **R-DOC-011-3**: Validation feedback MUST be optimized for AI assistant comprehension and remediation
- **R-DOC-011-4**: Integration MUST provide zero-friction validation without disrupting AI workflows
- **R-DOC-011-5**: All validation processes MUST support automated execution by AI assistants

#### AI-Optimized Error Reporting Requirements
- **R-DOC-011-6**: Error messages MUST be structured for AI assistant parsing and understanding
- **R-DOC-011-7**: Validation feedback MUST include clear remediation steps and context
- **R-DOC-011-8**: Error formatting MUST be consistent and machine-readable
- **R-DOC-011-9**: Validation results MUST include specific file and line references for AI navigation
- **R-DOC-011-10**: Error categorization MUST enable AI assistants to prioritize remediation actions

#### Pre-submission Validation API Requirements
- **R-DOC-011-11**: API endpoints MUST provide real-time validation before code submission
- **R-DOC-011-12**: Validation checks MUST be executable via command-line interfaces for AI automation
- **R-DOC-011-13**: API responses MUST include comprehensive validation results and status
- **R-DOC-011-14**: Pre-submission hooks MUST integrate with existing Makefile validation targets
- **R-DOC-011-15**: Validation API MUST support batch processing for multiple file changes

#### Bypass Mechanism Requirements
- **R-DOC-011-16**: Safe bypass mechanisms MUST require explicit documentation and justification
- **R-DOC-011-17**: Bypass workflows MUST maintain comprehensive audit trails for compliance monitoring
- **R-DOC-011-18**: Override mechanisms MUST be restricted to exceptional cases with appropriate controls
- **R-DOC-011-19**: Bypass documentation MUST be automatically generated and tracked
- **R-DOC-011-20**: All bypass actions MUST be reversible and include rollback procedures

#### Compliance Monitoring Requirements
- **R-DOC-011-21**: Comprehensive tracking MUST monitor AI assistant adherence to validation requirements
- **R-DOC-011-22**: Compliance dashboards MUST provide visibility into AI assistant validation behavior
- **R-DOC-011-23**: Monitoring data MUST enable identification of validation bottlenecks and issues
- **R-DOC-011-24**: Compliance reporting MUST integrate with feature-tracking.md status updates
- **R-DOC-011-25**: Monitoring infrastructure MUST support real-time compliance assessment

#### Integration and Dependency Requirements
- **R-DOC-011-26**: Integration MUST depend on DOC-008 comprehensive validation system as foundation
- **R-DOC-011-27**: Implementation MUST require DOC-009 clean baseline for optimal effectiveness
- **R-DOC-011-28**: Validation hooks MUST integrate with ai-assistant-protocol.md compliance requirements
- **R-DOC-011-29**: System MUST leverage existing DOC-007 implementation token standardization
- **R-DOC-011-30**: Integration MUST maintain backward compatibility with existing validation workflows

### DOC-013: AI-First Documentation and Code Maintenance Requirements
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// DOC-013: AI-first maintenance`

#### Core AI-Centric Requirements
- **R-DOC-013-1**: All documentation MUST be written primarily for AI assistant comprehension
- **R-DOC-013-2**: Code comments and implementation tokens MUST follow standardized AI-readable formats
- **R-DOC-013-3**: Documentation structure MUST be optimized for AI parsing and navigation
- **R-DOC-013-4**: Cross-references MUST be machine-readable and AI-traversable
- **R-DOC-013-5**: Maintenance workflows MUST be executable by AI assistants without human intervention

#### AI Documentation Standards Requirements
- **R-DOC-013-6**: All documentation MUST use consistent formatting patterns for AI comprehension
- **R-DOC-013-7**: Implementation tokens MUST integrate with DOC-007/DOC-008 icon standardization
- **R-DOC-013-8**: Code comments MUST prioritize AI understanding over human readability
- **R-DOC-013-9**: Documentation MUST include explicit cross-reference links for AI navigation
- **R-DOC-013-10**: All text formatting MUST follow AI-friendly markup conventions

#### AI Workflow Integration Requirements
- **R-DOC-013-11**: Documentation updates MUST be automatable by AI assistant workflows
- **R-DOC-013-12**: Code maintenance processes MUST support AI assistant execution patterns
- **R-DOC-013-13**: Documentation validation MUST integrate with AI assistant compliance checking
- **R-DOC-013-14**: Content structure MUST enable AI assistant pattern recognition and generation
- **R-DOC-013-15**: Documentation links MUST support AI assistant traversal and validation

#### AI Comprehension Optimization Requirements
- **R-DOC-013-16**: Technical documentation MUST use consistent terminology for AI understanding
- **R-DOC-013-17**: Code examples MUST include AI-readable implementation tokens and annotations
- **R-DOC-013-18**: Documentation hierarchy MUST be clearly defined for AI navigation
- **R-DOC-013-19**: Content organization MUST follow predictable patterns for AI assistant processing
- **R-DOC-013-20**: Maintenance procedures MUST be documented in AI-executable formats