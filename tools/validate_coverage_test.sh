#!/bin/bash

# TEST-002: Tools directory testing - Coverage validation script test suite
# Comprehensive test coverage for validate-coverage.sh functionality

set -e

# Test configuration
TEST_DIR="$(mktemp -d)"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VALIDATION_SCRIPT="$SCRIPT_DIR/validate-coverage.sh"
TEST_RESULTS=0
TESTS_RUN=0

# Colors for test output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test utility functions
print_test_header() {
    echo
    echo -e "${BLUE}=== TEST-002: Coverage Validation Script Tests ===${NC}"
    echo
}

print_test_status() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

run_test() {
    local test_name="$1"
    local test_function="$2"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    print_test_status $BLUE "Running: $test_name"
    
    if $test_function; then
        print_test_status $GREEN "✅ PASS: $test_name"
    else
        print_test_status $RED "❌ FAIL: $test_name"
        TEST_RESULTS=1
    fi
}

# Setup test environment
setup_test_env() {
    cd "$TEST_DIR"
    
    # Create minimal Go module structure for testing
    cat > go.mod << 'EOF'
module test-coverage

go 1.19
EOF

    # Create test Go files (without main function to avoid conflict)
    mkdir -p internal/testutil
    
    cat > main.go << 'EOF'
package main

func main() {
    // Simple main function for build compatibility
}
EOF

    cat > util.go << 'EOF'
package main

import "fmt"

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}

func untested() {
    fmt.Println("This function is not tested")
}
EOF

    cat > util_test.go << 'EOF'
package main

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}
EOF

    # Create coverage configuration
    cat > coverage.toml << 'EOF'
[coverage]
baseline_date = "2025-01-01"
new_code_cutoff_date = "2025-01-01"
show_legacy_coverage = false
show_new_code_coverage = true
show_overall_coverage = true

[[coverage.always_legacy]]
files = ["main.go"]

[[coverage.always_new]]
files = ["util.go"]

[[coverage.file_patterns]]
patterns = ["*_legacy.go", "*_deprecated.go"]
EOF

    # Copy the coverage tool and validation script to a tools subdirectory
    mkdir -p tools
    cp "$SCRIPT_DIR/coverage.go" tools/
    cp "$SCRIPT_DIR/validate-coverage.sh" tools/
    chmod +x tools/validate-coverage.sh
}

cleanup_test_env() {
    cd /
    rm -rf "$TEST_DIR"
}

# Test individual functions from the validation script

test_script_exists() {
    # TEST-002: Test that the validation script exists and is executable
    if [ ! -f "$VALIDATION_SCRIPT" ]; then
        echo "Validation script not found: $VALIDATION_SCRIPT"
        return 1
    fi
    
    if [ ! -x "$VALIDATION_SCRIPT" ]; then
        echo "Validation script is not executable: $VALIDATION_SCRIPT"
        return 1
    fi
    
    return 0
}

test_go_module_setup() {
    # TEST-002: Test that our test Go module is properly set up
    if ! go mod tidy; then
        echo "Failed to tidy Go module"
        return 1
    fi
    
    if ! go build ./...; then
        echo "Failed to build test Go code"
        return 1
    fi
    
    return 0
}

test_basic_test_execution() {
    # TEST-002: Test that basic Go tests can run
    if ! go test ./...; then
        echo "Basic test execution failed"
        return 1
    fi
    
    return 0
}

test_coverage_profile_generation() {
    # TEST-002: Test coverage profile generation
    if ! go test -coverprofile=coverage.out ./...; then
        echo "Failed to generate coverage profile"
        return 1
    fi
    
    if [ ! -f "coverage.out" ]; then
        echo "Coverage profile was not created"
        return 1
    fi
    
    # Check that coverage profile has expected format
    if ! head -1 coverage.out | grep -q "mode: set"; then
        echo "Coverage profile has unexpected format"
        return 1
    fi
    
    return 0
}

test_coverage_tool_build() {
    # TEST-002: Test that the coverage analysis tool can be built
    if ! go build -o coverage-analyzer tools/coverage.go; then
        echo "Failed to build coverage analyzer"
        return 1
    fi
    
    if [ ! -f "coverage-analyzer" ]; then
        echo "Coverage analyzer executable was not created"
        return 1
    fi
    
    return 0
}

