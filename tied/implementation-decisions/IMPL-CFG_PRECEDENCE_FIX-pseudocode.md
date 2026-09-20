# [IMPL-CFG_PRECEDENCE_FIX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION]

## Summary contract

Fix sequential config file loading so earlier discovery-order files keep explicitly set values even when equal to defaults, using explicitlySetFields instead of comparing only to DefaultConfig.

INPUT: ordered config file paths, raw YAML per file
OUTPUT: merged Config respecting first-wins precedence
DATA: explicitlySetFields map, fileProcessed flag, inheritContext boolean

## LOAD_CONFIG_FALLBACK

SPEC-ID: IMPL-CFG_PRECEDENCE_FIX::LOAD_CONFIG_FALLBACK
STEP T001: WALK search paths in order merging with inheritContext on first file
STEP T002: RECORD every raw key in explicitlySetFields after each merge

- [IMPL-CFG_PRECEDENCE_FIX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: walk search paths in order; first file uses inheritContext true; record every raw key in explicitlySetFields after each merge.

PROCEDURE LOAD_CONFIG_FALLBACK(root):
  cfg = DEFAULT_CONFIG
  explicitlySetFields = empty map
  fileProcessed = false
  FOR EACH configPath IN searchPaths:
    IF file missing THEN CONTINUE
    loadResult = LOAD_SINGLE_FILE(configPath)
    inheritContext = NOT fileProcessed
    cfg = APPLY_MERGE_STRATEGIES(cfg, loadResult.config, inheritContext, loadResult.rawMap, explicitlySetFields)
    fileProcessed = true
    FOR EACH key IN loadResult.rawMap:
      explicitlySetFields[stripMergePrefix(key)] = true
  RETURN cfg

## MERGE_CONFIGS

SPEC-ID: IMPL-CFG_PRECEDENCE_FIX::MERGE_CONFIGS
STEP T001: DELEGATE to category merge helpers with explicitlySetFields

- [IMPL-CFG_PRECEDENCE_FIX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: delegate to category merge helpers passing explicitlySetFields for precedence-aware scalar merges.

PROCEDURE MERGE_CONFIGS(dst, src, inheritContext, explicitlySetFields):
  MERGE_BASIC_SETTINGS(dst, src, inheritContext, explicitlySetFields)
  MERGE_FILE_BACKUP_SETTINGS(...)
  MERGE_STATUS_CODES(...)
  MERGE_FORMAT_STRINGS(...)
  MERGE_TEMPLATES(...)
  MERGE_PATTERNS(...)
  MERGE_GIT_SETTINGS(..., explicitlySetFields)

## MERGE_BASIC_SETTINGS

SPEC-ID: IMPL-CFG_PRECEDENCE_FIX::MERGE_BASIC_SETTINGS
STEP T001: SEQUENTIAL mode writes only keys in rawSrcMap not in explicitlySetFields
STEP T002: INHERITANCE mode allows child override when src differs from default

- [IMPL-CFG_PRECEDENCE_FIX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: sequential mode only writes fields present in rawSrcMap and not already in explicitlySetFields; inheritance mode allows child overrides when src differs from default.

PROCEDURE MERGE_BASIC_SETTINGS(dst, src, inheritContext, rawSrcMap, explicitlySetFields):
  IF inheritContext THEN
    FOR EACH scalar: IF src differs from default THEN dst = src
  ELSE
    FOR EACH scalar: IF key in rawSrcMap AND NOT explicitlySetFields[field] THEN dst = src
