// This file is part of bkpdir
//
// Package main provides concrete implementations of configuration interfaces.
// These implementations enable schema-agnostic configuration management.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v3"

	// 🔶 GIT-005: Import Git package for configuration integration
	"bkpdir/pkg/git"
)

// 🔻 REFACTOR-003: Config abstraction - Default configuration loader implementation - 🔧
// DefaultConfigLoader provides the default implementation of ConfigLoader interface.
// This implementation maintains backward compatibility while enabling schema abstraction.
type DefaultConfigLoader struct {
	fileOps   ConfigFileOperations
	merger    ConfigMerger
	validator ConfigValidator
}

// NewDefaultConfigLoader creates a new DefaultConfigLoader with default components.
func NewDefaultConfigLoader() *DefaultConfigLoader {
	return &DefaultConfigLoader{
		fileOps:   &FileSystemOperations{},
		merger:    &DefaultConfigMerger{},
		validator: &BackupAppValidator{},
	}
}

// LoadConfig loads configuration from the specified root directory.
// 🔻 REFACTOR-003: Config abstraction - Schema-agnostic configuration loading - 🔧
func (d *DefaultConfigLoader) LoadConfig(root string) (*Config, error) {
	cfg := DefaultConfig()

	paths := d.merger.GetConfigSearchPaths()
	for _, path := range paths {
		expandedPath := d.merger.ExpandPath(path)
		if d.fileOps.FileExists(expandedPath) {
			data, err := d.fileOps.ReadFile(expandedPath)
			if err != nil {
				continue
			}

			var fileCfg Config
			if err := yaml.Unmarshal(data, &fileCfg); err != nil {
				continue
			}

			d.merger.MergeConfigs(cfg, &fileCfg)
		}
	}

	return cfg, d.validator.ValidateSchema(cfg)
}

// LoadConfigValues loads configuration values with source tracking.
// 🔻 REFACTOR-003: Config abstraction - Source tracking for configuration values - 📝
func (d *DefaultConfigLoader) LoadConfigValues(root string) (map[string]ConfigValue, error) {
	cfg, err := d.LoadConfig(root)
	if err != nil {
		return nil, err
	}

	values := make(map[string]ConfigValue)
	for _, cv := range d.GetConfigValuesWithSources(cfg, root) {
		values[cv.Name] = cv
	}

	return values, nil
}

// GetConfigValues extracts configuration values from a Config struct.
// 🔻 REFACTOR-003: Config abstraction - Generic configuration value extraction - 📝
func (d *DefaultConfigLoader) GetConfigValues(cfg *Config) []ConfigValue {
	return GetConfigValues(cfg)
}

// GetConfigValuesWithSources extracts configuration values with source information.
// 🔻 REFACTOR-003: Config abstraction - Source-aware configuration value extraction - 📝
func (d *DefaultConfigLoader) GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue {
	return GetConfigValuesWithSources(cfg, root)
}

// ValidateConfig validates a configuration structure.
// 🔻 REFACTOR-003: Config abstraction - Pluggable configuration validation - 🔧
func (d *DefaultConfigLoader) ValidateConfig(cfg *Config) error {
	return d.validator.ValidateSchema(cfg)
}

// 🔻 REFACTOR-003: Config abstraction - Default configuration merger implementation - 🔧
// DefaultConfigMerger provides the default implementation of ConfigMerger interface.
// This implementation maintains backward compatibility while enabling schema abstraction.
type DefaultConfigMerger struct{}

// MergeConfigs merges source configuration into destination configuration.
// 🔻 REFACTOR-003: Config abstraction - Schema-agnostic configuration merging - 🔧
func (d *DefaultConfigMerger) MergeConfigs(dst, src *Config) {
	mergeConfigs(dst, src)
}

// MergeConfigValues merges configuration value maps.
// 🔻 REFACTOR-003: Config abstraction - Configuration value map merging - 📝
func (d *DefaultConfigMerger) MergeConfigValues(dst, src map[string]ConfigValue) {
	for key, value := range src {
		dst[key] = value
	}
}

// GetConfigSearchPaths returns the search paths for configuration files.
// 🔻 REFACTOR-003: Config abstraction - Configurable search path logic - 📝
func (d *DefaultConfigMerger) GetConfigSearchPaths() []string {
	return getConfigSearchPaths()
}

