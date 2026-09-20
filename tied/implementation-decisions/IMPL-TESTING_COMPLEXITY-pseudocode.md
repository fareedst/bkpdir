# [IMPL-TESTING_COMPLEXITY] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY]

## Summary contract

Testing strategy for large refactors: extract shared test utilities, add package-scoped tests, and keep legacy integration tests passing.

INPUT: monolithic test helpers and package boundaries
OUTPUT: pkg/*_test.go coverage plus unchanged root integration suite
DATA: testutil fixtures, interface mocks

## TESTING_COMPLEXITY_CHALLENGE

SPEC-ID: IMPL-TESTING_COMPLEXITY::TESTING_COMPLEXITY_CHALLENGE
STEP T001: extract shared test utilities first, add package-focused tests, and keep root integration tests green after package extraction

- [IMPL-TESTING_COMPLEXITY] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: extract shared test utilities first, add package-focused tests, and keep root integration tests green after package extraction.

PROCEDURE testing_complexity_challenge():
  EXTRACT test utilities to pkg/testutil
  ADD focused tests per extracted package
  RUN full integration suite to confirm no regressions
