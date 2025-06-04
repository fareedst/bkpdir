// Package config provides schema-agnostic configuration management for CLI applications.
//
// This test file verifies the extracted configuration system works independently
// from the original backup application.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestConfig represents a simple test configuration structure
type TestConfig struct {
	Name        string   `yaml:"name"`
	Port        int      `yaml:"port"`
	Enabled     bool     `yaml:"enabled"`
	Features    []string `yaml:"features"`
	Environment string   `yaml:"environment"`
}

// ⭐ EXTRACT-001: Package validation - Generic configuration loader test - 🧪
func TestGenericConfigLoader(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test configuration file
	configPath := filepath.Join(tempDir, "test.yml")
	configData := `
name: test-app
port: 8080
enabled: true
features:
  - logging
  - metrics
environment: testing
`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create path discovery with custom config
	discovery := NewPathDiscovery(DiscoveryConfig{
		EnvVarName:         "TEST_CONFIG",
		DefaultSearchPaths: []string{configPath},
		ConfigFileName:     "test.yml",
		HomeConfigPath:     "~/test.yml",
	})

	// Create a loader
	loader := NewGenericConfigLoader(
		discovery,
		NewDefaultEnvironmentProvider(),
		NewDefaultFileOperations(),
		NewGenericValidator(),
	)

	// Default configuration
	defaultConfig := &TestConfig{
		Name:        "default-app",
		Port:        3000,
		Enabled:     false,
		Features:    []string{"basic"},
		Environment: "development",
	}

	// Load configuration
	result, err := loader.LoadConfig(tempDir, defaultConfig)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	config, ok := result.(*TestConfig)
	if !ok {
		t.Fatalf("Expected *TestConfig, got %T", result)
	}

	// Verify configuration values
	if config.Name != "test-app" {
		t.Errorf("Expected name 'test-app', got '%s'", config.Name)
	}
	if config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Port)
	}
	if !config.Enabled {
		t.Errorf("Expected enabled true, got %v", config.Enabled)
	}
	if len(config.Features) != 2 {
		t.Errorf("Expected 2 features, got %d", len(config.Features))
	}
	if config.Environment != "testing" {
		t.Errorf("Expected environment 'testing', got '%s'", config.Environment)
	}
}

// ⭐ EXTRACT-001: Package validation - Configuration value extraction test - 🧪
func TestConfigValueExtraction(t *testing.T) {
	config := &TestConfig{
		Name:        "example-app",
		Port:        9000,
		Enabled:     true,
		Features:    []string{"auth", "cache"},
		Environment: "production",
	}

	loader := NewDefaultConfigLoader()
	values := loader.GetConfigValues(config)

	if len(values) != 5 {
		t.Errorf("Expected 5 config values, got %d", len(values))
	}

	// Check specific values
	for _, value := range values {
		switch value.Name {
		case "name":
			if value.Value != "example-app" {
				t.Errorf("Expected name 'example-app', got '%v'", value.Value)
			}
		case "port":
			if value.Value != 9000 {
				t.Errorf("Expected port 9000, got %v", value.Value)
			}
		case "enabled":
			if value.Value != true {
				t.Errorf("Expected enabled true, got %v", value.Value)
			}
		}
	}
}

// ⭐ EXTRACT-001: Package validation - Configuration merging test - 🧪
func TestConfigMerging(t *testing.T) {
	merger := NewGenericConfigMerger()

	base := &TestConfig{
		Name:        "base-app",
		Port:        3000,
		Enabled:     false,
		Features:    []string{"basic"},
		Environment: "development",
	}

	override := &TestConfig{
		Name: "override-app",
		Port: 8080,
		// Enabled not set (should keep base value)
		Features:    []string{"advanced", "premium"},
		Environment: "production",
	}

	// Merge configurations
	err := merger.MergeConfigs(base, override)
	if err != nil {
		t.Fatalf("Failed to merge configurations: %v", err)
	}

	// Verify merged values
	if base.Name != "override-app" {
		t.Errorf("Expected merged name 'override-app', got '%s'", base.Name)
	}
	if base.Port != 8080 {
		t.Errorf("Expected merged port 8080, got %d", base.Port)
	}
	if base.Enabled != false {
		t.Errorf("Expected enabled to remain false, got %v", base.Enabled)
	}
	if len(base.Features) != 2 {
		t.Errorf("Expected 2 features after merge, got %d", len(base.Features))
	}
	if base.Environment != "production" {
		t.Errorf("Expected merged environment 'production', got '%s'", base.Environment)
	}
}

