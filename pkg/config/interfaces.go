// Package config provides schema-agnostic configuration management for CLI applications.
//
// This package offers a complete configuration management system that supports:
// - Multiple configuration sources (files, environment variables, defaults)
// - Schema-agnostic configuration loading and merging
// - Source tracking for configuration values
// - Pluggable validation and customization
//
// The design uses interface-based abstractions to support different application
// schemas while preserving robust discovery, merging, and validation logic.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"os"
)

// ⭐ EXTRACT-001: Core interface extraction - Schema-agnostic configuration loading interface - 📝
// ConfigLoader provides schema-agnostic configuration management operations.
// This interface abstracts configuration loading from specific application schemas.
type ConfigLoader interface {
	// LoadConfig loads configuration from the specified root directory with default config
	LoadConfig(root string, defaultConfig interface{}) (interface{}, error)

	// LoadConfigValues loads configuration values with source tracking
	LoadConfigValues(root string, defaultConfig interface{}) (map[string]ConfigValue, error)

	// GetConfigValues extracts configuration values from a configuration struct
	GetConfigValues(cfg interface{}) []ConfigValue

	// GetConfigValuesWithSources extracts configuration values with source information
	GetConfigValuesWithSources(cfg interface{}, root string) []ConfigValue

	// ValidateConfig validates a configuration structure
	ValidateConfig(cfg interface{}) error
}

// ⭐ EXTRACT-001: Core interface extraction - Configuration merging and composition interface - 📝
// ConfigMerger provides schema-agnostic configuration merging and composition operations.
// This interface enables reusable configuration merging logic across different schemas.
type ConfigMerger interface {
	// MergeConfigs merges source configuration into destination configuration
	MergeConfigs(dst, src interface{}) error

	// MergeConfigValues merges configuration value maps
	MergeConfigValues(dst, src map[string]ConfigValue)

	// GetConfigSearchPaths returns the search paths for configuration files
	GetConfigSearchPaths() []string

	// ExpandPath expands path variables and returns absolute path
	ExpandPath(path string) string
}

// ⭐ EXTRACT-001: Core interface extraction - Configuration source abstraction interface - 📝
// ConfigSource abstracts different configuration sources (file, environment, defaults).
// This interface enables pluggable configuration sources for different environments.
type ConfigSource interface {
	// LoadFromFile loads configuration from a file source
	LoadFromFile(path string) (interface{}, error)

	// LoadFromEnvironment loads configuration from environment variables
	LoadFromEnvironment() (interface{}, error)

	// LoadDefaults returns the default configuration
	LoadDefaults() interface{}

	// GetSourceName returns the name of this configuration source
	GetSourceName() string

	// IsAvailable checks if this configuration source is available
	IsAvailable() bool
}

// ⭐ EXTRACT-001: Core interface extraction - Pluggable configuration validation interface - 📝
// ConfigValidator enables different applications to define their own configuration schemas.
// This interface allows for schema-specific validation while maintaining common validation logic.
type ConfigValidator interface {
	// ValidateSchema validates the configuration against the expected schema
	ValidateSchema(cfg interface{}) error

	// ValidateValues validates individual configuration values
	ValidateValues(values map[string]ConfigValue) error

	// GetRequiredFields returns the list of required configuration fields
	GetRequiredFields() []string

	// GetValidationRules returns validation rules for configuration fields
	GetValidationRules() map[string]ValidationRule
}

// ⭐ EXTRACT-001: Core interface extraction - Application-specific configuration interface - 📝
// ApplicationConfig provides access to application-specific configuration settings.
// This interface abstracts application-specific schema from generic configuration operations.
type ApplicationConfig interface {
	// GetSetting returns a configuration setting by name
	GetSetting(name string) interface{}

	// GetAllSettings returns all configuration settings
	GetAllSettings() map[string]interface{}

	// SetSetting sets a configuration setting
	SetSetting(name string, value interface{}) error

	// GetConfigSchema returns the configuration schema definition
	GetConfigSchema() interface{}
}

