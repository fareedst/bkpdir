package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ComplianceRateTestSuite focuses on decision compliance rate accuracy
type ComplianceRateTestSuite struct {
	WorkspaceRoot        string
	ValidationScriptPath string
	MetricsScriptPath    string
	TestScenariosPath    string
}

// ComplianceScenario represents a test scenario for compliance validation
type ComplianceScenario struct {
	Name               string                `json:"name"`
	Description        string                `json:"description"`
	DecisionComponents []DecisionComponent   `json:"decision_components"`
	ExpectedCompliance float64               `json:"expected_compliance"`
	ValidationCriteria []ValidationCriterion `json:"validation_criteria"`
	ScenarioType       string                `json:"scenario_type"`
}

// DecisionComponent represents a decision framework component
type DecisionComponent struct {
	HierarchyLevel   string  `json:"hierarchy_level"` // Safety Gates, Scope Boundaries, Quality Thresholds, Goal Alignment
	ComponentName    string  `json:"component_name"`
	ComplianceStatus string  `json:"compliance_status"` // COMPLIANT, NON_COMPLIANT, PARTIAL
	ValidationResult bool    `json:"validation_result"`
	ComplianceScore  float64 `json:"compliance_score"`
}

// ValidationCriterion represents criteria for compliance validation
type ValidationCriterion struct {
	CriterionName    string  `json:"criterion_name"`
	RequiredValue    string  `json:"required_value"`
	ActualValue      string  `json:"actual_value"`
	ComplianceWeight float64 `json:"compliance_weight"`
	ValidationPassed bool    `json:"validation_passed"`
}

// ComplianceAnalysisResult represents compliance analysis outcome
type ComplianceAnalysisResult struct {
	ScenarioName               string  `json:"scenario_name"`
	OverallComplianceRate      float64 `json:"overall_compliance_rate"`
	SafetyGateCompliance       float64 `json:"safety_gate_compliance"`
	ScopeBoundaryCompliance    float64 `json:"scope_boundary_compliance"`
	QualityThresholdCompliance float64 `json:"quality_threshold_compliance"`
	GoalAlignmentCompliance    float64 `json:"goal_alignment_compliance"`
	ComplianceAccuracy         float64 `json:"compliance_accuracy"`
	ValidationSuccess          bool    `json:"validation_success"`
	AnalysisTimestamp          string  `json:"analysis_timestamp"`
}

// TestDOC014ComplianceRateAccuracy - main test for compliance rate accuracy
func TestDOC014ComplianceRateAccuracy(t *testing.T) {
	suite := &ComplianceRateTestSuite{}

	// Initialize test suite
	if err := suite.Initialize(); err != nil {
		t.Fatalf("Failed to initialize compliance rate test suite: %v", err)
	}

	t.Logf("🎯 Starting DOC-014 Compliance Rate Accuracy Test Suite")
	t.Logf("📊 Workspace: %s", suite.WorkspaceRoot)
	t.Logf("🔍 Validation Script: %s", suite.ValidationScriptPath)

	// Create test scenarios
	scenarios := suite.CreateComplianceScenarios()

	// Run compliance rate accuracy tests
	t.Run("FullyCompliantScenario", func(t *testing.T) {
		suite.TestFullyCompliantScenario(t, scenarios[0])
	})
	t.Run("PartiallyCompliantScenario", func(t *testing.T) {
		suite.TestPartiallyCompliantScenario(t, scenarios[1])
	})
	t.Run("NonCompliantScenario", func(t *testing.T) {
		suite.TestNonCompliantScenario(t, scenarios[2])
	})
	t.Run("MixedComplianceScenario", func(t *testing.T) {
		suite.TestMixedComplianceScenario(t, scenarios[3])
	})
	t.Run("EdgeCaseScenarios", func(t *testing.T) {
		suite.TestEdgeCaseScenarios(t, scenarios[4:])
	})
	t.Run("ComplianceCalculationAccuracy", suite.TestComplianceCalculationAccuracy)
	t.Run("HierarchyLevelCompliance", suite.TestHierarchyLevelCompliance)
	t.Run("ComplianceRateConsistency", suite.TestComplianceRateConsistency)

	t.Logf("✅ DOC-014 Compliance Rate Accuracy Test Suite Completed")
}

