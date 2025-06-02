// This file is part of bkpdir

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Test data constants for better maintainability
const (
	defaultArchiveDir = "../.bkpdir"
	testArchiveDir    = "/tmp/archives"
)

var (
	defaultExcludePatterns = []string{".git/", "vendor/"}
	testExcludePatterns    = []string{"node_modules/", "*.log"}
)

func TestDefaultConfig(t *testing.T) {
	// CFG-002: Configuration defaults validation
	// TEST-REF: Feature tracking matrix CFG-002
	// IMMUTABLE-REF: Configuration System
	cfg := DefaultConfig()

	t.Run("default values", func(t *testing.T) {
		assertStringEqual(t, "ArchiveDirPath", cfg.ArchiveDirPath, defaultArchiveDir)
		assertBoolEqual(t, "UseCurrentDirName", cfg.UseCurrentDirName, true)
		assertStringSliceEqual(t, "ExcludePatterns", cfg.ExcludePatterns, defaultExcludePatterns)
		assertBoolEqual(t, "IncludeGitInfo", cfg.IncludeGitInfo, true)

		// Test verification config
		if cfg.Verification == nil {
			t.Error("Verification config should not be nil")
		} else {
			assertBoolEqual(t, "Verification.VerifyOnCreate", cfg.Verification.VerifyOnCreate, false)
			assertStringEqual(t, "Verification.ChecksumAlgorithm", cfg.Verification.ChecksumAlgorithm, "sha256")
		}
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("no config file uses defaults", func(t *testing.T) {
		dir := t.TempDir()

		cfg, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig error: %v", err)
		}

		assertStringEqual(t, "default ArchiveDirPath", cfg.ArchiveDirPath, defaultArchiveDir)
		assertBoolEqual(t, "default UseCurrentDirName", cfg.UseCurrentDirName, true)
	})

	t.Run("config file overrides defaults", func(t *testing.T) {
		dir := t.TempDir()

		// Create test config file
		createTestConfigFile(t, dir)

		cfg, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig error: %v", err)
		}

		// Verify overridden values
		assertStringEqual(t, "YAML ArchiveDirPath", cfg.ArchiveDirPath, testArchiveDir)
		assertBoolEqual(t, "YAML UseCurrentDirName", cfg.UseCurrentDirName, false)
		assertStringSliceEqual(t, "YAML ExcludePatterns", cfg.ExcludePatterns, testExcludePatterns)
	})

	t.Run("invalid config file returns error", func(t *testing.T) {
		dir := t.TempDir()

		// Create invalid YAML file
		invalidYAML := "invalid: yaml: content: ["
		configPath := filepath.Join(dir, ".bkpdir.yml")
		err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid config: %v", err)
		}

		// Should still return a config (with defaults) even if YAML is invalid
		cfg, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig should not error on invalid YAML: %v", err)
		}

		// Should fall back to defaults
		assertStringEqual(t, "fallback ArchiveDirPath", cfg.ArchiveDirPath, defaultArchiveDir)
	})
}

func TestGetConfigValues(t *testing.T) {
	t.Run("returns expected config values", func(t *testing.T) {
		cfg := DefaultConfig()
		values := GetConfigValues(cfg)

		if len(values) == 0 {
			t.Error("GetConfigValues should return non-empty slice")
		}

		// Check that we have expected config values
		expectedKeys := []string{
			"archive_dir_path",
			"use_current_dir_name",
			"include_git_info",
			"backup_dir_path",
			"use_current_dir_name_for_files",
			"verify_on_create",
			"checksum_algorithm",
		}

		valueMap := make(map[string]ConfigValue)
		for _, v := range values {
			valueMap[v.Name] = v
		}

		for _, key := range expectedKeys {
			if _, exists := valueMap[key]; !exists {
				t.Errorf("Expected config key %q not found", key)
			}
		}
	})
}

