package main

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type VerificationConfig struct {
	VerifyOnCreate    bool   `yaml:"verify_on_create"`
	ChecksumAlgorithm string `yaml:"checksum_algorithm"`
}

type Config struct {
	// Basic settings
	ArchiveDirPath    string              `yaml:"archive_dir_path"`
	UseCurrentDirName bool                `yaml:"use_current_dir_name"`
	ExcludePatterns   []string            `yaml:"exclude_patterns"`
	IncludeGitInfo    bool                `yaml:"include_git_info"`
	Verification      *VerificationConfig `yaml:"verification"`

	// Status codes
	StatusCreatedArchive                        int `yaml:"status_created_archive"`
	StatusFailedToCreateArchiveDirectory        int `yaml:"status_failed_to_create_archive_directory"`
	StatusDirectoryIsIdenticalToExistingArchive int `yaml:"status_directory_is_identical_to_existing_archive"`
	StatusDirectoryNotFound                     int `yaml:"status_directory_not_found"`
	StatusInvalidDirectoryType                  int `yaml:"status_invalid_directory_type"`
	StatusPermissionDenied                      int `yaml:"status_permission_denied"`
	StatusDiskFull                              int `yaml:"status_disk_full"`
	StatusConfigError                           int `yaml:"status_config_error"`

	// Printf-style format strings
	FormatCreatedArchive   string `yaml:"format_created_archive"`
	FormatIdenticalArchive string `yaml:"format_identical_archive"`
	FormatListArchive      string `yaml:"format_list_archive"`
	FormatConfigValue      string `yaml:"format_config_value"`
	FormatDryRunArchive    string `yaml:"format_dry_run_archive"`
	FormatError            string `yaml:"format_error"`

	// Template-based format strings
	TemplateCreatedArchive   string `yaml:"template_created_archive"`
	TemplateIdenticalArchive string `yaml:"template_identical_archive"`
	TemplateListArchive      string `yaml:"template_list_archive"`
	TemplateConfigValue      string `yaml:"template_config_value"`
	TemplateDryRunArchive    string `yaml:"template_dry_run_archive"`
	TemplateError            string `yaml:"template_error"`

	// Regex patterns
	PatternArchiveFilename string `yaml:"pattern_archive_filename"`
	PatternConfigLine      string `yaml:"pattern_config_line"`
	PatternTimestamp       string `yaml:"pattern_timestamp"`
}

type ConfigValue struct {
	Name   string
	Value  string
	Source string
}

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

		// Status codes
		StatusCreatedArchive:                        0,
		StatusFailedToCreateArchiveDirectory:        31,
		StatusDirectoryIsIdenticalToExistingArchive: 0,
		StatusDirectoryNotFound:                     20,
		StatusInvalidDirectoryType:                  21,
		StatusPermissionDenied:                      22,
		StatusDiskFull:                              30,
		StatusConfigError:                           10,

		// Printf-style format strings
		FormatCreatedArchive:   "Created archive: %s\n",
		FormatIdenticalArchive: "Directory is identical to existing archive: %s\n",
		FormatListArchive:      "%s (created: %s)\n",
		FormatConfigValue:      "%s: %s (source: %s)\n",
		FormatDryRunArchive:    "Would create archive: %s\n",
		FormatError:            "Error: %s\n",

		// Template-based format strings
		TemplateCreatedArchive:   "Created archive: %{path}\n",
		TemplateIdenticalArchive: "Directory is identical to existing archive: %{path}\n",
		TemplateListArchive:      "%{path} (created: %{creation_time})\n",
		TemplateConfigValue:      "%{name}: %{value} (source: %{source})\n",
		TemplateDryRunArchive:    "Would create archive: %{path}\n",
		TemplateError:            "Error: %{message}\n",

		// Regex patterns
		PatternArchiveFilename: `(?P<prefix>[^-]*)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})(?:=(?P<branch>[^=]+))?(?:=(?P<hash>[^=]+))?(?:=(?P<note>.+))?\.zip`,
		PatternConfigLine:      `(?P<name>[^:]+):\s*(?P<value>[^(]+)\s*\(source:\s*(?P<source>[^)]+)\)`,
		PatternTimestamp:       `(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})\s+(?P<hour>\d{2}):(?P<minute>\d{2}):(?P<second>\d{2})`,
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

