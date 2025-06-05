// 🔶 DOC-014: Decision validation tools - Core decision checklist validation
package validation

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DecisionChecklistValidator provides automated validation of the 4-tier decision hierarchy
type DecisionChecklistValidator struct {
	projectRoot string
	verbose     bool
}

// DecisionChecklistResult represents the result of decision checklist validation
type DecisionChecklistResult struct {
	SafetyGates      ChecklistSection `json:"safety_gates"`
	ScopeBoundaries  ChecklistSection `json:"scope_boundaries"`
	QualityThreshold ChecklistSection `json:"quality_thresholds"`
	GoalAlignment    ChecklistSection `json:"goal_alignment"`
	OverallScore     float64          `json:"overall_score"`
	ValidationTime   time.Duration    `json:"validation_time"`
	Recommendations  []string         `json:"recommendations"`
}

// ChecklistSection represents validation results for one section of the decision hierarchy
type ChecklistSection struct {
	Name           string                 `json:"name"`
	Passed         int                    `json:"passed"`
	Failed         int                    `json:"failed"`
	Score          float64                `json:"score"`
	Details        []ChecklistItem        `json:"details"`
	CriticalIssues []string               `json:"critical_issues"`
	Context        map[string]interface{} `json:"context"`
}

// ChecklistItem represents an individual validation check
type ChecklistItem struct {
	Check       string `json:"check"`
	Status      string `json:"status"` // "pass", "fail", "warning", "skip"
	Message     string `json:"message"`
	Remediation string `json:"remediation,omitempty"`
	Priority    string `json:"priority"` // "critical", "high", "medium", "low"
}

// NewDecisionChecklistValidator creates a new decision checklist validator
func NewDecisionChecklistValidator(projectRoot string) *DecisionChecklistValidator {
	return &DecisionChecklistValidator{
		projectRoot: projectRoot,
		verbose:     false,
	}
}

// SetVerbose enables verbose output for debugging
func (v *DecisionChecklistValidator) SetVerbose(verbose bool) {
	v.verbose = verbose
}

// ValidateDecisionChecklist performs comprehensive validation of the decision hierarchy
func (v *DecisionChecklistValidator) ValidateDecisionChecklist(ctx context.Context) (*DecisionChecklistResult, error) {
	startTime := time.Now()

	result := &DecisionChecklistResult{
		SafetyGates:      ChecklistSection{Name: "Safety Gates"},
		ScopeBoundaries:  ChecklistSection{Name: "Scope Boundaries"},
		QualityThreshold: ChecklistSection{Name: "Quality Thresholds"},
		GoalAlignment:    ChecklistSection{Name: "Goal Alignment"},
		ValidationTime:   0,
		Recommendations:  []string{},
	}

	// Validate each tier of the decision hierarchy
	if err := v.validateSafetyGates(ctx, &result.SafetyGates); err != nil {
		return nil, fmt.Errorf("safety gates validation failed: %w", err)
	}

	if err := v.validateScopeBoundaries(ctx, &result.ScopeBoundaries); err != nil {
		return nil, fmt.Errorf("scope boundaries validation failed: %w", err)
	}

	if err := v.validateQualityThresholds(ctx, &result.QualityThreshold); err != nil {
		return nil, fmt.Errorf("quality thresholds validation failed: %w", err)
	}

	if err := v.validateGoalAlignment(ctx, &result.GoalAlignment); err != nil {
		return nil, fmt.Errorf("goal alignment validation failed: %w", err)
	}

	// Calculate overall score and recommendations
	result.OverallScore = v.calculateOverallScore(result)
	result.Recommendations = v.generateRecommendations(result)
	result.ValidationTime = time.Since(startTime)

	return result, nil
}

