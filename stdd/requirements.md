# Requirements Directory

**STDD Methodology Version**: 1.0.1

## Overview
This document serves as the **central listing/registry** for all requirements in this project. Each requirement has a unique semantic token `[REQ:IDENTIFIER]` for traceability.

**For detailed information about how requirements are fulfilled, see:**
- **Architecture decisions**: See `architecture-decisions.md` for high-level design choices that fulfill requirements
- **Implementation decisions**: See `implementation-decisions.md` for detailed implementation approaches, APIs, and algorithms
- **Semantic tokens**: See `semantic-tokens.md` for the complete token registry

### Requirement Structure

Each requirement includes:
- **Description**: What the requirement specifies (WHAT)
- **Rationale**: Why the requirement exists (WHY)
- **Satisfaction Criteria**: How we know the requirement is satisfied (acceptance criteria, success conditions)
- **Validation Criteria**: How we verify/validate the requirement is met (testing approach, verification methods, success metrics)

**Note**: 
- Satisfaction and validation criteria that involve architectural or implementation details reference the appropriate layers
- Architecture decisions in `architecture-decisions.md` explain HOW requirements are fulfilled at a high level
- Implementation decisions in `implementation-decisions.md` explain HOW requirements are fulfilled at a detailed level

## Requirements Registry

### Core Functionality Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| `[REQ:CODE_QUALITY]` | Code Quality and Linting | P0 | ✅ | See `architecture-decisions.md` § Code Organization Principles | See `implementation-decisions.md` § Code Style and Conventions |
| `[REQ:RESOURCE_MANAGEMENT]` | Resource Management | P0 | ✅ | See `architecture-decisions.md` § Resource Management | See `implementation-decisions.md` § Resource Management with Cleanup |
| `[REQ:ERROR_HANDLING]` | Enhanced Error Handling | P0 | ✅ | See `architecture-decisions.md` § Error Handling Strategy | See `implementation-decisions.md` § Structured Error Handling |
| `[REQ:CONTEXT_SUPPORT]` | Context Support | P0 | ✅ | See `architecture-decisions.md` § Context-Aware Operations | See `implementation-decisions.md` § Context-Aware Operations |
| `[REQ:FILE_BACKUP]` | File Backup | P0 | ✅ | See `architecture-decisions.md` § Resource Management | See `implementation-decisions.md` § Atomic File Operations |
| `[REQ:OUTPUT_FORMATTING]` | Output Formatting | P0 | ✅ | See `architecture-decisions.md` § Output Formatting | See `implementation-decisions.md` § Dual Printf/Template Formatting |
| `[REQ:TEMPLATE_FORMATTING]` | Template Formatting | P0 | ✅ | See `architecture-decisions.md` § Output Formatting | See `implementation-decisions.md` § Dual Printf/Template Formatting |
| `[REQ:CONFIGURATION]` | Configuration Management | P0 | ✅ | See `architecture-decisions.md` § Configuration System | See `implementation-decisions.md` § Configuration Structure |
| `[REQ:GIT_INTEGRATION]` | Git Integration | P1 | ✅ | See `architecture-decisions.md` § Git Integration | See `implementation-decisions.md` § Git Command-line Integration |
| `[REQ:ARCHIVE_VERIFICATION]` | Archive Verification | P1 | ✅ | See `architecture-decisions.md` § Archive Format | See `implementation-decisions.md` § ZIP Archive Format Implementation |

### Configuration System Enhancement Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| `[REQ:CFG_005]` | Layered Configuration Inheritance | P0 | ✅ | See `architecture-decisions.md` § Layered Configuration Inheritance | See `implementation-decisions.md` § Configuration Structure |
| `[REQ:CFG_006]` | Complete Configuration Reflection and Visibility | P1 | ✅ | See `architecture-decisions.md` § Configuration System | See `implementation-decisions.md` § Configuration Structure |
| `[REQ:TEST_EXCLUDE_MERGE]` | Exclude Patterns Merge Testing | P1 | ✅ | See `architecture-decisions.md` § Configuration Testing Architecture [ARCH:TEST_EXCLUDE_MERGE] | See `implementation-decisions.md` § Exclude Patterns Merge Testing [IMPL:TEST_EXCLUDE_MERGE] |

### Non-Functional Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| `[REQ:PERFORMANCE]` | Performance | P1 | ✅ | See `architecture-decisions.md` § Language and Runtime | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:RELIABILITY]` | Reliability | P0 | ✅ | See `architecture-decisions.md` § Resource Management, Error Handling Strategy | See `implementation-decisions.md` § Resource Management with Cleanup, Structured Error Handling |
| `[REQ:MAINTAINABILITY]` | Maintainability | P1 | ✅ | See `architecture-decisions.md` § Code Organization Principles | See `implementation-decisions.md` § Code Style and Conventions |
| `[REQ:USABILITY]` | Usability | P1 | ✅ | See `architecture-decisions.md` § Output Formatting | See `implementation-decisions.md` § Dual Printf/Template Formatting |

### Immutable Requirements (Major Version Change Required)

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| `[REQ:IMMUTABLE_ARCHIVE_NAMING]` | Archive Naming Convention | P0 | ✅ | See `architecture-decisions.md` § Archive Format | See `implementation-decisions.md` § ZIP Archive Format Implementation |
| `[REQ:IMMUTABLE_FILE_BACKUP_NAMING]` | File Backup Naming Convention | P0 | ✅ | See `architecture-decisions.md` § File Operations Architecture | See `implementation-decisions.md` § Atomic File Operations |
| `[REQ:IMMUTABLE_DIRECTORY_OPERATIONS]` | Directory Operations | P0 | ✅ | See `architecture-decisions.md` § File Operations Architecture | See `implementation-decisions.md` § File Operations Implementation |
| `[REQ:IMMUTABLE_FILE_BACKUP_OPERATIONS]` | File Backup Operations | P0 | ✅ | See `architecture-decisions.md` § File Operations Architecture | See `implementation-decisions.md` § Atomic File Operations |
| `[REQ:IMMUTABLE_FILE_EXCLUSION]` | File Exclusion | P0 | ✅ | See `architecture-decisions.md` § Exclusion Patterns Architecture | See `implementation-decisions.md` § Exclusion Patterns Implementation |
| `[REQ:IMMUTABLE_GIT_INTEGRATION]` | Git Integration | P1 | ✅ | See `architecture-decisions.md` § Git Integration | See `implementation-decisions.md` § Git Command-line Integration |
| `[REQ:IMMUTABLE_ARCHIVE_VERIFICATION]` | Archive Verification | P1 | ✅ | See `architecture-decisions.md` § Verification Architecture | See `implementation-decisions.md` § Verification Implementation |
| `[REQ:IMMUTABLE_ERROR_HANDLING]` | Error Handling | P0 | ✅ | See `architecture-decisions.md` § Error Handling Strategy | See `implementation-decisions.md` § Structured Error Handling |
| `[REQ:IMMUTABLE_CODE_QUALITY]` | Code Quality Standards | P0 | ✅ | See `architecture-decisions.md` § Code Organization Principles | See `implementation-decisions.md` § Code Style and Conventions |
| `[REQ:IMMUTABLE_BUILD_SYSTEM]` | Build System | P0 | ✅ | See `architecture-decisions.md` § Build and Distribution | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:IMMUTABLE_OUTPUT_FORMATTING]` | Output Formatting | P0 | ✅ | See `architecture-decisions.md` § Output Formatting | See `implementation-decisions.md` § Dual Printf/Template Formatting |
| `[REQ:IMMUTABLE_TEMPLATE_FORMATTING]` | Template Formatting | P0 | ✅ | See `architecture-decisions.md` § Output Formatting | See `implementation-decisions.md` § Dual Printf/Template Formatting |
| `[REQ:IMMUTABLE_CLI_COMMANDS]` | CLI Commands Structure | P0 | ✅ | See `architecture-decisions.md` § CLI Commands Architecture | See `implementation-decisions.md` § CLI Framework Implementation |
| `[REQ:IMMUTABLE_CONFIGURATION_DEFAULTS]` | Configuration Defaults | P0 | ✅ | See `architecture-decisions.md` § Configuration System | See `implementation-decisions.md` § Configuration Structure |
| `[REQ:IMMUTABLE_PLATFORM_COMPATIBILITY]` | Platform Compatibility | P0 | ✅ | See `architecture-decisions.md` § Language and Runtime | See `implementation-decisions.md` § File Operations Implementation |
| `[REQ:IMMUTABLE_GLOBAL_OPTIONS]` | Global Options | P0 | ✅ | See `architecture-decisions.md` § CLI Commands Architecture | See `implementation-decisions.md` § CLI Framework Implementation |
| `[REQ:IMMUTABLE_RESOURCE_MANAGEMENT]` | Resource Management | P0 | ✅ | See `architecture-decisions.md` § Resource Management | See `implementation-decisions.md` § Resource Management with Cleanup |
| `[REQ:IMMUTABLE_PERFORMANCE]` | Performance | P1 | ✅ | See `architecture-decisions.md` § Performance Architecture | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:IMMUTABLE_FEATURE_PRESERVATION]` | Feature Preservation Rules | P0 | ✅ | See `architecture-decisions.md` § Code Organization Principles | See `implementation-decisions.md` § Code Style and Conventions |
| `[REQ:IMMUTABLE_TESTING_INFRASTRUCTURE]` | Testing Infrastructure | P0 | ✅ | See `architecture-decisions.md` § Testing Strategy | See `implementation-decisions.md` § Testing Implementation |

