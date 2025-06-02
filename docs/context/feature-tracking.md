# Feature Tracking Matrix

## ⚡ AI ASSISTANT INSTRUCTIONS - MANDATORY TOKEN REFERENCING

### 🚨 CRITICAL AI DIRECTIVE: ALWAYS REFERENCE CONTEXT FILE TOKENS WITH CODE CHANGES

**IMMEDIATE ACTION REQUIRED FOR ALL AI ASSISTANTS:**

When making ANY code changes, you MUST:

1. **🔍 SEARCH FOR EXISTING TOKENS**: Before making code changes, search the context directory files for relevant implementation tokens:
   - Use grep or search for patterns like `// ARCH-001:`, `// FILE-002:`, `// CFG-003:`, etc.
   - Look for feature IDs in the feature tracking tables below
   - Check all context files: `specification.md`, `requirements.md`, `architecture.md`, `testing.md`

2. **📋 REFERENCE RELEVANT TOKENS**: For every code change, explicitly mention related tokens:
   ```
   Example: "This change affects ARCH-001 (Archive naming), FILE-002 (File backup), and CFG-003 (Format strings)"
   ```

3. **🔗 LINK TO DOCUMENTATION**: Reference the specific sections in context files:
   ```
   Example: "See ARCH-001 in feature-tracking.md line 119, specification.md section 'Archive Naming', architecture.md 'Archive Service'"
   ```

4. **📝 UPDATE TOKEN REGISTRY**: Add new implementation tokens for new features:
   - Format: `// FEATURE-ID: Brief description`
   - Register in the feature tracking tables below
   - Update all affected context files

5. **✅ VALIDATE CROSS-REFERENCES**: Ensure all token references are consistent across:
   - feature-tracking.md (this file)
   - specification.md
   - requirements.md
   - architecture.md
   - testing.md

### 🎯 TOKEN SEARCH PATTERNS FOR AI ASSISTANTS

**Search these patterns in context files:**
- `ARCH-[0-9]+` (Archive operations)
- `FILE-[0-9]+` (File operations)
- `CFG-[0-9]+` (Configuration)
- `GIT-[0-9]+` (Git integration)
- `OUT-[0-9]+` (Output management)
- `TEST-[0-9]+` (Testing infrastructure)
- `DOC-[0-9]+` (Documentation system)
- `LINT-[0-9]+` (Code quality)
- `COV-[0-9]+` (Coverage)

**Token Reference Template for AI:**
```
This code change relates to:
- Primary token: [TOKEN-ID] ([Description])
- Secondary tokens: [TOKEN-ID], [TOKEN-ID] ([Brief descriptions])
- Documentation references:
  - feature-tracking.md: [Line numbers or sections]
  - specification.md: [Relevant sections]
  - requirements.md: [Requirement sections]
  - architecture.md: [Architecture components]
  - testing.md: [Test coverage areas]
```

### 🚨 AI ENFORCEMENT RULES

**REJECT code changes that do not:**
- Reference at least one existing feature token when modifying related functionality
- Propose new tokens for genuinely new features
- Link to relevant documentation sections
- Update the feature tracking tables when adding new tokens

**ACCEPT code changes that:**
- Clearly state which tokens are affected
- Reference specific documentation sections
- Propose updates to context files when needed
- Follow the token naming conventions

## Purpose
This document serves as a master index linking features across all documentation layers to ensure no unplanned changes occur during development.

## ⚠️ MANDATORY ENFORCEMENT: Context File Update Requirements

### 🚨 CRITICAL RULE: NO CODE CHANGES WITHOUT CONTEXT UPDATES
**ALL code modifications MUST include corresponding updates to relevant context documentation files. Failure to update context files invalidates the change.**

### 📋 MANDATORY CONTEXT FILE CHECKLIST
Before implementing ANY code change, developers MUST complete this checklist:

#### **Phase 1: Pre-Change Analysis (REQUIRED)**
- [ ] **Feature Impact Analysis**: Identify which existing features are affected by the change
- [ ] **Context File Mapping**: Determine which context files require updates based on change type:

| Change Type | Required Context File Updates |
|-------------|------------------------------|
| **New Feature** | feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md |
| **Feature Modification** | feature-tracking.md, + all files containing the feature |
| **Bug Fix** | feature-tracking.md (if affects documented behavior) |
| **Configuration Change** | feature-tracking.md, specification.md, requirements.md |
| **API/Interface Change** | feature-tracking.md, specification.md, architecture.md |
| **Test Addition** | feature-tracking.md, testing.md |
| **Error Handling Change** | feature-tracking.md, specification.md, architecture.md |
| **Performance Change** | feature-tracking.md, architecture.md |

#### **Phase 2: Documentation Updates (REQUIRED)**
- [ ] **Update feature-tracking.md**: Add/modify feature entries, update status, add implementation tokens
- [ ] **Update specification.md**: Modify user-facing behavior descriptions if applicable
- [ ] **Update requirements.md**: Add/modify implementation requirements with traceability
- [ ] **Update architecture.md**: Update technical implementation descriptions
- [ ] **Update testing.md**: Add/modify test coverage requirements and descriptions
- [ ] **Check immutable.md**: Verify no immutable requirements are violated
- [ ] **Update cross-references**: Ensure all documents link correctly to each other

#### **Phase 3: Implementation Tokens (REQUIRED)**
- [ ] **Add Implementation Tokens**: Mark ALL modified code with feature ID comments
- [ ] **Update Token Registry**: Record new tokens in feature-tracking.md
- [ ] **Verify Token Consistency**: Ensure tokens match feature IDs and descriptions

#### **Phase 4: Validation (REQUIRED)**
- [ ] **Run Documentation Validation**: Execute `make validate-docs` (if available)
- [ ] **Cross-Reference Check**: Verify all links between documents are valid
- [ ] **Feature Matrix Update**: Ensure feature tracking matrix reflects all changes
- [ ] **Test Coverage Verification**: Confirm test references are accurate

### 🔒 ENFORCEMENT MECHANISMS

#### **Automated Validation Rules**
1. **Feature ID Consistency**: All implementation tokens must correspond to feature matrix entries
2. **Cross-Reference Integrity**: All document links must resolve correctly
3. **Status Synchronization**: Feature status must be consistent across all documents
4. **Test Traceability**: All features must have corresponding test references

#### **Manual Review Requirements**
1. **Context Documentation Review**: Every code change must include review of context file updates
2. **Feature Impact Assessment**: Verify that all affected features are properly documented
3. **Immutable Requirement Check**: Confirm no immutable requirements are violated
4. **Backward Compatibility Verification**: Ensure changes preserve documented compatibility

#### **🤖 AI ASSISTANT ENFORCEMENT RULES**

**MANDATORY AI VALIDATION CHECKLIST:**

Before any AI assistant provides code changes, it MUST validate:

1. **🔍 TOKEN SEARCH PERFORMED**
   - [ ] Searched all context files for existing feature tokens
   - [ ] Identified all relevant feature IDs using grep patterns: `ARCH-[0-9]+`, `FILE-[0-9]+`, `CFG-[003:`, etc.
   - [ ] Listed affected tokens in response

2. **📋 TOKEN IMPACT DOCUMENTED**
   - [ ] Clearly stated which existing tokens are affected
   - [ ] Proposed new tokens for new functionality
   - [ ] Referenced specific line numbers in feature-tracking.md
   - [ ] Linked to relevant sections in other context files

3. **📄 DOCUMENTATION CROSS-REFERENCES**
   - [ ] Referenced specification.md sections
   - [ ] Referenced requirements.md areas
   - [ ] Referenced architecture.md components
   - [ ] Referenced testing.md coverage areas

4. **🔄 UPDATE REQUIREMENTS IDENTIFIED**
   - [ ] Listed which context files need updates
   - [ ] Specified what changes are needed in each file
   - [ ] Provided update templates or examples

**AI RESPONSE TEMPLATE (REQUIRED FORMAT):**

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Description of impact]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [X], [Y], [Z]
- **specification.md**: Section "[Section Name]"
- **requirements.md**: "[Requirement area]"
- **architecture.md**: "[Component name]"
- **testing.md**: "[Test coverage area]"

### 🔄 Required Context File Updates:
- [ ] feature-tracking.md: [Specific changes]
- [ ] specification.md: [Specific changes]
- [ ] requirements.md: [Specific changes]
- [ ] architecture.md: [Specific changes]
- [ ] testing.md: [Specific changes]

### 💻 Implementation Details:
- **New tokens**: [TOKEN-ID: Description]
- **Modified tokens**: [TOKEN-ID: Updated description]
- **Code files affected**: [List of files]
```

**🚨 AI REJECTION CRITERIA (ZERO TOLERANCE)**

AI assistants MUST REJECT their own responses if they:
- Fail to search for existing tokens before proposing changes
- Make code changes without referencing related feature tokens
- Don't identify which context files need updates
- Propose changes that conflict with immutable requirements
- Don't follow the required response template format

**✅ AI APPROVAL CRITERIA**

AI assistants should ONLY PROVIDE responses that:
- Include complete token impact analysis
- Reference specific documentation sections with line numbers
- Propose concrete updates to all affected context files
- Follow established naming conventions for new tokens
- Include implementation token placement in code changes

### 📁 CONTEXT FILE RESPONSIBILITIES

#### **feature-tracking.md** (THIS FILE)
- **ALWAYS UPDATE**: Must be updated for every code change
- **Content**: Feature registry, implementation tokens, status tracking, decision records
- **Triggers**: Any code modification, new features, status changes, architectural decisions

#### **specification.md** 
- **UPDATE WHEN**: User-facing behavior changes, new commands, configuration options, output format changes
- **Content**: User interface specifications, command behaviors, configuration schemas
- **Triggers**: New CLI commands, configuration changes, output modifications, user-visible behavior changes

#### **requirements.md**
- **UPDATE WHEN**: Implementation requirements change, new functional requirements, non-functional requirement modifications
- **Content**: Implementation requirements, traceability matrices, acceptance criteria
- **Triggers**: New features, requirement modifications, acceptance criteria changes

#### **architecture.md**
- **UPDATE WHEN**: Technical implementation changes, new components, interface modifications, design decisions
- **Content**: System architecture, component interactions, design patterns, technical decisions
- **Triggers**: Code structure changes, new components, interface modifications, architectural decisions

#### **testing.md**
- **UPDATE WHEN**: New tests added, test strategies change, coverage requirements modified
- **Content**: Test strategies, coverage requirements, test descriptions, validation approaches
- **Triggers**: New test files, test strategy changes, coverage modifications

#### **immutable.md**
- **UPDATE WHEN**: Never (immutable requirements), but must CHECK for violations
- **Content**: Unchangeable requirements, backward compatibility constraints
- **Triggers**: Check only - never modify

### 🚫 CHANGE REJECTION CRITERIA

Changes will be rejected if:
- [ ] Context files are not updated to reflect code changes
- [ ] Feature tracking matrix is not updated with new features or status changes
- [ ] Implementation tokens are missing from modified code
- [ ] Cross-references between documents are broken
- [ ] Immutable requirements are violated
- [ ] Test coverage documentation doesn't match actual tests

## Feature Tracking Format

Each feature must be documented across all relevant layers:
- **Immutable ID**: Unique identifier (e.g., `ARCH-001`, `CMD-001`, `CFG-001`)
- **Specification Reference**: Link to spec section
- **Requirements Reference**: Link to requirements section  
- **Architecture Reference**: Link to architecture section
- **Test Reference**: Link to test coverage
- **Implementation Tokens**: Code markers for traceability

## Current Feature Registry

### Core Archive Operations
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| ARCH-001 | Archive naming convention | Archive naming | ArchiveCreator | TestGenerateArchiveName | Implemented | `// ARCH-001: Archive naming` |
| ARCH-002 | Create archive command | Create archive ops | Archive Service | TestCreateFullArchive | Implemented | `// ARCH-002: Archive creation` |
| ARCH-003 | Incremental archives | Incremental logic | CompressionEngine | TestCreateIncremental | Implemented | `// ARCH-003: Incremental` |

### File Backup Operations  
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| FILE-001 | File backup naming | File backup naming | BackupCreator | TestGenerateBackupName | Implemented | `// FILE-001: Backup naming` |
| FILE-002 | Backup command | File backup ops | File Backup Service | TestCreateFileBackup | Implemented | `// FILE-002: File backup` |
| FILE-003 | File comparison | Identical detection | FileComparator | TestCompareFiles | Implemented | `// FILE-003: File comparison` |

### Configuration System
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| CFG-001 | Config discovery | Config discovery | Configuration Layer | TestGetConfigSearchPath | Implemented | `// CFG-001: Config discovery` |
| CFG-002 | Status codes | Status code config | Config object | TestDefaultConfig | Implemented | `// CFG-002: Status codes` |
| CFG-003 | Format strings | Output formatting | OutputFormatter | TestTemplateFormatter | Implemented | `// CFG-003: Format strings` |
| CFG-004 | Comprehensive string config | String externalization | String Management | TestStringExternalization | Completed | `// CFG-004: String externalization` |

### Git Integration
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| GIT-001 | Git info extraction | Git requirements | Git Service | TestGitIntegration | Completed | `// GIT-001: Git extraction` |
| GIT-002 | Branch/hash naming | Git naming | NamingService | TestGitNaming | Completed | `// GIT-002: Git naming` |
| GIT-003 | Git status detection | Git requirements | Git Service | TestGitStatus | Completed | `// GIT-003: Git status` |
| GIT-004 | Git submodule support | Git requirements | Git Service | TestGitSubmodules | Not Started | `// GIT-004: Git submodules` |
| GIT-005 | Git configuration integration | Git requirements | Git Service | TestGitConfig | Not Started | `// GIT-005: Git config` |
| GIT-006 | Configurable dirty status | Git requirements | Git Service | TestGitDirtyConfig | Completed | `// GIT-006: Git dirty config` |

### Output Management
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| OUT-001 | Delayed output management | Output control requirements | Output System | TestDelayedOutput | Completed | `// OUT-001: Delayed output` |

### Testing Infrastructure
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| TEST-001 | Comprehensive formatter testing | Test coverage requirements | Test Infrastructure | TestFormatterCoverage | Completed | `// TEST-001: Formatter testing` |
| TEST-002 | Tools directory test coverage | Test coverage requirements | Test Infrastructure | TestToolsCoverage | Not Started | `// TEST-002: Tools directory testing` |
| TEST-INFRA-001-B | Disk space simulation framework | Testing infrastructure requirements | Test Infrastructure | TestDiskSpaceSimulation | Completed | `// TEST-INFRA-001-B: Disk space simulation framework` |
| TEST-INFRA-001-E | Error injection framework | Testing infrastructure requirements | Test Infrastructure | TestErrorInjection | Completed | `// TEST-INFRA-001-E: Error injection framework` |

### Code Quality
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| LINT-001 | Code linting compliance | Code quality standards | Linting rules | TestLintCompliance | In Progress | `// LINT-001: Lint compliance` |
| COV-001 | Existing code coverage exclusion | Coverage control requirements | Coverage filtering | TestCoverageExclusion | ✅ Completed | `// COV-001: Coverage exclusion` |
| COV-002 | Coverage baseline establishment | Coverage metrics | Coverage tracking | TestCoverageBaseline | ✅ Completed | `// COV-002: Coverage baseline` |
| COV-003 | Selective coverage reporting | Coverage configuration | Coverage engine | TestSelectiveCoverage | Not Started | `// COV-003: Selective coverage` |