// Initialize sets up the compliance rate test environment
func (suite *ComplianceRateTestSuite) Initialize() error {
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
	suite.ValidationScriptPath = filepath.Join(workspaceRoot, "scripts", "validate-decision-framework.sh")
	suite.MetricsScriptPath = filepath.Join(workspaceRoot, "scripts", "track-decision-metrics.sh")
	suite.TestScenariosPath = filepath.Join(workspaceRoot, "test", "metrics", "compliance_scenarios")

	// Create test scenarios directory
	return os.MkdirAll(suite.TestScenariosPath, 0755)
}

// CreateComplianceScenarios creates test scenarios for compliance rate validation
func (suite *ComplianceRateTestSuite) CreateComplianceScenarios() []ComplianceScenario {
	scenarios := []ComplianceScenario{
		{
			Name:        "FullyCompliantNewFeature",
			Description: "New feature implementation with full compliance",
			DecisionComponents: []DecisionComponent{
				{HierarchyLevel: "Safety Gates", ComponentName: "No Breaking Changes", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Scope Boundaries", ComponentName: "Feature Scope Validation", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Quality Thresholds", ComponentName: "Test Coverage", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Goal Alignment", ComponentName: "Project Objectives", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
			},
			ExpectedCompliance: 100.0,
			ValidationCriteria: []ValidationCriterion{
				{CriterionName: "All Safety Gates", RequiredValue: "PASS", ActualValue: "PASS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Scope Boundaries", RequiredValue: "WITHIN_SCOPE", ActualValue: "WITHIN_SCOPE", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Quality Thresholds", RequiredValue: "MEETS_STANDARDS", ActualValue: "MEETS_STANDARDS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Goal Alignment", RequiredValue: "ALIGNED", ActualValue: "ALIGNED", ComplianceWeight: 1.0, ValidationPassed: true},
			},
			ScenarioType: "FULLY_COMPLIANT",
		},
		{
			Name:        "PartiallyCompliantModification",
			Description: "Feature modification with partial compliance",
			DecisionComponents: []DecisionComponent{
				{HierarchyLevel: "Safety Gates", ComponentName: "No Breaking Changes", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Scope Boundaries", ComponentName: "Feature Scope Validation", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 75.0},
				{HierarchyLevel: "Quality Thresholds", ComponentName: "Test Coverage", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Goal Alignment", ComponentName: "Project Objectives", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 80.0},
			},
			ExpectedCompliance: 88.75, // (100 + 75 + 100 + 80) / 4
			ValidationCriteria: []ValidationCriterion{
				{CriterionName: "All Safety Gates", RequiredValue: "PASS", ActualValue: "PASS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Scope Boundaries", RequiredValue: "WITHIN_SCOPE", ActualValue: "PARTIALLY_WITHIN_SCOPE", ComplianceWeight: 1.0, ValidationPassed: false},
				{CriterionName: "Quality Thresholds", RequiredValue: "MEETS_STANDARDS", ActualValue: "MEETS_STANDARDS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Goal Alignment", RequiredValue: "ALIGNED", ActualValue: "PARTIALLY_ALIGNED", ComplianceWeight: 1.0, ValidationPassed: false},
			},
			ScenarioType: "PARTIALLY_COMPLIANT",
		},
		{
			Name:        "NonCompliantBugFix",
			Description: "Bug fix with non-compliance issues",
			DecisionComponents: []DecisionComponent{
				{HierarchyLevel: "Safety Gates", ComponentName: "No Breaking Changes", ComplianceStatus: "NON_COMPLIANT", ValidationResult: false, ComplianceScore: 0.0},
				{HierarchyLevel: "Scope Boundaries", ComponentName: "Feature Scope Validation", ComplianceStatus: "NON_COMPLIANT", ValidationResult: false, ComplianceScore: 30.0},
				{HierarchyLevel: "Quality Thresholds", ComponentName: "Test Coverage", ComplianceStatus: "NON_COMPLIANT", ValidationResult: false, ComplianceScore: 45.0},
				{HierarchyLevel: "Goal Alignment", ComponentName: "Project Objectives", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 60.0},
			},
			ExpectedCompliance: 33.75, // (0 + 30 + 45 + 60) / 4
			ValidationCriteria: []ValidationCriterion{
				{CriterionName: "All Safety Gates", RequiredValue: "PASS", ActualValue: "FAIL", ComplianceWeight: 1.0, ValidationPassed: false},
				{CriterionName: "Scope Boundaries", RequiredValue: "WITHIN_SCOPE", ActualValue: "OUT_OF_SCOPE", ComplianceWeight: 1.0, ValidationPassed: false},
				{CriterionName: "Quality Thresholds", RequiredValue: "MEETS_STANDARDS", ActualValue: "BELOW_STANDARDS", ComplianceWeight: 1.0, ValidationPassed: false},
				{CriterionName: "Goal Alignment", RequiredValue: "ALIGNED", ActualValue: "PARTIALLY_ALIGNED", ComplianceWeight: 1.0, ValidationPassed: false},
			},
			ScenarioType: "NON_COMPLIANT",
		},
		{
			Name:        "MixedComplianceRefactoring",
			Description: "Refactoring with mixed compliance levels",
			DecisionComponents: []DecisionComponent{
				{HierarchyLevel: "Safety Gates", ComponentName: "No Breaking Changes", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 100.0},
				{HierarchyLevel: "Scope Boundaries", ComponentName: "Feature Scope Validation", ComplianceStatus: "COMPLIANT", ValidationResult: true, ComplianceScore: 95.0},
				{HierarchyLevel: "Quality Thresholds", ComponentName: "Test Coverage", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 70.0},
				{HierarchyLevel: "Goal Alignment", ComponentName: "Project Objectives", ComplianceStatus: "NON_COMPLIANT", ValidationResult: false, ComplianceScore: 40.0},
			},
			ExpectedCompliance: 76.25, // (100 + 95 + 70 + 40) / 4
			ValidationCriteria: []ValidationCriterion{
				{CriterionName: "All Safety Gates", RequiredValue: "PASS", ActualValue: "PASS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Scope Boundaries", RequiredValue: "WITHIN_SCOPE", ActualValue: "WITHIN_SCOPE", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Quality Thresholds", RequiredValue: "MEETS_STANDARDS", ActualValue: "PARTIALLY_MEETS_STANDARDS", ComplianceWeight: 1.0, ValidationPassed: false},
				{CriterionName: "Goal Alignment", RequiredValue: "ALIGNED", ActualValue: "NOT_ALIGNED", ComplianceWeight: 1.0, ValidationPassed: false},
			},
			ScenarioType: "MIXED_COMPLIANCE",
		},
		{
			Name:        "EdgeCaseMinimalCompliance",
			Description: "Edge case with minimal compliance",
			DecisionComponents: []DecisionComponent{
				{HierarchyLevel: "Safety Gates", ComponentName: "No Breaking Changes", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 51.0},
				{HierarchyLevel: "Scope Boundaries", ComponentName: "Feature Scope Validation", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 50.0},
				{HierarchyLevel: "Quality Thresholds", ComponentName: "Test Coverage", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 49.0},
				{HierarchyLevel: "Goal Alignment", ComponentName: "Project Objectives", ComplianceStatus: "PARTIAL", ValidationResult: true, ComplianceScore: 50.0},
			},
			ExpectedCompliance: 50.0, // (51 + 50 + 49 + 50) / 4
			ValidationCriteria: []ValidationCriterion{
				{CriterionName: "All Safety Gates", RequiredValue: "PASS", ActualValue: "MARGINAL_PASS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Scope Boundaries", RequiredValue: "WITHIN_SCOPE", ActualValue: "MARGINAL_SCOPE", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Quality Thresholds", RequiredValue: "MEETS_STANDARDS", ActualValue: "MARGINAL_STANDARDS", ComplianceWeight: 1.0, ValidationPassed: true},
				{CriterionName: "Goal Alignment", RequiredValue: "ALIGNED", ActualValue: "MARGINAL_ALIGNMENT", ComplianceWeight: 1.0, ValidationPassed: true},
			},
			ScenarioType: "EDGE_CASE",
		},
	}

	return scenarios
}

// TestFullyCompliantScenario tests compliance rate calculation for fully compliant scenarios
func (suite *ComplianceRateTestSuite) TestFullyCompliantScenario(t *testing.T, scenario ComplianceScenario) {
	t.Logf("✅ Testing fully compliant scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create scenario file: %v", err)
	}

	// Execute compliance rate calculation
	result, err := suite.ExecuteComplianceAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute compliance analysis: %v", err)
	}

	// Validate compliance rate accuracy
	tolerance := 2.0 // 2% tolerance for fully compliant scenarios
	if result.ComplianceAccuracy < (100.0 - tolerance) {
		t.Errorf("Compliance accuracy %.1f%% below expected (≥%.1f%%)", result.ComplianceAccuracy, 100.0-tolerance)
	}

	// Validate expected vs actual compliance rate
	actualCompliance := result.OverallComplianceRate
	expectedCompliance := scenario.ExpectedCompliance

	if actualCompliance < expectedCompliance-tolerance || actualCompliance > expectedCompliance+tolerance {
		t.Errorf("Compliance rate mismatch: expected %.1f%%, got %.1f%%", expectedCompliance, actualCompliance)
	}

	// For fully compliant scenarios, all hierarchy levels should be high
	if result.SafetyGateCompliance < 95.0 || result.ScopeBoundaryCompliance < 95.0 ||
		result.QualityThresholdCompliance < 95.0 || result.GoalAlignmentCompliance < 95.0 {
		t.Errorf("Hierarchy level compliance below expected for fully compliant scenario")
	}

	t.Logf("✅ Fully compliant scenario validation passed: %.1f%% compliance (accuracy: %.1f%%)",
		actualCompliance, result.ComplianceAccuracy)
}

// TestPartiallyCompliantScenario tests compliance rate calculation for partially compliant scenarios
func (suite *ComplianceRateTestSuite) TestPartiallyCompliantScenario(t *testing.T, scenario ComplianceScenario) {
	t.Logf("🔶 Testing partially compliant scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create scenario file: %v", err)
	}

	// Execute compliance rate calculation
	result, err := suite.ExecuteComplianceAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute compliance analysis: %v", err)
	}

	// Validate compliance rate accuracy
	tolerance := 3.0 // 3% tolerance for partially compliant scenarios
	if result.ComplianceAccuracy < (100.0 - tolerance) {
		t.Errorf("Compliance accuracy %.1f%% below expected (≥%.1f%%)", result.ComplianceAccuracy, 100.0-tolerance)
	}

	// Validate expected vs actual compliance rate
	actualCompliance := result.OverallComplianceRate
	expectedCompliance := scenario.ExpectedCompliance

	if actualCompliance < expectedCompliance-tolerance || actualCompliance > expectedCompliance+tolerance {
		t.Errorf("Compliance rate mismatch: expected %.1f%%, got %.1f%%", expectedCompliance, actualCompliance)
	}

	// For partially compliant scenarios, compliance should be moderate (50-95%)
	if result.OverallComplianceRate < 50.0 || result.OverallComplianceRate > 95.0 {
		t.Errorf("Partially compliant scenario compliance %.1f%% outside expected range (50-95%%)", result.OverallComplianceRate)
	}

	t.Logf("✅ Partially compliant scenario validation passed: %.1f%% compliance (accuracy: %.1f%%)",
		actualCompliance, result.ComplianceAccuracy)
}

// TestNonCompliantScenario tests compliance rate calculation for non-compliant scenarios
func (suite *ComplianceRateTestSuite) TestNonCompliantScenario(t *testing.T, scenario ComplianceScenario) {
	t.Logf("❌ Testing non-compliant scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create scenario file: %v", err)
	}

	// Execute compliance rate calculation
	result, err := suite.ExecuteComplianceAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute compliance analysis: %v", err)
	}

	// Validate compliance rate accuracy
	tolerance := 5.0 // 5% tolerance for non-compliant scenarios
	if result.ComplianceAccuracy < (100.0 - tolerance) {
		t.Errorf("Compliance accuracy %.1f%% below expected (≥%.1f%%)", result.ComplianceAccuracy, 100.0-tolerance)
	}

	// Validate expected vs actual compliance rate
	actualCompliance := result.OverallComplianceRate
	expectedCompliance := scenario.ExpectedCompliance

	if actualCompliance < expectedCompliance-tolerance || actualCompliance > expectedCompliance+tolerance {
		t.Errorf("Compliance rate mismatch: expected %.1f%%, got %.1f%%", expectedCompliance, actualCompliance)
	}

	// For non-compliant scenarios, compliance should be low (<50%)
	if result.OverallComplianceRate >= 50.0 {
		t.Errorf("Non-compliant scenario compliance %.1f%% above expected threshold (<50%%)", result.OverallComplianceRate)
	}

	t.Logf("✅ Non-compliant scenario validation passed: %.1f%% compliance (accuracy: %.1f%%)",
		actualCompliance, result.ComplianceAccuracy)
}

