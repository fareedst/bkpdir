// This file is part of bkpdir
//
// Package main provides error handling and resource management for BkpDir.
// It includes error types, resource cleanup, and context-aware operations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// REFACTOR-001: Error handling and resource management interface contracts defined
// REFACTOR-001: Config and formatter dependency interfaces required
// REFACTOR-004: Error standardization - Common interface for all structured errors
// REFACTOR-005: Structure optimization - Configuration interface abstraction
// REFACTOR-005: Structure optimization - Formatter interface abstraction
// REFACTOR-005: Structure optimization - Resource management interface

// ErrorConfig abstracts configuration dependencies for error handling
type ErrorConfig interface {
	GetStatusCodes() map[string]int
	GetErrorFormatStrings() map[string]string
	GetDirectoryPermissions() os.FileMode
	GetFilePermissions() os.FileMode
}

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

// ResourceManagerInterface defines clean resource management contract
type ResourceManagerInterface interface {
	AddResource(resource Resource)
	AddTempFile(path string)
	AddTempDir(path string)
	RemoveResource(resource Resource)
	Cleanup() error
	CleanupWithPanicRecovery() error
}

// ErrorInterface provides a common interface for all structured errors in the application
type ErrorInterface interface {
	Error() string
	GetStatusCode() int
	GetOperation() string
	GetPath() string
	GetMessage() string
	Unwrap() error
}

// ArchiveError represents a structured error with status code
type ArchiveError struct {
	Message    string
	StatusCode int
	Operation  string
	Path       string
	Err        error
}

func (e *ArchiveError) Error() string {
	// CFG-002: Structured error message formatting
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized error message format
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ArchiveError) Unwrap() error {
	return e.Err
}

// REFACTOR-004: Standardized error interface implementation
func (e *ArchiveError) GetStatusCode() int {
	return e.StatusCode
}

func (e *ArchiveError) GetOperation() string {
	return e.Operation
}

func (e *ArchiveError) GetPath() string {
	return e.Path
}

func (e *ArchiveError) GetMessage() string {
	return e.Message
}

// BackupError represents a structured error with status code for backup operations
// REFACTOR-004: Unified error pattern with ArchiveError
type BackupError struct {
	Message    string
	StatusCode int
	Operation  string
	Path       string
	Err        error
}

func (e *BackupError) Error() string {
	// REFACTOR-004: Consistent error message formatting across all error types
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BackupError) Unwrap() error {
	return e.Err
}

// REFACTOR-004: Standardized error interface implementation
func (e *BackupError) GetStatusCode() int {
	return e.StatusCode
}

func (e *BackupError) GetOperation() string {
	return e.Operation
}

func (e *BackupError) GetPath() string {
	return e.Path
}

func (e *BackupError) GetMessage() string {
	return e.Message
}

// REFACTOR-004: Standardized error constructors with consistent patterns

