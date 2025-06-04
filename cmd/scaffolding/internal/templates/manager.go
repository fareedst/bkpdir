// 🔺 EXTRACT-008: Template management system - 📝 File generation coordination
package templates

import (
	"fmt"
	"strings"

	"scaffolding/internal/ui"
)

// Manager handles template generation for different file types
type Manager struct{}

// NewManager creates a new template manager
func NewManager() *Manager {
	return &Manager{}
}

// GenerateGoMod creates the go.mod file content
func (m *Manager) GenerateGoMod(config *ui.ProjectConfig) string {
	var deps []string

	if m.containsPackage(config.SelectedPackages, "cli") {
		deps = append(deps, "\tgithub.com/spf13/cobra v1.9.1")
	}

	if m.containsPackage(config.SelectedPackages, "config") {
		deps = append(deps, "\tgopkg.in/yaml.v2 v2.4.0")
	}

	content := fmt.Sprintf(`module %s

go 1.24.3`, config.ModuleName)

	if len(deps) > 0 {
		content += "\n\nrequire (\n" + strings.Join(deps, "\n") + "\n)"
	}

	return content
}

// GenerateMain creates the main.go file content
func (m *Manager) GenerateMain(config *ui.ProjectConfig) string {
	return fmt.Sprintf(`// %s - %s
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"%s/cmd"
)

func main() {
	// Create context that cancels on interrupt signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Execute root command with context
	if err := cmd.Execute(ctx); err != nil {
		os.Exit(1)
	}
}`, config.ProjectName, config.Description, config.ModuleName)
}

// GenerateRootCommand creates the root command file
func (m *Manager) GenerateRootCommand(config *ui.ProjectConfig) string {
	imports := []string{
		`"context"`,
		`"fmt"`,
		`"os"`,
		`"github.com/spf13/cobra"`,
	}

	if m.containsPackage(config.SelectedPackages, "config") {
		imports = append(imports, `"path/filepath"`)
	}

	importBlock := strings.Join(imports, "\n\t")

	configVars := ""
	configFlags := ""
	configInit := ""

	if m.containsPackage(config.SelectedPackages, "config") {
		configVars = `
var (
	configFile string
	verbose    bool
	dryRun     bool
)`

		configFlags = `
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./config/default.yml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "dry run mode")`

		configInit = fmt.Sprintf(`
	// Initialize configuration
	if configFile == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			configFile = filepath.Join(homeDir, ".config", "%s", "config.yml")
		}
	}`, config.ProjectName)
	}

	return fmt.Sprintf(`package cmd

import (
	%s
)
%s
var rootCmd = &cobra.Command{
	Use:   "%s",
	Short: "%s",
	Long:  "%s",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command
func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {%s%s
}`, importBlock, configVars, config.ProjectName, config.Description, config.Description, configFlags, configInit)
}

// GenerateVersionCommand creates the version command
func (m *Manager) GenerateVersionCommand(config *ui.ProjectConfig) string {
	gitImport := ""
	gitLogic := `fmt.Printf("%s version: development\\n", cmd.Parent().Name())`

	if m.containsPackage(config.SelectedPackages, "git") {
		gitImport = `
	"os/exec"`
		gitLogic = `
	// Get git information if available
	version := "development"
	if gitHash, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		version = fmt.Sprintf("development-%s", strings.TrimSpace(string(gitHash)))
		
		if gitBranch, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
			version = fmt.Sprintf("%s (%s)", version, strings.TrimSpace(string(gitBranch)))
		}
	}
	
	fmt.Printf("%s version: %s\\n", cmd.Parent().Name(), version)`
	}

	return fmt.Sprintf(`package cmd

import (
	"fmt"
	"strings"%s

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version information for this application",
	RunE: func(cmd *cobra.Command, args []string) error {%s
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}`, gitImport, gitLogic)
}

