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
	// 🔺 CFG-002: Configuration defaults validation - 📝
	// TEST-REF: Feature tracking matrix CFG-002
	// IMMUTABLE-REF: Configuration System
	cfg := DefaultConfig()

	t.Run("default values", func(t *testing.T) {
		assertStringEqual(t, "ArchiveDirPath", cfg.ArchiveDirPath, defaultArchiveDir)
		assertBoolEqual(t, "UseCurrentDirName", cfg.UseCurrentDirName, true)
		assertStringSliceEqual(t, "ExcludePatterns", cfg.ExcludePatterns, defaultExcludePatterns)
		assertBoolEqual(t, "IncludeGitInfo", cfg.IncludeGitInfo, false)

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
	// 🔶 TEST-FIX-001: Set BKPDIR_CONFIG to avoid personal config interference - 🔧
	// Save original environment and set to non-existent path to avoid personal config
	origEnv := os.Getenv("BKPDIR_CONFIG")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("BKPDIR_CONFIG")
		} else {
			os.Setenv("BKPDIR_CONFIG", origEnv)
		}
	}()

	t.Run("no config file uses defaults", func(t *testing.T) {
		// Set BKPDIR_CONFIG to non-existent path to ensure only defaults are used
		os.Setenv("BKPDIR_CONFIG", "/nonexistent/path/config.yml")

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

		// Create test config file and set BKPDIR_CONFIG to use only this file
		createTestConfigFile(t, dir)
		configPath := filepath.Join(dir, ".bkpdir.yml")
		os.Setenv("BKPDIR_CONFIG", configPath)

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

		// Create invalid YAML file and set BKPDIR_CONFIG to use only this file
		invalidYAML := "invalid: yaml: content: ["
		configPath := filepath.Join(dir, ".bkpdir.yml")
		err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid config: %v", err)
		}
		os.Setenv("BKPDIR_CONFIG", configPath)

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
	// 🔺 CFG-001: Configuration discovery validation - 🔍
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

	t.Run("empty environment variable", func(t *testing.T) {
		os.Setenv("BKPDIR_CONFIG", "")

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
}

// 🔺 CFG-001: Test GetConfigValuesWithSources for comprehensive configuration value extraction - 🔍
func TestGetConfigValuesWithSources(t *testing.T) {
	// 🔶 TEST-FIX-001: Set BKPDIR_CONFIG to avoid personal config interference - 🔧
	// Save original environment and set to non-existent path to avoid personal config
	origEnv := os.Getenv("BKPDIR_CONFIG")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("BKPDIR_CONFIG")
		} else {
			os.Setenv("BKPDIR_CONFIG", origEnv)
		}
	}()

	t.Run("default configuration sources", func(t *testing.T) {
		// Set BKPDIR_CONFIG to non-existent path to ensure only defaults are used
		os.Setenv("BKPDIR_CONFIG", "/nonexistent/path/config.yml")

		dir := t.TempDir()
		cfg := DefaultConfig()

		values := GetConfigValuesWithSources(cfg, dir)

		if len(values) == 0 {
			t.Error("GetConfigValuesWithSources should return non-empty slice")
		}

		// Verify that all values are marked as default
		for _, value := range values {
			if value.Source != "default" {
				t.Errorf("Expected default source for value %q, got %q", value.Name, value.Source)
			}
		}

		// Check for expected configuration values
		valueMap := make(map[string]ConfigValue)
		for _, v := range values {
			valueMap[v.Name] = v
		}

		// Test basic configuration values
		if val, exists := valueMap["archive_dir_path"]; !exists {
			t.Error("Expected archive_dir_path in config values")
		} else if val.Value != defaultArchiveDir {
			t.Errorf("Expected archive_dir_path value %q, got %q", defaultArchiveDir, val.Value)
		}

		// Test boolean values
		if val, exists := valueMap["use_current_dir_name"]; !exists {
			t.Error("Expected use_current_dir_name in config values")
		} else if val.Value != "true" {
			t.Errorf("Expected use_current_dir_name value 'true', got %q", val.Value)
		}

		// Test verification values
		if val, exists := valueMap["verify_on_create"]; !exists {
			t.Error("Expected verify_on_create in config values")
		} else if val.Value != "false" {
			t.Errorf("Expected verify_on_create value 'false', got %q", val.Value)
		}

		if val, exists := valueMap["checksum_algorithm"]; !exists {
			t.Error("Expected checksum_algorithm in config values")
		} else if val.Value != "sha256" {
			t.Errorf("Expected checksum_algorithm value 'sha256', got %q", val.Value)
		}
	})

	t.Run("configuration with custom file", func(t *testing.T) {
		dir := t.TempDir()

		// Create custom config file and set BKPDIR_CONFIG to use only this file
		configPath := filepath.Join(dir, ".bkpdir.yml")
		configData := map[string]interface{}{
			"archive_dir_path":       "/custom/archives",
			"use_current_dir_name":   false,
			"include_git_info":       true,
			"backup_dir_path":        "/custom/backups",
			"status_created_archive": 42,
			"status_disk_full":       99,
		}
		createTestConfigFileWithData(t, configPath, configData)
		os.Setenv("BKPDIR_CONFIG", configPath)

		// Load config and get values with sources
		cfg, err := LoadConfig(dir)
		if err != nil {
			t.Fatalf("LoadConfig error: %v", err)
		}

		values := GetConfigValuesWithSources(cfg, dir)

		valueMap := make(map[string]ConfigValue)
		for _, v := range values {
			valueMap[v.Name] = v
		}

		// Verify custom values are marked with config file source
		if val, exists := valueMap["archive_dir_path"]; !exists {
			t.Error("Expected archive_dir_path in config values")
		} else {
			if val.Value != "/custom/archives" {
				t.Errorf("Expected archive_dir_path value '/custom/archives', got %q", val.Value)
			}
			if val.Source == "default" {
				t.Error("Expected archive_dir_path source to be config file, got 'default'")
			}
		}

		if val, exists := valueMap["use_current_dir_name"]; !exists {
			t.Error("Expected use_current_dir_name in config values")
		} else {
			if val.Value != "false" {
				t.Errorf("Expected use_current_dir_name value 'false', got %q", val.Value)
			}
			if val.Source == "default" {
				t.Error("Expected use_current_dir_name source to be config file, got 'default'")
			}
		}

		// Verify status codes
		if val, exists := valueMap["status_created_archive"]; !exists {
			t.Error("Expected status_created_archive in config values")
		} else {
			if val.Value != "42" {
				t.Errorf("Expected status_created_archive value '42', got %q", val.Value)
			}
			if val.Source == "default" {
				t.Error("Expected status_created_archive source to be config file, got 'default'")
			}
		}
	})

	t.Run("values are sorted alphabetically", func(t *testing.T) {
		// Set BKPDIR_CONFIG to non-existent path to ensure only defaults are used
		os.Setenv("BKPDIR_CONFIG", "/nonexistent/path/config.yml")

		dir := t.TempDir()
		cfg := DefaultConfig()

		values := GetConfigValuesWithSources(cfg, dir)

		// Verify sorting
		for i := 1; i < len(values); i++ {
			if values[i-1].Name >= values[i].Name {
				t.Errorf("Config values not sorted: %q should come before %q", values[i-1].Name, values[i].Name)
			}
		}
	})
}

