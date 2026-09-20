# !/usr/bin/env Markscope

- MARKSCOPE_CLIENT_DIR: ${MARKSCOPE_CLIENT_DIR}
- MARKSCOPE_DOCS_DIR: ${MARKSCOPE_DOCS_DIR}
- MARKSCOPE_SCRIPTS_DIR: ${MARKSCOPE_SCRIPTS_DIR}
<!-- do not generate stdout, it is part of every eval -->
```bash @on-load @hide
cd "$MARKSCOPE_DOCS_DIR/../.."
cd /Users/fareed/Documents/dev/go/bkpdir/docs/markscope/../..
source ~/.bash_profile
```
```bash @g:trace @hide
source "${MARKSCOPE_SCRIPTS_DIR:?}/trace-debug.sh"
```

# Formal pseudocode verification (specctl)

Machine-checkable IMPL sidecars across **four validation layers**:

| Layer | Scope | Tooling |
|-------|--------|---------|
| **A — TIED** | Index/detail YAML, token comments in pseudocode, REQ↔ARCH↔IMPL graph | `tied_validate_consistency` |
| **B — L0 formal DSL** | Parse, SPEC-ID/STEP shape, matrix coverage, literal block-lead copy | `specctl`, `run-spec-verification.sh` |
| **C — Logic L1–L3** | Sidecar ↔ code ↔ tests alignment, REQ audit, oracle conformance | `run_impl_logic_audit.py`, `oracle-registry.yaml` |
| **D — Projected** | STEP-behavior test strength, mutation score, semantic PRE/POST proof | mutation pilot, future oracle interpreter |

**Tokens:** `REQ-PSEUDOCODE_FORMAL_VERIFICATION` · `ARCH-SPEC_DSL_AND_ORACLE` · `IMPL-SPEC_CTL`  
**Grammar:** [tied/docs/pseudocode-dsl-grammar.md](../../tied/docs/pseudocode-dsl-grammar.md)  
**Logic runbook:** [tied/docs/impl-deep-sync-agent-guide.md §14](../../tied/docs/impl-deep-sync-agent-guide.md)

---

## All-green corpus (73/73 — current bar)

The **strongest enforced bar** today is **Layer B (L0) + Layer C tracking** on all 73 sidecars:

| Check | All-green means | Verified by |
|-------|-----------------|-------------|
| Corpus size | **73** sidecars == **73** registry entries | `spec-corpus-status` |
| Parse + shape (L0) | **0** `[error]`, **0** `[warning]` from validate | `specctl validate` |
| Traceability (L0) | No `COVER-001` on formal blocks | `matrix --check-coverage` |
| Block leads (L0) | No `TRACE-002` (Tier A/B Go sync; Tier C waived) | `check-leads` |
| Conformance (L3) | `specparse`, `specmodel`, `specconformance` tests pass | `go test` |
| Oracle (L3) | **48/48** runtime IMPLs in [`oracle-registry.yaml`](../../tied/spec/oracle-registry.yaml) | `spec-corpus-status` |
| Logic audit (L1–L3) | **73/73** `logic_verified` + **73/73** `req_audit_pass` (manual gate) | improvement queue |
| CRIT-001 scoped (L4) | **73/73** `crit001_pass` (scope registry + `--req-criteria-scoped`) | improvement queue + runbook |
| Three-way sync | Sidecar `// -` leads match Go for Tier A/B (`run_impl_logic_audit.py`) | logic audit script |
| Tracking | **73/73** `formal_spec: true` | improvement queue |
| Layer A (release) | `tied_validate_consistency` with pseudocode | `validate-tied-mcp` |

**Current logic-level distribution (73/73):** L3 **54**, L2 **7**, L1 **12** (Tier C doc/process + selected Tier B tokens).

**One-click status** — prints PASS/FAIL per row and `CORPUS STATUS: ALL GREEN` when every Layer B check passes:

```bash @n:corpus-status @r:trace
scripts/spec-corpus-status.sh
```

**Full Layer B gate** — same checks as CI plus conformance tests; ends with the status summary:

```bash @n:skip-tied-mcp-run-spec-verification @r:trace
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

**Layer A + Layer B (release)**:

```bash @n:validate-tied-mcp @r:trace
scripts/validate-tied-mcp.sh
```

When all-green, the status script exits **0** and `run-spec-verification.sh` prints:

```text
CORPUS STATUS: ALL GREEN
  Profile: formal_spec on all 73 IMPL sidecars
  Gate:    SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

Any non-zero exit or `FAIL` line means the corpus is not all-green — fix before merge.

---

## Logic verification levels (L0–L3 + projected L4)

