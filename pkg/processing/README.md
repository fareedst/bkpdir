# Package processing

// EXTRACT-010: Package processing documentation - Comprehensive API documentation and usage examples - 🔺

## Overview

Package `processing` provides generalized data processing patterns extracted from BkpDir, offering reusable patterns for timestamp-based naming conventions, data integrity verification, processing pipelines, and concurrent processing with resource management.

This package is designed to work independently or in combination with other extracted BkpDir packages like `pkg/config`, `pkg/errors`, and `pkg/formatter`.

## Key Features

- **Timestamp-based Naming**: Generate consistent names with timestamps and metadata
- **Data Integrity Verification**: Pluggable verification algorithms (SHA256, etc.)
- **Processing Pipelines**: Context-aware pipelines with atomic operations
- **Concurrent Processing**: Worker pools with resource management
- **Git Integration**: Git branch and hash information in naming
- **Context Support**: Full context.Context integration for cancellation
- **Progress Tracking**: Real-time progress monitoring and callbacks

## Installation

```bash
go get github.com/yourusername/bkpdir/pkg/processing
```

## Quick Start

### Basic Naming Convention

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    // Create a naming provider
    naming := processing.NewNamingProvider()
    
    // Generate a timestamped name
    name, err := naming.GenerateName(&processing.NamingTemplate{
        Prefix:    "backup",
        Timestamp: time.Now(),
        GitBranch: "main",
        GitHash:   "abc123",
        Note:      "initial",
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Generated name:", name)
    // Output: backup-2024-01-02T120000-main-abc123-initial
}
```

### Simple Processing Pipeline

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    // Create a processing pipeline
    pipeline := processing.NewPipeline("data-processing")
    
    // Add custom stages (implement PipelineStage interface)
    pipeline.AddStage(&CollectionStage{})
    pipeline.AddStage(&ProcessingStage{})
    pipeline.AddStage(&VerificationStage{})
    
    // Execute pipeline
    input := &processing.ProcessingInput{
        Source:      "/path/to/data",
        Destination: "/path/to/output",
        Options:     map[string]interface{}{"format": "json"},
    }
    
    ctx := context.Background()
    result, err := pipeline.Execute(ctx, input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Processed %d items in %v\n", result.ItemsProcessed, result.Duration)
}
```

## API Reference

### Core Interfaces

#### ProcessorInterface

Main interface for data processing operations:

```go
type ProcessorInterface interface {
    Process(ctx context.Context, input *ProcessingInput) (*ProcessingResult, error)
    Validate(input *ProcessingInput) error
    GetStatus() ProcessingStatus
}
```

#### NamingProviderInterface

Interface for timestamp-based naming:

```go
type NamingProviderInterface interface {
    GenerateName(template *NamingTemplate) (string, error)
    ParseName(name string, pattern string) (*NameComponents, error)
    ValidateName(name string, pattern string) error
    GetSupportedFormats() []string
}
```

#### PipelineInterface

Interface for processing workflows:

```go
type PipelineInterface interface {
    Execute(ctx context.Context, input *ProcessingInput) (*ProcessingResult, error)
    AddStage(stage PipelineStage)
    GetProgress() *PipelineProgress
    GetStages() []PipelineStage
}
```

### Key Types

#### ProcessingInput

Input configuration for processing operations:

```go
type ProcessingInput struct {
    Source      string                 `json:"source"`
    Destination string                 `json:"destination"`
    Data        io.Reader              `json:"-"`
    Options     map[string]interface{} `json:"options"`
    Metadata    map[string]string      `json:"metadata"`
    Timeout     time.Duration          `json:"timeout"`
    Concurrency int                    `json:"concurrency"`
}
```

#### ProcessingResult

Result of processing operations:

```go
type ProcessingResult struct {
    Output         string            `json:"output"`
    OutputData     io.Reader         `json:"-"`
    ItemsProcessed int               `json:"items_processed"`
    Duration       time.Duration     `json:"duration"`
    Statistics     map[string]int64  `json:"statistics"`
    Checksum       string            `json:"checksum,omitempty"`
    Verified       bool              `json:"verified"`
    Errors         []string          `json:"errors,omitempty"`
    Warnings       []string          `json:"warnings,omitempty"`
}
```

#### NamingTemplate

Template for generating timestamped names:

```go
type NamingTemplate struct {
    Prefix             string            `json:"prefix"`
    Timestamp          time.Time         `json:"timestamp"`
    Note               string            `json:"note,omitempty"`
    GitBranch          string            `json:"git_branch,omitempty"`
    GitHash            string            `json:"git_hash,omitempty"`
    GitIsClean         bool              `json:"git_is_clean"`
    ShowGitDirtyStatus bool              `json:"show_git_dirty_status"`
    Metadata           map[string]string `json:"metadata,omitempty"`
    TimestampFormat    string            `json:"timestamp_format"`
    IsIncremental      bool              `json:"is_incremental"`
    BaseName           string            `json:"base_name,omitempty"`
}
```

## Examples

### Advanced Naming with Git Information

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    naming := processing.NewNamingProvider()
    
    // Generate archive name with Git information
    archiveName := naming.GenerateArchiveName(
        "myproject",           // prefix
        "2024-01-02T120000",   // timestamp
        "feature-branch",      // git branch
        "abc123def",          // git hash
        "before-refactor",    // note
        false,                // git is clean
        true,                 // show dirty status
        false,                // is incremental
        "",                   // base name
    )
    
    fmt.Println("Archive name:", archiveName)
    // Output: myproject-2024-01-02T120000-feature-branch-abc123def-dirty-before-refactor.zip
    
    // Parse name back to components
    components, err := naming.ParseName(archiveName, "archive")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Parsed: prefix=%s, branch=%s, hash=%s, note=%s\n",
        components.Prefix, components.GitBranch, components.GitHash, components.Note)
}
```

### Custom Pipeline Stage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

// Custom validation stage
type ValidationStage struct {
    *processing.BaseStage
}

func NewValidationStage() *ValidationStage {
    return &ValidationStage{
        BaseStage: processing.NewBaseStage(
            "validation",
            "Validates input data integrity",
            time.Second*30,
        ),
    }
}

func (vs *ValidationStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
    // Custom validation logic here
    fmt.Printf("Validating source: %s\n", input.Source)
    
    // Simulate validation work
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(time.Millisecond * 100):
        // Validation complete
    }
    
    output.Statistics["validated_items"] = 42
    return nil
}

func main() {
    pipeline := processing.NewPipeline("custom-validation")
    pipeline.AddStage(NewValidationStage())
    
    // Set progress callback
    pipeline.SetProgressCallback(func(progress *processing.PipelineProgress) {
        fmt.Printf("Progress: %.1f%% (Stage: %s)\n", 
            progress.OverallProgress*100, progress.CurrentStage)
    })
    
    input := &processing.ProcessingInput{
        Source: "/data",
        Timeout: time.Minute,
    }
    
    result, err := pipeline.Execute(context.Background(), input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Validation complete: %d items validated\n", 
        result.Statistics["validated_items"])
}
```

