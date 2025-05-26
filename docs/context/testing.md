# Testing Requirements and Architecture

This document outlines the testing requirements, architecture, and approach for the BkpDir application.

## Test Categories

### Unit Tests
**Implementation**: `*_test.go` files
**Requirements**:
- `TestGenerateArchiveName`: Tests archive naming
  - Validates naming format with various inputs including Git information
  - Test cases:
    - Basic archive naming without Git info
    - Archive naming with Git branch and hash
    - Archive naming with notes
    - Archive naming with prefix
    - Incremental archive naming
    - Various timestamp formats
- `TestDefaultConfig` and `TestLoadConfig`: Tests configuration
  - Validates configuration loading and defaults
  - Tests configuration discovery with `BKPDIR_CONFIG` environment variable
  - Tests multiple configuration file precedence
  - Tests home directory expansion in paths
  - Tests status code configuration loading and defaults
  - Validates all status code fields are properly loaded from YAML
  - Tests status code configuration with custom values
  - Tests exclude patterns configuration
  - Tests verification configuration
  - Tests Git integration configuration
- `TestGetConfigSearchPath`: Tests configuration path discovery
  - Validates environment variable parsing
  - Tests hard-coded default path when environment variable not set
  - Tests colon-separated path list parsing
  - Tests home directory expansion
- `TestDisplayConfig`: Tests configuration value display
  - Validates configuration value computation and source tracking
  - Tests environment variable processing for configuration paths
  - Tests hard-coded default path handling
  - Tests default value handling and source attribution
  - Tests output format with name, value, and source
  - Tests display of status code configuration values
  - Test cases:
    - Configuration with default values only
    - Configuration from single file
    - Configuration from multiple files with precedence
    - Configuration with missing files
    - Configuration with invalid files
    - Environment variable override scenarios
    - Status code configuration display
    - Exclude patterns configuration display
    - Verification configuration display
- `TestShouldExcludeFile`: Tests file exclusion patterns
  - Validates doublestar glob pattern matching
  - Test cases:
    - Basic glob patterns
    - Directory exclusions
    - File extension exclusions
    - Complex nested patterns
    - Case sensitivity handling
- `TestCopyFileWithContext`: Tests context-aware file copying
  - Validates file copying with cancellation support
  - Test cases:
    - Successful copy with context
    - Copy cancelled via context
    - Copy with timeout
    - Copy with already cancelled context
- `TestCompareDirectories`: Tests directory comparison
  - Validates directory tree comparison functionality
  - Test cases:
    - Identical directories
    - Different directories
    - Directories with different file sizes
    - Empty directories
    - Large directories
    - Directories with special characters
    - Directories with excluded files
- `TestListArchives`: Tests archive listing
  - Validates archive listing and sorting
  - Tests verification status display
  - Test cases:
    - Empty archive directory
    - Multiple archives with different timestamps
    - Archives with Git information
    - Archives with notes
    - Archives with verification status
- `TestCreateFullArchive`: Tests full archive creation
  - Validates full archive creation with various scenarios
  - Tests status code exit behavior for different conditions
  - Tests panic recovery and error handling
  - Test cases:
    - Successful archive creation (should exit with `cfg.StatusCreatedArchive`)
    - Directory identical to existing archive (should exit with `cfg.StatusDirectoryIsIdenticalToExistingArchive`)
    - Directory not found (should exit with `cfg.StatusDirectoryNotFound`)
    - Invalid directory type (should exit with `cfg.StatusInvalidDirectoryType`)
    - Permission denied (should exit with `cfg.StatusPermissionDenied`)
    - Disk full scenarios (should exit with `cfg.StatusDiskFull`)
    - Archive directory creation failure (should exit with `cfg.StatusFailedToCreateArchiveDirectory`)
    - Configuration errors (should exit with `cfg.StatusConfigError`)
    - Panic recovery scenarios
- `TestCreateFullArchiveWithCleanup`: Tests archive creation with resource cleanup
  - Validates automatic resource cleanup functionality
  - Tests atomic operations with temporary files
  - Test cases:
    - Successful archive with cleanup verification
    - Archive failure with cleanup verification
    - No temporary files left after operations
    - Atomic file operations
