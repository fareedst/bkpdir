# [IMPL-CFG_QUOTED_KEY_PREFIX] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Parses YAML config maps so quoted merge-prefix keys (e.g. "+exclude_patterns") normalize to the correct strategy and base field name without duplicate unprefixed keys winning.

INPUT: config map[string]interface{} from unmarshaled YAML
OUTPUT: processedConfig with one mergeOperation per clean field key
DATA: explicitPrefixKeys map, extractStrategy prefix table

## PROCESSKEYS

SPEC-ID: IMPL-CFG_QUOTED_KEY_PREFIX::PROCESSKEYS
STEP T001: Two-pass scan prefers explicit-prefix keys over duplicate unprefixed entries when building operations per cleanKey

- [IMPL-CFG_QUOTED_KEY_PREFIX] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: two-pass scan prefers explicit-prefix keys over duplicate unprefixed entries when building operations per cleanKey.

PROCEDURE processKeys(config map):
  PASS 1: record explicitPrefixKeys for keys whose cleaned form starts with ! + ^ =
  PASS 2: FOR EACH key, value:
    strategy, cleanKey = extractStrategy(key)
    IF cleanKey has explicitPrefixKeys entry AND key != explicit entry THEN SKIP
    MERGE into result.operations[cleanKey] favoring explicit-prefix operations

## EXTRACTSTRATEGY

SPEC-ID: IMPL-CFG_QUOTED_KEY_PREFIX::EXTRACTSTRATEGY
STEP T001: EXECUTE PROCEDURE extractStrategy

- [IMPL-CFG_QUOTED_KEY_PREFIX] [ARCH-CFG_005] [REQ-CFG_005] [REQ-CONFIGURATION] — How: strip YAML quotes then map leading ! + ^ = to override merge prepend replace default strategies and return strategy plus base field name.

PROCEDURE extractStrategy(key):
  cleanKey = strip surrounding quotes from key
  IF first char is merge prefix THEN RETURN strategy name and key without prefix
  RETURN "override" and cleanKey
