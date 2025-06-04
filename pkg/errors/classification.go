// Error classification utilities for detecting and categorizing different types of errors.
// This file contains the error detection functions extracted from the original
// errors.go file, generalized for reuse across CLI applications.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package errors

import (
	"os"
	"strings"
)

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Disk space error detection
// IsDiskFullError checks if an error indicates disk space exhaustion
// Extracted from original errors.go with comprehensive disk space error detection
func IsDiskFullError(err error) bool {
	if err == nil {
		return false
	}

	errorText := strings.ToLower(err.Error())
	patterns := []string{
		"no space left on device",
		"disk full",
		"not enough space",
		"insufficient disk space",
		"device full",
		"quota exceeded",
		"file too large",
		"no space available",
		"disk quota exceeded",
		"storage full",
		"volume full",
		"filesystem full",
	}

	for _, pattern := range patterns {
		if strings.Contains(errorText, pattern) {
			return true
		}
	}
	return false
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Permission error detection
// IsPermissionError checks if an error indicates permission or access issues
// Extracted from original errors.go with enhanced permission error detection
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}

	errorText := strings.ToLower(err.Error())
	patterns := []string{
		"permission denied",
		"access denied",
		"operation not permitted",
		"insufficient privileges",
		"access is denied",
		"unauthorized",
		"forbidden",
		"not allowed",
		"access violation",
		"insufficient permissions",
		"privilege not held",
		"access restricted",
	}

	for _, pattern := range patterns {
		if strings.Contains(errorText, pattern) {
			return true
		}
	}
	return false
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Directory/file existence error detection
// IsDirectoryNotFoundError checks if an error indicates a directory doesn't exist
// Extracted from original errors.go with enhanced path existence detection
func IsDirectoryNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errorText := strings.ToLower(err.Error())
	patterns := []string{
		"no such file or directory",
		"directory not found",
		"path not found",
		"file not found",
		"does not exist",
		"cannot find",
		"not found",
		"no such directory",
		"invalid path",
		"path does not exist",
		"directory does not exist",
	}

	for _, pattern := range patterns {
		if strings.Contains(errorText, pattern) {
			return true
		}
	}
	return false
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 File existence error detection
// IsFileNotFoundError checks if an error indicates a file doesn't exist
func IsFileNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	// Use os.IsNotExist for Go standard library errors
	if os.IsNotExist(err) {
		return true
	}

	errorText := strings.ToLower(err.Error())
	patterns := []string{
		"no such file",
		"file not found",
		"cannot find file",
		"file does not exist",
		"no such file or directory",
	}

	for _, pattern := range patterns {
		if strings.Contains(errorText, pattern) {
			return true
		}
	}
	return false
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Network error detection
// IsNetworkError checks if an error indicates network connectivity issues
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	errorText := strings.ToLower(err.Error())
	patterns := []string{
		"connection refused",
		"network unreachable",
		"timeout",
		"connection timeout",
		"connection reset",
		"no route to host",
		"network is down",
		"host is down",
		"connection failed",
		"dns lookup failed",
		"temporary failure in name resolution",
	}

	for _, pattern := range patterns {
		if strings.Contains(errorText, pattern) {
			return true
		}
	}
	return false
}

// ⭐ EXTRACT-002: Error classification framework - 🔧 Default classifier implementation
// DefaultErrorClassifier provides a default implementation of ErrorClassifier
type DefaultErrorClassifier struct{}

// NewDefaultErrorClassifier creates a new default error classifier
func NewDefaultErrorClassifier() *DefaultErrorClassifier {
	return &DefaultErrorClassifier{}
}

// ⭐ EXTRACT-002: Error classification framework - 🔍 Error categorization
// ClassifyError categorizes an error into one of the predefined categories
func (c *DefaultErrorClassifier) ClassifyError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}

	switch {
	case IsDiskFullError(err):
		return ErrorCategoryDiskSpace
	case IsPermissionError(err):
		return ErrorCategoryPermission
	case IsDirectoryNotFoundError(err), IsFileNotFoundError(err):
		return ErrorCategoryFilesystem
	case IsNetworkError(err):
		return ErrorCategoryNetwork
	default:
		return ErrorCategoryUnknown
	}
}

