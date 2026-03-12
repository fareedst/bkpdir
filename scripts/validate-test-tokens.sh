#!/bin/bash

# 🔧 TEST-TOKEN-ENFORCEMENT: Test token validation script - 🔧
# Purpose: Validate proper token implementation in test files
# Usage: ./scripts/validate-test-tokens.sh [--strict] [--fix]

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
STRICT_MODE=false
FIX_MODE=false
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --strict)
            STRICT_MODE=true
            shift
            ;;
        --fix)
            FIX_MODE=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            echo "Usage: $0 [--strict] [--fix] [--verbose]"
            echo "  --strict   Enable strict validation mode"
            echo "  --fix      Attempt to fix common token issues"
            echo "  --verbose  Enable verbose output"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

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

# Function to validate token format
validate_token_format() {
    local file="$1"
    local line="$2"
    local line_num="$3"
    
    # Check for proper token format
    if [[ ! "$line" =~ $TOKEN_PATTERN ]]; then
        log_error "Invalid token format in $file:$line_num"
        log_error "  Expected: // [PRIORITY] FEATURE-ID: description [ACTION-category]"
        log_error "  Got: $line"
        ((ERRORS++))
        return 1
    fi
    
    # Extract components
    local priority_level=$(echo "$line" | grep -o '\[CRITICAL\|\[HIGH\|\[MEDIUM\|\[LOW' | head -1 | sed 's/\[//')
    local feature_id=$(echo "$line" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
    
    # Validate priority level
    if [[ ! " ${PRIORITY_LEVELS[@]} " =~ " ${priority_level} " ]]; then
        log_error "Invalid priority level '$priority_level' in $file:$line_num"
        ((ERRORS++))
        return 1
    fi
    
    # Validate feature ID
    if [[ ! " ${EXPECTED_FEATURE_IDS[@]} " =~ " ${feature_id} " ]]; then
        log_warning "Unknown feature ID '$feature_id' in $file:$line_num"
        ((WARNINGS++))
    fi
    
    # Look for action category in the same line or next few lines
    local action_category=""
    local action_found=false
    
    # Check current line for action category
    action_category=$(echo "$line" | grep -o '\[ACTION:(ANALYSIS|DOCUMENTATION|CONFIGURATION|VALIDATION)\]' | head -1 | sed 's/\[ACTION://;s/\]//')
    if [[ -n "$action_category" ]]; then
        action_found=true
    fi
    
    # If not found in current line, check next few lines
    if [[ "$action_found" == false ]]; then
        for ((i=line_num+1; i<=line_num+3; i++)); do
            local next_line=$(sed -n "${i}p" "$file" 2>/dev/null || echo "")
            action_category=$(echo "$next_line" | grep -o '\[ACTION:(ANALYSIS|DOCUMENTATION|CONFIGURATION|VALIDATION)\]' | head -1 | sed 's/\[ACTION://;s/\]//')
            if [[ -n "$action_category" ]]; then
                action_found=true
                break
            fi
        done
    fi
    
    # Validate action category if found
    if [[ "$action_found" == true ]]; then
        if [[ ! " ${ACTION_CATEGORIES[@]} " =~ " ${action_category} " ]]; then
            log_error "Invalid action category '$action_category' in $file:$line_num"
            ((ERRORS++))
            return 1
        fi
    else
        log_warning "Missing action category in $file:$line_num"
        ((WARNINGS++))
    fi
    
    return 0
}

# Function to check for required references
check_required_references() {
    local file="$1"
    local function_start="$2"
    local function_end="$3"
    
    local has_test_ref=false
    local has_immutable_ref=false
    
    # Check lines between function start and end
    for ((i=function_start; i<=function_end; i++)); do
        local line=$(sed -n "${i}p" "$file")
        
        if [[ "$line" =~ $TEST_REF_PATTERN ]]; then
            has_test_ref=true
        fi
        
        if [[ "$line" =~ $IMMUTABLE_REF_PATTERN ]]; then
            has_immutable_ref=true
        fi
    done
    
    if [[ "$has_test_ref" == false ]]; then
        log_warning "Missing TEST-REF comment in $file (function around line $function_start)"
        ((WARNINGS++))
    fi
    
    if [[ "$has_immutable_ref" == false ]]; then
        log_warning "Missing IMMUTABLE-REF comment in $file (function around line $function_start)"
        ((WARNINGS++))
    fi
}

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
    local line_num=1
    local in_test_function=false
    local function_start=0
    
    while IFS= read -r line; do
        # Check if this is a test function start
        if [[ "$line" =~ ^func[[:space:]]+Test ]]; then
            in_test_function=true
            function_start=$line_num
            
            # Check if the function has a token comment
            local has_token=false
            local token_line_num=0
            
            # Look for token comment in the previous few lines
            for ((i=line_num-1; i>=line_num-5; i--)); do
                if [[ $i -gt 0 ]]; then
                    local prev_line=$(sed -n "${i}p" "$file" 2>/dev/null || echo "")
                    if [[ "$prev_line" =~ $TOKEN_PATTERN ]]; then
                        has_token=true
                        token_line_num=$i
                        break
                    fi
                fi
            done
            
            if [[ "$has_token" == true ]]; then
                file_has_tokens=true
                file_test_functions_with_tokens=$((file_test_functions_with_tokens + 1))
                TEST_FUNCTIONS_WITH_TOKENS=$((TEST_FUNCTIONS_WITH_TOKENS + 1))
                
                # Validate token format
                if [[ $token_line_num -gt 0 ]]; then
                    local token_line=$(sed -n "${token_line_num}p" "$file")
                    validate_token_format "$file" "$token_line" "$token_line_num"
                fi
            else
                TEST_FUNCTIONS_WITHOUT_TOKENS=$((TEST_FUNCTIONS_WITHOUT_TOKENS + 1))
                log_error "Test function without token: $(echo "$line" | grep -o 'Test[[:alnum:]_]*') in $file:$line_num"
                ((ERRORS++))
            fi
        fi
        
        # Check for function end (next function or end of file)
        if [[ "$in_test_function" == true ]] && [[ "$line" =~ ^func[[:space:]]+Test ]]; then
            # This is the start of the next test function, so the previous one ended
            if [[ $function_start -gt 0 ]]; then
                check_required_references "$file" $function_start $((line_num - 1))
            fi
            function_start=$line_num
        fi
        
        ((line_num++))
    done < "$file"
    
    # Check the last function if we were in one
    if [[ "$in_test_function" == true ]] && [[ $function_start -gt 0 ]]; then
        check_required_references "$file" $function_start $((line_num - 1))
    fi
    
    # Update statistics
    TOTAL_FILES=$((TOTAL_FILES + 1))
    if [[ "$file_has_tokens" == true ]]; then
        FILES_WITH_TOKENS=$((FILES_WITH_TOKENS + 1))
    else
        FILES_WITHOUT_TOKENS=$((FILES_WITHOUT_TOKENS + 1))
    fi
    
    # Report file statistics
    if [[ "$VERBOSE" == true ]]; then
        log_info "  Test functions: $file_test_functions"
        log_info "  Functions with tokens: $file_test_functions_with_tokens"
        log_info "  Functions without tokens: $((file_test_functions - file_test_functions_with_tokens))"
    fi
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
    
    # Strict mode checks
    if [[ "$STRICT_MODE" == true ]]; then
        echo
        echo "🔒 Strict Mode Checks:"
        
        if [[ $WARNINGS -gt 0 ]]; then
            log_error "❌ Strict mode: Warnings are treated as errors"
            exit 1
        fi
        
        if [[ $TOTAL_FILES -gt 0 ]] && [[ $FILES_WITHOUT_TOKENS -gt 0 ]]; then
            log_error "❌ Strict mode: All files must have tokens"
            exit 1
        fi
        
        if [[ $TOTAL_TEST_FUNCTIONS -gt 0 ]] && [[ $TEST_FUNCTIONS_WITHOUT_TOKENS -gt 0 ]]; then
            log_error "❌ Strict mode: All test functions must have tokens"
            exit 1
        fi
        
        log_success "✅ Strict mode validation passed!"
    fi
}

# Main execution
main() {
    log_info "Starting test token validation..."
    log_info "Mode: $(if [[ "$STRICT_MODE" == true ]]; then echo "STRICT"; else echo "NORMAL"; fi)"
    log_info "Fix mode: $(if [[ "$FIX_MODE" == true ]]; then echo "ENABLED"; else echo "DISABLED"; fi)"
    echo
    
    # Find and validate all test files
    local test_files=$(find_test_files)
    
    if [[ -z "$test_files" ]]; then
        log_error "No test files found!"
        exit 1
    fi
    
    while IFS= read -r file; do
        if [[ -f "$file" ]]; then
            validate_test_file "$file"
        fi
    done <<< "$test_files"
    
    # Generate summary
    generate_summary
}

# Run main function
main "$@" 