// ⭐ EXTRACT-001: Package validation - Path discovery test - 🧪
func TestPathDiscovery(t *testing.T) {
	discovery := NewGenericPathDiscovery("myapp", ".myapp.yml")

	// Test environment variable name
	if discovery.GetEnvVarName() != "MYAPP_CONFIG" {
		t.Errorf("Expected env var 'MYAPP_CONFIG', got '%s'", discovery.GetEnvVarName())
	}

	// Test config file name
	if discovery.GetConfigFileName() != ".myapp.yml" {
		t.Errorf("Expected config file '.myapp.yml', got '%s'", discovery.GetConfigFileName())
	}

	// Test default paths
	defaultPaths := discovery.GetDefaultPaths()
	expectedPaths := []string{"./.myapp.yml", "~/.myapp.yml"}
	if len(defaultPaths) != len(expectedPaths) {
		t.Errorf("Expected %d default paths, got %d", len(expectedPaths), len(defaultPaths))
	}
}

// ⭐ EXTRACT-001: Package validation - Environment provider test - 🧪
func TestEnvironmentProvider(t *testing.T) {
	provider := NewDefaultEnvironmentProvider()

	// Set up field mapping
	mapping := map[string]string{
		"name":        "APP_NAME",
		"port":        "APP_PORT",
		"environment": "APP_ENV",
	}
	provider.SetEnvMapping(mapping)

	// Test mapping retrieval
	retrievedMapping := provider.GetEnvMapping()
	if len(retrievedMapping) != 3 {
		t.Errorf("Expected 3 mappings, got %d", len(retrievedMapping))
	}

	if retrievedMapping["name"] != "APP_NAME" {
		t.Errorf("Expected name mapping 'APP_NAME', got '%s'", retrievedMapping["name"])
	}

	// Test field-based environment access
	_ = provider.SetEnv("APP_NAME", "test-from-env")
	defer func() { _ = provider.SetEnv("APP_NAME", "") }()

	envValue := provider.GetEnvForField("name")
	if envValue != "test-from-env" {
		t.Errorf("Expected 'test-from-env', got '%s'", envValue)
	}
}

// ⭐ EXTRACT-001: Package validation - File operations test - 🧪
func TestFileOperations(t *testing.T) {
	fileOps := NewDefaultFileOperations()

	// Create a temporary file
	tempDir, err := os.MkdirTemp("", "fileops_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	testData := []byte("test file content")

	// Test write and read
	if err := fileOps.WriteFile(testFile, testData, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	if !fileOps.FileExists(testFile) {
		t.Error("File should exist after writing")
	}

	readData, err := fileOps.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(readData) != string(testData) {
		t.Errorf("Expected '%s', got '%s'", string(testData), string(readData))
	}

	// Test file info
	info, err := fileOps.GetFileInfo(testFile)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Size() != int64(len(testData)) {
		t.Errorf("Expected file size %d, got %d", len(testData), info.Size())
	}
}

// ⭐ EXTRACT-001: Package validation - Value extractor test - 🧪
func TestValueExtractor(t *testing.T) {
	extractor := NewGenericValueExtractor()

	config := &TestConfig{
		Name:        "test-app",
		Port:        8080,
		Enabled:     true,
		Features:    []string{"feature1"},
		Environment: "test",
	}

	defaultConfig := &TestConfig{
		Name:        "default-app",
		Port:        3000,
		Enabled:     false,
		Features:    []string{},
		Environment: "development",
	}

	sourceFunc := func(current, defaultValue interface{}) string {
		// Use reflect.DeepEqual for proper comparison
		if reflect.DeepEqual(current, defaultValue) {
			return "default"
		}
		return "custom"
	}

	values := extractor.ExtractValues(config, defaultConfig, sourceFunc)

	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	// All values should be marked as "custom" since they differ from defaults
	for _, value := range values {
		if value.Source != "custom" {
			t.Errorf("Expected source 'custom' for %s, got '%s'", value.Name, value.Source)
		}
	}
}

// BenchmarkConfigLoading benchmarks the configuration loading performance.
// ⭐ EXTRACT-001: Package validation - Performance benchmark - 🧪
func BenchmarkConfigLoading(b *testing.B) {
	// Create test setup
	tempDir, err := os.MkdirTemp("", "benchmark_config")
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "benchmark.yml")
	configData := `
name: benchmark-app
port: 9999
enabled: true
features: [feature1, feature2, feature3]
environment: benchmark
`
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		b.Fatalf("Failed to write config file: %v", err)
	}

	discovery := NewPathDiscovery(DiscoveryConfig{
		DefaultSearchPaths: []string{configPath},
		ConfigFileName:     "benchmark.yml",
	})

	loader := NewGenericConfigLoader(
		discovery,
		NewDefaultEnvironmentProvider(),
		NewDefaultFileOperations(),
		NewGenericValidator(),
	)

	defaultConfig := &TestConfig{
		Name:        "default",
		Port:        3000,
		Enabled:     false,
		Features:    []string{},
		Environment: "development",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.LoadConfig(tempDir, defaultConfig)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
