# [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] [REQ-OUTPUT_FORMATTING]

## Summary contract

Comprehensive formatter test coverage for OutputCollector, template/printf formatting, error templates, file statistics output, and interface compliance.

INPUT: formatter config, test archives, template strings, collector instances
OUTPUT: formatted strings, captured stdout/stderr messages
DATA: DefaultOutputFormatter, OutputCollector, TemplateFormatter test fixtures

## TEST_OUTPUT_COLLECTOR

SPEC-ID: IMPL-TEST_COVERAGE::TEST_OUTPUT_COLLECTOR
STEP T001: Verify AddStdout/AddStderr/Clear/GetMessages track delayed output by destination

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] — How: verify AddStdout/AddStderr/Clear/GetMessages track delayed output by destination.

PROCEDURE TEST_OUTPUT_COLLECTOR():
  collector = NewOutputCollector()
  collector.AddStdout("info line", "info")
  collector.AddStderr("error line", "error")
  ASSERT GetMessages returns both with correct destination
  collector.Clear(); ASSERT empty

## TEST_DELAYED_OUTPUT_MODE

SPEC-ID: IMPL-TEST_COVERAGE::TEST_DELAYED_OUTPUT_MODE
STEP T001: Verify formatter with collector buffers PrintCreatedArchive and PrintError until flush

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] — How: verify formatter with collector buffers PrintCreatedArchive and PrintError until flush.

PROCEDURE TEST_DELAYED_OUTPUT_MODE():
  formatter = NewDefaultOutputFormatterWithCollector(cfg, collector)
  ASSERT IsDelayedMode true
  formatter.PrintCreatedArchive(path); formatter.PrintError(err)
  ASSERT collector holds stdout and stderr messages

## TEST_TEMPLATE_FORMATTER

SPEC-ID: IMPL-TEST_COVERAGE::TEST_TEMPLATE_FORMATTER
STEP T001: Validate #{placeholder} replacement, list extraction, and mixed %s/#{...} format strings

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] [REQ-OUTPUT_FORMATTING] — How: validate #{placeholder} replacement, list extraction, and mixed %s/#{...} format strings.

PROCEDURE TEST_TEMPLATE_FORMATTER():
  FormatListArchiveWithExtraction returns filename in output
  FormatWithPlaceholders replaces #{path} and #{branch}
  FormatListArchive leaves no unprocessed #{...} patterns

## TEST_ERROR_AND_DEFAULT_FORMATTER

SPEC-ID: IMPL-TEST_COVERAGE::TEST_ERROR_AND_DEFAULT_FORMATTER
STEP T001: Verify default format strings and error template methods include path and error type labels

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] — How: verify default format strings and error template methods include path and error type labels.

PROCEDURE TEST_ERROR_AND_DEFAULT_FORMATTER():
  FormatCreatedArchive contains path and label
  FormatDiskFullError and permission errors include original err text
  custom cfg format strings override defaults

## TEST_FILE_STATS_AND_COMPAT

SPEC-ID: IMPL-TEST_COVERAGE::TEST_FILE_STATS_AND_COMPAT
STEP T001: Assert FormatCreatedArchiveWithStats and Print*WithStats include size placeholders and basic methods remain compatible

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] [REQ-OUTPUT_FORMATTING] — How: assert FormatCreatedArchiveWithStats and Print*WithStats include size placeholders and basic methods remain compatible.

PROCEDURE TEST_FILE_STATS_AND_COMPAT():
  FormatCreatedArchiveWithStats includes path and human size
  PrintCreatedArchiveWithStats writes expected stdout
  FormatCreatedArchive basic path still works alongside stats methods

## TEST_INTERFACE_COMPLIANCE

SPEC-ID: IMPL-TEST_COVERAGE::TEST_INTERFACE_COMPLIANCE
STEP T001: Assert DefaultOutputFormatter, DefaultTemplateFormatter, and DefaultPatternExtractor satisfy their interfaces

- [IMPL-TEST_COVERAGE] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY] — How: assert DefaultOutputFormatter, DefaultTemplateFormatter, and DefaultPatternExtractor satisfy their interfaces.

PROCEDURE TEST_INTERFACE_COMPLIANCE():
  ASSERT types implement OutputFormatterInterface
  ASSERT types implement TemplateFormatterInterface
  ASSERT types implement PatternExtractorInterface
