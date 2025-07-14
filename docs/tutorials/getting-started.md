# 🚀 Getting Started: Building Your First CLI Application

> **[HIGH] EXTRACT-008: Tutorial series creation - 📚 Step-by-step learning materials**

## 🎯 Overview

This tutorial walks you through building your first CLI application using the extracted packages from the bkpdir project. By the end of this tutorial, you'll have a working CLI application that demonstrates configuration management, error handling, and output formatting.

## 📋 Prerequisites

- Go 1.19 or later
- Basic understanding of Go programming
- Familiarity with command-line applications

## 🏗️ Project Setup

### Step 1: Initialize Your Project

Create a new Go module for your CLI application:

```bash
mkdir my-cli-app
cd my-cli-app
go mod init github.com/yourusername/my-cli-app
```

### Step 2: Add the Extracted Packages

For this tutorial, we'll use the extracted packages. In a real project, you would add them as dependencies:

```bash
# In the future, these will be available as separate modules
# go get github.com/extracted-packages/config
# go get github.com/extracted-packages/errors
# go get github.com/extracted-packages/cli
```

For now, copy the `pkg/` directory from the bkpdir project to your project root.

### Step 3: Create Project Structure

```
my-cli-app/
├── main.go
├── cmd/
│   ├── create.go
│   ├── list.go
│   └── delete.go
├── config.yaml
├── pkg/          # Copied from bkpdir
│   ├── config/
│   ├── errors/
│   ├── cli/
│   └── formatter/
└── go.mod
```

## [ACTION:core-functionality] Building the Application

### Step 1: Create the Main Entry Point

Create `main.go`:

```go
// [HIGH] EXTRACT-008: Tutorial series creation - 📚 Getting started main function
package main

import (
    "context"
    "fmt"
    "os"
    "slices"

    "./pkg/cli"
    "./pkg/config"
    "./pkg/errors"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig(".")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }

    // Create CLI application
    app := cli.NewApplication("my-cli-app", "1.0.0")
    app.SetConfig(cfg)

    // Check for dry-run flag
    if slices.Contains(os.Args, "--dry-run") {
        app.SetDryRun(true)
    }

    // Register commands
    app.AddCommand("create", "Create a new resource", createCommand)
    app.AddCommand("list", "List all resources", listCommand)
    app.AddCommand("delete", "Delete a resource", deleteCommand)

    // Execute with proper error handling
    ctx := context.Background()
    if err := app.Execute(ctx, os.Args[1:]); err != nil {
        handleError(err)
    }
}

func handleError(err error) {
    // [HIGH] EXTRACT-008: Tutorial series creation - 📚 Error handling demonstration
    if userErr := errors.AsUserError(err); userErr != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", userErr.UserMessage())
        if userErr.Code() != "" {
            fmt.Fprintf(os.Stderr, "Code: %s\n", userErr.Code())
        }
        os.Exit(userErr.ExitCode())
    }

    if sysErr := errors.AsSystemError(err); sysErr != nil {
        fmt.Fprintf(os.Stderr, "System error: %s\n", sysErr.UserMessage())
        if sysErr.Code() != "" {
            fmt.Fprintf(os.Stderr, "Code: %s\n", sysErr.Code())
        }
        fmt.Fprintf(os.Stderr, "Underlying error: %v\n", sysErr.Unwrap())
        os.Exit(sysErr.ExitCode())
    }

    if appErr := errors.AsApplicationError(err); appErr != nil {
        fmt.Fprintf(os.Stderr, "Application error: %s\n", appErr.UserMessage())
        os.Exit(appErr.ExitCode())
    }

    // Unknown error type
    fmt.Fprintf(os.Stderr, "Unexpected error: %v\n", err)
    os.Exit(1)
}
```

### Step 2: Create Configuration File

Create `config.yaml`:

```yaml
# [HIGH] EXTRACT-008: Tutorial series creation - 📚 Configuration example
app:
  name: "my-cli-app"
  version: "1.0.0"
  debug: false

output:
  format: "text"
  verbose: false
  color: true

resources:
  storage_path: "./data"
  max_items: 100
```

