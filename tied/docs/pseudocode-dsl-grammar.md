# IMPL Formal Pseudocode DSL Grammar

**Process tokens:** `[PROC-PSEUDOCODE_VALIDATION]`, `[PROC-IMPL_PSEUDOCODE_TOKENS]`, `[PROC-IMPL_CODE_TEST_SYNC]`  
**Traceability:** `[REQ-PSEUDOCODE_FORMAL_VERIFICATION]`, `[ARCH-SPEC_DSL_AND_ORACLE]`, `[IMPL-SPEC_CTL]`

Normative grammar for machine-verifiable `essence_pseudocode` sidecars. This document defines **syntax** (what parses) and **validation profiles** (what CI and logic audit enforce).

**Operational guide:** [docs/markscope/spec-verification.md](../../docs/markscope/spec-verification.md)  
**Logic audit runbook:** [impl-deep-sync-agent-guide.md §14](impl-deep-sync-agent-guide.md)

---

## 1. File structure

| Element | Syntax | Required at L0 (`formal_spec`) |
|---------|--------|--------------------------------|
| File title | `# [IMPL-X] [ARCH-Y] [REQ-Z]` (H1; bracket order IMPL, ARCH, REQ) | Yes |
| Summary block | `## Summary contract` with `INPUT:` / `OUTPUT:` / `DATA:` | Yes |
| Runtime block | `## <UPPER_SNAKE_NAME>` (H2) | ≥1 per IMPL |
| Block lead | `- [IMPL-…] [ARCH-…] [REQ-…] — How: <observable behavior>` | Every H2 except exempt |
| Formal identity | `SPEC-ID: <impl-token>::<block-name>` | Every runtime H2 |
| Preconditions | `PRE: <expr>` | Encouraged; **required** for L3 oracle-rich blocks |
| Postconditions | `POST: <expr>` | Encouraged; **required** for L3 oracle-rich blocks |
| Invariants | `INV: <expr>` | Optional |
| Steps | `STEP <id>: <action>` where `<id>` is `T` + digits (e.g. `T001`) | ≥1 per runtime H2 |
| Branches | `BRANCH <id>: <condition>` (`B001`, optional suffix `B001A`) | Encouraged on error paths |
| Errors | `ERROR <id>: <condition>` (`E001`) | Encouraged on fallible I/O |
| Procedure body | `PROCEDURE NAME(args):` + indented Algol body | ≥1 per runtime H2 (hybrid model) |
| Contracts | `INPUT:`, `OUTPUT:`, `DATA:`, `CONTROL:` at file or block scope | Summary + block scope as needed |
| Test embed | `## EMBEDDED_*` | Optional; waivable via `coverage-waivers.yaml` |

**Exempt H2 blocks (no SPEC-ID):** `Summary contract`, `EMBEDDED_*` test-spec blocks (register `INFRA-*` waiver if matrix fails).

**Anti-patterns (fail logic audit):** boilerplate `Block implements documented behavior`, generic `How:` placeholders, host-language pasted into PROCEDURE bodies.

---

## 2. Expression language (language-agnostic)

- Identifiers: `UPPER_SNAKE`, `camelCase`, or dotted paths (`writer.tempPath`)
- Operators: `==`, `!=`, `AND`, `OR`, `NOT`, `>=`, `<=`
- Calls: `NAME(arg1, arg2)` — oracle helpers include `VALIDATE_PATH(path)`
- Literals: strings in double quotes; booleans `true` / `false`; result token `ok`

**L4 executable subset** (evaluated by `internal/specmodel/contract` for pilot oracles):

- Field access on oracle state: `writer.committed`, `writer.closed`, `writer.HasHandle`
- Boolean connectives: `AND`, `OR`, `NOT`
- Comparisons: `==`, `!=` against literals and `ok`
- Calls: `VALIDATE_PATH(targetPath) == ok`

PROCEDURE bodies use TIED vocabulary (`IF`, `FOR`, `RETURN`, `PROCEDURE`) — not Go/Rust/Python syntax ([PROC-IMPL_PSEUDOCODE_TOKENS]).

