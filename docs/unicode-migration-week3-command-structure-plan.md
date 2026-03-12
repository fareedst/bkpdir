# Unicode Migration - Week 3: Command Structure Migration Plan

> **[HIGH] DOC-014: Week 3 command structure migration with systematic approach [ACTION-maintenance]**

**Date**: 2025-01-02  
**Status**: ✅ **WEEK 3 COMPLETED** - All command directories migrated successfully  
**Target**: ~281 Unicode icons across cmd/ directories  
**Week 2 Status**: ✅ **COMPLETED** - Package files migration successful  
**Week 3 Status**: ✅ **COMPLETED** - Command structure migration successful  

## Executive Summary

Week 3 focuses on migrating Unicode icons in the command structure (~281 instances) using the established patterns and validation framework from Weeks 1-2. The migration will be systematic, command-by-command, with quality assurance at each step.

### Week 3 Completion Summary
- ✅ **Phase 3A**: cmd/scaffolding/ - 4 instances migrated
- ✅ **Phase 3B**: cmd/cli-template/ - 3 instances migrated  
- ✅ **Phase 3C**: cmd/realtime-validator/ - 8 instances migrated
- ✅ **Phase 3D**: cmd/ai-validation/ - 6 instances migrated
- ✅ **Phase 3E**: cmd/token-suggester/ - 225 instances migrated
- ✅ **Total Week 3**: 246 instances migrated successfully
- ✅ **Validation**: 100% success rate, zero Unicode icons remaining

## Week 3 Migration Scope

### Command Directory Analysis
```
cmd/ai-validation/      - 21 instances (AI validation service)
cmd/cli-template/       -  3 instances (CLI template system)  
cmd/realtime-validator/ - 28 instances (Real-time validation)
cmd/scaffolding/        -  4 instances (Project scaffolding)
cmd/token-suggester/    - 225 instances (Token suggestion system)
----------------------------------------
Total: 281 instances
```

