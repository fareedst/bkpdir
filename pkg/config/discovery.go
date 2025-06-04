// Package config provides configuration discovery and path management.
//
// This file contains the generalized configuration discovery logic extracted
// from the original backup application, made schema-agnostic and configurable.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"os"
	"path/filepath"
	"strings"
)

// 🔺 EXTRACT-001: Discovery generalization - Configurable search path provider - 🔍
// DiscoveryConfig holds configuration for path discovery behavior.
// This allows applications to customize discovery without hardcoded paths.
type DiscoveryConfig struct {
	// EnvVarName is the environment variable name for custom config paths
	EnvVarName string

	// DefaultSearchPaths are the fallback paths when env var is not set
	DefaultSearchPaths []string

	// ConfigFileName is the name of the configuration file to search for
	ConfigFileName string

	// HomeConfigPath is the home directory config path template (uses ~ expansion)
	HomeConfigPath string
}

// 🔺 EXTRACT-001: Discovery generalization - Generic path discovery implementation - 🔍
// PathDiscovery provides generalized configuration file discovery.
// This replaces the hardcoded backup-specific search logic.
type PathDiscovery struct {
	config DiscoveryConfig
}

// NewPathDiscovery creates a new PathDiscovery with the specified configuration.
// 🔺 EXTRACT-001: Discovery generalization - Configurable discovery constructor - 🔧
func NewPathDiscovery(config DiscoveryConfig) *PathDiscovery {
	return &PathDiscovery{
		config: config,
	}
}

// NewDefaultPathDiscovery creates a PathDiscovery with backup application defaults.
// This maintains backward compatibility while enabling generalization.
// 🔺 EXTRACT-001: Discovery generalization - Backward compatible defaults - 🔧
func NewDefaultPathDiscovery() *PathDiscovery {
	return NewPathDiscovery(DiscoveryConfig{
		EnvVarName:         "BKPDIR_CONFIG",
		DefaultSearchPaths: []string{"./.bkpdir.yml", "~/.bkpdir.yml"},
		ConfigFileName:     ".bkpdir.yml",
		HomeConfigPath:     "~/.bkpdir.yml",
	})
}

// NewGenericPathDiscovery creates a configurable PathDiscovery for any application.
// 🔺 EXTRACT-001: Discovery generalization - Generic application support - 🔧
func NewGenericPathDiscovery(appName, configFile string) *PathDiscovery {
	envVar := strings.ToUpper(appName) + "_CONFIG"
	localPath := "./" + configFile
	homePath := "~/" + configFile

	return NewPathDiscovery(DiscoveryConfig{
		EnvVarName:         envVar,
		DefaultSearchPaths: []string{localPath, homePath},
		ConfigFileName:     configFile,
		HomeConfigPath:     homePath,
	})
}

// GetConfigSearchPaths returns the search paths for configuration files.
// This generalizes the original getConfigSearchPaths() function.
// 🔺 EXTRACT-001: Discovery generalization - Generic search path resolution - 🔍
func (p *PathDiscovery) GetConfigSearchPaths() []string {
	if envPath := os.Getenv(p.config.EnvVarName); envPath != "" {
		// Split on colon for multiple paths (Unix convention)
		paths := strings.Split(envPath, ":")
		var result []string
		for _, path := range paths {
			result = append(result, strings.TrimSpace(path))
		}
		return result
	}

	return p.config.DefaultSearchPaths
}

// ExpandPath expands path variables and returns absolute path.
// This generalizes the original expandPath() function.
// 🔺 EXTRACT-001: Discovery generalization - Generic path expansion - 🔍
func (p *PathDiscovery) ExpandPath(path string) string {
	// Handle home directory expansion
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}

	// Handle current directory relative paths
	if strings.HasPrefix(path, "./") {
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Join(cwd, path[2:])
		}
	}

	// Return absolute path or clean relative path
	if filepath.IsAbs(path) {
		return path
	}

	// For relative paths, resolve against current working directory
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}

	// Fallback to original path
	return path
}

