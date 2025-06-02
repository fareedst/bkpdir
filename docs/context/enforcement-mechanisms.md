# Enforcement Mechanisms

## Automated Validation Rules

1. **Feature ID Consistency**
   - All implementation tokens must correspond to feature matrix entries
   - Tokens must match feature IDs in feature-tracking.md
   - Token descriptions must align with feature descriptions

2. **Cross-Reference Integrity**
   - All document links must resolve correctly
   - Feature references must be valid across all context files
   - Implementation tokens must link to correct feature entries

3. **Status Synchronization**
   - Feature status must be consistent across all documents
   - Implementation status must match documentation status
   - Test coverage status must align with implementation status

4. **Test Traceability**
   - All features must have corresponding test references
   - Test coverage must be documented in testing.md
   - Test descriptions must match feature requirements

## Manual Review Requirements

1. **Context Documentation Review**
   - Every code change must include review of context file updates
   - Changes must be reflected in all relevant context files
   - Documentation must be updated before code changes

2. **Feature Impact Assessment**
   - Verify that all affected features are properly documented
   - Check cross-document consistency for affected features
   - Update all relevant context files for impacted features

3. **Immutable Requirement Check**
   - Confirm no immutable requirements are violated
   - Verify backward compatibility constraints
   - Check against immutable.md requirements

4. **Backward Compatibility Verification**
   - Ensure changes preserve documented compatibility
   - Verify no breaking changes for existing users
   - Document any necessary migration guidance

## Common Mistakes to Avoid

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

## Enforcement Reminders

- **NO CODE REVIEW** without context file updates
- **NO MERGE** without feature tracking matrix updates  
- **NO DEPLOYMENT** without documentation validation
- **NO EXCEPTIONS** - context files are as important as code
