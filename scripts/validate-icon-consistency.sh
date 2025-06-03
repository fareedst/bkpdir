#!/bin/bash

# 🔧 DOC-007: Source code icon integration - Validation script
# This script validates implementation token icon consistency between code and documentation

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
errors=0
warnings=0
successes=0

echo -e "${BLUE}🔍 DOC-007: Validating implementation token icon consistency...${NC}"
echo "========================================================="

# Check if we're in the right directory
if [[ ! -f "docs/context/feature-tracking.md" ]]; then
    echo -e "${RED}❌ Error: Must be run from project root directory${NC}"
    exit 1
fi

echo -e "${BLUE}📋 Step 1: Checking for legacy implementation tokens${NC}"
echo "---------------------------------------------------"

legacy_found=false
for go_file in *.go; do
    if [[ -f "$go_file" ]]; then
        echo -e "${BLUE}Checking ${go_file}...${NC}"
        
        # Look for old format tokens (no priority icons)
        while IFS= read -r line; do
            line_num=$(echo "$line" | cut -d: -f1)
            line_content=$(echo "$line" | cut -d: -f2-)
            
            if [[ $line_content =~ //[[:space:]]*([A-Z]+-[0-9]+):[[:space:]]*(.*)$ ]] && 
               [[ ! $line_content =~ (⭐|🔺|🔶|🔻) ]]; then
                echo -e "  ${YELLOW}⚠️ Legacy token format in $go_file:$line_num${NC}"
                echo -e "    ${YELLOW}   Current:$line_content${NC}"
                echo -e "    ${YELLOW}   Suggested: // 🔺 ${BASH_REMATCH[1]}: ${BASH_REMATCH[2]} - 🔧 [Action context]${NC}"
                legacy_found=true
                ((warnings++))
            elif [[ $line_content =~ //[[:space:]]*(⭐|🔺|🔶|🔻)[[:space:]]*([A-Z]+-[0-9]+):[[:space:]]* ]]; then
                echo -e "  ${GREEN}✅ Standardized token found: ${BASH_REMATCH[1]} ${BASH_REMATCH[2]}${NC}"
                ((successes++))
            fi
        done < <(grep -n "// .*[A-Z]\+-[0-9]\+:" "$go_file" 2>/dev/null || true)
    fi
done

if [[ "$legacy_found" = false ]] && [[ $successes -gt 0 ]]; then
    echo -e "  ${GREEN}✅ No legacy implementation token formats found${NC}"
elif [[ $successes -eq 0 ]]; then
    echo -e "  ${YELLOW}⚠️ No implementation tokens found in source code${NC}"
    ((warnings++))
fi

echo
echo -e "${BLUE}📋 Step 2: Checking icon usage in source code${NC}"
echo "---------------------------------------------------"

echo "Icon usage summary in source code:"
for icon in ⭐ 🔺 🔶 🔻 🚀 ⚡ 🔄 🏁 1️⃣ 2️⃣ 3️⃣ ✅ 🔍 📝 🔧 🛡️ 📑 📋 📊 📖; do
    count=0
    for go_file in *.go; do
        if [[ -f "$go_file" ]]; then
            file_count=$(grep -c "$icon" "$go_file" 2>/dev/null || echo "0")
            count=$((count + file_count))
        fi
    done
    
    if [[ $count -gt 0 ]]; then
        case "$icon" in
            "⭐") meaning="CRITICAL" ;;
            "🔺") meaning="HIGH" ;;
            "🔶") meaning="MEDIUM" ;;
            "🔻") meaning="LOW" ;;
            "🔍") meaning="SEARCH" ;;
            "📝") meaning="DOCUMENT" ;;
            "🔧") meaning="CONFIGURE" ;;
            "🛡️") meaning="PROTECT" ;;
            *) meaning="OTHER" ;;
        esac
        echo -e "  ${GREEN}$icon ($meaning): $count instances${NC}"
    fi
done

echo
echo -e "${BLUE}📋 Step 3: Validating DOC-007 implementation${NC}"
echo "---------------------------------------------------"

# Check if DOC-007 guidelines exist
if [[ -f "docs/context/source-code-icon-guidelines.md" ]]; then
    echo -e "  ${GREEN}✅ Source code icon guidelines document exists${NC}"
    ((successes++))
else
    echo -e "  ${RED}❌ Source code icon guidelines document missing${NC}"
    ((errors++))
fi

# Check if this validation script exists
if [[ -f "scripts/validate-icon-consistency.sh" ]]; then
    echo -e "  ${GREEN}✅ Icon validation script exists${NC}"
    ((successes++))
else
    echo -e "  ${RED}❌ Icon validation script missing${NC}"
    ((errors++))
fi

# Check if Makefile integration exists
if grep -q "validate-icons" Makefile 2>/dev/null; then
    echo -e "  ${GREEN}✅ Makefile integration for icon validation exists${NC}"
    ((successes++))
else
    echo -e "  ${RED}❌ Makefile integration for icon validation missing${NC}"
    ((errors++))
fi

echo
echo -e "${BLUE}📋 Step 4: Implementation recommendations${NC}"
echo "---------------------------------------------------"

if [[ $warnings -gt 0 ]]; then
    echo -e "${YELLOW}💡 Recommendations:${NC}"
    echo "1. Update legacy implementation tokens to use standardized format:"
    echo "   // [PRIORITY_ICON] FEATURE-ID: Description - [ACTION_ICON] Context"
    echo "2. Follow examples in docs/context/source-code-icon-guidelines.md"
    echo "3. Use priority icons: ⭐ (CRITICAL), 🔺 (HIGH), 🔶 (MEDIUM), 🔻 (LOW)"
    echo "4. Include action icons: 🔍 (SEARCH), 📝 (DOCUMENT), 🔧 (CONFIGURE), 🛡️ (PROTECT)"
fi

echo
echo "========================================================="
echo -e "${BLUE}📊 Validation Summary${NC}"
echo "========================================================="
echo -e "${GREEN}✅ Successes: $successes${NC}"
echo -e "${YELLOW}⚠️ Warnings: $warnings${NC}"
echo -e "${RED}❌ Errors: $errors${NC}"

if [[ $errors -eq 0 ]]; then
    echo -e "${GREEN}🎉 DOC-007 icon consistency validation passed!${NC}"
    echo -e "${GREEN}   Source code icon integration system is properly implemented.${NC}"
    exit 0
else
    echo -e "${RED}💥 DOC-007 validation failed with $errors errors${NC}"
    echo
    echo -e "${YELLOW}💡 To fix errors:${NC}"
    echo "1. Ensure all DOC-007 deliverables are in place"
    echo "2. Update implementation tokens to standardized format"
    echo "3. Follow guidelines in docs/context/source-code-icon-guidelines.md"
    exit 1
fi 