# Package Interdependency Mapping — EXTRACT-008

**Tokens:** [REQ:EXTRACT_008_INTERDEP_MAPPING] [ARCH:EXTRACT_008_INTERDEP] [IMPL:EXTRACT_008_DOC_MIGRATION]

> **[HIGH] EXTRACT-008: Package interdependency mapping - 📊 Clear usage patterns for extracted packages**

## 📑 Executive Summary

This document provides comprehensive guidance for integrating the 8 extracted packages from the BkpDir CLI application. These packages form a complete toolkit for building robust CLI applications with configuration management, error handling, file operations, and more.

### 🎯 Quick Reference

| Package | Purpose | Size | External Dependencies | Internal Dependencies |
|---------|---------|------|---------------------|---------------------|
| **pkg/config** | Configuration management | 1,322 lines | yaml.v3 | None |
| **pkg/errors** | Error handling & classification | 918 lines | None | None |
| **pkg/resources** | Resource management & cleanup | 486 lines | None | None |
| **pkg/formatter** | Output formatting & templates | 1,056 lines | None | None |
| **pkg/git** | Git integration utilities | 255 lines | None | None |
| **pkg/cli** | CLI command framework | 722 lines | cobra | None |
| **pkg/fileops** | File operations & utilities | 1,173 lines | doublestar | None |
| **pkg/processing** | Data processing patterns | 1,705 lines | None | None |

**Total Extracted Code**: 7,637 lines across 33 files

## 📊 Package Overview

### [ACTION:core-functionality] pkg/config - Configuration Management System
**Lines**: 1,322 | **Files**: 4 | **Maturity**: Production Ready

**Purpose**: Schema-agnostic configuration loading, merging, and validation
**Key Features**:
- Multi-source configuration (files, env vars, defaults)
- Schema-agnostic design using reflection
- Source tracking for debugging
- Environment variable override support

**Files**:
- `interfaces.go` (195 lines) - Core configuration interfaces
- `discovery.go` (212 lines) - Configuration file discovery
- `loader.go` (580 lines) - Generic configuration loading engine
- `utils.go` (339 lines) - Supporting utilities and implementations

**External Dependencies**: `gopkg.in/yaml.v3`
**Performance**: 24.3μs per configuration load operation

**Semantic Tokens**: [REQ:CONFIGURATION] [ARCH:CONFIG_SYSTEM] [IMPL:CONFIG_STRUCT]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Configuration package analysis

### [ACTION:validation] pkg/errors - Error Handling & Classification
**Lines**: 918 | **Files**: 3 | **Maturity**: Production Ready

**Purpose**: Structured error handling with context and classification
**Key Features**:
- Generic ApplicationError type
- Error classification framework
- Context preservation
- Recovery patterns

**Files**:
- `interfaces.go` - Error interfaces and contracts
- `classification.go` - Error type classification
- `handlers.go` - Error handling utilities

**External Dependencies**: None
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:ERROR_HANDLING] [ARCH:ERROR_HANDLING] [IMPL:STRUCTURED_ERRORS]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Error handling package analysis

### 🗂️ pkg/resources - Resource Management & Cleanup
**Lines**: 486 | **Files**: 2 | **Maturity**: Production Ready

**Purpose**: Thread-safe resource management with automatic cleanup
**Key Features**:
- ResourceManager with panic recovery
- Context-aware operations
- Atomic file operations
- Resource lifecycle management

**Files**:
- `interfaces.go` - Resource management interfaces
- `context.go` - Context-aware resource operations

**External Dependencies**: None
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:RESOURCE_MANAGEMENT] [ARCH:RESOURCE_MANAGEMENT] [IMPL:RESOURCE_MANAGER]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Resource management package analysis

### 🎨 pkg/formatter - Output Formatting System
**Lines**: 1,056 | **Files**: 5 | **Maturity**: Production Ready

**Purpose**: Comprehensive output formatting with templates and patterns
**Key Features**:
- Template-based formatting
- Printf-style formatting
- Pattern extraction engine
- Output collection system

**Files**:
- `interfaces.go` - Formatter interfaces
- `formatter.go` - Main formatter implementation
- `template.go` - Template processing engine
- `collector.go` - Output collection system
- `patterns.go` - Pattern extraction utilities

**External Dependencies**: None
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:OUTPUT_FORMATTING] [ARCH:OUTPUT_FORMATTING] [IMPL:DUAL_FORMATTING]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Output formatting package analysis

### 🔀 pkg/git - Git Integration Utilities
**Lines**: 255 | **Files**: 1 | **Maturity**: Production Ready

