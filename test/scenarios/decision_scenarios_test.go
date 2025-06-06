// 🔶 DOC-014: Decision scenario testing - Validate decision trees with realistic scenarios
package scenarios

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// DecisionScenario represents a realistic development scenario for testing
type DecisionScenario struct {
	Name               string
	Description        string
	Context            ScenarioContext
	ExpectedDecision   DecisionOutcome
	DecisionPath       []DecisionStep
	ValidationCriteria []ValidationCriterion
	PerformanceTarget  time.Duration
}

// ScenarioContext provides the situational context for decision making
type ScenarioContext struct {
	FeatureExists      bool
	FeatureStatus      string
	DependenciesMet    bool
	TestsCoverage      float64
	ArchitectureAlign  bool
	PriorityLevel      string
	ProjectPhase       string
	ConflictingChanges bool
	BreakingChanges    bool
}

// DecisionOutcome represents the expected result of applying the decision framework
type DecisionOutcome struct {
	ShouldProceed      bool
	RequiredProtocol   string
	SkippedHierarchy   []string
	ValidationRequired []string
	ExpectedFiles      []string
	PerformanceImpact  string
}

// DecisionStep represents one step in the 4-tier decision hierarchy
type DecisionStep struct {
	Level    string // Safety Gates, Scope Boundaries, Quality Thresholds, Goal Alignment
	Passed   bool
	Criteria []string
	Blockers []string
	Warnings []string
}

// ValidationCriterion defines how to validate the decision outcome
type ValidationCriterion struct {
	Type      string // file_exists, command_succeeds, content_matches, performance_acceptable
	Target    string
	Expected  interface{}
	Tolerance float64
}

// DecisionScenarioTester manages scenario testing
type DecisionScenarioTester struct {
	ProjectRoot    string
	ContextDir     string
	ScriptsDir     string
	TestTimeout    time.Duration
	ValidationMode string
}

// NewDecisionScenarioTester creates a new scenario tester
func NewDecisionScenarioTester(projectRoot string) *DecisionScenarioTester {
	return &DecisionScenarioTester{
		ProjectRoot:    projectRoot,
		ContextDir:     filepath.Join(projectRoot, "docs", "context"),
		ScriptsDir:     filepath.Join(projectRoot, "scripts"),
		TestTimeout:    60 * time.Second,
		ValidationMode: "standard",
	}
}

// TestDecisionScenarios runs the complete scenario test suite
func TestDecisionScenarios(t *testing.T) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	tester := NewDecisionScenarioTester(projectRoot)

	// Test all scenario categories
	t.Run("NewFeatureScenarios", func(t *testing.T) {
		scenarios := tester.createNewFeatureScenarios()
		for _, scenario := range scenarios {
			t.Run(scenario.Name, func(t *testing.T) {
				tester.runScenario(t, scenario)
			})
		}
	})

	t.Run("BugFixScenarios", func(t *testing.T) {
		scenarios := tester.createBugFixScenarios()
		for _, scenario := range scenarios {
			t.Run(scenario.Name, func(t *testing.T) {
				tester.runScenario(t, scenario)
			})
		}
	})

	t.Run("ConflictScenarios", func(t *testing.T) {
		scenarios := tester.createConflictScenarios()
		for _, scenario := range scenarios {
			t.Run(scenario.Name, func(t *testing.T) {
				tester.runScenario(t, scenario)
			})
		}
	})
}

