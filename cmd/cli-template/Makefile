# ⭐ EXTRACT-008: CLI Application Template - Build automation and development workflow
.PHONY: help build clean test lint run install demo examples

# Default target
help: ## Show this help message
	@echo "CLI Template Application - Build Targets"
	@echo "========================================"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the CLI template application
	@echo "Building CLI template application..."
	go build -o cli-template .
	@echo "✅ Build complete: ./cli-template"

build-all: ## Build for multiple platforms
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -o cli-template-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o cli-template-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o cli-template-windows-amd64.exe .
	@echo "✅ Multi-platform build complete"

# Development targets
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f cli-template cli-template-* example.txt demo.txt test.txt
	@echo "✅ Clean complete"

test: ## Run tests (when available)
	@echo "Running tests..."
	go test -v ./...
	@echo "✅ Tests complete"

lint: ## Run linting checks
	@echo "Running linting checks..."
	go vet ./...
	go fmt ./...
	@echo "✅ Linting complete"

# Runtime targets
run: build ## Build and run the application
	@echo "Running CLI template application..."
	./cli-template

run-help: build ## Show application help
	./cli-template --help

run-version: build ## Show version information
	./cli-template version

run-config: build ## Show configuration information
	./cli-template config

# Demo targets
demo: build ## Run a complete demo of all features
	@echo "CLI Template Application Demo"
	@echo "============================"
	@echo ""
	@echo "1. Showing help..."
	./cli-template --help
	@echo ""
	@echo "2. Showing version..."
	./cli-template version
	@echo ""
	@echo "3. Showing configuration..."
	./cli-template config
	@echo ""
	@echo "4. Creating example file..."
	./cli-template create demo.txt
	@echo ""
	@echo "5. Processing files (simulation)..."
	./cli-template process demo.txt example.txt
	@echo ""
	@echo "✅ Demo complete!"

demo-verbose: build ## Run demo with verbose output
	@echo "CLI Template Application - Verbose Demo"
	@echo "======================================"
	@echo ""
	./cli-template --verbose config
	@echo ""
	./cli-template --verbose create verbose-demo.txt
	@echo ""
	@echo "✅ Verbose demo complete!"

demo-dry-run: build ## Run demo in dry-run mode
	@echo "CLI Template Application - Dry Run Demo"
	@echo "======================================"
	@echo ""
	./cli-template --dry-run create dry-run-test.txt
	@echo ""
	@echo "✅ Dry run demo complete!"

# Example targets
examples: build ## Run all examples
	@echo "Running all examples..."
	@echo ""
	@echo "Example 1: Basic usage"
	./cli-template --help
	@echo ""
	@echo "Example 2: Configuration management"
	./cli-template config
	@echo ""
	@echo "Example 3: File operations"
	./cli-template create example.txt
	@echo ""
	@echo "Example 4: Concurrent processing"
	./cli-template process example.txt
	@echo ""
	@echo "✅ All examples complete!"

# Installation targets
install: build ## Install the application to GOPATH/bin
	@echo "Installing CLI template application..."
	go install .
	@echo "✅ Installation complete"

# Development workflow
dev: clean lint test build ## Complete development workflow
	@echo "✅ Development workflow complete"

# Quality checks
check: lint test ## Run quality checks
	@echo "✅ Quality checks complete"

# Package integration demos
demo-config: build ## Demo configuration package integration
	@echo "Configuration Package Demo (pkg/config)"
	@echo "======================================"
	./cli-template config --verbose

demo-fileops: build ## Demo file operations package integration
	@echo "File Operations Package Demo (pkg/fileops)"
	@echo "=========================================="
	./cli-template create fileops-demo.txt --verbose

demo-processing: build ## Demo concurrent processing package integration
	@echo "Concurrent Processing Package Demo (pkg/processing)"
	@echo "================================================="
	./cli-template process file1.txt file2.txt file3.txt

demo-git: build ## Demo Git integration package
	@echo "Git Integration Package Demo (pkg/git)"
	@echo "====================================="
	./cli-template version

# Documentation targets
docs: ## Generate documentation (placeholder)
	@echo "Documentation generation..."
	@echo "📚 README.md - Application documentation"
	@echo "📚 Package integration examples in commands"
	@echo "📚 Usage examples in Makefile targets"
	@echo "✅ Documentation ready"

# Cleanup targets
clean-examples: ## Clean example files created during demos
	@echo "Cleaning example files..."
	rm -f demo.txt example.txt verbose-demo.txt fileops-demo.txt
	@echo "✅ Example files cleaned"

# All-in-one targets
all: clean lint test build examples ## Run complete build and test cycle
	@echo "✅ Complete build cycle finished!"

quick: build run-help ## Quick build and help display
	@echo "✅ Quick start complete!" 