# CFG-005: Layered Configuration Inheritance - Implementation Plan

## 🚀 PHASE 1: CRITICAL VALIDATION [COMPLETED ✅]

### 1.1 Immutable Requirements Check ✅
- **Status**: VERIFIED - No conflicts with immutable.md
- **Key Findings**:
  - Configuration discovery pattern preserved (BKPDIR_CONFIG, search path order)
  - Backward compatibility required - existing configs must work unchanged
  - Default values must be preserved unless explicitly overridden
  - Performance requirements must be maintained

### 1.2 Feature Tracking Registry ✅  
- **Status**: VERIFIED - CFG-005 exists in feature-tracking.md
- **Feature ID**: CFG-005 
- **Priority**: ⭐ CRITICAL
- **Implementation Tokens**: `// ⭐ CFG-005: Layered configuration inheritance`

### 1.3 AI Assistant Compliance ✅
- **Status**: VERIFIED - Token requirements understood
- **Required Format**: `// ⭐ CFG-005: Description - 🔧 Action context`
- **Icons Required**: ⭐ (Critical priority), 🔧 (Configure/modify), 🔍 (Search/discover), 📝 (Document/update)

### 1.4 AI Assistant Protocol ✅
- **Protocol**: NEW FEATURE Protocol (followed)
- **Documentation Updated**: specification.md, requirements.md, architecture.md, testing.md

## 🔧 PHASE 2: IMPLEMENTATION ANALYSIS AND PLANNING

### 2.1 Existing Configuration System Analysis

**Current Architecture (from EXTRACT-001):**
- pkg/config package exists with schema-agnostic loading
- GenericConfigLoader with reflection-based merging  
- PathDiscovery for configurable search paths
- DefaultEnvironmentProvider for env var support
- Backward compatibility adapter in place

**Integration Points:**
- config.go - main configuration struct and loading
- config_adapter.go - backward compatibility layer
- pkg/config/ - extracted generic configuration system

### 2.2 Implementation Strategy

**Phase 2A: Core Inheritance Structure**
1. Extend InheritableConfig struct in pkg/config
2. Add inheritance parsing to GenericConfigLoader
3. Implement basic dependency chain building
4. Add circular dependency detection

**Phase 2B: Merge Strategy System**  
1. Implement prefix-based key processing
2. Create merge strategy interfaces and implementations
3. Add strategy-aware configuration merging
4. Integrate with existing reflection-based merging

**Phase 2C: Integration and Testing**
1. Update backward compatibility adapter
2. Add comprehensive test suite  
3. Update configuration discovery integration
4. Performance optimization and caching

## 📋 DETAILED SUBTASK BREAKDOWN

### Subtask 1: Design inheritance configuration structure (⭐ CRITICAL)
**Status**: ✅ COMPLETED
**Files**: pkg/config/interfaces.go, pkg/config/loader.go
**Implementation**:
```go
// ConfigInheritance metadata structure
type ConfigInheritance struct {
    Inherit []string `yaml:"inherit"` // Parent configuration files
}

// InheritableConfig wrapper
type InheritableConfig struct {
    Config interface{} `yaml:",inline"`         // Main configuration
    ConfigInheritance `yaml:",inline"`          // Inheritance metadata  
}

// Inheritance chain tracking
type InheritanceChain struct {
    Files   []string           // Files in dependency order
    Visited map[string]bool    // Circular dependency prevention
    Sources map[string]string  // Source tracking for debugging
}
```

**Dependencies**: None
**Estimated Time**: 2-3 hours
**Test Requirements**: Basic structure validation

### Subtask 2: Implement prefix-based merge strategies (⭐ CRITICAL)  
**Status**: ✅ COMPLETED
**Files**: pkg/config/merge_strategies.go (new), pkg/config/loader.go
**Implementation**:
```go
// Merge strategy interface
type MergeStrategy interface {
    Merge(dest, src reflect.Value) error
    GetPrefix() string
    GetDescription() string
}

// Strategy implementations
type StandardOverrideStrategy struct{}  // No prefix
type ArrayMergeStrategy struct{}        // + prefix  
type ArrayPrependStrategy struct{}      // ^ prefix
type ArrayReplaceStrategy struct{}      // ! prefix
type DefaultValueStrategy struct{}      // = prefix

// Key prefix processor
type PrefixedKeyProcessor struct {
    prefixMap map[string]MergeStrategy
}
```

