// [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [ARCH-PACKAGE_EXTRACTION] [REQ-CONFIGURATION]
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

// [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION]
// NewPatternMatcher creates a new PatternMatcher — delegates to fileops.NewPatternMatcher
func NewPatternMatcher(patterns []string) *PatternMatcher {
	return fileops.NewPatternMatcher(patterns)
}

// [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION]
// ShouldExcludeFile checks if a file should be excluded — delegates to fileops.ShouldExcludeFile
func ShouldExcludeFile(path string, patterns []string) bool {
	return fileops.ShouldExcludeFile(path, patterns)
}
