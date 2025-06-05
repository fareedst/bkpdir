#!/bin/bash
# Enhanced Implementation Token Migration Script
# 🔺 DOC-014: Enhanced implementation token system - token enhancement migration
# Extends existing token migration with decision context inference

set -euo pipefail

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
FEATURE_TRACKING_FILE="$PROJECT_ROOT/docs/context/feature-tracking.md"
VALIDATION_SCRIPT="$PROJECT_ROOT/scripts/validate-icon-enforcement.sh"

# Default values
DRY_RUN=false
PHASE=""
VERBOSE=false
OUTPUT_FORMAT="detailed"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage information
usage() {
    echo "Enhanced Implementation Token Migration Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --phase <1|2|3>     Migrate specific phase tokens"
    echo "                      Phase 1: ⭐ CRITICAL tokens (89 tokens)"
    echo "                      Phase 2: 🔺 HIGH tokens (207 tokens)"
    echo "                      Phase 3: 🔶🔻 MEDIUM/LOW tokens (296 tokens)"
    echo "  --dry-run           Preview changes without modifying files"
    echo "  --verbose           Enable verbose output"
    echo "  --format <format>   Output format: detailed, summary, json"
    echo "  --help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 --phase 1 --dry-run    Preview Phase 1 token enhancements"
    echo "  $0 --phase 1              Enhance Phase 1 tokens (CRITICAL priority)"
    echo "  $0 --dry-run              Preview all token enhancements"
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
        echo -e "${BLUE}[VERBOSE]${NC} $*"
    fi
}

# Decision context inference function
infer_decision_context() {
    local feature_id="$1"
    local description="$2"
    local priority="$3"
    
    local context_tags=()
    
    # Impact level inference based on feature category
    case "${feature_id%%-*}" in
        "ARCH"|"FILE") 
            context_tags+=("core-functionality") ;;
        "CFG"|"GIT"|"OUT") 
            context_tags+=("enhancement") ;;
        "DOC"|"TEST") 
            context_tags+=("developer-experience") ;;
        "COV"|"REFACTOR"|"LINT") 
            context_tags+=("infrastructure") ;;
        *) 
            # Default based on priority
            case "$priority" in
                "⭐") context_tags+=("core-functionality") ;;
                "🔺") context_tags+=("enhancement") ;;
                "🔶") context_tags+=("developer-experience") ;;
                "🔻") context_tags+=("infrastructure") ;;
            esac ;;
    esac
    
    # Dependency inference
    if [[ "$feature_id" =~ EXTRACT|REFACTOR ]]; then
        context_tags+=("blocks-extraction")
    elif [[ "$description" =~ configuration|template|inheritance ]]; then
        context_tags+=("enables-flexibility")
    elif [[ "$feature_id" =~ TEST.*-001|COV.*-001 ]]; then
        context_tags+=("prerequisite-for-testing")
    else
        context_tags+=("independent")
    fi
    
    # Constraint inference
    if [[ "$description" =~ CLI|command|output|user|naming ]]; then
        context_tags+=("user-facing")
    elif [[ "$description" =~ performance ]]; then
        context_tags+=("performance-critical")
    elif [[ "$description" =~ breaking ]]; then
        context_tags+=("breaking-change")
    else
        context_tags+=("backward-compatible")
    fi
    
    # Return first 3 tags
    local result=$(echo "${context_tags[@]:0:3}" | tr ' ' ',')
    echo "$result"
}

# Get feature priority from feature-tracking.md
get_feature_priority() {
    local feature_id="$1"
    
    if [[ ! -f "$FEATURE_TRACKING_FILE" ]]; then
        log_error "Feature tracking file not found: $FEATURE_TRACKING_FILE"
        return 1
    fi
    
    # Look for the feature ID in the feature tracking table
    local priority=$(grep "^| $feature_id " "$FEATURE_TRACKING_FILE" | head -1 | awk -F'|' '{print $NF}' | tr -d ' ')
    
    if [[ -z "$priority" ]]; then
        log_verbose "Priority not found for $feature_id, using default mapping"
        # Default priority mapping
        case "${feature_id%%-*}" in
            "ARCH"|"FILE") echo "⭐" ;;
            "CFG"|"GIT") echo "🔺" ;;
            "DOC"|"TEST") echo "🔶" ;;
            *) echo "🔻" ;;
        esac
    else
        echo "$priority"
    fi
}

