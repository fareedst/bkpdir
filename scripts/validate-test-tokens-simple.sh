#!/bin/bash

# 🔧 TEST-TOKEN-ENFORCEMENT: Simplified test token validation script - 🔧
# Purpose: Validate proper token implementation in test files
# Usage: ./scripts/validate-test-tokens-simple.sh

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Statistics
TOTAL_FILES=0
FILES_WITH_TOKENS=0
FILES_WITHOUT_TOKENS=0
TOTAL_TEST_FUNCTIONS=0
TEST_FUNCTIONS_WITH_TOKENS=0
TEST_FUNCTIONS_WITHOUT_TOKENS=0
ERRORS=0
WARNINGS=0

# Token validation patterns
TOKEN_PATTERN='^[[:space:]]*//[[:space:]]*\[(CRITICAL|HIGH|MEDIUM|LOW)\][[:space:]]*[A-Z]+-[0-9]+:'
TEST_REF_PATTERN='^[[:space:]]*//[[:space:]]*TEST-REF:'
IMMUTABLE_REF_PATTERN='^[[:space:]]*//[[:space:]]*IMMUTABLE-REF:'

# Semantic priority levels (AI-first)
PRIORITY_LEVELS=("CRITICAL" "HIGH" "MEDIUM" "LOW")
ACTION_CATEGORIES=("ANALYSIS" "DOCUMENTATION" "CONFIGURATION" "VALIDATION")

# Feature IDs from feature-tracking.md
EXPECTED_FEATURE_IDS=(
    "ARCH-001" "ARCH-002" "ARCH-003" "ARCH-004"
    "FILE-001" "FILE-002" "FILE-003"
    "CLI-015"
    "CFG-001" "CFG-002" "CFG-003" "CFG-004" "CFG-005" "CFG-006"
    "GIT-001" "GIT-002" "GIT-003" "GIT-004" "GIT-005" "GIT-006"
    "OUT-001" "OUT-002"
    "TEST-001" "TEST-002" "TEST-FIX-001" "TEST-INFRA-001-B" "TEST-INFRA-001-E"
    "LINT-001" "COV-001" "COV-002" "COV-003"
    "DOC-001" "DOC-006" "DOC-010" "DOC-012" "DOC-014" "DOC-015" "DOC-016" "DOC-017"
)

# Function to validate a single test file
validate_test_file() {
    local file="$1"
    local file_has_tokens=false
    local file_test_functions=0
    local file_test_functions_with_tokens=0
    
    log_info "Validating: $file"
    
    # Count test functions
    local test_functions=$(grep -c "^func Test" "$file" 2>/dev/null || echo "0")
    file_test_functions=${test_functions:-0}
    TOTAL_TEST_FUNCTIONS=$((TOTAL_TEST_FUNCTIONS + file_test_functions))
    
    if [[ $test_functions -eq 0 ]]; then
        log_warning "No test functions found in $file"
        return
    fi
    
    # Check each test function for tokens
    local test_functions_with_tokens=0
    local test_functions_without_tokens=0
    
    # Get all test function names
    local test_function_names=$(grep "^func Test" "$file" | sed 's/^func //;s/(.*$//')
    
    for func_name in $test_function_names; do
        # Find the line number of this function
        local func_line=$(grep -n "^func $func_name" "$file" | head -1 | cut -d: -f1)
        
        if [[ -n "$func_line" ]]; then
            # Check if there's a token comment in the previous 5 lines
            local has_token=false
            local token_line=""
            
            for ((i=func_line-1; i>=func_line-5; i--)); do
                if [[ $i -gt 0 ]]; then
                    local line=$(sed -n "${i}p" "$file" 2>/dev/null || echo "")
                    if [[ "$line" =~ $TOKEN_PATTERN ]]; then
                        has_token=true
                        token_line="$line"
                        break
                    fi
                fi
            done
            
            if [[ "$has_token" == true ]]; then
                file_has_tokens=true
                file_test_functions_with_tokens=$((file_test_functions_with_tokens + 1))
                TEST_FUNCTIONS_WITH_TOKENS=$((TEST_FUNCTIONS_WITH_TOKENS + 1))
                
                # Check for TEST-REF and IMMUTABLE-REF
                local has_test_ref=false
                local has_immutable_ref=false
                
                for ((i=func_line-5; i<=func_line+5; i++)); do
                    if [[ $i -gt 0 ]]; then
                        local line=$(sed -n "${i}p" "$file" 2>/dev/null || echo "")
                        if [[ "$line" =~ $TEST_REF_PATTERN ]]; then
                            has_test_ref=true
                        fi
                        if [[ "$line" =~ $IMMUTABLE_REF_PATTERN ]]; then
                            has_immutable_ref=true
                        fi
                    fi
                done
                
                if [[ "$has_test_ref" == false ]]; then
                    log_warning "Missing TEST-REF comment for $func_name in $file"
                    ((WARNINGS++))
                fi
                
                if [[ "$has_immutable_ref" == false ]]; then
                    log_warning "Missing IMMUTABLE-REF comment for $func_name in $file"
                    ((WARNINGS++))
                fi
            else
                TEST_FUNCTIONS_WITHOUT_TOKENS=$((TEST_FUNCTIONS_WITHOUT_TOKENS + 1))
                log_error "Test function without token: $func_name in $file:$func_line"
                ((ERRORS++))
            fi
        fi
    done
    
    # Update statistics
    TOTAL_FILES=$((TOTAL_FILES + 1))
    if [[ "$file_has_tokens" == true ]]; then
        FILES_WITH_TOKENS=$((FILES_WITH_TOKENS + 1))
    else
        FILES_WITHOUT_TOKENS=$((FILES_WITHOUT_TOKENS + 1))
    fi
    
    # Report file statistics
    log_info "  Test functions: $file_test_functions"
    log_info "  Functions with tokens: $file_test_functions_with_tokens"
    log_info "  Functions without tokens: $((file_test_functions - file_test_functions_with_tokens))"
}

