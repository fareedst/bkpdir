package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// GoalAlignmentTestSuite focuses on goal alignment rate measurement validation
type GoalAlignmentTestSuite struct {
	WorkspaceRoot        string
	MetricsScriptPath    string
	ValidationScriptPath string
	TestScenariosPath    string
	AlignmentTargets     AlignmentTargets
}

// AlignmentTargets defines expected alignment targets
type AlignmentTargets struct {
	OverallTarget       float64            `json:"overall_target"` // ≥95% goal alignment
	HierarchyTargets    map[string]float64 `json:"hierarchy_targets"`
	ProtocolTargets     map[string]float64 `json:"protocol_targets"`
	ProjectObjectives   []string           `json:"project_objectives"`
	MeasurementCriteria []string           `json:"measurement_criteria"`
}

// GoalAlignmentScenario represents a test scenario for goal alignment
type GoalAlignmentScenario struct {
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	ProjectObjectives []ProjectObjective   `json:"project_objectives"`
	DecisionOutcomes  []DecisionOutcome    `json:"decision_outcomes"`
	AlignmentCriteria []AlignmentCriterion `json:"alignment_criteria"`
	ExpectedAlignment float64              `json:"expected_alignment"`
	AlignmentType     string               `json:"alignment_type"`
}

// ProjectObjective represents a project goal or objective
type ProjectObjective struct {
	ObjectiveID     string  `json:"objective_id"`
	ObjectiveName   string  `json:"objective_name"`
	Priority        string  `json:"priority"` // HIGH, MEDIUM, LOW
	Success_Metric  string  `json:"success_metric"`
	Target_Value    string  `json:"target_value"`
	Weight          float64 `json:"weight"`          // 0.0-1.0
	Alignment_Score float64 `json:"alignment_score"` // 0.0-100.0
}

// DecisionOutcome represents the outcome of a decision framework application
type DecisionOutcome struct {
	DecisionID      string  `json:"decision_id"`
	DecisionType    string  `json:"decision_type"`
	Hierarchy_Level string  `json:"hierarchy_level"`
	Outcome_Result  string  `json:"outcome_result"`
	Goal_Impact     string  `json:"goal_impact"`     // POSITIVE, NEUTRAL, NEGATIVE
	Alignment_Score float64 `json:"alignment_score"` // 0.0-100.0
}

// AlignmentCriterion represents criteria for measuring goal alignment
type AlignmentCriterion struct {
	CriterionID     string  `json:"criterion_id"`
	CriterionName   string  `json:"criterion_name"`
	MeasurementType string  `json:"measurement_type"` // OBJECTIVE_BASED, OUTCOME_BASED, IMPACT_BASED
	Required_Value  string  `json:"required_value"`
	Actual_Value    string  `json:"actual_value"`
	Weight          float64 `json:"weight"`          // 0.0-1.0
	Alignment_Score float64 `json:"alignment_score"` // 0.0-100.0
	Validation_Pass bool    `json:"validation_pass"`
}

// GoalAlignmentResult represents goal alignment measurement result
type GoalAlignmentResult struct {
	ScenarioName           string  `json:"scenario_name"`
	OverallAlignmentRate   float64 `json:"overall_alignment_rate"`
	ObjectiveAlignmentRate float64 `json:"objective_alignment_rate"`
	OutcomeAlignmentRate   float64 `json:"outcome_alignment_rate"`
	CriteriaAlignmentRate  float64 `json:"criteria_alignment_rate"`
	AlignmentAccuracy      float64 `json:"alignment_accuracy"`
	MeasurementReliability float64 `json:"measurement_reliability"`
	ValidationSuccess      bool    `json:"validation_success"`
	AnalysisTimestamp      string  `json:"analysis_timestamp"`
}

