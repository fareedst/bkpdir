// [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING]
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
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with no preloaded exclusion patterns.
func NewPathValidator() Validator {
	return &PathValidator{}
}

// ValidatePath performs comprehensive path validation including security checks
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject empty paths, run IsSecurePath, clean path without failing on normalization drift.
func (pv *PathValidator) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !pv.IsSecurePath(path) {
		return fmt.Errorf("path contains unsafe elements: %s", path)
	}

	cleanPath := filepath.Clean(path)
	if cleanPath != path {
		// Allow this but could be flagged for security review
	}

	return nil
}

// ValidateExistence checks if a path exists
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then Stat path; return not-exist or access errors.
func (pv *PathValidator) ValidateExistence(path string) error {
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
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidateExistence then Open path for read; close on success.
func (pv *PathValidator) ValidateReadable(path string) error {
	if err := pv.ValidateExistence(path); err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("path is not readable: %s (%v)", path, err)
	}
	file.Close()

	return nil
}

// ValidateWritable checks if a path is writable
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then test write via temp file in directory or OpenFile WRONLY for files; recurse to parent when path missing.
func (pv *PathValidator) ValidateWritable(path string) error {
	if err := pv.ValidatePath(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			parentDir := filepath.Dir(path)
			return pv.ValidateWritable(parentDir)
		}
		return fmt.Errorf("cannot access path %s: %v", path, err)
	}

	if info.IsDir() {
		tempFile := filepath.Join(path, ".tmp_write_test")
		file, err := os.Create(tempFile)
		if err != nil {
			return fmt.Errorf("directory is not writable: %s (%v)", path, err)
		}
		file.Close()
		os.Remove(tempFile)
	} else {
		file, err := os.OpenFile(path, os.O_WRONLY, 0)
		if err != nil {
			return fmt.Errorf("file is not writable: %s (%v)", path, err)
		}
		file.Close()
	}

	return nil
}

// IsSecurePath checks if a path is secure (no path traversal, etc.)
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject .. traversal, null bytes, newlines, and shell-expansion characters in path strings.
func (pv *PathValidator) IsSecurePath(path string) bool {
	if strings.Contains(path, "..") {
		return false
	}

	if filepath.IsAbs(path) {
		// Absolute paths are allowed but should be flagged for review
		// in security-sensitive contexts
	}

	suspicious := []string{
		"~",
		"$",
		"\x00",
		"\r",
		"\n",
	}

	for _, pattern := range suspicious {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	return true
}

// ValidatePath validates a path using the default validator
func ValidatePath(path string) error {
	validator := NewPathValidator()
	return validator.ValidatePath(path)
}

// ValidateExistence validates path existence using the default validator
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then Stat path; return not-exist or access errors.
func ValidateExistence(path string) error {
	validator := NewPathValidator()
	return validator.ValidateExistence(path)
}

// ValidateReadable validates path readability using the default validator
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidateExistence then Open path for read; close on success.
func ValidateReadable(path string) error {
	validator := NewPathValidator()
	return validator.ValidateReadable(path)
}

// ValidateWritable validates path writability using the default validator
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then test write via temp file in directory or OpenFile WRONLY for files; recurse to parent when path missing.
func ValidateWritable(path string) error {
	validator := NewPathValidator()
	return validator.ValidateWritable(path)
}

// IsSecurePath checks path security using the default validator
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject .. traversal, null bytes, newlines, and shell-expansion characters in path strings.
func IsSecurePath(path string) bool {
	validator := NewPathValidator()
	return validator.IsSecurePath(path)
}
