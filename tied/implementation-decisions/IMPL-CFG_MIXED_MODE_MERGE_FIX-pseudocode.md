# [IMPL-CFG_MIXED_MODE_MERGE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Applies per-field merge behavior when combining config sources: accumulate fields default to merge unless an explicit YAML prefix (!, +, ^, =) is present, with special handling for exclude_patterns and inheritance vs sequential context.

INPUT: dst and src Config, inheritContext, rawSrcMap, initialDefaultCfg, explicitlySetFields, excludeUnprefixedReplacesBuiltinDefaults
OUTPUT: merged Config or error
DATA: processed merge operations map, fieldMergeBehaviors, hasExplicitPrefix map

## APPLYMERGESTRATEGIES

SPEC-ID: IMPL-CFG_MIXED_MODE_MERGE_FIX::APPLYMERGESTRATEGIES
STEP T001: EXECUTE PROCEDURE applyMergeStrategies

- [IMPL-CFG_MIXED_MODE_MERGE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: process prefixed keys into operations, honor explicit prefixes over accumulate defaults, and apply each operation via applyMergeOperation with registry-driven behavior.

PROCEDURE applyMergeStrategies(dst, src, inheritContext, rawSrcMap, initialDefaultCfg, explicitlySetFields, excludeUnprefixedReplacesBuiltinDefaults):
  processed = processor.processKeys(srcMap from rawSrcMap or configToMap(src))
  BUILD hasExplicitPrefix from rawSrcMap keys (quoted keys stripped)
  FOR EACH key, operation IN processed.operations:
    behavior = getFieldMergeBehavior(key)
    IF behavior == MergeBehaviorAccumulate AND NOT hasExplicitPrefix[key] THEN
      ADJUST operation.strategy for inheritContext / exclude_patterns rules
    applyMergeOperation(dst, key, operation, ...)