// TestDOC014GoalAlignmentMeasurement - main test for goal alignment measurement validation
func TestDOC014GoalAlignmentMeasurement(t *testing.T) {
	suite := &GoalAlignmentTestSuite{}

	// Initialize test suite
	if err := suite.Initialize(); err != nil {
		t.Fatalf("Failed to initialize goal alignment test suite: %v", err)
	}

	t.Logf("🎯 Starting DOC-014 Goal Alignment Measurement Test Suite")
	t.Logf("📊 Workspace: %s", suite.WorkspaceRoot)
	t.Logf("📈 Metrics Script: %s", suite.MetricsScriptPath)
	t.Logf("🎯 Target Goal Alignment: %.1f%%", suite.AlignmentTargets.OverallTarget)

	// Create test scenarios
	scenarios := suite.CreateGoalAlignmentScenarios()

	// Run goal alignment measurement tests
	t.Run("HighAlignmentScenario", func(t *testing.T) {
		suite.TestHighAlignmentScenario(t, scenarios[0])
	})
	t.Run("ModerateAlignmentScenario", func(t *testing.T) {
		suite.TestModerateAlignmentScenario(t, scenarios[1])
	})
	t.Run("LowAlignmentScenario", func(t *testing.T) {
		suite.TestLowAlignmentScenario(t, scenarios[2])
	})
	t.Run("MixedAlignmentScenario", func(t *testing.T) {
		suite.TestMixedAlignmentScenario(t, scenarios[3])
	})
	t.Run("ObjectiveBasedAlignment", suite.TestObjectiveBasedAlignment)
	t.Run("OutcomeBasedAlignment", suite.TestOutcomeBasedAlignment)
	t.Run("ImpactBasedAlignment", suite.TestImpactBasedAlignment)
	t.Run("AlignmentMeasurementAccuracy", suite.TestAlignmentMeasurementAccuracy)
	t.Run("AlignmentTargetValidation", suite.TestAlignmentTargetValidation)
	t.Run("AlignmentConsistency", suite.TestAlignmentConsistency)

	t.Logf("✅ DOC-014 Goal Alignment Measurement Test Suite Completed")
}

// Initialize sets up the goal alignment test environment
func (suite *GoalAlignmentTestSuite) Initialize() error {
	// Get workspace root
	workspaceRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get workspace root: %w", err)
	}

	// Navigate to workspace root (assuming we're in test/metrics)
	for !strings.HasSuffix(workspaceRoot, "bkpdir") && workspaceRoot != "/" {
		workspaceRoot = filepath.Dir(workspaceRoot)
	}

	suite.WorkspaceRoot = workspaceRoot
	suite.MetricsScriptPath = filepath.Join(workspaceRoot, "scripts", "track-decision-metrics.sh")
	suite.ValidationScriptPath = filepath.Join(workspaceRoot, "scripts", "validate-decision-framework.sh")
	suite.TestScenariosPath = filepath.Join(workspaceRoot, "test", "metrics", "goal_alignment_scenarios")

	// Set up alignment targets
	suite.AlignmentTargets = AlignmentTargets{
		OverallTarget: 95.0, // Target ≥95% goal alignment
		HierarchyTargets: map[string]float64{
			"Safety Gates":       95.0,
			"Scope Boundaries":   90.0,
			"Quality Thresholds": 85.0,
			"Goal Alignment":     98.0,
		},
		ProtocolTargets: map[string]float64{
			"NEW_FEATURE":   95.0,
			"MODIFICATION":  90.0,
			"BUG_FIX":       85.0,
			"CONFIG_CHANGE": 90.0,
			"REFACTORING":   85.0,
			"DOCUMENTATION": 80.0,
			"TESTING":       95.0,
			"DEPLOYMENT":    90.0,
		},
		ProjectObjectives: []string{
			"Improve AI assistant decision quality",
			"Reduce development rework rates",
			"Enhance documentation consistency",
			"Strengthen validation frameworks",
			"Optimize development workflows",
		},
		MeasurementCriteria: []string{
			"Objective achievement rate",
			"Decision outcome quality",
			"Framework impact assessment",
			"User satisfaction metrics",
			"Validation effectiveness",
		},
	}

	// Create test scenarios directory
	return os.MkdirAll(suite.TestScenariosPath, 0755)
}

