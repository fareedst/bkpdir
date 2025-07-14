// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically integration testing orchestration for combining multiple
// error conditions and operations in realistic test scenarios.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// TEST-INFRA-001-F: Integration testing orchestration
// DECISION-REF: DEC-008 (Testing infrastructure), DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Uses builder pattern for scenario composition with clear separation between setup, execution, and verification

package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ScenarioPhase represents different phases of scenario execution
type ScenarioPhase int

const (
	// PhaseSetup - initial setup phase
	PhaseSetup ScenarioPhase = iota
	// PhaseExecution - main execution phase
	PhaseExecution
	// PhaseVerification - verification and validation phase
	PhaseVerification
	// PhaseCleanup - cleanup and teardown phase
	PhaseCleanup
)

// ScenarioStep represents a single step in a test scenario
type ScenarioStep struct {
	ID            string
	Name          string
	Description   string
	Phase         ScenarioPhase
	Timeout       time.Duration
	Retryable     bool
	MaxRetries    int
	Prerequisites []string // IDs of steps that must complete first
	Function      ScenarioStepFunc
	Validation    ScenarioValidationFunc
	OnSuccess     ScenarioCallbackFunc
	OnFailure     ScenarioCallbackFunc
	OnTimeout     ScenarioCallbackFunc
}

// ScenarioStepFunc defines the function signature for scenario steps
type ScenarioStepFunc func(ctx context.Context, runtime *ScenarioRuntime) error

// ScenarioValidationFunc defines the function signature for step validation
type ScenarioValidationFunc func(ctx context.Context, runtime *ScenarioRuntime) error

// ScenarioCallbackFunc defines the function signature for step callbacks
type ScenarioCallbackFunc func(ctx context.Context, runtime *ScenarioRuntime, err error)

// ScenarioRuntime provides runtime context and utilities for scenario execution
type ScenarioRuntime struct {
	mu                  sync.RWMutex
	ScenarioID          string
	StartTime           time.Time
	WorkingDirectory    string
	ArchiveDirectory    string
	TempDirectory       string
	ConfigFile          string
	ErrorInjector       *ErrorInjector
	CorruptionManager   *ArchiveCorruptor
	PermissionSimulator *PermissionSimulator
	DiskSpaceSimulator  *DiskSpaceSimulator
	ContextController   *ContextController
	ExecutionLog        []ScenarioEvent
	StepResults         map[string]*StepResult
	SharedData          map[string]interface{}
	ResourcesAllocated  []string
	filesCreated        []string
	dirsCreated         []string
}

// ScenarioEvent represents an event during scenario execution
type ScenarioEvent struct {
	Timestamp time.Time
	Type      ScenarioEventType
	StepID    string
	Phase     ScenarioPhase
	Message   string
	Error     error
	Duration  time.Duration
	Data      map[string]interface{}
}

// ScenarioEventType defines types of scenario events
type ScenarioEventType int

const (
	// EventStepStarted indicates a step has started
	EventStepStarted ScenarioEventType = iota
	// EventStepCompleted indicates a step completed successfully
	EventStepCompleted
	// EventStepFailed indicates a step failed
	EventStepFailed
	// EventStepRetried indicates a step was retried
	EventStepRetried
	// EventStepSkipped indicates a step was skipped due to prerequisites
	EventStepSkipped
	// EventScenarioStarted indicates the scenario started
	EventScenarioStarted
	// EventScenarioCompleted indicates the scenario completed
	EventScenarioCompleted
	// EventScenarioFailed indicates the scenario failed
	EventScenarioFailed
	// EventResourceAllocated indicates a resource was allocated
	EventResourceAllocated
	// EventResourceReleased indicates a resource was released
	EventResourceReleased
	// EventPhaseStarted indicates a new phase started
	EventPhaseStarted
	// EventPhaseCompleted indicates a phase completed
	EventPhaseCompleted
)

// StepResult contains the result of executing a scenario step
type StepResult struct {
	StepID       string
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	Success      bool
	Error        error
	RetryCount   int
	Phase        ScenarioPhase
	ValidationOK bool
	Data         map[string]interface{}
}

