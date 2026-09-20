// [REQ-CONFIGURATION] [ARCH-EXCLUSION_PATTERNS] [IMPL-EXCLUSION_PATTERNS]
package fileops
import "testing"
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: store pattern slice on PatternMatcher for iterative ShouldExclude checks.
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: normalize path to slashes, return true on first pattern match via matchesPattern dispatch.
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: route trailing-/ to directory rules, glob * to doublestar, else exact path equality.
// - [IMPL-EXCLUSION_PATTERNS] [ARCH-EXCLUSION_PATTERNS] [REQ-CONFIGURATION] — How: convenience wrapper constructing matcher and calling ShouldExclude once.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with no preloaded exclusion patterns.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with PatternMatcher built from initial pattern list.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: delegate to WalkWithOptions with default options (no exclusions, unlimited depth).
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: WalkWithOptions with ExcludePatterns, no symlink follow, ignore permission errors.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: validate root exists, build PatternMatcher from options or traverser matcher, invoke walkRecursive.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: enforce depth limit, apply exclusion and hidden-file rules, optionally follow symlinks, call visitor, recurse into directories.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: collect regular file paths under root with optional recursion via ListFilesWithExclusions.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: walk tree applying exclusion patterns and append file paths only.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject empty paths, run IsSecurePath, clean path without failing on normalization drift.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then Stat path; return not-exist or access errors.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidateExistence then Open path for read; close on success.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then test write via temp file in directory or OpenFile WRONLY for files; recurse to parent when path missing.
// - [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject .. traversal, null bytes, newlines, and shell-expansion characters in path strings.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: validate target path, ensure parent directory exists, create temp file beside target for same-filesystem rename.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: reject writes when writer is closed or temp handle missing; otherwise write bytes to temp file.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: convert text to bytes and delegate to ATOMICWRITER_WRITE.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: close temp handle then ATOMIC_RENAME temp to target; cleanup temp on failure; mark committed and closed.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: forbid rollback after commit; remove temp file via CLEANUP.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: idempotent close; roll back when not committed.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: close open handle and remove temp path from disk ignoring not-exist.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: validate readable source and destination, stream to AtomicWriter, match permissions, commit.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: write byte slice through AtomicWriter, set permissions on temp, commit.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TESTING_STRATEGY] [REQ-RESOURCE_MANAGEMENT] — How: delegate text payload to ATOMICWRITEFILE.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: delegate to context-aware atomic write with background context.
// - [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT] — How: check cancellation before each stage, write temp path, rename to final, untrack temp on success.
func TestShouldExcludeFile_GlobAndExact_REQ_CONFIGURATION(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{name: "glob_star_txt", path: "logs/debug.log", patterns: []string{"*.log"}, want: true},
		{name: "no_match", path: "src/main.go", patterns: []string{"*.log"}, want: false},
		{name: "exact", path: "tmp/cache", patterns: []string{"tmp/cache"}, want: true},
		{name: "empty_patterns", path: "any.txt", patterns: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldExcludeFile(tt.path, tt.patterns); got != tt.want {
				t.Errorf("ShouldExcludeFile(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestPatternMatcher_ShouldExclude_REQ_CONFIGURATION(t *testing.T) {
	pm := NewPatternMatcher([]string{"vendor/", "*.tmp"})
	if !pm.ShouldExclude("vendor/foo/bar.go") {
		t.Error("expected vendor/ to exclude nested path")
	}
	if !pm.ShouldExclude("scratch.tmp") {
		t.Error("expected *.tmp to match")
	}
	if pm.ShouldExclude("cmd/main.go") {
		t.Error("did not expect cmd/main.go to be excluded")
	}
}