// 🔺 CFG-001: Test determineConfigSource for config file source detection - 🔍
func TestDetermineConfigSource(t *testing.T) {
	// 🔶 TEST-FIX-001: Set BKPDIR_CONFIG to avoid personal config interference - 🔧
	// Save original environment and set to non-existent path to avoid personal config
	origEnv := os.Getenv("BKPDIR_CONFIG")
	defer func() {
		if origEnv == "" {
			os.Unsetenv("BKPDIR_CONFIG")
		} else {
			os.Setenv("BKPDIR_CONFIG", origEnv)
		}
	}()

	t.Run("no config file returns default", func(t *testing.T) {
		// Set BKPDIR_CONFIG to non-existent path to ensure defaults are used
		os.Setenv("BKPDIR_CONFIG", "/nonexistent/path/config.yml")

		dir := t.TempDir()

		source := determineConfigSource(dir)

		if source != "default" {
			t.Errorf("Expected 'default' source, got %q", source)
		}
	})

	t.Run("existing config file returns path", func(t *testing.T) {
		dir := t.TempDir()

		// Create config file and set BKPDIR_CONFIG to use only this file
		configPath := filepath.Join(dir, ".bkpdir.yml")
		err := os.WriteFile(configPath, []byte("archive_dir_path: /test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create config file: %v", err)
		}
		os.Setenv("BKPDIR_CONFIG", configPath)

		source := determineConfigSource(dir)

		if source != configPath {
			t.Errorf("Expected config path %q, got %q", configPath, source)
		}
	})

	t.Run("home config file detection", func(t *testing.T) {
		dir := t.TempDir()

		// Simulate home directory config
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Skip("Unable to get home directory")
		}

		homeConfigPath := filepath.Join(homeDir, ".bkpdir.yml")
		// Set BKPDIR_CONFIG to check home config path
		os.Setenv("BKPDIR_CONFIG", homeConfigPath)

		// Only test if home config doesn't exist to avoid interfering with real config
		if _, err := os.Stat(homeConfigPath); os.IsNotExist(err) {
			source := determineConfigSource(dir)
			if source != "default" {
				t.Errorf("Expected 'default' source when no config exists, got %q", source)
			}
		}
	})
}

