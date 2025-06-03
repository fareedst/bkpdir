// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically error injection utilities for comprehensive error testing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-E: Error injection framework tests - 🔧
// DECISION-REF: DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Comprehensive testing of error injection, propagation, and recovery

package testutil

import (
	"errors"
	"fmt"
	"syscall"
	"testing"
	"time"
)

func TestNewErrorInjector(t *testing.T) {
	injector := NewErrorInjector()

	if injector == nil {
		t.Fatal("NewErrorInjector should not return nil")
	}

	if !injector.IsEnabled() {
		t.Error("ErrorInjector should be enabled by default")
	}

	stats := injector.GetStats()
	if stats.TotalInjections != 0 {
		t.Error("New injector should have zero injections")
	}

	if len(injector.GetPropagationTraces()) != 0 {
		t.Error("New injector should have no propagation traces")
	}

	if len(injector.GetRecoveryAttempts()) != 0 {
		t.Error("New injector should have no recovery attempts")
	}
}

func TestErrorInjector_AddRemoveInjectionPoint(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:    "WriteFile",
		Path:         "/test/path",
		TriggerCount: 1,
		Error:        CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "/test/path", nil),
		Probability:  1.0,
	}

	// Add injection point
	injector.AddInjectionPoint("test_point", point)

	// Test injection
	injectedError, shouldInject := injector.ShouldInjectError("WriteFile", "/test/path")
	if !shouldInject {
		t.Error("Should inject error for matching operation and path")
	}

	if injectedError == nil {
		t.Fatal("Injected error should not be nil")
	}

	if injectedError.Type != ErrorTypeTransient {
		t.Errorf("Expected ErrorTypeTransient, got %v", injectedError.Type)
	}

	if injectedError.Category != ErrorCategoryFilesystem {
		t.Errorf("Expected ErrorCategoryFilesystem, got %v", injectedError.Category)
	}

	// Remove injection point
	injector.RemoveInjectionPoint("test_point")

	// Test no injection after removal
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test/path")
	if shouldInject {
		t.Error("Should not inject error after injection point removed")
	}
}

func TestErrorInjector_EnableDisable(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("test", point)

	// Test enabled (default)
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/any/path")
	if !shouldInject {
		t.Error("Should inject when enabled")
	}

	// Test disabled
	injector.Enable(false)
	if injector.IsEnabled() {
		t.Error("Should be disabled after Enable(false)")
	}

	_, shouldInject = injector.ShouldInjectError("WriteFile", "/any/path")
	if shouldInject {
		t.Error("Should not inject when disabled")
	}

	// Test re-enabled
	injector.Enable(true)
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/any/path")
	if !shouldInject {
		t.Error("Should inject when re-enabled")
	}
}

func TestErrorInjector_TriggerCount(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:    "WriteFile",
		Path:         "*",
		TriggerCount: 3, // Inject on 3rd occurrence
		Error:        CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "", nil),
		Probability:  1.0,
	}

	injector.AddInjectionPoint("trigger_test", point)

	// First two calls should not inject
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/test")
	if shouldInject {
		t.Error("Should not inject on first call")
	}

	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test")
	if shouldInject {
		t.Error("Should not inject on second call")
	}

	// Third call should inject
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test")
	if !shouldInject {
		t.Error("Should inject on third call")
	}

	// Fourth call should not inject (TriggerCount is specific occurrence)
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test")
	if shouldInject {
		t.Error("Should not inject on fourth call")
	}
}

func TestErrorInjector_MaxInjections(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:     "WriteFile",
		Path:          "*",
		TriggerCount:  0, // Always trigger
		MaxInjections: 2, // Maximum 2 injections
		Error:         CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "", nil),
		Probability:   1.0,
	}

	injector.AddInjectionPoint("max_test", point)

	// First injection should work
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/test")
	if !shouldInject {
		t.Error("Should inject first time")
	}

	// Second injection should work
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test")
	if !shouldInject {
		t.Error("Should inject second time")
	}

	// Third injection should not work (exceeded max)
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/test")
	if shouldInject {
		t.Error("Should not inject third time (exceeded max)")
	}
}

