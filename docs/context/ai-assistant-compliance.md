# AI Assistant Compliance Guide

## ⚡ MANDATORY TOKEN REFERENCING FOR ALL AI ASSISTANTS

### 🚨 CRITICAL AI DIRECTIVE: ALWAYS REFERENCE CONTEXT FILE TOKENS WITH CODE CHANGES

**🤖 IMMEDIATE ACTION REQUIRED FOR ALL AI ASSISTANTS:**

When making ANY code changes, you MUST:

1. **[ACTION:discovery] SEARCH FOR EXISTING TOKENS** [PRIORITY: CRITICAL]: Before making code changes, search the context directory files for relevant implementation tokens:
   - Use grep or search for patterns like `// ARCH-001:`, `// FILE-002:`, `// CFG-003:`, etc.
   - Look for feature IDs in the feature tracking tables in `feature-tracking.md`
   - Check all context files: `specification.md`, `requirements.md`, `architecture.md`, `testing.md`

2. **📋 REFERENCE RELEVANT TOKENS** [PRIORITY: CRITICAL]: For every code change, explicitly mention related tokens:
   ```
   Example: "This change affects ARCH-001 (Archive naming), FILE-002 (File backup), and CFG-003 (Format strings)"
   ```

3. **🔗 LINK TO DOCUMENTATION** [PRIORITY: HIGH]: Reference the specific sections in context files:
   ```
   Example: "See ARCH-001 in feature-tracking.md line 119, specification.md section 'Archive Naming', architecture.md 'Archive Service'"
   ```

4. **[ACTION:format-processing] UPDATE TOKEN REGISTRY** [PRIORITY: HIGH]: Add new implementation tokens for new features:
   - Format: `// FEATURE-ID: Brief description`
   - Register in the feature tracking tables in `feature-tracking.md`
   - Update all affected context files

5. **✅ VALIDATE CROSS-REFERENCES** [PRIORITY: MEDIUM]: Ensure all token references are consistent across:
   - `feature-tracking.md` 
   - `specification.md`
   - `requirements.md`
   - `architecture.md`
   - `testing.md`

### 🎯 TOKEN SEARCH PATTERNS FOR AI ASSISTANTS

**[ACTION:discovery] Search these patterns in context files [Execute in Order]:**
- `ARCH-[0-9]+` [CRITICAL] (Archive operations - CRITICAL)
- `FILE-[0-9]+` [CRITICAL] (File operations - CRITICAL)
- `CFG-[0-9]+` [HIGH] (Configuration - HIGH)
- `GIT-[0-9]+` [HIGH] (Git integration - HIGH)
- `TEST-[0-9]+` [HIGH] (Testing infrastructure - HIGH)
- `OUT-[0-9]+` [MEDIUM] (Output management - MEDIUM)
- `DOC-[0-9]+` [HIGH] (Documentation system - HIGH)
- `LINT-[0-9]+` [HIGH] (Code quality - HIGH)
- `COV-[0-9]+` [HIGH] (Coverage - HIGH)

### [ACTION:core-functionality] DOC-007: STANDARDIZED IMPLEMENTATION TOKEN REQUIREMENTS

> **[ACTION:validation] MANDATORY**: All AI assistants MUST use the standardized icon system in implementation tokens as established by DOC-007 (Source Code Icon Integration).

#### ✨ **Required Implementation Token Format**
```go
// [PRIORITY_ICON] FEATURE-ID: Brief description [- ACTION_ICON Context]
```

**Examples:**
```go
// [CRITICAL] ARCH-001: Archive naming convention implementation - [ACTION:core-functionality] Core functionality
// [HIGH] CFG-003: Template formatting logic - [ACTION:format-processing] Format string processing  
// [MEDIUM] GIT-004: Git submodule support - [ACTION:discovery] Discovery and validation
// [LOW] REFACTOR-001: Interface standardization - 🏗️ Architecture preparation
```

#### 🔥 **Priority Icon Requirements**
**AI assistants MUST use correct priority icons based on feature priority in feature-tracking.md:**

