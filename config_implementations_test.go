// This file is part of bkpdir
//
// Package main provides tests for the configuration abstraction implementations.
// These tests validate that the abstraction layer works correctly and provides
// backward compatibility while enabling schema-agnostic configuration management.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"os"
	"path/filepath"
	"testing"
)

// REFACTOR-003: Config abstraction - Test configuration abstraction implementations
func TestConfigAbstraction(t *testing.T) {
	t.Run("GenericConfigLoader", func(t *testing.T) {
		// Test generic configuration loading
		t.Run("LoadConfigWithDefaultSchema", func(t *testing.T) {
			// Create temporary directory
			tempDir := t.TempDir()

			// Create configuration loader
			loader := NewGenericConfigLoader()

			// Load configuration (should use defaults)
			config, err := loader.LoadConfig(tempDir)
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}

			// Verify config is loaded
			if config == nil {
				t.Fatal("Config should not be nil")
			}

			// Verify it's a backup config
			if cfg, ok := config.(*Config); ok {
				if cfg.ArchiveDirPath == "" {
					t.Error("ArchiveDirPath should have default value")
				}
			} else {
				t.Errorf("Expected *Config, got %T", config)
			}
		})

		t.Run("LoadConfigWithYAMLFile", func(t *testing.T) {
			// Create temporary directory with config file
			tempDir := t.TempDir()
			configFile := filepath.Join(tempDir, ".bkpdir.yml")

			// Write test config file
			configContent := `archive_dir_path: "/test/archives"
use_current_dir_name: false
exclude_patterns:
  - "*.tmp"
  - "build/"
`
			if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Create configuration loader
			loader := NewGenericConfigLoader()
			schema := NewBackupConfigSchema()

			// Load configuration
			config, err := loader.LoadConfigWithSchema(tempDir, schema)
			if err != nil {
				t.Fatalf("LoadConfigWithSchema failed: %v", err)
			}

			// Verify config values were loaded from file
			if cfg, ok := config.(*Config); ok {
				if cfg.ArchiveDirPath != "/test/archives" {
					t.Errorf("Expected ArchiveDirPath=/test/archives, got %s", cfg.ArchiveDirPath)
				}
				if cfg.UseCurrentDirName != false {
					t.Error("Expected UseCurrentDirName=false")
				}
				if len(cfg.ExcludePatterns) != 2 {
					t.Errorf("Expected 2 exclude patterns, got %d", len(cfg.ExcludePatterns))
				}
			} else {
				t.Errorf("Expected *Config, got %T", config)
			}
		})

		t.Run("GetConfigSearchPaths", func(t *testing.T) {
			loader := NewGenericConfigLoader()
			paths := loader.GetConfigSearchPaths()

			if len(paths) == 0 {
				t.Error("Config search paths should not be empty")
			}

			// Should include ./.bkpdir.yml (first in search order)
			found := false
			for _, path := range paths {
				if path == "./.bkpdir.yml" {
					found = true
					break
				}
			}
			if !found {
				t.Error("Config search paths should include ./.bkpdir.yml")
			}
		})
	})

	t.Run("ConfigSources", func(t *testing.T) {
		t.Run("YAMLConfigSource", func(t *testing.T) {
			source := NewYAMLConfigSource()

			// Test source properties
			if source.GetSourceType() != "file" {
				t.Error("YAML source should be 'file' type")
			}

			if source.GetPriority() <= 0 {
				t.Error("YAML source should have positive priority")
			}

			// Test with temporary file
			tempFile := filepath.Join(t.TempDir(), "test.yml")
			testData := map[string]interface{}{
				"test_key": "test_value",
				"test_int": 42,
			}

			// Test Save
			if err := source.Save(tempFile, testData); err != nil {
				t.Fatalf("Save failed: %v", err)
			}

			// Test Exists
			if !source.Exists(tempFile) {
				t.Error("Exists should return true for saved file")
			}

			// Test Load
			loadedData, err := source.Load(tempFile)
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if loadedData["test_key"] != "test_value" {
				t.Error("Loaded data doesn't match saved data")
			}
		})

		t.Run("EnvConfigSource", func(t *testing.T) {
			source := NewEnvConfigSource()

			// Test source properties
			if source.GetSourceType() != "env" {
				t.Error("Env source should be 'env' type")
			}

			if !source.Exists("") {
				t.Error("Env source should always exist")
			}

			// Set test environment variable
			oldValue := os.Getenv("BKPDIR_TEST_KEY")
			defer func() {
				if oldValue == "" {
					os.Unsetenv("BKPDIR_TEST_KEY")
				} else {
					os.Setenv("BKPDIR_TEST_KEY", oldValue)
				}
			}()

			os.Setenv("BKPDIR_TEST_KEY", "test_value")

			// Test Load
			data, err := source.Load("")
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if data["test_key"] != "test_value" {
				t.Errorf("Expected test_key=test_value, got %v", data["test_key"])
			}
		})

		t.Run("DefaultConfigSource", func(t *testing.T) {
			source := NewDefaultConfigSource()

			// Test source properties
			if source.GetSourceType() != "default" {
				t.Error("Default source should be 'default' type")
			}

			if !source.Exists("") {
				t.Error("Default source should always exist")
			}

			// Test Load (should return empty map)
			data, err := source.Load("")
			if err != nil {
				t.Fatalf("Load failed: %v", err)
			}

			if len(data) != 0 {
				t.Error("Default source should return empty data")
			}
		})
	})

	t.Run("GenericConfigMerger", func(t *testing.T) {
		t.Run("OverwriteStrategy", func(t *testing.T) {
			merger := NewGenericConfigMerger(OverwriteStrategy)

			if merger.GetMergeStrategy() != OverwriteStrategy {
				t.Error("Merge strategy should be OverwriteStrategy")
			}

			// Test strategy change
			merger.SetMergeStrategy(PreserveStrategy)
			if merger.GetMergeStrategy() != PreserveStrategy {
				t.Error("Merge strategy should be updated to PreserveStrategy")
			}
		})

		t.Run("MergeWithPriority", func(t *testing.T) {
			merger := NewGenericConfigMerger(OverwriteStrategy)

			// Create test configs
			cfg1 := &Config{ArchiveDirPath: "path1", UseCurrentDirName: true}
			cfg2 := &Config{ArchiveDirPath: "path2", UseCurrentDirName: false}
			cfg3 := &Config{ArchiveDirPath: "path3"}

			configs := []interface{}{cfg1, cfg2, cfg3}
			priorities := []int{10, 30, 20} // cfg2 should win

			result, err := merger.MergeWithPriority(configs, priorities)
			if err != nil {
				t.Fatalf("MergeWithPriority failed: %v", err)
			}

			if resultCfg, ok := result.(*Config); ok {
				if resultCfg.ArchiveDirPath != "path2" {
					t.Errorf("Expected highest priority path (path2), got %s", resultCfg.ArchiveDirPath)
				}
			} else {
				t.Errorf("Expected *Config, got %T", result)
			}
		})

		t.Run("ErrorCases", func(t *testing.T) {
			merger := NewGenericConfigMerger(OverwriteStrategy)

			// Test empty configs
			_, err := merger.MergeWithPriority([]interface{}{}, []int{})
			if err == nil {
				t.Error("Should error on empty configs")
			}

			// Test mismatched lengths
			configs := []interface{}{&Config{}}
			priorities := []int{1, 2}
			_, err = merger.MergeWithPriority(configs, priorities)
			if err == nil {
				t.Error("Should error on mismatched lengths")
			}
		})
	})

	t.Run("ConfigAdapter", func(t *testing.T) {
		t.Run("ToLegacyConfig", func(t *testing.T) {
			// Create test config
			originalCfg := &Config{
				ArchiveDirPath:    "/test/path",
				UseCurrentDirName: false,
				ExcludePatterns:   []string{"*.tmp"},
				IncludeGitInfo:    true,
				BackupDirPath:     "/test/backup",
			}

			// Create provider and adapter
			schema := NewBackupConfigSchema()
			provider := NewBackupConfigProvider(originalCfg, schema)
			adapter := NewConfigAdapter(provider, schema)

			// Convert to legacy config
			legacyCfg := adapter.ToLegacyConfig()

			// Verify values are preserved
			if legacyCfg.ArchiveDirPath != originalCfg.ArchiveDirPath {
				t.Errorf("ArchiveDirPath mismatch: expected %s, got %s", originalCfg.ArchiveDirPath, legacyCfg.ArchiveDirPath)
			}
			if legacyCfg.UseCurrentDirName != originalCfg.UseCurrentDirName {
				t.Errorf("UseCurrentDirName mismatch: expected %v, got %v", originalCfg.UseCurrentDirName, legacyCfg.UseCurrentDirName)
			}
			if len(legacyCfg.ExcludePatterns) != len(originalCfg.ExcludePatterns) {
				t.Errorf("ExcludePatterns length mismatch: expected %d, got %d", len(originalCfg.ExcludePatterns), len(legacyCfg.ExcludePatterns))
			}
		})

		t.Run("FromLegacyConfig", func(t *testing.T) {
			schema := NewBackupConfigSchema()
			provider := NewBackupConfigProvider(&Config{}, schema)
			adapter := NewConfigAdapter(provider, schema)

			// Create legacy config
			legacyCfg := &Config{
				ArchiveDirPath: "/legacy/path",
			}

			// Convert from legacy config
			newAdapter := adapter.FromLegacyConfig(legacyCfg)

			// Verify adapter is created
			if newAdapter == nil {
				t.Error("FromLegacyConfig should return adapter")
			}

			// Verify provider and schema are accessible
			if newAdapter.GetProvider() == nil {
				t.Error("Adapter should have provider")
			}
			if newAdapter.GetSchema() == nil {
				t.Error("Adapter should have schema")
			}
		})
	})

	t.Run("BackwardCompatibility", func(t *testing.T) {
		t.Run("LoadConfigLegacy", func(t *testing.T) {
			// Test backward compatible loading
			tempDir := t.TempDir()

			// Load using legacy function
			cfg, err := LoadConfigLegacy(tempDir)
			if err != nil {
				t.Fatalf("LoadConfigLegacy failed: %v", err)
			}

			// Verify it returns legacy Config struct
			if cfg == nil {
				t.Fatal("Config should not be nil")
			}

			if cfg.ArchiveDirPath == "" {
				t.Error("Config should have default archive path")
			}
		})

		t.Run("ExistingLoadConfigStillWorks", func(t *testing.T) {
			// Test that existing LoadConfig function still works
			tempDir := t.TempDir()

			cfg, err := LoadConfig(tempDir)
			if err != nil {
				t.Fatalf("LoadConfig failed: %v", err)
			}

			if cfg == nil {
				t.Fatal("Config should not be nil")
			}
		})
	})
}

