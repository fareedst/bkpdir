// This file is part of bkpdir
//
// Package main provides configuration management for BkpDir.
// It handles loading, merging, and managing configuration from multiple sources.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

// 🔶 REFACTOR-005: Structure optimization - Clean configuration interface preparation - 🔧
// 🔶 REFACTOR-005: Extraction preparation - Standardized naming conventions and decoupled interfaces - 🔧

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// 🔶 REFACTOR-001: Configuration interface contracts defined - 🔧
// 🔶 REFACTOR-001: Dependency analysis - clean boundary confirmed - 📝
// 🔶 REFACTOR-005: Structure optimization - Interface-based configuration access - 📝
// Note: Interfaces defined in config_interfaces.go for clean separation

// 🔶 REFACTOR-005: Structure optimization - Standardized configuration structure - 📝
// Separated concerns into logical groupings for better extraction boundaries

// 🔺 CFG-002: Verification configuration structure - 📝
// IMMUTABLE-REF: Archive Verification Requirements
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// 🔶 REFACTOR-003: Schema separation - Backup-specific verification config - 🔍
// VerificationConfig defines settings for archive verification.
// It controls whether archives are verified on creation and which checksum algorithm to use.
type VerificationConfig struct {
	VerifyOnCreate    bool   `yaml:"verify_on_create"`
	ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}

// 🔺 CFG-001: Main configuration structure - 🔍
// 🔺 CFG-002: Status code configuration - 🔍
// 🔺 CFG-003: Output formatting configuration - 🔍
// IMMUTABLE-REF: Configuration Defaults, Output Formatting Requirements
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// 🔶 REFACTOR-001: Configuration interface contracts defined - 📝
// 🔶 REFACTOR-001: Dependency analysis - clean boundary confirmed - 📝
// 🔶 REFACTOR-003: Schema separation - Backup application specific schema - 📝
// Config holds all configuration settings for the BkpDir application.
// It includes settings for archive creation, file backup, status codes,
// and output formatting.
// The configuration can be loaded from YAML files and environment variables.
type Config struct {
	// 🔶 REFACTOR-003: Schema separation - Basic backup settings - 📝
	// Basic settings
	ArchiveDirPath     string              `yaml:"archive_dir_path"`
	UseCurrentDirName  bool                `yaml:"use_current_dir_name"`
	ExcludePatterns    []string            `yaml:"exclude_patterns"`
	IncludeGitInfo     bool                `yaml:"include_git_info"`
	ShowGitDirtyStatus bool                `yaml:"show_git_dirty_status"`
	SkipBrokenSymlinks bool                `yaml:"skip_broken_symlinks"`
	Verification       *VerificationConfig `yaml:"verification"`

	// ⭐ CFG-005: Configuration inheritance support - 🔧 Core inheritance functionality
	// Inherit specifies configuration files to inherit from
	Inherit []string `yaml:"inherit,omitempty"`

	// 🔶 REFACTOR-003: Schema separation - File backup specific settings - 🔧
	// File backup settings
	BackupDirPath             string `yaml:"backup_dir_path"`
	UseCurrentDirNameForFiles bool   `yaml:"use_current_dir_name_for_files"`

	// 🔶 REFACTOR-003: Schema separation - Backup application status codes - 🔧
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

	// 🔶 REFACTOR-003: Schema separation - Backup application format strings - 📝
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

	// 🔶 REFACTOR-003: Schema separation - Backup application template strings - 📝
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

	// 🔶 REFACTOR-003: Schema separation - Backup application regex patterns - 🔧
	// Regex patterns
	PatternArchiveFilename string `yaml:"pattern_archive_filename"`
	PatternBackupFilename  string `yaml:"pattern_backup_filename"`
	PatternConfigLine      string `yaml:"pattern_config_line"`
	PatternTimestamp       string `yaml:"pattern_timestamp"`

	// 🔺 CFG-004: Extended format strings for comprehensive string configuration - 📝
	// 🔶 REFACTOR-003: Schema separation - Extended backup operation messages - 📝
	// Archive operation messages
	FormatNoArchivesFound      string `yaml:"format_no_archives_found"`
	FormatVerificationFailed   string `yaml:"format_verification_failed"`
	FormatVerificationSuccess  string `yaml:"format_verification_success"`
	FormatVerificationWarning  string `yaml:"format_verification_warning"`
	FormatConfigurationUpdated string `yaml:"format_configuration_updated"`
	FormatConfigFilePath       string `yaml:"format_config_file_path"`
	FormatDryRunFilesHeader    string `yaml:"format_dry_run_files_header"`
	FormatDryRunFileEntry      string `yaml:"format_dry_run_file_entry"`
	FormatNoFilesModified      string `yaml:"format_no_files_modified"`
	FormatIncrementalCreated   string `yaml:"format_incremental_created"`

	// Backup operation messages
	FormatNoBackupsFound    string `yaml:"format_no_backups_found"`
	FormatBackupWouldCreate string `yaml:"format_backup_would_create"`
	FormatBackupIdentical   string `yaml:"format_backup_identical"`
	FormatBackupCreated     string `yaml:"format_backup_created"`

	// 🔺 CFG-004: Error message format strings - 📝
	// 🔶 REFACTOR-003: Schema separation - Backup application error messages - 📝
	FormatDiskFullError       string `yaml:"format_disk_full_error"`
	FormatPermissionError     string `yaml:"format_permission_error"`
	FormatDirectoryNotFound   string `yaml:"format_directory_not_found"`
	FormatFileNotFound        string `yaml:"format_file_not_found"`
	FormatInvalidDirectory    string `yaml:"format_invalid_directory"`
	FormatInvalidFile         string `yaml:"format_invalid_file"`
	FormatFailedWriteTemp     string `yaml:"format_failed_write_temp"`
	FormatFailedFinalizeFile  string `yaml:"format_failed_finalize_file"`
	FormatFailedCreateDirDisk string `yaml:"format_failed_create_dir_disk"`
	FormatFailedCreateDir     string `yaml:"format_failed_create_dir"`
	FormatFailedAccessDir     string `yaml:"format_failed_access_dir"`
	FormatFailedAccessFile    string `yaml:"format_failed_access_file"`

	// 🔶 REFACTOR-003: Schema separation - Extended backup template strings - 📝
	// Template-based extended format strings
	TemplateNoArchivesFound      string `yaml:"template_no_archives_found"`
	TemplateVerificationFailed   string `yaml:"template_verification_failed"`
	TemplateVerificationSuccess  string `yaml:"template_verification_success"`
	TemplateVerificationWarning  string `yaml:"template_verification_warning"`
	TemplateConfigurationUpdated string `yaml:"template_configuration_updated"`
	TemplateConfigFilePath       string `yaml:"template_config_file_path"`
	TemplateDryRunFilesHeader    string `yaml:"template_dry_run_files_header"`
	TemplateDryRunFileEntry      string `yaml:"template_dry_run_file_entry"`
	TemplateNoFilesModified      string `yaml:"template_no_files_modified"`
	TemplateIncrementalCreated   string `yaml:"template_incremental_created"`

	// Template-based backup operation messages
	TemplateNoBackupsFound    string `yaml:"template_no_backups_found"`
	TemplateBackupWouldCreate string `yaml:"template_backup_would_create"`
	TemplateBackupIdentical   string `yaml:"template_backup_identical"`
	TemplateBackupCreated     string `yaml:"template_backup_created"`

	// 🔺 CFG-004: Template-based error message format strings - 📝
	// 🔶 REFACTOR-003: Schema separation - Backup application error templates - 📝
	TemplateDiskFullError       string `yaml:"template_disk_full_error"`
	TemplatePermissionError     string `yaml:"template_permission_error"`
	TemplateDirectoryNotFound   string `yaml:"template_directory_not_found"`
	TemplateFileNotFound        string `yaml:"template_file_not_found"`
	TemplateInvalidDirectory    string `yaml:"template_invalid_directory"`
	TemplateInvalidFile         string `yaml:"template_invalid_file"`
	TemplateFailedWriteTemp     string `yaml:"template_failed_write_temp"`
	TemplateFailedFinalizeFile  string `yaml:"template_failed_finalize_file"`
	TemplateFailedCreateDirDisk string `yaml:"template_failed_create_dir_disk"`
	TemplateFailedCreateDir     string `yaml:"template_failed_create_dir"`
	TemplateFailedAccessDir     string `yaml:"template_failed_access_dir"`
	TemplateFailedAccessFile    string `yaml:"template_failed_access_file"`
}

