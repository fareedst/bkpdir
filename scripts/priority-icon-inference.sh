#!/bin/bash

# 🔺 DOC-009: Priority icon inference from feature-tracking.md - 🔍 Feature tracking analysis
# This script analyzes the feature tracking matrix to automatically assign priority icons

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

FEATURE_TRACKING_FILE="docs/context/feature-tracking.md"
OUTPUT_FORMAT=${1:-"console"}  # console, json, mapping

echo -e "${BLUE}🔺 DOC-009: Priority Icon Inference System${NC}"
echo "=============================================================="

# Check if feature tracking file exists
if [[ ! -f "$FEATURE_TRACKING_FILE" ]]; then
    echo -e "${RED}❌ Error: Feature tracking file not found: $FEATURE_TRACKING_FILE${NC}"
    exit 1
fi

# 🔺 DOC-009: Extract feature priorities from feature-tracking.md - 🔍 Priority extraction
extract_feature_priorities() {
    echo -e "${BLUE}📋 Extracting feature priorities from $FEATURE_TRACKING_FILE...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Extract feature table entries with priority icons
    local feature_count=0
    
    while IFS='|' read -r feature_id_col spec_col req_col arch_col test_col status_col token_col priority_col rest; do
        # Skip header rows and empty lines
        [[ "$feature_id_col" =~ ^[[:space:]]*Feature[[:space:]]*ID ]] && continue
        [[ "$feature_id_col" =~ ^[[:space:]]*-+ ]] && continue
        [[ -z "${feature_id_col// /}" ]] && continue
        
        # Extract feature ID (remove leading/trailing whitespace and markdown formatting)
        local feature_id=$(echo "$feature_id_col" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//;s/\*\*//g')
        
        # Skip if not a valid feature ID format
        [[ ! "$feature_id" =~ ^[A-Z]+-[0-9]+$ ]] && continue
        
        # Extract priority icon (look for priority icon in the last column)
        local priority_icon=""
        if [[ "$priority_col" =~ (⭐|🔺|🔶|🔻) ]]; then
            priority_icon=$(echo "$priority_col" | grep -o '[⭐🔺🔶🔻]' | head -1)
        elif [[ "$rest" =~ (⭐|🔺|🔶|🔻) ]]; then
            # Sometimes the priority is in the extended columns
            priority_icon=$(echo "$rest" | grep -o '[⭐🔺🔶🔻]' | head -1)
        fi
        
        # Default priority based on feature prefix if not found
        if [[ -z "$priority_icon" ]]; then
            local feature_prefix=$(echo "$feature_id" | cut -d'-' -f1)
            case "$feature_prefix" in
                "ARCH"|"FILE") priority_icon="⭐" ;;  # CRITICAL
                "CFG"|"GIT"|"TEST"|"DOC"|"LINT"|"COV") priority_icon="🔺" ;;  # HIGH
                "OUT"|"REFACTOR") priority_icon="🔶" ;;  # MEDIUM
                *) priority_icon="🔶" ;;  # Default MEDIUM
            esac
        fi
        
        # Store the mapping
        echo "$feature_id:$priority_icon"
        ((feature_count++))
        
    done < <(grep -E '^\|[[:space:]]*[A-Z]+-[0-9]+' "$FEATURE_TRACKING_FILE")
    
    echo -e "  ${GREEN}✅ Extracted priority mappings for $feature_count features${NC}"
}

# 🔺 DOC-009: Generate priority mapping file - 📝 Mapping file generation
generate_priority_mapping() {
    local output_file="priority-icon-mapping.txt"
    echo -e "${BLUE}📋 Generating priority mapping file: $output_file${NC}"
    echo "-------------------------------------------------------------------"
    
    # Create mapping file with header
    cat > "$output_file" << 'EOF'
# 🔺 DOC-009: Priority Icon Mapping for Implementation Tokens
# Format: FEATURE-ID:PRIORITY_ICON
# Generated from feature-tracking.md analysis
# 
# Priority Icons:
# ⭐ - CRITICAL (Core operations, blocking functionality)
# 🔺 - HIGH (Important features, significant configuration)
# 🔶 - MEDIUM (Secondary features, optional enhancements)
# 🔻 - LOW (Cleanup, optimization, debug utilities)
#
EOF
    
    # Add feature mappings
    extract_feature_priorities >> "$output_file"
    
    echo -e "  ${GREEN}✅ Priority mapping saved to: $output_file${NC}"
    echo -e "  ${CYAN}💡 Use this file with token migration scripts for consistent priority assignment${NC}"
}