**Purpose**: Git repository detection and metadata extraction
**Key Features**:
- Repository detection
- Branch and hash extraction
- Working directory status
- Configuration support

**Files**:
- `git.go` - Complete Git integration implementation

**External Dependencies**: None (uses command-line git)
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:GIT_INTEGRATION] [ARCH:GIT_INTEGRATION] [IMPL:GIT_CLI]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Git integration package analysis

### ⚡ pkg/cli - CLI Command Framework
**Lines**: 722 | **Files**: 6 | **Maturity**: Production Ready

**Purpose**: CLI command patterns and framework utilities
**Key Features**:
- Command patterns and builders
- Dry-run support
- Context-aware execution
- Version handling

**Files**:
- `types.go` - Core CLI types and interfaces
- `builder.go` - Command builder patterns
- `context.go` - Context management
- `dryrun.go` - Dry-run implementation
- `flags.go` - Flag handling utilities
- `version.go` - Version command support

**External Dependencies**: `github.com/spf13/cobra`
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:IMMUTABLE_CLI_COMMANDS] [ARCH:CLI_COMMANDS] [IMPL:CLI_FRAMEWORK]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] CLI framework package analysis

### 📁 pkg/fileops - File Operations & Utilities
**Lines**: 1,173 | **Files**: 6 | **Maturity**: Production Ready

**Purpose**: Comprehensive file system operations and utilities
**Key Features**:
- File comparison and validation
- Pattern-based exclusion
- Atomic operations
- Directory traversal

**Files**:
- `fileops.go` - Package documentation and interfaces
- `comparison.go` - File and directory comparison
- `exclusion.go` - Pattern-based file exclusion
- `validation.go` - Path validation and security
- `atomic.go` - Atomic file operations
- `traversal.go` - Safe directory traversal

**External Dependencies**: `github.com/bmatcuk/doublestar/v4`
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:FILE_BACKUP] [ARCH:FILE_OPERATIONS] [IMPL:FILE_OPERATIONS]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] File operations package analysis

### ⚙️ pkg/processing - Data Processing Patterns
**Lines**: 1,705 | **Files**: 6 | **Maturity**: Production Ready

**Purpose**: Reusable data processing workflows and patterns
**Key Features**:
- Concurrent processing patterns
- Pipeline processing
- Naming conventions
- Verification systems

**Files**:
- `doc.go` - Package documentation
- `processor.go` - Core processing interfaces
- `pipeline.go` - Pipeline processing patterns
- `concurrent.go` - Concurrent processing utilities
- `naming.go` - Naming convention utilities
- `verification.go` - Data verification patterns

**External Dependencies**: None
**Internal Dependencies**: None

**Semantic Tokens**: [REQ:PERFORMANCE] [ARCH:PROCESSING_PATTERNS] [IMPL:PROCESSING_PATTERNS]

// [HIGH] EXTRACT-008: Package interdependency mapping - [ACTION:discovery] Data processing package analysis 

## 🗺️ Dependency Matrix

### 📊 Package Relationship Overview

The extracted packages are designed with **zero circular dependencies** and minimal coupling. This matrix shows the relationships:

```
                config  errors  resources  formatter  git  cli  fileops  processing
config            -       -        -          -       -    -      -         -
errors            -       -        -          -       -    -      -         -
resources         -       -        -          -       -    -      -         -
formatter         -       -        -          -       -    -      -         -
git               -       -        -          -       -    -      -         -
cli               -       -        -          -       -    -      -         -
fileops           -       -        -          -       -    -      -         -
processing        -       -        -          -       -    -      -         -
```

**Legend**: `-` = No dependency | `✓` = Direct dependency | `○` = Optional integration

### 🔗 External Dependencies Summary

| Package | External Dependencies | Purpose |
|---------|---------------------|---------|
| **config** | `gopkg.in/yaml.v3` | YAML configuration file parsing |
| **cli** | `github.com/spf13/cobra` | Command-line interface framework |
| **fileops** | `github.com/bmatcuk/doublestar/v4` | Advanced glob pattern matching |
| **All Others** | None | Self-contained implementations |

// [HIGH] EXTRACT-008: Package interdependency mapping - 📊 Dependency matrix creation

## 🎯 Usage Pattern Catalog

### [ACTION:core-functionality] Single Package Usage Patterns

#### Configuration Management (pkg/config)
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Configuration usage pattern
import "bkpdir/pkg/config"

// Basic configuration loading
loader := config.NewGenericConfigLoader()
cfg, err := loader.LoadConfig(".", &MyConfig{})