// NewArchiveError creates a new structured error
func NewArchiveError(message string, statusCode int) *ArchiveError {
	// CFG-002: Basic structured error creation
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized constructor pattern
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewArchiveErrorWithCause creates a new structured error with underlying cause
func NewArchiveErrorWithCause(message string, statusCode int, err error) *ArchiveError {
	// CFG-002: Error creation with cause chain
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized constructor with cause pattern
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// NewArchiveErrorWithContext creates a new structured error with operation context
func NewArchiveErrorWithContext(
	message string,
	statusCode int,
	operation, path string,
	err error,
) *ArchiveError {
	// CFG-002: Error creation with operation context
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized constructor with full context pattern
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// NewBackupError creates a new structured backup error
// REFACTOR-004: Standardized constructor pattern matching ArchiveError
func NewBackupError(message string, statusCode int, operation, path string) *BackupError {
	return &BackupError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
	}
}

// NewBackupErrorWithCause creates a new structured backup error with underlying cause
// REFACTOR-004: Standardized constructor with cause pattern matching ArchiveError
func NewBackupErrorWithCause(message string, statusCode int, operation, path string, err error) *BackupError {
	return &BackupError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// REFACTOR-004: Enhanced error classification with standardized patterns

// IsDiskFullError checks if an error indicates disk space issues
func IsDiskFullError(err error) bool {
	// Enhanced disk space error detection
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized error classification pattern
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	diskSpaceIndicators := []string{
		"no space left",
		"disk full",
		"not enough space",
		"insufficient disk space",
		"device full",
		"quota exceeded",
		"file too large",
		"no space available",
		"disk quota exceeded",
		"space limit exceeded",
	}

	for _, indicator := range diskSpaceIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	return false
}

// IsPermissionError checks if an error indicates permission issues
func IsPermissionError(err error) bool {
	// Permission error classification
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized error classification pattern
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	permissionIndicators := []string{
		"permission denied",
		"access is denied",
		"operation not permitted",
		"insufficient privileges",
		"unauthorized",
		"access denied",
		"forbidden",
	}

	for _, indicator := range permissionIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	return false
}

func IsDirectoryNotFoundError(err error) bool {
	// Error classification for directory not found
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized error classification pattern
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	notFoundIndicators := []string{
		"no such file or directory",
		"cannot find the path",
		"directory not found",
		"path does not exist",
		"file not found",
	}

	for _, indicator := range notFoundIndicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	return false
}

// REFACTOR-004: Standardized resource management patterns

// Resource represents a resource that needs cleanup
type Resource interface {
	Cleanup() error
	String() string
}

// TempFile represents a temporary file resource
type TempFile struct {
	Path string
}

// Cleanup removes the temporary file from the filesystem.
func (tf *TempFile) Cleanup() error {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized cleanup pattern
	return os.Remove(tf.Path)
}

func (tf *TempFile) String() string {
	return fmt.Sprintf("temp file: %s", tf.Path)
}

// TempDir represents a temporary directory resource
type TempDir struct {
	Path string
}

// Cleanup removes the temporary directory and its contents from the filesystem.
func (td *TempDir) Cleanup() error {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized cleanup pattern
	return os.RemoveAll(td.Path)
}

func (td *TempDir) String() string {
	return fmt.Sprintf("temp dir: %s", td.Path)
}

// ResourceManager manages automatic cleanup of resources
// REFACTOR-004: Standardized resource management implementation
type ResourceManager struct {
	resources []Resource
	mutex     sync.RWMutex
}

// NewResourceManager creates a new ResourceManager for tracking resources.
func NewResourceManager() *ResourceManager {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized resource manager constructor
	return &ResourceManager{
		resources: make([]Resource, 0),
	}
}

// AddResource adds a resource to the ResourceManager for cleanup.
func (rm *ResourceManager) AddResource(resource Resource) {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Thread-safe resource management pattern
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = append(rm.resources, resource)
}

// AddTempFile adds a temporary file resource to the ResourceManager.
func (rm *ResourceManager) AddTempFile(path string) {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized temp file pattern
	rm.AddResource(&TempFile{Path: path})
}

// AddTempDir adds a temporary directory resource to the ResourceManager.
func (rm *ResourceManager) AddTempDir(path string) {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized temp directory pattern
	rm.AddResource(&TempDir{Path: path})
}

// RemoveResource removes a resource from the ResourceManager.
func (rm *ResourceManager) RemoveResource(resource Resource) {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Thread-safe resource removal pattern
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for i, r := range rm.resources {
		if r.String() == resource.String() {
			rm.resources = append(rm.resources[:i], rm.resources[i+1:]...)
			break
		}
	}
}

// Cleanup cleans up all tracked resources in the ResourceManager.
func (rm *ResourceManager) Cleanup() error {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized cleanup error handling pattern
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	var lastError error
	for _, resource := range rm.resources {
		if err := resource.Cleanup(); err != nil {
			lastError = err
			// Continue cleanup even if individual operations fail
		}
	}

	rm.resources = rm.resources[:0] // Clear the slice
	return lastError
}

// CleanupWithPanicRecovery cleans up all resources and recovers from panics during cleanup.
func (rm *ResourceManager) CleanupWithPanicRecovery() (err error) {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized panic recovery pattern
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during cleanup: %v", r)
		}
	}()

	return rm.Cleanup()
}

// REFACTOR-004: Standardized context propagation patterns

// ContextualOperation provides context and resource management for operations.
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

// NewContextualOperation creates a new ContextualOperation with the given context.
func NewContextualOperation(ctx context.Context) *ContextualOperation {
	// DECISION-REF: DEC-006, DEC-007
	// REFACTOR-004: Standardized contextual operation pattern
	return &ContextualOperation{
		ctx: ctx,
		rm:  NewResourceManager(),
	}
}

// Context returns the context associated with the ContextualOperation.
func (co *ContextualOperation) Context() context.Context {
	// DECISION-REF: DEC-007
	// REFACTOR-004: Standardized context access pattern
	return co.ctx
}

// ResourceManager returns the ResourceManager associated with the ContextualOperation.
func (co *ContextualOperation) ResourceManager() *ResourceManager {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized resource manager access pattern
	return co.rm
}

// IsCancelled checks if the operation has been cancelled.
func (co *ContextualOperation) IsCancelled() bool {
	// DECISION-REF: DEC-007
	// REFACTOR-004: Standardized cancellation check pattern
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

// CheckCancellation checks if the operation has been cancelled and returns an error if it has.
func (co *ContextualOperation) CheckCancellation() error {
	// DECISION-REF: DEC-007
	// REFACTOR-004: Standardized cancellation error pattern
	return co.ctx.Err()
}

// Cleanup cleans up all resources associated with the operation.
func (co *ContextualOperation) Cleanup() error {
	// DECISION-REF: DEC-006
	// REFACTOR-004: Standardized contextual cleanup pattern
	return co.rm.Cleanup()
}

// REFACTOR-005: Structure optimization - Interface-based error handling
// HandleError provides centralized error handling with interface abstractions
func HandleError(err error, cfg ErrorConfig, formatter ErrorFormatter) int {
	// REFACTOR-004: Error standardization
	if err == nil {
		return 0
	}

	if archiveErr, ok := err.(*ArchiveError); ok {
		return HandleArchiveErrorWithInterface(archiveErr, cfg, formatter)
	}
	if backupErr, ok := err.(*BackupError); ok {
		return HandleBackupErrorWithInterface(backupErr, cfg, formatter)
	}

	// Handle specific error types with configurable messages
	statusCodes := cfg.GetStatusCodes()
	switch {
	case IsDiskFullError(err):
		formatter.PrintDiskFullError(err)
		if code, ok := statusCodes["disk_full"]; ok {
			return code
		}
		return 1
	case IsPermissionError(err):
		formatter.PrintPermissionError(err)
		if code, ok := statusCodes["permission_denied"]; ok {
			return code
		}
		return 1
	case IsDirectoryNotFoundError(err):
		formatter.PrintDirectoryNotFound(err)
		if code, ok := statusCodes["directory_not_found"]; ok {
			return code
		}
		return 1
	default:
		// Handle generic errors
		formatter.PrintError(err.Error())
		return 1
	}
}

// REFACTOR-005: Structure optimization - Interface-based archive error handling
// HandleArchiveErrorWithInterface handles archive errors using interface abstractions
func HandleArchiveErrorWithInterface(err *ArchiveError, cfg ErrorConfig, formatter ErrorFormatter) int {
	formatter.PrintError(err.Error())
	return err.GetStatusCode()
}

// REFACTOR-005: Structure optimization - Interface-based backup error handling
// HandleBackupErrorWithInterface handles backup errors using interface abstractions
func HandleBackupErrorWithInterface(err *BackupError, cfg ErrorConfig, formatter ErrorFormatter) int {
	formatter.PrintError(err.Error())
	return err.GetStatusCode()
}

// REFACTOR-005: Extraction preparation - Backward compatibility layer
// HandleArchiveError provides backward compatibility with concrete types
func HandleArchiveError(err error, cfg *Config, formatter *OutputFormatter) int {
	return HandleError(err, cfg, formatter)
}

// REFACTOR-005: Extraction preparation - Backward compatibility layer
// HandleBackupError provides backward compatibility with concrete types
func HandleBackupError(err error, cfg *Config, formatter *OutputFormatter) int {
	return HandleError(err, cfg, formatter)
}

// REFACTOR-004: Standardized atomic operation patterns

// AtomicWriteFile writes data to a file atomically using a temporary file.
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	// DECISION-REF: DEC-006, DEC-008
	// CFG-004: Configurable error messages
	// REFACTOR-004: Standardized atomic operation pattern
	tempFile := path + ".tmp"
	rm.AddTempFile(tempFile)

	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to write temporary file",
			1,
			err,
		)
	}

	if err := os.Rename(tempFile, path); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to finalize file",
			1,
			err,
		)
	}

	rm.RemoveResource(&TempFile{Path: tempFile})
	return nil
}

