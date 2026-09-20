// [REQ-DOC_016]
// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
package main

// SPEC-ID: IMPL-TOKEN_SYSTEM::STRUCTURED_EXTRACTION
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: parse legacy token YAML groups and emit normalized JSON with descriptions, status, and source paths per category.
// SPEC-ID: IMPL-TOKEN_SYSTEM::CROSS_LINK_SYNTHESIS
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-GOV_REGISTRY_COMPLETENESS] [REQ-IMMUTABLE_DIRECTORY_OPERATIONS] — How: resolve canonical REQ/ARCH/IMPL references for each legacy token and flag gaps when links are missing.

// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: validate cwd, collect files, generate full archive name, and delegate atomic zip write unless dry-run.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: write zip to path.tmp, rename to final path, untrack temp file, print created stats.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: create zip writer with Deflate and add each collected file with context cancellation checks.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: write one entry with Deflate header, copying file content or symlink target and honoring skip_broken_symlinks.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: route to incremental or full archive name builder based on config flags.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: build {prefix}-{timestamp}[={branch}={hash}[-dirty]][={note}].zip from config git segments.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: build {base}_update={timestamp}[={branch}={hash}[-dirty]][={note}].zip for incremental archives.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: walk cwd tree, skip excluded paths and directories, collect relative file paths with cancellation checks.
// - [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: read archive directory entries and return Archive metadata for each .zip file.

// - [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] — How: scan each module file for REQ/ARCH/IMPL comments and record missing expected tokens.

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: format Message alone or Message WITH underlying Err when present.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: return wrapped Err for errors.Is/As chains.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: mirror ArchiveError Error formatting for backup operations.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message and status code only.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status code, and underlying Err.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status, operation, path, and optional Err.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message and status code only.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status code, and underlying Err.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status, operation, path, and optional Err.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match disk-space substrings in error text OR path error with OS no-space/quota/large-file codes.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match permission-denied substrings OR path error with EACCES/EPERM.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match directory-not-found substrings OR path error with ENOENT.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: route ApplicationError or classified errors to formatter and configured status codes.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: map error to ErrorCategory via disk, permission, filesystem, network detectors.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify ArchiveError.Error with and without underlying cause.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify IsDiskFullError matches message patterns and syscall path errors.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify IsPermissionError matches permission strings and EACCES/EPERM.
// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify HandleError returns configured status codes for classified errors.

// - [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: for each package in a phase move code to pkg/, define interfaces first, update imports, and require go build plus all tests green.
// - [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: expose type aliases and delegating wrappers at root with deprecation comments pointing to new pkg paths.

// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: register --limit/-n persistent flag with default 10 so list subcommand and --list share one limit variable.
// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: load config and formatter then pass listLimit to ListArchivesEnhanced for list subcommand.
// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] [REQ-OUTPUT_FORMATTING] — How: sort archives most-recent-first then truncate to limit when limit > 0 before formatted output.
// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: resolve file path from --list or args and pass listLimit to ListFileBackupsEnhanced.
// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] [REQ-OUTPUT_FORMATTING] — How: sort backups most-recent-first then truncate to limit when limit > 0 before formatted output.
// - [IMPL-LIST_LIMIT] [ARCH-LIST_LIMIT] [REQ-LIST_LIMIT] — How: apply hardcoded default limit of 10 for CommandHandler archive and file-backup listing paths.

// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: route template placeholders before printf and return plain format when no verbs.
// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: mirror FORMAT_LIST_ARCHIVE guard pattern using cfg.FormatListBackup.
// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: priority-based format selection with extraction data and guarded printf.
// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: simplified adapter path always gathers stats and selects format by priority.
// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: verify template-style FormatListArchiveWithExtraction replaces #{size_human} and includes path.
// - [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: verify printf-style FormatListArchiveWithExtraction matches FORMAT_WITH_PLACEHOLDERS output without EXTRA args.

// - [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING] — How: reconstruct archive state, CalculateDiff against cwd, skip creation and print skip message when diff has no added/modified/deleted entries.
// - [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING] — How: format cfg.FormatIncrementalSkippedNoChanges and emit via delayed collector or stdout.

