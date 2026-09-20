// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-FILE_OPERATIONS] [IMPL-DIFF_COMMAND] [IMPL-DIRECTORY_COMPARISON] [IMPL-EXCLUSION_PATTERNS]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
	"bkpdir/pkg/fileops"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestFileopsValidatePath_OracleMatchesPkg(t *testing.T) {
	v := fileops.NewPathValidator()
	cases := []struct {
		path string
		ok   bool
	}{
		{"", false},
		{"/tmp/safe.txt", true},
		{"/tmp/../etc/passwd", false},
	}
	for _, c := range cases {
		errPkg := v.ValidatePath(c.path)
		errOracle := specmodel.OracleValidatePath(c.path)
		if (errPkg == nil) != c.ok || (errOracle == nil) != c.ok {
			t.Fatalf("path %q pkg=%v oracle=%v want ok=%v", c.path, errPkg, errOracle, c.ok)
		}
	}
}

func TestFileopsSidecarFormal_IMPL_FILE_OPERATIONS(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-FILE_OPERATIONS")
	if err != nil {
		t.Fatal(err)
	}
	if missing := specmodel.BlocksMissingLead(doc); len(missing) > 0 {
		t.Fatalf("missing leads: %v", missing)
	}
}
