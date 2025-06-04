package cli

import (
	"fmt"
)

// ⭐ EXTRACT-005: Generic dry-run operation support extracted from main.go - 🔧

// DefaultDryRunManager provides standard dry-run functionality
type DefaultDryRunManager struct{}

// NewDryRunManager creates a new dry-run manager
func NewDryRunManager() DryRunManager {
	return &DefaultDryRunManager{}
}

// Execute runs the operation or simulates it based on dry-run flag
func (drm *DefaultDryRunManager) Execute(ctx CommandContext, op DryRunOperation) error {
	if ctx.DryRun {
		// In dry-run mode, log what would be done and skip execution
		drm.Log(ctx, op.Describe())
		return nil
	}

	// Execute the actual operation
	return op.Execute(ctx)
}

// Log outputs dry-run information to the appropriate writer
func (drm *DefaultDryRunManager) Log(ctx CommandContext, message string) {
	fmt.Fprintf(ctx.Output, "[DRY-RUN] %s\n", message)
}

// SimpleDryRunOperation provides a basic implementation of DryRunOperation
type SimpleDryRunOperation struct {
	description string
	operation   func(ctx CommandContext) error
}

// NewSimpleDryRunOperation creates a simple dry-run operation
func NewSimpleDryRunOperation(description string, op func(ctx CommandContext) error) DryRunOperation {
	return &SimpleDryRunOperation{
		description: description,
		operation:   op,
	}
}

// Execute performs the actual operation
func (op *SimpleDryRunOperation) Execute(ctx CommandContext) error {
	if op.operation == nil {
		return nil
	}
	return op.operation(ctx)
}

// Describe returns a description of what the operation would do
func (op *SimpleDryRunOperation) Describe() string {
	return op.description
}

// DryRunWrapper wraps a function to make it compatible with dry-run execution
func DryRunWrapper(description string, fn func(CommandContext) error) DryRunOperation {
	return NewSimpleDryRunOperation(description, fn)
}
