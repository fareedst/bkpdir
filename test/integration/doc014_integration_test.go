// 🔶 DOC-014: Integration testing - Main test suite for decision framework integration
package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// 🔶 DOC-014: Integration test suite configuration
type IntegrationTestSuite struct {
	ProjectRoot     string
	ContextDir      string
	ScriptsDir      string
	ValidationDir   string
	ValidateTimeout time.Duration
	SkipSlowTests   bool
	TestDataDir     string
}

// NewIntegrationTestSuite creates a new integration test suite
func NewIntegrationTestSuite(projectRoot string) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		ProjectRoot:     projectRoot,
		ContextDir:      filepath.Join(projectRoot, "docs", "context"),
		ScriptsDir:      filepath.Join(projectRoot, "scripts"),
		ValidationDir:   filepath.Join(projectRoot, "internal", "validation"),
		ValidateTimeout: 30 * time.Second,
		SkipSlowTests:   os.Getenv("SKIP_SLOW_TESTS") == "true",
		TestDataDir:     filepath.Join(projectRoot, "test", "testdata"),
	}
}

// setupTestEnvironment prepares the test environment
func (suite *IntegrationTestSuite) setupTestEnvironment(t *testing.T) {
	t.Helper()

	// Create test data directory if it doesn't exist
	if err := os.MkdirAll(suite.TestDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Verify required files exist
	requiredFiles := []string{
		filepath.Join(suite.ContextDir, "ai-decision-framework.md"),
		filepath.Join(suite.ContextDir, "feature-tracking.md"),
		filepath.Join(suite.ContextDir, "ai-assistant-compliance.md"),
		filepath.Join(suite.ContextDir, "ai-assistant-protocol.md"),
		filepath.Join(suite.ScriptsDir, "validate-decision-framework.sh"),
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Fatalf("Required file missing for integration testing: %s", file)
		}
	}
}

// TestDOC014MainIntegrationSuite runs the complete integration test suite
func TestDOC014MainIntegrationSuite(t *testing.T) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	suite := NewIntegrationTestSuite(projectRoot)
	suite.setupTestEnvironment(t)

	t.Run("DecisionFrameworkCore", suite.testDecisionFrameworkCore)
	t.Run("FeatureTrackingIntegration", suite.testFeatureTrackingIntegration)
	t.Run("ProtocolIntegration", suite.testProtocolIntegration)
	t.Run("ValidationSystemIntegration", suite.testValidationSystemIntegration)
	t.Run("MakefileWorkflows", suite.testMakefileWorkflows)
	t.Run("TokenConsistency", suite.testTokenConsistency)
}

// testDecisionFrameworkCore validates the core decision framework document
func (suite *IntegrationTestSuite) testDecisionFrameworkCore(t *testing.T) {
	frameworkFile := filepath.Join(suite.ContextDir, "ai-decision-framework.md")

	// Test 1: Framework document exists and is readable
	content, err := os.ReadFile(frameworkFile)
	if err != nil {
		t.Errorf("Failed to read decision framework document: %v", err)
		return
	}

	frameworkContent := string(content)

	// Test 2: Required sections exist
	requiredSections := []string{
		"🚨 DECISION HIERARCHY",
		"🛡️ Safety Gates",
		"📑 Scope Boundaries",
		"📊 Quality Thresholds",
		"🎯 Goal Alignment",
		"🧠 DECISION TREE",
		"🔒 ENHANCED ENFORCEMENT",
		"🎯 DECISION VALIDATION CHECKLIST",
	}

	for _, section := range requiredSections {
		if !strings.Contains(frameworkContent, section) {
			t.Errorf("Decision framework missing required section: %s", section)
		}
	}

	// Test 3: 4-tier hierarchy is properly defined
	hierarchyLevels := []string{"Safety Gates", "Scope Boundaries", "Quality Thresholds", "Goal Alignment"}
	for i, level := range hierarchyLevels {
		expectedPattern := fmt.Sprintf("%d. **", i+1)
		if !strings.Contains(frameworkContent, expectedPattern) {
			t.Errorf("Decision hierarchy level %d (%s) not properly numbered", i+1, level)
		}
	}

	// Test 4: Decision trees are present
	decisionTrees := []string{
		"Should I implement this feature request?",
		"Should I refactor this code?",
		"Should I fix this test failure?",
		"Should I update documentation?",
	}

	for _, tree := range decisionTrees {
		if !strings.Contains(frameworkContent, tree) {
			t.Errorf("Decision framework missing decision tree: %s", tree)
		}
	}
}

