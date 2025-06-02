# Documentation Validation Checklist

## Purpose
This checklist ensures consistency and completeness across all documentation layers before any feature changes are committed.

## Pre-Change Validation

### 1. Immutable Requirements Check
- [ ] Feature is not listed in `immutable.md` unless this is a major version change
- [ ] If modifying immutable requirements, version bump is planned
- [ ] Backward compatibility impact is assessed
- [ ] Default values and naming conventions are preserved

### 2. Cross-Document Consistency
- [ ] Feature appears in all relevant documents (spec, requirements, architecture, testing)
- [ ] Terminology is consistent across all documents
- [ ] Command examples match actual implementation
- [ ] Configuration options are documented consistently
- [ ] Status codes are defined in all relevant places

### 3. Requirements Traceability
- [ ] Each requirement has corresponding test coverage defined
- [ ] Implementation approach is documented in architecture
- [ ] User-facing behavior is described in specification
- [ ] Error handling is documented with status codes

### 4. Test Coverage Verification
- [ ] New features have comprehensive test plans
- [ ] Existing test references are still valid
- [ ] Test names match documentation references
- [ ] Both positive and negative test cases are covered

### 5. Implementation Tokens
- [ ] Code will be marked with feature ID comments
- [ ] Decision references are documented
- [ ] Test references link to actual test functions
- [ ] Immutable requirement references are accurate

## Post-Change Validation

### 1. Documentation Sync Check
- [ ] All documents updated simultaneously
- [ ] No orphaned references to old behavior
- [ ] Feature matrix is updated
- [ ] Version control includes all document changes

### 2. Impact Assessment
- [ ] Dependent features identified and verified
- [ ] Breaking changes are documented
- [ ] Migration path is defined if needed
- [ ] Status code changes are justified

### 3. Testing Alignment
- [ ] Tests cover documented behavior
- [ ] Test output matches documented examples
- [ ] Error conditions match documented status codes
- [ ] Configuration examples work as documented

## Automation Opportunities

### Documentation Lint Rules
```bash
# Check for orphaned references
grep -r "TestXXX" docs/ | grep -v "# Example"

# Verify feature IDs are consistently formatted
grep -r "ARCH-[0-9]" docs/ | wc -l
grep -r "FILE-[0-9]" docs/ | wc -l
grep -r "CFG-[0-9]" docs/ | wc -l

# Check immutable requirements consistency
diff <(grep "^##" docs/context/immutable.md) <(grep "immutable" docs/context/*)
```

### Validation Script Template
```bash
#!/bin/bash
# docs/validate.sh

echo "Validating documentation consistency..."

# Check for missing cross-references
echo "Checking cross-references..."
# Implementation specific to your needs

# Verify feature matrix completeness
echo "Checking feature matrix..."
# Compare feature IDs across documents

# Validate test references
echo "Checking test references..."
# Verify test function names exist

echo "Validation complete."
```

## Common Documentation Errors

### 1. Inconsistent Naming
- **Problem**: Same feature called different names in different documents
- **Solution**: Establish canonical names in feature matrix
- **Prevention**: Use feature IDs consistently

### 2. Orphaned References
- **Problem**: Documentation references tests or functions that don't exist
- **Solution**: Validate references during review process
- **Prevention**: Automated reference checking

### 3. Missing Test Coverage
- **Problem**: Features documented without corresponding tests
- **Solution**: Require test plan before feature approval
- **Prevention**: Test-first documentation approach

### 4. Configuration Drift
- **Problem**: Configuration examples become outdated
- **Solution**: Extract examples from working configuration files
- **Prevention**: Automated configuration validation

### 5. Immutable Violations
- **Problem**: Accidentally changing immutable requirements
- **Solution**: Clear marking and review process for immutable changes
- **Prevention**: Automated immutable requirement protection

## Review Process

### For New Features
1. Review feature matrix completeness
2. Verify all five documents are updated
3. Check test coverage definition
4. Validate immutable requirement compliance
5. Confirm implementation token strategy

### For Existing Feature Changes
1. Assess immutable requirement impact
2. Review cross-document consistency
3. Verify test alignment
4. Check backward compatibility
5. Update decision records if needed

### For Documentation-Only Changes
1. Verify accuracy against implementation
2. Check cross-references
3. Update feature matrix if needed
4. Maintain terminology consistency
5. Preserve immutable requirement integrity 