# BkpDir: Directory Archiving CLI Application

## Overview
BkpDir is a command-line application for macOS and Linux that archives directory hierarchies with configurable behavior. It supports full and incremental backups with customizable naming and exclusion patterns.

## Configuration File
- Configuration is stored in a YAML file at the root of the directory to be archived
- If the file is not present, default values are used

### Configuration Options
1. **Archive Directory Path**
   - Specifies where archives are stored
   - Default: `../.bkpdir` relative to current directory

2. **Use Current Directory Name**
   - Controls whether to include current directory name in the archive path
   - Default: `true`
   - Example: With directory 'HOME', archive path becomes '../.bkpdir/HOME/'

3. **Exclude Patterns**
   - List of patterns to exclude (similar to `.gitignore`)
   - Default: Excludes `.git/` and `vendor/` directories at the root of the source

## Commands

### 1. Create Full Archive
- Compresses current directory tree into an archive
- Excludes files matching patterns in configuration
- Usage: `bkpdir full [NOTE]`
- Archive naming format:
  - For Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm]=[GIT_BRANCH]=[GIT_SHORT_HASH][=[NOTE]].zip`
  - For non-Git repositories: `[PREFIX-][YYYY-MM-DD-hh-mm][=[NOTE]].zip`
- PREFIX is added if the source directory name is not already in the path
- NOTE is an optional positional argument provided by the user

### 2. Create Incremental Update
- Creates an update archive containing only files changed since the last full archive
- Usage: `bkpdir inc [NOTE]`
- Naming format: `[LAST_FULL_ARCHIVE_NAME]_update=[TIMESTAMP][=[GIT_INFO]][=[NOTE]].zip`
- Maintains relationships to original archive for tracking changes

### 3. List Archives
- Displays all archives associated with the current directory
- Shows full and incremental archives with timestamps

## Global Options
- **Dry-Run Mode**: When enabled, shows:
  - Files to be included in the archive
  - Files excluded based on patterns
  - Files changed since last backup (for incremental updates)
  - Archive filename that would be created
  - Target directory path
  - No actual archive is created

## Platform Compatibility
- Works on macOS and Linux systems
- Uses platform-independent path handling
- Preserves file permissions and ownership where applicable
- Handles file system differences between platforms
