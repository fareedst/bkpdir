## Data Objects

1. **Config**
   - `ArchiveDirPath`: string - Path where archives are stored
   - `UseCurrentDirName`: bool - Whether to use current directory name in archive path
   - `ExcludePatterns`: []string - Patterns for files to exclude
   - `Verification`: *VerificationConfig - Settings for verification

2. **VerificationConfig**
   - `VerifyOnCreate`: bool - Whether to automatically verify archives after creation
   - `ChecksumAlgorithm`: string - Algorithm used for checksums (sha256)

3. **Archive**
   - `Name`: string - Name of the archive file
   - `Path`: string - Full path to the archive file
   - `CreationTime`: time.Time - When the archive was created
   - `IsIncremental`: bool - Whether this is an incremental archive
   - `GitBranch`: string - Current git branch (if in git repo)
   - `GitHash`: string - Current git commit hash (if in git repo)
   - `Note`: string - Optional note for the archive
   - `BaseArchive`: string - Name of the base archive for incremental backups
   - `VerificationStatus`: *VerificationStatus - Result of the last verification, if any

4. **VerificationStatus**
   - `VerifiedAt`: time.Time - When the archive was last verified
   - `IsVerified`: bool - Whether the archive passed verification
   - `HasChecksums`: bool - Whether checksums were verified
   - `Errors`: []string - List of verification errors, if any

## Core Functions

1. **Configuration Management**
   - `DefaultConfig() *Config`: Creates default configuration
   - `LoadConfig(root string) (*Config, error)`: Loads config from YAML or uses defaults

2. **Git Integration**
   - `IsGitRepository(dir string) bool`: Checks if directory is a git repo
   - `GetGitBranch(dir string) string`: Gets current branch name
   - `GetGitShortHash(dir string) string`: Gets current commit hash

3. **File System Operations**
   - `ShouldExcludeFile(filePath string, patterns []string) bool`: Checks if file matches exclude patterns using doublestar glob matching
   - File walking is handled directly in `CreateFullArchive` and `CreateIncrementalArchive` using `filepath.Walk`
   - **Only regular files are included for archiving and checksums; directories are skipped.**

4. **Archive Management**
   - `GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, isGit, isIncremental bool, baseName string) string`: Generates archive filename
   - `ListArchives(archiveDir string) ([]Archive, error)`: Gets all archives in directory
   - `CreateFullArchive(cfg *Config, note string, dryRun bool) error`: Creates full archive
   - `CreateIncrementalArchive(cfg *Config, note string, dryRun bool) error`: Creates incremental archive

5. **Archive Verification**
   - `VerifyArchive(archivePath string) (*VerificationStatus, error)`: Verifies archive structure and integrity
   - `VerifyChecksums(archivePath string) (*VerificationStatus, error)`: Verifies file checksums against stored values
   - `GenerateChecksums(files []string, algorithm string) (map[string]string, error)`: Generates checksums for files (uses base filename as key, matching archive entries; skips directories; uses absolute paths for reading)
   - `StoreChecksums(archive *Archive, checksums map[string]string) error`: Stores checksums in archive
   - `ReadChecksums(archive *Archive) (map[string]string, error)`: Reads checksums from archive

## Main Application Structure

1. **CLI Interface**
   - Uses `cobra` for command-line interface
   - Global flags:
     - `--dry-run`: Show what would be done without creating archives
   - Commands:
     - `full [NOTE]`: Create full archive with optional note
     - `inc [NOTE]`: Create incremental archive with optional note
     - `list`: List all archives
     - `verify [ARCHIVE_NAME]`: Verify archive integrity
     - Flags for verify command:
       - `--checksum`: Include checksum verification

2. **Workflow Implementation**
   - For full archive:
     1. Load config
     2. Determine archive directory path
     3. Walk directory collecting files (excluding patterns, skipping directories)
     4. Get git info if applicable
     5. Generate archive name
     6. Create zip archive (or simulate in dry-run)
     7. If config.Verification.VerifyOnCreate is true, verify archive

   - For incremental archive:
     1. Load config
     2. Determine archive directory path
     3. Find most recent full archive
     4. Walk directory and collect files modified since the full archive (skipping directories)
     5. Get git info if applicable
     6. Generate incremental archive name
     7. Create zip archive (or simulate in dry-run)
     8. If config.Verification.VerifyOnCreate is true, verify archive

   - For listing archives:
     1. Load config
     2. Get archive directory path
     3. List all .zip files
     4. Display archive names with verification status if available (shows [VERIFIED], [FAILED], or [UNVERIFIED])

   - For verifying archives:
     1. Load config
     2. Get archive path
     3. Verify archive structure and integrity
     4. If --checksum flag is set, verify file checksums
     5. Store verification status
     6. Display verification results

3. **Utility Functions**
   - File exclusion uses `doublestar` for glob pattern matching
   - Archive naming follows format:
     - Full: `[prefix-]timestamp[=branch=hash][=note].zip`
     - Incremental: `baseName_update=timestamp[=branch=hash][=note].zip`
   - Checksum generation uses SHA-256 algorithm by default and stores checksums using the base filename as the key

## Testing Architecture

1. **Unit Tests**
   - `TestGenerateArchiveName`: Verifies archive naming for various scenarios
   - `TestShouldExcludeFile`: Tests file exclusion patterns
   - `TestDefaultConfig` and `TestLoadConfig`: Verify configuration management
   - `TestVerifyArchive`: Tests archive verification, including storing and loading verification status, and checksum verification (using base filename as key)
   - `TestVerifyChecksums`: Tests checksum verification
   - `TestGenerateChecksums`: Tests checksum generation, ensuring directories are skipped
   - `TestVerifyCorruptedArchive`: Ensures corrupted archives are detected and reported

2. **Integration Tests**
   - `TestFullCmdWithNote`: Tests full archive command with various notes
   - `TestIncCmdWithNote`: Tests incremental archive command with various notes
   - `TestCmdArgsValidation`: Validates command-line argument handling
   - `TestVerifyCmdWithChecksums`: Tests verify command with checksum option
   - `TestVerifyOnCreate`: Tests automatic verification after archive creation

3. **Test Approach**
   - Uses temporary directories for testing
   - Simulates user environment with test files and directories
   - Tests both dry-run and actual execution modes
   - Verifies correct behavior for edge cases and error conditions
   - Tests verification with both valid and corrupted archives
   - Ensures that checksum keys match archive file entries (base filename)
   - **Ensures directories are not included in checksums and that corrupted archives are detected**
