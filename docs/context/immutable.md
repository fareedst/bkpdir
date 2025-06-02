# Immutable Specifications

This document contains specifications that MUST NOT be changed without a major version bump. These are core behaviors that users and other systems depend on.

## Archive Naming Convention
- Format: `{prefix}-{timestamp}-{git_info}-{note}.zip`
- `prefix`: Optional prefix (default: "bkp")
- `timestamp`: ISO 8601 format (YYYY-MM-DDTHHmmss)
- `git_info`: Optional Git branch and short hash (e.g., "main-abc123")
- `note`: Optional note (default: empty)
- Example: `bkp-2024-03-20T143022-main-abc123-initial.zip`
- This naming convention is fixed and must not be modified

## File Backup Naming Convention
- Format: `{filename}-{timestamp}[={note}]`
- `filename`: Original filename without path
- `timestamp`: YYYY-MM-DD-HH-MM format
- `note`: Optional note appended with equals sign
- Examples:
  - `document.txt-2024-03-20-14-30`
  - `config.yml-2024-03-20-14-30=before-changes`
  - `script.sh-2024-03-20-14-30=working-version`
- This naming convention is fixed and must not be modified

## Directory Operations
- **Platform Independence**: Use platform-independent path handling
- **Permission Preservation**: Preserve file permissions and modification time
- **Atomic Operations**: All operations must be atomic to prevent corruption
- **Structure Preservation**: Maintain directory structure in archives
- **Special File Handling**: Handle symbolic links, devices, sockets, etc.
- **Automatic Creation**: Create archive directories if they don't exist
- **Path Display**: Display all paths relative to current directory
- These file operation rules are fundamental and must not be altered

## File Backup Operations
- **Atomic File Operations**: All file backup operations must be atomic
- **Identical File Detection**: Must compare files byte-by-byte to detect identical backups
- **Directory Structure Preservation**: Maintain source file's directory structure in backup path
- **Timestamp Preservation**: Preserve original file modification times
- **Permission Preservation**: Preserve file permissions where applicable
- **Path Resolution**: Handle both absolute and relative file paths
- **Backup Directory Creation**: Create backup directories if they don't exist
- These file backup operation rules are fundamental and must not be altered

## File Exclusion Requirements
- **Pattern Matching**: Doublestar glob pattern matching
- **Configuration**: Configurable exclusion patterns
- **System Files**: Default exclusions for system files
- **Case Sensitivity**: Case-sensitive matching
- **Directory Support**: Directory exclusion support
- **Precedence**: Pattern precedence rules
- These exclusion requirements are mandatory and must be preserved

## Git Integration Requirements
- **Repository Detection**: Automatic Git repository detection
- **Archive Naming**: Git information in archive names
- **Info Extraction**: Git branch and commit hash extraction
- **State Detection**: Dirty working directory detection
- **Submodule Support**: Submodule handling
- **Configuration**: Git configuration integration
- These Git integration requirements are mandatory and must be preserved

## Archive Verification Requirements
- **Structure Check**: ZIP structure verification
- **Checksum Support**: SHA-256 checksum verification
- **Status Tracking**: Verification status tracking
- **Atomic Operations**: Atomic verification operations
- **Result Storage**: Verification result persistence
- **Error Reporting**: Error reporting for verification failures
- These verification requirements are mandatory and must be preserved

## Error Handling Requirements
- **Structured Errors**: All operations must return structured errors with status codes
- **Resource Cleanup**: No temporary files or directories may remain after any operation
- **Panic Recovery**: Application must recover from panics without leaving temporary resources
- **Context Support**: Long-running operations must support cancellation via context
- **Enhanced Detection**: Must detect various disk space and permission error conditions
- **Error Formatting**: Consistent error message formatting
- **Error Logging**: Comprehensive error logging
- These error handling requirements are mandatory and must be preserved

## Code Quality Standards
- **Linting**: All Go code must pass `revive` linter checks
- **Error Handling**: All errors must be properly handled (no unhandled return values)
- **Testing**: All code must have comprehensive test coverage
- **Documentation**: All public functions must be documented
- **Backward Compatibility**: New features must not break existing functionality
- **Code Organization**: Consistent code organization and structure
- **Naming Conventions**: Standard naming conventions
- These quality standards are immutable and must be maintained

## Build System Requirements
- **Quality Gates**: Build system must enforce quality standards before compilation
- **Dependency Management**: Proper ordering of build steps (lint → test → build)
- **Artifact Management**: Clean target for removing build artifacts
- **Continuous Integration**: Automated quality checks in CI/CD pipeline
- **Build Dependencies**: `make build` must depend on `make lint` and `make test` passing
- **Error Propagation**: Non-zero exit codes from linting or testing must prevent build
- **CI/CD Commands**: Support for `make ci-lint`, `make ci-test`, `make ci-build`
- These build system requirements are mandatory and must be preserved

