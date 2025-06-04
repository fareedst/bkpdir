// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestDefaultTestUtilProvider tests the main provider functionality
func TestDefaultTestUtilProvider(t *testing.T) {
	provider := NewTestUtilProvider()

	// Test that all components are available
	if provider.GetConfigProvider() == nil {
		t.Error("ConfigProvider should not be nil")
	}

	if provider.GetEnvironmentManager() == nil {
		t.Error("EnvironmentManager should not be nil")
	}

	if provider.GetFileSystemHelper() == nil {
		t.Error("FileSystemTestHelper should not be nil")
	}

	if provider.GetCliHelper() == nil {
		t.Error("CliTestHelper should not be nil")
	}

	if provider.GetAssertionHelper() == nil {
		t.Error("AssertionHelper should not be nil")
	}

	if provider.GetFixtureManager() == nil {
		t.Error("TestFixtureManager should not be nil")
	}

	// Test scenario creation
	scenario := provider.CreateScenario("test", "Test scenario")
	if scenario == nil {
		t.Error("CreateScenario should not return nil")
	}

	if scenario.GetName() != "test" {
		t.Errorf("Expected scenario name 'test', got %q", scenario.GetName())
	}

	if scenario.GetDescription() != "Test scenario" {
		t.Errorf("Expected scenario description 'Test scenario', got %q", scenario.GetDescription())
	}
}

// TestAssertionHelper tests the assertion functionality
func TestAssertionHelper(t *testing.T) {
	helper := NewAssertionHelper()

	t.Run("StringEqual", func(t *testing.T) {
		// This should not fail
		helper.AssertStringEqual(t, "test", "hello", "hello")

		// Test with different strings - we can't easily test failure without causing test failure
		// So we just verify the method exists and can be called
	})

	t.Run("BoolEqual", func(t *testing.T) {
		helper.AssertBoolEqual(t, "test", true, true)
		helper.AssertBoolEqual(t, "test", false, false)
	})

	t.Run("IntEqual", func(t *testing.T) {
		helper.AssertIntEqual(t, "test", 42, 42)
		helper.AssertIntEqual(t, "test", 0, 0)
	})

	t.Run("SliceEqual", func(t *testing.T) {
		helper.AssertSliceEqual(t, "test", []string{"a", "b"}, []string{"a", "b"})
		helper.AssertSliceEqual(t, "test", []string{}, []string{})
	})

	t.Run("Error", func(t *testing.T) {
		helper.AssertError(t, nil, false)
		helper.AssertError(t, os.ErrNotExist, true)
	})

	t.Run("Contains", func(t *testing.T) {
		helper.AssertContains(t, "hello world", "world", "test")
		helper.AssertContains(t, "test string", "test", "test")
	})
}

// TestFileSystemTestHelper tests the file system utilities
func TestFileSystemTestHelper(t *testing.T) {
	helper := NewFileSystemTestHelper()

	t.Run("CreateTempDir", func(t *testing.T) {
		tempDir := helper.CreateTempDir(t, "testutil-test-")

		// Verify directory exists
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Errorf("Temp directory %q should exist", tempDir)
		}

		// Verify prefix
		if !strings.Contains(filepath.Base(tempDir), "testutil-test-") {
			t.Errorf("Temp directory should contain prefix 'testutil-test-', got %q", tempDir)
		}
	})

	t.Run("CreateTempFile", func(t *testing.T) {
		tempDir := t.TempDir()
		content := []byte("test content")

		filePath := helper.CreateTempFile(t, tempDir, "test.txt", content)

		// Verify file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Temp file %q should exist", filePath)
		}

		// Verify content
		readContent, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read temp file: %v", err)
		}

		if string(readContent) != string(content) {
			t.Errorf("Expected content %q, got %q", string(content), string(readContent))
		}
	})

	t.Run("CreateTestFiles", func(t *testing.T) {
		tempDir := t.TempDir()
		files := map[string]string{
			"file1.txt":     "content1",
			"dir/file2.txt": "content2",
		}

		helper.CreateTestFiles(t, tempDir, files)

		// Verify files exist and have correct content
		for relPath, expectedContent := range files {
			fullPath := filepath.Join(tempDir, relPath)

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				t.Errorf("File %q should exist", fullPath)
				continue
			}

			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read file %q: %v", fullPath, err)
				continue
			}

			if string(content) != expectedContent {
				t.Errorf("File %q: expected content %q, got %q", fullPath, expectedContent, string(content))
			}
		}
	})

	t.Run("CreateZipArchive", func(t *testing.T) {
		tempDir := t.TempDir()
		archivePath := filepath.Join(tempDir, "test.zip")
		files := map[string]string{
			"file1.txt": "content1",
			"file2.txt": "content2",
		}

		err := helper.CreateZipArchive(t, archivePath, files)
		if err != nil {
			t.Fatalf("Failed to create zip archive: %v", err)
		}

		// Verify archive exists
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			t.Errorf("Archive %q should exist", archivePath)
		}
	})
}

