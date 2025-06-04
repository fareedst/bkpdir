// Package fileops provides file operations and utilities for CLI applications.
//
// This file contains file exclusion functionality with pattern matching support.
package fileops

import (
	"path/filepath"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

// ⭐ EXTRACT-006: File exclusion system extracted - 🔧

// PatternMatcher handles matching paths against exclusion patterns
type PatternMatcher struct {
	patterns []string
}

// Excluder defines the interface for file exclusion operations
type Excluder interface {
	ShouldExclude(path string) bool
}

// NewPatternMatcher creates a new PatternMatcher with the given patterns
func NewPatternMatcher(patterns []string) *PatternMatcher {
	// ⭐ EXTRACT-006: Pattern matcher initialization extracted - 🔧
	return &PatternMatcher{patterns: patterns}
}

// ShouldExclude checks if a path should be excluded based on patterns
func (pm *PatternMatcher) ShouldExclude(path string) bool {
	// ⭐ EXTRACT-006: File exclusion logic implementation extracted - 🔍
	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range pm.patterns {
		if pm.matchesPattern(normalizedPath, pattern) {
			return true
		}
	}
	return false
}

// matchesPattern checks if a path matches a single pattern
func (pm *PatternMatcher) matchesPattern(path, pattern string) bool {
	// ⭐ EXTRACT-006: Pattern matching logic extracted - 🔧
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
	// ⭐ EXTRACT-006: Directory pattern matching extracted - 🔧
	patternsToTry := []string{
		pattern,                // e.g., node_modules/
		pattern + "**",         // e.g., node_modules/**
		"**/" + pattern,        // e.g., **/node_modules/
		"**/" + pattern + "**", // e.g., **/node_modules/**
	}

	for _, p := range patternsToTry {
		matched, err := doublestar.Match(p, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// matchesGlobPattern handles patterns containing *
func (pm *PatternMatcher) matchesGlobPattern(path, pattern string) bool {
	// ⭐ EXTRACT-006: Glob pattern matching extracted - 🔧
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
	// ⭐ EXTRACT-006: Root-level pattern matching extracted - 🔧
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
	// ⭐ EXTRACT-006: Directory-level pattern matching extracted - 🔧
	if len(pathParts) == len(patternParts) {
		matched, err := doublestar.Match(pattern, path)
		return err == nil && matched
	}
	return false
}

// ShouldExcludeFile checks if a file should be excluded based on patterns
// This is the main public interface for file exclusion
func ShouldExcludeFile(path string, patterns []string) bool {
	// ⭐ EXTRACT-006: Public file exclusion interface extracted - 🔍
	matcher := NewPatternMatcher(patterns)
	return matcher.ShouldExclude(path)
}
