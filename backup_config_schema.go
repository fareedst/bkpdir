// This file is part of bkpdir
//
// Package main provides the backup application specific configuration schema
// implementation. This demonstrates how the configuration abstraction interfaces
// can be used to create application-specific configuration schemas.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
)

// REFACTOR-003: Schema separation - Backup application schema implementation
// BackupConfigSchema implements ConfigSchema for the backup application.
// It defines the structure, validation rules, and default values for backup configuration.
type BackupConfigSchema struct {
	version string
}

// NewBackupConfigSchema creates a new backup application configuration schema.
func NewBackupConfigSchema() *BackupConfigSchema {
	return &BackupConfigSchema{
		version: "1.0.0",
	}
}

// REFACTOR-003: Schema separation - Schema identification
// GetSchemaName returns the unique name of the backup application schema.
func (s *BackupConfigSchema) GetSchemaName() string {
	return "backup-application"
}

// GetSchemaVersion returns the version of the backup application schema.
func (s *BackupConfigSchema) GetSchemaVersion() string {
	return s.version
}

// REFACTOR-003: Schema separation - Default configuration provision
// GetDefaultConfig returns a default configuration instance for the backup application.
func (s *BackupConfigSchema) GetDefaultConfig() interface{} {
	return DefaultConfig()
}

// REFACTOR-003: Schema separation - Field definitions for backup application
// GetFieldDefinitions returns field definitions for all backup application configuration fields.
func (s *BackupConfigSchema) GetFieldDefinitions() map[string]FieldDefinition {
	return map[string]FieldDefinition{
		// Basic backup settings
		"archive_dir_path": {
			Name:         "archive_dir_path",
			Type:         "string",
			Required:     true,
			DefaultValue: "../.bkpdir",
			Description:  "Directory path for storing archives",
			Validator:    validateDirectoryPath,
			Tags:         map[string]string{"category": "paths"},
		},
		"use_current_dir_name": {
			Name:         "use_current_dir_name",
			Type:         "bool",
			Required:     false,
			DefaultValue: true,
			Description:  "Use current directory name in archive names",
			Tags:         map[string]string{"category": "naming"},
		},
		"exclude_patterns": {
			Name:         "exclude_patterns",
			Type:         "[]string",
			Required:     false,
			DefaultValue: []string{".git/", "vendor/"},
			Description:  "Patterns to exclude from archives",
			Validator:    validateExcludePatterns,
			Tags:         map[string]string{"category": "filtering"},
		},
		"include_git_info": {
			Name:         "include_git_info",
			Type:         "bool",
			Required:     false,
			DefaultValue: false,
			Description:  "Include Git information in archive names",
			Tags:         map[string]string{"category": "git"},
		},
		"backup_dir_path": {
			Name:         "backup_dir_path",
			Type:         "string",
			Required:     true,
			DefaultValue: "../.bkpdir/backups",
			Description:  "Directory path for storing file backups",
			Validator:    validateDirectoryPath,
			Tags:         map[string]string{"category": "paths"},
		},
		// Status codes
		"status_created_archive": {
			Name:         "status_created_archive",
			Type:         "int",
			Required:     false,
			DefaultValue: 0,
			Description:  "Exit code when archive is created successfully",
			Validator:    validateStatusCode,
			Tags:         map[string]string{"category": "status_codes"},
		},
		"status_config_error": {
			Name:         "status_config_error",
			Type:         "int",
			Required:     false,
			DefaultValue: 1,
			Description:  "Exit code for configuration errors",
			Validator:    validateStatusCode,
			Tags:         map[string]string{"category": "status_codes"},
		},
		// Format strings
		"format_created_archive": {
			Name:         "format_created_archive",
			Type:         "string",
			Required:     false,
			DefaultValue: "Created archive: %s\n",
			Description:  "Printf format string for archive creation messages",
			Validator:    validateFormatString,
			Tags:         map[string]string{"category": "formats"},
		},
		"template_created_archive": {
			Name:         "template_created_archive",
			Type:         "string",
			Required:     false,
			DefaultValue: "Created archive: %{path}\n",
			Description:  "Template format string for archive creation messages",
			Validator:    validateTemplateString,
			Tags:         map[string]string{"category": "templates"},
		},
		// Patterns
		"pattern_archive_filename": {
			Name:         "pattern_archive_filename",
			Type:         "string",
			Required:     false,
			DefaultValue: defaultArchivePattern,
			Description:  "Regex pattern for parsing archive filenames",
			Validator:    validateRegexPattern,
			Tags:         map[string]string{"category": "patterns"},
		},
	}
}