// 🔺 EXTRACT-001: Configuration structure extraction - Configuration validation rule structure - 📝
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

// 🔺 EXTRACT-001: Configuration structure extraction - Configuration source determination interface - 📝
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

// 🔺 EXTRACT-001: Configuration structure extraction - Generic configuration value extractor interface - 📝
// ValueExtractor provides methods to extract configuration values from different structures.
// This interface enables schema-agnostic value extraction for different application types.
type ValueExtractor interface {
	// ExtractValues extracts configuration values from any configuration structure
	ExtractValues(cfg, defaultCfg interface{}, getSource func(interface{}, interface{}) string) []ConfigValue

	// ExtractValuesByCategory extracts values from a specific category
	ExtractValuesByCategory(cfg, defaultCfg interface{}, category string, getSource func(interface{}, interface{}) string) []ConfigValue

	// GetSupportedCategories returns the list of supported value categories
	GetSupportedCategories() []string
}

// 🔺 EXTRACT-001: Configuration structure extraction - File operation interface for configuration - 📝
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

// 🔺 EXTRACT-001: Configuration structure extraction - Environment variable provider interface - 📝
// EnvironmentProvider abstracts environment variable operations for configuration.
// This interface enables configurable environment variable mapping and testing.
type EnvironmentProvider interface {
	// GetEnv retrieves an environment variable value
	GetEnv(key string) string

	// SetEnv sets an environment variable (primarily for testing)
	SetEnv(key, value string) error

	// GetEnvMapping returns the mapping of config fields to environment variables
	GetEnvMapping() map[string]string

	// SetEnvMapping configures the environment variable mapping
	SetEnvMapping(mapping map[string]string)
}

// 🔶 EXTRACT-001: Value tracking abstraction - Configuration value representation - 📝
// ConfigValue represents a configuration value with its source information.
// This structure enables source tracking for any configuration type.
type ConfigValue struct {
	Name   string      // Configuration field name
	Value  interface{} // Configuration value (made generic)
	Source string      // Source of the configuration value
	Type   string      // Type of the configuration value
}

// ⭐ CFG-005: Configuration inheritance interfaces - 🔧 Core inheritance functionality

// InheritanceLoader extends ConfigLoader with inheritance support.
// This interface enables configuration files to inherit from other configuration files.
type InheritanceLoader interface {
	ConfigLoader

	// LoadConfigWithInheritance loads configuration with inheritance chain processing
	LoadConfigWithInheritance(root string, defaultConfig interface{}) (interface{}, error)

	// GetInheritanceChain returns the inheritance chain for a configuration file
	GetInheritanceChain(configPath string) (*InheritanceChain, error)

	// ValidateInheritanceChain validates the inheritance chain for circular dependencies
	ValidateInheritanceChain(chain *InheritanceChain) error
}

// MergeStrategy defines how configuration values should be merged during inheritance.
// This interface enables flexible merge strategies for different data types.
type MergeStrategy interface {
	// Merge applies the merge strategy to destination and source values
	Merge(dest, src interface{}) (interface{}, error)

	// GetPrefix returns the prefix that triggers this merge strategy
	GetPrefix() string

	// GetDescription returns a human-readable description of the merge strategy
	GetDescription() string

	// SupportsType checks if the strategy supports the given data type
	SupportsType(valueType string) bool
}

// InheritanceChainBuilder builds and validates inheritance dependency chains.
// This interface abstracts the process of building inheritance hierarchies.
type InheritanceChainBuilder interface {
	// BuildChain builds the inheritance chain starting from the given configuration file
	BuildChain(configPath string, pathResolver PathResolver) (*InheritanceChain, error)

	// ValidateChain validates the inheritance chain for circular dependencies
	ValidateChain(chain *InheritanceChain) error

	// GetChainMetadata returns metadata about the inheritance chain
	GetChainMetadata(chain *InheritanceChain) *InheritanceChainMetadata
}