// TestCliTestHelper tests the CLI utilities
func TestCliTestHelper(t *testing.T) {
	helper := NewCliTestHelper()

	t.Run("CreateTestCommand", func(t *testing.T) {
		var executed bool
		cmd := helper.CreateTestCommand("test", func(cmd *cobra.Command, args []string) {
			executed = true
		})

		if cmd.Use != "test" {
			t.Errorf("Expected command use 'test', got %q", cmd.Use)
		}

		// Execute the command
		cmd.Run(cmd, []string{})

		if !executed {
			t.Error("Command should have been executed")
		}
	})

	t.Run("ExecuteCommand", func(t *testing.T) {
		cmd := helper.CreateTestCommand("test", func(cmd *cobra.Command, args []string) {
			cmd.Print("test output")
		})

		output, err := helper.ExecuteCommand(t, cmd, []string{})
		if err != nil {
			t.Fatalf("Command execution should not fail: %v", err)
		}

		if !strings.Contains(output, "test output") {
			t.Errorf("Expected output to contain 'test output', got %q", output)
		}
	})

	t.Run("CaptureOutput", func(t *testing.T) {
		stdout, stderr := helper.CaptureOutput(t, func() {
			// This would normally print to stdout/stderr
			// For testing, we'll just verify the function can be called
		})

		// stdout and stderr should be strings (even if empty)
		_ = stdout
		_ = stderr
	})
}

// TestConfigProvider tests the configuration utilities
func TestConfigProvider(t *testing.T) {
	configData := map[string]interface{}{
		"test_key":    "test_value",
		"number_key":  42,
		"boolean_key": true,
	}

	provider := NewConfigProvider(configData)

	t.Run("GetConfigData", func(t *testing.T) {
		data := provider.GetConfigData()

		if data["test_key"] != "test_value" {
			t.Errorf("Expected test_key='test_value', got %v", data["test_key"])
		}

		if data["number_key"] != 42 {
			t.Errorf("Expected number_key=42, got %v", data["number_key"])
		}

		if data["boolean_key"] != true {
			t.Errorf("Expected boolean_key=true, got %v", data["boolean_key"])
		}
	})

	t.Run("SetConfigValue", func(t *testing.T) {
		provider.SetConfigValue("new_key", "new_value")

		data := provider.GetConfigData()
		if data["new_key"] != "new_value" {
			t.Errorf("Expected new_key='new_value', got %v", data["new_key"])
		}
	})

	t.Run("SaveToFile", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "test_config.yml")

		err := provider.SaveToFile(configPath)
		if err != nil {
			t.Fatalf("Failed to save config to file: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Errorf("Config file %q should exist", configPath)
		}
	})
}

// TestEnvironmentManager tests the environment management utilities
func TestEnvironmentManager(t *testing.T) {
	manager := NewEnvironmentManager()

	t.Run("SetAndGetEnv", func(t *testing.T) {
		manager.SetEnv("TEST_VAR", "test_value")

		value := manager.GetEnv("TEST_VAR")
		if value != "test_value" {
			t.Errorf("Expected TEST_VAR='test_value', got %q", value)
		}
	})

	t.Run("UnsetEnv", func(t *testing.T) {
		manager.SetEnv("TEMP_VAR", "temp_value")
		manager.UnsetEnv("TEMP_VAR")

		value := manager.GetEnv("TEMP_VAR")
		if value != "" {
			t.Errorf("Expected TEMP_VAR to be empty after unset, got %q", value)
		}
	})

	t.Run("RestoreEnv", func(t *testing.T) {
		// Set a variable, then restore
		manager.SetEnv("RESTORE_TEST", "new_value")
		manager.RestoreEnv()

		// After restore, the variable should be back to its original state
		// (which is likely empty since we're in a test environment)
	})
}

