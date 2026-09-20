// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestSpecctl_OracleValidateAtomicSidecar(t *testing.T) {
	root := repoRoot(t)
	errs, warns, err := specmodel.OracleSpecctlValidatePath(root, "tied/implementation-decisions/IMPL-ATOMIC_OPS-pseudocode.md")
	if err != nil {
		t.Fatal(err)
	}
	if errs != 0 || warns != 0 {
		t.Fatalf("errs=%d warns=%d", errs, warns)
	}
}

// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: parse each file and fail on PARSE/SHAPE/RESOLVE errors; optional CRIT-001 via run_impl_logic_audit.py --req-criteria-strict; L4 mutation pilot and contract eval documented in scripts (not global CI).
// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: build REQ/IMPL/STEP to test/code matrix; L4 mutation score on pilot packages via go-mutesting; property conformance tests in test/specconformance.
// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: verify [PROC-IMPL_PSEUDOCODE_TOKENS] literal copy in Go sources.
func TestSpecctlSidecar_IMPL_SPEC_CTL(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-SPEC_CTL")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}