- `TestCreateFullArchiveWithContext`: Tests context-aware archive creation
  - Validates archive creation with cancellation support
  - Test cases:
    - Successful archive with context
    - Archive cancelled via context
    - Archive with timeout
    - Context cancellation at various stages
- `TestCreateFullArchiveWithContextAndCleanup`: Tests context-aware archive with cleanup
  - Validates most robust archive creation functionality
  - Combines context support with resource cleanup
  - Test cases:
    - Successful archive with context and cleanup
    - Cancelled archive with proper cleanup
    - Timeout scenarios with cleanup verification
    - No resource leaks on cancellation
- `TestCreateIncrementalArchive`: Tests incremental archive creation
  - Validates incremental archive functionality
  - Test cases:
    - Successful incremental archive creation
    - No base archive available
    - No changes since base archive
    - Multiple incremental archives
    - Incremental archive with Git changes
- `TestResourceManager`: Tests resource management functionality
  - Validates thread-safe resource tracking
  - Tests automatic cleanup mechanisms
  - Test cases:
    - Basic resource registration and cleanup
    - Thread-safe concurrent access
    - Cleanup of both files and directories
    - Error-resilient cleanup (continues on individual failures)
    - Cleanup warnings logged to stderr
- `TestArchiveError`: Tests structured error handling
  - Validates ArchiveError functionality
  - Test cases:
    - Error creation with message and status code
    - Error interface implementation
    - Status code extraction
    - Error message formatting
- `TestIsDiskFullError`: Tests enhanced disk space detection
  - Validates disk full error detection
  - Test cases:
    - Various disk full error messages
    - Case-insensitive matching
    - Multiple disk space indicators
    - Non-disk-full errors (should return false)
    - Nil error handling
- `TestConfigurationDiscovery`: Tests configuration file discovery
  - Tests multiple configuration files with different precedence
  - Tests environment variable override behavior
  - Tests hard-coded default path behavior
  - Tests missing configuration files handling
  - Tests invalid configuration file handling
  - Tests configuration merging with defaults
  - Tests status code configuration precedence and merging
- `TestStatusCodeConfiguration`: Tests status code configuration
  - Validates status code loading from YAML
  - Tests status code defaults
  - Tests status code precedence with multiple configuration files
  - Tests invalid status code values handling
  - Test cases:
    - Default status codes (should match immutable specification defaults)
    - Custom status codes from configuration file
    - Status code precedence with multiple files
    - Invalid status code values (non-integer)
    - Missing status code fields (should use defaults)

### Git Integration Tests
**Implementation**: `git_test.go`
**Requirements**:
- `TestIsGitRepository`: Tests Git repository detection
  - Test cases:
    - Valid Git repository
    - Non-Git directory
    - Directory with .git file (submodule)
    - Permission denied scenarios
- `TestGetGitBranch`: Tests Git branch extraction
  - Test cases:
    - Standard branch names
    - Branch names with special characters
    - Detached HEAD state
    - Non-Git repository
- `TestGetGitShortHash`: Tests Git commit hash extraction
  - Test cases:
    - Valid commit hash
    - Initial commit
    - Dirty working directory
    - Non-Git repository
- `TestGitIntegrationInArchiveNaming`: Tests Git info in archive names
  - Test cases:
    - Archive naming with Git info enabled
    - Archive naming with Git info disabled
    - Archive naming in non-Git directory
    - Archive naming with special branch names

### Archive Verification Tests
**Implementation**: `verify_test.go`
**Requirements**:
- `TestVerifyArchive`: Tests archive structure verification
  - Validates ZIP file structure verification
  - Test cases:
    - Valid ZIP archive
    - Corrupted ZIP archive
    - Non-ZIP file
    - Empty file
    - Archive with special characters
- `TestVerifyChecksums`: Tests checksum verification
  - Validates file content verification against stored checksums
  - Test cases:
    - Valid checksums
    - Corrupted file content
    - Missing checksum file
    - Invalid checksum format
    - Checksum algorithm mismatch
