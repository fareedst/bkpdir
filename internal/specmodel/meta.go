// [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import "bkpdir/internal/specparse"

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
// - [IMPL-TRACEABILITY] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_003] — How: document behavioral contracts in TIED YAML and link each contract to requirement tokens.
// - [IMPL-TRACEABILITY] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_003] — How: build semantic-token dependency graph from registry cross-references for change-impact analysis.
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: generate Markdown registry tables from normalized token artifacts for semantic-tokens.md.
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: replace legacy project token YAML with pointer stubs after migration to tied/semantic-tokens.yaml.
// - [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] — How: aggregate module coverage scans into a cross-module audit report and remediation plan.

// OracleSidecarFormalReady verifies L1/L2 sidecar shape for doc/tooling IMPLs.
func OracleSidecarFormalReady(doc *specparse.Document) bool {
	if doc == nil || len(doc.Blocks) == 0 {
		return false
	}
	for _, b := range doc.Blocks {
		if isExemptBlock(b.Name) {
			continue
		}
		if !BlockHasFormalMeta(b) {
			return false
		}
		if b.Lead == "" {
			return false
		}
	}
	return true
}