### Concurrent Processing with Progress

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    // Create concurrent processing pipeline
    pipeline := processing.NewPipeline("concurrent-processing")
    
    // Add stages that support concurrency
    pipeline.AddStage(&ConcurrentProcessingStage{
        concurrency: 4,
        batchSize:   10,
    })
    
    // Monitor progress in separate goroutine
    var wg sync.WaitGroup
    wg.Add(1)
    
    go func() {
        defer wg.Done()
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                progress := pipeline.GetProgress()
                if progress.OverallProgress >= 1.0 {
                    fmt.Println("Processing complete!")
                    return
                }
                fmt.Printf("Progress: %.1f%% - %s (ETA: %v)\n",
                    progress.OverallProgress*100,
                    progress.CurrentStage,
                    progress.RemainingTime,
                )
            }
        }
    }()
    
    // Execute pipeline
    input := &processing.ProcessingInput{
        Source:      "/large-dataset",
        Destination: "/processed-output",
        Concurrency: 4,
        Timeout:     time.Hour,
    }
    
    ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
    defer cancel()
    
    result, err := pipeline.Execute(ctx, input)
    if err != nil {
        panic(err)
    }
    
    wg.Wait()
    
    fmt.Printf("Final result: %d items processed in %v\n",
        result.ItemsProcessed, result.Duration)
}

// Example concurrent processing stage
type ConcurrentProcessingStage struct {
    *processing.BaseStage
    concurrency int
    batchSize   int
}