# 🔺 DOC-009: Validate existing implementation tokens against mappings - 🛡️ Priority validation
validate_token_priorities() {
    echo -e "${BLUE}📋 Validating existing token priorities against feature tracking...${NC}"
    echo "-------------------------------------------------------------------"
    
    local correct_mappings=0
    local incorrect_mappings=0
    local missing_mappings=0
    
    # Create temporary mapping file
    local temp_mapping=$(mktemp)
    extract_feature_priorities > "$temp_mapping"
    
    # Find all implementation tokens in Go files
    while IFS= read -r line; do
        local file=$(echo "$line" | cut -d: -f1)
        local line_num=$(echo "$line" | cut -d: -f2)
        local line_content=$(echo "$line" | cut -d: -f3-)
        
        # Extract feature ID from token
        local feature_id=$(echo "$line_content" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
        [[ -z "$feature_id" ]] && continue
        
        # Extract current priority icon if present
        local current_icon=""
        if echo "$line_content" | grep -q "// *[⭐🔺🔶🔻]"; then
            current_icon=$(echo "$line_content" | grep -o '[⭐🔺🔶🔻]' | head -1)
        fi
        
        # Get expected priority from mapping
        local expected_icon=$(grep "^$feature_id:" "$temp_mapping" | cut -d: -f2)
        
        if [[ -z "$expected_icon" ]]; then
            echo -e "    ${YELLOW}⚠️ No mapping found for $feature_id in $file:$line_num${NC}"
            ((missing_mappings++))
        elif [[ -z "$current_icon" ]]; then
            echo -e "    ${YELLOW}⚠️ Missing priority icon for $feature_id in $file:$line_num (should be $expected_icon)${NC}"
            ((incorrect_mappings++))
        elif [[ "$current_icon" == "$expected_icon" ]]; then
            ((correct_mappings++))
        else
            echo -e "    ${RED}❌ Incorrect priority for $feature_id in $file:$line_num (is $current_icon, should be $expected_icon)${NC}"
            ((incorrect_mappings++))
        fi
        
    done < <(grep -rn "// [A-Z]\+-[0-9]\+:" . --include="*.go" 2>/dev/null || true)
    
    # Clean up
    rm -f "$temp_mapping"
    
    echo
    echo -e "${BLUE}📊 Priority Validation Summary:${NC}"
    echo "  Correct priorities: $correct_mappings"
    echo "  Incorrect/missing priorities: $incorrect_mappings"
    echo "  Features without mappings: $missing_mappings"
    
    local total=$((correct_mappings + incorrect_mappings + missing_mappings))
    if [[ $total -gt 0 ]]; then
        local accuracy=$((correct_mappings * 100 / total))
        echo "  Priority accuracy: $accuracy%"
        
        if [[ $accuracy -ge 90 ]]; then
            echo -e "  ${GREEN}✅ Excellent priority accuracy${NC}"
        elif [[ $accuracy -ge 70 ]]; then
            echo -e "  ${YELLOW}⚠️ Good priority accuracy, some corrections needed${NC}"
        else
            echo -e "  ${RED}❌ Poor priority accuracy, mass correction recommended${NC}"
        fi
    fi
}

# 🔺 DOC-009: Generate action icon suggestions - 🔍 Action icon analysis
suggest_action_icons() {
    echo -e "${BLUE}📋 Analyzing function context for action icon suggestions...${NC}"
    echo "-------------------------------------------------------------------"
    
    local suggestions=0
    
    # Find implementation tokens and analyze their context
    while IFS= read -r line; do
        local file=$(echo "$line" | cut -d: -f1)
        local line_num=$(echo "$line" | cut -d: -f2)
        local line_content=$(echo "$line" | cut -d: -f3-)
        
        # Extract feature ID
        local feature_id=$(echo "$line_content" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
        [[ -z "$feature_id" ]] && continue
        
        # Get surrounding context (5 lines before and after)
        local context=$(sed -n "$((line_num-5)),$((line_num+5))p" "$file" 2>/dev/null || echo "")
        
        # Suggest action icon based on context analysis
        local suggested_icon="🔧"  # Default
        local reason="Default configuration/modification"
        
        if echo "$context" | grep -qi "get\|find\|search\|discover\|detect\|analyze\|check\|validate\|parse"; then
            suggested_icon="🔍"
            reason="Search/discovery pattern detected"
        elif echo "$context" | grep -qi "format\|print\|write\|update\|log\|output\|render\|display"; then
            suggested_icon="📝"
            reason="Documentation/output pattern detected"
        elif echo "$context" | grep -qi "protect\|secure\|validate\|verify\|guard\|ensure\|handle.*error"; then
            suggested_icon="🛡️"
            reason="Protection/validation pattern detected"
        fi
        
        echo -e "    ${CYAN}$feature_id in $file:$line_num${NC}"
        echo -e "      ${GREEN}Suggested action icon: $suggested_icon ($reason)${NC}"
        ((suggestions++))
        
    done < <(grep -rn "// [A-Z]\+-[0-9]\+:" . --include="*.go" 2>/dev/null | head -10)  # Limit to first 10 for demo
    
    echo -e "  ${GREEN}✅ Generated $suggestions action icon suggestions${NC}"
    echo -e "  ${CYAN}💡 Run full migration to apply all suggestions automatically${NC}"
}

# 🔺 DOC-009: Command line interface - 📝 CLI processing
case "${1:-console}" in
    "mapping")
        generate_priority_mapping
        ;;
    "validate")
        validate_token_priorities
        ;;
    "suggest")
        suggest_action_icons
        ;;
    "json")
        echo "{"
        echo "  \"priority_mappings\": ["
        extract_feature_priorities | sed 's/\(.*\):\(.*\)/    { "feature_id": "\1", "priority_icon": "\2" },/' | sed '$ s/,$//'
        echo "  ]"
        echo "}"
        ;;
    "help"|"-h"|"--help")
        echo "🔺 DOC-009: Priority Icon Inference System"
        echo
        echo "Usage:"
        echo "  $0 [console]        - Extract and display priority mappings"
        echo "  $0 mapping          - Generate priority mapping file"
        echo "  $0 validate         - Validate existing token priorities"
        echo "  $0 suggest          - Suggest action icons for tokens"
        echo "  $0 json             - Output mappings in JSON format"
        echo "  $0 help             - Show this help"
        echo
        echo "Examples:"
        echo "  $0                  - Show priority mappings"
        echo "  $0 mapping          - Create priority-icon-mapping.txt"
        echo "  $0 validate         - Check current token priorities"
        ;;
    *)
        extract_feature_priorities | while IFS=':' read -r feature_id priority_icon; do
            echo -e "  ${GREEN}$feature_id${NC} → ${CYAN}$priority_icon${NC}"
        done
        ;;
esac 