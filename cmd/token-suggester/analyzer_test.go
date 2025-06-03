// 🔶 DOC-010: Token suggestion engine tests - 🧪 Comprehensive testing suite
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// 🔶 DOC-010: Test analyzer creation - 🧪 Analyzer initialization testing
func TestNewTokenAnalyzer(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created, got nil")
	}

	if analyzer.config == nil {
		t.Error("Expected config to be initialized")
	}

	if analyzer.fileSet == nil {
		t.Error("Expected fileSet to be initialized")
	}

	if analyzer.featureMap == nil {
		t.Error("Expected featureMap to be initialized")
	}
}

// 🔶 DOC-010: Test validator creation - 🧪 Validator initialization testing
func TestNewTokenValidator(t *testing.T) {
	validator := NewTokenValidator()

	if validator == nil {
		t.Fatal("Expected validator to be created, got nil")
	}

	if validator.config == nil {
		t.Error("Expected config to be initialized")
	}
}

// 🔶 DOC-010: Test batch processor creation - 🧪 Batch processor initialization testing
func TestNewBatchProcessor(t *testing.T) {
	processor := NewBatchProcessor()

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}

	if processor.validator == nil {
		t.Error("Expected validator to be initialized")
	}
}

// 🔶 DOC-010: Test priority determination - 🧪 Priority assignment testing
func TestDeterminePriority(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name           string
		functionName   string
		context        map[string]string
		expectedIcon   string
		expectedReason string
	}{
		{
			name:           "Critical function - main",
			functionName:   "MainFunction",
			context:        map[string]string{},
			expectedIcon:   "⭐",
			expectedReason: "Critical operation: main",
		},
		{
			name:           "High priority - config",
			functionName:   "LoadConfig",
			context:        map[string]string{},
			expectedIcon:   "🔺",
			expectedReason: "High priority: config",
		},
		{
			name:           "Medium priority - format",
			functionName:   "FormatOutput",
			context:        map[string]string{},
			expectedIcon:   "🔶",
			expectedReason: "Medium priority: format",
		},
		{
			name:           "Context-based high priority - error handling",
			functionName:   "ProcessData",
			context:        map[string]string{"error_handling": "true"},
			expectedIcon:   "🔺",
			expectedReason: "High priority: error handling function",
		},
		{
			name:           "Default low priority",
			functionName:   "UtilityFunction",
			context:        map[string]string{},
			expectedIcon:   "🔻",
			expectedReason: "Low priority: utility function",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			icon, reason := analyzer.determinePriority(signature, tt.context)

			if icon != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, icon)
			}

			if !strings.Contains(reason, strings.Split(tt.expectedReason, ":")[0]) {
				t.Errorf("Expected reason to contain %s, got %s", tt.expectedReason, reason)
			}
		})
	}
}

// 🔶 DOC-010: Test action determination - 🧪 Action assignment testing
func TestDetermineAction(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name           string
		functionName   string
		context        map[string]string
		expectedIcon   string
		expectedReason string
	}{
		{
			name:           "Analysis function - get",
			functionName:   "GetConfig",
			context:        map[string]string{},
			expectedIcon:   "🔍",
			expectedReason: "Analysis operation: get",
		},
		{
			name:           "Documentation function - format",
			functionName:   "FormatOutput",
			context:        map[string]string{},
			expectedIcon:   "📝",
			expectedReason: "Documentation operation: format",
		},
		{
			name:           "Configuration function - create",
			functionName:   "CreateArchive",
			context:        map[string]string{},
			expectedIcon:   "🔧",
			expectedReason: "Configuration operation: create",
		},
		{
			name:           "Protection function - validate",
			functionName:   "ValidateInput",
			context:        map[string]string{},
			expectedIcon:   "🛡️",
			expectedReason: "Protection operation: validate",
		},
		{
			name:           "Context-based protection - error handling",
			functionName:   "ProcessRequest",
			context:        map[string]string{"error_handling": "true"},
			expectedIcon:   "🛡️",
			expectedReason: "Protection: error handling",
		},
		{
			name:           "Default configuration",
			functionName:   "UtilityFunction",
			context:        map[string]string{},
			expectedIcon:   "🔧",
			expectedReason: "Configuration: general operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			icon, reason := analyzer.determineAction(signature, tt.context)

			if icon != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, icon)
			}

			if !strings.Contains(reason, strings.Split(tt.expectedReason, ":")[0]) {
				t.Errorf("Expected reason to contain %s, got %s", tt.expectedReason, reason)
			}
		})
	}
}

