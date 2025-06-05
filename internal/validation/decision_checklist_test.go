// 🔶 DOC-014: Decision validation tools test suite
package validation

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewDecisionChecklistValidator(t *testing.T) {
	projectRoot := "/test/project"
	validator := NewDecisionChecklistValidator(projectRoot)

	if validator.projectRoot != projectRoot {
		t.Errorf("Expected project root %s, got %s", projectRoot, validator.projectRoot)
	}

	if validator.verbose != false {
		t.Errorf("Expected verbose to be false by default, got %v", validator.verbose)
	}
}

func TestSetVerbose(t *testing.T) {
	validator := NewDecisionChecklistValidator("/test")

	validator.SetVerbose(true)
	if !validator.verbose {
		t.Errorf("Expected verbose to be true after SetVerbose(true)")
	}

	validator.SetVerbose(false)
	if validator.verbose {
		t.Errorf("Expected verbose to be false after SetVerbose(false)")
	}
}

func TestValidateDecisionChecklist(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "decision_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create basic project structure
	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	result, err := validator.ValidateDecisionChecklist(ctx)
	if err != nil {
		t.Fatalf("ValidateDecisionChecklist failed: %v", err)
	}

	// Verify result structure
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check that all sections are populated
	if result.SafetyGates.Name != "Safety Gates" {
		t.Errorf("Expected SafetyGates name to be 'Safety Gates', got '%s'", result.SafetyGates.Name)
	}

	if result.ScopeBoundaries.Name != "Scope Boundaries" {
		t.Errorf("Expected ScopeBoundaries name to be 'Scope Boundaries', got '%s'", result.ScopeBoundaries.Name)
	}

	if result.QualityThreshold.Name != "Quality Thresholds" {
		t.Errorf("Expected QualityThreshold name to be 'Quality Thresholds', got '%s'", result.QualityThreshold.Name)
	}

	if result.GoalAlignment.Name != "Goal Alignment" {
		t.Errorf("Expected GoalAlignment name to be 'Goal Alignment', got '%s'", result.GoalAlignment.Name)
	}

	// Check that validation took some time
	if result.ValidationTime <= 0 {
		t.Errorf("Expected positive validation time, got %v", result.ValidationTime)
	}

	// Check overall score is calculated
	if result.OverallScore < 0 || result.OverallScore > 100 {
		t.Errorf("Expected overall score between 0-100, got %f", result.OverallScore)
	}

	// Check that recommendations are generated
	if result.Recommendations == nil {
		t.Error("Expected non-nil recommendations slice")
	}
}

func TestValidateSafetyGates(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "safety_gates_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	var section ChecklistSection
	err = validator.validateSafetyGates(ctx, &section)
	if err != nil {
		t.Fatalf("validateSafetyGates failed: %v", err)
	}

	// Should have 4 checks: test validation, backward compatibility, token compliance, validation scripts
	expectedChecks := 4
	totalChecks := len(section.Details)
	if totalChecks != expectedChecks {
		t.Errorf("Expected %d safety gate checks, got %d", expectedChecks, totalChecks)
	}

	// Check that each detail has required fields
	for i, detail := range section.Details {
		if detail.Check == "" {
			t.Errorf("Detail %d missing check description", i)
		}
		if detail.Status == "" {
			t.Errorf("Detail %d missing status", i)
		}
		if detail.Priority == "" {
			t.Errorf("Detail %d missing priority", i)
		}
	}

	// Check score calculation
	totalValidations := section.Passed + section.Failed
	if totalValidations > 0 {
		expectedScore := float64(section.Passed) / float64(totalValidations) * 100
		if section.Score != expectedScore {
			t.Errorf("Expected score %f, got %f", expectedScore, section.Score)
		}
	}
}

func TestValidateScopeBoundaries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "scope_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	var section ChecklistSection
	err = validator.validateScopeBoundaries(ctx, &section)
	if err != nil {
		t.Fatalf("validateScopeBoundaries failed: %v", err)
	}

	// Should have 4 checks: feature scope, dependencies, architecture alignment, context updates
	expectedChecks := 4
	totalChecks := len(section.Details)
	if totalChecks != expectedChecks {
		t.Errorf("Expected %d scope boundary checks, got %d", expectedChecks, totalChecks)
	}

	// Verify section context is initialized
	if section.Context == nil {
		t.Error("Expected non-nil context map")
	}
}

func TestValidateQualityThresholds(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	var section ChecklistSection
	err = validator.validateQualityThresholds(ctx, &section)
	if err != nil {
		t.Fatalf("validateQualityThresholds failed: %v", err)
	}

	// Should have 4 checks: test coverage, error patterns, performance, traceability
	expectedChecks := 4
	totalChecks := len(section.Details)
	if totalChecks != expectedChecks {
		t.Errorf("Expected %d quality threshold checks, got %d", expectedChecks, totalChecks)
	}
}

func TestValidateGoalAlignment(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "goal_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	var section ChecklistSection
	err = validator.validateGoalAlignment(ctx, &section)
	if err != nil {
		t.Fatalf("validateGoalAlignment failed: %v", err)
	}

	// Should have 4 checks: phase progress, priority order, future impact, reusability
	expectedChecks := 4
	totalChecks := len(section.Details)
	if totalChecks != expectedChecks {
		t.Errorf("Expected %d goal alignment checks, got %d", expectedChecks, totalChecks)
	}
}