// CreateGoalAlignmentScenarios creates test scenarios for goal alignment measurement
func (suite *GoalAlignmentTestSuite) CreateGoalAlignmentScenarios() []GoalAlignmentScenario {
	scenarios := []GoalAlignmentScenario{
		{
			Name:        "HighAlignmentNewFeature",
			Description: "New feature development with high goal alignment",
			ProjectObjectives: []ProjectObjective{
				{ObjectiveID: "OBJ001", ObjectiveName: "Improve AI decision quality", Priority: "HIGH", Success_Metric: "Decision accuracy", Target_Value: "≥95%", Weight: 0.3, Alignment_Score: 98.0},
				{ObjectiveID: "OBJ002", ObjectiveName: "Reduce rework rates", Priority: "HIGH", Success_Metric: "Rework percentage", Target_Value: "≤5%", Weight: 0.25, Alignment_Score: 95.0},
				{ObjectiveID: "OBJ003", ObjectiveName: "Enhance documentation", Priority: "MEDIUM", Success_Metric: "Documentation coverage", Target_Value: "≥90%", Weight: 0.2, Alignment_Score: 92.0},
				{ObjectiveID: "OBJ004", ObjectiveName: "Strengthen validation", Priority: "HIGH", Success_Metric: "Validation coverage", Target_Value: "≥95%", Weight: 0.25, Alignment_Score: 97.0},
			},
			DecisionOutcomes: []DecisionOutcome{
				{DecisionID: "DEC001", DecisionType: "NEW_FEATURE", Hierarchy_Level: "Safety Gates", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 95.0},
				{DecisionID: "DEC002", DecisionType: "NEW_FEATURE", Hierarchy_Level: "Scope Boundaries", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 93.0},
				{DecisionID: "DEC003", DecisionType: "NEW_FEATURE", Hierarchy_Level: "Quality Thresholds", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 96.0},
				{DecisionID: "DEC004", DecisionType: "NEW_FEATURE", Hierarchy_Level: "Goal Alignment", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 98.0},
			},
			AlignmentCriteria: []AlignmentCriterion{
				{CriterionID: "CRIT001", CriterionName: "Objective Achievement", MeasurementType: "OBJECTIVE_BASED", Required_Value: "≥95%", Actual_Value: "96%", Weight: 0.4, Alignment_Score: 96.0, Validation_Pass: true},
				{CriterionID: "CRIT002", CriterionName: "Outcome Quality", MeasurementType: "OUTCOME_BASED", Required_Value: "HIGH", Actual_Value: "HIGH", Weight: 0.3, Alignment_Score: 95.0, Validation_Pass: true},
				{CriterionID: "CRIT003", CriterionName: "Impact Assessment", MeasurementType: "IMPACT_BASED", Required_Value: "POSITIVE", Actual_Value: "POSITIVE", Weight: 0.3, Alignment_Score: 97.0, Validation_Pass: true},
			},
			ExpectedAlignment: 95.8, // Weighted average
			AlignmentType:     "HIGH_ALIGNMENT",
		},
		{
			Name:        "ModerateAlignmentModification",
			Description: "Feature modification with moderate goal alignment",
			ProjectObjectives: []ProjectObjective{
				{ObjectiveID: "OBJ001", ObjectiveName: "Improve AI decision quality", Priority: "HIGH", Success_Metric: "Decision accuracy", Target_Value: "≥95%", Weight: 0.3, Alignment_Score: 88.0},
				{ObjectiveID: "OBJ002", ObjectiveName: "Reduce rework rates", Priority: "HIGH", Success_Metric: "Rework percentage", Target_Value: "≤5%", Weight: 0.25, Alignment_Score: 82.0},
				{ObjectiveID: "OBJ003", ObjectiveName: "Enhance documentation", Priority: "MEDIUM", Success_Metric: "Documentation coverage", Target_Value: "≥90%", Weight: 0.2, Alignment_Score: 78.0},
				{ObjectiveID: "OBJ004", ObjectiveName: "Strengthen validation", Priority: "HIGH", Success_Metric: "Validation coverage", Target_Value: "≥95%", Weight: 0.25, Alignment_Score: 85.0},
			},
			DecisionOutcomes: []DecisionOutcome{
				{DecisionID: "DEC001", DecisionType: "MODIFICATION", Hierarchy_Level: "Safety Gates", Outcome_Result: "APPROVED", Goal_Impact: "NEUTRAL", Alignment_Score: 80.0},
				{DecisionID: "DEC002", DecisionType: "MODIFICATION", Hierarchy_Level: "Scope Boundaries", Outcome_Result: "CONDITIONAL", Goal_Impact: "NEUTRAL", Alignment_Score: 75.0},
				{DecisionID: "DEC003", DecisionType: "MODIFICATION", Hierarchy_Level: "Quality Thresholds", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 85.0},
				{DecisionID: "DEC004", DecisionType: "MODIFICATION", Hierarchy_Level: "Goal Alignment", Outcome_Result: "CONDITIONAL", Goal_Impact: "NEUTRAL", Alignment_Score: 78.0},
			},
			AlignmentCriteria: []AlignmentCriterion{
				{CriterionID: "CRIT001", CriterionName: "Objective Achievement", MeasurementType: "OBJECTIVE_BASED", Required_Value: "≥95%", Actual_Value: "83%", Weight: 0.4, Alignment_Score: 83.0, Validation_Pass: false},
				{CriterionID: "CRIT002", CriterionName: "Outcome Quality", MeasurementType: "OUTCOME_BASED", Required_Value: "HIGH", Actual_Value: "MEDIUM", Weight: 0.3, Alignment_Score: 79.0, Validation_Pass: false},
				{CriterionID: "CRIT003", CriterionName: "Impact Assessment", MeasurementType: "IMPACT_BASED", Required_Value: "POSITIVE", Actual_Value: "NEUTRAL", Weight: 0.3, Alignment_Score: 80.0, Validation_Pass: false},
			},
			ExpectedAlignment: 80.9, // Weighted average
			AlignmentType:     "MODERATE_ALIGNMENT",
		},
		{
			Name:        "LowAlignmentBugFix",
			Description: "Bug fix with low goal alignment",
			ProjectObjectives: []ProjectObjective{
				{ObjectiveID: "OBJ001", ObjectiveName: "Improve AI decision quality", Priority: "HIGH", Success_Metric: "Decision accuracy", Target_Value: "≥95%", Weight: 0.3, Alignment_Score: 65.0},
				{ObjectiveID: "OBJ002", ObjectiveName: "Reduce rework rates", Priority: "HIGH", Success_Metric: "Rework percentage", Target_Value: "≤5%", Weight: 0.25, Alignment_Score: 60.0},
				{ObjectiveID: "OBJ003", ObjectiveName: "Enhance documentation", Priority: "MEDIUM", Success_Metric: "Documentation coverage", Target_Value: "≥90%", Weight: 0.2, Alignment_Score: 55.0},
				{ObjectiveID: "OBJ004", ObjectiveName: "Strengthen validation", Priority: "HIGH", Success_Metric: "Validation coverage", Target_Value: "≥95%", Weight: 0.25, Alignment_Score: 58.0},
			},
			DecisionOutcomes: []DecisionOutcome{
				{DecisionID: "DEC001", DecisionType: "BUG_FIX", Hierarchy_Level: "Safety Gates", Outcome_Result: "REJECTED", Goal_Impact: "NEGATIVE", Alignment_Score: 50.0},
				{DecisionID: "DEC002", DecisionType: "BUG_FIX", Hierarchy_Level: "Scope Boundaries", Outcome_Result: "CONDITIONAL", Goal_Impact: "NEGATIVE", Alignment_Score: 55.0},
				{DecisionID: "DEC003", DecisionType: "BUG_FIX", Hierarchy_Level: "Quality Thresholds", Outcome_Result: "CONDITIONAL", Goal_Impact: "NEUTRAL", Alignment_Score: 60.0},
				{DecisionID: "DEC004", DecisionType: "BUG_FIX", Hierarchy_Level: "Goal Alignment", Outcome_Result: "REJECTED", Goal_Impact: "NEGATIVE", Alignment_Score: 45.0},
			},
			AlignmentCriteria: []AlignmentCriterion{
				{CriterionID: "CRIT001", CriterionName: "Objective Achievement", MeasurementType: "OBJECTIVE_BASED", Required_Value: "≥95%", Actual_Value: "59%", Weight: 0.4, Alignment_Score: 59.0, Validation_Pass: false},
				{CriterionID: "CRIT002", CriterionName: "Outcome Quality", MeasurementType: "OUTCOME_BASED", Required_Value: "HIGH", Actual_Value: "LOW", Weight: 0.3, Alignment_Score: 52.0, Validation_Pass: false},
				{CriterionID: "CRIT003", CriterionName: "Impact Assessment", MeasurementType: "IMPACT_BASED", Required_Value: "POSITIVE", Actual_Value: "NEGATIVE", Weight: 0.3, Alignment_Score: 48.0, Validation_Pass: false},
			},
			ExpectedAlignment: 54.7, // Weighted average
			AlignmentType:     "LOW_ALIGNMENT",
		},
		{
			Name:        "MixedAlignmentRefactoring",
			Description: "Refactoring with mixed goal alignment levels",
			ProjectObjectives: []ProjectObjective{
				{ObjectiveID: "OBJ001", ObjectiveName: "Improve AI decision quality", Priority: "HIGH", Success_Metric: "Decision accuracy", Target_Value: "≥95%", Weight: 0.3, Alignment_Score: 92.0},
				{ObjectiveID: "OBJ002", ObjectiveName: "Reduce rework rates", Priority: "HIGH", Success_Metric: "Rework percentage", Target_Value: "≤5%", Weight: 0.25, Alignment_Score: 70.0},
				{ObjectiveID: "OBJ003", ObjectiveName: "Enhance documentation", Priority: "MEDIUM", Success_Metric: "Documentation coverage", Target_Value: "≥90%", Weight: 0.2, Alignment_Score: 95.0},
				{ObjectiveID: "OBJ004", ObjectiveName: "Strengthen validation", Priority: "HIGH", Success_Metric: "Validation coverage", Target_Value: "≥95%", Weight: 0.25, Alignment_Score: 65.0},
			},
			DecisionOutcomes: []DecisionOutcome{
				{DecisionID: "DEC001", DecisionType: "REFACTORING", Hierarchy_Level: "Safety Gates", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 90.0},
				{DecisionID: "DEC002", DecisionType: "REFACTORING", Hierarchy_Level: "Scope Boundaries", Outcome_Result: "CONDITIONAL", Goal_Impact: "NEUTRAL", Alignment_Score: 75.0},
				{DecisionID: "DEC003", DecisionType: "REFACTORING", Hierarchy_Level: "Quality Thresholds", Outcome_Result: "REJECTED", Goal_Impact: "NEGATIVE", Alignment_Score: 55.0},
				{DecisionID: "DEC004", DecisionType: "REFACTORING", Hierarchy_Level: "Goal Alignment", Outcome_Result: "APPROVED", Goal_Impact: "POSITIVE", Alignment_Score: 88.0},
			},
			AlignmentCriteria: []AlignmentCriterion{
				{CriterionID: "CRIT001", CriterionName: "Objective Achievement", MeasurementType: "OBJECTIVE_BASED", Required_Value: "≥95%", Actual_Value: "81%", Weight: 0.4, Alignment_Score: 81.0, Validation_Pass: false},
				{CriterionID: "CRIT002", CriterionName: "Outcome Quality", MeasurementType: "OUTCOME_BASED", Required_Value: "HIGH", Actual_Value: "MIXED", Weight: 0.3, Alignment_Score: 77.0, Validation_Pass: false},
				{CriterionID: "CRIT003", CriterionName: "Impact Assessment", MeasurementType: "IMPACT_BASED", Required_Value: "POSITIVE", Actual_Value: "MIXED", Weight: 0.3, Alignment_Score: 75.0, Validation_Pass: false},
			},
			ExpectedAlignment: 78.1, // Weighted average
			AlignmentType:     "MIXED_ALIGNMENT",
		},
	}

	return scenarios
}

