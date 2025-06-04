# 📚 BkpDir Package Integration Guide

> **🔺 EXTRACT-008: Integration documentation system - 📝 Comprehensive usage guides**

## 🎯 Overview

This guide demonstrates how to integrate and use the extracted BkpDir packages in your CLI applications. The packages are designed to work independently or together, providing a flexible foundation for building backup and file management tools.

## 📦 Package Ecosystem

### Available Packages

| Package | Purpose | Key Features |
|---------|---------|--------------|
| **pkg/config** | Configuration management | Schema-agnostic loading, environment integration |
| **pkg/errors** | Error handling | Structured errors, context preservation |
| **pkg/resources** | Resource management | Cleanup coordination, lifecycle management |
| **pkg/formatter** | Output formatting | Template-based, printf-style formatting |
| **pkg/git** | Git integration | Repository detection, branch/commit info |
| **pkg/cli** | CLI framework | Command patterns, argument handling |
| **pkg/fileops** | File operations | Safe operations, comparison utilities |
| **pkg/processing** | Concurrent processing | Worker pools, pipeline patterns |

### Package Dependencies

```
pkg/config ──┐
             ├─→ pkg/errors ──┐
pkg/git ─────┤                ├─→ pkg/resources
             └─→ pkg/formatter ┘
                      │
pkg/cli ──────────────┼─→ pkg/fileops
                      │
pkg/processing ───────┘
```

## 🚀 Quick Start

### Basic CLI Application

Here's a minimal CLI application using the core packages:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Basic CLI example
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/yourusername/bkpdir/pkg/config"
    "github.com/yourusername/bkpdir/pkg/errors"
    "github.com/yourusername/bkpdir/pkg/cli"
    "github.com/yourusername/bkpdir/pkg/formatter"
)

func main() {
    // Initialize configuration
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
        os.Exit(1)
    }

    // Create CLI application
    app := cli.NewApp("myapp", "My CLI Application")
    
    // Add commands
    app.AddCommand(&cli.Command{
        Name:        "backup",
        Description: "Create a backup",
        Handler:     backupHandler(cfg),
    })

    // Run application
    if err := app.Run(context.Background(), os.Args); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

// 🔺 EXTRACT-008: Integration documentation system - 📝 Command implementation example
func backupHandler(cfg *config.Config) cli.HandlerFunc {
    return func(ctx context.Context, args []string) error {
        // Use formatter for output
        formatter := formatter.New(cfg.OutputFormat)
        
        // 🔺 EXTRACT-008: Integration documentation system - 📝 Command implementation example
        result := map[string]interface{}{
            "status": "success",
            "files":  42,
            "size":   "1.2GB",
        }
        
        return formatter.Print(result)
    }
}
```

## 🔧 Integration Patterns

### 1. Configuration-Driven Pattern

Use configuration to drive application behavior across all packages:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Configuration-driven pattern
type Application struct {
    config    *config.Config
    formatter *formatter.Formatter
    git       *git.Repository
    processor *processing.Pool
}

func NewApplication() (*Application, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, errors.Wrap(err, "failed to load configuration")
    }

    app := &Application{
        config:    cfg,
        formatter: formatter.New(cfg.OutputFormat),
        processor: processing.NewPool(cfg.MaxWorkers),
    }

    // Initialize Git if in repository
    if repo, err := git.Discover("."); err == nil {
        app.git = repo
    }

    return app, nil
}

func (a *Application) Process(files []string) error {
    ctx := context.Background()
    
    // Use processor for concurrent operations
    return a.processor.Process(ctx, files, func(file string) error {
        // Process individual file
        return a.processFile(file)
    })
}
```

### 2. Error Propagation Pattern