Validation strength increases by level. **L0 and corpus-wide logic audit are complete (73/73).** L3 oracle coverage is complete for the **48** runtime domains in the oracle registry.

| Level | Meaning | Enforced by | Status (73 IMPLs) |
|-------|---------|-------------|-------------------|
| **L0** | `formal_spec`: SPEC-ID, STEP/PROCEDURE, block leads parse; matrix + check-leads | `specctl`, CI gate | **73/73 done** |
| **L1** | Meaningful `How:` leads; REQ tokens resolve; sidecar structure (Summary + blocks); Tier C target | `run_impl_logic_audit.py` | **12 at L1** (Tier C) |
| **L2** | Tests assert STEP behavior (not comment refs only); test-spec sidecars match tests | conformance + manual review | **7 at L2** |
| **L3** | `internal/specmodel` oracle + `test/specconformance` vs `pkg/*` | oracle registry + `go test` | **54 at L3**; **48/48** oracle domains |
| **L4** *(projected)* | Mutation score threshold; automated REQ criteria ↔ STEP mapping; PRE/POST interpreter | mutation pilot, future tooling | pilot on `pkg/fileops` only |

### REQ audit (required for `req_audit_pass: true`)

Per [impl-deep-sync-agent-guide §14](impl-deep-sync-agent-guide.md) and [crit001-corpus-alignment-runbook.md](crit001-corpus-alignment-runbook.md): for each block lead, every linked REQ `satisfaction_criteria` / `validation_criteria` item must map to at least one `STEP` or `PROCEDURE` branch. **Scoped mode** limits criteria to primary-owner REQs or scope-registry subsets. Gaps → LEAP or `issues[]` in the improvement queue.

**Do not** treat `seed_logic_verification_tracking.py` as proof of REQ audit — it refreshes oracle counts only. Use:

```bash
python3 scripts/run_impl_logic_audit.py --token "$IMPL_TOKEN"
python3 scripts/mark_impl_logic_audit_complete.py "$IMPL_TOKEN" [--logic-level L3] [--tier-c]
```

Batch audit: `python3 scripts/run_impl_logic_audit.py --batch all`  
Reset tracking: `python3 scripts/reset_impl_logic_audit_queue.py`

---

| Goal | Command | When |
|------|---------|------|
| **Corpus health snapshot (fast)** | @code-as-button(corpus-status) | L0 + oracle 48/48 + logic 73/73 + req_audit 73/73 |
| **Per-token logic audit** | `python3 scripts/run_impl_logic_audit.py --token IMPL-TOKEN` | Before marking `req_audit_pass` |
| **Mark logic audit complete** | `python3 scripts/mark_impl_logic_audit_complete.py IMPL-TOKEN` | After audit + specctl gates |
| **Seed oracle counts only** | @code-as-button(seed-logic-tracking) | After oracle registry edits (not REQ audit) |
| **Full local gate (fast)** | @code-as-button(skip-tied-mcp-run-spec-verification) | After editing sidecars; no Node/MCP needed |
| **Full gate + TIED Layer A** | @code-as-button(validate-tied-mcp) | Before commit; needs built `mcp-server` |
| **Layer A only** | @code-as-button(tied-validate-consistency) | Indexes, detail YAML, pseudocode token comments |
| **One IMPL sidecar only** | @code-as-button(validate-sidecar) | Narrow fix loop |
| **Coverage for formal IMPLs** | @code-as-button(check-coverage-all-sidecars) | After tests/code exist |
| **Block-lead copy audit** | @code-as-button(check-leads-sidecar) | After syncing `//` comments in `pkg/`, `internal/`, `cmd/`, or `test/` (or repo-root `*.go`) |
| **Property / model tests** | `go test ./test/specconformance/... -count=1` | Oracle vs `internal/specmodel` |
| **Mutation strength (L4 pilot)** | @code-as-button(run-mutation-pilot) | Requires `go-mutesting` on `PATH` |
| **L4 REQ scoped audit** | `python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-scoped` | CRIT-001 scoped check (corpus bar) |
| **L4 REQ strict debt** | `python3 scripts/report_crit001_orphans.py --batch all --mode strict --csv` | Informational; see `tied/spec/crit001-baseline-strict.csv` |
| **L4 corpus + pilots** | `RUN_L4_PILOT=1 scripts/spec-corpus-status.sh` | Mutation wave + optional REQ strict |
| **L4 mutation waves 0–4** | `scripts/run-mutation-waves.sh` | Sequential rollout (opt-in) |

---

## Recommended workflows

### Daily / PR check (simplest)

