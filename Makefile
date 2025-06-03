# Makefile for bkpdir - Directory archiving CLI
# 
# Available targets:
#   Development:
#     build-local     - Build for local development (current platform)
#     dev             - Build and run basic functionality test
#     install         - Install to ~/.local/bin (requires build-local first)
#   
#   Testing:
#     test            - Run all tests
#     test-verbose    - Run tests with verbose output
#     test-coverage   - Run tests with coverage report
#     test-coverage-new - Run tests with selective coverage (COV-001)
#     test-coverage-validate - Validate coverage with exclusion patterns
#     test-race       - Run tests with race detection
#     test-bench      - Run benchmark tests
#     test-all        - Run all test variants
#   
#   Code Quality:
#     lint            - Run linter (revive)
#     fmt             - Format code with gofmt
#     vet             - Run go vet
#     check           - Run all code quality checks
#   
#   Icon Validation:
#     validate-icons  - Validate implementation token icon consistency (DOC-007)
#     validate-icon-enforcement - Comprehensive icon validation and enforcement (DOC-008)
#     validate-icons-strict - Run DOC-008 validation in strict mode for CI/CD
#   
#   Production Builds:
#     build-all       - Build for all platforms
#     build-macos     - Build for macOS (ARM64 and AMD64)
#     build-ubuntu    - Build for Ubuntu (20.04, 22.04, 24.04)
#   
#   Utilities:
#     clean           - Clean build artifacts and test cache
#     deps            - Download and verify dependencies
#     help            - Show this help message

.PHONY: build-all build-ubuntu20 build-ubuntu22 build-ubuntu24 build-macos build-macos-arm64 build-macos-amd64 build-local clean
.PHONY: test test-verbose test-coverage test-coverage-new test-coverage-validate test-race test-bench test-all
.PHONY: lint validate-icons validate-icon-enforcement validate-icons-strict fmt vet check dev install deps help
.PHONY: token-migration-dry-run token-migration token-migration-rollback

# Variables
BINARY_NAME := bkpdir
VERSION := $(shell grep -o 'Version = "[^"]*"' main.go | cut -d'"' -f2)
BUILD_TIME := $(shell date -u +%Y-%m-%d\ %H:%M:%S\ UTC)
PLATFORM := $(shell go env GOOS)-$(shell go env GOARCH)
LDFLAGS := -X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=$(PLATFORM)'

# COV-001: Coverage configuration
COVERAGE_PROFILE := coverage.out
COVERAGE_NEW_PROFILE := coverage_new.out
COVERAGE_LEGACY_PROFILE := coverage_legacy.out
COVERAGE_THRESHOLD := 85.0

# COV-002: Enhanced coverage configuration
COVERAGE_BASELINE_FILE := docs/coverage-baseline.md
COVERAGE_HISTORY_FILE := docs/coverage-history.json
COVERAGE_REPORTS_DIR := coverage_reports
DIFFERENTIAL_TOOL := tools/coverage-differential

# Default target
all: check test build-local token-suggester

# Help target
help:
	@echo "bkpdir Makefile - Directory archiving CLI"
	@echo ""
	@echo "Development targets:"
	@echo "  build-local     Build for local development (current platform)"
	@echo "  dev             Build and run basic functionality test"
	@echo "  install         Install to ~/.local/bin"
	@echo ""
	@echo "Testing targets:"
	@echo "  test            Run all tests"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  test-coverage   Run tests with coverage report (legacy)"
	@echo "  test-coverage-new Run tests with selective coverage (COV-001)"
	@echo "  test-coverage-validate Validate coverage with exclusion patterns"
	@echo "  test-coverage-baseline Generate coverage baseline documentation (COV-002)"
	@echo "  test-coverage-differential Generate differential coverage report (COV-002)"
	@echo "  test-coverage-trends Update coverage trends and history (COV-002)"
	@echo "  test-coverage-full Run full coverage suite with COV-002 features"
	@echo "  test-coverage-quality-gates Validate quality gates with enhanced reporting (COV-002)"
	@echo "  test-race       Run tests with race detection"
	@echo "  test-bench      Run benchmark tests"
	@echo "  test-all        Run all test variants"
	@echo ""
	@echo "Code quality targets:"
	@echo "  lint            Run linter (revive) with icon validation"
	@echo "  validate-icons  Validate implementation token icon consistency (DOC-007)"
	@echo "  validate-icon-enforcement Comprehensive icon validation and enforcement (DOC-008)"
	@echo "  validate-icons-strict Run DOC-008 validation in strict mode for CI/CD"
	@echo "  fmt             Format code with gofmt"
	@echo "  vet             Run go vet"
	@echo "  check           Run all code quality checks"
	@echo ""
	@echo "Token standardization targets (DOC-009):"
	@echo "  standardize-tokens Complete token standardization workflow (analysis + dry run)"
	@echo "  token-migration-dry-run Run DOC-009 token migration in dry-run mode"
	@echo "  token-migration Run DOC-009 mass token standardization"
	@echo "  token-migration-rollback Rollback DOC-009 token migration"
	@echo "  analyze-priority-icons Generate priority icon mappings from feature tracking"
	@echo "  validate-token-priorities Validate token priorities against feature tracking"
	@echo "  suggest-action-icons Generate action icon suggestions for tokens"
	@echo ""
	@echo "Production build targets:"
	@echo "  build-all       Build for all platforms"
	@echo "  build-macos     Build for macOS (ARM64 and AMD64)"
	@echo "  build-ubuntu    Build for Ubuntu (20.04, 22.04, 24.04)"
	@echo ""
	@echo "Utility targets:"
	@echo "  clean           Clean build artifacts and test cache"
	@echo "  deps            Download and verify dependencies"
	@echo "  help            Show this help message"

