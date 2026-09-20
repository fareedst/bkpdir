# IMPL Deep Sync — Agent Operating Guide

**Audience:** AI coding agents completing the brownfield **Track C** pass: promote every project `IMPL-*` from bulk sidecar migration to **pilot-quality** three-way alignment (sidecar ↔ production code ↔ tests).

**Process tokens:** `[PROC-IMPL_CODE_TEST_SYNC]`, `[PROC-IMPL_PSEUDOCODE_TOKENS]`, `[PROC-PSEUDOCODE_VALIDATION]`, `[PROC-LEAP]`

**Master queue:** [`impl-pseudocode-sync-checklist.yaml`](impl-pseudocode-sync-checklist.yaml) — work strictly by `process_order`, **one token per work item**.

**Companion docs:** [pseudocode-writing-and-validation.md](pseudocode-writing-and-validation.md) (Track C theory), [AGENTS.md](../../AGENTS.md) (session rules), [`.cursor/skills/tied-yaml/SKILL.md`](../../.cursor/skills/tied-yaml/SKILL.md) (TIED YAML mutations).

---

## 1. What “done” means (pilot bar)

A token is **pilot-complete** only when **all** of the following hold. Checklist `status: validated` must mean this — not “bulk migration ran once.”

| Layer | Requirement |
| --- | --- |
| **Sidecar** | Language-agnostic `IMPL-*-pseudocode.md` with `## Summary contract`, one `## PROCEDURE` per logical block, **no** `"Block implements documented behavior for:"` or `per resource-management atomic I/O contract` boilerplate |
| **Block lead** | Exactly one line per block: `- [IMPL-…] [ARCH-…] [REQ-…] — How: <specific behavior>.` naming tokens **and** how the block implements them ([PROC-IMPL_PSEUDOCODE_TOKENS]) |
| **Code** | Immediately before each implementing `func`, a **byte-identical** copy: `// - {same lead as sidecar}` (not paraphrased `// [IMPL-…]` file banners) |
| **Tests** | Same literal `// - {lead}` before `Test*` (or test file header) when the test exercises that block |
| **Counts** | `blocks_total` = `blocks_synced` = number of `^- [IMPL-` leads in the sidecar (recounted after rewrite) |
| **Layer B** | Manually satisfied per [pseudocode-validation-checklist.yaml](pseudocode-validation-checklist.yaml) — do not set `layer_b_validated: true` from scripts alone |
| **Build** | Targeted packages + `go test ./...` green after Go changes |
| **TIED** | `tied_validate_consistency` (MCP) reports no issues for changed indexes |

**Gold reference (read first):** [`../implementation-decisions/IMPL-LIST_FORMAT_SAFETY-pseudocode.md`](../implementation-decisions/IMPL-LIST_FORMAT_SAFETY-pseudocode.md) and matching leads in `formatter.go`.

**Original pilots (skip re-work unless regression):** `IMPL-CLI_FRAMEWORK` (order 8), `IMPL-CONFIG_STRUCT` (12), `IMPL-LIST_FORMAT_SAFETY` (28), `IMPL-STRUCTURED_ERRORS` (37). Rows are in `PILOT_SYNCED` / early `MANUAL_DEEP_SYNCED` — do not bulk-overwrite.

---

## 2. What is *not* sufficient (bulk migration trap)

These steps were run project-wide earlier; they are **starting points only**:

| Artifact | Why insufficient |
| --- | --- |
| Sidecar file exists | Often weak structure, Go-isms, truncated bodies, or boilerplate `How:` lines |
| `sync_impl_block_leads.py` | Copies leads only after sidecar is correct; cannot fix wrong YAML↔function mapping; risky on broken sidecars |
| `rewrite_bulk_sidecar_to_pilot.py` | In-place boilerplate→`How:` only; **skipped** tokens with ≥2 “meaningful” leads; **damaged** structure-poor sidecars in at least one case |
| `run_impl_deep_sync_batch.py` + checklist `validated` | Set `layer_b_validated` / `blocks_synced` without literal three-way verification |
| `impl_pseudocode_remediation.py add-test-leads --all` | **Disabled.** Pasted entire sidecar lead lists after `package` in many `*_test.go` files; use func-scoped `--token` or `sync_impl_block_leads.py` on production only. Remediation: `scripts/strip_package_level_leads.py` + `scripts/audit_package_level_leads.py`. |
| Paraphrased `// [IMPL-*]` file banners (D15) | Obsolete when the same IMPL token has a literal `// - {sidecar lead}` in the same file. Remediation: `scripts/strip_redundant_impl_banners.py` (dry-run default; `--apply`). |
| Checklist row `validated` without `MANUAL_DEEP_SYNCED` | Bulk scripts may still rewrite notes on the next inventory pass |

