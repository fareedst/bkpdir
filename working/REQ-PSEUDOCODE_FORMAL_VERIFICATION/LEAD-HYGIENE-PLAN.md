# Lead hygiene plan — current state

**Request:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`  
**Working folder:** `working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/`  
**CITDP:** `CITDP-LEAD-HYGIENE.yaml`  
**Tracker:** `agent-req-implementation-checklist.yaml` (`execution_evidence.request` set)  
**Last updated:** 2026-09-20 (refine-plan sync with Cursor plan; close-out complete)  
**Cursor plan mirror:** [`repair_bulk_test_leads_3e9a0d0c.plan.md`](file:///Users/fareed/.cursor/plans/repair_bulk_test_leads_3e9a0d0c.plan.md)  
**Handoff media:** [`CLOSE-OUT-HANDOFF.md`](CLOSE-OUT-HANDOFF.md)

**Status:** Hygiene + integrated close-out **complete**. No further implementation waves unless a new CITDP is opened (see § Post-close-out).

---

## Goal

Remove misleading **package-level** `// - [IMPL-*] — How:` mega blocks pasted by bulk `add-test-leads --all`, restore **Track C** placement (leads immediately before implementing `func` or covering `Test*` / `Benchmark*`), and keep **Layer B** green (`specctl validate`, `check-leads`, matrix, hygiene audit).

**Out of scope (separate CITDP):** Expanding `specctl check-leads` scan roots to `internal/`, `cmd/`, `test/`; runtime behavior changes.

## Refined close-out contract

This pass closes evidence for work already applied; it does not authorize another
implementation wave. The strict close-out path selects `depth_tier: integrated`
with `gate_policy: advisory`, because the checklist treats strict close-out as an
integrated-inquiry trigger. The existing behavior-neutral rationale remains valid
for runtime risk, but it does not waive the close-out evidence contract.

The authoritative Tracker is
`working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/agent-req-implementation-checklist.yaml`.
Step dispositions must be written on the step rows with typed `evidence_refs`;
`execution_evidence.completed` is compatibility metadata and cannot substitute
for those dispositions. No step is marked complete from this plan alone.

### Inquiry identity (per phase)

| Phase | `run_id` (must match `evidence-provenance.json`) |
|-------|---------------------------------------------------|
| `pre_implementation` | `lead-hygiene-20260919-pre` |
| `verification` | `lead-hygiene-20260919-ver` |
| `close_out` | `lead-hygiene-20260919-close` |

`CITDP-LEAD-HYGIENE.yaml` `completion_criteria.activation` records the **close_out**
binding only. For replay, call `tied_checklist_activation_collect` with the **phase**
`run_id`, not the close-out id.

### Close-out sequence

1. Run the `pre_implementation` inquiry and persist its four bounded artifacts in
   the request-scoped phase directory.
2. Call `tied_checklist_activation_collect` with the phase `run_id`, then
   `tied_checklist_gate_validate` for `pre_implementation` (pass the full collect
   payload; partial activation fails closed).
3. Re-run the existing Layer B, regression, and lead-hygiene checks; attach
   command provenance and results to the verification evidence.
4. Run the verification inquiry and gate. Any discovered divergence follows
   IMPL → test → code LEAP order; this plan does not expand runtime scope.
5. Build and validate `request-evidence-envelope.v1.json` with
   `fail_on_error_gaps: true`, run the close-out inquiry, and call the
   `close_out` gate.
6. Report machine close-out, process contract, and adherence ledger separately.
   Corpus bookkeeping is an independent maintenance track and cannot be used as
   proof of lead hygiene.

---

## Root cause

[`scripts/impl_pseudocode_remediation.py`](../../scripts/impl_pseudocode_remediation.py) `add-test-leads --all` prepended **entire sidecar lead lists** after `package` in all `*_test.go` under `TOKEN_TEST_DIRS` entries mapping to `"."`, producing ~80–175 lines per file across **52** suspects.

---

## Progress summary

| Phase | Status | Notes |
|-------|--------|--------|
| 0 Bootstrap (Tracker, CITDP, audit) | **Done** | Evidence: `evidence/archive-lead-hygiene-20260919/bulk-lead-audit.json`, `.md` |
| 1 Strip (waves A/B/C) | **Done** | 51 files, **7,374** package-level lines removed |
| 2 Production backfill (`TRACE-002`) | **Done** | 0 `check-leads` errors on 73 sidecars; fixes in `comparison.go`, `pkg/testutil/doc.go` |
| 3 Three-way / test leads | **Done (classifier)** | `test_leads_missing` / `test_gap` → **0** tokens; func-scoped leads + sync batches |
| 4 Prevent recurrence | **Done** | `--all` disabled; audit in `run-spec-verification.sh`; impl-deep-sync guide note |
| 5 Close-out | **Done (hygiene scope)** | Integrated inquiry (3 phases); gates `allowed: true`; evidence-chain-profile + envelope rebuilt; Layer B regression re-run 2026-09-20; corpus tracking **73/73** restored 2026-09-20 |

