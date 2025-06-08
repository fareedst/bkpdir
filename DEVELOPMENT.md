# Development Guide for bkpdir

This guide explains how to build, test, and develop the `bkpdir` CLI application.

## Prerequisites

- Go 1.19 or later
- Make (for using the Makefile)
- Git (for version information in builds)

## Quick Start

```bash
# Build and test everything
make

# Build for local development
make build-local

# Run all tests
make test

# Show all available targets
make help
```

## Development Workflow

### 1. Building

```bash
# Build for local development (current platform)
make build-local

# Build and run basic functionality test
make dev

# Install to ~/.local/bin (make sure it's in your PATH)
make install
```

### 2. Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report (generates coverage.html)
make test-coverage

# Run tests with race detection
make test-race

# Run benchmark tests
make test-bench

# Run all test variants
make test-all
```

### 3. Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run linter (requires revive: go install github.com/mgechev/revive@latest)
make lint

# Run all code quality checks
make check
```

### 4. Production Builds

```bash
# Build for all platforms
make build-all

# Build for macOS (both ARM64 and Intel)
make build-macos

# Build for Ubuntu (20.04, 22.04, 24.04)
make build-ubuntu

# Build for specific platforms
make build-macos-arm64
make build-macos-amd64
make build-ubuntu20
make build-ubuntu22
make build-ubuntu24
```

### 5. Utilities

```bash
# Clean build artifacts and test cache
make clean

# Download and verify dependencies
make deps

# Show help
make help
```

## Testing Strategy

The project includes several types of tests:

1. **Unit Tests**: Test individual functions and components
2. **Integration Tests**: Test command-line interface and workflows
3. **Coverage Tests**: Ensure adequate test coverage
4. **Race Tests**: Detect race conditions in concurrent code

### Running Specific Tests

```bash
# Run tests for a specific file
go test -v -run TestVerifyArchive

# Run tests with coverage for specific packages
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Code Quality Standards

The project uses several tools to maintain code quality:

- **gofmt**: Code formatting
- **go vet**: Static analysis
- **revive**: Comprehensive linting with custom rules (see `.revive.toml`)

### Linting Configuration

The project includes a comprehensive `.revive.toml` configuration that enforces:

- Cyclomatic complexity limits
- Line length limits (120 characters)
- Proper documentation for exported functions
- Error handling best practices
- File header requirements

## Build Configuration

The Makefile includes several build-time variables:

- `compileDate`: Set to current UTC timestamp
- `platform`: Set to target OS and architecture
- `version`: Application version (currently 1.5.0)

These are injected using Go's `-ldflags` during compilation.

## Directory Structure

```
.
├── main.go              # Main application entry point
├── archive.go           # Archive creation and management
├── verify.go            # Archive verification
├── config.go            # Configuration management
├── formatter.go         # Output formatting
├── errors.go            # Error handling and resource management
├── comparison.go        # Directory comparison
├── exclude.go           # File exclusion patterns
├── git.go              # Git integration
├── *_test.go           # Test files
├── .bkpdir.yml         # Default configuration
├── .revive.toml        # Linting configuration
├── Makefile            # Build automation
└── docs/               # Documentation
```

## Continuous Integration

The Makefile is designed to support CI/CD workflows:

```bash
# Typical CI pipeline
make deps     # Download dependencies
make check    # Code quality checks
make test-all # All test variants
make build-all # Build for all platforms
```

## Troubleshooting

### Common Issues

1. **Tests failing**: Run `make clean` then `make test`
2. **Linter not found**: Install revive with `go install github.com/mgechev/revive@latest`
3. **Build failures**: Ensure Go 1.19+ is installed and `GOPATH` is set correctly

### Debug Builds

For debugging, you can build with additional flags:

```bash
go build -gcflags="all=-N -l" -o bkpdir-debug
```

### Verbose Output

Most Makefile targets support verbose output. For more detailed information:

```bash
make test-verbose
make build-local VERBOSE=1
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `make check test` to ensure quality
5. Submit a pull request

Make sure all tests pass and code quality checks succeed before submitting. 