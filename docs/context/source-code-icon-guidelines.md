# Source Code Icon Guidelines - Implementation Token Standardization

## 📑 Purpose
This document establishes the official guidelines for using standardized icons in source code implementation tokens. It ensures consistent icon usage between documentation and code, improving traceability and AI assistant understanding.

## [ACTION:validation] Standardized Implementation Token Format

### ✨ **Core Format Structure**
```go
// [PRIORITY_ICON] FEATURE-ID: Brief description [- DETAIL_ICON Context]
```

### 🆕 **Enhanced Format with Decision Context** (DOC-014)
```go
// [PRIORITY_ICON] FEATURE-ID: Brief description [DECISION: context1, context2, context3]
// [PRIORITY_ICON] FEATURE-ID: Brief description - [ACTION_ICON] Context [DECISION: context1, context2, context3]
```

**Standard Examples:**
```go
// [CRITICAL] ARCH-001: Archive naming convention implementation - [ACTION:core-functionality] Core functionality
// [HIGH] CFG-003: Template formatting logic - [ACTION:format-processing] Format string processing  
// [MEDIUM] GIT-004: Git submodule support - [ACTION:discovery] Discovery and validation
// [LOW] REFACTOR-001: Interface standardization - 🏗️ Architecture preparation
```

**Enhanced Examples with Decision Context:**
```go
// [CRITICAL] ARCH-001: Archive naming convention [DECISION: core-functionality, blocks-extraction, user-facing]
// [HIGH] CFG-005: Configuration inheritance [DECISION: enhancement, enables-flexibility, backward-compatible]
// [MEDIUM] DOC-010: Token format suggestions [DECISION: developer-experience, independent]
// [LOW] TEST-FIX-001: Config isolation [DECISION: infrastructure, prerequisite-for-testing]
```

#### 🏷️ **Decision Context Categories** (DOC-014 Integration)

**Impact Level Categories:**
- **`core-functionality`**: Essential system operations (archive creation, data integrity, critical error handling)
- **`enhancement`**: Feature improvements and extensions (Git integration, template formatting, configuration inheritance)  
- **`developer-experience`**: Tools and utilities for development workflow (testing infrastructure, validation tools, debugging)
- **`infrastructure`**: Build system, CI/CD, project structure (test framework, linting, coverage analysis)

**Dependency Categories:**
- **`blocks-extraction`**: Required for component extraction work (EXTRACT-001-010 series)
- **`enables-flexibility`**: Enables future feature development and system extensibility
- **`prerequisite-for-X`**: Required for specific feature X implementation (replace X with feature ID)
- **`independent`**: No blocking dependencies, can be implemented standalone

**Constraint Categories:**
- **`user-facing`**: Affects UI, command behavior, or visible functionality (CLI commands, output formatting, error messages)
- **`backward-compatible`**: Preserves existing functionality, doesn't break current behavior
- **`breaking-change`**: May break existing functionality, requires major version bump
- **`performance-critical`**: Significant performance impact, requires optimization consideration

#### 📋 **Decision Context Usage Guidelines**

**Mandatory Decision Context:**
- [CRITICAL] CRITICAL priority implementations
- Components involved in extraction planning (EXTRACT-001-010)
- User-facing functionality changes
- Breaking changes or major architectural decisions

**Recommended Decision Context:**
- [HIGH] HIGH priority implementations
- Features that enable future development
- Performance-critical implementations
- Infrastructure and tooling improvements

**Optional Decision Context:**
- [MEDIUM] MEDIUM priority implementations
- [LOW] LOW priority implementations
- Straightforward bug fixes
- Routine maintenance tasks

#### [ACTION:core-functionality] **Migration to Enhanced Format**

**Context Selection Process:**
1. **Identify Impact Level**: Core/enhancement/developer-experience/infrastructure
2. **Assess Dependencies**: Blocks-extraction/enables-flexibility/prerequisite-for-X/independent
3. **Evaluate Constraints**: User-facing/backward-compatible/breaking-change/performance-critical

**Enhanced Migration Examples:**
```go
// Before: Standard format
// [CRITICAL] ARCH-002: Archive creation with context handling

// After: Enhanced format with decision context
// [CRITICAL] ARCH-002: Archive creation with context handling [DECISION: core-functionality, performance-critical, backward-compatible]
```