### Incomplete Requirements

| Token | Requirement | Priority | Status | Architecture | Implementation |
|-------|------------|----------|--------|--------------|----------------|
| `[REQ:OUT_002]` | Enhanced Command Output with File Statistics | P0 | ⏳ | See `architecture-decisions.md` § File Statistics Architecture | See `implementation-decisions.md` § File Statistics Implementation |
| `[REQ:LINT_001]` | Code Linting Compliance | P1 | ⏳ | See `architecture-decisions.md` § Code Organization Principles | See `implementation-decisions.md` § Code Style and Conventions |
| `[REQ:DOC_015]` | Unicode to Semantic Token Mapping | P1 | ⏳ | See `architecture-decisions.md` § Documentation Architecture | See `implementation-decisions.md` § Documentation Implementation |
| `[REQ:DOC_016]` | AI-First Comprehensive Token System | P0 | ⏳ | See `architecture-decisions.md` § Token System Architecture | See `implementation-decisions.md` § Token System Implementation |
| `[REQ:COV_003]` | Selective Coverage Reporting | P2 | ⏳ | See `architecture-decisions.md` § Testing Strategy | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:CICD_001]` | AI-First Development Optimization | P2 | ⏳ | See `architecture-decisions.md` § CI/CD Pipeline Architecture | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:DOC_011]` | Token Validation Integration for AI Assistants | P1 | ⏳ | See `architecture-decisions.md` § AI-First Documentation Architecture | See `implementation-decisions.md` § Testing Implementation |
| `[REQ:DOC_013]` | AI-First Documentation and Code Maintenance | P2 | ⏳ | See `architecture-decisions.md` § AI-First Documentation Architecture | See `implementation-decisions.md` § Code Style and Conventions |

---

## Detailed Requirements

### Core Functionality

### [REQ:CODE_QUALITY] Code Quality and Linting Requirements

**Priority: P0 (Critical)**

- **Description**: All code must pass automated quality checks (linting and testing) before successful build. All errors must be properly handled with no unhandled errors allowed.
- **Rationale**: Ensures code quality standards and prevents bugs through automated validation
- **Satisfaction Criteria**:
  - All code passes automated linter checks
  - All function return values are checked
  - All operations handle errors appropriately
  - Build system enforces quality gates before compilation
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Code Organization Principles [ARCH:CODE_ORGANIZATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ✅ Implemented

### [REQ:RESOURCE_MANAGEMENT] Resource Management Requirements

**Priority: P0 (Critical)**

- **Description**: All temporary resources must be cleaned up. Archive and backup operations must be atomic. Operations must recover from panics. Resource management must be thread-safe.
- **Rationale**: Prevents resource leaks and ensures system reliability
- **Satisfaction Criteria**:
  - No temporary resources remain after operations
  - Archive and backup creation uses atomic operations
  - Panic recovery prevents resource leaks
  - Thread-safe resource tracking
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Resource Management [ARCH:RESOURCE_MANAGEMENT]
- **Implementation**: See `implementation-decisions.md` § Resource Management with Cleanup [IMPL:RESOURCE_MANAGER], § Atomic File Operations [IMPL:ATOMIC_OPS]

**Status**: ✅ Implemented

### [REQ:ERROR_HANDLING] Enhanced Error Handling Requirements

**Priority: P0 (Critical)**

- **Description**: Use structured errors with status codes and operation context. Enhanced disk space error detection. Operation context support for better debugging.
- **Rationale**: Consistent error handling with machine-readable status codes and enhanced debugging
- **Satisfaction Criteria**:
  - All operations return structured errors with status codes
  - Disk space errors detected reliably
  - Error messages include operation context
  - Error formatting supports operation context
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Error Handling Strategy [ARCH:ERROR_HANDLING]
- **Implementation**: See `implementation-decisions.md` § Structured Error Handling [IMPL:STRUCTURED_ERRORS]

**Status**: ✅ Implemented

### [REQ:CONTEXT_SUPPORT] Context Support Requirements

**Priority: P0 (Critical)**

- **Description**: All long-running operations must support cancellation and timeouts. Operations must check for cancellation at multiple points. File operations support context cancellation.
- **Rationale**: Enables operation timeouts and graceful cancellation for better user experience
- **Satisfaction Criteria**:
  - All operations accept context parameter
  - Operations check for cancellation periodically
  - Timeout support available
  - File operations support context cancellation
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Context-Aware Operations [ARCH:CONTEXT_SUPPORT]
- **Implementation**: See `implementation-decisions.md` § Context-Aware Operations [IMPL:CONTEXT_OPS]

**Status**: ✅ Implemented

### [REQ:FILE_BACKUP] File Backup Requirements

**Priority: P0 (Critical)**

- **Description**: Create backups of single files with comparison. Detect identical backups. List backups for specific files. File operations use atomic operations and context support.
- **Rationale**: Provides individual file backup capability with comparison and listing
- **Satisfaction Criteria**:
  - File backups created with timestamp and note
  - Identical file detection available
  - Backup listing sorted by creation time
  - Atomic operations prevent corruption
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Resource Management [ARCH:RESOURCE_MANAGEMENT]
- **Implementation**: See `implementation-decisions.md` § Atomic File Operations [IMPL:ATOMIC_OPS], § Context-Aware Operations [IMPL:CONTEXT_OPS]

**Status**: ✅ Implemented

### [REQ:OUTPUT_FORMATTING] Output Formatting Requirements

**Priority: P0 (Critical)**

- **Description**: Centralized formatting for all application output. Template-based formatting with data extraction. Support for multiple formatting syntaxes.
- **Rationale**: Consistent output formatting with rich template support
- **Satisfaction Criteria**:
  - All output uses configurable format strings
  - Template formatting supports data extraction
  - Multiple formatting approaches available
  - Named placeholders supported
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Output Formatting [ARCH:OUTPUT_FORMATTING]
- **Implementation**: See `implementation-decisions.md` § Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING]

**Status**: ✅ Implemented

### [REQ:TEMPLATE_FORMATTING] Template Formatting Requirements

**Priority: P0 (Critical)**

