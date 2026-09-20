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
func TestOracleRegistry_48Verified_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	root := repoRoot(t)
	reg, err := specmodel.LoadOracleRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	if reg.OracleTargets != 48 {
		t.Fatalf("oracle_targets=%d want 48", reg.OracleTargets)
	}
	if len(reg.Oracles) != 48 {
		t.Fatalf("oracles=%d want 48", len(reg.Oracles))
	}
	if n := reg.VerifiedCount(); n != 48 {
		t.Fatalf("verified=%d want 48", n)
	}
}
