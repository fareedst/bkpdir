#!/bin/bash

# 🔺 DOC-009: Mass implementation token standardization - Automated migration system
# This script provides safe, comprehensive migration of legacy implementation tokens
# to the standardized format with priority and action icons

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
DRY_RUN=false
VERBOSE=false
BACKUP_DIR="token-migration-backup"
CHECKPOINT_FILE="token-migration-checkpoint.json"
ROLLBACK_ENABLED=true

# Counters
files_processed=0
tokens_migrated=0
tokens_skipped=0
migration_errors=0

# Feature priority mapping from feature-tracking.md
declare -A FEATURE_PRIORITIES
declare -A FEATURE_ACTIONS

echo -e "${BLUE}🔺 DOC-009: Mass Implementation Token Standardization${NC}"
echo "================================================================="

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --no-backup)
            ROLLBACK_ENABLED=false
            shift
            ;;
        --checkpoint)
            CHECKPOINT_FILE="$2"
            shift 2
            ;;
        --help)
            echo "🔺 DOC-009: Mass implementation token standardization"
            echo ""
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --dry-run     Show what would be changed without modifying files"
            echo "  --verbose     Enable verbose output"
            echo "  --no-backup   Skip creating backup before migration"
            echo "  --checkpoint  Specify checkpoint file (default: token-migration-checkpoint.json)"
            echo "  --help        Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 --dry-run              # Preview changes"
            echo "  $0                        # Execute migration"
            echo "  $0 --verbose              # Detailed output"
            echo "  $0 --rollback             # Rollback previous migration"
            exit 0
            ;;
        --rollback)
            perform_rollback
            exit $?
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# 🔺 DOC-009: Load feature priorities from feature-tracking.md
load_feature_priorities() {
    echo -e "${BLUE}📋 Step 1: Loading feature priorities from feature-tracking.md...${NC}"
    echo "-------------------------------------------------------------------"
    
    if [[ ! -f "docs/context/feature-tracking.md" ]]; then
        echo -e "${RED}❌ Error: feature-tracking.md not found${NC}"
        return 1
    fi
    
    # Parse feature tracking table to extract priorities
    local in_table=false
    local feature_count=0
    
    while IFS= read -r line; do
        # Detect table start
        if [[ "$line" =~ ^\|.*Feature\ ID.*\|.*Priority.*\| ]]; then
            in_table=true
            continue
        fi
        
        # Skip table header separator
        if [[ "$line" =~ ^\|\-+\|\-+ ]]; then
            continue
        fi
        
        # Process table rows
        if [[ $in_table == true ]]; then
            # End of table
            if [[ ! "$line" =~ ^\| ]]; then
                in_table=false
                continue
            fi
            
            # Extract feature ID and priority
            if [[ "$line" =~ ^\|[[:space:]]*([A-Z]+-[0-9]+)[[:space:]]*\| ]]; then
                local feature_id="${BASH_REMATCH[1]}"
                
                # Determine priority based on content
                if echo "$line" | grep -q "🚨 CRITICAL\|CRITICAL"; then
                    FEATURE_PRIORITIES["$feature_id"]="⭐"
                elif echo "$line" | grep -q "🔺 HIGH\|HIGH"; then
                    FEATURE_PRIORITIES["$feature_id"]="🔺"
                elif echo "$line" | grep -q "🔶 MEDIUM\|MEDIUM"; then
                    FEATURE_PRIORITIES["$feature_id"]="🔶"
                elif echo "$line" | grep -q "🔻 LOW\|LOW"; then
                    FEATURE_PRIORITIES["$feature_id"]="🔻"
                else
                    # Default to medium if not specified
                    FEATURE_PRIORITIES["$feature_id"]="🔶"
                fi
                
                # Determine action icon based on feature type
                if echo "$line" | grep -q "validation\|test\|verify"; then
                    FEATURE_ACTIONS["$feature_id"]="🛡️"
                elif echo "$line" | grep -q "config\|setup\|extract"; then
                    FEATURE_ACTIONS["$feature_id"]="🔧"
                elif echo "$line" | grep -q "document\|format\|output"; then
                    FEATURE_ACTIONS["$feature_id"]="📝"
                else
                    FEATURE_ACTIONS["$feature_id"]="🔍"
                fi
                
                ((feature_count++))
                
                if [[ "$VERBOSE" == "true" ]]; then
                    echo -e "    ${CYAN}$feature_id → ${FEATURE_PRIORITIES[$feature_id]} ${FEATURE_ACTIONS[$feature_id]}${NC}"
                fi
            fi
        fi
    done < "docs/context/feature-tracking.md"
    
    echo -e "  ${GREEN}✅ Loaded $feature_count feature priorities${NC}"
    
    # Add some common patterns for features not in tracking
    FEATURE_PRIORITIES["REFACTOR-001"]="🔶"
    FEATURE_PRIORITIES["REFACTOR-002"]="🔶"
    FEATURE_PRIORITIES["REFACTOR-003"]="🔶"
    FEATURE_PRIORITIES["REFACTOR-004"]="🔶"
    FEATURE_PRIORITIES["REFACTOR-005"]="🔶"
    FEATURE_PRIORITIES["REFACTOR-006"]="🔶"
    
    FEATURE_ACTIONS["REFACTOR-001"]="🔧"
    FEATURE_ACTIONS["REFACTOR-002"]="🔧"
    FEATURE_ACTIONS["REFACTOR-003"]="🔧"
    FEATURE_ACTIONS["REFACTOR-004"]="🔧"
    FEATURE_ACTIONS["REFACTOR-005"]="🔧"
    FEATURE_ACTIONS["REFACTOR-006"]="🛡️"
    
    return 0
}

