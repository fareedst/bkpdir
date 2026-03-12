// [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION]
// Configuration loader implementation — links to STDD tokens for traceability.
// Package config provides schema-agnostic configuration loading and merging.
//
// This file contains the core configuration loading engine extracted from the
// original backup application, generalized to work with any configuration schema.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"fmt"
	"reflect"

	yaml "gopkg.in/yaml.v3"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GenericConfigLoader provides schema-agnostic configuration loading and merging.
// This replaces the backup-specific LoadConfig function with a generalized implementation.
type GenericConfigLoader struct {
	pathDiscovery *PathDiscovery
	envProvider   EnvironmentProvider
	fileOps       ConfigFileOperations
	validator     ConfigValidator
}

// NewGenericConfigLoader creates a new GenericConfigLoader with specified components.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewGenericConfigLoader(pathDiscovery *PathDiscovery, envProvider EnvironmentProvider, fileOps ConfigFileOperations, validator ConfigValidator) *GenericConfigLoader {
	return &GenericConfigLoader{
		pathDiscovery: pathDiscovery,
		envProvider:   envProvider,
		fileOps:       fileOps,
		validator:     validator,
	}
}

// NewDefaultConfigLoader creates a GenericConfigLoader with backup application defaults.
// This maintains backward compatibility while providing the new generic interface.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewDefaultConfigLoader() *GenericConfigLoader {
	return NewGenericConfigLoader(
		NewDefaultPathDiscovery(),
		NewBackupEnvironmentProvider(),
		&DefaultFileOperations{},
		&GenericValidator{},
	)
}

// LoadConfig loads configuration from multiple sources with the specified default.
// This generalizes the original LoadConfig function to work with any configuration type.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) LoadConfig(root string, defaultConfig interface{}) (interface{}, error) {
	// Start with default configuration (clone to avoid modifying original)
	config := g.cloneConfig(defaultConfig)

	// Get search paths and process each configuration file
	searchPaths := g.pathDiscovery.GetConfigSearchPaths()
	for _, path := range searchPaths {
		expandedPath := g.pathDiscovery.ExpandPath(path)
		if g.fileOps.FileExists(expandedPath) {
			if err := g.loadConfigFromFile(config, expandedPath); err != nil {
				// Continue with other files if one fails (non-fatal error)
				continue
			}
		}
	}

	// Apply environment variable overrides
	if err := g.applyEnvironmentOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	// Validate the final configuration
	if g.validator != nil {
		if err := g.validator.ValidateSchema(config); err != nil {
			return nil, fmt.Errorf("configuration validation failed: %w", err)
		}
	}

	return config, nil
}

// LoadConfigValues loads configuration values with source tracking.
// This generalizes the original LoadConfigValues function.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) LoadConfigValues(root string, defaultConfig interface{}) (map[string]ConfigValue, error) {
	config, err := g.LoadConfig(root, defaultConfig)
	if err != nil {
		return nil, err
	}

	// Extract configuration values with source tracking
	values := make(map[string]ConfigValue)
	configValues := g.GetConfigValuesWithSources(config, root)
	for _, cv := range configValues {
		values[cv.Name] = cv
	}

	return values, nil
}

// GetConfigValues extracts configuration values from a configuration struct.
// This generalizes configuration value extraction for any struct type.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) GetConfigValues(cfg interface{}) []ConfigValue {
	var values []ConfigValue

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return values
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

		values = append(values, ConfigValue{
			Name:   name,
			Value:  fieldValue.Interface(),
			Source: "default",
			Type:   fieldValue.Type().String(),
		})
	}

	return values
}

// GetConfigValuesWithSources extracts configuration values with source information.
// This provides enhanced source tracking for configuration debugging.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) GetConfigValuesWithSources(cfg interface{}, root string) []ConfigValue {
	// Start with basic values
	values := g.GetConfigValues(cfg)

	// Enhance with source information
	sourceDeterminer := NewGenericSourceDeterminer(root, g.pathDiscovery)
	defaultConfig := g.getDefaultForType(cfg)

	for i := range values {
		values[i].Source = sourceDeterminer.DetermineSource(values[i].Value, g.getDefaultValue(defaultConfig, values[i].Name))
	}

	return values
}

