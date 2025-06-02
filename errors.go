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
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

// REFACTOR-001: Error handling and resource management interface contracts defined
// REFACTOR-001: Config and formatter dependency interfaces required

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
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ArchiveError) Unwrap() error {
	return e.Err
}

// NewArchiveError creates a new structured error
func NewArchiveError(message string, statusCode int) *ArchiveError {
	// CFG-002: Basic structured error creation
	// DECISION-REF: DEC-004
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewArchiveErrorWithCause creates a new structured error with underlying cause
func NewArchiveErrorWithCause(message string, statusCode int, err error) *ArchiveError {
	// CFG-002: Error creation with cause chain
	// DECISION-REF: DEC-004
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
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// IsDiskFullError reports whether the error is due to disk full or quota exceeded conditions.
func IsDiskFullError(err error) bool {
	// Error classification for disk space issues
	// DECISION-REF: DEC-004
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	diskFullPatterns := []string{
		"no space left on device",
		"disk full",
		"not enough space",
		"insufficient disk space",
		"device full",
		"quota exceeded",
		"file too large",
	}

	for _, pattern := range diskFullPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// isDiskFullError provides enhanced disk space error detection with comprehensive pattern matching
func isDiskFullError(err error) bool {
	// Enhanced disk space error detection
	// DECISION-REF: DEC-004
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	diskFullPatterns := []string{
		"no space left",
		"disk full",
		"not enough space",
		"insufficient disk space",
		"device full",
		"quota exceeded",
		"file too large",
	}

	for _, pattern := range diskFullPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// IsPermissionError reports whether the error is due to permission or access denied conditions.
func IsPermissionError(err error) bool {
	// Error classification for permission issues
	// DECISION-REF: DEC-004
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	permissionPatterns := []string{
		"permission denied",
		"access denied",
		"operation not permitted",
		"insufficient privileges",
	}

	for _, pattern := range permissionPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// IsDirectoryNotFoundError reports whether the error is due to a missing directory.
func IsDirectoryNotFoundError(err error) bool {
	// Error classification for directory not found
	// DECISION-REF: DEC-004
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "no such file or directory") ||
		strings.Contains(errStr, "cannot find the path") ||
		strings.Contains(errStr, "directory not found")
}

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
	return os.RemoveAll(td.Path)
}

func (td *TempDir) String() string {
	return fmt.Sprintf("temp dir: %s", td.Path)
}

// ResourceManager manages automatic cleanup of resources
type ResourceManager struct {
	resources []Resource
	mutex     sync.RWMutex
}

// NewResourceManager creates a new ResourceManager for tracking resources.
func NewResourceManager() *ResourceManager {
	// DECISION-REF: DEC-006
	return &ResourceManager{
		resources: make([]Resource, 0),
	}
}

// AddResource adds a resource to the ResourceManager for cleanup.
func (rm *ResourceManager) AddResource(resource Resource) {
	// DECISION-REF: DEC-006
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = append(rm.resources, resource)
}

// AddTempFile adds a temporary file resource to the ResourceManager.
func (rm *ResourceManager) AddTempFile(path string) {
	// DECISION-REF: DEC-006
	rm.AddResource(&TempFile{Path: path})
}

// AddTempDir adds a temporary directory resource to the ResourceManager.
func (rm *ResourceManager) AddTempDir(path string) {
	// DECISION-REF: DEC-006
	rm.AddResource(&TempDir{Path: path})
}

// RemoveResource removes a resource from the ResourceManager.
func (rm *ResourceManager) RemoveResource(resource Resource) {
	// DECISION-REF: DEC-006
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
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during cleanup: %v", r)
		}
	}()

	return rm.Cleanup()
}

// ContextualOperation provides context and resource management for operations.
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

// NewContextualOperation creates a new ContextualOperation with the given context.
func NewContextualOperation(ctx context.Context) *ContextualOperation {
	// DECISION-REF: DEC-006, DEC-007
	return &ContextualOperation{
		ctx: ctx,
		rm:  NewResourceManager(),
	}
}

// Context returns the context associated with the ContextualOperation.
func (co *ContextualOperation) Context() context.Context {
	// DECISION-REF: DEC-007
	return co.ctx
}

// ResourceManager returns the ResourceManager associated with the ContextualOperation.
func (co *ContextualOperation) ResourceManager() *ResourceManager {
	// DECISION-REF: DEC-006
	return co.rm
}

// IsCancelled checks if the operation has been cancelled.
func (co *ContextualOperation) IsCancelled() bool {
	// DECISION-REF: DEC-007
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
	return co.ctx.Err()
}

// Cleanup cleans up all resources associated with the operation.
func (co *ContextualOperation) Cleanup() error {
	// DECISION-REF: DEC-006
	return co.rm.Cleanup()
}

// HandleArchiveError processes an archive error and returns the appropriate status code.
func HandleArchiveError(err error, cfg *Config, formatter *OutputFormatter) int {
	// CFG-002: Error handling with status code resolution
	// CFG-003: Error output formatting
	// CFG-004: Configurable error messages
	// DECISION-REF: DEC-004
	if err == nil {
		return 0
	}

	var archiveErr *ArchiveError
	if errors.As(err, &archiveErr) {
		formatter.PrintError(archiveErr.Error())
		return archiveErr.StatusCode
	}

	// Handle specific error types with configurable messages
	switch {
	case IsDiskFullError(err):
		formatter.PrintDiskFullError(err)
		return cfg.StatusDiskFull
	case IsPermissionError(err):
		formatter.PrintPermissionError(err)
		return cfg.StatusPermissionDenied
	case IsDirectoryNotFoundError(err):
		formatter.PrintDirectoryNotFound(err)
		return cfg.StatusDirectoryNotFound
	default:
		formatter.PrintError(err.Error())
		return 1
	}
}

// AtomicWriteFile writes data to a file atomically using a temporary file.
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	// DECISION-REF: DEC-006, DEC-008
	// CFG-004: Configurable error messages
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

// SafeMkdirAll creates a directory and all necessary parent directories.
func SafeMkdirAll(path string, perm os.FileMode, cfg *Config) error {
	// Directory creation with error handling
	// CFG-004: Configurable error messages
	// DECISION-REF: DEC-004
	if err := os.MkdirAll(path, perm); err != nil {
		if IsDiskFullError(err) {
			return NewArchiveErrorWithCause(
				"Failed to create directory: disk full",
				cfg.StatusDiskFull,
				err,
			)
		}
		return NewArchiveErrorWithCause(
			"Failed to create directory",
			cfg.StatusConfigError,
			err,
		)
	}
	return nil
}

// ValidateDirectoryPath checks if a path is a valid directory.
func ValidateDirectoryPath(path string, cfg *Config) error {
	// Directory path validation
	// CFG-004: Configurable error messages
	// DECISION-REF: DEC-004
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
