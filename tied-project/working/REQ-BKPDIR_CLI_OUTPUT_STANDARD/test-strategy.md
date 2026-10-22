# Test strategy — [REQ-BKPDIR_CLI_OUTPUT_STANDARD]

**Tracker:** `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/agent-req-implementation-checklist.yaml`  
**CITDP:** `tied-project/working/REQ-BKPDIR_CLI_OUTPUT_STANDARD/CITDP-REQ-BKPDIR_CLI_OUTPUT_STANDARD.yaml`  
**Profile reference:** `tied-project/working/CLI-OPERATOR-OUTPUT-TEXTURE/spec.md`

## Suite split

| Suite | Proves | Must not assert |
| --- | --- | --- |
| **Contract** | Stream routing, exit codes, help destination selection, stdin non-consumption, path parity | ANSI colors, symbols, rule width, exact help prose |
| **Profile (current)** | Working-spec decoration and control mapping | Machine output line formats from `format_*` templates |

## TDD order ([PROC-TIED_DEV_CYCLE])

1. **Unit RED — policy** (`pkg/clioutput/policy_*_test.go`): help writer resolution, plain triggers, TTY on selected writer, env/flag precedence.
2. **Unit RED — conformance harness** (`pkg/clioutput/conformance_*_test.go`): table cases shared by Go subprocess and shell tests.
3. **Unit RED — profile/renderer** (`pkg/clioutput/profile_*_test.go`, `renderer_*_test.go`): current appearance only.
4. **Composition RED — bkpdir** (`main_cli_composition_test.go`): separate stdout/stderr buffers; Cobra and auto-detect entry points.
5. **Subprocess contract — fleet**: one lightweight test file per adopted binary (see plan fleet table).
6. **Shell contract**: script invoking include + matrix cases vs Go harness expectations.

## Adopted surfaces (contract coverage)

| Surface | Class | Contract test |
| --- | --- | --- |
| `bkpdir` (root module) | Cobra + auto-detect | composition + subprocess |
| `token-suggester` | Cobra | subprocess |
| `cli-template` | Cobra | subprocess |
| `ai-validation` | Cobra | subprocess |
| `realtime-validator` | Cobra | subprocess |
| `specctl` | Manual usage (already stderr usage) | subprocess (usage errors; matrix stdout unchanged) |
| Seven shell scripts in plan | `usage()` / `--help` | shell harness |

**Out of contract scope:** `cmd/scaffolding` (no in-binary help; Makefile `help` only).

## CI / local gates

- `go test ./pkg/clioutput/...` and affected packages
- Full `go test ./...` before verification gate
- `pseudocode_validate` on [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] sidecar
- `tied_validate_consistency` after TIED records land
