// This file is part of bkpdir
//
// Package main provides concrete implementations of the configuration abstraction
// interfaces. These implementations demonstrate how the abstraction layer enables
// schema-agnostic configuration management.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// 🔶 REFACTOR-003: Config abstraction - Generic configuration loader implementation - 🔧
// GenericConfigLoader implements ConfigLoader for schema-agnostic configuration loading.
type GenericConfigLoader struct {
	sources []ConfigSource
}

// NewGenericConfigLoader creates a new generic configuration loader.
func NewGenericConfigLoader(sources ...ConfigSource) *GenericConfigLoader {
	if len(sources) == 0 {
		// Default sources: YAML files, environment variables, defaults
		sources = []ConfigSource{
			NewYAMLConfigSource(),
			NewEnvConfigSource(),
			NewDefaultConfigSource(),
		}
	}
	return &GenericConfigLoader{sources: sources}
}

// LoadConfig loads configuration using default schema detection.
func (loader *GenericConfigLoader) LoadConfig(root string) (interface{}, error) {
	// Use backup application schema as default
	schema := NewBackupConfigSchema()
	return loader.LoadConfigWithSchema(root, schema)
}

// LoadConfigWithSchema loads configuration with a specific schema.
func (loader *GenericConfigLoader) LoadConfigWithSchema(root string, schema ConfigSchema) (interface{}, error) {
	// Start with default configuration from schema
	config := schema.GetDefaultConfig()

	// Load and merge from all sources in priority order
	for _, source := range loader.sources {
		if source.GetSourceType() == "file" {
			// For file sources, look for config files
			paths := loader.GetConfigSearchPaths()
			for _, path := range paths {
				fullPath := filepath.Join(root, path)
				if source.Exists(fullPath) {
					data, err := source.Load(fullPath)
					if err != nil {
						return nil, fmt.Errorf("failed to load config from %s: %w", fullPath, err)
					}

					// Merge loaded data into config
					if err := loader.mergeIntoConfig(config, data, schema); err != nil {
						return nil, fmt.Errorf("failed to merge config from %s: %w", fullPath, err)
					}
					break // Use first found config file
				}
			}
		} else {
			// For other sources (env, defaults), load without path
			data, err := source.Load("")
			if err != nil {
				continue // Skip sources that fail to load
			}

			// Merge loaded data into config
			if err := loader.mergeIntoConfig(config, data, schema); err != nil {
				return nil, fmt.Errorf("failed to merge config from %s source: %w", source.GetSourceType(), err)
			}
		}
	}

	// Validate final configuration
	if err := schema.ValidateSchema(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// GetConfigSearchPaths returns the search paths for configuration files.
func (loader *GenericConfigLoader) GetConfigSearchPaths() []string {
	return getConfigSearchPaths() // Use existing function
}

// ExpandPath expands environment variables and tildes in paths.
func (loader *GenericConfigLoader) ExpandPath(path string) string {
	return expandPath(path) // Use existing function
}

// mergeIntoConfig merges data map into configuration struct using reflection.
func (loader *GenericConfigLoader) mergeIntoConfig(config interface{}, data map[string]interface{}, schema ConfigSchema) error {
	// Use reflection to set fields based on data map
	configValue := reflect.ValueOf(config)
	if configValue.Kind() == reflect.Ptr {
		configValue = configValue.Elem()
	}

	configType := configValue.Type()

	for key, value := range data {
		// Find field by YAML tag
		for i := 0; i < configType.NumField(); i++ {
			field := configType.Field(i)
			yamlTag := field.Tag.Get("yaml")
			if yamlTag == key || (yamlTag == "" && strings.ToLower(field.Name) == strings.ToLower(key)) {
				fieldValue := configValue.Field(i)
				if fieldValue.CanSet() {
					if err := loader.setFieldValue(fieldValue, value); err != nil {
						return fmt.Errorf("failed to set field %s: %w", field.Name, err)
					}
				}
				break
			}
		}
	}

	return nil
}

// setFieldValue sets a reflect.Value from an interface{} value.
func (loader *GenericConfigLoader) setFieldValue(fieldValue reflect.Value, value interface{}) error {
	switch fieldValue.Kind() {
	case reflect.String:
		if str, ok := value.(string); ok {
			fieldValue.SetString(str)
		} else {
			fieldValue.SetString(fmt.Sprintf("%v", value))
		}
	case reflect.Bool:
		if b, ok := value.(bool); ok {
			fieldValue.SetBool(b)
		} else {
			return fmt.Errorf("cannot convert %v to bool", value)
		}
	case reflect.Int:
		if i, ok := value.(int); ok {
			fieldValue.SetInt(int64(i))
		} else if str, ok := value.(string); ok {
			if i, err := strconv.Atoi(str); err == nil {
				fieldValue.SetInt(int64(i))
			} else {
				return fmt.Errorf("cannot convert %s to int", str)
			}
		} else {
			return fmt.Errorf("cannot convert %v to int", value)
		}
	case reflect.Slice:
		if slice, ok := value.([]interface{}); ok {
			sliceValue := reflect.MakeSlice(fieldValue.Type(), len(slice), len(slice))
			for i, item := range slice {
				if err := loader.setFieldValue(sliceValue.Index(i), item); err != nil {
					return err
				}
			}
			fieldValue.Set(sliceValue)
		} else {
			return fmt.Errorf("cannot convert %v to slice", value)
		}
	case reflect.Ptr:
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		return loader.setFieldValue(fieldValue.Elem(), value)
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}
	return nil
}

// 🔶 REFACTOR-003: Schema separation - YAML configuration source implementation - 🔧
// YAMLConfigSource implements ConfigSource for YAML files.
type YAMLConfigSource struct{}

// NewYAMLConfigSource creates a new YAML configuration source.
func NewYAMLConfigSource() *YAMLConfigSource {
	return &YAMLConfigSource{}
}

// Load loads configuration data from a YAML file.
func (source *YAMLConfigSource) Load(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves configuration data to a YAML file.
func (source *YAMLConfigSource) Save(path string, data map[string]interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, yamlData, 0644)
}

// Exists checks if the YAML configuration file exists.
func (source *YAMLConfigSource) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetSourceType returns the type of source.
func (source *YAMLConfigSource) GetSourceType() string {
	return "file"
}

// GetPriority returns the priority of this source for merging.
func (source *YAMLConfigSource) GetPriority() int {
	return 100 // High priority
}

// 🔶 REFACTOR-003: Config abstraction - Environment variable configuration source - 🔍
// EnvConfigSource implements ConfigSource for environment variables.
type EnvConfigSource struct{}

// NewEnvConfigSource creates a new environment variable configuration source.
func NewEnvConfigSource() *EnvConfigSource {
	return &EnvConfigSource{}
}

// Load loads configuration data from environment variables.
func (source *EnvConfigSource) Load(path string) (map[string]interface{}, error) {
	config := make(map[string]interface{})

	// Look for BKPDIR_ prefixed environment variables
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "BKPDIR_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := strings.ToLower(strings.TrimPrefix(parts[0], "BKPDIR_"))
				key = strings.ReplaceAll(key, "_", "_")
				config[key] = parts[1]
			}
		}
	}

	return config, nil
}