- **Description**: Advanced template processing with named placeholders and data extraction. Support for multiple template syntaxes. Template methods for all operations.
- **Rationale**: Rich template formatting with data extraction capabilities
- **Satisfaction Criteria**:
  - Template processing supports multiple syntaxes
  - Data extraction provides named groups
  - Template methods available for all operations
  - Graceful degradation on template errors
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Output Formatting [ARCH:OUTPUT_FORMATTING]
- **Implementation**: See `implementation-decisions.md` § Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING]

**Status**: ✅ Implemented

### [REQ:CONFIGURATION] Configuration Management Requirements

**Priority: P0 (Critical)**

- **Description**: Configuration stored in files with configurable discovery. Environment variable support. Default search path with precedence rules. Configuration merging with source tracking.
- **Rationale**: Flexible configuration system with multiple sources and precedence
- **Satisfaction Criteria**:
  - Configuration loaded from multiple sources
  - Environment variable support for search paths
  - Precedence: command line > environment > files > defaults
  - Source tracking shows value origins
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Configuration System [ARCH:CONFIG_SYSTEM]
- **Implementation**: See `implementation-decisions.md` § Configuration Structure [IMPL:CONFIG_STRUCT]

**Status**: ✅ Implemented

### [REQ:GIT_INTEGRATION] Git Integration Requirements

**Priority: P1 (Important)**

- **Description**: Automatic Git repository detection. Git branch and commit hash extraction. Git status detection for clean/dirty state. Integration with archive naming.
- **Rationale**: Provides version control information in archives
- **Satisfaction Criteria**:
  - Git repository detected automatically
  - Branch and hash extracted correctly
  - Status detection works for clean/dirty state
  - Archive names include Git info when enabled
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Git Integration [ARCH:GIT_INTEGRATION]
- **Implementation**: See `implementation-decisions.md` § Git Command-line Integration [IMPL:GIT_CLI]

**Status**: ✅ Implemented

### [REQ:ARCHIVE_VERIFICATION] Archive Verification Requirements

**Priority: P1 (Important)**

- **Description**: Archive verification with checksum validation. Verification on creation option. Comprehensive corruption testing framework. Archive repair detection.
- **Rationale**: Ensures archive integrity and data reliability
- **Satisfaction Criteria**:
  - Archives can be verified for integrity
  - Checksum validation available
  - Corruption testing framework supports various types
  - Verification status tracked
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Archive Format [ARCH:ARCHIVE_FORMAT]
- **Implementation**: See `implementation-decisions.md` § ZIP Archive Format Implementation [IMPL:ZIP_FORMAT]

**Status**: ✅ Implemented

## Configuration System Enhancement

### [REQ:CFG_005] Layered Configuration Inheritance Requirements

**Priority: P0 (Critical)**

- **Description**: Configuration files support explicit inheritance declarations using `inherit` field. Inheritance system prevents circular dependencies. Configuration loading processes inheritance chains in correct dependency order. Multiple merge strategies supported. Array configuration fields (such as `exclude_patterns`) default to merge (accumulate) strategy to preserve values from defaults and parent configs when child configs add values.
- **Rationale**: Enables hierarchical configuration management for complex project structures. Array fields accumulating by default ensures that default values (like `.git/`, `vendor/`) are preserved when users add local values.
- **Satisfaction Criteria**:
  - Configuration files can inherit from other files
  - Circular dependencies detected and prevented
  - Inheritance chains processed in correct order
  - Merge strategies (override, append, prepend, replace, default) work correctly
  - Array fields default to merge (accumulate) strategy when no prefix is specified
  - Array field values from defaults and parent configs are preserved when child configs add values
  - Duplicate values in arrays are deduplicated during merge
  - Order is preserved: defaults/parent values first, then child additions
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Layered Configuration Inheritance [ARCH:CFG_005]
- **Implementation**: See `implementation-decisions.md` § Configuration Structure [IMPL:CONFIG_STRUCT]

**Status**: ✅ Implemented

### [REQ:CFG_006] Complete Configuration Reflection and Visibility Requirements

**Priority: P1 (Important)**

- **Description**: Configuration command displays ALL configuration parameters automatically through reflection-based field discovery. Field discovery handles nested structures, embedded fields, collections, and complex types. Source tracking shows complete inheritance chain with merge strategy attribution. Multiple output formats (table, tree, JSON) supported with comprehensive filtering options.
- **Rationale**: Provides complete configuration visibility without manual maintenance. Enables zero-maintenance configuration inspection where new fields automatically appear. Supports debugging through complete source attribution and inheritance chain visualization.
- **Satisfaction Criteria**:
  - **Automatic Field Discovery**: All configuration fields (100+ fields) discovered automatically using Go reflection without manual maintenance
  - **Type Support**: Correct handling of strings, bools, ints, slices, pointers, structs, nested structures, embedded fields, and collections
  - **Source Tracking**: Complete source tracking showing inheritance chain with environment → inheritance → defaults resolution
  - **Inheritance Integration**: Integration with CFG-005 inheritance system showing merge strategies and override points
  - **Multiple Output Formats**: Support for table format (quick scanning), tree format (hierarchical view), and JSON format (programmatic use)
  - **Filtering Options**: Command-line filtering with --all, --overrides-only, --sources, --filter pattern, and --format flags
  - **Performance**: Configuration inspection completes in <100ms for typical usage, <10ms for cached results
  - **Backward Compatibility**: Existing GetConfigValuesWithSources() function works unchanged
  - **Test Coverage**: Comprehensive test suite with >95% coverage including edge cases
- **Validation Criteria**: 
  - **Test Coverage**: >95% coverage for all CFG-006 functionality including field discovery, source attribution, display formatting, filtering, and performance
  - **Performance Benchmarks**: Field discovery <50ms for cache miss, <10ms for cache hit; config command <100ms end-to-end; single field access <10ms
  - **Edge Case Testing**: Anonymous embedded fields, unexported field exclusion, interface{} handling, circular reference prevention, malformed struct handling, depth limitation, error recovery
  - **Integration Testing**: CFG-005 inheritance integration, source tracking accuracy, display formatting validation, filtering functionality
  - **Performance Validation**: Reflection result caching (60%+ improvement), lazy source evaluation, incremental resolution, memory allocation efficiency, concurrent access safety
  - See `architecture-decisions.md` § Configuration Reflection Architecture [ARCH:CFG_006] and `implementation-decisions.md` § Configuration Reflection Implementation [IMPL:CFG_006] for detailed testing approach
- **Architecture**: See `architecture-decisions.md` § Configuration Reflection Architecture [ARCH:CFG_006]
- **Implementation**: See `implementation-decisions.md` § Configuration Reflection Implementation [IMPL:CFG_006]

**Status**: ✅ Implemented

### [REQ:TEST_EXCLUDE_MERGE] Exclude Patterns Merge Testing Requirements

**Priority: P1 (Important)**

- **Description**: Test scenario that verifies `exclude_patterns` are correctly merged from multiple configuration files and that the `config` command accurately displays the source of merged exclude patterns. Test must validate that patterns from multiple config files (e.g., default config and local `.bkpdir.yml`) are properly accumulated and that source tracking shows all contributing sources.
- **Rationale**: Ensures configuration merging works correctly for exclude patterns and that users can debug configuration issues by seeing where each pattern comes from. Critical for validating the fix to exclude pattern merging behavior.
- **Satisfaction Criteria**:
  - Test creates scenario with default config and local config file
  - Test verifies exclude patterns from both sources are merged (accumulated, not replaced)
  - Test verifies merged patterns contain all patterns from all sources
  - Test verifies `config` command shows correct source attribution (not just "default")
  - Test validates that patterns are deduplicated during merge
  - Test validates that order is preserved (defaults first, then local additions)
- **Validation Criteria**: 
  - Test passes with merged exclude patterns containing patterns from all sources
  - Test validates config command output shows source as config file path (not "default") when local patterns are present
  - Test validates inheritance chain tracking shows all contributing files
  - Test covers edge cases: empty arrays, duplicate patterns, pattern order
