// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-CONTEXT_OPS] [IMPL-RESOURCE_MANAGER]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestContextTransition_Oracle(t *testing.T) {
	if !specmodel.OracleContextTransition(specmodel.ContextOpen, specmodel.ContextActive) {
		t.Fatal("open->active")
	}
	if specmodel.OracleContextTransition(specmodel.ContextClosed, specmodel.ContextActive) {
		t.Fatal("closed->active invalid")
	}
}

func TestResourcesSidecarsFormal(t *testing.T) {
	root := repoRoot(t)
	for _, impl := range []string{"IMPL-CONTEXT_OPS", "IMPL-RESOURCE_MANAGER"} {
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
