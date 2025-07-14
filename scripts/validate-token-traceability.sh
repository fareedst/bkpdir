#!/bin/bash
# validate-token-traceability.sh
# Validates comprehensive token consistency across all layers

# ⭐ DOC-016: Token traceability validation
# Ensures every feature has proper token coverage across documentation, source code, and tests

set -e

# [VALIDATION] Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
VALIDATION_REPORT="${PROJECT_ROOT}/docs/validation-reports/token-traceability-report.md"

# [TECHNICAL] Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# [METRICS] Validation counters
TOTAL_FEATURES=0
FEATURES_WITH_IMPL_TOKENS=0
FEATURES_WITH_TEST_TOKENS=0
FEATURES_WITH_ARCH_TOKENS=0
VALIDATION_ERRORS=0
VALIDATION_WARNINGS=0

echo -e "${BLUE}[VALIDATION] Starting comprehensive token traceability validation...${NC}"

# [TECHNICAL] Create validation report directory
mkdir -p "$(dirname "${VALIDATION_REPORT}")"

# [TECHNICAL] Initialize validation report
cat > "${VALIDATION_REPORT}" << EOF
# Token Traceability Validation Report

**Generated**: $(date)
**System**: DOC-016 Comprehensive Token System
**Validation Script**: validate-token-traceability.sh

## [VALIDATION] Summary

EOF

# [SEARCH_DISCOVER] Function to extract feature IDs from feature-tracking.md
extract_feature_ids() {
    echo -e "${BLUE}[SEARCH_DISCOVER] Extracting feature IDs from feature-tracking.md...${NC}"
    
    # Extract feature IDs from the feature registry table
    grep -E "^\| [A-Z]+-[0-9]+" "${PROJECT_ROOT}/docs/context/feature-tracking.md" | \
    sed -E 's/^\| ([A-Z]+-[0-9]+).*/\1/' | \
    sort -u
}

# [SEARCH_DISCOVER] Function to find implementation tokens in source code
find_implementation_tokens() {
    local feature_id="$1"
    echo -e "${BLUE}[SEARCH_DISCOVER] Checking implementation tokens for ${feature_id}...${NC}"
    
    # Find Go files with implementation tokens for this feature
    find "${PROJECT_ROOT}" -name "*.go" -exec grep -l "// [⭐🔺🔶🔻] ${feature_id}:" {} \; 2>/dev/null | sort
}

# [SEARCH_DISCOVER] Function to find test tokens
find_test_tokens() {
    local feature_id="$1"
    echo -e "${BLUE}[SEARCH_DISCOVER] Checking test tokens for ${feature_id}...${NC}"
    
    # Find test files with test tokens for this feature
    find "${PROJECT_ROOT}" -name "*_test.go" -exec grep -l "// [⭐🔺🔶🔻] TEST-${feature_id}:" {} \; 2>/dev/null | sort
}

# [SEARCH_DISCOVER] Function to find architecture decision tokens
find_architecture_tokens() {
    local feature_id="$1"
    echo -e "${BLUE}[SEARCH_DISCOVER] Checking architecture tokens for ${feature_id}...${NC}"
    
    # Find documentation files with architecture decision tokens for this feature
    find "${PROJECT_ROOT}/docs" -name "*.md" -exec grep -l "// [⭐🔺🔶🔻] ARCH-DECISION-.*${feature_id}" {} \; 2>/dev/null | sort
}

