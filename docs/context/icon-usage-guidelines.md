# Icon Usage Guidelines - Maintaining the Standardized System

## 📑 Purpose
This document establishes the official guidelines for using icons in the context documentation system. It ensures the standardized icon system remains conflict-free and provides clear rules for future documentation updates.

## [ACTION:validation] Core Principles

### ✨ **Fundamental Rules**
1. **One Icon, One Meaning** - Each icon represents exactly one semantic concept across ALL documentation
2. **Meaning Consistency** - The same icon MUST have the same meaning in every document
3. **Visual Hierarchy** - Icon appearance reflects importance and priority levels
4. **Semantic Clarity** - Icon visual appearance suggests its meaning
5. **AI-Friendly Design** - Unambiguous interpretation for AI assistants

### 🚫 **Prohibited Practices** 
- ❌ **Reusing Priority Icons** - Never use [CRITICAL][HIGH][MEDIUM][LOW] for non-priority meanings
- ❌ **Phase/Step Confusion** - Don't mix 🚀⚡[ACTION:migration]🏁 with 1️⃣2️⃣3️⃣✅
- ❌ **Navigation Conflicts** - Keep 📑📋📊📖 separate from actions
- ❌ **Context-Dependent Meanings** - Same icon must mean same thing everywhere

## 🔥 Priority Hierarchy Guidelines

### [CRITICAL] **CRITICAL PRIORITY** 
**When to Use:**
- ✅ Blocking issues that prevent progress
- ✅ Mandatory actions required for ALL changes
- ✅ Features essential for system operation
- ✅ Urgent actions requiring immediate attention

**When NOT to Use:**
- ❌ Important but non-blocking features
- ❌ High-priority items that aren't critical
- ❌ Phase indicators or step markers

**Examples:**
```markdown
[CRITICAL] CRITICAL: All code changes require feature tracking
[CRITICAL] MANDATORY: Follow protocol validation steps
[CRITICAL] BLOCKING: Resolve dependency conflicts before extraction
```

### [HIGH] **HIGH PRIORITY**
**When to Use:**
- ✅ Important features requiring documentation updates
- ✅ High-value tasks that should be completed soon  
- ✅ Features that significantly impact the system
- ✅ Tasks that affect multiple components

**When NOT to Use:**
- ❌ Critical/blocking items (use [CRITICAL] instead)
- ❌ Medium-importance conditional items
- ❌ Document purpose headers

**Examples:**
```markdown
[HIGH] HIGH: Configuration system requires comprehensive testing
[HIGH] IMPORTANT: Git integration affects multiple workflows
[HIGH] SIGNIFICANT: Format system impacts all user output
```

### [MEDIUM] **MEDIUM PRIORITY**
**When to Use:**
- ✅ Conditional updates that depend on other factors
- ✅ Secondary features with moderate importance
- ✅ Tasks that can be delayed without major impact
- ✅ Features requiring evaluation before implementation

**When NOT to Use:**
- ❌ High-importance features (use [HIGH] instead)
- ❌ Analysis document headers (keep as document type)
- ❌ Low-priority cleanup tasks

**Examples:**
```markdown
[MEDIUM] CONDITIONAL: Update architecture.md IF structural changes
[MEDIUM] EVALUATE: Consider implementing advanced Git features
[MEDIUM] SECONDARY: Performance optimizations for edge cases
```

### [LOW] **LOW PRIORITY**
**When to Use:**
- ✅ Cleanup tasks with minimal impact
- ✅ Optional improvements and enhancements
- ✅ Tasks that can be executed last
- ✅ Minor updates and refinements

**When NOT to Use:**
- ❌ Medium-importance items (use [MEDIUM] instead)  
- ❌ Documentation action headers
- ❌ Reference materials

**Examples:**
```markdown
[LOW] MINOR: Update code comments for clarity
[LOW] OPTIONAL: Add advanced configuration examples
[LOW] CLEANUP: Remove outdated implementation notes
```

