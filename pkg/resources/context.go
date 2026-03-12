// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// Context-aware operations for resource management and cancellation support.
// Provides ContextualOperation (context + ResourceManager bundle), cancellation
// checking, context key management, and context-aware cleanup.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package resources

import (
	"context"
	"os"
)

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// CONTEXTUAL_OPERATION: bundles a context.Context with a ResourceManager.
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// CONTEXTUAL_OPERATION: creates a new ContextualOperation with embedded ResourceManager.
func NewContextualOperation(ctx context.Context) *ContextualOperation {
	return &ContextualOperation{
		ctx: ctx,
		rm:  NewResourceManager(),
	}
}

// [IMPL-CONTEXT_OPS] — returns the wrapped context.
func (co *ContextualOperation) Context() context.Context {
	return co.ctx
}

// [IMPL-CONTEXT_OPS] [IMPL-RESOURCE_MANAGER] — returns the embedded ResourceManager.
func (co *ContextualOperation) ResourceManager() *ResourceManager {
	return co.rm
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// IS_CANCELLED: checks context Done channel for cancellation signal.
func (co *ContextualOperation) IsCancelled() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// CHECK_CANCELLATION: surfaces cancellation error from context.
func (co *ContextualOperation) CheckCancellation() error {
	return co.ctx.Err()
}

// [IMPL-CONTEXT_OPS] [IMPL-RESOURCE_MANAGER] — CLEANUP: delegates to ResourceManager.
func (co *ContextualOperation) Cleanup() error {
	return co.rm.Cleanup()
}

// [IMPL-CONTEXT_OPS] [IMPL-RESOURCE_MANAGER] — CLEANUP_WITH_PANIC_RECOVERY: delegates with panic safety.
func (co *ContextualOperation) CleanupWithPanicRecovery() error {
	return co.rm.CleanupWithPanicRecovery()
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// WITH_RESOURCE_MANAGER: embeds a ResourceManager into a context for downstream retrieval.
func WithResourceManager(ctx context.Context) (context.Context, *ResourceManager) {
	rm := NewResourceManager()
	return context.WithValue(ctx, ResourceManagerKey, rm), rm
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// CHECK_CONTEXT_AND_CLEANUP: checks cancellation, triggers cleanup on cancel.
func CheckContextAndCleanup(ctx context.Context, rm *ResourceManager) error {
	if err := ctx.Err(); err != nil {
		// Context is cancelled, perform cleanup
		if cleanupErr := rm.CleanupWithPanicRecovery(); cleanupErr != nil {
			// Return both the context error and cleanup error
			return CombineErrors(err, cleanupErr)
		}
		return err
	}
	return nil
}

// [IMPL-CONTEXT_OPS] [IMPL-STRUCTURED_ERRORS] — CombineErrors: merges multiple errors into one.
func CombineErrors(errors ...error) error {
	var validErrors []error
	for _, err := range errors {
		if err != nil {
			validErrors = append(validErrors, err)
		}
	}

	if len(validErrors) == 0 {
		return nil
	}

	if len(validErrors) == 1 {
		return validErrors[0]
	}

	// Combine multiple errors into a descriptive message
	var message string
	for i, err := range validErrors {
		if i == 0 {
			message = err.Error()
		} else {
			message += "; " + err.Error()
		}
	}

	return &CombinedError{
		Message: "multiple errors occurred: " + message,
		Errors:  validErrors,
	}
}

// [IMPL-CONTEXT_OPS] [IMPL-STRUCTURED_ERRORS] — CombinedError: multi-error container.
type CombinedError struct {
	Message string
	Errors  []error
}

// [IMPL-CONTEXT_OPS] — Error returns the combined error message.
func (ce *CombinedError) Error() string {
	return ce.Message
}

// [IMPL-CONTEXT_OPS] — Unwrap returns the first error for error chain compatibility.
func (ce *CombinedError) Unwrap() error {
	if len(ce.Errors) > 0 {
		return ce.Errors[0]
	}
	return nil
}

// [IMPL-CONTEXT_OPS] — GetAllErrors returns all constituent errors.
func (ce *CombinedError) GetAllErrors() []error {
	return ce.Errors
}

// [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// AtomicWriteFile writes data to a file atomically using a temporary file.
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	return AtomicWriteFileWithContext(context.Background(), path, data, rm)
}

// [IMPL-ATOMIC_OPS] [IMPL-CONTEXT_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] [REQ-CONTEXT_SUPPORT]
// AtomicWriteFileWithContext: atomic file write with cancellation checks at each stage.
func AtomicWriteFileWithContext(ctx context.Context, path string, data []byte, rm *ResourceManager) error {
	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	tempFile := path + ".tmp"
	rm.AddTempFile(tempFile)

	// Check for cancellation before writing
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}

	// Check for cancellation before finalizing
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.Rename(tempFile, path); err != nil {
		return err
	}

	// Remove from resource tracking since operation succeeded
	rm.RemoveResource(&TempFile{Path: tempFile})
	return nil
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// ContextualOperationWithTimeout: creates a ContextualOperation with a timeout.
func ContextualOperationWithTimeout(ctx context.Context, timeout int64) (*ContextualOperation, context.CancelFunc) {
	// For now, we create without timeout - in a real implementation this would use context.WithTimeout
	// The timeout parameter is included for future enhancement
	return NewContextualOperation(ctx), func() {}
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// ContextKey: typed key for context value storage.
type ContextKey string

const (
	// ResourceManagerKey is the context key for storing ResourceManager
	ResourceManagerKey ContextKey = "resourceManager"
	// OperationIDKey is the context key for storing operation IDs
	OperationIDKey ContextKey = "operationID"
)

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// GET_RESOURCE_MANAGER: retrieves a ResourceManager from context.
func GetResourceManagerFromContext(ctx context.Context) (*ResourceManager, bool) {
	rm, ok := ctx.Value(ResourceManagerKey).(*ResourceManager)
	return rm, ok
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// WITH_OPERATION_ID: stores an operation ID in the context.
func WithOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, OperationIDKey, operationID)
}

// [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT]
// GET_OPERATION_ID: retrieves an operation ID from context.
func GetOperationIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(OperationIDKey).(string)
	return id, ok
}