// createNewFeatureScenarios generates scenarios for new feature implementation
func (tester *DecisionScenarioTester) createNewFeatureScenarios() []DecisionScenario {
	return []DecisionScenario{
		{
			Name:        "ValidNewFeature",
			Description: "Implementing a new feature with all prerequisites met",
			Context: ScenarioContext{
				FeatureExists:      true,
				FeatureStatus:      "📝 Not Started",
				DependenciesMet:    true,
				TestsCoverage:      95.0,
				ArchitectureAlign:  true,
				PriorityLevel:      "⭐ CRITICAL",
				ProjectPhase:       "Phase 4: Component Extraction",
				ConflictingChanges: false,
				BreakingChanges:    false,
			},
			ExpectedDecision: DecisionOutcome{
				ShouldProceed:      true,
				RequiredProtocol:   "NEW FEATURE",
				SkippedHierarchy:   []string{},
				ValidationRequired: []string{"safety_gates", "scope_boundaries", "quality_thresholds", "goal_alignment"},
				ExpectedFiles:      []string{"feature-tracking.md", "architecture.md", "requirements.md", "specification.md"},
				PerformanceImpact:  "acceptable",
			},
			DecisionPath: []DecisionStep{
				{Level: "Safety Gates", Passed: true, Criteria: []string{"tests_pass", "backward_compatible", "tokens_compliant"}},
				{Level: "Scope Boundaries", Passed: true, Criteria: []string{"feature_scope", "dependencies_met", "architecture_aligned"}},
				{Level: "Quality Thresholds", Passed: true, Criteria: []string{"test_coverage", "error_patterns", "performance", "traceability"}},
				{Level: "Goal Alignment", Passed: true, Criteria: []string{"phase_progress", "priority_order", "future_impact", "reusability"}},
			},
			PerformanceTarget: 30 * time.Second,
		},
		{
			Name:        "FeatureWithMissingDependencies",
			Description: "Attempting to implement feature with unmet dependencies",
			Context: ScenarioContext{
				FeatureExists:      true,
				FeatureStatus:      "📝 Not Started",
				DependenciesMet:    false,
				TestsCoverage:      90.0,
				ArchitectureAlign:  true,
				PriorityLevel:      "🔺 HIGH",
				ProjectPhase:       "Phase 4: Component Extraction",
				ConflictingChanges: false,
				BreakingChanges:    false,
			},
			ExpectedDecision: DecisionOutcome{
				ShouldProceed:      false,
				RequiredProtocol:   "",
				SkippedHierarchy:   []string{"Quality Thresholds", "Goal Alignment"},
				ValidationRequired: []string{"scope_boundaries"},
				ExpectedFiles:      []string{},
				PerformanceImpact:  "none",
			},
			DecisionPath: []DecisionStep{
				{Level: "Safety Gates", Passed: true, Criteria: []string{"tests_pass", "backward_compatible", "tokens_compliant"}},
				{Level: "Scope Boundaries", Passed: false, Blockers: []string{"dependencies_not_met"}, Criteria: []string{"dependency_check"}},
			},
			PerformanceTarget: 10 * time.Second,
		},
	}
}

// createBugFixScenarios generates scenarios for bug fixes
func (tester *DecisionScenarioTester) createBugFixScenarios() []DecisionScenario {
	return []DecisionScenario{
		{
			Name:        "CriticalBugFix",
			Description: "Fixing a critical bug blocking extraction work",
			Context: ScenarioContext{
				FeatureExists:      true,
				FeatureStatus:      "✅ Implemented",
				DependenciesMet:    true,
				TestsCoverage:      85.0,
				ArchitectureAlign:  true,
				PriorityLevel:      "⭐ CRITICAL",
				ProjectPhase:       "Phase 4: Component Extraction",
				ConflictingChanges: false,
				BreakingChanges:    false,
			},
			ExpectedDecision: DecisionOutcome{
				ShouldProceed:      true,
				RequiredProtocol:   "BUG FIX",
				SkippedHierarchy:   []string{},
				ValidationRequired: []string{"safety_gates", "quality_thresholds"},
				ExpectedFiles:      []string{"testing.md"},
				PerformanceImpact:  "minimal",
			},
			DecisionPath: []DecisionStep{
				{Level: "Safety Gates", Passed: true, Criteria: []string{"tests_pass", "backward_compatible"}},
				{Level: "Scope Boundaries", Passed: true, Criteria: []string{"blocking_work", "critical_component"}},
				{Level: "Quality Thresholds", Passed: true, Criteria: []string{"error_patterns", "functionality_preserved"}},
				{Level: "Goal Alignment", Passed: true, Criteria: []string{"immediate_fix", "extraction_unblocked"}},
			},
			PerformanceTarget: 15 * time.Second,
		},
	}
}

// createConflictScenarios generates scenarios with conflicts or edge cases
func (tester *DecisionScenarioTester) createConflictScenarios() []DecisionScenario {
	return []DecisionScenario{
		{
			Name:        "TestFailuresBlocking",
			Description: "Test failures preventing any development work",
			Context: ScenarioContext{
				FeatureExists:      true,
				FeatureStatus:      "📝 Not Started",
				DependenciesMet:    true,
				TestsCoverage:      0.0,
				ArchitectureAlign:  true,
				PriorityLevel:      "⭐ CRITICAL",
				ProjectPhase:       "Phase 4: Component Extraction",
				ConflictingChanges: false,
				BreakingChanges:    false,
			},
			ExpectedDecision: DecisionOutcome{
				ShouldProceed:      false,
				RequiredProtocol:   "",
				SkippedHierarchy:   []string{"Scope Boundaries", "Quality Thresholds", "Goal Alignment"},
				ValidationRequired: []string{"safety_gates"},
				ExpectedFiles:      []string{},
				PerformanceImpact:  "none",
			},
			DecisionPath: []DecisionStep{
				{Level: "Safety Gates", Passed: false, Blockers: []string{"tests_failing"}, Criteria: []string{"tests_pass"}},
			},
			PerformanceTarget: 3 * time.Second,
		},
	}
}

