// [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION]
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
// [IMPL-EXCLUSION_PATTERNS] PatternMatcher — holds exclusion patterns for dispatch-based matching
type PatternMatcher struct {
	patterns []string
}

// Excluder defines the interface for file exclusion operations
// [IMPL-EXCLUSION_PATTERNS] Excluder — contract for path exclusion decision
type Excluder interface {
	ShouldExclude(path string) bool
}

// NewPatternMatcher creates a new PatternMatcher with the given patterns
// [IMPL-EXCLUSION_PATTERNS] NewPatternMatcher — constructs matcher from pattern slice
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return &PatternMatcher{patterns: patterns}
}

// ShouldExclude checks if a path should be excluded based on patterns
// [IMPL-EXCLUSION_PATTERNS] ShouldExclude — normalizes path, iterates patterns, returns on first match
func (pm *PatternMatcher) ShouldExclude(path string) bool {
	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range pm.patterns {
		if pm.matchesPattern(normalizedPath, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a path matches a single pattern
// [IMPL-EXCLUSION_PATTERNS] matchesPattern — dispatch: directory, glob, or exact match
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
// [IMPL-EXCLUSION_PATTERNS] matchesDirectoryPattern — trailing-/ patterns with **/ recursive support
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
// [IMPL-EXCLUSION_PATTERNS] Used for trailing-slash patterns that name a single directory component.
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
// [IMPL-EXCLUSION_PATTERNS] matchesGlobPattern — delegates ** to doublestar, dispatches single-* by depth
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
// [IMPL-EXCLUSION_PATTERNS] matchesRootLevelPattern — tries pattern as-is and with **/ prefix
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
// [IMPL-EXCLUSION_PATTERNS] matchesDirectoryLevelPattern — multi-segment glob, matches if part counts equal
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
// [IMPL-EXCLUSION_PATTERNS] ShouldExcludeFile — convenience: creates matcher and checks single path
func ShouldExcludeFile(path string, patterns []string) bool {
	matcher := NewPatternMatcher(patterns)
	return matcher.ShouldExclude(path)
}
