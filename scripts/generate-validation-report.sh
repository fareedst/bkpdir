#!/bin/bash

# Semantic Token Cross-Reference Validation Report Generator
# This script generates comprehensive reports on the semantic token system
# including statistics, cross-references, and validation results.

set -euo pipefail

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPORT_DIR="$PROJECT_ROOT/docs/validation-reports"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Ensure report directory exists
mkdir -p "$REPORT_DIR"

# Token extraction functions
extract_tokens_from_file() {
    local file="$1"
    local layer="$2"
    
    if [[ ! -f "$file" ]]; then
        return 1
    fi
    
    # Extract tokens based on layer
    case "$layer" in
        "docs")
            grep -E "^# [A-Z]+-[A-Z]+-[0-9]+:" "$file" || true
            ;;
        "code")
            grep -E "^// [A-Z]+-[A-Z]+-[0-9]+:" "$file" || true
            ;;
        "test")
            grep -E "^// TEST-[A-Z]+-[A-Z]+-[0-9]+:" "$file" || true
            ;;
        *)
            grep -E "^// [A-Z]+-[A-Z]+-[0-9]+:" "$file" || true
            ;;
    esac
}

# Statistics collection
collect_token_statistics() {
    local report_file="$1"
    
    echo "## Token Statistics" >> "$report_file"
    echo "" >> "$report_file"
    
    # Documentation layer statistics
    echo "### Documentation Layer" >> "$report_file"
    local doc_files=(
        "docs/context/immutable.md"
        "docs/context/requirements.md"
        "docs/context/specification.md"
        "docs/context/architecture.md"
    )
    
    local total_doc_tokens=0
    for doc_file in "${doc_files[@]}"; do
        if [[ -f "$PROJECT_ROOT/$doc_file" ]]; then
            local tokens=$(extract_tokens_from_file "$PROJECT_ROOT/$doc_file" "docs")
            local token_count=$(echo "$tokens" | wc -l)
            total_doc_tokens=$((total_doc_tokens + token_count))
            echo "- $doc_file: $token_count tokens" >> "$report_file"
        fi
    done
    echo "- **Total documentation tokens: $total_doc_tokens**" >> "$report_file"
    echo "" >> "$report_file"
    
    # Code layer statistics
    echo "### Code Layer" >> "$report_file"
    local code_files=(
        "main.go"
        "archive.go"
        "config.go"
        "errors.go"
        "backup.go"
        "formatter.go"
        "git.go"
        "comparison.go"
        "verify.go"
        "exclude.go"
        "file_stats.go"
        "config_adapter.go"
        "config_interfaces.go"
        "formatter_adapter.go"
    )
    
    local total_code_tokens=0
    for code_file in "${code_files[@]}"; do
        if [[ -f "$PROJECT_ROOT/$code_file" ]]; then
            local tokens=$(extract_tokens_from_file "$PROJECT_ROOT/$code_file" "code")
            local token_count=$(echo "$tokens" | wc -l)
            total_code_tokens=$((total_code_tokens + token_count))
            echo "- $code_file: $token_count tokens" >> "$report_file"
        fi
    done
    echo "- **Total code tokens: $total_code_tokens**" >> "$report_file"
    echo "" >> "$report_file"
    
    # Test layer statistics
    echo "### Test Layer" >> "$report_file"
    local test_files=(
        "main_test.go"
        "archive_test.go"
        "config_test.go"
        "backup_test.go"
        "errors_test.go"
        "verify_test.go"
        "git_test.go"
    )
    
    local total_test_tokens=0
    for test_file in "${test_files[@]}"; do
        if [[ -f "$PROJECT_ROOT/$test_file" ]]; then
            local tokens=$(extract_tokens_from_file "$PROJECT_ROOT/$test_file" "test")
            local token_count=$(echo "$tokens" | wc -l)
            total_test_tokens=$((total_test_tokens + token_count))
            echo "- $test_file: $token_count tokens" >> "$report_file"
        fi
    done
    echo "- **Total test tokens: $total_test_tokens**" >> "$report_file"
    echo "" >> "$report_file"
    
    # Package statistics
    echo "### Package Layer" >> "$report_file"
    local package_dirs=(
        "pkg/cli"
        "pkg/git"
        "pkg/config"
        "pkg/errors"
    )
    
    local total_package_tokens=0
    for pkg_dir in "${package_dirs[@]}"; do
        if [[ -d "$PROJECT_ROOT/$pkg_dir" ]]; then
            local test_file="$PROJECT_ROOT/$pkg_dir/cli_test.go"
            if [[ ! -f "$test_file" ]]; then
                test_file="$PROJECT_ROOT/$pkg_dir/git_test.go"
            fi
            if [[ ! -f "$test_file" ]]; then
                test_file="$PROJECT_ROOT/$pkg_dir/config_test.go"
            fi
            if [[ ! -f "$test_file" ]]; then
                test_file="$PROJECT_ROOT/$pkg_dir/errors_test.go"
            fi
            
            if [[ -f "$test_file" ]]; then
                local tokens=$(extract_tokens_from_file "$test_file" "test")
                local token_count=$(echo "$tokens" | wc -l)
                total_package_tokens=$((total_package_tokens + token_count))
                echo "- $pkg_dir: $token_count tokens" >> "$report_file"
            fi
        fi
    done
    echo "- **Total package tokens: $total_package_tokens**" >> "$report_file"
    echo "" >> "$report_file"
    
    # Integration test statistics
    echo "### Integration Tests" >> "$report_file"
    local integration_file="$PROJECT_ROOT/test/integration/doc014_integration_test.go"
    if [[ -f "$integration_file" ]]; then
        local tokens=$(extract_tokens_from_file "$integration_file" "test")
        local token_count=$(echo "$tokens" | wc -l)
        echo "- Integration tests: $token_count tokens" >> "$report_file"
    else
        echo "- Integration tests: 0 tokens (file not found)" >> "$report_file"
    fi
    echo "" >> "$report_file"
    
    # Overall statistics
    local total_tokens=$((total_doc_tokens + total_code_tokens + total_test_tokens + total_package_tokens))
    echo "### Overall Statistics" >> "$report_file"
    echo "- **Total semantic tokens: $total_tokens**" >> "$report_file"
    echo "- Documentation layer: $total_doc_tokens tokens" >> "$report_file"
    echo "- Code layer: $total_code_tokens tokens" >> "$report_file"
    echo "- Test layer: $total_test_tokens tokens" >> "$report_file"
    echo "- Package layer: $total_package_tokens tokens" >> "$report_file"
    echo "" >> "$report_file"
}

