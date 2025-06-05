# Configuration Inspection Guide

<!-- 🔻 CFG-006: Documentation - 📝 Configuration inspection guide -->

## 📑 Overview

The BkpDir configuration system provides comprehensive configuration inspection capabilities through the enhanced `config` command. With automatic field discovery using Go reflection, you can inspect over 100 configuration parameters without any manual maintenance.

This guide covers the powerful configuration visibility features introduced in CFG-006, including automatic field discovery, inheritance chain visualization, multiple output formats, and advanced filtering options.

## 🔍 Key Features

### Automatic Field Discovery
- **Zero Maintenance**: New configuration fields appear automatically
- **Complete Coverage**: 100+ configuration parameters discovered via reflection
- **Type Awareness**: Handles structs, slices, maps, pointers, and interfaces
- **Nested Support**: Traverses embedded and nested configuration structures

### Source Attribution and Inheritance
- **Complete Traceability**: Shows exactly where each value originates
- **Inheritance Visualization**: Displays full inheritance chain resolution
- **Source Priority**: Environment → config files → inheritance → defaults
- **Conflict Detection**: Identifies where values were overridden

### Performance Optimization
- **Reflection Caching**: 60%+ performance improvement through caching
- **Lazy Evaluation**: Only resolves sources for displayed fields
- **Sub-second Response**: Optimized for responsive development workflow
- **Memory Efficient**: Minimal overhead for configuration inspection

## 🚀 Basic Usage

### Display All Configuration Values

```bash
# Show all configuration parameters (100+ fields auto-discovered)
bkpdir config
```

This displays all configuration fields in a clean table format with their current values and sources.

### Show Only Overridden Values

```bash
# Display only non-default values
bkpdir config --overrides-only
```

Perfect for quickly seeing what customizations have been made to the default configuration.

### View Source Attribution

```bash
# Show detailed source information for each value
bkpdir config --sources
```

Displays the complete source chain showing where each configuration value originated.

## 🌳 Output Formats

### Table Format (Default)

```bash
# Clean table display for quick scanning
bkpdir config --format table
```

**Example Output:**
```
Configuration Values:
╭─────────────────────────┬────────────────────────┬─────────────────╮
│ Parameter               │ Value                  │ Source          │
├─────────────────────────┼────────────────────────┼─────────────────┤
│ archive_dir_path        │ ~/Archives             │ Config File     │
│ include_git_info        │ true                   │ Environment     │
│ exclude_patterns        │ [*.tmp, *.log]         │ Inheritance     │
│ verification_timeout    │ 30s                    │ Default         │
╰─────────────────────────┴────────────────────────┴─────────────────╯
```

### Tree Format

```bash
# Hierarchical display showing configuration categories
bkpdir config --format tree
```

**Example Output:**
```
Configuration Tree:
├── Archive Settings
│   ├── archive_dir_path: ~/Archives (Config File)
│   ├── archive_name_format: %{timestamp}_%{branch} (Default)
│   └── include_git_info: true (Environment: BKPDIR_INCLUDE_GIT_INFO)
├── Backup Settings
│   ├── backup_dir_path: ~/Backups (Config File)
│   └── file_comparison_method: hash (Default)
└── Format Settings
    ├── timestamp_format: 2006-01-02_15-04-05 (Default)
    └── printf_format: [%{level}] %{message} (Config File)
```

### JSON Format

```bash
# Machine-readable output for automation
bkpdir config --format json
```

**Example Output:**
```json
{
  "configuration": [
    {
      "name": "archive_dir_path",
      "value": "~/Archives",
      "source": "Config File",
      "category": "Archive Settings",
      "type": "string",
      "is_overridden": true,
      "inheritance_chain": ["~/.bkpdir.yml"],
      "merge_strategy": "override"
    },
    {
      "name": "include_git_info",
      "value": "true",
      "source": "Environment",
      "category": "Archive Settings", 
      "type": "bool",
      "is_overridden": true,
      "environment_variable": "BKPDIR_INCLUDE_GIT_INFO"
    }
  ]
}
```

## 🔍 Advanced Filtering

### Filter by Pattern

```bash
# Show only configuration fields matching a pattern
bkpdir config --filter "archive"
```

Displays only fields containing "archive" in their name.

### Filter by Category

```bash
# Show only Git-related configuration
bkpdir config --filter "git"

# Show only format-related settings
bkpdir config --filter "format"
```

### Combine Filters

```bash
# Show only overridden archive settings with sources
bkpdir config --overrides-only --filter "archive" --sources
```

## 🔗 Inheritance Chain Visualization

### Understanding Inheritance Resolution

The configuration system supports layered inheritance (CFG-005) with sophisticated merge strategies:

```bash
# View complete inheritance chains
bkpdir config --sources --format tree
```

