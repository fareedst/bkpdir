// [IMPL-FILE_OPERATIONS] [IMPL-AUTO_DETECTION] [IMPL-DIFF_COMMAND] [IMPL-DIRECTORY_COMPARISON] [IMPL-EXCLUSION_PATTERNS]
// [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"strings"
)

// OracleValidatePath models PRE: VALIDATE_PATH(path) == ok from sidecars.
func OracleValidatePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errEmptyPath
	}
	if strings.Contains(path, "..") {
		return errUnsafePath
	}
	return nil
}

var (
	errEmptyPath  = &pathError{"path cannot be empty"}
	errUnsafePath = &pathError{"path contains unsafe elements"}
)

type pathError struct{ msg string }

func (e *pathError) Error() string { return e.msg }

// OraclePathKind classifies paths for auto-detection semantics.
func OraclePathKind(path string) string {
	if OracleValidatePath(path) != nil {
		return "invalid"
	}
	if strings.HasSuffix(path, "/") {
		return "directory"
	}
	if strings.Contains(path, ".") && !strings.HasSuffix(path, ".") {
		base := path[strings.LastIndex(path, "/")+1:]
		if strings.Contains(base, ".") {
			return "file"
		}
	}
	return "unknown"
}
