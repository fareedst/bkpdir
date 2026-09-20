// [IMPL-PROCESSING_PATTERNS] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// OraclePipelineStageOrder validates stage indices are monotonic.
func OraclePipelineStageOrder(stages []int) bool {
	for i := 1; i < len(stages); i++ {
		if stages[i] < stages[i-1] {
			return false
		}
	}
	return true
}

// OracleConcurrentWorkers caps worker count per sidecar DATA contract.
func OracleConcurrentWorkers(requested, maxAllowed int) int {
	if requested <= 0 {
		return 1
	}
	if requested > maxAllowed {
		return maxAllowed
	}
	return requested
}