// ScenarioBuilder provides a fluent interface for building test scenarios
type ScenarioBuilder struct {
	scenario *TestScenario
}

// TestScenario represents a complete test scenario with multiple steps
type TestScenario struct {
	ID               string
	Name             string
	Description      string
	Tags             []string
	Timeout          time.Duration
	Steps            []*ScenarioStep
	SetupTimeout     time.Duration
	ExecutionTimeout time.Duration
	CleanupTimeout   time.Duration
	RequiredFeatures []string
	FailFast         bool
	ParallelSteps    map[string][]string // phase -> step IDs that can run in parallel
	Dependencies     []string            // Other scenario IDs this depends on
	ExpectedResults  ScenarioExpectation
}

// ScenarioExpectation defines what to expect from scenario execution
type ScenarioExpectation struct {
	ShouldSucceed       bool
	ExpectedFailures    []string // Step IDs expected to fail
	MinSuccessfulSteps  int
	MaxExecutionTime    time.Duration
	RequiredEvents      []ScenarioEventType
	RequiredArtifacts   []string // Files/directories that should be created
	ForbiddenArtifacts  []string // Files/directories that should NOT be created
	MemoryLeakThreshold int64    // Bytes
	ResourceCleanup     bool     // Whether all resources should be cleaned up
}

// TimingCoordinator handles timing and synchronization across scenario steps
type TimingCoordinator struct {
	mu           sync.RWMutex
	barriers     map[string]*sync.WaitGroup
	signals      map[string]chan struct{}
	delays       map[string]time.Duration
	checkpoints  map[string]time.Time
	synchronized bool
}

// ScenarioOrchestrator manages the execution of complex test scenarios
type ScenarioOrchestrator struct {
	mu               sync.RWMutex
	scenarios        map[string]*TestScenario
	activeRuntimes   map[string]*ScenarioRuntime
	coordinator      *TimingCoordinator
	globalTimeout    time.Duration
	maxConcurrency   int
	retryPolicy      RetryPolicy
	cleanupPolicy    CleanupPolicy
	resourceLimits   ResourceLimits
	executionHistory []ScenarioExecution
}

// ScenarioExecution tracks the execution of a scenario
type ScenarioExecution struct {
	ScenarioID    string
	ExecutionID   string
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	Success       bool
	StepsExecuted int
	StepsFailed   int
	StepsSkipped  int
	Events        []ScenarioEvent
	FinalState    map[string]interface{}
	ResourceUsage ResourceUsage
}

// ResourceUsage tracks resource consumption during scenario execution
type ResourceUsage struct {
	MaxMemoryUsed    int64
	TotalDiskUsed    int64
	FilesCreated     int
	DirectoriesUsed  int
	NetworkRequests  int
	ProcessesSpawned int
}

// CleanupPolicy defines how cleanup should be handled
type CleanupPolicy struct {
	AlwaysCleanup      bool
	CleanupOnSuccess   bool
	CleanupOnFailure   bool
	PreserveArtifacts  []string // Patterns for files to preserve
	CleanupTimeout     time.Duration
	ForceCleanupOnExit bool
}

// ResourceLimits defines limits on resource usage
type ResourceLimits struct {
	MaxMemoryMB        int64
	MaxDiskSpaceMB     int64
	MaxFiles           int
	MaxDirectories     int
	MaxExecutionTime   time.Duration
	MaxConcurrentSteps int
}

// NewScenarioBuilder creates a new scenario builder
func NewScenarioBuilder(id, name string) *ScenarioBuilder {
	return &ScenarioBuilder{
		scenario: &TestScenario{
			ID:               id,
			Name:             name,
			Tags:             make([]string, 0),
			Steps:            make([]*ScenarioStep, 0),
			ParallelSteps:    make(map[string][]string),
			Dependencies:     make([]string, 0),
			RequiredFeatures: make([]string, 0),
			Timeout:          30 * time.Minute,
			SetupTimeout:     5 * time.Minute,
			ExecutionTimeout: 20 * time.Minute,
			CleanupTimeout:   5 * time.Minute,
			FailFast:         true,
		},
	}
}

// WithDescription sets the scenario description
func (sb *ScenarioBuilder) WithDescription(description string) *ScenarioBuilder {
	sb.scenario.Description = description
	return sb
}

