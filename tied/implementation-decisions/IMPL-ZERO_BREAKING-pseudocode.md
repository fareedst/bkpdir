# [IMPL-ZERO_BREAKING] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Ensure extraction steps preserve existing application behavior by keeping tests green, maintaining coverage, and verifying behavior at each step.

INPUT: extraction step, existing test suite, behavioral baselines
OUTPUT: unchanged user-facing behavior after each extraction increment
DATA: test results, coverage reports, verification checkpoints

## PRESERVE_EXISTING_TESTS

SPEC-ID: IMPL-ZERO_BREAKING::PRESERVE_EXISTING_TESTS
STEP T001: run the full existing test suite after each extraction step without modifying test expectations

- [IMPL-ZERO_BREAKING] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: run the full existing test suite after each extraction step without modifying test expectations.

PROCEDURE preserve_existing_tests():
  BEFORE each extraction step:
    RUN full existing test suite and record pass status
  AFTER each extraction step:
    RE-RUN same tests without changing test code or expected outcomes
    ASSERT all previously passing tests still pass

## MAINTAIN_COMPREHENSIVE_COVERAGE

SPEC-ID: IMPL-ZERO_BREAKING::MAINTAIN_COMPREHENSIVE_COVERAGE
STEP T001: measure coverage before and after extraction and add tests when coverage drops on moved code

- [IMPL-ZERO_BREAKING] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: measure coverage before and after extraction and add tests when coverage drops on moved code.

PROCEDURE maintain_comprehensive_coverage():
  RECORD coverage baseline before extraction
  AFTER each move:
    MEASURE coverage on affected packages
    IF coverage drops on public surfaces THEN
      ADD tests until baseline is restored
    END IF
  END FOR

## VERIFY_BEHAVIOR_AT_EACH_STEP

SPEC-ID: IMPL-ZERO_BREAKING::VERIFY_BEHAVIOR_AT_EACH_STEP
STEP T001: compare observable CLI and API behavior against baselines before proceeding to the next extraction step

- [IMPL-ZERO_BREAKING] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: compare observable CLI and API behavior against baselines before proceeding to the next extraction step.

PROCEDURE verify_behavior_at_each_step():
  FOR each extraction increment:
    RUN behavioral smoke checks on backup, list, and config flows
    COMPARE outputs to pre-extraction baselines
    HALT extraction if any user-visible behavior differs
  END FOR
