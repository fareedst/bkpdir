// [REQ-CODE_QUALITY]
// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
package validation

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestRealTimeValidation tests the main real-time validation functionality for DOC-012
func TestRealTimeValidation(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()
	if validator == nil {
		t.Fatal("Failed to create real-time validator")
	}

	// Test the cache and subscription components instead of full validation
	// since full validation requires external make command setup

	// Test 1: Cache functionality
	testFile := "test.go"
	testContent := "package main"

	// Create a mock response
	mockResponse := &ValidationResponse{
		Status:         "pass",
		Errors:         []AIOptimizedError{},
		Warnings:       []AIOptimizedWarning{},
		ProcessingTime: 50 * time.Millisecond,
	}

	// Test cache operations
	validator.cache.Set(testFile, testContent, mockResponse)
	cached := validator.cache.Get(testFile, testContent)

	if cached == nil {
		t.Error("Cache should return stored response")
	}

	if cached.Status != "pass" {
		t.Errorf("Expected cached status 'pass', got %s", cached.Status)
	}

	// Test 2: Subscription management
	subscriberID := "test-subscriber"
	files := []string{testFile}

	updates := validator.SubscribeToValidation(subscriberID, files)
	if updates == nil {
		t.Fatal("Expected update channel, got nil")
	}

	// Verify subscriber was added
	validator.mu.RLock()
	subscriber, exists := validator.subscribers[subscriberID]
	validator.mu.RUnlock()

	if !exists {
		t.Error("Subscriber not found after subscription")
	}

	if subscriber.ID != subscriberID {
		t.Errorf("Expected subscriber ID %s, got %s", subscriberID, subscriber.ID)
	}

	// Test 3: Status indicator creation
	indicator := validator.GetValidationStatusIndicator(files)
	if indicator == nil {
		t.Fatal("Expected status indicator, got nil")
	}

	if indicator.OverallStatus == "" {
		t.Error("Expected overall status to be set")
	}

	// Cleanup
	validator.UnsubscribeFromValidation(subscriberID)

	// Verify performance metrics are being tracked
	if validator.metricsTracker == nil {
		t.Error("Expected metrics tracker to be initialized")
	}
}

// TestValidationCache tests the caching mechanism for performance optimization
func TestValidationCache(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	cache := NewValidationCache(time.Minute)
	if cache == nil {
		t.Fatal("Failed to create validation cache")
	}

	testFile := "cache_test.go"
	testContent := "package main"

	// Create mock response
	response := &ValidationResponse{
		Status:         "pass",
		Errors:         []AIOptimizedError{},
		Warnings:       []AIOptimizedWarning{},
		ProcessingTime: 50 * time.Millisecond,
	}

	// Test cache set and get
	cache.Set(testFile, testContent, response)

	cached := cache.Get(testFile, testContent)
	if cached == nil {
		t.Error("Expected cached response, got nil")
	}

	if cached.Status != response.Status {
		t.Errorf("Expected status %s, got %s", response.Status, cached.Status)
	}
}

// TestSubscriptionManagement tests the real-time subscription system
func TestSubscriptionManagement(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	subscriberID := "test-subscriber"
	files := []string{"file1.go", "file2.go"}

	// Test subscription
	updates := validator.SubscribeToValidation(subscriberID, files)
	if updates == nil {
		t.Fatal("Expected update channel, got nil")
	}

	// Verify subscriber was added
	validator.mu.RLock()
	subscriber, exists := validator.subscribers[subscriberID]
	validator.mu.RUnlock()

	if !exists {
		t.Error("Subscriber not found after subscription")
	}

	if subscriber.ID != subscriberID {
		t.Errorf("Expected subscriber ID %s, got %s", subscriberID, subscriber.ID)
	}

	// Test unsubscription
	validator.UnsubscribeFromValidation(subscriberID)

	validator.mu.RLock()
	_, exists = validator.subscribers[subscriberID]
	validator.mu.RUnlock()

	if exists {
		t.Error("Subscriber still exists after unsubscription")
	}
}