# Extract current tokens from a file
extract_tokens() {
    local file="$1"
    local tokens=()
    
    # Find lines with implementation tokens
    while IFS= read -r line; do
        if [[ "$line" =~ ^[[:space:]]*//[[:space:]]*([⭐🔺🔶🔻])[[:space:]]*([A-Z]+-[0-9]+(-[A-Z0-9]+)*)[[:space:]]*:[[:space:]]*(.+)$ ]]; then
            local priority="${BASH_REMATCH[1]}"
            local feature_id="${BASH_REMATCH[2]}"
            local description="${BASH_REMATCH[4]}"
            
            # Remove existing decision context if present
            description=$(echo "$description" | sed 's/\[DECISION:[^]]*\]$//' | sed 's/[[:space:]]*$//')
            
            tokens+=("$priority|$feature_id|$description")
            log_verbose "Found token: $priority $feature_id: $description"
        fi
    done < "$file"
    
    printf '%s\n' "${tokens[@]}"
}

# Filter tokens by phase
filter_tokens_by_phase() {
    local phase="$1"
    
    case "$phase" in
        "1")
            grep "^⭐" ;;
        "2") 
            grep "^🔺" ;;
        "3")
            grep "^🔶\|^🔻" ;;
        *)
            cat ;;
    esac
}

