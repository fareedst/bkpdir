// This file is part of bkpdir
//
// Package main implements the BkpDir CLI application for directory archiving and file backup.
// It provides commands for creating full and incremental directory archives, backing up individual files,
// and managing archive verification and configuration.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	version = "1.2.0"
)

var (
	compileDate = "2024-03-20" // This is a placeholder - actual value is set during build via -ldflags
	platform    = "unknown"    // This is a placeholder - actual value is set during build via -ldflags
)

var (
	dryRun     bool
	note       string
	showConfig bool
	listFile   string
)

// Long description for root command
const rootLongDesc = `bkpdir version %s (compiled %s) [%s]

BkpDir is a command-line tool for archiving directories and backing up individual files on macOS and Linux. 
It supports full and incremental directory backups, individual file backups, customizable exclusion patterns, 
Git-aware archive naming, and archive verification.`

func main() {
	rootCmd := &cobra.Command{
		Use:     "bkpdir",
		Short:   "Directory archiving and file backup CLI for macOS and Linux",
		Long:    fmt.Sprintf(rootLongDesc, version, compileDate, platform),
		Version: fmt.Sprintf("%s (compiled %s) [%s]", version, compileDate, platform),
		Example: `  # Create a full directory archive
  bkpdir create "Initial backup"
  bkpdir full -n "Initial backup"  # backward compatibility

  # Create an incremental directory archive with verification
  bkpdir create --incremental "Changes after feature X" -v
  bkpdir inc -n "Changes after feature X" -v  # backward compatibility

  # Create a file backup
  bkpdir backup myfile.txt "Before changes"

  # List all directory archives
  bkpdir list

  # List backups for a specific file
  bkpdir --list myfile.txt

  # Verify a specific archive with checksums
  bkpdir verify backup-2024-03-20.zip -c

  # Show configuration
  bkpdir config
  bkpdir --config  # backward compatibility`,
		Run: func(cmd *cobra.Command, args []string) {
			// Handle --config flag when no subcommand is provided (backward compatibility)
			if showConfig {
				handleConfigCommand()
				return
			}
			// Handle --list flag for file backups
			if cmd.Flags().Changed("list") {
				handleListFileBackupsCommand(args)
				return
			}
			// If no config flag and no subcommand, show help
			cmd.Help()
		},
	}

	// Set the version template to show version in help output
	versionTemplate := fmt.Sprintf("bkpdir version %s (compiled %s) [%s]\n",
		version, compileDate, platform)
	rootCmd.SetVersionTemplate(versionTemplate)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false,
		"Show what would be done without creating archives")
	rootCmd.PersistentFlags().BoolVar(&showConfig, "config", false,
		"Display configuration values and exit (backward compatibility)")
	rootCmd.PersistentFlags().StringVar(&listFile, "list", "",
		"List backups for a specific file")

	// Add commands - new specification-compliant commands first
	rootCmd.AddCommand(createCmd())
	rootCmd.AddCommand(configCmd())

	// Add backward compatibility commands
	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())

	// Add other commands
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(backupCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleConfigCommand() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	cfg, err := LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(cfg.StatusConfigError)
	}

	formatter := NewOutputFormatter(cfg)

	// Display configuration file paths
	configPaths := getConfigSearchPaths()
	expandedPaths := make([]string, len(configPaths))
	for i, path := range configPaths {
		expandedPath := expandPath(path)
		if !filepath.IsAbs(expandedPath) {
			expandedPath = filepath.Join(cwd, expandedPath)
		}
		expandedPaths[i] = expandedPath
	}

	configPathsStr := strings.Join(expandedPaths, ":")
	formatter.PrintConfigValue("config", configPathsStr, "default")

	// Get all configuration values with their sources
	configValues := GetConfigValuesWithSources(cfg, cwd)

	// Display each configuration value
	for _, cv := range configValues {
		formatter.PrintConfigValue(cv.Name, cv.Value, cv.Source)
	}
}

func handleCreateCommand() {
	// Implementation
}

func handleListCommand() {
	// Requirement: List Archives - Display all archives in the archive directory
	// Specification: Shows each archive with path and creation time using configurable format
	// Specification: Shows verification status if available: [VERIFIED], [FAILED], or [UNVERIFIED]
	// Specification: Archives are sorted by creation time (most recent first)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	cfg, err := LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(cfg.StatusConfigError)
	}

	formatter := NewOutputFormatter(cfg)

	if err := ListArchivesEnhanced(cfg, formatter); err != nil {
		exitCode := HandleArchiveError(err, cfg, formatter)
		os.Exit(exitCode)
	}
}