// TestMixedComplianceScenario tests compliance rate calculation for mixed compliance scenarios
func (suite *ComplianceRateTestSuite) TestMixedComplianceScenario(t *testing.T, scenario ComplianceScenario) {
	t.Logf("🔀 Testing mixed compliance scenario: %s", scenario.Name)

	// Create scenario test file
	scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
	if err := suite.CreateScenarioFile(scenarioFile, scenario); err != nil {
		t.Fatalf("Failed to create scenario file: %v", err)
	}

	// Execute compliance rate calculation
	result, err := suite.ExecuteComplianceAnalysis(scenarioFile)
	if err != nil {
		t.Fatalf("Failed to execute compliance analysis: %v", err)
	}

	// Validate compliance rate accuracy
	tolerance := 4.0 // 4% tolerance for mixed compliance scenarios
	if result.ComplianceAccuracy < (100.0 - tolerance) {
		t.Errorf("Compliance accuracy %.1f%% below expected (≥%.1f%%)", result.ComplianceAccuracy, 100.0-tolerance)
	}

	// Validate expected vs actual compliance rate
	actualCompliance := result.OverallComplianceRate
	expectedCompliance := scenario.ExpectedCompliance

	if actualCompliance < expectedCompliance-tolerance || actualCompliance > expectedCompliance+tolerance {
		t.Errorf("Compliance rate mismatch: expected %.1f%%, got %.1f%%", expectedCompliance, actualCompliance)
	}

	// For mixed compliance scenarios, validate hierarchy level variation
	hierarchyLevels := []float64{result.SafetyGateCompliance, result.ScopeBoundaryCompliance, result.QualityThresholdCompliance, result.GoalAlignmentCompliance}

	// Calculate variance to ensure there is actual variation
	var sum, variance float64
	for _, level := range hierarchyLevels {
		sum += level
	}
	mean := sum / float64(len(hierarchyLevels))

	for _, level := range hierarchyLevels {
		variance += (level - mean) * (level - mean)
	}
	variance /= float64(len(hierarchyLevels))

	if variance < 100.0 { // Expect significant variation in mixed compliance
		t.Errorf("Mixed compliance scenario shows insufficient variation (variance: %.1f)", variance)
	}

	t.Logf("✅ Mixed compliance scenario validation passed: %.1f%% compliance (accuracy: %.1f%%, variance: %.1f)",
		actualCompliance, result.ComplianceAccuracy, variance)
}

