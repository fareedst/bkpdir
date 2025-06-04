// Package fileops provides file operations and utilities for CLI applications.
//
// This file contains directory traversal functionality with exclusion support.
package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ⭐ EXTRACT-006: Directory traversal system extracted - 🔍

// TraversalOptions configures directory traversal behavior
type TraversalOptions struct {
	ExcludePatterns  []string
	FollowSymlinks   bool
	MaxDepth         int
	IgnoreHidden     bool
	IgnorePermErrors bool
}

// FileVisitor is called for each file or directory encountered during traversal
type FileVisitor func(path string, info os.FileInfo, err error) error

// Traverser defines the interface for directory traversal operations
type Traverser interface {
	Walk(root string, visitor FileVisitor) error
	WalkWithOptions(root string, options TraversalOptions, visitor FileVisitor) error
	ListFiles(root string, recursive bool) ([]string, error)
	ListFilesWithExclusions(root string, excludePatterns []string, recursive bool) ([]string, error)
}

// DefaultTraverser implements the Traverser interface
type DefaultTraverser struct {
	matcher *PatternMatcher
}

// NewTraverser creates a new DefaultTraverser instance
func NewTraverser() Traverser {
	// ⭐ EXTRACT-006: Traverser initialization - 🔍
	return &DefaultTraverser{}
}

// NewTraverserWithPatterns creates a traverser with exclusion patterns
func NewTraverserWithPatterns(patterns []string) Traverser {
	// ⭐ EXTRACT-006: Traverser with patterns initialization - 🔍
	return &DefaultTraverser{
		matcher: NewPatternMatcher(patterns),
	}
}

// Walk performs a basic directory walk
func (dt *DefaultTraverser) Walk(root string, visitor FileVisitor) error {
	// ⭐ EXTRACT-006: Basic directory walk implementation - 🔍
	options := TraversalOptions{
		FollowSymlinks:   false,
		MaxDepth:         -1, // No limit
		IgnoreHidden:     false,
		IgnorePermErrors: false,
	}
	return dt.WalkWithOptions(root, options, visitor)
}

// WalkWithOptions performs directory walk with specified options
func (dt *DefaultTraverser) WalkWithOptions(root string, options TraversalOptions, visitor FileVisitor) error {
	// ⭐ EXTRACT-006: Directory walk with options implementation - 🔍

	// Validate the root path
	if err := ValidatePath(root); err != nil {
		return fmt.Errorf("invalid root path: %v", err)
	}

	// Check if root exists and is accessible
	if err := ValidateExistence(root); err != nil {
		return err
	}

	// Create pattern matcher for exclusions
	var exclusionMatcher *PatternMatcher
	if len(options.ExcludePatterns) > 0 {
		exclusionMatcher = NewPatternMatcher(options.ExcludePatterns)
	} else if dt.matcher != nil {
		exclusionMatcher = dt.matcher
	}

	return dt.walkRecursive(root, root, 0, options, exclusionMatcher, visitor)
}

