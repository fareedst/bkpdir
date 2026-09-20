// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-GIT_CLI] [IMPL-GIT_DIRTY_CONFIG]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestGitDirty_Oracle(t *testing.T) {
	if !specmodel.OracleGitDirtyAllowed(true, true) {
		t.Fatal("allow dirty")
	}
	if specmodel.OracleGitDirtyAllowed(false, true) {
		t.Fatal("reject dirty")
	}
}

func TestGitSidecarsFormal(t *testing.T) {
	root := repoRoot(t)
	for _, impl := range []string{"IMPL-GIT_CLI", "IMPL-GIT_DIRTY_CONFIG"} {
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
