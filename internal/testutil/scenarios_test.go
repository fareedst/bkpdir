// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically comprehensive tests for integration testing orchestration.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// TEST-INFRA-001-F: Integration testing orchestration tests
// DECISION-REF: DEC-008 (Testing infrastructure), DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Comprehensive testing of scenario composition, execution, and orchestration

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestScenarioBuilder tests the scenario builder functionality
func TestScenarioBuilder(t *testing.T) {
	t.Run("BasicScenarioCreation", func(t *testing.T) {
		scenario := NewScenarioBuilder("test-001", "Basic Test Scenario").
			WithDescription("A basic test scenario").
			WithTags("unit", "basic").
			WithTimeout(10*time.Minute).
			WithFailFast(true).
			AddSetupStep("setup-001", "Create test files", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return runtime.CreateTestFile("test.txt", []byte("test content"))
			}).
			AddExecutionStep("exec-001", "Execute main operation", func(ctx context.Context, runtime *ScenarioRuntime) error {
				runtime.SetSharedData("operation_result", "success")
				return nil
			}).
			AddVerificationStep("verify-001", "Verify results", func(ctx context.Context, runtime *ScenarioRuntime) error {
				result, exists := runtime.GetSharedData("operation_result")
				if !exists || result != "success" {
					return fmt.Errorf("operation result not found or incorrect")
				}
				return nil
			}).
			AddCleanupStep("cleanup-001", "Clean up resources", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			Build()

		// Verify scenario structure
		if scenario.ID != "test-001" {
			t.Errorf("Expected scenario ID 'test-001', got '%s'", scenario.ID)
		}
		if scenario.Name != "Basic Test Scenario" {
			t.Errorf("Expected scenario name 'Basic Test Scenario', got '%s'", scenario.Name)
		}
		if len(scenario.Steps) != 4 {
			t.Errorf("Expected 4 steps, got %d", len(scenario.Steps))
		}
		if len(scenario.Tags) != 2 {
			t.Errorf("Expected 2 tags, got %d", len(scenario.Tags))
		}
	})

	t.Run("ScenarioWithPrerequisites", func(t *testing.T) {
		scenario := NewScenarioBuilder("test-002", "Prerequisite Test").
			AddSetupStep("setup-001", "Setup A", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			AddSetupStep("setup-002", "Setup B", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			WithStepPrerequisites("setup-001").
			AddExecutionStep("exec-001", "Main execution", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			WithStepPrerequisites("setup-001", "setup-002").
			Build()

		// Verify prerequisites
		setupStep := scenario.Steps[1] // setup-002
		if len(setupStep.Prerequisites) != 1 || setupStep.Prerequisites[0] != "setup-001" {
			t.Errorf("Expected setup-002 to have prerequisite setup-001")
		}

		execStep := scenario.Steps[2] // exec-001
		if len(execStep.Prerequisites) != 2 {
			t.Errorf("Expected exec-001 to have 2 prerequisites, got %d", len(execStep.Prerequisites))
		}
	})

	t.Run("ScenarioValidation", func(t *testing.T) {
		// Test empty scenario validation
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty scenario")
			}
		}()

		NewScenarioBuilder("", "").Build()
	})

	t.Run("InvalidPrerequisiteValidation", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid prerequisite")
			}
		}()

		NewScenarioBuilder("test-003", "Invalid Prerequisite Test").
			AddSetupStep("setup-001", "Setup", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			WithStepPrerequisites("nonexistent-step").
			Build()
	})
}