// WithTags adds tags to the scenario
func (sb *ScenarioBuilder) WithTags(tags ...string) *ScenarioBuilder {
	sb.scenario.Tags = append(sb.scenario.Tags, tags...)
	return sb
}

// WithTimeout sets the overall scenario timeout
func (sb *ScenarioBuilder) WithTimeout(timeout time.Duration) *ScenarioBuilder {
	sb.scenario.Timeout = timeout
	return sb
}

// WithFailFast sets whether the scenario should fail fast on first error
func (sb *ScenarioBuilder) WithFailFast(failFast bool) *ScenarioBuilder {
	sb.scenario.FailFast = failFast
	return sb
}

// WithDependencies adds scenario dependencies
func (sb *ScenarioBuilder) WithDependencies(scenarioIDs ...string) *ScenarioBuilder {
	sb.scenario.Dependencies = append(sb.scenario.Dependencies, scenarioIDs...)
	return sb
}

// WithRequiredFeatures adds required features
func (sb *ScenarioBuilder) WithRequiredFeatures(features ...string) *ScenarioBuilder {
	sb.scenario.RequiredFeatures = append(sb.scenario.RequiredFeatures, features...)
	return sb
}

// AddSetupStep adds a setup phase step
func (sb *ScenarioBuilder) AddSetupStep(id, name string, fn ScenarioStepFunc) *ScenarioBuilder {
	step := &ScenarioStep{
		ID:            id,
		Name:          name,
		Phase:         PhaseSetup,
		Function:      fn,
		Timeout:       5 * time.Minute,
		Retryable:     false,
		MaxRetries:    0,
		Prerequisites: make([]string, 0),
	}
	sb.scenario.Steps = append(sb.scenario.Steps, step)
	return sb
}

// AddExecutionStep adds an execution phase step
func (sb *ScenarioBuilder) AddExecutionStep(id, name string, fn ScenarioStepFunc) *ScenarioBuilder {
	step := &ScenarioStep{
		ID:            id,
		Name:          name,
		Phase:         PhaseExecution,
		Function:      fn,
		Timeout:       10 * time.Minute,
		Retryable:     true,
		MaxRetries:    3,
		Prerequisites: make([]string, 0),
	}
	sb.scenario.Steps = append(sb.scenario.Steps, step)
	return sb
}

// AddVerificationStep adds a verification phase step
func (sb *ScenarioBuilder) AddVerificationStep(id, name string, fn ScenarioStepFunc) *ScenarioBuilder {
	step := &ScenarioStep{
		ID:            id,
		Name:          name,
		Phase:         PhaseVerification,
		Function:      fn,
		Timeout:       5 * time.Minute,
		Retryable:     false,
		MaxRetries:    0,
		Prerequisites: make([]string, 0),
	}
	sb.scenario.Steps = append(sb.scenario.Steps, step)
	return sb
}

// AddCleanupStep adds a cleanup phase step
func (sb *ScenarioBuilder) AddCleanupStep(id, name string, fn ScenarioStepFunc) *ScenarioBuilder {
	step := &ScenarioStep{
		ID:            id,
		Name:          name,
		Phase:         PhaseCleanup,
		Function:      fn,
		Timeout:       2 * time.Minute,
		Retryable:     true,
		MaxRetries:    2,
		Prerequisites: make([]string, 0),
	}
	sb.scenario.Steps = append(sb.scenario.Steps, step)
	return sb
}

// WithStepTimeout sets timeout for the last added step
func (sb *ScenarioBuilder) WithStepTimeout(timeout time.Duration) *ScenarioBuilder {
	if len(sb.scenario.Steps) > 0 {
		sb.scenario.Steps[len(sb.scenario.Steps)-1].Timeout = timeout
	}
	return sb
}

// WithStepPrerequisites sets prerequisites for the last added step
func (sb *ScenarioBuilder) WithStepPrerequisites(stepIDs ...string) *ScenarioBuilder {
	if len(sb.scenario.Steps) > 0 {
		sb.scenario.Steps[len(sb.scenario.Steps)-1].Prerequisites = stepIDs
	}
	return sb
}

