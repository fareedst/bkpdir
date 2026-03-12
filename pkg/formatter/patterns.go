// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// DefaultPatternExtractor provides default pattern extraction functionality
type DefaultPatternExtractor struct {
	configProvider ConfigProvider
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// NewDefaultPatternExtractor creates a new DefaultPatternExtractor
func NewDefaultPatternExtractor(configProvider ConfigProvider) *DefaultPatternExtractor {
	return &DefaultPatternExtractor{
		configProvider: configProvider,
	}
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractArchiveFilenameData extracts data from archive filenames using configured patterns
func (pe *DefaultPatternExtractor) ExtractArchiveFilenameData(filename string) map[string]string {
	pattern := pe.configProvider.GetPattern("archive_filename")
	if pattern == "" {
		// Default pattern for archive filenames
		pattern = `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	}
	return pe.ExtractPatternData(pattern, filename)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractBackupFilenameData extracts data from backup filenames using configured patterns
func (pe *DefaultPatternExtractor) ExtractBackupFilenameData(filename string) map[string]string {
	pattern := pe.configProvider.GetPattern("backup_filename")
	if pattern == "" {
		// Default pattern for backup filenames
		pattern = `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	}
	return pe.ExtractPatternData(pattern, filename)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractConfigLineData extracts data from configuration lines
func (pe *DefaultPatternExtractor) ExtractConfigLineData(line string) map[string]string {
	pattern := pe.configProvider.GetPattern("config_line")
	if pattern == "" {
		// Default pattern for config lines (key=value format)
		pattern = `^(?P<key>[^=]+)=(?P<value>.*)$`
	}
	return pe.ExtractPatternData(pattern, strings.TrimSpace(line))
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractTimestampData extracts timestamp components
func (pe *DefaultPatternExtractor) ExtractTimestampData(timestamp string) map[string]string {
	pattern := pe.configProvider.GetPattern("timestamp")
	if pattern == "" {
		// Default pattern for timestamps (YYYY-MM-DD_HH-MM-SS format)
		pattern = `^(?P<year>\d{4})-(?P<month>\d{2})-(?P<day>\d{2})_(?P<hour>\d{2})-(?P<minute>\d{2})-(?P<second>\d{2})$`
	}
	return pe.ExtractPatternData(pattern, timestamp)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// GetFilenameFromPath extracts filename from a full path
func GetFilenameFromPath(path string) string {
	return filepath.Base(path)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// SimplePatternExtractor provides pattern extraction without configuration dependency
type SimplePatternExtractor struct{}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// NewSimplePatternExtractor creates a SimplePatternExtractor
func NewSimplePatternExtractor() *SimplePatternExtractor {
	return &SimplePatternExtractor{}
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractArchiveFilenameData extracts data using default archive pattern
func (spe *SimplePatternExtractor) ExtractArchiveFilenameData(filename string) map[string]string {
	pattern := `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	return spe.ExtractPatternData(pattern, filename)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// ExtractBackupFilenameData extracts data using default backup pattern
func (spe *SimplePatternExtractor) ExtractBackupFilenameData(filename string) map[string]string {
	pattern := `^(?P<name>.*?)_(?P<timestamp>\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2})(?P<suffix>\..*)?$`
	return spe.ExtractPatternData(pattern, filename)
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
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
