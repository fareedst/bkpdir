# Changelog

All notable changes to this project will be documented in this file.

## [1.6.1pre1212] - 2025-12-12

### Added
- List command output limit feature: `--limit` flag (short: `-n`) to control the number of items displayed
  - Default limit of 10 items for both `list` command (directory archives) and `--list` command (file backups)
  - Option `--limit 0` to show all items
  - Maintains sorting by creation time (most recent first)
  - Comprehensive test coverage with edge cases

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