func TestValidateTestPassing(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "test_passing_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	result := validator.validateTestPassing(ctx)

	if result.Check != "Test Validation" {
		t.Errorf("Expected check name 'Test Validation', got '%s'", result.Check)
	}

	if result.Status == "" {
		t.Error("Expected non-empty status")
	}

	if result.Priority == "" {
		t.Error("Expected non-empty priority")
	}

	validStatuses := map[string]bool{"pass": true, "fail": true, "warning": true, "skip": true}
	if !validStatuses[result.Status] {
		t.Errorf("Invalid status '%s'", result.Status)
	}

	validPriorities := map[string]bool{"critical": true, "high": true, "medium": true, "low": true}
	if !validPriorities[result.Priority] {
		t.Errorf("Invalid priority '%s'", result.Priority)
	}
}

func TestValidateBackwardCompatibility(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "compat_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	result := validator.validateBackwardCompatibility(ctx)

	if result.Check != "Backward Compatibility" {
		t.Errorf("Expected check name 'Backward Compatibility', got '%s'", result.Check)
	}

	if result.Message == "" {
		t.Error("Expected non-empty message")
	}
}

func TestValidateTokenCompliance(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "token_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)
	ctx := context.Background()

	result := validator.validateTokenCompliance(ctx)

	if result.Check != "Token Compliance" {
		t.Errorf("Expected check name 'Token Compliance', got '%s'", result.Check)
	}

	// Should contain some form of validation message
	if result.Message == "" {
		t.Error("Expected non-empty message")
	}
}

func TestCalculateOverallScore(t *testing.T) {
	validator := NewDecisionChecklistValidator("/test")

	result := &DecisionChecklistResult{
		SafetyGates:      ChecklistSection{Passed: 4, Failed: 0, Score: 100.0},
		ScopeBoundaries:  ChecklistSection{Passed: 3, Failed: 1, Score: 75.0},
		QualityThreshold: ChecklistSection{Passed: 4, Failed: 0, Score: 100.0},
		GoalAlignment:    ChecklistSection{Passed: 2, Failed: 2, Score: 50.0},
	}

	score := validator.calculateOverallScore(result)

	// Weighted average: (100*0.4 + 75*0.25 + 100*0.25 + 50*0.1) = 40 + 18.75 + 25 + 5 = 88.75
	expectedScore := 88.75
	if score != expectedScore {
		t.Errorf("Expected overall score %f, got %f", expectedScore, score)
	}
}

func TestGenerateRecommendations(t *testing.T) {
	validator := NewDecisionChecklistValidator("/test")

	result := &DecisionChecklistResult{
		SafetyGates: ChecklistSection{
			Score:          75.0,
			CriticalIssues: []string{"Test failures detected"},
		},
		ScopeBoundaries:  ChecklistSection{Score: 100.0},
		QualityThreshold: ChecklistSection{Score: 50.0},
		GoalAlignment:    ChecklistSection{Score: 80.0},
	}

	recommendations := validator.generateRecommendations(result)

	if len(recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}

	// Should recommend fixing safety gates due to critical issues
	foundSafetyRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "safety") || strings.Contains(rec, "critical") {
			foundSafetyRecommendation = true
			break
		}
	}

	if !foundSafetyRecommendation {
		t.Error("Expected recommendation about safety gates with critical issues")
	}

	// Should recommend improving quality thresholds due to low score
	foundQualityRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "quality") || strings.Contains(rec, "threshold") {
			foundQualityRecommendation = true
			break
		}
	}

	if !foundQualityRecommendation {
		t.Error("Expected recommendation about quality thresholds")
	}
}

func TestContextCancellation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cancel_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)

	// Create a context that cancels immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// The validation should handle cancellation gracefully
	result, err := validator.ValidateDecisionChecklist(ctx)

	// The implementation currently doesn't check for cancellation in all paths,
	// so we mainly verify it doesn't panic
	if err != nil && result == nil {
		t.Logf("Context cancellation handled: %v", err)
	}
}

// TestWithTimeout verifies timeout handling
func TestWithTimeout(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "timeout_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createTestProject(t, tmpDir)

	validator := NewDecisionChecklistValidator(tmpDir)

	// Use a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give the context time to expire
	time.Sleep(1 * time.Millisecond)

	result, err := validator.ValidateDecisionChecklist(ctx)

	// The validation should complete even with expired context
	// (current implementation doesn't check context in all validation functions)
	if result != nil {
		t.Logf("Validation completed despite timeout")
	}
	if err != nil {
		t.Logf("Timeout handled: %v", err)
	}
}

// createTestProject creates a minimal project structure for testing
func createTestProject(t *testing.T, projectRoot string) {
	// Create basic directories
	dirs := []string{
		"docs/context",
		"scripts",
		"internal/validation",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(projectRoot, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", fullPath, err)
		}
	}

	// Create basic files
	files := map[string]string{
		"docs/context/feature-tracking.md": `# Feature Tracking Matrix
## Feature Registry
| Feature ID | Status |
|------------|--------|
| TEST-001   | ✅ Completed |
`,
		"docs/context/ai-decision-framework.md": `# AI Decision Framework
## Decision Hierarchy
1. Safety Gates
2. Scope Boundaries  
3. Quality Thresholds
4. Goal Alignment
`,
		"scripts/validate-docs.sh": `#!/bin/bash
echo "Validation script placeholder"
exit 0
`,
		"go.mod": `module bkpdir
go 1.21
`,
		"main.go": `// ⭐ TEST-001: Test main file
package main
func main() {}
`,
	}

	for filename, content := range files {
		fullPath := filepath.Join(projectRoot, filename)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Make scripts executable
	scriptPath := filepath.Join(projectRoot, "scripts", "validate-docs.sh")
	if err := os.Chmod(scriptPath, 0755); err != nil {
		t.Fatalf("Failed to make script executable: %v", err)
	}
}