// 🔶 DOC-010: Test feature ID determination - 🧪 Feature mapping testing
func TestDetermineFeatureID(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name         string
		functionName string
		context      map[string]string
		expected     string
	}{
		{
			name:         "Config function",
			functionName: "LoadConfig",
			context:      map[string]string{},
			expected:     "CFG-NEW",
		},
		{
			name:         "Archive function",
			functionName: "CreateArchive",
			context:      map[string]string{},
			expected:     "ARCH-NEW",
		},
		{
			name:         "Backup function",
			functionName: "BackupFile",
			context:      map[string]string{},
			expected:     "FILE-NEW",
		},
		{
			name:         "Git function",
			functionName: "GitCommit",
			context:      map[string]string{},
			expected:     "GIT-NEW",
		},
		{
			name:         "Test function",
			functionName: "TestValidation",
			context:      map[string]string{},
			expected:     "TEST-NEW",
		},
		{
			name:         "Generic function",
			functionName: "UtilityFunction",
			context:      map[string]string{},
			expected:     "UTIL-NEW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := FunctionSignature{Name: tt.functionName}
			featureID := analyzer.determineFeatureID(signature, tt.context)

			if featureID != tt.expected {
				t.Errorf("Expected feature ID %s, got %s", tt.expected, featureID)
			}
		})
	}
}

// 🔶 DOC-010: Test confidence calculation - 🧪 Confidence scoring testing
func TestCalculateConfidence(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	tests := []struct {
		name       string
		signature  FunctionSignature
		context    map[string]string
		suggestion *TokenSuggestion
		minConf    float64
		maxConf    float64
	}{
		{
			name: "High confidence - exported with parameters",
			signature: FunctionSignature{
				Name:       "ProcessData",
				IsExported: true,
				Parameters: []Parameter{{Name: "data", Type: "string"}},
				ReturnType: "error",
			},
			context:    map[string]string{"error_handling": "true", "resource_management": "true"},
			suggestion: &TokenSuggestion{FeatureID: "ARCH-001", PriorityReason: "High priority"},
			minConf:    0.8,
			maxConf:    1.0,
		},
		{
			name: "Low confidence - unexported no parameters",
			signature: FunctionSignature{
				Name:       "helper",
				IsExported: false,
				Parameters: []Parameter{},
				ReturnType: "void",
			},
			context:    map[string]string{},
			suggestion: &TokenSuggestion{FeatureID: "UTIL-NEW", PriorityReason: "Low priority: utility function"},
			minConf:    0.0,
			maxConf:    0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confidence := analyzer.calculateConfidence(tt.signature, tt.context, tt.suggestion)

			if confidence < tt.minConf || confidence > tt.maxConf {
				t.Errorf("Expected confidence between %f and %f, got %f", tt.minConf, tt.maxConf, confidence)
			}
		})
	}
}

// 🔶 DOC-010: Test complexity analysis - 🧪 Complexity assessment testing
func TestAnalyzeComplexity(t *testing.T) {
	// This would require creating AST nodes, which is complex for testing
	// In a real implementation, we would mock the AST or use test fixtures
	t.Skip("Complexity analysis requires AST fixtures - implementation would include proper AST testing")
}

// 🔶 DOC-010: Test context analysis - 🧪 Context analysis testing
func TestAnalyzeContext(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// 🔶 DOC-010: Test source code - 📝 Sample code for analysis
	testCode := `package main

import "fmt"

func ProcessData(data string) error {
	if data == "" {
		return fmt.Errorf("empty data")
	}
	
	defer cleanup()
	
	config := loadConfig()
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	
	return nil
}`

	context, err := analyzer.analyzeContext("test.go", 5, []byte(testCode))
	if err != nil {
		t.Fatalf("Context analysis failed: %v", err)
	}

	// 🔶 DOC-010: Context validation - 🛡️ Expected patterns
	expectedPatterns := map[string]string{
		"error_handling":        "true",
		"resource_management":   "true",
		"configuration_related": "true",
		"file_operations":       "true",
	}

	for key, expected := range expectedPatterns {
		if context[key] != expected {
			t.Errorf("Expected %s to be %s, got %s", key, expected, context[key])
		}
	}

	if context["surrounding_code"] == "" {
		t.Error("Expected surrounding_code to be populated")
	}
}