- `TestGenerateChecksums`: Tests checksum generation
  - Validates checksum generation for files
  - Test cases:
    - Single file checksum
    - Multiple file checksums
    - Large file checksums
    - Empty file checksums
    - Files with special characters
    - Directory exclusion from checksums
- `TestStoreAndReadChecksums`: Tests checksum storage and retrieval
  - Validates checksum persistence in archives
  - Test cases:
    - Store checksums in archive
    - Read checksums from archive
    - Invalid checksum data
    - Missing checksum file
- `TestVerificationStatus`: Tests verification status tracking
  - Validates verification result storage and retrieval
  - Test cases:
    - Successful verification status
    - Failed verification status
    - Verification status persistence
    - Multiple verification attempts

### Output Formatting Tests
**Implementation**: `formatter_test.go`
**Requirements**:
- `TestOutputFormatter`: Tests output formatter functionality
  - Validates printf-style formatting with various inputs
  - Tests ANSI color code support
  - Test cases:
    - Basic format string application
    - Format strings with ANSI color codes
    - Format strings with special characters
    - Invalid format strings (should fall back to safe defaults)
    - Empty format strings
    - Format strings with multiple parameters
- `TestFormatCreatedArchive`: Tests archive creation message formatting
  - Validates format string application for successful archive messages
  - Test cases:
    - Default format string
    - Custom format string with colors
    - Format string with special characters
    - Various path lengths and characters
- `TestFormatIdenticalArchive`: Tests identical directory message formatting
  - Validates format string application for identical directory messages
  - Test cases:
    - Default format string
    - Custom format string with highlighting
    - Format string with symbols and colors
    - Various path formats
- `TestFormatListArchive`: Tests archive listing entry formatting
  - Validates format string application for archive list entries
  - Test cases:
    - Default format string with two parameters
    - Custom format string with color coding
    - Format string with timestamp formatting
    - Various path and time combinations
- `TestFormatConfigValue`: Tests configuration value display formatting
  - Validates format string application for configuration display
  - Test cases:
    - Default format string with three parameters
    - Custom format string with highlighting
    - Format string with source emphasis
    - Various configuration name/value combinations
- `TestFormatDryRunArchive`: Tests dry-run message formatting
  - Validates format string application for dry-run messages
  - Test cases:
    - Default format string
    - Custom format string with warning indicators
    - Format string with visual emphasis
    - Various path formats
- `TestFormatError`: Tests error message formatting
  - Validates format string application for error messages
  - Test cases:
    - Default format string
    - Custom format string with error highlighting
    - Format string with symbols and colors
    - Various error message types
- `TestTemplateFormatter`: Tests template formatter functionality
  - Validates template-based formatting with named placeholders
  - Tests regex pattern extraction and template application
  - Test cases:
    - Basic template string application
    - Template strings with named placeholders
    - Regex pattern extraction with named groups
    - Template strings with Go text/template syntax
    - Invalid template strings (should fall back to safe defaults)
    - Template strings with conditional formatting
- `TestTemplateCreatedArchive`: Tests template-based archive creation formatting
  - Validates template formatting for archive creation messages
  - Test cases:
    - Template with named placeholders
    - Template with regex pattern extraction
    - Template with conditional formatting
    - Template with archive metadata
- `TestTemplateListArchive`: Tests template-based archive listing formatting
  - Validates template formatting for archive list entries
  - Test cases:
    - Template with archive filename parsing
    - Template with Git information extraction
    - Template with timestamp formatting
    - Template with note display
- `TestFormatWithTemplate`: Tests template application with regex extraction
  - Validates template formatting with named regex groups
  - Test cases:
    - Archive filename pattern extraction
    - Configuration line pattern extraction
    - Timestamp pattern extraction
    - Invalid patterns and templates
- `TestFormatWithPlaceholders`: Tests placeholder replacement
  - Validates %{name} placeholder substitution
  - Test cases:
    - Basic placeholder replacement
    - Multiple placeholders
    - Missing placeholders
    - Special characters in placeholders

