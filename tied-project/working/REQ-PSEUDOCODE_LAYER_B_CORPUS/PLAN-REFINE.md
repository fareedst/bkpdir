# Plan refine handoff — REQ-PSEUDOCODE_LAYER_B_CORPUS

**Parent REQ:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`  
**Slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS`  
**Refined:** 2026-10-21 via `@refine-plan`

## RESOLVE (Touchpoint 1)

| Sponsor term | Canonical meaning |
|--------------|-------------------|
| Prior standard | L0 **formal DSL** (`SPEC-ID` / `STEP`, block leads, specctl green) — not Layer B `minimum_gating_rules` |
| New standard | Layer B at **verification-gate** + integrated **adversarial inquiry** (checklist `sub-adversarial-inquiry-pass`) |
| Hybrid workflow | **Track A:** scripts/CI/tests (may start before full Tracker green). **Track B:** sidecar/YAML edits only under full checklist after gates |
| Exemplar sidecar | `[IMPL-EXCLUSION_PATTERNS]` — `Grammar-Version: v2` + nested `Contract:` precision |

**Intent authority:** Sponsor goal is measurable corpus uplift, not cosmetic v2 headers.

## Depth and gate policy

| Field | Value |
|-------|--------|
| `profile_depth` / `depth_tier` | **integrated** |
| `gate_policy` | **advisory** (blocking escalation only if CITDP later sets strict) |
| Inquiry run_ids (planned) | `layer-b-corpus-20261021-{pre,ver,close}` |

## Pre-implementation gate (2026-10-21)

`tied_checklist_gate_validate` → **`allowed: false`** (expected for sparse bootstrap Tracker):

- `tracker_sparse` / `malformed_tracker:steps`
- Missing: `risk-assessment`, `sub-adversarial-inquiry-pass`, `gate-pseudocode-validation`
- `integrated_depth_requires_pairing` / `activation_pairing_incomplete`

**Before RED or sidecar edits:** Copy full step list from [`tied-bundle/docs/agent-req-implementation-checklist.yaml`](../../../tied-bundle/docs/agent-req-implementation-checklist.yaml) into this folder’s Tracker (or symlink workflow from `session-bootstrap` through `test-strategy`), complete `risk-assessment` dispositions, run structural `sub-adversarial-inquiry-pass`, then re-run gate with activation when phase dir exists.

## Invalidated downstream evidence

- `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]` lead-hygiene / scan-root **close_out** dispositions — **do not reuse** for this slice.
- Prior plan assumption of pure `non-tied-plan` for sidecars — superseded by **refine-plan** full TIED tracking for Track B.

## Artifacts

| Artifact | Path |
|----------|------|
| CITDP | [`CITDP-REQ-PSEUDOCODE_LAYER_B_CORPUS.yaml`](CITDP-REQ-PSEUDOCODE_LAYER_B_CORPUS.yaml) |
| Tracker (bootstrap) | [`agent-req-implementation-checklist.yaml`](agent-req-implementation-checklist.yaml) |
| Linked Cursor plan | `~/.cursor/plans/impl_pseudocode_modernization_59816848.plan.md` |
