// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// Context-aware operations for resource management and cancellation support.
// This file contains the context operations extracted from the original
// errors.go file, providing cancellation and timeout support for resource operations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package resources

import (
	"context"
	"os"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ContextualOperation provides context and resource management for operations
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// NewContextualOperation creates a new ContextualOperation with the given context
func NewContextualOperation(ctx context.Context) *ContextualOperation {
	return &ContextualOperation{
		ctx: ctx,
		rm:  NewResourceManager(),
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Context returns the context associated with the ContextualOperation
func (co *ContextualOperation) Context() context.Context {
	return co.ctx
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ResourceManager returns the ResourceManager associated with the ContextualOperation
func (co *ContextualOperation) ResourceManager() *ResourceManager {
	return co.rm
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// IsCancelled checks if the operation has been cancelled
func (co *ContextualOperation) IsCancelled() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CheckCancellation checks if the operation has been cancelled and returns an error if it has
func (co *ContextualOperation) CheckCancellation() error {
	return co.ctx.Err()
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Cleanup cleans up all resources associated with the operation
func (co *ContextualOperation) Cleanup() error {
	return co.rm.Cleanup()
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CleanupWithPanicRecovery cleans up all resources with panic recovery
func (co *ContextualOperation) CleanupWithPanicRecovery() error {
	return co.rm.CleanupWithPanicRecovery()
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// WithResourceManager creates a new context with an associated ResourceManager
func WithResourceManager(ctx context.Context) (context.Context, *ResourceManager) {
	rm := NewResourceManager()
	return context.WithValue(ctx, ResourceManagerKey, rm), rm
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CheckContextAndCleanup checks for context cancellation and performs cleanup if needed
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CombineErrors combines multiple errors into a single error message
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CombinedError represents multiple errors that occurred together
type CombinedError struct {
	Message string
	Errors  []error
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Error returns the combined error message
func (ce *CombinedError) Error() string {
	return ce.Message
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Unwrap returns the first error for compatibility with error unwrapping
func (ce *CombinedError) Unwrap() error {
	if len(ce.Errors) > 0 {
		return ce.Errors[0]
	}
	return nil
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetAllErrors returns all the constituent errors
func (ce *CombinedError) GetAllErrors() []error {
	return ce.Errors
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AtomicWriteFile writes data to a file atomically using a temporary file
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	return AtomicWriteFileWithContext(context.Background(), path, data, rm)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AtomicWriteFileWithContext writes data to a file atomically with context support
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ContextualOperationWithTimeout creates a ContextualOperation with a timeout
func ContextualOperationWithTimeout(ctx context.Context, timeout int64) (*ContextualOperation, context.CancelFunc) {
	// For now, we create without timeout - in a real implementation this would use context.WithTimeout
	// The timeout parameter is included for future enhancement
	return NewContextualOperation(ctx), func() {}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ContextKey represents a key for storing values in context
type ContextKey string

const (
	// ResourceManagerKey is the context key for storing ResourceManager
	ResourceManagerKey ContextKey = "resourceManager"
	// OperationIDKey is the context key for storing operation IDs
	OperationIDKey ContextKey = "operationID"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetResourceManagerFromContext retrieves a ResourceManager from context
func GetResourceManagerFromContext(ctx context.Context) (*ResourceManager, bool) {
	rm, ok := ctx.Value(ResourceManagerKey).(*ResourceManager)
	return rm, ok
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// WithOperationID adds an operation ID to the context
func WithOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, OperationIDKey, operationID)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetOperationIDFromContext retrieves an operation ID from context
func GetOperationIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(OperationIDKey).(string)
	return id, ok
}
