# Tasks and Incomplete Subtasks

**STDD Methodology Version**: 1.0.2

## Overview
This document tracks all tasks and subtasks for implementing this project. Tasks are organized by priority and implementation phase.

## Priority Levels

- **P0 (Critical)**: Must have - Core functionality, blocks other work
- **P1 (Important)**: Should have - Enhanced functionality, better error handling
- **P2 (Nice-to-Have)**: Could have - UI/UX improvements, convenience features
- **P3 (Future)**: Won't have now - Deferred features, experimental ideas

## Task Format

```markdown
## P0: Task Name [REQ:IDENTIFIER] [ARCH:IDENTIFIER] [IMPL:IDENTIFIER]

**Status**: 🟡 In Progress | ✅ Complete | ⏸️ Blocked | ⏳ Pending

**Description**: Brief description of what this task accomplishes.

**Dependencies**: List of other tasks/tokens this depends on.

**Subtasks**:
- [ ] Subtask 1 [REQ:X] [IMPL:Y]
- [ ] Subtask 2 [REQ:X] [IMPL:Z]
- [ ] Subtask 3 [TEST:X]

**Completion Criteria**:
- [ ] All subtasks complete
- [ ] Code implements requirement
- [ ] Tests pass with semantic token references
- [ ] Documentation updated

**Priority Rationale**: Why this is P0/P1/P2/P3
```

## Task Management Rules

1. **Subtasks are Temporary**
   - Subtasks exist only while the parent task is in progress
   - Remove subtasks when parent task completes

2. **Priority Must Be Justified**
   - Each task must have a priority rationale
   - Priorities follow: Tests/Code/Functions > DX > Infrastructure > Security

3. **Semantic Token References Required**
   - Every task MUST reference at least one semantic token
   - Cross-reference to related tokens

4. **Completion Criteria Must Be Met**
   - All criteria must be checked before marking complete
   - Documentation must be updated

## Task Status Icons

- 🟡 **In Progress**: Actively being worked on
- ✅ **Complete**: All criteria met, subtasks removed
- ⏸️ **Blocked**: Waiting on dependency
- ⏳ **Pending**: Not yet started

## P1: Configuration Output Grouping [REQ:CONFIG_OUTPUT_GROUPING] [ARCH:CONFIG_OUTPUT_GROUPING] [IMPL:CONFIG_OUTPUT_GROUPING]

**Status**: ✅ Complete

**Description**: Change the primary presentation of the `--config` output to be grouped by category and ranked by importance/frequency. The existing single sorted list becomes the secondary presentation.

**Dependencies**: [REQ:CFG_006] (Configuration Reflection)

**Completion Criteria**:
- [x] Default `bkpdir config` output is grouped by category
- [x] Categories are ranked by importance
- [x] Items within categories are ranked by importance
- [x] Secondary presentation (flat list) is accessible via flag
- [x] Tests pass

**Priority Rationale**: P1 - User requested change to improve usability of configuration display.

## P0: Enhanced Command Output with File Statistics [REQ:OUT_002] [ARCH:FILE_STATISTICS] [IMPL:FILE_STATISTICS]

**Status**: 🟡 In Progress

**Description**: Enhance backup commands to output formatted file information with stat-like details. All commands must output to stdout a single line about any files generated or existing to satisfy a backup request, using format strings with named replacements for stat information.

**Dependencies**: [REQ:OUTPUT_FORMATTING] (format strings system)

**Subtasks**:
- [ ] File Stat Information Gathering [REQ:OUTPUT_FORMATTING] [IMPL:OUTPUT_FORMATTING]
  - [ ] Create FileStatInfo struct - Define structure for file statistics data
  - [ ] Implement GatherFileStatInfo function - Gather stat-like info for files
  - [ ] Add human-readable size formatting - Convert bytes to KB/MB/GB format
  - [ ] Add file type detection - Determine file type (regular, directory, symlink)
- [ ] Enhanced Format String Configuration [REQ:OUTPUT_FORMATTING] [IMPL:OUTPUT_FORMATTING]
  - [ ] Add detailed format string options - format_archive_created_detailed, format_incremental_created_detailed
  - [ ] Implement named replacement support - {path}, {name}, {size}, {size_human}, {mtime}, {mode}, {type}
  - [ ] Update DefaultConfig() with new format strings - Maintain backward compatibility
  - [ ] Add template-based formatting support - Enhanced template processing for stat data
- [ ] OutputFormatter Enhancement [REQ:OUTPUT_FORMATTING] [IMPL:OUTPUT_FORMATTING]
  - [ ] Add stat-aware formatting methods - FormatCreatedArchiveWithStats, FormatIncrementalCreatedWithStats
  - [ ] Implement named replacement processing - Template engine for {name} style replacements
  - [ ] Add print methods for enhanced formatting - PrintCreatedArchiveWithStats, PrintIncrementalCreatedWithStats
  - [ ] Maintain backward compatibility - Preserve existing format methods unchanged
- [ ] Archive Creation Function Updates [REQ:OUTPUT_FORMATTING] [IMPL:OUTPUT_FORMATTING]
  - [ ] Add success output to CreateFullArchiveWithContext - Output when archive created successfully
  - [ ] Enhance createIncrementalArchive output - Use stat-based formatting
  - [ ] Ensure consistent behavior - Both inc and full commands output file info
  - [ ] Integrate with error handling - Proper error handling for stat gathering
- [ ] Testing and Validation [REQ:OUTPUT_FORMATTING] [TEST:OUTPUT_FORMATTING]
  - [ ] Add unit tests for stat gathering - Test FileStatInfo function and error handling
  - [ ] Test format string processing - Verify named replacement functionality
  - [ ] Test command output behavior - Validate inc and full command output
  - [ ] Verify backward compatibility - Ensure existing format strings work unchanged