// TestHighAlignmentScenario tests goal alignment measurement for high alignment scenarios
func (suite *GoalAlignmentTestSuite) TestHighAlignmentScenario(t *testing.T, scenario GoalAlignmentScenario) {
	t.Logf("🎯 Testing high alignment scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateAlignmentScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create alignment scenario file: %v", err)
	}

	// Execute goal alignment measurement
	result, err := suite.ExecuteGoalAlignmentAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute goal alignment analysis: %v", err)
	}

	// Validate alignment rate meets high standards (≥90%)
	if result.OverallAlignmentRate < 90.0 {
		t.Errorf("High alignment scenario rate %.1f%% below expected (≥90%%)", result.OverallAlignmentRate)
	}

	// Validate measurement accuracy
	tolerance := 3.0 // 3% tolerance for high alignment scenarios
	expectedAlignment := scenario.ExpectedAlignment
	actualAlignment := result.OverallAlignmentRate

	if actualAlignment < expectedAlignment-tolerance || actualAlignment > expectedAlignment+tolerance {
		t.Errorf("Goal alignment measurement error: expected %.1f%%, got %.1f%%", expectedAlignment, actualAlignment)
	}

	// Validate alignment components are consistently high
	if result.ObjectiveAlignmentRate < 85.0 || result.OutcomeAlignmentRate < 85.0 || result.CriteriaAlignmentRate < 85.0 {
		t.Errorf("High alignment scenario has inconsistent component alignment rates")
	}

	// Validate measurement reliability
	if result.MeasurementReliability < 90.0 {
		t.Errorf("High alignment scenario measurement reliability %.1f%% below target (≥90%%)", result.MeasurementReliability)
	}

	t.Logf("✅ High alignment scenario validation passed: %.1f%% alignment (expected: %.1f%%, reliability: %.1f%%)",
		actualAlignment, expectedAlignment, result.MeasurementReliability)
}

