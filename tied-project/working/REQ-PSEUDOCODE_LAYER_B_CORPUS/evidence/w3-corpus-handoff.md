# Continuation handoff — W3 corpus (close-out complete)

**Parent slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS`  
**Prior close-out:** W2 — [`w2-close-out-handoff.md`](w2-close-out-handoff.md) (commit `f8608c1…`)  
**Wave map:** [`w1-w3-corpus-handoff.md`](w1-w3-corpus-handoff.md)  
**Tracker:** `operator_evidence.status: w3_slice_close_out` — [`../agent-req-implementation-checklist.yaml`](../agent-req-implementation-checklist.yaml)  

**Close-out:** [`w3-close-out-handoff.md`](w3-close-out-handoff.md) (feature commit `9adc16b…`)

## W3 scope (27 IMPLs)

Tier C `lead_check_skip` doc/process tokens — registry `tied-project/spec/layer-b-wave-registry.yaml` → `waves.w3`

## Completed (W3 slice)

| Step | Evidence |
|------|----------|
| Structural + full Layer B | [`gate-pseudocode-validation-w3-structural.md`](gate-pseudocode-validation-w3-structural.md), `pseudocode-validate-w3-full/*.json` |
| persist-implementation-records | [`persist-implementation-records-w3.md`](persist-implementation-records-w3.md) |
| verification-gate | [`w3-verification-gate.md`](w3-verification-gate.md), [`verification-evidence-manifest-w3.v1.json`](verification-evidence-manifest-w3.v1.json) |

**Verification gate:** `allowed: true` (advisory inquiry diagnostics non-blocking).

## Re-verify before commit / PR

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode full --wave w3 --enforce \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-full
TIED_BASE_PATH="$PWD/tied-project" .cursor/skills/tied-yaml/scripts/tied-cli.sh tied_validate_consistency
```

## Next (post–W3)

- **Corpus Layer B waves done** — W0–W3 = 73/73 registry tokens; optional CI `RUN_LAYER_B` fleet enforce, Grammar-Version v2 on touch, D15 banner hygiene

## Remaining risks

- L0 D15 redundant banners (`pkg/fileops/comparison.go`) — dedicated hygiene pass
- Envelope hydration advisory (`hydration_envelope_gaps_missing`) until close-out sync

## Costly choice (sponsor)

Merging W3 sidecar bulk with W2 close-out in one PR — **costly**; prefer **W3-only PR** after close-out.