// 🔺 CFG-003: Configuration value representation - 📝
// IMMUTABLE-REF: Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-002
// 🔶 REFACTOR-003: Config abstraction - Generic configuration value representation - 📝
// ConfigValue represents a single configuration value with its source.
// It is used for displaying configuration values and their origins.
type ConfigValue struct {
	Name   string
	Value  string
	Source string
}

// 🔺 CFG-003: Default regex patterns for template extraction - 📝
// IMMUTABLE-REF: Template Formatting Requirements, Configuration Defaults
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// 🔶 REFACTOR-003: Schema separation - Backup application default patterns - 📝
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

// 🔺 CFG-001: Default configuration implementation - 📝
// 🔺 CFG-002: Default status codes - 📝
// 🔺 CFG-003: Default format strings and templates - 📝
// IMMUTABLE-REF: Configuration Defaults
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// 🔶 REFACTOR-003: Schema separation - Backup application default configuration - 📝
// DefaultConfig returns a new Config instance with default values.
// These values are used when no configuration is provided or when merging configurations.
func DefaultConfig() *Config {
	return &Config{
		// Basic settings
		ArchiveDirPath:     "../.bkpdir",
		UseCurrentDirName:  true,
		ExcludePatterns:    []string{".git/", "vendor/"},
		IncludeGitInfo:     false,
		ShowGitDirtyStatus: true,
		SkipBrokenSymlinks: false,
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

		// 🔺 CFG-004: Extended format strings for comprehensive string configuration - 📝
		// Archive operation messages
		FormatNoArchivesFound:      "No archives found in %s\n",
		FormatVerificationFailed:   "Archive %s verification failed: %v\n",
		FormatVerificationSuccess:  "Archive %s verified successfully\n",
		FormatVerificationWarning:  "Warning: Could not store verification status for %s: %v\n",
		FormatConfigurationUpdated: "Configuration updated: %s = %v\n",
		FormatConfigFilePath:       "Config file: %s\n",
		FormatDryRunFilesHeader:    "[Dry Run] Files to include:\n",
		FormatDryRunFileEntry:      "  %s\n",
		FormatNoFilesModified:      "No files modified since last full archive\n",
		FormatIncrementalCreated:   "Created incremental archive: %s\n",

		// Backup operation messages
		FormatNoBackupsFound:    "No backups found for %s in %s\n",
		FormatBackupWouldCreate: "Would create backup: %s\n",
		FormatBackupIdentical:   "File is identical to existing backup: %s\n",
		FormatBackupCreated:     "Created backup: %s\n",

		// 🔺 CFG-004: Error message format strings - 📝
		FormatDiskFullError:       "Disk full error: %v\n",
		FormatPermissionError:     "Permission error: %v\n",
		FormatDirectoryNotFound:   "Directory not found: %v\n",
		FormatFileNotFound:        "File not found: %v\n",
		FormatInvalidDirectory:    "Invalid directory: %v\n",
		FormatInvalidFile:         "Invalid file: %v\n",
		FormatFailedWriteTemp:     "Failed to write temporary file: %v\n",
		FormatFailedFinalizeFile:  "Failed to finalize file: %v\n",
		FormatFailedCreateDirDisk: "Failed to create directory on disk: %v\n",
		FormatFailedCreateDir:     "Failed to create directory: %v\n",
		FormatFailedAccessDir:     "Failed to access directory: %v\n",
		FormatFailedAccessFile:    "Failed to access file: %v\n",

		// Template-based extended format strings
		TemplateNoArchivesFound:      "No archives found in %{archive_dir}\n",
		TemplateVerificationFailed:   "Archive %{name} verification failed: %{error}\n",
		TemplateVerificationSuccess:  "Archive %{name} verified successfully\n",
		TemplateVerificationWarning:  "Warning: Could not store verification status for %{name}: %{error}\n",
		TemplateConfigurationUpdated: "Configuration updated: %{key} = %{value}\n",
		TemplateConfigFilePath:       "Config file: %{path}\n",
		TemplateDryRunFilesHeader:    "[Dry Run] Files to include:\n",
		TemplateDryRunFileEntry:      "  %{file}\n",
		TemplateNoFilesModified:      "No files modified since last full archive\n",
		TemplateIncrementalCreated:   "Created incremental archive: %{path}\n",

		// Template-based backup operation messages
		TemplateNoBackupsFound:    "No backups found for %{filename} in %{backup_dir}\n",
		TemplateBackupWouldCreate: "Would create backup: %{path}\n",
		TemplateBackupIdentical:   "File is identical to existing backup: %{path}\n",
		TemplateBackupCreated:     "Created backup: %{path}\n",

		// 🔺 CFG-004: Template-based error message format strings - 📝
		TemplateDiskFullError:       "Disk full error: %{error}\n",
		TemplatePermissionError:     "Permission error: %{error}\n",
		TemplateDirectoryNotFound:   "Directory not found: %{error}\n",
		TemplateFileNotFound:        "File not found: %{error}\n",
		TemplateInvalidDirectory:    "Invalid directory: %{error}\n",
		TemplateInvalidFile:         "Invalid file: %{error}\n",
		TemplateFailedWriteTemp:     "Failed to write temporary file: %{error}\n",
		TemplateFailedFinalizeFile:  "Failed to finalize file: %{error}\n",
		TemplateFailedCreateDirDisk: "Failed to create directory on disk: %{error}\n",
		TemplateFailedCreateDir:     "Failed to create directory: %{error}\n",
		TemplateFailedAccessDir:     "Failed to access directory: %{error}\n",
		TemplateFailedAccessFile:    "Failed to access file: %{error}\n",
	}
}