// TestScenarioOrchestrator tests the orchestrator functionality
func TestScenarioOrchestrator(t *testing.T) {
	t.Run("BasicOrchestration", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()

		// Create a simple scenario
		scenario := NewScenarioBuilder("orch-001", "Orchestration Test").
			AddSetupStep("setup", "Setup", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return runtime.CreateTestFile("orch-test.txt", []byte("orchestration test"))
			}).
			AddExecutionStep("execute", "Execute", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Verify file exists
				filePath := filepath.Join(runtime.WorkingDirectory, "orch-test.txt")
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					return fmt.Errorf("test file not found")
				}
				runtime.SetSharedData("file_verified", true)
				return nil
			}).
			AddVerificationStep("verify", "Verify", func(ctx context.Context, runtime *ScenarioRuntime) error {
				verified, exists := runtime.GetSharedData("file_verified")
				if !exists || verified != true {
					return fmt.Errorf("file verification failed")
				}
				return nil
			}).
			Build()

		// Register and execute scenario
		err := orchestrator.RegisterScenario(scenario)
		if err != nil {
			t.Fatalf("Failed to register scenario: %v", err)
		}

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "orch-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		// Verify execution results
		if !execution.Success {
			t.Errorf("Expected successful execution, got failure")
		}
		if execution.StepsExecuted != 3 {
			t.Errorf("Expected 3 steps executed, got %d", execution.StepsExecuted)
		}
		if execution.StepsFailed != 0 {
			t.Errorf("Expected 0 failed steps, got %d", execution.StepsFailed)
		}
	})

	t.Run("ScenarioWithFailure", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()

		scenario := NewScenarioBuilder("fail-001", "Failure Test").
			WithFailFast(true).
			AddSetupStep("setup", "Setup", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return nil
			}).
			AddExecutionStep("fail-step", "Failing Step", func(ctx context.Context, runtime *ScenarioRuntime) error {
				return fmt.Errorf("intentional failure")
			}).
			AddVerificationStep("verify", "Verify", func(ctx context.Context, runtime *ScenarioRuntime) error {
				t.Error("This step should not execute due to fail-fast")
				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "fail-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		// Verify failure handling
		if execution.Success {
			t.Errorf("Expected failed execution, got success")
		}
		if execution.StepsFailed == 0 {
			t.Errorf("Expected at least one failed step")
		}
	})

	t.Run("ScenarioWithRetries", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()
		var attemptCount int
		var mu sync.Mutex

		scenario := NewScenarioBuilder("retry-001", "Retry Test").
			AddExecutionStep("retry-step", "Retrying Step", func(ctx context.Context, runtime *ScenarioRuntime) error {
				mu.Lock()
				attemptCount++
				current := attemptCount
				mu.Unlock()

				if current < 3 {
					return fmt.Errorf("attempt %d failed", current)
				}
				return nil
			}).
			WithStepTimeout(1 * time.Minute).
			Build()

		// Manually set retryable for the step
		scenario.Steps[0].Retryable = true
		scenario.Steps[0].MaxRetries = 3

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "retry-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		// Verify retry behavior
		if !execution.Success {
			t.Errorf("Expected successful execution after retries")
		}

		mu.Lock()
		if attemptCount != 3 {
			t.Errorf("Expected 3 attempts, got %d", attemptCount)
		}
		mu.Unlock()
	})
}

