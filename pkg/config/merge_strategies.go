// Package config provides merge strategies for layered configuration inheritance.
//
// This file implements the merge strategies that enable flexible configuration
// inheritance with different behaviors for arrays, objects, and primitive values.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"fmt"
	"reflect"
	"strings"
)

// ⭐ CFG-005: Merge strategy implementations - 🔧 Core merge functionality

// StandardOverrideStrategy implements the default override behavior (no prefix).
// Child values completely replace parent values.
type StandardOverrideStrategy struct{}

// ⭐ CFG-005: Standard override merge strategy - 🔧 Default behavior
func (s *StandardOverrideStrategy) Merge(dest, src interface{}) (interface{}, error) {
	// Standard override: child value completely replaces parent value
	return src, nil
}

func (s *StandardOverrideStrategy) GetPrefix() string {
	return "" // No prefix for standard override
}

func (s *StandardOverrideStrategy) GetDescription() string {
	return "Child values completely replace parent values"
}

func (s *StandardOverrideStrategy) SupportsType(valueType string) bool {
	return true // Supports all types
}

// ArrayMergeStrategy implements array merge behavior (+ prefix).
// Child array elements are appended to parent array elements.
type ArrayMergeStrategy struct{}

// ⭐ CFG-005: Array merge strategy implementation - 🔧 Array append behavior
func (a *ArrayMergeStrategy) Merge(dest, src interface{}) (interface{}, error) {
	destVal := reflect.ValueOf(dest)
	srcVal := reflect.ValueOf(src)

	// Handle nil cases
	if !destVal.IsValid() || destVal.IsNil() {
		return src, nil
	}
	if !srcVal.IsValid() || srcVal.IsNil() {
		return dest, nil
	}

	// Both must be slices for array merge
	if destVal.Kind() != reflect.Slice || srcVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("array merge strategy requires slice types, got %s and %s",
			destVal.Kind(), srcVal.Kind())
	}

	// Create new slice with combined capacity
	destLen := destVal.Len()
	srcLen := srcVal.Len()
	newSlice := reflect.MakeSlice(destVal.Type(), destLen+srcLen, destLen+srcLen)

	// Copy destination elements first
	for i := 0; i < destLen; i++ {
		newSlice.Index(i).Set(destVal.Index(i))
	}

	// Append source elements
	for i := 0; i < srcLen; i++ {
		newSlice.Index(destLen + i).Set(srcVal.Index(i))
	}

	return newSlice.Interface(), nil
}

func (a *ArrayMergeStrategy) GetPrefix() string {
	return "+"
}

func (a *ArrayMergeStrategy) GetDescription() string {
	return "Child array elements are appended to parent array elements"
}

func (a *ArrayMergeStrategy) SupportsType(valueType string) bool {
	return valueType == "slice" || strings.Contains(valueType, "[]")
}

// ArrayPrependStrategy implements array prepend behavior (^ prefix).
// Child array elements are prepended to parent array elements (higher priority).
type ArrayPrependStrategy struct{}

// ⭐ CFG-005: Array prepend strategy implementation - 🔧 Array prepend behavior
func (a *ArrayPrependStrategy) Merge(dest, src interface{}) (interface{}, error) {
	destVal := reflect.ValueOf(dest)
	srcVal := reflect.ValueOf(src)

	// Handle nil cases
	if !destVal.IsValid() || destVal.IsNil() {
		return src, nil
	}
	if !srcVal.IsValid() || srcVal.IsNil() {
		return dest, nil
	}

	// Both must be slices for array prepend
	if destVal.Kind() != reflect.Slice || srcVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("array prepend strategy requires slice types, got %s and %s",
			destVal.Kind(), srcVal.Kind())
	}

	// Create new slice with combined capacity
	destLen := destVal.Len()
	srcLen := srcVal.Len()
	newSlice := reflect.MakeSlice(destVal.Type(), destLen+srcLen, destLen+srcLen)

	// Copy source elements first (prepend)
	for i := 0; i < srcLen; i++ {
		newSlice.Index(i).Set(srcVal.Index(i))
	}

	// Append destination elements after
	for i := 0; i < destLen; i++ {
		newSlice.Index(srcLen + i).Set(destVal.Index(i))
	}

	return newSlice.Interface(), nil
}

func (a *ArrayPrependStrategy) GetPrefix() string {
	return "^"
}

func (a *ArrayPrependStrategy) GetDescription() string {
	return "Child array elements are prepended to parent array elements (higher priority)"
}

func (a *ArrayPrependStrategy) SupportsType(valueType string) bool {
	return valueType == "slice" || strings.Contains(valueType, "[]")
}

// ArrayReplaceStrategy implements array replace behavior (! prefix).
// Child array completely replaces parent array.
type ArrayReplaceStrategy struct{}

// ⭐ CFG-005: Array replace strategy implementation - 🔧 Array replacement behavior
func (a *ArrayReplaceStrategy) Merge(dest, src interface{}) (interface{}, error) {
	// Array replace: child array completely replaces parent array
	return src, nil
}

func (a *ArrayReplaceStrategy) GetPrefix() string {
	return "!"
}

