# Migration Guide: From Monolithic to Package-Based Structure

> **📑 EXTRACT-008: Migration Guide** - Moving from monolithic to package-based structure using extracted packages

## 📑 Purpose

This guide provides step-by-step instructions for migrating existing CLI applications from a monolithic structure (like the original bkpdir application) to a modular package-based architecture using the extracted components.

## 🎯 Target Audience

- **Developers** migrating existing CLI applications
- **Project maintainers** looking to modularize their codebase
- **Teams** adopting the extracted package ecosystem

## 📊 Migration Overview

### Before: Monolithic Structure
```
myapp/
├── main.go                    # 1500+ lines with everything
├── config.go                  # Configuration logic mixed with schema
├── archive.go                 # Business logic + file operations
├── backup.go                  # More business logic + file operations
├── comparison.go              # Directory comparison logic
├── verify.go                  # Verification logic
├── formatter.go               # Output formatting (1677 lines)
├── git.go                     # Git integration
├── resources.go               # Resource management
└── errors.go                  # Error handling
```

### After: Package-Based Structure
```
myapp/
├── main.go                    # 20-30 lines - entry point only
├── go.mod                     # Dependency management
├── cmd/                       # Command definitions
│   ├── root.go               # Root command setup
│   ├── create.go             # Create command
│   └── version.go            # Version command
├── config/                    # Configuration files
│   ├── default.yml           # Default configuration
│   └── example.yml           # Example configuration
└── pkg/                       # Extracted packages (imported)
    ├── config/               # Schema-agnostic configuration
    ├── errors/               # Structured error handling
    ├── resources/            # Resource management
    ├── formatter/            # Output formatting
    ├── git/                  # Git integration
    ├── cli/                  # Command framework
    ├── fileops/              # File operations
    └── processing/           # Concurrent processing
```

## 🚀 Phase 1: Pre-Migration Analysis

### Step 1: Assess Current Structure

**[ACTION:discovery] Analyze your monolithic application:**

```bash
# 1. Count lines of code in major files
wc -l *.go

# 2. Identify component boundaries
grep -n "func " *.go | head -20

# 3. Map dependencies between components
grep -r "import" . --include="*.go"

# 4. Identify configuration patterns
grep -r "config\|Config" . --include="*.go"
```

**📋 Create assessment report:**

```markdown
## Current Application Assessment

### File Analysis
- main.go: [X] lines - Command orchestration + business logic
- config.go: [X] lines - Configuration management + schema
- [component].go: [X] lines - [Description]

### Dependency Complexity
- [ ] Circular dependencies identified
- [ ] Tightly coupled components identified
- [ ] Configuration coupling assessed

### Extract Candidates
- [ ] Configuration logic
- [ ] Error handling patterns
- [ ] File operation patterns
- [ ] Output formatting logic
- [ ] CLI command structure
```

### Step 2: Choose Migration Strategy

**[MEDIUM] Progressive Migration (Recommended)**
- Migrate one component at a time
- Maintain backward compatibility
- Validate each step before proceeding

**[MEDIUM] Complete Rewrite**
- Start fresh with package-based structure
- Faster but higher risk
- Suitable for smaller applications

**[MEDIUM] Hybrid Approach**
- Migrate core components progressively
- Rewrite complex monolithic sections
- Best for medium-sized applications

## ⚡ Phase 2: Component Migration

### Step 1: Setup New Project Structure

**[ACTION:core-functionality] Initialize new project:**

```bash
# 1. Create new project directory
mkdir myapp-v2
cd myapp-v2

# 2. Initialize Go module
go mod init myapp

# 3. Add extracted package dependencies
go mod edit -require=github.com/spf13/cobra@v1.9.1
go mod edit -require=gopkg.in/yaml.v2@v2.4.0

# 4. Create directory structure
mkdir -p cmd config pkg docs
```

**[ACTION:format-processing] Create basic structure:**

```bash
# Copy and adapt from CLI template
cp -r /path/to/cmd/cli-template/* .
# Customize for your application
```

### Step 2: Migrate Configuration

**[ACTION:core-functionality] Before (Monolithic):**
```go
// config.go - 800 lines mixed with business logic
type Config struct {
    // Mixed application-specific and generic settings
    MyAppSpecificSetting string `yaml:"my_app_setting"`
    ArchiveDirPath      string `yaml:"archive_dir_path"`
    // ... hundreds of lines of mixed concerns
}

func LoadConfig() (*Config, error) {
    // Hardcoded configuration loading logic
    // Mixed with schema validation
    // Tightly coupled to file system
}
```

