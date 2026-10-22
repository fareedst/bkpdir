# W3 verification-gate evidence (2026-10-21)

**Run ID:** `layer-b-corpus-w3-ver`  
**Commit (pre-slice):** `98c5a3ebbcc503fc14a8d883573189f25eb912f2`  
**Wave:** `w3` — 27 Tier C `lead_check_skip` IMPLs per `tied-project/spec/layer-b-wave-registry.yaml`

## Layer B (minimum_gating_rules)

```bash
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode full --wave w3 --enforce \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-full
```

**Result:** **27/27 ok**

## L0 regression

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

**Result:** pass (D15 redundant banner advisory on `pkg/fileops/comparison.go` — deferred hygiene)

## Executable tests

| Command | Result |
|---------|--------|
| COVER tests for W3 | **N/A** — `lead_check_skip` doc/process tier; waivers in `tied-project/spec/coverage-waivers.yaml` |

## TIED consistency

`tied_validate_consistency` → pass (2026-10-21 session)

## Layer C

Not required for W3 (structural SHAPE + meaningful How only).

## Adversarial inquiry

Mode B refresh — activation pairs with on-disk `run_id` **`layer-b-corpus-20261021-ver`** (W0/W1/W2 pattern). **Advisory** policy — gate receipt diagnostics `finding_unresolved` / `warn_not_success` non-blocking.

**Gate receipt:** `gates/verification-2026-10-22T02-29-46-257Z.json` → `allowed: true`

## Manifest

`tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/evidence/verification-evidence-manifest-w3.v1.json`

## Close-out

See [`w3-close-out-handoff.md`](w3-close-out-handoff.md) — pre-commit gates green (`close-out-gates-run-w3-precommit.json`).
