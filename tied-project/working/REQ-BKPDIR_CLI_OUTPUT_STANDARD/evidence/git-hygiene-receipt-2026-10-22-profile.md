# Git hygiene receipt — REQ-BKPDIR_CLI_OUTPUT_STANDARD (profile close-out)

**Date:** 2026-10-22  
**Tool:** `tied_git_hygiene` · `context_mode: plan_close_out` · `apply: false`

## Preview summary

- Untracked REQ-scoped paths staged in traceable commit (profile code, tests, working slice artifacts).
- Excluded from commit: `tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-ci/` (unrelated REQ drift).
- No destructive apply; no gitignore pattern changes required for `tmp_yaml_err/` (not present in preview).

## N/A

- Secret/credential paths: none in staging map.
