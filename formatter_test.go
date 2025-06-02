// This file is part of bkpdir
// LINT-001: Lint compliance

package main

import (
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
