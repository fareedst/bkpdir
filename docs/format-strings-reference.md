# Format Strings Reference

This document provides a comprehensive reference for all customizable format strings in bkpdir.

## Overview

All output format strings can be customized through your `.bkpdir.yml` configuration file. This allows you to:

- Customize message wording and formatting
- Add emoji or special characters
- Internationalize messages
- Integrate with other tools that expect specific output formats

## Format String Types

### Printf-Style Format Strings

Printf-style format strings use standard Go format verbs like `%s`, `%d`, `%v`:

```yaml
format_created_archive: "Created archive: %s\n"
```

**Common Format Verbs:**
- `%s` - String
- `%d` - Integer
- `%v` - Generic value (any type)
- `%f` - Floating point

### Template-Style Format Strings

Template-style format strings use named placeholders like `#{path}`, `#{name}`:

```yaml
template_created_archive: "Created archive: #{path}\n"
```

**Benefits:**
- Self-documenting (placeholder names show what data is inserted)
- More readable for complex formats
- Easier to reorder placeholders

### Special Placeholders

Some format strings support special placeholders:

- `#{size_human}` - Human-readable file size (e.g., "1.2MB", "455KB")
- `#{mtime}` - Modification time
- `#{hash}` - Git commit hash
- `#{branch}` - Git branch name

---

## Directory Operation Format Strings

### format_created_archive

**YAML Field:** `format_created_archive`

**Description:** Message displayed when an archive is successfully created.

**Placeholders:**
- `%s` - Archive path

**Default:**
```yaml
format_created_archive: "Created archive: %s\n"
```

**Examples:**
```yaml
# With emoji
format_created_archive: "✅ Archive created: %s\n"

# More verbose
format_created_archive: "Successfully created archive at: %s\n"

# Minimal
format_created_archive: "%s\n"
```

---

### format_identical_archive

**YAML Field:** `format_identical_archive`

**Description:** Message displayed when directory is identical to an existing archive (no changes detected).

**Placeholders:**
- `%s` - Archive path

**Default:**
```yaml
format_identical_archive: "Directory is identical to existing archive: %s\n"
```

**Examples:**
```yaml
# With emoji
format_identical_archive: "⚠️  No changes detected. Identical to: %s\n"

# Shorter
format_identical_archive: "Identical to: %s\n"

# More descriptive
format_identical_archive: "No changes since last archive: %s\n"
```

---

### format_list_archive

**YAML Field:** `format_list_archive`

**Description:** Format for listing archives.

**Placeholders:**
- `%s` - Archive path
- `#{size_human}` - Human-readable archive size

**Default:**
```yaml
format_list_archive: "%s (size: #{size_human})\n"
```

**Examples:**
```yaml
# With emoji
format_list_archive: "📦 %s (#{size_human})\n"

# More detailed
format_list_archive: "Archive: %s | Size: #{size_human}\n"

# Minimal
format_list_archive: "%s\n"
```

---

### format_config_value

**YAML Field:** `format_config_value`

**Description:** Format for displaying configuration values.

**Placeholders:**
- First `%s` - Configuration name
- Second `%s` - Configuration value
- Third `%s` - Source (where the value came from)

**Default:**
```yaml
format_config_value: "%s: %s (source: %s)\n"
```

**Examples:**
```yaml
# With emoji
format_config_value: "⚙️  %s: %s (from: %s)\n"

# More structured
format_config_value: "[%s] %s (source: %s)\n"

# Minimal
format_config_value: "%s=%s\n"
```

---

### format_dry_run_archive

**YAML Field:** `format_dry_run_archive`

**Description:** Message displayed during dry-run mode (when `--dry-run` flag is used).

**Placeholders:**
- `%s` - Archive path that would be created

**Default:**
```yaml
format_dry_run_archive: "Would create archive: %s\n"
```

**Examples:**
```yaml
# With emoji
format_dry_run_archive: "🔍 [DRY RUN] Would create: %s\n"

# More explicit
format_dry_run_archive: "[DRY RUN] Archive would be created at: %s\n"

# Minimal
format_dry_run_archive: "[DRY RUN] %s\n"
```

---

### format_error

**YAML Field:** `format_error`

