# Continuation handoff — [REQ-INCREMENTAL_ANCHOR_MODE]

**Parent plans:** Cursor plan `full-anchored_incrementals_1ebb1978.plan.md` (full-anchored incrementals)  
**Branch:** `chore/tied-two-folder-cleanup`  
**Corpus / CI context:** [`phase5-ci-handoff.md`](../REQ-PSEUDOCODE_LAYER_B_CORPUS/evidence/phase5-ci-handoff.md) (Phase 5 committed `8c62f6767dfd`; modernization baseline in Cursor plan `impl_pseudocode_modernization_59816848.plan.md`)

## Sponsor intent (resolved)

Default **full-anchored incremental**: each `bkpdir inc` zip lists paths that differ from the latest **full** snapshot; skip/diff effective state = full + **latest** `_update=` only. **Chained** mode (`incremental_anchor_mode: chained`) preserves legacy merge-all-incrementals behavior. Same-minute collisions use optional `-seqN` before `.zip`.

## Implementation status (committed this slice)

| Plan todo | Status |
|-----------|--------|
| TIED REQ/ARCH/IMPL + pseudo-code sidecar | **Done** — indexes + detail YAML; `IMPL-INCREMENTAL_ANCHOR_MODE-pseudocode.md`; `tied-project/citdp/CITDP-INCREMENTAL-ANCHOR-MODE.yaml` |
| RED tests | **Green** — `TestIncrementalAnchorModeSecondIncCumulativeVsFull_*`, `TestIncrementalArchivePathDisambiguationSameMinute_*`; chained tests gated with `cfg.IncrementalAnchorMode = "chained"` |
| Config + CLI | **Done** — `incremental_anchor_mode` default `full`; `--incremental-anchor-mode` on `inc`/`diff` |
| `ReconstructArchiveState` / inc baseline | **Done** — `incremental_anchor_mode.go`, `ReconstructArchiveStateWithMode`, branches in `archive.go` / `comparison.go` |
| Path disambiguation | **Done** — `resolveIncrementalArchivePath` in `archive.go` |
| Vocab + user spec | **Done** — `bkpdir-archives-backup-diff.md`, `bkpdir-configuration.md`, `docs/user/specification.md` |
| LEAP duplicate-prevention / diff REQ+ARCH | **Done** — mode-scoped reconstruction text; `REQ-INCREMENTAL_ANCHOR_MODE` in depends_on |
| Tracker copy | **Done** — `agent-req-implementation-checklist.yaml` + minimal inquiry waiver |
| Close-out gates | **Done** — `close-out-gates-run-2026-10-21.json`, gate receipts under `gates/` |

### Re-verify (baseline)

```bash
go test -run 'TestIncremental|TestReconstruct|TestIncCommand|TestCalculateDiff|IncrementalAnchor|Disambiguation' .
TIED_BASE_PATH="$PWD/tied-project" .cursor/skills/tied-yaml/scripts/tied-cli.sh tied_validate_consistency
TIED_BASE_PATH="$PWD/tied-project" .cursor/skills/tied-yaml/scripts/tied-cli.sh pseudocode_validate IMPL-INCREMENTAL_ANCHOR_MODE --mode structural
```

**Last run this session:** incremental test subset **pass**; `tied_validate_consistency` **ok**; structural `pseudocode_validate` on `IMPL-INCREMENTAL_ANCHOR_MODE` **ok**.

## Working tree (this slice)

| Area | Notes |
|------|--------|
| Go | `incremental_anchor_mode.go`, `incremental_anchor_mode_test.go`, edits to `archive.go`, `comparison.go`, `config.go`, `main.go`, `diff_command_test.go`, `incremental_duplicate_prevention_test.go` |
| TIED | New REQ/ARCH/IMPL detail files; index updates; working CITDP under `tied-project/working/REQ-INCREMENTAL_ANCHOR_MODE/` |
| Scratch | `tmp_yaml_err/` — do not commit |

**Suggested commit:** one logical unit separate from corpus evidence commits (behavior + TIED stack + docs).

## Costly choices (recorded)

- **Default `full`** — users with legacy chained-only archive sets must set `incremental_anchor_mode: chained` or take a new full (documented in user spec).
- **`-seqN` suffix** — extension to immutable naming; LEAP note on `[REQ-IMMUTABLE_ARCHIVE_NAMING]` if detail text not yet updated.

## Parallel track — corpus tail (unchanged priority)

From Phase 5 handoff; not blocking incremental anchor close-out:

1. Restore `tied-bundle/docs/impl-pseudocode-improvement-queue.yaml` → green `spec-corpus-status.sh`.
2. PR: confirm `layer-b-enforce` on sidecar-touching diff; branch protection **L0 only** until Layer B promoted.
3. Optional: pin `TIED_GIT_REF`, cache `mcp-server/node_modules`.

## Next session entry

1. **LEAP** — amend `REQ-INCREMENTAL_DUPLICATE_PREVENTION` (+ ARCH/IMPL if needed) for mode-scoped reconstruction; add `related_requirements` link to `REQ-INCREMENTAL_ANCHOR_MODE`.
2. **Tracker + gates** — bootstrap checklist in working folder; `pre_implementation` / `verification` `tied_checklist_gate_validate` per CITDP.
3. **`lint_yaml`** on changed project TIED files (if any direct edits remain).
4. **Commit** incremental anchor slice; manual replay of sponsor log (full → inc → inc dry-run) if not already captured in tests.
5. **Corpus** — improvement queue refresh when sponsor wants `spec-corpus-status.sh` green (independent PR).
