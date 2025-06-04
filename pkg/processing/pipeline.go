// ⭐ EXTRACT-007: Processing pipelines - Pipeline patterns extracted from context-aware functions - 🔧
package processing

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// PipelineInterface defines the interface for processing workflows
type PipelineInterface interface {
	Execute(ctx context.Context, input *ProcessingInput) (*ProcessingResult, error)
	AddStage(stage PipelineStage)
	GetProgress() *PipelineProgress
	GetStages() []PipelineStage
}

// PipelineStage represents a single stage in a processing pipeline
type PipelineStage interface {
	GetName() string
	GetDescription() string
	Execute(ctx context.Context, input *ProcessingInput, output *ProcessingResult) error
	CanSkip(input *ProcessingInput) bool
	GetEstimatedDuration() time.Duration
}

// PipelineProgress represents the progress of a pipeline execution
type PipelineProgress struct {
	TotalStages     int           `json:"total_stages"`
	CompletedStages int           `json:"completed_stages"`
	CurrentStage    string        `json:"current_stage"`
	OverallProgress float64       `json:"overall_progress"`
	StageProgress   float64       `json:"stage_progress"`
	EstimatedTime   time.Duration `json:"estimated_time"`
	ElapsedTime     time.Duration `json:"elapsed_time"`
	RemainingTime   time.Duration `json:"remaining_time"`
	StartedAt       time.Time     `json:"started_at"`
	LastUpdate      time.Time     `json:"last_update"`
}

// PipelineResult extends ProcessingResult with pipeline-specific information
type PipelineResult struct {
	*ProcessingResult
	StageResults     map[string]*StageResult `json:"stage_results"`
	SkippedStages    []string                `json:"skipped_stages"`
	FailedStage      string                  `json:"failed_stage,omitempty"`
	PipelineDuration time.Duration           `json:"pipeline_duration"`
}