// WithStepValidation sets validation function for the last added step
func (sb *ScenarioBuilder) WithStepValidation(validation ScenarioValidationFunc) *ScenarioBuilder {
	if len(sb.scenario.Steps) > 0 {
		sb.scenario.Steps[len(sb.scenario.Steps)-1].Validation = validation
	}
	return sb
}

// WithStepCallbacks sets callback functions for the last added step
func (sb *ScenarioBuilder) WithStepCallbacks(onSuccess, onFailure, onTimeout ScenarioCallbackFunc) *ScenarioBuilder {
	if len(sb.scenario.Steps) > 0 {
		step := sb.scenario.Steps[len(sb.scenario.Steps)-1]
		step.OnSuccess = onSuccess
		step.OnFailure = onFailure
		step.OnTimeout = onTimeout
	}
	return sb
}

// EnableParallelExecution enables parallel execution for steps in the specified phase
func (sb *ScenarioBuilder) EnableParallelExecution(phase ScenarioPhase, stepIDs ...string) *ScenarioBuilder {
	phaseKey := fmt.Sprintf("phase_%d", int(phase))
	if sb.scenario.ParallelSteps[phaseKey] == nil {
		sb.scenario.ParallelSteps[phaseKey] = make([]string, 0)
	}
	sb.scenario.ParallelSteps[phaseKey] = append(sb.scenario.ParallelSteps[phaseKey], stepIDs...)
	return sb
}

// WithExpectation sets the expected results for the scenario
func (sb *ScenarioBuilder) WithExpectation(expectation ScenarioExpectation) *ScenarioBuilder {
	sb.scenario.ExpectedResults = expectation
	return sb
}

// Build finalizes and returns the test scenario
func (sb *ScenarioBuilder) Build() *TestScenario {
	// Validate the scenario before returning
	if err := sb.validateScenario(); err != nil {
		panic(fmt.Sprintf("Invalid scenario configuration: %v", err))
	}
	return sb.scenario
}

// validateScenario validates the scenario configuration
func (sb *ScenarioBuilder) validateScenario() error {
	if sb.scenario.ID == "" {
		return fmt.Errorf("scenario ID cannot be empty")
	}
	if sb.scenario.Name == "" {
		return fmt.Errorf("scenario name cannot be empty")
	}
	if len(sb.scenario.Steps) == 0 {
		return fmt.Errorf("scenario must have at least one step")
	}

	// Validate step prerequisites
	stepIDs := make(map[string]bool)
	for _, step := range sb.scenario.Steps {
		stepIDs[step.ID] = true
	}

	for _, step := range sb.scenario.Steps {
		for _, prereq := range step.Prerequisites {
			if !stepIDs[prereq] {
				return fmt.Errorf("step %s has invalid prerequisite %s", step.ID, prereq)
			}
		}
	}

	return nil
}

// NewScenarioOrchestrator creates a new scenario orchestrator
func NewScenarioOrchestrator() *ScenarioOrchestrator {
	return &ScenarioOrchestrator{
		scenarios:      make(map[string]*TestScenario),
		activeRuntimes: make(map[string]*ScenarioRuntime),
		coordinator:    NewTimingCoordinator(),
		globalTimeout:  60 * time.Minute,
		maxConcurrency: 10,
		retryPolicy: RetryPolicy{
			MaxAttempts:   3,
			BaseDelay:     1 * time.Second,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
		},
		cleanupPolicy: CleanupPolicy{
			AlwaysCleanup:      true,
			CleanupOnSuccess:   true,
			CleanupOnFailure:   false,
			CleanupTimeout:     5 * time.Minute,
			ForceCleanupOnExit: true,
		},
		resourceLimits: ResourceLimits{
			MaxMemoryMB:        1024,
			MaxDiskSpaceMB:     5120,
			MaxFiles:           10000,
			MaxDirectories:     1000,
			MaxExecutionTime:   120 * time.Minute,
			MaxConcurrentSteps: 5,
		},
		executionHistory: make([]ScenarioExecution, 0),
	}
}