// 🔶 DOC-010: Test token validation - 🧪 Validation testing
func TestValidateFile(t *testing.T) {
	validator := NewTokenValidator()

	// 🔶 DOC-010: Create temporary test file - 📁 Test file setup
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testContent := `package main

// Valid token format
// ⭐ ARCH-001: Archive creation - 🔧 Core functionality
func CreateArchive() {}

// Invalid token format - missing icons
// ARCH-002: Archive validation
func ValidateArchive() {}

// Invalid token format - wrong structure
// Some random comment with ARCH-003
func ProcessArchive() {}

// Valid but different format
// 🔺 CFG-001: Configuration loading - 🔍 Discovery operation
func LoadConfig() {}`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 🔶 DOC-010: Validate test file - 🛡️ Validation execution
	violations, err := validator.validateFile(testFile)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	// 🔶 DOC-010: Violation analysis - 📊 Expected violations
	if len(violations) == 0 {
		t.Error("Expected violations to be found")
	}

	// Check specific violations
	foundInvalidFormat := false
	for _, violation := range violations {
		if violation.ViolationType == "INVALID_FORMAT" {
			foundInvalidFormat = true
		}

		if violation.Severity != "WARNING" {
			t.Errorf("Expected severity WARNING, got %s", violation.Severity)
		}

		if violation.RuleID != "FORMAT-001" {
			t.Errorf("Expected rule ID FORMAT-001, got %s", violation.RuleID)
		}
	}

	if !foundInvalidFormat {
		t.Error("Expected INVALID_FORMAT violation to be found")
	}
}

// 🔶 DOC-010: Test file analysis integration - 🧪 Integration testing
func TestAnalyzeTargetFile(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// 🔶 DOC-010: Create test Go file - 📁 Test file creation
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testContent := `package main

import "fmt"

// CreateBackup creates a backup of the specified file
func CreateBackup(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}
	
	config := loadConfig()
	return processBackup(filePath, config)
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	return &Config{}, nil
}

// validateInput is unexported helper
func validateInput(data string) bool {
	return data != ""
}`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 🔶 DOC-010: Analyze test file - 🔍 File analysis execution
	results, err := analyzer.AnalyzeTarget(testFile)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// 🔶 DOC-010: Results validation - 📊 Analysis verification
	if results.Target != testFile {
		t.Errorf("Expected target %s, got %s", testFile, results.Target)
	}

	if results.FunctionsAnalyzed != 2 { // Only exported functions
		t.Errorf("Expected 2 functions analyzed, got %d", results.FunctionsAnalyzed)
	}

	if len(results.Suggestions) != 2 {
		t.Errorf("Expected 2 suggestions, got %d", len(results.Suggestions))
	}

	if results.ProcessingTime == 0 {
		t.Error("Expected processing time to be recorded")
	}

	// 🔶 DOC-010: Suggestion content validation - 💡 Suggestion quality check
	for _, suggestion := range results.Suggestions {
		if suggestion.FilePath != testFile {
			t.Errorf("Expected file path %s, got %s", testFile, suggestion.FilePath)
		}

		if suggestion.Confidence < 0 || suggestion.Confidence > 1 {
			t.Errorf("Expected confidence between 0 and 1, got %f", suggestion.Confidence)
		}

		if suggestion.SuggestedToken == "" {
			t.Error("Expected suggested token to be generated")
		}

		if suggestion.PriorityIcon == "" {
			t.Error("Expected priority icon to be assigned")
		}

		if suggestion.ActionIcon == "" {
			t.Error("Expected action icon to be assigned")
		}

		if suggestion.FeatureID == "" {
			t.Error("Expected feature ID to be assigned")
		}
	}
}