// TestValidationStatusIndicator tests the visual status indicator system
func TestValidationStatusIndicator(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	files := []string{"test1.go", "test2.go"}
	indicator := validator.GetValidationStatusIndicator(files)

	if indicator == nil {
		t.Fatal("Expected status indicator, got nil")
	}

	if indicator.OverallStatus == "" {
		t.Error("Expected overall status to be set")
	}

	if indicator.VisualElements == nil {
		t.Error("Expected visual elements to be present")
	}

	if indicator.ComplianceLevel == "" {
		t.Error("Expected compliance level to be set")
	}
}

// TestIntelligentSuggestions tests the intelligent suggestion generation
func TestIntelligentSuggestions(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	// Create mock response with errors
	response := &ValidationResponse{
		Status: "fail",
		Errors: []AIOptimizedError{
			{
				ErrorID:  "test-error-1",
				Category: "icon_format",
				Message:  "Missing priority icon",
				Context:  map[string]string{"type": "token"},
			},
		},
		Warnings: []AIOptimizedWarning{
			{
				WarningID: "test-warning-1",
				Category:  "token_consistency",
				Message:   "Inconsistent token format",
			},
		},
	}

	suggestions := validator.generateIntelligentSuggestions("test.go", "test content", response)

	if len(suggestions) == 0 {
		t.Error("Expected suggestions to be generated")
	}

	// Verify error-based suggestions
	foundErrorSuggestion := false
	for _, suggestion := range suggestions {
		if suggestion.Type == "icon_fix" {
			foundErrorSuggestion = true
			break
		}
	}

	if !foundErrorSuggestion {
		t.Error("Expected icon fix suggestion based on error")
	}
}

// TestPerformanceMetrics tests the metrics tracking system
func TestPerformanceMetrics(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	metrics := &PerformanceMetrics{
		lastMetricsUpdate: time.Now(),
	}

	// Test cache hit recording
	metrics.RecordCacheHit()
	if metrics.cacheHitRate <= 0 {
		t.Error("Expected cache hit rate to be updated")
	}

	// Test validation recording
	testDuration := 100 * time.Millisecond
	metrics.RecordValidation(testDuration)

	if metrics.totalValidations != 1 {
		t.Errorf("Expected 1 validation, got %d", metrics.totalValidations)
	}

	if metrics.averageResponseTime != testDuration {
		t.Errorf("Expected response time %v, got %v", testDuration, metrics.averageResponseTime)
	}

	// Test subscriber tracking
	metrics.IncrementActiveSubscribers()
	if metrics.activeSubscribers != 1 {
		t.Errorf("Expected 1 active subscriber, got %d", metrics.activeSubscribers)
	}

	metrics.DecrementActiveSubscribers()
	if metrics.activeSubscribers != 0 {
		t.Errorf("Expected 0 active subscribers, got %d", metrics.activeSubscribers)
	}
}

// TestVisualElements tests the visual feedback element creation
func TestVisualElements(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	testCases := []struct {
		status        string
		compliance    string
		expectedColor string
		expectedIcon  string
		expectedBadge string
	}{
		{"pass", "excellent", "#28a745", "[SUCCESS]", "Perfect"},
		{"warning", "good", "#ffc107", "⚠️", "Good"},
		{"error", "poor", "#dc3545", "❌", "Poor"},
		{"unknown", "needs_work", "#6c757d", "❓", "Issues"},
	}

	for _, tc := range testCases {
		elements := validator.createVisualElements(tc.status, tc.compliance)

		if elements.StatusColor != tc.expectedColor {
			t.Errorf("Status %s: expected color %s, got %s", tc.status, tc.expectedColor, elements.StatusColor)
		}

		if elements.StatusIcon != tc.expectedIcon {
			t.Errorf("Status %s: expected icon %s, got %s", tc.status, tc.expectedIcon, elements.StatusIcon)
		}

		if elements.BadgeText != tc.expectedBadge {
			t.Errorf("Compliance %s: expected badge %s, got %s", tc.compliance, tc.expectedBadge, elements.BadgeText)
		}
	}
}