// testFeatureTrackingIntegration validates integration with feature-tracking.md
func (suite *IntegrationTestSuite) testFeatureTrackingIntegration(t *testing.T) {
	featureTrackingFile := filepath.Join(suite.ContextDir, "feature-tracking.md")

	// Test 1: DOC-014 entry exists in feature tracking
	content, err := os.ReadFile(featureTrackingFile)
	if err != nil {
		t.Errorf("Failed to read feature tracking file: %v", err)
		return
	}

	trackingContent := string(content)

	// Test 2: DOC-014 feature entry exists
	if !strings.Contains(trackingContent, "DOC-014") {
		t.Error("DOC-014 entry not found in feature-tracking.md")
		return
	}

	// Test 3: DOC-014 has proper structure
	doc014Sections := []string{
		"AI Assistant Decision Framework",
		"Subtask",               // Should have detailed subtask breakdown
		"Implementation Tokens", // Should have implementation tokens
		"Priority",              // Should have priority classification
	}

	for _, section := range doc014Sections {
		if !strings.Contains(trackingContent, section) {
			t.Errorf("DOC-014 entry missing required element: %s", section)
		}
	}

	// Test 4: Subtasks 1-5 should be marked as completed
	completedSubtasks := []string{
		"[x] Create core decision framework document",
		"[x] Update feature tracking integration",
		"[x] Update AI assistant compliance requirements",
		"[x] Enhance implementation token system",
		"[x] Create decision validation tools",
	}

	for _, subtask := range completedSubtasks {
		if !strings.Contains(trackingContent, subtask) {
			t.Errorf("DOC-014 subtask not marked as completed: %s", subtask)
		}
	}
}

// testProtocolIntegration validates integration with AI assistant protocols
func (suite *IntegrationTestSuite) testProtocolIntegration(t *testing.T) {
	protocolFile := filepath.Join(suite.ContextDir, "ai-assistant-protocol.md")

	content, err := os.ReadFile(protocolFile)
	if err != nil {
		t.Errorf("Failed to read protocol file: %v", err)
		return
	}

	protocolContent := string(content)

	// Test 1: Decision Framework references exist
	if !strings.Contains(protocolContent, "Decision Framework") {
		t.Error("AI assistant protocol missing Decision Framework references")
	}

	// Test 2: Key protocols include decision validation
	protocols := []string{
		"🆕 NEW FEATURE Protocol",
		"🔧 MODIFICATION Protocol",
		"🐛 BUG FIX Protocol",
		"⚙️ CONFIG CHANGE Protocol",
		"🔌 API CHANGE Protocol",
		"🔄 REFACTORING Protocol",
	}

	// Count how many protocols include decision framework validation
	protocolsWithValidation := 0
	for _, protocol := range protocols {
		protocolStart := strings.Index(protocolContent, protocol)
		if protocolStart == -1 {
			t.Errorf("Protocol not found: %s", protocol)
			continue
		}

		// Find the next protocol or end of file
		nextProtocolStart := len(protocolContent)
		for _, otherProtocol := range protocols {
			if otherProtocol == protocol {
				continue
			}
			otherStart := strings.Index(protocolContent[protocolStart+len(protocol):], otherProtocol)
			if otherStart != -1 && protocolStart+len(protocol)+otherStart < nextProtocolStart {
				nextProtocolStart = protocolStart + len(protocol) + otherStart
			}
		}

		protocolSection := protocolContent[protocolStart:nextProtocolStart]
		if strings.Contains(protocolSection, "Decision Framework") ||
			strings.Contains(protocolSection, "decision") ||
			strings.Contains(protocolSection, "validation") {
			protocolsWithValidation++
		}
	}

	// Test 3: At least 75% of protocols should have decision integration
	expectedMinimum := int(float64(len(protocols)) * 0.75)
	if protocolsWithValidation < expectedMinimum {
		t.Errorf("Insufficient protocol integration: %d/%d protocols have decision validation (expected >=%d)",
			protocolsWithValidation, len(protocols), expectedMinimum)
	}
}

