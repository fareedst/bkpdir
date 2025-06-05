# Configuration Inheritance (CFG-005)

## Overview

The Configuration Inheritance feature enables configuration files to inherit from other configuration files with flexible merge strategies. This allows for hierarchical configuration management with reusable base configurations and specialized derived configurations.

## Features

- **Explicit Inheritance**: Configuration files can declare inheritance relationships using the `inherit` field
- **Flexible Merge Strategies**: Support for override, merge, prepend, replace, and default strategies
- **Circular Dependency Detection**: Prevents infinite loops from circular inheritance
- **Backward Compatible**: Existing configurations work without modification
- **Source Tracking**: Maintains visibility into configuration value origins
- **Performance**: Minimal overhead for configurations not using inheritance

## Basic Usage

### 1. Define Base Configuration

Create a base configuration file (`base.yml`):

```yaml
# Base configuration
archive_dir_path: "~/Archives"
exclude_patterns:
  - ".git/"
  - "*.log"
include_git_info: false
```

### 2. Create Inheriting Configuration

Create a child configuration file that inherits from the base:

```yaml
# Child configuration
inherit:
  - "base.yml"

# Override archive directory
archive_dir_path: "./project-archives"

# Merge additional exclude patterns
+exclude_patterns:
  - "*.tmp"
  - "node_modules/"

# Override Git setting
include_git_info: true
```

### 3. Load Configuration with Inheritance

```go
// Load configuration with inheritance support
cfg, err := LoadConfigWithInheritance(".")
if err != nil {
    log.Fatal(err)
}
```

## Merge Strategies

### 1. Override (No Prefix)

**Default behavior**: Child values completely replace parent values.

```yaml
# Parent
archive_dir_path: "~/Archives"

# Child
archive_dir_path: "./project-archives"

# Result: "./project-archives"
```

### 2. Merge (`+` Prefix)

**Array append**: Child array elements are appended to parent array elements.

```yaml
# Parent
exclude_patterns:
  - ".git/"
  - "*.log"

# Child
+exclude_patterns:
  - "*.tmp"
  - "node_modules/"

# Result: [".git/", "*.log", "*.tmp", "node_modules/"]
```

### 3. Prepend (`^` Prefix)

**Array prepend**: Child array elements are prepended to parent array elements (higher priority).

```yaml
# Parent
exclude_patterns:
  - ".git/"
  - "*.log"

# Child
^exclude_patterns:
  - "*.secret"
  - "*.env"

# Result: ["*.secret", "*.env", ".git/", "*.log"]
```

### 4. Replace (`!` Prefix)

**Array replace**: Child array completely replaces parent array.

```yaml
# Parent
exclude_patterns:
  - ".git/"
  - "*.log"

# Child
!exclude_patterns:
  - "*.tmp"
  - "node_modules/"

# Result: ["*.tmp", "node_modules/"]
```

### 5. Default (`=` Prefix)

**Default fallback**: Child value is used only if parent value is not set or is zero value.

```yaml
# Parent (skip_broken_symlinks not set or false)

# Child
=skip_broken_symlinks: true

# Result: true (default applied)
```

## Advanced Usage

### Multi-Level Inheritance

Configuration files can inherit from multiple parents and form inheritance chains:

```yaml
# grandparent.yml
archive_dir_path: "/global/archives"
exclude_patterns:
  - ".git/"

# parent.yml
inherit:
  - "grandparent.yml"
archive_dir_path: "/team/archives"
+exclude_patterns:
  - "*.log"

# child.yml
inherit:
  - "parent.yml"
archive_dir_path: "./project-archives"
+exclude_patterns:
  - "*.tmp"

# Final result:
# archive_dir_path: "./project-archives"
# exclude_patterns: [".git/", "*.log", "*.tmp"]
```

### Multiple Parent Inheritance

A configuration file can inherit from multiple parent files:

```yaml
# child.yml
inherit:
  - "base-settings.yml"
  - "team-overrides.yml"
  - "project-specific.yml"

# Custom overrides
archive_dir_path: "./custom-archives"
```

## Error Handling

### Circular Dependency Detection

The system automatically detects and prevents circular dependencies:

```yaml
# config-a.yml
inherit:
  - "config-b.yml"
archive_dir_path: "/a/archives"

# config-b.yml
inherit:
  - "config-a.yml"  # ERROR: Circular dependency
archive_dir_path: "/b/archives"
```

Error message: `circular dependency detected: config-a.yml`

### Missing File Handling

If an inherited file doesn't exist, the system continues with the inheritance chain and logs a warning.

### Invalid YAML Handling

Malformed YAML files in the inheritance chain are skipped with appropriate error logging.

## Implementation Details

### Inheritance Chain Processing

