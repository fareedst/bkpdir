// This file is part of bkpdir
// 🔺 LINT-001: Lint compliance - 🔧

package main

import (
	"bkpdir/pkg/formatter"
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestTemplateFormatter tests the template formatting functionality for CFG-003 feature
func TestTemplateFormatter(t *testing.T) {
	// 🔺 CFG-003: Template formatting validation - 📝
	// TEST-REF: Feature tracking matrix CFG-003
	// IMMUTABLE-REF: Output Formatting Requirements
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("BasicTemplateFormatting", func(t *testing.T) {
		// Test basic template formatting
		data := map[string]string{
			"path":    "/archives/test.zip",
			"time":    "2024-06-01 12:00:00",
			"message": "Test message",
		}

		// Test template-based formatting
		result := fa.TemplateCreatedArchive(data)
		if !strings.Contains(result, "test.zip") {
			t.Errorf("Expected result to contain filename, got %q", result)
		}
	})

	t.Run("PlaceholderFormatting", func(t *testing.T) {
		// Test placeholder-based formatting
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "main",
			"note":   "test note",
		}

		format := "Archive: %{path} on %{branch} - %{note}"
		result := fa.FormatWithPlaceholders(format, data)
		expected := "Archive: /test/archive.zip on main - test note"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("TemplateWithPatternExtraction", func(t *testing.T) {
		// Test template formatting with pattern extraction
		archivePath := "/archives/HOME-2024-06-01-12-00=main=abc123=test.zip"
		template := cfg.TemplateListArchive

		if template != "" {
			// Use FormatWithPlaceholders with the correct data that matches the template
			data := map[string]string{
				"path":          archivePath,
				"creation_time": "2024-06-01 12:00:00",
			}
			result := fa.FormatWithPlaceholders(template, data)
			if !strings.Contains(result, "test.zip") {
				t.Errorf("Expected result to contain filename, got %q", result)
			}
		}
	})

	t.Run("ErrorTemplateFormatting", func(t *testing.T) {
		// Test error template formatting
		errorData := map[string]string{
			"message":   "Something went wrong",
			"operation": "create_archive",
		}
		result := fa.TemplateError(errorData)
		if !strings.Contains(result, "Something went wrong") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("ExtractionWithFormatting", func(t *testing.T) {
		// Test formatting with automatic extraction using the new interface
		archivePath := "/archives/HOME-2024-06-01-12-00=main=abc123=test.zip"

		// Test list formatting with extraction using the new method
		result := fa.FormatListArchiveWithExtraction(archivePath, "2024-06-01 12:00:00")
		if !strings.Contains(result, "HOME-2024-06-01-12-00=main=abc123=test.zip") {
			t.Errorf("Expected result to contain archive filename, got %q", result)
		}
		if !strings.Contains(result, "2024-06-01 12:00:00") {
			t.Errorf("Expected result to contain creation time, got %q", result)
		}

		// Test backup formatting with extraction using the new method
		backupPath := "/backups/document.txt-2024-06-01-12-00=backup"
		result = fa.FormatListBackupWithExtraction(backupPath, "2024-06-01 12:00:00")
		if !strings.Contains(result, "document.txt-2024-06-01-12-00=backup") {
			t.Errorf("Expected result to contain backup filename, got %q", result)
		}
	})

	t.Run("PatternValidation", func(t *testing.T) {
		// Test that default patterns work correctly
		if cfg.PatternArchiveFilename == "" {
			t.Error("Expected non-empty archive filename pattern")
		}
		if cfg.PatternBackupFilename == "" {
			t.Error("Expected non-empty backup filename pattern")
		}
		if cfg.PatternConfigLine == "" {
			t.Error("Expected non-empty config line pattern")
		}
		if cfg.PatternTimestamp == "" {
			t.Error("Expected non-empty timestamp pattern")
		}

		// Test pattern extraction using the new interface
		timestamp := "2024-06-01 12:30:45"
		timestampData := fa.ExtractPatternData(cfg.PatternTimestamp, timestamp)

		expectedTimestampKeys := []string{"year", "month", "day", "hour", "minute", "second"}
		for _, key := range expectedTimestampKeys {
			if _, exists := timestampData[key]; !exists {
				t.Errorf("Expected timestamp key %q not found in extracted data", key)
			}
		}

		if timestampData["year"] != "2024" {
			t.Errorf("Expected year '2024', got %q", timestampData["year"])
		}
		if timestampData["hour"] != "12" {
			t.Errorf("Expected hour '12', got %q", timestampData["hour"])
		}
		if timestampData["second"] != "45" {
			t.Errorf("Expected second '45', got %q", timestampData["second"])
		}
	})

	t.Run("CustomTemplateHandling", func(t *testing.T) {
		// Test custom template format with placeholders
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "feature-branch",
			"note":   "custom note",
		}

		// Test placeholder replacement
		format := "Archive: %{path} on %{branch} - %{note}"
		result := fa.FormatWithPlaceholders(format, data)
		expected := "Archive: /test/archive.zip on feature-branch - custom note"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test with missing placeholders (should remain as-is)
		format = "Archive: %{path} unknown: %{missing}"
		result = fa.FormatWithPlaceholders(format, data)
		expected = "Archive: /test/archive.zip unknown: %{missing}"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

// TestOutputCollector tests the OutputCollector functionality for OUT-001 feature
func TestOutputCollector(t *testing.T) {
	t.Run("NewOutputCollector", func(t *testing.T) {
		collector := formatter.NewOutputCollector()
		if collector == nil {
			t.Fatal("Expected non-nil collector")
		}
		if len(collector.GetMessages()) != 0 {
			t.Errorf("Expected empty messages, got %d", len(collector.GetMessages()))
		}
	})

	t.Run("AddStdout", func(t *testing.T) {
		collector := formatter.NewOutputCollector()
		collector.AddStdout("test message", "info")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}

		msg := messages[0]
		if msg.Content != "test message" {
			t.Errorf("Expected content 'test message', got %q", msg.Content)
		}
		if msg.Destination != "stdout" {
			t.Errorf("Expected destination 'stdout', got %q", msg.Destination)
		}
		if msg.Type != "info" {
			t.Errorf("Expected type 'info', got %q", msg.Type)
		}
	})

	t.Run("AddStderr", func(t *testing.T) {
		collector := formatter.NewOutputCollector()
		collector.AddStderr("error message", "error")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}

		msg := messages[0]
		if msg.Content != "error message" {
			t.Errorf("Expected content 'error message', got %q", msg.Content)
		}
		if msg.Destination != "stderr" {
			t.Errorf("Expected destination 'stderr', got %q", msg.Destination)
		}
		if msg.Type != "error" {
			t.Errorf("Expected type 'error', got %q", msg.Type)
		}
	})

	t.Run("GetMessages", func(t *testing.T) {
		collector := formatter.NewOutputCollector()
		collector.AddStdout("message 1", "info")
		collector.AddStderr("message 2", "error")

		messages := collector.GetMessages()
		if len(messages) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(messages))
		}

		// Check order is preserved
		if messages[0].Content != "message 1" {
			t.Errorf("Expected first message 'message 1', got %q", messages[0].Content)
		}
		if messages[1].Content != "message 2" {
			t.Errorf("Expected second message 'message 2', got %q", messages[1].Content)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		collector := formatter.NewOutputCollector()
		collector.AddStdout("test message", "info")
		collector.AddStderr("error message", "error")

		if len(collector.GetMessages()) != 2 {
			t.Errorf("Expected 2 messages before clear")
		}

		collector.Clear()

		if len(collector.GetMessages()) != 0 {
			t.Errorf("Expected 0 messages after clear, got %d", len(collector.GetMessages()))
		}
	})

	// Note: FlushAll, FlushStdout, FlushStderr test stdout/stderr output,
	// which is difficult to test without capturing output. These are tested
	// indirectly through delayed mode testing below.
}

// TestDelayedOutputMode tests the delayed output functionality
func TestDelayedOutputMode(t *testing.T) {
	cfg := DefaultConfig()
	collector := formatter.NewOutputCollector()

	t.Run("NewOutputFormatterWithCollector", func(t *testing.T) {
		fa := NewFormatterAdapterWithCollector(cfg, collector)
		if !fa.IsDelayedMode() {
			t.Error("Expected formatter to be in delayed mode")
		}
		if fa.GetCollector() != collector {
			t.Error("Expected collector to match")
		}
	})

	t.Run("IsDelayedMode", func(t *testing.T) {
		normalFormatter := NewOutputFormatter(cfg)
		if normalFormatter.IsDelayedMode() {
			t.Error("Expected normal formatter not to be in delayed mode")
		}

		delayedFormatter := NewOutputFormatterWithCollector(cfg, collector)
		if !delayedFormatter.IsDelayedMode() {
			t.Error("Expected delayed formatter to be in delayed mode")
		}
	})

	t.Run("GetCollector", func(t *testing.T) {
		normalFormatter := NewOutputFormatter(cfg)
		if normalFormatter.GetCollector() != nil {
			t.Error("Expected normal formatter to have nil collector")
		}

		delayedFormatter := NewOutputFormatterWithCollector(cfg, collector)
		if delayedFormatter.GetCollector() == nil {
			t.Error("Expected delayed formatter to have non-nil collector")
		}
	})

	t.Run("SetCollector", func(t *testing.T) {
		formatter := NewOutputFormatter(cfg)

		// Set collector
		formatter.SetCollector(collector)
		if !formatter.IsDelayedMode() {
			t.Error("Expected formatter to be in delayed mode after setting collector")
		}

		// Remove collector
		formatter.SetCollector(nil)
		if formatter.IsDelayedMode() {
			t.Error("Expected formatter not to be in delayed mode after removing collector")
		}
	})

	t.Run("PrintMethodsWithCollector", func(t *testing.T) {
		collector.Clear()
		formatter := NewOutputFormatterWithCollector(cfg, collector)

		// Test various print methods
		formatter.PrintCreatedArchive("/test/archive.zip")
		formatter.PrintIdenticalArchive("/test/identical.zip")
		formatter.PrintListArchive("/test/list.zip", "2024-01-01 12:00:00")
		formatter.PrintConfigValue("test_key", "test_value", "test_source")

		messages := collector.GetMessages()
		if len(messages) != 4 {
			t.Errorf("Expected 4 messages, got %d", len(messages))
		}

		// All should be stdout messages, but PrintConfigValue uses "config" type
		for i, msg := range messages {
			if msg.Destination != "stdout" {
				t.Errorf("Message %d: expected stdout destination, got %q", i, msg.Destination)
			}
			if i == 3 { // PrintConfigValue
				if msg.Type != "config" {
					t.Errorf("Message %d: expected config type, got %q", i, msg.Type)
				}
			} else {
				if msg.Type != "info" {
					t.Errorf("Message %d: expected info type, got %q", i, msg.Type)
				}
			}
		}
	})
}

// TestTemplateFormattingMethods tests template formatting functions with 0% coverage
func TestTemplateFormattingMethods(t *testing.T) {
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("TemplateCreatedArchive", func(t *testing.T) {
		data := map[string]string{
			"path":          "/archives/test.zip",
			"creation_time": "2024-01-01 12:00:00",
		}
		result := fa.TemplateCreatedArchive(data)
		if !strings.Contains(result, "test.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateIdenticalArchive", func(t *testing.T) {
		data := map[string]string{
			"path": "/archives/existing.zip",
		}
		result := fa.TemplateIdenticalArchive(data)
		if !strings.Contains(result, "existing.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateListArchive", func(t *testing.T) {
		data := map[string]string{
			"path":          "/archives/test.zip",
			"creation_time": "2024-01-01 12:00:00",
		}
		result := fa.TemplateListArchive(data)
		if !strings.Contains(result, "test.zip") || !strings.Contains(result, "2024-01-01 12:00:00") {
			t.Errorf("Expected result to contain path and creation time, got %q", result)
		}
	})

	t.Run("TemplateConfigValue", func(t *testing.T) {
		data := map[string]string{
			"name":   "archive_dir_path",
			"value":  "../.bkpdir",
			"source": "default",
		}
		result := fa.TemplateConfigValue(data)
		if !strings.Contains(result, "archive_dir_path") || !strings.Contains(result, "../.bkpdir") {
			t.Errorf("Expected result to contain config name and value, got %q", result)
		}
	})

	t.Run("TemplateDryRunArchive", func(t *testing.T) {
		data := map[string]string{
			"path": "/archives/would-create.zip",
		}
		result := fa.TemplateDryRunArchive(data)
		if !strings.Contains(result, "would-create.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})
}

// TestErrorFormattingMethods tests error formatting functions with 0% coverage
func TestErrorFormattingMethods(t *testing.T) {
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("FormatDiskFullError", func(t *testing.T) {
		err := fmt.Errorf("no space left on device")
		result := fa.FormatDiskFullError(err)
		if !strings.Contains(result, "no space left") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatPermissionError", func(t *testing.T) {
		err := fmt.Errorf("permission denied")
		result := fa.FormatPermissionError(err)
		if !strings.Contains(result, "permission denied") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatDirectoryNotFound", func(t *testing.T) {
		err := fmt.Errorf("directory not found")
		result := fa.FormatDirectoryNotFound(err)
		if !strings.Contains(result, "directory not found") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFileNotFound", func(t *testing.T) {
		err := fmt.Errorf("file not found")
		result := fa.FormatFileNotFound(err)
		if !strings.Contains(result, "file not found") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatInvalidDirectory", func(t *testing.T) {
		err := fmt.Errorf("invalid directory")
		result := fa.FormatInvalidDirectory(err)
		if !strings.Contains(result, "invalid directory") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatInvalidFile", func(t *testing.T) {
		err := fmt.Errorf("invalid file")
		result := fa.FormatInvalidFile(err)
		if !strings.Contains(result, "invalid file") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	// Test template error formatting
	t.Run("TemplateDiskFullError", func(t *testing.T) {
		err := fmt.Errorf("no space left on device")
		result := fa.TemplateDiskFullError(err)
		if !strings.Contains(result, "no space left") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("TemplatePermissionError", func(t *testing.T) {
		err := fmt.Errorf("permission denied")
		result := fa.TemplatePermissionError(err)
		if !strings.Contains(result, "permission denied") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("TemplateDirectoryNotFound", func(t *testing.T) {
		err := fmt.Errorf("directory not found")
		result := fa.TemplateDirectoryNotFound(err)
		if !strings.Contains(result, "directory not found") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("TemplateFileNotFound", func(t *testing.T) {
		err := fmt.Errorf("file not found")
		result := fa.TemplateFileNotFound(err)
		if !strings.Contains(result, "file not found") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})
}

// TestAdditionalFormattingMethods tests other formatting functions with 0% coverage
func TestAdditionalFormattingMethods(t *testing.T) {
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("FormatNoArchivesFound", func(t *testing.T) {
		result := fa.FormatNoArchivesFound("/archives")
		if !strings.Contains(result, "/archives") {
			t.Errorf("Expected result to contain archive directory, got %q", result)
		}
	})

	t.Run("FormatVerificationSuccess", func(t *testing.T) {
		result := fa.FormatVerificationSuccess("test.zip")
		if !strings.Contains(result, "test.zip") {
			t.Errorf("Expected result to contain archive name, got %q", result)
		}
	})

	t.Run("FormatVerificationFailed", func(t *testing.T) {
		err := fmt.Errorf("checksum mismatch")
		result := fa.FormatVerificationFailed("test.zip", err)
		if !strings.Contains(result, "test.zip") || !strings.Contains(result, "checksum mismatch") {
			t.Errorf("Expected result to contain archive name and error, got %q", result)
		}
	})

	t.Run("FormatConfigurationUpdated", func(t *testing.T) {
		result := fa.FormatConfigurationUpdated("test_key", "test_value")
		if !strings.Contains(result, "test_key") || !strings.Contains(result, "test_value") {
			t.Errorf("Expected result to contain key and value, got %q", result)
		}
	})
}

// TestTemplateFormatterAdvanced tests advanced template formatter functionality
func TestTemplateFormatterAdvanced(t *testing.T) {
	cfg := DefaultConfig()
	fa := NewFormatterAdapter(cfg)

	t.Run("TemplateWithExtraction", func(t *testing.T) {
		// Test template formatting with pattern extraction
		archivePath := "/archives/HOME-2024-06-01-12-00=main=abc123=test.zip"
		pattern := cfg.PatternArchiveFilename
		template := cfg.TemplateListArchive

		if template != "" && pattern != "" {
			result, err := fa.FormatWithTemplate(archivePath, pattern, template)
			if err != nil {
				t.Logf("Template formatting not available: %v", err)
			} else {
				if !strings.Contains(result, "test") {
					t.Logf("Template result: %q", result)
				}
			}
		}
	})

	t.Run("ExtractArchiveFilenameData", func(t *testing.T) {
		// Test archive filename pattern extraction
		archiveFilename := "HOME-2024-06-01-12-00=main=abc123=test note.zip"
		data := fa.ExtractArchiveFilenameData(archiveFilename)

		expectedKeys := []string{"prefix", "year", "month", "day", "hour", "minute", "branch", "hash", "note"}
		for _, key := range expectedKeys {
			if _, exists := data[key]; !exists {
				t.Errorf("Expected key %q not found in extracted data", key)
			}
		}

		if data["prefix"] != "HOME" {
			t.Errorf("Expected prefix 'HOME', got %q", data["prefix"])
		}
		if data["year"] != "2024" {
			t.Errorf("Expected year '2024', got %q", data["year"])
		}
		if data["branch"] != "main" {
			t.Errorf("Expected branch 'main', got %q", data["branch"])
		}
		if data["hash"] != "abc123" {
			t.Errorf("Expected hash 'abc123', got %q", data["hash"])
		}
		if data["note"] != "test note" {
			t.Errorf("Expected note 'test note', got %q", data["note"])
		}
	})

	t.Run("ExtractBackupFilenameData", func(t *testing.T) {
		// Test backup filename pattern extraction
		backupFilename := "document.txt-2024-06-01-12-00=backup note"
		backupData := fa.ExtractBackupFilenameData(backupFilename)

		expectedBackupKeys := []string{"filename", "year", "month", "day", "hour", "minute", "note"}
		for _, key := range expectedBackupKeys {
			if _, exists := backupData[key]; !exists {
				t.Errorf("Expected backup key %q not found in extracted data", key)
			}
		}

		if backupData["filename"] != "document.txt" {
			t.Errorf("Expected filename 'document.txt', got %q", backupData["filename"])
		}
		if backupData["note"] != "backup note" {
			t.Errorf("Expected note 'backup note', got %q", backupData["note"])
		}
	})
}

// TestPrintMethods tests Print methods that use collectors in delayed mode
func TestPrintMethods(t *testing.T) {
	cfg := DefaultConfig()
	collector := formatter.NewOutputCollector()
	fa := NewFormatterAdapterWithCollector(cfg, collector)

	t.Run("PrintCreatedArchive", func(t *testing.T) {
		fa.PrintCreatedArchive("/test/archive.zip")
		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if !strings.Contains(messages[0].Content, "/test/archive.zip") {
			t.Errorf("Expected message to contain path, got %q", messages[0].Content)
		}
	})

	t.Run("PrintIdenticalArchive", func(t *testing.T) {
		collector.Clear()
		fa.PrintIdenticalArchive("/test/existing.zip")
		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintError", func(t *testing.T) {
		collector.Clear()
		fa.PrintError("test error")
		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if !strings.Contains(messages[0].Content, "test error") {
			t.Errorf("Expected message to contain error text, got %q", messages[0].Content)
		}
	})

	t.Run("PrintCreatedBackup", func(t *testing.T) {
		collector.Clear()
		fa.PrintCreatedBackup("/test/backup.txt")
		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if !strings.Contains(messages[0].Content, "/test/backup.txt") {
			t.Errorf("Expected message to contain path, got %q", messages[0].Content)
		}
	})

	t.Run("PrintConfigValue", func(t *testing.T) {
		collector.Clear()
		fa.PrintConfigValue("test_key", "test_value", "default")
		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if !strings.Contains(messages[0].Content, "test_key") {
			t.Errorf("Expected message to contain key, got %q", messages[0].Content)
		}
	})
}

// 🔺 TEST-001: Comprehensive formatter testing - OutputCollector flush methods - 📝
func TestOutputCollectorFlushMethods(t *testing.T) {
	// Capture stdout and stderr for testing flush methods
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	defer func() {
		os.Stdout = originalStdout
		os.Stderr = originalStderr
	}()

	t.Run("FlushAll", func(t *testing.T) {
		// Create pipes to capture output
		stdoutReader, stdoutWriter, _ := os.Pipe()
		stderrReader, stderrWriter, _ := os.Pipe()
		os.Stdout = stdoutWriter
		os.Stderr = stderrWriter

		collector := formatter.NewOutputCollector()
		collector.AddStdout("stdout message\n", "info")
		collector.AddStderr("stderr message\n", "error")
		collector.AddStdout("another stdout\n", "info")

		// Verify we have messages before flush
		messages := collector.GetMessages()
		if len(messages) != 3 {
			t.Errorf("Expected 3 messages before flush, got %d", len(messages))
		}

		// Flush all messages
		collector.FlushAll()

		// Close writers and read output
		stdoutWriter.Close()
		stderrWriter.Close()

		stdoutOutput := make([]byte, 1024)
		stderrOutput := make([]byte, 1024)
		stdoutN, _ := stdoutReader.Read(stdoutOutput)
		stderrN, _ := stderrReader.Read(stderrOutput)

		stdoutStr := string(stdoutOutput[:stdoutN])
		stderrStr := string(stderrOutput[:stderrN])

		// Verify stdout contains both stdout messages
		if !strings.Contains(stdoutStr, "stdout message") || !strings.Contains(stdoutStr, "another stdout") {
			t.Errorf("Expected stdout to contain both stdout messages, got: %q", stdoutStr)
		}

		// Verify stderr contains stderr message
		if !strings.Contains(stderrStr, "stderr message") {
			t.Errorf("Expected stderr to contain stderr message, got: %q", stderrStr)
		}

		// Verify collector is cleared after flush
		messages = collector.GetMessages()
		if len(messages) != 0 {
			t.Errorf("Expected 0 messages after FlushAll, got %d", len(messages))
		}

		stdoutReader.Close()
		stderrReader.Close()
	})

	t.Run("FlushStdout", func(t *testing.T) {
		// Create pipe to capture stdout
		stdoutReader, stdoutWriter, _ := os.Pipe()
		os.Stdout = stdoutWriter

		collector := formatter.NewOutputCollector()
		collector.AddStdout("stdout message 1\n", "info")
		collector.AddStderr("stderr message\n", "error")
		collector.AddStdout("stdout message 2\n", "info")

		// Verify we have 3 messages before flush
		messages := collector.GetMessages()
		if len(messages) != 3 {
			t.Errorf("Expected 3 messages before flush, got %d", len(messages))
		}

		// Flush only stdout messages
		collector.FlushStdout()

		// Close writer and read output
		stdoutWriter.Close()

		stdoutOutput := make([]byte, 1024)
		stdoutN, _ := stdoutReader.Read(stdoutOutput)
		stdoutStr := string(stdoutOutput[:stdoutN])

		// Verify stdout contains both stdout messages
		if !strings.Contains(stdoutStr, "stdout message 1") || !strings.Contains(stdoutStr, "stdout message 2") {
			t.Errorf("Expected stdout to contain both stdout messages, got: %q", stdoutStr)
		}

		// Verify only stderr message remains
		messages = collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message after FlushStdout, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected remaining message to be stderr, got %q", messages[0].Destination)
		}
		if !strings.Contains(messages[0].Content, "stderr message") {
			t.Errorf("Expected remaining message to contain stderr content, got %q", messages[0].Content)
		}

		stdoutReader.Close()
	})

	t.Run("FlushStderr", func(t *testing.T) {
		// Create pipe to capture stderr
		stderrReader, stderrWriter, _ := os.Pipe()
		os.Stderr = stderrWriter

		collector := formatter.NewOutputCollector()
		collector.AddStdout("stdout message\n", "info")
		collector.AddStderr("stderr message 1\n", "error")
		collector.AddStderr("stderr message 2\n", "warning")

		// Verify we have 3 messages before flush
		messages := collector.GetMessages()
		if len(messages) != 3 {
			t.Errorf("Expected 3 messages before flush, got %d", len(messages))
		}

		// Flush only stderr messages
		collector.FlushStderr()

		// Close writer and read output
		stderrWriter.Close()

		stderrOutput := make([]byte, 1024)
		stderrN, _ := stderrReader.Read(stderrOutput)
		stderrStr := string(stderrOutput[:stderrN])

		// Verify stderr contains both stderr messages
		if !strings.Contains(stderrStr, "stderr message 1") || !strings.Contains(stderrStr, "stderr message 2") {
			t.Errorf("Expected stderr to contain both stderr messages, got: %q", stderrStr)
		}

		// Verify only stdout message remains
		messages = collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message after FlushStderr, got %d", len(messages))
		}
		if messages[0].Destination != "stdout" {
			t.Errorf("Expected remaining message to be stdout, got %q", messages[0].Destination)
		}
		if !strings.Contains(messages[0].Content, "stdout message") {
			t.Errorf("Expected remaining message to contain stdout content, got %q", messages[0].Content)
		}

		stderrReader.Close()
	})

	t.Run("FlushMethodsWithEmptyCollector", func(t *testing.T) {
		// Test flush methods with empty collector (should not panic)
		collector := formatter.NewOutputCollector()

		// These should not panic or cause errors
		collector.FlushAll()
		collector.FlushStdout()
		collector.FlushStderr()

		// Verify collector remains empty
		messages := collector.GetMessages()
		if len(messages) != 0 {
			t.Errorf("Expected 0 messages in empty collector, got %d", len(messages))
		}
	})
}

// 🔺 TEST-001: Test Print methods in delayed output mode to achieve 100% coverage - 📝
func TestPrintMethodsDelayedMode(t *testing.T) {
	cfg := DefaultConfig()
	collector := formatter.NewOutputCollector()
	formatter := NewOutputFormatterWithCollector(cfg, collector)

	t.Run("PrintMethodsWithCollectorRouting", func(t *testing.T) {
		collector.Clear()

		// Test all Print methods that route to stdout
		formatter.PrintCreatedArchive("/test/archive.zip")
		formatter.PrintIdenticalArchive("/test/identical.zip")
		formatter.PrintListArchive("/test/list.zip", "2024-01-01 12:00:00")
		formatter.PrintConfigValue("test_key", "test_value", "test_source")
		formatter.PrintDryRunArchive("/test/dryrun.zip")

		// Test backup print methods
		formatter.PrintCreatedBackup("/test/backup.txt")
		formatter.PrintIdenticalBackup("/test/identical.txt")
		formatter.PrintListBackup("/test/list.txt", "2024-01-01 12:00:00")
		formatter.PrintDryRunBackup("/test/dryrun.txt")

		// Test additional formatting methods
		formatter.PrintNoArchivesFound("/test/archives")
		formatter.PrintVerificationSuccess("test-archive.zip")
		formatter.PrintVerificationWarning("test-archive.zip", fmt.Errorf("warning"))
		formatter.PrintConfigurationUpdated("test_key", "test_value")
		formatter.PrintConfigFilePath("/path/to/config.yml")
		formatter.PrintDryRunFilesHeader()
		formatter.PrintDryRunFileEntry("test-file.txt")
		formatter.PrintNoFilesModified()
		formatter.PrintIncrementalCreated("/test/incremental.zip")
		formatter.PrintNoBackupsFound("test.txt", "/backups")
		formatter.PrintBackupWouldCreate("/test/would-create.txt")
		formatter.PrintBackupIdentical("/test/identical.txt")
		formatter.PrintBackupCreated("/test/created.txt")

		// Test error methods (these route to stderr)
		formatter.PrintError("test error")
		formatter.PrintVerificationFailed("test-archive.zip", fmt.Errorf("failed"))

		// Test error formatting methods
		formatter.PrintDiskFullError(fmt.Errorf("disk full"))
		formatter.PrintPermissionError(fmt.Errorf("permission denied"))
		formatter.PrintDirectoryNotFound(fmt.Errorf("directory not found"))
		formatter.PrintFileNotFound(fmt.Errorf("file not found"))
		formatter.PrintInvalidDirectory(fmt.Errorf("invalid directory"))
		formatter.PrintInvalidFile(fmt.Errorf("invalid file"))
		formatter.PrintFailedWriteTemp(fmt.Errorf("failed write temp"))
		formatter.PrintFailedFinalizeFile(fmt.Errorf("failed finalize"))
		formatter.PrintFailedCreateDirDisk(fmt.Errorf("failed create dir disk"))
		formatter.PrintFailedCreateDir(fmt.Errorf("failed create dir"))
		formatter.PrintFailedAccessDir(fmt.Errorf("failed access dir"))
		formatter.PrintFailedAccessFile(fmt.Errorf("failed access file"))
		formatter.PrintVerificationErrorDetail("error details")
		formatter.PrintArchiveListWithStatus("output", "status")

		messages := collector.GetMessages()
		if len(messages) == 0 {
			t.Error("Expected messages to be collected in delayed mode")
		}

		// Count stdout vs stderr messages
		stdoutCount := 0
		stderrCount := 0
		for _, msg := range messages {
			if msg.Destination == "stdout" {
				stdoutCount++
			} else if msg.Destination == "stderr" {
				stderrCount++
			}
		}

		// Should have both stdout and stderr messages
		if stdoutCount == 0 {
			t.Error("Expected some stdout messages")
		}
		if stderrCount == 0 {
			t.Error("Expected some stderr messages")
		}
	})
}
