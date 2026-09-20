// This file is part of bkpdir
//
// Package main provides error handling and resource management for BkpDir.
// It includes error types, resource cleanup, and context-aware operations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [REQ-ERROR_HANDLING] Enhanced error handling with structured error types
// [ARCH-ERROR_HANDLING] Structured error handling strategy
package main

import (
	"bkpdir/pkg/formatter"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

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

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: format Message alone or Message WITH underlying Err when present.
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

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: return wrapped Err for errors.Is/As chains.
func (e *ArchiveError) Unwrap() error {
	return e.Err
}

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

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: mirror ArchiveError Error formatting for backup operations.
// BackupError represents a structured error with status code for backup operations
type BackupError struct {
	Message    string
	StatusCode int
	Operation  string
	Path       string
	Err        error
}

func (e *BackupError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *BackupError) Unwrap() error {
	return e.Err
}

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

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message and status code only.
func NewArchiveError(message string, statusCode int) *ArchiveError {
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status code, and underlying Err.
func NewArchiveErrorWithCause(message string, statusCode int, err error) *ArchiveError {
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct ArchiveError with message, status, operation, path, and optional Err.
func NewArchiveErrorWithContext(
	message string,
	statusCode int,
	operation, path string,
	err error,
) *ArchiveError {
	return &ArchiveError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message and status code only.
func NewBackupError(message string, statusCode int) *BackupError {
	return &BackupError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status code, and underlying Err.
func NewBackupErrorWithCause(message string, statusCode int, err error) *BackupError {
	return &BackupError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: construct BackupError with message, status, operation, path, and optional Err.
func NewBackupErrorWithContext(
	message string,
	statusCode int,
	operation, path string,
	err error,
) *BackupError {
	return &BackupError{
		Message:    message,
		StatusCode: statusCode,
		Operation:  operation,
		Path:       path,
		Err:        err,
	}
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match disk-space substrings in error text OR path error with OS no-space/quota/large-file codes.
func IsDiskFullError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "no space left on device") ||
		strings.Contains(errStr, "disk full") ||
		strings.Contains(errStr, "device full") ||
		strings.Contains(errStr, "not enough space") ||
		strings.Contains(errStr, "insufficient disk space") ||
		strings.Contains(errStr, "disk quota exceeded") ||
		strings.Contains(errStr, "file too large") ||
		strings.Contains(errStr, "write: no space left on device") ||
		strings.Contains(errStr, "write: disk full") ||
		strings.Contains(errStr, "write: file too large") ||
		strings.Contains(errStr, "write: insufficient disk space") ||
		strings.Contains(errStr, "write error") ||
		strings.Contains(errStr, "out of disk space") ||
		strings.Contains(errStr, "device out of space") ||
		strings.Contains(errStr, "filesystem full") ||
		strings.Contains(errStr, "volume full") ||
		strings.Contains(errStr, "quota exceeded") ||
		strings.Contains(errStr, "enospc") ||
		strings.Contains(errStr, "edquot") ||
		strings.Contains(errStr, "efbig") ||
		func() bool {
			var pathErr *os.PathError
			if errors.As(err, &pathErr) {
				return pathErr.Err == syscall.ENOSPC ||
					pathErr.Err == syscall.EDQUOT ||
					pathErr.Err == syscall.EFBIG
			}
			return false
		}()
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match permission-denied substrings OR path error with EACCES/EPERM.
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "permission denied") ||
		strings.Contains(errStr, "access denied") ||
		strings.Contains(errStr, "operation not permitted") ||
		strings.Contains(errStr, "insufficient permissions") ||
		strings.Contains(errStr, "insufficient privileges") ||
		strings.Contains(errStr, "access is denied") ||
		strings.Contains(errStr, "cannot access") ||
		strings.Contains(errStr, "eacces") ||
		strings.Contains(errStr, "eperm") ||
		func() bool {
			var pathErr *os.PathError
			if errors.As(err, &pathErr) {
				return pathErr.Err == syscall.EACCES ||
					pathErr.Err == syscall.EPERM
			}
			return false
		}()
}

// - [IMPL-STRUCTURED_ERRORS] [ARCH-ERROR_HANDLING] [REQ-ERROR_HANDLING] — How: match directory-not-found substrings OR path error with ENOENT.
func IsDirectoryNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "no such file or directory") ||
		strings.Contains(errStr, "cannot find the path") ||
		strings.Contains(errStr, "directory not found") ||
		strings.Contains(errStr, "path does not exist") ||
		strings.Contains(errStr, "enoent") ||
		func() bool {
			var pathErr *os.PathError
			if errors.As(err, &pathErr) {
				return pathErr.Err == syscall.ENOENT
			}
			return false
		}()
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
	return os.RemoveAll(td.Path)
}

func (td *TempDir) String() string {
	return fmt.Sprintf("temp dir: %s", td.Path)
}

type ResourceManager struct {
	resources []Resource
	mutex     sync.RWMutex
}

// NewResourceManager creates a new ResourceManager for tracking resources.
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

	var errs []error
	for _, resource := range rm.resources {
		if err := resource.Cleanup(); err != nil {
			errs = append(errs, err)
		}
	}

	// Clear the resources list after cleanup
	rm.resources = make([]Resource, 0)

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}

func (rm *ResourceManager) CleanupWithPanicRecovery() (err error) {
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
	select {
	case <-co.ctx.Done():
		return co.ctx.Err()
	default:
		return nil
	}
}

func (co *ContextualOperation) Cleanup() error {
	return co.rm.Cleanup()
}

func HandleError(err error, cfg ErrorConfig, formatter ErrorFormatter) int {
	if err == nil {
		return 0
	}

	// Check for specific error types
	if IsDiskFullError(err) {
		formatter.PrintDiskFullError(err)
		return cfg.GetStatusCodes()["disk_full"]
	}

	if IsPermissionError(err) {
		formatter.PrintPermissionError(err)
		return cfg.GetStatusCodes()["permission_denied"]
	}

	if IsDirectoryNotFoundError(err) {
		formatter.PrintDirectoryNotFound(err)
		return cfg.GetStatusCodes()["directory_not_found"]
	}

	// Check for structured errors
	if archiveErr, ok := err.(*ArchiveError); ok {
		return HandleArchiveErrorWithInterface(archiveErr, cfg, formatter)
	}

	if backupErr, ok := err.(*BackupError); ok {
		return HandleBackupErrorWithInterface(backupErr, cfg, formatter)
	}

	// Generic error handling
	formatter.PrintError(err.Error())
	return cfg.GetStatusCodes()["general_error"]
}

func HandleArchiveErrorWithInterface(err *ArchiveError, cfg ErrorConfig, formatter ErrorFormatter) int {
	return err.GetStatusCode()
}

func HandleBackupErrorWithInterface(err *BackupError, cfg ErrorConfig, formatter ErrorFormatter) int {
	return err.GetStatusCode()
}

// - [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: classify errno-style errors, print via formatter Format* methods using cfg strings, and return matching cfg status codes.
func HandleArchiveError(err error, cfg *Config, formatter formatter.OutputFormatterInterface) int {
	if err == nil {
		return 0
	}

	// Check for specific error types
	if IsDiskFullError(err) {
		formatter.PrintError(formatter.FormatDiskFullError(err))
		return cfg.StatusDiskFull
	}

	if IsPermissionError(err) {
		formatter.PrintError(formatter.FormatPermissionError(err))
		return cfg.StatusPermissionDenied
	}

	if IsDirectoryNotFoundError(err) {
		formatter.PrintError(formatter.FormatDirectoryNotFound(err))
		return cfg.StatusDirectoryNotFound
	}

	// Check for structured errors
	if archiveErr, ok := err.(*ArchiveError); ok {
		return archiveErr.GetStatusCode()
	}

	if backupErr, ok := err.(*BackupError); ok {
		return backupErr.GetStatusCode()
	}

	// Generic error handling
	formatter.PrintError(err.Error())
	return 1
}

func HandleBackupError(err error, cfg *Config, formatter *OutputFormatter) int {
	return HandleError(err, cfg, formatter)
}

// AtomicWriteFile writes data to a file atomically using a temporary file.
func AtomicWriteFile(path string, data []byte, rm *ResourceManager) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return NewArchiveErrorWithContext("Failed to write temporary file", 1, "atomic_write", path, err)
	}

	tempFile, err := os.CreateTemp(dir, ".tmp-"+filepath.Base(path))
	if err != nil {
		return NewArchiveErrorWithContext("Failed to write temporary file", 1, "atomic_write", path, err)
	}

	tempPath := tempFile.Name()
	rm.AddTempFile(tempPath)

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return NewArchiveErrorWithContext("Failed to write temporary file", 1, "atomic_write", path, err)
	}

	if err := tempFile.Close(); err != nil {
		return NewArchiveErrorWithContext("Failed to write temporary file", 1, "atomic_write", path, err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return NewArchiveErrorWithContext("Failed to write temporary file", 1, "atomic_write", path, err)
	}

	return nil
}

