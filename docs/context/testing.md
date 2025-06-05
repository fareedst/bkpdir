# Testing Requirements and Architecture

This document outlines the testing requirements, architecture, and approach for the BkpDir application.

## Test Categories

### Unit Tests
**Implementation**: `*_test.go` files
**Requirements**:
- `TestGenerateArchiveName`: Tests archive naming
  - Validates naming format with various inputs including Git information
  - Test cases:
    - Basic archive naming without Git info
    - Archive naming with Git branch and hash
    - Archive naming with notes
    - Archive naming with prefix
    - Incremental archive naming
    - Various timestamp formats
- `TestGenerateBackupName`: Tests file backup naming
  - Validates backup filename generation with timestamps and notes
  - Test cases:
    - Basic backup naming without notes
    - Backup naming with notes
    - Various source file paths and extensions
    - Timestamp format validation
- `TestDefaultConfig` and `TestLoadConfig`: Tests configuration
  - Validates configuration loading and defaults
  - Tests configuration discovery with `BKPDIR_CONFIG` environment variable
  - Tests multiple configuration file precedence
  - Tests home directory expansion in paths
  - Tests status code configuration loading and defaults
  - Validates all status code fields are properly loaded from YAML
  - Tests status code configuration with custom values
  - Tests exclude patterns configuration
  - Tests verification configuration
  - Tests Git integration configuration
  - Tests format string configuration (printf-style and template-based)
  - Tests regex pattern configuration
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
  - Tests display of format string configuration values
  - Tests display of template configuration values
  - Tests display of regex pattern configuration values
  - Test cases:
    - Configuration with default values only
    - Configuration from single file
    - Configuration from multiple files with precedence
    - Configuration with missing files
    - Configuration with invalid files
    - Environment variable override scenarios
    - Status code configuration display
    - Exclude patterns configuration display
    - Verification configuration display
    - Format string configuration display
    - Template configuration display
    - Regex pattern configuration display
- `TestConfigModification`: Tests configuration value modification
  - Validates configuration setting and persistence
  - Tests type-safe value conversion and validation
  - Tests YAML file creation and updates
  - Tests nested configuration structure handling
  - Test cases:
    - Setting string configuration values (archive_dir_path, backup_dir_path, checksum_algorithm)
    - Setting boolean configuration values (use_current_dir_name, include_git_info, verify_on_create)
    - Setting integer configuration values (status codes)
    - Creating new configuration file when none exists
    - Updating existing configuration file
    - Preserving existing configuration values when setting new ones
    - Handling nested verification section (verify_on_create, checksum_algorithm)
    - Error handling for invalid configuration keys
    - Error handling for invalid boolean values
    - Error handling for invalid integer values
    - Error handling for file I/O errors
    - Error handling for YAML parsing errors
    - Configuration persistence verification
    - Type validation for all supported configuration types
- `TestShouldExcludeFile`: Tests file exclusion patterns
  - Validates doublestar glob pattern matching
  - Test cases:
    - Basic glob patterns
    - Directory exclusions
    - File extension exclusions
    - Complex nested patterns
    - Case sensitivity handling
- `TestCopyFileWithContext`: Tests context-aware file copying
  - Validates file copying with cancellation support
  - Test cases:
    - Successful copy with context
    - Copy cancelled via context
    - Copy with timeout
    - Copy with already cancelled context
    - Large file copying with periodic cancellation checks
- `TestCopyFile`: Tests standard file copying
  - Validates file copying without context
  - Test cases:
    - Successful file copy
    - Permission preservation
    - Error handling for missing source
    - Error handling for permission denied
- `TestCompareDirectories`: Tests directory comparison
  - Validates directory tree comparison functionality
  - Test cases:
    - Identical directories
    - Different directories
    - Directories with different file sizes
    - Empty directories
    - Large directories
    - Directories with special characters
    - Directories with excluded files
- `TestCompareFiles`: Tests file comparison
  - Validates byte-by-byte file comparison
  - Test cases:
    - Identical files
    - Different files
    - Files with different sizes
    - Empty files
    - Large files
    - Binary files
    - Files with special characters
- `TestListArchives`: Tests archive listing
  - Validates archive listing and sorting
  - Tests verification status display
  - Test cases:
    - Empty archive directory
    - Multiple archives with different timestamps
    - Archives with Git information
    - Archives with notes
    - Archives with verification status
    - Incremental archives
- `TestListFileBackups`: Tests file backup listing
  - Validates backup listing and sorting for specific files
  - Test cases:
    - Empty backup directory
    - Multiple backups with different timestamps
    - Backups with notes
    - Backups in nested directory structures
    - Sorting by creation time
- `TestCreateFullArchive`: Tests full archive creation
  - Validates full archive creation with various scenarios
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
- `TestCreateFullArchiveWithCleanup`: Tests archive creation with resource cleanup
  - Validates automatic resource cleanup functionality
  - Tests atomic operations with temporary files
  - Test cases:
    - Successful archive with cleanup verification
    - Archive failure with cleanup verification
    - No temporary files left after operations
    - Atomic file operations
- `TestCreateFullArchiveWithContext`: Tests context-aware archive creation
  - Validates archive creation with cancellation support
  - Test cases:
    - Successful archive with context
    - Archive cancelled via context
    - Archive with timeout
    - Context cancellation at various stages
    - Cancellation during file scanning
    - Cancellation during ZIP creation
- `TestCreateFullArchiveWithContextAndCleanup`: Tests context-aware archive with cleanup
  - Validates most robust archive creation functionality
  - Combines context support with resource cleanup
  - Test cases:
    - Successful archive with context and cleanup
    - Cancelled archive with proper cleanup
    - Timeout scenarios with cleanup verification
    - No resource leaks on cancellation
- `TestCreateFileBackup`: Tests file backup creation
  - Validates file backup creation with various scenarios
  - Tests status code exit behavior for different conditions
  - Test cases:
    - Successful backup creation (should exit with `cfg.StatusCreatedBackup`)
    - File identical to existing backup (should exit with `cfg.StatusFileIsIdenticalToExistingBackup`)
    - File not found (should exit with `cfg.StatusFileNotFound`)
    - Invalid file type (should exit with `cfg.StatusInvalidFileType`)
    - Permission denied (should exit with `cfg.StatusPermissionDenied`)
    - Disk full scenarios (should exit with `cfg.StatusDiskFull`)
    - Backup directory creation failure (should exit with `cfg.StatusFailedToCreateBackupDirectory`)
    - Configuration errors (should exit with `cfg.StatusConfigError`)
    - Panic recovery scenarios
- `TestCreateFileBackupWithContext`: Tests context-aware file backup creation
  - Validates file backup creation with cancellation support
  - Test cases:
    - Successful backup with context
    - Backup cancelled via context
    - Backup with timeout
    - Context cancellation during file copying
    - Context cancellation during file comparison
    - Already cancelled context handling
- `TestCreateFileBackupWithContextAndCleanup`: Tests context-aware backup with cleanup
  - Validates most robust backup creation functionality
  - Combines context support with resource cleanup
  - Test cases:
    - Successful backup with context and cleanup
    - Cancelled backup with proper cleanup
    - Timeout scenarios with cleanup verification
    - No resource leaks on cancellation
    - Atomic operations with temporary files
    - Cleanup verification after panic recovery
- `TestCheckForIdenticalFileBackup`: Tests file backup comparison
  - Validates file comparison with existing backups
  - Test cases:
    - File identical to most recent backup
    - File different from existing backups
    - No existing backups
    - Multiple backups with different content
    - Backup comparison with various file sizes
- `TestTemplateFormatter`: Tests template-based formatting
  - Validates template formatting with regex extraction
  - Tests both Go text/template and placeholder syntax
  - Test cases:
    - `FormatWithTemplate()` with valid regex patterns
    - `FormatWithTemplate()` with invalid regex patterns
    - `FormatWithPlaceholders()` with valid data
    - `FormatWithPlaceholders()` with missing data
    - Template formatting for archive operations
    - Template formatting for backup operations
    - Template formatting for configuration display
    - Template formatting for error messages with operation context
    - ANSI color code handling in templates
    - Fallback to placeholder formatting on template errors
    - Named regex group extraction for archives and backups
- `TestTemplateFormatterArchiveOperations`: Tests archive template formatting
  - Validates template formatting for directory operations
  - Test cases:
    - `TemplateCreatedArchive()` with archive filename extraction
    - `TemplateIdenticalArchive()` with archive metadata
    - `TemplateListArchive()` with creation time formatting
    - `TemplateDryRunArchive()` with planned archive information
    - Archive filename parsing with Git information
    - Archive filename parsing with notes
    - Archive filename parsing with branch and hash
    - Template error handling and fallback
- `TestTemplateFormatterBackupOperations`: Tests backup template formatting
  - Validates template formatting for file operations
  - Test cases:
    - `TemplateCreatedBackup()` with backup filename extraction
    - `TemplateIdenticalBackup()` with backup metadata
    - `TemplateListBackup()` with creation time formatting
    - `TemplateDryRunBackup()` with planned backup information
    - Backup filename parsing with notes
    - Backup filename parsing with timestamps
    - Template error handling and fallback
- `TestTemplateFormatterCommonOperations`: Tests common template formatting
  - Validates template formatting for configuration and errors
  - Test cases:
    - `TemplateConfigValue()` with conditional formatting
    - `TemplateError()` with operation context
    - Configuration source-based formatting
    - Error context integration
    - Template validation and error handling
- `TestOutputFormatter`: Tests printf-style formatting
  - Validates centralized output formatting
  - Test cases:
    - Format methods for all directory operations
    - Format methods for all file operations
    - Format methods for configuration and errors
    - ANSI color code handling
    - Format string validation
    - Print methods for stdout/stderr routing
    - Integration with template formatting
- `TestOutputFormatterArchiveOperations`: Tests archive output formatting
  - Validates printf-style formatting for directory operations
  - Test cases:
    - `FormatCreatedArchive()` with various paths
    - `FormatIdenticalArchive()` with archive paths
    - `FormatListArchive()` with timestamps
    - `FormatDryRunArchive()` with planned paths
    - Color formatting and ANSI codes
    - Format string customization
- `TestOutputFormatterBackupOperations`: Tests backup output formatting
  - Validates printf-style formatting for file operations
  - Test cases:
    - `FormatCreatedBackup()` with various paths
    - `FormatIdenticalBackup()` with backup paths
    - `FormatListBackup()` with timestamps
    - `FormatDryRunBackup()` with planned paths
    - Color formatting and ANSI codes
    - Format string customization
- `TestEnhancedErrorHandling`: Tests enhanced error handling
  - Validates structured error handling with operation context
  - Test cases:
    - `ArchiveError` creation with operation context
    - `BackupError` creation with operation context
    - Enhanced disk space error detection
    - Error context preservation
    - Template error formatting with operation context
    - Panic recovery with error context
        - Status code extraction from structured errors
    - File different from existing backups
    - No existing backups
    - Multiple backups with different content
    - Large file comparison
    - Binary file comparison
    - Files with special characters
    - Permission denied during comparison
- `TestCreateIncrementalArchive`: Tests incremental archive creation
  - Validates incremental archive functionality
  - Test cases:
    - Successful incremental archive creation
    - No base archive available
    - No changes since base archive
    - Multiple incremental archives
    - Incremental archive with Git changes
- `TestResourceManager`: Tests resource management functionality
  - Validates thread-safe resource tracking
  - Tests automatic cleanup mechanisms
  - Test cases:
    - Basic resource registration and cleanup
    - Thread-safe concurrent access
    - Cleanup of both files and directories
    - Error-resilient cleanup (continues on individual failures)
    - Cleanup warnings logged to stderr
    - Panic recovery during cleanup
- `TestArchiveError`: Tests structured error handling
  - Validates ArchiveError functionality
  - Test cases:
    - Error creation with message and status code
    - Error interface implementation
    - Status code extraction
    - Error message formatting
    - Error context preservation
    - Error unwrapping
- `TestIsDiskFullError`: Tests enhanced disk space detection
  - Validates disk full error detection
  - Test cases:
    - Various disk full error messages
    - Case-insensitive matching
    - Multiple disk space indicators
    - Non-disk-full errors (should return false)
    - Nil error handling
- `TestIsPermissionError`: Tests permission error detection
  - Validates permission error detection
  - Test cases:
    - Various permission error messages
    - Case-insensitive matching
    - Multiple permission indicators
    - Non-permission errors (should return false)
- `TestValidateDirectoryPath`: Tests directory validation
  - Validates directory path validation functionality
  - Test cases:
    - Valid directory paths
    - Non-existent directories
    - Files instead of directories
    - Permission denied scenarios
- `TestValidateFilePath`: Tests file validation
  - Validates file path validation functionality
  - Test cases:
    - Valid file paths
    - Non-existent files
    - Directories instead of files
    - Permission denied scenarios
- `TestConfigurationDiscovery`: Tests configuration file discovery
  - Tests multiple configuration files with different precedence
  - Tests environment variable override behavior
  - Tests hard-coded default path behavior
  - Tests missing configuration files handling
  - Tests invalid configuration file handling
  - Tests configuration merging with defaults
  - Tests status code configuration precedence and merging
  - Tests format string configuration precedence and merging
  - Tests template configuration precedence and merging
  - Tests regex pattern configuration precedence and merging
- `TestStatusCodeConfiguration`: Tests status code configuration
  - Validates status code loading from YAML
  - Tests status code defaults
  - Tests status code precedence with multiple configuration files
  - Tests invalid status code values handling
  - Test cases:
    - Default status codes (should match immutable specification defaults)
    - Custom status codes from configuration file
    - Status code precedence with multiple files
    - Invalid status code values (non-integer)
    - Missing status code fields (should use defaults)

### Git Integration Tests
**Implementation**: `git_test.go`
**Requirements**:
- `TestIsGitRepository`: Tests Git repository detection
  - Test cases:
    - Valid Git repository
    - Non-Git directory
    - Directory with .git file (submodule)
    - Permission denied scenarios
- `TestGetGitBranch`: Tests Git branch extraction
  - Test cases:
    - Standard branch names
    - Branch names with special characters
    - Detached HEAD state
    - Non-Git repository
- `TestGetGitShortHash`: Tests Git commit hash extraction
  - Test cases:
    - Valid commit hash
    - Initial commit
    - Dirty working directory
    - Non-Git repository
- `TestGitIntegrationInArchiveNaming`: Tests Git info in archive names
  - Test cases:
    - Archive naming with Git info enabled
    - Archive naming with Git info disabled
    - Archive naming in non-Git directory
    - Archive naming with special branch names
- `TestGitIntegrationInBackupNaming`: Tests Git info in backup operations
  - Test cases:
    - Backup creation in Git repository
    - Backup creation in non-Git directory
    - Git info handling in backup workflows

### Archive Verification Tests
**Implementation**: `verify_test.go`
**Requirements**:
- `TestVerifyArchive`: Tests archive structure verification
  - Validates ZIP file structure verification
  - Test cases:
    - Valid ZIP archive
    - Corrupted ZIP archive
    - Non-ZIP file
    - Empty file
    - Archive with special characters
