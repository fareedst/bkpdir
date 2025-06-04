package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// ⭐ EXTRACT-005: CLI framework usage example - 📝

// Example shows how to use the CLI framework
func ExampleCLIApp() {
	// Define application information
	appInfo := AppInfo{
		Name:  "myapp",
		Short: "My CLI application",
		Long:  "A comprehensive CLI application built with the extracted framework",
		Build: BuildInfo{
			Version:  "1.0.0",
			Date:     "2024-01-01",
			Commit:   "abc123",
			Platform: "linux/amd64",
		},
	}

	// Create a new CLI application
	app := NewCLIApp(appInfo)

	// Create a simple command using the framework
	var dryRun bool
	var note string

	helloCmd := &cobra.Command{
		Use:   "hello",
		Short: "Say hello",
		Long:  "Say hello to demonstrate the CLI framework",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the dry-run manager
			ctx := CommandContext{
				Context:     cmd.Context(),
				Output:      os.Stdout,
				ErrorOutput: os.Stderr,
				DryRun:      dryRun,
			}

			dryRunMgr := NewDryRunManager()
			op := DryRunWrapper("Say hello with note: "+note, func(ctx CommandContext) error {
				fmt.Fprintf(ctx.Output, "Hello! Note: %s\n", note)
				return nil
			})

			return dryRunMgr.Execute(ctx, op)
		},
	}

	// Add flags using the flag manager
	flagMgr := NewFlagManager()
	flagMgr.AddDryRunFlag(helloCmd, &dryRun)
	flagMgr.AddNoteFlag(helloCmd, &note)

	// Add the command to the application
	app.AddCommand(helloCmd)

	// Execute with context and signal handling
	if err := app.ExecuteWithContext(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// ExampleCommandTemplate shows how to use the command template
func ExampleCommandTemplate() {
	// Create flag manager
	flagMgr := NewFlagManager()

	// Define template variables
	var dryRun bool
	var note string

	// Create command manually
	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy application",
		Long:    "Deploy the application with various options",
		Example: "  myapp deploy --note 'Production deployment' --dry-run",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := CommandContext{
				Context:     cmd.Context(),
				Output:      os.Stdout,
				ErrorOutput: os.Stderr,
				DryRun:      dryRun,
			}

			dryRunMgr := NewDryRunManager()
			op := DryRunWrapper("Deploy application with note: "+note, func(ctx CommandContext) error {
				fmt.Fprintf(ctx.Output, "Deploying application... Note: %s\n", note)
				return nil
			})

			return dryRunMgr.Execute(ctx, op)
		},
	}

	// Add flags using flag manager
	flagMgr.AddDryRunFlag(cmd, &dryRun)
	flagMgr.AddNoteFlag(cmd, &note)

	fmt.Printf("Command created: %s\n", cmd.Use)
	fmt.Printf("Short description: %s\n", cmd.Short)
}

// ExampleCancellableOperation shows how to use cancellable operations
func ExampleCancellableOperation() {
	contextMgr := NewContextManager()

	// Create a context with signal handling
	ctx, cancel := contextMgr.Create(context.Background())
	defer cancel()

	// Create a cancellable operation
	op := NewCancellableOperation(func(ctx context.Context) error {
		// Simulate long-running operation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println("Operation completed successfully")
			return nil
		}
	})

	// Execute the operation
	if err := op.Execute(ctx); err != nil {
		fmt.Printf("Operation error: %v\n", err)
	}
}

// Example_versionHandling shows how to use version management
func Example_versionHandling() {
	versionMgr := NewVersionManager()

	buildInfo := BuildInfo{
		Version:  "2.1.0",
		Date:     "2024-01-15",
		Commit:   "def456",
		Platform: "darwin/amd64",
	}

	// Format version information
	version := versionMgr.FormatVersion(buildInfo)
	fmt.Printf("Formatted version: %s\n", version)

	// Create version command
	versionCmd := versionMgr.CreateVersionCommand(buildInfo)
	fmt.Printf("Version command: %s\n", versionCmd.Use)

	// Create version template
	template := versionMgr.CreateVersionTemplate(buildInfo)
	fmt.Printf("Version template: %s", template)
}