| Priority Icon | Feature Priority | When to Use in Code |
|---------------|------------------|---------------------|
| [CRITICAL] | CRITICAL | Core system functionality, blocking operations, essential data integrity |
| [HIGH] | HIGH | Important features, significant configuration, core business logic |
| [MEDIUM] | MEDIUM | Conditional features, optional enhancements, secondary processing |
| [LOW] | LOW | Cleanup functions, optional optimization, debug utilities |

#### [ACTION:core-functionality] **Action Category Icon Requirements**
**AI assistants MUST include appropriate action icons based on function behavior:**

| Action Icon | Action Category | When to Use in Code |
|-------------|-----------------|---------------------|
| [ACTION:discovery] | SEARCH/DISCOVER | File system discovery, configuration search, Git analysis, pattern matching |
| [ACTION:format-processing] | DOCUMENT/UPDATE | Configuration writing, log creation, status updates, output formatting |
| [ACTION:core-functionality] | CONFIGURE/MODIFY | System configuration, parameter adjustment, environment setup |
| [ACTION:validation] | PROTECT/VALIDATE | Input validation, security checks, data integrity, error detection |

#### 🏗️ **Component Category Icons**
**AI assistants SHOULD include component context when relevant:**

| Component Icon | Component Type | When to Use in Code |
|----------------|----------------|---------------------|
| 🏗️ | ARCHITECTURE | Core system components, service layers, interface definitions |
| 🧪 | TESTING | Test infrastructure, mock implementations, test data generation |

### [ACTION:validation] DOC-008: COMPREHENSIVE ICON VALIDATION AND ENFORCEMENT

> **[HIGH] MANDATORY**: All AI assistants MUST validate icon usage using the DOC-008 comprehensive validation and enforcement system before submitting any changes.

#### 🚨 **Mandatory Pre-Submit Validation**
**Before submitting ANY changes, AI assistants MUST:**

1. **[ACTION:discovery] Run Icon Consistency Check**: Execute `make validate-icons` (DOC-007 basic validation)
2. **[ACTION:validation] Run Comprehensive Validation**: Execute `make validate-icon-enforcement` (DOC-008 full validation)
3. **📊 Review Validation Report**: Check `docs/validation-reports/icon-validation-report.md` for warnings and errors
4. **✅ Achieve Zero Errors**: Fix all critical errors before proceeding

#### 🆕 **DOC-009 THROUGH DOC-012: AI ASSISTANT OPTIMIZATION TASKS**

> **⚡ CRITICAL for AI Assistant Workflows**: The following tasks are specifically designed to optimize the development experience for AI assistants and maintain icon system integrity in an AI-first development environment.

**[HIGH] DOC-009: Mass Implementation Token Standardization (CRITICAL)**
- **Status**: Not Started - **IMMEDIATE PRIORITY**
- **Impact**: Currently 562 validation warnings affect all AI assistant code changes
- **AI Benefit**: 90%+ standardization rate will enable seamless code navigation and comprehension
- **Timeline**: 2-4 weeks to process 592 legacy tokens across 47 files
- **AI Assistant Role**: May be requested to assist with batch token updates and validation

**[HIGH] DOC-011: Token Validation Integration for AI Assistants (HIGH)**
- **Status**: Not Started - **HIGH PRIORITY** 
- **Impact**: Zero-friction validation integration in AI assistant workflows
- **AI Benefit**: Pre-submission validation with intelligent error reporting optimized for AI comprehension
- **Dependencies**: DOC-009 completion for clean baseline
- **AI Assistant Role**: Direct beneficiary of seamless workflow integration

**[MEDIUM] DOC-010: Automated Token Format Suggestions (MEDIUM)**
- **Status**: Not Started - **QUALITY ENHANCEMENT**
- **Impact**: 95%+ accuracy in token format suggestions for new code
- **AI Benefit**: Reduced manual token formatting effort, consistent creation patterns
- **Dependencies**: DOC-009 provides training data from standardized patterns
- **AI Assistant Role**: Primary user of intelligent token suggestions

