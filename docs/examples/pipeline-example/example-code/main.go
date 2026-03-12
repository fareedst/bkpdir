// [REQ-PERFORMANCE] Processing patterns must support high-performance concurrent execution
// [ARCH-PROCESSING_PATTERNS] Processing pipeline architecture with worker pools
// [IMPL-PROCESSING_PATTERNS] Pipeline implementation (stages, processors)
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"bkpdir/pkg/fileops"
	"bkpdir/pkg/processing"
	"bkpdir/pkg/resources"
)

func main() {
	// Simple pipeline example using public interfaces
	ctx := context.Background()

	// Create resource manager for temporary artifacts
	rm := resources.NewResourceManager()
	defer func() {
		if err := rm.CleanupWithPanicRecovery(); err != nil {
			log.Printf("cleanup error: %v", err)
		}
	}()

	// Create processing pipeline
	pipeline := processing.NewPipeline("example-pipeline")

	// Add a simple stage that verifies input exists
	pipeline.AddStage(&verifyStage{})

	// Prepare input
	tinput := &processing.ProcessingInput{
		Source:      ".", // current directory
		Destination: "./out",
		Timeout:     30 * time.Second,
	}

	// Ensure destination exists
	if err := os.MkdirAll(tinput.Destination, 0o755); err != nil {
		log.Fatalf("failed to create output dir: %v", err)
	}

	// Execute pipeline
	result, err := pipeline.Execute(ctx, tinput)
	if err != nil {
		log.Fatalf("pipeline failed: %v", err)
	}

	fmt.Printf("Processed %d items in %v\n", result.ItemsProcessed, result.Duration)
}

// verifyStage is a minimal pipeline stage used for documentation examples.
type verifyStage struct{}

func (v *verifyStage) GetName() string        { return "verify" }
func (v *verifyStage) GetDescription() string { return "verify input exists" }
func (v *verifyStage) Execute(ctx context.Context, input *processing.ProcessingInput, output *processing.ProcessingResult) error {
	// Example: list files using fileops and count them
	files, err := fileops.ListFiles(input.Source, false)
	if err != nil {
		return err
	}

	output.ItemsProcessed = len(files)
	output.Duration = 0
	return nil
}

func (v *verifyStage) CanSkip(input *processing.ProcessingInput) bool { return false }
func (v *verifyStage) GetEstimatedDuration() time.Duration            { return 0 }

// Ensure verifyStage implements expected interface
var _ processing.PipelineStage = (*verifyStage)(nil)
