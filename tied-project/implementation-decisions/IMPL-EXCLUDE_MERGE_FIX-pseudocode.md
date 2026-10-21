# [IMPL-EXCLUDE_MERGE_FIX] [ARCH-EXCLUDE_MERGE_FIX] [REQ-CFG_005] [REQ-CONFIGURATION] [REQ-TEST_EXCLUDE_MERGE]

## Summary contract

Fix exclude_patterns inheritance so unprefixed YAML values merge and deduplicate across config files instead of replacing prior patterns, with typed slice conversion from YAML.

INPUT: dst/src Config, rawSrcMap with strategy prefixes
OUTPUT: merged Config.ExcludePatterns
DATA: hasExplicitPrefix map, MergeBehaviorAccumulate, merge operations

## APPLY_MERGE_STRATEGIES_EXCLUDE

SPEC-ID: IMPL-EXCLUDE_MERGE_FIX::APPLY_MERGE_STRATEGIES_EXCLUDE
STEP T001: CHANGE unprefixed override to merge for accumulate fields

- [IMPL-EXCLUDE_MERGE_FIX] [ARCH-EXCLUDE_MERGE_FIX] [REQ-CFG_005] [REQ-CONFIGURATION] [REQ-TEST_EXCLUDE_MERGE] — How: for accumulate fields without explicit prefix, change unprefixed override to merge so exclude_patterns append across inheritance files.

PROCEDURE APPLY_MERGE_STRATEGIES_EXCLUDE(dst, src, rawSrcMap):
  PROCESS raw keys preserving ! + ^ = prefixes
  IF field behavior is Accumulate AND strategy is override AND no explicit prefix:
    SET strategy = merge  # except first LoadConfig file may use replace for built-in defaults
  APPLY merge operations via applyMergeOperation

## APPLY_MERGE_DEDUPE

SPEC-ID: IMPL-EXCLUDE_MERGE_FIX::APPLY_MERGE_DEDUPE
STEP T001: APPEND unique source items to destination slice via setConfigField

- [IMPL-EXCLUDE_MERGE_FIX] [ARCH-EXCLUDE_MERGE_FIX] [REQ-CFG_005] [REQ-CONFIGURATION] [REQ-TEST_EXCLUDE_MERGE] — How: convert []interface{} to []string, append source items not already in destination slice, write via setConfigField.

PROCEDURE APPLY_MERGE_DEDUPE(result, key, value, dstValue):
  srcSlice = value as []string OR convert []interface{}
  dstSlice = current slice on result
  merged = dstSlice + unique items from srcSlice not in dstSlice
  setConfigField(result, key, merged)

## SET_CONFIG_FIELD

SPEC-ID: IMPL-EXCLUDE_MERGE_FIX::SET_CONFIG_FIELD
STEP T001: ASSIGN known keys WITH YAML slice-to-string conversion

- [IMPL-EXCLUDE_MERGE_FIX] [ARCH-EXCLUDE_MERGE_FIX] [REQ-CFG_005] [REQ-CONFIGURATION] [REQ-TEST_EXCLUDE_MERGE] — How: assign exclude_patterns and other known keys with YAML []interface{} to []string conversion.

PROCEDURE SET_CONFIG_FIELD(cfg, key, value):
  SWITCH key: exclude_patterns, archive_dir_path, ... known scalars
  ON unknown key: RETURN error

## IS_KNOWN_CONFIG_FIELD

SPEC-ID: IMPL-EXCLUDE_MERGE_FIX::IS_KNOWN_CONFIG_FIELD
STEP T001: RETURN whether key is in knownFields set

- [IMPL-EXCLUDE_MERGE_FIX] [ARCH-EXCLUDE_MERGE_FIX] [REQ-CFG_005] [REQ-CONFIGURATION] [REQ-TEST_EXCLUDE_MERGE] — How: guard merge pipeline against unknown YAML keys.

PROCEDURE IS_KNOWN_CONFIG_FIELD(key):
  RETURN key IN knownFields set
