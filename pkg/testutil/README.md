# Testing Utilities Package (testutil)

⭐ **EXTRACT-009: Testing Patterns and Utilities** - 🔧

This package provides comprehensive testing utilities extracted from the bkpdir codebase, offering reusable testing patterns and infrastructure for Go applications.

## Overview

The `testutil` package consolidates common testing patterns into a cohesive, interface-based testing framework. It provides utilities for:

- **Test Assertions**: Common assertion functions for various data types
- **File System Testing**: Temporary files, directories, and archive creation
- **CLI Testing**: Command creation, execution, and output capture
- **Configuration Testing**: Test configuration management and environment isolation
- **Test Fixtures**: Reusable test data and setup patterns
- **Test Scenarios**: Complex test orchestration and execution

## Quick Start

```go
package mypackage

import (
    "testing"
    "github.com/bkpdir/pkg/testutil"
)

func TestMyFunction(t *testing.T) {
    // Get the default test utility provider
    provider := testutil.GetDefaultProvider()
    
    // Use assertion helpers
    assert := provider.GetAssertionHelper()
    assert.AssertStringEqual(t, "test", "hello", "hello")
    
    // Use file system helpers
    fs := provider.GetFileSystemHelper()
    tempDir := fs.CreateTempDir(t, "mytest-")
    fs.CreateTempFile(t, tempDir, "test.txt", []byte("content"))
    
    // Use configuration helpers
    config := provider.GetConfigProvider()
    configData := map[string]interface{}{
        "key": "value",
    }
    config.SetConfigValue("new_key", "new_value")
}
```

## Core Interfaces

### TestUtilProvider

The main entry point that provides access to all testing utilities:

```go
provider := testutil.NewTestUtilProvider()

// Access individual helpers
configProvider := provider.GetConfigProvider()
envManager := provider.GetEnvironmentManager()
fsHelper := provider.GetFileSystemHelper()
cliHelper := provider.GetCliHelper()
assertHelper := provider.GetAssertionHelper()
fixtureManager := provider.GetFixtureManager()
```

### AssertionHelper

Provides common assertion functions:

```go
assert := testutil.NewAssertionHelper()

assert.AssertStringEqual(t, "name", got, want)
assert.AssertBoolEqual(t, "name", got, want)
assert.AssertIntEqual(t, "name", got, want)
assert.AssertSliceEqual(t, "name", got, want)
assert.AssertError(t, err, expectError)
assert.AssertContains(t, str, substr, "name")
```

### FileSystemTestHelper

File system testing utilities:

```go
fs := testutil.NewFileSystemTestHelper()

// Create temporary directories and files
tempDir := fs.CreateTempDir(t, "prefix-")
filePath := fs.CreateTempFile(t, tempDir, "file.txt", []byte("content"))

// Create multiple test files
files := map[string]string{
    "file1.txt": "content1",
    "dir/file2.txt": "content2",
}
fs.CreateTestFiles(t, tempDir, files)

// Create ZIP archives for testing
fs.CreateZipArchive(t, "archive.zip", files)
```

### CliTestHelper

CLI testing utilities:

```go
cli := testutil.NewCliTestHelper()

// Create test commands
cmd := cli.CreateTestCommand("test", func(cmd *cobra.Command, args []string) {
    // Command logic
})

// Execute commands and capture output
output, err := cli.ExecuteCommand(t, cmd, []string{"arg1", "arg2"})

// Capture stdout/stderr
stdout, stderr := cli.CaptureOutput(t, func() {
    // Code that prints to stdout/stderr
})
```

### ConfigProvider

Configuration testing utilities:

```go
configData := map[string]interface{}{
    "key": "value",
}
config := testutil.NewConfigProvider(configData)

// Get and set configuration values
data := config.GetConfigData()
config.SetConfigValue("new_key", "new_value")

// Save configuration to file
config.SaveToFile("config.yml")
```

### EnvironmentManager

Environment variable management:

```go
env := testutil.NewEnvironmentManager()

// Set and get environment variables
env.SetEnv("TEST_VAR", "test_value")
value := env.GetEnv("TEST_VAR")

// Unset variables
env.UnsetEnv("TEST_VAR")

// Restore original environment
env.RestoreEnv()
```

### TestFixtureManager

Test fixture and data management:

```go
fixtures := testutil.NewTestFixtureManager("testdata")

// Register test data
testData := map[string]interface{}{
    "key": "value",
}
fixtures.RegisterTestData("test_data", testData)

// Get registered data
data := fixtures.GetTestData("test_data")

// Save and load fixtures
fixtures.SaveFixture(t, "fixture_name", testData)
loaded, err := fixtures.LoadFixture(t, "fixture_name")
```

## Test Scenarios

The package provides a powerful scenario system for complex test orchestration:

