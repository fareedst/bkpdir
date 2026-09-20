# [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016]

## Summary contract

Audit modules for expected REQ/ARCH/IMPL token annotations and remediate gaps across the codebase.

INPUT: module file lists, expected token sets per module
OUTPUT: coverage report, added annotations in source and tests
DATA: scanned token comments, gap records, cross-module aggregates

## AUDIT_MODULE_COVERAGE

SPEC-ID: IMPL-TOKEN_COVERAGE_AUDIT::AUDIT_MODULE_COVERAGE
STEP T001: Scan each module file for REQ/ARCH/IMPL comments and record missing expected tokens

- [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] — How: scan each module file for REQ/ARCH/IMPL comments and record missing expected tokens.

PROCEDURE AUDIT_MODULE_COVERAGE(module):
  FOR EACH file IN module.files:
    SCAN for [REQ-*], [ARCH-*], [IMPL-*] comments
    COMPARE found tokens TO module.expected_tokens
    IF gaps THEN record missing tokens
  IF gaps THEN remediate by adding annotations to source and tests

## RUN_FULL_AUDIT (DOC-ONLY NOTE)

SPEC-ID: IMPL-TOKEN_COVERAGE_AUDIT::RUN_FULL_AUDIT
STEP T001: EXECUTE block PROCEDURE body

- [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] — How: aggregate module coverage scans into a cross-module audit report and remediation plan.

NOTE: Cross-module aggregation is currently tracked as documentation/process guidance; no dedicated Go `// - [IMPL-TOKEN_COVERAGE_AUDIT]` lead exists yet for an end-to-end runner.

## HYGIENE_POLICY (lead placement)

SPEC-ID: IMPL-TOKEN_COVERAGE_AUDIT::HYGIENE_POLICY
STEP T001: Enforce module-scoped literal block leads; forbid repo-wide test paste

- [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: place `// - {exact sidecar lead}` immediately before the implementing `func` or covering `Test*`/`Benchmark*` only; never paste full sidecar lead lists after `package` in `*_test.go`; disable bulk `add-test-leads --all`; run `audit_package_level_leads.py --fail-on-suspect` and D15 `strip_redundant_impl_banners.py` to drop paraphrased `// [IMPL-*]` banners when literal block leads for the same IMPL token exist in the same file.
