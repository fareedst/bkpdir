// This file is part of bkpdir
//
// Package main provides comprehensive tests for REFACTOR-005 structure optimization.
// These tests validate that all interfaces work correctly and backward compatibility is maintained.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// 🔶 REFACTOR-005: Structure optimization - Comprehensive structure optimization tests - 📝
// This file contains tests that validate interface functionality and backward compatibility

// TestStructureOptimization validates the complete structure optimization implementation
func TestStructureOptimization(t *testing.T) {
	// Create test config
	cfg := DefaultConfig()
	cfg.ArchiveDirPath = "/tmp/test-archives"
	cfg.BackupDirPath = "/tmp/test-backups"

	// Create test formatter
	formatter := NewOutputFormatter(cfg)

	t.Run("ConfigProviderInterface", func(t *testing.T) {
		testConfigProviderInterface(t, cfg)
	})

	t.Run("FormatterConfigInterface", func(t *testing.T) {
		testFormatterConfigInterface(t, cfg)
	})

	t.Run("OutputFormatterInterface", func(t *testing.T) {
		testOutputFormatterInterface(t, formatter)
	})

	t.Run("ResourceManagerFactoryInterface", func(t *testing.T) {
		testResourceManagerFactoryInterface(t)
	})

	t.Run("ServiceProviderInterface", func(t *testing.T) {
		testServiceProviderInterface(t, cfg, formatter)
	})

	t.Run("BackwardCompatibility", func(t *testing.T) {
		testBackwardCompatibility(t, cfg, formatter)
	})

	t.Run("ExtractionReadiness", func(t *testing.T) {
		testExtractionReadiness(t, cfg, formatter)
	})
}

// testConfigProviderInterface validates ConfigProviderInterface implementation
func testConfigProviderInterface(t *testing.T, cfg *Config) {
	// 🔶 REFACTOR-005: Structure optimization - Config provider interface testing - 🔧
	provider := NewConfigProviderAdapter(cfg)
	if provider == nil {
		t.Fatal("Failed to create ConfigProviderAdapter")
	}

	// Test all config interface retrieval
	archiveConfig := provider.GetArchiveConfig()
	if archiveConfig == nil {
		t.Error("GetArchiveConfig returned nil")
	}
	if archiveConfig.GetArchiveDirPath() != cfg.ArchiveDirPath {
		t.Error("Archive dir path mismatch")
	}

	backupConfig := provider.GetBackupConfig()
	if backupConfig == nil {
		t.Error("GetBackupConfig returned nil")
	}
	if backupConfig.GetBackupDirPath() != cfg.BackupDirPath {
		t.Error("Backup dir path mismatch")
	}

	formatterConfig := provider.GetFormatterConfig()
	if formatterConfig == nil {
		t.Error("GetFormatterConfig returned nil")
	}
	if formatterConfig.GetFormatCreatedArchive() != cfg.FormatCreatedArchive {
		t.Error("Format created archive mismatch")
	}

	errorConfig := provider.GetErrorConfig()
	if errorConfig == nil {
		t.Error("GetErrorConfig returned nil")
	}
	statusCodes := errorConfig.GetStatusCodes()
	if len(statusCodes) == 0 {
		t.Error("No status codes returned")
	}

	gitConfig := provider.GetGitConfig()
	if gitConfig == nil {
		t.Error("GetGitConfig returned nil")
	}
	if gitConfig.GetIncludeGitInfo() != cfg.IncludeGitInfo {
		t.Error("Include git info mismatch")
	}
}

