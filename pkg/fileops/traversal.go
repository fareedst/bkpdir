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
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with no preloaded exclusion patterns.
func NewTraverser() Traverser {
	return &DefaultTraverser{}
}

// NewTraverserWithPatterns creates a traverser with exclusion patterns
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with PatternMatcher built from initial pattern list.
func NewTraverserWithPatterns(patterns []string) Traverser {
	return &DefaultTraverser{
		matcher: NewPatternMatcher(patterns),
	}
}

// Walk performs a basic directory walk
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: delegate to WalkWithOptions with default options (no exclusions, unlimited depth).
func (dt *DefaultTraverser) Walk(root string, visitor FileVisitor) error {
	options := TraversalOptions{
		FollowSymlinks:   false,
		MaxDepth:         -1,
		IgnoreHidden:     false,
		IgnorePermErrors: false,
	}
	return dt.WalkWithOptions(root, options, visitor)
}

// WalkWithOptions performs directory traversal with custom options
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: validate root exists, build PatternMatcher from options or traverser matcher, invoke walkRecursive.
func (dt *DefaultTraverser) WalkWithOptions(root string, options TraversalOptions, visitor FileVisitor) error {
	if err := ValidatePath(root); err != nil {
		return fmt.Errorf("invalid root path: %v", err)
	}

	if err := ValidateExistence(root); err != nil {
		return err
	}

	var exclusionMatcher *PatternMatcher
	if len(options.ExcludePatterns) > 0 {
		exclusionMatcher = NewPatternMatcher(options.ExcludePatterns)
	} else if dt.matcher != nil {
		exclusionMatcher = dt.matcher
	}

	return dt.walkRecursive(root, root, 0, options, exclusionMatcher, visitor)
}

// walkRecursive performs the actual recursive directory traversal
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: enforce depth limit, apply exclusion and hidden-file rules, optionally follow symlinks, call visitor, recurse into directories.
func (dt *DefaultTraverser) walkRecursive(
	root, currentPath string,
	currentDepth int,
	options TraversalOptions,
	exclusionMatcher *PatternMatcher,
	visitor FileVisitor,
) error {
	if options.MaxDepth >= 0 && currentDepth > options.MaxDepth {
		return nil
	}

	info, err := os.Lstat(currentPath)
	if err != nil {
		if options.IgnorePermErrors && os.IsPermission(err) {
			return nil
		}
		return visitor(currentPath, nil, err)
	}

	relPath, err := filepath.Rel(root, currentPath)
	if err != nil {
		return err
	}

	if relPath == "." {
		relPath = ""
	}

	if exclusionMatcher != nil && relPath != "" && exclusionMatcher.ShouldExclude(relPath) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if options.IgnoreHidden && relPath != "" && strings.HasPrefix(filepath.Base(currentPath), ".") {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	if info.Mode()&os.ModeSymlink != 1 {
		if !options.FollowSymlinks {
			return visitor(currentPath, info, nil)
		}

		realPath, err := filepath.EvalSymlinks(currentPath)
		if err != nil {
			if options.IgnorePermErrors {
				return nil
			}
			return visitor(currentPath, info, err)
		}

		realInfo, err := os.Stat(realPath)
		if err != nil {
			if options.IgnorePermErrors && os.IsPermission(err) {
				return nil
			}
			return visitor(currentPath, info, err)
		}
		info = realInfo
	}

	if err := visitor(currentPath, info, nil); err != nil {
		if err == filepath.SkipDir && info.IsDir() {
			return nil
		}
		return err
	}

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
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: collect regular file paths under root with optional recursion via ListFilesWithExclusions.
func (dt *DefaultTraverser) ListFiles(root string, recursive bool) ([]string, error) {
	return dt.ListFilesWithExclusions(root, nil, recursive)
}

// ListFilesWithExclusions returns a list of files with exclusion patterns applied
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: walk tree applying exclusion patterns and append file paths only.
func (dt *DefaultTraverser) ListFilesWithExclusions(root string, excludePatterns []string, recursive bool) ([]string, error) {
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
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	}

	err := dt.WalkWithOptions(root, options, visitor)
	return files, err
}

// Walk performs directory traversal using the default traverser
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: delegate to WalkWithOptions with default options (no exclusions, unlimited depth).
func Walk(root string, visitor FileVisitor) error {
	traverser := NewTraverser()
	return traverser.Walk(root, visitor)
}

// WalkWithExclusions walks a directory tree applying exclusion patterns
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: WalkWithOptions with ExcludePatterns, no symlink follow, ignore permission errors.
func WalkWithExclusions(root string, excludePatterns []string, visitor FileVisitor) error {
	traverser := NewTraverser()
	options := TraversalOptions{
		ExcludePatterns:  excludePatterns,
		FollowSymlinks:   false,
		IgnorePermErrors: true,
	}
	return traverser.WalkWithOptions(root, options, visitor)
}

// ListFiles lists files under root using the default traverser
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: collect regular file paths under root with optional recursion via ListFilesWithExclusions.
func ListFiles(root string, recursive bool) ([]string, error) {
	traverser := NewTraverser()
	return traverser.ListFiles(root, recursive)
}

// ListFilesWithExclusions lists files under root with exclusion patterns
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: walk tree applying exclusion patterns and append file paths only.
func ListFilesWithExclusions(root string, excludePatterns []string, recursive bool) ([]string, error) {
	traverser := NewTraverser()
	return traverser.ListFilesWithExclusions(root, excludePatterns, recursive)
}
