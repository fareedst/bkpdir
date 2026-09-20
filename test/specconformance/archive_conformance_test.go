// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-ZIP_FORMAT] [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [IMPL-DATA_MODELS]
package specconformance_test

import (
	"strings"
	"testing"

	"bkpdir/internal/specmodel"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestArchiveNaming_Oracle(t *testing.T) {
	full := specmodel.OracleArchiveName("proj", "20260101", ".zip", false)
	if !strings.Contains(full, "proj-20260101.zip") {
		t.Fatalf("full name %q", full)
	}
	inc := specmodel.OracleArchiveName("proj", "20260101", ".zip", true)
	if !strings.Contains(inc, "-inc-") {
		t.Fatalf("inc name %q", inc)
	}
}

func TestDuplicatePrevention_Oracle(t *testing.T) {
	if !specmodel.OracleDuplicatePrevented([]string{"a.zip"}, "b.zip") {
		t.Fatal("expected unique candidate allowed")
	}
	if specmodel.OracleDuplicatePrevented([]string{"a.zip"}, "a.zip") {
		t.Fatal("expected duplicate rejected")
	}
}

func TestArchiveSidecarsFormal(t *testing.T) {
	root := repoRoot(t)
	for _, impl := range []string{"IMPL-ZIP_FORMAT", "IMPL-INCREMENTAL_DUPLICATE_PREVENTION", "IMPL-DATA_MODELS"} {
		t.Run(impl, func(t *testing.T) {
			doc, err := specmodel.SidecarForImpl(root, impl)
			if err != nil {
				t.Fatal(err)
			}
			if !specmodel.OracleSidecarFormalReady(doc) {
				t.Fatal("not formal-ready")
			}
		})
	}
}
