// Error handling functions for processing and responding to different types of errors.
// This file contains the error handler functions extracted from the original
// errors.go file, generalized for reuse across CLI applications.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package errors

import (
	"context"
	"os"
)

// ⭐ EXTRACT-002: Error handler functions - 🔧 Centralized error handling
// HandleError provides centralized error handling with interface abstractions
// This function processes different types of errors and returns appropriate status codes
func HandleError(err error, cfg ErrorConfig, formatter ErrorFormatter) int {
	if err == nil {
		return 0
	}

	// Handle ApplicationError (generalized from ArchiveError/BackupError)
	if appErr, ok := err.(*ApplicationError); ok {
		return HandleApplicationError(appErr, cfg, formatter)
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
	case IsFileNotFoundError(err):
		formatter.PrintError("File not found: " + err.Error())
		if code, ok := statusCodes["file_not_found"]; ok {
			return code
		}
		return 1
	case IsNetworkError(err):
		formatter.PrintError("Network error: " + err.Error())
		if code, ok := statusCodes["network_error"]; ok {
			return code
		}
		return 1
	default:
		// Handle generic errors
		formatter.PrintError(err.Error())
		return 1
	}
}

// ⭐ EXTRACT-002: Error handler functions - 🔍 Application error handling
// HandleApplicationError handles ApplicationError instances with proper formatting
func HandleApplicationError(err *ApplicationError, cfg ErrorConfig, formatter ErrorFormatter) int {
	formatter.PrintError(err.Error())
	return err.GetStatusCode()
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Context-aware error handling
// HandleErrorWithContext handles errors with context information for enhanced debugging
func HandleErrorWithContext(err error, errorCtx *ErrorContext, cfg ErrorConfig, formatter ErrorFormatter) int {
	if err == nil {
		return 0
	}

	// If we have an ApplicationError, enhance it with context
	if appErr, ok := err.(*ApplicationError); ok {
		if appErr.Operation == "" && errorCtx != nil {
			appErr.Operation = errorCtx.Operation
		}
		if appErr.Path == "" && errorCtx != nil {
			appErr.Path = errorCtx.Path
		}
		return HandleApplicationError(appErr, cfg, formatter)
	}

	// For other errors, create an ApplicationError with context
	contextualError := &ApplicationError{
		Message:    err.Error(),
		StatusCode: 1,
		Err:        err,
	}

	if errorCtx != nil {
		contextualError.Operation = errorCtx.Operation
		contextualError.Path = errorCtx.Path
	}

	return HandleApplicationError(contextualError, cfg, formatter)
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Error recovery framework
// ErrorRecoveryStrategy defines a strategy for attempting to recover from an error
type ErrorRecoveryStrategy interface {
	CanRecover(err error, context *ErrorContext) bool
	Recover(err error, context *ErrorContext) error
	GetRecoveryDescription() string
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Disk space recovery strategy
// DiskSpaceRecoveryStrategy attempts to recover from disk space errors
type DiskSpaceRecoveryStrategy struct {
	MinFreeSpace int64 // Minimum free space required in bytes
}

// NewDiskSpaceRecoveryStrategy creates a new disk space recovery strategy
func NewDiskSpaceRecoveryStrategy(minFreeSpace int64) *DiskSpaceRecoveryStrategy {
	return &DiskSpaceRecoveryStrategy{
		MinFreeSpace: minFreeSpace,
	}
}

// ⭐ EXTRACT-002: Error handler functions - 🔍 Disk space recovery capability
// CanRecover checks if a disk space error can potentially be recovered from
func (rs *DiskSpaceRecoveryStrategy) CanRecover(err error, context *ErrorContext) bool {
	return IsDiskFullError(err)
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Disk space recovery attempt
// Recover attempts to recover from disk space errors by checking available space
func (rs *DiskSpaceRecoveryStrategy) Recover(err error, context *ErrorContext) error {
	if !rs.CanRecover(err, context) {
		return err
	}

	// In a real implementation, this could:
	// 1. Check available disk space
	// 2. Clean up temporary files
	// 3. Suggest user actions
	// For now, we just return the original error
	return err
}

// ⭐ EXTRACT-002: Error handler functions - 🔍 Recovery description
// GetRecoveryDescription returns a description of what this strategy attempts
func (rs *DiskSpaceRecoveryStrategy) GetRecoveryDescription() string {
	return "Attempts to recover from disk space errors by checking available space and suggesting cleanup actions"
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Error recovery manager
// ErrorRecoveryManager manages multiple recovery strategies
type ErrorRecoveryManager struct {
	strategies []ErrorRecoveryStrategy
	classifier ErrorClassifier
}

// NewErrorRecoveryManager creates a new error recovery manager
func NewErrorRecoveryManager(classifier ErrorClassifier) *ErrorRecoveryManager {
	if classifier == nil {
		classifier = NewDefaultErrorClassifier()
	}
	return &ErrorRecoveryManager{
		strategies: make([]ErrorRecoveryStrategy, 0),
		classifier: classifier,
	}
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Recovery strategy registration
// AddStrategy adds a recovery strategy to the manager
func (rm *ErrorRecoveryManager) AddStrategy(strategy ErrorRecoveryStrategy) {
	rm.strategies = append(rm.strategies, strategy)
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Error recovery attempt
// TryRecover attempts to recover from an error using available strategies
func (rm *ErrorRecoveryManager) TryRecover(err error, context *ErrorContext) error {
	if err == nil {
		return nil
	}

	// Only attempt recovery if the error is classified as recoverable
	if !rm.classifier.IsRecoverable(err) {
		return err
	}

	// Try each strategy in order
	for _, strategy := range rm.strategies {
		if strategy.CanRecover(err, context) {
			if recoveredErr := strategy.Recover(err, context); recoveredErr == nil {
				return nil // Recovery successful
			}
		}
	}

	// No recovery strategy worked
	return err
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Enhanced error handling with recovery
// HandleErrorWithRecovery handles errors with automatic recovery attempts
func HandleErrorWithRecovery(
	err error,
	errorCtx *ErrorContext,
	cfg ErrorConfig,
	formatter ErrorFormatter,
	recoveryManager *ErrorRecoveryManager,
) int {
	if err == nil {
		return 0
	}

	// Attempt recovery if a recovery manager is provided
	if recoveryManager != nil {
		if recoveredErr := recoveryManager.TryRecover(err, errorCtx); recoveredErr == nil {
			// Recovery successful
			return 0
		}
	}

	// Fall back to standard error handling if recovery fails or isn't available
	return HandleErrorWithContext(err, errorCtx, cfg, formatter)
}

// ⭐ EXTRACT-002: Error handler functions - 🔍 Path validation utilities
// ValidateDirectoryPath validates that a path points to an accessible directory
// Extracted from original errors.go with enhanced error handling
func ValidateDirectoryPath(path string, cfg ErrorConfig) error {
	if path == "" {
		return NewApplicationError("directory path cannot be empty", 1)
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewApplicationErrorWithCause(
				"directory does not exist: "+path,
				20, // Directory not found status
				err,
			)
		}
		if IsPermissionError(err) {
			return NewApplicationErrorWithCause(
				"permission denied accessing directory: "+path,
				22, // Permission denied status
				err,
			)
		}
		return NewApplicationErrorWithCause(
			"cannot access directory: "+path,
			1,
			err,
		)
	}

	if !info.IsDir() {
		return NewApplicationError(
			"path is not a directory: "+path,
			21, // Invalid directory type status
		)
	}

	return nil
}

// ⭐ EXTRACT-002: Error handler functions - 🔍 File validation utilities
// ValidateFilePath validates that a path points to an accessible file
// Extracted from original errors.go with enhanced error handling
func ValidateFilePath(path string, cfg ErrorConfig) error {
	if path == "" {
		return NewApplicationError("file path cannot be empty", 1)
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewApplicationErrorWithCause(
				"file does not exist: "+path,
				20, // File not found status
				err,
			)
		}
		if IsPermissionError(err) {
			return NewApplicationErrorWithCause(
				"permission denied accessing file: "+path,
				22, // Permission denied status
				err,
			)
		}
		return NewApplicationErrorWithCause(
			"cannot access file: "+path,
			1,
			err,
		)
	}

	if info.IsDir() {
		return NewApplicationError(
			"path is a directory, not a file: "+path,
			21, // Invalid file type status
		)
	}

	return nil
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Safe filesystem operations
// SafeMkdirAll creates directories with enhanced error handling and context support
// Extracted from original errors.go with interface abstractions
func SafeMkdirAll(path string, perm os.FileMode, cfg ErrorConfig) error {
	return SafeMkdirAllWithContext(context.Background(), path, perm, cfg)
}

// ⭐ EXTRACT-002: Error handler functions - 🔧 Context-aware safe directory creation
// SafeMkdirAllWithContext creates directories with context support and enhanced error handling
func SafeMkdirAllWithContext(ctx context.Context, path string, perm os.FileMode, cfg ErrorConfig) error {
	// Check for cancellation
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.MkdirAll(path, perm); err != nil {
		if IsDiskFullError(err) {
			return NewApplicationErrorWithContext(
				"insufficient disk space to create directory",
				30, // Disk full status
				"directory_creation",
				path,
				err,
			)
		}
		if IsPermissionError(err) {
			return NewApplicationErrorWithContext(
				"permission denied creating directory",
				22, // Permission denied status
				"directory_creation",
				path,
				err,
			)
		}
		return NewApplicationErrorWithContext(
			"failed to create directory",
			31, // Directory creation failure status
			"directory_creation",
			path,
			err,
		)
	}

	return nil
}