// TestEdgeCaseScenarios tests compliance rate calculation for edge cases
func (suite *ComplianceRateTestSuite) TestEdgeCaseScenarios(t *testing.T, scenarios []ComplianceScenario) {
	t.Logf("🔍 Testing edge case scenarios")

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			// Create scenario test file
			scenarioFile := filepath.Join(suite.TestScenariosPath, fmt.Sprintf("%s.json", scenario.Name))
			if err := suite.CreateScenarioFile(scenarioFile, scenario); err != nil {
				t.Fatalf("Failed to create scenario file: %v", err)
			}

			// Execute compliance rate calculation
			result, err := suite.ExecuteComplianceAnalysis(scenarioFile)
			if err != nil {
				t.Fatalf("Failed to execute compliance analysis: %v", err)
			}

			// Validate compliance rate accuracy
			tolerance := 5.0 // 5% tolerance for edge cases
			if result.ComplianceAccuracy < (100.0 - tolerance) {
				t.Errorf("Edge case compliance accuracy %.1f%% below expected (≥%.1f%%)", result.ComplianceAccuracy, 100.0-tolerance)
			}

			// Validate expected vs actual compliance rate
			actualCompliance := result.OverallComplianceRate
			expectedCompliance := scenario.ExpectedCompliance

			if actualCompliance < expectedCompliance-tolerance || actualCompliance > expectedCompliance+tolerance {
				t.Errorf("Edge case compliance rate mismatch: expected %.1f%%, got %.1f%%", expectedCompliance, actualCompliance)
			}

			t.Logf("✅ Edge case scenario %s validation passed: %.1f%% compliance (accuracy: %.1f%%)",
				scenario.Name, actualCompliance, result.ComplianceAccuracy)
		})
	}
}

