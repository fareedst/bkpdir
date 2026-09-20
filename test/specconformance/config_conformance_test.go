// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-CFG_006] [IMPL-CFG_INHERITANCE_PATH_RESOLUTION]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: table-driven test asserts determineMergeStrategyForField paths map to expected strategies.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestConfigMergePrepend_Oracle(t *testing.T) {
	base := []string{"a", "b"}
	over := []string{"x"}
	got := specmodel.OracleMergeLists(specmodel.MergePrepend, base, over)
	want := []string{"x", "a", "b"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestConfigExtractStrategy_QuotedKey(t *testing.T) {
	if got := specmodel.OracleExtractStrategy(`"exclude":prepend`); got != specmodel.MergeReplace {
		t.Fatalf("quoted key got %v", got)
	}
	if got := specmodel.OracleExtractStrategy("exclude:prepend"); got != specmodel.MergePrepend {
		t.Fatalf("prepend got %v", got)
	}
}

func TestConfigSidecarsFormal(t *testing.T) {
	root := repoRoot(t)
	for _, impl := range []string{
		"IMPL-CFG_006",
		"IMPL-CFG_INHERITANCE_PATH_RESOLUTION",
		"IMPL-CFG_MERGE_BEHAVIOR_REGISTRY",
		"IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE",
		"IMPL-CFG_PRECEDENCE_FIX",
		"IMPL-CONFIG_OUTPUT_GROUPING",
		"IMPL-CONFIG_SCHEMA_FLEX",
		"IMPL-CONFIG_STRUCT",
		"IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX",
		"IMPL-CFG_MIXED_MODE_MERGE_FIX",
		"IMPL-CFG_QUOTED_KEY_PREFIX",
		"IMPL-EXCLUDE_MERGE_FIX",
		"IMPL-TEST_CFG_005_P1",
		"IMPL-TEST_DEFAULT_STRATEGY_EDGES",
		"IMPL-TEST_EMPTY_STRING_HANDLING",
		"IMPL-TEST_EXCLUDE_MERGE",
		"IMPL-TEST_PREPEND_ORDERING",
		"IMPL-TEST_UNICODE_HANDLING",
	} {
		t.Run(impl, func(t *testing.T) {
			doc, err := specmodel.SidecarForImpl(root, impl)
			if err != nil {
				t.Fatal(err)
			}
			if !specmodel.OracleSidecarFormalReady(doc) {
				t.Fatal("sidecar not formal-ready")
			}
		})
	}
}