### Documentation Enhancement System
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| DOC-001 | Semantic linking system | Cross-reference requirements | LinkingEngine | TestSemanticLinks | Completed | `// DOC-001: Semantic linking` |
| DOC-002 | Sync framework | Synchronization requirements | SyncFramework | TestDocumentSync | Far Future (Unlikely) | `// DOC-002: Document sync` |
| DOC-003 | Enhanced traceability | Traceability requirements | TraceabilitySystem | TestEnhancedTrace | Far Future (Unlikely) | `// DOC-003: Enhanced traceability` |
| DOC-004 | Automated validation | Validation requirements | ValidationEngine | TestAutomatedValidation | Far Future (Unlikely) | `// DOC-004: Automated validation` |
| DOC-005 | Change impact analysis | Impact analysis requirements | ImpactAnalyzer | TestChangeImpact | Far Future (Unlikely) | `// DOC-005: Change impact` |

### Pre-Extraction Refactoring
| Feature ID | Specification | Requirements | Architecture | Testing | Status | Implementation Tokens |
|------------|---------------|--------------|--------------|---------|--------|----------------------|
| REFACTOR-001 | Dependency analysis and interface standardization | Pre-extraction requirements | Component interfaces | TestDependencyAnalysis | Not Started | `// REFACTOR-001: Dependency analysis` |
| REFACTOR-002 | Large file decomposition preparation | Code structure requirements | Component boundaries | TestFormatterDecomposition | COMPLETED (2025-01-02) | `// REFACTOR-002: Formatter decomposition` |
| REFACTOR-003 | Configuration schema abstraction | Configuration extraction requirements | Config interfaces | TestConfigAbstraction | Not Started | `// REFACTOR-003: Config abstraction` |
| REFACTOR-004 | Error handling consolidation | Error handling standards | Error type patterns | TestErrorStandardization | Not Started | `// REFACTOR-004: Error standardization` |
| REFACTOR-005 | Code structure optimization | Extraction preparation requirements | Structure optimization | TestStructureOptimization | Not Started | `// REFACTOR-005: Structure optimization` |
| REFACTOR-006 | Refactoring impact validation | Quality assurance requirements | Validation framework | TestRefactoringValidation | Not Started | `// REFACTOR-006: Validation` |

## Feature Change Protocol

### 🔥 ENHANCED ACTIONABLE CHANGE PROCESS

#### **📝 Adding New Features (STEP-BY-STEP)**

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Assign Feature ID**: Use prefix (ARCH, FILE, CFG, GIT, etc.) + sequential number
   - Check existing feature IDs in this document to avoid conflicts
   - Use consistent prefixes: ARCH (archive ops), FILE (file ops), CFG (config), GIT (git integration), etc.

2. **🚨 MANDATORY: Update feature-tracking.md (THIS FILE)**
   - Add new row to appropriate feature category table
   - Include: Feature ID, Specification ref, Requirements ref, Architecture ref, Testing ref, Status, Implementation tokens
   - Add feature to Implementation Status Summary if applicable

3. **🚨 MANDATORY: Update specification.md**
   - Add user-facing behavior description
   - Include command syntax, configuration options, output formats
   - Add examples of user interaction
   - Update command reference section

4. **🚨 MANDATORY: Update requirements.md**
   - Define implementation requirements with traceability
   - Add functional and non-functional requirements
   - Include acceptance criteria
   - Link to feature ID from this document

5. **🚨 MANDATORY: Update architecture.md**
   - Specify technical implementation approach
   - Add component diagrams if applicable
   - Document design patterns and interfaces
   - Include integration points

6. **🚨 MANDATORY: Update testing.md**
   - Define comprehensive test coverage before implementation
   - Specify test strategies and approaches
   - Include performance testing requirements
   - Add integration test scenarios

7. **⚡ BEFORE CODING: Add Implementation Tokens**
   - Mark ALL new code with feature ID comments: `// FEATURE-ID: Description`
   - Update token registry in this document
   - Ensure tokens are added to every function, method, and significant code block

8. **📋 FINAL STEP: Update Feature Matrix**
   - Verify all cross-references are correct
   - Check that status reflects actual implementation state
   - Ensure all documentation links are valid

#### **🔧 Modifying Existing Features (STEP-BY-STEP)**

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **🚨 CHECK IMMUTABLE STATUS FIRST**
   - Open `immutable.md` and verify feature is not protected
   - If feature is immutable, STOP - change is not allowed
   - Document any immutable constraint violations

2. **📊 Impact Analysis**
   - Search ALL context files for references to the feature
   - Use `grep -r "FEATURE-ID" docs/context/` to find all mentions
   - List affected documents and sections

3. **🚨 MANDATORY: Update ALL Affected Context Files**
   - **feature-tracking.md**: Update status, add new implementation tokens, modify description
   - **specification.md**: Modify user-facing behavior if applicable
   - **requirements.md**: Update requirements and acceptance criteria
   - **architecture.md**: Update technical implementation details
   - **testing.md**: Update test coverage and strategies
   - **Any other files**: Update any document that references the feature

4. **🔄 Synchronize Changes**
   - Ensure all documents reflect the same change
   - Update cross-references between documents
   - Verify feature status is consistent everywhere

5. **⚡ Update Implementation Tokens**
   - Add tokens to newly modified code
   - Update existing tokens if descriptions change
   - Remove tokens from deleted code

6. **📋 Validation**
   - Run documentation validation if available
   - Check all cross-references manually
   - Verify test coverage reflects changes

#### **🗑️ Removing Features (STEP-BY-STEP)**

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **📢 Deprecation Notice**
   - Mark feature as deprecated in `specification.md`
   - Add deprecation notice with timeline
   - Update user documentation with migration guidance

2. **🔒 Backward Compatibility Check**
   - Verify `immutable.md` requirements are preserved
   - Ensure no breaking changes for existing users
   - Document compatibility preservation measures

3. **🧪 Test Maintenance**
   - Keep tests until complete removal
   - Mark tests as deprecated but functional
   - Plan test removal timeline

4. **🚨 MANDATORY: Documentation Cleanup**
   - **feature-tracking.md**: Mark as deprecated/removed, update status
   - **specification.md**: Remove or mark as deprecated
   - **requirements.md**: Mark requirements as deprecated
   - **architecture.md**: Update architecture to reflect removal
   - **testing.md**: Update test coverage documentation
   - Remove from all other context files simultaneously

5. **📋 Feature Matrix Update**
   - Update status to "Deprecated" or "Removed"
   - Preserve entry for historical tracking
   - Update implementation tokens to reflect removal

### 🎯 CONTEXT FILE UPDATE TEMPLATES

#### **For New Features:**
```
1. Open feature-tracking.md → Add feature to appropriate table
2. Open specification.md → Add to relevant command/config section
3. Open requirements.md → Add to functional requirements section
4. Open architecture.md → Add to relevant component section
5. Open testing.md → Add to test coverage section
6. Add implementation tokens to code
7. Update all cross-references
```

#### **For Feature Modifications:**
```
1. Search all context files for feature references
2. Update feature-tracking.md status and tokens
3. Update specification.md behavior descriptions
4. Update requirements.md if requirements change
5. Update architecture.md if implementation changes
6. Update testing.md if test strategy changes
7. Update implementation tokens in modified code
8. Verify all cross-references remain valid
```

#### **For Bug Fixes:**
```
1. Check if bug affects documented behavior
2. If YES: Update specification.md to clarify correct behavior
3. Update feature-tracking.md if implementation tokens change
4. Update testing.md if new tests are added
5. Add implementation tokens to fix code
```

### ⚠️ COMMON MISTAKES TO AVOID

1. **❌ Updating code without updating context files**
   - ALWAYS update context files BEFORE or DURING code changes

2. **❌ Forgetting to check immutable.md**
   - ALWAYS verify no immutable requirements are violated

3. **❌ Inconsistent feature status across documents**
   - ALWAYS ensure all documents show the same feature status

4. **❌ Missing implementation tokens**
   - ALWAYS add tokens to every modified function/method

5. **❌ Broken cross-references**
   - ALWAYS verify links between documents after changes

6. **❌ Incomplete impact analysis**
   - ALWAYS search all context files for feature references

### 🚀 ENFORCEMENT REMINDERS

- **NO CODE REVIEW** without context file updates
- **NO MERGE** without feature tracking matrix updates  
- **NO DEPLOYMENT** without documentation validation
- **NO EXCEPTIONS** - context files are as important as code

## Implementation Decision Tracking

### Decision Record Format
```
Decision ID: DEC-XXX
Feature IDs: [ARCH-001, CFG-002]
Date: YYYY-MM-DD
Decision: [Brief description]
Rationale: [Why this approach was chosen]
Alternatives: [Other options considered]
Impact: [Effect on existing features]
Code Markers: [Specific implementation tokens]
```

### Current Decisions

#### DEC-001: ZIP Archive Format
- **Feature IDs**: ARCH-001, ARCH-002, ARCH-003
- **Date**: 2024-03-01 (estimated)
- **Decision**: Use ZIP format for all archive operations
- **Rationale**: 
  - Cross-platform compatibility
  - Built-in compression
  - Wide tool support
  - Standard library support in Go
- **Alternatives**: tar.gz, tar.bz2, custom format
- **Impact**: All archive creation, listing, and verification must support ZIP
- **Code Markers**: `archive/zip` imports, `*.zip` file extensions

#### DEC-002: YAML Configuration
- **Feature IDs**: CFG-001, CFG-002, CFG-003
- **Date**: 2024-03-01 (estimated)
- **Decision**: Use YAML for configuration files with `.bkpdir.yml` naming
- **Rationale**:
  - Human-readable format
  - Supports complex nested structures
  - Good Go library support
  - Industry standard for configuration
- **Alternatives**: JSON, TOML, INI
- **Impact**: All configuration must be YAML-compatible
- **Code Markers**: `gopkg.in/yaml.v3` imports, `.yml` file extensions

#### DEC-003: Dual Printf/Template Formatting
- **Feature IDs**: CFG-003
- **Date**: 2024-03-01 (estimated)
- **Decision**: Support both printf-style and template-based output formatting
- **Rationale**:
  - Printf for simple, backward-compatible formatting
  - Templates for rich data extraction and advanced formatting
  - Graceful fallback from template to printf
  - Supports ANSI colors and structured output
- **Alternatives**: Printf only, template only, custom format
- **Impact**: All output must support both formatting modes
- **Code Markers**: `FormatXXX()` and `TemplateXXX()` functions, regex patterns

#### DEC-004: Structured Error Handling
- **Feature IDs**: FILE-002, ARCH-002
- **Date**: 2024-03-01 (estimated)
- **Decision**: Use structured error types with status codes and operation context
- **Rationale**:
  - Consistent error handling across operations
  - Machine-readable status codes for scripting
  - Enhanced debugging with operation context
  - Supports template-based error formatting
- **Alternatives**: Standard Go errors, error codes only, string-based errors
- **Impact**: All operations must return structured errors
- **Code Markers**: `ArchiveError`, `BackupError` types, `NewXXXError()` functions

#### DEC-005: Git Command-line Integration
- **Feature IDs**: GIT-001, GIT-002
- **Date**: 2024-03-01 (estimated)
- **Decision**: Use Git command-line interface for repository information
- **Rationale**:
  - Simplicity over Git library dependencies
  - Relies on user's Git installation
  - Consistent with user's Git configuration
  - Lightweight implementation
- **Alternatives**: Git libraries (go-git), Git C libraries
- **Impact**: Git functionality requires Git CLI installation
- **Code Markers**: `exec.Command("git", ...)` calls, Git command parsing

#### DEC-006: Resource Management with Cleanup
- **Feature IDs**: ARCH-002, FILE-002
- **Date**: 2024-03-01 (estimated)
- **Decision**: Implement ResourceManager for automatic cleanup with panic recovery
- **Rationale**:
  - Prevents resource leaks on operation failure
  - Handles panic scenarios gracefully
  - Thread-safe for concurrent operations
  - Simplifies error handling in operations
- **Alternatives**: Manual cleanup, defer-only cleanup, no cleanup
- **Impact**: All file operations must use ResourceManager
- **Code Markers**: `ResourceManager` type, `CleanupWithPanicRecovery()` calls

#### DEC-007: Context-Aware Operations
- **Feature IDs**: ARCH-002, FILE-002
- **Date**: 2024-03-01 (estimated)
- **Decision**: Support context cancellation for long-running operations
- **Rationale**:
  - Enables operation timeouts
  - Supports graceful cancellation
  - Better user experience for long operations
  - Standard Go pattern for cancellable operations
- **Alternatives**: No cancellation, timeout-only, signal-based cancellation
- **Impact**: Core operations must accept and check context
- **Code Markers**: `context.Context` parameters, `checkContextCancellation()` calls

#### DEC-008: Atomic File Operations
- **Feature IDs**: ARCH-002, FILE-002
- **Date**: 2024-03-01 (estimated)
- **Decision**: Use temporary files with atomic rename for all file creation
- **Rationale**:
  - Prevents corruption on operation failure
  - Ensures consistency during concurrent access
  - Standard pattern for safe file operations
  - Integrates with resource cleanup
- **Alternatives**: Direct file writing, lock files, database-style transactions
- **Impact**: All file creation must use temporary files
- **Code Markers**: `.tmp` file extensions, `os.Rename()` calls

#### GIT-006: Configurable Git Dirty Status
- **Status**: COMPLETED
- **Priority**: High
- **Description**: Make the Git dirty status indicator ('-dirty' suffix) configurable through the configuration system, allowing users to enable/disable this feature.
- **Requirements**: 
  - Add configuration option to enable/disable dirty status indicator
  - Maintain backward compatibility with existing behavior
  - Update documentation to reflect new configuration option
  - Add test coverage for new configuration option
- **Implementation Areas**:
  - Configuration system (add new option)
  - Git status detection (respect configuration)
  - Archive naming (conditional dirty status)
  - Documentation updates
- **Files Modified**:
  - `config.go` - Added `show_git_dirty_status` option to Config struct
  - `archive.go` - Updated archive naming to respect configuration
- **Implementation Details**:
  - Added `show_git_dirty_status` boolean option to Config struct
  - Default to true for backward compatibility
  - Updated Git status detection to respect configuration
  - Added configuration documentation
  - Added test coverage for new option
- **Testing**: 
  - Test with configuration enabled/disabled
  - Verify backward compatibility
  - Test archive naming with both settings
- **Implementation Notes**:
  - Added `show_git_dirty_status` to Config struct with YAML tag
  - Set default value to true to maintain backward compatibility
  - Updated `ArchiveNameConfig` struct to include the new option
  - Modified archive naming functions to check the option before adding "-dirty" suffix
  - Added proper merging of the new option in `mergeBasicSettings`
  - No changes needed to Git status detection as it still provides the status, just not always shown

#### CFG-004: Eliminate Hardcoded Strings and Enhance Template Formatting
- **Status**: COMPLETED
- **Priority**: High
- **Description**: Implement comprehensive string externalization system allowing all user-facing strings to be loaded from configuration files rather than hardcoded, with enhanced template formatting using named elements in data structures.
- **Requirements**: 
  - All user-facing strings must be configurable through YAML files
  - Support both printf-style (%s) and template-based (%{name}) formatting
  - Named elements in data structures for template formatting
  - Backward compatibility with existing configurations
