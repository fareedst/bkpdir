// This file is part of bkpdir
//
// BkpDir provides directory archiving and file backup functionality with Git integration.
// It supports creating ZIP archives of directories and individual file backups with timestamps,
// verification, and optional Git branch/hash information in filenames.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"bkpdir/pkg/formatter"
)

// 🔶 REFACTOR-001: CLI orchestration interface contracts defined - 📝
// 🔶 REFACTOR-001: Central dependency aggregation point identified - 📝

// Version, date, and commit are set at build time via ldflags
var version = "dev"
var date = "unknown"
var commit = "unknown"

// 🔺 CFG-003: Global variables for command configuration - 📝
var (
	createNote   string
	createDryRun bool
	createVerify bool
	listFile     string
	archiveName  string
	withChecksum bool
)

// Short description for the main application
var shortDesc = `bkpdir is a comprehensive backup and archiving solution for directories and files.`

// Long description for the main application
var longDesc = `bkpdir is a comprehensive backup and archiving solution for directories and files.
It supports full and incremental directory backups, individual file backups, customizable exclusion patterns, 
Git-aware archive naming, and archive verification.`

// Version information
const (
	// 🔶 REFACTOR-005: Structure optimization - Standardized version constants - 📝
	appVersion     = "1.5.0"
	appDescription = "Directory archiving and file backup tool with Git integration"
)

// Runtime compilation information (set by build flags)
var (
	// 🔶 REFACTOR-005: Structure optimization - Standardized build variables - 📝
	compileDate = "unknown"
	platform    = "unknown"
	Version     = appVersion // Public version for external access
)

var (
	dryRun     bool
	note       string
	showConfig bool
)

// ⭐ CLI-015: Path type detection for automatic command routing - 🔍
// isFile checks if the given path is a regular file
func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// ⭐ CLI-015: Path type detection for automatic command routing - 🔍
// isDirectory checks if the given path is a directory
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ⭐ CLI-015: Path validation for automatic command routing - 🛡️
// validatePath checks if the path exists and returns appropriate error message
func validatePath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied accessing path: %s", path)
		}
		return fmt.Errorf("error accessing path %s: %v", path, err)
	}
	return nil
}

// ⭐ CLI-015: Automatic command routing based on path type - 🔧
// handleAutoDetectedCommand routes to appropriate command based on first argument type
func handleAutoDetectedCommand(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no path provided\n")
		os.Exit(1)
	}

	path := args[0]

	// Validate path exists and is accessible
	if err := validatePath(path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Determine operation type based on path type
	if isFile(path) {
		// Route to file backup operation
		handleAutoDetectedFileBackup(args)
	} else if isDirectory(path) {
		// Route to directory archive operation
		handleAutoDetectedDirectoryArchive(args)
	} else {
		// Handle special file types (symlinks, devices, etc.)
		fmt.Fprintf(os.Stderr, "Error: unsupported file type for path: %s\n", path)
		fmt.Fprintf(os.Stderr, "Supported types: regular files and directories\n")
		os.Exit(1)
	}
}

// ⭐ CLI-015: Auto-detected file backup operation - 📝
// handleAutoDetectedFileBackup handles file backup when auto-detected
func handleAutoDetectedFileBackup(args []string) {
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

	// Extract note from second argument if provided
	backupNote := note // Use global note flag if set
	if backupNote == "" && len(args) > 1 {
		backupNote = args[1]
	}

	// Create file backup using existing functionality
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
}

// ⭐ CLI-015: Auto-detected directory archive operation - 📝
// handleAutoDetectedDirectoryArchive handles directory archive when auto-detected
func handleAutoDetectedDirectoryArchive(args []string) {
	ctx := context.Background()

	// Change to the specified directory for archiving
	dirPath := args[0]
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.Chdir(dirPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error changing to directory %s: %v\n", dirPath, err)
		os.Exit(1)
	}

	// Restore original directory on exit
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to restore original directory: %v\n", err)
		}
	}()

	cfg, err := LoadConfig(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1) // Use default exit code since cfg might be nil
	}

	formatter := NewOutputFormatter(cfg)

	// Extract note from second argument if provided
	archiveNote := note // Use global note flag if set
	if archiveNote == "" && len(args) > 1 {
		archiveNote = args[1]
	}

	// Create full archive using existing functionality
	if err := CreateFullArchiveWithContext(ctx, cfg, archiveNote, dryRun, false); err != nil {
		exitCode := HandleArchiveError(err, cfg, formatter)
		os.Exit(exitCode)
	}
}

// Long description for root command
const rootLongDesc = `bkpdir version %s (compiled %s) [%s]

BkpDir is a command-line tool for archiving directories and backing up individual files on macOS and Linux. 
It supports full and incremental directory backups, individual file backups, customizable exclusion patterns, 
Git-aware archive naming, and archive verification.`

