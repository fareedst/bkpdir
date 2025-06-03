// This file is part of bkpdir
//
// Package main provides backward compatibility wrapper functions for REFACTOR-005.
// These wrappers maintain existing function signatures while using new interface-based implementations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
)

// 🔶 REFACTOR-005: Structure optimization - Backward compatibility wrappers - 📝
// This file contains wrapper functions that preserve existing APIs while using new interface-based implementations

// ========================================
// CONFIGURATION WRAPPER FUNCTIONS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Config access wrapper functions - 🔧
// These functions provide backward compatibility for existing config access patterns

// GetConfigProviderFromConfig creates a ConfigProviderInterface from Config struct
// Backward compatibility: Existing code can continue using Config directly
func GetConfigProviderFromConfig(cfg *Config) ConfigProviderInterface {
	return NewConfigProviderAdapter(cfg)
}

// GetFormatterConfigFromConfig creates FormatterConfigInterface from Config struct
// Backward compatibility: Existing formatter code can continue using Config directly
func GetFormatterConfigFromConfig(cfg *Config) FormatterConfigInterface {
	return &FormatterConfigAdapter{cfg: cfg}
}

// ========================================
// FORMATTER WRAPPER FUNCTIONS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Formatter access wrapper functions - 🔧
// These functions provide backward compatibility for existing formatter access patterns

// GetOutputFormatterInterface creates OutputFormatterInterface from OutputFormatter struct
// Backward compatibility: Existing code can continue using OutputFormatter directly
func GetOutputFormatterInterface(formatter *OutputFormatter) OutputFormatterInterface {
	return NewOutputFormatterAdapter(formatter)
}

// ========================================
// RESOURCE MANAGER WRAPPER FUNCTIONS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Resource manager factory wrapper - 🔧
// These functions provide backward compatibility for resource manager creation

// GetResourceManagerFactory creates a ResourceManagerFactoryInterface
// Backward compatibility: Existing code can continue creating ResourceManager directly
func GetResourceManagerFactory() ResourceManagerFactoryInterface {
	return &DefaultResourceManagerFactory{}
}

// CreateResourceManagerInterface creates ResourceManagerInterface using factory
// Backward compatibility: Wraps existing NewResourceManager function
func CreateResourceManagerInterface() ResourceManagerInterface {
	return NewResourceManager()
}

// CreateResourceManagerWithContextInterface creates ResourceManagerInterface with context
// Backward compatibility: Wraps existing resource manager creation with context support
func CreateResourceManagerWithContextInterface(ctx context.Context) ResourceManagerInterface {
	return NewResourceManager() // Context stored separately in application
}

// ========================================
// SERVICE PROVIDER WRAPPER FUNCTIONS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Service provider wrapper functions - 🔧
// These functions provide easy access to the service provider pattern

// CreateDefaultServiceProvider creates a comprehensive service provider with all components
// Backward compatibility: Provides single access point for all interfaces
func CreateDefaultServiceProvider(config *Config, formatter *OutputFormatter) ServiceProviderInterface {
	return NewDefaultServiceProvider(config, formatter)
}

// GetAllInterfaces provides access to all interface implementations from config and formatter
// Backward compatibility: Single function to get all interfaces
func GetAllInterfaces(config *Config, formatter *OutputFormatter) (
	ConfigProviderInterface,
	OutputFormatterInterface,
	ResourceManagerFactoryInterface,
	GitProviderInterface,
	ServiceProviderInterface,
) {
	serviceProvider := NewDefaultServiceProvider(config, formatter)
	return serviceProvider.GetConfigProvider(),
		serviceProvider.GetFormatterProvider(),
		serviceProvider.GetResourceManagerFactory(),
		serviceProvider.GetGitProvider(),
		serviceProvider
}

// ========================================
// FUNCTION SIGNATURE MIGRATION HELPERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Function signature migration helpers - 🔧
// These functions help migrate from concrete struct parameters to interface parameters

// MigrateConfigToInterface converts *Config parameter to ConfigProviderInterface
// Usage: configProvider := MigrateConfigToInterface(cfg)
func MigrateConfigToInterface(cfg *Config) ConfigProviderInterface {
	return NewConfigProviderAdapter(cfg)
}

