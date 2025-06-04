# Basic Usage Examples

This directory contains basic examples of using the CLI template application.

## 🚀 Getting Started

### 1. Build the Application

```bash
cd cmd/cli-template
make build
```

### 2. Run Basic Commands

```bash
# Show help
./cli-template --help

# Show version
./cli-template version

# Show configuration
./cli-template config
```

## 📝 Basic Examples

### Example 1: Create a File

```bash
# Create a simple file
./cli-template create hello.txt

# Create with verbose output
./cli-template --verbose create detailed.txt

# Dry run (show what would happen)
./cli-template --dry-run create test.txt
```

### Example 2: Process Files

```bash
# Process a single file
./cli-template process hello.txt

# Process multiple files
./cli-template process file1.txt file2.txt file3.txt

# Process with verbose output
./cli-template --verbose process *.txt
```

### Example 3: Configuration Management

```bash
# Show current configuration
./cli-template config

# Use custom config file
./cli-template --config ./config/example.yml config

# Show verbose configuration details
./cli-template --verbose config
```

### Example 4: Using Global Flags

```bash
# Verbose mode
./cli-template --verbose version

# Dry run mode
./cli-template --dry-run create test-file.txt

# Custom config file
./cli-template --config custom.yml config
```

## 🎯 Package Demonstrations

Each command demonstrates different extracted packages:

- **version**: pkg/git integration
- **config**: pkg/config management
- **create**: pkg/fileops, pkg/errors, pkg/resources
- **process**: pkg/processing, concurrent patterns

## 📚 Next Steps

- See [Advanced Examples](../advanced/) for complex usage patterns
- See [Integration Examples](../integration/) for package integration patterns
- Read the main [README](../../README.md) for complete documentation 