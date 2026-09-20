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

// CONTEXTUAL_OPERATION: bundles a context.Context with a ResourceManager.
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: wrap ctx with a freshly allocated ResourceManager for scoped cleanup during the operation.
func NewContextualOperation(ctx context.Context) *ContextualOperation {
	return &ContextualOperation{
		ctx: ctx,
		rm:  NewResourceManager(),
	}
}

func (co *ContextualOperation) Context() context.Context {
	return co.ctx
}

func (co *ContextualOperation) ResourceManager() *ResourceManager {
	return co.rm
}

// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: non-blocking select on ctx.Done(); return true when cancellation signaled.
func (co *ContextualOperation) IsCancelled() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

func (co *ContextualOperation) CheckCancellation() error {
	return co.ctx.Err()
}

func (co *ContextualOperation) Cleanup() error {
	return co.rm.Cleanup()
}

func (co *ContextualOperation) CleanupWithPanicRecovery() error {
	return co.rm.CleanupWithPanicRecovery()
}

// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: store a new ResourceManager in context under ResourceManagerKey for downstream retrieval.
func WithResourceManager(ctx context.Context) (context.Context, *ResourceManager) {
	rm := NewResourceManager()
	return context.WithValue(ctx, ResourceManagerKey, rm), rm
}

// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: when ctx.Err() is set, run CleanupWithPanicRecovery and combine cleanup failure with context error.
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

type CombinedError struct {
	Message string
	Errors  []error
}

func (ce *CombinedError) Error() string {
	return ce.Message
}

func (ce *CombinedError) Unwrap() error {
	if len(ce.Errors) > 0 {
		return ce.Errors[0]
	}
	return nil
}

func (ce *CombinedError) GetAllErrors() []error {
	return ce.Errors
}

// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: delegate to context-aware atomic write with background context.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: write byte slice through AtomicWriter, set permissions on temp, commit.
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	return AtomicWriteFileWithContext(context.Background(), path, data, rm)
}

// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: check cancellation before each stage, write temp path, rename to final, untrack temp on success.
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

// ContextualOperationWithTimeout: creates a ContextualOperation with a timeout.
func ContextualOperationWithTimeout(ctx context.Context, timeout int64) (*ContextualOperation, context.CancelFunc) {
	// For now, we create without timeout - in a real implementation this would use context.WithTimeout
	// The timeout parameter is included for future enhancement
	return NewContextualOperation(ctx), func() {}
}

// ContextKey: typed key for context value storage.
type ContextKey string

const (
	// ResourceManagerKey is the context key for storing ResourceManager
	ResourceManagerKey ContextKey = "resourceManager"
	// OperationIDKey is the context key for storing operation IDs
	OperationIDKey ContextKey = "operationID"
)

// GET_RESOURCE_MANAGER: retrieves a ResourceManager from context.
func GetResourceManagerFromContext(ctx context.Context) (*ResourceManager, bool) {
	rm, ok := ctx.Value(ResourceManagerKey).(*ResourceManager)
	return rm, ok
}

// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: attach operationID to context under OperationIDKey for tracing nested work.
func WithOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, OperationIDKey, operationID)
}

// SPEC-ID: IMPL-CONTEXT_OPS::GET_OPERATION_ID
// - [IMPL-CONTEXT_OPS] [ARCH-CONTEXT_SUPPORT] [REQ-CONTEXT_SUPPORT] — How: read OperationIDKey from context and report whether a string ID was stored.
func GetOperationIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(OperationIDKey).(string)
	return id, ok
}
