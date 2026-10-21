# [IMPL-EXTRACTION_CHALLENGES] [ARCH-CODE_ORGANIZATION] [ARCH-CONFIG_SYSTEM] [ARCH-TESTING_STRATEGY] [REQ-CONFIGURATION] [REQ-MAINTAINABILITY]

## Summary contract

Address the hardest extraction problems: shared test utilities, per-package test isolation, and config coupling via minimal interfaces.

INPUT: shared test helpers, extracted packages, config structs
OUTPUT: testutil package, package-scoped tests, ConfigProvider adapters
DATA: fixtures, mocks, interface method sets

## EXTRACT_TEST_UTILITIES

SPEC-ID: IMPL-EXTRACTION_CHALLENGES::EXTRACT_TEST_UTILITIES
STEP T001: move shared builders, fixtures, and assertions into testutil before splitting production packages

- [IMPL-EXTRACTION_CHALLENGES] [ARCH-CODE_ORGANIZATION] [ARCH-TESTING_STRATEGY] [REQ-MAINTAINABILITY] — How: move shared builders, fixtures, and assertions into testutil before splitting production packages.

PROCEDURE extract_test_utilities():
  IDENTIFY shared test helpers used across multiple test files
  MOVE helpers to testutil package or package-internal test helpers
  UPDATE dependency references in all consuming test files
  VALIDATE all existing tests still pass after helper extraction

## CREATE_PACKAGE_SPECIFIC_TESTS

SPEC-ID: IMPL-EXTRACTION_CHALLENGES::CREATE_PACKAGE_SPECIFIC_TESTS
STEP T001: give each extracted package focused unit tests with mocked dependencies at its public API

- [IMPL-EXTRACTION_CHALLENGES] [ARCH-CODE_ORGANIZATION] [ARCH-TESTING_STRATEGY] [REQ-MAINTAINABILITY] — How: give each extracted package focused unit tests with mocked dependencies at its public API.

PROCEDURE create_package_specific_tests():
  FOR each extracted_package:
    WRITE unit tests for public API surface
    WRITE edge-case tests for error paths
    MOCK dependencies using interfaces, not concrete types
    VALIDATE package tests pass in isolation without root-level imports
  END FOR

## RESOLVE_CONFIG_COUPLING

SPEC-ID: IMPL-EXTRACTION_CHALLENGES::RESOLVE_CONFIG_COUPLING
STEP T001: replace direct config struct imports with per-package ConfigProvider interfaces and root adapters

- [IMPL-EXTRACTION_CHALLENGES] [ARCH-CODE_ORGANIZATION] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: replace direct config struct imports with per-package ConfigProvider interfaces and root adapters.

PROCEDURE resolve_config_coupling():
  FOR each extracted package:
    DEFINE ConfigProvider interface with minimal method set for that package
    IMPLEMENT adapter in root package that wraps full Config struct
    ENSURE extracted package depends only on ConfigProvider, not root config
  END FOR
  VALIDATE no extracted package depends on root config directly
