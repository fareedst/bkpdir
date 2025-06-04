// Tests for the pkg/errors package to validate extracted error handling functionality.
// These tests ensure the extracted error handling components work correctly
// and maintain backward compatibility with the original functionality.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package errors

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
)

// ⭐ EXTRACT-002: Error interface testing - 🧪 ApplicationError functionality
func TestApplicationError(t *testing.T) {
	// Test basic error creation
	err := NewApplicationError("test error", 42)
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
	if err.GetStatusCode() != 42 {
		t.Errorf("Expected status code 42, got %d", err.GetStatusCode())
	}

	// Test error with cause
	cause := errors.New("underlying error")
	errWithCause := NewApplicationErrorWithCause("wrapper error", 43, cause)
	if !strings.Contains(errWithCause.Error(), "wrapper error") {
		t.Errorf("Error message should contain wrapper error")
	}
	if !strings.Contains(errWithCause.Error(), "underlying error") {
		t.Errorf("Error message should contain underlying error")
	}
	if errWithCause.Unwrap() != cause {
		t.Errorf("Unwrap should return the original cause")
	}

	// Test error with full context
	contextErr := NewApplicationErrorWithContext("context error", 44, "test_operation", "/test/path", cause)
	if contextErr.GetOperation() != "test_operation" {
		t.Errorf("Expected operation 'test_operation', got '%s'", contextErr.GetOperation())
	}
	if contextErr.GetPath() != "/test/path" {
		t.Errorf("Expected path '/test/path', got '%s'", contextErr.GetPath())
	}
}

// ⭐ EXTRACT-002: Error classification testing - 🧪 Error detection functions
func TestErrorClassification(t *testing.T) {
	// Test disk full error detection
	diskFullErrors := []error{
		errors.New("no space left on device"),
		errors.New("disk full"),
		errors.New("insufficient disk space"),
		errors.New("quota exceeded"),
	}
	for _, err := range diskFullErrors {
		if !IsDiskFullError(err) {
			t.Errorf("Should detect disk full error: %s", err.Error())
		}
	}

	// Test permission error detection
	permissionErrors := []error{
		errors.New("permission denied"),
		errors.New("access denied"),
		errors.New("operation not permitted"),
		errors.New("insufficient privileges"),
	}
	for _, err := range permissionErrors {
		if !IsPermissionError(err) {
			t.Errorf("Should detect permission error: %s", err.Error())
		}
	}

	// Test directory not found error detection
	notFoundErrors := []error{
		errors.New("no such file or directory"),
		errors.New("directory not found"),
		errors.New("path does not exist"),
	}
	for _, err := range notFoundErrors {
		if !IsDirectoryNotFoundError(err) {
			t.Errorf("Should detect directory not found error: %s", err.Error())
		}
	}

	// Test that normal errors are not incorrectly classified
	normalErr := errors.New("some other error")
	if IsDiskFullError(normalErr) || IsPermissionError(normalErr) || IsDirectoryNotFoundError(normalErr) {
		t.Errorf("Normal error should not be classified as special error type")
	}
}

// ⭐ EXTRACT-002: Error classification framework testing - 🧪 Classifier functionality
func TestDefaultErrorClassifier(t *testing.T) {
	classifier := NewDefaultErrorClassifier()

	// Test disk space error classification
	diskErr := errors.New("no space left on device")
	if classifier.ClassifyError(diskErr) != ErrorCategoryDiskSpace {
		t.Errorf("Should classify disk space error correctly")
	}
	if classifier.GetSeverity(diskErr) != ErrorSeverityCritical {
		t.Errorf("Disk space errors should be critical severity")
	}
	if !classifier.IsRecoverable(diskErr) {
		t.Errorf("Disk space errors should be recoverable")
	}

	// Test permission error classification
	permErr := errors.New("permission denied")
	if classifier.ClassifyError(permErr) != ErrorCategoryPermission {
		t.Errorf("Should classify permission error correctly")
	}
	if classifier.GetSeverity(permErr) != ErrorSeverityError {
		t.Errorf("Permission errors should be error severity")
	}
	if classifier.IsRecoverable(permErr) {
		t.Errorf("Permission errors should not be recoverable")
	}

	// Test unknown error classification
	unknownErr := errors.New("some unknown error")
	if classifier.ClassifyError(unknownErr) != ErrorCategoryUnknown {
		t.Errorf("Should classify unknown error correctly")
	}
}

