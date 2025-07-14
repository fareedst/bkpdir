#!/bin/bash

# [CRITICAL] TOKEN-004: Simple semantic token migration script - [ACTION:migration]
# Purpose: Migrate from Unicode icons to semantic tokens using simple sed commands

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${PROJECT_ROOT}/backup-semantic-tokens-simple-$(date +%Y%m%d-%H%M%S)"
LOG_FILE="${PROJECT_ROOT}/semantic-token-migration-simple.log"

echo -e "${BLUE}🚀 Starting Simple Semantic Token Migration${NC}"
echo "Project Root: ${PROJECT_ROOT}"
echo "Backup Directory: ${BACKUP_DIR}"
echo "Log File: ${LOG_FILE}"
echo ""

# Create backup directory
mkdir -p "${BACKUP_DIR}"

# Function to log messages
log_message() {
    local level="$1"
    local message="$2"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo "[${timestamp}] [${level}] ${message}" | tee -a "${LOG_FILE}"
}

# Function to create backup
create_backup() {
    local file="$1"
    local backup_path="${BACKUP_DIR}/$(basename "${file}")"
    cp "${file}" "${backup_path}"
    log_message "BACKUP" "Created backup: ${backup_path}"
}

# Function to migrate a single file
migrate_file() {
    local file="$1"
    local temp_file="${file}.tmp"
    
    log_message "INFO" "Processing file: ${file}"
    
    # Check if file contains Unicode icons
    if ! grep -q "[⭐🔺🔶🔻🔧📝🔍🛠️🛡️🔄]" "${file}" 2>/dev/null; then
        log_message "SKIP" "No Unicode icons found in ${file}"
        echo -e "${YELLOW}⏭️  Skipping ${file} (no Unicode icons)${NC}"
        return 0
    fi
    
    # Create backup
    create_backup "${file}"
    
    # Copy original to temp file
    cp "${file}" "${temp_file}"
    
    # Apply semantic token replacements using perl for better Unicode handling
    local changes_made=0
    
    # Replace priority icons
    local count=$(grep -o "⭐" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/⭐/[CRITICAL]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of ⭐ with [CRITICAL] in ${file}"
    fi
    
    count=$(grep -o "🔺" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔺/[HIGH]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔺 with [HIGH] in ${file}"
    fi
    
    count=$(grep -o "🔶" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔶/[MEDIUM]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔶 with [MEDIUM] in ${file}"
    fi
    
    count=$(grep -o "🔻" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔻/[LOW]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔻 with [LOW] in ${file}"
    fi
    
    # Replace action icons
    count=$(grep -o "🔧" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔧/[ACTION:core-functionality]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔧 with [ACTION:core-functionality] in ${file}"
    fi
    
    count=$(grep -o "📝" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/📝/[ACTION:format-processing]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 📝 with [ACTION:format-processing] in ${file}"
    fi
    
    count=$(grep -o "🔍" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔍/[ACTION:discovery]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔍 with [ACTION:discovery] in ${file}"
    fi
    
    count=$(grep -o "🛠️" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🛠️/[ACTION:maintenance]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🛠️ with [ACTION:maintenance] in ${file}"
    fi
    
    count=$(grep -o "🛡️" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🛡️/[ACTION:validation]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🛡️ with [ACTION:validation] in ${file}"
    fi
    
    count=$(grep -o "🔄" "${temp_file}" 2>/dev/null | wc -l || echo "0")
    if [ "$count" -gt 0 ]; then
        perl -pi -e 's/🔄/[ACTION:migration]/g' "${temp_file}"
        changes_made=$((changes_made + count))
        log_message "REPLACE" "Replaced ${count} instances of 🔄 with [ACTION:migration] in ${file}"
    fi
    
    if [ "$changes_made" -gt 0 ]; then
        # Replace original with migrated file
        mv "${temp_file}" "${file}"
        log_message "SUCCESS" "Migrated ${file} with ${changes_made} changes"
        echo -e "${GREEN}✅ Migrated ${file} (${changes_made} changes)${NC}"
    else
        # No changes made, remove temp file
        rm -f "${temp_file}"
        log_message "INFO" "No changes needed for ${file}"
        echo -e "${YELLOW}⚠️  No changes needed for ${file}${NC}"
    fi
}

# Function to validate migration
validate_migration() {
    local file="$1"
    local validation_errors=0
    
    log_message "VALIDATE" "Validating migration for ${file}"
    
    # Check for remaining Unicode icons
    local remaining_icons=$(grep -o "[⭐🔺🔶🔻🔧📝🔍🛠️🛡️🔄]" "${file}" 2>/dev/null | wc -l || echo "0")
    if [ "$remaining_icons" -gt 0 ]; then
        log_message "ERROR" "Found ${remaining_icons} remaining Unicode icons in ${file}"
        echo -e "${RED}❌ Found ${remaining_icons} remaining Unicode icons in ${file}${NC}"
        validation_errors=$((validation_errors + 1))
    fi
    
    # Check for semantic tokens
    local semantic_tokens=$(grep -o "\[CRITICAL\]\|\[HIGH\]\|\[MEDIUM\]\|\[LOW\]\|\[ACTION:" "${file}" 2>/dev/null | wc -l || echo "0")
    if [ "$semantic_tokens" -gt 0 ]; then
        log_message "SUCCESS" "Found ${semantic_tokens} semantic tokens in ${file}"
        echo -e "${GREEN}✅ Found ${semantic_tokens} semantic tokens in ${file}${NC}"
    fi
    
    return $validation_errors
}

# Function to generate migration report
generate_report() {
    local report_file="${PROJECT_ROOT}/semantic-token-migration-simple-report.md"
    
    cat > "${report_file}" << EOF
# Simple Semantic Token Migration Report

**Migration Date**: $(date '+%Y-%m-%d %H:%M:%S')
**Project**: BkpDir
**Migration Script**: migrate-to-semantic-tokens-simple.sh

## Migration Summary

### Token Mapping Applied

#### Priority Tokens
- ⭐ → [CRITICAL]
- 🔺 → [HIGH]
- 🔶 → [MEDIUM]
- 🔻 → [LOW]

#### Action Tokens
- 🔧 → [ACTION:core-functionality]
- 📝 → [ACTION:format-processing]
- 🔍 → [ACTION:discovery]
- 🛠️ → [ACTION:maintenance]
- 🛡️ → [ACTION:validation]
- 🔄 → [ACTION:migration]

### Files Processed

\`\`\`
$(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "${BACKUP_DIR}" | sort)
\`\`\`

### Migration Log

\`\`\`
$(cat "${LOG_FILE}" 2>/dev/null || echo "No log file found")
\`\`\`

## Validation Results

### Files with Semantic Tokens
\`\`\`
$(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "${BACKUP_DIR}" | xargs grep -l "\[CRITICAL\]\|\[HIGH\]\|\[MEDIUM\]\|\[LOW\]\|\[ACTION:" 2>/dev/null || echo "No files with semantic tokens found")
\`\`\`

### Files with Remaining Unicode Icons
\`\`\`
$(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "${BACKUP_DIR}" | xargs grep -l "[⭐🔺🔶🔻🔧📝🔍🛠️🛡️🔄]" 2>/dev/null || echo "No files with remaining Unicode icons found")
\`\`\`

## Next Steps

1. **Review Migration**: Check all migrated files for accuracy
2. **Update Documentation**: Update any documentation referencing old tokens
3. **Update Tests**: Ensure all tests pass with new semantic tokens
4. **Update Validation**: Update any validation scripts to use semantic tokens
5. **AI Assistant Integration**: Update AI assistant protocols to use semantic tokens

## Backup Location

All original files have been backed up to: \`${BACKUP_DIR}\`

## Rollback Instructions

To rollback the migration:

\`\`\`bash
# Restore from backup
cp -r "${BACKUP_DIR}"/* "${PROJECT_ROOT}"/

# Or restore specific files
cp "${BACKUP_DIR}/filename.go" "${PROJECT_ROOT}/"
\`\`\`
EOF

    log_message "REPORT" "Generated migration report: ${report_file}"
    echo -e "${BLUE}📊 Migration report generated: ${report_file}${NC}"
}

# Main migration process
main() {
    echo -e "${BLUE}🔍 Scanning for files to migrate...${NC}"
    
    # Find all relevant files
    local files_to_migrate=($(find "${PROJECT_ROOT}" -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" | grep -v "${BACKUP_DIR}" | grep -v "node_modules" | grep -v ".git"))
    
    echo -e "${BLUE}📋 Found ${#files_to_migrate[@]} files to process${NC}"
    echo ""
    
    local total_files=${#files_to_migrate[@]}
    local processed_files=0
    local successful_migrations=0
    local validation_errors=0
    
    # Process each file
    for file in "${files_to_migrate[@]}"; do
        processed_files=$((processed_files + 1))
        echo -e "${BLUE}[${processed_files}/${total_files}] Processing: ${file}${NC}"
        
        # Check if file contains Unicode icons
        if grep -q "[⭐🔺🔶🔻🔧📝🔍🛠️🛡️🔄]" "${file}" 2>/dev/null; then
            migrate_file "${file}"
            successful_migrations=$((successful_migrations + 1))
        else
            log_message "SKIP" "No Unicode icons found in ${file}"
            echo -e "${YELLOW}⏭️  Skipping ${file} (no Unicode icons)${NC}"
        fi
        
        echo ""
    done
    
    echo -e "${BLUE}🔍 Validating migration...${NC}"
    
    # Validate all migrated files
    for file in "${files_to_migrate[@]}"; do
        if ! validate_migration "${file}"; then
            validation_errors=$((validation_errors + 1))
        fi
    done
    
    echo ""
    echo -e "${BLUE}📊 Migration Summary${NC}"
    echo -e "Total files processed: ${total_files}"
    echo -e "Successful migrations: ${successful_migrations}"
    echo -e "Validation errors: ${validation_errors}"
    echo -e "Backup location: ${BACKUP_DIR}"
    echo -e "Log file: ${LOG_FILE}"
    
    # Generate report
    generate_report
    
    if [ "$validation_errors" -eq 0 ]; then
        echo -e "${GREEN}🎉 Migration completed successfully!${NC}"
        exit 0
    else
        echo -e "${RED}⚠️  Migration completed with ${validation_errors} validation errors${NC}"
        exit 1
    fi
}

# Run main function
main "$@" 