// TestTestFixtureManager tests the fixture management utilities
func TestTestFixtureManager(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewTestFixtureManager(tempDir)

	t.Run("RegisterAndGetTestData", func(t *testing.T) {
		testData := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		}

		manager.RegisterTestData("test_data", testData)

		retrieved := manager.GetTestData("test_data")
		if retrieved == nil {
			t.Fatal("Retrieved test data should not be nil")
		}

		if retrieved["key1"] != "value1" {
			t.Errorf("Expected key1='value1', got %v", retrieved["key1"])
		}

		if retrieved["key2"] != 42 {
			t.Errorf("Expected key2=42, got %v", retrieved["key2"])
		}
	})

	t.Run("SaveAndLoadFixture", func(t *testing.T) {
		testData := map[string]interface{}{
			"fixture_key": "fixture_value",
			"number":      123,
		}

		// Save fixture
		err := manager.SaveFixture(t, "test_fixture", testData)
		if err != nil {
			t.Fatalf("Failed to save fixture: %v", err)
		}

		// Load fixture
		loaded, err := manager.LoadFixture(t, "test_fixture")
		if err != nil {
			t.Fatalf("Failed to load fixture: %v", err)
		}

		if loaded["fixture_key"] != "fixture_value" {
			t.Errorf("Expected fixture_key='fixture_value', got %v", loaded["fixture_key"])
		}

		// JSON unmarshaling converts numbers to float64
		if num, ok := loaded["number"].(float64); !ok || int(num) != 123 {
			t.Errorf("Expected number=123, got %v (type %T)", loaded["number"], loaded["number"])
		}
	})
}

// TestTestScenario tests the scenario functionality
func TestTestScenario(t *testing.T) {
	t.Run("BasicScenario", func(t *testing.T) {
		var setupCalled, executeCalled, verifyCalled, cleanupCalled bool

		scenario := NewScenarioBuilder("test_scenario", "Test scenario description").
			WithSetup(func(t *testing.T) error {
				setupCalled = true
				return nil
			}).
			WithExecute(func(t *testing.T) error {
				executeCalled = true
				return nil
			}).
			WithVerify(func(t *testing.T) error {
				verifyCalled = true
				return nil
			}).
			WithCleanup(func(t *testing.T) error {
				cleanupCalled = true
				return nil
			}).
			Build()

		// Test scenario properties
		if scenario.GetName() != "test_scenario" {
			t.Errorf("Expected name 'test_scenario', got %q", scenario.GetName())
		}

		if scenario.GetDescription() != "Test scenario description" {
			t.Errorf("Expected description 'Test scenario description', got %q", scenario.GetDescription())
		}

		// Execute scenario phases
		if err := scenario.Setup(t); err != nil {
			t.Fatalf("Setup should not fail: %v", err)
		}

		if err := scenario.Execute(t); err != nil {
			t.Fatalf("Execute should not fail: %v", err)
		}

		if err := scenario.Verify(t); err != nil {
			t.Fatalf("Verify should not fail: %v", err)
		}

		if err := scenario.Cleanup(t); err != nil {
			t.Fatalf("Cleanup should not fail: %v", err)
		}

		// Verify all phases were called
		if !setupCalled {
			t.Error("Setup should have been called")
		}

		if !executeCalled {
			t.Error("Execute should have been called")
		}

		if !verifyCalled {
			t.Error("Verify should have been called")
		}

		if !cleanupCalled {
			t.Error("Cleanup should have been called")
		}
	})

	t.Run("RunScenario", func(t *testing.T) {
		var executed bool

		scenario := CreateBasicScenario("run_test", "Run test scenario", func(t *testing.T) {
			executed = true
		})

		RunScenario(t, scenario)

		if !executed {
			t.Error("Scenario should have been executed")
		}
	})
}

// TestConvenienceFunctions tests package-level convenience functions
func TestConvenienceFunctions(t *testing.T) {
	t.Run("GetDefaultProvider", func(t *testing.T) {
		provider := GetDefaultProvider()
		if provider == nil {
			t.Error("GetDefaultProvider should not return nil")
		}
	})

	t.Run("GetProviderWithFixtures", func(t *testing.T) {
		tempDir := t.TempDir()
		provider := GetProviderWithFixtures(tempDir)
		if provider == nil {
			t.Error("GetProviderWithFixtures should not return nil")
		}
	})

	t.Run("CreateTestDataSet", func(t *testing.T) {
		dataSet := CreateTestDataSet()
		if dataSet == nil {
			t.Error("CreateTestDataSet should not return nil")
		}

		// Verify some expected data sets exist
		if dataSet["basic_config"] == nil {
			t.Error("basic_config data set should exist")
		}

		if dataSet["test_files"] == nil {
			t.Error("test_files data set should exist")
		}

		if dataSet["error_scenarios"] == nil {
			t.Error("error_scenarios data set should exist")
		}
	})
}
