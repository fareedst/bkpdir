# Configuration Schema Abstraction Analysis (REFACTOR-003)

**Task**: Configuration Schema Abstraction - HIGH PRIORITY  
**Status**: In Progress  
**Implementation Tokens**: `// REFACTOR-003: Config abstraction`, `// REFACTOR-003: Schema separation`

## Executive Summary

The current configuration system is tightly coupled to the backup application schema, making it impossible to reuse for other applications. This analysis identifies the coupling issues and provides a comprehensive abstraction strategy to enable schema-agnostic configuration management.

## Current Configuration Coupling Analysis

### 1. **Monolithic Config Struct** (Lines 35-165 in config.go)
The `Config` struct contains 80+ fields specific to backup operations:
- **Archive-specific**: `ArchiveDirPath`, `UseCurrentDirName`, `ExcludePatterns`
- **Backup-specific**: `BackupDirPath`, `UseCurrentDirNameForFiles`
- **Status codes**: 13 backup-specific status codes
- **Format strings**: 24 printf-style format strings for backup operations
- **Templates**: 24 template strings for backup operations
- **Patterns**: 4 regex patterns for backup filename parsing

**Problem**: Any application using this configuration system must define all backup-specific fields, even if irrelevant.

### 2. **Direct Config Injection** (80+ function signatures)
Functions throughout the codebase directly accept `*Config`:
```go
func CreateArchiveWithContext(ctx context.Context, cfg *Config, note string, dryRun bool, verify bool) error
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error
func NewOutputFormatter(cfg *Config) *OutputFormatter
func HandleArchiveError(err error, cfg *Config, formatter *OutputFormatter) int
```

**Problem**: Functions are tightly coupled to the specific Config struct schema.

### 3. **Field Access Patterns** (200+ direct field accesses)
Code directly accesses Config fields:
```go
f.cfg.FormatCreatedArchive    // Formatter accessing format strings
cfg.ArchiveDirPath           // Archive logic accessing paths
cfg.StatusCreatedArchive     // Error handling accessing status codes
```

**Problem**: No abstraction layer between business logic and configuration schema.

### 4. **Configuration Loading Coupling** (Lines 383-425 in config.go)
The `LoadConfig` function is hardcoded to load the backup-specific schema:
```go
func LoadConfig(root string) (*Config, error) {
    cfg := DefaultConfig()  // Returns backup-specific defaults
    // Hardcoded YAML unmarshaling to Config struct
}
```

**Problem**: Configuration loading is schema-specific and non-reusable.

## Configuration Abstraction Strategy

### Phase 1: Interface Abstraction Layer

#### 1.1 **ConfigLoader Interface**
```go
// REFACTOR-003: Config abstraction - Generic configuration loading
type ConfigLoader interface {
    LoadConfig(root string) (interface{}, error)
    LoadConfigWithSchema(root string, schema interface{}) error
    GetConfigSearchPaths() []string
    ExpandPath(path string) string
}
```

#### 1.2 **ConfigValidator Interface**
```go
// REFACTOR-003: Config abstraction - Pluggable validation
type ConfigValidator interface {
    ValidateConfig(config interface{}) error
    ValidateField(fieldName string, value interface{}) error
    GetValidationRules() map[string]ValidationRule
}

type ValidationRule struct {
    Required bool
    Type     string
    Validator func(interface{}) error
}
```

#### 1.3 **ConfigSource Interface**
```go
// REFACTOR-003: Schema separation - Source abstraction
type ConfigSource interface {
    Load(path string) (map[string]interface{}, error)
    Save(path string, data map[string]interface{}) error
    Exists(path string) bool
    GetSourceType() string // "file", "env", "default"
}
```

#### 1.4 **ConfigMerger Interface**
```go
// REFACTOR-003: Config abstraction - Generic merging
type ConfigMerger interface {
    MergeConfigs(dst, src interface{}) error
    MergeWithPriority(configs []interface{}, priorities []int) (interface{}, error)
    GetMergeStrategy() MergeStrategy
}

type MergeStrategy int
const (
    OverwriteStrategy MergeStrategy = iota
    PreserveStrategy
    AppendStrategy
)
```