- `TestVerifyChecksums`: Tests checksum verification
  - Validates file content verification against stored checksums
  - Test cases:
    - Valid checksums
    - Corrupted file content
    - Missing checksum file
    - Invalid checksum format
    - Checksum algorithm mismatch
- `TestGenerateChecksums`: Tests checksum generation
  - Validates checksum generation for files
  - Test cases:
    - Single file checksum
    - Multiple file checksums
    - Large file checksums
    - Empty file checksums
    - Files with special characters
    - Directory exclusion from checksums
- `TestStoreAndReadChecksums`: Tests checksum storage and retrieval
  - Validates checksum persistence in archives
  - Test cases:
    - Store checksums in archive
    - Read checksums from archive
    - Invalid checksum data
    - Missing checksum file
- `TestVerificationStatus`: Tests verification status tracking
  - Validates verification result storage and retrieval
  - Test cases:
    - Successful verification status
    - Failed verification status
    - Verification status persistence
    - Multiple verification attempts

### Output Formatting Tests
**Implementation**: `formatter_test.go`
**Requirements**:
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
  - Test cases:
    - Default format string
    - Custom format string with colors
    - Format string with special characters
    - Various path lengths and characters
- `TestFormatCreatedBackup`: Tests backup creation message formatting
  - Validates format string application for successful backup messages
  - Test cases:
    - Default format string
    - Custom format string with colors
    - Format string with special characters
    - Various path lengths and characters
- `TestFormatIdenticalArchive`: Tests identical directory message formatting
  - Validates format string application for identical directory messages
  - Test cases:
    - Default format string
    - Custom format string with highlighting
    - Format string with symbols and colors
    - Various path formats
- `TestFormatIdenticalBackup`: Tests identical file message formatting
  - Validates format string application for identical file messages
  - Test cases:
    - Default format string
    - Custom format string with highlighting
    - Format string with symbols and colors
    - Various path formats
- `TestFormatListArchive`: Tests archive listing entry formatting
  - Validates format string application for archive list entries
  - Test cases:
    - Default format string with two parameters
    - Custom format string with color coding
    - Format string with timestamp formatting
    - Various path and time combinations
- `TestFormatListBackup`: Tests backup listing entry formatting
  - Validates format string application for backup list entries
  - Test cases:
    - Default format string with two parameters
    - Custom format string with color coding
    - Format string with timestamp formatting
    - Various path and time combinations
- `TestFormatConfigValue`: Tests configuration value display formatting
  - Validates format string application for configuration display
  - Test cases:
    - Default format string with three parameters
    - Custom format string with highlighting
    - Format string with source emphasis
    - Various configuration name/value combinations
- `TestFormatDryRunArchive`: Tests dry-run archive message formatting
  - Validates format string application for dry-run archive messages
  - Test cases:
    - Default format string
    - Custom format string with warning indicators
    - Format string with visual emphasis
    - Various path formats
- `TestFormatDryRunBackup`: Tests dry-run backup message formatting
  - Validates format string application for dry-run backup messages
  - Test cases:
    - Default format string
    - Custom format string with warning indicators
    - Format string with visual emphasis
    - Various path formats
- `TestFormatError`: Tests error message formatting
  - Validates format string application for error messages
  - Test cases:
    - Default format string
    - Custom format string with error highlighting
    - Format string with symbols and colors
    - Various error message types
- `TestPrintFunctions`: Tests direct printing functions
  - Validates output to stdout and stderr
  - Test cases:
    - PrintCreatedArchive outputs to stdout
    - PrintCreatedBackup outputs to stdout
    - PrintIdenticalArchive outputs to stdout
    - PrintIdenticalBackup outputs to stdout
    - PrintListArchive outputs to stdout
    - PrintListBackup outputs to stdout
    - PrintConfigValue outputs to stdout
    - PrintDryRunArchive outputs to stdout
    - PrintDryRunBackup outputs to stdout
    - PrintError outputs to stderr
    - Output matches formatted strings
- `TestFormatStringConfiguration`: Tests format string configuration loading
  - Validates format string loading from YAML
  - Tests format string defaults
  - Tests format string precedence with multiple configuration files
  - Test cases:
    - Default format strings (should match immutable specification defaults)
    - Custom format strings from configuration file
    - Format string precedence with multiple files
    - Invalid format strings (should use safe defaults)
    - Missing format string fields (should use defaults)
    - Format strings with ANSI escape codes
    - Format strings with special characters
- `TestFormatStringValidation`: Tests format string validation
  - Validates printf-style format string compatibility
  - Test cases:
    - Valid format strings with correct parameter counts
    - Invalid format strings with wrong parameter counts
    - Format strings with unsupported format specifiers
    - Format strings with malformed syntax
    - Empty or nil format strings
    - Format strings with escape sequences
- `TestTemplateFormatter`: Tests template formatter functionality
  - Validates template-based formatting with named placeholders
  - Tests regex pattern extraction and template application
  - Tests both Go text/template syntax and %{name} placeholder syntax
  - Test cases:
    - Basic template string application with {{.name}} syntax
    - Placeholder formatting with %{name} syntax
    - Template strings with ANSI color codes
    - Template strings with conditional logic
    - Invalid template strings (should fall back to safe defaults)
    - Empty template strings
    - Template strings with multiple parameters
    - Regex pattern integration with template formatting
- `TestTemplateCreatedArchive`: Tests template-based archive creation formatting
  - Validates template formatting for archive creation messages
  - Test cases:
    - Template with named placeholders
    - Template with regex pattern extraction
    - Template with conditional formatting
    - Template with archive metadata
    - Default template string with path extraction
    - Custom template string with prefix and timestamp extraction
    - Template string with ANSI color codes
    - Archive paths with notes
    - Archive paths without notes
    - Invalid archive path formats (should fall back gracefully)
- `TestTemplateCreatedBackup`: Tests template-based backup creation formatting
  - Validates template formatting for backup creation messages
  - Test cases:
    - Template with named placeholders
    - Template with regex pattern extraction
    - Template with conditional formatting
    - Template with backup metadata
    - Default template string with path extraction
    - Custom template string with filename and timestamp extraction
    - Template string with ANSI color codes
    - Backup paths with notes
    - Backup paths without notes
    - Invalid backup path formats (should fall back gracefully)
- `TestTemplateIdenticalArchive`: Tests template-based identical directory formatting
  - Validates template formatting for identical directory messages
  - Test cases:
    - Template with archive filename parsing
    - Template with Git information extraction
    - Template with conditional formatting based on extracted data
    - Archive paths with and without notes
    - Various archive filename formats
- `TestTemplateIdenticalBackup`: Tests template-based identical file formatting
  - Validates template formatting for identical file messages
  - Test cases:
    - Default template string with path extraction
    - Custom template string with rich backup information
    - Template string with conditional formatting based on extracted data
    - Backup paths with and without notes
    - Various backup filename formats
- `TestTemplateListArchive`: Tests template-based archive listing formatting
  - Validates template formatting for archive list entries
  - Test cases:
    - Template with archive filename parsing
    - Template with Git information extraction
    - Template with timestamp formatting
    - Template with note display
    - Default template string with path and time parameters
    - Custom template string with extracted prefix and timestamp data
    - Template string with note extraction and display
    - Various timestamp formats
    - Archive paths with special characters
    - Template formatting with missing data (graceful degradation)
- `TestTemplateListBackup`: Tests template-based backup listing formatting
  - Validates template formatting for backup list entries
  - Test cases:
    - Template with backup filename parsing
    - Template with timestamp formatting
    - Template with note display
    - Template with file metadata
    - Default template string with path and time parameters
    - Custom template string with extracted filename and timestamp data
    - Template string with note extraction and display
    - Various timestamp formats
    - Backup paths with special characters
    - Template formatting with missing data (graceful degradation)
- `TestTemplateConfigValue`: Tests template-based configuration display formatting
  - Validates template formatting for configuration display
  - Test cases:
    - Default template string with three parameters
    - Custom template string with conditional formatting based on source
    - Template string with value type detection
    - Various configuration name/value combinations
    - Source highlighting and emphasis
- `TestTemplateDryRunArchive`: Tests template-based dry-run archive formatting
  - Validates template formatting for dry-run archive messages
  - Test cases:
    - Default template string with path extraction
    - Custom template string with planned archive information
    - Template string with prefix and timestamp extraction
    - Various planned archive path formats
- `TestTemplateDryRunBackup`: Tests template-based dry-run backup formatting
  - Validates template formatting for dry-run backup messages
  - Test cases:
    - Default template string with path extraction
    - Custom template string with planned backup information
    - Template string with filename and timestamp extraction
    - Various planned backup path formats
- `TestTemplateError`: Tests template-based error message formatting
  - Validates template formatting for error messages
  - Test cases:
    - Default template string with message parameter
    - Custom template string with operation context
    - Template string with error categorization
    - Various error message types
    - Operation context integration
- `TestFormatWithTemplate`: Tests Go text/template integration
  - Validates template application with named regex groups
  - Test cases:
    - Archive filename pattern extraction
    - Backup filename pattern extraction
    - Configuration line pattern extraction
    - Timestamp pattern extraction
    - Invalid patterns and templates
    - Valid regex patterns with named groups
    - Template strings with {{.name}} syntax
    - Complex template logic with conditionals
    - Invalid regex patterns (should return error)
    - No regex matches (should return error)
    - Invalid template syntax (should return error)
    - Template execution errors
- `TestFormatWithPlaceholders`: Tests placeholder replacement
  - Validates %{name} placeholder substitution
  - Test cases:
    - Basic placeholder replacement
    - Multiple placeholders
    - Missing placeholders
    - Special characters in placeholders
    - Multiple placeholders in single string
    - Unmatched placeholders (should be left intact)
    - Empty placeholder names
    - Placeholders with special characters
    - ANSI color codes in placeholder values
    - Nested placeholder-like strings
- `TestExtractArchiveFilenameData`: Tests archive filename data extraction
  - Validates regex pattern extraction from archive filenames
  - Test cases:
    - Basic archive filename parsing
    - Archive filename with Git information
    - Archive filename with notes
    - Invalid archive filenames
- `TestExtractBackupFilenameData`: Tests backup filename data extraction
  - Validates regex pattern extraction from backup filenames
  - Test cases:
    - Basic backup filename parsing
    - Backup filename with notes
    - Invalid backup filenames
    - Various file extensions
- `TestTemplateStringConfiguration`: Tests template string configuration loading
  - Validates template string loading from YAML
  - Tests template string defaults
  - Tests template string precedence with multiple configuration files
  - Test cases:
    - Default template strings (should match immutable specification defaults)
    - Custom template strings from configuration file
    - Template string precedence with multiple files
    - Invalid template strings (should use safe defaults)
    - Missing template string fields (should use defaults)
    - Template strings with Go text/template syntax
    - Template strings with %{name} placeholder syntax
    - Template strings with ANSI escape codes
- `TestRegexPatternConfiguration`: Tests regex pattern configuration loading
  - Validates regex pattern loading from YAML
  - Tests regex pattern defaults
  - Tests regex pattern precedence with multiple configuration files
  - Test cases:
    - Default regex patterns (should match immutable specification defaults)
    - Custom regex patterns from configuration file
    - Regex pattern precedence with multiple files
    - Invalid regex patterns (should use safe defaults)
    - Missing regex pattern fields (should use defaults)
    - Named capture group validation
    - Regex pattern compilation and performance
- `TestTemplateRegexIntegration`: Tests integration between templates and regex patterns
  - Validates data extraction and template application workflow
  - Test cases:
    - Archive filename parsing with template formatting
    - Backup filename parsing with template formatting
    - Timestamp parsing with template formatting
    - Configuration line parsing with template formatting
    - Complex data extraction scenarios
    - Error handling when regex patterns don't match
    - Fallback behavior for template formatting failures
    - Performance with large datasets

### Integration Tests
**Implementation**: `integration_test.go`
**Requirements**:
- `TestFullWorkflow`: Tests complete archive workflow
  - Validates end-to-end archive creation and listing
  - Test cases:
    - Create full archive and list
    - Create incremental archive and list
    - Verify archive after creation
    - Multiple archives with different configurations
- `TestFileBackupWorkflow`: Tests complete file backup workflow
  - Validates end-to-end file backup creation and listing
  - Test cases:
    - Create file backup and list
    - Multiple backups with different configurations
    - Backup with directory structure preservation
    - Backup with notes and timestamps
    - File backup with Git repository integration
    - File backup with custom configuration paths
    - File backup with environment variable configuration
    - File backup error handling and recovery
    - File backup with resource cleanup verification
    - File backup with context cancellation
    - File backup with template formatting
    - File backup with printf-style formatting
- `TestDryRunMode`: Tests dry-run functionality
  - Validates dry-run behavior across all commands
  - **Updated for printf-style formatting**:
    - Tests dry-run output with default format strings
    - Tests dry-run output with custom format strings
  - **Updated for template-based formatting**:
    - Tests template formatting with default template strings
    - Tests template formatting with custom template strings
    - Validates regex pattern integration with template formatting
    - Tests template string configuration loading and application
    - Validates data extraction and rich formatting capabilities
  - Test cases:
    - Dry-run full archive creation
    - Dry-run incremental archive creation
    - Dry-run file backup creation
    - No files created in dry-run mode
    - Proper output formatting in dry-run
    - Dry run with identical file (using `format_identical_backup`)
    - Dry run with modified file (using `format_dry_run_backup`)
    - Dry run with list flag (using `format_list_backup`)
    - Dry run with resource cleanup verification
    - Dry run with custom format string configuration
- `TestCLICommands`: Tests command-line interface
  - Validates all CLI commands and flags
  - **Updated for printf-style formatting**:
    - Tests error message formatting for invalid arguments
  - **Updated for template-based formatting**:
    - Tests template formatting with default template strings
    - Tests template formatting with custom template strings
    - Validates regex pattern integration with template formatting
    - Tests template string configuration loading and application
    - Validates data extraction and rich formatting capabilities
  - Test cases:
    - `bkpdir full` command
    - `bkpdir inc` command
    - `bkpdir list` command
    - `bkpdir verify` command
    - `bkpdir backup` command
    - `bkpdir --config` command (backward compatibility)
    - `bkpdir config` command (display configuration)
    - `bkpdir config KEY VALUE` command (set configuration)
    - `bkpdir --list [FILE_PATH]` command
    - `--dry-run` flag
    - `--checksum` flag for verify
    - Invalid flag combinations (using `format_error`)
    - Missing file path (using `format_error`)
    - Invalid file path (using `format_error`)
    - Config flag with other arguments (should be ignored)
    - Config set with invalid key (using `format_error`)
    - Config set with invalid value type (using `format_error`)
    - Config set with missing arguments (using `format_error`)
