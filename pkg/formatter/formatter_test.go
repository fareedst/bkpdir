// Tests for the pkg/formatter package to validate extracted formatting functionality.
// These tests ensure the extracted formatter components work correctly
// and maintain backward compatibility with the original functionality.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

import (
	"errors"
	"strings"
	"testing"
)

// ⭐ EXTRACT-003: Test configuration provider - 🧪 Mock for testing
// MockConfigProvider provides test configuration values
type MockConfigProvider struct {
	formatStrings   map[string]string
	templateStrings map[string]string
	patterns        map[string]string
	errorFormats    map[string]string
}

// NewMockConfigProvider creates a new MockConfigProvider with default values
func NewMockConfigProvider() *MockConfigProvider {
	return &MockConfigProvider{
		formatStrings:   make(map[string]string),
		templateStrings: make(map[string]string),
		patterns:        make(map[string]string),
		errorFormats:    make(map[string]string),
	}
}

func (mcp *MockConfigProvider) GetFormatString(formatType string) string {
	return mcp.formatStrings[formatType]
}

func (mcp *MockConfigProvider) GetTemplateString(templateType string) string {
	return mcp.templateStrings[templateType]
}

func (mcp *MockConfigProvider) GetPattern(patternType string) string {
	return mcp.patterns[patternType]
}

func (mcp *MockConfigProvider) GetErrorFormat(errorType string) string {
	return mcp.errorFormats[errorType]
}

// ⭐ EXTRACT-003: OutputCollector tests - 🧪 Delayed output functionality
func TestOutputCollector(t *testing.T) {
	collector := NewOutputCollector()

	// Test initial state
	if len(collector.GetMessages()) != 0 {
		t.Errorf("New collector should have no messages")
	}

	// Test adding messages
	collector.AddStdout("test stdout", "info")
	collector.AddStderr("test stderr", "error")

	messages := collector.GetMessages()
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(messages))
	}

	// Verify message content
	if messages[0].Content != "test stdout" || messages[0].Destination != "stdout" {
		t.Errorf("First message incorrect: %+v", messages[0])
	}
	if messages[1].Content != "test stderr" || messages[1].Destination != "stderr" {
		t.Errorf("Second message incorrect: %+v", messages[1])
	}

	// Test clear
	collector.Clear()
	if len(collector.GetMessages()) != 0 {
		t.Errorf("Collector should be empty after clear")
	}
}

// ⭐ EXTRACT-003: PatternExtractor tests - 🧪 Regex pattern extraction
func TestPatternExtractor(t *testing.T) {
	configProvider := NewMockConfigProvider()
	extractor := NewDefaultPatternExtractor(configProvider)

	// Test archive filename extraction with default pattern
	filename := "myproject_2024-01-15_14-30-45.zip"
	data := extractor.ExtractArchiveFilenameData(filename)

	expectedKeys := []string{"name", "timestamp", "suffix"}
	for _, key := range expectedKeys {
		if _, exists := data[key]; !exists {
			t.Errorf("Expected key %s not found in extracted data", key)
		}
	}

	if data["name"] != "myproject" {
		t.Errorf("Expected name 'myproject', got '%s'", data["name"])
	}
	if data["timestamp"] != "2024-01-15_14-30-45" {
		t.Errorf("Expected timestamp '2024-01-15_14-30-45', got '%s'", data["timestamp"])
	}

	// Test simple pattern extractor
	simpleExtractor := NewSimplePatternExtractor()
	simpleData := simpleExtractor.ExtractArchiveFilenameData(filename)

	if simpleData["name"] != data["name"] {
		t.Errorf("Simple extractor should produce same results")
	}
}