// ExpandPath expands path variables and returns absolute path.
// 🔻 REFACTOR-003: Config abstraction - Path expansion abstraction - 📝
func (d *DefaultConfigMerger) ExpandPath(path string) string {
	return expandPath(path)
}

// 🔻 REFACTOR-003: Config abstraction - File configuration source implementation - 🔧
// FileConfigSource provides file-based configuration loading.
// This implementation abstracts file-based configuration from other sources.
type FileConfigSource struct {
	path    string
	fileOps ConfigFileOperations
}

// NewFileConfigSource creates a new file-based configuration source.
func NewFileConfigSource(path string) *FileConfigSource {
	return &FileConfigSource{
		path:    path,
		fileOps: &FileSystemOperations{},
	}
}

// LoadFromFile loads configuration from a file source.
// 🔻 REFACTOR-003: Config abstraction - File-based configuration loading - 🔧
func (f *FileConfigSource) LoadFromFile(path string) (*Config, error) {
	if !f.fileOps.FileExists(path) {
		return DefaultConfig(), nil
	}

	data, err := f.fileOps.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return &cfg, nil
}

// LoadFromEnvironment loads configuration from environment variables.
// 🔻 REFACTOR-003: Config abstraction - Environment-based configuration loading - 🔧
func (f *FileConfigSource) LoadFromEnvironment() (*Config, error) {
	cfg := DefaultConfig()

	// Load environment variable overrides
	if archiveDir := os.Getenv("BKPDIR_ARCHIVE_DIR"); archiveDir != "" {
		cfg.ArchiveDirPath = archiveDir
	}
	if backupDir := os.Getenv("BKPDIR_BACKUP_DIR"); backupDir != "" {
		cfg.BackupDirPath = backupDir
	}
	if includeGit := os.Getenv("BKPDIR_INCLUDE_GIT"); includeGit != "" {
		cfg.IncludeGitInfo = strings.ToLower(includeGit) == "true"
	}

	// 🔶 GIT-005: Git configuration environment variable support - 📝
	// Git configuration overrides from environment variables
	if cfg.Git == nil {
		cfg.Git = DefaultGitConfig()
	}

	if enabled := os.Getenv("BKPDIR_GIT_ENABLED"); enabled != "" {
		cfg.Git.Enabled = strings.ToLower(enabled) == "true"
	}
	if includeInfo := os.Getenv("BKPDIR_GIT_INCLUDE_INFO"); includeInfo != "" {
		cfg.Git.IncludeInfo = strings.ToLower(includeInfo) == "true"
	}
	if showDirty := os.Getenv("BKPDIR_GIT_SHOW_DIRTY_STATUS"); showDirty != "" {
		cfg.Git.ShowDirtyStatus = strings.ToLower(showDirty) == "true"
	}
	if command := os.Getenv("BKPDIR_GIT_COMMAND"); command != "" {
		cfg.Git.Command = command
	}
	if workingDir := os.Getenv("BKPDIR_GIT_WORKING_DIRECTORY"); workingDir != "" {
		cfg.Git.WorkingDirectory = workingDir
	}
	if includeSubmodules := os.Getenv("BKPDIR_GIT_INCLUDE_SUBMODULES"); includeSubmodules != "" {
		cfg.Git.IncludeSubmodules = strings.ToLower(includeSubmodules) == "true"
	}

	return cfg, nil
}

// LoadDefaults returns the default configuration.
// 🔻 REFACTOR-003: Config abstraction - Default configuration provider - 📝
func (f *FileConfigSource) LoadDefaults() *Config {
	return DefaultConfig()
}

// GetSourceName returns the name of this configuration source.
// 🔻 REFACTOR-003: Config abstraction - Configuration source identification - 📝
func (f *FileConfigSource) GetSourceName() string {
	if f.path != "" {
		return fmt.Sprintf("file:%s", f.path)
	}
	return "file"
}

// IsAvailable checks if this configuration source is available.
// 🔻 REFACTOR-003: Config abstraction - Configuration source availability - 🔧
func (f *FileConfigSource) IsAvailable() bool {
	return f.path == "" || f.fileOps.FileExists(f.path)
}

// 🔻 REFACTOR-003: Config abstraction - Backup application validator implementation - 🔧
// BackupAppValidator provides validation specific to the backup application schema.
// This implementation demonstrates schema-specific validation while using the common interface.
type BackupAppValidator struct{}

