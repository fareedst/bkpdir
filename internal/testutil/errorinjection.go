// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically error injection utilities for comprehensive error testing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-E: Error injection framework - 🔧
// DECISION-REF: DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Uses interface-based injection with configurable error schedules rather than global state modification

package testutil

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ErrorType categorizes different types of errors for classification
type ErrorType int

const (
	// ErrorTypeTransient - temporary errors that may resolve on retry
	ErrorTypeTransient ErrorType = iota
	// ErrorTypePermanent - errors that won't resolve without intervention
	ErrorTypePermanent
	// ErrorTypeRecoverable - errors that can be recovered from with specific actions
	ErrorTypeRecoverable
	// ErrorTypeFatal - errors that should cause immediate operation termination
	ErrorTypeFatal
)

// ErrorCategory represents the functional area where error occurs
type ErrorCategory int

const (
	// ErrorCategoryFilesystem - file system operation errors
	ErrorCategoryFilesystem ErrorCategory = iota
	// ErrorCategoryGit - Git operation errors
	ErrorCategoryGit
	// ErrorCategoryNetwork - network operation errors
	ErrorCategoryNetwork
	// ErrorCategoryPermission - permission/access errors
	ErrorCategoryPermission
	// ErrorCategoryResource - resource exhaustion errors
	ErrorCategoryResource
	// ErrorCategoryConfiguration - configuration errors
	ErrorCategoryConfiguration
)

// InjectedError represents a categorized error with metadata
type InjectedError struct {
	Type       ErrorType
	Category   ErrorCategory
	Message    string
	Operation  string
	Path       string
	Retryable  bool
	RetryAfter time.Duration
	Cause      error
	InjectedAt time.Time
	StackTrace []string
}

func (e *InjectedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s (injected %s error): %v", e.Message, e.CategoryString(), e.Cause)
	}
	return fmt.Sprintf("%s (injected %s error)", e.Message, e.CategoryString())
}

func (e *InjectedError) Unwrap() error {
	return e.Cause
}