func mergeConfigs(dst, src *Config) {
	// Only merge non-default values from src into dst
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

	// Merge verification config
	if src.Verification != nil {
		if dst.Verification == nil {
			dst.Verification = &VerificationConfig{}
		}
		if src.Verification.VerifyOnCreate != DefaultConfig().Verification.VerifyOnCreate {
			dst.Verification.VerifyOnCreate = src.Verification.VerifyOnCreate
		}
		if src.Verification.ChecksumAlgorithm != DefaultConfig().Verification.ChecksumAlgorithm {
			dst.Verification.ChecksumAlgorithm = src.Verification.ChecksumAlgorithm
		}
	}

	// Merge status codes (only if different from default)
	defaultCfg := DefaultConfig()
	if src.StatusCreatedArchive != defaultCfg.StatusCreatedArchive {
		dst.StatusCreatedArchive = src.StatusCreatedArchive
	}
	if src.StatusFailedToCreateArchiveDirectory != defaultCfg.StatusFailedToCreateArchiveDirectory {
		dst.StatusFailedToCreateArchiveDirectory = src.StatusFailedToCreateArchiveDirectory
	}
	if src.StatusDirectoryIsIdenticalToExistingArchive != defaultCfg.StatusDirectoryIsIdenticalToExistingArchive {
		dst.StatusDirectoryIsIdenticalToExistingArchive = src.StatusDirectoryIsIdenticalToExistingArchive
	}
	if src.StatusDirectoryNotFound != defaultCfg.StatusDirectoryNotFound {
		dst.StatusDirectoryNotFound = src.StatusDirectoryNotFound
	}
	if src.StatusInvalidDirectoryType != defaultCfg.StatusInvalidDirectoryType {
		dst.StatusInvalidDirectoryType = src.StatusInvalidDirectoryType
	}
	if src.StatusPermissionDenied != defaultCfg.StatusPermissionDenied {
		dst.StatusPermissionDenied = src.StatusPermissionDenied
	}
	if src.StatusDiskFull != defaultCfg.StatusDiskFull {
		dst.StatusDiskFull = src.StatusDiskFull
	}
	if src.StatusConfigError != defaultCfg.StatusConfigError {
		dst.StatusConfigError = src.StatusConfigError
	}

	// Merge format strings (only if different from default)
	if src.FormatCreatedArchive != defaultCfg.FormatCreatedArchive {
		dst.FormatCreatedArchive = src.FormatCreatedArchive
	}
	if src.FormatIdenticalArchive != defaultCfg.FormatIdenticalArchive {
		dst.FormatIdenticalArchive = src.FormatIdenticalArchive
	}
	if src.FormatListArchive != defaultCfg.FormatListArchive {
		dst.FormatListArchive = src.FormatListArchive
	}
	if src.FormatConfigValue != defaultCfg.FormatConfigValue {
		dst.FormatConfigValue = src.FormatConfigValue
	}
	if src.FormatDryRunArchive != defaultCfg.FormatDryRunArchive {
		dst.FormatDryRunArchive = src.FormatDryRunArchive
	}
	if src.FormatError != defaultCfg.FormatError {
		dst.FormatError = src.FormatError
	}

	// Merge template strings (only if different from default)
	if src.TemplateCreatedArchive != defaultCfg.TemplateCreatedArchive {
		dst.TemplateCreatedArchive = src.TemplateCreatedArchive
	}
	if src.TemplateIdenticalArchive != defaultCfg.TemplateIdenticalArchive {
		dst.TemplateIdenticalArchive = src.TemplateIdenticalArchive
	}
	if src.TemplateListArchive != defaultCfg.TemplateListArchive {
		dst.TemplateListArchive = src.TemplateListArchive
	}
	if src.TemplateConfigValue != defaultCfg.TemplateConfigValue {
		dst.TemplateConfigValue = src.TemplateConfigValue
	}
	if src.TemplateDryRunArchive != defaultCfg.TemplateDryRunArchive {
		dst.TemplateDryRunArchive = src.TemplateDryRunArchive
	}
	if src.TemplateError != defaultCfg.TemplateError {
		dst.TemplateError = src.TemplateError
	}

	// Merge regex patterns (only if different from default)
	if src.PatternArchiveFilename != defaultCfg.PatternArchiveFilename {
		dst.PatternArchiveFilename = src.PatternArchiveFilename
	}
	if src.PatternConfigLine != defaultCfg.PatternConfigLine {
		dst.PatternConfigLine = src.PatternConfigLine
	}
	if src.PatternTimestamp != defaultCfg.PatternTimestamp {
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

func GetConfigValues(cfg *Config) []ConfigValue {
	// This would be used by the --config command to display all configuration values
	// For now, return basic values - this can be expanded
	return []ConfigValue{
		{Name: "archive_dir_path", Value: cfg.ArchiveDirPath, Source: "config"},
		{Name: "use_current_dir_name", Value: boolToString(cfg.UseCurrentDirName), Source: "config"},
		{Name: "include_git_info", Value: boolToString(cfg.IncludeGitInfo), Source: "config"},
	}
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
