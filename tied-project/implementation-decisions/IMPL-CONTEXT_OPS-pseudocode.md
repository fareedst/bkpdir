# [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]

## Summary contract

Bundle context.Context with ResourceManager for cancellable operations, expose non-blocking cancellation checks, embed managers in context values, and clean up tracked resources when contexts end.

INPUT: parent context, operation ID, timeout (optional)
OUTPUT: ContextualOperation, derived contexts, cleanup errors
DATA: ResourceManagerKey, OperationIDKey, ContextualOperation struct

## NEW_CONTEXTUAL_OPERATION

SPEC-ID: IMPL-CONTEXT_OPS::NEW_CONTEXTUAL_OPERATION
STEP T001: WRAP ctx with newly allocated ResourceManager

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: wrap ctx with a freshly allocated ResourceManager for scoped cleanup during the operation.

PROCEDURE NEW_CONTEXTUAL_OPERATION(ctx):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: ctx
OUTPUT: result of NEW_CONTEXTUAL_OPERATION
PRE: true
POST: true
EFFECTS: pure
  RETURN ContextualOperation{ctx: ctx, rm: NEW ResourceManager}

## IS_CANCELLED

SPEC-ID: IMPL-CONTEXT_OPS::IS_CANCELLED
STEP T001: NON-BLOCKING select on ctx.Done for cancellation signal

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: non-blocking select on ctx.Done(); return true when cancellation signaled.

PROCEDURE IS_CANCELLED(co):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: co
OUTPUT: result of IS_CANCELLED
PRE: true
POST: true
EFFECTS: pure
FAILURE_MODES: none
  SELECT co.ctx.Done:
    CASE signaled: RETURN true
    DEFAULT: RETURN false

## CHECK_CONTEXT_AND_CLEANUP

SPEC-ID: IMPL-CONTEXT_OPS::CHECK_CONTEXT_AND_CLEANUP
STEP T001: WHEN ctx cancelled RUN CleanupWithPanicRecovery and combine errors

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: when ctx.Err() is set, run CleanupWithPanicRecovery and combine cleanup failure with context error.

PROCEDURE CHECK_CONTEXT_AND_CLEANUP(ctx, rm):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: ctx, rm
OUTPUT: result of CHECK_CONTEXT_AND_CLEANUP
PRE: true
POST: true
EFFECTS: pure
  IF ctx.Err() != nil:
    cleanupErr = rm.CleanupWithPanicRecovery()
    IF cleanupErr != nil: RETURN CombineErrors(ctx.Err(), cleanupErr)
    RETURN ctx.Err()
  RETURN nil

## WITH_RESOURCE_MANAGER

SPEC-ID: IMPL-CONTEXT_OPS::WITH_RESOURCE_MANAGER
STEP T001: STORE new ResourceManager in context under ResourceManagerKey

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: store a new ResourceManager in context under ResourceManagerKey for downstream retrieval.

PROCEDURE WITH_RESOURCE_MANAGER(ctx):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: ctx
OUTPUT: result of WITH_RESOURCE_MANAGER
PRE: true
POST: true
EFFECTS: pure
  rm = NEW ResourceManager
  RETURN context.WithValue(ctx, ResourceManagerKey, rm), rm

## WITH_OPERATION_ID

SPEC-ID: IMPL-CONTEXT_OPS::WITH_OPERATION_ID
STEP T001: ATTACH operationID to context under OperationIDKey

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: attach operationID to context under OperationIDKey for tracing nested work.

PROCEDURE WITH_OPERATION_ID(ctx, operationID):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: ctx, operationID
OUTPUT: result of WITH_OPERATION_ID
PRE: true
POST: true
EFFECTS: pure
  RETURN context.WithValue(ctx, OperationIDKey, operationID)

## GET_OPERATION_ID

SPEC-ID: IMPL-CONTEXT_OPS::GET_OPERATION_ID
STEP T001: READ OperationIDKey from context as string

- [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: read OperationIDKey from context and report whether a string ID was stored.

PROCEDURE GET_OPERATION_ID(ctx):
  # [IMPL-CONTEXT_OPS] — Layer B procedure token link
Contract:
INPUT: ctx
OUTPUT: result of GET_OPERATION_ID
PRE: true
POST: true
EFFECTS: pure
  id, ok = ctx.Value(OperationIDKey).(string)
  RETURN id, ok
