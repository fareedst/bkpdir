# AI Assistant Decision Framework

> **[CRITICAL] CRITICAL PURPOSE**: This document codifies explicit decision-making principles to keep AI assistants aligned with project goals when working on implementation details.

## 🚨 DECISION HIERARCHY (Check in Order)

### 4-Tier Decision Hierarchy

### 1. **[ACTION:validation] Safety Gates** (NEVER Override)
- **🧪 Test Validation**: Are all tests passing? (`make test`)
- **[ACTION:core-functionality] Backward Compatibility**: Does the change preserve existing functionality?
- **🏷️ Token Compliance**: Are implementation tokens properly formatted? (DOC-007/008/009)
- **📋 Validation Scripts**: Do validation scripts pass? (`make validate-icons`)

### 2. **📑 Scope Boundaries** (Strict Limits)
- **📋 Feature Scope**: Is this change within the current feature's documented scope?
- **🔗 Dependency Check**: Are all blocking dependencies satisfied in feature-tracking.md?
- **🏗️ Architecture Alignment**: Does this align with documented architecture patterns?
- **[ACTION:format-processing] Context Updates**: Will this require updating multiple context files per protocol?

### 3. **📊 Quality Thresholds** (Must Meet)
- **📈 Test Coverage**: Does this maintain >90% test coverage for new code? (COV-001)
- **[ACTION:validation] Error Patterns**: Are error handling patterns consistent with existing code?
- **🚀 Performance**: Do performance benchmarks remain stable?
- **🏷️ Traceability**: Are implementation tokens added for bidirectional traceability?

### 4. **🎯 Goal Alignment** (Strategic Check)
- **⚡ Phase Progress**: Does this advance the current extraction/refactoring phase?
- **[CRITICAL] Priority Order**: Is this the highest priority task available in feature-tracking.md?
- **🔮 Future Impact**: Will this enable future work or create technical debt?
- **[ACTION:core-functionality] Reusability**: Does this preserve the component extraction and reusability goals?

## 🧠 DECISION TREE FOR COMMON SCENARIOS

### Decision Trees

### **"Should I implement this feature request?"**
```
1. ✅ Feature exists in feature-tracking.md with status "[ACTION:format-processing] Not Started"
2. ✅ All blocking dependencies marked "✅ Completed"  
3. ✅ Implementation aligns with documented architecture (architecture.md)
4. ✅ Required context files identified per ai-assistant-protocol.md
➡️ **PROCEED** with full NEW FEATURE Protocol execution
```

### **"Should I refactor this code?"**
```
1. ❓ Is this part of planned REFACTOR-001-006 tasks?
2. ❓ Does current code structure prevent extraction goals (EXTRACT-001-010)?
3. ❓ Will this improve AI assistant code comprehension?
4. ❓ Is technical debt blocking current phase progress?
⚠️ **EVALUATE** impact and get validation through proper protocol
```

### **"Should I fix this test failure?"**
```
1. 🚨 Is this blocking other work or extraction tasks?
2. 🚨 Is this in a critical component ([CRITICAL] or [HIGH] priority)?
3. ✅ Can I fix without changing documented functionality?
4. ✅ Does fix follow established error handling patterns?
➡️ **IMMEDIATE** fix with BUG FIX Protocol (minimal scope)
```

### **"Should I update documentation?"**
```
1. ✅ Is this required by ai-assistant-protocol.md for my change type?
2. ✅ Does the change affect user-facing behavior (specification.md)?
3. ✅ Does the change affect architecture or interfaces?
4. ✅ Are implementation tokens updated to reflect changes?
➡️ **REQUIRED** documentation updates per protocol
```

## 🔒 ENHANCED ENFORCEMENT THROUGH IMPLEMENTATION TOKENS

### 🔒 ENHANCED ENFORCEMENT

### **Current Token System Strengths**
- **🔗 Bidirectional Traceability**: `// FEATURE-ID: Description` links code ↔ documentation
- **[CRITICAL] Priority Hierarchy**: Visual priority icons ([CRITICAL][HIGH][MEDIUM][LOW]) communicate execution order
- **[ACTION:core-functionality] Action Categories**: Action icons ([ACTION:discovery][ACTION:format-processing][ACTION:core-functionality][ACTION:validation]) clarify implementation purpose
- **📋 Validation Integration**: DOC-008 automated validation prevents token drift

### **🆕 Enhancement: Decision Context in Tokens**
Add decision rationale to implementation tokens for enhanced AI assistant guidance:

```go
// [CRITICAL] ARCH-001: Archive naming [DECISION: core-functionality, blocks-extraction, user-facing]
// [HIGH] CFG-005: Inheritance [DECISION: enhancement, enables-flexibility, backward-compatible]
// [MEDIUM] DOC-010: Token suggestions [DECISION: developer-experience, ai-assistance]
// [LOW] TEST-FIX-001: Config isolation [DECISION: test-reliability, infrastructure]
```

**Decision Context Categories:**
- **Impact Level**: `core-functionality`, `enhancement`, `developer-experience`, `infrastructure`
- **Dependencies**: `blocks-extraction`, `enables-flexibility`, `prerequisite-for-X`
- **Constraints**: `user-facing`, `backward-compatible`, `breaking-change`, `performance-critical`

## 🎯 DECISION VALIDATION CHECKLIST

### Validation Checklists

### **Pre-Implementation Validation**
```bash
# [ACTION:validation] Safety Gates
make test                    # All tests must pass
make lint                    # All lint checks must pass  
make validate-icons          # Icon validation must pass
git status                   # Working directory must be clean

# 📋 Scope Validation
grep -r "FEATURE-ID" docs/context/feature-tracking.md  # Feature must exist
# Check dependencies in feature-tracking.md are "✅ Completed"
# Verify change aligns with current extraction/refactoring phase
```

