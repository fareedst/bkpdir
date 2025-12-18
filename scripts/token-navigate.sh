#!/bin/bash
# [REQ:DOC_016] Token Navigation Tool
# Purpose: Enable AI assistants to navigate codebase using semantic tokens

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
SCRIPT_NAME="token-navigate"

# Usage information
usage() {
    cat << EOF
${BLUE}[REQ:DOC_016] Token Navigation Tool${NC}

Usage: $SCRIPT_NAME <command> [options]

Commands:
  find-req <TOKEN>          Find all implementations of a requirement token
  find-arch <TOKEN>          Find all implementations of an architecture token
  find-impl <TOKEN>          Find all implementations of an implementation token
  find-tests <TOKEN>         Find all tests for a requirement token
  trace <TOKEN>              Trace a token across all layers (REQ → ARCH → IMPL → Code → Tests)
  list-req                   List all requirement tokens in the codebase
  list-arch                  List all architecture tokens in the codebase
  list-impl                  List all implementation tokens in the codebase
  coverage <TOKEN>           Show token coverage across all layers
  help                       Show this help message

Examples:
  $SCRIPT_NAME find-req FILE_BACKUP
  $SCRIPT_NAME find-arch RESOURCE_MANAGEMENT
  $SCRIPT_NAME trace FILE_BACKUP
  $SCRIPT_NAME coverage CONFIGURATION

EOF
}

# Find requirement token implementations
find_req() {
    local token="$1"
    echo -e "${CYAN}Finding implementations of [REQ:${token}]...${NC}\n"
    
    # Find in source code
    echo -e "${BLUE}Source Code:${NC}"
    grep -rn "\[REQ:${token}\]" --include="*.go" --exclude="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No implementations found"
    
    echo ""
    
    # Find in tests
    echo -e "${BLUE}Tests:${NC}"
    grep -rn "REQ_${token}\|\[REQ:${token}\]" --include="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No tests found"
    
    echo ""
    
    # Find in documentation
    echo -e "${BLUE}Documentation:${NC}"
    grep -rn "\[REQ:${token}\]" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | \
        sed 's|^|  |' || echo "  No documentation found"
}

# Find architecture token implementations
find_arch() {
    local token="$1"
    echo -e "${CYAN}Finding implementations of [ARCH:${token}]...${NC}\n"
    
    # Find in source code
    echo -e "${BLUE}Source Code:${NC}"
    grep -rn "\[ARCH:${token}\]" --include="*.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No implementations found"
    
    echo ""
    
    # Find in documentation
    echo -e "${BLUE}Documentation:${NC}"
    grep -rn "\[ARCH:${token}\]" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | \
        sed 's|^|  |' || echo "  No documentation found"
}

# Find implementation token usages
find_impl() {
    local token="$1"
    echo -e "${CYAN}Finding usages of [IMPL:${token}]...${NC}\n"
    
    # Find in source code
    echo -e "${BLUE}Source Code:${NC}"
    grep -rn "\[IMPL:${token}\]" --include="*.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No usages found"
    
    echo ""
    
    # Find in documentation
    echo -e "${BLUE}Documentation:${NC}"
    grep -rn "\[IMPL:${token}\]" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | \
        sed 's|^|  |' || echo "  No documentation found"
}

# Find tests for a requirement
find_tests() {
    local token="$1"
    echo -e "${CYAN}Finding tests for [REQ:${token}]...${NC}\n"
    
    # Find test functions
    echo -e "${BLUE}Test Functions:${NC}"
    grep -rn "func Test.*REQ_${token}\|func Test.*\[REQ:${token}\]" --include="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No test functions found"
    
    echo ""
    
    # Find test comments
    echo -e "${BLUE}Test Comments:${NC}"
    grep -rn "\[REQ:${token}\]" --include="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        sed 's|^|  |' || echo "  No test comments found"
}

# Trace token across all layers
trace_token() {
    local token="$1"
    echo -e "${CYAN}Tracing [REQ:${token}] across all layers...${NC}\n"
    
    # Requirements layer
    echo -e "${GREEN}[1] Requirements Layer:${NC}"
    grep -rn "\[REQ:${token}\]" --include="*.md" "${PROJECT_ROOT}/stdd/requirements.md" 2>/dev/null | \
        head -5 | sed 's|^|  |' || echo "  Not found in requirements.md"
    
    echo ""
    
    # Architecture layer
    echo -e "${GREEN}[2] Architecture Layer:${NC}"
    grep -rn "\[ARCH:.*\].*\[REQ:${token}\]\|\[REQ:${token}\].*\[ARCH:" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | \
        head -5 | sed 's|^|  |' || echo "  Not found in architecture decisions"
    
    echo ""
    
    # Implementation layer
    echo -e "${GREEN}[3] Implementation Layer:${NC}"
    grep -rn "\[IMPL:.*\].*\[REQ:${token}\]\|\[REQ:${token}\].*\[IMPL:" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | \
        head -5 | sed 's|^|  |' || echo "  Not found in implementation decisions"
    
    echo ""
    
    # Source code
    echo -e "${GREEN}[4] Source Code:${NC}"
    grep -rn "\[REQ:${token}\]" --include="*.go" --exclude="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        head -5 | sed 's|^|  |' || echo "  Not found in source code"
    
    echo ""
    
    # Tests
    echo -e "${GREEN}[5] Tests:${NC}"
    grep -rn "REQ_${token}\|\[REQ:${token}\]" --include="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | \
        head -5 | sed 's|^|  |' || echo "  Not found in tests"
}

