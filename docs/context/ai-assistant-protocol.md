# AI Assistant Feature Update Protocol

> **🤖 MANDATORY FOR ALL AI ASSISTANTS**: This document defines the REQUIRED protocol for all code changes. Failure to follow this protocol invalidates any code modifications.

## 🚨 CRITICAL COMPLIANCE REQUIREMENTS

### **[ACTION:validation] MANDATORY PRE-WORK VALIDATION** (Priority 1 - Must Execute First)
Before making ANY code changes:

1. **📋 Task Verification**: The task MUST exist in `feature-tracking.md` with a valid Feature ID
2. **[ACTION:discovery] Compliance Check**: Review `ai-assistant-compliance.md` for mandatory token referencing requirements
3. **📁 File Impact Analysis**: Use this guide to determine which documentation files require updates
4. **[ACTION:validation] Immutable Check**: Verify no changes violate requirements in `immutable.md`

### **✅ MANDATORY POST-WORK COMPLETION** (Priority 2 - Must Execute Last)
After ALL code changes are complete:

1. **[CRITICAL] Decision Framework Validation**: All changes must comply with DOC-014 Decision Framework (see [`ai-decision-framework.md`](ai-decision-framework.md))
2. **🧪 Full Test Suite**: All tests must pass (`make test`)
3. **[ACTION:core-functionality] Lint Compliance**: All lint checks must pass (`make lint`)
4. **[ACTION:format-processing] Documentation Updates**: All required documentation files updated per this protocol
5. **🏁 Task Completion**: Update task status to "Completed" in `feature-tracking.md`

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

## 🤖 AI Assistant Decision Engine

### [ACTION:discovery] Step 1: Task Identification and Validation (EXECUTE FIRST)
```
REQUIRED: Before any code modification
├─ 📋 Does Feature ID exist in feature-tracking.md? → YES: Continue | NO: STOP - Create task first
├─ [ACTION:discovery] Are implementation tokens required? → Check ai-assistant-compliance.md
├─ 📁 Which files need documentation updates? → Use workflows below
└─ [ACTION:validation] Any immutable requirement conflicts? → Check immutable.md
```

### 🎯 Step 2: Change Type Classification (EXECUTE SECOND)
```
What type of change are you making?
├─ 🆕 NEW FEATURE → Execute: NEW FEATURE Protocol [Priority: HIGH]
├─ [ACTION:core-functionality] MODIFY EXISTING → Execute: MODIFICATION Protocol [Priority: HIGH]
├─ 🐛 BUG FIX → Execute: BUG FIX Protocol [Priority: MEDIUM]
├─ ⚙️ CONFIG CHANGE → Execute: CONFIG CHANGE Protocol [Priority: MEDIUM]
├─ 🔌 API/INTERFACE → Execute: API CHANGE Protocol [Priority: HIGH]
├─ 🧪 TEST ONLY → Execute: TEST ADDITION Protocol [Priority: LOW]
├─ 🚀 PERFORMANCE → Execute: PERFORMANCE Protocol [Priority: MEDIUM]
└─ [ACTION:migration] REFACTORING → Execute: REFACTORING Protocol [Priority: LOW]
```

### ✅ Step 3: Execute Protocol and Validate (EXECUTE THIRD)
```
For each protocol:
├─ [ACTION:format-processing] Update REQUIRED files (marked with ✅ HIGH PRIORITY)
├─ 📊 Evaluate CONDITIONAL files (marked with ⚠️ MEDIUM PRIORITY)
├─ ❌ Skip files marked as SKIP (NO ACTION)
├─ [ACTION:core-functionality] Add implementation tokens to code (MANDATORY)
├─ ✅ Run validation checks (MANDATORY)
└─ 🏁 Mark task complete in feature-tracking.md (MANDATORY)
```

## [ACTION:core-functionality] AI Assistant Change Protocols

### 🆕 NEW FEATURE Protocol [PRIORITY: CRITICAL]
**🚨 MANDATORY ACTIONS (Execute in Order):**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I implement this feature request?" with safety gates and scope boundaries validation
2. **📋 Feature ID Verification**: Must exist in feature-tracking.md with valid `FEATURE-001` format
3. **[ACTION:validation] Immutable Check**: Verify no conflicts with immutable.md requirements