- **Implementation Areas**:
  - Archive operation messages (no archives found, verification results, configuration updates, dry-run operations)
  - Backup operation messages (no backups found, backup creation, identical detection)
  - Error handling messages (disk full, permission denied, file/directory not found, validation errors)
  - Configuration management (file paths, updates)
  - Template-based formatting with named data structure elements
- **Files Modified**: 
  - `config.go` - Extended Config struct with comprehensive format strings
  - `formatter.go` - Added OutputFormatter methods for all new format strings
  - `errors.go` - Updated error handling to use configurable messages
- **Implementation Details**:
  - Added 12 new error message format strings (FormatDiskFullError, FormatPermissionError, etc.)
  - Added corresponding template versions with named data structure support
  - Updated HandleArchiveError to use configurable error messages
  - Maintained full backward compatibility with existing format strings
  - All error messages now use OutputFormatter methods instead of hardcoded strings
- **Testing**: Manual compilation successful, all existing functionality preserved

#### OUT-001: Delayed Output Management
- **Status**: COMPLETED
- **Priority**: Medium
- **Description**: Implement delayed output functionality by returning output messages to calling functions instead of direct stdout/stderr printing, enabling better control over when and how output is displayed.
- **Requirements**: 
  - Return formatted messages from formatter methods instead of direct printing
  - Maintain backward compatibility with existing Print methods
  - Support buffering of multiple output messages
  - Enable calling functions to control output timing and destination
- **Implementation Areas**:
  - OutputFormatter methods - modify to return strings instead of direct printing
  - Command handlers - collect and manage output messages
  - Error handling - return error messages for delayed display
  - Archive and backup operations - return operation result messages
- **Design Decisions**:
  - Dual-mode approach: keep existing Print methods for backward compatibility
  - Add new Format methods that return strings without printing
  - Calling functions decide when to display collected messages
  - Preserve current error vs stdout routing behavior
- **Files Modified**:
  - `formatter.go` - Added OutputCollector and OutputMessage types, enhanced OutputFormatter with delayed output support
- **Implementation Details**:
  - Added OutputMessage struct with Content, Destination, and Type fields
  - Added OutputCollector with methods: AddStdout, AddStderr, FlushAll, FlushStdout, FlushStderr, Clear
  - Enhanced OutputFormatter with optional collector field and delayed mode support
  - Added NewOutputFormatterWithCollector constructor for delayed output mode
  - Updated all Print methods to check for collector and route messages accordingly
  - Maintained full backward compatibility - existing code works unchanged
  - Added utility methods: IsDelayedMode, GetCollector, SetCollector
- **Testing**: Manual compilation successful, application runs correctly with delayed output functionality
- **Implementation Tokens**: `// OUT-001: Delayed output`

#### TEST-001: Comprehensive formatter.go Test Coverage
- **Status**: COMPLETED
- **Priority**: Medium
- **Description**: Test all 0% coverage functions in formatter.go to improve code coverage and ensure reliability of OutputCollector, template methods, and error formatting functionality.
- **Requirements**: 
  - Test OutputCollector functionality (AddStdout, AddStderr, GetMessages, Clear, etc.)
  - Test template formatting methods with 0% coverage
  - Test error formatting functions (Format/Template error methods)
  - Test Print methods in delayed output mode
  - Verify correct message routing (stdout vs stderr)
- **Implementation Areas**:
  - OutputCollector methods - comprehensive testing of delayed output functionality
  - Template formatting methods - test TemplateIdenticalArchive, TemplateDryRunArchive, etc.
  - Error formatting - test all Format*Error and Template*Error methods
  - Print methods - test delayed mode behavior and message collection
  - Pattern extraction - test ExtractConfigLineData and other extraction methods
- **Files Modified**: 
  - `formatter_test.go` - Added comprehensive test suites for previously untested functions
- **Implementation Details**:
  - Added TestOutputCollector with tests for NewOutputCollector, AddStdout, AddStderr, GetMessages, Clear
  - Added TestDelayedOutputMode testing OutputFormatter with collector integration
  - Added TestTemplateFormattingMethods for all template functions with 0% coverage
  - Added TestErrorFormattingMethods for error formatting functions
  - Added TestPrintMethods testing Print methods in delayed mode
  - Added TestTemplateFormatterAdvanced for TemplateFormatter functionality
  - Added TestOutputCollectorFlushMethods for comprehensive testing of FlushAll, FlushStdout, FlushStderr
  - Added TestPrintMethodsDelayedMode for complete delayed output mode coverage
  - Eliminated ALL 0% coverage functions in formatter.go (152 total functions tracked)
  - Achieved 116 functions at 100% coverage, 36 functions with partial coverage (75%+ typical)
- **Testing**: All tests pass successfully, comprehensive coverage of OutputCollector, template methods, and error formatting
- **Implementation Notes**: 
  - Template error methods return template strings with %v placeholders rather than formatted error text
  - Print methods route messages correctly between stdout/stderr based on message type
  - FlushAll, FlushStdout, FlushStderr methods tested with stdout/stderr capture techniques
  - Successfully eliminated all zero-coverage functions in formatter.go through targeted testing
  - Delayed output mode now fully tested with proper message routing verification

#### DEC-009: Documentation Enhancement Framework
- **Feature IDs**: DOC-001, DOC-002, DOC-003, DOC-004, DOC-005
- **Date**: 2024-12-19
- **Decision**: Implement comprehensive documentation enhancement system with semantic linking, synchronization, and traceability
- **Rationale**:
  - Cross-document repetition identified as highly valuable for LLM consumption
  - Need automated systems to maintain consistency during changes
  - Enhanced traceability prevents functionality loss during evolution
  - Semantic linking improves understanding and change impact analysis
- **Alternatives**: Manual processes, simple validation scripts, external documentation tools
- **Impact**: Establishes foundation for maintaining documentation quality at scale
- **Code Markers**: `// DOC-XXX: Documentation enhancement` tokens

#### DEC-010: Semantic Cross-Referencing Strategy
- **Feature IDs**: DOC-001
- **Date**: 2024-12-19
- **Decision**: Use bi-directional linking with feature reference format including all document layers
- **Rationale**:
  - LLMs benefit from rich cross-references to understand relationships
  - Bi-directional links prevent orphaned references
  - Standardized format ensures consistency across documents
  - Forward/backward/sibling links provide comprehensive navigation
- **Alternatives**: Simple hyperlinks, external reference management, manual cross-referencing
- **Impact**: All features must be consistently referenced across all documents
- **Code Markers**: LinkingEngine implementation, automated link validation

#### DEC-011: Documentation Synchronization Approach
- **Feature IDs**: DOC-002
- **Date**: 2024-12-19
- **Decision**: Use synchronized checkpoints with automated validation and change templates
- **Rationale**:
  - Prevents documentation drift during rapid development
  - Automated validation catches inconsistencies early
  - Change templates ensure systematic updates across all layers
  - Clear conflict resolution hierarchy prevents documentation disagreements
- **Alternatives**: Manual synchronization, version-based sync, external sync tools
- **Impact**: All documentation changes must follow synchronized update process
- **Code Markers**: SyncFramework, pre-commit hooks, validation scripts

#### DEC-012: Enhanced Traceability Design
- **Feature IDs**: DOC-003
- **Date**: 2024-12-19
- **Decision**: Implement feature fingerprints with behavioral contracts and dependency mapping
- **Rationale**:
  - Feature fingerprints ensure stable identity through changes
  - Behavioral contracts define what cannot change without version bump
  - Dependency mapping shows change impact chains
  - Automated regression prevention protects against functionality loss
- **Alternatives**: Simple feature tracking, manual impact analysis, external traceability tools
- **Impact**: All features must have behavioral contracts and dependency mappings
- **Code Markers**: TraceabilitySystem, behavioral contract validation, dependency analysis

## Code Marker Strategy

### Recommended Tokens
- `// FEATURE-ID: Brief description` - Primary feature implementation
- `// IMMUTABLE-REF: [section]` - Links to immutable requirements
- `// TEST-REF: [test function]` - Links to test coverage
- `// DECISION-REF: DEC-XXX` - Links to implementation decisions

### Example Implementation
```go
// ARCH-001: Archive naming convention implementation
// IMMUTABLE-REF: Archive Naming Convention
// TEST-REF: TestGenerateArchiveName
// DECISION-REF: DEC-001
func GenerateArchiveName(prefix string, timestamp time.Time, gitInfo *GitInfo, note string) string {
    // Implementation here
}
```

## Implementation Status Summary

### Completed (Phase 1) ✅
- ✅ Feature matrix created with cross-references
- ✅ Implementation tokens added to core functions (66 tokens across 5 files)
- ✅ Decision records documented with rationale (8 decisions)
- ✅ Code markers linked to documentation
- ✅ Validation script implemented and tested
- ✅ Documentation consistency framework established

### Completed (Phase 2) ✅
- ✅ **Fix Missing Test Functions** - Successfully eliminated all 9 validation errors
- ✅ **Validation Script Enhancement** - Fixed false positives in test function detection
- ✅ **Zero-Error Validation Achieved** - Validation script now passes with 0 errors, 0 warnings
- ✅ **Complete Implementation Token Coverage** - Achieved 144 total tokens with comprehensive coverage
- ✅ **Strategic Test Token Enhancement** - Added targeted tokens to key feature validation functions

### Phase 3: Pre-Extraction Refactoring (CRITICAL FOUNDATION - MUST COMPLETE FIRST)

#### **⚠️ CRITICAL PATH: NO EXTRACTION WITHOUT REFACTORING**
**All EXTRACT-001 through EXTRACT-010 tasks are BLOCKED until refactoring phase completes successfully.**

**Week 0: Foundation Preparation (CRITICAL)**
1. **REFACTOR-001: Dependency Analysis and Interface Standardization** - **CRITICAL BLOCKER**
2. **REFACTOR-002: Large File Decomposition Preparation** - **✅ COMPLETED (2025-01-02)**
   - ✅ **Component Boundary Analysis**: Identified 5 distinct components in 1677-line formatter.go
   - ✅ **OutputCollector Component**: Ready for immediate extraction (lines 20-111) - zero dependencies
   - ✅ **Printf Formatter Component**: Identified with config dependency (lines 120-610)
   - ✅ **Template Formatter Component**: Identified with config dependency (lines 637-928)
   - ✅ **Extended Output Formatters**: Complex operations component (lines 929-1351)
   - ✅ **Error Formatting Component**: Specialized error handling (lines 1352-1677)
   - ✅ **Internal Interfaces Designed**: FormatProvider, OutputDestination, PatternExtractor, FormatterInterface, TemplateFormatterInterface
   - ✅ **Documentation Created**: Comprehensive `docs/formatter-decomposition.md` with extraction strategy
   - ✅ **Implementation Tokens Added**: REFACTOR-002 tokens throughout formatter.go marking boundaries
   - ✅ **Validation Successful**: Compilation successful, all tests passing, 168 tests with zero failures
   - ✅ **Backward Compatibility**: Functional integrity preserved, no breaking changes
   - **Status**: Ready for EXTRACT-003 (Output Formatting System) - requires config interface abstraction first
3. **REFACTOR-003: Configuration Schema Abstraction** - **HIGH PRIORITY**

**Week 0.5: Consolidation and Validation (HIGH PRIORITY)**
4. **REFACTOR-004: Error Handling and Resource Management Consolidation** - **MEDIUM PRIORITY**
5. **REFACTOR-005: Code Structure Optimization for Extraction** - **MEDIUM PRIORITY**
6. **REFACTOR-006: Refactoring Impact Validation** - **HIGH PRIORITY**

#### **🎯 REFACTORING SUCCESS GATES**
Before proceeding to Phase 4 (Component Extraction), ALL of these must be completed:
- ✅ Complete dependency analysis with zero circular dependency risks
- ✅ Formatter decomposition strategy validated (REFACTOR-002 COMPLETED)
- ❌ Configuration abstraction interfaces defined
- ❌ Error handling patterns standardized
- ❌ All refactoring changes validated with zero test failures
- ❌ Pre-extraction validation checklist passed

### **Phase 4: Component Extraction and Generalization (BLOCKED - READY AFTER PHASE 3)**

#### **🚫 EXTRACTION AUTHORIZATION REQUIRED**
**These tasks require explicit authorization after Phase 3 completion:**

**Weeks 1-2: Core Infrastructure Extraction (IMMEDIATE - After Refactoring)**
- **EXTRACT-001: Configuration Management System** (CRITICAL) - Requires REFACTOR-003 completion
- **EXTRACT-002: Error Handling and Resource Management** (CRITICAL) - Requires REFACTOR-004 completion

**Weeks 3-4: User Experience Layer (HIGH PRIORITY)**
- **EXTRACT-003: Output Formatting System** (HIGH) - Requires REFACTOR-002 completion
- **EXTRACT-004: Git Integration System** (MEDIUM) - Ready after dependency analysis

**Weeks 5-6: CLI Framework (HIGH PRIORITY)**  
- **EXTRACT-005: CLI Command Framework** (HIGH) - Requires core infrastructure completion
- **EXTRACT-006: File Operations and Utilities** (HIGH) - Ready after core infrastructure

**Weeks 7-8: Application Patterns and Quality (HIGH PRIORITY)**
- **EXTRACT-007: Data Processing Patterns** (MEDIUM) - Ready after framework completion
- **EXTRACT-008: CLI Application Template** (HIGH) - Requires all components
- **EXTRACT-009: Testing Patterns and Utilities** (HIGH) - Critical for quality
- **EXTRACT-010: Package Documentation and Examples** (HIGH) - Essential for adoption

#### **📋 EXTRACTION DEPENDENCY CHAIN**
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

### **⚠️ REVISED TIMELINE AND CRITICAL PATH**

#### **UPDATED TIMELINE (Total: 9 weeks)**
- **Week 0**: REFACTOR-001, REFACTOR-002, REFACTOR-003 (Parallel)
- **Week 0.5**: REFACTOR-004, REFACTOR-005, REFACTOR-006 (Sequential validation)
- **Week 1**: Extraction authorization and EXTRACT-001, EXTRACT-002 start
- **Weeks 1-8**: Original extraction timeline proceeds as planned

#### **CRITICAL SUCCESS FACTORS**
1. **Zero Compromise on Refactoring Quality**: All refactoring tasks must be completed to full specification
2. **Comprehensive Validation**: Each refactoring step must be validated before proceeding
3. **Interface Stability**: All interfaces must be finalized before extraction begins
4. **Dependency Clarity**: Complete dependency mapping prevents extraction failures
5. **Backward Compatibility**: All refactoring must preserve existing functionality

#### **RISK MITIGATION**
- **Parallel Refactoring**: REFACTOR-001, 002, 003 can proceed in parallel with coordination
- **Incremental Validation**: REFACTOR-006 validates each step to prevent compound failures  
- **Rollback Plan**: Each refactoring step includes rollback procedures if validation fails
- **Extraction Gate**: Hard stop before extraction until all criteria met

**🎯 FINAL RECOMMENDATION: PROCEED WITH REFACTORING PHASE IMMEDIATELY**

The codebase has excellent test coverage (73.5%) and comprehensive testing infrastructure, making it ideal for safe refactoring. The pre-extraction refactoring will ensure:
- **Clean Architecture**: Well-defined component boundaries
- **Maintainable Code**: Reduced complexity and improved organization  
- **Reliable Extraction**: Zero risk of circular dependencies or architectural issues
- **Future-Proof Design**: Extracted components will be robust and reusable

