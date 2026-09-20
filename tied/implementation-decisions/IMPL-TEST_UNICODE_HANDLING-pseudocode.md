# [IMPL-TEST_UNICODE_HANDLING] [ARCH-TESTING_STRATEGY] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Ensures Unicode and special characters survive YAML load, merge with defaults, and config paths containing spaces or non-ASCII directory names.

INPUT: UTF-8 config values and filesystem paths
OUTPUT: LoadConfig results preserving characters
DATA: archive_dir_path, exclude_patterns, temp directories

## TESTUNICODEHANDLING

SPEC-ID: IMPL-TEST_UNICODE_HANDLING::TESTUNICODEHANDLING
STEP T001: load config with Unicode archive paths and exclude_patterns and assert values round-trip unchanged after merge with defaults

- [IMPL-TEST_UNICODE_HANDLING] [ARCH-TESTING_STRATEGY] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: load config with Unicode archive paths and exclude_patterns and assert values round-trip unchanged after merge with defaults.

PROCEDURE TestUnicodeHandling():
  WRITE .bkpdir.yml with emoji and CJK paths and patterns
  LoadConfig
  ASSERT ArchiveDirPath and each ExcludePatterns entry unchanged

## TESTSPECIALCHARACTERSINPATHS

SPEC-ID: IMPL-TEST_UNICODE_HANDLING::TESTSPECIALCHARACTERSINPATHS
STEP T001: load config from directories whose paths contain spaces or Unicode without corrupting merged field values

- [IMPL-TEST_UNICODE_HANDLING] [ARCH-TESTING_STRATEGY] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: load config from directories whose paths contain spaces or Unicode without corrupting merged field values.

PROCEDURE TestSpecialCharactersInPaths():
  subtests with spaced paths, Unicode directory names, and combined
  LoadConfig from each directory layout
  ASSERT expected archive_dir_path and exclude_patterns
