// AI-First Output Manager Implementation
// Provides output handling with delayed output support for optimal AI assistant comprehension.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [CRITICAL] FMT-001: AI-first formatter refactoring - [ACTION:format-processing]
package formatter

import (
	"fmt"
	"os"
	"time"
)

// [CRITICAL] FMT-001: AI-first output manager implementation - [ACTION:format-processing]
type AIOutputManagerImpl struct {
	collector *OutputCollector
	delayed   bool
}

// [CRITICAL] FMT-001: AI-first output manager constructor - [ACTION:format-processing]
func NewAIOutputManager() *AIOutputManagerImpl {
	return &AIOutputManagerImpl{
		collector: nil,
		delayed:   false,
	}
}

// [CRITICAL] FMT-001: AI-first output manager with collector - [ACTION:format-processing]
func NewAIOutputManagerWithCollector(collector *OutputCollector) *AIOutputManagerImpl {
	return &AIOutputManagerImpl{
		collector: collector,
		delayed:   collector != nil,
	}
}

// [CRITICAL] FMT-001: Direct output operations - [ACTION:format-processing]
func (m *AIOutputManagerImpl) Print(message string) error {
	if m.IsDelayedMode() {
		aiMessage := AIOutputMessage{
			Content:     message,
			Destination: AIOutputDestinationStdout,
			Type:        AIMessageTypeInfo,
			Timestamp:   time.Now(),
		}
		return m.Collect(aiMessage)
	} else {
		_, err := fmt.Print(message)
		return err
	}
}

func (m *AIOutputManagerImpl) PrintError(message string) error {
	if m.IsDelayedMode() {
		aiMessage := AIOutputMessage{
			Content:     message,
			Destination: AIOutputDestinationStderr,
			Type:        AIMessageTypeError,
			Timestamp:   time.Now(),
		}
		return m.Collect(aiMessage)
	} else {
		_, err := fmt.Fprint(os.Stderr, message)
		return err
	}
}

// [CRITICAL] FMT-001: Delayed output operations - [ACTION:format-processing]
func (m *AIOutputManagerImpl) Collect(message AIOutputMessage) error {
	if m.collector == nil {
		return fmt.Errorf("no collector available for delayed output")
	}

	// Convert AI message to legacy format for compatibility
	legacyMessage := OutputMessage{
		Content:     message.Content,
		Destination: string(message.Destination),
		Type:        string(message.Type),
	}

	if message.Destination == AIOutputDestinationStderr {
		m.collector.AddStderr(legacyMessage.Content, legacyMessage.Type)
	} else {
		m.collector.AddStdout(legacyMessage.Content, legacyMessage.Type)
	}

	return nil
}

func (m *AIOutputManagerImpl) Flush() error {
	if m.collector == nil {
		return fmt.Errorf("no collector available for flush operation")
	}

	m.collector.FlushAll()
	return nil
}

func (m *AIOutputManagerImpl) FlushStdout() error {
	if m.collector == nil {
		return fmt.Errorf("no collector available for flush operation")
	}

	m.collector.FlushStdout()
	return nil
}

func (m *AIOutputManagerImpl) FlushStderr() error {
	if m.collector == nil {
		return fmt.Errorf("no collector available for flush operation")
	}

	m.collector.FlushStderr()
	return nil
}

func (m *AIOutputManagerImpl) Clear() error {
	if m.collector == nil {
		return fmt.Errorf("no collector available for clear operation")
	}

	m.collector.Clear()
	return nil
}

// [CRITICAL] FMT-001: Output state management - [ACTION:format-processing]
func (m *AIOutputManagerImpl) IsDelayedMode() bool {
	return m.delayed
}

func (m *AIOutputManagerImpl) SetDelayedMode(enabled bool) error {
	if enabled && m.collector == nil {
		return fmt.Errorf("cannot enable delayed mode without collector")
	}

	m.delayed = enabled
	return nil
}

func (m *AIOutputManagerImpl) GetCollectedMessages() []AIOutputMessage {
	if m.collector == nil {
		return []AIOutputMessage{}
	}

	legacyMessages := m.collector.GetMessages()
	aiMessages := make([]AIOutputMessage, len(legacyMessages))

	for i, msg := range legacyMessages {
		aiMessages[i] = AIOutputMessage{
			Content:     msg.Content,
			Destination: AIOutputDestination(msg.Destination),
			Type:        AIMessageType(msg.Type),
			Metadata:    make(map[string]string), // Legacy doesn't have metadata
			Timestamp:   time.Now(),              // Legacy doesn't have timestamp, use current time
		}
	}

	return aiMessages
}

// [CRITICAL] FMT-001: AI-friendly context printing - [ACTION:format-processing]
func (m *AIOutputManagerImpl) PrintWithContext(ctx PrintContext) error {
	if ctx.Options.Delayed {
		message := AIOutputMessage{
			Content:     ctx.Message,
			Destination: ctx.Destination,
			Type:        ctx.Type,
			Metadata:    ctx.Metadata,
			Timestamp:   time.Now(),
		}

		err := m.Collect(message)
		if err != nil {
			return err
		}

		if ctx.Options.Flush {
			return m.Flush()
		}

		return nil
	} else {
		if ctx.Destination == AIOutputDestinationStderr {
			return m.PrintError(ctx.Message)
		} else {
			return m.Print(ctx.Message)
		}
	}
}