# Enhance a single token with decision context
enhance_token() {
    local priority="$1"
    local feature_id="$2"
    local description="$3"
    
    # Check if token already has decision context
    if [[ "$description" =~ \[DECISION: ]]; then
        log_verbose "Token $feature_id already has decision context"
        echo "// $priority $feature_id: $description"
        return 0
    fi
    
    # Infer decision context
    local context=$(infer_decision_context "$feature_id" "$description" "$priority")
    
    # Create enhanced token
    local enhanced_token="// $priority $feature_id: $description [DECISION: $context]"
    
    log_verbose "Enhanced token: $enhanced_token"
    echo "$enhanced_token"
}

# Process a single file
process_file() {
    local file="$1"
    local changes_made=false
    
    log_info "Processing file: $file"
    
    # Create temporary file for modifications
    local temp_file=$(mktemp)
    
    # Process file line by line
    while IFS= read -r line; do
        if [[ "$line" =~ ^([[:space:]]*//)([[:space:]]*)([⭐🔺🔶🔻])([[:space:]]*)([A-Z]+-[0-9]+(-[A-Z0-9]+)*)([[:space:]]*:[[:space:]]*)(.+)$ ]]; then
            local prefix="${BASH_REMATCH[1]}${BASH_REMATCH[2]}"
            local priority="${BASH_REMATCH[3]}"
            local feature_id="${BASH_REMATCH[5]}"
            local separator="${BASH_REMATCH[7]}"
            local description="${BASH_REMATCH[8]}"
            
            # Remove existing decision context if present
            description=$(echo "$description" | sed 's/\[DECISION:[^]]*\]$//' | sed 's/[[:space:]]*$//')
            
            # Check phase filter
            if [[ -n "$PHASE" ]]; then
                case "$PHASE" in
                    "1") [[ "$priority" != "⭐" ]] && echo "$line" >> "$temp_file" && continue ;;
                    "2") [[ "$priority" != "🔺" ]] && echo "$line" >> "$temp_file" && continue ;;
                    "3") [[ "$priority" != "🔶" && "$priority" != "🔻" ]] && echo "$line" >> "$temp_file" && continue ;;
                esac
            fi
            
            # Check if should add decision context
            local should_enhance=false
            case "$priority" in
                "⭐") should_enhance=true ;;  # Mandatory for CRITICAL
                "🔺") should_enhance=true ;;  # Recommended for HIGH
                "🔶"|"🔻") 
                    # Optional for MEDIUM/LOW - only if not already present
                    [[ ! "$description" =~ \[DECISION: ]] && should_enhance=true ;;
            esac
            
            if [[ "$should_enhance" == "true" && ! "$description" =~ \[DECISION: ]]; then
                # Enhance token
                local context=$(infer_decision_context "$feature_id" "$description" "$priority")
                local enhanced_line="$prefix$priority $feature_id$separator$description [DECISION: $context]"
                
                echo "$enhanced_line" >> "$temp_file"
                changes_made=true
                
                log_verbose "Enhanced: $feature_id -> [DECISION: $context]"
            else
                # Keep original line
                echo "$line" >> "$temp_file"
            fi
        else
            # Non-token line, keep as-is
            echo "$line" >> "$temp_file"
        fi
    done < "$file"
    
    # Apply changes if not dry run
    if [[ "$changes_made" == "true" ]]; then
        if [[ "$DRY_RUN" == "true" ]]; then
            log_warning "DRY RUN: Would modify $file"
        else
            mv "$temp_file" "$file"
            log_success "Enhanced tokens in $file"
        fi
    else
        rm "$temp_file"
        log_info "No changes needed for $file"
    fi
}

# Find all Go files with implementation tokens
find_token_files() {
    find "$PROJECT_ROOT" -name "*.go" -type f | while read -r file; do
        if grep -q "^[[:space:]]*//[[:space:]]*[⭐🔺🔶🔻]" "$file"; then
            echo "$file"
        fi
    done
}

# Main processing function
process_tokens() {
    local files=($(find_token_files))
    local total_files=${#files[@]}
    local processed_files=0
    
    log_info "Found $total_files files with implementation tokens"
    
    if [[ -n "$PHASE" ]]; then
        log_info "Processing Phase $PHASE tokens only"
    fi
    
    for file in "${files[@]}"; do
        process_file "$file"
        ((processed_files++))
        
        if [[ "$OUTPUT_FORMAT" == "detailed" ]]; then
            echo "Progress: $processed_files/$total_files files processed"
        fi
    done
    
    log_success "Token enhancement complete"
}

# Validation function
validate_enhanced_tokens() {
    log_info "Validating enhanced tokens..."
    
    if [[ -x "$VALIDATION_SCRIPT" ]]; then
        if "$VALIDATION_SCRIPT" > /dev/null 2>&1; then
            log_success "Token validation passed"
        else
            log_warning "Token validation found issues"
            "$VALIDATION_SCRIPT"
        fi
    else
        log_warning "Validation script not found or not executable: $VALIDATION_SCRIPT"
    fi
}

# Statistics function
show_statistics() {
    local files=($(find_token_files))
    local total_tokens=0
    local enhanced_tokens=0
    
    for file in "${files[@]}"; do
        while IFS= read -r line; do
            if [[ "$line" =~ ^[[:space:]]*//[[:space:]]*[⭐🔺🔶🔻] ]]; then
                ((total_tokens++))
                if [[ "$line" =~ \[DECISION: ]]; then
                    ((enhanced_tokens++))
                fi
            fi
        done < "$file"
    done
    
    local enhancement_rate=0
    if [[ $total_tokens -gt 0 ]]; then
        enhancement_rate=$((enhanced_tokens * 100 / total_tokens))
    fi
    
    log_info "Token Enhancement Statistics:"
    echo "  Total tokens: $total_tokens"
    echo "  Enhanced tokens: $enhanced_tokens"
    echo "  Enhancement rate: $enhancement_rate%"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --phase)
            PHASE="$2"
            if [[ ! "$PHASE" =~ ^[123]$ ]]; then
                log_error "Invalid phase: $PHASE. Must be 1, 2, or 3."
                exit 1
            fi
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --format)
            OUTPUT_FORMAT="$2"
            if [[ ! "$OUTPUT_FORMAT" =~ ^(detailed|summary|json)$ ]]; then
                log_error "Invalid format: $OUTPUT_FORMAT. Must be detailed, summary, or json."
                exit 1
            fi
            shift 2
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
    log_info "🔺 DOC-014: Enhanced Implementation Token Migration"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log_warning "DRY RUN MODE: No files will be modified"
    fi
    
    # Show initial statistics
    if [[ "$OUTPUT_FORMAT" == "detailed" ]]; then
        show_statistics
        echo ""
    fi
    
    # Process tokens
    process_tokens
    
    # Show final statistics
    if [[ "$OUTPUT_FORMAT" != "summary" ]]; then
        echo ""
        show_statistics
    fi
    
    # Validate results if not dry run
    if [[ "$DRY_RUN" == "false" ]]; then
        echo ""
        validate_enhanced_tokens
    fi
    
    log_success "Enhanced token migration completed successfully"
}

# Execute main function
main 