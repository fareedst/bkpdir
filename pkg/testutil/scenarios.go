// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"fmt"
	"sync"
	"testing"
)

// DefaultTestScenario provides standard test scenario functionality.
type DefaultTestScenario struct {
	name        string
	description string
	setupFunc   func(t *testing.T) error
	executeFunc func(t *testing.T) error
	verifyFunc  func(t *testing.T) error
	cleanupFunc func(t *testing.T) error
	mu          sync.RWMutex
}

// NewTestScenario creates a new test scenario.
// Extracted from patterns in internal/testutil/scenarios.go.
//
// ⭐ EXTRACT-009: Test scenario creation - 🔧
func NewTestScenario(name, description string) TestScenario {
	return &DefaultTestScenario{
		name:        name,
		description: description,
	}
}

// Setup prepares the test scenario.
//
// ⭐ EXTRACT-009: Test scenario setup - 🔧
func (s *DefaultTestScenario) Setup(t *testing.T) error {
	t.Helper()

	s.mu.RLock()
	setupFunc := s.setupFunc
	s.mu.RUnlock()

	if setupFunc == nil {
		return nil // No setup required
	}

	return setupFunc(t)
}

// Execute runs the test scenario.
//
// ⭐ EXTRACT-009: Test scenario execution - 🔧
func (s *DefaultTestScenario) Execute(t *testing.T) error {
	t.Helper()

	s.mu.RLock()
	executeFunc := s.executeFunc
	s.mu.RUnlock()

	if executeFunc == nil {
		return fmt.Errorf("no execute function defined for scenario %q", s.name)
	}

	return executeFunc(t)
}

// Verify validates the test scenario results.
//
// ⭐ EXTRACT-009: Test scenario verification - 🔧
func (s *DefaultTestScenario) Verify(t *testing.T) error {
	t.Helper()

	s.mu.RLock()
	verifyFunc := s.verifyFunc
	s.mu.RUnlock()

	if verifyFunc == nil {
		return nil // No verification required
	}

	return verifyFunc(t)
}

// Cleanup cleans up the test scenario.
//
// ⭐ EXTRACT-009: Test scenario cleanup - 🔧
func (s *DefaultTestScenario) Cleanup(t *testing.T) error {
	t.Helper()

	s.mu.RLock()
	cleanupFunc := s.cleanupFunc
	s.mu.RUnlock()

	if cleanupFunc == nil {
		return nil // No cleanup required
	}

	return cleanupFunc(t)
}

// GetName returns the scenario name.
//
// ⭐ EXTRACT-009: Test scenario name retrieval - 🔧
func (s *DefaultTestScenario) GetName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.name
}

// GetDescription returns the scenario description.
//
// ⭐ EXTRACT-009: Test scenario description retrieval - 🔧
func (s *DefaultTestScenario) GetDescription() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.description
}

// SetSetupFunc sets the setup function for the scenario.
//
// ⭐ EXTRACT-009: Test scenario setup configuration - 🔧
func (s *DefaultTestScenario) SetSetupFunc(fn func(t *testing.T) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setupFunc = fn
}

// SetExecuteFunc sets the execute function for the scenario.
//
// ⭐ EXTRACT-009: Test scenario execution configuration - 🔧
func (s *DefaultTestScenario) SetExecuteFunc(fn func(t *testing.T) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executeFunc = fn
}

// SetVerifyFunc sets the verify function for the scenario.
//
// ⭐ EXTRACT-009: Test scenario verification configuration - 🔧
func (s *DefaultTestScenario) SetVerifyFunc(fn func(t *testing.T) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.verifyFunc = fn
}

// SetCleanupFunc sets the cleanup function for the scenario.
//
// ⭐ EXTRACT-009: Test scenario cleanup configuration - 🔧
func (s *DefaultTestScenario) SetCleanupFunc(fn func(t *testing.T) error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupFunc = fn
}

// ScenarioBuilder provides a fluent interface for building test scenarios.
// Extracted from patterns in internal/testutil/scenarios.go.
type ScenarioBuilder struct {
	scenario *DefaultTestScenario
}

// NewScenarioBuilder creates a new scenario builder.
//
// ⭐ EXTRACT-009: Test scenario builder creation - 🔧
func NewScenarioBuilder(name, description string) *ScenarioBuilder {
	return &ScenarioBuilder{
		scenario: &DefaultTestScenario{
			name:        name,
			description: description,
		},
	}
}