**[MEDIUM] DOC-012: Real-time Icon Validation Feedback (MEDIUM)**
- **Status**: Not Started - **ADVANCED ENHANCEMENT**
- **Impact**: Sub-second validation feedback during code writing
- **AI Benefit**: Proactive compliance maintenance, reduced validation friction
- **Dependencies**: DOC-008 engine, DOC-011 integration patterns
- **AI Assistant Role**: Enhanced development experience with immediate feedback

#### 📋 **DOC-008 Validation Categories**
**The comprehensive validation system checks:**

| Validation Area | Description | Failure Action |
|-----------------|-------------|----------------|
| **Master Icon Legend** | All icons conform to master legend in README.md | Fix invalid icons |
| **Documentation Consistency** | Icons used consistently across all context docs | Remove conflicting usage |
| **Implementation Tokens** | Code tokens follow standardized format | Update to standard format |
| **Cross-References** | Feature IDs consistent between docs and code | Sync documentation |
| **Enforcement Rules** | Validation infrastructure properly integrated | Complete setup |

#### ⚡ **Automation Integration Requirements**
**AI assistants MUST ensure:**

- ✅ Makefile integration: `validate-icon-enforcement` target exists and works
- ✅ Script accessibility: `scripts/validate-icon-enforcement.sh` is executable
- ✅ CI/CD readiness: `make validate-icons-strict` passes in strict mode
- ✅ Report generation: `docs/validation-reports/icon-validation-report.md` created after validation

#### [ACTION:core-functionality] **DOC-008 Enforcement Levels**

| Mode | Command | Usage | Fail Conditions |
|------|---------|-------|-----------------|
| **Standard** | `make validate-icon-enforcement` | Development validation | Errors > 0 |
| **Strict** | `make validate-icons-strict` | CI/CD pipelines | Errors > 0 OR Warnings > 5 |
| **Legacy** | `make validate-icons` | DOC-007 compatibility | DOC-007 specific checks |

### [HIGH] DOC-014: MANDATORY AI ASSISTANT DECISION FRAMEWORK

> **[CRITICAL] CRITICAL**: All AI assistants MUST use the AI Assistant Decision Framework (DOC-014) for ALL code changes. The framework provides explicit decision-making principles to ensure AI assistant alignment with project goals.

#### 🚨 **Mandatory Decision Framework Usage**
**AI assistants MUST follow the Decision Framework defined in [`ai-decision-framework.md`](ai-decision-framework.md) for every code change:**

1. **[ACTION:validation] Safety Gates**: Check test validation, backward compatibility, token compliance, validation scripts
2. **📑 Scope Boundaries**: Verify feature scope, dependency satisfaction, architecture alignment, context updates
3. **📊 Quality Thresholds**: Ensure test coverage, error pattern consistency, performance stability, traceability
4. **🎯 Goal Alignment**: Confirm phase progress, priority order, future impact, reusability preservation

#### 🧠 **Decision Tree Compliance**
**AI assistants MUST use appropriate decision trees for common scenarios:**

| Scenario | Decision Tree | Required Validation |
|----------|---------------|-------------------|
| **Feature Implementation** | "Should I implement this feature request?" | Feature exists in feature-tracking.md, dependencies satisfied, architecture alignment |
| **Code Refactoring** | "Should I refactor this code?" | Part of planned tasks, improves extraction goals, enhances AI comprehension |
| **Test Failures** | "Should I fix this test failure?" | Blocking work, critical component, minimal scope change |
| **Documentation Updates** | "Should I update documentation?" | Required by protocol, affects behavior, impacts architecture |

#### 🎯 **Decision Validation Requirements**
**AI assistants MUST execute decision validation checklists:**

**Pre-Implementation Validation**:
```bash
# [ACTION:validation] Safety Gates Validation
make test                    # All tests must pass
make lint                    # All lint checks must pass  
make validate-icons          # Icon validation must pass
git status                   # Working directory must be clean

# 📋 Scope Validation  
grep -r "FEATURE-ID" docs/context/feature-tracking.md  # Feature must exist
# Check dependencies in feature-tracking.md are "✅ Completed"
# Verify change aligns with current extraction/refactoring phase
```

