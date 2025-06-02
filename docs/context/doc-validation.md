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

## Validation Results Example: TEST-INFRA-001-A

### Archive Corruption Testing Framework Documentation Validation

This example demonstrates comprehensive documentation validation for a completed feature across all documentation layers.

#### Standard Validation Results
```bash
docs/validate-docs.sh

=== Documentation Validation Report ===

1. Feature Tracking Matrix Validation: ✅ PASSED
   - All feature IDs have valid references
   - No orphaned features detected  
   - Implementation tokens properly linked

2. Cross-Reference Validation: ✅ PASSED
   - All markdown links resolve correctly
   - No broken internal references
   - External URLs accessible

3. Implementation Token Validation: ✅ PASSED
   - All features have corresponding code markers
   - Token format matches required pattern
   - Implementation coverage complete

4. Test Function Validation: ✅ PASSED  
   - All referenced test functions exist
   - Test coverage matches documented requirements
   - No missing test implementations
```

#### DOC-001 Semantic Cross-Reference Validation
```bash
=== DOC-001: Semantic Cross-Reference Validation ===

Feature Reference Format Validation:
✅ TEST-INFRA-001-A: Complete feature reference block found
   - Feature: ✅ Link to feature-tracking.md
   - Spec: ✅ Link to specification.md  
   - Requirements: ✅ Link to requirements.md
   - Architecture: ✅ Link to architecture.md
   - Tests: ✅ Link to testing.md
   - Code: ✅ Implementation token reference

Cross-Document Consistency Analysis:
✅ TEST-INFRA-001-A: Found 8 references across documentation
   - feature-tracking.md: Feature definition and completion status
   - specification.md: User-facing behavior specification
   - requirements.md: Technical implementation requirements
   - architecture.md: System design and data models
   - testing.md: Comprehensive test coverage documentation
   - semantic-links-implementation.md: Cross-reference example
   - enhanced-traceability.md: Behavioral contracts and dependency mapping
   - sync-framework.md: Synchronization demonstration

Link Validation Results:
✅ All internal markdown links valid
✅ All anchor targets found
✅ No orphaned references detected
```

#### Feature-Specific Validation Checks

**Implementation Completeness**:
```
✅ Code Implementation: internal/testutil/corruption.go (736 lines)
✅ Test Implementation: internal/testutil/corruption_test.go (1004 lines)  
✅ Documentation Coverage: All 11 documentation files updated
✅ Cross-References: 8+ references across documentation layers
✅ Implementation Tokens: // TEST-INFRA-001-A markers in code
```

**Documentation Layer Verification**:
```
✅ Feature Tracking: Status marked as completed with implementation notes
✅ Specification: User-facing behavior documented with examples
✅ Requirements: Technical requirements with design decisions  
✅ Architecture: System design with data models and integration
✅ Testing: Comprehensive test coverage with all test functions
✅ Immutable: Behavioral contracts and stability guarantees
✅ Traceability: Dependency mapping and change impact analysis
✅ Semantic Links: Cross-reference example and validation
✅ Sync Framework: Synchronization demonstration
```

**Content Consistency Validation**:
```
✅ Performance Numbers: Consistent across all documents (CRC:763μs, Detection:49μs)
✅ Corruption Types: 8 types consistently documented across layers
✅ API Signatures: Function signatures match between architecture and implementation
✅ Example Code: Usage examples consistent with actual implementation
✅ Test Coverage: Test functions documented match actual implementation
```

#### Validation Framework Integration

The TEST-INFRA-001-A feature demonstrates how the validation framework ensures:

1. **Comprehensive Coverage**: Every documentation layer includes relevant information
2. **Cross-Reference Integrity**: Rich linking enables navigation between related concepts  
3. **Implementation Alignment**: Documentation matches actual code implementation
4. **Consistency Maintenance**: Information remains consistent across all documents
5. **Change Detection**: Updates propagate correctly across all affected layers

#### Validation Error Prevention

**Prevented Documentation Issues**:
- ❌ Orphaned feature references (cross-reference validation)
- ❌ Broken internal links (link validation)  
- ❌ Missing test functions (implementation validation)
- ❌ Inconsistent performance numbers (content validation)
- ❌ Incomplete feature documentation (coverage validation)

**Quality Assurance Metrics**:
```
Feature Coverage: 100% (8/8 documentation layers)
Cross-Reference Density: 8+ references per feature
Link Validation: 100% success rate
Implementation Alignment: 100% (code matches documentation)
Consistency Score: 100% (no conflicts detected)
```

This example demonstrates how the validation framework provides comprehensive quality assurance for documentation, ensuring that completed features are thoroughly documented across all layers with proper cross-references and consistent information. 