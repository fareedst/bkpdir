// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-STRUCTURED_ERRORS]
package specconformance_test

import (
	"errors"
	"testing"

	"bkpdir/internal/specmodel"
	pkgerrors "bkpdir/pkg/errors"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestDiskFull_OracleMatchesPkg(t *testing.T) {
	err := errors.New("no space left on device")
	if !specmodel.OracleIsDiskFull(err) || !pkgerrors.IsDiskFullError(err) {
		t.Fatal("disk full mismatch")
	}
	err2 := errors.New("permission denied")
	if specmodel.OracleIsDiskFull(err2) || pkgerrors.IsDiskFullError(err2) {
		t.Fatal("false positive")
	}
}

func TestStructuredErrorsSidecar_IMPL_STRUCTURED_ERRORS(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-STRUCTURED_ERRORS")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}