### Step 3: Implement Commands

Create `cmd/create.go`:

```go
// [HIGH] EXTRACT-008: Tutorial series creation - 📚 Create command implementation
package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "./pkg/cli"
    "./pkg/errors"
    "./pkg/formatter"
)

func createCommand(ctx context.Context, args []string) error {
    if len(args) == 0 {
        return errors.NewUserError(
            "Resource name is required. Usage: create <name>",
            "MISSING_RESOURCE_NAME")
    }

    resourceName := args[0]
    
    // Get CLI context
    cliCtx := cli.FromContext(ctx)
    
    // Create formatter for output
    output := formatter.New(cliCtx.Config)
    
    if cliCtx.Verbose {
        output.Printf("Creating resource '%s'...\n", resourceName)
    }

    // Check for dry-run mode
    if cliCtx.DryRun {
        output.Printf("DRY RUN: Would create resource '%s'\n", resourceName)
        return nil
    }

    // Create the resource
    if err := createResource(resourceName, cliCtx.Config); err != nil {
        return errors.WrapError(err, 
            fmt.Sprintf("Failed to create resource '%s'", resourceName),
            "RESOURCE_CREATE_ERROR")
    }

    // Success message
    output.Printf("✅ Successfully created resource '%s'\n", resourceName)
    
    return nil
}

func createResource(name string, cfg *config.Config) error {
    // Get storage path from configuration
    storagePath := cfg.GetString("resources.storage_path", "./data")
    
    // Ensure storage directory exists
    if err := os.MkdirAll(storagePath, 0755); err != nil {
        return errors.NewSystemError(
            "Failed to create storage directory",
            "STORAGE_DIR_ERROR",
            err)
    }

    // Create resource file
    resourceFile := filepath.Join(storagePath, name+".txt")
    content := fmt.Sprintf("Resource: %s\nCreated: %s\n", name, time.Now().Format(time.RFC3339))
    
    if err := os.WriteFile(resourceFile, []byte(content), 0644); err != nil {
        return errors.NewSystemError(
            "Failed to write resource file",
            "RESOURCE_WRITE_ERROR",
            err)
    }

    return nil
}
```

Create `cmd/list.go`:

```go
// [HIGH] EXTRACT-008: Tutorial series creation - 📚 List command implementation
package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "./pkg/cli"
    "./pkg/errors"
    "./pkg/formatter"
)

func listCommand(ctx context.Context, args []string) error {
    // Get CLI context
    cliCtx := cli.FromContext(ctx)
    
    // Create formatter for output
    output := formatter.New(cliCtx.Config)
    
    if cliCtx.Verbose {
        output.Printf("Listing resources...\n")
    }

    // Get storage path from configuration
    storagePath := cliCtx.Config.GetString("resources.storage_path", "./data")
    
    // Check if storage directory exists
    if _, err := os.Stat(storagePath); os.IsNotExist(err) {
        output.Printf("No resources found (storage directory doesn't exist)\n")
        return nil
    }

    // Read directory
    entries, err := os.ReadDir(storagePath)
    if err != nil {
        return errors.NewSystemError(
            "Failed to read storage directory",
            "STORAGE_READ_ERROR",
            err)
    }

    // Filter for resource files
    var resources []string
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".txt") {
            resourceName := strings.TrimSuffix(entry.Name(), ".txt")
            resources = append(resources, resourceName)
        }
    }

    // Display results
    if len(resources) == 0 {
        output.Printf("No resources found\n")
        return nil
    }

    // Use template formatting for better output
    template := `Resources ({{len .Resources}}):
{{range $i, $resource := .Resources}}{{add $i 1}}. {{$resource}}
{{end}}`

    data := map[string]interface{}{
        "Resources": resources,
    }

    return output.PrintTemplate(template, data)
}
```

Create `cmd/delete.go`:

```go
// [HIGH] EXTRACT-008: Tutorial series creation - 📚 Delete command implementation
package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "./pkg/cli"
    "./pkg/errors"
    "./pkg/formatter"
)

func deleteCommand(ctx context.Context, args []string) error {
    if len(args) == 0 {
        return errors.NewUserError(
            "Resource name is required. Usage: delete <name>",
            "MISSING_RESOURCE_NAME")
    }

    resourceName := args[0]
    
    // Get CLI context
    cliCtx := cli.FromContext(ctx)
    
    // Create formatter for output
    output := formatter.New(cliCtx.Config)
    
    if cliCtx.Verbose {
        output.Printf("Deleting resource '%s'...\n", resourceName)
    }

    // Check for dry-run mode
    if cliCtx.DryRun {
        output.Printf("DRY RUN: Would delete resource '%s'\n", resourceName)
        return nil
    }

    // Delete the resource
    if err := deleteResource(resourceName, cliCtx.Config); err != nil {
        return errors.WrapError(err,
            fmt.Sprintf("Failed to delete resource '%s'", resourceName),
            "RESOURCE_DELETE_ERROR")
    }

    // Success message
    output.Printf("✅ Successfully deleted resource '%s'\n", resourceName)
    
    return nil
}

func deleteResource(name string, cfg *config.Config) error {
    // Get storage path from configuration
    storagePath := cfg.GetString("resources.storage_path", "./data")
    
    // Build resource file path
    resourceFile := filepath.Join(storagePath, name+".txt")
    
    // Check if resource exists
    if _, err := os.Stat(resourceFile); os.IsNotExist(err) {
        return errors.NewUserError(
            fmt.Sprintf("Resource '%s' does not exist", name),
            "RESOURCE_NOT_FOUND")
    }

    // Delete the resource file
    if err := os.Remove(resourceFile); err != nil {
        return errors.NewSystemError(
            "Failed to delete resource file",
            "RESOURCE_DELETE_ERROR",
            err)
    }

    return nil
}
```

## 🧪 Testing Your Application

### Step 1: Build the Application

```bash
go build -o my-cli-app
```

### Step 2: Test Basic Commands

```bash
# Show help
./my-cli-app help

# Show version
./my-cli-app version

# Show configuration
./my-cli-app config
./my-cli-app config --verbose

# Test create command
./my-cli-app create myresource

# Test list command
./my-cli-app list
./my-cli-app list --verbose

# Test delete command
./my-cli-app delete myresource

# Test dry-run mode
./my-cli-app create test-resource --dry-run
./my-cli-app delete test-resource --dry-run
```

### Step 3: Test Error Handling

```bash
# Test missing arguments
./my-cli-app create

# Test deleting non-existent resource
./my-cli-app delete nonexistent

# Test invalid command
./my-cli-app invalid-command
```

## 🎯 Key Learning Points

### 1. Configuration Management

- Configuration is loaded automatically from `config.yaml`
- Environment variables can override configuration values
- Default values are provided for missing configuration

### 2. Error Handling

- Different error types for different scenarios
- Structured error codes for programmatic handling
- User-friendly error messages
- Proper exit codes

### 3. CLI Framework

- Automatic help and version commands
- Dry-run support built-in
- Context passing between commands
- Consistent flag handling

### 4. Output Formatting

- Printf-style formatting for simple output
- Template-based formatting for complex output
- Configuration-driven output behavior

## [ACTION:migration] Next Steps

Now that you have a basic CLI application working, you can:

1. **Add More Features**: Explore the advanced patterns tutorial
2. **Add Git Integration**: Include repository information in your output
3. **Add File Operations**: Use pkg/fileops for safe file handling
4. **Add Concurrent Processing**: Use pkg/processing for parallel operations
5. **Improve Error Handling**: Add resource management with pkg/resources

## 📚 Related Tutorials

- [Advanced Patterns](advanced-patterns.md) - Complex integration scenarios
- [Performance Optimization](performance-optimization.md) - Performance best practices
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

## 📖 References

- [Integration Guide](../integration-guide.md) - Comprehensive usage patterns
- [Package Reference](../package-reference.md) - Detailed API documentation
- [CLI Template](../../cmd/cli-template/README.md) - Complete working example

---

**[HIGH] EXTRACT-008: Tutorial series creation - 📚 Complete getting started guide with working examples** 