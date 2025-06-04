// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"gopkg.in/yaml.v3"
)

// DefaultTestFixtureManager provides standard test fixture management functionality.
type DefaultTestFixtureManager struct {
	fixtureDir string
	testData   map[string]map[string]interface{}
	mu         sync.RWMutex
}

// NewTestFixtureManager creates a new test fixture manager.
// If fixtureDir is empty, it uses a default "testdata" directory.
//
// ⭐ EXTRACT-009: Test fixture manager creation - 🔧
func NewTestFixtureManager(fixtureDir string) TestFixtureManager {
	if fixtureDir == "" {
		fixtureDir = "testdata"
	}

	return &DefaultTestFixtureManager{
		fixtureDir: fixtureDir,
		testData:   make(map[string]map[string]interface{}),
	}
}

// LoadFixture loads test data from a fixture file.
// Supports both JSON and YAML formats based on file extension.
// Extracted from patterns in config_test.go and main_test.go.
//
// ⭐ EXTRACT-009: Test fixture loading - 🔧
func (m *DefaultTestFixtureManager) LoadFixture(t *testing.T, name string) (map[string]interface{}, error) {
	t.Helper()

	// Try different file extensions
	extensions := []string{".json", ".yaml", ".yml"}
	var filePath string
	var found bool

	for _, ext := range extensions {
		candidate := filepath.Join(m.fixtureDir, name+ext)
		if _, err := os.Stat(candidate); err == nil {
			filePath = candidate
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("fixture file not found for %q in %q", name, m.fixtureDir)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %q: %w", filePath, err)
	}

	var result map[string]interface{}
	ext := filepath.Ext(filePath)

	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("failed to parse JSON fixture %q: %w", filePath, err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("failed to parse YAML fixture %q: %w", filePath, err)
		}
	default:
		return nil, fmt.Errorf("unsupported fixture format %q", ext)
	}

	return result, nil
}

// SaveFixture saves test data to a fixture file.
// Uses JSON format by default, but can save as YAML if name ends with .yaml or .yml.
//
// ⭐ EXTRACT-009: Test fixture saving - 🔧
func (m *DefaultTestFixtureManager) SaveFixture(t *testing.T, name string, data map[string]interface{}) error {
	t.Helper()

	// Ensure fixture directory exists
	if err := os.MkdirAll(m.fixtureDir, 0755); err != nil {
		return fmt.Errorf("failed to create fixture directory %q: %w", m.fixtureDir, err)
	}

	// Determine format and file path
	var filePath string
	var fileData []byte
	var err error

	if filepath.Ext(name) == "" {
		// Default to JSON if no extension provided
		filePath = filepath.Join(m.fixtureDir, name+".json")
		fileData, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON data: %w", err)
		}
	} else {
		filePath = filepath.Join(m.fixtureDir, name)
		ext := filepath.Ext(name)

		switch ext {
		case ".json":
			fileData, err = json.MarshalIndent(data, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON data: %w", err)
			}
		case ".yaml", ".yml":
			fileData, err = yaml.Marshal(data)
			if err != nil {
				return fmt.Errorf("failed to marshal YAML data: %w", err)
			}
		default:
			return fmt.Errorf("unsupported fixture format %q", ext)
		}
	}

	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return fmt.Errorf("failed to write fixture file %q: %w", filePath, err)
	}

	return nil
}

// GetTestData returns predefined test data sets.
// This provides access to commonly used test data without file I/O.
//
// ⭐ EXTRACT-009: Test data retrieval - 🔧
func (m *DefaultTestFixtureManager) GetTestData(name string) map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if data, exists := m.testData[name]; exists {
		// Return a copy to prevent modification
		result := make(map[string]interface{})
		for k, v := range data {
			result[k] = v
		}
		return result
	}

	return nil
}

