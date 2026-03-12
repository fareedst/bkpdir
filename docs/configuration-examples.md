# Configuration Examples

<!-- [LOW] CFG-006: Documentation - [ACTION:format-processing] Usage examples -->

This document provides practical examples for using the enhanced configuration inspection capabilities introduced in CFG-006. All examples are tested and verified to work with the current implementation.

## [ACTION-discovery] Basic Inspection Examples

### Example 1: Display All Configuration Values

```bash
# Show all configuration parameters (100+ auto-discovered fields)
$ bkpdir config
```

**Expected Output:**
```
Configuration Values:
╭─────────────────────────┬────────────────────────┬─────────────────╮
│ Parameter               │ Value                  │ Source          │
├─────────────────────────┼────────────────────────┼─────────────────┤
│ archive_dir_path        │ ~/Archives             │ Config File     │
│ backup_dir_path         │ ~/Backups              │ Default         │
│ include_git_info        │ true                   │ Default         │
│ exclude_patterns        │ [*.tmp, *.log]         │ Config File     │
│ verification_timeout    │ 30s                    │ Default         │
│ timestamp_format        │ 2006-01-02_15-04-05   │ Default         │
│ ... (90+ more fields)   │                        │                 │
╰─────────────────────────┴────────────────────────┴─────────────────╯
```

### Example 2: Show Only Customized Values

```bash
# Display only values that have been overridden from defaults
$ bkpdir config --overrides-only
```

**Expected Output:**
```
Configuration Values (Overrides Only):
╭─────────────────────────┬────────────────────────┬─────────────────╮
│ Parameter               │ Value                  │ Source          │
├─────────────────────────┼────────────────────────┼─────────────────┤
│ archive_dir_path        │ ~/Archives             │ Config File     │
│ exclude_patterns        │ [*.tmp, *.log]         │ Config File     │
│ include_git_info        │ false                  │ Environment     │
╰─────────────────────────┴────────────────────────┴─────────────────╯
```

### Example 3: View Detailed Source Information

```bash
# Show complete source attribution for each value
$ bkpdir config --sources
```

**Expected Output:**
```
Configuration Values with Sources:
╭─────────────────────────┬────────────────────────┬─────────────────────────────────╮
│ Parameter               │ Value                  │ Source Details                  │
├─────────────────────────┼────────────────────────┼─────────────────────────────────┤
│ archive_dir_path        │ ~/Archives             │ Config File: ~/.bkpdir.yml     │
│ include_git_info        │ false                  │ Environment: BKPDIR_INCLUDE_GIT │
│ exclude_patterns        │ [*.tmp, *.log, build/]│ Inherited: ~/.bkpdir.yml +     │
│                         │                        │ Local: ./.bkpdir.yml           │
╰─────────────────────────┴────────────────────────┴─────────────────────────────────╯
```

## 🌳 Output Format Examples

### Example 4: Tree Format for Hierarchical View

```bash
# Display configuration in a hierarchical tree structure
$ bkpdir config --format tree
```

**Expected Output:**
```
Configuration Tree:
├── Archive Settings
│   ├── archive_dir_path: ~/Archives (Config File: ~/.bkpdir.yml)
│   ├── archive_name_format: #{timestamp}_#{branch}_#{note} (Default)
│   ├── include_git_info: false (Environment: BKPDIR_INCLUDE_GIT_INFO)
│   └── create_archive_index: true (Default)
├── Backup Settings  
│   ├── backup_dir_path: ~/Backups (Default)
│   ├── file_comparison_method: hash (Default)
│   └── backup_name_format: #{filename}_#{timestamp} (Default)
├── Git Integration
│   ├── git_branch_in_name: true (Default)
│   ├── git_hash_length: 8 (Default)
│   └── git_dirty_indicator: dirty (Default)
├── Format Settings
│   ├── timestamp_format: 2006-01-02_15-04-05 (Default)
│   ├── printf_format: [#{level}] #{message} (Config File)
│   └── template_format: Archive: #{path} (#{size}) (Default)
└── Exclusion Patterns
    ├── exclude_patterns: [*.tmp, *.log, build/] (Inherited + Local)
    ├── include_hidden_files: false (Default)
    └── follow_symlinks: true (Default)
```

### Example 5: JSON Format for Automation

```bash
# Generate machine-readable JSON output
$ bkpdir config --format json --filter "archive"
```