// ⭐ EXTRACT-003: TemplateFormatter tests - 🧪 Template-based formatting
func TestTemplateFormatter(t *testing.T) {
	configProvider := NewMockConfigProvider()
	formatter := NewDefaultTemplateFormatter(configProvider)

	// Test placeholder formatting
	template := "Created archive: %{path} at %{timestamp}"
	data := map[string]string{
		"path":      "/test/archive.zip",
		"timestamp": "2024-01-15_14-30-45",
	}

	result := formatter.FormatWithPlaceholders(template, data)
	expected := "Created archive: /test/archive.zip at 2024-01-15_14-30-45"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	// Test template with pattern
	input := "myproject_2024-01-15_14-30-45.zip"
	pattern := `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	tmplStr := "Archive %{name} created at %{timestamp}"

	templateResult, err := formatter.FormatWithTemplate(input, pattern, tmplStr)
	if err != nil {
		t.Errorf("Template formatting failed: %v", err)
	}

	expectedTemplate := "Archive myproject created at 2024-01-15_14-30-45"
	if templateResult != expectedTemplate {
		t.Errorf("Expected '%s', got '%s'", expectedTemplate, templateResult)
	}

	// Test simple template formatter
	simpleFormatter := NewSimpleTemplateFormatter()
	simpleResult := simpleFormatter.FormatWithPlaceholders(template, data)
	if simpleResult != result {
		t.Errorf("Simple formatter should produce same results")
	}
}

// ⭐ EXTRACT-003: DefaultOutputFormatter tests - 🧪 Main formatter functionality
func TestDefaultOutputFormatter(t *testing.T) {
	configProvider := NewMockConfigProvider()
	formatter := NewDefaultOutputFormatter(configProvider)

	// Test basic formatting with defaults
	result := formatter.FormatCreatedArchive("/test/archive.zip")
	if !strings.Contains(result, "/test/archive.zip") {
		t.Errorf("Formatted message should contain path")
	}
	if !strings.Contains(result, "Created archive") {
		t.Errorf("Formatted message should contain 'Created archive'")
	}

	// Test with custom format string
	configProvider.formatStrings["created_archive"] = "Archive created: %s"
	customResult := formatter.FormatCreatedArchive("/test/archive.zip")
	if customResult != "Archive created: /test/archive.zip" {
		t.Errorf("Expected custom format, got '%s'", customResult)
	}

	// Test error formatting
	testErr := errors.New("test error")
	errorResult := formatter.FormatDiskFullError(testErr)
	if !strings.Contains(errorResult, "test error") {
		t.Errorf("Error message should contain original error")
	}
	if !strings.Contains(errorResult, "Disk full") {
		t.Errorf("Error message should contain error type")
	}
}

// ⭐ EXTRACT-003: Delayed output mode tests - 🧪 OutputCollector integration
func TestDelayedOutputMode(t *testing.T) {
	configProvider := NewMockConfigProvider()
	collector := NewOutputCollector()
	formatter := NewDefaultOutputFormatterWithCollector(configProvider, collector)

	// Test delayed mode detection
	if !formatter.IsDelayedMode() {
		t.Errorf("Formatter should be in delayed mode")
	}

	if formatter.GetCollector() != collector {
		t.Errorf("Formatter should return the same collector")
	}

	// Test print operations in delayed mode
	formatter.PrintCreatedArchive("/test/archive.zip")
	formatter.PrintError("test error")

	messages := collector.GetMessages()
	if len(messages) != 2 {
		t.Errorf("Expected 2 messages in collector, got %d", len(messages))
	}

	// Verify message destinations
	stdoutCount := 0
	stderrCount := 0
	for _, msg := range messages {
		if msg.Destination == "stdout" {
			stdoutCount++
		} else if msg.Destination == "stderr" {
			stderrCount++
		}
	}

	if stdoutCount != 1 || stderrCount != 1 {
		t.Errorf("Expected 1 stdout and 1 stderr message, got %d stdout, %d stderr", stdoutCount, stderrCount)
	}
}

// ⭐ EXTRACT-003: Template delegation tests - 🧪 Component integration
func TestTemplateDelegation(t *testing.T) {
	configProvider := NewMockConfigProvider()
	formatter := NewDefaultOutputFormatter(configProvider)

	// Test template operations delegation
	data := map[string]string{"path": "/test/archive.zip"}
	result := formatter.TemplateCreatedArchive(data)

	if !strings.Contains(result, "/test/archive.zip") {
		t.Errorf("Template result should contain path")
	}

	// Test pattern extraction delegation
	filename := "myproject_2024-01-15_14-30-45.zip"
	extractedData := formatter.ExtractArchiveFilenameData(filename)

	if extractedData["name"] != "myproject" {
		t.Errorf("Pattern extraction delegation failed")
	}
}

// ⭐ EXTRACT-003: Error template formatting tests - 🧪 Template error handling
func TestErrorTemplateFormatting(t *testing.T) {
	configProvider := NewMockConfigProvider()
	formatter := NewDefaultOutputFormatter(configProvider)

	testErr := errors.New("permission denied")

	// Test template error formatting
	result := formatter.TemplatePermissionError(testErr)
	if !strings.Contains(result, "permission denied") {
		t.Errorf("Template error should contain original error message")
	}
	if !strings.Contains(result, "Permission denied") {
		t.Errorf("Template error should contain error type")
	}

	// Test with custom template
	configProvider.templateStrings["error_permission"] = "Access denied: %{error}"
	customResult := formatter.TemplatePermissionError(testErr)
	if !strings.Contains(customResult, "Access denied") {
		t.Errorf("Custom template should be used")
	}
}

// ⭐ EXTRACT-003: Pattern extraction edge cases - 🧪 Robust pattern handling
func TestPatternExtractionEdgeCases(t *testing.T) {
	configProvider := NewMockConfigProvider()
	extractor := NewDefaultPatternExtractor(configProvider)

	// Test with invalid pattern
	configProvider.patterns["archive_filename"] = "[invalid regex"
	data := extractor.ExtractArchiveFilenameData("test.zip")
	if len(data) != 0 {
		t.Errorf("Invalid pattern should return empty data")
	}

	// Test with non-matching filename
	configProvider.patterns["archive_filename"] = `^(?P<name>.*?)_(?P<timestamp>\d{4})$`
	data = extractor.ExtractArchiveFilenameData("nomatch.zip")
	if len(data) != 0 {
		t.Errorf("Non-matching filename should return empty data")
	}

	// Test utility function
	path := "/path/to/file.txt"
	filename := GetFilenameFromPath(path)
	if filename != "file.txt" {
		t.Errorf("Expected 'file.txt', got '%s'", filename)
	}
}

// ⭐ EXTRACT-003: Interface compliance tests - 🧪 Interface implementation validation
func TestInterfaceCompliance(t *testing.T) {
	configProvider := NewMockConfigProvider()

	// Test that DefaultOutputFormatter implements all interfaces
	var formatter OutputFormatterInterface = NewDefaultOutputFormatter(configProvider)

	// Test Formatter interface
	_ = formatter.FormatCreatedArchive("/test")
	_ = formatter.FormatError("test")

	// Test TemplateFormatter interface
	data := map[string]string{"path": "/test"}
	_ = formatter.TemplateCreatedArchive(data)

	// Test PatternExtractor interface
	_ = formatter.ExtractArchiveFilenameData("test.zip")

	// Test PrintFormatter interface (should not panic)
	formatter.PrintCreatedArchive("/test")

	// Test ErrorFormatter interface
	testErr := errors.New("test")
	_ = formatter.FormatDiskFullError(testErr)
	_ = formatter.TemplateDiskFullError(testErr)

	// Test delayed output interface
	_ = formatter.IsDelayedMode()
	_ = formatter.GetCollector()
	formatter.SetCollector(NewOutputCollector())
}

// ⭐ EXTRACT-003: Performance benchmark tests - 🧪 Performance validation
func BenchmarkFormatterOperations(b *testing.B) {
	configProvider := NewMockConfigProvider()
	formatter := NewDefaultOutputFormatter(configProvider)

	b.Run("FormatCreatedArchive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			formatter.FormatCreatedArchive("/test/archive.zip")
		}
	})

	b.Run("TemplateFormatting", func(b *testing.B) {
		data := map[string]string{"path": "/test/archive.zip"}
		for i := 0; i < b.N; i++ {
			formatter.TemplateCreatedArchive(data)
		}
	})

	b.Run("PatternExtraction", func(b *testing.B) {
		filename := "myproject_2024-01-15_14-30-45.zip"
		for i := 0; i < b.N; i++ {
			formatter.ExtractArchiveFilenameData(filename)
		}
	})
}