**Post-Implementation Validation**:
```bash
# 📊 Quality Validation
make test-coverage-validate  # Coverage thresholds must be met
make test                    # All tests must still pass
make validate-icons          # Token format must be compliant

# [ACTION:format-processing] Documentation Validation
# Verify all required context files updated per protocol
# Check implementation tokens added to all modified code
# Confirm feature status updated in feature-tracking.md
```

#### [ACTION:core-functionality] **Framework Integration with Existing Systems**
**Decision Framework MUST integrate with:**

- **🔗 Feature Tracking**: Respect [CRITICAL][HIGH][MEDIUM][LOW] priority hierarchy, check blocking dependencies
- **[ACTION:core-functionality] Change Protocols**: Use decision trees to select appropriate ai-assistant-protocol.md workflow
- **[ACTION:validation] Validation Systems**: Leverage DOC-008/DOC-011 for compliance checking and real-time feedback
- **🏷️ Token System**: Include decision context in implementation tokens using DOC-007 standardized format

#### 📊 **Decision Quality Metrics**
**AI assistants MUST achieve:**

- **🎯 Goal Alignment Rate**: >95% of changes advance documented project goals
- **🔗 Traceability Compliance**: 100% of changes have Feature IDs and implementation tokens
- **📋 Protocol Adherence**: 100% of changes follow appropriate ai-assistant-protocol.md workflows
- **[ACTION:validation] Regression Prevention**: Zero test failures introduced by AI assistant changes

### 🚨 AI ENFORCEMENT RULES

**❌ REJECT code changes that do not:**
- Reference at least one existing feature token when modifying related functionality
- Propose new tokens for genuinely new features
- Link to relevant documentation sections
- Update the feature tracking tables when adding new tokens
- **🆕 DOC-007**: Use standardized icon format in implementation tokens
- **🆕 DOC-007**: Include correct priority icons matching feature-tracking.md
- **🆕 DOC-007**: Include appropriate action category icons
- **🆕 DOC-008**: Pass comprehensive icon validation (`make validate-icon-enforcement`)
- **🆕 DOC-008**: Generate clean validation report with zero critical errors
- **🚨 DOC-008**: Update BOTH feature registry table AND detailed subtask blocks when marking tasks completed
- **[CRITICAL] DOC-014**: Follow the AI Assistant Decision Framework for all code changes
- **[CRITICAL] DOC-014**: Execute pre-implementation validation checklists (safety gates, scope boundaries)
- **[CRITICAL] DOC-014**: Use appropriate decision trees for implementation scenarios
- **[CRITICAL] DOC-014**: Achieve decision quality metrics (>95% goal alignment, 100% traceability)

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

**✅ ACCEPT code changes that:**
- Clearly state which tokens are affected
- Reference specific documentation sections
- Propose updates to context files when needed
- Follow the token naming conventions
- **🆕 DOC-007**: Use standardized implementation token format with icons
- **🆕 DOC-007**: Maintain icon consistency between code and documentation
- **🆕 DOC-008**: Pass all validation categories in DOC-008 comprehensive system
- **🆕 DOC-008**: Include validation report summary in change description
- **🚨 DOC-008**: Demonstrate completion consistency between registry table and subtask blocks
- **[CRITICAL] DOC-014**: Demonstrate Decision Framework compliance with decision context and validation
- **[CRITICAL] DOC-014**: Include decision tree rationale in implementation token comments
- **[CRITICAL] DOC-014**: Show pre/post-implementation validation results
- **[CRITICAL] DOC-014**: Document goal alignment and project impact in change description

## 🤖 MANDATORY AI VALIDATION CHECKLIST

Before any AI assistant provides code changes, it MUST validate:

1. **[ACTION:discovery] TOKEN SEARCH PERFORMED** [PRIORITY: CRITICAL]
   - [ ] Searched all context files for existing feature tokens
   - [ ] Identified all relevant feature IDs using grep patterns: `ARCH-[0-9]+`, `FILE-[0-9]+`, `CFG-[0-9]+`, etc.
   - [ ] Listed affected tokens in response