// StageResult represents the result of a single pipeline stage
type StageResult struct {
	Name        string                 `json:"name"`
	Duration    time.Duration          `json:"duration"`
	Success     bool                   `json:"success"`
	Skipped     bool                   `json:"skipped"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt time.Time              `json:"completed_at"`
}

// Pipeline implements context-aware processing workflows
type Pipeline struct {
	name     string
	stages   []PipelineStage
	progress *PipelineProgress
	mutex    sync.RWMutex

	// Configuration
	stopOnError     bool
	enableRollback  bool
	maxStageRetries int
	retryDelay      time.Duration

	// Progress tracking
	progressCallback func(*PipelineProgress)

	// Atomic counters for thread-safe progress updates
	completedStages int64
	totalItems      int64
	processedItems  int64
}

// NewPipeline creates a new processing pipeline
func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		name:            name,
		stages:          make([]PipelineStage, 0),
		progress:        &PipelineProgress{},
		stopOnError:     true,
		enableRollback:  false,
		maxStageRetries: 3,
		retryDelay:      time.Second,
	}
}

// AddStage adds a stage to the pipeline
func (p *Pipeline) AddStage(stage PipelineStage) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.stages = append(p.stages, stage)
	p.progress.TotalStages = len(p.stages)
}

// GetStages returns all pipeline stages
func (p *Pipeline) GetStages() []PipelineStage {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	stages := make([]PipelineStage, len(p.stages))
	copy(stages, p.stages)
	return stages
}

// SetStopOnError configures whether to stop on first error
func (p *Pipeline) SetStopOnError(stop bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.stopOnError = stop
}

// SetProgressCallback sets a callback for progress updates
func (p *Pipeline) SetProgressCallback(callback func(*PipelineProgress)) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.progressCallback = callback
}

// Execute runs the pipeline with context support (extracted from context-aware patterns)
func (p *Pipeline) Execute(ctx context.Context, input *ProcessingInput) (*ProcessingResult, error) {
	start := time.Now()

	// Initialize pipeline execution
	p.initializeExecution(start)

	// Create pipeline result
	pipelineResult := &PipelineResult{
		ProcessingResult: &ProcessingResult{
			Statistics: make(map[string]int64),
			Errors:     []string{},
			Warnings:   []string{},
		},
		StageResults:  make(map[string]*StageResult),
		SkippedStages: []string{},
	}

	// Execute stages
	for i, stage := range p.stages {
		select {
		case <-ctx.Done():
			return p.handleCancellation(pipelineResult, ctx.Err())
		default:
			// Continue execution
		}

		// Update current stage
		p.updateCurrentStage(stage.GetName(), i)

		// Check if stage can be skipped
		if stage.CanSkip(input) {
			p.handleSkippedStage(stage, pipelineResult)
			continue
		}

		// Execute stage with retries
		stageResult := p.executeStageWithRetries(ctx, stage, input, pipelineResult.ProcessingResult)
		pipelineResult.StageResults[stage.GetName()] = stageResult

		// Handle stage failure
		if !stageResult.Success {
			pipelineResult.FailedStage = stage.GetName()
			if p.stopOnError {
				break
			}
		}

		// Update progress
		atomic.AddInt64(&p.completedStages, 1)
		p.updateProgress()
	}

	// Finalize execution
	p.finalizeExecution(pipelineResult, time.Since(start))

	return pipelineResult.ProcessingResult, nil
}

// GetProgress returns the current pipeline progress
func (p *Pipeline) GetProgress() *PipelineProgress {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	// Create a copy to avoid race conditions
	progress := *p.progress
	return &progress
}

// initializeExecution sets up the pipeline for execution
func (p *Pipeline) initializeExecution(startTime time.Time) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	atomic.StoreInt64(&p.completedStages, 0)
	atomic.StoreInt64(&p.processedItems, 0)

	p.progress = &PipelineProgress{
		TotalStages:     len(p.stages),
		CompletedStages: 0,
		CurrentStage:    "",
		OverallProgress: 0.0,
		StageProgress:   0.0,
		StartedAt:       startTime,
		LastUpdate:      startTime,
	}
}

// executeStageWithRetries executes a stage with retry logic
func (p *Pipeline) executeStageWithRetries(ctx context.Context, stage PipelineStage, input *ProcessingInput, output *ProcessingResult) *StageResult {
	stageResult := &StageResult{
		Name:      stage.GetName(),
		StartedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	var lastErr error

	// Retry loop
	for attempt := 0; attempt <= p.maxStageRetries; attempt++ {
		select {
		case <-ctx.Done():
			stageResult.Error = fmt.Sprintf("stage cancelled: %v", ctx.Err())
			stageResult.Success = false
			stageResult.CompletedAt = time.Now()
			stageResult.Duration = time.Since(stageResult.StartedAt)
			return stageResult
		default:
			// Continue with execution
		}

		// Execute stage
		err := stage.Execute(ctx, input, output)
		if err == nil {
			// Success
			stageResult.Success = true
			break
		}

		lastErr = err

		// Handle retry
		if attempt < p.maxStageRetries {
			select {
			case <-ctx.Done():
				stageResult.Error = fmt.Sprintf("stage cancelled during retry: %v", ctx.Err())
				stageResult.Success = false
				stageResult.CompletedAt = time.Now()
				stageResult.Duration = time.Since(stageResult.StartedAt)
				return stageResult
			case <-time.After(p.retryDelay):
				// Continue to next attempt
			}
		}
	}

	// Set final result
	if lastErr != nil {
		stageResult.Error = lastErr.Error()
		stageResult.Success = false
	}

	stageResult.CompletedAt = time.Now()
	stageResult.Duration = time.Since(stageResult.StartedAt)

	return stageResult
}

// updateCurrentStage updates the current stage in progress
func (p *Pipeline) updateCurrentStage(stageName string, stageIndex int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.progress.CurrentStage = stageName
	p.progress.CompletedStages = stageIndex
	p.progress.LastUpdate = time.Now()

	if p.progressCallback != nil {
		go p.progressCallback(p.progress)
	}
}

// updateProgress calculates and updates overall progress
func (p *Pipeline) updateProgress() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	completed := atomic.LoadInt64(&p.completedStages)
	total := int64(len(p.stages))

	if total > 0 {
		p.progress.OverallProgress = float64(completed) / float64(total)
	}

	p.progress.ElapsedTime = time.Since(p.progress.StartedAt)
	p.progress.LastUpdate = time.Now()

	// Estimate remaining time
	if p.progress.OverallProgress > 0 {
		estimatedTotal := time.Duration(float64(p.progress.ElapsedTime) / p.progress.OverallProgress)
		p.progress.RemainingTime = estimatedTotal - p.progress.ElapsedTime
	}

	if p.progressCallback != nil {
		go p.progressCallback(p.progress)
	}
}

// handleSkippedStage handles a skipped stage
func (p *Pipeline) handleSkippedStage(stage PipelineStage, result *PipelineResult) {
	stageName := stage.GetName()
	result.SkippedStages = append(result.SkippedStages, stageName)

	stageResult := &StageResult{
		Name:        stageName,
		Success:     true,
		Skipped:     true,
		StartedAt:   time.Now(),
		CompletedAt: time.Now(),
		Duration:    0,
		Metadata:    map[string]interface{}{"reason": "stage can be skipped"},
	}

	result.StageResults[stageName] = stageResult
	atomic.AddInt64(&p.completedStages, 1)
	p.updateProgress()
}

// handleCancellation handles pipeline cancellation
func (p *Pipeline) handleCancellation(result *PipelineResult, err error) (*ProcessingResult, error) {
	result.ProcessingResult.Errors = append(result.ProcessingResult.Errors, fmt.Sprintf("pipeline cancelled: %v", err))
	return result.ProcessingResult, NewProcessingError("PIPELINE_CANCELLED", "Execute", fmt.Sprintf("pipeline execution cancelled: %v", err))
}

// finalizeExecution completes the pipeline execution
func (p *Pipeline) finalizeExecution(result *PipelineResult, duration time.Duration) {
	result.PipelineDuration = duration
	result.ProcessingResult.Duration = duration

	// Count successful stages
	successCount := 0
	for _, stageResult := range result.StageResults {
		if stageResult.Success {
			successCount++
		}
	}

	result.ProcessingResult.Statistics["total_stages"] = int64(len(p.stages))
	result.ProcessingResult.Statistics["successful_stages"] = int64(successCount)
	result.ProcessingResult.Statistics["skipped_stages"] = int64(len(result.SkippedStages))
	result.ProcessingResult.Statistics["failed_stages"] = int64(len(p.stages) - successCount)

	p.updateProgress()
}

// BaseStage provides common functionality for pipeline stages
type BaseStage struct {
	name              string
	description       string
	estimatedDuration time.Duration
	skipCondition     func(*ProcessingInput) bool
}

// NewBaseStage creates a new base stage
func NewBaseStage(name, description string, estimatedDuration time.Duration) *BaseStage {
	return &BaseStage{
		name:              name,
		description:       description,
		estimatedDuration: estimatedDuration,
	}
}

// GetName returns the stage name
func (bs *BaseStage) GetName() string {
	return bs.name
}

// GetDescription returns the stage description
func (bs *BaseStage) GetDescription() string {
	return bs.description
}

// GetEstimatedDuration returns the estimated duration for this stage
func (bs *BaseStage) GetEstimatedDuration() time.Duration {
	return bs.estimatedDuration
}

// CanSkip checks if this stage can be skipped
func (bs *BaseStage) CanSkip(input *ProcessingInput) bool {
	if bs.skipCondition != nil {
		return bs.skipCondition(input)
	}
	return false
}

// SetSkipCondition sets a condition function for skipping this stage
func (bs *BaseStage) SetSkipCondition(condition func(*ProcessingInput) bool) {
	bs.skipCondition = condition
}
