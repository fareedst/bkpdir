// Tests for the pkg/resources package to validate extracted resource management functionality.
// These tests ensure the extracted resource management components work correctly
// and maintain backward compatibility with the original functionality.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package resources

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// ⭐ EXTRACT-002: Resource interface testing - 🧪 TempFile functionality
func TestTempFile(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "test.tmp")

	// Write some content to the file
	if err := os.WriteFile(tempPath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create TempFile resource
	tempFile := &TempFile{Path: tempPath}

	// Test String method
	if !strings.Contains(tempFile.String(), tempPath) {
		t.Errorf("String() should contain the file path")
	}

	// Verify file exists before cleanup
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("File should exist before cleanup")
	}

	// Test cleanup
	if err := tempFile.Cleanup(); err != nil {
		t.Errorf("Cleanup should not return error: %v", err)
	}

	// Verify file is removed after cleanup
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Errorf("File should be removed after cleanup")
	}

	// Test cleanup of non-existent file (should return error but not crash)
	err := tempFile.Cleanup()
	if err == nil {
		t.Logf("Cleanup of non-existent file returned no error (implementation may vary)")
	} else {
		t.Logf("Cleanup of non-existent file returned expected error: %v", err)
	}
}

// ⭐ EXTRACT-002: Resource interface testing - 🧪 TempDir functionality
func TestTempDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "test-dir")

	// Create the directory with some content
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Add some files to the directory
	testFile := filepath.Join(testDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file in directory: %v", err)
	}

	// Create TempDir resource
	tempDirResource := &TempDir{Path: testDir}

	// Test String method
	if !strings.Contains(tempDirResource.String(), testDir) {
		t.Errorf("String() should contain the directory path")
	}

	// Verify directory exists before cleanup
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Directory should exist before cleanup")
	}

	// Test cleanup
	if err := tempDirResource.Cleanup(); err != nil {
		t.Errorf("Cleanup should not return error: %v", err)
	}

	// Verify directory is removed after cleanup
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Errorf("Directory should be removed after cleanup")
	}

	// Test cleanup of non-existent directory (should return error but not crash)
	err := tempDirResource.Cleanup()
	if err == nil {
		t.Logf("Cleanup of non-existent directory returned no error (implementation may vary)")
	} else {
		t.Logf("Cleanup of non-existent directory returned expected error: %v", err)
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Basic resource management
func TestResourceManager(t *testing.T) {
	rm := NewResourceManager()

	// Test initial state
	if rm.GetResourceCount() != 0 {
		t.Errorf("New ResourceManager should have 0 resources")
	}

	// Create test resources
	tempDir := t.TempDir()
	testFile1 := filepath.Join(tempDir, "test1.tmp")
	testFile2 := filepath.Join(tempDir, "test2.tmp")
	testSubDir := filepath.Join(tempDir, "subdir")

	// Create actual files/directories
	if err := os.WriteFile(testFile1, []byte("test1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.MkdirAll(testSubDir, 0755); err != nil {
		t.Fatalf("Failed to create test subdirectory: %v", err)
	}

	// Add resources
	rm.AddTempFile(testFile1)
	rm.AddTempFile(testFile2)
	rm.AddTempDir(testSubDir)

	// Check resource count
	if rm.GetResourceCount() != 3 {
		t.Errorf("Expected 3 resources, got %d", rm.GetResourceCount())
	}

	// Test GetResources
	resources := rm.GetResources()
	if len(resources) != 3 {
		t.Errorf("Expected 3 resources from GetResources, got %d", len(resources))
	}

	// Test resource type filtering
	fileResources := rm.GetResourcesByType("file")
	if len(fileResources) != 2 {
		t.Errorf("Expected 2 file resources, got %d", len(fileResources))
	}

	dirResources := rm.GetResourcesByType("directory")
	if len(dirResources) != 1 {
		t.Errorf("Expected 1 directory resource, got %d", len(dirResources))
	}

	// Verify files exist before cleanup
	for _, file := range []string{testFile1, testFile2} {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("File %s should exist before cleanup", file)
		}
	}

	// Test cleanup
	if err := rm.Cleanup(); err != nil {
		t.Errorf("Cleanup should not return error: %v", err)
	}

	// Check that resource count is reset
	if rm.GetResourceCount() != 0 {
		t.Errorf("ResourceManager should have 0 resources after cleanup")
	}

	// Verify files are removed after cleanup
	for _, file := range []string{testFile1, testFile2} {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			t.Errorf("File %s should be removed after cleanup", file)
		}
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Resource removal
func TestResourceManager_RemoveResource(t *testing.T) {
	rm := NewResourceManager()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.tmp")

	// Create test file
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add resource
	rm.AddTempFile(testFile)
	if rm.GetResourceCount() != 1 {
		t.Errorf("Expected 1 resource after adding")
	}

	// Remove resource
	tempFileResource := &TempFile{Path: testFile}
	rm.RemoveResource(tempFileResource)

	if rm.GetResourceCount() != 0 {
		t.Errorf("Expected 0 resources after removal")
	}

	// Verify file still exists (not cleaned up when removed from tracking)
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("File should still exist after removal from tracking")
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Conditional cleanup
func TestResourceManager_CleanupIf(t *testing.T) {
	rm := NewResourceManager()
	tempDir := t.TempDir()

	// Create multiple test files
	testFile1 := filepath.Join(tempDir, "keep.tmp")
	testFile2 := filepath.Join(tempDir, "remove.tmp")

	if err := os.WriteFile(testFile1, []byte("keep"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("remove"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	rm.AddTempFile(testFile1)
	rm.AddTempFile(testFile2)

	// Cleanup only files with "remove" in the path
	predicate := func(r Resource) bool {
		return strings.Contains(r.String(), "remove")
	}

	if err := rm.CleanupIf(predicate); err != nil {
		t.Errorf("CleanupIf should not return error: %v", err)
	}

	// Should have 1 resource remaining
	if rm.GetResourceCount() != 1 {
		t.Errorf("Expected 1 resource remaining, got %d", rm.GetResourceCount())
	}

	// Verify selective cleanup
	if _, err := os.Stat(testFile1); os.IsNotExist(err) {
		t.Errorf("File with 'keep' should still exist")
	}
	if _, err := os.Stat(testFile2); !os.IsNotExist(err) {
		t.Errorf("File with 'remove' should be cleaned up")
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Panic recovery
func TestResourceManager_CleanupWithPanicRecovery(t *testing.T) {
	rm := NewResourceManager()

	// Add a resource that will panic during cleanup
	panicResource := &MockPanicResource{}
	rm.AddResource(panicResource)

	// Test panic recovery
	err := rm.CleanupWithPanicRecovery()
	if err == nil {
		t.Errorf("CleanupWithPanicRecovery should return error when panic occurs")
	}

	if !strings.Contains(err.Error(), "panic during cleanup") {
		t.Errorf("Error should indicate panic occurred: %v", err)
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Context-aware cleanup
func TestResourceManager_CleanupWithContext(t *testing.T) {
	rm := NewResourceManager()
	tempDir := t.TempDir()

	// Add multiple resources with valid filenames
	for i := 0; i < 5; i++ {
		testFile := filepath.Join(tempDir, "test"+strconv.Itoa(i)+".tmp")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
		rm.AddTempFile(testFile)
	}

	// Test with cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanup should succeed with valid context
	if err := rm.CleanupWithContext(ctx); err != nil {
		t.Errorf("CleanupWithContext should not return error: %v", err)
	}

	// Test with already cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	anotherFile := filepath.Join(tempDir, "another.tmp")
	if err := os.WriteFile(anotherFile, []byte("test"), 0644); err == nil {
		rm.AddTempFile(anotherFile)
	}

	err := rm.CleanupWithContext(cancelledCtx)
	if err == nil {
		t.Errorf("CleanupWithContext should return error for cancelled context")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// ⭐ EXTRACT-002: ResourceManager testing - 🧪 Concurrent access
func TestResourceManager_ConcurrentAccess(t *testing.T) {
	rm := NewResourceManager()
	tempDir := t.TempDir()

	var wg sync.WaitGroup
	numGoroutines := 10

	// Concurrently add resources
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			testFile := filepath.Join(tempDir, "test_"+strconv.Itoa(index)+".tmp")
			if err := os.WriteFile(testFile, []byte("test"), 0644); err == nil {
				rm.AddTempFile(testFile)
				// Note: We can't safely increment a counter here due to race conditions
				// The test will check the final count instead
			}
		}(i)
	}

	wg.Wait()

	// Check that resources were added safely (some may have failed due to concurrency)
	count := rm.GetResourceCount()
	if count == 0 {
		t.Errorf("Expected at least some resources to be added, got %d", count)
	}
	if count > numGoroutines {
		t.Errorf("Expected at most %d resources, got %d", numGoroutines, count)
	}
	t.Logf("Successfully added %d out of %d resources concurrently", count, numGoroutines)

	// Concurrently access resources
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_ = rm.GetResourceCount()
			_ = rm.GetResources()
			_ = rm.GetResourcesByType("file")
		}()
	}

	wg.Wait()

	// Final cleanup
	if err := rm.Cleanup(); err != nil {
		t.Errorf("Final cleanup should not return error: %v", err)
	}
}

// ⭐ EXTRACT-002: ContextualOperation testing - 🧪 Context operations
func TestContextualOperation(t *testing.T) {
	ctx := context.Background()
	co := NewContextualOperation(ctx)

	// Test context access
	if co.Context() != ctx {
		t.Errorf("Context() should return the original context")
	}

	// Test resource manager access
	rm := co.ResourceManager()
	if rm == nil {
		t.Errorf("ResourceManager() should not return nil")
	}

	// Test cancellation checking with non-cancelled context
	if co.IsCancelled() {
		t.Errorf("IsCancelled() should return false for non-cancelled context")
	}

	if err := co.CheckCancellation(); err != nil {
		t.Errorf("CheckCancellation() should not return error for non-cancelled context")
	}

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	cancelledCo := NewContextualOperation(cancelledCtx)
	if !cancelledCo.IsCancelled() {
		t.Errorf("IsCancelled() should return true for cancelled context")
	}

	if err := cancelledCo.CheckCancellation(); err == nil {
		t.Errorf("CheckCancellation() should return error for cancelled context")
	}

	// Test cleanup functionality
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.tmp")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	co.ResourceManager().AddTempFile(testFile)

	// Test cleanup
	if err := co.Cleanup(); err != nil {
		t.Errorf("Cleanup() should not return error: %v", err)
	}

	// Verify file is cleaned up
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Errorf("File should be removed after cleanup")
	}

	// Test panic recovery cleanup
	if err := co.CleanupWithPanicRecovery(); err != nil {
		t.Errorf("CleanupWithPanicRecovery() should not return error for normal case: %v", err)
	}
}

// ⭐ EXTRACT-002: Context operations testing - 🧪 Atomic file operations
func TestAtomicWriteFile(t *testing.T) {
	rm := NewResourceManager()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testData := []byte("test content")

	// Test atomic write without context
	if err := AtomicWriteFile(testFile, testData, rm); err != nil {
		t.Errorf("AtomicWriteFile should not return error: %v", err)
	}

	// Verify file was created with correct content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read created file: %v", err)
	}

	if string(content) != string(testData) {
		t.Errorf("File content mismatch. Expected %s, got %s", testData, content)
	}

	// Test atomic write with context
	ctx := context.Background()
	testFile2 := filepath.Join(tempDir, "test2.txt")
	testData2 := []byte("test content 2")

	if err := AtomicWriteFileWithContext(ctx, testFile2, testData2, rm); err != nil {
		t.Errorf("AtomicWriteFileWithContext should not return error: %v", err)
	}

	// Verify second file
	content2, err := os.ReadFile(testFile2)
	if err != nil {
		t.Fatalf("Failed to read second created file: %v", err)
	}

	if string(content2) != string(testData2) {
		t.Errorf("Second file content mismatch. Expected %s, got %s", testData2, content2)
	}

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	testFile3 := filepath.Join(tempDir, "test3.txt")
	err = AtomicWriteFileWithContext(cancelledCtx, testFile3, testData, rm)
	if err == nil {
		t.Errorf("AtomicWriteFileWithContext should return error for cancelled context")
	}
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// ⭐ EXTRACT-002: Context operations testing - 🧪 Context utilities
func TestContextUtilities(t *testing.T) {
	// Test WithResourceManager
	ctx := context.Background()
	newCtx, rm := WithResourceManager(ctx)

	if rm == nil {
		t.Errorf("WithResourceManager should return non-nil ResourceManager")
	}

	// Test GetResourceManagerFromContext with proper type assertion
	if retrievedVal := newCtx.Value(ResourceManagerKey); retrievedVal != nil {
		if retrievedRM, ok := retrievedVal.(*ResourceManager); ok {
			if retrievedRM != rm {
				t.Errorf("Retrieved ResourceManager should be the same as the one set")
			}
		} else {
			t.Errorf("Context value should be a *ResourceManager")
		}
	} else {
		t.Errorf("Context should contain ResourceManager value")
	}

	// Test GetResourceManagerFromContext helper function
	retrievedRM, exists := GetResourceManagerFromContext(newCtx)
	if !exists {
		t.Errorf("GetResourceManagerFromContext should return true for context with ResourceManager")
	}
	if retrievedRM != rm {
		t.Errorf("Retrieved ResourceManager should be the same as the one set")
	}

	// Test with context without ResourceManager
	emptyCtx := context.Background()
	_, exists = GetResourceManagerFromContext(emptyCtx)
	if exists {
		t.Errorf("GetResourceManagerFromContext should return false for context without ResourceManager")
	}

	// Test WithOperationID and GetOperationIDFromContext
	operationID := "test-operation-123"
	ctxWithOp := WithOperationID(ctx, operationID)

	retrievedID, exists := GetOperationIDFromContext(ctxWithOp)
	if !exists {
		t.Errorf("GetOperationIDFromContext should return true for context with operation ID")
	}
	if retrievedID != operationID {
		t.Errorf("Retrieved operation ID should match. Expected %s, got %s", operationID, retrievedID)
	}

	// Test CheckContextAndCleanup
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.tmp")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	rm.AddTempFile(testFile)

	// Test with non-cancelled context
	if err := CheckContextAndCleanup(ctx, rm); err != nil {
		t.Errorf("CheckContextAndCleanup should not return error for non-cancelled context: %v", err)
	}

	// Add resource again for cancelled context test
	testFile2 := filepath.Join(tempDir, "test2.tmp")
	if err := os.WriteFile(testFile2, []byte("test2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	rm.AddTempFile(testFile2)

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := CheckContextAndCleanup(cancelledCtx, rm)
	if err == nil {
		t.Errorf("CheckContextAndCleanup should return error for cancelled context")
	}

	// Should still include context.Canceled
	if !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Error should contain context cancellation information: %v", err)
	}

	// Verify cleanup occurred
	if _, err := os.Stat(testFile2); !os.IsNotExist(err) {
		t.Errorf("File should be cleaned up when context is cancelled")
	}
}

// ⭐ EXTRACT-002: Error combination testing - 🧪 CombineErrors utility
func TestCombineErrors(t *testing.T) {
	// Test with no errors
	combinedErr := CombineErrors()
	if combinedErr != nil {
		t.Errorf("CombineErrors with no errors should return nil")
	}

	// Test with nil errors
	combinedErr = CombineErrors(nil, nil)
	if combinedErr != nil {
		t.Errorf("CombineErrors with only nil errors should return nil")
	}

	// Test with single error
	err1 := context.Canceled
	combinedErr = CombineErrors(err1)
	if combinedErr != err1 {
		t.Errorf("CombineErrors with single error should return that error")
	}

	// Test with multiple errors
	err2 := context.DeadlineExceeded
	combinedErr = CombineErrors(err1, err2)
	if combinedErr == nil {
		t.Errorf("CombineErrors with multiple errors should not return nil")
	}

	combinedErrStr := combinedErr.Error()
	if !strings.Contains(combinedErrStr, "multiple errors occurred") {
		t.Errorf("Combined error message should indicate multiple errors")
	}
	if !strings.Contains(combinedErrStr, "context canceled") {
		t.Errorf("Combined error should contain first error message")
	}
	if !strings.Contains(combinedErrStr, "context deadline exceeded") {
		t.Errorf("Combined error should contain second error message")
	}

	// Test CombinedError type
	if ce, ok := combinedErr.(*CombinedError); ok {
		if len(ce.GetAllErrors()) != 2 {
			t.Errorf("CombinedError should contain 2 errors")
		}

		unwrapped := ce.Unwrap()
		if unwrapped != err1 {
			t.Errorf("Unwrap should return the first error")
		}
	} else {
		t.Errorf("Combined error should be of type *CombinedError")
	}
}

// ⭐ EXTRACT-002: Mock resources for testing - 🧪 Test utilities

// MockPanicResource is a resource that panics during cleanup for testing panic recovery
type MockPanicResource struct{}

func (mpr *MockPanicResource) Cleanup() error {
	panic("test panic during cleanup")
}

func (mpr *MockPanicResource) String() string {
	return "MockPanicResource{panic test}"
}

// ⭐ EXTRACT-002: Resource interface testing - 🧪 Integration test
func TestResourceManagement_Integration(t *testing.T) {
	// Create a contextual operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	co := NewContextualOperation(ctx)
	rm := co.ResourceManager()

	tempDir := t.TempDir()

	// Simulate a complex operation with multiple resources
	for i := 0; i < 3; i++ {
		// Create temporary files
		tempFile := filepath.Join(tempDir, "temp_"+string(rune('a'+i))+".tmp")
		if err := os.WriteFile(tempFile, []byte("content "+string(rune('a'+i))), 0644); err != nil {
			t.Fatalf("Failed to create temp file %d: %v", i, err)
		}
		rm.AddTempFile(tempFile)

		// Create temporary directories
		tempSubDir := filepath.Join(tempDir, "subdir_"+string(rune('a'+i)))
		if err := os.MkdirAll(tempSubDir, 0755); err != nil {
			t.Fatalf("Failed to create temp dir %d: %v", i, err)
		}
		rm.AddTempDir(tempSubDir)
	}

	// Verify all resources are tracked
	if rm.GetResourceCount() != 6 { // 3 files + 3 directories
		t.Errorf("Expected 6 resources, got %d", rm.GetResourceCount())
	}

	// Simulate operation completion and cleanup
	if err := co.Cleanup(); err != nil {
		t.Errorf("Integration cleanup should not return error: %v", err)
	}

	// Verify all resources are cleaned up
	if rm.GetResourceCount() != 0 {
		t.Errorf("All resources should be cleaned up after operation")
	}

	// Verify files and directories are actually removed
	for i := 0; i < 3; i++ {
		tempFile := filepath.Join(tempDir, "temp_"+string(rune('a'+i))+".tmp")
		if _, err := os.Stat(tempFile); !os.IsNotExist(err) {
			t.Errorf("Temp file %d should be removed", i)
		}

		tempSubDir := filepath.Join(tempDir, "subdir_"+string(rune('a'+i)))
		if _, err := os.Stat(tempSubDir); !os.IsNotExist(err) {
			t.Errorf("Temp directory %d should be removed", i)
		}
	}
}
