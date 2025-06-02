// This file is part of bkpdir

// Package main provides tests for file exclusion logic in BkpDir.
// It verifies pattern matching and exclusion behavior.
package main

import "testing"

// FILE-003: File exclusion pattern validation
// TEST-REF: TestShouldExcludeFile
func TestShouldExcludeFile(t *testing.T) {
	patterns := []string{".git/", "vendor/", "*.tmp", "foo/bar/*", "*.log", "build/*", "**/node_modules/", "*.DS_Store"}
	tests := []struct {
		file     string
		excluded bool
		name     string
	}{
		{".git/config", true, ".git subfile"},
		{"vendor/package/file.go", true, "vendor subfile"},
		{"main.go", false, "normal file"},
		{"temp.tmp", true, "tmp file at root"},
		{"foo/bar/baz.txt", true, "foo/bar/* glob"},
		{"foo/baz.txt", false, "foo/ not excluded"},
		{"error.log", true, "log file at root"},
		{"build/output.o", true, "build/* glob"},
		{"src/build/output.o", false, "nested build/ not excluded"},
		{"project/node_modules/module.js", true, "node_modules in subdir"},
		{".DS_Store", true, "hidden macOS file at root"},
		{"src/.DS_Store", true, "hidden macOS file in subdir"},
		{"foo/.gitignore", false, "dotfile not excluded unless pattern"},
		{"foo/bar/.gitkeep", true, "foo/bar/* glob with dotfile"},
	}
	for _, tt := range tests {
		if got := ShouldExcludeFile(tt.file, patterns); got != tt.excluded {
			t.Errorf("%s: ShouldExcludeFile(%q) = %v, want %v", tt.name, tt.file, got, tt.excluded)
		}
	}
}
