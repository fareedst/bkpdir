# CLI Template Application

A complete CLI application template showcasing all extracted packages from the bkpdir project. This template demonstrates how to build modern CLI applications using the extracted, reusable components.

## 🎯 Purpose

This template serves as:
- **Working Example**: Complete CLI application using extracted packages
- **Foundation**: Starting point for new CLI projects
- **Documentation**: Live examples of package integration
- **Best Practices**: Demonstrates proper usage patterns

## 📦 Extracted Packages Demonstrated

### 1. **pkg/config** - Configuration Management
- Schema-agnostic configuration loading
- Multiple configuration sources (files, environment, defaults)
- Configuration merging and validation
- Source tracking for configuration values

### 2. **pkg/errors** - Structured Error Handling
- Structured error types with context
- Error wrapping and unwrapping
- Status code integration
- Enhanced error reporting

### 3. **pkg/resources** - Resource Management
- Automatic resource cleanup
- Resource tracking and leak prevention
- Context-aware resource management
- Safe resource disposal patterns

### 4. **pkg/formatter** - Output Formatting
- Template-based formatting
- Printf-style formatting
- Structured output collection
- ANSI color support

### 5. **pkg/git** - Git Integration
- Repository detection and validation
- Branch and commit information extraction
- Dirty status detection
- Git configuration integration

### 6. **pkg/cli** - CLI Command Framework
- Command patterns and structure
- Flag handling and validation
- Context management
- Dry-run support

### 7. **pkg/fileops** - File Operations
- Safe file creation and manipulation
- File comparison and verification
- Path handling and validation
- Atomic file operations

### 8. **pkg/processing** - Concurrent Processing
- Worker pool patterns
- Context-aware processing
- Progress tracking and reporting
- Error handling and recovery

## 🚀 Quick Start

### Build and Run

```bash
# Build the application
cd cmd/cli-template
go build -o cli-template

# Run the application
./cli-template --help
```

### Available Commands

```bash
# Show help
./cli-template --help

# Show version information (demonstrates pkg/git)
./cli-template version

# Show configuration capabilities (demonstrates pkg/config)
./cli-template config

# Create example files (demonstrates pkg/fileops, pkg/errors, pkg/resources)
./cli-template create example.txt

# Process files concurrently (demonstrates pkg/processing)
./cli-template process file1.txt file2.txt file3.txt

# Use global flags
./cli-template --verbose config
./cli-template --dry-run create test.txt
```

## 🏗️ Architecture

### Command Structure

```
cli-template/
├── main.go              # Application entry point
├── cmd/                 # Command implementations
│   ├── root.go         # Root command and global setup
│   ├── version.go      # Version command (pkg/git demo)
│   ├── config.go       # Configuration command (pkg/config demo)
│   ├── create.go       # File creation command (pkg/fileops demo)
│   └── process.go      # Processing command (pkg/processing demo)
├── config/             # Configuration templates
├── examples/           # Usage examples
└── README.md           # This file
```

### Package Integration Patterns

1. **Foundation Layer**: pkg/config, pkg/errors (used by everything)
2. **Resource Layer**: pkg/resources (builds on errors)
3. **Service Layer**: pkg/formatter, pkg/git (domain-specific)
4. **Framework Layer**: pkg/cli (orchestration)
5. **Operation Layer**: pkg/fileops, pkg/processing (business logic)

## 📚 Usage Examples

### Basic Usage

```bash
# Show application help
./cli-template --help

# Show version with Git information
./cli-template version

# Display configuration management features
./cli-template config --verbose
```

### File Operations

```bash
# Create an example file
./cli-template create demo.txt

# Create with verbose output
./cli-template --verbose create detailed.txt

# Dry run (show what would be done)
./cli-template --dry-run create test.txt
```

### Concurrent Processing

```bash
# Process multiple files
./cli-template process file1.txt file2.txt file3.txt

# Process with verbose output
./cli-template --verbose process *.txt
```

### Configuration Management

```bash
# Show configuration capabilities
./cli-template config

# Use custom config file
./cli-template --config custom.yml config

# Show verbose configuration details
./cli-template --verbose config
```

## 🔧 Development

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd bkpdir/cmd/cli-template

# Build the application
go build -o cli-template

# Run tests (when available)
go test ./...
```

### Adding New Commands

1. Create a new file in `cmd/` directory
2. Implement the command using cobra patterns
3. Add the command to `root.go` in the `init()` function
4. Follow the existing patterns for package integration

### Package Integration

Each command demonstrates different aspects of the extracted packages:

- **Configuration**: Use pkg/config for settings management
- **Error Handling**: Use pkg/errors for structured error reporting
- **Resource Management**: Use pkg/resources for cleanup
- **Output**: Use pkg/formatter for consistent output
- **Git Integration**: Use pkg/git for repository information
- **File Operations**: Use pkg/fileops for safe file handling
- **Processing**: Use pkg/processing for concurrent operations

## 📖 Documentation

### Integration Guides

- [Package Integration Guide](../../docs/integration-guide.md) (Coming Soon)
- [Migration Guide](../../docs/migration-guide.md) (Coming Soon)
- [Package Reference](../../docs/package-reference.md) (Coming Soon)

### Tutorials

- [Building Your First CLI](../../docs/tutorials/first-cli.md) (Coming Soon)
- [Advanced Integration Patterns](../../docs/tutorials/advanced-patterns.md) (Coming Soon)
- [Performance Optimization](../../docs/tutorials/performance.md) (Coming Soon)

## 🎯 Next Steps

### Phase 1: Basic Template (Current)
- [x] Basic command structure
- [x] Example commands for each package
- [x] Documentation and README
- [ ] Configuration file examples
- [ ] Basic testing

### Phase 2: Full Integration (Next)
- [ ] Complete package integration
- [ ] Real configuration loading
- [ ] Actual Git information display
- [ ] Working file operations
- [ ] Functional concurrent processing

### Phase 3: Scaffolding System (Future)
- [ ] Interactive project generator
- [ ] Component selection system
- [ ] Project templates
- [ ] Automated dependency management

### Phase 4: Advanced Features (Future)
- [ ] Plugin system
- [ ] Advanced configuration schemas
- [ ] Performance monitoring
- [ ] Comprehensive testing suite

## 🤝 Contributing

This template is part of the bkpdir extraction project. Contributions should:

1. Follow the established package integration patterns
2. Maintain compatibility with extracted packages
3. Include appropriate documentation
4. Add tests for new functionality

## 📄 License

This template is part of the bkpdir project and follows the same licensing terms.

---

**Note**: This is a working template that demonstrates the extracted packages. Full integration with all package features is ongoing. See the [project roadmap](../../docs/extract-008-cli-template-plan.md) for detailed implementation plans. 