**Start with REFACTOR-001 (Dependency Analysis) immediately to begin the foundation for successful component extraction.**

## Prioritized Task List

### **Phase 2: Process Establishment and Validation (HIGH PRIORITY)**

#### **Immediate Tasks (Next 1-2 weeks)**

**🚨 CRITICAL FOUNDATION TASK:**
1. **Establish Solid Documentation Foundation** (TOP PRIORITY) ✅ **COMPLETED**
   - [x] **Eliminate the 9 validation errors** - Fix missing test functions to achieve zero-error validation ✅
   - [x] **Complete implementation token coverage** - Add tokens to remaining 5 core files ✅
   - [x] **Verify validation script passes cleanly** - Ensure solid foundation for automation ✅
   - **Rationale**: This establishes the solid foundation required for advanced automation and monitoring capabilities in later phases
   - **Status**: Foundation established with zero validation errors achieved!

**🔧 SPECIFIC IMPLEMENTATION TASKS:**
2. **Fix Missing Test Functions** (Critical - 9 validation errors) ✅ **COMPLETED**
   - [x] Add `TestCompareFiles` for FILE-003 feature
   - [x] Add `TestCreateFullArchive` for ARCH-002 feature  
   - [x] Add `TestCreateIncremental` for ARCH-003 feature
   - [x] Add `TestGetConfigSearchPath` for CFG-001 feature
   - [x] Add `TestGitIntegration` for GIT-001 feature
   - [x] Add `TestGitNaming` for GIT-002 feature
   - [x] Add `TestTemplateFormatter` for CFG-003 feature
   - [x] Clean up test reference parsing in feature-tracking.md
   - **Implementation Notes**: 
     - Created new test files: `comparison_test.go`, `git_test.go`, `formatter_test.go`
     - Added missing test functions to existing `archive_test.go` and `config_test.go`
     - All tests pass successfully with comprehensive coverage of core functionality
     - Fixed test isolation issues in comparison tests by using separate temporary directories

2.1. **Add Test Coverage for Tools Directory** (TEST-002) - **HIGH PRIORITY** ✅ **COMPLETED**
   - [x] **Create `tools/coverage_test.go`** - Add comprehensive test coverage for the coverage analyzer tool
   - [x] **Create `tools/validate_coverage_test.sh`** - Add test validation for the coverage validation script
   - [x] **Test coverage profile parsing** - Verify correct parsing of Go coverage profiles
   - [x] **Test coverage classification logic** - Verify legacy vs new code classification
   - [x] **Test coverage report generation** - Validate HTML and text report generation
   - [x] **Test coverage threshold validation** - Ensure proper enforcement of coverage thresholds
   - [x] **Test error handling scenarios** - Cover file not found, malformed profiles, permission errors
   - [x] **Add integration tests for script execution** - Test end-to-end workflow of coverage validation
   - **Rationale**: The tools directory currently shows `[no test files]` in test output, creating gaps in code quality validation
   - **Status**: Not Started
   - **Priority**: High - Tools are critical for maintaining code quality standards
   - **Implementation Areas**:
     - Coverage analysis tool (`tools/coverage.go`) - COV-001 implementation
     - Coverage validation script (`tools/validate-coverage.sh`) - COV-001 validation
     - Test infrastructure for shell script testing
     - Integration testing for coverage workflow
   - **Dependencies**: Requires existing test infrastructure (TEST-INFRA-001 components)
   - **Implementation Tokens**: `// TEST-002: Tools directory testing`
   - **Expected Outcomes**:
     - Eliminate `[no test files]` warning in test output
     - Achieve 85%+ coverage for tools directory code
     - Validate reliability of coverage analysis and validation tools
     - Ensure proper error handling in critical quality assurance tools
     - **Coverage Achievement**: Eliminated the `[no test files]` warning and achieved comprehensive test coverage for critical quality assurance tools
     - **Test Integration**: Successfully integrated with existing test infrastructure (TEST-INFRA-001 components) for robust testing patterns

3. **Review Git Integration Features** (HIGH PRIORITY - Feature Analysis) ✅ **COMPLETED**
   - [x] Review existing Git integration features (GIT-001, GIT-002)
   - [x] Analyze implementation against requirements
   - [x] Identify missing features and gaps
   - [x] Document findings and recommendations
   - [x] Implement GIT-003 (Git status detection)
   - **Rationale**: Ensure comprehensive Git integration coverage and identify any missing features
   - **Status**: ✅ **COMPLETED** - All planned features implemented and tested
   - **Implementation Notes**:
     - **Completed Features**:
       - GIT-001: Git info extraction (repository detection, branch/hash extraction)
       - GIT-002: Branch/hash naming (integration with archive naming)
       - GIT-003: Git status detection (clean/dirty working directory)
     - **Remaining Features**:
       - GIT-004: Git submodule support
       - GIT-005: Git configuration integration
     - **Next Steps**:
       - Consider implementing GIT-004 (Git submodule support)
       - Plan for GIT-005 (Git configuration integration)

4. **Implement Git Status Detection** (HIGH PRIORITY - Feature Implementation) ✅ **COMPLETED**
   - [x] Add `IsGitWorkingDirectoryClean` function to detect clean/dirty state
   - [x] Add status information to archive naming when enabled
   - [x] Add comprehensive test coverage for status detection
   - [x] Update documentation to reflect new functionality
   - **Rationale**: Enhance Git integration by detecting working directory state
   - **Status**: ✅ **COMPLETED** - Implementation finished
   - **Implementation Notes**:
     - Used `git status --porcelain` for reliable status detection
     - Integrated with existing Git info extraction
     - Added status indicator to archive names when enabled
     - Maintained backward compatibility with existing features
   - **Next Steps**:
     - Consider implementing GIT-004 (Git submodule support)
     - Plan for GIT-005 (Git configuration integration)

5. **Implement Code Coverage Exclusion for Existing Code** (COV-001) - **HIGH PRIORITY**
   - [x] **Create coverage exclusion configuration** - Add build tags or comments to exclude legacy code from coverage metrics
   - [x] **Establish coverage baseline** - Document current coverage levels for existing codebase before exclusion
   - [x] **Implement selective coverage reporting** - Configure Go tools to focus on new/modified code only
   - [x] **Add coverage validation for new code** - Ensure new development maintains high coverage standards
   - [x] **Update Makefile coverage targets** - Modify existing `test-coverage` target to support exclusion patterns
   - **Rationale**: Focus coverage metrics on new development while preserving testing of existing functionality
   - **Status**: ✅ **COMPLETED** - Comprehensive coverage exclusion system implemented successfully
   - **Priority**: High - Essential for maintaining quality standards on new code without penalizing legacy code ✅ **SATISFIED**
   - **Implementation Areas**:
     - Build system (Makefile modification) ✅ **COMPLETED**
     - Test configuration (Go coverage tools) ✅ **COMPLETED**
     - CI/CD integration (coverage reporting) ✅ **COMPLETED**
     - Documentation (coverage standards) ✅ **COMPLETED**
   - **Implementation Notes**:
     - ✅ **Coverage Configuration**: Created `coverage.toml` with comprehensive exclusion patterns for legacy code (main.go, config.go, formatter.go, backup.go, archive.go)
     - ✅ **Baseline Documentation**: Established baseline coverage levels (73.5% main package, 75.6% internal/testutil) in `docs/coverage-baseline.md`
     - ✅ **Selective Reporting Tool**: Implemented `tools/coverage.go` to parse coverage profiles and apply exclusion patterns with 85% threshold for new code
     - ✅ **Validation Framework**: Created `tools/validate-coverage.sh` script for comprehensive coverage validation with quality gates
     - ✅ **Makefile Integration**: Added new targets `test-coverage-new`, `test-coverage-validate`, `coverage-check`, and `dev-full` for development workflow
     - ✅ **Design Decisions**: 
       - Used file-based exclusion rather than build tags for simplicity and visibility
       - Maintained all existing test execution while hiding legacy code from coverage metrics  
       - Established 85% coverage threshold for new code and 70% for modified legacy code
       - Created separate HTML reports for overall vs new code coverage analysis
       - Implemented quality gates that fail builds when new code coverage is below threshold
     - ✅ **Testing Results**: System successfully identifies no new code in current state, reports overall 73.5% coverage, and passes all validation checks
     - ✅ **Files Created**:
       - `coverage.toml` - Coverage configuration with exclusion patterns
       - `docs/coverage-baseline.md` - Baseline documentation and analysis 
       - `tools/coverage.go` - Go tool for selective coverage reporting
       - `tools/validate-coverage.sh` - Comprehensive validation script
       - Updated `Makefile` with new coverage targets

6. **Establish Coverage Baseline and Selective Reporting** (COV-002) - **MEDIUM PRIORITY**
   - [x] **Document current coverage metrics** - Capture baseline coverage percentages for all existing files
   - [x] **Create coverage configuration file** - Define which files/functions should be excluded from coverage reporting
   - [x] **Implement coverage differential reporting** - Show coverage changes only for modified code in PRs
   - [x] **Add coverage trend tracking** - Monitor coverage evolution over time for new development
   - [x] **Create coverage quality gates** - Define minimum coverage thresholds for new/modified code
   - **Rationale**: Establish measurement framework for code quality while avoiding legacy code penalties
   - **Status**: ✅ **COMPLETED**
   - **Dependencies**: Requires COV-001 to be implemented first ✅ **SATISFIED**
   - **Implementation Notes**:
     - ✅ **Comprehensive baseline documentation**: Created detailed coverage baseline in `docs/coverage-baseline.md` with current metrics (73.3% overall, 73.5% main, 75.6% testutil)
     - ✅ **Enhanced configuration system**: Extended `coverage.toml` with COV-002 sections for quality gates, differential reporting, trend tracking, and integration settings
     - ✅ **Differential reporting tool**: Implemented `tools/coverage-differential.go` with Git integration, baseline comparison, quality gate validation, and HTML/JSON report generation
     - ✅ **Trend tracking system**: Added coverage history tracking in `docs/coverage-history.json` with automatic baseline updates
     - ✅ **Quality gates framework**: Implemented configurable thresholds (70% new code, 80% critical paths) with automatic validation and recommendations
     - ✅ **Makefile integration**: Added 5 new targets: `test-coverage-baseline`, `test-coverage-differential`, `test-coverage-trends`, `test-coverage-full`, `test-coverage-quality-gates`
     - ✅ **Report generation**: Automated HTML and JSON report generation in `coverage_reports/` directory
     - **Design Decisions**: 
       - Used TOML configuration for flexibility and human readability
       - Implemented Git integration for automatic modified file detection
       - Created separate build process for tools to avoid main function conflicts
       - Focused on modified files only to avoid penalizing legacy code
       - Added comprehensive recommendation system for coverage improvements

7. **Configure Advanced Coverage Controls** (COV-003) - **FUTURE ENHANCEMENT**
   - [ ] **Implement function-level coverage exclusion** - Granular control over which functions are included in coverage metrics
   - [ ] **Add coverage comment directives** - Support `//coverage:ignore` style annotations for specific code blocks
   - [ ] **Create coverage exception documentation** - Maintain list of excluded code with justification
   - [ ] **Integrate with development workflow** - Automatic coverage analysis in pre-commit hooks and CI
   - [ ] **Add coverage visualization tools** - Generate HTML reports with exclusion highlighting
   - **Rationale**: Provide fine-grained control over coverage reporting while maintaining comprehensive testing
   - **Status**: Not Started
   - **Priority**: Low - Advanced feature for mature coverage management
   - **Design Decisions**:
     - Comment-based exclusion provides visibility and intentionality
     - Function-level control allows excluding difficult-to-test error paths
     - Documentation ensures excluded code decisions are preserved and reviewable
     - Visualization helps developers understand coverage impact

### **🚨 CRITICAL PRE-EXTRACTION TESTING PHASE (IMMEDIATE - BLOCKS ALL EXTRACTION)**

**Testing Priority Summary:**
- **270 completely untested functions** identified in coverage analysis
- **47.2% overall coverage** - insufficient for extraction project
- **Critical shared code** in extraction targets has significant testing gaps
- **Testing must be completed** before any extraction work begins

8. **Comprehensive Testing of Extraction Target Functions** (TEST-EXTRACT-001) - **CRITICAL BLOCKER**
   - [x] **Test all 0% coverage functions in archive.go** - CreateArchiveWithContext, verifyArchive, CreateFullArchiveWithCleanup
   - [x] **Test all 0% coverage functions in backup.go** - Error handling, Enhanced functions, Context operations
   - [x] **Test all 0% coverage functions in config.go** - GetConfigValuesWithSources, determineConfigSource, mergeExtended*
   - [x] **Test all 0% coverage functions in errors.go** - ArchiveError methods, error classification, atomic operations
   - [x] **Test all 0% coverage functions in formatter.go** - OutputCollector, template methods, error formatting
   - [x] **Test all 0% coverage functions in comparison.go** - CreateArchiveSnapshot, tree summaries, verification
   - [x] **Test all 0% coverage functions in main.go** - Command handlers, enhanced CLI functions ✅ **COMPLETED**
   - **Rationale**: Cannot extract untested code into reusable components - creates unreliable foundation
   - **Status**: ✅ **COMPLETED** (Overall coverage improved from 47.2% to 73.5%)
   - **Priority**: CRITICAL - Extraction project cannot proceed without this ✅ **UNBLOCKED**
   - **Blocking**: All EXTRACT-001 through EXTRACT-010 tasks ✅ **UNBLOCKED**
   - **Implementation Notes**:
     - Current overall coverage improved from 47.2% to 73.5%
     - Focus on functions that will be shared across multiple applications
     - Prioritize error paths and edge cases that are difficult to test after extraction
     - Establish performance baselines before refactoring begins
     - **errors.go COMPLETED**: Achieved 100% coverage for most functions (87.5-100% for all), comprehensive test suite with 700+ lines covering ArchiveError methods, error classification, resource management, context operations, atomic file operations, and edge cases including panic recovery and permission scenarios
     - **comparison.go COMPLETED**: Achieved excellent coverage improvement (88-100% for all previously 0% functions), comprehensive test suite with 1000+ lines covering CreateArchiveSnapshot, archive file hashing, directory-to-archive comparison, archive discovery, tree summaries, and verification functions. Created robust test helpers for ZIP archive creation and directory structures. Performance benchmarks established for critical functions.
     - **config.go COMPLETED**: Successfully achieved 100% coverage for all previously untested functions including GetConfigValuesWithSources, determineConfigSource, createSourceDeterminer, getBasicConfigValues, getStatusCodeValues, getVerificationValues, mergeExtendedFormatStrings, and mergeExtendedTemplates. Comprehensive test suite covers configuration value extraction with source tracking, config file detection, value comparison logic, and template/format string merging. Tests include validation of default vs custom configuration sources, alphabetical sorting of config values, and proper handling of all configuration data types (strings, booleans, integers, slices).
     - **main.go COMPLETED**: Successfully tested all 30+ previously 0% coverage functions with comprehensive test suite of 29 new test functions. Addressed challenges including os.Exit() calls by testing component logic separately and using appropriate test skipping. Added tests for command handlers (handleConfigCommand, handleCreateCommand, handleVerifyCommand, handleVersionCommand), command creation functions (configCmd, createCmd, verifyCmd, versionCmd, backupCmd), enhanced archive functions (CreateFullArchiveEnhanced, CreateIncrementalArchiveEnhanced, VerifyArchiveEnhanced), configuration management functions (handleConfigSetCommand, loadExistingConfigData, convertConfigValue, convertBooleanValue, convertIntegerValue, updateConfigData, saveConfigData), and verification functions (verifySingleArchive, verifyAllArchives, performVerification, handleVerificationResult). Overall test coverage improved significantly with all tests passing successfully.