func handleVerifyCommand() {
	// Implementation
}

func handleVersionCommand() {
	// Implementation
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config [KEY] [VALUE]",
		Short: "Display or modify configuration values",
		Long: `Display configuration values or set a specific configuration value.

Examples:
  # Display all configuration values
  bkpdir config

  # Set a configuration value
  bkpdir config archive_dir_path /custom/archive/path
  bkpdir config include_git_info false`,
		Args: cobra.MaximumNArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				// Display configuration
				handleConfigCommand()
			} else if len(args) == 2 {
				// Set configuration value
				handleConfigSetCommand(args[0], args[1])
			} else {
				fmt.Fprintf(os.Stderr, "Error: config set requires both KEY and VALUE\n")
				fmt.Fprintf(os.Stderr, "Usage: bkpdir config [KEY] [VALUE]\n")
				os.Exit(1)
			}
		},
	}
	return cmd
}

func createCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new archive",
		Run: func(*cobra.Command, []string) {
			handleCreateCommand()
		},
	}
	return cmd
}

func fullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full [NOTE]",
		Short: "Create a full archive of the current directory",
		Long: `Create a complete ZIP archive of the current directory. The archive will be stored in the 
archive directory with a timestamp. If the directory is identical to the most recent archive, 
no new archive is created.

Before creating an archive, the command compares the directory with its most recent archive.
If the directory is identical to the most recent archive, no new archive is created.`,
		Example: `  # Create a full archive
  bkpdir full

  # Create a full archive with a note
  bkpdir full "Before changes"

  # Show what would be archived without creating archive
  bkpdir full -d`,
		Args: cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx := context.Background()
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			// Use note from flag if provided, otherwise use positional argument
			archiveNote := note
			if archiveNote == "" && len(args) > 0 {
				archiveNote = args[0]
			}

			if err := CreateFullArchiveWithContext(ctx, cfg, archiveNote, dryRun, false); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the archive name")
	return cmd
}

func incCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inc [NOTE]",
		Short: "Create an incremental archive of the current directory",
		Long: `Create an incremental ZIP archive containing only files changed since the last full archive.
The archive will be stored in the archive directory with a timestamp. If no files have changed,
no new archive is created.

Before creating an archive, the command compares the directory with its most recent archive.
If the directory is identical to the most recent archive, no new archive is created.`,
		Example: `  # Create an incremental archive
  bkpdir inc

  # Create an incremental archive with a note
  bkpdir inc "After changes"

  # Show what would be archived without creating archive
  bkpdir inc -d`,
		Args: cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx := context.Background()
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			// Use note from flag if provided, otherwise use positional argument
			archiveNote := note
			if archiveNote == "" && len(args) > 0 {
				archiveNote = args[0]
			}

			if err := CreateIncrementalArchiveWithContext(ctx, cfg, archiveNote, dryRun, false); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the archive name")
	return cmd
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List archives",
		Run: func(*cobra.Command, []string) {
			handleListCommand()
		},
	}
	return cmd
}

func verifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify archives",
		Run: func(*cobra.Command, []string) {
			handleVerifyCommand()
		},
	}
	return cmd
}

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(*cobra.Command, []string) {
			handleVersionCommand()
		},
	}
	return cmd
}

// ArchiveOptions holds parameters for archive creation functions
type ArchiveOptions struct {
	Context   context.Context
	Config    *Config
	Formatter *OutputFormatter
	Note      string
	DryRun    bool
	Verify    bool
}

// CreateFullArchiveEnhanced creates a full archive of the current directory with enhanced error handling
// and resource management. It supports dry-run mode and optional verification.
func CreateFullArchiveEnhanced(opts ArchiveOptions) error {
	return CreateFullArchive(opts.Config, opts.Note, opts.DryRun, opts.Verify)
}

// CreateIncrementalArchiveEnhanced creates an incremental archive containing only files changed since
// the last full archive. It supports dry-run mode and optional verification.
func CreateIncrementalArchiveEnhanced(opts ArchiveOptions) error {
	return CreateIncrementalArchive(opts.Config, opts.Note, opts.DryRun, opts.Verify)
}

