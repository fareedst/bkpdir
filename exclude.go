// This file is part of bkpdir
//
// Package main provides file exclusion functionality for BkpDir.
// It handles pattern matching for excluding files and directories from archives.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"path/filepath"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
)

// PatternMatcher handles matching paths against exclusion patterns
type PatternMatcher struct {
	patterns []string
}

// NewPatternMatcher creates a new PatternMatcher with the given patterns
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return &PatternMatcher{patterns: patterns}
}

// ShouldExclude checks if a path should be excluded based on patterns
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
func ShouldExcludeFile(path string, patterns []string) bool {
	matcher := NewPatternMatcher(patterns)
	return matcher.ShouldExclude(path)
}
