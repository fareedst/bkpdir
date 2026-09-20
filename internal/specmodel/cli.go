// [IMPL-CLI_FRAMEWORK] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// OracleDryRunSkipsWrite models dry-run POST: no filesystem mutation.
func OracleDryRunSkipsWrite(dryRun bool) bool {
	return dryRun
}

// OracleCommandRequiresArgs reports whether a command needs positional args per sidecar.
func OracleCommandRequiresArgs(command string) bool {
	switch command {
	case "backup", "archive", "diff":
		return true
	default:
		return false
	}
}
