# Semantic Tokens Directory

**TIED Methodology Version**: 2.2.0 (align with AGENTS.md)

## Overview
This document is a **human-readable companion** to the machine registry. **Authoritative** fields (rationale, satisfaction criteria, `essence_pseudocode`, traceability) live in **detail YAML**, not duplicated here.

| Layer | Index | Per-token detail |
|-------|--------|------------------|
| REQ | `tied/requirements.yaml` | `tied/requirements/REQ-*.yaml` |
| ARCH | `tied/architecture-decisions.yaml` | `tied/architecture-decisions/ARCH-*.yaml` |
| IMPL | `tied/implementation-decisions.yaml` | `tied/implementation-decisions/IMPL-*.yaml` |
| Registry | `tied/semantic-tokens.yaml` | (merged view with methodology) |

Narrative digests (secondary): `tied/requirements.md`, `tied/architecture-decisions.md`, `tied/implementation-decisions.md`. Agent rules: `AGENTS.md`, `ai-principles.md`.

## AI Assistant Integration Guidelines [REQ-DOC_016]

### Token Usage for AI Assistants

AI assistants should use semantic tokens for:

1. **Code Navigation**: Search for `[REQ-*]`, `[ARCH-*]`, `[IMPL-*]` tokens to find related code
2. **Feature Understanding**: Trace features from requirements through architecture to implementation
3. **Change Impact Analysis**: Use token cross-references to identify affected components
4. **Test Discovery**: Find tests for features using `[REQ-*]` tokens in test names

### Token-Based Code Navigation

```bash
# Find all implementations of a requirement
grep -r "\[REQ-FEATURE_NAME\]" --include="*.go" .

# Find all tests for a requirement
grep -r "REQ_FEATURE_NAME" --include="*_test.go" .

# Find architecture decisions for a feature
grep -r "\[ARCH-FEATURE_NAME\]" --include="*.md" .

# Find implementation details
grep -r "\[IMPL-FEATURE_NAME\]" --include="*.go" .
```

### Token Creation Requirements

When implementing features:
1. **ALWAYS** add or extend `[REQ-*]` detail under `tied/requirements/` and `tied/requirements.yaml` (TIED MCP preferred).
2. **ALWAYS** add `[ARCH-*]` detail under `tied/architecture-decisions/` and `tied/architecture-decisions.yaml`.
3. **ALWAYS** add `[IMPL-*]` detail (including `essence_pseudocode` with token comments) under `tied/implementation-decisions/`.
4. **ALWAYS** annotate code and tests with `[REQ-*]` / `[ARCH-*]` / `[IMPL-*]` per `AGENTS.md`.
5. **ALWAYS** update `tied/semantic-tokens.yaml` when introducing new tokens.
6. **ALWAYS** follow PROC tokens via `AGENTS.md` and `tied/docs/` / `tied/methodology/`.

### Token Audit Workflow `[PROC-TOKEN_AUDIT]`

- Map requirement → architecture → implementation tokens before touching code.
- Annotate every code edit with `[IMPL-*] [ARCH-*] [REQ-*]` (same triplet used in documentation).
- Require tests to include the `[REQ-*]` (and optional `[TEST-*]`) identifiers in both the test name and supporting comments.
- Record the audit result inside the relevant task/subtask so future agents can see when the chain was verified.

### Automated Validation `[PROC-TOKEN_VALIDATION]`

- Run `./scripts/validate_tokens.sh` (or repo-specific equivalent) after each audit to ensure every referenced token exists in the registry.
- Treat validation failures as blocking defects until the registry and documents are synchronized.
- Capture validation outputs in task notes or IMPL detail metadata as appropriate; registry truth is `tied/semantic-tokens.yaml` plus detail YAML.

### Token Validation Requirements

Before marking features complete:
1. **ALWAYS** run token validation scripts (e.g., `./scripts/validate_tokens.sh`) and retain evidence per project practice (see `README.md`).
2. **ALWAYS** ensure token consistency across all layers
3. **ALWAYS** verify token traceability in documentation
4. **ALWAYS** check that all cross-references are valid

