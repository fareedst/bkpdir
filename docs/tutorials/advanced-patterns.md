# 🔧 Advanced Integration Patterns

> **🔺 EXTRACT-008: Tutorial series creation - 📚 Advanced integration scenarios**

## 🎯 Overview

This tutorial covers advanced patterns for integrating the extracted packages in complex scenarios. You'll learn about concurrent processing, resource management, Git integration, and building production-ready CLI applications.

## 📋 Prerequisites

- Completed the [Getting Started](getting-started.md) tutorial
- Understanding of Go concurrency patterns
- Familiarity with Git workflows

## 🏗️ Advanced Application Architecture

### Multi-Package Integration

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Advanced application structure
package main

import (
    "context"
    "fmt"
    "os"
    "runtime"
    "time"

    "./pkg/cli"
    "./pkg/config"
    "./pkg/errors"
    "./pkg/formatter"
    "./pkg/git"
    "./pkg/processing"
    "./pkg/resources"
    "./pkg/fileops"
)

type Application struct {
    config    *config.Config
    formatter *formatter.Formatter
    git       *git.Repository
    processor *processing.Pool
    resources *resources.Manager
}

func NewApplication() (*Application, error) {
    // Load configuration with environment overrides
    cfg, err := config.LoadWithDefaults(map[string]interface{}{
        "max_workers": runtime.NumCPU(),
        "output_format": "json",
        "timeout": "30s",
    })
    if err != nil {
        return nil, errors.Wrap(err, "failed to load configuration")
    }

    app := &Application{
        config:    cfg,
        formatter: formatter.New(cfg.OutputFormat),
        processor: processing.NewPool(cfg.MaxWorkers),
        resources: resources.NewManager(),
    }

    // Initialize Git if in repository
    if repo, err := git.Discover("."); err == nil {
        app.git = repo
    }

    return app, nil
}

func (a *Application) Close() error {
    var errs []error
    
    if err := a.processor.Close(); err != nil {
        errs = append(errs, err)
    }
    
    if err := a.resources.Cleanup(); err != nil {
        errs = append(errs, err)
    }
    
    if len(errs) > 0 {
        return errors.New("cleanup failed").WithField("errors", errs)
    }
    
    return nil
}
```

## 🔄 Concurrent Processing Patterns

### Pattern 1: Pipeline Processing

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Pipeline processing pattern
func (a *Application) ProcessPipeline(ctx context.Context, files []string) error {
    // Stage 1: File validation
    validFiles := make(chan string, len(files))
    go func() {
        defer close(validFiles)
        for _, file := range files {
            if fileops.IsAccessible(file) {
                select {
                case validFiles <- file:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()

    // Stage 2: Processing
    processedFiles := make(chan ProcessResult, 100)
    go func() {
        defer close(processedFiles)
        err := a.processor.Process(ctx, channelToSlice(validFiles), func(file string) error {
            result, err := a.processFile(file)
            if err != nil {
                return err
            }
            
            select {
            case processedFiles <- result:
            case <-ctx.Done():
            }
            return nil
        })
        
        if err != nil {
            a.formatter.PrintError("Processing failed: " + err.Error())
        }
    }()

    // Stage 3: Output formatting
    return a.formatResults(ctx, processedFiles)
}

type ProcessResult struct {
    File     string        `json:"file"`
    Size     int64         `json:"size"`
    Duration time.Duration `json:"duration"`
    GitInfo  *git.Info     `json:"git_info,omitempty"`
}

func (a *Application) processFile(file string) (ProcessResult, error) {
    start := time.Now()
    
    info, err := fileops.Stat(file)
    if err != nil {
        return ProcessResult{}, errors.Wrap(err, "failed to stat file").
            WithField("file", file)
    }

    result := ProcessResult{
        File:     file,
        Size:     info.Size,
        Duration: time.Since(start),
    }

    // Add Git information if available
    if a.git != nil {
        if gitInfo, err := a.git.Info(); err == nil {
            result.GitInfo = gitInfo
        }
    }

    return result, nil
}
```

### Pattern 2: Worker Pool with Error Aggregation

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Worker pool with error handling
func (a *Application) ProcessWithErrorAggregation(ctx context.Context, files []string) error {
    errorCollector := errors.NewCollector()
    results := make([]ProcessResult, 0, len(files))
    resultsMutex := sync.Mutex{}

    err := a.processor.ProcessWithCallback(
        ctx,
        files,
        func(file string) error {
            result, err := a.processFile(file)
            if err != nil {
                return err
            }

            resultsMutex.Lock()
            results = append(results, result)
            resultsMutex.Unlock()

            return nil
        },
        func(file string, err error) {
            if err != nil {
                errorCollector.Add(errors.Wrap(err, "processing failed").
                    WithField("file", file).
                    WithField("timestamp", time.Now()))
            }
        },
    )

    // Handle processing errors
    if err != nil {
        return errors.Wrap(err, "worker pool failed")
    }

    // Report results even if some files failed
    a.reportResults(results)

    // Return aggregated errors if any
    if errorCollector.HasErrors() {
        return errorCollector.Error()
    }

    return nil
}

