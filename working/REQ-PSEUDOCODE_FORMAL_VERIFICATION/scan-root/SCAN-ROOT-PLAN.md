# check-leads scan-root CITDP — executed

**Request:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`  
**CITDP:** `CITDP-SCAN-ROOT.yaml` → `tied/citdp/CITDP-REQ-PSEUDOCODE_FORMAL_VERIFICATION-SCAN-ROOT.yaml`  
**Depth:** `minimal` (tooling-only; no integrated close-out envelope for this slice)

## Change

- `speccheck.LeadScanDirPaths(repoRoot)` — `pkg/`, `internal/`, `cmd/`, `test/` when directories exist.
- `specctl check-leads` uses that helper instead of hard-coded `pkg/` only.
- IMPL `CHECK_LEADS_COMMAND` STEP T002 aligned with implementation.
- Unit tests in `internal/speccheck/leads_test.go`.

## Verification

```bash
go test ./internal/speccheck/... -count=1
go run ./cmd/specctl check-leads tied/implementation-decisions/*-pseudocode.md
go test ./... -count=1
```

Log: `evidence/check-leads-regression.log`

## Integrated close-out (2026-09-20)

Full handoff: [`INTEGRATED-CLOSE-OUT-HANDOFF.md`](INTEGRATED-CLOSE-OUT-HANDOFF.md)

| Signal | Status |
|--------|--------|
| **Machine close-out** | **pass** — `merged_decision.allowed: true`; envelope **0** blocking gaps |
| **Process contract** | Scan-root Tracker + CITDP (`integrated`) + manifest |
| **Adherence ledger** | Reconcile grade **B** (75); advisory inquiry warns only |

`sub-adversarial-inquiry-pass`: **completed** (three-phase inquiry, `scan-root-20260920-*` run ids).  
`traceable-commit`: **waived** until sponsor requests commit.
