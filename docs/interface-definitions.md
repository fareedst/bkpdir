# Interface Definitions for Component Extraction

## Overview

This document defines the standardized interface contracts required for clean component extraction as identified in REFACTOR-001. These interfaces enable dependency inversion and prevent circular dependencies during the extraction process.

## Configuration Interfaces

### ConfigProvider
Primary interface for configuration management operations.

```go
type ConfigProvider interface {
    // Core configuration loading
    LoadConfig(root string) (*Config, error)
    LoadConfigValues(root string) (map[string]ConfigValue, error)
    
    // Configuration value extraction
    GetConfigValues(cfg *Config) []ConfigValue
    GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue
    
    // Configuration validation
    ValidateConfig(cfg *Config) error
}
```

### ConfigMerger
Interface for configuration merging and composition operations.

```go
type ConfigMerger interface {
    // Configuration merging
    MergeConfigs(dst, src *Config)
    MergeConfigValues(dst, src map[string]ConfigValue)
    
    // Path and file operations
    GetConfigSearchPaths() []string
    ExpandPath(path string) string
}
```

### ArchiveConfig
Configuration interface specific to archive operations.

```go
type ArchiveConfig interface {
    // Archive settings
    GetArchiveDirectory() string
    GetExclusionPatterns() []string
    GetUseCurrentDirName() bool
    GetIncludeGitInfo() bool
    GetShowGitDirtyStatus() bool
    
    // Verification settings
    GetVerificationSettings() VerificationSettings
    
    // Status codes
    GetArchiveStatusCodes() map[string]int
    
    // Format strings
    GetArchiveFormatStrings() map[string]string
    GetArchiveTemplateStrings() map[string]string
}
```

### BackupConfig
Configuration interface specific to backup operations.

```go
type BackupConfig interface {
    // Backup settings
    GetBackupDirectory() string
    GetUseCurrentDirNameForFiles() bool
    
    // File operation settings
    GetFileOperationSettings() FileOperationSettings
    
    // Status codes
    GetFileStatusCodes() map[string]int
    
    // Format strings
    GetBackupFormatStrings() map[string]string
    GetBackupTemplateStrings() map[string]string
}
```

### FormatterConfig
Configuration interface for output formatting.

```go
type FormatterConfig interface {
    // Format strings
    GetFormatStrings() map[string]string
    GetTemplateStrings() map[string]string
    
    // Patterns for data extraction
    GetPatterns() map[string]string
    
    // Error format strings
    GetErrorFormatStrings() map[string]string
}
```

### ErrorConfig
Configuration interface for error handling.

```go
type ErrorConfig interface {
    // Status codes
    GetStatusCodes() map[string]int
    GetErrorStatusCodes() map[string]int
    
    // Error format strings
    GetErrorFormatStrings() map[string]string
    GetErrorTemplateStrings() map[string]string
}
```

## Output Formatting Interfaces

### OutputFormatter
Primary interface for output formatting operations.

```go
type OutputFormatter interface {
    // Core formatting methods
    FormatCreatedArchive(path string) string
    FormatIdenticalArchive(path string) string
    FormatListArchive(path, creationTime string) string
    FormatConfigValue(name, value, source string) string
    FormatDryRunArchive(path string) string
    FormatError(message string) string
    
    // Backup formatting
    FormatCreatedBackup(path string) string
    FormatIdenticalBackup(path string) string
    FormatListBackup(path, creationTime string) string
    FormatDryRunBackup(path string) string
    
    // Template-based formatting
    FormatTemplate(templateStr string, data map[string]string) string
    ExtractPatternData(pattern, text string) map[string]string
    
    // Print methods
    PrintCreatedArchive(path string)
    PrintIdenticalArchive(path string)
    PrintListArchive(path, creationTime string)
    PrintConfigValue(name, value, source string)
    PrintDryRunArchive(path string)
    PrintError(message string)
    
    // Advanced formatting
    FormatWithExtraction(path string) string
    FormatListWithExtraction(path, creationTime string) string
}
```

### ErrorFormatter
Specialized interface for error formatting.

```go
type ErrorFormatter interface {
    // Error formatting
    FormatError(message string) string
    FormatErrorWithDetails(message string, details map[string]interface{}) string
    
    // Error printing
    PrintError(message string)
    PrintErrorWithDetails(message string, details map[string]interface{})
    
    // Specific error types
    FormatDiskFullError(err error) string
    FormatPermissionError(err error) string
    FormatDirectoryNotFound(err error) string
    FormatFileNotFound(err error) string
}
```