Maintain error context across package boundaries:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Error propagation pattern
func (a *Application) BackupFile(path string) error {
    // File operations with error context
    info, err := fileops.Stat(path)
    if err != nil {
        return errors.Wrap(err, "failed to stat file").
            WithField("path", path).
            WithField("operation", "backup")
    }

    // Git integration with error context
    if a.git != nil {
        status, err := a.git.FileStatus(path)
        if err != nil {
            return errors.Wrap(err, "failed to get git status").
                WithField("path", path).
                WithField("git_repo", a.git.Path())
        }
        
        if status.IsModified() {
            return errors.New("file has uncommitted changes").
                WithField("path", path).
                WithField("git_status", status.String())
        }
    }

    // Resource management with cleanup
    resource, err := resources.Acquire("backup_slot")
    if err != nil {
        return errors.Wrap(err, "failed to acquire backup resource")
    }
    defer resource.Release()

    // Perform backup operation
    return a.performBackup(path, info)
}
```

### 3. Output Formatting Pattern

Consistent output formatting across operations:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Output formatting pattern
type BackupResult struct {
    Path     string    `json:"path"`
    Size     int64     `json:"size"`
    Duration time.Duration `json:"duration"`
    GitInfo  *git.Info `json:"git_info,omitempty"`
}

func (a *Application) FormatResults(results []BackupResult) error {
    // Progress formatting
    // 🔺 EXTRACT-008: Integration documentation system - 📝 Progress formatting
    progress := map[string]interface{}{
        "completed": len(results),
        "total":     a.totalFiles,
        "percent":   float64(len(results)) / float64(a.totalFiles) * 100,
    }
    
    if err := a.formatter.PrintProgress(progress); err != nil {
        return errors.Wrap(err, "failed to format progress")
    }

    // Results formatting
    // 🔺 EXTRACT-008: Integration documentation system - 📝 Results formatting
    summary := map[string]interface{}{
        "results":    results,
        "total_size": a.calculateTotalSize(results),
        "duration":   a.totalDuration,
        "git_repo":   a.git != nil,
    }

    return a.formatter.Print(summary)
}
```

### 4. Git Integration Pattern

Leverage Git information across the application:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Git integration pattern
func (a *Application) EnhanceWithGitInfo(files []string) ([]FileInfo, error) {
    var enhanced []FileInfo
    
    for _, file := range files {
        info := FileInfo{Path: file}
        
        if a.git != nil {
            // Get Git status
            status, err := a.git.FileStatus(file)
            if err != nil {
                return nil, errors.Wrap(err, "failed to get git status").
                    WithField("file", file)
            }
            
            info.GitStatus = status.String()
            
            // Get commit information
            if !status.IsUntracked() {
                commit, err := a.git.LastCommit(file)
                if err == nil {
                    info.LastCommit = commit.Hash
                    info.LastAuthor = commit.Author
                }
            }
        }
        
        enhanced = append(enhanced, info)
    }
    
    return enhanced, nil
}

// 🔺 EXTRACT-008: Integration documentation system - 📝 Context display
func (a *Application) DisplayContext() error {
    if a.git == nil {
        return a.formatter.Print(map[string]interface{}{
            "context": "no git repository",
        })
    }

    info, err := a.git.Info()
    if err != nil {
        return errors.Wrap(err, "failed to get git info")
    }

    // 🔺 EXTRACT-008: Integration documentation system - 📝 Git-aware naming
    context := map[string]interface{}{
        "git_branch":     info.Branch,
        "git_commit":     info.Commit[:8],
        "git_dirty":      info.IsDirty,
        "git_repository": info.Repository,
    }

    return a.formatter.Print(context)
}
```

### 5. CLI Command Orchestration Pattern

Coordinate multiple packages within CLI commands:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 CLI orchestration pattern
func (a *Application) SetupCommands() {
    a.cli.AddCommand(&cli.Command{
        Name:        "backup",
        Description: "Create backup with Git awareness",
        Flags: []cli.Flag{
            {Name: "format", Value: "json", Description: "Output format"},
            {Name: "workers", Value: "4", Description: "Number of workers"},
        },
        Handler: a.backupCommand,
    })

    a.cli.AddCommand(&cli.Command{
        Name:        "status",
        Description: "Show backup status",
        Handler:     a.statusCommand,
    })
}

func (a *Application) backupCommand(ctx context.Context, args []string) error {
    // Parse flags and update configuration
    if format := cli.GetFlag(ctx, "format"); format != "" {
        a.formatter = formatter.New(format)
    }
    
    if workers := cli.GetIntFlag(ctx, "workers"); workers > 0 {
        a.processor = processing.NewPool(workers)
    }

    // Discover files
    files, err := a.discoverFiles(args)
    if err != nil {
        return errors.Wrap(err, "failed to discover files")
    }

    // Process with Git awareness
    return a.processWithGitAwareness(ctx, files)
}
```

### 6. Resource Management Pattern

