// This file is part of bkpdir
//
// Package main provides configuration loading, merging, and validation for BkpDir.
// It handles loading, merging, and validating configuration from YAML files
// and environment variables.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// VerificationConfig defines settings for archive verification.
type VerificationConfig struct {
	VerifyOnCreate    bool   `yaml:"verify_on_create"`
	ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}

// Config holds all configuration settings for the BkpDir application.
// It includes settings for archive creation, file backup, status codes,
// and output formatting.
type Config struct {
	// Basic settings
	ArchiveDirPath    string              `yaml:"archive_dir_path"`
	UseCurrentDirName bool                `yaml:"use_current_dir_name"`
	ExcludePatterns   []string            `yaml:"exclude_patterns"`
	IncludeGitInfo    bool                `yaml:"include_git_info"`
	Verification      *VerificationConfig `yaml:"verification"`

	// File backup settings
	BackupDirPath             string `yaml:"backup_dir_path"`
	UseCurrentDirNameForFiles bool   `yaml:"use_current_dir_name_for_files"`

	// Status codes for directory operations
	StatusCreatedArchive                        int `yaml:"status_created_archive"`
	StatusFailedToCreateArchiveDirectory        int `yaml:"status_failed_to_create_archive_directory"`
	StatusDirectoryIsIdenticalToExistingArchive int `yaml:"status_directory_is_identical_to_existing_archive"`
	StatusDirectoryNotFound                     int `yaml:"status_directory_not_found"`
	StatusInvalidDirectoryType                  int `yaml:"status_invalid_directory_type"`
	StatusPermissionDenied                      int `yaml:"status_permission_denied"`
	StatusDiskFull                              int `yaml:"status_disk_full"`
	StatusConfigError                           int `yaml:"status_config_error"`

	// Status codes for file operations
	StatusCreatedBackup                   int `yaml:"status_created_backup"`
	StatusFailedToCreateBackupDirectory   int `yaml:"status_failed_to_create_backup_directory"`
	StatusFileIsIdenticalToExistingBackup int `yaml:"status_file_is_identical_to_existing_backup"`
	StatusFileNotFound                    int `yaml:"status_file_not_found"`
	StatusInvalidFileType                 int `yaml:"status_invalid_file_type"`

	// Printf-style format strings for directory operations
	FormatCreatedArchive   string `yaml:"format_created_archive"`
	FormatIdenticalArchive string `yaml:"format_identical_archive"`
	FormatListArchive      string `yaml:"format_list_archive"`
	FormatConfigValue      string `yaml:"format_config_value"`
	FormatDryRunArchive    string `yaml:"format_dry_run_archive"`
	FormatError            string `yaml:"format_error"`

	// Printf-style format strings for file operations
	FormatCreatedBackup   string `yaml:"format_created_backup"`
	FormatIdenticalBackup string `yaml:"format_identical_backup"`
	FormatListBackup      string `yaml:"format_list_backup"`
	FormatDryRunBackup    string `yaml:"format_dry_run_backup"`

	// Template-based format strings for directory operations
	TemplateCreatedArchive   string `yaml:"template_created_archive"`
	TemplateIdenticalArchive string `yaml:"template_identical_archive"`
	TemplateListArchive      string `yaml:"template_list_archive"`
	TemplateConfigValue      string `yaml:"template_config_value"`
	TemplateDryRunArchive    string `yaml:"template_dry_run_archive"`
	TemplateError            string `yaml:"template_error"`

	// Template-based format strings for file operations
	TemplateCreatedBackup   string `yaml:"template_created_backup"`
	TemplateIdenticalBackup string `yaml:"template_identical_backup"`
	TemplateListBackup      string `yaml:"template_list_backup"`
	TemplateDryRunBackup    string `yaml:"template_dry_run_backup"`

	// Regex patterns
	PatternArchiveFilename string `yaml:"pattern_archive_filename"`
	PatternBackupFilename  string `yaml:"pattern_backup_filename"`
	PatternConfigLine      string `yaml:"pattern_config_line"`
	PatternTimestamp       string `yaml:"pattern_timestamp"`
}

