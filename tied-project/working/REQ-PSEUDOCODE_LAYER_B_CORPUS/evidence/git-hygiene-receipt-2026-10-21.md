# Git hygiene receipt — plan_close_out (2026-10-21)

**Tool:** `tied_git_hygiene` (`context_mode: plan_close_out`, `apply: false`)

## Untracked (in scope for W0 Layer B slice)

- `scripts/audit_corpus_layer_b_gaps.py`
- `scripts/run-pseudocode-layer-b-batch.sh`
- `scripts/test_audit_corpus_layer_b_gaps.py`
- `scripts/test_run_pseudocode_layer_b_batch.py`
- `scripts/testdata/layer_b_batch/`
- `tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/`

## Tracked dirty (staged for commit)

W0 sidecars, Layer B/C gate scripts, markscope note, block-lead sync in `archive.go` / `exclude_test.go`.

## Excluded from this commit

- `tied-project/working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/evidence/request-evidence-envelope.v1.json` (PFV slice; not Layer B corpus close-out)

## Gitignore N/A

No new ephemeral patterns required; `tied-bundle/working/` remains gitignored local inquiry root.
