// [REQ-INCREMENTAL_ANCHOR_MODE] Full-anchored incremental archive behavior
package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIncrementalArchivePathDisambiguationSameMinute_REQ_INCREMENTAL_ANCHOR_MODE(t *testing.T) {
	tmpDir := t.TempDir()
	candidate := filepath.Join(tmpDir, "proj-2026-01-01-12-00_update=2026-01-01-12-01.zip")
	if err := os.WriteFile(candidate, []byte("existing"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	got := resolveIncrementalArchivePath(candidate)
	want := filepath.Join(tmpDir, "proj-2026-01-01-12-00_update=2026-01-01-12-01-seq2.zip")
	if got != want {
		t.Errorf("resolveIncrementalArchivePath() = %q, want %q", got, want)
	}
}

// TestIncrementalAnchorModeSecondIncCumulativeVsFull replays sponsor log: second inc (full mode) includes all paths changed vs full.
func TestIncrementalAnchorModeSecondIncCumulativeVsFull_REQ_INCREMENTAL_ANCHOR_MODE(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")

	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatalf("Failed to create archive directory: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(sourceDir, "tmp_yaml_err"), 0755); err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "README.md"), []byte("base"), 0644); err != nil {
		t.Fatalf("Failed to create base file: %v", err)
	}

	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}
	defer os.Chdir(originalCwd)

	cfg := DefaultConfig()
	cfg.ArchiveDirPath = archiveDir
	cfg.UseCurrentDirName = false
	cfg.IncrementalAnchorMode = "full"

	if err := CreateFullArchive(cfg, "", false); err != nil {
		t.Fatalf("CreateFullArchive: %v", err)
	}

	if err := os.WriteFile(filepath.Join(sourceDir, "tmp_yaml_err", "1021.1"), []byte("one"), 0644); err != nil {
		t.Fatalf("write 1021.1: %v", err)
	}
	if err := createIncrementalArchive(IncrementalArchiveConfig{
		Config:  cfg,
		Context: context.Background(),
	}); err != nil {
		t.Fatalf("first inc: %v", err)
	}

	if err := os.WriteFile(filepath.Join(sourceDir, "tmp_yaml_err", "1021.2"), []byte("two"), 0644); err != nil {
		t.Fatalf("write 1021.2: %v", err)
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err = createIncrementalArchive(IncrementalArchiveConfig{
		Config:  cfg,
		DryRun:  true,
		Context: context.Background(),
	})
	w.Close()
	os.Stdout = originalStdout

	outBuf := make([]byte, 4096)
	n, _ := r.Read(outBuf)
	output := string(outBuf[:n])

	if err != nil {
		t.Fatalf("second inc dry-run: %v (output: %q)", err, output)
	}

	if !strings.Contains(output, "tmp_yaml_err/1021.1") {
		t.Errorf("expected cumulative file 1021.1 in dry-run output, got: %q", output)
	}
	if !strings.Contains(output, "tmp_yaml_err/1021.2") {
		t.Errorf("expected file 1021.2 in dry-run output, got: %q", output)
	}
}
