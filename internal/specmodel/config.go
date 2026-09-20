// [IMPL-CFG_*] [IMPL-CONFIG_*] [IMPL-TEST_*] config test IMPLs
// [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// MergeStrategy names config merge behaviors from sidecars.
type MergeStrategy string

const (
	MergeReplace MergeStrategy = "replace"
	MergeAppend  MergeStrategy = "append"
	MergePrepend MergeStrategy = "prepend"
)

// OracleMergeLists applies merge strategy per IMPL-CFG_006 / merge fix sidecars.
func OracleMergeLists(strategy MergeStrategy, base, overlay []string) []string {
	switch strategy {
	case MergePrepend:
		out := make([]string, 0, len(base)+len(overlay))
		out = append(out, overlay...)
		out = append(out, base...)
		return out
	case MergeAppend:
		out := make([]string, 0, len(base)+len(overlay))
		out = append(out, base...)
		out = append(out, overlay...)
		return out
	case MergeReplace:
		if len(overlay) == 0 {
			return append([]string(nil), base...)
		}
		return append([]string(nil), overlay...)
	default:
		return append([]string(nil), base...)
	}
}

// OracleExtractStrategy parses strategy token from config key notation.
func OracleExtractStrategy(key string) MergeStrategy {
	if len(key) > 0 && key[0] == '"' {
		return MergeReplace
	}
	if idx := indexStrategySuffix(key); idx >= 0 {
		switch key[idx:] {
		case ":prepend":
			return MergePrepend
		case ":append":
			return MergeAppend
		}
	}
	return MergeReplace
}

func indexStrategySuffix(key string) int {
	for _, suf := range []string{":prepend", ":append"} {
		if len(key) >= len(suf) && key[len(key)-len(suf):] == suf {
			return len(key) - len(suf)
		}
	}
	return -1
}
