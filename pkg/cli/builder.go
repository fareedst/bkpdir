// Package cli provides command and root builders for the CLI framework.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// DefaultCommandBuilder provides standard command building functionality
type DefaultCommandBuilder struct {
	flagManager FlagManager
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultCommandBuilder with a non-nil FlagManager (defaulting when nil).
func NewCommandBuilder(flagMgr FlagManager) CommandBuilder {
	if flagMgr == nil {
		flagMgr = NewFlagManager()
	}
	return &DefaultCommandBuilder{
		flagManager: flagMgr,
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create a subcommand with use, short, and long description fields.
func (cb *DefaultCommandBuilder) NewCommand(name, short, long string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		Long:  long,
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach RunE handler to command and return command for fluent chaining.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: register each child command on parent and return parent.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultRootCommandBuilder with default flag and version managers when nil.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: build root command from AppInfo with formatted version, version template, and help when invoked without subcommand.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: add persistent help and verbose flags via provided or embedded FlagManager.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: materialize command from CommandTemplate including handler, prerun, postrun, flags, and subcommands.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wire all managers and builders then create root command from AppInfo.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach subcommand to application root.
func (app *CLIApp) AddCommand(cmd *cobra.Command) {
	app.RootCommand.AddCommand(cmd)
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root command without extra signal context.
func (app *CLIApp) Execute() error {
	return app.RootCommand.Execute()
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root with signal-handled context on command; exit non-zero on error.
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