## Output Formatting Requirements
- **Printf-Style Formatting**: All standard output must use printf-style format specifications
- **Template-Based Formatting**: Must support text/template and placeholder-based formatting
- **Configuration-Driven**: Format strings and templates must be retrieved from configuration
- **Text Highlighting**: Must provide means to highlight/format text for structure and meaning
- **Data Separation**: All user-facing text must be extracted from code into data files
- **Named Placeholders**: Must support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
- **Regex Integration**: Must support named regex groups for data extraction
- **Backward Compatibility**: Default format strings must preserve existing output appearance
- **Immutable Defaults**: Default format specifications cannot be changed without major version bump
- These output formatting requirements are mandatory and must be preserved

## Template Formatting Requirements
- **Dual Syntax Support**: Must support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
- **Regex Data Extraction**: Must extract named groups from filenames using configurable regex patterns
- **Graceful Degradation**: Must fall back to placeholder formatting when template processing fails
- **Operation Context**: Error messages must include operation context for enhanced debugging
- **Rich Data Display**: Must extract and display rich information from archive and backup filenames
- **ANSI Color Support**: Must support ANSI color codes and text formatting in templates
- **Template Methods**: Must provide template methods for all operations (archives, backups, config, errors)
- **Print Methods**: Must provide direct printing methods for template-formatted output
- **Error Handling**: Must handle invalid regex patterns and template syntax gracefully
- **Configuration Integration**: All template strings and regex patterns must be configurable
- These template formatting requirements are mandatory and must be preserved

## Commands
1. Create Archive:
   - Command: `bkpdir create`
   - Create full archive of current directory
   - **Must use atomic operations with automatic cleanup**
   - **Must support context cancellation**
   - **Output formatting must use configurable format strings**
   - This command structure must be preserved

2. Create Incremental Archive:
   - Command: `bkpdir create --incremental`
   - Create incremental archive of changes
   - **Must use atomic operations with automatic cleanup**
   - **Must support context cancellation**
   - **Output formatting must use configurable format strings**
   - This command structure must be preserved

3. Create File Backup:
   - Command: `bkpdir backup [FILE_PATH] [NOTE]`
   - Create backup of single file with comparison
   - **Must use atomic operations with automatic cleanup**
   - **Must support context cancellation**
   - **Must detect identical files and report appropriately**
   - **Output formatting must use configurable format strings**
   - This command structure must be preserved

4. List Archives:
   - Command: `bkpdir list`
   - Sort by creation time (most recent first)
   - Display format: `{path} (created: {time})`
   - Show verification status: [VERIFIED], [FAILED], or [UNVERIFIED]
   - **Output formatting must use configurable printf-style format strings**
   - This command structure and output format must be preserved

5. List File Backups:
   - Command: `bkpdir --list [FILE_PATH]`
   - Sort by creation time (most recent first)
   - Display format: `{path} (created: {time})`
   - **Output formatting must use configurable printf-style format strings**
   - This command structure and output format must be preserved

6. Verify Archive:
   - Command: `bkpdir verify [ARCHIVE_NAME]`
   - Support `--checksum` flag for checksum verification
   - Verify ZIP structure and integrity
   - Store verification results for display in list command
   - **Output formatting must use configurable printf-style format strings**
   - This verification behavior must remain unchanged

7. Display Configuration:
   - Command: `bkpdir config`
   - Display computed configuration values with name, value, and source
   - Process configuration files from `BKPDIR_CONFIG` environment variable
   - Exit after displaying values
   - **Output formatting must use configurable printf-style format strings**
   - This command behavior must remain unchanged once implemented

8. Backward Compatibility Commands:
   - Command: `bkpdir full [NOTE]` (alias for `bkpdir create`)
   - Command: `bkpdir inc [NOTE]` (alias for `bkpdir create --incremental`)
   - Command: `bkpdir --config` (alias for `bkpdir config`)
   - These backward compatibility commands must be preserved

## Configuration Defaults
- Configuration discovery uses `BKPDIR_CONFIG` environment variable to specify search path
- Default configuration search path is hard-coded as `./.bkpdir.yml:~/.bkpdir.yml` (if `BKPDIR_CONFIG` not set)
- Configuration files are processed in order with earlier files taking precedence
- Default archive directory: `../.bkpdir` relative to current directory
- Default backup directory: `../.bkpdir` relative to current directory
- Default use_current_dir_name: true
- Default use_current_dir_name_for_files: true
- Default include_git_info: true
- Default exclude_patterns: `[".git/", "vendor/"]`
- Default verification.verify_on_create: false
- Default verification.checksum_algorithm: "sha256"
- Default status codes: All status codes default to `0` (success) if not specified
  - `status_config_error`: 10
  - `status_created_archive`: 0
  - `status_disk_full`: 30
  - `status_failed_to_create_archive_directory`: 31
  - `status_directory_is_identical_to_existing_archive`: 0
  - `status_directory_not_found`: 20
  - `status_invalid_directory_type`: 21
  - `status_permission_denied`: 22
  - `status_created_backup`: 0
  - `status_failed_to_create_backup_directory`: 31
  - `status_file_is_identical_to_existing_backup`: 0
  - `status_file_not_found`: 20
  - `status_invalid_file_type`: 21
