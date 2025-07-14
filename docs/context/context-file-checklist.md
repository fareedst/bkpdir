# [CHECKLIST] Context File Checklist: AI Assistant Guide

## [PURPOSE] Purpose
This checklist helps AI assistants determine which context documentation files need updates when making code changes.

> **[AI_ASSISTANT] For AI Assistants**: Use this checklist AFTER identifying your change type in [ai-assistant-protocol.md](ai-assistant-protocol.md) and BEFORE implementing changes.

## [ALERT] CRITICAL FILES [AI Must Always Check]

### [ACTION:validation] immutable.md [MANDATORY CHECK - NEVER MODIFY]
**When to Check**: Before ANY code change
**AI Action**: ✅ Always verify no conflicts
**Purpose**: Ensure changes don't violate unchangeable specifications
```
✅ ALWAYS CHECK: Verify no conflicts with core behaviors
❌ NEVER MODIFY: This file is read-only for validation
[ACTION:discovery] SEARCH FOR: Your change area to identify potential conflicts
```

### 📋 feature-tracking.md [MANDATORY UPDATE]
**When to Update**: Every single code change
**AI Action**: ✅ Always update
**Purpose**: Track all changes and feature implementations
```
✅ ALWAYS UPDATE: Add or update Feature ID entries
🆕 NEW FEATURES: Create new feature entry with status "In Progress"
[ACTION:core-functionality] MODIFICATIONS: Update existing entry with modification suffix
🏁 COMPLETION: Change status to "Completed" when done
```

> **🚨 CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

### [ACTION:discovery] ai-assistant-compliance.md [MANDATORY CHECK]
**When to Check**: Before ANY code change
**AI Action**: ✅ Always review for token requirements
**Purpose**: Ensure compliance with token referencing rules
```
✅ ALWAYS CHECK: Review token requirements before changes
[ACTION:discovery] REFERENCE: Follow response template format
🏷️ TOKENS: Ensure proper implementation token usage
```

## [HIGH_PRIORITY] HIGH-PRIORITY FILES [AI Update Based on Change Type]

### [TECHNICAL] specification.md [UPDATE: User-Facing Changes]
**When to Update**: 
- 🆕 New features that users interact with
- [ACTION:core-functionality] Changes to existing user-facing behavior
- 🔌 API/interface changes that affect users
- ⚙️ Configuration changes visible to users

**AI Action**: ✅ Update if user-facing behavior changes
**Skip When**: Internal-only implementation changes
```
✅ UPDATE FOR: New user features, behavior changes, visible configuration
❌ SKIP FOR: Internal refactoring, bug fixes without behavior change
[ACTION:format-processing] FOCUS ON: User-visible functionality and configuration options
```

### [DOCUMENT_UPDATE] requirements.md [UPDATE: Implementation Requirements]
**When to Update**:
- [NEW_FEATURE] New features requiring new implementation patterns
- [CONFIGURE_MODIFY] Changes that alter implementation requirements
- [INTEGRATION] API changes that require new integration patterns
- [EXECUTION] Performance requirements or constraints

**AI Action**: ✅ Update if implementation requirements change
**Skip When**: Code changes that don't alter requirements
```
✅ UPDATE FOR: New requirements, changed constraints, integration patterns
❌ SKIP FOR: Implementation-only changes within existing requirements
[DOCUMENT_UPDATE] FOCUS ON: Technical constraints, integration requirements, patterns
```

### [BUILD_CONSTRUCT] architecture.md [UPDATE: Technical Implementation]
**When to Update**:
- [NEW_FEATURE] New features requiring new components
- [CONFIGURE_MODIFY] Changes to existing component interfaces
- [INTEGRATION] API/interface architecture changes
- [EXECUTION] Performance-related architectural changes
- [PROCESS] Refactoring that changes component structure

**AI Action**: ✅ Update if technical architecture changes
**Skip When**: Internal implementation without architectural impact
```
✅ UPDATE FOR: Component changes, interface modifications, architectural decisions
❌ SKIP FOR: Internal implementation details within existing architecture
[DOCUMENT_UPDATE] FOCUS ON: Component interfaces, system design, architectural patterns
```

### [VALIDATION] testing.md [UPDATE: Test Requirements]
**When to Update**:
- [NEW_FEATURE] New features requiring new test coverage
- [CONFIGURE_MODIFY] Changes that require additional test validation
- [VALIDATION] Test-only changes adding new test requirements
- [INTEGRATION] API changes requiring interface testing

**AI Action**: ✅ Update if test coverage requirements change
**Skip When**: Changes that don't require new testing approaches
```
✅ UPDATE FOR: New test coverage areas, testing requirements, validation needs
❌ SKIP FOR: Changes covered by existing test patterns
[DOCUMENT_UPDATE] FOCUS ON: Test coverage requirements, validation strategies, test patterns
```

