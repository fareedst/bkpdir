// [REQ-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [IMPL-DIRECTORY_COMPARISON]
package fileops

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

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
