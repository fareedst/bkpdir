// [IMPL-LIST_FORMAT_SAFETY] [IMPL-FILE_STATISTICS] [IMPL-DELAYED_OUTPUT] [IMPL-DUAL_FORMATTING]
// [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// OracleListFormatSafe reports whether a format string is safe for list output (no bare %).
func OracleListFormatSafe(format string) bool {
	if format == "" {
		return true
	}
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		if i+1 >= len(format) {
			return false
		}
		next := format[i+1]
		if next == '%' {
			i++
			continue
		}
		if next == 's' || next == 'd' || next == 'v' {
			return true
		}
		return false
	}
	return true
}

// OracleDelayedOutputReady models delayed output collector state POST.
func OracleDelayedOutputReady(buffered int, flushed bool) bool {
	if flushed {
		return buffered == 0
	}
	return buffered >= 0
}

// OracleTruncateList applies list limit semantics from IMPL-LIST_LIMIT.
func OracleTruncateList(items []string, limit int) []string {
	if limit <= 0 || len(items) <= limit {
		return append([]string(nil), items...)
	}
	return append([]string(nil), items[:limit]...)
}

// OracleDualFormatPrimary returns whether AI format takes precedence when both set.
func OracleDualFormatPrimary(useAI, useLegacy bool) string {
	if useAI && !useLegacy {
		return "ai"
	}
	if useLegacy {
		return "legacy"
	}
	return "default"
}