// With source tracking
values, err := loader.LoadConfigValues(".", &MyConfig{})
for _, value := range values {
    fmt.Printf("%s: %v (from %s)\n", value.Name, value.Value, value.Source)
}
```

#### Error Handling (pkg/errors)
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Error handling usage pattern
import "bkpdir/pkg/errors"

// Create application-specific error
err := errors.NewApplicationError("operation", "failed to process", originalErr)

// Error classification
if errors.IsDiskFullError(err) {
    // Handle disk space issues
}
```

#### Resource Management (pkg/resources)
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Resource management usage pattern
import "bkpdir/pkg/resources"

// Create resource manager
rm := resources.NewResourceManager()
defer rm.CleanupWithPanicRecovery()

// Create managed temporary file
tempFile, err := rm.CreateTempFile("prefix", ".tmp")
// Automatic cleanup on defer
```

### [ACTION:migration] Multi-Package Integration Patterns

#### Configuration + Formatter Pattern
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Config+Formatter integration pattern
import (
    "bkpdir/pkg/config"
    "bkpdir/pkg/formatter"
)

// Load configuration
loader := config.NewGenericConfigLoader()
cfg, err := loader.LoadConfig(".", &AppConfig{})

// Create formatter with config
formatter := formatter.NewTemplateFormatter(cfg)
output := formatter.FormatTemplate("Status: #{status}", data)
```

#### Errors + Resources Coordination
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Error+Resource coordination pattern
import (
    "bkpdir/pkg/errors"
    "bkpdir/pkg/resources"
)

func ProcessWithCleanup() error {
    rm := resources.NewResourceManager()
    defer rm.CleanupWithPanicRecovery()
    
    // Create resources
    tempFile, err := rm.CreateTempFile("work", ".tmp")
    if err != nil {
        return errors.NewApplicationError("setup", "failed to create temp file", err)
    }
    
    // Process with automatic cleanup
    return processFile(tempFile)
}
```

#### Complete CLI Application Assembly
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Complete CLI assembly pattern
import (
    "bkpdir/pkg/cli"
    "bkpdir/pkg/config"
    "bkpdir/pkg/errors"
    "bkpdir/pkg/formatter"
    "bkpdir/pkg/git"
    "bkpdir/pkg/resources"
)

func BuildCLIApp() *cobra.Command {
    // Initialize components
    configLoader := config.NewGenericConfigLoader()
    formatter := formatter.NewTemplateFormatter(nil)
    gitInfo := git.NewRepository(".")
    
    // Build command with all components
    cmd := cli.NewCommandBuilder("myapp").
        WithConfig(configLoader).
        WithFormatter(formatter).
        WithGitIntegration(gitInfo).
        Build()
    
    return cmd
}
```

// [HIGH] EXTRACT-008: Package interdependency mapping - 🎯 Usage pattern documentation

## 💡 Integration Examples

### 🚀 Basic CLI Application
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 💡 Basic CLI integration example
package main

import (
    "context"
    "fmt"
    "os"
    
    "bkpdir/pkg/cli"
    "bkpdir/pkg/config"
    "bkpdir/pkg/errors"
    "bkpdir/pkg/formatter"
    "github.com/spf13/cobra"
)

type AppConfig struct {
    OutputFormat string `yaml:"output_format"`
    Verbose      bool   `yaml:"verbose"`
}

