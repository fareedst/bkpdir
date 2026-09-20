# [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING]

## Summary contract

Skip creating incremental archives when reconstructed archive state shows no added, modified, or deleted files versus the working tree, using diff command primitives and a configurable skip message.

INPUT: IncrementalArchiveConfig, archiveDir, cwd, exclude patterns
OUTPUT: nil when skipped, archive path when changes exist
DATA: reconstructedState, DirectoryDiff, FormatIncrementalSkippedNoChanges

## CREATE_INCREMENTAL_ARCHIVE

SPEC-ID: IMPL-INCREMENTAL_DUPLICATE_PREVENTION::CREATE_INCREMENTAL_ARCHIVE
STEP T001: RECONSTRUCT archive state AND skip when diff has no changes
STEP T002: CONTINUE incremental creation with modified file list

- [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING] — How: reconstruct archive state, CalculateDiff against cwd, skip creation and print skip message when diff has no added/modified/deleted entries.

PROCEDURE CREATE_INCREMENTAL_ARCHIVE(config):
  archiveDir = prepareArchiveDirectory(...)
  reconstructedState, err = ReconstructArchiveState(archiveDir)
  IF err: FALLBACK collectModifiedFiles from latest full archive
  ELSE:
    diff = CalculateDiff(cwd, reconstructedState, excludePatterns)
    IF diff empty: PrintIncrementalSkippedNoChanges(); RETURN nil
    modifiedFiles = diff.Added + diff.Modified
  CONTINUE incremental archive creation with modifiedFiles

## PRINT_INCREMENTAL_SKIPPED_NO_CHANGES

SPEC-ID: IMPL-INCREMENTAL_DUPLICATE_PREVENTION::PRINT_INCREMENTAL_SKIPPED_NO_CHANGES
STEP T001: FORMAT skip message AND emit via collector or stdout

- [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [ARCH-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-DIFF_COMMAND] [REQ-INCREMENTAL_DUPLICATE_PREVENTION] [REQ-OUTPUT_FORMATTING] — How: format cfg.FormatIncrementalSkippedNoChanges and emit via delayed collector or stdout.

PROCEDURE PRINT_INCREMENTAL_SKIPPED_NO_CHANGES(formatter):
  message = cfg.FormatIncrementalSkippedNoChanges
  IF collector THEN AddStdout(message) ELSE print message
