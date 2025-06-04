# Makefile for CLI Project Scaffolding Generator

.PHONY: help build test clean demo run install

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the scaffolding application
	@echo "🔨 Building scaffolding generator..."
	@mkdir -p bin
	@go build -o bin/scaffolding .
	@echo "✅ Build completed: bin/scaffolding"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

run: build ## Build and run interactively
	@echo "🚀 Starting scaffolding generator..."
	@./bin/scaffolding

demo: build ## Build and run demonstration
	@echo "🚀 Running scaffolding demonstration..."
	@echo "This will create a demo project in /tmp/demo-cli-app"
	@rm -rf /tmp/demo-cli-app
	@mkdir -p /tmp/demo-output
	@echo "my-demo-app" | \
	 echo "/tmp/demo-output/my-demo-app" | \
	 echo "github.com/demo/my-demo-app" | \
	 echo "Demo Author" | \
	 echo "A demonstration CLI application" | \
	 echo "standard" | \
	 echo "y" | \
	 echo "n" | \
	 echo "y" | \
	 ./bin/scaffolding || echo "Interactive demo - run 'make run' for full experience"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

install: build ## Install the scaffolding generator globally
	@echo "📦 Installing scaffolding generator..."
	@go install .
	@echo "✅ Installation completed - use 'scaffolding' command globally"

# Development targets
lint: ## Run linting
	@echo "🔍 Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run

format: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Generate a test project for validation
test-generate: build ## Generate a test project for validation
	@echo "🧪 Generating test project..."
	@rm -rf /tmp/test-cli-project
	@mkdir -p /tmp/test-outputs
	@echo "Creating standard template project..."
	@echo "test-app\n/tmp/test-outputs/test-app\ngithub.com/test/test-app\nTest Author\nA test CLI application\n2\ny\nn\ny" | ./bin/scaffolding || true
	@if [ -d "/tmp/test-outputs/test-app" ]; then \
		echo "✅ Test project generated successfully"; \
		echo "📁 Testing build in generated project..."; \
		cd /tmp/test-outputs/test-app && make build && echo "✅ Generated project builds successfully"; \
	else \
		echo "❌ Test project generation failed"; \
	fi 