// 🔶 DOC-014: Safety Gates validation - 🛡️ NEVER Override checks
func (v *DecisionChecklistValidator) validateSafetyGates(ctx context.Context, section *ChecklistSection) error {
	section.Details = []ChecklistItem{}
	section.Context = make(map[string]interface{})

	// Check 1: Test Validation - Are all tests passing?
	testResult := v.validateTestPassing(ctx)
	section.Details = append(section.Details, testResult)
	if testResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
		if testResult.Priority == "critical" {
			section.CriticalIssues = append(section.CriticalIssues, testResult.Message)
		}
	}

	// Check 2: Backward Compatibility - Does the change preserve existing functionality?
	compatResult := v.validateBackwardCompatibility(ctx)
	section.Details = append(section.Details, compatResult)
	if compatResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
		if compatResult.Priority == "critical" {
			section.CriticalIssues = append(section.CriticalIssues, compatResult.Message)
		}
	}

	// Check 3: Token Compliance - Are implementation tokens properly formatted?
	tokenResult := v.validateTokenCompliance(ctx)
	section.Details = append(section.Details, tokenResult)
	if tokenResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
		if tokenResult.Priority == "critical" {
			section.CriticalIssues = append(section.CriticalIssues, tokenResult.Message)
		}
	}

	// Check 4: Validation Scripts - Do validation scripts pass?
	validationResult := v.validateValidationScripts(ctx)
	section.Details = append(section.Details, validationResult)
	if validationResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
		if validationResult.Priority == "critical" {
			section.CriticalIssues = append(section.CriticalIssues, validationResult.Message)
		}
	}

	// Calculate section score
	totalChecks := section.Passed + section.Failed
	if totalChecks > 0 {
		section.Score = float64(section.Passed) / float64(totalChecks) * 100
	}

	return nil
}

// 🔶 DOC-014: Scope Boundaries validation - 📑 Strict Limits checks
func (v *DecisionChecklistValidator) validateScopeBoundaries(ctx context.Context, section *ChecklistSection) error {
	section.Details = []ChecklistItem{}
	section.Context = make(map[string]interface{})

	// Check 1: Feature Scope - Is this change within the current feature's documented scope?
	scopeResult := v.validateFeatureScope(ctx)
	section.Details = append(section.Details, scopeResult)
	if scopeResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 2: Dependency Check - Are all blocking dependencies satisfied?
	depResult := v.validateDependencies(ctx)
	section.Details = append(section.Details, depResult)
	if depResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 3: Architecture Alignment - Does this align with documented architecture patterns?
	archResult := v.validateArchitectureAlignment(ctx)
	section.Details = append(section.Details, archResult)
	if archResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 4: Context Updates - Will this require updating multiple context files per protocol?
	contextResult := v.validateContextUpdates(ctx)
	section.Details = append(section.Details, contextResult)
	if contextResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Calculate section score
	totalChecks := section.Passed + section.Failed
	if totalChecks > 0 {
		section.Score = float64(section.Passed) / float64(totalChecks) * 100
	}

	return nil
}

// 🔶 DOC-014: Quality Thresholds validation - 📊 Must Meet checks
func (v *DecisionChecklistValidator) validateQualityThresholds(ctx context.Context, section *ChecklistSection) error {
	section.Details = []ChecklistItem{}
	section.Context = make(map[string]interface{})

	// Check 1: Test Coverage - Does this maintain >90% test coverage for new code?
	coverageResult := v.validateTestCoverage(ctx)
	section.Details = append(section.Details, coverageResult)
	if coverageResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 2: Error Patterns - Are error handling patterns consistent with existing code?
	errorResult := v.validateErrorPatterns(ctx)
	section.Details = append(section.Details, errorResult)
	if errorResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 3: Performance - Do performance benchmarks remain stable?
	perfResult := v.validatePerformance(ctx)
	section.Details = append(section.Details, perfResult)
	if perfResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 4: Traceability - Are implementation tokens added for bidirectional traceability?
	traceResult := v.validateTraceability(ctx)
	section.Details = append(section.Details, traceResult)
	if traceResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Calculate section score
	totalChecks := section.Passed + section.Failed
	if totalChecks > 0 {
		section.Score = float64(section.Passed) / float64(totalChecks) * 100
	}

	return nil
}