// 🔶 DOC-010: Test batch processing - 🧪 Batch processing testing
func TestBatchProcessing(t *testing.T) {
	processor := NewBatchProcessor()

	// 🔶 DOC-010: Create test directory structure - 📁 Test setup
	tmpDir := t.TempDir()

	// Create multiple test files
	testFiles := map[string]string{
		"file1.go": `package main
func CreateArchive() error { return nil }
func LoadConfig() error { return nil }`,

		"file2.go": `package main
func ProcessBackup(path string) error { return nil }
func ValidateData(data []byte) bool { return true }`,

		"file3_test.go": `package main
func TestSomething() {}`, // This should be skipped
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// 🔶 DOC-010: Execute batch processing - 🚀 Batch execution
	results, err := processor.ProcessDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Batch processing failed: %v", err)
	}

	// 🔶 DOC-010: Batch results validation - 📊 Results verification
	if results.Directory != tmpDir {
		t.Errorf("Expected directory %s, got %s", tmpDir, results.Directory)
	}

	if results.FilesProcessed != 2 { // Should skip test file
		t.Errorf("Expected 2 files processed, got %d", results.FilesProcessed)
	}

	if results.TotalFunctions != 4 { // 2 + 2 exported functions
		t.Errorf("Expected 4 functions total, got %d", results.TotalFunctions)
	}

	if results.TotalSuggestions != 4 {
		t.Errorf("Expected 4 suggestions total, got %d", results.TotalSuggestions)
	}

	if results.ProcessingTime == 0 {
		t.Error("Expected processing time to be recorded")
	}

	// 🔶 DOC-010: Priority and action breakdown validation - 📊 Category verification
	totalPriority := results.PriorityBreakdown.Critical + results.PriorityBreakdown.High +
		results.PriorityBreakdown.Medium + results.PriorityBreakdown.Low
	if totalPriority != results.TotalSuggestions {
		t.Errorf("Priority breakdown doesn't match total suggestions: %d vs %d", totalPriority, results.TotalSuggestions)
	}

	totalAction := results.ActionBreakdown.Analysis + results.ActionBreakdown.Documentation +
		results.ActionBreakdown.Configuration + results.ActionBreakdown.Protection
	if totalAction != results.TotalSuggestions {
		t.Errorf("Action breakdown doesn't match total suggestions: %d vs %d", totalAction, results.TotalSuggestions)
	}

	// 🔶 DOC-010: Top suggestions validation - 💡 Quality ranking verification
	if len(results.TopSuggestions) > 10 {
		t.Error("Expected top suggestions to be limited to 10")
	}

	// Check suggestions are sorted by confidence
	for i := 1; i < len(results.TopSuggestions); i++ {
		if results.TopSuggestions[i-1].Confidence < results.TopSuggestions[i].Confidence {
			t.Error("Expected top suggestions to be sorted by confidence (descending)")
		}
	}
}

