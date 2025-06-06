#!/usr/bin/env bash
# Decision Quality Metrics Tracker
# 🔶 DOC-014: Decision validation tools - Decision quality metrics tracking
# Tracks and analyzes decision quality metrics for AI Assistant Decision Framework

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
METRICS_DIR="$DOCS_DIR/decision-metrics"
HISTORY_FILE="$METRICS_DIR/decision-metrics-history.json"
CONFIG_FILE="$METRICS_DIR/metrics-config.toml"

# Default configuration
COLLECTION_MODE="full"      # full, incremental, analysis-only
OUTPUT_FORMAT="detailed"    # detailed, summary, json, dashboard
TIME_RANGE="30d"           # 7d, 30d, 90d, all
TRACK_TRENDS=true
GENERATE_DASHBOARD=true
VERBOSE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Metrics data
declare -A CURRENT_METRICS
declare -A TREND_DATA
declare -A QUALITY_INDICATORS

# Usage information
usage() {
    echo "Decision Quality Metrics Tracker (DOC-014)"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --mode <full|incremental|analysis-only>  Collection mode (default: full)"
    echo "  --format <detailed|summary|json|dashboard>  Output format (default: detailed)"
    echo "  --time-range <7d|30d|90d|all>           Time range for analysis (default: 30d)"
    echo "  --no-trends                              Skip trend analysis"
    echo "  --no-dashboard                           Skip dashboard generation"
    echo "  --verbose                                Enable verbose output"
    echo "  --help                                   Show this help message"
    echo ""
    echo "Collection Modes:"
    echo "  full          Collect all metrics and analyze trends"
    echo "  incremental   Update existing metrics with new data"
    echo "  analysis-only Only analyze existing data without collection"
    echo ""
    echo "Examples:"
    echo "  $0                                       Full metrics collection and analysis"
    echo "  $0 --mode incremental --format summary  Incremental update with summary"
    echo "  $0 --mode analysis-only --time-range 7d Short-term analysis only"
}

# Logging functions
log_info() {
    [[ "$OUTPUT_FORMAT" != "json" ]] && echo -e "${BLUE}[INFO]${NC} $*"
}

log_success() {
    [[ "$OUTPUT_FORMAT" != "json" ]] && echo -e "${GREEN}[SUCCESS]${NC} $*"
}

log_warning() {
    [[ "$OUTPUT_FORMAT" != "json" ]] && echo -e "${YELLOW}[WARNING]${NC} $*"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $*" >&2
}

log_verbose() {
    if [[ "$VERBOSE" == "true" && "$OUTPUT_FORMAT" != "json" ]]; then
        echo -e "${CYAN}[VERBOSE]${NC} $*"
    fi
}

log_metric() {
    [[ "$OUTPUT_FORMAT" != "json" ]] && echo -e "${PURPLE}[METRIC]${NC} $*"
}

# Initialize metrics tracking
init_metrics() {
    mkdir -p "$METRICS_DIR"
    
    # Initialize metrics arrays
    CURRENT_METRICS[goal_alignment_rate]=0
    CURRENT_METRICS[protocol_compliance_rate]=0
    CURRENT_METRICS[token_enhancement_rate]=0
    CURRENT_METRICS[rework_rate]=0
    CURRENT_METRICS[decision_compliance_rate]=0
    CURRENT_METRICS[validation_pass_rate]=0
    CURRENT_METRICS[framework_adoption_rate]=0
    
    # Quality indicators
    QUALITY_INDICATORS[excellent_threshold]=95
    QUALITY_INDICATORS[good_threshold]=85
    QUALITY_INDICATORS[acceptable_threshold]=70
    QUALITY_INDICATORS[poor_threshold]=50
    
    # Create configuration file if it doesn't exist
    if [[ ! -f "$CONFIG_FILE" ]]; then
        create_metrics_config
    fi
}

# Create metrics configuration
create_metrics_config() {
    log_info "Creating metrics configuration..."
    
    cat > "$CONFIG_FILE" << 'EOF'
# Decision Quality Metrics Configuration
# 🔶 DOC-014: Decision validation tools - Metrics tracking configuration

[thresholds]
# Quality thresholds for decision metrics (percentages)
excellent = 95
good = 85
acceptable = 70
poor = 50

[targets]
# Target values for DOC-014 success criteria
goal_alignment_rate = 95
traceability_rate = 100
rework_rate_max = 5
decision_compliance_rate = 90

[collection]
# Metrics collection settings
update_frequency = "daily"
retention_period = "1y"
auto_cleanup = true
trend_analysis_window = 30

[alerts]
# Alert conditions for quality degradation
goal_alignment_below = 90
rework_rate_above = 10
compliance_below = 85
consecutive_failures = 3

[dashboard]
# Dashboard generation settings
generate_charts = true
include_trends = true
historical_comparison = true
export_formats = ["html", "json"]
EOF

    log_verbose "Metrics configuration created: $CONFIG_FILE"
}