// ⭐ EXTRACT-002: Error classification framework - 🔍 Recoverability assessment
// IsRecoverable determines if an error can potentially be recovered from
func (c *DefaultErrorClassifier) IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	category := c.ClassifyError(err)
	switch category {
	case ErrorCategoryDiskSpace:
		// Disk space errors might be recoverable if space is freed
		return true
	case ErrorCategoryPermission:
		// Permission errors are typically not recoverable without intervention
		return false
	case ErrorCategoryFilesystem:
		// File/directory not found errors are typically not recoverable
		return false
	case ErrorCategoryNetwork:
		// Network errors might be recoverable with retry
		return true
	case ErrorCategoryConfiguration:
		// Configuration errors require fixing configuration
		return false
	case ErrorCategoryValidation:
		// Validation errors require fixing input
		return false
	default:
		// Unknown errors default to not recoverable
		return false
	}
}

// ⭐ EXTRACT-002: Error classification framework - 🔍 Severity assessment
// GetSeverity determines the severity level of an error
func (c *DefaultErrorClassifier) GetSeverity(err error) ErrorSeverity {
	if err == nil {
		return ErrorSeverityInfo
	}

	category := c.ClassifyError(err)
	switch category {
	case ErrorCategoryDiskSpace:
		// Disk space errors are critical as they prevent operations
		return ErrorSeverityCritical
	case ErrorCategoryPermission:
		// Permission errors are errors that prevent operations
		return ErrorSeverityError
	case ErrorCategoryFilesystem:
		// File/directory not found errors are errors
		return ErrorSeverityError
	case ErrorCategoryNetwork:
		// Network errors vary in severity, default to error
		return ErrorSeverityError
	case ErrorCategoryConfiguration:
		// Configuration errors are critical for application startup
		return ErrorSeverityCritical
	case ErrorCategoryValidation:
		// Validation errors are warnings/errors depending on context
		return ErrorSeverityWarning
	default:
		// Unknown errors default to error severity
		return ErrorSeverityError
	}
}

// ⭐ EXTRACT-002: Error classification utilities - 🔧 Pattern matching framework
// ErrorPattern represents a configurable error detection pattern
type ErrorPattern struct {
	Name        string   // Name of the error pattern
	Patterns    []string // Text patterns to match (case-insensitive)
	Category    ErrorCategory
	Severity    ErrorSeverity
	Recoverable bool
}

// ⭐ EXTRACT-002: Error classification utilities - 🔧 Configurable classifier
// ConfigurableErrorClassifier allows customizing error classification patterns
type ConfigurableErrorClassifier struct {
	patterns []ErrorPattern
	fallback ErrorClassifier
}

// NewConfigurableErrorClassifier creates a new configurable error classifier
func NewConfigurableErrorClassifier(patterns []ErrorPattern, fallback ErrorClassifier) *ConfigurableErrorClassifier {
	if fallback == nil {
		fallback = NewDefaultErrorClassifier()
	}
	return &ConfigurableErrorClassifier{
		patterns: patterns,
		fallback: fallback,
	}
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Pattern-based classification
// ClassifyError classifies errors using configured patterns, falling back to default
func (c *ConfigurableErrorClassifier) ClassifyError(err error) ErrorCategory {
	if err == nil {
		return ErrorCategoryUnknown
	}

	errorText := strings.ToLower(err.Error())

	// Check custom patterns first
	for _, pattern := range c.patterns {
		for _, textPattern := range pattern.Patterns {
			if strings.Contains(errorText, strings.ToLower(textPattern)) {
				return pattern.Category
			}
		}
	}

	// Fall back to default classifier
	return c.fallback.ClassifyError(err)
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Pattern-based recoverability
// IsRecoverable determines recoverability using configured patterns
func (c *ConfigurableErrorClassifier) IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	errorText := strings.ToLower(err.Error())

	// Check custom patterns first
	for _, pattern := range c.patterns {
		for _, textPattern := range pattern.Patterns {
			if strings.Contains(errorText, strings.ToLower(textPattern)) {
				return pattern.Recoverable
			}
		}
	}

	// Fall back to default classifier
	return c.fallback.IsRecoverable(err)
}

// ⭐ EXTRACT-002: Error classification utilities - 🔍 Pattern-based severity
// GetSeverity determines severity using configured patterns
func (c *ConfigurableErrorClassifier) GetSeverity(err error) ErrorSeverity {
	if err == nil {
		return ErrorSeverityInfo
	}

	errorText := strings.ToLower(err.Error())

	// Check custom patterns first
	for _, pattern := range c.patterns {
		for _, textPattern := range pattern.Patterns {
			if strings.Contains(errorText, strings.ToLower(textPattern)) {
				return pattern.Severity
			}
		}
	}

	// Fall back to default classifier
	return c.fallback.GetSeverity(err)
}