// RegisterTestData registers a test data set for later retrieval.
// This allows tests to register commonly used data sets.
//
// ⭐ EXTRACT-009: Test data registration - 🔧
func (m *DefaultTestFixtureManager) RegisterTestData(name string, data map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Store a copy to prevent external modification
	dataCopy := make(map[string]interface{})
	for k, v := range data {
		dataCopy[k] = v
	}

	m.testData[name] = dataCopy
}

// Package-level convenience functions

// LoadTestFixture loads a test fixture using the default fixture manager.
// Extracted from patterns in config_test.go and main_test.go.
//
// ⭐ EXTRACT-009: Convenience fixture loading - 🔧
func LoadTestFixture(t *testing.T, name string) (map[string]interface{}, error) {
	manager := NewTestFixtureManager("")
	return manager.LoadFixture(t, name)
}

// SaveTestFixture saves test data using the default fixture manager.
//
// ⭐ EXTRACT-009: Convenience fixture saving - 🔧
func SaveTestFixture(t *testing.T, name string, data map[string]interface{}) error {
	manager := NewTestFixtureManager("")
	return manager.SaveFixture(t, name, data)
}

// WithTestFixture loads a fixture and executes a function with the data.
// Provides automatic error handling and cleanup.
//
// ⭐ EXTRACT-009: Fixture-based test execution - 🔧
func WithTestFixture(t *testing.T, fixtureName string, fn func(data map[string]interface{})) {
	t.Helper()

	data, err := LoadTestFixture(t, fixtureName)
	if err != nil {
		t.Fatalf("Failed to load test fixture %q: %v", fixtureName, err)
	}

	fn(data)
}

// CreateTestDataSet creates a predefined test data set with common test values.
// Extracted from patterns across multiple test files.
//
// ⭐ EXTRACT-009: Common test data creation - 🔧
func CreateTestDataSet() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"basic_config": {
			"archive_dir_path":     "archives",
			"use_current_dir_name": true,
			"compression_level":    6,
			"exclude_patterns":     []string{".git", "*.tmp"},
			"include_hidden_files": false,
		},
		"advanced_config": {
			"archive_dir_path":     "/custom/archives",
			"use_current_dir_name": false,
			"custom_archive_name":  "custom-backup",
			"compression_level":    9,
			"exclude_patterns":     []string{".git", "node_modules", "*.log"},
			"include_hidden_files": true,
			"verify_archives":      true,
			"max_archive_size":     "1GB",
		},
		"minimal_config": {
			"archive_dir_path": ".",
		},
		"test_files": {
			"small_file":  "test content",
			"medium_file": "This is a medium-sized test file with more content for testing purposes.",
			"large_file":  generateLargeTestContent(1024),
		},
		"error_scenarios": {
			"invalid_path":      "/nonexistent/path",
			"permission_denied": "/root/restricted",
			"disk_full":         "/dev/full",
			"network_timeout":   "timeout_simulation",
		},
		"edge_cases": {
			"empty_string":     "",
			"unicode_content":  "测试内容 🚀 émojis",
			"special_chars":    "!@#$%^&*()[]{}|\\:;\"'<>?,./",
			"very_long_string": generateLongString(10000),
			"null_values":      nil,
		},
	}
}

// generateLargeTestContent creates test content of specified size.
func generateLargeTestContent(size int) string {
	content := make([]byte, size)
	for i := range content {
		content[i] = byte('A' + (i % 26))
	}
	return string(content)
}

// generateLongString creates a very long test string.
func generateLongString(length int) string {
	const pattern = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = pattern[i%len(pattern)]
	}
	return string(result)
}

// RegisterCommonTestData registers commonly used test data sets.
// This should be called during test setup to make common data available.
//
// ⭐ EXTRACT-009: Common test data registration - 🔧
func RegisterCommonTestData(manager TestFixtureManager) {
	testData := CreateTestDataSet()
	for name, data := range testData {
		manager.RegisterTestData(name, data)
	}
}