### Phase 2: Schema Abstraction Layer

#### 2.1 **ConfigSchema Interface**
```go
// REFACTOR-003: Schema separation - Schema definition
type ConfigSchema interface {
    GetSchemaName() string
    GetDefaultConfig() interface{}
    GetFieldDefinitions() map[string]FieldDefinition
    ValidateSchema(config interface{}) error
    MigrateSchema(oldVersion, newVersion string, config interface{}) (interface{}, error)
}

type FieldDefinition struct {
    Name         string
    Type         string
    Required     bool
    DefaultValue interface{}
    Description  string
    Validator    func(interface{}) error
}
```

#### 2.2 **Application-Specific Schema Implementation**
```go
// REFACTOR-003: Schema separation - Backup application schema
type BackupConfigSchema struct{}

func (s *BackupConfigSchema) GetSchemaName() string {
    return "backup-application-v1"
}

func (s *BackupConfigSchema) GetDefaultConfig() interface{} {
    return DefaultConfig() // Current backup-specific defaults
}

func (s *BackupConfigSchema) GetFieldDefinitions() map[string]FieldDefinition {
    return map[string]FieldDefinition{
        "archive_dir_path": {
            Name:         "archive_dir_path",
            Type:         "string",
            Required:     true,
            DefaultValue: "../.bkpdir",
            Description:  "Directory path for storing archives",
            Validator:    validateDirectoryPath,
        },
        // ... other backup-specific fields
    }
}
```

### Phase 3: Configuration Provider Abstraction

#### 3.1 **ConfigProvider Interface**
```go
// REFACTOR-003: Config abstraction - Provider pattern
type ConfigProvider interface {
    GetString(key string) string
    GetBool(key string) bool
    GetInt(key string) int
    GetStringSlice(key string) []string
    GetDefault(key string) interface{}
    HasKey(key string) bool
    GetAllKeys() []string
}
```

#### 3.2 **Specialized Provider Interfaces**
```go
// REFACTOR-003: Config abstraction - Specialized providers
type FormatProvider interface {
    GetFormatString(operation string) string
    GetTemplateString(operation string) string
    GetPattern(patternType string) string
}

type StatusProvider interface {
    GetStatusCode(operation string) int
    GetErrorStatusCode(errorType string) int
}

type PathProvider interface {
    GetArchivePath() string
    GetBackupPath() string
    GetConfigPath() string
}
```

### Phase 4: Backward Compatibility Layer

#### 4.1 **Config Adapter Pattern**
```go
// REFACTOR-003: Config abstraction - Backward compatibility
type ConfigAdapter struct {
    provider ConfigProvider
    schema   ConfigSchema
}

func (a *ConfigAdapter) ToLegacyConfig() *Config {
    // Convert from generic provider to legacy Config struct
    return &Config{
        ArchiveDirPath:     a.provider.GetString("archive_dir_path"),
        UseCurrentDirName:  a.provider.GetBool("use_current_dir_name"),
        ExcludePatterns:    a.provider.GetStringSlice("exclude_patterns"),
        // ... map all fields
    }
}

func NewConfigAdapter(provider ConfigProvider, schema ConfigSchema) *ConfigAdapter {
    return &ConfigAdapter{
        provider: provider,
        schema:   schema,
    }
}
```

#### 4.2 **Legacy Function Wrappers**
```go
// REFACTOR-003: Config abstraction - Legacy compatibility
func LoadConfigLegacy(root string) (*Config, error) {
    loader := NewGenericConfigLoader()
    schema := &BackupConfigSchema{}
    provider, err := loader.LoadConfigWithSchema(root, schema)
    if err != nil {
        return nil, err
    }
    
    adapter := NewConfigAdapter(provider, schema)
    return adapter.ToLegacyConfig(), nil
}
```

## Implementation Plan