// TestScenarioExecution tests scenario execution with different configurations
func TestScenarioExecution(t *testing.T) {
	t.Run("SequentialExecution", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()
		var executionOrder []string
		var mu sync.Mutex

		scenario := NewScenarioBuilder("seq-001", "Sequential Test").
			AddSetupStep("step-1", "Step 1", func(ctx context.Context, runtime *ScenarioRuntime) error {
				mu.Lock()
				executionOrder = append(executionOrder, "step-1")
				mu.Unlock()
				time.Sleep(10 * time.Millisecond)
				return nil
			}).
			AddSetupStep("step-2", "Step 2", func(ctx context.Context, runtime *ScenarioRuntime) error {
				mu.Lock()
				executionOrder = append(executionOrder, "step-2")
				mu.Unlock()
				time.Sleep(10 * time.Millisecond)
				return nil
			}).
			AddSetupStep("step-3", "Step 3", func(ctx context.Context, runtime *ScenarioRuntime) error {
				mu.Lock()
				executionOrder = append(executionOrder, "step-3")
				mu.Unlock()
				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "seq-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}

		// Verify sequential order
		mu.Lock()
		expectedOrder := []string{"step-1", "step-2", "step-3"}
		if len(executionOrder) != len(expectedOrder) {
			t.Errorf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
		}
		for i, expected := range expectedOrder {
			if i >= len(executionOrder) || executionOrder[i] != expected {
				t.Errorf("Expected step %s at position %d, got %s", expected, i, executionOrder[i])
			}
		}
		mu.Unlock()
	})

	t.Run("ParallelExecution", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()
		var startTimes sync.Map
		var endTimes sync.Map

		scenario := NewScenarioBuilder("par-001", "Parallel Test").
			AddExecutionStep("par-1", "Parallel Step 1", func(ctx context.Context, runtime *ScenarioRuntime) error {
				startTimes.Store("par-1", time.Now())
				time.Sleep(100 * time.Millisecond)
				endTimes.Store("par-1", time.Now())
				return nil
			}).
			AddExecutionStep("par-2", "Parallel Step 2", func(ctx context.Context, runtime *ScenarioRuntime) error {
				startTimes.Store("par-2", time.Now())
				time.Sleep(100 * time.Millisecond)
				endTimes.Store("par-2", time.Now())
				return nil
			}).
			AddExecutionStep("par-3", "Parallel Step 3", func(ctx context.Context, runtime *ScenarioRuntime) error {
				startTimes.Store("par-3", time.Now())
				time.Sleep(100 * time.Millisecond)
				endTimes.Store("par-3", time.Now())
				return nil
			}).
			EnableParallelExecution(PhaseExecution, "par-1", "par-2", "par-3").
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "par-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}

		// Verify parallel execution (steps should overlap in time)
		start1, _ := startTimes.Load("par-1")
		start2, _ := startTimes.Load("par-2")
		start3, _ := startTimes.Load("par-3")

		// All should start within a reasonable time window
		times := []time.Time{start1.(time.Time), start2.(time.Time), start3.(time.Time)}
		minTime := times[0]
		maxTime := times[0]
		for _, t := range times {
			if t.Before(minTime) {
				minTime = t
			}
			if t.After(maxTime) {
				maxTime = t
			}
		}

		if maxTime.Sub(minTime) > 50*time.Millisecond {
			t.Errorf("Parallel steps should start closer together, time span: %v", maxTime.Sub(minTime))
		}
	})

	t.Run("TimeoutHandling", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()

		scenario := NewScenarioBuilder("timeout-001", "Timeout Test").
			AddExecutionStep("slow-step", "Slow Step", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Use context-aware sleep that can be cancelled
				select {
				case <-time.After(2 * time.Second):
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			}).
			WithStepTimeout(100 * time.Millisecond).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "timeout-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		// Should fail due to timeout
		if execution.Success {
			t.Errorf("Expected execution to fail due to timeout")
		}
	})
}

// TestTimingCoordinator tests the timing coordination functionality
func TestTimingCoordinator(t *testing.T) {
	t.Run("BasicBarrier", func(t *testing.T) {
		coordinator := NewTimingCoordinator()
		coordinator.CreateBarrier("test-barrier", 3)

		var finishTimes []time.Time
		var mu sync.Mutex

		// Start 3 goroutines
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				// Simulate different work durations
				time.Sleep(time.Duration(id*50) * time.Millisecond)
				coordinator.WaitForBarrier("test-barrier")

				mu.Lock()
				finishTimes = append(finishTimes, time.Now())
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		// Verify all goroutines finished at approximately the same time
		mu.Lock()
		if len(finishTimes) != 3 {
			t.Errorf("Expected 3 finish times, got %d", len(finishTimes))
		}

		// Check that all finish times are close together
		if len(finishTimes) >= 2 {
			maxDiff := time.Duration(0)
			for i := 1; i < len(finishTimes); i++ {
				diff := finishTimes[i].Sub(finishTimes[i-1])
				if diff < 0 {
					diff = -diff
				}
				if diff > maxDiff {
					maxDiff = diff
				}
			}

			if maxDiff > 10*time.Millisecond {
				t.Errorf("Finish times too spread out: %v", maxDiff)
			}
		}
		mu.Unlock()
	})

	t.Run("SignalCoordination", func(t *testing.T) {
		coordinator := NewTimingCoordinator()
		coordinator.CreateSignal("test-signal")

		received := false
		var mu sync.Mutex

		go func() {
			coordinator.WaitForSignal("test-signal")
			mu.Lock()
			received = true
			mu.Unlock()
		}()

		// Wait a bit to ensure goroutine is waiting
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		if received {
			t.Error("Signal received too early")
		}
		mu.Unlock()

		// Send signal
		coordinator.SendSignal("test-signal")

		// Wait for signal to be processed
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		if !received {
			t.Error("Signal not received")
		}
		mu.Unlock()
	})

	t.Run("DelayApplication", func(t *testing.T) {
		coordinator := NewTimingCoordinator()
		coordinator.SetDelay("test-operation", 100*time.Millisecond)

		start := time.Now()
		coordinator.ApplyDelay("test-operation")
		elapsed := time.Since(start)

		if elapsed < 100*time.Millisecond {
			t.Errorf("Delay too short: %v", elapsed)
		}
		if elapsed > 150*time.Millisecond {
			t.Errorf("Delay too long: %v", elapsed)
		}
	})
}