**[ACTION:core-functionality] After (Package-Based):**
```go
// cmd/root.go - Clean configuration usage
import "myapp/pkg/config"

type AppConfig struct {
    // Only application-specific settings
    MyAppSpecificSetting string `yaml:"my_app_setting"`
    DatabaseURL         string `yaml:"database_url"`
}

func initConfig() {
    // Use extracted config package
    loader := config.NewConfigLoader(config.DefaultOptions())
    appCfg := &AppConfig{}
    
    if err := loader.LoadInto(appCfg); err != nil {
        log.Fatal(err)
    }
}
```

**[ACTION:discovery] Migration Steps:**

1. **Extract Schema-Specific Settings:**
   ```go
   // Identify application-specific configuration
   type AppConfig struct {
       AppName     string `yaml:"app_name"`
       DatabaseURL string `yaml:"database_url"`
       LogLevel    string `yaml:"log_level"`
   }
   ```

2. **Implement ConfigLoader Interface:**
   ```go
   import "myapp/pkg/config"
   
   loader := config.NewConfigLoader(config.Options{
       SearchPaths: []string{"./.myapp.yml", "~/.myapp.yml"},
       EnvPrefix:   "MYAPP",
   })
   ```

3. **Migrate Configuration Files:**
   ```yaml
   # config/default.yml
   app:
     name: "myapp"
     version: "1.0.0"
   
   database:
     url: "sqlite://./data.db"
   
   logging:
     level: "info"
   ```

### Step 3: Migrate Error Handling

**[ACTION:core-functionality] Before (Monolithic):**
```go
// Scattered error handling throughout codebase
func ProcessFile(path string) error {
    if _, err := os.Stat(path); err != nil {
        return fmt.Errorf("file not found: %v", err)
    }
    // No structured error types
    // No consistent error handling patterns
}
```

**[ACTION:core-functionality] After (Package-Based):**
```go
import "myapp/pkg/errors"

func ProcessFile(path string) error {
    if _, err := os.Stat(path); err != nil {
        return errors.NewFileError("file not found", errors.FileNotFound, err).
               WithOperation("ProcessFile").
               WithPath(path)
    }
    return nil
}
```

**[ACTION:discovery] Migration Steps:**

1. **Identify Error Categories:**
   - File operation errors
   - Configuration errors
   - Business logic errors
   - Network/external service errors

2. **Replace Error Patterns:**
   ```go
   // Old: Simple error messages
   return fmt.Errorf("failed to process: %v", err)
   
   // New: Structured errors
   return errors.NewProcessingError("failed to process", err).
          WithContext("operation", "file-processing")
   ```

3. **Add Error Recovery:**
   ```go
   if err := ProcessFile(path); err != nil {
       if errors.IsFileNotFound(err) {
           // Handle missing file specifically
           return handleMissingFile(path)
       }
       return err
   }
   ```

### Step 4: Migrate Output Formatting

**[ACTION:core-functionality] Before (Monolithic):**
```go
// formatter.go - 1677 lines of mixed formatting logic
func PrintResult(message string, args ...interface{}) {
    fmt.Printf(message, args...)
}

func PrintError(err error) {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
```

**[ACTION:core-functionality] After (Package-Based):**
```go
import "myapp/pkg/formatter"

func setupFormatter() formatter.FormatterInterface {
    return formatter.NewTemplateFormatter(formatter.Options{
        Templates: map[string]string{
            "success": "✅ {{.operation}} completed: {{.result}}\n",
            "error":   "❌ {{.operation}} failed: {{.error}}\n",
        },
    })
}

func PrintResult(fmt formatter.FormatterInterface, operation, result string) {
    fmt.PrintTemplate("success", map[string]interface{}{
        "operation": operation,
        "result":    result,
    })
}
```

### Step 5: Migrate File Operations

**[ACTION:core-functionality] Before (Monolithic):**
```go
// Scattered file operations with inconsistent error handling
func CopyFile(src, dst string) error {
    // Basic file copy without proper error handling
    // No atomic operations
    // No cleanup on failure
}
```

**[ACTION:core-functionality] After (Package-Based):**
```go
import "myapp/pkg/fileops"

func CopyFile(src, dst string) error {
    copier := fileops.NewFileCopier(fileops.DefaultOptions())
    return copier.CopyWithValidation(src, dst)
}
```

### Step 6: Migrate CLI Commands

**[ACTION:core-functionality] Before (Monolithic):**
```go
// main.go - Everything mixed together
func main() {
    if len(os.Args) > 1 && os.Args[1] == "process" {
        // Inline command handling
        processFiles(os.Args[2:])
    }
    // ... hundreds of lines of command logic
}
```

**[ACTION:core-functionality] After (Package-Based):**
```go
// cmd/process.go
import "myapp/pkg/cli"

func NewProcessCommand() *cobra.Command {
    return cli.NewCommandBuilder().
        NewCommand("process", "Process files", "Process files with extracted components").
        WithHandler(handleProcess).
        WithFlags([]string{"dry-run", "verbose"})
}

func handleProcess(cmd *cobra.Command, args []string) error {
    ctx := cli.CommandContext{
        Context: cmd.Context(),
        Output:  os.Stdout,
        DryRun:  dryRun,
    }
    return processFiles(ctx, args)
}
```

