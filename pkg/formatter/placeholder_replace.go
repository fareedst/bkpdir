// [REQ-CUSTOMIZABLE_FORMAT_STRINGS] Simple placeholder replacement
// This package provides a simple, testable function for replacing #{key} placeholders
// Using #{...} instead of %{...} (old syntax) to avoid conflicts with Go's fmt package
package formatter

import (
	"fmt"
	"strings"
)

// ReplacePlaceholders replaces all #{key} patterns in formatStr with values from data map
// Returns the string with all placeholders replaced
// Unknown placeholders are left as-is
// Uses #{...} syntax to avoid fmt package conflicts
func ReplacePlaceholders(formatStr string, data map[string]string) string {
	result := formatStr

	// Replace all #{key} patterns with their values
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
