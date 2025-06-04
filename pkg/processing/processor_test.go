// ⭐ EXTRACT-007: Processing package structure - Comprehensive test suite for all components - 🧪
package processing

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Test naming provider functionality
func TestNamingProvider(t *testing.T) {
	np := NewNamingProvider()

	// Test archive name generation
	template := &NamingTemplate{
		Prefix:             "test",
		Timestamp:          time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		GitBranch:          "main",
		GitHash:            "abc123",
		Note:               "backup",
		GitIsClean:         true,
		ShowGitDirtyStatus: false,
	}

	name, err := np.GenerateName(template)
	if err != nil {
		t.Fatalf("Failed to generate name: %v", err)
	}

	expected := "test-2024-01-01T120000-main-abc123-backup"
	if name != expected {
		t.Errorf("Expected %s, got %s", expected, name)
	}

	// Test name parsing
	fullName := name + ".zip"
	t.Logf("Generated full name: %s", fullName) // Debug output
	components, err := np.ParseName(fullName, "archive")
	if err != nil {
		t.Fatalf("Failed to parse name: %v", err)
	}

	if components.Prefix != "test" {
		t.Errorf("Expected prefix 'test', got '%s'", components.Prefix)
	}

	if components.GitBranch != "main" {
		t.Errorf("Expected git branch 'main', got '%s'", components.GitBranch)
	}

	// Test supported formats
	formats := np.GetSupportedFormats()
	if len(formats) == 0 {
		t.Error("Expected supported formats, got none")
	}
}

// Test verification provider functionality
func TestVerificationProvider(t *testing.T) {
	verifier := NewSHA256Verifier()

	// Test checksum calculation
	data := strings.NewReader("test data")
	checksum, err := verifier.Calculate(data)
	if err != nil {
		t.Fatalf("Failed to calculate checksum: %v", err)
	}

	if checksum == "" {
		t.Error("Expected checksum, got empty string")
	}

	// Test verification
	data2 := strings.NewReader("test data")
	isValid, err := verifier.Verify(data2, checksum)
	if err != nil {
		t.Fatalf("Failed to verify checksum: %v", err)
	}

	if !isValid {
		t.Error("Expected valid checksum verification")
	}

	// Test algorithm info
	if verifier.GetAlgorithm() != "sha256" {
		t.Errorf("Expected algorithm 'sha256', got '%s'", verifier.GetAlgorithm())
	}
}

// Test verification manager functionality
func TestVerificationManager(t *testing.T) {
	vm := NewVerificationManager()

	// Test supported algorithms
	algorithms := vm.GetSupportedAlgorithms()
	if len(algorithms) == 0 {
		t.Error("Expected supported algorithms, got none")
	}

	// Test provider retrieval
	provider, err := vm.GetProvider("sha256")
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}

	if provider.GetAlgorithm() != "sha256" {
		t.Errorf("Expected algorithm 'sha256', got '%s'", provider.GetAlgorithm())
	}

	// Test default provider
	defaultProvider, err := vm.GetProvider("")
	if err != nil {
		t.Fatalf("Failed to get default provider: %v", err)
	}

	if defaultProvider.GetAlgorithm() != "sha256" {
		t.Errorf("Expected default algorithm 'sha256', got '%s'", defaultProvider.GetAlgorithm())
	}
}

// Test pipeline functionality
func TestPipeline(t *testing.T) {
	pipeline := NewPipeline("test-pipeline")

	// Create test stages
	stage1 := NewTestStage("stage1", "Test stage 1", 100*time.Millisecond)
	stage2 := NewTestStage("stage2", "Test stage 2", 50*time.Millisecond)

	pipeline.AddStage(stage1)
	pipeline.AddStage(stage2)

	// Test stage count
	stages := pipeline.GetStages()
	if len(stages) != 2 {
		t.Errorf("Expected 2 stages, got %d", len(stages))
	}

	// Test execution
	ctx := context.Background()
	input := &ProcessingInput{
		Source: "test",
		Options: map[string]interface{}{
			"test": true,
		},
	}

	result, err := pipeline.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Pipeline had errors: %v", result.Errors)
	}

	// Test progress
	progress := pipeline.GetProgress()
	if progress.OverallProgress != 1.0 {
		t.Errorf("Expected progress 1.0, got %f", progress.OverallProgress)
	}
}

// Test concurrent processing functionality
func TestConcurrentProcessor(t *testing.T) {
	// Create test processor function
	processFunc := func(ctx context.Context, item *ProcessingItem) (interface{}, error) {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return map[string]interface{}{
			"processed": item.ID,
			"data":      item.Data,
		}, nil
	}

	processor := NewConcurrentProcessor(processFunc)
	processor.SetWorkerCount(2)

	// Create test items
	items := make([]ProcessingItem, 10)
	for i := 0; i < 10; i++ {
		items[i] = ProcessingItem{
			ID:   fmt.Sprintf("item-%d", i),
			Data: fmt.Sprintf("data-%d", i),
		}
	}

	// Process items
	ctx := context.Background()
	result, err := processor.Process(ctx, items)
	if err != nil {
		t.Fatalf("Concurrent processing failed: %v", err)
	}

	// Validate results
	if result.TotalItems != 10 {
		t.Errorf("Expected 10 total items, got %d", result.TotalItems)
	}

	if result.SuccessfulItems != 10 {
		t.Errorf("Expected 10 successful items, got %d", result.SuccessfulItems)
	}

	if result.FailedItems != 0 {
		t.Errorf("Expected 0 failed items, got %d", result.FailedItems)
	}

	if len(result.Results) != 10 {
		t.Errorf("Expected 10 results, got %d", len(result.Results))
	}
}

