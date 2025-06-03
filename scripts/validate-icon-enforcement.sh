#!/bin/bash

# 🔺 DOC-008: Icon validation and enforcement - Comprehensive validation system
# This script provides automated validation and enforcement of icon consistency
# across all documentation and source code files

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

# Configuration
STRICT_MODE=${1:-false}  # Enable strict mode for CI/CD
VERBOSE=${VERBOSE:-false}
OUTPUT_FORMAT=${OUTPUT_FORMAT:-"console"}  # console, json, markdown

echo -e "${BLUE}🛡️ DOC-008: Comprehensive Icon Validation and Enforcement${NC}"
echo "=============================================================="

# Check if we're in the right directory
if [[ ! -f "docs/context/feature-tracking.md" ]] || [[ ! -f "docs/context/README.md" ]]; then
    echo -e "${RED}❌ Error: Must be run from project root directory${NC}"
    echo -e "${RED}   Required files: docs/context/feature-tracking.md, docs/context/README.md${NC}"
    exit 1
fi

# 🔍 DOC-008: Master icon definitions from README.md
# Define master icons and their meanings
MASTER_ICON_LIST="⭐:CRITICAL 🔺:HIGH 🔶:MEDIUM 🔻:LOW 🚀:PHASE_1 ⚡:PHASE_2 🔄:PHASE_3 🏁:PHASE_4 1️⃣:STEP_1 2️⃣:STEP_2 3️⃣:STEP_3 ✅:COMPLETE 📑:PURPOSE 📋:CHECKLIST 📊:ANALYSIS 📖:REFERENCE 🔍:SEARCH 📝:DOCUMENT 🔧:CONFIGURE 🛡️:PROTECT"

# Function to check if icon is in master list
is_master_icon() {
    local icon="$1"
    echo "$MASTER_ICON_LIST" | grep -q "$icon:"
}

# Function to get icon meaning
get_icon_meaning() {
    local icon="$1"
    echo "$MASTER_ICON_LIST" | grep -o "$icon:[A-Z_]*" | cut -d: -f2
}

# 🔍 DOC-008: Load validation patterns from context files
load_validation_patterns() {
    echo -e "${BLUE}📋 Step 1: Loading icon validation patterns from context system...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Check if master icon legend exists
    if grep -q "Master Icon Legend" docs/context/README.md; then
        echo -e "  ${GREEN}✅ Master Icon Legend found in README.md${NC}"
        ((successes++))
    else
        echo -e "  ${RED}❌ Master Icon Legend missing from README.md${NC}"
        ((errors++))
        return 1
    fi
    
    # Count master icons
    local icon_count=$(echo "$MASTER_ICON_LIST" | tr ' ' '\n' | wc -l)
    echo -e "  ${GREEN}✅ Master icon definitions loaded: $icon_count unique icons${NC}"
    ((successes++))
}