---

## 3. Validation profiles

Profiles stack; higher levels assume lower levels pass.

| Profile | Alias | Enforces | Tooling |
|---------|-------|----------|---------|
| `legacy` | — | H2 + PROCEDURE/IF; SPEC-ID optional | Warnings only |
| `formal_spec` | **L0** | SPEC-ID, STEP/PROCEDURE, block leads, registry membership | `specctl validate`, `matrix --check-coverage`, `check-leads` |
| `logic_verified` | **L1–L3** | L0 + three-way `// -` lead copy + REQ audit + target `logic_level` | `run_impl_logic_audit.py`, improvement queue |
| `oracle_verified` | **L3** | L0 + reference oracle in `internal/specmodel` + conformance tests | `oracle-registry.yaml`, `test/specconformance` |
| `mutation_verified` | **L4** *(projected)* | L3 + mutation score on implementing package | `run-mutation-pilot.sh` |

Set `formal_spec: true` in [`impl-pseudocode-improvement-queue.yaml`](impl-pseudocode-improvement-queue.yaml) when an IMPL passes L0. Set `req_audit_pass: true` only after manual REQ criteria mapping and `mark_impl_logic_audit_complete.py`.

**Current corpus (bkpdir):** **73/73** `formal_spec`, **73/73** `req_audit_pass`, **48/48** oracle domains, logic levels L3 **54** / L2 **7** / L1 **12**.

---

## 4. Profile variants (authoring)

| Variant | When | Required lines | Matrix / leads / oracle |
|---------|------|----------------|-------------------------|
| **Rich** | State machines, I/O, errors (`ATOMIC_OPS`, `FILE_OPERATIONS`, `STRUCTURED_ERRORS`) | SPEC-ID, STEP, PRE/POST, BRANCH/ERROR | Full check-leads + matrix; **L3 oracle** |
| **Standard** | Most Tier A/B runtime | SPEC-ID, STEP, PROCEDURE, block lead | Full check-leads + matrix; L3 where registry lists token |
| **Light** | Tier C doc/process (orders 61–72) | SPEC-ID, STEP, block lead, PROCEDURE or NOTE | Matrix/check-leads waiver when no `.go` refs; **L1 target** |
| **Meta** | `IMPL-SPEC_CTL` | SPEC-ID, STEP on tooling blocks | `INFRA-*` waivers; L3 via specmodel |

Gold references: [`IMPL-ATOMIC_OPS-pseudocode.md`](../implementation-decisions/IMPL-ATOMIC_OPS-pseudocode.md) (rich), [`IMPL-CFG_006-pseudocode.md`](../implementation-decisions/IMPL-CFG_006-pseudocode.md) (minimum viable), [`IMPL-LIST_FORMAT_SAFETY-pseudocode.md`](../implementation-decisions/IMPL-LIST_FORMAT_SAFETY-pseudocode.md) (block-lead sync).

---

## 5. Three-way alignment (L1+)

For Tier A/B IMPLs with Go implementations:

1. Sidecar block lead (line starting `- [IMPL-`) is **byte-identical** to production `// - {lead}` immediately before the implementing `func`.
2. Same lead on covering `Test*` when tests exercise the block.
3. `specctl check-leads` emits no `TRACE-002`.
4. `run_impl_logic_audit.py` reports no `missing_in_go` / `extra_in_go` for tokens with `go_refs > 0`.

Sidecar is authoritative on behavior drift; propagate with LEAP (IMPL → ARCH → REQ) when code or tests differ.

---

## 6. REQ audit (L1+)

For each `[REQ-*]` in a block lead:

1. REQ must exist in `tied/requirements.yaml` (or methodology merge).
2. Each `satisfaction_criteria` / `validation_criteria` item in the REQ detail must trace to ≥1 `STEP` or `PROCEDURE` branch in that block (or related block named in the lead).
3. Document gaps in improvement queue `issues[]` or fix via LEAP.

