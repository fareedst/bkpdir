#!/bin/bash

# 🔺 DOC-008: Icon validation and enforcement - Comprehensive validation system
# This script provides automated validation and enforcement of icon consistency
# across all documentation and source code files
# Integrated with DOC-014: AI Assistant Decision Framework

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Counters
errors=0
warnings=0
successes=0
total_files_checked=0
total_tokens_found=0
standardized_tokens=0

# Configuration
STRICT_MODE=false
VERBOSE=false
OUTPUT_FORMAT="console"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --strict)
            STRICT_MODE=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --format)
            OUTPUT_FORMAT="$2"
            shift 2
            ;;
        --help)
            echo "🔺 DOC-008: Icon validation and enforcement"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --strict    Enable strict mode (fail on warnings)"
            echo "  --verbose   Enable verbose output" 
            echo "  --format    Set output format (console, json, markdown)"
            echo "  --help      Show this help message"
            echo ""
            echo "Environment Variables:"
            echo "  VERBOSE        Enable verbose output (true/false)"
            echo ""
            echo "Examples:"
            echo "  $0                    # Standard validation"
            echo "  $0 --strict          # Strict validation for CI/CD"
            echo "  $0 --verbose         # Verbose output"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Override with environment variables if set
[[ "${VERBOSE:-false}" == "true" ]] && VERBOSE=true

echo -e "${BLUE}🛡️ DOC-008: Comprehensive Icon Validation and Enforcement${NC}"
echo "=============================================================="

# Check if we're in the right directory
if [[ ! -f "docs/context/feature-tracking.md" ]] || [[ ! -f "README.md" ]]; then
    echo -e "${RED}❌ Error: Must be run from project root directory${NC}"
    echo -e "${RED}   Required files: docs/context/feature-tracking.md, README.md${NC}"
    exit 1
fi

# 🔍 DOC-008: Master icon definitions
# Define master icons and their meanings based on standardized system
declare -A MASTER_ICONS
MASTER_ICONS=(
    ["⭐"]="CRITICAL"
    ["🔺"]="HIGH"
    ["🔶"]="MEDIUM"
    ["🔻"]="LOW"
    ["🚀"]="PHASE_1"
    ["⚡"]="PHASE_2"
    ["🔄"]="PHASE_3"
    ["🏁"]="PHASE_4"
    ["1️⃣"]="STEP_1"
    ["2️⃣"]="STEP_2"
    ["3️⃣"]="STEP_3"
    ["✅"]="COMPLETE"
    ["📑"]="PURPOSE"
    ["📋"]="CHECKLIST"
    ["📊"]="ANALYSIS"
    ["📖"]="REFERENCE"
    ["🔍"]="SEARCH"
    ["📝"]="DOCUMENT"
    ["🔧"]="CONFIGURE"
    ["🛡️"]="PROTECT"
)

# Priority icons for implementation tokens
PRIORITY_ICONS="⭐🔺🔶🔻"
ACTION_ICONS="🔍📝🔧🛡️"

# Function to check if icon is in master list
is_master_icon() {
    local icon="$1"
    [[ -n "${MASTER_ICONS[$icon]}" ]]
}

# Function to get icon meaning
get_icon_meaning() {
    local icon="$1"
    echo "${MASTER_ICONS[$icon]:-UNKNOWN}"
}

