// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically context cancellation testing helpers for reliable testing
// of context-aware operations.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-D: Context cancellation testing helpers - 🔧
// DECISION-REF: DEC-007 (Context-aware operations)
// IMPLEMENTATION-NOTES: Use ticker-based timing control and goroutine coordination for deterministic cancellation testing

package testutil

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// CancellationPoint represents a specific point in an operation where cancellation can be injected
type CancellationPoint struct {
	ID               string
	Name             string
	Description      string
	OperationStage   string
	CancellationFunc func(ctx context.Context) error
	ExecutionCount   int64
	InjectionEnabled bool
}

// ContextController provides controlled timing and cancellation for context testing
type ContextController struct {
	mu                sync.RWMutex
	baseContext       context.Context
	cancelFunc        context.CancelFunc
	timeout           time.Duration
	cancellationDelay time.Duration
	ticker            *time.Ticker
	stopTicker        chan bool
	isActive          bool
	events            []ContextEvent
	stats             ContextTestStats
}

// ContextEvent represents an event during context testing
type ContextEvent struct {
	Timestamp   time.Time
	Type        ContextEventType
	Description string
	PointID     string
	Error       error
}

// ContextEventType defines the type of context event
type ContextEventType int

// Context event types for tracking operations and cancellations
const (
	// EventOperationStart indicates an operation has started
	EventOperationStart ContextEventType = iota
	// EventCancellationTriggered indicates cancellation was triggered
	EventCancellationTriggered
	// EventTimeoutTriggered indicates a timeout occurred
	EventTimeoutTriggered
	// EventOperationComplete indicates an operation completed
	EventOperationComplete
	// EventOperationCancelled indicates an operation was cancelled
	EventOperationCancelled
	// EventPropagationVerified indicates propagation was verified
	EventPropagationVerified
	// EventPropagationFailed indicates propagation verification failed
	EventPropagationFailed
)

// ContextTestStats tracks statistics for context testing
type ContextTestStats struct {
	TotalOperations      int64
	CancelledOperations  int64
	TimeoutOperations    int64
	SuccessfulOperations int64
	FailedOperations     int64
	PropagationTests     int64
	PropagationFailures  int64
	AverageDuration      time.Duration
	TotalDuration        time.Duration
}

// CancellationManager manages cancellation point injection and timing
type CancellationManager struct {
	mu                 sync.RWMutex
	cancellationPoints map[string]*CancellationPoint
	controllers        map[string]*ContextController
	globalEnabled      bool
	stats              CancellationStats
}

// CancellationStats tracks cancellation testing statistics
type CancellationStats struct {
	TotalInjections      int64
	SuccessfulInjections int64
	FailedInjections     int64
	PointsRegistered     int64
	ControllersActive    int64
}

// ConcurrentTestConfig configures concurrent operation testing
type ConcurrentTestConfig struct {
	NumOperations       int
	MaxConcurrency      int
	OperationTimeout    time.Duration
	CancellationDelay   time.Duration
	StaggerStart        time.Duration
	EnableRandomization bool
	TestDuration        time.Duration
}

// ConcurrentTestResult contains results from concurrent testing
type ConcurrentTestResult struct {
	TotalOperations     int
	CompletedOperations int
	CancelledOperations int
	TimeoutOperations   int
	FailedOperations    int
	AverageDuration     time.Duration
	MaxDuration         time.Duration
	MinDuration         time.Duration
	ConcurrencyAchieved int
	ResourceLeaks       int
	PropagationFailures int
}

// PropagationTestConfig configures context propagation testing
type PropagationTestConfig struct {
	ChainDepth         int
	PropagationDelay   time.Duration
	VerificationPoints []string
	TimeoutDuration    time.Duration
	EnableDeepTrace    bool
}

// PropagationChain represents a chain of operations for propagation testing
type PropagationChain struct {
	ID             string
	Depth          int
	Operations     []PropagationOperation
	StartTime      time.Time
	EndTime        time.Time
	PropagationLog []PropagationEvent
	Success        bool
	Error          error
}

// PropagationOperation represents an operation in the propagation chain
type PropagationOperation struct {
	ID         string
	Name       string
	Depth      int
	StartTime  time.Time
	EndTime    time.Time
	Context    context.Context
	Success    bool
	Error      error
	Propagated bool
}

// PropagationEvent tracks propagation events
type PropagationEvent struct {
	Timestamp time.Time
	Depth     int
	Operation string
	EventType string
	Message   string
	Error     error
}

// NewContextController creates a new context controller for testing
func NewContextController(timeout time.Duration) *ContextController {
	ctx, cancel := context.WithCancel(context.Background())

	controller := &ContextController{
		baseContext:       ctx,
		cancelFunc:        cancel,
		timeout:           timeout,
		cancellationDelay: 0,
		stopTicker:        make(chan bool, 1),
		isActive:          false,
		events:            make([]ContextEvent, 0),
		stats:             ContextTestStats{},
	}

	return controller
}

