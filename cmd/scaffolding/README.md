# 🚀 Go CLI Project Scaffolding Generator

Interactive CLI tool for generating new Go CLI applications using extracted packages from the bkpdir project.

## 🎯 Purpose

This scaffolding generator demonstrates the value of extracted packages by creating new CLI projects with:
- **Interactive configuration** - User-friendly prompts for project setup
- **Multiple templates** - Minimal, standard, advanced, and custom configurations
- **Package selection** - Choose which extracted packages to include
- **Complete project generation** - Ready-to-build CLI applications with documentation

## 📦 Available Packages

| Package | Description | Required |
|---------|-------------|----------|
| **config** | Configuration management - schema-agnostic YAML/JSON loading | ✅ Yes |
| **errors** | Structured error handling with status codes | ✅ Yes |
| **cli** | CLI framework and command patterns | ✅ Yes |
| **resources** | Resource management with cleanup coordination | ❌ Optional |
| **formatter** | Template-based and printf-style output formatting | ❌ Optional |
| **git** | Git repository detection and information extraction | ❌ Optional |
| **fileops** | Safe file operations and comparisons | ❌ Optional |
| **processing** | Concurrent processing with worker pools | ❌ Optional |

## 🏗️ Project Templates

### 1. Minimal Template
- **Packages**: config, errors, cli
- **Features**: Basic CLI with configuration and error handling
- **Use Case**: Simple command-line tools

### 2. Standard Template  
- **Packages**: config, errors, cli, formatter, git
- **Features**: Enhanced output formatting and Git integration
- **Use Case**: Most common CLI applications

### 3. Advanced Template
- **Packages**: All 8 packages included
- **Features**: Full feature set with file operations and concurrent processing
- **Use Case**: Complex CLI applications with advanced requirements

### 4. Custom Template
- **Packages**: User-selected combination
- **Features**: Tailored to specific needs
- **Use Case**: Specific requirements or learning individual packages

## 🚀 Quick Start

### Installation

```bash
# Build the scaffolding generator
make build

# Or install globally
make install
```

### Usage

```bash
# Interactive mode (recommended)
make run

# Or run directly
./bin/scaffolding
```

### Example Session

```
🚀 Go CLI Project Scaffolding Generator
=====================================

✔ Project name: my-awesome-cli
✔ Output directory: /Users/me/projects/my-awesome-cli
✔ Go module name: github.com/me/my-awesome-cli
✔ Author name: Your Name
✔ Project description: An awesome CLI tool
✔ Select project template: standard
📦 Selected packages: config, errors, cli, formatter, git
✔ Include example tests: Yes
✔ Include Dockerfile: No
✔ Include .gitignore: Yes

📁 Project structure created:
my-awesome-cli/
├── main.go
├── go.mod
├── Makefile
├── README.md
├── .gitignore
├── cmd/
│   ├── root.go
│   ├── version.go
│   └── config.go
└── config/
    ├── default.yml
    └── example.yml

🔗 Included packages: config, errors, cli, formatter, git

✅ Successfully generated project 'my-awesome-cli'
📁 Project location: /Users/me/projects/my-awesome-cli

🎯 Next steps:
   cd my-awesome-cli
   make build    # Build the application
   make demo     # Run demonstration
   make help     # See all available commands
```

## 🏗️ Generated Project Structure

```
my-project/
├── main.go                 # Application entry point with context handling
├── go.mod                  # Go module with selected dependencies
├── Makefile               # Build automation (build, test, demo, lint)
├── README.md              # Generated documentation
├── .gitignore             # Standard Go .gitignore (optional)
├── Dockerfile             # Multi-stage Docker build (optional)
├── cmd/                   # Command structure
│   ├── root.go           # Root command with flags
│   ├── version.go        # Version command
│   ├── config.go         # Config command (if config package selected)
│   ├── create.go         # File creation command (if fileops selected)
│   └── process.go        # Processing command (if processing selected)
├── config/               # Configuration files (if config package selected)
│   ├── default.yml       # Default settings
│   └── example.yml       # Example configuration
└── test/                 # Test files (optional)
    └── main_test.go      # Basic application tests
```

## [ACTION:maintenance] Development

### Build and Test

```bash
# Build the scaffolding generator
make build

# Run tests
make test

# Generate a test project and validate it builds
make test-generate

# Run linting
make lint

# Format code
make format
```

### Creating Templates

The template system uses Go's text templates with package configuration:

1. **Template Manager** (`internal/templates/manager.go`): Generates all file content
2. **UI System** (`internal/ui/config.go`): Interactive configuration collection
3. **Generator** (`internal/generator/generator.go`): Orchestrates project creation

### Adding New Packages

To add support for a new package:

1. Add package definition to `ui.AvailablePackages`
2. Update template methods to include package-specific content
3. Add package-specific commands if needed
4. Update documentation and examples

## 📋 Generated Project Features

### Make Targets
Every generated project includes comprehensive build automation:

```bash
make help          # Show available targets
make build         # Build the application
make test          # Run tests
make clean         # Clean build artifacts
make demo          # Run demonstration
make install       # Install globally
make lint          # Run code linting
make dev           # Development mode with file watching
```

### Configuration System
Generated projects include configuration management:

- YAML-based configuration files
- Environment variable support
- Configuration file discovery
- Verbose configuration display

### Command Structure
CLI applications follow consistent patterns:

- Context-aware command execution
- Graceful interrupt handling
- Consistent flag patterns (--verbose, --dry-run, --config)
- Help and version commands

### Error Handling
Structured error handling throughout:

- Consistent error formatting
- Resource cleanup on errors
- Context cancellation support
- Exit code management

## 🎯 Integration Examples

### Basic CLI Application
```bash
# Generate minimal project
./bin/scaffolding
# Select: minimal template
# Result: Basic CLI with config and error handling
```

### File Processing Tool
```bash
# Generate with fileops and processing
./bin/scaffolding  
# Select: custom template → include fileops, processing
# Result: CLI with file operations and concurrent processing
```

### Git-Aware Application
```bash
# Generate with Git integration
./bin/scaffolding
# Select: standard template (includes git package)
# Result: CLI with Git branch/commit detection
```

## [ACTION:migration] [HIGH] EXTRACT-008: Project scaffolding system [ACTION:maintenance]

This scaffolding system demonstrates:

1. **Package Value** - Shows immediate utility of extracted packages
2. **Integration Patterns** - Demonstrates how packages work together  
3. **Adoption Acceleration** - Reduces barrier to entry for using packages
4. **Documentation by Example** - Generated projects serve as working examples

## 🤝 Contributing

The scaffolding generator is part of EXTRACT-008 (CLI Application Template). Contributions should:

1. Follow the established template patterns
2. Include appropriate implementation tokens
3. Update documentation and examples
4. Test generated projects build successfully

## [ACTION:format-processing] Implementation Tokens

Key implementation areas:

- ` [HIGH] EXTRACT-008: Project scaffolding system [ACTION:maintenance]` - Core generator functionality
- ` [HIGH] EXTRACT-008: Component selection interface [ACTION:maintenance]` - Package selection UI
- ` [HIGH] EXTRACT-008: Project template generation [ACTION:maintenance]` - File generation system
- ` [HIGH] EXTRACT-008: Template management system [ACTION:maintenance]` - Template coordination

---

**Created**: 2025-01-02  
**Part of**: EXTRACT-008 (CLI Application Template)  
**Status**: Subtask 2 Implementation  
**Next**: Integration documentation, migration guide, dependency mapping 