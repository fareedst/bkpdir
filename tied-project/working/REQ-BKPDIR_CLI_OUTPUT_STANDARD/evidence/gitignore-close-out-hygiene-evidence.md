# Gitignore close-out hygiene — REQ-BKPDIR_CLI_OUTPUT_STANDARD

**Date:** 2026-10-22

Reviewed untracked paths from implementation pass:

- `tmp_yaml_err/` — scratch YAML lint; **do not commit**; no `.gitignore` change (local only).
- Contract test artifacts under `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/evidence/` — **commit** as close-out evidence.

No new ephemeral patterns required in root `.gitignore` for this REQ.