// Save saves configuration data to environment variables (not implemented).
func (source *EnvConfigSource) Save(path string, data map[string]interface{}) error {
	return fmt.Errorf("saving to environment variables not supported")
}

// Exists checks if environment variables exist (always true).
func (source *EnvConfigSource) Exists(path string) bool {
	return true
}

// GetSourceType returns the type of source.
func (source *EnvConfigSource) GetSourceType() string {
	return "env"
}

// GetPriority returns the priority of this source for merging.
func (source *EnvConfigSource) GetPriority() int {
	return 50 // Medium priority
}

// 🔶 REFACTOR-003: Config abstraction - Default configuration source - 🔍
// DefaultConfigSource implements ConfigSource for default values.
type DefaultConfigSource struct{}

// NewDefaultConfigSource creates a new default configuration source.
func NewDefaultConfigSource() *DefaultConfigSource {
	return &DefaultConfigSource{}
}

// Load loads default configuration data.
func (source *DefaultConfigSource) Load(path string) (map[string]interface{}, error) {
	// Return empty map - defaults are handled by schema
	return make(map[string]interface{}), nil
}

// Save saves configuration data (not applicable for defaults).
func (source *DefaultConfigSource) Save(path string, data map[string]interface{}) error {
	return fmt.Errorf("saving defaults not supported")
}