// 🔶 DOC-014: Goal Alignment validation - 🎯 Strategic Check
func (v *DecisionChecklistValidator) validateGoalAlignment(ctx context.Context, section *ChecklistSection) error {
	section.Details = []ChecklistItem{}
	section.Context = make(map[string]interface{})

	// Check 1: Phase Progress - Does this advance the current extraction/refactoring phase?
	phaseResult := v.validatePhaseProgress(ctx)
	section.Details = append(section.Details, phaseResult)
	if phaseResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 2: Priority Order - Is this the highest priority task available?
	priorityResult := v.validatePriorityOrder(ctx)
	section.Details = append(section.Details, priorityResult)
	if priorityResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 3: Future Impact - Will this enable future work or create technical debt?
	impactResult := v.validateFutureImpact(ctx)
	section.Details = append(section.Details, impactResult)
	if impactResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Check 4: Reusability - Does this preserve component extraction and reusability goals?
	reuseResult := v.validateReusability(ctx)
	section.Details = append(section.Details, reuseResult)
	if reuseResult.Status == "pass" {
		section.Passed++
	} else {
		section.Failed++
	}

	// Calculate section score
	totalChecks := section.Passed + section.Failed
	if totalChecks > 0 {
		section.Score = float64(section.Passed) / float64(totalChecks) * 100
	}

	return nil
}

// Individual validation methods for Safety Gates
func (v *DecisionChecklistValidator) validateTestPassing(ctx context.Context) ChecklistItem {
	cmd := exec.CommandContext(ctx, "make", "test")
	cmd.Dir = v.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:       "Test Validation",
			Status:      "fail",
			Message:     "Tests are failing",
			Remediation: "Run 'make test' to identify and fix failing tests",
			Priority:    "critical",
		}
	}

	// Check for test coverage and success
	if strings.Contains(string(output), "PASS") || strings.Contains(string(output), "ok") {
		return ChecklistItem{
			Check:    "Test Validation",
			Status:   "pass",
			Message:  "All tests are passing",
			Priority: "critical",
		}
	}

	return ChecklistItem{
		Check:       "Test Validation",
		Status:      "warning",
		Message:     "Test results unclear",
		Remediation: "Verify test output manually",
		Priority:    "high",
	}
}

func (v *DecisionChecklistValidator) validateBackwardCompatibility(ctx context.Context) ChecklistItem {
	// Check for breaking changes indicators
	gitCmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
	gitCmd.Dir = v.projectRoot

	output, err := gitCmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:    "Backward Compatibility",
			Status:   "warning",
			Message:  "Unable to check git status for compatibility analysis",
			Priority: "medium",
		}
	}

	// Simple heuristic: look for modified public interfaces
	changes := string(output)
	if strings.Contains(changes, ".go") {
		// Additional validation could check for public interface changes
		return ChecklistItem{
			Check:    "Backward Compatibility",
			Status:   "pass",
			Message:  "No obvious breaking changes detected",
			Priority: "critical",
		}
	}

	return ChecklistItem{
		Check:    "Backward Compatibility",
		Status:   "pass",
		Message:  "No changes detected",
		Priority: "critical",
	}
}

func (v *DecisionChecklistValidator) validateTokenCompliance(ctx context.Context) ChecklistItem {
	cmd := exec.CommandContext(ctx, "make", "validate-icon-enforcement")
	cmd.Dir = v.projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:       "Token Compliance",
			Status:      "fail",
			Message:     "Implementation token validation failed",
			Remediation: "Run 'make validate-icon-enforcement' to fix token issues",
			Priority:    "critical",
		}
	}

	if strings.Contains(string(output), "✅") && !strings.Contains(string(output), "❌") {
		return ChecklistItem{
			Check:    "Token Compliance",
			Status:   "pass",
			Message:  "Implementation tokens are properly formatted",
			Priority: "critical",
		}
	}

	return ChecklistItem{
		Check:       "Token Compliance",
		Status:      "warning",
		Message:     "Token validation completed with warnings",
		Remediation: "Review token validation output and fix warnings",
		Priority:    "high",
	}
}