// ⭐ CLI-015: Custom command execution with auto-detection fallback - 🔧
// executeWithAutoDetection handles Cobra command resolution issues by implementing
// custom argument parsing that allows auto-detection to work alongside existing commands
func executeWithAutoDetection(rootCmd *cobra.Command) error {
	args := os.Args[1:] // Skip program name

	// If no arguments, execute normally (will show help)
	if len(args) == 0 {
		return rootCmd.Execute()
	}

	// Check if first argument is a known command
	firstArg := args[0]

	// List of known commands that should be handled by Cobra normally
	knownCommands := []string{
		"create", "config", "template", "full", "inc", "list", "verify", "backup", "version",
		"help", "--help", "-h", "--version", "-v",
	}

	// Check for global flags that should be handled normally
	globalFlags := []string{
		"--config", "--dry-run", "-d", "--list",
	}

	// If first argument is a known command or global flag, execute normally
	for _, cmd := range knownCommands {
		if firstArg == cmd {
			return rootCmd.Execute()
		}
	}

	for _, flag := range globalFlags {
		if firstArg == flag {
			return rootCmd.Execute()
		}
	}

	// Check if first argument starts with a flag (-)
	if strings.HasPrefix(firstArg, "-") {
		return rootCmd.Execute()
	}

	// At this point, we assume it's a path for auto-detection
	// Let handleAutoDetectedCommand handle path validation and provide appropriate errors
	// We need to manually handle the global flags that might be present
	var filteredArgs []string
	var dryRunFlag bool
	var noteFlag string

	// Parse arguments to extract global flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--dry-run" || arg == "-d" {
			dryRunFlag = true
		} else if arg == "--note" || arg == "-n" {
			if i+1 < len(args) {
				noteFlag = args[i+1]
				i++ // Skip the next argument as it's the note value
			}
		} else if strings.HasPrefix(arg, "--note=") {
			noteFlag = strings.TrimPrefix(arg, "--note=")
		} else if strings.HasPrefix(arg, "-n=") {
			noteFlag = strings.TrimPrefix(arg, "-n=")
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
		i++
	}

	// Set global variables for the handlers to use
	dryRun = dryRunFlag
	if noteFlag != "" {
		note = noteFlag
	}

	// Execute auto-detection with filtered arguments
	handleAutoDetectedCommand(filteredArgs)
	return nil
}