// TestModerateAlignmentScenario tests goal alignment measurement for moderate alignment scenarios
func (suite *GoalAlignmentTestSuite) TestModerateAlignmentScenario(t *testing.T, scenario GoalAlignmentScenario) {
	t.Logf("🔶 Testing moderate alignment scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateAlignmentScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create alignment scenario file: %v", err)
	}

	// Execute goal alignment measurement
	result, err := suite.ExecuteGoalAlignmentAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute goal alignment analysis: %v", err)
	}

	// Validate alignment rate is in moderate range (70-90%)
	if result.OverallAlignmentRate < 70.0 || result.OverallAlignmentRate >= 90.0 {
		t.Errorf("Moderate alignment scenario rate %.1f%% outside expected range (70-90%%)", result.OverallAlignmentRate)
	}

	// Validate measurement accuracy
	tolerance := 4.0 // 4% tolerance for moderate alignment scenarios
	expectedAlignment := scenario.ExpectedAlignment
	actualAlignment := result.OverallAlignmentRate

	if actualAlignment < expectedAlignment-tolerance || actualAlignment > expectedAlignment+tolerance {
		t.Errorf("Goal alignment measurement error: expected %.1f%%, got %.1f%%", expectedAlignment, actualAlignment)
	}

	// Validate measurement reliability is reasonable for moderate scenarios
	if result.MeasurementReliability < 85.0 {
		t.Errorf("Moderate alignment scenario measurement reliability %.1f%% below target (≥85%%)", result.MeasurementReliability)
	}

	t.Logf("✅ Moderate alignment scenario validation passed: %.1f%% alignment (expected: %.1f%%, reliability: %.1f%%)",
		actualAlignment, expectedAlignment, result.MeasurementReliability)
}