func main() {
    ctx := context.Background()
    
    // Initialize configuration
    loader := config.NewGenericConfigLoader()
    cfg := &AppConfig{
        OutputFormat: "text",
        Verbose:      false,
    }
    
    // Load configuration
    loadedCfg, err := loader.LoadConfig(".", cfg)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Config error: %v\n", err)
        os.Exit(1)
    }
    
    appCfg := loadedCfg.(*AppConfig)
    
    // Initialize formatter
    formatter := formatter.NewTemplateFormatter(appCfg)
    
    // Create root command
    rootCmd := &cobra.Command{
        Use:   "myapp",
        Short: "Example CLI application",
        Run: func(cmd *cobra.Command, args []string) {
            output := formatter.FormatTemplate("App running in #{format} mode", map[string]interface{}{
                "format": appCfg.OutputFormat,
            })
            fmt.Println(output)
        },
    }
    
    // Execute
    if err := rootCmd.ExecuteContext(ctx); err != nil {
        appErr := errors.NewApplicationError("execution", "command failed", err)
        fmt.Fprintf(os.Stderr, "Error: %v\n", appErr)
        os.Exit(1)
    }
}
```

### [ACTION:migration] File Processing Pipeline
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - 💡 File processing pipeline example
import (
    "bkpdir/pkg/fileops"
    "bkpdir/pkg/processing"
    "bkpdir/pkg/resources"
    "bkpdir/pkg/errors"
)

func ProcessFiles(inputDir, outputDir string) error {
    // Resource management
    rm := resources.NewResourceManager()
    defer rm.CleanupWithPanicRecovery()
    
    // File operations
    comparer := fileops.NewDefaultComparer()
    excluder := fileops.NewPatternMatcher([]string{"*.tmp", "*.log"})
    
    // Processing pipeline
    processor := processing.NewConcurrentProcessor(4) // 4 workers
    
    // Traverse input directory
    traverser := fileops.NewTraverser(excluder)
    files, err := traverser.TraverseDirectory(inputDir)
    if err != nil {
        return errors.NewApplicationError("traversal", "failed to traverse input", err)
    }
    
    // Process files concurrently
    results := processor.ProcessBatch(files, func(file string) error {
        // Create temporary output file
        tempFile, err := rm.CreateTempFile("process", ".tmp")
        if err != nil {
            return err
        }
        
        // Process file (example: copy with validation)
        return processFile(file, tempFile, outputDir)
    })
    
    // Handle results
    for _, result := range results {
        if result.Error != nil {
            fmt.Printf("Failed to process %s: %v\n", result.Input, result.Error)
        }
    }
    
    return nil
}
```

// [HIGH] EXTRACT-008: Package interdependency mapping - 💡 Integration examples creation

## ⚡ Performance Considerations

### 📊 Resource Usage Patterns

| Package | Memory Impact | CPU Impact | I/O Impact | Concurrency Safe |
|---------|---------------|------------|------------|------------------|
| **config** | Low (24.3μs/load) | Low | Medium (file reads) | Yes |
| **errors** | Minimal | Minimal | None | Yes |
| **resources** | Low (tracking overhead) | Low | High (cleanup) | Yes |
| **formatter** | Medium (templates) | Medium | Low | Yes |
| **git** | Low | Medium (exec) | Medium (git commands) | Yes |
| **cli** | Low | Low | Low | Yes |
| **fileops** | Medium (buffers) | Medium | High (file ops) | Yes |
| **processing** | High (workers) | High | Variable | Yes |

### [ACTION:core-functionality] Optimization Guidelines

#### Memory Optimization
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - ⚡ Memory optimization patterns
// Reuse formatters and config loaders
var (
    globalFormatter = formatter.NewTemplateFormatter(config)
    globalLoader    = config.NewGenericConfigLoader()
)

// Pool resources for high-throughput scenarios
resourcePool := sync.Pool{
    New: func() interface{} {
        return resources.NewResourceManager()
    },
}
```

#### Concurrent Processing
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - ⚡ Concurrent processing optimization
// Use processing package for CPU-intensive work
processor := processing.NewConcurrentProcessor(runtime.NumCPU())

// Coordinate with context for cancellation
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

results := processor.ProcessBatchWithContext(ctx, items, workerFunc)
```

#### Resource Management
```go
// [HIGH] EXTRACT-008: Package interdependency mapping - ⚡ Resource management optimization
// Use single resource manager per operation
func OptimizedOperation() error {
    rm := resources.NewResourceManager()
    defer rm.CleanupWithPanicRecovery()
    
    // Batch resource creation
    tempFiles := make([]*os.File, 0, 10)
    for i := 0; i < 10; i++ {
        file, err := rm.CreateTempFile("batch", ".tmp")
        if err != nil {
            return err
        }
        tempFiles = append(tempFiles, file)
    }
    
    // Process all files
    return processBatch(tempFiles)
}
```

// [HIGH] EXTRACT-008: Package interdependency mapping - ⚡ Performance analysis and optimization

## 📋 Best Practices

### ✅ Recommended Patterns

1. **Configuration First**: Always initialize configuration before other components
2. **Resource Management**: Use ResourceManager for any temporary resources
3. **Error Context**: Wrap errors with ApplicationError for better debugging
4. **Interface Usage**: Depend on interfaces, not concrete implementations
5. **Context Propagation**: Pass context through all operations for cancellation

### ❌ Anti-Patterns to Avoid

1. **Direct File Operations**: Use pkg/fileops instead of direct os package calls
2. **Ignoring Errors**: Always handle errors with pkg/errors patterns
3. **Resource Leaks**: Always use defer with ResourceManager cleanup
4. **Blocking Operations**: Use context for cancellation in long-running operations
5. **Hardcoded Paths**: Use pkg/config for all configuration values

