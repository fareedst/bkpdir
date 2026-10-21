# [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY]

## Summary contract

Defines testing strategy: table-driven unit tests, shared testutil helpers, internal scenario utilities, and quality gates for bkpdir.

INPUT: test cases, temp paths, fixtures, scenario configs
OUTPUT: assertions, temp artifacts, integration scenario results
DATA: pkg/testutil helpers, internal/testutil frameworks, make test targets

## TABLE_DRIVEN_UNIT_TESTS

SPEC-ID: IMPL-TESTING::TABLE_DRIVEN_UNIT_TESTS
STEP T001: Express unit tests as struct slices with name/input/want fields and t.Run subtests per case

- [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: express unit tests as struct slices with name/input/want fields and t.Run subtests per case.

PROCEDURE TABLE_DRIVEN_UNIT_TESTS():
  FOR EACH test case IN cases:
    t.Run(case.name, func(t):
      SETUP inputs from case
      result = EXECUTE function under test
      ASSERT result matches case.want

## TESTUTIL_FILESYSTEM

SPEC-ID: IMPL-TESTING::TESTUTIL_FILESYSTEM
STEP T001: Provide CreateTempDir and CreateTempFile with t.Cleanup registration for isolated filesystem tests

- [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: provide CreateTempDir and CreateTempFile with t.Cleanup registration for isolated filesystem tests.

PROCEDURE TESTUTIL_FILESYSTEM():
  CreateTempDir(t, prefix) -> MkdirTemp + t.Cleanup RemoveAll
  CreateTempFile(t, dir, name, content) -> WriteFile + t.Cleanup Remove

## TESTUTIL_FIXTURES

SPEC-ID: IMPL-TESTING::TESTUTIL_FIXTURES
STEP T001: Build reusable git repos, archives, backups, and config fixtures for integration tests

- [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: build reusable git repos, archives, backups, and config fixtures for integration tests.

PROCEDURE TESTUTIL_FIXTURES():
  CreateTestGitRepo(t, dir) -> init repo with commits
  CreateTestArchive(t, dir, files) -> zip fixture
  CreateTestBackup(t, path) -> timestamped backup file
  SetupTestConfig(t, overrides) -> temp config YAML

## TESTUTIL_ASSERTIONS

SPEC-ID: IMPL-TESTING::TESTUTIL_ASSERTIONS
STEP T001: Centralize AssertNoTempFiles, AssertArchiveContents, and output capture helpers for tests

- [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: centralize AssertNoTempFiles, AssertArchiveContents, and output capture helpers for tests.

PROCEDURE TESTUTIL_ASSERTIONS():
  AssertNoTempFiles(t, dir) -> verify cleanup
  AssertArchiveContents(t, zipPath, expected) -> read zip entries
  CaptureOutput(t, fn) -> redirect stdout/stderr during fn

## INTERNAL_SCENARIO_HELPERS

SPEC-ID: IMPL-TESTING::INTERNAL_SCENARIO_HELPERS
STEP T001: Simulate disk space, permissions, corruption, and error injection in internal/testutil for stress paths

- [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: simulate disk space, permissions, corruption, and error injection in internal/testutil for stress paths.

PROCEDURE INTERNAL_SCENARIO_HELPERS():
  SimulateDiskFull(t, path)
  SimulatePermissionDenied(t, path)
  InjectArchiveCorruption(t, zipPath)
  RunWithErrorInjection(t, scenario)
