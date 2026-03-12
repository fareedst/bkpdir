// AI-First Formatter Tests
// Comprehensive tests for AI-first formatter implementation to validate
// clear component separation and interface standardization.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:core-functionality]
package formatter

import (
	"errors"
	"testing"
	"time"
)

// [CRITICAL] FMT-001: Mock configuration for testing - [ACTION:configure-modify]
type MockFormatterConfig struct {
	formatStrings   map[FormatType]string
	templateStrings map[FormatType]string
	errorFormats    map[ErrorType]string
	patterns        map[PatternType]string
}

func NewMockFormatterConfig() *MockFormatterConfig {
	return &MockFormatterConfig{
		formatStrings:   make(map[FormatType]string),
		templateStrings: make(map[FormatType]string),
		errorFormats:    make(map[ErrorType]string),
		patterns:        make(map[PatternType]string),
	}
}

func (m *MockFormatterConfig) GetFormatString(formatType FormatType) (string, error) {
	if format, ok := m.formatStrings[formatType]; ok {
		return format, nil
	}
	return "", nil
}

func (m *MockFormatterConfig) GetTemplateString(templateType FormatType) (string, error) {
	if template, ok := m.templateStrings[templateType]; ok {
		return template, nil
	}
	return "", nil
}

func (m *MockFormatterConfig) GetErrorFormat(errorType ErrorType) (string, error) {
	if format, ok := m.errorFormats[errorType]; ok {
		return format, nil
	}
	return "", nil
}

func (m *MockFormatterConfig) GetPattern(patternType PatternType) (string, error) {
	if pattern, ok := m.patterns[patternType]; ok {
		return pattern, nil
	}
	return "", nil
}

func (m *MockFormatterConfig) Validate() error {
	return nil
}

func (m *MockFormatterConfig) GetValidationErrors() []ConfigError {
	return []ConfigError{}
}

// [CRITICAL] FMT-001: Test AI-first formatter creation - [ACTION:core-functionality]
func TestNewAIFirstFormatter(t *testing.T) {
	config := NewMockFormatterConfig()
	formatter := NewAIFirstFormatter(config)

	if formatter == nil {
		t.Fatal("NewAIFirstFormatter returned nil")
	}

	if formatter.GetConfig() != config {
		t.Error("formatter config does not match provided config")
	}
}

// [CRITICAL] FMT-001: Test AI-first formatter with collector - [ACTION:core-functionality]
func TestNewAIFirstFormatterWithCollector(t *testing.T) {
	config := NewMockFormatterConfig()
	collector := NewOutputCollector()
	formatter := NewAIFirstFormatterWithCollector(config, collector)

	if formatter == nil {
		t.Fatal("NewAIFirstFormatterWithCollector returned nil")
	}

	if !formatter.IsDelayedMode() {
		t.Error("formatter should be in delayed mode with collector")
	}
}

// [CRITICAL] FMT-001: Test configuration management - [ACTION:configure-modify]
func TestAIFirstFormatterConfigManagement(t *testing.T) {
	config1 := NewMockFormatterConfig()
	config2 := NewMockFormatterConfig()
	formatter := NewAIFirstFormatter(config1)

	// Test initial config
	if formatter.GetConfig() != config1 {
		t.Error("initial config does not match")
	}

	// Test config update
	err := formatter.SetConfig(config2)
	if err != nil {
		t.Errorf("SetConfig failed: %v", err)
	}

	if formatter.GetConfig() != config2 {
		t.Error("updated config does not match")
	}

	// Test nil config rejection
	err = formatter.SetConfig(nil)
	if err == nil {
		t.Error("SetConfig should reject nil config")
	}
}

