// Package errors provides structured error handling and error classification utilities
// for CLI applications. It includes generic error types, error classification,
// and integration patterns for reliable error handling across applications.
//
// This package extracts and generalizes error handling patterns from the BkpDir
// application to provide reusable error handling infrastructure for Go CLI applications.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package errors

import (
	"context"
	"fmt"
	"os"
)

// ⭐ EXTRACT-002: Error interfaces and types - 🔍 Core error handling contracts
// ErrorInterface provides a common interface for all structured errors in CLI applications
type ErrorInterface interface {
	Error() string
	GetStatusCode() int
	GetOperation() string
	GetPath() string
	GetMessage() string
	Unwrap() error
}

// ⭐ EXTRACT-002: Error configuration interface - 🔧 Configuration abstraction
// ErrorConfig abstracts configuration dependencies for error handling
type ErrorConfig interface {
	GetStatusCodes() map[string]int
	GetErrorFormatStrings() map[string]string
	GetDirectoryPermissions() os.FileMode
	GetFilePermissions() os.FileMode
}

// ⭐ EXTRACT-002: Error formatter interface - 📝 Formatter abstraction
// ErrorFormatter abstracts formatter dependencies for error handling
type ErrorFormatter interface {
	FormatError(message string) string
	PrintError(message string)
	FormatDiskFullError(err error) string
	FormatPermissionError(err error) string
	FormatDirectoryNotFound(err error) string
	FormatFileNotFound(err error) string
	PrintDiskFullError(err error)
	PrintPermissionError(err error)
	PrintDirectoryNotFound(err error)
}

// ⭐ EXTRACT-002: Generic ApplicationError type - 🔧 Generalized error structure
// ApplicationError represents a structured error with status code and operation context
// This generalizes the ArchiveError/BackupError pattern for any CLI application
type ApplicationError struct {
	Message    string // Human-readable error message
	StatusCode int    // Exit status code for the application
	Operation  string // Operation context (e.g., "create", "verify", "list")
	Path       string // File or directory path involved in the error
	Err        error  // Underlying error for debugging and error chaining
}

// ⭐ EXTRACT-002: Error interface implementation - 📝 Standard error interface
func (e *ApplicationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ⭐ EXTRACT-002: Error unwrapping support - 🔍 Error chain access
func (e *ApplicationError) Unwrap() error {
	return e.Err
}

// ⭐ EXTRACT-002: Status code access - 🔍 Exit code retrieval
func (e *ApplicationError) GetStatusCode() int {
	return e.StatusCode
}

// ⭐ EXTRACT-002: Operation context access - 🔍 Operation tracking
func (e *ApplicationError) GetOperation() string {
	return e.Operation
}

// ⭐ EXTRACT-002: Path context access - 🔍 Path information
func (e *ApplicationError) GetPath() string {
	return e.Path
}

// ⭐ EXTRACT-002: Message access - 🔍 Error message retrieval
func (e *ApplicationError) GetMessage() string {
	return e.Message
}

// ⭐ EXTRACT-002: Error factory functions - 🔧 Error creation utilities

// NewApplicationError creates a new structured error with message and status code
func NewApplicationError(message string, statusCode int) *ApplicationError {
	return &ApplicationError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewApplicationErrorWithCause creates a new structured error with underlying cause
func NewApplicationErrorWithCause(message string, statusCode int, err error) *ApplicationError {
	return &ApplicationError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// NewApplicationErrorWithContext creates a new structured error with full context
func NewApplicationErrorWithContext(
	message string,
	statusCode int,
	operation, path string,
	err error,
) *ApplicationError {
	return &ApplicationError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// ⭐ EXTRACT-002: Error classification framework - 🔍 Error categorization

// ErrorClassifier defines an interface for classifying different types of errors
type ErrorClassifier interface {
	ClassifyError(err error) ErrorCategory
	IsRecoverable(err error) bool
	GetSeverity(err error) ErrorSeverity
}

// ErrorCategory represents different categories of errors
type ErrorCategory int

const (
	// ErrorCategoryUnknown represents unclassified errors
	ErrorCategoryUnknown ErrorCategory = iota
	// ErrorCategoryFilesystem represents file system related errors
	ErrorCategoryFilesystem
	// ErrorCategoryPermission represents permission and access errors
	ErrorCategoryPermission
	// ErrorCategoryDiskSpace represents disk space and storage errors
	ErrorCategoryDiskSpace
	// ErrorCategoryNetwork represents network and connectivity errors
	ErrorCategoryNetwork
	// ErrorCategoryConfiguration represents configuration and setup errors
	ErrorCategoryConfiguration
	// ErrorCategoryValidation represents input validation errors
	ErrorCategoryValidation
)

// ErrorSeverity represents the severity level of errors
type ErrorSeverity int

const (
	// ErrorSeverityInfo represents informational messages
	ErrorSeverityInfo ErrorSeverity = iota
	// ErrorSeverityWarning represents warning conditions
	ErrorSeverityWarning
	// ErrorSeverityError represents error conditions
	ErrorSeverityError
	// ErrorSeverityCritical represents critical system errors
	ErrorSeverityCritical
	// ErrorSeverityFatal represents fatal application errors
	ErrorSeverityFatal
)

// ⭐ EXTRACT-002: Error context abstraction - 🔧 Operation context
// ErrorContext provides context information for error handling operations
type ErrorContext struct {
	Operation string                 // The operation being performed
	Path      string                 // File or directory path involved
	Context   context.Context        // Cancellation context
	Metadata  map[string]interface{} // Additional context metadata
	StartTime int64                  // Operation start time for timeout tracking
}

// NewErrorContext creates a new error context with the specified operation and path
func NewErrorContext(operation, path string, ctx context.Context) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Path:      path,
		Context:   ctx,
		Metadata:  make(map[string]interface{}),
		StartTime: 0, // Will be set when operation starts
	}
}

// WithMetadata adds metadata to the error context
func (ec *ErrorContext) WithMetadata(key string, value interface{}) *ErrorContext {
	ec.Metadata[key] = value
	return ec
}

// GetMetadata retrieves metadata from the error context
func (ec *ErrorContext) GetMetadata(key string) (interface{}, bool) {
	value, exists := ec.Metadata[key]
	return value, exists
}
