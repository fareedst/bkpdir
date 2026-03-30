// [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY]
// Composition tests: production CLI wiring (newRootCommand + executeWithAutoDetection)
// without duplicating a partial Cobra tree. See tied/docs/composition-coverage.md.
package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func sortedSubcommandNames(root *cobra.Command) []string {
	var names []string
	for _, c := range root.Commands() {
		names = append(names, c.Name())
	}
	sort.Strings(names)
	return names
}

// [REQ-USABILITY] [IMPL-AUTO_DETECTION] Production root registers all subcommands (no drift vs createTestRootCmd).
func TestComposition_NewRootCommand_Subcommands_REQ_USABILITY(t *testing.T) {
	root := newRootCommand()
	got := sortedSubcommandNames(root)
	want := []string{"backup", "config", "create", "diff", "full", "inc", "list", "template", "version"}
	if !slices.Equal(got, want) {
		t.Errorf("subcommand names\ngot  %q\nwant %q", got, want)
	}
}

// [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] First token is a known subcommand → Cobra Execute (not path auto-detect).
func TestComposition_ExecuteWithAutoDetection_KnownSubcommandHelp_REQ_USABILITY(t *testing.T) {
	root := newRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	if err := executeWithAutoDetection(root, []string{"template", "--help"}); err != nil {
		t.Fatalf("executeWithAutoDetection: %v", err)
	}
	if !strings.Contains(out.String(), "template") && !strings.Contains(out.String(), "Generate") {
		t.Errorf("expected template help output, got: %q", out.String())
	}
}

// [IMPL-AUTO_DETECTION] [ARCH-AUTO_DETECTION] [REQ-USABILITY] First token is a filesystem path → auto-detect path (dry-run file backup).
func TestComposition_ExecuteWithAutoDetection_PathArg_FileDryRun_REQ_USABILITY(t *testing.T) {
	origDry, origNote := dryRun, note
	defer func() {
		dryRun, note = origDry, origNote
	}()

	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "f.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(tmp, ".bkpdir.yml")
	content := "archive_dir_path: ./ar\nbackup_dir_path: ./bk\nuse_current_dir_name: true\nuse_current_dir_name_for_files: true\n"
	if err := os.WriteFile(cfg, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "ar"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "bk"), 0755); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	if err := executeWithAutoDetection(root, []string{"f.txt", "--dry-run"}); err != nil {
		t.Fatalf("executeWithAutoDetection: %v", err)
	}
}

// [REQ-CONFIGURATION] [IMPL-AUTO_DETECTION] Root persistent --config goes through root Run and handleConfigCommand.
func TestComposition_RootPersistentConfigFlag_REQ_CONFIGURATION(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".bkpdir.yml")
	if err := os.WriteFile(cfgPath, []byte("archive_dir_path: ./archives\n"), 0644); err != nil {
		t.Fatal(err)
	}
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	root := newRootCommand()
	root.SetOut(w)
	root.SetErr(w)
	done := make(chan struct{})
	var buf bytes.Buffer
	go func() {
		_, _ = io.Copy(&buf, r)
		close(done)
	}()

	errExec := executeWithAutoDetection(root, []string{"--config"})
	_ = w.Close()
	os.Stdout = oldStdout
	<-done
	_ = r.Close()

	if errExec != nil {
		t.Fatalf("executeWithAutoDetection: %v", errExec)
	}
	out := buf.String()
	if !strings.Contains(out, "archive_dir_path") {
		t.Errorf("expected config output to mention archive_dir_path, got: %q", out)
	}
}

// [REQ-LIST_LIMIT] [ARCH-LIST_LIMIT] [IMPL-LIST_LIMIT] Root --list with --limit uses handleListFileBackupsCommand (no backups OK).
func TestComposition_RootListAndLimit_REQ_LIST_LIMIT(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "listed.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(tmp, ".bkpdir.yml")
	if err := os.WriteFile(cfg, []byte("backup_dir_path: ./backups\nuse_current_dir_name_for_files: false\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "backups"), 0755); err != nil {
		t.Fatal(err)
	}
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	if err := executeWithAutoDetection(root, []string{"--list", "listed.txt", "--limit", "2"}); err != nil {
		t.Fatalf("executeWithAutoDetection: %v", err)
	}
	// No backups: should print a friendly message, not crash
	if out.Len() == 0 {
		t.Log("no stdout captured (formatter may use os.Stdout); still OK if no exit")
	}
}

// [REQ-DIFF_COMMAND] [ARCH-DIFF_COMMAND] [IMPL-AUTO_DETECTION] diff runs via production root + executeWithAutoDetection (requires archives; avoids os.Exit(0) no-archive path).
func TestComposition_DiffViaExecuteWithAutoDetection_REQ_DIFF_COMMAND(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	archiveDir := filepath.Join(tmpDir, "archives")
	if err := os.MkdirAll(sourceDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"file1.txt": "content1"}
	for rel, c := range files {
		p := filepath.Join(sourceDir, rel)
		if err := os.WriteFile(p, []byte(c), 0644); err != nil {
			t.Fatal(err)
		}
	}
	fullArchivePath := filepath.Join(archiveDir, "bkp-2024-01-01T120000.zip")
	if err := createTestArchive(fullArchivePath, sourceDir, files); err != nil {
		t.Fatal(err)
	}
	cfgContent := "archive_dir_path: " + archiveDir + "\nuse_current_dir_name: false\n"
	if err := os.WriteFile(filepath.Join(sourceDir, ".bkpdir.yml"), []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}

	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)
	if err := os.Chdir(sourceDir); err != nil {
		t.Fatal(err)
	}

	root := newRootCommand()
	var stdout, stderr bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stderr)
	if err := executeWithAutoDetection(root, []string{"diff"}); err != nil {
		t.Fatalf("executeWithAutoDetection: %v", err)
	}
	combined := stdout.String() + stderr.String()
	if strings.Contains(combined, "Error calculating diff") {
		t.Errorf("unexpected diff error output: %q", combined)
	}
}
