#!/bin/bash

# [CRITICAL] TOKEN-002: Semantic token validation script - [ACTION:validation]
# Purpose: Validate semantic token system implementation and usage

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_FILE="${PROJECT_ROOT}/semantic-token-validation-report.md"

# [REQ:DOC_015] Semantic token patterns - Extended to include REQ/ARCH/IMPL tokens
PRIORITY_TOKENS=(
    "\[CRITICAL\]"
    "\[HIGH\]"
    "\[MEDIUM\]"
    "\[LOW\]"
)

ACTION_TOKENS=(
    "\[ACTION:core-functionality\]"
    "\[ACTION:format-processing\]"
    "\[ACTION:discovery\]"
    "\[ACTION:maintenance\]"
    "\[ACTION:validation\]"
    "\[ACTION:migration\]"
)

# [REQ:DOC_015] STDD semantic tokens (REQ/ARCH/IMPL)
STDD_TOKENS=(
    "\[REQ:[A-Z0-9_]+\]"
    "\[ARCH:[A-Z0-9_]+\]"
    "\[IMPL:[A-Z0-9_]+\]"
)

# Unicode icons that should be replaced
UNICODE_ICONS=(
    "⭐"
    "🔺"
    "🔶"
    "🔻"
    "🔧"
    "📝"
    "🔍"
    "🛠️"
    "🛡️"
    "🔄"
)

# Function to log messages
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}"
}

# Function to count occurrences
count_occurrences() {
    local pattern="$1"
    local file="$2"
    grep -o "${pattern}" "${file}" 2>/dev/null | wc -l || echo "0"
}

