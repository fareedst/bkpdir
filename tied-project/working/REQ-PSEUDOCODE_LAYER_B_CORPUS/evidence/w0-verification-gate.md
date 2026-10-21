# W0 verification-gate evidence (2026-10-21)

**Run ID:** `layer-b-corpus-20261021-ver`  
**Commit:** `8fcb9215c9b55a5ef2a84d8af80718ebf6ef812c`

## Layer B (minimum_gating_rules)

`run-pseudocode-layer-b-batch.sh --mode full --wave w0 --enforce` → **6/6 ok**  
Reports: `pseudocode-validate-full/*.json`

## Layer C

`scripts/run_layer_c_gate.sh` → **6/6 PASS** (includes IMPL-EXCLUSION_PATTERNS)

## Executable tests

| Command | Result |
|---------|--------|
| `go test ./test/specconformance/... ./internal/specmodel/... -count=1` | pass |
| `go test ./internal/specparse/... -count=1` | pass |

## L0 / traceability (W0 subset)

| Command | Scope | Result |
|---------|-------|--------|
| `specctl validate` | W0 six sidecars | pass |
| `specctl check-leads` | W0 six sidecars | pass |
| `run-spec-verification.sh` | full formal registry (73) | **fail** TRACE-002 on non-W0 IMPLs (DIFF_COMMAND, DIRECTORY_COMPARISON, INCREMENTAL_DUPLICATE_PREVENTION) — pre-existing corpus debt; **out of W0 wave scope** |

## Adversarial inquiry (verification)

Artifacts: `adversarial-inquiry/phase-verification/` (run_id `layer-b-corpus-20261021-ver`, policy advisory)

## TIED consistency

`tied_validate_consistency` → index valid (2026-10-21 session)
