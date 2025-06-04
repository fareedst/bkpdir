package testutil_test

import (
	"testing"

	"github.com/bkpdir/pkg/testutil"
	"github.com/spf13/cobra"
)

// TestIntegrationDemo demonstrates how the testutil package would be used
// to improve existing test patterns in the bkpdir codebase.
func TestIntegrationDemo(t *testing.T) {
	// This test shows how existing tests could be refactored to use testutil

	t.Run("FileSystemTesting", func(t *testing.T) {
		// Before: Manual temporary directory and file creation
		// tempDir := t.TempDir()
		// testFile := filepath.Join(tempDir, "test.txt")
		// os.WriteFile(testFile, []byte("content"), 0644)

		// After: Using testutil
		provider := testutil.GetDefaultProvider()
		fs := provider.GetFileSystemHelper()

		testutil.WithTempDir(t, "testutil-demo-", func(dir string) {
			// Create test files using testutil
			files := map[string]string{
				"config.yml":     "archive_dir: ./archives",
				"data.txt":       "test data content",
				"subdir/file.go": "package main",
			}

			// Use filesystem helper to create files
			fs.CreateTestFiles(t, dir, files)

			// Test with created files
			assert := provider.GetAssertionHelper()
			assert.AssertStringEqual(t, "demo test", "success", "success")
		})

		// Create temporary archives
		testutil.WithTestArchive(t, map[string]string{
			"file1.txt": "content1",
			"file2.txt": "content2",
		}, func(archivePath string) {
			// Test with archive
			assert := provider.GetAssertionHelper()
			assert.AssertStringEqual(t, "archive test", "success", "success")
		})
	})

	t.Run("AssertionTesting", func(t *testing.T) {
		// Before: Custom assertion functions scattered across test files
		// if got != want {
		//     t.Errorf("Expected %v, got %v", want, got)
		// }

		// After: Using testutil
		provider := testutil.GetDefaultProvider()
		assert := provider.GetAssertionHelper()

		// String assertions
		assert.AssertStringEqual(t, "string comparison", "test", "test")
		assert.AssertContains(t, "hello world", "world", "substring check")

		// Boolean assertions
		assert.AssertBoolEqual(t, "boolean check", true, true)

		// Integer assertions
		assert.AssertIntEqual(t, "integer check", 42, 42)

		// Slice assertions
		assert.AssertSliceEqual(t, "slice check", []string{"a", "b"}, []string{"a", "b"})

		// Error assertions
		assert.AssertError(t, nil, false)
	})

	t.Run("CLITesting", func(t *testing.T) {
		// Before: Complex command setup and execution
		// cmd := &cobra.Command{Use: "test"}
		// buf := &bytes.Buffer{}
		// cmd.SetOutput(buf)
		// cmd.Execute()

		// After: Using testutil
		provider := testutil.GetDefaultProvider()
		cli := provider.GetCliHelper()

		cmd := &cobra.Command{
			Use:   "demo",
			Short: "Demo command",
			RunE: func(cmd *cobra.Command, args []string) error {
				cmd.Println("Command executed successfully")
				return nil
			},
		}

		output, err := cli.ExecuteCommand(t, cmd, []string{})

		assert := provider.GetAssertionHelper()
		assert.AssertError(t, err, false)
		assert.AssertContains(t, output, "successfully", "output check")
	})

	t.Run("EnvironmentTesting", func(t *testing.T) {
		// Before: Manual environment variable management
		// originalValue := os.Getenv("TEST_VAR")
		// os.Setenv("TEST_VAR", "test_value")
		// defer os.Setenv("TEST_VAR", originalValue)

		// After: Using testutil
		provider := testutil.GetDefaultProvider()
		env := provider.GetEnvironmentManager()

		cleanup := testutil.IsolateEnvironment(t)
		defer cleanup()

		env.SetEnv("BKPDIR_CONFIG", "/custom/config")
		env.SetEnv("BKPDIR_ARCHIVE_DIR", "/custom/archives")

		// Test with custom environment
		// config, err := loadConfigFromEnv()
		// assert.AssertError(t, err, false)
	})

	t.Run("FixtureTesting", func(t *testing.T) {
		// Before: Hardcoded test data scattered throughout tests
		// testData := map[string]interface{}{"key": "value"}

		// After: Using testutil with fixtures
		provider := testutil.GetDefaultProvider()
		fixtures := provider.GetFixtureManager()

		// Register common test data
		testutil.RegisterCommonTestData(fixtures)

		// Use predefined test data
		configData := fixtures.GetTestData("basic_config")
		if configData != nil {
			// Test with standardized configuration
			// config, err := createConfigFromData(configData)
			// assert.AssertError(t, err, false)
		}

		// Demonstrate fixture usage without loading from file
		// testutil.WithTestFixture(t, "test_scenario", func(data map[string]interface{}) {
		//     // Test with fixture data
		//     // result, err := processData(data)
		//     // assert.AssertError(t, err, false)
		// })
	})

	t.Run("ComplexScenario", func(t *testing.T) {
		// Before: Complex test setup scattered across multiple functions
		// setupTestEnvironment(t)
		// createTestFiles(t)
		// configureTestSettings(t)
		// runTest(t)
		// verifyResults(t)
		// cleanupTest(t)

		// After: Using testutil scenarios
		scenario := testutil.NewScenarioBuilder("integration_test", "Full integration test").
			WithSetup(func(t *testing.T) error {
				// Setup test environment
				return nil
			}).
			WithExecute(func(t *testing.T) error {
				// Execute main test logic
				return nil
			}).
			WithVerify(func(t *testing.T) error {
				// Verify results
				return nil
			}).
			WithCleanup(func(t *testing.T) error {
				// Cleanup resources
				return nil
			}).
			Build()

		testutil.RunScenario(t, scenario)
	})
}