func main() {
	// 🔺 CFG-001: CLI application initialization and command structure - 📝
	// DECISION-REF: DEC-002
	rootCmd := &cobra.Command{
		Use:     "bkpdir",
		Short:   "Directory archiving and file backup CLI for macOS and Linux",
		Long:    fmt.Sprintf(rootLongDesc, Version, compileDate, platform),
		Version: fmt.Sprintf("%s (compiled %s) [%s]", Version, compileDate, platform),
		// ⭐ CLI-015: Disable default command suggestions to enable auto-detection - 🔧
		DisableSuggestions: true,
		Example: `  # Auto-detect operation based on path type (NEW)
  bkpdir myfile.txt "Before changes"        # Creates file backup automatically
  bkpdir mydirectory "Initial backup"       # Creates directory archive automatically
  bkpdir /path/to/file.txt                  # Auto-detects file backup
  bkpdir /path/to/directory                 # Auto-detects directory archive

  # Create a full directory archive (explicit commands)
  bkpdir create "Initial backup"
  bkpdir full -n "Initial backup"  # backward compatibility

  # Create an incremental directory archive with verification
  bkpdir create --incremental "Changes after feature X" -v
  bkpdir inc -n "Changes after feature X" -v  # backward compatibility

  # Create a file backup (explicit command)
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

			// ⭐ CLI-015: Automatic file/directory command detection - 🔧
			// If positional arguments provided, auto-detect operation type
			if len(args) > 0 {
				handleAutoDetectedCommand(args)
				return
			}

			// If no config flag and no subcommand, show help
			cmd.Help()
		},
	}

	// Set the version template to show version in help output
	versionTemplate := fmt.Sprintf("bkpdir version %s (compiled %s) [%s]\n",
		Version, compileDate, platform)
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
	rootCmd.AddCommand(templateCmd())

	// Add backward compatibility commands
	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())

	// Add other commands
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(backupCmd())
	rootCmd.AddCommand(versionCmd())

	// ⭐ CLI-015: Custom command execution with auto-detection fallback - 🔧
	if err := executeWithAutoDetection(rootCmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleConfigCommand() {
	// 🔺 CFG-001: Configuration display handling - 🔍
	// 🔺 CFG-003: Configuration output formatting - 🔍
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

// 🔺 CFG-006: Enhanced config command interface implementation - 🔧
// IMPLEMENTATION-REF: CFG-006 Subtask 4: Enhanced config command interface
// handleEnhancedConfigCommand provides comprehensive configuration inspection with filtering options.
func handleEnhancedConfigCommand(showAll, showOverrides, showSources bool, outputFormat, filterPattern string) {
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

	// Get enhanced configuration values with metadata
	configValues := GetAllConfigValuesWithSources(cfg, cwd)

	// Apply filtering
	filteredValues := applyConfigFiltering(configValues, showOverrides, filterPattern)

	// Display based on format
	switch outputFormat {
	case "tree":
		displayConfigTree(filteredValues, showSources)
	case "json":
		displayConfigJSON(filteredValues, showSources)
	default: // "table"
		displayConfigTable(filteredValues, showSources)
	}
}

// 🔺 CFG-006: Configuration filtering implementation - 🔧
// IMPLEMENTATION-REF: CFG-006 Subtask 5: Command-line options and filtering
// applyConfigFiltering filters configuration values based on user criteria.
func applyConfigFiltering(values []ConfigValueWithMetadata, showOverridesOnly bool, filterPattern string) []ConfigValueWithMetadata {
	var filtered []ConfigValueWithMetadata

	for _, value := range values {
		// Apply overrides-only filter
		if showOverridesOnly && !value.IsOverridden {
			continue
		}

		// Apply pattern filter
		if filterPattern != "" {
			if !strings.Contains(strings.ToLower(value.ConfigValue.Name), strings.ToLower(filterPattern)) &&
				!strings.Contains(strings.ToLower(value.FieldInfo.Category), strings.ToLower(filterPattern)) {
				continue
			}
		}

		filtered = append(filtered, value)
	}

	return filtered
}

// 🔺 CFG-006: Hierarchical value resolution display - 🔧
// IMPLEMENTATION-REF: CFG-006 Subtask 3: Hierarchical value resolution display
// displayConfigTree shows configuration in tree format with inheritance chains.
func displayConfigTree(values []ConfigValueWithMetadata, showSources bool) {
	// Group by category for tree display
	categories := make(map[string][]ConfigValueWithMetadata)
	for _, value := range values {
		category := value.FieldInfo.Category
		categories[category] = append(categories[category], value)
	}

	// Display each category as a tree branch
	for category, categoryValues := range categories {
		fmt.Printf("📁 %s\n", strings.Title(strings.ReplaceAll(category, "_", " ")))

		for i, value := range categoryValues {
			prefix := "├── "
			if i == len(categoryValues)-1 {
				prefix = "└── "
			}

			fmt.Printf("%s%s: %s", prefix, value.ConfigValue.Name, value.ConfigValue.Value)

			if showSources {
				fmt.Printf(" (source: %s)", value.ConfigValue.Source)
				if len(value.InheritanceChain) > 1 {
					fmt.Printf(" [chain: %s]", strings.Join(value.InheritanceChain, " → "))
				}
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

// 🔺 CFG-006: JSON output format implementation - 🔧
// IMPLEMENTATION-REF: CFG-006 Subtask 4: Multiple display formats
// displayConfigJSON outputs configuration in JSON format for programmatic use.
func displayConfigJSON(values []ConfigValueWithMetadata, showSources bool) {
	type ConfigOutput struct {
		Name             string   `json:"name"`
		Value            string   `json:"value"`
		Source           string   `json:"source"`
		Category         string   `json:"category"`
		Type             string   `json:"type"`
		IsOverridden     bool     `json:"is_overridden"`
		InheritanceChain []string `json:"inheritance_chain,omitempty"`
		MergeStrategy    string   `json:"merge_strategy,omitempty"`
	}

	var output []ConfigOutput
	for _, value := range values {
		configOutput := ConfigOutput{
			Name:         value.ConfigValue.Name,
			Value:        value.ConfigValue.Value,
			Source:       value.ConfigValue.Source,
			Category:     value.FieldInfo.Category,
			Type:         value.FieldInfo.Type,
			IsOverridden: value.IsOverridden,
		}

		if showSources {
			configOutput.InheritanceChain = value.InheritanceChain
			configOutput.MergeStrategy = value.MergeStrategy
		}

		output = append(output, configOutput)
	}

	// Pretty print JSON
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

// 🔺 CFG-006: Table display format implementation - 🔧
// IMPLEMENTATION-REF: CFG-006 Subtask 4: Enhanced display with source attribution
// displayConfigTable shows configuration in traditional table format with enhanced metadata.
func displayConfigTable(values []ConfigValueWithMetadata, showSources bool) {
	// Display configuration file paths first
	cwd, _ := os.Getwd()
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
	fmt.Printf("config: %s (source: default)\n", configPathsStr)

	// Display each configuration value
	for _, value := range values {
		fmt.Printf("%s: %s (source: %s)", value.ConfigValue.Name, value.ConfigValue.Value, value.ConfigValue.Source)

		if showSources {
			if value.IsOverridden {
				fmt.Printf(" [overridden]")
			}
			if value.FieldInfo.Category != "basic_settings" {
				fmt.Printf(" [%s]", value.FieldInfo.Category)
			}
			if len(value.InheritanceChain) > 1 {
				fmt.Printf(" [chain: %s]", strings.Join(value.InheritanceChain, " → "))
			}
		}
		fmt.Println()
	}
}

func handleCreateCommand() {
	// Implementation
}

// ⭐ CFG-TEMPLATE-001: Template command implementation - 🔧
func handleTemplateCommand(cmd *cobra.Command, args []string) {
	// Get flag values
	outputFile, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	force, _ := cmd.Flags().GetBool("force")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// ⭐ CFG-TEMPLATE-001: Configuration reflection - 🔧
	// Load current configuration to populate template values
	cfg, err := LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// ⭐ CFG-TEMPLATE-001: File management - 🔧
	// Determine output filename
	targetFile := determineTemplateFileName(outputFile)

	// Check if file exists and handle conflicts
	if !force && !dryRun {
		if _, err := os.Stat(targetFile); err == nil {
			fmt.Printf("File %s already exists. Use --force to overwrite or choose a different name.\n", targetFile)
			os.Exit(1)
		}
	}

	// ⭐ CFG-TEMPLATE-001: Template generation - 🔧
	// Generate template content
	templateContent, err := generateConfigurationTemplate(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating template: %v\n", err)
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf("Would create file: %s\n", targetFile)
		fmt.Printf("Template content:\n")
		fmt.Printf("=================\n")
		fmt.Print(templateContent)
		return
	}

	// Write template to file
	err = writeTemplateToFile(targetFile, templateContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing template file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Configuration template created: %s\n", targetFile)
	fmt.Printf("📝 Edit the file to customize your configuration options\n")
}

func handleListCommand() {
	// ⭐ ARCH-002: Archive listing command implementation - 📝
	// 🔺 CFG-003: Archive listing output formatting - 📝
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
	// 🔺 CFG-001: Configuration command implementation - 📝
	// 🔺 CFG-003: Configuration command interface - 📝
	// 🔺 CFG-006: Enhanced config command interface - 🔧
	// 🔻 CFG-006: Documentation - 📝 Enhanced help text

	var (
		showAll       bool
		showOverrides bool
		showSources   bool
		outputFormat  string
		filterPattern string
	)

	cmd := &cobra.Command{
		Use:   "config [KEY] [VALUE]",
		Short: "Display or modify configuration values",
		Long: `Display configuration values or set a specific configuration value.

The config command supports comprehensive configuration inspection with automatic field discovery.
All configuration parameters (100+) are visible without manual maintenance using Go reflection.

Key Features:
  • Automatic field discovery - Zero maintenance, new fields appear automatically
  • Source attribution - Shows complete inheritance chain (environment → config → defaults)
  • Performance optimized - Sub-100ms response with reflection caching
  • Multiple output formats - Table, tree, and JSON for different use cases
  • Advanced filtering - Pattern-based field filtering and override detection

Display Options:
  --all             Show all configuration fields (default behavior)
  --overrides-only  Display only non-default values  
  --sources         Show detailed source attribution with inheritance chains
  --format FORMAT   Choose output format: table (default), tree, json
  --filter PATTERN  Filter fields by name pattern

Examples:
  # Display all configuration values (100+ fields auto-discovered)
  bkpdir config

  # Show only customized values with sources
  bkpdir config --overrides-only --sources

  # Debug inheritance chain for specific field
  bkpdir config exclude_patterns --sources --format tree

  # Display in hierarchical tree format
  bkpdir config --format tree

  # Filter for archive-related settings
  bkpdir config --filter "archive" --sources

  # Export configuration for automation
  bkpdir config --format json | jq '.configuration[]'

  # Set a configuration value
  bkpdir config archive_dir_path /custom/archive/path
  bkpdir config include_git_info false

Troubleshooting:
  # Check why a value isn't being applied
  bkpdir config [field_name] --sources --format tree
  
  # Find all environment variable overrides
  bkpdir config --sources | grep Environment
  
  # Performance check (should be <100ms)
  time bkpdir config >/dev/null

For detailed documentation, see docs/configuration-inspection-guide.md`,
		Args: cobra.MaximumNArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			if len(args) == 0 {
				// Enhanced configuration display with filtering options
				handleEnhancedConfigCommand(showAll, showOverrides, showSources, outputFormat, filterPattern)
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

	// 🔺 CFG-006: Command-line options and filtering - 🔧
	cmd.Flags().BoolVar(&showAll, "all", true, "Show all configuration fields")
	cmd.Flags().BoolVar(&showOverrides, "overrides-only", false, "Display only non-default values")
	cmd.Flags().BoolVar(&showSources, "sources", false, "Show detailed source attribution")
	cmd.Flags().StringVar(&outputFormat, "format", "table", "Output format: table, tree, json")
	cmd.Flags().StringVar(&filterPattern, "filter", "", "Filter fields by name pattern")

	return cmd
}

func createCmd() *cobra.Command {
	// ⭐ ARCH-002: Archive creation command implementation - 🔧
	// 🔺 CFG-003: Command interface for archive creation - 🔧
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new archive",
		Run: func(*cobra.Command, []string) {
			handleCreateCommand()
		},
	}
	return cmd
}

func templateCmd() *cobra.Command {
	// ⭐ CFG-TEMPLATE-001: CLI command implementation - 🔧
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Generate configuration template file",
		Long: `Generate a comprehensive configuration template with all available options.

Creates a YAML configuration file with all possible configuration keys, properly 
commented and organized by category. Values are populated from the current 
configuration after loading BKPDIR_CONFIG.

File naming:
- Creates .bkpdir.yml if it doesn't exist
- Creates .bkpdir.default-YYYY-MM-DD.yml if .bkpdir.yml already exists`,
		Example: `  # Generate template in current directory
  bkpdir template

  # Generate template with custom filename
  bkpdir template --output custom-config.yml

  # Preview template without creating file
  bkpdir template --dry-run`,
		Run: func(cmd *cobra.Command, args []string) {
			handleTemplateCommand(cmd, args)
		},
	}

	// Add flags
	cmd.Flags().StringP("output", "o", "", "Custom output filename (default: .bkpdir.yml or .bkpdir.default-YYYY-MM-DD.yml)")
	cmd.Flags().BoolP("dry-run", "d", false, "Show what would be written without creating the file")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing files without confirmation")

	return cmd
}

func fullCmd() *cobra.Command {
	// ⭐ ARCH-002: Full archive creation command (backward compatibility) - 🔧
	// 🔺 CFG-003: Backward compatibility command interface - 🔧
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
	// ⭐ ARCH-003: Incremental archive creation command - 🔧
	// 🔺 CFG-003: Incremental command interface - 🔧
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
	// ⭐ ARCH-002: Archive listing command - 🔧
	// 🔺 CFG-003: List command interface - 🔧
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
	// Archive verification command
	// 🔺 CFG-003: Verify command interface - 🛡️
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
	// Version display command
	// 🔺 CFG-003: Version command interface - 📝
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
	Formatter formatter.OutputFormatterInterface
	Note      string
	DryRun    bool
	Verify    bool
}

// CreateFullArchiveEnhanced creates a full archive of the current directory with enhanced error handling
// and resource management. It supports dry-run mode and optional verification.
func CreateFullArchiveEnhanced(opts ArchiveOptions) error {
	// ⭐ ARCH-002: Enhanced full archive creation - 🔧
	// DECISION-REF: DEC-006, DEC-007
	return CreateFullArchive(opts.Config, opts.Note, opts.DryRun, opts.Verify)
}

// CreateIncrementalArchiveEnhanced creates an incremental archive containing only files changed since
// the last full archive. It supports dry-run mode and optional verification.
func CreateIncrementalArchiveEnhanced(opts ArchiveOptions) error {
	// ⭐ ARCH-003: Enhanced incremental archive creation - 📝
	// DECISION-REF: DEC-006, DEC-007
	return CreateIncrementalArchive(opts.Config, opts.Note, opts.DryRun, opts.Verify)
}

// ListArchivesEnhanced displays all archives in the archive directory with enhanced formatting
// and error handling.
func ListArchivesEnhanced(cfg *Config, formatter formatter.OutputFormatterInterface) error {
	// ⭐ ARCH-002: Enhanced archive listing with formatting - 🔍
	// 🔺 CFG-003: Template-based archive listing - 🔍
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
		// Cast to FormatterAdapter to access extended methods
		if formatterAdapter, ok := formatter.(*FormatterAdapter); ok {
			formatterAdapter.PrintNoArchivesFound(archiveDir)
		} else {
			formatter.PrintError(fmt.Sprintf("No archives found in %s", archiveDir))
		}
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
		if formatterAdapter, ok := formatter.(*FormatterAdapter); ok {
			output := formatterAdapter.FormatListArchiveWithExtraction(a.Name, creationTime)
			// Remove trailing newline from output to add status on same line
			output = strings.TrimSuffix(output, "\n")
			formatterAdapter.PrintArchiveListWithStatus(output, status)
		} else {
			output := formatter.FormatListArchive(a.Name, creationTime)
			fmt.Printf("%s%s\n", strings.TrimSuffix(output, "\n"), status)
		}
	}

	return nil
}

// VerifyOptions holds parameters for archive verification functions
type VerifyOptions struct {
	Config       *Config
	Formatter    formatter.OutputFormatterInterface
	ArchiveName  string
	WithChecksum bool
}

// VerifyArchiveEnhanced verifies the integrity of an archive with optional checksum verification.
// It provides enhanced error handling and reporting.
func VerifyArchiveEnhanced(opts VerifyOptions) error {
	// Archive verification implementation
	// 🔺 CFG-003: Verification output formatting - 🔍
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
	// 🔺 CFG-001: Archive directory resolution - 🔍
	// DECISION-REF: DEC-002
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
	// Single archive verification
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
	// All archives verification
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to list archives", 1, err)
	}

	allPassed := true
	for _, archive := range archives {
		status, err := performVerification(archive.Path, opts.WithChecksum)
		if err != nil {
			// Cast to FormatterAdapter to access extended methods
			if formatterAdapter, ok := opts.Formatter.(*FormatterAdapter); ok {
				formatterAdapter.PrintVerificationFailed(archive.Name, err)
			} else {
				opts.Formatter.PrintError(fmt.Sprintf("Verification failed for %s: %v", archive.Name, err))
			}
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
	// Archive verification execution
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
	// Get config and formatter for proper output formatting
	cwd, _ := os.Getwd()
	cfg, _ := LoadConfig(cwd)
	formatter := NewOutputFormatter(cfg)

	// Store verification status
	if err := StoreVerificationStatus(archive, status); err != nil {
		// Don't fail if we can't store status, just warn
		formatter.PrintVerificationWarning(name, err)
	}

	if status.IsVerified {
		formatter.PrintVerificationSuccess(name)
		return nil
	}

	formatter.PrintVerificationFailed(name, fmt.Errorf("verification failed"))
	for _, errMsg := range status.Errors {
		formatter.PrintVerificationErrorDetail(errMsg)
	}
	return NewArchiveError("Archive verification failed", 1)
}

func handleListFileBackupsCommand(args []string) {
	// ⭐ FILE-002: File backup listing command implementation - 📝
	// 🔺 CFG-003: File backup listing output formatting - 📝
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
	// ⭐ FILE-002: File backup command implementation - 🔧
	// 🔺 CFG-003: Backup command interface - 🔧
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
	// 🔺 CFG-001: Configuration modification command - 🔍
	// 🔺 CFG-002: Configuration value setting - 🔍
	// DECISION-REF: DEC-002
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
	configPath := filepath.Join(cwd, ".bkpdir.yml")
	configData := loadExistingConfigData(configPath)
	convertedValue := convertConfigValue(key, value)
	updateConfigData(configData, key, convertedValue)
	saveConfigData(configPath, configData)

	formatter.PrintConfigurationUpdated(key, convertedValue)
	formatter.PrintConfigFilePath(configPath)
}

func loadExistingConfigData(configPath string) map[string]interface{} {
	// 🔺 CFG-001: Configuration file loading - 📝
	// DECISION-REF: DEC-002
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
	// 🔺 CFG-002: Configuration value type conversion - 🔧
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
	// 🔺 CFG-002: Boolean configuration value conversion - 🔧
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
	// 🔺 CFG-002: Integer configuration value conversion - 📝
	intVal, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s requires an integer value, got: %s\n", key, value)
		os.Exit(1)
	}
	return intVal
}

func updateConfigData(configData map[string]interface{}, key string, convertedValue interface{}) {
	// 🔺 CFG-001: Configuration data updating - 🔍
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
	// 🔺 CFG-001: Configuration data persistence - 📝
	// DECISION-REF: DEC-002, DEC-008
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

// 🔶 REFACTOR-005: Structure optimization - Standardized command configuration - 📝
// CommandConfig holds configuration for CLI command execution
type CommandConfig struct {
	Config    *Config
	Formatter formatter.OutputFormatterInterface
	Context   context.Context
}

// 🔶 REFACTOR-005: Structure optimization - Interface-based command handler abstraction - 📝
// CommandHandlerInterface abstracts command execution for better testability
type CommandHandlerInterface interface {
	HandleFullArchive(args []string, note string, dryRun bool, verify bool) error
	HandleIncrementalArchive(args []string, note string, dryRun bool, verify bool) error
	HandleListArchives(args []string) error
	HandleVerifyArchive(args []string, archiveName string, withChecksum bool) error
	HandleFileBackup(args []string, note string, dryRun bool) error
	HandleListFileBackups(args []string, filePath string) error
	HandleDisplayConfig(args []string) error
}

// 🔶 REFACTOR-005: Structure optimization - Concrete command handler implementation - 📝
// CommandHandler provides the concrete implementation of CommandHandlerInterface
type CommandHandler struct {
	config CommandConfig
}

// 🔶 REFACTOR-005: Extraction preparation - Interface-based command handler factory - 🔧
// NewCommandHandler creates a new CommandHandler with the given configuration
func NewCommandHandler(cfg CommandConfig) CommandHandlerInterface {
	return &CommandHandler{config: cfg}
}

// 🔶 REFACTOR-005: Structure optimization - Interface method implementations - 🔧
// HandleFullArchive handles full archive creation command
func (h *CommandHandler) HandleFullArchive(args []string, note string, dryRun bool, verify bool) error {
	return CreateFullArchiveWithContext(h.config.Context, h.config.Config, note, dryRun, verify)
}

// HandleIncrementalArchive handles incremental archive creation command
func (h *CommandHandler) HandleIncrementalArchive(args []string, note string, dryRun bool, verify bool) error {
	return CreateIncrementalArchiveWithContext(h.config.Context, h.config.Config, note, dryRun, verify)
}

// HandleListArchives handles archive listing command
func (h *CommandHandler) HandleListArchives(args []string) error {
	archiveDir, err := getArchiveDirectory(h.config.Config)
	if err != nil {
		return err
	}

	archives, err := ListArchives(archiveDir)
	if err != nil {
		return err
	}

	if len(archives) == 0 {
		// Cast to FormatterAdapter to access extended methods
		if formatterAdapter, ok := h.config.Formatter.(*FormatterAdapter); ok {
			formatterAdapter.PrintNoArchivesFound(archiveDir)
		} else {
			h.config.Formatter.PrintError(fmt.Sprintf("No archives found in %s", archiveDir))
		}
		return nil
	}

	for _, archive := range archives {
		creationTime := archive.CreationTime.Format("2006-01-02 15:04:05")
		if formatterAdapter, ok := h.config.Formatter.(*FormatterAdapter); ok {
			output := formatterAdapter.FormatListArchiveWithExtraction(archive.Name, creationTime)
			fmt.Print(output)
		} else {
			output := h.config.Formatter.FormatListArchive(archive.Name, creationTime)
			fmt.Print(output)
		}
	}

	return nil
}

// HandleVerifyArchive handles archive verification command
func (h *CommandHandler) HandleVerifyArchive(args []string, archiveName string, withChecksum bool) error {
	opts := VerifyOptions{
		Config:       h.config.Config,
		Formatter:    h.config.Formatter,
		ArchiveName:  archiveName,
		WithChecksum: withChecksum,
	}
	return VerifyArchiveEnhanced(opts)
}

// HandleFileBackup handles file backup creation command
func (h *CommandHandler) HandleFileBackup(args []string, note string, dryRun bool) error {
	if len(args) == 0 {
		return fmt.Errorf("file path required")
	}
	filePath := args[0]

	opts := BackupOptions{
		Context:   h.config.Context,
		Config:    h.config.Config,
		Formatter: h.config.Formatter,
		FilePath:  filePath,
		Note:      note,
		DryRun:    dryRun,
	}
	return CreateFileBackupEnhanced(opts)
}

// HandleListFileBackups handles file backup listing command
func (h *CommandHandler) HandleListFileBackups(args []string, filePath string) error {
	targetFile := filePath
	if targetFile == "" && len(args) > 0 {
		targetFile = args[0]
	}
	if targetFile == "" {
		return fmt.Errorf("file path required for listing backups")
	}

	return ListFileBackupsEnhanced(h.config.Config, h.config.Formatter, targetFile)
}

// HandleDisplayConfig handles configuration display command
func (h *CommandHandler) HandleDisplayConfig(args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	configValues := GetConfigValuesWithSources(h.config.Config, cwd)
	for _, cv := range configValues {
		h.config.Formatter.PrintConfigValue(cv.Name, cv.Value, cv.Source)
	}
	return nil
}

// ⭐ CFG-TEMPLATE-001: File management - 🔧
// determineTemplateFileName decides the target filename for the template based on user input and existing files
func determineTemplateFileName(customOutput string) string {
	if customOutput != "" {
		return customOutput
	}

	// Check if .bkpdir.yml exists
	defaultFile := ".bkpdir.yml"
	if _, err := os.Stat(defaultFile); os.IsNotExist(err) {
		return defaultFile
	}

	// Generate date-based filename if default exists
	currentTime := time.Now()
	dateBasedFile := fmt.Sprintf(".bkpdir.default-%s.yml", currentTime.Format("2006-01-02"))
	return dateBasedFile
}

// ⭐ CFG-TEMPLATE-001: Template generation - 🔧
// generateConfigurationTemplate creates a comprehensive YAML template with all Config fields
func generateConfigurationTemplate(cfg *Config) (string, error) {
	var template strings.Builder

	// Add header
	template.WriteString("# BkpDir Configuration Template\n")
	template.WriteString("# Generated on: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	template.WriteString("#\n")
	template.WriteString("# This template includes all available configuration options organized by category.\n")
	template.WriteString("# Uncomment and modify the values you want to customize.\n")
	template.WriteString("# Values shown are the current effective configuration after loading BKPDIR_CONFIG.\n")
	template.WriteString("#\n")
	template.WriteString("# For more information, see: docs/configuration.md\n")
	template.WriteString("\n")

	// ⭐ CFG-TEMPLATE-001: Configuration reflection - 🔧
	// Use existing CFG-006 field discovery system
	allFields := GetAllConfigFields(cfg)

	// Group fields by category
	categories := make(map[string][]configFieldInfo)
	for _, field := range allFields {
		category := field.Category
		if category == "" {
			category = "other"
		}
		categories[category] = append(categories[category], field)
	}

	// Define category order for better organization
	categoryOrder := []string{
		"basic_settings", "archive_settings", "backup_settings", "verification", "inheritance",
		"status_codes", "format_strings", "template_strings", "regex_patterns",
	}

	// Generate template sections by category
	for _, categoryName := range categoryOrder {
		fields, exists := categories[categoryName]
		if !exists || len(fields) == 0 {
			continue
		}

		// Category header
		categoryTitle := strings.Title(strings.ReplaceAll(categoryName, "_", " "))
		template.WriteString(fmt.Sprintf("# %s Configuration\n", categoryTitle))
		template.WriteString(strings.Repeat("#", len(categoryTitle)+16) + "\n")

		// Add category description
		addCategoryDescription(&template, categoryName)
		template.WriteString("\n")

		// Generate fields for this category
		for _, field := range fields {
			generateFieldTemplate(&template, field, cfg)
		}

		template.WriteString("\n")
	}

	return template.String(), nil
}

// ⭐ CFG-TEMPLATE-001: Template generation - 🔧
// addCategoryDescription adds helpful descriptions for each configuration category
func addCategoryDescription(template *strings.Builder, categoryName string) {
	descriptions := map[string]string{
		"basic_settings":   "# Basic backup and archive settings including paths, patterns, and file handling",
		"archive_settings": "# Directory archive specific settings for archive creation and management",
		"backup_settings":  "# File backup specific settings for individual file backup operations",
		"verification":     "# Archive verification settings including checksum algorithms and validation",
		"inheritance":      "# Configuration inheritance settings for loading configuration from multiple files",
		"status_codes":     "# Exit status codes for various operation outcomes and error conditions",
		"format_strings":   "# Printf-style format strings for output messages and display formatting",
		"template_strings": "# Template-based format strings with placeholders for dynamic content",
		"regex_patterns":   "# Regular expression patterns for filename parsing and data extraction",
	}

	if desc, exists := descriptions[categoryName]; exists {
		template.WriteString(desc + "\n")
	}
}

// ⭐ CFG-TEMPLATE-001: Template generation - 🔧
// generateFieldTemplate creates a commented YAML entry for a single configuration field
func generateFieldTemplate(template *strings.Builder, field configFieldInfo, cfg *Config) {
	// Get current value
	currentValue := formatFieldValue(field.Value, field.Kind)

	// Format the field with proper indentation
	yamlKey := field.YAMLName

	// Handle nested fields (like git.enabled)
	if strings.Contains(field.Path, ".") {
		// For nested fields, we need special handling
		parts := strings.Split(field.Path, ".")
		if len(parts) == 2 {
			parentKey := parts[0]
			childKey := field.YAMLName

			// Check if this is the first field of a nested struct
			if field.Name == getFirstFieldOfStruct(field.Path) {
				template.WriteString(fmt.Sprintf("# %s:\n", parentKey))
			}
			template.WriteString(fmt.Sprintf("#   %s: %s\n", childKey, currentValue))
			return
		}
	}

	// Regular field
	template.WriteString(fmt.Sprintf("# %s: %s\n", yamlKey, currentValue))
}

// ⭐ CFG-TEMPLATE-001: Template generation - 🔧
// getFirstFieldOfStruct determines if this is the first field of a nested struct (for proper formatting)
func getFirstFieldOfStruct(fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	if len(parts) != 2 {
		return ""
	}

	parent := parts[0]
	switch parent {
	case "verification":
		return "VerifyOnCreate"
	case "git":
		return "Enabled"
	default:
		return ""
	}
}

// ⭐ CFG-TEMPLATE-001: File management - 🔧
// writeTemplateToFile safely writes the template content to the specified file
func writeTemplateToFile(filename, content string) error {
	// Create backup of existing file if it exists
	if _, err := os.Stat(filename); err == nil {
		backupName := filename + ".backup"
		if err := copyFile(filename, backupName); err != nil {
			return fmt.Errorf("failed to create backup: %v", err)
		}
	}

	// Write template content
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write template file: %v", err)
	}

	return nil
}