- `TestConfigurationIntegration`: Tests configuration in real scenarios
  - Validates configuration loading and application
  - **Updated for printf-style formatting**:
    - Tests configuration display with default format strings
    - Tests configuration display with custom format strings
    - Validates format string configuration display
  - **Updated for template-based formatting**:
    - Tests template formatting with default template strings
    - Tests template formatting with custom template strings
    - Validates regex pattern integration with template formatting
    - Tests template string configuration loading and application
    - Validates data extraction and rich formatting capabilities
  - Test cases:
    - Default configuration behavior
    - Custom configuration files
    - Environment variable configuration
    - Configuration precedence
    - Invalid configuration handling
    - Status code configuration integration
    - Format string configuration integration
    - Template configuration integration
    - Regex pattern configuration integration
    - Display config with default values only (including status codes and format strings)
    - Display config with values from single configuration file
    - Display config with values from multiple configuration files
    - Display config with `BKPDIR_CONFIG` environment variable set
    - Display config with hard-coded default path when environment variable not set
    - Display config with invalid configuration files (error handling using `format_error`)
    - Display config with custom status code values
    - Display config with custom format string values
    - Verify application exits after displaying configuration
    - Verify format string configuration is displayed with proper formatting
    - Tests backup operations with custom configuration paths
    - Tests environment variable override in real scenarios
    - Tests hard-coded default path behavior
    - Tests configuration precedence with actual file operations
    - Tests status code configuration in real application scenarios
- `TestConfigCommandIntegration`: Tests config command in real application scenarios
  - Validates configuration display and modification functionality
  - Tests configuration persistence and file handling
  - Tests error handling and validation in real scenarios
  - Test cases:
    - `bkpdir config` displays configuration correctly
    - `bkpdir config KEY VALUE` sets configuration values
    - Configuration changes persist across command invocations
    - Setting string values (archive_dir_path, backup_dir_path, checksum_algorithm)
    - Setting boolean values (use_current_dir_name, include_git_info, verify_on_create)
    - Setting integer values (status codes)
    - Creating new configuration file when none exists
    - Updating existing configuration file preserves other values
    - Nested configuration handling (verification section)
    - Error handling for invalid configuration keys
    - Error handling for invalid boolean values (not true/false)
    - Error handling for invalid integer values (non-numeric)
    - Error handling for file permission issues
    - Error handling for YAML parsing errors
    - Configuration sorting (alphabetical display)
    - Backward compatibility with `bkpdir --config`
    - Integration with other commands using modified configuration
    - Configuration file creation with proper permissions
    - YAML format preservation and readability
- `TestGitIntegrationWorkflow`: Tests Git integration in workflows
  - Validates Git information in archive and backup creation
  - Test cases:
    - Archive creation in Git repository
    - Archive creation in non-Git directory
    - Archive creation with different branches
    - Archive creation with dirty working directory
    - Backup creation in Git repository
- `TestVerificationWorkflow`: Tests verification in complete workflows
  - Validates verification integration with archive creation
  - Test cases:
    - Automatic verification after creation
    - Manual verification of existing archives
    - Verification with checksums
    - Verification failure scenarios
- `TestErrorHandlingIntegration`: Tests error handling in real scenarios
  - Validates error handling across all operations
  - Test cases:
    - Permission denied scenarios
    - Disk full scenarios
    - Invalid input scenarios
    - Network interruption scenarios (if applicable)
    - Context cancellation scenarios
- `TestConcurrentOperations`: Tests concurrent archive and backup operations
  - Validates thread safety and resource management
  - Test cases:
    - Multiple archive creations
    - Multiple backup creations
    - Concurrent verification operations
    - Resource cleanup under concurrency
    - No race conditions
- `TestContextIntegration`: Tests context-aware operations in workflows
  - Validates context cancellation and timeout handling
  - Test cases:
    - Archive creation with context cancellation
    - Backup creation with context cancellation
    - Long-running operations with timeouts
    - Resource cleanup after cancellation
- `TestResourceManagementIntegration`: Tests resource management in workflows
  - Validates resource cleanup across all operations
  - Test cases:
    - Archive creation with resource cleanup
    - Backup creation with resource cleanup
    - Error scenarios with cleanup verification
    - No temporary files left after operations
- `TestFormatStringIntegration`: Tests format string configuration in full application context
  - Tests all operations with custom format string configurations
  - Tests format string precedence with multiple configuration files
  - Tests format string error handling and fallbacks
  - Test cases:
    - Archive creation with custom format strings
    - Backup creation with custom format strings
    - List operation with custom format strings
    - Configuration display with custom format strings
    - Error handling with custom error format strings
    - Format string precedence with multiple configuration files
    - Invalid format strings falling back to safe defaults
    - ANSI color code rendering in various terminal environments
    - Format string validation and error handling
- `TestStatusCodeIntegration`: Tests status code configuration in full application context
  - Validates application exit codes match configuration
  - Tests status code behavior with various error conditions
  - Test cases:
    - Application exits with correct status code for successful archive
    - Application exits with correct status code for successful backup
    - Application exits with correct status code for identical directory
    - Application exits with correct status code for identical file
    - Application exits with correct status code for directory not found
    - Application exits with correct status code for file not found
    - Application exits with correct status code for permission denied
    - Application exits with correct status code for invalid directory type
    - Application exits with correct status code for invalid file type
    - Status code configuration from multiple files with precedence
    - Default status codes when no configuration provided
- `TestResourceCleanupIntegration`: Tests resource cleanup in full application context
  - Validates resource cleanup across entire application workflow
  - Tests cleanup with various error scenarios
  - Test cases:
    - Successful operations with cleanup verification
    - Failed operations with cleanup verification
    - Interrupted operations with cleanup verification
    - No temporary files left in any scenario
- `TestTemplateFormattingIntegration`: Tests template formatting configuration in full application context
  - Tests all operations with custom template string configurations
  - Tests template string precedence with multiple configuration files
  - Tests template string error handling and fallbacks
  - Tests regex pattern integration with template formatting
  - Test cases:
    - Archive creation with custom template strings and regex patterns
    - Backup creation with custom template strings and regex patterns
    - List operation with custom template strings and data extraction
    - Configuration display with custom template strings
    - Error handling with custom error template strings
    - Template string precedence with multiple configuration files
    - Invalid template strings falling back to safe defaults
    - Regex pattern matching and data extraction in various scenarios
    - Template formatting with missing or incomplete data
    - ANSI color code rendering in template-formatted output
    - Template validation and error handling
    - Performance with complex template strings and large datasets
- `TestRegexPatternIntegration`: Tests regex pattern configuration in full application context
  - Tests regex pattern loading and compilation
  - Tests named capture group extraction and usage
  - Tests pattern precedence with multiple configuration files
  - Test cases:
    - Archive filename parsing with custom regex patterns
    - Backup filename parsing with custom regex patterns
    - Timestamp parsing with custom regex patterns
    - Configuration line parsing with custom regex patterns
    - Invalid regex patterns falling back to safe defaults
    - Regex pattern precedence with multiple configuration files
    - Named capture group validation and extraction
    - Performance with complex regex patterns and large datasets
    - Error handling for malformed regex patterns

### Performance Tests
**Implementation**: `performance_test.go`
**Requirements**:
- `BenchmarkArchiveCreation`: Benchmarks archive creation performance
  - Measures performance with various directory sizes
  - Test cases:
    - Small directories (< 100 files)
    - Medium directories (100-1000 files)
    - Large directories (> 1000 files)
    - Deep directory hierarchies
- `BenchmarkFileBackupCreation`: Benchmarks file backup creation performance
  - Measures performance with various file sizes
  - Test cases:
    - Small files (< 1MB)
    - Medium files (1-100MB)
    - Large files (> 100MB)
    - Multiple file backups
    - File backup with context cancellation overhead
    - File backup with resource cleanup overhead
    - File backup directory structure preservation
    - File backup with notes and metadata
- `BenchmarkDirectoryComparison`: Benchmarks directory comparison performance
  - Measures comparison speed for various scenarios
  - Test cases:
    - Identical directories
    - Completely different directories
    - Directories with few changes
    - Large directories
- `BenchmarkFileComparison`: Benchmarks file comparison performance
  - Measures file comparison speed for various scenarios
  - Test cases:
    - Identical files
    - Different files
    - Large files
    - Multiple file comparisons
    - File comparison with size optimization
    - File comparison with byte-by-byte verification
- `BenchmarkVerification`: Benchmarks verification performance
  - Measures verification speed for various archive sizes
  - Test cases:
    - Structure verification only
    - Checksum verification
    - Large archives
    - Many small files vs few large files
- `BenchmarkGitOperations`: Benchmarks Git integration performance
  - Measures Git information extraction speed
  - Test cases:
    - Git repository detection
    - Branch name extraction
    - Commit hash extraction
    - Large Git repositories
- `BenchmarkTemplateFormatting`: Benchmarks template formatting performance
  - Measures template processing speed
  - Test cases:
    - Simple template formatting
    - Complex template with regex extraction
    - Multiple template applications
    - Large data sets
- `BenchmarkResourceManagement`: Benchmarks resource management performance
  - Measures resource tracking and cleanup speed
  - Test cases:
    - Resource registration
    - Resource cleanup
    - Concurrent resource operations
    - Large resource sets
- `BenchmarkCreateBackup`: Benchmarks backup creation performance
  - Tests performance with various file sizes
  - Measures resource cleanup overhead
  - Validates memory usage patterns
- `BenchmarkResourceManager`: Benchmarks resource management performance
  - Tests performance with many temporary resources
  - Measures cleanup time with large resource lists
  - Validates thread-safe performance
- `TestMemoryUsage`: Tests memory consumption
  - Validates memory usage stays within reasonable bounds
  - Test cases:
    - Large directory archiving
    - Large file backup creation
    - Long-running operations
    - Multiple concurrent operations
    - Memory leak detection
    - Resource cleanup memory impact

### Stress Tests
**Implementation**: `stress_test.go`
**Requirements**:
- `TestLargeDirectories`: Tests with very large directories
  - Validates handling of directories with many files
  - Test cases:
    - Directories with > 10,000 files
    - Deep directory hierarchies (> 20 levels)
    - Files with very long names
    - Directories with special characters
- `TestLargeFiles`: Tests with very large files
  - Validates handling of large file backups
  - Test cases:
    - Files > 1GB in size
    - Multiple large files
    - Large files with context cancellation
    - Large file comparison operations
    - Large file backup creation with resource cleanup
    - Large file backup listing and sorting
    - Memory usage during large file backup operations
    - Disk space monitoring during large file backups
- `TestLongRunningOperations`: Tests operations that take significant time
  - Validates timeout and cancellation handling
  - Test cases:
    - Archive creation with timeout
    - Backup creation with timeout
    - Verification with timeout
    - Context cancellation during operations
    - Resource cleanup after cancellation
    - Large file backup operations with timeout
    - File comparison operations with timeout
    - Backup listing operations with timeout
- `TestResourceExhaustion`: Tests behavior under resource constraints
  - Validates graceful handling of resource limitations
  - Test cases:
    - Low disk space scenarios
    - Memory pressure scenarios
    - File descriptor exhaustion
    - Temporary directory full
- `TestConcurrentStress`: Tests high concurrency scenarios
  - Validates thread safety under stress
  - Test cases:
    - Many concurrent archive operations
    - Many concurrent backup operations
    - Concurrent read/write operations
    - Resource contention scenarios
    - Cleanup under high concurrency
    - Concurrent file backup creation and listing
    - Concurrent file comparison operations
    - Mixed archive and backup operations under load
    - Resource manager thread safety under stress
- `TestContextStress`: Tests context handling under stress
  - Validates context cancellation under high load
  - Test cases:
    - Many concurrent operations with cancellation
    - Rapid context cancellation scenarios
    - Context timeout under load
    - Resource cleanup with context stress
    - Concurrent backup operations with cancellation
    - Mixed operation types with context cancellation
    - Context timeout during file backup operations
- Stress tests:
  - Concurrent backup operations
  - Large file handling
  - Many temporary resources
  - Resource cleanup under load

### Linting and Code Quality Tests
**Implementation**: `Makefile`, CI/CD pipeline
**Requirements**:
- `make lint`: Runs revive linter
  - Validates all Go code passes linting standards
  - Checks error handling compliance
  - Validates code style and formatting
  - Test cases:
    - All source files pass revive checks
    - No unhandled errors
    - Proper function and variable naming
    - Adequate documentation
- Code quality validation:
  - All `fmt.Printf`, `fmt.Fprintf` return values checked
  - All file operations handle errors appropriately
  - Consistent error handling patterns
  - Proper resource cleanup in all code paths

## Test Infrastructure

### Test Utilities
**Implementation**: `testutil/` package
**Requirements**:
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

### Test Data Management
**Requirements**:
- Temporary directories created for each test
- Test files with various sizes and types
- Test Git repositories with different states
- Test configuration files with various settings
- Test backup scenarios with different file types
- Cleanup verification after each test
- No test artifacts left in file system
- Resource tracking for all test operations
- Test files use consistent naming patterns
- Test data covers various file sizes and types
- Test data includes special characters and edge cases
- Test configuration files cover various YAML structures and edge cases
- Test configuration files for modification testing with various initial states
- Test configuration files with nested structures (verification section)
- Test configuration files with missing sections for creation testing
- Test configuration files with invalid YAML for error testing
- **Test configuration files include various format string configurations**
- **Test data includes format strings with ANSI color codes**
- **Test data includes invalid format strings for error handling**
- **Test data covers format string precedence scenarios**
- **Test data includes template strings with Go text/template syntax**
- **Test data includes template strings with %{name} placeholder syntax**
- **Test data includes regex patterns with named capture groups**
- **Test data includes invalid template strings for error handling**
- **Test data includes invalid regex patterns for error handling**
- **Test data covers template string precedence scenarios**
- **Test data covers regex pattern precedence scenarios**
- **Test data includes complex data extraction scenarios**
- Temporary resources are tracked and verified for cleanup
- Test data includes scenarios for context cancellation
- Test files simulate various error conditions
- Test data covers atomic operation scenarios

### Continuous Integration Requirements
**Requirements**:
- All tests must pass before merge
- Test coverage must be > 90%
- Performance benchmarks must not regress
- Memory usage tests must pass
- Concurrent operation tests must pass
- Resource cleanup verification required
- Context cancellation tests must pass
- Cross-platform testing (Linux, macOS)
- Linting must pass before code merge
- Code coverage must meet minimum thresholds
- Performance benchmarks must not regress
- Resource cleanup must be verified in CI environment
- Tests must run on multiple platforms
- Memory leak detection must be performed
- Static analysis must pass
- **Format string validation must pass in CI environment**
- **Output formatting tests must pass on all supported platforms**
- **ANSI color code rendering must be tested in CI environment**
- **Template string validation must pass in CI environment**
- **Template formatting tests must pass on all supported platforms**
- **Regex pattern compilation must be tested in CI environment**
- **Data extraction and template integration must be validated in CI**

## Test Execution

