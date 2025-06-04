# 🔧 Troubleshooting: Common Issues and Solutions

> **🔺 EXTRACT-008: Tutorial series creation - 📚 Troubleshooting guide**

## 🎯 Overview

This guide covers common issues you might encounter when using the extracted packages and provides practical solutions. Each issue includes symptoms, root causes, and step-by-step resolution.

## 📋 Common Issues by Package

### pkg/config Issues

#### Issue: Configuration File Not Found

**Symptoms:**
```
Failed to load config: configuration file not found
```

**Root Cause:** Configuration discovery can't find a valid config file.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Config troubleshooting
func loadConfigWithFallback() (*config.Config, error) {
    // Try current directory first
    cfg, err := config.LoadConfig(".")
    if err != nil {
        // Try home directory
        homeDir, _ := os.UserHomeDir()
        if homeDir != "" {
            cfg, err = config.LoadConfig(homeDir)
            if err == nil {
                return cfg, nil
            }
        }
        
        // Use defaults as last resort
        fmt.Fprintf(os.Stderr, "Warning: Using default configuration\n")
        return config.NewDefault(), nil
    }
    return cfg, nil
}
```

#### Issue: Environment Variables Not Working

**Symptoms:**
```
Environment variable MYAPP_DEBUG=true not affecting config
```

**Root Cause:** Environment overrides not applied or wrong prefix.

**Solution:**
```go
// Apply environment overrides after loading
cfg, err := config.LoadConfig(".")
if err != nil {
    cfg = config.NewDefault()
}

// Apply with correct prefix (must match your app)
cfg.ApplyEnvironmentOverrides("MYAPP_")

// Verify the override worked
debug := cfg.GetBool("debug", false)
fmt.Printf("Debug mode: %v\n", debug)
```

### pkg/errors Issues

#### Issue: Error Context Lost

**Symptoms:**
```
Error: file operation failed
// No information about which file or why
```

**Root Cause:** Errors not properly wrapped with context.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Error context preservation
func processFileWithContext(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        // Wrap with context - don't just return err
        return errors.WrapError(err,
            fmt.Sprintf("Failed to open file %s", filename),
            "FILE_OPEN_ERROR")
    }
    defer file.Close()

    if err := processData(file); err != nil {
        // Add more context at each level
        return errors.WrapError(err,
            fmt.Sprintf("Failed to process data from %s", filename),
            "DATA_PROCESS_ERROR")
    }

    return nil
}
```

#### Issue: Wrong Error Type Used

**Symptoms:**
```
System error for user input validation
```

**Root Cause:** Using SystemError for user input issues.

**Solution:**
```go
func validateInput(input string) error {
    if input == "" {
        // User error - not system error
        return errors.NewUserError(
            "Input cannot be empty",
            "EMPTY_INPUT")
    }
    
    if len(input) > 100 {
        // User error - input too long
        return errors.NewUserError(
            "Input must be 100 characters or less",
            "INPUT_TOO_LONG")
    }
    
    // System error for file operations
    if err := validateInputFile(input); err != nil {
        return errors.NewSystemError(
            "Failed to validate input file",
            "FILE_VALIDATION_ERROR",
            err)
    }
    
    return nil
}
```

### pkg/resources Issues

#### Issue: Resources Not Cleaned Up

**Symptoms:**
```
File handles left open, temporary files not deleted
```

**Root Cause:** Resource manager not used or cleanup not called.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Resource cleanup patterns
func properResourceManagement() error {
    rm := resources.NewManager()
    defer rm.Cleanup() // Always defer cleanup

    // Register resources as they're created
    file, err := os.Open("data.txt")
    if err != nil {
        return err
    }
    rm.Add(func() error { return file.Close() })

    tmpDir, err := os.MkdirTemp("", "app-")
    if err != nil {
        return err
    }
    rm.Add(func() error { return os.RemoveAll(tmpDir) })

    // Even if this fails, cleanup still happens
    return doWork(file, tmpDir)
}
```

#### Issue: Cleanup Timeout

**Symptoms:**
```
Application hangs during shutdown
```

**Root Cause:** Cleanup functions taking too long.

**Solution:**
```go
func gracefulCleanupWithTimeout(rm *resources.Manager) error {
    // Use cleanup with timeout
    err := rm.CleanupWithTimeout(30 * time.Second)
    if err != nil {
        log.Printf("Cleanup timeout or error: %v", err)
        // Force exit if needed
        return err
    }
    return nil
}
```

### pkg/formatter Issues

#### Issue: Template Parsing Errors

**Symptoms:**
```
template: parse error at line 3: unexpected "}"
```

**Root Cause:** Invalid template syntax.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Template debugging
func debugTemplate() {
    // Test template syntax separately
    templateStr := `