### **Post-Implementation Validation**
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

## Integration Points

### **🔗 Feature Tracking Integration** (feature-tracking.md)
- **Priority System**: Respect [CRITICAL][HIGH][MEDIUM][LOW] priority hierarchy in execution order
- **Dependency Chains**: Never work on tasks with incomplete blocking dependencies
- **Status Tracking**: Update both registry table AND detailed subtask blocks (DOC-008 requirement)

### **[ACTION:core-functionality] Protocol Integration** (ai-assistant-protocol.md)
- **Change Classification**: Use decision tree to select appropriate protocol (NEW FEATURE, MODIFICATION, etc.)
- **Documentation Cascade**: Follow required context file updates per change type
- **Validation Requirements**: Execute all mandatory validation steps

### **[ACTION:validation] Validation Integration** (DOC-008, DOC-011)
- **Pre-submission Validation**: Use ai-validation CLI for zero-friction compliance checking
- **Icon Compliance**: Maintain 100% standardization rate achieved by DOC-009
- **Real-time Feedback**: Leverage DOC-012 for immediate compliance feedback

## Success Metrics

### **📈 Decision Quality Metrics**
- **🎯 Goal Alignment Rate**: >95% of implemented changes advance documented project goals
- **🔗 Traceability Compliance**: 100% of code changes have corresponding Feature IDs and implementation tokens
- **📋 Protocol Adherence**: 100% of changes follow appropriate ai-assistant-protocol.md workflows
- **[ACTION:validation] Regression Prevention**: Zero test failures introduced by AI assistant changes

### **⚡ Efficiency Metrics**
- **🚀 First-Time Success Rate**: >90% of AI assistant changes pass validation on first submission
- **[ACTION:migration] Rework Minimization**: <5% of changes require significant rework due to scope/goal misalignment
- **[ACTION:format-processing] Documentation Completeness**: 100% of required context files updated per protocol
- **⏱️ Decision Time**: <2 minutes to classify change type and select appropriate protocol

### **🎯 Strategic Alignment Metrics**
- **📊 Phase Progress**: AI assistant changes consistently advance current project phase (extraction/refactoring)
- **[ACTION:core-functionality] Architecture Consistency**: Zero architectural decisions that conflict with documented patterns
- **🔮 Technical Debt**: AI assistant changes reduce rather than increase technical debt
- **🤖 AI Comprehension**: Enhanced code navigation and understanding through consistent token usage

## [ACTION:discovery] TROUBLESHOOTING COMMON DECISION SCENARIOS

### **🚨 Conflicting Priorities**
**Scenario**: Multiple [CRITICAL] CRITICAL tasks available
**Decision Process**:
1. Check feature-tracking.md for current project phase
2. Prioritize extraction/refactoring tasks over new features
3. Choose tasks that unblock the most downstream work
4. Consider AI assistant expertise and change complexity

### **⚠️ Scope Uncertainty**
**Scenario**: Change request spans multiple features
**Decision Process**:
1. Break down into individual Feature IDs
2. Identify minimum viable change that provides value
3. Check if change requires new Feature ID creation
4. Follow MODIFICATION Protocol for existing features

### **[ACTION:core-functionality] Technical Debt vs Feature Work**
**Scenario**: Technical debt blocks feature implementation
**Decision Process**:
1. Classify as REFACTORING if improves extraction readiness
2. Classify as BUG FIX if fixes broken functionality
3. Create separate Feature ID if substantial architectural change
4. Prioritize based on extraction timeline impact

### **[ACTION:format-processing] Documentation Scope Questions**
**Scenario**: Uncertain which context files to update
**Decision Process**:
1. Use ai-assistant-protocol.md change type classification
2. Follow documentation cascade requirements exactly
3. When in doubt, update more rather than fewer files
4. Validate with ai-validation CLI before submission

## 🎯 IMPLEMENTATION ROADMAP

### **⚡ Phase 1: Decision Framework Integration (Week 1)**
- [x] Create ai-decision-framework.md document ✅
- [ ] Add DOC-014 entry to feature-tracking.md
- [ ] Update ai-assistant-compliance.md to reference decision framework
- [ ] Integrate decision checklist into ai-assistant-protocol.md

### **[ACTION:core-functionality] Phase 2: Enhanced Token Context (Week 2)**
- [ ] Define decision context categories and syntax
- [ ] Update implementation token guidelines (DOC-007)
- [ ] Create migration script for enhanced token format
- [ ] Update validation scripts to check decision context

### **📊 Phase 3: Validation and Metrics (Week 3)**
- [ ] Implement decision validation tools
- [ ] Create decision quality metrics dashboard
- [ ] Add decision context to real-time validation (DOC-012)
- [ ] Establish baseline metrics for success measurement

### **🚀 Phase 4: Integration and Testing (Week 4)**
- [ ] Comprehensive testing of decision framework
- [ ] Integration testing with existing validation systems
- [ ] Documentation review and refinement
- [ ] Training material creation for AI assistant onboarding

## [CRITICAL] CRITICAL SUCCESS FACTORS

1. **[ACTION:validation] Mandatory Adoption**: All AI assistants MUST use this framework for ANY code changes
2. **📋 Documentation Integration**: Framework must be referenced from ai-assistant-compliance.md
3. **[ACTION:core-functionality] Tool Integration**: Decision validation must be integrated into existing validation workflows
4. **📊 Continuous Improvement**: Framework must evolve based on AI assistant usage patterns and effectiveness metrics
5. **🎯 Goal Alignment**: Framework success measured by improved project goal achievement and reduced rework

**🤖 This framework transforms implicit decision-making into explicit, teachable processes that ensure AI assistant compliance with project goals while maintaining development velocity and code quality.** 