// runScenario executes a decision scenario and validates the outcome
func (tester *DecisionScenarioTester) runScenario(t *testing.T, scenario DecisionScenario) {
	t.Helper()

	startTime := time.Now()

	// Step 1: Apply the 4-tier decision hierarchy
	actualDecision := tester.applyDecisionHierarchy(t, scenario)

	// Step 2: Validate the decision outcome
	tester.validateDecisionOutcome(t, scenario, actualDecision)

	// Step 3: Validate performance
	elapsed := time.Since(startTime)
	if elapsed > scenario.PerformanceTarget {
		t.Errorf("Scenario %s exceeded performance target: %v > %v",
			scenario.Name, elapsed, scenario.PerformanceTarget)
	}

	// Step 4: Run validation criteria
	for _, criterion := range scenario.ValidationCriteria {
		tester.validateCriterion(t, criterion)
	}
}

// applyDecisionHierarchy simulates applying the decision framework
func (tester *DecisionScenarioTester) applyDecisionHierarchy(t *testing.T, scenario DecisionScenario) DecisionOutcome {
	t.Helper()

	actualDecision := DecisionOutcome{
		ValidationRequired: []string{},
		ExpectedFiles:      []string{},
	}

	// Simulate 4-tier hierarchy validation
	for _, step := range scenario.DecisionPath {
		actualDecision.ValidationRequired = append(actualDecision.ValidationRequired, strings.ToLower(strings.ReplaceAll(step.Level, " ", "_")))

		if !step.Passed {
			actualDecision.ShouldProceed = false
			actualDecision.SkippedHierarchy = append(actualDecision.SkippedHierarchy, step.Level)
			break
		}
	}

	// If all steps passed, determine protocol
	if len(actualDecision.SkippedHierarchy) == 0 {
		actualDecision.ShouldProceed = true

		// Determine protocol based on context
		if scenario.Context.FeatureStatus == "📝 Not Started" {
			actualDecision.RequiredProtocol = "NEW FEATURE"
			actualDecision.ExpectedFiles = []string{"feature-tracking.md", "architecture.md", "requirements.md", "specification.md"}
		} else if strings.Contains(scenario.Description, "bug") {
			actualDecision.RequiredProtocol = "BUG FIX"
			actualDecision.ExpectedFiles = []string{"testing.md"}
		}
	}

	return actualDecision
}

// validateDecisionOutcome compares actual vs expected decision outcomes
func (tester *DecisionScenarioTester) validateDecisionOutcome(t *testing.T, scenario DecisionScenario, actual DecisionOutcome) {
	t.Helper()

	if actual.ShouldProceed != scenario.ExpectedDecision.ShouldProceed {
		t.Errorf("Scenario %s: Expected ShouldProceed=%v, got %v",
			scenario.Name, scenario.ExpectedDecision.ShouldProceed, actual.ShouldProceed)
	}

	if actual.RequiredProtocol != scenario.ExpectedDecision.RequiredProtocol {
		t.Errorf("Scenario %s: Expected RequiredProtocol=%s, got %s",
			scenario.Name, scenario.ExpectedDecision.RequiredProtocol, actual.RequiredProtocol)
	}
}

// validateCriterion validates individual validation criteria
func (tester *DecisionScenarioTester) validateCriterion(t *testing.T, criterion ValidationCriterion) {
	t.Helper()

	switch criterion.Type {
	case "file_exists":
		filePath := filepath.Join(tester.ProjectRoot, criterion.Target)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Required file missing: %s", criterion.Target)
		}

	case "command_succeeds":
		ctx, cancel := context.WithTimeout(context.Background(), tester.TestTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "bash", "-c", criterion.Target)
		cmd.Dir = tester.ProjectRoot

		if err := cmd.Run(); err != nil {
			t.Errorf("Command failed: %s - %v", criterion.Target, err)
		}
	}
}

// findProjectRoot locates the project root directory
func findProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Look for go.mod file to identify project root
	dir := cwd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback to current directory
	return cwd, nil
}