// ValidateSchema validates the configuration against the expected schema.
// 🔻 REFACTOR-003: Config abstraction - Backup application schema validation - 🔧
func (b *BackupAppValidator) ValidateSchema(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	// Validate required paths
	if cfg.ArchiveDirPath == "" {
		return fmt.Errorf("archive_dir_path cannot be empty")
	}
	if cfg.BackupDirPath == "" {
		return fmt.Errorf("backup_dir_path cannot be empty")
	}

	// Validate verification configuration
	if cfg.Verification != nil {
		if cfg.Verification.ChecksumAlgorithm != "" {
			allowed := []string{"sha256", "md5", "sha1"}
			valid := false
			for _, alg := range allowed {
				if cfg.Verification.ChecksumAlgorithm == alg {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid checksum algorithm: %s", cfg.Verification.ChecksumAlgorithm)
			}
		}
	}

	return nil
}

// ValidateValues validates individual configuration values.
// 🔻 REFACTOR-003: Config abstraction - Individual value validation - 📝
func (b *BackupAppValidator) ValidateValues(values map[string]ConfigValue) error {
	for name, value := range values {
		if err := b.validateValue(name, value.Value); err != nil {
			return fmt.Errorf("invalid value for %s: %w", name, err)
		}
	}
	return nil
}

// GetRequiredFields returns the list of required configuration fields.
// 🔻 REFACTOR-003: Config abstraction - Required field specification - 📝
func (b *BackupAppValidator) GetRequiredFields() []string {
	return []string{
		"archive_dir_path",
		"backup_dir_path",
	}
}

// GetValidationRules returns validation rules for configuration fields.
// 🔻 REFACTOR-003: Config abstraction - Validation rule specification - 📝
func (b *BackupAppValidator) GetValidationRules() map[string]ValidationRule {
	return map[string]ValidationRule{
		"archive_dir_path": {
			Required:  true,
			Type:      "string",
			MinLength: 1,
		},
		"backup_dir_path": {
			Required:  true,
			Type:      "string",
			MinLength: 1,
		},
		"exclude_patterns": {
			Type: "[]string",
		},
		"verification.checksum_algorithm": {
			Type:        "string",
			ValidValues: []string{"sha256", "md5", "sha1"},
		},
	}
}

// validateValue validates an individual configuration value.
// 🔻 REFACTOR-003: Config abstraction - Single value validation logic - 🔧
func (b *BackupAppValidator) validateValue(name, value string) error {
	rules := b.GetValidationRules()
	rule, exists := rules[name]
	if !exists {
		return nil // No validation rule for this field
	}

	if rule.Required && value == "" {
		return fmt.Errorf("required field cannot be empty")
	}

	if rule.MinLength > 0 && len(value) < rule.MinLength {
		return fmt.Errorf("value must be at least %d characters", rule.MinLength)
	}

	if rule.MaxLength > 0 && len(value) > rule.MaxLength {
		return fmt.Errorf("value must be at most %d characters", rule.MaxLength)
	}

	if rule.Pattern != "" {
		matched, err := regexp.MatchString(rule.Pattern, value)
		if err != nil {
			return fmt.Errorf("invalid pattern: %w", err)
		}
		if !matched {
			return fmt.Errorf("value does not match required pattern")
		}
	}

	if len(rule.ValidValues) > 0 {
		valid := false
		for _, validValue := range rule.ValidValues {
			if value == validValue {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("value must be one of: %s", strings.Join(rule.ValidValues, ", "))
		}
	}

	return nil
}

// 🔻 REFACTOR-003: Config abstraction - Application-specific configuration implementation - 🔧
// BackupApplicationConfig provides access to backup-specific configuration settings.
// This implementation demonstrates schema separation from generic configuration operations.
type BackupApplicationConfig struct {
	cfg *Config
}

// NewBackupApplicationConfig creates a new backup application configuration wrapper.
func NewBackupApplicationConfig(cfg *Config) *BackupApplicationConfig {
	return &BackupApplicationConfig{cfg: cfg}
}

// GetArchiveSettings returns archive-related configuration.
// 🔻 REFACTOR-003: Schema separation - Archive settings extraction - 📝
func (b *BackupApplicationConfig) GetArchiveSettings() ArchiveSettings {
	return ArchiveSettings{
		DirectoryPath:      b.cfg.ArchiveDirPath,
		UseCurrentDirName:  b.cfg.UseCurrentDirName,
		ExcludePatterns:    b.cfg.ExcludePatterns,
		IncludeGitInfo:     b.cfg.IncludeGitInfo,
		ShowGitDirtyStatus: b.cfg.ShowGitDirtyStatus,
		Verification:       b.cfg.Verification,
	}
}

// GetBackupSettings returns backup-related configuration.
// 🔻 REFACTOR-003: Schema separation - Backup settings extraction - 📝
func (b *BackupApplicationConfig) GetBackupSettings() BackupSettings {
	return BackupSettings{
		DirectoryPath:             b.cfg.BackupDirPath,
		UseCurrentDirNameForFiles: b.cfg.UseCurrentDirNameForFiles,
	}
}

// GetStatusCodes returns application status codes.
// 🔻 REFACTOR-003: Schema separation - Status code extraction - 📝
func (b *BackupApplicationConfig) GetStatusCodes() map[string]int {
	return b.cfg.GetStatusCodes()
}

// GetFormatSettings returns output formatting configuration.
// 🔻 REFACTOR-003: Schema separation - Format settings extraction - 📝
func (b *BackupApplicationConfig) GetFormatSettings() FormatSettings {
	return FormatSettings{
		FormatStrings: map[string]string{
			"created_archive":   b.cfg.FormatCreatedArchive,
			"identical_archive": b.cfg.FormatIdenticalArchive,
			"list_archive":      b.cfg.FormatListArchive,
			"config_value":      b.cfg.FormatConfigValue,
			"dry_run_archive":   b.cfg.FormatDryRunArchive,
			"error":             b.cfg.FormatError,
			"created_backup":    b.cfg.FormatCreatedBackup,
			"identical_backup":  b.cfg.FormatIdenticalBackup,
			"list_backup":       b.cfg.FormatListBackup,
			"dry_run_backup":    b.cfg.FormatDryRunBackup,
		},
		TemplateStrings: map[string]string{
			"created_archive":   b.cfg.TemplateCreatedArchive,
			"identical_archive": b.cfg.TemplateIdenticalArchive,
			"list_archive":      b.cfg.TemplateListArchive,
			"config_value":      b.cfg.TemplateConfigValue,
			"dry_run_archive":   b.cfg.TemplateDryRunArchive,
			"error":             b.cfg.TemplateError,
			"created_backup":    b.cfg.TemplateCreatedBackup,
			"identical_backup":  b.cfg.TemplateIdenticalBackup,
			"list_backup":       b.cfg.TemplateListBackup,
			"dry_run_backup":    b.cfg.TemplateDryRunBackup,
		},
		PatternStrings: map[string]string{
			"archive_filename": b.cfg.PatternArchiveFilename,
			"backup_filename":  b.cfg.PatternBackupFilename,
			"config_line":      b.cfg.PatternConfigLine,
			"timestamp":        b.cfg.PatternTimestamp,
		},
		ErrorFormatStrings: b.cfg.GetErrorFormatStrings(),
	}
}

// 🔻 REFACTOR-003: Config abstraction - File system operations implementation - 🔧
// FileSystemOperations provides concrete file system operations for configuration.
// This implementation enables testing and different storage backends.
type FileSystemOperations struct{}

// FileExists checks if a configuration file exists.
// 🔻 REFACTOR-003: Config abstraction - File existence checking - 📝
func (f *FileSystemOperations) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads configuration file contents.
// 🔻 REFACTOR-003: Config abstraction - File content reading - 📝
func (f *FileSystemOperations) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes configuration file contents.
// 🔻 REFACTOR-003: Config abstraction - File content writing - 📝
func (f *FileSystemOperations) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// GetFileInfo returns file information for configuration files.
// 🔻 REFACTOR-003: Config abstraction - File information retrieval - 📝
func (f *FileSystemOperations) GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// 🔻 REFACTOR-003: Config abstraction - Configuration source determiner implementation - 🔧
// DefaultSourceDeterminer provides default source determination logic.
// This implementation enables source tracking across different configuration providers.
type DefaultSourceDeterminer struct {
	configSource string
}

// NewDefaultSourceDeterminer creates a new source determiner.
func NewDefaultSourceDeterminer(configSource string) *DefaultSourceDeterminer {
	return &DefaultSourceDeterminer{configSource: configSource}
}

// DetermineSource determines the source of a configuration value.
// 🔻 REFACTOR-003: Config abstraction - Configuration value source tracking - 📝
func (d *DefaultSourceDeterminer) DetermineSource(current, defaultValue interface{}) string {
	if current == defaultValue {
		return "default"
	}
	if d.configSource != "" {
		return d.configSource
	}
	return "computed"
}

// GetConfigSource returns the primary configuration source.
// 🔻 REFACTOR-003: Config abstraction - Primary source identification - 📝
func (d *DefaultSourceDeterminer) GetConfigSource() string {
	return d.configSource
}

// GetSourcePriority returns the priority order of configuration sources.
// 🔻 REFACTOR-003: Config abstraction - Source priority definition - 📝
func (d *DefaultSourceDeterminer) GetSourcePriority() []string {
	return []string{"file", "environment", "default"}
}

// 🔻 REFACTOR-003: Config abstraction - Configuration value extractor implementation - 🔧
// DefaultValueExtractor provides default configuration value extraction logic.
// This implementation enables schema-agnostic value extraction for different application types.
type DefaultValueExtractor struct{}

// ExtractBasicValues extracts basic configuration values.
// 🔻 REFACTOR-003: Config abstraction - Basic value extraction - 📝
func (d *DefaultValueExtractor) ExtractBasicValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return getBasicConfigValues(cfg, defaultCfg, getSource)
}

// ExtractStatusCodeValues extracts status code configuration values.
// 🔻 REFACTOR-003: Config abstraction - Status code value extraction - 📝
func (d *DefaultValueExtractor) ExtractStatusCodeValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return getStatusCodeValues(cfg, defaultCfg, getSource)
}