- **Default output format strings**: All format strings default to preserve existing output appearance
  - `format_created_archive`: "Created archive: %s\n"
  - `format_identical_archive`: "Directory is identical to existing archive: %s\n"
  - `format_list_archive`: "%s (created: %s)\n"
  - `format_config_value`: "%s: %s (source: %s)\n"
  - `format_dry_run_archive`: "Would create archive: %s\n"
  - `format_error`: "Error: %s\n"
  - `format_created_backup`: "Created backup: %s\n"
  - `format_identical_backup`: "File is identical to existing backup: %s\n"
  - `format_list_backup`: "%s (created: %s)\n"
  - `format_dry_run_backup`: "Would create backup: %s\n"
- **Default template format strings**: All template strings default to preserve existing output appearance
  - `template_created_archive`: "Created archive: %{path}\n"
  - `template_identical_archive`: "Directory is identical to existing archive: %{path}\n"
  - `template_list_archive`: "%{path} (created: %{creation_time})\n"
  - `template_config_value`: "%{name}: %{value} (source: %{source})\n"
  - `template_dry_run_archive`: "Would create archive: %{path}\n"
  - `template_error`: "Error: %{message}\n"
  - `template_created_backup`: "Created backup: %{path}\n"
  - `template_identical_backup`: "File is identical to existing backup: %{path}\n"
  - `template_list_backup`: "%{path} (created: %{creation_time})\n"
  - `template_dry_run_backup`: "Would create backup: %{path}\n"
- **Default regex patterns**: All regex patterns default to support data extraction
  - `pattern_archive_filename`: "(?P<prefix>[^-]*)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\\.zip"
  - `pattern_backup_filename`: "(?P<filename>[^/]+)-(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})-(?P<hour>\\d{2})-(?P<minute>\\d{2})(?:=(?P<note>.+))?"
  - `pattern_config_line`: "(?P<name>[^:]+):\\s*(?P<value>[^(]+)\\s*\\(source:\\s*(?P<source>[^)]+)\\)"
  - `pattern_timestamp`: "(?P<year>\\d{4})-(?P<month>\\d{2})-(?P<day>\\d{2})\\s+(?P<hour>\\d{2}):(?P<minute>\\d{2}):(?P<second>\\d{2})"
  - `template_identical_archive`: "Directory is identical to existing archive: %{path}\n"
  - `template_list_archive`: "%{path} (created: %{creation_time})\n"
  - `template_config_value`: "%{name}: %{value} (source: %{source})\n"
  - `template_dry_run_archive`: "Would create archive: %{path}\n"
  - `template_error`: "Error: %{message}\n"
  - `template_created_backup`: "Created backup: %{path}\n"
  - `template_identical_backup`: "File is identical to existing backup: %{path}\n"
  - `template_list_backup`: "%{path} (created: %{creation_time})\n"
  - `template_dry_run_backup`: "Would create backup: %{path}\n"
- **Default regex patterns**: Named regex patterns for data extraction
  - `pattern_archive_filename`: `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`
  - `pattern_backup_filename`: `(?P<filename>[^/]+)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<note>.+))?`
  - `pattern_config_line`: `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`
  - `pattern_timestamp`: `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`
- These configuration defaults must never be changed without explicit user override

## Platform Compatibility Requirements
- **OS Support**: Support macOS and Linux systems
- **File System**: Handle platform-specific file system differences
- **Permissions**: Preserve file permissions and ownership where applicable
- **Thread Safety**: Thread-safe operations for concurrent access
- **Resource Management**: Efficient resource management across platforms
- **Character Encoding**: Handle platform-specific character encodings
- **Line Endings**: Handle platform-specific line endings
- **Path Separators**: Handle platform-specific path separators
- Platform support must never be reduced or modified

## Global Options
- Support `--dry-run` flag for previewing archive operations
- **Dry-run must include resource cleanup verification**
- **Output formatting must use configurable printf-style format strings**
- Existing flag behavior must be maintained

