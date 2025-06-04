# Git-Aware Backup Tool Example

// EXTRACT-010: Git-aware backup tool integration example - Demonstrates config + cli + git + fileops package integration - 🔺

This example demonstrates how to build a Git-aware backup tool using the extracted BkpDir packages:

- **pkg/config**: Multi-source configuration management
- **pkg/cli**: CLI framework with context support
- **pkg/git**: Git repository integration and metadata extraction
- **pkg/fileops**: File operations with atomic operations and exclusion patterns
- **pkg/processing**: Naming conventions with Git information
- **pkg/formatter**: Output formatting and progress reporting

## Features Demonstrated

### Git Repository Integration
- Git repository detection and validation
- Branch and commit hash extraction
- Dirty state detection and handling
- Git-aware naming conventions

### File Operations
- Directory listing with pattern matching
- File exclusion and inclusion patterns
- Atomic file operations
- Path validation and creation

### Processing Patterns
- Timestamp-based naming with Git metadata
- Backup name generation and parsing
- Progress tracking and reporting
- Error handling with context

### Configuration Management
- Multi-source configuration loading
- Environment variable integration
- Default value handling
- Configuration validation

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

The tool loads configuration from multiple sources in priority order:

### Environment Variables

```bash
export GITBACKUP_SOURCE_PATH=/path/to/project
export GITBACKUP_OUTPUT_DIRECTORY=/backups
export GITBACKUP_OUTPUT_FORMAT=json
export GITBACKUP_GIT_INCLUDE_HASH=true
export GITBACKUP_GIT_MAX_HASH_LENGTH=8
```

### Configuration File

Create a `backup-config.yaml` file:

```yaml
# Source and output configuration
source_path: ./myproject
output_directory: ./backups
output_format: table

# Git configuration
git:
  include_hash: true
  max_hash_length: 8
  include_branch: true
  include_dirty_state: true
  require_clean: false

# File operations
file_ops:
  ignore_patterns:
    - ".git"
    - "node_modules"
    - "*.tmp"
    - "*.log"
  follow_symlinks: false
  preserve_metadata: true

# Output configuration
output:
  progress: true
  verbose: false
```

## Usage Examples

### Creating Git-Aware Backups

```bash
# Basic backup creation
./git-backup create

# Create backup from specific directory
GITBACKUP_SOURCE_PATH=/path/to/project ./git-backup create

# Create backup with JSON output
./git-backup create --format=json
```

### Listing Backups with Git Information

```bash
# List all backups
./git-backup list

# List backups in JSON format
./git-backup list --format=json

# List backups in YAML format
./git-backup list --format=yaml
```

### Configuration Management

```bash
# Show version information
./git-backup version

# Run with custom config file
./git-backup --config=custom-backup.yaml create
```

## Example Output

### Git Repository Detection

```
Creating Git-aware backup...
Git Repository Found:
  Branch: feature-branch
  Hash: abc123de
  Clean: false
Backup will be created: ./backups/myproject-2024-01-02T120000-feature-branch-abc123de-dirty.zip
Found 156 files to backup
Creating backup...
Progress: 100.0% (156/156) README.md

✅ Backup created successfully!
Name: myproject-2024-01-02T120000-feature-branch-abc123de-dirty.zip
Files: 156
Git: feature-branch@abc123de (clean: false)
```

### JSON Output Format

```json
{
  "backup_name": "myproject-2024-01-02T120000-feature-branch-abc123de-dirty.zip",
  "source_path": "./myproject",
  "file_count": 156,
  "created_at": "2024-01-02T12:00:00Z",
  "git": {
    "branch": "feature-branch",
    "hash": "abc123def456789",
    "is_clean": false
  }
}
```

### Backup Listing

```
Name                                                           Branch         Hash     Size     Created
----                                                           ------         ----     ----     -------
myproject-2024-01-02T120000-main-abc123de.zip                main           abc123de 2.5 MB   2024-01-02 12:00:00
myproject-2024-01-02T110000-feature-branch-def456ab-dirty.zip feature-branch def456ab 2.3 MB   2024-01-02 11:00:00
myproject-2024-01-01T180000-main-789cdef1.zip                main           789cdef1 2.4 MB   2024-01-01 18:00:00
```

## Git-Aware Naming Convention

The tool generates backup names that include Git information:

### Naming Pattern

```
{project-name}-{timestamp}-{git-branch}-{git-hash}[-dirty].zip
```

### Examples

- `myproject-2024-01-02T120000-main-abc123de.zip` - Clean repository on main branch
- `myproject-2024-01-02T120000-feature-branch-def456ab-dirty.zip` - Dirty repository on feature branch
- `myproject-2024-01-02T120000.zip` - Non-Git repository

### Parsing Backup Names

The tool can parse backup names to extract Git information:

```go
components, err := naming.ParseName(backupName, "archive")
if err == nil {
    fmt.Printf("Branch: %s\n", components.GitBranch)
    fmt.Printf("Hash: %s\n", components.GitHash)
    fmt.Printf("Timestamp: %s\n", components.Timestamp.Format(time.RFC3339))
}
```

