// [IMPL:DELAYED_OUTPUT] [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]
// OutputCollector implementation — links to STDD tokens for traceability.
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// OutputMessage represents a message that can be displayed later
type OutputMessage struct {
	Content     string
	Destination string // "stdout" or "stderr"
	Type        string // "info", "error", "warning", etc.
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// OutputCollector collects output messages for delayed display
type OutputCollector struct {
	messages []OutputMessage
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// NewOutputCollector creates a new OutputCollector
func NewOutputCollector() *OutputCollector {
	return &OutputCollector{
		messages: make([]OutputMessage, 0),
	}
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AddStdout adds a stdout message to the collector
func (oc *OutputCollector) AddStdout(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stdout",
		Type:        messageType,
	})
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// AddStderr adds a stderr message to the collector
func (oc *OutputCollector) AddStderr(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stderr",
		Type:        messageType,
	})
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GetMessages returns all collected messages
func (oc *OutputCollector) GetMessages() []OutputMessage {
	return oc.messages
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
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

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// Clear removes all collected messages without displaying them
func (oc *OutputCollector) Clear() {
	oc.messages = make([]OutputMessage, 0)
}