```bash @n:skip-tied-mcp-run-spec-verification @r:trace
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
```

Pipeline: `specctl validate` → `matrix` → `matrix --check-coverage` → `check-leads` → `go test` spec packages → **`spec-corpus-status.sh`** (includes `req audit pass` row).

For logic drift after code changes: `python3 scripts/run_impl_logic_audit.py --batch all` before merge.

### Release / TIED-complete check

```bash @n:validate-tied-mcp @r:trace
scripts/validate-tied-mcp.sh
```

Runs: `yaml_index_validate` + `tied_validate_consistency` (Node MCP), then the same spec script as above. Fails if either layer fails.

### After editing one sidecar

```ux
act: :allow
allow:
  command: 'ls -1 tied/implementation-decisions/IMPL-*-pseudocode.md'
  separator: "\n"
init: false
name: SIDECAR
```

```bash @n:validate-sidecar @r:trace @req(SIDECAR)
go run ./cmd/specctl validate "$SIDECAR"
```
```bash @n:check-leads-sidecar @r:trace @req(SIDECAR)
go run ./cmd/specctl check-leads "$SIDECAR"
```

Add `matrix --check-coverage` only if that IMPL is listed under `formal_spec` in the registry.

### CRIT-001 scoped audit (L4 corpus)

Scoped mode checks only **in-scope** REQ criteria per [`req-criteria-scope-registry.yaml`](../../tied/spec/req-criteria-scope-registry.yaml): primary-owner REQs need full coverage; auxiliary cites use the registry subset (empty list = traceability-only). Runbook: [`crit001-corpus-alignment-runbook.md`](../../tied/docs/crit001-corpus-alignment-runbook.md).

**Corpus-wide gate** — all 73 sidecars; `--skip-specctl` keeps the loop fast (structure + CRIT-001 only):

```bash
python3 scripts/run_impl_logic_audit.py --batch all --req-criteria-scoped --skip-specctl
```

**Per-token diagnostic** — before fixing one IMPL; lists orphan criteria with REQ id and suggested action:

```bash @r:trace
.cursor/skills/tied-yaml/scripts/tied-cli.sh yaml_index_list_tokens '{"index":"implementation"}' | yq '.[] | .' | head
```
```ux
allow:
  command: >-
    .cursor/skills/tied-yaml/scripts/tied-cli.sh yaml_index_list_tokens '{"index":"implementation"}' | yq '.[] | .'
  separator: "\n"
init: false
name: IMPL_TOKEN
```
$(echo $IMPL_TOKEN | hexdump -C)
```bash @req(IMPL_TOKEN) @r:trace
python3 scripts/report_crit001_orphans.py --token "$IMPL_TOKEN" --mode scoped
```

**L4 dashboard** — extends `spec-corpus-status.sh` with scoped CRIT-001 re-check and optional mutation pilot (not part of the release bar unless promoted):

```bash
RUN_L4_PILOT=1 RUN_L4_REQ_STRICT=1 scripts/spec-corpus-status.sh
```

Mark a token complete after scoped audit passes: `python3 scripts/mark_crit001_complete.py "$IMPL_TOKEN" --mode scoped`.

---

## `specctl` subcommands

Build/run from repo root:

```bash @readonly
go run ./cmd/specctl validate  [paths...]
go run ./cmd/specctl matrix    [paths...]
go run ./cmd/specctl matrix --check-coverage [paths...]
go run ./cmd/specctl check-leads [paths...]
```

**Paths:** Glob or explicit files. Default if omitted: `tied/implementation-decisions/*-pseudocode.md`.

**Examples:**

```bash
# All sidecars — parse + schema + symbols (formal rules only for registry IMPLs)
go run ./cmd/specctl validate tied/implementation-decisions/*-pseudocode.md

# Traceability table (stdout: IMPL, SPEC-ID, STEP, test/code flags)
go run ./cmd/specctl matrix tied/implementation-decisions/IMPL-CFG_006-pseudocode.md

# Fail on uncovered formal blocks (COVER-001)
go run ./cmd/specctl matrix --check-coverage "$SIDECAR"
```

```bash @n:check-coverage-all-sidecars @r:trace
go run ./cmd/specctl matrix --check-coverage tied/implementation-decisions/*-pseudocode.md
```

---

## Environment variables