## Token Format

```
[TYPE-IDENTIFIER]
```

## Inherited tokens (TIED/LEAP methodology)

All TIED projects **inherit** a core set of REQ/ARCH/IMPL and PROC tokens via `copy_files.sh` (from `templates/`). These tokens are **mandatory for TIED success** and enforce the methodology; they must not be removed from the client's `tied/`. The inherited set includes:

- **REQ**: REQ-TIED_SETUP, REQ-MODULE_VALIDATION; optionally REQ-FEEDBACK_TO_TIED
- **ARCH**: ARCH-TIED_STRUCTURE, ARCH-MODULE_VALIDATION; optionally ARCH-FEEDBACK_STORAGE
- **IMPL**: IMPL-TIED_FILES, IMPL-MODULE_VALIDATION; optionally IMPL-MCP_FEEDBACK_TOOLS
- **PROC**: e.g. PROC-LEAP, PROC-TOKEN_AUDIT, PROC-TOKEN_VALIDATION, PROC-TIED_DEV_CYCLE, PROC-IMPL_PSEUDOCODE_TOKENS (see `processes.md` and `semantic-tokens.yaml`)

For structure and sample records, agents refer to **`templates/`** in the TIED repository (see AGENTS.md § Client inheritance of LEAP R+A+I).

## Token Types

- `[REQ-*]` - Requirements (functional/non-functional) - **The source of intent**
- `[ARCH-*]` - Architecture decisions - **High-level design choices that preserve intent**
- `[IMPL-*]` - Implementation decisions - **Low-level choices that preserve intent**
- `[TEST-*]` - Test specifications - **Validation of intent**
- `[PROC-*]` - Process definitions for survey/build/test/deploy work that stay linked to `[REQ-*]`

## Token Naming Convention

- Use UPPER_SNAKE_CASE for identifiers
- Be descriptive but concise
- Example: `[REQ-DUPLICATE_PREVENTION]` not `[REQ-DP]`

## Cross-Reference Format

When referencing other tokens:

```markdown
[IMPL-EXAMPLE] Description [ARCH-DESIGN] [REQ-REQUIREMENT]
```

## Token Registry

**📖 The canonical registry of all tokens lives in `semantic-tokens.yaml`.**

This YAML index serves as the single source of truth for "does this token exist?" and provides structured metadata for all semantic tokens across all types (REQ, ARCH, IMPL, TEST, PROC).

**Quick lookup commands:**
- List all tokens: `yq 'keys' tied/semantic-tokens.yaml`
- Filter by type: `yq 'to_entries | map(select(.value.type == "REQ")) | from_entries' tied/semantic-tokens.yaml`
- Check existence: `yq '.["REQ-TIED_SETUP"]' tied/semantic-tokens.yaml`
- Get token details: `yq '.REQ-TIED_SETUP' tied/semantic-tokens.yaml`

**For full details on each token:**
- **Requirements tokens**: See `requirements.yaml` (YAML index) and `requirements/` (detail files)
- **Architecture tokens**: See `architecture-decisions.yaml` (YAML index) and `architecture-decisions/` (detail files)
- **Implementation tokens**: See `implementation-decisions.yaml` (YAML index) and `implementation-decisions/` (detail files)
- **Process tokens**: See `processes.md`

## Token Relationships

### Hierarchical Relationships
- `[REQ-PARENT_FEATURE]` contains `[REQ-SUB_FEATURE_1]`, `[REQ-SUB_FEATURE_2]`
- `[ARCH-FEATURE]` includes `[ARCH-COMPONENT_1]`, `[ARCH-COMPONENT_2]`

### Flow Relationships
- `[REQ-FEATURE]` → `[ARCH-DESIGN]` → `[IMPL-IMPLEMENTATION]` → Code + Tests