Coordinate resource cleanup across packages:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Resource creation pattern
func (a *Application) ProcessWithResources(ctx context.Context, files []string) error {
    // Acquire processing resources
    procResource, err := resources.Acquire("processing")
    if err != nil {
        return errors.Wrap(err, "failed to acquire processing resource")
    }
    defer procResource.Release()

    // Acquire file system resources
    fsResource, err := resources.Acquire("filesystem")
    if err != nil {
        return errors.Wrap(err, "failed to acquire filesystem resource")
    }
    defer fsResource.Release()

    // Create resource-aware processor
    processor := processing.NewPoolWithResources(
        a.config.MaxWorkers,
        []resources.Resource{procResource, fsResource},
    )
    defer processor.Close()

    return processor.Process(ctx, files, a.processFile)
}
```

### 7. File Operations Pattern

Safe file operations with error handling and resource management:

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 File operations pattern
func (a *Application) SafeFileOperations(files []string) error {
    for _, file := range files {
        // Check file accessibility
        if !fileops.IsAccessible(file) {
            return errors.New("file not accessible").
                WithField("file", file)
        }

        // Compare with existing backup
        backupPath := a.getBackupPath(file)
        if fileops.Exists(backupPath) {
            same, err := fileops.Compare(file, backupPath)
            if err != nil {
                return errors.Wrap(err, "failed to compare files").
                    WithField("original", file).
                    WithField("backup", backupPath)
            }
            
            if same {
                a.formatter.PrintInfo("File unchanged, skipping: " + file)
                continue
            }
        }

        // Perform safe copy
        if err := fileops.SafeCopy(file, backupPath); err != nil {
            return errors.Wrap(err, "failed to copy file").
                WithField("source", file).
                WithField("destination", backupPath)
        }
    }

    return nil
}
```

## 🔄 Advanced Integration Patterns

### Concurrent Processing with Error Aggregation

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Concurrent processing pattern
func (a *Application) ProcessConcurrently(ctx context.Context, files []string) error {
    // Create error collector
    var errorCollector errors.Collector
    
    // Process files concurrently
    err := a.processor.ProcessWithCallback(ctx, files, 
        func(file string) error {
            return a.processFile(file)
        },
        func(file string, err error) {
            if err != nil {
                errorCollector.Add(errors.Wrap(err, "processing failed").
                    WithField("file", file))
            }
        },
    )

    // Check for processing errors
    if err != nil {
        return errors.Wrap(err, "concurrent processing failed")
    }

    // Return aggregated errors
    if errorCollector.HasErrors() {
        return errorCollector.Error()
    }

    return nil
}
```

## 🛡️ Error Handling Best Practices

### Structured Error Context

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Error handling patterns
func (a *Application) HandleOperationError(operation string, err error) error {
    // Wrap with operation context
    wrappedErr := errors.Wrap(err, "operation failed").
        WithField("operation", operation).
        WithField("timestamp", time.Now())

    // Add Git context if available
    if a.git != nil {
        if info, gitErr := a.git.Info(); gitErr == nil {
            wrappedErr = wrappedErr.
                WithField("git_branch", info.Branch).
                WithField("git_commit", info.Commit[:8])
        }
    }

    // Add configuration context
    wrappedErr = wrappedErr.
        WithField("config_file", a.config.Source()).
        WithField("output_format", a.config.OutputFormat)

    return wrappedErr
}

// 🔺 EXTRACT-008: Integration documentation system - 📝 Error context patterns
func (a *Application) RecoverableOperation(fn func() error) error {
    const maxRetries = 3
    
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }

        // Check if error is recoverable
        if !errors.IsRecoverable(err) {
            return errors.Wrap(err, "non-recoverable error").
                WithField("attempt", attempt)
        }

        // Log retry attempt
        a.formatter.PrintWarning(fmt.Sprintf(
            "Operation failed (attempt %d/%d): %v", 
            attempt, maxRetries, err,
        ))

        // Wait before retry
        time.Sleep(time.Duration(attempt) * time.Second)
    }

    return errors.New("operation failed after retries").
        WithField("max_retries", maxRetries)
}
```

## 🚀 Performance Optimization

### Memory Management

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Memory management patterns
func (a *Application) OptimizedProcessing(files []string) error {
    // Use buffered processing for large file sets
    const batchSize = 100
    
    for i := 0; i < len(files); i += batchSize {
        end := i + batchSize
        if end > len(files) {
            end = len(files)
        }
        
        batch := files[i:end]
        if err := a.processBatch(batch); err != nil {
            return errors.Wrap(err, "batch processing failed").
                WithField("batch_start", i).
                WithField("batch_size", len(batch))
        }
        
        // Force garbage collection between batches
        runtime.GC()
    }
    
    return nil
}

