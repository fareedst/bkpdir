// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-AUTO_DETECTION]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestAutoDetect_CommandMapping(t *testing.T) {
	if cmd := specmodel.OracleAutoCommand("file"); cmd != "backup" {
		t.Fatalf("got %q", cmd)
	}
	if cmd := specmodel.OracleAutoCommand("directory"); cmd != "archive" {
		t.Fatalf("got %q", cmd)
	}
}

func TestAutoDetectSidecar_IMPL_AUTO_DETECTION(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-AUTO_DETECTION")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("sidecar not formal-ready")
	}
}
