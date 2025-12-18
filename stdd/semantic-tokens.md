# Semantic Tokens Directory

**STDD Methodology Version**: 1.0.2

## Overview
This document serves as the **central directory/registry** for all semantic tokens used in the project. Semantic tokens (`[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`) provide a consistent vocabulary and traceability mechanism that ties together all documentation, code, and tests.

**For detailed information about tokens, see:**
- **Requirements tokens**: See `requirements.md` for full descriptions, rationale, satisfaction criteria, and validation criteria
- **Architecture tokens**: See `architecture-decisions.md` for architectural decisions, rationale, and alternatives considered
- **Implementation tokens**: See `implementation-decisions.md` for implementation details, code structures, and algorithms

## AI Assistant Integration Guidelines [REQ:DOC_016]

### Token Usage for AI Assistants

AI assistants should use semantic tokens for:

1. **Code Navigation**: Search for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]` tokens to find related code
2. **Feature Understanding**: Trace features from requirements through architecture to implementation
3. **Change Impact Analysis**: Use token cross-references to identify affected components
4. **Test Discovery**: Find tests for features using `[REQ:*]` tokens in test names

### Token-Based Code Navigation

```bash
# Find all implementations of a requirement
grep -r "\[REQ:FEATURE_NAME\]" --include="*.go" .

# Find all tests for a requirement
grep -r "REQ_FEATURE_NAME" --include="*_test.go" .

# Find architecture decisions for a feature
grep -r "\[ARCH:FEATURE_NAME\]" --include="*.md" .