- **Architecture**: See `architecture-decisions.md` § Configuration Testing Architecture [ARCH:TEST_EXCLUDE_MERGE]
- **Implementation**: See `implementation-decisions.md` § Exclude Patterns Merge Testing [IMPL:TEST_EXCLUDE_MERGE]

**Status**: ✅ Implemented

**Cross-References**: [REQ:CONFIGURATION] (Configuration Management), [REQ:CFG_006] (Configuration Reflection and Visibility), [REQ:IMMUTABLE_FILE_EXCLUSION] (File Exclusion Requirements)


## Non-Functional Requirements

### [REQ:PERFORMANCE] Performance Requirements

**Priority: P1 (Important)**

- **Description**: Operations must complete within reasonable time limits. Configuration loading must be responsive. Reflection results cached for performance. Lazy source evaluation for displayed fields only.
- **Rationale**: Ensures system responsiveness and good user experience
- **Satisfaction Criteria**:
  - Archive creation completes within acceptable time
  - Configuration loading remains responsive
  - Reflection caching reduces overhead
  - Lazy evaluation improves performance
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Language and Runtime [ARCH:LANGUAGE_SELECTION]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ✅ Implemented

### [REQ:RELIABILITY] Reliability Requirements

**Priority: P0 (Critical)**

- **Description**: System must handle errors gracefully. Resource cleanup must be error-resilient. Panic recovery prevents resource leaks. Thread-safe operations for concurrent access.
- **Rationale**: Ensures system stability and prevents data loss
- **Satisfaction Criteria**:
  - Errors handled without crashing
  - Resource cleanup continues on failures
  - Panic recovery prevents leaks
  - Thread-safety verified
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Resource Management [ARCH:RESOURCE_MANAGEMENT], § Error Handling Strategy [ARCH:ERROR_HANDLING]
- **Implementation**: See `implementation-decisions.md` § Resource Management with Cleanup [IMPL:RESOURCE_MANAGER], § Structured Error Handling [IMPL:STRUCTURED_ERRORS]

**Status**: ✅ Implemented

### [REQ:MAINTAINABILITY] Maintainability Requirements

**Priority: P1 (Important)**

- **Description**: Code must follow Go standards. Documentation must be comprehensive. Tests must be maintainable. Configuration system must be extensible.
- **Rationale**: Ensures long-term maintainability and evolution
- **Satisfaction Criteria**:
  - Code follows Go conventions
  - Documentation is complete and up-to-date
  - Tests are easy to understand and modify
  - Configuration system supports new fields automatically
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Code Organization Principles [ARCH:CODE_ORGANIZATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ✅ Implemented

### [REQ:USABILITY] Usability Requirements

**Priority: P1 (Important)**

- **Description**: CLI interface must be intuitive. Error messages must be clear. Configuration must be discoverable. Output must be readable.
- **Rationale**: Ensures good user experience
- **Satisfaction Criteria**:
  - CLI commands are intuitive
  - Error messages are clear and actionable
  - Configuration options are discoverable
  - Output formatting is readable
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Output Formatting [ARCH:OUTPUT_FORMATTING]
- **Implementation**: See `implementation-decisions.md` § Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING]

**Status**: ✅ Implemented

## Incomplete Requirements

### [REQ:OUT_002] Enhanced Command Output with File Statistics Requirements

**Priority: P0 (Critical)**

- **Description**: Enhance backup commands to output formatted file information with stat-like details. All commands must output to stdout a single line about any files generated or existing to satisfy a backup request, using format strings with named replacements for stat information.
- **Rationale**: Provides consistent, informative output across all backup commands with detailed file statistics
- **Satisfaction Criteria**:
  - Full command outputs single line to stdout when archive created successfully with file info
  - Inc command output includes stat-like information, not just path
  - Format string support with named replacements for file statistics ({path}, {size}, {mtime}, etc.)
  - Backward compatibility maintained with existing format strings
  - Both inc and full commands behave consistently
  - Configuration control allows users to customize output format
- **Validation Criteria**: 
  - Unit tests for stat gathering functionality
  - Format string processing tests with named replacement validation
  - Command output behavior tests for inc and full commands
  - Backward compatibility tests ensuring existing format strings work unchanged
  - See `architecture-decisions.md` § File Statistics Architecture [ARCH:FILE_STATISTICS] and `implementation-decisions.md` § File Statistics Implementation [IMPL:FILE_STATISTICS] for detailed testing approach
- **Architecture**: See `architecture-decisions.md` § File Statistics Architecture [ARCH:FILE_STATISTICS]
- **Implementation**: See `implementation-decisions.md` § File Statistics Implementation [IMPL:FILE_STATISTICS]

**Status**: ⏳ In Progress

### [REQ:LINT_001] Code Linting Compliance Requirements

**Priority: P1 (Important)**

