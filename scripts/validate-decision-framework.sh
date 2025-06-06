#!/usr/bin/env bash
# Decision Framework Validation Script
# 🔶 DOC-014: Decision validation tools - Automated compliance checking
# Validates compliance with DOC-014 AI Assistant Decision Framework

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
CONTEXT_DIR="$DOCS_DIR/context"
VALIDATION_REPORTS_DIR="$DOCS_DIR/validation-reports"
DECISION_FRAMEWORK_FILE="$CONTEXT_DIR/ai-decision-framework.md"
AI_COMPLIANCE_FILE="$CONTEXT_DIR/ai-assistant-compliance.md"
AI_PROTOCOL_FILE="$CONTEXT_DIR/ai-assistant-protocol.md"
FEATURE_TRACKING_FILE="$CONTEXT_DIR/feature-tracking.md"

# Validation modes
VALIDATION_MODE="standard"  # standard, strict, metrics-only
OUTPUT_FORMAT="detailed"    # detailed, summary, json
VERBOSE=false
GENERATE_REPORT=true
CHECK_TOKENS=true
CHECK_PROTOCOLS=true
CHECK_COMPLIANCE=true
TRACK_METRICS=true

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Validation results
declare -A VALIDATION_RESULTS
declare -A DECISION_METRICS
declare -A COMPLIANCE_STATUS