Name: {{.Name}}
{{if .Active}}Status: Active{{else}}Status: Inactive{{end}}
Count: {{len .Items}}
`

    // Parse template to check for errors
    tmpl, err := template.New("test").Parse(templateStr)
    if err != nil {
        fmt.Printf("Template error: %v\n", err)
        return
    }

    // Test with sample data
    data := map[string]interface{}{
        "Name":   "Test",
        "Active": true,
        "Items":  []string{"a", "b", "c"},
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        fmt.Printf("Template execution error: %v\n", err)
        return
    }

    fmt.Printf("Template output:\n%s", buf.String())
}
```

#### Issue: Missing Template Functions

**Symptoms:**
```
template: function "formatBytes" not defined
```

**Root Cause:** Using undefined template functions.

**Solution:**
```go
// Check available template functions in formatter
formatter := formatter.New(cfg)

// Use built-in functions:
// - formatTime: Format time values
// - percent: Format as percentage  
// - bytes: Format byte counts
// - duration: Format durations
// - upper/lower/title: String case conversion

template := `
Size: {{.Size | bytes}}
Progress: {{.Progress | percent}}
Duration: {{.Duration | duration}}
Name: {{.Name | title}}
`
```

### pkg/git Issues

#### Issue: Git Commands Failing

**Symptoms:**
```
Failed to get Git info: git command not found
```

**Root Cause:** Git not installed or not in PATH.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Git troubleshooting
func checkGitAvailability() error {
    // Check if git is available
    if !git.DetectRepository(".") {
        return errors.NewUserError(
            "Not in a Git repository",
            "NOT_GIT_REPO")
    }

    // Try to get basic info to test git availability
    _, err := git.GetBranch(".")
    if err != nil {
        return errors.WrapError(err,
            "Git command failed - ensure Git is installed and in PATH",
            "GIT_COMMAND_ERROR")
    }

    return nil
}

// Use custom git command if needed
func useCustomGitCommand() {
    config := git.Config{
        WorkingDirectory:   ".",
        IncludeDirtyStatus: true,
        GitCommand:         "/usr/local/bin/git", // Custom path
    }

    repo := git.NewRepository(config)
    info, err := repo.GetInfo()
    if err != nil {
        log.Printf("Git error: %v", err)
        return
    }

    fmt.Printf("Branch: %s\n", info.Branch)
}
```

### pkg/cli Issues

#### Issue: Commands Not Found

**Symptoms:**
```
Unknown command: mycommand
```

**Root Cause:** Command not registered or wrong name.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 CLI command debugging
func debugCommandRegistration() {
    app := cli.NewApplication("myapp", "1.0.0")
    
    // Register commands with exact names
    app.AddCommand("create", "Create a resource", createCommand)
    app.AddCommand("list", "List resources", listCommand)
    
    // Debug: Print registered commands
    fmt.Println("Registered commands:")
    // Use help command to see what's registered
    ctx := context.Background()
    app.Execute(ctx, []string{"help"})
}
```

#### Issue: Context Not Available

**Symptoms:**
```
panic: CLI context not found in context
```

**Root Cause:** Trying to access CLI context outside command execution.

**Solution:**
```go
func commandWithContext(ctx context.Context, args []string) error {
    // Check if CLI context is available
    cliCtx := cli.FromContext(ctx)
    if cliCtx == nil {
        return errors.NewApplicationError(
            "CLI context not available",
            "NO_CLI_CONTEXT")
    }

    // Now safe to use
    if cliCtx.Verbose {
        fmt.Println("Verbose mode enabled")
    }

    return nil
}
```

### pkg/processing Issues

#### Issue: Worker Pool Deadlock

**Symptoms:**
```
Processing hangs indefinitely
```

**Root Cause:** Workers waiting for results that never come.