9. **Create Testing Infrastructure for Complex Scenarios** (TEST-INFRA-001) - **CRITICAL ENABLER**
   
   **9.1 Archive Corruption Testing Framework** (TEST-INFRA-001-A) - **HIGH PRIORITY** ✅ **COMPLETED**
   - [x] **Create controlled ZIP corruption utilities** - Systematic header/data corruption for verification testing
   - [x] **Implement corruption type enumeration** - CRC errors, header corruption, truncation, invalid central directory
   - [x] **Add corruption reproducibility** - Deterministic corruption patterns for consistent test results
   - [x] **Create archive repair detection** - Test recovery behavior from various corruption types
   - **Implementation Areas**: Test utilities for `verify.go`, `comparison.go` archive validation
   - **Files Created**: `internal/testutil/corruption.go`, `internal/testutil/corruption_test.go`
   - **Dependencies**: None (foundational) ✅ **SATISFIED**
   - **Design Decision**: Use Go's ZIP library knowledge to corrupt specific sections (local headers, central directory, file data)
   - **Status**: ✅ **COMPLETED**
   - **Implementation Notes**: 
     - **Comprehensive corruption types**: Implemented 8 different corruption types (CRC, Header, Truncate, Central Directory, Local Header, Data, Signature, Comment) with systematic approach
     - **Reproducibility achieved**: Used deterministic seeding with offset-based variation to ensure identical corruption across identical archives
     - **Detection logic**: Created CorruptionDetector that can identify different corruption types through systematic analysis
     - **Resilient testing**: Go's ZIP reader is surprisingly resilient - adjusted tests to verify corruption effects rather than complete failure
     - **Performance**: Benchmarks show CRC corruption at ~763μs and detection at ~49μs for typical archives
     - **Recovery support**: Framework tracks original bytes and corruption locations for potential recovery scenarios
     - **Edge case handling**: Handles empty archives, very small archives, and various error conditions gracefully

   **9.2 Disk Space Simulation Framework** (TEST-INFRA-001-B) - **HIGH PRIORITY**
   - [x] **Create mock filesystem with quota limits** - Controlled disk space simulation without affecting real system
   - [x] **Implement progressive space exhaustion** - Gradually reduce available space during operations
   - [x] **Add disk full error injection** - Trigger ENOSPC at specific operation points
   - [x] **Create space recovery testing** - Test behavior when disk space becomes available again
   - **Implementation Areas**: Error handling in `archive.go`, `backup.go`, atomic file operations in `errors.go`
   - **Files to Create**: `internal/testutil/diskspace.go`, `internal/testutil/diskspace_test.go`
   - **Dependencies**: TEST-INFRA-001-D (error injection framework)
   - **Design Decision**: Use filesystem interface wrapper to simulate space constraints without requiring large files
   - **Status**: ✅ **COMPLETED**
   - **Implementation Notes**: 
     - **Comprehensive disk space simulation**: Implemented 4 exhaustion modes (Linear, Progressive, Random, Immediate) with configurable rates and thresholds
     - **Mock filesystem interface**: Created FileSystemInterface abstraction with RealFileSystem and SimulatedFileSystem implementations
     - **Space constraint enforcement**: Implements quota limits, progressive exhaustion, and space recovery without requiring large test files
     - **Error injection framework**: Supports operation-specific failure points, file-path-specific injection, and space recovery points
     - **Thread-safe operations**: All operations protected with mutex for concurrent access
     - **Comprehensive statistics**: Tracks operations, failures, space changes, and file sizes for test validation
     - **Predefined scenarios**: Includes "GradualExhaustion", "ImmediateFailure", "SpaceRecovery", and "ProgressiveExhaustion" test scenarios
     - **Helper functions**: SimulateArchiveCreation and SimulateBackupOperation for realistic testing
     - **Dynamic configuration**: Support for adding/removing failure points, recovery points, and injection points at runtime
     - **Performance optimized**: Benchmarks show efficient space checking and concurrent access patterns
     - **Error type support**: Creates standard disk full (ENOSPC), quota exceeded, and device full errors
     - **Space accounting**: Proper tracking of available and used space with accurate recovery from file deletions

   **9.3 Permission Testing Framework** (TEST-INFRA-001-C) - **HIGH PRIORITY** ✅ **COMPLETED**
   - [x] **Create permission scenario generator** - Systematic permission combinations for comprehensive testing
   - [x] **Implement cross-platform permission simulation** - Handle Unix/Windows permission differences
   - [x] **Add permission restoration utilities** - Safely restore original permissions after tests
   - [x] **Create permission change detection** - Test behavior when permissions change during operations
   - **Implementation Areas**: File operations in `comparison.go`, config file handling in `config.go`, atomic operations
   - **Files Created**: `internal/testutil/permissions.go` (756 lines), `internal/testutil/permissions_test.go` (609 lines)
   - **Dependencies**: None (foundational)
   - **Design Decision**: Use temporary directories with controlled permissions rather than modifying system files
   - **Implementation Notes**: 
     - **COMPLETED** - Critical for testing permission error paths in file operations and configuration management
     - **Key Features Implemented**:
       - PermissionSimulator with systematic permission combinations (0000-0777)
       - Cross-platform support for Unix/Windows permission differences
       - Restoration utilities with atomic permission rollback
       - Permission change detection and monitoring
       - Built-in scenarios: basic_permission_denial, directory_access_denied, mixed_permissions
       - High-level PermissionTestHelper for easy testing integration
       - Thread-safe operations with mutex protection
       - Statistics tracking and error aggregation
     - **Test Coverage**: 14 comprehensive tests including benchmarks, all passing
     - **Design Patterns**: Temporary directory isolation, internal lock-free methods to avoid deadlocks

   **9.4 Context Cancellation Testing Helpers** (TEST-INFRA-001-D) - **MEDIUM PRIORITY** ✅ **COMPLETED**
   - [x] **Create controlled timeout scenarios** - Precise timing control for context cancellation testing
   - [x] **Implement cancellation point injection** - Trigger cancellation at specific operation stages
   - [x] **Add concurrent operation testing** - Test cancellation during concurrent archive/backup operations
   - [x] **Create cancellation propagation verification** - Ensure proper context propagation through operation chains
   - **Implementation Areas**: Context handling in `archive.go`, `backup.go`, long-running operations, ResourceManager cleanup
   - **Files Created**: `internal/testutil/context.go` (623 lines), `internal/testutil/context_test.go` (832 lines)
   - **Dependencies**: None (foundational) ✅ **SATISFIED**
   - **Design Decision**: Use ticker-based timing control and goroutine coordination for deterministic cancellation testing
   - **Implementation Notes**: 
     - **COMPLETED** - Comprehensive context cancellation testing infrastructure for reliable testing of context-aware operations
     - **Key Features Implemented**:
       - ContextController with precise timing control for delayed cancellation scenarios
       - CancellationManager for injection point management and statistics tracking
       - ConcurrentTestConfig for testing cancellation during concurrent operations
       - PropagationTestConfig for verifying context propagation through operation chains
       - Helper functions for common scenarios: archive, backup, and resource cleanup cancellation points
       - Comprehensive event tracking and statistics collection for test analysis
       - Integration testing utilities combining multiple cancellation scenarios
     - **Core Components**:
       - ContextController: Provides controlled cancellation with configurable delays and timeout management
       - CancellationManager: Manages injection points with enable/disable functionality and execution counting
       - PropagationChain: Tests multi-level context propagation with detailed operation tracking
       - Helper utilities: CreateTimeoutContext, CreateCancelledContext, SimulateSlowOperation, VerifyContextPropagation
     - **Test Coverage**: 18 comprehensive test functions including integration tests, stress tests, and benchmarks
     - **Design Patterns**: Ticker-based timing, goroutine coordination, atomic counters for thread safety
     - **Notable Features**: Event logging, statistics tracking, concurrent operation testing, propagation verification

   **9.5 Error Injection Framework** (TEST-INFRA-001-E) - **HIGH PRIORITY** ✅ **COMPLETED**
   - [x] **Create systematic error injection points** - Configurable error insertion at filesystem, Git, and network operations
   - [x] **Implement error type classification** - Categorize errors (transient, permanent, recoverable, fatal)
   - [x] **Add error propagation tracing** - Track error flow through operation chains
   - [x] **Create error recovery testing** - Test retry logic and graceful degradation
   - **Implementation Areas**: Error handling patterns in `errors.go`, Git operations in `git.go`, file operations across all modules
   - **Files Created**: `internal/testutil/errorinjection.go` (720 lines), `internal/testutil/errorinjection_test.go` (650 lines)
   - **Dependencies**: TEST-INFRA-001-B (disk space simulation), TEST-INFRA-001-C (permission testing) ✅ **SATISFIED**
   - **Design Decision**: Use interface-based injection with configurable error schedules rather than global state modification
   - **Implementation Notes**: 
     - **Interface-based design**: Created comprehensive interfaces for filesystem, Git, and network operations with injectable wrappers
     - **Error classification system**: Implemented 4 error types (Transient, Permanent, Recoverable, Fatal) and 6 categories (Filesystem, Git, Network, Permission, Resource, Configuration)
     - **Advanced injection control**: Support for trigger counts, max injections, probability-based injection, conditional functions, and timing delays
     - **Propagation tracing**: Complete error flow tracking through operation chains with timestamps and stack depth
     - **Recovery testing**: Built-in recovery attempt tracking with success/failure statistics and timing analysis
     - **Built-in scenarios**: Pre-configured scenarios for filesystem and Git error testing with expected results validation
     - **Performance**: Excellent performance with error injection at ~896ns/op and propagation tracking at ~214ns/op
     - **Thread-safe operations**: All operations protected with RWMutex for concurrent access
     - **Comprehensive testing**: 16 test functions plus 3 benchmarks achieving 100% coverage of core functionality

   **9.6 Integration Testing Orchestration** (TEST-INFRA-001-F) - **COMPLETED**
   - [x] **Create complex scenario composition** - Combine multiple error conditions for realistic testing
   - [x] **Implement test scenario scripting** - Declarative scenario definition for complex multi-step tests
   - [x] **Add timing and coordination utilities** - Synchronize multiple error conditions and operations
   - [x] **Create regression test suite integration** - Plug infrastructure into existing test suites
   - **Implementation Areas**: Integration with existing `*_test.go` files, comprehensive scenario testing
   - **Files Created**: `internal/testutil/scenarios.go` (1,100+ lines), `internal/testutil/scenarios_test.go` (900+ lines)
   - **Dependencies**: All previous TEST-INFRA-001 subtasks ✓
   - **Design Decision**: Use builder pattern for scenario composition with clear separation between setup, execution, and verification
   - **Implementation Notes**: 
     - **Builder Pattern**: Implemented fluent interface with `ScenarioBuilder` for intuitive scenario construction
     - **Phase-based Execution**: Clear separation into Setup, Execution, Verification, and Cleanup phases
     - **Parallel Step Support**: Configurable parallel execution within phases for performance testing
     - **Timing Coordination**: `TimingCoordinator` with barriers, signals, and delays for synchronized testing
     - **Runtime Environment**: Complete `ScenarioRuntime` with temp directories, config files, and shared data
     - **Error Integration**: Full integration with existing error injection, corruption, permission, and disk space components
     - **Comprehensive Testing**: 15 test functions covering builder patterns, orchestration, parallel execution, timing coordination
     - **Regression Integration**: Demonstrated integration patterns with existing test suite structure
     - **Performance Benchmarks**: Included benchmark tests for scenario execution performance

#### **Process Implementation (Next 2-4 weeks)**

### **Phase 3: Component Extraction and Generalization (HIGH PRIORITY)**

#### **Strategy Overview**
Extract components from the backup application to create reusable CLI building blocks that align with the generalized CLI specification. This extraction will create a foundation library that can be used by future Go CLI applications while maintaining the existing backup application.

#### **Refactoring Phases**

**🔧 PHASE 5A: Core Infrastructure Extraction (IMMEDIATE - Weeks 1-2)**

22. **Extract Configuration Management System** (EXTRACT-001)
    - [ ] **Create `pkg/config` package** - Extract configuration loading, merging, and validation
    - [ ] **Generalize configuration discovery** - Remove backup-specific paths and make configurable
    - [ ] **Extract environment variable support** - Generic env var override system
    - [ ] **Create configuration interfaces** - Define contracts for configuration providers
    - [ ] **Add configuration value tracking** - Generic source tracking for any config type
    - **Priority**: CRITICAL - Foundation for all other CLI apps
    - **Files to Extract**: `config.go` (1097 lines) → `pkg/config/`
    - **Design Decision**: Use interface-based design to support different config schema while keeping the robust discovery and merging logic
    - **Implementation Notes**: 
      - Maintain existing YAML support but make schema-agnostic
      - Extract search path logic, file merging, environment variable overrides
      - Create ConfigLoader interface that can be implemented for different schemas
      - Preserve source tracking and validation patterns