// Test context cancellation
func TestContextCancellation(t *testing.T) {
	// Create a slow processor function that always checks for cancellation
	processFunc := func(ctx context.Context, item *ProcessingItem) (interface{}, error) {
		// Simulate some initial work
		time.Sleep(10 * time.Millisecond)

		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			// Continue with more work
			time.Sleep(100 * time.Millisecond) // This should be interrupted
			return "processed", nil
		}
	}

	processor := NewConcurrentProcessor(processFunc)
	processor.SetWorkerCount(2) // Use 2 workers to increase chance of concurrent processing

	// Create more test items to ensure some are being processed when timeout occurs
	items := make([]ProcessingItem, 10)
	for i := 0; i < 10; i++ {
		items[i] = ProcessingItem{
			ID:   fmt.Sprintf("item-%d", i),
			Data: "test",
		}
	}

	// Create context with very short timeout to ensure cancellation during processing
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	// Process items (should be cancelled)
	result, err := processor.Process(ctx, items)

	// Context cancellation should cause one of these conditions:
	// 1. The Process method returns an error, OR
	// 2. Some tasks fail due to context cancellation, OR
	// 3. Fewer items are processed than submitted (due to cancellation during submission)

	if err != nil {
		// Case 1: Process method returned an error - this is valid cancellation behavior
		t.Logf("Process method returned error due to cancellation: %v", err)
		return
	}

	if result == nil {
		t.Fatal("Result is nil but no error was returned")
	}

	if result.FailedItems > 0 {
		// Case 2: Some tasks failed due to cancellation - this is valid
		t.Logf("Context cancellation caused %d failed items", result.FailedItems)
		return
	}

	if result.TotalItems < len(items) {
		// Case 3: Fewer items processed than submitted - this means cancellation during submission
		t.Logf("Context cancellation stopped submission early. Processed %d out of %d items",
			result.TotalItems, len(items))
		return
	}

	// If we reach here, cancellation didn't have the expected effect
	t.Errorf("Expected cancellation to have an effect. Got err=%v, failed=%d, successful=%d, total=%d, submitted=%d",
		err, result.FailedItems, result.SuccessfulItems, result.TotalItems, len(items))
}

// Test utility functions
func TestUtilityFunctions(t *testing.T) {
	// Test ProcessItems convenience function
	processFunc := func(ctx context.Context, item *ProcessingItem) (interface{}, error) {
		return "processed", nil
	}

	items := []ProcessingItem{
		{ID: "1", Data: "test1"},
		{ID: "2", Data: "test2"},
	}

	ctx := context.Background()
	result, err := ProcessItems(ctx, items, processFunc)
	if err != nil {
		t.Fatalf("ProcessItems failed: %v", err)
	}

	if result.TotalItems != 2 {
		t.Errorf("Expected 2 items, got %d", result.TotalItems)
	}

	// Test ProcessItemsWithOptions
	result2, err := ProcessItemsWithOptions(ctx, items, processFunc, 1, 50)
	if err != nil {
		t.Fatalf("ProcessItemsWithOptions failed: %v", err)
	}

	if result2.WorkerCount != 1 {
		t.Errorf("Expected 1 worker, got %d", result2.WorkerCount)
	}
}

// Test helper types and functions

// TestStage implements PipelineStage for testing
type TestStage struct {
	*BaseStage
	shouldFail bool
}

func NewTestStage(name, description string, duration time.Duration) *TestStage {
	return &TestStage{
		BaseStage: NewBaseStage(name, description, duration),
	}
}

func (ts *TestStage) Execute(ctx context.Context, input *ProcessingInput, output *ProcessingResult) error {
	// Simulate work
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(ts.GetEstimatedDuration()):
		// Work completed
	}

	if ts.shouldFail {
		return NewProcessingError("TEST_FAILURE", "Execute", "test stage failure")
	}

	// Add some test output
	output.ItemsProcessed++
	return nil
}

func (ts *TestStage) SetShouldFail(fail bool) {
	ts.shouldFail = fail
}

// Benchmark tests
func BenchmarkNamingProvider(b *testing.B) {
	np := NewNamingProvider()
	template := &NamingTemplate{
		Prefix:    "bench",
		Timestamp: time.Now(),
		Note:      "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = np.GenerateName(template)
	}
}

func BenchmarkVerification(b *testing.B) {
	verifier := NewSHA256Verifier()
	data := strings.NewReader("benchmark data for testing")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data.Seek(0, 0)
		_, _ = verifier.Calculate(data)
	}
}

func BenchmarkConcurrentProcessing(b *testing.B) {
	processFunc := func(ctx context.Context, item *ProcessingItem) (interface{}, error) {
		return "processed", nil
	}

	items := make([]ProcessingItem, 100)
	for i := 0; i < 100; i++ {
		items[i] = ProcessingItem{
			ID:   fmt.Sprintf("item-%d", i),
			Data: "benchmark",
		}
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ProcessItems(ctx, items, processFunc)
	}
}