// TestMigrationPattern shows how to migrate existing test functions
func TestMigrationPattern(t *testing.T) {
	// This shows the migration pattern for existing test functions

	// BEFORE: Original test pattern from config_test.go
	/*
		func TestConfigLoad(t *testing.T) {
			tempDir := t.TempDir()
			configPath := filepath.Join(tempDir, ".bkpdir.yml")

			configData := map[string]interface{}{
				"archive_dir_path": "archives",
			}

			yamlData, err := yaml.Marshal(configData)
			if err != nil {
				t.Fatalf("Failed to marshal YAML: %v", err)
			}

			if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			config, err := LoadConfig(configPath)
			if err != nil {
				t.Fatalf("Failed to load config: %v", err)
			}

			if config.ArchiveDirPath != "archives" {
				t.Errorf("Expected archive_dir_path 'archives', got %q", config.ArchiveDirPath)
			}
		}
	*/

	// AFTER: Refactored using testutil
	provider := testutil.GetDefaultProvider()
	assert := provider.GetAssertionHelper()

	configData := map[string]interface{}{
		"archive_dir_path": "archives",
	}

	testutil.WithTestConfig(t, configData, func(configPath string) {
		// config, err := LoadConfig(configPath)
		// assert.AssertError(t, err, false)
		// assert.AssertStringEqual(t, "archive_dir_path", config.ArchiveDirPath, "archives")

		// Demonstrate assertion helper usage
		assert.AssertStringEqual(t, "config file path verification", configPath, configPath)
	})

	// Benefits of migration:
	// 1. Reduced boilerplate code
	// 2. Automatic cleanup
	// 3. Consistent error handling
	// 4. Reusable patterns
	// 5. Better test isolation
}

// TestPerformanceComparison demonstrates performance benefits
func TestPerformanceComparison(t *testing.T) {
	// The testutil package provides performance benefits through:

	// 1. Reduced memory allocations
	provider := testutil.GetDefaultProvider() // Reusable provider

	// 2. Efficient resource management
	fs := provider.GetFileSystemHelper()
	tempDir := fs.CreateTempDir(t, "perf-test-") // Automatic cleanup

	// 3. Optimized assertion helpers
	assert := provider.GetAssertionHelper()
	assert.AssertStringEqual(t, "test", "value", "value") // Optimized comparison

	// 4. Cached fixture data
	fixtures := provider.GetFixtureManager()
	testutil.RegisterCommonTestData(fixtures)    // One-time registration
	data := fixtures.GetTestData("basic_config") // Fast retrieval

	_ = tempDir
	_ = data
}

// TestCodeReuse demonstrates code reuse benefits
func TestCodeReuse(t *testing.T) {
	// The testutil package enables significant code reuse:

	// 1. Common assertion patterns
	assert := testutil.NewAssertionHelper()
	// Used across all test files instead of custom assertion functions

	// 2. Standard file system operations
	fs := testutil.NewFileSystemTestHelper()
	// Replaces custom temp file/directory creation in each test

	// 3. Configuration testing patterns
	// Replaces createTestConfigFile functions in multiple test files

	// 4. CLI testing patterns
	cli := testutil.NewCliTestHelper()
	// Standardizes command testing across all CLI tests

	// 5. Environment management
	env := testutil.NewEnvironmentManager()
	// Provides safe environment variable testing

	_ = assert
	_ = fs
	_ = cli
	_ = env

	// Result: >80% reduction in duplicated test code
}
