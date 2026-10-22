# Continuation handoff — W2 corpus (close-out complete)

**Parent slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS`  
**Prior close-out:** W1 — [`close-out-handoff.md`](close-out-handoff.md) (commit `db5fb43913c8d615da7a12c84275bf4d86a87e95` on branch at session start)  
**Wave map:** [`w1-w3-corpus-handoff.md`](w1-w3-corpus-handoff.md)  
**Tracker:** `operator_evidence.status: w2_slice_close_out` — [`../agent-req-implementation-checklist.yaml`](../agent-req-implementation-checklist.yaml)

**Close-out:** [`w2-close-out-handoff.md`](w2-close-out-handoff.md) (feature commit `f8608c1…`)

## W2 scope (2 IMPLs)

| Token | Role |
|-------|------|
| `IMPL-CONFIGURABLE_STRINGS` | Tier B formal — error format strings, `GetErrorFormatStrings`, `HandleArchiveError` |
| `IMPL-CONFIG_DISPLAY_FLATTENING` | Tier B formal — flat config display, `reflectConfigFields`, `mergeGitSettings` |

Registry: `tied-project/spec/layer-b-wave-registry.yaml` → `waves.w2`

## Completed (W2 slice)

| Step | Evidence |
|------|----------|
| Tracker → `w2_corpus_in_progress` … verification | checklist `execution_evidence.completed` through `verification-gate` |
| Pseudocode gates (structural + Layer C + full) | `gate-pseudocode-validation-w2-*.md`, `pseudocode-validate-w2-full/` |
| `persist-implementation-records` | [`persist-implementation-records-w2.md`](persist-implementation-records-w2.md) |
| `unit-test-red` / `unit-test-green` | [`unit-test-red-w2.md`](unit-test-red-w2.md), [`unit-test-green-w2.md`](unit-test-green-w2.md) |
| Adversarial inquiry (W2 Mode B refresh) | `adversarial-inquiry/phase-pre_implementation/`, `phase-verification/` |
| `verification-gate` | [`w2-verification-gate.md`](w2-verification-gate.md), [`verification-evidence-manifest-w2.v1.json`](verification-evidence-manifest-w2.v1.json) |

**Verification gate:** `merged_decision.allowed: true` (advisory `finding_unresolved` / `warn_not_success` non-blocking).

**Sidecars (canonical):**

- `tied-project/implementation-decisions/IMPL-CONFIGURABLE_STRINGS-pseudocode.md`
- `tied-project/implementation-decisions/IMPL-CONFIG_DISPLAY_FLATTENING-pseudocode.md`

**Tests:**

- `config_pseudocode_w2_test.go` (new)
- `errors_test.go` — `TestHandleArchiveError` aligned to STEP / FAILURE_MODES

## Machine receipts (W2)

| Artifact | Path |
|----------|------|
| Layer B structural | `pseudocode-validate-w2/*.json` |
| Layer B full | `pseudocode-validate-w2-full/*.json` |
| Layer C gate_mode | `pseudocode-analysis/*.v1.json` (+ run copies `pseudocode-analysis-w2/*-layer-c.json`) |
| Verification manifest | `evidence/verification-evidence-manifest-w2.v1.json` |

**Layer C `input_identity.hash` (skip re-PSA if sidecar unchanged):**

- `IMPL-CONFIGURABLE_STRINGS`: `dd27b11eb465cf9afd6d802b436905ee0c488f990ea15e5647234999471c9b06`
- `IMPL-CONFIG_DISPLAY_FLATTENING`: `b76cfa476f05ee4978fe05ce5609b344b8639f59ba73a88e605a06a3dc24e2f2`

## Re-verify before commit / PR

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
go test -count=1 -run 'TestIMPL_CONFIG|TestHandleArchiveError' .
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode full --wave w2 --enforce \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w2-full
TIED_BASE_PATH="$PWD/tied-project" .cursor/skills/tied-yaml/scripts/tied-cli.sh tied_validate_consistency
```

## Next slice

- **W3** — 27 `lead_check_skip` IMPLs; new session handoff when W3 work starts (do not mix with W2 PR).

## Remaining corpus after W2

- **W3:** 27 `lead_check_skip` IMPLs per `layer-b-wave-registry.yaml`
- Layer C CI gate still **W0 six tokens** unless expanded (optional follow-up)

## Costly choice (sponsor)

Committing W2 verification + close-out in the same PR as W3 sidecar bulk — **costly**; prefer **one wave per PR**.
