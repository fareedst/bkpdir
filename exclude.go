// Package main provides backward-compatible file exclusion wrappers.
// Delegates to pkg/fileops for pattern matching implementation.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"bkpdir/pkg/fileops"
)

// Legacy type aliases for backward compatibility
type PatternMatcher = fileops.PatternMatcher

// Legacy function wrappers for backward compatibility

// NewPatternMatcher creates a new PatternMatcher — delegates to fileops.NewPatternMatcher
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: store pattern slice on PatternMatcher for iterative ShouldExclude checks.
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return fileops.NewPatternMatcher(patterns)
}

// ShouldExcludeFile checks if a file should be excluded — delegates to fileops.ShouldExcludeFile
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: convenience wrapper constructing matcher and calling ShouldExclude once.
func ShouldExcludeFile(path string, patterns []string) bool {
	return fileops.ShouldExcludeFile(path, patterns)
}
