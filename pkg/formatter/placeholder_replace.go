// Provides a simple, testable function for replacing #{key} placeholders.
package formatter

import (
	"fmt"
	"strings"
)

// ReplacePlaceholders replaces all #{key} patterns in formatStr with values from data map
// Returns the string with all placeholders replaced
// Unknown placeholders are left as-is
// Uses #{...} syntax to avoid fmt package conflicts
// SPEC-ID: IMPL-DUAL_FORMATTING::SIMPLE_PLACEHOLDER_REPLACE
// - [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING] — How: replace every #{key} in formatStr with data map values; leave unknown placeholders unchanged.
func ReplacePlaceholders(formatStr string, data map[string]string) string {
	result := formatStr

	// Replace all #{key} patterns with their values
	for key, value := range data {
		placeholder := fmt.Sprintf("#{%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}