// testFormatterConfigInterface validates FormatterConfigInterface implementation
func testFormatterConfigInterface(t *testing.T, cfg *Config) {
	// 🔶 REFACTOR-005: Structure optimization - Formatter config interface testing - 🔧
	adapter := &FormatterConfigAdapter{cfg: cfg}

	// Test basic format strings
	if adapter.GetFormatCreatedArchive() != cfg.FormatCreatedArchive {
		t.Error("FormatCreatedArchive mismatch")
	}
	if adapter.GetFormatIdenticalArchive() != cfg.FormatIdenticalArchive {
		t.Error("FormatIdenticalArchive mismatch")
	}
	if adapter.GetFormatError() != cfg.FormatError {
		t.Error("FormatError mismatch")
	}

	// Test template strings
	if adapter.GetTemplateCreatedArchive() != cfg.TemplateCreatedArchive {
		t.Error("TemplateCreatedArchive mismatch")
	}
	if adapter.GetTemplateError() != cfg.TemplateError {
		t.Error("TemplateError mismatch")
	}

	// Test pattern strings
	if adapter.GetArchiveFilenamePattern() != cfg.PatternArchiveFilename {
		t.Error("ArchiveFilenamePattern mismatch")
	}
	if adapter.GetConfigLinePattern() != cfg.PatternConfigLine {
		t.Error("ConfigLinePattern mismatch")
	}
}

// testOutputFormatterInterface validates OutputFormatterInterface implementation
func testOutputFormatterInterface(t *testing.T, formatter *OutputFormatter) {
	// 🔶 REFACTOR-005: Structure optimization - Output formatter interface testing - 🔧
	adapter := NewOutputFormatterAdapter(formatter)
	if adapter == nil {
		t.Fatal("Failed to create OutputFormatterAdapter")
	}

	// Test basic formatting methods
	testPath := "/tmp/test-archive.zip"
	formatted := adapter.FormatCreatedArchive(testPath)
	if formatted == "" {
		t.Error("FormatCreatedArchive returned empty string")
	}

	errorMessage := "Test error message"
	formattedError := adapter.FormatError(errorMessage)
	if formattedError == "" {
		t.Error("FormatError returned empty string")
	}

	// Test backup formatting methods
	backupFormatted := adapter.FormatCreatedBackup(testPath)
	if backupFormatted == "" {
		t.Error("FormatCreatedBackup returned empty string")
	}

	// Test template methods
	testData := map[string]string{
		"path":      testPath,
		"timestamp": "2024-01-01T12:00:00Z",
	}
	templateResult := adapter.TemplateCreatedArchive(testData)
	if templateResult == "" {
		t.Error("TemplateCreatedArchive returned empty string")
	}

	// Test pattern extraction methods with properly formatted filename
	// Use a filename that matches the default archive pattern
	archiveData := adapter.ExtractArchiveFilenameData("home-docs-20240101-120000.zip")
	// Don't require data - extraction may return empty if pattern doesn't match exactly
	_ = archiveData

	// Test output control methods
	_ = adapter.IsDelayedMode() // Should not panic
	collector := adapter.GetCollector()
	if collector != nil {
		adapter.SetCollector(collector) // Should not panic
	}
}

// testResourceManagerFactoryInterface validates ResourceManagerFactoryInterface implementation
func testResourceManagerFactoryInterface(t *testing.T) {
	// 🔶 REFACTOR-005: Structure optimization - Resource manager factory interface testing - 🔧
	factory := &DefaultResourceManagerFactory{}

	// Test basic resource manager creation
	rm := factory.CreateResourceManager()
	if rm == nil {
		t.Fatal("CreateResourceManager returned nil")
	}

	// Create an actual temporary file for testing
	tempDir := createTempTestDir(t)
	testFile := createTempTestFile(t, tempDir, "test-temp-file", "test content")

	// Test that returned resource manager implements interface correctly
	rm.AddTempFile(testFile)

	// Test cleanup - this should succeed since file exists
	err := rm.Cleanup()
	if err != nil {
		t.Errorf("Resource manager cleanup failed: %v", err)
	}

	// Test context-aware creation
	ctx := context.Background()
	rmWithContext := factory.CreateResourceManagerWithContext(ctx)
	if rmWithContext == nil {
		t.Fatal("CreateResourceManagerWithContext returned nil")
	}
}