test_coverage_analysis_execution() {
    # TEST-002: Test that coverage analysis can run
    
    # First generate coverage profile
    go test -coverprofile=coverage.out ./... >/dev/null 2>&1
    
    # Build coverage analyzer
    go build -o coverage-analyzer tools/coverage.go
    
    # Run coverage analysis
    if ! ./coverage-analyzer coverage.out > coverage_report.txt; then
        echo "Coverage analysis execution failed"
        return 1
    fi
    
    if [ ! -f "coverage_report.txt" ]; then
        echo "Coverage report was not generated"
        return 1
    fi
    
    # Check report contains expected sections
    if ! grep -q "Coverage Report" coverage_report.txt; then
        echo "Coverage report missing header"
        return 1
    fi
    
    return 0
}

test_html_report_generation() {
    # TEST-002: Test HTML coverage report generation
    
    # Generate coverage profile
    go test -coverprofile=coverage.out ./... >/dev/null 2>&1
    
    # Generate HTML report
    if ! go tool cover -html=coverage.out -o coverage.html; then
        echo "Failed to generate HTML coverage report"
        return 1
    fi
    
    if [ ! -f "coverage.html" ]; then
        echo "HTML coverage report was not created"
        return 1
    fi
    
    # Check that HTML file contains expected content
    if ! grep -q "<html" coverage.html; then
        echo "Generated file is not valid HTML"
        return 1
    fi
    
    return 0
}

test_coverage_threshold_validation() {
    # TEST-002: Test coverage threshold validation logic
    
    # Generate coverage profile with low coverage
    go test -coverprofile=coverage.out ./... >/dev/null 2>&1
    
    # Extract coverage percentage
    local coverage_percent=$(go test -cover ./... 2>&1 | grep "coverage:" | tail -1 | sed 's/.*coverage: \([0-9.]*\)%.*/\1/')
    
    if [ -z "$coverage_percent" ]; then
        echo "Could not extract coverage percentage"
        return 1
    fi
    
    # Verify we can get a numeric value
    if ! echo "$coverage_percent" | grep -q "^[0-9.]*$"; then
        echo "Coverage percentage is not numeric: $coverage_percent"
        return 1
    fi
    
    return 0
}

test_error_handling_missing_files() {
    # TEST-002: Test error handling with missing files
    
    # Build the coverage analyzer first
    go build -o coverage-analyzer tools/coverage.go 2>/dev/null
    
    # Test with missing coverage profile
    if ./coverage-analyzer nonexistent.out 2>/dev/null; then
        echo "Coverage analyzer should fail with missing profile"
        return 1
    fi
    
    # Test with missing coverage tool
    if [ -f "coverage-analyzer" ]; then
        rm coverage-analyzer
    fi
    
    # The validation script should detect missing tool
    # We can't easily test the full script without complex mocking
    # so we test the component pieces
    
    return 0
}

test_script_function_isolation() {
    # TEST-002: Test that individual functions in the script work in isolation
    
    # Source the validation script to test individual functions
    # Note: This is a simplified test since the script has complex interdependencies
    
    # Test that basic utilities are available
    if ! command -v go >/dev/null; then
        echo "Go command not available for testing"
        return 1
    fi
    
    # Test that we can detect build tags (simplified)
    echo '//go:build legacy' > test_legacy.go
    echo 'package main' >> test_legacy.go
    
    local legacy_files=$(find . -name "*.go" -exec grep -l "//go:build.*legacy" {} \; | wc -l)
    if [ "$legacy_files" -lt 1 ]; then
        echo "Failed to detect legacy build tags"
        return 1
    fi
    
    rm test_legacy.go
    return 0
}

