// ⭐ EXTRACT-008: CLI Application Template - Version command demonstrating pkg/git integration
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version information including Git details (demonstrates pkg/git integration)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CLI Template v1.0.0")
		fmt.Println("Built with extracted packages from bkpdir")

		// TODO: Integrate pkg/git to show Git information
		fmt.Println("Git integration: Available via pkg/git")
		fmt.Println("- Repository detection")
		fmt.Println("- Branch and commit information")
		fmt.Println("- Dirty status detection")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