// TestLowAlignmentScenario tests goal alignment measurement for low alignment scenarios
func (suite *GoalAlignmentTestSuite) TestLowAlignmentScenario(t *testing.T, scenario GoalAlignmentScenario) {
	t.Logf("🔻 Testing low alignment scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateAlignmentScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create alignment scenario file: %v", err)
	}

	// Execute goal alignment measurement
	result, err := suite.ExecuteGoalAlignmentAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute goal alignment analysis: %v", err)
	}

	// Validate alignment rate is in low range (<70%)
	if result.OverallAlignmentRate >= 70.0 {
		t.Errorf("Low alignment scenario rate %.1f%% above expected threshold (<70%%)", result.OverallAlignmentRate)
	}

	// Validate measurement accuracy
	tolerance := 5.0 // 5% tolerance for low alignment scenarios
	expectedAlignment := scenario.ExpectedAlignment
	actualAlignment := result.OverallAlignmentRate

	if actualAlignment < expectedAlignment-tolerance || actualAlignment > expectedAlignment+tolerance {
		t.Errorf("Goal alignment measurement error: expected %.1f%%, got %.1f%%", expectedAlignment, actualAlignment)
	}

	// Validate measurement reliability is maintained even for low alignment
	if result.MeasurementReliability < 80.0 {
		t.Errorf("Low alignment scenario measurement reliability %.1f%% below target (≥80%%)", result.MeasurementReliability)
	}

	t.Logf("✅ Low alignment scenario validation passed: %.1f%% alignment (expected: %.1f%%, reliability: %.1f%%)",
		actualAlignment, expectedAlignment, result.MeasurementReliability)
}

// TestMixedAlignmentScenario tests goal alignment measurement for mixed alignment scenarios
func (suite *GoalAlignmentTestSuite) TestMixedAlignmentScenario(t *testing.T, scenario GoalAlignmentScenario) {
	t.Logf("🔀 Testing mixed alignment scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateAlignmentScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create alignment scenario file: %v", err)
	}

	// Execute goal alignment measurement
	result, err := suite.ExecuteGoalAlignmentAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute goal alignment analysis: %v", err)
	}

	// Validate measurement accuracy
	tolerance := 5.0 // 5% tolerance for mixed alignment scenarios
	expectedAlignment := scenario.ExpectedAlignment
	actualAlignment := result.OverallAlignmentRate

	if actualAlignment < expectedAlignment-tolerance || actualAlignment > expectedAlignment+tolerance {
		t.Errorf("Goal alignment measurement error: expected %.1f%%, got %.1f%%", expectedAlignment, actualAlignment)
	}

	// Validate there is variation in component alignment rates
	alignmentRates := []float64{result.ObjectiveAlignmentRate, result.OutcomeAlignmentRate, result.CriteriaAlignmentRate}

	// Calculate variance to ensure mixed nature
	var sum, variance float64
	for _, rate := range alignmentRates {
		sum += rate
	}
	mean := sum / float64(len(alignmentRates))

	for _, rate := range alignmentRates {
		variance += (rate - mean) * (rate - mean)
	}
	variance /= float64(len(alignmentRates))

	if variance < 50.0 { // Expect significant variation in mixed alignment
		t.Errorf("Mixed alignment scenario shows insufficient variation (variance: %.1f)", variance)
	}

	t.Logf("✅ Mixed alignment scenario validation passed: %.1f%% alignment (expected: %.1f%%, variance: %.1f)",
		actualAlignment, expectedAlignment, variance)
}

