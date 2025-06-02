// This file is part of bkpdir
// LINT-001: Lint compliance

package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

// TestTemplateFormatter tests the template formatting functionality for CFG-003 feature
func TestTemplateFormatter(t *testing.T) {
	// CFG-003: Template formatting validation
	// TEST-REF: Feature tracking matrix CFG-003
	// IMMUTABLE-REF: Output Formatting Requirements
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	t.Run("PrintfStyleFormatting", func(t *testing.T) {
		// Test basic printf-style formatting
		result := formatter.FormatCreatedArchive("/path/to/archive.zip")
		expected := "Created archive: /path/to/archive.zip\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test identical archive formatting
		result = formatter.FormatIdenticalArchive("/path/to/existing.zip")
		expected = "Directory is identical to existing archive: /path/to/existing.zip\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test list archive formatting
		result = formatter.FormatListArchive("/path/to/archive.zip", "2024-01-01 12:00:00")
		expected = "/path/to/archive.zip (created: 2024-01-01 12:00:00)\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test config value formatting
		result = formatter.FormatConfigValue("archive_dir_path", "../.bkpdir", "default")
		expected = "archive_dir_path: ../.bkpdir (source: default)\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test dry run formatting
		result = formatter.FormatDryRunArchive("/path/to/would-create.zip")
		expected = "Would create archive: /path/to/would-create.zip\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test error formatting
		result = formatter.FormatError("Something went wrong")
		expected = "Error: Something went wrong\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("BackupFormatting", func(t *testing.T) {
		// Test backup creation formatting
		result := formatter.FormatCreatedBackup("/path/to/backup.txt")
		expected := "Created backup: /path/to/backup.txt\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test identical backup formatting
		result = formatter.FormatIdenticalBackup("/path/to/existing-backup.txt")
		expected = "File is identical to existing backup: /path/to/existing-backup.txt\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test list backup formatting
		result = formatter.FormatListBackup("/path/to/backup.txt", "2024-01-01 12:00:00")
		expected = "/path/to/backup.txt (created: 2024-01-01 12:00:00)\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test dry run backup formatting
		result = formatter.FormatDryRunBackup("/path/to/would-backup.txt")
		expected = "Would create backup: /path/to/would-backup.txt\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("PatternExtraction", func(t *testing.T) {
		// Test archive filename pattern extraction
		archiveFilename := "HOME-2024-06-01-12-00=main=abc123=test note.zip"
		data := formatter.ExtractArchiveFilenameData(archiveFilename)

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

		// Test backup filename pattern extraction
		backupFilename := "document.txt-2024-06-01-12-00=backup note"
		backupData := formatter.ExtractBackupFilenameData(backupFilename)

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

	t.Run("TemplateFormatting", func(t *testing.T) {
		templateFormatter := NewTemplateFormatter(cfg)

		// Test template-based archive formatting
		archiveData := map[string]string{
			"path":          "/path/to/archive.zip",
			"creation_time": "2024-01-01 12:00:00",
			"prefix":        "HOME",
			"branch":        "main",
			"hash":          "abc123",
			"note":          "test note",
		}

		result := templateFormatter.TemplateCreatedArchive(archiveData)
		if !strings.Contains(result, "/path/to/archive.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}

		// Test template list formatting
		result = templateFormatter.TemplateListArchive(archiveData)
		if !strings.Contains(result, "/path/to/archive.zip") || !strings.Contains(result, "2024-01-01 12:00:00") {
			t.Errorf("Expected result to contain path and creation time, got %q", result)
		}

		// Test config value template formatting
		configData := map[string]string{
			"name":   "archive_dir_path",
			"value":  "../.bkpdir",
			"source": "default",
		}
		result = templateFormatter.TemplateConfigValue(configData)
		if !strings.Contains(result, "archive_dir_path") || !strings.Contains(result, "../.bkpdir") {
			t.Errorf("Expected result to contain config name and value, got %q", result)
		}

		// Test error template formatting
		errorData := map[string]string{
			"message":   "Something went wrong",
			"operation": "create_archive",
		}
		result = templateFormatter.TemplateError(errorData)
		if !strings.Contains(result, "Something went wrong") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("ExtractionWithFormatting", func(t *testing.T) {
		// Test formatting with automatic extraction
		archivePath := "/archives/HOME-2024-06-01-12-00=main=abc123=test.zip"
		result := formatter.FormatArchiveWithExtraction(archivePath)

		// Should contain the original path
		if !strings.Contains(result, "HOME-2024-06-01-12-00=main=abc123=test.zip") {
			t.Errorf("Expected result to contain archive filename, got %q", result)
		}

		// Test list formatting with extraction
		result = formatter.FormatListArchiveWithExtraction(archivePath, "2024-06-01 12:00:00")
		if !strings.Contains(result, "HOME-2024-06-01-12-00=main=abc123=test.zip") {
			t.Errorf("Expected result to contain archive filename, got %q", result)
		}
		if !strings.Contains(result, "2024-06-01 12:00:00") {
			t.Errorf("Expected result to contain creation time, got %q", result)
		}

		// Test backup formatting with extraction
		backupPath := "/backups/document.txt-2024-06-01-12-00=backup"
		result = formatter.FormatBackupWithExtraction(backupPath)
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

		// Test timestamp pattern extraction
		timestamp := "2024-06-01 12:30:45"
		timestampData := formatter.ExtractTimestampData(timestamp)

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
		templateFormatter := NewTemplateFormatter(cfg)

		// Test custom template format with placeholders
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "feature-branch",
			"note":   "custom note",
		}

		// Test placeholder replacement
		format := "Archive: %{path} on %{branch} - %{note}"
		result := templateFormatter.FormatWithPlaceholders(format, data)
		expected := "Archive: /test/archive.zip on feature-branch - custom note"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		// Test with missing placeholders (should remain as-is)
		format = "Archive: %{path} unknown: %{missing}"
		result = templateFormatter.FormatWithPlaceholders(format, data)
		expected = "Archive: /test/archive.zip unknown: %{missing}"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}

// TestOutputCollector tests the OutputCollector functionality for OUT-001 feature
func TestOutputCollector(t *testing.T) {
	t.Run("NewOutputCollector", func(t *testing.T) {
		collector := NewOutputCollector()
		if collector == nil {
			t.Fatal("Expected non-nil collector")
		}
		if len(collector.GetMessages()) != 0 {
			t.Errorf("Expected empty messages, got %d", len(collector.GetMessages()))
		}
	})

	t.Run("AddStdout", func(t *testing.T) {
		collector := NewOutputCollector()
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
		collector := NewOutputCollector()
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
		collector := NewOutputCollector()
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
		collector := NewOutputCollector()
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
	collector := NewOutputCollector()

	t.Run("NewOutputFormatterWithCollector", func(t *testing.T) {
		formatter := NewOutputFormatterWithCollector(cfg, collector)
		if !formatter.IsDelayedMode() {
			t.Error("Expected formatter to be in delayed mode")
		}
		if formatter.GetCollector() != collector {
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
	formatter := NewOutputFormatter(cfg)
	templateFormatter := NewTemplateFormatter(cfg)

	t.Run("TemplateIdenticalArchive", func(t *testing.T) {
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "main",
			"hash":   "abc123",
		}
		result := templateFormatter.TemplateIdenticalArchive(data)
		if !strings.Contains(result, "/test/archive.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateDryRunArchive", func(t *testing.T) {
		data := map[string]string{
			"path":   "/test/would-create.zip",
			"branch": "feature",
		}
		result := templateFormatter.TemplateDryRunArchive(data)
		if !strings.Contains(result, "/test/would-create.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateCreatedBackup", func(t *testing.T) {
		data := map[string]string{
			"path":     "/test/backup.txt",
			"filename": "backup.txt",
		}
		result := templateFormatter.TemplateCreatedBackup(data)
		if !strings.Contains(result, "/test/backup.txt") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateIdenticalBackup", func(t *testing.T) {
		data := map[string]string{
			"path":     "/test/identical.txt",
			"filename": "identical.txt",
		}
		result := templateFormatter.TemplateIdenticalBackup(data)
		if !strings.Contains(result, "/test/identical.txt") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("TemplateListBackup", func(t *testing.T) {
		data := map[string]string{
			"path":          "/test/backup.txt",
			"creation_time": "2024-01-01 12:00:00",
			"filename":      "backup.txt",
		}
		result := templateFormatter.TemplateListBackup(data)
		if !strings.Contains(result, "/test/backup.txt") || !strings.Contains(result, "2024-01-01 12:00:00") {
			t.Errorf("Expected result to contain path and time, got %q", result)
		}
	})

	t.Run("TemplateDryRunBackup", func(t *testing.T) {
		data := map[string]string{
			"path":     "/test/would-backup.txt",
			"filename": "would-backup.txt",
		}
		result := templateFormatter.TemplateDryRunBackup(data)
		if !strings.Contains(result, "/test/would-backup.txt") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("FormatIdenticalArchiveTemplate", func(t *testing.T) {
		data := map[string]string{
			"path":   "/test/archive.zip",
			"branch": "main",
		}
		result := formatter.FormatIdenticalArchiveTemplate(data)
		if !strings.Contains(result, "/test/archive.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("FormatConfigValueTemplate", func(t *testing.T) {
		data := map[string]string{
			"name":   "test_config",
			"value":  "test_value",
			"source": "test_source",
		}
		result := formatter.FormatConfigValueTemplate(data)
		if !strings.Contains(result, "test_config") || !strings.Contains(result, "test_value") {
			t.Errorf("Expected result to contain name and value, got %q", result)
		}
	})

	t.Run("FormatDryRunArchiveTemplate", func(t *testing.T) {
		data := map[string]string{
			"path":   "/test/dry-run.zip",
			"branch": "feature",
		}
		result := formatter.FormatDryRunArchiveTemplate(data)
		if !strings.Contains(result, "/test/dry-run.zip") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("FormatIdenticalBackupTemplate", func(t *testing.T) {
		data := map[string]string{
			"path":     "/test/identical-backup.txt",
			"filename": "identical-backup.txt",
		}
		result := formatter.FormatIdenticalBackupTemplate(data)
		if !strings.Contains(result, "/test/identical-backup.txt") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})

	t.Run("FormatListBackupTemplate", func(t *testing.T) {
		data := map[string]string{
			"path":          "/test/list-backup.txt",
			"creation_time": "2024-01-01 12:00:00",
			"filename":      "list-backup.txt",
		}
		result := formatter.FormatListBackupTemplate(data)
		if !strings.Contains(result, "/test/list-backup.txt") || !strings.Contains(result, "2024-01-01 12:00:00") {
			t.Errorf("Expected result to contain path and time, got %q", result)
		}
	})

	t.Run("FormatDryRunBackupTemplate", func(t *testing.T) {
		data := map[string]string{
			"path":     "/test/dry-run-backup.txt",
			"filename": "dry-run-backup.txt",
		}
		result := formatter.FormatDryRunBackupTemplate(data)
		if !strings.Contains(result, "/test/dry-run-backup.txt") {
			t.Errorf("Expected result to contain path, got %q", result)
		}
	})
}

// TestErrorFormattingMethods tests error formatting functions with 0% coverage
func TestErrorFormattingMethods(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	testError := fmt.Errorf("test error message")

	t.Run("FormatFileNotFound", func(t *testing.T) {
		result := formatter.FormatFileNotFound(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatInvalidDirectory", func(t *testing.T) {
		result := formatter.FormatInvalidDirectory(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatInvalidFile", func(t *testing.T) {
		result := formatter.FormatInvalidFile(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedWriteTemp", func(t *testing.T) {
		result := formatter.FormatFailedWriteTemp(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedFinalizeFile", func(t *testing.T) {
		result := formatter.FormatFailedFinalizeFile(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedCreateDirDisk", func(t *testing.T) {
		result := formatter.FormatFailedCreateDirDisk(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedCreateDir", func(t *testing.T) {
		result := formatter.FormatFailedCreateDir(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedAccessDir", func(t *testing.T) {
		result := formatter.FormatFailedAccessDir(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	t.Run("FormatFailedAccessFile", func(t *testing.T) {
		result := formatter.FormatFailedAccessFile(testError)
		if !strings.Contains(result, "test error message") {
			t.Errorf("Expected result to contain error message, got %q", result)
		}
	})

	// Template error formatting
	t.Run("TemplateDiskFullError", func(t *testing.T) {
		result := formatter.TemplateDiskFullError(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplatePermissionError", func(t *testing.T) {
		result := formatter.TemplatePermissionError(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateDirectoryNotFound", func(t *testing.T) {
		result := formatter.TemplateDirectoryNotFound(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFileNotFound", func(t *testing.T) {
		result := formatter.TemplateFileNotFound(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateInvalidDirectory", func(t *testing.T) {
		result := formatter.TemplateInvalidDirectory(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateInvalidFile", func(t *testing.T) {
		result := formatter.TemplateInvalidFile(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedWriteTemp", func(t *testing.T) {
		result := formatter.TemplateFailedWriteTemp(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedFinalizeFile", func(t *testing.T) {
		result := formatter.TemplateFailedFinalizeFile(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedCreateDirDisk", func(t *testing.T) {
		result := formatter.TemplateFailedCreateDirDisk(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedCreateDir", func(t *testing.T) {
		result := formatter.TemplateFailedCreateDir(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedAccessDir", func(t *testing.T) {
		result := formatter.TemplateFailedAccessDir(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})

	t.Run("TemplateFailedAccessFile", func(t *testing.T) {
		result := formatter.TemplateFailedAccessFile(testError)
		if !strings.Contains(result, "%v") {
			t.Errorf("Expected result to contain template placeholder %%v, got %q", result)
		}
	})
}

// TestAdditionalFormattingMethods tests other formatting functions with 0% coverage
func TestAdditionalFormattingMethods(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	t.Run("FormatNoArchivesFound", func(t *testing.T) {
		result := formatter.FormatNoArchivesFound("/test/archives")
		if !strings.Contains(result, "/test/archives") {
			t.Errorf("Expected result to contain archive directory, got %q", result)
		}
	})

	t.Run("FormatVerificationFailed", func(t *testing.T) {
		testErr := fmt.Errorf("verification failed")
		result := formatter.FormatVerificationFailed("test-archive.zip", testErr)
		if !strings.Contains(result, "test-archive.zip") {
			t.Errorf("Expected result to contain archive name, got %q", result)
		}
	})

	t.Run("FormatVerificationSuccess", func(t *testing.T) {
		result := formatter.FormatVerificationSuccess("test-archive.zip")
		if !strings.Contains(result, "test-archive.zip") {
			t.Errorf("Expected result to contain archive name, got %q", result)
		}
	})

	t.Run("FormatVerificationWarning", func(t *testing.T) {
		testErr := fmt.Errorf("warning message")
		result := formatter.FormatVerificationWarning("test-archive.zip", testErr)
		if !strings.Contains(result, "test-archive.zip") {
			t.Errorf("Expected result to contain archive name, got %q", result)
		}
	})

	t.Run("FormatNoBackupsFound", func(t *testing.T) {
		result := formatter.FormatNoBackupsFound("test.txt", "/backups")
		if !strings.Contains(result, "test.txt") || !strings.Contains(result, "/backups") {
			t.Errorf("Expected result to contain filename and backup dir, got %q", result)
		}
	})

	t.Run("FormatBackupWouldCreate", func(t *testing.T) {
		result := formatter.FormatBackupWouldCreate("/test/backup.txt")
		if !strings.Contains(result, "/test/backup.txt") {
			t.Errorf("Expected result to contain backup path, got %q", result)
		}
	})

	t.Run("FormatBackupIdentical", func(t *testing.T) {
		result := formatter.FormatBackupIdentical("/test/identical.txt")
		if !strings.Contains(result, "/test/identical.txt") {
			t.Errorf("Expected result to contain backup path, got %q", result)
		}
	})

	t.Run("ExtractConfigLineData", func(t *testing.T) {
		configLine := "archive_dir_path: ../archives (source: config)"
		data := formatter.ExtractConfigLineData(configLine)

		if data["name"] != "archive_dir_path" {
			t.Errorf("Expected name 'archive_dir_path', got %q", data["name"])
		}
		// The extracted value may have trailing spaces - trim for comparison
		value := strings.TrimSpace(data["value"])
		if value != "../archives" {
			t.Errorf("Expected value '../archives', got %q", value)
		}
		if data["source"] != "config" {
			t.Errorf("Expected source 'config', got %q", data["source"])
		}
	})

	t.Run("FormatListBackupWithExtraction", func(t *testing.T) {
		backupPath := "/backups/test.txt-2024-01-01-12-00=backup note"
		result := formatter.FormatListBackupWithExtraction(backupPath, "2024-01-01 12:00:00")
		if !strings.Contains(result, "test.txt-2024-01-01-12-00=backup note") {
			t.Errorf("Expected result to contain backup filename, got %q", result)
		}
		if !strings.Contains(result, "2024-01-01 12:00:00") {
			t.Errorf("Expected result to contain creation time, got %q", result)
		}
	})
}

// TestTemplateFormatterAdvanced tests advanced template formatter functionality
func TestTemplateFormatterAdvanced(t *testing.T) {
	cfg := DefaultConfig()
	templateFormatter := NewTemplateFormatter(cfg)

	t.Run("FormatWithTemplate", func(t *testing.T) {
		input := "HOME-2024-01-01-12-00=main=abc123=test note.zip"
		pattern := `^(?P<prefix>[^-]+)-(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})-(?P<hour>\d{2})-(?P<minute>\d{2})=(?P<branch>[^=]+)=(?P<hash>[^=]+)=(?P<note>.*)\.zip$`
		template := "Archive: {{.prefix}} created on {{.year}}-{{.month}}-{{.day}} from branch {{.branch}}\n"

		result, err := templateFormatter.FormatWithTemplate(input, pattern, template)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expected := "Archive: HOME created on 2024-01-01 from branch main\n"
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("PrintTemplateCreatedArchive", func(t *testing.T) {
		// This test verifies the method runs without error
		// Actual output testing would require capturing stdout
		templateFormatter.PrintTemplateCreatedArchive("/test/archive.zip")
	})

	t.Run("PrintTemplateCreatedBackup", func(t *testing.T) {
		// This test verifies the method runs without error
		// Actual output testing would require capturing stdout
		templateFormatter.PrintTemplateCreatedBackup("/test/backup.txt")
	})

	t.Run("PrintTemplateListBackup", func(t *testing.T) {
		// This test verifies the method runs without error
		// Actual output testing would require capturing stdout
		templateFormatter.PrintTemplateListBackup("/test/backup.txt", "2024-01-01 12:00:00")
	})

	t.Run("PrintTemplateError", func(t *testing.T) {
		// This test verifies the method runs without error
		// Actual output testing would require capturing stderr
		templateFormatter.PrintTemplateError("test error", "test_operation")
	})

	t.Run("extractArchiveData", func(t *testing.T) {
		filename := "HOME-2024-01-01-12-00=main=abc123=test note.zip"
		data := templateFormatter.extractArchiveData(filename)

		if data["prefix"] != "HOME" {
			t.Errorf("Expected prefix 'HOME', got %q", data["prefix"])
		}
		if data["year"] != "2024" {
			t.Errorf("Expected year '2024', got %q", data["year"])
		}
		if data["branch"] != "main" {
			t.Errorf("Expected branch 'main', got %q", data["branch"])
		}
	})

	t.Run("extractBackupData", func(t *testing.T) {
		filename := "document.txt-2024-01-01-12-00=backup note"
		data := templateFormatter.extractBackupData(filename)

		if data["filename"] != "document.txt" {
			t.Errorf("Expected filename 'document.txt', got %q", data["filename"])
		}
		if data["year"] != "2024" {
			t.Errorf("Expected year '2024', got %q", data["year"])
		}
		if data["note"] != "backup note" {
			t.Errorf("Expected note 'backup note', got %q", data["note"])
		}
	})
}

// TestPrintMethods tests Print methods that use collectors in delayed mode
func TestPrintMethods(t *testing.T) {
	cfg := DefaultConfig()
	collector := NewOutputCollector()
	formatter := NewOutputFormatterWithCollector(cfg, collector)

	// Clear collector before each test
	beforeEach := func() {
		collector.Clear()
	}

	t.Run("PrintCreatedBackup", func(t *testing.T) {
		beforeEach()
		formatter.PrintCreatedBackup("/test/backup.txt")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stdout" {
			t.Errorf("Expected stdout destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintIdenticalBackup", func(t *testing.T) {
		beforeEach()
		formatter.PrintIdenticalBackup("/test/identical.txt")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintListBackup", func(t *testing.T) {
		beforeEach()
		formatter.PrintListBackup("/test/backup.txt", "2024-01-01 12:00:00")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintDryRunBackup", func(t *testing.T) {
		beforeEach()
		formatter.PrintDryRunBackup("/test/dry-run.txt")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintNoArchivesFound", func(t *testing.T) {
		beforeEach()
		formatter.PrintNoArchivesFound("/test/archives")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintVerificationFailed", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("verification failed")
		formatter.PrintVerificationFailed("test.zip", testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintVerificationSuccess", func(t *testing.T) {
		beforeEach()
		formatter.PrintVerificationSuccess("test.zip")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintVerificationWarning", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("warning message")
		formatter.PrintVerificationWarning("test.zip", testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		// PrintVerificationWarning uses stdout, not stderr
		if messages[0].Destination != "stdout" {
			t.Errorf("Expected stdout destination, got %q", messages[0].Destination)
		}
		if messages[0].Type != "warning" {
			t.Errorf("Expected warning type, got %q", messages[0].Type)
		}
	})

	t.Run("PrintNoBackupsFound", func(t *testing.T) {
		beforeEach()
		formatter.PrintNoBackupsFound("test.txt", "/backups")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintBackupWouldCreate", func(t *testing.T) {
		beforeEach()
		formatter.PrintBackupWouldCreate("/test/backup.txt")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	t.Run("PrintBackupIdentical", func(t *testing.T) {
		beforeEach()
		formatter.PrintBackupIdentical("/test/identical.txt")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
	})

	// Error print methods
	t.Run("PrintFileNotFound", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("file not found")
		formatter.PrintFileNotFound(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintInvalidDirectory", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("invalid directory")
		formatter.PrintInvalidDirectory(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintInvalidFile", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("invalid file")
		formatter.PrintInvalidFile(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedWriteTemp", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to write temp")
		formatter.PrintFailedWriteTemp(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedFinalizeFile", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to finalize file")
		formatter.PrintFailedFinalizeFile(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedCreateDirDisk", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to create dir - disk full")
		formatter.PrintFailedCreateDirDisk(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedCreateDir", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to create dir")
		formatter.PrintFailedCreateDir(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedAccessDir", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to access dir")
		formatter.PrintFailedAccessDir(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintFailedAccessFile", func(t *testing.T) {
		beforeEach()
		testErr := fmt.Errorf("failed to access file")
		formatter.PrintFailedAccessFile(testErr)

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		if messages[0].Destination != "stderr" {
			t.Errorf("Expected stderr destination, got %q", messages[0].Destination)
		}
	})

	t.Run("PrintVerificationErrorDetail", func(t *testing.T) {
		beforeEach()
		formatter.PrintVerificationErrorDetail("verification error details")

		messages := collector.GetMessages()
		if len(messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(messages))
		}
		// PrintVerificationErrorDetail uses stdout, not stderr
		if messages[0].Destination != "stdout" {
			t.Errorf("Expected stdout destination, got %q", messages[0].Destination)
		}
		if messages[0].Type != "error" {
			t.Errorf("Expected error type, got %q", messages[0].Type)
		}
	})
}

// TestTemplateFormattingWithData tests template formatting methods with various data combinations
func TestTemplateFormattingWithData(t *testing.T) {
	cfg := DefaultConfig()
	formatter := NewOutputFormatter(cfg)

	// Test all template methods that weren't covered in other tests
	t.Run("AllTemplateFormattingMethods", func(t *testing.T) {
		testData := map[string]string{
			"archive_dir":   "/test/archives",
			"archive_name":  "test-archive.zip",
			"filename":      "test.txt",
			"backup_dir":    "/test/backups",
			"path":          "/test/path",
			"creation_time": "2024-01-01 12:00:00",
			"key":           "test_key",
			"value":         "test_value",
			"message":       "test message",
			"error":         "test error",
		}

		// Test all template formatting methods that had 0% coverage
		templateMethods := []struct {
			name   string
			method func(map[string]string) string
		}{
			{"FormatNoArchivesFoundTemplate", formatter.FormatNoArchivesFoundTemplate},
			{"FormatVerificationFailedTemplate", formatter.FormatVerificationFailedTemplate},
			{"FormatVerificationSuccessTemplate", formatter.FormatVerificationSuccessTemplate},
			{"FormatVerificationWarningTemplate", formatter.FormatVerificationWarningTemplate},
			{"FormatConfigurationUpdatedTemplate", formatter.FormatConfigurationUpdatedTemplate},
			{"FormatConfigFilePathTemplate", formatter.FormatConfigFilePathTemplate},
			{"FormatDryRunFilesHeaderTemplate", formatter.FormatDryRunFilesHeaderTemplate},
			{"FormatDryRunFileEntryTemplate", formatter.FormatDryRunFileEntryTemplate},
			{"FormatNoFilesModifiedTemplate", formatter.FormatNoFilesModifiedTemplate},
			{"FormatIncrementalCreatedTemplate", formatter.FormatIncrementalCreatedTemplate},
			{"FormatNoBackupsFoundTemplate", formatter.FormatNoBackupsFoundTemplate},
			{"FormatBackupWouldCreateTemplate", formatter.FormatBackupWouldCreateTemplate},
			{"FormatBackupIdenticalTemplate", formatter.FormatBackupIdenticalTemplate},
			{"FormatBackupCreatedTemplate", formatter.FormatBackupCreatedTemplate},
		}

		for _, tm := range templateMethods {
			t.Run(tm.name, func(t *testing.T) {
				result := tm.method(testData)
				// Basic validation - result should not be empty and should be a string
				if result == "" {
					t.Errorf("Expected non-empty result from %s", tm.name)
				}
			})
		}
	})
}

// TEST-001: Comprehensive formatter testing - OutputCollector flush methods
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

		collector := NewOutputCollector()
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

		collector := NewOutputCollector()
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

		collector := NewOutputCollector()
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
		collector := NewOutputCollector()

		// These should not panic or cause errors
		collector.FlushAll()
		collector.FlushStdout()
		collector.FlushStderr()

		// Verify collector remains empty
		messages := collector.GetMessages()
		if len(messages) != 0 {
			t.Errorf("Expected 0 messages after flushing empty collector, got %d", len(messages))
		}
	})
}

// TEST-001: Test Print methods in delayed output mode to achieve 100% coverage
func TestPrintMethodsDelayedMode(t *testing.T) {
	cfg := DefaultConfig()
	collector := NewOutputCollector()
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