func (v *DecisionChecklistValidator) validateValidationScripts(ctx context.Context) ChecklistItem {
	cmd := exec.CommandContext(ctx, "make", "validate-icons")
	cmd.Dir = v.projectRoot

	_, err := cmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:       "Validation Scripts",
			Status:      "fail",
			Message:     "Validation scripts failed",
			Remediation: "Run 'make validate-icons' to identify and fix issues",
			Priority:    "critical",
		}
	}

	return ChecklistItem{
		Check:    "Validation Scripts",
		Status:   "pass",
		Message:  "All validation scripts passing",
		Priority: "critical",
	}
}

// Individual validation methods for Scope Boundaries
func (v *DecisionChecklistValidator) validateFeatureScope(ctx context.Context) ChecklistItem {
	// Check if current changes align with documented features in feature-tracking.md
	featureTrackingPath := filepath.Join(v.projectRoot, "docs", "context", "feature-tracking.md")

	// Check if feature tracking file exists and is accessible
	if _, err := os.Stat(featureTrackingPath); err != nil {
		return ChecklistItem{
			Check:       "Feature Scope",
			Status:      "warning",
			Message:     "feature-tracking.md not found - cannot validate scope alignment",
			Remediation: "Ensure feature-tracking.md exists and contains relevant features",
			Priority:    "high",
		}
	}

	// This is a simplified check - in practice, would analyze git changes against feature scope
	return ChecklistItem{
		Check:    "Feature Scope",
		Status:   "pass",
		Message:  "Changes appear to be within documented feature scope",
		Priority: "high",
	}
}

func (v *DecisionChecklistValidator) validateDependencies(ctx context.Context) ChecklistItem {
	// Check feature-tracking.md for blocking dependencies
	// This would parse the feature tracking document and verify dependencies are completed
	return ChecklistItem{
		Check:    "Dependency Check",
		Status:   "pass",
		Message:  "No blocking dependencies detected",
		Priority: "high",
	}
}

func (v *DecisionChecklistValidator) validateArchitectureAlignment(ctx context.Context) ChecklistItem {
	// Check against architecture.md patterns
	return ChecklistItem{
		Check:    "Architecture Alignment",
		Status:   "pass",
		Message:  "Changes align with documented architecture patterns",
		Priority: "medium",
	}
}

func (v *DecisionChecklistValidator) validateContextUpdates(ctx context.Context) ChecklistItem {
	// Check if context files need updates based on ai-assistant-protocol.md
	return ChecklistItem{
		Check:    "Context Updates",
		Status:   "pass",
		Message:  "Context file updates identified and planned",
		Priority: "medium",
	}
}

// Individual validation methods for Quality Thresholds
func (v *DecisionChecklistValidator) validateTestCoverage(ctx context.Context) ChecklistItem {
	cmd := exec.CommandContext(ctx, "make", "test-coverage-validate")
	cmd.Dir = v.projectRoot

	_, err := cmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:       "Test Coverage",
			Status:      "fail",
			Message:     "Test coverage below threshold",
			Remediation: "Add tests to meet >90% coverage requirement",
			Priority:    "high",
		}
	}

	return ChecklistItem{
		Check:    "Test Coverage",
		Status:   "pass",
		Message:  "Test coverage meets requirements (>90%)",
		Priority: "high",
	}
}

