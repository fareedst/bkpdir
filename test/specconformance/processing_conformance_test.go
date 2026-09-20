// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-PROCESSING_PATTERNS]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestProcessingPipeline_Oracle(t *testing.T) {
	if !specmodel.OraclePipelineStageOrder([]int{1, 2, 3}) {
		t.Fatal("valid order")
	}
	if specmodel.OraclePipelineStageOrder([]int{2, 1}) {
		t.Fatal("invalid order")
	}
}

func TestProcessingSidecar_IMPL_PROCESSING_PATTERNS(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-PROCESSING_PATTERNS")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}
