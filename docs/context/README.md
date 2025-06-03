# 🤖 AI Assistant Context Documentation Index

> **🚨 CRITICAL**: This is the MASTER INDEX for all AI assistants working on the BkpDir project. All AI assistants MUST start here and follow the priority order below.

## 🎯 AI Assistant Quick Start Guide

### 🚨 PHASE 1: CRITICAL VALIDATION [Execute FIRST - MANDATORY]

Before making ANY code changes, AI assistants MUST:

1. **🛡️ [Immutable Requirements Check](immutable.md)** - Verify no conflicts with unchangeable specifications
2. **📋 [Feature Tracking Registry](feature-tracking.md)** - Find or create Feature ID for your task
3. **🔍 [AI Assistant Compliance](ai-assistant-compliance.md)** - Review token requirements and response format
4. **🚨 [AI Assistant Protocol](ai-assistant-protocol.md)** - Follow the appropriate change protocol

### 🎯 PHASE 2: CORE DOCUMENTATION [Execute SECOND - HIGH PRIORITY]

For understanding and implementing changes:

5. **🏗️ [Architecture](architecture.md)** - System design and technical components
6. **📝 [Requirements](requirements.md)** - Implementation requirements and constraints  
7. **🧪 [Testing](testing.md)** - Test coverage requirements and validation standards
8. **📖 [Specification](specification.md)** - User-facing features and behaviors

### 📊 PHASE 3: CONDITIONAL REFERENCES [Execute THIRD - MEDIUM PRIORITY]

Reference only if your changes affect these areas:

9. **⚙️ [Implementation Decisions](implementation-decisions.md)** - IF making architectural decisions
10. **🔧 [Validation Automation](validation-automation.md)** - IF adding validation processes
11. **🔄 [Sync Framework](sync-framework.md)** - IF affecting data synchronization
12. **📈 [Implementation Status](implementation-status.md)** - For progress tracking

### 📝 PHASE 4: PROCESS REFERENCES [Execute LAST - LOW PRIORITY]

Reference only for process understanding (NEVER modify):

13. **📋 [Context File Checklist](context-file-checklist.md)** - File update guidelines
14. **🔒 [Enforcement Mechanisms](enforcement-mechanisms.md)** - Validation rules
15. **🚫 [Change Rejection Criteria](change-rejection-criteria.md)** - Common rejection scenarios
16. **📚 [Feature Documentation Standards](feature-documentation-standards.md)** - Documentation guidelines

## 🚨 AI Assistant Priority Matrix

### 🔥 CRITICAL ACTIONS [NEVER SKIP]
```
┌─ 🛡️ Check immutable.md for conflicts
├─ 📋 Verify Feature ID exists in feature-tracking.md  
├─ 🔍 Review ai-assistant-compliance.md requirements
└─ 🚨 Follow ai-assistant-protocol.md for change type
```

### ⚡ HIGH PRIORITY [ALWAYS EXECUTE]
```
┌─ 🏗️ Update architecture.md if technical changes
├─ 📝 Update requirements.md if new requirements
├─ 🧪 Update testing.md if test changes
└─ 📖 Update specification.md if user-facing changes
```

### 📊 MEDIUM PRIORITY [EVALUATE CONDITIONALLY]
```
┌─ ⚙️ Update implementation-decisions.md if architectural
├─ 🔧 Update validation-automation.md if validation changes
├─ 🔄 Update sync-framework.md if sync changes
└─ 📈 Update implementation-status.md for progress
```

### ❌ NEVER MODIFY [REFERENCE ONLY]
```
┌─ 🛡️ immutable.md - Only check for conflicts
├─ 🔒 enforcement-mechanisms.md - Process reference only
├─ 📋 context-file-checklist.md - Guidelines only
└─ 🚫 change-rejection-criteria.md - Reference only
```

## 🎯 Change Type Quick Reference

When you know your change type, jump directly to the protocol:

