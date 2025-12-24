# Changelog

All notable changes to this project will be documented in this file.

## [1.7.1] - 2025-12-24

### Fixed
- **Configuration Hierarchy Preservation**: Fixed bug where local config files would stop reading values from home directory config files
  - When a local config file (`.bkpdir.yml`) exists but doesn't set a field (e.g., `archive_dir_path`), the system now correctly preserves values from home config files (`~/.bkpdir.yml`) instead of falling back to compiled defaults
  - Compiled defaults now only apply when a value is not set anywhere in the configuration hierarchy
  - Ensures proper configuration hierarchy: home config → local config → compiled defaults
  - Updated `mergeBasicSettings` and `applyOverride` functions to correctly preserve values from earlier files when later files don't set those fields
  - Added comprehensive test case to validate the fix

## [1.7.0] - 2025-12-13

### Added
- **Diff Command** (`bkpdir diff`): New CLI command that reports changes between the current directory and the reconstructed archive state (most recent incremental applied on top of most recent full archive)
  - Shows added, modified, and deleted files
  - Handles edge cases: no archives, no incremental archive, no full archive
  - Respects exclude patterns from configuration
  - Supports context cancellation for long-running operations
  - Uses configurable format strings for output
  - Supports dry-run mode for preview
  - Comprehensive test coverage with integration tests

- **Incremental Archive Duplicate Prevention**: Prevents creation of duplicate incremental archives when no observable changes exist
  - Compares current directory state against the union of most recent incremental + most recent full archive
  - Skips archive creation with appropriate message when no changes detected
  - Reuses diff command analysis for consistency
  - Handles edge cases: no incremental exists (compares against full), no archives exist (errors appropriately)
  - Respects exclude patterns from configuration
  - Behavior is consistent with diff command output
  - Comprehensive test coverage with integration tests

- List command output limit feature: `--limit` flag (short: `-n`) to control the number of items displayed
  - Default limit of 10 items for both `list` command (directory archives) and `--list` command (file backups)
  - Option `--limit 0` to show all items
  - Maintains sorting by creation time (most recent first)
  - Comprehensive test coverage with edge cases

### Changed
- **Archive Selection Strategy**: Changed archive selection to use name-based sorting instead of file modification times
  - Archive names include timestamps in ISO 8601 format (`YYYY-MM-DDTHHmmss`) which are alphabetically sortable
  - More reliable than file system modification times
  - Consistent with archive naming conventions
  - Simpler implementation (no file system stat calls needed)
  - Affects: `diff` command, incremental archive duplicate prevention, and archive selection functions

## [1.6.0] - 2025-12-12

### Added
- User-customizable format strings for output formatting
- Diagnostic output control functionality
- Configuration output grouping feature
- Comprehensive configuration merge tests

### Fixed
- Fixed precedence handling for configuration files in merge operations

### Changed
- Removed archive verification feature
- Updated to STDD (Semantic Token-Driven Development) methodology v1.0.1

## [1.5.1] - 2025-11-19

### Fixed
- Fixed date formatting in `bkpdir list` command output (was showing `%!s(MISSING)`).

## [1.5.0] - 2025-06-07

### Added
- Initial release of version 1.5.0.