- **Description**: Ensure all code passes linting checks and maintains code quality standards. All code must pass automated quality checks before successful build.
- **Rationale**: Ensures code quality standards and prevents bugs through automated validation
- **Satisfaction Criteria**:
  - All code passes automated linter checks
  - Linting integrated into build process
  - CI/CD pipeline enforces linting
  - No linting errors in codebase
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Code Organization Principles [ARCH:CODE_ORGANIZATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ⏳ In Progress

### [REQ:DOC_015] Unicode to Semantic Token Mapping Requirements

**Priority: P1 (Important)**

- **Description**: Replace all unicode icons with AI-first semantic tokens to improve AI assistant comprehension, enable automated validation, and create searchable documentation structure. Maintains backward compatibility with existing priority system ([CRITICAL][HIGH][MEDIUM][LOW]) while replacing decorative unicode icons with machine-readable semantic tokens.
- **Rationale**: Enhances AI assistant navigation and comprehension through searchable semantic tokens
- **Satisfaction Criteria**:
  - Zero unicode icons in documentation headers and structure
  - All validation systems recognize semantic tokens
  - Semantic tokens used consistently across all documentation
  - AI assistants use semantic tokens for navigation and comprehension
  - Backward compatibility preserved with existing priority system
- **Validation Criteria**: 
  - Validation integration with DOC-008 system
  - Semantic token validation for consistency
  - AI assistant comprehension testing
  - See `architecture-decisions.md` § Documentation Architecture and `implementation-decisions.md` § Documentation Implementation for detailed approach
- **Architecture**: See `architecture-decisions.md` § Documentation Architecture
- **Implementation**: See `implementation-decisions.md` § Documentation Implementation

**Status**: ⏳ In Progress (4/5 subtasks complete)

### [REQ:DOC_016] AI-First Comprehensive Token System Requirements

**Priority: P0 (Critical)**

- **Description**: Establish comprehensive token-based traceability for features, architecture decisions, and implementation consistency as a core AI-first principle. Every feature, architecture decision, and implementation must be traced through documentation, source code, tests, and validation layers.
- **Rationale**: Ensures complete traceability from specification to implementation to testing
- **Satisfaction Criteria**:
  - 100% token coverage in source code and tests
  - Cross-layer consistency with token validation passing
  - AI assistant effectiveness >95% feature navigation accuracy using tokens
  - Automated validation integration with existing systems
  - Complete feature traceability from specification to implementation
- **Validation Criteria**: 
  - Token consistency validation across all layers
  - AI assistant navigation accuracy testing
  - Automated validation integration testing
  - See `architecture-decisions.md` § Token System Architecture and `implementation-decisions.md` § Token System Implementation for detailed approach
- **Architecture**: See `architecture-decisions.md` § Token System Architecture
- **Implementation**: See `implementation-decisions.md` § Token System Implementation

**Status**: ⏳ In Progress (1/4 subtasks complete)

### [REQ:COV_003] Selective Coverage Reporting Requirements

**Priority: P2 (Nice-to-Have)**

- **Description**: Implement function-level coverage exclusion and coverage comment directives for granular control over coverage metrics. Support `//coverage:ignore` style annotations for specific code blocks.
- **Rationale**: Provides fine-grained control over coverage reporting while maintaining comprehensive testing
- **Satisfaction Criteria**:
  - Function-level exclusion works correctly
  - Comment directives supported with `//coverage:ignore` style
  - Coverage exception documentation maintained
  - Integration with development workflow complete
  - Coverage visualization tools available
- **Validation Criteria**: See `architecture-decisions.md` § Testing Strategy and `implementation-decisions.md` § Testing Implementation for testing approach
- **Architecture**: See `architecture-decisions.md` § Testing Strategy [ARCH:TESTING_STRATEGY]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ⏳ Not Started

## Immutable Requirements (Major Version Change Required)

These requirements define core behaviors that MUST NOT be changed without a major version bump. These are fundamental specifications that users and other systems depend on.

### [REQ:IMMUTABLE_ARCHIVE_NAMING] Archive Naming Convention (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Archive naming convention is fixed and must not be modified. Format: `{prefix}-{timestamp}-{git_info}-{note}.zip`
  - `prefix`: Optional prefix (default: "bkp")
  - `timestamp`: ISO 8601 format (YYYY-MM-DDTHHmmss)
  - `git_info`: Optional Git branch and short hash (e.g., "main-abc123")
  - `note`: Optional note (default: empty)
  - Example: `bkp-2024-03-20T143022-main-abc123-initial.zip`
- **Rationale**: Users and scripts depend on this naming convention. Changing it would break existing workflows and integrations.
- **Satisfaction Criteria**:
  - Archive names follow the exact format specified
  - All components (prefix, timestamp, git_info, note) are correctly formatted
  - Default prefix is "bkp" when not specified
- **Validation Criteria**: Archive naming tests verify format compliance. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Archive Format [ARCH:ARCHIVE_FORMAT]
- **Implementation**: See `implementation-decisions.md` § ZIP Archive Format Implementation [IMPL:ZIP_FORMAT]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_FILE_BACKUP_NAMING] File Backup Naming Convention (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: File backup naming convention is fixed and must not be modified. Format: `{filename}-{timestamp}[={note}]`
  - `filename`: Original filename without path
  - `timestamp`: YYYY-MM-DD-HH-MM format
  - `note`: Optional note appended with equals sign
  - Examples:
    - `document.txt-2024-03-20-14-30`
    - `config.yml-2024-03-20-14-30=before-changes`
    - `script.sh-2024-03-20-14-30=working-version`
- **Rationale**: Users and scripts depend on this naming convention. Changing it would break existing workflows and integrations.
- **Satisfaction Criteria**:
  - File backup names follow the exact format specified
  - Timestamp format is YYYY-MM-DD-HH-MM
  - Optional notes are appended with equals sign
- **Validation Criteria**: File backup naming tests verify format compliance. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § File Operations Architecture [ARCH:FILE_OPERATIONS]
- **Implementation**: See `implementation-decisions.md` § Atomic File Operations [IMPL:ATOMIC_OPS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_DIRECTORY_OPERATIONS] Directory Operations Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Directory operation rules are fundamental and must not be altered:
  - **Platform Independence**: Use platform-independent path handling
  - **Permission Preservation**: Preserve file permissions and modification time
  - **Atomic Operations**: All operations must be atomic to prevent corruption
  - **Structure Preservation**: Maintain directory structure in archives
  - **Special File Handling**: Handle symbolic links, devices, sockets, etc.
  - **Automatic Creation**: Create archive directories if they don't exist
  - **Path Display**: Display all paths relative to current directory
- **Rationale**: These file operation rules are fundamental behaviors that users depend on. Changing them would break existing functionality.
- **Satisfaction Criteria**:
  - All directory operations use platform-independent paths
  - File permissions and modification times are preserved
  - All operations are atomic
  - Directory structure is maintained in archives
  - Special files are handled correctly
  - Archive directories are created automatically
  - Paths are displayed relative to current directory
- **Validation Criteria**: Directory operation tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § File Operations Architecture [ARCH:FILE_OPERATIONS]
- **Implementation**: See `implementation-decisions.md` § File Operations Implementation [IMPL:FILE_OPERATIONS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_FILE_BACKUP_OPERATIONS] File Backup Operations Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: File backup operation rules are fundamental and must not be altered:
  - **Atomic File Operations**: All file backup operations must be atomic
  - **Identical File Detection**: Must compare files byte-by-byte to detect identical backups
  - **Directory Structure Preservation**: Maintain source file's directory structure in backup path
  - **Timestamp Preservation**: Preserve original file modification times
  - **Permission Preservation**: Preserve file permissions where applicable
  - **Path Resolution**: Handle both absolute and relative file paths
  - **Backup Directory Creation**: Create backup directories if they don't exist
- **Rationale**: These file backup operation rules are fundamental behaviors that users depend on. Changing them would break existing functionality.
- **Satisfaction Criteria**:
  - All file backup operations are atomic
  - Identical files are detected via byte-by-byte comparison
  - Directory structure is preserved in backup paths
  - File modification times are preserved
  - File permissions are preserved where applicable
  - Both absolute and relative paths are handled
  - Backup directories are created automatically
- **Validation Criteria**: File backup operation tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § File Operations Architecture [ARCH:FILE_OPERATIONS]
- **Implementation**: See `implementation-decisions.md` § Atomic File Operations [IMPL:ATOMIC_OPS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_FILE_EXCLUSION] File Exclusion Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: File exclusion requirements are mandatory and must be preserved:
  - **Pattern Matching**: Doublestar glob pattern matching
  - **Configuration**: Configurable exclusion patterns
  - **System Files**: Default exclusions for system files
  - **Case Sensitivity**: Case-sensitive matching
  - **Directory Support**: Directory exclusion support
  - **Precedence**: Pattern precedence rules
- **Rationale**: These exclusion requirements are mandatory behaviors that users depend on. Changing them would break existing configurations.
- **Satisfaction Criteria**:
  - Doublestar glob patterns are supported
  - Exclusion patterns are configurable
  - Default system file exclusions are applied
  - Pattern matching is case-sensitive
  - Directories can be excluded
  - Pattern precedence rules are followed
- **Validation Criteria**: File exclusion tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Exclusion Patterns Architecture [ARCH:EXCLUSION_PATTERNS]
- **Implementation**: See `implementation-decisions.md` § Exclusion Patterns Implementation [IMPL:EXCLUSION_PATTERNS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_GIT_INTEGRATION] Git Integration Requirements (Immutable)

**Priority: P1 (Important - Immutable)**

- **Description**: Git integration requirements are mandatory and must be preserved:
  - **Repository Detection**: Automatic Git repository detection
  - **Archive Naming**: Git information in archive names
  - **Info Extraction**: Git branch and commit hash extraction
  - **State Detection**: Dirty working directory detection
  - **Submodule Support**: Submodule handling
  - **Configuration**: Git configuration integration
- **Rationale**: These Git integration requirements are mandatory behaviors that users depend on. Changing them would break existing workflows.
- **Satisfaction Criteria**:
  - Git repositories are detected automatically
  - Git information is included in archive names when enabled
  - Branch and commit hash are extracted correctly
  - Dirty working directory state is detected
  - Submodules are handled appropriately
  - Git configuration is integrated
- **Validation Criteria**: Git integration tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Git Integration [ARCH:GIT_INTEGRATION]
- **Implementation**: See `implementation-decisions.md` § Git Command-line Integration [IMPL:GIT_CLI]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_ARCHIVE_VERIFICATION] Archive Verification Requirements (Immutable)

**Priority: P1 (Important - Immutable)**

- **Description**: Archive verification requirements are mandatory and must be preserved:
  - **Structure Check**: ZIP structure verification
  - **Checksum Support**: SHA-256 checksum verification
  - **Status Tracking**: Verification status tracking
  - **Atomic Operations**: Atomic verification operations
  - **Result Storage**: Verification result persistence
  - **Error Reporting**: Error reporting for verification failures
