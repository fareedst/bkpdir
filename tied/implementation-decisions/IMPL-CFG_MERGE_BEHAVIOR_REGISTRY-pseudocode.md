# [IMPL-CFG_MERGE_BEHAVIOR_REGISTRY] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Per-field registry choosing accumulate versus precedence merge behavior, used when applying replace and other strategies during multi-file config load.

INPUT: field name, dst value, explicitlySetFields map
OUTPUT: FieldMergeBehavior enum, merged Config fields
DATA: fieldMergeBehaviors map, MergeBehaviorAccumulate, MergeBehaviorPrecedence

## GET_FIELD_MERGE_BEHAVIOR

SPEC-ID: IMPL-CFG_MERGE_BEHAVIOR_REGISTRY::GET_FIELD_MERGE_BEHAVIOR
STEP T001: LOOKUP field in registry OR default to precedence

- [IMPL-CFG_MERGE_BEHAVIOR_REGISTRY] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: return registry entry for field or default to precedence behavior for unknown fields.

PROCEDURE GET_FIELD_MERGE_BEHAVIOR(fieldName):
  IF fieldName in fieldMergeBehaviors THEN RETURN mapped behavior
  RETURN MergeBehaviorPrecedence

## APPLY_REPLACE

SPEC-ID: IMPL-CFG_MERGE_BEHAVIOR_REGISTRY::APPLY_REPLACE
STEP T001: SKIP replace when earlier file set field or dst non-default
STEP T002: SET result field to new value

- [IMPL-CFG_MERGE_BEHAVIOR_REGISTRY] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: on sequential merge skip replace when earlier file set field or dst differs from default; otherwise set field to new value.

PROCEDURE APPLY_REPLACE(result, key, value, dstValue, inheritContext, defaultCfg, explicitlySetFields):
  IF NOT inheritContext AND dst non-default OR explicitlySetFields[key] THEN RETURN without change
  SET result.key = value
