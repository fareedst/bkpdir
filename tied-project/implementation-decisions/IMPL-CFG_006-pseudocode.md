# [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006]

## Summary contract

Reflection-based configuration field discovery with thread-safe schema cache, value formatting, and merge-strategy inference. Flat YAML-tag reflection and display grouping are owned by [IMPL-CONFIG_DISPLAY_FLATTENING] (reflectConfigFields, determineFieldCategory). Backward compatibility: existing GetConfigValuesWithSources callers unchanged. Source tracking shows inheritance chain origins per field.

INPUT: Config instance, field paths, reflect types
OUTPUT: configFieldInfo slices, display strings, merge strategy names
DATA: globalFieldCache, configFieldInfo, ConfigValueWithMetadata

## GET_ALL_CONFIG_FIELDS

SPEC-ID: IMPL-CFG_006::GET_ALL_CONFIG_FIELDS
STEP T001: RETURN cached fields WHEN schema hash matches
STEP T002: reflectConfigFields AND cache metadata
STEP T003: PRESERVE backward compatibility for GetConfigValuesWithSources callers

- [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: return cached field metadata when schema hash matches; on miss call reflectConfigFields (IMPL-CONFIG_DISPLAY_FLATTENING), sort, strip values, cache metadata.

PROCEDURE GET_ALL_CONFIG_FIELDS(cfg):
Contract:
PRE: true
POST: true
INPUT: cfg
OUTPUT: result of GET_ALL_CONFIG_FIELDS
EFFECTS: State.Config, State.FieldCache, pure
  cached = globalFieldCache.getCachedFields()
  IF cached != nil THEN RETURN updateFieldValues(cached, cfg)
  fields = reflectConfigFields(TYPE_OF(cfg), VALUE_OF(cfg), "", "")  # see IMPL-CONFIG_DISPLAY_FLATTENING
  SORT fields BY Name
  metadata = COPY fields WITH Value = nil AND Importance filled
  globalFieldCache.setCachedFields(metadata)
  RETURN fields WITH live values

## UPDATE_FIELD_VALUES

SPEC-ID: IMPL-CFG_006::UPDATE_FIELD_VALUES
STEP T001: FOR EACH field SET Value from getFieldValueByPath

- [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: copy cached metadata and populate Value via getFieldValueByPath or zero value on failure.

PROCEDURE UPDATE_FIELD_VALUES(cachedFields, cfg):
Contract:
PRE: true
POST: true
INPUT: cachedFields, cfg
OUTPUT: result of UPDATE_FIELD_VALUES
EFFECTS: State.Config, State.FieldCache, pure
  FOR EACH field IN cachedFields:
    COPY metadata
    SET Value = getFieldValueByPath(cfg, field.Path) OR zero for Kind

## GET_FIELD_VALUE_BY_PATH

SPEC-ID: IMPL-CFG_006::GET_FIELD_VALUE_BY_PATH
STEP T001: WALK struct fields by dot path

- [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: split dot path, dereference pointers, walk struct fields by name, return leaf interface.

PROCEDURE GET_FIELD_VALUE_BY_PATH(structValue, path):
Contract:
PRE: true
POST: true
INPUT: structValue, path
OUTPUT: result of GET_FIELD_VALUE_BY_PATH
EFFECTS: State.Config, State.FieldCache, pure
  FOR EACH part IN SPLIT(path, "."):
    IF pointer THEN dereference OR error if nil
    IF not struct THEN error
    current = FieldByName(part)
  RETURN interface value

## FORMAT_FIELD_VALUE

SPEC-ID: IMPL-CFG_006::FORMAT_FIELD_VALUE
STEP T001: SWITCH on reflect kind for display string

- [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: format nil, bool, numeric, string, string-slice, pointer, and default kinds for config display.

PROCEDURE FORMAT_FIELD_VALUE(value, kind):
Contract:
PRE: true
POST: true
INPUT: value, kind
OUTPUT: result of FORMAT_FIELD_VALUE
EFFECTS: State.Config, State.FieldCache, pure
  IF value nil THEN RETURN "<nil>"
  SWITCH kind FOR bool, ints, string, slice, pointer recurse, default stringify

## DETERMINE_MERGE_STRATEGY_FOR_FIELD

SPEC-ID: IMPL-CFG_006::DETERMINE_MERGE_STRATEGY_FOR_FIELD
STEP T001: MATCH path patterns to merge strategy

- [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: infer append/prepend/merge/override/default from field path patterns and value comparison to old value.

PROCEDURE DETERMINE_MERGE_STRATEGY_FOR_FIELD(fieldPath, newValue, oldValue):
Contract:
PRE: true
POST: true
INPUT: fieldPath, newValue, oldValue
OUTPUT: result of DETERMINE_MERGE_STRATEGY_FOR_FIELD
EFFECTS: State.Config, State.FieldCache, pure
  IF path matches exclude_patterns OR inherit AND slice THEN append
  IF path matches patterns AND NOT exclude_patterns THEN prepend
  IF path contains git THEN merge
  IF path contains verification THEN override
  IF values differ with zero transitions THEN default OR override
  DEFAULT override

## DETECT_MERGE_STRATEGY_FROM_FIELD

SPEC-ID: IMPL-CFG_006::DETECT_MERGE_STRATEGY_FROM_FIELD
STEP T001: MAP field shape flags to strategy name

- [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: map field shape flags and name hints to append, merge, prepend, or override.

PROCEDURE DETECT_MERGE_STRATEGY_FROM_FIELD(field):
Contract:
PRE: true
POST: true
INPUT: field
OUTPUT: result of DETECT_MERGE_STRATEGY_FROM_FIELD
EFFECTS: State.Config, State.FieldCache, pure
  IF field.IsSlice THEN append
  IF field.IsStruct THEN merge
  IF field.IsPointer THEN override
  IF name contains Pattern THEN prepend
  IF name contains Git THEN merge
  IF name contains Verification THEN override
  DEFAULT override

## CONFIG_FIELD_CACHE

SPEC-ID: IMPL-CFG_006::CONFIG_FIELD_CACHE
STEP T001: RWMutex get/set with struct hash invalidation

- [IMPL-CFG_006] [ARCH-CFG_006] [ARCH-SYSTEM_COMPONENTS] [REQ-CFG_006] — How: RWMutex-backed cache invalidated when Config struct type hash changes.

PROCEDURE CONFIG_FIELD_CACHE:
  getCachedFields: IF valid AND structHash matches THEN return copy ELSE nil
  setCachedFields: store copy, hash, timestamp, valid=true

## EMBEDDED_MINITEST_MERGE_STRATEGY

SPEC-ID: IMPL-CFG_006::EMBEDDED_MINITEST_MERGE_STRATEGY
STEP T001: TABLE test determineMergeStrategyForField

- [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: table-driven test asserts determineMergeStrategyForField paths map to expected strategies.

PROCEDURE EMBEDDED_MINITEST_MERGE_STRATEGY():
Contract:
PRE: true
POST: true
INPUT: context for EMBEDDED_MINITEST_MERGE_STRATEGY
OUTPUT: result of EMBEDDED_MINITEST_MERGE_STRATEGY
EFFECTS: State.Config, State.FieldCache, pure
  FOR EACH case CALL determineMergeStrategyForField AND ASSERT expected
