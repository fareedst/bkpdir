# Continuation handoff — W1–W3 corpus (post–W0 close-out)

**Parent slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS`  
**Prior close-out:** [`close-out-handoff.md`](close-out-handoff.md) (`layer-b-corpus-20261021-close-v2`, W0 + Track A only)  
**Plan:** `~/.cursor/plans/impl_pseudocode_modernization_59816848.plan.md` Phase 3–5  
**Run ID (suggested):** `layer-b-corpus-w1-ver` (W1 complete — see [`w1-verification-gate.md`](w1-verification-gate.md))

## RESOLVE

| Phrase | Meaning |
|--------|---------|
| **73/73 L0** | All **73** `formal_spec` sidecars in [`tied-project/spec/formal-spec-registry.yaml`](../../../spec/formal-spec-registry.yaml) stay on the L0 merge gate (`specctl validate` / matrix / `check-leads`). |
| **TRACE-002 debt** | Literal block-lead strings from each sidecar must appear in Go comments under `pkg/`, `internal/`, `cmd/`, `test/`, and repo-root `*.go`. Not the same as Layer B SHAPE/COVER/CONTRACT. |
| **W1–W3** | Sidecar **Track B** waves after W0 (Layer B contract depth + tests), not a fourth close-out on W0 evidence. |

## L0 baseline (2026-10-21 continuation)

After W0 close-out, **`specctl check-leads` had four TRACE-002 errors** on diff/inc/directory blocks (wording drift vs sidecar). Those leads were realigned in `comparison.go`, `pkg/fileops/comparison.go`, and `archive.go`.

**Verify before W1 sidecar edits:**

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

Expect **zero** `[error] TRACE-002` lines. Keep this green on every wave PR.

## Wave map (authoritative intent)

| Wave | IMPL scope (initial) | Verification-gate bar |
|------|----------------------|------------------------|
| **W0** | 6 tokens (5 Tier A Layer C + `IMPL-EXCLUSION_PATTERNS`) | Done — see W0 gate receipts |
| **W1** | **38** tokens — post-W0 oracle runtime (`layer-b-wave-registry.yaml` `waves.w1`) | Layer B **full** + FAILURE_MODES / DATA_TRANSITION / TERMINATION where triggered |
| **W2** | **2** tokens — formal without oracle (not in `lead_check_skip`) | COVER-001/002: tests assert STEP behavior |
| **W3** | **27** tokens — `lead_check_skip` Tier C doc/process | Structural SHAPE + meaningful How; waivers reviewed at close-out |

**W0 enforce set (batch driver):**  
`scripts/run-pseudocode-layer-b-batch.sh --mode full --wave w0 --enforce`

**W1+:** Define `--wave w1|w2|w3` token lists in [`tied-project/spec/layer-b-wave-registry.yaml`](../../../spec/layer-b-wave-registry.yaml) (TIED pass) before expanding CI `RUN_LAYER_B` enforcement.

## Track B entry (per wave)

1. Copy or reset Tracker from [`agent-req-implementation-checklist.yaml`](../agent-req-implementation-checklist.yaml) — clear W0 `close_out_evidence`; set `operator_evidence.status` to `w{N}_corpus_in_progress`.
2. CITDP working copy → update `depth_tier: integrated`, wave-scoped `success_criteria`, then `persist-citdp-record` when behavior changes.
3. `gate-pseudocode-validation` (structural) → edit sidecars (token comments on every block; v2 **on touch only**; exemplar [`IMPL-EXCLUSION_PATTERNS-pseudocode.md`](../../../implementation-decisions/IMPL-EXCLUSION_PATTERNS-pseudocode.md)).
4. `sub-adversarial-inquiry-pass` at **pre_implementation** and **verification**; four artifacts under `adversarial-inquiry/phase-{phase}/`.
5. Sync Go/test block leads; add **failure-path tests** when FAILURE_MODES expand (LEAP if behavior shifts).
6. `verification-gate`: `run-pseudocode-layer-b-batch.sh --mode full --wave w{N} --enforce` + `tied_validate_consistency`.
7. Close-out: new envelope + `close_out` gate — do **not** reuse W0 receipts.

## Evidence anchors

| Artifact | Path |
|----------|------|
| Gap report (Phase 0) | [`corpus-layer-b-gap-report.json`](corpus-layer-b-gap-report.json) |
| Layer B reports dir | `../pseudocode-validate/` |
| Persisted CITDP | [`../../../citdp/CITDP-REQ-PSEUDOCODE_LAYER_B_CORPUS.yaml`](../../../citdp/CITDP-REQ-PSEUDOCODE_LAYER_B_CORPUS.yaml) |

## Out of scope for this handoff

- PFV parent [`REQ-PSEUDOCODE_FORMAL_VERIFICATION`](../REQ-PSEUDOCODE_FORMAL_VERIFICATION/) scan-root close-out dispositions
- Bulk `Grammar-Version: v2` without contract edits
- Replacing MCP Layer B with duplicated SHAPE logic in `specctl`

## Costly choices (sponsor)

| Choice | Ladder |
|--------|--------|
| Enforce Layer B on all 67 post-W0 IMPLs in one PR | **Costly** — prefer one wave per PR |
| Merge gate switches from L0-only to L0+Layer B fleet-wide before W1 completes | **Costly** — use opt-in `RUN_LAYER_B` per wave |