// Exists checks if defaults exist (always true).
func (source *DefaultConfigSource) Exists(path string) bool {
	return true
}

// GetSourceType returns the type of source.
func (source *DefaultConfigSource) GetSourceType() string {
	return "default"
}

// GetPriority returns the priority of this source for merging.
func (source *DefaultConfigSource) GetPriority() int {
	return 10 // Low priority
}

// 🔶 REFACTOR-003: Config abstraction - Generic configuration merger implementation - 🔍
// GenericConfigMerger implements ConfigMerger for schema-agnostic merging.
type GenericConfigMerger struct {
	strategy MergeStrategy
}

// NewGenericConfigMerger creates a new generic configuration merger.
func NewGenericConfigMerger(strategy MergeStrategy) *GenericConfigMerger {
	return &GenericConfigMerger{strategy: strategy}
}

// MergeConfigs merges two configuration objects.
func (merger *GenericConfigMerger) MergeConfigs(dst, src interface{}) error {
	// Use reflection to merge configurations
	dstValue := reflect.ValueOf(dst)
	srcValue := reflect.ValueOf(src)

	if dstValue.Kind() == reflect.Ptr {
		dstValue = dstValue.Elem()
	}
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	// Ensure both are structs
	if dstValue.Kind() != reflect.Struct || srcValue.Kind() != reflect.Struct {
		return fmt.Errorf("can only merge struct types")
	}

	// Iterate through all fields in the struct
	for i := 0; i < dstValue.NumField(); i++ {
		dstField := dstValue.Field(i)
		srcField := srcValue.Field(i)

		if err := merger.mergeValues(dstField, srcField); err != nil {
			return err
		}
	}

	return nil
}

// MergeWithPriority merges multiple configurations with priority ordering.
func (merger *GenericConfigMerger) MergeWithPriority(configs []interface{}, priorities []int) (interface{}, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("no configurations to merge")
	}

	if len(priorities) != len(configs) {
		return nil, fmt.Errorf("priorities length must match configs length")
	}

	// Sort configs by priority (highest first)
	type configPair struct {
		config   interface{}
		priority int
	}

	pairs := make([]configPair, len(configs))
	for i, config := range configs {
		pairs[i] = configPair{config: config, priority: priorities[i]}
	}

	// Simple sorting by priority (highest first)
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].priority > pairs[i].priority {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// Start with first (highest priority) config as base
	result := pairs[0].config

	// Create a copy to avoid modifying the original
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() == reflect.Ptr {
		// Create a new instance and copy values
		resultType := resultValue.Type().Elem()
		newResult := reflect.New(resultType)
		newResult.Elem().Set(resultValue.Elem())
		result = newResult.Interface()
	}

	// Save original strategy and temporarily use PreserveStrategy
	// to preserve high priority values when merging lower priority configs
	originalStrategy := merger.strategy
	merger.strategy = PreserveStrategy

	// Merge remaining configs into result (lower priority configs)
	// PreserveStrategy ensures high priority values are preserved
	for i := 1; i < len(pairs); i++ {
		if err := merger.MergeConfigs(result, pairs[i].config); err != nil {
			merger.strategy = originalStrategy // Restore original strategy
			return nil, err
		}
	}

	// Restore original strategy
	merger.strategy = originalStrategy

	return result, nil
}

// GetMergeStrategy returns the current merge strategy.
func (merger *GenericConfigMerger) GetMergeStrategy() MergeStrategy {
	return merger.strategy
}

