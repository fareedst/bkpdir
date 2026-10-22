# BKPDIR-TWO-FOLDER-CLEANUP verification

- `tied_validate_consistency`: ok (2026-10-21 `@clear-workspace`, TIED_BASE_PATH=bkpdir/tied-project)
- `tied_project_lint`: ok (11 token_reference_issues advisory)
- `run-tied-new-client-audit.mjs --client-root`: ok (SC-TFL-NO-ROOT-WORKING: repo-root `working/` alongside `tied-project/working/`)
- `rg 'tied/'` excluding `tied-bundle/`, citdp, working: **0 matches** (2026-10-21; archive oracle test uses `nested/foo.yaml` fixture)
- `make test`: ok (2026-10-21)
- `make check`: not re-run this pass — run before traceable commit
- **Git:** ~130 untracked under `tied-project/`; 64 tracked dirty path/doc edits vs HEAD
