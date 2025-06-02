# BkpDir: Directory Archiving and File Backup CLI Application

## Overview
BkpDir is a command-line application for macOS and Linux that creates ZIP-based archives of directories and backups of individual files. It supports Git integration, customizable naming patterns, file exclusion patterns, maintains a history of directory archives and file backups, and provides robust error handling with automatic resource cleanup. It also features configurable printf-style and template-based output formatting for enhanced user experience.

> **Important**: This document describes the user-facing features and behaviors. For immutable specifications that cannot be changed without a major version bump, see [Immutable Specifications](immutable.md).

## Documentation Navigation

### For Users
- Start with this [Specification](specification.md) document
- Refer to [Immutable Specifications](immutable.md) for core behaviors that cannot change

### For Developers
- Begin with [Architecture](architecture.md) for system design
- Follow [Requirements](requirements.md) for implementation details
- Use [Testing](testing.md) for test coverage requirements

### For Contributors
- Review [Immutable Specifications](immutable.md) first to understand constraints
- Follow [Testing](testing.md) requirements for all changes
- Ensure all code passes linting requirements before submission

### Document Maintenance
- Keep [Specification](specification.md) and [Immutable Specifications](immutable.md) in sync
- Update [Requirements](requirements.md) with new features
- Maintain test coverage as per [Testing](testing.md)
- All changes must preserve existing functionality per [Immutable Specifications](immutable.md)

## Quality Assurance and Code Standards

### Linting Requirements
- All Go code must pass `revive` linter checks before commit
- Linting configuration is maintained in `.revive.toml`
- Run linting with `make lint` command
- Code must follow Go best practices and naming conventions
- All errors must be properly handled (no unhandled return values)

### Error Handling Standards
- All archive and backup operations return structured errors with status codes
- Enhanced disk space detection for various storage conditions
- Panic recovery mechanisms prevent application crashes
- Context support for operation cancellation and timeouts
- Comprehensive error logging without exposing sensitive information

### Resource Management
- Automatic cleanup of temporary files and directories
- Thread-safe resource tracking for concurrent operations
- Atomic file operations to prevent data corruption
- No resource leaks in any error scenario
- Comprehensive cleanup testing and verification

## Configuration Discovery
- Configuration files are discovered using a configurable search path
- The search path is controlled by the `BKPDIR_CONFIG` environment variable
- If `BKPDIR_CONFIG` is not set, the default search path is hard-coded as: `./.bkpdir.yml:~/.bkpdir.yml`
- Configuration files are processed in order, with values from earlier files taking precedence
- If multiple configuration files exist, settings in earlier files override settings in later files

### Environment Variable: BKPDIR_CONFIG
- Specifies a colon-separated list of configuration file paths to search
- Example: `BKPDIR_CONFIG="/etc/bkpdir.yml:~/.config/bkpdir.yml:./.bkpdir.yml"`
- Paths can be absolute or relative
- Relative paths are resolved from the current working directory
- Home directory (`~`) expansion is supported