// ExtractVerificationValues extracts verification configuration values.
// 🔻 REFACTOR-003: Config abstraction - Verification value extraction - 📝
func (d *DefaultValueExtractor) ExtractVerificationValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return getVerificationValues(cfg, defaultCfg, getSource)
}

// ExtractFormatValues extracts format string configuration values.
// 🔻 REFACTOR-003: Config abstraction - Format value extraction - 📝
func (d *DefaultValueExtractor) ExtractFormatValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	// This would extract format strings - implementation depends on schema
	// For now, return empty slice as format extraction is complex
	return []ConfigValue{}
}

// 🔶 GIT-005: Git configuration conversion - 📝
// ToGitConfig converts the main application's GitConfig to pkg/git Config
func (gc *GitConfig) ToGitConfig() *git.Config {
	if gc == nil {
		return git.DefaultConfig()
	}

	gitConfig := &git.Config{
		// Basic Git integration settings
		Enabled:         gc.Enabled,
		IncludeInfo:     gc.IncludeInfo,
		ShowDirtyStatus: gc.ShowDirtyStatus,

		// Git command configuration
		Command:          gc.Command,
		WorkingDirectory: gc.WorkingDirectory,

		// Git behavior settings
		RequireCleanRepo:  gc.RequireCleanRepo,
		AutoDetectRepo:    gc.AutoDetectRepo,
		IncludeSubmodules: gc.IncludeSubmodules,

		// Git information inclusion
		IncludeBranch: gc.IncludeBranch,
		IncludeHash:   gc.IncludeHash,
		IncludeStatus: gc.IncludeStatus,

		// Git command timeouts and limits
		CommandTimeout:    gc.CommandTimeout,
		MaxSubmoduleDepth: gc.MaxSubmoduleDepth,

		// Legacy compatibility fields
		IncludeDirtyStatus: gc.ShowDirtyStatus, // Map to legacy field
		GitCommand:         gc.Command,         // Map to legacy field
	}

	return gitConfig
}

// 🔶 GIT-005: Git configuration integration helper - 📝
// GetGitConfig returns a properly configured pkg/git Config from the application config
func GetGitConfig(cfg *Config) *git.Config {
	if cfg.Git != nil {
		return cfg.Git.ToGitConfig()
	}

	// Fallback to legacy configuration fields
	gitConfig := git.DefaultConfig()
	gitConfig.IncludeInfo = cfg.IncludeGitInfo
	gitConfig.ShowDirtyStatus = cfg.ShowGitDirtyStatus
	gitConfig.IncludeDirtyStatus = cfg.ShowGitDirtyStatus

	return gitConfig
}
