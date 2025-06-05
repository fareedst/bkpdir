#!/usr/bin/env bash
# Decision Context Validation Script
# 🔶 DOC-014: Decision validation tools - Enhanced token decision context validation
# Validates decision context format in enhanced implementation tokens

set -euo pipefail

# Check for bash version 4+ (required for associative arrays)
if [[ ${BASH_VERSION%%.*} -lt 4 ]]; then
    echo "Error: This script requires bash 4.0 or later for associative arrays"
    echo "Current bash version: $BASH_VERSION"
    exit 1
fi

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DOCS_DIR="$PROJECT_ROOT/docs"
VALIDATION_REPORTS_DIR="$DOCS_DIR/validation-reports"
CONTEXT_SPECS_FILE="$DOCS_DIR/context/enhanced-token-syntax-specification.md"

# Validation settings
VALIDATION_MODE="standard"    # standard, strict, permissive
OUTPUT_FORMAT="detailed"      # detailed, summary, json
CHECK_FORMAT=true
CHECK_CATEGORIES=true
CHECK_USAGE=true
GENERATE_SUGGESTIONS=true
VERBOSE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Validation results
declare -A VALIDATION_STATS
declare -A FORMAT_ERRORS
declare -A CATEGORY_VIOLATIONS
declare -A USAGE_ISSUES

# Decision context categories (from enhanced token specification)
declare -a IMPACT_CATEGORIES=(
    "core-functionality"
    "enhancement" 
    "developer-experience"
    "infrastructure"
)

declare -a DEPENDENCY_CATEGORIES=(
    "blocks-extraction"
    "enables-flexibility" 
    "prerequisite-for-testing"
    "independent"
)

declare -a CONSTRAINT_CATEGORIES=(
    "user-facing"
    "backward-compatible"
    "breaking-change"
    "performance-critical"
)

# Usage information
usage() {
    echo "Decision Context Validation Script (DOC-014)"
    echo ""
    echo "Usage: $0 [OPTIONS] [FILES...]"
    echo ""
    echo "Options:"
    echo "  --mode <standard|strict|permissive>  Validation mode (default: standard)"
    echo "  --format <detailed|summary|json>     Output format (default: detailed)"
    echo "  --no-format                          Skip format validation"
    echo "  --no-categories                      Skip category validation"
    echo "  --no-usage                           Skip usage pattern validation"
    echo "  --no-suggestions                     Skip generating improvement suggestions"
    echo "  --verbose                            Enable verbose output"
    echo "  --help                               Show this help message"
    echo ""
    echo "Validation Modes:"
    echo "  standard     Standard validation with reasonable tolerance"
    echo "  strict       Strict validation with zero tolerance for deviations"
    echo "  permissive   Permissive validation focusing on major issues only"
    echo ""
    echo "Examples:"
    echo "  $0                                   Validate all Go files"
    echo "  $0 --mode strict main.go            Strict validation of specific file"
    echo "  $0 --format json --no-suggestions   JSON output without suggestions"
}

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*"
}

log_verbose() {
    if [[ "$VERBOSE" == "true" ]]; then
        echo -e "${CYAN}[VERBOSE]${NC} $*"
    fi
}

log_validation() {
    echo -e "${PURPLE}[VALIDATION]${NC} $*"
}

# Initialize validation
init_validation() {
    VALIDATION_STATS[files_processed]=0
    VALIDATION_STATS[tokens_found]=0
    VALIDATION_STATS[tokens_with_context]=0
    VALIDATION_STATS[format_errors]=0
    VALIDATION_STATS[category_errors]=0
    VALIDATION_STATS[usage_warnings]=0
    VALIDATION_STATS[suggestions_generated]=0
}