# Collect current decision metrics
collect_current_metrics() {
    log_info "📊 Collecting current decision metrics..."
    
    # Run decision framework validation to get current metrics
    local validation_output
    if validation_output=$("$SCRIPT_DIR/validate-decision-framework.sh" --format json 2>/dev/null); then
        log_verbose "Decision framework validation completed successfully"
        
        # Parse JSON output to extract metrics
        if command -v jq >/dev/null 2>&1; then
            CURRENT_METRICS[goal_alignment_rate]=$(echo "$validation_output" | jq -r '.metrics.goal_alignment_rate // 0')
            CURRENT_METRICS[protocol_compliance_rate]=$(echo "$validation_output" | jq -r '.metrics.goal_alignment_rate // 0')
            CURRENT_METRICS[token_enhancement_rate]=$(echo "$validation_output" | jq -r '.metrics.token_enhancement_rate // 0')
            CURRENT_METRICS[framework_maturity]=$(echo "$validation_output" | jq -r '.metrics.maturity_score // 0')
        else
            log_warning "jq not available, using fallback metric extraction"
            extract_metrics_fallback "$validation_output"
        fi
    else
        log_warning "Decision framework validation failed, using estimated metrics"
        estimate_current_metrics
    fi
    
    # Calculate additional metrics
    calculate_rework_rate
    calculate_decision_compliance_rate
    calculate_validation_pass_rate
    calculate_framework_adoption_rate
    calculate_protocol_coverage
    
    log_verbose "Current metrics collected successfully"
}

# Fallback method for extracting metrics without jq
extract_metrics_fallback() {
    local json_output="$1"
    
    # Extract metrics using grep and sed
    CURRENT_METRICS[goal_alignment_rate]=$(echo "$json_output" | grep -o '"goal_alignment_rate": [0-9]*' | sed 's/.*: //' || echo "0")
    CURRENT_METRICS[token_enhancement_rate]=$(echo "$json_output" | grep -o '"token_enhancement_rate": [0-9]*' | sed 's/.*: //' || echo "0")
    CURRENT_METRICS[framework_maturity]=$(echo "$json_output" | grep -o '"maturity_score": [0-9]*' | sed 's/.*: //' || echo "0")
    
    log_verbose "Metrics extracted using fallback method"
}

# Estimate metrics when validation fails
estimate_current_metrics() {
    log_verbose "Estimating current metrics..."
    
    # Count protocols with Decision Framework validation
    local protocols_with_validation=0
    local total_protocols=8
    
    if [[ -f "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" ]]; then
        protocols_with_validation=$(grep -c "Decision Framework Validation" "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" || echo "0")
    fi
    
    CURRENT_METRICS[goal_alignment_rate]=$((protocols_with_validation * 100 / total_protocols))
    CURRENT_METRICS[protocol_compliance_rate]=$((protocols_with_validation * 100 / total_protocols))
    
    # Count enhanced tokens
    local enhanced_tokens=0
    local total_tokens=0
    
    while IFS= read -r -d '' file; do
        local file_tokens
        local file_enhanced
        file_tokens=$(grep -c "^[[:space:]]*//[[:space:]]*[⭐🔺🔶🔻]" "$file" 2>/dev/null || echo "0")
        file_enhanced=$(grep -c "\[DECISION:" "$file" 2>/dev/null || echo "0")
        # Clean any whitespace/newlines from command substitution
        file_tokens=$(echo "$file_tokens" | tr -d '\n\r' | head -1)
        file_enhanced=$(echo "$file_enhanced" | tr -d '\n\r' | head -1)
        # Ensure we have valid numbers
        [[ "$file_tokens" =~ ^[0-9]+$ ]] || file_tokens=0
        [[ "$file_enhanced" =~ ^[0-9]+$ ]] || file_enhanced=0
        total_tokens=$((total_tokens + file_tokens))
        enhanced_tokens=$((enhanced_tokens + file_enhanced))
    done < <(find "$PROJECT_ROOT" -name "*.go" -type f -print0)
    
    if [[ $total_tokens -gt 0 ]]; then
        CURRENT_METRICS[token_enhancement_rate]=$((enhanced_tokens * 100 / total_tokens))
    else
        CURRENT_METRICS[token_enhancement_rate]=0
    fi
    
    log_verbose "Current metrics estimated from direct analysis"
}

