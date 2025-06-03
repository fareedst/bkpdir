#!/bin/bash

# 🔺 DOC-009: Mass implementation token standardization - 🔧 Token migration automation
# This script provides automated migration of legacy implementation tokens to standardized format
# with priority icons, action icons, and validation capabilities

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
BACKUP_DIR="./token-migration-backups"
DRY_RUN=${1:-true}  # Default to dry run for safety
BATCH_SIZE=${BATCH_SIZE:-50}  # Number of files to process in each batch
ROLLBACK_FILE=""
CHECKPOINT_FILE="./token-migration-checkpoint.json"

echo -e "${BLUE}🔺 DOC-009: Mass Implementation Token Standardization${NC}"
echo "=============================================================="

# Check if we're in the right directory
if [[ ! -f "docs/context/feature-tracking.md" ]] || [[ ! -f "Makefile" ]]; then
    echo -e "${RED}❌ Error: Must be run from project root directory${NC}"
    echo -e "${RED}   Required files: docs/context/feature-tracking.md, Makefile${NC}"
    exit 1
fi

# 🔺 DOC-009: Feature priority mapping from feature-tracking.md - 🔍 Priority analysis
get_feature_priority() {
    local feature_id="$1"
    # Extract the primary prefix (first part before any hyphen)
    local feature_prefix=$(echo "$feature_id" | cut -d'-' -f1)
    
    case "$feature_prefix" in
        "ARCH")  echo "⭐" ;;  # Archive operations - CRITICAL
        "FILE")  echo "⭐" ;;  # File operations - CRITICAL  
        "CFG")   echo "🔺" ;;  # Configuration - HIGH
        "GIT")   echo "🔺" ;;  # Git integration - HIGH
        "TEST")  echo "🔺" ;;  # Testing infrastructure - HIGH
        "OUT")   echo "🔶" ;;  # Output management - MEDIUM
        "DOC")   echo "🔺" ;;  # Documentation - HIGH
        "LINT")  echo "🔺" ;;  # Code quality - HIGH
        "COV")   echo "🔺" ;;  # Coverage - HIGH
        "REFACTOR") echo "🔶" ;;  # Refactoring - MEDIUM
        *)       echo "🔶" ;;  # Default to MEDIUM
    esac
}

# 🔺 DOC-009: Action icon inference from function context - 🔍 Function behavior analysis
get_action_icon() {
    local context="$1"
    local line_content="$2"
    
    # Analyze surrounding context and function names for action type
    if echo "$context" | grep -qi "get\|find\|search\|discover\|detect\|analyze\|check\|validate\|parse"; then
        echo "🔍"  # SEARCH/DISCOVER
    elif echo "$context" | grep -qi "format\|print\|write\|update\|log\|output\|render\|display"; then
        echo "📝"  # DOCUMENT/UPDATE
    elif echo "$context" | grep -qi "config\|setup\|init\|create\|build\|generate\|establish"; then
        echo "🔧"  # CONFIGURE/MODIFY
    elif echo "$context" | grep -qi "protect\|secure\|validate\|verify\|guard\|ensure\|handle.*error"; then
        echo "🛡️"  # PROTECT/VALIDATE
    else
        # Default based on common patterns
        if echo "$line_content" | grep -qi "func.*get\|func.*find\|func.*search"; then
            echo "🔍"
        elif echo "$line_content" | grep -qi "func.*format\|func.*print\|func.*write"; then
            echo "📝"
        elif echo "$line_content" | grep -qi "func.*config\|func.*setup\|func.*create"; then
            echo "🔧"
        else
            echo "🔧"  # Default to CONFIGURE/MODIFY
        fi
    fi
}