func TestErrorInjector_PatternMatching(t *testing.T) {
	injector := NewErrorInjector()

	// Operation pattern matching
	point1 := &ErrorInjectionPoint{
		Operation:   "Write",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "write error", "Write", "", nil),
		Probability: 1.0,
	}

	// Path pattern matching
	point2 := &ErrorInjectionPoint{
		Operation:   "*",
		Path:        "/tmp",
		Error:       CreateFilesystemError(ErrorTypeTransient, "tmp error", "", "/tmp", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("op_pattern", point1)
	injector.AddInjectionPoint("path_pattern", point2)

	// Test operation pattern matching
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/any/path")
	if !shouldInject {
		t.Error("Should match operation pattern 'Write' in 'WriteFile'")
	}

	// Test path pattern matching
	_, shouldInject = injector.ShouldInjectError("AnyOperation", "/tmp/file")
	if !shouldInject {
		t.Error("Should match path pattern '/tmp' in '/tmp/file'")
	}

	// Test no matching
	_, shouldInject = injector.ShouldInjectError("Read", "/home/user")
	if shouldInject {
		t.Error("Should not match when neither operation nor path patterns match")
	}
}

func TestErrorInjector_ConditionalFunction(t *testing.T) {
	injector := NewErrorInjector()

	// Custom condition: only inject for .txt files
	point := &ErrorInjectionPoint{
		Operation:   "*",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "txt file error", "", "", nil),
		Probability: 1.0,
		ConditionalFunc: func(operation, path string) bool {
			return len(path) >= 4 && path[len(path)-4:] == ".txt"
		},
	}

	injector.AddInjectionPoint("conditional", point)

	// Should inject for .txt files
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/path/file.txt")
	if !shouldInject {
		t.Error("Should inject for .txt files")
	}

	// Should not inject for other files
	_, shouldInject = injector.ShouldInjectError("WriteFile", "/path/file.dat")
	if shouldInject {
		t.Error("Should not inject for non-.txt files")
	}
}

func TestErrorInjector_Statistics(t *testing.T) {
	injector := NewErrorInjector()

	// Add multiple injection points for different types and categories
	point1 := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "fs error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	point2 := &ErrorInjectionPoint{
		Operation:   "GitCommand",
		Path:        "*",
		Error:       CreateGitError(ErrorTypePermanent, "git error", "GitCommand", "", fmt.Errorf("git failed")),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("fs_error", point1)
	injector.AddInjectionPoint("git_error", point2)

	// Trigger injections
	injector.ShouldInjectError("WriteFile", "/test1")
	injector.ShouldInjectError("WriteFile", "/test2")
	injector.ShouldInjectError("GitCommand", "/repo")

	// Check statistics
	stats := injector.GetStats()

	if stats.TotalInjections != 3 {
		t.Errorf("Expected 3 total injections, got %d", stats.TotalInjections)
	}

	if stats.InjectionsByType[ErrorTypeTransient] != 2 {
		t.Errorf("Expected 2 transient errors, got %d", stats.InjectionsByType[ErrorTypeTransient])
	}

	if stats.InjectionsByType[ErrorTypePermanent] != 1 {
		t.Errorf("Expected 1 permanent error, got %d", stats.InjectionsByType[ErrorTypePermanent])
	}

	if stats.InjectionsByCategory[ErrorCategoryFilesystem] != 2 {
		t.Errorf("Expected 2 filesystem errors, got %d", stats.InjectionsByCategory[ErrorCategoryFilesystem])
	}

	if stats.InjectionsByCategory[ErrorCategoryGit] != 1 {
		t.Errorf("Expected 1 git error, got %d", stats.InjectionsByCategory[ErrorCategoryGit])
	}

	if stats.OperationsAffected["WriteFile"] != 2 {
		t.Errorf("Expected WriteFile to be affected 2 times, got %d", stats.OperationsAffected["WriteFile"])
	}

	if stats.OperationsAffected["GitCommand"] != 1 {
		t.Errorf("Expected GitCommand to be affected 1 time, got %d", stats.OperationsAffected["GitCommand"])
	}
}

