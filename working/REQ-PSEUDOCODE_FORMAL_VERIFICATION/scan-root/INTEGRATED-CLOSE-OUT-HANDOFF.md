# Scan-root CITDP — integrated close-out handoff

**Slice:** `check-leads` scan roots (`LeadScanDirPaths`)  
**Parent REQ:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`  
**CITDP:** `CITDP-SCAN-ROOT.yaml` → `tied/citdp/CITDP-REQ-PSEUDOCODE_FORMAL_VERIFICATION-SCAN-ROOT.yaml`  
**Tracker:** `scan-root/agent-req-implementation-checklist.yaml`  
**Date:** 2026-09-20

## Completion signals

| Signal | Status | Evidence |
|--------|--------|----------|
| **Machine close-out** | **pass** | `close_out` gate `allowed: true`; `request-evidence-envelope.v1.json` validate with `fail_on_error_gaps: true` → **0 blocking gaps** (6 advisory `finding_unresolved`) |
| **Process contract** | **pass** | Scan-root Tracker dispositions synced; verification manifest collected; CITDP integrated depth recorded |
| **Adherence ledger** | **pass (grade A)** | Top-level `request_token: REQ-PSEUDOCODE_FORMAL_VERIFICATION` on parent + scan-root trackers (rubric uses `request_token`, not `execution_evidence.request`). Reconcile 2026-09-20: score **100**, band **A**, all dimensions 100. |

## Inquiry identity (scan-root run)

| Phase | `run_id` |
|-------|----------|
| `pre_implementation` | `scan-root-20260920-pre` |
| `verification` | `scan-root-20260920-ver` |
| `close_out` | `scan-root-20260920-close` |

Artifacts: `working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/adversarial-inquiry/phase-{phase}/`  
Scope: `IMPL-SPEC_CTL` + `internal/speccheck/leads_test.go` / `leads.go`; advisory policy (UNRESOLVED findings are **warn-only**, not LEAP triggers).

**Note:** This inquiry pass **supersedes** hygiene-era provenance in those phase directories for gate/envelope purposes. Hygiene history remains in `CLOSE-OUT-HANDOFF.md` and `gates/archive-hygiene-20260919/`.

**Working layout:** Use `scan-root/` under this REQ only. The mistaken bootstrap copy `working/REQ-PSEUDOCODE_FORMAL_VERIFICATION-SCAN-ROOT/` was removed 2026-09-20 (invalid `REQ-*-SCAN-ROOT` token; gates used parent `REQ-PSEUDOCODE_FORMAL_VERIFICATION`).

## Replay

```bash
cd /Users/fareed/Documents/dev/go/bkpdir
TIED_BASE_PATH="$PWD/tied" node /Users/fareed/Documents/dev/chatgpt/stdd/tools/bootstrap/templates/run-close-out-gates.mjs \
  --project-root "$PWD" \
  --request-token REQ-PSEUDOCODE_FORMAL_VERIFICATION \
  --tracker-path working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/scan-root/agent-req-implementation-checklist.yaml \
  --citdp-path working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/scan-root/CITDP-SCAN-ROOT.yaml \
  --phase close_out \
  --run-id scan-root-20260920-close \
  --envelope-blocking \
  --sync-dispositions \
  --reconcile

go test ./internal/speccheck/... -count=1
go run ./cmd/specctl check-leads tied/implementation-decisions/*-pseudocode.md
```

## Key artifacts

| Artifact | Path |
|----------|------|
| Close-out run JSON | `scan-root/evidence/close-out-gates-run-integrated-20260920.json` |
| Post–`request_token` sync | `scan-root/evidence/close-out-gates-sync-20260920.json` (grade **A**, envelope rebuilt) |
| Envelope (copy) | `scan-root/evidence/request-evidence-envelope.v1.json` |
| Canonical envelope | `../evidence/request-evidence-envelope.v1.json` (slim index; hygiene discovery → `../evidence/archive-lead-hygiene-20260919/`) |
| Regression log | `scan-root/evidence/check-leads-regression.log` |
| Verification manifest | `evidence/verification-evidence-manifest.v1.json` (`run_id: scan-root-manifest-20260920`; command stdout under `scan-root/evidence/manifest-run/`) |
| Adherence ledger | `../gates/ledger.jsonl` |

## Proof boundary

Scan-root close-out proves **check-leads scan directory expansion** and regression compatibility for this slice. It does **not** re-certify the 2026-09-19 lead-hygiene strip/corpus maintenance (see parent `CLOSE-OUT-HANDOFF.md`).