func (a *Application) reportResults(results []ProcessResult) {
    summary := map[string]interface{}{
        "total_files":    len(results),
        "total_size":     a.calculateTotalSize(results),
        "average_time":   a.calculateAverageTime(results),
        "git_repository": a.git != nil,
    }

    if a.git != nil {
        if info, err := a.git.Info(); err == nil {
            summary["git_branch"] = info.Branch
            summary["git_commit"] = info.Commit[:8]
        }
    }

    a.formatter.Print(summary)
}
```

## 🛡️ Advanced Resource Management

### Pattern 1: Resource Pools

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Resource pool management
type ResourcePool struct {
    manager   *resources.Manager
    pools     map[string]chan resources.Resource
    limits    map[string]int
    mutex     sync.RWMutex
}

func NewResourcePool() *ResourcePool {
    return &ResourcePool{
        manager: resources.NewManager(),
        pools:   make(map[string]chan resources.Resource),
        limits:  map[string]int{
            "file_handle":    100,
            "network_conn":   50,
            "memory_buffer":  20,
        },
    }
}

func (rp *ResourcePool) AcquireFromPool(resourceType string) (resources.Resource, error) {
    rp.mutex.RLock()
    pool, exists := rp.pools[resourceType]
    rp.mutex.RUnlock()

    if !exists {
        return rp.createPool(resourceType)
    }

    select {
    case resource := <-pool:
        return resource, nil
    default:
        // Pool empty, create new resource
        return rp.manager.Acquire(resourceType)
    }
}

func (rp *ResourcePool) ReturnToPool(resource resources.Resource) {
    resourceType := resource.Type()
    
    rp.mutex.RLock()
    pool, exists := rp.pools[resourceType]
    rp.mutex.RUnlock()

    if !exists {
        resource.Release()
        return
    }

    select {
    case pool <- resource:
        // Successfully returned to pool
    default:
        // Pool full, release resource
        resource.Release()
    }
}
```

### Pattern 2: Context-Aware Resource Management

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Context-aware resources
func (a *Application) ProcessWithContextResources(ctx context.Context, files []string) error {
    // Create context-specific resource manager
    ctxManager := resources.NewManagerWithContext(ctx)
    defer ctxManager.Cleanup()

    // Acquire resources with context cancellation
    procResource, err := ctxManager.AcquireWithTimeout("processing", 30*time.Second)
    if err != nil {
        return errors.Wrap(err, "failed to acquire processing resource")
    }

    fsResource, err := ctxManager.AcquireWithTimeout("filesystem", 10*time.Second)
    if err != nil {
        return errors.Wrap(err, "failed to acquire filesystem resource")
    }

    // Create processor with resource constraints
    processor := processing.NewPoolWithResources(
        a.config.MaxWorkers,
        []resources.Resource{procResource, fsResource},
    )

    // Process with automatic cleanup on context cancellation
    return processor.ProcessWithContext(ctx, files, a.processFile)
}
```

## 🔌 Git-Aware Processing

### Pattern 1: Git-Driven File Selection

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Git-aware file processing
func (a *Application) ProcessGitAwareFiles(ctx context.Context, patterns []string) error {
    if a.git == nil {
        return errors.New("not in a git repository")
    }

    // Get files based on Git status
    files, err := a.selectFilesByGitStatus(patterns)
    if err != nil {
        return errors.Wrap(err, "failed to select files by git status")
    }

    // Process with Git context
    return a.processWithGitContext(ctx, files)
}

func (a *Application) selectFilesByGitStatus(patterns []string) ([]string, error) {
    var selectedFiles []string

    for _, pattern := range patterns {
        matches, err := filepath.Glob(pattern)
        if err != nil {
            return nil, errors.Wrap(err, "invalid pattern").
                WithField("pattern", pattern)
        }

        for _, file := range matches {
            status, err := a.git.FileStatus(file)
            if err != nil {
                continue // Skip files not in git
            }

            // Only process modified or new files
            if status.IsModified || status.IsUntracked {
                selectedFiles = append(selectedFiles, file)
            }
        }
    }

    return selectedFiles, nil
}

func (a *Application) processWithGitContext(ctx context.Context, files []string) error {
    gitInfo, err := a.git.Info()
    if err != nil {
        return errors.Wrap(err, "failed to get git info")
    }

    // Add Git context to all operations
    ctxWithGit := context.WithValue(ctx, "git_info", gitInfo)

    return a.processor.Process(ctxWithGit, files, func(file string) error {
        return a.processFileWithGitContext(ctxWithGit, file)
    })
}

func (a *Application) processFileWithGitContext(ctx context.Context, file string) error {
    gitInfo := ctx.Value("git_info").(*git.Info)
    
    // Include Git information in processing
    result := map[string]interface{}{
        "file":       file,
        "git_branch": gitInfo.Branch,
        "git_commit": gitInfo.Commit[:8],
        "timestamp":  time.Now(),
    }

    // Check if file has uncommitted changes
    if status, err := a.git.FileStatus(file); err == nil {
        result["git_status"] = status.Status
        
        if status.IsModified {
            a.formatter.PrintWarning(fmt.Sprintf(
                "File %s has uncommitted changes", file))
        }
    }

    return a.formatter.Print(result)
}
```

