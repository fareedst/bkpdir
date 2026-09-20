# [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION]

## Summary contract

Externalize user-facing error and status format strings on Config with YAML defaults, expose them via GetErrorFormatStrings, format through OutputFormatter methods, and route archive failures through HandleArchiveError instead of hardcoded text.

INPUT: err, cfg format fields, formatter interface
OUTPUT: printed message and process exit status
DATA: Format*Error strings, Template*Error strings, status code map

## DEFAULT_CONFIG_ERROR_FORMATS

SPEC-ID: IMPL-CONFIGURABLE_STRINGS::DEFAULT_CONFIG_ERROR_FORMATS
STEP T001: DECLARE Format and Template error fields on Config struct
STEP T002: SEED DefaultConfig with backward-compatible printf defaults

- [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: declare Format* and Template* error message fields on Config and seed DefaultConfig with backward-compatible printf defaults.

PROCEDURE DEFAULT_CONFIG_ERROR_FORMATS():
  Config struct holds format_disk_full_error, format_permission_error, and sibling YAML keys
  DefaultConfig assigns literal defaults such as "Disk full error: %v\n" for each Format* field
  Template* fields receive named-placeholder defaults for template-based error rendering

## GET_ERROR_FORMAT_STRINGS

SPEC-ID: IMPL-CONFIGURABLE_STRINGS::GET_ERROR_FORMAT_STRINGS
STEP T001: BUILD name-to-format-string map from cfg Format fields

- [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: return a name-to-format-string map so error handlers and tests read configurable messages without hardcoded literals.

PROCEDURE GET_ERROR_FORMAT_STRINGS(cfg):
  BUILD map with keys disk_full, permission, directory_not_found, file_not_found, invalid_directory, invalid_file, failed_write_temp, failed_finalize_file, failed_create_dir_disk, failed_create_dir, failed_access_dir, failed_access_file
  EACH value = corresponding cfg.Format* field
  RETURN map

## FORMAT_DISK_FULL_ERROR

SPEC-ID: IMPL-CONFIGURABLE_STRINGS::FORMAT_DISK_FULL_ERROR
STEP T001: SPRINT configured FormatDiskFullError pattern with err

- [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: sprintf the configured FormatDiskFullError pattern with the underlying error value.

PROCEDURE FORMAT_DISK_FULL_ERROR(formatter, err):
  RETURN fmt.Sprintf(formatter.cfg.FormatDiskFullError, err)

## HANDLE_ARCHIVE_ERROR

SPEC-ID: IMPL-CONFIGURABLE_STRINGS::HANDLE_ARCHIVE_ERROR
STEP T001: CLASSIFY errno-style errors and print via formatter Format methods
STEP T002: RETURN matching cfg status code or generic failure

- [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: classify errno-style errors, print via formatter Format* methods using cfg strings, and return matching cfg status codes.

PROCEDURE HANDLE_ARCHIVE_ERROR(err, cfg, formatter):
  IF err == nil THEN RETURN 0
  IF disk full THEN PrintError(FormatDiskFullError(err)); RETURN cfg.StatusDiskFull
  IF permission denied THEN PrintError(FormatPermissionError(err)); RETURN cfg.StatusPermissionDenied
  IF directory not found THEN PrintError(FormatDirectoryNotFound(err)); RETURN cfg.StatusDirectoryNotFound
  IF structured ArchiveError or BackupError THEN RETURN embedded status code
  PrintError(err.Error()); RETURN 1
