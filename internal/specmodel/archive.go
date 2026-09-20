// [IMPL-ZIP_FORMAT] [IMPL-INCREMENTAL_DUPLICATE_PREVENTION] [IMPL-DATA_MODELS]
// [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"fmt"
	"strings"
)

// OracleArchiveName builds incremental archive basename per sidecar naming rules.
func OracleArchiveName(prefix, timestamp, ext string, incremental bool) string {
	base := strings.TrimSuffix(prefix, ext)
	if incremental {
		return fmt.Sprintf("%s-inc-%s%s", base, timestamp, ext)
	}
	return fmt.Sprintf("%s-%s%s", base, timestamp, ext)
}

// OracleDuplicatePrevented reports whether an archive name would duplicate an existing set.
func OracleDuplicatePrevented(existing []string, candidate string) bool {
	for _, e := range existing {
		if e == candidate {
			return false
		}
	}
	return true
}

// OracleZipEntryPath normalizes relative paths inside archives.
func OracleZipEntryPath(rel string) string {
	rel = strings.TrimPrefix(rel, "./")
	rel = strings.ReplaceAll(rel, "\\", "/")
	return rel
}