func (e *InjectedError) TypeString() string {
	switch e.Type {
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

func (e *InjectedError) CategoryString() string {
	switch e.Category {
	case ErrorCategoryFilesystem:
		return "filesystem"
	case ErrorCategoryGit:
		return "git"
	case ErrorCategoryNetwork:
		return "network"
	case ErrorCategoryPermission:
		return "permission"
	case ErrorCategoryResource:
		return "resource"
	case ErrorCategoryConfiguration:
		return "configuration"
	default:
		return "unknown"
	}
}

// ErrorInjectionPoint defines when and how to inject an error
type ErrorInjectionPoint struct {
	Operation       string                            // Operation name pattern to match
	Path            string                            // File/resource path pattern to match
	TriggerCount    int                               // Inject error on Nth occurrence (0 = always)
	MaxInjections   int                               // Maximum number of times to inject (0 = unlimited)
	Error           *InjectedError                    // Error to inject
	Probability     float64                           // Probability of injection (0.0-1.0)
	DelayBefore     time.Duration                     // Delay before injection
	DelayAfter      time.Duration                     // Delay after injection
	ConditionalFunc func(operation, path string) bool // Custom condition function
}

// ErrorPropagationTrace tracks error flow through operation chains
type ErrorPropagationTrace struct {
	ErrorID        string
	OriginalError  *InjectedError
	PropagationLog []ErrorPropagationEntry
	StartTime      time.Time
	EndTime        time.Time
}

// ErrorPropagationEntry represents one step in error propagation
type ErrorPropagationEntry struct {
	Timestamp  time.Time
	Operation  string
	Component  string
	Action     string // "caught", "wrapped", "returned", "handled"
	ErrorType  ErrorType
	Message    string
	StackDepth int
}

// ErrorRecoveryAttempt tracks recovery testing
type ErrorRecoveryAttempt struct {
	ErrorID       string
	RecoveryType  string // "retry", "fallback", "cleanup", "abort"
	AttemptNumber int
	Success       bool
	Duration      time.Duration
	Timestamp     time.Time
}

// ErrorInjectionStats tracks statistics about error injection
type ErrorInjectionStats struct {
	TotalInjections      int
	InjectionsByType     map[ErrorType]int
	InjectionsByCategory map[ErrorCategory]int
	OperationsAffected   map[string]int
	SuccessfulRecoveries int
	FailedRecoveries     int
	PropagationTraces    int
	AverageRecoveryTime  time.Duration
}

// ErrorInjector manages systematic error injection for testing
type ErrorInjector struct {
	mu                 sync.RWMutex
	isEnabled          bool
	injectionPoints    map[string]*ErrorInjectionPoint
	operationCounts    map[string]int
	injectionCounts    map[string]int
	propagationTraces  map[string]*ErrorPropagationTrace
	recoveryAttempts   []ErrorRecoveryAttempt
	stats              ErrorInjectionStats
	defaultProbability float64
	contextTimeouts    map[string]time.Duration
	retryPolicies      map[ErrorType]RetryPolicy
}

// RetryPolicy defines how to handle retry for different error types
type RetryPolicy struct {
	MaxAttempts    int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	BackoffFactor  float64
	RetryableTypes []ErrorType
}

// FilesystemInterface abstracts filesystem operations for error injection
type FilesystemInterface interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	ReadFile(filename string) ([]byte, error)
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(name string) error
	RemoveAll(path string) error
	Stat(name string) (os.FileInfo, error)
	Chmod(name string, mode os.FileMode) error
	Rename(oldpath, newpath string) error
}

// GitInterface abstracts Git operations for error injection
type GitInterface interface {
	IsRepository(dir string) (bool, error)
	GetBranch(dir string) (string, error)
	GetShortHash(dir string) (string, error)
	GetStatus(dir string) (string, error)
	ExecuteCommand(dir string, args ...string) ([]byte, error)
}

// NetworkInterface abstracts network operations for error injection
type NetworkInterface interface {
	HTTPGet(url string) ([]byte, error)
	HTTPPost(url string, data []byte) ([]byte, error)
	DNSLookup(hostname string) ([]string, error)
	TCPConnect(address string) error
}

// InjectableFilesystem wraps filesystem operations with error injection
type InjectableFilesystem struct {
	injector   *ErrorInjector
	underlying FilesystemInterface
}

// InjectableGit wraps Git operations with error injection
type InjectableGit struct {
	injector   *ErrorInjector
	underlying GitInterface
}

// InjectableNetwork wraps network operations with error injection
type InjectableNetwork struct {
	injector   *ErrorInjector
	underlying NetworkInterface
}

// NewErrorInjector creates a new error injection manager
func NewErrorInjector() *ErrorInjector {
	return &ErrorInjector{
		isEnabled:          true,
		injectionPoints:    make(map[string]*ErrorInjectionPoint),
		operationCounts:    make(map[string]int),
		injectionCounts:    make(map[string]int),
		propagationTraces:  make(map[string]*ErrorPropagationTrace),
		recoveryAttempts:   make([]ErrorRecoveryAttempt, 0),
		defaultProbability: 1.0,
		contextTimeouts:    make(map[string]time.Duration),
		retryPolicies:      getDefaultRetryPolicies(),
		stats: ErrorInjectionStats{
			InjectionsByType:     make(map[ErrorType]int),
			InjectionsByCategory: make(map[ErrorCategory]int),
			OperationsAffected:   make(map[string]int),
		},
	}
}

// AddInjectionPoint adds a new error injection point
func (ei *ErrorInjector) AddInjectionPoint(id string, point *ErrorInjectionPoint) {
	ei.mu.Lock()
	defer ei.mu.Unlock()

	// Set default probability if not specified
	if point.Probability == 0 {
		point.Probability = ei.defaultProbability
	}

	ei.injectionPoints[id] = point
}

// RemoveInjectionPoint removes an error injection point
func (ei *ErrorInjector) RemoveInjectionPoint(id string) {
	ei.mu.Lock()
	defer ei.mu.Unlock()
	delete(ei.injectionPoints, id)
}

// ShouldInjectError determines if an error should be injected for the given operation
func (ei *ErrorInjector) ShouldInjectError(operation, path string) (*InjectedError, bool) {
	ei.mu.Lock()
	defer ei.mu.Unlock()

	if !ei.isEnabled {
		return nil, false
	}

	// Increment operation count
	key := fmt.Sprintf("%s:%s", operation, path)
	ei.operationCounts[key]++

	// Check all injection points
	for id, point := range ei.injectionPoints {
		if ei.matchesInjectionPoint(point, operation, path) {
			// Check trigger count
			if point.TriggerCount > 0 && ei.operationCounts[key] != point.TriggerCount {
				continue
			}

			// Check max injections
			injectionKey := fmt.Sprintf("%s:%s", id, key)
			if point.MaxInjections > 0 && ei.injectionCounts[injectionKey] >= point.MaxInjections {
				continue
			}

			// Check probability
			if ei.shouldInjectBasedOnProbability(point.Probability) {
				// Apply delays
				if point.DelayBefore > 0 {
					time.Sleep(point.DelayBefore)
				}

				// Increment injection count
				ei.injectionCounts[injectionKey]++

				// Create error with enhanced metadata
				injectedError := &InjectedError{
					Type:       point.Error.Type,
					Category:   point.Error.Category,
					Message:    point.Error.Message,
					Operation:  operation,
					Path:       path,
					Retryable:  point.Error.Retryable,
					RetryAfter: point.Error.RetryAfter,
					Cause:      point.Error.Cause,
					InjectedAt: time.Now(),
				}

				// Update statistics
				ei.stats.TotalInjections++
				ei.stats.InjectionsByType[injectedError.Type]++
				ei.stats.InjectionsByCategory[injectedError.Category]++
				ei.stats.OperationsAffected[operation]++

				// Start propagation trace
				ei.startPropagationTrace(injectedError)

				// Apply delay after injection
				if point.DelayAfter > 0 {
					time.Sleep(point.DelayAfter)
				}

				return injectedError, true
			}
		}
	}

	return nil, false
}

// TrackErrorPropagation records error flow through operations
func (ei *ErrorInjector) TrackErrorPropagation(errorID, operation, component, action string, errorType ErrorType, message string) {
	ei.mu.Lock()
	defer ei.mu.Unlock()

	trace, exists := ei.propagationTraces[errorID]
	if !exists {
		return
	}

	entry := ErrorPropagationEntry{
		Timestamp:  time.Now(),
		Operation:  operation,
		Component:  component,
		Action:     action,
		ErrorType:  errorType,
		Message:    message,
		StackDepth: ei.getCurrentStackDepth(),
	}

	trace.PropagationLog = append(trace.PropagationLog, entry)
}

// TrackRecoveryAttempt records error recovery testing
func (ei *ErrorInjector) TrackRecoveryAttempt(errorID, recoveryType string, attemptNumber int, success bool, duration time.Duration) {
	ei.mu.Lock()
	defer ei.mu.Unlock()

	attempt := ErrorRecoveryAttempt{
		ErrorID:       errorID,
		RecoveryType:  recoveryType,
		AttemptNumber: attemptNumber,
		Success:       success,
		Duration:      duration,
		Timestamp:     time.Now(),
	}

	ei.recoveryAttempts = append(ei.recoveryAttempts, attempt)

	if success {
		ei.stats.SuccessfulRecoveries++
	} else {
		ei.stats.FailedRecoveries++
	}

	// Update average recovery time
	totalTime := time.Duration(0)
	successCount := 0
	for _, ra := range ei.recoveryAttempts {
		if ra.Success {
			totalTime += ra.Duration
			successCount++
		}
	}
	if successCount > 0 {
		ei.stats.AverageRecoveryTime = totalTime / time.Duration(successCount)
	}
}

// GetStats returns error injection statistics
func (ei *ErrorInjector) GetStats() ErrorInjectionStats {
	ei.mu.RLock()
	defer ei.mu.RUnlock()
	return ei.stats
}

// GetPropagationTraces returns all propagation traces
func (ei *ErrorInjector) GetPropagationTraces() map[string]*ErrorPropagationTrace {
	ei.mu.RLock()
	defer ei.mu.RUnlock()

	traces := make(map[string]*ErrorPropagationTrace)
	for k, v := range ei.propagationTraces {
		traces[k] = v
	}
	return traces
}

// GetRecoveryAttempts returns all recovery attempts
func (ei *ErrorInjector) GetRecoveryAttempts() []ErrorRecoveryAttempt {
	ei.mu.RLock()
	defer ei.mu.RUnlock()

	attempts := make([]ErrorRecoveryAttempt, len(ei.recoveryAttempts))
	copy(attempts, ei.recoveryAttempts)
	return attempts
}

// Enable enables or disables error injection
func (ei *ErrorInjector) Enable(enabled bool) {
	ei.mu.Lock()
	defer ei.mu.Unlock()
	ei.isEnabled = enabled
}

// IsEnabled returns whether error injection is enabled
func (ei *ErrorInjector) IsEnabled() bool {
	ei.mu.RLock()
	defer ei.mu.RUnlock()
	return ei.isEnabled
}

// Reset clears all injection points and statistics
func (ei *ErrorInjector) Reset() {
	ei.mu.Lock()
	defer ei.mu.Unlock()

	ei.injectionPoints = make(map[string]*ErrorInjectionPoint)
	ei.operationCounts = make(map[string]int)
	ei.injectionCounts = make(map[string]int)
	ei.propagationTraces = make(map[string]*ErrorPropagationTrace)
	ei.recoveryAttempts = make([]ErrorRecoveryAttempt, 0)
	ei.stats = ErrorInjectionStats{
		InjectionsByType:     make(map[ErrorType]int),
		InjectionsByCategory: make(map[ErrorCategory]int),
		OperationsAffected:   make(map[string]int),
	}
}

// Helper methods

func (ei *ErrorInjector) matchesInjectionPoint(point *ErrorInjectionPoint, operation, path string) bool {
	// Check operation pattern
	if point.Operation != "" && point.Operation != "*" {
		if !strings.Contains(operation, point.Operation) {
			return false
		}
	}

	// Check path pattern
	if point.Path != "" && point.Path != "*" {
		if !strings.Contains(path, point.Path) {
			return false
		}
	}

	// Check custom condition
	if point.ConditionalFunc != nil {
		return point.ConditionalFunc(operation, path)
	}

	return true
}

func (ei *ErrorInjector) shouldInjectBasedOnProbability(probability float64) bool {
	// For testing, we'll use a simple deterministic approach
	// In real usage, this would use proper randomization
	return probability >= 1.0
}

func (ei *ErrorInjector) startPropagationTrace(err *InjectedError) {
	errorID := fmt.Sprintf("err_%d", time.Now().UnixNano())

	trace := &ErrorPropagationTrace{
		ErrorID:        errorID,
		OriginalError:  err,
		PropagationLog: make([]ErrorPropagationEntry, 0),
		StartTime:      time.Now(),
	}

	ei.propagationTraces[errorID] = trace
	ei.stats.PropagationTraces++
}

func (ei *ErrorInjector) getCurrentStackDepth() int {
	// Simplified stack depth calculation
	// In real implementation, would use runtime.Callers
	return 3
}

func getDefaultRetryPolicies() map[ErrorType]RetryPolicy {
	return map[ErrorType]RetryPolicy{
		ErrorTypeTransient: {
			MaxAttempts:    3,
			BaseDelay:      100 * time.Millisecond,
			MaxDelay:       5 * time.Second,
			BackoffFactor:  2.0,
			RetryableTypes: []ErrorType{ErrorTypeTransient},
		},
		ErrorTypeRecoverable: {
			MaxAttempts:    2,
			BaseDelay:      500 * time.Millisecond,
			MaxDelay:       10 * time.Second,
			BackoffFactor:  1.5,
			RetryableTypes: []ErrorType{ErrorTypeRecoverable},
		},
	}
}

// Pre-defined error creators for common scenarios

// CreateFilesystemError creates a filesystem-related injected error
func CreateFilesystemError(errorType ErrorType, message, operation, path string, cause error) *InjectedError {
	return &InjectedError{
		Type:      errorType,
		Category:  ErrorCategoryFilesystem,
		Message:   message,
		Operation: operation,
		Path:      path,
		Cause:     cause,
		Retryable: errorType == ErrorTypeTransient || errorType == ErrorTypeRecoverable,
	}
}

// CreateGitError creates a Git-related injected error
func CreateGitError(errorType ErrorType, message, operation, path string, cause error) *InjectedError {
	return &InjectedError{
		Type:      errorType,
		Category:  ErrorCategoryGit,
		Message:   message,
		Operation: operation,
		Path:      path,
		Cause:     cause,
		Retryable: errorType == ErrorTypeTransient,
	}
}

// CreatePermissionError creates a permission-related injected error
func CreatePermissionError(errorType ErrorType, message, operation, path string) *InjectedError {
	return &InjectedError{
		Type:      errorType,
		Category:  ErrorCategoryPermission,
		Message:   message,
		Operation: operation,
		Path:      path,
		Cause:     syscall.EACCES,
		Retryable: false, // Permission errors typically aren't retryable
	}
}

// CreateResourceError creates a resource exhaustion injected error
func CreateResourceError(errorType ErrorType, message, operation, path string) *InjectedError {
	return &InjectedError{
		Type:      errorType,
		Category:  ErrorCategoryResource,
		Message:   message,
		Operation: operation,
		Path:      path,
		Cause:     syscall.ENOSPC,
		Retryable: errorType == ErrorTypeTransient,
	}
}

// CreateNetworkError creates a network-related injected error
func CreateNetworkError(errorType ErrorType, message, operation, path string, cause error) *InjectedError {
	return &InjectedError{
		Type:      errorType,
		Category:  ErrorCategoryNetwork,
		Message:   message,
		Operation: operation,
		Path:      path,
		Cause:     cause,
		Retryable: errorType == ErrorTypeTransient,
	}
}

// ErrorInjectionScenario defines complete error injection test scenarios
type ErrorInjectionScenario struct {
	Name            string
	Description     string
	InjectionPoints map[string]*ErrorInjectionPoint
	ExpectedResults ScenarioExpectedResults
	SetupFunc       func(*ErrorInjector) error
	CleanupFunc     func(*ErrorInjector) error
}

// ScenarioExpectedResults defines what to expect from a scenario
type ScenarioExpectedResults struct {
	ShouldInjectErrors       bool
	ExpectedInjectionCount   int
	ExpectedRecoveryCount    int
	ExpectedPropagationDepth int
	AllowedErrorTypes        []ErrorType
	RequiredErrorCategories  []ErrorCategory
}

// Built-in scenarios for common testing patterns

func GetFilesystemErrorScenario() ErrorInjectionScenario {
	return ErrorInjectionScenario{
		Name:        "filesystem_errors",
		Description: "Comprehensive filesystem error injection testing",
		InjectionPoints: map[string]*ErrorInjectionPoint{
			"write_permission_denied": {
				Operation:    "WriteFile",
				Path:         "*",
				TriggerCount: 1,
				Error:        CreatePermissionError(ErrorTypePermanent, "write permission denied", "WriteFile", ""),
				Probability:  1.0,
			},
			"disk_full": {
				Operation:    "WriteFile",
				Path:         "*",
				TriggerCount: 2,
				Error:        CreateResourceError(ErrorTypeTransient, "disk full", "WriteFile", ""),
				Probability:  1.0,
			},
		},
		ExpectedResults: ScenarioExpectedResults{
			ShouldInjectErrors:      true,
			ExpectedInjectionCount:  2,
			ExpectedRecoveryCount:   1, // Only transient disk full should be retryable
			AllowedErrorTypes:       []ErrorType{ErrorTypePermanent, ErrorTypeTransient},
			RequiredErrorCategories: []ErrorCategory{ErrorCategoryPermission, ErrorCategoryResource},
		},
	}
}

func GetGitErrorScenario() ErrorInjectionScenario {
	return ErrorInjectionScenario{
		Name:        "git_errors",
		Description: "Git operation error injection testing",
		InjectionPoints: map[string]*ErrorInjectionPoint{
			"git_command_failed": {
				Operation:    "ExecuteCommand",
				Path:         "*",
				TriggerCount: 1,
				Error:        CreateGitError(ErrorTypeTransient, "git command failed", "ExecuteCommand", "", fmt.Errorf("command execution failed")),
				Probability:  1.0,
			},
			"repository_not_found": {
				Operation:    "IsRepository",
				Path:         "*",
				TriggerCount: 1,
				Error:        CreateGitError(ErrorTypePermanent, "not a git repository", "IsRepository", "", fmt.Errorf("fatal: not a git repository")),
				Probability:  1.0,
			},
		},
		ExpectedResults: ScenarioExpectedResults{
			ShouldInjectErrors:      true,
			ExpectedInjectionCount:  2,
			ExpectedRecoveryCount:   1, // Only transient command failure should be retryable
			AllowedErrorTypes:       []ErrorType{ErrorTypePermanent, ErrorTypeTransient},
			RequiredErrorCategories: []ErrorCategory{ErrorCategoryGit},
		},
	}
}

// RunErrorInjectionScenario executes a complete error injection scenario
func RunErrorInjectionScenario(injector *ErrorInjector, scenario ErrorInjectionScenario) error {
	// Setup scenario
	if scenario.SetupFunc != nil {
		if err := scenario.SetupFunc(injector); err != nil {
			return fmt.Errorf("scenario setup failed: %w", err)
		}
	}

	// Add injection points
	for id, point := range scenario.InjectionPoints {
		injector.AddInjectionPoint(id, point)
	}

	// Cleanup on exit
	defer func() {
		if scenario.CleanupFunc != nil {
			scenario.CleanupFunc(injector)
		}
	}()

	return nil
}