### Integration Tests
**Implementation**: `integration_test.go`
**Requirements**:
- `TestFullWorkflow`: Tests complete archive workflow
  - Validates end-to-end archive creation and listing
  - Test cases:
    - Create full archive and list
    - Create incremental archive and list
    - Verify archive after creation
    - Multiple archives with different configurations
- `TestDryRunMode`: Tests dry-run functionality
  - Validates dry-run behavior across all commands
  - Test cases:
    - Dry-run full archive creation
    - Dry-run incremental archive creation
    - No files created in dry-run mode
    - Proper output formatting in dry-run
- `TestCLICommands`: Tests command-line interface
  - Validates all CLI commands and flags
  - Test cases:
    - `bkpdir full` command
    - `bkpdir inc` command
    - `bkpdir list` command
    - `bkpdir verify` command
    - `bkpdir --config` command
    - `--dry-run` flag
    - `--checksum` flag for verify
- `TestConfigurationIntegration`: Tests configuration in real scenarios
  - Validates configuration loading and application
  - Test cases:
    - Default configuration behavior
    - Custom configuration files
    - Environment variable configuration
    - Configuration precedence
    - Invalid configuration handling
- `TestGitIntegrationWorkflow`: Tests Git integration in workflows
  - Validates Git information in archive creation
  - Test cases:
    - Archive creation in Git repository
    - Archive creation in non-Git directory
    - Archive creation with different branches
    - Archive creation with dirty working directory
- `TestVerificationWorkflow`: Tests verification in complete workflows
  - Validates verification integration with archive creation
  - Test cases:
    - Automatic verification after creation
    - Manual verification of existing archives
    - Verification with checksums
    - Verification failure scenarios
- `TestErrorHandlingIntegration`: Tests error handling in real scenarios
  - Validates error handling across all operations
  - Test cases:
    - Permission denied scenarios
    - Disk full scenarios
    - Invalid input scenarios
    - Network interruption scenarios (if applicable)
- `TestConcurrentOperations`: Tests concurrent archive operations
  - Validates thread safety and resource management
  - Test cases:
    - Multiple archive creations
    - Concurrent verification operations
    - Resource cleanup under concurrency
    - No race conditions

### Performance Tests
**Implementation**: `performance_test.go`
**Requirements**:
- `BenchmarkArchiveCreation`: Benchmarks archive creation performance
  - Measures performance with various directory sizes
  - Test cases:
    - Small directories (< 100 files)
    - Medium directories (100-1000 files)
    - Large directories (> 1000 files)
    - Deep directory hierarchies
- `BenchmarkDirectoryComparison`: Benchmarks directory comparison performance
  - Measures comparison speed for various scenarios
  - Test cases:
    - Identical directories
    - Completely different directories
    - Directories with few changes
    - Large directories
- `BenchmarkVerification`: Benchmarks verification performance
  - Measures verification speed for various archive sizes
  - Test cases:
    - Structure verification only
    - Checksum verification
    - Large archives
    - Many small files vs few large files
- `BenchmarkGitOperations`: Benchmarks Git integration performance
  - Measures Git information extraction speed
  - Test cases:
    - Git repository detection
    - Branch name extraction
    - Commit hash extraction
    - Large Git repositories
- `TestMemoryUsage`: Tests memory consumption
  - Validates memory usage stays within reasonable bounds
  - Test cases:
    - Large directory archiving
    - Long-running operations
    - Multiple concurrent operations
    - Memory leak detection

### Stress Tests
**Implementation**: `stress_test.go`
**Requirements**:
- `TestLargeDirectories`: Tests with very large directories
  - Validates handling of directories with many files
  - Test cases:
    - Directories with > 10,000 files
    - Deep directory hierarchies (> 20 levels)
    - Files with very long names
    - Directories with special characters
- `TestLongRunningOperations`: Tests operations that take significant time
  - Validates timeout and cancellation handling
  - Test cases:
    - Archive creation with timeout
    - Verification with timeout
    - Context cancellation during operations
    - Resource cleanup after cancellation