# Development targets
build-local:
	@echo "Building $(BINARY_NAME) for local development..."
	go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME)
	@echo "✓ Built $(BINARY_NAME) successfully"

dev: build-local
	@echo "Running basic functionality test..."
	@./$(BINARY_NAME) --help > /dev/null && echo "✓ Help command works"
	@./$(BINARY_NAME) --config > /dev/null && echo "✓ Config command works"
	@./$(BINARY_NAME) full --dry-run > /dev/null && echo "✓ Dry-run works"
	@echo "✓ Basic functionality test passed"

install: build-local
	@echo "Installing $(BINARY_NAME) to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp $(BINARY_NAME) ~/.local/bin/$(BINARY_NAME)
	@echo "✓ Installed $(BINARY_NAME) to ~/.local/bin/$(BINARY_NAME)"
	@echo "  Make sure ~/.local/bin is in your PATH"

# Testing targets
test:
	@echo "Running tests..."
	go test ./...
	@echo "✓ All tests passed"

test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Legacy coverage target (preserved for compatibility)
test-coverage:
	@echo "Running tests with coverage (legacy mode)..."
	go test -cover ./...
	@echo "Generating detailed coverage report..."
	go test -coverprofile=$(COVERAGE_PROFILE) ./...
	go tool cover -html=$(COVERAGE_PROFILE) -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# COV-001: New selective coverage target
test-coverage-new:
	@echo "Running tests with selective coverage (COV-001)..."
	@echo "Focusing on new code development while preserving legacy test execution..."
	go test -coverprofile=$(COVERAGE_PROFILE) ./...
	@echo "Building coverage analysis tools..."
	@mkdir -p tools
	@if [ -f "tools/coverage.go" ]; then \
		go build -o tools/coverage-analyzer tools/coverage.go; \
		echo "Analyzing coverage with exclusion patterns..."; \
		./tools/coverage-analyzer $(COVERAGE_PROFILE); \
		echo "Generating HTML reports..."; \
		go tool cover -html=$(COVERAGE_PROFILE) -o coverage.html; \
		cp coverage.html coverage_new.html; \
		echo "✓ Selective coverage analysis complete"; \
		echo "  - Overall report: coverage.html"; \
		echo "  - New code focus: coverage_new.html"; \
		rm -f tools/coverage-analyzer; \
	else \
		echo "⚠️  Coverage exclusion tool not found, falling back to standard coverage"; \
		go test -cover ./...; \
		go tool cover -html=$(COVERAGE_PROFILE) -o coverage.html; \
	fi

# COV-001: Coverage validation with quality gates
test-coverage-validate:
	@echo "Validating coverage with exclusion patterns and quality gates..."
	@if [ -f "tools/validate-coverage.sh" ]; then \
		chmod +x tools/validate-coverage.sh; \
		./tools/validate-coverage.sh; \
	else \
		echo "⚠️  Coverage validation script not found, running basic coverage check"; \
		$(MAKE) test-coverage-new; \
	fi