```go
// Create a basic scenario
scenario := testutil.CreateBasicScenario("test_name", "description", func(t *testing.T) {
    // Test logic
})

// Create a file operation scenario
scenario := testutil.CreateFileOperationScenario("file_test", 
    map[string]string{"file.txt": "content"}, 
    func(t *testing.T, tempDir string) {
        // Test with files in tempDir
    })

// Create a configuration scenario
scenario := testutil.CreateConfigScenario("config_test",
    map[string]interface{}{"key": "value"},
    func(t *testing.T, configPath string) {
        // Test with configuration file
    })

// Run scenarios
testutil.RunScenario(t, scenario)
```

### Advanced Scenario Building

```go
scenario := testutil.NewScenarioBuilder("complex_test", "Complex test scenario").
    WithSetup(func(t *testing.T) error {
        // Setup phase
        return nil
    }).
    WithExecute(func(t *testing.T) error {
        // Execution phase
        return nil
    }).
    WithVerify(func(t *testing.T) error {
        // Verification phase
        return nil
    }).
    WithCleanup(func(t *testing.T) error {
        // Cleanup phase
        return nil
    }).
    Build()

testutil.RunScenario(t, scenario)
```

## Convenience Functions

The package provides many convenience functions for common operations:

```go
// Configuration testing
testutil.WithTestConfig(t, configData, func(configPath string) {
    // Test with configuration file
})

// Environment isolation
testutil.IsolateEnvironment(t, func() {
    // Test with isolated environment
})

// Fixture-based testing
testutil.WithTestFixture(t, "fixture_name", func(data map[string]interface{}) {
    // Test with fixture data
})

// File system operations
testutil.WithTempDir(t, func(tempDir string) {
    // Test with temporary directory
})
```

## Common Test Data

The package includes predefined test data sets:

```go
dataSet := testutil.CreateTestDataSet()

// Available data sets:
// - basic_config: Basic configuration data
// - advanced_config: Advanced configuration with all options
// - minimal_config: Minimal configuration
// - test_files: Various test file contents
// - error_scenarios: Error condition data
// - edge_cases: Edge case test data

// Register common data with a fixture manager
testutil.RegisterCommonTestData(fixtureManager)
```

## Integration Examples

### Testing with Configuration

```go
func TestWithConfiguration(t *testing.T) {
    configData := map[string]interface{}{
        "archive_dir_path": "test_archives",
        "compression_level": 6,
    }
    
    testutil.WithTestConfig(t, configData, func(configPath string) {
        // Test your function with the configuration file
        result, err := myFunction(configPath)
        
        assert := testutil.NewAssertionHelper()
        assert.AssertError(t, err, false)
        assert.AssertStringEqual(t, "result", result, "expected")
    })
}
```

### Testing CLI Commands

```go
func TestCliCommand(t *testing.T) {
    cli := testutil.NewCliTestHelper()
    
    cmd := cli.CreateTestCommand("mycommand", func(cmd *cobra.Command, args []string) {
        cmd.Print("Command executed successfully")
    })
    
    output, err := cli.ExecuteCommand(t, cmd, []string{})
    
    assert := testutil.NewAssertionHelper()
    assert.AssertError(t, err, false)
    assert.AssertContains(t, output, "successfully", "output check")
}
```

### Testing File Operations

```go
func TestFileOperations(t *testing.T) {
    files := map[string]string{
        "input.txt": "test input",
        "config.yml": "test: true",
    }
    
    scenario := testutil.CreateFileOperationScenario("file_test", files, 
        func(t *testing.T, tempDir string) {
            // Test your file operations
            result, err := processFiles(tempDir)
            
            assert := testutil.NewAssertionHelper()
            assert.AssertError(t, err, false)
            // Additional assertions...
        })
    
    testutil.RunScenario(t, scenario)
}
```

## Design Principles

1. **Interface-Based**: All utilities implement interfaces for maximum flexibility
2. **Automatic Cleanup**: All resources are automatically cleaned up using `t.Cleanup()`
3. **Thread-Safe**: All utilities are safe for concurrent use
4. **Composable**: Utilities can be combined and composed as needed
5. **Extracted Patterns**: All utilities are based on real patterns from the codebase

## Dependencies

- `github.com/spf13/cobra` v1.8.0 - CLI testing support
- `gopkg.in/yaml.v3` v3.0.1 - YAML fixture support

## Migration Guide

When migrating existing tests to use this package:

1. Replace direct assertion calls with `AssertionHelper` methods
2. Replace manual temp file/directory creation with `FileSystemTestHelper`
3. Replace manual environment variable management with `EnvironmentManager`
4. Replace manual configuration file creation with `ConfigProvider`
5. Use scenarios for complex test orchestration

## Performance

The package is designed for testing performance:

- Minimal overhead for simple operations
- Efficient resource management with automatic cleanup
- Concurrent-safe operations
- Optimized for test execution speed

## Contributing

When adding new utilities:

1. Follow the interface-based design pattern
2. Ensure thread safety
3. Implement automatic cleanup
4. Add comprehensive tests
5. Update documentation

## License

This package is part of the bkpdir project and follows the same license terms. 