// walkRecursive performs the actual recursive directory traversal
func (dt *DefaultTraverser) walkRecursive(
	root, currentPath string,
	currentDepth int,
	options TraversalOptions,
	exclusionMatcher *PatternMatcher,
	visitor FileVisitor,
) error {
	// ⭐ EXTRACT-006: Recursive directory walk implementation - 🔍

	// Check depth limit
	if options.MaxDepth >= 0 && currentDepth > options.MaxDepth {
		return nil
	}

	// Get file info
	info, err := os.Lstat(currentPath)
	if err != nil {
		if options.IgnorePermErrors && os.IsPermission(err) {
			return nil
		}
		return visitor(currentPath, nil, err)
	}

	// Get relative path for exclusion checking
	relPath, err := filepath.Rel(root, currentPath)
	if err != nil {
		return err
	}

	// Normalize path for exclusion checking
	if relPath == "." {
		relPath = ""
	}

	// Check exclusion patterns
	if exclusionMatcher != nil && relPath != "" && exclusionMatcher.ShouldExclude(relPath) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Check hidden files
	if options.IgnoreHidden && relPath != "" && strings.HasPrefix(filepath.Base(currentPath), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	// Handle symlinks
	if info.Mode()&os.ModeSymlink != 0 {
		if !options.FollowSymlinks {
			return visitor(currentPath, info, nil)
		}

		// Resolve symlink and get real info
		realPath, err := filepath.EvalSymlinks(currentPath)
		if err != nil {
			if options.IgnorePermErrors {
				return nil
			}
			return visitor(currentPath, info, err)
		}

		// Get info from real path
		realInfo, err := os.Stat(realPath)
		if err != nil {
			if options.IgnorePermErrors && os.IsPermission(err) {
				return nil
			}
			return visitor(currentPath, info, err)
		}
		info = realInfo
	}

	// Call visitor for current path
	if err := visitor(currentPath, info, nil); err != nil {
		if err == filepath.SkipDir && info.IsDir() {
			return nil
		}
		return err
	}

	// If it's a directory, recurse into it
	if info.IsDir() {
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			if options.IgnorePermErrors && os.IsPermission(err) {
				return nil
			}
			return visitor(currentPath, info, err)
		}

		for _, entry := range entries {
			childPath := filepath.Join(currentPath, entry.Name())
			if err := dt.walkRecursive(root, childPath, currentDepth+1, options, exclusionMatcher, visitor); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListFiles returns a list of files in a directory
func (dt *DefaultTraverser) ListFiles(root string, recursive bool) ([]string, error) {
	// ⭐ EXTRACT-006: File listing implementation - 🔍
	return dt.ListFilesWithExclusions(root, nil, recursive)
}

// ListFilesWithExclusions returns a list of files with exclusion patterns applied
func (dt *DefaultTraverser) ListFilesWithExclusions(root string, excludePatterns []string, recursive bool) ([]string, error) {
	// ⭐ EXTRACT-006: File listing with exclusions implementation - 🔍
	var files []string

	options := TraversalOptions{
		ExcludePatterns:  excludePatterns,
		FollowSymlinks:   false,
		IgnorePermErrors: true,
	}

	if !recursive {
		options.MaxDepth = 0
	}

	visitor := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore errors in listing mode
		}

		// Only include regular files
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	}

	err := dt.WalkWithOptions(root, options, visitor)
	return files, err
}

// Convenience functions

// Walk performs a basic directory walk using the default traverser
func Walk(root string, visitor FileVisitor) error {
	// ⭐ EXTRACT-006: Convenience walk function - 🔍
	traverser := NewTraverser()
	return traverser.Walk(root, visitor)
}

// WalkWithExclusions performs directory walk with exclusion patterns
func WalkWithExclusions(root string, excludePatterns []string, visitor FileVisitor) error {
	// ⭐ EXTRACT-006: Convenience walk with exclusions function - 🔍
	traverser := NewTraverser()
	options := TraversalOptions{
		ExcludePatterns:  excludePatterns,
		FollowSymlinks:   false,
		IgnorePermErrors: true,
	}
	return traverser.WalkWithOptions(root, options, visitor)
}

// ListFiles returns all files in a directory
func ListFiles(root string, recursive bool) ([]string, error) {
	// ⭐ EXTRACT-006: Convenience file listing function - 🔍
	traverser := NewTraverser()
	return traverser.ListFiles(root, recursive)
}

// ListFilesWithExclusions returns files with exclusion patterns applied
func ListFilesWithExclusions(root string, excludePatterns []string, recursive bool) ([]string, error) {
	// ⭐ EXTRACT-006: Convenience file listing with exclusions function - 🔍
	traverser := NewTraverser()
	return traverser.ListFilesWithExclusions(root, excludePatterns, recursive)
}