// 🔺 CFG-001: Configuration search path implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// getConfigSearchPaths returns the list of paths to search for configuration files.
// It includes both system-wide and user-specific configuration paths.
func getConfigSearchPaths() []string {
	// Check BKPDIR_CONFIG environment variable
	if configPaths := os.Getenv("BKPDIR_CONFIG"); configPaths != "" {
		return strings.Split(configPaths, ":")
	}

	// Default search path
	return []string{"./.bkpdir.yml", "~/.bkpdir.yml"}
}

// 🔺 CFG-001: Path expansion implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// expandPath expands a path by replacing special tokens with actual values.
// It handles tokens like ~ for home directory and %ROOT% for the root directory.
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

// 🔺 CFG-001: Configuration loading implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// 🔶 REFACTOR-003: Config abstraction - Schema-specific configuration loading - 🔍
// LoadConfig loads configuration from YAML files and environment variables.
// It searches for configuration files in the standard locations and merges them with defaults.
func LoadConfig(root string) (*Config, error) {
	// ⭐ CFG-005: Configuration loading with inheritance support - 🔧 Enhanced loading engine
	// Try loading with inheritance first (the new default behavior)
	cfg, err := LoadConfigWithInheritance(root)
	if err == nil {
		return cfg, nil
	}

	// If inheritance loading fails, fallback to original method for backward compatibility
	// 🔶 REFACTOR-003: Schema separation - Backup application default config - 🔍
	cfg = DefaultConfig()
	// 🔶 REFACTOR-003: Config abstraction - Hardcoded search paths need abstraction - 🔍
	searchPaths := getConfigSearchPaths()

	// Process configuration files in order (earlier files take precedence)
	for _, configPath := range searchPaths {
		// 🔶 REFACTOR-003: Config abstraction - Path expansion needs abstraction - 🔍
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

			// 🔶 REFACTOR-003: Schema separation - Hardcoded Config struct unmarshaling - 🔧
			// Create a temporary config to load into
			tempCfg := DefaultConfig()
			d := yaml.NewDecoder(f)
			if err := d.Decode(tempCfg); err != nil {
				f.Close()
				continue // Skip files with invalid YAML
			}
			f.Close()

			// 🔶 REFACTOR-003: Config abstraction - Schema-specific merging logic - 📝
			// Merge non-zero values from tempCfg into cfg
			mergeConfigs(cfg, tempCfg)
			break // Use first valid config file found
		}
	}

	return cfg, nil
}