# Function to find all test files
find_test_files() {
    find . -name "*_test.go" -type f | grep -v "/vendor/" | sort
}

# Function to generate summary report
generate_summary() {
    echo
    echo "=========================================="
    echo "           TOKEN VALIDATION SUMMARY"
    echo "=========================================="
    echo
    
    # File statistics
    echo "📁 File Statistics:"
    echo "  Total test files: $TOTAL_FILES"
    echo "  Files with tokens: $FILES_WITH_TOKENS"
    echo "  Files without tokens: $FILES_WITHOUT_TOKENS"
    
    if [[ $TOTAL_FILES -gt 0 ]]; then
        local file_coverage=$((FILES_WITH_TOKENS * 100 / TOTAL_FILES))
        echo "  File token coverage: ${file_coverage}%"
    fi
    
    echo
    
    # Function statistics
    echo "🧪 Test Function Statistics:"
    echo "  Total test functions: $TOTAL_TEST_FUNCTIONS"
    echo "  Functions with tokens: $TEST_FUNCTIONS_WITH_TOKENS"
    echo "  Functions without tokens: $TEST_FUNCTIONS_WITHOUT_TOKENS"
    
    if [[ $TOTAL_TEST_FUNCTIONS -gt 0 ]]; then
        local function_coverage=$((TEST_FUNCTIONS_WITH_TOKENS * 100 / TOTAL_TEST_FUNCTIONS))
        echo "  Function token coverage: ${function_coverage}%"
    fi
    
    echo
    
    # Error and warning statistics
    echo "⚠️  Issues Found:"
    echo "  Errors: $ERRORS"
    echo "  Warnings: $WARNINGS"
    echo
    
    # Overall status
    if [[ $ERRORS -eq 0 ]] && [[ $WARNINGS -eq 0 ]]; then
        log_success "✅ All test files have proper token implementation!"
    elif [[ $ERRORS -eq 0 ]]; then
        log_warning "⚠️  Token validation completed with warnings"
    else
        log_error "❌ Token validation failed with errors"
        exit 1
    fi
}

# Main execution
main() {
    log_info "Starting test token validation..."
    
    # Find and validate all test files
    local test_files=$(find_test_files)
    
    for file in $test_files; do
        validate_test_file "$file"
    done
    
    # Generate summary
    generate_summary
}

# Run main function
main "$@" 