# COV-002: Generate coverage baseline documentation
test-coverage-baseline:
	@echo "Generating coverage baseline (COV-002)..."
	@echo "Capturing current coverage metrics..."
	go test -coverprofile=$(COVERAGE_PROFILE) . ./internal/testutil
	@echo "Updating baseline documentation..."
	@mkdir -p docs
	@echo "# Coverage Baseline Documentation (COV-002)" > $(COVERAGE_BASELINE_FILE)
	@echo "" >> $(COVERAGE_BASELINE_FILE)
	@echo "> **Generated on:** \`$(shell date -u +%Y-%m-%d\ %H:%M:%S)\ UTC\`" >> $(COVERAGE_BASELINE_FILE)
	@echo "> **Overall Coverage:** \`$(shell go tool cover -func=$(COVERAGE_PROFILE) | tail -n 1 | awk '{print $$3}')\` of statements" >> $(COVERAGE_BASELINE_FILE)
	@echo "" >> $(COVERAGE_BASELINE_FILE)
	@echo "## Detailed Coverage by File" >> $(COVERAGE_BASELINE_FILE)
	@echo "" >> $(COVERAGE_BASELINE_FILE)
	@go tool cover -func=$(COVERAGE_PROFILE) | grep -v "total:" | awk 'BEGIN{print "| File | Function | Coverage |"; print "|------|----------|----------|"} {print "| " $$1 " | " $$2 " | " $$3 " |"}' >> $(COVERAGE_BASELINE_FILE)
	@echo "✓ Baseline documentation updated: $(COVERAGE_BASELINE_FILE)"

# COV-002: Generate differential coverage report
test-coverage-differential:
	@echo "Generating differential coverage report (COV-002)..."
	@echo "Building differential coverage tool..."
	@mkdir -p $(COVERAGE_REPORTS_DIR)
	@cd tools && go build -o ../$(DIFFERENTIAL_TOOL) coverage-differential.go
	@echo "Running current coverage analysis..."
	go test -coverprofile=$(COVERAGE_PROFILE) . ./internal/testutil
	@echo "Generating differential report..."
	@./$(DIFFERENTIAL_TOOL)
	@echo "✓ Differential coverage analysis complete"

# COV-002: Update coverage trends and history
test-coverage-trends:
	@echo "Updating coverage trends and history (COV-002)..."
	@echo "Running coverage analysis..."
	go test -coverprofile=$(COVERAGE_PROFILE) . ./internal/testutil
	@echo "Building trend tracking tool..."
	@mkdir -p $(COVERAGE_REPORTS_DIR)
	@cd tools && go build -o ../$(DIFFERENTIAL_TOOL) coverage-differential.go
	@echo "Updating coverage history..."
	@./$(DIFFERENTIAL_TOOL)
	@if [ -f "$(COVERAGE_HISTORY_FILE)" ]; then \
		echo "✓ Coverage history updated: $(COVERAGE_HISTORY_FILE)"; \
	else \
		echo "⚠️  Coverage history file not found"; \
	fi

# COV-002: Full coverage suite including baseline, differential, and trends
test-coverage-full:
	@echo "Running full coverage suite (COV-002)..."
	$(MAKE) test-coverage-baseline
	$(MAKE) test-coverage-differential
	$(MAKE) test-coverage-trends
	@echo "✓ Full coverage analysis complete"

# COV-002: Quality gate validation with enhanced reporting
test-coverage-quality-gates:
	@echo "Validating quality gates with enhanced reporting (COV-002)..."
	@echo "Building differential coverage tool..."
	@mkdir -p $(COVERAGE_REPORTS_DIR)
	@cd tools && go build -o ../$(DIFFERENTIAL_TOOL) coverage-differential.go
	@echo "Running coverage analysis..."
	go test -coverprofile=$(COVERAGE_PROFILE) . ./internal/testutil
	@echo "Checking quality gates..."
	@./$(DIFFERENTIAL_TOOL)
	@echo "Quality gate validation complete"

test-race:
	@echo "Running tests with race detection..."
	go test -race ./...
	@echo "✓ Race detection tests passed"

test-bench:
	@echo "Running benchmark tests..."
	go test -bench=. ./...

test-all: test test-race test-coverage-new
	@echo "✓ All test variants completed"

# Code quality targets
fmt:
	@echo "Formatting code with gofmt..."
	gofmt -w .
	@echo "✓ Code formatting completed"

vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ go vet completed"

lint: validate-icon-enforcement
	@echo "Running linter (revive)..."
	@if command -v revive >/dev/null 2>&1; then \
		revive -config revive.toml -formatter friendly ./...; \
		echo "✓ Linting completed"; \
	else \
		echo "⚠️  revive not installed, skipping linting"; \
		echo "  Install with: go install github.com/mgechev/revive@latest"; \
	fi

