#!/bin/bash
# [REQ:DOC_016] Token Coverage Analysis Tool
# Purpose: Analyze token coverage across source code and identify gaps

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_FILE="${PROJECT_ROOT}/docs/validation-reports/token-coverage-analysis.md"

# Statistics
TOTAL_FUNCTIONS=0
FUNCTIONS_WITH_TOKENS=0
FUNCTIONS_WITHOUT_TOKENS=0
TOTAL_TEST_FUNCTIONS=0
TEST_FUNCTIONS_WITH_TOKENS=0
TEST_FUNCTIONS_WITHOUT_TOKENS=0
OVERALL_COVERAGE=0

echo -e "${BLUE}[REQ:DOC_016] Token Coverage Analysis${NC}\n"

# Create report directory
mkdir -p "$(dirname "${REPORT_FILE}")"

# Initialize report
cat > "${REPORT_FILE}" << EOF
# Token Coverage Analysis Report

**Generated**: $(date)
**System**: DOC-016 Comprehensive Token System
**Analysis Script**: token-coverage-analysis.sh

## Summary

EOF

# Analyze source files
analyze_source_files() {
    echo -e "${CYAN}Analyzing source files...${NC}"
    
    local files_without_tokens=()
    
    while IFS= read -r file; do
        if [[ ! -f "$file" ]]; then
            continue
        fi
        
        # Count functions in file (exported functions)
        local func_count
        func_count=$(grep -cE "^func [A-Z]" "$file" 2>/dev/null || true)
        if [[ -z "$func_count" || "$func_count" -eq 0 ]]; then
            continue
        fi
        
        TOTAL_FUNCTIONS=$((TOTAL_FUNCTIONS + func_count))
        
        # Check for tokens
        if grep -qE "\[(REQ|ARCH|IMPL):[A-Z0-9_]+\]" "$file" 2>/dev/null; then
            FUNCTIONS_WITH_TOKENS=$((FUNCTIONS_WITH_TOKENS + func_count))
        else
            FUNCTIONS_WITHOUT_TOKENS=$((FUNCTIONS_WITHOUT_TOKENS + func_count))
            files_without_tokens+=("$file")
        fi
    done < <(find "${PROJECT_ROOT}" -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" -not -path "*/.git/*")
    
    # Report files without tokens
    if [[ ${#files_without_tokens[@]} -gt 0 ]]; then
        echo -e "\n${YELLOW}Files without tokens:${NC}"
        for file in "${files_without_tokens[@]}"; do
            echo "  - ${file#${PROJECT_ROOT}/}"
        done
    fi
}

# Analyze test files
analyze_test_files() {
    echo -e "\n${CYAN}Analyzing test files...${NC}"
    
    local tests_without_tokens=()
    
    while IFS= read -r file; do
        if [[ ! -f "$file" ]]; then
            continue
        fi
        
        # Count test functions
        local test_count
        test_count=$(grep -cE "^func Test" "$file" 2>/dev/null || true)
        if [[ -z "$test_count" || "$test_count" -eq 0 ]]; then
            continue
        fi
        
        TOTAL_TEST_FUNCTIONS=$((TOTAL_TEST_FUNCTIONS + test_count))
        
        # Check for tokens
        if grep -qE "\[REQ:[A-Z0-9_]+\]|REQ_[A-Z0-9_]+" "$file" 2>/dev/null; then
            TEST_FUNCTIONS_WITH_TOKENS=$((TEST_FUNCTIONS_WITH_TOKENS + test_count))
        else
            TEST_FUNCTIONS_WITHOUT_TOKENS=$((TEST_FUNCTIONS_WITHOUT_TOKENS + test_count))
            tests_without_tokens+=("$file")
        fi
    done < <(find "${PROJECT_ROOT}" -name "*_test.go" -not -path "*/vendor/*" -not -path "*/.git/*")
    
    # Report tests without tokens
    if [[ ${#tests_without_tokens[@]} -gt 0 ]]; then
        echo -e "\n${YELLOW}Test files without tokens:${NC}"
        for file in "${tests_without_tokens[@]}"; do
            echo "  - ${file#${PROJECT_ROOT}/}"
        done
    fi
}

# Generate report
generate_report() {
    local source_coverage test_coverage
    
    if [[ $TOTAL_FUNCTIONS -gt 0 ]]; then
        source_coverage=$((FUNCTIONS_WITH_TOKENS * 100 / TOTAL_FUNCTIONS))
    else
        source_coverage=0
    fi
    
    if [[ $TOTAL_TEST_FUNCTIONS -gt 0 ]]; then
        test_coverage=$((TEST_FUNCTIONS_WITH_TOKENS * 100 / TOTAL_TEST_FUNCTIONS))
    else
        test_coverage=0
    fi
    
    if [[ $((TOTAL_FUNCTIONS + TOTAL_TEST_FUNCTIONS)) -gt 0 ]]; then
        OVERALL_COVERAGE=$(((FUNCTIONS_WITH_TOKENS + TEST_FUNCTIONS_WITH_TOKENS) * 100 / (TOTAL_FUNCTIONS + TOTAL_TEST_FUNCTIONS)))
    else
        OVERALL_COVERAGE=0
    fi
    
    # Update report
    cat >> "${REPORT_FILE}" << EOF
### Source Code Coverage
- **Total Functions**: ${TOTAL_FUNCTIONS}
- **Functions with Tokens**: ${FUNCTIONS_WITH_TOKENS}
- **Functions without Tokens**: ${FUNCTIONS_WITHOUT_TOKENS}
- **Coverage**: ${source_coverage}%

### Test Coverage
- **Total Test Functions**: ${TOTAL_TEST_FUNCTIONS}
- **Tests with Tokens**: ${TEST_FUNCTIONS_WITH_TOKENS}
- **Tests without Tokens**: ${TEST_FUNCTIONS_WITHOUT_TOKENS}
- **Coverage**: ${test_coverage}%

### Overall Coverage
- **Overall Coverage**: ${OVERALL_COVERAGE}%
- **Target**: 100%

## Recommendations

1. **Priority**: Add tokens to functions without coverage
2. **Focus**: Start with public functions and test functions
3. **Validation**: Run \`./scripts/token-navigate.sh coverage TOKEN\` to verify
4. **Documentation**: Update \`stdd/semantic-tokens.md\` when adding new tokens

EOF
    
    echo -e "\n${BLUE}Report generated: ${REPORT_FILE}${NC}"
    echo -e "${GREEN}Source Code Coverage: ${source_coverage}%${NC}"
    echo -e "${GREEN}Test Coverage: ${test_coverage}%${NC}"
    echo -e "${GREEN}Overall Coverage: ${OVERALL_COVERAGE}%${NC}"
}

# Main execution
main() {
    analyze_source_files
    analyze_test_files
    generate_report
    
    if [[ ${OVERALL_COVERAGE:-0} -lt 100 ]]; then
        echo -e "\n${YELLOW}Token coverage is below 100%. Review the report for details.${NC}"
        exit 0
    else
        echo -e "\n${GREEN}100% token coverage achieved!${NC}"
        exit 0
    fi
}

main "$@"