// TestObjectiveBasedAlignment tests objective-based alignment measurement
func (suite *GoalAlignmentTestSuite) TestObjectiveBasedAlignment(t *testing.T) {
	t.Logf("🎯 Testing objective-based alignment measurement")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute objective-based alignment measurement
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--alignment-type", "objective-based", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute objective-based alignment measurement: %v", err)
	}

	// Parse alignment result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse objective-based alignment JSON: %v", err)
	}

	// Validate objective alignment rate meets target
	if result.ObjectiveAlignmentRate < suite.AlignmentTargets.OverallTarget {
		t.Errorf("Objective alignment rate %.1f%% below target of %.1f%%",
			result.ObjectiveAlignmentRate, suite.AlignmentTargets.OverallTarget)
	}

	// Validate measurement accuracy
	if result.AlignmentAccuracy < 90.0 {
		t.Errorf("Objective-based alignment accuracy %.1f%% below target of 90%%", result.AlignmentAccuracy)
	}

	t.Logf("✅ Objective-based alignment: %.1f%% (target: %.1f%%, accuracy: %.1f%%)",
		result.ObjectiveAlignmentRate, suite.AlignmentTargets.OverallTarget, result.AlignmentAccuracy)
}

// TestOutcomeBasedAlignment tests outcome-based alignment measurement
func (suite *GoalAlignmentTestSuite) TestOutcomeBasedAlignment(t *testing.T) {
	t.Logf("📊 Testing outcome-based alignment measurement")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute outcome-based alignment measurement
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--alignment-type", "outcome-based", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute outcome-based alignment measurement: %v", err)
	}

	// Parse alignment result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse outcome-based alignment JSON: %v", err)
	}

	// Validate outcome alignment rate is reasonable
	if result.OutcomeAlignmentRate < 80.0 || result.OutcomeAlignmentRate > 100.0 {
		t.Errorf("Outcome alignment rate %.1f%% outside reasonable range (80-100%%)", result.OutcomeAlignmentRate)
	}

	// Validate measurement reliability
	if result.MeasurementReliability < 85.0 {
		t.Errorf("Outcome-based alignment reliability %.1f%% below target of 85%%", result.MeasurementReliability)
	}

	t.Logf("✅ Outcome-based alignment: %.1f%% (reliability: %.1f%%)",
		result.OutcomeAlignmentRate, result.MeasurementReliability)
}

// TestImpactBasedAlignment tests impact-based alignment measurement
func (suite *GoalAlignmentTestSuite) TestImpactBasedAlignment(t *testing.T) {
	t.Logf("📈 Testing impact-based alignment measurement")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute impact-based alignment measurement
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--alignment-type", "impact-based", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute impact-based alignment measurement: %v", err)
	}

	// Parse alignment result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse impact-based alignment JSON: %v", err)
	}

	// Validate criteria alignment rate is reasonable
	if result.CriteriaAlignmentRate < 75.0 || result.CriteriaAlignmentRate > 100.0 {
		t.Errorf("Criteria alignment rate %.1f%% outside reasonable range (75-100%%)", result.CriteriaAlignmentRate)
	}

	// Validate measurement accuracy
	if result.AlignmentAccuracy < 85.0 {
		t.Errorf("Impact-based alignment accuracy %.1f%% below target of 85%%", result.AlignmentAccuracy)
	}

	t.Logf("✅ Impact-based alignment: %.1f%% (accuracy: %.1f%%)",
		result.CriteriaAlignmentRate, result.AlignmentAccuracy)
}