// RegisterScenario registers a scenario with the orchestrator
func (so *ScenarioOrchestrator) RegisterScenario(scenario *TestScenario) error {
	so.mu.Lock()
	defer so.mu.Unlock()

	if _, exists := so.scenarios[scenario.ID]; exists {
		return fmt.Errorf("scenario with ID %s already registered", scenario.ID)
	}

	so.scenarios[scenario.ID] = scenario
	return nil
}

// ExecuteScenario executes a registered scenario
func (so *ScenarioOrchestrator) ExecuteScenario(ctx context.Context, scenarioID string) (*ScenarioExecution, error) {
	so.mu.RLock()
	scenario, exists := so.scenarios[scenarioID]
	so.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("scenario %s not found", scenarioID)
	}

	// Create runtime environment
	runtime, err := so.createScenarioRuntime(scenarioID)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	// Store active runtime
	so.mu.Lock()
	so.activeRuntimes[scenarioID] = runtime
	so.mu.Unlock()

	// Ensure cleanup
	defer func() {
		so.mu.Lock()
		delete(so.activeRuntimes, scenarioID)
		so.mu.Unlock()
		so.cleanupRuntime(runtime)
	}()

	execution := &ScenarioExecution{
		ScenarioID:  scenarioID,
		ExecutionID: fmt.Sprintf("%s_%d", scenarioID, time.Now().Unix()),
		StartTime:   time.Now(),
		Events:      make([]ScenarioEvent, 0),
		FinalState:  make(map[string]interface{}),
	}

	// Set up context with timeout
	execCtx, cancel := context.WithTimeout(ctx, scenario.Timeout)
	defer cancel()

	// Record scenario start
	runtime.recordEvent(EventScenarioStarted, "", PhaseSetup, "Scenario execution started", nil, 0, nil)

	// Execute scenario phases
	success := true
	phases := []ScenarioPhase{PhaseSetup, PhaseExecution, PhaseVerification, PhaseCleanup}

	for _, phase := range phases {
		phaseSuccess := so.executePhase(execCtx, scenario, runtime, phase)
		if !phaseSuccess && scenario.FailFast && phase != PhaseCleanup {
			success = false
			break
		}
		if !phaseSuccess {
			success = false
		}
	}

	// Finalize execution
	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime)
	execution.Success = success
	execution.Events = runtime.ExecutionLog
	execution.FinalState = runtime.SharedData

	// Calculate statistics
	execution.StepsExecuted = len(runtime.StepResults)
	for _, result := range runtime.StepResults {
		if !result.Success {
			execution.StepsFailed++
		}
	}

	// Record execution in history
	so.mu.Lock()
	so.executionHistory = append(so.executionHistory, *execution)
	so.mu.Unlock()

	// Record scenario completion
	eventType := EventScenarioCompleted
	if !success {
		eventType = EventScenarioFailed
	}
	runtime.recordEvent(eventType, "", PhaseCleanup, "Scenario execution completed", nil, execution.Duration, nil)

	return execution, nil
}

// createScenarioRuntime creates a runtime environment for scenario execution
func (so *ScenarioOrchestrator) createScenarioRuntime(scenarioID string) (*ScenarioRuntime, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("bkpdir_scenario_%s_", scenarioID))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	workingDir := filepath.Join(tempDir, "working")
	archiveDir := filepath.Join(tempDir, "archives")

	// Create directories
	for _, dir := range []string{workingDir, archiveDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create config file
	configFile := filepath.Join(workingDir, ".bkpdir.yml")
	configContent := fmt.Sprintf("archive_dir_path: %s\nuse_current_dir_name: true\n", archiveDir)
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create config file: %w", err)
	}

	runtime := &ScenarioRuntime{
		ScenarioID:         scenarioID,
		StartTime:          time.Now(),
		WorkingDirectory:   workingDir,
		ArchiveDirectory:   archiveDir,
		TempDirectory:      tempDir,
		ConfigFile:         configFile,
		ErrorInjector:      NewErrorInjector(),
		ContextController:  NewContextController(30 * time.Minute),
		ExecutionLog:       make([]ScenarioEvent, 0),
		StepResults:        make(map[string]*StepResult),
		SharedData:         make(map[string]interface{}),
		ResourcesAllocated: make([]string, 0),
		filesCreated:       make([]string, 0),
		dirsCreated:        make([]string, 0),
	}

	return runtime, nil
}

