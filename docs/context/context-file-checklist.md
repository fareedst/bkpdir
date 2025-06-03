# 📋 Context File Checklist: AI Assistant Guide

## 🎯 Purpose
This checklist helps AI assistants determine which context documentation files need updates when making code changes.

> **🤖 For AI Assistants**: Use this checklist AFTER identifying your change type in [ai-assistant-protocol.md](ai-assistant-protocol.md) and BEFORE implementing changes.

## 🚨 CRITICAL FILES [AI Must Always Check]

### 🛡️ immutable.md [MANDATORY CHECK - NEVER MODIFY]
**When to Check**: Before ANY code change
**AI Action**: ✅ Always verify no conflicts
**Purpose**: Ensure changes don't violate unchangeable specifications
```
✅ ALWAYS CHECK: Verify no conflicts with core behaviors
❌ NEVER MODIFY: This file is read-only for validation
🔍 SEARCH FOR: Your change area to identify potential conflicts
```

### 📋 feature-tracking.md [MANDATORY UPDATE]
**When to Update**: Every single code change
**AI Action**: ✅ Always update
**Purpose**: Track all changes and feature implementations
```
✅ ALWAYS UPDATE: Add or update Feature ID entries
🆕 NEW FEATURES: Create new feature entry with status "In Progress"
🔧 MODIFICATIONS: Update existing entry with modification suffix
🏁 COMPLETION: Change status to "Completed" when done
```

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

### 🔍 ai-assistant-compliance.md [MANDATORY CHECK]
**When to Check**: Before ANY code change
**AI Action**: ✅ Always review for token requirements
**Purpose**: Ensure compliance with token referencing rules
```
✅ ALWAYS CHECK: Review token requirements before changes
🔍 REFERENCE: Follow response template format
🏷️ TOKENS: Ensure proper implementation token usage
```

## 🎯 HIGH-PRIORITY FILES [AI Update Based on Change Type]

### 📖 specification.md [UPDATE: User-Facing Changes]
**When to Update**: 
- 🆕 New features that users interact with
- 🔧 Changes to existing user-facing behavior
- 🔌 API/interface changes that affect users
- ⚙️ Configuration changes visible to users

**AI Action**: ✅ Update if user-facing behavior changes
**Skip When**: Internal-only implementation changes
```
✅ UPDATE FOR: New user features, behavior changes, visible configuration
❌ SKIP FOR: Internal refactoring, bug fixes without behavior change
📝 FOCUS ON: User-visible functionality and configuration options
```

### 📝 requirements.md [UPDATE: Implementation Requirements]
**When to Update**:
- 🆕 New features requiring new implementation patterns
- 🔧 Changes that alter implementation requirements
- 🔌 API changes that require new integration patterns
- 🚀 Performance requirements or constraints

**AI Action**: ✅ Update if implementation requirements change
**Skip When**: Code changes that don't alter requirements
```
✅ UPDATE FOR: New requirements, changed constraints, integration patterns
❌ SKIP FOR: Implementation-only changes within existing requirements
📝 FOCUS ON: Technical constraints, integration requirements, patterns
```

### 🏗️ architecture.md [UPDATE: Technical Implementation]
**When to Update**:
- 🆕 New features requiring new components
- 🔧 Changes to existing component interfaces
- 🔌 API/interface architecture changes
- 🚀 Performance-related architectural changes
- 🔄 Refactoring that changes component structure

**AI Action**: ✅ Update if technical architecture changes
**Skip When**: Internal implementation without architectural impact
```
✅ UPDATE FOR: Component changes, interface modifications, architectural decisions
❌ SKIP FOR: Internal implementation details within existing architecture
📝 FOCUS ON: Component interfaces, system design, architectural patterns
```

### 🧪 testing.md [UPDATE: Test Requirements]
**When to Update**:
- 🆕 New features requiring new test coverage
- 🔧 Changes that require additional test validation
- 🧪 Test-only changes adding new test requirements
- 🔌 API changes requiring interface testing

**AI Action**: ✅ Update if test coverage requirements change
**Skip When**: Changes that don't require new testing approaches
```
✅ UPDATE FOR: New test coverage areas, testing requirements, validation needs
❌ SKIP FOR: Changes covered by existing test patterns
📝 FOCUS ON: Test coverage requirements, validation strategies, test patterns
```