**Phase 2: Documentation Updates [ACTION:format-processing]**
3. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Add feature entry with status "In Progress"
   - ✅ `specification.md` - Document user-facing behavior
   - ✅ `requirements.md` - Add implementation requirements
   - ✅ `architecture.md` - Document technical implementation
   - ✅ `testing.md` - Add test coverage requirements

4. **⚠️ MEDIUM PRIORITY - Evaluate CONDITIONAL files:**
   - ⚠️ `implementation-decisions.md` - IF making significant architectural decisions
   - ⚠️ `validation-automation.md` - IF adding new validation processes
   - ⚠️ `sync-framework.md` - IF feature affects data synchronization
   - ⚠️ `ai-assistant-compliance.md` - IF feature involves AI assistant guidance

**Phase 3: Implementation 💻**
5. **🏷️ Add implementation tokens:**
   ```go
   // FEATURE-ID: Brief description of implementation
   ```

**Phase 4: Quality Assurance ✅**
6. **🚫 NEVER modify these files:**
   - ❌ `immutable.md` - Only check for conflicts
   - ❌ `enforcement-mechanisms.md` - Reference only
   - ❌ Process files (`feature-change-protocol.md`, etc.) - Reference only

**🏁 COMPLETION CRITERIA (All Must Pass):**
- ✅ All tests pass
- ✅ All lint checks pass
- ✅ All required documentation updated
- ✅ Feature status set to "Completed" in `feature-tracking.md`

### [ACTION:core-functionality] MODIFICATION Protocol [PRIORITY: CRITICAL]
**🚨 MANDATORY ACTIONS (Execute in Order):**