## Build System Requirements
- **Linting**: `make lint` must pass before any code commit
- **Testing**: `make test` must pass with comprehensive coverage
- **Building**: `make build` must depend on successful linting and testing
- **Cleaning**: `make clean` must remove all build artifacts
- **Dependencies**: Proper dependency management
- **Version Info**: Version information in binaries
- **Binary Compatibility**: Maintain binary compatibility
- **Resource Embedding**: Proper resource embedding
- **Documentation**: Documentation generation
- **Release Packaging**: Proper release packaging
- These build requirements are immutable and must be enforced

## Resource Management Requirements
- **Automatic Cleanup**: All temporary resources must be cleaned up automatically
- **Thread Safety**: Resource management must be thread-safe
- **Atomic Operations**: File operations must use temporary files for atomicity
- **Leak Prevention**: No resource leaks allowed in any scenario
- **Error Resilience**: Cleanup must continue even if individual operations fail
- **Resource Tracking**: Comprehensive resource tracking
- **Memory Management**: Efficient memory management
- **Disk Space**: Efficient disk space management
- **Cleanup Warnings**: Proper cleanup warnings
- **Resource Limits**: Enforced resource limits
- These resource management requirements are mandatory and cannot be relaxed

## Performance Requirements
- **Minimal Overhead**: Resource tracking must have minimal performance impact
- **Efficient Operations**: File comparison must check length before byte comparison
- **Scalability**: Must handle large files and many archives efficiently
- **Memory Management**: Must maintain low memory footprint
- **CPU Usage**: Enforced CPU usage limits
- **Disk I/O**: Enforced disk I/O limits
- **Response Time**: Enforced response time limits
- **Cleanup Timing**: Efficient resource cleanup timing
- **Concurrency**: Enforced concurrent operation limits
- **Verification**: Efficient verification speed
- These performance characteristics must be preserved

## Feature Preservation Rules
1. New Features:
   - Must not interfere with existing functionality
   - Must maintain all current behaviors
   - Must be optional and not affect existing workflows
   - **Must include automatic resource cleanup**
   - **Must support context cancellation where appropriate**
   - **Must pass all linting and testing requirements**

2. Modifications:
   - Must preserve all existing command-line interfaces
   - Must maintain current directory handling behaviors
   - Must keep existing configuration options
   - Must not change established archive naming patterns
   - **Must not introduce resource leaks**
   - **Must maintain error handling standards**
   - **Must preserve atomic operation guarantees**
   - **Must maintain Git integration compatibility**
   - **Must preserve verification functionality**

3. Testing Requirements:
   - All new code must include tests for existing functionality
   - Regression tests must verify no existing features are broken
   - Platform compatibility tests must be maintained
   - **Resource cleanup must be verified in all test scenarios**
   - **Context cancellation and timeout handling must be tested**
   - **Performance benchmarks must not regress**
   - **All code must pass linting before commit**
   - **Git integration tests must be maintained**
   - **Archive verification tests must be comprehensive**

4. Quality Assurance:
   - **Code must pass revive linter with zero warnings**
   - **All errors must be properly handled**
   - **All public functions must be documented**
   - **Test coverage must meet minimum thresholds**
   - **No temporary files may remain after any operation**
   - **Memory leaks are strictly prohibited**
   - **Archive integrity must be guaranteed**
   - **Git integration must be reliable**

## Testing Infrastructure Immutable Requirements

### Archive Corruption Testing Framework (TEST-INFRA-001-A)

#### Corruption Type Enumeration Stability
- **Immutable**: The 8 corruption types (CRC, Header, Truncate, CentralDir, LocalHeader, Data, Signature, Comment) must remain stable across versions
- **Rationale**: Test code depends on these specific corruption types for systematic verification testing
- **Impact**: New corruption types may be added but existing types cannot be renamed or removed
- **Version**: Established in initial implementation, immutable from v1.0.0

#### Deterministic Corruption Behavior  
- **Immutable**: Identical seeds must produce identical corruption across versions
- **Rationale**: CI/CD systems depend on reproducible test results for regression detection
- **Impact**: Corruption algorithms cannot be changed in ways that break reproducibility
- **Version**: Established in initial implementation, immutable from v1.0.0
- **Exception**: Bug fixes that improve correctness may change behavior if documented in release notes

#### Safe Testing Guarantees
- **Immutable**: Backup/restore functionality must prevent data loss during corruption testing
- **Rationale**: Testing infrastructure must never risk data integrity
- **Impact**: Any changes to backup/restore logic must maintain safety guarantees
- **Version**: Established in initial implementation, immutable from v1.0.0

#### Performance Baseline Stability
- **Immutable**: Performance characteristics must not degrade by more than 20% across versions
- **Rationale**: Testing infrastructure performance affects overall test suite execution time
- **Baseline**: CRC corruption ~763μs, detection ~49μs (±20% acceptable variance)
- **Impact**: Performance optimizations welcome, but regressions require justification
- **Version**: Established in initial implementation, immutable from v1.0.0 