**Rule:** Treat any row with `blocks_synced: 0`, boilerplate in the sidecar, or mismatched `// -` leads as **not done**, even if `status: validated`.

---

## 3. Session bootstrap (every agent run)

1. Preface responses per [AGENTS.md](../../AGENTS.md); confirm `tied/docs/ai-principles.md` read this session.
2. **MCP:** `tied_config_get_base_path` → must resolve to **this repo’s** `tied/` directory.
3. Open [`impl-pseudocode-sync-checklist.yaml`](impl-pseudocode-sync-checklist.yaml); pick the **lowest `process_order`** row that is not pilot-complete.
4. Set that row `status: in_progress` (or `manual_deep_sync`) while working; revert others to `validated` only when exit criteria pass.
5. Read the gold reference sidecar + this token’s `IMPL-{TOKEN}.yaml` before editing.

**Stop rule:** Finish **one token completely** (phases A–H, or Tier C variant) before starting the next.

---

## 4. Work queue: waves and `process_order`

Work **top to bottom** by `process_order` in the checklist. Waves group calendar batches; **order is canonical**, not tier clusters.

| Wave | `process_order` | Skip (already pilot) | Focus |
| --- | --- | --- | --- |
| **1** | 1–7 | 8 | CFG / atomic / precedence cluster |
| **2** | 9–27 | 12, 28 | Tier A runtime (large: DATA_MODELS, DELAYED_OUTPUT, EXCLUSION_PATTERNS, FILE_OPERATIONS, GIT_CLI) |
| **3** | 29–42 | 37 | Tier A tail, meta, vis, zip |
| **4** | 43–60 | — | Tier B (CFG fixes, test IMPL specs, doc migration) |
| **5** | 61–72 | — | Tier C doc/process (sidecar quality; code sync only if `rg` finds refs) |

**Suggested batch size:** 1 token per session for large tokens (`blocks_total` ≥ 8 or `go_refs` ≥ 6); up to 3–5 small tokens per session when each has ≤ 4 blocks and familiar code paths.

**Progress tracking:** After each completed token:

```bash
python3 scripts/mark_impl_deep_sync_complete.py IMPL-TOKEN
# Tier C only (no code/test leads required):
python3 scripts/mark_impl_deep_sync_complete.py IMPL-TOKEN --tier-c
```

That script updates the checklist row, sets `blocks_*` from the sidecar, and appends the token to `MANUAL_DEEP_SYNCED` in [`scripts/update_impl_sync_checklist.py`](../scripts/update_impl_sync_checklist.py).

---

## 5. Repeatable 12-step process (per token)

Expanded from the checklist header. Map to [pseudocode-writing-and-validation.md](pseudocode-writing-and-validation.md) phases A–I where useful.

### Phase A — Discovery (steps 1–3)

| Step | Action |
| --- | --- |
| A1 | Read `tied/implementation-decisions/IMPL-{TOKEN}.yaml` and `IMPL-{TOKEN}-pseudocode.md` |
| A2 | Read `related_decisions`, shared REQ/ARCH neighbors; note merge with adjacent IMPL tokens if one function serves multiple decisions |
| A3 | `rg '\[IMPL-{TOKEN}\]'` repo-wide; build inventory: block ↔ file ↔ function ↔ test |

Normalize `code_locations.files` / `functions` in IMPL YAML via **tied-yaml MCP** or `tied-cli` when entries are stale (wrong file, string-only function names). Prefer MCP for index YAML ([PROC-YAML_EDIT_LOOP]).

### Phase B — Read behavior (steps 4–6)

| Step | Action |
| --- | --- |
| B4 | Read `traceability.tests` and discovered `*_test.go` files |
| B5 | Read production loci (functions, types, adapters) — understand real control flow |
| B6 | List **logical blocks** → one future `## PROCEDURE_NAME` each (not one block per line of code) |

### Phase C — Rewrite sidecar (steps 7–10)

