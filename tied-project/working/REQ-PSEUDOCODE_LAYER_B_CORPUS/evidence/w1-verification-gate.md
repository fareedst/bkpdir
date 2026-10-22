# W1 verification-gate evidence (2026-10-21)

**Run ID:** `layer-b-corpus-w1-ver`  
**Commit:** `2e5986df600d5001a58d4d4abf62ee175a5a620f`  
**Wave:** `w1` — 38 oracle/Tier-A runtime IMPLs per [`layer-b-wave-registry.yaml`](../../../spec/layer-b-wave-registry.yaml)

## Layer B (minimum_gating_rules)

```bash
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh --mode full --wave w1 --enforce
```

**Result:** **38/38 ok** — reports under `pseudocode-validate/{IMPL-TOKEN}.json`

Mechanical pass: `python3 scripts/enrich_tier_a_procedure_contracts.py --wave w1 --apply` (procedure `Contract:` stanzas, `# [IMPL-*]` token links, TERMINATION/FAILURE_MODES where triggered) plus nine manual FAILURE_MODES rows for inter-procedure validator bleed.

## L0 regression (full formal registry)

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

**Result:** pass (zero TRACE-002)

## Executable tests

| Command | Result |
|---------|--------|
| `go test ./test/specconformance/... ./internal/specmodel/... ./internal/specparse/... -count=1` | pass |

## TIED consistency

`tied_validate_consistency` (include_pseudocode) → index valid (2026-10-21 session)

## Adversarial inquiry

Reuse W0 integrated verification artifacts (`adversarial-inquiry/phase-verification/`, run_id `layer-b-corpus-20261021-ver`) under advisory policy; W1 sidecar edits are contract-precision only (no runtime LEAP in this wave).

## Deferred

- **W2** (2 tokens) and **W3** (27 tokens) in wave registry
- Layer C `gate_mode` expansion beyond W0 six tokens (optional follow-up)
