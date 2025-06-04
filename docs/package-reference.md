# 📖 BkpDir Package Reference

> **🔺 EXTRACT-008: Package reference documentation - 📖 API and configuration reference**

## 🎯 Overview

This document provides comprehensive API reference documentation for all extracted BkpDir packages. Each package is designed to be used independently or in combination with others.

## 📦 Package Index

| Package | Purpose | Key Types | Main Functions |
|---------|---------|-----------|----------------|
| [pkg/config](#pkgconfig) | Configuration management | `Config`, `Loader` | `Load()`, `Save()`, `Validate()` |
| [pkg/errors](#pkgerrors) | Error handling | `Error`, `Collector` | `Wrap()`, `New()`, `WithField()` |
| [pkg/resources](#pkgresources) | Resource management | `Manager`, `Resource` | `Acquire()`, `Release()`, `Cleanup()` |
| [pkg/formatter](#pkgformatter) | Output formatting | `Formatter`, `Template` | `Print()`, `PrintProgress()`, `Format()` |
| [pkg/git](#pkggit) | Git integration | `Repository`, `Info` | `Discover()`, `Info()`, `FileStatus()` |
| [pkg/cli](#pkgcli) | CLI framework | `App`, `Command` | `NewApp()`, `AddCommand()`, `Run()` |
| [pkg/fileops](#pkgfileops) | File operations | `FileInfo`, `Comparer` | `Compare()`, `SafeCopy()`, `IsAccessible()` |
| [pkg/processing](#pkgprocessing) | Concurrent processing | `Pool`, `Worker` | `NewPool()`, `Process()`, `Close()` |

---

## pkg/config

### 🎯 Purpose
Schema-agnostic configuration loading with support for multiple formats and environment variable overrides.

### 📋 Core Types

#### Config
```go
type Config struct {
    OutputFormat string        `json:"output_format" yaml:"output_format"`
    MaxWorkers   int          `json:"max_workers" yaml:"max_workers"`
    Timeout      time.Duration `json:"timeout" yaml:"timeout"`
    DatabaseURL  string        `json:"database_url" yaml:"database_url"`
    HTTPTimeout  int          `json:"http_timeout" yaml:"http_timeout"`
    SyncURL      string        `json:"sync_url" yaml:"sync_url"`
}
```

#### Loader
```go
type Loader interface {
    Load(path string) (*Config, error)
    LoadWithDefaults(defaults map[string]interface{}) (*Config, error)
    Save(config *Config, path string) error
}
```

### 🔧 Core Functions

#### Load Configuration
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Config API examples
func Load() (*Config, error)
func LoadFrom(path string) (*Config, error)
func LoadWithDefaults(defaults map[string]interface{}) (*Config, error)

// Example usage
cfg, err := config.Load()
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// Load with custom defaults
cfg, err := config.LoadWithDefaults(map[string]interface{}{
    "max_workers": runtime.NumCPU(),
    "output_format": "json",
    "timeout": "30s",
})
```

#### Save Configuration
```go
func Save(cfg *Config, path string) error
func SaveAs(cfg *Config, path string, format string) error

// Example usage
err := config.Save(cfg, "config.yaml")
if err != nil {
    return fmt.Errorf("failed to save config: %w", err)
}
```

#### Validation
```go
func Validate(cfg *Config) error
func ValidateField(cfg *Config, field string) error

// Example usage
if err := config.Validate(cfg); err != nil {
    return fmt.Errorf("invalid configuration: %w", err)
}
```

### ⚙️ Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `output_format` | string | `"text"` | Output format: text, json, yaml |
| `max_workers` | int | `4` | Maximum concurrent workers |
| `timeout` | duration | `"30s"` | Operation timeout |
| `database_url` | string | `""` | Database connection URL |
| `http_timeout` | int | `30` | HTTP client timeout in seconds |
| `sync_url` | string | `""` | Remote sync service URL |

### 📁 Configuration File Locations

1. `./config.yaml` (current directory)
2. `~/.config/bkpdir/config.yaml` (user config)
3. `/etc/bkpdir/config.yaml` (system config)
4. Environment variables with `BKPDIR_` prefix

### 🌍 Environment Variables

```bash
BKPDIR_OUTPUT_FORMAT=json
BKPDIR_MAX_WORKERS=8
BKPDIR_TIMEOUT=60s
BKPDIR_DATABASE_URL=postgres://user:pass@localhost/db
```

---

## pkg/errors

### 🎯 Purpose
Structured error handling with context preservation and error chaining.

### 📋 Core Types

#### Error
```go
type Error interface {
    error
    WithField(key string, value interface{}) Error
    WithFields(fields map[string]interface{}) Error
    Fields() map[string]interface{}
    Unwrap() error
}
```

#### Collector
```go
type Collector interface {
    Add(err error)
    HasErrors() bool
    Error() error
    Errors() []error
    Count() int
}
```

### 🔧 Core Functions

#### Error Creation
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Error handling patterns
func New(message string) Error
func Wrap(err error, message string) Error
func Wrapf(err error, format string, args ...interface{}) Error

// Example usage
err := errors.New("operation failed").
    WithField("operation", "backup").
    WithField("file", "/path/to/file")

wrappedErr := errors.Wrap(originalErr, "backup operation failed").
    WithField("timestamp", time.Now())
```

#### Error Checking
```go
func Is(err, target error) bool
func As(err error, target interface{}) bool
func IsRecoverable(err error) bool

// Example usage
if errors.Is(err, os.ErrNotExist) {
    // Handle file not found
}

var pathErr *os.PathError
if errors.As(err, &pathErr) {
    // Handle path-specific error
}
```

#### Error Collection
```go
func NewCollector() Collector

// Example usage
collector := errors.NewCollector()
for _, file := range files {
    if err := processFile(file); err != nil {
        collector.Add(errors.Wrap(err, "processing failed").
            WithField("file", file))
    }
}

if collector.HasErrors() {
    return collector.Error()
}
```

### 🏷️ Error Categories

| Category | Description | Example |
|----------|-------------|---------|
| **System** | OS-level errors | File I/O, network, permissions |
| **User** | User input errors | Invalid arguments, missing files |
| **Application** | Business logic errors | Validation failures, state errors |
| **Recoverable** | Temporary failures | Network timeouts, resource busy |

---

## pkg/resources

### 🎯 Purpose
Resource lifecycle management with automatic cleanup and coordination.

### 📋 Core Types

#### Resource
```go
type Resource interface {
    ID() string
    Type() string
    Release() error
    IsReleased() bool
}
```

#### Manager
```go
type Manager interface {
    Acquire(resourceType string) (Resource, error)
    AcquireWithTimeout(resourceType string, timeout time.Duration) (Resource, error)
    Release(resource Resource) error
    Cleanup() error
    Usage() map[string]int64
}
```

### 🔧 Core Functions

#### Resource Management
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Resource management patterns
func Acquire(resourceType string) (Resource, error)
func AcquireWithTimeout(resourceType string, timeout time.Duration) (Resource, error)
func Release(resource Resource) error

// Example usage
resource, err := resources.Acquire("file_handle")
if err != nil {
    return fmt.Errorf("failed to acquire resource: %w", err)
}
defer resource.Release()

// Use resource
err = performOperation(resource)
```

#### Manager Operations
```go
func NewManager() Manager
func (m Manager) Cleanup() error
func (m Manager) Usage() map[string]int64

// Example usage
manager := resources.NewManager()
defer manager.Cleanup()

resource, err := manager.Acquire("processing_slot")
if err != nil {
    return err
}

// Resource automatically cleaned up on defer
```

### 🏷️ Resource Types

| Type | Description | Limit | Timeout |
|------|-------------|-------|---------|
| `file_handle` | File system handles | 100 | 30s |
| `processing_slot` | Processing workers | CPU count | 60s |
| `network_connection` | Network connections | 50 | 30s |
| `memory_buffer` | Memory allocations | 1GB | 10s |

---

## pkg/formatter

### 🎯 Purpose
Consistent output formatting with template support and multiple output formats.

### 📋 Core Types

#### Formatter
```go
type Formatter interface {
    Print(data interface{}) error
    PrintProgress(progress interface{}) error
    PrintInfo(message string) error
    PrintWarning(message string) error
    PrintError(message string) error
    Format(template string, data interface{}) (string, error)
}
```

#### Template
```go
type Template interface {
    Execute(data interface{}) (string, error)
    ExecuteToWriter(w io.Writer, data interface{}) error
}
```

### 🔧 Core Functions

#### Formatter Creation
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Formatter API examples
func New(format string) Formatter
func NewWithOptions(format string, options Options) Formatter

// Example usage
formatter := formatter.New("json")
err := formatter.Print(map[string]interface{}{
    "status": "success",
    "files": 42,
    "duration": "1.5s",
})
```

#### Output Methods
```go
func (f Formatter) Print(data interface{}) error
func (f Formatter) PrintProgress(progress interface{}) error
func (f Formatter) PrintInfo(message string) error

// Example usage
// JSON output
formatter.Print(map[string]interface{}{
    "operation": "backup",
    "files": []string{"file1.txt", "file2.txt"},
    "timestamp": time.Now(),
})

// Progress output
formatter.PrintProgress(map[string]interface{}{
    "completed": 75,
    "total": 100,
    "percent": 75.0,
})
```

#### Template Operations
```go
func NewTemplate(templateStr string) (Template, error)
func (t Template) Execute(data interface{}) (string, error)

// Example usage
tmpl, err := formatter.NewTemplate("Processing {{.Count}} files...")
if err != nil {
    return err
}

output, err := tmpl.Execute(map[string]interface{}{
    "Count": 42,
})
```

### 🎨 Output Formats

| Format | Description | Example |
|--------|-------------|---------|
| `text` | Human-readable text | `Processing 42 files...` |
| `json` | JSON format | `{"status":"success","files":42}` |
| `yaml` | YAML format | `status: success\nfiles: 42` |
| `table` | Tabular format | ASCII table output |
| `csv` | CSV format | `status,files\nsuccess,42` |

---

## pkg/git

### 🎯 Purpose
Git repository integration with branch detection, commit information, and file status.

### 📋 Core Types

#### Repository
```go
type Repository interface {
    Path() string
    Info() (*Info, error)
    FileStatus(path string) (*FileStatus, error)
    LastCommit(path string) (*Commit, error)
    IsClean() (bool, error)
}
```

#### Info
```go
type Info struct {
    Repository string `json:"repository"`
    Branch     string `json:"branch"`
    Commit     string `json:"commit"`
    IsDirty    bool   `json:"is_dirty"`
    Remote     string `json:"remote,omitempty"`
}
```

#### FileStatus
```go
type FileStatus struct {
    Path       string `json:"path"`
    Status     string `json:"status"`
    IsModified bool   `json:"is_modified"`
    IsStaged   bool   `json:"is_staged"`
    IsUntracked bool  `json:"is_untracked"`
}
```

### 🔧 Core Functions

#### Repository Discovery
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Git API examples
func Discover(path string) (Repository, error)
func Init(path string) (Repository, error)
func Clone(url, path string) (Repository, error)

// Example usage
repo, err := git.Discover(".")
if err != nil {
    // Not in a git repository
    return nil
}

info, err := repo.Info()
if err != nil {
    return fmt.Errorf("failed to get git info: %w", err)
}
```

#### Repository Information
```go
func (r Repository) Info() (*Info, error)
func (r Repository) IsClean() (bool, error)
func (r Repository) Branch() (string, error)

// Example usage
info, err := repo.Info()
if err != nil {
    return err
}

fmt.Printf("Branch: %s\n", info.Branch)
fmt.Printf("Commit: %s\n", info.Commit[:8])
fmt.Printf("Dirty: %v\n", info.IsDirty)
```

#### File Operations
```go
func (r Repository) FileStatus(path string) (*FileStatus, error)
func (r Repository) LastCommit(path string) (*Commit, error)
func (r Repository) Add(paths ...string) error

// Example usage
status, err := repo.FileStatus("README.md")
if err != nil {
    return err
}

if status.IsModified {
    fmt.Println("File has uncommitted changes")
}
```

### 🏷️ Git Status Types

| Status | Description | IsModified | IsStaged | IsUntracked |
|--------|-------------|------------|----------|-------------|
| `clean` | No changes | false | false | false |
| `modified` | Working tree changes | true | false | false |
| `staged` | Staged changes | false | true | false |
| `untracked` | New file | false | false | true |
| `deleted` | File deleted | true | false | false |

---

## pkg/cli

### 🎯 Purpose
CLI framework with command patterns, argument handling, and context management.

### 📋 Core Types

#### App
```go
type App interface {
    AddCommand(cmd *Command) error
    Run(ctx context.Context, args []string) error
    SetGlobalFlags(flags []Flag) error
    SetConfig(config interface{}) error
}
```

#### Command
```go
type Command struct {
    Name        string
    Description string
    Usage       string
    Flags       []Flag
    Handler     HandlerFunc
    Subcommands []*Command
}
```

#### Flag
```go
type Flag struct {
    Name        string
    Value       string
    Description string
    Required    bool
    Hidden      bool
}
```

### 🔧 Core Functions

#### Application Setup
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 CLI API examples
func NewApp(name, description string) App
func (a App) AddCommand(cmd *Command) error
func (a App) Run(ctx context.Context, args []string) error

// Example usage
app := cli.NewApp("myapp", "My CLI Application")

app.AddCommand(&cli.Command{
    Name:        "backup",
    Description: "Create a backup",
    Flags: []cli.Flag{
        {Name: "output", Value: "backup.tar", Description: "Output file"},
        {Name: "compress", Value: "true", Description: "Compress output"},
    },
    Handler: backupHandler,
})

err := app.Run(context.Background(), os.Args)
```

#### Command Handlers
```go
type HandlerFunc func(ctx context.Context, args []string) error

func backupHandler(ctx context.Context, args []string) error {
    outputFile := cli.GetFlag(ctx, "output")
    compress := cli.GetBoolFlag(ctx, "compress")
    
    // Implement backup logic
    return performBackup(outputFile, compress, args)
}
```

#### Flag Operations
```go
func GetFlag(ctx context.Context, name string) string
func GetBoolFlag(ctx context.Context, name string) bool
func GetIntFlag(ctx context.Context, name string) int

// Example usage
func commandHandler(ctx context.Context, args []string) error {
    verbose := cli.GetBoolFlag(ctx, "verbose")
    workers := cli.GetIntFlag(ctx, "workers")
    format := cli.GetFlag(ctx, "format")
    
    if verbose {
        fmt.Printf("Using %d workers with %s format\n", workers, format)
    }
    
    return nil
}
```

### 🏷️ Command Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| **Simple** | Single action command | `backup`, `restore` |
| **Subcommand** | Nested commands | `git commit`, `docker run` |
| **Pipeline** | Chainable operations | `process \| filter \| output` |
| **Interactive** | User input prompts | Configuration wizard |

---

## pkg/fileops

### 🎯 Purpose
Safe file operations with atomic writes, comparisons, and error handling.

### 📋 Core Types

#### FileInfo
```go
type FileInfo struct {
    Path     string      `json:"path"`
    Size     int64       `json:"size"`
    Mode     os.FileMode `json:"mode"`
    ModTime  time.Time   `json:"mod_time"`
    IsDir    bool        `json:"is_dir"`
    Checksum string      `json:"checksum,omitempty"`
}
```

#### Comparer
```go
type Comparer interface {
    Compare(path1, path2 string) (bool, error)
    CompareContent(path1, path2 string) (bool, error)
    CompareChecksum(path1, path2 string) (bool, error)
}
```

### 🔧 Core Functions

#### File Operations
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 File operations examples
func Exists(path string) bool
func IsAccessible(path string) bool
func Stat(path string) (*FileInfo, error)
func Compare(path1, path2 string) (bool, error)

// Example usage
if !fileops.Exists("input.txt") {
    return errors.New("input file does not exist")
}

if !fileops.IsAccessible("input.txt") {
    return errors.New("input file is not accessible")
}

info, err := fileops.Stat("input.txt")
if err != nil {
    return fmt.Errorf("failed to stat file: %w", err)
}
```

#### Safe Operations
```go
func SafeCopy(src, dst string) error
func SafeMove(src, dst string) error
func SafeWrite(path string, data []byte) error
func AtomicWrite(path string, data []byte) error

// Example usage
err := fileops.SafeCopy("source.txt", "backup.txt")
if err != nil {
    return fmt.Errorf("failed to create backup: %w", err)
}

// Atomic write ensures file is complete or not created
err = fileops.AtomicWrite("output.txt", data)
if err != nil {
    return fmt.Errorf("failed to write file: %w", err)
}
```

#### Comparison Operations
```go
func Compare(path1, path2 string) (bool, error)
func CompareContent(path1, path2 string) (bool, error)
func CompareSize(path1, path2 string) (bool, error)
func Checksum(path string) (string, error)

// Example usage
same, err := fileops.Compare("file1.txt", "file2.txt")
if err != nil {
    return err
}

if same {
    fmt.Println("Files are identical")
} else {
    fmt.Println("Files differ")
}
```

### 🛡️ Safety Features

| Feature | Description | Benefit |
|---------|-------------|---------|
| **Atomic Writes** | All-or-nothing file creation | Prevents partial files |
| **Backup Creation** | Automatic backup before overwrite | Data protection |
| **Permission Checks** | Verify access before operations | Early error detection |
| **Checksum Validation** | Content verification | Data integrity |

---

## pkg/processing

### 🎯 Purpose
Concurrent processing with worker pools, progress tracking, and cancellation support.

### 📋 Core Types

#### Pool
```go
type Pool interface {
    Process(ctx context.Context, items []string, fn ProcessFunc) error
    ProcessWithCallback(ctx context.Context, items []string, fn ProcessFunc, callback CallbackFunc) error
    Close() error
    Workers() int
    SetWorkers(count int) error
}
```

#### ProcessFunc
```go
type ProcessFunc func(item string) error
```

#### CallbackFunc
```go
type CallbackFunc func(item string, err error)
```

### 🔧 Core Functions

#### Pool Management
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Processing API examples
func NewPool(workers int) Pool
func NewPoolWithResources(workers int, resources []Resource) Pool
func (p Pool) Process(ctx context.Context, items []string, fn ProcessFunc) error

// Example usage
pool := processing.NewPool(runtime.NumCPU())
defer pool.Close()

err := pool.Process(context.Background(), files, func(file string) error {
    return processFile(file)
})
```

#### Advanced Processing
```go
func (p Pool) ProcessWithCallback(ctx context.Context, items []string, fn ProcessFunc, callback CallbackFunc) error
func (p Pool) ProcessBatch(ctx context.Context, items []string, batchSize int, fn BatchProcessFunc) error

// Example usage with callback
err := pool.ProcessWithCallback(
    ctx, 
    files,
    func(file string) error {
        return processFile(file)
    },
    func(file string, err error) {
        if err != nil {
            log.Printf("Failed to process %s: %v", file, err)
        } else {
            log.Printf("Successfully processed %s", file)
        }
    },
)
```

#### Progress Monitoring
```go
func (p Pool) Progress() *Progress
func (p Pool) SetProgressCallback(callback ProgressCallback)

type Progress struct {
    Total     int64 `json:"total"`
    Completed int64 `json:"completed"`
    Failed    int64 `json:"failed"`
    Percent   float64 `json:"percent"`
}

// Example usage
pool.SetProgressCallback(func(progress *Progress) {
    fmt.Printf("Progress: %.1f%% (%d/%d)\n", 
        progress.Percent, progress.Completed, progress.Total)
})
```

### ⚡ Performance Tuning

| Parameter | Description | Recommended |
|-----------|-------------|-------------|
| **Workers** | Number of concurrent workers | CPU count for CPU-bound, 2x CPU for I/O-bound |
| **Batch Size** | Items per batch | 50-200 for I/O, 10-50 for CPU |
| **Buffer Size** | Channel buffer size | 2x worker count |
| **Timeout** | Operation timeout | 30s-5m depending on operation |

---

## 🔗 Cross-Package Integration

### Error Handling Across Packages
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Cross-package error handling
func IntegratedOperation() error {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        return errors.Wrap(err, "configuration loading failed").
            WithField("operation", "integrated_operation")
    }

    // Acquire resources
    resource, err := resources.Acquire("processing")
    if err != nil {
        return errors.Wrap(err, "resource acquisition failed").
            WithField("resource_type", "processing")
    }
    defer resource.Release()

    // Process with Git awareness
    if repo, err := git.Discover("."); err == nil {
        info, _ := repo.Info()
        return errors.New("operation context").
            WithField("git_branch", info.Branch).
            WithField("git_commit", info.Commit[:8])
    }

    return nil
}
```

### Testing Utilities
```go
// 🔺 EXTRACT-008: Package reference documentation - 📖 Testing utilities
func TestIntegration(t *testing.T) {
    // Create test configuration
    cfg := &config.Config{
        OutputFormat: "json",
        MaxWorkers:   2,
    }

    // Create test formatter
    formatter := formatter.New(cfg.OutputFormat)

    // Create test processor
    processor := processing.NewPool(cfg.MaxWorkers)
    defer processor.Close()

    // Test integration
    err := processor.Process(context.Background(), []string{"test"}, func(item string) error {
        return formatter.PrintInfo("Processing: " + item)
    })

    assert.NoError(t, err)
}
```

---

## 📚 Summary

This reference documentation covers all 8 extracted packages:

1. **pkg/config**: Schema-agnostic configuration management
2. **pkg/errors**: Structured error handling with context
3. **pkg/resources**: Resource lifecycle management
4. **pkg/formatter**: Consistent output formatting
5. **pkg/git**: Git repository integration
6. **pkg/cli**: CLI framework and command patterns
7. **pkg/fileops**: Safe file operations
8. **pkg/processing**: Concurrent processing with worker pools

Each package is designed to work independently while providing seamless integration when used together. The packages follow consistent patterns for error handling, resource management, and configuration.

---

**🔺 EXTRACT-008: Package reference documentation - 📖 Comprehensive API and configuration reference** 