// TestScenarioRuntime tests the runtime environment functionality
func TestScenarioRuntime(t *testing.T) {
	orchestrator := NewScenarioOrchestrator()

	t.Run("FileCreation", func(t *testing.T) {
		scenario := NewScenarioBuilder("runtime-001", "Runtime File Test").
			AddSetupStep("create-files", "Create test files", func(ctx context.Context, runtime *ScenarioRuntime) error {
				if err := runtime.CreateTestFile("test1.txt", []byte("content1")); err != nil {
					return err
				}
				if err := runtime.CreateTestFile("subdir/test2.txt", []byte("content2")); err != nil {
					return err
				}
				return nil
			}).
			AddVerificationStep("verify-files", "Verify files exist", func(ctx context.Context, runtime *ScenarioRuntime) error {
				files := runtime.GetCreatedFiles()
				if len(files) != 2 {
					return fmt.Errorf("expected 2 files, got %d", len(files))
				}

				// Verify file contents
				for _, file := range files {
					if _, err := os.Stat(file); os.IsNotExist(err) {
						return fmt.Errorf("file %s does not exist", file)
					}
				}
				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "runtime-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}
	})

	t.Run("SharedData", func(t *testing.T) {
		scenario := NewScenarioBuilder("runtime-002", "Runtime Data Test").
			AddSetupStep("set-data", "Set shared data", func(ctx context.Context, runtime *ScenarioRuntime) error {
				runtime.SetSharedData("test-key", "test-value")
				runtime.SetSharedData("number", 42)
				return nil
			}).
			AddExecutionStep("use-data", "Use shared data", func(ctx context.Context, runtime *ScenarioRuntime) error {
				value, exists := runtime.GetSharedData("test-key")
				if !exists {
					return fmt.Errorf("test-key not found in shared data")
				}
				if value != "test-value" {
					return fmt.Errorf("expected 'test-value', got %v", value)
				}

				number, exists := runtime.GetSharedData("number")
				if !exists {
					return fmt.Errorf("number not found in shared data")
				}
				if number != 42 {
					return fmt.Errorf("expected 42, got %v", number)
				}

				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "runtime-002")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}
	})

	t.Run("ExecutionSummary", func(t *testing.T) {
		scenario := NewScenarioBuilder("runtime-003", "Runtime Summary Test").
			AddSetupStep("setup", "Setup", func(ctx context.Context, runtime *ScenarioRuntime) error {
				runtime.CreateTestFile("summary-test.txt", []byte("test"))
				runtime.CreateTestDirectory("test-dir")

				// Check summary during execution
				summary := runtime.GetExecutionSummary()
				if summary["scenario_id"] != "runtime-003" {
					return fmt.Errorf("incorrect scenario ID in summary")
				}

				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "runtime-003")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}
	})
}

// TestComplexScenarios tests complex multi-step scenarios with error injection
func TestComplexScenarios(t *testing.T) {
	t.Run("ErrorInjectionIntegration", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()

		scenario := NewScenarioBuilder("complex-001", "Error Injection Test").
			AddSetupStep("setup-error-injector", "Setup error injection", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Configure error injection
				runtime.ErrorInjector.Enable(true)
				point := &ErrorInjectionPoint{
					Operation:     "test-operation",
					TriggerCount:  1,
					MaxInjections: 1,
					Error: &InjectedError{
						Type:     ErrorTypeTransient,
						Category: ErrorCategoryFilesystem,
						Message:  "Simulated filesystem error",
					},
				}
				runtime.ErrorInjector.AddInjectionPoint("fs-error", point)
				return nil
			}).
			AddExecutionStep("trigger-error", "Trigger error injection", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Simulate an operation that should trigger error injection
				if injectedErr, shouldInject := runtime.ErrorInjector.ShouldInjectError("test-operation", "/test/path"); shouldInject {
					return injectedErr
				}
				return nil
			}).
			WithStepValidation(func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Validation should fail because we expect an error
				return fmt.Errorf("expected error injection to occur")
			}).
			AddVerificationStep("verify-error-stats", "Verify error injection stats", func(ctx context.Context, runtime *ScenarioRuntime) error {
				stats := runtime.ErrorInjector.GetStats()
				if stats.TotalInjections == 0 {
					return fmt.Errorf("expected at least one error injection")
				}
				return nil
			}).
			Build()

		// This scenario expects the error injection step to fail
		scenario.ExpectedResults.ShouldSucceed = false
		scenario.ExpectedResults.ExpectedFailures = []string{"trigger-error"}

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "complex-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		// The overall scenario should "succeed" because we expected the failure
		if execution.StepsFailed == 0 {
			t.Errorf("Expected at least one step to fail due to error injection")
		}
	})

	t.Run("MultiPhaseComplexScenario", func(t *testing.T) {
		orchestrator := NewScenarioOrchestrator()

		scenario := NewScenarioBuilder("complex-002", "Multi-Phase Complex Test").
			WithDescription("A complex scenario testing multiple components").
			WithTags("complex", "integration", "multi-phase").
			// Setup phase
			AddSetupStep("setup-env", "Setup environment", func(ctx context.Context, runtime *ScenarioRuntime) error {
				runtime.SetSharedData("setup_complete", true)
				return runtime.CreateTestDirectory("workspace")
			}).
			AddSetupStep("setup-files", "Setup test files", func(ctx context.Context, runtime *ScenarioRuntime) error {
				if err := runtime.CreateTestFile("workspace/input.txt", []byte("input data")); err != nil {
					return err
				}
				if err := runtime.CreateTestFile("workspace/config.yml", []byte("test: true")); err != nil {
					return err
				}
				return nil
			}).
			WithStepPrerequisites("setup-env").
			// Execution phase
			AddExecutionStep("process-data", "Process input data", func(ctx context.Context, runtime *ScenarioRuntime) error {
				inputPath := filepath.Join(runtime.WorkingDirectory, "workspace", "input.txt")
				outputPath := filepath.Join(runtime.WorkingDirectory, "workspace", "output.txt")

				// Read input
				data, err := os.ReadFile(inputPath)
				if err != nil {
					return err
				}

				// Process (uppercase)
				processed := strings.ToUpper(string(data))

				// Write output
				if err := os.WriteFile(outputPath, []byte(processed), 0644); err != nil {
					return err
				}

				runtime.SetSharedData("processed_data", processed)
				return nil
			}).
			WithStepPrerequisites("setup-files").
			AddExecutionStep("validate-config", "Validate configuration", func(ctx context.Context, runtime *ScenarioRuntime) error {
				configPath := filepath.Join(runtime.WorkingDirectory, "workspace", "config.yml")
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return fmt.Errorf("config file not found")
				}
				runtime.SetSharedData("config_valid", true)
				return nil
			}).
			WithStepPrerequisites("setup-files").
			// Verification phase
			AddVerificationStep("verify-output", "Verify output data", func(ctx context.Context, runtime *ScenarioRuntime) error {
				processed, exists := runtime.GetSharedData("processed_data")
				if !exists {
					return fmt.Errorf("processed data not found")
				}
				if processed != "INPUT DATA" {
					return fmt.Errorf("expected 'INPUT DATA', got '%s'", processed)
				}

				outputPath := filepath.Join(runtime.WorkingDirectory, "workspace", "output.txt")
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					return fmt.Errorf("output file not found")
				}
				return nil
			}).
			WithStepPrerequisites("process-data").
			AddVerificationStep("verify-environment", "Verify environment state", func(ctx context.Context, runtime *ScenarioRuntime) error {
				setupComplete, exists := runtime.GetSharedData("setup_complete")
				if !exists || setupComplete != true {
					return fmt.Errorf("setup not completed")
				}

				configValid, exists := runtime.GetSharedData("config_valid")
				if !exists || configValid != true {
					return fmt.Errorf("config not validated")
				}

				return nil
			}).
			WithStepPrerequisites("validate-config", "verify-output").
			// Cleanup phase
			AddCleanupStep("cleanup-files", "Clean up test files", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Files will be cleaned up automatically by orchestrator
				runtime.SetSharedData("cleanup_complete", true)
				return nil
			}).
			// Configure parallel execution for some steps
			EnableParallelExecution(PhaseExecution, "process-data", "validate-config").
			WithExpectation(ScenarioExpectation{
				ShouldSucceed:      true,
				MinSuccessfulSteps: 6,
				RequiredArtifacts:  []string{"workspace/output.txt"},
				ResourceCleanup:    true,
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "complex-002")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}

		if execution.StepsExecuted < 6 {
			t.Errorf("Expected at least 6 steps to execute, got %d", execution.StepsExecuted)
		}

		if execution.StepsFailed > 0 {
			t.Errorf("Expected no failed steps, got %d", execution.StepsFailed)
		}

		// Verify that all phases were executed
		phasesSeen := make(map[ScenarioPhase]bool)
		for _, event := range execution.Events {
			if event.Type == EventPhaseStarted {
				phasesSeen[event.Phase] = true
			}
		}

		expectedPhases := []ScenarioPhase{PhaseSetup, PhaseExecution, PhaseVerification, PhaseCleanup}
		for _, phase := range expectedPhases {
			if !phasesSeen[phase] {
				t.Errorf("Phase %d was not executed", int(phase))
			}
		}
	})
}