// TestComplianceCalculationAccuracy tests the accuracy of compliance calculation algorithms
func (suite *ComplianceRateTestSuite) TestComplianceCalculationAccuracy(t *testing.T) {
	t.Logf("🧮 Testing compliance calculation accuracy")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute compliance calculation validation
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--validate-calculations", "--compliance-rate", "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute compliance calculation validation: %v", err)
	}

	// Parse validation output
	var result ComplianceAnalysisResult
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatalf("Failed to parse compliance calculation JSON: %v", err)
	}

	// Validate calculation accuracy meets target (≥95%)
	if result.ComplianceAccuracy < 95.0 {
		t.Errorf("Compliance calculation accuracy %.1f%% below target of 95%%", result.ComplianceAccuracy)
	}

	// Validate all hierarchy levels are calculated
	if result.SafetyGateCompliance < 0 || result.ScopeBoundaryCompliance < 0 ||
		result.QualityThresholdCompliance < 0 || result.GoalAlignmentCompliance < 0 {
		t.Errorf("Negative compliance values detected in hierarchy levels")
	}

	// Validate overall compliance is calculated correctly
	expectedOverall := (result.SafetyGateCompliance + result.ScopeBoundaryCompliance +
		result.QualityThresholdCompliance + result.GoalAlignmentCompliance) / 4.0

	if result.OverallComplianceRate < expectedOverall-1.0 || result.OverallComplianceRate > expectedOverall+1.0 {
		t.Errorf("Overall compliance calculation error: expected %.1f%%, got %.1f%%", expectedOverall, result.OverallComplianceRate)
	}

	t.Logf("✅ Compliance calculation accuracy validated: %.1f%% (target: ≥95%%)", result.ComplianceAccuracy)
}

