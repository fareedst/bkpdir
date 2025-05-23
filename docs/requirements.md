# Architecture and Specification Requirements Traceability

This document maps code components to their corresponding architecture requirements and specification requirements.

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

**Example Usage**:
```go
// Create default configuration
cfg := DefaultConfig()

// Load configuration from YAML file
cfg, err := LoadConfig(".")
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

**Example Usage**:
```go
// Check if file should be excluded
patterns := []string{".git/", "vendor/", "*.tmp"}
if ShouldExcludeFile("path/to/file.tmp", patterns) {
    fmt.Println("File should be excluded")
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
  - Error Cases:
    - Invalid configuration
    - File system errors
    - Archive creation errors
    - Verification failures

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

// Create full archive
err = CreateFullArchive(cfg, "backup note", false)
if err != nil {
    log.Fatal(err)
}

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
- Commands:
  - `full`: `fullCmd()`
    - Spec: "Create full archive with optional note"
  - `inc`: `incCmd()`
    - Spec: "Create incremental archive with optional note"
  - `list`: `listCmd()`
    - Spec: "List all archives"
  - `verify`: `verifyCmd()`
    - Spec: "Verify archive integrity"

### Workflow Implementation
**Implementation**: `archive.go` and `verify.go`
**Specification Requirements**:
- Full archive workflow: `CreateFullArchive()`
  - Spec: "Compresses current directory tree into ZIP archive"
- Incremental archive workflow: `CreateIncrementalArchive()`
  - Spec: "Creates update archive containing only changed files"
- Archive listing workflow: `ListArchives()`
  - Spec: "Displays all archives with verification status"
- Archive verification workflow: `VerifyArchive()` and `VerifyChecksums()`
  - Spec: "Verifies archive integrity and optionally checksums"

### Utility Functions
**Implementation**: Various files
**Specification Requirements**:
- File exclusion: `ShouldExcludeFile()` in `archive.go`
  - Spec: "Uses doublestar glob pattern matching"
- Archive naming: `GenerateArchiveName()` in `archive.go`
  - Spec: "Follows specified naming format"
- Checksum generation: `GenerateChecksums()` in `verify.go`
  - Spec: "Uses SHA-256 by default"

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
- `TestVerifyArchive`: Tests archive verification
  - Spec: "Validates archive structure verification"
- `TestVerifyChecksums`: Tests checksum verification
  - Spec: "Validates file content verification"
- `TestGenerateChecksums`: Tests checksum generation
  - Spec: "Validates checksum generation and storage"
- `TestVerifyCorruptedArchive`: Tests corruption detection
  - Spec: "Validates error handling for corrupted archives"

### Integration Tests
**Implementation**: `*_test.go` files
**Specification Requirements**:
- `TestFullCmdWithNote`: Tests full archive command
  - Spec: "Validates full archive creation with notes"
- `TestIncCmdWithNote`: Tests incremental archive command
  - Spec: "Validates incremental archive creation with notes"
- `TestCmdArgsValidation`: Tests command-line arguments
  - Spec: "Validates command-line interface"
- `TestVerifyCmdWithChecksums`: Tests verify command
  - Spec: "Validates archive verification workflow"
- `TestVerifyOnCreate`: Tests automatic verification
  - Spec: "Validates automatic verification after creation" 