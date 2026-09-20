# [IMPL-CONFIG_OUTPUT_GROUPING] [ARCH-CONFIG_OUTPUT_GROUPING] [REQ-CONFIG_OUTPUT_GROUPING]

## Summary contract

Rank configuration fields by importance and category priority so grouped config display lists critical paths first within each category section.

INPUT: field name, category, ConfigValueWithMetadata list, showSources flag
OUTPUT: grouped stdout sections sorted by category priority then field importance
DATA: Importance* constants, CategoryPriority map, categories map

## GET_FIELD_IMPORTANCE

SPEC-ID: IMPL-CONFIG_OUTPUT_GROUPING::GET_FIELD_IMPORTANCE
STEP T001: MAP field name patterns to Importance level

- [IMPL-CONFIG_OUTPUT_GROUPING] [ARCH-CONFIG_OUTPUT_GROUPING] [REQ-CONFIG_OUTPUT_GROUPING] — How: map archive and backup dir paths to critical, key toggles to high, format/template/status names to low, otherwise medium.

PROCEDURE GET_FIELD_IMPORTANCE(name, category):
  IF name IN {archive_dir_path, backup_dir_path} THEN RETURN ImportanceCritical
  IF name IN {use_current_dir_name, use_current_dir_name_for_files, git.enabled} THEN RETURN ImportanceHigh
  IF name CONTAINS "format" OR "template" THEN RETURN ImportanceLow
  IF name CONTAINS "status" THEN RETURN ImportanceLow
  RETURN ImportanceMedium

## DISPLAY_CONFIG_GROUPED

SPEC-ID: IMPL-CONFIG_OUTPUT_GROUPING::DISPLAY_CONFIG_GROUPED
STEP T001: GROUP values by category and sort by priority and importance
STEP T002: PRINT section headers and formatted field rows

- [IMPL-CONFIG_OUTPUT_GROUPING] [ARCH-CONFIG_OUTPUT_GROUPING] [REQ-CONFIG_OUTPUT_GROUPING] — How: bucket values by category, sort categories by CategoryPriority, sort fields by importance then name, print sources header and section headers.

PROCEDURE DISPLAY_CONFIG_GROUPED(values, showSources):
  categories = GROUP values BY FieldInfo.Category
  sortedCategories = SORT category keys BY CategoryPriority (unknown = 999) THEN alphabetically
  PRINT expanded configuration search paths header
  FOR EACH category IN sortedCategories:
    SORT category values BY (Importance, Name)
    PRINT "## " + titleCase(category)
    FOR EACH value: PRINT indented name, value, optional source/override/chain annotations
