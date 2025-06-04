// Pattern extraction and regex-based data extraction for the formatter package.
// Provides functionality to extract structured data from filenames and text
// using named regex groups for template processing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

import (
	"path/filepath"
	"regexp"
	"strings"
)

// ⭐ EXTRACT-003: PatternExtractor component - 🔧 Generic pattern extraction implementation
// DefaultPatternExtractor provides default pattern extraction functionality
type DefaultPatternExtractor struct {
	configProvider ConfigProvider
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔧 Constructor
// NewDefaultPatternExtractor creates a new DefaultPatternExtractor
func NewDefaultPatternExtractor(configProvider ConfigProvider) *DefaultPatternExtractor {
	return &DefaultPatternExtractor{
		configProvider: configProvider,
	}
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Archive filename data extraction
// ExtractArchiveFilenameData extracts data from archive filenames using configured patterns
func (pe *DefaultPatternExtractor) ExtractArchiveFilenameData(filename string) map[string]string {
	pattern := pe.configProvider.GetPattern("archive_filename")
	if pattern == "" {
		// Default pattern for archive filenames
		pattern = `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	}
	return pe.ExtractPatternData(pattern, filename)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Backup filename data extraction
// ExtractBackupFilenameData extracts data from backup filenames using configured patterns
func (pe *DefaultPatternExtractor) ExtractBackupFilenameData(filename string) map[string]string {
	pattern := pe.configProvider.GetPattern("backup_filename")
	if pattern == "" {
		// Default pattern for backup filenames
		pattern = `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	}
	return pe.ExtractPatternData(pattern, filename)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Generic pattern data extraction
// ExtractPatternData extracts named groups from text using a regex pattern
func (pe *DefaultPatternExtractor) ExtractPatternData(pattern, text string) map[string]string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return make(map[string]string)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Additional extraction utilities
// ExtractConfigLineData extracts data from configuration lines
func (pe *DefaultPatternExtractor) ExtractConfigLineData(line string) map[string]string {
	pattern := pe.configProvider.GetPattern("config_line")
	if pattern == "" {
		// Default pattern for config lines (key=value format)
		pattern = `^(?P<key>[^=]+)=(?P<value>.*)$`
	}
	return pe.ExtractPatternData(pattern, strings.TrimSpace(line))
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Timestamp data extraction
// ExtractTimestampData extracts timestamp components
func (pe *DefaultPatternExtractor) ExtractTimestampData(timestamp string) map[string]string {
	pattern := pe.configProvider.GetPattern("timestamp")
	if pattern == "" {
		// Default pattern for timestamps (YYYY-MM-DD_HH-MM-SS format)
		pattern = `^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})_(?P<hour>\d{2})-(?P<minute>\d{2})-(?P<second>\d{2})$`
	}
	return pe.ExtractPatternData(pattern, timestamp)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔧 Utility functions
// GetFilenameFromPath extracts filename from a full path
func GetFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔧 Simple pattern extractor without config
// SimplePatternExtractor provides pattern extraction without configuration dependency
type SimplePatternExtractor struct{}

// ⭐ EXTRACT-003: PatternExtractor component - 🔧 Simple constructor
// NewSimplePatternExtractor creates a SimplePatternExtractor
func NewSimplePatternExtractor() *SimplePatternExtractor {
	return &SimplePatternExtractor{}
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Simple archive filename extraction
// ExtractArchiveFilenameData extracts data using default archive pattern
func (spe *SimplePatternExtractor) ExtractArchiveFilenameData(filename string) map[string]string {
	pattern := `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	return spe.ExtractPatternData(pattern, filename)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Simple backup filename extraction
// ExtractBackupFilenameData extracts data using default backup pattern
func (spe *SimplePatternExtractor) ExtractBackupFilenameData(filename string) map[string]string {
	pattern := `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	return spe.ExtractPatternData(pattern, filename)
}

// ⭐ EXTRACT-003: PatternExtractor component - 🔍 Simple pattern extraction
// ExtractPatternData extracts named groups from text using a regex pattern
func (spe *SimplePatternExtractor) ExtractPatternData(pattern, text string) map[string]string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return make(map[string]string)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i < len(matches) && name != "" {
			result[name] = matches[i]
		}
	}

	return result
}
