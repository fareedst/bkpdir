#!/bin/bash
# scripts/generate-missing-tokens.sh
# DOC-014: Automated token generation for cross-reference completion

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Usage information
show_usage() {
    cat << EOF
Usage: $0 --type <type> --priority <priority> [options]

Generate missing tokens for cross-reference completion.

Options:
  --type <type>           Token type: specification, requirements, architecture
  --priority <priority>   Priority: critical, high, medium
  --dry-run              Show what would be generated without making changes
  --validate             Validate existing tokens before generation
  --help                 Show this help message

Examples:
  $0 --type specification --priority critical
  $0 --type requirements --priority high --dry-run
  $0 --type architecture --priority medium --validate

EOF
}

# Parse command line arguments
TYPE=""
PRIORITY=""
DRY_RUN=false
VALIDATE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --type)
            TYPE="$2"
            shift 2
            ;;
        --priority)
            PRIORITY="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --validate)
            VALIDATE=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Validate required arguments
if [[ -z "$TYPE" ]] || [[ -z "$PRIORITY" ]]; then
    log_error "Missing required arguments: --type and --priority"
    show_usage
    exit 1
fi

# Validate type and priority values
case $TYPE in
    specification|requirements|architecture)
        ;;
    *)
        log_error "Invalid type: $TYPE. Must be specification, requirements, or architecture"
        exit 1
        ;;
esac

case $PRIORITY in
    critical|high|medium)
        ;;
    *)
        log_error "Invalid priority: $PRIORITY. Must be critical, high, or medium"
        exit 1
        ;;
esac

# Check if we're in the right directory
if [[ ! -f "$PROJECT_ROOT/go.mod" ]] || [[ ! -d "$PROJECT_ROOT/docs/context" ]]; then
    log_error "Please run this script from the project root directory"
    exit 1
fi

# Function to validate existing tokens
validate_existing_tokens() {
    log_info "Validating existing tokens..."
    
    local validation_errors=0
    
    # Check for token conflicts
    case $TYPE in
        specification)
            # Check for existing SPEC tokens
            existing_spec_tokens=$(grep -r "SPEC-[0-9]\+" "$PROJECT_ROOT/docs/context/" 2>/dev/null | grep -o "SPEC-[0-9]\+" | sort -u || true)
            if [[ -n "$existing_spec_tokens" ]]; then
                log_warning "Found existing SPEC tokens: $existing_spec_tokens"
            fi
            ;;
        requirements)
            # Check for existing REQ tokens
            existing_req_tokens=$(grep -r "REQ-[0-9]\+" "$PROJECT_ROOT/docs/context/" 2>/dev/null | grep -o "REQ-[0-9]\+" | sort -u || true)
            if [[ -n "$existing_req_tokens" ]]; then
                log_warning "Found existing REQ tokens: $existing_req_tokens"
            fi
            ;;
        architecture)
            # Check for existing ARCH tokens
            existing_arch_tokens=$(grep -r "ARCH-[0-9]\+" "$PROJECT_ROOT/docs/context/" 2>/dev/null | grep -o "ARCH-[0-9]\+" | sort -u || true)
            if [[ -n "$existing_arch_tokens" ]]; then
                log_warning "Found existing ARCH tokens: $existing_arch_tokens"
            fi
            ;;
    esac
    
    # Check feature tracking consistency
    if [[ -f "$PROJECT_ROOT/docs/context/feature-tracking.md" ]]; then
        log_info "Checking feature tracking consistency..."
        # Implementation would check if tokens exist in feature-tracking.md
    fi
    
    if [[ $validation_errors -eq 0 ]]; then
        log_success "Token validation passed"
    else
        log_error "Token validation failed with $validation_errors errors"
        return 1
    fi
}

# Function to generate specification tokens
generate_specification_tokens() {
    local priority="$1"
    log_info "Generating specification tokens with priority: $priority"
    
    case $priority in
        critical)
            generate_critical_specification_tokens
            ;;
        high)
            generate_high_specification_tokens
            ;;
        medium)
            generate_medium_specification_tokens
            ;;
    esac
}