// [CRITICAL] FMT-001: Test format with context - [ACTION:core-functionality]
func TestAIFirstFormatterFormatWithContext(t *testing.T) {
	config := NewMockFormatterConfig()
	formatter := NewAIFirstFormatter(config)

	// Test created archive formatting
	ctx := FormatContext{
		FormatType: FormatTypeCreated,
		Data: map[string]interface{}{
			"path": "/test/archive.tar.gz",
		},
		Options:  FormatOptions{},
		Metadata: make(map[string]string),
	}

	result, err := formatter.FormatWithContext(ctx)
	if err != nil {
		t.Errorf("FormatWithContext failed: %v", err)
	}

	if result == "" {
		t.Error("FormatWithContext returned empty result")
	}

	// Test missing path error
	ctx.Data = make(map[string]interface{})
	_, err = formatter.FormatWithContext(ctx)
	if err == nil {
		t.Error("FormatWithContext should fail with missing path")
	}

	// Test list archive formatting (Regression test for date formatting)
	// [TEST-FORMAT_LIST] [REQ-OUTPUT_FORMATTING]
	ctx = FormatContext{
		FormatType: FormatTypeList,
		Data: map[string]interface{}{
			"path":         "/test/archive.zip",
			"creationTime": "2023-10-27 10:00:00",
		},
		Options:  FormatOptions{},
		Metadata: make(map[string]string),
	}

	result, err = formatter.FormatWithContext(ctx)
	if err != nil {
		t.Errorf("FormatWithContext failed for list archive: %v", err)
	}

	expected := "/test/archive.zip (created: 2023-10-27 10:00:00)\n"
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

// [CRITICAL] FMT-001: Test extract with context - [ACTION-discovery]
func TestAIFirstFormatterExtractWithContext(t *testing.T) {
	config := NewMockFormatterConfig()
	formatter := NewAIFirstFormatter(config)

	// Test archive data extraction
	ctx := ExtractContext{
		PatternType: PatternTypeArchiveFilename,
		Input:       "test-2024-01-15-14-30-main-abc123-note.tar.gz",
		Options:     ExtractOptions{},
		Metadata:    make(map[string]string),
	}

	result, err := formatter.ExtractWithContext(ctx)
	if err != nil {
		t.Errorf("ExtractWithContext failed: %v", err)
	}

	if result == nil {
		t.Error("ExtractWithContext returned nil result")
	}

	// Verify result is ArchiveData
	if archiveData, ok := result.(AIArchiveData); ok {
		if archiveData.Prefix != "test" {
			t.Errorf("expected prefix 'test', got '%s'", archiveData.Prefix)
		}
		if archiveData.Year != "2024" {
			t.Errorf("expected year '2024', got '%s'", archiveData.Year)
		}
	} else {
		t.Error("ExtractWithContext did not return ArchiveData")
	}
}

// [CRITICAL] FMT-001: Test print with context - [ACTION:format-processing]
func TestAIFirstFormatterPrintWithContext(t *testing.T) {
	config := NewMockFormatterConfig()
	collector := NewOutputCollector()
	formatter := NewAIFirstFormatterWithCollector(config, collector)

	// Test delayed printing
	ctx := PrintContext{
		Message:     "Test message",
		Destination: AIOutputDestinationStdout,
		Type:        AIMessageTypeInfo,
		Options: PrintOptions{
			Delayed: true,
			Flush:   false,
		},
		Metadata: make(map[string]string),
	}

	err := formatter.PrintWithContext(ctx)
	if err != nil {
		t.Errorf("PrintWithContext failed: %v", err)
	}

	messages := formatter.GetCollectedMessages()
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	if messages[0].Content != "Test message" {
		t.Errorf("expected message 'Test message', got '%s'", messages[0].Content)
	}
}

// [CRITICAL] FMT-001: Test core formatter operations - [ACTION:core-functionality]
func TestAICoreFormatterOperations(t *testing.T) {
	config := NewMockFormatterConfig()
	formatter := NewAICoreFormatter(config)

	// Test archive formatting
	result, err := formatter.FormatArchive("/test/archive.tar.gz", FormatTypeCreated)
	if err != nil {
		t.Errorf("FormatArchive failed: %v", err)
	}

	if result == "" {
		t.Error("FormatArchive returned empty result")
	}

	// Test backup formatting
	result, err = formatter.FormatBackup("/test/backup.bak", FormatTypeCreated)
	if err != nil {
		t.Errorf("FormatBackup failed: %v", err)
	}

	if result == "" {
		t.Error("FormatBackup returned empty result")
	}

	// Test config formatting
	result, err = formatter.FormatConfig("test_key", "test_value", "test_source")
	if err != nil {
		t.Errorf("FormatConfig failed: %v", err)
	}

	if result == "" {
		t.Error("FormatConfig returned empty result")
	}

	// Test error formatting
	testErr := errors.New("test error")
	result, err = formatter.FormatError(testErr, ErrorTypeGeneric)
	if err != nil {
		t.Errorf("FormatError failed: %v", err)
	}

	if result == "" {
		t.Error("FormatError returned empty result")
	}
}

// [CRITICAL] FMT-001: Test pattern extractor operations - [ACTION-discovery]
func TestAIPatternExtractorOperations(t *testing.T) {
	config := NewMockFormatterConfig()
	extractor := NewAIPatternExtractor(config)

	// Test archive data extraction
	archiveData, err := extractor.ExtractArchiveData("test-2024-01-15-14-30-main-abc123-note.tar.gz")
	if err != nil {
		t.Errorf("ExtractArchiveData failed: %v", err)
	}

	if archiveData.Prefix != "test" {
		t.Errorf("expected prefix 'test', got '%s'", archiveData.Prefix)
	}

	if archiveData.Year != "2024" {
		t.Errorf("expected year '2024', got '%s'", archiveData.Year)
	}

	// Test backup data extraction
	backupData, err := extractor.ExtractBackupData("testfile-2024-01-15-14-30-note.bak")
	if err != nil {
		t.Errorf("ExtractBackupData failed: %v", err)
	}

	if backupData.Filename != "testfile" {
		t.Errorf("expected filename 'testfile', got '%s'", backupData.Filename)
	}

	// Test config data extraction
	configData, err := extractor.ExtractConfigData("test_key=test_value (test_source)")
	if err != nil {
		t.Errorf("ExtractConfigData failed: %v", err)
	}

	if configData.Name != "test_key" {
		t.Errorf("expected name 'test_key', got '%s'", configData.Name)
	}

	if configData.Value != "test_value" {
		t.Errorf("expected value 'test_value', got '%s'", configData.Value)
	}
}

// [CRITICAL] FMT-001: Test output manager operations - [ACTION:format-processing]
func TestAIOutputManagerOperations(t *testing.T) {
	collector := NewOutputCollector()
	manager := NewAIOutputManagerWithCollector(collector)

	// Test delayed mode
	if !manager.IsDelayedMode() {
		t.Error("manager should be in delayed mode with collector")
	}

	// Test message collection
	message := AIOutputMessage{
		Content:     "Test message",
		Destination: AIOutputDestinationStdout,
		Type:        AIMessageTypeInfo,
		Metadata:    make(map[string]string),
		Timestamp:   time.Now(),
	}

	err := manager.Collect(message)
	if err != nil {
		t.Errorf("Collect failed: %v", err)
	}

	messages := manager.GetCollectedMessages()
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	// Test flush
	err = manager.Flush()
	if err != nil {
		t.Errorf("Flush failed: %v", err)
	}

	// Test clear
	err = manager.Clear()
	if err != nil {
		t.Errorf("Clear failed: %v", err)
	}

	messages = manager.GetCollectedMessages()
	if len(messages) != 0 {
		t.Errorf("expected 0 messages after clear, got %d", len(messages))
	}
}

// [CRITICAL] FMT-001: Test AI-first formatter integration - [ACTION:core-functionality]
func TestAIFirstFormatterIntegration(t *testing.T) {
	config := NewMockFormatterConfig()
	collector := NewOutputCollector()
	formatter := NewAIFirstFormatterWithCollector(config, collector)

	// Test complete workflow
	ctx := FormatContext{
		FormatType: FormatTypeCreated,
		Data: map[string]interface{}{
			"path": "/test/archive.tar.gz",
		},
		Options:  FormatOptions{},
		Metadata: make(map[string]string),
	}

	// Format message
	result, err := formatter.FormatWithContext(ctx)
	if err != nil {
		t.Errorf("FormatWithContext failed: %v", err)
	}

	// Print message
	printCtx := PrintContext{
		Message:     result,
		Destination: AIOutputDestinationStdout,
		Type:        AIMessageTypeInfo,
		Options: PrintOptions{
			Delayed: true,
			Flush:   false,
		},
		Metadata: make(map[string]string),
	}

	err = formatter.PrintWithContext(printCtx)
	if err != nil {
		t.Errorf("PrintWithContext failed: %v", err)
	}

	// Verify message was collected
	messages := formatter.GetCollectedMessages()
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}

	if messages[0].Content != result {
		t.Errorf("collected message does not match formatted result")
	}
}