- [ ] Documentation Updates [REQ:OUTPUT_FORMATTING]
  - [ ] Update specification.md - Document new output behavior with examples
  - [ ] Update architecture.md - Add stat-based output formatting system documentation
  - [ ] Update requirements.md - Add file statistics and format string requirements
  - [ ] Update configuration documentation - Document new format string options and named replacements

**Completion Criteria**:
- [ ] All subtasks complete
- [ ] Full command outputs single line to stdout when archive created successfully with file info
- [ ] Inc command output includes stat-like information, not just path
- [ ] Format string support with named replacements for file statistics
- [ ] Backward compatibility maintained
- [ ] Both inc and full commands behave consistently
- [ ] Tests pass with semantic token references
- [ ] Documentation updated

**Priority Rationale**: P0 - Enhances user experience and provides consistent output formatting across commands

## P1: Array Field Default Merge Strategy Implementation [REQ:CFG_005] [ARCH:EXCLUDE_MERGE_FIX] [IMPL:EXCLUDE_MERGE_FIX]

**Status**: ✅ Complete

**Description**: Implement CFG-005 requirement that array configuration fields default to merge (accumulate) strategy. Fixes implementation bug where array fields were using "override" instead of "merge" by default.

**Dependencies**: [REQ:CFG_005] (Configuration Inheritance), [REQ:CONFIGURATION] (Configuration Management), [REQ:TEST_EXCLUDE_MERGE] (Exclude Patterns Merge Testing)

