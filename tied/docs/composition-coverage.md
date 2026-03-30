# Composition test coverage (bkpdir)

This inventory maps **bindings** between tested units to **tests** that exercise those bindings without duplicating production wiring. Process reference: `[PROC-AGENT_REQ_CHECKLIST]` S10, `tied/processes.md` (composition tests for bindings).

## CLI entry and Cobra tree

| Binding | Description | Coverage |
|--------|-------------|----------|
| Production root command | `newRootCommand()` in `main.go` registers the same flags, root `Run`, and subcommands as the binary. | `createTestRootCmd()` delegates to `newRootCommand()` (`main_test.go`). |
| Dispatch before Cobra | `executeWithAutoDetection(root, args)` routes known subcommands/flags to `rootCmd.SetArgs` + `Execute`, else path auto-detect. | `main_cli_composition_test.go` (`TestComposition_*`). |
| Root persistent flags | `--config`, `--list`, `--limit` on the root command. | `main_cli_composition_test.go`, `config_integration_test.go`. |
| `diff` subcommand registration | `diff` is on the production root (not added ad hoc in tests). | `TestComposition_DiffViaExecuteWithAutoDetection_REQ_DIFF_COMMAND`; `inc_diff_integration_test.go` uses `createTestRootCmd()` only. |

## Config and handlers

| Binding | Description | Coverage |
|--------|-------------|----------|
| Config load → output | `handleConfigCommand` / `LoadConfig` with real filesystem. | `TestComposition_RootPersistentConfigFlag_REQ_CONFIGURATION`, `config_integration_test.go`. |
| List backups + limit | `handleListFileBackupsCommand` + `listLimit`. | `TestComposition_RootListAndLimit_REQ_LIST_LIMIT`, `backup_test.go`. |

## `pkg/fileops` (module `bkpdir/pkg/fileops`)

| Binding | Description | Coverage |
|--------|-------------|----------|
| Snapshots and compare | `CreateDirectorySnapshot`, `CreateArchiveSnapshot`, `CompareSnapshots`, `IsDirectoryIdenticalToArchive` used by `comparison.go` in `main`. | `pkg/fileops/comparison_test.go`. |
| Exclusions | `ShouldExcludeFile` / `PatternMatcher` used when building snapshots. | `pkg/fileops/exclusion_test.go`. |

## Residual risks

- Handlers under `main` that call `os.Exit` cannot be fully exercised in-process when they take the exit path; composition tests use **dry-run**, **existing archives**, or **help** to stay on non-exit paths where possible.
- `make test` runs `go test ./...` plus `cd pkg/fileops && go test ./...` so the fileops module is covered.

## Related tokens

- `[IMPL-AUTO_DETECTION]`, `[ARCH-AUTO_DETECTION]`, `[REQ-USABILITY]` — CLI dispatch and auto-detect.
- `[REQ-DIFF_COMMAND]`, `[IMPL-DIRECTORY_COMPARISON]` — diff and fileops comparison.