# Cross-reference analysis
analyze_cross_references() {
    local report_file="$1"
    
    echo "## Cross-Reference Analysis" >> "$report_file"
    echo "" >> "$report_file"
    
    # Test to Code references
    echo "### Test → Code References" >> "$report_file"
    local test_files=(
        "main_test.go"
        "archive_test.go"
        "config_test.go"
        "backup_test.go"
        "errors_test.go"
        "verify_test.go"
        "git_test.go"
    )
    
    local total_references=0
    local valid_references=0
    local invalid_references=0
    
    for test_file in "${test_files[@]}"; do
        if [[ -f "$PROJECT_ROOT/$test_file" ]]; then
            local tokens=$(extract_tokens_from_file "$PROJECT_ROOT/$test_file" "test")
            while IFS= read -r token; do
                if [[ -n "$token" ]]; then
                    local source_ref=$(echo "$token" | grep -o "Source: [^[:space:]]*" | cut -d' ' -f2)
                    if [[ -n "$source_ref" ]]; then
                        total_references=$((total_references + 1))
                        local source_file=$(echo "$source_ref" | cut -d'/' -f1)
                        local source_token_id=$(echo "$source_ref" | cut -d'/' -f2)
                        
                        if [[ -f "$PROJECT_ROOT/$source_file" ]] && grep -q "$source_token_id" "$PROJECT_ROOT/$source_file"; then
                            valid_references=$((valid_references + 1))
                            echo "- ✅ $test_file → $source_ref" >> "$report_file"
                        else
                            invalid_references=$((invalid_references + 1))
                            echo "- ❌ $test_file → $source_ref (invalid)" >> "$report_file"
                        fi
                    fi
                fi
            done <<< "$tokens"
        fi
    done
    
    echo "" >> "$report_file"
    echo "**Cross-Reference Summary:**" >> "$report_file"
    echo "- Total references: $total_references" >> "$report_file"
    echo "- Valid references: $valid_references" >> "$report_file"
    echo "- Invalid references: $invalid_references" >> "$report_file"
    if [[ $total_references -gt 0 ]]; then
        echo "- Success rate: $((valid_references * 100 / total_references))%" >> "$report_file"
    else
        echo "- Success rate: N/A (no references found)" >> "$report_file"
    fi
    echo "" >> "$report_file"
}

