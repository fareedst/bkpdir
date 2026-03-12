# AI Assistant Compliance Guide

**STDD Methodology Version**: 1.0.2  
**Related**: [REQ:DOC_015] (Unicode to Semantic Token Mapping), [REQ:DOC_016] (AI-First Comprehensive Token System)

## Overview

This document provides compliance guidelines for AI assistants working on this project. All AI assistants must follow these guidelines to ensure consistency, traceability, and maintainability.

## Semantic Token Requirements [REQ:DOC_015]

### Mandatory Token Usage

AI assistants **MUST** use semantic tokens in all code, documentation, and comments:

1. **Requirements Tokens** `[REQ:IDENTIFIER]`
   - Use when referencing or implementing requirements
   - Example: `// [REQ:FILE_BACKUP] Create backup of single file`
   - Must reference tokens defined in `tied/requirements.md`

2. **Architecture Tokens** `[ARCH:IDENTIFIER]`
   - Use when referencing architecture decisions
   - Example: `// [ARCH:RESOURCE_MANAGEMENT] Use atomic operations`
   - Must reference tokens defined in `tied/architecture-decisions.md`
   - Must cross-reference `[REQ:*]` tokens

3. **Implementation Tokens** `[IMPL:IDENTIFIER]`
   - Use when documenting implementation details
   - Example: `// [IMPL:ATOMIC_OPS] [ARCH:RESOURCE_MANAGEMENT] [REQ:FILE_BACKUP]`
   - Must cross-reference both `[ARCH:*]` and `[REQ:*]` tokens

### Token Format Rules

- **Format**: `[TYPE:IDENTIFIER]` where TYPE is REQ, ARCH, or IMPL
- **Identifier**: Must use UPPER_SNAKE_CASE (e.g., `FILE_BACKUP`, `RESOURCE_MANAGEMENT`)
- **Cross-References**: Implementation tokens must reference architecture and requirement tokens
- **Consistency**: Use the same identifier across all layers (requirements → architecture → implementation)

### Token Validation

All semantic tokens are validated automatically:
- Format validation: Ensures tokens follow `[TYPE:IDENTIFIER]` pattern
- Consistency validation: Ensures cross-references are valid
- Registry validation: Ensures tokens are documented in `tied/semantic-tokens.md`

## Code Comment Requirements

### Required Format

All code comments that reference features, decisions, or requirements **MUST** include semantic tokens:

```go
// [REQ:FILE_BACKUP] Create backup of single file with comparison
// [IMPL:ATOMIC_OPS] [ARCH:RESOURCE_MANAGEMENT] [REQ:FILE_BACKUP]
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
    // Implementation
}
```

### Test Function Requirements

All test functions **MUST** reference the requirement they validate:

```go
// Test validates [REQ:FILE_BACKUP] is met
func TestCreateFileBackup_REQ_FILE_BACKUP(t *testing.T) {
    // Test implementation
}
```

## Documentation Requirements

### Requirements Documentation

- All requirements **MUST** have `[REQ:IDENTIFIER]` tokens
- Requirements **MUST** include satisfaction and validation criteria
- Requirements **MUST** be documented in `tied/requirements.md`

### Architecture Documentation

- All architecture decisions **MUST** have `[ARCH:IDENTIFIER]` tokens
- Architecture decisions **MUST** cross-reference `[REQ:*]` tokens
- Architecture decisions **MUST** be documented in `tied/architecture-decisions.md`

### Implementation Documentation

- All implementation decisions **MUST** have `[IMPL:IDENTIFIER]` tokens
- Implementation decisions **MUST** cross-reference both `[ARCH:*]` and `[REQ:*]` tokens
- Implementation decisions **MUST** be documented in `tied/implementation-decisions.md`

## Task Management Requirements

### Task Documentation

- All tasks **MUST** reference semantic tokens
- Tasks **MUST** be documented in `tied/tasks.md`
- Tasks **MUST** have priorities (P0/P1/P2/P3)
- Subtasks **MUST** be removed when parent task completes

### Task Format

```markdown
## P0: Task Name [REQ:IDENTIFIER] [ARCH:IDENTIFIER] [IMPL:IDENTIFIER]

**Status**: 🟡 In Progress | ✅ Complete | ⏸️ Blocked | ⏳ Pending

**Description**: Brief description referencing semantic tokens.

**Dependencies**: List of other tasks/tokens this depends on.
```

## Validation Requirements

### Pre-Submission Validation

Before submitting changes, AI assistants **MUST**:

1. Run semantic token validation: `make validate-tokens`
2. Ensure all tokens are properly formatted
3. Ensure all cross-references are valid
4. Ensure all tokens are documented in `tied/semantic-tokens.md`

### Validation Commands

- `make validate-tokens` - Validate semantic token usage
- `make validate-tokens-strict` - Strict validation for CI/CD
- `./scripts/validate-semantic-tokens.sh` - Direct script execution

## AI Assistant Workflow

### Standard Workflow

1. **Acknowledge Principles**: Start every response with "Observing AI principles!"
2. **Read Documentation**: Review `tied/ai-principles.md` and `tied/semantic-tokens.md`
3. **Use Semantic Tokens**: Include tokens in all code, comments, and documentation
4. **Cross-Reference**: Link tokens across layers (REQ → ARCH → IMPL)
5. **Validate**: Run validation before completing work

### Code Navigation

AI assistants should use semantic tokens for:
- **Feature Discovery**: Search for `[REQ:FEATURE_NAME]` to find all related code
- **Architecture Understanding**: Search for `[ARCH:DESIGN_NAME]` to understand design decisions
- **Implementation Details**: Search for `[IMPL:IMPLEMENTATION_NAME]` to find specific implementations
- **Test Discovery**: Search for `REQ_FEATURE_NAME` in test names to find related tests

### Token-Based Search Examples

**Using Token Navigation Tool [REQ:DOC_016]:**

```bash
# Find all implementations of a requirement
./scripts/token-navigate.sh find-req FEATURE_NAME

# Find all tests for a requirement
./scripts/token-navigate.sh find-tests FEATURE_NAME

# Find architecture decisions for a feature
./scripts/token-navigate.sh find-arch FEATURE_NAME

# Find implementation details
./scripts/token-navigate.sh find-impl FEATURE_NAME

# Trace a token across all layers (REQ → ARCH → IMPL → Code → Tests)
./scripts/token-navigate.sh trace FEATURE_NAME

# Show token coverage across all layers
./scripts/token-navigate.sh coverage FEATURE_NAME

# List all tokens
./scripts/token-navigate.sh list-req
./scripts/token-navigate.sh list-arch
./scripts/token-navigate.sh list-impl
```

**Using grep (alternative):**

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

## Compliance Checklist

Before marking work complete, verify:

- [ ] All code comments include semantic tokens
- [ ] All test functions reference requirements
- [ ] All requirements documented with `[REQ:*]` tokens
- [ ] All architecture decisions documented with `[ARCH:*]` tokens
- [ ] All implementation decisions documented with `[IMPL:*]` tokens
- [ ] All tokens cross-referenced correctly
- [ ] All tokens documented in `tied/semantic-tokens.md`
- [ ] Validation passes: `make validate-tokens`
- [ ] No unicode icons in documentation headers
- [ ] All validation systems recognize semantic tokens

## Token Navigation Tools [REQ:DOC_016]

### Token Navigation Script

The project includes a comprehensive token navigation tool at `scripts/token-navigate.sh`:

**Features:**
- Find implementations by token type (REQ/ARCH/IMPL)
- Trace tokens across all layers
- Show token coverage statistics
- List all tokens in the codebase

**Usage:**
```bash
./scripts/token-navigate.sh <command> [options]
```

See "Token-Based Search Examples" section above for detailed usage.

### Integration with AI Assistants

AI assistants should use the token navigation tool for:
- **Feature Discovery**: `./scripts/token-navigate.sh find-req FEATURE_NAME`
- **Impact Analysis**: `./scripts/token-navigate.sh trace FEATURE_NAME`
- **Coverage Validation**: `./scripts/token-navigate.sh coverage FEATURE_NAME`
- **Token Discovery**: `./scripts/token-navigate.sh list-req`

## Related Documentation

- `tied/ai-principles.md` - Complete AI-first principles and process guide
- `tied/semantic-tokens.md` - Central registry of all semantic tokens
- `tied/requirements.md` - Requirements with `[REQ:*]` tokens
- `tied/architecture-decisions.md` - Architecture decisions with `[ARCH:*]` tokens
- `tied/implementation-decisions.md` - Implementation decisions with `[IMPL:*]` tokens
- `tied/tasks.md` - Active task tracking
- `scripts/token-navigate.sh` - Token navigation tool [REQ:DOC_016]

## Version History

- **1.0.0** (2025-01-XX): Initial version with semantic token requirements
- **1.0.1** (2025-01-XX): Added validation requirements and compliance checklist
- **1.0.2** (2025-01-XX): Extended for DOC-015 completion with STDD token integration