### 🔥 Priority Icon Integration

#### [CRITICAL] **CRITICAL PRIORITY Tokens**
**When to Use:**
- ✅ Core system functionality that blocks other operations
- ✅ Critical error handling and recovery paths
- ✅ Essential data integrity operations
- ✅ Mandatory security and validation functions

**Format:**
```go
// [CRITICAL] ARCH-001: Archive creation with atomic operations
func CreateArchiveWithCleanup(ctx context.Context, cfg *Config) error {
    // [CRITICAL] ARCH-001: Critical resource management for atomic operations
    rm := NewResourceManager()
    defer rm.Cleanup()
    
    // [CRITICAL] ARCH-001: Essential compression and file operations
    return processArchiveCreation(ctx, cfg, rm)
}
```

#### [HIGH] **HIGH PRIORITY Tokens**
**When to Use:**
- ✅ Important feature implementations
- ✅ Significant configuration processing
- ✅ Core business logic functions
- ✅ Major API and interface implementations

**Format:**
```go
// [HIGH] CFG-003: Template formatting engine implementation
func ProcessTemplateFormat(template string, data interface{}) (string, error) {
    // [HIGH] CFG-003: High-priority template parsing and validation
    tmpl, err := template.New("format").Parse(template)
    if err != nil {
        return "", fmt.Errorf("template parsing failed: %w", err)
    }
    
    // [HIGH] CFG-003: Template execution with error handling
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", fmt.Errorf("template execution failed: %w", err)
    }
    
    return buf.String(), nil
}
```

#### [MEDIUM] **MEDIUM PRIORITY Tokens**
**When to Use:**
- ✅ Conditional feature implementations
- ✅ Optional enhancement functions
- ✅ Secondary processing logic
- ✅ Configuration helpers and utilities

**Format:**
```go
// [MEDIUM] GIT-004: Git submodule detection - conditional feature
func DetectGitSubmodules(repoPath string) ([]string, error) {
    // [MEDIUM] GIT-004: Medium-priority submodule discovery
    gitmodulesPath := filepath.Join(repoPath, ".gitmodules")
    if !fileExists(gitmodulesPath) {
        return nil, nil // No submodules
    }
    
    // [MEDIUM] GIT-004: Parse .gitmodules file for submodule paths
    return parseGitmodulesFile(gitmodulesPath)
}
```

#### [LOW] **LOW PRIORITY Tokens**
**When to Use:**
- ✅ Cleanup and maintenance functions
- ✅ Optional optimization code
- ✅ Debug and logging utilities
- ✅ Backward compatibility helpers

**Format:**
```go
// [LOW] REFACTOR-001: Legacy compatibility helper - low priority cleanup
func LegacyArchiveNameWrapper(prefix string, timestamp time.Time) string {
    // [LOW] REFACTOR-001: Backward compatibility for old naming convention
    return fmt.Sprintf("%s-%s.zip", prefix, timestamp.Format("2006-01-02-15-04"))
}
```

## [ACTION:core-functionality] Action Category Icons in Implementation Tokens

### [ACTION:discovery] **SEARCH/DISCOVER Action Tokens**
**When to Use:**
- ✅ File system discovery and scanning
- ✅ Configuration file search and detection
- ✅ Git repository analysis and information extraction
- ✅ Pattern matching and validation logic

**Format:**
```go
// [HIGH] CFG-001: Configuration discovery - [ACTION:discovery] Search and detection
func GetConfigSearchPath() []string {
    // [HIGH] CFG-001: Environment variable search - [ACTION:discovery] Discovery
    if envPath := os.Getenv("BKPDIR_CONFIG"); envPath != "" {
        return strings.Split(envPath, ":")
    }
    
    // [HIGH] CFG-001: Default path search - [ACTION:discovery] Fallback discovery
    return []string{"./.bkpdir.yml", "~/.bkpdir.yml"}
}
```

### [ACTION:format-processing] **DOCUMENT/UPDATE Action Tokens**
**When to Use:**
- ✅ Configuration file writing and updates
- ✅ Log file creation and maintenance
- ✅ Status file updates and persistence
- ✅ Output formatting and generation