**Phase 1: Feature Discovery [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I modify this existing functionality?" with safety gates and scope boundaries validation
2. **📋 Lookup existing Feature ID** in feature-tracking.md
3. **🆔 Create modification entry** - Add `-MOD-001` suffix to Feature ID

**Phase 2: Impact Analysis 📊**
3. **[ACTION:discovery] Impact Analysis:**
   - Find "Implementation Tokens" for the feature
   - Identify all files listed in the feature's registry entry
   - Update ALL files containing the modified feature

**Phase 3: Documentation Updates [ACTION:format-processing]**
4. **✅ HIGH PRIORITY - Update files based on impact:**
   - ✅ `feature-tracking.md` - Update feature entry with modification suffix
   - ⚠️ `specification.md` - IF user-facing behavior changes
   - ⚠️ `requirements.md` - IF implementation requirements change
   - ⚠️ `architecture.md` - IF technical implementation changes
   - ⚠️ `testing.md` - IF test requirements change
   - ⚠️ `implementation-decisions.md` - IF architectural decisions change

**Phase 4: Implementation 💻**
5. **🏷️ Update implementation tokens:**
   ```go
   // FEATURE-ID-MOD-001: Modification description
   ```

**🏁 COMPLETION CRITERIA:**
- ✅ All affected files updated per feature registry
- ✅ All tests pass
- ✅ All lint checks pass
- ✅ Modification documented in `feature-tracking.md`

### 🐛 BUG FIX Protocol [PRIORITY: MEDIUM]
**🚨 MINIMAL ACTIONS (Streamlined Process):**

**Phase 1: Impact Assessment [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I fix this test failure?" with safety gates and scope boundaries validation
2. **📊 Check if bug affects documented behavior**

**Phase 2: Selective Updates [ACTION:format-processing]**
2. **⚠️ Update files ONLY if necessary:**
   - ⚠️ `feature-tracking.md` - ONLY if fix affects documented behavior
   - ❌ SKIP: `specification.md` - Unless bug affects documented behavior
   - ❌ SKIP: `requirements.md` - Unless bug reveals requirement gap
   - ❌ SKIP: `architecture.md` - Unless architectural issue
   - ❌ SKIP: `testing.md` - Unless adding regression tests

**Phase 3: Implementation 💻**
3. **🏷️ Add implementation tokens** to code changes

**🏁 COMPLETION CRITERIA:**
- ✅ All tests pass
- ✅ All lint checks pass
- ✅ Minimal documentation updates only

### ⚙️ CONFIG CHANGE Protocol [PRIORITY: MEDIUM]
**🚨 MANDATORY ACTIONS:**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I make this configuration change?" with safety gates and scope boundaries validation

**Phase 2: Documentation Updates [ACTION:format-processing]**
2. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update configuration-related features
   - ✅ `specification.md` - IF user-visible configuration changes
   - ✅ `requirements.md` - Update configuration requirements

**Phase 3: Optional Updates 📋**
3. **❌ Usually SKIP:**
   - ❌ `architecture.md` - Unless configuration architecture changes
   - ❌ `testing.md` - Unless configuration testing changes

**🏁 COMPLETION CRITERIA:**
- ✅ Configuration changes documented
- ✅ All tests pass
- ✅ All lint checks pass

### 🔌 API CHANGE Protocol [PRIORITY: CRITICAL]
**🚨 MANDATORY ACTIONS:**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I make this API/interface change?" with safety gates and scope boundaries validation

**Phase 2: Core Documentation [ACTION:format-processing]**
2. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update API-related features
   - ✅ `specification.md` - Document API behavior changes
   - ✅ `architecture.md` - Update interface documentation

**Phase 3: Impact Evaluation 📊**
3. **⚠️ MEDIUM PRIORITY - Evaluate CONDITIONAL files:**
   - ⚠️ `sync-framework.md` - IF changes affect synchronization
   - ⚠️ `requirements.md` - IF API requirements change
   - ⚠️ `testing.md` - IF API testing requirements change

**🏁 COMPLETION CRITERIA:**
- ✅ API changes fully documented
- ✅ All interface changes reflected in architecture
- ✅ All tests pass
- ✅ All lint checks pass

### 🧪 TEST ADDITION Protocol [PRIORITY: LOW]
**🚨 MANDATORY ACTIONS:**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I add these tests?" with safety gates and scope boundaries validation

**Phase 2: Documentation Updates [ACTION:format-processing]**
2. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update test-related features
   - ✅ `testing.md` - Document new test coverage

**Phase 3: Usually Skip 📋**
3. **❌ Usually SKIP:**
   - ❌ `specification.md` - Unless tests reveal specification gaps
   - ❌ `requirements.md` - Unless tests reveal requirement gaps
   - ❌ `architecture.md` - Unless test infrastructure changes

**🏁 COMPLETION CRITERIA:**
- ✅ New tests documented
- ✅ All tests pass (including new ones)
- ✅ All lint checks pass

### 🚀 PERFORMANCE Protocol [PRIORITY: MEDIUM]
**🚨 MANDATORY ACTIONS:**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I make this performance optimization?" with safety gates and scope boundaries validation

**Phase 2: Core Documentation [ACTION:format-processing]**
2. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update performance-related features
   - ✅ `architecture.md` - Document performance improvements

**Phase 3: Impact Evaluation 📊**
3. **⚠️ MEDIUM PRIORITY - Evaluate CONDITIONAL files:**
   - ⚠️ `requirements.md` - IF performance requirements change
   - ⚠️ `testing.md` - IF performance testing requirements change
   - ⚠️ `specification.md` - IF user-visible performance characteristics change

**🏁 COMPLETION CRITERIA:**
- ✅ Performance improvements documented
- ✅ All tests pass

### [ACTION:migration] REFACTORING Protocol [PRIORITY: LOW]
**🚨 MANDATORY ACTIONS:**

**Phase 1: Pre-Implementation Validation [ACTION:discovery]**
1. **[CRITICAL] Decision Framework Validation**: Execute DOC-014 decision tree "Should I refactor this code?" with safety gates and scope boundaries validation

**Phase 2: Core Documentation [ACTION:format-processing]**
2. **✅ HIGH PRIORITY - Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update refactoring-related features
   - ✅ `architecture.md` - Document structural changes

**Phase 3: Usually Skip 📋**
3. **❌ Usually SKIP:**
   - ❌ `specification.md` - Unless user-facing behavior changes
   - ❌ `requirements.md` - Unless requirements are clarified
   - ❌ `testing.md` - Unless test structure changes

**🏁 COMPLETION CRITERIA:**
- ✅ Structural changes documented
- ✅ All tests pass (no regressions)
- ✅ All lint checks pass

## 🎯 AI Assistant Priority Matrix & Execution Order

### 🚨 CRITICAL PRIORITY [Execute IMMEDIATELY]
```
┌─ [ACTION:validation] IMMUTABLE CHECK
├─ 📋 TASK VERIFICATION  
├─ [ACTION:discovery] COMPLIANCE REVIEW
└─ 📁 IMPACT ANALYSIS
   ↓
PROCEED TO IMPLEMENTATION
```

### ⚡ HIGH PRIORITY [Execute FIRST in Implementation]
```
🆕 NEW FEATURE.......... CRITICAL - Full documentation cascade
[ACTION:core-functionality] MODIFY EXISTING...... CRITICAL - Impact analysis required  
🔌 API/INTERFACE........ CRITICAL - Interface documentation required
```

### 📊 MEDIUM PRIORITY [Execute with Evaluation]
```
🐛 BUG FIX.............. MEDIUM - Minimal documentation
⚙️ CONFIG CHANGE........ MEDIUM - Configuration documentation
🚀 PERFORMANCE.......... MEDIUM - Architecture documentation
```

### [ACTION:format-processing] LOW PRIORITY [Execute Last]
```
🧪 TEST ONLY............ LOW - Testing documentation only
[ACTION:migration] REFACTORING.......... LOW - Structural documentation only
```

## 🤖 AI Decision Tree for File Updates

### [ACTION:discovery] PRIMARY EVALUATION (Execute for ALL Changes)
```
Does this change affect:
├─ User-facing behavior? → UPDATE: specification.md
├─ Technical implementation? → UPDATE: architecture.md  
├─ New requirements? → UPDATE: requirements.md
├─ Test coverage? → UPDATE: testing.md
└─ Feature registry? → UPDATE: feature-tracking.md (ALWAYS)
```

### 📊 SECONDARY EVALUATION (Execute Conditionally)
```
Does this change involve:
├─ Architectural decisions? → EVALUATE: implementation-decisions.md
├─ Validation processes? → EVALUATE: validation-automation.md
├─ Data synchronization? → EVALUATE: sync-framework.md
└─ AI assistant guidance? → EVALUATE: ai-assistant-compliance.md
```

### ❌ SKIP EVALUATION (Never Update These)
```
Process files (Reference Only):
├─ immutable.md
├─ enforcement-mechanisms.md
├─ feature-change-protocol.md
└─ context-file-checklist.md
```

## 📊 AI Assistant File Priority Matrix

### 🔥 CRITICAL FILES (Always Check/Update)
| File | AI Assistant Action | When to Skip |
|------|---------------------|--------------|
| `feature-tracking.md` | **ALWAYS UPDATE** - Every change must be tracked | Never skip |
| `ai-assistant-compliance.md` | **ALWAYS CHECK** - Before any code change | Never skip checking |
| `immutable.md` | **ALWAYS CHECK** - Before any change | Never modify |

### 📋 HIGH-PRIORITY FILES (Update Based on Change Type)
| File | Update When | AI Skip When |
|------|-------------|--------------|
| `specification.md` | User-facing behavior changes | Internal-only changes |
| `requirements.md` | New/changed requirements | Implementation-only changes |
| `architecture.md` | Technical implementation changes | UI-only changes |
| `testing.md` | New test requirements | Bug fixes without new tests |

### 📄 MEDIUM-PRIORITY FILES (Evaluate Conditionally)
| File | Update When | AI Skip When |
|------|-------------|--------------|
| `implementation-decisions.md` | Significant architectural decisions | Minor implementation details |
| `validation-automation.md` | New validation processes | Standard feature additions |
| `sync-framework.md` | Synchronization logic changes | Unrelated features |

### 📚 LOW-PRIORITY FILES (Rarely Update)
| File | Update When | AI Skip When |
|------|-------------|--------------|
| `enforcement-mechanisms.md` | Process changes only | Feature development |
| `feature-change-protocol.md` | Process improvements | Individual features |
| `context-file-checklist.md` | Process improvements | Individual features |

## 🤖 AI Assistant Implementation Checklist

### [ACTION:validation] Pre-Change Validation [MANDATORY - Execute FIRST]
```bash
# MANDATORY: Execute before any code modification
echo "🤖 AI Assistant Pre-Change Validation"

# 1. Execute Decision Framework validation (DOC-014)
echo "[CRITICAL] Decision Framework: Selecting appropriate decision tree"
echo "[CRITICAL] Decision Framework: Validating safety gates and scope boundaries"

# 2. Verify task exists in feature tracking
grep -q "FEATURE-ID" docs/context/feature-tracking.md || echo "❌ CRITICAL: Task not found in feature-tracking.md"

# 3. Check ai-assistant-compliance requirements
echo "📋 Review ai-assistant-compliance.md for token requirements"

# 4. Check immutable requirements
echo "[ACTION:validation] Verify no immutable.md conflicts"

# 5. Determine file impact using this guide
echo "📁 Files requiring updates determined"
```

### [ACTION:core-functionality] During Implementation [MANDATORY - Execute SECOND]
```bash
# MANDATORY: Execute during code changes
echo "[ACTION:core-functionality] AI Assistant Implementation Phase"

# 1. Add implementation tokens to ALL modified code
echo "🏷️ Implementation tokens added"

# 2. Update documentation incrementally
echo "[ACTION:format-processing] Documentation updates in progress"

# 3. Maintain cross-references
echo "🔗 Cross-references maintained"
```

### ✅ Post-Change Validation [MANDATORY - Execute LAST]
```bash
# MANDATORY: Execute after all changes complete
echo "✅ AI Assistant Post-Change Validation"

# 1. Execute Decision Framework post-implementation validation (DOC-014)
echo "[CRITICAL] Decision Framework: Validating quality thresholds and goal alignment"
echo "[CRITICAL] Decision Framework: Confirming >95% goal alignment rate and 100% traceability"

# 2. Run full test suite
make test || echo "❌ CRITICAL: Tests failing"

# 3. Run lint checks
make lint || echo "❌ CRITICAL: Lint failures"

# 4. Validate documentation consistency
./scripts/validate-docs.sh || echo "⚠️ WARNING: Documentation inconsistencies"

# 5. Update task status in feature-tracking.md
echo "🏁 Task marked complete in feature-tracking.md"
```

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

### ✅ After Code Changes [PRIORITY 3]
1. **🧪 Full Test Validation**: All tests must pass
2. **[ACTION:core-functionality] Lint Compliance**: All lint checks must pass  
3. **📊 Documentation Validation**: Run validation scripts
4. **🏁 Task Completion**: Mark task complete in `feature-tracking.md`

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

## 🚨 AI Assistant Mandatory Compliance

### [ACTION:validation] Before ANY Code Change [PRIORITY 1]
1. **📋 Task Verification**: Feature ID must exist in `feature-tracking.md`
2. **[ACTION:discovery] Compliance Review**: Check `ai-assistant-compliance.md` for specific requirements
3. **📊 Impact Analysis**: Use appropriate protocol from this guide
4. **[ACTION:validation] Immutable Check**: Verify no conflicts with `immutable.md`

### [ACTION:core-functionality] During Code Changes [PRIORITY 2]
1. **🏷️ Implementation Tokens**: Add to ALL modified code sections
2. **[ACTION:format-processing] Incremental Documentation**: Update documentation as changes are made
3. **🔗 Cross-Reference Maintenance**: Keep links between documents valid

### ✅ After Code Changes [PRIORITY 3]
1. **🧪 Full Test Validation**: All tests must pass
2. **[ACTION:core-functionality] Lint Compliance**: All lint checks must pass  
3. **📊 Documentation Validation**: Run validation scripts
4. **🏁 Task Completion**: Mark task complete in `feature-tracking.md`

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

### 🎯 Validation Commands
```bash
# Pre-change validation
grep "FEATURE-ID" docs/context/feature-tracking.md
cat docs/context/ai-assistant-compliance.md

# Post-change validation  
make test
make lint
./scripts/validate-docs.sh

# Task completion
# Update feature-tracking.md status to "Completed"
```

## [ACTION:core-functionality] AI Assistant Quick Reference

### 🆔 Feature ID Format
```
PREFIX-###
Where PREFIX = ARCH|CMD|CFG|API|TEST|GIT|FILE|OUT
```

### 🏷️ Implementation Token Format
```go
// FEATURE-ID: Brief description
// FEATURE-ID-MOD-001: Modification description (for changes)
```

### 📊 Required Status Updates in feature-tracking.md
```
Status Options:
- "Not Started" → "In Progress" → "Completed"
- "Not Started" → "BLOCKED" (if dependencies missing)
```

### ✅ Critical Validation Gates [Execute in Order]
1. ✅ Task exists in feature-tracking.md
2. ✅ All required files updated per protocol
3. ✅ Implementation tokens added to code
4. ✅ All tests pass
5. ✅ All lint checks pass
6. ✅ Task marked "Completed" in feature-tracking.md

**🚨 FAILURE TO FOLLOW THIS PROTOCOL INVALIDATES ANY CODE CHANGES**