### Local Development
```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test categories
go test -run TestUnit ./...
go test -run TestIntegration ./...
go test -run Benchmark ./...

# Run tests with race detection
go test -race ./...

# Run stress tests
go test -tags=stress ./...

# Run context tests
go test -run TestContext ./...

# Run resource management tests
go test -run TestResource ./...

# Run template formatting tests
go test -run TestTemplate ./...

# Run output formatting tests
go test -run TestOutputFormatter ./...
go test -run TestFormatString ./...

# Run config command tests
go test -run TestConfigModification ./...
go test -run TestConfigCommand ./...

# Run benchmarks
go test -v -bench=.

# Run linting
make lint

# Run all quality checks
make test
```

### Test Categories by Build Tags
- **Unit Tests**: Default, no build tags required
- **Integration Tests**: `-tags=integration`
- **Performance Tests**: `-tags=performance`
- **Stress Tests**: `-tags=stress`
- **Context Tests**: `-tags=context`

### Test Environment Setup
**Requirements**:
- Temporary directory for test files
- Mock Git repositories
- Test configuration files
- Test backup directories
- Isolated test environment
- Cleanup verification
- Resource tracking
- Context management

## Quality Gates

### Code Coverage Requirements
- **Minimum Coverage**: 90% overall
- **Critical Functions**: 100% coverage required
  - Archive creation functions
  - File backup creation functions
  - Error handling functions
  - Resource cleanup functions
  - Configuration loading functions
  - Configuration modification functions
  - Context-aware functions
  - Template formatting functions
- **Integration Tests**: Must cover all CLI commands
- **Error Paths**: All error conditions must be tested
- **Context Paths**: All cancellation scenarios must be tested
- **Output formatting functionality must be comprehensively tested**
- **Format string configuration must be thoroughly tested**
- **Printf-style formatting must be validated for all output types**
- **ANSI color code support must be verified**
- **Format string validation and error handling must be tested**
- **Format string precedence and merging must be tested**
- **Default format string behavior must be verified**
- **Text highlighting and structure formatting must be validated**
- **Template-based formatting functionality must be comprehensively tested**
- **Template string configuration must be thoroughly tested**
- **Go text/template syntax must be validated for all output types**
- **Placeholder syntax (%{name}) must be validated for all output types**
- **Regex pattern integration must be verified**
- **Named capture group extraction must be tested**
- **Template string validation and error handling must be tested**
- **Template string precedence and merging must be tested**
- **Default template string behavior must be verified**
- **Data extraction and rich formatting capabilities must be validated**
- **Template formatting fallback behavior must be tested**
- Resource cleanup must be tested in all scenarios
- Context cancellation and timeout handling must be verified
- Enhanced error handling must be thoroughly tested
- Panic recovery must be validated
- Atomic operations must be tested
- Thread safety must be verified
- Performance characteristics must be measured
- Memory usage patterns must be validated

### Performance Requirements
- **Archive Creation**: < 1 second for 1000 files
- **File Backup Creation**: < 500ms for 100MB file
- **Directory Comparison**: < 500ms for 1000 files
- **File Comparison**: < 200ms for 100MB file
- **Verification**: < 2 seconds for 100MB archive
- **Template Formatting**: < 10ms for complex templates
- **Memory Usage**: < 100MB for 10,000 files
- **No Memory Leaks**: All resources properly cleaned up
- **Context Cancellation**: < 100ms response time

### Reliability Requirements
- **Resource Cleanup**: 100% cleanup verification
- **Error Recovery**: All error scenarios tested
- **Concurrent Safety**: No race conditions
- **Data Integrity**: Archive and backup corruption detection
- **Platform Compatibility**: Tests pass on Linux and macOS
- **Context Handling**: All cancellation scenarios handled
- **Template Robustness**: Invalid templates handled gracefully

## Test Maintenance

### Test Documentation
- Each test function must have clear documentation
- Test cases must be documented with expected outcomes
- Performance benchmarks must have baseline measurements
- Error scenarios must be documented with expected behavior
- Context scenarios must be documented with cancellation points
- Resource cleanup must be documented and verified

### Test Data Management
- Test data must be generated programmatically
- No hardcoded file paths or system dependencies
- Test isolation must be maintained
- Cleanup must be verified for all tests
- Resource tracking must be comprehensive
- Context management must be consistent

### Regression Testing
- All bug fixes must include regression tests
- Performance regressions must be detected
- Configuration changes must be tested
- Backward compatibility must be verified
- Context handling regressions must be detected
- Resource leak regressions must be detected

### Test Categories by Functionality

#### Archive Tests
- Archive creation (full and incremental)
- Archive listing and verification
- Archive naming with Git integration
- Archive comparison and deduplication

#### File Backup Tests
- File backup creation and listing
- File comparison and deduplication
- Backup directory structure preservation
- Backup naming with timestamps and notes

#### Configuration Tests
- Configuration loading and precedence
- Environment variable support
- Status code configuration
- Format string and template configuration
- Regex pattern configuration

#### Error Handling Tests
- Structured error creation and handling
- Enhanced error detection (disk space, permissions)
- Error context preservation
- Panic recovery and cleanup

#### Resource Management Tests
- Automatic resource cleanup
- Thread-safe resource tracking
- Atomic operations with temporary files
- Context-aware resource management

#### Template Formatting Tests
- Printf-style formatting
- Template-based formatting with placeholders
- Regex pattern extraction
- ANSI color support and text highlighting

#### Context Tests
- Context cancellation support
- Timeout handling
- Resource cleanup with cancellation
- Long-running operation cancellation

## Test Approach
**Implementation**: `*_test.go` files
**Requirements**:
- Uses temporary directories for testing
- Simulates user environment with test files
- Tests both dry-run and actual execution modes
- Verifies correct behavior for edge cases and error conditions
- Uses `CreateBackupWithTime` for consistent timestamps in tests
- Tests path handling for both absolute and relative paths
- Tests file comparison with various file types and sizes
- Verifies correct exit behavior when files are identical
- Tests flag-based command structure
- Validates proper reporting of existing backup names
- Tests environment variable handling with various configurations
- Creates multiple temporary configuration files for precedence testing
- Tests home directory expansion and path resolution
- Tests configuration discovery error conditions and edge cases
- Validates configuration display output format and source tracking
- Tests application exit behavior after configuration display
- Tests status code configuration loading and application
- Validates application exit codes match configured values
- Tests status code behavior with various error conditions and scenarios
- Creates test configuration files with custom status code values
- Verifies status code precedence with multiple configuration files
- Tests resource cleanup in all scenarios including failures
- Validates context cancellation and timeout handling
- Tests panic recovery and error resilience
- Verifies no resource leaks in any test scenario
- Tests atomic operations and data integrity
- Validates enhanced error detection and handling

## Test Environment
- Uses Go's testing package
- Leverages temporary directories for file operations
- Mocks time functions for consistent testing
- Simulates file system operations
- Tests on both macOS and Linux platforms
- Mocks environment variables for configuration testing
- Creates temporary configuration files with various content
- Uses context with timeouts for cancellation testing
- Simulates disk full and permission denied scenarios
- Tests with various file sizes and types
- Validates cleanup in all test scenarios

This testing architecture ensures comprehensive coverage of BkpDir functionality while maintaining high quality standards and reliable operation across different environments and use cases, including the new file backup capabilities, context-aware operations, enhanced error handling, and advanced template formatting features. 

### File Backup Integration Tests
**Implementation**: `backup_test.go`
**Requirements**:
- `TestFileBackupWorkflow`: Tests complete file backup workflow
  - Validates end-to-end file backup functionality
  - Test cases:
    - Complete backup creation workflow
    - Backup listing workflow
    - Identical file detection workflow
    - Multiple backup creation and management
    - Backup with various file types (text, binary, executable)
- `TestFileBackupWithGitIntegration`: Tests file backup in Git repositories
  - Validates Git integration with file backup operations
  - Test cases:
    - Backup creation in Git repository
    - Backup creation in non-Git directory
    - Git info handling in backup workflows
    - Branch and hash information preservation
- `TestFileBackupConcurrency`: Tests concurrent file backup operations
  - Validates thread-safety of file backup operations
  - Test cases:
    - Concurrent backup creation for different files
    - Concurrent backup creation for same file
    - Resource manager thread safety
    - No race conditions in backup listing
- `TestFileBackupResourceCleanup`: Tests resource cleanup in file backup operations
  - Validates comprehensive resource cleanup
  - Test cases:
    - Successful backup with cleanup verification
    - Failed backup with cleanup verification
    - Cancelled backup with cleanup verification
    - No temporary files left after operations
    - Atomic file operations verification
- `TestFileBackupErrorScenarios`: Tests error handling in file backup operations
  - Validates comprehensive error handling
  - Test cases:
    - Source file not found
    - Source path is directory
    - Permission denied on source file
    - Permission denied on backup directory
    - Disk full during backup creation
    - Backup directory creation failure
    - Invalid file types (devices, sockets, etc.)

### New Tests
- `TestGenerateBackupName`: Tests file backup naming
  - Validates backup filename generation with timestamps and notes
  - Test cases:
    - Basic backup naming without notes
    - Backup naming with notes
    - Various source file paths and extensions
    - Timestamp format validation
    - Special characters in filenames
    - Long filenames
    - Files with multiple extensions
- `TestListFileBackups`: Tests file backup listing
  - Validates backup listing and sorting for specific files
  - Test cases:
    - Empty backup directory
    - Multiple backups with different timestamps
    - Backups with notes
    - Backups in nested directory structures
    - Sorting by creation time (most recent first)
    - Files with special characters
    - Large number of backups
    - Permission denied scenarios
- `TestBackupDirectoryStructure`: Tests backup directory structure preservation
  - Validates that backup paths maintain source file directory structure
  - Test cases:
    - Files in root directory
    - Files in nested subdirectories
    - Files with complex directory paths
    - Files with special characters in paths
    - Absolute vs relative path handling
    - Current directory name inclusion/exclusion
- `TestBackupError`: Tests structured backup error handling
  - Validates BackupError functionality
  - Test cases:
    - Error creation with message and status code
    - Error interface implementation
    - Status code extraction
    - Error message formatting
    - Error context preservation (operation, path)
    - Error unwrapping
    - Backup-specific error scenarios 

### Testing Infrastructure and Complex Scenarios

#### Archive Corruption Testing Framework (TEST-INFRA-001-A) ✅ COMPLETED
**Implementation**: `internal/testutil/corruption.go`, `internal/testutil/corruption_test.go`
**Requirements**:
- `TestCorruptionType`: Tests corruption type enumeration (8 types: CRC, Header, Truncate, Central Directory, Local Header, Data, Signature, Comment)
- `TestCreateTestArchive`: Tests test archive creation utility with specified files
- `TestNewArchiveCorruptor`: Tests archive corruptor instantiation and validation
- `TestArchiveCorruptor_Backup`: Tests backup/restore functionality for safe corruption testing
- `TestArchiveCorruptor_RestoreFromBackup`: Tests restoration from backup after corruption
- `TestCorruptionCRC`: Tests systematic CRC checksum corruption
- `TestCorruptionHeader`: Tests ZIP file header corruption for unreadability
- `TestCorruptionTruncate`: Tests file truncation corruption with size validation
- `TestCorruptionCentralDirectory`: Tests central directory corruption (usually fatal)
- `TestCorruptionData`: Tests file data corruption within archives
- `TestCorruptionSignature`: Tests ZIP signature corruption for immediate failure
- `TestCorruptionComment`: Tests archive comment corruption (non-fatal)
- `TestCorruptionReproducibility`: Tests deterministic corruption with identical seeds
- `TestCorruptionDetector`: Tests automatic corruption type detection
- `TestCreateCorruptedTestArchive`: Tests convenience function for creating pre-corrupted archives
- `TestUnsupportedCorruptionType`: Tests error handling for invalid corruption types
- `TestEdgeCases`: Tests edge cases (empty archives, very small archives)
- `BenchmarkCorruptionCRC`: Performance benchmarks for CRC corruption (~763μs)
- `BenchmarkCorruptionDetection`: Performance benchmarks for corruption detection (~49μs)

**Test Coverage**:
- **Comprehensive corruption types**: All 8 corruption types tested systematically
- **Reproducibility validation**: Deterministic corruption using seed-based generation
- **Detection capabilities**: Automatic classification of corruption types
- **Performance benchmarks**: Baseline performance metrics for corruption operations
- **Edge case handling**: Empty archives, small archives, invalid operations
- **Error scenarios**: Comprehensive error handling and recovery testing
- **Cross-platform compatibility**: Works on Unix and Windows systems
- **Integration ready**: Can be imported and used by verification logic testing

**Key Features Implemented**:
1. **Controlled ZIP Corruption Utilities**:
   - `ArchiveCorruptor` class with backup/restore functionality
   - 8 systematic corruption types targeting different ZIP structures
   - Deterministic corruption patterns for reproducible test results
   - Safe corruption with automatic cleanup and restoration

2. **Corruption Type Enumeration**:
   - `CorruptionCRC`: Corrupt file checksums (recoverable)
   - `CorruptionHeader`: Corrupt ZIP file headers (fatal)
   - `CorruptionTruncate`: Cut off end of archive (permanent data loss)
   - `CorruptionCentralDir`: Corrupt ZIP central directory (fatal)
   - `CorruptionLocalHeader`: Corrupt individual file headers (sometimes recoverable)
   - `CorruptionData`: Corrupt actual file data (detectable/recoverable)
   - `CorruptionSignature`: Corrupt ZIP file signatures (fatal)
   - `CorruptionComment`: Corrupt archive comments (non-fatal)

3. **Deterministic Corruption Patterns**:
   - Seed-based corruption for consistent test results
   - Offset-based variation for multiple corruption points
   - Reproducible corruption across identical archives
   - Tracked original bytes for potential recovery scenarios

4. **Archive Repair Detection**:
   - `CorruptionDetector` for automatic corruption type identification
   - Tests recovery behavior from various corruption types
   - Classification of corruption severity (recoverable vs fatal)
   - Integration with Go's ZIP reader resilience testing

**Implementation Notes**:
- **Go ZIP Reader Resilience**: Go's `zip.OpenReader` proved surprisingly resilient to corruption, requiring test adjustments to verify corruption effects rather than complete failure
- **Reproducibility Challenges**: Initially had problems with identical corruption due to archive structure differences; fixed by ensuring byte-identical archives before corruption
- **Detection Logic Refinement**: Required handling multiple corruption types being detected simultaneously
- **Performance Characteristics**: CRC corruption ~763μs, detection ~49μs for typical archives
- **Cross-Platform Testing**: Successfully tested on Unix systems with proper handling of file permissions and temporary directories

**Usage in Verification Testing**:
```go
// Create corrupted archive for testing verification logic
config := CorruptionConfig{
    Type: CorruptionCRC,
    Seed: 12345, // For reproducible tests
}

result, err := CreateCorruptedTestArchive(archivePath, testFiles, config)
if err != nil {
    t.Fatalf("Failed to create corrupted archive: %v", err)
}

// Test verification behavior with corrupted archive
detector := NewCorruptionDetector(archivePath)
detected, err := detector.DetectCorruption()
// ... verify detection results match expected corruption
```