# Calculate rework rate based on git history
calculate_rework_rate() {
    log_verbose "Calculating rework rate..."
    
    # Analyze recent commits for rework indicators
    local recent_commits
    local rework_commits=0
    local total_commits=0
    
    if recent_commits=$(git log --oneline --since="30 days ago" 2>/dev/null); then
        total_commits=$(echo "$recent_commits" | wc -l | tr -d ' \n\r')
        [[ "$total_commits" =~ ^[0-9]+$ ]] || total_commits=0
        
        # Look for rework indicators in commit messages
        rework_commits=$(echo "$recent_commits" | grep -i -E "(fix|rework|correct|undo|revert)" | wc -l | tr -d ' \n\r')
        [[ "$rework_commits" =~ ^[0-9]+$ ]] || rework_commits=0
        
        if [[ $total_commits -gt 0 ]]; then
            CURRENT_METRICS[rework_rate]=$((rework_commits * 100 / total_commits))
        else
            CURRENT_METRICS[rework_rate]=0
        fi
    else
        log_warning "Git history not available, estimating rework rate"
        CURRENT_METRICS[rework_rate]=5  # Conservative estimate
    fi
    
    log_verbose "Rework rate: ${CURRENT_METRICS[rework_rate]}% ($rework_commits/$total_commits commits)"
}

# Calculate decision compliance rate
calculate_decision_compliance_rate() {
    log_verbose "Calculating decision compliance rate..."
    
    # Check for compliance with decision framework requirements
    local compliance_score=0
    local total_checks=5
    
    # Check 1: Decision framework document exists
    [[ -f "$PROJECT_ROOT/docs/context/ai-decision-framework.md" ]] && ((compliance_score++))
    
    # Check 2: AI compliance updated
    if [[ -f "$PROJECT_ROOT/docs/context/ai-assistant-compliance.md" ]]; then
        grep -q "DOC-014" "$PROJECT_ROOT/docs/context/ai-assistant-compliance.md" && ((compliance_score++))
    fi
    
    # Check 3: Protocols integrated
    if [[ -f "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" ]]; then
        local integrated_protocols=$(grep -c "Decision Framework Validation" "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" 2>/dev/null | tr -d ' \n\r' || echo "0")
        [[ "$integrated_protocols" =~ ^[0-9]+$ ]] || integrated_protocols=0
        [[ $integrated_protocols -ge 6 ]] && ((compliance_score++))  # At least 75% of protocols
    fi
    
    # Check 4: Enhanced tokens present
    local enhanced_tokens=$(find "$PROJECT_ROOT" -name "*.go" -exec grep -l "\[DECISION:" {} \; 2>/dev/null | wc -l | tr -d ' \n\r')
    [[ "$enhanced_tokens" =~ ^[0-9]+$ ]] || enhanced_tokens=0
    [[ $enhanced_tokens -gt 0 ]] && ((compliance_score++))
    
    # Check 5: Validation tools present
    [[ -f "$SCRIPT_DIR/validate-decision-framework.sh" ]] && ((compliance_score++))
    
    CURRENT_METRICS[decision_compliance_rate]=$((compliance_score * 100 / total_checks))
    
    log_verbose "Decision compliance rate: ${CURRENT_METRICS[decision_compliance_rate]}% ($compliance_score/$total_checks checks)"
}