**Format:**
```go
// [HIGH] CFG-002: Configuration persistence - [ACTION:format-processing] File updates
func SaveConfigData(configPath string, data *ConfigData) error {
    // [HIGH] CFG-002: YAML serialization - [ACTION:format-processing] Format and write
    yamlData, err := yaml.Marshal(data)
    if err != nil {
        return fmt.Errorf("YAML marshaling failed: %w", err)
    }
    
    // [HIGH] CFG-002: Atomic file writing - [ACTION:format-processing] Safe persistence
    return writeFileAtomic(configPath, yamlData, 0644)
}
```

### [ACTION:core-functionality] **CONFIGURE/MODIFY Action Tokens**
**When to Use:**
- ✅ System configuration and setup
- ✅ Parameter adjustment and tuning
- ✅ Environment preparation and initialization
- ✅ Tool and utility configuration

**Format:**
```go
// [HIGH] CFG-003: Output formatter configuration - [ACTION:core-functionality] System setup
func ConfigureFormatter(cfg *Config) (*OutputFormatter, error) {
    // [HIGH] CFG-003: Printf formatter setup - [ACTION:core-functionality] Configuration
    printf := &PrintfFormatter{
        StatusCodeConfig: cfg.StatusCodes,
        FormatStrings:   cfg.FormatStrings,
    }
    
    // [HIGH] CFG-003: Template formatter setup - [ACTION:core-functionality] Configuration  
    template := &TemplateFormatter{
        Templates:     cfg.Templates,
        RegexPatterns: cfg.RegexPatterns,
    }
    
    return &OutputFormatter{
        PrintfFormatter:   printf,
        TemplateFormatter: template,
    }, nil
}
```

### [ACTION:validation] **PROTECT/VALIDATE Action Tokens**
**When to Use:**
- ✅ Input validation and sanitization
- ✅ Security checks and enforcement
- ✅ Data integrity verification
- ✅ Error detection and prevention

**Format:**
```go
// [CRITICAL] ARCH-002: Archive verification - [ACTION:validation] Integrity protection
func VerifyArchiveIntegrity(archivePath string) error {
    // [CRITICAL] ARCH-002: ZIP structure validation - [ACTION:validation] Corruption detection
    zipReader, err := zip.OpenReader(archivePath)
    if err != nil {
        return fmt.Errorf("archive validation failed: %w", err)
    }
    defer zipReader.Close()
    
    // [CRITICAL] ARCH-002: Archive integrity verification - [ACTION:validation] Data protection
    return validateArchiveStructure(zipReader)
}
```

## 🆕 Enhanced Token Format with Decision Context (DOC-014)

### 🎯 **Purpose of Decision Context**

Decision context enhances implementation tokens by providing explicit rationale for implementation decisions. This improves AI assistant guidance by clarifying:
- **Why** the implementation was chosen
- **What** constraints influenced the decision  
- **How** the implementation fits into the overall architecture

### 🏷️ **Decision Context Categories**

#### **Impact Level Categories**
- **`core-functionality`**: Essential system operations fundamental to application purpose
- **`enhancement`**: Feature improvements that add value but aren't core requirements
- **`developer-experience`**: Tools and utilities that improve development workflow
- **`infrastructure`**: Build system, CI/CD, project structure, foundational support

#### **Dependency Categories**  
- **`blocks-extraction`**: Required for component extraction work (EXTRACT-001-010)
- **`enables-flexibility`**: Enables future feature development and extensibility
- **`prerequisite-for-X`**: Required for specific feature X implementation
- **`independent`**: No blocking dependencies, can be implemented standalone

#### **Constraint Categories**
- **`user-facing`**: Affects user interface, command behavior, or visible functionality
- **`backward-compatible`**: Preserves existing functionality and behavior
- **`breaking-change`**: May break existing functionality (requires major version)
- **`performance-critical`**: Impacts system performance significantly

### [ACTION:format-processing] **Decision Context Usage Guidelines**

#### **Mandatory Decision Context**
```go
// [CRITICAL] ARCH-001: Archive naming convention [DECISION: core-functionality, blocks-extraction, user-facing]
func GenerateArchiveName(cfg *Config, timestamp time.Time) string {
    // [CRITICAL] ARCH-001: Core naming logic with user-facing implications
    return fmt.Sprintf("%s-%s.zip", cfg.Prefix, timestamp.Format("2006-01-02T150405"))
}
```

