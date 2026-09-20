# [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION]

## Summary contract

Match paths against configured exclusion patterns using directory suffix, glob, and exact rules with normalized forward-slash paths.

INPUT: path string, pattern list
OUTPUT: boolean excluded
DATA: PatternMatcher.patterns, doublestar globs

## NEW_PATTERN_MATCHER

SPEC-ID: IMPL-EXCLUSION_PATTERNS::NEW_PATTERN_MATCHER
STEP T001: RETURN PatternMatcher storing pattern slice

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: store pattern slice on PatternMatcher for iterative ShouldExclude checks.

PROCEDURE NEW_PATTERN_MATCHER(patterns):
  RETURN &PatternMatcher{patterns: patterns}

## SHOULD_EXCLUDE

SPEC-ID: IMPL-EXCLUSION_PATTERNS::SHOULD_EXCLUDE
STEP T001: NORMALIZE path AND return true on first pattern match

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: normalize path to slashes, return true on first pattern match via matchesPattern dispatch.

PROCEDURE SHOULD_EXCLUDE(path):
  FOR EACH pattern: IF matchesPattern(path, pattern): RETURN true
  RETURN false

## MATCHES_PATTERN

SPEC-ID: IMPL-EXCLUSION_PATTERNS::MATCHES_PATTERN
STEP T001: DISPATCH directory suffix, glob, or exact path matching

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: route trailing-/ to directory rules, glob * to doublestar, else exact path equality.

PROCEDURE MATCHES_PATTERN(path, pattern):
  IF pattern ends with /: directory match WITH **/ support
  IF CONTAINS *: glob match
  RETURN path == pattern

## SHOULD_EXCLUDE_FILE

SPEC-ID: IMPL-EXCLUSION_PATTERNS::SHOULD_EXCLUDE_FILE
STEP T001: CONSTRUCT matcher AND call ShouldExclude once

- [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: convenience wrapper constructing matcher and calling ShouldExclude once.

PROCEDURE SHOULD_EXCLUDE_FILE(path, patterns):
  RETURN NewPatternMatcher(patterns).ShouldExclude(path)