// ConfigValue represents a single configuration value with its source.
type ConfigValue struct {
	Name   string
	Value  string
	Source string
}

// Default regex patterns
const (
	defaultArchivePattern = `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-` +
		`(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`
	defaultBackupPattern = `(?P<filename>[^/]+)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-` +
		`(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<note>.+))?`
	defaultConfigPattern    = `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`
	defaultTimestampPattern = `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+` +
		`(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`
)

// DefaultConfig returns a new Config instance with default values.
func DefaultConfig() *Config {
	return &Config{
		// Basic settings
		ArchiveDirPath:    "../.bkpdir",
		UseCurrentDirName: true,
		ExcludePatterns:   []string{".git/", "vendor/"},
		IncludeGitInfo:    true,
		Verification: &VerificationConfig{
			VerifyOnCreate:    false,
			ChecksumAlgorithm: "sha256",
		},

		// File backup settings
		BackupDirPath:             "../.bkpdir",
		UseCurrentDirNameForFiles: true,

		// Status codes for directory operations
		StatusCreatedArchive:                        0,
		StatusFailedToCreateArchiveDirectory:        31,
		StatusDirectoryIsIdenticalToExistingArchive: 0,
		StatusDirectoryNotFound:                     20,
		StatusInvalidDirectoryType:                  21,
		StatusPermissionDenied:                      22,
		StatusDiskFull:                              30,
		StatusConfigError:                           10,

		// Status codes for file operations
		StatusCreatedBackup:                   0,
		StatusFailedToCreateBackupDirectory:   31,
		StatusFileIsIdenticalToExistingBackup: 0,
		StatusFileNotFound:                    20,
		StatusInvalidFileType:                 21,

		// Printf-style format strings for directory operations
		FormatCreatedArchive:   "Created archive: %s\n",
		FormatIdenticalArchive: "Directory is identical to existing archive: %s\n",
		FormatListArchive:      "%s (created: %s)\n",
		FormatConfigValue:      "%s: %s (source: %s)\n",
		FormatDryRunArchive:    "Would create archive: %s\n",
		FormatError:            "Error: %s\n",

		// Printf-style format strings for file operations
		FormatCreatedBackup:   "Created backup: %s\n",
		FormatIdenticalBackup: "File is identical to existing backup: %s\n",
		FormatListBackup:      "%s (created: %s)\n",
		FormatDryRunBackup:    "Would create backup: %s\n",

		// Template-based format strings for directory operations
		TemplateCreatedArchive:   "Created archive: %{path}\n",
		TemplateIdenticalArchive: "Directory is identical to existing archive: %{path}\n",
		TemplateListArchive:      "%{path} (created: %{creation_time})\n",
		TemplateConfigValue:      "%{name}: %{value} (source: %{source})\n",
		TemplateDryRunArchive:    "Would create archive: %{path}\n",
		TemplateError:            "Error: %{message}\n",

		// Template-based format strings for file operations
		TemplateCreatedBackup:   "Created backup: %{path}\n",
		TemplateIdenticalBackup: "File is identical to existing backup: %{path}\n",
		TemplateListBackup:      "%{path} (created: %{creation_time})\n",
		TemplateDryRunBackup:    "Would create backup: %{path}\n",

		// Regex patterns
		PatternArchiveFilename: defaultArchivePattern,
		PatternBackupFilename:  defaultBackupPattern,
		PatternConfigLine:      defaultConfigPattern,
		PatternTimestamp:       defaultTimestampPattern,
	}
}

