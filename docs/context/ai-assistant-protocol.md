# AI Assistant Feature Update Protocol

> **🤖 MANDATORY FOR ALL AI ASSISTANTS**: This document defines the REQUIRED protocol for all code changes. Failure to follow this protocol invalidates any code modifications.

## 🚨 CRITICAL COMPLIANCE REQUIREMENTS

### **MANDATORY PRE-WORK VALIDATION**
Before making ANY code changes:

1. **📋 Task Verification**: The task MUST exist in `feature-tracking.md` with a valid Feature ID
2. **🔍 Compliance Check**: Review `ai-assistant-compliance.md` for mandatory token referencing requirements
3. **📁 File Impact Analysis**: Use this guide to determine which documentation files require updates
4. **🛡️ Immutable Check**: Verify no changes violate requirements in `immutable.md`

### **MANDATORY POST-WORK COMPLETION**
After ALL code changes are complete:

1. **✅ Full Test Suite**: All tests must pass (`make test`)
2. **🔧 Lint Compliance**: All lint checks must pass (`make lint`)
3. **📝 Documentation Updates**: All required documentation files updated per this protocol
4. **🏁 Task Completion**: Update task status to "Completed" in `feature-tracking.md`

## 🤖 AI Assistant Decision Engine

### Step 1: Task Identification and Validation
```
REQUIRED: Before any code modification
├─ 📋 Does Feature ID exist in feature-tracking.md? → YES: Continue | NO: STOP - Create task first
├─ 🔍 Are implementation tokens required? → Check ai-assistant-compliance.md
├─ 📁 Which files need documentation updates? → Use workflows below
└─ 🛡️ Any immutable requirement conflicts? → Check immutable.md
```

### Step 2: Change Type Classification
```
What type of change are you making?
├─ 🆕 NEW FEATURE → Execute: NEW FEATURE Protocol
├─ 🔧 MODIFY EXISTING → Execute: MODIFICATION Protocol  
├─ 🐛 BUG FIX → Execute: BUG FIX Protocol
├─ ⚙️ CONFIG CHANGE → Execute: CONFIG CHANGE Protocol
├─ 🔌 API/INTERFACE → Execute: API CHANGE Protocol
├─ 🧪 TEST ONLY → Execute: TEST ADDITION Protocol
├─ 🚀 PERFORMANCE → Execute: PERFORMANCE Protocol
└─ 🔄 REFACTORING → Execute: REFACTORING Protocol
```

### Step 3: Execute Protocol and Validate
```
For each protocol:
├─ 📝 Update REQUIRED files (marked with ✅)
├─ 📊 Evaluate CONDITIONAL files (marked with ⚠️)
├─ ❌ Skip files marked as SKIP
├─ 🔧 Add implementation tokens to code
├─ ✅ Run validation checks
└─ 🏁 Mark task complete in feature-tracking.md
```

## 🔧 AI Assistant Change Protocols

### 🆕 NEW FEATURE Protocol
**🚨 MANDATORY ACTIONS:**
1. **Verify Feature ID exists in feature-tracking.md** - Must have valid `FEATURE-001` format
2. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Add feature entry with status "In Progress"
   - ✅ `specification.md` - Document user-facing behavior
   - ✅ `requirements.md` - Add implementation requirements
   - ✅ `architecture.md` - Document technical implementation
   - ✅ `testing.md` - Add test coverage requirements

3. **Evaluate CONDITIONAL files:**
   - ⚠️ `implementation-decisions.md` - IF making significant architectural decisions
   - ⚠️ `validation-automation.md` - IF adding new validation processes
   - ⚠️ `sync-framework.md` - IF feature affects data synchronization
   - ⚠️ `ai-assistant-compliance.md` - IF feature involves AI assistant guidance

4. **Add implementation tokens:**
   ```go
   // FEATURE-ID: Brief description of implementation
   ```

5. **NEVER modify these files:**
   - ❌ `immutable.md` - Only check for conflicts
   - ❌ `enforcement-mechanisms.md` - Reference only
   - ❌ Process files (`feature-change-protocol.md`, etc.) - Reference only