## Configuration File
- Configuration is stored in YAML files with names specified by the configuration discovery system
- If no configuration files are found, default values are used (see [Immutable Specifications](immutable.md#configuration-defaults))
- Configuration files use the `.yml` extension by convention

### Configuration Options
1. **Archive Directory Path**
   - Specifies where archives are stored
   - Default: `../.bkpdir` relative to current directory
   - YAML key: `archive_dir_path`
   - Archives maintain the source directory's name in the archive path

2. **Use Current Directory Name**
   - Controls whether to include current directory name in the archive path
   - Default: `true`
   - YAML key: `use_current_dir_name`
   - Example: With directory 'myproject', archive path becomes '../.bkpdir/myproject-2025-05-12-13-49.zip'

3. **Exclude Patterns**
   - List of glob patterns to exclude from archiving
   - Default: `[".git/", "vendor/"]`
   - YAML key: `exclude_patterns`
   - Uses doublestar glob pattern matching
   - Example: `["*.tmp", "node_modules/", ".DS_Store"]`

4. **Git Integration**
   - Automatically detects Git repositories using `git rev-parse --is-inside-work-tree`
   - Extracts current branch name using `git rev-parse --abbrev-ref HEAD`
   - Extracts commit hash using `git rev-parse --short HEAD`
   - Includes Git information in archive names when `include_git_info` is enabled
   - Gracefully handles non-Git directories by returning empty Git info
   - Configurable via `include_git_info` setting in configuration
   - Works with both clean and dirty working directories

5. **File Backup Configuration**
   - **Backup Directory Path**
     - Specifies where file backups are stored
     - Default: `../.bkpdir` relative to current directory
     - YAML key: `backup_dir_path`
     - Backups maintain the source file's directory structure in the backup path
   - **Use Current Directory Name for Files**
     - Controls whether to include current directory name in the file backup path
     - Default: `true`
     - YAML key: `use_current_dir_name_for_files`
     - Example: With file 'cmd/bkpdir/main.go', backup path becomes '../.bkpdir/cmd/bkpdir/main.go-2025-05-12-13-49'

6. **Verification Configuration**
   - Controls archive verification behavior
   - YAML key: `verification`
   - Sub-options:
     - `verify_on_create`: Automatically verify archives after creation (default: false)
     - `checksum_algorithm`: Algorithm used for checksums (default: "sha256")

7. **Status Code Configuration**
   - Configures exit status codes returned for different application conditions
   - Status codes have specific defaults if not specified (see [Immutable Specifications](immutable.md#configuration-defaults))
   - YAML keys for directory operation status codes:
     - `status_created_archive`: Exit code when a new archive is successfully created (default: 0)
     - `status_failed_to_create_archive_directory`: Exit code when archive directory creation fails (default: 31)
     - `status_directory_is_identical_to_existing_archive`: Exit code when directory is identical to most recent archive (default: 0)
     - `status_directory_not_found`: Exit code when source directory does not exist (default: 20)
     - `status_invalid_directory_type`: Exit code when source path is not a directory (default: 21)
     - `status_permission_denied`: Exit code when directory access is denied (default: 22)
     - `status_disk_full`: Exit code when disk space is insufficient (default: 30)
     - `status_config_error`: Exit code when configuration is invalid (default: 10)
   - YAML keys for file operation status codes:
     - `status_created_backup`: Exit code when a new file backup is successfully created (default: 0)
     - `status_failed_to_create_backup_directory`: Exit code when backup directory creation fails (default: 31)
     - `status_file_is_identical_to_existing_backup`: Exit code when file is identical to most recent backup (default: 0)
     - `status_file_not_found`: Exit code when source file does not exist (default: 20)
     - `status_invalid_file_type`: Exit code when source file is not a regular file (default: 21)
   - Example configuration:
     ```yaml
     # Directory operation status codes
     status_created_archive: 0
     status_failed_to_create_archive_directory: 31
     status_directory_is_identical_to_existing_archive: 0
     status_directory_not_found: 20
     status_invalid_directory_type: 21
     status_permission_denied: 22
     status_disk_full: 30
     status_config_error: 10
     
     # File operation status codes
     status_created_backup: 0
     status_failed_to_create_backup_directory: 31
     status_file_is_identical_to_existing_backup: 0
     status_file_not_found: 20
     status_invalid_file_type: 21
     ```

8. **Output Format Configuration**
   - Configures printf-style format strings for all standard output
   - Configures template-based formatting with named placeholders and regex patterns
   - Format strings support text highlighting and structure formatting
   - Templates support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
   - All user-facing text is extracted from code into configuration data
   - Format strings have specific defaults if not specified (see [Immutable Specifications](immutable.md#configuration-defaults))
   - YAML keys for printf-style format strings for directory operations:
     - `format_created_archive`: Format for successful archive creation messages (default: "Created archive: %s\n")
     - `format_identical_archive`: Format for identical directory messages (default: "Directory is identical to existing archive: %s\n")
     - `format_list_archive`: Format for archive listing entries (default: "%s (created: %s)\n")
     - `format_config_value`: Format for configuration value display (default: "%s: %s (source: %s)\n")
     - `format_dry_run_archive`: Format for dry-run archive messages (default: "Would create archive: %s\n")
     - `format_error`: Format for error messages (default: "Error: %s\n")
   - YAML keys for printf-style format strings for file operations:
     - `format_created_backup`: Format for successful file backup creation messages (default: "Created backup: %s\n")
     - `format_identical_backup`: Format for identical file messages (default: "File is identical to existing backup: %s\n")
     - `format_list_backup`: Format for file backup listing entries (default: "%s (created: %s)\n")
     - `format_dry_run_backup`: Format for dry-run file backup messages (default: "Would create backup: %s\n")
   - YAML keys for template-based format strings for directory operations:
     - `template_created_archive`: Template for successful archive creation messages (default: "Created archive: %{path}\n")
     - `template_identical_archive`: Template for identical directory messages (default: "Directory is identical to existing archive: %{path}\n")
     - `template_list_archive`: Template for archive listing entries (default: "%{path} (created: %{creation_time})\n")
     - `template_config_value`: Template for configuration value display (default: "%{name}: %{value} (source: %{source})\n")
     - `template_dry_run_archive`: Template for dry-run archive messages (default: "Would create archive: %{path}\n")
     - `template_error`: Template for error messages (default: "Error: %{message}\n")
   - YAML keys for template-based format strings for file operations:
     - `template_created_backup`: Template for successful file backup creation messages (default: "Created backup: %{path}\n")
     - `template_identical_backup`: Template for identical file messages (default: "File is identical to existing backup: %{path}\n")
     - `template_list_backup`: Template for file backup listing entries (default: "%{path} (created: %{creation_time})\n")
     - `template_dry_run_backup`: Template for dry-run file backup messages (default: "Would create backup: %{path}\n")
   - YAML keys for regex patterns:
     - `pattern_archive_filename`: Named regex for parsing archive filenames (default: `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`)
     - `pattern_backup_filename`: Named regex for parsing file backup filenames (default: `(?P<filename>[^/]+)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<note>.+))?`)
     - `pattern_config_line`: Named regex for parsing configuration display lines (default: `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`)
     - `pattern_timestamp`: Named regex for parsing timestamps (default: `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`)
   - Format strings support ANSI color codes and text formatting for enhanced readability
   - Template strings support conditional formatting and advanced text processing
   - Template-based formatting provides rich data extraction from filenames using named regex groups
   - Templates support both Go text/template syntax for complex logic and simple placeholder replacement
   - Error messages in templates can include operation context for enhanced debugging information
   - Example configuration:
     ```yaml
     # Printf-style formatting for directory operations
     format_created_archive: "\033[32m✓ Created archive:\033[0m %s\n"
     format_identical_archive: "\033[33m≡ Directory is identical to existing archive:\033[0m %s\n"
     format_list_archive: "\033[36m%s\033[0m (created: \033[90m%s\033[0m)\n"
     format_config_value: "\033[1m%s:\033[0m %s \033[90m(source: %s)\033[0m\n"
     format_dry_run_archive: "\033[35m⚠ Would create archive:\033[0m %s\n"
     format_error: "\033[31m✗ Error:\033[0m %s\n"
     
     # Printf-style formatting for file operations
     format_created_backup: "\033[32m✓ Created backup:\033[0m %s\n"
     format_identical_backup: "\033[33m≡ File is identical to existing backup:\033[0m %s\n"
     format_list_backup: "\033[36m%s\033[0m (created: \033[90m%s\033[0m)\n"
     format_dry_run_backup: "\033[35m⚠ Would create backup:\033[0m %s\n"
     
     # Template-based formatting with named placeholders and data extraction
     template_created_archive: "\033[32m✓ Created archive:\033[0m %{path}\n"
     template_identical_archive: "\033[33m≡ Directory %{prefix} is identical to archive from %{year}-%{month}-%{day}:\033[0m %{path}\n"
     template_list_archive: "\033[36m%{path}\033[0m (created: \033[90m%{creation_time}\033[0m) %{note}\n"
     template_config_value: "\033[1m%{name}:\033[0m %{value} \033[90m(from %{source})\033[0m\n"
     template_dry_run_archive: "\033[35m⚠ Would create archive for %{prefix} on %{year}-%{month}-%{day}:\033[0m %{path}\n"
     template_error: "\033[31m✗ Error in %{operation}:\033[0m %{message}\n"
     
     # Template-based formatting for file operations
     template_created_backup: "\033[32m✓ Created backup:\033[0m %{path}\n"
     template_identical_backup: "\033[33m≡ File %{filename} is identical to backup from %{year}-%{month}-%{day}:\033[0m %{path}\n"
     template_list_backup: "\033[36m%{path}\033[0m (created: \033[90m%{creation_time}\033[0m) %{note}\n"
     template_dry_run_backup: "\033[35m⚠ Would create backup for %{filename} on %{year}-%{month}-%{day}:\033[0m %{path}\n"
     
           # Named regex patterns for data extraction
      pattern_archive_filename: "(?P<prefix>[^-]*)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\\.zip"
      pattern_backup_filename: "(?P<filename>[^/]+)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<note>.+))?"
      pattern_timestamp: "(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})\\s+(?P<hour>\\d{2}):(?P<minute>\\d{2}):(?P<second>\\d{2})"
     
     # Printf-style formatting for file operations
     format_created_backup: "\033[32m✓ Created backup:\033[0m %s\n"
     format_identical_backup: "\033[33m≡ File is identical to existing backup:\033[0m %s\n"
     format_list_backup: "\033[36m%s\033[0m (created: \033[90m%s\033[0m)\n"
     format_dry_run_backup: "\033[35m⚠ Would create backup:\033[0m %s\n"
     
     # Template-based formatting for directory operations
     template_created_archive: "\033[32m✓ Created archive:\033[0m %{path}\n"
     template_identical_archive: "\033[33m≡ Directory %{prefix} is identical to archive from %{year}-%{month}-%{day} (%{branch}@%{hash}):\033[0m %{path}\n"
     template_list_archive: "\033[36m%{path}\033[0m (created: \033[90m%{creation_time}\033[0m) %{note}\n"
     template_config_value: "\033[1m%{name}:\033[0m %{value} \033[90m(from %{source})\033[0m\n"
     template_dry_run_archive: "\033[35m⚠ Would create archive for %{prefix} on %{year}-%{month}-%{day}:\033[0m %{path}\n"
     template_error: "\033[31m✗ Error in %{operation}:\033[0m %{message}\n"
     
     # Template-based formatting for file operations
     template_created_backup: "\033[32m✓ Created backup:\033[0m %{path}\n"
     template_identical_backup: "\033[33m≡ File %{filename} is identical to backup from %{year}-%{month}-%{day}:\033[0m %{path}\n"
     template_list_backup: "\033[36m%{path}\033[0m (created: \033[90m%{creation_time}\033[0m) %{note}\n"
     template_dry_run_backup: "\033[35m⚠ Would create backup for %{filename} on %{year}-%{month}-%{day}:\033[0m %{path}\n"
     
     # Named regex patterns for data extraction
     pattern_archive_filename: "(?P<prefix>[^-]*)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\\.zip"
     pattern_backup_filename: "(?P<filename>[^/]+)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<note>.+))?"
     pattern_timestamp: "(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})\\s+(?P<hour>\\d{2}):(?P<minute>\\d{2}):(?P<second>\\d{2})"
     ```

## Commands

### 1. Create Full Archive
- Creates a complete ZIP archive of the current directory
- Usage: `bkpdir full [NOTE]`
- Before creating an archive:
  - Compares the directory with its most recent archive using directory tree comparison
  - If the directory is identical to the most recent archive:
    - Reports the existing archive path using `format_identical_archive` or `template_identical_archive` configuration
    - Template formatting can extract and display rich information from archive filename using `pattern_archive_filename`
    - Default format: "Directory is identical to existing archive: [PATH]"
    - Template format can show: "Directory [prefix] is identical to archive from [year]-[month]-[day] ([branch]@[hash]): [path]"
    - Exits with `status_directory_is_identical_to_existing_archive` status code
- When a new archive is created:
  - Reports success using `format_created_archive` or `template_created_archive` configuration
  - Default format: "Created archive: [PATH]"
  - Template format can include extracted directory name and timestamp information
  - Exits with `status_created_archive` status code
- Archive naming format: `[PREFIX-]YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip`
  - PREFIX is the current directory name (if `use_current_dir_name` is true)
  - YYYY-MM-DD-hh-mm is the timestamp of the archive
  - BRANCH and HASH are Git information (if in a Git repository and `include_git_info` is true)
  - NOTE is an optional note appended with an equals sign
- The archive excludes files and directories matching patterns in `exclude_patterns`
- NOTE is an optional positional argument provided by the user
- All output uses configurable printf-style format strings or template-based formatting for consistency and customization

### 2. Create Incremental Archive
- Creates an incremental ZIP archive containing only files changed since the last full archive
- Usage: `bkpdir inc [NOTE]`
- Requires an existing full archive as a base
- Archive naming format: `BASENAME_update=YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip`
  - BASENAME is the name of the base full archive
  - Timestamp and Git info follow the same format as full archives
- Only includes files modified since the base archive creation time
- Reports success using the same formatting configuration as full archives
- Exits with `status_created_archive` status code on success

### 3. List Archives
- Displays all archives in the archive directory
- Usage: `bkpdir list`
- Shows each archive with its path (relative to current directory) and creation time using configurable format string or template
- Default format: `.bkpdir/project-2024-03-21-15-30=main=abc123=note.zip (created: 2024-03-21 15:30:00)`
- Output formatting uses `format_list_archive` configuration setting with printf-style specifications
- Alternative template formatting uses `template_list_archive` with named placeholders and `pattern_archive_filename` for data extraction
- Supports text highlighting and color formatting through ANSI escape codes in format strings and templates
- Template-based formatting allows rich data extraction from archive filenames using named regex groups
- Archives are sorted by creation time (most recent first)
- Shows verification status if available: [VERIFIED], [FAILED], or [UNVERIFIED]
- Handles errors gracefully with appropriate status codes using `format_error` or `template_error` configuration

### 4. Verify Archive
- Verifies archive integrity and optionally checksums
- Usage: `bkpdir verify [ARCHIVE_NAME]`
- Flags:
  - `--checksum`: Include checksum verification of archive contents
- Performs ZIP archive structure and integrity verification
- With --checksum flag: verifies file contents against stored checksums
- Stores verification results for display in list command
- Reports verification status using configurable format strings
- Uses appropriate status codes for verification results

### 5. Create File Backup
- Creates a backup of a single file with robust error handling and resource cleanup
- Usage: `bkpdir backup [FILE_PATH] [NOTE]`
- Before creating a backup:
  - Compares the file with its most recent backup using byte comparison
  - If the file is identical to the most recent backup:
    - Reports the existing backup path using `format_identical_backup` or `template_identical_backup` configuration
    - Template formatting can extract and display rich information from backup filename using `pattern_backup_filename`
    - Default format: "File is identical to existing backup: [PATH]"
    - Template format can show: "File [filename] is identical to backup from [year]-[month]-[day]: [path]"
    - Exits with `status_file_is_identical_to_existing_backup` status code
- When a new backup is created:
  - Reports success using `format_created_backup` or `template_created_backup` configuration
  - Default format: "Created backup: [PATH]"
  - Template format can include extracted filename and timestamp information
  - Exits with `status_created_backup` status code
- Backup naming format: `SOURCE_FILENAME-YYYY-MM-DD-hh-mm[=NOTE]`
  - SOURCE_FILENAME is the base name of the original file
  - YYYY-MM-DD-hh-mm is the timestamp of the backup
  - NOTE is an optional note appended with an equals sign
- The backup maintains the original file's directory structure in the backup path
- NOTE is an optional positional argument provided by the user
- All output uses configurable printf-style format strings or template-based formatting for consistency and customization
- Error messages use `format_error` or `template_error` configuration setting with optional operation context
- Enhanced backup features:
  - **Atomic Operations**: Uses temporary files to ensure backup integrity
  - **Resource Cleanup**: Automatically cleans up temporary files on success or failure
  - **Context Support**: Supports operation cancellation and timeouts
  - **Enhanced Error Detection**: Detects various disk space and permission conditions
  - **Panic Recovery**: Recovers from unexpected errors without leaving temporary files
  - **Thread Safety**: Safe for concurrent operations
  - **Template Formatting**: Rich data extraction from backup filenames for enhanced display
  - **Operation Context**: Error messages include operation context for better debugging
  - **Enhanced File Operations**: Complete file system operations with comprehensive error handling

### 6. List File Backups
- Displays all backups associated with the specified file
- Usage: `bkpdir --list [FILE_PATH]`
- Shows each backup with its path (relative to current directory) and creation time using configurable format string or template
- Default format: `.bkpdir/path/to/file.txt-2024-03-21-15-30=note (created: 2024-03-21 15:30:00)`
- Output formatting uses `format_list_backup` configuration setting with printf-style specifications
- Alternative template formatting uses `template_list_backup` with named placeholders and `pattern_backup_filename` for data extraction
- Supports text highlighting and color formatting through ANSI escape codes in format strings and templates
- Template-based formatting allows rich data extraction from backup filenames using named regex groups
- Backups are sorted by creation time (most recent first)
- Backups are organized by their source file paths
- Handles errors gracefully with appropriate status codes using `format_error` or `template_error` configuration

### 7. Display Configuration
- Displays computed configuration values after processing configuration files
- Usage: `bkpdir --config`
- Shows each configuration value with its name, computed value (including defaults), and source file
- Output formatting uses `format_config_value` configuration setting with printf-style specifications
- Alternative template formatting uses `template_config_value` with named placeholders for enhanced display
- Default format: `archive_dir_path: ../.bkpdir (source: default)`
- Supports text highlighting and color formatting for enhanced readability
- Template-based formatting allows conditional formatting based on configuration source and value types
- The application exits after displaying the configuration values
- Configuration files are processed from the `BKPDIR_CONFIG` environment variable path list
- If `BKPDIR_CONFIG` is not set, uses the default search path
- Includes display of all format string configurations, template configurations, and regex patterns with their current values and sources
- Shows both directory archiving and file backup configuration options

## Global Options
- **Dry-Run Mode**: When enabled with `--dry-run` flag:
  - For directory operations: Shows the archive filename that would be created using `format_dry_run_archive` or `template_dry_run_archive` configuration
    - Default format: "Would create archive: [PATH]"
    - Template format can show: "Would create archive for [prefix] on [year]-[month]-[day]: [path]"
    - Template-based formatting can extract and display rich information about the planned archive
  - For file operations: Shows the backup filename that would be created using `format_dry_run_backup` or `template_dry_run_backup` configuration
    - Default format: "Would create backup: [PATH]"
    - Template format can show: "Would create backup for [filename] on [year]-[month]-[day]: [path]"
    - Template-based formatting can extract and display rich information about the planned backup
  - No actual archive or backup is created
  - Includes resource cleanup verification in dry-run mode
  - All output uses configurable printf-style format strings or template-based formatting

## Archive Features

### Git Integration
- Automatically detects Git repositories using `git rev-parse --is-inside-work-tree`
- Extracts current branch name using `git rev-parse --abbrev-ref HEAD`
- Extracts commit hash using `git rev-parse --short HEAD`
- Includes Git information in archive names when `include_git_info` is enabled
- Gracefully handles non-Git directories by returning empty Git info
- Configurable via `include_git_info` setting in configuration
- Works with both clean and dirty working directories

### File Exclusion
- Supports glob patterns for excluding files and directories
- Default exclusions: `.git/`, `vendor/`
- Uses doublestar glob matching for advanced patterns
- Configurable via `exclude_patterns` setting
- Patterns are applied recursively to directory tree

### Archive Verification
- ZIP structure integrity checking
- Optional SHA-256 checksum verification
- Verification status stored and displayed in listings
- Configurable automatic verification after creation
- Supports manual verification of existing archives

### Incremental Archives
- Creates archives containing only changed files
- Based on modification time comparison with base archive
- Maintains relationship to base archive in naming
- Efficient for large directories with few changes
- Supports same Git integration and exclusion patterns

## Error Handling and Recovery

### Structured Error Reporting
- All operations return structured errors with specific status codes
- Human-readable error messages for common scenarios
- Technical details logged to stderr when appropriate
- No sensitive information exposed in error messages

### Enhanced Error Detection
- **Disk Space**: Detects various disk full conditions including:
  - "no space left on device"
  - "disk full"
  - "not enough space"
  - "insufficient disk space"
  - "device full"
  - "quota exceeded"
  - "file too large"
- **Permission Errors**: Proper handling of directory access permissions
- **Directory Type Validation**: Ensures only directories are archived
- **Path Resolution**: Handles both absolute and relative paths securely

### Panic Recovery
- Critical operations include panic recovery mechanisms
- Panics are logged to stderr without exposing internal details
- Resource cleanup still occurs even when panics happen
- Application doesn't crash on unexpected errors

### Context and Cancellation Support
- Long-running operations support cancellation via context
- Timeout handling for operations that might hang
- Graceful shutdown with proper resource cleanup
- Context cancellation checked at multiple operation points

## Resource Management

### Automatic Cleanup
- All temporary files and directories are automatically cleaned up
- Cleanup occurs on success, failure, and cancellation
- Thread-safe resource tracking for concurrent operations
- Error-resilient cleanup continues even if individual operations fail

### Atomic Operations
- Archive operations use temporary files to ensure atomicity
- Atomic rename operations prevent data corruption
- Temporary files are registered for automatic cleanup
- Successful operations remove files from cleanup lists

### Leak Prevention
- Comprehensive testing verifies no resource leaks
- All temporary resources are tracked and cleaned
- No orphaned files remain in any error scenario
- Memory usage is properly managed

## Build and Development Requirements

### Code Quality Standards
- All code must pass `revive` linter before commit
- Comprehensive test coverage required for all features
- All errors must be properly handled
- Documentation required for all public functions
- Backward compatibility must be maintained

### Build System
- `make lint`: Run code linting
- `make test`: Run all tests with verbose output
- `make build`: Build application (depends on lint and test passing)
- `make clean`: Remove build artifacts

### Testing Requirements
- Unit tests for all core functions
- Integration tests for complete workflows
- Resource cleanup verification in all test scenarios
- Context cancellation and timeout testing
- Performance benchmarks for critical operations
- Stress testing for concurrent operations

## Testing Infrastructure

### Archive Corruption Testing Framework (TEST-INFRA-001-A) ✅ COMPLETED

The application includes comprehensive testing infrastructure for systematic archive corruption testing to validate verification logic under various failure scenarios.

#### Controlled ZIP Corruption Utilities
- **Systematic Corruption**: Provides 8 distinct corruption types targeting different ZIP format structures
  - **CRC Corruption**: Corrupt file checksums while maintaining readable archive structure
  - **Header Corruption**: Corrupt ZIP file headers to render archives unreadable
  - **Truncation**: Cut off end of archive to simulate incomplete transfers
  - **Central Directory Corruption**: Corrupt ZIP central directory (usually fatal)
  - **Local Header Corruption**: Corrupt individual file headers (sometimes recoverable)
  - **Data Corruption**: Corrupt actual file data within archives
  - **Signature Corruption**: Corrupt ZIP file signatures for immediate failure detection
  - **Comment Corruption**: Corrupt archive comments (non-fatal corruption type)

- **Safe Testing**: Automatic backup/restore functionality prevents data loss during testing
  - Creates backup before applying any corruption
  - Automatic restoration on test failure or completion
  - Panic recovery ensures cleanup occurs even during unexpected failures
  - Thread-safe resource tracking for concurrent testing

#### Deterministic Corruption Patterns
- **Reproducible Testing**: Seed-based corruption generation ensures consistent test results
  - Identical archives with identical seeds produce identical corruption
  - Offset-based variation generates different but deterministic corruption per location
  - Essential for CI/CD testing where reproducible results are required
  - Enables regression testing for verification logic improvements

- **Configurable Corruption**: Precise control over corruption parameters
  - **Type**: Specific corruption type from enumerated set
  - **Seed**: Reproducibility seed for deterministic corruption
  - **Size**: Number of bytes to corrupt
  - **Offset**: Specific byte offset or random placement
  - **Severity**: Corruption intensity from 0.0 (minimal) to 1.0 (maximum)

#### Archive Repair Detection
- **Automatic Classification**: Detects and classifies different types of archive corruption
  - Analyzes ZIP signatures, headers, and structure for corruption indicators
  - Classifies corruption as recoverable vs fatal for appropriate test expectations
  - Supports multiple simultaneous corruption types in single archive
  - Returns detailed corruption analysis for verification testing

- **Recovery Testing**: Framework supports testing archive repair mechanisms
  - Tracks original bytes for potential recovery scenarios
  - Tests verification behavior with different corruption severities
  - Validates detection logic matches applied corruption types
  - Enables testing of graceful degradation and recovery paths

#### Performance Characteristics
- **Corruption Speed**: Optimized for fast test execution
  - CRC corruption: ~763μs average for typical archives
  - Header corruption: ~45μs average (minimal data modification)
  - Truncation: ~12μs average (simple file size operation)
  - Detection analysis: ~49μs average for classification

- **Resource Efficiency**: Minimal overhead during testing
  - <1KB additional memory for corruption data
  - 1:1 archive size ratio for backup storage during testing
  - <1ms cleanup overhead with automatic resource management
  - Cross-platform compatibility with proper file permission handling

#### Integration with Verification Logic
- **Testing Utilities**: Designed for use with existing verification systems
  - Integrates with `verify.go` for testing archive verification logic
  - Compatible with `comparison.go` for testing archive comparison under corruption
  - Uses existing `ResourceManager` for cleanup and error handling
  - Supports existing `OutputFormatter` for test result display

- **Test Coverage**: Enables comprehensive testing of previously untestable scenarios
  - Error paths that are difficult to reproduce in normal testing
  - Edge cases with various corruption combinations
  - Performance baselines for verification operations
  - Cross-platform behavior validation

#### Example Usage

```go
// Create systematic corruption for testing
config := CorruptionConfig{
    Type:           CorruptionCRC,
    Seed:           12345, // For reproducible tests
    CorruptionSize: 4,     // Corrupt 4 bytes
    Severity:       0.5,   // Moderate corruption
}

// Apply controlled corruption with automatic backup
result, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
if err != nil {
    t.Fatalf("Failed to create corrupted archive: %v", err)
}

// Test verification behavior with corrupted archive
detector := NewCorruptionDetector(archivePath)
detected, err := detector.DetectCorruption()
if err != nil {
    t.Fatalf("Corruption detection failed: %v", err)
}

// Validate detection matches applied corruption
if !contains(detected, CorruptionCRC) {
    t.Errorf("CRC corruption not detected: got %v", detected)
}
```

This testing infrastructure provides the foundation for comprehensive verification testing, enabling systematic validation of archive verification logic against real-world corruption scenarios that would otherwise be difficult to reproduce reliably.

## Implementation Details
For detailed implementation requirements and constraints, see:
- [Immutable Specifications](immutable.md) for core behaviors that cannot be changed
- [Architecture](architecture.md) for system design and implementation details
- [Requirements](requirements.md) for technical requirements and test coverage
- [Testing](testing.md) for comprehensive testing requirements

## Platform Compatibility
- Works on macOS and Linux systems
- Uses platform-independent path handling
- Preserves file permissions and ownership where applicable
- Handles file system differences between platforms
- Thread-safe operations for concurrent access
- Efficient resource management across platforms

## Performance Characteristics
- Minimal overhead for resource tracking
- Efficient directory comparison with early termination
- Optimized cleanup operations
- Low memory footprint for large directories
- Fast atomic archive operations
- Scalable for large directories and many archives
- Streaming ZIP creation for memory efficiency 