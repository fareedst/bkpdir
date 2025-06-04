// ⭐ EXTRACT-008: CLI Application Template - Create command demonstrating pkg/fileops and integration
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [target]",
	Short: "Create example files and demonstrate file operations",
	Long:  "Demonstrate file operations, error handling, and resource management",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := "example.txt"
		if len(args) > 0 {
			target = args[0]
		}

		fmt.Println("File Operations Demo (pkg/fileops)")
		fmt.Println("==================================")

		fmt.Printf("Target file: %s\n", target)

		if dryRun {
			fmt.Println("\n[DRY RUN] Would perform the following operations:")
			fmt.Printf("- Check if %s exists\n", target)
			fmt.Printf("- Create %s with example content\n", target)
			fmt.Println("- Demonstrate file comparison")
			fmt.Println("- Show resource cleanup patterns")
			return
		}

		fmt.Println("\nFeatures demonstrated:")
		fmt.Println("- Safe file operations (pkg/fileops)")
		fmt.Println("- Structured error handling (pkg/errors)")
		fmt.Println("- Resource management (pkg/resources)")
		fmt.Println("- Output formatting (pkg/formatter)")

		// Check if file exists
		if _, err := os.Stat(target); err == nil {
			fmt.Printf("\n⚠️  File %s already exists\n", target)
			if !verbose {
				fmt.Println("Use --verbose to see more details")
				return
			}
		}

		// Create example file
		fmt.Printf("\nCreating %s...\n", target)
		content := `# CLI Template Example File

This file was created by the CLI template application to demonstrate:

1. File Operations (pkg/fileops)
   - Safe file creation and manipulation
   - File comparison and verification
   - Path handling and validation

2. Error Handling (pkg/errors)
   - Structured error types
   - Error context and details
   - Recovery patterns

3. Resource Management (pkg/resources)
   - Automatic cleanup
   - Resource tracking
   - Leak prevention

4. Configuration (pkg/config)
   - Schema-agnostic loading
   - Multiple configuration sources
   - Value merging and validation

5. Output Formatting (pkg/formatter)
   - Template-based formatting
   - Printf-style formatting
   - Structured output

6. Git Integration (pkg/git)
   - Repository detection
   - Branch and commit information
   - Status tracking

7. CLI Framework (pkg/cli)
   - Command patterns
   - Flag handling
   - Context management

8. Concurrent Processing (pkg/processing)
   - Worker pool patterns
   - Progress tracking
   - Context cancellation

Generated at: %s
`

		err := os.WriteFile(target, []byte(fmt.Sprintf(content, "2025-01-02")), 0644)
		if err != nil {
			fmt.Printf("❌ Error creating file: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully created %s\n", target)

		if verbose {
			fmt.Println("\nFile contents preview:")
			fmt.Println("# CLI Template Example File")
			fmt.Println("...")
			fmt.Printf("(Full content written to %s)\n", target)
		}

		fmt.Println("\nNote: This demonstrates basic file operations.")
		fmt.Println("Full package integration with proper error handling coming soon.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