## 📋 Process Execution Guidelines

### 🚀 **PHASE 1** (Foundation)
**When to Use:**
- ✅ Initial validation requirements
- ✅ Pre-work that must happen first
- ✅ Foundation setup and preparation
- ✅ Critical validation before other work

**Usage Pattern:**
```markdown
### 🚀 PHASE 1: CRITICAL VALIDATION [Execute FIRST - MANDATORY]
Before making ANY changes...
```

### ⚡ **PHASE 2** (Core Work)
**When to Use:**
- ✅ Main implementation tasks
- ✅ Core work and primary activities
- ✅ Central features and functionality
- ✅ Primary documentation updates

**Usage Pattern:**
```markdown
### ⚡ PHASE 2: CORE IMPLEMENTATION [Execute SECOND - HIGH PRIORITY]
For main development work...
```

### [ACTION:migration] **PHASE 3** (Follow-up)
**When to Use:**
- ✅ Conditional references and evaluations
- ✅ Follow-up work after core completion
- ✅ Iterative processes and refinements
- ✅ Secondary considerations

**Usage Pattern:**
```markdown
### [ACTION:migration] PHASE 3: CONDITIONAL REFERENCES [Execute THIRD - MEDIUM PRIORITY]
Reference only if changes affect...
```

### 🏁 **PHASE 4** (Completion)
**When to Use:**
- ✅ Final wrap-up tasks
- ✅ Process completion activities
- ✅ Final validation and verification
- ✅ End-of-process documentation

**Usage Pattern:**
```markdown
### 🏁 PHASE 4: PROCESS COMPLETION [Execute LAST - LOW PRIORITY]
For final wrap-up and completion...
```

## 📑 Process Steps Guidelines

### Step Numbering Rules
- **1️⃣ STEP 1**: Always first step, discovery/analysis phase
- **2️⃣ STEP 2**: Always second step, planning/documentation phase
- **3️⃣ STEP 3**: Always third step, execution/implementation phase
- **✅ COMPLETE**: Step or process completion indicator

### Step Usage Patterns
```markdown
### 1️⃣ Step 1: Task Identification and Validation (EXECUTE FIRST)
### 2️⃣ Step 2: Change Type Classification (EXECUTE SECOND)  
### 3️⃣ Step 3: Execute Protocol and Validate (EXECUTE THIRD)
### ✅ Completion Criteria (All Must Pass)
```

## 🗂️ Document Category Guidelines

### 📑 **PURPOSE SECTION**
**When to Use:**
- ✅ Document purpose statements
- ✅ Overview and introduction sections
- ✅ Goal and objective descriptions
- ✅ Section purpose explanations

**Usage Pattern:**
```markdown
## 📑 Purpose
This document provides...

## 📑 Overview  
The system architecture includes...
```

### 📋 **CHECKLIST DOCUMENT**
**When to Use:**
- ✅ Process checklists and step-by-step guides
- ✅ Validation checklists
- ✅ Task lists and action items
- ✅ Procedural documentation

**Usage Pattern:**
```markdown
# 📋 Context File Checklist: AI Assistant Guide
## 📋 AI Assistant Validation Checklist
```

### 📊 **ANALYSIS DOCUMENT**
**When to Use:**
- ✅ Data analysis and assessment sections
- ✅ Metrics and measurement documentation
- ✅ Evaluation and comparison content
- ✅ Statistical or analytical presentations

**Usage Pattern:**
```markdown
## 📊 Implementation Impact Assessment
## 📊 Complete Icon Inventory by Usage Frequency
```

### 📖 **REFERENCE DOCUMENT**
**When to Use:**
- ✅ Specification documents
- ✅ Reference materials and guides
- ✅ Technical documentation
- ✅ API and interface references

**Usage Pattern:**
```markdown
## 📖 Technical Specifications
## 📖 API Reference Guide
```

