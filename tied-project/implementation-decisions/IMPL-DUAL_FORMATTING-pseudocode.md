# [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING]

## Summary contract

OutputFormatter supports printf cfg strings and template #{key} strings with guarded selection, shared formatTemplate helper, and interface contracts for format providers and destinations.

INPUT: format string, path, creation time, data map
OUTPUT: rendered message string
DATA: FormatProvider, TemplateFormatter, OutputCollector routing

## FORMAT_CREATED_ARCHIVE_PRINTF

SPEC-ID: IMPL-DUAL_FORMATTING::FORMAT_CREATED_ARCHIVE_PRINTF
STEP T001: SPRINTF cfg.FormatCreatedArchive WITH path

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: sprintf cfg.FormatCreatedArchive with path for simple archive-created messages.

PROCEDURE FORMAT_CREATED_ARCHIVE_PRINTF(path):
  RETURN sprintf(cfg.FormatCreatedArchive, path)

## FORMAT_LIST_ARCHIVE_DUAL_MODE

SPEC-ID: IMPL-DUAL_FORMATTING::FORMAT_LIST_ARCHIVE_DUAL_MODE
STEP T001: IF format contains #{ THEN formatTemplate WITH stats
STEP T002: ELIF contains % THEN sprintf WITH path and creationTime

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: if format contains #{ use formatTemplate with GatherFileStatInfo data; elif contains % use sprintf; else return literal format string.

PROCEDURE FORMAT_LIST_ARCHIVE_DUAL_MODE(path, creationTime):
  IF CONTAINS(formatStr, "#{"): BUILD data map WITH stats; RETURN formatTemplate
  IF CONTAINS(formatStr, "%"): RETURN sprintf WITH path and creationTime
  RETURN formatStr unchanged

## FORMAT_CONFIG_VALUE_PRINTF

SPEC-ID: IMPL-DUAL_FORMATTING::FORMAT_CONFIG_VALUE_PRINTF
STEP T001: SPRINTF cfg.FormatConfigValue WITH name, value, source

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: sprintf cfg.FormatConfigValue with name, value, and source fields.

PROCEDURE FORMAT_CONFIG_VALUE_PRINTF(name, value, source):
  RETURN sprintf(cfg.FormatConfigValue, name, value, source)

## FORMAT_PROVIDER_INTERFACES

SPEC-ID: IMPL-DUAL_FORMATTING::FORMAT_PROVIDER_INTERFACES
STEP T001: DECLARE FormatProvider, OutputDestination, PatternExtractor contracts
STEP T002: DECLARE FormatterInterface AND TemplateFormatterInterface

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: declare FormatProvider, OutputDestination, PatternExtractor, FormatterInterface, and TemplateFormatterInterface contracts for extraction and delayed output.

PROCEDURE FORMAT_PROVIDER_INTERFACES():
  DEFINE interfaces for format/template/pattern operations and IsDelayedMode/SetCollector

## TEMPLATE_FORMATTER_PLACEHOLDERS

SPEC-ID: IMPL-DUAL_FORMATTING::TEMPLATE_FORMATTER_PLACEHOLDERS
STEP T001: REPLACE #{key} FROM data WITH known defaults
STEP T002: HANDLE residual %s AND optional Go text/template

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: TemplateFormatter.FormatWithPlaceholders replaces #{key} from data with defaults for missing stat fields.

PROCEDURE TEMPLATE_FORMATTER_PLACEHOLDERS(format, data):
  REPLACE #{key}; APPLY known defaults; HANDLE residual %s; optional Go text/template when no #{ remain

## FORMAT_WITH_TEMPLATE

SPEC-ID: IMPL-DUAL_FORMATTING::FORMAT_WITH_TEMPLATE
STEP T001: REGEX_FIND_SUBMATCH pattern ON input
STEP T002: DELEGATE to FormatWithPlaceholders on template string

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: compile regex pattern, extract named submatches into data map, delegate to FormatWithPlaceholders on template string.

PROCEDURE FORMAT_WITH_TEMPLATE(input, pattern, tmplStr):
  matches = REGEX_FIND_SUBMATCH(pattern, input)
  BUILD data FROM named groups
  RETURN FormatWithPlaceholders(tmplStr, data)

## EXTRACT_PATTERN_DATA

SPEC-ID: IMPL-DUAL_FORMATTING::EXTRACT_PATTERN_DATA
STEP T001: COMPILE configured regex ON input text
STEP T002: STORE named capture groups IN result map

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: compile configured regex and return map of named capture groups from input text.

PROCEDURE EXTRACT_PATTERN_DATA(pattern, text):
  IF compile fails THEN RETURN empty map
  FOR EACH named submatch: STORE in result map
  RETURN result

## SIMPLE_PLACEHOLDER_REPLACE

SPEC-ID: IMPL-DUAL_FORMATTING::SIMPLE_PLACEHOLDER_REPLACE
STEP T001: FOR EACH key IN data REPLACE #{key} WITH value

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: replace every #{key} in formatStr with data map values; leave unknown placeholders unchanged.

PROCEDURE SIMPLE_PLACEHOLDER_REPLACE(formatStr, data):
  FOR EACH key IN data: REPLACE #{key} WITH value
  RETURN result

## PRINT_WITH_DELAYED_ROUTING

SPEC-ID: IMPL-DUAL_FORMATTING::PRINT_WITH_DELAYED_ROUTING
STEP T001: FORMAT message THEN route via collector OR immediate print

- [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: Print* methods format message then AddStdout when collector attached else immediate print.

PROCEDURE PRINT_WITH_DELAYED_ROUTING(message):
  IF collector: AddStdout ELSE immediate print
