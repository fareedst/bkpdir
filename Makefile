# Makefile for bkpdir - Directory archiving CLI
# 
# Available targets:
#   Development:
#     build-local     - Build for local development (current platform)
#     dev             - Build and run basic functionality test
#     install         - Install to ~/.local/bin (requires build-local first)
#   
#   Testing:
#     test            - Run all tests
#     test-verbose    - Run tests with verbose output
#     test-coverage   - Run tests with coverage report
#     test-race       - Run tests with race detection
#     test-bench      - Run benchmark tests
#     test-all        - Run all test variants
#   
#   Code Quality:
#     lint            - Run linter (revive)
#     fmt             - Format code with gofmt
#     vet             - Run go vet
#     check           - Run all code quality checks
#   
#   Production Builds:
#     build-all       - Build for all platforms
#     build-macos     - Build for macOS (ARM64 and AMD64)
#     build-ubuntu    - Build for Ubuntu (20.04, 22.04, 24.04)
#   
#   Utilities:
#     clean           - Clean build artifacts and test cache
#     deps            - Download and verify dependencies
#     help            - Show this help message

.PHONY: build-all build-ubuntu20 build-ubuntu22 build-ubuntu24 build-macos build-macos-arm64 build-macos-amd64 build-local clean
.PHONY: test test-verbose test-coverage test-race test-bench test-all
.PHONY: lint fmt vet check dev install deps help

# Variables
BINARY_NAME := bkpdir
VERSION := 1.3.0
BUILD_TIME := $(shell date -u +%Y-%m-%d\ %H:%M:%S\ UTC)
PLATFORM := $(shell go env GOOS)-$(shell go env GOARCH)
LDFLAGS := -X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=$(PLATFORM)'

# Default target
all: check test build-local

# Help target
help:
	@echo "bkpdir Makefile - Directory archiving CLI"
	@echo ""
	@echo "Development targets:"
	@echo "  build-local     Build for local development (current platform)"
	@echo "  dev             Build and run basic functionality test"
	@echo "  install         Install to ~/.local/bin"
	@echo ""
	@echo "Testing targets:"
	@echo "  test            Run all tests"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  test-race       Run tests with race detection"
	@echo "  test-bench      Run benchmark tests"
	@echo "  test-all        Run all test variants"
	@echo ""
	@echo "Code quality targets:"
	@echo "  lint            Run linter (revive)"
	@echo "  fmt             Format code with gofmt"
	@echo "  vet             Run go vet"
	@echo "  check           Run all code quality checks"
	@echo ""
	@echo "Production build targets:"
	@echo "  build-all       Build for all platforms"
	@echo "  build-macos     Build for macOS (ARM64 and AMD64)"
	@echo "  build-ubuntu    Build for Ubuntu (20.04, 22.04, 24.04)"
	@echo ""
	@echo "Utility targets:"
	@echo "  clean           Clean build artifacts and test cache"
	@echo "  deps            Download and verify dependencies"
	@echo "  help            Show this help message"

# Development targets
build-local:
	@echo "Building $(BINARY_NAME) for local development..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)
	@echo "✓ Built $(BINARY_NAME) successfully"

dev: build-local
	@echo "Running basic functionality test..."
	@./$(BINARY_NAME) --help > /dev/null && echo "✓ Help command works"
	@./$(BINARY_NAME) --config > /dev/null && echo "✓ Config command works"
	@./$(BINARY_NAME) full --dry-run > /dev/null && echo "✓ Dry-run works"
	@echo "✓ Basic functionality test passed"

install: build-local
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp $(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME) to ~/.local/bin/$(BINARY_NAME)"
	@echo "  Make sure ~/.local/bin is in your PATH"

# Testing targets
test:
	@echo "Running tests..."
	go test ./...
	@echo "✓ All tests passed"

test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	@echo "Generating detailed coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -race ./...
	@echo "✓ Race detection tests passed"

test-bench:
	@echo "Running benchmark tests..."
	go test -bench=. ./...

test-all: test test-race test-coverage
	@echo "✓ All test variants completed"

# Code quality targets
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ go vet passed"

lint:
	@echo "Running linter..."
	@if command -v revive >/dev/null 2>&1; then \
		revive -config .revive.toml ./...; \
		echo "✓ Linter passed"; \
	else \
		echo "⚠ revive not installed, skipping lint check"; \
		echo "  Install with: go install github.com/mgechev/revive@latest"; \
	fi

check: fmt vet lint
	@echo "✓ All code quality checks passed"

# Production build targets
build-all: build-local build-macos build-ubuntu
	@echo "✓ Built for all platforms"

build-ubuntu: build-ubuntu20 build-ubuntu22 build-ubuntu24
	@echo "✓ Built for all Ubuntu versions"

build-ubuntu20:
	@echo "Building for Ubuntu 20.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u20'" -o bin/$(BINARY_NAME)-ubuntu20.04

build-ubuntu22:
	@echo "Building for Ubuntu 22.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u22'" -o bin/$(BINARY_NAME)-ubuntu22.04

build-ubuntu24:
	@echo "Building for Ubuntu 24.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u24'" -o bin/$(BINARY_NAME)-ubuntu24.04

build-macos: build-macos-arm64 build-macos-amd64
	@echo "✓ Built for all macOS architectures"

build-macos-arm64:
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=darwin-arm64'" -o bin/$(BINARY_NAME)-macos-arm64

build-macos-amd64:
	@echo "Building for macOS AMD64 (Intel)..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=darwin-amd64'" -o bin/$(BINARY_NAME)-macos-amd64

# Utility targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -testcache
	@echo "✓ Cleaned build artifacts and test cache"

deps:
	@echo "Downloading and verifying dependencies..."
	go mod download
	go mod verify
	@echo "✓ Dependencies verified"

# Legacy targets for compatibility
build: build-local

# Example usage comment for symbolic link
# To create a symbolic link in your local bin:
# ln -s "$(pwd)/bin/$(BINARY_NAME)-macos-arm64" ~/.local/bin/$(BINARY_NAME)