**Integration with Existing Test Suite**:
- Located in `internal/testutil/` to provide reusable testing infrastructure
- Can be imported by `verify.go` and `comparison.go` tests for complex verification scenarios
- Provides foundation for testing edge cases that are difficult to reproduce reliably
- Enables systematic testing of archive corruption scenarios previously untestable

#### Disk Space Simulation Framework (TEST-INFRA-001-B) ✅ COMPLETED
**Implementation**: `internal/testutil/diskspace.go`, `internal/testutil/diskspace_test.go`
**Requirements**:
- `TestDiskSpaceSimulator`: Tests disk space simulation with quota limits and progressive exhaustion
- `TestFileSystemInterface`: Tests filesystem interface abstraction and wrapper functionality
- `TestExhaustionModes`: Tests 4 exhaustion modes (Linear, Progressive, Random, Immediate)
- `TestSpaceConstraintEnforcement`: Tests quota limits and space recovery scenarios
- `TestErrorInjectionFramework`: Tests operation-specific failure points and file-path-specific injection
- `TestPredefinedScenarios`: Tests built-in scenarios ("GradualExhaustion", "ImmediateFailure", "SpaceRecovery", "ProgressiveExhaustion")
- `TestDynamicConfiguration`: Tests adding/removing failure points, recovery points, and injection points at runtime
- `TestConcurrentAccess`: Tests thread-safe operations with mutex protection and concurrent access patterns
- `TestSpaceAccounting`: Tests tracking of available and used space with accurate recovery from file deletions
- `BenchmarkSpaceChecking`: Performance benchmarks for space checking operations
- `BenchmarkConcurrentAccess`: Performance benchmarks for concurrent access patterns

**Test Coverage**:
- **Comprehensive disk space simulation**: 4 exhaustion modes with configurable rates and thresholds
- **Mock filesystem interface**: FileSystemInterface abstraction with RealFileSystem and SimulatedFileSystem implementations
- **Space constraint enforcement**: Quota limits, progressive exhaustion, and space recovery without requiring large test files
- **Error injection framework**: Operation-specific failure points, file-path-specific injection, and space recovery points
- **Thread-safe operations**: All operations protected with mutex for concurrent access
- **Comprehensive statistics**: Operations, failures, space changes, and file sizes tracking for test validation
- **Predefined scenarios**: Built-in test scenarios with expected results validation
- **Helper functions**: SimulateArchiveCreation and SimulateBackupOperation for realistic testing
- **Dynamic configuration**: Runtime support for adding/removing failure points, recovery points, and injection points
- **Performance optimized**: Efficient space checking and concurrent access patterns
- **Error type support**: Standard disk full (ENOSPC), quota exceeded, and device full errors
- **Space accounting**: Proper tracking of available and used space with accurate recovery from file deletions

**Key Features Implemented**:
1. **Mock Filesystem with Quota Limits**:
   - `FileSystemInterface` abstraction with injectable wrappers
   - `SimulatedFileSystem` implementation with configurable space constraints
   - Controlled disk space simulation without affecting real system
   - Support for quota limits and progressive space exhaustion

2. **Progressive Space Exhaustion**:
   - Linear exhaustion: Constant rate space reduction
   - Progressive exhaustion: Accelerating space reduction over time
   - Random exhaustion: Unpredictable space reduction patterns
   - Immediate exhaustion: Instant space exhaustion at trigger points

3. **Disk Full Error Injection**:
   - ENOSPC error triggering at specific operation points
   - File-path-specific injection for targeted testing
   - Operation-specific failure points for precise control
   - Configurable failure rates and timing

4. **Space Recovery Testing**:
   - Test behavior when disk space becomes available again
   - Space recovery points for testing resilience
   - Accurate space accounting after file deletions
   - Recovery scenario validation

**Implementation Notes**:
- **Thread-Safe Design**: All operations protected with mutex for concurrent access safety
- **Performance Characteristics**: Efficient space checking operations with minimal overhead
- **Cross-Platform Support**: Works on Unix and Windows systems with appropriate error type mapping
- **Integration Ready**: Can be imported and used by archive creation and backup testing

**Usage in Archive and Backup Testing**:
```go
// Create disk space simulator for testing space exhaustion
simulator := NewDiskSpaceSimulator(QuotaLimit{MaxBytes: 1024 * 1024}) // 1MB limit
fs := NewSimulatedFileSystem(simulator)

// Add failure point for archive creation
simulator.AddFailurePoint("archive_creation", 0.8) // Fail at 80% capacity

// Test archive creation with space constraints
result := SimulateArchiveCreation(fs, testDirectory, archivePath)
// ... verify behavior under space constraints
```

#### Permission Testing Framework (TEST-INFRA-001-C) ✅ COMPLETED
**Implementation**: `internal/testutil/permissions.go` (756 lines), `internal/testutil/permissions_test.go` (609 lines)
**Requirements**:
- `TestPermissionSimulator`: Tests permission scenario generator with systematic permission combinations
- `TestCrossPlatformPermissions`: Tests Unix/Windows permission differences and simulation
- `TestPermissionRestoration`: Tests safe permission restoration utilities and atomic rollback
- `TestPermissionChangeDetection`: Tests behavior when permissions change during operations
- `TestPermissionScenarios`: Tests built-in scenarios (basic_permission_denial, directory_access_denied, mixed_permissions)
- `TestPermissionTestHelper`: Tests high-level PermissionTestHelper for easy testing integration
- `TestThreadSafeOperations`: Tests mutex protection and statistics tracking
- `TestPermissionCombinations`: Tests systematic permission combinations (0000-0777)
- `TestPermissionErrorHandling`: Tests permission error scenarios and error aggregation
- `BenchmarkPermissionSimulation`: Performance benchmarks for permission simulation operations
- `BenchmarkPermissionRestoration`: Performance benchmarks for permission restoration operations

**Test Coverage**:
- **Permission scenario generator**: Systematic permission combinations for comprehensive testing
- **Cross-platform support**: Unix/Windows permission differences handling
- **Restoration utilities**: Safe permission restoration with atomic rollback
- **Permission change detection**: Testing behavior when permissions change during operations
- **Built-in scenarios**: Predefined permission scenarios for common testing patterns
- **High-level integration**: PermissionTestHelper for easy testing integration
- **Thread-safe operations**: Mutex protection and statistics tracking
- **Error handling**: Comprehensive permission error scenarios and aggregation
- **Performance optimization**: Efficient permission simulation and restoration operations
- **Temporary directory isolation**: Safe testing without modifying system files

**Key Features Implemented**:
1. **Permission Scenario Generator**:
   - `PermissionSimulator` with systematic permission combinations (0000-0777)
   - Configurable permission patterns for files and directories
   - Support for mixed permission scenarios within directory trees
   - Temporary directory isolation for safe testing

2. **Cross-Platform Permission Simulation**:
   - Unix permission handling with octal notation (rwxrwxrwx)
   - Windows permission simulation using Go's os.Chmod limitations
   - Platform-specific permission validation and error handling
   - Consistent permission behavior across operating systems

3. **Permission Restoration Utilities**:
   - Atomic permission rollback for safe test cleanup
   - Original permission preservation and restoration
   - Comprehensive permission restoration with error handling
   - Restoration utilities that safely restore original permissions after tests

4. **Permission Change Detection**:
   - Monitor permission changes during operations
   - Detect unauthorized permission modifications
   - Test behavior when permissions change during file operations
   - Validation of permission preservation in file operations

**Implementation Notes**:
- **COMPLETED**: Critical for testing permission error paths in file operations and configuration management
- **Thread-Safe Design**: All operations protected with mutex for concurrent access
- **Temporary Directory Isolation**: Uses temporary directories with controlled permissions rather than modifying system files
- **Design Patterns**: Internal lock-free methods to avoid deadlocks, comprehensive error aggregation
- **Test Coverage**: 14 comprehensive tests including benchmarks, all passing

**Usage in File Operations and Configuration Testing**:
```go
// Create permission test helper for configuration file testing
helper := NewPermissionTestHelper()
defer helper.Cleanup()

// Create scenario with denied directory access
scenario := PermissionScenario{
    Path: configDir,
    Mode: 0000, // No access
    Type: PermissionTypeBoth,
}

// Apply permission restrictions and test behavior
helper.ApplyScenario(scenario)
result := LoadConfigFromDirectory(configDir)
// ... verify proper error handling for permission denied scenarios
```

#### Context Cancellation Testing Helpers (TEST-INFRA-001-D) ✅ COMPLETED
**Implementation**: `internal/testutil/context.go` (623 lines), `internal/testutil/context_test.go` (832 lines)
**Requirements**:
- `TestContextController`: Tests controlled timeout scenarios with precise timing control
- `TestCancellationManager`: Tests cancellation point injection and management
- `TestConcurrentOperations`: Tests cancellation during concurrent archive/backup operations
- `TestPropagationVerification`: Tests context propagation through operation chains
- `TestCancellationScenarios`: Tests helper functions for common scenarios (archive, backup, resource cleanup)
- `TestEventTracking`: Tests comprehensive event tracking and statistics collection
- `TestIntegrationUtilities`: Tests integration testing utilities combining multiple cancellation scenarios
- `TestTickerBasedTiming`: Tests ticker-based timing control and goroutine coordination
- `TestAtomicCounters`: Tests atomic counters for thread safety
- `TestConcurrentCancellation`: Tests cancellation during concurrent operations
- `BenchmarkContextController`: Performance benchmarks for context control operations
- `BenchmarkCancellationManager`: Performance benchmarks for cancellation management
- `BenchmarkPropagationChain`: Performance benchmarks for context propagation testing

**Test Coverage**:
- **Controlled timeout scenarios**: Precise timing control for context cancellation testing
- **Cancellation point injection**: Trigger cancellation at specific operation stages
- **Concurrent operation testing**: Test cancellation during concurrent archive/backup operations
- **Cancellation propagation verification**: Ensure proper context propagation through operation chains
- **Helper functions**: Common scenarios for archive, backup, and resource cleanup cancellation points
- **Event tracking and statistics**: Comprehensive analysis of cancellation behavior
- **Integration testing utilities**: Combining multiple cancellation scenarios for complex testing
- **Performance optimization**: Efficient context control and cancellation management
- **Thread safety**: Atomic counters and goroutine coordination for reliable testing

**Key Features Implemented**:
1. **Controlled Timeout Scenarios**:
   - `ContextController` with precise timing control for delayed cancellation scenarios
   - Configurable delays and timeout management
   - Ticker-based timing control for deterministic cancellation testing
   - Support for multiple timeout patterns and cancellation triggers

2. **Cancellation Point Injection**:
   - `CancellationManager` for injection point management and statistics tracking
   - Enable/disable functionality for specific cancellation points
   - Execution counting and cancellation trigger management
   - Flexible injection point configuration for targeted testing

3. **Concurrent Operation Testing**:
   - `ConcurrentTestConfig` for testing cancellation during concurrent operations
   - Coordination of multiple goroutines with cancellation scenarios
   - Thread-safe cancellation testing with proper synchronization
   - Concurrent operation patterns for archive and backup testing

4. **Cancellation Propagation Verification**:
   - `PropagationTestConfig` for verifying context propagation through operation chains
   - Multi-level context propagation with detailed operation tracking
   - `PropagationChain` for testing complex operation hierarchies
   - Verification of proper context handling across component boundaries

**Implementation Notes**:
- **COMPLETED**: Comprehensive context cancellation testing infrastructure for reliable testing of context-aware operations
- **Core Components**: ContextController, CancellationManager, PropagationChain, and helper utilities
- **Design Patterns**: Ticker-based timing, goroutine coordination, atomic counters for thread safety
- **Notable Features**: Event logging, statistics tracking, concurrent operation testing, propagation verification
- **Test Coverage**: 18 comprehensive test functions including integration tests, stress tests, and benchmarks

**Usage in Context-Aware Operations Testing**:
```go
// Create context controller for testing delayed cancellation
controller := NewContextController()
ctx := controller.CreateCancellableContext(context.Background(), 100*time.Millisecond)

// Test archive creation with context cancellation
result := CreateArchiveWithContext(ctx, sourceDir, archivePath)
// ... verify proper cancellation handling and resource cleanup

// Test context propagation through operation chain
chain := NewPropagationChain()
chain.AddOperation("validate", ValidateDirectoryWithContext)
chain.AddOperation("create", CreateArchiveWithContext)
chain.AddOperation("verify", VerifyArchiveWithContext)

propagationResult := chain.ExecuteWithCancellation(ctx, testParams)
// ... verify context is properly propagated through all operations
```

#### Error Injection Framework (TEST-INFRA-001-E) ✅ COMPLETED
**Implementation**: `internal/testutil/errorinjection.go` (720 lines), `internal/testutil/errorinjection_test.go` (650 lines)
**Requirements**:
- `TestErrorInjectionInterfaces`: Tests systematic error injection points for filesystem, Git, and network operations
- `TestErrorClassification`: Tests error type categorization (transient, permanent, recoverable, fatal)
- `TestErrorPropagationTracing`: Tests error flow tracking through operation chains
- `TestErrorRecoveryTesting`: Tests retry logic and graceful degradation
- `TestInjectionControl`: Tests trigger counts, max injections, probability-based injection, conditional functions
- `TestBuiltInScenarios`: Tests pre-configured scenarios for filesystem and Git error testing
- `TestPerformanceCharacteristics`: Tests error injection performance (~896ns/op) and propagation tracking (~214ns/op)
- `TestThreadSafeOperations`: Tests RWMutex protection for concurrent access
- `TestErrorTypes`: Tests 4 error types (Transient, Permanent, Recoverable, Fatal) and 6 categories
- `TestInjectionSchedules`: Tests configurable error schedules and timing delays
- `BenchmarkErrorInjection`: Performance benchmarks for error injection operations
- `BenchmarkPropagationTracking`: Performance benchmarks for error propagation tracking
- `BenchmarkRecoveryTesting`: Performance benchmarks for error recovery operations

**Test Coverage**:
- **Systematic error injection points**: Configurable error insertion at filesystem, Git, and network operations
- **Error type classification**: Categorize errors (transient, permanent, recoverable, fatal)
- **Error propagation tracing**: Track error flow through operation chains
- **Error recovery testing**: Test retry logic and graceful degradation
- **Advanced injection control**: Trigger counts, max injections, probability-based injection, conditional functions, timing delays
- **Built-in scenarios**: Pre-configured scenarios for filesystem and Git error testing with expected results validation
- **Performance optimization**: Excellent performance with error injection and propagation tracking
- **Thread-safe operations**: RWMutex protection for concurrent access
- **Comprehensive testing**: 16 test functions plus 3 benchmarks achieving 100% coverage

**Key Features Implemented**:
1. **Systematic Error Injection Points**:
   - Interface-based design with comprehensive interfaces for filesystem, Git, and network operations
   - Injectable wrappers for controlled error insertion
   - Configurable error insertion at specific operation points
   - Support for operation-specific and file-path-specific error injection