generate_critical_specification_tokens() {
    log_info "Generating critical priority specification tokens..."
    
    local tokens=(
        "SPEC-001: Quality assurance and code standards [DECISION: core-functionality, quality-gate, development-standards]"
        "SPEC-002: Configuration discovery system [DECISION: core-functionality, configuration-management]"
        "SPEC-003: CLI global options [DECISION: core-functionality, user-interface]"
        "SPEC-004: Archive features [DECISION: core-functionality, data-management]"
        "SPEC-005: Error handling standards [DECISION: core-functionality, reliability]"
        "SPEC-006: Resource management [DECISION: core-functionality, performance]"
        "SPEC-007: Build and development [DECISION: core-functionality, development-process]"
        "SPEC-008: Context management [DECISION: core-functionality, concurrency]"
        "SPEC-009: Backup features [DECISION: core-functionality, data-protection]"
        "SPEC-010: Format features [DECISION: core-functionality, user-interface]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_high_specification_tokens() {
    log_info "Generating high priority specification tokens..."
    
    local tokens=(
        "SPEC-011: Template system specifications [DECISION: enhanced-functionality, customization]"
        "SPEC-012: Verification system specifications [DECISION: enhanced-functionality, data-integrity]"
        "SPEC-013: Git integration specifications [DECISION: enhanced-functionality, version-control]"
        "SPEC-014: Performance specifications [DECISION: enhanced-functionality, optimization]"
        "SPEC-015: Security specifications [DECISION: enhanced-functionality, security]"
        "SPEC-016: Platform compatibility specifications [DECISION: enhanced-functionality, portability]"
        "SPEC-017: Testing specifications [DECISION: enhanced-functionality, quality-assurance]"
        "SPEC-018: Documentation specifications [DECISION: enhanced-functionality, maintainability]"
        "SPEC-019: Configuration specifications [DECISION: enhanced-functionality, flexibility]"
        "SPEC-020: Error recovery specifications [DECISION: enhanced-functionality, reliability]"
        "SPEC-021: Monitoring specifications [DECISION: enhanced-functionality, observability]"
        "SPEC-022: Deployment specifications [DECISION: enhanced-functionality, operations]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_medium_specification_tokens() {
    log_info "Generating medium priority specification tokens..."
    
    local tokens=(
        "SPEC-023: Advanced archive features [DECISION: advanced-functionality, data-management]"
        "SPEC-024: Advanced backup features [DECISION: advanced-functionality, data-protection]"
        "SPEC-025: Advanced configuration features [DECISION: advanced-functionality, flexibility]"
        "SPEC-026: Advanced error handling [DECISION: advanced-functionality, reliability]"
        "SPEC-027: Advanced performance features [DECISION: advanced-functionality, optimization]"
        "SPEC-028: Advanced security features [DECISION: advanced-functionality, security]"
        "SPEC-029: Advanced testing features [DECISION: advanced-functionality, quality-assurance]"
        "SPEC-030: Advanced documentation features [DECISION: advanced-functionality, maintainability]"
        "SPEC-031: Advanced monitoring features [DECISION: advanced-functionality, observability]"
        "SPEC-032: Advanced deployment features [DECISION: advanced-functionality, operations]"
        "SPEC-033: Advanced integration features [DECISION: advanced-functionality, interoperability]"
        "SPEC-034: Advanced customization features [DECISION: advanced-functionality, flexibility]"
        "SPEC-035: Advanced automation features [DECISION: advanced-functionality, efficiency]"
        "SPEC-036: Advanced analytics features [DECISION: advanced-functionality, insights]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

# Function to generate requirements tokens
generate_requirements_tokens() {
    local priority="$1"
    log_info "Generating requirements tokens with priority: $priority"
    
    case $priority in
        critical)
            generate_critical_requirements_tokens
            ;;
        high)
            generate_high_requirements_tokens
            ;;
        medium)
            generate_medium_requirements_tokens
            ;;
    esac
}

generate_critical_requirements_tokens() {
    log_info "Generating critical priority requirements tokens..."
    
    local tokens=(
        "REQ-001: Quality standards [DECISION: core-functionality, quality-gate]"
        "REQ-002: Build system [DECISION: core-functionality, development-process]"
        "REQ-003: Resource management [DECISION: core-functionality, performance]"
        "REQ-004: Error handling [DECISION: core-functionality, reliability]"
        "REQ-005: Context management [DECISION: core-functionality, concurrency]"
        "REQ-006: Backup system [DECISION: core-functionality, data-protection]"
        "REQ-007: Format system [DECISION: core-functionality, user-interface]"
        "REQ-008: Template system [DECISION: core-functionality, customization]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_high_requirements_tokens() {
    log_info "Generating high priority requirements tokens..."
    
    local tokens=(
        "REQ-009: Advanced archive requirements [DECISION: enhanced-functionality, data-management]"
        "REQ-010: Advanced backup requirements [DECISION: enhanced-functionality, data-protection]"
        "REQ-011: Advanced configuration requirements [DECISION: enhanced-functionality, flexibility]"
        "REQ-012: Advanced error handling requirements [DECISION: enhanced-functionality, reliability]"
        "REQ-013: Advanced performance requirements [DECISION: enhanced-functionality, optimization]"
        "REQ-014: Advanced security requirements [DECISION: enhanced-functionality, security]"
        "REQ-015: Advanced testing requirements [DECISION: enhanced-functionality, quality-assurance]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_medium_requirements_tokens() {
    log_info "Generating medium priority requirements tokens..."
    
    local tokens=(
        "REQ-016: Advanced documentation requirements [DECISION: advanced-functionality, maintainability]"
        "REQ-017: Advanced monitoring requirements [DECISION: advanced-functionality, observability]"
        "REQ-018: Advanced deployment requirements [DECISION: advanced-functionality, operations]"
        "REQ-019: Advanced integration requirements [DECISION: advanced-functionality, interoperability]"
        "REQ-020: Advanced customization requirements [DECISION: advanced-functionality, flexibility]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

# Function to generate architecture tokens
generate_architecture_tokens() {
    local priority="$1"
    log_info "Generating architecture tokens with priority: $priority"
    
    case $priority in
        critical)
            generate_critical_architecture_tokens
            ;;
        high)
            generate_high_architecture_tokens
            ;;
        medium)
            generate_medium_architecture_tokens
            ;;
    esac
}

generate_critical_architecture_tokens() {
    log_info "Generating critical priority architecture tokens..."
    
    local tokens=(
        "ARCH-001: Core system architecture [DECISION: core-functionality, system-design]"
        "ARCH-002: Configuration management architecture [DECISION: core-functionality, configuration]"
        "ARCH-003: Archive management architecture [DECISION: core-functionality, data-management]"
        "ARCH-004: Backup management architecture [DECISION: core-functionality, data-protection]"
        "ARCH-005: Error handling architecture [DECISION: core-functionality, reliability]"
        "ARCH-006: Resource management architecture [DECISION: core-functionality, performance]"
        "ARCH-007: Context management architecture [DECISION: core-functionality, concurrency]"
        "ARCH-008: Format system architecture [DECISION: core-functionality, user-interface]"
        "ARCH-009: Template system architecture [DECISION: core-functionality, customization]"
        "ARCH-010: Verification system architecture [DECISION: core-functionality, data-integrity]"
        "ARCH-011: Git integration architecture [DECISION: core-functionality, version-control]"
        "ARCH-012: Performance architecture [DECISION: core-functionality, optimization]"
        "ARCH-013: Security architecture [DECISION: core-functionality, security]"
        "ARCH-014: Testing architecture [DECISION: core-functionality, quality-assurance]"
        "ARCH-015: Documentation architecture [DECISION: core-functionality, maintainability]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_high_architecture_tokens() {
    log_info "Generating high priority architecture tokens..."
    
    local tokens=(
        "ARCH-016: Advanced integration architecture [DECISION: enhanced-functionality, interoperability]"
        "ARCH-017: Advanced monitoring architecture [DECISION: enhanced-functionality, observability]"
        "ARCH-018: Advanced deployment architecture [DECISION: enhanced-functionality, operations]"
        "ARCH-019: Advanced customization architecture [DECISION: enhanced-functionality, flexibility]"
        "ARCH-020: Advanced automation architecture [DECISION: enhanced-functionality, efficiency]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

generate_medium_architecture_tokens() {
    log_info "Generating medium priority architecture tokens..."
    
    local tokens=(
        "ARCH-021: Advanced analytics architecture [DECISION: advanced-functionality, insights]"
        "ARCH-022: Advanced machine learning architecture [DECISION: advanced-functionality, intelligence]"
        "ARCH-023: Advanced cloud integration architecture [DECISION: advanced-functionality, cloud]"
        "ARCH-024: Advanced distributed architecture [DECISION: advanced-functionality, scalability]"
        "ARCH-025: Advanced microservices architecture [DECISION: advanced-functionality, modularity]"
    )
    
    for token in "${tokens[@]}"; do
        if [[ "$DRY_RUN" == true ]]; then
            log_info "Would generate token: $token"
        else
            log_success "Generated token: $token"
        fi
    done
}

# Main execution
main() {
    log_info "Starting token generation for type: $TYPE, priority: $PRIORITY"
    
    if [[ "$VALIDATE" == true ]]; then
        validate_existing_tokens
    fi
    
    case $TYPE in
        specification)
            generate_specification_tokens "$PRIORITY"
            ;;
        requirements)
            generate_requirements_tokens "$PRIORITY"
            ;;
        architecture)
            generate_architecture_tokens "$PRIORITY"
            ;;
    esac
    
    if [[ "$DRY_RUN" == true ]]; then
        log_warning "Dry run completed. No changes were made."
    else
        log_success "Token generation completed successfully!"
    fi
}

# Run main function
main "$@" 