### Pattern 2: Branch-Aware Operations

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Branch-aware processing
func (a *Application) ProcessByBranch(ctx context.Context, files []string) error {
    if a.git == nil {
        return a.processRegular(ctx, files)
    }

    branch, err := a.git.Branch()
    if err != nil {
        return errors.Wrap(err, "failed to get current branch")
    }

    switch branch {
    case "main", "master":
        return a.processProduction(ctx, files)
    case "develop":
        return a.processDevelopment(ctx, files)
    default:
        return a.processFeatureBranch(ctx, files, branch)
    }
}

func (a *Application) processProduction(ctx context.Context, files []string) error {
    // Production processing with extra validation
    a.formatter.PrintInfo("Processing in PRODUCTION mode")
    
    // Ensure repository is clean
    isClean, err := a.git.IsClean()
    if err != nil {
        return errors.Wrap(err, "failed to check git status")
    }
    
    if !isClean {
        return errors.New("repository has uncommitted changes").
            WithField("branch", "production")
    }

    // Use conservative processing settings
    conservativeProcessor := processing.NewPool(1) // Single worker for safety
    defer conservativeProcessor.Close()

    return conservativeProcessor.Process(ctx, files, a.processFileConservatively)
}

func (a *Application) processDevelopment(ctx context.Context, files []string) error {
    // Development processing with enhanced logging
    a.formatter.PrintInfo("Processing in DEVELOPMENT mode")
    
    return a.processor.ProcessWithCallback(
        ctx,
        files,
        a.processFile,
        func(file string, err error) {
            if err != nil {
                a.formatter.PrintError(fmt.Sprintf("DEV: Failed %s: %v", file, err))
            } else {
                a.formatter.PrintInfo(fmt.Sprintf("DEV: Processed %s", file))
            }
        },
    )
}

func (a *Application) processFeatureBranch(ctx context.Context, files []string, branch string) error {
    // Feature branch processing with branch-specific settings
    a.formatter.PrintInfo(fmt.Sprintf("Processing in FEATURE mode: %s", branch))
    
    // Create branch-specific output directory
    outputDir := fmt.Sprintf("output/%s", branch)
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return errors.Wrap(err, "failed to create branch output directory").
            WithField("branch", branch)
    }

    return a.processor.Process(ctx, files, func(file string) error {
        return a.processFileForBranch(file, branch, outputDir)
    })
}
```

## 🔧 Configuration-Driven Behavior

### Pattern 1: Dynamic Configuration

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Dynamic configuration
type DynamicConfig struct {
    base     *config.Config
    overrides map[string]interface{}
    mutex    sync.RWMutex
}

func NewDynamicConfig(base *config.Config) *DynamicConfig {
    return &DynamicConfig{
        base:      base,
        overrides: make(map[string]interface{}),
    }
}

func (dc *DynamicConfig) SetOverride(key string, value interface{}) {
    dc.mutex.Lock()
    defer dc.mutex.Unlock()
    dc.overrides[key] = value
}

func (dc *DynamicConfig) GetWithOverride(key string) interface{} {
    dc.mutex.RLock()
    defer dc.mutex.RUnlock()
    
    if value, exists := dc.overrides[key]; exists {
        return value
    }
    
    return dc.base.Get(key)
}

func (a *Application) ProcessWithDynamicConfig(ctx context.Context, files []string) error {
    dynConfig := NewDynamicConfig(a.config)
    
    // Adjust configuration based on file count
    if len(files) > 1000 {
        dynConfig.SetOverride("max_workers", runtime.NumCPU()*2)
        dynConfig.SetOverride("batch_size", 200)
    } else if len(files) < 10 {
        dynConfig.SetOverride("max_workers", 1)
        dynConfig.SetOverride("batch_size", 1)
    }

    // Create processor with dynamic configuration
    workers := dynConfig.GetWithOverride("max_workers").(int)
    processor := processing.NewPool(workers)
    defer processor.Close()

    return processor.Process(ctx, files, a.processFile)
}
```

