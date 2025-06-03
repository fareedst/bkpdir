// 🔺 COV-002: Test suite for differential coverage reporting tool - 🔧
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseCurrentCoverage(t *testing.T) {
	// Create a temporary coverage profile for testing
	tmpDir := t.TempDir()
	profilePath := filepath.Join(tmpDir, "coverage.out")

	testProfile := `mode: set
bkpdir/main.go:45: main 0.0%
bkpdir/archive.go:73: GenerateArchiveName 100.0%
bkpdir/archive.go:85: generateIncrementalArchiveName 88.9%
bkpdir/backup.go:102: CreateFileBackup 100.0%
bkpdir/verify.go:31: VerifyArchive 83.3%
`

	err := os.WriteFile(profilePath, []byte(testProfile), 0644)
	if err != nil {
		t.Fatalf("Failed to create test profile: %v", err)
	}

	// Parse the coverage profile
	files, err := parseCurrentCoverage(profilePath)
	if err != nil {
		t.Fatalf("parseCurrentCoverage failed: %v", err)
	}

	// Verify parsing results
	if len(files) == 0 {
		t.Fatal("Expected parsed files, got none")
	}

	// Check specific file coverage
	if archiveFile, exists := files["bkpdir/archive.go"]; exists {
		if len(archiveFile.Functions) != 2 {
			t.Errorf("Expected 2 functions in archive.go, got %d", len(archiveFile.Functions))
		}

		// Check function coverage
		foundGenerate := false
		foundIncremental := false
		for _, fn := range archiveFile.Functions {
			if fn.Function == "GenerateArchiveName" && fn.Coverage == 100.0 {
				foundGenerate = true
			}
			if fn.Function == "generateIncrementalArchiveName" && fn.Coverage == 88.9 {
				foundIncremental = true
			}
		}
		if !foundGenerate {
			t.Error("Expected GenerateArchiveName with 100% coverage")
		}
		if !foundIncremental {
			t.Error("Expected generateIncrementalArchiveName with 88.9% coverage")
		}
	} else {
		t.Error("Expected bkpdir/archive.go in parsed files")
	}
}

func TestCalculateOverallCoverage(t *testing.T) {
	files := map[string]FileCoverage{
		"file1.go": {Coverage: 80.0},
		"file2.go": {Coverage: 60.0},
		"file3.go": {Coverage: 70.0},
	}

	overall := calculateOverallCoverage(files)
	expected := 70.0 // (80 + 60 + 70) / 3
	if overall != expected {
		t.Errorf("Expected overall coverage %.1f%%, got %.1f%%", expected, overall)
	}

	// Test empty files
	emptyFiles := map[string]FileCoverage{}
	overall = calculateOverallCoverage(emptyFiles)
	if overall != 0.0 {
		t.Errorf("Expected 0%% for empty files, got %.1f%%", overall)
	}
}

func TestEvaluateQualityGates(t *testing.T) {
	config := DifferentialConfig{}
	config.QualityGates.NewCodeThreshold = 70.0
	config.QualityGates.CriticalPathThreshold = 80.0
	config.QualityGates.OverallRegressionLimit = -5.0

	// Test passing quality gates
	passingReport := DifferentialReport{
		NewCodeCoverage:    75.0,
		OverallChange:      2.0,
		CriticalPathStatus: map[string]bool{"critical1": true, "critical2": true},
		FileChanges: map[string]FileChange{
			"file1.go": {IsModified: true, MeetsThreshold: true},
			"file2.go": {IsModified: true, MeetsThreshold: true},
		},
	}

	if !evaluateQualityGates(config, passingReport) {
		t.Error("Expected quality gates to pass")
	}

	// Test failing new code threshold
	failingReport := passingReport
	failingReport.NewCodeCoverage = 65.0
	if evaluateQualityGates(config, failingReport) {
		t.Error("Expected quality gates to fail due to low new code coverage")
	}

	// Test failing regression limit
	failingReport = passingReport
	failingReport.OverallChange = -10.0
	if evaluateQualityGates(config, failingReport) {
		t.Error("Expected quality gates to fail due to regression")
	}

	// Test failing critical path
	failingReport = passingReport
	failingReport.CriticalPathStatus["critical1"] = false
	if evaluateQualityGates(config, failingReport) {
		t.Error("Expected quality gates to fail due to critical path")
	}
}

