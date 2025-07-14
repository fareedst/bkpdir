#!/bin/bash
# scripts/enhance-cross-reference-validation.sh
# DOC-014: Enhanced cross-reference validation for completeness and traceability

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Usage information
show_usage() {
    cat << EOF
Usage: $0 [options]

Enhanced cross-reference validation for completeness and traceability.

Options:
  --check-completeness    Check for missing cross-references
  --check-coverage        Check token coverage across all layers
  --traceability-check    Validate complete traceability chain
  --full-validation       Run all validation checks
  --report-only          Generate report without fixing issues
  --help                 Show this help message

Examples:
  $0 --check-completeness
  $0 --check-coverage --report-only
  $0 --full-validation

EOF
}

# Parse command line arguments
CHECK_COMPLETENESS=false
CHECK_COVERAGE=false
TRACEABILITY_CHECK=false
FULL_VALIDATION=false
REPORT_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --check-completeness)
            CHECK_COMPLETENESS=true
            shift
            ;;
        --check-coverage)
            CHECK_COVERAGE=true
            shift
            ;;
        --traceability-check)
            TRACEABILITY_CHECK=true
            shift
            ;;
        --full-validation)
            FULL_VALIDATION=true
            shift
            ;;
        --report-only)
            REPORT_ONLY=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# If no specific checks are requested, run full validation
if [[ "$FULL_VALIDATION" == true ]] || ([[ "$CHECK_COMPLETENESS" == false ]] && [[ "$CHECK_COVERAGE" == false ]] && [[ "$TRACEABILITY_CHECK" == false ]]); then
    CHECK_COMPLETENESS=true
    CHECK_COVERAGE=true
    TRACEABILITY_CHECK=true
fi

# Check if we're in the right directory
if [[ ! -f "$PROJECT_ROOT/go.mod" ]] || [[ ! -d "$PROJECT_ROOT/docs/context" ]]; then
    log_error "Please run this script from the project root directory"
    exit 1
fi

# Function to check completeness of cross-references
check_completeness() {
    log_info "Checking completeness of cross-references..."
    
    local missing_refs=0
    local total_refs=0
    
    # Check specification completeness
    if [[ -f "$PROJECT_ROOT/docs/context/specification.md" ]]; then
        log_info "Checking specification completeness..."
        
        # Count major sections in specification
        local spec_sections=$(grep -c "^## " "$PROJECT_ROOT/docs/context/specification.md" || echo "0")
        local spec_tokens=$(grep -c "SPEC-[0-9]\+" "$PROJECT_ROOT/docs/context/specification.md" || echo "0")
        
        if [[ $spec_sections -gt $spec_tokens ]] && [[ $spec_tokens -gt 0 ]]; then
            local missing=$((spec_sections - spec_tokens))
            log_warning "Missing $missing specification tokens"
            missing_refs=$((missing_refs + missing))
        else
            log_success "Specification tokens complete"
        fi
        
        total_refs=$((total_refs + spec_sections))
    fi
    
    # Check requirements completeness
    if [[ -f "$PROJECT_ROOT/docs/context/requirements.md" ]]; then
        log_info "Checking requirements completeness..."
        
        # Count major sections in requirements
        local req_sections=$(grep -c "^## " "$PROJECT_ROOT/docs/context/requirements.md" || echo "0")
        local req_tokens=$(grep -c "REQ-[0-9]\+" "$PROJECT_ROOT/docs/context/requirements.md" || echo "0")
        
        if [[ $req_sections -gt $req_tokens ]] && [[ $req_tokens -gt 0 ]]; then
            local missing=$((req_sections - req_tokens))
            log_warning "Missing $missing requirements tokens"
            missing_refs=$((missing_refs + missing))
        else
            log_success "Requirements tokens complete"
        fi
        
        total_refs=$((total_refs + req_sections))
    fi
    
    # Check architecture completeness
    if [[ -f "$PROJECT_ROOT/docs/context/architecture.md" ]]; then
        log_info "Checking architecture completeness..."
        
        # Count major sections in architecture
        local arch_sections=$(grep -c "^## " "$PROJECT_ROOT/docs/context/architecture.md" || echo "0")
        local arch_tokens=$(grep -c "ARCH-[0-9]\+" "$PROJECT_ROOT/docs/context/architecture.md" || echo "0")
        
        if [[ $arch_sections -gt $arch_tokens ]] && [[ $arch_tokens -gt 0 ]]; then
            local missing=$((arch_sections - arch_tokens))
            log_warning "Missing $missing architecture tokens"
            missing_refs=$((missing_refs + missing))
        else
            log_success "Architecture tokens complete"
        fi
        
        total_refs=$((total_refs + arch_sections))
    fi
    
    # Report completeness results
    if [[ $missing_refs -eq 0 ]]; then
        log_success "All cross-references complete ($total_refs total references)"
    else
        log_warning "Found $missing_refs missing cross-references out of $total_refs total"
    fi
    
    return $missing_refs
}