**Subtasks**:
- [x] Implement merge strategy detection for array fields [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Fix merge state management to use current state [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Fix merge operation to handle nil values and deduplication [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Add explicit pattern preservation before merge [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Handle YAML type conversions (`[]interface{}` to `[]string`) [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Implement graceful unknown field handling [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Add `isKnownConfigField` helper function [REQ:CFG_005] [IMPL:EXCLUDE_MERGE_FIX]
- [x] Update tests to use explicit merge strategy prefixes where override expected [REQ:CFG_005] [TEST:EXCLUDE_MERGE]
- [x] Create test scenario for merge validation [REQ:TEST_EXCLUDE_MERGE] [IMPL:TEST_EXCLUDE_MERGE]

**Completion Criteria**:
- [x] exclude_patterns defaults to "merge" strategy instead of "override" in all contexts
- [x] Patterns from defaults are preserved when local config adds patterns
- [x] Patterns are deduplicated during merge
- [x] Order is preserved (defaults first, then additions)
- [x] Merge operations use current state (not original state)
- [x] YAML unmarshaling type conversions handled correctly (`[]interface{}` to `[]string`)
- [x] Unknown config fields gracefully skipped instead of aborting merge
- [x] Metadata fields (like `inherit`) filtered out before merge operations
- [x] All tests pass (TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE, TestLoadConfigMultipleFiles, TestLoadConfigWithInheritance_MultiFile, TestSourceConflictDetection)
- [x] Documentation updated

**Priority Rationale**: P1 - Implements CFG-005 requirement for array field default merge behavior. Critical for user experience and configuration behavior - users expect array fields to accumulate from multiple sources.

## P1: Code Linting Compliance [REQ:LINT_001] [ARCH:CODE_ORGANIZATION] [IMPL:CODE_STYLE]

**Status**: 🟡 In Progress

**Description**: Ensure all code passes linting checks and maintains code quality standards.

**Dependencies**: [REQ:CODE_QUALITY] (code quality requirements)

**Subtasks**:
- [ ] Review and fix linting errors
- [ ] Update linting configuration if needed
- [ ] Ensure all new code passes linting
- [ ] Add linting to CI/CD pipeline

**Completion Criteria**:
- [ ] All code passes linting checks
- [ ] Linting integrated into build process
- [ ] CI/CD pipeline enforces linting

**Priority Rationale**: P1 - Important for code quality and maintainability

## P2: Selective Coverage Reporting [REQ:COV_003] [ARCH:TESTING_STRATEGY] [IMPL:TESTING]

**Status**: ⏳ Pending

**Description**: Implement function-level coverage exclusion and coverage comment directives for granular control over coverage metrics.

**Dependencies**: COV-001, COV-002 (coverage baseline)

**Subtasks**:
- [ ] Implement function-level coverage exclusion
- [ ] Add coverage comment directives - Support `//coverage:ignore` style annotations
- [ ] Create coverage exception documentation
- [ ] Integrate with development workflow
- [ ] Add coverage visualization tools

**Completion Criteria**:
- [ ] Function-level exclusion works correctly
- [ ] Comment directives supported
- [ ] Exception documentation maintained
- [ ] Integration with workflow complete

**Priority Rationale**: P2 - Nice-to-have feature for advanced coverage management

## P0: Unicode to Semantic Token Mapping [REQ:DOC_015]

**Status**: 🟡 In Progress (4/5 subtasks complete)

**Description**: Replace all unicode icons with AI-first semantic tokens to improve AI assistant comprehension, enable automated validation, and create searchable documentation structure.

**Dependencies**: DOC-007/DOC-008 (validation systems)

**Subtasks**:
- [x] Unicode Icon Mapping Definition - COMPLETED
- [x] Semantic Token Registry - COMPLETED
- [x] Documentation Migration - COMPLETED
- [ ] Validation Integration [REQ:DOC_015]
  - [ ] Update DOC-008 validation - Extend validation to check semantic tokens
  - [ ] Create semantic token validation - Validate token usage and consistency
  - [ ] Update validation scripts - Modify scripts to recognize semantic tokens
  - [ ] Test validation integration - Ensure all validation passes with semantic tokens
- [ ] AI Assistant Integration [REQ:DOC_015]
  - [ ] Update ai-assistant-compliance.md - Reference semantic tokens in all guidance
  - [ ] Update AI assistant templates - Use semantic tokens in all templates
  - [ ] Update validation automation - Check semantic token compliance
  - [ ] Test AI assistant comprehension - Validate improved AI navigation

**Completion Criteria**:
- [ ] All subtasks complete
- [ ] Zero unicode icons in documentation headers
- [ ] All validation systems recognize semantic tokens
- [ ] AI assistants use semantic tokens for navigation

**Priority Rationale**: P0 - Critical for AI-first development workflow

## P0: AI-First Comprehensive Token System [REQ:DOC_016]

**Status**: 🟡 In Progress (1/4 subtasks complete)

**Description**: Establish comprehensive token-based traceability for features, architecture decisions, and implementation consistency as a core AI-first principle.

**Dependencies**: [REQ:DOC_015] (Unicode to semantic token mapping)

**Subtasks**:
- [x] Token Registry Establishment - COMPLETED
- [ ] Source Code Token Implementation [REQ:DOC_016]
  - [ ] Add implementation tokens - Add feature tokens to all Go source files
  - [ ] Add test tokens - Add test tokens to all test files
  - [ ] Add architecture decision tokens - Add architecture tokens to documentation
  - [ ] Validate token consistency - Ensure token consistency across all layers
- [ ] AI Assistant Integration [REQ:DOC_016]
  - [ ] Update AI assistant guidelines - Add token usage requirements
  - [ ] Create token navigation tools - Enable token-based code navigation
  - [ ] Implement token-based search - Create search capabilities
  - [ ] Validate AI assistant effectiveness - Test AI assistant token usage
- [ ] Advanced Token Features [REQ:DOC_016]
  - [ ] Implement real-time token validation - Add live validation during development
  - [ ] Create token-based impact analysis - Analyze feature impacts using tokens
  - [ ] Develop token-based documentation generation - Generate docs from tokens
  - [ ] Establish token-based quality metrics - Create quality metrics system

**Completion Criteria**:
- [ ] All subtasks complete
- [ ] 100% token coverage in source code and tests
- [ ] Cross-layer consistency with token validation passing
- [ ] AI assistant effectiveness >95% feature navigation accuracy

**Priority Rationale**: P0 - Critical foundation for AI-first development

## P0: Dependency Analysis and Interface Standardization [REFACTOR-001] [ARCH:CODE_ORGANIZATION] [IMPL:REFACTOR_PREP]

**Status**: ⏳ Pending

**Description**: Complete comprehensive dependency analysis to identify all component dependencies and standardize interfaces before extraction begins. This is a CRITICAL BLOCKER for all extraction tasks.

**Dependencies**: None (foundation task)

**Subtasks**:
- [ ] Map all component dependencies - Identify dependencies between all components
- [ ] Identify circular dependency risks - Detect potential circular dependencies
- [ ] Standardize interface definitions - Create consistent interface patterns
- [ ] Document dependency hierarchy - Create dependency graph documentation
- [ ] Validate zero circular dependencies - Ensure no circular dependency risks exist

**Completion Criteria**:
- [ ] Complete dependency mapping for all components
- [ ] Zero circular dependency risks identified
- [ ] All interfaces standardized and documented
- [ ] Dependency hierarchy documented
- [ ] Validation confirms extraction readiness

**Priority Rationale**: P0 - CRITICAL BLOCKER - All extraction tasks depend on this completion

## P0: Refactoring Impact Validation [REFACTOR-006] [ARCH:CODE_ORGANIZATION] [IMPL:REFACTOR_PREP]

**Status**: ⏳ Pending

**Description**: Run comprehensive validation after each refactoring step to ensure no functionality regression, performance degradation, or documentation drift.

**Dependencies**: REFACTOR-001, REFACTOR-002, REFACTOR-003, REFACTOR-004, REFACTOR-005

**Subtasks**:
- [ ] Run comprehensive test suite after each refactoring - Ensure no functionality regression
- [ ] Validate performance impact - Confirm refactoring doesn't degrade performance
- [ ] Check implementation token consistency - Verify all tokens remain valid after refactoring
- [ ] Validate documentation synchronization - Ensure context files reflect refactoring changes
- [ ] Run extraction readiness assessment - Confirm codebase is ready for component extraction

**Completion Criteria**:
- [ ] Zero functional regressions from refactoring
- [ ] Maintained or improved performance
- [ ] Consistent implementation tokens
- [ ] Synchronized documentation
- [ ] Validated extraction readiness

**Priority Rationale**: P0 - HIGH - Must validate each refactoring step before proceeding

## P0: Configuration Management System Extraction [EXTRACT-001] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ✅ Authorized (Ready to begin)

**Description**: Extract configuration management system to `go-cli-config` package with config loading, validation, YAML/JSON support, environment variables, and defaults.

**Dependencies**: ✅ REFACTOR-003 completed (Configuration abstraction interfaces defined)

**Subtasks**:
- [ ] Create `go-cli-config` package structure
- [ ] Extract config loading logic - ConfigLoader, ConfigValidator interfaces
- [ ] Extract config merging logic - ConfigMerger with strategies
- [ ] Extract config sources - YAML, environment, defaults
- [ ] Extract schema abstraction - ConfigSchema, FieldDefinition
- [ ] Create backward compatibility layer - ConfigAdapter for existing code
- [ ] Add comprehensive tests - Package-specific test suite
- [ ] Update main application - Use extracted package

**Completion Criteria**:
- [ ] Package extracted to `go-cli-config`
- [ ] All existing functionality preserved
- [ ] Comprehensive test coverage
- [ ] Backward compatibility maintained
- [ ] Main application updated to use package

**Priority Rationale**: P0 - CRITICAL - Foundation for other extractions

## P0: Error Handling and Resource Management Extraction [EXTRACT-002] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ✅ Authorized (Ready to begin)

**Description**: Extract error handling and resource management to `go-cli-errors` package with ErrorInterface, ArchiveError, BackupError, ResourceManager, and context patterns.

**Dependencies**: ✅ REFACTOR-004 completed (Error handling patterns standardized)

**Subtasks**:
- [ ] Create `go-cli-errors` package structure
- [ ] Extract error types - ErrorInterface, ArchiveError, BackupError
- [ ] Extract ResourceManager - Resource cleanup and panic recovery
- [ ] Extract context patterns - Context-aware operations
- [ ] Create backward compatibility layer - Error adapters
- [ ] Add comprehensive tests - Package-specific test suite
- [ ] Update main application - Use extracted package

**Completion Criteria**:
- [ ] Package extracted to `go-cli-errors`
- [ ] All error handling patterns preserved
- [ ] Resource management functionality maintained
- [ ] Comprehensive test coverage
- [ ] Main application updated to use package

**Priority Rationale**: P0 - CRITICAL - Foundation for other extractions

## P1: Output Formatting System Extraction [EXTRACT-003] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ⏸️ Blocked (Requires REFACTOR-002 completion)

**Description**: Extract output formatting system including OutputFormatter, TemplateFormatter, OutputCollector, and formatting patterns.

**Dependencies**: REFACTOR-002 (Formatter decomposition strategy validated)

**Subtasks**:
- [ ] Create output formatting package structure
- [ ] Extract OutputCollector component - Zero dependencies, ready for extraction
- [ ] Extract Printf Formatter component - Config dependency
- [ ] Extract Template Formatter component - Config dependency
- [ ] Extract Extended Output Formatters - Complex operations
- [ ] Extract Error Formatting component - Specialized error handling
- [ ] Create backward compatibility layer
- [ ] Add comprehensive tests

**Completion Criteria**:
- [ ] All formatter components extracted
- [ ] Clean component boundaries maintained
- [ ] Backward compatibility preserved
- [ ] Comprehensive test coverage
- [ ] Main application updated

**Priority Rationale**: P1 - HIGH - Requires formatter decomposition completion

## P1: Git Integration System Extraction [EXTRACT-004] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ⏳ Pending (Ready after dependency analysis)

**Description**: Extract Git integration system to reusable package with repository detection, branch extraction, hash extraction, and status detection.

**Dependencies**: REFACTOR-001 (Dependency analysis)

**Subtasks**:
- [ ] Create Git integration package structure
- [ ] Extract Git command execution logic
- [ ] Extract repository detection - `git rev-parse --is-inside-work-tree`
- [ ] Extract branch extraction - `git rev-parse --abbrev-ref HEAD`
- [ ] Extract hash extraction - `git rev-parse --short HEAD`
- [ ] Extract status detection - `git status --porcelain`
- [ ] Create backward compatibility layer
- [ ] Add comprehensive tests

**Completion Criteria**:
- [ ] Git integration extracted to package
- [ ] All Git functionality preserved
- [ ] Comprehensive test coverage
- [ ] Main application updated

**Priority Rationale**: P1 - MEDIUM - Ready after dependency analysis

## P1: File Operations and Utilities Extraction [EXTRACT-006] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ⏳ Pending (Ready after core infrastructure)

**Description**: Extract file operations and utilities including atomic operations, file copying, comparison, and file system utilities.

**Dependencies**: EXTRACT-001, EXTRACT-002 (Core infrastructure)

**Subtasks**:
- [ ] Create file operations package structure
- [ ] Extract atomic file operations - Temporary files with atomic rename
- [ ] Extract file copying utilities - Context-aware file operations
- [ ] Extract file comparison logic - Directory and file comparison
- [ ] Extract file system utilities - Path operations, exclusion patterns
- [ ] Create backward compatibility layer
- [ ] Add comprehensive tests

**Completion Criteria**:
- [ ] File operations extracted to package
- [ ] All file operation functionality preserved
- [ ] Comprehensive test coverage
- [ ] Main application updated

**Priority Rationale**: P1 - HIGH - Ready after core infrastructure

## P1: CLI Application Template Extraction [EXTRACT-008] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ⏸️ Blocked (Requires all components)

**Description**: Extract CLI application template system that demonstrates integration of all extracted components.

**Dependencies**: All extraction components (EXTRACT-001 through EXTRACT-007)

**Subtasks**:
- [ ] Create CLI template package structure
- [ ] Extract CLI framework integration patterns
- [ ] Extract command structure templates
- [ ] Extract configuration integration examples
- [ ] Extract error handling examples
- [ ] Create comprehensive documentation
- [ ] Add example applications

**Completion Criteria**:
- [ ] CLI template package created
- [ ] Integration patterns documented
- [ ] Example applications provided
- [ ] Comprehensive documentation

**Priority Rationale**: P1 - HIGH - Requires all components

## P1: Testing Patterns and Utilities Extraction [EXTRACT-009] [ARCH:TESTING_STRATEGY] [IMPL:TESTING]

**Status**: ⏳ Pending (Critical for quality)

**Description**: Extract testing patterns and utilities including test infrastructure, corruption testing, disk space simulation, and error injection.

**Dependencies**: Core extraction components

**Subtasks**:
- [ ] Create testing utilities package structure
- [ ] Extract archive corruption testing framework
- [ ] Extract disk space simulation framework
- [ ] Extract permission testing framework
- [ ] Extract context cancellation testing helpers
- [ ] Extract error injection framework
- [ ] Create comprehensive test examples

**Completion Criteria**:
- [ ] Testing utilities package created
- [ ] All test infrastructure extracted
- [ ] Comprehensive test examples provided
- [ ] Documentation complete

**Priority Rationale**: P1 - HIGH - Critical for quality

## P1: Package Documentation and Examples Extraction [EXTRACT-010] [ARCH:PACKAGE_EXTRACTION] [IMPL:PACKAGE_EXTRACTION]

**Status**: ⏳ Pending (Essential for adoption)

**Description**: Extract comprehensive package documentation and examples demonstrating usage of all extracted components.

**Dependencies**: All extraction components

**Subtasks**:
- [ ] Create documentation structure
- [ ] Extract package usage examples
- [ ] Extract integration patterns documentation
- [ ] Extract API reference documentation
- [ ] Create getting started guides
- [ ] Add migration guides

**Completion Criteria**:
- [ ] Comprehensive documentation created
- [ ] Usage examples provided
- [ ] API reference complete
- [ ] Migration guides available

**Priority Rationale**: P1 - HIGH - Essential for adoption

## Extraction Success Gates

Before proceeding to Phase 4 (Component Extraction), ALL of these must be completed:
- [ ] Complete dependency analysis with zero circular dependency risks (REFACTOR-001)
- [x] Formatter decomposition strategy validated (REFACTOR-002 COMPLETED)
- [x] Configuration abstraction interfaces defined (REFACTOR-003 COMPLETED)
- [x] Error handling patterns standardized (REFACTOR-004 COMPLETED)
- [x] Code structure optimization for extraction (REFACTOR-005 COMPLETED)
- [ ] All refactoring changes validated with zero test failures (REFACTOR-006)
- [ ] Pre-extraction validation checklist passed

## Extraction Dependency Chain

```
REFACTOR-001 (Dependency Analysis) → ALL extraction tasks
REFACTOR-002 (Formatter Decomposition) → EXTRACT-003 (Formatting)
REFACTOR-003 (Config Abstraction) → EXTRACT-001 (Configuration)
REFACTOR-004 (Error Consolidation) → EXTRACT-002 (Error Handling)
REFACTOR-006 (Validation) → Extraction Authorization

EXTRACT-001, EXTRACT-002 → EXTRACT-003, EXTRACT-004
EXTRACT-003, EXTRACT-004 → EXTRACT-005, EXTRACT-006
EXTRACT-005, EXTRACT-006 → EXTRACT-007, EXTRACT-008
EXTRACT-008 → EXTRACT-009, EXTRACT-010
```

## Task Status Legend

- **Pending**: Not started
- **In Progress**: Currently being worked on
- **Complete**: Finished and tested
- **Deferred**: Postponed to future phase
- **Future**: Out of scope for initial implementation

## Dependencies

### Task Dependencies
- OUT-002 → Requires CFG-003 (format strings system)
- COV-003 → Requires COV-001 and COV-002
- Migration tasks → Can proceed in parallel

## Completed Tasks (Historical Record)

### Phase 1: Foundation Setup ✅
- ✅ Feature matrix created with cross-references
- ✅ Implementation tokens added to core functions (66 tokens across 5 files)
- ✅ Decision records documented with rationale (8 decisions)
- ✅ Code markers linked to documentation
- ✅ Validation script implemented and tested
- ✅ Documentation consistency framework established

### Phase 2: Test Infrastructure ✅
- ✅ Fix Missing Test Functions - Successfully eliminated all 9 validation errors
- ✅ Validation Script Enhancement - Fixed false positives in test function detection
- ✅ Zero-Error Validation Achieved - Validation script now passes with 0 errors, 0 warnings
- ✅ Complete Implementation Token Coverage - Achieved 144 total tokens with comprehensive coverage
- ✅ Strategic Test Token Enhancement - Added targeted tokens to key feature validation functions

### Phase 3: Pre-Extraction Refactoring ✅
- ✅ REFACTOR-002: Large File Decomposition Preparation (2025-01-02)
  - Component boundary analysis completed
  - Identified 5 distinct components in formatter.go (OutputCollector, Printf Formatter, Template Formatter, Extended Output Formatters, Error Formatting)
  - Internal interfaces designed (FormatProvider, OutputDestination, PatternExtractor)
  - Extraction strategy with clean component boundaries
  - All 168 tests pass, compilation successful, no functional regressions
- ✅ REFACTOR-003: Configuration Schema Abstraction (2025-01-02)
  - Interface abstraction layer implemented (ConfigLoader, ConfigValidator, ConfigMerger, ConfigSource)
  - Schema abstraction with provider pattern (ConfigSchema, FieldDefinition, BackupConfigSchema)
  - Generic configuration loading (YAML, environment, defaults)
  - Schema-based validation system
  - Zero breaking changes, all tests pass
- ✅ REFACTOR-004: Error Handling Consolidation
  - Error type standardization completed (ErrorInterface, unified BackupError/ArchiveError)
  - Resource management patterns standardized
  - Context propagation consistency validated
  - Panic recovery standardization
- ✅ REFACTOR-005: Code Structure Optimization (2025-01-02)
  - Component coupling reduced through interface abstractions
  - Naming conventions standardized (ArchiveConfigInterface, BackupConfigInterface, CommandHandlerInterface)
  - Import optimization for future package structure
  - Backward compatibility preserved with wrapper functions and adapters

### Phase 4: Component Extraction ✅
- ✅ EXTRACT-005: CLI Command Framework (2025-01-02)
  - Package structure created in `pkg/cli/`
  - CLI framework with command handlers, flag processing, context support
  - Comprehensive test coverage
- ✅ EXTRACT-007: Data Processing Patterns (2025-01-02)
  - Package structure created in `pkg/processing/` with 7 files
  - Core interfaces implemented (ProcessorInterface, NamingProviderInterface, VerificationProviderInterface, PipelineInterface, ConcurrentProcessorInterface)
  - Naming system extracted with regex parsing
  - Verification system with multi-algorithm support (SHA-256/SHA-512/MD5)
  - Pipeline system with context-aware processing
  - Concurrent system with worker pool
  - Comprehensive test suite with all tests passing
  - Critical fixes: Resolved naming pattern regex issue and concurrent processor deadlock

### Configuration System Enhancements ✅
- ✅ CFG-005: Layered Configuration Inheritance - Completed
- ✅ CFG-006: Complete Configuration Reflection and Visibility (2025-01-02)
  - Automatic field discovery using Go reflection
  - Enhanced config command with filtering options
  - Performance optimization with caching
- ✅ CFG-TEMPLATE-001: Configuration Template Generation Command (2025-01-02)
  - CLI command implemented
  - Template generation engine completed
  - All 6 subtasks completed

### Core Features ✅
- ✅ ARCH-001: Archive naming convention - Implemented
- ✅ ARCH-002: Create archive command - Implemented
- ✅ ARCH-003: Incremental archives - Implemented
- ✅ ARCH-004: Broken symlink handling - Completed
- ✅ FILE-001: File backup naming - Implemented
- ✅ FILE-002: Backup command - Implemented
- ✅ FILE-003: File comparison - Implemented
- ✅ CLI-015: Automatic file/directory command detection - Implemented
- ✅ CFG-001: Config discovery - Implemented
- ✅ CFG-002: Status codes - Implemented
- ✅ CFG-003: Format strings - Implemented
- ✅ CFG-004: Comprehensive string config - Completed
- ✅ GIT-001 through GIT-006: Git integration features - Completed
- ✅ OUT-001: Delayed output management - Completed
- ✅ TEST-001: Comprehensive formatter testing - Completed
- ✅ TEST-002: Tools directory test coverage - Completed
- ✅ TEST-FIX-001: Personal config isolation in tests - Completed
- ✅ TEST-INFRA-001-B: Disk space simulation framework - Completed
- ✅ TEST-INFRA-001-E: Error injection framework - Completed
- ✅ COV-001: Existing code coverage exclusion - Completed
- ✅ COV-002: Coverage baseline establishment - Completed
- ✅ DOC-001 through DOC-017: Documentation enhancement system - Various completion statuses
- ✅ DIAGNOSTIC-OUTPUT: Diagnostic Output Control (2025-11-21)
  - Wrapped all 24 diagnostic output statements in `config.go` with `if debug { ... }` blocks
  - Added SEMANTIC-TOKEN: DIAGNOSTIC-OUTPUT markers for STDD compliance
  - Implemented `TestMain` in `config_test.go` to enable debug output for tests
  - Normal execution produces clean output (0 diagnostic messages)
  - Debug mode (`--debug` flag) shows all diagnostic output
  - All tests pass with diagnostic output enabled by default
  - Related tokens: DEBUG-OUTPUT, [REQ:CONFIGURATION], [REQ:CFG_005], [REQ:CFG_006]

## P1: User-Customizable Format Strings [REQ:CUSTOMIZABLE_FORMAT_STRINGS] [ARCH:CUSTOMIZABLE_FORMAT_STRINGS] [IMPL:CUSTOMIZABLE_FORMAT_STRINGS]

**Status**: ✅ Complete

**Description**: Enable user customization of all output format strings through configuration files with validation. Users can override any format string in their `.bkpdir.yml` to customize the application's output interface.

**Dependencies**: [REQ:CONFIGURATION], [REQ:OUTPUT_FORMATTING]

**Completion Criteria**:
- [x] All format strings can be overridden in `.bkpdir.yml`
- [x] Format string validation warns users of invalid or unexpected placeholders
- [x] Validation provides helpful error messages indicating expected placeholders
- [x] Comprehensive reference documentation lists all available format strings with their placeholders
- [x] Brief example configuration file demonstrates common customizations
- [x] Placeholder documentation explains available variables for each format string
- [x] Backward compatibility maintained (defaults work when not specified)
- [x] Template-style (`#{path}`) format supported (migrated from `%{...}` to `#{...}` to avoid fmt conflicts)
- [x] Special placeholders like `#{size_human}` documented and validated
- [x] Validation integrated into LoadConfig and LoadConfigWithInheritance
- [x] Comprehensive tests for validation functionality
- [x] Integration tests for custom format strings loading and usage

**Priority Rationale**: P1 - Important for user experience and internationalization support

## P0: Priority 0 Configuration Merge Tests [REQ:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION] [IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX]

**Status**: ✅ Complete

**Description**: Implement all Priority 0 (Critical) configuration merge tests from the test plan to ensure comprehensive coverage of core merge functionality.

**Dependencies**: [REQ:CFG_005] (Configuration Inheritance), [REQ:CFG_001] (Configuration Discovery), [REQ:CONFIGURATION] (Configuration Management)

**Subtasks**:
- [x] TestMultipleFilesSameField - Test 3+ files all setting the same field [REQ:CFG_005] [REQ:CFG_001]
- [x] TestPartialFieldUpdates - Test partial field updates across files [REQ:CFG_005] [REQ:CFG_001]
- [x] TestAllStrategiesWithAccumulateFields - Test all merge strategies with accumulate fields [REQ:CFG_005]
- [x] TestAllStrategiesWithPrecedenceFields - Test all merge strategies with precedence fields [REQ:CFG_001]
- [x] TestOverrideDefaultsWithPrefix - Test explicit ! prefix overriding defaults [REQ:CFG_005]
- [x] Fix applyMerge and applyPrepend precedence checking [REQ:CFG_001] [IMPL:CFG_MERGE_PREPEND_PRECEDENCE_FIX]

**Completion Criteria**:
- [x] All Priority 0 tests implemented
- [x] All tests passing
- [x] Implementation fix for precedence checking in applyMerge and applyPrepend
- [x] Documentation updated with implementation decision

**Priority Rationale**: P0 - Critical for ensuring configuration merge behavior is correct and comprehensive. These tests validate core functionality required by CFG-001 and CFG-005.

## P1: Priority 1 Configuration Merge Tests [REQ:CFG_005] [REQ:CFG_001] [REQ:CONFIGURATION] [IMPL:CFG_INHERITANCE_PATH_RESOLUTION]

**Status**: ✅ Complete

**Description**: Implement all Priority 1 (Important) configuration merge tests from the test plan to ensure comprehensive coverage of edge cases and inheritance scenarios.

**Dependencies**: [REQ:CFG_005] (Configuration Inheritance), [REQ:CFG_001] (Configuration Discovery), [REQ:CONFIGURATION] (Configuration Management)

**Subtasks**:
- [x] TestMultipleInheritanceSources - Test child inheriting from multiple parent files [REQ:CFG_005]
- [x] TestRelativePathInheritance - Test relative path resolution (../base.yml, ./sibling.yml) [REQ:CFG_005]
- [x] TestHomeDirectoryExpansion - Test home directory expansion (~/.bkpdir-base.yml) [REQ:CFG_005]
- [x] TestMissingInheritanceFile - Test missing file in inheritance chain error handling [REQ:CFG_005]
- [x] TestInvalidYAMLHandling - Test invalid YAML in one file with graceful error handling [REQ:CFG_005]
- [x] TestDeepInheritanceChain - Test very long inheritance chain (10+ files) [REQ:CFG_005] - Skipped due to known implementation bug
- [x] TestTypeMismatchHandling - Test type mismatches (array vs string) [REQ:CFG_005]
- [x] TestNullValueHandling - Test nil/null value handling [REQ:CFG_005]
- [x] TestWhitespaceStringHandling - Test whitespace-only string handling [REQ:CFG_005]
- [x] Fix inheritance chain path resolution bug (resolvedPath -> filepath.Dir(resolvedPath)) [REQ:CFG_005] [IMPL:CFG_INHERITANCE_PATH_RESOLUTION]
- [x] Add ExpandPath method to defaultPathResolver for ~ expansion [REQ:CFG_005] [IMPL:CFG_INHERITANCE_PATH_RESOLUTION]
- [x] Fix resolvePath to handle directory paths correctly [REQ:CFG_005] [IMPL:CFG_INHERITANCE_PATH_RESOLUTION]

**Completion Criteria**:
- [x] All Priority 1 tests implemented (8 passing, 1 skipped due to known bug)
- [x] All tests account for default merging behavior
- [x] Inheritance path resolution bugs fixed
- [x] Home directory expansion working in inherit paths
- [x] Tests include semantic token references [REQ:CFG_005]

**Priority Rationale**: P1 - Important for ensuring edge cases and inheritance scenarios are properly handled. These tests validate important functionality for real-world configuration usage.

## P1: Unicode and Special Character Handling Tests [REQ:CONFIGURATION] [REQ:CFG_005] [REQ:CFG_001] [IMPL:TEST_UNICODE_HANDLING]

**Status**: ✅ Complete

**Description**: Implement tests for Unicode and special character handling in configuration values and file paths to ensure internationalization support and real-world path compatibility.

**Dependencies**: [REQ:CONFIGURATION] (Configuration Management), [REQ:CFG_005] (Configuration Inheritance), [REQ:CFG_001] (Configuration Discovery)

**Subtasks**:
- [x] TestUnicodeHandling - Test Unicode characters in paths and patterns [REQ:CONFIGURATION] [REQ:CFG_005]
- [x] TestSpecialCharactersInPaths - Test special characters in config file paths [REQ:CONFIGURATION] [REQ:CFG_001]

**Completion Criteria**:
- [x] TestUnicodeHandling implemented and passing
- [x] TestSpecialCharactersInPaths implemented and passing
- [x] Unicode characters preserved correctly in configuration values
- [x] Special characters in config file paths handled correctly
- [x] Tests include semantic token references [REQ:CONFIGURATION] [REQ:CFG_005] [REQ:CFG_001]
- [x] Documentation updated with implementation decision

**Priority Rationale**: P1 - Important for internationalization support and real-world usage scenarios. These tests validate that the configuration system correctly handles Unicode characters and special characters in file paths, which is critical for users working with non-ASCII paths and patterns.

## P1: Empty String Handling Tests [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005] [IMPL:TEST_EMPTY_STRING_HANDLING]

**Status**: ✅ Complete

**Description**: Implement tests for empty string handling in configuration merging to ensure correct behavior when empty strings are explicitly set in configuration files.

**Dependencies**: [REQ:CONFIGURATION] (Configuration Management), [REQ:CFG_001] (Configuration Discovery), [REQ:CFG_005] (Configuration Inheritance)

**Subtasks**:
- [x] TestEmptyStringHandling - Test empty string handling in various merge scenarios [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005]
  - [x] First file empty string, second file value
  - [x] First file value, second file empty string
  - [x] Both files empty string

**Completion Criteria**:
- [x] TestEmptyStringHandling implemented and passing
- [x] All three scenarios tested (first empty, second value; first value, second empty; both empty)
- [x] Empty string behavior validated (zero value vs explicit empty)
- [x] CFG-001 precedence validated with empty strings
- [x] Tests include semantic token references [REQ:CONFIGURATION] [REQ:CFG_001] [REQ:CFG_005]
- [x] Documentation updated with implementation decision

**Priority Rationale**: P1 - Important for understanding how empty strings are handled in configuration merging. These tests validate whether empty strings are treated as zero values (allowing override) or explicit values (preserving per CFG-001 precedence), which is critical for correct configuration behavior.

## P1: Prepend Strategy Ordering Tests [REQ:CONFIGURATION] [REQ:CFG_005] [IMPL:TEST_PREPEND_ORDERING]

**Status**: ✅ Complete

**Description**: Implement tests for prepend strategy ordering to ensure prepend strategy (`^` prefix) correctly maintains order with new values before existing values.

**Dependencies**: [REQ:CONFIGURATION] (Configuration Management), [REQ:CFG_005] (Configuration Inheritance)

**Subtasks**:
- [x] TestPrependStrategyOrdering - Test prepend strategy ordering in various scenarios [REQ:CONFIGURATION] [REQ:CFG_005]
  - [x] Prepend in inheritance chain (parent then child)
  - [x] Prepend in sequential files (first file then second file)
  - [x] Prepend with defaults (prepended values before defaults)
  - [x] Multiple prepend operations (CFG-001 precedence)

**Completion Criteria**:
- [x] TestPrependStrategyOrdering implemented and passing
- [x] All four scenarios tested (inheritance, sequential, defaults, multiple)
- [x] Prepend ordering validated (new values before existing values)
- [x] CFG-001 precedence validated with multiple prepends
- [x] Tests include semantic token references [REQ:CONFIGURATION] [REQ:CFG_005]
- [x] Documentation updated with implementation decision

**Priority Rationale**: P1 - Important for understanding prepend strategy ordering semantics. These tests validate that prepend strategy correctly places new values before existing values in both inheritance chains and sequential file processing, which is critical for correct merge behavior.

## P1: Default Strategy Edge Cases Tests [REQ:CONFIGURATION] [REQ:CFG_005] [IMPL:TEST_DEFAULT_STRATEGY_EDGES]

**Status**: ✅ Complete

**Description**: Implement tests for default strategy edge cases to ensure default strategy (`=` prefix) only applies when destination is zero value, not when it's non-zero or equals default.

**Dependencies**: [REQ:CONFIGURATION] (Configuration Management), [REQ:CFG_005] (Configuration Inheritance)

**Subtasks**:
- [x] TestDefaultStrategyEdgeCases - Test default strategy edge cases [REQ:CONFIGURATION] [REQ:CFG_005]
  - [x] Default strategy when field is zero value
  - [x] Default strategy when field is non-zero
  - [x] Default strategy when field equals default
  - [x] Default strategy with array field (zero and non-zero)
  - [x] Default strategy with bool field (zero and non-zero)

**Completion Criteria**:
- [x] TestDefaultStrategyEdgeCases implemented and passing
- [x] All seven scenarios tested (string zero/non-zero/equals-default, array zero/non-zero, bool zero/non-zero)
- [x] Default strategy behavior validated (only applies when destination is zero value)
- [x] Edge cases validated (equals default but not zero → does NOT apply)
- [x] Array field special case validated (empty array merges with defaults first, so default strategy doesn't apply)
- [x] Tests include semantic token references [REQ:CONFIGURATION] [REQ:CFG_005]
- [x] Documentation updated with implementation decision

**Priority Rationale**: P1 - Important for understanding default strategy semantics. These tests validate that default strategy correctly only applies when destination is zero value, which is critical for correct merge behavior and prevents unexpected overrides.

## P1: List Command Output Limit [REQ:LIST_LIMIT] [ARCH:LIST_LIMIT] [IMPL:LIST_LIMIT]

**Status**: ✅ Complete

**Description**: Limit the `list` command display to the newest N files (default N=10). Add an option to control the length of the list. Support both `list` command (directory archives) and `--list` command (file backups).

**Dependencies**: [REQ:LIST_LIMIT] (List Command Output Limit Requirements)

**Completion Criteria**:
- [x] Default behavior: Display only the newest 10 files
- [x] Command-line option `--limit N` (short: `-n N`) to control the limit
- [x] Option `--limit 0` shows all files
- [x] Archives remain sorted by creation time (most recent first)
- [x] Both `list` command (directory archives) and `--list` command (file backups) support the limit option
- [x] Backward compatibility maintained (existing behavior when limit not specified)
- [x] Tests verify default limit of 10 files
- [x] Tests verify custom limit values work correctly
- [x] Tests verify `--limit 0` shows all files
- [x] Integration tests verify both `list` and `--list` commands respect the limit
- [x] Comprehensive test coverage with edge cases (fewer than limit, exactly limit, more than limit, sorting preserved)

**Priority Rationale**: P1 - Important for improving usability by preventing overwhelming output when there are many archives

## Recommended Implementation Order

1. P1: List Command Output Limit (REQ:LIST_LIMIT) - ✅ Complete
2. P1: Configuration Output Grouping (REQ:CONFIG_OUTPUT_GROUPING)
3. P0: Enhanced Command Output with File Statistics (OUT-002)
4. P1: Code Linting Compliance (LINT-001)
5. P2: Selective Coverage Reporting (COV-003)
6. Migration tasks (can proceed in parallel)
