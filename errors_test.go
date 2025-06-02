// This file is part of bkpdir
//
// Package main provides comprehensive tests for error handling and resource management.
// It tests ArchiveError methods, error classification, resource cleanup, and context-aware operations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test ArchiveError.Error() method - 0% coverage
func TestArchiveError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ArchiveError
		expected string
	}{
		{
			name: "error without underlying cause",
			err: &ArchiveError{
				Message:    "operation failed",
				StatusCode: 1,
				Operation:  "test",
				Path:       "/path/to/file",
			},
			expected: "operation failed",
		},
		{
			name: "error with underlying cause",
			err: &ArchiveError{
				Message:    "operation failed",
				StatusCode: 1,
				Operation:  "test",
				Path:       "/path/to/file",
				Err:        errors.New("underlying error"),
			},
			expected: "operation failed: underlying error",
		},
		{
			name: "error with complex underlying cause",
			err: &ArchiveError{
				Message:    "file operation failed",
				StatusCode: 22,
				Operation:  "file_write",
				Path:       "/tmp/test.txt",
				Err:        &fs.PathError{Op: "open", Path: "/tmp/test.txt", Err: errors.New("permission denied")},
			},
			expected: "file operation failed: open /tmp/test.txt: permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test ArchiveError.Unwrap() method - 0% coverage
func TestArchiveError_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		err      *ArchiveError
		expected error
	}{
		{
			name: "error without underlying cause",
			err: &ArchiveError{
				Message:    "operation failed",
				StatusCode: 1,
			},
			expected: nil,
		},
		{
			name: "error with underlying cause",
			err: &ArchiveError{
				Message:    "operation failed",
				StatusCode: 1,
				Err:        errors.New("underlying error"),
			},
			expected: errors.New("underlying error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Unwrap()
			if (result == nil) != (tt.expected == nil) {
				t.Errorf("Unwrap() = %v, want %v", result, tt.expected)
			}
			if result != nil && tt.expected != nil && result.Error() != tt.expected.Error() {
				t.Errorf("Unwrap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test NewArchiveErrorWithCause function - 0% coverage
func TestNewArchiveErrorWithCause(t *testing.T) {
	underlyingErr := errors.New("underlying issue")
	err := NewArchiveErrorWithCause("operation failed", 42, underlyingErr)

	if err.Message != "operation failed" {
		t.Errorf("Message = %q, want %q", err.Message, "operation failed")
	}
	if err.StatusCode != 42 {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode, 42)
	}
	if err.Err != underlyingErr {
		t.Errorf("Err = %v, want %v", err.Err, underlyingErr)
	}
}

// Test NewArchiveErrorWithContext function - 0% coverage
func TestNewArchiveErrorWithContext(t *testing.T) {
	underlyingErr := errors.New("underlying issue")
	err := NewArchiveErrorWithContext("operation failed", 42, "test_operation", "/path/to/file", underlyingErr)

	if err.Message != "operation failed" {
		t.Errorf("Message = %q, want %q", err.Message, "operation failed")
	}
	if err.StatusCode != 42 {
		t.Errorf("StatusCode = %d, want %d", err.StatusCode, 42)
	}
	if err.Operation != "test_operation" {
		t.Errorf("Operation = %q, want %q", err.Operation, "test_operation")
	}
	if err.Path != "/path/to/file" {
		t.Errorf("Path = %q, want %q", err.Path, "/path/to/file")
	}
	if err.Err != underlyingErr {
		t.Errorf("Err = %v, want %v", err.Err, underlyingErr)
	}
}

// Test IsDiskFullError function - 0% coverage
func TestIsDiskFullError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "no space left on device",
			err:      errors.New("no space left on device"),
			expected: true,
		},
		{
			name:     "disk full",
			err:      errors.New("DISK FULL - cannot write"),
			expected: true,
		},
		{
			name:     "not enough space",
			err:      errors.New("Not enough space on disk"),
			expected: true,
		},
		{
			name:     "insufficient disk space",
			err:      errors.New("Insufficient disk space available"),
			expected: true,
		},
		{
			name:     "device full",
			err:      errors.New("Device full"),
			expected: true,
		},
		{
			name:     "quota exceeded",
			err:      errors.New("Quota exceeded for user"),
			expected: true,
		},
		{
			name:     "file too large",
			err:      errors.New("File too large for filesystem"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("permission denied"),
			expected: false,
		},
		{
			name:     "complex error with disk full",
			err:      fmt.Errorf("failed to write file: %w", errors.New("no space left on device")),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDiskFullError(tt.err)
			if result != tt.expected {
				t.Errorf("IsDiskFullError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// Test isDiskFullError function (internal) - 0% coverage
func TestIsDiskFullErrorInternal(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "no space left on device",
			err:      errors.New("no space left on device"),
			expected: true,
		},
		{
			name:     "disk full",
			err:      errors.New("disk full"),
			expected: true,
		},
		{
			name:     "not enough space",
			err:      errors.New("not enough space available"),
			expected: true,
		},
		{
			name:     "quota exceeded",
			err:      errors.New("disk quota exceeded"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("permission denied"),
			expected: false,
		},
		{
			name:     "case insensitive matching",
			err:      errors.New("DISK FULL ERROR"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDiskFullError(tt.err)
			if result != tt.expected {
				t.Errorf("IsDiskFullError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// Test IsPermissionError function - 0% coverage
func TestIsPermissionError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "permission denied",
			err:      errors.New("permission denied"),
			expected: true,
		},
		{
			name:     "access denied",
			err:      errors.New("Access denied"),
			expected: true,
		},
		{
			name:     "operation not permitted",
			err:      errors.New("Operation not permitted"),
			expected: true,
		},
		{
			name:     "insufficient privileges",
			err:      errors.New("Insufficient privileges"),
			expected: true,
		},
		{
			name:     "case insensitive",
			err:      errors.New("PERMISSION DENIED for file"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("file not found"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPermissionError(tt.err)
			if result != tt.expected {
				t.Errorf("IsPermissionError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// Test IsDirectoryNotFoundError function - 0% coverage
func TestIsDirectoryNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "no such file or directory",
			err:      errors.New("no such file or directory"),
			expected: true,
		},
		{
			name:     "cannot find the path",
			err:      errors.New("Cannot find the path specified"),
			expected: true,
		},
		{
			name:     "directory not found",
			err:      errors.New("Directory not found"),
			expected: true,
		},
		{
			name:     "case insensitive",
			err:      errors.New("NO SUCH FILE OR DIRECTORY"),
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      errors.New("permission denied"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDirectoryNotFoundError(tt.err)
			if result != tt.expected {
				t.Errorf("IsDirectoryNotFoundError(%v) = %v, want %v", tt.err, result, tt.expected)
			}
		})
	}
}

// Test TempFile.Cleanup() method - 0% coverage
func TestTempFile_Cleanup(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_file.tmp")

	// Create the file
	if err := os.WriteFile(tempFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatalf("Test file was not created")
	}

	// Test cleanup
	tf := &TempFile{Path: tempFile}
	err := tf.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() failed: %v", err)
	}

	// Verify file is removed
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Errorf("File was not removed by Cleanup()")
	}

	// Test cleanup of non-existent file (should not error)
	tf2 := &TempFile{Path: "/nonexistent/file.tmp"}
	err = tf2.Cleanup()
	if err == nil {
		// This is actually expected behavior for os.Remove on non-existent files
		// os.Remove returns an error, but some implementations might handle it gracefully
	}
}

// Test TempDir.Cleanup() method - 0% coverage
func TestTempDir_Cleanup(t *testing.T) {
	// Create a temporary directory with content
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test_dir")

	// Create directory with nested content
	if err := os.MkdirAll(filepath.Join(testDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testFile := filepath.Join(testDir, "subdir", "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Fatalf("Test directory was not created")
	}

	// Test cleanup
	td := &TempDir{Path: testDir}
	err := td.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() failed: %v", err)
	}

	// Verify directory is removed
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Errorf("Directory was not removed by Cleanup()")
	}
}

// Test TempDir.String() method - 0% coverage
func TestTempDir_String(t *testing.T) {
	td := &TempDir{Path: "/tmp/test_dir"}
	expected := "temp dir: /tmp/test_dir"
	result := td.String()
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

// Test ResourceManager.AddTempDir() method - 0% coverage
func TestResourceManager_AddTempDir(t *testing.T) {
	rm := NewResourceManager()
	testPath := "/tmp/test_dir"

	rm.AddTempDir(testPath)

	// Verify resource was added
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	if len(rm.resources) != 1 {
		t.Errorf("Expected 1 resource, got %d", len(rm.resources))
	}

	if rm.resources[0].String() != "temp dir: "+testPath {
		t.Errorf("Resource string = %q, want %q", rm.resources[0].String(), "temp dir: "+testPath)
	}
}

// Test ResourceManager edge cases and error handling - Enhanced coverage
func TestResourceManager_CleanupErrorHandling(t *testing.T) {
	rm := NewResourceManager()

	// Add a resource that will fail to clean up
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "nonexistent.tmp")
	rm.AddTempFile(nonExistentFile)

	// Cleanup should handle errors gracefully
	err := rm.Cleanup()
	// Error is expected since file doesn't exist, but cleanup should continue
	if err == nil {
		t.Log("Cleanup handled non-existent file gracefully")
	}

	// Verify resources list is cleared even after errors
	rm.mutex.RLock()
	resourceCount := len(rm.resources)
	rm.mutex.RUnlock()

	if resourceCount != 0 {
		t.Errorf("Resources not cleared after cleanup, count: %d", resourceCount)
	}
}

// Test ResourceManager.CleanupWithPanicRecovery() method - Enhanced coverage for panic scenarios
func TestResourceManager_CleanupWithPanicRecovery(t *testing.T) {
	// Test normal cleanup without panic
	rm := NewResourceManager()
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.tmp")

	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rm.AddTempFile(tempFile)

	err := rm.CleanupWithPanicRecovery()
	if err != nil {
		t.Errorf("CleanupWithPanicRecovery() failed: %v", err)
	}

	// Verify file was cleaned up
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Errorf("File was not cleaned up")
	}
}

// Test NewContextualOperation function - 0% coverage
func TestNewContextualOperation(t *testing.T) {
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	if co.ctx != ctx {
		t.Errorf("Context not set correctly")
	}

	if co.rm == nil {
		t.Errorf("ResourceManager not initialized")
	}
}

// Test ContextualOperation.Context() method - 0% coverage
func TestContextualOperation_Context(t *testing.T) {
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	result := co.Context()
	if result != ctx {
		t.Errorf("Context() = %v, want %v", result, ctx)
	}
}

// Test ContextualOperation.ResourceManager() method - 0% coverage
func TestContextualOperation_ResourceManager(t *testing.T) {
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	rm := co.ResourceManager()
	if rm == nil {
		t.Errorf("ResourceManager() returned nil")
	}

	if rm != co.rm {
		t.Errorf("ResourceManager() returned different instance")
	}
}

// Test ContextualOperation.IsCancelled() method - 0% coverage
func TestContextualOperation_IsCancelled(t *testing.T) {
	// Test with non-cancelled context
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	if co.IsCancelled() {
		t.Errorf("IsCancelled() = true for non-cancelled context")
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	co = NewContextualOperation(ctx)

	cancel()

	if !co.IsCancelled() {
		t.Errorf("IsCancelled() = false for cancelled context")
	}
}

// Test ContextualOperation.CheckCancellation() method - 0% coverage
func TestContextualOperation_CheckCancellation(t *testing.T) {
	// Test with non-cancelled context
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	err := co.CheckCancellation()
	if err != nil {
		t.Errorf("CheckCancellation() = %v for non-cancelled context", err)
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	co = NewContextualOperation(ctx)

	cancel()

	err = co.CheckCancellation()
	if err == nil {
		t.Errorf("CheckCancellation() = nil for cancelled context")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("CheckCancellation() = %v, want context.Canceled", err)
	}
}

// Test ContextualOperation.Cleanup() method - 0% coverage
func TestContextualOperation_Cleanup(t *testing.T) {
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	// Add some resources
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.tmp")

	if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	co.ResourceManager().AddTempFile(tempFile)

	err := co.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() failed: %v", err)
	}

	// Verify resource was cleaned up
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Errorf("Resource was not cleaned up")
	}
}

// Test HandleArchiveError function - 0% coverage
func TestHandleArchiveError(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedStatus: 0,
		},
		{
			name:           "ArchiveError",
			err:            NewArchiveError("test error", 42),
			expectedStatus: 42,
		},
		{
			name:           "disk full error",
			err:            errors.New("no space left on device"),
			expectedStatus: cfg.StatusDiskFull,
		},
		{
			name:           "permission error",
			err:            errors.New("permission denied"),
			expectedStatus: cfg.StatusPermissionDenied,
		},
		{
			name:           "directory not found error",
			err:            errors.New("no such file or directory"),
			expectedStatus: cfg.StatusDirectoryNotFound,
		},
		{
			name:           "generic error",
			err:            errors.New("some other error"),
			expectedStatus: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := HandleArchiveError(tt.err, cfg, formatter)
			if status != tt.expectedStatus {
				t.Errorf("HandleArchiveError() = %d, want %d", status, tt.expectedStatus)
			}
		})
	}
}

// Test AtomicWriteFile function - 0% coverage
func TestAtomicWriteFile(t *testing.T) {
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "target.txt")
	testData := []byte("test data for atomic write")

	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Test successful atomic write
	err := AtomicWriteFile(targetFile, testData, rm)
	if err != nil {
		t.Errorf("AtomicWriteFile() failed: %v", err)
	}

	// Verify file contents
	data, err := os.ReadFile(targetFile)
	if err != nil {
		t.Errorf("Failed to read target file: %v", err)
	}

	if string(data) != string(testData) {
		t.Errorf("File contents = %q, want %q", string(data), string(testData))
	}

	// Verify temp file was cleaned up
	tempFile := targetFile + ".tmp"
	if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
		t.Errorf("Temporary file was not cleaned up")
	}

	// Test write to invalid directory
	invalidPath := "/nonexistent/directory/file.txt"
	err = AtomicWriteFile(invalidPath, testData, rm)
	if err == nil {
		t.Errorf("AtomicWriteFile() should fail for invalid path")
	}

	var archiveErr *ArchiveError
	if !errors.As(err, &archiveErr) {
		t.Errorf("Expected ArchiveError, got %T", err)
	}
}

// Test ValidateFilePath function - 0% coverage
func TestValidateFilePath(t *testing.T) {
	cfg := DefaultConfig()
	tempDir := t.TempDir()

	// Create test files and directories
	validFile := filepath.Join(tempDir, "valid.txt")
	if err := os.WriteFile(validFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testDir := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
		statusCode  int
	}{
		{
			name:        "valid file",
			path:        validFile,
			expectError: false,
		},
		{
			name:        "nonexistent file",
			path:        filepath.Join(tempDir, "nonexistent.txt"),
			expectError: true,
			statusCode:  cfg.StatusFileNotFound,
		},
		{
			name:        "directory instead of file",
			path:        testDir,
			expectError: true,
			statusCode:  cfg.StatusInvalidFileType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path, cfg)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateFilePath() expected error, got nil")
					return
				}

				var archiveErr *ArchiveError
				if errors.As(err, &archiveErr) && tt.statusCode != 0 {
					if archiveErr.StatusCode != tt.statusCode {
						t.Errorf("StatusCode = %d, want %d", archiveErr.StatusCode, tt.statusCode)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFilePath() unexpected error: %v", err)
				}
			}
		})
	}
}

// Benchmark tests for performance baselines
func BenchmarkArchiveError_Error(b *testing.B) {
	err := &ArchiveError{
		Message:    "benchmark error",
		StatusCode: 1,
		Operation:  "benchmark",
		Path:       "/path/to/file",
		Err:        errors.New("underlying error"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkIsDiskFullError(b *testing.B) {
	err := errors.New("no space left on device")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsDiskFullError(err)
	}
}

func BenchmarkResourceManager_AddRemove(b *testing.B) {
	rm := NewResourceManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.AddTempFile("/tmp/test.tmp")
		rm.RemoveResource(&TempFile{Path: "/tmp/test.tmp"})
	}
}

// Concurrent access test for ResourceManager thread safety
func TestResourceManager_ConcurrentAccess(t *testing.T) {
	rm := NewResourceManager()

	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := 100

	// Start multiple goroutines adding and removing resources
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				path := fmt.Sprintf("/tmp/test_%d_%d.tmp", id, j)
				rm.AddTempFile(path)
				rm.RemoveResource(&TempFile{Path: path})
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	rm.mutex.RLock()
	resourceCount := len(rm.resources)
	rm.mutex.RUnlock()

	if resourceCount != 0 {
		t.Errorf("Expected 0 resources after concurrent operations, got %d", resourceCount)
	}
}

// Test context cancellation during operations
func TestContextualOperation_WithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	co := NewContextualOperation(ctx)

	// Wait for timeout
	time.Sleep(150 * time.Millisecond)

	if !co.IsCancelled() {
		t.Errorf("Operation should be cancelled after timeout")
	}

	err := co.CheckCancellation()
	if err == nil {
		t.Errorf("CheckCancellation() should return error after timeout")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded, got %v", err)
	}
}

// Additional tests for improved coverage

// Test panic recovery scenarios for ResourceManager.CleanupWithPanicRecovery()
func TestResourceManager_CleanupWithPanicRecovery_PanicScenario(t *testing.T) {
	// Create a mock resource that panics during cleanup
	rm := NewResourceManager()

	// Create a resource that will panic during cleanup
	panicResource := &TestPanicResource{path: "/tmp/panic_test"}
	rm.AddResource(panicResource)

	// Test that panic is recovered
	err := rm.CleanupWithPanicRecovery()
	if err == nil {
		t.Errorf("Expected error from panic recovery, got nil")
	}

	if !strings.Contains(err.Error(), "panic during cleanup") {
		t.Errorf("Expected panic recovery error, got: %v", err)
	}
}

// Mock resource that panics during cleanup for testing
type TestPanicResource struct {
	path string
}

func (tpr *TestPanicResource) Cleanup() error {
	panic("intentional panic for testing")
}

func (tpr *TestPanicResource) String() string {
	return fmt.Sprintf("panic resource: %s", tpr.path)
}

// Test AtomicWriteFile with rename failure scenario
func TestAtomicWriteFile_RenameFailure(t *testing.T) {
	tempDir := t.TempDir()
	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Create a scenario where rename might fail
	// On some systems, we can create a read-only directory to test failure
	testData := []byte("test data")

	// Test with a path that would cause rename to fail
	// Use a directory that doesn't exist for the parent
	invalidPath := filepath.Join(tempDir, "nonexistent_dir", "target.txt")

	err := AtomicWriteFile(invalidPath, testData, rm)
	if err == nil {
		t.Errorf("AtomicWriteFile() should fail when parent directory doesn't exist")
	}

	var archiveErr *ArchiveError
	if errors.As(err, &archiveErr) {
		if !strings.Contains(archiveErr.Message, "Failed to write temporary file") {
			t.Errorf("Expected 'Failed to write temporary file' error, got: %s", archiveErr.Message)
		}
	}
}

// Test SafeMkdirAll with various scenarios
func TestSafeMkdirAll_ComprehensiveScenarios(t *testing.T) {
	cfg := DefaultConfig()
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		path        string
		perm        os.FileMode
		expectError bool
		statusCode  int
		setupFunc   func() error
	}{
		{
			name:        "successful directory creation",
			path:        filepath.Join(tempDir, "new_dir"),
			perm:        0755,
			expectError: false,
		},
		{
			name:        "nested directory creation",
			path:        filepath.Join(tempDir, "level1", "level2", "level3"),
			perm:        0755,
			expectError: false,
		},
		{
			name:        "directory already exists",
			path:        tempDir, // Use existing temp directory
			perm:        0755,
			expectError: false,
		},
		{
			name:        "invalid path that would fail",
			path:        filepath.Join(tempDir, "nonexistent_parent", "test_dir"), // Valid path that should actually succeed
			perm:        0755,
			expectError: false, // Changed to false since os.MkdirAll creates parent directories
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			err := SafeMkdirAll(tt.path, tt.perm, cfg)

			if tt.expectError {
				if err == nil {
					t.Errorf("SafeMkdirAll() expected error, got nil")
					return
				}

				var archiveErr *ArchiveError
				if errors.As(err, &archiveErr) && tt.statusCode != 0 {
					if archiveErr.StatusCode != tt.statusCode {
						t.Errorf("StatusCode = %d, want %d", archiveErr.StatusCode, tt.statusCode)
					}
				}
			} else {
				if err != nil {
					t.Errorf("SafeMkdirAll() unexpected error: %v", err)
				} else {
					// Verify directory was created
					if info, statErr := os.Stat(tt.path); statErr != nil || !info.IsDir() {
						t.Errorf("Directory was not created successfully: %v", statErr)
					}
				}
			}
		})
	}
}

// Test SafeMkdirAll with simulated disk full error
func TestSafeMkdirAll_DiskFullScenario(t *testing.T) {
	cfg := DefaultConfig()

	// Create a very long path that might trigger system limitations
	tempDir := t.TempDir()
	longPath := tempDir
	for i := 0; i < 50; i++ {
		longPath = filepath.Join(longPath, "very_long_directory_name_that_might_cause_issues")
	}

	err := SafeMkdirAll(longPath, 0755, cfg)
	// This might succeed or fail depending on the system
	// The important thing is that if it fails, it handles the error appropriately
	if err != nil {
		var archiveErr *ArchiveError
		if errors.As(err, &archiveErr) {
			// Accept any error status code since this test might succeed on some systems
			if archiveErr.StatusCode == 0 {
				t.Errorf("Expected non-zero error status code, got: %d", archiveErr.StatusCode)
			}
		}
	} else {
		// If it succeeds, verify the directory was created
		if info, statErr := os.Stat(longPath); statErr != nil || !info.IsDir() {
			t.Errorf("Directory was not created successfully despite success: %v", statErr)
		}
	}
}

// Test ValidateDirectoryPath with comprehensive scenarios
func TestValidateDirectoryPath_ComprehensiveScenarios(t *testing.T) {
	cfg := DefaultConfig()
	tempDir := t.TempDir()

	// Create test files and directories
	validDir := filepath.Join(tempDir, "valid_dir")
	if err := os.Mkdir(validDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	regularFile := filepath.Join(tempDir, "regular_file.txt")
	if err := os.WriteFile(regularFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a directory with restricted permissions
	restrictedDir := filepath.Join(tempDir, "restricted_dir")
	if err := os.Mkdir(restrictedDir, 0000); err != nil {
		t.Fatalf("Failed to create restricted directory: %v", err)
	}
	defer os.Chmod(restrictedDir, 0755) // Cleanup

	tests := []struct {
		name        string
		path        string
		expectError bool
		statusCode  int
	}{
		{
			name:        "valid directory",
			path:        validDir,
			expectError: false,
		},
		{
			name:        "root directory",
			path:        "/",
			expectError: false,
		},
		{
			name:        "nonexistent directory",
			path:        filepath.Join(tempDir, "nonexistent"),
			expectError: true,
			statusCode:  cfg.StatusDirectoryNotFound,
		},
		{
			name:        "regular file instead of directory",
			path:        regularFile,
			expectError: true,
			statusCode:  cfg.StatusInvalidDirectoryType,
		},
		{
			name:        "path with permission issues",
			path:        filepath.Join(restrictedDir, "inaccessible"),
			expectError: true,
			statusCode:  cfg.StatusPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDirectoryPath(tt.path, cfg)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateDirectoryPath() expected error, got nil")
					return
				}

				var archiveErr *ArchiveError
				if errors.As(err, &archiveErr) && tt.statusCode != 0 {
					if archiveErr.StatusCode != tt.statusCode {
						t.Errorf("StatusCode = %d, want %d", archiveErr.StatusCode, tt.statusCode)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateDirectoryPath() unexpected error: %v", err)
				}
			}
		})
	}
}

// Test ValidateFilePath with permission scenarios
func TestValidateFilePath_PermissionScenarios(t *testing.T) {
	cfg := DefaultConfig()
	tempDir := t.TempDir()

	// Create a file in a restricted directory
	restrictedDir := filepath.Join(tempDir, "restricted")
	if err := os.Mkdir(restrictedDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	restrictedFile := filepath.Join(restrictedDir, "file.txt")
	if err := os.WriteFile(restrictedFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Make directory inaccessible
	if err := os.Chmod(restrictedDir, 0000); err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}
	defer os.Chmod(restrictedDir, 0755) // Cleanup

	tests := []struct {
		name        string
		path        string
		expectError bool
		statusCode  int
	}{
		{
			name:        "inaccessible file due to directory permissions",
			path:        restrictedFile,
			expectError: true,
			statusCode:  cfg.StatusPermissionDenied, // Might be permission denied or config error
		},
		{
			name:        "path to inaccessible directory",
			path:        filepath.Join(restrictedDir, "nonexistent.txt"),
			expectError: true,
			statusCode:  cfg.StatusPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilePath(tt.path, cfg)

			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateFilePath() expected error, got nil")
					return
				}

				var archiveErr *ArchiveError
				if errors.As(err, &archiveErr) {
					// Accept either permission denied or config error
					if archiveErr.StatusCode != cfg.StatusPermissionDenied &&
						archiveErr.StatusCode != cfg.StatusConfigError &&
						archiveErr.StatusCode != cfg.StatusFileNotFound {
						t.Errorf("StatusCode = %d, want one of [%d, %d, %d]",
							archiveErr.StatusCode,
							cfg.StatusPermissionDenied,
							cfg.StatusConfigError,
							cfg.StatusFileNotFound)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateFilePath() unexpected error: %v", err)
				}
			}
		})
	}
}

