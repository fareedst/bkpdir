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

check() {
    echo "🔹 $1"
    ((CHECKS++))
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
    feature_ids=$(grep -o '\(ARCH\|FILE\|CFG\|GIT\)-[0-9]\+' docs/context/feature-tracking.md | sort -u)
    
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
        token_count=$(grep -c "// \(ARCH\|FILE\|CFG\|GIT\)-[0-9]\+" "$go_file" 2>/dev/null || true)
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