**Example Inheritance Resolution:**
```
exclude_patterns:
├── Environment: BKPDIR_EXCLUDE_PATTERNS → [*.secret] (Highest Priority)
├── Local Config: ./project/.bkpdir.yml → [node_modules/] (+ merge strategy)
├── Inherited: ~/.bkpdir.yml → [*.tmp, *.log] (Base values)
└── Default: [] (Fallback)
│
Final Value: [*.secret, *.tmp, *.log, node_modules/] (Merged result)
```

### Merge Strategy Indicators

- **Override** (default): Child value replaces parent value
- **Merge (+)**: Child values append to parent values  
- **Prepend (^)**: Child values prepend to parent values
- **Replace (!)**: Child value completely replaces, ignoring inheritance
- **Default (=)**: Use child value only if parent doesn't define it

## ⚡ Performance Features

### Reflection Caching

```bash
# First run: ~24.3μs (with reflection)
bkpdir config

# Subsequent runs: ~10.0μs (cached reflection results)
bkpdir config
```

### Lazy Source Evaluation

```bash
# Only resolves sources for displayed fields
bkpdir config --filter "archive" --sources
```

This optimization means filtering operations are extremely fast, even with 100+ configuration fields.

### Incremental Resolution

```bash
# Fast single-field access
bkpdir config archive_dir_path
```

Sub-100ms response time for single field access through incremental resolution.

## 🔧 Integration Features

### CFG-005 Inheritance Integration

The configuration inspection seamlessly integrates with the layered inheritance system:

```bash
# Debug inheritance resolution
bkpdir config --sources --format tree
```

Shows complete inheritance chains with merge strategy application.

### EXTRACT-001 Package Integration

Works with the extracted `pkg/config` architecture:

```bash
# All reflection and discovery works through pkg/config
bkpdir config --format json | jq '.configuration[].source'
```

### Environment Variable Support

```bash
# Override any field via environment
BKPDIR_INCLUDE_GIT_INFO=false bkpdir config include_git_info

# View environment variable names
bkpdir config --sources | grep Environment
```

## 📊 Practical Examples

### Development Workflow

```bash
# Quick configuration overview
bkpdir config --overrides-only

# Debug specific setting
bkpdir config --filter "timeout" --sources

# Check inheritance chain
bkpdir config exclude_patterns --sources --format tree
```

### Automation Integration

```bash
# Extract configuration for scripts
bkpdir config --format json | jq -r '.configuration[] | select(.is_overridden) | "\(.name)=\(.value)"'

# Check if specific setting is overridden
bkpdir config archive_dir_path --format json | jq '.configuration[0].is_overridden'
```

### Configuration Audit

```bash
# Find all environment variable overrides
bkpdir config --sources | grep "Environment"

# List all inherited values
bkpdir config --sources --format tree | grep "Inherited"

# Show configuration conflicts
bkpdir config --sources | grep "Overridden"
```

## 🎯 Best Practices

### Regular Configuration Review

```bash
# Weekly configuration review
bkpdir config --overrides-only --sources
```

Helps maintain awareness of your configuration customizations.

### Before Making Changes

```bash
# Document current state before configuration changes
bkpdir config --format json > config-backup-$(date +%Y%m%d).json
```

### Performance Monitoring

```bash
# Monitor configuration inspection performance
time bkpdir config --format table
```

Should complete in sub-100ms for optimal development experience.

### Troubleshooting Configuration Issues

```bash
# Complete configuration diagnosis
bkpdir config --sources --format tree > config-debug.txt
```

Provides complete configuration state for troubleshooting.

## 🚀 Advanced Use Cases

### Configuration Validation

```bash
# Verify all required fields are set
bkpdir config --format json | jq '.configuration[] | select(.value == "") | .name'
```

### Environment Comparison

```bash
# Development vs Production configuration diff
bkpdir config --format json > dev-config.json
BKPDIR_ENV=production bkpdir config --format json > prod-config.json
diff dev-config.json prod-config.json
```

### Configuration Documentation

```bash
# Generate configuration documentation
bkpdir config --format json | jq -r '.configuration[] | "- **\(.name)**: \(.value) (\(.source))"'
```

## 📝 Summary

The CFG-006 configuration inspection system provides:

- **🔍 Zero-maintenance discovery**: 100+ fields automatically discovered
- **📊 Complete visibility**: Full inheritance chain and source attribution
- **⚡ High performance**: Sub-second response with reflection caching
- **🎨 Flexible display**: Table, tree, and JSON output formats
- **🔧 Powerful filtering**: Pattern-based and category-based filtering
- **🔗 Seamless integration**: Works with CFG-005 inheritance and EXTRACT-001 packages

This powerful configuration inspection capability makes it easy to understand, debug, and maintain complex configuration setups in your BkpDir workflow.

---

*For troubleshooting specific configuration issues, see the [Configuration Troubleshooting Guide](configuration-troubleshooting.md).*
*For practical examples and use cases, see the [Configuration Examples](configuration-examples.md).* 