// TestHierarchyLevelCompliance tests compliance calculation for each hierarchy level
func (suite *ComplianceRateTestSuite) TestHierarchyLevelCompliance(t *testing.T) {
	t.Logf("🏗️ Testing hierarchy level compliance calculation")

	hierarchyLevels := []string{"safety-gates", "scope-boundaries", "quality-thresholds", "goal-alignment"}

	for _, level := range hierarchyLevels {
		t.Run(level, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// Execute hierarchy level compliance calculation
			cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--hierarchy-level", level, "--compliance-rate", "--json")
			cmd.Dir = suite.WorkspaceRoot

			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Failed to execute hierarchy level compliance calculation: %v", err)
			}

			// Parse compliance output
			var result ComplianceAnalysisResult
			if err := json.Unmarshal(output, &result); err != nil {
				t.Fatalf("Failed to parse hierarchy level compliance JSON: %v", err)
			}

			// Validate compliance rate is reasonable (0-100%)
			var levelCompliance float64
			switch level {
			case "safety-gates":
				levelCompliance = result.SafetyGateCompliance
			case "scope-boundaries":
				levelCompliance = result.ScopeBoundaryCompliance
			case "quality-thresholds":
				levelCompliance = result.QualityThresholdCompliance
			case "goal-alignment":
				levelCompliance = result.GoalAlignmentCompliance
			}

			if levelCompliance < 0 || levelCompliance > 100 {
				t.Errorf("Hierarchy level %s compliance %.1f%% outside valid range (0-100%%)", level, levelCompliance)
			}

			// Validate calculation accuracy
			if result.ComplianceAccuracy < 90.0 {
				t.Errorf("Hierarchy level %s compliance accuracy %.1f%% below target of 90%%", level, result.ComplianceAccuracy)
			}

			t.Logf("✅ Hierarchy level %s compliance: %.1f%% (accuracy: %.1f%%)", level, levelCompliance, result.ComplianceAccuracy)
		})
	}
}