// executePhase executes all steps in a specific phase
func (so *ScenarioOrchestrator) executePhase(ctx context.Context, scenario *TestScenario, runtime *ScenarioRuntime, phase ScenarioPhase) bool {
	runtime.recordEvent(EventPhaseStarted, "", phase, fmt.Sprintf("Starting phase %d", int(phase)), nil, 0, nil)

	// Get steps for this phase
	phaseSteps := make([]*ScenarioStep, 0)
	for _, step := range scenario.Steps {
		if step.Phase == phase {
			phaseSteps = append(phaseSteps, step)
		}
	}

	if len(phaseSteps) == 0 {
		runtime.recordEvent(EventPhaseCompleted, "", phase, fmt.Sprintf("Phase %d completed (no steps)", int(phase)), nil, 0, nil)
		return true
	}

	// Check for parallel execution
	phaseKey := fmt.Sprintf("phase_%d", int(phase))
	parallelStepIDs := scenario.ParallelSteps[phaseKey]

	success := true
	if len(parallelStepIDs) > 0 {
		// Execute parallel steps
		success = so.executeStepsInParallel(ctx, phaseSteps, parallelStepIDs, runtime)
	} else {
		// Execute steps sequentially
		success = so.executeStepsSequentially(ctx, phaseSteps, runtime, scenario.FailFast)
	}

	eventType := EventPhaseCompleted
	message := fmt.Sprintf("Phase %d completed successfully", int(phase))
	if !success {
		message = fmt.Sprintf("Phase %d failed", int(phase))
	}

	runtime.recordEvent(eventType, "", phase, message, nil, 0, nil)
	return success
}

// executeStepsSequentially executes steps one after another
func (so *ScenarioOrchestrator) executeStepsSequentially(ctx context.Context, steps []*ScenarioStep, runtime *ScenarioRuntime, failFast bool) bool {
	allSuccess := true

	for _, step := range steps {
		if !so.arePrerequisitesMet(step, runtime) {
			runtime.recordEvent(EventStepSkipped, step.ID, step.Phase, fmt.Sprintf("Prerequisites not met for step %s", step.ID), nil, 0, nil)
			continue
		}

		success := so.executeStep(ctx, step, runtime)
		if !success {
			allSuccess = false
			if failFast {
				break
			}
		}
	}

	return allSuccess
}

// executeStepsInParallel executes specified steps in parallel
func (so *ScenarioOrchestrator) executeStepsInParallel(ctx context.Context, steps []*ScenarioStep, parallelStepIDs []string, runtime *ScenarioRuntime) bool {
	// Separate parallel and sequential steps
	parallelSteps := make([]*ScenarioStep, 0)
	sequentialSteps := make([]*ScenarioStep, 0)

	parallelSet := make(map[string]bool)
	for _, id := range parallelStepIDs {
		parallelSet[id] = true
	}

	for _, step := range steps {
		if parallelSet[step.ID] {
			parallelSteps = append(parallelSteps, step)
		} else {
			sequentialSteps = append(sequentialSteps, step)
		}
	}

	allSuccess := true

	// Execute sequential steps first
	if len(sequentialSteps) > 0 {
		if !so.executeStepsSequentially(ctx, sequentialSteps, runtime, false) {
			allSuccess = false
		}
	}

	// Execute parallel steps
	if len(parallelSteps) > 0 {
		var wg sync.WaitGroup
		resultChan := make(chan bool, len(parallelSteps))

		for _, step := range parallelSteps {
			if !so.arePrerequisitesMet(step, runtime) {
				runtime.recordEvent(EventStepSkipped, step.ID, step.Phase, fmt.Sprintf("Prerequisites not met for step %s", step.ID), nil, 0, nil)
				resultChan <- true
				continue
			}

			wg.Add(1)
			go func(s *ScenarioStep) {
				defer wg.Done()
				success := so.executeStep(ctx, s, runtime)
				resultChan <- success
			}(step)
		}

		wg.Wait()
		close(resultChan)

		// Collect results
		for success := range resultChan {
			if !success {
				allSuccess = false
			}
		}
	}

	return allSuccess
}