// - [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: append -dirty to full archive name only when repo is dirty and ShowGitDirtyStatus is enabled.
// - [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: apply same conditional -dirty suffix on incremental update archive names.
// - [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: copy cfg.ShowGitDirtyStatus into ArchiveConfig before delegating to name builder with git info when enabled.
// - [IMPL-GIT_DIRTY_CONFIG] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION] — How: merge show_git_dirty_status in mergeBasicSettings and mergeGitSettings respecting CFG-001 explicit-set precedence.

// - [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: GatherFileStatInfo, build data map, replace #{key} in TemplateCreatedArchiveDetailed, fall back on error or leftover placeholders.
// - [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: same as created variant using TemplateIncrementalCreatedDetailed and FormatIncrementalCreated fallback.
// - [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: format with stats then Print via collector or stdout.
// - [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: format incremental with stats then Print via collector or stdout.

// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: stat path and populate FileStatInfo with name, size, human size, mtime, mode, and type.
// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: scale bytes to TB/GB/MB/KB/B with one decimal for large units.
// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: classify regular, directory, symlink, device, pipe, socket, or other from mode bits.

// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: sprintf cfg.FormatCreatedArchive with path for simple archive-created messages.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: if format contains #{ use formatTemplate with GatherFileStatInfo data; elif contains % use sprintf; else return literal format string.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: sprintf cfg.FormatConfigValue with name, value, and source fields.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: declare FormatProvider, OutputDestination, PatternExtractor, FormatterInterface, and TemplateFormatterInterface contracts for extraction and delayed output.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: TemplateFormatter.FormatWithPlaceholders replaces #{key} from data with defaults for missing stat fields.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: compile regex pattern, extract named submatches into data map, delegate to FormatWithPlaceholders on template string.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: compile configured regex and return map of named capture groups from input text.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: replace every #{key} in formatStr with data map values; leave unknown placeholders unchanged.
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: Print* methods format message then AddStdout when collector attached else immediate print.

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateDirectorySnapshot with exclusion patterns applied during walk.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateArchiveSnapshot to enumerate zip member files with metadata.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate equality check of two DirectorySnapshot values to fileops.CompareSnapshots.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: build dir and archive snapshots with exclusions and compare for full structural equality.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: list archives, keep full archives only, sort by name, return path of last entry.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: FindMostRecentArchive then IsDirectoryIdenticalToArchive against cwd.
// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: build human-readable file listing from CreateArchiveSnapshot for diagnostics and tests.

// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: list archives, filter non-incremental, sort by name ascending, return latest full Archive.
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: select incrementals whose name prefix matches base full archive _update= pattern; return latest by name.
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: load full zip snapshot, overlay incremental zip files by RelativePath, return merged DirectorySnapshot.
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: snapshot cwd, compare file maps by path/size/hash; classify added, modified, deleted (files only).
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-CONTEXT_SUPPORT] [REQ-DIFF_COMMAND] — How: load config, reconstruct state, handle no-archive gracefully, CalculateDiff, PrintDiffResult with context cancellation checks.
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-OUTPUT_FORMATTING] [REQ-DIFF_COMMAND] — How: use cfg FormatDiff* strings for no-changes header and per-category file lines.
// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-OUTPUT_FORMATTING] [REQ-DIFF_COMMAND] — How: format diff then route through delayed collector or stdout.

// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return an empty OutputCollector ready to append messages.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: append OutputMessage with Destination stdout and given message type.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: append OutputMessage with Destination stderr and given message type.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return a snapshot of all buffered messages without flushing.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write every message to stdout or stderr then clear the buffer.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write only stdout messages and retain stderr entries in the buffer.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: write only stderr messages and retain stdout entries in the buffer.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: discard all buffered messages without printing.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: report whether OutputFormatter has a non-nil collector attached.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: return the attached collector pointer for tests and flush orchestration.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: attach or detach delayed-mode collector (nil disables buffering).
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: when collector present, AddStdout/AddStderr instead of immediate fmt.Print for formatted messages.
// - [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: construct AIFormatterAdapter with cfg and pre-wired OutputCollector for delayed CLI output.

// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: hold naming inputs (prefix, timestamp, git segments, note, incremental base) for archive filename generation.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: represent a discovered zip archive with path, creation time, incremental flag, git metadata, and base archive link.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract Config field accessors needed by archive creation without importing main.Config in tests.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: bundle context, cwd, target path, file list, config interface, and resource manager for create-archive entry points.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract dry-run and incremental print operations for archive workflows.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *Config and delegate each ArchiveConfigInterface getter to the matching Config field.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: type-assert OutputFormatterInterface to FormatterAdapter or AIFormatterAdapter for extended archive print methods.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: group Config, note, dry-run flag, and context for incremental archive API entry.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract backup directory path, naming toggle, and backup-related status codes.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract backup dry-run, identical, created, and not-found output methods.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *Config and delegate BackupConfigInterface getters to Config fields.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *OutputFormatter and forward each BackupFormatterInterface call.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: capture list/compare metadata (name, path, creation time, size) for backup discovery.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: represent one backup file with source path and optional note segment.
// - [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: bundle context, config, formatter, file path, note, and dry-run for createFileBackupInternal.

// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: compare extracted placeholders to getExpectedPlaceholders and emit non-fatal warnings for unknown tokens.
// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: return allowed printf or template placeholders per Config field name from placeholderMap.
// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: regex-collect %[sdvfbtxX] verbs and #{name} template tokens from a format string.
// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: iterate all Format* and Template* cfg fields and aggregate ValidateFormatString warnings.
// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: replace #{key} from data, apply known defaults, handle residual %s with path/time, strip remaining #{...} before optional Go template pass.
// - [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: after YAML merge print validateAllFormatStrings warnings to stderr without failing load.

// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: define aggregate holding archive, backup, status, format, template, pattern, git, and inheritance settings for serialization and reflection.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: represent one displayed configuration entry with name, string value, and source label.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: capture reflection metadata for one config field including YAML name, type, path, category, and importance.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: return new Config populated with preset defaults for all fields including nested git config.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: expose subset of key fields as ConfigValue rows for simple --config display.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] [REQ-CFG_006] — How: delegate to reflection-based GetAllConfigValuesWithSources and convert to legacy ConfigValue format sorted by name.
// - [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: verify DefaultConfig preset values for archive path, flags, and exclude patterns.

// - [IMPL-CONFIG_SCHEMA_FLEX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: clone defaultConfig, merge each existing search-path YAML file, apply environment overrides, then validate schema.
// - [IMPL-CONFIG_SCHEMA_FLEX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: delegate to GenericConfigLoader with DefaultConfig() then type-assert to *Config or fall back to defaults.

// - [IMPL-CONFIG_OUTPUT_GROUPING] [ARCH-CONFIG_OUTPUT_GROUPING] [REQ-CONFIG_OUTPUT_GROUPING] — How: map archive and backup dir paths to critical, key toggles to high, format/template/status names to low, otherwise medium.
// - [IMPL-CONFIG_OUTPUT_GROUPING] [ARCH-CONFIG_OUTPUT_GROUPING] [REQ-CONFIG_OUTPUT_GROUPING] — How: bucket values by category, sort categories by CategoryPriority, sort fields by importance then name, print sources header and section headers.

// - [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: return cached field metadata when schema hash matches; on miss call reflectConfigFields (IMPL-CONFIG_DISPLAY_FLATTENING), sort, strip values, cache metadata.
// - [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: copy cached metadata and populate Value via getFieldValueByPath or zero value on failure.
// - [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: split dot path, dereference pointers, walk struct fields by name, return leaf interface.
// - [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: format nil, bool, numeric, string, string-slice, pointer, and default kinds for config display.
// - [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: infer append/prepend/merge/override/default from field path patterns and value comparison to old value.
// - [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: map field shape flags and name hints to append, merge, prepend, or override.
// - [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: RWMutex-backed cache invalidated when Config struct type hash changes.
// - [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: table-driven test asserts determineMergeStrategyForField paths map to expected strategies.

// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return true only when mode is a regular file.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return true only when entry is a directory.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: STAT path and return specific errors for not-exist, permission denied, or other access failures.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: require path arg, validate accessibility, route file to file backup handler else directory archive else unsupported type exit.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: load config from cwd, build formatter, resolve note from global flag or second arg, invoke enhanced file backup with dry-run.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: chdir to target directory, load config from dot, resolve note, run full archive with context and restore cwd on exit.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: build production root with persistent flags and subcommands; root Run dispatches config/list flags or auto-detect when positional args present.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: bypass Cobra when first token is not known command/flag; parse global dry-run and note from argv; else SetArgs and Execute.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: integration test validates validatePath, isFile, and isDirectory for file, directory, missing, and absolute paths.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: explicit backup, create, and full subcommands still execute via Cobra without auto-detect path routing.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: production newRootCommand registers fixed subcommand set matching composition contract.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: first token template routes through Cobra Execute and emits help output.
// - [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] — How: path-first argv with dry-run exercises auto-detect file backup without writing archives.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestNewTokenAnalyzer(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created, got nil")
	}

	if analyzer.config == nil {
		t.Error("Expected config to be initialized")
	}

	if analyzer.fileSet == nil {
		t.Error("Expected fileSet to be initialized")
	}

	if analyzer.featureMap == nil {
		t.Error("Expected featureMap to be initialized")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestNewTokenValidator(t *testing.T) {
	validator := NewTokenValidator()

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}

	if validator.config == nil {
		t.Error("Expected config to be initialized")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestNewBatchProcessor(t *testing.T) {
	processor := NewBatchProcessor()

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}

	if processor.validator == nil {
		t.Error("Expected validator to be initialized")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestDeterminePriority(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name           string
		functionName   string
		context        map[string]string
		expectedIcon   string
		expectedReason string
	}{
		{
			name:           "Critical function - main",
			functionName:   "MainFunction",
			context:        map[string]string{},
			expectedIcon:   "[CRITICAL]",
			expectedReason: "Critical operation: main",
		},
		{
			name:           "High priority - config",
			functionName:   "LoadConfig",
			context:        map[string]string{},
			expectedIcon:   "[HIGH]",
			expectedReason: "High priority: config",
		},
		{
			name:           "Medium priority - format",
			functionName:   "FormatOutput",
			context:        map[string]string{},
			expectedIcon:   "[MEDIUM]",
			expectedReason: "Medium priority: format",
		},
		{
			name:           "Context-based high priority - error handling",
			functionName:   "ProcessData",
			context:        map[string]string{"error_handling": "true"},
			expectedIcon:   "[HIGH]",
			expectedReason: "High priority: error handling function",
		},
		{
			name:           "Default low priority",
			functionName:   "UtilityFunction",
			context:        map[string]string{},
			expectedIcon:   "[LOW]",
			expectedReason: "Low priority: utility function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			icon, reason := analyzer.determinePriority(signature, tt.context)

			if icon != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, icon)
			}

			if !strings.Contains(reason, strings.Split(tt.expectedReason, ":")[0]) {
				t.Errorf("Expected reason to contain %s, got %s", tt.expectedReason, reason)
			}
		})
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestDetermineAction(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name           string
		functionName   string
		context        map[string]string
		expectedIcon   string
		expectedReason string
	}{
		{
			name:           "Analysis function - get",
			functionName:   "GetConfig",
			context:        map[string]string{},
			expectedIcon:   "[DECISION:discovery]",
			expectedReason: "Analysis operation: get",
		},
		{
			name:           "Documentation function - format",
			functionName:   "FormatOutput",
			context:        map[string]string{},
			expectedIcon:   "[DECISION:format-processing]",
			expectedReason: "Documentation operation: format",
		},
		{
			name:           "Configuration function - create",
			functionName:   "CreateArchive",
			context:        map[string]string{},
			expectedIcon:   "[DECISION:core-functionality]",
			expectedReason: "Configuration operation: create",
		},
		{
			name:           "Protection function - validate",
			functionName:   "ValidateInput",
			context:        map[string]string{},
			expectedIcon:   "[DECISION:validation]",
			expectedReason: "Protection operation: validate",
		},
		{
			name:           "Context-based protection - error handling",
			functionName:   "ProcessRequest",
			context:        map[string]string{"error_handling": "true"},
			expectedIcon:   "[DECISION:validation]",
			expectedReason: "Protection: error handling",
		},
		{
			name:           "Default configuration",
			functionName:   "UtilityFunction",
			context:        map[string]string{},
			expectedIcon:   "[DECISION:core-functionality]",
			expectedReason: "Configuration: general operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			icon, reason := analyzer.determineAction(signature, tt.context)

			if icon != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, icon)
			}

			if !strings.Contains(reason, strings.Split(tt.expectedReason, ":")[0]) {
				t.Errorf("Expected reason to contain %s, got %s", tt.expectedReason, reason)
			}
		})
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestDetermineFeatureID(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name         string
		functionName string
		context      map[string]string
		expected     string
	}{
		{
			name:         "Config function",
			functionName: "LoadConfig",
			context:      map[string]string{},
			expected:     "CFG-NEW",
		},
		{
			name:         "Archive function",
			functionName: "CreateArchive",
			context:      map[string]string{},
			expected:     "ARCH-NEW",
		},
		{
			name:         "Backup function",
			functionName: "BackupFile",
			context:      map[string]string{},
			expected:     "FILE-NEW",
		},
		{
			name:         "Git function",
			functionName: "GitCommit",
			context:      map[string]string{},
			expected:     "GIT-NEW",
		},
		{
			name:         "Test function",
			functionName: "TestValidation",
			context:      map[string]string{},
			expected:     "TEST-NEW",
		},
		{
			name:         "Generic function",
			functionName: "UtilityFunction",
			context:      map[string]string{},
			expected:     "UTIL-NEW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			featureID := analyzer.determineFeatureID(signature, tt.context)

			if featureID != tt.expected {
				t.Errorf("Expected feature ID %s, got %s", tt.expected, featureID)
			}
		})
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestCalculateConfidence(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name       string
		signature  FunctionSignature
		context    map[string]string
		suggestion *TokenSuggestion
		minConf    float64
		maxConf    float64
	}{
		{
			name: "High confidence - exported with parameters",
			signature: FunctionSignature{
				Name:       "ProcessData",
				IsExported: true,
				Parameters: []Parameter{{Name: "data", Type: "string"}},
				ReturnType: "error",
			},
			context:    map[string]string{"error_handling": "true", "resource_management": "true"},
			suggestion: &TokenSuggestion{FeatureID: "ARCH-001", PriorityReason: "High priority"},
			minConf:    0.8,
			maxConf:    1.0,
		},
		{
			name: "Low confidence - unexported no parameters",
			signature: FunctionSignature{
				Name:       "helper",
				IsExported: false,
				Parameters: []Parameter{},
				ReturnType: "void",
			},
			context:    map[string]string{},
			suggestion: &TokenSuggestion{FeatureID: "UTIL-NEW", PriorityReason: "Low priority: utility function"},
			minConf:    0.0,
			maxConf:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := analyzer.calculateConfidence(tt.signature, tt.context, tt.suggestion)

			if confidence < tt.minConf || confidence > tt.maxConf {
				t.Errorf("Expected confidence between %f and %f, got %f", tt.minConf, tt.maxConf, confidence)
			}
		})
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestAnalyzeComplexity(t *testing.T) {
	// This would require creating AST nodes, which is complex for testing
	// In a real implementation, we would mock the AST or use test fixtures
	t.Skip("Complexity analysis requires AST fixtures - implementation would include proper AST testing")
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestAnalyzeContext(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	testCode := `package main

import "fmt"

func ProcessData(data string) error {
	if data == "" {
		return fmt.Errorf("empty data")
	}
	
	defer cleanup()
	
	config := loadConfig()
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	
	return nil
}`

	context, err := analyzer.analyzeContext("test.go", 5, []byte(testCode))
	if err != nil {
		t.Fatalf("Context analysis failed: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	expectedPatterns := map[string]string{
		"error_handling":        "true",
		"resource_management":   "true",
		"configuration_related": "true",
		"file_operations":       "true",
	}

	for key, expected := range expectedPatterns {
		if context[key] != expected {
			t.Errorf("Expected %s to be %s, got %s", key, expected, context[key])
		}
	}

	if context["surrounding_code"] == "" {
		t.Error("Expected surrounding_code to be populated")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestValidateFile(t *testing.T) {
	validator := NewTokenValidator()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testContent := `package main

// Valid token format
// ARCH-001: See architecture.md - Core Architecture [DECISION:maintenance]
func CreateArchive() {}

// Invalid token format - missing icons
// ARCH-002: See architecture.md - Archive Validation [DECISION:maintenance]
func ValidateArchive() {}

// Invalid token format - wrong structure
// Some random comment with ARCH-003
func ProcessArchive() {}

// Valid but different format
// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
func LoadConfig() {}`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	violations, err := validator.validateFile(testFile)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(violations) == 0 {
		t.Error("Expected violations to be found")
	}

	// Check specific violations
	foundInvalidFormat := false
	for _, violation := range violations {
		if violation.ViolationType == "INVALID_FORMAT" {
			foundInvalidFormat = true
		}

		if violation.Severity != "WARNING" {
			t.Errorf("Expected severity WARNING, got %s", violation.Severity)
		}

		if violation.RuleID != "FORMAT-001" {
			t.Errorf("Expected rule ID FORMAT-001, got %s", violation.RuleID)
		}
	}

	if !foundInvalidFormat {
		t.Error("Expected INVALID_FORMAT violation to be found")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestAnalyzeTargetFile(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testContent := `package main

import "fmt"

// CreateBackup creates a backup of the specified file
func CreateBackup(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}
	
	config := loadConfig()
	return processBackup(filePath, config)
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	return &Config{}, nil
}

// validateInput is unexported helper
func validateInput(data string) bool {
	return data != ""
}`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	results, err := analyzer.AnalyzeTarget(testFile)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if results.Target != testFile {
		t.Errorf("Expected target %s, got %s", testFile, results.Target)
	}

	if results.FunctionsAnalyzed != 2 { // Only exported functions
		t.Errorf("Expected 2 functions analyzed, got %d", results.FunctionsAnalyzed)
	}

	if len(results.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(results.Suggestions))
	}

	if results.ProcessingTime == 0 {
		t.Error("Expected processing time to be recorded")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	for _, suggestion := range results.Suggestions {
		if suggestion.FilePath != testFile {
			t.Errorf("Expected file path %s, got %s", testFile, suggestion.FilePath)
		}

		if suggestion.Confidence < 0 || suggestion.Confidence > 1 {
			t.Errorf("Expected confidence between 0 and 1, got %f", suggestion.Confidence)
		}

		if suggestion.SuggestedToken == "" {
			t.Error("Expected suggested token to be generated")
		}

		if suggestion.PriorityIcon == "" {
			t.Error("Expected priority icon to be assigned")
		}

		if suggestion.ActionIcon == "" {
			t.Error("Expected action icon to be assigned")
		}

		if suggestion.FeatureID == "" {
			t.Error("Expected feature ID to be assigned")
		}
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestBatchProcessing(t *testing.T) {
	processor := NewBatchProcessor()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := t.TempDir()

	// Create multiple test files
	testFiles := map[string]string{
		"file1.go": `package main
func CreateArchive() error { return nil }
func LoadConfig() error { return nil }`,

		"file2.go": `package main
func ProcessBackup(path string) error { return nil }
func ValidateData(data []byte) bool { return true }`,

		"file3_test.go": `package main
func TestSomething() {}`, // This should be skipped
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	results, err := processor.ProcessDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Batch processing failed: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if results.Directory != tmpDir {
		t.Errorf("Expected directory %s, got %s", tmpDir, results.Directory)
	}

	if results.FilesProcessed != 2 { // Should skip test file
		t.Errorf("Expected 2 files processed, got %d", results.FilesProcessed)
	}

	if results.TotalFunctions != 4 { // 2 + 2 exported functions
		t.Errorf("Expected 4 functions total, got %d", results.TotalFunctions)
	}

	if results.TotalSuggestions != 4 {
		t.Errorf("Expected 4 suggestions total, got %d", results.TotalSuggestions)
	}

	if results.ProcessingTime == 0 {
		t.Error("Expected processing time to be recorded")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	totalPriority := results.PriorityBreakdown.Critical + results.PriorityBreakdown.High +
		results.PriorityBreakdown.Medium + results.PriorityBreakdown.Low
	if totalPriority != results.TotalSuggestions {
		t.Errorf("Priority breakdown doesn't match total suggestions: %d vs %d", totalPriority, results.TotalSuggestions)
	}

	totalAction := results.ActionBreakdown.Analysis + results.ActionBreakdown.Documentation +
		results.ActionBreakdown.Configuration + results.ActionBreakdown.Protection
	if totalAction != results.TotalSuggestions {
		t.Errorf("Action breakdown doesn't match total suggestions: %d vs %d", totalAction, results.TotalSuggestions)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(results.TopSuggestions) > 10 {
		t.Error("Expected top suggestions to be limited to 10")
	}

	// Check suggestions are sorted by confidence
	for i := 1; i < len(results.TopSuggestions); i++ {
		if results.TopSuggestions[i-1].Confidence < results.TopSuggestions[i].Confidence {
			t.Error("Expected top suggestions to be sorted by confidence (descending)")
		}
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestErrorHandling(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	_, err := analyzer.AnalyzeTarget("non-existent-file.go")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.go")

	err = os.WriteFile(invalidFile, []byte("invalid go syntax {{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	_, err = analyzer.AnalyzeTarget(invalidFile)
	if err == nil {
		t.Error("Expected error for invalid Go file")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestDefaultAnalysisConfig(t *testing.T) {
	config := DefaultAnalysisConfig()

	if config == nil {
		t.Fatal("Expected config to be created")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(config.PriorityRules.CriticalPatterns) == 0 {
		t.Error("Expected critical patterns to be defined")
	}

	if len(config.PriorityRules.HighPatterns) == 0 {
		t.Error("Expected high patterns to be defined")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(config.ActionRules.AnalysisPatterns) == 0 {
		t.Error("Expected analysis patterns to be defined")
	}

	if len(config.ActionRules.ConfigurationPatterns) == 0 {
		t.Error("Expected configuration patterns to be defined")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	totalWeight := config.ConfidenceWeights.SignatureMatch +
		config.ConfidenceWeights.PatternMatch +
		config.ConfidenceWeights.ContextMatch +
		config.ConfidenceWeights.FeatureMapping +
		config.ConfidenceWeights.ComplexityFactor

	if totalWeight != 1.0 {
		t.Errorf("Expected confidence weights to sum to 1.0, got %f", totalWeight)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if len(config.ValidationRules.ValidPriorityIcons) != 4 {
		t.Errorf("Expected 4 priority icons, got %d", len(config.ValidationRules.ValidPriorityIcons))
	}

	if len(config.ValidationRules.ValidActionIcons) != 4 {
		t.Errorf("Expected 4 action icons, got %d", len(config.ValidationRules.ValidActionIcons))
	}

	if config.ValidationRules.MinConfidence < 0 || config.ValidationRules.MinConfidence > 1 {
		t.Errorf("Expected min confidence between 0 and 1, got %f", config.ValidationRules.MinConfidence)
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func TestUtilityFunctions(t *testing.T) {
	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if min(5, 3) != 3 {
		t.Error("min function failed for integers")
	}

	if max(5, 3) != 5 {
		t.Error("max function failed for integers")
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	if minFloat64(5.5, 3.3) != 3.3 {
		t.Error("minFloat64 function failed")
	}

	if maxFloat64(5.5, 3.3) != 5.5 {
		t.Error("maxFloat64 function failed")
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func BenchmarkAnalyzeTarget(b *testing.B) {
	analyzer := NewTokenAnalyzer()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "benchmark.go")

	testContent := `package main

func CreateArchive(path string) error { return nil }
func LoadConfig() (*Config, error) { return nil, nil }
func ProcessBackup(data []byte) error { return nil }
func ValidateInput(input string) bool { return true }
func FormatOutput(data interface{}) string { return "" }`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeTarget(testFile)
		if err != nil {
			b.Fatalf("Analysis failed: %v", err)
		}
	}
}

// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
func BenchmarkBatchProcessing(b *testing.B) {
	processor := NewBatchProcessor()

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	tmpDir := b.TempDir()

	// Create multiple files for batch processing
	for i := 0; i < 10; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("file%d.go", i))
		content := fmt.Sprintf(`package main
func Function%dA() error { return nil }
func Function%dB() bool { return true }`, i, i)

		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			b.Fatalf("Failed to create benchmark file: %v", err)
		}
	}

	// DOC-010: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.ProcessDirectory(tmpDir)
		if err != nil {
			b.Fatalf("Batch processing failed: %v", err)
		}
	}
}