**Solution:**
```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Processing deadlock prevention
func preventProcessingDeadlock() error {
    // Always use context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Handle cancellation gracefully
    processor := processing.NewConcurrentProcessor(func(ctx context.Context, item *processing.ProcessingItem) (interface{}, error) {
        // Check for cancellation frequently
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        // Do work...
        time.Sleep(100 * time.Millisecond) // Simulate work

        // Check again before returning
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        return "processed", nil
    })

    items := []processing.ProcessingItem{
        {ID: "1", Data: "item1"},
        {ID: "2", Data: "item2"},
    }

    result, err := processor.Process(ctx, items)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return errors.NewApplicationError(
                "Processing timed out",
                "PROCESSING_TIMEOUT")
        }
        return err
    }

    fmt.Printf("Processed %d items\n", result.TotalItems)
    return nil
}
```

## 🔍 Debugging Techniques

### 1. Enable Debug Logging

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Debug logging patterns
func enableDebugLogging(cfg *config.Config) {
    if cfg.GetBool("debug", false) {
        log.SetLevel(log.DebugLevel)
        log.Debug("Debug logging enabled")
    }
}
```

### 2. Trace Error Chains

```go
func traceErrorChain(err error) {
    fmt.Printf("Error: %v\n", err)
    
    // Unwrap error chain
    current := err
    level := 0
    for current != nil {
        fmt.Printf("  Level %d: %v\n", level, current)
        
        if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
            current = unwrapper.Unwrap()
        } else {
            break
        }
        level++
    }
}
```

### 3. Monitor Resource Usage

```go
func monitorResources() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Memory: Alloc=%d KB, Sys=%d KB, NumGC=%d\n",
        m.Alloc/1024, m.Sys/1024, m.NumGC)
    
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}
```

## 🛠️ Development Tools

### Configuration Validation

```bash
# Validate YAML syntax
yamllint config.yaml

# Test configuration loading
go run -ldflags="-X main.debug=true" . config --verbose
```

### Template Testing

```bash
# Test templates separately
echo '{"Name":"Test","Count":5}' | go run template-test.go
```

### Git Integration Testing

```bash
# Test in different Git states
git checkout -b test-branch
echo "test" > test.txt
git add test.txt
# Test with staged changes

git commit -m "test"
# Test with clean state

echo "more" >> test.txt
# Test with dirty state
```

## 📊 Performance Troubleshooting

### Memory Leaks

```go
// 🔺 EXTRACT-008: Tutorial series creation - 📚 Memory leak detection
func detectMemoryLeaks() {
    // Force GC and get baseline
    runtime.GC()
    var m1 runtime.MemStats
    runtime.ReadMemStats(&m1)
    baseline := m1.Alloc

    // Run your code
    for i := 0; i < 1000; i++ {
        // Your processing code here
    }

    // Force GC and check memory
    runtime.GC()
    var m2 runtime.MemStats
    runtime.ReadMemStats(&m2)
    
    if m2.Alloc > baseline*2 {
        fmt.Printf("Potential memory leak: %d -> %d bytes\n", baseline, m2.Alloc)
    }
}
```

### Goroutine Leaks

```go
func detectGoroutineLeaks() {
    baseline := runtime.NumGoroutine()
    
    // Run your concurrent code
    // ...
    
    time.Sleep(1 * time.Second) // Allow cleanup
    current := runtime.NumGoroutine()
    
    if current > baseline+5 { // Allow some variance
        fmt.Printf("Potential goroutine leak: %d -> %d\n", baseline, current)
    }
}
```

## 🚨 Emergency Procedures

### Force Cleanup

```go
func forceCleanup() {
    // Set short timeout for cleanup
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Try graceful cleanup first
    rm := resources.NewManager()
    done := make(chan error, 1)
    go func() {
        done <- rm.Cleanup()
    }()

    select {
    case err := <-done:
        if err != nil {
            log.Printf("Cleanup error: %v", err)
        }
    case <-ctx.Done():
        log.Println("Cleanup timeout, forcing exit")
        os.Exit(1)
    }
}
```

### Recovery from Panics

```go
func recoverFromPanic() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            log.Printf("Stack trace:\n%s", debug.Stack())
            
            // Try to cleanup
            forceCleanup()
            os.Exit(1)
        }
    }()

    // Your code here
}
```

## 📚 Additional Resources

- [Integration Guide](../integration-guide.md) - Complete usage patterns
- [Package Reference](../package-reference.md) - Detailed API documentation
- [Getting Started](getting-started.md) - Basic usage tutorial
- [Advanced Patterns](advanced-patterns.md) - Complex integration scenarios

---

**🔺 EXTRACT-008: Tutorial series creation - 📚 Comprehensive troubleshooting guide** 