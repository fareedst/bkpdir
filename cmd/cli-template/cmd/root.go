// ⭐ EXTRACT-008: CLI Application Template - Root command demonstrating extracted packages
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Root command
	rootCmd = &cobra.Command{
		Use:   "cli-template",
		Short: "CLI Template demonstrating extracted packages",
		Long: `
A complete CLI application template showcasing all extracted packages:
- pkg/config: Configuration management
- pkg/errors: Structured error handling  
- pkg/resources: Resource management
- pkg/formatter: Output formatting
- pkg/git: Git integration
- pkg/cli: Command framework
- pkg/fileops: File operations
- pkg/processing: Concurrent processing

This template serves as a foundation for building new CLI applications.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CLI Template Application")
			fmt.Println("Run 'cli-template --help' for available commands")
		},
	}

	// Global flags
	configFile string
	verbose    bool
	dryRun     bool
)

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.cli-template.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would be done without executing")
}

// Execute runs the root command
func Execute(ctx context.Context) error {
	// Set context for all commands
	rootCmd.SetContext(ctx)

	// Execute command
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return err
	}

	return nil
}