## [ACTION:core-functionality] Action Category Guidelines

### [ACTION:discovery] **SEARCH/DISCOVER**
**When to Use:**
- ✅ Search actions and discovery tasks
- ✅ Analysis and investigation activities
- ✅ Research and exploration tasks
- ✅ Finding and locating activities

**Usage Pattern:**
```markdown
1. **[ACTION:discovery] SEARCH FOR EXISTING TOKENS** [PRIORITY: CRITICAL]
## [ACTION:discovery] Token Search Quick Commands
```

### [ACTION:format-processing] **DOCUMENT/UPDATE**
**When to Use:**
- ✅ Writing and documentation tasks
- ✅ Update and modification actions
- ✅ Record and note-taking activities
- ✅ Content creation and editing

**Usage Pattern:**
```markdown
**Phase 2: Documentation Updates [ACTION:format-processing]**
- [ ] [ACTION:format-processing] Updated all required documentation files
```

### [ACTION:core-functionality] **CONFIGURE/MODIFY**
**When to Use:**
- ✅ Configuration and setup tasks
- ✅ System modification and adjustment
- ✅ Tool and environment configuration
- ✅ Parameter and setting changes

**Usage Pattern:**
```markdown
**[ACTION:core-functionality] MODIFY EXISTING → Execute: MODIFICATION Protocol**
## [ACTION:core-functionality] Configuration Management
```

### [ACTION:validation] **PROTECT/VALIDATE**
**When to Use:**
- ✅ Security and protection tasks
- ✅ Validation and verification activities
- ✅ Guard and safety functions
- ✅ Compliance and enforcement

**Usage Pattern:**
```markdown
**[ACTION:validation] MANDATORY PRE-WORK VALIDATION**
1. **[ACTION:validation] Immutable Check**: Verify no conflicts
```

## 🚫 Common Mistakes and How to Avoid Them

### ❌ **Mistake 1: Priority Confusion**
```markdown
# WRONG - Using wrong priority icons
🎯 HIGH PRIORITY  # Old system - creates confusion
📋 MANDATORY     # Wrong - mixes mandatory with checklist

# CORRECT - Using standardized priority icons  
[HIGH] HIGH PRIORITY  # Clear priority indication
[CRITICAL] MANDATORY     # Clear critical priority
```

### ❌ **Mistake 2: Phase/Step Mixing**
```markdown
# WRONG - Mixing phases and steps
🚨 PHASE 1       # Old system - creates confusion
📋 STEP 2        # Wrong - mixes step with checklist

# CORRECT - Using proper phase/step icons
🚀 PHASE 1       # Clear phase indication
2️⃣ STEP 2        # Clear step indication
```

### ❌ **Mistake 3: Context-Dependent Meanings**
```markdown
# WRONG - Same icon, different meanings
🎯 Purpose       # In one document
🎯 HIGH PRIORITY # In another document - CONFLICT!

# CORRECT - Unique meanings
📑 Purpose       # Always means purpose
[HIGH] HIGH PRIORITY # Always means high priority
```

### ❌ **Mistake 4: Semantic Overloading**
```markdown
# WRONG - One icon for multiple concepts
📋 Checklist document
📋 Mandatory requirement  
📋 Task verification

# CORRECT - Distinct icons for distinct concepts
📋 Checklist document    # Document type only
[CRITICAL] Mandatory requirement # Priority level only
1️⃣ Task verification     # Process step only
```

## ✅ Icon Selection Decision Tree

When choosing an icon, follow this decision process:

### 🔥 **Priority Classification**
```
Is this about priority/importance?
├─ YES: Use priority hierarchy ([CRITICAL][HIGH][MEDIUM][LOW])
│   ├─ Blocking/Critical? → [CRITICAL]
│   ├─ Important/High? → [HIGH]  
│   ├─ Conditional/Medium? → [MEDIUM]
│   └─ Minor/Low? → [LOW]
└─ NO: Continue to next category
```