// ValidateConfig validates a configuration structure.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) ValidateConfig(cfg interface{}) error {
	if g.validator == nil {
		return nil
	}
	return g.validator.ValidateSchema(cfg)
}

// GetConfigSearchPaths returns the search paths for configuration files.
// This implements the ConfigMerger interface method.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) GetConfigSearchPaths() []string {
	return g.pathDiscovery.GetConfigSearchPaths()
}

// ExpandPath expands path variables and returns absolute path.
// This implements the ConfigMerger interface method.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) ExpandPath(path string) string {
	return g.pathDiscovery.ExpandPath(path)
}

// MergeConfigValues merges configuration value maps.
// This implements the ConfigMerger interface method.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) MergeConfigValues(dst, src map[string]ConfigValue) {
	for key, value := range src {
		dst[key] = value
	}
}

// MergeConfigs merges source configuration into destination configuration.
// This provides generic configuration merging using reflection.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) MergeConfigs(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	if dstVal.Kind() != reflect.Ptr || srcVal.Kind() != reflect.Ptr {
		return fmt.Errorf("both dst and src must be pointers")
	}

	dstElem := dstVal.Elem()
	srcElem := srcVal.Elem()

	if dstElem.Type() != srcElem.Type() {
		return fmt.Errorf("dst and src must be the same type")
	}

	return g.mergeValues(dstElem, srcElem)
}

// mergeValues recursively merges values using reflection.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) mergeValues(dst, src reflect.Value) error {
	switch src.Kind() {
	case reflect.Struct:
		return g.mergeStruct(dst, src)
	case reflect.Slice:
		return g.mergeSlice(dst, src)
	case reflect.Map:
		return g.mergeMap(dst, src)
	case reflect.Ptr:
		return g.mergePointer(dst, src)
	default:
		// For basic types, copy non-zero values
		if !src.IsZero() {
			dst.Set(src)
		}
		return nil
	}
}

// mergeStruct merges struct fields.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) mergeStruct(dst, src reflect.Value) error {
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)

		if !dstField.CanSet() {
			continue
		}

		if err := g.mergeValues(dstField, srcField); err != nil {
			return err
		}
	}
	return nil
}

// mergeSlice merges slice values (src overwrites dst if non-empty).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) mergeSlice(dst, src reflect.Value) error {
	if src.Len() > 0 {
		dst.Set(src)
	}
	return nil
}

// mergeMap merges map values.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) mergeMap(dst, src reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	for _, key := range src.MapKeys() {
		dst.SetMapIndex(key, src.MapIndex(key))
	}
	return nil
}

// mergePointer merges pointer values.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) mergePointer(dst, src reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}

	return g.mergeValues(dst.Elem(), src.Elem())
}

// loadConfigFromFile loads configuration from a specific file.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) loadConfigFromFile(config interface{}, path string) error {
	data, err := g.fileOps.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Create a temporary config of the same type for unmarshaling
	tempConfig := g.cloneConfig(config)
	if err := yaml.Unmarshal(data, tempConfig); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Merge the loaded config into the main config
	return g.MergeConfigs(config, tempConfig)
}

// applyEnvironmentOverrides applies environment variable overrides to configuration.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) applyEnvironmentOverrides(config interface{}) error {
	if g.envProvider == nil {
		return nil
	}

	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		// Get field name from YAML tag or field name
		fieldName := field.Name
		if yamlTag := field.Tag.Get("yaml"); yamlTag != "" && yamlTag != "-" {
			fieldName = yamlTag
		}

		// Check if there's an environment variable for this field
		if envValue := GetEnvForField(g.envProvider, fieldName); envValue != "" {
			if err := g.setFieldFromString(fieldValue, envValue); err != nil {
				return fmt.Errorf("failed to set field %s from environment: %w", fieldName, err)
			}
		}
	}

	return nil
}

// cloneConfig creates a deep copy of the configuration.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) cloneConfig(config interface{}) interface{} {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Create new instance of the same type
	newConfig := reflect.New(v.Type())
	newConfig.Elem().Set(v)

	return newConfig.Interface()
}

// getDefaultForType creates a default instance of the given configuration type.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) getDefaultForType(config interface{}) interface{} {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return reflect.New(v.Type()).Interface()
}

