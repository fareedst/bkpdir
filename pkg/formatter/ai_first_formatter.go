// AI-First Formatter Implementation
// Provides clear component separation and interface standardization for optimal AI assistant comprehension.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:core-functionality]
package formatter

import (
	"fmt"
	"time"
)

// [CRITICAL] FMT-001: AI-first formatter implementation - [ACTION:core-functionality]
type AIFirstFormatterImpl struct {
	config           FormatterConfig
	coreFormatter    CoreFormatter
	patternExtractor AIPatternExtractor
	outputManager    AIOutputManager
}

// [CRITICAL] FMT-001: AI-first formatter constructor - [ACTION:core-functionality]
func NewAIFirstFormatter(config FormatterConfig) *AIFirstFormatterImpl {
	return &AIFirstFormatterImpl{
		config:           config,
		coreFormatter:    NewAICoreFormatter(config),
		patternExtractor: NewAIPatternExtractor(config),
		outputManager:    NewAIOutputManager(),
	}
}

// [CRITICAL] FMT-001: AI-first formatter with collector - [ACTION:core-functionality]
func NewAIFirstFormatterWithCollector(config FormatterConfig, collector *OutputCollector) *AIFirstFormatterImpl {
	return &AIFirstFormatterImpl{
		config:           config,
		coreFormatter:    NewAICoreFormatter(config),
		patternExtractor: NewAIPatternExtractor(config),
		outputManager:    NewAIOutputManagerWithCollector(collector),
	}
}

// [CRITICAL] FMT-001: Configuration management - [ACTION:configure-modify]
func (f *AIFirstFormatterImpl) SetConfig(config FormatterConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	f.config = config
	return nil
}

func (f *AIFirstFormatterImpl) GetConfig() FormatterConfig {
	return f.config
}

// [CRITICAL] FMT-001: AI-friendly operations - [ACTION:core-functionality]
func (f *AIFirstFormatterImpl) FormatWithContext(ctx FormatContext) (string, error) {
	switch ctx.FormatType {
	case FormatTypeCreated:
		if path, ok := ctx.Data["path"].(string); ok {
			return f.FormatArchive(path, FormatTypeCreated)
		}
		return "", fmt.Errorf("missing path in format context")
	case FormatTypeIdentical:
		if path, ok := ctx.Data["path"].(string); ok {
			return f.FormatArchive(path, FormatTypeIdentical)
		}
		return "", fmt.Errorf("missing path in format context")
	case FormatTypeList:
		// [IMPL-DUAL_FORMATTING] [REQ-OUTPUT_FORMATTING] Delegate to core formatter for correct list formatting
		return f.coreFormatter.FormatWithContext(ctx)
	case FormatTypeError:
		if err, ok := ctx.Data["error"].(error); ok {
			if errorType, ok := ctx.Data["errorType"].(ErrorType); ok {
				return f.FormatError(err, errorType)
			}
		}
		return "", fmt.Errorf("missing error or errorType in format context")
	default:
		return "", fmt.Errorf("unsupported format type: %s", ctx.FormatType)
	}
}

func (f *AIFirstFormatterImpl) ExtractWithContext(ctx ExtractContext) (interface{}, error) {
	switch ctx.PatternType {
	case PatternTypeArchiveFilename:
		return f.ExtractArchiveData(ctx.Input)
	case PatternTypeBackupFilename:
		return f.ExtractBackupData(ctx.Input)
	case PatternTypeConfigLine:
		return f.ExtractConfigData(ctx.Input)
	case PatternTypeTimestamp:
		return f.ExtractTimestampData(ctx.Input)
	default:
		return f.ExtractPattern(ctx.Input, ctx.Input)
	}
}

func (f *AIFirstFormatterImpl) PrintWithContext(ctx PrintContext) error {
	if ctx.Options.Delayed {
		message := AIOutputMessage{
			Content:     ctx.Message,
			Destination: ctx.Destination,
			Type:        ctx.Type,
			Metadata:    ctx.Metadata,
			Timestamp:   time.Now(),
		}
		return f.Collect(message)
	} else {
		if ctx.Destination == AIOutputDestinationStderr {
			return f.PrintError(ctx.Message)
		} else {
			return f.Print(ctx.Message)
		}
	}
}

// [CRITICAL] FMT-001: Delegate to core formatter - [ACTION:core-functionality]
func (f *AIFirstFormatterImpl) FormatArchive(path string, formatType FormatType) (string, error) {
	return f.coreFormatter.FormatArchive(path, formatType)
}

func (f *AIFirstFormatterImpl) FormatBackup(path string, formatType FormatType) (string, error) {
	return f.coreFormatter.FormatBackup(path, formatType)
}

func (f *AIFirstFormatterImpl) FormatConfig(name, value, source string) (string, error) {
	return f.coreFormatter.FormatConfig(name, value, source)
}

func (f *AIFirstFormatterImpl) FormatError(err error, errorType ErrorType) (string, error) {
	return f.coreFormatter.FormatError(err, errorType)
}

func (f *AIFirstFormatterImpl) FormatWithTemplate(template string, data map[string]interface{}) (string, error) {
	return f.coreFormatter.FormatWithTemplate(template, data)
}

func (f *AIFirstFormatterImpl) FormatWithPlaceholders(format string, data map[string]string) (string, error) {
	return f.coreFormatter.FormatWithPlaceholders(format, data)
}

// [CRITICAL] FMT-001: Delegate to pattern extractor - [ACTION-discovery]
func (f *AIFirstFormatterImpl) ExtractArchiveData(filename string) (AIArchiveData, error) {
	return f.patternExtractor.ExtractArchiveData(filename)
}

func (f *AIFirstFormatterImpl) ExtractBackupData(filename string) (AIBackupData, error) {
	return f.patternExtractor.ExtractBackupData(filename)
}

func (f *AIFirstFormatterImpl) ExtractConfigData(line string) (AIConfigData, error) {
	return f.patternExtractor.ExtractConfigData(line)
}

func (f *AIFirstFormatterImpl) ExtractTimestampData(timestamp string) (AITimestampData, error) {
	return f.patternExtractor.ExtractTimestampData(timestamp)
}

func (f *AIFirstFormatterImpl) ExtractPattern(pattern, text string) (map[string]string, error) {
	return f.patternExtractor.ExtractPattern(pattern, text)
}

// [CRITICAL] FMT-001: Delegate to output manager - [ACTION:format-processing]
func (f *AIFirstFormatterImpl) Print(message string) error {
	return f.outputManager.Print(message)
}

func (f *AIFirstFormatterImpl) PrintError(message string) error {
	return f.outputManager.PrintError(message)
}

func (f *AIFirstFormatterImpl) Collect(message AIOutputMessage) error {
	return f.outputManager.Collect(message)
}

func (f *AIFirstFormatterImpl) Flush() error {
	return f.outputManager.Flush()
}

func (f *AIFirstFormatterImpl) FlushStdout() error {
	return f.outputManager.FlushStdout()
}

func (f *AIFirstFormatterImpl) FlushStderr() error {
	return f.outputManager.FlushStderr()
}

func (f *AIFirstFormatterImpl) Clear() error {
	return f.outputManager.Clear()
}

func (f *AIFirstFormatterImpl) IsDelayedMode() bool {
	return f.outputManager.IsDelayedMode()
}

func (f *AIFirstFormatterImpl) SetDelayedMode(enabled bool) error {
	return f.outputManager.SetDelayedMode(enabled)
}

func (f *AIFirstFormatterImpl) GetCollectedMessages() []AIOutputMessage {
	return f.outputManager.GetCollectedMessages()
}