# 🔺 DOC-009: Create backup before migration - 🛡️ Safe backup operations
create_backup() {
    echo -e "${BLUE}📋 Step 1: Creating migration backups...${NC}"
    echo "-------------------------------------------------------------------"
    
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    BACKUP_DIR="./token-migration-backups/$timestamp"
    
    mkdir -p "$BACKUP_DIR"
    
    # Find all Go files with implementation tokens
    local files_with_tokens=($(grep -l "// [A-Z].*-[0-9]\+.*:" . --include="*.go" -r 2>/dev/null || true))
    
    echo -e "  ${CYAN}Backing up ${#files_with_tokens[@]} files with implementation tokens...${NC}"
    
    for file in "${files_with_tokens[@]}"; do
        local backup_path="$BACKUP_DIR$(dirname "$file")"
        mkdir -p "$backup_path"
        cp "$file" "$backup_path/$(basename "$file")"
        echo -e "    ${GREEN}✅ Backed up: $file${NC}"
    done
    
    # Store backup location for rollback
    echo "$BACKUP_DIR" > .token-migration-backup-location
    
    echo -e "  ${GREEN}✅ Backup completed: $BACKUP_DIR${NC}"
    echo -e "  ${YELLOW}💾 Backup location saved for rollback capability${NC}"
}

# 🔺 DOC-009: Checkpoint system for incremental progress - 📊 Progress validation
save_checkpoint() {
    local processed_files="$1"
    local total_files="$2"
    local current_batch="$3"
    
    cat > "$CHECKPOINT_FILE" << EOF
{
    "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
    "processed_files": $processed_files,
    "total_files": $total_files,
    "current_batch": $current_batch,
    "backup_dir": "$BACKUP_DIR",
    "batch_size": $BATCH_SIZE
}
EOF
    echo -e "    ${CYAN}💾 Checkpoint saved: $processed_files/$total_files files processed${NC}"
}

# 🔺 DOC-009: Standardize a single implementation token - 🔧 Token format standardization
standardize_token() {
    local file="$1"
    local line_num="$2"
    local original_line="$3"
    local dry_run="$4"
    
    # Extract feature ID from original token
    local feature_id=$(echo "$original_line" | grep -o '[A-Z][A-Z0-9_-]*-[0-9]\+[A-Z0-9_-]*' | head -1)
    if [[ -z "$feature_id" ]]; then
        echo -e "      ${RED}❌ Could not extract feature ID from: $original_line${NC}"
        return 1
    fi
    
    # Extract description (everything after the colon, removing the feature ID prefix)
    local description=$(echo "$original_line" | sed "s/.*\/\/ *$feature_id: *//" | sed 's/^[⭐🔺🔶🔻🔍📝🔧🛡️ ]*//')
    
    # Get priority icon from feature mapping
    local priority_icon=$(get_feature_priority "$feature_id")
    
    # Get function context for action icon inference
    local context=$(sed -n "$((line_num-5)),$((line_num+5))p" "$file" 2>/dev/null || echo "")
    local action_icon=$(get_action_icon "$context" "$original_line")
    
    # Create standardized token
    local standardized_token="// $priority_icon $feature_id: $description - $action_icon"
    
    if [[ "$dry_run" == "true" ]]; then
        echo -e "      ${CYAN}🔍 DRY RUN - Would update:${NC}"
        echo -e "        ${YELLOW}OLD: $original_line${NC}"
        echo -e "        ${GREEN}NEW: $standardized_token${NC}"
        return 0
    else
        # Actually update the file
        sed -i.bak "${line_num}s|.*|$standardized_token|" "$file"
        echo -e "      ${GREEN}✅ Updated: $feature_id with priority $priority_icon and action $action_icon${NC}"
        return 0
    fi
}

# 🔺 DOC-009: Process a single file for token standardization - 🔧 File processing
process_file() {
    local file="$1"
    local dry_run="$2"
    
    echo -e "    ${CYAN}Processing: $file${NC}"
    
    local tokens_updated=0
    local tokens_failed=0
    
    # Find all implementation tokens in the file
    while IFS= read -r line; do
        local line_num=$(echo "$line" | cut -d: -f1)
        local line_content=$(echo "$line" | cut -d: -f2-)
        
        # Skip already standardized tokens (those with priority icons)
        if echo "$line_content" | grep -q "// *[⭐🔺🔶🔻]"; then
            echo -e "      ${GREEN}✅ Already standardized: $line_content${NC}"
            continue
        fi
        
        # Standardize the token
        if standardize_token "$file" "$line_num" "$line_content" "$dry_run"; then
            ((tokens_updated++))
        else
            ((tokens_failed++))
        fi
        
    done < <(grep -n "// [A-Z].*-[0-9]\+.*:" "$file" 2>/dev/null || true)
    
    echo -e "      ${GREEN}✅ File completed: $tokens_updated updated, $tokens_failed failed${NC}"
    return $tokens_failed
}