// TestComplianceRateConsistency tests compliance rate consistency across multiple runs
func (suite *ComplianceRateTestSuite) TestComplianceRateConsistency(t *testing.T) {
	t.Logf("🔄 Testing compliance rate consistency across multiple runs")

	const numRuns = 3
	var complianceRates []float64

	for i := 0; i < numRuns; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Execute compliance rate calculation
		cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--compliance-rate", "--json")
		cmd.Dir = suite.WorkspaceRoot

		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Failed to execute compliance rate calculation (run %d): %v", i+1, err)
		}

		// Parse compliance output
		var result ComplianceAnalysisResult
		if err := json.Unmarshal(output, &result); err != nil {
			t.Fatalf("Failed to parse compliance rate JSON (run %d): %v", i+1, err)
		}

		complianceRates = append(complianceRates, result.OverallComplianceRate)

		// Small delay between runs
		time.Sleep(100 * time.Millisecond)
	}

	// Validate consistency across runs (within 1% variation)
	baseRate := complianceRates[0]
	for i, rate := range complianceRates[1:] {
		if rate < baseRate-1.0 || rate > baseRate+1.0 {
			t.Errorf("Compliance rate inconsistent between runs: %.1f%% vs %.1f%% (run %d)", baseRate, rate, i+2)
		}
	}

	t.Logf("✅ Compliance rate consistency validated across %d runs", numRuns)
}

// CreateScenarioFile creates a test scenario file
func (suite *ComplianceRateTestSuite) CreateScenarioFile(filename string, scenario ComplianceScenario) error {
	scenarioJSON, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scenario: %w", err)
	}

	return os.WriteFile(filename, scenarioJSON, 0644)
}

// ExecuteComplianceAnalysis executes compliance analysis for a scenario
func (suite *ComplianceRateTestSuite) ExecuteComplianceAnalysis(scenarioFile string) (*ComplianceAnalysisResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Execute compliance analysis
	cmd := exec.CommandContext(ctx, "bash", suite.MetricsScriptPath, "--analyze-scenario", scenarioFile, "--json")
	cmd.Dir = suite.WorkspaceRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute compliance analysis: %w", err)
	}

	// Parse analysis result
	var result ComplianceAnalysisResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse compliance analysis JSON: %w", err)
	}

	return &result, nil
}