func getConfigSearchPaths() []string {
	// Check BKPDIR_CONFIG environment variable
	if configPaths := os.Getenv("BKPDIR_CONFIG"); configPaths != "" {
		return strings.Split(configPaths, ":")
	}

	// Default search path
	return []string{"./.bkpdir.yml", "~/.bkpdir.yml"}
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// LoadConfig loads and merges configuration from YAML files in the search path.
// It processes configuration files in order, with earlier files taking precedence.
// The root parameter is used to resolve relative paths in the configuration.
// Returns a Config instance with merged values from all valid configuration files.
func LoadConfig(root string) (*Config, error) {
	cfg := DefaultConfig()
	searchPaths := getConfigSearchPaths()

	// Process configuration files in order (earlier files take precedence)
	for _, configPath := range searchPaths {
		expandedPath := expandPath(configPath)

		// Make relative paths relative to root directory
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
		}

		if _, err := os.Stat(expandedPath); err == nil {
			f, err := os.Open(expandedPath)
			if err != nil {
				continue // Skip files we can't open
			}
			defer f.Close()

			// Create a temporary config to load into
			tempCfg := DefaultConfig()
			d := yaml.NewDecoder(f)
			if err := d.Decode(tempCfg); err != nil {
				f.Close()
				continue // Skip files with invalid YAML
			}
			f.Close()

			// Merge non-zero values from tempCfg into cfg
			mergeConfigs(cfg, tempCfg)
			break // Use first valid config file found
		}
	}

	return cfg, nil
}

// mergeConfigs merges non-default values from src into dst.
// It uses helper functions to reduce cyclomatic complexity.
func mergeConfigs(dst, src *Config) {
	mergeBasicSettings(dst, src)
	mergeFileBackupSettings(dst, src)
	mergeStatusCodes(dst, src)
	mergeFormatStrings(dst, src)
	mergeTemplates(dst, src)
	mergePatterns(dst, src)
}

// mergeBasicSettings merges basic configuration settings.
func mergeBasicSettings(dst, src *Config) {
	if src.ArchiveDirPath != DefaultConfig().ArchiveDirPath {
		dst.ArchiveDirPath = src.ArchiveDirPath
	}
	if src.UseCurrentDirName != DefaultConfig().UseCurrentDirName {
		dst.UseCurrentDirName = src.UseCurrentDirName
	}
	if len(src.ExcludePatterns) > 0 && !equalStringSlices(src.ExcludePatterns, DefaultConfig().ExcludePatterns) {
		dst.ExcludePatterns = src.ExcludePatterns
	}
	if src.IncludeGitInfo != DefaultConfig().IncludeGitInfo {
		dst.IncludeGitInfo = src.IncludeGitInfo
	}
	if src.Verification != nil {
		dst.Verification = src.Verification
	}
}

// mergeFileBackupSettings merges file backup configuration settings.
func mergeFileBackupSettings(dst, src *Config) {
	if src.BackupDirPath != DefaultConfig().BackupDirPath {
		dst.BackupDirPath = src.BackupDirPath
	}
	if src.UseCurrentDirNameForFiles != DefaultConfig().UseCurrentDirNameForFiles {
		dst.UseCurrentDirNameForFiles = src.UseCurrentDirNameForFiles
	}
}

// mergeStatusCodes merges status code configuration settings.
func mergeStatusCodes(dst, src *Config) {
	mergeDirectoryStatusCodes(dst, src)
	mergeFileStatusCodes(dst, src)
}

// mergeDirectoryStatusCodes merges directory operation status codes.
func mergeDirectoryStatusCodes(dst, src *Config) {
	statusCodes := map[string]struct {
		src *int
		dst *int
	}{
		"created_archive": {
			&src.StatusCreatedArchive,
			&dst.StatusCreatedArchive,
		},
		"failed_to_create_archive_directory": {
			&src.StatusFailedToCreateArchiveDirectory,
			&dst.StatusFailedToCreateArchiveDirectory,
		},
		"directory_is_identical_to_existing": {
			&src.StatusDirectoryIsIdenticalToExistingArchive,
			&dst.StatusDirectoryIsIdenticalToExistingArchive,
		},
		"directory_not_found": {
			&src.StatusDirectoryNotFound,
			&dst.StatusDirectoryNotFound,
		},
		"invalid_directory_type": {
			&src.StatusInvalidDirectoryType,
			&dst.StatusInvalidDirectoryType,
		},
		"permission_denied": {
			&src.StatusPermissionDenied,
			&dst.StatusPermissionDenied,
		},
		"disk_full": {
			&src.StatusDiskFull,
			&dst.StatusDiskFull,
		},
		"config_error": {
			&src.StatusConfigError,
			&dst.StatusConfigError,
		},
	}

	for _, codes := range statusCodes {
		if *codes.src != *codes.dst && *codes.src != 0 {
			*codes.dst = *codes.src
		}
	}
}

