// 🔺 EXTRACT-008: Project template generation - 📝 Template processing and file creation
package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"scaffolding/internal/templates"
	"scaffolding/internal/ui"
)

// Generator handles project generation from templates
type Generator struct {
	templates *templates.Manager
}

// New creates a new generator instance
func New() *Generator {
	return &Generator{
		templates: templates.NewManager(),
	}
}

// GenerateProject creates a new CLI project based on the configuration
func (g *Generator) GenerateProject(config *ui.ProjectConfig) error {
	// 🔺 EXTRACT-008: Project directory creation - 🔧 Directory structure setup
	if err := g.createProjectDirectory(config); err != nil {
		return fmt.Errorf("creating project directory: %w", err)
	}

	// 🔺 EXTRACT-008: Core file generation - 📝 Template-based file creation
	if err := g.generateCoreFiles(config); err != nil {
		return fmt.Errorf("generating core files: %w", err)
	}

	// 🔺 EXTRACT-008: Package-specific generation - 🔧 Component integration
	if err := g.generatePackageFiles(config); err != nil {
		return fmt.Errorf("generating package files: %w", err)
	}

	// 🔺 EXTRACT-008: Optional feature generation - 📝 Additional features
	if err := g.generateOptionalFiles(config); err != nil {
		return fmt.Errorf("generating optional files: %w", err)
	}

	fmt.Printf("📁 Project structure created:\n")
	g.printProjectStructure(config)

	return nil
}

func (g *Generator) createProjectDirectory(config *ui.ProjectConfig) error {
	// Create main project directory
	if err := os.MkdirAll(config.OutputPath, 0755); err != nil {
		return fmt.Errorf("creating project directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{"cmd", "config", "internal"}

	if config.IncludeTests {
		subdirs = append(subdirs, "test")
	}

	for _, subdir := range subdirs {
		dirPath := filepath.Join(config.OutputPath, subdir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("creating subdirectory %s: %w", subdir, err)
		}
	}

	return nil
}

func (g *Generator) generateCoreFiles(config *ui.ProjectConfig) error {
	// Generate go.mod
	goModContent := g.templates.GenerateGoMod(config)
	if err := g.writeFile(config.OutputPath, "go.mod", goModContent); err != nil {
		return fmt.Errorf("writing go.mod: %w", err)
	}

	// Generate main.go
	mainContent := g.templates.GenerateMain(config)
	if err := g.writeFile(config.OutputPath, "main.go", mainContent); err != nil {
		return fmt.Errorf("writing main.go: %w", err)
	}

	// Generate root command
	rootCmdContent := g.templates.GenerateRootCommand(config)
	if err := g.writeFile(filepath.Join(config.OutputPath, "cmd"), "root.go", rootCmdContent); err != nil {
		return fmt.Errorf("writing cmd/root.go: %w", err)
	}

	// Generate version command
	versionCmdContent := g.templates.GenerateVersionCommand(config)
	if err := g.writeFile(filepath.Join(config.OutputPath, "cmd"), "version.go", versionCmdContent); err != nil {
		return fmt.Errorf("writing cmd/version.go: %w", err)
	}

	// Generate Makefile
	makefileContent := g.templates.GenerateMakefile(config)
	if err := g.writeFile(config.OutputPath, "Makefile", makefileContent); err != nil {
		return fmt.Errorf("writing Makefile: %w", err)
	}

	// Generate README.md
	readmeContent := g.templates.GenerateReadme(config)
	if err := g.writeFile(config.OutputPath, "README.md", readmeContent); err != nil {
		return fmt.Errorf("writing README.md: %w", err)
	}

	return nil
}

func (g *Generator) generatePackageFiles(config *ui.ProjectConfig) error {
	// Generate configuration files if config package is included
	if g.containsPackage(config.SelectedPackages, "config") {
		configContent := g.templates.GenerateDefaultConfig(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "config"), "default.yml", configContent); err != nil {
			return fmt.Errorf("writing config/default.yml: %w", err)
		}

		exampleConfigContent := g.templates.GenerateExampleConfig(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "config"), "example.yml", exampleConfigContent); err != nil {
			return fmt.Errorf("writing config/example.yml: %w", err)
		}
	}

	// Generate additional commands based on selected packages
	if g.containsPackage(config.SelectedPackages, "config") {
		configCmdContent := g.templates.GenerateConfigCommand(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "cmd"), "config.go", configCmdContent); err != nil {
			return fmt.Errorf("writing cmd/config.go: %w", err)
		}
	}

	if g.containsPackage(config.SelectedPackages, "fileops") {
		createCmdContent := g.templates.GenerateCreateCommand(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "cmd"), "create.go", createCmdContent); err != nil {
			return fmt.Errorf("writing cmd/create.go: %w", err)
		}
	}

	if g.containsPackage(config.SelectedPackages, "processing") {
		processCmdContent := g.templates.GenerateProcessCommand(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "cmd"), "process.go", processCmdContent); err != nil {
			return fmt.Errorf("writing cmd/process.go: %w", err)
		}
	}

	return nil
}