// GenerateConfigCommand creates the config command
func (m *Manager) GenerateConfigCommand(config *ui.ProjectConfig) string {
	return `package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display configuration information",
	Long:  "Display current configuration values and sources",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPaths := []string{
			"./config/default.yml",
			"./config/example.yml",
		}
		
		if configFile != "" {
			configPaths = append([]string{configFile}, configPaths...)
		}
		
		fmt.Printf("Configuration for %s\n", cmd.Parent().Name())
		fmt.Printf("====================\n\n")
		
		for _, path := range configPaths {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("📁 Config file: %s\n", path)
				
				data, err := os.ReadFile(path)
				if err != nil {
					fmt.Printf("   ❌ Error reading: %v\n", err)
					continue
				}
				
				var config map[string]interface{}
				if err := yaml.Unmarshal(data, &config); err != nil {
					fmt.Printf("   ❌ Error parsing: %v\n", err)
					continue
				}
				
				fmt.Printf("   ✅ Loaded successfully\n")
				if verbose {
					fmt.Printf("   📊 Content:\n")
					for key, value := range config {
						fmt.Printf("      %s: %v\n", key, value)
					}
				}
				fmt.Println()
			}
		}
		
		fmt.Printf("🔧 Flags:\n")
		fmt.Printf("   verbose: %v\n", verbose)
		fmt.Printf("   dry-run: %v\n", dryRun)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}`
}

// GenerateCreateCommand creates the create command for file operations
func (m *Manager) GenerateCreateCommand(config *ui.ProjectConfig) string {
	return `package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [filename]",
	Short: "Create a new file with example content",
	Long:  "Create a new file with example content demonstrating file operations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := "example.txt"
		if len(args) > 0 {
			filename = args[0]
		}
		
		// Check if file exists
		if _, err := os.Stat(filename); err == nil {
			if !dryRun {
				return fmt.Errorf("file already exists: %s", filename)
			}
			fmt.Printf("🔍 [DRY RUN] File exists: %s\n", filename)
		}
		
		content := fmt.Sprintf("# Example file created by %s\n\nCreated at: %s\nContent: This is an example file demonstrating file operations.\n", 
			cmd.Parent().Name(), time.Now().Format(time.RFC3339))
		
		if dryRun {
			fmt.Printf("🔍 [DRY RUN] Would create file: %s\n", filename)
			fmt.Printf("📝 [DRY RUN] Content:\n%s\n", content)
			return nil
		}
		
		// Ensure directory exists
		if dir := filepath.Dir(filename); dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating directory: %w", err)
			}
		}
		
		// Write file
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		
		fmt.Printf("✅ Created file: %s\n", filename)
		if verbose {
			fmt.Printf("📁 Full path: %s\n", filepath.Join(".", filename))
			fmt.Printf("📊 Size: %d bytes\n", len(content))
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}`
}

// GenerateProcessCommand creates the process command for concurrent operations
func (m *Manager) GenerateProcessCommand(config *ui.ProjectConfig) string {
	return `package cmd

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var processCmd = &cobra.Command{
	Use:   "process [count]",
	Short: "Demonstrate concurrent processing",
	Long:  "Run concurrent processing simulation to demonstrate processing patterns",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		count := 5
		if len(args) > 0 {
			fmt.Sscanf(args[0], "%d", &count)
		}
		
		fmt.Printf("🔄 Starting processing simulation with %d workers\n", count)
		
		if dryRun {
			fmt.Printf("🔍 [DRY RUN] Would process %d items\n", count)
			return nil
		}
		
		ctx := cmd.Context()
		return runProcessingSimulation(ctx, count)
	},
}

func runProcessingSimulation(ctx context.Context, workerCount int) error {
	// Create work items
	workItems := make(chan int, workerCount*2)
	results := make(chan string, workerCount*2)
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, i, workItems, results, &wg)
	}
	
	// Send work
	go func() {
		defer close(workItems)
		for i := 0; i < workerCount*2; i++ {
			select {
			case workItems <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	
	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Display results
	for result := range results {
		fmt.Println(result)
	}
	
	fmt.Printf("✅ Processing simulation completed\n")
	return nil
}

func worker(ctx context.Context, id int, work <-chan int, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for {
		select {
		case item, ok := <-work:
			if !ok {
				return
			}
			
			// Simulate work
			duration := time.Duration(rand.Intn(1000)+500) * time.Millisecond
			time.Sleep(duration)
			
			results <- fmt.Sprintf("🔧 Worker %d processed item %d (took %v)", id, item, duration)
			
		case <-ctx.Done():
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(processCmd)
}`
}

