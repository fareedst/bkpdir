// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically helper functions for easy integration with existing test patterns.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// TEST-INFRA-001-F: Integration testing orchestration helpers
// DECISION-REF: DEC-008 (Testing infrastructure), DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Convenience functions for common scenario patterns and integration with existing tests

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// errorTypeToString converts ErrorType to string representation
func errorTypeToString(errorType ErrorType) string {
	switch errorType {
	case ErrorTypeTransient:
		return "transient"
	case ErrorTypePermanent:
		return "permanent"
	case ErrorTypeRecoverable:
		return "recoverable"
	case ErrorTypeFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// CommonScenarioTemplates provides pre-built scenario templates for common testing patterns
type CommonScenarioTemplates struct {
	orchestrator *ScenarioOrchestrator
}

// NewCommonScenarioTemplates creates a new common scenario templates helper
func NewCommonScenarioTemplates() *CommonScenarioTemplates {
	return &CommonScenarioTemplates{
		orchestrator: NewScenarioOrchestrator(),
	}
}

// CreateBasicArchiveScenario creates a basic archive testing scenario
func (cst *CommonScenarioTemplates) CreateBasicArchiveScenario(scenarioID, description string) *TestScenario {
	return NewScenarioBuilder(scenarioID, "Basic Archive Test").
		WithDescription(description).
		WithTags("archive", "basic").
		AddSetupStep("setup-env", "Setup test environment", func(ctx context.Context, runtime *ScenarioRuntime) error {
			// Create test files
			if err := runtime.CreateTestFile("test.txt", []byte("test content")); err != nil {
				return err
			}
			if err := runtime.CreateTestFile("config.yml", []byte("test: true")); err != nil {
				return err
			}
			return runtime.CreateTestDirectory("subdir")
		}).
		AddExecutionStep("create-archive", "Create archive", func(ctx context.Context, runtime *ScenarioRuntime) error {
			// Simulate archive creation
			runtime.SetSharedData("archive_created", true)
			runtime.SetSharedData("archive_path", filepath.Join(runtime.ArchiveDirectory, "test-archive.tar.gz"))
			return nil
		}).
		AddVerificationStep("verify-archive", "Verify archive", func(ctx context.Context, runtime *ScenarioRuntime) error {
			created, exists := runtime.GetSharedData("archive_created")
			if !exists || created != true {
				return fmt.Errorf("archive was not created")
			}
			return nil
		}).
		AddCleanupStep("cleanup", "Clean up", func(ctx context.Context, runtime *ScenarioRuntime) error {
			return nil
		}).
		Build()
}

// CreateErrorInjectionScenario creates a scenario with error injection
func (cst *CommonScenarioTemplates) CreateErrorInjectionScenario(scenarioID, operation string, errorType ErrorType) *TestScenario {
	return NewScenarioBuilder(scenarioID, "Error Injection Test").
		WithDescription(fmt.Sprintf("Tests error injection for %s operation", operation)).
		WithTags("error-injection", "testing").
		AddSetupStep("setup-error-injection", "Setup error injection", func(ctx context.Context, runtime *ScenarioRuntime) error {
			runtime.ErrorInjector.Enable(true)
			point := &ErrorInjectionPoint{
				Operation:     operation,
				TriggerCount:  1,
				MaxInjections: 1,
				Error: &InjectedError{
					Type:     errorType,
					Category: ErrorCategoryFilesystem,
					Message:  fmt.Sprintf("Simulated %s error for %s", errorTypeToString(errorType), operation),
				},
			}
			runtime.ErrorInjector.AddInjectionPoint("test-error", point)
			return nil
		}).
		AddExecutionStep("trigger-operation", "Trigger operation with error", func(ctx context.Context, runtime *ScenarioRuntime) error {
			if injectedErr, shouldInject := runtime.ErrorInjector.ShouldInjectError(operation, "/test/path"); shouldInject {
				return injectedErr
			}
			return nil
		}).
		AddVerificationStep("verify-error-stats", "Verify error injection occurred", func(ctx context.Context, runtime *ScenarioRuntime) error {
			stats := runtime.ErrorInjector.GetStats()
			if stats.TotalInjections == 0 {
				return fmt.Errorf("expected error injection to occur")
			}
			return nil
		}).
		WithExpectation(ScenarioExpectation{
			ShouldSucceed:    false,
			ExpectedFailures: []string{"trigger-operation"},
			ResourceCleanup:  true,
		}).
		Build()
}

// CreatePermissionScenario creates a scenario that tests permission handling
func (cst *CommonScenarioTemplates) CreatePermissionScenario(scenarioID, testPath string, deniedMode os.FileMode) *TestScenario {
	return NewScenarioBuilder(scenarioID, "Permission Test").
		WithDescription(fmt.Sprintf("Tests permission handling for %s", testPath)).
		WithTags("permissions", "filesystem").
		AddSetupStep("setup-permissions", "Setup permission test", func(ctx context.Context, runtime *ScenarioRuntime) error {
			// Create test file/directory
			fullPath := filepath.Join(runtime.WorkingDirectory, testPath)
			if err := runtime.CreateTestFile(testPath, []byte("permission test")); err != nil {
				return err
			}

			// Store original permissions
			info, err := os.Stat(fullPath)
			if err != nil {
				return err
			}
			runtime.SetSharedData("original_mode", info.Mode())

			// Set restricted permissions
			if err := os.Chmod(fullPath, deniedMode); err != nil {
				return err
			}
			runtime.SetSharedData("restricted_mode", deniedMode)

			return nil
		}).
		AddExecutionStep("test-operation", "Test operation with restricted permissions", func(ctx context.Context, runtime *ScenarioRuntime) error {
			fullPath := filepath.Join(runtime.WorkingDirectory, testPath)

			// Try to read the file (should fail with restricted permissions)
			_, err := os.ReadFile(fullPath)
			runtime.SetSharedData("operation_error", err)

			if err != nil {
				return err // Expected to fail
			}
			return nil
		}).
		AddVerificationStep("verify-permission-error", "Verify permission error occurred", func(ctx context.Context, runtime *ScenarioRuntime) error {
			opError, exists := runtime.GetSharedData("operation_error")
			if !exists {
				return fmt.Errorf("expected operation to fail with permission error")
			}
			if opError == nil {
				return fmt.Errorf("expected permission error, but operation succeeded")
			}
			return nil
		}).
		AddCleanupStep("restore-permissions", "Restore original permissions", func(ctx context.Context, runtime *ScenarioRuntime) error {
			fullPath := filepath.Join(runtime.WorkingDirectory, testPath)
			originalMode, exists := runtime.GetSharedData("original_mode")
			if exists {
				return os.Chmod(fullPath, originalMode.(os.FileMode))
			}
			return nil
		}).
		WithExpectation(ScenarioExpectation{
			ShouldSucceed:    false,
			ExpectedFailures: []string{"test-operation"},
			ResourceCleanup:  true,
		}).
		Build()
}

// CreateConcurrentScenario creates a scenario that tests concurrent operations
func (cst *CommonScenarioTemplates) CreateConcurrentScenario(scenarioID string, numSteps int) *TestScenario {
	builder := NewScenarioBuilder(scenarioID, "Concurrent Operations Test").
		WithDescription(fmt.Sprintf("Tests %d concurrent operations", numSteps)).
		WithTags("concurrent", "performance")

	// Add concurrent execution steps
	stepIDs := make([]string, numSteps)
	for i := 0; i < numSteps; i++ {
		stepID := fmt.Sprintf("concurrent-step-%d", i+1)
		stepIDs[i] = stepID

		builder.AddExecutionStep(stepID, fmt.Sprintf("Concurrent Step %d", i+1),
			func(stepNum int) ScenarioStepFunc {
				return func(ctx context.Context, runtime *ScenarioRuntime) error {
					// Create unique file for this step
					filename := fmt.Sprintf("concurrent-%d.txt", stepNum)
					content := fmt.Sprintf("content from step %d", stepNum)

					if err := runtime.CreateTestFile(filename, []byte(content)); err != nil {
						return err
					}

					// Simulate some work
					time.Sleep(100 * time.Millisecond)

					runtime.SetSharedData(fmt.Sprintf("step_%d_complete", stepNum), true)
					return nil
				}
			}(i+1))
	}

	// Enable parallel execution for all steps
	builder.EnableParallelExecution(PhaseExecution, stepIDs...)

	// Add verification step
	builder.AddVerificationStep("verify-concurrent-results", "Verify all concurrent operations completed",
		func(ctx context.Context, runtime *ScenarioRuntime) error {
			for i := 1; i <= numSteps; i++ {
				key := fmt.Sprintf("step_%d_complete", i)
				completed, exists := runtime.GetSharedData(key)
				if !exists || completed != true {
					return fmt.Errorf("step %d did not complete", i)
				}
			}

			// Verify all files were created
			files := runtime.GetCreatedFiles()
			if len(files) < numSteps {
				return fmt.Errorf("expected at least %d files, got %d", numSteps, len(files))
			}

			return nil
		})

	return builder.Build()
}

// TestingHelpers provides utility functions for integrating scenarios with Go testing
type TestingHelpers struct {
	orchestrator *ScenarioOrchestrator
}

// NewTestingHelpers creates a new testing helpers instance
func NewTestingHelpers() *TestingHelpers {
	return &TestingHelpers{
		orchestrator: NewScenarioOrchestrator(),
	}
}

// RunScenarioTest runs a scenario as a test and reports results
func (th *TestingHelpers) RunScenarioTest(t *testing.T, scenario *TestScenario) *ScenarioExecution {
	t.Helper()

	// Register scenario
	if err := th.orchestrator.RegisterScenario(scenario); err != nil {
		t.Fatalf("Failed to register scenario %s: %v", scenario.ID, err)
	}

	// Execute scenario
	ctx := context.Background()
	execution, err := th.orchestrator.ExecuteScenario(ctx, scenario.ID)
	if err != nil {
		t.Fatalf("Failed to execute scenario %s: %v", scenario.ID, err)
	}

	// Report results
	t.Logf("Scenario %s completed in %v", scenario.ID, execution.Duration)
	t.Logf("Steps executed: %d, failed: %d, skipped: %d",
		execution.StepsExecuted, execution.StepsFailed, execution.StepsSkipped)

	// Check expectations
	expectation := scenario.ExpectedResults
	if expectation.ShouldSucceed && !execution.Success {
		t.Errorf("Expected scenario %s to succeed, but it failed", scenario.ID)
	}
	if !expectation.ShouldSucceed && execution.Success {
		t.Errorf("Expected scenario %s to fail, but it succeeded", scenario.ID)
	}

	if expectation.MinSuccessfulSteps > 0 {
		successfulSteps := execution.StepsExecuted - execution.StepsFailed
		if successfulSteps < expectation.MinSuccessfulSteps {
			t.Errorf("Expected at least %d successful steps, got %d",
				expectation.MinSuccessfulSteps, successfulSteps)
		}
	}

	if expectation.MaxExecutionTime > 0 && execution.Duration > expectation.MaxExecutionTime {
		t.Errorf("Expected execution time under %v, got %v",
			expectation.MaxExecutionTime, execution.Duration)
	}

	return execution
}

// RunScenarioSubtest runs a scenario as a subtest
func (th *TestingHelpers) RunScenarioSubtest(t *testing.T, scenario *TestScenario) {
	t.Helper()

	t.Run(scenario.Name, func(t *testing.T) {
		th.RunScenarioTest(t, scenario)
	})
}

// CreateQuickScenario creates a simple scenario for quick testing
func CreateQuickScenario(id, name string, testFunc ScenarioStepFunc) *TestScenario {
	return NewScenarioBuilder(id, name).
		WithDescription("Quick test scenario").
		WithTags("quick", "simple").
		AddExecutionStep("test", "Execute test", testFunc).
		AddVerificationStep("verify", "Verify results", func(ctx context.Context, runtime *ScenarioRuntime) error {
			// Basic verification - check that no errors occurred
			return nil
		}).
		Build()
}

// CreateFileOperationScenario creates a scenario for testing file operations
func CreateFileOperationScenario(id, operation string, files map[string][]byte) *TestScenario {
	return NewScenarioBuilder(id, fmt.Sprintf("File Operation: %s", operation)).
		WithDescription(fmt.Sprintf("Tests file operation: %s", operation)).
		WithTags("file-operations", operation).
		AddSetupStep("create-files", "Create test files", func(ctx context.Context, runtime *ScenarioRuntime) error {
			for filename, content := range files {
				if err := runtime.CreateTestFile(filename, content); err != nil {
					return fmt.Errorf("failed to create file %s: %w", filename, err)
				}
			}
			runtime.SetSharedData("files_created", len(files))
			return nil
		}).
		AddExecutionStep("perform-operation", fmt.Sprintf("Perform %s operation", operation),
			func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Store file list for verification
				createdFiles := runtime.GetCreatedFiles()
				runtime.SetSharedData("operation_files", createdFiles)
				runtime.SetSharedData("operation_completed", true)
				return nil
			}).
		AddVerificationStep("verify-files", "Verify file operation results", func(ctx context.Context, runtime *ScenarioRuntime) error {
			completed, exists := runtime.GetSharedData("operation_completed")
			if !exists || completed != true {
				return fmt.Errorf("operation was not completed")
			}

			filesCreated, exists := runtime.GetSharedData("files_created")
			if !exists {
				return fmt.Errorf("files_created not found")
			}

			if filesCreated != len(files) {
				return fmt.Errorf("expected %d files, operation indicates %v", len(files), filesCreated)
			}

			return nil
		}).
		Build()
}

// CreateBenchmarkScenario creates a scenario optimized for benchmarking
func CreateBenchmarkScenario(id, operation string, iterations int) *TestScenario {
	return NewScenarioBuilder(id, fmt.Sprintf("Benchmark: %s", operation)).
		WithDescription(fmt.Sprintf("Benchmarks %s operation with %d iterations", operation, iterations)).
		WithTags("benchmark", "performance").
		WithTimeout(10*time.Minute). // Longer timeout for benchmarks
		AddSetupStep("benchmark-setup", "Setup benchmark environment", func(ctx context.Context, runtime *ScenarioRuntime) error {
			runtime.SetSharedData("iterations", iterations)
			runtime.SetSharedData("start_time", time.Now())
			return nil
		}).
		AddExecutionStep("benchmark-execution", fmt.Sprintf("Execute %s benchmark", operation),
			func(ctx context.Context, runtime *ScenarioRuntime) error {
				startTime := time.Now()

				// Perform iterations
				for i := 0; i < iterations; i++ {
					// Simulate work
					runtime.SetSharedData(fmt.Sprintf("iteration_%d", i), time.Now())
				}

				endTime := time.Now()
				duration := endTime.Sub(startTime)

				runtime.SetSharedData("execution_duration", duration)
				runtime.SetSharedData("iterations_completed", iterations)

				return nil
			}).
		AddVerificationStep("benchmark-results", "Calculate benchmark results", func(ctx context.Context, runtime *ScenarioRuntime) error {
			duration, exists := runtime.GetSharedData("execution_duration")
			if !exists {
				return fmt.Errorf("execution duration not recorded")
			}

			completed, exists := runtime.GetSharedData("iterations_completed")
			if !exists || completed != iterations {
				return fmt.Errorf("not all iterations completed")
			}

			avgDuration := duration.(time.Duration) / time.Duration(iterations)
			runtime.SetSharedData("average_duration", avgDuration)

			// Store results for external access
			runtime.SetSharedData("benchmark_results", map[string]interface{}{
				"total_duration":   duration,
				"iterations":       iterations,
				"average_duration": avgDuration,
				"ops_per_second":   float64(iterations) / duration.(time.Duration).Seconds(),
			})

			return nil
		}).
		Build()
}

// ValidateScenarioIntegration validates that a scenario integrates properly with existing test patterns
func ValidateScenarioIntegration(scenario *TestScenario) error {
	if scenario.ID == "" {
		return fmt.Errorf("scenario must have an ID")
	}

	if scenario.Name == "" {
		return fmt.Errorf("scenario must have a name")
	}

	if len(scenario.Steps) == 0 {
		return fmt.Errorf("scenario must have at least one step")
	}

	// Check that phases are properly used
	phases := make(map[ScenarioPhase]bool)
	for _, step := range scenario.Steps {
		phases[step.Phase] = true
	}

	// Recommend having at least setup and execution phases
	if !phases[PhaseSetup] {
		return fmt.Errorf("scenario should have at least one setup step")
	}

	if !phases[PhaseExecution] {
		return fmt.Errorf("scenario should have at least one execution step")
	}

	return nil
}