## Integration Patterns Demonstrated

### 1. Git Repository Detection

```go
// Check if directory is a Git repository
isRepo, err := gitRepo.IsRepository(sourcePath)
if err != nil {
    return fmt.Errorf("failed to check Git repository: %w", err)
}

if isRepo {
    // Extract Git metadata
    gitInfo, err := gitRepo.GetRepositoryInfo(sourcePath)
    if err != nil {
        return fmt.Errorf("failed to get Git info: %w", err)
    }
}
```

### 2. Configuration-Driven File Operations

```go
// Load exclusion patterns from configuration
ignorePatterns := cfg.GetStringSlice("file_ops.ignore_patterns", []string{".git"})

// List files with configuration-based options
files, err := fileOps.ListFiles(sourcePath, fileops.ListOptions{
    IgnorePatterns: ignorePatterns,
    Recursive:      true,
})
```

### 3. Git-Aware Naming

```go
// Create naming template with Git information
template := &processing.NamingTemplate{
    Prefix:          filepath.Base(sourcePath),
    Timestamp:       time.Now(),
    GitBranch:       gitInfo.Branch,
    GitHash:         gitInfo.Hash[:8],
    GitIsClean:      gitInfo.IsClean,
    ShowGitDirtyStatus: true,
}

// Generate backup name
name, err := naming.GenerateName(template)
```

### 4. Progress Reporting

```go
// Report progress during file processing
for i, file := range files {
    progress := float64(i+1) / float64(len(files)) * 100
    fmt.Printf("\rProgress: %.1f%% (%d/%d) %s", 
        progress, i+1, len(files), filepath.Base(file))
}
```

### 5. Multi-Format Output

```go
// Format output based on configuration
switch format {
case "json":
    output, _ := formatter.FormatJSON(metadata)
    fmt.Println(output)
case "yaml":
    output, _ := formatter.FormatYAML(metadata)
    fmt.Println(output)
default:
    // Human-readable format
    fmt.Printf("✅ Backup created successfully!\n")
}
```

## Error Handling Patterns

### Git Repository Errors

```go
if isRepo {
    gitInfo, err := gitRepo.GetRepositoryInfo(sourcePath)
    if err != nil {
        // Handle different error scenarios
        if cfg.GetBool("git.require_clean", false) {
            return fmt.Errorf("failed to get Git info: %w", err)
        }
        // Continue without Git info if not critical
        gitInfo = nil
        if cfg.GetBool("output.verbose", false) {
            fmt.Printf("Warning: Could not get Git information: %v\n", err)
        }
    }
}
```

### File Operation Errors

```go
// Check if source exists before processing
if exists, err := fileOps.PathExists(sourcePath); err != nil {
    return fmt.Errorf("failed to check source path: %w", err)
} else if !exists {
    return errors.NewApplicationError(
        errors.CategoryFilesystem,
        errors.SeverityError,
        "SOURCE_NOT_FOUND",
        fmt.Sprintf("Source path does not exist: %s", sourcePath),
        map[string]interface{}{"source": sourcePath},
    )
}
```

## Package Integration Benefits

### Git-Aware Operations
- Automatically include Git context in backup names
- Detect repository state for validation
- Support both Git and non-Git projects

### Flexible Configuration
- Environment variable support for CI/CD integration
- File-based configuration for complex setups
- Sensible defaults for quick usage

### Robust File Handling
- Pattern-based file inclusion/exclusion
- Atomic operations for data integrity
- Cross-platform path handling

### Consistent Output
- Multiple output formats for different use cases
- Progress reporting for long operations
- Structured error messages

## Testing

Test the integration with different scenarios:

```bash
# Test with Git repository
cd /path/to/git/repo
./git-backup create

# Test with non-Git directory
cd /path/to/regular/dir
./git-backup create

# Test with dirty Git repository
cd /path/to/dirty/git/repo
echo "test" >> file.txt
./git-backup create

# Test configuration loading
GITBACKUP_OUTPUT_FORMAT=json ./git-backup create

# Test backup listing
./git-backup list --format=table
```

## Real-World Use Cases

### CI/CD Integration

```bash
# In CI/CD pipeline
export GITBACKUP_SOURCE_PATH=${CI_PROJECT_DIR}
export GITBACKUP_OUTPUT_DIRECTORY=/backups
export GITBACKUP_OUTPUT_FORMAT=json
./git-backup create
```

### Development Workflow

```bash
# Before major refactoring
./git-backup create --note="before-refactor"

# Before merging branches
git checkout feature-branch
./git-backup create --note="before-merge"
```

### Automated Backups

```bash
# Cron job for automated backups
0 2 * * * /usr/local/bin/git-backup create --config=/etc/backup-config.yaml
```

This example demonstrates a real-world application of the extracted packages, showing how they work together to create a sophisticated Git-aware backup tool with minimal code duplication and maximum reusability. 