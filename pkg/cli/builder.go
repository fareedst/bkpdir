package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ⭐ EXTRACT-005: Command structure templates and builder patterns - 🔧

// DefaultCommandBuilder provides standard command building functionality
type DefaultCommandBuilder struct {
	flagManager FlagManager
}

// NewCommandBuilder creates a new command builder
func NewCommandBuilder(flagMgr FlagManager) CommandBuilder {
	if flagMgr == nil {
		flagMgr = NewFlagManager()
	}
	return &DefaultCommandBuilder{
		flagManager: flagMgr,
	}
}

// NewCommand creates a new command with standard setup
func (cb *DefaultCommandBuilder) NewCommand(name, short, long string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		Long:  long,
	}
}

// WithHandler sets the command handler function
func (cb *DefaultCommandBuilder) WithHandler(cmd *cobra.Command, handler func(*cobra.Command, []string) error) *cobra.Command {
	cmd.RunE = handler
	return cmd
}

// WithFlags adds flags using the flag manager
func (cb *DefaultCommandBuilder) WithFlags(cmd *cobra.Command, flags []string) *cobra.Command {
	for _, flag := range flags {
		switch flag {
		case "dry-run":
			// This would need a target variable - simplified for now
		case "note":
			// This would need a target variable - simplified for now
		case "config":
			// This would need a target variable - simplified for now
		}
	}
	return cmd
}

// WithSubcommands adds subcommands to a parent command
func (cb *DefaultCommandBuilder) WithSubcommands(parent *cobra.Command, children ...*cobra.Command) *cobra.Command {
	for _, child := range children {
		parent.AddCommand(child)
	}
	return parent
}

// DefaultRootCommandBuilder provides standard root command building functionality
type DefaultRootCommandBuilder struct {
	flagManager    FlagManager
	versionManager VersionManager
}

// NewRootCommandBuilder creates a new root command builder
func NewRootCommandBuilder(flagMgr FlagManager, versionMgr VersionManager) RootCommandBuilder {
	if flagMgr == nil {
		flagMgr = NewFlagManager()
	}
	if versionMgr == nil {
		versionMgr = NewVersionManager()
	}
	return &DefaultRootCommandBuilder{
		flagManager:    flagMgr,
		versionManager: versionMgr,
	}
}

// NewRootCommand creates the root command with application info
func (rb *DefaultRootCommandBuilder) NewRootCommand(info AppInfo) *cobra.Command {
	// Create the long description
	longDesc := info.Long

	cmd := &cobra.Command{
		Use:     info.Name,
		Short:   info.Short,
		Long:    longDesc,
		Version: rb.versionManager.FormatVersion(info.Build),
		Run: func(cmd *cobra.Command, args []string) {
			// Default behavior: show help if no subcommand is provided
			cmd.Help()
		},
	}

	// Set version template
	template := rb.versionManager.CreateVersionTemplate(info.Build)
	cmd.SetVersionTemplate(template)

	return cmd
}

// WithVersionTemplate sets custom version template
func (rb *DefaultRootCommandBuilder) WithVersionTemplate(cmd *cobra.Command, template string) *cobra.Command {
	cmd.SetVersionTemplate(template)
	return cmd
}

// WithGlobalFlags adds global flags to root command
func (rb *DefaultRootCommandBuilder) WithGlobalFlags(cmd *cobra.Command, flagMgr FlagManager) *cobra.Command {
	if flagMgr != nil {
		flagMgr.AddGlobalFlags(cmd)
	} else {
		rb.flagManager.AddGlobalFlags(cmd)
	}
	return cmd
}

// WithExampleUsage adds example usage to root command
func (rb *DefaultRootCommandBuilder) WithExampleUsage(cmd *cobra.Command, examples string) *cobra.Command {
	cmd.Example = examples
	return cmd
}

// CommandTemplate provides a template for common command patterns
type CommandTemplate struct {
	Name        string
	Short       string
	Long        string
	Example     string
	FlagSet     FlagSet
	Handler     func(*cobra.Command, []string) error
	PreRun      func(*cobra.Command, []string) error
	PostRun     func(*cobra.Command, []string) error
	Subcommands []*cobra.Command
}

// BuildCommand creates a cobra command from a template
func (rb *DefaultRootCommandBuilder) BuildCommand(template CommandTemplate) *cobra.Command {
	cmd := &cobra.Command{
		Use:     template.Name,
		Short:   template.Short,
		Long:    template.Long,
		Example: template.Example,
	}

	if template.Handler != nil {
		cmd.RunE = template.Handler
	}

	if template.PreRun != nil {
		cmd.PreRunE = template.PreRun
	}

	if template.PostRun != nil {
		cmd.PostRunE = template.PostRun
	}

	// Add flags manually for now - interface needs to be extended
	flagMgr := rb.flagManager.(*DefaultFlagManager)
	flagMgr.AddFlags(cmd, template.FlagSet)

	// Add subcommands
	for _, subcmd := range template.Subcommands {
		cmd.AddCommand(subcmd)
	}

	return cmd
}

// CLIApp represents a complete CLI application
type CLIApp struct {
	Info           AppInfo
	RootCommand    *cobra.Command
	rootBuilder    RootCommandBuilder
	commandBuilder CommandBuilder
	versionManager VersionManager
	flagManager    FlagManager
	contextManager ContextManager
	dryRunManager  DryRunManager
}

// NewCLIApp creates a new CLI application
func NewCLIApp(info AppInfo) *CLIApp {
	versionMgr := NewVersionManager()
	flagMgr := NewFlagManager()
	rootBuilder := NewRootCommandBuilder(flagMgr, versionMgr)
	commandBuilder := NewCommandBuilder(flagMgr)
	contextMgr := NewContextManager()
	dryRunMgr := NewDryRunManager()

	rootCmd := rootBuilder.NewRootCommand(info)

	return &CLIApp{
		Info:           info,
		RootCommand:    rootCmd,
		rootBuilder:    rootBuilder,
		commandBuilder: commandBuilder,
		versionManager: versionMgr,
		flagManager:    flagMgr,
		contextManager: contextMgr,
		dryRunManager:  dryRunMgr,
	}
}

// AddCommand adds a command to the CLI application
func (app *CLIApp) AddCommand(cmd *cobra.Command) {
	app.RootCommand.AddCommand(cmd)
}

// Execute runs the CLI application
func (app *CLIApp) Execute() error {
	return app.RootCommand.Execute()
}

// ExecuteWithContext runs the CLI application with signal handling
func (app *CLIApp) ExecuteWithContext() error {
	ctx, cancel := WithSignalHandling(nil)
	defer cancel()

	// Store context for use by commands
	app.RootCommand.SetContext(ctx)

	if err := app.RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}
