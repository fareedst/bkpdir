// Package config provides utility components for configuration management.
//
// This file contains supporting utilities for the configuration system including
// file operations, validation, and source determination.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// 🔺 EXTRACT-001: Configuration structure extraction - Default file operations implementation - 📝
// DefaultFileOperations provides the default implementation of ConfigFileOperations.
// This implementation uses the standard os package for file system operations.
type DefaultFileOperations struct{}

// FileExists checks if a configuration file exists.
// 🔺 EXTRACT-001: Configuration structure extraction - File existence checking - 🔧
func (d *DefaultFileOperations) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads configuration file contents.
// 🔺 EXTRACT-001: Configuration structure extraction - File content reading - 🔧
func (d *DefaultFileOperations) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes configuration file contents.
// 🔺 EXTRACT-001: Configuration structure extraction - File content writing - 🔧
func (d *DefaultFileOperations) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// GetFileInfo returns file information for configuration files.
// 🔺 EXTRACT-001: Configuration structure extraction - File information access - 🔧
func (d *DefaultFileOperations) GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// 🔺 EXTRACT-001: Loading engine extraction - Generic configuration validator - 🛡️
// GenericValidator provides basic validation for any configuration structure.
// This validator performs type checking and basic validation rules.
type GenericValidator struct {
	rules map[string]ValidationRule
}

// NewGenericValidator creates a new GenericValidator.
// 🔺 EXTRACT-001: Loading engine extraction - Validator constructor - 🔧
func NewGenericValidator() *GenericValidator {
	return &GenericValidator{
		rules: make(map[string]ValidationRule),
	}
}

// ValidateSchema validates the configuration against the expected schema.
// 🔺 EXTRACT-001: Loading engine extraction - Schema validation implementation - 🛡️
func (g *GenericValidator) ValidateSchema(cfg interface{}) error {
	// Basic type validation - ensure it's a struct or pointer to struct
	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("configuration must be a struct, got %s", v.Kind())
	}

	// Additional validation can be added here
	return nil
}

// ValidateValues validates individual configuration values.
// 🔺 EXTRACT-001: Loading engine extraction - Value validation implementation - 🛡️
func (g *GenericValidator) ValidateValues(values map[string]ConfigValue) error {
	for name, value := range values {
		if rule, exists := g.rules[name]; exists {
			if err := g.validateValue(value, rule); err != nil {
				return fmt.Errorf("validation failed for field %s: %w", name, err)
			}
		}
	}
	return nil
}

// GetRequiredFields returns the list of required configuration fields.
// 🔶 EXTRACT-001: Value tracking abstraction - Required field identification - 📝
func (g *GenericValidator) GetRequiredFields() []string {
	var required []string
	for name, rule := range g.rules {
		if rule.Required {
			required = append(required, name)
		}
	}
	return required
}

// GetValidationRules returns validation rules for configuration fields.
// 🔶 EXTRACT-001: Value tracking abstraction - Validation rule access - 📝
func (g *GenericValidator) GetValidationRules() map[string]ValidationRule {
	result := make(map[string]ValidationRule)
	for k, v := range g.rules {
		result[k] = v
	}
	return result
}

// SetValidationRule sets a validation rule for a field.
// 🔺 EXTRACT-001: Loading engine extraction - Validation rule configuration - 🔧
func (g *GenericValidator) SetValidationRule(fieldName string, rule ValidationRule) {
	g.rules[fieldName] = rule
}

// validateValue validates a single configuration value against a rule.
// 🔺 EXTRACT-001: Loading engine extraction - Individual value validation - 🛡️
func (g *GenericValidator) validateValue(value ConfigValue, rule ValidationRule) error {
	// Required field validation
	if rule.Required && (value.Value == nil || isZeroValue(value.Value)) {
		return fmt.Errorf("required field is missing or empty")
	}

	// Type validation
	if rule.Type != "" && value.Type != rule.Type {
		return fmt.Errorf("expected type %s, got %s", rule.Type, value.Type)
	}

	// String length validation
	if strVal, ok := value.Value.(string); ok {
		if rule.MinLength > 0 && len(strVal) < rule.MinLength {
			return fmt.Errorf("string length %d is less than minimum %d", len(strVal), rule.MinLength)
		}
		if rule.MaxLength > 0 && len(strVal) > rule.MaxLength {
			return fmt.Errorf("string length %d exceeds maximum %d", len(strVal), rule.MaxLength)
		}
	}

	// Valid values validation
	if len(rule.ValidValues) > 0 {
		strVal := fmt.Sprintf("%v", value.Value)
		for _, valid := range rule.ValidValues {
			if strVal == valid {
				return nil
			}
		}
		return fmt.Errorf("value %s is not in allowed values %v", strVal, rule.ValidValues)
	}

	return nil
}

// 🔶 EXTRACT-001: Value tracking abstraction - Generic source determination - 📝
// GenericSourceDeterminer provides generic configuration source determination.
// This component tracks where configuration values originate from.
type GenericSourceDeterminer struct {
	root          string
	pathDiscovery *PathDiscovery
}

// NewGenericSourceDeterminer creates a new GenericSourceDeterminer.
// 🔶 EXTRACT-001: Value tracking abstraction - Source determiner constructor - 🔧
func NewGenericSourceDeterminer(root string, pathDiscovery *PathDiscovery) *GenericSourceDeterminer {
	return &GenericSourceDeterminer{
		root:          root,
		pathDiscovery: pathDiscovery,
	}
}