// GetConfigFileName returns the configuration file name.
// 🔶 EXTRACT-001: Value tracking abstraction - Configuration metadata access - 📝
func (p *PathDiscovery) GetConfigFileName() string {
	return p.config.ConfigFileName
}

// GetEnvVarName returns the environment variable name used for configuration.
// 🔶 EXTRACT-001: Value tracking abstraction - Environment variable name access - 📝
func (p *PathDiscovery) GetEnvVarName() string {
	return p.config.EnvVarName
}

// GetDefaultPaths returns the default search paths.
// 🔶 EXTRACT-001: Value tracking abstraction - Default path access - 📝
func (p *PathDiscovery) GetDefaultPaths() []string {
	return p.config.DefaultSearchPaths
}

// 🔺 EXTRACT-001: Environment variable abstraction - Generic environment provider - 🔧
// DefaultEnvironmentProvider provides the default implementation of EnvironmentProvider.
// This implementation uses the standard os package for environment variable access.
type DefaultEnvironmentProvider struct {
	mapping map[string]string // Field name to environment variable mapping
}

// NewDefaultEnvironmentProvider creates a new DefaultEnvironmentProvider.
// 🔺 EXTRACT-001: Environment variable abstraction - Provider constructor - 🔧
func NewDefaultEnvironmentProvider() *DefaultEnvironmentProvider {
	return &DefaultEnvironmentProvider{
		mapping: make(map[string]string),
	}
}

// NewBackupEnvironmentProvider creates an environment provider with backup app mappings.
// This maintains backward compatibility for the backup application.
// 🔺 EXTRACT-001: Environment variable abstraction - Backward compatible provider - 🔧
func NewBackupEnvironmentProvider() *DefaultEnvironmentProvider {
	provider := NewDefaultEnvironmentProvider()
	provider.SetEnvMapping(map[string]string{
		"archive_dir_path":               "BKPDIR_ARCHIVE_DIR",
		"backup_dir_path":                "BKPDIR_BACKUP_DIR",
		"include_git_info":               "BKPDIR_INCLUDE_GIT",
		"show_git_dirty_status":          "BKPDIR_SHOW_GIT_DIRTY",
		"use_current_dir_name":           "BKPDIR_USE_CURRENT_DIR_NAME",
		"use_current_dir_name_for_files": "BKPDIR_USE_CURRENT_DIR_NAME_FOR_FILES",
	})
	return provider
}

// GetEnv retrieves an environment variable value.
// 🔺 EXTRACT-001: Environment variable abstraction - Generic environment access - 🔧
func (d *DefaultEnvironmentProvider) GetEnv(key string) string {
	return os.Getenv(key)
}

// SetEnv sets an environment variable (primarily for testing).
// 🔺 EXTRACT-001: Environment variable abstraction - Environment modification - 🔧
func (d *DefaultEnvironmentProvider) SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// GetEnvMapping returns the mapping of config fields to environment variables.
// 🔶 EXTRACT-001: Value tracking abstraction - Mapping access - 📝
func (d *DefaultEnvironmentProvider) GetEnvMapping() map[string]string {
	result := make(map[string]string)
	for k, v := range d.mapping {
		result[k] = v
	}
	return result
}

// SetEnvMapping configures the environment variable mapping.
// 🔺 EXTRACT-001: Environment variable abstraction - Mapping configuration - 🔧
func (d *DefaultEnvironmentProvider) SetEnvMapping(mapping map[string]string) {
	d.mapping = make(map[string]string)
	for k, v := range mapping {
		d.mapping[k] = v
	}
}

// GetEnvForField returns the environment variable value for a specific config field.
// 🔺 EXTRACT-001: Environment variable abstraction - Field-based environment access - 🔧
func (d *DefaultEnvironmentProvider) GetEnvForField(fieldName string) string {
	if envVar, exists := d.mapping[fieldName]; exists {
		return d.GetEnv(envVar)
	}
	return ""
}