// TestAlignmentMeasurementAccuracy tests overall alignment measurement accuracy
func (suite *GoalAlignmentTestSuite) TestAlignmentMeasurementAccuracy(t *testing.T) {
	t.Logf("🔬 Testing alignment measurement accuracy")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute alignment measurement accuracy validation
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--validate-alignment-accuracy", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute alignment accuracy validation: %v", err)
	}

	// Parse accuracy result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse alignment accuracy JSON: %v", err)
	}

	// Validate measurement accuracy meets target (≥95%)
	if result.AlignmentAccuracy < 95.0 {
		t.Errorf("Alignment measurement accuracy %.1f%% below target of 95%%", result.AlignmentAccuracy)
	}

	// Validate measurement reliability is high
	if result.MeasurementReliability < 90.0 {
		t.Errorf("Alignment measurement reliability %.1f%% below target of 90%%", result.MeasurementReliability)
	}

	// Validate overall alignment meets target
	if result.OverallAlignmentRate < suite.AlignmentTargets.OverallTarget {
		t.Errorf("Overall alignment rate %.1f%% below target of %.1f%%",
			result.OverallAlignmentRate, suite.AlignmentTargets.OverallTarget)
	}

	t.Logf("✅ Alignment measurement accuracy validated: %.1f%% (target: ≥95%%, reliability: %.1f%%)",
		result.AlignmentAccuracy, result.MeasurementReliability)
}

// TestAlignmentTargetValidation validates that alignment targets are being met
func (suite *GoalAlignmentTestSuite) TestAlignmentTargetValidation(t *testing.T) {
	t.Logf("🎯 Testing alignment target validation")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute alignment target validation
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--validate-alignment-targets", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute alignment target validation: %v", err)
	}

	// Parse target validation result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse alignment target validation JSON: %v", err)
	}

	// Validate overall target is met
	if result.OverallAlignmentRate < suite.AlignmentTargets.OverallTarget {
		t.Errorf("Overall alignment target not met: %.1f%% < %.1f%%",
			result.OverallAlignmentRate, suite.AlignmentTargets.OverallTarget)
	}

	// Check component alignment rates against targets
	componentRates := map[string]float64{
		"objective": result.ObjectiveAlignmentRate,
		"outcome":   result.OutcomeAlignmentRate,
		"criteria":  result.CriteriaAlignmentRate,
	}

	for component, rate := range componentRates {
		if rate < 80.0 { // Minimum expected for components
			t.Errorf("%s alignment rate %.1f%% below minimum threshold of 80%%", component, rate)
		}
	}

	t.Logf("✅ Alignment targets validated: overall %.1f%% (target: %.1f%%)",
		result.OverallAlignmentRate, suite.AlignmentTargets.OverallTarget)
}

// TestAlignmentConsistency tests alignment measurement consistency across multiple runs
func (suite *GoalAlignmentTestSuite) TestAlignmentConsistency(t *testing.T) {
	t.Logf("🔄 Testing alignment measurement consistency")

	const numRuns = 3
	var alignmentRates []float64

	for i := 0; i < numRuns; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Execute alignment measurement
		cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--goal-alignment", "--json")
		cmd.Dir = suite.WorkspaceRoot

		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to execute alignment measurement (run %d): %v", i+1, err)
		}

		// Parse alignment result
		var result GoalAlignmentResult
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("Failed to parse alignment JSON (run %d): %v", i+1, err)
		}

		alignmentRates = append(alignmentRates, result.OverallAlignmentRate)

		// Small delay between runs
		time.Sleep(100 * time.Millisecond)
	}

	// Validate consistency across runs (within 2% variation)
	baseRate := alignmentRates[0]
	for i, rate := range alignmentRates[1:] {
		if math.Abs(rate-baseRate) > 2.0 {
			t.Errorf("Alignment rate inconsistent between runs: %.1f%% vs %.1f%% (run %d)", baseRate, rate, i+2)
		}
	}

	t.Logf("✅ Alignment measurement consistency validated across %d runs", numRuns)
}

// CreateAlignmentScenarioFile creates a test scenario file for goal alignment
func (suite *GoalAlignmentTestSuite) CreateAlignmentScenarioFile(filename string, scenario GoalAlignmentScenario) error {
	scenarioJSON, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal alignment scenario: %w", err)
	}

	return os.WriteFile(filename, scenarioJSON, 0644)
}

// ExecuteGoalAlignmentAnalysis executes goal alignment analysis for a scenario
func (suite *GoalAlignmentTestSuite) ExecuteGoalAlignmentAnalysis(scenarioFile string) (*GoalAlignmentResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute goal alignment analysis
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--analyze-goal-alignment", scenarioFile, "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute goal alignment analysis: %w", err)
	}

	// Parse analysis result
	var result GoalAlignmentResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse goal alignment analysis JSON: %w", err)
	}

	return &result, nil
}
