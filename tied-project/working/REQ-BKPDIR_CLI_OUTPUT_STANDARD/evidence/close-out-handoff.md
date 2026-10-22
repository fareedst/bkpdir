# Close-out handoff — REQ-BKPDIR_CLI_OUTPUT_STANDARD

**Run ID:** `cli-output-close-20261022`  
**Feature commit SHA:** _(filled post-commit)_  
**Scope:** CLI stream contract + fleet/shell adoption; profile suite deferred.

### Completion signals

- **Machine close-out:** pass — `close-out-gates-run.json` → `merged_decision.allowed: true`; envelope `request-evidence-envelope.v1.json` blocking gaps **0** (`depth_tier: minimal`).
- **Process contract:** pass — Tracker dual-write through verification-gate; manifest `verification-evidence-manifest.v1.json`; `sub-adversarial-inquiry-pass` **not_applicable** with `not-applicable-receipt.v1.json`; gate runner `--sync-dispositions --reconcile`.
- **Adherence ledger:** pass — reconcile `process_grade` band **A** after evidence stubs + N/A receipt; ledger rows in `gates/ledger.jsonl`.

### Remaining risks

- Profile/ANSI TTY suite not in contract harness (`CLI-OPERATOR-OUTPUT-TEXTURE/spec.md`).
- `cmd/cli-template` separate module — no `clioutput` adoption yet.
- stderr-default help migration (`COSTLY-STDERR-DEFAULT-HELP`).