// REFACTOR-003: Schema separation - Configuration validation
// ValidateSchema validates a configuration against the backup application schema.
func (s *BackupConfigSchema) ValidateSchema(config interface{}) error {
	cfg, ok := config.(*Config)
	if !ok {
		return fmt.Errorf("invalid configuration type: expected *Config, got %T", config)
	}

	// Validate required fields
	if cfg.ArchiveDirPath == "" {
		return fmt.Errorf("archive_dir_path is required")
	}
	if cfg.BackupDirPath == "" {
		return fmt.Errorf("backup_dir_path is required")
	}

	// Validate field values using field definitions
	fieldDefs := s.GetFieldDefinitions()

	// Validate archive directory path
	if def, exists := fieldDefs["archive_dir_path"]; exists && def.Validator != nil {
		if err := def.Validator(cfg.ArchiveDirPath); err != nil {
			return fmt.Errorf("archive_dir_path validation failed: %w", err)
		}
	}

	// Validate backup directory path
	if def, exists := fieldDefs["backup_dir_path"]; exists && def.Validator != nil {
		if err := def.Validator(cfg.BackupDirPath); err != nil {
			return fmt.Errorf("backup_dir_path validation failed: %w", err)
		}
	}

	// Validate exclude patterns
	if def, exists := fieldDefs["exclude_patterns"]; exists && def.Validator != nil {
		if err := def.Validator(cfg.ExcludePatterns); err != nil {
			return fmt.Errorf("exclude_patterns validation failed: %w", err)
		}
	}

	return nil
}

// REFACTOR-003: Migration support - Schema migration for backup application
// MigrateSchema migrates backup configuration from old version to new version.
func (s *BackupConfigSchema) MigrateSchema(oldVersion, newVersion string, config interface{}) (interface{}, error) {
	// For now, no migration is needed as this is version 1.0.0
	if oldVersion == newVersion {
		return config, nil
	}

	return nil, fmt.Errorf("migration from version %s to %s not supported", oldVersion, newVersion)
}

// GetRequiredFields returns a list of required field names for the backup application.
func (s *BackupConfigSchema) GetRequiredFields() []string {
	var required []string
	for name, def := range s.GetFieldDefinitions() {
		if def.Required {
			required = append(required, name)
		}
	}
	return required
}

// REFACTOR-003: Provider pattern - Backup application configuration provider
// BackupConfigProvider implements ConfigProvider for the backup application.
// It provides type-safe access to backup configuration values.
type BackupConfigProvider struct {
	config *Config
	schema *BackupConfigSchema
}

// NewBackupConfigProvider creates a new backup configuration provider.
func NewBackupConfigProvider(config *Config, schema *BackupConfigSchema) *BackupConfigProvider {
	return &BackupConfigProvider{
		config: config,
		schema: schema,
	}
}

// REFACTOR-003: Config abstraction - Type-safe configuration access
// GetString returns a string configuration value.
func (p *BackupConfigProvider) GetString(key string) string {
	switch key {
	case "archive_dir_path":
		return p.config.ArchiveDirPath
	case "backup_dir_path":
		return p.config.BackupDirPath
	case "format_created_archive":
		return p.config.FormatCreatedArchive
	case "template_created_archive":
		return p.config.TemplateCreatedArchive
	case "pattern_archive_filename":
		return p.config.PatternArchiveFilename
	default:
		return ""
	}
}