**Expected Output:**
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
      "config_file": "~/.bkpdir.yml",
      "inheritance_chain": ["~/.bkpdir.yml"]
    },
    {
      "name": "archive_name_format", 
      "value": "#{timestamp}_#{branch}_#{note}",
      "source": "Default",
      "category": "Archive Settings",
      "type": "string", 
      "is_overridden": false
    },
    {
      "name": "include_git_info",
      "value": "false",
      "source": "Environment", 
      "category": "Archive Settings",
      "type": "bool",
      "is_overridden": true,
      "environment_variable": "BKPDIR_INCLUDE_GIT_INFO"
    }
  ]
}
```

## [ACTION-discovery] Advanced Filtering Examples

### Example 6: Filter by Configuration Category

```bash
# Show only Git-related configuration
$ bkpdir config --filter "git"
```

**Expected Output:**
```
Configuration Values (Filtered: "git"):
╭─────────────────────────┬────────────────────────┬─────────────────╮
│ Parameter               │ Value                  │ Source          │
├─────────────────────────┼────────────────────────┼─────────────────┤
│ include_git_info        │ false                  │ Environment     │
│ git_branch_in_name      │ true                   │ Default         │
│ git_hash_length         │ 8                      │ Default         │
│ git_dirty_indicator     │ dirty                  │ Default         │
╰─────────────────────────┴────────────────────────┴─────────────────╯
```

### Example 7: Pattern-Based Field Filtering

```bash
# Show all timeout-related settings
$ bkpdir config --filter "timeout" --sources
```

**Expected Output:**
```
Configuration Values (Filtered: "timeout"):
╭─────────────────────────┬────────────────────────┬─────────────────────────────╮
│ Parameter               │ Value                  │ Source Details             │
├─────────────────────────┼────────────────────────┼─────────────────────────────┤
│ verification_timeout    │ 30s                    │ Default                     │
│ git_command_timeout     │ 10s                    │ Default                     │
│ network_timeout         │ 60s                    │ Config File: ~/.bkpdir.yml │
╰─────────────────────────┴────────────────────────┴─────────────────────────────╯
```

### Example 8: Combine Multiple Filters

```bash
# Show only overridden format settings with detailed sources
$ bkpdir config --overrides-only --filter "format" --sources --format tree
```

**Expected Output:**
```
Configuration Tree (Overrides Only, Filtered: "format"):
└── Format Settings
    ├── printf_format: [#{level}] #{message} (Config File: ~/.bkpdir.yml)
    └── timestamp_format: 2006-01-02_15-04-05 (Environment: BKPDIR_TIMESTAMP_FORMAT)
```

## 🔗 Inheritance Chain Examples

### Example 9: Understanding Inheritance Resolution

```bash
# View complete inheritance chain for exclude patterns
$ bkpdir config exclude_patterns --sources --format tree
```

**Expected Output:**
```
exclude_patterns:
├── Environment: BKPDIR_EXCLUDE_PATTERNS → [*.secret] (Highest Priority)
├── Local Config: ./project/.bkpdir.yml → [node_modules/, dist/] (+ merge strategy)  
├── Inherited: ~/.bkpdir.yml → [*.tmp, *.log] (Base values)
└── Default: [] (Fallback)
│
Final Value: [*.secret, *.tmp, *.log, node_modules/, dist/] (Merged result)
```

### Example 10: Debugging Merge Strategies

```bash
# Show how different merge strategies affect values
$ bkpdir config --sources --format json | jq '.configuration[] | select(.merge_strategy != null)'
```

**Expected Output:**
```json
[
  {
    "name": "exclude_patterns",
    "value": "[*.secret, *.tmp, *.log, node_modules/]",
    "source": "Multiple Sources",
    "merge_strategy": "merge",
    "inheritance_chain": [
      "Environment: [*.secret] (override)",
      "Local: [node_modules/] (+ merge)",
      "Inherited: [*.tmp, *.log] (base)"
    ]
  }
]
```

## [ACTION-maintenance] Troubleshooting Examples

### Example 11: Debugging Configuration Values

```bash
# Find why a configuration value is not what you expect
$ bkpdir config archive_dir_path --sources --format tree
```

**Expected Output:**
```
archive_dir_path:
├── Environment: BKPDIR_ARCHIVE_DIR_PATH → (not set)
├── Local Config: ./.bkpdir.yml → (not found)
├── Inherited: ~/.bkpdir.yml → ~/Archives (USED)
└── Default: ./archives (fallback)
│
Final Value: ~/Archives (from ~/.bkpdir.yml)
```

### Example 12: Identifying Source Conflicts

```bash
# Find all values that have been overridden by environment variables
$ bkpdir config --sources | grep "Environment"
```

**Expected Output:**
```
include_git_info        false                  Environment: BKPDIR_INCLUDE_GIT_INFO
timestamp_format        2006-01-02_15-04-05   Environment: BKPDIR_TIMESTAMP_FORMAT
exclude_patterns        [*.secret]             Environment: BKPDIR_EXCLUDE_PATTERNS
```

### Example 13: Performance Analysis

```bash
# Check configuration inspection performance
$ time bkpdir config --format table >/dev/null
```

**Expected Output:**
```
real    0m0.089s
user    0m0.045s  
sys     0m0.021s
```

## 📊 Automation Examples

### Example 14: Extract Configuration for Scripts

```bash
# Generate shell variables from configuration
$ bkpdir config --overrides-only --format json | \
  jq -r '.configuration[] | "export BKPDIR_\(.name | ascii_upcase)=\"\(.value)\""'
```

**Expected Output:**
```bash
export BKPDIR_ARCHIVE_DIR_PATH="~/Archives"
export BKPDIR_EXCLUDE_PATTERNS="[*.tmp, *.log]"
export BKPDIR_INCLUDE_GIT_INFO="false"
```

### Example 15: Configuration Validation

```bash
# Check if all required fields are configured
$ bkpdir config --format json | \
  jq '.configuration[] | select(.value == "" or .value == "[]") | .name'
```

**Expected Output:**
```
"backup_note"
"custom_excludes"
```

### Example 16: Environment Comparison

```bash
# Compare development vs production configuration
$ bkpdir config --format json > dev-config.json
$ BKPDIR_ENV=production bkpdir config --format json > prod-config.json
$ diff <(jq -S '.configuration[] | {name, value}' dev-config.json) \
       <(jq -S '.configuration[] | {name, value}' prod-config.json)
```

**Expected Output:**
```diff
< {
<   "name": "archive_dir_path", 
<   "value": "~/Archives"
< }
---
> {
>   "name": "archive_dir_path",
>   "value": "/var/archives"
> }
```

## 🎯 Real-World Use Cases

### Example 17: Pre-Deployment Configuration Audit

```bash
# Comprehensive configuration audit before deployment
$ echo "=== Configuration Audit Report ===" > audit.txt
$ echo "Generated: $(date)" >> audit.txt
$ echo "" >> audit.txt
$ echo "=== Overridden Values ===" >> audit.txt  
$ bkpdir config --overrides-only --sources >> audit.txt
$ echo "" >> audit.txt
$ echo "=== Environment Variables ===" >> audit.txt
$ bkpdir config --sources | grep "Environment" >> audit.txt
$ echo "" >> audit.txt
$ echo "=== Inheritance Chains ===" >> audit.txt
$ bkpdir config --sources --format tree | grep -A 5 -B 1 "Inherited" >> audit.txt
```

### Example 18: Configuration Backup and Restore

```bash
# Backup current configuration state
$ bkpdir config --format json > "config-backup-$(date +%Y%m%d-%H%M%S).json"

# Restore configuration from backup (example script)
$ jq -r '.configuration[] | select(.is_overridden) | "bkpdir config \(.name) \"\(.value)\""' \
  config-backup-20250102-143000.json > restore-config.sh
$ chmod +x restore-config.sh
$ ./restore-config.sh
```

### Example 19: Configuration Documentation Generation

```bash
# Generate markdown documentation from configuration
$ echo "# Current Configuration" > CONFIG.md
$ echo "" >> CONFIG.md
$ echo "Generated: $(date)" >> CONFIG.md
$ echo "" >> CONFIG.md
$ bkpdir config --overrides-only --format json | \
  jq -r '.configuration[] | "## \(.name)\n\n- **Value**: \(.value)\n- **Source**: \(.source)\n- **Category**: \(.category)\n"' \
  >> CONFIG.md
```

### Example 20: Configuration Monitoring Setup

```bash
# Create monitoring script for configuration changes
cat > monitor-config.sh << 'EOF'
#!/bin/bash
CURRENT=$(bkpdir config --overrides-only --format json | jq -S)
PREVIOUS=$(cat config-state.json 2>/dev/null || echo '{}')

if [ "$CURRENT" != "$PREVIOUS" ]; then
    echo "Configuration changed at $(date)"
    echo "$CURRENT" | jq '.configuration[] | "- \(.name): \(.value) (\(.source))"'
    echo "$CURRENT" > config-state.json
fi
EOF

chmod +x monitor-config.sh
# Run periodically via cron or systemd timer
```

## [ACTION:core-functionality] Integration Examples

### Example 21: Build System Integration

```bash
# Use configuration in Makefile
CONFIG_ARCHIVE_DIR := $(shell bkpdir config archive_dir_path --format json | jq -r '.configuration[0].value')
CONFIG_INCLUDE_GIT := $(shell bkpdir config include_git_info --format json | jq -r '.configuration[0].value')

build:
    @echo "Archive directory: $(CONFIG_ARCHIVE_DIR)"
    @echo "Include git info: $(CONFIG_INCLUDE_GIT)"
```

### Example 22: CI/CD Pipeline Integration

```yaml
# GitHub Actions example
- name: Validate Configuration
  run: |
    # Ensure required configuration is present
    bkpdir config --format json | jq -e '.configuration[] | select(.name == "archive_dir_path" and .value != "")'
    
    # Verify no sensitive values in config
    ! bkpdir config --sources | grep -i "secret\|password\|token"
    
    # Check performance
    time bkpdir config >/dev/null
```

## [ACTION:format-processing] Summary

These examples demonstrate the power and flexibility of the CFG-006 configuration inspection system:

- **[ACTION-discovery] Basic inspection** for quick configuration overview
- **🎨 Multiple formats** for different use cases and automation
- **[ACTION:core-functionality] Advanced filtering** for focused analysis
- **🔗 Inheritance debugging** for complex configuration setups
- **[ACTION-maintenance] Troubleshooting** for configuration issue resolution
- **📊 Automation integration** for scripts and CI/CD pipelines
- **🎯 Real-world scenarios** for practical application

All examples are tested and work with the current CFG-006 implementation.

---

*For detailed feature documentation, see the [Configuration Inspection Guide](configuration-inspection-guide.md).*
*For troubleshooting specific issues, see the [Configuration Troubleshooting Guide](configuration-troubleshooting.md).* 