// PathResolver resolves and expands configuration file paths.
// This interface handles path resolution for inheritance hierarchies.
type PathResolver interface {
	// ResolvePath resolves a configuration file path (absolute, relative, ~)
	ResolvePath(path string, basePath string) (string, error)

	// ExpandPath expands variables in configuration file paths
	ExpandPath(path string) (string, error)

	// ValidatePath validates that a configuration file path is safe and accessible
	ValidatePath(path string) error
}

// CircularDependencyDetector detects circular dependencies in inheritance chains.
// This interface provides cycle detection for configuration inheritance.
type CircularDependencyDetector interface {
	// DetectCycle performs cycle detection on an inheritance chain
	DetectCycle(startFile string, resolver PathResolver) error

	// GetCyclePath returns the path of detected circular dependency
	GetCyclePath() []string

	// Reset resets the detector state for new cycle detection
	Reset()
}

// MergeStrategyProcessor processes configuration keys with merge strategy prefixes.
// This interface handles prefix-based merge strategy application.
type MergeStrategyProcessor interface {
	// ProcessKeys processes configuration keys and applies appropriate merge strategies
	ProcessKeys(config map[string]interface{}) (*ProcessedConfig, error)

	// RegisterStrategy registers a new merge strategy
	RegisterStrategy(strategy MergeStrategy) error

	// GetAvailableStrategies returns all available merge strategies
	GetAvailableStrategies() map[string]MergeStrategy
}

// InheritanceSourceTracker tracks the source of configuration values through inheritance.
// This interface provides enhanced source tracking for inherited configurations.
type InheritanceSourceTracker interface {
	SourceDeterminer

	// TrackInheritanceSource tracks the source of a value through inheritance chain
	TrackInheritanceSource(key string, value interface{}, sourceFile string, strategy string) *ValueOrigin

	// GetValueOrigin returns detailed origin information for a configuration value
	GetValueOrigin(key string) *ValueOrigin

	// GetChainSources returns all source files in the inheritance chain
	GetChainSources() []string
}

// ⭐ CFG-005: Configuration inheritance data structures - 📝 Core data models

// InheritanceChain represents a configuration inheritance dependency chain.
type InheritanceChain struct {
	Files        []string            // Configuration files in dependency order
	Visited      map[string]bool     // Circular dependency prevention
	Sources      map[string]string   // Source tracking for debugging
	Dependencies map[string][]string // File dependencies mapping
}

// InheritanceChainMetadata provides metadata about an inheritance chain.
type InheritanceChainMetadata struct {
	ChainLength  int                 // Total files in chain
	MaxDepth     int                 // Maximum inheritance depth
	SourceFiles  []string            // All source files
	Dependencies map[string][]string // Dependency mapping
	BuildTime    int64               // Time to build chain (nanoseconds)
}

// ProcessedConfig represents configuration after merge strategy processing.
type ProcessedConfig struct {
	Config     map[string]interface{} // Processed configuration
	MergeOps   []MergeOperation       // Applied merge operations
	SourceMap  map[string]string      // Value source mapping
	Strategies map[string]string      // Applied strategies per key
}

// MergeOperation represents a single merge operation applied during inheritance.
type MergeOperation struct {
	Key        string        // Configuration key
	Value      interface{}   // Configuration value
	Strategy   MergeStrategy // Applied merge strategy
	SourceFile string        // Source file for this operation
	TargetType string        // Target data type
}

// ValueOrigin provides detailed origin information for configuration values.
type ValueOrigin struct {
	SourceFile    string // Configuration file containing the value
	SourceType    string // "inherit", "override", "merge", "prepend", "replace", "default"
	MergeStrategy string // Merge strategy applied
	LineNumber    int    // Line number in source file (if available)
	ChainDepth    int    // Depth in inheritance chain
	Resolution    string // How the value was resolved
}

// ConfigInheritance represents inheritance metadata in configuration files.
type ConfigInheritance struct {
	Inherit []string `yaml:"inherit"` // List of parent configuration files
}