// mergeFileStatusCodes merges file operation status codes.
func mergeFileStatusCodes(dst, src *Config) {
	statusCodes := map[string]struct {
		src *int
		dst *int
	}{
		"created_backup": {
			&src.StatusCreatedBackup,
			&dst.StatusCreatedBackup,
		},
		"failed_to_create_backup_directory": {
			&src.StatusFailedToCreateBackupDirectory,
			&dst.StatusFailedToCreateBackupDirectory,
		},
		"file_is_identical_to_existing": {
			&src.StatusFileIsIdenticalToExistingBackup,
			&dst.StatusFileIsIdenticalToExistingBackup,
		},
		"file_not_found": {
			&src.StatusFileNotFound,
			&dst.StatusFileNotFound,
		},
		"invalid_file_type": {
			&src.StatusInvalidFileType,
			&dst.StatusInvalidFileType,
		},
	}

	for _, codes := range statusCodes {
		if *codes.src != *codes.dst && *codes.src != 0 {
			*codes.dst = *codes.src
		}
	}
}

// mergeFormatStrings merges printf-style format string settings.
func mergeFormatStrings(dst, src *Config) {
	mergeDirectoryFormatStrings(dst, src)
	mergeFileFormatStrings(dst, src)
}

// mergeDirectoryFormatStrings merges directory operation format strings.
func mergeDirectoryFormatStrings(dst, src *Config) {
	formats := map[string]struct {
		src *string
		dst *string
	}{
		"created_archive": {
			&src.FormatCreatedArchive,
			&dst.FormatCreatedArchive,
		},
		"identical_archive": {
			&src.FormatIdenticalArchive,
			&dst.FormatIdenticalArchive,
		},
		"list_archive": {
			&src.FormatListArchive,
			&dst.FormatListArchive,
		},
		"config_value": {
			&src.FormatConfigValue,
			&dst.FormatConfigValue,
		},
		"dry_run_archive": {
			&src.FormatDryRunArchive,
			&dst.FormatDryRunArchive,
		},
		"error": {
			&src.FormatError,
			&dst.FormatError,
		},
	}

	for _, format := range formats {
		if *format.src != *format.dst && *format.src != "" {
			*format.dst = *format.src
		}
	}
}

// mergeFileFormatStrings merges file operation format strings.
func mergeFileFormatStrings(dst, src *Config) {
	formats := map[string]struct {
		src *string
		dst *string
	}{
		"created_backup": {
			&src.FormatCreatedBackup,
			&dst.FormatCreatedBackup,
		},
		"identical_backup": {
			&src.FormatIdenticalBackup,
			&dst.FormatIdenticalBackup,
		},
		"list_backup": {
			&src.FormatListBackup,
			&dst.FormatListBackup,
		},
		"dry_run_backup": {
			&src.FormatDryRunBackup,
			&dst.FormatDryRunBackup,
		},
	}

	for _, format := range formats {
		if *format.src != *format.dst && *format.src != "" {
			*format.dst = *format.src
		}
	}
}

// mergeTemplates merges template-based format string settings.
func mergeTemplates(dst, src *Config) {
	mergeDirectoryTemplates(dst, src)
	mergeFileTemplates(dst, src)
}