2. **Error Type Classification**:
   - 4 error types: Transient (temporary), Permanent (persistent), Recoverable (can retry), Fatal (cannot continue)
   - 6 error categories: Filesystem, Git, Network, Permission, Resource, Configuration
   - Systematic error categorization for testing different error handling strategies
   - Error classification system for appropriate error response testing

3. **Error Propagation Tracing**:
   - Complete error flow tracking through operation chains with timestamps and stack depth
   - Error propagation analysis for understanding error handling patterns
   - Stack depth tracking for complex operation hierarchies
   - Timestamp tracking for error timing analysis

4. **Error Recovery Testing**:
   - Built-in recovery attempt tracking with success/failure statistics
   - Timing analysis for recovery operations
   - Test retry logic and graceful degradation under error conditions
   - Recovery scenario validation with expected results

**Implementation Notes**:
- **Interface-based design**: Created comprehensive interfaces for filesystem, Git, and network operations with injectable wrappers
- **Advanced injection control**: Support for trigger counts, max injections, probability-based injection, conditional functions, and timing delays
- **Performance**: Excellent performance with error injection at ~896ns/op and propagation tracking at ~214ns/op
- **Thread-safe operations**: All operations protected with RWMutex for concurrent access
- **Comprehensive testing**: 16 test functions plus 3 benchmarks achieving 100% coverage of core functionality

**Usage in Error Handling Pattern Testing**:
```go
// Create error injector for filesystem operations
injector := NewFilesystemErrorInjector()
injector.AddInjectionPoint("file_write", ErrorTypeTransient, 0.3) // 30% chance of transient error

// Test archive creation with error injection
fs := NewInjectableFileSystem(realFS, injector)
result := CreateArchiveWithErrorInjection(fs, sourceDir, archivePath)

// Analyze error propagation and recovery
trace := injector.GetPropagationTrace()
recoveryStats := injector.GetRecoveryStatistics()
// ... verify proper error handling and recovery behavior
```

#### Integration Testing Orchestration (TEST-INFRA-001-F) ✅ COMPLETED
**Implementation**: `internal/testutil/scenarios.go` (1,100+ lines), `internal/testutil/scenarios_test.go` (900+ lines)
**Requirements**:
- `TestScenarioBuilder`: Tests complex scenario composition with builder pattern
- `TestPhaseBasedExecution`: Tests clear separation into Setup, Execution, Verification, and Cleanup phases
- `TestParallelStepSupport`: Tests configurable parallel execution within phases for performance testing
- `TestTimingCoordination`: Tests TimingCoordinator with barriers, signals, and delays for synchronized testing
- `TestScenarioRuntime`: Tests complete runtime environment with temp directories, config files, and shared data
- `TestErrorIntegration`: Tests full integration with existing error injection, corruption, permission, and disk space components
- `TestRegressionIntegration`: Tests demonstrated integration patterns with existing test suite structure
- `TestScenarioExecution`: Tests end-to-end scenario execution with comprehensive validation
- `TestBuilderPatterns`: Tests fluent interface with ScenarioBuilder for intuitive scenario construction
- `TestOrchestration`: Tests coordination of multiple error conditions and operations
- `BenchmarkScenarioExecution`: Performance benchmarks for scenario execution
- `BenchmarkBuilderConstruction`: Performance benchmarks for scenario builder operations
- `BenchmarkTimingCoordination`: Performance benchmarks for timing coordination operations

**Test Coverage**:
- **Complex scenario composition**: Combine multiple error conditions for realistic testing
- **Test scenario scripting**: Declarative scenario definition for complex multi-step tests
- **Timing and coordination utilities**: Synchronize multiple error conditions and operations
- **Regression test suite integration**: Plug infrastructure into existing test suites
- **Builder pattern implementation**: Fluent interface for intuitive scenario construction
- **Phase-based execution**: Clear separation between setup, execution, verification, and cleanup
- **Parallel execution support**: Configurable parallel execution for performance testing
- **Runtime environment**: Complete scenario runtime with comprehensive testing support
- **Error integration**: Full integration with all existing testing infrastructure components
- **Performance optimization**: Efficient scenario execution and coordination operations

**Key Features Implemented**:
1. **Complex Scenario Composition**:
   - `ScenarioBuilder` with fluent interface for intuitive scenario construction
   - Support for combining multiple error conditions, corruption scenarios, permission restrictions
   - Declarative scenario definition for complex multi-step tests
   - Flexible scenario composition with reusable components

2. **Test Scenario Scripting**:
   - Declarative scenario definition language for complex testing workflows
   - Multi-step test scenarios with dependencies and conditional execution
   - Scenario scripting with clear separation of concerns
   - Reusable scenario templates for common testing patterns

3. **Timing and Coordination Utilities**:
   - `TimingCoordinator` with barriers, signals, and delays for synchronized testing
   - Precise timing control for complex multi-component testing
   - Synchronization primitives for coordinating multiple error conditions and operations
   - Timing utilities for performance testing and stress testing

4. **Regression Test Suite Integration**:
   - Demonstrated integration patterns with existing test suite structure
   - Plug infrastructure into existing `*_test.go` files
   - Backward compatibility with existing testing patterns
   - Integration examples for adopting scenario-based testing

**Implementation Notes**:
- **Builder Pattern**: Implemented fluent interface with `ScenarioBuilder` for intuitive scenario construction
- **Phase-based Execution**: Clear separation into Setup, Execution, Verification, and Cleanup phases
- **Parallel Step Support**: Configurable parallel execution within phases for performance testing
- **Runtime Environment**: Complete `ScenarioRuntime` with temp directories, config files, and shared data
- **Error Integration**: Full integration with existing error injection, corruption, permission, and disk space components
- **Comprehensive Testing**: 15 test functions covering builder patterns, orchestration, parallel execution, timing coordination
- **Performance Benchmarks**: Included benchmark tests for scenario execution performance

**Usage in Integration Testing**:
```go
// Create complex integration scenario combining multiple testing infrastructure components
scenario := NewScenarioBuilder("archive_stress_test").
    WithCorruption(CorruptionConfig{Type: CorruptionCRC, Seed: 12345}).
    WithDiskSpaceLimit(1024 * 1024). // 1MB limit
    WithPermissionRestrictions(0644). // Limited permissions
    WithErrorInjection(ErrorInjectionConfig{
        Type: ErrorTypeTransient,
        Probability: 0.1, // 10% error rate
    }).
    WithCancellationScenario(100 * time.Millisecond). // Cancel after 100ms
    Build()

// Execute complex scenario with full orchestration
runtime := NewScenarioRuntime()
result := scenario.Execute(runtime)

// Validate comprehensive scenario results
assert.True(t, result.AllPhasesCompleted())
assert.NoError(t, result.GetExecutionError())
assert.True(t, result.GetPerformanceMetrics().WithinExpectedBounds())
```

// ... existing code ... 

## 🔻 CI/CD Pipeline Testing Requirements

### CICD-001: AI-First Development Optimization Testing
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// CICD-001: AI-first CI/CD optimization`  
**Test Function**: `TestAIWorkflow`

#### Core Testing Requirements

**🧪 Pipeline Validation Testing**
- **Test Coverage**: All CI/CD pipeline stages and components
- **Priority Testing**: Low-priority task scheduling validation
- **Background Processing**: Non-blocking execution verification
- **Resource Management**: CPU/memory usage within limits

**🛡️ AI Protocol Integration Testing**
- **DOC-008 Integration**: Icon validation system integration tests
- **Token Compliance**: Implementation token validation testing
- **Documentation Sync**: Cross-reference consistency testing
- **Protocol Adherence**: ai-assistant-protocol.md compliance verification

**🔧 Quality Gate Testing**
- **Automated Gates**: Zero-human intervention validation
- **Failure Recovery**: Self-healing pipeline configuration tests
- **Threshold Configuration**: AI-optimized quality metrics testing
- **Error Reporting**: AI-friendly error message format validation

#### Test Implementation Strategy

**📋 Test Categories**
```go
// 🔻 CICD-001: AI workflow pipeline testing - 🧪 Test infrastructure
func TestAIWorkflow(t *testing.T) {
    // Test suite for AI-optimized CI/CD pipeline
}

func TestPipelinePriorityScheduling(t *testing.T) {
    // Verify low-priority task execution
}

func TestBackgroundProcessing(t *testing.T) {
    // Validate non-blocking execution
}

func TestAIProtocolValidation(t *testing.T) {
    // Test DOC-008 integration and token compliance
}

func TestQualityGateAutomation(t *testing.T) {
    // Verify automated quality gates without human intervention
}

func TestAIOptimizedReporting(t *testing.T) {
    // Test standardized icon output and AI-readable status messages
}
```

**🔍 Integration Test Scenarios**
- **Full Pipeline Execution**: End-to-end AI workflow testing
- **Failure Scenarios**: Error handling and recovery testing
- **Performance Testing**: Background resource usage validation
- **Concurrency Testing**: Parallel AI assistant workflow support

**📊 Validation Test Requirements**
- **Icon Validation**: DOC-007/DOC-008 compliance testing
- **Token Validation**: Implementation token format verification
- **Documentation Consistency**: Cross-file reference validation
- **Output Format Testing**: AI-readable report format validation

#### Test Environment Configuration

**🚀 Test Infrastructure**
- **Isolated Environment**: Separate CI/CD pipeline testing environment
- **Mock AI Workflows**: Simulated AI assistant development cycles
- **Resource Monitoring**: Background task resource usage tracking
- **Automated Cleanup**: Test artifact removal and environment reset

**📈 Test Metrics and Coverage**
- **Pipeline Coverage**: All CI/CD stages tested
- **Error Path Coverage**: All failure scenarios validated
- **Integration Coverage**: All external system integrations tested
- **Performance Metrics**: Resource usage and execution time validation

#### Compliance Testing

**🛡️ AI Assistant Protocol Compliance**
- Validation against ai-assistant-compliance.md requirements
- Token referencing requirement testing
- Documentation update requirement validation
- Cross-reference consistency verification

**📋 Feature Tracking Integration**
- Automatic status update testing
- Token-based change attribution validation
- Documentation synchronization testing

## 🔻 AI-First Documentation and Code Maintenance Testing

### DOC-013: AI-First Documentation Strategy Testing
**Priority**: 🔻 LOW  
**Implementation Tokens**: `// DOC-013: AI-first maintenance`  
**Test Function**: `TestAIDocumentation`

#### Core Testing Requirements

**🔍 AI Comprehension Testing**
- **Test Coverage**: All AI-centric documentation components and workflows
- **Icon Integration**: DOC-007/DOC-008 icon standardization validation
- **Token Format**: Implementation token format compliance testing
- **Cross-Reference**: Machine-readable link validation and traversal testing

**📋 Documentation Pattern Testing**
- **Pattern Consistency**: Validation of predictable documentation patterns
- **Terminology Consistency**: Verification of standardized vocabulary usage
- **Hierarchy Testing**: Information architecture validation for AI navigation
- **Organization Testing**: Content structure optimization verification

**🔧 Workflow Integration Testing**
- **Automation Testing**: AI assistant maintenance procedure validation
- **Update Testing**: Automated documentation update verification
- **Validation Testing**: Cross-reference consistency checking
- **Generation Testing**: AI-friendly content creation pattern validation

#### Test Implementation Strategy

**📋 Test Categories**
```go
// 🔻 DOC-013: AI documentation strategy testing - 🧪 Test infrastructure
func TestAIDocumentation(t *testing.T) {
    // Test suite for AI-first documentation system
}

func TestAIComprehensionOptimization(t *testing.T) {
    // Verify AI-friendly formatting and structure
}

func TestDocumentationPatternConsistency(t *testing.T) {
    // Validate predictable documentation patterns
}

func TestCrossReferenceValidation(t *testing.T) {
    // Test machine-readable link validation and traversal
}

func TestAIWorkflowIntegration(t *testing.T) {
    // Verify AI assistant maintenance procedures
}

func TestContentGenerationPatterns(t *testing.T) {
    // Test AI-friendly content creation patterns
}

func TestTerminologyConsistency(t *testing.T) {
    // Validate standardized vocabulary usage
}
```

**🔍 AI Integration Test Scenarios**
- **Documentation Creation**: AI assistant content generation testing
- **Cross-Reference Maintenance**: Automated link validation and repair testing
- **Pattern Recognition**: AI comprehension and pattern generation testing
- **Workflow Execution**: AI assistant maintenance task execution testing

**📊 AI Comprehension Validation Requirements**
- **Icon Standardization**: DOC-007/DOC-008 integration testing
- **Token Format Compliance**: Implementation token standardization validation
- **Cross-Reference Integrity**: Machine-readable link validation
- **Content Pattern Testing**: AI-friendly documentation pattern verification

#### Test Environment Configuration

**🚀 AI Documentation Test Infrastructure**
- **Mock AI Assistant Environment**: Simulated AI assistant documentation workflows
- **Pattern Recognition Testing**: Validation of AI comprehension patterns
- **Cross-Reference Testing**: Automated link validation and traversal
- **Content Generation Testing**: AI-friendly template and pattern testing

**📈 AI Documentation Metrics and Coverage**
- **Comprehension Coverage**: All AI-centric documentation features tested
- **Pattern Coverage**: All documentation patterns validated for AI understanding
- **Integration Coverage**: All AI assistant workflow integrations tested
- **Quality Metrics**: AI comprehension optimization and efficiency validation

#### AI Assistant Compliance Testing

**🛡️ AI-First Documentation Compliance**
- Validation against ai-assistant-compliance.md requirements for documentation
- AI-friendly formatting requirement testing
- Cross-reference traversal requirement validation
- Implementation token integration verification

**📋 Documentation Quality Assurance**
- AI comprehension optimization testing
- Content pattern consistency validation
- Cross-reference integrity verification
- AI assistant workflow integration testing

#### Performance and Efficiency Testing

**🚀 AI Assistant Performance**
- Documentation navigation efficiency testing
- Cross-reference traversal performance validation
- Content generation speed and accuracy testing
- Pattern recognition and application performance

**📊 AI Comprehension Analytics**
- Content accessibility metrics validation
- Pattern recognition success rate testing
- Cross-reference integrity monitoring
- AI assistant task completion efficiency measurement

## 🔻 AI-First Documentation and Code Maintenance Testing

### DOC-011: Token Validation Integration for AI Assistants Testing
**Priority**: 🔺 HIGH  
**Implementation Tokens**: `// DOC-011: AI validation integration`  
**Test Function**: `TestAIValidation`

#### Core Testing Requirements

**🔧 AI Workflow Integration Testing**
- **Test Coverage**: All AI workflow validation hooks and integration points
- **DOC-008 Integration**: Seamless integration with existing validation framework
- **Pre-submission Validation**: API endpoint functionality and workflow integration
- **Zero-friction Integration**: Validation workflow non-disruption verification