# Find implementation details
grep -r "\[IMPL:FEATURE_NAME\]" --include="*.go" .
```

### Token Creation Requirements

When implementing features:
1. **ALWAYS** create `[REQ:*]` token in `requirements.md` first
2. **ALWAYS** create `[ARCH:*]` token in `architecture-decisions.md` for design decisions
3. **ALWAYS** add `[IMPL:*]` tokens to code comments
4. **ALWAYS** reference `[REQ:*]` tokens in test names/comments
5. **ALWAYS** update `semantic-tokens.md` registry when creating new tokens

### Token Validation Requirements

Before marking features complete:
1. **ALWAYS** run token validation scripts
2. **ALWAYS** ensure token consistency across all layers
3. **ALWAYS** verify token traceability in documentation
4. **ALWAYS** check that all cross-references are valid

## Token Format

```
[TYPE:IDENTIFIER]
```

## Token Types

- `[REQ:*]` - Requirements (functional/non-functional) - **The source of intent**
- `[ARCH:*]` - Architecture decisions - **High-level design choices that preserve intent**
- `[IMPL:*]` - Implementation decisions - **Low-level choices that preserve intent**
- `[TEST:*]` - Test specifications - **Validation of intent**

## Token Naming Convention

- Use UPPER_SNAKE_CASE for identifiers
- Be descriptive but concise
- Example: `[REQ:DUPLICATE_PREVENTION]` not `[REQ:DP]`

## Cross-Reference Format

When referencing other tokens:

```markdown
[IMPL:EXAMPLE] Description [ARCH:DESIGN] [REQ:REQUIREMENT]
```

## Usage Examples

### In Code Comments
```go
// [REQ:FILE_BACKUP] Create backup of single file with comparison
// [IMPL:ATOMIC_OPS] [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
    // ...
}
```

### In Tests
```go
// Test validates [REQ:FILE_BACKUP] is met
func TestCreateFileBackup_REQ_FILE_BACKUP(t *testing.T) {
    // ...
}
```

### In Documentation
```markdown
The file backup feature uses [ARCH:RESOURCE_MANAGEMENT] to fulfill [REQ:FILE_BACKUP].
Implementation details are documented in [IMPL:ATOMIC_OPS].
```

## Token Validation Guidelines [REQ:DOC_016]

### Cross-Layer Token Consistency

Every feature must have proper token coverage across all layers:

1. **Requirements Layer**: Feature must have `[REQ:*]` token in `requirements.md`
2. **Architecture Layer**: Architecture decisions must have `[ARCH:*]` tokens in `architecture-decisions.md`
3. **Implementation Layer**: Implementation must have `[IMPL:*]` tokens in code comments
4. **Test Layer**: Tests must reference `[REQ:*]` tokens in test names/comments
5. **Documentation Layer**: All documentation must cross-reference tokens consistently

### Token Format Validation

1. **Token Format**: Must follow `[TYPE:IDENTIFIER]` pattern exactly
2. **Token Types**: Must use valid types (`REQ`, `ARCH`, `IMPL`, `TEST`)
3. **Identifier Format**: Must use UPPER_SNAKE_CASE
4. **Cross-References**: Implementation tokens must reference architecture and requirement tokens

### Token Traceability Validation

1. Every requirement in `requirements.md` must have corresponding implementation tokens
2. Every architecture decision must have corresponding implementation tokens
3. Every test must link to specific requirements via `[REQ:*]` tokens
4. All tokens must be discoverable through automated validation

### Automated Validation

Validation scripts check:
- Token format compliance across all files
- Cross-layer token consistency
- Token traceability completeness
- Missing token detection

See `scripts/validate-token-traceability.sh` for validation implementation.

## Token Relationships

### Hierarchical Relationships
- `[REQ:CONFIGURATION]` contains `[REQ:CFG_005]`, `[REQ:CFG_006]`
- `[ARCH:CONFIG_SYSTEM]` includes `[ARCH:CFG_005]`, `[ARCH:CFG_006]`
- `[IMPL:CONFIG_STRUCT]` implements `[ARCH:CONFIG_SYSTEM]`
- `[IMPL:CFG_006]` implements `[ARCH:CFG_006]` and fulfills `[REQ:CFG_006]`

### Flow Relationships
- `[REQ:FEATURE]` → `[ARCH:DESIGN]` → `[IMPL:IMPLEMENTATION]` → Code + Tests

### Dependency Relationships
- `[IMPL:FEATURE]` depends on `[ARCH:DESIGN]` and `[REQ:FEATURE]`
- `[ARCH:DESIGN]` depends on `[REQ:FEATURE]`
- `[ARCH:CFG_006]` depends on `[REQ:CFG_006]` and integrates with `[ARCH:CFG_005]`
- `[IMPL:CFG_006]` depends on `[ARCH:CFG_006]`, `[REQ:CFG_006]`, and `[IMPL:CONFIG_STRUCT]`

## Requirements Tokens Registry

**📖 Full details**: See `requirements.md`

### Core Functionality
- `[REQ:CODE_QUALITY]` → See `requirements.md` § Code Quality and Linting Requirements
- `[REQ:RESOURCE_MANAGEMENT]` → See `requirements.md` § Resource Management Requirements
- `[REQ:ERROR_HANDLING]` → See `requirements.md` § Enhanced Error Handling Requirements
- `[REQ:CONTEXT_SUPPORT]` → See `requirements.md` § Context Support Requirements
- `[REQ:FILE_BACKUP]` → See `requirements.md` § File Backup Requirements
- `[REQ:OUTPUT_FORMATTING]` → See `requirements.md` § Output Formatting Requirements
- `[REQ:TEMPLATE_FORMATTING]` → See `requirements.md` § Template Formatting Requirements
- `[REQ:CONFIGURATION]` → See `requirements.md` § Configuration Management Requirements
- `[REQ:GIT_INTEGRATION]` → See `requirements.md` § Git Integration Requirements
- `[REQ:LIST_LIMIT]` → See `requirements.md` § List Command Output Limit Requirements
- `[REQ:DIFF_COMMAND]` → See `requirements.md` § Diff Command Requirements
- `[REQ:INCREMENTAL_DUPLICATE_PREVENTION]` → See `requirements.md` § Incremental Archive Duplicate Prevention Requirements

### Incomplete Requirements
- `[REQ:INCREMENTAL_DUPLICATE_PREVENTION]` → See `requirements.md` § Incremental Archive Duplicate Prevention Requirements
  - Architecture: `[ARCH:INCREMENTAL_DUPLICATE_PREVENTION]`
  - Implementation: `[IMPL:INCREMENTAL_DUPLICATE_PREVENTION]`
- `[REQ:OUT_002]` → See `requirements.md` § Enhanced Command Output with File Statistics Requirements
- `[REQ:LINT_001]` → See `requirements.md` § Code Linting Compliance Requirements
- `[REQ:DOC_015]` → See `requirements.md` § Unicode to Semantic Token Mapping Requirements
- `[REQ:DOC_016]` → See `requirements.md` § AI-First Comprehensive Token System Requirements
- `[REQ:COV_003]` → See `requirements.md` § Selective Coverage Reporting Requirements
- `[REQ:CICD_001]` → See `requirements.md` § AI-First Development Optimization Requirements
- `[REQ:DOC_011]` → See `requirements.md` § Token Validation Integration for AI Assistants Requirements
- `[REQ:DOC_013]` → See `requirements.md` § AI-First Documentation and Code Maintenance Requirements

### Configuration System Enhancement
- `[REQ:CFG_005]` → See `requirements.md` § Layered Configuration Inheritance Requirements
- `[REQ:CFG_006]` → See `requirements.md` § Complete Configuration Reflection and Visibility Requirements
  - Architecture: `[ARCH:CFG_006]`
  - Implementation: `[IMPL:CFG_006]`
- `[REQ:CONFIG_OUTPUT_GROUPING]` → See `requirements.md` § Configuration Output Grouping Requirements
  - Architecture: `[ARCH:CONFIG_OUTPUT_GROUPING]`
  - Implementation: `[IMPL:CONFIG_OUTPUT_GROUPING]`
- `[REQ:CUSTOMIZABLE_FORMAT_STRINGS]` → See `requirements.md` § User-Customizable Format Strings Requirements
  - Architecture: `[ARCH:CUSTOMIZABLE_FORMAT_STRINGS]`
  - Implementation: `[IMPL:CUSTOMIZABLE_FORMAT_STRINGS]`
- `[REQ:TEST_EXCLUDE_MERGE]` → See `requirements.md` § Exclude Patterns Merge Testing Requirements
  - Architecture: `[ARCH:TEST_EXCLUDE_MERGE]`
  - Implementation: `[IMPL:TEST_EXCLUDE_MERGE]`

### Governance Requirements

- [REQ:GOVERNANCE_DOCUMENTATION] Governance rules must be documented with explicit anchors and cross-links to ARCH/IMPL decisions.
- [REQ:GOV_POLICY_MIGRATION] Migration of governance policies into STDD decision documents.
- [REQ:GOV_CHANGE_CONTROL] Requirements for proposing, approving, and rolling out governance changes.
- [REQ:GOV_AUDIT] Regular auditing of governance rules; ensure alignment with decisions and tests.
- [REQ:GOV_REGISTRY_COMPLETENESS] Registry must capture all governance-related tokens and their intents.
- [REQ:GOV_DISCOVERABILITY] Centralized discoverability of governance content via docs/index.md with links to tokens and decisions.
- [REQ:GOV_NAMING_CONVENTIONS] Establish and enforce token naming conventions for governance-related tokens.
- [REQ:GOV_CROSS_LINKING] Each REQ token must be cross-linked to relevant ARCH/IMPL decisions.
- [REQ:GOV_MIGRATION_STATUS] Track the progress/status of governance-content migration to ARCH/IMPL docs.
- [REQ:REGISTRY_COMPLETION_VERIFICATION] Verification rule that ensures registry remains complete over time.

### Non-Functional Requirements
- `[REQ:PERFORMANCE]` → See `requirements.md` § Performance Requirements
- `[REQ:RELIABILITY]` → See `requirements.md` § Reliability Requirements
- `[REQ:MAINTAINABILITY]` → See `requirements.md` § Maintainability Requirements
- `[REQ:USABILITY]` → See `requirements.md` § Usability Requirements

### Immutable Requirements (Major Version Change Required)
- `[REQ:IMMUTABLE_ARCHIVE_NAMING]` → See `requirements.md` § Archive Naming Convention (Immutable)
- `[REQ:IMMUTABLE_FILE_BACKUP_NAMING]` → See `requirements.md` § File Backup Naming Convention (Immutable)
- `[REQ:IMMUTABLE_DIRECTORY_OPERATIONS]` → See `requirements.md` § Directory Operations Requirements (Immutable)
- `[REQ:IMMUTABLE_FILE_BACKUP_OPERATIONS]` → See `requirements.md` § File Backup Operations Requirements (Immutable)
- `[REQ:IMMUTABLE_FILE_EXCLUSION]` → See `requirements.md` § File Exclusion Requirements (Immutable)
- `[REQ:IMMUTABLE_GIT_INTEGRATION]` → See `requirements.md` § Git Integration Requirements (Immutable)
- `[REQ:IMMUTABLE_ERROR_HANDLING]` → See `requirements.md` § Error Handling Requirements (Immutable)
- `[REQ:IMMUTABLE_CODE_QUALITY]` → See `requirements.md` § Code Quality Standards (Immutable)
- `[REQ:IMMUTABLE_BUILD_SYSTEM]` → See `requirements.md` § Build System Requirements (Immutable)
- `[REQ:IMMUTABLE_OUTPUT_FORMATTING]` → See `requirements.md` § Output Formatting Requirements (Immutable)
- `[REQ:IMMUTABLE_TEMPLATE_FORMATTING]` → See `requirements.md` § Template Formatting Requirements (Immutable)
- `[REQ:IMMUTABLE_CLI_COMMANDS]` → See `requirements.md` § CLI Commands Structure (Immutable)
- `[REQ:IMMUTABLE_CONFIGURATION_DEFAULTS]` → See `requirements.md` § Configuration Defaults (Immutable)
- `[REQ:IMMUTABLE_PLATFORM_COMPATIBILITY]` → See `requirements.md` § Platform Compatibility Requirements (Immutable)
- `[REQ:IMMUTABLE_GLOBAL_OPTIONS]` → See `requirements.md` § Global Options (Immutable)
- `[REQ:IMMUTABLE_RESOURCE_MANAGEMENT]` → See `requirements.md` § Resource Management Requirements (Immutable)
- `[REQ:IMMUTABLE_PERFORMANCE]` → See `requirements.md` § Performance Requirements (Immutable)
- `[REQ:IMMUTABLE_FEATURE_PRESERVATION]` → See `requirements.md` § Feature Preservation Rules (Immutable)
- `[REQ:IMMUTABLE_TESTING_INFRASTRUCTURE]` → See `requirements.md` § Testing Infrastructure Immutable Requirements

## Architecture Tokens Registry

**📖 Full details**: See `architecture-decisions.md`

### Core Architecture Decisions
- `[ARCH:LANGUAGE_SELECTION]` → See `architecture-decisions.md` § Language and Runtime
- `[ARCH:PROJECT_STRUCTURE]` → See `architecture-decisions.md` § Project Structure
- `[ARCH:ARCHIVE_FORMAT]` → See `architecture-decisions.md` § Archive Format
- `[ARCH:CONFIG_SYSTEM]` → See `architecture-decisions.md` § Configuration System
- `[ARCH:CFG_005]` → See `architecture-decisions.md` § Layered Configuration Inheritance
- `[ARCH:CFG_006]` → See `architecture-decisions.md` § Configuration Reflection Architecture [REQ:CFG_006]
- `[ARCH:TEST_EXCLUDE_MERGE]` → See `architecture-decisions.md` § Configuration Testing Architecture [REQ:TEST_EXCLUDE_MERGE]
- `[ARCH:EXCLUDE_MERGE_FIX]` → See `architecture-decisions.md` § Array Field Default Merge Strategy Implementation [REQ:CFG_005]
- `[ARCH:ERROR_HANDLING]` → See `architecture-decisions.md` § Error Handling Strategy
- `[ARCH:RESOURCE_MANAGEMENT]` → See `architecture-decisions.md` § Resource Management
- `[ARCH:CONTEXT_SUPPORT]` → See `architecture-decisions.md` § Context-Aware Operations
- `[ARCH:OUTPUT_FORMATTING]` → See `architecture-decisions.md` § Output Formatting
- `[ARCH:GIT_INTEGRATION]` → See `architecture-decisions.md` § Git Integration
- `[ARCH:TESTING_STRATEGY]` → See `architecture-decisions.md` § Testing Strategy
- `[ARCH:BUILD_DISTRIBUTION]` → See `architecture-decisions.md` § Build and Distribution
- `[ARCH:CODE_ORGANIZATION]` → See `architecture-decisions.md` § Code Organization Principles
- `[ARCH:PACKAGE_EXTRACTION]` → See `architecture-decisions.md` § Package Extraction Architecture
- `[ARCH:CLI_FRAMEWORK]` → See `architecture-decisions.md` § CLI Framework Architecture
- `[ARCH:FILE_OPERATIONS]` → See `architecture-decisions.md` § File Operations Architecture
- `[ARCH:PROCESSING_PATTERNS]` → See `architecture-decisions.md` § Processing Patterns Architecture
- `[ARCH:AUTO_DETECTION]` → See `architecture-decisions.md` § Auto-Detection Architecture
- `[ARCH:FILE_STATISTICS]` → See `architecture-decisions.md` § File Statistics Architecture
- `[ARCH:DIRECTORY_COMPARISON]` → See `architecture-decisions.md` § Directory Comparison Architecture
- `[ARCH:EXCLUSION_PATTERNS]` → See `architecture-decisions.md` § Exclusion Patterns Architecture
- `[ARCH:SYSTEM_COMPONENTS]` → See `architecture-decisions.md` § System Components Architecture
- `[ARCH:SECURITY]` → See `architecture-decisions.md` § Security Architecture
- `[ARCH:EXTENSIBILITY]` → See `architecture-decisions.md` § Extensibility Architecture
- `[ARCH:DEPLOYMENT]` → See `architecture-decisions.md` § Deployment Architecture
- `[ARCH:PERFORMANCE]` → See `architecture-decisions.md` § Performance Architecture
- `[ARCH:CLI_COMMANDS]` → See `architecture-decisions.md` § CLI Commands Architecture
- `[ARCH:CICD_PIPELINE]` → See `architecture-decisions.md` § CI/CD Pipeline Architecture
- `[ARCH:AI_DOCUMENTATION]` → See `architecture-decisions.md` § AI-First Documentation Architecture
- `[ARCH:PERF_VALIDATION]` → See `architecture-decisions.md` § Validation Performance Optimization Architecture
- `[ARCH:TOKEN_SYSTEM]` → See `architecture-decisions.md` § Token System Architecture [REQ:DOC_016]
- `[ARCH:CONFIG_OUTPUT_GROUPING]` → See `architecture-decisions.md` § Configuration Output Grouping [REQ:CONFIG_OUTPUT_GROUPING]
- `[ARCH:LIST_LIMIT]` → See `architecture-decisions.md` § List Command Limit Architecture [REQ:LIST_LIMIT]
- `[ARCH:DIFF_COMMAND]` → See `architecture-decisions.md` § Diff Command Architecture [REQ:DIFF_COMMAND]
- `[ARCH:INCREMENTAL_DUPLICATE_PREVENTION]` → See `architecture-decisions.md` § Incremental Archive Duplicate Prevention Architecture [REQ:INCREMENTAL_DUPLICATE_PREVENTION]

## Implementation Tokens Registry

**📖 Full details**: See `implementation-decisions.md`

### Core Implementation Decisions
- `[IMPL:CONFIG_STRUCT]` → See `implementation-decisions.md` § Configuration Structure [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]
- `[IMPL:CONFIG_DISPLAY_FLATTENING]` → See `implementation-decisions.md` § Configuration Display Flattening [ARCH:CFG_006] [REQ:CFG_006]
- `[IMPL:CFG_006]` → See `implementation-decisions.md` § Configuration Reflection Implementation [ARCH:CFG_006] [REQ:CFG_006]
- `[IMPL:TEST_EXCLUDE_MERGE]` → See `implementation-decisions.md` § Exclude Patterns Merge Testing [ARCH:TEST_EXCLUDE_MERGE] [REQ:TEST_EXCLUDE_MERGE]
- `[IMPL:EXCLUDE_MERGE_FIX]` → See `implementation-decisions.md` § Array Field Default Merge Strategy Implementation [ARCH:EXCLUDE_MERGE_FIX] [REQ:CFG_005]
- `[IMPL:CFG_MIXED_MODE_MERGE_FIX]` → See `implementation-decisions.md` § Mixed-Mode Merge Strategy Fix [ARCH:CFG_005] [ARCH:CFG_001] [REQ:CFG_005]
- `[IMPL:CFG_MIXED_SEQUENTIAL_INHERITANCE]` → See `implementation-decisions.md` § Mixed Sequential and Inheritance File Processing [ARCH:CFG_005] [REQ:CFG_005] [REQ:CFG_001]
- `[IMPL:CFG_QUOTED_KEY_PREFIX]` → See `implementation-decisions.md` § Quoted YAML Key Prefix Support [ARCH:CFG_005] [REQ:CFG_005]
- `[IMPL:ZIP_FORMAT]` → See `implementation-decisions.md` § ZIP Archive Format Implementation [ARCH:ARCHIVE_FORMAT]
- `[IMPL:DUAL_FORMATTING]` → See `implementation-decisions.md` § Dual Printf/Template Formatting [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]
- `[IMPL:STRUCTURED_ERRORS]` → See `implementation-decisions.md` § Structured Error Handling [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING]
- `[IMPL:GIT_CLI]` → See `implementation-decisions.md` § Git Command-line Integration [ARCH:GIT_INTEGRATION] [REQ:GIT_INTEGRATION]
- `[IMPL:RESOURCE_MANAGER]` → See `implementation-decisions.md` § Resource Management with Cleanup [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]
- `[IMPL:CONTEXT_OPS]` → See `implementation-decisions.md` § Context-Aware Operations [ARCH:CONTEXT_SUPPORT] [REQ:CONTEXT_SUPPORT]
- `[IMPL:ATOMIC_OPS]` → See `implementation-decisions.md` § Atomic File Operations [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]
- `[IMPL:TESTING]` → See `implementation-decisions.md` § Testing Implementation [ARCH:TESTING_STRATEGY] [REQ:*]
- `[IMPL:CODE_STYLE]` → See `implementation-decisions.md` § Code Style and Conventions
- `[IMPL:PACKAGE_EXTRACTION]` → See `implementation-decisions.md` § Package Extraction Implementation [ARCH:PACKAGE_EXTRACTION] [REQ:MAINTAINABILITY]
- `[IMPL:CLI_FRAMEWORK]` → See `implementation-decisions.md` § CLI Framework Implementation [ARCH:CLI_FRAMEWORK] [REQ:USABILITY]
- `[IMPL:FILE_OPERATIONS]` → See `implementation-decisions.md` § File Operations Implementation [ARCH:FILE_OPERATIONS] [REQ:RELIABILITY]
- `[IMPL:PROCESSING_PATTERNS]` → See `implementation-decisions.md` § Processing Patterns Implementation [ARCH:PROCESSING_PATTERNS] [REQ:PERFORMANCE]
- `[IMPL:AUTO_DETECTION]` → See `implementation-decisions.md` § Auto-Detection Implementation [ARCH:AUTO_DETECTION] [REQ:USABILITY]
- `[IMPL:FILE_STATISTICS]` → See `implementation-decisions.md` § File Statistics Implementation [ARCH:FILE_STATISTICS] [REQ:OUTPUT_FORMATTING]
- `[IMPL:FILE_STATISTICS_TEMPLATE_FIX]` → See `implementation-decisions.md` § File Statistics Template Processing Fix [ARCH:FILE_STATISTICS] [REQ:OUT_002] [REQ:OUTPUT_FORMATTING]
- `[IMPL:DIRECTORY_COMPARISON]` → See `implementation-decisions.md` § Directory Comparison Implementation [ARCH:DIRECTORY_COMPARISON]
- `[IMPL:EXCLUSION_PATTERNS]` → See `implementation-decisions.md` § Exclusion Patterns Implementation [ARCH:EXCLUSION_PATTERNS] [REQ:CONFIGURATION]
- `[IMPL:CFG_006]` → See `implementation-decisions.md` § Configuration Reflection Implementation [ARCH:CFG_006] [REQ:CFG_006]
- `[IMPL:DATA_MODELS]` → See `implementation-decisions.md` § Data Models [ARCH:SYSTEM_COMPONENTS] [REQ:*]
- `[IMPL:DOC_ENHANCEMENT]` → See `implementation-decisions.md` § Documentation Enhancement Framework [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_001]
- `[IMPL:SEMANTIC_CROSS_REF]` → See `implementation-decisions.md` § Semantic Cross-Referencing Strategy [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_001]
- `[IMPL:TRACEABILITY]` → See `implementation-decisions.md` § Enhanced Traceability Design [ARCH:DOCUMENTATION_ARCHITECTURE] [REQ:DOC_003]
- `[IMPL:REFACTOR_PREP]` → See `implementation-decisions.md` § Pre-Extraction Refactoring Strategy [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:INTERFACE_FIRST]` → See `implementation-decisions.md` § Interface-First Extraction Approach [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:LARGE_FILE_DECOMP]` → See `implementation-decisions.md` § Large File Decomposition Strategy [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:GIT_DIRTY_CONFIG]` → See `implementation-decisions.md` § GIT-006: Configurable Git Dirty Status [ARCH:GIT_INTEGRATION] [REQ:GIT_INTEGRATION]
- `[IMPL:CONFIGURABLE_STRINGS]` → See `implementation-decisions.md` § CFG-004: Eliminate Hardcoded Strings [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]
- `[IMPL:DELAYED_OUTPUT]` → See `implementation-decisions.md` § OUT-001: Delayed Output Management [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]
- `[IMPL:TEST_COVERAGE]` → See `implementation-decisions.md` § TEST-001: Comprehensive formatter.go Test Coverage [ARCH:TESTING_STRATEGY] [REQ:*]
- `[IMPL:EXTRACTION_PRINCIPLES]` → See `implementation-decisions.md` § Extraction Principles and Design Decisions [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:BACKWARD_COMPAT]` → See `implementation-decisions.md` § Maintain Backward Compatibility [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:INTERFACE_DRIVEN]` → See `implementation-decisions.md` § Interface-Driven Design [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:ZERO_BREAKING]` → See `implementation-decisions.md` § Zero-Breaking-Change Extraction [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:LAYERED_EXTRACTION]` → See `implementation-decisions.md` § Layered Extraction Approach [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:EXTRACTION_CHALLENGES]` → See `implementation-decisions.md` § Extraction Challenges and Solutions [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:LARGE_FILE_CHALLENGE]` → See `implementation-decisions.md` § Large File Decomposition Challenge [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:CONFIG_SCHEMA_FLEX]` → See `implementation-decisions.md` § Configuration Schema Flexibility [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]
- `[IMPL:DEPENDENCY_MGMT]` → See `implementation-decisions.md` § Dependency Management [ARCH:CODE_ORGANIZATION] [REQ:MAINTAINABILITY]
- `[IMPL:TESTING_COMPLEXITY]` → See `implementation-decisions.md` § Testing Complexity [ARCH:TESTING_STRATEGY] [REQ:*]
- `[IMPL:TOKEN_SYSTEM]` → See `implementation-decisions.md` § Token System Implementation [ARCH:TOKEN_SYSTEM] [REQ:DOC_016]
- `[IMPL:CONFIG_OUTPUT_GROUPING]` → See `implementation-decisions.md` § Configuration Output Grouping Implementation [ARCH:CONFIG_OUTPUT_GROUPING] [REQ:CONFIG_OUTPUT_GROUPING]
- `[IMPL:CFG_PRECEDENCE_FIX]` → See `implementation-decisions.md` § Configuration File Precedence Fix [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]
- `[IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX]` → See `implementation-decisions.md` § Merge and Prepend Strategy Precedence Fix [ARCH:CFG_001] [ARCH:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION]
- `[IMPL:TEST_CFG_005_P1]` → See `implementation-decisions.md` § Priority 1 Configuration Merge Tests Implementation [ARCH:CFG_005] [REQ:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION]
- `[IMPL:TEST_UNICODE_HANDLING]` → See `implementation-decisions.md` § Unicode and Special Character Handling Tests [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]
- `[IMPL:TEST_EMPTY_STRING_HANDLING]` → See `implementation-decisions.md` § Empty String Handling Tests [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005]
- `[IMPL:TEST_PREPEND_ORDERING]` → See `implementation-decisions.md` § Prepend Strategy Ordering Tests [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]
- `[IMPL:TEST_DEFAULT_STRATEGY_EDGES]` → See `implementation-decisions.md` § Default Strategy Edge Cases Tests [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]
- `[IMPL:CFG_INHERITANCE_PATH_RESOLUTION]` → See `implementation-decisions.md` § Inheritance Path Resolution Fix [ARCH:CFG_005] [REQ:CFG_005] [REQ:CONFIGURATION]
- `[IMPL:TEST_UNICODE_HANDLING]` → See `implementation-decisions.md` § Unicode and Special Character Handling Tests [ARCH:TESTING_STRATEGY] [REQ:CONFIGURATION] [REQ:CFG_005]
- `[IMPL:LIST_LIMIT]` → See `implementation-decisions.md` § List Command Limit Implementation [ARCH:LIST_LIMIT] [REQ:LIST_LIMIT]
- `[IMPL:DIFF_COMMAND]` → See `implementation-decisions.md` § Diff Command Implementation [ARCH:DIFF_COMMAND] [REQ:DIFF_COMMAND]
- `[IMPL:INCREMENTAL_DUPLICATE_PREVENTION]` → See `implementation-decisions.md` § Incremental Archive Duplicate Prevention Implementation [ARCH:INCREMENTAL_DUPLICATE_PREVENTION] [REQ:INCREMENTAL_DUPLICATE_PREVENTION]


## Legacy Token Migration Inventory [REQ:DOC_016] [REQ:GOV_REGISTRY_COMPLETENESS]

**Status**: ✅ Completed (2025-12-18)

The legacy `project-tokens.yaml` registry has been migrated into this document so that every identifier inherits the same traceability guarantees as the existing `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`, and `[TEST:*]` tokens. Each entry below cross-links its governing requirement, architectural anchor, implementation reference, and validating tests/code. The migration was validated using `scripts/token-coverage-analysis.sh` and `scripts/validate-token-traceability.sh` and achieved 100% coverage.

Refer to `docs/validation-reports/token-coverage-analysis.md` for the validation report and `project-tokens.yaml` for the pointer stub to the canonical registry.

### Action Tokens
- `[ACTION:core-functionality]` (`core-functionality`, Essential system operations) → Requirements: [REQ:FILE_BACKUP], [REQ:RESOURCE_MANAGEMENT], [REQ:CONFIGURATION]; Architecture: [ARCH:SYSTEM_COMPONENTS], [ARCH:RESOURCE_MANAGEMENT]; Implementation: [IMPL:RESOURCE_MANAGER], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `backup_test.go`, `comparison_test.go`, `config_test.go`.
- `[ACTION:format-processing]` (`format-processing`, Text formatting and output generation) → Requirements: [REQ:OUTPUT_FORMATTING], [REQ:TEMPLATE_FORMATTING]; Architecture: [ARCH:OUTPUT_FORMATTING], [ARCH:FILE_STATISTICS]; Implementation: [IMPL:DUAL_FORMATTING], [IMPL:FILE_STATISTICS]; Tests/Code: `formatter_test.go`, `formatter_adapter_test.go`.
- `[ACTION:discovery]` (`discovery`, File system and configuration discovery) → Requirements: [REQ:CONFIGURATION], [REQ:GIT_INTEGRATION]; Architecture: [ARCH:CONFIG_SYSTEM], [ARCH:GIT_INTEGRATION]; Implementation: [IMPL:CONFIG_STRUCT], [IMPL:GIT_CLI]; Tests/Code: `config_test.go`, `git_test.go`.
- `[ACTION:maintenance]` (`maintenance`, Code cleanup and refactoring) → Requirements: [REQ:MAINTAINABILITY], [REQ:CODE_QUALITY]; Architecture: [ARCH:CODE_ORGANIZATION], [ARCH:PACKAGE_EXTRACTION]; Implementation: [IMPL:REFACTOR_PREP], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `refactoring_validation_test.go`, `inc_diff_integration_test.go`.
- `[ACTION:validation]` (`validation`, Input validation and error checking) → Requirements: [REQ:CODE_QUALITY], [REQ:ERROR_HANDLING]; Architecture: [ARCH:TESTING_STRATEGY], [ARCH:ERROR_HANDLING]; Implementation: [IMPL:TESTING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `main_test.go`, `errors_test.go`.
- `[ACTION:migration]` (`migration`, System migration and token replacement) → Requirements: [REQ:DOC_015], [REQ:DOC_016]; Architecture: [ARCH:TOKEN_SYSTEM]; Implementation: [IMPL:TOKEN_SYSTEM]; Tests/Code: `scripts/validate-token-traceability.sh`, `scripts/token-coverage-analysis.sh`.

### Feature Tokens – Immutable Operations & Services
- `DIR-001` (Directory Operations Immutable) → [REQ:IMMUTABLE_DIRECTORY_OPERATIONS], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `comparison_test.go`, `archive.go`; Source: `docs/context/immutable.md` § Directory Operations.
- `EXCLUDE-001` (File Exclusion Requirements Immutable) → [REQ:IMMUTABLE_FILE_EXCLUSION], [ARCH:EXCLUSION_PATTERNS], [IMPL:EXCLUSION_PATTERNS]; Tests/Code: `exclude_test.go`, `exclude.go`; Source: `docs/context/immutable.md` § File Exclusion.
- `CONTEXT-001` (Context Support Requirements) → [REQ:CONTEXT_SUPPORT], [ARCH:CONTEXT_SUPPORT], [IMPL:CONTEXT_OPS]; Tests/Code: `internal/testutil/context_test.go`, `config.go`; Source: `docs/context/requirements.md` § Context Support.
- `CONFIG-DISCOVERY-001` (Configuration Discovery Specification) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`; Source: `docs/context/specification.md` § Configuration Discovery.
- `CONFIG-FILE-001` (Configuration File Specification) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_integration_test.go`, `config.go`; Source: `docs/context/specification.md` § Configuration File.
- `CLI-GLOBAL-001` (Global Options Specification) → [REQ:IMMUTABLE_GLOBAL_OPTIONS], [ARCH:CLI_COMMANDS], [IMPL:CLI_FRAMEWORK]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`; Source: `docs/context/specification.md` § Global Options.
- `ARCHIVE-FEATURES-001` (Archive Features Specification) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`; Source: `docs/context/specification.md` § Archive Features.
- `ERROR-HANDLING-001` (Error Handling & Recovery Specification) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`; Source: `docs/context/specification.md` § Error Handling.
- `RESOURCE-MGMT-001` (Resource Management Specification) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `file_stats.go`; Source: `docs/context/specification.md` § Resource Management.
- `CORE-ARCH-001` (Core Architecture Decision) → [REQ:MAINTAINABILITY], [ARCH:PROJECT_STRUCTURE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `main_test.go`, `main.go`; Source: `docs/context/architecture.md` § Core Architecture.
- `SYSTEM-COMPONENTS-001` (System Components Architecture Decision) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, `architecture-decisions.md`; Source: `docs/context/architecture.md` § System Components.
- `DATA-MODELS-001` (Data Models Architecture Decision) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, `pkg/cli/builder.go`; Source: `docs/context/architecture.md` § Data Models.
- `SERVICE-ARCH-001` (Service Architecture Decision) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, service layer files; Source: `docs/context/architecture.md` § Service Architecture.
- `SERVICE-ARCHIVE-001` (Archive Service Architecture Decision) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`; Source: `docs/context/architecture.md` § Archive Service.
- `SERVICE-BACKUP-001` (File Backup Service Architecture Decision) → [REQ:FILE_BACKUP], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `backup_test.go`, `backup.go`; Source: `docs/context/architecture.md` § Backup Service.
- `SERVICE-GIT-001` (Git Integration Service Architecture Decision) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`; Source: `docs/context/architecture.md` § Git Integration Service.
- `SERVICE-RESOURCE-001` (Resource Management Service Architecture Decision) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `file_stats.go`; Source: `docs/context/architecture.md` § Resource Management Service.
- `SERVICE-TEMPLATE-001` (Template Formatting Service Architecture Decision) → [REQ:TEMPLATE_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`; Source: `docs/context/architecture.md` § Template Formatting Service.
- `SERVICE-OUTPUT-001` (Output Formatting Service Architecture Decision) → [REQ:OUTPUT_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`; Source: `docs/context/architecture.md` § Output Formatting Service.
- `SERVICE-ERROR-001` (Error Handling Service Architecture Decision) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`; Source: `docs/context/architecture.md` § Error Handling Service.

### Feature Tokens – Data, Formatting & Platform Infrastructure
- `ARCH-001` (Archive naming convention implementation) → [REQ:IMMUTABLE_ARCHIVE_NAMING], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `CFG-003` (Template formatting logic) → [REQ:TEMPLATE_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_adapter_test.go`, `formatter.go`.
- `GIT-004` (Git submodule support) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `DOC-014` (AI-first documentation system) → [REQ:DOC_013], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/index.md`, `docs/ai-assistant-compliance.md`.
- `DOC-015` (Unicode to semantic token mapping system) → [REQ:DOC_015], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `scripts/validate-token-traceability.sh`, `docs/semantic-token-system-implementation-complete.md`.
- `DOC-016` (AI-first comprehensive token system) → [REQ:DOC_016], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `scripts/token-coverage-analysis.sh`, `docs/semantic-token-system-implementation-complete.md`.
- `DOC-017` (AI assistant token protocol) → [REQ:DOC_016], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `.ai-agent-instructions`, `scripts/validate-token-traceability.sh`.
- `TOKEN-001` (Semantic token migration system) → [REQ:DOC_016], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `scripts/validate-token-traceability.sh`, `tools/coverage_differential_test.go`.
- `QUALITY-001` (Code quality standards) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:CODE_STYLE]; Tests/Code: `revive.toml`, `main_test.go`.
- `BUILD-001` (Build system requirements) → [REQ:IMMUTABLE_BUILD_SYSTEM], [ARCH:BUILD_DISTRIBUTION], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `Makefile`, `tools/validate-coverage.sh`.
- `FORMAT-001` (Output formatting requirements) → [REQ:OUTPUT_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`.
- `TEMPLATE-001` (Template formatting requirements) → [REQ:TEMPLATE_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`.
- `CMD-001` (Commands interface) → [REQ:IMMUTABLE_CLI_COMMANDS], [ARCH:CLI_COMMANDS], [IMPL:CLI_FRAMEWORK]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`.
- `TEMPLATE-002` (Enhanced template formatting requirements) → [REQ:TEMPLATE_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_adapter_test.go`, `formatter.go`.
- `DATA-CONFIG-001` (Config data object requirements) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `DATA-CONFIGVALUE-001` (ConfigValue data object requirements) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `DATA-BACKUP-001` (Backup data object requirements) → [REQ:FILE_BACKUP], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `backup_test.go`, `backup.go`.
- `DATA-BACKUPERROR-001` (BackupError data object requirements) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `backup_test.go`, `errors.go`.
- `DATA-ARCHIVEERROR-001` (ArchiveError data object requirements) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `archive_test.go`, `errors.go`.
- `DATA-ARCHIVE-001` (Archive data object requirements) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `DATA-BACKUPINFO-001` (BackupInfo data object requirements) → [REQ:FILE_BACKUP], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `backup_test.go`, `backup.go`.
- `DATA-RESOURCEMANAGER-001` (ResourceManager data object requirements) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `file_stats.go`.
- `DATA-TEMPLATEFORMATTER-001` (TemplateFormatter data object requirements) → [REQ:TEMPLATE_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`.

### Feature Tokens – Standards, Specifications & Architecture References
- `BUILD-DEV-001` (Build and development requirements specification) → [REQ:IMMUTABLE_BUILD_SYSTEM], [ARCH:BUILD_DISTRIBUTION], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `Makefile`, `tools/validate-coverage.sh`.
- `QUALITY-STANDARDS-001` (Quality assurance and code standards specification) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:CODE_STYLE]; Tests/Code: `revive.toml`, `main_test.go`.
- `LINTING-001` (Linting requirements specification) → [REQ:LINT_001], [ARCH:TESTING_STRATEGY], [IMPL:CODE_STYLE]; Tests/Code: `revive.toml`, `config_test.go`.
- `ERROR-STANDARDS-001` (Error handling standards specification) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`.
- `RESOURCE-STANDARDS-001` (Resource management standards specification) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `file_stats.go`.
- `CONFIG-ARCH-001` (Configuration architecture decision) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `CONFIG-DISCOVERY-ARCH-001` (Configuration discovery architecture decision) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `CONFIG-SOURCES-ARCH-001` (Configuration sources architecture decision) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `DATA-CONFIG-CORE-001` (Core configuration object architecture decision) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `DATA-ENHANCED-001` (Enhanced data objects architecture decision) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, `pkg/cli`.
- `SPEC-001` (Quality assurance and code standards specification) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:TESTING]; Tests/Code: `main_test.go`, `revive.toml`.
- `SPEC-002` (Configuration discovery system specification) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `SPEC-003` (Configuration file specification) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_integration_test.go`, `config.go`.
- `SPEC-004` (Command interface specification) → [REQ:IMMUTABLE_CLI_COMMANDS], [ARCH:CLI_COMMANDS], [IMPL:CLI_FRAMEWORK]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`.
- `SPEC-005` (Global options specification) → [REQ:IMMUTABLE_GLOBAL_OPTIONS], [ARCH:CLI_COMMANDS], [IMPL:CLI_FRAMEWORK]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`.
- `SPEC-006` (Archive features specification) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `SPEC-007` (Error handling and recovery specification) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`.
- `SPEC-008` (Resource management specification) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `file_stats.go`.
- `SPEC-009` (Build and development requirements specification) → [REQ:IMMUTABLE_BUILD_SYSTEM], [ARCH:BUILD_DISTRIBUTION], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `Makefile`, `scripts/validate-test-tokens.sh`.
- `SPEC-010` (Testing infrastructure specification) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:TESTING]; Tests/Code: `main_test.go`, `test/scenarios`.
- `SPEC-011` (Implementation details specification) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `refactoring_validation_test.go`, `docs/implementation-decisions.md`.
- `SPEC-012` (Platform compatibility specification) → [REQ:IMMUTABLE_PLATFORM_COMPATIBILITY], [ARCH:PROJECT_STRUCTURE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `main_test.go`, `Makefile`.
- `SPEC-013` (Performance characteristics specification) → [REQ:PERFORMANCE], [ARCH:PERFORMANCE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `inc_diff_integration_test.go`, `comparison_test.go`.
- `SPEC-014` (CI/CD pipeline optimization specification) → [REQ:CICD_001], [ARCH:CICD_PIPELINE], [IMPL:TESTING]; Tests/Code: `tools/validate-coverage.sh`, `scripts/validate-token-traceability.sh`.
- `SPEC-015` (AI-first documentation specification) → [REQ:DOC_013], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/index.md`, `docs/ai-assistant-compliance.md`.

### Feature Tokens – Requirement Identifiers & Architecture Specifications
- `REQ-001` (Code quality and linting requirements) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:CODE_STYLE]; Tests/Code: `main_test.go`, `revive.toml`.
- `REQ-002` (Data objects requirements) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `REQ-003` (Core functions requirements) → [REQ:RESOURCE_MANAGEMENT], [ARCH:RESOURCE_MANAGEMENT], [IMPL:RESOURCE_MANAGER]; Tests/Code: `file_stats_test.go`, `main.go`.
- `REQ-004` (Main application structure requirements) → [REQ:MAINTAINABILITY], [ARCH:PROJECT_STRUCTURE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `main_test.go`, `main.go`.
- `REQ-005` (Git integration requirements) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `REQ-007` (Configuration system enhancement requirements) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `ARCH-SPEC-001` (Core architecture specification) → [REQ:MAINTAINABILITY], [ARCH:PROJECT_STRUCTURE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `main_test.go`, `architecture-decisions.md`.
- `ARCH-SPEC-002` (Data models architecture) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, `pkg/cli`.
- `ARCH-SPEC-003` (Service architecture) → [REQ:MAINTAINABILITY], [ARCH:SYSTEM_COMPONENTS], [IMPL:DATA_MODELS]; Tests/Code: `main_test.go`, `architecture-decisions.md`.
- `ARCH-SPEC-004` (Configuration architecture) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `ARCH-SPEC-005` (Archive format architecture) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `ARCH-SPEC-006` (Output formatting architecture) → [REQ:OUTPUT_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DUAL_FORMATTING]; Tests/Code: `formatter_test.go`, `formatter.go`.
- `ARCH-SPEC-007` (Error handling architecture) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`.
- `ARCH-SPEC-008` (Concurrency architecture) → [REQ:PERFORMANCE], [ARCH:PROCESSING_PATTERNS], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `inc_diff_integration_test.go`, `main_test.go`.
- `ARCH-SPEC-009` (Testing architecture) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:TESTING]; Tests/Code: `main_test.go`, `test/scenarios`.
- `ARCH-SPEC-010` (Security architecture) → [REQ:RELIABILITY], [ARCH:SECURITY], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `docs/security`.
- `ARCH-SPEC-011` (Extensibility architecture) → [REQ:MAINTAINABILITY], [ARCH:EXTENSIBILITY], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `refactoring_validation_test.go`, `architecture-decisions.md`.
- `ARCH-SPEC-012` (Deployment architecture) → [REQ:MAINTAINABILITY], [ARCH:DEPLOYMENT], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `Makefile`, `run-ubuntu-with-backup.sh`.
- `ARCH-SPEC-013` (Performance considerations architecture) → [REQ:PERFORMANCE], [ARCH:PERFORMANCE], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `inc_diff_integration_test.go`, `comparison_test.go`.
- `ARCH-SPEC-014` (CLI commands architecture) → [REQ:IMMUTABLE_CLI_COMMANDS], [ARCH:CLI_COMMANDS], [IMPL:CLI_FRAMEWORK]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`.

### Feature Tokens – Configuration Output, Refactoring & Extraction
- `CFG-004` (Comprehensive string config) → [REQ:CUSTOMIZABLE_FORMAT_STRINGS], [ARCH:OUTPUT_FORMATTING], [IMPL:CONFIGURABLE_STRINGS]; Tests/Code: `formatter_test.go`, `config_test.go`.
- `OUT-001` (Delayed output management) → [REQ:OUTPUT_FORMATTING], [ARCH:OUTPUT_FORMATTING], [IMPL:DELAYED_OUTPUT]; Tests/Code: `formatter_test.go`, `formatter.go`.
- `OUT-002` (Enhanced command output with file statistics) → [REQ:OUT_002], [ARCH:FILE_STATISTICS], [IMPL:FILE_STATISTICS]; Tests/Code: `formatter_stats_test.go`, `formatter_stats.go`.
- `REFACTOR-001` (Dependency analysis and interface standardization) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:REFACTOR_PREP]; Tests/Code: `refactoring_validation_test.go`, `docs/refactoring-validation-report.md`.
- `REFACTOR-003` (Configuration schema abstraction) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:CONFIG_SCHEMA_FLEX]; Tests/Code: `config_test.go`, `config_impl.go`.
- `REFACTOR-004` (Error handling consolidation) → [REQ:ERROR_HANDLING], [ARCH:ERROR_HANDLING], [IMPL:STRUCTURED_ERRORS]; Tests/Code: `errors_test.go`, `errors.go`.
- `REFACTOR-005` (Code structure optimization) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `refactoring_validation_test.go`, `docs/refactoring-validation-report.md`.
- `REFACTOR-006` (Refactoring impact validation) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:REFACTOR_PREP]; Tests/Code: `refactoring_validation_test.go`, `inc_diff_integration_test.go`.
- `EXTRACT-001` (Configuration management system extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `config_test.go`, `docs/extract-008-cli-template-plan.md`.
- `EXTRACT-002` (Error handling system extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `errors_test.go`, `docs/extract-008-subtask-2-completion.md`.
- `EXTRACT-003` (Formatter system extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `formatter_test.go`, `docs/extract-008-subtask-3-completion.md`.
- `EXTRACT-004` (Git integration extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `git_test.go`, `docs/extract-008-subtask-5-completion.md`.
- `EXTRACT-005` (CLI system extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `cmd/ai-validation/main.go`, `docs/extract-008-remaining-subtasks-plan.md`.
- `EXTRACT-006` (File operations extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `file_stats_test.go`, `docs/extract-009-completion-summary.md`.
- `EXTRACT-009` (Test utilities extraction) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `internal/testutil/context_test.go`, `docs/extract-008-session-summary.md`.
- `EXTRACT-010` (Package documentation and examples) → [REQ:MAINTAINABILITY], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/examples`, `docs/extract-010.md`.

### Feature Tokens – Comparison & Diff Capabilities
- `COMPARISON-001` (File comparison functionality) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-002` (Comparison utilities) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-003` (Comparison testing) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-004` (Comparison edge cases) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-005` (Comparison performance) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-006` (Comparison integration) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-007` (Comparison validation) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-009` (Comparison utilities) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-010` (Comparison edge cases) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `COMPARISON-011` (Comparison testing) → [REQ:DIFF_COMMAND], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.

### Feature Tokens – CLI, Archive, Configuration & Test Coverage
- `FEATURE-001` (Feature implementation) → [REQ:IMMUTABLE_FEATURE_PRESERVATION], [ARCH:SYSTEM_COMPONENTS], [IMPL:PROCESSING_PATTERNS]; Tests/Code: `main_test.go`, `main.go`.
- `CLI-015` (Automatic file/directory command detection) → [REQ:IMMUTABLE_CLI_COMMANDS], [ARCH:AUTO_DETECTION], [IMPL:AUTO_DETECTION]; Tests/Code: `cmd/ai-validation/main.go`, `main.go`.
- `ARCH-002` (Create archive command) → [REQ:FILE_BACKUP], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `ARCH-003` (Incremental archives) → [REQ:INCREMENTAL_DUPLICATE_PREVENTION], [ARCH:INCREMENTAL_DUPLICATE_PREVENTION], [IMPL:INCREMENTAL_DUPLICATE_PREVENTION]; Tests/Code: `incremental_duplicate_prevention_test.go`, `archive.go`.
- `ARCH-004` (Broken symlink handling) → [REQ:FILE_BACKUP], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `archive_test.go`, `archive.go`.
- `FILE-001` (File backup naming) → [REQ:IMMUTABLE_FILE_BACKUP_NAMING], [ARCH:ARCHIVE_FORMAT], [IMPL:ZIP_FORMAT]; Tests/Code: `archive_test.go`, `archive.go`.
- `FILE-002` (Backup command) → [REQ:FILE_BACKUP], [ARCH:FILE_OPERATIONS], [IMPL:FILE_OPERATIONS]; Tests/Code: `backup_test.go`, `backup.go`.
- `FILE-003` (File comparison) → [REQ:FILE_BACKUP], [ARCH:DIRECTORY_COMPARISON], [IMPL:DIRECTORY_COMPARISON]; Tests/Code: `comparison_test.go`, `comparison.go`.
- `CFG-001` (Config discovery) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `CFG-002` (Status codes) → [REQ:CONFIGURATION], [ARCH:CONFIG_SYSTEM], [IMPL:CONFIG_STRUCT]; Tests/Code: `config_test.go`, `config.go`.
- `CFG-005` (Layered configuration inheritance) → [REQ:CFG_005], [ARCH:CFG_005], [IMPL:EXCLUDE_MERGE_FIX]; Tests/Code: `config_integration_test.go`, `config_impl.go`.
- `CFG-006` (Complete configuration reflection and visibility) → [REQ:CFG_006], [ARCH:CFG_006], [IMPL:CFG_006]; Tests/Code: `config_test.go`, `config.go`.
- `GIT-001` (Git info extraction) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `GIT-003` (Git status detection) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `GIT-005` (Git configuration integration) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `TEST-001` (Comprehensive formatter testing) → [REQ:CODE_QUALITY], [ARCH:TESTING_STRATEGY], [IMPL:TESTING]; Tests/Code: `formatter_test.go`, `formatter_adapter_test.go`.

### Semantic Token Metadata Entries
- `document_structure` (Document section markers) → [REQ:DOC_015], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/index.md`, `docs/ai-assistant-compliance.md`.
- `status_indicators` (Status icon mapping) → [REQ:DOC_015], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/context/index.md`, `docs/unicode-migration-phase2-critical-tokens-completion-summary.md`.
- `action_categories` (Action taxonomy) → [REQ:DOC_016], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `.ai-agent-instructions`, `scripts/token-coverage-analysis.sh`.
- `context_specific` (Contextual documentation tags) → [REQ:DOC_015], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/context/specification.md`, `docs/integration-guide.md`.
- `DOC-010` (Automated token format suggestions) → [REQ:DOC_016], [ARCH:AI_DOCUMENTATION], [IMPL:DOC_ENHANCEMENT]; Tests/Code: `docs/format-strings-reference.md`, `scripts/validate-token-traceability.sh`.
- `DOC-011` (Token validation integration for AI assistants) → [REQ:DOC_011], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `scripts/validate-token-traceability.sh`, `tools/validate-token-traceability.sh`.
- `DOC-012` (Real-time icon validation feedback) → [REQ:DOC_016], [ARCH:TOKEN_SYSTEM], [IMPL:TOKEN_SYSTEM]; Tests/Code: `tools/coverage_differential_test.go`, `docs/format-strings-reference.md`.
- `EXTRACT-008` (Session extraction utilities) → [REQ:MAINTAINABILITY], [ARCH:PACKAGE_EXTRACTION], [IMPL:PACKAGE_EXTRACTION]; Tests/Code: `docs/extract-008-session-summary.md`, `refactoring_validation_test.go`.
- `REFACTOR-002` (Large file decomposition preparation) → [REQ:MAINTAINABILITY], [ARCH:CODE_ORGANIZATION], [IMPL:LARGE_FILE_DECOMP]; Tests/Code: `formatter_adapter_simple_test.go`, `docs/formatter-decomposition.md`.
- `GIT-002` (Branch/hash naming) → [REQ:GIT_INTEGRATION], [ARCH:GIT_INTEGRATION], [IMPL:GIT_CLI]; Tests/Code: `git_test.go`, `git.go`.
- `COV-002` (Coverage baseline establishment) → [REQ:COV_003], [ARCH:TESTING_STRATEGY], [IMPL:TEST_COVERAGE]; Tests/Code: `tools/coverage_differential_test.go`, `scripts/validate-coverage.sh`.

## Quick Reference Index

### Requirements Tokens by Category
- **Core Functionality**: `CODE_QUALITY`, `RESOURCE_MANAGEMENT`, `ERROR_HANDLING`, `CONTEXT_SUPPORT`, `FILE_BACKUP`, `OUTPUT_FORMATTING`, `TEMPLATE_FORMATTING`, `CONFIGURATION`, `GIT_INTEGRATION`, `LIST_LIMIT`
- **Configuration Enhancement**: `CFG_005`, `CFG_006`
- **Non-Functional**: `PERFORMANCE`, `RELIABILITY`, `MAINTAINABILITY`, `USABILITY`
- **Immutable Requirements**: `IMMUTABLE_ARCHIVE_NAMING`, `IMMUTABLE_FILE_BACKUP_NAMING`, `IMMUTABLE_DIRECTORY_OPERATIONS`, `IMMUTABLE_FILE_BACKUP_OPERATIONS`, `IMMUTABLE_FILE_EXCLUSION`, `IMMUTABLE_GIT_INTEGRATION`, `IMMUTABLE_ERROR_HANDLING`, `IMMUTABLE_CODE_QUALITY`, `IMMUTABLE_BUILD_SYSTEM`, `IMMUTABLE_OUTPUT_FORMATTING`, `IMMUTABLE_TEMPLATE_FORMATTING`, `IMMUTABLE_CLI_COMMANDS`, `IMMUTABLE_CONFIGURATION_DEFAULTS`, `IMMUTABLE_PLATFORM_COMPATIBILITY`, `IMMUTABLE_GLOBAL_OPTIONS`, `IMMUTABLE_RESOURCE_MANAGEMENT`, `IMMUTABLE_PERFORMANCE`, `IMMUTABLE_FEATURE_PRESERVATION`, `IMMUTABLE_TESTING_INFRASTRUCTURE`
- **Incomplete Requirements**: `OUT_002`, `LINT_001`, `DOC_015`, `DOC_016`, `COV_003`, `CICD_001`, `DOC_011`, `DOC_013`

### Architecture Tokens
- **Core**: `LANGUAGE_SELECTION`, `PROJECT_STRUCTURE`, `SYSTEM_COMPONENTS`, `ARCHIVE_FORMAT`, `CONFIG_SYSTEM`, `CFG_005`, `CFG_006`, `ERROR_HANDLING`, `RESOURCE_MANAGEMENT`, `CONTEXT_SUPPORT`, `OUTPUT_FORMATTING`, `GIT_INTEGRATION`, `TESTING_STRATEGY`, `BUILD_DISTRIBUTION`, `CODE_ORGANIZATION`
- **Archive Operations**: `ARCH_001`, `ARCH_002`, `ARCH_003`, `ARCH_004`
- **CLI Features**: `CLI_015`, `CLI_COMMANDS`, `LIST_LIMIT`
- **Extraction & Modularization**: `PACKAGE_EXTRACTION`, `CLI_FRAMEWORK`, `FILE_OPERATIONS`, `PROCESSING_PATTERNS`
- **Features**: `AUTO_DETECTION`, `FILE_STATISTICS`, `DIRECTORY_COMPARISON`, `EXCLUSION_PATTERNS`
- **Security & Quality**: `SECURITY`, `EXTENSIBILITY`, `DEPLOYMENT`, `PERFORMANCE`, `PERF_VALIDATION`
- **AI & CI/CD**: `CICD_PIPELINE`, `AI_DOCUMENTATION`, `TOKEN_SYSTEM`

### Implementation Tokens
- **Core**: `CONFIG_STRUCT`, `CFG_006`, `ZIP_FORMAT`, `DUAL_FORMATTING`, `STRUCTURED_ERRORS`, `GIT_CLI`, `RESOURCE_MANAGER`, `CONTEXT_OPS`, `ATOMIC_OPS`, `TESTING`, `CODE_STYLE`, `DATA_MODELS`
- **File Operations**: `FILE_001`, `FILE_002`, `FILE_003`
- **Configuration**: `CFG_001`, `CFG_002`, `CFG_003`, `CFG_004`, `CFG_005`, `CFG_TEMPLATE_001`, `CONFIGURABLE_STRINGS`, `CFG_PRECEDENCE_FIX`, `CFG_MERGE_PREPEND_PRECEDENCE_FIX`
- **Git Integration**: `GIT_001`, `GIT_002`, `GIT_003`, `GIT_004`, `GIT_005`, `GIT_006`, `GIT_DIRTY_CONFIG`
- **Output Management**: `OUT_001`, `OUT_002`, `DELAYED_OUTPUT`, `LIST_LIMIT`
- **Testing Infrastructure**: `TEST_001`, `TEST_002`, `TEST_FIX_001`, `TEST_INFRA_001_B`, `TEST_INFRA_001_E`, `TEST_COVERAGE`, `TESTING_COMPLEXITY`, `TEST_UNICODE_HANDLING`, `TEST_EMPTY_STRING_HANDLING`, `TEST_PREPEND_ORDERING`, `TEST_DEFAULT_STRATEGY_EDGES`
- **Code Quality**: `LINT_001`, `COV_001`, `COV_002`, `COV_003`
- **Documentation**: `DOC_001` through `DOC_017`, `INSTALL_001`, `DOC_ENHANCEMENT`, `SEMANTIC_CROSS_REF`, `TRACEABILITY`, `TOKEN_SYSTEM`
- **Extraction & Modularization**: `PACKAGE_EXTRACTION`, `CLI_FRAMEWORK`, `FILE_OPERATIONS`, `PROCESSING_PATTERNS`, `EXTRACT_001` through `EXTRACT_010`, `EXTRACTION_PRINCIPLES`, `BACKWARD_COMPAT`, `INTERFACE_DRIVEN`, `ZERO_BREAKING`, `LAYERED_EXTRACTION`, `EXTRACTION_CHALLENGES`, `LARGE_FILE_CHALLENGE`, `DEPENDENCY_MGMT`
- **Refactoring**: `REFACTOR_001` through `REFACTOR_006`, `REFACTOR_PREP`, `INTERFACE_FIRST`, `LARGE_FILE_DECOMP`
- **Performance**: `PERF_001`
- **Features**: `AUTO_DETECTION`, `FILE_STATISTICS`, `DIRECTORY_COMPARISON`, `EXCLUSION_PATTERNS`
- **Configuration Schema**: `CONFIG_SCHEMA_FLEX`
- **Output Control**: `DEBUG-OUTPUT`, `DIAGNOSTIC-OUTPUT`
