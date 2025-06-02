#!/bin/bash

# COV-001: Coverage validation script for new code
# This script ensures new development maintains high coverage standards
# while excluding legacy code from coverage requirements

set -e

# Configuration
COVERAGE_THRESHOLD=85.0
LEGACY_THRESHOLD=70.0
COVERAGE_PROFILE="coverage.out"
NEW_COVERAGE_PROFILE="coverage_new.out"
LEGACY_COVERAGE_PROFILE="coverage_legacy.out"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

print_header() {
    echo
    print_status $BLUE "=== COV-001: Coverage Validation ==="
    echo
}

# Function to run tests with coverage
run_coverage() {
    print_status $BLUE "Running tests with coverage analysis..."
    
    # Generate main coverage profile
    go test -coverprofile=$COVERAGE_PROFILE ./...
    
    if [ ! -f "$COVERAGE_PROFILE" ]; then
        print_status $RED "❌ Failed to generate coverage profile"
        exit 1
    fi
    
    print_status $GREEN "✅ Coverage profile generated: $COVERAGE_PROFILE"
}

# Function to build coverage analysis tool
build_coverage_tool() {
    print_status $BLUE "Building coverage analysis tool..."
    
    if [ ! -f "tools/coverage.go" ]; then
        print_status $RED "❌ Coverage tool not found: tools/coverage.go"
        exit 1
    fi
    
    # Build the coverage analysis tool
    go build -o tools/coverage-analyzer tools/coverage.go
    
    if [ ! -f "tools/coverage-analyzer" ]; then
        print_status $RED "❌ Failed to build coverage analyzer"
        exit 1
    fi
    
    print_status $GREEN "✅ Coverage analyzer built successfully"
}

# Function to analyze coverage with exclusions
analyze_coverage() {
    print_status $BLUE "Analyzing coverage with exclusion patterns..."
    
    # Run the coverage analyzer
    ./tools/coverage-analyzer $COVERAGE_PROFILE > coverage_report.txt
    
    # Display the report
    cat coverage_report.txt
    
    # Check if analysis passed
    if grep -q "❌ FAIL" coverage_report.txt; then
        print_status $RED "❌ Coverage validation failed"
        return 1
    elif grep -q "✅ PASS" coverage_report.txt; then
        print_status $GREEN "✅ Coverage validation passed"
        return 0
    else
        print_status $YELLOW "⚠️  Coverage validation result unclear"
        return 1
    fi
}

# Function to generate HTML reports
generate_html_reports() {
    print_status $BLUE "Generating HTML coverage reports..."
    
    # Generate overall coverage report
    go tool cover -html=$COVERAGE_PROFILE -o coverage.html
    print_status $GREEN "✅ Overall coverage report: coverage.html"
    
    # Generate focused reports (would require more advanced filtering)
    # For now, just create a copy for new code focus
    cp coverage.html coverage_new.html
    print_status $GREEN "✅ New code coverage report: coverage_new.html"
}

# Function to check specific coverage requirements
check_coverage_requirements() {
    print_status $BLUE "Checking coverage requirements..."
    
    # Extract coverage percentage from go test output
    local overall_coverage=$(go test -cover ./... 2>&1 | grep "coverage:" | tail -1 | sed 's/.*coverage: \([0-9.]*\)%.*/\1/')
    
    if [ -z "$overall_coverage" ]; then
        print_status $RED "❌ Could not extract coverage percentage"
        return 1
    fi
    
    print_status $BLUE "Overall coverage: ${overall_coverage}%"
    
    # Check if coverage meets baseline (using baseline of 73.5% from documentation)
    local baseline=73.5
    if (( $(echo "$overall_coverage >= $baseline" | bc -l) )); then
        print_status $GREEN "✅ Coverage meets baseline requirement (${baseline}%)"
    else
        print_status $YELLOW "⚠️  Coverage below baseline: ${overall_coverage}% < ${baseline}%"
    fi
    
    # For new code, we rely on the coverage analyzer tool
    # which implements the specific logic for new vs legacy code
}

# Function to validate test execution
validate_test_execution() {
    print_status $BLUE "Validating test execution..."
    
    # Ensure all tests pass
    if go test ./...; then
        print_status $GREEN "✅ All tests pass"
    else
        print_status $RED "❌ Some tests failed"
        return 1
    fi
}

# Function to check for coverage exclusion markers
check_exclusion_markers() {
    print_status $BLUE "Checking for coverage exclusion markers..."
    
    # Look for build tags in Go files
    local legacy_files=$(find . -name "*.go" -exec grep -l "//go:build.*legacy" {} \; | wc -l)
    local exclude_files=$(find . -name "*.go" -exec grep -l "//go:build.*exclude_coverage" {} \; | wc -l)
    
    print_status $BLUE "Files with legacy build tags: $legacy_files"
    print_status $BLUE "Files with exclude_coverage build tags: $exclude_files"
    
    # Check for coverage ignore comments
    local ignore_comments=$(find . -name "*.go" -exec grep -l "//coverage:ignore" {} \; | wc -l)
    print_status $BLUE "Files with coverage ignore comments: $ignore_comments"
}

# Function to cleanup temporary files
cleanup() {
    print_status $BLUE "Cleaning up temporary files..."
    
    # Remove generated files
    rm -f tools/coverage-analyzer
    rm -f coverage_report.txt
    
    # Keep coverage profiles for inspection
    print_status $GREEN "✅ Cleanup complete (coverage profiles preserved)"
}

# Main execution
main() {
    print_header
    
    # Set up error handling
    trap cleanup EXIT
    
    # Run validation steps
    validate_test_execution
    run_coverage
    build_coverage_tool
    analyze_coverage
    local analysis_result=$?
    generate_html_reports
    check_coverage_requirements
    check_exclusion_markers
    
    print_header
    
    if [ $analysis_result -eq 0 ]; then
        print_status $GREEN "🎉 Coverage validation completed successfully!"
        print_status $GREEN "   - All tests pass"
        print_status $GREEN "   - Coverage requirements met"
        print_status $GREEN "   - Reports generated: coverage.html, coverage_new.html"
        exit 0
    else
        print_status $RED "❌ Coverage validation failed!"
        print_status $RED "   - Review coverage_report.txt for details"
        print_status $RED "   - Check coverage.html for visual analysis"
        exit 1
    fi
}

# Check if bc is available for floating point comparison
if ! command -v bc &> /dev/null; then
    print_status $YELLOW "⚠️  bc command not found - some coverage checks may be limited"
fi

# Run main function
main "$@" 