// Package cli provides a reusable command-line interface framework
// extracted from the bkpdir application. It includes Cobra command patterns,
// flag handling, dry-run support, context-aware execution, and version management.
//
// This package is designed to accelerate development of new CLI applications
// by providing tested patterns and reusable components.
package cli

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

// ⭐ EXTRACT-005: Core CLI framework types and interfaces - 🔧

// BuildInfo contains version and build-time information for the application
type BuildInfo struct {
	// Version is the application version (e.g., "1.0.0")
	Version string
	// Date is the compilation date
	Date string
	// Commit is the Git commit hash
	Commit string
	// Platform contains platform information
	Platform string
}

// AppInfo contains application metadata and configuration
type AppInfo struct {
	// Name is the application name
	Name string
	// Short is a brief description
	Short string
	// Long is a detailed description
	Long string
	// Build contains version and build information
	Build BuildInfo
}

// CommandContext provides shared context and dependencies for commands
type CommandContext struct {
	// Context for cancellation and timeouts
	Context context.Context
	// Output writer for command output
	Output io.Writer
	// Error writer for error messages
	ErrorOutput io.Writer
	// DryRun indicates if operations should be simulated
	DryRun bool
}

// DryRunOperation represents an operation that can be executed or simulated
type DryRunOperation interface {
	// Execute performs the actual operation
	Execute(ctx CommandContext) error
	// Describe returns a description of what the operation would do
	Describe() string
}

// DryRunManager handles dry-run execution logic
type DryRunManager interface {
	// Execute runs the operation or simulates it based on dry-run flag
	Execute(ctx CommandContext, op DryRunOperation) error
	// Log outputs dry-run information
	Log(ctx CommandContext, message string)
}

// CancellableOperation represents an operation that supports cancellation
type CancellableOperation interface {
	// Execute performs the operation with cancellation support
	Execute(ctx context.Context) error
	// Cancel requests cancellation of the operation
	Cancel() error
}

// ContextManager handles context lifecycle for commands
type ContextManager interface {
	// Create creates a new context with cancellation
	Create(parent context.Context) (context.Context, context.CancelFunc)
	// WithTimeout creates a context with timeout
	WithTimeout(parent context.Context, timeout string) (context.Context, context.CancelFunc)
	// HandleSignals sets up signal handling for graceful shutdown
	HandleSignals(cancel context.CancelFunc)
}

// FlagManager handles consistent flag registration and binding
type FlagManager interface {
	// AddGlobalFlags adds common global flags to a command
	AddGlobalFlags(cmd *cobra.Command) error
	// AddDryRunFlag adds dry-run flag with consistent naming
	AddDryRunFlag(cmd *cobra.Command, target *bool) error
	// AddNoteFlag adds note flag for operations
	AddNoteFlag(cmd *cobra.Command, target *string) error
	// AddConfigFlag adds configuration flag
	AddConfigFlag(cmd *cobra.Command, target *bool) error
}

// CommandBuilder provides patterns for building Cobra commands
type CommandBuilder interface {
	// NewCommand creates a new command with standard setup
	NewCommand(name, short, long string) *cobra.Command
	// WithHandler sets the command handler function
	WithHandler(cmd *cobra.Command, handler func(*cobra.Command, []string) error) *cobra.Command
	// WithFlags adds flags using the flag manager
	WithFlags(cmd *cobra.Command, flags []string) *cobra.Command
	// WithSubcommands adds subcommands to a parent command
	WithSubcommands(parent *cobra.Command, children ...*cobra.Command) *cobra.Command
}

// RootCommandBuilder builds the root command with version and global configuration
type RootCommandBuilder interface {
	// NewRootCommand creates the root command with application info
	NewRootCommand(info AppInfo) *cobra.Command
	// WithVersionTemplate sets custom version template
	WithVersionTemplate(cmd *cobra.Command, template string) *cobra.Command
	// WithGlobalFlags adds global flags to root command
	WithGlobalFlags(cmd *cobra.Command, flagMgr FlagManager) *cobra.Command
	// WithExampleUsage adds example usage to root command
	WithExampleUsage(cmd *cobra.Command, examples string) *cobra.Command
}

// VersionManager handles version display and build information
type VersionManager interface {
	// FormatVersion formats version information for display
	FormatVersion(info BuildInfo) string
	// CreateVersionCommand creates a version subcommand
	CreateVersionCommand(info BuildInfo) *cobra.Command
	// CreateVersionTemplate creates a version template string
	CreateVersionTemplate(info BuildInfo) string
}