# 🔺 DOC-007: Implementation token icon consistency validation
validate-icons:
	@echo "🔧 DOC-007: Validating implementation token icon consistency..."
	@if [ -f "scripts/validate-icon-consistency.sh" ]; then \
		chmod +x scripts/validate-icon-consistency.sh; \
		./scripts/validate-icon-consistency.sh; \
	else \
		echo "❌ Icon consistency validation script not found"; \
		echo "   Expected: scripts/validate-icon-consistency.sh"; \
		exit 1; \
	fi

# 🔺 DOC-008: Comprehensive icon validation and enforcement
validate-icon-enforcement:
	@echo "🛡️ DOC-008: Running comprehensive icon validation and enforcement..."
	@if [ -f "scripts/validate-icon-enforcement.sh" ]; then \
		chmod +x scripts/validate-icon-enforcement.sh; \
		./scripts/validate-icon-enforcement.sh; \
		echo "✓ DOC-008 validation completed - see icon-validation-report.md for details"; \
	else \
		echo "❌ DOC-008 icon enforcement script not found"; \
		echo "   Expected: scripts/validate-icon-enforcement.sh"; \
		exit 1; \
	fi

# 🔺 DOC-008: Strict mode validation for CI/CD pipelines
validate-icons-strict:
	@echo "🛡️ DOC-008: Running icon validation in strict mode (CI/CD)..."
	@if [ -f "scripts/validate-icon-enforcement.sh" ]; then \
		chmod +x scripts/validate-icon-enforcement.sh; \
		./scripts/validate-icon-enforcement.sh --strict; \
		echo "✓ DOC-008 strict validation passed"; \
	else \
		echo "❌ DOC-008 icon enforcement script not found"; \
		echo "   Expected: scripts/validate-icon-enforcement.sh"; \
		exit 1; \
	fi

# 🔺 DOC-009: Mass implementation token standardization targets
token-migration-dry-run:
	@echo "🔺 DOC-009: Running token migration in dry-run mode..."
	@if [ -f "scripts/token-migration.sh" ]; then \
		chmod +x scripts/token-migration.sh; \
		./scripts/token-migration.sh true; \
	else \
		echo "❌ Token migration script not found"; \
		echo "   Expected: scripts/token-migration.sh"; \
		exit 1; \
	fi

token-migration:
	@echo "🔺 DOC-009: Running mass token standardization..."
	@if [ -f "scripts/token-migration.sh" ]; then \
		chmod +x scripts/token-migration.sh; \
		./scripts/token-migration.sh false; \
	else \
		echo "❌ Token migration script not found"; \
		echo "   Expected: scripts/token-migration.sh"; \
		exit 1; \
	fi

token-migration-rollback:
	@echo "🔺 DOC-009: Rolling back token migration..."
	@if [ -f "scripts/token-migration.sh" ]; then \
		chmod +x scripts/token-migration.sh; \
		./scripts/token-migration.sh rollback; \
	else \
		echo "❌ Token migration script not found"; \
		echo "   Expected: scripts/token-migration.sh"; \
		exit 1; \
	fi

# 🔺 DOC-009: Priority icon analysis and mapping
analyze-priority-icons:
	@echo "🔺 DOC-009: Analyzing priority icon mappings..."
	@if [ -f "scripts/priority-icon-inference.sh" ]; then \
		chmod +x scripts/priority-icon-inference.sh; \
		./scripts/priority-icon-inference.sh mapping; \
		echo "✓ DOC-009 priority mapping generated"; \
	else \
		echo "❌ DOC-009 priority inference script not found"; \
		echo "   Expected: scripts/priority-icon-inference.sh"; \
		exit 1; \
	fi

# 🔺 DOC-009: Validate token priorities against feature tracking
validate-token-priorities:
	@echo "🔺 DOC-009: Validating token priorities against feature tracking..."
	@if [ -f "scripts/priority-icon-inference.sh" ]; then \
		chmod +x scripts/priority-icon-inference.sh; \
		./scripts/priority-icon-inference.sh validate; \
		echo "✓ DOC-009 priority validation completed"; \
	else \
		echo "❌ DOC-009 priority inference script not found"; \
		echo "   Expected: scripts/priority-icon-inference.sh"; \
		exit 1; \
	fi

