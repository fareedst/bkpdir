# Semantic Token System Requirements

> **DOC-014: Core development requirement - Semantic token system as foundation of AI-first development [ACTION:core-functionality]**

## Overview

The semantic token system is a **core requirement** for all development work on this project. It replaces the legacy Unicode icon system with a machine-readable, AI-optimized approach that ensures consistent, maintainable, and searchable code annotations.

## Core Requirements

### 1. Mandatory Token Format

All implementation work MUST use the standardized semantic token format:

```go
// FEATURE-ID: Description [ACTION:context]
```

**Example:**
```go
// ARCH-001: Archive naming convention [ACTION:core-functionality]
```

### 2. Token Validation

All code changes MUST pass semantic token validation:

```bash
make validate-token-enforcement
```

This validation is automatically run as part of:
- `make check` - Full code quality checks
- `make lint` - Linting process
- CI/CD pipeline (strict mode)

### 3. Registry Compliance

All tokens MUST be defined in the central registry:

**File:** `project-tokens.yaml`

The registry defines:
- Valid feature IDs: `ARCH-001`, `CFG-003`, `GIT-004`, `DOC-014`
- Valid actions: `core-functionality`, `format-processing`, `discovery`, `maintenance`, `validation`

### 4. Migration from Legacy Icons

The system has migrated from Unicode icons to semantic tokens:

```bash
# Legacy (deprecated)
// [CRITICAL] ARCH-001: Archive naming convention - [ACTION:core-functionality] Core functionality

# Semantic (current)
// ARCH-001: Archive naming convention [ACTION:core-functionality]
```

## Development Workflow

### 1. Token Creation

When implementing new features:

1. **Check registry** - Verify feature ID exists in `project-tokens.yaml`
2. **Use correct format** - `FEATURE-ID: Description [ACTION:context]`
3. **Validate immediately** - Run `make validate-token-enforcement`

### 2. Token Maintenance

- **Registry updates** - Add new feature IDs to `project-tokens.yaml`
- **Consistent usage** - Use same feature ID for related work
- **Migration planning** - Convert legacy icons systematically

### 3. Quality Gates

The semantic token system provides multiple quality gates:

- **Development** - Fast validation during coding
- **Pre-commit** - Validation before commits
- **CI/CD** - Strict validation in pipeline
- **Documentation** - Consistency across all files

## AI-First Benefits

### 1. Perfect Searchability

```bash
# Find all archive-related work
grep "ARCH-001" *.go

# Find core functionality
grep "ACTION:core-functionality" *.go

# Find configuration features
grep "CFG-" *.go
```

### 2. Semantic Understanding

AI assistants can:
- Parse priority levels semantically
- Understand feature relationships
- Maintain consistency automatically
- Generate accurate suggestions

### 3. Cross-Platform Reliability

- **No Unicode issues** - Pure ASCII text
- **No font dependencies** - Works everywhere
- **No encoding problems** - Universal compatibility
- **No rendering issues** - Terminal-safe

## Enforcement Mechanisms

### 1. Automated Validation

```bash
# Standard validation
make validate-token-enforcement

# Strict mode (CI/CD)
make validate-tokens-strict

# Development mode
make validate-tokens
```

### 2. Build Integration

Token validation is integrated into:
- `make check` - Full quality checks
- `make lint` - Linting process
- All CI/CD pipelines
- Pre-commit hooks (recommended)

### 3. Documentation Requirements

All documentation must:
- Use semantic tokens instead of Unicode icons
- Reference the central registry
- Maintain consistency with code tokens
- Pass validation checks

## Migration Strategy

### Phase 1: System Implementation ✅
- [x] Replace icon validation with semantic validation
- [x] Update Makefile integration
- [x] Create central registry
- [x] Implement validation script

### Phase 2: Legacy Migration (In Progress)
- [ ] Convert 215 legacy Unicode icons
- [ ] Update all documentation files
- [ ] Migrate all Go source files
- [ ] Update shell scripts

### Phase 3: Enforcement
- [ ] Enable strict validation in CI/CD
- [ ] Add pre-commit hooks
- [ ] Create migration tools
- [ ] Update developer documentation

## Registry Management

### Adding New Feature IDs

```yaml
# In project-tokens.yaml
features:
  NEW-001:
    priority: HIGH
    description: "New feature description"
    action: core-functionality
    status: active
```

### Adding New Actions

```yaml
# In project-tokens.yaml
actions:
  new-action:
    description: "New action description"
    search_pattern: "\\[ACTION:new-action\\]"
```

## Performance Benefits

- **30x faster validation** - <1 second vs 25+ seconds
- **Reliable search** - 100% accuracy vs encoding issues
- **Simple maintenance** - Single registry vs multiple files
- **AI-optimized** - Semantic parsing vs visual interpretation

## Compliance Requirements

### For Developers
- **Mandatory usage** - All new code must use semantic tokens
- **Validation required** - Must pass `make validate-token-enforcement`
- **Registry compliance** - Use only registered feature IDs and actions
- **Documentation updates** - Keep tokens consistent across docs and code

### For AI Assistants
- **Token system first** - Always use semantic tokens
- **Registry awareness** - Check valid tokens in `project-tokens.yaml`
- **Validation integration** - Run validation after changes
- **Consistency maintenance** - Keep tokens aligned across all files

## Success Metrics

- **Error rate**: 0 validation errors required
- **Coverage**: 100% of new code uses semantic tokens
- **Performance**: Sub-second validation time
- **Maintainability**: Single source of truth in registry

## Conclusion

The semantic token system is not optional - it is a **core requirement** for maintaining code quality, AI-first development, and long-term project success. All development work must comply with these requirements to ensure consistent, maintainable, and AI-optimized codebase. 