// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// Package resources provides resource management utilities for CLI applications.
// It includes resource tracking, cleanup, and lifecycle management patterns
// extracted from the BkpDir application for reuse across Go CLI applications.
//
// This package handles temporary files, directories, and other resources
// that need automatic cleanup with support for panic recovery and context cancellation.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package resources

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// Resource interface: any resource that can be cleaned up.
type Resource interface {
	Cleanup() error
	String() string
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// ResourceManagerInterface: contract for resource tracking and cleanup.
type ResourceManagerInterface interface {
	AddResource(resource Resource)
	AddTempFile(path string)
	AddTempDir(path string)
	RemoveResource(resource Resource)
	Cleanup() error
	CleanupWithPanicRecovery() error
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// TempFile: cleans up via os.Remove.
type TempFile struct {
	Path string
}

// [IMPL-RESOURCE_MANAGER] — TempFile.Cleanup: removes file from filesystem.
func (tf *TempFile) Cleanup() error {
	return os.Remove(tf.Path)
}

// [IMPL-RESOURCE_MANAGER] — TempFile.String: path-based identity for resource matching.
func (tf *TempFile) String() string {
	return fmt.Sprintf("TempFile{Path: %s}", tf.Path)
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// TempDir: cleans up via os.RemoveAll.
type TempDir struct {
	Path string
}

// [IMPL-RESOURCE_MANAGER] — TempDir.Cleanup: removes directory recursively.
func (td *TempDir) Cleanup() error {
	return os.RemoveAll(td.Path)
}

// [IMPL-RESOURCE_MANAGER] — TempDir.String: path-based identity for resource matching.
func (td *TempDir) String() string {
	return fmt.Sprintf("TempDir{Path: %s}", td.Path)
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// ResourceManager: thread-safe collection of resources for automatic cleanup.
type ResourceManager struct {
	resources []Resource
	mutex     sync.RWMutex
}

// [IMPL-RESOURCE_MANAGER] — creates a new ResourceManager with empty resource list.
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make([]Resource, 0),
	}
}

// [IMPL-RESOURCE_MANAGER] — ADD_RESOURCE: thread-safe resource registration.
func (rm *ResourceManager) AddResource(resource Resource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = append(rm.resources, resource)
}

// [IMPL-RESOURCE_MANAGER] — ADD_TEMP_FILE: convenience wrapper for TempFile.
func (rm *ResourceManager) AddTempFile(path string) {
	rm.AddResource(&TempFile{Path: path})
}

// [IMPL-RESOURCE_MANAGER] — ADD_TEMP_DIR: convenience wrapper for TempDir.
func (rm *ResourceManager) AddTempDir(path string) {
	rm.AddResource(&TempDir{Path: path})
}

// [IMPL-RESOURCE_MANAGER] — REMOVE_RESOURCE: removes by String() match without cleanup.
func (rm *ResourceManager) RemoveResource(resource Resource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for i, r := range rm.resources {
		// Use string comparison for resource matching
		if r.String() == resource.String() {
			// Remove the resource from the slice
			rm.resources = append(rm.resources[:i], rm.resources[i+1:]...)
			break
		}
	}
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// CLEANUP: cleans all tracked resources, continues on individual failure.
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

	rm.resources = make([]Resource, 0) // Clear the slice
	return lastError
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// CLEANUP_WITH_PANIC_RECOVERY: wraps Cleanup with defer/recover.
func (rm *ResourceManager) CleanupWithPanicRecovery() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during cleanup: %v", r)
		}
	}()

	return rm.Cleanup()
}

// [IMPL-RESOURCE_MANAGER] — GetResourceCount: read-locked resource count.
func (rm *ResourceManager) GetResourceCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return len(rm.resources)
}

// [IMPL-RESOURCE_MANAGER] — GetResources: returns defensive copy of tracked resources.
func (rm *ResourceManager) GetResources() []Resource {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	// Return a copy to prevent external modification
	resources := make([]Resource, len(rm.resources))
	copy(resources, rm.resources)
	return resources
}

// [IMPL-RESOURCE_MANAGER] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// CLEANUP_IF: predicate-based selective cleanup.
func (rm *ResourceManager) CleanupIf(predicate func(Resource) bool) error {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	var lastError error
	var remainingResources []Resource

	for _, resource := range rm.resources {
		if predicate(resource) {
			if err := resource.Cleanup(); err != nil {
				lastError = err
				// Continue with other resources even if this one fails
			}
		} else {
			remainingResources = append(remainingResources, resource)
		}
	}

	rm.resources = remainingResources
	return lastError
}

// [IMPL-RESOURCE_MANAGER] — GetResourcesByType: filters by "file" or "directory" type.
func (rm *ResourceManager) GetResourcesByType(resourceType string) []Resource {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	var matchingResources []Resource
	for _, resource := range rm.resources {
		switch resourceType {
		case "file":
			if _, ok := resource.(*TempFile); ok {
				matchingResources = append(matchingResources, resource)
			}
		case "directory":
			if _, ok := resource.(*TempDir); ok {
				matchingResources = append(matchingResources, resource)
			}
		default:
			// For unknown types, use string matching
			if resource.String() == resourceType {
				matchingResources = append(matchingResources, resource)
			}
		}
	}

	return matchingResources
}

// [IMPL-RESOURCE_MANAGER] [IMPL-CONTEXT_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-CONTEXT_SUPPORT]
// CLEANUP_WITH_CONTEXT: checks cancellation between resource cleanup operations.
func (rm *ResourceManager) CleanupWithContext(ctx context.Context) error {
	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	var lastError error
	for _, resource := range rm.resources {
		// Check for cancellation between resource cleanup operations
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := resource.Cleanup(); err != nil {
			lastError = err
			// Continue cleanup even if individual operations fail
		}
	}

	rm.resources = make([]Resource, 0) // Clear the slice
	return lastError
}
