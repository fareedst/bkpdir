// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-TESTING]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestTestutil_Oracle(t *testing.T) {
	if !specmodel.OracleTestFixtureValid("/tmp/fixture") {
		t.Fatal("valid fixture")
	}
	if specmodel.OracleTestFixtureValid("") {
		t.Fatal("empty invalid")
	}
}

func TestTestingSidecar_IMPL_TESTING(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-TESTING")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}

// - [IMPL-TESTING_COMPLEXITY] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: extract shared test utilities first, add package-focused tests, and keep root integration tests green after package extraction.
func TestTestingComplexitySidecarFormal(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-TESTING_COMPLEXITY")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
	if !specmodel.OracleTestFixtureValid("/tmp/bkpdir-fixture") {
		t.Fatal("testutil oracle PRE/POST boundary")
	}
}