// 🔺 EXTRACT-008: Integration documentation system - 📝 Performance optimization patterns
func (a *Application) processBatch(files []string) error {
    // Create batch-specific resources
    batchResource, err := resources.AcquireWithTimeout("batch", 30*time.Second)
    if err != nil {
        return errors.Wrap(err, "failed to acquire batch resource")
    }
    defer batchResource.Release()

    // Process batch with dedicated worker pool
    batchProcessor := processing.NewPool(a.config.BatchWorkers)
    defer batchProcessor.Close()

    return batchProcessor.Process(context.Background(), files, a.processFile)
}
```

## 🔌 External Library Integration

### Database Integration

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Database integration patterns
import (
    "database/sql"
    _ "github.com/lib/pq"
)

func (a *Application) WithDatabase() error {
    // Use configuration for database connection
    dbURL := a.config.DatabaseURL
    if dbURL == "" {
        return errors.New("database URL not configured")
    }

    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return errors.Wrap(err, "failed to connect to database")
    }
    defer db.Close()

    // Store backup metadata
    return a.storeBackupMetadata(db)
}

func (a *Application) storeBackupMetadata(db *sql.DB) error {
    // Get Git information for metadata
    var gitInfo *git.Info
    if a.git != nil {
        info, err := a.git.Info()
        if err == nil {
            gitInfo = info
        }
    }

    // Insert backup record
    query := `
        INSERT INTO backups (timestamp, git_branch, git_commit, file_count)
        VALUES ($1, $2, $3, $4)
    `
    
    _, err := db.Exec(query, 
        time.Now(),
        gitInfo.Branch,
        gitInfo.Commit,
        a.processedFiles,
    )
    
    return errors.Wrap(err, "failed to store backup metadata")
}
```

### HTTP Client Integration

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 HTTP client integration patterns
func (a *Application) WithRemoteSync() error {
    client := &http.Client{
        Timeout: time.Duration(a.config.HTTPTimeout) * time.Second,
    }

    // Sync backup status to remote service
    return a.syncToRemote(client)
}

func (a *Application) syncToRemote(client *http.Client) error {
    // Prepare sync data
    syncData := map[string]interface{}{
        "timestamp": time.Now(),
        "files":     a.processedFiles,
        "git_info":  a.getGitInfo(),
    }

    // Format as JSON
    data, err := json.Marshal(syncData)
    if err != nil {
        return errors.Wrap(err, "failed to marshal sync data")
    }

    // Send to remote service
    resp, err := client.Post(a.config.SyncURL, "application/json", bytes.NewReader(data))
    if err != nil {
        return errors.Wrap(err, "failed to sync to remote")
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return errors.New("remote sync failed").
            WithField("status_code", resp.StatusCode)
    }

    return nil
}
```

## 🧪 Testing Integration

### Unit Testing with Packages

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Testing patterns
func TestApplicationIntegration(t *testing.T) {
    // Create test configuration
    cfg := &config.Config{
        OutputFormat: "json",
        MaxWorkers:   2,
    }

    // Create test application
    app := &Application{
        config:    cfg,
        formatter: formatter.New(cfg.OutputFormat),
        processor: processing.NewPool(cfg.MaxWorkers),
    }

    // Test file processing
    testFiles := []string{"test1.txt", "test2.txt"}
    err := app.Process(testFiles)
    
    assert.NoError(t, err)
    assert.Equal(t, len(testFiles), app.processedFiles)
}

// 🔺 EXTRACT-008: Integration documentation system - 📝 Integration testing patterns
func TestGitIntegration(t *testing.T) {
    // Create temporary git repository
    tempDir := t.TempDir()
    repo, err := git.Init(tempDir)
    require.NoError(t, err)

    // Create application with git
    app := &Application{
        git: repo,
    }

    // Test git-aware operations
    err = app.ProcessWithGitAwareness(context.Background(), []string{"test.txt"})
    assert.NoError(t, err)
}
```

## 🚀 Production Deployment

