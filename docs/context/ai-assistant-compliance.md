# AI Assistant Compliance Guide

## ⚡ MANDATORY TOKEN REFERENCING FOR ALL AI ASSISTANTS

### 🚨 CRITICAL AI DIRECTIVE: ALWAYS REFERENCE CONTEXT FILE TOKENS WITH CODE CHANGES

**IMMEDIATE ACTION REQUIRED FOR ALL AI ASSISTANTS:**

When making ANY code changes, you MUST:

1. **🔍 SEARCH FOR EXISTING TOKENS**: Before making code changes, search the context directory files for relevant implementation tokens:
   - Use grep or search for patterns like `// ARCH-001:`, `// FILE-002:`, `// CFG-003:`, etc.
   - Look for feature IDs in the feature tracking tables in `feature-tracking.md`
   - Check all context files: `specification.md`, `requirements.md`, `architecture.md`, `testing.md`

2. **📋 REFERENCE RELEVANT TOKENS**: For every code change, explicitly mention related tokens:
   ```
   Example: "This change affects ARCH-001 (Archive naming), FILE-002 (File backup), and CFG-003 (Format strings)"
   ```

3. **🔗 LINK TO DOCUMENTATION**: Reference the specific sections in context files:
   ```
   Example: "See ARCH-001 in feature-tracking.md line 119, specification.md section 'Archive Naming', architecture.md 'Archive Service'"
   ```

4. **📝 UPDATE TOKEN REGISTRY**: Add new implementation tokens for new features:
   - Format: `// FEATURE-ID: Brief description`
   - Register in the feature tracking tables in `feature-tracking.md`
   - Update all affected context files

5. **✅ VALIDATE CROSS-REFERENCES**: Ensure all token references are consistent across:
   - `feature-tracking.md` 
   - `specification.md`
   - `requirements.md`
   - `architecture.md`
   - `testing.md`

### 🎯 TOKEN SEARCH PATTERNS FOR AI ASSISTANTS

**Search these patterns in context files:**
- `ARCH-[0-9]+` (Archive operations)
- `FILE-[0-9]+` (File operations)
- `CFG-[0-9]+` (Configuration)
- `GIT-[0-9]+` (Git integration)
- `OUT-[0-9]+` (Output management)
- `TEST-[0-9]+` (Testing infrastructure)
- `DOC-[0-9]+` (Documentation system)
- `LINT-[0-9]+` (Code quality)
- `COV-[0-9]+` (Coverage)

### 🚨 AI ENFORCEMENT RULES

**REJECT code changes that do not:**
- Reference at least one existing feature token when modifying related functionality
- Propose new tokens for genuinely new features
- Link to relevant documentation sections
- Update the feature tracking tables when adding new tokens

**ACCEPT code changes that:**
- Clearly state which tokens are affected
- Reference specific documentation sections
- Propose updates to context files when needed
- Follow the token naming conventions

## 🤖 MANDATORY AI VALIDATION CHECKLIST

Before any AI assistant provides code changes, it MUST validate:

1. **🔍 TOKEN SEARCH PERFORMED**
   - [ ] Searched all context files for existing feature tokens
   - [ ] Identified all relevant feature IDs using grep patterns: `ARCH-[0-9]+`, `FILE-[0-9]+`, `CFG-[0-9]+`, etc.
   - [ ] Listed affected tokens in response

2. **📋 TOKEN IMPACT DOCUMENTED**
   - [ ] Clearly stated which existing tokens are affected
   - [ ] Proposed new tokens for new functionality
   - [ ] Referenced specific line numbers in feature-tracking.md
   - [ ] Linked to relevant sections in other context files

3. **📄 DOCUMENTATION CROSS-REFERENCES**
   - [ ] Referenced specification.md sections
   - [ ] Referenced requirements.md areas
   - [ ] Referenced architecture.md components
   - [ ] Referenced testing.md coverage areas

4. **🔄 UPDATE REQUIREMENTS IDENTIFIED**
   - [ ] Listed which context files need updates
   - [ ] Specified what changes are needed in each file
   - [ ] Provided update templates or examples

## 📋 AI RESPONSE TEMPLATE (REQUIRED FORMAT)

```
## Token Impact Analysis

### 🔍 Affected Tokens:
- **Primary**: [TOKEN-ID] - [Description of impact]
- **Secondary**: [TOKEN-ID], [TOKEN-ID] - [Brief descriptions]

### 📚 Documentation References:
- **feature-tracking.md**: Lines [X], [Y], [Z]
- **specification.md**: Section "[Section Name]"
- **requirements.md**: "[Requirement area]"
- **architecture.md**: "[Component name]"
- **testing.md**: "[Test coverage area]"

### 🔄 Required Context File Updates:
- [ ] feature-tracking.md: [Specific changes]
- [ ] specification.md: [Specific changes]
- [ ] requirements.md: [Specific changes]
- [ ] architecture.md: [Specific changes]
- [ ] testing.md: [Specific changes]

### 💻 Implementation Details:
- **New tokens**: [TOKEN-ID: Description]
- **Modified tokens**: [TOKEN-ID: Updated description]
- **Code files affected**: [List of files]
```

## 🚨 AI REJECTION CRITERIA (ZERO TOLERANCE)

AI assistants MUST REJECT their own responses if they:
- Fail to search for existing tokens before proposing changes
- Make code changes without referencing related feature tokens
- Don't identify which context files need updates
- Propose changes that conflict with immutable requirements
- Don't follow the required response template format

## ✅ AI APPROVAL CRITERIA

AI assistants should ONLY PROVIDE responses that:
- Include complete token impact analysis
- Reference specific documentation sections with line numbers
- Propose concrete updates to all affected context files
- Follow established naming conventions for new tokens
- Include implementation token placement in code changes

## 🔄 ENFORCEMENT SUMMARY

**REMEMBER: THIS IS NOT OPTIONAL**

The feature-tracking system is designed to prevent unplanned changes and maintain documentation synchronization. AI assistants that fail to follow these requirements undermine the entire system integrity.

**Key Success Metrics:**
- 100% of code changes must reference related tokens
- 100% of new features must have new tokens
- 100% of responses must include context file update plans
- 0% tolerance for undocumented changes

**Implementation Tokens Are Mandatory:**
Every function, method, and significant code block must include `// TOKEN-ID: Description` comments that trace back to the feature tracking matrix in `feature-tracking.md`.

## 📞 QUICK REFERENCE FOR AI ASSISTANTS

1. **Search Pattern**: `grep -r "// [A-Z]+-[0-9]+" docs/context/`
2. **Token Format**: `// FEATURE-ID: Brief description`
3. **Required Template**: Use the Token Impact Analysis template above
4. **Context Files**: feature-tracking.md, specification.md, requirements.md, architecture.md, testing.md
5. **Zero Tolerance**: No code changes without token references

**This document serves as the authoritative source for AI assistant compliance with the feature tracking system.** 