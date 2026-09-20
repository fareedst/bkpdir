// [IMPL-CONTEXT_OPS] [IMPL-RESOURCE_MANAGER] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// ContextPhase models resource context lifecycle from sidecar POST conditions.
type ContextPhase int

const (
	ContextOpen ContextPhase = iota
	ContextActive
	ContextClosed
)

// OracleContextTransition validates allowed context phase changes.
func OracleContextTransition(from, to ContextPhase) bool {
	switch from {
	case ContextOpen:
		return to == ContextActive || to == ContextClosed
	case ContextActive:
		return to == ContextClosed
	case ContextClosed:
		return false
	default:
		return false
	}
}