### Configuration Management

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Production configuration patterns
func (a *Application) LoadProductionConfig() error {
    // Load configuration with environment overrides
    cfg, err := config.LoadWithDefaults(map[string]interface{}{
        "max_workers":    runtime.NumCPU(),
        "output_format":  "json",
        "log_level":      "info",
        "timeout":        "30s",
    })
    if err != nil {
        return errors.Wrap(err, "failed to load production configuration")
    }

    // Validate production requirements
    if cfg.MaxWorkers < 1 || cfg.MaxWorkers > 100 {
        return errors.New("invalid worker count for production").
            WithField("workers", cfg.MaxWorkers)
    }

    a.config = cfg
    return nil
}
```

### Graceful Shutdown

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Graceful shutdown patterns
func (a *Application) RunWithGracefulShutdown(ctx context.Context) error {
    // Create cancellable context
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Handle shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        a.formatter.PrintInfo("Received shutdown signal, stopping gracefully...")
        cancel()
    }()

    // Run application with context
    return a.Run(ctx)
}

func (a *Application) Run(ctx context.Context) error {
    // Initialize resources
    if err := a.initialize(); err != nil {
        return errors.Wrap(err, "failed to initialize application")
    }
    defer a.cleanup()

    // Main application loop
    for {
        select {
        case <-ctx.Done():
            a.formatter.PrintInfo("Shutting down gracefully...")
            return nil
        default:
            if err := a.processNext(); err != nil {
                return errors.Wrap(err, "processing error")
            }
        }
    }
}
```

## 📊 Monitoring and Observability

### Metrics Collection

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Performance monitoring patterns
type Metrics struct {
    ProcessedFiles   int64
    ProcessingTime   time.Duration
    ErrorCount       int64
    GitOperations    int64
    ResourceUsage    map[string]int64
}

func (a *Application) CollectMetrics() *Metrics {
    return &Metrics{
        ProcessedFiles:   atomic.LoadInt64(&a.processedFiles),
        ProcessingTime:   a.totalProcessingTime,
        ErrorCount:       atomic.LoadInt64(&a.errorCount),
        GitOperations:    atomic.LoadInt64(&a.gitOperations),
        ResourceUsage:    a.resources.Usage(),
    }
}

func (a *Application) ReportMetrics() error {
    metrics := a.CollectMetrics()
    
    report := map[string]interface{}{
        "metrics":   metrics,
        "timestamp": time.Now(),
        "version":   a.version,
    }

    return a.formatter.Print(report)
}
```

## 🔄 Migration and Compatibility

### Legacy System Integration

```go
// 🔺 EXTRACT-008: Integration documentation system - 📝 Migration compatibility patterns
func (a *Application) MigrateLegacyConfig(legacyPath string) error {
    // Read legacy configuration
    legacyData, err := os.ReadFile(legacyPath)
    if err != nil {
        return errors.Wrap(err, "failed to read legacy config")
    }

    // Convert to new format
    newConfig, err := a.convertLegacyConfig(legacyData)
    if err != nil {
        return errors.Wrap(err, "failed to convert legacy config")
    }

    // Validate new configuration
    if err := config.Validate(newConfig); err != nil {
        return errors.Wrap(err, "converted config is invalid")
    }

    // Save new configuration
    return config.Save(newConfig, a.config.Path())
}

func (a *Application) convertLegacyConfig(data []byte) (*config.Config, error) {
    // Parse legacy format
    var legacy map[string]interface{}
    if err := json.Unmarshal(data, &legacy); err != nil {
        return nil, errors.Wrap(err, "failed to parse legacy config")
    }

    // Convert to new structure
    cfg := &config.Config{
        OutputFormat: getString(legacy, "format", "text"),
        MaxWorkers:   getInt(legacy, "threads", 1),
        Timeout:      getDuration(legacy, "timeout", "30s"),
    }

    return cfg, nil
}
```

## 📚 Summary

This integration guide demonstrates how to effectively use the extracted BkpDir packages:

1. **Individual Package Usage**: Each package can be used independently
2. **Integration Patterns**: Common patterns for combining packages
3. **Error Handling**: Structured error management across packages
4. **Performance**: Optimization techniques and resource management
5. **Production**: Deployment and monitoring considerations
6. **Testing**: Integration testing strategies
7. **Migration**: Legacy system compatibility

The packages are designed to be composable, allowing you to use only what you need while maintaining consistency and reliability across your CLI applications.

---

**🔺 EXTRACT-008: Integration documentation system - 📝 Complete integration guide with comprehensive examples** 