// SetCancellationDelay sets the delay before cancellation is triggered
func (cc *ContextController) SetCancellationDelay(delay time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.cancellationDelay = delay
}

// StartControlledCancellation starts controlled cancellation with timing
func (cc *ContextController) StartControlledCancellation() context.Context {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if cc.isActive {
		return cc.baseContext
	}

	cc.isActive = true

	// Create context with timeout if specified
	var ctx context.Context
	if cc.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(cc.baseContext, cc.timeout)
		// Store the cancel function so it can be called in Stop()
		cc.cancelFunc = cancel
	} else {
		ctx = cc.baseContext
	}

	// Start cancellation timer if delay is specified
	if cc.cancellationDelay > 0 {
		cc.ticker = time.NewTicker(cc.cancellationDelay)
		go cc.cancellationTimer()
	}

	cc.recordEvent(EventOperationStart, "Controlled cancellation started", "", nil)
	atomic.AddInt64(&cc.stats.TotalOperations, 1)

	return ctx
}

// cancellationTimer handles delayed cancellation
func (cc *ContextController) cancellationTimer() {
	select {
	case <-cc.ticker.C:
		cc.triggerCancellation()
	case <-cc.stopTicker:
		cc.ticker.Stop()
		return
	}
}

// triggerCancellation triggers the cancellation
func (cc *ContextController) triggerCancellation() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if cc.cancelFunc != nil {
		cc.cancelFunc()
		cc.recordEvent(EventCancellationTriggered, "Cancellation triggered by timer", "", nil)
		atomic.AddInt64(&cc.stats.CancelledOperations, 1)
	}
}

// Stop stops the controller and cleans up resources
func (cc *ContextController) Stop() {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if cc.ticker != nil {
		cc.stopTicker <- true
		cc.ticker.Stop()
	}

	if cc.cancelFunc != nil {
		cc.cancelFunc()
	}

	cc.isActive = false
	cc.recordEvent(EventOperationComplete, "Controller stopped", "", nil)
}

// recordEvent records a context event
func (cc *ContextController) recordEvent(eventType ContextEventType, description, pointID string, err error) {
	event := ContextEvent{
		Timestamp:   time.Now(),
		Type:        eventType,
		Description: description,
		PointID:     pointID,
		Error:       err,
	}
	cc.events = append(cc.events, event)
}

// GetEvents returns all recorded events
func (cc *ContextController) GetEvents() []ContextEvent {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	events := make([]ContextEvent, len(cc.events))
	copy(events, cc.events)
	return events
}

// GetStats returns current statistics
func (cc *ContextController) GetStats() ContextTestStats {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	stats := cc.stats
	stats.TotalOperations = atomic.LoadInt64(&cc.stats.TotalOperations)
	stats.CancelledOperations = atomic.LoadInt64(&cc.stats.CancelledOperations)
	stats.TimeoutOperations = atomic.LoadInt64(&cc.stats.TimeoutOperations)
	stats.SuccessfulOperations = atomic.LoadInt64(&cc.stats.SuccessfulOperations)
	stats.FailedOperations = atomic.LoadInt64(&cc.stats.FailedOperations)
	stats.PropagationTests = atomic.LoadInt64(&cc.stats.PropagationTests)
	stats.PropagationFailures = atomic.LoadInt64(&cc.stats.PropagationFailures)

	return stats
}

// NewCancellationManager creates a new cancellation manager
func NewCancellationManager() *CancellationManager {
	return &CancellationManager{
		cancellationPoints: make(map[string]*CancellationPoint),
		controllers:        make(map[string]*ContextController),
		globalEnabled:      true,
		stats:              CancellationStats{},
	}
}

// RegisterCancellationPoint registers a new cancellation point
func (cm *CancellationManager) RegisterCancellationPoint(point *CancellationPoint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cancellationPoints[point.ID] = point
	atomic.AddInt64(&cm.stats.PointsRegistered, 1)
}

// InjectCancellation injects cancellation at a specific point
func (cm *CancellationManager) InjectCancellation(pointID string, ctx context.Context) error {
	cm.mu.RLock()
	point, exists := cm.cancellationPoints[pointID]
	cm.mu.RUnlock()

	if !exists || !cm.globalEnabled || !point.InjectionEnabled {
		return nil
	}

	atomic.AddInt64(&point.ExecutionCount, 1)
	atomic.AddInt64(&cm.stats.TotalInjections, 1)

	if point.CancellationFunc != nil {
		err := point.CancellationFunc(ctx)
		if err != nil {
			atomic.AddInt64(&cm.stats.FailedInjections, 1)
			return err
		}
		atomic.AddInt64(&cm.stats.SuccessfulInjections, 1)
	}

	return ctx.Err() // Return cancellation error if context is cancelled
}

