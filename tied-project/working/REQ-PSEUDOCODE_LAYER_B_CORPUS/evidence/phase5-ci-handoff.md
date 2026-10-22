# Continuation handoff — Phase 5 CI (Path B) + corpus tail

**Parent slice:** `REQ-PSEUDOCODE_LAYER_B_CORPUS`  
**Branch:** `chore/tied-two-folder-cleanup`  
**Prior wave close-out:** [`w3-close-out-handoff.md`](w3-close-out-handoff.md) (feature `9adc16b…`)  
**Session commits on branch:** `ed88dc9` (D15 + incremental sidecar lead fix)

## Corpus (Track B)

| Item | Status |
|------|--------|
| W0–W3 Layer B waves | **Complete** — 73/73 in `tied-project/spec/layer-b-wave-registry.yaml` |
| W3 tracker | `operator_evidence.status: w3_slice_close_out` |
| L0 D15 (`comparison.go`) | **Fixed** in `ed88dc9` |
| L0 full `run-spec-verification.sh` | Passes through check-leads + spec conformance; **`spec-corpus-status.sh`** may still fail until improvement queue YAML exists under `tied-bundle/docs/` (see below) |

### Re-verify (baseline)

```bash
SKIP_TIED_MCP=1 scripts/run-spec-verification.sh
TIED_BASE_PATH="$PWD/tied-project" scripts/run-pseudocode-layer-b-batch.sh \
  --mode full --wave w3 --enforce \
  --out-dir tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-full
TIED_BASE_PATH="$PWD/tied-project" .cursor/skills/tied-yaml/scripts/tied-cli.sh tied_validate_consistency
```

## Phase 5 — Path B CI (implemented, **not committed**)

Fail-closed Layer B on PRs: clone TIED store, build MCP, enforce changed sidecars in `enforce_in_ci: true` waves.

| Artifact | Role |
|----------|------|
| [`.github/workflows/spec-verification.yml`](../../../../.github/workflows/spec-verification.yml) | `l0-spec` + `layer-b-enforce` jobs |
| [`scripts/ci/setup-tied-mcp.sh`](../../../../scripts/ci/setup-tied-mcp.sh) | Clone `fareedst/tied`, `npm run build` in `mcp-server/`, export `TIED_MCP_BIN` / methodology path |
| [`scripts/layer_b_ci_enforcement.py`](../../../../scripts/layer_b_ci_enforcement.py) | Diff + registry → enforce token list |
| [`scripts/run-layer-b-ci-enforcement.sh`](../../../../scripts/run-layer-b-ci-enforcement.sh) | Batch `pseudocode_validate` full + `--enforce` |
| [`scripts/run-spec-verification.sh`](../../../../scripts/run-spec-verification.sh) | `RUN_LAYER_B=1` / `RUN_LAYER_B_IF_CHANGED=1` hooks |
| [`tied-project/spec/layer-b-wave-registry.yaml`](../../../spec/layer-b-wave-registry.yaml) | `enforce_in_ci: true` on w0–w3 |
| [`docs/markscope/spec-verification.md`](../../../../docs/markscope/spec-verification.md) | Path B env + Actions table |

**Local Path B smoke:**

```bash
bash scripts/ci/setup-tied-mcp.sh   # or TIED_CHECKOUT=/path/to/tied clone
export TIED_BASE_PATH="$PWD/tied-project"
export TIED_METHODOLOGY_BUNDLE_PATH="$TIED_REPO_ROOT/tied-bundle"
RUN_LAYER_B_IF_CHANGED=1 scripts/run-layer-b-ci-enforcement.sh --if-changed
```

**Tests:** `python3 scripts/test_layer_b_ci_enforcement.py`

### Sponsor actions (Path B operational)

1. **Commit** Phase 5 file set (separate from W3 close-out evidence if you want a clean PR).
2. **Branch protection** — on `main`, require **only** the L0 status check until PR path coverage matches how you merge (see [`docs/markscope/spec-verification.md`](../../../../docs/markscope/spec-verification.md) § Branch protection). Required check label is usually **`spec-verification / L0 specctl + conformance`** (YAML job id `l0-spec`). **Do not** require **`layer-b-enforce`** yet — it stays advisory on PRs that trigger the workflow. After paths are broadened or L0 is split to run on every PR, add **`spec-verification / Layer B enforce (fail-closed)`** as required.
3. **Pin TIED in CI** — optional `TIED_GIT_REF` in workflow or `setup-tied-mcp.sh` if `main` drift breaks validate.
4. **Cache** — add Actions cache for `mcp-server/node_modules` if build time hurts.

## Uncommitted working tree (this branch)

| Area | Notes |
|------|--------|
| Phase 5 scripts + `.github/` | Ready to commit as one logical unit |
| W3 close-out evidence | `request-evidence-envelope.v1.json`, handoffs, `close-out-gates-run-w3-postcommit.json` |
| `pseudocode-validate-ci/` | Regenerable CI report dir (optional commit or gitignore) |

## Still open vs modernization plan

| Item | Notes |
|------|--------|
| **Improvement queue refresh** | Restore `tied-bundle/docs/impl-pseudocode-improvement-queue.yaml` + sync checklist (deleted in `9271bab` restructure; was under `tied/docs/`). Restore from git parent or `impl_pseudocode_remediation.py seed-queue` after checklist exists. |
| **W3 / REQ slice close-out commit** | Machine evidence on disk; not all committed |
| **Grammar-Version v2 on touch** | 6/73 sidecars; policy unchanged |
| **W4 non-formal IMPLs** | Not in wave registry |
| **Deep Tier A contracts** | W1 mechanical `PRE/POST`; optional deepening + failure-path tests |
| **5 advisory envelope gaps** | Inquiry under advisory policy |

## Costly choices (recorded)

- **Path B chosen** — PR green implies Layer B MCP for touched enforced sidecars (infra: TIED clone + npm per job until cached).
- **All waves `enforce_in_ci: true`** — full 73-token blast radius on sidecar edits once Layer B is a **required** branch check (deferred; L0-only protection first).

## Close-out (Phase 5 commit)

**Commit:** _(filled post-commit)_  
**Evidence:** `close-out-gates-run-phase5.json`, `traceable-commit-evidence-phase5.md`, `git-hygiene-receipt-2026-10-21-phase5.md`

### Completion signals

- **Machine close-out:** pass — gate `close_out` `allowed=true`; envelope `tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/evidence/request-evidence-envelope.v1.json` validate `fail_on_error_gaps=true` → `blocking_gap_count=0`, `advisory_gap_count=5`
- **Process contract:** pass — gate runner `--sync-dispositions --reconcile`; `python3 scripts/test_layer_b_ci_enforcement.py` OK; W0–W3 corpus unchanged (73/73 registry)
- **Adherence ledger:** pass with note — reconcile `process_grade` band **B** (`evidence_stale` on historical gate receipt paths); non-blocking under advisory policy

### Remaining risks

- Branch protection: require **L0 only** until Layer B job promoted (see Sponsor actions above)
- `spec-corpus-status.sh` still red until improvement queue YAML restored under `tied-bundle/docs/`
- 5 advisory envelope gaps (inquiry under advisory policy)

## Next session entry

1. Refresh improvement queue → green `spec-corpus-status.sh`.
2. Open PR and confirm Actions `layer-b-enforce` on a sidecar-touching diff.
3. Optional: pin `TIED_GIT_REF`, cache `mcp-server/node_modules`.