### Step 1: Create Interface Definitions (Week 1)
- [ ] Define all configuration interfaces in `config_interfaces.go`
- [ ] Create schema abstraction interfaces
- [ ] Define provider pattern interfaces
- [ ] Add implementation tokens throughout

### Step 2: Implement Generic Configuration System (Week 2)
- [ ] Implement `GenericConfigLoader` with schema support
- [ ] Create `YAMLConfigSource` and `EnvConfigSource` implementations
- [ ] Implement `GenericConfigMerger` with strategy support
- [ ] Create `GenericConfigValidator` with rule-based validation

### Step 3: Create Backup Application Schema (Week 3)
- [ ] Implement `BackupConfigSchema` with current field definitions
- [ ] Create `BackupConfigValidator` with backup-specific rules
- [ ] Implement specialized providers (`BackupFormatProvider`, `BackupStatusProvider`)
- [ ] Add schema migration support for future changes

### Step 4: Implement Backward Compatibility (Week 4)
- [ ] Create `ConfigAdapter` for legacy Config struct conversion
- [ ] Implement wrapper functions for existing APIs
- [ ] Add deprecation warnings for direct Config usage
- [ ] Create migration guide for consumers

### Step 5: Gradual Migration (Week 5-6)
- [ ] Update formatter to use `FormatProvider` interface
- [ ] Update error handling to use `StatusProvider` interface
- [ ] Update archive/backup logic to use `PathProvider` interface
- [ ] Maintain backward compatibility throughout

## Expected Outcomes

### 1. **Schema-Agnostic Configuration Loading**
```go
// Generic usage for any application
loader := NewGenericConfigLoader()
schema := &MyAppConfigSchema{}
config, err := loader.LoadConfigWithSchema(".", schema)

// Backup application usage (backward compatible)
config, err := LoadConfig(".") // Still works
```

### 2. **Pluggable Validation System**
```go
// Custom validation rules for different applications
validator := NewGenericConfigValidator()
validator.AddRule("custom_field", ValidationRule{
    Required: true,
    Type:     "string",
    Validator: func(v interface{}) error {
        // Custom validation logic
        return nil
    },
})
```

### 3. **Reusable Configuration Merging Logic**
```go
// Generic merging for any configuration schema
merger := NewGenericConfigMerger(OverwriteStrategy)
result, err := merger.MergeConfigs(defaultConfig, userConfig)
```

### 4. **Source-Independent Configuration Management**
```go
// Support multiple configuration sources
sources := []ConfigSource{
    NewYAMLConfigSource(),
    NewEnvConfigSource(),
    NewDefaultConfigSource(),
}
loader := NewGenericConfigLoader(sources...)
```

## Validation Criteria

### Technical Validation
- [ ] All existing tests pass with new abstraction layer
- [ ] Zero breaking changes to existing APIs during transition
- [ ] New configuration system supports multiple schemas
- [ ] Performance impact < 5% for configuration loading

### Functional Validation
- [ ] Backup application works identically with new system
- [ ] Configuration validation works for custom schemas
- [ ] Configuration merging works across different sources
- [ ] Schema migration works for version changes

### Extraction Readiness Validation
- [ ] Configuration interfaces are ready for extraction
- [ ] Zero circular dependencies in configuration system
- [ ] Clean separation between generic and application-specific code
- [ ] Backward compatibility preserved for smooth transition

## Implementation Tokens

Throughout the implementation, the following tokens will be used:

- `// REFACTOR-003: Config abstraction` - Generic configuration interfaces
- `// REFACTOR-003: Schema separation` - Schema-specific implementations
- `// REFACTOR-003: Provider pattern` - Configuration provider implementations
- `// REFACTOR-003: Backward compatibility` - Legacy compatibility layer
- `// REFACTOR-003: Migration support` - Schema migration functionality

## Dependencies

- **REFACTOR-001**: Dependency analysis must identify config coupling ✅
- **REFACTOR-002**: Large file decomposition must be complete ✅

## Blocking

- **EXTRACT-001**: Configuration Management System extraction

This abstraction will enable the configuration system to be extracted as a reusable component while maintaining full backward compatibility with the existing backup application. 