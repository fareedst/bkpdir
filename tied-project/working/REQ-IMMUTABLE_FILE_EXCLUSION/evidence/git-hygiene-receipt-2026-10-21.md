# Git hygiene receipt — plan_close_out (preview)

- **Date:** 2026-10-21
- **Mode:** preview (`apply: false`)
- **Context:** `@clear-workspace` — REQ-IMMUTABLE_FILE_EXCLUSION
- **MCP:** `tied_git_hygiene` (`context_mode: plan_close_out`, `apply: false`)
- **Untracked (count):** 130 paths (full list in MCP preview; includes entire `tied-project/` tree migration artifacts and legacy `working/REQ-IMMUTABLE_FILE_EXCLUSION/` duplicate envelope)
- **Tracked dirty (count):** 58 paths (TIED restructure deletions at repo root, docs, scripts, specctl, tests — not scoped to exclude wave)
- **Apply:** none — no sponsor `allowed_paths` allowlist
- **Scoped wave:** Staging map paths in `proposed-commit-message.txt` are **clean vs HEAD** (no pending diff on exclude implementation files)
- **Note:** Repo-wide two-folder / bootstrap drift remains; do not bulk-delete untracked `tied-project/` without a named migration commit plan (`BKPDIR-TWO-FOLDER-CLEANUP`).