**🏁 COMPLETION CRITERIA:**
- All tests pass
- All lint checks pass
- All required documentation updated
- Feature status set to "Completed" in `feature-tracking.md`

### 🔧 MODIFICATION Protocol
**🚨 MANDATORY ACTIONS:**
1. **Lookup existing Feature ID in feature-tracking.md**
2. **Create modification entry** - Add `-MOD-001` suffix to Feature ID
3. **Impact Analysis:**
   - Find "Implementation Tokens" for the feature
   - Identify all files listed in the feature's registry entry
   - Update ALL files containing the modified feature

4. **Update files based on impact:**
   - ✅ `feature-tracking.md` - Update feature entry with modification suffix
   - ⚠️ `specification.md` - IF user-facing behavior changes
   - ⚠️ `requirements.md` - IF implementation requirements change
   - ⚠️ `architecture.md` - IF technical implementation changes
   - ⚠️ `testing.md` - IF test requirements change
   - ⚠️ `implementation-decisions.md` - IF architectural decisions change

5. **Update implementation tokens:**
   ```go
   // FEATURE-ID-MOD-001: Modification description
   ```

**🏁 COMPLETION CRITERIA:**
- All affected files updated per feature registry
- All tests pass
- All lint checks pass
- Modification documented in `feature-tracking.md`

### 🐛 BUG FIX Protocol
**🚨 MINIMAL ACTIONS (Most changes skip documentation):**
1. **Check if bug affects documented behavior**
2. **Update files ONLY if necessary:**
   - ⚠️ `feature-tracking.md` - ONLY if fix affects documented behavior
   - ❌ SKIP: `specification.md` - Unless bug affects documented behavior
   - ❌ SKIP: `requirements.md` - Unless bug reveals requirement gap
   - ❌ SKIP: `architecture.md` - Unless architectural issue
   - ❌ SKIP: `testing.md` - Unless adding regression tests

3. **Add implementation tokens to code changes**

**🏁 COMPLETION CRITERIA:**
- All tests pass
- All lint checks pass
- Minimal documentation updates only

### ⚙️ CONFIG CHANGE Protocol
**🚨 MANDATORY ACTIONS:**
1. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update configuration-related features
   - ✅ `specification.md` - IF user-visible configuration changes
   - ✅ `requirements.md` - Update configuration requirements

2. **Usually SKIP:**
   - ❌ `architecture.md` - Unless configuration architecture changes
   - ❌ `testing.md` - Unless configuration testing changes

**🏁 COMPLETION CRITERIA:**
- Configuration changes documented
- All tests pass
- All lint checks pass

### 🔌 API CHANGE Protocol
**🚨 MANDATORY ACTIONS:**
1. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update API-related features
   - ✅ `specification.md` - Document API behavior changes
   - ✅ `architecture.md` - Update interface documentation

2. **Evaluate CONDITIONAL files:**
   - ⚠️ `sync-framework.md` - IF changes affect synchronization
   - ⚠️ `requirements.md` - IF API requirements change
   - ⚠️ `testing.md` - IF API testing requirements change

**🏁 COMPLETION CRITERIA:**
- API changes fully documented
- All interface changes reflected in architecture
- All tests pass
- All lint checks pass

### 🧪 TEST ADDITION Protocol
**🚨 MANDATORY ACTIONS:**
1. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update test-related features
   - ✅ `testing.md` - Document new test coverage

2. **Usually SKIP:**
   - ❌ `specification.md` - Unless tests reveal specification gaps
   - ❌ `requirements.md` - Unless tests reveal requirement gaps
   - ❌ `architecture.md` - Unless test infrastructure changes

**🏁 COMPLETION CRITERIA:**
- New tests documented
- All tests pass (including new ones)
- All lint checks pass

### 🚀 PERFORMANCE Protocol
**🚨 MANDATORY ACTIONS:**
1. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update performance-related features
   - ✅ `architecture.md` - Document performance improvements

2. **Evaluate CONDITIONAL files:**
   - ⚠️ `requirements.md` - IF performance requirements change
   - ⚠️ `testing.md` - IF performance testing requirements change
   - ⚠️ `specification.md` - IF user-visible performance characteristics change

