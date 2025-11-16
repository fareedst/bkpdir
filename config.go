// This file is part of bkpdir
//
// Package main provides configuration management for BkpDir.
// It handles loading, merging, and managing configuration from multiple sources.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// CONFIG-DISCOVERY-001: Configuration discovery specification - Configuration discovery and loading [ACTION:discovery]
// Source: config.go - CONFIG-DISCOVERY-001
// Impact: Core functionality requirement for configuration discovery

// CONFIG-FILE-001: Configuration file specification - Configuration file handling [ACTION:core-functionality]
// Source: config.go - CONFIG-FILE-001
// Impact: Core functionality requirement for configuration file handling

// SERVICE-ARCH-001: Service architecture decision - Configuration service implementation [ACTION:core-functionality]
// Source: config.go - SERVICE-ARCH-001
// Impact: Configuration service implementation decision
package main

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]

import (
	"fmt"
	"hash"
	"hash/fnv"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]
// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]
// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
// Note: Interfaces defined in config_interfaces.go for clean separation

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
// Separated concerns into logical groupings for better extraction boundaries

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Archive Verification Requirements
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
// VerificationConfig defines settings for archive verification.
// It controls whether archives are verified on creation and which checksum algorithm to use.
type VerificationConfig struct {
	VerifyOnCreate    bool   `yaml:"verify_on_create"`
	ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Configuration Defaults, Output Formatting Requirements
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]
// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
// Config holds all configuration settings for the BkpDir application.
// It includes settings for archive creation, file backup, status codes,
// and output formatting.
// The configuration can be loaded from YAML files and environment variables.
type Config struct {
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	// Basic settings
	ArchiveDirPath     string              `yaml:"archive_dir_path"`
	UseCurrentDirName  bool                `yaml:"use_current_dir_name"`
	ExcludePatterns    []string            `yaml:"exclude_patterns"`
	IncludeGitInfo     bool                `yaml:"include_git_info"`      // Legacy - use Git.IncludeInfo
	ShowGitDirtyStatus bool                `yaml:"show_git_dirty_status"` // Legacy - use Git.ShowDirtyStatus
	SkipBrokenSymlinks bool                `yaml:"skip_broken_symlinks"`
	Verification       *VerificationConfig `yaml:"verification"`

	// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality] Core inheritance functionality
	// Inherit specifies configuration files to inherit from
	Inherit []string `yaml:"inherit,omitempty"`

	// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
	// Git configuration for repository detection and information extraction
	Git *GitConfig `yaml:"git,omitempty"`

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	// File backup settings
	BackupDirPath             string `yaml:"backup_dir_path"`
	UseCurrentDirNameForFiles bool   `yaml:"use_current_dir_name_for_files"`

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	// Regex patterns
	PatternArchiveFilename string `yaml:"pattern_archive_filename"`
	PatternBackupFilename  string `yaml:"pattern_backup_filename"`
	PatternConfigLine      string `yaml:"pattern_config_line"`
	PatternTimestamp       string `yaml:"pattern_timestamp"`

	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	// Enhanced format strings with stat information support
	FormatCreatedArchiveDetailed     string `yaml:"format_created_archive_detailed"`
	FormatIncrementalCreatedDetailed string `yaml:"format_incremental_created_detailed"`

	// Backup operation messages
	FormatNoBackupsFound    string `yaml:"format_no_backups_found"`
	FormatBackupWouldCreate string `yaml:"format_backup_would_create"`
	FormatBackupIdentical   string `yaml:"format_backup_identical"`
	FormatBackupCreated     string `yaml:"format_backup_created"`

	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	// Enhanced template strings with stat information support
	TemplateCreatedArchiveDetailed     string `yaml:"template_created_archive_detailed"`
	TemplateIncrementalCreatedDetailed string `yaml:"template_incremental_created_detailed"`

	// Template-based backup operation messages
	TemplateNoBackupsFound    string `yaml:"template_no_backups_found"`
	TemplateBackupWouldCreate string `yaml:"template_backup_would_create"`
	TemplateBackupIdentical   string `yaml:"template_backup_identical"`
	TemplateBackupCreated     string `yaml:"template_backup_created"`

	// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Commands - Display Configuration
// TEST-REF: TestDisplayConfig
// DECISION-REF: DEC-002
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
// ConfigValue represents a single configuration value with its source.
// It is used for displaying configuration values and their origins.
type ConfigValue struct {
	Name   string
	Value  string
	Source string
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Template Formatting Requirements, Configuration Defaults
// TEST-REF: TestTemplateFormatter
// DECISION-REF: DEC-003
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Configuration Defaults
// TEST-REF: TestDefaultConfig
// DECISION-REF: DEC-002
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
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

		// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
		Git: DefaultGitConfig(),

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

		// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
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

		// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
		// Enhanced format strings with stat information (backward compatible defaults)
		FormatCreatedArchiveDetailed:     "Created archive: %s (%s, %s)\n",
		FormatIncrementalCreatedDetailed: "Created incremental archive: %s (%s, %s)\n",

		// Backup operation messages
		FormatNoBackupsFound:    "No backups found for %s in %s\n",
		FormatBackupWouldCreate: "Would create backup: %s\n",
		FormatBackupIdentical:   "File is identical to existing backup: %s\n",
		FormatBackupCreated:     "Created backup: %s\n",

		// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
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

		// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
		// Enhanced template strings with stat information support
		TemplateCreatedArchiveDetailed:     "Created archive: %{path} (size: %{size_human}, modified: %{mtime})\n",
		TemplateIncrementalCreatedDetailed: "Created incremental archive: %{path} (size: %{size_human}, modified: %{mtime})\n",

		// Template-based backup operation messages
		TemplateNoBackupsFound:    "No backups found for %{filename} in %{backup_dir}\n",
		TemplateBackupWouldCreate: "Would create backup: %{path}\n",
		TemplateBackupIdentical:   "File is identical to existing backup: %{path}\n",
		TemplateBackupCreated:     "Created backup: %{path}\n",

		// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
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

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
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

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
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

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
// LoadConfig loads configuration from YAML files and environment variables.
// It searches for configuration files in the standard locations and merges them with defaults.
func LoadConfig(root string) (*Config, error) {
	if debug {
		fmt.Printf("DEBUG: Entered LoadConfig with root: %s\n", root)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]
	// Check if any config files exist first
	searchPaths := getConfigSearchPaths()
	foundAnyFile := false
	for _, configPath := range searchPaths {
		expandedPath := expandPath(configPath)
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
		}
		if _, err := os.Stat(expandedPath); err == nil {
			foundAnyFile = true
			break
		}
	}
	
	// If no files found, try inheritance loading (might have inheritance chains)
	// Otherwise, use fallback method directly for better compatibility
	if foundAnyFile {
		// Files exist - use fallback method directly for now
		// TODO: Fix LoadConfigWithInheritance to handle single files correctly
		if debug {
			fmt.Printf("DEBUG: Config files found, using fallback loading method\n")
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		fmt.Printf("DIAGNOSTIC: LoadConfig - Files found, using fallback method\n")
	} else {
		// No files found - try inheritance loading
		if debug {
			fmt.Printf("DEBUG: No config files found, trying inheritance loading\n")
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		fmt.Printf("DIAGNOSTIC: LoadConfig - No files found, trying inheritance loading\n")
		cfg, err := LoadConfigWithInheritance(root)
		if err == nil && cfg != nil {
			fmt.Printf("DIAGNOSTIC: LoadConfig - Inheritance loading succeeded\n")
			return cfg, nil
		}
		fmt.Printf("DIAGNOSTIC: LoadConfig - Inheritance loading failed or returned nil, using fallback\n")
	}

	// If inheritance loading fails, fallback to original method for backward compatibility
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	cfg := DefaultConfig()
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	// searchPaths already declared above
	if debug {
		fmt.Printf("DEBUG: searchPaths = %v\n", searchPaths)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT

	// Process configuration files in order (earlier files take precedence)
	fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Processing %d search paths\n", len(searchPaths))
	for i, configPath := range searchPaths {
		// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
		expandedPath := expandPath(configPath)
		if debug {
			fmt.Printf("DEBUG: configPath = %s, expandedPath = %s\n", configPath, expandedPath)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT

		// Make relative paths relative to root directory
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
			if debug {
				fmt.Printf("DEBUG: expandedPath made absolute: %s\n", expandedPath)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		}

		if debug {
			fmt.Printf("DEBUG: Checking existence of: %s\n", expandedPath)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		if _, err := os.Stat(expandedPath); err == nil {
			if debug {
				fmt.Printf("DEBUG: File exists: %s\n", expandedPath)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Processing file[%d]: %s\n", i, expandedPath)
			// Use loadSingleConfigFile to preserve merge strategy prefixes
			loadResult, err := loadSingleConfigFile(expandedPath)
			if err != nil {
				if debug {
					fmt.Printf("DEBUG: Failed to load config file %s: %v\n", expandedPath, err)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Failed to load file[%d] %s: %v\n", i, expandedPath, err)
				continue // Skip files with errors
			}
			tempCfg := loadResult.config

			// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
			// Merge non-zero values from tempCfg into cfg
			// Earlier files take precedence over later files as per specification
			// For fallback path: first file merges with defaults, subsequent files override
			if debug {
				fmt.Printf("DEBUG: Loading config from: %s\n", expandedPath)
				if len(tempCfg.ExcludePatterns) > 0 {
					fmt.Printf("DEBUG:   Exclude patterns in this file: %v\n", tempCfg.ExcludePatterns)
				}
				fmt.Printf("DEBUG:   Current exclude patterns before merge: %v\n", cfg.ExcludePatterns)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - File[%d] exclude_patterns: %v\n", i, tempCfg.ExcludePatterns)
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Current cfg exclude_patterns before merge: %v\n", cfg.ExcludePatterns)
			
			// Determine merge context: first file merges with defaults, subsequent files override
			isFirstFile := cfg.ArchiveDirPath == DefaultConfig().ArchiveDirPath && 
				len(cfg.ExcludePatterns) == len(DefaultConfig().ExcludePatterns)
			inheritContext := isFirstFile // First file merges with defaults
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - File[%d] isFirstFile=%v, inheritContext=%v\n", i, isFirstFile, inheritContext)
			
			mergedCfg, err := applyMergeStrategies(cfg, tempCfg, inheritContext, loadResult.rawMap)
			if err != nil {
				if debug {
					fmt.Printf("DEBUG: Failed to merge config from %s: %v\n", expandedPath, err)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Failed to merge file[%d] %s: %v\n", i, expandedPath, err)
				continue // Skip problematic merges
			}
			cfg = mergedCfg
			if debug {
				fmt.Printf("DEBUG:   Exclude patterns after merge: %v\n", cfg.ExcludePatterns)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - File[%d] exclude_patterns after merge: %v\n", i, cfg.ExcludePatterns)
		} else {
			fmt.Printf("DIAGNOSTIC: LoadConfig fallback - File[%d] %s does not exist\n", i, expandedPath)
		}
	}
	fmt.Printf("DIAGNOSTIC: LoadConfig fallback - Final cfg exclude_patterns: %v\n", cfg.ExcludePatterns)

	return cfg, nil
}

// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
// mergeConfigs merges source configuration into destination configuration.
// It preserves non-zero values from the source configuration.
func mergeConfigs(dst, src *Config) {
	// REFACTOR-003: See architecture.md - Configuration Abstraction [DECISION:format-processing]
	mergeBasicSettings(dst, src)
	mergeFileBackupSettings(dst, src)
	mergeStatusCodes(dst, src)
	mergeFormatStrings(dst, src)
	mergeTemplates(dst, src)
	mergePatterns(dst, src)
	mergeExtendedFormatStrings(dst, src)
	mergeExtendedTemplates(dst, src)
	// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
	mergeGitSettings(dst, src)
}

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
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
	// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
	// Merge exclude patterns by appending (earlier files take precedence, but patterns accumulate)
	if len(src.ExcludePatterns) > 0 && !equalStringSlices(src.ExcludePatterns, DefaultConfig().ExcludePatterns) {
		if debug {
			fmt.Printf("DEBUG: Merging exclude patterns - existing: %v, new: %v\n", dst.ExcludePatterns, src.ExcludePatterns)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		// Merge by appending new patterns that aren't already present
		existingMap := make(map[string]bool)
		for _, pattern := range dst.ExcludePatterns {
			existingMap[pattern] = true
		}
		for _, pattern := range src.ExcludePatterns {
			if !existingMap[pattern] {
				dst.ExcludePatterns = append(dst.ExcludePatterns, pattern)
				existingMap[pattern] = true
			}
		}
		if debug {
			fmt.Printf("DEBUG: Merged exclude patterns result: %v\n", dst.ExcludePatterns)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
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
	// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
	// Support legacy Git fields for backward compatibility
	if src.Git != nil {
		if dst.Git == nil {
			dst.Git = DefaultGitConfig()
		}
		// Merge explicit Git configuration over defaults
		if src.Git.IncludeInfo != dst.Git.IncludeInfo {
			dst.Git.IncludeInfo = src.Git.IncludeInfo
		}
		if src.Git.ShowDirtyStatus != dst.Git.ShowDirtyStatus {
			dst.Git.ShowDirtyStatus = src.Git.ShowDirtyStatus
		}
	}
}

// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
// mergeGitSettings merges Git configuration settings between configs.
// It handles both the new Git configuration and legacy fields for backward compatibility.
func mergeGitSettings(dst, src *Config) {
	defaultGit := DefaultGitConfig()

	// Initialize destination Git config if needed
	if dst.Git == nil {
		dst.Git = DefaultGitConfig()
	}

	// If source has Git configuration, merge it
	if src.Git != nil {
		mergeGitConfigStruct(dst.Git, src.Git, defaultGit)
	}

	// Handle legacy fields for backward compatibility
	// Legacy fields take precedence over Git struct for compatibility
	if src.IncludeGitInfo != defaultGit.IncludeInfo {
		dst.Git.IncludeInfo = src.IncludeGitInfo
		dst.IncludeGitInfo = src.IncludeGitInfo // Keep legacy field in sync
	}
	if src.ShowGitDirtyStatus != defaultGit.ShowDirtyStatus {
		dst.Git.ShowDirtyStatus = src.ShowGitDirtyStatus
		dst.ShowGitDirtyStatus = src.ShowGitDirtyStatus // Keep legacy field in sync
	}
}

// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
// mergeGitConfigStruct merges GitConfig struct fields
func mergeGitConfigStruct(dst, src, defaultCfg *GitConfig) {
	if src.Enabled != defaultCfg.Enabled {
		dst.Enabled = src.Enabled
	}
	if src.IncludeInfo != defaultCfg.IncludeInfo {
		dst.IncludeInfo = src.IncludeInfo
	}
	if src.ShowDirtyStatus != defaultCfg.ShowDirtyStatus {
		dst.ShowDirtyStatus = src.ShowDirtyStatus
	}
	if src.Command != defaultCfg.Command {
		dst.Command = src.Command
	}
	if src.WorkingDirectory != defaultCfg.WorkingDirectory {
		dst.WorkingDirectory = src.WorkingDirectory
	}
	if src.RequireCleanRepo != defaultCfg.RequireCleanRepo {
		dst.RequireCleanRepo = src.RequireCleanRepo
	}
	if src.AutoDetectRepo != defaultCfg.AutoDetectRepo {
		dst.AutoDetectRepo = src.AutoDetectRepo
	}
	if src.IncludeSubmodules != defaultCfg.IncludeSubmodules {
		dst.IncludeSubmodules = src.IncludeSubmodules
	}
	if src.IncludeBranch != defaultCfg.IncludeBranch {
		dst.IncludeBranch = src.IncludeBranch
	}
	if src.IncludeHash != defaultCfg.IncludeHash {
		dst.IncludeHash = src.IncludeHash
	}
	if src.IncludeStatus != defaultCfg.IncludeStatus {
		dst.IncludeStatus = src.IncludeStatus
	}
	if src.CommandTimeout != defaultCfg.CommandTimeout {
		dst.CommandTimeout = src.CommandTimeout
	}
	if src.MaxSubmoduleDepth != defaultCfg.MaxSubmoduleDepth {
		dst.MaxSubmoduleDepth = src.MaxSubmoduleDepth
	}
}

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
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

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeStatusCodes merges all status code settings.
// It handles both directory and file operation status codes.
func mergeStatusCodes(dst, src *Config) {
	mergeDirectoryStatusCodes(dst, src)
	mergeFileStatusCodes(dst, src)
}

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
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

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
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
// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Step 4.5: Preserve backward compatibility
// GetConfigValuesWithSources maintains backward compatibility while using the new reflection system.
// This function now uses automatic field discovery instead of manual enumeration.
func GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue {
	// Use the new reflection-based system and convert to legacy format
	enhancedValues := GetAllConfigValuesWithSources(cfg, root)

	var legacyValues []ConfigValue
	for _, enhanced := range enhancedValues {
		legacyValues = append(legacyValues, enhanced.ConfigValue)
	}

	return legacyValues
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
			// For arrays, check if current contains all defaults AND has additional items
			defaultSlice := defaultVal.([]string)
			if equalStringSlices(v, defaultSlice) {
				return "default"
			}
			// If current array contains all default values plus more, it's from config
			// Check if all defaults are present in current
			defaultMap := make(map[string]bool)
			for _, d := range defaultSlice {
				defaultMap[d] = true
			}
			hasAllDefaults := true
			for _, d := range defaultSlice {
				found := false
				for _, c := range v {
					if c == d {
						found = true
						break
					}
				}
				if !found {
					hasAllDefaults = false
					break
				}
			}
			// If it has all defaults but is longer, it's merged from config
			if hasAllDefaults && len(v) > len(defaultSlice) {
				return configSource
			}
			// If arrays are different, it's from config
			return configSource
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

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// LoadConfigValues loads configuration values from YAML files and environment variables.
// It returns a map of configuration values with their sources.
func LoadConfigValues(root string) (map[string]ConfigValue, error) {
	// Implementation of LoadConfigValues function
	return nil, nil // Placeholder return, actual implementation needed
}

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeConfigValues merges configuration values from source into destination.
// It preserves values from the source configuration.
func mergeConfigValues(dst, src map[string]ConfigValue) {
	// Implementation of mergeConfigValues function
}

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeBasicSettingValues merges basic configuration setting values.
// It handles archive directory path, Git integration, and verification settings.
func mergeBasicSettingValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeBasicSettingValues function
}

// CFG-001: See specification.md - Configuration Discovery [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileBackupSettingValues merges file backup configuration setting values.
// It handles backup directory path and naming settings.
func mergeFileBackupSettingValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeFileBackupSettingValues function
}

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeStatusCodeValues merges all status code setting values.
// It handles both directory and file operation status codes.
func mergeStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeStatusCodeValues function
}

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeDirectoryStatusCodeValues merges directory operation status code values.
// It handles archive creation and verification status codes.
func mergeDirectoryStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeDirectoryStatusCodeValues function
}

// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
// IMMUTABLE-REF: Configuration Discovery
// TEST-REF: TestGetConfigSearchPath
// DECISION-REF: DEC-002
// mergeFileStatusCodeValues merges file operation status code values.
// It handles file backup and verification status codes.
func mergeFileStatusCodeValues(dst, src map[string]ConfigValue, srcCfg *Config) {
	// Implementation of mergeFileStatusCodeValues function
}

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
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

	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	// Merge enhanced format strings
	if src.FormatCreatedArchiveDetailed != defaultCfg.FormatCreatedArchiveDetailed {
		dst.FormatCreatedArchiveDetailed = src.FormatCreatedArchiveDetailed
	}
	if src.FormatIncrementalCreatedDetailed != defaultCfg.FormatIncrementalCreatedDetailed {
		dst.FormatIncrementalCreatedDetailed = src.FormatIncrementalCreatedDetailed
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

// CFG-004: See specification.md - Configuration Format [DECISION:format-processing]
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

	// OUT-002: See specification.md - Output Formatting [DECISION:format-processing]
	// Merge enhanced template strings
	if src.TemplateCreatedArchiveDetailed != defaultCfg.TemplateCreatedArchiveDetailed {
		dst.TemplateCreatedArchiveDetailed = src.TemplateCreatedArchiveDetailed
	}
	if src.TemplateIncrementalCreatedDetailed != defaultCfg.TemplateIncrementalCreatedDetailed {
		dst.TemplateIncrementalCreatedDetailed = src.TemplateIncrementalCreatedDetailed
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

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
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

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
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

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
// GetDirectoryPermissions returns the default directory permissions
func (c *Config) GetDirectoryPermissions() os.FileMode {
	return 0755 // Standard directory permissions
}

// REFACTOR-005: See architecture.md - Structure Optimization [DECISION:maintenance]
// GetFilePermissions returns the default file permissions
func (c *Config) GetFilePermissions() os.FileMode {
	return 0644 // Standard file permissions
}

// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]
// LoadConfigWithInheritance loads configuration with inheritance chain processing.
// This extends the original LoadConfig function to support layered configuration inheritance.
func LoadConfigWithInheritance(root string) (*Config, error) {
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	searchPaths := getConfigSearchPaths()
	finalCfg := DefaultConfig()
	var foundAny bool
	var isFirstFile = true

	for _, configPath := range searchPaths {
		expandedPath := expandPath(configPath)
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(root, expandedPath)
		}
		if _, err := os.Stat(expandedPath); err == nil {
			if debug {
				fmt.Printf("DEBUG: Processing config file with inheritance: %s\n", expandedPath)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			// Build and process inheritance chain for this file
			chain, err := chainBuilder.buildChain(expandedPath, pathResolver)
			if err != nil {
				if debug {
					fmt.Printf("DEBUG: Failed to build inheritance chain for %s: %v\n", expandedPath, err)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				// If buildChain fails, try loading the file directly as a single-file chain
				// This handles files without inheritance or with parse errors in inherit field
				if debug {
					fmt.Printf("DEBUG: Attempting to load %s as single file\n", expandedPath)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				loadResult, err2 := loadSingleConfigFile(expandedPath)
				if err2 != nil {
					if debug {
						fmt.Printf("DEBUG: Failed to load %s directly: %v\n", expandedPath, err2)
					} // SEMANTIC-TOKEN: DEBUG-OUTPUT
					continue // Skip files we can't load
				}
				// Create a single-file chain
				chain = &inheritanceChain{
					files:   []string{expandedPath},
					visited: make(map[string]bool),
				}
				chain.visited[expandedPath] = true
				// Process this single file
				tempCfg := loadResult.config
				if debug {
					fmt.Printf("DEBUG: Loading config from single file: %s\n", expandedPath)
					if len(tempCfg.ExcludePatterns) > 0 {
						fmt.Printf("DEBUG:   Exclude patterns in this file: %v\n", tempCfg.ExcludePatterns)
					}
					fmt.Printf("DEBUG:   Current exclude patterns before merge: %v\n", finalCfg.ExcludePatterns)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				
				inheritContext := isFirstFile // First file merges with defaults
				mergedCfg, err2 := applyMergeStrategies(finalCfg, tempCfg, inheritContext, loadResult.rawMap)
				if err2 != nil {
					if debug {
						fmt.Printf("DEBUG: Failed to merge config from %s: %v\n", expandedPath, err2)
					} // SEMANTIC-TOKEN: DEBUG-OUTPUT
					continue // Skip problematic merges
				}
				finalCfg = mergedCfg
				if debug {
					fmt.Printf("DEBUG:   Exclude patterns after merge: %v\n", finalCfg.ExcludePatterns)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				foundAny = true
				isFirstFile = false
				continue // Skip the normal chain processing below
			}
			if debug {
				fmt.Printf("DEBUG: Inheritance chain: %v\n", chain.files)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			
			// Determine if this is an inheritance chain (multiple files) or single file
			isInheritanceChain := len(chain.files) > 1
			
			for _, filePath := range chain.files {
				loadResult, err := loadSingleConfigFile(filePath)
				if err != nil {
					if debug {
						fmt.Printf("DEBUG: Failed to load config file %s: %v\n", filePath, err)
					} // SEMANTIC-TOKEN: DEBUG-OUTPUT
					continue // Skip files with errors
				}
				tempCfg := loadResult.config
				if debug {
					fmt.Printf("DEBUG: Loading config from inheritance chain: %s\n", filePath)
					if len(tempCfg.ExcludePatterns) > 0 {
						fmt.Printf("DEBUG:   Exclude patterns in this file: %v\n", tempCfg.ExcludePatterns)
					}
					fmt.Printf("DEBUG:   Current exclude patterns before merge: %v\n", finalCfg.ExcludePatterns)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
				
				// Determine merge context:
				// - Within inheritance chain (multiple files): always merge (inheritContext=true)
				// - First file from searchPaths: merge with defaults per CFG-005 (inheritContext=true)
				// - Subsequent files from searchPaths: override previous files (inheritContext=false)
				// Note: For backward compatibility, use !exclude_patterns prefix to explicitly override
				inheritContext := isInheritanceChain || isFirstFile
				
				mergedCfg, err := applyMergeStrategies(finalCfg, tempCfg, inheritContext, loadResult.rawMap)
				if err != nil {
					if debug {
						fmt.Printf("DEBUG: Failed to merge config from %s: %v\n", filePath, err)
					} // SEMANTIC-TOKEN: DEBUG-OUTPUT
					continue // Skip problematic merges
				}
				fmt.Printf("DIAGNOSTIC: LoadConfigWithInheritance - After merge from %s: exclude_patterns = %v (len=%d)\n", filePath, mergedCfg.ExcludePatterns, len(mergedCfg.ExcludePatterns))
				finalCfg = mergedCfg
				if debug {
					fmt.Printf("DEBUG:   Exclude patterns after merge: %v\n", finalCfg.ExcludePatterns)
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			}
			foundAny = true
			isFirstFile = false // After processing first file, subsequent files override
		}
	}

	if debug {
		fmt.Printf("DEBUG: Final merged exclude patterns (with inheritance): %v\n", finalCfg.ExcludePatterns)
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT

	if !foundAny {
		return DefaultConfig(), nil
	}
	return finalCfg, nil
}

// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]
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
		loadResult, err := loadSingleConfigFile(filePath)
		if err != nil {
			continue // Skip files with errors, continue with chain
		}
		tempCfg := loadResult.config

		// Apply merge strategies and merge into main config
		// Within inheritance chain, array fields default to merge
		mergedCfg, err := applyMergeStrategies(cfg, tempCfg, true, loadResult.rawMap)
		if err != nil {
			return nil, fmt.Errorf("failed to merge config from %s: %w", filePath, err)
		}
		cfg = mergedCfg
	}

	return cfg, nil
}

// configFileLoadResult holds both the Config and the raw map with prefixes preserved
type configFileLoadResult struct {
	config  *Config
	rawMap  map[string]interface{}
}

// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]
// loadSingleConfigFile loads a single configuration file.
// Unmarshals into a map first to preserve merge strategy prefixes (!, +, ^, =),
// then processes merge strategies and converts to Config.
// Returns both the Config and the raw map with prefixes preserved.
func loadSingleConfigFile(configPath string) (*configFileLoadResult, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}
	defer f.Close()

	// Unmarshal into map first to preserve merge strategy prefixes
	var rawMap map[string]interface{}
	d := yaml.NewDecoder(f)
	if err := d.Decode(&rawMap); err != nil {
		return nil, fmt.Errorf("failed to decode config file %s: %w", configPath, err)
	}

	// Process merge strategies to extract clean keys and strategies
	processor := newMergeStrategyProcessor()
	processed, err := processor.processKeys(rawMap)
	if err != nil {
		return nil, fmt.Errorf("failed to process merge strategies: %w", err)
	}

	// Convert processed operations back to a map with clean keys
	cleanMap := make(map[string]interface{})
	for key, op := range processed.operations {
		cleanMap[key] = op.value
	}

	// Now unmarshal the clean map into Config struct
	cfg := DefaultConfig()
	cleanYAML, err := yaml.Marshal(cleanMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal processed config: %w", err)
	}
	if err := yaml.Unmarshal(cleanYAML, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal processed config: %w", err)
	}

	return &configFileLoadResult{config: cfg, rawMap: rawMap}, nil
}

// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]
// applyMergeStrategies applies merge strategies when combining configurations.
// inheritContext indicates if this merge is within an inheritance chain (true) or between sequential files (false).
// Array fields default to merge only within inheritance chains, not between sequential files.
// rawSrcMap is the original map with merge strategy prefixes preserved (can be nil if not available)
func applyMergeStrategies(dst, src *Config, inheritContext bool, rawSrcMap map[string]interface{}) (*Config, error) {
	processor := newMergeStrategyProcessor()

	// Use rawSrcMap if available (preserves merge strategy prefixes), otherwise convert Config to map
	var srcMap map[string]interface{}
	if rawSrcMap != nil {
		srcMap = rawSrcMap
		// Ensure exclude_patterns is in srcMap even if not in rawSrcMap (for backward compatibility)
		if _, exists := srcMap["exclude_patterns"]; !exists && len(src.ExcludePatterns) > 0 {
			srcMap["exclude_patterns"] = src.ExcludePatterns
		}
	} else {
		srcMap = configToMap(src)
	}

	// Process merge strategies
	processed, err := processor.processKeys(srcMap)
	if err != nil {
		return nil, fmt.Errorf("failed to process merge strategies: %w", err)
	}
	
	// Ensure exclude_patterns is in processed.operations if it exists in src but wasn't processed
	// This handles cases where exclude_patterns might not be in rawSrcMap but is in src.Config
	if _, exists := processed.operations["exclude_patterns"]; !exists && len(src.ExcludePatterns) > 0 {
		processed.operations["exclude_patterns"] = &mergeOperation{
			strategy: "override", // Will be changed to "merge" below if inheritContext is true
			value:    src.ExcludePatterns,
			key:      "exclude_patterns",
		}
	}
	
	if debug {
		keys := make([]string, 0, len(processed.operations))
		for k := range processed.operations {
			keys = append(keys, k)
		}
		fmt.Printf("DEBUG: applyMergeStrategies - processed.operations keys: %v\n", keys)
		if op, ok := processed.operations["exclude_patterns"]; ok {
			fmt.Printf("DEBUG: applyMergeStrategies - exclude_patterns operation: strategy=%s, value=%v\n", op.strategy, op.value)
		} else {
			fmt.Printf("DEBUG: applyMergeStrategies - exclude_patterns NOT in processed.operations\n")
		}
	} // SEMANTIC-TOKEN: DEBUG-OUTPUT
	
	// CFG-005: Array fields default to merge (accumulate) strategy to preserve values
	// from defaults and parent configs when child configs add values.
	// This applies to both inheritance chains and sequential file processing.
	// But respect explicit merge strategy prefixes (!, +, ^, =) in all contexts
	// Check if explicit prefix was used in rawSrcMap
	hasExplicitPrefix := false
	if rawSrcMap != nil {
		for origKey := range rawSrcMap {
			if origKey == "!exclude_patterns" || origKey == "+exclude_patterns" || 
			   origKey == "^exclude_patterns" || origKey == "=exclude_patterns" {
				hasExplicitPrefix = true
				break
			}
		}
	}
	
	for key, operation := range processed.operations {
		if key == "exclude_patterns" {
			if debug {
				fmt.Printf("DEBUG: applyMergeStrategies - exclude_patterns before: strategy=%s, hasExplicitPrefix=%v, inheritContext=%v\n", operation.strategy, hasExplicitPrefix, inheritContext)
			} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			// If explicit prefix was used (!, +, ^, =), respect it (don't change to merge)
			// Note: ! gives "replace", + gives "merge", ^ gives "prepend", = gives "default"
			// CFG-005: Array fields default to merge (no prefix) in all contexts
			// If strategy is already "replace", "merge", "prepend", or "default", keep it
			if operation.strategy == "override" && !hasExplicitPrefix {
				// Default override (no prefix), change to merge per CFG-005 requirement
				operation.strategy = "merge"
				if debug {
					fmt.Printf("DEBUG: applyMergeStrategies - exclude_patterns changed to merge strategy (CFG-005: array fields default to merge)\n")
				} // SEMANTIC-TOKEN: DEBUG-OUTPUT
			}
			// For "replace" (!), "merge" (+), "prepend" (^), "default" (=), keep as-is
		}
	}

	// Apply processed configuration
	// Start with a copy of dst (current state), not DefaultConfig()
	// This ensures we preserve any accumulated values from previous merges
	result := &Config{}
	*result = *dst  // Copy the config
	
	// Save dst's exclude_patterns - we'll handle exclude_patterns via applyMergeOperation based on strategy
	dstExcludePatterns := make([]string, len(dst.ExcludePatterns))
	copy(dstExcludePatterns, dst.ExcludePatterns)
	
	// Create a copy of src with exclude_patterns cleared to prevent mergeConfigs from merging them
	// (we'll handle exclude_patterns via merge strategies instead)
	srcCopy := &Config{}
	*srcCopy = *src
	srcCopy.ExcludePatterns = nil
	
	// Merge other settings from srcCopy into result (exclude_patterns skipped)
	mergeConfigs(result, srcCopy)
	
	// For exclude_patterns, restore dst's original patterns so applyMergeOperation can merge based on strategy
	// This ensures we start with the current state (dst) and merge src's patterns into it
	result.ExcludePatterns = dstExcludePatterns

	// Apply source values with merge strategies
	// Use result's current values (not dst's) for merge operations
	// Create resultMap AFTER restoring exclude_patterns so it reflects the correct state
	resultMap := configToMap(result)
	
	// DIAGNOSTIC: Always log exclude_patterns processing (not just when debug is enabled)
	fmt.Printf("\n=== DIAGNOSTIC: applyMergeStrategies ===\n")
	fmt.Printf("inheritContext: %v\n", inheritContext)
	fmt.Printf("dst.ExcludePatterns: %v (len=%d)\n", dst.ExcludePatterns, len(dst.ExcludePatterns))
	fmt.Printf("src.ExcludePatterns: %v (len=%d)\n", src.ExcludePatterns, len(src.ExcludePatterns))
	fmt.Printf("dstExcludePatterns (saved): %v (len=%d)\n", dstExcludePatterns, len(dstExcludePatterns))
	fmt.Printf("result.ExcludePatterns (after copy): %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
	fmt.Printf("result.ExcludePatterns (after restore): %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
	
	keys := make([]string, 0, len(processed.operations))
	for k := range processed.operations {
		keys = append(keys, k)
	}
	fmt.Printf("processed.operations keys: %v\n", keys)
	if _, hasExcludePatterns := processed.operations["exclude_patterns"]; hasExcludePatterns {
		fmt.Printf("✓ exclude_patterns found in processed.operations\n")
		op := processed.operations["exclude_patterns"]
		fmt.Printf("  operation.strategy: %s\n", op.strategy)
		fmt.Printf("  operation.value: %v (type: %T)\n", op.value, op.value)
		if opSlice, ok := op.value.([]string); ok {
			fmt.Printf("  operation.value as []string: %v (len=%d)\n", opSlice, len(opSlice))
		}
		if opSlice, ok := op.value.([]interface{}); ok {
			fmt.Printf("  operation.value as []interface{}: %v (len=%d)\n", opSlice, len(opSlice))
		}
	} else {
		fmt.Printf("✗ exclude_patterns NOT found in processed.operations. Keys: %v\n", keys)
	}
	
	fmt.Printf("resultMap[\"exclude_patterns\"]: %v\n", resultMap["exclude_patterns"])
	
	for key, operation := range processed.operations {
		fmt.Printf("DIAGNOSTIC: Processing key: %s, strategy: %s\n", key, operation.strategy)
		// Skip metadata fields that are not part of the Config struct
		// These fields are used for inheritance processing but shouldn't be merged into the config
		if key == "inherit" {
			fmt.Printf("DIAGNOSTIC: Skipping metadata field: %s\n", key)
			continue
		}
		// Check if this is a known config field before processing
		if !isKnownConfigField(key) {
			fmt.Printf("DIAGNOSTIC: Skipping unknown config field: %s\n", key)
			continue
		}
		// Get current value from result (which has dst merged in)
		currentValue := resultMap[key]
		if key == "exclude_patterns" {
			fmt.Printf("\n--- Processing exclude_patterns ---\n")
			fmt.Printf("currentValue: %v (type: %T)\n", currentValue, currentValue)
			if cvSlice, ok := currentValue.([]string); ok {
				fmt.Printf("currentValue as []string: %v (len=%d)\n", cvSlice, len(cvSlice))
			}
			fmt.Printf("operation.value: %v (type: %T)\n", operation.value, operation.value)
			fmt.Printf("operation.strategy: %s\n", operation.strategy)
			fmt.Printf("result.ExcludePatterns BEFORE applyMergeOperation: %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
		}
		if debug && key == "exclude_patterns" {
			fmt.Printf("DEBUG: applyMergeStrategies - key: %s, currentValue: %v, operation.value: %v, strategy: %s\n", 
				key, currentValue, operation.value, operation.strategy)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		err := applyMergeOperation(result, key, operation, currentValue)
		if err != nil {
			// Check if error is due to unknown field (shouldn't happen if isKnownConfigField works correctly)
			if strings.Contains(err.Error(), "unknown config field") {
				fmt.Printf("DIAGNOSTIC: Skipping unknown config field (from error): %s\n", key)
				continue
			}
			fmt.Printf("DIAGNOSTIC: Error applying merge operation for %s: %v\n", key, err)
			return nil, fmt.Errorf("failed to apply merge operation for %s: %w", key, err)
		}
		// Update resultMap to reflect the merge for subsequent operations
		resultMap = configToMap(result)
		if key == "exclude_patterns" {
			fmt.Printf("result.ExcludePatterns AFTER applyMergeOperation: %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
			fmt.Printf("--- End processing exclude_patterns ---\n\n")
		}
		if debug && key == "exclude_patterns" {
			fmt.Printf("DEBUG: applyMergeStrategies - after merge, result.ExcludePatterns: %v\n", result.ExcludePatterns)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT
		fmt.Printf("DIAGNOSTIC: Completed processing key: %s\n", key)
	}
	
	fmt.Printf("=== END DIAGNOSTIC: applyMergeStrategies ===\n\n")

	return result, nil
}

// CFG-005: See specification.md - Configuration Inheritance [DECISION:core-functionality]

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
		"backup_dir_path":      cfg.BackupDirPath,
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
	// CFG-002: See specification.md - Configuration Merging [DECISION:discovery]
	// For arrays, merge by appending with deduplication
	if key == "exclude_patterns" {
		fmt.Printf("=== DIAGNOSTIC: applyMerge for exclude_patterns ===\n")
		fmt.Printf("value: %v (type: %T)\n", value, value)
		fmt.Printf("dstValue: %v (type: %T)\n", dstValue, dstValue)
		fmt.Printf("result.ExcludePatterns: %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
	}
	
	// Convert []interface{} to []string if needed (common from YAML unmarshaling)
	var srcSlice []string
	var ok bool
	if srcSlice, ok = value.([]string); !ok {
		// Try to convert from []interface{}
		if ifaceSlice, ok2 := value.([]interface{}); ok2 {
			srcSlice = make([]string, 0, len(ifaceSlice))
			for _, item := range ifaceSlice {
				if str, ok3 := item.(string); ok3 {
					srcSlice = append(srcSlice, str)
				} else {
					// If conversion fails, fall back to setConfigField
					if key == "exclude_patterns" {
						fmt.Printf("Failed to convert []interface{} element to string, using setConfigField directly\n")
						fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
					}
					return setConfigField(result, key, value)
				}
			}
			if key == "exclude_patterns" {
				fmt.Printf("Converted []interface{} to []string: %v (len=%d)\n", srcSlice, len(srcSlice))
			}
		} else {
			// Not a slice at all, fall back to setConfigField
			if key == "exclude_patterns" {
				fmt.Printf("value is not []string or []interface{}, using setConfigField directly\n")
				fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
			}
			return setConfigField(result, key, value)
		}
	}
	
	// Only process if we have source values to merge
	if len(srcSlice) > 0 {
		if key == "exclude_patterns" {
			fmt.Printf("srcSlice: %v (len=%d)\n", srcSlice, len(srcSlice))
		}
		// Handle nil or missing dstValue by getting current value from result
		if dstValue == nil {
			if key == "exclude_patterns" {
				fmt.Printf("dstValue is nil, getting from result config\n")
			}
			// Get current value from result config
			resultMap := configToMap(result)
			dstValue = resultMap[key]
			if key == "exclude_patterns" {
				fmt.Printf("dstValue from resultMap: %v (type: %T)\n", dstValue, dstValue)
			}
		}
		
		var dstSlice []string
		if dstValue != nil {
			if dstSlice, ok = dstValue.([]string); !ok {
				// Try to convert from []interface{}
				if ifaceSlice, ok2 := dstValue.([]interface{}); ok2 {
					dstSlice = make([]string, 0, len(ifaceSlice))
					for _, item := range ifaceSlice {
						if str, ok3 := item.(string); ok3 {
							dstSlice = append(dstSlice, str)
						}
					}
					if key == "exclude_patterns" {
						fmt.Printf("Converted dstValue []interface{} to []string: %v (len=%d)\n", dstSlice, len(dstSlice))
					}
				} else {
					// dstValue is not a slice, start with srcSlice
					if key == "exclude_patterns" {
						fmt.Printf("dstValue is not []string or []interface{}, starting with srcSlice\n")
						fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
					}
					return setConfigField(result, key, srcSlice)
				}
			}
		}
		
		// Merge srcSlice into dstSlice (or start with srcSlice if dstSlice is nil/empty)
		if len(dstSlice) > 0 {
			if key == "exclude_patterns" {
				fmt.Printf("dstSlice: %v (len=%d)\n", dstSlice, len(dstSlice))
			}
			// Create map to track existing patterns for deduplication
			existingMap := make(map[string]bool)
			for _, pattern := range dstSlice {
				existingMap[pattern] = true
			}
			
			// Append new patterns that aren't already present
			merged := make([]string, len(dstSlice), len(dstSlice)+len(srcSlice))
			copy(merged, dstSlice)
			for _, pattern := range srcSlice {
				if !existingMap[pattern] {
					merged = append(merged, pattern)
					existingMap[pattern] = true
					if key == "exclude_patterns" {
						fmt.Printf("  Added pattern: %q\n", pattern)
					}
				} else {
					if key == "exclude_patterns" {
						fmt.Printf("  Skipped duplicate pattern: %q\n", pattern)
					}
				}
			}
			if key == "exclude_patterns" {
				fmt.Printf("merged result: %v (len=%d)\n", merged, len(merged))
				fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
			}
			return setConfigField(result, key, merged)
		}
		// If dstSlice is nil or empty, start with srcSlice (no existing patterns)
		if key == "exclude_patterns" {
			fmt.Printf("dstSlice is nil or empty, starting with srcSlice\n")
			fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
		}
		return setConfigField(result, key, srcSlice)
	}
	// If srcSlice is empty, don't merge - preserve existing value
	if key == "exclude_patterns" {
		fmt.Printf("srcSlice is empty, preserving existing value\n")
		fmt.Printf("=== END DIAGNOSTIC: applyMerge ===\n\n")
	}
	return nil // No merge needed, preserve existing value
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
	if key == "exclude_patterns" {
		fmt.Printf("=== DIAGNOSTIC: applyReplace for exclude_patterns ===\n")
		fmt.Printf("value: %v (type: %T)\n", value, value)
		fmt.Printf("result.ExcludePatterns BEFORE: %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
	}
	err := setConfigField(result, key, value)
	if key == "exclude_patterns" {
		fmt.Printf("result.ExcludePatterns AFTER: %v (len=%d)\n", result.ExcludePatterns, len(result.ExcludePatterns))
		if err != nil {
			fmt.Printf("setConfigField error: %v\n", err)
		}
		fmt.Printf("=== END DIAGNOSTIC: applyReplace ===\n\n")
	}
	return err
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
		} else if ifaceSlice, ok := value.([]interface{}); ok {
			// Convert []interface{} to []string (common from YAML unmarshaling)
			strSlice := make([]string, 0, len(ifaceSlice))
			for _, item := range ifaceSlice {
				if str, ok := item.(string); ok {
					strSlice = append(strSlice, str)
				}
			}
			cfg.ExcludePatterns = strSlice
		}
	case "include_git_info":
		if b, ok := value.(bool); ok {
			cfg.IncludeGitInfo = b
		}
	case "backup_dir_path":
		if s, ok := value.(string); ok {
			cfg.BackupDirPath = s
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

// isKnownConfigField checks if a field name is a known config field
func isKnownConfigField(key string) bool {
	knownFields := map[string]bool{
		"archive_dir_path":     true,
		"use_current_dir_name": true,
		"exclude_patterns":     true,
		"include_git_info":     true,
		"backup_dir_path":      true,
		"skip_broken_symlinks": true,
		"status_created_archive": true,
		"status_disk_full":      true,
	}
	return knownFields[key]
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

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Automatic field discovery
// configFieldInfo represents metadata about a configuration field discovered through reflection.
// It provides complete information about field structure, type, and documentation.
type configFieldInfo struct {
	Name      string       // Field name in Go struct
	YAMLName  string       // YAML tag name for configuration file
	Type      string       // Go type as string
	Kind      reflect.Kind // Reflect kind for type handling
	Value     interface{}  // Current field value
	IsPointer bool         // Whether field is a pointer type
	IsSlice   bool         // Whether field is a slice type
	IsStruct  bool         // Whether field is a struct type
	Category  string       // Field category (basic, status, format, template, etc.)
	Path      string       // Full path for nested fields (e.g., "verification.verify_on_create")
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Source tracking extension
// ConfigValueWithMetadata extends ConfigValue with complete field information and inheritance tracking.
type ConfigValueWithMetadata struct {
	ConfigValue
	FieldInfo        configFieldInfo
	InheritanceChain []string // Chain of files that contributed to this value
	MergeStrategy    string   // Strategy used to merge this value (override, append, prepend, etc.)
	IsOverridden     bool     // Whether this value was overridden from default
	ConflictSources  []string // Sources that had conflicting values
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 1: Automatic Field Discovery System
// IMPLEMENTATION-REF: CFG-006 Subtask 6.1: Add reflection result caching
// GetAllConfigFields discovers all configuration fields using reflection.
// It provides comprehensive field enumeration without manual maintenance.
// Performance optimized with caching to reduce reflection overhead.
func GetAllConfigFields(cfg *Config) []configFieldInfo {
	// Try to get cached results first for performance
	if cachedFields := globalFieldCache.getCachedFields(); cachedFields != nil {
		// Update values in cached fields with current config values
		return updateFieldValues(cachedFields, cfg)
	}

	// Cache miss - perform expensive reflection
	var fields []configFieldInfo

	// Get reflection information about the Config struct
	configType := reflect.TypeOf(*cfg)
	configValue := reflect.ValueOf(*cfg)

	// Recursively discover all fields
	fields = append(fields, reflectConfigFields(configType, configValue, "", "")...)

	// Sort fields by name for consistent output
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	// Cache the field metadata (without values) for future use
	fieldMetadata := make([]configFieldInfo, len(fields))
	for i, field := range fields {
		// Store field metadata without specific values for caching
		fieldMetadata[i] = configFieldInfo{
			Name:      field.Name,
			YAMLName:  field.YAMLName,
			Type:      field.Type,
			Kind:      field.Kind,
			Value:     nil, // Values will be updated per call
			IsPointer: field.IsPointer,
			IsSlice:   field.IsSlice,
			IsStruct:  field.IsStruct,
			Category:  field.Category,
			Path:      field.Path,
		}
	}
	globalFieldCache.setCachedFields(fieldMetadata)

	return fields
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.1: Add reflection result caching
// updateFieldValues updates cached field metadata with current config values.
// This avoids expensive reflection while keeping values current.
func updateFieldValues(cachedFields []configFieldInfo, cfg *Config) []configFieldInfo {
	result := make([]configFieldInfo, len(cachedFields))
	configValue := reflect.ValueOf(*cfg)

	for i, field := range cachedFields {
		// Copy cached metadata
		result[i] = field

		// Update value using field path
		if fieldValue, err := getFieldValueByPath(configValue, field.Path); err == nil {
			result[i].Value = fieldValue
		} else {
			// Fallback to zero value if path resolution fails
			result[i].Value = getZeroValueForKind(field.Kind)
		}
	}

	return result
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.1: Add reflection result caching
// getFieldValueByPath retrieves a field value from a struct using dot-separated path.
func getFieldValueByPath(structValue reflect.Value, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	currentValue := structValue

	for _, part := range parts {
		// Handle pointer dereferencing
		if currentValue.Kind() == reflect.Ptr {
			if currentValue.IsNil() {
				return nil, fmt.Errorf("nil pointer in path %s at %s", path, part)
			}
			currentValue = currentValue.Elem()
		}

		// Get field by name
		if currentValue.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct in path %s at %s, got %s", path, part, currentValue.Kind())
		}

		field := currentValue.FieldByName(part)
		if !field.IsValid() {
			return nil, fmt.Errorf("field %s not found in path %s", part, path)
		}

		currentValue = field
	}

	// Handle final value extraction
	if currentValue.Kind() == reflect.Ptr {
		if currentValue.IsNil() {
			return nil, nil
		}
		return currentValue.Elem().Interface(), nil
	}

	return currentValue.Interface(), nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Step 1.2: Implement reflectConfigFields() function
// reflectConfigFields recursively discovers fields in structs, handling nested types.
func reflectConfigFields(structType reflect.Type, structValue reflect.Value, prefix string, category string) []configFieldInfo {
	var fields []configFieldInfo

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Build field path for nested structures
		fieldPath := field.Name
		if prefix != "" {
			fieldPath = prefix + "." + field.Name
		}

		// Get YAML tag name
		yamlTag := field.Tag.Get("yaml")
		yamlName := strings.Split(yamlTag, ",")[0] // Remove options like "omitempty"
		if yamlName == "" {
			yamlName = strings.ToLower(field.Name) // Default to lowercase field name
		}

		// Determine field category
		fieldCategory := determineFieldCategory(field.Name, category)

		// Check if this is a struct type (including pointer to struct)
		isStruct := field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct)

		// Handle nested structs
		if isStruct && field.Type != reflect.TypeOf(yaml.Node{}) {
			// For pointer to struct, get the actual struct
			var nestedType reflect.Type
			var nestedValue reflect.Value

			if field.Type.Kind() == reflect.Ptr {
				if !fieldValue.IsNil() {
					nestedType = field.Type.Elem()
					nestedValue = fieldValue.Elem()
				} else {
					// Create zero value for nil pointer
					nestedType = field.Type.Elem()
					nestedValue = reflect.Zero(nestedType)
				}
			} else {
				nestedType = field.Type
				nestedValue = fieldValue
			}

			// Recursively process nested struct fields
			nestedFields := reflectConfigFields(nestedType, nestedValue, fieldPath, fieldCategory)
			fields = append(fields, nestedFields...)
		} else {
			// For non-struct fields, create field info with proper type detection
			var actualType string
			var actualKind reflect.Kind
			var actualValue interface{}

			// Handle pointer types properly
			if field.Type.Kind() == reflect.Ptr {
				actualType = field.Type.Elem().String()
				actualKind = field.Type.Elem().Kind()
				if !fieldValue.IsNil() {
					actualValue = fieldValue.Elem().Interface()
				} else {
					actualValue = nil
				}
			} else {
				actualType = field.Type.String()
				actualKind = field.Type.Kind()
				actualValue = fieldValue.Interface()
			}

			fieldInfo := configFieldInfo{
				Name:      field.Name,
				YAMLName:  yamlName,
				Type:      actualType,
				Kind:      actualKind,
				Value:     actualValue,
				IsPointer: field.Type.Kind() == reflect.Ptr,
				IsSlice:   field.Type.Kind() == reflect.Slice,
				IsStruct:  isStruct,
				Category:  fieldCategory,
				Path:      fieldPath,
			}

			fields = append(fields, fieldInfo)
		}
	}

	return fields
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Step 1.5: Create field filtering and categorization
// determineFieldCategory categorizes configuration fields by their purpose and type.
func determineFieldCategory(fieldName string, parentCategory string) string {
	if parentCategory != "" {
		return parentCategory
	}

	// Categorize based on field name patterns
	switch {
	case strings.HasPrefix(fieldName, "Status"):
		return "status_codes"
	case strings.HasPrefix(fieldName, "Format"):
		return "format_strings"
	case strings.HasPrefix(fieldName, "Template"):
		return "template_strings"
	case strings.HasPrefix(fieldName, "Pattern"):
		return "regex_patterns"
	case fieldName == "Verification":
		return "verification"
	case fieldName == "Inherit":
		return "inheritance"
	case strings.Contains(fieldName, "Backup"):
		return "backup_settings"
	case strings.Contains(fieldName, "Archive"):
		return "archive_settings"
	default:
		return "basic_settings"
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 2: Enhanced Source Tracking Extension
// GetAllConfigValuesWithSources provides comprehensive configuration visibility using reflection.
// It replaces the manual field enumeration with automatic discovery and enhanced source tracking.
func GetAllConfigValuesWithSources(cfg *Config, root string) []ConfigValueWithMetadata {
	// Get all fields using reflection
	fields := GetAllConfigFields(cfg)
	defaultCfg := DefaultConfig()
	defaultFields := GetAllConfigFields(defaultCfg)

	// Create mapping of field names to default values
	defaultValues := make(map[string]interface{})
	for _, field := range defaultFields {
		defaultValues[field.Path] = field.Value
	}

	// Determine configuration source
	configSource := determineConfigSource(root)
	getSource := createSourceDeterminer(configSource)

	var results []ConfigValueWithMetadata

	for _, field := range fields {
		// Skip struct fields themselves (keep their children)
		if field.IsStruct && !field.IsPointer {
			continue
		}

		// Get default value for comparison
		defaultValue, hasDefault := defaultValues[field.Path]
		if !hasDefault {
			defaultValue = getZeroValueForKind(field.Kind)
		}

		// Determine source of this field
		source := getSource(field.Value, defaultValue)
		if debug && field.YAMLName == "exclude_patterns" {
			fmt.Printf("DEBUG: Source determination for exclude_patterns:\n")
			fmt.Printf("DEBUG:   Current value: %v (type: %T)\n", field.Value, field.Value)
			fmt.Printf("DEBUG:   Default value: %v (type: %T)\n", defaultValue, defaultValue)
			fmt.Printf("DEBUG:   Determined source: %s\n", source)
			fmt.Printf("DEBUG:   Config source path: %s\n", configSource)
		} // SEMANTIC-TOKEN: DEBUG-OUTPUT

		// Format value as string
		valueStr := formatFieldValue(field.Value, field.Kind)

		// Track inheritance chain for this field
		inheritanceChain, mergeStrategy, conflictSources := trackInheritanceChain(field.Path, cfg, root)

		// Ensure conflict sources is always a slice
		if conflictSources == nil {
			conflictSources = []string{}
		}

		// Create enhanced config value
		configValue := ConfigValueWithMetadata{
			ConfigValue: ConfigValue{
				Name:   field.YAMLName,
				Value:  valueStr,
				Source: source,
			},
			FieldInfo:        field,
			InheritanceChain: inheritanceChain,
			MergeStrategy:    mergeStrategy,
			IsOverridden:     source != "default",
			ConflictSources:  conflictSources,
		}

		results = append(results, configValue)
	}

	// Sort results alphabetically by YAML name
	sort.Slice(results, func(i, j int) bool {
		return results[i].ConfigValue.Name < results[j].ConfigValue.Name
	})

	return results
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Step 3.2: Implement type-aware value formatting
// formatFieldValue formats configuration values based on their Go type.
func formatFieldValue(value interface{}, kind reflect.Kind) string {
	if value == nil {
		return "<nil>"
	}

	switch kind {
	case reflect.Bool:
		return boolToString(value.(bool))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", value)
	case reflect.String:
		return value.(string)
	case reflect.Slice:
		// Handle string slices specifically
		if slice, ok := value.([]string); ok {
			if len(slice) == 0 {
				return "[]"
			}
			return fmt.Sprintf("[%s]", strings.Join(slice, ", "))
		}
		return fmt.Sprintf("%v", value)
	case reflect.Ptr:
		// For pointers, format the pointed-to value
		ptrValue := reflect.ValueOf(value)
		if ptrValue.IsNil() {
			return "<nil>"
		}
		return formatFieldValue(ptrValue.Elem().Interface(), ptrValue.Elem().Kind())
	default:
		return fmt.Sprintf("%v", value)
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Helper functions for reflection system
// getZeroValueForKind returns the appropriate zero value for a given reflect.Kind.
func getZeroValueForKind(kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Bool:
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(0)
	case reflect.String:
		return ""
	case reflect.Slice:
		return []string{}
	case reflect.Ptr:
		return nil
	default:
		return nil
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.1: Add reflection result caching

// ConfigFieldCache provides thread-safe caching of configuration field discovery results.
// It significantly reduces reflection overhead for repeated GetAllConfigFields() calls.
type ConfigFieldCache struct {
	mu         sync.RWMutex
	fields     []configFieldInfo
	structHash uint64    // Hash of Config struct to detect schema changes
	lastUpdate time.Time // When cache was last updated
	valid      bool      // Whether cached data is valid
}

// ConfigFilter provides filtering options for configuration field enumeration.
// It enables lazy evaluation by only processing fields that match filter criteria.
type ConfigFilter struct {
	FieldPatterns []string // Field name patterns to include (glob patterns)
	Categories    []string // Field categories to include
	OverridesOnly bool     // Show only non-default values
	SourceTypes   []string // Show only specific source types (environment, config, default)
}

// Global cache instance for configuration field discovery
var globalFieldCache = &ConfigFieldCache{}

// getConfigStructHash computes a hash of the Config struct type to detect schema changes.
// This enables cache invalidation when the struct definition is modified.
func getConfigStructHash() uint64 {
	h := fnv.New64a()
	configType := reflect.TypeOf(Config{})
	writeTypeToHash(h, configType)
	return h.Sum64()
}

// writeTypeToHash recursively writes type information to hash for cache validation.
func writeTypeToHash(h hash.Hash64, t reflect.Type) {
	// Write type name and kind
	h.Write([]byte(t.String()))
	h.Write([]byte{byte(t.Kind())})

	// For struct types, hash all field information
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			h.Write([]byte(field.Name))
			h.Write([]byte(field.Tag))
			writeTypeToHash(h, field.Type)
		}
	}

	// For pointer and slice types, hash element type
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		writeTypeToHash(h, t.Elem())
	}
}

// getCachedFields retrieves fields from cache if valid, otherwise returns nil.
// Thread-safe read operation with minimal lock contention.
func (c *ConfigFieldCache) getCachedFields() []configFieldInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.valid {
		return nil
	}

	// Validate cache against current struct hash
	currentHash := getConfigStructHash()
	if c.structHash != currentHash {
		// Struct has changed, cache is invalid
		return nil
	}

	// Return copy of cached fields to prevent modification
	result := make([]configFieldInfo, len(c.fields))
	copy(result, c.fields)
	return result
}

// setCachedFields stores fields in cache with current struct hash.
// Thread-safe write operation that invalidates and repopulates cache.
func (c *ConfigFieldCache) setCachedFields(fields []configFieldInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store copy of fields to prevent external modification
	c.fields = make([]configFieldInfo, len(fields))
	copy(c.fields, fields)

	c.structHash = getConfigStructHash()
	c.lastUpdate = time.Now()
	c.valid = true
}

// invalidateCache marks the cache as invalid, forcing refresh on next access.
func (c *ConfigFieldCache) invalidateCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.valid = false
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.2: Implement lazy source evaluation
// GetConfigValuesWithSourcesFiltered provides configuration visibility with filtering.
// It only resolves sources for fields that match the filter criteria for better performance.
func GetConfigValuesWithSourcesFiltered(cfg *Config, root string, filter *ConfigFilter) []ConfigValueWithMetadata {
	// Get all fields using reflection
	fields := GetAllConfigFields(cfg)

	// Pre-filter fields to avoid expensive source resolution for unwanted fields
	filteredFields := applyConfigFilter(fields, filter)

	if len(filteredFields) == 0 {
		return []ConfigValueWithMetadata{}
	}

	// Only resolve sources for filtered fields
	defaultCfg := DefaultConfig()
	defaultFields := GetAllConfigFields(defaultCfg)

	// Create mapping of field names to default values (for all fields for comparison)
	defaultValues := make(map[string]interface{})
	for _, defaultField := range defaultFields {
		defaultValues[defaultField.Path] = defaultField.Value
	}

	// Determine configuration source
	configSource := determineConfigSource(root)
	getSource := createSourceDeterminer(configSource)

	var results []ConfigValueWithMetadata

	for _, field := range filteredFields {
		// Skip struct fields themselves (keep their children)
		if field.IsStruct && !field.IsPointer {
			continue
		}

		// Get default value for comparison
		defaultValue, hasDefault := defaultValues[field.Path]
		if !hasDefault {
			defaultValue = getZeroValueForKind(field.Kind)
		}

		// Determine source of this field
		source := getSource(field.Value, defaultValue)

		// Apply overrides-only filter if specified
		if filter != nil && filter.OverridesOnly {
			if source == "default" {
				continue // Skip default values when only showing overrides
			}
		}

		// Apply source type filter if specified
		if filter != nil && len(filter.SourceTypes) > 0 {
			sourceMatches := false
			for _, allowedSource := range filter.SourceTypes {
				if source == allowedSource {
					sourceMatches = true
					break
				}
			}
			if !sourceMatches {
				continue // Skip sources not in the allowed list
			}
		}

		// Format value as string
		valueStr := formatFieldValue(field.Value, field.Kind)

		// Create enhanced config value
		configValue := ConfigValueWithMetadata{
			ConfigValue: ConfigValue{
				Name:   field.YAMLName,
				Value:  valueStr,
				Source: source,
			},
			FieldInfo: field,
			// TODO: Add inheritance chain and merge strategy tracking
			InheritanceChain: []string{}, // Placeholder for future CFG-005 integration
			MergeStrategy:    "override", // Placeholder for future merge strategy tracking
			IsOverridden:     source != "default",
			ConflictSources:  []string{}, // Placeholder for conflict detection
		}

		results = append(results, configValue)
	}

	return results
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.2: Implement lazy source evaluation
// applyConfigFilter filters configuration fields based on filter criteria.
func applyConfigFilter(fields []configFieldInfo, filter *ConfigFilter) []configFieldInfo {
	if filter == nil {
		return fields // No filtering
	}

	var filtered []configFieldInfo

	for _, field := range fields {
		// Apply field pattern filter
		if len(filter.FieldPatterns) > 0 {
			patternMatches := false
			for _, pattern := range filter.FieldPatterns {
				if matched, err := filepath.Match(pattern, field.Name); err == nil && matched {
					patternMatches = true
					break
				}
				// Also check YAML name and path
				if matched, err := filepath.Match(pattern, field.YAMLName); err == nil && matched {
					patternMatches = true
					break
				}
				if matched, err := filepath.Match(pattern, field.Path); err == nil && matched {
					patternMatches = true
					break
				}
			}
			if !patternMatches {
				continue
			}
		}

		// Apply category filter
		if len(filter.Categories) > 0 {
			categoryMatches := false
			for _, allowedCategory := range filter.Categories {
				if field.Category == allowedCategory {
					categoryMatches = true
					break
				}
			}
			if !categoryMatches {
				continue
			}
		}

		filtered = append(filtered, field)
	}

	return filtered
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.3: Create incremental resolution support
// GetConfigFieldByPattern retrieves specific configuration fields matching a pattern.
// This enables efficient single-field or pattern-based queries without full enumeration.
func GetConfigFieldByPattern(cfg *Config, pattern string) ([]configFieldInfo, error) {
	// Use caching for performance
	allFields := GetAllConfigFields(cfg)

	var matchingFields []configFieldInfo

	for _, field := range allFields {
		// Check if field matches pattern
		if matched, err := filepath.Match(pattern, field.Name); err != nil {
			return nil, fmt.Errorf("invalid pattern %s: %v", pattern, err)
		} else if matched {
			matchingFields = append(matchingFields, field)
			continue
		}

		// Also check YAML name
		if matched, err := filepath.Match(pattern, field.YAMLName); err == nil && matched {
			matchingFields = append(matchingFields, field)
			continue
		}

		// Also check field path for nested field access
		if matched, err := filepath.Match(pattern, field.Path); err == nil && matched {
			matchingFields = append(matchingFields, field)
		}
	}

	return matchingFields, nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.3: Create incremental resolution support
// GetConfigFieldValue retrieves a single configuration field value with complete metadata.
// This is the most efficient way to access a specific configuration field.
func GetConfigFieldValue(cfg *Config, fieldPath string) (ConfigValueWithMetadata, error) {
	// Try to get the field value directly using path resolution
	configValue := reflect.ValueOf(*cfg)
	fieldValue, err := getFieldValueByPath(configValue, fieldPath)
	if err != nil {
		return ConfigValueWithMetadata{}, fmt.Errorf("field path %s not found: %v", fieldPath, err)
	}

	// Get field metadata from cache or reflection
	allFields := GetAllConfigFields(cfg)
	var targetField *configFieldInfo

	for _, field := range allFields {
		if field.Path == fieldPath {
			targetField = &field
			break
		}
	}

	if targetField == nil {
		return ConfigValueWithMetadata{}, fmt.Errorf("field metadata not found for path %s", fieldPath)
	}

	// Get default value for source determination
	defaultCfg := DefaultConfig()
	defaultValue := getZeroValueForKind(targetField.Kind)

	if defaultConfigValue, err := getFieldValueByPath(reflect.ValueOf(*defaultCfg), fieldPath); err == nil {
		defaultValue = defaultConfigValue
	}

	// Determine source (simplified for single field access)
	var source string
	if reflect.DeepEqual(fieldValue, defaultValue) {
		source = "default"
	} else {
		// Check if environment variable exists
		envVarName := strings.ToUpper(strings.ReplaceAll(fieldPath, ".", "_"))
		if _, exists := os.LookupEnv(envVarName); exists {
			source = "environment"
		} else {
			source = "config"
		}
	}

	// Format value
	valueStr := formatFieldValue(fieldValue, targetField.Kind)

	return ConfigValueWithMetadata{
		ConfigValue: ConfigValue{
			Name:   targetField.YAMLName,
			Value:  valueStr,
			Source: source,
		},
		FieldInfo:        *targetField,
		InheritanceChain: []string{}, // Simplified for single field access
		MergeStrategy:    "override",
		IsOverridden:     source != "default",
		ConflictSources:  []string{},
	}, nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6.3: Create incremental resolution support
// HasConfigField checks if a configuration field exists without full field enumeration.
func HasConfigField(cfg *Config, fieldPath string) bool {
	configValue := reflect.ValueOf(*cfg)
	_, err := getFieldValueByPath(configValue, fieldPath)
	return err == nil
}

// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
// GitConfig defines Git integration configuration options.
// It controls Git repository detection, information extraction, and behavior.
type GitConfig struct {
	// Basic Git integration settings
	Enabled         bool `yaml:"enabled"`           // Enable/disable Git integration
	IncludeInfo     bool `yaml:"include_info"`      // Include Git info in operations (legacy: include_git_info)
	ShowDirtyStatus bool `yaml:"show_dirty_status"` // Show dirty status indicator (legacy: show_git_dirty_status)

	// Git command configuration
	Command          string `yaml:"command"`           // Git command path (default: "git")
	WorkingDirectory string `yaml:"working_directory"` // Working directory for Git operations (default: ".")

	// Git behavior settings
	RequireCleanRepo  bool `yaml:"require_clean_repo"` // Fail operations if repository is dirty
	AutoDetectRepo    bool `yaml:"auto_detect_repo"`   // Automatically detect Git repositories
	IncludeSubmodules bool `yaml:"include_submodules"` // Include submodule information

	// Git information inclusion
	IncludeBranch bool `yaml:"include_branch"` // Include branch name in operations
	IncludeHash   bool `yaml:"include_hash"`   // Include commit hash in operations
	IncludeStatus bool `yaml:"include_status"` // Include working directory status

	// Git command timeouts and limits
	CommandTimeout    string `yaml:"command_timeout"`     // Timeout for Git commands (default: "30s")
	MaxSubmoduleDepth int    `yaml:"max_submodule_depth"` // Maximum submodule recursion depth
}

// GIT-005: See specification.md - Git Configuration Integration [DECISION:maintenance]
// DefaultGitConfig returns a GitConfig with sensible defaults
func DefaultGitConfig() *GitConfig {
	return &GitConfig{
		Enabled:           true,
		IncludeInfo:       false, // Legacy compatibility
		ShowDirtyStatus:   false, // Legacy compatibility
		Command:           "git",
		WorkingDirectory:  ".",
		RequireCleanRepo:  false,
		AutoDetectRepo:    true,
		IncludeSubmodules: false,
		IncludeBranch:     true,
		IncludeHash:       true,
		IncludeStatus:     true,
		CommandTimeout:    "30s",
		MaxSubmoduleDepth: 3,
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 5: Full Inheritance Chain Tracking
// trackInheritanceChain tracks the complete inheritance chain for a configuration field.
// It analyzes the inheritance hierarchy to determine which files contributed to each value.
func trackInheritanceChain(fieldPath string, cfg *Config, root string) ([]string, string, []string) {
	// Get inheritance chain using existing system
	fileOps := &configFileOperations{}
	pathResolver := newPathResolver(fileOps)
	chainBuilder := newInheritanceChainBuilder(fileOps)

	// Build the inheritance chain
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

	if primaryConfigPath == "" {
		// No config file found, return default source
		return []string{"default"}, "default", []string{}
	}

	// Build inheritance chain
	chain, err := chainBuilder.buildChain(primaryConfigPath, pathResolver)
	if err != nil {
		// Fallback to single file
		return []string{primaryConfigPath}, "config", []string{}
	}

	// Track value through inheritance chain
	var inheritanceChain []string
	var mergeStrategy string
	var conflictSources []string

	// Start with default configuration
	currentCfg := DefaultConfig()
	defaultValue, _ := getFieldValueByPath(reflect.ValueOf(currentCfg), fieldPath)

	// Process files in inheritance order (parents first)
	for i, filePath := range chain.files {
		loadResult, err := loadSingleConfigFile(filePath)
		if err != nil {
			continue // Skip files with errors
		}
		tempCfg := loadResult.config

		// Get value from this file
		tempValue, _ := getFieldValueByPath(reflect.ValueOf(tempCfg), fieldPath)

		// Check if this file contributes to the final value
		if !reflect.DeepEqual(tempValue, defaultValue) {
			inheritanceChain = append(inheritanceChain, filePath)

			// Determine merge strategy based on field characteristics
			mergeStrategy = determineMergeStrategyForField(fieldPath, tempValue, defaultValue)

			// Check for conflicts with previous values
			if i > 0 && !reflect.DeepEqual(tempValue, defaultValue) {
				// This is a simplified conflict detection
				// In a full implementation, this would track actual conflicts
				if hasConflict(tempValue, defaultValue) {
					conflictSources = append(conflictSources, filePath)
				}
			}
		}

		// Apply merge strategies and merge into current config
		// Within inheritance chain, array fields default to merge
		mergedCfg, err := applyMergeStrategies(currentCfg, tempCfg, true, loadResult.rawMap)
		if err != nil {
			continue // Skip problematic merges
		}
		currentCfg = mergedCfg
	}

	// If no inheritance chain found, use default
	if len(inheritanceChain) == 0 {
		return []string{"default"}, "default", []string{}
	}

	return inheritanceChain, mergeStrategy, conflictSources
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 5: Helper function for inheritance tracking
// determineMergeStrategyForField determines the merge strategy used for a specific field.
func determineMergeStrategyForField(fieldPath string, newValue, oldValue interface{}) string {
	// Check for exclude_patterns first (should use append strategy)
	if strings.Contains(fieldPath, "exclude_patterns") || strings.Contains(fieldPath, "inherit") {
		// Slice fields typically use append strategy
		if reflect.TypeOf(newValue).Kind() == reflect.Slice {
			return "append"
		}
	}

	// Check for prepend strategy indicators - patterns fields should use prepend
	// But avoid matching exclude_patterns
	if strings.Contains(fieldPath, "patterns") && !strings.Contains(fieldPath, "exclude_patterns") {
		// Pattern fields use prepend strategy
		return "prepend"
	}

	// Check for specific field patterns that indicate different strategies
	if strings.Contains(fieldPath, "git") {
		// Git configuration fields use merge strategy for nested structs
		return "merge"
	}

	if strings.Contains(fieldPath, "verification") {
		// Verification fields use override strategy
		return "override"
	}

	// Check if value is different from default (override strategy)
	if oldValue != nil && !reflect.DeepEqual(newValue, oldValue) {
		// Determine if this is a default value application
		if isZeroValue(oldValue) && !isZeroValue(newValue) {
			return "default"
		}
		return "override"
	}

	// Default strategy for non-nil values
	if newValue != nil {
		return "override"
	}

	// Default strategy
	return "override"
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 5: Helper function for conflict detection
// hasConflict checks if there's a conflict between two values.
func hasConflict(newValue, oldValue interface{}) bool {
	// This is a simplified conflict detection
	// In a full implementation, this would check for actual conflicts
	// For now, we consider it a conflict if values are different and both non-zero
	if reflect.DeepEqual(newValue, oldValue) {
		return false
	}

	// Check if both values are non-zero
	newIsZero := isZeroValue(newValue)
	oldIsZero := isZeroValue(oldValue)

	// Conflict if both are non-zero and different
	return !newIsZero && !oldIsZero
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 6: Enhanced merge strategy detection
// detectMergeStrategyFromField analyzes field characteristics to determine merge strategy.
func detectMergeStrategyFromField(field configFieldInfo) string {
	// Analyze field type and characteristics
	switch {
	case field.IsSlice:
		// Slice fields typically use append strategy
		return "append"
	case field.IsStruct:
		// Struct fields use merge strategy
		return "merge"
	case field.IsPointer:
		// Pointer fields use override strategy
		return "override"
	case strings.Contains(field.Name, "Pattern"):
		// Pattern fields might use prepend strategy
		return "prepend"
	case strings.Contains(field.Name, "Git"):
		// Git fields use merge strategy
		return "merge"
	case strings.Contains(field.Name, "Verification"):
		// Verification fields use override strategy
		return "override"
	default:
		// Default to override strategy
		return "override"
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Configuration Validation and Documentation
// ConfigValidationRule represents a validation rule for a configuration field.
type ConfigValidationRule struct {
	FieldPath   string
	RuleType    string // "range", "pattern", "required", "custom"
	MinValue    interface{}
	MaxValue    interface{}
	Pattern     string
	CustomFunc  func(interface{}) error
	Description string
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Field documentation structure
// ConfigFieldDocumentation provides documentation for a configuration field.
type ConfigFieldDocumentation struct {
	FieldPath     string
	Description   string
	DefaultValue  string
	ValidRange    string
	Examples      []string
	RelatedFields []string
	Deprecated    bool
	DeprecatedIn  string
	Replacement   string
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Validation and documentation functions
// validateConfigField validates a configuration field against defined rules.
func validateConfigField(fieldPath string, value interface{}, rules []ConfigValidationRule) []error {
	var errors []error

	for _, rule := range rules {
		if rule.FieldPath != fieldPath {
			continue
		}

		switch rule.RuleType {
		case "range":
			if err := validateRange(value, rule.MinValue, rule.MaxValue); err != nil {
				errors = append(errors, fmt.Errorf("field %s: %w", fieldPath, err))
			}
		case "pattern":
			if err := validatePattern(value, rule.Pattern); err != nil {
				errors = append(errors, fmt.Errorf("field %s: %w", fieldPath, err))
			}
		case "required":
			if err := validateRequired(value); err != nil {
				errors = append(errors, fmt.Errorf("field %s: %w", fieldPath, err))
			}
		case "custom":
			if rule.CustomFunc != nil {
				if err := rule.CustomFunc(value); err != nil {
					errors = append(errors, fmt.Errorf("field %s: %w", fieldPath, err))
				}
			}
		}
	}

	return errors
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Validation helper functions
// validateRange validates that a value is within the specified range.
func validateRange(value interface{}, min, max interface{}) error {
	switch v := value.(type) {
	case int:
		if minInt, ok := min.(int); ok && v < minInt {
			return fmt.Errorf("value %d is below minimum %d", v, minInt)
		}
		if maxInt, ok := max.(int); ok && v > maxInt {
			return fmt.Errorf("value %d is above maximum %d", v, maxInt)
		}
	case string:
		if minStr, ok := min.(string); ok && v < minStr {
			return fmt.Errorf("value %s is below minimum %s", v, minStr)
		}
		if maxStr, ok := max.(string); ok && v > maxStr {
			return fmt.Errorf("value %s is above maximum %s", v, maxStr)
		}
	}
	return nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Pattern validation
// validatePattern validates a value against a regex pattern.
func validatePattern(value interface{}, pattern string) error {
	if pattern == "" {
		return nil
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("pattern validation only supports string values")
	}

	matched, err := regexp.MatchString(pattern, strValue)
	if err != nil {
		return fmt.Errorf("invalid pattern %s: %w", pattern, err)
	}

	if !matched {
		return fmt.Errorf("value %s does not match pattern %s", strValue, pattern)
	}

	return nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Required field validation
// validateRequired validates that a required field is not empty.
func validateRequired(value interface{}) error {
	if isZeroValue(value) {
		return fmt.Errorf("field is required but has zero value")
	}
	return nil
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Documentation generation
// GenerateConfigDocumentation generates comprehensive documentation for all configuration fields.
func GenerateConfigDocumentation(cfg *Config) []ConfigFieldDocumentation {
	fields := GetAllConfigFields(cfg)
	var docs []ConfigFieldDocumentation

	for _, field := range fields {
		doc := ConfigFieldDocumentation{
			FieldPath:     field.YAMLName, // Use YAMLName instead of Go struct Name
			Description:   generateFieldDescription(field),
			DefaultValue:  formatFieldValue(field.Value, field.Kind),
			ValidRange:    generateValidRange(field),
			Examples:      generateFieldExamples(field),
			RelatedFields: generateRelatedFields(field),
		}
		docs = append(docs, doc)
	}

	return docs
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Documentation helper functions
// generateFieldDescription generates a description for a configuration field.
func generateFieldDescription(field configFieldInfo) string {
	switch {
	case strings.Contains(field.Name, "Archive"):
		return "Archive operation configuration"
	case strings.Contains(field.Name, "Backup"):
		return "File backup operation configuration"
	case strings.Contains(field.Name, "Status"):
		return "Status code for operation result"
	case strings.Contains(field.Name, "Format"):
		return "Output formatting string"
	case strings.Contains(field.Name, "Template"):
		return "Template-based output formatting"
	case strings.Contains(field.Name, "Pattern"):
		return "Regex pattern for matching"
	case strings.Contains(field.Name, "Git"):
		return "Git integration configuration"
	case strings.Contains(field.Name, "Verification"):
		return "Archive verification configuration"
	default:
		return "Configuration parameter"
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Range generation
// generateValidRange generates valid range information for a field.
func generateValidRange(field configFieldInfo) string {
	switch field.Kind {
	case reflect.Int:
		return "Integer value"
	case reflect.String:
		return "String value"
	case reflect.Bool:
		return "Boolean value (true/false)"
	case reflect.Slice:
		return "Array of values"
	default:
		return "Any valid value"
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Example generation
// generateFieldExamples generates example values for a field.
func generateFieldExamples(field configFieldInfo) []string {
	switch {
	case strings.Contains(field.Name, "ArchiveDirPath"):
		return []string{"./archives", "/var/backups"}
	case strings.Contains(field.Name, "ExcludePatterns"):
		return []string{"*.tmp", "node_modules/", ".git/"}
	case strings.Contains(field.Name, "Format"):
		return []string{"Created archive: %s", "Archive %s created successfully"}
	case strings.Contains(field.Name, "Status"):
		return []string{"0", "1", "2"}
	default:
		return []string{}
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 7: Related fields generation
// generateRelatedFields generates a list of related configuration fields.
func generateRelatedFields(field configFieldInfo) []string {
	var related []string

	// Add related fields based on field characteristics
	switch {
	case strings.Contains(field.Name, "Archive"):
		related = append(related, "archive_dir_path", "use_current_dir_name")
	case strings.Contains(field.Name, "Backup"):
		related = append(related, "backup_dir_path", "use_current_dir_name_for_files")
	case strings.Contains(field.Name, "Git"):
		related = append(related, "git.enabled", "git.include_info")
	case strings.Contains(field.Name, "Verification"):
		related = append(related, "verification.verify_on_create", "verification.checksum_algorithm")
	}

	return related
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: Documentation Plan
// GenerateMarkdownDocumentation generates markdown documentation for configuration.
func GenerateMarkdownDocumentation(cfg *Config) string {
	docs := GenerateConfigDocumentation(cfg)

	var buf strings.Builder

	// Write header
	buf.WriteString("# BkpDir Configuration Reference\n\n")
	buf.WriteString("This document provides a comprehensive reference for all BkpDir configuration options.\n\n")

	// Group fields by category
	categories := groupFieldsByCategory(docs)

	for category, fields := range categories {
		buf.WriteString(fmt.Sprintf("## %s\n\n", category))

		for _, field := range fields {
			buf.WriteString(fmt.Sprintf("### %s\n\n", field.FieldPath))
			buf.WriteString(fmt.Sprintf("%s\n\n", field.Description))

			if field.DefaultValue != "" {
				buf.WriteString(fmt.Sprintf("**Default:** `%s`\n\n", field.DefaultValue))
			}

			if field.ValidRange != "" {
				buf.WriteString(fmt.Sprintf("**Valid Range:** %s\n\n", field.ValidRange))
			}

			if len(field.Examples) > 0 {
				buf.WriteString("**Examples:**\n")
				for _, example := range field.Examples {
					buf.WriteString(fmt.Sprintf("- `%s`\n", example))
				}
				buf.WriteString("\n")
			}

			if len(field.RelatedFields) > 0 {
				buf.WriteString("**Related Fields:**\n")
				for _, related := range field.RelatedFields {
					buf.WriteString(fmt.Sprintf("- `%s`\n", related))
				}
				buf.WriteString("\n")
			}

			if field.Deprecated {
				buf.WriteString(fmt.Sprintf("**Deprecated:** This field is deprecated since %s\n", field.DeprecatedIn))
				if field.Replacement != "" {
					buf.WriteString(fmt.Sprintf("**Replacement:** Use `%s` instead\n", field.Replacement))
				}
				buf.WriteString("\n")
			}
		}
	}

	return buf.String()
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: Documentation helper functions
// groupFieldsByCategory groups configuration fields by their category.
func groupFieldsByCategory(docs []ConfigFieldDocumentation) map[string][]ConfigFieldDocumentation {
	categories := make(map[string][]ConfigFieldDocumentation)

	for _, doc := range docs {
		category := determineFieldCategory(doc.FieldPath, "")
		categories[category] = append(categories[category], doc)
	}

	return categories
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: JSON Schema generation
// GenerateJSONSchema generates a JSON schema for the configuration structure.
func GenerateJSONSchema(cfg *Config) string {
	fields := GetAllConfigFields(cfg)

	var buf strings.Builder
	buf.WriteString("{\n")
	buf.WriteString(`  "$schema": "http://json-schema.org/draft-07/schema#",` + "\n")
	buf.WriteString(`  "title": "BkpDir Configuration Schema",` + "\n")
	buf.WriteString(`  "type": "object",` + "\n")
	buf.WriteString(`  "properties": {` + "\n")

	for i, field := range fields {
		if i > 0 {
			buf.WriteString(",\n")
		}

		buf.WriteString(fmt.Sprintf(`    "%s": {`, field.YAMLName))
		buf.WriteString(fmt.Sprintf(`"type": "%s"`, getJSONSchemaType(field.Kind)))

		description := generateFieldDescription(field)
		if description != "" {
			buf.WriteString(fmt.Sprintf(`, "description": "%s"`, description))
		}

		defaultValue := formatFieldValue(field.Value, field.Kind)
		if defaultValue != "" {
			buf.WriteString(fmt.Sprintf(`, "default": %s`, formatJSONValue(defaultValue, field.Kind)))
		}

		buf.WriteString("}")
	}

	buf.WriteString("\n  },\n")
	buf.WriteString(`  "additionalProperties": false` + "\n")
	buf.WriteString("}\n")

	return buf.String()
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: JSON Schema helper functions
// getJSONSchemaType converts Go reflect.Kind to JSON schema type.
func getJSONSchemaType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice:
		return "array"
	case reflect.Struct:
		return "object"
	default:
		return "string"
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: JSON value formatting
// formatJSONValue formats a value for JSON schema.
func formatJSONValue(value string, kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return fmt.Sprintf(`"%s"`, value)
	case reflect.Bool:
		return value
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value
	default:
		return fmt.Sprintf(`"%s"`, value)
	}
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: Configuration validation
// ValidateConfiguration validates the entire configuration against defined rules.
func ValidateConfiguration(cfg *Config, rules []ConfigValidationRule) []error {
	var errors []error
	fields := GetAllConfigFields(cfg)

	for _, field := range fields {
		fieldErrors := validateConfigField(field.Path, field.Value, rules)
		errors = append(errors, fieldErrors...)
	}

	return errors
}

// CFG-006: See specification.md - Configuration Performance [DECISION:maintenance]
// IMPLEMENTATION-REF: CFG-006 Subtask 8: Default validation rules
// GetDefaultValidationRules returns default validation rules for configuration fields.
func GetDefaultValidationRules() []ConfigValidationRule {
	return []ConfigValidationRule{
		{
			FieldPath:   "archive_dir_path",
			RuleType:    "required",
			Description: "Archive directory path is required",
		},
		{
			FieldPath:   "status_created_archive",
			RuleType:    "range",
			MinValue:    0,
			MaxValue:    255,
			Description: "Status code must be between 0 and 255",
		},
		{
			FieldPath:   "pattern_archive_filename",
			RuleType:    "pattern",
			Pattern:     `^[^*?:"<>|]+$`,
			Description: "Archive filename pattern must not contain invalid characters",
		},
	}
}
