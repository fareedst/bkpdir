// ⭐ EXTRACT-007: Concurrent processing - Worker pool patterns extracted from file collection functions - 🔧
package processing

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// ConcurrentProcessorInterface defines the interface for concurrent processing
type ConcurrentProcessorInterface interface {
	Process(ctx context.Context, items []ProcessingItem) (*ConcurrentResult, error)
	SetWorkerCount(count int)
	SetBatchSize(size int)
	GetStatus() *ConcurrentStatus
}

// ProcessingItem represents an item to be processed concurrently
type ProcessingItem struct {
	ID       string                 `json:"id"`
	Data     interface{}            `json:"data"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Priority int                    `json:"priority,omitempty"`
}

// ProcessingTask represents a processing task with context and cancellation
type ProcessingTask struct {
	Item      *ProcessingItem
	Context   context.Context
	Cancel    context.CancelFunc
	Result    chan *TaskResult
	StartTime time.Time
}

// TaskResult represents the result of processing a single item
type TaskResult struct {
	Item      *ProcessingItem        `json:"item"`
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Duration  time.Duration          `json:"duration"`
	Output    interface{}            `json:"output,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	WorkerID  int                    `json:"worker_id"`
	Timestamp time.Time              `json:"timestamp"`
}

// ConcurrentResult represents the result of concurrent processing
type ConcurrentResult struct {
	TotalItems      int              `json:"total_items"`
	SuccessfulItems int              `json:"successful_items"`
	FailedItems     int              `json:"failed_items"`
	Duration        time.Duration    `json:"duration"`
	WorkerCount     int              `json:"worker_count"`
	AverageItemTime time.Duration    `json:"average_item_time"`
	Results         []*TaskResult    `json:"results"`
	Errors          []string         `json:"errors,omitempty"`
	Statistics      map[string]int64 `json:"statistics"`
}

// ConcurrentStatus represents the status of concurrent processing
type ConcurrentStatus struct {
	State           ProcessingState `json:"state"`
	TotalItems      int             `json:"total_items"`
	ProcessedItems  int64           `json:"processed_items"`
	SuccessfulItems int64           `json:"successful_items"`
	FailedItems     int64           `json:"failed_items"`
	ActiveWorkers   int             `json:"active_workers"`
	QueuedItems     int             `json:"queued_items"`
	StartedAt       time.Time       `json:"started_at"`
	ElapsedTime     time.Duration   `json:"elapsed_time"`
	EstimatedTime   time.Duration   `json:"estimated_time"`
	Progress        float64         `json:"progress"`
}

// ConcurrentProcessor implements worker pool patterns for concurrent processing
type ConcurrentProcessor struct {
	// Configuration
	workerCount int
	batchSize   int

	// Processing function
	processFunc func(ctx context.Context, item *ProcessingItem) (interface{}, error)

	// Status and synchronization
	status   *ConcurrentStatus
	statusMu sync.RWMutex

	// Atomic counters
	processedItems  int64
	successfulItems int64
	failedItems     int64

	// Channels and workers
	taskQueue   chan *ProcessingTask
	resultQueue chan *TaskResult
	workers     []*Worker
	workerGroup sync.WaitGroup

	// Context management
	ctx    context.Context
	cancel context.CancelFunc
}

// Worker represents a single worker in the pool
type Worker struct {
	id          int
	processor   *ConcurrentProcessor
	taskQueue   <-chan *ProcessingTask
	resultQueue chan<- *TaskResult
	ctx         context.Context
	cancel      context.CancelFunc
	wg          *sync.WaitGroup
}

// NewConcurrentProcessor creates a new concurrent processor
func NewConcurrentProcessor(processFunc func(ctx context.Context, item *ProcessingItem) (interface{}, error)) *ConcurrentProcessor {
	workerCount := runtime.NumCPU()

	return &ConcurrentProcessor{
		workerCount: workerCount,
		batchSize:   100,
		processFunc: processFunc,
		status: &ConcurrentStatus{
			State: StateInitializing,
		},
		taskQueue:   make(chan *ProcessingTask, workerCount*2),
		resultQueue: make(chan *TaskResult, workerCount*2),
	}
}

// SetWorkerCount sets the number of workers
func (cp *ConcurrentProcessor) SetWorkerCount(count int) {
	if count <= 0 {
		count = runtime.NumCPU()
	}
	cp.workerCount = count
}

// SetBatchSize sets the batch size for processing
func (cp *ConcurrentProcessor) SetBatchSize(size int) {
	if size <= 0 {
		size = 100
	}
	cp.batchSize = size
}