**🏁 COMPLETION CRITERIA:**
- Performance improvements documented
- All tests pass
- All lint checks pass

### 🔄 REFACTORING Protocol
**🚨 MANDATORY ACTIONS:**
1. **Update REQUIRED files:**
   - ✅ `feature-tracking.md` - Update refactoring-related features
   - ✅ `architecture.md` - Document structural changes

2. **Usually SKIP:**
   - ❌ `specification.md` - Unless user-facing behavior changes
   - ❌ `requirements.md` - Unless requirements are clarified
   - ❌ `testing.md` - Unless test structure changes

**🏁 COMPLETION CRITERIA:**
- Structural changes documented
- All tests pass (no regressions)
- All lint checks pass

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

### Pre-Change Validation
```bash
# MANDATORY: Execute before any code modification
echo "🤖 AI Assistant Pre-Change Validation"

# 1. Verify task exists in feature tracking
grep -q "FEATURE-ID" docs/context/feature-tracking.md || echo "❌ CRITICAL: Task not found in feature-tracking.md"

# 2. Check ai-assistant-compliance requirements
echo "📋 Review ai-assistant-compliance.md for token requirements"

# 3. Check immutable requirements
echo "🛡️ Verify no immutable.md conflicts"

# 4. Determine file impact using this guide
echo "📁 Files requiring updates determined"
```

### During Implementation
```bash
# MANDATORY: Execute during code changes
echo "🔧 AI Assistant Implementation Phase"

# 1. Add implementation tokens to ALL modified code
echo "🏷️ Implementation tokens added"

# 2. Update documentation incrementally
echo "📝 Documentation updates in progress"

# 3. Maintain cross-references
echo "🔗 Cross-references maintained"
```

### Post-Change Validation
```bash
# MANDATORY: Execute after all changes complete
echo "✅ AI Assistant Post-Change Validation"

# 1. Run full test suite
make test || echo "❌ CRITICAL: Tests failing"

# 2. Run lint checks
make lint || echo "❌ CRITICAL: Lint failures"

# 3. Validate documentation consistency
./scripts/validate-docs.sh || echo "⚠️ WARNING: Documentation inconsistencies"

# 4. Update task status in feature-tracking.md
echo "🏁 Task marked complete in feature-tracking.md"
```

## 🚨 AI Assistant Mandatory Compliance

### Before ANY Code Change
1. **Task Verification**: Feature ID must exist in `feature-tracking.md`
2. **Compliance Review**: Check `ai-assistant-compliance.md` for specific requirements
3. **Impact Analysis**: Use appropriate protocol from this guide
4. **Immutable Check**: Verify no conflicts with `immutable.md`

### During Code Changes
1. **Implementation Tokens**: Add to ALL modified code sections
2. **Incremental Documentation**: Update documentation as changes are made
3. **Cross-Reference Maintenance**: Keep links between documents valid

### After Code Changes
1. **Full Test Validation**: All tests must pass
2. **Lint Compliance**: All lint checks must pass  
3. **Documentation Validation**: Run validation scripts
4. **Task Completion**: Mark task complete in `feature-tracking.md`

### Validation Commands
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

## 🔧 AI Assistant Quick Reference

### Feature ID Format
```
PREFIX-###
Where PREFIX = ARCH|CMD|CFG|API|TEST|GIT|FILE|OUT
```

### Implementation Token Format
```go
// FEATURE-ID: Brief description
// FEATURE-ID-MOD-001: Modification description (for changes)
```

### Required Status Updates in feature-tracking.md
```
Status Options:
- "Not Started" → "In Progress" → "Completed"
- "Not Started" → "BLOCKED" (if dependencies missing)
```

### Critical Validation Gates
1. ✅ Task exists in feature-tracking.md
2. ✅ All required files updated per protocol
3. ✅ Implementation tokens added to code
4. ✅ All tests pass
5. ✅ All lint checks pass
6. ✅ Task marked "Completed" in feature-tracking.md

**🚨 FAILURE TO FOLLOW THIS PROTOCOL INVALIDATES ANY CODE CHANGES**