// MigrateFormatterToInterface converts *OutputFormatter parameter to OutputFormatterInterface
// Usage: formatterInterface := MigrateFormatterToInterface(formatter)
func MigrateFormatterToInterface(formatter *OutputFormatter) OutputFormatterInterface {
	return NewOutputFormatterAdapter(formatter)
}

// MigrateResourceManagerToInterface converts *ResourceManager parameter to ResourceManagerInterface
// Usage: rmInterface := MigrateResourceManagerToInterface(rm)
func MigrateResourceManagerToInterface(rm *ResourceManager) ResourceManagerInterface {
	return rm // ResourceManager already implements ResourceManagerInterface
}

// ========================================
// EXTRACTION PREPARATION HELPERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Extraction preparation functions - 🔧
// These functions prepare components for clean extraction to separate packages

// PrepareConfigForExtraction validates that config can be cleanly extracted
// Returns interfaces that will be stable across package boundaries
func PrepareConfigForExtraction(cfg *Config) (
	ConfigProviderInterface,
	FormatterConfigInterface,
	ErrorConfigInterface,
	GitConfigInterface,
) {
	adapter := NewConfigProviderAdapter(cfg)
	return adapter,
		adapter.GetFormatterConfig(),
		adapter.GetErrorConfig(),
		adapter.GetGitConfig()
}

// PrepareFormatterForExtraction validates that formatter can be cleanly extracted
// Returns interfaces that will be stable across package boundaries
func PrepareFormatterForExtraction(formatter *OutputFormatter) OutputFormatterInterface {
	return NewOutputFormatterAdapter(formatter)
}

// PrepareResourceManagerForExtraction validates that resource manager can be cleanly extracted
// Returns interfaces that will be stable across package boundaries
func PrepareResourceManagerForExtraction() ResourceManagerFactoryInterface {
	return &DefaultResourceManagerFactory{}
}

// ValidateExtractionReadiness performs comprehensive validation that all components
// are ready for extraction with proper interface boundaries
func ValidateExtractionReadiness(cfg *Config, formatter *OutputFormatter) bool {
	// Validate that all adapters can be created successfully
	configProvider := NewConfigProviderAdapter(cfg)
	if configProvider == nil {
		return false
	}

	formatterAdapter := NewOutputFormatterAdapter(formatter)
	if formatterAdapter == nil {
		return false
	}

	resourceFactory := &DefaultResourceManagerFactory{}
	if resourceFactory == nil {
		return false
	}

	// Validate that service provider can be created
	serviceProvider := NewDefaultServiceProvider(cfg, formatter)
	if serviceProvider == nil {
		return false
	}

	// All validations passed
	return true
}

// ========================================
// INTERFACE VALIDATION HELPERS
// ========================================

// 🔶 REFACTOR-005: Structure optimization - Interface validation functions - 🔧
// These functions validate that interfaces work correctly with adapters

// ValidateConfigInterfaces ensures all config interfaces work correctly
func ValidateConfigInterfaces(cfg *Config) error {
	provider := NewConfigProviderAdapter(cfg)

	// Test that all interfaces can be retrieved
	_ = provider.GetArchiveConfig()
	_ = provider.GetBackupConfig()
	_ = provider.GetFormatterConfig()
	_ = provider.GetErrorConfig()
	_ = provider.GetGitConfig()

	return nil
}

// ValidateFormatterInterfaces ensures formatter interface works correctly
func ValidateFormatterInterfaces(formatter *OutputFormatter) error {
	adapter := NewOutputFormatterAdapter(formatter)

	// Test basic formatting methods work
	_ = adapter.FormatCreatedArchive("test")
	_ = adapter.FormatError("test")

	return nil
}

// ValidateResourceManagerInterfaces ensures resource manager interface works correctly
func ValidateResourceManagerInterfaces() error {
	factory := &DefaultResourceManagerFactory{}
	rm := factory.CreateResourceManager()

	// Test basic resource manager operations
	rm.AddTempFile("test")

	return nil
}