# Function to validate semantic tokens in a file
validate_file() {
    local file="$1"
    local validation_errors=0
    local validation_warnings=0
    
    log_message "VALIDATE" "Validating semantic tokens in ${file}"
    
    # Check for remaining Unicode icons
    for icon in "${UNICODE_ICONS[@]}"; do
        local count=$(count_occurrences "${icon}" "${file}")
        if [ "$count" -gt 0 ]; then
            echo -e "${RED}❌ Found ${count} instances of Unicode icon '${icon}' in ${file}${NC}"
            validation_errors=$((validation_errors + 1))
        fi
    done
    
    # Check for semantic priority tokens
    local priority_count=0
    for token in "${PRIORITY_TOKENS[@]}"; do
        local count=$(count_occurrences "${token}" "${file}")
        priority_count=$((priority_count + count))
    done
    
    if [ "$priority_count" -gt 0 ]; then
        echo -e "${GREEN}✅ Found ${priority_count} semantic priority tokens in ${file}${NC}"
    fi
    
    # Check for semantic action tokens
    local action_count=0
    for token in "${ACTION_TOKENS[@]}"; do
        local count=$(count_occurrences "${token}" "${file}")
        action_count=$((action_count + count))
    done
    
    if [ "$action_count" -gt 0 ]; then
        echo -e "${GREEN}✅ Found ${action_count} semantic action tokens in ${file}${NC}"
    fi
    
    # [REQ:DOC_015] Check for STDD semantic tokens (REQ/ARCH/IMPL)
    local stdd_count=0
    for token_pattern in "${STDD_TOKENS[@]}"; do
        local count=$(grep -oE "${token_pattern}" "${file}" 2>/dev/null | wc -l || echo "0")
        stdd_count=$((stdd_count + count))
    done
    
    if [ "$stdd_count" -gt 0 ]; then
        echo -e "${GREEN}✅ Found ${stdd_count} STDD semantic tokens (REQ/ARCH/IMPL) in ${file}${NC}"
    fi
    
    # [REQ:DOC_015] Check for proper STDD token format
    local malformed_stdd_tokens=$(grep -oE "\[(REQ|ARCH|IMPL):[^]]+\]" "${file}" 2>/dev/null | grep -vE "\[(REQ|ARCH|IMPL):[A-Z0-9_]+\]" | wc -l || echo "0")
    if [ "$malformed_stdd_tokens" -gt 0 ]; then
        echo -e "${RED}❌ Found ${malformed_stdd_tokens} malformed STDD tokens in ${file}${NC}"
        validation_errors=$((validation_errors + 1))
    fi
    
    # Check for proper token format (legacy)
    local malformed_tokens=$(grep -o "\[[A-Z]\+[^]]*\]" "${file}" 2>/dev/null | grep -v "\[CRITICAL\]\|\[HIGH\]\|\[MEDIUM\]\|\[LOW\]\|\[ACTION:" | grep -vE "\[(REQ|ARCH|IMPL):" | wc -l || echo "0")
    if [ "$malformed_tokens" -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Found ${malformed_tokens} potentially malformed tokens in ${file}${NC}"
        validation_warnings=$((validation_warnings + 1))
    fi
    
    return $validation_errors
}

# Function to generate validation report
generate_validation_report() {
    local total_files=0
    local files_with_semantic_tokens=0
    local files_with_unicode_icons=0
    local total_semantic_tokens=0
    local total_unicode_icons=0
    local validation_errors=0
    local validation_warnings=0
    
    # Collect statistics
    for file in $(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "node_modules" | grep -v ".git" | grep -v "backup-"); do
        total_files=$((total_files + 1))
        
        # Count semantic tokens (including STDD tokens)
        local semantic_count=0
        for token in "${PRIORITY_TOKENS[@]}" "${ACTION_TOKENS[@]}"; do
            local count=$(count_occurrences "${token}" "${file}")
            semantic_count=$((semantic_count + count))
        done
        # [REQ:DOC_015] Count STDD tokens
        for token_pattern in "${STDD_TOKENS[@]}"; do
            local count=$(grep -oE "${token_pattern}" "${file}" 2>/dev/null | wc -l || echo "0")
            semantic_count=$((semantic_count + count))
        done
        
        if [ "$semantic_count" -gt 0 ]; then
            files_with_semantic_tokens=$((files_with_semantic_tokens + 1))
            total_semantic_tokens=$((total_semantic_tokens + semantic_count))
        fi
        
        # Count Unicode icons
        local unicode_count=0
        for icon in "${UNICODE_ICONS[@]}"; do
            local count=$(count_occurrences "${icon}" "${file}")
            unicode_count=$((unicode_count + count))
        done
        
        if [ "$unicode_count" -gt 0 ]; then
            files_with_unicode_icons=$((files_with_unicode_icons + 1))
            total_unicode_icons=$((total_unicode_icons + unicode_count))
        fi
    done
    
    # Generate report
    cat > "${REPORT_FILE}" << EOF
# Semantic Token Validation Report

**Validation Date**: $(date '+%Y-%m-%d %H:%M:%S')
**Project**: BkpDir
**Validation Script**: validate-semantic-tokens.sh

## Validation Summary

### Statistics
- **Total files scanned**: ${total_files}
- **Files with semantic tokens**: ${files_with_semantic_tokens}
- **Files with Unicode icons**: ${files_with_unicode_icons}
- **Total semantic tokens found**: ${total_semantic_tokens}
- **Total Unicode icons found**: ${total_unicode_icons}

### Semantic Token Patterns

#### Priority Tokens
$(for token in "${PRIORITY_TOKENS[@]}"; do
    echo "- \`${token}\`"
done)

#### Action Tokens
$(for token in "${ACTION_TOKENS[@]}"; do
    echo "- \`${token}\`"
done)

#### STDD Semantic Tokens [REQ:DOC_015]
$(for token_pattern in "${STDD_TOKENS[@]}"; do
    echo "- \`${token_pattern}\` (REQ/ARCH/IMPL tokens)"
done)

### Unicode Icons to Replace
$(for icon in "${UNICODE_ICONS[@]}"; do
    echo "- \`${icon}\`"
done)

## Detailed File Analysis

### Files with Semantic Tokens
\`\`\`
$(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "node_modules" | grep -v ".git" | grep -v "backup-" | xargs grep -lE "\[CRITICAL\]|\[HIGH\]|\[MEDIUM\]|\[LOW\]|\[ACTION:|\[REQ:|\[ARCH:|\[IMPL:" 2>/dev/null || echo "No files with semantic tokens found")
\`\`\`

### Files with Remaining Unicode Icons
\`\`\`
$(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "node_modules" | grep -v ".git" | grep -v "backup-" | xargs grep -l "[⭐🔺🔶🔻🔧📝🔍🛠️🛡️🔄]" 2>/dev/null || echo "No files with remaining Unicode icons found")
\`\`\`

## Validation Results

### Success Criteria
- ✅ **Semantic tokens present**: ${files_with_semantic_tokens} files
- ✅ **Total semantic tokens**: ${total_semantic_tokens} tokens
- ❌ **Unicode icons remaining**: ${files_with_unicode_icons} files
- ❌ **Total Unicode icons**: ${total_unicode_icons} icons

### Recommendations

1. **Token Standardization**: Ensure all new code uses semantic tokens
2. **Documentation Update**: Update documentation to reference semantic tokens
3. **AI Assistant Integration**: Update AI assistant protocols to use semantic tokens
4. **Validation Automation**: Integrate validation into CI/CD pipeline
5. **Continuous Validation**: Add validation to development workflow

## Token Usage Guidelines

### Priority Tokens
- Use \`[CRITICAL]\` for essential system operations
- Use \`[HIGH]\` for important features
- Use \`[MEDIUM]\` for standard features
- Use \`[LOW]\` for minor improvements

### Action Tokens
- Use \`[ACTION:core-functionality]\` for essential operations
- Use \`[ACTION:format-processing]\` for text formatting
- Use \`[ACTION:discovery]\` for file system operations
- Use \`[ACTION:maintenance]\` for code cleanup
- Use \`[ACTION:validation]\` for error checking
- Use \`[ACTION:migration]\` for system migrations

## Compliance Status

**Overall Status**: $(if [ "$files_with_unicode_icons" -eq 0 ]; then echo "✅ COMPLIANT"; else echo "❌ NON-COMPLIANT"; fi)

**Migration Progress**: $(( (files_with_semantic_tokens * 100) / total_files ))%

**Remaining Work**: ${files_with_unicode_icons} files need migration
EOF

    log_message "REPORT" "Generated validation report: ${REPORT_FILE}"
    echo -e "${BLUE}📊 Validation report generated: ${REPORT_FILE}${NC}"
}

# Main validation process
main() {
    echo -e "${BLUE}🔍 Starting Semantic Token Validation${NC}"
    echo "Project Root: ${PROJECT_ROOT}"
    echo "Report File: ${REPORT_FILE}"
    echo ""
    
    local total_files=0
    local files_with_errors=0
    local files_with_warnings=0
    
    # Find all relevant files
    local files_to_validate=($(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "node_modules" | grep -v ".git" | grep -v "backup-"))
    
    echo -e "${BLUE}📋 Found ${#files_to_validate[@]} files to validate${NC}"
    echo ""
    
    # Validate each file
    for file in "${files_to_validate[@]}"; do
        total_files=$((total_files + 1))
        echo -e "${BLUE}[${total_files}/${#files_to_validate[@]}] Validating: ${file}${NC}"
        
        if ! validate_file "${file}"; then
            files_with_errors=$((files_with_errors + 1))
        fi
        
        echo ""
    done
    
    echo -e "${BLUE}📊 Validation Summary${NC}"
    echo -e "Total files validated: ${total_files}"
    echo -e "Files with errors: ${files_with_errors}"
    echo -e "Files with warnings: ${files_with_warnings}"
    echo -e "Report file: ${REPORT_FILE}"
    
    # Generate report
    generate_validation_report
    
    if [ "$files_with_errors" -eq 0 ]; then
        echo -e "${GREEN}🎉 Validation completed successfully!${NC}"
        exit 0
    else
        echo -e "${RED}⚠️  Validation completed with ${files_with_errors} files with errors${NC}"
        exit 1
    fi
}

# Run main function
main "$@" 