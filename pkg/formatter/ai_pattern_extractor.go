// AI-First Pattern Extractor Implementation
// Provides structured data extraction for optimal AI assistant comprehension.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:discovery]
package formatter

import (
	"fmt"
	"regexp"
	"strings"
)

// [CRITICAL] FMT-001: AI-first pattern extractor implementation - [ACTION:discovery]
type AIPatternExtractorImpl struct {
	config FormatterConfig
}

// [CRITICAL] FMT-001: AI-first pattern extractor constructor - [ACTION:discovery]
func NewAIPatternExtractor(config FormatterConfig) *AIPatternExtractorImpl {
	return &AIPatternExtractorImpl{
		config: config,
	}
}

// [CRITICAL] FMT-001: Archive data extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractArchiveData(filename string) (AIArchiveData, error) {
	pattern, err := e.config.GetPattern(PatternTypeArchiveFilename)
	if err != nil {
		return AIArchiveData{}, fmt.Errorf("failed to get archive pattern: %w", err)
	}

	if pattern == "" {
		// Default archive pattern: prefix-YYYY-MM-DD-HH-MM-branch-hash-note.tar.gz
		pattern = `^(.+)-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(.+)-([a-f0-9]+)(?:-(.+))?\.tar\.gz$`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return AIArchiveData{}, fmt.Errorf("failed to compile archive pattern: %w", err)
	}

	matches := re.FindStringSubmatch(filename)
	if matches == nil {
		return AIArchiveData{}, fmt.Errorf("filename does not match archive pattern: %s", filename)
	}

	// Expected groups: [full, prefix, year, month, day, hour, minute, branch, hash, note]
	if len(matches) < 9 {
		return AIArchiveData{}, fmt.Errorf("insufficient matches for archive pattern: %d", len(matches))
	}

	return AIArchiveData{
		Prefix: matches[1],
		Year:   matches[2],
		Month:  matches[3],
		Day:    matches[4],
		Hour:   matches[5],
		Minute: matches[6],
		Branch: matches[7],
		Hash:   matches[8],
		Note:   matches[9], // Optional
	}, nil
}

// [CRITICAL] FMT-001: Backup data extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractBackupData(filename string) (AIBackupData, error) {
	pattern, err := e.config.GetPattern(PatternTypeBackupFilename)
	if err != nil {
		return AIBackupData{}, fmt.Errorf("failed to get backup pattern: %w", err)
	}

	if pattern == "" {
		// Default backup pattern: filename-YYYY-MM-DD-HH-MM-note.bak
		pattern = `^(.+)-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})(?:-(.+))?\.bak$`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return AIBackupData{}, fmt.Errorf("failed to compile backup pattern: %w", err)
	}

	matches := re.FindStringSubmatch(filename)
	if matches == nil {
		return AIBackupData{}, fmt.Errorf("filename does not match backup pattern: %s", filename)
	}

	// Expected groups: [full, filename, year, month, day, hour, minute, note]
	if len(matches) < 7 {
		return AIBackupData{}, fmt.Errorf("insufficient matches for backup pattern: %d", len(matches))
	}

	return AIBackupData{
		Filename: matches[1],
		Year:     matches[2],
		Month:    matches[3],
		Day:      matches[4],
		Hour:     matches[5],
		Minute:   matches[6],
		Note:     matches[7], // Optional
	}, nil
}

// [CRITICAL] FMT-001: Config data extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractConfigData(line string) (AIConfigData, error) {
	pattern, err := e.config.GetPattern(PatternTypeConfigLine)
	if err != nil {
		return AIConfigData{}, fmt.Errorf("failed to get config pattern: %w", err)
	}

	if pattern == "" {
		// Default config pattern: name=value (source)
		pattern = `^([^=]+)=(.+?)(?:\s+\((.+)\))?$`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return AIConfigData{}, fmt.Errorf("failed to compile config pattern: %w", err)
	}

	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return AIConfigData{}, fmt.Errorf("line does not match config pattern: %s", line)
	}

	// Expected groups: [full, name, value, source]
	if len(matches) < 3 {
		return AIConfigData{}, fmt.Errorf("insufficient matches for config pattern: %d", len(matches))
	}

	return AIConfigData{
		Name:   strings.TrimSpace(matches[1]),
		Value:  strings.TrimSpace(matches[2]),
		Source: strings.TrimSpace(matches[3]), // Optional
	}, nil
}

// [CRITICAL] FMT-001: Timestamp data extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractTimestampData(timestamp string) (AITimestampData, error) {
	pattern, err := e.config.GetPattern(PatternTypeTimestamp)
	if err != nil {
		return AITimestampData{}, fmt.Errorf("failed to get timestamp pattern: %w", err)
	}

	if pattern == "" {
		// Default timestamp pattern: YYYY-MM-DD-HH-MM-SS
		pattern = `^(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})$`
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return AITimestampData{}, fmt.Errorf("failed to compile timestamp pattern: %w", err)
	}

	matches := re.FindStringSubmatch(timestamp)
	if matches == nil {
		return AITimestampData{}, fmt.Errorf("timestamp does not match pattern: %s", timestamp)
	}

	// Expected groups: [full, year, month, day, hour, minute, second]
	if len(matches) < 7 {
		return AITimestampData{}, fmt.Errorf("insufficient matches for timestamp pattern: %d", len(matches))
	}

	return AITimestampData{
		Year:   matches[1],
		Month:  matches[2],
		Day:    matches[3],
		Hour:   matches[4],
		Minute: matches[5],
		Second: matches[6],
	}, nil
}

// [CRITICAL] FMT-001: Generic pattern extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractPattern(pattern, text string) (map[string]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile pattern: %w", err)
	}

	matches := re.FindStringSubmatch(text)
	if matches == nil {
		return nil, fmt.Errorf("text does not match pattern: %s", text)
	}

	result := make(map[string]string)
	for i, match := range matches {
		if i == 0 {
			continue // Skip the full match
		}
		result[fmt.Sprintf("group_%d", i)] = match
	}

	return result, nil
}

// [CRITICAL] FMT-001: AI-friendly context extraction - [ACTION:discovery]
func (e *AIPatternExtractorImpl) ExtractWithContext(ctx ExtractContext) (interface{}, error) {
	switch ctx.PatternType {
	case PatternTypeArchiveFilename:
		return e.ExtractArchiveData(ctx.Input)
	case PatternTypeBackupFilename:
		return e.ExtractBackupData(ctx.Input)
	case PatternTypeConfigLine:
		return e.ExtractConfigData(ctx.Input)
	case PatternTypeTimestamp:
		return e.ExtractTimestampData(ctx.Input)
	default:
		return e.ExtractPattern(ctx.Input, ctx.Input)
	}
}