// EnableGlobal enables or disables global cancellation injection
func (cm *CancellationManager) EnableGlobal(enabled bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.globalEnabled = enabled
}

// EnablePoint enables or disables a specific cancellation point
func (cm *CancellationManager) EnablePoint(pointID string, enabled bool) {
	cm.mu.RLock()
	point, exists := cm.cancellationPoints[pointID]
	cm.mu.RUnlock()

	if exists {
		point.InjectionEnabled = enabled
	}
}

// GetStats returns cancellation manager statistics
func (cm *CancellationManager) GetStats() CancellationStats {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	stats := cm.stats
	stats.TotalInjections = atomic.LoadInt64(&cm.stats.TotalInjections)
	stats.SuccessfulInjections = atomic.LoadInt64(&cm.stats.SuccessfulInjections)
	stats.FailedInjections = atomic.LoadInt64(&cm.stats.FailedInjections)
	stats.PointsRegistered = atomic.LoadInt64(&cm.stats.PointsRegistered)
	stats.ControllersActive = atomic.LoadInt64(&cm.stats.ControllersActive)

	return stats
}

// RunConcurrentTest runs a concurrent operation test with cancellation
func (cm *CancellationManager) RunConcurrentTest(
	config ConcurrentTestConfig,
	operationFunc func(ctx context.Context) error,
) *ConcurrentTestResult {
	result := &ConcurrentTestResult{
		TotalOperations: config.NumOperations,
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxConcurrency)
	results := make(chan operationResult, config.NumOperations)

	// Start operations
	for i := 0; i < config.NumOperations; i++ {
		wg.Add(1)
		go func(opID int) {
			defer wg.Done()

			// Acquire semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Optional stagger start
			if config.StaggerStart > 0 {
				time.Sleep(time.Duration(opID) * config.StaggerStart)
			}

			// Create context with timeout
			ctx, cancel := context.WithTimeout(
				context.Background(),
				config.OperationTimeout,
			)
			defer cancel()

			// Optional cancellation delay
			if config.CancellationDelay > 0 {
				go func() {
					time.Sleep(config.CancellationDelay)
					cancel()
				}()
			}

			opStart := time.Now()
			err := operationFunc(ctx)
			duration := time.Since(opStart)

			results <- operationResult{
				ID:        opID,
				Duration:  duration,
				Error:     err,
				Cancelled: ctx.Err() == context.Canceled,
				Timeout:   ctx.Err() == context.DeadlineExceeded,
			}
		}(i)
	}

	// Close results channel when all operations complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	durations := make([]time.Duration, 0, config.NumOperations)

	for opResult := range results {
		durations = append(durations, opResult.Duration)

		if opResult.Error != nil {
			result.FailedOperations++
		} else {
			result.CompletedOperations++
		}

		if opResult.Cancelled {
			result.CancelledOperations++
		}

		if opResult.Timeout {
			result.TimeoutOperations++
		}
	}

	// Calculate statistics
	if len(durations) > 0 {
		var totalDuration time.Duration
		result.MinDuration = durations[0]
		result.MaxDuration = durations[0]

		for _, d := range durations {
			totalDuration += d
			if d < result.MinDuration {
				result.MinDuration = d
			}
			if d > result.MaxDuration {
				result.MaxDuration = d
			}
		}

		result.AverageDuration = totalDuration / time.Duration(len(durations))
	}

	result.ConcurrencyAchieved = config.MaxConcurrency

	return result
}

// operationResult represents the result of a single operation
type operationResult struct {
	ID        int
	Duration  time.Duration
	Error     error
	Cancelled bool
	Timeout   bool
}

// TestPropagation tests context propagation through a chain of operations
func (cm *CancellationManager) TestPropagation(config PropagationTestConfig) *PropagationChain {
	chain := &PropagationChain{
		ID:             fmt.Sprintf("chain-%d", time.Now().UnixNano()),
		Depth:          config.ChainDepth,
		Operations:     make([]PropagationOperation, 0, config.ChainDepth),
		StartTime:      time.Now(),
		PropagationLog: make([]PropagationEvent, 0),
		Success:        true,
	}

	// Create base context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.TimeoutDuration)
	defer cancel()

	// Start propagation chain
	chain.Error = cm.executePropagationChain(ctx, chain, config, 0)
	chain.EndTime = time.Now()

	if chain.Error != nil {
		chain.Success = false
	}

	return chain
}