### TemplateFormatter
Interface for template-based formatting operations.

```go
type TemplateFormatter interface {
    // Template formatting
    FormatWithTemplate(input, pattern, tmplStr string) (string, error)
    FormatWithPlaceholders(format string, data map[string]string) string
    
    // Data extraction
    ExtractArchiveData(filename string) map[string]string
    ExtractBackupData(filename string) map[string]string
    ExtractConfigData(line string) map[string]string
    ExtractTimestampData(timestamp string) map[string]string
    
    // Template methods
    TemplateCreatedArchive(data map[string]string) string
    TemplateIdenticalArchive(data map[string]string) string
    TemplateListArchive(data map[string]string) string
    TemplateError(data map[string]string) string
}
```

## Git Integration Interfaces

### GitProvider
Primary interface for Git repository operations.

```go
type GitProvider interface {
    // Repository detection
    IsGitRepository(dir string) bool
    
    // Git information retrieval
    GetCurrentBranch(dir string) (string, error)
    GetCurrentCommitHash(dir string) (string, error)
    IsWorkingTreeClean(dir string) (bool, error)
    
    // Composite operations
    GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool)
    
    // Git command execution
    ExecuteGitCommand(dir string, args ...string) (string, error)
}
```

### GitInfoProvider
Interface for Git information extraction.

```go
type GitInfoProvider interface {
    // Git metadata
    GetBranchName(dir string) (string, error)
    GetCommitHash(dir string) (string, error)
    GetShortHash(dir string) (string, error)
    
    // Status information
    GetWorkingTreeStatus(dir string) (bool, error)
    GetStatusSummary(dir string) (GitStatus, error)
}
```

## Error Handling Interfaces

### ErrorHandler
Primary interface for error handling operations.

```go
type ErrorHandler interface {
    // Error handling
    HandleArchiveError(err error, cfg ErrorConfig, formatter ErrorFormatter) int
    HandleBackupError(err error, cfg ErrorConfig, formatter ErrorFormatter) int
    HandleGenericError(err error, cfg ErrorConfig, formatter ErrorFormatter) int
    
    // Error creation
    NewArchiveError(message string, statusCode int) *ArchiveError
    NewArchiveErrorWithCause(message string, statusCode int, err error) *ArchiveError
    NewBackupError(message string, statusCode int, operation, path string) *BackupError
    NewBackupErrorWithCause(message string, statusCode int, operation, path string, err error) *BackupError
    
    // Error analysis
    IsUserError(err error) bool
    IsSystemError(err error) bool
    GetErrorCategory(err error) ErrorCategory
}
```

### ResourceManager
Interface for resource management and cleanup.

```go
type ResourceManager interface {
    // Resource management
    AddResource(resource Resource)
    AddTempFile(path string)
    AddTempDir(path string)
    RemoveResource(resource Resource)
    
    // Cleanup operations
    Cleanup() error
    CleanupWithPanicRecovery() error
    
    // Resource queries
    GetManagedResources() []Resource
    HasResource(resource Resource) bool
}
```

## File Operations Interfaces

### FileComparator
Interface for file comparison operations.

```go
type FileComparator interface {
    // File comparison
    CompareFiles(file1, file2 string) (bool, error)
    CompareFilesWithContext(ctx context.Context, file1, file2 string) (bool, error)
    CompareFileContents(f1, f2 *os.File) (bool, error)
    
    // Advanced comparison
    CompareFileContentsWithContext(ctx context.Context, f1, f2 *os.File) (bool, error)
    GetFileHash(filePath string) (string, error)
    GetFileHashWithContext(ctx context.Context, filePath string) (string, error)
}
```

### ExclusionManager
Interface for file exclusion pattern management.

```go
type ExclusionManager interface {
    // Exclusion checks
    ShouldExclude(path string, excludePatterns []string) bool
    FilterFiles(files []string, excludePatterns []string) []string
    
    // Pattern management
    LoadExclusionPatterns() ([]string, error)
    AddExclusionPattern(pattern string)
    RemoveExclusionPattern(pattern string)
    
    // Pattern validation
    ValidatePattern(pattern string) error
    CompilePatterns(patterns []string) ([]*regexp.Regexp, error)
}
```