# Token consistency analysis
analyze_token_consistency() {
    local report_file="$1"
    
    echo "## Token Consistency Analysis" >> "$report_file"
    echo "" >> "$report_file"
    
    # Check for duplicate token IDs
    local all_tokens=$(find "$PROJECT_ROOT" -name "*.go" -o -name "*.md" | xargs grep -h -E "^# [A-Z]+-[A-Z]+-[0-9]+:|^// [A-Z]+-[A-Z]+-[0-9]+:|^// TEST-[A-Z]+-[A-Z]+-[0-9]+:" || true)
    local token_ids=$(echo "$all_tokens" | grep -o "[A-Z]+-[A-Z]+-[0-9]+" | sort | uniq -d)
    
    if [[ -n "$token_ids" ]]; then
        echo "### Duplicate Token IDs" >> "$report_file"
        echo "The following token IDs appear multiple times:" >> "$report_file"
        echo "$token_ids" | while read -r token_id; do
            echo "- $token_id" >> "$report_file"
        done
        echo "" >> "$report_file"
    else
        echo "### Duplicate Token IDs" >> "$report_file"
        echo "✅ No duplicate token IDs found" >> "$report_file"
        echo "" >> "$report_file"
    fi
    
    # Token format analysis
    echo "### Token Format Analysis" >> "$report_file"
    local total_tokens=$(echo "$all_tokens" | wc -l)
    local valid_format=0
    local invalid_format=0
    
    while IFS= read -r token; do
        if [[ -n "$token" ]]; then
            if [[ "$token" =~ ^#\ [A-Z]+-[A-Z]+-[0-9]+:\ .* ]] || [[ "$token" =~ ^//\ [A-Z]+-[A-Z]+-[0-9]+:\ .* ]] || [[ "$token" =~ ^//\ TEST-[A-Z]+-[A-Z]+-[0-9]+:\ .* ]]; then
                valid_format=$((valid_format + 1))
            else
                invalid_format=$((invalid_format + 1))
            fi
        fi
    done <<< "$all_tokens"
    
    echo "- Total tokens: $total_tokens" >> "$report_file"
    echo "- Valid format: $valid_format" >> "$report_file"
    echo "- Invalid format: $invalid_format" >> "$report_file"
    if [[ $total_tokens -gt 0 ]]; then
        echo "- Format success rate: $((valid_format * 100 / total_tokens))%" >> "$report_file"
    else
        echo "- Format success rate: N/A (no tokens found)" >> "$report_file"
    fi
    echo "" >> "$report_file"
}

# Generate comprehensive report
generate_report() {
    local report_file="$REPORT_DIR/semantic-token-validation-report-$TIMESTAMP.md"
    
    echo "# Semantic Token Cross-Reference Validation Report" > "$report_file"
    echo "" >> "$report_file"
    echo "**Generated:** $(date)" >> "$report_file"
    echo "**Project:** bkpdir" >> "$report_file"
    echo "**Report ID:** $TIMESTAMP" >> "$report_file"
    echo "" >> "$report_file"
    
    echo "## Executive Summary" >> "$report_file"
    echo "" >> "$report_file"
    echo "This report provides a comprehensive analysis of the semantic token cross-reference system" >> "$report_file"
    echo "implemented across all layers of the bkpdir project: Documentation, Code, Test, and Validation." >> "$report_file"
    echo "" >> "$report_file"
    
    # Collect statistics
    collect_token_statistics "$report_file"
    
    # Analyze cross-references
    analyze_cross_references "$report_file"
    
    # Analyze token consistency
    analyze_token_consistency "$report_file"
    
    # Recommendations
    echo "## Recommendations" >> "$report_file"
    echo "" >> "$report_file"
    echo "### Immediate Actions" >> "$report_file"
    echo "1. Review and fix any invalid cross-references identified above" >> "$report_file"
    echo "2. Ensure all tokens have proper source references and impact assessments" >> "$report_file"
    echo "3. Validate token format consistency across all files" >> "$report_file"
    echo "" >> "$report_file"
    
    echo "### Long-term Improvements" >> "$report_file"
    echo "1. Implement automated validation in CI/CD pipeline" >> "$report_file"
    echo "2. Create token management tools for easier maintenance" >> "$report_file"
    echo "3. Develop visualization tools for cross-reference analysis" >> "$report_file"
    echo "" >> "$report_file"
    
    echo "## Conclusion" >> "$report_file"
    echo "" >> "$report_file"
    echo "The semantic token cross-reference system provides a solid foundation for" >> "$report_file"
    echo "traceability across all layers of the bkpdir project. Regular validation and" >> "$report_file"
    echo "maintenance will ensure continued effectiveness." >> "$report_file"
    echo "" >> "$report_file"
    
    echo "---" >> "$report_file"
    echo "Report generated by semantic token validation system" >> "$report_file"
    
    echo -e "${GREEN}✅ Validation report generated: $report_file${NC}"
    return 0
}

# Main function
main() {
    echo "Semantic Token Cross-Reference Validation Report Generator"
    echo "======================================================="
    echo "Project root: $PROJECT_ROOT"
    echo "Report directory: $REPORT_DIR"
    echo ""
    
    generate_report
    
    echo ""
    echo "Report generation completed successfully!"
}

# Run main function
main "$@" 