- **Rationale**: These verification requirements are mandatory behaviors that users depend on. Changing them would break existing verification workflows.
- **Satisfaction Criteria**:
  - ZIP structure is verified
  - SHA-256 checksums are supported
  - Verification status is tracked
  - Verification operations are atomic
  - Verification results are persisted
  - Verification errors are reported appropriately
- **Validation Criteria**: Archive verification tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Verification Architecture [ARCH:VERIFICATION]
- **Implementation**: See `implementation-decisions.md` § Verification Implementation [IMPL:VERIFICATION]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_ERROR_HANDLING] Error Handling Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Error handling requirements are mandatory and must be preserved:
  - **Structured Errors**: All operations must return structured errors with status codes
  - **Resource Cleanup**: No temporary files or directories may remain after any operation
  - **Panic Recovery**: Application must recover from panics without leaving temporary resources
  - **Context Support**: Long-running operations must support cancellation via context
  - **Enhanced Detection**: Must detect various disk space and permission error conditions
  - **Error Formatting**: Consistent error message formatting
  - **Error Logging**: Comprehensive error logging
- **Rationale**: These error handling requirements are mandatory behaviors that users and scripts depend on. Changing them would break error handling workflows.
- **Satisfaction Criteria**:
  - All operations return structured errors with status codes
  - No temporary resources remain after operations
  - Panics are recovered without leaving temporary resources
  - Long-running operations support context cancellation
  - Disk space and permission errors are detected
  - Error messages are consistently formatted
  - Errors are comprehensively logged
- **Validation Criteria**: Error handling tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Error Handling Strategy [ARCH:ERROR_HANDLING]
- **Implementation**: See `implementation-decisions.md` § Structured Error Handling [IMPL:STRUCTURED_ERRORS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_CODE_QUALITY] Code Quality Standards (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Code quality standards are immutable and must be maintained:
  - **Linting**: All Go code must pass `revive` linter checks
  - **Error Handling**: All errors must be properly handled (no unhandled return values)
  - **Testing**: All code must have comprehensive test coverage
  - **Documentation**: All public functions must be documented
  - **Backward Compatibility**: New features must not break existing functionality
  - **Code Organization**: Consistent code organization and structure
  - **Naming Conventions**: Standard naming conventions
- **Rationale**: These quality standards are immutable behaviors that ensure code quality. Changing them would reduce code quality.
- **Satisfaction Criteria**:
  - All code passes `revive` linter checks
  - All errors are properly handled
  - All code has comprehensive test coverage
  - All public functions are documented
  - New features maintain backward compatibility
  - Code organization is consistent
  - Naming conventions are standard
- **Validation Criteria**: Code quality tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Code Organization Principles [ARCH:CODE_ORGANIZATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_BUILD_SYSTEM] Build System Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Build system requirements are mandatory and must be preserved:
  - **Quality Gates**: Build system must enforce quality standards before compilation
  - **Dependency Management**: Proper ordering of build steps (lint → test → build)
  - **Artifact Management**: Clean target for removing build artifacts
  - **Continuous Integration**: Automated quality checks in CI/CD pipeline
  - **Build Dependencies**: `make build` must depend on `make lint` and `make test` passing
  - **Error Propagation**: Non-zero exit codes from linting or testing must prevent build
  - **CI/CD Commands**: Support for `make ci-lint`, `make ci-test`, `make ci-build`
- **Rationale**: These build system requirements are mandatory behaviors that ensure code quality. Changing them would reduce build quality.
- **Satisfaction Criteria**:
  - Quality gates are enforced before compilation
  - Build steps are ordered correctly (lint → test → build)
  - Clean target removes build artifacts
  - CI/CD pipeline includes automated quality checks
  - `make build` depends on lint and test passing
  - Non-zero exit codes prevent build
  - CI/CD commands are supported
- **Validation Criteria**: Build system tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Build and Distribution [ARCH:BUILD_DISTRIBUTION]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_OUTPUT_FORMATTING] Output Formatting Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Output formatting requirements are mandatory and must be preserved:
  - **Printf-Style Formatting**: All standard output must use printf-style format specifications
  - **Template-Based Formatting**: Must support text/template and placeholder-based formatting
  - **Configuration-Driven**: Format strings and templates must be retrieved from configuration
  - **Text Highlighting**: Must provide means to highlight/format text for structure and meaning
  - **Data Separation**: All user-facing text must be extracted from code into data files
  - **Named Placeholders**: Must support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
  - **Regex Integration**: Must support named regex groups for data extraction
  - **Backward Compatibility**: Default format strings must preserve existing output appearance
  - **Immutable Defaults**: Default format specifications cannot be changed without major version bump
- **Rationale**: These output formatting requirements are mandatory behaviors that users and scripts depend on. Changing them would break output parsing.
- **Satisfaction Criteria**:
  - All standard output uses printf-style format specifications
  - Template-based formatting supports text/template and placeholder syntax
  - Format strings and templates are retrieved from configuration
  - Text highlighting is available
  - User-facing text is extracted from code
  - Both template syntaxes are supported
  - Named regex groups are supported
  - Default format strings preserve existing output appearance
  - Default format specifications are immutable
- **Validation Criteria**: Output formatting tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Output Formatting [ARCH:OUTPUT_FORMATTING]
- **Implementation**: See `implementation-decisions.md` § Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_TEMPLATE_FORMATTING] Template Formatting Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Template formatting requirements are mandatory and must be preserved:
  - **Dual Syntax Support**: Must support both Go text/template syntax ({{.name}}) and placeholder syntax (%{name})
  - **Regex Data Extraction**: Must extract named groups from filenames using configurable regex patterns
  - **Graceful Degradation**: Must fall back to placeholder formatting when template processing fails
  - **Operation Context**: Error messages must include operation context for enhanced debugging
  - **Rich Data Display**: Must extract and display rich information from archive and backup filenames
  - **ANSI Color Support**: Must support ANSI color codes and text formatting in templates
  - **Template Methods**: Must provide template methods for all operations (archives, backups, config, errors)
  - **Print Methods**: Must provide direct printing methods for template-formatted output
  - **Error Handling**: Must handle invalid regex patterns and template syntax gracefully
  - **Configuration Integration**: All template strings and regex patterns must be configurable
- **Rationale**: These template formatting requirements are mandatory behaviors that users depend on. Changing them would break template-based output.
- **Satisfaction Criteria**:
  - Both template syntaxes are supported
  - Named regex groups are extracted from filenames
  - Graceful fallback occurs on template processing failure
  - Error messages include operation context
  - Rich data is extracted and displayed
  - ANSI color codes are supported
  - Template methods are available for all operations
  - Direct printing methods are available
  - Invalid patterns and syntax are handled gracefully
  - Template strings and regex patterns are configurable