// BenchmarkScenarioExecution benchmarks scenario execution performance
func BenchmarkScenarioExecution(b *testing.B) {
	orchestrator := NewScenarioOrchestrator()

	scenario := NewScenarioBuilder("bench-001", "Benchmark Test").
		AddSetupStep("setup", "Setup", func(ctx context.Context, runtime *ScenarioRuntime) error {
			return runtime.CreateTestFile("bench.txt", []byte("benchmark"))
		}).
		AddExecutionStep("execute", "Execute", func(ctx context.Context, runtime *ScenarioRuntime) error {
			runtime.SetSharedData("result", "complete")
			return nil
		}).
		AddVerificationStep("verify", "Verify", func(ctx context.Context, runtime *ScenarioRuntime) error {
			_, exists := runtime.GetSharedData("result")
			if !exists {
				return fmt.Errorf("result not found")
			}
			return nil
		}).
		Build()

	orchestrator.RegisterScenario(scenario)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		_, err := orchestrator.ExecuteScenario(ctx, "bench-001")
		if err != nil {
			b.Fatalf("Failed to execute scenario: %v", err)
		}
	}
}

// TestRegressionIntegration tests integration with existing test patterns
func TestRegressionIntegration(t *testing.T) {
	t.Run("IntegrationWithExistingTests", func(t *testing.T) {
		// This test demonstrates how the scenario framework can be integrated
		// with existing test patterns in the codebase

		orchestrator := NewScenarioOrchestrator()

		scenario := NewScenarioBuilder("regression-001", "Regression Integration Test").
			WithDescription("Tests integration with existing test patterns").
			AddSetupStep("setup-test-env", "Setup test environment", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Simulate setting up environment similar to existing tests
				if err := runtime.CreateTestFile(".bkpdir.yml", []byte("archive_dir_path: archives\nuse_current_dir_name: true\n")); err != nil {
					return err
				}
				if err := runtime.CreateTestFile("test.txt", []byte("test content")); err != nil {
					return err
				}
				return nil
			}).
			AddExecutionStep("simulate-archive-operation", "Simulate archive operation", func(ctx context.Context, runtime *ScenarioRuntime) error {
				// Simulate an archive operation that might be tested in existing test files
				runtime.SetSharedData("archive_created", true)
				runtime.SetSharedData("archive_name", "test-archive.zip")
				return nil
			}).
			AddVerificationStep("verify-archive-result", "Verify archive result", func(ctx context.Context, runtime *ScenarioRuntime) error {
				created, exists := runtime.GetSharedData("archive_created")
				if !exists || created != true {
					return fmt.Errorf("archive was not created")
				}

				name, exists := runtime.GetSharedData("archive_name")
				if !exists || name == "" {
					return fmt.Errorf("archive name not set")
				}

				return nil
			}).
			Build()

		orchestrator.RegisterScenario(scenario)

		ctx := context.Background()
		execution, err := orchestrator.ExecuteScenario(ctx, "regression-001")
		if err != nil {
			t.Fatalf("Failed to execute scenario: %v", err)
		}

		if !execution.Success {
			t.Errorf("Expected successful execution")
		}

		// Verify the scenario can be used to test patterns similar to existing tests
		summary := execution.FinalState
		if archiveName, exists := summary["archive_name"]; !exists || archiveName != "test-archive.zip" {
			t.Errorf("Expected archive name in final state")
		}
	})
}
