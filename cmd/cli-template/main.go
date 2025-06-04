// ⭐ EXTRACT-008: CLI Application Template - Main entry point demonstrating extracted packages
package main

import (
	"context"
	"os"

	"cli-template/cmd"
)

func main() {
	ctx := context.Background()

	// Execute the root command
	if err := cmd.Execute(ctx); err != nil {
		os.Exit(1)
	}
}
