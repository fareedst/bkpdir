// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-CLI_FRAMEWORK]
package specconformance_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCommand sets use/short and WithHandler invokes RunE.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestCLI_DryRunOracle(t *testing.T) {
	if !specmodel.OracleDryRunSkipsWrite(true) {
		t.Fatal("dry run should skip writes")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewRootCommand sets name, short, and non-empty version from AppInfo.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCLIApp preserves AppInfo and AddCommand registers subcommand on root.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify dry-run skips Execute and logs [DRY-RUN] prefix; non-dry-run runs operation.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Create yields active context until cancel; WithTimeout ends after duration.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify AddDryRunFlag and AddNoteFlag register expected flags.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify FormatVersion, CreateVersionTemplate, and CreateVersionCommand outputs.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Execute runs until Cancel then returns canceled.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify WithSignalHandling context is active until manual cancel.
func TestCLISidecar_IMPL_CLI_FRAMEWORK(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-CLI_FRAMEWORK")
	if err != nil {
		t.Fatal(err)
	}
	if !specmodel.OracleSidecarFormalReady(doc) {
		t.Fatal("not formal-ready")
	}
}
