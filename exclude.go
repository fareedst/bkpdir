// This file is part of bkpdir
//
// Package main provides file exclusion functionality for BkpDir.
// It handles pattern matching for excluding files and directories from archives.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"bkpdir/pkg/fileops"
)

// ⭐ EXTRACT-006: Updated to use extracted fileops package - 🔧

// Legacy type aliases for backward compatibility
type PatternMatcher = fileops.PatternMatcher

// Legacy function wrappers for backward compatibility

// NewPatternMatcher creates a new PatternMatcher with the given patterns using the extracted package
func NewPatternMatcher(patterns []string) *PatternMatcher {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔧
	return fileops.NewPatternMatcher(patterns)
}

// ShouldExcludeFile checks if a file should be excluded based on patterns using the extracted package
func ShouldExcludeFile(path string, patterns []string) bool {
	// ⭐ EXTRACT-006: Delegating to extracted package - 🔍
	return fileops.ShouldExcludeFile(path, patterns)
}