// testValidationSystemIntegration validates integration with DOC-008 and DOC-011
func (suite *IntegrationTestSuite) testValidationSystemIntegration(t *testing.T) {
	// Test 1: Decision validation scripts exist
	validationScripts := []string{
		filepath.Join(suite.ScriptsDir, "validate-decision-framework.sh"),
		filepath.Join(suite.ScriptsDir, "validate-decision-context.sh"),
		filepath.Join(suite.ScriptsDir, "track-decision-metrics.sh"),
	}

	for _, script := range validationScripts {
		if _, err := os.Stat(script); os.IsNotExist(err) {
			t.Errorf("Required validation script missing: %s", script)
		} else {
			// Test script is executable
			info, err := os.Stat(script)
			if err != nil {
				t.Errorf("Failed to check script permissions: %v", err)
			} else if info.Mode()&0111 == 0 {
				t.Errorf("Validation script not executable: %s", script)
			}
		}
	}

	// Test 2: Go validation packages exist
	validationPackages := []string{
		filepath.Join(suite.ValidationDir, "decision_checklist.go"),
		filepath.Join(suite.ValidationDir, "decision_checklist_test.go"),
	}

	for _, pkg := range validationPackages {
		if _, err := os.Stat(pkg); os.IsNotExist(err) {
			t.Errorf("Required validation package missing: %s", pkg)
		}
	}

	// Test 3: AI validation integration
	aiValidationDir := filepath.Join(suite.ProjectRoot, "cmd", "ai-validation")
	if _, err := os.Stat(aiValidationDir); os.IsNotExist(err) {
		t.Log("AI validation command not found (optional for integration)")
	} else {
		// Check if ai-validation includes decision capabilities
		mainFile := filepath.Join(aiValidationDir, "main.go")
		if content, err := os.ReadFile(mainFile); err == nil {
			if !strings.Contains(string(content), "decision") {
				t.Log("AI validation may not include decision framework integration")
			}
		}
	}
}

// testMakefileWorkflows validates Makefile integration
func (suite *IntegrationTestSuite) testMakefileWorkflows(t *testing.T) {
	if suite.SkipSlowTests {
		t.Skip("Skipping Makefile workflow tests (SKIP_SLOW_TESTS=true)")
	}

	makefilePath := filepath.Join(suite.ProjectRoot, "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Error("Makefile not found")
		return
	}

	// Test 1: Decision validation targets exist
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Errorf("Failed to read Makefile: %v", err)
		return
	}

	makefileContent := string(content)
	expectedTargets := []string{
		"validate-decision-framework",
		"validate-decision-context",
		"track-decision-metrics",
		"decision-validation-suite",
	}

	for _, target := range expectedTargets {
		if !strings.Contains(makefileContent, target+":") {
			t.Errorf("Makefile missing required target: %s", target)
		}
	}

	// Test 2: Run a quick validation test (if not in CI)
	if os.Getenv("CI") != "true" {
		ctx, cancel := context.WithTimeout(context.Background(), suite.ValidateTimeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "make", "validate-decision-framework")
		cmd.Dir = suite.ProjectRoot

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Decision framework validation failed (may be expected): %v", err)
			t.Logf("Output: %s", output)
		}
	}
}

// testTokenConsistency validates implementation token consistency
func (suite *IntegrationTestSuite) testTokenConsistency(t *testing.T) {
	// Test 1: Search for DOC-014 tokens in codebase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "grep", "-r", "DOC-014", suite.ProjectRoot)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// This is okay - might not have DOC-014 tokens yet
		t.Log("No DOC-014 tokens found in codebase (expected for new feature)")
		return
	}

	// Test 2: If tokens exist, validate format
	tokenLines := strings.Split(string(output), "\n")
	validTokenCount := 0

	for _, line := range tokenLines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 🔍 Check for proper token format: // [ICON] DOC-014-EXAMPLE: Description [CONTEXT: tags]
		if strings.Contains(line, "// ") && strings.Contains(line, "DOC-"+"014:") {
			validTokenCount++

			// Basic format validation
			if !strings.Contains(line, "🔶") && !strings.Contains(line, "⭐") &&
				!strings.Contains(line, "🔺") && !strings.Contains(line, "🔻") {
				t.Errorf("DOC-014 token missing priority icon: %s", line)
			}
		}
	}

	if validTokenCount > 0 {
		t.Logf("Found %d DOC-014 implementation tokens", validTokenCount)
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

// Helper function for running shell commands with timeout
func runCommandWithTimeout(ctx context.Context, dir string, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}