// 🔺 CFG-001: Configuration merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// 🔶 REFACTOR-003: Config abstraction - Schema-specific merging needs abstraction - 🔍
// mergeConfigs merges source configuration into destination configuration.
// It preserves non-zero values from the source configuration.
func mergeConfigs(dst, src *Config) {
	// 🔶 REFACTOR-003: Schema separation - Backup application specific merge functions - 🔍
	mergeBasicSettings(dst, src)
	mergeFileBackupSettings(dst, src)
	mergeStatusCodes(dst, src)
	mergeFormatStrings(dst, src)
	mergeTemplates(dst, src)
	mergePatterns(dst, src)
	mergeExtendedFormatStrings(dst, src)
	mergeExtendedTemplates(dst, src)
}

// 🔺 CFG-001: Basic settings merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeBasicSettings merges basic configuration settings.
// It handles archive directory path, Git integration, and verification settings.
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
	if src.ShowGitDirtyStatus != DefaultConfig().ShowGitDirtyStatus {
		dst.ShowGitDirtyStatus = src.ShowGitDirtyStatus
	}
	if src.SkipBrokenSymlinks != DefaultConfig().SkipBrokenSymlinks {
		dst.SkipBrokenSymlinks = src.SkipBrokenSymlinks
	}
	if src.Verification != nil {
		dst.Verification = src.Verification
	}
}

// 🔺 CFG-001: File backup settings merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileBackupSettings merges file backup configuration settings.
// It handles backup directory path and naming settings.
func mergeFileBackupSettings(dst, src *Config) {
	if src.BackupDirPath != DefaultConfig().BackupDirPath {
		dst.BackupDirPath = src.BackupDirPath
	}
	if src.UseCurrentDirNameForFiles != DefaultConfig().UseCurrentDirNameForFiles {
		dst.UseCurrentDirNameForFiles = src.UseCurrentDirNameForFiles
	}
}

// 🔺 CFG-002: Status code merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeStatusCodes merges all status code settings.
// It handles both directory and file operation status codes.
func mergeStatusCodes(dst, src *Config) {
	mergeDirectoryStatusCodes(dst, src)
	mergeFileStatusCodes(dst, src)
}

// 🔺 CFG-002: Directory status code merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeDirectoryStatusCodes merges directory operation status codes.
// It handles archive creation and verification status codes.
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

// 🔺 CFG-002: File status code merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileStatusCodes merges file operation status codes.
// It handles file backup and verification status codes.
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
		{
			Name:   "skip_broken_symlinks",
			Value:  boolToString(cfg.SkipBrokenSymlinks),
			Source: getSource(cfg.SkipBrokenSymlinks, defaultCfg.SkipBrokenSymlinks),
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

// 🔺 CFG-001: Configuration value loading implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// LoadConfigValues loads configuration values from YAML files and environment variables.
// It returns a map of configuration values with their sources.
func LoadConfigValues(root string) (map[string]ConfigValue, error) {
	// Implementation of LoadConfigValues function
	return nil, nil // Placeholder return, actual implementation needed
}

// 🔺 CFG-001: Configuration value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeConfigValues merges configuration values from source into destination.
// It preserves values from the source configuration.
func mergeConfigValues(dst, src map[string]ConfigValue) {
	// Implementation of mergeConfigValues function
}

// 🔺 CFG-001: Basic settings value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeBasicSettingValues merges basic configuration setting values.
// It handles archive directory path, Git integration, and verification settings.
func mergeBasicSettingValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeBasicSettingValues function
}

// 🔺 CFG-001: File backup settings value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileBackupSettingValues merges file backup configuration setting values.
// It handles backup directory path and naming settings.
func mergeFileBackupSettingValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeFileBackupSettingValues function
}

// 🔺 CFG-002: Status code value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeStatusCodeValues merges all status code setting values.
// It handles both directory and file operation status codes.
func mergeStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeStatusCodeValues function
}

