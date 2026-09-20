# [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING]

## Summary contract

Structured domain errors with status codes, operation and path context, unwrap support, constructors at increasing detail levels, and classifiers that map low-level failures to disk, permission, and not-found conditions.

INPUT: message, statusCode, operation, path, underlying err
OUTPUT: ErrorInterface values, boolean classification, exit status from handlers
DATA: ArchiveError, BackupError, error pattern lists, ErrorCategory

## ARCHIVEERROR_ERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::ARCHIVEERROR_ERROR
STEP T001: Format Message alone or Message WITH underlying Err when present

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: format Message alone or Message WITH underlying Err when present.

PROCEDURE ARCHIVEERROR_ERROR(e):
  IF e.Err != nil THEN RETURN FORMAT(message, underlying)
  RETURN e.Message

## ARCHIVEERROR_UNWRAP

SPEC-ID: IMPL-STRUCTURED_ERRORS::ARCHIVEERROR_UNWRAP
STEP T001: Return wrapped Err for errors.Is/As chains

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: return wrapped Err for errors.Is/As chains.

PROCEDURE ARCHIVEERROR_UNWRAP(e):
  IF e == nil THEN
    RETURN nil
  END IF
  RETURN underlying error stored in e.Err

## BACKUPERROR_ERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::BACKUPERROR_ERROR
STEP T001: Mirror ArchiveError Error formatting for backup operations

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: mirror ArchiveError Error formatting for backup operations.

PROCEDURE BACKUPERROR_ERROR(e):
  IF e.Err != nil THEN RETURN FORMAT(message, underlying)
  RETURN e.Message

## NEWARCHIVEERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWARCHIVEERROR
STEP T001: Construct ArchiveError with message and status code only

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message and status code only.

PROCEDURE NEWARCHIVEERROR(message, statusCode):
  DECLARE result as ArchiveError
  SET result.Message = message
  SET result.StatusCode = statusCode
  LEAVE Operation, Path, Err empty
  RETURN pointer to result

## NEWARCHIVEERRORWITHCAUSE

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWARCHIVEERRORWITHCAUSE
STEP T001: Construct ArchiveError with message, status code, and underlying Err

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status code, and underlying Err.

PROCEDURE NEWARCHIVEERRORWITHCAUSE(message, statusCode, err):
  DECLARE result as ArchiveError
  SET result.Message = message
  SET result.StatusCode = statusCode
  SET result.Err = err
  LEAVE Operation and Path empty
  RETURN pointer to result

## NEWARCHIVEERRORWITHCONTEXT

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWARCHIVEERRORWITHCONTEXT
STEP T001: Construct ArchiveError with message, status, operation, path, and optional Err

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status, operation, path, and optional Err.

PROCEDURE NEWARCHIVEERRORWITHCONTEXT(message, statusCode, operation, path, err):
  RETURN ArchiveError{Message, StatusCode, Operation: operation, Path: path, Err: err}

## NEWBACKUPERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWBACKUPERROR
STEP T001: Construct BackupError with message and status code only

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message and status code only.

PROCEDURE NEWBACKUPERROR(message, statusCode):
  DECLARE result as BackupError
  SET result.Message = message
  SET result.StatusCode = statusCode
  LEAVE Operation, Path, Err empty
  RETURN pointer to result

## NEWBACKUPERRORWITHCAUSE

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWBACKUPERRORWITHCAUSE
STEP T001: Construct BackupError with message, status code, and underlying Err

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status code, and underlying Err.

PROCEDURE NEWBACKUPERRORWITHCAUSE(message, statusCode, err):
  RETURN BackupError{Message, StatusCode, Err: err}

## NEWBACKUPERRORWITHCONTEXT

SPEC-ID: IMPL-STRUCTURED_ERRORS::NEWBACKUPERRORWITHCONTEXT
STEP T001: Construct BackupError with message, status, operation, path, and optional Err

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status, operation, path, and optional Err.

PROCEDURE NEWBACKUPERRORWITHCONTEXT(message, statusCode, operation, path, err):
  RETURN BackupError{Message, StatusCode, Operation: operation, Path: path, Err: err}

## ISDISKFULLERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::ISDISKFULLERROR
STEP T001: Match disk-space substrings in error text OR path error with OS no-space/quota/large-file codes

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match disk-space substrings in error text OR path error with OS no-space/quota/large-file codes.

PROCEDURE ISDISKFULLERROR(err):
  IF err == nil THEN RETURN false
  lowered = ToLower(err.Error())
  RETURN matches disk-full patterns OR pathErr has ENOSPC/EDQUOT/EFBIG

## ISPERMISSIONERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::ISPERMISSIONERROR
STEP T001: Match permission-denied substrings OR path error with EACCES/EPERM

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match permission-denied substrings OR path error with EACCES/EPERM.

PROCEDURE ISPERMISSIONERROR(err):
  IF err == nil THEN RETURN false
  lowered = ToLower(err.Error())
  RETURN matches permission patterns OR pathErr has EACCES/EPERM

## ISDIRECTORYNOTFOUNDERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::ISDIRECTORYNOTFOUNDERROR
STEP T001: Match directory-not-found substrings OR path error with ENOENT

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match directory-not-found substrings OR path error with ENOENT.

PROCEDURE ISDIRECTORYNOTFOUNDERROR(err):
  IF err == nil THEN RETURN false
  lowered = ToLower(err.Error())
  RETURN matches not-found patterns OR pathErr has ENOENT

## HANDLEERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::HANDLEERROR
STEP T001: Route ApplicationError or classified errors to formatter and configured status codes

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: route ApplicationError or classified errors to formatter and configured status codes.

PROCEDURE HANDLEERROR(err, cfg, formatter):
  IF err == nil THEN RETURN 0
  IF ApplicationError THEN RETURN HandleApplicationError
  IF ISDISKFULLERROR THEN print disk full; RETURN status disk_full OR 1
  IF ISPERMISSIONERROR THEN print permission; RETURN status permission_denied OR 1
  IF ISDIRECTORYNOTFOUND THEN print directory not found; RETURN matching status OR 1
  IF ISFILENOTFOUND THEN print file not found; RETURN matching status OR 1
  ELSE print generic error; RETURN 1

## CLASSIFYERROR

SPEC-ID: IMPL-STRUCTURED_ERRORS::CLASSIFYERROR
STEP T001: Map error to ErrorCategory via disk, permission, filesystem, network detectors

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: map error to ErrorCategory via disk, permission, filesystem, network detectors.

PROCEDURE CLASSIFYERROR(err):
  IF err == nil THEN RETURN Unknown
  IF ISDISKFULLERROR THEN RETURN DiskSpace
  IF ISPERMISSIONERROR THEN RETURN Permission
  IF ISDIRECTORYNOTFOUND OR ISFILENOTFOUND THEN RETURN Filesystem
  IF ISNETWORKERROR THEN RETURN Network
  RETURN Unknown

## EMBEDDED_MINITEST: archive error string

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify ArchiveError.Error with and without underlying cause.

## EMBEDDED_MINITEST: disk full classifier

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify IsDiskFullError matches message patterns and syscall path errors.

## EMBEDDED_MINITEST: permission classifier

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify IsPermissionError matches permission strings and EACCES/EPERM.

## EMBEDDED_MINITEST: pkg handle error

- [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: verify HandleError returns configured status codes for classified errors.
