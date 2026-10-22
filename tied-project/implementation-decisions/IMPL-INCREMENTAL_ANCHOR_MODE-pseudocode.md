# IMPL-INCREMENTAL_ANCHOR_MODE essence pseudocode

## PARSE_INCREMENTAL_ANCHOR_MODE

- [IMPL-INCREMENTAL_ANCHOR_MODE] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] — How: validate config/CLI values full or chained; default full when empty.

```
INPUT: raw string from config or CLI override
OUTPUT: IncrementalAnchorMode
PRE: raw is trimmed string or empty
POST: returns full or chained, or error if unknown
```

## RECONSTRUCT_ARCHIVE_STATE_WITH_MODE

- [IMPL-INCREMENTAL_ANCHOR_MODE] [IMPL-DIFF_COMMAND] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] [REQ-DIFF_COMMAND] — How: load latest full snapshot; apply incrementals per mode (full: latest only; chained: all in name order).

```
INPUT: archiveDir, mode
OUTPUT: DirectorySnapshot effective archived state
```

## CREATE_INCREMENTAL_ARCHIVE

- [IMPL-INCREMENTAL_ANCHOR_MODE] [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] — How: skip when diff vs effective state empty; zip file list from diff vs full snapshot in full mode, else vs chained effective state.

```
INPUT: IncrementalArchiveConfig including CLIAnchorMode
OUTPUT: inc zip or skip message
EFFECTS: resolveIncrementalArchivePath before create/dry-run
```

## RESOLVE_INCREMENTAL_ARCHIVE_PATH

- [IMPL-INCREMENTAL_ANCHOR_MODE] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] [REQ-IMMUTABLE_ARCHIVE_NAMING] — How: if candidate exists, try base-seqN.zip for N from 2 upward until free path.

```
INPUT: candidate incremental path
OUTPUT: unique path on disk
```
