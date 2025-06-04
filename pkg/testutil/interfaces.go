// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"testing"

	"github.com/spf13/cobra"
)

// ConfigProvider defines the interface for providing configuration data in tests.
// This allows the testutil package to work with any configuration system.
type ConfigProvider interface {
	// GetConfigData returns the configuration data as a map
	GetConfigData() map[string]interface{}

	// SetConfigValue sets a configuration value
	SetConfigValue(key string, value interface{})

	// SaveToFile saves the configuration to a file
	SaveToFile(path string) error
}

// EnvironmentManager defines the interface for managing environment variables in tests.
type EnvironmentManager interface {
	// SetEnv sets an environment variable for the test
	SetEnv(key, value string)

	// UnsetEnv removes an environment variable
	UnsetEnv(key string)

	// GetEnv gets an environment variable value
	GetEnv(key string) string

	// RestoreEnv restores all environment variables to their original state
	RestoreEnv()
}

// FileSystemTestHelper defines the interface for file system testing utilities.
type FileSystemTestHelper interface {
	// CreateTempDir creates a temporary directory and returns its path
	CreateTempDir(t *testing.T, prefix string) string

	// CreateTempFile creates a temporary file with content and returns its path
	CreateTempFile(t *testing.T, dir, name string, content []byte) string

	// CreateDirectory creates a directory structure from a map of paths to content
	CreateDirectory(t *testing.T, root string, files map[string]string) error

	// CreateZipArchive creates a ZIP archive from a map of paths to content
	CreateZipArchive(t *testing.T, archivePath string, files map[string]string) error

	// CreateTestFiles creates multiple test files in a directory from a map
	CreateTestFiles(t *testing.T, baseDir string, files map[string]string)
}

// CliTestHelper defines the interface for CLI testing utilities.
type CliTestHelper interface {
	// CreateTestCommand creates a test cobra command
	CreateTestCommand(name string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command

	// ExecuteCommand executes a command with arguments and returns output and error
	ExecuteCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error)

	// CaptureOutput captures stdout/stderr during command execution
	CaptureOutput(t *testing.T, fn func()) (stdout, stderr string)
}

// AssertionHelper defines the interface for test assertion utilities.
type AssertionHelper interface {
	// AssertEqual asserts that two values are equal
	AssertEqual(t *testing.T, name string, got, want interface{})

	// AssertStringEqual asserts that two strings are equal
	AssertStringEqual(t *testing.T, name, got, want string)

	// AssertBoolEqual asserts that two booleans are equal
	AssertBoolEqual(t *testing.T, name string, got, want bool)

	// AssertIntEqual asserts that two integers are equal
	AssertIntEqual(t *testing.T, name string, got, want int)

	// AssertSliceEqual asserts that two slices are equal
	AssertSliceEqual(t *testing.T, name string, got, want []string)

	// AssertError asserts that an error occurred
	AssertError(t *testing.T, err error, expectError bool)

	// AssertContains asserts that a string contains a substring
	AssertContains(t *testing.T, str, substr, name string)
}

// TestFixtureManager defines the interface for managing test fixtures and data.
type TestFixtureManager interface {
	// LoadFixture loads test data from a fixture file
	LoadFixture(t *testing.T, name string) (map[string]interface{}, error)

	// SaveFixture saves test data to a fixture file
	SaveFixture(t *testing.T, name string, data map[string]interface{}) error

	// GetTestData returns predefined test data sets
	GetTestData(name string) map[string]interface{}

	// RegisterTestData registers a test data set
	RegisterTestData(name string, data map[string]interface{})
}

// TestScenario defines the interface for complex test scenarios.
type TestScenario interface {
	// Setup prepares the test scenario
	Setup(t *testing.T) error

	// Execute runs the test scenario
	Execute(t *testing.T) error

	// Verify validates the test scenario results
	Verify(t *testing.T) error

	// Cleanup cleans up the test scenario
	Cleanup(t *testing.T) error

	// GetName returns the scenario name
	GetName() string

	// GetDescription returns the scenario description
	GetDescription() string
}

// TestUtilProvider defines the main interface that provides access to all testing utilities.
// This is the primary interface that applications should use to access testutil functionality.
type TestUtilProvider interface {
	// GetConfigProvider returns a configuration provider for testing
	GetConfigProvider() ConfigProvider

	// GetEnvironmentManager returns an environment manager for testing
	GetEnvironmentManager() EnvironmentManager

	// GetFileSystemHelper returns a file system helper for testing
	GetFileSystemHelper() FileSystemTestHelper

	// GetCliHelper returns a CLI helper for testing
	GetCliHelper() CliTestHelper

	// GetAssertionHelper returns an assertion helper for testing
	GetAssertionHelper() AssertionHelper

	// GetFixtureManager returns a fixture manager for testing
	GetFixtureManager() TestFixtureManager

	// CreateScenario creates a new test scenario
	CreateScenario(name, description string) TestScenario
}
