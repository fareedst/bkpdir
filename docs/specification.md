# BkpDir: Directory Archiving CLI Application

## Overview
BkpDir is a command-line application for macOS and Linux that archives directory hierarchies with configurable behavior. It supports full and incremental backups with customizable naming and exclusion patterns.

## Configuration File
- Configuration is stored in a YAML file named `.bkpdir.yaml` at the root of the directory to be archived
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

## Commands

### 1. Create Full Archive
- Compresses current directory tree into a ZIP archive
- Excludes files matching patterns in configuration
- Usage: `bkpdir full [NOTE]`
- Archive naming format:
  - For Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm]=[GIT_BRANCH]=[GIT_SHORT_HASH][=[NOTE]].zip`
  - For non-Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm][=[NOTE]].zip`
- PREFIX is added if `use_current_dir_name` is false (prefix is the source directory name)
- NOTE is an optional positional argument provided by the user
- Git information is only included if the directory is a valid Git repository

### 2. Create Incremental Update
- Creates an update archive containing only files changed since the last full archive
- Based on file modification times compared to the last full archive
- Requires at least one full archive to exist
- Usage: `bkpdir inc [NOTE]`
- Naming format: `[LAST_FULL_ARCHIVE_NAME]_update=[TIMESTAMP][=[GIT_INFO]][=[NOTE]].zip`
- Git information is only included if the directory is a valid Git repository
- NOTE is an optional positional argument provided by the user

### 3. List Archives
- Displays all archives associated with the current directory
- Usage: `bkpdir list`
- Shows filenames of all archives in the archive directory

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

## Platform Compatibility
- Works on macOS and Linux systems
- Uses platform-independent path handling
- Preserves file permissions and ownership where applicable
- Handles file system differences between platforms
