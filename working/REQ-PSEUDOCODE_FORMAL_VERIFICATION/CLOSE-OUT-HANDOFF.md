# Lead hygiene — close-out handoff media

**Request:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`  
**CITDP:** `CITDP-LEAD-HYGIENE.yaml` → persisted `tied/citdp/CITDP-REQ-PSEUDOCODE_FORMAL_VERIFICATION.yaml` (evidence paths → `evidence/archive-lead-hygiene-20260919/` as of 2026-09-20)  
**Date:** 2026-09-20

## Completion signals

| Signal | Status | Evidence |
|--------|--------|----------|
| **Machine close-out** | **pass** | `close_out` gate `allowed: true`; `request-evidence-envelope.v1.json` validate with `fail_on_error_gaps: true` → **0 blocking gaps** (2026-09-20 via `run-close-out-gates.mjs --envelope-blocking`) |
| **Process contract** | **pass** | Tracker dispositions + slug evidence stubs; verification manifest; evidence-chain-profile; CITDP persist; corpus tracking 73/73 restored (bookkeeping, not hygiene proof) |
| **Adherence ledger** | **pass (grade C at hygiene time)** | Hygiene sync report archived: `evidence/archive-lead-hygiene-20260919/close-out-gates-run-20260920.json`. Current ledger: `gates/ledger.jsonl` (includes scan-root rows; reconcile grade **A** after `request_token`). |

## Proof boundary

- **In scope:** Package-level test lead removal, func-scoped leads, `check-leads`, audit `--fail-on-suspect`, hygiene scripts/tests, recurrence guard on `add-test-leads --all`.
- **Out of scope:** Runtime backup/config behavior; corpus queue flags prove maintenance state only after explicit restore.

## Evidence index

| Artifact | Path |
|----------|------|
| Bulk audit | `evidence/archive-lead-hygiene-20260919/bulk-lead-audit.json`, `.md` |
| Regression (close-out) | `evidence/archive-lead-hygiene-20260919/verification-regression-20260920-closeout.log` |
| Verification manifest (hygiene era) | superseded — see scan-root manifest at `evidence/verification-evidence-manifest.v1.json` (`scan-root-manifest-20260920`) |
| Evidence chain profile | `evidence/evidence-chain-profile.v1.json` |
| Request envelope | `evidence/request-evidence-envelope.v1.json` |
| Gates | `gates/close_out-2026-09-20T02-52-14-546Z.json` (latest close_out at sync time) |
| Adversarial inquiry | `adversarial-inquiry/phase-{pre_implementation,verification,close_out}/` |
| Tracker | `agent-req-implementation-checklist.yaml` |

## Replay (canonical)

```bash
cd /path/to/bkpdir
TIED_BASE_PATH="$PWD/tied" node /path/to/tied-repo/tools/bootstrap/templates/run-close-out-gates.mjs \
  --project-root "$PWD" \
  --request-token REQ-PSEUDOCODE_FORMAL_VERIFICATION \
  --tracker-path working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/agent-req-implementation-checklist.yaml \
  --citdp-path working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/CITDP-LEAD-HYGIENE.yaml \
  --phase close_out \
  --run-id lead-hygiene-20260919-close \
  --envelope-blocking \
  --sync-dispositions \
  --reconcile

SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
scripts/spec-corpus-status.sh
python3 scripts/audit_package_level_leads.py --fail-on-suspect
```

## Proposed commit message

```
fix(traceability): remove package-level IMPL test leads and guard recurrence

Strip bulk add-test-leads paste from *_test.go, restore func-scoped leads,
add audit/strip tooling, disable --all, and restore spec corpus tracking
bookkeeping. [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
```

**Do not commit** unless the sponsor requests `traceable-commit`.

## Gitignore close-out

**N/A** — `working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/` evidence is intentional per-request artifacts; no new ephemeral patterns required beyond existing `.gitignore`.
