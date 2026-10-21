# bkpdir CLI (canonical)

**Scope:** The **`bkpdir`** command-line interface: binary name, subcommands, global persistent flags, auto-detection routing, and backward-compatibility aliases. **Vocabulary only** — command dispatch algorithms live in [`../../main.go`](../../main.go) and IMPL detail files.

**Traceability:** [REQ-IMMUTABLE_CLI_COMMANDS](../requirements/REQ-IMMUTABLE_CLI_COMMANDS.yaml) · [REQ-IMMUTABLE_GLOBAL_OPTIONS](../requirements/REQ-IMMUTABLE_GLOBAL_OPTIONS.yaml) · [REQ-USABILITY](../requirements/REQ-USABILITY.yaml) · [REQ-LIST_LIMIT](../requirements/REQ-LIST_LIMIT.yaml) · [REQ-DIFF_COMMAND](../requirements/REQ-DIFF_COMMAND.yaml) · [REQ-DEBUG_OUTPUT_CONTROL](../requirements/REQ-DEBUG_OUTPUT_CONTROL.yaml) · [ARCH-CLI_COMMANDS](../architecture-decisions/ARCH-CLI_COMMANDS.yaml) · [ARCH-CLI_FRAMEWORK](../architecture-decisions/ARCH-CLI_FRAMEWORK.yaml) · [ARCH-AUTO_DETECTION](../architecture-decisions/ARCH-AUTO_DETECTION.yaml) · [ARCH-LIST_LIMIT](../architecture-decisions/ARCH-LIST_LIMIT.yaml) · [ARCH-DIFF_COMMAND](../architecture-decisions/ARCH-DIFF_COMMAND.yaml) · [IMPL-CLI_FRAMEWORK](../implementation-decisions/IMPL-CLI_FRAMEWORK.yaml) · [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) · [IMPL-LIST_LIMIT](../implementation-decisions/IMPL-LIST_LIMIT.yaml) · [IMPL-DIFF_COMMAND](../implementation-decisions/IMPL-DIFF_COMMAND.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-configuration.md`](bkpdir-configuration.md) · [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) · [`domain-references.md`](domain-references.md) · [`../../docs/user/specification.md`](../../docs/user/specification.md)

**Split/merge rule:** CLI surface terms stay in this file; YAML keys and archive naming live in sibling glossaries. Split if a new CLI binary is added; do not merge config or formatting terms here.

---

## Preferred terms vs synonyms

| Preferred | Avoid in UI/docs | Notes |
|-----------|------------------|-------|
| **bkpdir** | backup dir, bkp | Binary and Cobra root `Use` name |
| **auto-detection** | magic args, implicit mode | Positional path routes to file backup or directory archive |
| **subcommand** | command (alone) | Explicit verbs registered on root |
| **persistent flag** | global option (alone) | Flags on root available to all subcommands |
| **backward compatibility** | legacy alias | e.g. `full`/`inc`/`--config` on root |
| **dry-run** | preview, simulate | `-d` / `--dry-run`; no archive/backup created |
| **list limit** | pagination limit | `-n` / `--limit`; `0` = show all |

---

## Subcommands (catalog)

| Subcommand | Preferred use | Backward alias |
|------------|---------------|----------------|
| `create` | Explicit full or incremental archive | — |
| `config` | Inspect merged configuration | root `--config` |
| `template` | Emit default config YAML | — |
| `full` | Full directory archive | `create` without `--incremental` |
| `inc` | Incremental directory archive | `create --incremental` |
| `list` | List directory archives | — |
| `diff` | Compare directory to effective archive state | — |
| `backup` | Explicit single-file backup | auto-detect on file path |
| `version` | Print version string | `--version` / `-v` |

---

## Global persistent flags

| Flag | Short | Default | Notes |
|------|-------|---------|-------|
| `--dry-run` | `-d` | `false` | Show actions without writing archives/backups |
| `--config` | — | `false` | Display configuration and exit (root only) |
| `--list` | — | `""` | List backups for file path |
| `--limit` | `-n` | `10` | Max list items; `0` = all |
| `--debug` | — | `false` | Debug output ([REQ-DEBUG_OUTPUT_CONTROL](../requirements/REQ-DEBUG_OUTPUT_CONTROL.yaml)) |

---

## Subcommand flags (catalog)

| Subcommand | Flag | Short | Purpose |
|------------|------|-------|---------|
| `create` | `--incremental` | — | Incremental archive |
| `create`, `full`, `inc` | `--note` | — | Note suffix in archive name |
| `config` | `--all` | — | Show all fields (default) |
| `config` | `--overrides-only` | — | Non-default values only |
| `config` | `--sources` | — | Source attribution per value |
| `config` | `--format` | — | `table` \| `tree` \| `json` |
| `config` | `--filter` | — | Field name pattern filter |
| `config` | `--flat` | — | Flat list vs grouped output |
| `template` | `--output` | `-o` | Output path for template file |
| `template` | `--dry-run` | `-d` | Preview without writing |
| `template` | `--force` | `-f` | Overwrite existing template |
| `backup` | `--note` | — | Note suffix in backup name |

---

## Naming bridge: CLI vs code

| Concept | CLI surface | Go symbol | TIED token |
|---------|-------------|-----------|------------|
| Root command tree | `bkpdir` | `newRootCommand()` | [IMPL-CLI_FRAMEWORK](../implementation-decisions/IMPL-CLI_FRAMEWORK.yaml) |
| Auto-detect router | positional path | `handleAutoDetectedCommand` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| Cobra bypass for paths | `executeWithAutoDetection` | `executeWithAutoDetection` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| Known command guard | `create`, `config`, … | `knownCommands` slice | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| List limit | `--limit` / `-n` | `listLimit` | [IMPL-LIST_LIMIT](../implementation-decisions/IMPL-LIST_LIMIT.yaml) |
| Diff subcommand | `diff` | `diffCmd()` | [IMPL-DIFF_COMMAND](../implementation-decisions/IMPL-DIFF_COMMAND.yaml) |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Root command tree | `NEW_ROOT_COMMAND` | [IMPL-CLI_FRAMEWORK](../implementation-decisions/IMPL-CLI_FRAMEWORK.yaml) |
| Auto-detect routing | `HANDLE_AUTO_DETECTED_COMMAND` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| File path auto-detect | `HANDLE_AUTO_DETECTED_FILE_BACKUP` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| Directory auto-detect | `HANDLE_AUTO_DETECTED_DIRECTORY_ARCHIVE` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| Cobra execute bypass | `EXECUTE_WITH_AUTO_DETECTION` | [IMPL-AUTO_DETECTION](../implementation-decisions/IMPL-AUTO_DETECTION.yaml) |
| Diff command handler | `DIFF_CMD` | [IMPL-DIFF_COMMAND](../implementation-decisions/IMPL-DIFF_COMMAND.yaml) |

---

## Alphabetical index

| Term | Section |
|------|---------|
| auto-detection | Preferred terms |
| backward compatibility | Preferred terms |
| bkpdir | Preferred terms |
| config | Subcommands |
| create | Subcommands |
| diff | Subcommands |
| dry-run | Global persistent flags |
| full | Subcommands |
| inc | Subcommands |
| list | Subcommands |
| list limit | Preferred terms |
| persistent flag | Preferred terms |
| subcommand | Preferred terms |
| template | Subcommands |
| version | Subcommands |
| --config | Global persistent flags |
| --debug | Global persistent flags |
| --dry-run | Global persistent flags |
| --flat | Subcommand flags |
| --format | Subcommand flags |
| --incremental | Subcommand flags |
| --limit | Global persistent flags |
| --list | Global persistent flags |
| --note | Subcommand flags |
| --sources | Subcommand flags |
