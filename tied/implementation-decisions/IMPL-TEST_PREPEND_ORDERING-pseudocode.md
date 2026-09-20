# [IMPL-TEST_PREPEND_ORDERING] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_005] [REQ-CFG_001]

## Summary contract

Tests ^prepend merge strategy ordering for exclude_patterns in inheritance and sequential file chains.

INPUT: parent/child YAML with ^exclude_patterns and multi-file BKPDIR_CONFIG
OUTPUT: ordered ExcludePatterns after LoadConfig / LoadConfigWithInheritance
DATA: prepend before destination slice semantics

## TESTPREPENDSTRATEGYORDERING

SPEC-ID: IMPL-TEST_PREPEND_ORDERING::TESTPREPENDSTRATEGYORDERING
STEP T001: verify ^prepend places source exclude_patterns before destination patterns in inheritance and sequential chains preserving order

- [IMPL-TEST_PREPEND_ORDERING] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-CFG_005] [REQ-CFG_001] — How: verify ^prepend places source exclude_patterns before destination patterns in inheritance and sequential chains preserving order.

PROCEDURE TestPrependStrategyOrdering():
  subtest inheritance parent then child prepend
  subtest sequential multiple prepend operations
  ASSERT source segments appear before earlier destination segments
