# bkpdir archives, backup, and diff (canonical)

**Scope:** Directory archives (full and incremental ZIP), single-file backups, archive/backup naming grammar, directory snapshots, diff results, identical-skip behavior, and duplicate prevention. **Vocabulary only** — ZIP and comparison algorithms live in [`../../archive.go`](../../archive.go), [`../../backup.go`](../../backup.go), [`../../comparison.go`](../../comparison.go), and [`../../pkg/fileops/`](../../pkg/fileops/).

**Traceability:** [REQ-FILE_BACKUP](../requirements/REQ-FILE_BACKUP.yaml) · [REQ-INCREMENTAL_DUPLICATE_PREVENTION](../requirements/REQ-INCREMENTAL_DUPLICATE_PREVENTION.yaml) · [REQ-DIFF_COMMAND](../requirements/REQ-DIFF_COMMAND.yaml) · [REQ-IMMUTABLE_ARCHIVE_NAMING](../requirements/REQ-IMMUTABLE_ARCHIVE_NAMING.yaml) · [REQ-IMMUTABLE_FILE_BACKUP_NAMING](../requirements/REQ-IMMUTABLE_FILE_BACKUP_NAMING.yaml) · [REQ-IMMUTABLE_FILE_BACKUP_OPERATIONS](../requirements/REQ-IMMUTABLE_FILE_BACKUP_OPERATIONS.yaml) · [REQ-IMMUTABLE_DIRECTORY_OPERATIONS](../requirements/REQ-IMMUTABLE_DIRECTORY_OPERATIONS.yaml) · [REQ-IMMUTABLE_FILE_EXCLUSION](../requirements/REQ-IMMUTABLE_FILE_EXCLUSION.yaml) · [ARCH-ARCHIVE_FORMAT](../architecture-decisions/ARCH-ARCHIVE_FORMAT.yaml) · [ARCH-DIRECTORY_COMPARISON](../architecture-decisions/ARCH-DIRECTORY_COMPARISON.yaml) · [ARCH-INCREMENTAL_DUPLICATE_PREVENTION](../architecture-decisions/ARCH-INCREMENTAL_DUPLICATE_PREVENTION.yaml) · [ARCH-DIFF_COMMAND](../architecture-decisions/ARCH-DIFF_COMMAND.yaml) · [IMPL-ZIP_FORMAT](../implementation-decisions/IMPL-ZIP_FORMAT.yaml) · [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) · [IMPL-INCREMENTAL_DUPLICATE_PREVENTION](../implementation-decisions/IMPL-INCREMENTAL_DUPLICATE_PREVENTION.yaml) · [IMPL-DIFF_COMMAND](../implementation-decisions/IMPL-DIFF_COMMAND.yaml) · [IMPL-DATA_MODELS](../implementation-decisions/IMPL-DATA_MODELS.yaml) · [IMPL-ATOMIC_OPS](../implementation-decisions/IMPL-ATOMIC_OPS.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-cli.md`](bkpdir-cli.md) · [`bkpdir-configuration.md`](bkpdir-configuration.md) · [`bkpdir-git-integration.md`](bkpdir-git-integration.md) · [`../../docs/user/specification.md`](../../docs/user/specification.md)

**Split/merge rule:** Archive/backup/diff vocabulary stays here; CLI flags in [`bkpdir-cli.md`](bkpdir-cli.md); exclusion matching rules in [`bkpdir-configuration.md`](bkpdir-configuration.md).

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **full archive** | full backup (for dirs) | ZIP of entire directory tree |
| **incremental archive** | delta backup | ZIP of changed files since base full archive |
| **file backup** | file archive | Single-file timestamped copy (not ZIP) |
| **base archive** | parent archive | Most recent full archive for incremental chain |
| **effective archive state** | merged archive | Full + latest incremental combined view |
| **identical directory** | no changes | Skip create when dir matches archive |
| **identical file** | unchanged file | Skip backup when content matches |
| **directory snapshot** | dir scan | Path → hash map for comparison |
| **diff result** | change list | Added / modified / deleted file sets |

---

## Archive naming grammar

Full archive pattern:

```text
[PREFIX-]YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip
```

Incremental archive pattern:

```text
BASENAME_update=YYYY-MM-DD-hh-mm[=BRANCH=HASH][=NOTE].zip
```

| Token | Meaning |
|-------|---------|
| `PREFIX-` | Optional directory name prefix when `use_current_dir_name` |
| `YYYY-MM-DD-hh-mm` | Timestamp segment |
| `=BRANCH=HASH` | Git branch and short hash when Git info enabled |
| `=NOTE` | User note from `--note` or positional arg |
| `_update=` | Incremental marker substring |
| `-dirty` | Appended when repo dirty and dirty status enabled |

File backup pattern:

```text
SOURCE_FILENAME-YYYY-MM-DD-hh-mm[=NOTE]
```

---

## Operations catalog

| Operation | CLI entry | Go concept |
|-----------|-----------|------------|
| Full directory archive | `bkpdir .`, `create`, `full` | `Archive`, `ArchiveConfig` |
| Incremental archive | `bkpdir inc`, `create --incremental` | `IncrementalArchiveConfig` |
| File backup | `bkpdir backup`, auto-detect file | `Backup`, `BackupOptions` |
| List archives | `list` | archive directory scan |
| List file backups | `--list FILE` | `ListFileBackupsEnhanced` |
| Diff | `diff` | `CalculateDiff`, `DiffResult` |

---

## Snapshot and diff types

| Type | Package | Role |
|------|---------|------|
| `DirectorySnapshot` | `pkg/fileops` | Files map for comparison |
| `FileInfo` | `pkg/fileops` | Path, size, mod time, content hash |
| `DiffResult` | `main` | `Added`, `Modified`, `Deleted` slices |
| `ArchiveConfig` | `main` | Full archive naming inputs |
| `IncrementalArchiveConfig` | `main` | Incremental naming + base reference |
| `BackupInfo` | `main` | File backup metadata |

Key functions: `CreateDirectorySnapshot`, `CreateArchiveSnapshot`, `CompareSnapshots`, `IsDirectoryIdenticalToArchive`, `ReconstructArchiveState`, `FindMostRecentArchive`, `CalculateDiff`.

---

## Naming bridge

| Concept | User label | Config key | Code symbol |
|---------|------------|------------|-------------|
| Archive storage | archive directory | `archive_dir_path` | `ArchiveDirPath` |
| File backup storage | backup directory | `backup_dir_path` | `BackupDirPath` |
| Incremental detector | `_update=` in filename | — | incremental substring match |
| Duplicate prevention | skip identical | — | [REQ-INCREMENTAL_DUPLICATE_PREVENTION](../requirements/REQ-INCREMENTAL_DUPLICATE_PREVENTION.yaml) |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Create directory snapshot | `CREATE_DIRECTORY_SNAPSHOT` | [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) |
| Create archive snapshot | `CREATE_ARCHIVE_SNAPSHOT` | [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) |
| Compare snapshots | `COMPARE_SNAPSHOTS` | [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) |
| Identical directory check | `IS_DIRECTORY_IDENTICAL_TO_ARCHIVE` | [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) |
| Reconstruct archive state | `RECONSTRUCT_ARCHIVE_STATE` | [IMPL-INCREMENTAL_DUPLICATE_PREVENTION](../implementation-decisions/IMPL-INCREMENTAL_DUPLICATE_PREVENTION.yaml) |
| Calculate diff | `CALCULATE_DIFF` | [IMPL-DIFF_COMMAND](../implementation-decisions/IMPL-DIFF_COMMAND.yaml) |
| Find most recent archive | `FIND_MOST_RECENT_ARCHIVE` | [IMPL-DIRECTORY_COMPARISON](../implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml) |
| ZIP archive format | `ZIP_FORMAT` | [IMPL-ZIP_FORMAT](../implementation-decisions/IMPL-ZIP_FORMAT.yaml) |

---

## Alphabetical index

| Term | Section |
|------|---------|
| base archive | Preferred terms |
| CalculateDiff | Snapshot and diff types |
| diff result | Preferred terms |
| DirectorySnapshot | Snapshot and diff types |
| effective archive state | Preferred terms |
| file backup | Preferred terms |
| full archive | Preferred terms |
| identical directory | Preferred terms |
| identical file | Preferred terms |
| incremental archive | Preferred terms |
| _update= | Archive naming grammar |
| =BRANCH=HASH | Archive naming grammar |
| =NOTE | Archive naming grammar |
| ReconstructArchiveState | Snapshot and diff types |
