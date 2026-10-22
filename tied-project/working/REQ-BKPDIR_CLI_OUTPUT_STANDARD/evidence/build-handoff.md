# Handoff — [REQ-BKPDIR_CLI_OUTPUT_STANDARD]

**Date:** 2026-10-22  
**Request:** CLI operator output standard (stream contract + replaceable profile)  
**Depth:** `minimal` · **Gate policy:** `advisory`

---

## Completion signals (honest status)

| Signal | Status | Notes |
| --- | --- | --- |
| **Machine — pre_implementation** | **Pass** | `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/gates/pre_implementation-2026-10-22T05-23-05-535Z.json` (`allowed: true`) |
| **Machine — verification** | **Not run** | Requires green full test scope + `tied_checklist_gate_validate` `phase: verification` + `tied_verify` |
| **Machine — close_out** | **Not run** | Requires envelope sync (`run-close-out-gates.mjs --envelope-blocking --sync-dispositions --reconcile`) |
| **Process contract** | **Partial** | TIED chain authored; pseudocode validated; Tracker dispositions in `checklist-dispositions.yaml` — unit-test-red+ steps not fully dispositioned |
| **Adherence ledger** | **Started** | `gates/ledger.jsonl` has pre_implementation receipt; no verification/close_out events |

**Do not claim “complete” or “ready to merge” until verification and close_out gates pass.**

---

## Delivered behavior

### Stream contract (primary binary + fleet)

- Default **`--help`**: empty **stdout**, help on **stderr** (`bkpdir`, Cobra tools using `clioutput.ExecuteRoot`).
- **`--stdout`** / **`TIED_CLI_HELP_STDOUT=1`**: help on stdout, success stderr empty.
- Execution errors: **stderr** (`PrintExecError`).
- **Peel** `--plain`, `-p`, `--no-color`, `--stdout` before Cobra / auto-detect argv.

### Code seams

| Area | Path |
| --- | --- |
| Shared module | `pkg/clioutput/` (`policy`, `peel`, `cobra`, `run`, `tty`) |
| Primary wire | `main.go` |
| Fleet | `cmd/token-suggester`, `cmd/ai-validation`, `cmd/realtime-validator` → `ExecuteRoot` |
| specctl | Unchanged manual `usage()` on stderr (no Cobra) |

### Tests (contract)

| Test file | Scope |
| --- | --- |
| `cli_output_contract_test.go` | `bkpdir` subprocess |
| `cli_output_fleet_contract_test.go` | token-suggester, ai-validation, realtime-validator, specctl usage error |
| `shell_cli_output_contract_test.go` | Seven migrated scripts + `TIED_CLI_HELP_STDOUT` |
| `main_cli_composition_test.go` | Subcommand help on stderr (in-process) |
| `pkg/clioutput/*_test.go` | Policy / peel unit |

Quick check:

```bash
go test -vet=off . -count=1 -run 'CLIOutput|ShellCLI|FleetCLI|Composition_ExecuteWithAutoDetection_KnownSubcommandHelp'
go test -vet=off ./pkg/clioutput/... -count=1
```

### Shell migration

- Include: `scripts/_cli_operator_output_standard.inc.sh`
- Migrated: `run-mutation-pilot.sh`, `run-pseudocode-layer-b-batch.sh`, `run-layer-b-ci-enforcement.sh`, `token-navigate.sh`, `validate-test-tokens.sh`, `generate-missing-tokens.sh`, `enhance-cross-reference-validation.sh`

### TIED + working artifacts

| Artifact | Path |
| --- | --- |
| REQ / ARCH / IMPL | `tied-project/requirements/REQ-BKPDIR_CLI_OUTPUT_STANDARD.yaml`, `architecture-decisions/ARCH-BKPDIR_CLI_OUTPUT_STANDARD.yaml`, `implementation-decisions/IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE.yaml` |
| Sidecar | `implementation-decisions/IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE-pseudocode.md` |
| CITDP | `tied-project/citdp/CITDP-REQ-BKPDIR_CLI_OUTPUT_STANDARD.yaml` + working copy |
| Portable next-client prompt | `tied-project/working/CLI-OPERATOR-OUTPUT-TEXTURE/spec.md` |
| Vocab | `tied-project/vocab/bkpdir-cli.md` (help / stream section) |

---

## Known gaps (next owner)

1. **Profile suite** — ANSI / decorated TTY snapshots from working spec (contract suite intentionally excludes).
2. **`cmd/cli-template`** — separate `go.mod`; no `clioutput` adoption or subprocess test yet.
3. **`cmd/scaffolding`** — no in-binary help (intentionally skipped).
4. **REQ/IMPL status** — still “In Progress”; run verification-gated `tied_verify` if project uses it.
5. **Full regression** — `go test ./...` (root may need `-vet=off` on legacy vet noise in unrelated files).
6. **Close-out** — CALL `sub-close-out-evidence-sync` or `tools/bootstrap/templates/run-close-out-gates.mjs` with required flags; validate envelope `fail_on_error_gaps: true`.

---

## Implementation lessons (for maintainers)

- **Cobra help:** delegate default `HelpFunc` with temporary `SetOut(helpWriter)` — do not re-parse help templates without Cobra’s full func map.
- **`-p` shorthand:** removed from **persistent** `--plain` — conflicts with subcommand shorthands (e.g. realtime-validator `server -p` / port). **`-p` still peeled from argv** before Execute on the main binary.
- **Dual path:** auto-detect and Cobra share peel + policy in `main`; extend peel if new routers added.
- **Subprocess tests** required where code writes directly to `os.Stdout`/`Stderr`.

---

## Costly choices (sponsor)

1. **stderr-default help** — `bkpdir --help | …` needs `2>&1` or explicit stdout help.
2. **Auto-detect `-n`** — still note vs list limit; texture must not repurpose `-n`.

---

## Suggested next session entry

1. `go test ./...` (or project Makefile target) and fix any regressions.
2. Update Tracker dispositions through `unit-test-green` / `verification-gate`.
3. `tied_checklist_gate_validate` `phase: verification` → `tied_verify` → close_out gate + envelope.
4. Optional: adopt `cli-template` via module `replace`; add profile tests behind one renderer seam.

**Prompt for next agent:** Link `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/evidence/build-handoff.md` and `/close-out` (or “finish verification gate for REQ-BKPDIR_CLI_OUTPUT_STANDARD”).
