# [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE]

## Summary contract

Snapshot-based directory and zip archive comparison via pkg/fileops wrappers, plus helpers to find the latest full archive and detect directory-archive identity.

INPUT: root path, archive zip path, exclude patterns
OUTPUT: DirectorySnapshot, boolean identical, most recent archive path
DATA: FileInfo RelativePath, Size, Hash

## CREATE_DIRECTORY_SNAPSHOT

SPEC-ID: IMPL-DIRECTORY_COMPARISON::CREATE_DIRECTORY_SNAPSHOT
STEP T001: DELEGATE to fileops.CreateDirectorySnapshot with exclude patterns

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateDirectorySnapshot with exclusion patterns applied during walk.

PROCEDURE CREATE_DIRECTORY_SNAPSHOT(root, excludePatterns):
  RETURN fileops.CreateDirectorySnapshot(root, excludePatterns)

## CREATE_ARCHIVE_SNAPSHOT

SPEC-ID: IMPL-DIRECTORY_COMPARISON::CREATE_ARCHIVE_SNAPSHOT
STEP T001: DELEGATE to fileops.CreateArchiveSnapshot for zip members

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateArchiveSnapshot to enumerate zip member files with metadata.

PROCEDURE CREATE_ARCHIVE_SNAPSHOT(archivePath):
  RETURN fileops.CreateArchiveSnapshot(archivePath)

## COMPARE_SNAPSHOTS

SPEC-ID: IMPL-DIRECTORY_COMPARISON::COMPARE_SNAPSHOTS
STEP T001: DELEGATE equality check to fileops.CompareSnapshots

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate equality check of two DirectorySnapshot values to fileops.CompareSnapshots.

PROCEDURE COMPARE_SNAPSHOTS(a, b):
  RETURN fileops.CompareSnapshots(a, b)

## IS_DIRECTORY_IDENTICAL_TO_ARCHIVE

SPEC-ID: IMPL-DIRECTORY_COMPARISON::IS_DIRECTORY_IDENTICAL_TO_ARCHIVE
STEP T001: BUILD dir and archive snapshots WITH exclusions
STEP T002: COMPARE snapshots for structural equality

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: build dir and archive snapshots with exclusions and compare for full structural equality.

PROCEDURE IS_DIRECTORY_IDENTICAL_TO_ARCHIVE(dir, archive, excludePatterns):
  RETURN fileops.IsDirectoryIdenticalToArchive(dir, archive, excludePatterns)

## FIND_MOST_RECENT_ARCHIVE

SPEC-ID: IMPL-DIRECTORY_COMPARISON::FIND_MOST_RECENT_ARCHIVE
STEP T001: FILTER full archives only
STEP T002: SORT by name AND return last entry path

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: list archives, keep full archives only, sort by name, return path of last entry.

PROCEDURE FIND_MOST_RECENT_ARCHIVE(archiveDir):
  fullArchives = FILTER NOT incremental; SORT Name; RETURN last.Path

## CHECK_FOR_IDENTICAL_ARCHIVE

SPEC-ID: IMPL-DIRECTORY_COMPARISON::CHECK_FOR_IDENTICAL_ARCHIVE
STEP T001: FIND_MOST_RECENT_ARCHIVE in archiveDir
STEP T002: IS_DIRECTORY_IDENTICAL_TO_ARCHIVE against cwd

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: FindMostRecentArchive then IsDirectoryIdenticalToArchive against cwd.

PROCEDURE CHECK_FOR_IDENTICAL_ARCHIVE(dir, archiveDir, excludePatterns):
  recent = FIND_MOST_RECENT_ARCHIVE
  IF none: RETURN false
  RETURN IS_DIRECTORY_IDENTICAL_TO_ARCHIVE(dir, recent, excludePatterns)

## GET_DIRECTORY_TREE_SUMMARY

SPEC-ID: IMPL-DIRECTORY_COMPARISON::GET_DIRECTORY_TREE_SUMMARY
STEP T001: CREATE_DIRECTORY_SNAPSHOT for dirPath with exclusions
STEP T002: EMIT header AND list each RelativePath and Size

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: build human-readable directory listing from CreateDirectorySnapshot for diagnostics.

PROCEDURE GET_DIRECTORY_TREE_SUMMARY(dirPath, excludePatterns):
  snapshot = CREATE_DIRECTORY_SNAPSHOT(dirPath, excludePatterns)
  EMIT header with path and file count; LIST each RelativePath and Size

## GET_ARCHIVE_TREE_SUMMARY

SPEC-ID: IMPL-DIRECTORY_COMPARISON::GET_ARCHIVE_TREE_SUMMARY
STEP T001: CREATE_ARCHIVE_SNAPSHOT for archivePath
STEP T002: EMIT header AND list each RelativePath and Size

- [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: build human-readable file listing from CreateArchiveSnapshot for diagnostics and tests.

PROCEDURE GET_ARCHIVE_TREE_SUMMARY(archivePath):
  snapshot = CREATE_ARCHIVE_SNAPSHOT(archivePath)
  EMIT header with path and file count; LIST each RelativePath and Size