# Validate decision context format
validate_decision_context_format() {
    local file="$1"
    local line_num="$2"
    local line_content="$3"
    local context_string="$4"
    
    log_verbose "Validating decision context format: $context_string"
    
    # Check basic format: [DECISION: context1, context2, context3]
    if [[ ! "$context_string" =~ ^\[DECISION:[[:space:]]*([^]]+)\]$ ]]; then
        FORMAT_ERRORS["$file:$line_num"]="Invalid decision context format: $context_string"
        ((VALIDATION_STATS[format_errors]++))
        return 1
    fi
    
    local context_content="${BASH_REMATCH[1]}"
    
    # Parse individual context tags
    local IFS=','
    local tags=($context_content)
    local tag_count=${#tags[@]}
    
    # Check tag count (max 3 recommended)
    if [[ $tag_count -gt 3 ]]; then
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            FORMAT_ERRORS["$file:$line_num"]="Too many decision context tags ($tag_count > 3): $context_string"
            ((VALIDATION_STATS[format_errors]++))
            return 1
        else
            log_warning "Many decision context tags ($tag_count) in $file:$line_num"
        fi
    fi
    
    # Check minimum tags (at least 1)
    if [[ $tag_count -eq 0 ]]; then
        FORMAT_ERRORS["$file:$line_num"]="Empty decision context: $context_string"
        ((VALIDATION_STATS[format_errors]++))
        return 1
    fi
    
    # Validate individual tags
    local invalid_tags=()
    for tag in "${tags[@]}"; do
        # Clean up tag (remove whitespace)
        tag=$(echo "$tag" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        if ! is_valid_context_category "$tag"; then
            invalid_tags+=("$tag")
        fi
    done
    
    if [[ ${#invalid_tags[@]} -gt 0 ]]; then
        CATEGORY_VIOLATIONS["$file:$line_num"]="Invalid decision context categories: ${invalid_tags[*]}"
        ((VALIDATION_STATS[category_errors]++))
        return 1
    fi
    
    return 0
}

# Check if context category is valid
is_valid_context_category() {
    local category="$1"
    
    # Check impact categories
    for valid_category in "${IMPACT_CATEGORIES[@]}"; do
        if [[ "$category" == "$valid_category" ]]; then
            return 0
        fi
    done
    
    # Check dependency categories
    for valid_category in "${DEPENDENCY_CATEGORIES[@]}"; do
        if [[ "$category" == "$valid_category" ]]; then
            return 0
        fi
    done
    
    # Check constraint categories
    for valid_category in "${CONSTRAINT_CATEGORIES[@]}"; do
        if [[ "$category" == "$valid_category" ]]; then
            return 0
        fi
    done
    
    return 1
}

# Validate decision context usage patterns
validate_decision_context_usage() {
    local file="$1"
    local line_num="$2"
    local full_token="$3"
    local context_string="$4"
    
    log_verbose "Validating decision context usage: $full_token"
    
    # Extract feature ID from token
    local feature_id=""
    if [[ "$full_token" =~ [⭐🔺🔶🔻][[:space:]]+([A-Z]+-[0-9]+): ]]; then
        feature_id="${BASH_REMATCH[1]}"
    fi
    
    # Check if context matches token purpose
    local token_purpose=""
    if [[ "$full_token" =~ [⭐🔺🔶🔻][[:space:]]+[A-Z]+-[0-9]+:[[:space:]]*([^[]+) ]]; then
        token_purpose="${BASH_REMATCH[1]}"
    fi
    
    # Validate context consistency
    validate_context_consistency "$file" "$line_num" "$feature_id" "$token_purpose" "$context_string"
    
    # Check for common usage patterns
    check_common_usage_patterns "$file" "$line_num" "$full_token" "$context_string"
}

# Validate context consistency with token purpose
validate_context_consistency() {
    local file="$1"
    local line_num="$2"
    local feature_id="$3"
    local token_purpose="$4"
    local context_string="$5"
    
    # Extract context tags
    local context_content
    if [[ "$context_string" =~ ^\[DECISION:[[:space:]]*([^]]+)\]$ ]]; then
        context_content="${BASH_REMATCH[1]}"
    else
        return 1
    fi
    
    local IFS=','
    local tags=($context_content)
    
    # Check for logical consistency
    local has_core_functionality=false
    local has_breaking_change=false
    local has_backward_compatible=false
    
    for tag in "${tags[@]}"; do
        tag=$(echo "$tag" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        case "$tag" in
            "core-functionality") has_core_functionality=true ;;
            "breaking-change") has_breaking_change=true ;;
            "backward-compatible") has_backward_compatible=true ;;
        esac
    done
    
    # Check for contradictory tags
    if [[ "$has_breaking_change" == "true" && "$has_backward_compatible" == "true" ]]; then
        USAGE_ISSUES["$file:$line_num"]="Contradictory context tags: breaking-change and backward-compatible"
        ((VALIDATION_STATS[usage_warnings]++))
    fi
    
    # Check critical feature consistency
    if [[ "$feature_id" =~ ^(ARCH|CFG|GIT)-[0-9]+$ && "$has_core_functionality" == "false" ]]; then
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            USAGE_ISSUES["$file:$line_num"]="Core feature $feature_id should include core-functionality context"
            ((VALIDATION_STATS[usage_warnings]++))
        fi
    fi
}

# Check common usage patterns
check_common_usage_patterns() {
    local file="$1"
    local line_num="$2"
    local full_token="$3"
    local context_string="$4"
    
    # Pattern 1: Test-related tokens should have appropriate context
    if [[ "$full_token" =~ TEST.*[Tt]est ]]; then
        if [[ ! "$context_string" =~ prerequisite-for-testing ]]; then
            log_verbose "Test-related token might benefit from 'prerequisite-for-testing' context"
        fi
    fi
    
    # Pattern 2: Refactoring tokens should indicate dependencies
    if [[ "$full_token" =~ REFACTOR ]]; then
        if [[ ! "$context_string" =~ (blocks-extraction|enables-flexibility|independent) ]]; then
            log_verbose "Refactoring token should specify dependency impact"
        fi
    fi
    
    # Pattern 3: Configuration changes should specify compatibility
    if [[ "$full_token" =~ CFG.*[Cc]onfig ]]; then
        if [[ ! "$context_string" =~ (backward-compatible|breaking-change) ]]; then
            log_verbose "Configuration change should specify compatibility impact"
        fi
    fi
}

# Generate improvement suggestions
generate_improvement_suggestions() {
    local file="$1"
    local line_num="$2"
    local full_token="$3"
    local context_string="$4"
    
    if [[ "$GENERATE_SUGGESTIONS" != "true" ]]; then
        return 0
    fi
    
    local suggestions=()
    
    # Suggest missing impact category
    if [[ ! "$context_string" =~ (core-functionality|enhancement|developer-experience|infrastructure) ]]; then
        suggestions+=("Consider adding impact category: core-functionality, enhancement, developer-experience, or infrastructure")
    fi
    
    # Suggest dependency information for complex features
    if [[ "$full_token" =~ (EXTRACT|REFACTOR) && ! "$context_string" =~ (blocks-extraction|enables-flexibility|prerequisite-for-testing|independent) ]]; then
        suggestions+=("Consider adding dependency context: blocks-extraction, enables-flexibility, prerequisite-for-testing, or independent")
    fi
    
    # Suggest compatibility information for user-facing changes
    if [[ "$full_token" =~ (API|CLI|CONFIG) && ! "$context_string" =~ (user-facing|backward-compatible|breaking-change) ]]; then
        suggestions+=("Consider adding compatibility context: user-facing, backward-compatible, or breaking-change")
    fi
    
    # Store suggestions
    if [[ ${#suggestions[@]} -gt 0 ]]; then
        local suggestion_text=$(IFS=$'\n'; echo "${suggestions[*]}")
        USAGE_ISSUES["$file:$line_num:suggestions"]="$suggestion_text"
        ((VALIDATION_STATS[suggestions_generated]++))
    fi
}

# Validate a single file
validate_file() {
    local file="$1"
    
    if [[ ! -f "$file" ]]; then
        log_error "File not found: $file"
        return 1
    fi
    
    if [[ ! "$file" =~ \.go$ ]]; then
        log_verbose "Skipping non-Go file: $file"
        return 0
    fi
    
    ((VALIDATION_STATS[files_processed]++))
    log_verbose "Processing file: $file"
    
    local line_num=0
    while IFS= read -r line; do
        ((line_num++))
        
        # Look for implementation tokens
        if [[ "$line" =~ ^[[:space:]]*//[[:space:]]*[⭐🔺🔶🔻] ]]; then
            ((VALIDATION_STATS[tokens_found]++))
            log_verbose "Found token at line $line_num: $line"
            
            # Check for decision context
            if [[ "$line" =~ \[DECISION:[^]]*\] ]]; then
                ((VALIDATION_STATS[tokens_with_context]++))
                
                # Extract decision context
                local context_match
                if [[ "$line" =~ (\[DECISION:[^]]*\]) ]]; then
                    context_match="${BASH_REMATCH[1]}"
                    
                    # Validate format
                    if [[ "$CHECK_FORMAT" == "true" ]]; then
                        validate_decision_context_format "$file" "$line_num" "$line" "$context_match"
                    fi
                    
                    # Validate categories
                    if [[ "$CHECK_CATEGORIES" == "true" ]]; then
                        # Categories are validated within format validation
                        :
                    fi
                    
                    # Validate usage patterns
                    if [[ "$CHECK_USAGE" == "true" ]]; then
                        validate_decision_context_usage "$file" "$line_num" "$line" "$context_match"
                    fi
                    
                    # Generate suggestions
                    generate_improvement_suggestions "$file" "$line_num" "$line" "$context_match"
                fi
            fi
        fi
    done < "$file"
    
    return 0
}

# Generate validation report
generate_validation_report() {
    log_info "📄 Generating decision context validation report..."
    
    mkdir -p "$VALIDATION_REPORTS_DIR"
    local report_file="$VALIDATION_REPORTS_DIR/decision-context-validation-report.md"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    cat > "$report_file" << EOF
# Decision Context Validation Report (DOC-014)

> **Generated on:** \`$timestamp\`  
> **Validation Mode:** \`$VALIDATION_MODE\`  
> **Report Format:** \`$OUTPUT_FORMAT\`

## 🎯 Executive Summary

Decision context validation has been completed for enhanced implementation tokens.

### 📊 Validation Statistics
- **Files Processed:** ${VALIDATION_STATS[files_processed]}
- **Implementation Tokens Found:** ${VALIDATION_STATS[tokens_found]}
- **Tokens with Decision Context:** ${VALIDATION_STATS[tokens_with_context]}
- **Enhancement Rate:** $(( VALIDATION_STATS[tokens_found] > 0 ? VALIDATION_STATS[tokens_with_context] * 100 / VALIDATION_STATS[tokens_found] : 0 ))%

### 🚨 Issues Found
- **Format Errors:** ${VALIDATION_STATS[format_errors]}
- **Category Violations:** ${VALIDATION_STATS[category_errors]}
- **Usage Warnings:** ${VALIDATION_STATS[usage_warnings]}
- **Suggestions Generated:** ${VALIDATION_STATS[suggestions_generated]}

## 📋 Detailed Validation Results

### 🔍 Format Validation
EOF

    if [[ ${VALIDATION_STATS[format_errors]} -eq 0 ]]; then
        echo "**Status:** ✅ PASSED - No format errors found" >> "$report_file"
    else
        echo "**Status:** ❌ FAILED - ${VALIDATION_STATS[format_errors]} format errors found" >> "$report_file"
        echo "" >> "$report_file"
        echo "**Format Errors:**" >> "$report_file"
        for error_key in "${!FORMAT_ERRORS[@]}"; do
            echo "- \`$error_key\`: ${FORMAT_ERRORS[$error_key]}" >> "$report_file"
        done
    fi

    cat >> "$report_file" << EOF

### 🏷️ Category Validation
EOF

    if [[ ${VALIDATION_STATS[category_errors]} -eq 0 ]]; then
        echo "**Status:** ✅ PASSED - All decision context categories are valid" >> "$report_file"
    else
        echo "**Status:** ❌ FAILED - ${VALIDATION_STATS[category_errors]} category violations found" >> "$report_file"
        echo "" >> "$report_file"
        echo "**Category Violations:**" >> "$report_file"
        for error_key in "${!CATEGORY_VIOLATIONS[@]}"; do
            echo "- \`$error_key\`: ${CATEGORY_VIOLATIONS[$error_key]}" >> "$report_file"
        done
    fi

    cat >> "$report_file" << EOF

### 🎯 Usage Pattern Analysis
EOF

    if [[ ${VALIDATION_STATS[usage_warnings]} -eq 0 ]]; then
        echo "**Status:** ✅ GOOD - No usage pattern issues found" >> "$report_file"
    else
        echo "**Status:** ⚠️ WARNINGS - ${VALIDATION_STATS[usage_warnings]} usage issues found" >> "$report_file"
        echo "" >> "$report_file"
        echo "**Usage Issues:**" >> "$report_file"
        for issue_key in "${!USAGE_ISSUES[@]}"; do
            if [[ ! "$issue_key" =~ :suggestions$ ]]; then
                echo "- \`$issue_key\`: ${USAGE_ISSUES[$issue_key]}" >> "$report_file"
            fi
        done
    fi

    if [[ ${VALIDATION_STATS[suggestions_generated]} -gt 0 ]]; then
        cat >> "$report_file" << EOF

### 💡 Improvement Suggestions

EOF
        for suggestion_key in "${!USAGE_ISSUES[@]}"; do
            if [[ "$suggestion_key" =~ :suggestions$ ]]; then
                local file_line="${suggestion_key%:suggestions}"
                echo "**$file_line:**" >> "$report_file"
                echo "${USAGE_ISSUES[$suggestion_key]}" | sed 's/^/- /' >> "$report_file"
                echo "" >> "$report_file"
            fi
        done
    fi

    cat >> "$report_file" << EOF

## 📚 Valid Decision Context Categories

### Impact Level
$(printf "- %s\n" "${IMPACT_CATEGORIES[@]}")

### Dependencies
$(printf "- %s\n" "${DEPENDENCY_CATEGORIES[@]}")

### Constraints
$(printf "- %s\n" "${CONSTRAINT_CATEGORIES[@]}")

## 🔧 Usage Examples

### Standard Format
\`\`\`go
// ⭐ ARCH-001: Archive naming convention [DECISION: core-functionality, user-facing, backward-compatible]
// 🔺 CFG-005: Configuration inheritance [DECISION: enhancement, enables-flexibility]
// 🔶 DOC-010: Token format suggestions [DECISION: developer-experience, independent]
\`\`\`

### Context Selection Guidelines
1. **Always include one Impact Level category**
2. **Include Dependencies for complex features**
3. **Specify Constraints for user-facing changes**
4. **Maximum 3 categories per decision context**

---

*This report was generated automatically by the DOC-014 Decision Context validation system.*
EOF

    log_success "Validation report generated: $report_file"
}

# Output results
output_results() {
    case "$OUTPUT_FORMAT" in
        "json")
            output_json_results
            ;;
        "summary")
            output_summary_results
            ;;
        "detailed"|*)
            output_detailed_results
            ;;
    esac
}

# Output detailed results
output_detailed_results() {
    echo ""
    echo "🎯 Decision Context Validation Results"
    echo "====================================="
    echo ""
    
    echo "📊 Statistics:"
    echo "  Files Processed:           ${VALIDATION_STATS[files_processed]}"
    echo "  Implementation Tokens:     ${VALIDATION_STATS[tokens_found]}"
    echo "  Tokens with Context:       ${VALIDATION_STATS[tokens_with_context]}"
    echo "  Enhancement Rate:          $(( VALIDATION_STATS[tokens_found] > 0 ? VALIDATION_STATS[tokens_with_context] * 100 / VALIDATION_STATS[tokens_found] : 0 ))%"
    echo ""
    
    echo "🚨 Issues:"
    echo "  Format Errors:             ${VALIDATION_STATS[format_errors]}"
    echo "  Category Violations:       ${VALIDATION_STATS[category_errors]}"
    echo "  Usage Warnings:            ${VALIDATION_STATS[usage_warnings]}"
    echo "  Suggestions Generated:     ${VALIDATION_STATS[suggestions_generated]}"
    echo ""
    
    if [[ ${VALIDATION_STATS[format_errors]} -gt 0 ]]; then
        echo "❌ Format Errors:"
        for error_key in "${!FORMAT_ERRORS[@]}"; do
            echo "  $error_key: ${FORMAT_ERRORS[$error_key]}"
        done
        echo ""
    fi
    
    if [[ ${VALIDATION_STATS[category_errors]} -gt 0 ]]; then
        echo "❌ Category Violations:"
        for error_key in "${!CATEGORY_VIOLATIONS[@]}"; do
            echo "  $error_key: ${CATEGORY_VIOLATIONS[$error_key]}"
        done
        echo ""
    fi
    
    if [[ ${VALIDATION_STATS[usage_warnings]} -gt 0 ]]; then
        echo "⚠️ Usage Issues:"
        for issue_key in "${!USAGE_ISSUES[@]}"; do
            if [[ ! "$issue_key" =~ :suggestions$ ]]; then
                echo "  $issue_key: ${USAGE_ISSUES[$issue_key]}"
            fi
        done
        echo ""
    fi
}

# Output summary results
output_summary_results() {
    local enhancement_rate=0
    if [[ ${VALIDATION_STATS[tokens_found]} -gt 0 ]]; then
        enhancement_rate=$(( VALIDATION_STATS[tokens_with_context] * 100 / VALIDATION_STATS[tokens_found] ))
    fi
    
    local overall_status="PASS"
    if [[ ${VALIDATION_STATS[format_errors]} -gt 0 || ${VALIDATION_STATS[category_errors]} -gt 0 ]]; then
        overall_status="FAIL"
    elif [[ ${VALIDATION_STATS[usage_warnings]} -gt 0 ]]; then
        overall_status="WARN"
    fi
    
    echo "Decision Context Validation: $overall_status"
    echo "Enhancement Rate: $enhancement_rate%"
    echo "Errors: $((VALIDATION_STATS[format_errors] + VALIDATION_STATS[category_errors]))"
    echo "Warnings: ${VALIDATION_STATS[usage_warnings]}"
}

# Output JSON results
output_json_results() {
    local enhancement_rate=0
    if [[ ${VALIDATION_STATS[tokens_found]} -gt 0 ]]; then
        enhancement_rate=$(( VALIDATION_STATS[tokens_with_context] * 100 / VALIDATION_STATS[tokens_found] ))
    fi
    
    cat << EOF
{
  "validation_timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "validation_mode": "$VALIDATION_MODE",
  "statistics": {
    "files_processed": ${VALIDATION_STATS[files_processed]},
    "tokens_found": ${VALIDATION_STATS[tokens_found]},
    "tokens_with_context": ${VALIDATION_STATS[tokens_with_context]},
    "enhancement_rate": $enhancement_rate,
    "format_errors": ${VALIDATION_STATS[format_errors]},
    "category_errors": ${VALIDATION_STATS[category_errors]},
    "usage_warnings": ${VALIDATION_STATS[usage_warnings]},
    "suggestions_generated": ${VALIDATION_STATS[suggestions_generated]}
  },
  "overall_status": "$([ $((VALIDATION_STATS[format_errors] + VALIDATION_STATS[category_errors])) -eq 0 ] && echo "PASS" || echo "FAIL")"
}
EOF
}

# Parse command line arguments
FILES=()
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode)
            VALIDATION_MODE="$2"
            if [[ ! "$VALIDATION_MODE" =~ ^(standard|strict|permissive)$ ]]; then
                log_error "Invalid validation mode: $VALIDATION_MODE"
                usage
                exit 1
            fi
            shift 2
            ;;
        --format)
            OUTPUT_FORMAT="$2"
            if [[ ! "$OUTPUT_FORMAT" =~ ^(detailed|summary|json)$ ]]; then
                log_error "Invalid output format: $OUTPUT_FORMAT"
                usage
                exit 1
            fi
            shift 2
            ;;
        --no-format)
            CHECK_FORMAT=false
            shift
            ;;
        --no-categories)
            CHECK_CATEGORIES=false
            shift
            ;;
        --no-usage)
            CHECK_USAGE=false
            shift
            ;;
        --no-suggestions)
            GENERATE_SUGGESTIONS=false
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            usage
            exit 0
            ;;
        -*)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
        *)
            FILES+=("$1")
            shift
            ;;
    esac
done

# Main execution
main() {
    log_info "🔶 DOC-014: Decision Context Validation"
    log_info "Mode: $VALIDATION_MODE | Format: $OUTPUT_FORMAT"
    echo ""
    
    # Initialize validation
    init_validation
    
    # Determine files to process
    if [[ ${#FILES[@]} -eq 0 ]]; then
        # Process all Go files
        while IFS= read -r -d '' file; do
            FILES+=("$file")
        done < <(find "$PROJECT_ROOT" -name "*.go" -type f -not -path "*/vendor/*" -print0)
        log_verbose "Processing ${#FILES[@]} Go files"
    else
        log_verbose "Processing ${#FILES[@]} specified files"
    fi
    
    # Validate files
    local validation_failed=false
    for file in "${FILES[@]}"; do
        if ! validate_file "$file"; then
            validation_failed=true
        fi
    done
    
    # Generate report
    generate_validation_report
    
    # Output results
    echo ""
    output_results
    
    # Exit with appropriate status
    local total_errors=$((VALIDATION_STATS[format_errors] + VALIDATION_STATS[category_errors]))
    if [[ $total_errors -gt 0 && "$VALIDATION_MODE" == "strict" ]]; then
        log_error "Validation failed with $total_errors errors"
        exit 1
    elif [[ $total_errors -gt 0 ]]; then
        log_warning "Validation completed with $total_errors errors"
        exit 0
    else
        log_success "Decision context validation completed successfully"
        exit 0
    fi
}

# Execute main function
main 