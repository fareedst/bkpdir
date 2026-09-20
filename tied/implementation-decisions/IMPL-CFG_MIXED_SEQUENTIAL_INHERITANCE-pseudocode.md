# [IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Load configs from search paths mixing single-file sequential discovery with multi-file inheritance chains, tracking explicitlySetFields so sequential files win over later files but inheritance children can override parents.

INPUT: search paths, inherit chains, raw YAML maps
OUTPUT: merged Config, explicitlySetFields map
DATA: inheritContext flag, chain.files length

## LOAD_CONFIG_WITH_INHERITANCE

SPEC-ID: IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE::LOAD_CONFIG_WITH_INHERITANCE
STEP T001: FOR each search path BUILD chain and MERGE files in order
STEP T002: TRACK explicitlySetFields for single-file chains only

- [IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: for each search path build inheritance chain, merge files in order, track explicitlySetFields only for single-file chains.

PROCEDURE LOAD_CONFIG_WITH_INHERITANCE(root):
  finalCfg = DEFAULT_CONFIG
  explicitlySetFields = empty map
  FOR EACH searchPath IN discovery order:
    IF file missing THEN CONTINUE
    chain = BUILD_CHAIN(searchPath)
    isSingleFile = len(chain.files) == 1
    FOR EACH file IN chain.files:
      tempCfg, rawMap = LOAD_FILE(file)
      inheritContext = len(chain.files) > 1 OR is first file
      finalCfg = APPLY_MERGE_STRATEGIES(finalCfg, tempCfg, inheritContext, rawMap, explicitlySetFields)
    IF isSingleFile AND rawMap present THEN
      FOR EACH key IN rawMap: explicitlySetFields[stripPrefix(key)] = true
  RETURN finalCfg

## MERGE_BASIC_SETTINGS

SPEC-ID: IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE::MERGE_BASIC_SETTINGS
STEP T001: APPLY inheritance or sequential precedence using explicitlySetFields

- [IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: in inheritance context skip precedence scalars when earlier sequential file set field; in sequential context skip when explicitlySetFields or dst differs from default.

PROCEDURE MERGE_BASIC_SETTINGS(dst, src, inheritContext, explicitlySetFields, rawSrcMap):
  FOR EACH precedence scalar field:
    IF inheritContext THEN
      IF explicitlySetFields[field] THEN SKIP
      ELSE apply when src differs from default
    ELSE
      IF explicitlySetFields[field] OR dst differs from default THEN SKIP
      ELSE apply when key in rawSrcMap

## APPLY_OVERRIDE

SPEC-ID: IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE::APPLY_OVERRIDE
STEP T001: PRESERVE dst when earlier sequential file set key OR apply new value

- [IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: preserve dst when earlier sequential file set key or sequential dst differs from default; inheritance only checks explicitlySetFields.

PROCEDURE APPLY_OVERRIDE(result, key, value, dstValue, inheritContext, explicitlySetFields):
  IF shouldPreserve(dstValue, default, explicitlySetFields, inheritContext) THEN
    RETURN setConfigField(result, key, dstValue)
  RETURN setConfigField(result, key, value)
