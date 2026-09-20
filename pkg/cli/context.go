package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DefaultContextManager provides standard context management functionality
type DefaultContextManager struct{}

// NewContextManager creates a new context manager
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultCommandBuilder with a non-nil FlagManager (defaulting when nil).
func NewContextManager() ContextManager {
	return &DefaultContextManager{}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: return cancellable child context defaulting parent to background.
func (cm *DefaultContextManager) Create(parent context.Context) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithCancel(parent)
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: parse duration string; on parse failure return cancel-only context without timeout.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: listen for interrupt/terminate and invoke cancel when signal received.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wrap function as operation that returns canceled when Cancel was called.
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

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create cancelable context that cancels on INT/TERM or parent done.
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