**Dependencies**: Subtask 1
**Estimated Time**: 4-5 hours  
**Test Requirements**: All merge strategies with nested data

### Subtask 3: Extend configuration loading engine (🔺 HIGH)
**Status**: 📝 NOT STARTED  
**Files**: pkg/config/loader.go, pkg/config/inheritance.go (new)
**Implementation**:
```go
// Enhanced loader with inheritance support
type InheritanceConfigLoader struct {
    baseLoader       ConfigLoader
    pathResolver     PathResolver  
    chainBuilder     InheritanceChainBuilder
    strategyEngine   MergeStrategyProcessor
    circularDetector CircularDependencyDetector
}

func (l *InheritanceConfigLoader) LoadConfigWithInheritance(root string, config interface{}) (interface{}, error) {
    // 1. Build inheritance dependency graph
    // 2. Load configurations in dependency order  
    // 3. Apply merge strategies and resolve final configuration
}
```

**Dependencies**: Subtasks 1, 2
**Estimated Time**: 6-8 hours
**Test Requirements**: Multi-level inheritance chains

### Subtask 4: Implement inheritance resolution logic (🔺 HIGH)
**Status**: ✅ COMPLETED
**Files**: pkg/config/inheritance.go, pkg/config/resolver.go (new)  
**Implementation**:
```go
// Inheritance chain builder
type InheritanceChainBuilder struct {
    pathDiscovery *PathDiscovery
    fileOps       ConfigFileOperations
}

// Circular dependency detector  
type CircularDependencyDetector struct {
    visited map[string]bool
    inStack map[string]bool
    path    []string
}

func (d *CircularDependencyDetector) DetectCycle(startFile string, resolver PathResolver) error {
    // Depth-first search for cycle detection
}
```

**Dependencies**: Subtask 3
**Estimated Time**: 4-6 hours
**Test Requirements**: Circular dependency detection, missing files

### Subtask 5: Create merge strategy processor (🔺 HIGH)
**Status**: 📝 NOT STARTED
**Files**: pkg/config/merge_strategies.go, pkg/config/processor.go (new)
**Implementation**:
```go
// Strategy processor coordination
type MergeStrategyProcessor struct {
    strategies map[string]MergeStrategy
    processor  *PrefixedKeyProcessor
}

func (p *MergeStrategyProcessor) ProcessMergeOperations(config map[string]interface{}) (ProcessedConfig, error) {
    // Parse prefixes, apply strategies, return merged config
}
```

**Dependencies**: Subtasks 2, 4  
**Estimated Time**: 3-4 hours
**Test Requirements**: Complex nested structure merging

### Subtask 6: Add configuration debugging and tracing (🔶 MEDIUM)
**Status**: 📝 NOT STARTED
**Files**: pkg/config/tracing.go (new), pkg/config/debug.go (new)
**Implementation**:
```go
// Enhanced source tracking with inheritance  
type InheritanceSourceTracker struct {
    configSources map[string]ConfigSource
    valueOrigins  map[string]ValueOrigin
    chainMetadata InheritanceChainMetadata
}

type ValueOrigin struct {
    SourceFile    string    // Configuration file containing value
    SourceType    string    // "inherit", "override", "merge", etc.
    MergeStrategy string    // Strategy applied
    LineNumber    int       // Line number in source file
    ChainDepth    int       // Depth in inheritance chain
}
```

**Dependencies**: Subtask 5
**Estimated Time**: 3-4 hours  
**Test Requirements**: Source tracking accuracy

### Subtask 7: Implement comprehensive testing (🔶 MEDIUM)
**Status**: 📝 NOT STARTED
**Files**: pkg/config/inheritance_test.go (new), config_test.go updates
**Implementation**:
```go
// Test functions from testing.md specification
func TestConfigInheritance(t *testing.T) {}
func TestInheritanceChainBuilding(t *testing.T) {}  
func TestMergeStrategyProcessing(t *testing.T) {}
func TestCircularDependencyDetection(t *testing.T) {}
func TestInheritancePathResolution(t *testing.T) {}
func TestConfigurationSourceTracking(t *testing.T) {}
func TestInheritanceErrorHandling(t *testing.T) {}
func TestInheritancePerformance(t *testing.T) {}
```

