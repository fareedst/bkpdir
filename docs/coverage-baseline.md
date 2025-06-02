# Coverage Baseline Documentation

> **COV-001**: Baseline coverage levels established before implementing selective coverage reporting

## Baseline Information

- **Date Established**: 2025-06-02
- **Go Version**: go1.24.3 darwin/arm64
- **Measurement Method**: `go test -cover ./...`

## Current Coverage Levels

### Main Package Coverage: 73.5%

**Test Command**: `go test -cover .`
**Result**: `coverage: 73.5% of statements`

### Internal/Testutil Package Coverage: 75.6%

**Test Command**: `go test -cover ./internal/testutil`
**Result**: `coverage: 75.6% of statements`

## Coverage Analysis by File

### Files with High Coverage (>80%)
- Most test utility functions in `internal/testutil/`
- Core error handling in `errors.go`
- Configuration management functions
- Archive verification utilities

### Files with Medium Coverage (60-80%)
- Main application logic in `main.go`
- Archive creation functions in `archive.go`
- Backup management in `backup.go`

### Files with Lower Coverage (<60%)
- Complex formatting logic in `formatter.go`
- Git integration in `git.go`
- Configuration edge cases in `config.go`

## Testing Infrastructure Status

### Completed Testing Infrastructure
- **Archive Corruption Testing**: ✅ Complete (75.6% coverage)
- **Disk Space Simulation**: ✅ Complete 
- **Permission Testing**: ✅ Complete
- **Context Cancellation**: ✅ Complete
- **Error Injection Framework**: ✅ Complete

### Test Coverage Improvements Made
- **errors.go**: Achieved 100% coverage for most functions (87.5-100%)
- **comparison.go**: Achieved excellent coverage (88-100% for all functions)
- **config.go**: Successfully achieved 100% coverage for all previously untested functions
- **main.go**: Comprehensive test suite covering 30+ previously 0% coverage functions

## Exclusion Strategy

### Files Marked as Legacy (Stable)
These files have stable functionality and should be excluded from new coverage requirements:

1. **main.go** - CLI interface, well-tested through integration
2. **config.go** - Configuration management, comprehensive existing tests
3. **formatter.go** - Output formatting, complex template logic
4. **backup.go** - Backup operations, stable functionality
5. **archive.go** - Archive creation, core stable functions

### New Development Focus Areas
Coverage requirements will apply strictly to:

- New features added after 2025-06-02
- Bug fixes that modify existing logic
- Performance improvements
- API changes or extensions

## Coverage Quality Gates

### For New Code
- **Minimum Coverage**: 85%
- **Test Types Required**: Unit tests, integration tests
- **Edge Case Coverage**: Required for error paths

### For Modified Legacy Code
- **Minimum Coverage**: 70%
- **Preserve Existing Tests**: All existing tests must continue to pass
- **Document Changes**: Coverage impact must be documented

## Implementation Notes

### Build Tags Strategy
- Use `//go:build legacy` for legacy code exclusion
- Use `//go:build exclude_coverage` for specific function exclusion
- Maintain test execution while hiding from coverage metrics

### Coverage Tooling
- `go test -coverpkg` for selective package coverage
- HTML reports with exclusion highlighting
- Differential coverage for pull requests

## Next Steps

1. ✅ **Baseline Established**: Current coverage levels documented
2. 🔄 **Implement Exclusion**: Modify Makefile coverage targets
3. ⏳ **Validate System**: Test selective coverage reporting
4. ⏳ **Document Standards**: Update development documentation 