**🛡️ AI Error Processing Testing**
- **Error Message Formatting**: AI-readable error structure and content validation
- **Remediation Guidance**: Context-aware guidance generation and accuracy testing
- **Error Categorization**: Priority-based error classification validation
- **Response API Testing**: Structured API response format and content verification

**📊 Bypass Mechanism Testing**
- **Safe Override Testing**: Controlled bypass mechanism functionality validation
- **Audit Trail Testing**: Comprehensive tracking and logging verification
- **Approval Workflow Testing**: Automated approval process validation
- **Rollback Testing**: Reversible bypass action verification and recovery testing

**🔍 Compliance Monitoring Testing**
- **Behavior Tracking**: AI assistant validation behavior pattern analysis
- **Adherence Monitoring**: Real-time compliance measurement and reporting
- **Dashboard Testing**: Compliance dashboard generation and accuracy validation
- **Reporting Testing**: Automated compliance report generation and content verification

#### Test Implementation Strategy

**📋 Test Categories**
```go
// 🔺 DOC-011: AI validation integration testing - 🧪 Test infrastructure
func TestAIValidation(t *testing.T) {
    // Test suite for AI validation integration framework
}

func TestAIWorkflowValidationHooks(t *testing.T) {
    // Verify seamless DOC-008 integration and pre-submission validation
}

func TestAIErrorProcessingSystem(t *testing.T) {
    // Test AI-optimized error formatting and remediation guidance
}

func TestBypassMechanismFramework(t *testing.T) {
    // Validate safe override controls and audit trail management
}

func TestComplianceMonitoringInfrastructure(t *testing.T) {
    // Test AI behavior tracking and adherence monitoring
}

func TestValidationAPIGateway(t *testing.T) {
    // Test pre-submission validation API endpoints and responses
}

func TestAIOptimizedErrorReporting(t *testing.T) {
    // Verify AI-readable error messages and structured responses
}

func TestValidationWorkflowIntegration(t *testing.T) {
    // Test Makefile, CLI, and API integration points
}
```

**🔍 AI Validation Integration Test Scenarios**
- **Pre-submission Validation**: End-to-end validation workflow testing
- **Error Processing Scenarios**: AI error reporting and remediation testing
- **Bypass Workflow Testing**: Safe override mechanism and audit trail testing
- **Compliance Monitoring**: AI assistant behavior tracking and reporting testing

**📊 Validation API Testing Requirements**
- **Request Processing**: ValidationRequest structure and parameter validation
- **Response Generation**: ValidationResponse format and content verification
- **Error Handling**: API error response format and AI-readable content
- **Batch Processing**: Multiple file validation request handling

#### API Integration Testing

**🚀 Validation Gateway Testing**
```go
func TestValidationRequestProcessing(t *testing.T) {
    request := ValidationRequest{
        SourceFiles:    []string{"main.go", "config.go"},
        ChangedTokens:  []string{"ARCH-001", "CFG-002"},
        ValidationMode: "standard",
        RequestContext: &AIRequestContext{
            AssistantID: "test-ai",
            SessionID:   "session-123",
        },
    }
    
    response, err := validationGateway.ProcessRequest(request)
    if err != nil {
        t.Fatalf("Validation request failed: %v", err)
    }
    
    // Verify response structure and content
    assert.Equal(t, "pass", response.Status)
    assert.NotNil(t, response.RemediationSteps)
    assert.GreaterOrEqual(t, response.ComplianceScore, 0.0)
}
```

## ⭐ Configuration System Enhancement Testing

### CFG-005: Layered Configuration Inheritance Testing
**Priority**: ⭐ CRITICAL  
**Implementation Tokens**: `// ⭐ CFG-005: Layered configuration inheritance`  
**Test Function**: `TestConfigInheritance`

#### Core Testing Requirements

**🔧 Inheritance Mechanism Testing**
- **Test Coverage**: All inheritance declaration parsing and chain building functionality
- **Circular Dependency Detection**: Comprehensive cycle detection with various inheritance patterns
- **Path Resolution Testing**: Absolute, relative, and home directory expansion validation
- **Backward Compatibility**: Non-inheritance configuration loading preserved

**🛡️ Merge Strategy Testing**
- **Standard Override Testing**: Child values completely replace parent values
- **Array Merge Testing**: Child arrays append to parent arrays (`+` prefix)
- **Array Prepend Testing**: Child arrays prepend to parent arrays (`^` prefix)
- **Array Replace Testing**: Child arrays replace parent arrays (`!` prefix)
- **Default Value Testing**: Child values used only if parent not set (`=` prefix)

**📊 Configuration Processing Testing**
- **Dependency Order Testing**: Parent configurations loaded before children
- **Type Safety Testing**: Configuration merging preserves all data types
- **Source Tracking Testing**: Value origin visibility throughout inheritance chain
- **Error Handling Testing**: Missing files, malformed syntax, invalid declarations

**🔍 Performance and Reliability Testing**
- **Performance Testing**: O(n) complexity validation for inheritance processing
- **Thread Safety Testing**: Concurrent configuration loading verification
- **Memory Usage Testing**: Linear scaling with inheritance chain depth
- **Cache Testing**: Inheritance caching functionality and invalidation

#### Test Implementation Strategy

**📋 Test Categories**
```go
// ⭐ CFG-005: Configuration inheritance testing - 🧪 Test infrastructure
func TestConfigInheritance(t *testing.T) {
    // Test suite for layered configuration inheritance system
}

func TestInheritanceChainBuilding(t *testing.T) {
    // Verify inheritance dependency graph construction and ordering
}

func TestMergeStrategyProcessing(t *testing.T) {
    // Test all merge strategies with complex data structures
}

func TestCircularDependencyDetection(t *testing.T) {
    // Validate cycle detection with various inheritance patterns
}

func TestInheritancePathResolution(t *testing.T) {
    // Test absolute, relative, and home directory path expansion
}

func TestConfigurationSourceTracking(t *testing.T) {
    // Verify value origin tracking throughout inheritance chains
}

func TestInheritanceErrorHandling(t *testing.T) {
    // Test error scenarios and recovery mechanisms
}

func TestInheritancePerformance(t *testing.T) {
    // Validate O(n) complexity and memory usage characteristics
}
```

**🔍 Inheritance Test Scenarios**
- **Simple Inheritance**: Parent-child configuration relationships
- **Multi-level Inheritance**: Deep inheritance chains (3+ levels)
- **Multiple Parents**: Configuration inheriting from multiple sources
- **Complex Merge Scenarios**: Nested structures with various merge strategies

**📊 Merge Strategy Validation Requirements**
- **Array Operations**: All array merge strategies with complex nested arrays
- **Object Merging**: Nested object inheritance with strategy application
- **Primitive Overrides**: String, boolean, and numeric value inheritance
- **Mixed Type Handling**: Configuration with diverse data type combinations

#### Configuration Test Data

**🚀 Test Configuration Hierarchy**
```yaml
# test/configs/base.yml (root configuration)
archive_dir_path: "~/Archives"
backup_dir_path: "~/Backups"
exclude_patterns:
  - "*.tmp"
  - "*.log"
include_git_info: true
verification:
  verify_on_create: false
  checksum_algorithm: "sha256"

# test/configs/development.yml (inherits from base)
inherit: "base.yml"
archive_dir_path: "./dev-archives"
+exclude_patterns:
  - "node_modules/"
  - "dist/"
=verification:
  verify_on_create: true

# test/configs/project.yml (inherits from development)
inherit: "development.yml"
^exclude_patterns:
  - "*.secret"
!status_codes:
  - created: 0
  - identical: 0
```

**Expected Resolution Result Testing**
```go
func TestConfigurationResolution(t *testing.T) {
    loader := NewInheritanceConfigLoader()
    config, err := loader.LoadConfigWithInheritance("test/configs/project.yml")
    
    assert.NoError(t, err)
    assert.Equal(t, "./dev-archives", config.ArchiveDirPath)
    
    expectedPatterns := []string{
        "*.secret",        // ^ prepend strategy
        "*.tmp",           // base configuration
        "*.log",           // base configuration
        "node_modules/",   // + merge strategy
        "dist/",           // + merge strategy
    }
    assert.Equal(t, expectedPatterns, config.ExcludePatterns)
    
    assert.True(t, config.Verification.VerifyOnCreate) // = default strategy
}
```

#### Error Handling Test Cases

**🚨 Error Scenario Testing**
```go
func TestInheritanceErrorScenarios(t *testing.T) {
    testCases := []struct {
        name          string
        configFile    string
        expectedError string
    }{
        {
            name:          "Circular Dependency",
            configFile:    "circular-a.yml",
            expectedError: "circular inheritance detected: a.yml -> b.yml -> a.yml",
        },
        {
            name:          "Missing Parent File",
            configFile:    "missing-parent.yml",
            expectedError: "inheritance file not found: nonexistent.yml",
        },
        {
            name:          "Invalid Syntax",
            configFile:    "invalid-syntax.yml",
            expectedError: "invalid inheritance declaration",
        },
        {
            name:          "Invalid Merge Prefix",
            configFile:    "invalid-prefix.yml",
            expectedError: "unknown merge strategy prefix: @",
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            loader := NewInheritanceConfigLoader()
            _, err := loader.LoadConfigWithInheritance(tc.configFile)
            
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tc.expectedError)
        })
    }
}
```

#### Performance and Scalability Testing

**⚡ Performance Validation**
```go
func TestInheritancePerformanceCharacteristics(t *testing.T) {
    // Test inheritance chain depth scaling
    depths := []int{1, 5, 10, 20, 50}
    
    for _, depth := range depths {
        t.Run(fmt.Sprintf("depth_%d", depth), func(t *testing.T) {
            configs := createInheritanceChain(depth)
            
            start := time.Now()
            loader := NewInheritanceConfigLoader()
            _, err := loader.LoadConfigWithInheritance(configs[depth-1])
            duration := time.Since(start)
            
            assert.NoError(t, err)
            
            // Verify O(n) complexity - processing time should scale linearly
            expectedMaxDuration := time.Duration(depth) * 10 * time.Millisecond
            assert.Less(t, duration, expectedMaxDuration,
                "Inheritance processing time exceeded O(n) expectations")
        })
    }
}

func TestInheritanceMemoryUsage(t *testing.T) {
    // Verify memory usage scales linearly with chain depth
    var memBefore, memAfter runtime.MemStats
    
    runtime.GC()
    runtime.ReadMemStats(&memBefore)
    
    // Load deep inheritance chain
    loader := NewInheritanceConfigLoader()
    configs := createInheritanceChain(100)
    _, err := loader.LoadConfigWithInheritance(configs[99])
    
    runtime.ReadMemStats(&memAfter)
    
    assert.NoError(t, err)
    
    memUsage := memAfter.Alloc - memBefore.Alloc
    maxExpectedUsage := uint64(1024 * 1024) // 1MB max for 100-level chain
    
    assert.Less(t, memUsage, maxExpectedUsage,
        "Memory usage exceeded linear scaling expectations")
}
```

#### Integration Testing with Existing Features

**🔗 Feature Integration Testing**
```go
func TestInheritanceWithEnvironmentVariables(t *testing.T) {
    // Test inheritance combined with environment variable overrides
    os.Setenv("BKPDIR_ARCHIVE_DIR", "/env/override")
    defer os.Unsetenv("BKPDIR_ARCHIVE_DIR")
    
    loader := NewInheritanceConfigLoader()
    config, err := loader.LoadConfigWithInheritance("inherit-test.yml")
    
    assert.NoError(t, err)
    assert.Equal(t, "/env/override", config.ArchiveDirPath,
        "Environment variables should override inheritance chain")
}

func TestInheritanceWithConfigDiscovery(t *testing.T) {
    // Test inheritance with existing configuration discovery system (CFG-001)
    os.Setenv("BKPDIR_CONFIG", "base.yml:child.yml")
    defer os.Unsetenv("BKPDIR_CONFIG")
    
    loader := NewInheritanceConfigLoader()
    config, err := loader.LoadConfigWithInheritance(".")
    
    assert.NoError(t, err)
    // Verify inheritance works with discovery system
}

func TestInheritanceBackwardCompatibility(t *testing.T) {
    // Verify non-inheritance configurations continue working
    loader := NewInheritanceConfigLoader()
    config, err := loader.LoadConfigWithInheritance("legacy-config.yml")
    
    assert.NoError(t, err)
    assert.Equal(t, "legacy-behavior", config.SomeField,
        "Legacy configurations must continue working without inheritance")
}
```

This comprehensive testing framework ensures the layered configuration inheritance system works reliably across all scenarios while maintaining backward compatibility and performance characteristics.

### CFG-006: Complete Configuration Reflection and Visibility Testing
**Priority**: 🔺 HIGH  
**Implementation Tokens**: `// 🔺 CFG-006: Complete config visibility`  
**Test Function**: `TestConfigReflection`

#### Core Testing Requirements

**🔍 Automatic Field Discovery Testing**
- **Reflection Accuracy**: Comprehensive validation of Go struct field enumeration for all data types
- **Nested Structure Handling**: Deep nested structs, embedded fields, slices, maps, and complex types
- **Type Safety Testing**: Correct Go type information extraction and preservation
- **Performance Testing**: Field discovery caching and reflection optimization validation
- **Edge Case Handling**: Anonymous fields, unexported fields, and circular references

**📊 Source Attribution and Integration Testing**
- **CFG-005 Integration**: Complete integration with inheritance system source tracking
- **Source Chain Accuracy**: Correct attribution of configuration values to their sources
- **Merge Strategy Visibility**: Display of applied merge strategies and resolution paths
- **Conflict Detection**: Identification and reporting of configuration value conflicts
- **Hierarchical Resolution**: Environment → inheritance chain → defaults resolution tracking

**🎛️ Configuration Command Interface Testing**
- **Multiple Output Formats**: Table, tree, JSON, and YAML output format validation
- **Filtering Functionality**: All filtering options (--all, --overrides-only, --sources, --filter)
- **Type-Aware Formatting**: Correct display formatting for different configuration data types
- **Backward Compatibility**: Existing config command functionality preservation
- **CLI Flag Processing**: Command-line flag parsing and validation

**⚡ Performance and Optimization Testing**
- **Reflection Caching**: Field metadata caching accuracy and invalidation
- **Lazy Evaluation**: Source evaluation performance for displayed fields only
- **Memory Usage**: Configuration inspection memory efficiency validation
- **Incremental Resolution**: Partial configuration inspection performance

#### Test Implementation Strategy

**📋 Test Categories**
```go
// 🔺 CFG-006: Complete config visibility testing - 🧪 Test infrastructure  
func TestConfigReflection(t *testing.T) {
    // Test suite for complete configuration reflection and visibility system
}

func TestFieldDiscoveryAccuracy(t *testing.T) {
    // Verify Go reflection field enumeration for all data types
}

func TestSourceAttributionIntegration(t *testing.T) {
    // Test integration with CFG-005 inheritance source tracking
}

func TestConfigurationDisplayFormats(t *testing.T) {
    // Validate all output formats (table, tree, JSON, YAML)
}

func TestFilteringAndDisplayOptions(t *testing.T) {
    // Test all configuration filtering and display options
}

func TestTypeAwareFormatting(t *testing.T) {
    // Verify correct display formatting for different data types
}

func TestConfigurationCaching(t *testing.T) {
    // Validate field metadata and source attribution caching
}

func TestLazyEvaluationPerformance(t *testing.T) {
    // Test performance of lazy source evaluation
}

func TestConfigInspectionErrorHandling(t *testing.T) {
    // Test error scenarios and graceful recovery mechanisms
}
```

