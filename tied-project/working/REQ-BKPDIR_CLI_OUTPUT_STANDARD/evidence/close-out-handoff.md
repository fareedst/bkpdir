# Close-out handoff — REQ-BKPDIR_CLI_OUTPUT_STANDARD

**Run ID:** `cli-output-profile-close-20261022`  
**Feature commit SHA:** *(filled post-commit)*  
**Scope:** Stream contract + appearance profile (`RenderHelp`, profile tests, shell profile helpers).

### Completion signals

- **Machine close-out:** pass — `close-out-gates-run-profile.json` → `merged_decision.allowed: true`; envelope `request-evidence-envelope.v1.json` blocking gaps **0** (`depth_tier: minimal`).
- **Process contract:** pass — Tracker through `verification-gate`; manifest `verification-evidence-manifest.v1.json` run `cli-output-profile-ver-20261022`; gate runner `--sync-dispositions --reconcile`.
- **Adherence ledger:** pass — reconcile `process_grade` band **A** in profile close-out run.

### Remaining risks

- `cmd/cli-template` separate module — no `clioutput` adoption yet.
- Shell decorated proof Go-primary (PTY optional).
- stderr-default help migration (`COSTLY-STDERR-DEFAULT-HELP`).