---

## Tooling (delivered)

| Script | Purpose |
|--------|---------|
| [`scripts/lead_hygiene.py`](../../scripts/lead_hygiene.py) | Shared audit/strip logic |
| [`scripts/audit_package_level_leads.py`](../../scripts/audit_package_level_leads.py) | Inventory; `--fail-on-suspect` (CI) |
| [`scripts/strip_package_level_leads.py`](../../scripts/strip_package_level_leads.py) | `--wave a\|b\|c\|all`, dry-run default, `--apply` |
| [`scripts/strip_redundant_impl_banners.py`](../../scripts/strip_redundant_impl_banners.py) | D15: drop paraphrased `// [IMPL-*]` when `// -` block leads exist; `--fail-on-redundant` |
| [`scripts/test_lead_hygiene.py`](../../scripts/test_lead_hygiene.py) | Unit tests + `scripts/testdata/lead_hygiene/` |

**Policy:** All `*_test.go` files — **no** package-level `// - [IMPL-*]` blocks; test leads only before `func Test*` / `Benchmark*`.

**`add-test-leads`:** Func-scoped only; `--all` exits with error. Extended `TOKEN_TEST_DIRS` for `IMPL-SPEC_CTL`, `IMPL-DOC_ENHANCEMENT`.

---

## Verification (current)

Run from repo root:

```bash
python3 scripts/audit_package_level_leads.py --fail-on-suspect
go run ./cmd/specctl check-leads tied/implementation-decisions/*-pseudocode.md
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
go test ./... -count=1
scripts/spec-corpus-status.sh
```

| Check | Result (2026-09-19) |
|-------|-------------------|
| Package-level lead audit | **PASS** — 0 suspect files |
| `specctl validate` (73 sidecars) | **PASS** |
| `specctl check-leads` (73 registry) | **PASS** |
| Matrix `--check-coverage` | **PASS** |
| `go test ./...` | **PASS** |
| `classify_token` test gaps | **PASS** — 0 tokens with `test_leads_missing` / `test_gap` |
| `spec-corpus-status` ALL GREEN | **PASS** (2026-09-20) — tracking restored via logic audit + queue field sync (`formal_spec`, `deep_sync_complete`) |

**Registry grep fix:** [`scripts/run-spec-verification.sh`](../../scripts/run-spec-verification.sh) and [`scripts/spec-corpus-status.sh`](../../scripts/spec-corpus-status.sh) now parse `formal_spec` entries as `- "IMPL-..."`.

---

## Test strategy and proof boundary

- Run the package-level audit with `--fail-on-suspect`; zero suspects is the
  acceptance criterion for the repaired hygiene defect.
- Run `specctl validate`, `specctl check-leads`, matrix coverage, and the fixture
  tests for `scripts/lead_hygiene.py` and
  `scripts/strip_package_level_leads.py`.
- Run `SKIP_TIED_MCP=1 scripts/run-spec-verification.sh` and
  `go test ./... -count=1` as regression evidence. These checks validate
  placement, formal-spec coverage, and runtime compatibility; they do not prove
  corpus bookkeeping.
- No composition or E2E tests apply: this pass changes comments, analysis
  tooling, and evidence records only; it introduces no UI, IPC, network,
  persistence, or runtime binding behavior.
- The 73 corpus rows are a separate maintenance denominator. Restore them only
  through their dedicated scripts and report before/after counts independently.

---

## Notable file changes

- **Strip:** All root/pkg/internal/test/cmd/tools `*_test.go` mega blocks removed; exemplar [`config_bench_test.go`](../../config_bench_test.go) — `[REQ-CFG_006] [REQ-PERFORMANCE]`, benchmarks only.
- **Prod backfill:** [`comparison.go`](../../comparison.go) (PACKAGE_EXTRACTION); [`pkg/testutil/doc.go`](../../pkg/testutil/doc.go) (TESTING stress-path lead).
- **Tests:** [`pkg/config/config_test.go`](../../pkg/config/config_test.go) — stripped then func-scoped CFG merge leads; [`test/specconformance/specctl_conformance_test.go`](../../test/specconformance/specctl_conformance_test.go) — SPEC_CTL sidecar leads.
- **Docs:** [`tied/docs/impl-deep-sync-agent-guide.md`](../../tied/docs/impl-deep-sync-agent-guide.md) — bulk `add-test-leads --all` trap row.

---

## Incidents / lessons

1. **`git checkout pkg/config/config_test.go`** after strip reintroduced mega block — re-ran `strip_package_level_leads.py --apply` on that file.
2. **Broken `add_test_leads`** (insert mid-function) — fixed to insert **once per lead** before a **single** matching `Test`/`Benchmark`; always run `go test ./pkg/config/...` after batch inserts.

---

## Remaining work

### A. Checklist / gates (REQ workflow)

