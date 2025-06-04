# Package errors

[![Go Reference](https://pkg.go.dev/badge/github.com/bkpdir/pkg/errors.svg)](https://pkg.go.dev/github.com/bkpdir/pkg/errors)

## Overview

Package `errors` provides structured error handling and error classification utilities for CLI applications. It includes generic error types, error classification, and integration patterns for reliable error handling across applications, extracted and generalized from the BkpDir application.

### Key Features

- **Structured Errors**: ApplicationError type with status codes, operation context, and path information
- **Error Classification**: Categorize errors by type (filesystem, permission, disk space, etc.)
- **Error Context**: Rich context information for error handling operations
- **Severity Levels**: Error severity classification from info to fatal
- **Error Chaining**: Full support for error wrapping and unwrapping
- **CLI Integration**: Designed specifically for command-line application error handling

### Design Philosophy

The `errors` package was extracted to provide a standardized approach to error handling in CLI applications. It emphasizes structured error information, proper error classification, and context preservation to enable better error reporting and debugging.

## Installation

```bash
go get github.com/bkpdir/pkg/errors
```

## Quick Start

### Basic Structured Errors

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/bkpdir/pkg/errors"
)

func main() {
    // Create a structured error
    err := errors.NewApplicationError("Failed to create archive", 1)
    
    // Create error with underlying cause
    underlyingErr := fmt.Errorf("disk full")
    err = errors.NewApplicationErrorWithCause("Archive creation failed", 2, underlyingErr)
    
    // Create error with full context
    err = errors.NewApplicationErrorWithContext(
        "Archive creation failed",
        2,                    // status code
        "create",            // operation
        "/path/to/archive",  // path
        underlyingErr,       // underlying error
    )
    
    // Use the error
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
        fmt.Printf("Status Code: %d\n", err.GetStatusCode())
        fmt.Printf("Operation: %s\n", err.GetOperation())
        fmt.Printf("Path: %s\n", err.GetPath())
        os.Exit(err.GetStatusCode())
    }
}
```

### Error Classification

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/bkpdir/pkg/errors"
)

func main() {
    // Create error classifier
    classifier := errors.NewDefaultErrorClassifier()
    
    // Classify different types of errors
    diskErr := fmt.Errorf("no space left on device")
    permErr := fmt.Errorf("permission denied")
    fileErr := fmt.Errorf("file not found")
    
    // Classify errors
    diskCategory := classifier.ClassifyError(diskErr)
    permCategory := classifier.ClassifyError(permErr)
    fileCategory := classifier.ClassifyError(fileErr)
    
    fmt.Printf("Disk error category: %v\n", diskCategory)     // ErrorCategoryDiskSpace
    fmt.Printf("Permission error category: %v\n", permCategory) // ErrorCategoryPermission
    fmt.Printf("File error category: %v\n", fileCategory)     // ErrorCategoryFilesystem
    
    // Check if errors are recoverable
    fmt.Printf("Disk error recoverable: %t\n", classifier.IsRecoverable(diskErr))
    fmt.Printf("Permission error recoverable: %t\n", classifier.IsRecoverable(permErr))
    
    // Get error severity
    diskSeverity := classifier.GetSeverity(diskErr)
    fmt.Printf("Disk error severity: %v\n", diskSeverity) // ErrorSeverityCritical
}
```

## API Reference

### Core Types

#### ApplicationError

Structured error with status code and operation context:

```go
type ApplicationError struct {
    Message    string // Human-readable error message
    StatusCode int    // Exit status code for the application
    Operation  string // Operation context (e.g., "create", "verify", "list")
    Path       string // File or directory path involved in the error
    Err        error  // Underlying error for debugging and error chaining
}
```

**Methods:**
- `Error() string` - Returns formatted error message
- `GetStatusCode() int` - Returns exit status code
- `GetOperation() string` - Returns operation context
- `GetPath() string` - Returns path information
- `GetMessage() string` - Returns error message
- `Unwrap() error` - Returns underlying error

#### ErrorContext

Provides context information for error handling operations:

```go
type ErrorContext struct {
    Operation string                 // The operation being performed
    Path      string                 // File or directory path involved
    Context   context.Context        // Cancellation context
    Metadata  map[string]interface{} // Additional context metadata
    StartTime int64                  // Operation start time
}
```

### Core Interfaces

#### ErrorInterface

Common interface for all structured errors:

```go
type ErrorInterface interface {
    Error() string
    GetStatusCode() int
    GetOperation() string
    GetPath() string
    GetMessage() string
    Unwrap() error
}
```

#### ErrorClassifier

Interface for classifying different types of errors:

