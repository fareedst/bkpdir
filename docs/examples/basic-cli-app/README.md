# Basic CLI Application Example

// EXTRACT-010: Basic CLI application integration example - Demonstrates config + cli + formatter package integration - 🔺

This example demonstrates how to build a CLI application using the extracted BkpDir packages:

- **pkg/config**: Configuration management from multiple sources
- **pkg/cli**: CLI framework with dry-run and context support
- **pkg/formatter**: Output formatting with templates and multiple formats
- **pkg/errors**: Structured error handling

## Features Demonstrated

### Configuration Management
- Loading configuration from multiple sources (defaults, environment, files)
- Schema-agnostic configuration handling
- Configuration validation
- Environment variable mapping with prefixes

### CLI Framework
- Command structure with subcommands
- Global and command-specific flags
- Dry-run mode integration
- Context-aware execution
- Version command integration

### Output Formatting
- Multiple output formats (table, JSON, YAML, templates)
- Template-based custom formatting
- Structured data formatting
- Format-specific error handling

### Error Handling
- Structured error types with categories
- Error context and metadata
- Recoverable vs non-recoverable errors
- Application-specific error codes

## Installation

This is a **conceptual example** that demonstrates integration patterns. The code in `example-code/main.go` shows the intended API usage patterns but may require adjustments to match the actual package interfaces.

To use these patterns in your own project:

```bash
# Copy the example code as a starting point
cp example-code/main.go your-project/
# Update import paths and adapt to actual package APIs
go mod tidy
```

## Configuration

The application loads configuration from multiple sources in priority order:

1. **Default values** (hardcoded in the application)
2. **Environment variables** (prefixed with `BASICAPP_`)
3. **Configuration file** (`config.yaml`)

### Environment Variables

```bash
export BASICAPP_OUTPUT_FORMAT=json
export BASICAPP_VERBOSE=true
export BASICAPP_DRY_RUN=false
export BASICAPP_PROCESSING_CONCURRENCY=8
```

### Configuration File

Create a `config.yaml` file:

```yaml
# Output configuration
output_format: table
verbose: false
dry_run: false

# Processing configuration
processing:
  batch_size: 150
  timeout: 15m
  concurrency: 6

# Output templates
output:
  directory: ./results
  template: "{{.name}}: {{.status}} (took {{.duration}})"
```

## Usage Examples

### Basic Command Help

```bash
./basic-cli-app --help
```

### Process Command with Different Formats

```bash
# Table format (default)
./basic-cli-app process --input=/data --output=/results

# JSON format
./basic-cli-app process --format=json --verbose

# YAML format with dry-run
./basic-cli-app process --format=yaml --dry-run

# Custom template format
./basic-cli-app process --format=template
```

### Status Command

```bash
# Show application status
./basic-cli-app status

# Status in JSON format
./basic-cli-app status --format=json

# Status with verbose output
./basic-cli-app status --verbose
```

### Configuration Management

```bash
# Show current configuration
./basic-cli-app config show

# Validate configuration
./basic-cli-app config validate
```

### Version Information

```bash
./basic-cli-app version
```

## Example Output

### Table Format Output

```
Name      Size      Status     Duration  Timestamp
------    ------    --------   --------  ---------
file1.txt 1024 bytes processed 15ms      2024-01-02T12:00:00Z
file2.txt 2048 bytes processed 23ms      2024-01-02T12:00:01Z
file3.txt 512 bytes  skipped    5ms       2024-01-02T12:00:02Z

Processing Summary:
- Input Path: /data
- Output Path: /results
- Processed Items: 3
- Dry Run: false
- Format: table
```

### JSON Format Output

```json
[
  {
    "name": "file1.txt",
    "size": 1024,
    "status": "processed",
    "duration": "15ms",
    "timestamp": "2024-01-02T12:00:00Z"
  },
  {
    "name": "file2.txt",
    "size": 2048,
    "status": "processed",
    "duration": "23ms",
    "timestamp": "2024-01-02T12:00:01Z"
  },
  {
    "name": "file3.txt",
    "size": 512,
    "status": "skipped",
    "duration": "5ms",
    "timestamp": "2024-01-02T12:00:02Z"
  }
]
```

### Status Output

```
Application Status
==================
Name:         basic-cli-app
Version:      1.0.0
Status:       running
Uptime:       45 minutes

Configuration
=============
Format:       table
Verbose:      false
Dry Run:      false
Config File:  config.yaml

Statistics
==========
Processed:    156
Errors:       3
Warnings:     12
Last Run:     2024-01-02T11:00:00Z
```

## Integration Patterns Demonstrated

### 1. Configuration-Driven Behavior

The application shows how configuration affects behavior across all packages:

```go
// Configuration affects output formatting
outputFormat := cfg.GetString("output_format", "table")

// Configuration affects processing parameters
concurrency := cfg.GetInt("processing.concurrency", 4)
timeout := cfg.GetDuration("processing.timeout", time.Minute*10)
```

### 2. CLI Context Propagation

Context flows through the entire application:

```go
// CLI context creation
cmdCtx := &cli.CommandContext{
    Context: ctx,
    DryRun:  false,
    Verbose: false,
}

// Context usage in processing
if cmdCtx.DryRun {
    fmt.Println("DRY RUN MODE - No changes will be made")
}
```

### 3. Structured Error Handling

Errors are handled consistently across packages:

```go
return errors.NewApplicationError(
    errors.CategoryConfiguration,
    errors.SeverityError,
    "CONFIG_LOAD_FAILED",
    fmt.Sprintf("Failed to load configuration: %v", err),
    map[string]interface{}{"command": "process"},
)
```

### 4. Format-Agnostic Output

The same data can be formatted in multiple ways:

```go
switch outputFormat {
case "json":
    jsonOutput, err := app.formatter.FormatJSON(data)
case "yaml":
    yamlOutput, err := app.formatter.FormatYAML(data)
case "table":
    tableOutput, err := app.formatter.FormatTable(headers, rows)
case "template":
    output, err := app.formatter.FormatTemplate(template, data)
}
```

## Package Integration Benefits

### Modularity
Each package can be used independently or together, allowing for flexible application architecture.

### Consistency
All packages follow similar patterns for configuration, error handling, and context management.

### Testability
Interface-based design makes it easy to mock and test individual components.

### Extensibility
New output formats, configuration sources, and CLI commands can be added easily.

## Building and Running

```bash
# Build the application
go build -o basic-cli-app main.go

# Run with default configuration
./basic-cli-app process

# Run with custom configuration
./basic-cli-app --config=myconfig.yaml process --format=json

# Run in dry-run mode
./basic-cli-app process --dry-run --verbose
```

## Testing

The example includes patterns for testing integrated functionality:

```bash
# Test configuration loading
./basic-cli-app config validate

# Test different output formats
./basic-cli-app status --format=json
./basic-cli-app status --format=yaml

# Test error handling
./basic-cli-app process --input=/nonexistent
```

This example serves as a template for building CLI applications using the extracted BkpDir packages, demonstrating best practices for package integration and real-world usage patterns. 