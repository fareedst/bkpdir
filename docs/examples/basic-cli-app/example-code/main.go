// EXTRACT-010: Basic CLI Application Example - Integration of config + cli + formatter packages - 🔺
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"bkpdir/pkg/cli"
	"bkpdir/pkg/config"
	"bkpdir/pkg/errors"
	"bkpdir/pkg/formatter"
)

// BasicCliApp demonstrates integration of core extracted packages
type BasicCliApp struct {
	configLoader *config.ConfigLoader
	formatter    formatter.FormatterInterface
	appInfo      *cli.AppInfo
}

// NewBasicCliApp creates a new basic CLI application
func NewBasicCliApp() *BasicCliApp {
	return &BasicCliApp{
		configLoader: config.NewConfigLoader(),
		formatter:    formatter.NewTemplateFormatter(),
		appInfo: &cli.AppInfo{
			Name:        "basic-cli-app",
			Version:     "1.0.0",
			Description: "Example CLI application using extracted packages",
		},
	}
}

// Configuration structure for the application
type AppConfig struct {
	// Output configuration
	OutputFormat string `json:"output_format" yaml:"output_format"`
	Verbose      bool   `json:"verbose" yaml:"verbose"`
	DryRun       bool   `json:"dry_run" yaml:"dry_run"`

	// Processing configuration
	Processing struct {
		BatchSize   int           `json:"batch_size" yaml:"batch_size"`
		Timeout     time.Duration `json:"timeout" yaml:"timeout"`
		Concurrency int           `json:"concurrency" yaml:"concurrency"`
	} `json:"processing" yaml:"processing"`

	// Output configuration
	Output struct {
		Directory string `json:"directory" yaml:"directory"`
		Template  string `json:"template" yaml:"template"`
	} `json:"output" yaml:"output"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *AppConfig {
	return &AppConfig{
		OutputFormat: "table",
		Verbose:      false,
		DryRun:       false,
		Processing: struct {
			BatchSize   int           `json:"batch_size" yaml:"batch_size"`
			Timeout     time.Duration `json:"timeout" yaml:"timeout"`
			Concurrency int           `json:"concurrency" yaml:"concurrency"`
		}{
			BatchSize:   100,
			Timeout:     time.Minute * 10,
			Concurrency: 4,
		},
		Output: struct {
			Directory string `json:"directory" yaml:"directory"`
			Template  string `json:"template" yaml:"template"`
		}{
			Directory: "./output",
			Template:  "{{.Name}}: {{.Status}} ({{.Duration}})",
		},
	}
}

// buildRootCommand creates the root command with integrated packages
func (app *BasicCliApp) buildRootCommand() *cobra.Command {
	var configFile string
	var outputFormat string
	var verbose bool
	var dryRun bool

	rootCmd := &cobra.Command{
		Use:   app.appInfo.Name,
		Short: app.appInfo.Description,
		Long: fmt.Sprintf(`%s

This example demonstrates integration of extracted BkpDir packages:
- pkg/config: Configuration management from multiple sources
- pkg/cli: CLI framework with dry-run and context support  
- pkg/formatter: Output formatting with templates
- pkg/errors: Structured error handling

Example usage:
  %s process --config=config.yaml --format=json
  %s process --dry-run --verbose
  %s status --format=table`,
			app.appInfo.Description, app.appInfo.Name, app.appInfo.Name, app.appInfo.Name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Global flags using pkg/cli patterns
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Configuration file path")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "Output format (table, json, yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without executing")

	// Add commands
	rootCmd.AddCommand(app.buildProcessCommand())
	rootCmd.AddCommand(app.buildStatusCommand())
	rootCmd.AddCommand(app.buildConfigCommand())

	return rootCmd
}

// buildProcessCommand creates the process command
func (app *BasicCliApp) buildProcessCommand() *cobra.Command {
	var inputPath string
	var outputPath string

	cmd := &cobra.Command{
		Use:   "process [path]",
		Short: "Process files using integrated packages",
		Long: `Process files demonstrating package integration:

This command shows how config, cli, and formatter packages work together:
1. Load configuration from multiple sources (file, env, flags)
2. Use CLI framework for command structure and dry-run support
3. Format output using templates and different formats
4. Handle errors with structured error types`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runProcessCommand(cmd.Context(), inputPath, outputPath, args)
		},
	}

	cmd.Flags().StringVar(&inputPath, "input", ".", "Input path to process")
	cmd.Flags().StringVar(&outputPath, "output", "./output", "Output path for results")

	return cmd
}

// buildStatusCommand creates the status command
func (app *BasicCliApp) buildStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show application status with formatted output",
		Long: `Display application status demonstrating formatter integration:

This command shows different output formats:
- Table format for human-readable output
- JSON format for machine processing
- YAML format for configuration-like output
- Custom templates for specialized formatting`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runStatusCommand(cmd.Context())
		},
	}

	return cmd
}

// buildConfigCommand creates the config command
func (app *BasicCliApp) buildConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management commands",
		Long:  `Commands for managing application configuration using pkg/config`,
	}

	// Add subcommands
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runConfigShowCommand(cmd.Context())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "validate",
		Short: "Validate configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.runConfigValidateCommand(cmd.Context())
		},
	})

	return cmd
}

// runProcessCommand demonstrates integrated processing
func (app *BasicCliApp) runProcessCommand(ctx context.Context, inputPath, outputPath string, args []string) error {
	// Load configuration using pkg/config
	cfg, err := app.loadConfiguration()
	if err != nil {
		return errors.NewApplicationError(
			errors.CategoryConfiguration,
			errors.SeverityError,
			"CONFIG_LOAD_FAILED",
			fmt.Sprintf("Failed to load configuration: %v", err),
			map[string]interface{}{"command": "process"},
		)
	}

	// Get CLI context using pkg/cli patterns
	cmdCtx := cli.GetCommandContext(ctx)

	// Check dry-run mode
	if cmdCtx.DryRun {
		fmt.Println("DRY RUN MODE - No changes will be made")
	}

	// Create processing data
	processData := []map[string]interface{}{
		{
			"name":      "file1.txt",
			"size":      1024,
			"status":    "processed",
			"duration":  "15ms",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		{
			"name":      "file2.txt",
			"size":      2048,
			"status":    "processed",
			"duration":  "23ms",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		{
			"name":      "file3.txt",
			"size":      512,
			"status":    "skipped",
			"duration":  "5ms",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	}

	// Format output using pkg/formatter
	outputFormat := cfg.GetString("output_format", "table")

	switch outputFormat {
	case "json":
		jsonOutput, err := app.formatter.FormatJSON(processData)
		if err != nil {
			return fmt.Errorf("failed to format JSON output: %w", err)
		}
		fmt.Println(jsonOutput)

	case "yaml":
		yamlOutput, err := app.formatter.FormatYAML(processData)
		if err != nil {
			return fmt.Errorf("failed to format YAML output: %w", err)
		}
		fmt.Println(yamlOutput)

	case "table":
		headers := []string{"Name", "Size", "Status", "Duration", "Timestamp"}
		var rows [][]string
		for _, item := range processData {
			rows = append(rows, []string{
				item["name"].(string),
				fmt.Sprintf("%d bytes", item["size"].(int)),
				item["status"].(string),
				item["duration"].(string),
				item["timestamp"].(string),
			})
		}

		tableOutput, err := app.formatter.FormatTable(headers, rows)
		if err != nil {
			return fmt.Errorf("failed to format table output: %w", err)
		}
		fmt.Println(tableOutput)

	case "template":
		template := cfg.GetString("output.template", "{{.name}}: {{.status}}")
		for _, item := range processData {
			output, err := app.formatter.FormatTemplate(template, item)
			if err != nil {
				return fmt.Errorf("failed to format template output: %w", err)
			}
			fmt.Println(output)
		}

	default:
		return errors.NewApplicationError(
			errors.CategoryValidation,
			errors.SeverityError,
			"INVALID_FORMAT",
			fmt.Sprintf("Unsupported output format: %s", outputFormat),
			map[string]interface{}{"format": outputFormat, "supported": []string{"table", "json", "yaml", "template"}},
		)
	}

	// Show processing summary
	if cfg.GetBool("verbose", false) {
		fmt.Printf("\nProcessing Summary:\n")
		fmt.Printf("- Input Path: %s\n", inputPath)
		fmt.Printf("- Output Path: %s\n", outputPath)
		fmt.Printf("- Processed Items: %d\n", len(processData))
		fmt.Printf("- Dry Run: %v\n", cmdCtx.DryRun)
		fmt.Printf("- Format: %s\n", outputFormat)
	}

	return nil
}

// runStatusCommand demonstrates formatter integration
func (app *BasicCliApp) runStatusCommand(ctx context.Context) error {
	cfg, err := app.loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create status data
	statusData := map[string]interface{}{
		"application": map[string]interface{}{
			"name":    app.appInfo.Name,
			"version": app.appInfo.Version,
			"uptime":  "45 minutes",
			"status":  "running",
		},
		"configuration": map[string]interface{}{
			"format":      cfg.GetString("output_format", "table"),
			"verbose":     cfg.GetBool("verbose", false),
			"dry_run":     cfg.GetBool("dry_run", false),
			"config_file": cfg.GetString("config_file", "default"),
		},
		"statistics": map[string]interface{}{
			"total_processed": 156,
			"errors":          3,
			"warnings":        12,
			"last_run":        time.Now().Add(-time.Hour).Format(time.RFC3339),
		},
	}

	// Format using requested format
	outputFormat := cfg.GetString("output_format", "table")

	switch outputFormat {
	case "json":
		output, err := app.formatter.FormatJSON(statusData)
		if err != nil {
			return fmt.Errorf("failed to format status as JSON: %w", err)
		}
		fmt.Println(output)

	case "yaml":
		output, err := app.formatter.FormatYAML(statusData)
		if err != nil {
			return fmt.Errorf("failed to format status as YAML: %w", err)
		}
		fmt.Println(output)

	default: // table format
		fmt.Printf("Application Status\n")
		fmt.Printf("==================\n")
		fmt.Printf("Name:         %s\n", statusData["application"].(map[string]interface{})["name"])
		fmt.Printf("Version:      %s\n", statusData["application"].(map[string]interface{})["version"])
		fmt.Printf("Status:       %s\n", statusData["application"].(map[string]interface{})["status"])
		fmt.Printf("Uptime:       %s\n", statusData["application"].(map[string]interface{})["uptime"])
		fmt.Printf("\nConfiguration\n")
		fmt.Printf("=============\n")
		fmt.Printf("Format:       %s\n", statusData["configuration"].(map[string]interface{})["format"])
		fmt.Printf("Verbose:      %v\n", statusData["configuration"].(map[string]interface{})["verbose"])
		fmt.Printf("Dry Run:      %v\n", statusData["configuration"].(map[string]interface{})["dry_run"])
		fmt.Printf("Config File:  %s\n", statusData["configuration"].(map[string]interface{})["config_file"])
		fmt.Printf("\nStatistics\n")
		fmt.Printf("==========\n")
		fmt.Printf("Processed:    %d\n", statusData["statistics"].(map[string]interface{})["total_processed"])
		fmt.Printf("Errors:       %d\n", statusData["statistics"].(map[string]interface{})["errors"])
		fmt.Printf("Warnings:     %d\n", statusData["statistics"].(map[string]interface{})["warnings"])
		fmt.Printf("Last Run:     %s\n", statusData["statistics"].(map[string]interface{})["last_run"])
	}

	return nil
}

// runConfigShowCommand shows current configuration
func (app *BasicCliApp) runConfigShowCommand(ctx context.Context) error {
	cfg, err := app.loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Get all configuration values
	configData := cfg.GetAll()

	// Format and display
	output, err := app.formatter.FormatYAML(configData)
	if err != nil {
		return fmt.Errorf("failed to format configuration: %w", err)
	}

	fmt.Printf("Current Configuration:\n")
	fmt.Printf("======================\n")
	fmt.Println(output)

	return nil
}

// runConfigValidateCommand validates configuration
func (app *BasicCliApp) runConfigValidateCommand(ctx context.Context) error {
	_, err := app.loadConfiguration()
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		return err
	}

	fmt.Printf("✅ Configuration is valid\n")
	return nil
}

// loadConfiguration loads configuration from multiple sources
func (app *BasicCliApp) loadConfiguration() (config.ConfigInterface, error) {
	// Load default configuration
	defaultCfg := DefaultConfig()

	// Create config sources
	sources := []config.ConfigSource{
		// Default configuration
		config.NewMapConfigSource("defaults", map[string]interface{}{
			"output_format":          defaultCfg.OutputFormat,
			"verbose":                defaultCfg.Verbose,
			"dry_run":                defaultCfg.DryRun,
			"processing.batch_size":  defaultCfg.Processing.BatchSize,
			"processing.timeout":     defaultCfg.Processing.Timeout,
			"processing.concurrency": defaultCfg.Processing.Concurrency,
			"output.directory":       defaultCfg.Output.Directory,
			"output.template":        defaultCfg.Output.Template,
		}),

		// Environment variables
		config.NewEnvConfigSource("BASICAPP"),

		// Configuration file (if exists)
		config.NewFileConfigSource("config.yaml"),
	}

	// Load and merge configuration
	cfg, err := app.configLoader.LoadConfig(sources...)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	validator := config.NewConfigValidator()
	if err := validator.Validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// main function demonstrates CLI integration
func main() {
	app := NewBasicCliApp()

	// Build root command using pkg/cli patterns
	rootCmd := app.buildRootCommand()

	// Add version command using pkg/cli
	cli.AddVersionCommand(rootCmd, app.appInfo)

	// Create CLI context
	ctx := context.Background()
	cmdCtx := &cli.CommandContext{
		Context: ctx,
		DryRun:  false,
		Verbose: false,
	}

	// Execute with context
	if err := cli.ExecuteWithContext(rootCmd, cmdCtx); err != nil {
		// Handle errors using pkg/errors
		if appErr := errors.AsApplicationError(err); appErr != nil {
			log.Printf("Application error [%s]: %s", appErr.Code, appErr.Message)
			if appErr.Context != nil {
				log.Printf("Context: %+v", appErr.Context)
			}
		} else {
			log.Printf("Error: %v", err)
		}
		os.Exit(1)
	}
}