// executeStep executes a single scenario step
func (so *ScenarioOrchestrator) executeStep(ctx context.Context, step *ScenarioStep, runtime *ScenarioRuntime) bool {
	startTime := time.Now()

	runtime.recordEvent(EventStepStarted, step.ID, step.Phase, fmt.Sprintf("Starting step %s", step.Name), nil, 0, nil)

	result := &StepResult{
		StepID:    step.ID,
		StartTime: startTime,
		Phase:     step.Phase,
		Success:   false,
		Data:      make(map[string]interface{}),
	}

	// Set up step context with timeout
	stepCtx, cancel := context.WithTimeout(ctx, step.Timeout)
	defer cancel()

	var stepErr error

	// Execute with retries if enabled
	for attempt := 0; attempt <= step.MaxRetries; attempt++ {
		if attempt > 0 {
			runtime.recordEvent(EventStepRetried, step.ID, step.Phase, fmt.Sprintf("Retrying step %s (attempt %d)", step.Name, attempt+1), stepErr, 0, nil)
		}

		stepErr = step.Function(stepCtx, runtime)
		if stepErr == nil {
			result.Success = true
			break
		}

		result.RetryCount = attempt + 1

		if !step.Retryable || attempt >= step.MaxRetries {
			break
		}

		// Wait before retry
		select {
		case <-time.After(time.Duration(attempt+1) * time.Second):
		case <-stepCtx.Done():
			stepErr = stepCtx.Err()
			break
		}
	}

	// Run validation if provided and step succeeded
	if result.Success && step.Validation != nil {
		if validationErr := step.Validation(stepCtx, runtime); validationErr != nil {
			result.Success = false
			result.ValidationOK = false
			stepErr = fmt.Errorf("validation failed: %w", validationErr)
		} else {
			result.ValidationOK = true
		}
	}

	// Finalize result
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Error = stepErr

	// Store result
	runtime.mu.Lock()
	runtime.StepResults[step.ID] = result
	runtime.mu.Unlock()

	// Execute callbacks
	if result.Success && step.OnSuccess != nil {
		step.OnSuccess(stepCtx, runtime, nil)
	} else if !result.Success && step.OnFailure != nil {
		step.OnFailure(stepCtx, runtime, stepErr)
	}

	if stepCtx.Err() == context.DeadlineExceeded && step.OnTimeout != nil {
		step.OnTimeout(stepCtx, runtime, stepCtx.Err())
	}

	// Record completion
	eventType := EventStepCompleted
	message := fmt.Sprintf("Step %s completed successfully", step.Name)
	if !result.Success {
		eventType = EventStepFailed
		message = fmt.Sprintf("Step %s failed", step.Name)
	}

	runtime.recordEvent(eventType, step.ID, step.Phase, message, stepErr, result.Duration, nil)

	return result.Success
}

// arePrerequisitesMet checks if all prerequisites for a step are met
func (so *ScenarioOrchestrator) arePrerequisitesMet(step *ScenarioStep, runtime *ScenarioRuntime) bool {
	runtime.mu.RLock()
	defer runtime.mu.RUnlock()

	for _, prereqID := range step.Prerequisites {
		result, exists := runtime.StepResults[prereqID]
		if !exists || !result.Success {
			return false
		}
	}
	return true
}

// cleanupRuntime cleans up scenario runtime resources
func (so *ScenarioOrchestrator) cleanupRuntime(runtime *ScenarioRuntime) {
	if runtime.TempDirectory != "" {
		os.RemoveAll(runtime.TempDirectory)
	}

	if runtime.ContextController != nil {
		runtime.ContextController.Stop()
	}

	if runtime.ErrorInjector != nil {
		runtime.ErrorInjector.Reset()
	}
}

// NewTimingCoordinator creates a new timing coordinator
func NewTimingCoordinator() *TimingCoordinator {
	return &TimingCoordinator{
		barriers:     make(map[string]*sync.WaitGroup),
		signals:      make(map[string]chan struct{}),
		delays:       make(map[string]time.Duration),
		checkpoints:  make(map[string]time.Time),
		synchronized: false,
	}
}

// CreateBarrier creates a synchronization barrier
func (tc *TimingCoordinator) CreateBarrier(name string, count int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	wg := &sync.WaitGroup{}
	wg.Add(count)
	tc.barriers[name] = wg
}