// ListArchivesEnhanced displays all archives in the archive directory with enhanced formatting
// and error handling.
func ListArchivesEnhanced(cfg *Config, formatter *OutputFormatter) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}

	archives, err := ListArchives(archiveDir)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to list archives", 1, err)
	}

	if len(archives) == 0 {
		fmt.Printf("No archives found in %s\n", archiveDir)
		return nil
	}

	// Requirement: Archives are sorted by creation time (most recent first)
	sort.Slice(archives, func(i, j int) bool {
		return archives[i].CreationTime.After(archives[j].CreationTime)
	})

	for _, a := range archives {
		status := ""
		if a.VerificationStatus != nil {
			if a.VerificationStatus.IsVerified {
				status = " [VERIFIED]"
			} else {
				status = " [FAILED]"
			}
		} else {
			status = " [UNVERIFIED]"
		}

		// Use enhanced formatting with extraction if possible
		creationTime := a.CreationTime.Format("2006-01-02 15:04:05")
		output := formatter.FormatListArchiveWithExtraction(a.Name, creationTime)
		// Remove trailing newline from output to add status on same line
		output = strings.TrimSuffix(output, "\n")
		fmt.Printf("%s%s\n", output, status)
	}

	return nil
}

// VerifyOptions holds parameters for archive verification functions
type VerifyOptions struct {
	Config       *Config
	Formatter    *OutputFormatter
	ArchiveName  string
	WithChecksum bool
}

// VerifyArchiveEnhanced verifies the integrity of an archive with optional checksum verification.
// It provides enhanced error handling and reporting.
func VerifyArchiveEnhanced(opts VerifyOptions) error {
	archiveDir, err := getArchiveDirectory(opts.Config)
	if err != nil {
		return err
	}

	if opts.ArchiveName != "" {
		return verifySingleArchive(opts, archiveDir)
	}
	return verifyAllArchives(opts, archiveDir)
}

// getArchiveDirectory determines the archive directory path
func getArchiveDirectory(cfg *Config) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", NewArchiveErrorWithCause("Failed to get current directory",
			cfg.StatusDirectoryNotFound, err)
	}

	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	return archiveDir, nil
}

// verifySingleArchive verifies a specific archive
func verifySingleArchive(opts VerifyOptions, archiveDir string) error {
	archivePath := filepath.Join(archiveDir, opts.ArchiveName)
	archive := &Archive{
		Name: opts.ArchiveName,
		Path: archivePath,
	}

	status, err := performVerification(archive.Path, opts.WithChecksum)
	if err != nil {
		return err
	}

	return handleVerificationResult(archive, status, opts.ArchiveName)
}

// verifyAllArchives verifies all archives in the directory
func verifyAllArchives(opts VerifyOptions, archiveDir string) error {
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to list archives", 1, err)
	}

	allPassed := true
	for _, archive := range archives {
		status, err := performVerification(archive.Path, opts.WithChecksum)
		if err != nil {
			fmt.Printf("Archive %s verification failed: %v\n", archive.Name, err)
			allPassed = false
			continue
		}

		if err := handleVerificationResult(&archive, status, archive.Name); err != nil {
			allPassed = false
		}
	}

	if !allPassed {
		return NewArchiveError("Some archives failed verification", 1)
	}
	return nil
}

// performVerification performs the actual verification based on type
func performVerification(archivePath string, withChecksum bool) (*VerificationStatus, error) {
	if withChecksum {
		status, err := VerifyChecksums(archivePath)
		if err != nil {
			return nil, NewArchiveErrorWithCause("Archive checksum verification failed", 1, err)
		}
		return status, nil
	}

	status, err := VerifyArchive(archivePath)
	if err != nil {
		return nil, NewArchiveErrorWithCause("Archive verification failed", 1, err)
	}
	return status, nil
}

// handleVerificationResult handles the result of verification
func handleVerificationResult(archive *Archive, status *VerificationStatus, name string) error {
	// Store verification status
	if err := StoreVerificationStatus(archive, status); err != nil {
		// Don't fail if we can't store status, just warn
		fmt.Printf("Warning: Could not store verification status for %s: %v\n", name, err)
	}

	if status.IsVerified {
		fmt.Printf("Archive %s verified successfully\n", name)
		return nil
	}

	fmt.Printf("Archive %s verification failed:\n", name)
	for _, errMsg := range status.Errors {
		fmt.Printf("  - %s\n", errMsg)
	}
	return NewArchiveError("Archive verification failed", 1)
}

