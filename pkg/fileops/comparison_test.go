// [REQ-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [IMPL-DIRECTORY_COMPARISON]
package fileops
import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)
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
func TestCreateDirectorySnapshot_NoExclusions_REQ_DIFF_COMMAND(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("alpha"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("beta"), 0644); err != nil {
		t.Fatal(err)
	}

	snap, err := CreateDirectorySnapshot(dir, nil)
	if err != nil {
		t.Fatalf("CreateDirectorySnapshot: %v", err)
	}
	if len(snap.Files) != 2 {
		t.Fatalf("want 2 files in snapshot, got %d", len(snap.Files))
	}
}

func TestCreateDirectorySnapshot_WithExclusions_REQ_DIFF_COMMAND(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("k"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "drop.log"), []byte("l"), 0644); err != nil {
		t.Fatal(err)
	}

	snap, err := CreateDirectorySnapshot(dir, []string{"*.log"})
	if err != nil {
		t.Fatalf("CreateDirectorySnapshot: %v", err)
	}
	if len(snap.Files) != 1 {
		t.Fatalf("want 1 file after exclusion, got %d (%+v)", len(snap.Files), snap.Files)
	}
	if snap.Files[0].RelativePath != "keep.txt" {
		t.Errorf("want keep.txt, got %q", snap.Files[0].RelativePath)
	}
}

func TestCompareSnapshots_IdenticalAndDifferent_REQ_DIFF_COMMAND(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "x.txt")
	if err := os.WriteFile(p, []byte("same"), 0644); err != nil {
		t.Fatal(err)
	}
	s1, err := CreateDirectorySnapshot(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := CreateDirectorySnapshot(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !CompareSnapshots(s1, s2) {
		t.Error("expected identical snapshots")
	}

	if err := os.WriteFile(p, []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	s3, err := CreateDirectorySnapshot(dir, nil)
	if err != nil {
		t.Fatal(err)
	}
	if CompareSnapshots(s1, s3) {
		t.Error("expected snapshots to differ after content change")
	}
}

func writeTestZip(t *testing.T, zipPath string, entries map[string]string) {
	t.Helper()
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateArchiveSnapshot_AndCompareToDir_REQ_DIFF_COMMAND(t *testing.T) {
	dir := t.TempDir()
	content := "hello-archive"
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(t.TempDir(), "snap.zip")
	writeTestZip(t, zipPath, map[string]string{"f.txt": content})

	archSnap, err := CreateArchiveSnapshot(zipPath)
	if err != nil {
		t.Fatalf("CreateArchiveSnapshot: %v", err)
	}
	dirSnap, err := CreateDirectorySnapshot(dir, nil)
	if err != nil {
		t.Fatalf("CreateDirectorySnapshot: %v", err)
	}
	if !CompareSnapshots(dirSnap, archSnap) {
		t.Error("directory and archive snapshots should match")
	}

	ok, err := IsDirectoryIdenticalToArchive(dir, zipPath, nil)
	if err != nil {
		t.Fatalf("IsDirectoryIdenticalToArchive: %v", err)
	}
	if !ok {
		t.Error("expected directory identical to archive")
	}
}

func TestIsDirectoryIdenticalToArchive_ExtraFileInDir_REQ_DIFF_COMMAND(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "extra.txt"), []byte("y"), 0644); err != nil {
		t.Fatal(err)
	}
	zipPath := filepath.Join(t.TempDir(), "only.zip")
	writeTestZip(t, zipPath, map[string]string{"f.txt": "x"})

	ok, err := IsDirectoryIdenticalToArchive(dir, zipPath, nil)
	if err != nil {
		t.Fatalf("IsDirectoryIdenticalToArchive: %v", err)
	}
	if ok {
		t.Error("expected not identical when directory has extra file")
	}
}