// GetStatus returns the current processing status
func (cp *ConcurrentProcessor) GetStatus() *ConcurrentStatus {
	cp.statusMu.RLock()
	defer cp.statusMu.RUnlock()

	// Create a copy to avoid race conditions
	status := *cp.status
	status.ProcessedItems = atomic.LoadInt64(&cp.processedItems)
	status.SuccessfulItems = atomic.LoadInt64(&cp.successfulItems)
	status.FailedItems = atomic.LoadInt64(&cp.failedItems)

	if status.TotalItems > 0 {
		status.Progress = float64(status.ProcessedItems) / float64(status.TotalItems)
	}

	if !status.StartedAt.IsZero() {
		status.ElapsedTime = time.Since(status.StartedAt)

		// Estimate remaining time
		if status.Progress > 0 {
			estimatedTotal := time.Duration(float64(status.ElapsedTime) / status.Progress)
			status.EstimatedTime = estimatedTotal - status.ElapsedTime
		}
	}

	return &status
}

// Process processes items concurrently using worker pools (extracted from collectFilesToArchive patterns)
func (cp *ConcurrentProcessor) Process(ctx context.Context, items []ProcessingItem) (*ConcurrentResult, error) {
	start := time.Now()

	// Initialize processing
	cp.initializeProcessing(ctx, len(items), start)

	// Start workers
	err := cp.startWorkers()
	if err != nil {
		return nil, err
	}
	defer cp.stopWorkers()

	// Start result collector
	resultCollector := cp.startResultCollector()

	// Submit tasks and wait for completion
	taskSubmissionDone := make(chan struct{})
	go func() {
		defer close(taskSubmissionDone)
		cp.submitTasks(items)
	}()

	// Wait for task submission to complete
	<-taskSubmissionDone

	// Close task queue to signal no more tasks
	close(cp.taskQueue)

	// Wait for completion and collect results
	results := cp.waitForCompletion(resultCollector, len(items))

	// Create final result
	return cp.createFinalResult(results, time.Since(start)), nil
}

// initializeProcessing sets up the processor for execution
func (cp *ConcurrentProcessor) initializeProcessing(ctx context.Context, itemCount int, startTime time.Time) {
	cp.ctx, cp.cancel = context.WithCancel(ctx)

	// Reset counters
	atomic.StoreInt64(&cp.processedItems, 0)
	atomic.StoreInt64(&cp.successfulItems, 0)
	atomic.StoreInt64(&cp.failedItems, 0)

	// Update status
	cp.statusMu.Lock()
	cp.status.State = StateProcessing
	cp.status.TotalItems = itemCount
	cp.status.StartedAt = startTime
	cp.status.ActiveWorkers = cp.workerCount
	cp.statusMu.Unlock()
}

// startWorkers creates and starts the worker pool
func (cp *ConcurrentProcessor) startWorkers() error {
	cp.workers = make([]*Worker, cp.workerCount)

	for i := 0; i < cp.workerCount; i++ {
		worker := &Worker{
			id:          i,
			processor:   cp,
			taskQueue:   cp.taskQueue,
			resultQueue: cp.resultQueue,
			wg:          &cp.workerGroup,
		}

		worker.ctx, worker.cancel = context.WithCancel(cp.ctx)
		cp.workers[i] = worker

		cp.workerGroup.Add(1)
		go worker.run()
	}

	return nil
}

// stopWorkers stops all workers and cleans up resources
func (cp *ConcurrentProcessor) stopWorkers() {
	// Cancel all workers
	if cp.cancel != nil {
		cp.cancel()
	}

	// Wait for workers to finish
	cp.workerGroup.Wait()

	// Update status
	cp.statusMu.Lock()
	cp.status.State = StateCompleted
	cp.status.ActiveWorkers = 0
	cp.statusMu.Unlock()
}

// submitTasks submits items for processing
func (cp *ConcurrentProcessor) submitTasks(items []ProcessingItem) {
	defer func() {
		// Update queued items when done
		cp.statusMu.Lock()
		cp.status.QueuedItems = 0
		cp.statusMu.Unlock()
	}()

	for i := range items {
		select {
		case <-cp.ctx.Done():
			return
		default:
			// Create task
			taskCtx, taskCancel := context.WithCancel(cp.ctx)
			task := &ProcessingTask{
				Item:      &items[i],
				Context:   taskCtx,
				Cancel:    taskCancel,
				Result:    make(chan *TaskResult, 1),
				StartTime: time.Now(),
			}

			// Submit task
			select {
			case cp.taskQueue <- task:
				// Update queued items
				cp.statusMu.Lock()
				cp.status.QueuedItems++
				cp.statusMu.Unlock()
			case <-cp.ctx.Done():
				taskCancel()
				return
			}
		}
	}
}