### Pattern 2: Environment-Aware Processing

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Environment-aware processing
func (a *Application) ProcessByEnvironment(ctx context.Context, files []string) error {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    switch env {
    case "production":
        return a.processForProduction(ctx, files)
    case "staging":
        return a.processForStaging(ctx, files)
    case "development":
        return a.processForDevelopment(ctx, files)
    default:
        return errors.New("unknown environment").
            WithField("environment", env)
    }
}

func (a *Application) processForProduction(ctx context.Context, files []string) error {
    // Production: Maximum safety, logging, monitoring
    prodConfig := &ProcessingConfig{
        Workers:     runtime.NumCPU(),
        BatchSize:   50,
        Timeout:     5 * time.Minute,
        RetryCount:  3,
        Monitoring:  true,
        Validation:  true,
    }

    return a.processWithConfig(ctx, files, prodConfig)
}

func (a *Application) processForStaging(ctx context.Context, files []string) error {
    // Staging: Production-like with more verbose logging
    stagingConfig := &ProcessingConfig{
        Workers:     runtime.NumCPU(),
        BatchSize:   100,
        Timeout:     3 * time.Minute,
        RetryCount:  2,
        Monitoring:  true,
        Validation:  true,
        VerboseLog:  true,
    }

    return a.processWithConfig(ctx, files, stagingConfig)
}

func (a *Application) processForDevelopment(ctx context.Context, files []string) error {
    // Development: Fast processing with detailed debugging
    devConfig := &ProcessingConfig{
        Workers:     2,
        BatchSize:   10,
        Timeout:     1 * time.Minute,
        RetryCount:  1,
        Monitoring:  false,
        Validation:  false,
        VerboseLog:  true,
        DebugMode:   true,
    }

    return a.processWithConfig(ctx, files, devConfig)
}

type ProcessingConfig struct {
    Workers     int
    BatchSize   int
    Timeout     time.Duration
    RetryCount  int
    Monitoring  bool
    Validation  bool
    VerboseLog  bool
    DebugMode   bool
}

func (a *Application) processWithConfig(ctx context.Context, files []string, cfg *ProcessingConfig) error {
    // Create processor with configuration
    processor := processing.NewPoolWithConfig(cfg.Workers, &processing.Config{
        BatchSize:  cfg.BatchSize,
        Timeout:    cfg.Timeout,
        RetryCount: cfg.RetryCount,
    })
    defer processor.Close()

    // Setup monitoring if enabled
    if cfg.Monitoring {
        go a.monitorProcessing(processor)
    }

    // Process with configuration-specific behavior
    return processor.Process(ctx, files, func(file string) error {
        if cfg.VerboseLog {
            a.formatter.PrintInfo(fmt.Sprintf("Processing: %s", file))
        }

        if cfg.Validation {
            if err := a.validateFile(file); err != nil {
                return errors.Wrap(err, "validation failed").
                    WithField("file", file)
            }
        }

        result, err := a.processFile(file)
        if err != nil && cfg.DebugMode {
            a.formatter.PrintError(fmt.Sprintf("Debug: %s failed: %v", file, err))
        }

        return err
    })
}
```

## 🚀 Production Patterns

### Pattern 1: Graceful Shutdown

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Graceful shutdown pattern
func (a *Application) RunWithGracefulShutdown() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Setup graceful shutdown
    shutdownChan := make(chan struct{})
    go func() {
        <-sigChan
        a.formatter.PrintInfo("Received shutdown signal, stopping gracefully...")
        
        // Cancel context to stop processing
        cancel()
        
        // Wait for cleanup with timeout
        cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cleanupCancel()
        
        if err := a.gracefulCleanup(cleanupCtx); err != nil {
            a.formatter.PrintError("Cleanup failed: " + err.Error())
        }
        
        close(shutdownChan)
    }()

    // Run main application
    err := a.Run(ctx)
    
    // Wait for shutdown completion
    select {
    case <-shutdownChan:
        a.formatter.PrintInfo("Shutdown completed")
    case <-time.After(35 * time.Second):
        a.formatter.PrintError("Shutdown timeout, forcing exit")
        return errors.New("shutdown timeout")
    }

    return err
}

func (a *Application) gracefulCleanup(ctx context.Context) error {
    var errs []error

    // Stop processor
    if err := a.processor.StopGracefully(ctx); err != nil {
        errs = append(errs, errors.Wrap(err, "processor shutdown failed"))
    }

    // Cleanup resources
    if err := a.resources.CleanupWithContext(ctx); err != nil {
        errs = append(errs, errors.Wrap(err, "resource cleanup failed"))
    }

    // Save state if needed
    if err := a.saveState(); err != nil {
        errs = append(errs, errors.Wrap(err, "state save failed"))
    }

    if len(errs) > 0 {
        return errors.New("cleanup partially failed").
            WithField("error_count", len(errs)).
            WithField("errors", errs)
    }

    return nil
}
```

