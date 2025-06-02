// This file is part of bkpdir
//
// Package main provides configuration abstraction interfaces for schema-agnostic
// configuration management. These interfaces enable the configuration system
// to be reused across different applications with different schemas.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

// REFACTOR-003: Config abstraction - Core configuration loading interface
// ConfigLoader provides schema-agnostic configuration loading capabilities.
// It abstracts the loading process from specific configuration schemas.
type ConfigLoader interface {
	// LoadConfig loads configuration using default schema detection
	LoadConfig(root string) (interface{}, error)

	// LoadConfigWithSchema loads configuration with a specific schema
	LoadConfigWithSchema(root string, schema ConfigSchema) (interface{}, error)

	// GetConfigSearchPaths returns the search paths for configuration files
	GetConfigSearchPaths() []string

	// ExpandPath expands environment variables and tildes in paths
	ExpandPath(path string) string
}

// REFACTOR-003: Config abstraction - Pluggable validation interface
// ConfigValidator provides pluggable validation for different configuration schemas.
// It allows applications to define their own validation rules.
type ConfigValidator interface {
	// ValidateConfig validates an entire configuration object
	ValidateConfig(config interface{}) error

	// ValidateField validates a specific field value
	ValidateField(fieldName string, value interface{}) error

	// GetValidationRules returns all validation rules for the schema
	GetValidationRules() map[string]ValidationRule

	// AddValidationRule adds a new validation rule
	AddValidationRule(fieldName string, rule ValidationRule) error
}

// REFACTOR-003: Config abstraction - Validation rule definition
// ValidationRule defines validation constraints for configuration fields.
type ValidationRule struct {
	Required     bool                    // Whether the field is required
	Type         string                  // Expected type (string, bool, int, etc.)
	Validator    func(interface{}) error // Custom validation function
	DefaultValue interface{}             // Default value if not provided
	Description  string                  // Human-readable description
}

// REFACTOR-003: Schema separation - Configuration source abstraction
// ConfigSource abstracts different configuration sources (files, environment, etc.).
// It enables loading configuration from multiple sources with consistent interface.
type ConfigSource interface {
	// Load loads configuration data from the source
	Load(path string) (map[string]interface{}, error)

	// Save saves configuration data to the source
	Save(path string, data map[string]interface{}) error

	// Exists checks if the configuration source exists
	Exists(path string) bool

	// GetSourceType returns the type of source (file, env, default)
	GetSourceType() string

	// GetPriority returns the priority of this source for merging
	GetPriority() int
}

// REFACTOR-003: Config abstraction - Generic configuration merging
// ConfigMerger provides generic configuration merging logic that works
// with any configuration schema.
type ConfigMerger interface {
	// MergeConfigs merges two configuration objects
	MergeConfigs(dst, src interface{}) error

	// MergeWithPriority merges multiple configurations with priority ordering
	MergeWithPriority(configs []interface{}, priorities []int) (interface{}, error)

	// GetMergeStrategy returns the current merge strategy
	GetMergeStrategy() MergeStrategy

	// SetMergeStrategy sets the merge strategy
	SetMergeStrategy(strategy MergeStrategy)
}

// REFACTOR-003: Config abstraction - Merge strategy enumeration
// MergeStrategy defines how configuration values should be merged.
type MergeStrategy int

const (
	// OverwriteStrategy overwrites destination values with source values
	OverwriteStrategy MergeStrategy = iota

	// PreserveStrategy preserves destination values, ignoring source values
	PreserveStrategy

	// AppendStrategy appends source values to destination (for slices/arrays)
	AppendStrategy
)

// REFACTOR-003: Schema separation - Schema definition interface
// ConfigSchema defines the structure and validation rules for a specific
// application's configuration schema.
type ConfigSchema interface {
	// GetSchemaName returns the unique name of this schema
	GetSchemaName() string

	// GetSchemaVersion returns the version of this schema
	GetSchemaVersion() string

	// GetDefaultConfig returns a default configuration instance
	GetDefaultConfig() interface{}

	// GetFieldDefinitions returns field definitions for all schema fields
	GetFieldDefinitions() map[string]FieldDefinition

	// ValidateSchema validates a configuration against this schema
	ValidateSchema(config interface{}) error

	// MigrateSchema migrates configuration from old version to new version
	MigrateSchema(oldVersion, newVersion string, config interface{}) (interface{}, error)

	// GetRequiredFields returns a list of required field names
	GetRequiredFields() []string
}

// REFACTOR-003: Schema separation - Field definition structure
// FieldDefinition describes a single configuration field within a schema.
type FieldDefinition struct {
	Name         string                  // Field name
	Type         string                  // Field type (string, bool, int, []string, etc.)
	Required     bool                    // Whether the field is required
	DefaultValue interface{}             // Default value for the field
	Description  string                  // Human-readable description
	Validator    func(interface{}) error // Custom validation function
	Tags         map[string]string       // Additional metadata tags
}