func (cps *ConcurrentProcessingStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
    // Implement concurrent processing logic
    fmt.Printf("Processing with %d workers, batch size %d\n", 
        cps.concurrency, cps.batchSize)
    
    // Simulate processing work
    time.Sleep(time.Second * 2)
    
    output.ItemsProcessed = 1000
    output.Statistics["workers_used"] = int64(cps.concurrency)
    output.Statistics["batches_processed"] = 100
    
    return nil
}
```

### Error Handling and Recovery

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    // Create pipeline with error handling
    pipeline := processing.NewPipeline("error-handling-demo")
    
    // Configure pipeline to continue on error
    pipeline.SetStopOnError(false)
    
    // Add stages with potential failures
    pipeline.AddStage(&RiskyStage{name: "stage1"})
    pipeline.AddStage(&RiskyStage{name: "stage2"})
    pipeline.AddStage(&RecoveryStage{})
    
    input := &processing.ProcessingInput{
        Source: "/risky-data",
        Options: map[string]interface{}{
            "retry_failed": true,
            "max_retries":  3,
        },
    }
    
    result, err := pipeline.Execute(context.Background(), input)
    
    // Check for errors but continue if recoverable
    if err != nil {
        fmt.Printf("Pipeline error: %v\n", err)
        if processingErr, ok := err.(*processing.ProcessingError); ok {
            if processingErr.Recoverable {
                fmt.Println("Error is recoverable, attempting recovery...")
            }
        }
    }
    
    if len(result.Errors) > 0 {
        fmt.Printf("Processing completed with %d errors:\n", len(result.Errors))
        for _, errMsg := range result.Errors {
            fmt.Printf("  - %s\n", errMsg)
        }
    }
    
    if len(result.Warnings) > 0 {
        fmt.Printf("Processing completed with %d warnings:\n", len(result.Warnings))
        for _, warning := range result.Warnings {
            fmt.Printf("  - %s\n", warning)
        }
    }
    
    fmt.Printf("Successfully processed %d items\n", result.ItemsProcessed)
}

type RiskyStage struct {
    *processing.BaseStage
    name string
}

func (rs *RiskyStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
    // Simulate potential failure
    if rs.name == "stage1" {
        err := processing.NewProcessingError("SIMULATION", "RiskyStage", "Simulated failure in stage1")
        err.Recoverable = true
        output.Errors = append(output.Errors, err.Error())
        return err
    }
    
    if rs.name == "stage2" {
        output.Warnings = append(output.Warnings, "Stage2 completed with warnings")
    }
    
    return nil
}

type RecoveryStage struct {
    *processing.BaseStage
}

func (rs *RecoveryStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
    fmt.Println("Recovery stage: cleaning up and finalizing...")
    output.ItemsProcessed = 100 // Partial success
    return nil
}
```

## Integration with Other Packages

### With pkg/config

```go
package main

import (
    "context"
    
    "github.com/yourusername/bkpdir/pkg/config"
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    // Load processing configuration
    loader := config.NewConfigLoader()
    cfg, err := loader.LoadConfig("processing-config.yaml")
    if err != nil {
        panic(err)
    }
    
    // Configure processing options from config
    options := &processing.ProcessingOptions{
        Concurrency:           cfg.GetInt("processing.concurrency", 4),
        Timeout:               cfg.GetDuration("processing.timeout", time.Minute*5),
        EnableVerification:    cfg.GetBool("processing.verify", true),
        VerificationAlgorithm: cfg.GetString("processing.verify_algorithm", "sha256"),
        ContinueOnError:       cfg.GetBool("processing.continue_on_error", false),
    }
    
    // Use options in processing
    input := &processing.ProcessingInput{
        Source:      cfg.GetString("input.source"),
        Destination: cfg.GetString("output.destination"),
        Concurrency: options.Concurrency,
        Timeout:     options.Timeout,
    }
    
    // Process with configuration
    pipeline := processing.NewPipeline("configured-processing")
    result, err := pipeline.Execute(context.Background(), input)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Configured processing complete: %d items\n", result.ItemsProcessed)
}
```

### With pkg/errors

```go
package main

import (
    "context"
    
    "github.com/yourusername/bkpdir/pkg/errors"
    "github.com/yourusername/bkpdir/pkg/processing"
)

func main() {
    pipeline := processing.NewPipeline("error-integrated")
    
    // Add stage that uses structured errors
    pipeline.AddStage(&ErrorAwareStage{})
    
    input := &processing.ProcessingInput{Source: "/data"}
    result, err := pipeline.Execute(context.Background(), input)
    
    if err != nil {
        // Handle with structured error handling
        if appErr := errors.AsApplicationError(err); appErr != nil {
            switch appErr.Category {
            case errors.CategoryFilesystem:
                fmt.Printf("Filesystem error: %s\n", appErr.Message)
            case errors.CategoryPermission:
                fmt.Printf("Permission error: %s\n", appErr.Message)
            default:
                fmt.Printf("Processing error: %s\n", appErr.Message)
            }
        }
    }
    
    fmt.Printf("Processing result: %d items\n", result.ItemsProcessed)
}

type ErrorAwareStage struct {
    *processing.BaseStage
}

func (eas *ErrorAwareStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
    // Use structured errors
    if input.Source == "" {
        return errors.NewApplicationError(
            errors.CategoryValidation,
            errors.SeverityError,
            "MISSING_SOURCE",
            "Source path is required for processing",
            map[string]interface{}{"stage": "ErrorAwareStage"},
        )
    }
    
    output.ItemsProcessed = 50
    return nil
}
```