### Pattern 2: Health Monitoring

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Health monitoring pattern
type HealthMonitor struct {
    app       *Application
    metrics   *Metrics
    alerts    chan Alert
    interval  time.Duration
}

type Metrics struct {
    ProcessedFiles   int64         `json:"processed_files"`
    FailedFiles      int64         `json:"failed_files"`
    AverageTime      time.Duration `json:"average_time"`
    MemoryUsage      int64         `json:"memory_usage"`
    GoroutineCount   int           `json:"goroutine_count"`
    LastUpdate       time.Time     `json:"last_update"`
}

type Alert struct {
    Level     string    `json:"level"`
    Message   string    `json:"message"`
    Timestamp time.Time `json:"timestamp"`
    Metrics   *Metrics  `json:"metrics"`
}

func NewHealthMonitor(app *Application) *HealthMonitor {
    return &HealthMonitor{
        app:      app,
        metrics:  &Metrics{},
        alerts:   make(chan Alert, 100),
        interval: 30 * time.Second,
    }
}

func (hm *HealthMonitor) Start(ctx context.Context) {
    ticker := time.NewTicker(hm.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            hm.collectMetrics()
            hm.checkHealth()
        }
    }
}

func (hm *HealthMonitor) collectMetrics() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    hm.metrics.MemoryUsage = int64(m.Alloc)
    hm.metrics.GoroutineCount = runtime.NumGoroutine()
    hm.metrics.LastUpdate = time.Now()

    // Collect application-specific metrics
    if appMetrics := hm.app.GetMetrics(); appMetrics != nil {
        hm.metrics.ProcessedFiles = appMetrics.ProcessedFiles
        hm.metrics.FailedFiles = appMetrics.FailedFiles
        hm.metrics.AverageTime = appMetrics.AverageTime
    }
}

func (hm *HealthMonitor) checkHealth() {
    // Check memory usage
    if hm.metrics.MemoryUsage > 1024*1024*1024 { // 1GB
        hm.sendAlert("warning", "High memory usage detected")
    }

    // Check goroutine count
    if hm.metrics.GoroutineCount > 1000 {
        hm.sendAlert("warning", "High goroutine count detected")
    }

    // Check failure rate
    if hm.metrics.FailedFiles > 0 {
        total := hm.metrics.ProcessedFiles + hm.metrics.FailedFiles
        failureRate := float64(hm.metrics.FailedFiles) / float64(total)
        
        if failureRate > 0.1 { // 10% failure rate
            hm.sendAlert("error", fmt.Sprintf("High failure rate: %.2f%%", failureRate*100))
        }
    }
}

func (hm *HealthMonitor) sendAlert(level, message string) {
    alert := Alert{
        Level:     level,
        Message:   message,
        Timestamp: time.Now(),
        Metrics:   hm.metrics,
    }

    select {
    case hm.alerts <- alert:
    default:
        // Alert channel full, log to stderr
        fmt.Fprintf(os.Stderr, "ALERT: %s - %s\n", level, message)
    }
}
```

## 📚 Summary

This advanced tutorial covered:

1. **Multi-Package Integration**: Building complex applications with all packages
2. **Concurrent Processing**: Pipeline patterns and worker pools with error handling
3. **Resource Management**: Advanced resource pools and context-aware management
4. **Git Integration**: Git-aware processing and branch-specific behavior
5. **Configuration Patterns**: Dynamic and environment-aware configuration
6. **Production Patterns**: Graceful shutdown and health monitoring

These patterns provide a foundation for building robust, production-ready CLI applications using the extracted packages.

---

**🔺 EXTRACT-008: Tutorial series creation - 📚 Advanced integration patterns and production-ready applications**