// REFACTOR-003: Config abstraction - Generic configuration provider
// ConfigProvider provides type-safe access to configuration values
// without exposing the underlying configuration structure.
type ConfigProvider interface {
	// GetString returns a string configuration value
	GetString(key string) string

	// GetBool returns a boolean configuration value
	GetBool(key string) bool

	// GetInt returns an integer configuration value
	GetInt(key string) int

	// GetStringSlice returns a string slice configuration value
	GetStringSlice(key string) []string

	// GetDefault returns the default value for a key
	GetDefault(key string) interface{}

	// HasKey checks if a configuration key exists
	HasKey(key string) bool

	// GetAllKeys returns all available configuration keys
	GetAllKeys() []string

	// GetWithDefault returns a value with a fallback default
	GetWithDefault(key string, defaultValue interface{}) interface{}
}

// REFACTOR-003: Config abstraction - Configuration format provider interface
// ConfigFormatProvider provides access to format strings and templates
// used for output formatting (renamed to avoid conflict with existing FormatProvider).
type ConfigFormatProvider interface {
	// GetFormatString returns a printf-style format string
	GetFormatString(operation string) string

	// GetTemplateString returns a template-based format string
	GetTemplateString(operation string) string

	// GetPattern returns a regex pattern for data extraction
	GetPattern(patternType string) string

	// HasFormat checks if a format string exists for an operation
	HasFormat(operation string) bool

	// HasTemplate checks if a template exists for an operation
	HasTemplate(operation string) bool

	// HasPattern checks if a pattern exists for a type
	HasPattern(patternType string) bool
}

// REFACTOR-003: Config abstraction - Status code provider interface
// StatusProvider provides access to status codes used for exit codes
// and error handling.
type StatusProvider interface {
	// GetStatusCode returns a status code for a successful operation
	GetStatusCode(operation string) int

	// GetErrorStatusCode returns a status code for an error type
	GetErrorStatusCode(errorType string) int

	// HasStatusCode checks if a status code exists for an operation
	HasStatusCode(operation string) bool

	// HasErrorStatusCode checks if an error status code exists
	HasErrorStatusCode(errorType string) bool

	// GetAllStatusCodes returns all available status codes
	GetAllStatusCodes() map[string]int
}

// REFACTOR-003: Config abstraction - Path provider interface
// PathProvider provides access to path configuration values
// used throughout the application.
type PathProvider interface {
	// GetArchivePath returns the archive directory path
	GetArchivePath() string

	// GetBackupPath returns the backup directory path
	GetBackupPath() string

	// GetConfigPath returns the configuration file path
	GetConfigPath() string

	// GetTempPath returns the temporary directory path
	GetTempPath() string

	// HasPath checks if a path configuration exists
	HasPath(pathType string) bool

	// GetPath returns a generic path by type
	GetPath(pathType string) string
}

// REFACTOR-003: Backward compatibility - Configuration adapter interface
// ConfigAdapter provides backward compatibility by converting between
// generic configuration providers and legacy configuration structures.
type ConfigAdapter interface {
	// ToLegacyConfig converts to the legacy Config struct
	ToLegacyConfig() *Config

	// FromLegacyConfig creates an adapter from a legacy Config struct
	FromLegacyConfig(cfg *Config) ConfigAdapter

	// GetProvider returns the underlying configuration provider
	GetProvider() ConfigProvider

	// GetSchema returns the configuration schema
	GetSchema() ConfigSchema
}

// REFACTOR-003: Migration support - Configuration migrator interface
// ConfigMigrator handles migration of configuration between different
// schema versions.
type ConfigMigrator interface {
	// CanMigrate checks if migration is possible between versions
	CanMigrate(fromVersion, toVersion string) bool

	// Migrate performs the migration from one version to another
	Migrate(fromVersion, toVersion string, config interface{}) (interface{}, error)

	// GetSupportedVersions returns all supported schema versions
	GetSupportedVersions() []string

	// GetMigrationPath returns the migration path between versions
	GetMigrationPath(fromVersion, toVersion string) ([]string, error)
}

// REFACTOR-003: Config abstraction - Configuration factory interface
// ConfigFactory creates configuration objects and their associated
// providers, validators, and other components.
type ConfigFactory interface {
	// CreateLoader creates a configuration loader for a schema
	CreateLoader(schema ConfigSchema) ConfigLoader

	// CreateValidator creates a validator for a schema
	CreateValidator(schema ConfigSchema) ConfigValidator

	// CreateProvider creates a provider for a configuration
	CreateProvider(config interface{}, schema ConfigSchema) ConfigProvider

	// CreateAdapter creates an adapter for backward compatibility
	CreateAdapter(provider ConfigProvider, schema ConfigSchema) ConfigAdapter

	// CreateMigrator creates a migrator for a schema
	CreateMigrator(schema ConfigSchema) ConfigMigrator
}

// REFACTOR-003: Config abstraction - Configuration registry interface
// ConfigRegistry manages multiple configuration schemas and provides
// schema discovery and registration capabilities.
type ConfigRegistry interface {
	// RegisterSchema registers a new configuration schema
	RegisterSchema(schema ConfigSchema) error

	// GetSchema returns a schema by name
	GetSchema(name string) (ConfigSchema, error)

	// GetSchemaByVersion returns a schema by name and version
	GetSchemaByVersion(name, version string) (ConfigSchema, error)

	// ListSchemas returns all registered schema names
	ListSchemas() []string

	// HasSchema checks if a schema is registered
	HasSchema(name string) bool

	// UnregisterSchema removes a schema from the registry
	UnregisterSchema(name string) error
}