func TestErrorInjector_PropagationTracking(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("prop_test", point)

	// Inject error
	injectedError, shouldInject := injector.ShouldInjectError("WriteFile", "/test")
	if !shouldInject {
		t.Fatal("Should inject error for propagation test")
	}

	if injectedError == nil {
		t.Fatal("Injected error should not be nil")
	}

	// Get traces before tracking
	traces := injector.GetPropagationTraces()
	if len(traces) != 1 {
		t.Errorf("Expected 1 propagation trace, got %d", len(traces))
	}

	// Find the trace for our error
	var errorID string
	for id := range traces {
		errorID = id
		break
	}

	// Track propagation
	injector.TrackErrorPropagation(errorID, "HandleError", "ErrorHandler", "caught", ErrorTypeTransient, "error caught in handler")
	injector.TrackErrorPropagation(errorID, "ProcessError", "ErrorProcessor", "wrapped", ErrorTypeTransient, "error wrapped with context")

	// Check propagation log
	updatedTraces := injector.GetPropagationTraces()
	trace := updatedTraces[errorID]

	if len(trace.PropagationLog) != 2 {
		t.Errorf("Expected 2 propagation entries, got %d", len(trace.PropagationLog))
	}

	if trace.PropagationLog[0].Action != "caught" {
		t.Errorf("Expected first action to be 'caught', got '%s'", trace.PropagationLog[0].Action)
	}

	if trace.PropagationLog[1].Action != "wrapped" {
		t.Errorf("Expected second action to be 'wrapped', got '%s'", trace.PropagationLog[1].Action)
	}
}

func TestErrorInjector_RecoveryTracking(t *testing.T) {
	injector := NewErrorInjector()

	// Track recovery attempts
	injector.TrackRecoveryAttempt("error_1", "retry", 1, false, 100*time.Millisecond)
	injector.TrackRecoveryAttempt("error_1", "retry", 2, true, 200*time.Millisecond)
	injector.TrackRecoveryAttempt("error_2", "fallback", 1, true, 50*time.Millisecond)

	// Check recovery attempts
	attempts := injector.GetRecoveryAttempts()
	if len(attempts) != 3 {
		t.Errorf("Expected 3 recovery attempts, got %d", len(attempts))
	}

	// Check statistics
	stats := injector.GetStats()
	if stats.SuccessfulRecoveries != 2 {
		t.Errorf("Expected 2 successful recoveries, got %d", stats.SuccessfulRecoveries)
	}

	if stats.FailedRecoveries != 1 {
		t.Errorf("Expected 1 failed recovery, got %d", stats.FailedRecoveries)
	}

	// Check average recovery time (should be average of successful recoveries)
	expectedAverage := (200*time.Millisecond + 50*time.Millisecond) / 2
	if stats.AverageRecoveryTime != expectedAverage {
		t.Errorf("Expected average recovery time %v, got %v", expectedAverage, stats.AverageRecoveryTime)
	}
}

func TestErrorInjector_Reset(t *testing.T) {
	injector := NewErrorInjector()

	// Add some data
	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "test error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("test", point)
	injector.ShouldInjectError("WriteFile", "/test")
	injector.TrackRecoveryAttempt("error_1", "retry", 1, true, 100*time.Millisecond)

	// Verify data exists
	stats := injector.GetStats()
	if stats.TotalInjections == 0 {
		t.Error("Should have injections before reset")
	}

	if len(injector.GetRecoveryAttempts()) == 0 {
		t.Error("Should have recovery attempts before reset")
	}

	// Reset
	injector.Reset()

	// Verify data cleared
	stats = injector.GetStats()
	if stats.TotalInjections != 0 {
		t.Error("Should have no injections after reset")
	}

	if len(injector.GetRecoveryAttempts()) != 0 {
		t.Error("Should have no recovery attempts after reset")
	}

	if len(injector.GetPropagationTraces()) != 0 {
		t.Error("Should have no propagation traces after reset")
	}

	// Test that injection doesn't work after reset (injection points cleared)
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/test")
	if shouldInject {
		t.Error("Should not inject after reset (injection points cleared)")
	}
}

