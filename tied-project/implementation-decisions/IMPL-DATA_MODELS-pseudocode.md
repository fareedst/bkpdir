# [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP]

## Summary contract

Type-safe archive and backup domain models with configuration and formatter interfaces plus adapters that decouple operations from concrete Config and OutputFormatter types.

INPUT: Config, OutputFormatter, operation parameters
OUTPUT: structured value objects and interface contracts for archive/backup pipelines
DATA: Archive, ArchiveConfig, Backup, BackupInfo, adapter structs

## ARCHIVE_CONFIG

SPEC-ID: IMPL-DATA_MODELS::ARCHIVE_CONFIG
STEP T001: DEFINE ArchiveConfig WITH naming inputs for filename generation

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: hold naming inputs (prefix, timestamp, git segments, note, incremental base) for archive filename generation.

PROCEDURE ARCHIVE_CONFIG():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for ARCHIVE_CONFIG
OUTPUT: result of ARCHIVE_CONFIG
PRE: true
POST: true
EFFECTS: pure
  DEFINE ArchiveConfig with Prefix, Timestamp, GitBranch, GitHash, GitIsClean, ShowGitDirtyStatus, Note, IsGit, IsIncremental, BaseName

## ARCHIVE

SPEC-ID: IMPL-DATA_MODELS::ARCHIVE
STEP T001: DEFINE Archive WITH path, metadata, and incremental flag

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: represent a discovered zip archive with path, creation time, incremental flag, git metadata, and base archive link.

PROCEDURE ARCHIVE():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for ARCHIVE
OUTPUT: result of ARCHIVE
PRE: true
POST: true
EFFECTS: pure
  DEFINE Archive with Name, Path, CreationTime, IsIncremental, GitBranch, GitHash, Note, BaseArchive

## ARCHIVE_CONFIG_INTERFACE

SPEC-ID: IMPL-DATA_MODELS::ARCHIVE_CONFIG_INTERFACE
STEP T001: DECLARE getters for archive dir, exclude patterns, git flags, status codes

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract Config field accessors needed by archive creation without importing main.Config in tests.

PROCEDURE ARCHIVE_CONFIG_INTERFACE():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for ARCHIVE_CONFIG_INTERFACE
OUTPUT: result of ARCHIVE_CONFIG_INTERFACE
PRE: true
POST: true
EFFECTS: pure
  DECLARE getters for archive dir, exclude patterns, git flags, status codes

## ARCHIVE_CREATION_OPTIONS

SPEC-ID: IMPL-DATA_MODELS::ARCHIVE_CREATION_OPTIONS
STEP T001: DEFINE ArchiveCreationOptions value object

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: bundle context, cwd, target path, file list, config interface, and resource manager for create-archive entry points.

PROCEDURE ARCHIVE_CREATION_OPTIONS():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for ARCHIVE_CREATION_OPTIONS
OUTPUT: result of ARCHIVE_CREATION_OPTIONS
PRE: true
POST: true
EFFECTS: pure
  DEFINE ArchiveCreationOptions value object

## ARCHIVE_FORMATTER_INTERFACE

SPEC-ID: IMPL-DATA_MODELS::ARCHIVE_FORMATTER_INTERFACE
STEP T001: DECLARE PrintDryRun AND PrintIncrementalCreated methods

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract dry-run and incremental print operations for archive workflows.

PROCEDURE ARCHIVE_FORMATTER_INTERFACE():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for ARCHIVE_FORMATTER_INTERFACE
OUTPUT: result of ARCHIVE_FORMATTER_INTERFACE
PRE: true
POST: true
EFFECTS: pure
  DECLARE PrintDryRun* and PrintIncrementalCreated methods

## CONFIG_TO_ARCHIVE_CONFIG_ADAPTER

SPEC-ID: IMPL-DATA_MODELS::CONFIG_TO_ARCHIVE_CONFIG_ADAPTER
STEP T001: WRAP *Config AND delegate ArchiveConfigInterface getters

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *Config and delegate each ArchiveConfigInterface getter to the matching Config field.

PROCEDURE CONFIG_TO_ARCHIVE_CONFIG_ADAPTER(cfg):
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: cfg
OUTPUT: result of CONFIG_TO_ARCHIVE_CONFIG_ADAPTER
PRE: true
POST: true
EFFECTS: pure
  RETURN adapter forwarding GetArchiveDirPath, GetExcludePatterns, status helpers, etc.

## OUTPUT_FORMATTER_TO_ARCHIVE_FORMATTER_ADAPTER

SPEC-ID: IMPL-DATA_MODELS::OUTPUT_FORMATTER_TO_ARCHIVE_FORMATTER_ADAPTER
STEP T001: TYPE-ASSERT OutputFormatterInterface to concrete formatter
STEP T002: DELEGATE Print* archive methods

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: type-assert OutputFormatterInterface to FormatterAdapter or AIFormatterAdapter for extended archive print methods.

PROCEDURE OUTPUT_FORMATTER_TO_ARCHIVE_FORMATTER_ADAPTER(formatter):
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: formatter
OUTPUT: result of OUTPUT_FORMATTER_TO_ARCHIVE_FORMATTER_ADAPTER
PRE: true
POST: true
EFFECTS: pure
  DELEGATE Print* to concrete formatter implementations

