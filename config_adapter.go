// This file is part of bkpdir
//
// Package main provides backward compatibility adapter for the extracted configuration system.
// It bridges the original Config struct with the extracted pkg/config package.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration adapter for extracted package - 🔧

import (
	"fmt"
	"strings"

	"bkpdir/pkg/config"
)

// ⭐ EXTRACT-001: Backward compatibility layer - Application-specific configuration adapter - 📝
// ConfigAdapter provides backward compatibility between the original Config struct
// and the extracted pkg/config package. This allows the original application to
// continue working unchanged while using the extracted configuration system.
type ConfigAdapter struct {
	loader      config.ConfigLoader
	merger      config.ConfigMerger
	validator   config.ConfigValidator
	extractor   config.ValueExtractor
	fileOps     config.ConfigFileOperations
	envProvider config.EnvironmentProvider
}

// ⭐ EXTRACT-001: Backward compatibility layer - Adapter constructor - 📝
// NewConfigAdapter creates a new configuration adapter using the extracted package.
func NewConfigAdapter() *ConfigAdapter {
	pathDiscovery := config.NewDefaultPathDiscovery()
	envProvider := config.NewBackupEnvironmentProvider()
	fileOps := config.NewDefaultFileOperations()
	validator := config.NewGenericValidator()

	return &ConfigAdapter{
		loader:      config.NewGenericConfigLoader(pathDiscovery, envProvider, fileOps, validator),
		merger:      config.NewGenericConfigMerger(),
		validator:   validator,
		extractor:   config.NewGenericValueExtractor(),
		fileOps:     fileOps,
		envProvider: envProvider,
	}
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration loading adapter - 🔧
// LoadConfig loads configuration using the extracted package but returns the original Config struct.
// This maintains backward compatibility with existing code.
func (a *ConfigAdapter) LoadConfig(root string) (*Config, error) {
	// Use the extracted package to load configuration with default config
	defaultConfig := DefaultConfig()
	genericCfg, err := a.loader.LoadConfig(root, defaultConfig)
	if err != nil {
		return nil, err
	}

	// Convert the generic configuration to the original Config struct
	if cfg, ok := genericCfg.(*Config); ok {
		return cfg, nil
	}

	// If conversion fails, create a new Config with defaults
	return DefaultConfig(), nil
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration values adapter - 🔧
// LoadConfigValues loads configuration values using the extracted package.
func (a *ConfigAdapter) LoadConfigValues(root string) (map[string]ConfigValue, error) {
	// Use the extracted package to load configuration values with default config
	defaultConfig := DefaultConfig()
	genericValues, err := a.loader.LoadConfigValues(root, defaultConfig)
	if err != nil {
		return nil, err
	}

	// Convert generic ConfigValue to original ConfigValue
	result := make(map[string]ConfigValue)
	for name, genericValue := range genericValues {
		result[name] = ConfigValue{
			Name:   genericValue.Name,
			Value:  convertToString(genericValue.Value),
			Source: genericValue.Source,
		}
	}

	return result, nil
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration values with sources adapter - 🔧
// GetConfigValuesWithSources extracts configuration values with sources using the extracted package.
func (a *ConfigAdapter) GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue {
	// Use the extracted package to get configuration values with sources
	genericValues := a.loader.GetConfigValuesWithSources(cfg, root)

	// Convert generic ConfigValue to original ConfigValue
	result := make([]ConfigValue, len(genericValues))
	for i, genericValue := range genericValues {
		result[i] = ConfigValue{
			Name:   genericValue.Name,
			Value:  convertToString(genericValue.Value),
			Source: genericValue.Source,
		}
	}

	return result
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration values extraction adapter - 📝
// GetConfigValues extracts configuration values using the extracted package.
func (a *ConfigAdapter) GetConfigValues(cfg *Config) []ConfigValue {
	// Use the extracted package to get configuration values
	genericValues := a.loader.GetConfigValues(cfg)

	// Convert generic ConfigValue to original ConfigValue
	result := make([]ConfigValue, len(genericValues))
	for i, genericValue := range genericValues {
		result[i] = ConfigValue{
			Name:   genericValue.Name,
			Value:  convertToString(genericValue.Value),
			Source: genericValue.Source,
		}
	}

	return result
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration validation adapter - 📝
// ValidateConfig validates configuration using the extracted package.
func (a *ConfigAdapter) ValidateConfig(cfg *Config) error {
	return a.validator.ValidateSchema(cfg)
}

// ⭐ EXTRACT-001: Backward compatibility layer - Configuration search paths adapter - 📝
// GetConfigSearchPaths returns configuration search paths using the extracted package.
func (a *ConfigAdapter) GetConfigSearchPaths() []string {
	return a.merger.GetConfigSearchPaths()
}

// ⭐ EXTRACT-001: Backward compatibility layer - Path expansion adapter - 📝
// ExpandPath expands path variables using the extracted package.
func (a *ConfigAdapter) ExpandPath(path string) string {
	return a.merger.ExpandPath(path)
}

// 🔺 EXTRACT-001: Backward compatibility layer - Value conversion utility - 📝
// convertToString converts interface{} values to strings for backward compatibility.
func convertToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case []string:
		return strings.Join(v, ", ")
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ⭐ EXTRACT-001: Backward compatibility layer - Global adapter instance - 📝
// Global adapter instance for backward compatibility
var globalConfigAdapter = NewConfigAdapter()

// ⭐ EXTRACT-001: Backward compatibility layer - Backward compatible function wrappers - 🔧
// These functions maintain the original API while using the extracted package internally.

// LoadConfig loads configuration using the extracted package (backward compatible).
func LoadConfigWithAdapter(root string) (*Config, error) {
	return globalConfigAdapter.LoadConfig(root)
}

// LoadConfigValues loads configuration values using the extracted package (backward compatible).
func LoadConfigValuesWithAdapter(root string) (map[string]ConfigValue, error) {
	return globalConfigAdapter.LoadConfigValues(root)
}

// GetConfigValuesWithSources extracts configuration values with sources (backward compatible).
func GetConfigValuesWithSourcesAdapter(cfg *Config, root string) []ConfigValue {
	return globalConfigAdapter.GetConfigValuesWithSources(cfg, root)
}

// GetConfigValues extracts configuration values (backward compatible).
func GetConfigValuesAdapter(cfg *Config) []ConfigValue {
	return globalConfigAdapter.GetConfigValues(cfg)
}

// ValidateConfig validates configuration (backward compatible).
func ValidateConfigAdapter(cfg *Config) error {
	return globalConfigAdapter.ValidateConfig(cfg)
}

// getConfigSearchPaths returns configuration search paths (backward compatible).
func getConfigSearchPathsAdapter() []string {
	return globalConfigAdapter.GetConfigSearchPaths()
}

// expandPath expands path variables (backward compatible).
func expandPathAdapter(path string) string {
	return globalConfigAdapter.ExpandPath(path)
}