func TestInjectedError_Methods(t *testing.T) {
	cause := errors.New("underlying error")

	err := &InjectedError{
		Type:       ErrorTypeTransient,
		Category:   ErrorCategoryFilesystem,
		Message:    "test error message",
		Operation:  "WriteFile",
		Path:       "/test/path",
		Retryable:  true,
		RetryAfter: 100 * time.Millisecond,
		Cause:      cause,
		InjectedAt: time.Now(),
	}

	// Test Error() method
	errorStr := err.Error()
	expectedSubstr := "test error message (injected filesystem error): underlying error"
	if errorStr != expectedSubstr {
		t.Errorf("Expected error string to be '%s', got '%s'", expectedSubstr, errorStr)
	}

	// Test Unwrap() method
	if err.Unwrap() != cause {
		t.Error("Unwrap should return the underlying cause")
	}

	// Test TypeString() method
	if err.TypeString() != "transient" {
		t.Errorf("Expected TypeString() to be 'transient', got '%s'", err.TypeString())
	}

	// Test CategoryString() method
	if err.CategoryString() != "filesystem" {
		t.Errorf("Expected CategoryString() to be 'filesystem', got '%s'", err.CategoryString())
	}

	// Test with no cause
	errNoCause := &InjectedError{
		Type:     ErrorTypePermanent,
		Category: ErrorCategoryGit,
		Message:  "no cause error",
	}

	errorStrNoCause := errNoCause.Error()
	expectedNoCause := "no cause error (injected git error)"
	if errorStrNoCause != expectedNoCause {
		t.Errorf("Expected error string to be '%s', got '%s'", expectedNoCause, errorStrNoCause)
	}
}

func TestCreateErrorFunctions(t *testing.T) {
	// Test CreateFilesystemError
	fsErr := CreateFilesystemError(ErrorTypeTransient, "fs error", "WriteFile", "/path", errors.New("cause"))

	if fsErr.Type != ErrorTypeTransient {
		t.Error("Filesystem error should have correct type")
	}

	if fsErr.Category != ErrorCategoryFilesystem {
		t.Error("Filesystem error should have filesystem category")
	}

	if !fsErr.Retryable {
		t.Error("Transient filesystem error should be retryable")
	}

	// Test CreateGitError
	gitErr := CreateGitError(ErrorTypePermanent, "git error", "GetBranch", "/repo", errors.New("git cause"))

	if gitErr.Type != ErrorTypePermanent {
		t.Error("Git error should have correct type")
	}

	if gitErr.Category != ErrorCategoryGit {
		t.Error("Git error should have git category")
	}

	if gitErr.Retryable {
		t.Error("Permanent git error should not be retryable")
	}

	// Test CreatePermissionError
	permErr := CreatePermissionError(ErrorTypePermanent, "permission denied", "Chmod", "/file")

	if permErr.Category != ErrorCategoryPermission {
		t.Error("Permission error should have permission category")
	}

	if permErr.Cause != syscall.EACCES {
		t.Error("Permission error should have EACCES cause")
	}

	if permErr.Retryable {
		t.Error("Permission error should not be retryable")
	}

	// Test CreateResourceError
	resErr := CreateResourceError(ErrorTypeTransient, "disk full", "WriteFile", "/file")

	if resErr.Category != ErrorCategoryResource {
		t.Error("Resource error should have resource category")
	}

	if resErr.Cause != syscall.ENOSPC {
		t.Error("Resource error should have ENOSPC cause")
	}

	if !resErr.Retryable {
		t.Error("Transient resource error should be retryable")
	}

	// Test CreateNetworkError
	netErr := CreateNetworkError(ErrorTypeTransient, "connection failed", "HTTPGet", "http://test", errors.New("network cause"))

	if netErr.Category != ErrorCategoryNetwork {
		t.Error("Network error should have network category")
	}

	if !netErr.Retryable {
		t.Error("Transient network error should be retryable")
	}
}

