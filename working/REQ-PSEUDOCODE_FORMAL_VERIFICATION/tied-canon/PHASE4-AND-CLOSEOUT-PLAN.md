# Phase 4 and close-out (refined 2026-09-20)

Canonical plan: Cursor plan **`pseudocode_tied_canon_aa64174a.plan.md`** (Phase 4 + close-out sections). This file is a working-slice extract.

## Phase 4A — Agent canon (docs)

1. `AGENTS.md` — sidecar-only behavior; `sync_tied_yaml_projections.py`
2. `tied/docs/impl-deep-sync-agent-guide.md` — § Sidecar-canonical YAML projection
3. Cross-links in `docs/markscope/spec-verification.md` and client development index

## Phase 4B — Mutation wave 4

```bash
MUTATION_WAVE=4 scripts/run-mutation-pilot.sh
# evidence → working/.../tied-canon/evidence/mutation-wave4-*.json
```

## Phase 4C — CRIT-001 strict (optional archive)

```bash
python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-strict
```

## Close-out

1. `go test ./... -count=1`
2. Verification gate + `run-close-out-gates.mjs --envelope-blocking --sync-dispositions --reconcile`
3. Close-out gate; persist CITDP to `tied/citdp/`