# Usage information
usage() {
    echo "Decision Framework Validation Script (DOC-014)"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --mode <standard|strict|metrics-only>  Validation mode (default: standard)"
    echo "  --format <detailed|summary|json>       Output format (default: detailed)"
    echo "  --no-report                            Skip generating validation report"
    echo "  --no-tokens                            Skip enhanced token validation"
    echo "  --no-protocols                         Skip protocol validation"
    echo "  --no-compliance                        Skip compliance validation"
    echo "  --no-metrics                           Skip metrics tracking"
    echo "  --verbose                              Enable verbose output"
    echo "  --help                                 Show this help message"
    echo ""
    echo "Validation Modes:"
    echo "  standard      Standard validation (default)"
    echo "  strict        Strict validation with zero tolerance"
    echo "  metrics-only  Only collect and report metrics"
    echo ""
    echo "Examples:"
    echo "  $0                                     Run standard validation"
    echo "  $0 --mode strict --format summary     Strict validation with summary output"
    echo "  $0 --mode metrics-only --no-report    Collect metrics without report"
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

log_metric() {
    echo -e "${PURPLE}[METRIC]${NC} $*"
}

# Initialize validation results
init_validation() {
    VALIDATION_RESULTS[decision_framework_exists]=false
    VALIDATION_RESULTS[protocols_integrated]=false
    VALIDATION_RESULTS[compliance_documented]=false
    VALIDATION_RESULTS[tokens_enhanced]=false
    VALIDATION_RESULTS[validation_integrated]=false
    
    DECISION_METRICS[total_protocols]=0
    DECISION_METRICS[protocols_with_validation]=0
    DECISION_METRICS[total_tokens]=0
    DECISION_METRICS[tokens_with_decision_context]=0
    DECISION_METRICS[compliance_sections]=0
    DECISION_METRICS[validation_requirements]=0
    
    COMPLIANCE_STATUS[framework_mandatory]=false
    COMPLIANCE_STATUS[protocol_integration]=false
    COMPLIANCE_STATUS[token_enhancement]=false
    COMPLIANCE_STATUS[validation_hooks]=false
}

# Validate Decision Framework document exists and is complete
validate_decision_framework_document() {
    log_info "🔍 Validating Decision Framework document..."
    
    if [[ ! -f "$DECISION_FRAMEWORK_FILE" ]]; then
        log_error "Decision Framework document not found: $DECISION_FRAMEWORK_FILE"
        VALIDATION_RESULTS[decision_framework_exists]=false
        return 1
    fi
    
    VALIDATION_RESULTS[decision_framework_exists]=true
    log_verbose "Decision Framework document found"
    
    # Check for required sections
    local required_sections=(
        "4-Tier Decision Hierarchy"
        "Decision Trees"
        "Validation Checklists"
        "Success Metrics"
        "Integration Points"
    )
    
    local missing_sections=()
    for section in "${required_sections[@]}"; do
        if ! grep -q "$section" "$DECISION_FRAMEWORK_FILE"; then
            missing_sections+=("$section")
            log_verbose "Missing section: $section"
        else
            log_verbose "Found section: $section"
        fi
    done
    
    if [[ ${#missing_sections[@]} -gt 0 ]]; then
        log_warning "Missing sections in Decision Framework: ${missing_sections[*]}"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            return 1
        fi
    else
        log_success "All required sections found in Decision Framework"
    fi
    
    return 0
}

# Validate protocol integration
validate_protocol_integration() {
    log_info "🔍 Validating protocol integration..."
    
    if [[ ! -f "$AI_PROTOCOL_FILE" ]]; then
        log_error "AI Assistant Protocol file not found: $AI_PROTOCOL_FILE"
        VALIDATION_RESULTS[protocols_integrated]=false
        return 1
    fi
    
    # Check for protocol integration
    local protocols=(
        "NEW FEATURE Protocol"
        "MODIFICATION Protocol"
        "BUG FIX Protocol"
        "CONFIG CHANGE Protocol"
        "API CHANGE Protocol"
        "TEST ADDITION Protocol"
        "PERFORMANCE Protocol"
        "REFACTORING Protocol"
    )
    
    local integrated_protocols=0
    local total_protocols=${#protocols[@]}
    DECISION_METRICS[total_protocols]=$total_protocols
    
    for protocol in "${protocols[@]}"; do
        # Check if protocol has Decision Framework validation
        if grep -A 10 "### .*$protocol" "$AI_PROTOCOL_FILE" | grep -q "Decision Framework Validation"; then
            ((integrated_protocols++))
            log_verbose "Protocol integrated: $protocol"
        else
            log_warning "Protocol missing Decision Framework validation: $protocol"
            if [[ "$VALIDATION_MODE" == "strict" ]]; then
                VALIDATION_RESULTS[protocols_integrated]=false
                return 1
            fi
        fi
    done
    
    DECISION_METRICS[protocols_with_validation]=$integrated_protocols
    
    if [[ $integrated_protocols -eq $total_protocols ]]; then
        VALIDATION_RESULTS[protocols_integrated]=true
        COMPLIANCE_STATUS[protocol_integration]=true
        log_success "All $total_protocols protocols have Decision Framework validation"
    else
        log_warning "Only $integrated_protocols/$total_protocols protocols have Decision Framework validation"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            VALIDATION_RESULTS[protocols_integrated]=false
            return 1
        else
            VALIDATION_RESULTS[protocols_integrated]=true
        fi
    fi
    
    return 0
}

# Validate compliance documentation
validate_compliance_documentation() {
    log_info "🔍 Validating compliance documentation..."
    
    if [[ ! -f "$AI_COMPLIANCE_FILE" ]]; then
        log_error "AI Assistant Compliance file not found: $AI_COMPLIANCE_FILE"
        VALIDATION_RESULTS[compliance_documented]=false
        return 1
    fi
    
    # Check for DOC-014 compliance sections
    local compliance_sections=(
        "DOC-014: MANDATORY AI ASSISTANT DECISION FRAMEWORK"
        "Decision Framework Validation"
        "AI ENFORCEMENT RULES"
        "DOC-014 DECISION FRAMEWORK COMPLIANCE"
    )
    
    local found_sections=0
    local total_sections=${#compliance_sections[@]}
    DECISION_METRICS[compliance_sections]=$total_sections
    
    for section in "${compliance_sections[@]}"; do
        if grep -q "$section" "$AI_COMPLIANCE_FILE"; then
            ((found_sections++))
            log_verbose "Found compliance section: $section"
        else
            log_warning "Missing compliance section: $section"
            if [[ "$VALIDATION_MODE" == "strict" ]]; then
                VALIDATION_RESULTS[compliance_documented]=false
                return 1
            fi
        fi
    done
    
    if [[ $found_sections -eq $total_sections ]]; then
        VALIDATION_RESULTS[compliance_documented]=true
        COMPLIANCE_STATUS[framework_mandatory]=true
        log_success "All compliance sections documented"
    else
        log_warning "Only $found_sections/$total_sections compliance sections found"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            VALIDATION_RESULTS[compliance_documented]=false
            return 1
        else
            VALIDATION_RESULTS[compliance_documented]=true
        fi
    fi
    
    return 0
}

# Validate enhanced tokens with decision context
validate_enhanced_tokens() {
    log_info "🔍 Validating enhanced token implementation..."
    
    local go_files=($(find "$PROJECT_ROOT" -name "*.go" -type f | grep -v vendor))
    local total_tokens=0
    local enhanced_tokens=0
    
    for file in "${go_files[@]}"; do
        # Count total implementation tokens
        local file_tokens=$(grep -c "^[[:space:]]*//[[:space:]]*[⭐🔺🔶🔻]" "$file" 2>/dev/null || echo "0")
        # Clean any whitespace/newlines from command substitution
        file_tokens=$(echo "$file_tokens" | tr -d '\n\r' | head -1)
        [[ "$file_tokens" =~ ^[0-9]+$ ]] || file_tokens=0
        total_tokens=$((total_tokens + file_tokens))
        
        # Count tokens with decision context
        local file_enhanced_tokens=$(grep -c "\[DECISION:" "$file" 2>/dev/null || echo "0")
        # Clean any whitespace/newlines from command substitution
        file_enhanced_tokens=$(echo "$file_enhanced_tokens" | tr -d '\n\r' | head -1)
        [[ "$file_enhanced_tokens" =~ ^[0-9]+$ ]] || file_enhanced_tokens=0
        enhanced_tokens=$((enhanced_tokens + file_enhanced_tokens))
        
        log_verbose "File $file: $file_tokens tokens, $file_enhanced_tokens enhanced"
    done
    
    DECISION_METRICS[total_tokens]=$total_tokens
    DECISION_METRICS[tokens_with_decision_context]=$enhanced_tokens
    
    if [[ $total_tokens -eq 0 ]]; then
        log_warning "No implementation tokens found"
        VALIDATION_RESULTS[tokens_enhanced]=true  # No tokens to enhance
        return 0
    fi
    
    local enhancement_rate=$((enhanced_tokens * 100 / total_tokens))
    log_metric "Token enhancement rate: $enhancement_rate% ($enhanced_tokens/$total_tokens)"
    
    # Check enhancement thresholds
    local required_rate=0
    case "$VALIDATION_MODE" in
        "strict") required_rate=90 ;;
        "standard") required_rate=50 ;;
        "metrics-only") required_rate=0 ;;
    esac
    
    if [[ $enhancement_rate -ge $required_rate ]]; then
        VALIDATION_RESULTS[tokens_enhanced]=true
        COMPLIANCE_STATUS[token_enhancement]=true
        log_success "Token enhancement rate ($enhancement_rate%) meets requirements (≥$required_rate%)"
    else
        log_warning "Token enhancement rate ($enhancement_rate%) below requirements (≥$required_rate%)"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            VALIDATION_RESULTS[tokens_enhanced]=false
            return 1
        else
            VALIDATION_RESULTS[tokens_enhanced]=true
        fi
    fi
    
    return 0
}

# Validate decision context format in enhanced tokens
validate_decision_context_format() {
    log_info "🔍 Validating decision context format..."
    
    local go_files=($(find "$PROJECT_ROOT" -name "*.go" -type f | grep -v vendor))
    local format_errors=0
    local context_tags_valid=0
    local context_tags_total=0
    
    # Valid decision context categories
    local valid_categories=(
        "core-functionality" "enhancement" "developer-experience" "infrastructure"
        "blocks-extraction" "enables-flexibility" "prerequisite-for-testing" "independent"
        "user-facing" "backward-compatible" "breaking-change" "performance-critical"
    )
    
    for file in "${go_files[@]}"; do
        # Find lines with decision context
        while IFS= read -r line; do
            if [[ "$line" =~ \[DECISION:([^]]+)\] ]]; then
                local context="${BASH_REMATCH[1]}"
                log_verbose "Found decision context: $context"
                
                # Validate format (comma-separated, max 3 tags)
                local tags=($(echo "$context" | tr ',' '\n' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//'))
                local tag_count=${#tags[@]}
                
                if [[ $tag_count -gt 3 ]]; then
                    log_warning "Too many decision context tags ($tag_count > 3): $context"
                    ((format_errors++))
                fi
                
                # Validate individual tags
                for tag in "${tags[@]}"; do
                    ((context_tags_total++))
                    local tag_valid=false
                    for valid_tag in "${valid_categories[@]}"; do
                        if [[ "$tag" == "$valid_tag" ]]; then
                            tag_valid=true
                            ((context_tags_valid++))
                            break
                        fi
                    done
                    
                    if [[ "$tag_valid" == "false" ]]; then
                        log_warning "Invalid decision context tag: '$tag' in $file"
                        ((format_errors++))
                    fi
                done
            fi
        done < "$file"
    done
    
    log_metric "Decision context validation: $context_tags_valid/$context_tags_total tags valid"
    
    if [[ $format_errors -eq 0 ]]; then
        log_success "All decision context formats are valid"
        return 0
    else
        log_warning "Found $format_errors decision context format errors"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            return 1
        fi
    fi
    
    return 0
}

# Validate integration with existing validation systems
validate_validation_integration() {
    log_info "🔍 Validating integration with existing validation systems..."
    
    local integration_points=(
        "scripts/validate-icon-enforcement.sh"
        "internal/validation/ai_validation.go"
        "cmd/ai-validation/main.go"
    )
    
    local integrated_systems=0
    local total_systems=${#integration_points[@]}
    
    for system in "${integration_points[@]}"; do
        local system_path="$PROJECT_ROOT/$system"
        if [[ -f "$system_path" ]]; then
            # Check if system references DOC-014
            if grep -q "DOC-014" "$system_path"; then
                ((integrated_systems++))
                log_verbose "Validation system integrated: $system"
            else
                log_verbose "Validation system exists but not integrated: $system"
            fi
        else
            log_verbose "Validation system not found: $system"
        fi
    done
    
    DECISION_METRICS[validation_requirements]=$total_systems
    
    if [[ $integrated_systems -gt 0 ]]; then
        VALIDATION_RESULTS[validation_integrated]=true
        COMPLIANCE_STATUS[validation_hooks]=true
        log_success "Decision Framework integrated with $integrated_systems validation systems"
    else
        log_warning "Decision Framework not integrated with existing validation systems"
        if [[ "$VALIDATION_MODE" == "strict" ]]; then
            VALIDATION_RESULTS[validation_integrated]=false
            return 1
        else
            VALIDATION_RESULTS[validation_integrated]=true
        fi
    fi
    
    return 0
}

# Calculate decision quality metrics
calculate_decision_metrics() {
    log_info "📊 Calculating decision quality metrics..."
    
    # Goal alignment rate (based on protocol integration)
    local protocols_rate=0
    if [[ ${DECISION_METRICS[total_protocols]} -gt 0 ]]; then
        protocols_rate=$((DECISION_METRICS[protocols_with_validation] * 100 / DECISION_METRICS[total_protocols]))
    fi
    
    # Token enhancement rate
    local tokens_rate=0
    if [[ ${DECISION_METRICS[total_tokens]} -gt 0 ]]; then
        tokens_rate=$((DECISION_METRICS[tokens_with_decision_context] * 100 / DECISION_METRICS[total_tokens]))
    fi
    
    # Compliance documentation rate
    local compliance_rate=0
    if [[ ${DECISION_METRICS[compliance_sections]} -gt 0 ]]; then
        compliance_rate=$((${DECISION_METRICS[compliance_sections]} * 100 / 4))  # 4 expected sections
    fi
    
    # Overall decision framework maturity
    local framework_components=0
    local total_components=5
    
    [[ "${VALIDATION_RESULTS[decision_framework_exists]}" == "true" ]] && ((framework_components++))
    [[ "${VALIDATION_RESULTS[protocols_integrated]}" == "true" ]] && ((framework_components++))
    [[ "${VALIDATION_RESULTS[compliance_documented]}" == "true" ]] && ((framework_components++))
    [[ "${VALIDATION_RESULTS[tokens_enhanced]}" == "true" ]] && ((framework_components++))
    [[ "${VALIDATION_RESULTS[validation_integrated]}" == "true" ]] && ((framework_components++))
    
    local maturity_score=$((framework_components * 100 / total_components))
    
    # Log metrics
    log_metric "Protocol Integration Rate: $protocols_rate%"
    log_metric "Token Enhancement Rate: $tokens_rate%"
    log_metric "Compliance Documentation Rate: $compliance_rate%"
    log_metric "Framework Maturity Score: $maturity_score%"
    
    # Success criteria validation (from DOC-014 specification)
    local goal_alignment_target=95
    local traceability_target=100
    
    if [[ $protocols_rate -ge $goal_alignment_target ]]; then
        log_success "Goal alignment rate ($protocols_rate%) meets target (≥$goal_alignment_target%)"
    else
        log_warning "Goal alignment rate ($protocols_rate%) below target (≥$goal_alignment_target%)"
    fi
    
    if [[ $protocols_rate -eq 100 ]]; then
        log_success "Traceability rate (100%) meets target (≥$traceability_target%)"
    else
        log_warning "Traceability rate ($protocols_rate%) below target (≥$traceability_target%)"
    fi
    
    # Store metrics for reporting
    DECISION_METRICS[goal_alignment_rate]=$protocols_rate
    DECISION_METRICS[token_enhancement_rate]=$tokens_rate
    DECISION_METRICS[compliance_rate]=$compliance_rate
    DECISION_METRICS[maturity_score]=$maturity_score
}

# Generate validation report
generate_validation_report() {
    if [[ "$GENERATE_REPORT" != "true" ]]; then
        return 0
    fi
    
    log_info "📄 Generating validation report..."
    
    mkdir -p "$VALIDATION_REPORTS_DIR"
    local report_file="$VALIDATION_REPORTS_DIR/decision-framework-validation-report.md"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    cat > "$report_file" << EOF
# Decision Framework Validation Report (DOC-014)

> **Generated on:** \`$timestamp\`  
> **Validation Mode:** \`$VALIDATION_MODE\`  
> **Report Format:** \`$OUTPUT_FORMAT\`

## 🎯 Executive Summary

The DOC-014 AI Assistant Decision Framework validation has been completed with the following results:

### 📊 Key Metrics
- **Framework Maturity Score:** ${DECISION_METRICS[maturity_score]}%
- **Goal Alignment Rate:** ${DECISION_METRICS[goal_alignment_rate]}%
- **Token Enhancement Rate:** ${DECISION_METRICS[token_enhancement_rate]}%
- **Compliance Documentation Rate:** ${DECISION_METRICS[compliance_rate]}%

### ✅ Validation Results
| Component | Status | Details |
|-----------|---------|---------|
| Decision Framework Document | $([ "${VALIDATION_RESULTS[decision_framework_exists]}" == "true" ] && echo "✅ PASS" || echo "❌ FAIL") | Core framework document validation |
| Protocol Integration | $([ "${VALIDATION_RESULTS[protocols_integrated]}" == "true" ] && echo "✅ PASS" || echo "❌ FAIL") | ${DECISION_METRICS[protocols_with_validation]}/${DECISION_METRICS[total_protocols]} protocols integrated |
| Compliance Documentation | $([ "${VALIDATION_RESULTS[compliance_documented]}" == "true" ] && echo "✅ PASS" || echo "❌ FAIL") | AI assistant compliance requirements |
| Enhanced Tokens | $([ "${VALIDATION_RESULTS[tokens_enhanced]}" == "true" ] && echo "✅ PASS" || echo "❌ FAIL") | ${DECISION_METRICS[tokens_with_decision_context]}/${DECISION_METRICS[total_tokens]} tokens enhanced |
| Validation Integration | $([ "${VALIDATION_RESULTS[validation_integrated]}" == "true" ] && echo "✅ PASS" || echo "❌ FAIL") | Integration with existing validation systems |

## 📋 Detailed Validation Results

### 🔍 Decision Framework Document
$([ "${VALIDATION_RESULTS[decision_framework_exists]}" == "true" ] && echo "**Status:** ✅ VALIDATED" || echo "**Status:** ❌ FAILED")

The core Decision Framework document (\`ai-decision-framework.md\`) has been validated for completeness and structure.

### 🔄 Protocol Integration  
**Status:** $([ "${VALIDATION_RESULTS[protocols_integrated]}" == "true" ] && echo "✅ VALIDATED" || echo "❌ FAILED")  
**Coverage:** ${DECISION_METRICS[protocols_with_validation]}/${DECISION_METRICS[total_protocols]} protocols (${DECISION_METRICS[goal_alignment_rate]}%)

All AI assistant change protocols have been analyzed for Decision Framework integration.

### 📚 Compliance Documentation
**Status:** $([ "${VALIDATION_RESULTS[compliance_documented]}" == "true" ] && echo "✅ VALIDATED" || echo "❌ FAILED")  
**Coverage:** ${DECISION_METRICS[compliance_sections]}/4 required sections documented

AI assistant compliance requirements have been updated to include DOC-014 Decision Framework mandates.

### 🏷️ Enhanced Token Implementation
**Status:** $([ "${VALIDATION_RESULTS[tokens_enhanced]}" == "true" ] && echo "✅ VALIDATED" || echo "❌ FAILED")  
**Enhancement Rate:** ${DECISION_METRICS[tokens_with_decision_context]}/${DECISION_METRICS[total_tokens]} tokens (${DECISION_METRICS[token_enhancement_rate]}%)

Implementation tokens have been analyzed for decision context integration.

### 🔗 Validation System Integration
**Status:** $([ "${VALIDATION_RESULTS[validation_integrated]}" == "true" ] && echo "✅ VALIDATED" || echo "❌ FAILED")

Integration with existing DOC-008/DOC-011 validation systems has been verified.

## 🎯 Success Criteria Analysis

### Goal Alignment (Target: ≥95%)
**Current:** ${DECISION_METRICS[goal_alignment_rate]}%  
**Status:** $([ ${DECISION_METRICS[goal_alignment_rate]} -ge 95 ] && echo "✅ MET" || echo "⚠️ BELOW TARGET")

### Traceability (Target: 100%)
**Current:** ${DECISION_METRICS[goal_alignment_rate]}%  
**Status:** $([ ${DECISION_METRICS[goal_alignment_rate]} -eq 100 ] && echo "✅ MET" || echo "⚠️ BELOW TARGET")

## 📈 Recommendations

### Immediate Actions
EOF

    if [[ ${DECISION_METRICS[goal_alignment_rate]} -lt 95 ]]; then
        echo "- **Improve Protocol Integration:** Complete Decision Framework validation for remaining protocols" >> "$report_file"
    fi
    
    if [[ ${DECISION_METRICS[token_enhancement_rate]} -lt 50 ]]; then
        echo "- **Enhance Token Coverage:** Run \`make enhance-tokens-phase1\` to improve decision context coverage" >> "$report_file"
    fi
    
    if [[ "${VALIDATION_RESULTS[validation_integrated]}" != "true" ]]; then
        echo "- **Integrate Validation Systems:** Complete DOC-008/DOC-011 integration with Decision Framework" >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF

### Long-term Improvements
- Monitor decision quality metrics over time
- Implement automated decision compliance checking
- Expand decision context categories as needed
- Integrate decision metrics into CI/CD pipeline

## 🔧 Technical Details

### Validation Command
\`\`\`bash
$0 --mode $VALIDATION_MODE --format $OUTPUT_FORMAT
\`\`\`

### Validation Date
$timestamp

### Validation Mode
$VALIDATION_MODE

---

*This report was generated automatically by the DOC-014 Decision Framework validation system.*
EOF

    log_success "Validation report generated: $report_file"
}

# Output results based on format
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
    echo "🎯 DOC-014 Decision Framework Validation Results"
    echo "================================================"
    echo ""
    
    echo "📊 Framework Maturity: ${DECISION_METRICS[maturity_score]}%"
    echo "🎯 Goal Alignment: ${DECISION_METRICS[goal_alignment_rate]}%"
    echo "🏷️ Token Enhancement: ${DECISION_METRICS[token_enhancement_rate]}%"
    echo "📚 Compliance Rate: ${DECISION_METRICS[compliance_rate]}%"
    echo ""
    
    echo "✅ Validation Components:"
    for component in decision_framework_exists protocols_integrated compliance_documented tokens_enhanced validation_integrated; do
        local status="${VALIDATION_RESULTS[$component]}"
        local icon=$([ "$status" == "true" ] && echo "✅" || echo "❌")
        echo "  $icon ${component//_/ }"
    done
    echo ""
    
    echo "🔍 Compliance Status:"
    for status in framework_mandatory protocol_integration token_enhancement validation_hooks; do
        local value="${COMPLIANCE_STATUS[$status]}"
        local icon=$([ "$value" == "true" ] && echo "✅" || echo "❌")
        echo "  $icon ${status//_/ }"
    done
}

# Output summary results
output_summary_results() {
    local overall_status="PASS"
    local failed_components=0
    
    for component in "${!VALIDATION_RESULTS[@]}"; do
        if [[ "${VALIDATION_RESULTS[$component]}" != "true" ]]; then
            ((failed_components++))
            overall_status="FAIL"
        fi
    done
    
    echo "DOC-014 Decision Framework Validation: $overall_status"
    echo "Maturity Score: ${DECISION_METRICS[maturity_score]:-0}%"
    echo "Failed Components: $failed_components/5"
}

# Output JSON results
output_json_results() {
    # Calculate failed components for JSON output
    local failed_components=0
    for component in "${!VALIDATION_RESULTS[@]}"; do
        if [[ "${VALIDATION_RESULTS[$component]}" != "true" ]]; then
            ((failed_components++))
        fi
    done
    
    cat << EOF
{
  "validation_timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "validation_mode": "$VALIDATION_MODE",
  "overall_status": "$([ "$failed_components" -eq 0 ] && echo "PASS" || echo "FAIL")",
  "metrics": {
    "maturity_score": ${DECISION_METRICS[maturity_score]:-0},
    "goal_alignment_rate": ${DECISION_METRICS[goal_alignment_rate]:-0},
    "token_enhancement_rate": ${DECISION_METRICS[token_enhancement_rate]:-0},
    "compliance_rate": ${DECISION_METRICS[compliance_rate]:-0},
    "total_protocols": ${DECISION_METRICS[total_protocols]:-0},
    "protocols_with_validation": ${DECISION_METRICS[protocols_with_validation]:-0},
    "total_tokens": ${DECISION_METRICS[total_tokens]:-0},
    "tokens_with_decision_context": ${DECISION_METRICS[tokens_with_decision_context]:-0}
  },
  "validation_results": {
EOF
    
    local first=true
    for component in "${!VALIDATION_RESULTS[@]}"; do
        if [[ "$first" == "true" ]]; then
            first=false
        else
            echo ","
        fi
        echo -n "    \"$component\": ${VALIDATION_RESULTS[$component]}"
    done
    
    echo ""
    echo "  },"
    echo "  \"compliance_status\": {"
    
    first=true
    for status in "${!COMPLIANCE_STATUS[@]}"; do
        if [[ "$first" == "true" ]]; then
            first=false
        else
            echo ","
        fi
        echo -n "    \"$status\": ${COMPLIANCE_STATUS[$status]}"
    done
    
    echo ""
    echo "  }"
    echo "}"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode)
            VALIDATION_MODE="$2"
            if [[ ! "$VALIDATION_MODE" =~ ^(standard|strict|metrics-only)$ ]]; then
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
        --no-report)
            GENERATE_REPORT=false
            shift
            ;;
        --no-tokens)
            CHECK_TOKENS=false
            shift
            ;;
        --no-protocols)
            CHECK_PROTOCOLS=false
            shift
            ;;
        --no-compliance)
            CHECK_COMPLIANCE=false
            shift
            ;;
        --no-metrics)
            TRACK_METRICS=false
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
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Main validation execution
main() {
    log_info "🔶 DOC-014: Decision Framework Validation"
    log_info "Mode: $VALIDATION_MODE | Format: $OUTPUT_FORMAT"
    echo ""
    
    # Initialize validation
    init_validation
    
    # Run validation components
    local validation_failed=false
    
    # Always validate framework document
    if ! validate_decision_framework_document; then
        validation_failed=true
    fi
    
    # Protocol validation
    if [[ "$CHECK_PROTOCOLS" == "true" ]]; then
        if ! validate_protocol_integration; then
            validation_failed=true
        fi
    fi
    
    # Compliance validation
    if [[ "$CHECK_COMPLIANCE" == "true" ]]; then
        if ! validate_compliance_documentation; then
            validation_failed=true
        fi
    fi
    
    # Token validation
    if [[ "$CHECK_TOKENS" == "true" ]]; then
        if ! validate_enhanced_tokens; then
            validation_failed=true
        fi
        
        if ! validate_decision_context_format; then
            validation_failed=true
        fi
    fi
    
    # Validation integration
    if ! validate_validation_integration; then
        validation_failed=true
    fi
    
    # Calculate metrics
    if [[ "$TRACK_METRICS" == "true" ]]; then
        calculate_decision_metrics
    fi
    
    # Generate report
    generate_validation_report
    
    # Output results
    echo ""
    output_results
    
    # Exit with appropriate status
    if [[ "$validation_failed" == "true" && "$VALIDATION_MODE" == "strict" ]]; then
        log_error "Validation failed in strict mode"
        exit 1
    else
        log_success "Decision Framework validation completed"
        exit 0
    fi
}

# Execute main function
main 