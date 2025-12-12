# Merge Strategy Conflict Resolution Proposal

**Date**: 2025-12-11  
**Status**: Proposal  
**Related**: [ARCH:CFG_005], [ARCH:CFG_001], [REQ:CFG_005], [REQ:CONFIGURATION]

## Problem Statement

There is a conflict between two requirements:

1. **CFG-005**: Array fields (like `exclude_patterns`) should default to "merge" (accumulate) strategy in all contexts - both inheritance chains and sequential file processing.

2. **CFG-001**: Earlier files take precedence over later files when processing sequential config files.

**The Conflict:**
- For `exclude_patterns`, CFG-005 says to merge/accumulate by default
- But CFG-001 says earlier files take precedence, which would mean replace, not merge
- Currently, `applyReplace` respects earlier file precedence even with `!` prefix, which conflicts with the merge behavior

**Current Behavior:**
- `exclude_patterns` defaults to "merge" strategy (CFG-005)
- But `applyReplace` respects earlier file precedence even with `!` prefix (CFG-001)
- This creates inconsistent behavior where explicit `!exclude_patterns` in a later file doesn't work as expected

## Proposed Solution

### Option 1: Field-Level Merge Behavior Registry (Recommended)

Create a registry that specifies merge behavior per field, allowing different fields to have different precedence rules.

**Implementation:**
```go
// Field merge behavior configuration
type FieldMergeBehavior int

const (
    // MergeBehaviorAccumulate: Field accumulates values from all files (CFG-005 style)
    MergeBehaviorAccumulate FieldMergeBehavior = iota
    // MergeBehaviorPrecedence: Field respects earlier file precedence (CFG-001 style)
    MergeBehaviorPrecedence
)

// Field merge behavior registry
var fieldMergeBehaviors = map[string]FieldMergeBehavior{
    "exclude_patterns": MergeBehaviorAccumulate,  // CFG-005: Merge by default
    "archive_dir_path": MergeBehaviorPrecedence,   // CFG-001: Earlier files win
    "include_git_info": MergeBehaviorPrecedence,   // CFG-001: Earlier files win
    // ... other fields
}
```

**Behavior:**
- Fields with `MergeBehaviorAccumulate`: Always merge, even in sequential file processing
- Fields with `MergeBehaviorPrecedence`: Respect earlier file precedence, even for array fields
- Explicit prefixes (`!`, `+`, `^`, `=`) override the default behavior

**Pros:**
- Clear, explicit configuration per field
- Easy to understand and maintain
- Allows different fields to have different behaviors
- Backward compatible (can default to current behavior)

**Cons:**
- Requires maintaining a registry
- Need to update registry when adding new fields

### Option 2: Semantic Token-Based Configuration

Use semantic tokens to mark fields with their merge behavior.

**Implementation:**
```go
// Field metadata with semantic tokens
type FieldMetadata struct {
    Name           string
    MergeBehavior  FieldMergeBehavior
    SemanticTokens []string  // e.g., [REQ:CFG_005], [REQ:CFG_001]
}

var fieldMetadata = map[string]FieldMetadata{
    "exclude_patterns": {
        Name:          "exclude_patterns",
        MergeBehavior: MergeBehaviorAccumulate,
        SemanticTokens: []string{"[REQ:CFG_005]"},
    },
    "archive_dir_path": {
        Name:          "archive_dir_path",
        MergeBehavior: MergeBehaviorPrecedence,
        SemanticTokens: []string{"[REQ:CFG_001]"},
    },
}
```

**Pros:**
- Links merge behavior to requirements via semantic tokens
- Traceable to requirements
- Self-documenting

**Cons:**
- More complex than Option 1
- Requires semantic token system integration

### Option 3: Type-Based Default with Field Overrides

Default behavior based on field type (array vs. scalar), with field-level overrides.

**Implementation:**
```go
// Default: Arrays accumulate, scalars respect precedence
// Override: Specific fields can override default

var fieldMergeOverrides = map[string]FieldMergeBehavior{
    // Only specify fields that need to override default
    // "exclude_patterns" uses default (accumulate for arrays)
    // "archive_dir_path" uses default (precedence for scalars)
}
```

**Pros:**
- Minimal configuration needed
- Intuitive defaults

**Cons:**
- Less explicit
- Harder to understand edge cases

## Recommended Approach

**Option 1 (Field-Level Registry)** is recommended because:
1. It's explicit and clear
2. Easy to maintain and understand
3. Allows fine-grained control per field
4. Backward compatible
5. Can be extended with semantic tokens later (Option 2)

## Implementation Plan

1. **Create field merge behavior registry** in `config.go`
2. **Update `applyMerge` and `applyReplace`** to check registry
3. **Update `applyMergeStrategies`** to respect field-specific behavior
4. **Document field behaviors** in architecture decisions
5. **Update tests** to reflect new behavior
6. **Update STDD documentation** with new decision

## Field Behavior Matrix

| Field | Type | Default Behavior | Sequential Files | Inheritance Chains |
|-------|------|------------------|------------------|-------------------|
| `exclude_patterns` | Array | Accumulate (CFG-005) | Merge | Merge |
| `archive_dir_path` | Scalar | Precedence (CFG-001) | Earlier wins | Child overrides |
| `include_git_info` | Scalar | Precedence (CFG-001) | Earlier wins | Child overrides |
| `!exclude_patterns` | Array | Replace | Earlier wins* | Replace |

*Note: Even with `!` prefix, earlier file precedence applies for sequential files per CFG-001, but this conflicts with CFG-005. The registry would resolve this by allowing field-specific behavior.

## Resolution Strategy

For `exclude_patterns` specifically:
- **Default (no prefix)**: Always merge/accumulate (CFG-005)
- **Explicit `!` prefix**: Replace, but still respect earlier file precedence if field was set by earlier file (CFG-001)
- **Explicit `+` prefix**: Explicit merge (same as default for this field)
- **Explicit `^` prefix**: Prepend to existing values
- **Explicit `=` prefix**: Use only if not set by earlier file

This allows both requirements to coexist:
- CFG-005: Default behavior is merge (accumulate)
- CFG-001: Earlier files take precedence (even with `!` prefix, if earlier file set it)

## Next Steps

1. Review and approve this proposal
2. Implement Option 1 (Field-Level Registry)
3. Update documentation
4. Update tests
5. Verify backward compatibility