#### **Recommended Decision Context**
```go
// [HIGH] CFG-005: Configuration inheritance system [DECISION: enhancement, enables-flexibility]
func LoadConfigWithInheritance(configPath string) (*Config, error) {
    // [HIGH] CFG-005: Enhanced config loading enabling future extensibility
    return loadConfigRecursive(configPath, make(map[string]bool))
}
```

#### **Optional Decision Context**
```go
// [MEDIUM] DOC-010: Token format suggestions [DECISION: developer-experience]
func SuggestTokenFormat(functionName string) (*TokenSuggestion, error) {
    // [MEDIUM] DOC-010: AI assistant guidance utility
    return analyzeFunction(functionName)
}
```

### [ACTION:core-functionality] **Integration with Existing Token Systems**

#### **Compatibility with Action Icons**
```go
// [HIGH] CFG-003: Template formatting - [ACTION:format-processing] Output generation [DECISION: enhancement, user-facing, backward-compatible]
func ProcessTemplateFormat(template string, data interface{}) (string, error) {
    // [HIGH] CFG-003: Template parsing with backward compatibility - [ACTION:format-processing] Format processing
    tmpl, err := template.New("format").Parse(template)
    return executeTemplate(tmpl, data)
}
```

#### **Migration from Standard Format**
**Before (Standard Format):**
```go
// [HIGH] GIT-001: Git info extraction - [ACTION:discovery] Repository analysis
```

**After (Enhanced Format):**
```go
// [HIGH] GIT-001: Git info extraction [DECISION: enhancement, independent, user-facing]
```

### 🚨 **Validation Requirements**

#### **Format Validation**
- Decision context must use square brackets: `[DECISION: ...]`
- Context tags must be comma-separated with single spaces
- Maximum 3 context tags per token
- All context tags must be from approved categories

#### **Quality Assurance**
- Priority icon must match implementation complexity ([CRITICAL] for core, [LOW] for infrastructure)
- Feature ID must exist in feature-tracking.md
- Decision context must add meaningful information for AI assistance
- Context tags must accurately reflect implementation purpose

### 📊 **Example Decision Context Applications**

#### **Core Functionality with Extraction Blocking**
```go
// [CRITICAL] REFACTOR-001: Interface standardization [DECISION: core-functionality, blocks-extraction]
func CreateStandardizedInterface() InterfaceContract {
    // [CRITICAL] REFACTOR-001: Critical interface prep blocking component extraction
    return defineComponentBoundaries()
}
```

#### **Enhancement Enabling Future Development** 
```go
// [HIGH] CFG-005: Configuration inheritance [DECISION: enhancement, enables-flexibility, backward-compatible]
func ProcessConfigInheritance(configs []string) (*Config, error) {
    // [HIGH] CFG-005: Flexible config system enabling future extensibility
    return mergeConfigurationChain(configs)
}
```

#### **Developer Experience with Independence**
```go
// [MEDIUM] DOC-012: Real-time validation feedback [DECISION: developer-experience, independent]
func ProvideRealTimeValidation(code string) (*ValidationResult, error) {
    // [MEDIUM] DOC-012: Standalone development experience enhancement
    return validateInRealTime(code)
}
```

#### **Infrastructure with Testing Prerequisites**
```go
// [LOW] TEST-FIX-001: Personal config isolation [DECISION: infrastructure, prerequisite-for-testing]
func IsolateTestConfiguration() error {
    // [LOW] TEST-FIX-001: Test infrastructure foundation for reliable testing
    return setupConfigurationIsolation()
}
    // [CRITICAL] ARCH-002: Checksum verification - [ACTION:validation] Data integrity
    return validateArchiveChecksums(zipReader)
}
```

## 🏗️ Component Category Icons

### 🏗️ **ARCHITECTURE Component Tokens**
**When to Use:**
- ✅ Core system architecture components
- ✅ Service layer implementations
- ✅ Interface definitions and contracts
- ✅ System boundary and integration points

