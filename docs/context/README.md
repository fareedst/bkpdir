# AI Assistant Context Documentation Index

> **⚠️ UPDATED**: This project follows **Semantic Token-Driven Development (STDD)** methodology. This index provides navigation to STDD documents and historical context files.

## Current Development Process

This project follows **Semantic Token-Driven Development (STDD)**. Please refer to:

1. **`.cursorrules`** - Complete STDD methodology and rules (auto-loaded by Cursor IDE)
2. **`ai-principles.md`** - Complete AI-First Principles and process guide
3. **STDD Documents** - See below for primary STDD documentation

## STDD Documents (Primary References)

### Core STDD Documents

1. **`stdd/ai-principles.md`** - Complete STDD methodology and process guide
2. **`stdd/requirements.md`** - All functional and non-functional requirements with `[REQ:*]` tokens
3. **`stdd/architecture-decisions.md`** - Architecture decisions with `[ARCH:*]` tokens
4. **`stdd/implementation-decisions.md`** - Implementation decisions with `[IMPL:*]` tokens
5. **`stdd/tasks.md`** - Active tasks and feature tracking
6. **`stdd/semantic-tokens.md`** - Central registry of all semantic tokens

### User-Facing Documentation

- **`../user/specification.md`** - User-facing features and behaviors (moved to `docs/user/`)
- **`testing.md`** - Test coverage requirements and validation standards (migrated to implementation-decisions.md)

### Historical Context Files

All historical context files have been merged into STDD (`ai-principles.md`) and removed. The STDD methodology now includes:

- Complete cross-referencing guidance with feature documentation format
- Change impact tracking with pre/post validation checklists
- Behavioral contracts and dependency mapping
- Code standards and documentation templates
- All validation and traceability guidance

## STDD Development Process

### Phase 1: Requirements → Pseudo-Code (NO CODE YET)

**MANDATORY**: Before any code changes:

1. Read `stdd/ai-principles.md` - Understand the complete STDD process
2. Check `stdd/semantic-tokens.md` - Review existing tokens
3. Review `stdd/requirements.md` - Check existing requirements
4. Review `stdd/architecture-decisions.md` - Check existing architecture decisions
5. Review `stdd/implementation-decisions.md` - Check existing implementation decisions
6. Expand requirements into pseudo-code and decisions
7. Document architecture decisions IMMEDIATELY in `stdd/architecture-decisions.md` with `[ARCH:*]` tokens
8. Document implementation decisions IMMEDIATELY in `stdd/implementation-decisions.md` with `[IMPL:*]` tokens
9. Update `stdd/semantic-tokens.md` with new tokens
10. Create tasks in `stdd/tasks.md` with priorities BEFORE implementation

### Phase 2: Tasks → Implementation

1. Work on highest priority tasks first (P0 > P1 > P2 > P3)
2. Update documentation AS YOU WORK
3. Mark tasks complete when all criteria met

## Documentation Update Requirements

For ANY code changes, you MUST:

1. **Reference Requirements**: Link code changes to `[REQ:*]` tokens in `stdd/requirements.md`
2. **Reference Architecture**: Link code changes to `[ARCH:*]` tokens in `stdd/architecture-decisions.md`
3. **Reference Implementation**: Link code changes to `[IMPL:*]` tokens in `stdd/implementation-decisions.md`
4. **Update Task Tracking**: Update `stdd/tasks.md` with task status
5. **Update Semantic Tokens**: Update `stdd/semantic-tokens.md` if creating new tokens

### Code Comment Format

```go
// [REQ-FILE_BACKUP] Create backup of single file with comparison
// [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
    // ...
}
```

### Test Name Format

```go
func TestCreateFileBackup_REQ_FILE_BACKUP(t *testing.T) {
    // ...
}
```

## Token Search Quick Commands

```bash
# Search for semantic tokens
grep -r "\[REQ:" . --include="*.go"   # Requirements tokens
grep -r "\[ARCH:" . --include="*.go"   # Architecture tokens
grep -r "\[IMPL:" . --include="*.go"  # Implementation tokens

# Find tokens in documentation
grep -r "\[REQ:" *.md                  # Requirements in docs
grep -r "\[ARCH:" *.md                 # Architecture in docs
grep -r "\[IMPL:" *.md                 # Implementation in docs

# Validate all changes
make test && make lint
```

## AI Assistant Validation Checklist

Before submitting any code changes, verify:

- [ ] Read `stdd/ai-principles.md` and understand STDD process
- [ ] Check `stdd/semantic-tokens.md` for existing tokens
- [ ] Review `stdd/requirements.md` for related requirements
- [ ] Review `stdd/architecture-decisions.md` for related architecture
- [ ] Review `stdd/implementation-decisions.md` for related implementation
- [ ] Added semantic token references (`[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`) to code
- [ ] Updated `stdd/tasks.md` with task status
- [ ] Updated `stdd/semantic-tokens.md` if creating new tokens
- [ ] All tests pass (`make test`)
- [ ] All lint checks pass (`make lint`)

## Critical Reminders

1. **STDD Process**: Follow the 3-phase STDD development process (Requirements → Tasks → Implementation)
2. **Documentation First**: Document requirements, architecture, and implementation decisions BEFORE coding
3. **Semantic Tokens**: Use `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]` tokens for all cross-references
4. **Task Planning**: Create tasks in `tasks.md` with priorities BEFORE implementation
5. **Validation**: Run tests and linting before marking tasks complete

## Quick Reference

If you're unsure about anything:

1. **Read `stdd/ai-principles.md`** - Complete STDD methodology guide
2. **Check `stdd/semantic-tokens.md`** - Token registry
3. **Review `stdd/tasks.md`** - Active tasks
4. **Review STDD documents** - Requirements, architecture, implementation decisions
5. **Validate Changes**: Run `make test && make lint`

---

**Last Updated**: 2025-01-13  
**Status**: Updated to reflect STDD methodology migration 