// 🔺 CFG-001: Test createSourceDeterminer for config value source determination - 🔧
func TestCreateSourceDeterminer(t *testing.T) {
	t.Run("string value comparison", func(t *testing.T) {
		determiner := createSourceDeterminer("/test/config.yml")

		// Test matching values
		source := determiner("default_value", "default_value")
		if source != "default" {
			t.Errorf("Expected 'default' for matching strings, got %q", source)
		}

		// Test different values
		source = determiner("custom_value", "default_value")
		if source != "/test/config.yml" {
			t.Errorf("Expected config path for different strings, got %q", source)
		}
	})

	t.Run("boolean value comparison", func(t *testing.T) {
		determiner := createSourceDeterminer("/test/config.yml")

		// Test matching booleans
		source := determiner(true, true)
		if source != "default" {
			t.Errorf("Expected 'default' for matching booleans, got %q", source)
		}

		source = determiner(false, false)
		if source != "default" {
			t.Errorf("Expected 'default' for matching booleans, got %q", source)
		}

		// Test different booleans
		source = determiner(true, false)
		if source != "/test/config.yml" {
			t.Errorf("Expected config path for different booleans, got %q", source)
		}
	})

	t.Run("integer value comparison", func(t *testing.T) {
		determiner := createSourceDeterminer("/test/config.yml")

		// Test matching integers
		source := determiner(42, 42)
		if source != "default" {
			t.Errorf("Expected 'default' for matching integers, got %q", source)
		}

		// Test different integers
		source = determiner(42, 24)
		if source != "/test/config.yml" {
			t.Errorf("Expected config path for different integers, got %q", source)
		}
	})

	t.Run("string slice comparison", func(t *testing.T) {
		determiner := createSourceDeterminer("/test/config.yml")

		// Test matching slices
		slice1 := []string{"a", "b", "c"}
		slice2 := []string{"a", "b", "c"}
		source := determiner(slice1, slice2)
		if source != "default" {
			t.Errorf("Expected 'default' for matching slices, got %q", source)
		}

		// Test different slices
		slice3 := []string{"x", "y", "z"}
		source = determiner(slice1, slice3)
		if source != "/test/config.yml" {
			t.Errorf("Expected config path for different slices, got %q", source)
		}
	})
}