```go
type ErrorClassifier interface {
    ClassifyError(err error) ErrorCategory
    IsRecoverable(err error) bool
    GetSeverity(err error) ErrorSeverity
}
```

### Error Categories

```go
const (
    ErrorCategoryUnknown       ErrorCategory = iota
    ErrorCategoryFilesystem    // File system related errors
    ErrorCategoryPermission    // Permission and access errors
    ErrorCategoryDiskSpace     // Disk space and storage errors
    ErrorCategoryNetwork       // Network and connectivity errors
    ErrorCategoryConfiguration // Configuration and setup errors
    ErrorCategoryValidation    // Input validation errors
)
```

### Error Severity Levels

```go
const (
    ErrorSeverityInfo     ErrorSeverity = iota
    ErrorSeverityWarning  // Warning conditions
    ErrorSeverityError    // Error conditions
    ErrorSeverityCritical // Critical system errors
    ErrorSeverityFatal    // Fatal application errors
)
```

## Examples

### Advanced Error Handling

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    
    "github.com/bkpdir/pkg/errors"
)

func performOperation(ctx context.Context, path string) error {
    // Create error context
    errorCtx := errors.NewErrorContext("backup", path, ctx)
    errorCtx.WithMetadata("start_time", time.Now())
    errorCtx.WithMetadata("user", os.Getenv("USER"))
    
    // Simulate operation that might fail
    if _, err := os.Stat(path); err != nil {
        // Create structured error with context
        return errors.NewApplicationErrorWithContext(
            "Failed to access file for backup",
            3,                    // status code
            errorCtx.Operation,   // operation
            errorCtx.Path,        // path
            err,                  // underlying error
        )
    }
    
    // Check for cancellation
    select {
    case <-ctx.Done():
        return errors.NewApplicationErrorWithContext(
            "Operation cancelled",
            130, // SIGINT exit code
            errorCtx.Operation,
            errorCtx.Path,
            ctx.Err(),
        )
    default:
    }
    
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err := performOperation(ctx, "/nonexistent/file.txt")
    if err != nil {
        // Handle structured error
        if appErr, ok := err.(*errors.ApplicationError); ok {
            fmt.Printf("Operation '%s' failed\n", appErr.GetOperation())
            fmt.Printf("Path: %s\n", appErr.GetPath())
            fmt.Printf("Error: %s\n", appErr.GetMessage())
            
            // Classify the error
            classifier := errors.NewDefaultErrorClassifier()
            category := classifier.ClassifyError(appErr)
            severity := classifier.GetSeverity(appErr)
            
            fmt.Printf("Category: %v\n", category)
            fmt.Printf("Severity: %v\n", severity)
            fmt.Printf("Recoverable: %t\n", classifier.IsRecoverable(appErr))
            
            os.Exit(appErr.GetStatusCode())
        }
        
        fmt.Printf("Unexpected error: %v\n", err)
        os.Exit(1)
    }
}
```

### Error Handler with Recovery

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/bkpdir/pkg/errors"
)

type ErrorHandler struct {
    classifier errors.ErrorClassifier
    maxRetries int
    retryDelay time.Duration
}

func NewErrorHandler() *ErrorHandler {
    return &ErrorHandler{
        classifier: errors.NewDefaultErrorClassifier(),
        maxRetries: 3,
        retryDelay: time.Second,
    }
}

func (eh *ErrorHandler) HandleWithRetry(operation func() error) error {
    var lastErr error
    
    for attempt := 0; attempt <= eh.maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        lastErr = err
        
        // Check if error is recoverable
        if !eh.classifier.IsRecoverable(err) {
            return err
        }
        
        // Check severity - don't retry fatal errors
        severity := eh.classifier.GetSeverity(err)
        if severity >= errors.ErrorSeverityFatal {
            return err
        }
        
        if attempt < eh.maxRetries {
            fmt.Printf("Attempt %d failed, retrying in %v: %v\n", 
                attempt+1, eh.retryDelay, err)
            time.Sleep(eh.retryDelay)
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", eh.maxRetries+1, lastErr)
}

func main() {
    handler := NewErrorHandler()
    
    // Example operation that might fail
    operation := func() error {
        // Simulate intermittent failure
        if time.Now().UnixNano()%3 == 0 {
            return errors.NewApplicationError("Temporary network error", 1)
        }
        return nil
    }
    
    err := handler.HandleWithRetry(operation)
    if err != nil {
        fmt.Printf("Final error: %v\n", err)
    } else {
        fmt.Println("Operation succeeded")
    }
}
```

### Custom Error Types