// REFACTOR-003: Config abstraction - Test configuration provider implementation
func TestBackupConfigProvider(t *testing.T) {
	t.Run("GetMethods", func(t *testing.T) {
		cfg := &Config{
			ArchiveDirPath:    "/test/archive",
			UseCurrentDirName: true,
			ExcludePatterns:   []string{"*.tmp", "build/"},
			StatusConfigError: 42,
		}

		schema := NewBackupConfigSchema()
		provider := NewBackupConfigProvider(cfg, schema)

		// Test GetString
		if provider.GetString("archive_dir_path") != "/test/archive" {
			t.Error("GetString failed for archive_dir_path")
		}

		// Test GetBool
		if !provider.GetBool("use_current_dir_name") {
			t.Error("GetBool failed for use_current_dir_name")
		}

		// Test GetInt
		if provider.GetInt("status_config_error") != 42 {
			t.Error("GetInt failed for status_config_error")
		}

		// Test GetStringSlice
		patterns := provider.GetStringSlice("exclude_patterns")
		if len(patterns) != 2 {
			t.Errorf("GetStringSlice failed: expected 2 patterns, got %d", len(patterns))
		}

		// Test HasKey
		if !provider.HasKey("archive_dir_path") {
			t.Error("HasKey should return true for existing key")
		}

		if provider.HasKey("nonexistent_key") {
			t.Error("HasKey should return false for nonexistent key")
		}

		// Test GetAllKeys
		keys := provider.GetAllKeys()
		if len(keys) == 0 {
			t.Error("GetAllKeys should return keys")
		}

		// Test GetWithDefault
		defaultValue := provider.GetWithDefault("nonexistent_key", "default_value")
		if defaultValue != "default_value" {
			t.Error("GetWithDefault should return default for nonexistent key")
		}

		existingValue := provider.GetWithDefault("archive_dir_path", "default_value")
		if existingValue != "/test/archive" {
			t.Error("GetWithDefault should return existing value")
		}
	})
}