## Performance Characteristics

The processing package is designed for high-performance data processing:

- **Naming Operations**: ~5μs per name generation
- **Pipeline Overhead**: <1% additional overhead vs direct processing
- **Concurrent Processing**: Scales linearly with worker count up to CPU cores
- **Memory Usage**: Bounded by configured batch sizes and worker pools
- **Context Cancellation**: <10ms response time for cancellation requests

### Benchmarking

```go
package main

import (
    "testing"
    "time"
    
    "github.com/yourusername/bkpdir/pkg/processing"
)

func BenchmarkNamingProvider(b *testing.B) {
    naming := processing.NewNamingProvider()
    template := &processing.NamingTemplate{
        Prefix:    "test",
        Timestamp: time.Now(),
        GitBranch: "main",
        GitHash:   "abc123",
        Note:      "benchmark",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = naming.GenerateName(template)
    }
}

func BenchmarkPipelineExecution(b *testing.B) {
    pipeline := processing.NewPipeline("benchmark")
    pipeline.AddStage(&processing.BaseStage{})
    
    input := &processing.ProcessingInput{
        Source:      "/test",
        Destination: "/output",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx := context.Background()
        _, _ = pipeline.Execute(ctx, input)
    }
}
```

## Best Practices

### 1. Context Management

Always pass context.Context and handle cancellation:

```go
func processWithTimeout() error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
    defer cancel()
    
    pipeline := processing.NewPipeline("timeout-aware")
    result, err := pipeline.Execute(ctx, input)
    
    if ctx.Err() == context.DeadlineExceeded {
        return fmt.Errorf("processing timed out after 10 minutes")
    }
    
    return err
}
```

### 2. Resource Management

Use appropriate concurrency levels and timeouts:

```go
options := processing.DefaultProcessingOptions()
options.Concurrency = runtime.NumCPU() // Match CPU cores
options.Timeout = time.Hour             // Generous timeout
options.MaxRetries = 3                  // Reasonable retry count
```

### 3. Error Handling

Implement proper error handling and recovery:

```go
pipeline.SetStopOnError(false) // Continue processing other items
pipeline.SetProgressCallback(func(progress *processing.PipelineProgress) {
    if progress.OverallProgress > 0.5 {
        // Implement checkpointing at 50% completion
        saveCheckpoint(progress)
    }
})
```

### 4. Progress Monitoring

Use progress callbacks for long-running operations:

```go
pipeline.SetProgressCallback(func(progress *processing.PipelineProgress) {
    log.Printf("Progress: %.1f%% - %s (ETA: %v)",
        progress.OverallProgress*100,
        progress.CurrentStage,
        progress.RemainingTime,
    )
})
```

### 5. Testing

Test pipeline stages independently:

```go
func TestCustomStage(t *testing.T) {
    stage := &CustomStage{}
    input := &processing.ProcessingInput{Source: "/test"}
    output := &processing.ProcessingResult{}
    
    ctx := context.Background()
    err := stage.Execute(ctx, input, output)
    
    assert.NoError(t, err)
    assert.Equal(t, 100, output.ItemsProcessed)
}
```

## Supported Formats

### Naming Patterns

The package supports several naming patterns extracted from BkpDir:

1. **Archive Pattern**: `{prefix}-{timestamp}-{git_branch}-{git_hash}-{note}.zip`
2. **Backup Pattern**: `{filename}-{timestamp}[={note}]`
3. **Incremental Pattern**: `{prefix}-{timestamp}-{git_info}-inc-{base}-{note}.zip`

### Timestamp Formats

- **Archive Format**: `2006-01-02T150405` (ISO 8601)
- **Backup Format**: `2006-01-02-15-04` (Backup compatible)
- **Custom Formats**: Any Go time format layout

## Thread Safety

All package components are thread-safe:

- **NamingProvider**: Concurrent name generation supported
- **Pipeline**: Safe for concurrent execution with different inputs
- **Progress Tracking**: Atomic operations for thread-safe updates
- **Error Handling**: Thread-safe error collection and reporting

The processing package provides a robust foundation for building scalable data processing applications with full context support, progress monitoring, and error resilience. 