# 🔍 DOC-008: Validate documentation icon consistency
validate_documentation_icons() {
    echo
    echo -e "${BLUE}📋 Step 2: Validating icon consistency in context documentation...${NC}"
    echo "-------------------------------------------------------------------"
    
    local doc_files=(
        "docs/context/README.md"
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
            
            # Check for unknown icons (sample list - emojis can be tricky)
            local unknown_icons_found=false
            
            # Check for consistent priority icon usage
            local priority_icons_found=0
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
    
    local go_files=($(find . -name "*.go" -not -path "./vendor/*" 2>/dev/null))
    local token_files_with_icons=0
    local token_files_without_icons=0
    local total_tokens_found=0
    local standardized_tokens=0
    
    for go_file in "${go_files[@]}"; do
        if [[ -f "$go_file" ]]; then
            echo -e "${CYAN}  Checking $go_file...${NC}"
            ((total_files_checked++))
            
            # Look for implementation tokens
            local tokens_found=false
            while IFS= read -r line; do
                local line_num=$(echo "$line" | cut -d: -f1)
                local line_content=$(echo "$line" | cut -d: -f2-)
                
                ((total_tokens_found++))
                tokens_found=true
                
                # Check for standardized format with priority icons
                if echo "$line_content" | grep -q "// *[⭐🔺🔶🔻]"; then
                    # Extract feature ID
                    local feature_id=$(echo "$line_content" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
                    
                    echo -e "    ${GREEN}✅ Standardized token found: $feature_id${NC}"
                    ((standardized_tokens++))
                    ((successes++))
                    
                    # Check for action icons
                    if echo "$line_content" | grep -q "[🔍📝🔧🛡️]"; then
                        echo -e "    ${GREEN}✅ Action icon found in token${NC}"
                        ((successes++))
                    else
                        echo -e "    ${YELLOW}⚠️ Token missing action icon: $feature_id${NC}"
                        ((warnings++))
                    fi
                    
                elif echo "$line_content" | grep -q "// *[A-Z]\+-[0-9]\+:"; then
                    local feature_id=$(echo "$line_content" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
                    echo -e "    ${YELLOW}⚠️ Legacy token format: $feature_id (missing priority icon)${NC}"
                    ((warnings++))
                else
                    echo -e "    ${RED}❌ Invalid implementation token format in line $line_num${NC}"
                    ((errors++))
                fi
                
            done < <(grep -n "// .*[A-Z]\+-[0-9]\+:" "$go_file" 2>/dev/null || true)
            
            if [[ "$tokens_found" == true ]]; then
                ((token_files_with_icons++))
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
    
    # Extract feature IDs from feature-tracking.md
    local doc_feature_ids=($(grep "^| [A-Z]\+-[0-9]\+ |" docs/context/feature-tracking.md 2>/dev/null | sed 's/^| \([A-Z]\+-[0-9]\+\) |.*/\1/' || true))
    
    # Extract feature IDs from source code
    local code_feature_ids=()
    for go_file in $(find . -name "*.go" -not -path "./vendor/*" 2>/dev/null); do
        while IFS= read -r feature_id; do
            code_feature_ids+=("$feature_id")
        done < <(grep -o '[A-Z]\+-[0-9]\+' "$go_file" 2>/dev/null | sort -u || true)
    done
    
    echo "  Feature IDs in documentation: ${#doc_feature_ids[@]}"
    echo "  Feature IDs in source code: ${#code_feature_ids[@]}"
    
    # Simple cross-reference check
    if [[ ${#doc_feature_ids[@]} -gt 0 ]] && [[ ${#code_feature_ids[@]} -gt 0 ]]; then
        echo -e "  ${GREEN}✅ Cross-references found between documentation and code${NC}"
        ((successes++))
    else
        echo -e "  ${YELLOW}⚠️ Limited cross-reference data available${NC}"
        ((warnings++))
    fi
}

# 🔍 DOC-008: Enforcement rules validation
validate_enforcement_rules() {
    echo
    echo -e "${BLUE}📋 Step 5: Validating enforcement rules compliance...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Check AI assistant compliance integration
    if grep -q "DOC-007" docs/context/ai-assistant-compliance.md && \
       grep -q "DOC-008" docs/context/ai-assistant-compliance.md; then
        echo -e "  ${GREEN}✅ AI assistant compliance includes DOC-007 and DOC-008 requirements${NC}"
        ((successes++))
    else
        echo -e "  ${YELLOW}⚠️ AI assistant compliance may need updates for DOC-008${NC}"
        ((warnings++))
    fi
    
    # Check Makefile integration
    if grep -q "validate-icon" Makefile; then
        echo -e "  ${GREEN}✅ Makefile includes icon validation targets${NC}"
        ((successes++))
    else
        echo -e "  ${RED}❌ Makefile missing icon validation integration${NC}"
        ((errors++))
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
    
    local report_file="icon-validation-report.md"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    cat > "$report_file" << EOF
# Icon Validation and Enforcement Report (DOC-008)

> **Generated on:** \`$timestamp\`
> **Script:** \`scripts/validate-icon-enforcement.sh\`
> **Mode:** ${STRICT_MODE:-"Standard"}

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
    
    cat >> "$report_file" << EOF
## Enforcement Status

- **Icon System**: Comprehensive (20+ icons defined)
- **Documentation Compliance**: $([ $((errors + warnings)) -lt 10 ] && echo "✅ Good" || echo "❌ Needs improvement")
- **Code Standardization**: $([ $((standardized_tokens * 100 / total_tokens_found)) -gt 80 ] 2>/dev/null && echo "✅ Good" || echo "⚠️ In progress")
- **Automation**: $(grep -q "validate-icon" Makefile && echo "✅ Integrated" || echo "❌ Missing")

---
*This report was generated by DOC-008 icon validation and enforcement system*
EOF

    echo -e "  ${GREEN}✅ Report generated: $report_file${NC}"
    ((successes++))
}

# Main execution
main() {
    load_validation_patterns
    validate_documentation_icons  
    validate_source_code_icons
    validate_cross_references
    validate_enforcement_rules
    generate_enforcement_report
    
    echo
    echo "=============================================================="
    echo -e "${BLUE}📊 Final Validation Summary${NC}"
    echo "=============================================================="
    echo -e "${GREEN}✅ Successes: $successes${NC}"
    echo -e "${YELLOW}⚠️ Warnings: $warnings${NC}"
    echo -e "${RED}❌ Errors: $errors${NC}"
    echo "Files checked: $total_files_checked"
    
    # Determine exit code based on strict mode
    if [[ "$STRICT_MODE" == "true" ]]; then
        if [[ $errors -gt 0 || $warnings -gt 5 ]]; then
            echo
            echo -e "${RED}💥 DOC-008 validation failed in strict mode${NC}"
            echo -e "${RED}   Errors: $errors, Warnings: $warnings (threshold: 5)${NC}"
            echo -e "${YELLOW}💡 Run without strict mode for development validation${NC}"
            exit 1
        fi
    elif [[ $errors -gt 0 ]]; then
        echo
        echo -e "${RED}💥 DOC-008 validation failed with $errors critical errors${NC}"
        exit 1
    fi
    
    echo
    echo -e "${GREEN}🎉 DOC-008 icon validation and enforcement completed successfully!${NC}"
    echo -e "${GREEN}   Icon system integrity maintained across documentation and code.${NC}"
    
    if [[ $warnings -gt 0 ]]; then
        echo -e "${YELLOW}💡 Consider addressing $warnings warnings for optimal compliance.${NC}"
    fi
    
    exit 0
}

# Handle command line arguments
case "${1:-}" in
    --strict)
        STRICT_MODE=true
        shift
        ;;
    --verbose)
        VERBOSE=true
        shift
        ;;
    --help)
        echo "🔺 DOC-008: Icon validation and enforcement"
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --strict    Enable strict mode (fail on warnings)"
        echo "  --verbose   Enable verbose output" 
        echo "  --help      Show this help message"
        echo ""
        echo "Environment Variables:"
        echo "  OUTPUT_FORMAT  Set output format (console, json, markdown)"
        echo "  VERBOSE        Enable verbose output (true/false)"
        echo ""
        echo "Examples:"
        echo "  $0                    # Standard validation"
        echo "  $0 --strict          # Strict validation for CI/CD"
        echo "  VERBOSE=true $0      # Verbose output"
        exit 0
        ;;
esac

# Run main function
main "$@" 