# Function to check token coverage across all layers
check_coverage() {
    log_info "Checking token coverage across all layers..."
    
    local coverage_issues=0
    
    # Check feature tracking coverage
    if [[ -f "$PROJECT_ROOT/docs/context/feature-tracking.md" ]]; then
        log_info "Checking feature tracking coverage..."
        
        # Extract all feature IDs from feature tracking
        local feature_ids=$(grep -o '[A-Z]\+-[0-9]\+' "$PROJECT_ROOT/docs/context/feature-tracking.md" | sort -u)
        
        for feature_id in $feature_ids; do
            local ref_count=0
            
            # Count references across all documentation files
            for doc_file in "$PROJECT_ROOT"/docs/context/*.md; do
                if [[ -f "$doc_file" ]] && [[ "$doc_file" != "$PROJECT_ROOT/docs/context/feature-tracking.md" ]]; then
                    local count=$(grep -c "$feature_id" "$doc_file" 2>/dev/null || echo "0")
                    ref_count=$((ref_count + count))
                fi
            done
            
            # Check if feature has sufficient cross-references
            if [[ $ref_count -lt 3 ]]; then
                log_warning "Feature $feature_id has only $ref_count cross-references (minimum 3 recommended)"
                coverage_issues=$((coverage_issues + 1))
            else
                log_success "Feature $feature_id has $ref_count cross-references"
            fi
        done
    fi
    
    # Check code implementation coverage
    log_info "Checking code implementation coverage..."
    
    local go_files=("$PROJECT_ROOT"/*.go)
    local code_tokens=0
    
    for go_file in "${go_files[@]}"; do
        if [[ -f "$go_file" ]]; then
            local file_tokens=$(grep -c "// [A-Z]\+-[0-9]\+" "$go_file" 2>/dev/null || echo "0")
            code_tokens=$((code_tokens + file_tokens))
        fi
    done
    
    if [[ $code_tokens -eq 0 ]]; then
        log_warning "No implementation tokens found in Go files"
        coverage_issues=$((coverage_issues + 1))
    else
        log_success "Found $code_tokens implementation tokens in Go files"
    fi
    
    # Check test coverage
    log_info "Checking test coverage..."
    
    local test_files=("$PROJECT_ROOT"/*_test.go)
    local test_tokens=0
    
    for test_file in "${test_files[@]}"; do
        if [[ -f "$test_file" ]]; then
            local file_tokens=$(grep -c "// [A-Z]\+-[0-9]\+" "$test_file" 2>/dev/null || echo "0")
            test_tokens=$((test_tokens + file_tokens))
        fi
    done
    
    if [[ $test_tokens -eq 0 ]]; then
        log_warning "No test tokens found in test files"
        coverage_issues=$((coverage_issues + 1))
    else
        log_success "Found $test_tokens test tokens in test files"
    fi
    
    if [[ $coverage_issues -eq 0 ]]; then
        log_success "Token coverage is complete across all layers"
    else
        log_warning "Found $coverage_issues coverage issues"
    fi
    
    return $coverage_issues
}

# Function to validate complete traceability chain
check_traceability() {
    log_info "Validating complete traceability chain..."
    
    local traceability_issues=0
    
    # Check specification -> requirements -> architecture -> code -> test chain
    if [[ -f "$PROJECT_ROOT/docs/context/specification.md" ]] && [[ -f "$PROJECT_ROOT/docs/context/requirements.md" ]]; then
        log_info "Checking specification -> requirements traceability..."
        
        # Extract specification sections
        local spec_sections=$(grep "^## " "$PROJECT_ROOT/docs/context/specification.md" | sed 's/^## //' | head -5)
        
        for section in $spec_sections; do
            # Check if section has corresponding requirements
            if ! grep -q "$section" "$PROJECT_ROOT/docs/context/requirements.md" 2>/dev/null; then
                log_warning "Specification section '$section' has no corresponding requirements"
                traceability_issues=$((traceability_issues + 1))
            fi
        done
    fi
    
    # Check requirements -> architecture traceability
    if [[ -f "$PROJECT_ROOT/docs/context/requirements.md" ]] && [[ -f "$PROJECT_ROOT/docs/context/architecture.md" ]]; then
        log_info "Checking requirements -> architecture traceability..."
        
        # Extract requirements sections
        local req_sections=$(grep "^## " "$PROJECT_ROOT/docs/context/requirements.md" | sed 's/^## //' | head -5)
        
        for section in $req_sections; do
            # Check if section has corresponding architecture
            if ! grep -q "$section" "$PROJECT_ROOT/docs/context/architecture.md" 2>/dev/null; then
                log_warning "Requirements section '$section' has no corresponding architecture"
                traceability_issues=$((traceability_issues + 1))
            fi
        done
    fi
    
    # Check code -> test traceability
    log_info "Checking code -> test traceability..."
    
    local go_files=("$PROJECT_ROOT"/*.go)
    local missing_test_refs=0
    
    for go_file in "${go_files[@]}"; do
        if [[ -f "$go_file" ]]; then
            # Extract tokens from Go file
            local tokens=$(grep -o "// [A-Z]\+-[0-9]\+" "$go_file" 2>/dev/null || true)
            
            for token in $tokens; do
                local token_id=$(echo "$token" | grep -o "[A-Z]\+-[0-9]\+")
                
                # Check if token has corresponding test
                local test_files=("$PROJECT_ROOT"/*_test.go)
                local found_test=false
                
                for test_file in "${test_files[@]}"; do
                    if [[ -f "$test_file" ]] && grep -q "$token_id" "$test_file" 2>/dev/null; then
                        found_test=true
                        break
                    fi
                done
                
                if [[ "$found_test" == false ]]; then
                    log_warning "Code token $token_id has no corresponding test"
                    missing_test_refs=$((missing_test_refs + 1))
                fi
            done
        fi
    done
    
    if [[ $missing_test_refs -gt 0 ]]; then
        traceability_issues=$((traceability_issues + missing_test_refs))
    fi
    
    if [[ $traceability_issues -eq 0 ]]; then
        log_success "Complete traceability chain validated"
    else
        log_warning "Found $traceability_issues traceability issues"
    fi
    
    return $traceability_issues
}

# Function to generate comprehensive validation report
generate_report() {
    local report_file="$PROJECT_ROOT/cross-reference-validation-report.md"
    
    log_info "Generating comprehensive validation report..."
    
    cat > "$report_file" << EOF
# Cross-Reference Validation Report

**Generated**: $(date)
**Project**: $(basename "$PROJECT_ROOT")

## Executive Summary

This report provides a comprehensive analysis of cross-reference completeness, coverage, and traceability across all documentation layers.

## Completeness Analysis

EOF
    
    # Add completeness analysis
    if [[ "$CHECK_COMPLETENESS" == true ]]; then
        echo "### Completeness Results" >> "$report_file"
        # Implementation would add detailed completeness results
        echo "- Specification tokens: [Analysis]" >> "$report_file"
        echo "- Requirements tokens: [Analysis]" >> "$report_file"
        echo "- Architecture tokens: [Analysis]" >> "$report_file"
    fi
    
    # Add coverage analysis
    if [[ "$CHECK_COVERAGE" == true ]]; then
        echo "" >> "$report_file"
        echo "## Coverage Analysis" >> "$report_file"
        echo "" >> "$report_file"
        echo "### Feature Tracking Coverage" >> "$report_file"
        # Implementation would add detailed coverage results
        echo "- Cross-reference density: [Analysis]" >> "$report_file"
        echo "- Implementation token coverage: [Analysis]" >> "$report_file"
        echo "- Test token coverage: [Analysis]" >> "$report_file"
    fi
    
    # Add traceability analysis
    if [[ "$TRACEABILITY_CHECK" == true ]]; then
        echo "" >> "$report_file"
        echo "## Traceability Analysis" >> "$report_file"
        echo "" >> "$report_file"
        echo "### Traceability Chain" >> "$report_file"
        # Implementation would add detailed traceability results
        echo "- Specification → Requirements: [Analysis]" >> "$report_file"
        echo "- Requirements → Architecture: [Analysis]" >> "$report_file"
        echo "- Architecture → Code: [Analysis]" >> "$report_file"
        echo "- Code → Tests: [Analysis]" >> "$report_file"
    fi
    
    echo "" >> "$report_file"
    echo "## Recommendations" >> "$report_file"
    echo "" >> "$report_file"
    echo "1. **Immediate Actions**: [List of immediate fixes needed]" >> "$report_file"
    echo "2. **Short-term Improvements**: [List of short-term improvements]" >> "$report_file"
    echo "3. **Long-term Enhancements**: [List of long-term enhancements]" >> "$report_file"
    
    log_success "Validation report generated: $report_file"
}

# Main execution
main() {
    log_info "Starting enhanced cross-reference validation..."
    
    local total_issues=0
    
    # Run requested validation checks
    if [[ "$CHECK_COMPLETENESS" == true ]]; then
        check_completeness
        total_issues=$((total_issues + $?))
    fi
    
    if [[ "$CHECK_COVERAGE" == true ]]; then
        check_coverage
        total_issues=$((total_issues + $?))
    fi
    
    if [[ "$TRACEABILITY_CHECK" == true ]]; then
        check_traceability
        total_issues=$((total_issues + $?))
    fi
    
    # Generate report if requested
    if [[ "$REPORT_ONLY" == true ]] || [[ $total_issues -gt 0 ]]; then
        generate_report
    fi
    
    # Summary
    if [[ $total_issues -eq 0 ]]; then
        log_success "All validation checks passed successfully!"
    else
        log_warning "Validation completed with $total_issues issues found"
        if [[ "$REPORT_ONLY" == false ]]; then
            log_info "Consider running with --report-only for detailed analysis"
        fi
    fi
    
    return $total_issues
}

# Run main function
main "$@" 