// TestComplianceCalculation tests the compliance level calculation logic
func TestComplianceCalculation(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	testCases := []struct {
		errors   int
		warnings int
		files    int
		expected string
	}{
		{0, 0, 1, "excellent"},
		{0, 1, 1, "good"},
		{1, 0, 1, "needs_work"},
		{2, 3, 1, "poor"},
	}

	for _, tc := range testCases {
		result := validator.calculateComplianceLevel(tc.errors, tc.warnings, tc.files)
		if result != tc.expected {
			t.Errorf("Errors: %d, Warnings: %d, Files: %d - expected %s, got %s",
				tc.errors, tc.warnings, tc.files, tc.expected, result)
		}
	}
}

// TestProactiveSuggestions tests the proactive suggestion generation
func TestProactiveSuggestions(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	// Test content with tokens that could be enhanced
	testContent := `
// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
// FEATURE-001: See specification.md - Feature Implementation [DECISION:maintenance]
package main
`

	suggestions := validator.generateProactiveSuggestions("test.go", testContent)

	if len(suggestions) == 0 {
		t.Error("Expected proactive suggestions for unenhanced tokens")
	}

	// Verify suggestion types
	for _, suggestion := range suggestions {
		if suggestion.Type != "icon_enhancement" {
			t.Errorf("Expected icon_enhancement suggestion, got %s", suggestion.Type)
		}

		if suggestion.Confidence <= 0 || suggestion.Confidence > 1 {
			t.Errorf("Expected confidence between 0 and 1, got %f", suggestion.Confidence)
		}
	}
}

// TestCacheExpiration tests the cache TTL functionality
func TestCacheExpiration(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	cache := NewValidationCache(100 * time.Millisecond) // Short TTL for testing

	testFile := "expire_test.go"
	testContent := "package main"
	response := &ValidationResponse{Status: "pass"}

	// Set cache entry
	cache.Set(testFile, testContent, response)

	// Verify it exists
	cached := cache.Get(testFile, testContent)
	if cached == nil {
		t.Error("Expected cached response immediately after set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Verify it's expired
	expired := cache.Get(testFile, testContent)
	if expired != nil {
		t.Error("Expected cached response to be expired")
	}
}

// TestSuggestionConfidence tests the confidence scoring for suggestions
func TestSuggestionConfidence(t *testing.T) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()

	// Test different error types and their confidence levels
	testErrors := []AIOptimizedError{
		{Category: "icon_format", Message: "Missing priority icon"},
		{Category: "token_format", Message: "Invalid token format"},
		{Category: "priority_mismatch", Message: "Priority mismatch"},
		{Category: "unknown", Message: "Unknown error"},
	}

	expectedConfidences := map[string]float64{
		"icon_format":       0.9,
		"token_format":      0.85,
		"priority_mismatch": 0.8,
		"unknown":           0.5,
	}

	for _, error := range testErrors {
		suggestion := validator.createSuggestionFromError("test.go", "content", error)

		if suggestion == nil {
			t.Errorf("Expected suggestion for error category %s", error.Category)
			continue
		}

		expected := expectedConfidences[error.Category]
		if suggestion.Confidence != expected {
			t.Errorf("Category %s: expected confidence %f, got %f",
				error.Category, expected, suggestion.Confidence)
		}
	}
}

// BenchmarkRealTimeValidation benchmarks the real-time validation performance
func BenchmarkRealTimeValidation(b *testing.B) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	validator := NewRealTimeValidator()
	ctx := context.Background()
	testFile := "benchmark.go"
	testContent := `
// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
package main

func main() {
	println("Benchmark test")
}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := validator.ValidateRealtimeFile(ctx, testFile, testContent)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
	}
}

// BenchmarkCachePerformance benchmarks the cache performance
func BenchmarkCachePerformance(b *testing.B) {
	// DOC-012: See ai-decision-framework.md - Documentation Standards [DECISION:validation]
	cache := NewValidationCache(time.Hour)
	response := &ValidationResponse{Status: "pass"}

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("file%d.go", i), "content", response)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("file%d.go", i%1000), "content")
	}
}
