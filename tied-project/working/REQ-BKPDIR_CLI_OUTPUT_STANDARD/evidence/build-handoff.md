# Handoff — [REQ-BKPDIR_CLI_OUTPUT_STANDARD]

**Date:** 2026-10-22  
**Request:** CLI operator output standard (stream contract + appearance profile)  
**Depth:** `minimal` · **Gate policy:** `advisory`  
**Close-out run:** `cli-output-profile-close-20261022`

---

## Completion signals (honest status)

| Signal | Status | Notes |
| --- | --- | --- |
| **Machine — pre_implementation** | **Pass** | Stream receipt `gates/pre_implementation-2026-10-22T05-23-05-535Z.json`; profile logic covered by post-stream pseudocode + verification re-run |
| **Machine — verification** | **Pass** | `verification-evidence-manifest.v1.json` run `cli-output-profile-ver-20261022`; gate runner verification phase `allowed: true` |
| **Machine — close_out** | **Pass** | `evidence/close-out-gates-run-profile.json` → `merged_decision.allowed: true`; envelope blocking gaps **0** |
| **Process contract** | **Pass** | Tracker through `verification-gate`; profile test logs on disk |
| **Adherence ledger** | **Pass** | Reconcile band **A** in profile close-out run |

---

## Delivered behavior

### Stream contract (unchanged)

- Default **`--help`**: empty **stdout**, help on **stderr**.
- **`--stdout`** / **`TIED_CLI_HELP_STDOUT=1`**: help on stdout, success stderr empty.
- Execution errors: **stderr** (`PrintExecError`).
- **Peel** global texture flags before Cobra / auto-detect argv.

### Appearance profile (slice `appearance-profile-2026-10`)

- **`RenderHelp`** in `pkg/clioutput/help_render.go` — post-processes Cobra help (decorated TTY vs plain).
- **Cobra wiring** — buffer default `HelpFunc`, render, write to policy help writer.
- **Shell** — `_cli_out_profile_*` helpers in `scripts/_cli_operator_output_standard.inc.sh`.
- **Tests** — profile unit/composition + `shell_cli_output_profile_test.go` (pipe/plain); contract suite unchanged.

### Code seams

| Area | Path |
| --- | --- |
| Shared module | `pkg/clioutput/` (`policy`, `peel`, `cobra`, `help_render`, `run`, `tty`) |
| Primary wire | `main.go` |
| Fleet | `cmd/token-suggester`, `cmd/ai-validation`, `cmd/realtime-validator` → `ExecuteRoot` |

Quick check:

```bash
go test -vet=off ./pkg/clioutput/... -count=1 -run 'RenderHelp|Profile|Policy|Peel'
go test -vet=off . -count=1 -run 'CLIOutput|ShellCLI|FleetCLI|ShellCLIOutputProfile|Composition_ExecuteWithAutoDetection'
go test -vet=off ./... -count=1
```

---

## Known gaps (next owner)

1. **`cmd/cli-template`** — separate `go.mod`; no `clioutput` adoption yet.
2. **`cmd/scaffolding`** — no in-binary help (intentionally skipped).
3. **Shell PTY decorated test** — optional; Go profile suite is primary proof.
4. **Script bodies** — may adopt `_cli_out_profile_*` incrementally (helpers available in include).

---

## Costly choices (sponsor)

1. **stderr-default help** — `bkpdir --help | …` needs `2>&1` or explicit stdout help.
2. **Auto-detect `-n`** — still note vs list limit; texture must not repurpose `-n`.