23. **Extract Error Handling and Resource Management** (EXTRACT-002)
    - [ ] **Create `pkg/errors` package** - Extract structured error types and handling
    - [ ] **Create `pkg/resources` package** - Extract ResourceManager and cleanup logic
    - [ ] **Generalize error context** - Remove backup-specific operation names
    - [ ] **Extract context-aware operations** - Generic cancellation and timeout support
    - [ ] **Create error classification utilities** - Disk space, permission, file system errors
    - **Priority**: CRITICAL - Foundation for reliable CLI operations
    - **Files to Extract**: `errors.go` (494 lines) → `pkg/errors/`, `pkg/resources/`
    - **Design Decision**: Separate error handling from resource management but maintain tight integration
    - **Implementation Notes**:
      - Keep ArchiveError pattern but generalize to ApplicationError
      - Extract ResourceManager as-is (it's already generic)
      - Preserve panic recovery and atomic operations
      - Maintain disk space and permission error classification

24. **Extract Output Formatting System** (EXTRACT-003)
    - [ ] **Create `pkg/formatter` package** - Extract printf and template formatting
    - [ ] **Generalize template engine** - Remove backup-specific template variables
    - [ ] **Extract regex pattern system** - Generic named pattern extraction
    - [ ] **Create output collector system** - Delayed output management
    - [ ] **Extract ANSI color support** - Terminal capability detection
    - **Priority**: HIGH - Critical for user experience consistency
    - **Files to Extract**: `formatter.go` (1675 lines) → `pkg/formatter/`
    - **Design Decision**: Create pluggable formatter interfaces with printf and template implementations
    - **Implementation Notes**:
      - Massive 1675-line file with rich functionality perfect for extraction
      - Template system with regex patterns is highly reusable
      - Output collector system (delayed output) is innovative and valuable
      - ANSI color support and formatting utilities are universally useful

25. **Extract Git Integration System** (EXTRACT-004)
    - [ ] **Create `pkg/git` package** - Extract Git repository detection and info extraction
    - [ ] **Generalize Git command execution** - Flexible Git operation framework
    - [ ] **Extract branch and hash utilities** - Reusable Git metadata extraction
    - [ ] **Create Git status detection** - Working directory state management
    - [ ] **Add Git configuration support** - Repository-specific configuration
    - **Priority**: MEDIUM - Valuable for many CLI apps but not universal
    - **Files to Extract**: `git.go` (122 lines) → `pkg/git/`
    - **Design Decision**: Keep command-line Git approach but make operations more flexible
    - **Implementation Notes**:
      - Small but complete Git integration suitable for many CLI apps
      - Command-line approach is simple and reliable
      - Status detection and repository validation are well-implemented
      - Could be extended with more Git operations in the future

**🔧 PHASE 5B: CLI Framework Extraction (Weeks 3-4)**

26. **Extract CLI Command Framework** (EXTRACT-005)
    - [ ] **Create `pkg/cli` package** - Extract cobra command patterns and flag handling
    - [ ] **Generalize command structure** - Template for common CLI patterns
    - [ ] **Extract dry-run implementation** - Generic dry-run operation support
    - [ ] **Create context-aware command execution** - Cancellation support for commands
    - [ ] **Extract version and build info handling** - Standard version command support
    - **Priority**: HIGH - Accelerates new CLI development
    - **Files to Extract**: `main.go` (816 lines) → `pkg/cli/`
    - **Design Decision**: Extract command patterns while leaving application-specific logic
    - **Implementation Notes**:
      - Rich command structure with backward compatibility patterns
      - Dry-run support is valuable pattern for many operations
      - Version handling with build-time information is common need
      - Context integration throughout command chain is well-implemented

27. **Extract File Operations and Utilities** (EXTRACT-006)
    - [ ] **Create `pkg/fileops` package** - Extract file comparison, copying, validation
    - [ ] **Generalize path validation** - Security and existence checking
    - [ ] **Extract atomic file operations** - Safe file writing patterns
    - [ ] **Create file exclusion system** - Generic pattern-based exclusion
    - [ ] **Extract directory traversal** - Safe directory walking with exclusions
    - **Priority**: HIGH - Common operations for many CLI apps
    - **Files to Extract**: `comparison.go` (329 lines), `exclude.go` (134 lines) → `pkg/fileops/`
    - **Design Decision**: Combine related file operations into cohesive package
    - **Implementation Notes**:
      - File comparison logic is robust and reusable
      - Path validation includes security considerations
      - Exclusion patterns using doublestar are valuable
      - Atomic operations integrate well with resource management

**🔧 PHASE 5C: Application-Specific Utilities (Weeks 5-6)**

28. **Extract Data Processing Patterns** (EXTRACT-007)
    - [ ] **Create `pkg/processing` package** - Extract data processing workflows
    - [ ] **Generalize naming conventions** - Timestamp-based naming patterns
    - [ ] **Extract verification systems** - Generic data integrity checking
    - [ ] **Create processing pipelines** - Template for data transformation workflows
    - [ ] **Extract concurrent processing** - Worker pool patterns
    - **Priority**: MEDIUM - Useful for data-focused CLI applications
    - **Files to Extract**: `archive.go` (679 lines), `backup.go` (869 lines), `verify.go` (442 lines) → `pkg/processing/`
    - **Design Decision**: Extract patterns while leaving domain-specific logic in original app
    - **Implementation Notes**:
      - Naming conventions with timestamps and metadata are broadly useful
      - Verification patterns can be adapted to different data types
      - Concurrent processing with context support is valuable
      - Pipeline pattern could accelerate development of data processing CLIs

29. **Create CLI Application Template** (EXTRACT-008)
    - [ ] **Create `cmd/cli-template` example - Working example using extracted packages
    - [ ] **Develop project scaffolding** - Generator for new CLI projects
    - [ ] **Create integration documentation** - How to use extracted components
    - [ ] **Add migration guide** - Moving from monolithic to package-based structure
    - [ ] **Create package interdependency mapping** - Clear usage patterns
    - **Priority**: HIGH - Demonstrates value and accelerates adoption
    - **Design Decision**: Create complete working example that showcases extracted components
    - **Implementation Notes**:
      - Template should demonstrate configuration, formatting, error handling, Git integration
      - Scaffolding could generate project structure with selected components
      - Documentation needs to show both individual package usage and integration patterns
      - Migration guide helps transition existing projects

**🔧 PHASE 5D: Testing and Documentation (Weeks 7-8)**

30. **Extract Testing Patterns and Utilities** (EXTRACT-009)
    - [ ] **Create `pkg/testutil` package** - Extract common testing utilities
    - [ ] **Generalize test fixtures** - Reusable test data management
    - [ ] **Extract test helpers** - Configuration, temporary files, assertions
    - [ ] **Create testing patterns documentation** - Best practices and examples
    - [ ] **Add package testing integration** - Ensure extracted packages are well-tested
    - **Priority**: HIGH - Critical for maintaining quality in extracted components
    - **Files to Extract**: Common patterns from `*_test.go` files → `pkg/testutil/`
    - **Design Decision**: Extract testing utilities while maintaining existing test coverage
    - **Implementation Notes**:
      - Rich testing patterns across 10 test files provide good extraction candidates
      - Temporary file/directory management patterns are valuable
      - Configuration testing patterns with multiple files and env vars are complex
      - Testing patterns ensure extracted packages are production-ready

31. **Create Package Documentation and Examples** (EXTRACT-010)
    - [ ] **Document each extracted package** - API documentation and usage examples
    - [ ] **Create integration examples** - How packages work together
    - [ ] **Add performance benchmarks** - Performance characteristics of extracted components
    - [ ] **Create troubleshooting guides** - Common issues and solutions
    - [ ] **Document design decisions** - Rationale for extraction choices
    - **Priority**: HIGH - Essential for adoption and maintenance
    - **Design Decision**: Comprehensive documentation following the existing high-quality documentation patterns
    - **Implementation Notes**:
      - Each package needs godoc-compatible documentation
      - Integration examples show real-world usage patterns
      - Performance documentation helps users understand trade-offs
      - Design decision documentation preserves knowledge

#### **Implementation Strategy and Design Decisions**

**📋 EXTRACTION PRINCIPLES:**

32. **Maintain Backward Compatibility** (DESIGN-001)
    - **Decision**: Extract components without breaking existing backup application
    - **Rationale**: Allows gradual migration and maintains stability of existing functionality
    - **Implementation**: Use Go modules and interfaces to isolate extracted code
    - **Design Note**: Original application continues to work while providing reusable components

33. **Interface-Driven Design** (DESIGN-002)
    - **Decision**: Create clear interfaces for all extracted components
    - **Rationale**: Enables flexibility, testing, and future evolution of implementations
    - **Implementation**: Define contracts before extracting concrete implementations
    - **Design Note**: Prevents tight coupling between extracted packages

34. **Zero-Breaking-Change Extraction** (DESIGN-003)
    - **Decision**: Extraction must not change existing application behavior
    - **Rationale**: Maintains confidence in extraction process and preserves existing functionality
    - **Implementation**: Comprehensive test coverage and behavioral verification
    - **Design Note**: Existing tests must continue to pass without modification

35. **Layered Extraction Approach** (DESIGN-004)
    - **Decision**: Extract in dependency order - core utilities first, then higher-level components
    - **Rationale**: Prevents circular dependencies and ensures stable foundation
    - **Implementation**: Infrastructure (config, errors) → Utilities (formatter, git) → Framework (cli) → Patterns (processing)
    - **Design Note**: Each layer builds on previous layers without circular references

**⚠️ EXTRACTION CHALLENGES AND SOLUTIONS:**

36. **Large File Decomposition** (CHALLENGE-001)
    - **Challenge**: `formatter.go` (1675 lines) is large and complex
    - **Solution**: Break into multiple focused packages while maintaining interface compatibility
    - **Approach**: Extract template engine, printf formatter, output collector, and ANSI support separately
    - **Design Note**: Size indicates rich functionality perfect for reuse, but needs careful decomposition

37. **Configuration Schema Flexibility** (CHALLENGE-002)
    - **Challenge**: Current config is backup-specific but needs to be generic
    - **Solution**: Extract configuration loading/merging logic with pluggable schema validation
    - **Approach**: Create ConfigLoader interface that can handle different struct types
    - **Design Note**: Preserve powerful discovery and merging while enabling different schemas

38. **Dependency Management** (CHALLENGE-003)
    - **Challenge**: Extracted packages will have interdependencies
    - **Solution**: Design clear dependency hierarchy and use Go modules for versioning
    - **Approach**: Core packages (config, errors) have no internal dependencies, higher-level packages compose core packages
    - **Design Note**: Clear layering prevents circular dependencies and simplifies usage

39. **Testing Complexity** (CHALLENGE-004)
    - **Challenge**: Extracted components need comprehensive testing without duplicating existing tests
    - **Solution**: Create package-specific tests while ensuring original integration tests still pass
    - **Approach**: Extract test utilities first, then create focused tests for each package
    - **Design Note**: Comprehensive testing ensures extracted components are production-ready

#### **Success Metrics for Extraction Project**

**📊 TECHNICAL METRICS:**
- **Code Reuse**: >80% of extracted code successfully reused in template application
- **Test Coverage**: >95% test coverage maintained for all extracted packages
- **Performance**: <5% performance degradation in original application
- **Interface Stability**: Zero breaking changes in extracted package interfaces after initial release

**🎯 ADOPTION METRICS:**
- **Template Usage**: Complete CLI template application built using extracted packages
- **Documentation Quality**: All packages have comprehensive godoc and usage examples
- **Migration Success**: Original backup application successfully migrated to use extracted packages
- **External Usage**: Framework suitable for building other CLI applications

**📈 QUALITY METRICS:**
- **Backward Compatibility**: All existing tests pass without modification
- **Package Independence**: Extracted packages can be used individually or in combination
- **Error Handling**: Comprehensive error propagation and context preservation
- **Resource Management**: Zero resource leaks in extracted components

#### **Timeline and Dependencies**

**Week 1-2 (Foundation)**: EXTRACT-001, EXTRACT-002 (config, errors, resources)
**Week 3-4 (Framework)**: EXTRACT-003, EXTRACT-004 (formatter, git) → EXTRACT-005 (cli framework)
**Week 5-6 (Patterns)**: EXTRACT-006, EXTRACT-007 (file ops, processing) → EXTRACT-008 (template)
**Week 7-8 (Quality)**: EXTRACT-009, EXTRACT-010 (testing, documentation)

**Critical Path**: Configuration and error handling must be completed before formatter and CLI framework extraction can begin. All core components must be ready before creating the template application.

**Risk Mitigation**: Each phase includes validation that existing application continues to work unchanged. Comprehensive testing ensures extraction doesn't introduce regressions.

This extraction project will create a powerful foundation for future Go CLI applications while preserving and enhancing the existing backup application. The extracted components embody years of refinement and testing, making them highly suitable for reuse.

## 🔧 VALIDATION AND AUTOMATION TOOLS

### 📋 PRE-COMMIT VALIDATION CHECKLIST

Before ANY commit that includes code changes, run this complete validation checklist:

```bash
# 1. MANDATORY: Check for feature tracking updates
echo "🔍 Checking feature tracking compliance..."
git diff --cached --name-only | grep -E '\.(go|yaml|yml)$' && echo "✅ Code changes detected - context files MUST be updated" || echo "ℹ️  No code changes"

# 2. MANDATORY: Verify context file updates
echo "📋 Verifying context documentation updates..."
git diff --cached --name-only | grep -E '^docs/context/' && echo "✅ Context files updated" || echo "❌ ERROR: Code changes without context updates"

# 3. MANDATORY: Check implementation tokens
echo "🏷️  Checking implementation tokens..."
git diff --cached | grep -E '^\+.*//.*[A-Z]+-[0-9]+:' && echo "✅ Implementation tokens found" || echo "⚠️  Warning: No new implementation tokens"

# 4. MANDATORY: Validate feature matrix
echo "📊 Validating feature matrix..."
grep -c "| [A-Z]+-[0-9]" docs/context/feature-tracking.md && echo "✅ Feature matrix has entries" || echo "❌ ERROR: Feature matrix empty"

# 5. MANDATORY: Check cross-references
echo "🔗 Checking cross-references..."
./scripts/validate-docs.sh || echo "❌ ERROR: Documentation validation failed"
```

### 🚨 AUTOMATED VALIDATION COMMANDS

#### **Quick Context Check**
```bash
# Check if context files need updates for current changes
function check_context_updates() {
    local code_changes=$(git diff --cached --name-only | grep -E '\.(go|yaml|yml)$' | wc -l)
    local context_changes=$(git diff --cached --name-only | grep -E '^docs/context/' | wc -l)
    
    if [ $code_changes -gt 0 ] && [ $context_changes -eq 0 ]; then
        echo "❌ ERROR: $code_changes code file(s) changed but no context files updated"
        echo "📋 Required context files to check:"
        echo "   - docs/context/feature-tracking.md (ALWAYS)"
        echo "   - docs/context/specification.md (if user-facing changes)"
        echo "   - docs/context/requirements.md (if requirements change)"
        echo "   - docs/context/architecture.md (if technical changes)"
        echo "   - docs/context/testing.md (if test changes)"
        return 1
    else
        echo "✅ Context file requirements satisfied"
        return 0
    fi
}
```

#### **Feature ID Validation**
```bash
# Validate all feature IDs are properly registered
function validate_feature_ids() {
    echo "🔍 Validating feature ID consistency..."
    
    # Extract feature IDs from code
    local code_features=$(grep -r "// [A-Z]+-[0-9]" --include="*.go" . | sed 's/.*\/\/ \([A-Z]*-[0-9]*\).*/\1/' | sort -u)
    
    # Extract feature IDs from tracking matrix
    local matrix_features=$(grep "| [A-Z]+-[0-9]" docs/context/feature-tracking.md | sed 's/|.*\([A-Z]*-[0-9]*\).*/\1/' | sort -u)
    
    # Find unregistered features
    local unregistered=$(comm -23 <(echo "$code_features") <(echo "$matrix_features"))
    
    if [ -n "$unregistered" ]; then
        echo "❌ ERROR: Unregistered feature IDs found in code:"
        echo "$unregistered"
        echo "💡 Add these features to docs/context/feature-tracking.md"
        return 1
    else
        echo "✅ All feature IDs properly registered"
        return 0
    fi
}
```

#### **Implementation Token Audit**
```bash
# Audit implementation tokens for completeness
function audit_implementation_tokens() {
    echo "🏷️  Auditing implementation tokens..."
    
    local missing_tokens=""
    
    # Check each .go file for feature-related functions without tokens
    for file in $(find . -name "*.go" -not -path "./vendor/*"); do
        # Look for functions that might need tokens but don't have them
        local suspicious_functions=$(grep -n "^func.*\(Archive\|Backup\|Config\|Git\|Format\|Test\)" "$file" | grep -v "// [A-Z]+-[0-9]")
        
        if [ -n "$suspicious_functions" ]; then
            missing_tokens="$missing_tokens\n$file:\n$suspicious_functions"
        fi
    done
    
    if [ -n "$missing_tokens" ]; then
        echo "⚠️  Functions potentially missing implementation tokens:"
        echo -e "$missing_tokens"
        echo "💡 Consider adding feature ID tokens to these functions"
    else
        echo "✅ Implementation token coverage looks good"
    fi
}
```

### 📄 DOCUMENTATION SYNC VALIDATION

#### **Cross-Reference Checker**
```bash
# Check all cross-references between context documents
function check_cross_references() {
    echo "🔗 Checking cross-references between context documents..."
    
    local broken_refs=""
    
    # Check references to feature IDs
    for doc in docs/context/*.md; do
        local doc_name=$(basename "$doc")
        local refs=$(grep -o "[A-Z]+-[0-9]\+" "$doc" | sort -u)
        
        for ref in $refs; do
            if ! grep -q "| $ref |" docs/context/feature-tracking.md; then
                broken_refs="$broken_refs\n$doc_name references undefined $ref"
            fi
        done
    done
    
    if [ -n "$broken_refs" ]; then
        echo "❌ ERROR: Broken cross-references found:"
        echo -e "$broken_refs"
        return 1
    else
        echo "✅ All cross-references valid"
        return 0
    fi
}
```

#### **Status Consistency Checker**
```bash
# Check feature status consistency across documents
function check_status_consistency() {
    echo "📊 Checking feature status consistency..."
    
    local inconsistencies=""
    
    # Extract features and their status from feature-tracking.md
    while IFS='|' read -r feature_id spec req arch test status token; do
        # Clean up the values
        feature_id=$(echo "$feature_id" | xargs)
        status=$(echo "$status" | xargs)
        
        if [[ "$feature_id" =~ ^[A-Z]+-[0-9]+$ ]]; then
            # Check if status matches in other documents
            # This is a simplified check - real implementation would be more thorough
            if [ "$status" = "Completed" ] || [ "$status" = "Implemented" ]; then
                if ! grep -q "$feature_id" docs/context/specification.md; then
                    inconsistencies="$inconsistencies\n$feature_id marked $status but not in specification.md"
                fi
            fi
        fi
    done < <(grep "| [A-Z]+-[0-9]" docs/context/feature-tracking.md)
    
    if [ -n "$inconsistencies" ]; then
        echo "⚠️  Status inconsistencies found:"
        echo -e "$inconsistencies"
    else
        echo "✅ Status consistency looks good"
    fi
}
```

### 🎯 DEVELOPER WORKFLOW INTEGRATION

#### **Git Hooks Setup**
```bash
# Add to .git/hooks/pre-commit
#!/bin/bash
echo "🚨 MANDATORY: Context Documentation Validation"

# Source the validation functions
source scripts/context-validation.sh

# Run all validation checks
check_context_updates || exit 1
validate_feature_ids || exit 1
check_cross_references || exit 1

echo "✅ All context documentation requirements satisfied"
```

#### **Make Target Integration**
```bash
# Add to Makefile
.PHONY: validate-context
validate-context:
	@echo "🔍 Validating context documentation..."
	@bash scripts/context-validation.sh
	@echo "✅ Context validation complete"

.PHONY: pre-commit
pre-commit: validate-context test
	@echo "✅ Pre-commit validation passed"

# Ensure context validation runs before tests
test: validate-context
	go test ./...
```

#### **IDE Integration Hints**
```bash
# VSCode settings.json snippet
{
    "go.buildTags": "integration",
    "files.watcherExclude": {
        "**/docs/context/**": false
    },
    "search.exclude": {
        "**/docs/context/**": false
    },
    "todo-tree.regex.regex": "((//|#|<!--|;|/\\*|^)\\s*($TAGS)|^\\s*- \\[ \\])",
    "todo-tree.regex.regexFlags": "gim",
    "todo-tree.highlights.customHighlight": {
        "FEATURE-ID": {
            "icon": "tag",
            "type": "tag",
            "foreground": "#FF6B6B"
        }
    }
}
```

### 📈 METRICS AND MONITORING

#### **Documentation Coverage Metrics**
```bash
# Generate context documentation coverage report
function generate_coverage_report() {
    echo "📊 Context Documentation Coverage Report"
    echo "========================================"
    
    local total_features=$(grep -c "| [A-Z]+-[0-9]" docs/context/feature-tracking.md)
    local implemented_features=$(grep -c "| [A-Z]+-[0-9].*Implemented\|Completed" docs/context/feature-tracking.md)
    local coverage_percentage=$((implemented_features * 100 / total_features))
    
    echo "Total Features: $total_features"
    echo "Implemented Features: $implemented_features"
    echo "Implementation Coverage: $coverage_percentage%"
    
    local code_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
    local files_with_tokens=$(grep -l "// [A-Z]+-[0-9]" --include="*.go" -r . | wc -l)
    local token_coverage=$((files_with_tokens * 100 / code_files))
    
    echo "Code Files: $code_files"
    echo "Files with Tokens: $files_with_tokens"
    echo "Token Coverage: $token_coverage%"
}
```

### 🔄 CONTINUOUS IMPROVEMENT

#### **Weekly Validation Audit**
```bash
# Run comprehensive weekly audit
function weekly_audit() {
    echo "📅 Weekly Context Documentation Audit"
    echo "====================================="
    
    check_context_updates
    validate_feature_ids
    audit_implementation_tokens
    check_cross_references
    check_status_consistency
    generate_coverage_report
    
    echo "📋 Action Items:"
    echo "- Review any warnings or errors above"
    echo "- Update missing implementation tokens"
    echo "- Verify feature status accuracy"
    echo "- Check for outdated documentation"
}
```

### 💡 AUTOMATION RECOMMENDATIONS

1. **CI/CD Integration**: Add context validation to GitHub Actions/Jenkins
2. **Documentation Generator**: Create tools to auto-update parts of context files
3. **Template Generation**: Auto-generate boilerplate for new features
4. **Dependency Tracking**: Monitor changes to shared components
5. **Compliance Dashboard**: Web interface showing documentation health

This comprehensive validation framework ensures that the context documentation remains synchronized with code changes and maintains the high quality standards established in this feature tracking system. 

### 📝 IMPLEMENTATION TOKEN REQUIREMENTS (DETAILED)

#### **Token Format Standards**
- **Standard Format**: `// FEATURE-ID: Brief description`
- **Examples**:
  ```go
  // ARCH-001: Archive naming convention implementation
  func GenerateArchiveName(cfg *Config, dir string) string {
      // ARCH-001: Include timestamp in archive name
      timestamp := time.Now().Format("2006-01-02-15-04")
      
      // GIT-002: Add Git branch and hash if available
      gitInfo := getGitInfo()
      if gitInfo.Branch != "" {
          return fmt.Sprintf("%s-%s=%s=%s.zip", dir, timestamp, gitInfo.Branch, gitInfo.Hash)
      }
      
      return fmt.Sprintf("%s-%s.zip", dir, timestamp)
  }
  
  // CFG-001: Configuration discovery implementation
  func GetConfigSearchPath() []string {
      // CFG-001: Check environment variable first
      if envPath := os.Getenv("BKPDIR_CONFIG"); envPath != "" {
          return strings.Split(envPath, ":")
      }
      
      // CFG-001: Use default search path
      return []string{"./.bkpdir.yml", "~/.bkpdir.yml"}
  }
  ```

#### **Token Placement Requirements**
1. **Function Level**: Every public function must have a feature token comment
2. **Method Level**: Every method must reference the relevant feature
3. **Code Block Level**: Significant logic blocks within functions should have tokens
4. **Error Handling**: Error paths should include tokens for traceability
5. **Configuration Handling**: All config reading/writing must have tokens

#### **AI ASSISTANT TOKEN HANDLING INSTRUCTIONS**

**🤖 FOR AI ASSISTANTS: AUTOMATIC TOKEN DETECTION AND REFERENCING**

When making code changes, AI assistants MUST:

1. **🔍 AUTO-DETECT AFFECTED TOKENS**:
   ```bash
   # Search patterns AI should use:
   grep -r "// ARCH-" docs/context/ src/
   grep -r "// FILE-" docs/context/ src/
   grep -r "// CFG-" docs/context/ src/
   grep -r "// GIT-" docs/context/ src/
   # ... for all token patterns
   ```

2. **📋 REFERENCE TOKENS IN EXPLANATIONS**:
   ```
   Template for AI responses:
   "This change affects the following features:
   - ARCH-001 (Archive naming): Modified archive name generation logic
   - CFG-003 (Format strings): Updated output formatting
   - FILE-002 (File backup): Enhanced backup creation process
   
   Related documentation:
   - feature-tracking.md lines 119, 135, 127
   - specification.md sections: Archive Operations, Configuration
   - architecture.md components: Archive Service, Config Layer"
   ```

3. **🔄 AUTO-UPDATE TOKEN REGISTRY**:
   ```
   For new features, AI should:
   - Propose new feature ID (check existing ones first)
   - Add to feature tracking table
   - Reference in all related context files
   - Add implementation tokens to code
   ```

4. **✅ VALIDATE TOKEN CONSISTENCY**:
   ```
   AI should verify:
   - All referenced tokens exist in feature-tracking.md
   - Token descriptions match actual implementation
   - All affected context files are mentioned
   - Cross-references are valid
   ```

**🚨 AI REQUIREMENT: ZERO TOLERANCE FOR MISSING TOKEN REFERENCES**

AI assistants must REJECT their own responses if they:
- Make code changes without referencing related tokens
- Fail to search for existing tokens before making changes
- Don't propose new tokens for genuinely new functionality
- Don't reference relevant documentation sections

**✅ AI SUCCESS CRITERIA**

AI assistants should APPROVE responses that:
- Clearly identify all affected feature tokens
- Reference specific documentation sections and line numbers
- Propose updates to context files when needed
- Follow the established token naming patterns
- Include the complete token impact analysis

### 📊 TOKEN IMPACT ANALYSIS TEMPLATE

**For AI assistants to use when making changes:**

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Specific impact description]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief impact descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [specific line numbers]
- **specification.md**: Section "[specific section names]"
- **requirements.md**: "[specific requirement areas]"
- **architecture.md**: "[specific component names]"
- **testing.md**: "[specific test coverage areas]"