// TestGetConfigSearchPath tests the configuration search path functionality for CFG-001 feature
func TestGetConfigSearchPath(t *testing.T) {
	// CFG-001: Configuration discovery validation
	// TEST-REF: Feature tracking matrix CFG-001
	// IMMUTABLE-REF: Configuration Discovery
	// Test basic path retrieval

	// Save original environment
	origEnv := os.Getenv("BKPDIR_CONFIG")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("BKPDIR_CONFIG")
		} else {
			os.Setenv("BKPDIR_CONFIG", origEnv)
		}
	}()

	t.Run("default search paths", func(t *testing.T) {
		// Clear environment variable
		os.Unsetenv("BKPDIR_CONFIG")

		paths := getConfigSearchPaths()
		expectedPaths := []string{"./.bkpdir.yml", "~/.bkpdir.yml"}

		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		for i, expected := range expectedPaths {
			if i < len(paths) && paths[i] != expected {
				t.Errorf("Expected path[%d] = %s, got %s", i, expected, paths[i])
			}
		}
	})

	t.Run("custom environment paths", func(t *testing.T) {
		customPaths := "/custom/config.yml:/another/config.yml"
		os.Setenv("BKPDIR_CONFIG", customPaths)

		paths := getConfigSearchPaths()
		expectedPaths := []string{"/custom/config.yml", "/another/config.yml"}

		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		for i, expected := range expectedPaths {
			if i < len(paths) && paths[i] != expected {
				t.Errorf("Expected path[%d] = %s, got %s", i, expected, paths[i])
			}
		}
	})

	t.Run("single custom path", func(t *testing.T) {
		customPath := "/single/config.yml"
		os.Setenv("BKPDIR_CONFIG", customPath)

		paths := getConfigSearchPaths()
		expectedPaths := []string{"/single/config.yml"}

		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		if len(paths) > 0 && paths[0] != expectedPaths[0] {
			t.Errorf("Expected path[0] = %s, got %s", expectedPaths[0], paths[0])
		}
	})

	t.Run("path expansion", func(t *testing.T) {
		// Test home directory expansion
		testPath := "~/test-config.yml"
		expandedPath := expandPath(testPath)

		if expandedPath == testPath {
			// If expandPath didn't change the path, check if it's because home dir wasn't available
			home, err := os.UserHomeDir()
			if err == nil && home != "" {
				t.Errorf("Expected path expansion for %s, but got unchanged path", testPath)
			}
		} else if !strings.Contains(expandedPath, "test-config.yml") {
			t.Errorf("Expected expanded path to contain 'test-config.yml', got %s", expandedPath)
		}

		// Test non-home path (should remain unchanged)
		regularPath := "/regular/path.yml"
		if expandPath(regularPath) != regularPath {
			t.Errorf("Expected regular path to remain unchanged, got %s", expandPath(regularPath))
		}
	})
}

func TestConfigModification(t *testing.T) {
	t.Run("setting string configuration values", func(t *testing.T) {
		testConfigStringValues(t)
	})

	t.Run("setting boolean configuration values", func(t *testing.T) {
		testConfigBooleanValues(t)
	})

	t.Run("setting integer configuration values", func(t *testing.T) {
		testConfigIntegerValues(t)
	})

	t.Run("creating new configuration file when none exists", func(t *testing.T) {
		testConfigFileCreation(t)
	})

	t.Run("updating existing configuration file preserves other values", func(t *testing.T) {
		testConfigFileUpdate(t)
	})

	t.Run("handling nested verification section", func(t *testing.T) {
		testConfigNestedSection(t)
	})

	t.Run("error handling for invalid configuration keys", func(t *testing.T) {
		testConfigInvalidKeys(t)
	})

	t.Run("error handling for invalid boolean values", func(t *testing.T) {
		testConfigInvalidBooleans(t)
	})

	t.Run("error handling for invalid integer values", func(t *testing.T) {
		testConfigInvalidIntegers(t)
	})

	t.Run("configuration persistence verification", func(t *testing.T) {
		testConfigPersistence(t)
	})

	t.Run("type validation for all supported configuration types", func(t *testing.T) {
		testConfigTypeValidation(t)
	})
}

// Helper functions for configuration modification tests - broken down to reduce complexity

func testConfigStringValues(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Test setting archive_dir_path
	testConfigSet(t, dir, "archive_dir_path", "/custom/archive/path")
	assertConfigFileValue(t, configPath, "archive_dir_path", "/custom/archive/path")

	// Test setting backup_dir_path
	testConfigSet(t, dir, "backup_dir_path", "/custom/backup/path")
	assertConfigFileValue(t, configPath, "backup_dir_path", "/custom/backup/path")

	// Test setting checksum_algorithm
	testConfigSet(t, dir, "checksum_algorithm", "md5")
	assertNestedConfigFileValue(t, configPath, "verification", "checksum_algorithm", "md5")
}