// AtomicWriteFileWithContext writes data to a file atomically with context support
// REFACTOR-004: Context-aware atomic operation pattern
func AtomicWriteFileWithContext(ctx context.Context, path string, data []byte, rm *ResourceManager) error {
	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	tempFile := path + ".tmp"
	rm.AddTempFile(tempFile)

	// Check for cancellation before write
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to write temporary file",
			1,
			err,
		)
	}

	// Check for cancellation before rename
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.Rename(tempFile, path); err != nil {
		return NewArchiveErrorWithCause(
			"Failed to finalize file",
			1,
			err,
		)
	}

	rm.RemoveResource(&TempFile{Path: tempFile})
	return nil
}

// REFACTOR-005: Structure optimization - Interface-based directory operations
// SafeMkdirAllWithInterface creates directories using interface abstractions
func SafeMkdirAllWithInterface(path string, perm os.FileMode, cfg ErrorConfig) error {
	// REFACTOR-004: Resource consolidation
	if err := os.MkdirAll(path, perm); err != nil {
		if IsDiskFullError(err) {
			return NewArchiveError(
				fmt.Sprintf("Failed to create directory due to insufficient disk space: %s", path),
				cfg.GetStatusCodes()["disk_full"],
			)
		}
		if IsPermissionError(err) {
			return NewArchiveError(
				fmt.Sprintf("Failed to create directory due to permission denied: %s", path),
				cfg.GetStatusCodes()["permission_denied"],
			)
		}
		return NewArchiveError(
			fmt.Sprintf("Failed to create directory: %s", path),
			1,
		)
	}
	return nil
}

