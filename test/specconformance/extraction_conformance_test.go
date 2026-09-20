// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-PACKAGE_EXTRACTION]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestExtraction_OraclePkgSet(t *testing.T) {
	if !specmodel.OraclePkgExtracted("pkg/fileops") {
		t.Fatal("fileops expected")
	}
	if specmodel.OraclePkgExtracted("pkg/nope") {
		t.Fatal("unexpected pkg")
	}
}

func TestExtractionSidecar_IMPL_PACKAGE_EXTRACTION(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-PACKAGE_EXTRACTION")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}