// getDefaultValue gets the default value for a field name.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) getDefaultValue(defaultConfig interface{}, fieldName string) interface{} {
	v := reflect.ValueOf(defaultConfig)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	// Find field by name or YAML tag
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		name := field.Name
		if yamlTag := field.Tag.Get("yaml"); yamlTag != "" && yamlTag != "-" {
			name = yamlTag
		}

		if name == fieldName {
			return v.Field(i).Interface()
		}
	}

	return nil
}

// setFieldFromString sets a field value from a string representation.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigLoader) setFieldFromString(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Bool:
		if value == "true" || value == "1" {
			field.SetBool(true)
		} else if value == "false" || value == "0" {
			field.SetBool(false)
		} else {
			return fmt.Errorf("invalid boolean value: %s", value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := parseInt64(value); err == nil {
			field.SetInt(intVal)
		} else {
			return err
		}
	default:
		return fmt.Errorf("unsupported field type for environment override: %s", field.Kind())
	}

	return nil
}

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// GenericConfigMerger provides standalone configuration merging functionality.
// This is a separate component that can be used independently of the loader.
type GenericConfigMerger struct {
	pathDiscovery *PathDiscovery
}

// NewGenericConfigMerger creates a new GenericConfigMerger.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewGenericConfigMerger() *GenericConfigMerger {
	return &GenericConfigMerger{
		pathDiscovery: NewDefaultPathDiscovery(),
	}
}

// NewGenericConfigMergerWithDiscovery creates a new GenericConfigMerger with custom path discovery.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewGenericConfigMergerWithDiscovery(pathDiscovery *PathDiscovery) *GenericConfigMerger {
	return &GenericConfigMerger{
		pathDiscovery: pathDiscovery,
	}
}

// MergeConfigs merges source configuration into destination configuration.
// This provides generic configuration merging using reflection.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) MergeConfigs(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)

	if dstVal.Kind() != reflect.Ptr || srcVal.Kind() != reflect.Ptr {
		return fmt.Errorf("both dst and src must be pointers")
	}

	dstElem := dstVal.Elem()
	srcElem := srcVal.Elem()

	if dstElem.Type() != srcElem.Type() {
		return fmt.Errorf("dst and src must be the same type")
	}

	return g.mergeValues(dstElem, srcElem)
}

// MergeConfigValues merges configuration value maps.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) MergeConfigValues(dst, src map[string]ConfigValue) {
	for key, value := range src {
		dst[key] = value
	}
}

// GetConfigSearchPaths returns the search paths for configuration files.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) GetConfigSearchPaths() []string {
	return g.pathDiscovery.GetConfigSearchPaths()
}

// ExpandPath expands path variables and returns absolute path.
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) ExpandPath(path string) string {
	return g.pathDiscovery.ExpandPath(path)
}

// mergeValues recursively merges values using reflection (from GenericConfigMerger).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) mergeValues(dst, src reflect.Value) error {
	switch src.Kind() {
	case reflect.Struct:
		return g.mergeStruct(dst, src)
	case reflect.Slice:
		return g.mergeSlice(dst, src)
	case reflect.Map:
		return g.mergeMap(dst, src)
	case reflect.Ptr:
		return g.mergePointer(dst, src)
	default:
		// For basic types, copy non-zero values
		if !src.IsZero() {
			dst.Set(src)
		}
		return nil
	}
}

// mergeStruct merges struct fields (from GenericConfigMerger).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) mergeStruct(dst, src reflect.Value) error {
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)

		if !dstField.CanSet() {
			continue
		}

		if err := g.mergeValues(dstField, srcField); err != nil {
			return err
		}
	}
	return nil
}

// mergeSlice merges slice values (from GenericConfigMerger).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) mergeSlice(dst, src reflect.Value) error {
	if src.Len() > 0 {
		dst.Set(src)
	}
	return nil
}

// mergeMap merges map values (from GenericConfigMerger).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) mergeMap(dst, src reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}

	for _, key := range src.MapKeys() {
		dst.SetMapIndex(key, src.MapIndex(key))
	}
	return nil
}

// mergePointer merges pointer values (from GenericConfigMerger).
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (g *GenericConfigMerger) mergePointer(dst, src reflect.Value) error {
	if src.IsNil() {
		return nil
	}

	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}

	return g.mergeValues(dst.Elem(), src.Elem())
}