func (g *Generator) generateOptionalFiles(config *ui.ProjectConfig) error {
	// Generate .gitignore if requested
	if config.IncludeGitIgnore {
		gitignoreContent := g.templates.GenerateGitignore(config)
		if err := g.writeFile(config.OutputPath, ".gitignore", gitignoreContent); err != nil {
			return fmt.Errorf("writing .gitignore: %w", err)
		}
	}

	// Generate Dockerfile if requested
	if config.IncludeDockerfile {
		dockerfileContent := g.templates.GenerateDockerfile(config)
		if err := g.writeFile(config.OutputPath, "Dockerfile", dockerfileContent); err != nil {
			return fmt.Errorf("writing Dockerfile: %w", err)
		}
	}

	// Generate test files if requested
	if config.IncludeTests {
		testContent := g.templates.GenerateTestFile(config)
		if err := g.writeFile(filepath.Join(config.OutputPath, "test"), "main_test.go", testContent); err != nil {
			return fmt.Errorf("writing test/main_test.go: %w", err)
		}
	}

	return nil
}

func (g *Generator) writeFile(dir, filename, content string) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func (g *Generator) containsPackage(packages []string, target string) bool {
	for _, pkg := range packages {
		if pkg == target {
			return true
		}
	}
	return false
}

func (g *Generator) printProjectStructure(config *ui.ProjectConfig) {
	fmt.Printf("%s/\n", config.ProjectName)
	fmt.Printf("├── main.go\n")
	fmt.Printf("├── go.mod\n")
	fmt.Printf("├── Makefile\n")
	fmt.Printf("├── README.md\n")

	if config.IncludeGitIgnore {
		fmt.Printf("├── .gitignore\n")
	}

	if config.IncludeDockerfile {
		fmt.Printf("├── Dockerfile\n")
	}

	fmt.Printf("├── cmd/\n")
	fmt.Printf("│   ├── root.go\n")
	fmt.Printf("│   └── version.go\n")

	if g.containsPackage(config.SelectedPackages, "config") {
		fmt.Printf("│   └── config.go\n")
		fmt.Printf("├── config/\n")
		fmt.Printf("│   ├── default.yml\n")
		fmt.Printf("│   └── example.yml\n")
	}

	if g.containsPackage(config.SelectedPackages, "fileops") {
		fmt.Printf("│   └── create.go\n")
	}

	if g.containsPackage(config.SelectedPackages, "processing") {
		fmt.Printf("│   └── process.go\n")
	}

	if config.IncludeTests {
		fmt.Printf("└── test/\n")
		fmt.Printf("    └── main_test.go\n")
	}

	fmt.Printf("\n🔗 Included packages: %s\n", strings.Join(config.SelectedPackages, ", "))
}