**Description:** Generic error message format.

**Placeholders:**
- `%s` - Error message

**Default:**
```yaml
format_error: "Error: %s\n"
```

**Examples:**
```yaml
# With emoji
format_error: "❌ Error: %s\n"

# More formal
format_error: "ERROR: %s\n"

# Minimal
format_error: "%s\n"
```

---

## File Backup Format Strings

### format_created_backup

**YAML Field:** `format_created_backup`

**Description:** Message displayed when a file backup is successfully created.

**Placeholders:**
- `%s` - Backup file path

**Default:**
```yaml
format_created_backup: "Created backup: %s\n"
```

**Examples:**
```yaml
format_created_backup: "✅ Backup created: %s\n"
format_created_backup: "File backed up to: %s\n"
```

---

### format_identical_backup

**YAML Field:** `format_identical_backup`

**Description:** Message displayed when file is identical to existing backup.

**Placeholders:**
- `%s` - Backup file path

**Default:**
```yaml
format_identical_backup: "File is identical to existing backup: %s\n"
```

**Examples:**
```yaml
format_identical_backup: "⚠️  File unchanged. Identical to: %s\n"
format_identical_backup: "No changes since: %s\n"
```

---

### format_list_backup

**YAML Field:** `format_list_backup`

**Description:** Format for listing file backups.

**Placeholders:**
- `%s` - Backup file path
- `#{size_human}` - Human-readable file size

**Default:**
```yaml
format_list_backup: "%s (size: #{size_human})\n"
```

**Examples:**
```yaml
format_list_backup: "📄 %s (#{size_human})\n"
format_list_backup: "Backup: %s | Size: #{size_human}\n"
```

---

### format_dry_run_backup

**YAML Field:** `format_dry_run_backup`

**Description:** Message displayed during dry-run mode for file backups.

**Placeholders:**
- `%s` - Backup file path that would be created

**Default:**
```yaml
format_dry_run_backup: "Would create backup: %s\n"
```

**Examples:**
```yaml
format_dry_run_backup: "🔍 [DRY RUN] Would backup: %s\n"
format_dry_run_backup: "[DRY RUN] Backup would be created at: %s\n"
```

---

## Template-Based Format Strings

Template-based format strings use named placeholders for better readability.

### template_created_archive

**YAML Field:** `template_created_archive`

**Description:** Template version of `format_created_archive`.

**Placeholders:**
- `#{path}` - Archive path

**Default:**
```yaml
template_created_archive: "Created archive: #{path}\n"
```

---

### template_identical_archive

**YAML Field:** `template_identical_archive`

**Description:** Template version of `format_identical_archive`.

**Placeholders:**
- `#{path}` - Archive path

**Default:**
```yaml
template_identical_archive: "Directory is identical to existing archive: #{path}\n"
```

---

### template_list_archive

**YAML Field:** `template_list_archive`

**Description:** Template version of `format_list_archive`.

**Placeholders:**
- `#{path}` - Archive path
- `#{size_human}` - Human-readable archive size

**Default:**
```yaml
template_list_archive: "#{path} (size: #{size_human})\n"
```

---

### template_config_value

**YAML Field:** `template_config_value`

**Description:** Template version of `format_config_value`.

**Placeholders:**
- `#{name}` - Configuration name
- `#{value}` - Configuration value
- `#{source}` - Source (where the value came from)

**Default:**
```yaml
template_config_value: "#{name}: #{value} (source: #{source})\n"
```

---

### template_dry_run_archive

**YAML Field:** `template_dry_run_archive`

**Description:** Template version of `format_dry_run_archive`.

**Placeholders:**
- `#{path}` - Archive path that would be created

**Default:**
```yaml
template_dry_run_archive: "Would create archive: #{path}\n"
```

---

### template_error

**YAML Field:** `template_error`

**Description:** Template version of `format_error`.

**Placeholders:**
- `#{message}` - Error message

**Default:**
```yaml
template_error: "Error: #{message}\n"
```

---

## Extended Format Strings

### Archive Operations

#### format_no_archives_found

**YAML Field:** `format_no_archives_found`

**Placeholders:** `%s` - Archive directory path

**Default:** `"No archives found in %s\n"`