| Step | Action |
| --- | --- |
| C7 | `## Summary contract`: INPUT / OUTPUT / DATA in vocabulary terms |
| C8 | De-Go-ify: use `PROCEDURE`, `IF`/`FOR`/`RETURN`, not host-language syntax |
| C9 | Per block: `## NAME` then lead line then pseudocode body |
| C10 | Remove all bulk boilerplate; fix corrupted `How:` fragments (e.g. `CONFIG_OUTPUT_GROUPING]` typos) |

**Sidecar edit:** Direct edit of `IMPL-*-pseudocode.md` is allowed and preferred for multi-block rewrites. Use `impl_detail_set_essence_pseudocode` / `tied-cli` when setting the whole body from a file.

**Template:** [`../../templates/impl-essence-pseudocode-template.md`](../../templates/impl-essence-pseudocode-template.md)

### Sidecar-canonical YAML projection

Sidecar markdown is **canonical**; IMPL detail and index `implementation_approach.summary` fields are **lean projections** of the sidecar `## Summary contract` section ([ARCH-SPEC_DSL_AND_ORACLE] `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`).

| Step | Action |
| --- | --- |
| 1 | Author or repair behavior in the sidecar (Phase C above)—especially `## Summary contract`. |
| 2 | Run literal three-way sync (Phase D–F): sidecar leads ↔ `// -` comments in Go. |
| 3 | Project summary into YAML: `python3 scripts/sync_tied_yaml_projections.py --apply` from repo root (use `--check` before commit or in CI). |
| 4 | Update other IMPL metadata (`code_locations`, `traceability`, cross-refs) via tied-yaml MCP as in Phase G—**not** by pasting summary prose into YAML by hand. |

**Methodology read-only:** Tokens in `INDEX_ONLY_IMPL_TOKENS` inside `scripts/sync_tied_yaml_projections.py` (`IMPL-MCP_FEEDBACK_TOOLS`, `IMPL-MODULE_VALIDATION`, `IMPL-TIED_FILES`) sync **project index** summaries only; inherited detail under `tied/methodology/` is not edited in client repos ([PROC-TIED_METHODOLOGY_READONLY]).

**Lead format (must match in Go):**

```markdown
- [IMPL-TOKEN] [ARCH-…] [REQ-…] — How: <one specific sentence; no placeholder "implement X per contract">.
```

### Phase D–F — Literal three-way sync (steps 11–15)

For each block **in sidecar order**:

| Step | Action |
| --- | --- |
| D11 | Finalize the sidecar lead text |
| D12 | Place `// - {exact lead}` on the line immediately before the implementing `func` |
| D13 | Same before `Test*` (or dedicated test file comment) when tests cover the block |
| D14 | Byte-compare: sidecar lead === `// -` comment (character-for-character after the `// - ` prefix) |
| D15 | Remove obsolete paraphrased file-level `// [IMPL-…]` banners that duplicate or contradict block leads |

**Do not** run `sync_impl_block_leads.py` until the sidecar and YAML function map are correct.

**Multiple IMPL tokens on one function:** Prefer a **single** lead whose token set matches the block’s primary contract, or split into separate procedures in the sidecar so each `func` gets one lead.

### Phase G — Metadata (steps 16–17)

| Step | Action |
| --- | --- |
| G16 | Update `code_locations`, `traceability`, cross-refs via tied-yaml MCP |
| G17 | `lint_yaml` on changed TIED YAML; fix quoting issues |

### Phase H — Validate (steps 18–20)

| Step | Action |
| --- | --- |
| H18 | Layer A: `tied_validate_consistency` |
| H19 | Layer B: walk [pseudocode-validation-checklist.yaml](pseudocode-validation-checklist.yaml) for this token |
| H20 | `go test ./...` (or packages touched); fix failures before marking complete |

---

## 6. Tier variants

### Tier A / B (runtime IMPL with Go)

- Full phases A–H.
- Every production block with `go_refs` or listed in `code_locations` needs a literal `// -` lead.
- Large tokens: work block-by-block; run `go test` on affected packages mid-pass if refactoring comments across many files.

### Tier C (orders 61–72, doc/process)

- Phases **C7–C10** and **H18–H19** are mandatory (sidecar quality + token comments in pseudocode).
- Phases **D11–D15** apply **only if** `rg '\[IMPL-{TOKEN}\]' --glob '*.go'` finds references.
- Mark complete with `mark_impl_deep_sync_complete.py --tier-c` and note: `Manual deep sync: Tier C sidecar quality (no code comment sync)`.

---

## 7. Per-token exit checklist (copy before marking validated)

Use this list immediately before `mark_impl_deep_sync_complete.py`:

```
[ ] Sidecar: Summary contract + one ## per block; no boilerplate phrases
[ ] Each block: exactly one ^- [IMPL- lead with meaningful How:
[ ] blocks_total counted from sidecar (^- [IMPL- lines)
[ ] For each block with Go implementation: // - lead byte-matches sidecar
[ ] For each block with tests: // - lead on Test* (if applicable)
[ ] No stray "Block implements documented behavior" in sidecar or Go
[ ] IMPL YAML code_locations sane (MCP-updated if needed)
[ ] lint_yaml on changed TIED YAML
[ ] tied_validate_consistency ok
[ ] go test ./... ok (if Go touched)
[ ] mark_impl_deep_sync_complete.py run → MANUAL_DEEP_SYNCED includes token
```

---

## 8. Scripts reference

| Script | When to use | When **not** to use |
| --- | --- | --- |
| [`mark_impl_deep_sync_complete.py`](../../scripts/mark_impl_deep_sync_complete.py) | **After** manual exit criteria pass (deep sync) | As a substitute for manual sync |
| [`mark_impl_logic_audit_complete.py`](../../scripts/mark_impl_logic_audit_complete.py) | **After** REQ audit + gates pass; sets `req_audit_pass` | Before automated audit passes |
| [`run_impl_logic_audit.py`](../../scripts/run_impl_logic_audit.py) | Per-token or `--batch N` automated audit | Substitute for sidecar/code review |
| [`reset_impl_logic_audit_queue.py`](../../scripts/reset_impl_logic_audit_queue.py) | Batch 0 tracking reset (`pending_audit`, P0–P2) | Mid-program without `--preserve-audited` |
| [`mark_all_impl_logic_audit_complete.py`](../../scripts/mark_all_impl_logic_audit_complete.py) | Mark all passing tokens after full audit | When audit not yet green |
| [`update_impl_sync_checklist.py`](../../scripts/update_impl_sync_checklist.py) | Inventory / regen block counts from sidecars | Trusting its `layer_b_validated` for non-`MANUAL_DEEP_SYNCED` rows |
| [`reset_impl_checklist_for_deep_sync.py`](../../scripts/reset_impl_checklist_for_deep_sync.py) | One-time or scoped re-queue of non-pilot rows | On tokens already in `MANUAL_DEEP_SYNCED` |
| `sync_impl_block_leads.py` | Optional assist **after** sidecar + YAML are correct | First step on broken sidecars |
| `rewrite_bulk_sidecar_to_pilot.py` | Emergency boilerplate strip only | Primary quality pass |
| `run_impl_deep_sync_batch.py` | **Do not use** for pilot queue | Entire remaining queue |

**Guards:** `MANUAL_DEEP_SYNCED` and `PILOT_SYNCED` in `update_impl_sync_checklist.py` protect completed rows from bulk overwrites.

---

## 9. Verification commands

```bash
# Boilerplate must be zero in sidecars and Go when project gate runs
rg 'Block implements documented behavior|per resource-management atomic I/O' \
  tied/implementation-decisions/*-pseudocode.md
rg 'Block implements documented behavior|per resource-management atomic I/O' \
  --glob '*.go'

# Count block leads for checklist reconciliation
rg -c '^- \[IMPL-' tied/implementation-decisions/IMPL-TOKEN-pseudocode.md

# Spot-check literal match (example)
# Sidecar line must equal Go line after "// - "
rg -n '// - \[IMPL-TOKEN\]' path/to/file.go

# Tests
go test ./...

# TIED (MCP tool)
tied_validate_consistency
```

---

## 10. Anti-patterns (learned from waves 1–2)

| Anti-pattern | Consequence | Remedy |
| --- | --- | --- |
| Marking `validated` from bulk script | False progress; `blocks_synced: 0` with validated status | Re-run exit checklist; reset row if needed |
| Generic `How: implement X per resource-management…` | Fails [PROC-IMPL_PSEUDOCODE_TOKENS]; untestable alignment | Rewrite How: with observable behavior |
| Duplicate/wrong leads on `NewOutputCollector` vs `AddStdout` | Misaligned MCP/tooling story | One lead per **procedure**; match function actually implementing step |
| File-level `// [IMPL-*]` only | Layer B false positive | Add/replace with block `// -` leads |
| Truncated sidecar from automated rewrite | Lost procedures | Restore from git: `git show :tied/implementation-decisions/IMPL-TOKEN-pseudocode.md` |
| Editing checklist YAML by hand | Invalid YAML / wrong counts | Use `mark_impl_deep_sync_complete.py` |
| Many tokens in one commit without per-token H phase | Hard to review; mixed regression | One token complete per work item; test before next |