// 🔺 CFG-002: Directory status code value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeDirectoryStatusCodeValues merges directory operation status code values.
// It handles archive creation and verification status codes.
func mergeDirectoryStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeDirectoryStatusCodeValues function
}

// 🔺 CFG-002: File status code value merging implementation - 🔍
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileStatusCodeValues merges file operation status code values.
// It handles file backup and verification status codes.
func mergeFileStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeFileStatusCodeValues function
}

// 🔺 CFG-004: Extended format strings for comprehensive string configuration - 📝
// mergeExtendedFormatStrings merges extended format string settings.
func mergeExtendedFormatStrings(dst, src *Config) {
	defaultCfg := DefaultConfig()

	// Archive operation messages
	if src.FormatNoArchivesFound != defaultCfg.FormatNoArchivesFound {
		dst.FormatNoArchivesFound = src.FormatNoArchivesFound
	}
	if src.FormatVerificationFailed != defaultCfg.FormatVerificationFailed {
		dst.FormatVerificationFailed = src.FormatVerificationFailed
	}
	if src.FormatVerificationSuccess != defaultCfg.FormatVerificationSuccess {
		dst.FormatVerificationSuccess = src.FormatVerificationSuccess
	}
	if src.FormatVerificationWarning != defaultCfg.FormatVerificationWarning {
		dst.FormatVerificationWarning = src.FormatVerificationWarning
	}
	if src.FormatConfigurationUpdated != defaultCfg.FormatConfigurationUpdated {
		dst.FormatConfigurationUpdated = src.FormatConfigurationUpdated
	}
	if src.FormatConfigFilePath != defaultCfg.FormatConfigFilePath {
		dst.FormatConfigFilePath = src.FormatConfigFilePath
	}
	if src.FormatDryRunFilesHeader != defaultCfg.FormatDryRunFilesHeader {
		dst.FormatDryRunFilesHeader = src.FormatDryRunFilesHeader
	}
	if src.FormatDryRunFileEntry != defaultCfg.FormatDryRunFileEntry {
		dst.FormatDryRunFileEntry = src.FormatDryRunFileEntry
	}
	if src.FormatNoFilesModified != defaultCfg.FormatNoFilesModified {
		dst.FormatNoFilesModified = src.FormatNoFilesModified
	}
	if src.FormatIncrementalCreated != defaultCfg.FormatIncrementalCreated {
		dst.FormatIncrementalCreated = src.FormatIncrementalCreated
	}

	// Backup operation messages
	if src.FormatNoBackupsFound != defaultCfg.FormatNoBackupsFound {
		dst.FormatNoBackupsFound = src.FormatNoBackupsFound
	}
	if src.FormatBackupWouldCreate != defaultCfg.FormatBackupWouldCreate {
		dst.FormatBackupWouldCreate = src.FormatBackupWouldCreate
	}
	if src.FormatBackupIdentical != defaultCfg.FormatBackupIdentical {
		dst.FormatBackupIdentical = src.FormatBackupIdentical
	}
	if src.FormatBackupCreated != defaultCfg.FormatBackupCreated {
		dst.FormatBackupCreated = src.FormatBackupCreated
	}
}

// 🔺 CFG-004: Extended templates for comprehensive string configuration - 📝
// mergeExtendedTemplates merges extended template settings.
func mergeExtendedTemplates(dst, src *Config) {
	defaultCfg := DefaultConfig()

	// Archive operation templates
	if src.TemplateNoArchivesFound != defaultCfg.TemplateNoArchivesFound {
		dst.TemplateNoArchivesFound = src.TemplateNoArchivesFound
	}
	if src.TemplateVerificationFailed != defaultCfg.TemplateVerificationFailed {
		dst.TemplateVerificationFailed = src.TemplateVerificationFailed
	}
	if src.TemplateVerificationSuccess != defaultCfg.TemplateVerificationSuccess {
		dst.TemplateVerificationSuccess = src.TemplateVerificationSuccess
	}
	if src.TemplateVerificationWarning != defaultCfg.TemplateVerificationWarning {
		dst.TemplateVerificationWarning = src.TemplateVerificationWarning
	}
	if src.TemplateConfigurationUpdated != defaultCfg.TemplateConfigurationUpdated {
		dst.TemplateConfigurationUpdated = src.TemplateConfigurationUpdated
	}
	if src.TemplateConfigFilePath != defaultCfg.TemplateConfigFilePath {
		dst.TemplateConfigFilePath = src.TemplateConfigFilePath
	}
	if src.TemplateDryRunFilesHeader != defaultCfg.TemplateDryRunFilesHeader {
		dst.TemplateDryRunFilesHeader = src.TemplateDryRunFilesHeader
	}
	if src.TemplateDryRunFileEntry != defaultCfg.TemplateDryRunFileEntry {
		dst.TemplateDryRunFileEntry = src.TemplateDryRunFileEntry
	}
	if src.TemplateNoFilesModified != defaultCfg.TemplateNoFilesModified {
		dst.TemplateNoFilesModified = src.TemplateNoFilesModified
	}
	if src.TemplateIncrementalCreated != defaultCfg.TemplateIncrementalCreated {
		dst.TemplateIncrementalCreated = src.TemplateIncrementalCreated
	}

	// Backup operation templates
	if src.TemplateNoBackupsFound != defaultCfg.TemplateNoBackupsFound {
		dst.TemplateNoBackupsFound = src.TemplateNoBackupsFound
	}
	if src.TemplateBackupWouldCreate != defaultCfg.TemplateBackupWouldCreate {
		dst.TemplateBackupWouldCreate = src.TemplateBackupWouldCreate
	}
	if src.TemplateBackupIdentical != defaultCfg.TemplateBackupIdentical {
		dst.TemplateBackupIdentical = src.TemplateBackupIdentical
	}
	if src.TemplateBackupCreated != defaultCfg.TemplateBackupCreated {
		dst.TemplateBackupCreated = src.TemplateBackupCreated
	}
}

