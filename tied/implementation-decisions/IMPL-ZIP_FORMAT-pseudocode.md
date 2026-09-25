# [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP]

## Summary contract

Creates and lists ZIP archives with atomic temp-file writes, deflate compression, exclusion-aware file collection, and git-aware naming.

INPUT: context, Config, note, dryRun flag, archive directory paths
OUTPUT: zip file on disk, Archive metadata list
DATA: file paths relative to cwd, ArchiveConfig adapters, ResourceManager temp files

## CREATE_FULL_ARCHIVE

SPEC-ID: IMPL-ZIP_FORMAT::CREATE_FULL_ARCHIVE
STEP T001: Validate cwd, collect files, generate full archive name, and delegate atomic zip write unless dry-run

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: validate cwd, collect files, generate full archive name, and delegate atomic zip write unless dry-run.

PROCEDURE CREATE_FULL_ARCHIVE(ctx, cfg, note, dryRun):
Contract:
PRE: true
POST: true
INPUT: ctx, cfg, note, dryRun
OUTPUT: result of CREATE_FULL_ARCHIVE
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  validate directory and context
  rm = NEW_RESOURCE_MANAGER(); DEFER cleanup_with_panic_recovery
  files = COLLECT_FILES_TO_ARCHIVE(ctx, cwd, excludePatterns)
  archivePath = join(archiveDir, GENERATE_FULL_ARCHIVE_NAME(...))
  IF dryRun THEN print dry-run info; RETURN
  RETURN CREATE_AND_VERIFY_ARCHIVE(ctx, cwd, archivePath, files, rm)

## CREATE_AND_VERIFY_ARCHIVE

SPEC-ID: IMPL-ZIP_FORMAT::CREATE_AND_VERIFY_ARCHIVE
STEP T001: Write zip to path.tmp, rename to final path, untrack temp file, print created stats

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: write zip to path.tmp, rename to final path, untrack temp file, print created stats.

PROCEDURE CREATE_AND_VERIFY_ARCHIVE(ctx, cwd, path, files, rm):
Contract:
PRE: true
POST: true
INPUT: ctx, cwd, path, files, rm
OUTPUT: result of CREATE_AND_VERIFY_ARCHIVE
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  tempFile = path + ".tmp"
  rm.AddTempFile(tempFile)
  CREATE_ZIP_ARCHIVE(ctx, cwd, tempFile, files)
  rename(tempFile, path)
  rm.RemoveResource(tempFile)
  print created archive with stats

## CREATE_ZIP_ARCHIVE

SPEC-ID: IMPL-ZIP_FORMAT::CREATE_ZIP_ARCHIVE
STEP T001: Create zip writer with Deflate and add each collected file with context cancellation checks

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: create zip writer with Deflate and add each collected file with context cancellation checks.

PROCEDURE CREATE_ZIP_ARCHIVE(ctx, sourceDir, archivePath, files, config):
Contract:
PRE: true
POST: true
INPUT: ctx, sourceDir, archivePath, files, config
OUTPUT: result of CREATE_ZIP_ARCHIVE
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  CHECK context cancellation
  OPEN archivePath; CREATE zip.NewWriter
  FOR EACH rel IN files:
    CHECK context cancellation
    ADD_FILE_TO_ZIP(sourceDir, rel, zipWriter, config)

## ADD_FILE_TO_ZIP

SPEC-ID: IMPL-ZIP_FORMAT::ADD_FILE_TO_ZIP
STEP T001: Write one entry with Deflate header, copying file content or symlink target and honoring skip_broken_symlinks

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] [REQ-IMMUTABLE_PLATFORM_COMPATIBILITY] — How: normalize rel to OracleZipEntryPath for zip member name, then write one Deflate entry copying file content or symlink target and honoring skip_broken_symlinks.