### 📋 **Process Classification**  
```
Is this about process execution?
├─ YES: Determine type
│   ├─ Sequential phases? → 🚀⚡[ACTION:migration]🏁
│   ├─ Numbered steps? → 1️⃣2️⃣3️⃣
│   └─ Completion status? → ✅
└─ NO: Continue to next category
```

### 🗂️ **Document Classification**
```
Is this about document navigation?
├─ YES: Use document categories
│   ├─ Purpose/Overview? → 📑
│   ├─ Checklist/Process? → 📋
│   ├─ Analysis/Data? → 📊
│   └─ Reference/Spec? → 📖
└─ NO: Continue to next category
```

### [ACTION:core-functionality] **Action Classification**
```
Is this about actions/tasks?
├─ YES: Use action categories
│   ├─ Search/Discovery? → [ACTION:discovery]
│   ├─ Documentation? → [ACTION:format-processing]
│   ├─ Configuration? → [ACTION:core-functionality]
│   └─ Validation? → [ACTION:validation]
└─ NO: Use specialized icons (🧪🆕🔌🐛🚫⚠️)
```

## [ACTION:migration] Validation and Enforcement

### ✅ **Self-Validation Checklist**
Before using any icon, verify:

- [ ] **Unique Meaning**: Icon has exactly one meaning across all documentation
- [ ] **Correct Category**: Icon belongs to appropriate category (priority/process/document/action)  
- [ ] **Consistent Usage**: Icon usage matches established patterns
- [ ] **Clear Semantics**: Icon meaning is obvious to AI assistants and humans
- [ ] **Future-Proof**: Usage doesn't create potential conflicts

### [ACTION:validation] **Review Process Integration**
All documentation changes must verify:

1. **Icon Compliance**: All icons follow these guidelines
2. **Conflict Prevention**: No new conflicts introduced
3. **Meaning Consistency**: Same icons mean same things everywhere
4. **Category Adherence**: Icons used in correct categories

### [ACTION:core-functionality] **Enforcement Tools**
```bash
# Validate icon consistency across documents
grep -r "🚨\|🎯\|📋" docs/context/ # Should find zero legacy conflicts
grep -r "[CRITICAL]\|[HIGH]\|[MEDIUM]\|[LOW]" docs/context/ # Should find only priority usage

# Check for proper phase/step usage  
grep -r "🚀\|⚡\|[ACTION:migration]\|🏁" docs/context/ # Should find only phase usage
grep -r "1️⃣\|2️⃣\|3️⃣" docs/context/ # Should find only step usage
```

## 📋 Quick Reference

### 🔥 **Priority Icons**
| Icon | Meaning | Usage |
|------|---------|-------|
| [CRITICAL] | CRITICAL | Blocking, mandatory, urgent |
| [HIGH] | HIGH | Important, significant impact |
| [MEDIUM] | MEDIUM | Conditional, balanced consideration |
| [LOW] | LOW | Minor, cleanup, optional |

### 📋 **Process Icons** 
| Icon | Meaning | Usage |
|------|---------|-------|
| 🚀⚡[ACTION:migration]🏁 | PHASES 1-4 | Sequential phases |
| 1️⃣2️⃣3️⃣ | STEPS 1-3 | Numbered steps |
| ✅ | COMPLETE | Success/completion |

### 🗂️ **Category Icons**
| Icon | Meaning | Usage |
|------|---------|-------|
| 📑📋📊📖 | DOCUMENT TYPES | Navigation categories |
| [ACTION:discovery][ACTION:format-processing][ACTION:core-functionality][ACTION:validation] | ACTION TYPES | Task categories |

---

**📋 Guidelines Status**: Official - Mandatory for all documentation
**[ACTION:validation] Enforcement**: Integrated with documentation review process
**[CRITICAL] Critical**: Zero tolerance for icon conflicts

// DOC-006: Icon standardization - Usage guidelines complete 