### Dependency Relationships
- `[IMPL-FEATURE]` depends on `[ARCH-DESIGN]` and `[REQ-FEATURE]`
- `[ARCH-DESIGN]` depends on `[REQ-FEATURE]`

## Usage Examples

### In Code Comments
```[your-language]
// [REQ-EXAMPLE_FEATURE] Implementation of example feature
// [IMPL-EXAMPLE_IMPLEMENTATION] [ARCH-EXAMPLE_DECISION] [REQ-EXAMPLE_FEATURE]
function exampleFunction() {
    // ...
}
```
> **NOTE**: Code merged without these annotations is considered incomplete because it fails `[PROC-TOKEN_AUDIT]`.

### In Tests
```[your-language]
// Test validates [REQ-EXAMPLE_FEATURE] is met
function testExampleFeature_REQ_EXAMPLE_FEATURE() {
    // Test implementation
}
```
> **NOTE**: Tests without `[REQ-*]` markers are rejected during `[PROC-TOKEN_VALIDATION]` because they cannot prove intent.

### In Documentation
```markdown
The feature uses [ARCH-ARCHITECTURE_NAME] to fulfill [REQ-FEATURE_NAME].
Implementation details are documented in [IMPL-IMPLEMENTATION_NAME].
```

## Token Validation Guidelines

### Cross-Layer Token Consistency

Every feature must have proper token coverage across all layers:

1. **Requirements Layer**: Feature must have `[REQ-*]` detail in `tied/requirements/` (indexed in `tied/requirements.yaml`)
2. **Architecture Layer**: Architecture decisions must have `[ARCH-*]` detail in `tied/architecture-decisions/`
3. **Implementation Layer**: Implementation must have `[IMPL-*]` detail in `tied/implementation-decisions/` and tokens in code comments
4. **Test Layer**: Tests must reference `[REQ-*]` tokens in test names/comments
5. **Documentation Layer**: All documentation must cross-reference tokens consistently

### Token Format Validation

1. **Token Format**: Must follow `[TYPE-IDENTIFIER]` pattern exactly
2. **Token Types**: Must use valid types (`REQ`, `ARCH`, `IMPL`, `TEST`, `PROC`)
3. **Identifier Format**: Must use UPPER_SNAKE_CASE
4. **Cross-References**: Implementation tokens must reference architecture and requirement tokens

### Token Traceability Validation

1. Every requirement detail YAML must chain to corresponding `[ARCH-*]` / `[IMPL-*]` where applicable
2. Every architecture decision detail YAML must chain to implementation and tests as documented in traceability fields
3. Every test must link to specific requirements via `[REQ-*]` tokens
4. All tokens must be discoverable through automated validation

## Token Creation Requirements

When implementing features:
1. **ALWAYS** add or extend `[REQ-*]` detail under `tied/requirements/` and `tied/requirements.yaml` (TIED MCP preferred).
2. **ALWAYS** add `[ARCH-*]` detail under `tied/architecture-decisions/` and `tied/architecture-decisions.yaml`.
3. **ALWAYS** add `[IMPL-*]` detail (including `essence_pseudocode` with token comments) under `tied/implementation-decisions/`.
4. **ALWAYS** annotate code and tests with `[REQ-*]` / `[ARCH-*]` / `[IMPL-*]` per `AGENTS.md`.
5. **ALWAYS** update `tied/semantic-tokens.yaml` when introducing new tokens.
6. **ALWAYS** follow PROC tokens via `AGENTS.md` and `tied/docs/` / `tied/methodology/`.

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

## REQ, ARCH, and IMPL registries (no per-token duplicate list)

Long per-token bullet lists that pointed at `requirements.md` / `architecture-decisions.md` / `implementation-decisions.md` sections duplicated the **detail YAML** and drifted from TIED.

- **REQ**: `tied/requirements.yaml` and `tied/requirements/REQ-*.yaml`
- **ARCH**: `tied/architecture-decisions.yaml` and `tied/architecture-decisions/ARCH-*.yaml`
- **IMPL**: `tied/implementation-decisions.yaml` and `tied/implementation-decisions/IMPL-*.yaml`