// mergeDirectoryTemplates merges directory operation templates.
func mergeDirectoryTemplates(dst, src *Config) {
	templates := map[string]struct {
		src *string
		dst *string
	}{
		"created_archive": {
			&src.TemplateCreatedArchive,
			&dst.TemplateCreatedArchive,
		},
		"identical_archive": {
			&src.TemplateIdenticalArchive,
			&dst.TemplateIdenticalArchive,
		},
		"list_archive": {
			&src.TemplateListArchive,
			&dst.TemplateListArchive,
		},
		"config_value": {
			&src.TemplateConfigValue,
			&dst.TemplateConfigValue,
		},
		"dry_run_archive": {
			&src.TemplateDryRunArchive,
			&dst.TemplateDryRunArchive,
		},
		"error": {
			&src.TemplateError,
			&dst.TemplateError,
		},
	}

	for _, tmpl := range templates {
		if *tmpl.src != *tmpl.dst && *tmpl.src != "" {
			*tmpl.dst = *tmpl.src
		}
	}
}

// mergeFileTemplates merges file operation templates.
func mergeFileTemplates(dst, src *Config) {
	templates := map[string]struct {
		src *string
		dst *string
	}{
		"created_backup": {
			&src.TemplateCreatedBackup,
			&dst.TemplateCreatedBackup,
		},
		"identical_backup": {
			&src.TemplateIdenticalBackup,
			&dst.TemplateIdenticalBackup,
		},
		"list_backup": {
			&src.TemplateListBackup,
			&dst.TemplateListBackup,
		},
		"dry_run_backup": {
			&src.TemplateDryRunBackup,
			&dst.TemplateDryRunBackup,
		},
	}

	for _, tmpl := range templates {
		if *tmpl.src != *tmpl.dst && *tmpl.src != "" {
			*tmpl.dst = *tmpl.src
		}
	}
}

// mergePatterns merges regex pattern settings.
func mergePatterns(dst, src *Config) {
	if src.PatternArchiveFilename != DefaultConfig().PatternArchiveFilename {
		dst.PatternArchiveFilename = src.PatternArchiveFilename
	}
	if src.PatternBackupFilename != DefaultConfig().PatternBackupFilename {
		dst.PatternBackupFilename = src.PatternBackupFilename
	}
	if src.PatternConfigLine != DefaultConfig().PatternConfigLine {
		dst.PatternConfigLine = src.PatternConfigLine
	}
	if src.PatternTimestamp != DefaultConfig().PatternTimestamp {
		dst.PatternTimestamp = src.PatternTimestamp
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// GetConfigValues returns a slice of ConfigValue containing all configuration
// values and their sources.
func GetConfigValues(cfg *Config) []ConfigValue {
	// This would be used by the --config command to display all configuration values
	// For now, return basic values - this can be expanded
	return []ConfigValue{
		{Name: "archive_dir_path", Value: cfg.ArchiveDirPath, Source: "config"},
		{Name: "use_current_dir_name", Value: boolToString(cfg.UseCurrentDirName), Source: "config"},
		{Name: "include_git_info", Value: boolToString(cfg.IncludeGitInfo), Source: "config"},
		{Name: "backup_dir_path", Value: cfg.BackupDirPath, Source: "config"},
		{Name: "use_current_dir_name_for_files", Value: boolToString(cfg.UseCurrentDirNameForFiles), Source: "config"},
		{Name: "verify_on_create", Value: boolToString(cfg.Verification.VerifyOnCreate), Source: "config"},
		{Name: "checksum_algorithm", Value: cfg.Verification.ChecksumAlgorithm, Source: "config"},
	}
}

// GetConfigValuesWithSources returns a slice of ConfigValue containing all configuration
// values with their actual sources (default, config file, etc.).
// The returned values are sorted alphabetically by configuration name.
func GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue {
	defaultCfg := DefaultConfig()
	configSource := determineConfigSource(root)
	getSource := createSourceDeterminer(configSource)

	var configValues []ConfigValue
	configValues = append(configValues, getBasicConfigValues(cfg, defaultCfg, getSource)...)
	configValues = append(configValues, getStatusCodeValues(cfg, defaultCfg, getSource)...)
	configValues = append(configValues, getVerificationValues(cfg, defaultCfg, getSource)...)

	// Requirement: Sort configuration values alphabetically by name
	sort.Slice(configValues, func(i, j int) bool {
		return configValues[i].Name < configValues[j].Name
	})

	return configValues
}

// determineConfigSource finds which config file was actually loaded
func determineConfigSource(root string) string {
	searchPaths := getConfigSearchPaths()
	for _, configPath := range searchPaths {
		expandedPath := expandPath(configPath)
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
		}
		if _, err := os.Stat(expandedPath); err == nil {
			return expandedPath
		}
	}
	return "default"
}

// createSourceDeterminer creates a function to determine if a value is default or from config
func createSourceDeterminer(configSource string) func(interface{}, interface{}) string {
	return func(current, defaultVal interface{}) string {
		switch v := current.(type) {
		case string:
			if v == defaultVal.(string) {
				return "default"
			}
		case bool:
			if v == defaultVal.(bool) {
				return "default"
			}
		case int:
			if v == defaultVal.(int) {
				return "default"
			}
		case []string:
			if equalStringSlices(v, defaultVal.([]string)) {
				return "default"
			}
		}
		return configSource
	}
}

// getBasicConfigValues returns basic configuration values
func getBasicConfigValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return []ConfigValue{
		{
			Name:   "archive_dir_path",
			Value:  cfg.ArchiveDirPath,
			Source: getSource(cfg.ArchiveDirPath, defaultCfg.ArchiveDirPath),
		},
		{
			Name:   "backup_dir_path",
			Value:  cfg.BackupDirPath,
			Source: getSource(cfg.BackupDirPath, defaultCfg.BackupDirPath),
		},
		{
			Name:   "use_current_dir_name",
			Value:  boolToString(cfg.UseCurrentDirName),
			Source: getSource(cfg.UseCurrentDirName, defaultCfg.UseCurrentDirName),
		},
		{
			Name:   "use_current_dir_name_for_files",
			Value:  boolToString(cfg.UseCurrentDirNameForFiles),
			Source: getSource(cfg.UseCurrentDirNameForFiles, defaultCfg.UseCurrentDirNameForFiles),
		},
		{
			Name:   "include_git_info",
			Value:  boolToString(cfg.IncludeGitInfo),
			Source: getSource(cfg.IncludeGitInfo, defaultCfg.IncludeGitInfo),
		},
	}
}