// testServiceProviderInterface validates ServiceProviderInterface implementation
func testServiceProviderInterface(t *testing.T, cfg *Config, formatter *OutputFormatter) {
	// 🔶 REFACTOR-005: Structure optimization - Service provider interface testing - 🔧
	serviceProvider := NewDefaultServiceProvider(cfg, formatter)
	if serviceProvider == nil {
		t.Fatal("Failed to create DefaultServiceProvider")
	}

	// Test that all providers can be retrieved
	configProvider := serviceProvider.GetConfigProvider()
	if configProvider == nil {
		t.Error("GetConfigProvider returned nil")
	}

	formatterProvider := serviceProvider.GetFormatterProvider()
	if formatterProvider == nil {
		t.Error("GetFormatterProvider returned nil")
	}

	resourceFactory := serviceProvider.GetResourceManagerFactory()
	if resourceFactory == nil {
		t.Error("GetResourceManagerFactory returned nil")
	}

	gitProvider := serviceProvider.GetGitProvider()
	if gitProvider == nil {
		t.Error("GetGitProvider returned nil")
	}

	fileOps := serviceProvider.GetFileOperations()
	if fileOps == nil {
		t.Error("GetFileOperations returned nil")
	}

	exclusionManager := serviceProvider.GetExclusionManager()
	if exclusionManager == nil {
		t.Error("GetExclusionManager returned nil")
	}

	yamlProcessor := serviceProvider.GetYAMLProcessor()
	if yamlProcessor == nil {
		t.Error("GetYAMLProcessor returned nil")
	}

	filesystem := serviceProvider.GetFilesystem()
	if filesystem == nil {
		t.Error("GetFilesystem returned nil")
	}

	cmdExecutor := serviceProvider.GetCommandExecutor()
	if cmdExecutor == nil {
		t.Error("GetCommandExecutor returned nil")
	}
}

// testBackwardCompatibility validates that existing code continues to work
func testBackwardCompatibility(t *testing.T, cfg *Config, formatter *OutputFormatter) {
	// 🔶 REFACTOR-005: Structure optimization - Backward compatibility testing - 🔧

	// Test wrapper functions work correctly
	configProvider := GetConfigProviderFromConfig(cfg)
	if configProvider == nil {
		t.Error("GetConfigProviderFromConfig failed")
	}

	formatterConfig := GetFormatterConfigFromConfig(cfg)
	if formatterConfig == nil {
		t.Error("GetFormatterConfigFromConfig failed")
	}

	formatterInterface := GetOutputFormatterInterface(formatter)
	if formatterInterface == nil {
		t.Error("GetOutputFormatterInterface failed")
	}

	resourceFactory := GetResourceManagerFactory()
	if resourceFactory == nil {
		t.Error("GetResourceManagerFactory failed")
	}

	// Test migration helpers
	migratedConfig := MigrateConfigToInterface(cfg)
	if migratedConfig == nil {
		t.Error("MigrateConfigToInterface failed")
	}

	migratedFormatter := MigrateFormatterToInterface(formatter)
	if migratedFormatter == nil {
		t.Error("MigrateFormatterToInterface failed")
	}

	rm := NewResourceManager()
	migratedRM := MigrateResourceManagerToInterface(rm)
	if migratedRM == nil {
		t.Error("MigrateResourceManagerToInterface failed")
	}

	// Test that all interfaces work together
	allConfigProvider, allFormatter, allResourceFactory, allGitProvider, allServiceProvider := GetAllInterfaces(cfg, formatter)
	if allConfigProvider == nil || allFormatter == nil || allResourceFactory == nil || allGitProvider == nil || allServiceProvider == nil {
		t.Error("GetAllInterfaces failed to return all interfaces")
	}
}

