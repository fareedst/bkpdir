package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// ArchiveError represents a structured error with status code
type ArchiveError struct {
	Message    string
	StatusCode int
	Operation  string
	Path       string
	Err        error
}

func (e *ArchiveError) Error() string {
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
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewArchiveErrorWithCause creates a new structured error with underlying cause
func NewArchiveErrorWithCause(message string, statusCode int, err error) *ArchiveError {
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// NewArchiveErrorWithContext creates a new structured error with operation context
func NewArchiveErrorWithContext(message string, statusCode int, operation, path string, err error) *ArchiveError {
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// Enhanced error detection functions
func IsDiskFullError(err error) bool {
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

func IsPermissionError(err error) bool {
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

func IsDirectoryNotFoundError(err error) bool {
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

func (tf *TempFile) Cleanup() error {
	return os.Remove(tf.Path)
}

func (tf *TempFile) String() string {
	return fmt.Sprintf("temp file: %s", tf.Path)
}

// TempDir represents a temporary directory resource
type TempDir struct {
	Path string
}

func (td *TempDir) Cleanup() error {
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

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make([]Resource, 0),
	}
}

func (rm *ResourceManager) AddResource(resource Resource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = append(rm.resources, resource)
}

func (rm *ResourceManager) AddTempFile(path string) {
	rm.AddResource(&TempFile{Path: path})
}

func (rm *ResourceManager) AddTempDir(path string) {
	rm.AddResource(&TempDir{Path: path})
}

func (rm *ResourceManager) RemoveResource(resource Resource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for i, r := range rm.resources {
		if r.String() == resource.String() {
			rm.resources = append(rm.resources[:i], rm.resources[i+1:]...)
			break
		}
	}
}

func (rm *ResourceManager) Cleanup() error {
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

func (rm *ResourceManager) CleanupWithPanicRecovery() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during cleanup: %v", r)
		}
	}()

	return rm.Cleanup()
}

// Context-aware operations
type ContextualOperation struct {
	ctx context.Context
	rm  *ResourceManager
}

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

func (co *ContextualOperation) IsCancelled() bool {
	select {
	case <-co.ctx.Done():
		return true
	default:
		return false
	}
}

func (co *ContextualOperation) CheckCancellation() error {
	if co.IsCancelled() {
		return co.ctx.Err()
	}
	return nil
}

func (co *ContextualOperation) Cleanup() error {
	return co.rm.CleanupWithPanicRecovery()
}

// Error handling utilities
func HandleArchiveError(err error, cfg *Config, formatter *OutputFormatter) int {
	if err == nil {
		return 0
	}

	// Check if it's already an ArchiveError
	if archiveErr, ok := err.(*ArchiveError); ok {
		formatter.PrintError(archiveErr.Message)
		return archiveErr.StatusCode
	}

	// Detect error type and create appropriate ArchiveError
	var statusCode int
	var message string

	if IsDiskFullError(err) {
		statusCode = cfg.StatusDiskFull
		message = "Insufficient disk space"
	} else if IsPermissionError(err) {
		statusCode = cfg.StatusPermissionDenied
		message = "Permission denied"
	} else if IsDirectoryNotFoundError(err) {
		statusCode = cfg.StatusDirectoryNotFound
		message = "Directory not found"
	} else {
		statusCode = 1 // Generic error
		message = err.Error()
	}

	formatter.PrintError(message)
	return statusCode
}

// Atomic file operations
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	tempPath := path + ".tmp"
	rm.AddTempFile(tempPath)

	// Write to temporary file
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return err
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}

	// Remove from cleanup list since operation succeeded
	rm.RemoveResource(&TempFile{Path: tempPath})
	return nil
}

// Safe directory creation
func SafeMkdirAll(path string, perm os.FileMode, cfg *Config) error {
	if err := os.MkdirAll(path, perm); err != nil {
		if IsDiskFullError(err) {
			return NewArchiveError("Failed to create archive directory: insufficient disk space", cfg.StatusDiskFull)
		}
		if IsPermissionError(err) {
			return NewArchiveError("Failed to create archive directory: permission denied", cfg.StatusPermissionDenied)
		}
		return NewArchiveErrorWithCause("Failed to create archive directory", cfg.StatusFailedToCreateArchiveDirectory, err)
	}
	return nil
}

// Validate directory path
func ValidateDirectoryPath(path string, cfg *Config) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewArchiveError("Directory does not exist", cfg.StatusDirectoryNotFound)
		}
		if IsPermissionError(err) {
			return NewArchiveError("Permission denied accessing directory", cfg.StatusPermissionDenied)
		}
		return NewArchiveErrorWithCause("Failed to access directory", cfg.StatusDirectoryNotFound, err)
	}

	if !info.IsDir() {
		return NewArchiveError("Path is not a directory", cfg.StatusInvalidDirectoryType)
	}

	return nil
}