// getStatusCodeValues returns status code configuration values
func getStatusCodeValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return []ConfigValue{
		{
			Name:   "status_config_error",
			Value:  fmt.Sprintf("%d", cfg.StatusConfigError),
			Source: getSource(cfg.StatusConfigError, defaultCfg.StatusConfigError),
		},
		{
			Name:   "status_created_archive",
			Value:  fmt.Sprintf("%d", cfg.StatusCreatedArchive),
			Source: getSource(cfg.StatusCreatedArchive, defaultCfg.StatusCreatedArchive),
		},
		{
			Name:   "status_created_backup",
			Value:  fmt.Sprintf("%d", cfg.StatusCreatedBackup),
			Source: getSource(cfg.StatusCreatedBackup, defaultCfg.StatusCreatedBackup),
		},
		{
			Name:   "status_disk_full",
			Value:  fmt.Sprintf("%d", cfg.StatusDiskFull),
			Source: getSource(cfg.StatusDiskFull, defaultCfg.StatusDiskFull),
		},
		{
			Name:   "status_permission_denied",
			Value:  fmt.Sprintf("%d", cfg.StatusPermissionDenied),
			Source: getSource(cfg.StatusPermissionDenied, defaultCfg.StatusPermissionDenied),
		},
	}
}

// getVerificationValues returns verification configuration values
func getVerificationValues(cfg, defaultCfg *Config, getSource func(interface{}, interface{}) string) []ConfigValue {
	return []ConfigValue{
		{
			Name:   "verify_on_create",
			Value:  boolToString(cfg.Verification.VerifyOnCreate),
			Source: getSource(cfg.Verification.VerifyOnCreate, defaultCfg.Verification.VerifyOnCreate),
		},
		{
			Name:   "checksum_algorithm",
			Value:  cfg.Verification.ChecksumAlgorithm,
			Source: getSource(cfg.Verification.ChecksumAlgorithm, defaultCfg.Verification.ChecksumAlgorithm),
		},
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
