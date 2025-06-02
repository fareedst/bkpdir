# Feature Change Protocol

## 🔥 ENHANCED ACTIONABLE CHANGE PROCESS

### 📝 Adding New Features (STEP-BY-STEP)

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Assign Feature ID**: Use prefix (ARCH, FILE, CFG, GIT, etc.) + sequential number
2. Open specification.md → Add to relevant command/config section
3. Open requirements.md → Add to functional requirements section
4. Open architecture.md → Add to relevant component section
5. Open testing.md → Add to test coverage section
6. Add implementation tokens to code

### 🔄 Modifying Existing Features

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Update Feature ID**: Add modification suffix (e.g., ARCH-001-MOD-001)
2. Open specification.md → Update affected sections
3. Open requirements.md → Update affected requirements
4. Open architecture.md → Update affected components
5. Open testing.md → Update test coverage
6. Add modification tokens to code

### 🐛 Fixing Bugs

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Create Bug ID**: Use BUG- prefix + sequential number
2. Open specification.md → Document bug fix if affects behavior
3. Open requirements.md → Add bug fix requirements
4. Open architecture.md → Document technical fix
5. Open testing.md → Add regression tests
6. Add bug fix tokens to code

### 🛠️ Configuration Changes

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Update CFG- Feature ID**: Add modification suffix
2. Open specification.md → Update config sections
3. Open requirements.md → Update config requirements
4. Open architecture.md → Update config implementation
5. Open testing.md → Update config tests
6. Add config modification tokens

### 🔄 API/Interface Changes

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Update ARCH- Feature ID**: Add interface suffix
2. Open specification.md → Update interface specs
3. Open requirements.md → Update interface requirements
4. Open architecture.md → Update interface implementation
5. Open testing.md → Update interface tests
6. Add interface modification tokens

### ✅ Test Addition

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Create TEST- Feature ID**: Add sequential number
2. Open specification.md → Update test requirements
3. Open requirements.md → Add test requirements
4. Open architecture.md → Update test infrastructure
5. Open testing.md → Add test descriptions
6. Add test implementation tokens

### 🔧 Error Handling Changes

**⚡ IMMEDIATE ACTIONS REQUIRED:**
1. **Update ARCH- Feature ID**: Add error suffix
2. Open specification.md → Update error behavior
3. Open requirements.md → Update error requirements
4. Open architecture.md → Update error handling
5. Open testing.md → Update error tests
6. Add error handling tokens

### 📋 MANDATORY CONTEXT FILE CHECKLIST

#### **Phase 1: Pre-Change Analysis (REQUIRED)**
- [ ] **Feature Impact Analysis**: Identify which existing features are affected by the change
- [ ] **Context File Mapping**: Determine which context files require updates based on change type

| Change Type | Required Context File Updates |
|-------------|------------------------------|
| **New Feature** | feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md |
| **Feature Modification** | feature-tracking.md, + all files containing the feature |
| **Bug Fix** | feature-tracking.md (if affects documented behavior) |
| **Configuration Change** | feature-tracking.md, specification.md, requirements.md |
| **API/Interface Change** | feature-tracking.md, specification.md, architecture.md |
| **Test Addition** | feature-tracking.md, testing.md |
| **Error Handling Change** | feature-tracking.md, specification.md, architecture.md |

#### **Phase 2: Documentation Updates (REQUIRED)**
- [ ] **Update feature-tracking.md**: Add/modify feature entries, update status, add implementation tokens
- [ ] **Update specification.md**: Modify user-facing behavior descriptions if applicable
- [ ] **Update requirements.md**: Add/modify implementation requirements with traceability
- [ ] **Update architecture.md**: Update technical implementation descriptions
- [ ] **Update testing.md**: Add/modify test coverage requirements and descriptions
- [ ] **Check immutable.md**: Verify no immutable requirements are violated
- [ ] **Update cross-references**: Ensure all documents link correctly to each other

#### **Phase 3: Implementation Tokens (REQUIRED)**
- [ ] **Add Implementation Tokens**: Mark ALL modified code with feature ID comments
- [ ] **Update Token Registry**: Record new tokens in feature-tracking.md
- [ ] **Verify Token Consistency**: Ensure tokens match feature IDs and descriptions

#### **Phase 4: Validation (REQUIRED)**
- [ ] **Run Documentation Validation**: Execute `make validate-docs` (if available)
- [ ] **Cross-Reference Check**: Verify all links between documents are valid
- [ ] **Feature Matrix Update**: Ensure feature tracking matrix reflects all changes
- [ ] **Test Coverage Verification**: Confirm test references are accurate