2. **📋 TOKEN IMPACT DOCUMENTED** [PRIORITY: CRITICAL]
   - [ ] Clearly stated which existing tokens are affected
   - [ ] Proposed new tokens for new functionality
   - [ ] Referenced specific line numbers in feature-tracking.md
   - [ ] Linked to relevant sections in other context files

3. **[ACTION:core-functionality] DOC-007 ICON COMPLIANCE** [PRIORITY: CRITICAL]
   - [ ] Used standardized implementation token format with priority and action icons
   - [ ] Verified priority icons match feature priority in feature-tracking.md
   - [ ] Included appropriate action category icons ([ACTION:discovery][ACTION:format-processing][ACTION:core-functionality][ACTION:validation]) based on function behavior
   - [ ] Ensured icon consistency between code and documentation system

4. **[ACTION:validation] DOC-008 COMPREHENSIVE VALIDATION** [PRIORITY: CRITICAL]
   - [ ] Executed `make validate-icons` for basic DOC-007 validation
   - [ ] Executed `make validate-icon-enforcement` for comprehensive DOC-008 validation
   - [ ] Reviewed generated `docs/validation-reports/icon-validation-report.md` for issues
   - [ ] Achieved zero critical errors in validation report
   - [ ] Addressed all warnings or documented acceptable risk
   - [ ] Verified all icons conform to master legend in README.md

5. **📄 DOCUMENTATION CROSS-REFERENCES** [PRIORITY: HIGH]
   - [ ] Referenced specification.md sections
   - [ ] Referenced requirements.md areas
   - [ ] Referenced architecture.md components
   - [ ] Referenced testing.md coverage areas

6. **[ACTION:migration] UPDATE REQUIREMENTS IDENTIFIED** [PRIORITY: MEDIUM]
   - [ ] Listed which context files need updates
   - [ ] Specified what changes are needed in each file
   - [ ] Provided update templates or examples

7. **[CRITICAL] DOC-014 DECISION FRAMEWORK COMPLIANCE** [PRIORITY: CRITICAL]
   - [ ] Selected and documented appropriate decision tree for implementation scenario
   - [ ] Executed pre-implementation validation checklist (safety gates, scope boundaries)
   - [ ] Verified feature exists in feature-tracking.md with satisfied dependencies
   - [ ] Confirmed change aligns with current project phase and priorities
   - [ ] Documented decision rationale and goal alignment in implementation tokens
   - [ ] Executed post-implementation validation (quality thresholds, documentation updates)
   - [ ] Achieved >95% goal alignment rate and 100% traceability compliance
   - [ ] Integrated framework validation with existing DOC-007/008 systems

8. **🏁 TASK COMPLETION CONSISTENCY** [PRIORITY: CRITICAL]
   - [ ] Updated feature registry table status to "Completed"
   - [ ] Updated ALL detailed subtask blocks with checkmarks [x]
   - [ ] Verified consistency between table and subtask completion status
   - [ ] Ensured DOC-008 enforcement compliance

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

## 📋 AI RESPONSE TEMPLATE (REQUIRED FORMAT)

```
## 🎯 Token Impact Analysis

### [ACTION:discovery] Affected Tokens:
- **🚨 Primary**: [TOKEN-ID] - [Description of impact]
- **📊 Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [X], [Y], [Z]
- **specification.md**: Section "[Section Name]"
- **requirements.md**: "[Requirement area]"
- **architecture.md**: "[Component name]"
- **testing.md**: "[Test coverage area]"

### [ACTION:validation] DOC-008 Validation Results:
- **Basic Validation (DOC-007)**: [✅ PASSED | ❌ FAILED - Details]
- **Comprehensive Validation (DOC-008)**: [✅ PASSED | ❌ FAILED - Details] 
- **Critical Errors**: [Count and brief description]
- **Warnings**: [Count and action plan]
- **Validation Report**: [Reference to docs/validation-reports/icon-validation-report.md]

### [ACTION:migration] Required Context File Updates:
- [ ] feature-tracking.md: [Specific changes]
- [ ] specification.md: [Specific changes]
- [ ] requirements.md: [Specific changes]
- [ ] architecture.md: [Specific changes]
- [ ] testing.md: [Specific changes]

### 💻 Implementation Details:
- **🆕 New tokens**: [TOKEN-ID: Description]
- **[ACTION:core-functionality] Modified tokens**: [TOKEN-ID: Updated description]
- **📁 Code files affected**: [List of files]
```