### 🔄 Context File Updates Required:
- [ ] feature-tracking.md: [specific changes needed]
- [ ] specification.md: [specific changes needed]
- [ ] requirements.md: [specific changes needed]
- [ ] architecture.md: [specific changes needed]
- [ ] testing.md: [specific changes needed]

### 💻 Implementation Changes:
- **New tokens added**: [TOKEN-ID: Description]
- **Existing tokens modified**: [TOKEN-ID: Updated description]
- **Code files affected**: [list of specific files]
```

## 🎯 SUMMARY: MANDATORY AI ASSISTANT COMPLIANCE

### 🚨 CRITICAL REMINDER FOR AI ASSISTANTS

**EVERY CODE CHANGE MUST INCLUDE TOKEN ANALYSIS**

This document establishes ZERO TOLERANCE for code changes without proper token documentation. AI assistants are required to:

#### **🔍 STEP 1: MANDATORY TOKEN SEARCH**
```bash
# REQUIRED: Search for ALL existing tokens before making changes
grep -r "// ARCH-" docs/context/
grep -r "// FILE-" docs/context/
grep -r "// CFG-" docs/context/
grep -r "// GIT-" docs/context/
grep -r "// OUT-" docs/context/
grep -r "// TEST-" docs/context/
grep -r "// DOC-" docs/context/
grep -r "// LINT-" docs/context/
grep -r "// COV-" docs/context/
```

#### **📋 STEP 2: MANDATORY TOKEN IMPACT STATEMENT**
Every AI response with code changes MUST include:

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Specific impact description]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief impact descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [specific line numbers]
- **specification.md**: Section "[specific section names]"
- **requirements.md**: "[specific requirement areas]"
- **architecture.md**: "[specific component names]"
- **testing.md**: "[specific test coverage areas]"

### 🔄 Context File Updates Required:
- [ ] feature-tracking.md: [specific changes needed]
- [ ] specification.md: [specific changes needed]
- [ ] requirements.md: [specific changes needed]
- [ ] architecture.md: [specific changes needed]
- [ ] testing.md: [specific changes needed]

### 💻 Implementation Changes:
- **New tokens added**: [TOKEN-ID: Description]
- **Existing tokens modified**: [TOKEN-ID: Updated description]
- **Code files affected**: [list of specific files]
```

#### **🚫 STEP 3: MANDATORY REJECTION CRITERIA**

AI assistants MUST REFUSE to provide code changes if:
- No token search was performed
- No existing tokens are referenced for related functionality
- No new tokens are proposed for new functionality
- Context file updates are not specified
- Documentation cross-references are missing

#### **✅ STEP 4: MANDATORY APPROVAL CRITERIA**

AI assistants should ONLY proceed if their response:
- Includes complete token impact analysis using the required template
- References specific line numbers in feature-tracking.md
- Links to relevant sections in all context files
- Proposes concrete updates to affected documentation
- Follows established token naming conventions (ARCH-, FILE-, CFG-, etc.)

### 🔄 ENFORCEMENT SUMMARY

**REMEMBER: THIS IS NOT OPTIONAL**

The feature-tracking system is designed to prevent unplanned changes and maintain documentation synchronization. AI assistants that fail to follow these requirements undermine the entire system integrity.

**Key Success Metrics:**
- 100% of code changes must reference related tokens
- 100% of new features must have new tokens
- 100% of responses must include context file update plans
- 0% tolerance for undocumented changes

**Implementation Tokens Are Mandatory:**
Every function, method, and significant code block must include `// TOKEN-ID: Description` comments that trace back to this feature tracking matrix.

### 📞 QUICK REFERENCE FOR AI ASSISTANTS

1. **Search Pattern**: `grep -r "// [A-Z]+-[0-9]+" docs/context/`
2. **Token Format**: `// FEATURE-ID: Brief description`
3. **Required Template**: Use the Token Impact Analysis template above
4. **Context Files**: feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md
5. **Zero Tolerance**: No code changes without token references

**This document serves as the authoritative source for feature tracking and token management. Any AI assistant working with this codebase must comply with these requirements.**

### **🔧 PHASE 4: PRE-EXTRACTION REFACTORING (CRITICAL FOUNDATION)**

**Priority Summary:**
- **Code structure preparation**: Essential before extraction to ensure clean package boundaries
- **Interface stabilization**: Must be completed before extraction to prevent breaking changes
- **Dependency cleanup**: Critical for avoiding circular dependencies in extracted packages
- **Architecture preparation**: Required for successful component extraction

#### **IMMEDIATE PRE-EXTRACTION TASKS (Week 0 - Before Extraction Begins)**