// WaitForBarrier waits for a synchronization barrier
func (tc *TimingCoordinator) WaitForBarrier(name string) {
	tc.mu.RLock()
	barrier, exists := tc.barriers[name]
	tc.mu.RUnlock()

	if exists {
		barrier.Done()
		barrier.Wait()
	}
}

// SetDelay sets a delay for a specific operation
func (tc *TimingCoordinator) SetDelay(operation string, delay time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.delays[operation] = delay
}

// ApplyDelay applies any configured delay for an operation
func (tc *TimingCoordinator) ApplyDelay(operation string) {
	tc.mu.RLock()
	delay, exists := tc.delays[operation]
	tc.mu.RUnlock()

	if exists && delay > 0 {
		time.Sleep(delay)
	}
}

// CreateSignal creates a signal channel
func (tc *TimingCoordinator) CreateSignal(name string) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.signals[name] = make(chan struct{})
}

// WaitForSignal waits for a signal
func (tc *TimingCoordinator) WaitForSignal(name string) {
	tc.mu.RLock()
	signal, exists := tc.signals[name]
	tc.mu.RUnlock()

	if exists {
		<-signal
	}
}

// SendSignal sends a signal
func (tc *TimingCoordinator) SendSignal(name string) {
	tc.mu.RLock()
	signal, exists := tc.signals[name]
	tc.mu.RUnlock()

	if exists {
		select {
		case signal <- struct{}{}:
		default:
		}
	}
}

// recordEvent records an event in the scenario execution log
func (sr *ScenarioRuntime) recordEvent(eventType ScenarioEventType, stepID string, phase ScenarioPhase, message string, err error, duration time.Duration, data map[string]interface{}) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	event := ScenarioEvent{
		Timestamp: time.Now(),
		Type:      eventType,
		StepID:    stepID,
		Phase:     phase,
		Message:   message,
		Error:     err,
		Duration:  duration,
		Data:      data,
	}

	sr.ExecutionLog = append(sr.ExecutionLog, event)
}

// SetSharedData sets shared data accessible to all steps
func (sr *ScenarioRuntime) SetSharedData(key string, value interface{}) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	sr.SharedData[key] = value
}

// GetSharedData gets shared data
func (sr *ScenarioRuntime) GetSharedData(key string) (interface{}, bool) {
	sr.mu.RLock()
	defer sr.mu.RUnlock()
	value, exists := sr.SharedData[key]
	return value, exists
}

// CreateTestFile creates a test file in the working directory
func (sr *ScenarioRuntime) CreateTestFile(filename string, content []byte) error {
	filePath := filepath.Join(sr.WorkingDirectory, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return err
	}

	sr.mu.Lock()
	sr.filesCreated = append(sr.filesCreated, filePath)
	sr.mu.Unlock()

	return nil
}

// CreateTestDirectory creates a test directory in the working directory
func (sr *ScenarioRuntime) CreateTestDirectory(dirname string) error {
	dirPath := filepath.Join(sr.WorkingDirectory, dirname)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	sr.mu.Lock()
	sr.dirsCreated = append(sr.dirsCreated, dirPath)
	sr.mu.Unlock()

	return nil
}

// GetCreatedFiles returns list of files created during scenario
func (sr *ScenarioRuntime) GetCreatedFiles() []string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	result := make([]string, len(sr.filesCreated))
	copy(result, sr.filesCreated)
	return result
}

// GetExecutionSummary returns a summary of scenario execution
func (sr *ScenarioRuntime) GetExecutionSummary() map[string]interface{} {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	successful := 0
	failed := 0
	for _, result := range sr.StepResults {
		if result.Success {
			successful++
		} else {
			failed++
		}
	}

	return map[string]interface{}{
		"scenario_id":      sr.ScenarioID,
		"start_time":       sr.StartTime,
		"duration":         time.Since(sr.StartTime),
		"steps_successful": successful,
		"steps_failed":     failed,
		"total_events":     len(sr.ExecutionLog),
		"files_created":    len(sr.filesCreated),
		"dirs_created":     len(sr.dirsCreated),
	}
}
