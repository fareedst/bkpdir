// Package fileops provides file operations and utilities for CLI applications.
//
// This file contains file exclusion functionality with pattern matching support.
package fileops

import (
	"path/filepath"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

// PatternMatcher handles matching paths against exclusion patterns
type PatternMatcher struct {
	patterns []string
}

// Excluder defines the interface for file exclusion operations
type Excluder interface {
	ShouldExclude(path string) bool
}

// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: store pattern slice on PatternMatcher for iterative ShouldExclude checks.
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return &PatternMatcher{patterns: patterns}
}

// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: normalize path to slashes, return true on first pattern match via matchesPattern dispatch.
func (pm *PatternMatcher) ShouldExclude(path string) bool {
	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range pm.patterns {
		if pm.matchesPattern(normalizedPath, pattern) {
			return true
		}
	}
	return false
}

// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: route trailing-/ to directory rules, glob * to doublestar, else exact path equality.
func (pm *PatternMatcher) matchesPattern(path, pattern string) bool {
	normalizedPattern := filepath.ToSlash(pattern)

	if strings.HasSuffix(normalizedPattern, "/") {
		return pm.matchesDirectoryPattern(path, normalizedPattern)
	}

	if strings.Contains(normalizedPattern, "*") {
		return pm.matchesGlobPattern(path, normalizedPattern)
	}

	return path == normalizedPattern
}

// matchesDirectoryPattern handles patterns ending with /
func (pm *PatternMatcher) matchesDirectoryPattern(path, pattern string) bool {
	patternBase := strings.TrimSuffix(pattern, "/")

	// Patterns ending with / should only match directories:
	// - Path must have a / after the pattern base (e.g., "demo/batches/file.txt")
	// - Path must end with / (e.g., "demo/batches/")
	// - Path must NOT be exactly the pattern base without / (e.g., "demo/batches" as a file)

	// Check if pattern starts with **/ (matches anywhere)
	if strings.HasPrefix(pattern, "**/") {
		// For **/ patterns, extract the directory name (after **/)
		// Pattern is like "**/node_modules/", patternBase is "**/node_modules"
		// We need to find "node_modules/" in the path
		dirName := strings.TrimPrefix(patternBase, "**/")
		dirPattern := dirName + "/"

		// Check if path contains dirPattern anywhere
		if strings.Contains(path, dirPattern) {
			// Verify it's not just the directory name as a file
			// Path must have content after dirPattern or end with "/"
			idx := strings.Index(path, dirPattern)
			if idx >= 0 {
				remaining := path[idx+len(dirPattern):]
				// If there's content after dirPattern, it's a directory match
				if len(remaining) > 0 || strings.HasSuffix(path, "/") {
					return true
				}
			}
		}
		return false
	}

	// Single segment (e.g. node_modules/, vendor/): match that directory name at any depth in the tree
	if !strings.Contains(patternBase, "/") {
		return pathHasDirSegment(path, patternBase)
	}

	// Multi-segment non-** patterns: path must be under that directory prefix
	if !strings.HasPrefix(path, patternBase+"/") && path != patternBase+"/" {
		return false
	}
	return true
}

// pathHasDirSegment reports whether path contains a path segment exactly equal to seg (e.g. node_modules).
func pathHasDirSegment(path, seg string) bool {
	if seg == "" {
		return false
	}
	for _, p := range strings.Split(path, "/") {
		if p == seg {
			return true
		}
	}
	return false
}

// matchesGlobPattern handles patterns containing *
func (pm *PatternMatcher) matchesGlobPattern(path, pattern string) bool {
	if strings.Contains(pattern, "**") {
		matched, err := doublestar.Match(pattern, path)
		return err == nil && matched
	}

	pathParts := strings.Split(path, "/")
	patternParts := strings.Split(pattern, "/")

	if len(patternParts) == 1 {
		return pm.matchesRootLevelPattern(path, pattern)
	}

	return pm.matchesDirectoryLevelPattern(path, pattern, pathParts, patternParts)
}

// matchesRootLevelPattern handles patterns without directory parts
func (pm *PatternMatcher) matchesRootLevelPattern(path, pattern string) bool {
	patternsToTry := []string{pattern}
	if !strings.HasPrefix(pattern, "**/") {
		patternsToTry = append(patternsToTry, "**/"+pattern)
	}

	for _, p := range patternsToTry {
		matched, err := doublestar.Match(p, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// matchesDirectoryLevelPattern handles patterns with directory parts
func (pm *PatternMatcher) matchesDirectoryLevelPattern(
	path, pattern string,
	pathParts, patternParts []string,
) bool {
	if len(pathParts) == len(patternParts) {
		matched, err := doublestar.Match(pattern, path)
		return err == nil && matched
	}
	return false
}

// ShouldExcludeFile checks if a file should be excluded based on patterns
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: convenience wrapper constructing matcher and calling ShouldExclude once.
func ShouldExcludeFile(path string, patterns []string) bool {
	matcher := NewPatternMatcher(patterns)
	return matcher.ShouldExclude(path)
}