```go
package main

import (
    "fmt"
    
    "github.com/bkpdir/pkg/errors"
)

// Custom error type that implements ErrorInterface
type ValidationError struct {
    *errors.ApplicationError
    Field string
    Value interface{}
}

func NewValidationError(field string, value interface{}, message string) *ValidationError {
    return &ValidationError{
        ApplicationError: errors.NewApplicationErrorWithContext(
            message,
            22, // Invalid argument exit code
            "validate",
            field,
            nil,
        ),
        Field: field,
        Value: value,
    }
}

func (ve *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s' with value '%v': %s", 
        ve.Field, ve.Value, ve.GetMessage())
}

// Custom error classifier
type AppErrorClassifier struct {
    *errors.DefaultErrorClassifier
}

func (aec *AppErrorClassifier) ClassifyError(err error) errors.ErrorCategory {
    if _, ok := err.(*ValidationError); ok {
        return errors.ErrorCategoryValidation
    }
    
    // Fall back to default classification
    return aec.DefaultErrorClassifier.ClassifyError(err)
}

func main() {
    // Create custom validation error
    err := NewValidationError("email", "invalid-email", "must be valid email address")
    
    fmt.Printf("Error: %s\n", err.Error())
    fmt.Printf("Field: %s\n", err.Field)
    fmt.Printf("Value: %v\n", err.Value)
    fmt.Printf("Status Code: %d\n", err.GetStatusCode())
    
    // Use custom classifier
    classifier := &AppErrorClassifier{
        DefaultErrorClassifier: errors.NewDefaultErrorClassifier(),
    }
    
    category := classifier.ClassifyError(err)
    fmt.Printf("Category: %v\n", category) // ErrorCategoryValidation
}
```

## Integration

### Integration with Other Packages

#### With pkg/cli

```go
import (
    "github.com/bkpdir/pkg/errors"
    "github.com/bkpdir/pkg/cli"
)

func createErrorAwareCLI() *cobra.Command {
    cmd := cli.NewCommandBuilder().NewCommand(
        "process", 
        "Process files", 
        "Process files with comprehensive error handling",
    )
    
    cli.NewCommandBuilder().WithHandler(cmd, func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return errors.NewApplicationError("No files specified", 22)
        }
        
        for _, file := range args {
            if err := processFile(file); err != nil {
                // Return structured error for CLI handling
                if appErr, ok := err.(*errors.ApplicationError); ok {
                    return appErr
                }
                
                // Wrap unknown errors
                return errors.NewApplicationErrorWithContext(
                    "File processing failed",
                    1,
                    "process",
                    file,
                    err,
                )
            }
        }
        
        return nil
    })
    
    return cmd
}
```

#### With pkg/formatter

```go
import (
    "github.com/bkpdir/pkg/errors"
    "github.com/bkpdir/pkg/formatter"
)

func formatErrors(err error, formatter formatter.ErrorFormatter) {
    classifier := errors.NewDefaultErrorClassifier()
    category := classifier.ClassifyError(err)
    
    switch category {
    case errors.ErrorCategoryDiskSpace:
        formatter.PrintDiskFullError(err)
    case errors.ErrorCategoryPermission:
        formatter.PrintPermissionError(err)
    case errors.ErrorCategoryFilesystem:
        formatter.PrintDirectoryNotFound(err)
    default:
        formatter.PrintError(err.Error())
    }
}
```

## Best Practices

### 1. Use Structured Errors

Always create structured errors with proper context:

```go
// Good
err := errors.NewApplicationErrorWithContext(
    "Failed to create backup",
    1,
    "backup",
    "/path/to/file",
    underlyingErr,
)

// Avoid
err := fmt.Errorf("backup failed: %v", underlyingErr)
```

### 2. Classify Errors for Better Handling

Use error classification to determine appropriate responses:

```go
classifier := errors.NewDefaultErrorClassifier()
if classifier.IsRecoverable(err) {
    // Retry logic
} else {
    // Fail fast
}
```

### 3. Preserve Error Context

Always preserve error context when wrapping errors:

```go
if err != nil {
    return errors.NewApplicationErrorWithContext(
        "Operation failed",
        statusCode,
        operation,
        path,
        err, // Preserve original error
    )
}
```

### 4. Use Appropriate Exit Codes

Follow standard exit code conventions:

```go
const (
    ExitSuccess         = 0
    ExitGeneralError    = 1
    ExitMisusage        = 2
    ExitCannotExecute   = 126
    ExitCommandNotFound = 127
    ExitInvalidExit     = 128
    ExitSIGINT          = 130
)
```

## License

Licensed under the MIT License. See LICENSE file for details.

---

**// EXTRACT-010: Package errors comprehensive documentation - Structured error handling with classification and context - 🔺** 