// startResultCollector starts the result collection goroutine
func (cp *ConcurrentProcessor) startResultCollector() chan []*TaskResult {
	resultCollector := make(chan []*TaskResult, 1)

	go func() {
		defer close(resultCollector)

		var results []*TaskResult

		for result := range cp.resultQueue {
			results = append(results, result)

			// Update counters
			atomic.AddInt64(&cp.processedItems, 1)
			if result.Success {
				atomic.AddInt64(&cp.successfulItems, 1)
			} else {
				atomic.AddInt64(&cp.failedItems, 1)
			}

			// Update queued items
			cp.statusMu.Lock()
			if cp.status.QueuedItems > 0 {
				cp.status.QueuedItems--
			}
			cp.statusMu.Unlock()
		}

		resultCollector <- results
	}()

	return resultCollector
}

// waitForCompletion waits for all tasks to complete and returns results
func (cp *ConcurrentProcessor) waitForCompletion(resultCollector chan []*TaskResult, expectedCount int) []*TaskResult {
	// Wait for all workers to finish
	cp.workerGroup.Wait()

	// Close result queue so result collector can finish
	close(cp.resultQueue)

	// Get collected results
	results := <-resultCollector

	return results
}

// createFinalResult creates the final concurrent processing result
func (cp *ConcurrentProcessor) createFinalResult(results []*TaskResult, duration time.Duration) *ConcurrentResult {
	successCount := 0
	failCount := 0
	var totalItemTime time.Duration
	var errors []string

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failCount++
			if result.Error != "" {
				errors = append(errors, result.Error)
			}
		}
		totalItemTime += result.Duration
	}

	averageItemTime := time.Duration(0)
	if len(results) > 0 {
		averageItemTime = totalItemTime / time.Duration(len(results))
	}

	return &ConcurrentResult{
		TotalItems:      len(results),
		SuccessfulItems: successCount,
		FailedItems:     failCount,
		Duration:        duration,
		WorkerCount:     cp.workerCount,
		AverageItemTime: averageItemTime,
		Results:         results,
		Errors:          errors,
		Statistics: map[string]int64{
			"total_items":      int64(len(results)),
			"successful_items": int64(successCount),
			"failed_items":     int64(failCount),
			"worker_count":     int64(cp.workerCount),
		},
	}
}

// Worker implementation

// run executes the worker main loop
func (w *Worker) run() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		case task, ok := <-w.taskQueue:
			if !ok {
				return // Queue closed
			}

			// Process task
			result := w.processTask(task)

			// Send result
			select {
			case w.resultQueue <- result:
			case <-w.ctx.Done():
				return
			}
		}
	}
}

// processTask processes a single task
func (w *Worker) processTask(task *ProcessingTask) *TaskResult {
	start := time.Now()

	result := &TaskResult{
		Item:      task.Item,
		WorkerID:  w.id,
		Timestamp: start,
		Metadata:  make(map[string]interface{}),
	}

	// Check for cancellation
	select {
	case <-task.Context.Done():
		result.Success = false
		result.Error = fmt.Sprintf("task cancelled: %v", task.Context.Err())
		result.Duration = time.Since(start)
		return result
	default:
		// Continue processing
	}

	// Process item
	output, err := w.processor.processFunc(task.Context, task.Item)
	result.Duration = time.Since(start)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.Success = true
		result.Output = output
	}

	return result
}

// Utility functions for concurrent processing

// ProcessItems is a convenience function for processing items concurrently
func ProcessItems(ctx context.Context, items []ProcessingItem, processFunc func(ctx context.Context, item *ProcessingItem) (interface{}, error)) (*ConcurrentResult, error) {
	processor := NewConcurrentProcessor(processFunc)
	return processor.Process(ctx, items)
}

// ProcessItemsWithOptions processes items with custom options
func ProcessItemsWithOptions(ctx context.Context, items []ProcessingItem, processFunc func(ctx context.Context, item *ProcessingItem) (interface{}, error), workerCount, batchSize int) (*ConcurrentResult, error) {
	processor := NewConcurrentProcessor(processFunc)
	processor.SetWorkerCount(workerCount)
	processor.SetBatchSize(batchSize)
	return processor.Process(ctx, items)
}