## [ACTION:migration] Phase 3: Integration and Testing

### Step 1: Validate Component Integration

**🧪 Create integration tests:**

```go
// integration_test.go
func TestMigrationComponents(t *testing.T) {
    tests := []struct {
        name      string
        component string
        test      func(t *testing.T)
    }{
        {"Config Loading", "config", testConfigIntegration},
        {"Error Handling", "errors", testErrorIntegration},
        {"File Operations", "fileops", testFileOpsIntegration},
        {"CLI Commands", "cli", testCLIIntegration},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, tt.test)
    }
}
```

**[ACTION:discovery] Verify functionality:**

```bash
# 1. Run integration tests
go test ./... -v

# 2. Compare outputs with original
./original-app process file.txt > original.out
./new-app process file.txt > new.out
diff original.out new.out

# 3. Performance benchmarking
go test -bench=. ./...
```

### Step 2: Migration Validation Checklist

**✅ Functional Validation**
- [ ] All original commands work identically
- [ ] Error messages maintain consistency
- [ ] Configuration loading behavior preserved
- [ ] Output formatting matches original
- [ ] File operations produce identical results

**✅ Performance Validation**
- [ ] No significant performance regression (<5%)
- [ ] Memory usage within acceptable limits
- [ ] Startup time not degraded
- [ ] Large file operations perform similarly

**✅ Quality Validation**
- [ ] Test coverage maintained or improved
- [ ] Linting passes for all new code
- [ ] Documentation updated appropriately
- [ ] No breaking changes introduced

## 🏁 Phase 4: Deployment and Rollback

### Step 1: Gradual Rollout Strategy

**[ACTION:migration] Blue-Green Deployment:**

```bash
# 1. Deploy new version alongside old
cp original-app app-v1
cp new-app app-v2

# 2. Route subset of traffic to new version
# Use feature flags or environment switching

# 3. Monitor metrics and error rates
# Compare logs and performance data

# 4. Full cutover once validated
mv app-v2 app
```

**[ACTION:migration] Feature Flag Migration:**

```go
// Use feature flags for gradual migration
func ProcessFiles(files []string) error {
    if useExtractedPackages() {
        return processFilesV2(files) // New package-based implementation
    }
    return processFilesV1(files) // Original monolithic implementation
}
```

### Step 2: Rollback Plan

**🚨 Emergency Rollback:**

```bash
# 1. Immediate rollback procedure
cp app-v1 app
systemctl restart myapp

# 2. Check service health
curl -f http://localhost:8080/health

# 3. Monitor logs for stability
tail -f /var/log/myapp.log
```

**📊 Rollback Triggers:**
- Error rate increase >5%
- Performance degradation >10%
- Critical functionality broken
- User complaints or support tickets

## [ACTION:maintenance] Common Migration Patterns

### Pattern 1: Configuration Migration

```go
// Step 1: Create adapter for backward compatibility
type ConfigAdapter struct {
    newConfig *AppConfig
    loader    config.LoaderInterface
}

func (a *ConfigAdapter) GetLegacySetting(key string) interface{} {
    // Map old configuration keys to new structure
    mapping := map[string]func() interface{}{
        "old_key": func() interface{} { return a.newConfig.NewKey },
    }
    return mapping[key]()
}

// Step 2: Gradual replacement
func LoadConfigWithMigration() *Config {
    adapter := &ConfigAdapter{
        loader: config.NewConfigLoader(config.DefaultOptions()),
    }
    
    // Load new config structure
    if err := adapter.loader.LoadInto(&adapter.newConfig); err != nil {
        // Fallback to old config loading
        return loadLegacyConfig()
    }
    
    return adapter
}
```

### Pattern 2: Error Handling Migration

```go
// Step 1: Create error adapter
func WrapLegacyError(err error, operation string) error {
    if err == nil {
        return nil
    }
    
    // Convert legacy errors to structured errors
    return errors.NewLegacyError(err.Error(), err).
           WithOperation(operation).
           WithTimestamp(time.Now())
}

// Step 2: Replace error handling incrementally
func ProcessFileWithMigration(path string) error {
    if err := legacyProcessFile(path); err != nil {
        return WrapLegacyError(err, "ProcessFile")
    }
    return nil
}
```

### Pattern 3: CLI Command Migration

```go
// Step 1: Wrap legacy commands
func WrapLegacyCommand(name string, legacyHandler func([]string)) *cobra.Command {
    return &cobra.Command{
        Use:   name,
        Short: fmt.Sprintf("Legacy %s command", name),
        Run: func(cmd *cobra.Command, args []string) {
            legacyHandler(args)
        },
    }
}

// Step 2: Gradual replacement with new framework
func NewMigratedCommand(name string) *cobra.Command {
    return cli.NewCommandBuilder().
        NewCommand(name, "Migrated command", "New implementation").
        WithHandler(newImplementation)
}
```