// GenerateMakefile creates the Makefile content
func (m *Manager) GenerateMakefile(config *ui.ProjectConfig) string {
	return fmt.Sprintf(`# Makefile for %s

.PHONY: help build test clean demo install lint

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %%s\n", $$1, $$2}'

build: ## Build the application
	@echo "🔨 Building %s..."
	@go build -o bin/%s .
	@echo "✅ Build completed: bin/%s"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

demo: build ## Build and run demonstration
	@echo "🚀 Running demonstration..."
	@./bin/%s version
	@echo ""
	@./bin/%s config --verbose
	@echo ""
	@./bin/%s create example-demo.txt
	@echo ""
	@./bin/%s process 3

install: build ## Install the application
	@echo "📦 Installing %s..."
	@go install .
	@echo "✅ Installation completed"

lint: ## Run linting
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

dev: ## Development mode with file watching
	@echo "👁️  Starting development mode..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	@air

# Build variables
VERSION ?= development
LDFLAGS = -X main.version=$(VERSION)

build-release: ## Build release version
	@echo "🔨 Building release version $(VERSION)..."
	@go build -ldflags "$(LDFLAGS)" -o bin/%s .
	@echo "✅ Release build completed: bin/%s"`,
		config.ProjectName, config.ProjectName, config.ProjectName, config.ProjectName,
		config.ProjectName, config.ProjectName, config.ProjectName, config.ProjectName,
		config.ProjectName, config.ProjectName, config.ProjectName)
}

// GenerateReadme creates the README.md content
func (m *Manager) GenerateReadme(config *ui.ProjectConfig) string {
	packageList := strings.Join(config.SelectedPackages, ", ")
	commandHelp := m.generateCommandHelp(config)
	structureExtras := m.generateStructureExtras(config)

	return fmt.Sprintf(`# %s

%s

## Features

This CLI application demonstrates the use of extracted packages:
- **Packages Used**: %s
- **Template**: %s

## Installation

### Prerequisites

- Go 1.21 or later

### Build from Source

`+"```bash"+`
# Clone the repository
git clone <repository-url>
cd %s

# Build the application
make build

# Install globally
make install
`+"```"+`

## Usage

### Basic Commands

`+"```bash"+`
# Show help
%s help

# Show version
%s version

# Display configuration
%s config --verbose
`+"```"+`

### Available Commands

- **help**: Show help information
- **version**: Display version information
- **config**: Show configuration values%s

## Configuration

The application looks for configuration files in the following order:
1. File specified by --config flag
2. ./config/default.yml
3. ./config/example.yml

## Development

### Build and Test

`+"```bash"+`
# Build the application
make build

# Run tests
make test

# Run demonstration
make demo

# Development mode with file watching
make dev
`+"```"+`

### Project Structure

`+"```"+`
%s/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── Makefile               # Build automation
├── README.md              # This file
├── cmd/                   # Command definitions
│   ├── root.go           # Root command
│   └── version.go        # Version command%s
└── config/               # Configuration files
    ├── default.yml       # Default configuration
    └── example.yml       # Example configuration
`+"```"+`

## Author

%s

## License

This project is generated from a template. Please add your license here.

---

*Generated by Go CLI Project Scaffolding Generator*`,
		config.ProjectName, config.Description, packageList, config.Template,
		config.ProjectName, config.ProjectName, config.ProjectName, config.ProjectName,
		commandHelp, config.ProjectName, structureExtras, config.Author)
}

// GenerateDefaultConfig creates the default configuration file
func (m *Manager) GenerateDefaultConfig(config *ui.ProjectConfig) string {
	packageConfig := m.generatePackageConfig(config)

	return fmt.Sprintf(`# Default configuration for %s
app:
  name: "%s"
  version: "development"
  description: "%s"

# Logging configuration
logging:
  level: "info"
  format: "text"

# Feature flags
features:
  verbose: false
  dry_run: false

# Package-specific settings%s`,
		config.ProjectName, config.ProjectName, config.Description, packageConfig)
}

