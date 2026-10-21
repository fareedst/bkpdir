# [IMPL-CFG_HIERARCHY_PRESERVATION] [ARCH-CONFIG_SYSTEM] [ARCH-STDD_VIS_FLOW] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Preserve configuration field values across sequential and inheritance merge paths so earlier explicit settings are not overwritten by later files or defaults.

INPUT: dst config, src config, inherit_context, field name
OUTPUT: merged config with correct precedence per merge mode
DATA: explicitly_set_tracker, default values, inheritContext flag

## MERGE_BASIC_SETTINGS_WITH_PRESERVATION

SPEC-ID: IMPL-CFG_HIERARCHY_PRESERVATION::MERGE_BASIC_SETTINGS_WITH_PRESERVATION
STEP T001: for sequential files, override a field only when source is explicit and no earlier file already set it

- [IMPL-CFG_HIERARCHY_PRESERVATION] [ARCH-CONFIG_SYSTEM] [ARCH-STDD_VIS_FLOW] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: for sequential files, override a field only when source is explicit and no earlier file already set it.

PROCEDURE merge_basic_settings_with_preservation(dst, src, inherit_context):
  FOR each field in basic_settings:
    IF inherit_context THEN
      IF src.field is explicitly set THEN
        dst.field = src.field
      END IF
    ELSE
      IF src.field is explicitly_set AND NOT dst_differs_from_default(dst.field) AND NOT explicitly_set_by_earlier(field) THEN
        dst.field = src.field
        RECORD field in explicitly_set_tracker
      END IF
    END IF
  END FOR

## APPLY_OVERRIDE_WITH_PRESERVATION

SPEC-ID: IMPL-CFG_HIERARCHY_PRESERVATION::APPLY_OVERRIDE_WITH_PRESERVATION
STEP T001: preserve dst when field was set by an earlier file or differs from compiled default; otherwise apply src

- [IMPL-CFG_HIERARCHY_PRESERVATION] [ARCH-CONFIG_SYSTEM] [ARCH-STDD_VIS_FLOW] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: preserve dst when field was set by an earlier file or differs from compiled default; otherwise apply src.

PROCEDURE apply_override_with_preservation(dst, src, field):
  IF field was explicitly set by an earlier config file THEN
    PRESERVE dst.field
  ELSE IF dst.field differs from compiled default THEN
    PRESERVE dst.field
  ELSE
    dst.field = src.field
  END IF
