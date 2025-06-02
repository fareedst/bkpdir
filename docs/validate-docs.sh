#!/bin/bash

# Documentation Validation Script for BkpDir
# Checks consistency between documentation layers and implementation

set -e

echo "🔍 BkpDir Documentation Validation"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
ERRORS=0
WARNINGS=0
CHECKS=0

# Helper functions
error() {
    echo -e "${RED}❌ ERROR: $1${NC}"
    ((ERRORS++))
}

warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
    ((WARNINGS++))
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

check() {
    echo "🔹 $1"
    ((CHECKS++))
}

# DOC-001: Semantic linking validation functions
validate_markdown_link() {
    local link="$1"
    local source_file="$2"
    
    # Extract file path and anchor from markdown link [text](file#anchor)
    # Use grep for more reliable pattern matching
    local link_text=$(echo "$link" | grep -o '\[[^\]]*\]' | sed 's/\[//g' | sed 's/\]//g')
    local link_target=$(echo "$link" | grep -o '([^)]\+)' | sed 's/(//g' | sed 's/)//g')
    
    if [[ -z "$link_target" ]]; then
        echo "INVALID_FORMAT:$source_file:$link:Invalid markdown link format"
        return 1
    fi
    
    # Split file and anchor
    local file_path=""
    local anchor=""
    if [[ $link_target == *"#"* ]]; then
        file_path="${link_target%#*}"
        anchor="${link_target#*#}"
    else
        file_path="$link_target"
    fi
    
    # Resolve relative path from source file location
    local source_dir=$(dirname "$source_file")
    local resolved_path=""
    if [[ $file_path == /* ]]; then
        # Absolute path
        resolved_path="$file_path"
    elif [[ $file_path == http* ]]; then
        # URL - skip validation
        return 0
    else
        # Relative path
        resolved_path="$source_dir/$file_path"
    fi
    
    # Check if target file exists
    if [[ -f "$resolved_path" ]]; then
        # Check if anchor exists in target file (if specified)
        if [[ -n "$anchor" ]]; then
            # Look for the anchor in the target file
            if grep -q "$anchor" "$resolved_path"; then
                return 0  # Link is valid
            else
                echo "INVALID_ANCHOR:$source_file:$link_target:Anchor '$anchor' not found in $resolved_path"
                return 1
            fi
        else
            return 0  # File exists, no anchor to check
        fi
    else
        echo "INVALID_FILE:$source_file:$link_target:File not found: $resolved_path"
        return 1
    fi
}

validate_feature_reference_format() {
    local file="$1"
    
    # Look for feature reference blocks in the format specified in cross-reference-template.md
    local feature_refs=$(grep -n "^\*\*Feature\*\*:" "$file" 2>/dev/null || true)
    
    while IFS= read -r line; do
        if [[ -n "$line" ]]; then
            local line_num=$(echo "$line" | cut -d: -f1)
            local feature_line=$(echo "$line" | cut -d: -f2-)
            
            # Extract feature ID from the line using grep
            local feature_id=$(echo "$feature_line" | grep -o '[A-Z]\+-[0-9]\+' | head -1)
            
            if [[ -n "$feature_id" ]]; then
                # Check if complete feature reference block exists
                local block_start=$line_num
                local block_end=$((line_num + 5))
                
                local has_spec=false
                local has_requirements=false
                local has_architecture=false
                local has_tests=false
                local has_code=false
                
                for ((i=block_start; i<=block_end; i++)); do
                    local check_line=$(sed -n "${i}p" "$file" 2>/dev/null || true)
                    if echo "$check_line" | grep -q "^\*\*Spec\*\*:"; then
                        has_spec=true
                    elif echo "$check_line" | grep -q "^\*\*Requirements\*\*:"; then
                        has_requirements=true
                    elif echo "$check_line" | grep -q "^\*\*Architecture\*\*:"; then
                        has_architecture=true
                    elif echo "$check_line" | grep -q "^\*\*Tests\*\*:"; then
                        has_tests=true
                    elif echo "$check_line" | grep -q "^\*\*Code\*\*:"; then
                        has_code=true
                    fi
                done
                
                # Report missing components
                if ! $has_spec; then
                    echo "MISSING_SPEC_REF:$file:$line_num:Feature $feature_id missing Spec reference"
                fi
                if ! $has_requirements; then
                    echo "MISSING_REQ_REF:$file:$line_num:Feature $feature_id missing Requirements reference"
                fi
                if ! $has_architecture; then
                    echo "MISSING_ARCH_REF:$file:$line_num:Feature $feature_id missing Architecture reference"
                fi
                if ! $has_tests; then
                    echo "MISSING_TEST_REF:$file:$line_num:Feature $feature_id missing Tests reference"
                fi
                if ! $has_code; then
                    echo "MISSING_CODE_REF:$file:$line_num:Feature $feature_id missing Code reference"
                fi
                
                if $has_spec && $has_requirements && $has_architecture && $has_tests && $has_code; then
                    echo "VALID_FEATURE_REF:$file:$line_num:Feature $feature_id has complete reference block"
                fi
            fi
        fi
    done <<< "$feature_refs"
}

check_orphaned_references() {
    local file="$1"
    
    # Extract all markdown links from the file
    local links=$(grep -o '\[[^\]]*\]([^)]+)' "$file" 2>/dev/null || true)
    
    local orphaned_count=0
    local valid_count=0
    
    while IFS= read -r link; do
        if [[ -n "$link" ]]; then
            if validate_markdown_link "$link" "$file" >/dev/null 2>&1; then
                ((valid_count++))
            else
                ((orphaned_count++))
            fi
        fi
    done <<< "$links"
    
    echo "LINK_SUMMARY:$file:$valid_count:$orphaned_count"
}

validate_cross_document_consistency() {
    # Check that feature IDs are consistently named across documents
    local feature_tracking="docs/context/feature-tracking.md"
    
    if [[ -f "$feature_tracking" ]]; then
        # Extract feature IDs and their descriptions from feature tracking
        local features=$(grep -E '^\| [A-Z]+-[0-9]+ \|' "$feature_tracking" | cut -d'|' -f2 | sed 's/^ *//;s/ *$//')
        
        for feature_id in $features; do
            local ref_count=0
            local files_with_refs=""
            
            # Count references across all documentation files
            for doc_file in docs/context/*.md; do
                if [[ -f "$doc_file" ]] && [[ "$doc_file" != "$feature_tracking" ]]; then
                    local count=$(grep -c "$feature_id" "$doc_file" 2>/dev/null || echo "0")
                    if [[ $count -gt 0 ]]; then
                        ref_count=$((ref_count + count))
                        files_with_refs="$files_with_refs $(basename "$doc_file")"
                    fi
                fi
            done
            
            echo "FEATURE_CONSISTENCY:$feature_id:$ref_count:$files_with_refs"
        done
    fi
}

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -d "docs/context" ]]; then
    error "Please run this script from the project root directory"
    exit 1
fi

echo
check "Checking documentation files exist"

# Check required documentation files
required_docs=(
    "docs/context/immutable.md"
    "docs/context/requirements.md"
    "docs/context/specification.md"
    "docs/context/architecture.md"
    "docs/context/testing.md"
    "docs/context/feature-tracking.md"
    "docs/context/doc-validation.md"
)

for doc in "${required_docs[@]}"; do
    if [[ -f "$doc" ]]; then
        success "Found $doc"
    else
        error "Missing required documentation: $doc"
    fi
done

echo
check "Checking feature ID consistency"

# Extract feature IDs from feature tracking
if [[ -f "docs/context/feature-tracking.md" ]]; then
    feature_ids=$(grep -o '\(ARCH\|FILE\|CFG\|GIT\|DOC\)-[0-9]\+' docs/context/feature-tracking.md | sort -u)
    
    # Check if feature IDs are referenced in documentation
    for feature_id in $feature_ids; do
        count=$(grep -r "$feature_id" docs/context/ | wc -l)
        if [[ $count -gt 1 ]]; then
            success "Feature ID $feature_id referenced $count times"
        else
            warning "Feature ID $feature_id only referenced once"
        fi
    done
else
    error "Feature tracking document not found"
fi

echo
check "Checking implementation tokens"

# Check if feature IDs are present in Go code
go_files=(*.go)
for go_file in "${go_files[@]}"; do
    if [[ -f "$go_file" ]]; then
        token_count=$(grep -c "// \(ARCH\|FILE\|CFG\|GIT\|DOC\)-[0-9]\+" "$go_file" 2>/dev/null || true)
        if [[ -z "$token_count" ]]; then
            token_count=0
        fi
        if [[ $token_count -gt 0 ]]; then
            success "Found $token_count implementation tokens in $go_file"
        else
            warning "No implementation tokens found in $go_file"
        fi
    fi
done

echo
check "Checking test references"

# Check if test functions referenced in feature tracking exist
if [[ -f "docs/context/feature-tracking.md" ]]; then
    test_refs=$(grep -o 'Test[A-Za-z][A-Za-z0-9]*' docs/context/feature-tracking.md | grep -v '^Testing$' | grep -v '^Tests$' | sort -u)
    
    for test_ref in $test_refs; do
        if grep -r "func $test_ref" *_test.go >/dev/null 2>&1; then
            success "Test function $test_ref exists"
        else
            error "Referenced test function $test_ref not found"
        fi
    done
fi

echo
check "Checking decision references"

# Check if decision IDs are properly documented
decision_ids=$(grep -o 'DEC-[0-9]\+' docs/context/feature-tracking.md 2>/dev/null | sort -u || echo "")

if [[ -n "$decision_ids" ]]; then
    for decision_id in $decision_ids; do
        if grep -q "#### $decision_id:" docs/context/feature-tracking.md; then
            success "Decision $decision_id is documented"
        else
            error "Decision $decision_id referenced but not documented"
        fi
    done
else
    warning "No decision IDs found in feature tracking"
fi

echo
check "Checking immutable requirements consistency"

# Check if immutable sections are referenced properly
if [[ -f "docs/context/immutable.md" ]]; then
    immutable_sections=$(grep "^## " docs/context/immutable.md | sed 's/^## //' | head -5)
    
    for section in $immutable_sections; do
        ref_count=$(grep -r "IMMUTABLE-REF.*$section" . 2>/dev/null | wc -l || echo "0")
        if [[ $ref_count -gt 0 ]]; then
            success "Immutable section '$section' referenced $ref_count times"
        else
            warning "Immutable section '$section' not referenced in code"
        fi
    done
fi

echo
check "Checking configuration consistency"

# Check if configuration options in immutable.md match defaults in code
if [[ -f "docs/context/immutable.md" ]] && [[ -f "config.go" ]]; then
    # This is a basic check - could be expanded
    if grep -q "DefaultConfig" config.go; then
        success "DefaultConfig function found in config.go"
    else
        error "DefaultConfig function not found in config.go"
    fi
fi

echo
check "DOC-001: Validating semantic cross-references"

# DOC-001: Semantic cross-reference validation
link_validation_results=""
feature_ref_validation_results=""
orphaned_ref_results=""
consistency_results=""

for doc_file in docs/context/*.md; do
    if [[ -f "$doc_file" ]]; then
        info "Validating links in $(basename "$doc_file")"
        
        # Validate markdown links
        local links=$(grep -o '\[[^\]]*\]([^)]+)' "$doc_file" 2>/dev/null || true)
        local link_errors=0
        local link_count=0
        
        while IFS= read -r link; do
            if [[ -n "$link" ]]; then
                ((link_count++))
                if ! validate_markdown_link "$link" "$doc_file"; then
                    ((link_errors++))
                fi
            fi
        done <<< "$links"
        
        if [[ $link_errors -eq 0 ]] && [[ $link_count -gt 0 ]]; then
            success "All $link_count links valid in $(basename "$doc_file")"
        elif [[ $link_count -eq 0 ]]; then
            info "No links found in $(basename "$doc_file")"
        else
            error "$link_errors invalid links found in $(basename "$doc_file")"
        fi
        
        # Validate feature reference format
        local feature_ref_output=$(validate_feature_reference_format "$doc_file")
        if echo "$feature_ref_output" | grep -q "VALID_FEATURE_REF"; then
            local valid_refs=$(echo "$feature_ref_output" | grep -c "VALID_FEATURE_REF" || echo "0")
            success "Found $valid_refs complete feature reference blocks in $(basename "$doc_file")"
        fi
        
        if echo "$feature_ref_output" | grep -q "MISSING_.*_REF"; then
            local missing_refs=$(echo "$feature_ref_output" | grep "MISSING_.*_REF" | wc -l)
            warning "$missing_refs incomplete feature references in $(basename "$doc_file")"
        fi
        
        # Check orphaned references
        check_orphaned_references "$doc_file"
    fi
done

# Cross-document consistency validation
info "Validating cross-document consistency"
consistency_output=$(validate_cross_document_consistency)

# Parse and report consistency results
echo "$consistency_output" | while IFS=: read -r type feature_id ref_count files; do
    if [[ "$type" == "FEATURE_CONSISTENCY" ]]; then
        if [[ $ref_count -ge 3 ]]; then
            success "Feature $feature_id consistently referenced ($ref_count times) in:$files"
        elif [[ $ref_count -ge 1 ]]; then
            warning "Feature $feature_id minimally referenced ($ref_count times) in:$files"
        else
            error "Feature $feature_id not found in other documentation files"
        fi
    fi
done

echo
echo "📊 Validation Summary"
echo "===================="
echo "Total checks: $CHECKS"
echo -e "${GREEN}✅ Successful validations${NC}"
echo -e "${YELLOW}⚠️  Warnings: $WARNINGS${NC}"
echo -e "${RED}❌ Errors: $ERRORS${NC}"

if [[ $ERRORS -gt 0 ]]; then
    echo
    echo -e "${RED}❌ Validation failed with $ERRORS errors${NC}"
    exit 1
elif [[ $WARNINGS -gt 0 ]]; then
    echo
    echo -e "${YELLOW}⚠️  Validation completed with $WARNINGS warnings${NC}"
    exit 0
else
    echo
    echo -e "${GREEN}✅ All validations passed successfully!${NC}"
    exit 0
fi 