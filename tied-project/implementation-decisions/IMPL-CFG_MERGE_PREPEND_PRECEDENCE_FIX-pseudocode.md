# [IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CONFIGURATION]

## Summary contract

Ensures merge and prepend strategies respect CFG-001 earlier-file precedence for scalar fields registered as MergeBehaviorPrecedence when processing sequential config files (non-inheritance context).

INPUT: result config, field key, operation value, dstValue, inheritContext, defaultCfg, explicitlySetFields
OUTPUT: updated result config field or preserved dstValue
DATA: mergeOperation strategy, fieldMergeBehaviors registry, default vs dst comparison

## APPLYMERGEOPERATION

SPEC-ID: IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX::APPLYMERGEOPERATION
STEP T001: Dispatch override, merge, prepend, replace, and default handlers for one processed key

- [IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CONFIGURATION] — How: dispatch override, merge, prepend, replace, and default handlers for one processed key.

PROCEDURE applyMergeOperation(result, key, operation, dstValue, originalDstValue, inheritContext, defaultCfg, explicitlySetFields):
  SWITCH operation.strategy:
    "override" -> applyOverride(...)
    "merge"    -> applyMerge(...)
    "prepend"  -> applyPrepend(...)
    "replace"  -> applyReplace(...)
    "default"  -> applyDefault(...)

## APPLYMERGE_SCALAR_PRECEDENCE

SPEC-ID: IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX::APPLYMERGE_SCALAR_PRECEDENCE
STEP T001: EXECUTE PROCEDURE applyMerge

- [IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CONFIGURATION] — How: when + merge hits a non-array MergeBehaviorPrecedence field in sequential loading, keep dstValue if an earlier file already set it or dst differs from default.

PROCEDURE applyMerge(result, key, value, dstValue, inheritContext, defaultCfg, explicitlySetFields):
  IF value is not array AND getFieldMergeBehavior(key) == MergeBehaviorPrecedence AND NOT inheritContext AND dstValue != nil THEN
    IF dstValue differs from default OR explicitlySetFields[key] OR earlier files processed THEN
      RETURN setConfigField(result, key, dstValue)
  ENDIF
  (array path handled by accumulate merge / deduplication)

## APPLYPREPEND

SPEC-ID: IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX::APPLYPREPEND
STEP T001: EXECUTE PROCEDURE applyPrepend

- [IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX] [ARCH-CFG_001] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CONFIGURATION] — How: when ^ prepend hits a non-array MergeBehaviorPrecedence field sequentially, preserve dstValue under the same earlier-file rules as applyMerge scalars.

PROCEDURE applyPrepend(result, key, value, dstValue, inheritContext, defaultCfg, explicitlySetFields):
  IF value is not array AND getFieldMergeBehavior(key) == MergeBehaviorPrecedence AND NOT inheritContext AND dstValue != nil THEN
    IF dstValue differs from default OR explicitlySetFields[key] OR earlier files processed THEN
      RETURN setConfigField(result, key, dstValue)
  ENDIF
  (array path prepends source slice before destination slice)