**Dependencies**: Subtasks 1-6
**Estimated Time**: 8-10 hours
**Test Requirements**: All scenarios from testing.md

### Subtask 8: Create backward compatibility layer (🔶 MEDIUM)  
**Status**: 📝 NOT STARTED
**Files**: config_adapter.go, config.go updates
**Implementation**:
```go
// Enhanced adapter with inheritance support
type ConfigAdapterWithInheritance struct {
    inheritanceLoader *InheritanceConfigLoader
    legacyLoader      *DefaultConfigLoader  
    inheritanceEnabled bool
}

func (a *ConfigAdapterWithInheritance) LoadConfig(root string) (*Config, error) {
    // Check if inheritance is being used
    if a.hasInheritanceDeclarations(root) {
        return a.inheritanceLoader.LoadConfigWithInheritance(root, DefaultConfig())
    }
    // Fall back to legacy loading
    return a.legacyLoader.LoadConfig(root)
}
```

**Dependencies**: Subtask 7
**Estimated Time**: 2-3 hours
**Test Requirements**: Legacy config compatibility

### Subtask 9: Update documentation and examples (🔻 LOW)
**Status**: 📝 NOT STARTED  
**Files**: README.md updates, example configurations
**Implementation**:
- Update main README with inheritance examples
- Create example configuration hierarchies  
- Add troubleshooting guide for common inheritance issues
- Document best practices for inheritance design

**Dependencies**: Subtask 8
**Estimated Time**: 2-3 hours
**Test Requirements**: Documentation accuracy

### Subtask 10: Update task status in feature tracking (🔻 LOW)
**Status**: 📝 NOT STARTED
**Files**: docs/context/feature-tracking.md
**Implementation**:
- Mark completed subtasks with checkmarks [x]
- Update overall feature status to completed
- Document implementation details and performance metrics
- Add notes about notable implementation decisions

**Dependencies**: Subtask 9  
**Estimated Time**: 1 hour
**Test Requirements**: Documentation consistency

## 🎯 SUCCESS CRITERIA CHECKLIST

- [ ] **Explicit Inheritance**: Configuration files can declare inheritance relationships
- [ ] **Flexible Merge Strategies**: Support for override, merge, prepend, replace, and default strategies  
- [ ] **Circular Detection**: Prevent infinite loops from circular inheritance
- [ ] **Backward Compatible**: Existing configurations work without modification
- [ ] **Source Tracking**: Maintain visibility into configuration value origins
- [ ] **Performance**: Minimal overhead for configurations not using inheritance
- [ ] **Comprehensive Testing**: Full test coverage for all inheritance scenarios

## 🚨 RISK MITIGATION

**Technical Risks**:
- **Circular Dependencies**: Mitigated by depth-first search detection
- **Performance Impact**: Mitigated by caching and lazy loading
- **Backward Compatibility**: Mitigated by adapter pattern and feature detection

**Implementation Risks**:  
- **Complex Merging Logic**: Mitigated by strategy pattern and comprehensive testing
- **Type Safety**: Mitigated by reflection-based validation and runtime checks
- **Error Handling**: Mitigated by descriptive error messages and recovery strategies

## 📊 ESTIMATED TIMELINE

**Total Estimated Time**: 35-50 hours
**Phases**:
- **Phase 2A (Subtasks 1-2)**: 6-8 hours
- **Phase 2B (Subtasks 3-5)**: 13-18 hours  
- **Phase 2C (Subtasks 6-10)**: 16-24 hours

**Critical Path**: Subtasks 1 → 2 → 3 → 4 → 5 → 7 → 8

## 🔧 VALIDATION COMMANDS

```bash
# Run comprehensive tests
make test

# Run linting validation  
make lint

# Run icon validation (DOC-008)
make validate-icon-enforcement

# Run inheritance-specific tests
go test ./pkg/config -run TestInheritance -v

# Performance validation
go test ./pkg/config -run TestInheritancePerformance -bench=.
```

## 📋 IMPLEMENTATION NOTES

This plan follows the NEW FEATURE protocol and ensures:
- All immutable requirements preserved
- Comprehensive documentation maintained  
- Full backward compatibility
- Performance characteristics maintained
- Complete test coverage implemented

The implementation leverages the existing pkg/config architecture from EXTRACT-001, ensuring clean integration with minimal disruption to existing functionality. 