# [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING]

## Summary contract

Enhance created and incremental archive messages by gathering FileStatInfo and substituting #{path}, #{size_human}, #{mtime}, and related placeholders in detailed template format strings with safe fallbacks.

INPUT: archive file path, cfg template strings
OUTPUT: formatted message with stats or fallback to basic printf format
DATA: TemplateCreatedArchiveDetailed, GatherFileStatInfo map

## FORMAT_CREATED_ARCHIVE_WITH_STATS

SPEC-ID: IMPL-FILE_STATISTICS_TEMPLATE_FIX::FORMAT_CREATED_ARCHIVE_WITH_STATS
STEP T001: GATHER FileStatInfo AND replace #{key} placeholders OR fallback

- [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: GatherFileStatInfo, build data map, replace #{key} in TemplateCreatedArchiveDetailed, fall back on error or leftover placeholders.

PROCEDURE FORMAT_CREATED_ARCHIVE_WITH_STATS(path):
  statInfo = GatherFileStatInfo(path) OR RETURN FormatCreatedArchive
  REPLACE placeholders in TemplateCreatedArchiveDetailed FROM data
  IF still #{ in result: RETURN FormatCreatedArchive
  RETURN result

## FORMAT_INCREMENTAL_CREATED_WITH_STATS

SPEC-ID: IMPL-FILE_STATISTICS_TEMPLATE_FIX::FORMAT_INCREMENTAL_CREATED_WITH_STATS
STEP T001: MIRROR created variant using incremental templates and fallback

- [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: same as created variant using TemplateIncrementalCreatedDetailed and FormatIncrementalCreated fallback.

PROCEDURE FORMAT_INCREMENTAL_CREATED_WITH_STATS(path):
  MIRROR FORMAT_CREATED_ARCHIVE_WITH_STATS with incremental templates

## PRINT_CREATED_ARCHIVE_WITH_STATS

SPEC-ID: IMPL-FILE_STATISTICS_TEMPLATE_FIX::PRINT_CREATED_ARCHIVE_WITH_STATS
STEP T001: FORMAT with stats AND route print via collector or stdout

- [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: format with stats then Print via collector or stdout.

PROCEDURE PRINT_CREATED_ARCHIVE_WITH_STATS(path):
  message = FORMAT_CREATED_ARCHIVE_WITH_STATS(path); ROUTE print

## PRINT_INCREMENTAL_CREATED_WITH_STATS

SPEC-ID: IMPL-FILE_STATISTICS_TEMPLATE_FIX::PRINT_INCREMENTAL_CREATED_WITH_STATS
STEP T001: FORMAT incremental with stats AND route print via collector or stdout

- [IMPL-FILE_STATISTICS_TEMPLATE_FIX] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: format incremental with stats then Print via collector or stdout.

PROCEDURE PRINT_INCREMENTAL_CREATED_WITH_STATS(path):
  message = FORMAT_INCREMENTAL_CREATED_WITH_STATS(path); ROUTE print