---

#### format_verification_failed

**YAML Field:** `format_verification_failed`

**Placeholders:** `%s` - Archive name, `%v` - Error

**Default:** `"Archive %s verification failed: %v\n"`

---

#### format_verification_success

**YAML Field:** `format_verification_success`

**Placeholders:** `%s` - Archive name

**Default:** `"Archive %s verified successfully\n"`

---

#### format_incremental_created

**YAML Field:** `format_incremental_created`

**Placeholders:** `%s` - Archive path

**Default:** `"Created incremental archive: %s\n"`

---

### Detailed Output with Statistics

#### format_created_archive_detailed

**YAML Field:** `format_created_archive_detailed`

**Placeholders:** `%s` - Path, `%s` - Size, `%s` - Modified time

**Default:** `"Created archive: %s (%s, %s)\n"`

---

#### format_incremental_created_detailed

**YAML Field:** `format_incremental_created_detailed`

**Placeholders:** `%s` - Path, `%s` - Size, `%s` - Modified time

**Default:** `"Created incremental archive: %s (%s, %s)\n"`

---

### Error Messages

#### format_disk_full_error

**YAML Field:** `format_disk_full_error`

**Placeholders:** `%v` - Error details

**Default:** `"Disk full error: %v\n"`

---

#### format_permission_error

**YAML Field:** `format_permission_error`

**Placeholders:** `%v` - Error details

**Default:** `"Permission error: %v\n"`

---

#### format_directory_not_found

**YAML Field:** `format_directory_not_found`

**Placeholders:** `%v` - Error details

**Default:** `"Directory not found: %v\n"`

---

#### format_file_not_found

**YAML Field:** `format_file_not_found`

**Placeholders:** `%v` - Error details

**Default:** `"File not found: %v\n"`

---

## Complete Example

Here's a complete example showing customization of all common format strings:

```yaml
# Basic settings
archive_dir_path: "~/Archives"
use_current_dir_name: true

# Directory operation format strings with emoji
format_created_archive: "✅ Archive created: %s\n"
format_identical_archive: "⚠️  No changes detected. Identical to: %s\n"
format_list_archive: "📦 %s (size: #{size_human})\n"
format_config_value: "⚙️  %s: %s (from: %s)\n"
format_dry_run_archive: "🔍 [DRY RUN] Would create: %s\n"
format_error: "❌ Error: %s\n"

# File backup format strings
format_created_backup: "✅ Backup created: %s\n"
format_identical_backup: "⚠️  File unchanged. Identical to: %s\n"
format_list_backup: "📄 %s (size: #{size_human})\n"
format_dry_run_backup: "🔍 [DRY RUN] Would backup: %s\n"

# Template-based format strings
template_created_archive: "✅ Created: #{path}\n"
template_identical_archive: "⚠️  Identical: #{path}\n"
template_list_archive: "📦 #{path} (#{size_human})\n"
template_config_value: "⚙️  #{name}: #{value} (source: #{source})\n"
template_dry_run_archive: "🔍 Would create: #{path}\n"
template_error: "❌ #{message}\n"
```

## Tips and Best Practices

1. **Keep it consistent**: Use similar formatting across related messages
2. **Include newlines**: Most format strings should end with `\n` for proper line breaks
3. **Test your changes**: Run `bkpdir` commands to verify your custom formats work as expected
4. **Use emoji sparingly**: While emoji can make output more visual, too many can be distracting
5. **Consider automation**: If you're parsing output programmatically, use simple formats without emoji

## Validation

When you load a configuration file with custom format strings, bkpdir will validate that you're using the correct placeholders. If you use an unexpected placeholder, you'll see a warning like:

```
Warning: Field 'format_created_archive': unexpected placeholder '#{path}'. 
Expected one of: [%s]
Hint: Use printf-style placeholders for this field, or use 'template_created_archive' for template-style formatting.
```

These warnings are informational and won't prevent your configuration from loading.

## See Also

- [Configuration Management](../README.md#configuration)
- [Example Configuration Files](../example-custom-formats.yml)
- [Output formatting (REQ-OUTPUT_FORMATTING)](../../tied/requirements/REQ-OUTPUT_FORMATTING.yaml)
