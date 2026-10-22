# Test strategy — appearance profile slice

**Parent REQ:** `[REQ-BKPDIR_CLI_OUTPUT_STANDARD]`  
**Slice:** `appearance-profile-2026-10`  
**Profile reference:** [`../CLI-OPERATOR-OUTPUT-TEXTURE/spec.md`](../CLI-OPERATOR-OUTPUT-TEXTURE/spec.md) §4  
**Stream contract strategy (unchanged):** [`test-strategy.md`](test-strategy.md)

## Suite split (mandatory)

| Suite | Package / files | Proves | Must not assert |
| --- | --- | --- | --- |
| **Profile** | `pkg/clioutput/help_render_test.go`, `pkg/clioutput/cobra_profile_test.go` (names flexible) | ANSI/layout when plain off + TTY; plain ASCII when plain on or non-TTY | Help destination, stdout emptiness, exit codes |
| **Contract** | `cli_output_contract_test.go`, `cli_output_fleet_contract_test.go`, `shell_cli_output_contract_test.go`, `main_cli_composition_test.go` | Streams only | ESC, `▶`, rule width, exact help copy |

**Rule:** After any profile-only edit, run **profile suite first**, then **contract suite unchanged**.

## Fixtures

- Store fixed Cobra help text under `pkg/clioutput/testdata/help_*.txt` (synthetic minimal trees: root, subcommand, empty sections).
- Tests call `RenderHelp(fixture, opts)` — do not depend on live `bkpdir` version strings or full command catalog prose.

## TDD order ([PROC-TIED_DEV_CYCLE])

1. **RED — `RENDER_HELP` unit:** decorated vs plain; edge inputs; no panics.
2. **GREEN — renderer implementation.**
3. **RED — Cobra wiring:** in-process buffers + injected TTY; subcommand help; no ESC on plain/`NO_COLOR`/non-TTY.
4. **GREEN — `installHelpFunc` buffer + `RenderHelp`; fix `CaptureHelp` stacking.
5. **RED — fleet smoke (optional table):** one decorated-or-plain check per `ExecuteRoot` binary via composition helper (not full subprocess ANSI unless already harnessed).
6. **Shell — plain pipe test:** script sourcing include emits no ESC when stdout/stderr not TTY.
7. **Shell — decorated (optional):** PTY-backed on `unix` only; skip elsewhere with documented reason.

## Assertion patterns

| Check | Decorated TTY | Plain / non-TTY |
| --- | --- | --- |
| ESC present | Allow `\033[` or `\x1b[` in section/header lines | `!strings.Contains(out, "\033")` |
| Rule character | Unicode `─` dim (decorated) | 56× `-` |
| Section marker | Optional `▶` | Absent |
| Command row | Bold name + padded description | Fixed columns, no bold codes |

## Local / CI commands

```bash
# Profile-focused
go test -vet=off ./pkg/clioutput/... -count=1 -run 'HelpRender|Profile|Renderer'

# Contract regression (unchanged selectors)
go test -vet=off . -count=1 -run 'CLIOutput|ShellCLI|FleetCLI|Composition_ExecuteWithAutoDetection_KnownSubcommandHelp'

# Full tree before verification gate
go test -vet=off ./...
```

## Pre-implementation gate

Do not write RED tests until:

- IMPL sidecar includes `RENDER_HELP` and updated `APPLY_TO_ROOT`.
- `pseudocode_validate` passes on the sidecar.
- `tied_checklist_gate_validate` `phase: pre_implementation` returns `allowed: true` for the appearance-profile Tracker slice.