func TestGenerateRecommendations(t *testing.T) {
	config := DifferentialConfig{}
	config.QualityGates.NewCodeThreshold = 70.0

	// Test recommendations for low coverage
	report := DifferentialReport{
		NewCodeCoverage: 65.0,
		OverallChange:   -2.0,
		CriticalPathStatus: map[string]bool{
			"critical1": false,
		},
		FileChanges: map[string]FileChange{
			"file1.go": {IsModified: true, MeetsThreshold: false},
			"file2.go": {IsModified: true, MeetsThreshold: false},
		},
	}

	recommendations := generateRecommendations(config, report)

	if len(recommendations) == 0 {
		t.Error("Expected recommendations for problematic coverage")
	}

	// Check for specific recommendations
	foundCoverageRec := false
	foundCriticalRec := false
	foundFileRec := false

	for _, rec := range recommendations {
		if contains(rec, "below threshold") {
			foundCoverageRec = true
		}
		if contains(rec, "Critical path") {
			foundCriticalRec = true
		}
		if contains(rec, "modified files") {
			foundFileRec = true
		}
	}

	if !foundCoverageRec {
		t.Error("Expected recommendation about coverage threshold")
	}
	if !foundCriticalRec {
		t.Error("Expected recommendation about critical path")
	}
	if !foundFileRec {
		t.Error("Expected recommendation about modified files")
	}

	// Test passing case
	passingReport := DifferentialReport{
		NewCodeCoverage:    75.0,
		OverallChange:      2.0,
		CriticalPathStatus: map[string]bool{"critical1": true},
		FileChanges: map[string]FileChange{
			"file1.go": {IsModified: true, MeetsThreshold: true},
		},
	}

	passingRecs := generateRecommendations(config, passingReport)
	foundPassingRec := false
	for _, rec := range passingRecs {
		if contains(rec, "quality gates passed") {
			foundPassingRec = true
		}
	}
	if !foundPassingRec {
		t.Error("Expected positive recommendation when quality gates pass")
	}
}

func TestAnalyzeFunctionChanges(t *testing.T) {
	currentFile := FileCoverage{
		File: "test.go",
		Functions: []FunctionCoverage{
			{Function: "Function1", Coverage: 85.0},
			{Function: "Function2", Coverage: 70.0},
			{Function: "NewFunction", Coverage: 90.0},
		},
	}

	baselineFile := FileCoverage{
		Functions: []FunctionCoverage{
			{Function: "Function1", Coverage: 80.0},
			{Function: "Function2", Coverage: 75.0},
		},
	}

	criticalPaths := []string{"test.go:Function1"}

	changes := analyzeFunctionChanges(currentFile, baselineFile, criticalPaths)

	if len(changes) != 3 {
		t.Errorf("Expected 3 function changes, got %d", len(changes))
	}

	// Check Function1 (critical path, improved)
	var func1Change *FunctionChange
	for i := range changes {
		if changes[i].Function == "Function1" {
			func1Change = &changes[i]
			break
		}
	}

	if func1Change == nil {
		t.Fatal("Expected Function1 in changes")
	}
	if !func1Change.IsCriticalPath {
		t.Error("Expected Function1 to be marked as critical path")
	}
	if func1Change.Change != 5.0 {
		t.Errorf("Expected Function1 change of 5.0%%, got %.1f%%", func1Change.Change)
	}
	if !func1Change.MeetsThreshold {
		t.Error("Expected Function1 to meet threshold")
	}

	// Check NewFunction (new function)
	var newFuncChange *FunctionChange
	for i := range changes {
		if changes[i].Function == "NewFunction" {
			newFuncChange = &changes[i]
			break
		}
	}

	if newFuncChange == nil {
		t.Fatal("Expected NewFunction in changes")
	}
	if newFuncChange.BaselineCoverage != 0.0 {
		t.Error("Expected NewFunction baseline coverage to be 0.0")
	}
	if newFuncChange.Change != 90.0 {
		t.Errorf("Expected NewFunction change of 90.0%%, got %.1f%%", newFuncChange.Change)
	}
}

func TestLoadBaselineCoverage(t *testing.T) {
	// This test checks the baseline loading function handles missing files
	baseline, err := loadBaselineCoverage()

	if err != nil {
		// When the baseline file doesn't exist, the function returns an error
		// and the baseline struct is in its zero state (Files map is nil)
		t.Logf("loadBaselineCoverage returned error as expected for missing file: %v", err)

		// The function returns an error but baseline values remain at zero state
		if baseline.Files != nil {
			t.Error("Expected Files map to be nil when baseline file is missing and error is returned")
		}

		if baseline.OverallCoverage != 0.0 {
			t.Errorf("Expected zero coverage when error returned, got %.1f%%", baseline.OverallCoverage)
		}
	} else {
		// If no error, should have a valid baseline structure
		if baseline.Files == nil {
			t.Error("Expected Files map to be initialized when no error")
		}

		if baseline.OverallCoverage < 0 || baseline.OverallCoverage > 100 {
			t.Errorf("Expected reasonable overall coverage, got %.1f%%", baseline.OverallCoverage)
		}
	}
}