## [METRICS] MEDIUM-PRIORITY FILES [AI Evaluate Conditionally]

### ⚙️ implementation-decisions.md [CONDITIONAL UPDATE]
**When to Update**: 
- Making significant architectural decisions
- Choosing between multiple implementation approaches
- Establishing new patterns or conventions

**AI Action**: ⚠️ Update only for significant architectural decisions
```
⚠️ UPDATE IF: Significant architectural choices, pattern establishment
❌ SKIP IF: Standard implementation following existing patterns
[ACTION:format-processing] FOCUS ON: Decision rationale, alternative approaches considered
```

### [ACTION:core-functionality] validation-automation.md [CONDITIONAL UPDATE]
**When to Update**:
- Adding new validation processes
- Modifying existing validation logic
- Introducing new quality assurance mechanisms

**AI Action**: ⚠️ Update only for validation process changes
```
⚠️ UPDATE IF: New validation processes, quality assurance changes
❌ SKIP IF: Standard feature development
[ACTION:format-processing] FOCUS ON: Validation processes, quality assurance mechanisms
```

### [ACTION:migration] sync-framework.md [CONDITIONAL UPDATE]
**When to Update**:
- Changes affecting data synchronization
- Cross-system integration modifications
- Document synchronization process changes

**AI Action**: ⚠️ Update only for synchronization-related changes
```
⚠️ UPDATE IF: Data sync changes, cross-system integration modifications
❌ SKIP IF: Changes unrelated to synchronization
[ACTION:format-processing] FOCUS ON: Synchronization processes, data consistency mechanisms
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
[ACTION:discovery] USE FOR: Understanding requirements and avoiding common mistakes
```

## [OBJECTIVE] AI Assistant Quick Decision Tree

```
[AI_ASSISTANT] AI CHANGE DECISION TREE:

1. [SEARCH_DISCOVER] What type of change am I making?
   ├─ [NEW_FEATURE] NEW FEATURE → Update: specification.md, requirements.md, architecture.md, testing.md
   ├─ [CONFIGURE_MODIFY] MODIFY EXISTING → Evaluate: specification.md*, requirements.md*, architecture.md*, testing.md*
   ├─ [VALIDATION] BUG FIX → Minimal: Usually skip most documentation
   ├─ [CONFIGURATION] CONFIG CHANGE → Focus: specification.md, requirements.md
   ├─ [INTEGRATION] API/INTERFACE → Update: specification.md, architecture.md
   ├─ [VALIDATION] TEST ONLY → Update: testing.md
   ├─ [EXECUTION] PERFORMANCE → Update: architecture.md
   └─ [PROCESS] REFACTORING → Update: architecture.md

2. [METRICS] Does this involve special considerations?
   ├─ Architectural decisions → Consider: implementation-decisions.md
   ├─ Validation processes → Consider: validation-automation.md
   └─ Data synchronization → Consider: sync-framework.md

3. [COMPLETED] ALWAYS UPDATE:
   ├─ feature-tracking.md (Feature registry)
   └─ Add implementation tokens to code

4. [VALIDATION] ALWAYS CHECK:
   ├─ immutable.md (Conflict verification)
   └─ ai-assistant-compliance.md (Token requirements)

* = Update only if the change affects that area
```

## [COMPLETED] AI Assistant Final Checklist

Before marking any task complete:

- [ ] [VALIDATION] Checked immutable.md for conflicts
- [ ] [CHECKLIST] Updated feature-tracking.md with Feature ID
- [ ] [SEARCH_DISCOVER] Followed ai-assistant-compliance.md requirements
- [ ] [CLASSIFICATION] Added implementation tokens to modified code
- [ ] [DOCUMENT_UPDATE] Updated all required context files per change type
- [ ] [VALIDATION] All tests pass (`make test`)
- [ ] [CONFIGURE_MODIFY] All lint checks pass (`make lint`)
- [ ] [ACHIEVEMENT] Marked Feature ID as "Completed" in feature-tracking.md
- [ ] [ALERT] **CRITICAL**: Updated BOTH feature registry table AND detailed subtask blocks with completion status

> **[ALERT] CRITICAL NOTE FOR AI ASSISTANTS**: When marking a task as completed in the feature registry table, you MUST also update the detailed subtask blocks to show all subtasks as completed with checkmarks [x]. Failure to update both locations creates documentation inconsistency and violates DOC-008 enforcement requirements.

**[AI_ASSISTANT] Remember: Better to update too much documentation than too little. When in doubt, include the update!**
