# [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]

## Summary contract

Thread-safe tracking of temporary files and directories with automatic cleanup, panic recovery, and context-aware cancellation. Archive and backup creation uses atomic operations via coordinated temp-file lifecycle.

INPUT: Resource instances, temp paths, optional context
OUTPUT: cleanup errors, updated resource list
DATA: ResourceManager slice guarded by mutex, TempFile and TempDir implementations

## ADD_RESOURCE

SPEC-ID: IMPL-RESOURCE_MANAGER::ADD_RESOURCE
STEP T001: Append a Resource to the tracked slice under write lock for later cleanup

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: append a Resource to the tracked slice under write lock for later cleanup.

PROCEDURE ADD_RESOURCE(rm, resource):
  LOCK rm.mutex
  APPEND resource TO rm.resources
  UNLOCK rm.mutex

## ADD_TEMP_RESOURCES

SPEC-ID: IMPL-RESOURCE_MANAGER::ADD_TEMP_RESOURCES
STEP T001: Wrap AddResource with TempFile or TempDir for common temp path registration

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: wrap AddResource with TempFile or TempDir for common temp path registration.

PROCEDURE ADD_TEMP_FILE(rm, path):
  CALL ADD_RESOURCE(rm, TempFile{path})
PROCEDURE ADD_TEMP_DIR(rm, path):
  CALL ADD_RESOURCE(rm, TempDir{path})

## REMOVE_RESOURCE

SPEC-ID: IMPL-RESOURCE_MANAGER::REMOVE_RESOURCE
STEP T001: Remove matching resource from tracking by String() identity without invoking Cleanup

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: remove matching resource from tracking by String() identity without invoking Cleanup.

PROCEDURE REMOVE_RESOURCE(rm, resource):
  LOCK rm.mutex
  FIND index WHERE resources[i].String() == resource.String()
  REMOVE element at index
  UNLOCK rm.mutex

## CLEANUP

SPEC-ID: IMPL-RESOURCE_MANAGER::CLEANUP
STEP T001: Invoke Cleanup on every tracked resource, retain last error, then clear the slice

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: invoke Cleanup on every tracked resource, retain last error, then clear the slice.

PROCEDURE CLEANUP(rm):
  LOCK rm.mutex
  FOR EACH resource IN rm.resources:
    err = resource.Cleanup()
    IF err != nil THEN lastError = err
  CLEAR rm.resources
  UNLOCK rm.mutex
  RETURN lastError

## CLEANUP_WITH_PANIC_RECOVERY

SPEC-ID: IMPL-RESOURCE_MANAGER::CLEANUP_WITH_PANIC_RECOVERY
STEP T001: Defer recover around CLEANUP and convert panics into returned errors

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: defer recover around CLEANUP and convert panics into returned errors.

PROCEDURE CLEANUP_WITH_PANIC_RECOVERY(rm):
  DEFER recover panic INTO err
  RETURN CLEANUP(rm)

## CLEANUP_WITH_CONTEXT

SPEC-ID: IMPL-RESOURCE_MANAGER::CLEANUP_WITH_CONTEXT
STEP T001: Return ctx.Err() immediately or between resources when context is cancelled during cleanup

- [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] [REQ-CONTEXT_SUPPORT] — How: return ctx.Err() immediately or between resources when context is cancelled during cleanup.

PROCEDURE CLEANUP_WITH_CONTEXT(rm, ctx):
  IF ctx.Err() != nil THEN RETURN ctx.Err()
  LOCK rm.mutex
  FOR EACH resource IN rm.resources:
    IF ctx cancelled THEN RETURN ctx.Err()
    resource.Cleanup()
  CLEAR rm.resources
  UNLOCK rm.mutex
