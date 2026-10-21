# [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS]

## Summary contract

Validate configurable format and template strings against per-field placeholder rules, warn on unexpected placeholders during config load, and render #{key} and mixed printf templates with safe defaults. Brief example configuration file demonstrates common customization patterns; documentation is complete, accurate, and tested.

INPUT: fieldName, formatString, cfg format fields, data map
OUTPUT: warning strings, rendered message text
DATA: placeholderMap, regex extractors, known default placeholders

## VALIDATE_FORMAT_STRING

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::VALIDATE_FORMAT_STRING
STEP T001: COMPARE extracted placeholders to expected AND warn on unknown tokens
STEP T002: EMIT helpful validation error messages when placeholders are invalid

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: compare extracted placeholders to getExpectedPlaceholders and emit non-fatal warnings for unknown tokens.

PROCEDURE VALIDATE_FORMAT_STRING(fieldName, formatString):
  expected = GET_EXPECTED_PLACEHOLDERS(fieldName)
  IF expected empty THEN RETURN nil
  found = EXTRACT_PLACEHOLDERS(formatString)
  FOR EACH ph IN found NOT IN expected: APPEND warning
  RETURN warnings

## GET_EXPECTED_PLACEHOLDERS

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::GET_EXPECTED_PLACEHOLDERS
STEP T001: RETURN allowed placeholders from placeholderMap for fieldName

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: return allowed printf or template placeholders per Config field name from placeholderMap.

PROCEDURE GET_EXPECTED_PLACEHOLDERS(fieldName):
  RETURN placeholderMap[fieldName]  # e.g. FormatListArchive allows #{path}, #{size_human}, ...

## EXTRACT_PLACEHOLDERS

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::EXTRACT_PLACEHOLDERS
STEP T001: COLLECT printf verbs AND #{name} template tokens via regex

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: regex-collect %[sdvfbtxX] verbs and #{name} template tokens from a format string.

PROCEDURE EXTRACT_PLACEHOLDERS(formatString):
  RETURN printfMatches + templateMatches

## VALIDATE_ALL_FORMAT_STRINGS

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::VALIDATE_ALL_FORMAT_STRINGS
STEP T001: ITERATE Format* and Template* fields AND aggregate warnings

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: iterate all Format* and Template* cfg fields and aggregate ValidateFormatString warnings.

PROCEDURE VALIDATE_ALL_FORMAT_STRINGS(cfg):
  warnings = []
  FOR EACH (fieldName, value) IN cfg format fields:
    APPEND warnings FROM VALIDATE_FORMAT_STRING(fieldName, value)
  RETURN warnings

## FORMAT_TEMPLATE

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::FORMAT_TEMPLATE
STEP T001: REPLACE #{key} from data AND apply known defaults
STEP T002: EXECUTE Go text/template WHEN {{.}} patterns remain

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: replace #{key} from data, apply known defaults, handle residual %s with path/time, strip remaining #{...} before optional Go template pass.

PROCEDURE FORMAT_TEMPLATE(templateStr, data):
  result = REPLACE ALL #{key} FROM data
  FILL known defaults for #{size_human}, #{size}, #{mtime}
  IF %s remains: substitute path/name then creation_time
  REPLACE any leftover #{...} with data or "unknown"
  IF Go {{.}} patterns AND no #{ left: execute text/template ELSE return result

## LOAD_CONFIG_VALIDATION_HOOK

SPEC-ID: IMPL-CUSTOMIZABLE_FORMAT_STRINGS::LOAD_CONFIG_VALIDATION_HOOK
STEP T001: PRINT validateAllFormatStrings warnings to stderr without failing load
STEP T002: VERIFY example configuration files work as documented when loaded

- [IMPL-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-CUSTOMIZABLE_FORMAT_STRINGS] [ARCH-OUTPUT_FORMATTING] [REQ-CUSTOMIZABLE_FORMAT_STRINGS] — How: after YAML merge print validateAllFormatStrings warnings to stderr without failing load.

PROCEDURE LOAD_CONFIG_VALIDATION_HOOK(cfg):
  FOR w IN VALIDATE_ALL_FORMAT_STRINGS(cfg): PRINT w to stderr