# 🔺 DOC-009: Batch processing with progress tracking - 📊 Batch operations
process_batch() {
    local files=("$@")
    local dry_run="$DRY_RUN"
    
    echo -e "${BLUE}📋 Step 2: Processing batch of ${#files[@]} files...${NC}"
    echo "-------------------------------------------------------------------"
    
    local batch_success=0
    local batch_failed=0
    
    for file in "${files[@]}"; do
        if [[ -f "$file" ]]; then
            if process_file "$file" "$dry_run"; then
                ((batch_success++))
            else
                ((batch_failed++))
            fi
        else
            echo -e "    ${YELLOW}⚠️ File not found: $file${NC}"
            ((batch_failed++))
        fi
    done
    
    echo -e "  ${GREEN}✅ Batch completed: $batch_success successful, $batch_failed failed${NC}"
    return $batch_failed
}

# 🔺 DOC-009: Validation after migration - 🛡️ Migration validation
validate_migration() {
    echo -e "${BLUE}📋 Step 3: Validating migration results...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Run the icon validation to check results
    echo -e "  ${CYAN}Running comprehensive icon validation...${NC}"
    
    if make validate-icon-enforcement > /tmp/validation-output.log 2>&1; then
        echo -e "  ${GREEN}✅ Icon validation passed!${NC}"
        
        # Extract standardization rate from validation output
        local std_rate=$(grep "Standardization rate:" /tmp/validation-output.log | tail -1 | grep -o '[0-9]\+%' || echo "Unknown")
        echo -e "  ${GREEN}✅ Current standardization rate: $std_rate${NC}"
        
        return 0
    else
        echo -e "  ${RED}❌ Icon validation failed${NC}"
        echo -e "  ${YELLOW}📋 Validation details:${NC}"
        tail -20 /tmp/validation-output.log | sed 's/^/    /'
        return 1
    fi
}