// GenerateExampleConfig creates the example configuration file
func (m *Manager) GenerateExampleConfig(config *ui.ProjectConfig) string {
	packageConfig := m.generatePackageConfig(config)

	return fmt.Sprintf(`# Example configuration for %s
# Copy this file and modify as needed

app:
  name: "%s"
  version: "1.0.0"
  description: "%s"

# Enhanced logging configuration
logging:
  level: "debug"
  format: "json"
  output: "stdout"

# Feature flags
features:
  verbose: true
  dry_run: false
  enable_colors: true

# Development settings
development:
  auto_reload: true
  debug_mode: true

# Package-specific settings%s`,
		config.ProjectName, config.ProjectName, config.Description, packageConfig)
}

// GenerateGitignore creates the .gitignore file
func (m *Manager) GenerateGitignore(config *ui.ProjectConfig) string {
	return `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
build/

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out
coverage.html

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Application specific
config/local.yml
*.log
tmp/
temp/
.env
.env.local`
}

// GenerateDockerfile creates the Dockerfile
func (m *Manager) GenerateDockerfile(config *ui.ProjectConfig) string {
	return fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o %s .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/%s .
COPY --from=builder /app/config ./config

CMD ["./%s"]`, config.ProjectName, config.ProjectName, config.ProjectName)
}

// GenerateTestFile creates the test file
func (m *Manager) GenerateTestFile(config *ui.ProjectConfig) string {
	return `package test

import (
	"context"
	"testing"
	"time"
)

func TestApplicationBasics(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{"Context handling", testContextHandling},
		{"Configuration loading", testConfigurationLoading},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func testContextHandling(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	if ctx.Err() != nil {
		t.Errorf("Context should not be cancelled initially")
	}
	
	// Test context cancellation
	cancel()
	if ctx.Err() == nil {
		t.Errorf("Context should be cancelled after cancel()")
	}
}

func testConfigurationLoading(t *testing.T) {
	// Test basic configuration loading
	// This would test the actual configuration loading logic
	// when the config package is integrated
	t.Log("Configuration loading test placeholder")
}

// BenchmarkApplicationStartup benchmarks application startup time
func BenchmarkApplicationStartup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		// Benchmark application initialization
		_ = ctx
	}
}`
}

// Helper methods

func (m *Manager) containsPackage(packages []string, target string) bool {
	for _, pkg := range packages {
		if pkg == target {
			return true
		}
	}
	return false
}

func (m *Manager) generateCommandHelp(config *ui.ProjectConfig) string {
	var commands []string

	if m.containsPackage(config.SelectedPackages, "fileops") {
		commands = append(commands, "\n- **create**: Create files with example content")
	}

	if m.containsPackage(config.SelectedPackages, "processing") {
		commands = append(commands, "\n- **process**: Demonstrate concurrent processing")
	}

	return strings.Join(commands, "")
}

func (m *Manager) generateStructureExtras(config *ui.ProjectConfig) string {
	var extras []string

	if m.containsPackage(config.SelectedPackages, "config") {
		extras = append(extras, "\n│   └── config.go      # Configuration command")
	}

	if m.containsPackage(config.SelectedPackages, "fileops") {
		extras = append(extras, "\n│   └── create.go      # File creation command")
	}

	if m.containsPackage(config.SelectedPackages, "processing") {
		extras = append(extras, "\n│   └── process.go     # Processing command")
	}

	return strings.Join(extras, "")
}

func (m *Manager) generatePackageConfig(config *ui.ProjectConfig) string {
	var configs []string

	if m.containsPackage(config.SelectedPackages, "formatter") {
		configs = append(configs, "\nformatter:\n  template_dir: \"./templates\"\n  default_format: \"text\"")
	}

	if m.containsPackage(config.SelectedPackages, "git") {
		configs = append(configs, "\ngit:\n  auto_detect: true\n  include_branch: true\n  include_hash: true")
	}

	if m.containsPackage(config.SelectedPackages, "fileops") {
		configs = append(configs, "\nfileops:\n  backup_dir: \"./backups\"\n  preserve_permissions: true")
	}

	if m.containsPackage(config.SelectedPackages, "processing") {
		configs = append(configs, "\nprocessing:\n  worker_count: 4\n  buffer_size: 100\n  timeout: \"30s\"")
	}

	if len(configs) == 0 {
		return "\n# Add package-specific configuration here"
	}

	return strings.Join(configs, "")
}