func handleListFileBackupsCommand(args []string) {
	var filePath string
	if listFile != "" {
		filePath = listFile
	} else if len(args) > 0 {
		filePath = args[0]
	} else {
		fmt.Fprintf(os.Stderr, "Error: file path required for --list command\n")
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	cfg, err := LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(cfg.StatusConfigError)
	}

	formatter := NewOutputFormatter(cfg)

	if err := ListFileBackupsEnhanced(cfg, formatter, filePath); err != nil {
		exitCode := HandleArchiveError(err, cfg, formatter)
		os.Exit(exitCode)
	}
}

func backupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup [FILE_PATH] [NOTE]",
		Short: "Create a backup of a single file",
		Long: `Create a backup of the specified file. The backup will be stored in the backup directory
with a timestamp. If the file is identical to the most recent backup, no new backup is created.

Before creating a backup, the command compares the file with its most recent backup.
If the file is identical to the most recent backup, no new backup is created.`,
		Example: `  # Create a file backup
  bkpdir backup myfile.txt

  # Create a file backup with a note
  bkpdir backup myfile.txt -n "Before changes"

  # Show what would be backed up without creating backup
  bkpdir backup -d myfile.txt`,
		Args: cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ctx := context.Background()
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			filePath := args[0]

			// Use note from flag if provided, otherwise use positional argument
			backupNote := note
			if backupNote == "" && len(args) > 1 {
				backupNote = args[1]
			}

			if err := CreateFileBackupEnhanced(BackupOptions{
				Context:   ctx,
				Config:    cfg,
				Formatter: formatter,
				FilePath:  filePath,
				Note:      backupNote,
				DryRun:    dryRun,
			}); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the backup name")
	return cmd
}

func handleConfigSetCommand(key, value string) {
	// Requirement: Simple configuration modification functionality
	// Note: This is a basic implementation that modifies the local .bkpdir.yml file

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(cwd, ".bkpdir.yml")
	configData := loadExistingConfigData(configPath)
	convertedValue := convertConfigValue(key, value)
	updateConfigData(configData, key, convertedValue)
	saveConfigData(configPath, configData)

	fmt.Printf("Configuration updated: %s = %v\n", key, convertedValue)
	fmt.Printf("Config file: %s\n", configPath)
}

func loadExistingConfigData(configPath string) map[string]interface{} {
	var configData map[string]interface{}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &configData); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing existing config file: %v\n", err)
			os.Exit(1)
		}
	} else {
		configData = make(map[string]interface{})
	}

	return configData
}

func convertConfigValue(key, value string) interface{} {
	switch key {
	case "use_current_dir_name", "use_current_dir_name_for_files", "include_git_info", "verify_on_create":
		return convertBooleanValue(key, value)
	case "status_config_error", "status_created_archive", "status_created_backup",
		"status_disk_full", "status_permission_denied":
		return convertIntegerValue(key, value)
	case "archive_dir_path", "backup_dir_path", "checksum_algorithm":
		return value
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown configuration key: %s\n", key)
		fmt.Fprintf(os.Stderr, "Valid keys: archive_dir_path, backup_dir_path, use_current_dir_name, "+
			"use_current_dir_name_for_files, include_git_info, verify_on_create, checksum_algorithm, "+
			"status_config_error, status_created_archive, status_created_backup, status_disk_full, "+
			"status_permission_denied\n")
		os.Exit(1)
		return nil
	}
}

func convertBooleanValue(key, value string) bool {
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}
	fmt.Fprintf(os.Stderr, "Error: %s requires a boolean value (true/false), got: %s\n", key, value)
	os.Exit(1)
	return false
}

func convertIntegerValue(key, value string) int {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s requires an integer value, got: %s\n", key, value)
		os.Exit(1)
	}
	return intVal
}

func updateConfigData(configData map[string]interface{}, key string, convertedValue interface{}) {
	if key == "verify_on_create" || key == "checksum_algorithm" {
		if configData["verification"] == nil {
			configData["verification"] = make(map[string]interface{})
		}
		verificationMap := configData["verification"].(map[string]interface{})
		if key == "verify_on_create" {
			verificationMap["verify_on_create"] = convertedValue
		} else {
			verificationMap["checksum_algorithm"] = convertedValue
		}
	} else {
		configData[key] = convertedValue
	}
}

func saveConfigData(configPath string, configData map[string]interface{}) {
	yamlData, err := yaml.Marshal(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling config data: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
		os.Exit(1)
	}
}