// REFACTOR-005: Structure optimization - Interface-based context-aware directory operations
// SafeMkdirAllWithContextAndInterface creates directories with context using interface abstractions
func SafeMkdirAllWithContextAndInterface(ctx context.Context, path string, perm os.FileMode, cfg ErrorConfig) error {
	// REFACTOR-004: Context propagation
	if err := ctx.Err(); err != nil {
		return NewArchiveError(
			fmt.Sprintf("Context cancelled before creating directory: %s", path),
			1,
		)
	}

	return SafeMkdirAllWithInterface(path, perm, cfg)
}

// REFACTOR-005: Extraction preparation - Backward compatibility layer
// SafeMkdirAll provides backward compatibility with concrete types
func SafeMkdirAll(path string, perm os.FileMode, cfg *Config) error {
	return SafeMkdirAllWithInterface(path, perm, cfg)
}

// REFACTOR-005: Extraction preparation - Backward compatibility layer
// SafeMkdirAllWithContext provides backward compatibility with concrete types
func SafeMkdirAllWithContext(ctx context.Context, path string, perm os.FileMode, cfg *Config) error {
	return SafeMkdirAllWithContextAndInterface(ctx, path, perm, cfg)
}

// ValidateDirectoryPath checks if a path is a valid directory.
func ValidateDirectoryPath(path string, cfg *Config) error {
	// Directory path validation
	// CFG-004: Configurable error messages
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized validation pattern
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewArchiveError(
				"Directory not found",
				cfg.StatusDirectoryNotFound,
			)
		}
		if IsPermissionError(err) {
			return NewArchiveErrorWithCause(
				"Permission denied",
				cfg.StatusPermissionDenied,
				err,
			)
		}
		return NewArchiveErrorWithCause(
			"Failed to access directory",
			cfg.StatusConfigError,
			err,
		)
	}

	if !info.IsDir() {
		return NewArchiveError(
			"Path is not a directory",
			cfg.StatusInvalidDirectoryType,
		)
	}

	return nil
}

// ValidateFilePath checks if a path is a valid file.
func ValidateFilePath(path string, cfg *Config) error {
	// File path validation
	// CFG-004: Configurable error messages
	// DECISION-REF: DEC-004
	// REFACTOR-004: Standardized validation pattern
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewArchiveError(
				"File not found",
				cfg.StatusFileNotFound,
			)
		}
		if IsPermissionError(err) {
			return NewArchiveErrorWithCause(
				"Permission denied",
				cfg.StatusPermissionDenied,
				err,
			)
		}
		return NewArchiveErrorWithCause(
			"Failed to access file",
			cfg.StatusConfigError,
			err,
		)
	}

	if !info.Mode().IsRegular() {
		return NewArchiveError(
			"Path is not a regular file",
			cfg.StatusInvalidFileType,
		)
	}

	return nil
}

// REFACTOR-004: Context propagation utilities

// WithResourceManager creates a new context with an embedded resource manager
func WithResourceManager(ctx context.Context) (context.Context, *ResourceManager) {
	// REFACTOR-004: Standardized context enhancement pattern
	rm := NewResourceManager()
	return ctx, rm
}

// CheckContextAndCleanup checks context cancellation and performs cleanup if needed
func CheckContextAndCleanup(ctx context.Context, rm *ResourceManager) error {
	// REFACTOR-004: Standardized context check with cleanup pattern
	if err := ctx.Err(); err != nil {
		// Cleanup resources on cancellation
		if cleanupErr := rm.CleanupWithPanicRecovery(); cleanupErr != nil {
			// Return the original context error, but log cleanup error if needed
			return fmt.Errorf("context cancelled: %w (cleanup also failed: %v)", err, cleanupErr)
		}
		return err
	}
	return nil
}