// WithSetup adds a setup function to the scenario.
//
// ⭐ EXTRACT-009: Scenario builder setup configuration - 🔧
func (b *ScenarioBuilder) WithSetup(fn func(t *testing.T) error) *ScenarioBuilder {
	b.scenario.SetSetupFunc(fn)
	return b
}

// WithExecute adds an execute function to the scenario.
//
// ⭐ EXTRACT-009: Scenario builder execution configuration - 🔧
func (b *ScenarioBuilder) WithExecute(fn func(t *testing.T) error) *ScenarioBuilder {
	b.scenario.SetExecuteFunc(fn)
	return b
}

// WithVerify adds a verify function to the scenario.
//
// ⭐ EXTRACT-009: Scenario builder verification configuration - 🔧
func (b *ScenarioBuilder) WithVerify(fn func(t *testing.T) error) *ScenarioBuilder {
	b.scenario.SetVerifyFunc(fn)
	return b
}

// WithCleanup adds a cleanup function to the scenario.
//
// ⭐ EXTRACT-009: Scenario builder cleanup configuration - 🔧
func (b *ScenarioBuilder) WithCleanup(fn func(t *testing.T) error) *ScenarioBuilder {
	b.scenario.SetCleanupFunc(fn)
	return b
}

// Build returns the constructed scenario.
//
// ⭐ EXTRACT-009: Scenario builder finalization - 🔧
func (b *ScenarioBuilder) Build() TestScenario {
	return b.scenario
}

// Package-level convenience functions

// RunScenario executes a complete test scenario with proper error handling.
// Extracted from patterns in internal/testutil/scenarios_test.go.
//
// ⭐ EXTRACT-009: Complete scenario execution - 🔧
func RunScenario(t *testing.T, scenario TestScenario) {
	t.Helper()

	t.Run(scenario.GetName(), func(t *testing.T) {
		// Setup phase
		if err := scenario.Setup(t); err != nil {
			t.Fatalf("Scenario setup failed: %v", err)
		}

		// Ensure cleanup runs even if other phases fail
		t.Cleanup(func() {
			if err := scenario.Cleanup(t); err != nil {
				t.Errorf("Scenario cleanup failed: %v", err)
			}
		})

		// Execute phase
		if err := scenario.Execute(t); err != nil {
			t.Fatalf("Scenario execution failed: %v", err)
		}

		// Verify phase
		if err := scenario.Verify(t); err != nil {
			t.Fatalf("Scenario verification failed: %v", err)
		}
	})
}

// CreateBasicScenario creates a basic test scenario with common patterns.
// Extracted from patterns across multiple test files.
//
// ⭐ EXTRACT-009: Basic scenario creation - 🔧
func CreateBasicScenario(name, description string, testFunc func(t *testing.T)) TestScenario {
	return NewScenarioBuilder(name, description).
		WithExecute(func(t *testing.T) error {
			testFunc(t)
			return nil
		}).
		Build()
}

// CreateFileOperationScenario creates a scenario for file operations.
// Extracted from patterns in comparison_test.go and verify_test.go.
//
// ⭐ EXTRACT-009: File operation scenario creation - 🔧
func CreateFileOperationScenario(name string, files map[string]string, operation func(t *testing.T, tempDir string)) TestScenario {
	var tempDir string

	return NewScenarioBuilder(name, "File operation test scenario").
		WithSetup(func(t *testing.T) error {
			tempDir = t.TempDir()

			// Create test files
			fsHelper := NewFileSystemTestHelper()
			fsHelper.CreateTestFiles(t, tempDir, files)

			return nil
		}).
		WithExecute(func(t *testing.T) error {
			operation(t, tempDir)
			return nil
		}).
		Build()
}

// CreateConfigScenario creates a scenario for configuration testing.
// Extracted from patterns in config_test.go and main_test.go.
//
// ⭐ EXTRACT-009: Configuration scenario creation - 🔧
func CreateConfigScenario(name string, configData map[string]interface{}, testFunc func(t *testing.T, configPath string)) TestScenario {
	var configPath string
	var cleanup func()

	return NewScenarioBuilder(name, "Configuration test scenario").
		WithSetup(func(t *testing.T) error {
			tempDir := t.TempDir()
			var err error
			configPath, cleanup = CreateTestConfig(t, tempDir, configData)
			return err
		}).
		WithExecute(func(t *testing.T) error {
			testFunc(t, configPath)
			return nil
		}).
		WithCleanup(func(t *testing.T) error {
			if cleanup != nil {
				cleanup()
			}
			return nil
		}).
		Build()
}