func (a *ArrayReplaceStrategy) GetDescription() string {
	return "Child array completely replaces parent array"
}

func (a *ArrayReplaceStrategy) SupportsType(valueType string) bool {
	return valueType == "slice" || strings.Contains(valueType, "[]")
}

// DefaultValueStrategy implements default value behavior (= prefix).
// Child value is used only if parent value is not set or is zero value.
type DefaultValueStrategy struct{}

// ⭐ CFG-005: Default value strategy implementation - 🔧 Default fallback behavior
func (d *DefaultValueStrategy) Merge(dest, src interface{}) (interface{}, error) {
	destVal := reflect.ValueOf(dest)

	// If destination is nil or zero value, use source
	if !destVal.IsValid() || destVal.IsZero() {
		return src, nil
	}

	// Otherwise keep destination value
	return dest, nil
}

func (d *DefaultValueStrategy) GetPrefix() string {
	return "="
}

func (d *DefaultValueStrategy) GetDescription() string {
	return "Child value is used only if parent value is not set or is zero value"
}

func (d *DefaultValueStrategy) SupportsType(valueType string) bool {
	return true // Supports all types
}

// PrefixedKeyProcessor processes configuration keys with merge strategy prefixes.
type PrefixedKeyProcessor struct {
	strategies map[string]MergeStrategy
}

// ⭐ CFG-005: Prefixed key processor implementation - 🔍 Key prefix analysis
func NewPrefixedKeyProcessor() *PrefixedKeyProcessor {
	processor := &PrefixedKeyProcessor{
		strategies: make(map[string]MergeStrategy),
	}

	// Register default strategies
	processor.RegisterStrategy(&StandardOverrideStrategy{})
	processor.RegisterStrategy(&ArrayMergeStrategy{})
	processor.RegisterStrategy(&ArrayPrependStrategy{})
	processor.RegisterStrategy(&ArrayReplaceStrategy{})
	processor.RegisterStrategy(&DefaultValueStrategy{})

	return processor
}

func (p *PrefixedKeyProcessor) RegisterStrategy(strategy MergeStrategy) error {
	if strategy == nil {
		return fmt.Errorf("merge strategy cannot be nil")
	}

	prefix := strategy.GetPrefix()
	p.strategies[prefix] = strategy
	return nil
}

func (p *PrefixedKeyProcessor) GetAvailableStrategies() map[string]MergeStrategy {
	// Return copy to prevent external modification
	strategies := make(map[string]MergeStrategy)
	for prefix, strategy := range p.strategies {
		strategies[prefix] = strategy
	}
	return strategies
}

// ⭐ CFG-005: Key processing implementation - 🔧 Prefix extraction and strategy selection
func (p *PrefixedKeyProcessor) ProcessKeys(config map[string]interface{}) (*ProcessedConfig, error) {
	processed := &ProcessedConfig{
		Config:     make(map[string]interface{}),
		MergeOps:   make([]MergeOperation, 0),
		SourceMap:  make(map[string]string),
		Strategies: make(map[string]string),
	}

	for key, value := range config {
		strategy, cleanKey := p.extractMergeStrategy(key)

		// Store the processed key-value pair
		processed.Config[cleanKey] = value

		// Record the merge operation
		processed.MergeOps = append(processed.MergeOps, MergeOperation{
			Key:        cleanKey,
			Value:      value,
			Strategy:   strategy,
			SourceFile: "", // Will be filled by caller
			TargetType: reflect.TypeOf(value).String(),
		})

		// Record the strategy used
		processed.Strategies[cleanKey] = strategy.GetDescription()
	}

	return processed, nil
}

// extractMergeStrategy extracts the merge strategy and clean key from a prefixed key.
func (p *PrefixedKeyProcessor) extractMergeStrategy(key string) (MergeStrategy, string) {
	// Check for strategy prefixes
	for prefix, strategy := range p.strategies {
		if prefix != "" && strings.HasPrefix(key, prefix) {
			// Return strategy and key without prefix
			return strategy, strings.TrimPrefix(key, prefix)
		}
	}

	// Default to standard override strategy (no prefix)
	return p.strategies[""], key
}

// DefaultMergeStrategyProcessor provides a configured merge strategy processor.
type DefaultMergeStrategyProcessor struct {
	*PrefixedKeyProcessor
}

// ⭐ CFG-005: Default merge strategy processor - 🔧 Pre-configured processor
func NewDefaultMergeStrategyProcessor() *DefaultMergeStrategyProcessor {
	return &DefaultMergeStrategyProcessor{
		PrefixedKeyProcessor: NewPrefixedKeyProcessor(),
	}
}

func (d *DefaultMergeStrategyProcessor) ProcessKeys(config map[string]interface{}) (*ProcessedConfig, error) {
	return d.PrefixedKeyProcessor.ProcessKeys(config)
}

func (d *DefaultMergeStrategyProcessor) RegisterStrategy(strategy MergeStrategy) error {
	return d.PrefixedKeyProcessor.RegisterStrategy(strategy)
}

func (d *DefaultMergeStrategyProcessor) GetAvailableStrategies() map[string]MergeStrategy {
	return d.PrefixedKeyProcessor.GetAvailableStrategies()
}
