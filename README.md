# BkpDir - AI-First Directory Archiving CLI

> **[CRITICAL] ARCH-001: Archive naming convention implementation [ACTION:core-functionality]**

A powerful command-line tool for creating, managing, and verifying directory archives with AI-optimized development workflow.

## Quick Start

### Installation

```bash
# Build and install locally
make build-local
make install

# Or run directly
go run main.go --help
```

### Basic Usage

```bash
# Create full archive
bkpdir create /path/to/directory

# Create incremental archive
bkpdir create /path/to/directory --incremental

# Verify archive integrity
bkpdir verify /path/to/archive.tar.gz

# List all archives
bkpdir list
```

## AI-First Development

### Semantic Token System

This project uses a **semantic token system** for AI-optimized development:

```go
// [CRITICAL] ARCH-001: Archive naming convention [ACTION:core-functionality]
// [HIGH] CFG-003: Template formatting logic [ACTION:format-processing]
// [MEDIUM] GIT-004: Git submodule support [ACTION:discovery]
```

**Key Benefits:**
- **Perfect searchability**: `grep "\[CRITICAL\]" *.go` works everywhere
- **AI-native**: Semantic parsing vs visual interpretation
- **Cross-platform**: No Unicode or font dependencies
- **Maintainable**: Single source of truth in `project-tokens.yaml`

### Development Requirements

All development work must comply with:

1. **Semantic tokens** - Use `[CRITICAL|HIGH|MEDIUM|LOW] FEATURE-ID: Description [ACTION:context]`
2. **Validation** - Pass `make validate-token-enforcement`
3. **Registry compliance** - Use tokens from `project-tokens.yaml`
4. **Documentation consistency** - Keep tokens aligned across all files
5. **Milestone tracking** - Document completion of major milestones during development

See [Semantic Token System Requirements](docs/semantic-token-system-requirements.md) for complete details.

## Project Structure

```
bkpdir/
├── main.go                    # CLI entry point
├── archive.go                 # Archive operations
├── backup.go                  # Backup operations
├── config.go                  # Configuration management
├── pkg/                       # Extracted packages
│   ├── cli/                   # CLI framework
│   ├── config/                # Configuration loading
│   ├── errors/                # Error handling
│   ├── fileops/               # File operations
│   ├── formatter/             # Output formatting
│   ├── git/                   # Git integration
│   ├── processing/            # Processing pipelines
│   ├── resources/             # Resource management
│   └── testutil/              # Testing utilities
├── scripts/                   # Build and validation scripts
├── docs/                      # Documentation
└── project-tokens.yaml        # Semantic token registry
```

## Features

### Core Functionality

- **Archive Creation**: Full and incremental directory archiving
- **Backup Management**: File-level backup with deduplication
- **Verification**: Integrity checking with checksums
- **Git Integration**: Git-aware archiving with branch detection
- **Configuration**: Flexible YAML-based configuration
- **Template System**: Customizable output formatting

### AI-First Architecture

- **Semantic Tokens**: Machine-readable code annotations
- **Validation Integration**: Automated consistency checking
- **Package Extraction**: Modular, reusable components
- **Test Coverage**: Comprehensive testing framework
- **Documentation**: AI-optimized documentation system

## Configuration

### Basic Configuration

```yaml
# .bkpdir.yml
archive:
  base_directory: "/path/to/archives"
  format: "tar.gz"
  incremental: true

backup:
  directory: "/path/to/backups"
  max_versions: 5

git:
  include_branch: true
  include_status: true
```

### Advanced Configuration

```yaml
# Inheritance support
inherit_from: "../base-config.yml"

# Template customization
template:
  archive_name: "{{.Name}}-{{.Date}}-{{.Branch}}"
  output_format: "detailed"

# Exclusion patterns
exclude:
  - "*.tmp"
  - "node_modules/"
  - ".git/"
```

## Development

### Build Targets

```bash
# Development
make dev                    # Build and test
make build-local           # Build for current platform
make install              # Install to ~/.local/bin

# Testing
make test                 # Run all tests
make test-coverage        # Run with coverage
make test-performance     # Run performance tests

# Quality
make check                # All quality checks
make lint                 # Lint code
make validate-token-enforcement  # Validate semantic tokens

# Production
make build-all           # Build for all platforms
```

### Validation

```bash
# Semantic token validation
make validate-token-enforcement

# Strict validation (CI/CD)
make validate-tokens-strict

# Development validation
make validate-tokens
```

## Token Registry

The project uses a central registry for semantic tokens:

**File:** `project-tokens.yaml`

### Valid Priorities
- `CRITICAL` - Blocking operations, core system integrity
- `HIGH` - Important features, significant business logic
- `MEDIUM` - Secondary features, conditional processing
- `LOW` - Maintenance tasks, cleanup, documentation

### Valid Actions
- `core-functionality` - Essential system operations
- `format-processing` - Text formatting and output generation
- `discovery` - File system and configuration discovery
- `maintenance` - Code cleanup and refactoring
- `validation` - Input validation and error checking

### Feature IDs
- `ARCH-001` - Archive naming convention implementation
- `CFG-003` - Template formatting logic
- `GIT-004` - Git submodule support

## Migration from Legacy Icons

The project is migrating from Unicode icons to semantic tokens:

```bash
# Legacy (deprecated)
// [CRITICAL] ARCH-001: Archive naming convention [ACTION:core-functionality]

# Semantic (required)
// [CRITICAL] ARCH-001: Archive naming convention [ACTION:core-functionality]
```

**Status:** ~2715 legacy icons identified, ~30% migrated. See [Unicode Migration Status](docs/unicode-migration-current-status.md) for current progress.

## Contributing

### Requirements

1. **Semantic tokens** - All code must use semantic token format
2. **Validation** - Must pass `make validate-token-enforcement`
3. **Testing** - Include comprehensive tests
4. **Documentation** - Update docs with semantic tokens

### Process

1. Fork the repository
2. Create feature branch
3. Implement with semantic tokens
4. Run validation: `make check`
5. Submit pull request

### AI Assistant Guidelines

- **Token system first** - Always use semantic tokens
- **Registry compliance** - Check `project-tokens.yaml` for valid tokens
- **Validation integration** - Run validation after changes
- **Consistency maintenance** - Keep tokens aligned across all files

## License

MIT License - see LICENSE file for details.

## Links

- [Documentation](docs/)
- [Semantic Token Requirements](docs/semantic-token-system-requirements.md)
- [AI-First Development Guide](docs/context/AI-First-Development-Procedure-Complete-Guide.md)
- [Package Reference](docs/package-reference.md)
- [Migration Guide](docs/migration-guide.md)
- [Unicode Migration Status](docs/unicode-migration-current-status.md)
- [Migration Next Steps](docs/unicode-migration-next-steps.md)
- [Traceability-Enhanced Migration Plan](docs/unicode-migration-traceability-plan.md) 