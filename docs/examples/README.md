# Integration Examples

// EXTRACT-010: Integration examples overview - Package integration examples and patterns - 🔺

This directory contains comprehensive integration examples demonstrating how to use the extracted BkpDir packages together to build real-world applications.

## Available Examples

### 1. Basic CLI Application (`basic-cli-app/`)

**Demonstrates**: `pkg/config` + `pkg/cli` + `pkg/formatter` + `pkg/errors`

A complete CLI application showing:
- Multi-source configuration management (files, environment, defaults)
- CLI framework with commands, flags, and context support
- Multiple output formats (table, JSON, YAML, templates)
- Structured error handling with categories and context
- Dry-run mode and verbose output

**Key Features**:
- Configuration-driven behavior
- Context-aware execution
- Format-agnostic output
- Error handling patterns

### 2. Git-Aware Backup Tool (`git-aware-backup/`)

**Demonstrates**: `pkg/config` + `pkg/cli` + `pkg/git` + `pkg/fileops` + `pkg/processing` + `pkg/formatter`

A sophisticated backup tool featuring:
- Git repository detection and metadata extraction
- Git-aware naming conventions with timestamps
- File operations with pattern matching and exclusions
- Progress reporting and multi-format output
- Configuration-driven file operations

**Key Features**:
- Git integration patterns
- File operation patterns
- Naming conventions with metadata
- Progress tracking
- Cross-package integration

## Integration Patterns Demonstrated

### Configuration Management
- Loading from multiple sources (defaults → environment → files)
- Type-safe configuration access
- Validation and error handling
- Environment variable mapping

### CLI Framework Integration
- Command structure with subcommands
- Global and command-specific flags
- Context propagation throughout the application
- Dry-run mode support
- Version command integration

### Output Formatting
- Multiple format support (table, JSON, YAML, templates)
- Format-agnostic data structures
- Template-based custom formatting
- Progress reporting and user feedback

### Error Handling
- Structured error types with categories
- Error context and metadata
- Recoverable vs non-recoverable errors
- Consistent error reporting across packages

### File Operations
- Pattern-based inclusion/exclusion
- Atomic operations for data integrity
- Path validation and creation
- Cross-platform compatibility

### Git Integration
- Repository detection and validation
- Metadata extraction (branch, hash, status)
- Git-aware naming and operations
- Support for both Git and non-Git projects

## Usage Guidelines

### Running Examples

These are **conceptual examples** that demonstrate integration patterns. The code shows intended API usage but may require adjustments to match actual package interfaces.

```bash
# Navigate to an example directory
cd basic-cli-app/

# Review the README for specific usage instructions
cat README.md

# Examine the example code
cat example-code/main.go

# Use as a starting point for your own project
cp example-code/main.go your-project/
```

### Adapting Examples

When using these examples as templates:

1. **Update import paths** to match your module structure
2. **Verify package APIs** against current interfaces
3. **Adjust configuration structures** to match your needs
4. **Customize error handling** for your use cases
5. **Add tests** for your specific functionality

## Package Integration Benefits

### Modularity
- Each package can be used independently
- Flexible composition for different use cases
- Clear separation of concerns

### Consistency
- Similar patterns across all packages
- Consistent error handling
- Uniform configuration management

### Testability
- Interface-based design enables mocking
- Independent testing of components
- Integration testing patterns

### Extensibility
- Easy to add new output formats
- Pluggable configuration sources
- Extensible command structures
- Custom error types and handlers

## Best Practices Demonstrated

### 1. Configuration Management
```go
// Multi-source configuration loading
sources := []config.ConfigSource{
    config.NewMapConfigSource("defaults", defaults),
    config.NewEnvConfigSource("MYAPP"),
    config.NewFileConfigSource("config.yaml"),
}
cfg, err := loader.LoadConfig(sources...)
```

### 2. Context Propagation
```go
// CLI context flows through the application
cmdCtx := &cli.CommandContext{
    Context: ctx,
    DryRun:  dryRun,
    Verbose: verbose,
}
```

### 3. Error Handling
```go
// Structured errors with context
return errors.NewApplicationError(
    errors.CategoryFilesystem,
    errors.SeverityError,
    "FILE_NOT_FOUND",
    "Source file does not exist",
    map[string]interface{}{"path": sourcePath},
)
```

### 4. Format-Agnostic Output
```go
// Same data, multiple formats
switch outputFormat {
case "json":
    return formatter.FormatJSON(data)
case "yaml":
    return formatter.FormatYAML(data)
case "table":
    return formatter.FormatTable(headers, rows)
}
```

### 5. Resource Management
```go
// Proper resource cleanup
defer func() {
    if err := resources.Cleanup(); err != nil {
        log.Printf("Cleanup failed: %v", err)
    }
}()
```

## Development Workflow

### Creating New Examples

1. **Identify integration scenario** - Which packages work together?
2. **Create example directory** - Follow naming convention
3. **Write example code** - Focus on integration patterns
4. **Document usage patterns** - Comprehensive README
5. **Add configuration examples** - Show real-world configs
6. **Include test patterns** - How to test integrated functionality

### Example Structure

```
my-example/
├── README.md              # Comprehensive documentation
├── config.yaml           # Example configuration
├── example-code/
│   ├── go.mod            # Separate module for example
│   └── main.go           # Integration example code
└── .testignore           # Exclude from main test suite
```

These examples serve as both documentation and starting points for building applications with the extracted BkpDir packages, demonstrating real-world integration patterns and best practices. 