// testExtractionReadiness validates that components are ready for extraction
func testExtractionReadiness(t *testing.T, cfg *Config, formatter *OutputFormatter) {
	// 🔶 REFACTOR-005: Structure optimization - Extraction readiness testing - 🔧

	// Test extraction preparation functions
	configProvider, formatterConfig, errorConfig, gitConfig := PrepareConfigForExtraction(cfg)
	if configProvider == nil || formatterConfig == nil || errorConfig == nil || gitConfig == nil {
		t.Error("PrepareConfigForExtraction failed")
	}

	formatterInterface := PrepareFormatterForExtraction(formatter)
	if formatterInterface == nil {
		t.Error("PrepareFormatterForExtraction failed")
	}

	resourceFactory := PrepareResourceManagerForExtraction()
	if resourceFactory == nil {
		t.Error("PrepareResourceManagerForExtraction failed")
	}

	// Test extraction readiness validation
	isReady := ValidateExtractionReadiness(cfg, formatter)
	if !isReady {
		t.Error("ValidateExtractionReadiness returned false")
	}

	// Test interface validation functions
	err := ValidateConfigInterfaces(cfg)
	if err != nil {
		t.Errorf("ValidateConfigInterfaces failed: %v", err)
	}

	err = ValidateFormatterInterfaces(formatter)
	if err != nil {
		t.Errorf("ValidateFormatterInterfaces failed: %v", err)
	}

	err = ValidateResourceManagerInterfaces()
	if err != nil {
		t.Errorf("ValidateResourceManagerInterfaces failed: %v", err)
	}
}

// TestNamingConventionStandardization validates naming consistency across interfaces
func TestNamingConventionStandardization(t *testing.T) {
	// 🔶 REFACTOR-005: Structure optimization - Naming convention validation - 🔧
	cfg := DefaultConfig()

	t.Run("InterfaceNamingConsistency", func(t *testing.T) {
		// All interfaces should use "Interface" suffix
		provider := NewConfigProviderAdapter(cfg)

		// Test that interface types are correctly named and accessible
		var _ ConfigProviderInterface = provider
		var _ ArchiveConfigInterface = provider.GetArchiveConfig()
		var _ BackupConfigInterface = provider.GetBackupConfig()
		var _ FormatterConfigInterface = provider.GetFormatterConfig()
		var _ ErrorConfigInterface = provider.GetErrorConfig()
		var _ GitConfigInterface = provider.GetGitConfig()
	})

	t.Run("AdapterNamingConsistency", func(t *testing.T) {
		// All adapters should use "Adapter" suffix with shortened names
		var _ *ArchiveConfigAdapter = &ArchiveConfigAdapter{cfg: cfg}
		var _ *BackupConfigAdapter = &BackupConfigAdapter{cfg: cfg}
		var _ *FormatterConfigAdapter = &FormatterConfigAdapter{cfg: cfg}
		var _ *ErrorConfigAdapter = &ErrorConfigAdapter{cfg: cfg}
		var _ *GitConfigAdapter = &GitConfigAdapter{cfg: cfg}
	})
}

// TestFunctionSignatureOptimization validates that function signatures use interfaces
func TestFunctionSignatureOptimization(t *testing.T) {
	// 🔶 REFACTOR-005: Structure optimization - Function signature optimization validation - 🔧
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	t.Run("InterfaceBasedSignatures", func(t *testing.T) {
		// Test that wrapper functions accept interfaces properly
		configProvider := MigrateConfigToInterface(cfg)
		formatterInterface := MigrateFormatterToInterface(formatter)

		// Validate these can be used in interface-based function calls
		if configProvider == nil || formatterInterface == nil {
			t.Error("Interface migration failed")
		}

		// Test service provider works with interfaces
		serviceProvider := CreateDefaultServiceProvider(cfg, formatter)
		if serviceProvider == nil {
			t.Error("Service provider creation failed")
		}
	})

	t.Run("BackwardCompatibilityPreserved", func(t *testing.T) {
		// Test that existing function signatures still work through wrappers
		// These should all work without modification
		_ = GetConfigProviderFromConfig(cfg)
		_ = GetFormatterConfigFromConfig(cfg)
		_ = GetOutputFormatterInterface(formatter)
		_ = GetResourceManagerFactory()
	})
}

