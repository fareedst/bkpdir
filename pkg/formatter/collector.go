// [IMPL-DELAYED_OUTPUT] [ARCH-OUTPUT_FORMATTING] [REQ-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// OutputMessage represents a message that can be displayed later
type OutputMessage struct {
	Content     string
	Destination string // "stdout" or "stderr"
	Type        string // "info", "error", "warning", etc.
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// OutputCollector collects output messages for delayed display
type OutputCollector struct {
	messages []OutputMessage
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// NewOutputCollector creates a new OutputCollector
func NewOutputCollector() *OutputCollector {
	return &OutputCollector{
		messages: make([]OutputMessage, 0),
	}
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// AddStdout adds a stdout message to the collector
func (oc *OutputCollector) AddStdout(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stdout",
		Type:        messageType,
	})
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// AddStderr adds a stderr message to the collector
func (oc *OutputCollector) AddStderr(content, messageType string) {
	oc.messages = append(oc.messages, OutputMessage{
		Content:     content,
		Destination: "stderr",
		Type:        messageType,
	})
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// GetMessages returns all collected messages
func (oc *OutputCollector) GetMessages() []OutputMessage {
	return oc.messages
}

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
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

// [IMPL-DUAL_FORMATTING] [ARCH-OUTPUT_FORMATTING]
// Clear removes all collected messages without displaying them
func (oc *OutputCollector) Clear() {
	oc.messages = make([]OutputMessage, 0)
}