**Format:**
```go
// [HIGH] ARCH-001: Archive service architecture - 🏗️ Core component
type ArchiveService struct {
    config      *Config
    formatter   OutputFormatter
    git         GitService
    resourceMgr *ResourceManager
}

// [HIGH] ARCH-001: Archive service constructor - 🏗️ Component initialization
func NewArchiveService(cfg *Config) *ArchiveService {
    return &ArchiveService{
        config:      cfg,
        formatter:   NewOutputFormatter(cfg),
        git:         NewGitService(),
        resourceMgr: NewResourceManager(),
    }
}
```

### 🧪 **TESTING Component Tokens**
**When to Use:**
- ✅ Test infrastructure and utilities
- ✅ Mock implementations and fixtures
- ✅ Test data generation and setup
- ✅ Validation and assertion helpers

**Format:**
```go
// [HIGH] TEST-001: Archive testing infrastructure - 🧪 Test utilities
func CreateTestArchive(t *testing.T, files map[string]string) string {
    // [HIGH] TEST-001: Temporary archive creation - 🧪 Test fixture
    tempDir := t.TempDir()
    archivePath := filepath.Join(tempDir, "test.zip")
    
    // [HIGH] TEST-001: Test archive population - 🧪 Data generation
    return createArchiveWithFiles(archivePath, files)
}
```

## 📋 Special Implementation Token Patterns

### 🚀 **PHASE Implementation Tokens**
**Used for process-oriented implementations:**
```go
// 🚀 Phase 1: Validation and setup
func performValidationPhase(ctx context.Context) error {
    // [HIGH] CFG-001: Config validation - 🚀 Phase 1 setup
    if err := validateConfiguration(); err != nil {
        return err
    }
    return nil
}

// ⚡ Phase 2: Core processing
func performProcessingPhase(ctx context.Context) error {
    // [CRITICAL] ARCH-002: Archive creation - ⚡ Phase 2 execution
    return createArchiveFiles(ctx)
}

// [ACTION:migration] Phase 3: Cleanup and finalization
func performCleanupPhase(ctx context.Context) error {
    // [LOW] REFACTOR-001: Resource cleanup - [ACTION:migration] Phase 3 finalization
    return cleanupTemporaryResources()
}
```

### 1️⃣2️⃣3️⃣ **STEP Implementation Tokens**
**Used for sequential step implementations:**
```go
// [HIGH] CFG-003: Template processing workflow
func ProcessTemplate(template string, data interface{}) (string, error) {
    // 1️⃣ Step 1: Template parsing and validation
    tmpl, err := parseAndValidateTemplate(template)
    if err != nil {
        return "", err
    }
    
    // 2️⃣ Step 2: Data preparation and formatting
    formattedData, err := prepareTemplateData(data)
    if err != nil {
        return "", err
    }
    
    // 3️⃣ Step 3: Template execution and output generation
    result, err := executeTemplate(tmpl, formattedData)
    if err != nil {
        return "", err
    }
    
    // ✅ Complete: Successful template processing
    return result, nil
}
```

## [ACTION:migration] Migration from Legacy Implementation Tokens

### ❌ **Legacy Patterns to Replace**
```go
// OLD PATTERN - No priority indication
// ARCH-001: Archive naming convention implementation

// OLD PATTERN - No action context
// CFG-003: Format string processing

// OLD PATTERN - No component categorization  
// GIT-002: Git branch and hash naming
```

### ✅ **New Standardized Patterns**
```go
// NEW PATTERN - Clear priority and context
// [CRITICAL] ARCH-001: Archive naming convention implementation - [ACTION:core-functionality] Core functionality

// NEW PATTERN - Action category included
// [HIGH] CFG-003: Format string processing - [ACTION:format-processing] Template formatting  

// NEW PATTERN - Component architecture context
// [HIGH] GIT-002: Git branch and hash naming - 🏗️ Naming service
```

### [ACTION:migration] **Migration Command Pattern**
```bash
# Search and replace legacy patterns
sed -i 's|// ARCH-\([0-9]\+\): \(.*\)|// [CRITICAL] ARCH-\1: \2 - [ACTION:core-functionality] Core functionality|g' *.go
sed -i 's|// CFG-\([0-9]\+\): \(.*\)|// [HIGH] CFG-\1: \2 - [ACTION:core-functionality] Configuration|g' *.go
sed -i 's|// FILE-\([0-9]\+\): \(.*\)|// [CRITICAL] FILE-\1: \2 - 📁 File operations|g' *.go
sed -i 's|// GIT-\([0-9]\+\): \(.*\)|// [HIGH] GIT-\1: \2 - 🏗️ Git integration|g' *.go
```