# List all requirement tokens
list_req() {
    echo -e "${CYAN}All Requirement Tokens:${NC}\n"
    grep -roh "\[REQ:[A-Z0-9_]\+\]" "${PROJECT_ROOT}" --include="*.go" --include="*.md" 2>/dev/null | \
        sort -u | sed 's/\[REQ://;s/\]//' | sed 's/^/  [REQ:/;s/$/]/' | sort
}

# List all architecture tokens
list_arch() {
    echo -e "${CYAN}All Architecture Tokens:${NC}\n"
    grep -roh "\[ARCH:[A-Z0-9_]\+\]" "${PROJECT_ROOT}" --include="*.go" --include="*.md" 2>/dev/null | \
        sort -u | sed 's/\[ARCH://;s/\]//' | sed 's/^/  [ARCH:/;s/$/]/' | sort
}

# List all implementation tokens
list_impl() {
    echo -e "${CYAN}All Implementation Tokens:${NC}\n"
    grep -roh "\[IMPL:[A-Z0-9_]\+\]" "${PROJECT_ROOT}" --include="*.go" --include="*.md" 2>/dev/null | \
        sort -u | sed 's/\[IMPL://;s/\]//' | sed 's/^/  [IMPL:/;s/$/]/' | sort
}

# Show token coverage
coverage() {
    local token="$1"
    echo -e "${CYAN}Token Coverage for [REQ:${token}]:${NC}\n"
    
    local req_count arch_count impl_count code_count test_count
    
    req_count=$(grep -r "\[REQ:${token}\]" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | wc -l | tr -d ' ')
    arch_count=$(grep -r "\[ARCH:.*\].*\[REQ:${token}\]\|\[REQ:${token}\].*\[ARCH:" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | wc -l | tr -d ' ')
    impl_count=$(grep -r "\[IMPL:.*\].*\[REQ:${token}\]\|\[REQ:${token}\].*\[IMPL:" --include="*.md" "${PROJECT_ROOT}/stdd" 2>/dev/null | wc -l | tr -d ' ')
    code_count=$(grep -r "\[REQ:${token}\]" --include="*.go" --exclude="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | wc -l | tr -d ' ')
    test_count=$(grep -r "REQ_${token}\|\[REQ:${token}\]" --include="*_test.go" "${PROJECT_ROOT}" 2>/dev/null | wc -l | tr -d ' ')
    
    echo -e "${BLUE}Requirements:${NC} ${req_count} references"
    echo -e "${BLUE}Architecture:${NC} ${arch_count} references"
    echo -e "${BLUE}Implementation:${NC} ${impl_count} references"
    echo -e "${BLUE}Source Code:${NC} ${code_count} references"
    echo -e "${BLUE}Tests:${NC} ${test_count} references"
    
    echo ""
    
    # Calculate coverage percentage
    local total=$((req_count + arch_count + impl_count + code_count + test_count))
    if [ "$total" -gt 0 ]; then
        local coverage_pct=$(( (req_count + arch_count + impl_count + code_count + test_count) * 100 / 5 ))
        echo -e "${GREEN}Overall Coverage: ${coverage_pct}%${NC}"
    else
        echo -e "${RED}No coverage found${NC}"
    fi
}

# Main command dispatcher
main() {
    if [ $# -eq 0 ]; then
        usage
        exit 0
    fi
    
    case "$1" in
        find-req)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            find_req "$2"
            ;;
        find-arch)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            find_arch "$2"
            ;;
        find-impl)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            find_impl "$2"
            ;;
        find-tests)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            find_tests "$2"
            ;;
        trace)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            trace_token "$2"
            ;;
        list-req)
            list_req
            ;;
        list-arch)
            list_arch
            ;;
        list-impl)
            list_impl
            ;;
        coverage)
            if [ $# -lt 2 ]; then
                echo -e "${RED}Error: Token required${NC}"
                usage
                exit 1
            fi
            coverage "$2"
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            echo -e "${RED}Error: Unknown command '$1'${NC}"
            usage
            exit 1
            ;;
    esac
}

main "$@"