# 🔺 DOC-009: Generate action icon suggestions
suggest-action-icons:
	@echo "🔺 DOC-009: Generating action icon suggestions..."
	@if [ -f "scripts/priority-icon-inference.sh" ]; then \
		chmod +x scripts/priority-icon-inference.sh; \
		./scripts/priority-icon-inference.sh suggest; \
		echo "✓ DOC-009 action icon suggestions generated"; \
	else \
		echo "❌ DOC-009 priority inference script not found"; \
		echo "   Expected: scripts/priority-icon-inference.sh"; \
		exit 1; \
	fi

# 🔺 DOC-009: Complete token standardization workflow
standardize-tokens: analyze-priority-icons validate-token-priorities migrate-tokens-dry-run
	@echo "🔺 DOC-009: Token standardization workflow completed"
	@echo "  Next steps:"
	@echo "  1. Review dry run output above"
	@echo "  2. Run 'make migrate-tokens' to execute actual migration"
	@echo "  3. Run 'make validate-icon-enforcement' to verify results"

check: fmt vet lint validate-icon-enforcement
	@echo "✓ All code quality checks completed (including DOC-008 icon validation)"

# Production build targets
build-all: build-local build-macos build-ubuntu
	@echo "✓ Built for all platforms"

build-ubuntu: build-ubuntu20 build-ubuntu22 build-ubuntu24
	@echo "✓ Built for all Ubuntu versions"

build-ubuntu20:
	@echo "Building for Ubuntu 20.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u20'" -o bin/$(BINARY_NAME)-ubuntu20.04

build-ubuntu22:
	@echo "Building for Ubuntu 22.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u22'" -o bin/$(BINARY_NAME)-ubuntu22.04

build-ubuntu24:
	@echo "Building for Ubuntu 24.04..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=linux-amd64-u24'" -o bin/$(BINARY_NAME)-ubuntu24.04

build-macos: build-macos-arm64 build-macos-amd64
	@echo "✓ Built for all macOS architectures"

build-macos-arm64:
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	@mkdir -p bin
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=darwin-arm64'" -o bin/$(BINARY_NAME)-macos-arm64

build-macos-amd64:
	@echo "Building for macOS AMD64 (Intel)..."
	@mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.compileDate=$(BUILD_TIME)' -X 'main.platform=darwin-amd64'" -o bin/$(BINARY_NAME)-macos-amd64

# Utility targets
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f $(BINARY_NAME)
	@rm -f $(COVERAGE_PROFILE) $(COVERAGE_NEW_PROFILE) $(COVERAGE_LEGACY_PROFILE)
	@rm -f coverage.html coverage_new.html coverage_legacy.html
	@rm -f coverage_report.txt
	@rm -f tools/coverage-analyzer
	@go clean -testcache
	@echo "✓ Clean completed"

deps:
	@echo "Downloading and verifying dependencies..."
	go mod download
	go mod verify
	@echo "✓ Dependencies verified"

# Legacy targets for compatibility
build: build-local

# COV-001: Coverage-focused targets for development workflow
coverage-check: test-coverage-validate
	@echo "✓ Coverage validation complete"

# Development workflow target that includes new coverage validation
dev-full: check test-coverage-validate build-local token-suggester
	@echo "🚀 Full development workflow completed with token suggestion integration"

# Example usage comment for symbolic link
# To create a symbolic link in your local bin:
# ln -s "$(pwd)/bin/$(BINARY_NAME)-macos-arm64" ~/.local/bin/$(BINARY_NAME)

# 🔶 DOC-012: Real-time icon validation feedback targets
build-realtime-validator:
	@echo "🔶 DOC-012: Building real-time validation service..."
	@mkdir -p bin
	go build -ldflags="$(LDFLAGS)" -o bin/realtime-validator cmd/realtime-validator/main.go
	@echo "✓ DOC-012: Real-time validator built successfully"

test-realtime-validation:
	@echo "🔶 DOC-012: Testing real-time validation system..."
	go test -v ./internal/validation/... -run TestRealTimeValidation
	@echo "✓ DOC-012: Real-time validation tests completed"

benchmark-realtime-validation:
	@echo "🔶 DOC-012: Benchmarking real-time validation performance..."
	go test -bench=BenchmarkRealTimeValidation -benchmem ./internal/validation/...
	@echo "✓ DOC-012: Performance benchmarks completed"

start-realtime-server:
	@echo "🔶 DOC-012: Starting real-time validation server..."
	@if [ -f "bin/realtime-validator" ]; then \
		./bin/realtime-validator server --port 8080; \
	else \
		echo "❌ Real-time validator not built. Run 'make build-realtime-validator' first"; \
		exit 1; \
	fi