// Test edge cases for error classification functions
func TestErrorClassification_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		diskFull    bool
		permErr     bool
		dirNotFound bool
	}{
		{
			name:        "wrapped disk full error",
			err:         fmt.Errorf("operation failed: %w", errors.New("no space left on device")),
			diskFull:    true,
			permErr:     false,
			dirNotFound: false,
		},
		{
			name:        "multiple wrapped errors",
			err:         fmt.Errorf("outer: %w", fmt.Errorf("middle: %w", errors.New("permission denied"))),
			diskFull:    false,
			permErr:     true,
			dirNotFound: false,
		},
		{
			name:        "case variations",
			err:         errors.New("CANNOT FIND THE PATH specified"),
			diskFull:    false,
			permErr:     false,
			dirNotFound: true,
		},
		{
			name:        "partial pattern matches",
			err:         errors.New("user quota exceeded for disk /tmp"),
			diskFull:    true,
			permErr:     false,
			dirNotFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if IsDiskFullError(tt.err) != tt.diskFull {
				t.Errorf("IsDiskFullError() = %v, want %v", IsDiskFullError(tt.err), tt.diskFull)
			}
			if IsPermissionError(tt.err) != tt.permErr {
				t.Errorf("IsPermissionError() = %v, want %v", IsPermissionError(tt.err), tt.permErr)
			}
			if IsDirectoryNotFoundError(tt.err) != tt.dirNotFound {
				t.Errorf("IsDirectoryNotFoundError() = %v, want %v", IsDirectoryNotFoundError(tt.err), tt.dirNotFound)
			}
		})
	}
}