// 🔶 DOC-010: Test error handling - 🧪 Error condition testing
func TestErrorHandling(t *testing.T) {
	analyzer := NewTokenAnalyzer()

	// 🔶 DOC-010: Test non-existent file - 🛡️ Error case testing
	_, err := analyzer.AnalyzeTarget("non-existent-file.go")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// 🔶 DOC-010: Test invalid Go file - 🛡️ Parser error testing
	tmpDir := t.TempDir()
	invalidFile := filepath.Join(tmpDir, "invalid.go")

	err = os.WriteFile(invalidFile, []byte("invalid go syntax {{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	_, err = analyzer.AnalyzeTarget(invalidFile)
	if err == nil {
		t.Error("Expected error for invalid Go file")
	}
}

// 🔶 DOC-010: Test configuration defaults - 🧪 Configuration testing
func TestDefaultAnalysisConfig(t *testing.T) {
	config := DefaultAnalysisConfig()

	if config == nil {
		t.Fatal("Expected config to be created")
	}

	// 🔶 DOC-010: Validate priority rules - ⭐ Priority configuration validation
	if len(config.PriorityRules.CriticalPatterns) == 0 {
		t.Error("Expected critical patterns to be defined")
	}

	if len(config.PriorityRules.HighPatterns) == 0 {
		t.Error("Expected high patterns to be defined")
	}

	// 🔶 DOC-010: Validate action rules - 🔧 Action configuration validation
	if len(config.ActionRules.AnalysisPatterns) == 0 {
		t.Error("Expected analysis patterns to be defined")
	}

	if len(config.ActionRules.ConfigurationPatterns) == 0 {
		t.Error("Expected configuration patterns to be defined")
	}

	// 🔶 DOC-010: Validate confidence weights - 📊 Weight configuration validation
	totalWeight := config.ConfidenceWeights.SignatureMatch +
		config.ConfidenceWeights.PatternMatch +
		config.ConfidenceWeights.ContextMatch +
		config.ConfidenceWeights.FeatureMapping +
		config.ConfidenceWeights.ComplexityFactor

	if totalWeight != 1.0 {
		t.Errorf("Expected confidence weights to sum to 1.0, got %f", totalWeight)
	}

	// 🔶 DOC-010: Validate validation rules - 🛡️ Validation configuration
	if len(config.ValidationRules.ValidPriorityIcons) != 4 {
		t.Errorf("Expected 4 priority icons, got %d", len(config.ValidationRules.ValidPriorityIcons))
	}

	if len(config.ValidationRules.ValidActionIcons) != 4 {
		t.Errorf("Expected 4 action icons, got %d", len(config.ValidationRules.ValidActionIcons))
	}

	if config.ValidationRules.MinConfidence < 0 || config.ValidationRules.MinConfidence > 1 {
		t.Errorf("Expected min confidence between 0 and 1, got %f", config.ValidationRules.MinConfidence)
	}
}

// 🔶 DOC-010: Test utility functions - 🧪 Utility function testing
func TestUtilityFunctions(t *testing.T) {
	// 🔶 DOC-010: Test integer min/max - 🔧 Integer utilities
	if min(5, 3) != 3 {
		t.Error("min function failed for integers")
	}

	if max(5, 3) != 5 {
		t.Error("max function failed for integers")
	}

	// 🔶 DOC-010: Test float64 min/max - 🔧 Float utilities
	if minFloat64(5.5, 3.3) != 3.3 {
		t.Error("minFloat64 function failed")
	}

	if maxFloat64(5.5, 3.3) != 5.5 {
		t.Error("maxFloat64 function failed")
	}
}

// 🔶 DOC-010: Benchmark analysis performance - ⚡ Performance testing
func BenchmarkAnalyzeTarget(b *testing.B) {
	analyzer := NewTokenAnalyzer()

	// 🔶 DOC-010: Create benchmark test file - 📁 Performance test setup
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "benchmark.go")

	testContent := `package main

func CreateArchive(path string) error { return nil }
func LoadConfig() (*Config, error) { return nil, nil }
func ProcessBackup(data []byte) error { return nil }
func ValidateInput(input string) bool { return true }
func FormatOutput(data interface{}) string { return "" }`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	// 🔶 DOC-010: Execute benchmark - ⚡ Performance measurement
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := analyzer.AnalyzeTarget(testFile)
		if err != nil {
			b.Fatalf("Analysis failed: %v", err)
		}
	}
}

// 🔶 DOC-010: Benchmark batch processing performance - ⚡ Batch performance testing
func BenchmarkBatchProcessing(b *testing.B) {
	processor := NewBatchProcessor()

	// 🔶 DOC-010: Create benchmark directory - 📁 Batch performance setup
	tmpDir := b.TempDir()

	// Create multiple files for batch processing
	for i := 0; i < 10; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("file%d.go", i))
		content := fmt.Sprintf(`package main
func Function%dA() error { return nil }
func Function%dB() bool { return true }`, i, i)

		err := os.WriteFile(filename, []byte(content), 0644)
		if err != nil {
			b.Fatalf("Failed to create benchmark file: %v", err)
		}
	}

	// 🔶 DOC-010: Execute batch benchmark - ⚡ Batch performance measurement
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.ProcessDirectory(tmpDir)
		if err != nil {
			b.Fatalf("Batch processing failed: %v", err)
		}
	}
}