## ⚠️ Common Pitfalls and Solutions

### Pitfall 1: Circular Dependencies

**❌ Problem:**
```go
// Package A imports Package B
import "myapp/pkg/config"

// Package B imports Package A
import "myapp/pkg/formatter" // which depends on config
```

**✅ Solution:**
```go
// Use dependency injection and interfaces
type FormatterInterface interface {
    Format(template string, data interface{}) string
}

// Inject dependencies at application level
func NewApp() *App {
    cfg := config.NewLoader()
    fmt := formatter.NewFormatter(cfg)
    return &App{config: cfg, formatter: fmt}
}
```

### Pitfall 2: Configuration Coupling

**❌ Problem:**
```go
// Tight coupling to specific configuration schema
func ProcessFiles(cfg *BackupConfig) error {
    // Only works with backup-specific configuration
}
```

**✅ Solution:**
```go
// Use interfaces for configuration contracts
type ProcessingConfig interface {
    GetWorkDir() string
    GetMaxWorkers() int
}

func ProcessFiles(cfg ProcessingConfig) error {
    // Works with any configuration implementing the interface
}
```

### Pitfall 3: State Management

**❌ Problem:**
```go
// Global state shared across components
var globalFormatter *Formatter
var globalConfig *Config
```

**✅ Solution:**
```go
// Dependency injection with proper initialization
type App struct {
    config    config.LoaderInterface
    formatter formatter.FormatterInterface
    cli       cli.CLIApp
}

func (a *App) Initialize() error {
    // Initialize components with dependencies
    a.formatter = formatter.NewFormatter(a.config)
    a.cli = cli.NewCLIApp(cli.AppInfo{...})
    return nil
}
```

## 📈 Success Metrics

### Technical Metrics

**📊 Code Quality:**
- **Lines of Code**: Expect 30-50% reduction in main files
- **Cyclomatic Complexity**: Reduced complexity per function
- **Test Coverage**: Maintain or improve coverage >90%
- **Linting**: Zero linting violations

**🚀 Performance:**
- **Startup Time**: <5% degradation acceptable
- **Memory Usage**: Monitor for memory leaks
- **CPU Usage**: Ensure no performance regression
- **File I/O**: Maintain throughput for file operations

### Operational Metrics

**⚡ Deployment:**
- **Build Time**: Should improve with better caching
- **Binary Size**: May increase slightly due to dependencies
- **Deploy Success Rate**: Target 100% successful deployments
- **Rollback Time**: <5 minutes for emergency rollback

**[ACTION:validation] Reliability:**
- **Error Rate**: No increase in error rates
- **Uptime**: Maintain current SLA levels
- **Resource Usage**: Monitor memory and CPU patterns
- **Log Quality**: Improved structured logging

## 📚 Additional Resources

### Documentation References
- [Package Interdependency Mapping](package-interdependency-mapping.md)
- [CLI Template Application](../cmd/cli-template/README.md)
- [Extracted Package Documentation](../pkg/README.md)

### Migration Tools
- [Scaffolding Generator](../cmd/scaffolding/README.md) - Generate new projects
- [Token Suggester](../cmd/token-suggester/README.md) - Implementation tokens
- [Real-time Validator](../cmd/realtime-validator/README.md) - Validation during migration

### Best Practices
- **Start Small**: Begin with least complex components
- **Test Thoroughly**: Validate each migration step
- **Document Changes**: Keep detailed migration logs
- **Monitor Continuously**: Watch metrics during and after migration
- **Plan Rollback**: Always have a rollback strategy ready

## 🎯 Migration Timeline Template

### Week 1: Preparation
- [ ] Complete application assessment
- [ ] Choose migration strategy
- [ ] Setup new project structure
- [ ] Create migration plan

### Week 2-3: Core Components
- [ ] Migrate configuration system
- [ ] Migrate error handling
- [ ] Migrate output formatting
- [ ] Create integration tests

### Week 4: Advanced Components
- [ ] Migrate file operations
- [ ] Migrate CLI framework
- [ ] Integrate with Git operations
- [ ] Add concurrent processing

### Week 5: Testing and Validation
- [ ] Complete integration testing
- [ ] Performance benchmarking
- [ ] User acceptance testing
- [ ] Documentation updates

### Week 6: Deployment
- [ ] Gradual rollout
- [ ] Monitor metrics
- [ ] Address any issues
- [ ] Complete migration

---

**[ACTION:core-functionality] EXTRACT-008: Migration guide for package-based development patterns** - This comprehensive guide enables successful migration from monolithic to modular architecture using the extracted packages. 