- [x] Authoritative Tracker step dispositions with typed evidence refs (profile,
  CITDP persist, close-out sync sub); do not rely on `execution_evidence.completed`
  alone.
- [x] Integrated `pre_implementation` inquiry artifacts, activation collection,
  and `tied_checklist_gate_validate` (2026-09-20 receipts under `gates/`).
- [x] Verification inquiry and gate with fresh Layer B/regression evidence
  (`verification-regression-20260920.log`).
- [x] Close-out envelope with zero **blocking** gaps (`run-close-out-gates.mjs
  --envelope-blocking`, 2026-09-20); advisory inquiry warns remain.
- [x] `outcome_verified` ledger rows via `sync-tracker-dispositions.mjs`; reconcile
  report in `evidence/close-out-gates-run-20260920.json` (process grade band C:
  stale historical gate receipts, ledger correlation 100%).

### B. Corpus tracking rows (optional for ALL GREEN banner)

- [x] Restored **2026-09-20**: `run_impl_logic_audit.py --batch all --mark-complete
  --mark-crit001 --req-criteria-scoped` after excluding `scripts/testdata` from
  `go_block_leads` (hygiene fixture had caused `extra_in_go` on
  `IMPL-CFG_QUOTED_KEY_PREFIX`). **After:** 73/73 `req_audit_pass`,
  `logic_verified`, `deep_sync_complete`, `crit001_pass`; `spec-corpus-status.sh`
  ALL GREEN for tracking lines.

### C. Hygiene polish (optional)

- [x] **D15:** `strip_redundant_impl_banners.py --apply` removed **424** paraphrased banner lines in **37** files (2026-09-20); audit `evidence/d15-banner-audit.json` post-run **0** redundant.
- [x] **LEAP:** `IMPL-TOKEN_COVERAGE_AUDIT-pseudocode.md` **HYGIENE_POLICY** block (module-scoped `// -` leads; no package paste; audit/strip scripts).
- [x] **Follow-up CITDP:** Scan roots widened **2026-09-20** — see `scan-root/SCAN-ROOT-PLAN.md` and `tied/citdp/CITDP-REQ-PSEUDOCODE_FORMAL_VERIFICATION-SCAN-ROOT.yaml` (`LeadScanDirPaths`: pkg/, internal/, cmd/, test/).

---

## Post-close-out (separate work items)

| Item | Owner | Proof if executed |
|------|-------|-------------------|
| `check-leads` scan-root expansion | New CITDP | Conformance tests + matrix; no claim from hygiene envelope |
| Broad `IMPL-TOKEN_COVERAGE_AUDIT` rewrite | Optional LEAP | Pseudo-code + three-way sync |
| Sponsor `traceable-commit` | Sponsor | Proposed message in [`CLOSE-OUT-HANDOFF.md`](CLOSE-OUT-HANDOFF.md) |

---

## Original phased plan (reference)

### Phase 0 — Bootstrap

1. Copy Tracker; write `CITDP-LEAD-HYGIENE.yaml`.  
2. Pre-implementation gate.  
3. `audit_package_level_leads.py` → `evidence/archive-lead-hygiene-20260919/bulk-lead-audit.json`.

### Phase 1 — Strip

`strip_package_level_leads.py --apply --wave all` (or a → b → c). Re-run audit after each wave.

### Phase 2 — Production backfill

After strip: `specctl check-leads`; `sync_impl_block_leads.py --token IMPL-*` on verified YAML maps.

### Phase 3 — Three-way by token

D11–D15 per [`formal-spec-registry.yaml`](../../tied/spec/formal-spec-registry.yaml); prioritize P0/P1; pilots (`CLI_FRAMEWORK`, `CONFIG_STRUCT`, `LIST_FORMAT_SAFETY`, `STRUCTURED_ERRORS`) — do not regress.

**Completed in repo:** Classifier shows no test lead gaps; production `check-leads` clean; func-scoped insertion for remaining P1/P3 tokens (batch 2026-09-19).

### Phase 4 — Prevent recurrence

Disable package `--all`; CI audit; document in impl-deep-sync guide. **Done.**

### Phase 5 — Close-out

`run-spec-verification.sh`, full tests, corpus status, CITDP persist, verification/close-out gates, envelope. **Done (2026-09-20)** — see progress table and [`CLOSE-OUT-HANDOFF.md`](CLOSE-OUT-HANDOFF.md).

---

## Key references

- Process: [`tied/docs/impl-deep-sync-agent-guide.md`](../../tied/docs/impl-deep-sync-agent-guide.md)  
- Layer B: [`docs/markscope/spec-verification.md`](../../docs/markscope/spec-verification.md)  
- Pollution source / fix: [`scripts/impl_pseudocode_remediation.py`](../../scripts/impl_pseudocode_remediation.py)  
- Lead gate: [`internal/speccheck/leads.go`](../../internal/speccheck/leads.go), [`cmd/specctl/main.go`](../../cmd/specctl/main.go)