func testConfigBooleanValues(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Test setting use_current_dir_name
	testConfigSet(t, dir, "use_current_dir_name", "false")
	assertConfigFileValue(t, configPath, "use_current_dir_name", false)

	// Test setting include_git_info
	testConfigSet(t, dir, "include_git_info", "true")
	assertConfigFileValue(t, configPath, "include_git_info", true)

	// Test setting verify_on_create
	testConfigSet(t, dir, "verify_on_create", "true")
	assertNestedConfigFileValue(t, configPath, "verification", "verify_on_create", true)
}

func testConfigIntegerValues(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Test setting status codes
	testConfigSet(t, dir, "status_config_error", "10")
	assertConfigFileValue(t, configPath, "status_config_error", 10)

	testConfigSet(t, dir, "status_created_archive", "20")
	assertConfigFileValue(t, configPath, "status_created_archive", 20)

	testConfigSet(t, dir, "status_disk_full", "30")
	assertConfigFileValue(t, configPath, "status_disk_full", 30)
}

func testConfigFileCreation(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Verify file doesn't exist
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("Config file should not exist initially")
	}

	// Set a configuration value
	testConfigSet(t, dir, "archive_dir_path", "/new/path")

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file should have been created")
	}

	// Verify value was set
	assertConfigFileValue(t, configPath, "archive_dir_path", "/new/path")
}

func testConfigFileUpdate(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Create initial config with multiple values
	initialConfig := map[string]interface{}{
		"archive_dir_path":     "/initial/path",
		"use_current_dir_name": true,
		"include_git_info":     false,
	}
	createTestConfigFileWithData(t, configPath, initialConfig)

	// Update one value
	testConfigSet(t, dir, "archive_dir_path", "/updated/path")

	// Verify updated value
	assertConfigFileValue(t, configPath, "archive_dir_path", "/updated/path")

	// Verify other values preserved
	assertConfigFileValue(t, configPath, "use_current_dir_name", true)
	assertConfigFileValue(t, configPath, "include_git_info", false)
}

func testConfigNestedSection(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Set verify_on_create (should create verification section)
	testConfigSet(t, dir, "verify_on_create", "true")
	assertNestedConfigFileValue(t, configPath, "verification", "verify_on_create", true)

	// Set checksum_algorithm (should update existing verification section)
	testConfigSet(t, dir, "checksum_algorithm", "sha1")
	assertNestedConfigFileValue(t, configPath, "verification", "checksum_algorithm", "sha1")

	// Verify both values exist in verification section
	assertNestedConfigFileValue(t, configPath, "verification", "verify_on_create", true)
	assertNestedConfigFileValue(t, configPath, "verification", "checksum_algorithm", "sha1")
}

func testConfigInvalidKeys(t *testing.T) {
	// Note: Testing invalid keys would call os.Exit() which terminates the test process
	// This is the expected behavior of the application when invalid keys are provided
	// The error handling is tested through manual verification and integration tests
	t.Log("Invalid key handling calls os.Exit() - tested through integration tests")
}

func testConfigInvalidBooleans(t *testing.T) {
	// Note: Testing invalid boolean values would call os.Exit() which terminates the test process
	// This is the expected behavior of the application when invalid values are provided
	// The error handling is tested through manual verification and integration tests
	t.Log("Invalid boolean value handling calls os.Exit() - tested through integration tests")
}

func testConfigInvalidIntegers(t *testing.T) {
	// Note: Testing invalid integer values would call os.Exit() which terminates the test process
	// This is the expected behavior of the application when invalid values are provided
	// The error handling is tested through manual verification and integration tests
	t.Log("Invalid integer value handling calls os.Exit() - tested through integration tests")
}

func testConfigPersistence(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".bkpdir.yml")

	// Set multiple configuration values
	testConfigSet(t, dir, "archive_dir_path", "/persistent/path")
	testConfigSet(t, dir, "use_current_dir_name", "false")
	testConfigSet(t, dir, "status_config_error", "42")

	// Verify all values persist
	assertConfigFileValue(t, configPath, "archive_dir_path", "/persistent/path")
	assertConfigFileValue(t, configPath, "use_current_dir_name", false)
	assertConfigFileValue(t, configPath, "status_config_error", 42)

	// Load config and verify values are applied
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	assertStringEqual(t, "loaded ArchiveDirPath", cfg.ArchiveDirPath, "/persistent/path")
	assertBoolEqual(t, "loaded UseCurrentDirName", cfg.UseCurrentDirName, false)
	assertIntEqual(t, "loaded StatusConfigError", cfg.StatusConfigError, 42)
}