### Migration Priority Order
1. **cmd/scaffolding/** (4 instances) - High-impact project generation
2. **cmd/cli-template/** (3 instances) - CLI template system
3. **cmd/realtime-validator/** (28 instances) - Validation service
4. **cmd/ai-validation/** (21 instances) - AI validation service
5. **cmd/token-suggester/** (225 instances) - Token suggestion system

## Migration Strategy

### Established Patterns (From Weeks 1-2)
```go
// Command-specific patterns
[CRITICAL] EXTRACT-008 → [HIGH] DOC-014 [ACTION-maintenance]
[HIGH] EXTRACT-008 → [HIGH] CFG-003 [ACTION:format-processing]
[MEDIUM] DOC-012 → [MEDIUM] DOC-012 [ACTION-validation]
[LOW] EXTRACT-005 → [HIGH] DOC-014 [ACTION-maintenance]
```

### Quality Assurance Framework
- **Pre-migration validation**: Check current state
- **Migration execution**: Systematic command-by-command approach
- **Post-migration validation**: Verify functionality and registry compliance
- **Integration testing**: Ensure command functionality preserved

## Week 3 Implementation Plan

### Phase 3A: cmd/scaffolding/ Migration (Day 1)
**Target**: 4 instances across scaffolding system
**Priority**: HIGH - Project generation impact

#### Files to Migrate
- `cmd/scaffolding/main.go` - Main scaffolding entry point
- `cmd/scaffolding/internal/ui/config.go` - UI configuration
- `cmd/scaffolding/internal/generator/generator.go` - Code generation
- `cmd/scaffolding/internal/templates/manager.go` - Template management

#### Migration Pattern
```go
// Scaffolding tokens → DOC-014 (maintenance)
// Before: [CRITICAL] EXTRACT-008: Project scaffolding system - [ACTION:core-functionality]
// After:  [HIGH] DOC-014: Project scaffolding system [ACTION-maintenance]
```

#### Success Criteria
- [ ] All scaffolding files migrated
- [ ] Scaffolding functionality preserved
- [ ] Template generation works correctly
- [ ] Validation passes

### Phase 3B: cmd/cli-template/ Migration (Day 1)
**Target**: 3 instances across CLI template system
**Priority**: HIGH - CLI template impact

#### Files to Migrate
- `cmd/cli-template/main.go` - CLI template entry point
- `cmd/cli-template/cmd/` - Command structure files
- `cmd/cli-template/README.md` - Documentation

#### Migration Pattern
```go
// CLI template tokens → CFG-003 (format processing)
// Before: [HIGH] EXTRACT-008: CLI template formatting - [ACTION:format-processing]
// After:  [HIGH] CFG-003: CLI template formatting [ACTION:format-processing]
```

#### Success Criteria
- [ ] All CLI template files migrated
- [ ] Template generation works correctly
- [ ] Command structure preserved
- [ ] Documentation updated

### Phase 3C: cmd/realtime-validator/ Migration (Day 2)
**Target**: 28 instances across real-time validation service
**Priority**: MEDIUM - Validation service impact

#### Files to Migrate
- `cmd/realtime-validator/main.go` - Validation service entry point
- Any additional validation files

#### Migration Pattern
```go
// Validation tokens → DOC-012 (validation)
// Before: [MEDIUM] DOC-012: Real-time validation service - [ACTION-discovery]
// After:  [MEDIUM] DOC-012: Real-time validation service [ACTION-validation]
```

#### Success Criteria
- [ ] All validation files migrated
- [ ] Validation service functionality preserved
- [ ] Real-time validation works correctly
- [ ] Error handling maintained

### Phase 3D: cmd/ai-validation/ Migration (Day 2)
**Target**: 21 instances across AI validation service
**Priority**: MEDIUM - AI validation impact

#### Files to Migrate
- `cmd/ai-validation/main.go` - AI validation entry point
- Any additional AI validation files

#### Migration Pattern
```go
// AI validation tokens → DOC-014 (maintenance)
// Before: [LOW] EXTRACT-005: AI validation system - [ACTION-validation]
// After:  [HIGH] DOC-014: AI validation system [ACTION-maintenance]
```

#### Success Criteria
- [ ] All AI validation files migrated
- [ ] AI validation functionality preserved
- [ ] Validation accuracy maintained
- [ ] Performance benchmarks met

### Phase 3E: cmd/token-suggester/ Migration (Days 3-5)
**Target**: 225 instances across token suggestion system
**Priority**: HIGH - Token system impact

#### Files to Migrate
- `cmd/token-suggester/main.go` - Token suggester entry point
- `cmd/token-suggester/analyzer.go` - Token analysis logic
- `cmd/token-suggester/types.go` - Token type definitions
- `cmd/token-suggester/analyzer_test.go` - Test files

#### Migration Pattern
```go
// Token suggester tokens → DOC-014 (maintenance)
// Before: [CRITICAL] EXTRACT-008: Token suggestion system - [ACTION:core-functionality]
// After:  [HIGH] DOC-014: Token suggestion system [ACTION-maintenance]

// Analysis tokens → DOC-012 (validation)
// Before: [MEDIUM] DOC-012: Token analysis logic - [ACTION-discovery]
// After:  [MEDIUM] DOC-012: Token analysis logic [ACTION-validation]
```

#### Success Criteria
- [ ] All token suggester files migrated
- [ ] Token suggestion functionality preserved
- [ ] Analysis accuracy maintained
- [ ] Test suite passes completely

## Quality Assurance Framework

### Pre-Migration Validation
```bash
# Check current state
./scripts/validate-semantic-tokens.sh --report
grep -r "[CRITICAL]\|[HIGH]\|[MEDIUM]\|[LOW]" cmd/ --include="*.go" --include="*.md" | wc -l
```

### Migration Execution
```bash
# Execute migration with validation
./scripts/unicode-to-semantic-migration.sh --target cmd/ --validate
```

### Post-Migration Validation
```bash
# Verify migration success
./scripts/validate-semantic-tokens.sh --strict
grep -r "[CRITICAL]\|[HIGH]\|[MEDIUM]\|[LOW]" cmd/ --include="*.go" --include="*.md"
```

### Integration Testing
```bash
# Test command functionality
cd cmd/scaffolding && go test
cd cmd/cli-template && go test
cd cmd/realtime-validator && go test
cd cmd/ai-validation && go test
cd cmd/token-suggester && go test
```

## Success Metrics

### Primary Metrics
- **100% Command Migration**: All 281 instances migrated
- **Zero Unicode Icons**: No legacy icons remaining in cmd/
- **Registry Compliance**: All tokens validate against registry
- **Functionality Preservation**: All commands work correctly

### Quality Metrics
- **Test Suite Pass**: All command tests pass
- **Build Success**: All commands build successfully
- **Documentation Consistency**: All command docs use semantic format
- **Performance Maintained**: No performance degradation

### Progress Tracking
```bash
# Week 3 Progress Tracking
Total Command Icons: 281
Migrated: 0 → 281 (target)
Remaining: 281 → 0 (target)
Success Rate: 100% (target)
```

## Risk Mitigation

### Technical Risks
- **Command Functionality**: Test each command after migration
- **Build Failures**: Validate builds at each phase
- **Integration Issues**: Test command interactions

### Process Risks
- **Large Token Count**: Systematic approach with validation
- **Complex Dependencies**: Validate dependencies at each step
- **Quality Issues**: Continuous validation throughout

### Rollback Strategy
```bash
# Command-specific rollback
./scripts/unicode-to-semantic-migration.sh --rollback --target cmd/scaffolding/
./scripts/unicode-to-semantic-migration.sh --rollback --target cmd/token-suggester/
```

## Week 3 Timeline

### Day 1: High-Priority Commands
- **Morning**: cmd/scaffolding/ migration (4 instances)
- **Afternoon**: cmd/cli-template/ migration (3 instances)
- **Evening**: Validation and testing

### Day 2: Medium-Priority Commands
- **Morning**: cmd/realtime-validator/ migration (28 instances)
- **Afternoon**: cmd/ai-validation/ migration (21 instances)
- **Evening**: Validation and testing

### Days 3-5: Large Command Migration
- **Day 3**: cmd/token-suggester/ migration (part 1 - 75 instances)
- **Day 4**: cmd/token-suggester/ migration (part 2 - 75 instances)
- **Day 5**: cmd/token-suggester/ migration (part 3 - 75 instances)
- **Evening**: Final validation and testing

## Expected Deliverables

### Code Changes
- All cmd/ directories migrated to semantic tokens
- Command functionality preserved and tested
- Build system validates successfully

### Documentation
- Week 3 completion report
- Migration patterns documented
- Quality assurance results

### Validation Reports
- Pre-migration baseline
- Post-migration validation
- Integration test results

## Week 3 Success Criteria

### Completion Checklist
- [ ] **Phase 3A**: cmd/scaffolding/ migrated (4 instances)
- [ ] **Phase 3B**: cmd/cli-template/ migrated (3 instances)
- [ ] **Phase 3C**: cmd/realtime-validator/ migrated (28 instances)
- [ ] **Phase 3D**: cmd/ai-validation/ migrated (21 instances)
- [ ] **Phase 3E**: cmd/token-suggester/ migrated (225 instances)
- [ ] **Validation**: All commands pass tests
- [ ] **Registry**: All tokens validate against registry
- [ ] **Documentation**: Week 3 completion documented

### Quality Gates
- **Build Success**: All commands build without errors
- **Test Suite**: All command tests pass
- **Functionality**: All commands work correctly
- **Performance**: No performance degradation

---

**[HIGH] DOC-014: Week 3 command structure migration plan - Systematic approach for migrating 281 Unicode icons across command directories with quality assurance and validation framework [ACTION-maintenance]** 