// GetBool returns a boolean configuration value.
func (p *BackupConfigProvider) GetBool(key string) bool {
	switch key {
	case "use_current_dir_name":
		return p.config.UseCurrentDirName
	case "use_current_dir_name_for_files":
		return p.config.UseCurrentDirNameForFiles
	case "include_git_info":
		return p.config.IncludeGitInfo
	case "show_git_dirty_status":
		return p.config.ShowGitDirtyStatus
	default:
		return false
	}
}

// GetInt returns an integer configuration value.
func (p *BackupConfigProvider) GetInt(key string) int {
	switch key {
	case "status_created_archive":
		return p.config.StatusCreatedArchive
	case "status_config_error":
		return p.config.StatusConfigError
	case "status_created_backup":
		return p.config.StatusCreatedBackup
	default:
		return 0
	}
}

// GetStringSlice returns a string slice configuration value.
func (p *BackupConfigProvider) GetStringSlice(key string) []string {
	switch key {
	case "exclude_patterns":
		return p.config.ExcludePatterns
	default:
		return nil
	}
}

// GetDefault returns the default value for a key.
func (p *BackupConfigProvider) GetDefault(key string) interface{} {
	fieldDefs := p.schema.GetFieldDefinitions()
	if def, exists := fieldDefs[key]; exists {
		return def.DefaultValue
	}
	return nil
}

// HasKey checks if a configuration key exists.
func (p *BackupConfigProvider) HasKey(key string) bool {
	fieldDefs := p.schema.GetFieldDefinitions()
	_, exists := fieldDefs[key]
	return exists
}

// GetAllKeys returns all available configuration keys.
func (p *BackupConfigProvider) GetAllKeys() []string {
	fieldDefs := p.schema.GetFieldDefinitions()
	keys := make([]string, 0, len(fieldDefs))
	for key := range fieldDefs {
		keys = append(keys, key)
	}
	return keys
}

// GetWithDefault returns a value with a fallback default.
func (p *BackupConfigProvider) GetWithDefault(key string, defaultValue interface{}) interface{} {
	if !p.HasKey(key) {
		return defaultValue
	}

	// Try to get the value based on the expected type
	fieldDefs := p.schema.GetFieldDefinitions()
	if def, exists := fieldDefs[key]; exists {
		switch def.Type {
		case "string":
			if value := p.GetString(key); value != "" {
				return value
			}
		case "bool":
			return p.GetBool(key)
		case "int":
			return p.GetInt(key)
		case "[]string":
			if value := p.GetStringSlice(key); value != nil {
				return value
			}
		}
	}

	return defaultValue
}

// REFACTOR-003: Config abstraction - Validation functions for backup application
// Validation functions for backup application configuration fields

func validateDirectoryPath(value interface{}) error {
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if path == "" {
		return fmt.Errorf("directory path cannot be empty")
	}
	// Additional path validation could be added here
	return nil
}

func validateExcludePatterns(value interface{}) error {
	patterns, ok := value.([]string)
	if !ok {
		return fmt.Errorf("expected []string, got %T", value)
	}
	for _, pattern := range patterns {
		if pattern == "" {
			return fmt.Errorf("exclude pattern cannot be empty")
		}
	}
	return nil
}

func validateStatusCode(value interface{}) error {
	code, ok := value.(int)
	if !ok {
		return fmt.Errorf("expected int, got %T", value)
	}
	if code < 0 || code > 255 {
		return fmt.Errorf("status code must be between 0 and 255, got %d", code)
	}
	return nil
}

func validateFormatString(value interface{}) error {
	format, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if format == "" {
		return fmt.Errorf("format string cannot be empty")
	}
	// Could add printf format validation here
	return nil
}

func validateTemplateString(value interface{}) error {
	template, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if template == "" {
		return fmt.Errorf("template string cannot be empty")
	}
	// Could add template format validation here
	return nil
}

func validateRegexPattern(value interface{}) error {
	pattern, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}
	if pattern == "" {
		return fmt.Errorf("regex pattern cannot be empty")
	}
	// Could add regex compilation validation here
	return nil
}
