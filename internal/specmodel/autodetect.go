// [IMPL-AUTO_DETECTION] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// OracleAutoCommand maps detected path kind to subcommand name from sidecar.
func OracleAutoCommand(kind string) string {
	switch kind {
	case "file":
		return "backup"
	case "directory":
		return "archive"
	default:
		return ""
	}
}
