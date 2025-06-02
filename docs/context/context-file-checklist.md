# Context File Checklist

## Phase 1: Pre-Change Analysis (REQUIRED)

- [ ] **Feature Impact Analysis**: Identify which existing features are affected by the change
- [ ] **Context File Mapping**: Determine which context files require updates based on change type

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

## Phase 2: Documentation Updates (REQUIRED)

- [ ] **Update feature-tracking.md**: Add/modify feature entries, update status, add implementation tokens
- [ ] **Update specification.md**: Modify user-facing behavior descriptions if applicable
- [ ] **Update requirements.md**: Add/modify implementation requirements with traceability
- [ ] **Update architecture.md**: Update technical implementation descriptions
- [ ] **Update testing.md**: Add/modify test coverage requirements and descriptions
- [ ] **Check immutable.md**: Verify no immutable requirements are violated
- [ ] **Update cross-references**: Ensure all documents link correctly to each other

## Phase 3: Implementation Tokens (REQUIRED)

- [ ] **Add Implementation Tokens**: Mark ALL modified code with feature ID comments
- [ ] **Update Token Registry**: Record new tokens in feature-tracking.md
- [ ] **Verify Token Consistency**: Ensure tokens match feature IDs and descriptions

## Phase 4: Validation (REQUIRED)

- [ ] **Run Documentation Validation**: Execute `make validate-docs` (if available)
- [ ] **Cross-Reference Check**: Verify all links between documents are valid
- [ ] **Feature Matrix Update**: Ensure feature tracking matrix reflects all changes
- [ ] **Test Coverage Verification**: Confirm test references are accurate
