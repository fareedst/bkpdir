# Validation and Automation Tools

## Purpose
This document contains comprehensive validation and automation tools for maintaining context documentation synchronization and enforcing development standards established in the feature tracking system.

> **🔗 Related Documents**: See [`feature-tracking.md`](feature-tracking.md) for the core feature matrix and [`ai-assistant-compliance.md`](ai-assistant-compliance.md) for AI assistant requirements.

## 🔧 VALIDATION AND AUTOMATION TOOLS

### 📋 PRE-COMMIT VALIDATION CHECKLIST

Before ANY commit that includes code changes, run this complete validation checklist:

```bash
# 1. MANDATORY: Check for feature tracking updates
echo "🔍 Checking feature tracking compliance..."
git diff --cached --name-only | grep -E '\.(go|yaml|yml)$' && echo "✅ Code changes detected - context files MUST be updated" || echo "ℹ️  No code changes"

# 2. MANDATORY: Verify context file updates
echo "📋 Verifying context documentation updates..."
git diff --cached --name-only | grep -E '^docs/context/' && echo "✅ Context files updated" || echo "❌ ERROR: Code changes without context updates"

# 3. MANDATORY: Check implementation tokens
echo "🏷️  Checking implementation tokens..."
git diff --cached | grep -E '^\+.*//.*[A-Z]+-[0-9]+:' && echo "✅ Implementation tokens found" || echo "⚠️  Warning: No new implementation tokens"

# 4. MANDATORY: Validate feature matrix
echo "📊 Validating feature matrix..."
grep -c "| [A-Z]+-[0-9]" docs/context/feature-tracking.md && echo "✅ Feature matrix has entries" || echo "❌ ERROR: Feature matrix empty"

# 5. MANDATORY: Check cross-references
echo "🔗 Checking cross-references..."
./scripts/validate-docs.sh || echo "❌ ERROR: Documentation validation failed"
```

### 🚨 AUTOMATED VALIDATION COMMANDS

#### **Quick Context Check**
```bash
# Check if context files need updates for current changes
function check_context_updates() {
    local code_changes=$(git diff --cached --name-only | grep -E '\.(go|yaml|yml)$' | wc -l)
    local context_changes=$(git diff --cached --name-only | grep -E '^docs/context/' | wc -l)
    
    if [ $code_changes -gt 0 ] && [ $context_changes -eq 0 ]; then
        echo "❌ ERROR: $code_changes code file(s) changed but no context files updated"
        echo "📋 Required context files to check:"
        echo "   - docs/context/feature-tracking.md (ALWAYS)"
        echo "   - docs/context/specification.md (if user-facing changes)"
        echo "   - docs/context/requirements.md (if requirements change)"
        echo "   - docs/context/architecture.md (if technical changes)"
        echo "   - docs/context/testing.md (if test changes)"
        return 1
    else
        echo "✅ Context file requirements satisfied"
        return 0
    fi
}
```

#### **Feature ID Validation**
```bash
# Validate all feature IDs are properly registered
function validate_feature_ids() {
    echo "🔍 Validating feature ID consistency..."
    
    # Extract feature IDs from code
    local code_features=$(grep -r "// [A-Z]+-[0-9]" --include="*.go" . | sed 's/.*\/\/ \([A-Z]*-[0-9]*\).*/\1/' | sort -u)
    
    # Extract feature IDs from tracking matrix
    local matrix_features=$(grep "| [A-Z]+-[0-9]" docs/context/feature-tracking.md | sed 's/|.*\([A-Z]*-[0-9]*\).*/\1/' | sort -u)
    
    # Find unregistered features
    local unregistered=$(comm -23 <(echo "$code_features") <(echo "$matrix_features"))
    
    if [ -n "$unregistered" ]; then
        echo "❌ ERROR: Unregistered feature IDs found in code:"
        echo "$unregistered"
        echo "💡 Add these features to docs/context/feature-tracking.md"
        return 1
    else
        echo "✅ All feature IDs properly registered"
        return 0
    fi
}
```

#### **Implementation Token Audit**
```bash
# Audit implementation tokens for completeness
function audit_implementation_tokens() {
    echo "🏷️  Auditing implementation tokens..."
    
    local missing_tokens=""
    
    # Check each .go file for feature-related functions without tokens
    for file in $(find . -name "*.go" -not -path "./vendor/*"); do
        # Look for functions that might need tokens but don't have them
        local suspicious_functions=$(grep -n "^func.*\(Archive\|Backup\|Config\|Git\|Format\|Test\)" "$file" | grep -v "// [A-Z]+-[0-9]")
        
        if [ -n "$suspicious_functions" ]; then
            missing_tokens="$missing_tokens\n$file:\n$suspicious_functions"
        fi
    done
    
    if [ -n "$missing_tokens" ]; then
        echo "⚠️  Functions potentially missing implementation tokens:"
        echo -e "$missing_tokens"
        echo "💡 Consider adding feature ID tokens to these functions"
    else
        echo "✅ Implementation token coverage looks good"
    fi
}
```

