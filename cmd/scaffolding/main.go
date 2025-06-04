// 🔺 EXTRACT-008: Project scaffolding system - 🔧 Interactive generator core
package main

import (
	"fmt"
	"os"

	"scaffolding/internal/generator"
	"scaffolding/internal/ui"
)

func main() {
	// 🔺 EXTRACT-008: Interactive CLI configuration - 🔍 User input collection and validation
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

	fmt.Printf("\n✅ Successfully generated project '%s'\n", config.ProjectName)
	fmt.Printf("📁 Project location: %s\n", config.OutputPath)
	fmt.Println("\n🎯 Next steps:")
	fmt.Printf("   cd %s\n", config.ProjectName)
	fmt.Println("   make build    # Build the application")
	fmt.Println("   make demo     # Run demonstration")
	fmt.Println("   make help     # See all available commands")
}
