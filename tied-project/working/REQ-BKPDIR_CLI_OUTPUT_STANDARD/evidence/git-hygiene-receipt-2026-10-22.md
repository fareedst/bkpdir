# Git hygiene receipt — REQ-BKPDIR_CLI_OUTPUT_STANDARD

**Date:** 2026-10-22  
**Tool:** `tied_git_hygiene` · `context_mode: plan_close_out` · `apply: false`

## Preview summary

- **Untracked (REQ-scoped commit):** `pkg/clioutput/`, CLI contract tests, shell include, TIED R+A+I + working folder for this REQ.
- **Excluded from REQ commit:** `tmp_yaml_err/` (local scratch), unrelated `REQ-PSEUDOCODE_LAYER_B_CORPUS` drift.
- **Gitignore:** No new patterns required (`gitignore-close-out-hygiene.md`).