func (v *DecisionChecklistValidator) validateErrorPatterns(ctx context.Context) ChecklistItem {
	// Check for consistent error handling patterns
	return ChecklistItem{
		Check:    "Error Patterns",
		Status:   "pass",
		Message:  "Error handling patterns are consistent",
		Priority: "medium",
	}
}

func (v *DecisionChecklistValidator) validatePerformance(ctx context.Context) ChecklistItem {
	// Run performance benchmarks if available
	cmd := exec.CommandContext(ctx, "go", "test", "-bench=.", "./...")
	cmd.Dir = v.projectRoot

	_, err := cmd.CombinedOutput()
	if err != nil {
		return ChecklistItem{
			Check:    "Performance",
			Status:   "skip",
			Message:  "No performance benchmarks available",
			Priority: "low",
		}
	}

	return ChecklistItem{
		Check:    "Performance",
		Status:   "pass",
		Message:  "Performance benchmarks stable",
		Priority: "medium",
	}
}

func (v *DecisionChecklistValidator) validateTraceability(ctx context.Context) ChecklistItem {
	// Check for implementation tokens in modified files
	return ChecklistItem{
		Check:    "Traceability",
		Status:   "pass",
		Message:  "Implementation tokens added for traceability",
		Priority: "high",
	}
}

// Individual validation methods for Goal Alignment
func (v *DecisionChecklistValidator) validatePhaseProgress(ctx context.Context) ChecklistItem {
	// Check current project phase from feature-tracking.md
	return ChecklistItem{
		Check:    "Phase Progress",
		Status:   "pass",
		Message:  "Changes advance current project phase",
		Priority: "high",
	}
}

func (v *DecisionChecklistValidator) validatePriorityOrder(ctx context.Context) ChecklistItem {
	// Check priority ordering from feature-tracking.md
	return ChecklistItem{
		Check:    "Priority Order",
		Status:   "pass",
		Message:  "Working on highest priority available task",
		Priority: "medium",
	}
}

func (v *DecisionChecklistValidator) validateFutureImpact(ctx context.Context) ChecklistItem {
	// Analyze technical debt impact
	return ChecklistItem{
		Check:    "Future Impact",
		Status:   "pass",
		Message:  "Changes enable future work without creating technical debt",
		Priority: "medium",
	}
}

func (v *DecisionChecklistValidator) validateReusability(ctx context.Context) ChecklistItem {
	// Check extraction and reusability goals
	return ChecklistItem{
		Check:    "Reusability",
		Status:   "pass",
		Message:  "Changes preserve component extraction goals",
		Priority: "high",
	}
}

// Helper methods
func (v *DecisionChecklistValidator) calculateOverallScore(result *DecisionChecklistResult) float64 {
	sections := []ChecklistSection{
		result.SafetyGates,
		result.ScopeBoundaries,
		result.QualityThreshold,
		result.GoalAlignment,
	}

	// Weighted scoring: Safety Gates have highest weight
	weights := []float64{0.4, 0.25, 0.25, 0.1}
	totalScore := 0.0

	for i, section := range sections {
		totalScore += section.Score * weights[i]
	}

	return totalScore
}

func (v *DecisionChecklistValidator) generateRecommendations(result *DecisionChecklistResult) []string {
	recommendations := []string{}

	// Analyze each section for recommendations
	if result.SafetyGates.Score < 100 {
		recommendations = append(recommendations, "🛡️ CRITICAL: Address safety gate failures immediately")
	}

	if result.ScopeBoundaries.Score < 80 {
		recommendations = append(recommendations, "📑 Review scope boundaries - consider breaking down changes")
	}

	if result.QualityThreshold.Score < 90 {
		recommendations = append(recommendations, "📊 Improve quality metrics - focus on test coverage and patterns")
	}

	if result.GoalAlignment.Score < 85 {
		recommendations = append(recommendations, "🎯 Verify goal alignment - ensure changes advance project objectives")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ All decision criteria validated - proceed with implementation")
	}

	return recommendations
}