| Variable | Effect |
|----------|--------|
| `SKIP_TIED_MCP=1` | `run-spec-verification.sh` does not call `validate-tied-mcp.sh` |
| `MUTATION_SCORE_THRESHOLD=80` | Mutation pilot pass threshold (default `70`) |
| `RUN_L4_PILOT=1` | `spec-corpus-status.sh` runs L4 pilot section (mutation + optional REQ strict) |
| `RUN_L4_REQ_STRICT=1` | With `RUN_L4_PILOT=1`, run `--req-criteria-scoped` CRIT-001 audit |
| `MUTATION_WAVE=0`–`4` / `all` | Wave filter for `run-mutation-pilot.sh` (see `tied/spec/mutation-wave-registry.yaml`) |
| `MUTATION_FROM_REGISTRY=1` | Mutate all unique `pkg_paths` from oracle registry |
| `MUTATION_PKG=pkg/errors` | Single-package mutation target |

---

## Projected strongest bar (L4 roadmap)

**Not globally CI-gating.** Pilot commands below; release bar remains Layer A + L0 + 73/73 + 48/48 + `go test ./...`.

| Capability | Pilot enforcement | Current state |
|------------|-------------------|---------------|
| REQ criteria ↔ STEP auto-map | `--req-criteria-scoped` → CRIT-001 orphans (73/73 corpus) | `scripts/req_criteria_loader.py`, scope registry, runbook |
| L2 STEP-behavior tests | Every L3 block has test asserting STEP outcome | Comment/matrix refs; partial table tests |
| PRE/POST semantic proof | `AssertBlockContracts` + `internal/specmodel/contract` | Pilot: `IMPL-ATOMIC_OPS` POST/BRANCH in `atomic_pbt_test.go` |
| Mutation testing (L4) | `MUTATION_SCORE_THRESHOLD` per wave/pkg | Waves 0–4 in `mutation-wave-registry.yaml`; `run-mutation-waves.sh` |
| L4 corpus row | `RUN_L4_PILOT=1 scripts/spec-corpus-status.sh` | Fifth section (SKIP unless env set) |

### L4 promotion bar (future global CI)

All of the following before wiring L4 into `run-spec-verification.sh`:

1. All **9** oracle `pkg_paths` ≥ `MUTATION_SCORE_THRESHOLD` under `run-mutation-waves.sh`
2. `--req-criteria-scoped` clean on all 73 sidecars (corpus complete)
3. Primary REQs pass `--req-criteria-strict` (strict debt report → zero)
4. Contract interpreter green on ≥4 rich-profile IMPLs
5. `make lint` / `validate-tied-mcp` duration acceptable with L4 pilots enabled

### Mutation waves (opt-in)

| Wave | Packages |
|------|----------|
| 0 | `pkg/fileops` |
| 1 | `pkg/errors` |
| 2 | `pkg/git`, `pkg/cli` |
| 3 | `pkg/resources`, `pkg/processing`, `pkg/testutil` |
| 4 | `pkg/config`, `pkg/formatter` |

Registry: [`tied/spec/mutation-wave-registry.yaml`](../../tied/spec/mutation-wave-registry.yaml)

When L4 is promoted globally, extend the default `spec-corpus-status.sh` exit criteria; until then, **release bar = Layer A + all-green Layer B/C rows above**.

---

## Mutation pilot (L4 — optional today)

```bash @r:trace
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```
```bash @n:run-mutation-pilot @r:trace
scripts/run-mutation-pilot.sh
```
```bash @n:run-mutation-waves @r:trace
scripts/run-mutation-waves.sh
```

Wave / registry examples:

```bash
MUTATION_WAVE=1 scripts/run-mutation-pilot.sh
MUTATION_FROM_REGISTRY=1 scripts/run-mutation-pilot.sh
RUN_L4_PILOT=1 RUN_L4_REQ_STRICT=1 scripts/spec-corpus-status.sh
```

If `go-mutesting` is missing, the script exits **0** with waiver `INFRA-MUTATION-TOOL` (pilot only targets `./pkg/fileops/...`).

---

## Conformance tests

```bash
go test ./internal/specparse/... ./internal/specmodel/... ./test/specconformance/... -count=1
```

Oracle: `internal/specmodel` (e.g. `AtomicWriter` state machine). Property tests: `test/specconformance` (gopter vs `pkg/fileops`).

---

## Layer A — TIED (without full validate script)

```bash @n:tied-validate-consistency @r:trace
.cursor/skills/tied-yaml/scripts/tied-cli.sh tied_validate_consistency \
  '{"include_detail_files":true,"include_pseudocode":true,"require_detail_record":true}'
```

Or:

```bash
node scripts/validate-tied-mcp.mjs
```

Checks indexes, detail YAML, and **token comments** in `essence_pseudocode` — not behavioral equivalence with Go code.

---

## Data and configuration

