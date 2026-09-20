// [IMPL-STRUCTURED_ERRORS] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"errors"
	"strings"
)

// OracleIsDiskFull mirrors pkg/errors disk-full classification semantics.
func OracleIsDiskFull(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	for _, p := range []string{
		"no space left on device", "disk full", "not enough space",
		"quota exceeded", "filesystem full",
	} {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}

// OracleIsPermission mirrors permission error classification.
func OracleIsPermission(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "permission denied") ||
		strings.Contains(text, "access denied") ||
		errors.Is(err, errPermission)
}

var errPermission = errors.New("permission denied")
