# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.7.2] - 2026-03-29

### Changed (behavior)

- **Default `exclude_patterns`**: Built-in default is now an **empty** list. `.git/`, `vendor/`, and other paths are no longer excluded automatically; add them under `exclude_patterns` in `.bkpdir.yml` (or use `bkpdir template`) when you want those directories skipped.

### Fixed
- **Directory exclusion patterns** (`[IMPL-EXCLUSION_PATTERNS]`): single-segment trailing-slash patterns (e.g. `node_modules/`) match that directory name at any depth.
- **Config merge** (`[IMPL-CFG_PRECEDENCE_FIX]` / CFG-005): the first loaded config file’s unprefixed `exclude_patterns` replaces built-in defaults instead of merging into them.
- **Config load**: stderr warnings and a YAML glob/alias parsing hint when a config file fails to load.

### Changed
- **Makefile**: added `validate-tied-mcp`; `lint` runs MCP validation before `revive`; `check` is fmt + vet + lint only; `validate-token-enforcement` (DOC-008) remains optional.
- **TIED**: removed invalid nested `detail_file` from REQ/ARCH detail YAML; widespread IMPL decision YAML normalization and pseudocode alignment; `tied/semantic-tokens.yaml` and agent implementation checklist docs updated.
- **`project-tokens.yaml`** and **`semantic-token-validation-report.md`**: registry and report refreshed.

### Removed
- Tracked prebuilt platform binaries under `bin/` (macOS amd64/arm64, Ubuntu 20.04/22.04/24.04).

## [2.0.0] - 2026-03-12

### Changed
- **STDD to TIED migration**: Converted entire project from STDD monolithic markdown files (v1.0.2) to TIED YAML-driven methodology (v2.2.0) with MCP server support
- **Token format**: All semantic tokens switched from colon delimiter (`[REQ:TOKEN]`) to hyphen delimiter (`[REQ-TOKEN]`) across 238 files (149 Go source/test files, 89 non-Go files), totaling 1,026 line changes
- **Semantic token validator**: Updated `internal/validation/semantic_token_validator.go` to detect both colon and hyphen token formats for backward compatibility
- **Shell scripts**: Updated `scripts/validate-semantic-tokens.sh`, `scripts/token-coverage-analysis.sh`, and `scripts/token-navigate.sh` to support both token formats and reference `tied/` paths instead of `stdd/`
- **Legacy comment replacement**: Replaced ad-hoc comment styles (`ARCH-001:`, `REFACTOR-*:`, `SERVICE-*-001:`, `IMMUTABLE-REF:`, etc.) with standardized TIED semantic token annotations (`[IMPL-*]`, `[ARCH-*]`, `[REQ-*]`) across all Go source and test files (~68 files); no functional code or test logic changes

### Added
- **TIED YAML database**: 177 detail YAML files with matching index records
  - 59 REQ records (56 project-specific + 3 inherited methodology tokens)
  - 46 ARCH records (43 project-specific + 3 inherited)
  - 72 IMPL records (69 project-specific + 3 inherited)
- **Semantic tokens registry**: 147 tokens registered in `tied/semantic-tokens.yaml`
- **IMPL pseudo-code**: All 72 IMPL detail files enriched with `essence_pseudocode` containing semantic token comments per `[PROC-IMPL_PSEUDOCODE_TOKENS]`
- **Conversion scripts** (Ruby 3.0.3):
  - `scripts/convert_tokens.rb` — converts colon-format tokens to hyphen-format in Go files
  - `scripts/convert_tokens_nongo.rb` — same for non-Go files (shell, markdown, etc.)
  - `scripts/generate_pseudocode.rb` — generates `essence_pseudocode` for IMPL detail files
  - `scripts/register_semantic_tokens.rb` — registers tokens in `semantic-tokens.yaml` from indexes
  - `scripts/fix_stdd_refs.rb` — cleans stale `stdd/` path references in `tied/` markdown files
- **IMPL-to-code-and-tests linkage**: All ~50 IMPL detail YAML files enriched with detailed `essence_pseudocode`, `code_locations`, and `traceability.tests` per `[PROC-IMPL_CODE_TEST_SYNC]`
- **Process `[PROC-IMPL_CODE_TEST_SYNC]`**: New 9-phase IMPL-to-Code-and-Tests linkage workflow in `tied/processes.md`
- **Linkage guide**: New `tied/impl-code-test-linkage.md` — practical reference for the three-way alignment principle (pseudo-code, tests, code)

### Removed
- **STDD directory**: `stdd/` monolithic markdown files replaced by structured `tied/` YAML files

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