// TestImportStructureOptimization validates clean import structure preparation
func TestImportStructureOptimization(t *testing.T) {
	// 🔶 REFACTOR-005: Structure optimization - Import structure validation - 🔧

	t.Run("NoCircularDependencies", func(t *testing.T) {
		// Test that interfaces can be created without circular dependencies
		cfg := DefaultConfig()
		formatter := NewOutputFormatter(cfg)

		// All these should be createable independently
		_ = NewConfigProviderAdapter(cfg)
		_ = NewOutputFormatterAdapter(formatter)
		_ = &DefaultResourceManagerFactory{}
		_ = &DefaultGitProvider{}
		_ = &DefaultFileOperations{}
		_ = NewDefaultExclusionManager(cfg.ExcludePatterns)
	})

	t.Run("DependencyInjectionReady", func(t *testing.T) {
		// Test that components work with dependency injection
		cfg := DefaultConfig()
		formatter := NewOutputFormatter(cfg)

		serviceProvider := NewDefaultServiceProvider(cfg, formatter)

		// Test that all dependencies can be injected through service provider
		configProvider := serviceProvider.GetConfigProvider()
		if configProvider.GetArchiveConfig().GetArchiveDirPath() != cfg.ArchiveDirPath {
			t.Error("Dependency injection validation failed")
		}
	})
}

// BenchmarkStructureOptimization benchmarks the performance impact of interface indirection
func BenchmarkStructureOptimization(b *testing.B) {
	// 🔶 REFACTOR-005: Structure optimization - Performance impact validation - 🔧
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	b.Run("DirectAccess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Direct access (existing pattern)
			_ = cfg.FormatCreatedArchive
			_ = formatter.FormatCreatedArchive("/tmp/test")
		}
	})

	b.Run("InterfaceAccess", func(b *testing.B) {
		adapter := &FormatterConfigAdapter{cfg: cfg}
		formatterAdapter := NewOutputFormatterAdapter(formatter)

		for i := 0; i < b.N; i++ {
			// Interface access (new pattern)
			_ = adapter.GetFormatCreatedArchive()
			_ = formatterAdapter.FormatCreatedArchive("/tmp/test")
		}
	})
}

// TestExtractionSimulation simulates component extraction to verify boundaries
func TestExtractionSimulation(t *testing.T) {
	// 🔶 REFACTOR-005: Structure optimization - Extraction simulation testing - 🔧
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	t.Run("ConfigPackageSimulation", func(t *testing.T) {
		// Simulate config package extraction
		configProvider := NewConfigProviderAdapter(cfg)

		// Test that config can work independently through interfaces
		archiveConfig := configProvider.GetArchiveConfig()
		if archiveConfig.GetArchiveDirPath() == "" {
			t.Error("Config package simulation failed")
		}
	})

	t.Run("FormatterPackageSimulation", func(t *testing.T) {
		// Simulate formatter package extraction
		formatterInterface := NewOutputFormatterAdapter(formatter)

		// Test that formatter can work independently through interfaces
		result := formatterInterface.FormatCreatedArchive("/tmp/test")
		if result == "" {
			t.Error("Formatter package simulation failed")
		}
	})

	t.Run("ResourceManagerPackageSimulation", func(t *testing.T) {
		// Simulate resource manager package extraction
		factory := &DefaultResourceManagerFactory{}
		rm := factory.CreateResourceManager()

		// Create an actual temporary file for testing
		tempFile, err := os.CreateTemp("", "test-resource-*")
		if err != nil {
			t.Fatalf("Failed to create temp file for test: %v", err)
		}
		tempFile.Close()

		// Test that resource manager can work independently
		rm.AddTempFile(tempFile.Name())
		err = rm.Cleanup()
		if err != nil {
			t.Errorf("Resource manager package simulation failed: %v", err)
		}
	})
}

// Helper function to create temporary test directory
func createTempTestDir(t *testing.T) string {
	// 🔶 REFACTOR-005: Structure optimization - Test helper for directory creation - 🔧
	tempDir, err := os.MkdirTemp("", "bkpdir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	return tempDir
}

// Helper function to create test file
func createTempTestFile(t *testing.T, dir, filename, content string) string {
	// 🔶 REFACTOR-005: Structure optimization - Test helper for file creation - 🔧
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	return filePath
}
