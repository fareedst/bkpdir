// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// DefaultAssertionHelper provides standard assertion functionality.
type DefaultAssertionHelper struct{}

// NewAssertionHelper creates a new assertion helper.
func NewAssertionHelper() AssertionHelper {
	return &DefaultAssertionHelper{}
}

// AssertEqual asserts that two values are equal using reflection.
//
// ⭐ EXTRACT-009: Generic equality assertion - 🔧
func (h *DefaultAssertionHelper) AssertEqual(t *testing.T, name string, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

// AssertStringEqual asserts that two strings are equal.
// Extracted from config_test.go assertStringEqual function.
//
// ⭐ EXTRACT-009: String assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertStringEqual(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

// AssertBoolEqual asserts that two booleans are equal.
// Extracted from config_test.go assertBoolEqual function.
//
// ⭐ EXTRACT-009: Boolean assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertBoolEqual(t *testing.T, name string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %t, want %t", name, got, want)
	}
}

// AssertIntEqual asserts that two integers are equal.
// Extracted from config_test.go assertIntEqual function.
//
// ⭐ EXTRACT-009: Integer assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertIntEqual(t *testing.T, name string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %d, want %d", name, got, want)
	}
}

// AssertSliceEqual asserts that two string slices are equal.
// Extracted from config_test.go assertStringSliceEqual function.
//
// ⭐ EXTRACT-009: Slice assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertSliceEqual(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s length = %d, want %d", name, len(got), len(want))
		return
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("%s[%d] = %q, want %q", name, i, v, want[i])
		}
	}
}

// AssertError asserts whether an error should or should not have occurred.
//
// ⭐ EXTRACT-009: Error assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertError(t *testing.T, err error, expectError bool) {
	t.Helper()
	if expectError && err == nil {
		t.Error("Expected error, got nil")
	} else if !expectError && err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// AssertContains asserts that a string contains a substring.
//
// ⭐ EXTRACT-009: String contains assertion helper - 🔧
func (h *DefaultAssertionHelper) AssertContains(t *testing.T, str, substr, name string) {
	t.Helper()
	if !strings.Contains(str, substr) {
		t.Errorf("%s: expected %q to contain %q", name, str, substr)
	}
}

// AssertNotContains asserts that a string does not contain a substring.
// Additional utility not in the interface but useful for testing.
//
// ⭐ EXTRACT-009: String not contains assertion helper - 🔧
func AssertNotContains(t *testing.T, str, substr, name string) {
	t.Helper()
	if strings.Contains(str, substr) {
		t.Errorf("%s: expected %q to not contain %q", name, str, substr)
	}
}

// AssertFileExists asserts that a file exists at the given path.
//
// ⭐ EXTRACT-009: File existence assertion helper - 🔧
func AssertFileExists(t *testing.T, path, name string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("%s: file %q should exist", name, path)
	} else if err != nil {
		t.Errorf("%s: error checking file %q: %v", name, path, err)
	}
}

// AssertFileNotExists asserts that a file does not exist at the given path.
//
// ⭐ EXTRACT-009: File non-existence assertion helper - 🔧
func AssertFileNotExists(t *testing.T, path, name string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err == nil {
			t.Errorf("%s: file %q should not exist", name, path)
		} else {
			t.Errorf("%s: error checking file %q: %v", name, path, err)
		}
	}
}

// AssertFileContent asserts that a file contains the expected content.
//
// ⭐ EXTRACT-009: File content assertion helper - 🔧
func AssertFileContent(t *testing.T, path, expectedContent, name string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("%s: failed to read file %q: %v", name, path, err)
		return
	}

	if string(content) != expectedContent {
		t.Errorf("%s: file content = %q, want %q", name, string(content), expectedContent)
	}
}

// Package-level convenience functions that use the default assertion helper

var defaultHelper = NewAssertionHelper()

// AssertStringEqual is a package-level convenience function for string assertions.
//
// ⭐ EXTRACT-009: Package-level string assertion - 🔧
func AssertStringEqual(t *testing.T, name, got, want string) {
	defaultHelper.AssertStringEqual(t, name, got, want)
}

// AssertBoolEqual is a package-level convenience function for boolean assertions.
//
// ⭐ EXTRACT-009: Package-level boolean assertion - 🔧
func AssertBoolEqual(t *testing.T, name string, got, want bool) {
	defaultHelper.AssertBoolEqual(t, name, got, want)
}

// AssertIntEqual is a package-level convenience function for integer assertions.
//
// ⭐ EXTRACT-009: Package-level integer assertion - 🔧
func AssertIntEqual(t *testing.T, name string, got, want int) {
	defaultHelper.AssertIntEqual(t, name, got, want)
}

// AssertSliceEqual is a package-level convenience function for slice assertions.
//
// ⭐ EXTRACT-009: Package-level slice assertion - 🔧
func AssertSliceEqual(t *testing.T, name string, got, want []string) {
	defaultHelper.AssertSliceEqual(t, name, got, want)
}

// AssertError is a package-level convenience function for error assertions.
//
// ⭐ EXTRACT-009: Package-level error assertion - 🔧
func AssertError(t *testing.T, err error, expectError bool) {
	defaultHelper.AssertError(t, err, expectError)
}

// AssertContains is a package-level convenience function for string contains assertions.
//
// ⭐ EXTRACT-009: Package-level contains assertion - 🔧
func AssertContains(t *testing.T, str, substr, name string) {
	defaultHelper.AssertContains(t, str, substr, name)
}
