// ⭐ EXTRACT-007: Processing package structure - Core interfaces and types - 🔧
package processing

import (
	"context"
	"fmt"
	"io"
	"time"
)

// ProcessorInterface defines the main data processing interface
type ProcessorInterface interface {
	Process(ctx context.Context, input *ProcessingInput) (*ProcessingResult, error)
	Validate(input *ProcessingInput) error
	GetStatus() ProcessingStatus
}

// ProcessingInput represents input data for processing operations
type ProcessingInput struct {
	// Data source configuration
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Data        io.Reader `json:"-"`

	// Processing options
	Options  map[string]interface{} `json:"options"`
	Metadata map[string]string      `json:"metadata"`

	// Context and timing
	Timeout     time.Duration `json:"timeout"`
	Concurrency int           `json:"concurrency"`
}

// ProcessingResult represents the result of a processing operation
type ProcessingResult struct {
	// Output information
	Output     string    `json:"output"`
	OutputData io.Reader `json:"-"`

	// Processing statistics
	ItemsProcessed int              `json:"items_processed"`
	Duration       time.Duration    `json:"duration"`
	Statistics     map[string]int64 `json:"statistics"`

	// Verification and integrity
	Checksum string `json:"checksum,omitempty"`
	Verified bool   `json:"verified"`

	// Error handling
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// ProcessingStatus represents the current status of a processing operation
type ProcessingStatus struct {
	State          ProcessingState `json:"state"`
	Progress       float64         `json:"progress"`
	CurrentItem    string          `json:"current_item,omitempty"`
	ItemsRemaining int             `json:"items_remaining"`
	EstimatedTime  time.Duration   `json:"estimated_time"`
	StartedAt      time.Time       `json:"started_at"`
	LastUpdate     time.Time       `json:"last_update"`
}

// ProcessingState represents the state of a processing operation
type ProcessingState string

const (
	StateInitializing ProcessingState = "initializing"
	StateProcessing   ProcessingState = "processing"
	StateVerifying    ProcessingState = "verifying"
	StateCompleted    ProcessingState = "completed"
	StateFailed       ProcessingState = "failed"
	StateCancelled    ProcessingState = "cancelled"
)

// ProcessingError represents errors that occur during processing
type ProcessingError struct {
	Code        string    `json:"code"`
	Message     string    `json:"message"`
	Operation   string    `json:"operation"`
	Item        string    `json:"item,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Recoverable bool      `json:"recoverable"`
}

// Error implements the error interface
func (e *ProcessingError) Error() string {
	if e.Item != "" {
		return fmt.Sprintf("%s: %s (item: %s)", e.Operation, e.Message, e.Item)
	}
	return fmt.Sprintf("%s: %s", e.Operation, e.Message)
}

// NewProcessingError creates a new processing error
func NewProcessingError(code, operation, message string) *ProcessingError {
	return &ProcessingError{
		Code:      code,
		Operation: operation,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// ProcessingOptions provides common configuration for processors
type ProcessingOptions struct {
	// General options
	DryRun      bool          `json:"dry_run"`
	Verbose     bool          `json:"verbose"`
	Timeout     time.Duration `json:"timeout"`
	Concurrency int           `json:"concurrency"`

	// Verification options
	EnableVerification    bool   `json:"enable_verification"`
	VerificationAlgorithm string `json:"verification_algorithm"`

	// Error handling options
	ContinueOnError bool          `json:"continue_on_error"`
	MaxRetries      int           `json:"max_retries"`
	RetryDelay      time.Duration `json:"retry_delay"`
}

// DefaultProcessingOptions returns sensible defaults for processing operations
func DefaultProcessingOptions() *ProcessingOptions {
	return &ProcessingOptions{
		DryRun:                false,
		Verbose:               false,
		Timeout:               5 * time.Minute,
		Concurrency:           4,
		EnableVerification:    true,
		VerificationAlgorithm: "sha256",
		ContinueOnError:       false,
		MaxRetries:            3,
		RetryDelay:            time.Second,
	}
}
