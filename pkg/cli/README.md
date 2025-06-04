# Package cli

[![Go Reference](https://pkg.go.dev/badge/github.com/bkpdir/pkg/cli.svg)](https://pkg.go.dev/github.com/bkpdir/pkg/cli)

## Overview

Package `cli` provides a reusable command-line interface framework extracted from the BkpDir application. It includes Cobra command patterns, flag handling, dry-run support, context-aware execution, and version management, designed to accelerate development of new CLI applications through tested patterns and reusable components.

### Key Features

- **Cobra Integration**: Built on top of the popular Cobra CLI framework with enhanced patterns
- **Dry-Run Support**: Built-in dry-run functionality for safe operation simulation
- **Context-Aware Execution**: Proper context propagation for cancellation and timeouts
- **Version Management**: Comprehensive version display and build information handling
- **Flag Standardization**: Consistent flag patterns across commands
- **Signal Handling**: Graceful shutdown support with signal handling
- **Command Building**: Builder patterns for rapid command construction

### Design Philosophy

The `cli` package was extracted to provide CLI application building blocks that have been proven in production. It emphasizes consistency, testability, and developer experience by providing high-level interfaces while maintaining flexibility for custom requirements.

## Installation

```bash
go get github.com/bkpdir/pkg/cli
```

## Quick Start

### Basic CLI Application

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    "github.com/bkpdir/pkg/cli"
)

func main() {
    // Define application information
    appInfo := cli.AppInfo{
        Name:  "myapp",
        Short: "My awesome CLI application",
        Long:  "A comprehensive CLI application built with pkg/cli framework",
        Build: cli.BuildInfo{
            Version:  "1.0.0",
            Date:     "2024-01-01",
            Commit:   "abc123",
            Platform: "linux/amd64",
        },
    }
    
    // Create root command builder
    rootBuilder := cli.NewRootCommandBuilder()
    flagMgr := cli.NewFlagManager()
    
    // Build root command
    rootCmd := rootBuilder.NewRootCommand(appInfo)
    rootBuilder.WithGlobalFlags(rootCmd, flagMgr)
    
    // Create subcommand
    cmdBuilder := cli.NewCommandBuilder()
    helloCmd := cmdBuilder.NewCommand("hello", "Say hello", "Print a greeting message")
    
    // Add handler
    cmdBuilder.WithHandler(helloCmd, func(cmd *cobra.Command, args []string) error {
        fmt.Println("Hello, World!")
        return nil
    })
    
    // Add to root
    rootCmd.AddCommand(helloCmd)
    
    // Execute
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

## API Reference

### Core Types

#### AppInfo

Application metadata and configuration:

```go
type AppInfo struct {
    Name  string    // Application name
    Short string    // Brief description
    Long  string    // Detailed description
    Build BuildInfo // Version and build information
}
```

#### BuildInfo

Version and build-time information:

```go
type BuildInfo struct {
    Version  string // Application version (e.g., "1.0.0")
    Date     string // Compilation date
    Commit   string // Git commit hash
    Platform string // Platform information
}
```

#### CommandContext

Shared context and dependencies for commands:

```go
type CommandContext struct {
    Context     context.Context // Context for cancellation and timeouts
    Output      io.Writer       // Output writer for command output
    ErrorOutput io.Writer       // Error writer for error messages
    DryRun      bool           // Indicates if operations should be simulated
}
```

### Core Interfaces

#### CommandBuilder

Builder pattern for creating Cobra commands:

```go
type CommandBuilder interface {
    NewCommand(name, short, long string) *cobra.Command
    WithHandler(cmd *cobra.Command, handler func(*cobra.Command, []string) error) *cobra.Command
    WithFlags(cmd *cobra.Command, flags []string) *cobra.Command
    WithSubcommands(parent *cobra.Command, children ...*cobra.Command) *cobra.Command
}
```

#### FlagManager

Consistent flag registration and binding:

```go
type FlagManager interface {
    AddGlobalFlags(cmd *cobra.Command) error
    AddDryRunFlag(cmd *cobra.Command, target *bool) error
    AddNoteFlag(cmd *cobra.Command, target *string) error
    AddConfigFlag(cmd *cobra.Command, target *bool) error
}
```

#### DryRunOperation

Interface for operations that support dry-run mode:

```go
type DryRunOperation interface {
    Execute(ctx CommandContext) error  // Performs the actual operation
    Describe() string                  // Returns a description of the operation
}
```

## Examples

### Using Dry-Run Support

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/bkpdir/pkg/cli"
)

type FileOperation struct {
    filename string
}

func (op *FileOperation) Execute(ctx cli.CommandContext) error {
    if ctx.DryRun {
        fmt.Printf("[DRY RUN] Would create file: %s\n", op.filename)
        return nil
    }
    
    // Actually create the file
    fmt.Printf("Creating file: %s\n", op.filename)
    return nil
}

func (op *FileOperation) Describe() string {
    return fmt.Sprintf("Create file %s", op.filename)
}

func main() {
    // Set up command with dry-run support
    var dryRun bool
    
    cmdBuilder := cli.NewCommandBuilder()
    flagMgr := cli.NewFlagManager()
    
    cmd := cmdBuilder.NewCommand("create", "Create file", "Create a new file")
    flagMgr.AddDryRunFlag(cmd, &dryRun)
    
    cmdBuilder.WithHandler(cmd, func(cmd *cobra.Command, args []string) error {
        if len(args) == 0 {
            return fmt.Errorf("filename required")
        }
        
        ctx := cli.CommandContext{
            Context: context.Background(),
            Output:  os.Stdout,
            DryRun:  dryRun,
        }
        
        operation := &FileOperation{filename: args[0]}
        dryRunMgr := cli.NewDryRunManager()
        
        return dryRunMgr.Execute(ctx, operation)
    })
    
    cmd.Execute()
}
```

## Integration

### Integration with Other Packages

#### With pkg/config

```go
import (
    "github.com/bkpdir/pkg/cli"
    "github.com/bkpdir/pkg/config"
)

func setupConfigAwareCLI() *cobra.Command {
    appInfo := cli.AppInfo{
        Name:  "configapp",
        Short: "Configuration-driven CLI app",
        Long:  "CLI application with configuration management",
    }
    
    rootBuilder := cli.NewRootCommandBuilder()
    flagMgr := cli.NewFlagManager()
    
    rootCmd := rootBuilder.NewRootCommand(appInfo)
    
    // Add config flag
    var showConfig bool
    flagMgr.AddConfigFlag(rootCmd, &showConfig)
    
    return rootCmd
}
```

## Best Practices

### 1. Use Builder Patterns

Leverage the builder patterns for consistent command creation:

```go
cmdBuilder := cli.NewCommandBuilder()
cmd := cmdBuilder.NewCommand("mycommand", "Short description", "Long description")
cmdBuilder.WithHandler(cmd, myHandler)
```

### 2. Implement Dry-Run Support

Always implement the DryRunOperation interface for operations that modify state:

```go
type MyOperation struct{}

func (op *MyOperation) Execute(ctx cli.CommandContext) error {
    if ctx.DryRun {
        fmt.Println("[DRY RUN] Would perform operation")
        return nil
    }
    // Actual implementation
    return nil
}

func (op *MyOperation) Describe() string {
    return "Description of what the operation does"
}
```

## License

Licensed under the MIT License. See LICENSE file for details.

---

**// EXTRACT-010: Package cli comprehensive documentation - Complete CLI framework guide with examples and patterns - 🔺** 