### 📄 DOCUMENTATION SYNC VALIDATION

#### **Cross-Reference Checker**
```bash
# Check all cross-references between context documents
function check_cross_references() {
    echo "🔗 Checking cross-references between context documents..."
    
    local broken_refs=""
    
    # Check references to feature IDs
    for doc in docs/context/*.md; do
        local doc_name=$(basename "$doc")
        local refs=$(grep -o "[A-Z]+-[0-9]\+" "$doc" | sort -u)
        
        for ref in $refs; do
            if ! grep -q "| $ref |" docs/context/feature-tracking.md; then
                broken_refs="$broken_refs\n$doc_name references undefined $ref"
            fi
        done
    done
    
    if [ -n "$broken_refs" ]; then
        echo "❌ ERROR: Broken cross-references found:"
        echo -e "$broken_refs"
        return 1
    else
        echo "✅ All cross-references valid"
        return 0
    fi
}
```

#### **Status Consistency Checker**
```bash
# Check feature status consistency across documents
function check_status_consistency() {
    echo "📊 Checking feature status consistency..."
    
    local inconsistencies=""
    
    # Extract features and their status from feature-tracking.md
    while IFS='|' read -r feature_id spec req arch test status token; do
        # Clean up the values
        feature_id=$(echo "$feature_id" | xargs)
        status=$(echo "$status" | xargs)
        
        if [[ "$feature_id" =~ ^[A-Z]+-[0-9]+$ ]]; then
            # Check if status matches in other documents
            # This is a simplified check - real implementation would be more thorough
            if [ "$status" = "Completed" ] || [ "$status" = "Implemented" ]; then
                if ! grep -q "$feature_id" docs/context/specification.md; then
                    inconsistencies="$inconsistencies\n$feature_id marked $status but not in specification.md"
                fi
            fi
        fi
    done < <(grep "| [A-Z]+-[0-9]" docs/context/feature-tracking.md)
    
    if [ -n "$inconsistencies" ]; then
        echo "⚠️  Status inconsistencies found:"
        echo -e "$inconsistencies"
    else
        echo "✅ Status consistency looks good"
    fi
}
```

### 🎯 DEVELOPER WORKFLOW INTEGRATION

#### **Git Hooks Setup**
```bash
# Add to .git/hooks/pre-commit
#!/bin/bash
echo "🚨 MANDATORY: Context Documentation Validation"

# Source the validation functions
source scripts/context-validation.sh

# Run all validation checks
check_context_updates || exit 1
validate_feature_ids || exit 1
check_cross_references || exit 1

echo "✅ All context documentation requirements satisfied"
```

#### **Make Target Integration**
```bash
# Add to Makefile
.PHONY: validate-context
validate-context:
	@echo "🔍 Validating context documentation..."
	@bash scripts/context-validation.sh
	@echo "✅ Context validation complete"

.PHONY: pre-commit
pre-commit: validate-context test
	@echo "✅ Pre-commit validation passed"

# Ensure context validation runs before tests
test: validate-context
	go test ./...
```

#### **IDE Integration Hints**
```bash
# VSCode settings.json snippet
{
    "go.buildTags": "integration",
    "files.watcherExclude": {
        "**/docs/context/**": false
    },
    "search.exclude": {
        "**/docs/context/**": false
    },
    "todo-tree.regex.regex": "((//|#|<!--|;|/\\*|^)\\s*($TAGS)|^\\s*- \\[ \\])",
    "todo-tree.regex.regexFlags": "gim",
    "todo-tree.highlights.customHighlight": {
        "FEATURE-ID": {
            "icon": "tag",
            "type": "tag",
            "foreground": "#FF6B6B"
        }
    }
}
```

### 📈 METRICS AND MONITORING

#### **Documentation Coverage Metrics**
```bash
# Generate context documentation coverage report
function generate_coverage_report() {
    echo "📊 Context Documentation Coverage Report"
    echo "========================================"
    
    local total_features=$(grep -c "| [A-Z]+-[0-9]" docs/context/feature-tracking.md)
    local implemented_features=$(grep -c "| [A-Z]+-[0-9].*Implemented\|Completed" docs/context/feature-tracking.md)
    local coverage_percentage=$((implemented_features * 100 / total_features))
    
    echo "Total Features: $total_features"
    echo "Implemented Features: $implemented_features"
    echo "Implementation Coverage: $coverage_percentage%"
    
    local code_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
    local files_with_tokens=$(grep -l "// [A-Z]+-[0-9]" --include="*.go" -r . | wc -l)
    local token_coverage=$((files_with_tokens * 100 / code_files))
    
    echo "Code Files: $code_files"
    echo "Files with Tokens: $files_with_tokens"
    echo "Token Coverage: $token_coverage%"
}
```

