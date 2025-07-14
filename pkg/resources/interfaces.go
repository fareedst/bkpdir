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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Resource represents any resource that can be cleaned up
type Resource interface {
	Cleanup() error
	String() string
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ResourceManagerInterface defines clean resource management contract
type ResourceManagerInterface interface {
	AddResource(resource Resource)
	AddTempFile(path string)
	AddTempDir(path string)
	RemoveResource(resource Resource)
	Cleanup() error
	CleanupWithPanicRecovery() error
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TempFile represents a temporary file that can be cleaned up
type TempFile struct {
	Path string
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Cleanup removes the temporary file from the filesystem
func (tf *TempFile) Cleanup() error {
	return os.Remove(tf.Path)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// String returns a string representation of the temporary file
func (tf *TempFile) String() string {
	return fmt.Sprintf("TempFile{Path: %s}", tf.Path)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// TempDir represents a temporary directory that can be cleaned up
type TempDir struct {
	Path string
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Cleanup removes the temporary directory and all its contents from the filesystem
func (td *TempDir) Cleanup() error {
	return os.RemoveAll(td.Path)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// String returns a string representation of the temporary directory
func (td *TempDir) String() string {
	return fmt.Sprintf("TempDir{Path: %s}", td.Path)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// ResourceManager manages a collection of resources for automatic cleanup
// Extracted from original errors.go with enhanced thread safety
type ResourceManager struct {
	resources []Resource
	mutex     sync.RWMutex
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// NewResourceManager creates a new ResourceManager instance
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make([]Resource, 0),
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AddResource adds a resource to be tracked for cleanup
func (rm *ResourceManager) AddResource(resource Resource) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.resources = append(rm.resources, resource)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AddTempFile adds a temporary file to be tracked for cleanup
func (rm *ResourceManager) AddTempFile(path string) {
	rm.AddResource(&TempFile{Path: path})
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AddTempDir adds a temporary directory to be tracked for cleanup
func (rm *ResourceManager) AddTempDir(path string) {
	rm.AddResource(&TempDir{Path: path})
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// RemoveResource removes a resource from tracking (typically after successful completion)
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Cleanup cleans up all tracked resources in the ResourceManager
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CleanupWithPanicRecovery cleans up all resources and recovers from panics during cleanup
func (rm *ResourceManager) CleanupWithPanicRecovery() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during cleanup: %v", r)
		}
	}()

	return rm.Cleanup()
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetResourceCount returns the number of currently tracked resources
func (rm *ResourceManager) GetResourceCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()
	return len(rm.resources)
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetResources returns a copy of the currently tracked resources
func (rm *ResourceManager) GetResources() []Resource {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	// Return a copy to prevent external modification
	resources := make([]Resource, len(rm.resources))
	copy(resources, rm.resources)
	return resources
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CleanupIf cleans up resources that match a given predicate
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetResourcesByType returns resources of a specific type
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// CleanupWithContext cleans up resources with context cancellation support
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