- **Validation Criteria**: Template formatting tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Output Formatting [ARCH:OUTPUT_FORMATTING]
- **Implementation**: See `implementation-decisions.md` § Dual Printf/Template Formatting [IMPL:DUAL_FORMATTING]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_CLI_COMMANDS] CLI Commands Structure (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: CLI command structure must be preserved:
  1. Create Archive: `bkpdir create` - Create full archive of current directory. Must use atomic operations with automatic cleanup. Must support context cancellation. Output formatting must use configurable format strings.
  2. Create Incremental Archive: `bkpdir create --incremental` - Create incremental archive of changes. Must use atomic operations with automatic cleanup. Must support context cancellation. Output formatting must use configurable format strings.
  3. Create File Backup: `bkpdir backup [FILE_PATH] [NOTE]` - Create backup of single file with comparison. Must use atomic operations with automatic cleanup. Must support context cancellation. Must detect identical files and report appropriately. Output formatting must use configurable format strings.
  4. List Archives: `bkpdir list` - Sort by creation time (most recent first). Display format: `{path} (created: {time})`. Show verification status: [VERIFIED], [FAILED], or [UNVERIFIED]. Output formatting must use configurable printf-style format strings.
  5. List File Backups: `bkpdir --list [FILE_PATH]` - Sort by creation time (most recent first). Display format: `{path} (created: {time})`. Output formatting must use configurable printf-style format strings.
  6. Verify Archive: `bkpdir verify [ARCHIVE_NAME]` - Support `--checksum` flag for checksum verification. Verify ZIP structure and integrity. Store verification results for display in list command. Output formatting must use configurable printf-style format strings.
  7. Display Configuration: `bkpdir config` - Display computed configuration values with name, value, and source. Process configuration files from `BKPDIR_CONFIG` environment variable. Exit after displaying values. Output formatting must use configurable printf-style format strings.
  8. Backward Compatibility Commands: `bkpdir full [NOTE]` (alias for `bkpdir create`), `bkpdir inc [NOTE]` (alias for `bkpdir create --incremental`), `bkpdir --config` (alias for `bkpdir config`).
- **Rationale**: CLI command structure is fundamental to user experience. Changing it would break user workflows and scripts.
- **Satisfaction Criteria**:
  - All commands work as specified
  - Command aliases are preserved
  - Output formatting uses configurable format strings
  - Atomic operations and context cancellation are supported
- **Validation Criteria**: CLI command tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § CLI Commands Architecture [ARCH:CLI_COMMANDS]
- **Implementation**: See `implementation-decisions.md` § CLI Framework Implementation [IMPL:CLI_FRAMEWORK]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_CONFIGURATION_DEFAULTS] Configuration Defaults (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Configuration defaults must never be changed without explicit user override:
  - Configuration discovery uses `BKPDIR_CONFIG` environment variable to specify search path
  - Default configuration search path is hard-coded as `./.bkpdir.yml:~/.bkpdir.yml` (if `BKPDIR_CONFIG` not set)
  - Configuration files are processed in order with earlier files taking precedence
  - Default archive directory: `../.bkpdir` relative to current directory
  - Default backup directory: `../.bkpdir` relative to current directory
  - Default use_current_dir_name: true
  - Default use_current_dir_name_for_files: true
  - Default include_git_info: true
  - Default exclude_patterns: `[".git/", "vendor/"]`
  - Default verification.verify_on_create: false
  - Default verification.checksum_algorithm: "sha256"
  - Default status codes: All status codes default to `0` (success) if not specified
  - Default output format strings: All format strings default to preserve existing output appearance
  - Default template format strings: All template strings default to preserve existing output appearance
  - Default regex patterns: All regex patterns default to support data extraction
- **Rationale**: Configuration defaults are fundamental behaviors that users depend on. Changing them would break existing configurations.
- **Satisfaction Criteria**:
  - All defaults are as specified
  - Defaults can be overridden by user configuration
  - Configuration discovery works as specified
  - Precedence rules are followed
- **Validation Criteria**: Configuration tests verify all defaults. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Configuration System [ARCH:CONFIG_SYSTEM]
- **Implementation**: See `implementation-decisions.md` § Configuration Structure [IMPL:CONFIG_STRUCT]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_PLATFORM_COMPATIBILITY] Platform Compatibility Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Platform support must never be reduced or modified:
  - **OS Support**: Support macOS and Linux systems
  - **File System**: Handle platform-specific file system differences
  - **Permissions**: Preserve file permissions and ownership where applicable
  - **Thread Safety**: Thread-safe operations for concurrent access
  - **Resource Management**: Efficient resource management across platforms
  - **Character Encoding**: Handle platform-specific character encodings
  - **Line Endings**: Handle platform-specific line endings
  - **Path Separators**: Handle platform-specific path separators
- **Rationale**: Platform compatibility is fundamental to user experience. Reducing support would break existing users.
- **Satisfaction Criteria**:
  - macOS and Linux are supported
  - Platform-specific differences are handled
  - File permissions are preserved
  - Operations are thread-safe
  - Resource management is efficient
  - Character encodings are handled
  - Line endings are handled
  - Path separators are handled
- **Validation Criteria**: Platform compatibility tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Language and Runtime [ARCH:LANGUAGE_SELECTION]
- **Implementation**: See `implementation-decisions.md` § File Operations Implementation [IMPL:FILE_OPERATIONS]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_GLOBAL_OPTIONS] Global Options (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Global options must be maintained:
  - Support `--dry-run` flag for previewing archive operations
  - Dry-run must include resource cleanup verification
  - Output formatting must use configurable printf-style format strings
  - Existing flag behavior must be maintained
- **Rationale**: Global options are fundamental to user experience. Changing them would break user workflows.
- **Satisfaction Criteria**:
  - `--dry-run` flag works as specified
  - Resource cleanup is verified in dry-run
  - Output formatting uses configurable format strings
  - Flag behavior is maintained
- **Validation Criteria**: Global options tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § CLI Commands Architecture [ARCH:CLI_COMMANDS]
- **Implementation**: See `implementation-decisions.md` § CLI Framework Implementation [IMPL:CLI_FRAMEWORK]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_RESOURCE_MANAGEMENT] Resource Management Requirements (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Resource management requirements are mandatory and cannot be relaxed:
  - **Automatic Cleanup**: All temporary resources must be cleaned up automatically
  - **Thread Safety**: Resource management must be thread-safe
  - **Atomic Operations**: File operations must use temporary files for atomicity
  - **Leak Prevention**: No resource leaks allowed in any scenario
  - **Error Resilience**: Cleanup must continue even if individual operations fail
  - **Resource Tracking**: Comprehensive resource tracking
  - **Memory Management**: Efficient memory management
  - **Disk Space**: Efficient disk space management
  - **Cleanup Warnings**: Proper cleanup warnings
  - **Resource Limits**: Enforced resource limits
- **Rationale**: Resource management requirements are mandatory behaviors that ensure system reliability. Relaxing them would cause resource leaks.
- **Satisfaction Criteria**:
  - All temporary resources are cleaned up automatically
  - Resource management is thread-safe
  - File operations use temporary files for atomicity
  - No resource leaks occur
  - Cleanup continues on individual operation failures
  - Resource tracking is comprehensive
  - Memory management is efficient
  - Disk space management is efficient
  - Cleanup warnings are proper
  - Resource limits are enforced
- **Validation Criteria**: Resource management tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Resource Management [ARCH:RESOURCE_MANAGEMENT]
- **Implementation**: See `implementation-decisions.md` § Resource Management with Cleanup [IMPL:RESOURCE_MANAGER]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_PERFORMANCE] Performance Requirements (Immutable)

**Priority: P1 (Important - Immutable)**

- **Description**: Performance characteristics must be preserved:
  - **Minimal Overhead**: Resource tracking must have minimal performance impact
  - **Efficient Operations**: File comparison must check length before byte comparison
  - **Scalability**: Must handle large files and many archives efficiently
  - **Memory Management**: Must maintain low memory footprint
  - **CPU Usage**: Enforced CPU usage limits
  - **Disk I/O**: Enforced disk I/O limits
  - **Response Time**: Enforced response time limits
  - **Cleanup Timing**: Efficient resource cleanup timing
  - **Concurrency**: Enforced concurrent operation limits
  - **Verification**: Efficient verification speed
- **Rationale**: Performance characteristics are fundamental to user experience. Degrading them would reduce usability.
- **Satisfaction Criteria**:
  - Resource tracking has minimal performance impact
  - File comparison checks length before byte comparison
  - Large files and many archives are handled efficiently
  - Memory footprint is low
  - CPU usage limits are enforced
  - Disk I/O limits are enforced
  - Response time limits are enforced
  - Resource cleanup timing is efficient
  - Concurrent operation limits are enforced
  - Verification speed is efficient
- **Validation Criteria**: Performance tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Performance Architecture [ARCH:PERFORMANCE]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_FEATURE_PRESERVATION] Feature Preservation Rules (Immutable)

**Priority: P0 (Critical - Immutable)**

