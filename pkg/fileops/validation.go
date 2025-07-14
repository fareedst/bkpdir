// Package fileops provides file operations and utilities for CLI applications.
//
// This file contains path validation functionality including security checks.
package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ARCH-001: See architecture.md - Core Architecture [DECISION:maintenance]

// PathValidator provides path validation and security checking functionality
type PathValidator struct{}

// Validator defines the interface for path validation operations
type Validator interface {
	ValidatePath(path string) error
	ValidateExistence(path string) error
	ValidateReadable(path string) error
	ValidateWritable(path string) error
	IsSecurePath(path string) bool
}

// NewPathValidator creates a new PathValidator instance
func NewPathValidator() Validator {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	return &PathValidator{}
}

// ValidatePath performs comprehensive path validation including security checks
func (pv *PathValidator) ValidatePath(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Security validation
	if !pv.IsSecurePath(path) {
		return fmt.Errorf("path contains unsafe elements: %s", path)
	}

	// Clean and resolve the path
	cleanPath := filepath.Clean(path)
	if cleanPath != path {
		// Allow this but could be flagged for security review
	}

	return nil
}

// ValidateExistence checks if a path exists
func (pv *PathValidator) ValidateExistence(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	if err := pv.ValidatePath(path); err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		return fmt.Errorf("cannot access path %s: %v", path, err)
	}

	return nil
}

// ValidateReadable checks if a path is readable
func (pv *PathValidator) ValidateReadable(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	if err := pv.ValidateExistence(path); err != nil {
		return err
	}

	// Try to open the file/directory for reading
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("path is not readable: %s (%v)", path, err)
	}
	file.Close()

	return nil
}

// ValidateWritable checks if a path is writable
func (pv *PathValidator) ValidateWritable(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	if err := pv.ValidatePath(path); err != nil {
		return err
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Check if parent directory is writable
			parentDir := filepath.Dir(path)
			return pv.ValidateWritable(parentDir)
		}
		return fmt.Errorf("cannot access path %s: %v", path, err)
	}

	// Check permissions
	if info.IsDir() {
		// For directories, try to create a temporary file
		tempFile := filepath.Join(path, ".tmp_write_test")
		file, err := os.Create(tempFile)
		if err != nil {
			return fmt.Errorf("directory is not writable: %s (%v)", path, err)
		}
		file.Close()
		os.Remove(tempFile)
	} else {
		// For files, try to open in write mode
		file, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("file is not writable: %s (%v)", path, err)
		}
		file.Close()
	}

	return nil
}

// IsSecurePath checks if a path is secure (no path traversal, etc.)
func (pv *PathValidator) IsSecurePath(path string) bool {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]

	// Check for path traversal attempts
	if strings.Contains(path, "..") {
		return false
	}

	// Check for absolute path attempts that might escape boundaries
	if filepath.IsAbs(path) {
		// Absolute paths are allowed but should be flagged for review
		// in security-sensitive contexts
	}

	// Check for suspicious patterns
	suspicious := []string{
		"~",    // Home directory references
		"$",    // Environment variable references
		"\x00", // Null bytes
		"\r",   // Carriage returns
		"\n",   // Line feeds
	}

	for _, pattern := range suspicious {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	return true
}

// Convenience functions for common validation tasks

// ValidatePath validates a path using the default validator
func ValidatePath(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	validator := NewPathValidator()
	return validator.ValidatePath(path)
}

// ValidateExistence validates path existence using the default validator
func ValidateExistence(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	validator := NewPathValidator()
	return validator.ValidateExistence(path)
}

// ValidateReadable validates path readability using the default validator
func ValidateReadable(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	validator := NewPathValidator()
	return validator.ValidateReadable(path)
}

// ValidateWritable validates path writability using the default validator
func ValidateWritable(path string) error {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	validator := NewPathValidator()
	return validator.ValidateWritable(path)
}

// IsSecurePath checks path security using the default validator
func IsSecurePath(path string) bool {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	validator := NewPathValidator()
	return validator.IsSecurePath(path)
}