# Calculate validation pass rate
calculate_validation_pass_rate() {
    log_verbose "Calculating validation pass rate..."
    
    # Run validation in non-verbose mode to get pass/fail status
    local validation_exit_code=0
    "$SCRIPT_DIR/validate-decision-framework.sh" --format summary >/dev/null 2>&1 || validation_exit_code=$?
    
    if [[ $validation_exit_code -eq 0 ]]; then
        CURRENT_METRICS[validation_pass_rate]=100
        log_verbose "Current validation: PASS"
    else
        # Partial validation based on component status
        local passing_components=0
        local total_components=5
        
        # Check each component individually
        [[ -f "$PROJECT_ROOT/docs/context/ai-decision-framework.md" ]] && ((passing_components++))
        
        local protocols_integrated=$(grep -c "Decision Framework Validation" "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" 2>/dev/null | tr -d ' \n\r' || echo "0")
        [[ "$protocols_integrated" =~ ^[0-9]+$ ]] || protocols_integrated=0
        [[ $protocols_integrated -ge 4 ]] && ((passing_components++))
        
        grep -q "DOC-014" "$PROJECT_ROOT/docs/context/ai-assistant-compliance.md" 2>/dev/null && ((passing_components++))
        
        local enhanced_tokens=$(find "$PROJECT_ROOT" -name "*.go" -exec grep -l "\[DECISION:" {} \; 2>/dev/null | wc -l | tr -d ' \n\r' || echo "0")
        [[ "$enhanced_tokens" =~ ^[0-9]+$ ]] || enhanced_tokens=0
        [[ $enhanced_tokens -gt 0 ]] && ((passing_components++))
        
        [[ -f "$SCRIPT_DIR/validate-decision-framework.sh" ]] && ((passing_components++))
        
        CURRENT_METRICS[validation_pass_rate]=$((passing_components * 100 / total_components))
        log_verbose "Current validation: PARTIAL ($passing_components/$total_components components)"
    fi
}

# Calculate framework adoption rate
calculate_framework_adoption_rate() {
    log_verbose "Calculating framework adoption rate..."
    
    # Measure adoption across different areas
    local adoption_score=0
    local total_areas=6
    
    # Area 1: Core documentation complete
    [[ -f "$PROJECT_ROOT/docs/context/ai-decision-framework.md" ]] && ((adoption_score++))
    
    # Area 2: Feature tracking integration
    if grep -q "DOC-014" "$PROJECT_ROOT/docs/context/feature-tracking.md" 2>/dev/null; then
        ((adoption_score++))
    fi
    
    # Area 3: Protocol integration (>75%)
    local integrated_protocols=$(grep -c "Decision Framework Validation" "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" 2>/dev/null | tr -d ' \n\r' || echo "0")
    [[ "$integrated_protocols" =~ ^[0-9]+$ ]] || integrated_protocols=0
    [[ $integrated_protocols -ge 6 ]] && ((adoption_score++))
    
    # Area 4: Compliance requirements updated
    grep -q "DOC-014" "$PROJECT_ROOT/docs/context/ai-assistant-compliance.md" 2>/dev/null && ((adoption_score++))
    
    # Area 5: Enhanced tokens in use
    local enhanced_tokens=$(find "$PROJECT_ROOT" -name "*.go" -exec grep -l "\[DECISION:" {} \; 2>/dev/null | wc -l | tr -d ' \n\r' || echo "0")
    [[ "$enhanced_tokens" =~ ^[0-9]+$ ]] || enhanced_tokens=0
    [[ $enhanced_tokens -gt 0 ]] && ((adoption_score++))
    
    # Area 6: Validation tools operational
    [[ -f "$SCRIPT_DIR/validate-decision-framework.sh" && -x "$SCRIPT_DIR/validate-decision-framework.sh" ]] && ((adoption_score++))
    
    CURRENT_METRICS[framework_adoption_rate]=$((adoption_score * 100 / total_areas))
    
    log_verbose "Framework adoption rate: ${CURRENT_METRICS[framework_adoption_rate]}% ($adoption_score/$total_areas areas)"
}