| Resource | Path | Role |
|----------|------|------|
| Sidecar source | `tied/implementation-decisions/IMPL-{TOKEN}-pseudocode.md` | Authoritative pseudocode body |
| Formal IMPL list | `tied/spec/formal-spec-registry.yaml` | **73** tokens under `formal_spec:` |
| Coverage waivers | `tied/spec/coverage-waivers.yaml` | `INFRA-*` blocks exempt from coverage; `lead_check_skip` exempt from check-leads |
| Status script | `scripts/spec-corpus-status.sh` | One-screen all-green dashboard (exit 0 = pass) |
| Gate script | `scripts/run-spec-verification.sh` | Full Layer B pipeline + status summary |
| Oracle registry | `tied/spec/oracle-registry.yaml` | **48** L3 domains (`status: verified`) |
| Logic audit | `scripts/run_impl_logic_audit.py` | Per-token / batch structural + lead sync audit |
| Logic complete | `scripts/mark_impl_logic_audit_complete.py` | Sets `req_audit_pass` after gates (manual sign-off); `--logic-level L4` sets `mutation_verified` |
| Mutation tracking | `scripts/seed_mutation_tracking.py` | Adds `mutation_verified` / `mutation_score` fields (not proof) |
| Mutation waves | `scripts/run-mutation-waves.sh` | Sequential L4 wave rollout |
| Wave registry | `tied/spec/mutation-wave-registry.yaml` | Packages per wave 0–4 |
| REQ criteria loader | `scripts/req_criteria_loader.py` | CRIT-001 orphan detection (scoped + strict) |
| CRIT-001 runbook | `tied/docs/crit001-corpus-alignment-runbook.md` | Corpus alignment program |
| Scope registry | `tied/spec/req-criteria-scope-registry.yaml` | Primary/auxiliary REQ criteria scope |
| Strict debt CSV | `tied/spec/crit001-baseline-strict.csv` | Informational full-strict orphans |
| Queue reset | `scripts/reset_impl_logic_audit_queue.py` | Batch 0: `pending_audit`, P0–P2 priorities |
| Oracle refresh | `scripts/seed_logic_verification_tracking.py` | Oracle counts only — **not** REQ audit proof |
| AST JSON schema | `tied/spec/grammar/schema.json` | Normalized block shape (tooling) |
| Layer B checklist | `tied/docs/pseudocode-validation-checklist.yaml` | PARSE/SHAPE/RESOLVE/COVER codes |
| Sync checklist | `tied/docs/impl-pseudocode-sync-checklist.yaml` | Per-token `formal_spec`, `req_audit_pass` (73/73) |
| Improvement queue | `tied/docs/impl-pseudocode-improvement-queue.yaml` | `logic_level`, `target_logic_level`, `priority`, `audit_notes` |

**Adding a new formal IMPL:** add `SPEC-ID` / `STEP` per grammar, add token under `formal_spec:` in the registry, sync Go refs, update tracking YAML, then run @code-as-button(corpus-status).

---

## Output reference

Diagnostics print to stdout:

```text
[severity] CODE file:line block=NAME: message
```

| Code | Severity | Meaning |
|------|----------|---------|
| `PARSE-001` | error | Sidecar did not parse or has no H2 blocks |
| `SHAPE-001` | error/warning | Missing `SPEC-ID`, `STEP`/`PROCEDURE`, or block lead |
| `RESOLVE-002` | error | Duplicate `SPEC-ID` or `STEP` |
| `COVER-001` | error | Formal block has no matching Go symbol/comment (with `--check-coverage`) |
| `TRACE-002` | error | Block lead not found literally in Go `//` comments |
| `CRIT-001` | error *(L4 pilot)* | REQ satisfaction/validation criterion has no STEP/PROCEDURE coverage (`--req-criteria-strict`) |

`matrix` tab-separated columns:

```text
IMPL-TOKEN    SPEC-ID    STEP_ID    test=true|false    code=true|false
```

Exit code **0** = pass; **1** = at least one error-severity finding (validate/matrix/check-leads/status script).

---

## Related docs

- [tied/docs/pseudocode-writing-and-validation.md](../../tied/docs/pseudocode-writing-and-validation.md) — three-way alignment, Layer A/B order  
- [tied/docs/pseudocode-format-and-practices.md](../../tied/docs/pseudocode-format-and-practices.md) — vocabulary and formal IDs  
- [tied/docs/pseudocode-dsl-grammar.md](../../tied/docs/pseudocode-dsl-grammar.md) — normative formal DSL + validation profiles (L0–L4)  
- [tied/docs/impl-deep-sync-agent-guide.md](../../tied/docs/impl-deep-sync-agent-guide.md) — Track C sync, REQ audit §14, logic scripts