| Change Type | Icon | Protocol Section | Priority | Documentation Impact |
|-------------|------|------------------|----------|---------------------|
| **New Feature** | 🆕 | [NEW FEATURE Protocol](ai-assistant-protocol.md#-new-feature-protocol-priority-critical) | 🚨 CRITICAL | Full documentation cascade |
| **Modify Existing** | 🔧 | [MODIFICATION Protocol](ai-assistant-protocol.md#-modification-protocol-priority-critical) | 🚨 CRITICAL | Impact analysis required |
| **Bug Fix** | 🐛 | [BUG FIX Protocol](ai-assistant-protocol.md#-bug-fix-protocol-priority-medium) | 📊 MEDIUM | Minimal documentation |
| **Config Change** | ⚙️ | [CONFIG CHANGE Protocol](ai-assistant-protocol.md#-config-change-protocol-priority-medium) | 📊 MEDIUM | Configuration focus |
| **API/Interface** | 🔌 | [API CHANGE Protocol](ai-assistant-protocol.md#-api-change-protocol-priority-critical) | 🚨 CRITICAL | Interface documentation |
| **Test Only** | 🧪 | [TEST ADDITION Protocol](ai-assistant-protocol.md#-test-addition-protocol-priority-low) | 📝 LOW | Testing documentation |
| **Performance** | 🚀 | [PERFORMANCE Protocol](ai-assistant-protocol.md#-performance-protocol-priority-medium) | 📊 MEDIUM | Architecture documentation |
| **Refactoring** | 🔄 | [REFACTORING Protocol](ai-assistant-protocol.md#-refactoring-protocol-priority-low) | 📝 LOW | Structural documentation |

## 🔍 Token Search Quick Commands

```bash
# Search for existing feature tokens
grep -r "// [A-Z]+-[0-9]+" docs/context/

# Find specific token patterns
grep -r "ARCH-[0-9]+" docs/context/   # Archive operations
grep -r "FILE-[0-9]+" docs/context/   # File operations  
grep -r "CFG-[0-9]+" docs/context/    # Configuration
grep -r "GIT-[0-9]+" docs/context/    # Git integration
grep -r "TEST-[0-9]+" docs/context/   # Testing infrastructure

# Validate all changes
make test && make lint
```

## ✅ AI Assistant Validation Checklist

Before submitting any code changes, verify:

- [ ] 🔍 Searched for existing tokens related to your changes
- [ ] 📋 Referenced specific Feature ID from feature-tracking.md
- [ ] 🎯 Followed appropriate protocol from ai-assistant-protocol.md
- [ ] 🏷️ Added implementation tokens to all modified code
- [ ] 📝 Updated all required documentation files
- [ ] 🧪 All tests pass (`make test`)
- [ ] 🔧 All lint checks pass (`make lint`)
- [ ] 🏁 Marked task complete in feature-tracking.md

## 🚨 Critical Reminders for AI Assistants

1. **🛡️ IMMUTABLE REQUIREMENTS**: Never modify anything that conflicts with `immutable.md`
2. **📋 FEATURE TRACKING**: Every code change MUST have a corresponding Feature ID
3. **🏷️ IMPLEMENTATION TOKENS**: Every modified function/method needs `// FEATURE-ID: Description`
4. **📝 DOCUMENTATION SYNC**: Update ALL affected context files, not just code
5. **✅ VALIDATION**: Run tests and linting before marking tasks complete

## 📞 Emergency Quick Reference

If you're unsure about anything:

1. **🔍 Search First**: Use `grep -r "relevant-term" docs/context/`
2. **📋 Check Feature Registry**: Look in feature-tracking.md for existing features
3. **🚨 Follow Protocol**: Use ai-assistant-protocol.md for your change type
4. **🛡️ Check Immutable**: Verify no conflicts with immutable.md
5. **🧪 Validate Changes**: Run `make test && make lint`

**🤖 This index serves as your primary navigation hub. Bookmark this document and always start here!** 