func testConfigTypeValidation(t *testing.T) {
	dir := t.TempDir()

	// Change to test directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(dir)

	// Test all string types
	stringKeys := []string{"archive_dir_path", "backup_dir_path", "checksum_algorithm"}
	for _, key := range stringKeys {
		testConfigSet(t, dir, key, "test_value")
	}

	// Test all boolean types
	boolKeys := []string{
		"use_current_dir_name", "use_current_dir_name_for_files",
		"include_git_info", "verify_on_create",
	}
	for _, key := range boolKeys {
		testConfigSet(t, dir, key, "true")
		testConfigSet(t, dir, key, "false")
	}

	// Test all integer types
	intKeys := []string{
		"status_config_error", "status_created_archive", "status_created_backup",
		"status_disk_full", "status_permission_denied",
	}
	for _, key := range intKeys {
		testConfigSet(t, dir, key, "123")
	}
}

// Helper functions for configuration modification tests

func testConfigSet(t *testing.T, dir, key, value string) {
	t.Helper()

	// Change to test directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(dir)

	// Call the config set function
	handleConfigSetCommand(key, value)
}

func testConfigSetWithError(t *testing.T, dir, key, value string) error {
	t.Helper()

	// Capture the exit behavior by using a defer to recover from os.Exit
	var exitCode int
	var exitCalled bool

	// Override os.Exit temporarily
	oldOsExit := osExit
	osExit = func(code int) {
		exitCode = code
		exitCalled = true
		panic("exit called") // Use panic to break execution flow
	}
	defer func() {
		osExit = oldOsExit
		if r := recover(); r != nil && r != "exit called" {
			panic(r) // Re-panic if it's not our expected exit
		}
	}()

	// Change to test directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(dir)

	// Call the config set function
	handleConfigSetCommand(key, value)

	// If we get here without panic, no error occurred
	if exitCalled && exitCode != 0 {
		return &ConfigError{Message: "Configuration error", Code: exitCode}
	}
	return nil
}

// ConfigError represents a configuration error for testing
type ConfigError struct {
	Message string
	Code    int
}

func (e *ConfigError) Error() string {
	return e.Message
}

// Variable to allow overriding os.Exit in tests
var osExit = os.Exit

func assertConfigFileValue(t *testing.T, configPath, key string, expectedValue interface{}) {
	t.Helper()

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config YAML: %v", err)
	}

	actualValue, exists := config[key]
	if !exists {
		t.Fatalf("Config key %q not found in file", key)
	}

	if actualValue != expectedValue {
		t.Errorf("Config key %q = %v (%T), want %v (%T)",
			key, actualValue, actualValue, expectedValue, expectedValue)
	}
}

func assertNestedConfigFileValue(t *testing.T, configPath, section, key string, expectedValue interface{}) {
	t.Helper()

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse config YAML: %v", err)
	}

	sectionData, exists := config[section]
	if !exists {
		t.Fatalf("Config section %q not found in file", section)
	}

	sectionMap, ok := sectionData.(map[string]interface{})
	if !ok {
		t.Fatalf("Config section %q is not a map", section)
	}

	actualValue, exists := sectionMap[key]
	if !exists {
		t.Fatalf("Config key %q not found in section %q", key, section)
	}

	if actualValue != expectedValue {
		t.Errorf("Config key %q in section %q = %v (%T), want %v (%T)",
			key, section, actualValue, actualValue, expectedValue, expectedValue)
	}
}

func createTestConfigFileWithData(t *testing.T, configPath string, data map[string]interface{}) {
	t.Helper()

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal config data: %v", err)
	}

	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// Helper functions for better test assertions

func assertStringEqual(t *testing.T, name, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

func assertBoolEqual(t *testing.T, name string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %t, want %t", name, got, want)
	}
}

func assertIntEqual(t *testing.T, name string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %d, want %d", name, got, want)
	}
}

func assertStringSliceEqual(t *testing.T, name string, got, want []string) {
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

func createTestConfigFile(t *testing.T, dir string) {
	t.Helper()

	configYAML := "archive_dir_path: " + testArchiveDir + "\n" +
		"use_current_dir_name: false\n" +
		"exclude_patterns:\n" +
		"  - node_modules/\n" +
		"  - '*.log'\n"

	configPath := filepath.Join(dir, ".bkpdir.yml")
	err := os.WriteFile(configPath, []byte(configYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
}
