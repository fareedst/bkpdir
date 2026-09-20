# [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-CONTEXT_SUPPORT] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING]

## Summary contract

CLI diff command reconstructs effective archive state from latest full plus incremental zip snapshots, compares cwd snapshot to that state, and prints configurable added/modified/deleted reports.

INPUT: archiveDir, cwd, exclude patterns, context
OUTPUT: DiffResult, stdout diff text
DATA: DirectorySnapshot maps by RelativePath, name-sorted Archive list

## FIND_LATEST_FULL_ARCHIVE

SPEC-ID: IMPL-DIFF_COMMAND::FIND_LATEST_FULL_ARCHIVE
STEP T001: FILTER non-incremental archives AND return latest by name

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: list archives, filter non-incremental, sort by name ascending, return latest full Archive.

PROCEDURE FIND_LATEST_FULL_ARCHIVE(archiveDir):
  fullArchives = FILTER ListArchives WHERE NOT IsIncremental
  SORT BY Name; RETURN last entry OR error if empty

## FIND_LATEST_INCREMENTAL_ARCHIVE

SPEC-ID: IMPL-DIFF_COMMAND::FIND_LATEST_INCREMENTAL_ARCHIVE
STEP T001: FILTER incrementals by base prefix AND return latest by name

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: select incrementals whose name prefix matches base full archive _update= pattern; return latest by name.

PROCEDURE FIND_LATEST_INCREMENTAL_ARCHIVE(archiveDir, baseFull):
  matching = FILTER incrementals WITH prefix baseName + "_update="
  SORT BY Name; RETURN last OR nil

## RECONSTRUCT_ARCHIVE_STATE

SPEC-ID: IMPL-DIFF_COMMAND::RECONSTRUCT_ARCHIVE_STATE
STEP T001: LOAD full zip snapshot AND overlay incremental entries by path

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: load full zip snapshot, overlay incremental zip files by RelativePath, return merged DirectorySnapshot.

PROCEDURE RECONSTRUCT_ARCHIVE_STATE(archiveDir):
  fullSnapshot = CreateArchiveSnapshot(latestFull.Path)
  IF no incremental: RETURN fullSnapshot
  incrementalSnapshot = CreateArchiveSnapshot(latestIncremental.Path)
  MERGE maps: incremental entries override full by RelativePath
  RETURN sorted file list snapshot

## CALCULATE_DIFF

SPEC-ID: IMPL-DIFF_COMMAND::CALCULATE_DIFF
STEP T001: SNAPSHOT cwd AND classify added, modified, deleted file paths

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: snapshot cwd, compare file maps by path/size/hash; classify added, modified, deleted (files only).

PROCEDURE CALCULATE_DIFF(cwd, reconstructed, excludePatterns):
  current = CreateDirectorySnapshot(cwd, excludePatterns)
  BUILD currentMap and reconstructedMap (non-dir files only)
  added = paths in current not in reconstructed
  modified = same path with size or hash mismatch
  deleted = paths in reconstructed not in current
  RETURN DiffResult

## DIFF_CMD

SPEC-ID: IMPL-DIFF_COMMAND::DIFF_CMD
STEP T001: LOAD config, reconstruct state, calculate diff, AND print result

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-CONTEXT_SUPPORT] [REQ-DIFF_COMMAND] — How: load config, reconstruct state, handle no-archive gracefully, CalculateDiff, PrintDiffResult with context cancellation checks.

PROCEDURE DIFF_CMD():
  cfg = LoadConfig; archiveDir = prepareArchiveDirectory
  state = ReconstructArchiveState OR PrintNoArchivesFound on empty
  diff = CalculateDiff; formatter.PrintDiffResult(diff)

## FORMAT_DIFF_RESULT

SPEC-ID: IMPL-DIFF_COMMAND::FORMAT_DIFF_RESULT
STEP T001: EMIT no-changes header OR per-category added/modified/deleted lines

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-OUTPUT_FORMATTING] [REQ-DIFF_COMMAND] — How: use cfg FormatDiff* strings for no-changes header and per-category file lines.

PROCEDURE FORMAT_DIFF_RESULT(diff):
  IF empty diff: RETURN FormatDiffNoChanges
  EMIT FormatDiffChanges + FormatDiffAdded/Modified/Deleted per path

## PRINT_DIFF_RESULT

SPEC-ID: IMPL-DIFF_COMMAND::PRINT_DIFF_RESULT
STEP T001: FORMAT diff AND route through delayed collector or stdout

- [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [REQ-OUTPUT_FORMATTING] [REQ-DIFF_COMMAND] — How: format diff then route through delayed collector or stdout.

PROCEDURE PRINT_DIFF_RESULT(diff):
  message = FORMAT_DIFF_RESULT(diff); print or AddStdout