### 🔄 CONTINUOUS IMPROVEMENT

#### **Weekly Validation Audit**
```bash
# Run comprehensive weekly audit
function weekly_audit() {
    echo "📅 Weekly Context Documentation Audit"
    echo "====================================="
    
    check_context_updates
    validate_feature_ids
    audit_implementation_tokens
    check_cross_references
    check_status_consistency
    generate_coverage_report
    
    echo "📋 Action Items:"
    echo "- Review any warnings or errors above"
    echo "- Update missing implementation tokens"
    echo "- Verify feature status accuracy"
    echo "- Check for outdated documentation"
}
```

### 💡 AUTOMATION RECOMMENDATIONS

1. **CI/CD Integration**: Add context validation to GitHub Actions/Jenkins
2. **Documentation Generator**: Create tools to auto-update parts of context files
3. **Template Generation**: Auto-generate boilerplate for new features
4. **Dependency Tracking**: Monitor changes to shared components
5. **Compliance Dashboard**: Web interface showing documentation health

This comprehensive validation framework ensures that the context documentation remains synchronized with code changes and maintains the high quality standards established in the feature tracking system.

### 📝 IMPLEMENTATION TOKEN REQUIREMENTS (DETAILED)

#### **Token Format Standards**
- **Standard Format**: `// FEATURE-ID: Brief description`
- **Examples**:
  ```go
  // ARCH-001: Archive naming convention implementation
  func GenerateArchiveName(cfg *Config, dir string) string {
      // ARCH-001: Include timestamp in archive name
      timestamp := time.Now().Format("2006-01-02-15-04")
      
      // GIT-002: Add Git branch and hash if available
      gitInfo := getGitInfo()
      if gitInfo.Branch != "" {
          return fmt.Sprintf("%s-%s=%s=%s.zip", dir, timestamp, gitInfo.Branch, gitInfo.Hash)
      }
      
      return fmt.Sprintf("%s-%s.zip", dir, timestamp)
  }
  
  // CFG-001: Configuration discovery implementation
  func GetConfigSearchPath() []string {
      // CFG-001: Check environment variable first
      if envPath := os.Getenv("BKPDIR_CONFIG"); envPath != "" {
          return strings.Split(envPath, ":")
      }
      
      // CFG-001: Use default search path
      return []string{"./.bkpdir.yml", "~/.bkpdir.yml"}
  }
  ```

#### **Token Placement Requirements**
1. **Function Level**: Every public function must have a feature token comment
2. **Method Level**: Every method must reference the relevant feature
3. **Code Block Level**: Significant logic blocks within functions should have tokens
4. **Error Handling**: Error paths should include tokens for traceability
5. **Configuration Handling**: All config reading/writing must have tokens

#### **AI ASSISTANT TOKEN HANDLING INSTRUCTIONS**

**🤖 FOR AI ASSISTANTS: AUTOMATIC TOKEN DETECTION AND REFERENCING**

When making code changes, AI assistants MUST:

1. **🔍 AUTO-DETECT AFFECTED TOKENS**:
   ```bash
   # Search patterns AI should use:
   grep -r "// ARCH-" docs/context/ src/
   grep -r "// FILE-" docs/context/ src/
   grep -r "// CFG-" docs/context/ src/
   grep -r "// GIT-" docs/context/ src/
   # ... for all token patterns
   ```

2. **📋 REFERENCE TOKENS IN EXPLANATIONS**:
   ```
   Template for AI responses:
   "This change affects the following features:
   - ARCH-001 (Archive naming): Modified archive name generation logic
   - CFG-003 (Format strings): Updated output formatting
   - FILE-002 (File backup): Enhanced backup creation process
   
   Related documentation:
   - feature-tracking.md lines 119, 135, 127
   - specification.md sections: Archive Operations, Configuration
   - architecture.md components: Archive Service, Config Layer"
   ```

3. **🔄 AUTO-UPDATE TOKEN REGISTRY**:
   ```
   For new features, AI should:
   - Propose new feature ID (check existing ones first)
   - Add to feature tracking table
   - Reference in all related context files
   - Add implementation tokens to code
   ```

4. **✅ VALIDATE TOKEN CONSISTENCY**:
   ```
   AI should verify:
   - All referenced tokens exist in feature-tracking.md
   - Token descriptions match actual implementation
   - All affected context files are mentioned
   - Cross-references are valid
   ```

**🚨 AI REQUIREMENT: ZERO TOLERANCE FOR MISSING TOKEN REFERENCES**

AI assistants must REJECT their own responses if they:
- Make code changes without referencing related tokens
- Fail to search for existing tokens before making changes
- Don't propose new tokens for genuinely new functionality
- Don't reference relevant documentation sections

