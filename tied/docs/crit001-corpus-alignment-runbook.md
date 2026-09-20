# CRIT-001 Corpus Alignment Runbook

**Process tokens:** `[PROC-IMPL_CODE_TEST_SYNC]`, `[PROC-PSEUDOCODE_VALIDATION]`  
**Traceability:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`, `[IMPL-SPEC_CTL]`

Execute scoped REQ criteria alignment across all **73** `formal_spec` sidecars. Companion to [impl-deep-sync-agent-guide.md §14](impl-deep-sync-agent-guide.md).

**Operator guide:** [docs/markscope/spec-verification.md](../../docs/markscope/spec-verification.md) Layer D  
**Grammar:** [pseudocode-dsl-grammar.md §9](pseudocode-dsl-grammar.md)

---

## 1. Scope and non-goals

| In scope | Out of scope (unchanged) |
|----------|---------------------------|
| Scoped CRIT-001 (`--req-criteria-scoped`) on 73 sidecars | Release bar: Layer A + L0 + 73/73 `req_audit_pass` + 48/48 oracle |
| Scope registry + optional `CRIT-REQ-*` STEP tags | Clearing or resetting `req_audit_pass` during this program |
| `crit001_pass` tracking on improvement queue | Global CI gating of L4 until promotion decision |

---

## 2. Definitions

| Term | Meaning |
|------|---------|
| **Primary REQ** | IMPL listed in REQ `traceability.implementation` (or sole owner) — all satisfaction criteria in scope unless scope registry narrows |
| **Auxiliary REQ** | REQ cited for cross-traceability but not owner — only `applicable_crit_ids` from scope registry apply; **empty list = traceability-only** |
| **Orphan (CRIT-001)** | In-scope criterion with no STEP/PROCEDURE keyword overlap and no explicit `CRIT-REQ-*` tag |
| **Meta validation waiver** | Validation criteria referencing `.yaml` paths or external performance benchmarks — auto-waived by loader |

Optional STEP suffix for explicit linkage:

```markdown
STEP T001: ENSURE_DIRECTORY_EXISTS(dir)  CRIT-REQ-RESOURCE_MANAGEMENT-001
```

---

## 3. Scoped completion bar (73/73)

**Hybrid policy:**

1. **Primary owner REQ** → all satisfaction criteria must map to sidecar STEP/PROCEDURE/PRE/POST text (sidecar-wide corpus).
2. **Auxiliary REQ** → only criteria in [`tied/spec/req-criteria-scope-registry.yaml`](../spec/req-criteria-scope-registry.yaml) for `(impl, req)`; empty `applicable_crit_ids` = no criteria required.
3. **Mis-attributed REQ** → LEAP: remove from lead/H1 or reclassify in scope registry with `rationale`.

**Corpus target:** `python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-scoped` → **73/73 PASS**.

**Informational debt:** `--req-criteria-strict` (all criteria for every cited REQ) — not gating until L4 promotion.

---

## 4. Per-IMPL workflow (12 steps)

1. Set improvement queue row `status: crit001_in_progress` (optional).
2. Baseline: `python3 scripts/report_crit001_orphans.py --token IMPL-TOKEN --mode scoped`
3. List `[REQ-*]` from block leads and H1.
4. Build §14 matrix: Block × REQ × criterion × STEP/PROCEDURE.
5. Classify each REQ as primary / auxiliary / mis-attributed (compare REQ detail + scope registry).
6. Fix sidecar: expand STEP text, add `CRIT-REQ-*` tags, clean mis-attributed tokens.
7. Update scope registry for auxiliary rows (`scripts/seed_req_criteria_scopes.py` then manual edits; set `manual: true` to preserve).
8. LEAP (IMPL → ARCH → REQ) if ownership is wrong.
9. `go run ./cmd/specctl validate tied/implementation-decisions/IMPL-TOKEN-pseudocode.md`
10. `python3 scripts/run_impl_logic_audit.py --token IMPL-TOKEN --req-criteria-scoped --skip-specctl`
11. `python3 scripts/mark_crit001_complete.py IMPL-TOKEN --mode scoped`
12. Batch gate + release bar: `SKIP_TIED_MCP=1 scripts/run-spec-verification.sh`

### Decision tree

| Situation | Action |
|-----------|--------|
| Keyword gap on primary REQ | Expand STEP/PROCEDURE or add `CRIT-REQ-*` |
| Auxiliary, traceability-only | Scope registry `applicable_crit_ids: []` |
| Auxiliary, partial criteria | Scope registry lists subset + map in sidecar |
| Mis-attributed REQ | Remove from lead/H1 or LEAP split REQ |
| Meta validation criterion | Loader waiver; note in `crit001_scope_notes` if disputed |

---

## 5. Batch schedule

Work by `process_order` from [impl-pseudocode-improvement-queue.yaml](impl-pseudocode-improvement-queue.yaml).

| Batch | Orders | Theme |
|-------|--------|-------|
| 1 | 1–12 | Core runtime, atomic, config foundations |
| 2 | 13–24 | CLI, errors, formatting |
| 3 | 25–36 | Processing, git, fileops |
| 4 | 37–48 | Config fixes, tests |
| 5 | 49–60 | Test IMPLs, extraction |
| 6 | 61–73 | Tier C doc/process + SPEC_CTL |

**Batch gate:**

```bash
python3 scripts/run_impl_logic_audit.py --batch N --req-criteria-scoped --skip-specctl
go run ./cmd/specctl validate tied/implementation-decisions/*-pseudocode.md
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

---

## 6. Tooling reference

| Script | Role |
|--------|------|
| `scripts/seed_req_criteria_scopes.py` | Seed scope registry from REQ/IMPL traceability |
| `scripts/report_crit001_orphans.py` | Per-token or batch orphan report (`--csv`) |
| `scripts/run_impl_logic_audit.py --req-criteria-scoped` | Automated scoped gate |
| `scripts/mark_crit001_complete.py` | Set `crit001_pass: true` |
| `scripts/run_impl_logic_audit.py --mark-crit001 --req-criteria-scoped --batch all` | Batch mark passing tokens |

| Data | Path |
|------|------|
| Scope registry | `tied/spec/req-criteria-scope-registry.yaml` |
| Baseline report | `tied/spec/crit001-baseline-strict.csv` (informational strict debt) |
| Loader | `scripts/req_criteria_loader.py` |

---

## 7. LEAP triggers

Apply **IMPL → ARCH → REQ** when:

- REQ cited but criteria belong to a different REQ token
- REQ `satisfaction_criteria` too monolithic for one IMPL (split REQ or per-IMPL criteria)
- Scope registry would waive >50% of criteria for a “primary” REQ — reclassify ownership in TIED

Use TIED MCP / [tied-yaml skill](../../.cursor/skills/tied-yaml/SKILL.md); run `tied_validate_consistency` after LEAP batches.

---

## 8. L4 promotion path (future)

Do **not** wire to global CI until all of:

1. All 9 oracle `pkg_paths` ≥ `MUTATION_SCORE_THRESHOLD` (`scripts/run-mutation-waves.sh`)
2. Scoped CRIT-001 **73/73** (this runbook)
3. Primary REQs pass `--req-criteria-strict` (informational debt → zero)
4. Contract interpreter green on ≥4 rich IMPLs
5. Acceptable `make lint` duration with L4 pilots

Then `RUN_L4_REQ_STRICT=1` may use strict mode in `spec-corpus-status.sh`.

---

## 9. Corpus status

**Complete (2026-06-03):** 73/73 scoped CRIT-001 PASS after scope registry seed + sidecar alignment on 7 primary-owner tokens.

Gold references:

- [`IMPL-ATOMIC_OPS-pseudocode.md`](../implementation-decisions/IMPL-ATOMIC_OPS-pseudocode.md) — passes strict
- [`IMPL-CONFIG_STRUCT-pseudocode.md`](../implementation-decisions/IMPL-CONFIG_STRUCT-pseudocode.md) — primary owner REQ-CONFIGURATION

---

**Last updated:** 2026-06-03