# [VALIDATION] Function to validate token format consistency
validate_token_format() {
    local feature_id="$1"
    echo -e "${BLUE}[VALIDATION] Validating token format consistency for ${feature_id}...${NC}"
    
    local impl_tokens arch_tokens test_tokens
    local priority_icons=()
    local format_errors=0
    
    # Check implementation token format
    impl_tokens=$(find "${PROJECT_ROOT}" -name "*.go" -exec grep -H "// [⭐🔺🔶🔻] ${feature_id}:" {} \; 2>/dev/null || true)
    
    # Check test token format
    test_tokens=$(find "${PROJECT_ROOT}" -name "*_test.go" -exec grep -H "// [⭐🔺🔶🔻] TEST-${feature_id}:" {} \; 2>/dev/null || true)
    
    # Check architecture token format
    arch_tokens=$(find "${PROJECT_ROOT}/docs" -name "*.md" -exec grep -H "// [⭐🔺🔶🔻] ARCH-DECISION-.*${feature_id}" {} \; 2>/dev/null || true)
    
    # Validate priority icon consistency
    if [[ -n "$impl_tokens" ]]; then
        local impl_priority=$(echo "$impl_tokens" | head -1 | sed -E 's/.*\/\/ ([⭐🔺🔶🔻]).*/\1/')
        priority_icons+=("$impl_priority")
    fi
    
    if [[ -n "$test_tokens" ]]; then
        local test_priority=$(echo "$test_tokens" | head -1 | sed -E 's/.*\/\/ ([⭐🔺🔶🔻]).*/\1/')
        priority_icons+=("$test_priority")
    fi
    
    # Check if all priority icons are consistent
    if [[ ${#priority_icons[@]} -gt 1 ]]; then
        local first_priority="${priority_icons[0]}"
        for priority in "${priority_icons[@]}"; do
            if [[ "$priority" != "$first_priority" ]]; then
                echo -e "${RED}[ALERT] Priority icon inconsistency for ${feature_id}: ${priority} != ${first_priority}${NC}"
                format_errors=$((format_errors + 1))
            fi
        done
    fi
    
    return $format_errors
}

# [METRICS] Function to generate feature coverage report
generate_coverage_report() {
    local feature_id="$1"
    local impl_files test_files arch_files
    
    impl_files=$(find_implementation_tokens "$feature_id")
    test_files=$(find_test_tokens "$feature_id")
    arch_files=$(find_architecture_tokens "$feature_id")
    
    # Update counters
    TOTAL_FEATURES=$((TOTAL_FEATURES + 1))
    
    if [[ -n "$impl_files" ]]; then
        FEATURES_WITH_IMPL_TOKENS=$((FEATURES_WITH_IMPL_TOKENS + 1))
    fi
    
    if [[ -n "$test_files" ]]; then
        FEATURES_WITH_TEST_TOKENS=$((FEATURES_WITH_TEST_TOKENS + 1))
    fi
    
    if [[ -n "$arch_files" ]]; then
        FEATURES_WITH_ARCH_TOKENS=$((FEATURES_WITH_ARCH_TOKENS + 1))
    fi
    
    # Generate individual feature report
    cat >> "${VALIDATION_REPORT}" << EOF

### Feature: ${feature_id}

#### [IMPLEMENTATION] Implementation Token Coverage
EOF
    
    if [[ -n "$impl_files" ]]; then
        cat >> "${VALIDATION_REPORT}" << EOF
✅ **Implementation tokens found**:
$(echo "$impl_files" | sed 's/^/- /')

EOF
    else
        cat >> "${VALIDATION_REPORT}" << EOF
❌ **No implementation tokens found**

EOF
        echo -e "${RED}[ALERT] Missing implementation tokens for ${feature_id}${NC}"
        VALIDATION_ERRORS=$((VALIDATION_ERRORS + 1))
    fi
    
    cat >> "${VALIDATION_REPORT}" << EOF
#### [TECHNICAL] Test Token Coverage
EOF
    
    if [[ -n "$test_files" ]]; then
        cat >> "${VALIDATION_REPORT}" << EOF
✅ **Test tokens found**:
$(echo "$test_files" | sed 's/^/- /')

EOF
    else
        cat >> "${VALIDATION_REPORT}" << EOF
⚠️ **No test tokens found**

EOF
        echo -e "${YELLOW}[WARNING] Missing test tokens for ${feature_id}${NC}"
        VALIDATION_WARNINGS=$((VALIDATION_WARNINGS + 1))
    fi
    
    cat >> "${VALIDATION_REPORT}" << EOF
#### [BUILD_CONSTRUCT] Architecture Token Coverage
EOF
    
    if [[ -n "$arch_files" ]]; then
        cat >> "${VALIDATION_REPORT}" << EOF
✅ **Architecture tokens found**:
$(echo "$arch_files" | sed 's/^/- /')

EOF
    else
        cat >> "${VALIDATION_REPORT}" << EOF
⚠️ **No architecture decision tokens found**

EOF
        echo -e "${YELLOW}[WARNING] Missing architecture tokens for ${feature_id}${NC}"
        VALIDATION_WARNINGS=$((VALIDATION_WARNINGS + 1))
    fi
    
    # Validate token format consistency
    if ! validate_token_format "$feature_id"; then
        VALIDATION_ERRORS=$((VALIDATION_ERRORS + 1))
    fi
}

# [EXECUTION] Main validation logic
main() {
    echo -e "${BLUE}[EXECUTION] Starting token traceability validation...${NC}"
    
    # Extract all feature IDs
    local feature_ids
    feature_ids=$(extract_feature_ids)
    
    if [[ -z "$feature_ids" ]]; then
        echo -e "${RED}[ALERT] No feature IDs found in feature-tracking.md${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}[SUCCESS] Found $(echo "$feature_ids" | wc -l) features to validate${NC}"
    
    # Validate each feature
    while IFS= read -r feature_id; do
        if [[ -n "$feature_id" ]]; then
            echo -e "${BLUE}[TECHNICAL] Validating feature: ${feature_id}${NC}"
            generate_coverage_report "$feature_id"
        fi
    done <<< "$feature_ids"
    
    # Generate summary report
    local impl_coverage test_coverage arch_coverage
    impl_coverage=$(( FEATURES_WITH_IMPL_TOKENS * 100 / TOTAL_FEATURES ))
    test_coverage=$(( FEATURES_WITH_TEST_TOKENS * 100 / TOTAL_FEATURES ))
    arch_coverage=$(( FEATURES_WITH_ARCH_TOKENS * 100 / TOTAL_FEATURES ))
    
    # Update report summary
    sed -i.bak '/^## \[VALIDATION\] Summary/,/^$/c\
## [VALIDATION] Summary\
\
| [METRICS] Metric | [METRICS] Count | [METRICS] Percentage |\
|------------------|-----------------|----------------------|\
| Total Features | '"$TOTAL_FEATURES"' | 100% |\
| Implementation Token Coverage | '"$FEATURES_WITH_IMPL_TOKENS"' | '"$impl_coverage"'% |\
| Test Token Coverage | '"$FEATURES_WITH_TEST_TOKENS"' | '"$test_coverage"'% |\
| Architecture Token Coverage | '"$FEATURES_WITH_ARCH_TOKENS"' | '"$arch_coverage"'% |\
| Validation Errors | '"$VALIDATION_ERRORS"' | - |\
| Validation Warnings | '"$VALIDATION_WARNINGS"' | - |\
\
### [OBJECTIVE] Coverage Goals\
- **Target**: 100% implementation token coverage\
- **Target**: 90% test token coverage\
- **Target**: 80% architecture token coverage\
\
### [VALIDATION] Validation Status\
' "$VALIDATION_REPORT"
    
    # Clean up backup file
    rm -f "${VALIDATION_REPORT}.bak"
    
    # Display final results
    echo -e "${BLUE}[METRICS] Validation Complete${NC}"
    echo -e "  Total Features: ${TOTAL_FEATURES}"
    echo -e "  Implementation Coverage: ${FEATURES_WITH_IMPL_TOKENS}/${TOTAL_FEATURES} (${impl_coverage}%)"
    echo -e "  Test Coverage: ${FEATURES_WITH_TEST_TOKENS}/${TOTAL_FEATURES} (${test_coverage}%)"
    echo -e "  Architecture Coverage: ${FEATURES_WITH_ARCH_TOKENS}/${TOTAL_FEATURES} (${arch_coverage}%)"
    echo -e "  Errors: ${VALIDATION_ERRORS}"
    echo -e "  Warnings: ${VALIDATION_WARNINGS}"
    
    # Exit with appropriate code
    if [[ $VALIDATION_ERRORS -gt 0 ]]; then
        echo -e "${RED}[ALERT] Token traceability validation FAILED with ${VALIDATION_ERRORS} errors${NC}"
        exit 1
    elif [[ $VALIDATION_WARNINGS -gt 0 ]]; then
        echo -e "${YELLOW}[WARNING] Token traceability validation completed with ${VALIDATION_WARNINGS} warnings${NC}"
        exit 0
    else
        echo -e "${GREEN}[SUCCESS] Token traceability validation PASSED${NC}"
        exit 0
    fi
}

# [EXECUTION] Run main function
main "$@" 