validate-realtime-performance:
	@echo "🔶 DOC-012: Validating real-time performance targets..."
	@if [ -f "bin/realtime-validator" ]; then \
		echo "Testing sub-second validation performance..."; \
		time ./bin/realtime-validator validate main.go 2>&1 | grep -E "(Processing Time|real)" || true; \
		echo "✓ DOC-012: Performance validation completed"; \
	else \
		echo "❌ Real-time validator not built. Run 'make build-realtime-validator' first"; \
		exit 1; \
	fi

demo-realtime-validation:
	@echo "🔶 DOC-012: Running real-time validation demonstration..."
	@if [ -f "bin/realtime-validator" ]; then \
		echo "1. Testing file validation with intelligent suggestions:"; \
		./bin/realtime-validator validate main.go --format summary; \
		echo ""; \
		echo "2. Testing visual status indicators:"; \
		./bin/realtime-validator status main.go internal/validation/ai_validation.go; \
		echo ""; \
		echo "3. Displaying performance metrics:"; \
		./bin/realtime-validator metrics; \
		echo "✓ DOC-012: Real-time validation demonstration completed"; \
	else \
		echo "❌ Real-time validator not built. Run 'make build-realtime-validator' first"; \
		exit 1; \
	fi

# 🔶 DOC-012: Complete real-time validation setup
setup-realtime-validation: build-realtime-validator test-realtime-validation benchmark-realtime-validation
	@echo "🔶 DOC-012: Real-time validation system setup completed"
	@echo "  Available commands:"
	@echo "  - make start-realtime-server    # Start HTTP server on port 8080"
	@echo "  - make demo-realtime-validation # Run demonstration"
	@echo "  - make validate-realtime-performance # Test performance targets"
	@echo "  - bin/realtime-validator validate <files>  # Validate files"
	@echo "  - bin/realtime-validator status <files>    # Show status indicators"
	@echo "  - bin/realtime-validator watch <files>     # Watch files for changes"

# 🔶 DOC-010: Token suggestion engine - 🔧 Automated token format suggestion system
token-suggester:
	@echo "🔶 DOC-010: Building token suggestion engine..."
	@cd cmd/token-suggester && go build -o ../../bin/token-suggester .
	@echo "✓ DOC-010 token-suggester built successfully"

token-suggester-test:
	@echo "🔶 DOC-010: Running token suggestion engine tests..."
	@cd cmd/token-suggester && go test -v ./...
	@echo "✓ DOC-010 token suggestion tests completed"

token-suggester-benchmark:
	@echo "🔶 DOC-010: Running token suggestion benchmarks..."
	@cd cmd/token-suggester && go test -bench=. -benchmem
	@echo "✓ DOC-010 token suggestion benchmarks completed"

# 🔶 DOC-010: Token analysis commands - 💡 Suggestion generation
analyze-tokens:
	@echo "🔶 DOC-010: Analyzing codebase for token suggestions..."
	@if [ -f "bin/token-suggester" ]; then \
		./bin/token-suggester analyze . --verbose; \
	else \
		echo "❌ token-suggester not built. Run 'make token-suggester' first"; \
		exit 1; \
	fi

suggest-tokens-batch:
	@echo "🔶 DOC-010: Generating batch token suggestions..."
	@if [ -f "bin/token-suggester" ]; then \
		./bin/token-suggester batch-suggest . --json > token-suggestions.json; \
		echo "✓ DOC-010 batch suggestions saved to token-suggestions.json"; \
	else \
		echo "❌ token-suggester not built. Run 'make token-suggester' first"; \
		exit 1; \
	fi

validate-token-formats:
	@echo "🔶 DOC-010: Validating existing token formats..."
	@if [ -f "bin/token-suggester" ]; then \
		./bin/token-suggester validate-tokens . --verbose; \
	else \
		echo "❌ token-suggester not built. Run 'make token-suggester' first"; \
		exit 1; \
	fi

# 🔶 DOC-010: Comprehensive token workflow - 🚀 Complete suggestion pipeline
token-workflow: token-suggester analyze-tokens validate-token-formats
	@echo "🔶 DOC-010: Complete token suggestion workflow completed"

# 🔶 DOC-010: Development targets - 🔧 Development utilities
token-dev: token-suggester-test token-suggester analyze-tokens
	@echo "🔶 DOC-010: Development workflow completed"