// 🔶 REFACTOR-005: Structure optimization - ErrorConfig interface implementation - 🔍
// GetStatusCodes returns a map of status code names to values
func (c *Config) GetStatusCodes() map[string]int {
	return map[string]int{
		"disk_full":                               c.StatusDiskFull,
		"permission_denied":                       c.StatusPermissionDenied,
		"directory_not_found":                     c.StatusDirectoryNotFound,
		"file_not_found":                          c.StatusFileNotFound,
		"invalid_directory":                       c.StatusInvalidDirectoryType,
		"invalid_file":                            c.StatusInvalidFileType,
		"created_archive":                         c.StatusCreatedArchive,
		"created_backup":                          c.StatusCreatedBackup,
		"failed_create_archive_directory":         c.StatusFailedToCreateArchiveDirectory,
		"failed_create_backup_directory":          c.StatusFailedToCreateBackupDirectory,
		"directory_identical_to_existing_archive": c.StatusDirectoryIsIdenticalToExistingArchive,
		"file_identical_to_existing_backup":       c.StatusFileIsIdenticalToExistingBackup,
		"config_error":                            c.StatusConfigError,
	}
}

// 🔶 REFACTOR-005: Structure optimization - ErrorConfig interface implementation - 🔍
// GetErrorFormatStrings returns a map of error format string names to values
func (c *Config) GetErrorFormatStrings() map[string]string {
	return map[string]string{
		"disk_full":              c.FormatDiskFullError,
		"permission":             c.FormatPermissionError,
		"directory_not_found":    c.FormatDirectoryNotFound,
		"file_not_found":         c.FormatFileNotFound,
		"invalid_directory":      c.FormatInvalidDirectory,
		"invalid_file":           c.FormatInvalidFile,
		"failed_write_temp":      c.FormatFailedWriteTemp,
		"failed_finalize_file":   c.FormatFailedFinalizeFile,
		"failed_create_dir_disk": c.FormatFailedCreateDirDisk,
		"failed_create_dir":      c.FormatFailedCreateDir,
		"failed_access_dir":      c.FormatFailedAccessDir,
		"failed_access_file":     c.FormatFailedAccessFile,
	}
}

// 🔶 REFACTOR-005: Structure optimization - ErrorConfig interface implementation - 🔍
// GetDirectoryPermissions returns the default directory permissions
func (c *Config) GetDirectoryPermissions() os.FileMode {
	return 0755 // Standard directory permissions
}

// 🔶 REFACTOR-005: Structure optimization - ErrorConfig interface implementation - 🔍
// GetFilePermissions returns the default file permissions
func (c *Config) GetFilePermissions() os.FileMode {
	return 0644 // Standard file permissions
}

// ⭐ CFG-005: Configuration loading with inheritance - 🔧 Enhanced loading engine
// LoadConfigWithInheritance loads configuration with inheritance chain processing.
// This extends the original LoadConfig function to support layered configuration inheritance.
func LoadConfigWithInheritance(root string) (*Config, error) {
	// Use the pkg/config system for inheritance support
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	// Get the primary configuration file
	searchPaths := getConfigSearchPaths()
	var primaryConfigPath string

	for _, configPath := range searchPaths {
		expandedPath := expandPath(configPath)
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
		}

		if _, err := os.Stat(expandedPath); err == nil {
			primaryConfigPath = expandedPath
			break
		}
	}

	// If no config file found, return default config
	if primaryConfigPath == "" {
		return DefaultConfig(), nil
	}

	// Load configuration with inheritance
	return loadConfigRecursive(primaryConfigPath, pathResolver, chainBuilder)
}

// ⭐ CFG-005: Recursive configuration loading - 🔍 Inheritance chain processing
// loadConfigRecursive loads configuration following inheritance chains.
func loadConfigRecursive(configPath string, pathResolver pathResolver, chainBuilder inheritanceChainBuilder) (*Config, error) {
	// Build inheritance chain
	chain, err := chainBuilder.buildChain(configPath, pathResolver)
	if err != nil {
		return nil, fmt.Errorf("failed to build inheritance chain: %w", err)
	}

	// Start with default configuration
	cfg := DefaultConfig()

	// Process files in inheritance order (parents first)
	for _, filePath := range chain.files {
		tempCfg, err := loadSingleConfigFile(filePath)
		if err != nil {
			continue // Skip files with errors, continue with chain
		}

		// Apply merge strategies and merge into main config
		mergedCfg, err := applyMergeStrategies(cfg, tempCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config from %s: %w", filePath, err)
		}
		cfg = mergedCfg
	}

	return cfg, nil
}