**🔬 Field Discovery Validation Tests**
```go
func TestComplexTypeReflection(t *testing.T) {
    // Test configuration with all Go data types
    testConfig := &TestConfig{
        StringField:    "test",
        IntField:      42,
        BoolField:     true,
        SliceField:    []string{"a", "b", "c"},
        MapField:      map[string]int{"key": 1},
        StructField:   NestedStruct{Value: "nested"},
        EmbeddedField: EmbeddedStruct{EmbeddedValue: "embedded"},
        PointerField:  &SomeStruct{PointerValue: "pointer"},
    }
    
    inspector := NewConfigurationInspector()
    fields, err := inspector.GetAllConfigFields(testConfig)
    
    assert.NoError(t, err)
    assert.Len(t, fields, 8) // All fields discovered
    
    // Verify specific field metadata
    stringField := findField(fields, "StringField")
    assert.Equal(t, reflect.String, stringField.Kind)
    assert.Equal(t, "test", stringField.CurrentValue)
    
    sliceField := findField(fields, "SliceField")
    assert.Equal(t, reflect.Slice, sliceField.Kind)
    assert.Equal(t, []string{"a", "b", "c"}, sliceField.CurrentValue)
    
    mapField := findField(fields, "MapField")
    assert.Equal(t, reflect.Map, mapField.Kind)
    assert.Equal(t, map[string]int{"key": 1}, mapField.CurrentValue)
}

func TestNestedStructReflection(t *testing.T) {
    testConfig := &Config{
        Archive: ArchiveConfig{
            DirPath: "/archives",
            Settings: ArchiveSettings{
                VerifyOnCreate: true,
                Compression: CompressionConfig{
                    Level:     9,
                    Algorithm: "gzip",
                },
            },
        },
    }
    
    inspector := NewConfigurationInspector()
    fields, err := inspector.GetAllConfigFields(testConfig)
    
    assert.NoError(t, err)
    
    // Verify nested field paths
    archiveDirField := findFieldByPath(fields, "Archive.DirPath")
    assert.NotNil(t, archiveDirField)
    assert.Equal(t, "/archives", archiveDirField.CurrentValue)
    
    compressionLevelField := findFieldByPath(fields, "Archive.Settings.Compression.Level")
    assert.NotNil(t, compressionLevelField)
    assert.Equal(t, 9, compressionLevelField.CurrentValue)
}
```

**🔗 Source Attribution Integration Tests**
```go
func TestInheritanceSourceIntegration(t *testing.T) {
    // Create inheritance chain for testing
    createTestInheritanceChain(t)
    
    inspector := NewConfigurationInspector()
    configValue, err := inspector.GetConfigValueWithCompleteSource("exclude_patterns")
    
    assert.NoError(t, err)
    assert.Len(t, configValue.Sources, 3) // Base + Development + Project
    
    // Verify source chain attribution
    sources := configValue.Sources
    assert.Equal(t, "config_file", sources[0].Type)
    assert.Contains(t, sources[0].Location, "base.yml")
    assert.Equal(t, []string{"*.tmp", "*.log"}, sources[0].Value)
    
    assert.Equal(t, "inheritance", sources[1].Type)
    assert.Contains(t, sources[1].Location, "development.yml")
    assert.Equal(t, "+", sources[1].MergeStrategy) // Append strategy
    
    assert.Equal(t, "inheritance", sources[2].Type)
    assert.Contains(t, sources[2].Location, "project.yml")
    assert.Equal(t, "^", sources[2].MergeStrategy) // Prepend strategy
    
    // Verify final resolved value
    expectedValue := []string{
        "*.secret",      // ^ prepend from project.yml
        "*.tmp",         // base from base.yml
        "*.log",         // base from base.yml
        "node_modules/", // + append from development.yml
        "dist/",         // + append from development.yml
    }
    assert.Equal(t, expectedValue, configValue.CurrentValue)
}

func TestResolutionPathTracking(t *testing.T) {
    inspector := NewConfigurationInspector()
    configValue, err := inspector.GetConfigValueWithCompleteSource("archive_dir_path")
    
    assert.NoError(t, err)
    assert.Len(t, configValue.ResolutionPath, 2) // Base → Override
    
    steps := configValue.ResolutionPath
    assert.Equal(t, "~/Archives", steps[0].PreviousValue)
    assert.Equal(t, "./project-archives", steps[0].NewValue)
    assert.Equal(t, "override", steps[0].Operation)
    assert.Contains(t, steps[0].Explanation, "Child configuration overrides parent value")
}
```

**🎨 Display Format Validation Tests**
```go
func TestTableFormatOutput(t *testing.T) {
    inspector := NewConfigurationInspector()
    displayEngine := NewConfigurationDisplayEngine()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        ShowAll:      true,
        ShowSources:  true,
        OutputFormat: "table",
    })
    
    assert.NoError(t, err)
    
    output := displayEngine.RenderTable(tree)
    
    // Verify table format structure
    lines := strings.Split(output, "\n")
    assert.Contains(t, lines[0], "Field Name") // Header row
    assert.Contains(t, lines[0], "Current Value")
    assert.Contains(t, lines[0], "Source")
    
    // Verify data rows
    assert.Contains(t, output, "archive_dir_path")
    assert.Contains(t, output, "./project-archives")  
    assert.Contains(t, output, "project.yml")
}

func TestTreeFormatOutput(t *testing.T) {
    inspector := NewConfigurationInspector()
    displayEngine := NewConfigurationDisplayEngine()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        ShowAll:      true,
        ShowSources:  true,
        OutputFormat: "tree",
    })
    
    assert.NoError(t, err)
    
    output := displayEngine.RenderTree(tree)
    
    // Verify tree structure symbols
    assert.Contains(t, output, "├─") // Tree branch symbols
    assert.Contains(t, output, "└─") // Tree end symbols
    assert.Contains(t, output, "│ ") // Tree continuation symbols
    
    // Verify hierarchical structure
    lines := strings.Split(output, "\n")
    archiveLine := findLineContaining(lines, "archive_dir_path")
    sourceLine := findLineContaining(lines, "Source:")
    
    // Source lines should be indented more than field lines
    archiveIndent := len(archiveLine) - len(strings.TrimLeft(archiveLine, " │├└─"))
    sourceIndent := len(sourceLine) - len(strings.TrimLeft(sourceLine, " │├└─"))
    assert.Greater(t, sourceIndent, archiveIndent)
}

func TestJSONFormatOutput(t *testing.T) {
    inspector := NewConfigurationInspector()
    displayEngine := NewConfigurationDisplayEngine()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        ShowAll:      true,
        OutputFormat: "json",
    })
    
    assert.NoError(t, err)
    
    output := displayEngine.RenderJSON(tree)
    
    // Verify valid JSON structure
    var jsonResult map[string]interface{}
    err = json.Unmarshal([]byte(output), &jsonResult)
    assert.NoError(t, err)
    
    // Verify JSON content structure
    assert.Contains(t, jsonResult, "fields")
    assert.Contains(t, jsonResult, "sources")
    
    fields := jsonResult["fields"].([]interface{})
    assert.Greater(t, len(fields), 0)
    
    firstField := fields[0].(map[string]interface{})
    assert.Contains(t, firstField, "name")
    assert.Contains(t, firstField, "value")
    assert.Contains(t, firstField, "type")
}
```

**🎯 Filtering and Display Option Tests**
```go
func TestOverridesOnlyFilter(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        OverridesOnly: true,
        ShowSources:   true,
    })
    
    assert.NoError(t, err)
    
    // Should only show fields that have been overridden from defaults
    for _, node := range tree.Root.Children {
        hasOverride := false
        for _, source := range node.Value.Sources {
            if source.Type != "default" {
                hasOverride = true
                break
            }
        }
        assert.True(t, hasOverride, 
            "OverridesOnly filter should only show non-default values")
    }
}

func TestFieldPatternFilter(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        FieldPattern: "exclude_*",
        ShowAll:      true,
    })
    
    assert.NoError(t, err)
    
    // Should only show fields matching the pattern
    for _, node := range tree.Root.Children {
        assert.True(t, strings.HasPrefix(node.Field.Name, "exclude_"),
            "Field pattern filter should only show matching fields")
    }
}

func TestMaxDepthFilter(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    tree, err := inspector.GetConfigurationTree(&DisplayFilter{
        MaxDepth: 2,
        ShowAll:  true,
    })
    
    assert.NoError(t, err)
    
    // Verify depth limitation
    maxDepthFound := calculateTreeDepth(tree.Root)
    assert.LessOrEqual(t, maxDepthFound, 2,
        "MaxDepth filter should limit tree depth")
}
```

**⚡ Performance and Caching Tests**
```go
func TestFieldMetadataCaching(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    // First call - should populate cache
    start1 := time.Now()
    fields1, err1 := inspector.GetAllConfigFields(&Config{})
    duration1 := time.Since(start1)
    
    assert.NoError(t, err1)
    
    // Second call - should use cache
    start2 := time.Now()
    fields2, err2 := inspector.GetAllConfigFields(&Config{})
    duration2 := time.Since(start2)
    
    assert.NoError(t, err2)
    assert.Equal(t, fields1, fields2)
    
    // Cached call should be significantly faster
    assert.Less(t, duration2, duration1/2,
        "Cached field discovery should be at least 2x faster")
}

func TestLazySourceEvaluation(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    // Request only specific field - should not evaluate others
    start := time.Now()
    configValue, err := inspector.GetConfigValueWithCompleteSource("archive_dir_path")
    duration := time.Since(start)
    
    assert.NoError(t, err)
    assert.NotNil(t, configValue)
    
    // Should be fast since only one field was evaluated
    maxExpectedDuration := 50 * time.Millisecond
    assert.Less(t, duration, maxExpectedDuration,
        "Lazy evaluation should be fast for single field requests")
}

func TestIncrementalResolution(t *testing.T) {
    inspector := NewConfigurationInspector()
    resolver := NewIncrementalResolver(&Config{})
    
    // Request fields incrementally
    fields := []string{"archive_dir_path", "backup_dir_path", "exclude_patterns"}
    
    for i, field := range fields {
        configValue, err := resolver.ResolveField(field)
        assert.NoError(t, err)
        assert.NotNil(t, configValue)
        
        // Verify cache grows incrementally
        assert.Equal(t, i+1, len(resolver.resolvedFields))
    }
}
```

**🔧 Error Handling and Recovery Tests**
```go
func TestConfigInspectionErrorRecovery(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    // Test with malformed configuration
    malformedConfig := &MalformedConfig{
        BrokenField: make(chan int), // Unsupported type
    }
    
    fields, err := inspector.GetAllConfigFields(malformedConfig)
    
    // Should recover gracefully and provide partial results
    assert.Error(t, err)
    assert.NotNil(t, fields) // Partial results should be available
    
    // Error should be descriptive
    configErr, ok := err.(*ConfigInspectionError)
    assert.True(t, ok)
    assert.Equal(t, "reflection_error", configErr.Type)
    assert.Contains(t, configErr.Recovery, "Skip unsupported field types")
}

func TestPartialSuccessHandling(t *testing.T) {
    inspector := NewConfigurationInspector()
    
    // Configuration with mix of valid and problematic fields
    mixedConfig := &MixedConfig{
        ValidField:   "valid",
        InvalidField: func() {}, // Unsupported function type
        AnotherValid: 42,
    }
    
    fields, err := inspector.GetAllConfigFields(mixedConfig)
    
    // Should succeed with warnings
    assert.Error(t, err) // Error for problematic fields
    assert.Len(t, fields, 2) // Should get valid fields
    
    // Check for specific valid fields
    validField := findField(fields, "ValidField")
    assert.NotNil(t, validField)
    assert.Equal(t, "valid", validField.CurrentValue)
    
    anotherValidField := findField(fields, "AnotherValid")
    assert.NotNil(t, anotherValidField)
    assert.Equal(t, 42, anotherValidField.CurrentValue)
}
```

#### Integration Testing with Existing Systems

**🔗 System Integration Tests**
```go
func TestCFG005IntegrationCompliance(t *testing.T) {
    // Verify complete integration with CFG-005 inheritance system
    inspector := NewConfigurationInspector()
    
    // Load configuration with inheritance
    configValue, err := inspector.GetConfigValueWithCompleteSource("exclude_patterns")
    
    assert.NoError(t, err)
    
    // Verify CFG-005 source tracking integration
    assert.Greater(t, len(configValue.Sources), 1) // Multiple sources from inheritance
    assert.NotEmpty(t, configValue.MergeHistory)   // Merge operations recorded
    assert.NotEmpty(t, configValue.ResolutionPath) // Resolution steps tracked
    
    // Verify inheritance chain visibility
    for _, source := range configValue.Sources {
        assert.NotEmpty(t, source.Location)     // Source file identified
        assert.NotEmpty(t, source.MergeStrategy) // Merge strategy recorded
        assert.GreaterOrEqual(t, source.ChainDepth, 0) // Chain depth tracked
    }
}

func TestEXTRACT001PkgConfigIntegration(t *testing.T) {
    // Verify integration with EXTRACT-001 pkg/config architecture
    inspector := NewConfigurationInspector()
    
    // Should work with existing pkg/config structures
    fields, err := inspector.GetAllConfigFields(&pkg.Config{})
    
    assert.NoError(t, err)
    assert.Greater(t, len(fields), 0)
    
    // Verify pkg/config field compatibility
    for _, field := range fields {
        assert.NotEmpty(t, field.Name)
        assert.NotNil(t, field.Type)
        // Should handle pkg/config specific types correctly
    }
}

func TestCLIFrameworkIntegration(t *testing.T) {
    // Test CLI command integration
    cmd := NewConfigCommand()
    
    // Test various command flag combinations
    testCases := []struct {
        args     []string
        expected string
    }{
        {[]string{"config", "--all"}, "all fields"},
        {[]string{"config", "--overrides-only"}, "overridden values only"},
        {[]string{"config", "--format=json"}, "JSON format"},
        {[]string{"config", "--filter=exclude_*"}, "filtered fields"},
    }
    
    for _, tc := range testCases {
        t.Run(strings.Join(tc.args, "_"), func(t *testing.T) {
            output, err := executeCommand(cmd, tc.args...)
            
            assert.NoError(t, err)
            assert.Contains(t, output, tc.expected)
        })
    }
}
```

This comprehensive testing framework ensures the complete configuration reflection and visibility system provides accurate, performant, and reliable configuration inspection capabilities while maintaining full integration with existing systems.