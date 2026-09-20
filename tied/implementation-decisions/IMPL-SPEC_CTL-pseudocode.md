# [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]

## Summary contract

CLI and libraries that parse formal IMPL pseudocode, emit Layer B diagnostics, build traceability matrices, and orchestrate conformance checks.

INPUT: sidecar paths, repo root, formal_spec registry
OUTPUT: exit code 0 or diagnostic report
DATA: parsed AST per file, symbol table, coverage maps

## VALIDATE_COMMAND

SPEC-ID: IMPL-SPEC_CTL::VALIDATE_COMMAND
PRE: paths non-empty
STEP T001: FOR each path PARSE sidecar INTO document
STEP T002: RUN schema and symbol checks EMIT diagnostics
STEP T003: IF --check-req-criteria THEN EMIT CRIT-001 orphan criteria diagnostics
POST: exit code 0 IFF no error-severity findings

- [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: parse each file and fail on PARSE/SHAPE/RESOLVE errors; optional CRIT-001 via run_impl_logic_audit.py --req-criteria-strict; L4 mutation pilot and contract eval documented in scripts (not global CI).

PROCEDURE VALIDATE(paths):
  FOR path IN paths:
    doc = PARSE(path)
    EMIT diagnostics FROM validateDocument(doc)
  RETURN failure IF any error severity

## MATRIX_COMMAND

SPEC-ID: IMPL-SPEC_CTL::MATRIX_COMMAND
STEP T001: COLLECT SPEC-ID and STEP from all formal_spec documents
STEP T002: SCAN Go tests and production for trace references
STEP T003: IF check_coverage THEN FAIL on uncovered STEP without INFRA waiver
STEP T004: RUN mutation pilot scripts/run-mutation-pilot.sh per MUTATION_WAVE (L4 opt-in; score threshold on pilot packages)
STEP T005: EVAL PRE POST via internal/specmodel/contract for pilot oracle blocks (AssertBlockContracts)

- [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: build REQ/IMPL/STEP to test/code matrix; L4 mutation score on pilot packages via go-mutesting; property conformance tests in test/specconformance.

## CHECK_LEADS_COMMAND

SPEC-ID: IMPL-SPEC_CTL::CHECK_LEADS_COMMAND
STEP T001: EXTRACT block lead lines from each sidecar block
STEP T002: FIND matching literal comment in Go files under pkg/, internal/, cmd/, and test/ (each when present) plus repo-root top-level *.go
POST: every block lead has at least one Go comment match

- [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: verify [PROC-IMPL_PSEUDOCODE_TOKENS] literal copy in Go sources.