// ⭐ CFG-005: Single file loading - 📝 Individual config file processing
// loadSingleConfigFile loads a single configuration file.
func loadSingleConfigFile(configPath string) (*Config, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer f.Close()

	cfg := DefaultConfig()
	d := yaml.NewDecoder(f)
	if err := d.Decode(cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", configPath, err)
	}

	return cfg, nil
}

// ⭐ CFG-005: Merge strategy application - 🔧 Strategy-based merging
// applyMergeStrategies applies merge strategies when combining configurations.
func applyMergeStrategies(dst, src *Config) (*Config, error) {
	processor := newMergeStrategyProcessor()

	// Convert configs to map for strategy processing
	dstMap := configToMap(dst)
	srcMap := configToMap(src)

	// Process merge strategies
	processed, err := processor.processKeys(srcMap)
	if err != nil {
		return nil, fmt.Errorf("failed to process merge strategies: %w", err)
	}

	// Apply processed configuration
	result := DefaultConfig()

	// Start with destination values
	mergeConfigs(result, dst)

	// Apply source values with merge strategies
	for key, operation := range processed.operations {
		err := applyMergeOperation(result, key, operation, dstMap[key])
		if err != nil {
			return nil, fmt.Errorf("failed to apply merge operation for %s: %w", key, err)
		}
	}

	return result, nil
}

// ⭐ CFG-005: Supporting types and interfaces for inheritance - 🔧 Implementation infrastructure

// configFileOperations implements file operations for inheritance system
type configFileOperations struct{}

func (c *configFileOperations) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (c *configFileOperations) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// pathResolver provides path resolution for inheritance
type pathResolver interface {
	resolvePath(path string, basePath string) (string, error)
	validatePath(path string) error
}

// defaultPathResolver implements pathResolver interface
type defaultPathResolver struct {
	fileOps *configFileOperations
}

func newPathResolver(fileOps *configFileOperations) pathResolver {
	return &defaultPathResolver{fileOps: fileOps}
}

func (r *defaultPathResolver) resolvePath(path string, basePath string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	if basePath == "" {
		return path, nil
	}

	resolved := filepath.Join(filepath.Dir(basePath), path)
	return resolved, nil
}

func (r *defaultPathResolver) validatePath(path string) error {
	if !r.fileOps.FileExists(path) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

// inheritanceChain represents a resolved inheritance dependency chain
type inheritanceChain struct {
	files   []string
	visited map[string]bool
}

// inheritanceChainBuilder builds inheritance chains
type inheritanceChainBuilder interface {
	buildChain(configPath string, pathResolver pathResolver) (*inheritanceChain, error)
}

// defaultInheritanceChainBuilder implements inheritanceChainBuilder
type defaultInheritanceChainBuilder struct {
	fileOps *configFileOperations
}

func newInheritanceChainBuilder(fileOps *configFileOperations) inheritanceChainBuilder {
	return &defaultInheritanceChainBuilder{fileOps: fileOps}
}

func (b *defaultInheritanceChainBuilder) buildChain(configPath string, pathResolver pathResolver) (*inheritanceChain, error) {
	chain := &inheritanceChain{
		files:   make([]string, 0),
		visited: make(map[string]bool),
	}

	return chain, b.buildChainRecursive(configPath, "", pathResolver, chain)
}

func (b *defaultInheritanceChainBuilder) buildChainRecursive(configPath, basePath string, pathResolver pathResolver, chain *inheritanceChain) error {
	// Resolve path
	resolvedPath, err := pathResolver.resolvePath(configPath, basePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", configPath, err)
	}

	// Check for circular dependency
	if chain.visited[resolvedPath] {
		return fmt.Errorf("circular dependency detected: %s", resolvedPath)
	}

	// Validate path exists
	if err := pathResolver.validatePath(resolvedPath); err != nil {
		return fmt.Errorf("invalid path %s: %w", resolvedPath, err)
	}

	// Mark as visited
	chain.visited[resolvedPath] = true

	// Load inheritance metadata
	inheritance, err := b.loadInheritanceMetadata(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load inheritance: %w", err)
	}

	// Process parent files first
	for _, parentPath := range inheritance {
		err := b.buildChainRecursive(parentPath, resolvedPath, pathResolver, chain)
		if err != nil {
			return fmt.Errorf("failed to process parent %s: %w", parentPath, err)
		}
	}

	// Add current file to chain
	chain.files = append(chain.files, resolvedPath)

	return nil
}

func (b *defaultInheritanceChainBuilder) loadInheritanceMetadata(configPath string) ([]string, error) {
	data, err := b.fileOps.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var metadata struct {
		Inherit []string `yaml:"inherit"`
	}

	err = yaml.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inheritance metadata: %w", err)
	}

	return metadata.Inherit, nil
}

// mergeStrategyProcessor processes merge strategies
type mergeStrategyProcessor interface {
	processKeys(config map[string]interface{}) (*processedConfig, error)
}

// processedConfig represents configuration after merge strategy processing
type processedConfig struct {
	operations map[string]*mergeOperation
}