**10. Dependency Analysis and Interface Standardization** (REFACTOR-001) - **CRITICAL BLOCKER**
   - [x] **Complete dependency mapping analysis** - Map all current dependencies between major components
   - [x] **Identify circular dependency risks** - Find potential circular imports before extraction
   - [x] **Standardize interface contracts** - Define clear interfaces for all major components
   - [x] **Create interface compatibility layer** - Ensure backward compatibility during extraction
   - [x] **Validate package boundary design** - Confirm clean separation between future packages
   - **Rationale**: Must identify and resolve dependency issues before extraction to prevent circular imports and ensure clean package boundaries
   - **Status**: ✅ **COMPLETED**
   - **Priority**: CRITICAL - BLOCKS ALL EXTRACTION WORK ⚠️
   - **Blocking**: All EXTRACT-001 through EXTRACT-010 tasks
   - **Implementation Areas**:
     - ✅ Complete dependency graph analysis across all core files
     - ✅ Interface definition for Config, OutputFormatter, ResourceManager, Git integration
     - ✅ Cross-cutting concern identification (logging, error handling, context management)
     - ✅ Package boundary validation using Go module dependency analysis
   - **Dependencies**: None (foundational pre-extraction work)
   - **Implementation Tokens**: `// REFACTOR-001: Dependency analysis`, `// REFACTOR-001: Interface standardization`
   - **Expected Outcomes**:
     - ✅ Clear dependency map showing extraction order requirements
     - ✅ Standardized interfaces preventing tight coupling
     - ✅ Validated package boundaries ensuring clean extraction
     - ✅ Zero circular dependency risks identified
   - **Deliverables**:
     - ✅ `docs/extraction-dependencies.md` - Complete dependency analysis
     - ✅ `docs/interface-definitions.md` - Interface definitions for all major components
     - ✅ Package boundary validation report
     - ✅ Circular dependency risk assessment
   - **Implementation Notes**:
     - **Dependency Analysis**: Comprehensive analysis revealed clean extraction boundaries for config, git, and comparison components
     - **Interface Contracts**: Defined 25+ interfaces covering all major component interactions
     - **Package Boundaries**: Validated 8 distinct components with clear separation of concerns
     - **Zero Circular Dependencies**: Confirmed unidirectional dependency flow enables safe extraction
     - **Implementation Tokens**: Added to all 8 core Go files marking interface standardization points
     - **Extraction Order**: Established 4-phase extraction plan based on dependency complexity
     - **Critical Finding**: Config.go has clean boundary (no internal dependencies) - ready for immediate extraction

**11. Large File Decomposition Preparation** (REFACTOR-002) - **HIGH PRIORITY**
   - [x] **Analyze formatter.go structure** (1677 lines) - Break down large file for extraction readiness
   - [x] **Identify component boundaries in formatter.go** - Separate template engine, printf formatter, output collector
   - [x] **Create internal interfaces within large files** - Prepare for clean extraction
   - [x] **Validate decomposition strategy** - Ensure each component can be extracted independently
   - [x] **Prepare extraction interfaces** - Define contracts for extracted formatter components
   - **Rationale**: The 1677-line formatter.go needs decomposition analysis before extraction to ensure proper package boundaries
   - **Status**: ✅ **COMPLETED** (2025-01-02)
   - **Priority**: HIGH - Required for EXTRACT-003 (Output Formatting System)
   - **Blocking**: EXTRACT-003 (Output Formatting System)
   - **Implementation Areas**:
     - OutputFormatter component analysis and interface definition ✅
     - TemplateFormatter separation and interface design ✅
     - OutputCollector isolation and contract definition ✅
     - ANSI color support component identification ✅ (None found - no ANSI color code in current implementation)
     - Pattern extraction engine component boundaries ✅
   - **Dependencies**: REFACTOR-001 (dependency analysis must be completed first) ✅
   - **Implementation Tokens**: `// REFACTOR-002: Formatter decomposition`, `// REFACTOR-002: Component boundary` ✅
   - **Expected Outcomes**:
     - Clear component boundaries within formatter.go ✅
     - Interface contracts for each formatter component ✅
     - Validated extraction strategy for large file ✅
     - Reduced complexity through logical separation ✅
   - **Deliverables**:
     - `docs/formatter-decomposition.md` - Component analysis and extraction plan ✅
     - Interface definitions for formatter components ✅
     - Extraction strategy document ✅
   - **Implementation Notes**:
     - **5 Component Boundaries Identified**: OutputCollector (ready for immediate extraction), PrintfFormatter, TemplateFormatter, PatternExtractor, ErrorFormatter
     - **Config Dependency Challenge**: All components except OutputCollector require config interface abstraction due to tight coupling with Config struct
     - **Clean Extraction Path**: OutputCollector has zero dependencies and can be extracted immediately; other components require FormatProvider interface
     - **No ANSI Color Support**: Analysis revealed no ANSI color handling in current formatter implementation
     - **Interface Design**: Created internal interfaces (FormatProvider, OutputDestination, PatternExtractor, FormatterInterface, TemplateFormatterInterface) to prepare for extraction
     - **Backward Compatibility**: Extraction strategy includes wrapper pattern to preserve existing method signatures during transition

**12. Configuration Schema Abstraction** (REFACTOR-003) - **HIGH PRIORITY**
   - [ ] **Create configuration loader interface** - Abstract configuration loading from specific schema
   - [ ] **Separate configuration logic from backup-specific schema** - Enable schema-agnostic configuration
   - [ ] **Design pluggable configuration validation** - Allow different applications to define their own schemas
   - [ ] **Create configuration source abstraction** - Abstract file, environment, and default sources
   - [ ] **Prepare configuration merging interfaces** - Enable generic configuration merging logic
   - **Rationale**: Current configuration is tightly coupled to backup application schema; must be abstracted for reuse
   - **Status**: Not Started
   - **Priority**: HIGH - Required for EXTRACT-001 (Configuration Management System)
   - **Blocking**: EXTRACT-001 (Configuration Management System)
   - **Implementation Areas**:
     - ConfigLoader interface definition
     - ConfigValidator interface for pluggable validation
     - ConfigSource interface for different configuration sources
     - ConfigMerger interface for generic merging logic
     - Schema abstraction layer for different application types
   - **Dependencies**: REFACTOR-001 (dependency analysis must identify config coupling)
   - **Implementation Tokens**: `// REFACTOR-003: Config abstraction`, `// REFACTOR-003: Schema separation`
   - **Expected Outcomes**:
     - Schema-agnostic configuration loading
     - Pluggable validation system
     - Reusable configuration merging logic
     - Source-independent configuration management
   - **Deliverables**:
     - Configuration interface definitions
     - Schema abstraction design document
     - Configuration extraction plan

**13. Error Handling and Resource Management Consolidation** (REFACTOR-004) - **MEDIUM PRIORITY**
   - [ ] **Standardize error type patterns** - Ensure consistent error handling across components
   - [ ] **Consolidate resource management patterns** - Standardize ResourceManager usage
   - [ ] **Create context propagation standards** - Ensure consistent context handling
   - [ ] **Validate atomic operation patterns** - Confirm consistent atomic file operations
   - [ ] **Prepare error handling for extraction** - Design extractable error handling patterns
   - **Rationale**: Error handling and resource management must be consistent before extraction to ensure reliable extracted components
   - **Status**: Not Started
   - **Priority**: MEDIUM - Required for EXTRACT-002 (Error Handling and Resource Management)
   - **Blocking**: EXTRACT-002 (Error Handling and Resource Management)
   - **Implementation Areas**:
     - Error type standardization across ArchiveError, BackupError patterns
     - ResourceManager usage pattern validation
     - Context propagation consistency checking
     - Atomic operation pattern validation
     - Panic recovery standardization
   - **Dependencies**: REFACTOR-001 (dependency analysis must identify error handling patterns)
   - **Implementation Tokens**: `// REFACTOR-004: Error standardization`, `// REFACTOR-004: Resource consolidation`
   - **Expected Outcomes**:
     - Consistent error handling patterns
     - Standardized resource management
     - Reliable context propagation
     - Uniform atomic operations
   - **Deliverables**:
     - Error handling standardization report
     - Resource management pattern documentation
     - Context propagation guidelines

**14. Code Structure Optimization for Extraction** (REFACTOR-005) - **MEDIUM PRIORITY**
   - [ ] **Remove tight coupling between components** - Identify and resolve unnecessary dependencies
   - [ ] **Standardize naming conventions** - Ensure consistent naming across extractable components
   - [ ] **Optimize import structure** - Prepare for clean package imports after extraction
   - [ ] **Validate function signatures for extraction** - Ensure extractable functions have clean signatures
   - [ ] **Prepare backward compatibility layer** - Plan compatibility preservation during extraction
   - **Rationale**: Code structure must be optimized for clean extraction without breaking existing functionality
   - **Status**: Not Started
   - **Priority**: MEDIUM - Enhances extraction quality but not blocking
   - **Implementation Areas**:
     - Component coupling analysis and reduction
     - Naming convention standardization across codebase
     - Import optimization for future package structure
     - Function signature validation for extractability
     - Backward compatibility planning
   - **Dependencies**: REFACTOR-001, REFACTOR-002, REFACTOR-003 (prior refactoring must be completed)
   - **Implementation Tokens**: `// REFACTOR-005: Structure optimization`, `// REFACTOR-005: Extraction preparation`
   - **Expected Outcomes**:
     - Reduced coupling between components
     - Consistent naming conventions
     - Optimized import structure
     - Clean function signatures
     - Preserved backward compatibility
   - **Deliverables**:
     - Code structure optimization report
     - Naming convention guidelines
     - Extraction compatibility assessment

#### **REFACTORING VALIDATION AND QUALITY ASSURANCE (Week 0.5)**

**15. Refactoring Impact Validation** (REFACTOR-006) - **HIGH PRIORITY**
   - [ ] **Run comprehensive test suite after each refactoring** - Ensure no functionality regression
   - [ ] **Validate performance impact** - Confirm refactoring doesn't degrade performance
   - [ ] **Check implementation token consistency** - Verify all tokens remain valid after refactoring
   - [ ] **Validate documentation synchronization** - Ensure context files reflect refactoring changes
   - [ ] **Run extraction readiness assessment** - Confirm codebase is ready for component extraction
   - **Rationale**: All refactoring must be validated to ensure it improves extraction readiness without breaking functionality
   - **Status**: Not Started
   - **Priority**: HIGH - Must validate each refactoring step
   - **Implementation Areas**:
     - Automated test suite execution after each refactoring
     - Performance benchmarking and comparison
     - Implementation token validation and updating
     - Documentation synchronization checking
     - Extraction readiness criteria validation
   - **Dependencies**: All REFACTOR-001 through REFACTOR-005 tasks
   - **Implementation Tokens**: `// REFACTOR-006: Validation`, `// REFACTOR-006: Quality assurance`
   - **Expected Outcomes**:
     - Zero functional regressions from refactoring
     - Maintained or improved performance
     - Consistent implementation tokens
     - Synchronized documentation
     - Validated extraction readiness
   - **Deliverables**:
     - Refactoring validation report
     - Performance impact assessment
     - Extraction readiness certification

### **🎯 REFACTORING SUCCESS CRITERIA**

#### **MANDATORY PRE-EXTRACTION REQUIREMENTS**
Before proceeding with any EXTRACT-001 through EXTRACT-010 tasks, these criteria must be met:

1. **✅ Dependency Analysis Complete** (REFACTOR-001)
   - Complete dependency map created
   - Zero circular dependency risks identified
   - Clean package boundaries validated
   - Interface contracts defined

2. **✅ Large File Analysis Complete** (REFACTOR-002)  
   - Formatter.go decomposition strategy defined
   - Component boundaries identified
   - Extraction interfaces prepared

3. **✅ Configuration Abstraction Ready** (REFACTOR-003)
   - Schema-agnostic configuration interfaces defined
   - Pluggable validation system designed
   - Configuration extraction plan created

4. **✅ Error Handling Standardized** (REFACTOR-004)
   - Consistent error patterns across codebase
   - Standardized resource management
   - Uniform context propagation

5. **✅ Refactoring Validated** (REFACTOR-006)
   - All tests pass after refactoring
   - Performance maintained or improved
   - Documentation synchronized

#### **REFACTORING QUALITY GATES**

**Technical Quality Gates:**
- Zero test failures after each refactoring step
- <5% performance impact from refactoring changes
- 100% implementation token consistency maintained
- All context documentation updated to reflect changes

**Extraction Readiness Gates:**
- Clear component interfaces defined
- Zero circular dependency risks
- Clean package boundary validation
- Backward compatibility preservation plan

#### **REFACTORING ENFORCEMENT**

**Pre-Extraction Validation Checklist:**
```bash
# MANDATORY: Run before any extraction work
echo "🔧 Pre-Extraction Refactoring Validation"

# 1. Verify dependency analysis completion
[ -f "docs/extraction-dependencies.md" ] && echo "✅ Dependency analysis complete" || echo "❌ BLOCKER: Dependency analysis missing"

# 2. Verify large file decomposition analysis  
[ -f "docs/formatter-decomposition.md" ] && echo "✅ Formatter decomposition complete" || echo "❌ BLOCKER: Formatter analysis missing"

# 3. Verify configuration abstraction
grep -q "ConfigLoader interface" *.go && echo "✅ Config abstraction ready" || echo "❌ BLOCKER: Config abstraction missing"

# 4. Run all tests
go test ./... && echo "✅ All tests pass" || echo "❌ BLOCKER: Test failures"

# 5. Check implementation tokens
./scripts/validate-docs.sh && echo "✅ Documentation consistent" || echo "❌ WARNING: Documentation sync needed"
```

**Extraction Authorization:**
Component extraction is ONLY authorized after:
- All REFACTOR-001 through REFACTOR-005 tasks completed
- REFACTOR-006 validation passed
- Pre-extraction validation checklist passed
- Zero critical blockers remaining

### **📋 REFACTORING TASK INTEGRATION**

#### **Context File Update Requirements**
Each refactoring task MUST update:
- **feature-tracking.md**: Task status and implementation tokens
- **architecture.md**: Component interface changes and design decisions  
- **requirements.md**: Any requirement impacts from refactoring
- **testing.md**: Test coverage for refactored components

#### **AI Assistant Integration**
AI assistants working on refactoring MUST:
- Reference existing REFACTOR-XXX tokens when making changes
- Update dependency documentation when modifying component relationships
- Validate extraction readiness after each refactoring step
- Ensure backward compatibility preservation

#### **Implementation Decision Tracking**
Add new decision records for refactoring:

#### DEC-013: Pre-Extraction Refactoring Strategy
- **Feature IDs**: REFACTOR-001, REFACTOR-002, REFACTOR-003, REFACTOR-004, REFACTOR-005
- **Date**: 2024-12-19
- **Decision**: Implement comprehensive pre-extraction refactoring to ensure clean component boundaries
- **Rationale**:
  - Prevents circular dependencies in extracted packages
  - Ensures clean interfaces and minimal coupling
  - Preserves backward compatibility during extraction
  - Enables reliable, maintainable extracted components
- **Alternatives**: Extract components as-is and refactor later, minimal refactoring approach
- **Impact**: Establishes foundation for successful component extraction with clean architecture
- **Code Markers**: `// REFACTOR-XXX: Preparation for extraction` tokens

#### DEC-014: Interface-First Extraction Approach
- **Feature IDs**: REFACTOR-001, REFACTOR-003
- **Date**: 2024-12-19
- **Decision**: Define all component interfaces before extraction begins
- **Rationale**:
  - Prevents tight coupling between extracted packages
  - Enables independent testing and development of components
  - Provides clear contracts for component interaction
  - Facilitates future evolution of implementations
- **Alternatives**: Extract implementations first and define interfaces later
- **Impact**: All extracted packages will have well-defined, stable interfaces
- **Code Markers**: Interface definitions with `// REFACTOR-001: Interface standardization`

#### DEC-015: Large File Decomposition Strategy
- **Feature IDs**: REFACTOR-002
- **Date**: 2024-12-19  
- **Decision**: Decompose large files (>1000 lines) before extraction for better package boundaries
- **Rationale**:
  - 1675-line formatter.go contains multiple logical components
  - Clean separation enables focused extracted packages
  - Reduces complexity and improves maintainability
  - Enables independent evolution of formatter components
- **Alternatives**: Extract large files as single packages, post-extraction decomposition
- **Impact**: Extracted packages will be focused and maintainable rather than monolithic
- **Code Markers**: `// REFACTOR-002: Component boundary` markings for logical separations