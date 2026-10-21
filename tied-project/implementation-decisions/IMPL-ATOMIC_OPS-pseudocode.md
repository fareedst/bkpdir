# [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT]

## Summary contract

Atomic file I/O via write-to-temporary-file then rename on the same filesystem, with rollback on failure and convenience helpers for copy and full-file writes.

INPUT: targetPath, src, dst, data, permissions
OUTPUT: committed file at target or ERROR
DATA: AtomicWriter state (targetPath, tempPath, temp handle, committed, closed flags)

## NEWATOMICWRITER

SPEC-ID: IMPL-ATOMIC_OPS::NEWATOMICWRITER
INPUT: targetPath string
OUTPUT: AtomicWriter handle | ERROR
PRE: VALIDATE_PATH(targetPath) == ok AND PARENT_DIRECTORY(targetPath) is creatable
POST: writer.targetPath == targetPath AND writer.HasHandle == true AND writer.committed == false
EFFECTS: FS.CreateTemp, FS.MkdirAll
FAILURE_MODES: ErrInvalidPath, ErrCreateTemp, ErrMkdir
STEP T001: ENSURE_DIRECTORY_EXISTS(PARENT_DIRECTORY(targetPath))
STEP T002: CREATE_TEMP_FILE_IN(dir, pattern)

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: validate target path, ensure parent directory exists, create temp file beside target for same-filesystem rename.

PROCEDURE NEWATOMICWRITER(targetPath):
Contract:
PRE: true
POST: true
INPUT: targetPath string
OUTPUT: AtomicWriter | ERROR
EFFECTS: FS.CreateTemp, FS.MkdirAll
  IF VALIDATE_PATH(targetPath) fails THEN RETURN ERROR
  dir = PARENT_DIRECTORY(targetPath)
  ENSURE_DIRECTORY_EXISTS(dir)
  temp = CREATE_TEMP_FILE_IN(dir, pattern FROM base name + ".tmp.*")
  RETURN AtomicWriter WITH targetPath, temp path, open handle, not committed, not closed

## ATOMICWRITER_WRITE

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_WRITE
STEP T001: REJECT if writer.closed
STEP T002: WRITE_BYTES_TO_TEMP(data)
BRANCH B001: writer.closed == true
ERROR E001: RETURN ERROR "writer is closed"

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: reject writes when writer is closed or temp handle missing; otherwise write bytes to temp file.

PROCEDURE ATOMICWRITER_WRITE(writer, data):
Contract:
PRE: true
POST: true
INPUT: writer, data
OUTPUT: result of ATOMICWRITER_WRITE
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  IF writer.closed THEN RETURN ERROR "writer is closed"
  IF writer.temp handle missing THEN RETURN ERROR
  RETURN WRITE_BYTES_TO_TEMP(data)

## ATOMICWRITER_WRITESTRING

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_WRITESTRING
STEP T001: DELEGATE to ATOMICWRITER_WRITE(BYTES(text))

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: convert text to bytes and delegate to ATOMICWRITER_WRITE.

PROCEDURE ATOMICWRITER_WRITESTRING(writer, text):
Contract:
PRE: true
POST: true
INPUT: writer, text
OUTPUT: result of ATOMICWRITER_WRITESTRING
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  RETURN ATOMICWRITER_WRITE(writer, BYTES(text))

## ATOMICWRITER_COMMIT

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_COMMIT
STEP T001: CLOSE_TEMP_HANDLE(writer)
STEP T002: ATOMIC_RENAME(writer.tempPath, writer.targetPath)
POST: writer.committed == true AND writer.closed == true
BRANCH B001: writer.committed == true
ERROR E001: RETURN ERROR "already committed"

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: close temp handle then ATOMIC_RENAME temp to target; cleanup temp on failure; mark committed and closed.

PROCEDURE ATOMICWRITER_COMMIT(writer):
Contract:
PRE: true
POST: true
INPUT: writer
OUTPUT: result of ATOMICWRITER_COMMIT
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  IF writer.committed THEN RETURN ERROR "already committed"
  CLOSE_TEMP_HANDLE(writer)
  IF ATOMIC_RENAME(writer.tempPath, writer.targetPath) fails THEN CLEANUP(writer); RETURN ERROR
  SET writer.committed AND writer.closed

## ATOMICWRITER_ROLLBACK

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_ROLLBACK
STEP T001: CLEANUP(writer)
BRANCH B001: writer.committed == true
ERROR E001: RETURN ERROR "cannot rollback after commit"

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: forbid rollback after commit; remove temp file via CLEANUP.

PROCEDURE ATOMICWRITER_ROLLBACK(writer):
Contract:
PRE: true
POST: true
INPUT: writer
OUTPUT: result of ATOMICWRITER_ROLLBACK
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  IF writer.committed THEN RETURN ERROR "cannot rollback after commit"
  RETURN CLEANUP(writer)

## ATOMICWRITER_CLOSE

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_CLOSE
STEP T001: IF NOT writer.committed THEN ATOMICWRITER_ROLLBACK(writer)
POST: writer.closed == true

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: idempotent close; roll back when not committed.