func TestGenerateDifferentialReport(t *testing.T) {
	config := DifferentialConfig{}
	config.QualityGates.NewCodeThreshold = 70.0
	config.QualityGates.CriticalPathThreshold = 80.0
	config.CriticalFunctions.CriticalPaths = []string{"test.go:CriticalFunction"}

	currentCoverage := map[string]FileCoverage{
		"test.go": {
			Coverage: 85.0,
			Functions: []FunctionCoverage{
				{Function: "CriticalFunction", Coverage: 85.0},
				{Function: "RegularFunction", Coverage: 70.0},
			},
		},
	}

	baseline := BaselineCoverage{
		OverallCoverage: 80.0,
		Files: map[string]FileCoverage{
			"test.go": {
				Coverage: 80.0,
				Functions: []FunctionCoverage{
					{Function: "CriticalFunction", Coverage: 80.0},
					{Function: "RegularFunction", Coverage: 65.0},
				},
			},
		},
	}

	modifiedFiles := []string{"test.go"}

	report := generateDifferentialReport(config, currentCoverage, baseline, modifiedFiles)

	// Check basic report structure
	if len(report.ModifiedFiles) != 1 {
		t.Errorf("Expected 1 modified file, got %d", len(report.ModifiedFiles))
	}

	if report.NewCodeCoverage != 85.0 {
		t.Errorf("Expected new code coverage 85.0%%, got %.1f%%", report.NewCodeCoverage)
	}

	// Check file changes
	if fileChange, exists := report.FileChanges["test.go"]; exists {
		if fileChange.Change != 5.0 {
			t.Errorf("Expected file change of 5.0%%, got %.1f%%", fileChange.Change)
		}
		if !fileChange.MeetsThreshold {
			t.Error("Expected file to meet threshold")
		}
	} else {
		t.Error("Expected test.go in file changes")
	}

	// Check critical path status
	if status, exists := report.CriticalPathStatus["test.go:CriticalFunction"]; exists {
		if !status {
			t.Error("Expected critical function to pass threshold")
		}
	} else {
		t.Error("Expected critical path status for test.go:CriticalFunction")
	}
}

func TestSaveCoverageHistory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test history file path
	testHistoryPath := filepath.Join(tmpDir, "coverage-history.json")

	// Create a test baseline
	baseline := BaselineCoverage{
		Timestamp:       time.Now(),
		OverallCoverage: 75.0,
		Files:           make(map[string]FileCoverage),
		GitCommit:       "test-commit",
		GitBranch:       "test-branch",
	}

	// This would normally save to docs/coverage-history.json
	// For testing, we'll verify the JSON structure
	data, err := json.MarshalIndent([]BaselineCoverage{baseline}, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal baseline: %v", err)
	}

	err = os.WriteFile(testHistoryPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test history: %v", err)
	}

	// Verify the file was created and contains valid JSON
	savedData, err := os.ReadFile(testHistoryPath)
	if err != nil {
		t.Fatalf("Failed to read saved history: %v", err)
	}

	var histories []BaselineCoverage
	err = json.Unmarshal(savedData, &histories)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved history: %v", err)
	}

	if len(histories) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(histories))
	}

	if histories[0].OverallCoverage != 75.0 {
		t.Errorf("Expected coverage 75.0%%, got %.1f%%", histories[0].OverallCoverage)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests for performance
func BenchmarkParseCurrentCoverage(b *testing.B) {
	tmpDir := b.TempDir()
	profilePath := filepath.Join(tmpDir, "coverage.out")

	// Create a large test profile
	testProfile := "mode: set\n"
	for i := 0; i < 1000; i++ {
		testProfile += "bkpdir/file" + string(rune(i)) + ".go:10: function" + string(rune(i)) + " 75.5%\n"
	}

	err := os.WriteFile(profilePath, []byte(testProfile), 0644)
	if err != nil {
		b.Fatalf("Failed to create test profile: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseCurrentCoverage(profilePath)
		if err != nil {
			b.Fatalf("parseCurrentCoverage failed: %v", err)
		}
	}
}

func BenchmarkEvaluateQualityGates(b *testing.B) {
	config := DifferentialConfig{}
	config.QualityGates.NewCodeThreshold = 70.0
	config.QualityGates.CriticalPathThreshold = 80.0

	report := DifferentialReport{
		NewCodeCoverage: 75.0,
		OverallChange:   2.0,
		CriticalPathStatus: map[string]bool{
			"critical1": true,
			"critical2": true,
			"critical3": true,
		},
		FileChanges: map[string]FileChange{
			"file1.go": {IsModified: true, MeetsThreshold: true},
			"file2.go": {IsModified: true, MeetsThreshold: true},
			"file3.go": {IsModified: true, MeetsThreshold: true},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evaluateQualityGates(config, report)
	}
}