# 🔺 DOC-009: Rollback capability - 🛡️ Safe rollback operations
rollback_migration() {
    echo -e "${BLUE}🔄 DOC-009: Rolling back migration...${NC}"
    echo "=============================================================="
    
    if [[ ! -f ".token-migration-backup-location" ]]; then
        echo -e "${RED}❌ Error: No backup location found${NC}"
        echo -e "${RED}   Cannot perform rollback without backup location${NC}"
        exit 1
    fi
    
    local backup_location=$(cat .token-migration-backup-location)
    
    if [[ ! -d "$backup_location" ]]; then
        echo -e "${RED}❌ Error: Backup directory not found: $backup_location${NC}"
        exit 1
    fi
    
    echo -e "${CYAN}Restoring files from backup: $backup_location${NC}"
    
    # Restore all files
    find "$backup_location" -name "*.go" -type f | while read -r backup_file; do
        local relative_path=${backup_file#$backup_location/}
        local original_file="./$relative_path"
        
        if cp "$backup_file" "$original_file"; then
            echo -e "  ${GREEN}✅ Restored: $original_file${NC}"
        else
            echo -e "  ${RED}❌ Failed to restore: $original_file${NC}"
        fi
    done
    
    # Clean up rollback tracking
    rm -f .token-migration-backup-location
    rm -f "$CHECKPOINT_FILE"
    
    echo -e "${GREEN}✅ Rollback completed${NC}"
}

# 🔺 DOC-009: Main migration orchestration - 🔧 Migration coordination
main_migration() {
    local dry_run="$DRY_RUN"
    
    if [[ "$dry_run" == "true" ]]; then
        echo -e "${YELLOW}🔍 DRY RUN MODE - No actual changes will be made${NC}"
        echo -e "${YELLOW}   Run with 'false' as first argument to perform actual migration${NC}"
        echo
    fi
    
    # Create backup (only in actual run mode)
    if [[ "$dry_run" == "false" ]]; then
        create_backup
    else
        echo -e "${YELLOW}📋 Step 1: Backup skipped in dry run mode${NC}"
    fi
    
    # Find all Go files with implementation tokens  
    echo -e "${BLUE}📋 Discovering files with implementation tokens...${NC}"
    local files_with_tokens=($(grep -l "// [A-Z].*-[0-9]\+.*:" . --include="*.go" -r 2>/dev/null || true))
    
    echo -e "  ${GREEN}✅ Found ${#files_with_tokens[@]} files with implementation tokens${NC}"
    
    # Process files in batches
    local total_files=${#files_with_tokens[@]}
    local processed_files=0
    local current_batch=1
    local total_failed=0
    
    for ((i=0; i<$total_files; i+=$BATCH_SIZE)); do
        echo
        echo -e "${MAGENTA}📦 Processing Batch $current_batch...${NC}"
        
        # Create batch array
        local batch=("${files_with_tokens[@]:$i:$BATCH_SIZE}")
        
        # Process the batch
        if process_batch "${batch[@]}"; then
            processed_files=$((processed_files + ${#batch[@]}))
            
            # Save checkpoint (only in actual run mode)
            if [[ "$dry_run" == "false" ]]; then
                save_checkpoint $processed_files $total_files $current_batch
            fi
        else
            echo -e "  ${RED}❌ Batch $current_batch had failures${NC}"
            ((total_failed++))
        fi
        
        ((current_batch++))
        
        # Progress update
        local progress=$((processed_files * 100 / total_files))
        echo -e "  ${CYAN}📊 Progress: $processed_files/$total_files files ($progress%)${NC}"
    done
    
    echo
    echo -e "${BLUE}📋 Migration Summary:${NC}"
    echo "  Total files processed: $processed_files/$total_files"
    echo "  Batches with failures: $total_failed"
    
    # Validation (only in actual run mode)
    if [[ "$dry_run" == "false" ]]; then
        if validate_migration; then
            echo -e "  ${GREEN}✅ Migration validation passed${NC}"
            
            # Clean up checkpoint file on success
            rm -f "$CHECKPOINT_FILE"
            
            echo
            echo -e "${GREEN}🎉 DOC-009: Mass token standardization completed successfully!${NC}"
            echo -e "${GREEN}   Run 'make validate-icon-enforcement' to see updated metrics${NC}"
        else
            echo -e "  ${RED}❌ Migration validation failed${NC}"
            echo -e "  ${YELLOW}💡 Consider running rollback: $0 rollback${NC}"
            exit 1
        fi
    else
        echo -e "  ${YELLOW}🔍 Dry run completed - no validation performed${NC}"
        echo -e "  ${GREEN}💡 Run '$0 false' to perform actual migration${NC}"
    fi
}

# 🔺 DOC-009: Command line interface - 📝 CLI processing
case "${1:-}" in
    "rollback")
        rollback_migration
        ;;
    "help"|"-h"|"--help")
        echo "🔺 DOC-009: Mass Implementation Token Standardization"
        echo
        echo "Usage:"
        echo "  $0 [true|false]     - Run migration (true=dry run, false=actual)"
        echo "  $0 rollback         - Rollback last migration"
        echo "  $0 help             - Show this help"
        echo
        echo "Environment Variables:"
        echo "  BATCH_SIZE=50       - Number of files per batch"
        echo
        echo "Examples:"
        echo "  $0                  - Dry run (default)"
        echo "  $0 true             - Dry run (explicit)"
        echo "  $0 false            - Actual migration"
        echo "  BATCH_SIZE=25 $0 false - Smaller batches"
        ;;
    *)
        main_migration
        ;;
esac 