func AtomicWriteFileWithContext(ctx context.Context, path string, data []byte, rm *ResourceManager) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Check for cancellation before file operations
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tempFile, err := os.CreateTemp(dir, ".tmp-"+filepath.Base(path))
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}

	tempPath := tempFile.Name()
	rm.AddTempFile(tempPath)

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to rename temporary file to %s: %w", path, err)
	}

	return nil
}

func SafeMkdirAllWithInterface(path string, perm os.FileMode, cfg ErrorConfig) error {
	if err := os.MkdirAll(path, perm); err != nil {
		if IsPermissionError(err) {
			return fmt.Errorf("permission denied creating directory %s: %w", path, err)
		}
		if IsDiskFullError(err) {
			return fmt.Errorf("disk full creating directory %s: %w", path, err)
		}
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

func SafeMkdirAllWithContextAndInterface(ctx context.Context, path string, perm os.FileMode, cfg ErrorConfig) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return SafeMkdirAllWithInterface(path, perm, cfg)
}

func SafeMkdirAll(path string, perm os.FileMode, cfg *Config) error {
	return SafeMkdirAllWithInterface(path, perm, cfg)
}

func SafeMkdirAllWithContext(ctx context.Context, path string, perm os.FileMode, cfg *Config) error {
	return SafeMkdirAllWithContextAndInterface(ctx, path, perm, cfg)
}

func ValidateDirectoryPath(path string, cfg *Config) error {
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return fmt.Errorf("failed to access directory %s: %w", path, err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	// Check permissions
	if info.Mode()&0200 == 0 {
		return fmt.Errorf("directory is not writable: %s", path)
	}

	return nil
}

func ValidateFilePath(path string, cfg *Config) error {
	if path == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", path)
		}
		return fmt.Errorf("failed to access file %s: %w", path, err)
	}

	// Check if it's actually a file
	if info.IsDir() {
		return fmt.Errorf("path exists but is a directory: %s", path)
	}

	// Check permissions
	if info.Mode()&0400 == 0 {
		return fmt.Errorf("file is not readable: %s", path)
	}

	return nil
}

func WithResourceManager(ctx context.Context) (context.Context, *ResourceManager) {
	rm := NewResourceManager()
	return context.WithValue(ctx, "resource_manager", rm), rm
}

func CheckContextAndCleanup(ctx context.Context, rm *ResourceManager) error {
	select {
	case <-ctx.Done():
		rm.Cleanup()
		return ctx.Err()
	default:
		return nil
	}
}