PROCEDURE ADD_FILE_TO_ZIP(sourceDir, rel, zipWriter, config):
Contract:
PRE: true
POST: hdr.Name equals OracleZipEntryPath(rel) with forward slashes only and no ./ prefix
INPUT: sourceDir, rel, zipWriter, config
OUTPUT: result of ADD_FILE_TO_ZIP
EFFECTS: FS.Write, State.ResourceManager; zip central directory entry name normalized
  SET entryName = OracleZipEntryPath(rel)
  SET absPath = JOIN(sourceDir, FromSlash(entryName)) for filesystem access
  BUILD zip header with Deflate method; hdr.Name = entryName
  IF symlink AND broken AND config.skip_broken_symlinks THEN RETURN
  WRITE file or link target bytes to zip entry

## GENERATE_ARCHIVE_NAME

SPEC-ID: IMPL-ZIP_FORMAT::GENERATE_ARCHIVE_NAME
STEP T001: Route to incremental or full archive name builder based on config flags

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: route to incremental or full archive name builder based on config flags.

PROCEDURE GENERATE_ARCHIVE_NAME(config):
Contract:
PRE: true
POST: true
INPUT: config
OUTPUT: result of GENERATE_ARCHIVE_NAME
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  IF config.is_incremental AND config.base_name != "" THEN
    RETURN GENERATE_INCREMENTAL_ARCHIVE_NAME(config)
  RETURN GENERATE_FULL_ARCHIVE_NAME(config)

## GENERATE_FULL_ARCHIVE_NAME

SPEC-ID: IMPL-ZIP_FORMAT::GENERATE_FULL_ARCHIVE_NAME
STEP T001: Build {prefix}-{timestamp}[={branch}={hash}[-dirty]][={note}].zip from config git segments

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: build {prefix}-{timestamp}[={branch}={hash}[-dirty]][={note}].zip from config git segments.

PROCEDURE GENERATE_FULL_ARCHIVE_NAME(config):
Contract:
PRE: true
POST: true
INPUT: config
OUTPUT: result of GENERATE_FULL_ARCHIVE_NAME
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  ASSEMBLE prefix, timestamp, optional git info, optional note
  RETURN name + ".zip"

## GENERATE_INCREMENTAL_ARCHIVE_NAME

SPEC-ID: IMPL-ZIP_FORMAT::GENERATE_INCREMENTAL_ARCHIVE_NAME
STEP T001: Build {base}_update={timestamp}[={branch}={hash}[-dirty]][={note}].zip for incremental archives

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: build {base}_update={timestamp}[={branch}={hash}[-dirty]][={note}].zip for incremental archives.

PROCEDURE GENERATE_INCREMENTAL_ARCHIVE_NAME(config):
Contract:
PRE: true
POST: true
INPUT: config
OUTPUT: result of GENERATE_INCREMENTAL_ARCHIVE_NAME
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  base = strip .zip suffix from base_name
  ASSEMBLE base + "_update=" + timestamp + git segments + note
  RETURN name + ".zip"

## COLLECT_FILES_TO_ARCHIVE

SPEC-ID: IMPL-ZIP_FORMAT::COLLECT_FILES_TO_ARCHIVE
STEP T001: Walk cwd tree, skip excluded paths and directories, collect relative file paths with cancellation checks

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: walk cwd tree, skip excluded paths and directories, collect relative file paths with cancellation checks.

PROCEDURE COLLECT_FILES_TO_ARCHIVE(ctx, cwd, excludePatterns):
Contract:
PRE: true
POST: true
INPUT: ctx, cwd, excludePatterns
OUTPUT: result of COLLECT_FILES_TO_ARCHIVE
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  files = []
  WALK cwd:
    CHECK context cancellation
    IF should_exclude(rel) THEN skip dir or file
    IF regular file THEN append rel
  RETURN files

## LIST_ARCHIVES

SPEC-ID: IMPL-ZIP_FORMAT::LIST_ARCHIVES
STEP T001: Read archive directory entries and return Archive metadata for each .zip file

- [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP] — How: read archive directory entries and return Archive metadata for each .zip file.

PROCEDURE LIST_ARCHIVES(archiveDir):
Contract:
PRE: true
POST: true
INPUT: archiveDir
OUTPUT: result of LIST_ARCHIVES
EFFECTS: FS.Write, FS.Rename, State.ResourceManager
  FOR EACH directory entry:
    IF name ends with .zip THEN append createArchiveFromEntry result
  RETURN archives
