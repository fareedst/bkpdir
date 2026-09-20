# [IMPL-TEST_EXCLUDE_MERGE] [ARCH-TEST_EXCLUDE_MERGE] [REQ-TEST_EXCLUDE_MERGE] [REQ-CONFIGURATION] [REQ-CFG_006]

## Summary contract

Integration tests for exclude_patterns merging across sequential loading, inheritance chains, and config command source attribution.

INPUT: multi-file BKPDIR_CONFIG layouts and inherit directives
OUTPUT: merged ExcludePatterns and config display metadata
DATA: built-in defaults, merge vs override strategies

## TESTEXCLUDEPATTERNSMERGE

SPEC-ID: IMPL-TEST_EXCLUDE_MERGE::TESTEXCLUDEPATTERNSMERGE
STEP T001: verify exclude_patterns merge order across sequential files and inheritance with expected deduplicated slice and config source labels

- [IMPL-TEST_EXCLUDE_MERGE] [ARCH-TEST_EXCLUDE_MERGE] [REQ-TEST_EXCLUDE_MERGE] [REQ-CONFIGURATION] [REQ-CFG_006] — How: verify exclude_patterns merge order across sequential files and inheritance with expected deduplicated slice and config source labels.

PROCEDURE TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE():
  subtest first file replaces built-ins when unprefixed
  subtest sequential plus inheritance chain ordering
  ASSERT expected pattern order and source attribution
