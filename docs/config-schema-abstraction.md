# Configuration Schema Abstraction Design

## [LOW] REFACTOR-003: Configuration Schema Abstraction Implementation

This document describes the configuration schema abstraction design implemented to prepare the configuration system for extraction while maintaining backward compatibility with the existing backup application.

### 📑 Purpose

The configuration schema abstraction enables:
- **Schema-agnostic configuration loading** - Abstract configuration loading from specific backup application schema
- **Pluggable validation system** - Allow different applications to define their own schemas  
- **Reusable configuration merging logic** - Enable generic configuration merging across different schemas
- **Source-independent configuration management** - Abstract file, environment, and default sources
- **Configuration extraction preparation** - Prepare configuration system for clean extraction to `pkg/config`

### 🏗️ Architecture Overview

The abstraction introduces several layers to separate concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  (Backup Application - Uses ApplicationConfig interface)    │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Interface Layer                            │
│  ConfigLoader │ ConfigMerger │ ConfigValidator │ ConfigSource │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                Implementation Layer                         │
│  DefaultConfigLoader │ BackupAppValidator │ FileConfigSource │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   Core Layer                                │
│        (Existing Config struct and functions)               │
└─────────────────────────────────────────────────────────────┘
```

### [ACTION:core-functionality] Interface Definitions

#### ConfigLoader Interface
**Purpose**: Schema-agnostic configuration management operations
**File**: `config_interfaces.go`

```go
type ConfigLoader interface {
    LoadConfig(root string) (*Config, error)
    LoadConfigValues(root string) (map[string]ConfigValue, error)
    GetConfigValues(cfg *Config) []ConfigValue
    GetConfigValuesWithSources(cfg *Config, root string) []ConfigValue
    ValidateConfig(cfg *Config) error
}
```

#### ConfigMerger Interface  
**Purpose**: Schema-agnostic configuration merging and composition
**File**: `config_interfaces.go`

```go
type ConfigMerger interface {
    MergeConfigs(dst, src *Config)
    MergeConfigValues(dst, src map[string]ConfigValue)
    GetConfigSearchPaths() []string
    ExpandPath(path string) string
}
```

#### ConfigSource Interface
**Purpose**: Abstract different configuration sources (file, environment, defaults)
**File**: `config_interfaces.go`

```go
type ConfigSource interface {
    LoadFromFile(path string) (*Config, error)
    LoadFromEnvironment() (*Config, error)  
    LoadDefaults() *Config
    GetSourceName() string
    IsAvailable() bool
}
```

#### ConfigValidator Interface
**Purpose**: Pluggable configuration validation for different application schemas
**File**: `config_interfaces.go`

```go
type ConfigValidator interface {
    ValidateSchema(cfg *Config) error
    ValidateValues(values map[string]ConfigValue) error
    GetRequiredFields() []string
    GetValidationRules() map[string]ValidationRule
}
```

#### ApplicationConfig Interface
**Purpose**: Abstract backup-specific schema from generic configuration operations
**File**: `config_interfaces.go`

```go
type ApplicationConfig interface {
    GetArchiveSettings() ArchiveSettings
    GetBackupSettings() BackupSettings
    GetStatusCodes() map[string]int
    GetFormatSettings() FormatSettings
}
```

### [ACTION:core-functionality] Concrete Implementations

#### DefaultConfigLoader
**Purpose**: Default implementation maintaining backward compatibility
**File**: `config_impl.go`

- Uses dependency injection for file operations, merging, and validation
- Maintains existing configuration loading behavior
- Provides schema-agnostic interface for future extraction

#### BackupAppValidator
**Purpose**: Backup application-specific validation
**File**: `config_impl.go`

- Validates backup application schema requirements
- Demonstrates schema-specific validation within common interface
- Provides validation rules and required fields for backup application

#### BackupApplicationConfig
**Purpose**: Backup application configuration wrapper
**File**: `config_impl.go`

- Wraps existing Config struct with application-specific interface
- Separates archive, backup, status code, and format concerns
- Enables clean separation of backup-specific logic from generic configuration

#### FileSystemOperations
**Purpose**: File system operations abstraction
**File**: `config_impl.go`

- Abstracts file operations to enable testing and different storage backends
- Provides clean interface for configuration file management

### [ACTION:format-processing] Schema Separation Structures

#### ArchiveSettings
**Purpose**: Archive-specific configuration separation
```go
type ArchiveSettings struct {
    DirectoryPath      string
    UseCurrentDirName  bool
    ExcludePatterns    []string
    IncludeGitInfo     bool
    ShowGitDirtyStatus bool
    Verification       *VerificationConfig
}
```

#### BackupSettings  
**Purpose**: Backup-specific configuration separation
```go
type BackupSettings struct {
    DirectoryPath             string
    UseCurrentDirNameForFiles bool
}
```

#### FormatSettings
**Purpose**: Format configuration separation
```go
type FormatSettings struct {
    FormatStrings      map[string]string
    TemplateStrings    map[string]string
    PatternStrings     map[string]string
    ErrorFormatStrings map[string]string
}
```

### [ACTION:migration] Backward Compatibility

The abstraction maintains complete backward compatibility:

1. **Existing Config struct unchanged** - All existing fields and methods preserved
2. **Existing functions still work** - `LoadConfig()`, `GetConfigValues()`, etc. continue to function
3. **No breaking changes** - All existing code continues to work without modification
4. **Gradual migration** - New interfaces can be adopted incrementally

### 🚀 Extraction Preparation

The abstraction prepares for extraction to `pkg/config` by:

#### Clean Interface Boundaries
- All interfaces defined independently of backup application logic
- Clear separation between generic and application-specific concerns
- No circular dependencies between components

#### Dependency Injection
- File operations abstracted through `ConfigFileOperations` interface
- Validation logic separated through `ConfigValidator` interface  
- Merging logic abstracted through `ConfigMerger` interface

#### Schema Independence
- Configuration loading logic separated from backup application schema
- Validation rules externalized and configurable
- Source management abstracted from specific configuration types

### 📋 Implementation Tokens

All code changes include implementation tokens for traceability:

- `// [LOW] REFACTOR-003: Config abstraction - [Description] - [ACTION:core-functionality]`
- `// [LOW] REFACTOR-003: Schema separation - [Description] - [ACTION:format-processing]`