### FileOperations
Interface for general file operations.

```go
type FileOperations interface {
    // Directory operations
    SafeMkdirAll(path string, perm os.FileMode) error
    ValidateDirectoryPath(path string) error
    
    // File operations
    CopyFile(src, dst string) error
    CopyFileWithContext(ctx context.Context, src, dst string) error
    ValidateFilePath(path string) error
    
    // Atomic operations
    AtomicWriteFile(path string, data []byte, rm ResourceManager) error
    AtomicCopyFile(src, dst string, rm ResourceManager) error
}
```

## Archive Operations Interfaces

### ArchiveManager
Primary interface for archive operations.

```go
type ArchiveManager interface {
    // Archive creation
    CreateArchive(cfg ArchiveConfig) error
    CreateFullArchive(cfg ArchiveConfig) error
    CreateIncrementalArchive(cfg ArchiveConfig) error
    CreateArchiveWithContext(ctx context.Context, cfg ArchiveConfig) error
    
    // Archive listing
    ListArchives(archiveDir string) ([]Archive, error)
    FindLatestFullArchive(archiveDir string) (*Archive, error)
    
    // Archive naming
    GenerateArchiveName(cfg ArchiveNameConfig) string
    ParseArchiveName(name string) (*ArchiveMetadata, error)
    
    // Archive verification
    VerifyArchive(archivePath string) (*VerificationStatus, error)
    VerifyArchiveWithChecksum(archivePath string) (*VerificationStatus, error)
}
```

### ArchiveCreator
Interface for archive creation operations.

```go
type ArchiveCreator interface {
    // Archive creation
    CreateZipArchive(sourceDir, archivePath string, files []string) error
    CreateZipArchiveWithContext(ctx context.Context, sourceDir, archivePath string, files []string) error
    
    // File collection
    CollectFiles(ctx context.Context, cwd string, excludePatterns []string) ([]string, error)
    CollectModifiedFiles(cwd string, since time.Time, excludePatterns []string) ([]string, error)
    
    // Archive validation
    ValidateArchiveCreation(cfg ArchiveConfig) error
    PrepareArchiveDirectory(cfg ArchiveConfig, cwd string, dryRun bool) (string, error)
}
```

### ArchiveVerifier
Interface for archive verification operations.

```go
type ArchiveVerifier interface {
    // Verification operations
    VerifyArchive(archivePath string) (*VerificationStatus, error)
    VerifyArchiveWithChecksum(archivePath string) (*VerificationStatus, error)
    VerifyArchiveIntegrity(archivePath string) (*VerificationStatus, error)
    
    // Verification status management
    StoreVerificationStatus(archive *Archive, status *VerificationStatus) error
    LoadVerificationStatus(archivePath string) (*VerificationStatus, error)
    
    // Verification utilities
    CalculateArchiveChecksum(archivePath string) (string, error)
    ValidateArchiveStructure(archivePath string) error
}
```

## Backup Operations Interfaces

### BackupManager
Primary interface for backup operations.

```go
type BackupManager interface {
    // Backup creation
    CreateFileBackup(cfg BackupConfig, filePath string, note string, dryRun bool) error
    CreateFileBackupWithContext(ctx context.Context, cfg BackupConfig, filePath string, note string, dryRun bool) error
    CreateFileBackupEnhanced(opts BackupOptions) error
    
    // Backup listing
    ListFileBackups(backupDir string, baseFilename string) ([]BackupInfo, error)
    ListFileBackupsWithContext(ctx context.Context, backupDir, sourceFile string) ([]*Backup, error)
    
    // Backup validation
    CheckForIdenticalFileBackup(filePath, backupDir, baseFilename string) (bool, string, error)
    ValidateFileForBackup(filePath string) error
    
    // Backup utilities
    GenerateBackupPath(cfg BackupConfig, filePath, note string) (string, error)
    GetMostRecentBackup(backupDir, sourceFile string) (*Backup, error)
}
```

### BackupCreator
Interface for backup creation operations.