### [ACTION:core-functionality] Integration Checklist

- [ ] Configuration loaded before component initialization
- [ ] ResourceManager created for temporary resources
- [ ] Errors wrapped with ApplicationError
- [ ] Context passed to all operations
- [ ] Cleanup deferred for all resources
- [ ] Interfaces used for component dependencies
- [ ] Tests written for integration scenarios

// [HIGH] EXTRACT-008: Package interdependency mapping - 📋 Best practices documentation

## 🚨 Troubleshooting Guide

### Common Integration Issues

#### Configuration Loading Failures
```go
// Problem: Configuration not found
// Solution: Check search paths and file permissions
loader := config.NewGenericConfigLoader()
paths := loader.GetConfigSearchPaths()
fmt.Printf("Searching in: %v\n", paths)
```

#### Resource Cleanup Issues
```go
// Problem: Resources not cleaned up
// Solution: Always use defer with panic recovery
rm := resources.NewResourceManager()
defer rm.CleanupWithPanicRecovery() // This handles panics too
```

#### Performance Issues
```go
// Problem: Slow file operations
// Solution: Use concurrent processing
processor := processing.NewConcurrentProcessor(runtime.NumCPU())
// Process files in parallel
```

#### Error Context Loss
```go
// Problem: Losing error context
// Solution: Wrap with ApplicationError
if err != nil {
    return errors.NewApplicationError("operation", "context description", err)
}
```

### [ACTION:discovery] Debugging Tips

1. **Enable Verbose Logging**: Use formatter for structured output
2. **Check Resource Usage**: Monitor ResourceManager for leaks
3. **Validate Configuration**: Use config source tracking
4. **Profile Performance**: Use processing package benchmarks
5. **Test Integration**: Write integration tests for package combinations

// [HIGH] EXTRACT-008: Package interdependency mapping - 📋 Troubleshooting guide

## 🎯 Conclusion

The 8 extracted packages provide a comprehensive foundation for building robust CLI applications. With **zero circular dependencies** and **minimal coupling**, they can be used individually or in combination to accelerate development.

### 📈 Adoption Strategy

1. **Start Simple**: Begin with pkg/config and pkg/cli for basic applications
2. **Add Robustness**: Integrate pkg/errors and pkg/resources for production use
3. **Enhance UX**: Add pkg/formatter for better user experience
4. **Scale Up**: Use pkg/processing and pkg/fileops for complex operations
5. **Version Control**: Integrate pkg/git for development tools

### 🚀 Next Steps

- Review the [CLI Template](../cmd/cli-template/) for working examples
- Check the [Migration Guide](migration-guide.md) for transitioning existing applications
- Explore individual package documentation for detailed APIs
- Consider contributing improvements back to the extracted packages

## ✅ Deliverables (EXTRACT-008) — Migration status

- `docs/package-interdependency-mapping.md` — Canonical interdependency mapping and narrative (this document) ✅ Completed
- `docs/images/package-interdependency-mapping.svg` — Visual diagram illustrating package relationships ✅ Completed
- `pkg/*/README.md` — Package READMEs exist and were reviewed for public interface documentation (pkg/config, pkg/errors, pkg/resources, pkg/formatter, pkg/git, pkg/cli, pkg/fileops, pkg/processing, pkg/testutil) ✅ Reviewed
- `docs/examples/` — Example integrations reviewed (basic-cli-app, git-aware-backup); pipeline example validated in narrative ✅ Reviewed
- `stdd/` tokens — `[REQ:EXTRACT_008_INTERDEP_MAPPING]`, `[ARCH:EXTRACT_008_INTERDEP]`, `[IMPL:EXTRACT_008_DOC_MIGRATION]` present in `stdd/semantic-tokens.md` and cross-referenced ✅ Validated

## 🔁 Remaining Actions

- Add CI check for import-cycle detection and diagram presence (`go list -deps` / custom script) — Pending (postponed until GitHub integration)
- Add automated token cross-reference validation into CI — Pending
- Final verification run of token validation scripts (`scripts/validate-token-traceability.sh`) after CI check added — Pending

**Migration Note**: Per `[IMPL:EXTRACT_008_DOC_MIGRATION]`, the working-plan artifacts were reviewed and their content consolidated into STDD canonical documents. No supplemental working-plan files remain as active canonical sources.

**[HIGH] EXTRACT-008: Package interdependency mapping - 📊 Complete package interdependency mapping with comprehensive usage patterns and integration guidance**
