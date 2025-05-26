# Immutable Specifications

This document contains specifications that MUST NOT be changed without a major version bump. These are core behaviors that users and other systems depend on.

## Archive Naming Convention
- Format: `[PREFIX-]YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip`
- PREFIX is the current directory name (if `use_current_dir_name` is true)
- BRANCH and HASH are Git information (if in a Git repository and `include_git_info` is true)
- Optional note can be appended with equals sign
- Incremental archives: `BASENAME_update=YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip`
- This naming convention is fixed and must not be modified

## Directory Operations
- Use platform-independent path handling
- Preserve file permissions and modification time in archives
- Handle both absolute and relative directory paths
- Source paths must be directories (not files or special files)
- Create archive directories automatically if they don't exist
- Display all paths relative to current directory
- **Atomic Operations**: All archive operations must be atomic to prevent corruption
- **Resource Cleanup**: All temporary files must be cleaned up automatically
- **ZIP Format**: Archives must use ZIP format with compression
- These directory operation rules are fundamental and must not be altered

## File Exclusion Requirements
- **Glob Pattern Matching**: Must use doublestar glob pattern matching
- **Default Exclusions**: Default exclude patterns are `[".git/", "vendor/"]`
- **Configurable Patterns**: Exclude patterns must be configurable via `exclude_patterns`
- **Recursive Application**: Patterns must be applied recursively to directory tree
- These exclusion requirements are mandatory and must be preserved

## Git Integration Requirements
- **Repository Detection**: Must automatically detect Git repositories
- **Branch Information**: Must extract current branch name when available
- **Commit Hash**: Must extract current commit hash when available
- **Non-Git Handling**: Must gracefully handle non-Git directories
- **Configurable**: Git integration must be configurable via `include_git_info`
- **Archive Naming**: Git information must be included in archive names when enabled
- These Git integration requirements are immutable and must be maintained

## Archive Verification Requirements
- **ZIP Structure**: Must verify ZIP file structure and integrity
- **Checksum Support**: Must support SHA-256 checksum verification
- **Status Tracking**: Must track and store verification status
- **Optional Verification**: Verification must be optional and configurable
- **Manual Verification**: Must support manual verification of existing archives
- These verification requirements are mandatory and must be preserved

## Error Handling Requirements
- **Structured Errors**: All archive operations must return structured errors with status codes
- **No Resource Leaks**: No temporary files or directories may remain after any operation
- **Panic Recovery**: Application must recover from panics without leaving temporary resources
- **Context Support**: Long-running operations must support cancellation via context
- **Enhanced Detection**: Must detect various disk space and permission error conditions
- These error handling requirements are mandatory and must be preserved

## Code Quality Standards
- **Linting**: All Go code must pass `revive` linter checks
- **Error Handling**: All errors must be properly handled (no unhandled return values)
- **Testing**: All code must have comprehensive test coverage
- **Documentation**: All public functions must be documented
- **Backward Compatibility**: New features must not break existing functionality
- These quality standards are immutable and must be maintained

## Output Formatting Requirements
- **Printf-Style Formatting**: All standard output must use printf-style format specifications
- **Template-Based Formatting**: Must support text/template and placeholder-based formatting for named data extraction
- **Configuration-Driven**: Format strings and templates must be retrieved from application configuration
- **Text Highlighting**: Must provide means to highlight/format text for structure and meaning
- **Data Separation**: All user-facing text must be extracted from code into data files
- **Named Placeholders**: Must support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
- **Regex Integration**: Must support named regex groups for data extraction and template formatting
- **Backward Compatibility**: Default format strings must preserve existing output appearance
- **Immutable Defaults**: Default format specifications cannot be changed without major version bump
- These output formatting requirements are mandatory and must be preserved

## Commands
1. Create Full Archive:
   - Command: `bkpdir full [NOTE]`
   - Compare with most recent archive before creating
   - Skip if identical to most recent archive
   - **Must use atomic operations with automatic cleanup**
   - **Must support context cancellation**
   - **Output formatting must use configurable printf-style format strings**
   - This archive creation logic must remain unchanged