// DetermineSource determines the source of a configuration value.
// 🔶 EXTRACT-001: Value tracking abstraction - Source determination logic - 📝
func (g *GenericSourceDeterminer) DetermineSource(current, defaultValue interface{}) string {
	// Check if value is different from default
	if !reflect.DeepEqual(current, defaultValue) {
		// Check if it might be from environment variable
		if g.pathDiscovery != nil {
			envVar := g.pathDiscovery.GetEnvVarName()
			if os.Getenv(envVar) != "" {
				return "environment"
			}
		}
		return "config file"
	}
	return "default"
}

// GetConfigSource returns the primary configuration source.
// 🔶 EXTRACT-001: Value tracking abstraction - Primary source identification - 📝
func (g *GenericSourceDeterminer) GetConfigSource() string {
	if g.pathDiscovery != nil {
		searchPaths := g.pathDiscovery.GetConfigSearchPaths()
		for _, path := range searchPaths {
			expandedPath := g.pathDiscovery.ExpandPath(path)
			if _, err := os.Stat(expandedPath); err == nil {
				return expandedPath
			}
		}
	}
	return "default"
}

// GetSourcePriority returns the priority order of configuration sources.
// 🔶 EXTRACT-001: Value tracking abstraction - Source priority definition - 📝
func (g *GenericSourceDeterminer) GetSourcePriority() []string {
	return []string{"environment", "config file", "default"}
}

// Utility functions

// parseInt64 parses a string to int64.
// 🔺 EXTRACT-001: Environment variable abstraction - String to integer conversion utility - 🔧
func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// isZeroValue checks if a value is the zero value for its type.
// 🔺 EXTRACT-001: Loading engine extraction - Zero value detection utility - 🔧
func isZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	return v.IsZero()
}

// GetEnvForField returns the environment variable value for a specific config field.
// This method is added to support the GenericConfigLoader's environment integration.
// 🔺 EXTRACT-001: Environment variable abstraction - Field-based environment access utility - 🔧
func GetEnvForField(envProvider EnvironmentProvider, fieldName string) string {
	if provider, ok := envProvider.(*DefaultEnvironmentProvider); ok {
		return provider.GetEnvForField(fieldName)
	}

	// Fallback: try to get the environment variable directly
	mapping := envProvider.GetEnvMapping()
	if envVar, exists := mapping[fieldName]; exists {
		return envProvider.GetEnv(envVar)
	}

	return ""
}

// NewDefaultFileOperations creates a new DefaultFileOperations instance.
// 🔺 EXTRACT-001: Configuration structure extraction - Default file operations constructor - 🔧
func NewDefaultFileOperations() *DefaultFileOperations {
	return &DefaultFileOperations{}
}

// 🔺 EXTRACT-001: Value tracking abstraction - Generic value extractor implementation - 📝
// GenericValueExtractor provides generic configuration value extraction.
// This extractor works with any configuration structure using reflection.
type GenericValueExtractor struct{}

// NewGenericValueExtractor creates a new GenericValueExtractor.
// 🔺 EXTRACT-001: Value tracking abstraction - Value extractor constructor - 🔧
func NewGenericValueExtractor() *GenericValueExtractor {
	return &GenericValueExtractor{}
}

// ExtractValues extracts configuration values from any configuration structure.
// 🔺 EXTRACT-001: Value tracking abstraction - Generic value extraction - 📝
func (g *GenericValueExtractor) ExtractValues(cfg, defaultCfg interface{}, getSource func(interface{}, interface{}) string) []ConfigValue {
	var values []ConfigValue

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return values
	}

	defaultV := reflect.ValueOf(defaultCfg)
	if defaultV.Kind() == reflect.Ptr {
		defaultV = defaultV.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		// Get YAML tag name or use field name
		name := field.Name
		if yamlTag := field.Tag.Get("yaml"); yamlTag != "" && yamlTag != "-" {
			name = yamlTag
		}

		// Get default value for comparison
		var defaultValue interface{}
		if defaultV.IsValid() && i < defaultV.NumField() {
			defaultFieldValue := defaultV.Field(i)
			if defaultFieldValue.CanInterface() {
				defaultValue = defaultFieldValue.Interface()
			}
		}

		// Determine source using provided function
		source := "default"
		if getSource != nil {
			source = getSource(fieldValue.Interface(), defaultValue)
		}

		values = append(values, ConfigValue{
			Name:   name,
			Value:  fieldValue.Interface(),
			Source: source,
			Type:   fieldValue.Type().String(),
		})
	}

	return values
}

// ExtractValuesByCategory extracts values from a specific category.
// 🔺 EXTRACT-001: Value tracking abstraction - Category-based value extraction - 📝
func (g *GenericValueExtractor) ExtractValuesByCategory(cfg, defaultCfg interface{}, category string, getSource func(interface{}, interface{}) string) []ConfigValue {
	// For now, return all values regardless of category
	// Can be enhanced later for specific categorization
	return g.ExtractValues(cfg, defaultCfg, getSource)
}

// GetSupportedCategories returns the list of supported value categories.
// 🔶 EXTRACT-001: Value tracking abstraction - Category support definition - 📝
func (g *GenericValueExtractor) GetSupportedCategories() []string {
	return []string{"all", "basic", "advanced", "format", "paths"}
}