- **Description**: Feature preservation rules must be followed:
  1. **New Features**: Must not interfere with existing functionality. Must maintain all current behaviors. Must be optional and not affect existing workflows. Must include automatic resource cleanup. Must support context cancellation where appropriate. Must pass all linting and testing requirements.
  2. **Modifications**: Must preserve all existing command-line interfaces. Must maintain current directory handling behaviors. Must keep existing configuration options. Must not change established archive naming patterns. Must not introduce resource leaks. Must maintain error handling standards. Must preserve atomic operation guarantees. Must maintain Git integration compatibility. Must preserve verification functionality.
  3. **Testing Requirements**: All new code must include tests for existing functionality. Regression tests must verify no existing features are broken. Platform compatibility tests must be maintained. Resource cleanup must be verified in all test scenarios. Context cancellation and timeout handling must be tested. Performance benchmarks must not regress. All code must pass linting before commit. Git integration tests must be maintained. Archive verification tests must be comprehensive.
  4. **Quality Assurance**: Code must pass revive linter with zero warnings. All errors must be properly handled. All public functions must be documented. Test coverage must meet minimum thresholds. No temporary files may remain after any operation. Memory leaks are strictly prohibited. Archive integrity must be guaranteed. Git integration must be reliable.
- **Rationale**: Feature preservation rules ensure backward compatibility and system reliability. Violating them would break existing functionality.
- **Satisfaction Criteria**:
  - New features don't interfere with existing functionality
  - Modifications preserve existing interfaces and behaviors
  - Testing requirements are met
  - Quality assurance standards are maintained
- **Validation Criteria**: Feature preservation tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Code Organization Principles [ARCH:CODE_ORGANIZATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ✅ Implemented (Immutable)

### [REQ:IMMUTABLE_TESTING_INFRASTRUCTURE] Testing Infrastructure Immutable Requirements

**Priority: P0 (Critical - Immutable)**

- **Description**: Testing infrastructure requirements are immutable:
  - **Corruption Type Enumeration Stability**: The 8 corruption types (CRC, Header, Truncate, CentralDir, LocalHeader, Data, Signature, Comment) must remain stable across versions. Test code depends on these specific corruption types for systematic verification testing. New corruption types may be added but existing types cannot be renamed or removed.
  - **Deterministic Corruption Behavior**: Identical seeds must produce identical corruption across versions. CI/CD systems depend on reproducible test results for regression detection. Corruption algorithms cannot be changed in ways that break reproducibility. Exception: Bug fixes that improve correctness may change behavior if documented in release notes.
  - **Safe Testing Guarantees**: Backup/restore functionality must prevent data loss during corruption testing. Testing infrastructure must never risk data integrity. Any changes to backup/restore logic must maintain safety guarantees.
  - **Performance Baseline Stability**: Performance characteristics must not degrade by more than 20% across versions. Testing infrastructure performance affects overall test suite execution time. Baseline: CRC corruption ~763μs, detection ~49μs (±20% acceptable variance). Performance optimizations welcome, but regressions require justification.
- **Rationale**: Testing infrastructure requirements ensure test reliability and reproducibility. Changing them would break test infrastructure.
- **Satisfaction Criteria**:
  - Corruption types remain stable
  - Corruption behavior is deterministic
  - Testing is safe (no data loss)
  - Performance baselines are maintained
- **Validation Criteria**: Testing infrastructure tests verify all requirements. Breaking changes require major version bump.
- **Architecture**: See `architecture-decisions.md` § Testing Strategy [ARCH:TESTING_STRATEGY]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ✅ Implemented (Immutable)

## AI-First Development Requirements

### [REQ:CICD_001] AI-First Development Optimization Requirements

**Priority: P2 (Nice-to-Have)**

- **Description**: CI/CD pipeline must be optimized for AI-first development workflows. All CI/CD tasks must be configured with low priority scheduling. Pipeline configurations must assume zero human developer intervention. All documentation and code must be maintained for AI assistant comprehension. CI/CD outputs must use standardized icon system. Pipeline must integrate with implementation token validation system.
- **Rationale**: Enables efficient AI-assisted development workflows without blocking human developers
- **Satisfaction Criteria**:
  - CI/CD tasks configured with low priority scheduling
  - Pipeline assumes zero human intervention
  - Documentation and code maintained for AI comprehension
  - CI/CD outputs use standardized icon system
  - Pipeline integrates with token validation system
  - Background task execution prevents blocking AI workflows
  - Automated quality gates without human approval
  - AI-optimized error reporting and status messaging
  - Self-documenting pipeline configurations
  - Streamlined feedback loops for iterative AI development
- **Validation Criteria**: CI/CD pipeline tests verify all requirements. Integration tests verify AI workflow compatibility.
- **Architecture**: See `architecture-decisions.md` § CI/CD Pipeline Architecture [ARCH:CICD_PIPELINE]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ⏳ In Progress

### [REQ:DOC_011] Token Validation Integration for AI Assistants Requirements

**Priority: P1 (Important)**

- **Description**: AI workflow validation hooks must integrate seamlessly with existing validation framework. Pre-submission validation must prevent non-compliant changes from being submitted. Validation feedback must be optimized for AI assistant comprehension and remediation. Integration must provide zero-friction validation without disrupting AI workflows. All validation processes must support automated execution by AI assistants.
- **Rationale**: Enables AI assistants to validate changes automatically before submission, reducing errors and improving code quality
- **Satisfaction Criteria**:
  - AI workflow validation hooks integrate with existing validation framework
  - Pre-submission validation prevents non-compliant changes
  - Validation feedback optimized for AI comprehension
  - Zero-friction validation without disrupting workflows
  - Automated execution support for AI assistants
  - Error messages structured for AI parsing
  - Validation feedback includes clear remediation steps
  - Error formatting is consistent and machine-readable
  - Validation results include file and line references
  - Error categorization enables prioritization
- **Validation Criteria**: Validation integration tests verify all requirements. AI workflow tests verify automated execution.
- **Architecture**: See `architecture-decisions.md` § AI-First Documentation Architecture [ARCH:AI_DOCUMENTATION]
- **Implementation**: See `implementation-decisions.md` § Testing Implementation [IMPL:TESTING]

**Status**: ⏳ In Progress

### [REQ:DOC_013] AI-First Documentation and Code Maintenance Requirements

**Priority: P2 (Nice-to-Have)**

- **Description**: All documentation must be written primarily for AI assistant comprehension. Code comments and implementation tokens must follow standardized AI-readable formats. Documentation structure must be optimized for AI parsing and navigation. Cross-references must be machine-readable and AI-traversable. Maintenance workflows must be executable by AI assistants without human intervention.
- **Rationale**: Enables AI assistants to understand and maintain codebase effectively, reducing human maintenance burden
- **Satisfaction Criteria**:
  - Documentation written primarily for AI comprehension
  - Code comments and tokens follow standardized formats
  - Documentation structure optimized for AI parsing
  - Cross-references are machine-readable and AI-traversable
  - Maintenance workflows executable by AI assistants
  - Consistent formatting patterns for AI comprehension
  - Implementation tokens integrate with icon standardization
  - Code comments prioritize AI understanding
  - Documentation includes explicit cross-reference links
  - Text formatting follows AI-friendly markup conventions
- **Validation Criteria**: Documentation validation tests verify all requirements. AI comprehension tests verify effectiveness.
- **Architecture**: See `architecture-decisions.md` § AI-First Documentation Architecture [ARCH:AI_DOCUMENTATION]
- **Implementation**: See `implementation-decisions.md` § Code Style and Conventions [IMPL:CODE_STYLE]

**Status**: ⏳ In Progress

## Future Enhancements (Out of Scope)

The following features are documented but marked as future enhancements:
- Advanced performance monitoring
- Real-time validation feedback integration (DOC-012 completed)
- Enhanced AI assistant optimization
- Advanced coverage controls (COV-003 planned but not started)

These may be considered for future iterations but are not required for the initial implementation.