# 🔺 DOC-009: Create backup before migration
create_migration_backup() {
    if [[ "$ROLLBACK_ENABLED" == "false" ]]; then
        echo -e "  ${YELLOW}⚠️ Backup disabled, skipping...${NC}"
        return 0
    fi
    
    echo -e "${BLUE}📋 Step 2: Creating migration backup...${NC}"
    echo "-------------------------------------------------------------------"
    
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    BACKUP_DIR="token-migration-backup-$timestamp"
    
    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "  ${CYAN}[DRY RUN] Would create backup in: $BACKUP_DIR${NC}"
        return 0
    fi
    
    mkdir -p "$BACKUP_DIR"
    
    # Backup all Go files
    local go_files=($(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" 2>/dev/null))
    local backup_count=0
    
    for go_file in "${go_files[@]}"; do
        local backup_path="$BACKUP_DIR/$go_file"
        mkdir -p "$(dirname "$backup_path")"
        cp "$go_file" "$backup_path"
        ((backup_count++))
    done
    
    # Create checkpoint metadata
    cat > "$BACKUP_DIR/migration-metadata.json" << EOF
{
    "timestamp": "$timestamp",
    "files_backed_up": $backup_count,
    "backup_directory": "$BACKUP_DIR",
    "rollback_command": "./scripts/token-migration.sh --rollback --checkpoint $BACKUP_DIR/migration-metadata.json"
}
EOF

    echo -e "  ${GREEN}✅ Backup created: $BACKUP_DIR ($backup_count files)${NC}"
    echo -e "  ${GREEN}✅ Rollback available: ./scripts/token-migration.sh --rollback${NC}"
    
    return 0
}

# 🔺 DOC-009: Analyze and migrate tokens in a file
migrate_file_tokens() {
    local file_path="$1"
    local file_tokens=0
    local file_migrated=0
    local file_errors=0
    
    if [[ ! -f "$file_path" ]]; then
        return 1
    fi
    
    # Create temporary file for modifications
    local temp_file=$(mktemp)
    cp "$file_path" "$temp_file"
    
    # Process each line
    local line_num=0
    while IFS= read -r line; do
        ((line_num++))
        
        # Check for legacy implementation token pattern
        if echo "$line" | grep -qE '// .*[A-Z]+-[0-9]+:' && ! echo "$line" | grep -qE '// *[⭐🔺🔶🔻]'; then
            ((file_tokens++))
            
            # Extract feature ID
            local feature_id=$(echo "$line" | grep -oE '[A-Z]+-[0-9]+' | head -1)
            
            if [[ -n "$feature_id" ]]; then
                # Get priority and action icons
                local priority_icon="${FEATURE_PRIORITIES[$feature_id]:-🔶}"
                local action_icon="${FEATURE_ACTIONS[$feature_id]:-🔍}"
                
                # Extract the comment content after the colon
                local comment_content=$(echo "$line" | sed -E 's|^(.*)// .*[A-Z]+-[0-9]+: *(.*)|\2|')
                
                # Build new standardized token
                local new_token="// $priority_icon $feature_id: $comment_content - $action_icon"
                
                # Replace the line
                local escaped_line=$(echo "$line" | sed 's/[[\.*^$()+?{|]/\\&/g')
                sed -i.bak "${line_num}s|.*|$new_token|" "$temp_file"
                
                if [[ "$VERBOSE" == "true" ]]; then
                    echo -e "    ${CYAN}Line $line_num: $feature_id → $priority_icon $action_icon${NC}"
                fi
                
                ((file_migrated++))
            else
                echo -e "    ${RED}❌ Could not extract feature ID from line $line_num${NC}"
                ((file_errors++))
            fi
        fi
    done < "$file_path"
    
    # Apply changes if not dry run
    if [[ "$DRY_RUN" == "false" && $file_migrated -gt 0 ]]; then
        mv "$temp_file" "$file_path"
        echo -e "  ${GREEN}✅ Updated $file_path: $file_migrated tokens migrated${NC}"
    else
        rm -f "$temp_file"
        if [[ "$DRY_RUN" == "true" && $file_migrated -gt 0 ]]; then
            echo -e "  ${CYAN}[DRY RUN] Would update $file_path: $file_migrated tokens${NC}"
        fi
    fi
    
    # Clean up backup file created by sed
    rm -f "${temp_file}.bak"
    
    tokens_migrated=$((tokens_migrated + file_migrated))
    migration_errors=$((migration_errors + file_errors))
    
    return 0
}

# 🔺 DOC-009: Process all Go files for token migration
process_token_migration() {
    echo -e "${BLUE}📋 Step 3: Processing implementation token migration...${NC}"
    echo "-------------------------------------------------------------------"
    
    local go_files=($(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" 2>/dev/null))
    
    echo -e "  ${CYAN}Found ${#go_files[@]} Go files to process${NC}"
    
    for go_file in "${go_files[@]}"; do
        if [[ "$VERBOSE" == "true" ]] || [[ "$DRY_RUN" == "true" ]]; then
            echo -e "${CYAN}  Processing $go_file...${NC}"
        fi
        
        migrate_file_tokens "$go_file"
        ((files_processed++))
    done
    
    echo
    echo -e "${BLUE}📊 Migration Summary:${NC}"
    echo "  Files processed: $files_processed"
    echo "  Tokens migrated: $tokens_migrated"
    echo "  Errors encountered: $migration_errors"
    
    if [[ $tokens_migrated -gt 0 ]]; then
        echo -e "  ${GREEN}✅ Successfully migrated $tokens_migrated implementation tokens${NC}"
    fi
    
    if [[ $migration_errors -gt 0 ]]; then
        echo -e "  ${YELLOW}⚠️ $migration_errors tokens could not be migrated${NC}"
    fi
}

# 🔺 DOC-009: Validation after migration
validate_migration() {
    echo -e "${BLUE}📋 Step 4: Validating migration results...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Run the icon validation to check results
    echo -e "  ${CYAN}Running comprehensive icon validation...${NC}"
    
    if [[ -x "scripts/validate-icon-enforcement.sh" ]]; then
        if ./scripts/validate-icon-enforcement.sh > /tmp/validation-output.log 2>&1; then
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
    else
        echo -e "  ${YELLOW}⚠️ Icon validation script not found, skipping validation${NC}"
        return 0
    fi
}

# 🔺 DOC-009: Rollback capability
perform_rollback() {
    echo -e "${BLUE}📋 DOC-009: Rolling back token migration...${NC}"
    echo "-------------------------------------------------------------------"
    
    # Find most recent backup
    local backup_dirs=($(ls -dt token-migration-backup-* 2>/dev/null || true))
    
    if [[ ${#backup_dirs[@]} -eq 0 ]]; then
        echo -e "${RED}❌ No migration backups found${NC}"
        exit 1
    fi
    
    local latest_backup="${backup_dirs[0]}"
    
    if [[ ! -f "$latest_backup/migration-metadata.json" ]]; then
        echo -e "${RED}❌ Backup metadata not found in $latest_backup${NC}"
        exit 1
    fi
    
    echo -e "  ${CYAN}Rolling back from: $latest_backup${NC}"
    
    # Restore files
    local restored_count=0
    while IFS= read -r backup_file; do
        local original_file=$(echo "$backup_file" | sed "s|^$latest_backup/||")
        
        if [[ -f "$backup_file" ]]; then
            cp "$backup_file" "$original_file"
            ((restored_count++))
        fi
    done < <(find "$latest_backup" -name "*.go" -type f)
    
    echo -e "  ${GREEN}✅ Restored $restored_count files${NC}"
    echo -e "  ${GREEN}✅ Rollback completed successfully${NC}"
    
    return 0
}

# 🔺 DOC-009: Main migration orchestrator
main() {
    local start_time=$(date +%s)
    
    echo -e "${BLUE}🎯 Migration Configuration:${NC}"
    echo "  Dry run: $([ "$DRY_RUN" == "true" ] && echo "Yes" || echo "No")"
    echo "  Verbose: $([ "$VERBOSE" == "true" ] && echo "Yes" || echo "No")"
    echo "  Backup: $([ "$ROLLBACK_ENABLED" == "true" ] && echo "Yes" || echo "No")"
    echo
    
    # Check prerequisites
    if [[ ! -f "docs/context/feature-tracking.md" ]]; then
        echo -e "${RED}❌ Error: Must be run from project root directory${NC}"
        echo -e "${RED}   Required: docs/context/feature-tracking.md${NC}"
        exit 1
    fi
    
    # Execute migration steps
    load_feature_priorities || { echo -e "${RED}❌ Failed to load feature priorities${NC}"; exit 1; }
    create_migration_backup || { echo -e "${RED}❌ Failed to create backup${NC}"; exit 1; }
    process_token_migration || { echo -e "${RED}❌ Migration processing failed${NC}"; exit 1; }
    
    # Validate results if not dry run
    if [[ "$DRY_RUN" == "false" ]]; then
        validate_migration || { echo -e "${YELLOW}⚠️ Migration validation warnings detected${NC}"; }
    fi
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo
    echo -e "${BLUE}🎯 Migration Complete${NC}"
    echo "=============================================================="
    echo -e "  Duration: ${duration}s"
    echo -e "  Files processed: $files_processed"
    echo -e "  Tokens migrated: ${GREEN}$tokens_migrated${NC}"
    
    if [[ $migration_errors -gt 0 ]]; then
        echo -e "  Errors: ${RED}$migration_errors${NC}"
    fi
    
    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${CYAN}📋 This was a dry run. To execute migration, run:${NC}"
        echo -e "${CYAN}   ./scripts/token-migration.sh${NC}"
    else
        echo -e "${GREEN}✅ Token migration completed successfully!${NC}"
        
        if [[ "$ROLLBACK_ENABLED" == "true" ]]; then
            echo -e "${CYAN}📋 To rollback if needed, run:${NC}"
            echo -e "${CYAN}   ./scripts/token-migration.sh --rollback${NC}"
        fi
    fi
    
    exit 0
}

# Run main function
main "$@" 