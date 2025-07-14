// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
package main

import (
	"fmt"
	"os"

	"scaffolding/internal/generator"
	"scaffolding/internal/ui"
)

func main() {
	// CFG-003: See specification.md - Configuration Management [DECISION:maintenance]
	fmt.Println("🚀 Go CLI Project Scaffolding Generator")
	fmt.Println("=====================================")
	fmt.Println()

	// Collect project configuration from user
	config, err := ui.CollectProjectConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error collecting project configuration: %v\n", err)
		os.Exit(1)
	}

	// Generate the project
	generator := generator.New()
	if err := generator.GenerateProject(config); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n[SUCCESS] Successfully generated project '%s'\n", config.ProjectName)
	fmt.Printf("[DIR] Project location: %s\n", config.OutputPath)
	fmt.Println("\n🎯 Next steps:")
	fmt.Printf("   cd %s\n", config.ProjectName)
	fmt.Println("   make build    # Build the application")
	fmt.Println("   make demo     # Run demonstration")
	fmt.Println("   make help     # See all available commands")
}