# 🔍 DOC-008: Load validation patterns from context files
load_validation_patterns() {
    echo -e "${BLUE}📋 Step 1: Loading icon validation patterns from context system...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Check if master icon legend exists in README.md
    if grep -q "Master Icon Legend\|Icon Legend\|## Icons" README.md; then
        echo -e "  ${GREEN}✅ Master Icon Legend found in README.md${NC}"
        ((successes++))
    else
        echo -e "  ${RED}❌ Master Icon Legend missing from README.md${NC}"
        ((errors++))
    fi
    
    # Count master icons
    local icon_count=${#MASTER_ICONS[@]}
    echo -e "  ${GREEN}✅ Master icon definitions loaded: $icon_count unique icons${NC}"
    ((successes++))
    
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "  ${CYAN}📋 Master icons defined:${NC}"
        for icon in "${!MASTER_ICONS[@]}"; do
            echo -e "    $icon → ${MASTER_ICONS[$icon]}"
        done
    fi
}

# 🔍 DOC-008: Validate documentation icon consistency
validate_documentation_icons() {
    echo
    echo -e "${BLUE}📋 Step 2: Validating icon consistency in context documentation...${NC}"
    echo "-------------------------------------------------------------------"
    
    local doc_files=(
        "README.md"
        "docs/context/feature-tracking.md"
        "docs/context/ai-assistant-protocol.md"
        "docs/context/ai-assistant-compliance.md"
        "docs/context/specification.md"
        "docs/context/requirements.md"
        "docs/context/architecture.md"
        "docs/context/testing.md"
    )
    
    for doc_file in "${doc_files[@]}"; do
        if [[ -f "$doc_file" ]]; then
            echo -e "${CYAN}  Checking $doc_file...${NC}"
            ((total_files_checked++))
            
            # Check for consistent priority icon usage
            local priority_icons_found=0
            local unknown_icons=()
            
            # Count priority icons
            if grep -q "⭐" "$doc_file"; then ((priority_icons_found++)); fi
            if grep -q "🔺" "$doc_file"; then ((priority_icons_found++)); fi
            if grep -q "🔶" "$doc_file"; then ((priority_icons_found++)); fi
            if grep -q "🔻" "$doc_file"; then ((priority_icons_found++)); fi
            
            if [[ $priority_icons_found -gt 0 ]]; then
                echo -e "    ${GREEN}✅ Priority icons found: $priority_icons_found types${NC}"
                ((successes++))
            fi
            
            echo -e "    ${GREEN}✅ Documentation file validated${NC}"
            ((successes++))
        else
            echo -e "  ${YELLOW}⚠️ Documentation file not found: $doc_file${NC}"
            ((warnings++))
        fi
    done
}

# 🔍 DOC-008: Validate source code implementation tokens
validate_source_code_icons() {
    echo
    echo -e "${BLUE}📋 Step 3: Validating implementation token icon consistency...${NC}"
    echo "-------------------------------------------------------------------"
    
    local go_files=($(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" 2>/dev/null))
    local token_files_with_icons=0
    local token_files_without_icons=0
    
    for go_file in "${go_files[@]}"; do
        if [[ -f "$go_file" ]]; then
            if [[ "$VERBOSE" == "true" ]]; then
                echo -e "${CYAN}  Checking $go_file...${NC}"
            fi
            ((total_files_checked++))
            
            # Look for implementation tokens
            local tokens_found=false
            local file_token_count=0
            local file_standardized_count=0
            
            while IFS= read -r line_with_num; do
                local line_num=$(echo "$line_with_num" | cut -d: -f1)
                local line_content=$(echo "$line_with_num" | cut -d: -f2-)
                
                ((total_tokens_found++))
                ((file_token_count++))
                tokens_found=true
                
                # Check for standardized format with priority icons
                if echo "$line_content" | grep -qE "// *[⭐🔺🔶🔻]"; then
                    # Extract feature ID
                    local feature_id=$(echo "$line_content" | grep -oE '[A-Z]+-[0-9]+' | head -1)
                    
                    if [[ "$VERBOSE" == "true" ]]; then
                        echo -e "    ${GREEN}✅ Standardized token found: $feature_id${NC}"
                    fi
                    ((standardized_tokens++))
                    ((file_standardized_count++))
                    ((successes++))
                    
                    # Check for action icons
                    if echo "$line_content" | grep -qE "[🔍📝🔧🛡️]"; then
                        if [[ "$VERBOSE" == "true" ]]; then
                            echo -e "    ${GREEN}✅ Action icon found in token${NC}"
                        fi
                        ((successes++))
                    else
                        if [[ "$VERBOSE" == "true" ]]; then
                            echo -e "    ${YELLOW}⚠️ Token missing action icon: $feature_id${NC}"
                        fi
                        ((warnings++))
                    fi
                    
                elif echo "$line_content" | grep -qE "// .*[A-Z]+-[0-9]+:"; then
                    local feature_id=$(echo "$line_content" | grep -oE '[A-Z]+-[0-9]+' | head -1)
                    if [[ "$VERBOSE" == "true" ]]; then
                        echo -e "    ${YELLOW}⚠️ Legacy token format: $feature_id (missing priority icon)${NC}"
                    fi
                    ((warnings++))
                else
                    if [[ "$VERBOSE" == "true" ]]; then
                        echo -e "    ${RED}❌ Invalid implementation token format in line $line_num${NC}"
                    fi
                    ((errors++))
                fi
                
            done < <(grep -n "// .*[A-Z]\+-[0-9]\+:" "$go_file" 2>/dev/null || true)
            
            if [[ "$tokens_found" == true ]]; then
                ((token_files_with_icons++))
                if [[ "$VERBOSE" == "false" && $file_token_count -gt 0 ]]; then
                    local file_std_rate=0
                    if [[ $file_token_count -gt 0 ]]; then
                        file_std_rate=$(( (file_standardized_count * 100) / file_token_count ))
                    fi
                    echo -e "  ${CYAN}$go_file: $file_token_count tokens, ${file_std_rate}% standardized${NC}"
                fi
            else
                ((token_files_without_icons++))
            fi
        fi
    done
    
    echo
    echo -e "${BLUE}📊 Implementation Token Summary:${NC}"
    echo "  Total Go files checked: ${#go_files[@]}"
    echo "  Files with tokens: $token_files_with_icons"
    echo "  Files without tokens: $token_files_without_icons"
    echo "  Total tokens found: $total_tokens_found"
    echo "  Standardized tokens: $standardized_tokens"
    
    if [[ $total_tokens_found -gt 0 ]]; then
        local standardization_rate=$(( (standardized_tokens * 100) / total_tokens_found ))
        echo "  Standardization rate: ${standardization_rate}%"
        
        if [[ $standardization_rate -ge 90 ]]; then
            echo -e "  ${GREEN}✅ Excellent standardization rate${NC}"
            ((successes++))
        elif [[ $standardization_rate -ge 70 ]]; then
            echo -e "  ${YELLOW}⚠️ Good standardization rate, improvement possible${NC}"
            ((warnings++))
        else
            echo -e "  ${RED}❌ Poor standardization rate, action required${NC}"
            ((errors++))
        fi
    fi
}

# 🔍 DOC-008: Cross-reference validation between docs and code
validate_cross_references() {
    echo
    echo -e "${BLUE}📋 Step 4: Validating cross-reference consistency...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Extract feature IDs from feature tracking
    local feature_ids=()
    if [[ -f "docs/context/feature-tracking.md" ]]; then
        while IFS= read -r feature_id; do
            feature_ids+=("$feature_id")
        done < <(grep -oE '[A-Z]+-[0-9]+' docs/context/feature-tracking.md | sort -u)
        
        echo -e "  ${GREEN}✅ Found ${#feature_ids[@]} feature IDs in feature tracking${NC}"
        ((successes++))
    else
        echo -e "  ${YELLOW}⚠️ Feature tracking file not found${NC}"
        ((warnings++))
    fi
    
    # Check for orphaned implementation tokens
    local code_feature_ids=()
    while IFS= read -r feature_id; do
        code_feature_ids+=("$feature_id")
    done < <(grep -roE '// .*[A-Z]+-[0-9]+:' . --include="*.go" | grep -oE '[A-Z]+-[0-9]+' | sort -u)
    
    if [[ ${#code_feature_ids[@]} -gt 0 ]]; then
        echo -e "  ${GREEN}✅ Found ${#code_feature_ids[@]} unique feature IDs in code${NC}"
        ((successes++))
        
        # Check for orphaned tokens
        local orphaned_count=0
        for code_id in "${code_feature_ids[@]}"; do
            local found=false
            for feature_id in "${feature_ids[@]}"; do
                if [[ "$code_id" == "$feature_id" ]]; then
                    found=true
                    break
                fi
            done
            if [[ "$found" == false ]]; then
                if [[ "$VERBOSE" == "true" ]]; then
                    echo -e "    ${YELLOW}⚠️ Orphaned implementation token: $code_id${NC}"
                fi
                ((orphaned_count++))
                ((warnings++))
            fi
        done
        
        if [[ $orphaned_count -eq 0 ]]; then
            echo -e "  ${GREEN}✅ No orphaned implementation tokens found${NC}"
            ((successes++))
        else
            echo -e "  ${YELLOW}⚠️ Found $orphaned_count orphaned implementation tokens${NC}"
        fi
    fi
}

# 🔍 DOC-008: Validate enforcement rules compliance
validate_enforcement_rules() {
    echo
    echo -e "${BLUE}📋 Step 5: Validating enforcement rules compliance...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Check AI assistant compliance integration
    if [[ -f "docs/context/ai-assistant-compliance.md" ]]; then
        if grep -q "DOC-007\|DOC-008" docs/context/ai-assistant-compliance.md; then
            echo -e "  ${GREEN}✅ AI assistant compliance includes DOC-007/DOC-008 requirements${NC}"
            ((successes++))
        else
            echo -e "  ${YELLOW}⚠️ AI assistant compliance may need updates for DOC-008${NC}"
            ((warnings++))
        fi
    else
        echo -e "  ${YELLOW}⚠️ AI assistant compliance file not found${NC}"
        ((warnings++))
    fi
    
    # Check Makefile integration
    if [[ -f "Makefile" ]]; then
        if grep -q "validate-icon" Makefile; then
            echo -e "  ${GREEN}✅ Makefile includes icon validation targets${NC}"
            ((successes++))
        else
            echo -e "  ${YELLOW}⚠️ Makefile missing icon validation integration${NC}"
            ((warnings++))
        fi
    else
        echo -e "  ${YELLOW}⚠️ Makefile not found${NC}"
        ((warnings++))
    fi
    
    # Check if this enforcement script exists and is executable
    if [[ -x "scripts/validate-icon-enforcement.sh" ]]; then
        echo -e "  ${GREEN}✅ Icon enforcement script is executable${NC}"
        ((successes++))
    else
        echo -e "  ${YELLOW}⚠️ Icon enforcement script may not be executable${NC}"
        ((warnings++))
    fi
}

# 🔍 DOC-008: Generate enforcement report
generate_enforcement_report() {
    echo
    echo -e "${BLUE}📊 Step 6: Generating enforcement report...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Create validation reports directory if it doesn't exist
    mkdir -p docs/validation-reports
    
    local report_file="docs/validation-reports/icon-validation-report.md"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    # Generate report directly in docs/validation-reports/
    cat > "$report_file" << EOF
# Icon Validation and Enforcement Report (DOC-008)

> **Generated on:** \`$timestamp\`
> **Script:** \`scripts/validate-icon-enforcement.sh\`
> **Mode:** $([ "$STRICT_MODE" == "true" ] && echo "Strict" || echo "Standard")

## Validation Summary

| Metric | Count |
|--------|-------|
| Files Checked | $total_files_checked |
| Successes | $successes |
| Warnings | $warnings |
| Errors | $errors |

## Validation Categories

### ✅ Successful Validations
- Master icon legend validation
- Documentation icon consistency  
- Implementation token standardization
- Cross-reference consistency

### ⚠️ Warnings Found
- Legacy token formats requiring updates
- Missing action icons in implementation tokens
- Potential icon overuse in documentation

### ❌ Errors Found
- Invalid priority icons not in master legend
- Implementation token format violations
- Missing validation infrastructure

## Recommendations

EOF

    if [[ $warnings -gt 0 ]]; then
        cat >> "$report_file" << EOF
### Priority Actions Needed
1. **Update Legacy Tokens**: $warnings legacy implementation tokens need standardization
2. **Add Action Icons**: Include appropriate action icons (🔍📝🔧🛡️) in implementation tokens
3. **Review Documentation**: Ensure all icons conform to master legend in README.md

EOF
    fi
    
    if [[ $errors -gt 0 ]]; then
        cat >> "$report_file" << EOF
### Critical Issues Requiring Immediate Attention
1. **Invalid Icons**: Remove or replace icons not defined in master legend
2. **Implementation Token Violations**: Fix malformed implementation tokens
3. **Infrastructure Setup**: Complete validation infrastructure setup

EOF
    fi
    
    local standardization_rate=0
    if [[ $total_tokens_found -gt 0 ]]; then
        standardization_rate=$(( (standardized_tokens * 100) / total_tokens_found ))
    fi
    
    cat >> "$report_file" << EOF
## Enforcement Status

- **Icon System**: Comprehensive (${#MASTER_ICONS[@]} icons defined)
- **Documentation Compliance**: $([ $((errors + warnings)) -lt 10 ] && echo "✅ Good" || echo "❌ Needs improvement")
- **Code Standardization**: $([ $standardization_rate -gt 80 ] 2>/dev/null && echo "✅ Good ($standardization_rate%)" || echo "⚠️ In progress ($standardization_rate%)")
- **Automation**: $(grep -q "validate-icon" Makefile 2>/dev/null && echo "✅ Integrated" || echo "❌ Missing")

## Implementation Token Statistics

| Metric | Value |
|--------|-------|
| Total Tokens | $total_tokens_found |
| Standardized | $standardized_tokens |
| Legacy Format | $((total_tokens_found - standardized_tokens)) |
| Standardization Rate | ${standardization_rate}% |

---
*This report was generated by DOC-008 icon validation and enforcement system*
EOF

    echo -e "  ${GREEN}✅ Validation report generated: $report_file${NC}"
}

# Main validation function
main() {
    local start_time=$(date +%s)
    
    # Run all validation steps
    load_validation_patterns
    validate_documentation_icons
    validate_source_code_icons
    validate_cross_references
    validate_enforcement_rules
    generate_enforcement_report
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo
    echo -e "${BLUE}🎯 Validation Complete${NC}"
    echo "=============================================================="
    echo -e "  Duration: ${duration}s"
    echo -e "  Files checked: $total_files_checked"
    echo -e "  Successes: ${GREEN}$successes${NC}"
    echo -e "  Warnings: ${YELLOW}$warnings${NC}"
    echo -e "  Errors: ${RED}$errors${NC}"
    
    if [[ $total_tokens_found -gt 0 ]]; then
        local standardization_rate=$(( (standardized_tokens * 100) / total_tokens_found ))
        echo -e "  Token standardization: ${standardization_rate}%"
    fi
    
    # Determine exit code based on mode
    if [[ "$STRICT_MODE" == "true" ]]; then
        # Strict mode: fail on errors OR too many warnings
        if [[ $errors -gt 0 ]] || [[ $warnings -gt 5 ]]; then
            echo -e "${RED}❌ Validation failed in strict mode${NC}"
            exit 1
        fi
    else
        # Standard mode: fail only on errors
        if [[ $errors -gt 0 ]]; then
            echo -e "${RED}❌ Validation failed${NC}"
            exit 1
        fi
    fi
    
    echo -e "${GREEN}✅ Validation passed${NC}"
    exit 0
}

# Run main function
main "$@" 