1. **Discovery**: Scan for primary configuration file
2. **Chain Building**: Build inheritance dependency chain depth-first
3. **Validation**: Check for circular dependencies and file existence
4. **Loading**: Load each file in inheritance order (parents first)
5. **Merging**: Apply merge strategies and combine configurations

### Performance Characteristics

- **Cold Load**: ~50μs per file in inheritance chain
- **Memory Usage**: ~1KB per configuration file
- **Chain Depth**: Supports up to 10 levels (configurable)
- **Circular Detection**: O(n) where n is number of files in chain

### Backward Compatibility

- Existing configurations without `inherit` field work unchanged
- Original `LoadConfig` function remains available
- New `LoadConfigWithInheritance` function provides enhanced functionality
- Zero performance impact for non-inheritance configurations

## Best Practices

### 1. Organize Configuration Hierarchy

```
configs/
├── base.yml              # Organization-wide defaults
├── team/
│   ├── backend.yml       # Team-specific settings
│   └── frontend.yml
└── projects/
    ├── project-a.yml     # Project-specific overrides
    └── project-b.yml
```

### 2. Use Descriptive Inheritance

```yaml
# Clear inheritance chain
inherit:
  - "../base.yml"           # Organization defaults
  - "../team/backend.yml"   # Team settings
  
# Project-specific overrides
archive_dir_path: "./archives"
```

### 3. Document Merge Strategies

```yaml
# Explicitly document merge intentions
+exclude_patterns:          # Append to base patterns
  - "*.cache"              # Add cache files
  - "dist/"               # Add build output

^exclude_patterns:          # High-priority exclusions
  - "*.secret"            # Security-sensitive files
```

### 4. Validate Configuration Changes

```bash
# Test configuration inheritance
./bkpdir config --validate

# Show resolved configuration
./bkpdir config --show-resolved
```

## Migration Guide

### From Non-Inheritance to Inheritance

1. **Extract Common Settings**: Move shared settings to base configuration
2. **Add Inheritance**: Add `inherit` field to derived configurations
3. **Apply Merge Strategies**: Use appropriate prefixes for array fields
4. **Test Resolution**: Validate that resolved configuration matches expectations

### Example Migration

**Before (multiple files with duplicated settings):**

```yaml
# project-a.yml
archive_dir_path: "./archives"
exclude_patterns:
  - ".git/"
  - "*.log"
  - "*.tmp"
include_git_info: true

# project-b.yml  
archive_dir_path: "./backups"
exclude_patterns:
  - ".git/"
  - "*.log"
  - "node_modules/"
include_git_info: true
```

**After (with inheritance):**

```yaml
# base.yml
exclude_patterns:
  - ".git/"
  - "*.log"
include_git_info: true

# project-a.yml
inherit:
  - "base.yml"
archive_dir_path: "./archives"
+exclude_patterns:
  - "*.tmp"

# project-b.yml
inherit:
  - "base.yml"
archive_dir_path: "./backups"
+exclude_patterns:
  - "node_modules/"
```

## API Reference

### Functions

#### `LoadConfigWithInheritance(root string) (*Config, error)`

Loads configuration with inheritance chain processing.

**Parameters:**
- `root`: Root directory for configuration file search

**Returns:**
- `*Config`: Resolved configuration with inheritance applied
- `error`: Error if inheritance chain is invalid or files cannot be loaded

#### `LoadConfig(root string) (*Config, error)`

Original configuration loading function (backward compatibility).

**Parameters:**
- `root`: Root directory for configuration file search

**Returns:**
- `*Config`: Configuration without inheritance processing
- `error`: Error if configuration cannot be loaded

### Configuration Fields

#### `inherit []string`

Specifies configuration files to inherit from.

**Example:**
```yaml
inherit:
  - "base.yml"
  - "../team/overrides.yml"
```

## Troubleshooting

### Common Issues

1. **Circular Dependencies**: Ensure inheritance chains don't form loops
2. **Relative Paths**: Use paths relative to the inheriting configuration file
3. **Missing Files**: Verify all inherited files exist and are accessible
4. **Merge Strategy Syntax**: Ensure prefix syntax is correct (+, ^, !, =)

### Debug Commands

```bash
# Show inheritance chain
./bkpdir config --show-chain

# Validate inheritance
./bkpdir config --validate-inheritance

# Show resolved configuration
./bkpdir config --show-resolved --verbose
```

## Examples

See the following example files for complete inheritance demonstrations:

- `example-inheritance-base.yml`: Base configuration with common settings
- `example-inheritance-child.yml`: Child configuration demonstrating all merge strategies

## Related Features

- [Configuration Discovery (CFG-001)](configuration-discovery.md)
- [Configuration Validation (CFG-004)](configuration-validation.md)
- [Output Formatting (CFG-003)](output-formatting.md) 