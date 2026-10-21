# [IMPL-TEST_DEFAULT_STRATEGY_EDGES] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_005]

## Summary contract

Tests that =default merge strategy applies only when the destination field is the zero value (empty string, empty slice, false bool), not when the field equals the built-in default.

INPUT: paired YAML config files via BKPDIR_CONFIG
OUTPUT: LoadConfig results and assertions on ArchiveDirPath, ExcludePatterns, IncludeGitInfo
DATA: =prefixed keys, default merge behavior for accumulate fields

## TESTDEFAULTSTRATEGYEDGECASES

SPEC-ID: IMPL-TEST_DEFAULT_STRATEGY_EDGES::TESTDEFAULTSTRATEGYEDGECASES
STEP T001: assert =default strategy applies only when destination field is zero value for strings, arrays, and bools across sequential files

- [IMPL-TEST_DEFAULT_STRATEGY_EDGES] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_005] — How: assert =default strategy applies only when destination field is zero value for strings, arrays, and bools across sequential files.

PROCEDURE TestDefaultStrategyEdgeCases():
  FOR subtests covering zero vs non-zero strings, arrays, and bools:
    CREATE file1 and file2 with =field defaults
    LoadConfig via BKPDIR_CONFIG=file1:file2
    ASSERT default applied only when first file left field at zero value