Use MCP or your editor to open detail files by token name. Human digests remain in `tied/requirements.md`, `tied/architecture-decisions.md`, and `tied/implementation-decisions.md` for narrative only.

**Governance**: Project TIED currently includes `[REQ-GOV_REGISTRY_COMPLETENESS]` and `[REQ-GOV_DISCOVERABILITY]` under `tied/requirements/`; do not treat informal `REQ-GOV_*` lists elsewhere as authoritative.

## Legacy feature and action aliases (`project-tokens`)

Historical **features**, **actions**, and semantic metadata migrated from `project-tokens.yaml` are enumerated in **`tied/semantic-tokens.yaml`**. Keeping a second full matrix here duplicated the registry and pointed at removed or moved docs (`docs/context/*`, old analysis files).

- **Canonical list**: `tied/semantic-tokens.yaml` (with methodology merge as documented in `AGENTS.md`).
- **Behavior and R+A+I**: `tied/requirements/`, `tied/architecture-decisions/`, `tied/implementation-decisions/`.
- **Validation**: `scripts/validate-token-traceability.sh`, `./scripts/validate_tokens.sh`, `scripts/token-coverage-analysis.sh` as applicable.

### EXTRACT-008 (still tracked here)

- `[REQ:EXTRACT_008_INTERDEP_MAPPING]` — Package interdependency mapping; deliverable: `docs/package-interdependency-mapping.md`
- `[ARCH:EXTRACT_008_INTERDEP]` — Inter-package dependency mapping architecture
- `[IMPL:EXTRACT_008_DOC_MIGRATION]` — Documentation migration for EXTRACT-008 working plans

## Quick Reference Index

### New Implementation Tokens
- `[IMPL:LIST_FORMAT_SAFETY]` → `tied/implementation-decisions/IMPL-LIST_FORMAT_SAFETY.yaml` ([ARCH:OUTPUT_FORMATTING] [REQ:OUT_002])



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
- **Visualization**: `STDD_VIS_FLOW`

### Implementation Tokens
- **Core**: `CONFIG_STRUCT`, `CFG_006`, `ZIP_FORMAT`, `DUAL_FORMATTING`, `STRUCTURED_ERRORS`, `GIT_CLI`, `RESOURCE_MANAGER`, `CONTEXT_OPS`, `ATOMIC_OPS`, `TESTING`, `CODE_STYLE`, `DATA_MODELS`
- **File Operations**: `FILE_001`, `FILE_002`, `FILE_003`
- **Configuration**: `CFG_001`, `CFG_002`, `CFG_003`, `CFG_004`, `CFG_005`, `CFG_TEMPLATE_001`, `CONFIGURABLE_STRINGS`, `CFG_PRECEDENCE_FIX`, `CFG_HIERARCHY_PRESERVATION`, `CFG_MERGE_PREPEND_PRECEDENCE_FIX`
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
- **Visualization**: `STDD_VIS_DATA_PIPELINE`, `STDD_VIS_ASSETS`

### Test Tokens
- `[TEST:EXCLUDE_MERGE]` → `config_integration_test.go::TestExcludePatternsMerge_REQ_TEST_EXCLUDE_MERGE` validates `[REQ:TEST_EXCLUDE_MERGE]`, `[ARCH:TEST_EXCLUDE_MERGE]`, `[IMPL:TEST_EXCLUDE_MERGE]`
- `[TEST:LIST_LIMIT]` → `list_limit_test.go::TestListLimitDefault` validates `[REQ:LIST_LIMIT]`, `[ARCH:LIST_LIMIT]`, `[IMPL:LIST_LIMIT]`

### New Implementation Tokens
- `[IMPL:LIST_FORMAT_SAFETY]` → `tied/implementation-decisions/IMPL-LIST_FORMAT_SAFETY.yaml` ([ARCH:OUTPUT_FORMATTING] [REQ:OUT_002])