## ✅ Icon Consistency Validation

### [ACTION:validation] **Validation Rules**
1. **Priority Icon Consistency**: Priority icons in code must match feature priority in feature-tracking.md
2. **Action Category Appropriateness**: Action icons must reflect actual function behavior
3. **Component Classification**: Architecture icons must align with component boundaries
4. **No Conflicting Meanings**: Same icon must have same meaning everywhere

### [ACTION:core-functionality] **Validation Script Pattern**
```bash
#!/bin/bash
# Validate implementation token icon consistency

echo "[ACTION:discovery] Validating implementation token icon consistency..."

# Check for consistent priority icons
grep -r "// [CRITICAL]" *.go | while read line; do
    feature_id=$(echo "$line" | grep -o '[A-Z]\+-[0-9]\+')
    # Verify this feature has CRITICAL priority in feature-tracking.md
    if ! grep -q "$feature_id.*[CRITICAL].*CRITICAL" docs/context/feature-tracking.md; then
        echo "❌ Priority mismatch: $feature_id uses [CRITICAL] in code but not CRITICAL in docs"
    fi
done

# Check for appropriate action icons
grep -r "// .*[ACTION:discovery]" *.go | while read line; do
    # Verify function actually performs search/discovery operations
    if ! echo "$line" | grep -qE "(find|search|discover|detect|scan|lookup)"; then
        echo "⚠️ Action icon mismatch: [ACTION:discovery] used but function may not be search-related"
    fi
done

echo "✅ Icon consistency validation complete"
```

## 📋 Implementation Token Checklist

### ✅ **Before Committing Code**
- [ ] **Priority Icons**: All implementation tokens use correct priority icons ([CRITICAL][HIGH][MEDIUM][LOW])
- [ ] **Action Context**: Functions include appropriate action icons ([ACTION:discovery][ACTION:format-processing][ACTION:core-functionality][ACTION:validation])
- [ ] **Component Categories**: Architecture components marked with 🏗️ when appropriate
- [ ] **Phase/Step Indicators**: Process implementations use 🚀⚡[ACTION:migration]🏁 or 1️⃣2️⃣3️⃣✅
- [ ] **Consistency Check**: Icons match meanings in documentation system
- [ ] **No Legacy Patterns**: All old implementation tokens updated to new format

### [ACTION:core-functionality] **Automated Validation Integration**
```makefile
# Add to Makefile
validate-icons:
    @echo "[ACTION:discovery] Validating implementation token icons..."
    @./scripts/validate-icon-consistency.sh
    @echo "✅ Icon validation complete"

# Integrate with existing targets
lint: validate-icons
    @revive -config revive.toml -formatter friendly ./...

ci-lint: validate-icons
    @revive -config revive.toml -formatter github_actions ./...
```

## 🎯 Expected Outcomes

### 📈 **Benefits of Standardized Implementation Tokens**
1. **Improved Traceability**: Direct visual connection between code and documentation
2. **Enhanced AI Assistant Understanding**: Clear semantic meaning for automated analysis
3. **Better Code Navigation**: Priority and context immediately visible in code
4. **Consistent Developer Experience**: Same icon meanings across all project artifacts
5. **Simplified Code Review**: Reviewers can quickly assess priority and purpose

### 🚀 **Integration with Development Workflow**
1. **IDE Support**: Icons visible in code editors and IDEs
2. **Git Integration**: Implementation tokens visible in diffs and blame
3. **Documentation Generation**: Automated cross-referencing between code and docs
4. **CI/CD Validation**: Automated checking of icon consistency
5. **Code Quality Gates**: Implementation token compliance as quality requirement

---

**📋 Guidelines Status**: Official implementation standard for DOC-007
**[ACTION:validation] Enforcement**: Integrated with linting and validation processes  
**[CRITICAL] Critical**: Mandatory for all new code and recommended for legacy updates

// DOC-007: Source code icon integration - Implementation guidelines complete 