# [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002]

## Summary contract

Guards list-command formatting so printf-style formatting runs only when the format string contains printf verbs, preventing EXTRA-arg errors and misinterpreting template placeholders.

INPUT: path, creationTime, archivePath, cfg format strings
OUTPUT: formatted list line string
DATA: formatStr, data map with path, creation_time, file statistics

## FORMAT_LIST_ARCHIVE

SPEC-ID: IMPL-LIST_FORMAT_SAFETY::FORMAT_LIST_ARCHIVE
STEP T001: Route template placeholders before printf and return plain format when no verbs

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: route template placeholders before printf and return plain format when no verbs.

PROCEDURE FORMAT_LIST_ARCHIVE(path, creationTime):
  formatStr = cfg.FormatListArchive
  IF CONTAINS(formatStr, "#{") THEN
    data = BUILD_DATA_MAP(path, creationTime)
    statInfo, err = GATHER_FILE_STAT_INFO(path)
    IF err == nil THEN
      MERGE statInfo fields into data
    ELSE
      FILL data with safe defaults ("unknown", "0")
    END IF
    RETURN FORMAT_TEMPLATE(formatStr, data)
  END IF
  IF CONTAINS(formatStr, "%") THEN
    RETURN FORMAT_WITH_PLACEHOLDERS(formatStr, path, creationTime)
  END IF
  RETURN formatStr

## FORMAT_LIST_BACKUP

SPEC-ID: IMPL-LIST_FORMAT_SAFETY::FORMAT_LIST_BACKUP
STEP T001: Mirror FORMAT_LIST_ARCHIVE guard pattern using cfg.FormatListBackup

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: mirror FORMAT_LIST_ARCHIVE guard pattern using cfg.FormatListBackup.

PROCEDURE FORMAT_LIST_BACKUP(path, creationTime):
  formatStr = cfg.FormatListBackup
  IF CONTAINS(formatStr, "#{") THEN
    data = BUILD_DATA_MAP(path, creationTime)
    statInfo, err = GATHER_FILE_STAT_INFO(path)
    IF err == nil THEN
      MERGE statInfo fields (size, size_human, mtime, mtime_unix, mode, type, name) into data
    ELSE
      SET data["size"] = "0"
      SET data["size_human"] = "unknown"
      SET data["mtime"] = creationTime
      SET data["mtime_unix"] = "0"
      SET data["mode"] = "unknown"
      SET data["type"] = "unknown"
      SET data["name"] = BASENAME(path)
    END IF
    RETURN FORMAT_TEMPLATE(formatStr, data)
  END IF
  IF CONTAINS(formatStr, "%") THEN
    RETURN FORMAT_WITH_PLACEHOLDERS(formatStr, path, creationTime)
  END IF
  RETURN formatStr

## FORMAT_LIST_ARCHIVE_WITH_EXTRACTION

SPEC-ID: IMPL-LIST_FORMAT_SAFETY::FORMAT_LIST_ARCHIVE_WITH_EXTRACTION
STEP T001: Priority-based format selection with extraction data and guarded printf

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: priority-based format selection with extraction data and guarded printf.

PROCEDURE FORMAT_LIST_ARCHIVE_WITH_EXTRACTION(archivePath, creationTime):
  data = EXTRACT_ARCHIVE_FILENAME_DATA(archivePath)
  data["path"] = archivePath
  data["creation_time"] = creationTime
  GATHER file stats into data with safe defaults on error
  formatStr = cfg.FormatListArchive
  IF formatStr != "" AND CONTAINS(formatStr, "#{") THEN
    RETURN FORMAT_TEMPLATE(formatStr, data)
  END IF
  IF formatStr != "" AND CONTAINS(formatStr, "%") THEN
    RETURN FORMAT_WITH_PLACEHOLDERS(formatStr, archivePath, creationTime)
  END IF
  IF formatStr != "" THEN
    RETURN formatStr
  END IF
  IF cfg.TemplateListArchive != "" THEN
    RETURN FORMAT_LIST_ARCHIVE_TEMPLATE(data)
  END IF
  RETURN FORMAT_LIST_ARCHIVE(archivePath, creationTime)

## FORMAT_LIST_ARCHIVE_SIMPLE

SPEC-ID: IMPL-LIST_FORMAT_SAFETY::FORMAT_LIST_ARCHIVE_SIMPLE
STEP T001: Simplified adapter path always gathers stats and selects format by priority

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: simplified adapter path always gathers stats and selects format by priority.

PROCEDURE FORMAT_LIST_ARCHIVE_SIMPLE(cfg, formatter, archivePath, creationTime):
  data = MAP with path=archivePath and creation_time=creationTime
  statInfo, err = GATHER_FILE_STAT_INFO(archivePath)
  IF err == nil THEN
    MERGE statInfo fields into data
  ELSE
    SET missing stat fields in data to safe defaults ("unknown", "0", BASENAME(archivePath))
  END IF
  IF cfg.FormatListArchive != "" AND CONTAINS(cfg.FormatListArchive, "#{") THEN
    SET selectedFormat = cfg.FormatListArchive
  ELSE IF cfg.TemplateListArchive != "" THEN
    SET selectedFormat = cfg.TemplateListArchive
  ELSE
    SET selectedFormat = "#{path} (created: #{creation_time})\n"
  END IF
  RETURN REPLACE_PLACEHOLDERS(selectedFormat, data)

## FORMAT_WITH_CONTEXT

SPEC-ID: IMPL-LIST_FORMAT_SAFETY::FORMAT_WITH_CONTEXT
STEP T001: Route FormatTypeList through template detection before printf guards

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: guard list FormatWithContext branch with template detection before printf.

PROCEDURE FORMAT_WITH_CONTEXT(ctx):
  SWITCH ctx.FormatType
  CASE FormatTypeList:
    APPLY template and printf guards before formatting list output
  DEFAULT:
    DELEGATE to archive format helpers

## EMBEDDED_MINITEST: template placeholder list archive

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: verify template-style FormatListArchiveWithExtraction replaces #{size_human} and includes path.

## EMBEDDED_MINITEST: printf-style list archive

- [IMPL-LIST_FORMAT_SAFETY] [ARCH-OUTPUT_FORMATTING] [REQ-OUT_002] — How: verify printf-style FormatListArchiveWithExtraction matches FORMAT_WITH_PLACEHOLDERS output without EXTRA args.