**✅ AI SUCCESS CRITERIA**

AI assistants should APPROVE responses that:
- Clearly identify all affected feature tokens
- Reference specific documentation sections and line numbers
- Propose updates to context files when needed
- Follow the established token naming patterns
- Include the complete token impact analysis

### 📊 TOKEN IMPACT ANALYSIS TEMPLATE

**For AI assistants to use when making changes:**

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Specific impact description]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief impact descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [specific line numbers]
- **specification.md**: Section "[specific section names]"
- **requirements.md**: "[specific requirement areas]"
- **architecture.md**: "[specific component names]"
- **testing.md**: "[specific test coverage areas]"

### 🔄 Context File Updates Required:
- [ ] feature-tracking.md: [specific changes needed]
- [ ] specification.md: [specific changes needed]
- [ ] requirements.md: [specific changes needed]
- [ ] architecture.md: [specific changes needed]
- [ ] testing.md: [specific changes needed]

### 💻 Implementation Changes:
- **New tokens added**: [TOKEN-ID: Description]
- **Existing tokens modified**: [TOKEN-ID: Updated description]
- **Code files affected**: [list of specific files]
```

## 🎯 SUMMARY: MANDATORY AI ASSISTANT COMPLIANCE

### 🚨 CRITICAL REMINDER FOR AI ASSISTANTS

**EVERY CODE CHANGE MUST INCLUDE TOKEN ANALYSIS**

This document establishes ZERO TOLERANCE for code changes without proper token documentation. AI assistants are required to:

#### **🔍 STEP 1: MANDATORY TOKEN SEARCH**
```bash
# REQUIRED: Search for ALL existing tokens before making changes
grep -r "// ARCH-" docs/context/
grep -r "// FILE-" docs/context/
grep -r "// CFG-" docs/context/
grep -r "// GIT-" docs/context/
grep -r "// OUT-" docs/context/
grep -r "// TEST-" docs/context/
grep -r "// DOC-" docs/context/
grep -r "// LINT-" docs/context/
grep -r "// COV-" docs/context/
```

#### **📋 STEP 2: MANDATORY TOKEN IMPACT STATEMENT**
Every AI response with code changes MUST include:

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Specific impact description]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief impact descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [specific line numbers]
- **specification.md**: Section "[specific section names]"
- **requirements.md**: "[specific requirement areas]"
- **architecture.md**: "[specific component names]"
- **testing.md**: "[specific test coverage areas]"

### 🔄 Context File Updates Required:
- [ ] feature-tracking.md: [specific changes needed]
- [ ] specification.md: [specific changes needed]
- [ ] requirements.md: [specific changes needed]
- [ ] architecture.md: [specific changes needed]
- [ ] testing.md: [specific changes needed]

### 💻 Implementation Changes:
- **New tokens added**: [TOKEN-ID: Description]
- **Existing tokens modified**: [TOKEN-ID: Updated description]
- **Code files affected**: [list of specific files]
```

#### **🚫 STEP 3: MANDATORY REJECTION CRITERIA**

AI assistants MUST REFUSE to provide code changes if:
- No token search was performed
- No existing tokens are referenced for related functionality
- No new tokens are proposed for new functionality
- Context file updates are not specified
- Documentation cross-references are missing

#### **✅ STEP 4: MANDATORY APPROVAL CRITERIA**

AI assistants should ONLY proceed if their response:
- Includes complete token impact analysis using the required template
- References specific line numbers in feature-tracking.md
- Links to relevant sections in all context files
- Proposes concrete updates to affected documentation
- Follows established token naming conventions (ARCH-, FILE-, CFG-, etc.)

### 🔄 ENFORCEMENT SUMMARY

**REMEMBER: THIS IS NOT OPTIONAL**

The feature-tracking system is designed to prevent unplanned changes and maintain documentation synchronization. AI assistants that fail to follow these requirements undermine the entire system integrity.

**Key Success Metrics:**
- 100% of code changes must reference related tokens
- 100% of new features must have new tokens
- 100% of responses must include context file update plans
- 0% tolerance for undocumented changes

**Implementation Tokens Are Mandatory:**
Every function, method, and significant code block must include `// TOKEN-ID: Description` comments that trace back to this feature tracking matrix.

### 📞 QUICK REFERENCE FOR AI ASSISTANTS

1. **Search Pattern**: `grep -r "// [A-Z]+-[0-9]+" docs/context/`
2. **Token Format**: `// FEATURE-ID: Brief description`
3. **Required Template**: Use the Token Impact Analysis template above
4. **Context Files**: feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md
5. **Zero Tolerance**: No code changes without token references

**This document serves as the authoritative source for validation automation and token management compliance. Any AI assistant working with this codebase must comply with these requirements.** 