# Close-out blockers — BKPDIR-TWO-FOLDER-CLEANUP

**Status:** incomplete — not commit-ready (`traceable-commit` checklist step still **pending**)

## Gate / envelope

- **`request_evidence_envelope_build`:** failed — `InvalidRequestToken` (working folder uses `change_request_id`, not `REQ-*`; envelope runner expects REQ token layout)
- **`tied_checklist_gate_validate` (`close_out`):** not green — needs full CITDP hydration + envelope artifacts (`hydration_envelope_gaps_missing`)
- **`@close-out`:** not run end-to-end — run after path hygiene and envelope wiring (or adopt REQ-prefixed working alias per methodology)

## CITDP success criteria (current)

| Criterion | State |
|-----------|--------|
| `tied_validate_consistency` | **ok** (2026-10-21 clear-workspace) |
| `tied_project_lint` | **ok** (11 token_reference_issues advisory) |
| `run-tied-new-client-audit.mjs` | **ok** (root `working/` layout warning) |
| `rg 'tied/'` outside `tied-bundle/` | **ok** (citdp/working/CHANGELOG excluded from hygiene scan) |
| `make test` / `make check` | **test ok** this run; re-run `make check` before commit |

### Remaining `tied/` references (excluding citdp/working/CHANGELOG)

- **None** in product/docs/tests (2026-10-21) — `docs/user/specification.md` path rewrite; archive conformance fixture uses `nested/foo.yaml`

## Git / staging

- **`tied-project/`:** only **16** paths tracked; bulk migration YAML/detail **untracked** — migration commit must `git add tied-project/` (indexes, citdp, working history) per CITDP desired behavior
- **Tracked edits:** 64 files vs HEAD (path string updates, `.gitignore` install block, spec tooling)
- **No staged changes** — traceable commit not prepared

## Recommended next steps

1. Rewrite `docs/user/specification.md` (and any other product docs) to `tied-project/` / `tied-bundle/` paths; reconcile archive test exception
2. Run **`@close-out`** with `--request-token` strategy resolved (envelope + gate receipts under `evidence/`)
3. Stage scoped tree: untracked `tied-project/**` + tracked path/doc edits; exclude ephemeral `working/REQ-IMMUTABLE_FILE_EXCLUSION/` duplicate if consolidating under `tied-project/working/`
4. Explicit sponsor ask for **traceable commit** after gates green