// executePropagationChain recursively executes the propagation chain
func (cm *CancellationManager) executePropagationChain(ctx context.Context, chain *PropagationChain, config PropagationTestConfig, depth int) error {
	if depth >= config.ChainDepth {
		return nil
	}

	operation := PropagationOperation{
		ID:        fmt.Sprintf("op-%d-%d", depth, time.Now().UnixNano()),
		Name:      fmt.Sprintf("Operation Level %d", depth),
		Depth:     depth,
		StartTime: time.Now(),
		Context:   ctx,
		Success:   true,
	}

	// Check context at this level
	if err := ctx.Err(); err != nil {
		operation.Success = false
		operation.Error = err
		operation.Propagated = true
		chain.Operations = append(chain.Operations, operation)

		chain.PropagationLog = append(chain.PropagationLog, PropagationEvent{
			Timestamp: time.Now(),
			Depth:     depth,
			Operation: operation.Name,
			EventType: "cancellation_propagated",
			Message:   fmt.Sprintf("Context cancellation propagated to depth %d", depth),
			Error:     err,
		})

		return err
	}

	// Simulate work with propagation delay
	if config.PropagationDelay > 0 {
		select {
		case <-time.After(config.PropagationDelay):
			// Continue
		case <-ctx.Done():
			operation.Success = false
			operation.Error = ctx.Err()
			operation.Propagated = true
			operation.EndTime = time.Now()
			chain.Operations = append(chain.Operations, operation)
			return ctx.Err()
		}
	}

	// Continue to next level
	err := cm.executePropagationChain(ctx, chain, config, depth+1)
	operation.EndTime = time.Now()

	if err != nil {
		operation.Error = err
		operation.Success = false
		operation.Propagated = true
	}

	chain.Operations = append(chain.Operations, operation)
	return err
}

// CreateTimeoutContext creates a context with specified timeout for testing
func CreateTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// CreateCancelledContext creates a pre-cancelled context for testing
func CreateCancelledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

// CreateDeadlineContext creates a context with a deadline in the past for testing
func CreateDeadlineContext() context.Context {
	deadline := time.Now().Add(-1 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	// Call cancel immediately to avoid leak, since this is a test function
	// and the context is meant to be already expired
	cancel()
	return ctx
}

// SimulateSlowOperation simulates a slow operation that can be cancelled
func SimulateSlowOperation(ctx context.Context, duration time.Duration, checkInterval time.Duration) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Continue operation
		}
	}

	return nil
}

// VerifyContextPropagation verifies that context cancellation properly propagates
func VerifyContextPropagation(
	ctx context.Context,
	operationFunc func(context.Context) error,
) error {
	// Run operation and verify it respects context cancellation
	done := make(chan error, 1)

	go func() {
		done <- operationFunc(ctx)
	}()

	select {
	case err := <-done:
		if ctx.Err() != nil && err != ctx.Err() {
			return fmt.Errorf(
				"operation did not properly handle context cancellation: expected %v, got %v",
				ctx.Err(), err,
			)
		}
		return err
	case <-time.After(5 * time.Second):
		return fmt.Errorf(
			"operation did not complete or respect cancellation within timeout",
		)
	}
}

// Helper functions for common cancellation testing scenarios

// CreateArchiveCancellationPoint creates a cancellation point for archive operations
func CreateArchiveCancellationPoint(stage string) *CancellationPoint {
	return &CancellationPoint{
		ID:               fmt.Sprintf("archive_%s", stage),
		Name:             fmt.Sprintf("Archive %s", stage),
		Description:      fmt.Sprintf("Cancellation point during archive %s", stage),
		OperationStage:   stage,
		InjectionEnabled: true,
		CancellationFunc: func(ctx context.Context) error {
			return ctx.Err()
		},
	}
}

// CreateBackupCancellationPoint creates a cancellation point for backup operations
func CreateBackupCancellationPoint(stage string) *CancellationPoint {
	return &CancellationPoint{
		ID:               fmt.Sprintf("backup_%s", stage),
		Name:             fmt.Sprintf("Backup %s", stage),
		Description:      fmt.Sprintf("Cancellation point during backup %s", stage),
		OperationStage:   stage,
		InjectionEnabled: true,
		CancellationFunc: func(ctx context.Context) error {
			return ctx.Err()
		},
	}
}

// CreateResourceCleanupCancellationPoint creates a cancellation point for resource cleanup
func CreateResourceCleanupCancellationPoint() *CancellationPoint {
	return &CancellationPoint{
		ID:               "resource_cleanup",
		Name:             "Resource Cleanup",
		Description:      "Cancellation point during resource cleanup",
		OperationStage:   "cleanup",
		InjectionEnabled: true,
		CancellationFunc: func(ctx context.Context) error {
			return ctx.Err()
		},
	}
}
