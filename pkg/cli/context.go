package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ⭐ EXTRACT-005: Context-aware command execution with cancellation support - 🔧

// DefaultContextManager provides standard context management functionality
type DefaultContextManager struct{}

// NewContextManager creates a new context manager
func NewContextManager() ContextManager {
	return &DefaultContextManager{}
}

// Create creates a new context with cancellation
func (cm *DefaultContextManager) Create(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithCancel(parent)
}

// WithTimeout creates a context with timeout
func (cm *DefaultContextManager) WithTimeout(parent context.Context, timeout string) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		// If parsing fails, return a context without timeout
		return context.WithCancel(parent)
	}

	return context.WithTimeout(parent, duration)
}

// HandleSignals sets up signal handling for graceful shutdown
func (cm *DefaultContextManager) HandleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		if cancel != nil {
			cancel()
		}
	}()
}

// SimpleCancellableOperation provides a basic implementation of CancellableOperation
type SimpleCancellableOperation struct {
	operation func(ctx context.Context) error
	cancelled bool
}

// NewCancellableOperation creates a simple cancellable operation
func NewCancellableOperation(op func(ctx context.Context) error) CancellableOperation {
	return &SimpleCancellableOperation{
		operation: op,
		cancelled: false,
	}
}

// Execute performs the operation with cancellation support
func (op *SimpleCancellableOperation) Execute(ctx context.Context) error {
	if op.cancelled {
		return context.Canceled
	}

	if op.operation == nil {
		return nil
	}

	return op.operation(ctx)
}

// Cancel requests cancellation of the operation
func (op *SimpleCancellableOperation) Cancel() error {
	op.cancelled = true
	return nil
}

// WithSignalHandling creates a context with automatic signal handling
func WithSignalHandling(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	return ctx, cancel
}