// REFACTOR-003: Schema separation - Test backup config schema implementation
func TestBackupConfigSchema(t *testing.T) {
	t.Run("SchemaProperties", func(t *testing.T) {
		schema := NewBackupConfigSchema()

		// Test schema name and version
		if schema.GetSchemaName() != "backup-application" {
			t.Error("Schema name should be 'backup-application'")
		}

		if schema.GetSchemaVersion() == "" {
			t.Error("Schema version should not be empty")
		}

		// Test default config
		defaultConfig := schema.GetDefaultConfig()
		if defaultConfig == nil {
			t.Error("Default config should not be nil")
		}

		if cfg, ok := defaultConfig.(*Config); !ok {
			t.Error("Default config should be *Config type")
		} else if cfg.ArchiveDirPath == "" {
			t.Error("Default config should have archive path")
		}

		// Test field definitions
		fieldDefs := schema.GetFieldDefinitions()
		if len(fieldDefs) == 0 {
			t.Error("Field definitions should not be empty")
		}

		// Test specific field definition
		if archivePathDef, exists := fieldDefs["archive_dir_path"]; exists {
			if archivePathDef.Required != true {
				t.Error("archive_dir_path should be required")
			}
			if archivePathDef.Type != "string" {
				t.Error("archive_dir_path should be string type")
			}
		} else {
			t.Error("archive_dir_path field definition should exist")
		}

		// Test required fields
		requiredFields := schema.GetRequiredFields()
		if len(requiredFields) == 0 {
			t.Error("Should have required fields")
		}
	})

	t.Run("ValidateSchema", func(t *testing.T) {
		schema := NewBackupConfigSchema()

		// Test valid configuration
		validCfg := &Config{
			ArchiveDirPath: "/valid/path",
			BackupDirPath:  "/valid/backup",
		}

		if err := schema.ValidateSchema(validCfg); err != nil {
			t.Errorf("Valid config should pass validation: %v", err)
		}

		// Test invalid configuration type
		if err := schema.ValidateSchema("invalid"); err == nil {
			t.Error("Invalid config type should fail validation")
		}

		// Test missing required field
		invalidCfg := &Config{
			ArchiveDirPath: "", // Required field missing
			BackupDirPath:  "/valid/backup",
		}

		if err := schema.ValidateSchema(invalidCfg); err == nil {
			t.Error("Config with missing required field should fail validation")
		}
	})

	t.Run("MigrateSchema", func(t *testing.T) {
		schema := NewBackupConfigSchema()

		// Test migration with same version (should pass through)
		cfg := &Config{ArchiveDirPath: "/test"}
		migrated, err := schema.MigrateSchema("1.0.0", "1.0.0", cfg)
		if err != nil {
			t.Errorf("Same version migration should succeed: %v", err)
		}
		if migrated != cfg {
			t.Error("Same version migration should return same config")
		}

		// Test migration with different version (should fail for now)
		_, err = schema.MigrateSchema("1.0.0", "2.0.0", cfg)
		if err == nil {
			t.Error("Different version migration should fail (not implemented)")
		}
	})
}