func TestBuiltinScenarios(t *testing.T) {
	// Test filesystem error scenario
	fsScenario := GetFilesystemErrorScenario()

	if fsScenario.Name != "filesystem_errors" {
		t.Error("Filesystem scenario should have correct name")
	}

	if len(fsScenario.InjectionPoints) != 2 {
		t.Errorf("Expected 2 injection points, got %d", len(fsScenario.InjectionPoints))
	}

	if fsScenario.ExpectedResults.ExpectedInjectionCount != 2 {
		t.Error("Filesystem scenario should expect 2 injections")
	}

	// Test git error scenario
	gitScenario := GetGitErrorScenario()

	if gitScenario.Name != "git_errors" {
		t.Error("Git scenario should have correct name")
	}

	if len(gitScenario.InjectionPoints) != 2 {
		t.Errorf("Expected 2 injection points, got %d", len(gitScenario.InjectionPoints))
	}

	if gitScenario.ExpectedResults.ExpectedInjectionCount != 2 {
		t.Error("Git scenario should expect 2 injections")
	}
}

func TestRunErrorInjectionScenario(t *testing.T) {
	injector := NewErrorInjector()

	// Test with filesystem scenario
	scenario := GetFilesystemErrorScenario()

	err := RunErrorInjectionScenario(injector, scenario)
	if err != nil {
		t.Fatalf("Failed to run scenario: %v", err)
	}

	// Test that injection points were added
	_, shouldInject1 := injector.ShouldInjectError("WriteFile", "/test1")
	if !shouldInject1 {
		t.Error("Should inject for first WriteFile operation")
	}

	_, shouldInject2 := injector.ShouldInjectError("WriteFile", "/test2")
	if !shouldInject2 {
		t.Error("Should inject for second WriteFile operation")
	}

	// Verify statistics match expected results
	stats := injector.GetStats()
	if stats.TotalInjections != scenario.ExpectedResults.ExpectedInjectionCount {
		t.Errorf("Expected %d injections, got %d",
			scenario.ExpectedResults.ExpectedInjectionCount, stats.TotalInjections)
	}
}

func TestErrorInjector_DelayFunctionality(t *testing.T) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		DelayBefore: 10 * time.Millisecond,
		DelayAfter:  5 * time.Millisecond,
		Error:       CreateFilesystemError(ErrorTypeTransient, "delayed error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("delay_test", point)

	start := time.Now()
	_, shouldInject := injector.ShouldInjectError("WriteFile", "/test")
	elapsed := time.Since(start)

	if !shouldInject {
		t.Error("Should inject error")
	}

	// Should take at least DelayBefore + DelayAfter time
	expectedMinimum := point.DelayBefore + point.DelayAfter
	if elapsed < expectedMinimum {
		t.Errorf("Expected minimum delay of %v, got %v", expectedMinimum, elapsed)
	}
}

// Benchmarks for performance testing

func BenchmarkErrorInjector_ShouldInjectError(b *testing.B) {
	injector := NewErrorInjector()

	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "bench error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("bench", point)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		injector.ShouldInjectError("WriteFile", "/test")
	}
}

func BenchmarkErrorInjector_TrackErrorPropagation(b *testing.B) {
	injector := NewErrorInjector()

	// Add a trace
	point := &ErrorInjectionPoint{
		Operation:   "WriteFile",
		Path:        "*",
		Error:       CreateFilesystemError(ErrorTypeTransient, "bench error", "WriteFile", "", nil),
		Probability: 1.0,
	}

	injector.AddInjectionPoint("bench", point)
	injector.ShouldInjectError("WriteFile", "/test")

	// Get error ID
	traces := injector.GetPropagationTraces()
	var errorID string
	for id := range traces {
		errorID = id
		break
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		injector.TrackErrorPropagation(errorID, "TestOp", "TestComponent", "caught", ErrorTypeTransient, "test message")
	}
}

func BenchmarkCreateFilesystemError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CreateFilesystemError(ErrorTypeTransient, "benchmark error", "WriteFile", "/test", nil)
	}
}
