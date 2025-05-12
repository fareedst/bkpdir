package main

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// ShouldExcludeFile returns true if the file matches any of the exclude patterns.
func ShouldExcludeFile(filePath string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(filePath, pattern) {
				return true
			}
			continue
		}
		matched, err := doublestar.PathMatch(pattern, filePath)
		if err == nil && matched {
			return true
		}
		// Also try matching against the base name for patterns like 'vendor/'
		if matched, err := doublestar.PathMatch(pattern, filepath.Base(filePath)); err == nil && matched {
			return true
		}
	}
	return false
}