## 📊 MEDIUM-PRIORITY FILES [AI Evaluate Conditionally]

### ⚙️ implementation-decisions.md [CONDITIONAL UPDATE]
**When to Update**: 
- Making significant architectural decisions
- Choosing between multiple implementation approaches
- Establishing new patterns or conventions

**AI Action**: ⚠️ Update only for significant architectural decisions
```
⚠️ UPDATE IF: Significant architectural choices, pattern establishment
❌ SKIP IF: Standard implementation following existing patterns
📝 FOCUS ON: Decision rationale, alternative approaches considered
```

### 🔧 validation-automation.md [CONDITIONAL UPDATE]
**When to Update**:
- Adding new validation processes
- Modifying existing validation logic
- Introducing new quality assurance mechanisms

**AI Action**: ⚠️ Update only for validation process changes
```
⚠️ UPDATE IF: New validation processes, quality assurance changes
❌ SKIP IF: Standard feature development
📝 FOCUS ON: Validation processes, quality assurance mechanisms
```

### 🔄 sync-framework.md [CONDITIONAL UPDATE]
**When to Update**:
- Changes affecting data synchronization
- Cross-system integration modifications
- Document synchronization process changes

**AI Action**: ⚠️ Update only for synchronization-related changes
```
⚠️ UPDATE IF: Data sync changes, cross-system integration modifications
❌ SKIP IF: Changes unrelated to synchronization
📝 FOCUS ON: Synchronization processes, data consistency mechanisms
```

## ❌ NEVER MODIFY [AI Reference Only]

### 🔒 Process Documentation Files
These files define processes and should only be referenced, never modified during feature development:

- **🔒 enforcement-mechanisms.md** - Validation rules and enforcement
- **🚫 change-rejection-criteria.md** - Common rejection scenarios
- **📚 feature-documentation-standards.md** - Documentation guidelines
- **🎯 feature-change-protocol.md** - Change management protocols

**AI Action**: 📚 Reference only for understanding processes
```
📚 REFERENCE ONLY: Use for understanding processes and guidelines
❌ NEVER MODIFY: These define the development process itself
🔍 USE FOR: Understanding requirements and avoiding common mistakes
```

## 🎯 AI Assistant Quick Decision Tree

```
🤖 AI CHANGE DECISION TREE:

1. 🔍 What type of change am I making?
   ├─ 🆕 NEW FEATURE → Update: specification.md, requirements.md, architecture.md, testing.md
   ├─ 🔧 MODIFY EXISTING → Evaluate: specification.md*, requirements.md*, architecture.md*, testing.md*
   ├─ 🐛 BUG FIX → Minimal: Usually skip most documentation
   ├─ ⚙️ CONFIG CHANGE → Focus: specification.md, requirements.md
   ├─ 🔌 API/INTERFACE → Update: specification.md, architecture.md
   ├─ 🧪 TEST ONLY → Update: testing.md
   ├─ 🚀 PERFORMANCE → Update: architecture.md
   └─ 🔄 REFACTORING → Update: architecture.md

2. 📊 Does this involve special considerations?
   ├─ Architectural decisions → Consider: implementation-decisions.md
   ├─ Validation processes → Consider: validation-automation.md
   └─ Data synchronization → Consider: sync-framework.md

3. ✅ ALWAYS UPDATE:
   ├─ feature-tracking.md (Feature registry)
   └─ Add implementation tokens to code

4. 🛡️ ALWAYS CHECK:
   ├─ immutable.md (Conflict verification)
   └─ ai-assistant-compliance.md (Token requirements)

* = Update only if the change affects that area
```

## ✅ AI Assistant Final Checklist

Before marking any task complete:

- [ ] 🛡️ Checked immutable.md for conflicts
- [ ] 📋 Updated feature-tracking.md with Feature ID
- [ ] 🔍 Followed ai-assistant-compliance.md requirements
- [ ] 🏷️ Added implementation tokens to modified code
- [ ] 📝 Updated all required context files per change type
- [ ] 🧪 All tests pass (`make test`)
- [ ] 🔧 All lint checks pass (`make lint`)
- [ ] 🏁 Marked Feature ID as "Completed" in feature-tracking.md
- [ ] 🚨 **CRITICAL**: Updated BOTH feature registry table AND detailed subtask blocks with completion status

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

**🤖 Remember: Better to update too much documentation than too little. When in doubt, include the update!**