# Calculate protocol coverage metrics
calculate_protocol_coverage() {
    log_verbose "Calculating protocol coverage metrics..."
    
    # Define expected protocols
    local protocols=("NEW_FEATURE" "MODIFICATION" "BUG_FIX" "CONFIG_CHANGE" "API_CHANGE" "TEST_ADDITION" "PERFORMANCE" "REFACTORING")
    local total_protocols=${#protocols[@]}
    local covered_protocols=0
    
    # Check protocol integration in ai-assistant-protocol.md
    if [[ -f "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" ]]; then
        for protocol in "${protocols[@]}"; do
            if grep -q "Decision Framework Validation" "$PROJECT_ROOT/docs/context/ai-assistant-protocol.md" 2>/dev/null; then
                ((covered_protocols++))
            fi
        done
    fi
    
    # Calculate overall coverage percentage
    if [[ $total_protocols -gt 0 ]]; then
        CURRENT_METRICS[protocol_coverage]=$((covered_protocols * 100 / total_protocols))
    else
        CURRENT_METRICS[protocol_coverage]=0
    fi
    
    log_verbose "Protocol coverage: ${CURRENT_METRICS[protocol_coverage]}% ($covered_protocols/$total_protocols protocols)"
}

# Store metrics in history
store_metrics_history() {
    log_info "💾 Storing metrics in history..."
    
    local timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    local metrics_entry
    
    # Create JSON entry for current metrics
    metrics_entry=$(cat << EOF
{
  "timestamp": "$timestamp",
  "goal_alignment_rate": ${CURRENT_METRICS[goal_alignment_rate]},
  "protocol_compliance_rate": ${CURRENT_METRICS[protocol_compliance_rate]},
  "token_enhancement_rate": ${CURRENT_METRICS[token_enhancement_rate]},
  "rework_rate": ${CURRENT_METRICS[rework_rate]},
  "decision_compliance_rate": ${CURRENT_METRICS[decision_compliance_rate]},
  "validation_pass_rate": ${CURRENT_METRICS[validation_pass_rate]},
  "framework_adoption_rate": ${CURRENT_METRICS[framework_adoption_rate]},
  "framework_maturity": ${CURRENT_METRICS[framework_maturity]:-0}
}
EOF
)
    
    # Initialize history file if it doesn't exist
    if [[ ! -f "$HISTORY_FILE" ]]; then
        echo '{"metrics_history": []}' > "$HISTORY_FILE"
    fi
    
    # Add new entry to history
    if command -v jq >/dev/null 2>&1; then
        local temp_file
        temp_file=$(mktemp)
        jq ".metrics_history += [$metrics_entry]" "$HISTORY_FILE" > "$temp_file"
        mv "$temp_file" "$HISTORY_FILE"
        log_verbose "Metrics stored using jq"
    else
        # Fallback without jq
        local temp_file
        temp_file=$(mktemp)
        head -n -1 "$HISTORY_FILE" > "$temp_file"
        echo "  ," >> "$temp_file"
        echo "  $metrics_entry" >> "$temp_file"
        echo "]}" >> "$temp_file"
        mv "$temp_file" "$HISTORY_FILE"
        log_verbose "Metrics stored using fallback method"
    fi
    
    log_success "Metrics history updated: $HISTORY_FILE"
}

# Analyze trends from historical data
analyze_trends() {
    if [[ "$TRACK_TRENDS" != "true" || ! -f "$HISTORY_FILE" ]]; then
        return 0
    fi
    
    log_info "📈 Analyzing trends..."
    
    # Simple trend analysis without jq dependency
    local recent_entries
    local trend_window=7  # Days
    
    # Get recent entries (last week)
    local cutoff_date
    if command -v gdate >/dev/null 2>&1; then
        cutoff_date=$(gdate -u -d "$trend_window days ago" +"%Y-%m-%d" 2>/dev/null || echo "1970-01-01")
    else
        cutoff_date=$(date -u -d "$trend_window days ago" +"%Y-%m-%d" 2>/dev/null || date -u -v-"$trend_window"d +"%Y-%m-%d" 2>/dev/null || echo "1970-01-01")
    fi
    
    # Extract trend data for each metric
    for metric in goal_alignment_rate rework_rate framework_adoption_rate; do
        local values=()
        while IFS= read -r line; do
            if [[ "$line" =~ \"$metric\":[[:space:]]*([0-9]+) ]]; then
                values+=("${BASH_REMATCH[1]}")
            fi
        done < "$HISTORY_FILE"
        
        if [[ ${#values[@]} -ge 2 ]]; then
            local latest=${values[-1]}
            local previous=${values[-2]}
            local trend_direction=""
            local trend_magnitude=$((latest - previous))
            
            if [[ $trend_magnitude -gt 2 ]]; then
                trend_direction="📈 IMPROVING"
            elif [[ $trend_magnitude -lt -2 ]]; then
                trend_direction="📉 DECLINING"
            else
                trend_direction="➡️ STABLE"
            fi
            
            TREND_DATA[$metric]="$trend_direction ($trend_magnitude%)"
            log_verbose "Trend for $metric: ${TREND_DATA[$metric]}"
        else
            TREND_DATA[$metric]="📊 INSUFFICIENT DATA"
        fi
    done
    
    log_success "Trend analysis completed"
}

# Generate dashboard
generate_dashboard() {
    if [[ "$GENERATE_DASHBOARD" != "true" ]]; then
        return 0
    fi
    
    log_info "📊 Generating metrics dashboard..."
    
    local dashboard_file="$METRICS_DIR/decision-metrics-dashboard.html"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    cat > "$dashboard_file" << EOF
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DOC-014 Decision Framework Metrics Dashboard</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metrics-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 20px; }
        .metric-card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metric-value { font-size: 2em; font-weight: bold; margin: 10px 0; }
        .metric-label { color: #666; font-size: 0.9em; text-transform: uppercase; letter-spacing: 1px; }
        .metric-trend { font-size: 0.8em; margin-top: 10px; }
        .excellent { color: #22c55e; }
        .good { color: #84cc16; }
        .acceptable { color: #eab308; }
        .poor { color: #ef4444; }
        .status-excellent { background: #dcfce7; border-left: 4px solid #22c55e; }
        .status-good { background: #f0fdf4; border-left: 4px solid #84cc16; }
        .status-acceptable { background: #fefce8; border-left: 4px solid #eab308; }
        .status-poor { background: #fef2f2; border-left: 4px solid #ef4444; }
        .summary { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .timestamp { color: #666; font-size: 0.8em; text-align: center; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎯 DOC-014 Decision Framework Metrics Dashboard</h1>
            <p>Real-time monitoring of AI Assistant Decision Framework effectiveness and adoption</p>
        </div>
        
        <div class="summary">
            <h2>📊 Current Status Summary</h2>
            <p><strong>Framework Maturity:</strong> ${CURRENT_METRICS[framework_maturity]:-0}%</p>
            <p><strong>Overall Health:</strong> $(get_overall_health_status)</p>
            <p><strong>Key Recommendation:</strong> $(get_primary_recommendation)</p>
        </div>
        
        <div class="metrics-grid">
            <div class="metric-card $(get_metric_status_class ${CURRENT_METRICS[goal_alignment_rate]})">
                <div class="metric-label">Goal Alignment Rate</div>
                <div class="metric-value $(get_metric_color_class ${CURRENT_METRICS[goal_alignment_rate]})">${CURRENT_METRICS[goal_alignment_rate]}%</div>
                <div class="metric-trend">Target: ≥95% | ${TREND_DATA[goal_alignment_rate]:-"No trend data"}</div>
            </div>
            
            <div class="metric-card $(get_metric_status_class ${CURRENT_METRICS[framework_adoption_rate]})">
                <div class="metric-label">Framework Adoption Rate</div>
                <div class="metric-value $(get_metric_color_class ${CURRENT_METRICS[framework_adoption_rate]})">${CURRENT_METRICS[framework_adoption_rate]}%</div>
                <div class="metric-trend">${TREND_DATA[framework_adoption_rate]:-"No trend data"}</div>
            </div>
            
            <div class="metric-card $(get_rework_status_class ${CURRENT_METRICS[rework_rate]})">
                <div class="metric-label">Rework Rate</div>
                <div class="metric-value $(get_rework_color_class ${CURRENT_METRICS[rework_rate]})">${CURRENT_METRICS[rework_rate]}%</div>
                <div class="metric-trend">Target: <5% | ${TREND_DATA[rework_rate]:-"No trend data"}</div>
            </div>
            
            <div class="metric-card $(get_metric_status_class ${CURRENT_METRICS[token_enhancement_rate]})">
                <div class="metric-label">Token Enhancement Rate</div>
                <div class="metric-value $(get_metric_color_class ${CURRENT_METRICS[token_enhancement_rate]})">${CURRENT_METRICS[token_enhancement_rate]}%</div>
                <div class="metric-trend">Enhanced implementation tokens with decision context</div>
            </div>
            
            <div class="metric-card $(get_metric_status_class ${CURRENT_METRICS[validation_pass_rate]})">
                <div class="metric-label">Validation Pass Rate</div>
                <div class="metric-value $(get_metric_color_class ${CURRENT_METRICS[validation_pass_rate]})">${CURRENT_METRICS[validation_pass_rate]}%</div>
                <div class="metric-trend">Current validation framework compliance</div>
            </div>
            
            <div class="metric-card $(get_metric_status_class ${CURRENT_METRICS[decision_compliance_rate]})">
                <div class="metric-label">Decision Compliance Rate</div>
                <div class="metric-value $(get_metric_color_class ${CURRENT_METRICS[decision_compliance_rate]})">${CURRENT_METRICS[decision_compliance_rate]}%</div>
                <div class="metric-trend">Target: ≥90% | Framework requirement adherence</div>
            </div>
        </div>
        
        <div class="timestamp">
            Last updated: $timestamp
        </div>
    </div>
</body>
</html>
EOF

    log_success "Dashboard generated: $dashboard_file"
}

# Helper functions for dashboard
get_overall_health_status() {
    local avg_score=$(((CURRENT_METRICS[goal_alignment_rate] + CURRENT_METRICS[framework_adoption_rate] + CURRENT_METRICS[decision_compliance_rate]) / 3))
    
    if [[ $avg_score -ge 95 ]]; then
        echo "🟢 EXCELLENT"
    elif [[ $avg_score -ge 85 ]]; then
        echo "🟡 GOOD"
    elif [[ $avg_score -ge 70 ]]; then
        echo "🟠 ACCEPTABLE"
    else
        echo "🔴 NEEDS ATTENTION"
    fi
}

get_primary_recommendation() {
    if [[ ${CURRENT_METRICS[goal_alignment_rate]} -lt 95 ]]; then
        echo "Complete protocol integration to improve goal alignment"
    elif [[ ${CURRENT_METRICS[token_enhancement_rate]} -lt 50 ]]; then
        echo "Enhance implementation tokens with decision context"
    elif [[ ${CURRENT_METRICS[rework_rate]} -gt 5 ]]; then
        echo "Focus on decision quality to reduce rework"
    else
        echo "Maintain current excellent performance"
    fi
}

get_metric_status_class() {
    local value=$1
    if [[ $value -ge 95 ]]; then echo "status-excellent"
    elif [[ $value -ge 85 ]]; then echo "status-good"
    elif [[ $value -ge 70 ]]; then echo "status-acceptable"
    else echo "status-poor"
    fi
}

get_metric_color_class() {
    local value=$1
    if [[ $value -ge 95 ]]; then echo "excellent"
    elif [[ $value -ge 85 ]]; then echo "good"
    elif [[ $value -ge 70 ]]; then echo "acceptable"
    else echo "poor"
    fi
}

get_rework_status_class() {
    local value=$1
    if [[ $value -le 2 ]]; then echo "status-excellent"
    elif [[ $value -le 5 ]]; then echo "status-good"
    elif [[ $value -le 10 ]]; then echo "status-acceptable"
    else echo "status-poor"
    fi
}

get_rework_color_class() {
    local value=$1
    if [[ $value -le 2 ]]; then echo "excellent"
    elif [[ $value -le 5 ]]; then echo "good"
    elif [[ $value -le 10 ]]; then echo "acceptable"
    else echo "poor"
    fi
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
        "dashboard")
            log_info "Dashboard format requested - see generated HTML file"
            ;;
        "detailed"|*)
            output_detailed_results
            ;;
    esac
}

# Output detailed results
output_detailed_results() {
    echo ""
    echo "🎯 Decision Quality Metrics Report"
    echo "=================================="
    echo ""
    
    echo "📊 Current Metrics:"
    echo "  Goal Alignment Rate:       ${CURRENT_METRICS[goal_alignment_rate]}%"
    echo "  Framework Adoption Rate:   ${CURRENT_METRICS[framework_adoption_rate]}%"
    echo "  Token Enhancement Rate:    ${CURRENT_METRICS[token_enhancement_rate]}%"
    echo "  Decision Compliance Rate:  ${CURRENT_METRICS[decision_compliance_rate]}%"
    echo "  Validation Pass Rate:      ${CURRENT_METRICS[validation_pass_rate]}%"
    echo "  Rework Rate:               ${CURRENT_METRICS[rework_rate]}%"
    echo ""
    
    if [[ "$TRACK_TRENDS" == "true" && ${#TREND_DATA[@]} -gt 0 ]]; then
        echo "📈 Trends:"
        for metric in "${!TREND_DATA[@]}"; do
            echo "  $metric: ${TREND_DATA[$metric]}"
        done
        echo ""
    fi
    
    echo "🎯 Success Criteria Status:"
    if [[ ${CURRENT_METRICS[goal_alignment_rate]} -ge 95 ]]; then
        echo "  ✅ Goal Alignment: ${CURRENT_METRICS[goal_alignment_rate]}% (≥95%)"
    else
        echo "  ❌ Goal Alignment: ${CURRENT_METRICS[goal_alignment_rate]}% (≥95% required)"
    fi
    
    if [[ ${CURRENT_METRICS[rework_rate]} -le 5 ]]; then
        echo "  ✅ Rework Rate: ${CURRENT_METRICS[rework_rate]}% (≤5%)"
    else
        echo "  ❌ Rework Rate: ${CURRENT_METRICS[rework_rate]}% (≤5% required)"
    fi
}

# Output summary results
output_summary_results() {
    local health_score=$(((CURRENT_METRICS[goal_alignment_rate] + CURRENT_METRICS[framework_adoption_rate] + CURRENT_METRICS[decision_compliance_rate]) / 3))
    
    echo "Decision Framework Health: $health_score%"
    echo "Goal Alignment: ${CURRENT_METRICS[goal_alignment_rate]}%"
    echo "Rework Rate: ${CURRENT_METRICS[rework_rate]}%"
    echo "Framework Adoption: ${CURRENT_METRICS[framework_adoption_rate]}%"
}

# Output JSON results
output_json_results() {
    cat << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "collection_mode": "$COLLECTION_MODE",
  "current_metrics": {
    "goal_alignment_rate": ${CURRENT_METRICS[goal_alignment_rate]},
    "framework_adoption_rate": ${CURRENT_METRICS[framework_adoption_rate]},
    "token_enhancement_rate": ${CURRENT_METRICS[token_enhancement_rate]},
    "decision_compliance_rate": ${CURRENT_METRICS[decision_compliance_rate]},
    "validation_pass_rate": ${CURRENT_METRICS[validation_pass_rate]},
    "rework_rate": ${CURRENT_METRICS[rework_rate]},
    "framework_maturity_score": ${CURRENT_METRICS[framework_maturity]:-0},
    "protocol_coverage": ${CURRENT_METRICS[protocol_coverage]:-0}
  },
  "success_criteria": {
    "goal_alignment_target": 95,
    "goal_alignment_met": $([ ${CURRENT_METRICS[goal_alignment_rate]} -ge 95 ] && echo "true" || echo "false"),
    "rework_rate_target": 5,
    "rework_rate_met": $([ ${CURRENT_METRICS[rework_rate]} -le 5 ] && echo "true" || echo "false")
  }
}
EOF
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --mode)
            COLLECTION_MODE="$2"
            if [[ ! "$COLLECTION_MODE" =~ ^(full|incremental|analysis-only)$ ]]; then
                log_error "Invalid collection mode: $COLLECTION_MODE"
                usage
                exit 1
            fi
            shift 2
            ;;
        --format)
            OUTPUT_FORMAT="$2"
            if [[ ! "$OUTPUT_FORMAT" =~ ^(detailed|summary|json|dashboard)$ ]]; then
                log_error "Invalid output format: $OUTPUT_FORMAT"
                usage
                exit 1
            fi
            shift 2
            ;;
        --time-range)
            TIME_RANGE="$2"
            if [[ ! "$TIME_RANGE" =~ ^(7d|30d|90d|all)$ ]]; then
                log_error "Invalid time range: $TIME_RANGE"
                usage
                exit 1
            fi
            shift 2
            ;;
        --no-trends)
            TRACK_TRENDS=false
            shift
            ;;
        --no-dashboard)
            GENERATE_DASHBOARD=false
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

# Main execution
main() {
    # Only show log messages if not in JSON mode
    if [[ "$OUTPUT_FORMAT" != "json" ]]; then
        log_info "🔶 DOC-014: Decision Quality Metrics Tracker"
        log_info "Mode: $COLLECTION_MODE | Format: $OUTPUT_FORMAT | Range: $TIME_RANGE"
        echo ""
    fi
    
    # Initialize
    init_metrics
    
    # Collect metrics unless analysis-only
    if [[ "$COLLECTION_MODE" != "analysis-only" ]]; then
        if ! collect_current_metrics; then
            [[ "$OUTPUT_FORMAT" != "json" ]] && log_warning "Metrics collection failed, using default values"
        fi
        if ! store_metrics_history; then
            [[ "$OUTPUT_FORMAT" != "json" ]] && log_warning "Failed to store metrics history"
        fi
    fi
    
    # Analyze trends
    if ! analyze_trends; then
        [[ "$OUTPUT_FORMAT" != "json" ]] && log_warning "Trend analysis failed"
    fi
    
    # Generate dashboard
    if ! generate_dashboard; then
        [[ "$OUTPUT_FORMAT" != "json" ]] && log_warning "Dashboard generation failed"
    fi
    
    # Output results
    output_results
    
    [[ "$OUTPUT_FORMAT" != "json" ]] && log_success "Decision quality metrics tracking completed"
}

# Execute main function
main 