PROCEDURE ATOMICWRITER_CLOSE(writer):
Contract:
PRE: true
POST: true
INPUT: writer
OUTPUT: result of ATOMICWRITER_CLOSE
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  IF writer.closed THEN RETURN success
  IF NOT writer.committed THEN RETURN ATOMICWRITER_ROLLBACK(writer)
  SET writer.closed

## ATOMICWRITER_CLEANUP

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITER_CLEANUP
STEP T001: CLOSE_HANDLE_IF_OPEN(writer)
STEP T002: REMOVE_FILE(writer.tempPath)

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: close open handle and remove temp path from disk ignoring not-exist.

PROCEDURE ATOMICWRITER_CLEANUP(writer):
Contract:
PRE: true
POST: true
INPUT: writer
OUTPUT: result of ATOMICWRITER_CLEANUP
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  CLOSE_HANDLE_IF_OPEN(writer)
  REMOVE_FILE(writer.tempPath) ignoring not-exist
  SET writer.closed

## ATOMICCOPY

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICCOPY
STEP T001: VALIDATE_READABLE(src)
STEP T002: STREAM to AtomicWriter and COMMIT

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: validate readable source and destination, stream to AtomicWriter, match permissions, commit.

PROCEDURE ATOMICCOPY(src, dst):
Contract:
PRE: true
POST: true
INPUT: src, dst
OUTPUT: result of ATOMICCOPY
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  VALIDATE_READABLE(src)
  VALIDATE_PATH(dst)
  srcInfo = STAT(src)
  writer = NEWATOMICWRITER(dst) WITH defer CLOSE
  COPY_STREAM(src, writer)
  SET_PERMISSIONS(writer.tempPath, srcInfo.mode)
  RETURN ATOMICWRITER_COMMIT(writer)

## ATOMICWRITEFILE

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITEFILE
STEP T001: ATOMICWRITER_WRITE(data)
STEP T002: ATOMICWRITER_COMMIT(writer)

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: write byte slice through AtomicWriter, set permissions on temp, commit.

PROCEDURE ATOMICWRITEFILE(filename, data, perm):
Contract:
PRE: true
POST: true
INPUT: filename, data, perm
OUTPUT: result of ATOMICWRITEFILE
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  writer = NEWATOMICWRITER(filename) WITH defer CLOSE
  ATOMICWRITER_WRITE(writer, data)
  SET_PERMISSIONS(writer.tempPath, perm)
  RETURN ATOMICWRITER_COMMIT(writer)

## ATOMICWRITESTRING

SPEC-ID: IMPL-ATOMIC_OPS::ATOMICWRITESTRING
STEP T001: DELEGATE to ATOMICWRITEFILE(BYTES(text))

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: delegate text payload to ATOMICWRITEFILE.

PROCEDURE ATOMICWRITESTRING(filename, text, perm):
Contract:
PRE: true
POST: true
INPUT: filename, text, perm
OUTPUT: result of ATOMICWRITESTRING
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  RETURN ATOMICWRITEFILE(filename, BYTES(text), perm)

## RESOURCE_ATOMICWRITEFILE

SPEC-ID: IMPL-ATOMIC_OPS::RESOURCE_ATOMICWRITEFILE
STEP T001: DELEGATE RESOURCE_ATOMICWRITEFILEWITHCONTEXT(BACKGROUND_CONTEXT)

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: delegate to context-aware atomic write with background context.

PROCEDURE RESOURCE_ATOMICWRITEFILE(path, data, resourceManager):
Contract:
PRE: true
POST: true
INPUT: path, data, resourceManager
OUTPUT: result of RESOURCE_ATOMICWRITEFILE
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  RETURN RESOURCE_ATOMICWRITEFILEWITHCONTEXT(BACKGROUND_CONTEXT, path, data, resourceManager)

## RESOURCE_ATOMICWRITEFILEWITHCONTEXT

SPEC-ID: IMPL-ATOMIC_OPS::RESOURCE_ATOMICWRITEFILEWITHCONTEXT
STEP T001: CHECK ctx cancelled
STEP T002: WRITE_FILE tempPath THEN ATOMIC_RENAME

- [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: check cancellation before each stage, write temp path, rename to final, untrack temp on success.

PROCEDURE RESOURCE_ATOMICWRITEFILEWITHCONTEXT(ctx, path, data, resourceManager):
Contract:
PRE: true
POST: true
INPUT: ctx, path, data, resourceManager
OUTPUT: result of RESOURCE_ATOMICWRITEFILEWITHCONTEXT
EFFECTS: FS.Write, FS.Rename, State.AtomicWriter
  IF ctx cancelled THEN RETURN cancellation ERROR
  tempPath = path + ".tmp"
  REGISTER_TEMP_WITH(resourceManager, tempPath)
  IF ctx cancelled THEN RETURN cancellation ERROR
  WRITE_FILE(tempPath, data)
  IF ctx cancelled THEN RETURN cancellation ERROR
  ATOMIC_RENAME(tempPath, path)
  UNREGISTER_TEMP(resourceManager, tempPath)
