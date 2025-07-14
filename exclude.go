// This file is part of bkpdir
//
// Package main provides file exclusion functionality for BkpDir.
// It handles pattern matching for excluding files and directories from archives.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// EXCLUDE-001: File Exclusion Requirements Immutable - File exclusion and pattern matching [ACTION:core-functionality]
// Source: docs/context/immutable.md - File Exclusion Requirements section
// Impact: Core functionality requirement for file exclusion

// EXCLUSION-FEATURES-001: Exclusion features specification - File exclusion and pattern matching [ACTION:core-functionality]
// Source: exclude.go - EXCLUSION-FEATURES-001
// Impact: Core functionality requirement for exclusion features

// SERVICE-EXCLUSION-001: Exclusion service architecture decision - Exclusion service implementation [ACTION:core-functionality]
// Source: exclude.go - SERVICE-EXCLUSION-001
// Impact: Exclusion service implementation decision
package main

import (
	"bkpdir/pkg/fileops"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]

// Legacy type aliases for backward compatibility
type PatternMatcher = fileops.PatternMatcher

// Legacy function wrappers for backward compatibility

// NewPatternMatcher creates a new PatternMatcher with the given patterns using the extracted package
func NewPatternMatcher(patterns []string) *PatternMatcher {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	return fileops.NewPatternMatcher(patterns)
}

// ShouldExcludeFile checks if a file should be excluded based on patterns using the extracted package
func ShouldExcludeFile(path string, patterns []string) bool {
	// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
	return fileops.ShouldExcludeFile(path, patterns)
}