These tokens enable tracking during the extraction process.

### 🧪 Testing Strategy

The abstraction enables enhanced testing through:

#### Interface Mocking
- `ConfigFileOperations` can be mocked for file system testing
- `ConfigValidator` can be mocked for validation testing
- `ConfigSource` can be mocked for source testing

#### Schema Testing
- Different schemas can be tested through `ConfigValidator` implementations
- Validation rules can be tested independently
- Configuration merging can be tested with different schemas

#### Integration Testing
- Full configuration loading can be tested end-to-end
- Source priority can be tested with multiple sources
- Error handling can be tested at each layer

### 📈 Future Extraction Plan

#### Phase 1: Interface Stabilization (Complete)
- ✅ Define all configuration interfaces
- ✅ Create concrete implementations
- ✅ Maintain backward compatibility

#### Phase 2: Code Migration (EXTRACT-001)
- Move interfaces to `pkg/config/interfaces.go`
- Move implementations to `pkg/config/loader.go`, `pkg/config/validator.go`, etc.
- Update imports in existing code

#### Phase 3: Schema Generalization
- Create generic schema definition framework
- Enable registration of different application schemas
- Support schema migration and versioning

#### Phase 4: Advanced Features
- Add configuration caching and performance optimization
- Implement configuration change notifications
- Add configuration templating and preprocessing

### ⚠️ Dependencies and Constraints

#### Dependencies Satisfied
- ✅ **REFACTOR-001**: Dependency analysis confirmed clean configuration boundaries
- ✅ **Immutable Requirements**: No conflicts with configuration defaults or discovery

#### Constraints Maintained
- Configuration discovery via `BKPDIR_CONFIG` environment variable preserved
- Default search path `./.bkpdir.yml:~/.bkpdir.yml` maintained
- YAML configuration format preserved
- Configuration merging behavior unchanged

### 🎯 Success Criteria

#### Technical Criteria
- ✅ **Schema Independence**: Configuration loading abstracted from backup application schema
- ✅ **Pluggable Validation**: Multiple validation implementations possible
- ✅ **Source Abstraction**: File, environment, and default sources abstracted
- ✅ **Interface Contracts**: Clear interfaces defined for all components

#### Quality Criteria
- ✅ **Backward Compatibility**: All existing functionality preserved
- ✅ **No Breaking Changes**: Existing code works without modification
- ✅ **Clean Boundaries**: Clear separation between generic and application-specific logic
- ✅ **Testability**: Enhanced testing through interface mocking

#### Extraction Readiness Criteria
- ✅ **Interface Definitions**: All required interfaces defined
- ✅ **Concrete Implementations**: Working implementations created
- ✅ **Dependency Injection**: External dependencies abstracted
- ✅ **Documentation**: Complete design and implementation documentation

### 📊 Implementation Status

**Status**: ✅ **COMPLETED**

- **Interface Definitions**: ✅ Complete (`config_interfaces.go`)
- **Concrete Implementations**: ✅ Complete (`config_impl.go`)
- **Schema Separation**: ✅ Complete (ArchiveSettings, BackupSettings, FormatSettings)
- **Backward Compatibility**: ✅ Maintained
- **Documentation**: ✅ Complete (this document)

The configuration schema abstraction is ready for extraction and provides a solid foundation for reusable configuration management across different CLI applications. 