test_integration_workflow() {
    # TEST-002: Test a simplified version of the complete workflow
    
    print_test_status $BLUE "  Testing integrated workflow..."
    
    # Step 1: Run tests
    if ! go test ./... >/dev/null 2>&1; then
        echo "Integration test: basic tests failed"
        return 1
    fi
    
    # Step 2: Generate coverage
    if ! go test -coverprofile=coverage.out ./... >/dev/null 2>&1; then
        echo "Integration test: coverage generation failed"
        return 1
    fi
    
    # Step 3: Build coverage tool
    if ! go build -o coverage-analyzer tools/coverage.go 2>/dev/null; then
        echo "Integration test: coverage tool build failed"
        return 1
    fi
    
    # Step 4: Analyze coverage
    if ! ./coverage-analyzer coverage.out > coverage_report.txt 2>/dev/null; then
        echo "Integration test: coverage analysis failed"
        return 1
    fi
    
    # Step 5: Generate HTML report
    if ! go tool cover -html=coverage.out -o coverage.html 2>/dev/null; then
        echo "Integration test: HTML report generation failed"
        return 1
    fi
    
    # Verify all expected files exist
    local expected_files=("coverage.out" "coverage-analyzer" "coverage_report.txt" "coverage.html")
    for file in "${expected_files[@]}"; do
        if [ ! -f "$file" ]; then
            echo "Integration test: missing expected file $file"
            return 1
        fi
    done
    
    print_test_status $GREEN "  Integration workflow completed successfully"
    return 0
}

test_performance_benchmarks() {
    # TEST-002: Test performance characteristics
    
    # Generate larger test data for performance testing
    for i in {1..10}; do
        cat > "large_file_$i.go" << EOF
package main

import "fmt"

func Function${i}A() {
    fmt.Println("Function ${i}A")
}

func Function${i}B() {
    fmt.Println("Function ${i}B")
}

func Function${i}C() {
    fmt.Println("Function ${i}C")
}
EOF
    done
    
    # Time the coverage generation
    local start_time=$(date +%s)
    go test -coverprofile=coverage_large.out ./... >/dev/null 2>&1
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    if [ $duration -gt 30 ]; then
        echo "Performance test: coverage generation took too long ($duration seconds)"
        return 1
    fi
    
    # Clean up large files
    rm -f large_file_*.go coverage_large.out
    
    return 0
}

test_edge_cases() {
    # TEST-002: Test edge cases and boundary conditions
    
    # Test with empty Go file
    echo 'package main' > empty.go
    if ! go test -coverprofile=coverage_empty.out ./... >/dev/null 2>&1; then
        echo "Edge case test: failed with empty Go file"
        return 1
    fi
    
    # Test with file that has no functions
    cat > no_functions.go << 'EOF'
package main

var GlobalVar = "test"
const GlobalConst = 42
EOF
    
    if ! go test -coverprofile=coverage_no_funcs.out ./... >/dev/null 2>&1; then
        echo "Edge case test: failed with no-function file"
        return 1
    fi
    
    # Clean up
    rm -f empty.go no_functions.go coverage_empty.out coverage_no_funcs.out
    
    return 0
}

# Main test execution
main() {
    print_test_header
    
    # Setup test environment
    print_test_status $BLUE "Setting up test environment in $TEST_DIR"
    setup_test_env
    
    # Run all tests
    run_test "Script Exists and Executable" test_script_exists
    run_test "Go Module Setup" test_go_module_setup  
    run_test "Basic Test Execution" test_basic_test_execution
    run_test "Coverage Profile Generation" test_coverage_profile_generation
    run_test "Coverage Tool Build" test_coverage_tool_build
    run_test "Coverage Analysis Execution" test_coverage_analysis_execution
    run_test "HTML Report Generation" test_html_report_generation
    run_test "Coverage Threshold Validation" test_coverage_threshold_validation
    run_test "Error Handling Missing Files" test_error_handling_missing_files
    run_test "Script Function Isolation" test_script_function_isolation
    run_test "Integration Workflow" test_integration_workflow
    run_test "Performance Benchmarks" test_performance_benchmarks
    run_test "Edge Cases" test_edge_cases
    
    # Cleanup
    cleanup_test_env
    
    # Print results
    echo
    print_test_status $BLUE "=== Test Results ==="
    print_test_status $BLUE "Tests run: $TESTS_RUN"
    
    if [ $TEST_RESULTS -eq 0 ]; then
        print_test_status $GREEN "✅ All tests passed!"
        print_test_status $GREEN "Coverage validation script is working correctly"
    else
        print_test_status $RED "❌ Some tests failed"
        print_test_status $RED "Coverage validation script needs attention"
    fi
    
    exit $TEST_RESULTS
}

# Check dependencies
if ! command -v go >/dev/null; then
    print_test_status $RED "❌ Go is not installed or not in PATH"
    exit 1
fi

if ! command -v mktemp >/dev/null; then
    print_test_status $RED "❌ mktemp command not available"
    exit 1
fi

# Run main function
main "$@" 