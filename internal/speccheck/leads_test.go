// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
package speccheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"bkpdir/internal/speccheck"
)

func TestLeadScanDirPaths_includesStandardRoots(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{"pkg", "internal", "cmd", "test"} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got := speccheck.LeadScanDirPaths(root)
	want := []string{"pkg", "internal", "cmd", "test"}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d: %v", len(got), len(want), got)
	}
	for i, rel := range want {
		if got[i] != filepath.Join(root, rel) {
			t.Fatalf("got[%d]=%q want suffix %q", i, got[i], rel)
		}
	}
}

func TestCheckLeads_findsLeadInInternalOnly(t *testing.T) {
	root := t.TempDir()
	sidecarDir := filepath.Join(root, "tied", "implementation-decisions")
	if err := os.MkdirAll(sidecarDir, 0o755); err != nil {
		t.Fatal(err)
	}
	lead := "[IMPL-SCAN_FIXTURE] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: fixture internal-only lead."
	sidecar := filepath.Join(sidecarDir, "IMPL-SCAN_FIXTURE-pseudocode.md")
	body := "# [IMPL-SCAN_FIXTURE]\n\n## BLOCK\n\nSPEC-ID: IMPL-SCAN_FIXTURE::BLOCK\n\n- " + lead + "\n\nPROCEDURE BLOCK:\n  RETURN ok\n"
	if err := os.WriteFile(sidecar, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	internalGo := filepath.Join(root, "internal", "fixture")
	if err := os.MkdirAll(internalGo, 0o755); err != nil {
		t.Fatal(err)
	}
	goSrc := "// - " + lead + "\npackage fixture\n"
	if err := os.WriteFile(filepath.Join(internalGo, "doc.go"), []byte(goSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}

	narrow := []string{filepath.Join(root, "pkg")}
	diags, err := speccheck.CheckLeads([]string{sidecar}, narrow, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) == 0 {
		t.Fatal("expected TRACE-002 when lead exists only under internal/")
	}

	wide := speccheck.LeadScanDirPaths(root)
	diags, err = speccheck.CheckLeads([]string{sidecar}, wide, root)
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 0 {
		t.Fatalf("expected pass with default scan dirs, got %v", diags)
	}
}
