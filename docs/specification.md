# BkpDir: Directory Archiving CLI Application

## Overview
BkpDir is a command-line application for macOS and Linux that archives directory hierarchies with configurable behavior. It supports full and incremental backups with customizable naming and exclusion patterns, as well as backup verification to ensure data integrity.

## Configuration File
- Configuration is stored in a YAML file named `.bkpdir.yml` at the root of the directory to be archived
- If the file is not present, default values are used

### Configuration Options
1. **Archive Directory Path**
   - Specifies where archives are stored
   - Default: `../.bkpdir` relative to current directory
   - YAML key: `archive_dir_path`

2. **Use Current Directory Name**
   - Controls whether to include current directory name in the archive path
   - Default: `true`
   - YAML key: `use_current_dir_name`
   - Example: With directory 'HOME', archive path becomes '../.bkpdir/HOME/'

3. **Exclude Patterns**
   - List of patterns to exclude from archiving
   - Uses doublestar glob pattern matching for powerful exclusion rules
   - Default: Excludes `.git/` and `vendor/` directories
   - YAML key: `exclude_patterns`
   - Supports patterns like `*.tmp`, `build/*`, and `**/node_modules/`

4. **Verification Settings**
   - Controls archive verification behavior
   - YAML key: `verification`
   - Sub-keys:
     - `verify_on_create`: bool - Automatically verify archives after creation (default: `false`)
     - `checksum_algorithm`: string - Algorithm used for checksums (default: `sha256`)

## Commands

### 1. Create Full Archive
- Compresses current directory tree into a ZIP archive
- Excludes files matching patterns in configuration
- **Only regular files are included in the archive and for checksums; directories are skipped.**
- Usage: `bkpdir full [NOTE]`
- Archive naming format:
  - For Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm]=[GIT_BRANCH]=[GIT_SHORT_HASH][=[NOTE]].zip`
  - For non-Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm][=[NOTE]].zip`
- PREFIX is added if `use_current_dir_name` is false (prefix is the source directory name)
- NOTE is an optional positional argument provided by the user
- Git information is only included if the directory is a valid Git repository
- With `--verify` flag: will verify archive after creation

### 2. Create Incremental Update
- Creates an update archive containing only files changed since the last full archive
- Based on file modification times compared to the last full archive
- Requires at least one full archive to exist
- Usage: `bkpdir inc [NOTE]`
- Naming format: `[LAST_FULL_ARCHIVE_NAME]_update=[TIMESTAMP][=[GIT_INFO]][=[NOTE]].zip`
- Git information is only included if the directory is a valid Git repository
- NOTE is an optional positional argument provided by the user
- With `--verify` flag: will verify archive after creation

### 3. List Archives
- Displays all archives associated with the current directory
- Usage: `bkpdir list`
- Shows filenames of all archives in the archive directory
- Includes verification status for each archive if available, indicated by:
  - `[VERIFIED]`: Archive passed verification
  - `[UNVERIFIED]`: Archive has not been verified
  - `[FAILED]`: Archive failed verification with specific errors

### 4. Verify Archive
- Verifies the integrity of an archive
- Usage: `bkpdir verify [ARCHIVE_NAME]`
- Performs the following checks:
  - Validates ZIP file structure
  - Checks for corrupted or incomplete archives
  - Verifies file headers and compression
  - With `--checksum` flag: verifies file contents against stored checksums
- Stores verification results for display in the `list` command
- Reports verification status with the following exit codes:
  - 0: Verification successful
  - 1: Verification failed (with error details)
- Verification Status includes:
  - Verification timestamp
  - Whether the archive passed basic verification
  - Whether checksums were verified (if applicable)
  - Any errors encountered during verification

## Global Options
- **Dry-Run Mode**: When enabled with `--dry-run` flag:
  - Shows files to be included in the archive
  - Displays the archive filename that would be created
  - No actual archive is created
  - Works with both full and incremental commands

## Implementation Details
- Archives are created using the standard ZIP format
- File permissions are preserved in the ZIP archive
- The application checks whether the directory is a Git repository and includes relevant Git information
- For incremental backups, finds the most recent full archive automatically
- Nested directories are preserved in the archive structure
- File paths in archives are relative to the source directory
- **Checksum generation uses absolute file paths for reading, but stores checksums using the base filename as the key.**
- **Directories are never included in the checksum process.**

### Verification Implementation
- Checksums are stored in a `.checksums` file within the archive
- **Checksums use the base filename as the key** (e.g., `test.txt`), matching the file entries in the archive. This ensures correct mapping during verification and is validated by tests.
- SHA-256 is used as the default checksum algorithm
- Checksum file format is JSON with filename as key and checksum as value
- When verifying incremental archives, only checks structure, not relationship to base archive
- Verification status is stored in a metadata file within the archive directory
- Output of verification displays:
  - Archive name
  - Verification status (pass/fail)
  - Timestamp of verification
  - Any errors found during verification
  - Summary of checksum verification (if performed)
- **Tests ensure correct handling of both valid and corrupted archives, that directories are skipped for checksums, and that verification status and checksum verification are robust.**

## Platform Compatibility
- Works on macOS and Linux systems
- Uses platform-independent path handling
- Preserves file permissions and ownership where applicable
- Handles file system differences between platforms

## Binary Compatibility Requirements
- **macOS Support**:
  - Must maintain two separate binaries:
    - ARM64 (Apple Silicon) binary
    - AMD64 (Intel) binary
  - Both binaries must be built and tested for each release
  - Binary naming format: `bkpdir-macos-[arch]`

- **Ubuntu Support**:
  - Single binary compatible with Ubuntu 20.04, 22.04, and 24.04
  - Binary must be built on Ubuntu 20.04 to ensure maximum compatibility
  - Binary naming format: `bkpdir-ubuntu20.04`
  - No version-specific optimizations that would break compatibility
  - Must maintain backward compatibility with Ubuntu 20.04 system libraries

- **Build Requirements**:
  - All binaries must be built with Go 1.21 or later
  - Build process must be automated via Makefile
  - Each binary must include compile date in version information
  - All binaries must pass verification tests on their target platforms
