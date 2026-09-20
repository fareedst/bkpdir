# [IMPL-TEST_EMPTY_STRING_HANDLING] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_001] [REQ-CFG_005]

## Summary contract

Verifies empty-string configuration values merge correctly across sequential files without blocking later overrides.

INPUT: sequential YAML files with archive_dir_path "" or values
OUTPUT: merged Config from LoadConfig
DATA: BKPDIR_CONFIG file list, precedence fields

## TESTEMPTYSTRINGHANDLING

SPEC-ID: IMPL-TEST_EMPTY_STRING_HANDLING::TESTEMPTYSTRINGHANDLING
STEP T001: verify empty-string archive_dir_path merges with later files and later non-empty values override per CFG-001 precedence

- [IMPL-TEST_EMPTY_STRING_HANDLING] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_001] [REQ-CFG_005] — How: verify empty-string archive_dir_path merges with later files and later non-empty values override per CFG-001 precedence.

PROCEDURE TestEmptyStringHandling():
  subtest first file empty second file value -> expect second value
  subtest both empty -> expect empty string retained
  subtest empty then non-empty override paths