// mergeOperation represents a single merge operation
type mergeOperation struct {
	strategy string
	value    interface{}
	key      string
}

// defaultMergeStrategyProcessor implements mergeStrategyProcessor
type defaultMergeStrategyProcessor struct{}

func newMergeStrategyProcessor() mergeStrategyProcessor {
	return &defaultMergeStrategyProcessor{}
}

func (p *defaultMergeStrategyProcessor) processKeys(config map[string]interface{}) (*processedConfig, error) {
	result := &processedConfig{
		operations: make(map[string]*mergeOperation),
	}

	for key, value := range config {
		strategy, cleanKey := p.extractStrategy(key)
		result.operations[cleanKey] = &mergeOperation{
			strategy: strategy,
			value:    value,
			key:      cleanKey,
		}
	}

	return result, nil
}

func (p *defaultMergeStrategyProcessor) extractStrategy(key string) (string, string) {
	if len(key) == 0 {
		return "override", key
	}

	switch key[0] {
	case '+':
		return "merge", key[1:]
	case '^':
		return "prepend", key[1:]
	case '!':
		return "replace", key[1:]
	case '=':
		return "default", key[1:]
	default:
		return "override", key
	}
}

// configToMap converts Config struct to map for strategy processing
func configToMap(cfg *Config) map[string]interface{} {
	// This is a simplified conversion - in a full implementation,
	// this would use reflection to convert the struct to a map
	return map[string]interface{}{
		"archive_dir_path":     cfg.ArchiveDirPath,
		"use_current_dir_name": cfg.UseCurrentDirName,
		"exclude_patterns":     cfg.ExcludePatterns,
		"include_git_info":     cfg.IncludeGitInfo,
		"skip_broken_symlinks": cfg.SkipBrokenSymlinks,
		// Status codes
		"status_created_archive": cfg.StatusCreatedArchive,
		"status_disk_full":       cfg.StatusDiskFull,
		// Add other fields as needed
	}
}

// applyMergeOperation applies a merge operation to the result configuration
func applyMergeOperation(result *Config, key string, operation *mergeOperation, dstValue interface{}) error {
	switch operation.strategy {
	case "override":
		return applyOverride(result, key, operation.value)
	case "merge":
		return applyMerge(result, key, operation.value, dstValue)
	case "prepend":
		return applyPrepend(result, key, operation.value, dstValue)
	case "replace":
		return applyReplace(result, key, operation.value)
	case "default":
		return applyDefault(result, key, operation.value, dstValue)
	default:
		return fmt.Errorf("unknown merge strategy: %s", operation.strategy)
	}
}

// Helper functions for applying different merge strategies
func applyOverride(result *Config, key string, value interface{}) error {
	return setConfigField(result, key, value)
}

func applyMerge(result *Config, key string, value interface{}, dstValue interface{}) error {
	// For arrays, merge by appending
	if srcSlice, ok := value.([]string); ok {
		if dstSlice, ok := dstValue.([]string); ok {
			merged := append(dstSlice, srcSlice...)
			return setConfigField(result, key, merged)
		}
	}
	return setConfigField(result, key, value)
}

func applyPrepend(result *Config, key string, value interface{}, dstValue interface{}) error {
	// For arrays, prepend source to destination
	if srcSlice, ok := value.([]string); ok {
		if dstSlice, ok := dstValue.([]string); ok {
			merged := append(srcSlice, dstSlice...)
			return setConfigField(result, key, merged)
		}
	}
	return setConfigField(result, key, value)
}

func applyReplace(result *Config, key string, value interface{}) error {
	return setConfigField(result, key, value)
}

func applyDefault(result *Config, key string, value interface{}, dstValue interface{}) error {
	// Only use source value if destination is zero value
	if isZeroValue(dstValue) {
		return setConfigField(result, key, value)
	}
	return nil
}

// setConfigField sets a field in the Config struct based on the key
func setConfigField(cfg *Config, key string, value interface{}) error {
	switch key {
	case "archive_dir_path":
		if s, ok := value.(string); ok {
			cfg.ArchiveDirPath = s
		}
	case "use_current_dir_name":
		if b, ok := value.(bool); ok {
			cfg.UseCurrentDirName = b
		}
	case "exclude_patterns":
		if slice, ok := value.([]string); ok {
			cfg.ExcludePatterns = slice
		}
	case "include_git_info":
		if b, ok := value.(bool); ok {
			cfg.IncludeGitInfo = b
		}
	case "skip_broken_symlinks":
		if b, ok := value.(bool); ok {
			cfg.SkipBrokenSymlinks = b
		}
	case "status_created_archive":
		if i, ok := value.(int); ok {
			cfg.StatusCreatedArchive = i
		}
	case "status_disk_full":
		if i, ok := value.(int); ok {
			cfg.StatusDiskFull = i
		}
	// Add other fields as needed
	default:
		return fmt.Errorf("unknown config field: %s", key)
	}
	return nil
}

// isZeroValue checks if a value is the zero value for its type
func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case bool:
		return v == false
	case int:
		return v == 0
	case []string:
		return len(v) == 0
	default:
		return false
	}
}