// ⭐ EXTRACT-002: Error context testing - 🧪 Error context functionality
func TestErrorContext(t *testing.T) {
	ctx := context.Background()
	errorCtx := NewErrorContext("test_operation", "/test/path", ctx)

	if errorCtx.Operation != "test_operation" {
		t.Errorf("Expected operation 'test_operation', got '%s'", errorCtx.Operation)
	}
	if errorCtx.Path != "/test/path" {
		t.Errorf("Expected path '/test/path', got '%s'", errorCtx.Path)
	}

	// Test metadata operations
	errorCtx.WithMetadata("key1", "value1")
	if value, exists := errorCtx.GetMetadata("key1"); !exists || value != "value1" {
		t.Errorf("Metadata should be stored and retrieved correctly")
	}

	if _, exists := errorCtx.GetMetadata("nonexistent"); exists {
		t.Errorf("Nonexistent metadata should return false")
	}
}

// ⭐ EXTRACT-002: Path validation testing - 🧪 Validation functions
func TestPathValidation(t *testing.T) {
	// Mock configuration for testing
	mockConfig := &mockErrorConfig{}

	// Test empty path validation
	err := ValidateDirectoryPath("", mockConfig)
	if err == nil {
		t.Errorf("Empty directory path should return error")
	}

	err = ValidateFilePath("", mockConfig)
	if err == nil {
		t.Errorf("Empty file path should return error")
	}

	// Test nonexistent path validation
	err = ValidateDirectoryPath("/nonexistent/path", mockConfig)
	if err == nil {
		t.Errorf("Nonexistent directory path should return error")
	}

	err = ValidateFilePath("/nonexistent/file", mockConfig)
	if err == nil {
		t.Errorf("Nonexistent file path should return error")
	}
}

// ⭐ EXTRACT-002: Error handler testing - 🧪 Handler functions
func TestHandleError(t *testing.T) {
	mockConfig := &mockErrorConfig{}
	mockFormatter := &mockErrorFormatter{}

	// Test handling nil error
	statusCode := HandleError(nil, mockConfig, mockFormatter)
	if statusCode != 0 {
		t.Errorf("Nil error should return status code 0")
	}

	// Test handling ApplicationError
	appErr := NewApplicationError("test error", 42)
	statusCode = HandleError(appErr, mockConfig, mockFormatter)
	if statusCode != 42 {
		t.Errorf("ApplicationError should return its status code")
	}

	// Test handling disk full error
	diskErr := errors.New("no space left on device")
	statusCode = HandleError(diskErr, mockConfig, mockFormatter)
	if statusCode != 30 { // Mock config returns 30 for disk_full
		t.Errorf("Disk full error should return configured status code")
	}
}

// Mock implementations for testing

// ⭐ EXTRACT-002: Mock configuration for testing - 🧪 Test utilities
type mockErrorConfig struct{}

func (m *mockErrorConfig) GetStatusCodes() map[string]int {
	return map[string]int{
		"disk_full":           30,
		"permission_denied":   22,
		"directory_not_found": 20,
		"file_not_found":      20,
		"network_error":       1,
	}
}

func (m *mockErrorConfig) GetErrorFormatStrings() map[string]string {
	return map[string]string{
		"error": "Error: %s\n",
	}
}

func (m *mockErrorConfig) GetDirectoryPermissions() os.FileMode {
	return 0755
}

func (m *mockErrorConfig) GetFilePermissions() os.FileMode {
	return 0644
}

// ⭐ EXTRACT-002: Mock formatter for testing - 🧪 Test utilities
type mockErrorFormatter struct {
	lastMessage string
}

func (m *mockErrorFormatter) FormatError(message string) string {
	m.lastMessage = message
	return message
}

func (m *mockErrorFormatter) PrintError(message string) {
	m.lastMessage = message
}

func (m *mockErrorFormatter) FormatDiskFullError(err error) string {
	return "Disk full: " + err.Error()
}

func (m *mockErrorFormatter) FormatPermissionError(err error) string {
	return "Permission denied: " + err.Error()
}

func (m *mockErrorFormatter) FormatDirectoryNotFound(err error) string {
	return "Directory not found: " + err.Error()
}

func (m *mockErrorFormatter) FormatFileNotFound(err error) string {
	return "File not found: " + err.Error()
}

func (m *mockErrorFormatter) PrintDiskFullError(err error) {
	m.lastMessage = "Disk full: " + err.Error()
}

func (m *mockErrorFormatter) PrintPermissionError(err error) {
	m.lastMessage = "Permission denied: " + err.Error()
}

func (m *mockErrorFormatter) PrintDirectoryNotFound(err error) {
	m.lastMessage = "Directory not found: " + err.Error()
}
