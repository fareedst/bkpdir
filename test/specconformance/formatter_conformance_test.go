// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-LIST_FORMAT_SAFETY] [IMPL-FILE_STATISTICS]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestListFormatSafe_Oracle(t *testing.T) {
	if !specmodel.OracleListFormatSafe("file %s") {
		t.Fatal("expected safe")
	}
	if specmodel.OracleListFormatSafe("file %") {
		t.Fatal("expected unsafe bare percent")
	}
}

func TestFormatterSidecarsFormal(t *testing.T) {
	root := repoRoot(t)
	for _, impl := range []string{
		"IMPL-LIST_FORMAT_SAFETY",
		"IMPL-FILE_STATISTICS",
		"IMPL-FILE_STATISTICS_TEMPLATE_FIX",
		"IMPL-DELAYED_OUTPUT",
		"IMPL-DUAL_FORMATTING",
		"IMPL-CUSTOMIZABLE_FORMAT_STRINGS",
		"IMPL-LIST_LIMIT",
		"IMPL-TEST_COVERAGE",
	} {
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
