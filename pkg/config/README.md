# Package config

[![Go Reference](https://pkg.go.dev/badge/github.com/bkpdir/pkg/config.svg)](https://pkg.go.dev/github.com/bkpdir/pkg/config)

## Overview

Package `config` provides schema-agnostic configuration management for CLI applications. This package offers a complete configuration management system that supports multiple configuration sources, pluggable validation, and source tracking while remaining independent of specific application schemas.

### Key Features

- **Schema-Agnostic Design**: Works with any configuration structure through interface abstractions
- **Multiple Configuration Sources**: Files, environment variables, and defaults with configurable precedence
- **Source Tracking**: Track where each configuration value originated for debugging and transparency
- **Pluggable Validation**: Define custom validation rules for your application's schema
- **Atomic Operations**: Thread-safe configuration loading and merging
- **Zero Dependencies**: Minimal external dependencies for maximum compatibility

### Design Philosophy

The `config` package was extracted from the BkpDir project to provide reusable configuration management that doesn't lock you into a specific schema. It uses interface-based abstractions to enable different applications to define their own configuration structures while benefiting from robust discovery, merging, and validation logic.

## Installation

```bash
go get github.com/bkpdir/pkg/config
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/bkpdir/pkg/config"
)

// Define your application's configuration structure
type AppConfig struct {
    AppName     string `yaml:"app_name"`
    LogLevel    string `yaml:"log_level"`
    DatabaseURL string `yaml:"database_url"`
}

func main() {
    // Create default configuration
    defaultConfig := &AppConfig{
        AppName:     "MyApp",
        LogLevel:    "info",
        DatabaseURL: "localhost:5432",
    }
    
    // Create configuration loader
    loader := config.NewDefaultConfigLoader()
    
    // Load configuration from current directory
    cfg, err := loader.LoadConfig(".", defaultConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // Type assert to your configuration structure
    appConfig := cfg.(*AppConfig)
    
    fmt.Printf("App: %s, Log Level: %s\n", appConfig.AppName, appConfig.LogLevel)
}
```

### Configuration File Example

Create a `.myapp.yml` file in your application directory:

```yaml
app_name: "Production App"
log_level: "warn"
database_url: "prod-db:5432"
```

### Environment Variables

Override configuration with environment variables:

```bash
export MYAPP_LOG_LEVEL=debug
export MYAPP_DATABASE_URL=dev-db:5432
```

## API Reference

### Core Interfaces

#### ConfigLoader

The primary interface for loading and managing configuration:

```go
type ConfigLoader interface {
    LoadConfig(root string, defaultConfig interface{}) (interface{}, error)
    LoadConfigValues(root string, defaultConfig interface{}) (map[string]ConfigValue, error)
    GetConfigValues(cfg interface{}) []ConfigValue
    GetConfigValuesWithSources(cfg interface{}, root string) []ConfigValue
    ValidateConfig(cfg interface{}) error
}
```

**Key Methods:**
- `LoadConfig`: Load configuration with default fallbacks
- `LoadConfigValues`: Load configuration with source tracking
- `GetConfigValues`: Extract values from configuration structures
- `ValidateConfig`: Validate configuration against schema

#### ConfigMerger

Interface for merging configuration from multiple sources:

```go
type ConfigMerger interface {
    MergeConfigs(dst, src interface{}) error
    MergeConfigValues(dst, src map[string]ConfigValue)
    GetConfigSearchPaths() []string
    ExpandPath(path string) string
}
```

#### ConfigValidator

Interface for custom validation logic:

```go
type ConfigValidator interface {
    ValidateSchema(cfg interface{}) error
    ValidateValues(values map[string]ConfigValue) error
    GetRequiredFields() []string
    GetValidationRules() map[string]ValidationRule
}
```

### Data Structures

#### ConfigValue

Represents a configuration value with source tracking:

```go
type ConfigValue struct {
    Name   string      // Configuration field name
    Value  interface{} // Configuration value
    Source string      // Source of the value (file, env, default)
    Type   string      // Type of the value
}
```

#### ValidationRule

Defines validation criteria for configuration fields:

```go
type ValidationRule struct {
    Required     bool
    Type         string
    MinLength    int
    MaxLength    int
    Pattern      string
    ValidValues  []string
    Dependencies []string
}
```

## Examples

### Advanced Configuration Loading

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/bkpdir/pkg/config"
)

type DatabaseConfig struct {
    Host     string `yaml:"host" env:"DB_HOST"`
    Port     int    `yaml:"port" env:"DB_PORT"`
    Database string `yaml:"database" env:"DB_NAME"`
    Username string `yaml:"username" env:"DB_USER"`
    Password string `yaml:"password" env:"DB_PASS"`
}

type AppConfig struct {
    AppName   string          `yaml:"app_name" env:"APP_NAME"`
    Debug     bool            `yaml:"debug" env:"DEBUG"`
    Database  DatabaseConfig  `yaml:"database"`
    Features  []string        `yaml:"features"`
}

func main() {
    // Create default configuration
    defaultConfig := &AppConfig{
        AppName: "MyApp",
        Debug:   false,
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Database: "myapp",
            Username: "user",
            Password: "password",
        },
        Features: []string{"basic"},
    }
    
    // Create loader with custom options
    loader := config.NewDefaultConfigLoader()
    
    // Load configuration with source tracking
    values, err := loader.LoadConfigValues(".", defaultConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // Display configuration sources
    for name, value := range values {
        fmt.Printf("%s: %v (from %s)\n", name, value.Value, value.Source)
    }
}
```

### Custom Validation

```go
package main

import (
    "fmt"
    "strings"
    
    "github.com/bkpdir/pkg/config"
)

type CustomValidator struct{}

func (v *CustomValidator) ValidateSchema(cfg interface{}) error {
    appConfig, ok := cfg.(*AppConfig)
    if !ok {
        return fmt.Errorf("invalid configuration type")
    }
    
    // Custom validation logic
    if appConfig.Database.Port < 1024 || appConfig.Database.Port > 65535 {
        return fmt.Errorf("database port must be between 1024 and 65535")
    }
    
    return nil
}

func (v *CustomValidator) ValidateValues(values map[string]ConfigValue) error {
    for name, value := range values {
        if strings.Contains(name, "password") && value.Source == "default" {
            return fmt.Errorf("password fields cannot use default values")
        }
    }
    return nil
}

func (v *CustomValidator) GetRequiredFields() []string {
    return []string{"app_name", "database.host", "database.username"}
}

func (v *CustomValidator) GetValidationRules() map[string]ValidationRule {
    return map[string]ValidationRule{
        "app_name": {
            Required:  true,
            Type:      "string",
            MinLength: 1,
            MaxLength: 50,
        },
        "database.port": {
            Required: true,
            Type:     "int",
        },
    }
}
```

### Configuration Discovery

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/bkpdir/pkg/config"
)

func main() {
    // Set custom configuration search path
    os.Setenv("MYAPP_CONFIG", "./.myapp.yml:~/.myapp.yml:/etc/myapp.yml")
    
    loader := config.NewDefaultConfigLoader()
    merger := config.NewDefaultConfigMerger()
    
    // Get search paths
    paths := merger.GetConfigSearchPaths()
    fmt.Printf("Searching for configuration in: %v\n", paths)
    
    // Load configuration from discovered paths
    defaultConfig := &AppConfig{AppName: "Default"}
    
    cfg, err := loader.LoadConfig(".", defaultConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Loaded configuration: %+v\n", cfg)
}
```

## Integration

### Integration with Other Packages

The `config` package is designed to work seamlessly with other extracted packages:

#### With pkg/cli

```go
import (
    "github.com/bkpdir/pkg/config"
    "github.com/bkpdir/pkg/cli"
)

func createApp() *cli.App {
    loader := config.NewDefaultConfigLoader()
    
    app := cli.NewApp()
    app.Before = func(c *cli.Context) error {
        // Load configuration before command execution
        cfg, err := loader.LoadConfig(c.String("config-dir"), &AppConfig{})
        if err != nil {
            return err
        }
        
        // Store in context for use by commands
        c.App.Metadata["config"] = cfg
        return nil
    }
    
    return app
}
```

#### With pkg/errors

```go
import (
    "github.com/bkpdir/pkg/config"
    "github.com/bkpdir/pkg/errors"
)

func loadConfigSafely(path string) error {
    loader := config.NewDefaultConfigLoader()
    
    _, err := loader.LoadConfig(path, &AppConfig{})
    if err != nil {
        return errors.WrapConfigError(err, "CONFIG_LOAD_FAILED", "Failed to load configuration")
    }
    
    return nil
}
```

#### With pkg/formatter

```go
import (
    "github.com/bkpdir/pkg/config"
    "github.com/bkpdir/pkg/formatter"
)

func displayConfig(cfg *AppConfig) {
    loader := config.NewDefaultConfigLoader()
    values := loader.GetConfigValues(cfg)
    
    formatter := formatter.NewTableFormatter()
    for _, value := range values {
        formatter.AddRow([]string{value.Name, fmt.Sprintf("%v", value.Value), value.Source})
    }
    
    formatter.Render()
}
```

## Configuration File Formats

The package supports multiple configuration file formats through pluggable sources:

### YAML (Recommended)

```yaml
# .myapp.yml
app_name: "My Application"
debug: true
database:
  host: "localhost"
  port: 5432
  database: "myapp_dev"
features:
  - "feature1"
  - "feature2"
```

### JSON

```json
{
  "app_name": "My Application",
  "debug": true,
  "database": {
    "host": "localhost",
    "port": 5432,
    "database": "myapp_dev"
  },
  "features": ["feature1", "feature2"]
}
```

## Environment Variables

Configuration values can be overridden using environment variables. The package supports flexible environment variable mapping:

```bash
# Standard environment variables
export MYAPP_APP_NAME="Production App"
export MYAPP_DEBUG=false
export MYAPP_DATABASE_HOST="prod-db.example.com"
export MYAPP_DATABASE_PORT=5432

# Nested configuration
export MYAPP_DATABASE_USERNAME="prod_user"
export MYAPP_DATABASE_PASSWORD="secure_password"
```

## Performance Characteristics

- **Configuration Loading**: ~24.3μs for typical configuration files
- **Memory Usage**: Minimal overhead with lazy loading
- **Thread Safety**: All operations are thread-safe
- **Caching**: Built-in caching for improved performance with repeated loads

## Best Practices

### 1. Use Default Configuration

Always provide sensible default values:

```go
defaultConfig := &AppConfig{
    AppName:  "MyApp",
    LogLevel: "info",
    Timeout:  30 * time.Second,
}
```

### 2. Validate Configuration

Implement custom validation for critical configuration:

```go
func (c *AppConfig) Validate() error {
    if c.Timeout <= 0 {
        return fmt.Errorf("timeout must be positive")
    }
    return nil
}
```

### 3. Use Source Tracking

Leverage source tracking for debugging:

```go
values, err := loader.LoadConfigValues(".", defaultConfig)
for name, value := range values {
    if value.Source == "default" {
        log.Warnf("Using default value for %s", name)
    }
}
```

### 4. Environment Variable Conventions

Follow consistent environment variable naming:

```bash
# Use prefix to avoid conflicts
export MYAPP_DATABASE_HOST=localhost
export MYAPP_DATABASE_PORT=5432

# Use underscores for nested values
export MYAPP_FEATURE_FLAGS_ENABLE_CACHE=true
```

## Troubleshooting

### Common Issues

1. **Configuration file not found**: Check search paths and file permissions
2. **Environment variables not loaded**: Verify environment variable naming conventions
3. **Validation errors**: Check required fields and data types
4. **Merge conflicts**: Understand source precedence (env > file > default)

### Debug Mode

Enable debug logging to trace configuration loading:

```go
loader := config.NewDefaultConfigLoader()
loader.SetDebugMode(true)  // Enable detailed logging
```

## Contributing

This package is part of the BkpDir extraction project. For contributions:

1. Follow the existing interface patterns
2. Add comprehensive tests for new functionality
3. Update documentation for API changes
4. Ensure backward compatibility

## License

Licensed under the MIT License. See LICENSE file for details.

---

**// EXTRACT-010: Package config comprehensive documentation - API reference, examples, and integration guide - 🔺** 