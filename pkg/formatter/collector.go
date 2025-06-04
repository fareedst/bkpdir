// Output collection and delayed display functionality for the formatter package.
// Provides the ability to collect output messages and display them later,
// supporting both stdout and stderr destinations with message typing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package formatter

import (
	"fmt"
	"os"
)

// ⭐ EXTRACT-003: OutputCollector component - 🔧 Output message structure
// OutputMessage represents a message that can be displayed later
type OutputMessage struct {
	Content     string
	Destination string // "stdout" or "stderr"
	Type        string // "info", "error", "warning", etc.
}

// ⭐ EXTRACT-003: OutputCollector component - 🔧 Delayed output management
// OutputCollector collects output messages for delayed display
type OutputCollector struct {
	messages []OutputMessage
}

// ⭐ EXTRACT-003: OutputCollector component - 🔧 Constructor
// NewOutputCollector creates a new OutputCollector
func NewOutputCollector() *OutputCollector {
	return &OutputCollector{
		messages: make([]OutputMessage, 0),
	}
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Add stdout message
// AddStdout adds a stdout message to the collector
func (oc *OutputCollector) AddStdout(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stdout",
		Type:        messageType,
	})
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Add stderr message
// AddStderr adds a stderr message to the collector
func (oc *OutputCollector) AddStderr(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stderr",
		Type:        messageType,
	})
}

// ⭐ EXTRACT-003: OutputCollector component - 🔍 Get all messages
// GetMessages returns all collected messages
func (oc *OutputCollector) GetMessages() []OutputMessage {
	return oc.messages
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Flush all messages
// FlushAll displays all collected messages and clears the collector
func (oc *OutputCollector) FlushAll() {
	for _, msg := range oc.messages {
		if msg.Destination == "stderr" {
			fmt.Fprint(os.Stderr, msg.Content)
		} else {
			fmt.Print(msg.Content)
		}
	}
	oc.messages = make([]OutputMessage, 0)
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Flush stdout only
// FlushStdout displays only stdout messages and removes them from the collector
func (oc *OutputCollector) FlushStdout() {
	remaining := make([]OutputMessage, 0)
	for _, msg := range oc.messages {
		if msg.Destination == "stdout" {
			fmt.Print(msg.Content)
		} else {
			remaining = append(remaining, msg)
		}
	}
	oc.messages = remaining
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Flush stderr only
// FlushStderr displays only stderr messages and removes them from the collector
func (oc *OutputCollector) FlushStderr() {
	remaining := make([]OutputMessage, 0)
	for _, msg := range oc.messages {
		if msg.Destination == "stderr" {
			fmt.Fprint(os.Stderr, msg.Content)
		} else {
			remaining = append(remaining, msg)
		}
	}
	oc.messages = remaining
}

// ⭐ EXTRACT-003: OutputCollector component - 📝 Clear all messages
// Clear removes all collected messages without displaying them
func (oc *OutputCollector) Clear() {
	oc.messages = make([]OutputMessage, 0)
}