---

## 11. LEAP and drift

If production code or tests **differ** from the rewritten sidecar during the pass:

1. Update **IMPL** sidecar and block leads first (source of truth for Track C).
2. Update literal `// -` comments and tests to match.
3. If REQ/ARCH scope changed, propagate **IMPL → ARCH → REQ** in the same work item ([PROC-LEAP]).
4. Re-run H18–H20.

Do not leave sidecar, code, and tests in three different stories.

---

## 12. Project completion gate (after order 72)

All waves done when:

1. All 72 checklist rows: `status: validated` with note `Manual deep sync: pilot-quality…` or Tier C note.
2. `rg 'Block implements documented behavior' tied/implementation-decisions/*-pseudocode.md` → **0 hits**.
3. `rg 'per resource-management atomic I/O' tied/implementation-decisions/*-pseudocode.md` → **0 hits**.
4. Same two greps on `--glob '*.go'` → **0 hits**.
4. Spot-check: random sample of Tier A tokens — sidecar lead equals `// -` before `func` and test.
5. Every completed token ∈ `MANUAL_DEEP_SYNCED` (or `PILOT_SYNCED` for original four).
6. `tied_validate_consistency` clean.
7. `go test ./...` green.

---

## 14. REQ audit template (logic alignment program)

After L0 formal DSL is green, each token needs a **REQ satisfaction audit** before `req_audit_pass: true`.

### Criteria coverage matrix (per block)

| Block (H2) | REQ tokens in lead | satisfaction_criteria item | Mapped STEP/PROCEDURE branch | Gap? |
| --- | --- | --- | --- | --- |
| `FORMAT_LIST_ARCHIVE` | REQ-OUT_002 | … | STEP T001, IF CONTAINS… | — |

**Procedure:**

1. From sidecar block leads, list every `[REQ-*]` token.
2. Load each REQ detail via TIED MCP (`tied://requirement/{token}/detail`) or `tied-cli`.
3. For each `satisfaction_criteria` / `validation_criteria` bullet, trace to at least one `STEP Tnnn:` or `PROCEDURE` branch in that block (or a related block named in the lead).
4. Document gaps in [`impl-pseudocode-improvement-queue.yaml`](impl-pseudocode-improvement-queue.yaml) `issues[]` or fix via LEAP (IMPL → ARCH → REQ).
5. Run automated gates: `python3 scripts/run_impl_logic_audit.py --token IMPL-TOKEN`.
6. Mark complete: `python3 scripts/mark_impl_logic_audit_complete.py IMPL-TOKEN [--logic-level L3]`.

**Automated CRIT-001 corpus pass:** see [crit001-corpus-alignment-runbook.md](crit001-corpus-alignment-runbook.md) for scoped criteria alignment (`--req-criteria-scoped`, scope registry, `mark_crit001_complete.py`).

**Tracking reset (batch start):** `python3 scripts/reset_impl_logic_audit_queue.py` sets all rows to `pending_audit` and `req_audit_pass: false`.

**Do not** set `req_audit_pass` via `seed_logic_verification_tracking.py` — that script refreshes oracle counts only.

---

## 13. Quick reference links

| Resource | Path |
| --- | --- |
| Master checklist | [`impl-pseudocode-sync-checklist.yaml`](impl-pseudocode-sync-checklist.yaml) |
| Pilot sidecar | [`../implementation-decisions/IMPL-LIST_FORMAT_SAFETY-pseudocode.md`](../implementation-decisions/IMPL-LIST_FORMAT_SAFETY-pseudocode.md) |
| Pseudocode Track C | [`pseudocode-writing-and-validation.md`](pseudocode-writing-and-validation.md) |
| Layer B checklist | [`pseudocode-validation-checklist.yaml`](pseudocode-validation-checklist.yaml) |
| TIED YAML index | [`tied-yaml-agent-index.md`](tied-yaml-agent-index.md) |
| Agent session rules | [`AGENTS.md`](../../AGENTS.md) |
| Cursor plan (waves detail) | `.cursor/plans/impl_remaining_deep_sync_a9a6e270.plan.md` (local; optional) |

---

**Last updated:** 2026-06-02 — reflects waves 1–2 pilot process and bulk-migration lessons for bkpdir.