## INCREMENTAL_ARCHIVE_CONFIG

SPEC-ID: IMPL-DATA_MODELS::INCREMENTAL_ARCHIVE_CONFIG
STEP T001: DEFINE IncrementalArchiveConfig struct

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: group Config, note, dry-run flag, and context for incremental archive API entry.

PROCEDURE INCREMENTAL_ARCHIVE_CONFIG():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for INCREMENTAL_ARCHIVE_CONFIG
OUTPUT: result of INCREMENTAL_ARCHIVE_CONFIG
PRE: true
POST: true
EFFECTS: pure
  DEFINE IncrementalArchiveConfig struct

## BACKUP_CONFIG_INTERFACE

SPEC-ID: IMPL-DATA_MODELS::BACKUP_CONFIG_INTERFACE
STEP T001: DECLARE BackupConfigInterface getters for paths and status codes

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract backup directory path, naming toggle, and backup-related status codes.

PROCEDURE BACKUP_CONFIG_INTERFACE():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for BACKUP_CONFIG_INTERFACE
OUTPUT: result of BACKUP_CONFIG_INTERFACE
PRE: true
POST: true
EFFECTS: pure
  DECLARE BackupConfigInterface getters

## BACKUP_FORMATTER_INTERFACE

SPEC-ID: IMPL-DATA_MODELS::BACKUP_FORMATTER_INTERFACE
STEP T001: DECLARE BackupFormatterInterface Print* methods

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: abstract backup dry-run, identical, created, and not-found output methods.

PROCEDURE BACKUP_FORMATTER_INTERFACE():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for BACKUP_FORMATTER_INTERFACE
OUTPUT: result of BACKUP_FORMATTER_INTERFACE
PRE: true
POST: true
EFFECTS: pure
  DECLARE BackupFormatterInterface Print* methods

## CONFIG_TO_BACKUP_CONFIG_ADAPTER

SPEC-ID: IMPL-DATA_MODELS::CONFIG_TO_BACKUP_CONFIG_ADAPTER
STEP T001: WRAP *Config AND delegate BackupConfigInterface getters

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *Config and delegate BackupConfigInterface getters to Config fields.

PROCEDURE CONFIG_TO_BACKUP_CONFIG_ADAPTER(cfg):
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: cfg
OUTPUT: result of CONFIG_TO_BACKUP_CONFIG_ADAPTER
PRE: true
POST: true
EFFECTS: pure
  RETURN adapter forwarding backup paths and status codes

## OUTPUT_FORMATTER_TO_BACKUP_FORMATTER_ADAPTER

SPEC-ID: IMPL-DATA_MODELS::OUTPUT_FORMATTER_TO_BACKUP_FORMATTER_ADAPTER
STEP T001: WRAP *OutputFormatter AND forward BackupFormatterInterface calls

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: wrap *OutputFormatter and forward each BackupFormatterInterface call.

PROCEDURE OUTPUT_FORMATTER_TO_BACKUP_FORMATTER_ADAPTER(formatter):
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: formatter
OUTPUT: result of OUTPUT_FORMATTER_TO_BACKUP_FORMATTER_ADAPTER
PRE: true
POST: true
EFFECTS: pure
  DELEGATE PrintDryRunBackup, PrintBackupCreated, etc.

## BACKUP_INFO

SPEC-ID: IMPL-DATA_MODELS::BACKUP_INFO
STEP T001: DEFINE BackupInfo struct WITH list/compare metadata

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: capture list/compare metadata (name, path, creation time, size) for backup discovery.

PROCEDURE BACKUP_INFO():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for BACKUP_INFO
OUTPUT: result of BACKUP_INFO
PRE: true
POST: true
EFFECTS: pure
  DEFINE BackupInfo struct

## BACKUP

SPEC-ID: IMPL-DATA_MODELS::BACKUP
STEP T001: DEFINE Backup struct WITH source path AND optional note

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: represent one backup file with source path and optional note segment.

PROCEDURE BACKUP():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for BACKUP
OUTPUT: result of BACKUP
PRE: true
POST: true
EFFECTS: pure
  DEFINE Backup struct with Name, Path, CreationTime, SourceFile, Note

## BACKUP_OPTIONS

SPEC-ID: IMPL-DATA_MODELS::BACKUP_OPTIONS
STEP T001: DEFINE BackupOptions value object for createFileBackupInternal

- [IMPL-DATA_MODELS] [ARCH-SYSTEM_COMPONENTS] [REQ-CODE_QUALITY] [REQ-FILE_BACKUP] — How: bundle context, config, formatter, file path, note, and dry-run for createFileBackupInternal.

PROCEDURE BACKUP_OPTIONS():
  # [IMPL-DATA_MODELS] — Layer B procedure token link
Contract:
INPUT: context for BACKUP_OPTIONS
OUTPUT: result of BACKUP_OPTIONS
PRE: true
POST: true
EFFECTS: pure
  DEFINE BackupOptions value object