```go
type BackupCreator interface {
    // Backup creation
    CreateBackup(opts BackupOptions) error
    CreateBackupWithCleanup(opts BackupOptions) error
    
    // Backup path management
    GenerateBackupName(sourcePath, timestamp, note string) string
    DetermineBackupPath(cfg BackupConfig, filePath string) (string, error)
    
    // Backup validation
    ValidateBackupOperation(opts BackupOptions) error
    CheckBackupDirectorySpace(backupDir string, requiredSpace int64) error
}
```

## CLI Orchestration Interfaces

### CLIOrchestrator
Primary interface for CLI command orchestration.

```go
type CLIOrchestrator interface {
    // Command execution
    ExecuteCommand(args []string) error
    HandleCreateCommand(args []string) error
    HandleListCommand(args []string) error
    HandleVerifyCommand(args []string) error
    HandleBackupCommand(args []string) error
    HandleConfigCommand(args []string) error
    
    // Command utilities
    InitializeServices(cfg *Config) (*ApplicationServices, error)
    ValidateCommandArgs(cmd string, args []string) error
    SetupCommandContext(cfg *Config) (context.Context, error)
}
```

### ApplicationServices
Aggregate interface for all application services.

```go
type ApplicationServices interface {
    // Core services
    ConfigProvider() ConfigProvider
    OutputFormatter() OutputFormatter
    ErrorHandler() ErrorHandler
    
    // Operational services
    ArchiveManager() ArchiveManager
    BackupManager() BackupManager
    GitProvider() GitProvider
    
    // Utility services
    FileComparator() FileComparator
    ExclusionManager() ExclusionManager
    FileOperations() FileOperations
    ResourceManager() ResourceManager
}
```

## Supporting Types and Enums

### Configuration Types

```go
type VerificationSettings struct {
    VerifyOnCreate    bool
    ChecksumAlgorithm string
}

type FileOperationSettings struct {
    UseAtomicOperations bool
    TempFilePrefix      string
    BackupFilePermissions os.FileMode
}

type ArchiveSettings struct {
    UseCompression    bool
    CompressionLevel  int
    MaxArchiveSize    int64
    IncludeMetadata   bool
}

type BackupSettings struct {
    MaxBackups         int
    BackupRotation     bool
    CompressionEnabled bool
}
```

### Error Types

```go
type ErrorCategory int

const (
    UserError ErrorCategory = iota
    SystemError
    ConfigurationError
    PermissionError
    DiskSpaceError
    NetworkError
)

type GitStatus struct {
    IsClean      bool
    Branch       string
    Hash         string
    ShortHash    string
    HasUntracked bool
    HasModified  bool
    HasStaged    bool
}
```

### Archive Types

```go
type ArchiveMetadata struct {
    Name          string
    Prefix        string
    Timestamp     time.Time
    GitBranch     string
    GitHash       string
    Note          string
    IsIncremental bool
    BaseArchive   string
}

type ArchiveNameConfig struct {
    Prefix             string
    Timestamp          string
    GitBranch          string
    GitHash            string
    GitIsClean         bool
    ShowGitDirtyStatus bool
    Note               string
    IsGit              bool
    IsIncremental      bool
    BaseName           string
}
```

### Backup Types

```go
type BackupOptions struct {
    Context   context.Context
    Config    BackupConfig
    Formatter OutputFormatter
    FilePath  string
    Note      string
    DryRun    bool
}

type BackupInfo struct {
    Name         string
    Path         string
    CreationTime time.Time
    Size         int64
}

type Backup struct {
    Name         string
    Path         string
    CreationTime time.Time
    SourceFile   string
    Note         string
}
```

## Implementation Guidelines

### Interface Segregation
- Keep interfaces focused on specific concerns
- Avoid large, monolithic interfaces
- Group related methods logically

### Dependency Inversion
- Depend on abstractions, not concrete types
- Inject dependencies through interfaces
- Use constructor injection for required dependencies

### Error Handling
- Return errors as values, not through interfaces
- Use structured errors with context
- Provide error categorization for appropriate handling

### Context Support
- Accept context.Context for cancellable operations
- Respect context cancellation in long-running operations
- Use context for request-scoped values when appropriate

### Resource Management
- Use ResourceManager for cleanup operations
- Implement proper resource disposal patterns
- Handle panics gracefully in cleanup operations

---

**Document Version**: 1.0  
**Created**: 2025-01-02  
**Last Updated**: 2025-01-02  
**Related**: REFACTOR-001, extraction-dependencies.md 