## 🚨 AI REJECTION CRITERIA (ZERO TOLERANCE)

AI assistants MUST REJECT their own responses if they:
- ❌ Fail to search for existing tokens before proposing changes
- ❌ Make code changes without referencing related feature tokens
- ❌ Don't identify which context files need updates
- ❌ Propose changes that conflict with immutable requirements
- ❌ Don't follow the required response template format
- ❌ **DOC-008**: Fail comprehensive icon validation with critical errors
- ❌ **DOC-008**: Skip mandatory validation steps before submitting changes
- ❌ **DOC-008**: Use icons not defined in the master legend (README.md)
- ❌ **DOC-014**: Skip Decision Framework validation and decision tree selection
- ❌ **DOC-014**: Proceed without pre-implementation safety gates validation
- ❌ **DOC-014**: Make changes that don't advance documented project goals
- ❌ **DOC-014**: Fail to achieve decision quality metrics (goal alignment, traceability)

## ✅ AI APPROVAL CRITERIA

AI assistants should ONLY PROVIDE responses that:
- ✅ Include complete token impact analysis
- ✅ Reference specific documentation sections with line numbers
- ✅ Propose concrete updates to all affected context files
- ✅ Follow established naming conventions for new tokens
- ✅ Include implementation token placement in code changes
- ✅ **DOC-008**: Pass comprehensive icon validation with zero critical errors
- ✅ **DOC-008**: Include validation report summary in response
- ✅ **DOC-008**: Demonstrate icon consistency across documentation and code
- ✅ **DOC-014**: Execute Decision Framework validation with documented decision tree selection
- ✅ **DOC-014**: Complete pre-implementation safety gates and scope boundary validation
- ✅ **DOC-014**: Demonstrate goal alignment and project impact with >95% alignment rate
- ✅ **DOC-014**: Include decision context and rationale in implementation tokens

## [ACTION:migration] ENFORCEMENT SUMMARY

**[ACTION:validation] REMEMBER: THIS IS NOT OPTIONAL**

The feature-tracking system is designed to prevent unplanned changes and maintain documentation synchronization. AI assistants that fail to follow these requirements undermine the entire system integrity.

**🎯 Key Success Metrics:**
- ✅ 100% of code changes must reference related tokens
- ✅ 100% of new features must have new tokens
- ✅ 100% of responses must include context file update plans
- ❌ 0% tolerance for undocumented changes

**🏷️ Implementation Tokens Are Mandatory:**
Every function, method, and significant code block must include `// TOKEN-ID: Description` comments that trace back to the feature tracking matrix in `feature-tracking.md`.

## 📞 QUICK REFERENCE FOR AI ASSISTANTS

### [ACTION:discovery] Search Commands
1. **🔎 Search Pattern**: `grep -r "// [A-Z]+-[0-9]+" docs/context/`
2. **🏷️ Token Format**: `// FEATURE-ID: Brief description`
3. **📋 Required Template**: Use the Token Impact Analysis template above
4. **📁 Context Files**: feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md
5. **🚨 Zero Tolerance**: No code changes without token references

### 🎯 Priority Execution Order
1. **🚨 CRITICAL** - Search and identify existing tokens
2. **📋 HIGH** - Document token impact analysis
3. **📊 MEDIUM** - Cross-reference documentation
4. **[ACTION:format-processing] LOW** - Propose file updates

### 🚀 Quick Validation
```bash
# Pre-change token search
grep -r "FEATURE-ID" docs/context/feature-tracking.md

# Post-change validation
make test && make lint
```

**🤖 This document serves as the authoritative source for AI assistant compliance with the feature tracking system.** 