// 🔺 CFG-001: Test helper functions for config value extraction - 🔍
func TestGetBasicConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	defaultCfg := DefaultConfig()

	// Modify some values
	cfg.ArchiveDirPath = "/custom/path"
	cfg.UseCurrentDirName = false

	determiner := createSourceDeterminer("/test/config.yml")
	values := getBasicConfigValues(cfg, defaultCfg, determiner)

	if len(values) == 0 {
		t.Error("getBasicConfigValues should return non-empty slice")
	}

	// Check specific values
	valueMap := make(map[string]ConfigValue)
	for _, v := range values {
		valueMap[v.Name] = v
	}

	// Test modified values
	if val, exists := valueMap["archive_dir_path"]; !exists {
		t.Error("Expected archive_dir_path in basic config values")
	} else {
		if val.Value != "/custom/path" {
			t.Errorf("Expected archive_dir_path value '/custom/path', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}

	if val, exists := valueMap["use_current_dir_name"]; !exists {
		t.Error("Expected use_current_dir_name in basic config values")
	} else {
		if val.Value != "false" {
			t.Errorf("Expected use_current_dir_name value 'false', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}
}

// 🔺 CFG-002: Test status code value extraction - 🔍
func TestGetStatusCodeValues(t *testing.T) {
	cfg := DefaultConfig()
	defaultCfg := DefaultConfig()

	// Modify some status codes
	cfg.StatusCreatedArchive = 100
	cfg.StatusDiskFull = 200

	determiner := createSourceDeterminer("/test/config.yml")
	values := getStatusCodeValues(cfg, defaultCfg, determiner)

	if len(values) == 0 {
		t.Error("getStatusCodeValues should return non-empty slice")
	}

	// Check specific values
	valueMap := make(map[string]ConfigValue)
	for _, v := range values {
		valueMap[v.Name] = v
	}

	// Test modified values
	if val, exists := valueMap["status_created_archive"]; !exists {
		t.Error("Expected status_created_archive in status code values")
	} else {
		if val.Value != "100" {
			t.Errorf("Expected status_created_archive value '100', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}

	if val, exists := valueMap["status_disk_full"]; !exists {
		t.Error("Expected status_disk_full in status code values")
	} else {
		if val.Value != "200" {
			t.Errorf("Expected status_disk_full value '200', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}
}

// 🔺 CFG-001: Test verification value extraction - 🔍
func TestGetVerificationValues(t *testing.T) {
	cfg := DefaultConfig()
	defaultCfg := DefaultConfig()

	// Modify verification values
	cfg.Verification.VerifyOnCreate = true
	cfg.Verification.ChecksumAlgorithm = "md5"

	determiner := createSourceDeterminer("/test/config.yml")
	values := getVerificationValues(cfg, defaultCfg, determiner)

	if len(values) == 0 {
		t.Error("getVerificationValues should return non-empty slice")
	}

	// Check specific values
	valueMap := make(map[string]ConfigValue)
	for _, v := range values {
		valueMap[v.Name] = v
	}

	// Test modified values
	if val, exists := valueMap["verify_on_create"]; !exists {
		t.Error("Expected verify_on_create in verification values")
	} else {
		if val.Value != "true" {
			t.Errorf("Expected verify_on_create value 'true', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}

	if val, exists := valueMap["checksum_algorithm"]; !exists {
		t.Error("Expected checksum_algorithm in verification values")
	} else {
		if val.Value != "md5" {
			t.Errorf("Expected checksum_algorithm value 'md5', got %q", val.Value)
		}
		if val.Source != "/test/config.yml" {
			t.Errorf("Expected config file source, got %q", val.Source)
		}
	}
}

// 🔺 CFG-004: Test mergeExtendedFormatStrings for extended string configuration - 📝
func TestMergeExtendedFormatStrings(t *testing.T) {
	t.Run("merge non-default format strings", func(t *testing.T) {
		dst := DefaultConfig()
		src := DefaultConfig()

		// Modify source config with custom format strings
		src.FormatNoArchivesFound = "Custom: No archives found in %s"
		src.FormatVerificationSuccess = "Custom: Archive %s verified successfully"
		src.FormatConfigurationUpdated = "Custom: Configuration %s updated to %s"
		src.FormatNoBackupsFound = "Custom: No backups found for %s in %s"
		src.FormatBackupCreated = "Custom: Backup created at %s"

		// Keep some as default to test selective merging
		// src.FormatVerificationFailed remains default

		mergeExtendedFormatStrings(dst, src)

		// Verify custom values were merged
		if dst.FormatNoArchivesFound != "Custom: No archives found in %s" {
			t.Errorf("Expected custom FormatNoArchivesFound, got %q", dst.FormatNoArchivesFound)
		}

		if dst.FormatVerificationSuccess != "Custom: Archive %s verified successfully" {
			t.Errorf("Expected custom FormatVerificationSuccess, got %q", dst.FormatVerificationSuccess)
		}

		if dst.FormatConfigurationUpdated != "Custom: Configuration %s updated to %s" {
			t.Errorf("Expected custom FormatConfigurationUpdated, got %q", dst.FormatConfigurationUpdated)
		}

		if dst.FormatNoBackupsFound != "Custom: No backups found for %s in %s" {
			t.Errorf("Expected custom FormatNoBackupsFound, got %q", dst.FormatNoBackupsFound)
		}

		if dst.FormatBackupCreated != "Custom: Backup created at %s" {
			t.Errorf("Expected custom FormatBackupCreated, got %q", dst.FormatBackupCreated)
		}

		// Verify default values were not overridden
		defaultCfg := DefaultConfig()
		if dst.FormatVerificationFailed != defaultCfg.FormatVerificationFailed {
			t.Error("Default FormatVerificationFailed should not be overridden")
		}
	})

	t.Run("all extended format strings covered", func(t *testing.T) {
		dst := DefaultConfig()
		src := DefaultConfig()

		// Set all format strings to custom values
		src.FormatNoArchivesFound = "Custom1"
		src.FormatVerificationFailed = "Custom2"
		src.FormatVerificationSuccess = "Custom3"
		src.FormatVerificationWarning = "Custom4"
		src.FormatConfigurationUpdated = "Custom5"
		src.FormatConfigFilePath = "Custom6"
		src.FormatDryRunFilesHeader = "Custom7"
		src.FormatDryRunFileEntry = "Custom8"
		src.FormatNoFilesModified = "Custom9"
		src.FormatIncrementalCreated = "Custom10"
		src.FormatNoBackupsFound = "Custom11"
		src.FormatBackupWouldCreate = "Custom12"
		src.FormatBackupIdentical = "Custom13"
		src.FormatBackupCreated = "Custom14"

		mergeExtendedFormatStrings(dst, src)

		// Verify all custom values were merged
		if dst.FormatNoArchivesFound != "Custom1" {
			t.Errorf("Expected 'Custom1', got %q", dst.FormatNoArchivesFound)
		}
		if dst.FormatVerificationFailed != "Custom2" {
			t.Errorf("Expected 'Custom2', got %q", dst.FormatVerificationFailed)
		}
		if dst.FormatVerificationSuccess != "Custom3" {
			t.Errorf("Expected 'Custom3', got %q", dst.FormatVerificationSuccess)
		}
		if dst.FormatVerificationWarning != "Custom4" {
			t.Errorf("Expected 'Custom4', got %q", dst.FormatVerificationWarning)
		}
		if dst.FormatConfigurationUpdated != "Custom5" {
			t.Errorf("Expected 'Custom5', got %q", dst.FormatConfigurationUpdated)
		}
		if dst.FormatConfigFilePath != "Custom6" {
			t.Errorf("Expected 'Custom6', got %q", dst.FormatConfigFilePath)
		}
		if dst.FormatDryRunFilesHeader != "Custom7" {
			t.Errorf("Expected 'Custom7', got %q", dst.FormatDryRunFilesHeader)
		}
		if dst.FormatDryRunFileEntry != "Custom8" {
			t.Errorf("Expected 'Custom8', got %q", dst.FormatDryRunFileEntry)
		}
		if dst.FormatNoFilesModified != "Custom9" {
			t.Errorf("Expected 'Custom9', got %q", dst.FormatNoFilesModified)
		}
		if dst.FormatIncrementalCreated != "Custom10" {
			t.Errorf("Expected 'Custom10', got %q", dst.FormatIncrementalCreated)
		}
		if dst.FormatNoBackupsFound != "Custom11" {
			t.Errorf("Expected 'Custom11', got %q", dst.FormatNoBackupsFound)
		}
		if dst.FormatBackupWouldCreate != "Custom12" {
			t.Errorf("Expected 'Custom12', got %q", dst.FormatBackupWouldCreate)
		}
		if dst.FormatBackupIdentical != "Custom13" {
			t.Errorf("Expected 'Custom13', got %q", dst.FormatBackupIdentical)
		}
		if dst.FormatBackupCreated != "Custom14" {
			t.Errorf("Expected 'Custom14', got %q", dst.FormatBackupCreated)
		}
	})
}

// 🔺 CFG-004: Test mergeExtendedTemplates for extended template configuration - 📝
func TestMergeExtendedTemplates(t *testing.T) {
	t.Run("merge non-default template strings", func(t *testing.T) {
		dst := DefaultConfig()
		src := DefaultConfig()

		// Modify source config with custom template strings
		src.TemplateNoArchivesFound = "{{.custom}} template for no archives"
		src.TemplateVerificationSuccess = "{{.archive}} verification success template"
		src.TemplateConfigurationUpdated = "{{.key}} updated to {{.value}} template"
		src.TemplateNoBackupsFound = "{{.file}} no backups template"
		src.TemplateBackupCreated = "{{.path}} backup created template"

		// Keep some as default to test selective merging
		// src.TemplateVerificationFailed remains default

		mergeExtendedTemplates(dst, src)

		// Verify custom values were merged
		if dst.TemplateNoArchivesFound != "{{.custom}} template for no archives" {
			t.Errorf("Expected custom TemplateNoArchivesFound, got %q", dst.TemplateNoArchivesFound)
		}

		if dst.TemplateVerificationSuccess != "{{.archive}} verification success template" {
			t.Errorf("Expected custom TemplateVerificationSuccess, got %q", dst.TemplateVerificationSuccess)
		}

		if dst.TemplateConfigurationUpdated != "{{.key}} updated to {{.value}} template" {
			t.Errorf("Expected custom TemplateConfigurationUpdated, got %q", dst.TemplateConfigurationUpdated)
		}

		if dst.TemplateNoBackupsFound != "{{.file}} no backups template" {
			t.Errorf("Expected custom TemplateNoBackupsFound, got %q", dst.TemplateNoBackupsFound)
		}

		if dst.TemplateBackupCreated != "{{.path}} backup created template" {
			t.Errorf("Expected custom TemplateBackupCreated, got %q", dst.TemplateBackupCreated)
		}

		// Verify default values were not overridden
		defaultCfg := DefaultConfig()
		if dst.TemplateVerificationFailed != defaultCfg.TemplateVerificationFailed {
			t.Error("Default TemplateVerificationFailed should not be overridden")
		}
	})

	t.Run("all extended template strings covered", func(t *testing.T) {
		dst := DefaultConfig()
		src := DefaultConfig()

		// Set all template strings to custom values
		src.TemplateNoArchivesFound = "Template1"
		src.TemplateVerificationFailed = "Template2"
		src.TemplateVerificationSuccess = "Template3"
		src.TemplateVerificationWarning = "Template4"
		src.TemplateConfigurationUpdated = "Template5"
		src.TemplateConfigFilePath = "Template6"
		src.TemplateDryRunFilesHeader = "Template7"
		src.TemplateDryRunFileEntry = "Template8"
		src.TemplateNoFilesModified = "Template9"
		src.TemplateIncrementalCreated = "Template10"
		src.TemplateNoBackupsFound = "Template11"
		src.TemplateBackupWouldCreate = "Template12"
		src.TemplateBackupIdentical = "Template13"
		src.TemplateBackupCreated = "Template14"

		mergeExtendedTemplates(dst, src)

		// Verify all custom values were merged
		if dst.TemplateNoArchivesFound != "Template1" {
			t.Errorf("Expected 'Template1', got %q", dst.TemplateNoArchivesFound)
		}
		if dst.TemplateVerificationFailed != "Template2" {
			t.Errorf("Expected 'Template2', got %q", dst.TemplateVerificationFailed)
		}
		if dst.TemplateVerificationSuccess != "Template3" {
			t.Errorf("Expected 'Template3', got %q", dst.TemplateVerificationSuccess)
		}
		if dst.TemplateVerificationWarning != "Template4" {
			t.Errorf("Expected 'Template4', got %q", dst.TemplateVerificationWarning)
		}
		if dst.TemplateConfigurationUpdated != "Template5" {
			t.Errorf("Expected 'Template5', got %q", dst.TemplateConfigurationUpdated)
		}
		if dst.TemplateConfigFilePath != "Template6" {
			t.Errorf("Expected 'Template6', got %q", dst.TemplateConfigFilePath)
		}
		if dst.TemplateDryRunFilesHeader != "Template7" {
			t.Errorf("Expected 'Template7', got %q", dst.TemplateDryRunFilesHeader)
		}
		if dst.TemplateDryRunFileEntry != "Template8" {
			t.Errorf("Expected 'Template8', got %q", dst.TemplateDryRunFileEntry)
		}
		if dst.TemplateNoFilesModified != "Template9" {
			t.Errorf("Expected 'Template9', got %q", dst.TemplateNoFilesModified)
		}
		if dst.TemplateIncrementalCreated != "Template10" {
			t.Errorf("Expected 'Template10', got %q", dst.TemplateIncrementalCreated)
		}
		if dst.TemplateNoBackupsFound != "Template11" {
			t.Errorf("Expected 'Template11', got %q", dst.TemplateNoBackupsFound)
		}
		if dst.TemplateBackupWouldCreate != "Template12" {
			t.Errorf("Expected 'Template12', got %q", dst.TemplateBackupWouldCreate)
		}
		if dst.TemplateBackupIdentical != "Template13" {
			t.Errorf("Expected 'Template13', got %q", dst.TemplateBackupIdentical)
		}
		if dst.TemplateBackupCreated != "Template14" {
			t.Errorf("Expected 'Template14', got %q", dst.TemplateBackupCreated)
		}
	})
}

// 🔺 TEST-CONFIG-001: Test placeholder functions with basic implementation - 🔧
func TestPlaceholderConfigFunctions(t *testing.T) {
	t.Run("LoadConfigValues placeholder", func(t *testing.T) {
		dir := t.TempDir()

		// Test the placeholder implementation
		values, err := LoadConfigValues(dir)

		// Currently returns nil, nil - placeholder implementation
		if values != nil {
			t.Error("LoadConfigValues placeholder should return nil values")
		}
		if err != nil {
			t.Error("LoadConfigValues placeholder should return nil error")
		}
	})

	t.Run("mergeConfigValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)

		// Test that function doesn't panic
		mergeConfigValues(dst, src)
		// Function is placeholder, so no assertions on behavior
	})

	t.Run("mergeBasicSettingValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)
		cfg := DefaultConfig()

		// Test that function doesn't panic
		mergeBasicSettingValues(dst, src, cfg)
		// Function is placeholder, so no assertions on behavior
	})

	t.Run("mergeFileBackupSettingValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)
		cfg := DefaultConfig()

		// Test that function doesn't panic
		mergeFileBackupSettingValues(dst, src, cfg)
		// Function is placeholder, so no assertions on behavior
	})

	t.Run("mergeStatusCodeValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)
		cfg := DefaultConfig()

		// Test that function doesn't panic
		mergeStatusCodeValues(dst, src, cfg)
		// Function is placeholder, so no assertions on behavior
	})

	t.Run("mergeDirectoryStatusCodeValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)
		cfg := DefaultConfig()

		// Test that function doesn't panic
		mergeDirectoryStatusCodeValues(dst, src, cfg)
		// Function is placeholder, so no assertions on behavior
	})

	t.Run("mergeFileStatusCodeValues placeholder", func(t *testing.T) {
		dst := make(map[string]ConfigValue)
		src := make(map[string]ConfigValue)
		cfg := DefaultConfig()

		// Test that function doesn't panic
		mergeFileStatusCodeValues(dst, src, cfg)
		// Function is placeholder, so no assertions on behavior
	})
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

// ⭐ CFG-005: Configuration inheritance testing - 🧪 Comprehensive inheritance test coverage

// TestConfigInheritance tests basic configuration inheritance functionality
func TestConfigInheritance(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "config_inheritance_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create base configuration file
	baseConfig := `
archive_dir_path: "/base/archives"
exclude_patterns:
  - "*.log"
  - "*.tmp"
include_git_info: true
skip_broken_symlinks: true
`
	baseConfigPath := filepath.Join(tempDir, "base.yml")
	err = os.WriteFile(baseConfigPath, []byte(baseConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write base config: %v", err)
	}

	// Create child configuration file that inherits from base
	childConfig := `
inherit:
  - "base.yml"
archive_dir_path: "/child/archives"
+exclude_patterns:
  - "*.cache"
include_git_info: false
`
	childConfigPath := filepath.Join(tempDir, "child.yml")
	err = os.WriteFile(childConfigPath, []byte(childConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write child config: %v", err)
	}

	// Test inheritance chain building
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	chain, err := chainBuilder.buildChain(childConfigPath, pathResolver)
	if err != nil {
		t.Fatalf("Failed to build inheritance chain: %v", err)
	}

	// Verify chain contains both files in correct order
	if len(chain.files) != 2 {
		t.Errorf("Expected 2 files in chain, got %d", len(chain.files))
	}

	// Base file should come first (parent before child)
	if !strings.HasSuffix(chain.files[0], "base.yml") {
		t.Errorf("Expected base.yml to be first in chain, got %s", chain.files[0])
	}
	if !strings.HasSuffix(chain.files[1], "child.yml") {
		t.Errorf("Expected child.yml to be second in chain, got %s", chain.files[1])
	}
}

// TestMergeStrategies tests different merge strategies for configuration inheritance
func TestMergeStrategies(t *testing.T) {
	// Test merge strategy extraction
	processor := newMergeStrategyProcessor()

	tests := []struct {
		key      string
		expected string
		cleanKey string
	}{
		{"normal_key", "override", "normal_key"},
		{"+merge_key", "merge", "merge_key"},
		{"^prepend_key", "prepend", "prepend_key"},
		{"!replace_key", "replace", "replace_key"},
		{"=default_key", "default", "default_key"},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			config := map[string]interface{}{
				test.key: "test_value",
			}

			processed, err := processor.processKeys(config)
			if err != nil {
				t.Fatalf("Failed to process keys: %v", err)
			}

			operation, exists := processed.operations[test.cleanKey]
			if !exists {
				t.Errorf("Expected operation for key %s", test.cleanKey)
				return
			}

			if operation.strategy != test.expected {
				t.Errorf("Expected strategy %s, got %s", test.expected, operation.strategy)
			}
		})
	}
}

// TestArrayMergeStrategies tests array-specific merge strategies
func TestArrayMergeStrategies(t *testing.T) {
	// Create test configuration with array fields
	result := DefaultConfig()
	result.ExcludePatterns = []string{"base1", "base2"}

	tests := []struct {
		name     string
		strategy string
		srcValue []string
		expected []string
	}{
		{
			name:     "merge strategy appends",
			strategy: "merge",
			srcValue: []string{"child1", "child2"},
			expected: []string{"base1", "base2", "child1", "child2"},
		},
		{
			name:     "prepend strategy prepends",
			strategy: "prepend",
			srcValue: []string{"child1", "child2"},
			expected: []string{"child1", "child2", "base1", "base2"},
		},
		{
			name:     "replace strategy replaces",
			strategy: "replace",
			srcValue: []string{"child1", "child2"},
			expected: []string{"child1", "child2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset result for each test
			testResult := DefaultConfig()
			testResult.ExcludePatterns = []string{"base1", "base2"}

			operation := &mergeOperation{
				strategy: test.strategy,
				value:    test.srcValue,
				key:      "exclude_patterns",
			}

			dstValue := testResult.ExcludePatterns
			err := applyMergeOperation(testResult, "exclude_patterns", operation, dstValue)
			if err != nil {
				t.Fatalf("Failed to apply merge operation: %v", err)
			}

			if !equalStringSlices(testResult.ExcludePatterns, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, testResult.ExcludePatterns)
			}
		})
	}
}

// TestCircularDependencyDetection tests detection of circular inheritance dependencies
func TestCircularDependencyDetection(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "circular_dependency_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file A that inherits from B
	configA := `
inherit:
  - "config_b.yml"
archive_dir_path: "/a/archives"
`
	configAPath := filepath.Join(tempDir, "config_a.yml")
	err = os.WriteFile(configAPath, []byte(configA), 0644)
	if err != nil {
		t.Fatalf("Failed to write config A: %v", err)
	}

	// Create file B that inherits from A (circular dependency)
	configB := `
inherit:
  - "config_a.yml"
archive_dir_path: "/b/archives"
`
	configBPath := filepath.Join(tempDir, "config_b.yml")
	err = os.WriteFile(configBPath, []byte(configB), 0644)
	if err != nil {
		t.Fatalf("Failed to write config B: %v", err)
	}

	// Test that circular dependency is detected
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	_, err = chainBuilder.buildChain(configAPath, pathResolver)
	if err == nil {
		t.Error("Expected circular dependency error, but got nil")
	}

	if !strings.Contains(err.Error(), "circular dependency detected") {
		t.Errorf("Expected circular dependency error, got: %v", err)
	}
}

// TestDefaultValueStrategy tests the default value merge strategy
func TestDefaultValueStrategy(t *testing.T) {
	result := DefaultConfig()
	result.ArchiveDirPath = "" // Zero value

	// Test that default strategy only applies when destination is zero value
	operation := &mergeOperation{
		strategy: "default",
		value:    "/new/path",
		key:      "archive_dir_path",
	}

	err := applyMergeOperation(result, "archive_dir_path", operation, "")
	if err != nil {
		t.Fatalf("Failed to apply default operation: %v", err)
	}

	if result.ArchiveDirPath != "/new/path" {
		t.Errorf("Expected /new/path, got %s", result.ArchiveDirPath)
	}

	// Test that default strategy doesn't override non-zero values
	result.ArchiveDirPath = "/existing/path"
	operation.value = "/should/not/apply"

	err = applyMergeOperation(result, "archive_dir_path", operation, "/existing/path")
	if err != nil {
		t.Fatalf("Failed to apply default operation: %v", err)
	}

	if result.ArchiveDirPath != "/existing/path" {
		t.Errorf("Expected /existing/path to remain, got %s", result.ArchiveDirPath)
	}
}

// TestComplexInheritanceChain tests a complex multi-level inheritance chain
func TestComplexInheritanceChain(t *testing.T) {
	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "complex_inheritance_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create grandparent configuration
	grandparentConfig := `
archive_dir_path: "/grandparent/archives"
exclude_patterns:
  - "*.log"
include_git_info: true
`
	grandparentPath := filepath.Join(tempDir, "grandparent.yml")
	err = os.WriteFile(grandparentPath, []byte(grandparentConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write grandparent config: %v", err)
	}

	// Create parent configuration that inherits from grandparent
	parentConfig := `
inherit:
  - "grandparent.yml"
archive_dir_path: "/parent/archives"
+exclude_patterns:
  - "*.tmp"
skip_broken_symlinks: true
`
	parentPath := filepath.Join(tempDir, "parent.yml")
	err = os.WriteFile(parentPath, []byte(parentConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write parent config: %v", err)
	}

	// Create child configuration that inherits from parent
	childConfig := `
inherit:
  - "parent.yml"
+exclude_patterns:
  - "*.cache"
include_git_info: false
=show_git_dirty_status: true
`
	childPath := filepath.Join(tempDir, "child.yml")
	err = os.WriteFile(childPath, []byte(childConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write child config: %v", err)
	}

	// Test that complex inheritance chain builds correctly
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	chain, err := chainBuilder.buildChain(childPath, pathResolver)
	if err != nil {
		t.Fatalf("Failed to build complex inheritance chain: %v", err)
	}

	// Verify chain contains all files in correct order
	if len(chain.files) != 3 {
		t.Errorf("Expected 3 files in chain, got %d", len(chain.files))
	}

	// Files should be in order: grandparent, parent, child
	expectedOrder := []string{"grandparent.yml", "parent.yml", "child.yml"}
	for i, expected := range expectedOrder {
		if !strings.HasSuffix(chain.files[i], expected) {
			t.Errorf("Expected %s at position %d, got %s", expected, i, chain.files[i])
		}
	}
}
