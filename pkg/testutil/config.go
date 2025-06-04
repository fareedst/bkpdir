// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// DefaultConfigProvider provides standard configuration testing functionality.
type DefaultConfigProvider struct {
	data map[string]interface{}
}

// NewConfigProvider creates a new configuration provider with initial data.
//
// ⭐ EXTRACT-009: Configuration provider creation - 🔧
func NewConfigProvider(data map[string]interface{}) ConfigProvider {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &DefaultConfigProvider{
		data: data,
	}
}

// GetConfigData returns the configuration data as a map.
//
// ⭐ EXTRACT-009: Configuration data access - 🔧
func (cp *DefaultConfigProvider) GetConfigData() map[string]interface{} {
	return cp.data
}

// SetConfigValue sets a configuration value.
//
// ⭐ EXTRACT-009: Configuration value setting - 🔧
func (cp *DefaultConfigProvider) SetConfigValue(key string, value interface{}) {
	cp.data[key] = value
}

// SaveToFile saves the configuration to a YAML file.
//
// ⭐ EXTRACT-009: Configuration file saving - 🔧
func (cp *DefaultConfigProvider) SaveToFile(path string) error {
	yamlData, err := yaml.Marshal(cp.data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, yamlData, 0644)
}

// DefaultEnvironmentManager provides environment variable management for tests.
type DefaultEnvironmentManager struct {
	originalValues map[string]string
	unsetKeys      []string
}

// NewEnvironmentManager creates a new environment manager.
//
// ⭐ EXTRACT-009: Environment manager creation - 🔧
func NewEnvironmentManager() EnvironmentManager {
	return &DefaultEnvironmentManager{
		originalValues: make(map[string]string),
		unsetKeys:      make([]string, 0),
	}
}

// SetEnv sets an environment variable for the test.
// Extracted from TEST-FIX-001 environment variable isolation patterns.
//
// ⭐ EXTRACT-009: Environment variable setting - 🔧
func (em *DefaultEnvironmentManager) SetEnv(key, value string) {
	// Store original value if not already stored
	if _, exists := em.originalValues[key]; !exists {
		if originalValue, wasSet := os.LookupEnv(key); wasSet {
			em.originalValues[key] = originalValue
		} else {
			em.unsetKeys = append(em.unsetKeys, key)
		}
	}

	os.Setenv(key, value)
}

// UnsetEnv removes an environment variable.
//
// ⭐ EXTRACT-009: Environment variable unsetting - 🔧
func (em *DefaultEnvironmentManager) UnsetEnv(key string) {
	// Store original value if not already stored
	if _, exists := em.originalValues[key]; !exists {
		if originalValue, wasSet := os.LookupEnv(key); wasSet {
			em.originalValues[key] = originalValue
		}
	}

	os.Unsetenv(key)
}

// GetEnv gets an environment variable value.
//
// ⭐ EXTRACT-009: Environment variable getting - 🔧
func (em *DefaultEnvironmentManager) GetEnv(key string) string {
	return os.Getenv(key)
}

// RestoreEnv restores all environment variables to their original state.
//
// ⭐ EXTRACT-009: Environment variable restoration - 🔧
func (em *DefaultEnvironmentManager) RestoreEnv() {
	// Restore original values
	for key, value := range em.originalValues {
		os.Setenv(key, value)
	}

	// Unset keys that were originally unset
	for _, key := range em.unsetKeys {
		os.Unsetenv(key)
	}

	// Clear tracking
	em.originalValues = make(map[string]string)
	em.unsetKeys = make([]string, 0)
}

// CreateTestConfig creates a test configuration file with the given data.
// Extracted from config_test.go createTestConfigFile and createTestConfigFileWithData functions.
//
// ⭐ EXTRACT-009: Test configuration creation - 🔧
func CreateTestConfig(t *testing.T, dir string, data map[string]interface{}) (string, func()) {
	t.Helper()

	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Create config provider and save to file
	provider := NewConfigProvider(data)
	if err := provider.SaveToFile(configPath); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Create environment manager for isolation
	envManager := NewEnvironmentManager()

	// Set environment variable to use only this config file
	envManager.SetEnv("BKPDIR_CONFIG", configPath)

	// Return cleanup function
	cleanup := func() {
		envManager.RestoreEnv()
		os.Remove(configPath)
	}

	return configPath, cleanup
}

// CreateTestConfigWithDefaults creates a test configuration with common default values.
// This provides a convenient way to create test configs with typical settings.
//
// ⭐ EXTRACT-009: Default test configuration - 🔧
func CreateTestConfigWithDefaults(t *testing.T, dir string, overrides map[string]interface{}) (string, func()) {
	t.Helper()

	// Default configuration values extracted from test patterns
	defaults := map[string]interface{}{
		"archive_dir_path":               "../.bkpdir",
		"use_current_dir_name":           true,
		"include_git_info":               false,
		"backup_dir_path":                "../.bkpdir/files",
		"use_current_dir_name_for_files": false,
		"exclude_patterns":               []string{".git/", "vendor/"},
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	return CreateTestConfig(t, dir, defaults)
}

// WithTestConfig executes a function with a test configuration and cleans up automatically.
//
// ⭐ EXTRACT-009: Test configuration with cleanup - 🔧
func WithTestConfig(t *testing.T, data map[string]interface{}, fn func(configPath string)) {
	t.Helper()

	WithTempDir(t, "testutil-config-", func(dir string) {
		configPath, cleanup := CreateTestConfig(t, dir, data)
		defer cleanup()
		fn(configPath)
	})
}

// WithTestConfigDefaults executes a function with default test configuration.
//
// ⭐ EXTRACT-009: Default test configuration with cleanup - 🔧
func WithTestConfigDefaults(t *testing.T, overrides map[string]interface{}, fn func(configPath string)) {
	t.Helper()

	WithTempDir(t, "testutil-config-", func(dir string) {
		configPath, cleanup := CreateTestConfigWithDefaults(t, dir, overrides)
		defer cleanup()
		fn(configPath)
	})
}

// IsolateEnvironment creates an environment manager and returns a cleanup function.
// This is useful for tests that need to modify environment variables.
//
// ⭐ EXTRACT-009: Environment isolation - 🔧
func IsolateEnvironment(t *testing.T) func() {
	t.Helper()

	envManager := NewEnvironmentManager()

	// Return cleanup function that restores environment
	cleanup := func() {
		envManager.RestoreEnv()
	}

	// Register cleanup with test
	t.Cleanup(cleanup)

	return cleanup
}

// WithEnvironment executes a function with modified environment variables and restores them.
//
// ⭐ EXTRACT-009: Environment modification with cleanup - 🔧
func WithEnvironment(t *testing.T, envVars map[string]string, fn func()) {
	t.Helper()

	envManager := NewEnvironmentManager()
	defer envManager.RestoreEnv()

	// Set all environment variables
	for key, value := range envVars {
		envManager.SetEnv(key, value)
	}

	// Execute function
	fn()
}

// Package-level convenience functions

// CreateSimpleTestConfig creates a simple test configuration with minimal data.
//
// ⭐ EXTRACT-009: Simple test configuration - 🔧
func CreateSimpleTestConfig(t *testing.T, dir string) (string, func()) {
	return CreateTestConfigWithDefaults(t, dir, nil)
}