Automated pre-check: `run_impl_logic_audit.py` (REQ token resolution, boilerplate, structure). Manual sign-off: `mark_impl_logic_audit_complete.py`.

---

## 7. Tool mapping

| Checklist ID | Command | Layer |
|--------------|---------|-------|
| PARSE-001, PARSE-002 | `specctl validate` | L0 |
| SHAPE-001, SHAPE-002 | `specctl validate` | L0 |
| RESOLVE-001, RESOLVE-002 | `specctl validate` | L0 |
| TRACE-001–003, COVER-001–002 | `specctl matrix --check-coverage` | L0 |
| Block leads (TRACE-002) | `specctl check-leads` | L0 |
| Three-way + REQ structure | `run_impl_logic_audit.py` | L1–L3 |
| REQ criteria orphans (CRIT-001) | `run_impl_logic_audit.py --req-criteria-scoped` | L4 corpus (73/73) |
| REQ criteria strict debt | `report_crit001_orphans.py --mode strict` | Informational |
| Oracle conformance | `go test ./test/specconformance/...` | L3 |
| TIED indexes + pseudocode comments | `tied_validate_consistency` | Layer A |
| Mutation pilot | `run-mutation-pilot.sh` | L4 (projected) |

CI gate: `SKIP_TIED_MCP=1 scripts/run-spec-verification.sh` → `scripts/spec-corpus-status.sh`.

---

## 8. EBNF (subset)

```ebnf
file        ::= h1_line block* ;
block       ::= h2_line block_lead? spec_meta? contract* body* ;
spec_meta   ::= spec_id | pre_line | post_line | inv_line | step_line | branch_line | error_line ;
spec_id     ::= "SPEC-ID:" spec_ref ;
step_line   ::= "STEP" step_id ":" action ;
block_lead  ::= "- [" impl "] [" arch "] [" req "]" — How: " text ;
h2_line     ::= "##" name ;
```

Reference parser: `internal/specparse`. Normalized AST schema: `tied/spec/grammar/schema.json`.

---

## 9. L4 extensions (pilot — opt-in, not global CI)

Pilot tooling exists; promotion to global CI gate is documented separately.

| Extension | Status | Tooling |
|-----------|--------|---------|
| REQ criteria IDs (`CRIT-REQ-…` on STEP lines) | Optional authoring | `scripts/req_criteria_loader.py`, `--req-criteria-strict` on `run_impl_logic_audit.py` |
| Orphan criteria diagnostic `CRIT-001` | Corpus complete (scoped) | `python3 scripts/run_impl_logic_audit.py --req-criteria-scoped` |
| Executable PRE/POST subset | Pilot (`IMPL-ATOMIC_OPS`) | `internal/specmodel/contract`, `AssertBlockContracts` in `test/specconformance` |
| Mutation score by package | Pilot waves 0–4 | `scripts/run-mutation-pilot.sh`, `scripts/run-mutation-waves.sh`, `tied/spec/mutation-wave-registry.yaml` |
| Profile manifest export | Planned | `{ formal_spec, logic_level, oracle, mutation_score }` per IMPL |

Optional STEP suffix for explicit criteria linkage:

```markdown
STEP T001: ENSURE_DIRECTORY_EXISTS(dir)  CRIT-REQ-RESOURCE_MANAGEMENT-001
```

**Promotion bar (future global L4 CI):** all 9 oracle `pkg_paths` ≥ `MUTATION_SCORE_THRESHOLD`; strict REQ audit clean on L3 tokens; contract interpreter green on ≥4 rich IMPLs; acceptable `make lint` duration.

**Runbook:** [crit001-corpus-alignment-runbook.md](crit001-corpus-alignment-runbook.md)

Until promoted, **strongest enforced profile = Layer A + L0 + manual REQ audit + L3 oracle (release bar)**; scoped CRIT-001 is L4 corpus-complete (73/73).

---

**Last updated:** 2026-06-03 — L4 pilot tooling (contract eval, CRIT-001, mutation waves); release bar unchanged.