- `TestResourceExhaustion`: Tests behavior under resource constraints
  - Validates graceful handling of resource limitations
  - Test cases:
    - Low disk space scenarios
    - Memory pressure scenarios
    - File descriptor exhaustion
    - Temporary directory full
- `TestConcurrentStress`: Tests high concurrency scenarios
  - Validates thread safety under stress
  - Test cases:
    - Many concurrent archive operations
    - Concurrent read/write operations
    - Resource contention scenarios
    - Cleanup under high concurrency

## Test Infrastructure

### Test Utilities
**Implementation**: `testutil/` package
**Requirements**:
- `CreateTestDirectory()`: Creates temporary test directories with files
- `CreateTestGitRepo()`: Creates temporary Git repositories for testing
- `CreateTestArchive()`: Creates test ZIP archives
- `AssertNoTempFiles()`: Verifies no temporary files remain
- `AssertArchiveContents()`: Validates archive contents
- `MockFileSystem()`: Provides file system mocking capabilities
- `CaptureOutput()`: Captures stdout/stderr for testing
- `SetupTestConfig()`: Creates test configuration files
- `CleanupTestResources()`: Ensures test cleanup

### Test Data Management
**Requirements**:
- Temporary directories created for each test
- Test files with various sizes and types
- Test Git repositories with different states
- Test configuration files with various settings
- Cleanup verification after each test
- No test artifacts left in file system

### Continuous Integration Requirements
**Requirements**:
- All tests must pass before merge
- Test coverage must be > 90%
- Performance benchmarks must not regress
- Memory usage tests must pass
- Concurrent operation tests must pass
- Resource cleanup verification required
- Cross-platform testing (Linux, macOS)

## Test Execution

### Local Development
```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test categories
go test -run TestUnit ./...
go test -run TestIntegration ./...
go test -run Benchmark ./...

# Run tests with race detection
go test -race ./...

# Run stress tests
go test -tags=stress ./...
```

### Test Categories by Build Tags
- **Unit Tests**: Default, no build tags required
- **Integration Tests**: `-tags=integration`
- **Performance Tests**: `-tags=performance`
- **Stress Tests**: `-tags=stress`

### Test Environment Setup
**Requirements**:
- Temporary directory for test files
- Mock Git repositories
- Test configuration files
- Isolated test environment
- Cleanup verification
- Resource tracking

## Quality Gates

### Code Coverage Requirements
- **Minimum Coverage**: 90% overall
- **Critical Functions**: 100% coverage required
  - Archive creation functions
  - Error handling functions
  - Resource cleanup functions
  - Configuration loading functions
- **Integration Tests**: Must cover all CLI commands
- **Error Paths**: All error conditions must be tested

### Performance Requirements
- **Archive Creation**: < 1 second for 1000 files
- **Directory Comparison**: < 500ms for 1000 files
- **Verification**: < 2 seconds for 100MB archive
- **Memory Usage**: < 100MB for 10,000 files
- **No Memory Leaks**: All resources properly cleaned up

### Reliability Requirements
- **Resource Cleanup**: 100% cleanup verification
- **Error Recovery**: All error scenarios tested
- **Concurrent Safety**: No race conditions
- **Data Integrity**: Archive corruption detection
- **Platform Compatibility**: Tests pass on Linux and macOS

## Test Maintenance

### Test Documentation
- Each test function must have clear documentation
- Test cases must be documented with expected outcomes
- Performance benchmarks must have baseline measurements
- Error scenarios must be documented with expected behavior

### Test Data Management
- Test data must be generated programmatically
- No hardcoded file paths or system dependencies
- Test isolation must be maintained
- Cleanup must be verified for all tests

### Regression Testing
- All bug fixes must include regression tests
- Performance regressions must be detected
- Configuration changes must be tested
- Backward compatibility must be verified

This testing architecture ensures comprehensive coverage of BkpDir functionality while maintaining high quality standards and reliable operation across different environments and use cases. 