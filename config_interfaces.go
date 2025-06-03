// This file is part of bkpdir
//
// Package main provides configuration interface abstractions for BkpDir.
// These interfaces enable schema-agnostic configuration management for extraction.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"os"
)

// 🔻 REFACTOR-003: Config abstraction - Schema-agnostic configuration loading interface - 🔧
// ConfigLoader provides schema-agnostic configuration management operations.
// This interface abstracts configuration loading from specific application schemas.
type ConfigLoader interface {
	// LoadConfig loads configuration from the specified root directory
	LoadConfig(root string) (*Config, error)

	// LoadConfigValues loads configuration values with source tracking
	LoadConfigValues(root string) (map[string]ConfigValue, error)

	// GetConfigValues extracts configuration values from a Config struct
	GetConfigValues(cfg *Config) []ConfigValue

	// GetConfigValuesWithSources extracts configuration values with source information
	GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue

	// ValidateConfig validates a configuration structure
	ValidateConfig(cfg *Config) error
}

// 🔻 REFACTOR-003: Config abstraction - Configuration merging and composition interface - 🔧
// ConfigMerger provides schema-agnostic configuration merging and composition operations.
// This interface enables reusable configuration merging logic across different schemas.
type ConfigMerger interface {
	// MergeConfigs merges source configuration into destination configuration
	MergeConfigs(dst, src *Config)

	// MergeConfigValues merges configuration value maps
	MergeConfigValues(dst, src map[string]ConfigValue)

	// GetConfigSearchPaths returns the search paths for configuration files
	GetConfigSearchPaths() []string

	// ExpandPath expands path variables and returns absolute path
	ExpandPath(path string) string
}

// 🔻 REFACTOR-003: Config abstraction - Configuration source abstraction interface - 🔧
// ConfigSource abstracts different configuration sources (file, environment, defaults).
// This interface enables pluggable configuration sources for different environments.
type ConfigSource interface {
	// LoadFromFile loads configuration from a file source
	LoadFromFile(path string) (*Config, error)

	// LoadFromEnvironment loads configuration from environment variables
	LoadFromEnvironment() (*Config, error)

	// LoadDefaults returns the default configuration
	LoadDefaults() *Config

	// GetSourceName returns the name of this configuration source
	GetSourceName() string

	// IsAvailable checks if this configuration source is available
	IsAvailable() bool
}

// 🔻 REFACTOR-003: Config abstraction - Pluggable configuration validation interface - 🔧
// ConfigValidator enables different applications to define their own configuration schemas.
// This interface allows for schema-specific validation while maintaining common validation logic.
type ConfigValidator interface {
	// ValidateSchema validates the configuration against the expected schema
	ValidateSchema(cfg *Config) error

	// ValidateValues validates individual configuration values
	ValidateValues(values map[string]ConfigValue) error

	// GetRequiredFields returns the list of required configuration fields
	GetRequiredFields() []string

	// GetValidationRules returns validation rules for configuration fields
	GetValidationRules() map[string]ValidationRule
}

// 🔻 REFACTOR-003: Schema separation - Application-specific configuration interface - 🔧
// ApplicationConfig provides access to application-specific configuration settings.
// This interface abstracts the backup-specific schema from generic configuration operations.
type ApplicationConfig interface {
	// GetArchiveSettings returns archive-related configuration
	GetArchiveSettings() ArchiveSettings

	// GetBackupSettings returns backup-related configuration
	GetBackupSettings() BackupSettings

	// GetStatusCodes returns application status codes
	GetStatusCodes() map[string]int

	// GetFormatSettings returns output formatting configuration
	GetFormatSettings() FormatSettings
}

// 🔻 REFACTOR-003: Schema separation - Backup-specific archive settings structure - 📝
// ArchiveSettings contains archive-specific configuration settings.
// This structure separates archive concerns from generic configuration.
type ArchiveSettings struct {
	DirectoryPath      string
	UseCurrentDirName  bool
	ExcludePatterns    []string
	IncludeGitInfo     bool
	ShowGitDirtyStatus bool
	Verification       *VerificationConfig
}

// 🔻 REFACTOR-003: Schema separation - Backup-specific file backup settings structure - 📝
// BackupSettings contains file backup-specific configuration settings.
// This structure separates backup concerns from generic configuration.
type BackupSettings struct {
	DirectoryPath             string
	UseCurrentDirNameForFiles bool
}

// 🔻 REFACTOR-003: Schema separation - Application-specific format settings structure - 📝
// FormatSettings contains output formatting configuration settings.
// This structure separates formatting concerns from generic configuration.
type FormatSettings struct {
	FormatStrings      map[string]string
	TemplateStrings    map[string]string
	PatternStrings     map[string]string
	ErrorFormatStrings map[string]string
}

// 🔻 REFACTOR-003: Config abstraction - Configuration validation rule structure - 📝
// ValidationRule defines validation criteria for configuration fields.
// This structure enables flexible validation rules for different application schemas.
type ValidationRule struct {
	Required     bool
	Type         string
	MinLength    int
	MaxLength    int
	Pattern      string
	ValidValues  []string
	Dependencies []string
}

// 🔻 REFACTOR-003: Config abstraction - Configuration source determination interface - 🔧
// SourceDeterminer provides methods to determine configuration value sources.
// This interface enables source tracking across different configuration providers.
type SourceDeterminer interface {
	// DetermineSource determines the source of a configuration value
	DetermineSource(current, defaultValue interface{}) string

	// GetConfigSource returns the primary configuration source
	GetConfigSource() string

	// GetSourcePriority returns the priority order of configuration sources
	GetSourcePriority() []string
}

// 🔻 REFACTOR-003: Config abstraction - Generic configuration value extractor interface - 🔧
// ValueExtractor provides methods to extract configuration values from different structures.
// This interface enables schema-agnostic value extraction for different application types.
type ValueExtractor interface {
	// ExtractBasicValues extracts basic configuration values
	ExtractBasicValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue

	// ExtractStatusCodeValues extracts status code configuration values
	ExtractStatusCodeValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue

	// ExtractVerificationValues extracts verification configuration values
	ExtractVerificationValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue

	// ExtractFormatValues extracts format string configuration values
	ExtractFormatValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue
}

// 🔻 REFACTOR-003: Config abstraction - File operation interface for configuration - 🔧
// ConfigFileOperations provides file system operations for configuration management.
// This interface abstracts file operations to enable testing and different storage backends.
type ConfigFileOperations interface {
	// FileExists checks if a configuration file exists
	FileExists(path string) bool

	// ReadFile reads configuration file contents
	ReadFile(path string) ([]byte, error)

	// WriteFile writes configuration file contents
	WriteFile(path string, data []byte, perm os.FileMode) error

	// GetFileInfo returns file information for configuration files
	GetFileInfo(path string) (os.FileInfo, error)
}