// SetMergeStrategy sets the merge strategy.
func (merger *GenericConfigMerger) SetMergeStrategy(strategy MergeStrategy) {
	merger.strategy = strategy
}

// mergeValues merges source value into destination value using reflection.
func (merger *GenericConfigMerger) mergeValues(dst, src reflect.Value) error {
	if !dst.CanSet() {
		return nil // Skip read-only fields
	}

	switch merger.strategy {
	case OverwriteStrategy:
		// For OverwriteStrategy, we want to overwrite dst with src if src is not zero
		if !src.IsZero() {
			dst.Set(src)
		}
	case PreserveStrategy:
		// For PreserveStrategy, only set dst if it's zero and src is not zero
		if dst.IsZero() && !src.IsZero() {
			dst.Set(src)
		}
	case AppendStrategy:
		if dst.Kind() == reflect.Slice && src.Kind() == reflect.Slice {
			// Append slices
			newSlice := reflect.AppendSlice(dst, src)
			dst.Set(newSlice)
		} else if !src.IsZero() {
			dst.Set(src)
		}
	}

	return nil
}

// 🔶 REFACTOR-003: Backward compatibility - Configuration adapter implementation - 🔧
// ConfigAdapterImpl implements ConfigAdapter for backward compatibility.
type ConfigAdapterImpl struct {
	provider ConfigProvider
	schema   ConfigSchema
}

// NewConfigAdapter creates a new configuration adapter.
func NewConfigAdapter(provider ConfigProvider, schema ConfigSchema) *ConfigAdapterImpl {
	return &ConfigAdapterImpl{
		provider: provider,
		schema:   schema,
	}
}

// ToLegacyConfig converts to the legacy Config struct.
func (adapter *ConfigAdapterImpl) ToLegacyConfig() *Config {
	// Create new config with defaults
	cfg := DefaultConfig()

	// Override with provider values
	if adapter.provider.HasKey("archive_dir_path") {
		cfg.ArchiveDirPath = adapter.provider.GetString("archive_dir_path")
	}
	if adapter.provider.HasKey("use_current_dir_name") {
		cfg.UseCurrentDirName = adapter.provider.GetBool("use_current_dir_name")
	}
	if adapter.provider.HasKey("exclude_patterns") {
		cfg.ExcludePatterns = adapter.provider.GetStringSlice("exclude_patterns")
	}
	if adapter.provider.HasKey("include_git_info") {
		cfg.IncludeGitInfo = adapter.provider.GetBool("include_git_info")
	}
	if adapter.provider.HasKey("backup_dir_path") {
		cfg.BackupDirPath = adapter.provider.GetString("backup_dir_path")
	}

	// Add more field mappings as needed...

	return cfg
}

// FromLegacyConfig creates an adapter from a legacy Config struct.
func (adapter *ConfigAdapterImpl) FromLegacyConfig(cfg *Config) ConfigAdapter {
	provider := NewBackupConfigProvider(cfg, NewBackupConfigSchema())
	return NewConfigAdapter(provider, adapter.schema)
}

// GetProvider returns the underlying configuration provider.
func (adapter *ConfigAdapterImpl) GetProvider() ConfigProvider {
	return adapter.provider
}

// GetSchema returns the configuration schema.
func (adapter *ConfigAdapterImpl) GetSchema() ConfigSchema {
	return adapter.schema
}

// 🔶 REFACTOR-003: Backward compatibility - Legacy configuration loading wrapper - 🔍
// LoadConfigLegacy provides backward compatible configuration loading.
func LoadConfigLegacy(root string) (*Config, error) {
	loader := NewGenericConfigLoader()
	schema := NewBackupConfigSchema()

	configInterface, err := loader.LoadConfigWithSchema(root, schema)
	if err != nil {
		return nil, err
	}

	// Convert to legacy Config struct
	if cfg, ok := configInterface.(*Config); ok {
		return cfg, nil
	}

	return nil, fmt.Errorf("unexpected configuration type: %T", configInterface)
}
