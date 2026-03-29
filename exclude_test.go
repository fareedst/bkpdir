// [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION]
// Tests for file exclusion logic: validates pattern matching and exclusion behavior
// for directory patterns, glob patterns, and exact matches.
package main

import "testing"

// [REQ-CONFIGURATION] [ARCH-EXCLUSION_PATTERNS] [IMPL-EXCLUSION_PATTERNS]
// TestShouldExcludeFile validates pattern-based file exclusion
func TestShouldExcludeFile_REQ_CONFIGURATION(t *testing.T) {
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

// [REQ-CONFIGURATION] [ARCH-EXCLUSION_PATTERNS] [IMPL-EXCLUSION_PATTERNS]
// TestSingleSegmentDirPatternAnyDepth verifies names like node_modules/ match at any tree depth (no **/ required).
func TestSingleSegmentDirPatternAnyDepth_REQ_CONFIGURATION(t *testing.T) {
	patterns := []string{"node_modules/"}
	cases := []struct {
		relPath  string
		excluded bool
		name     string
	}{
		{"node_modules/lib/foo.js", true, "root node_modules"},
		{"packages/a/node_modules/pkg/x.js", true, "nested node_modules"},
		{"deep/nested/node_modules", true, "node_modules directory path"},
		{"src/not_node_modules/file.go", false, "path segment must match exactly"},
		{"node_modules_extra/file.go", false, "prefix of segment must not match"},
	}
	for _, tc := range cases {
		if got := ShouldExcludeFile(tc.relPath, patterns); got != tc.excluded {
			t.Errorf("%s: ShouldExcludeFile(%q) = %v, want %v", tc.name, tc.relPath, got, tc.excluded)
		}
	}
}

// [REQ-CONFIGURATION] [ARCH-EXCLUSION_PATTERNS] [IMPL-EXCLUSION_PATTERNS]
// TestDirectoryExclusionPattern validates directory exclusion patterns ending with /
func TestDirectoryExclusionPattern_REQ_CONFIGURATION(t *testing.T) {
	patterns := []string{"demo/batches/", "*.log"}
	tests := []struct {
		file     string
		excluded bool
		name     string
	}{
		{"demo/batches/file1.txt", true, "file in demo/batches/ directory"},
		{"demo/batches/subdir/file2.txt", true, "file in subdirectory of demo/batches/"},
		{"demo/batches", false, "demo/batches as file (not directory)"},
		{"demo/other/file.txt", false, "file in demo/other (not excluded)"},
		{"other/batches/file.txt", false, "batches in different parent"},
		{"error.log", true, "log file excluded"},
		{"demo/batches.log", true, "log file matching *.log"},
	}
	for _, tt := range tests {
		if got := ShouldExcludeFile(tt.file, patterns); got != tt.excluded {
			t.Errorf("%s: ShouldExcludeFile(%q) = %v, want %v", tt.name, tt.file, got, tt.excluded)
		}
	}
}
