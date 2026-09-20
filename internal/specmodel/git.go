// [IMPL-GIT_CLI] [IMPL-GIT_DIRTY_CONFIG] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// OracleGitDirtyAllowed models archive-when-dirty config interpretation.
func OracleGitDirtyAllowed(allowDirty, isDirty bool) bool {
	if !isDirty {
		return true
	}
	return allowDirty
}

// OracleGitStatusClean reports working tree clean per sidecar PRE.
func OracleGitStatusClean(modified, untracked int) bool {
	return modified == 0 && untracked == 0
}