2. Create Incremental Archive:
   - Command: `bkpdir inc [NOTE]`
   - Requires existing full archive as base
   - Only includes files modified since base archive
   - **Must use atomic operations with automatic cleanup**
   - **Must support context cancellation**
   - This incremental archive logic must remain unchanged

3. List Archives:
   - Command: `bkpdir list`
   - Sort by creation time (most recent first)
   - Display format: `.bkpdir/project-2024-03-21-15-30=main=abc123=note.zip (created: 2024-03-21 15:30:00)`
   - Show verification status: [VERIFIED], [FAILED], or [UNVERIFIED]
   - **Output formatting must use configurable printf-style format strings**
   - This command structure and output format must be preserved

4. Verify Archive:
   - Command: `bkpdir verify [ARCHIVE_NAME]`
   - Support `--checksum` flag for checksum verification
   - Verify ZIP structure and integrity
   - Store verification results for display in list command
   - **Output formatting must use configurable printf-style format strings**
   - This verification behavior must remain unchanged

5. Display Configuration:
   - Command: `bkpdir --config`
   - Display computed configuration values with name, value, and source
   - Process configuration files from `BKPDIR_CONFIG` environment variable
   - Exit after displaying values
   - **Output formatting must use configurable printf-style format strings**
   - This command behavior must remain unchanged once implemented

## Configuration Defaults
- Configuration discovery uses `BKPDIR_CONFIG` environment variable to specify search path
- Default configuration search path is hard-coded as `./.bkpdir.yml:~/.bkpdir.yml` (if `BKPDIR_CONFIG` not set)
- Configuration files are processed in order with earlier files taking precedence
- Default archive directory: `../.bkpdir` relative to current directory
- Default use_current_dir_name: true
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
- **Default output format strings**: All format strings default to preserve existing output appearance
  - `format_created_archive`: "Created archive: %s\n"
  - `format_identical_archive`: "Directory is identical to existing archive: %s\n"
  - `format_list_archive`: "%s (created: %s)\n"
  - `format_config_value`: "%s: %s (source: %s)\n"
  - `format_dry_run_archive`: "Would create archive: %s\n"
  - `format_error`: "Error: %s\n"
- **Default template format strings**: Template-based formatting with named placeholders
  - `template_created_archive`: "Created archive: %{path}\n"
  - `template_identical_archive`: "Directory is identical to existing archive: %{path}\n"
  - `template_list_archive`: "%{path} (created: %{creation_time})\n"
  - `template_config_value`: "%{name}: %{value} (source: %{source})\n"
  - `template_dry_run_archive`: "Would create archive: %{path}\n"
  - `template_error`: "Error: %{message}\n"
- **Default regex patterns**: Named regex patterns for data extraction
  - `pattern_archive_filename`: `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`
  - `pattern_config_line`: `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`
  - `pattern_timestamp`: `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`
- These configuration defaults must never be changed without explicit user override

## Platform Compatibility
- Support macOS and Linux systems
- Handle platform-specific file system differences
- Preserve file permissions and ownership where applicable
- **Thread-safe operations for concurrent access**
- **Efficient resource management across platforms**
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
- These build requirements are immutable and must be enforced

## Resource Management Requirements
- **Automatic Cleanup**: All temporary resources must be cleaned up automatically
- **Thread Safety**: Resource management must be thread-safe
- **Atomic Operations**: Archive operations must use temporary files for atomicity
- **Leak Prevention**: No resource leaks allowed in any scenario
- **Error Resilience**: Cleanup must continue even if individual operations fail
- These resource management requirements are mandatory and cannot be relaxed

## Performance Requirements
- **Minimal Overhead**: Resource tracking must have minimal performance impact
- **Efficient Operations**: Directory comparison must use early termination for differences
- **Scalability**: Must handle large directories and many archives efficiently
- **Memory Management**: Must maintain low memory footprint for large directories
- **Streaming Operations**: Must use streaming ZIP creation for memory efficiency
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