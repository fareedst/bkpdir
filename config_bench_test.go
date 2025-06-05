// This file is part of bkpdir
//
// Package main provides performance benchmarks for configuration reflection and visibility.
// It validates the performance optimizations implemented in CFG-006.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"testing"
	"time"
)

// 🔶 CFG-006: Performance optimization - Benchmark validation
// IMPLEMENTATION-REF: CFG-006 Subtask 6.4: Add benchmark validation

// BenchmarkGetAllConfigFields benchmarks field discovery performance.
// Target: <50ms for cache miss, <10ms for cache hit.
func BenchmarkGetAllConfigFields(b *testing.B) {
	cfg := DefaultConfig()

	// Clear cache to ensure we're measuring cold performance
	globalFieldCache.invalidateCache()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fields := GetAllConfigFields(cfg)
		if len(fields) == 0 {
			b.Fatal("Expected fields, got none")
		}
	}
}

// BenchmarkGetAllConfigFieldsCached benchmarks cached field discovery performance.
// This measures the cache hit performance after warming up the cache.
func BenchmarkGetAllConfigFieldsCached(b *testing.B) {
	cfg := DefaultConfig()

	// Warm up the cache
	GetAllConfigFields(cfg)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fields := GetAllConfigFields(cfg)
		if len(fields) == 0 {
			b.Fatal("Expected fields, got none")
		}
	}
}

// BenchmarkGetConfigValuesWithSources benchmarks full source resolution performance.
// Target: <100ms for full resolution.
func BenchmarkGetConfigValuesWithSources(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		values := GetAllConfigValuesWithSources(cfg, ".")
		if len(values) == 0 {
			b.Fatal("Expected values, got none")
		}
	}
}

// BenchmarkGetConfigValuesWithSourcesFiltered benchmarks filtered source resolution.
// Target: 50%+ improvement over unfiltered when filtering enabled.
func BenchmarkGetConfigValuesWithSourcesFiltered(b *testing.B) {
	cfg := DefaultConfig()
	filter := &ConfigFilter{
		Categories: []string{"basic_settings"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		values := GetConfigValuesWithSourcesFiltered(cfg, ".", filter)
		if len(values) == 0 {
			b.Fatal("Expected values, got none")
		}
	}
}

// BenchmarkGetConfigValuesOverridesOnly benchmarks overrides-only filtering.
// This tests the lazy evaluation of source resolution.
func BenchmarkGetConfigValuesOverridesOnly(b *testing.B) {
	cfg := DefaultConfig()
	// Create multiple obvious overrides that should be detected
	cfg.ArchiveDirPath = "/definitely/custom/override/path" // Override default "../.bkpdir"
	cfg.UseCurrentDirName = false                           // Override default true
	cfg.IncludeGitInfo = true                               // Override default false

	filter := &ConfigFilter{
		OverridesOnly: true,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		values := GetConfigValuesWithSourcesFiltered(cfg, ".", filter)
		// The exact number will depend on whether there's a config file present
		// But we should always get some values since we modified the config
		if len(values) == 0 {
			// Log the values we got for debugging
			var sources []string
			for _, v := range values {
				sources = append(sources, fmt.Sprintf("%s=%s(source:%s)", v.Name, v.Value, v.Source))
			}
			b.Fatalf("Expected some override values, got %d: %v", len(values), sources)
		}
	}
}

// BenchmarkGetConfigFieldByPattern benchmarks pattern-based field retrieval.
// Target: Sub-10ms response for single field queries.
func BenchmarkGetConfigFieldByPattern(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fields, err := GetConfigFieldByPattern(cfg, "ArchiveDirPath")
		if err != nil {
			b.Fatal("Pattern search failed:", err)
		}
		if len(fields) == 0 {
			b.Fatal("Expected fields, got none")
		}
	}
}

// BenchmarkGetConfigFieldValue benchmarks single field value retrieval.
// Target: Sub-10ms response for single field access.
func BenchmarkGetConfigFieldValue(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		value, err := GetConfigFieldValue(cfg, "ArchiveDirPath")
		if err != nil {
			b.Fatal("Field value retrieval failed:", err)
		}
		if value.Value == "" {
			b.Fatal("Expected field value, got empty")
		}
	}
}

// BenchmarkHasConfigField benchmarks field existence checking.
// This should be very fast since it only checks existence.
func BenchmarkHasConfigField(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		exists := HasConfigField(cfg, "ArchiveDirPath")
		if !exists {
			b.Fatal("Expected field to exist")
		}
	}
}

// BenchmarkConfigCommandResponse benchmarks end-to-end config command performance.
// This simulates the actual config command usage pattern.
func BenchmarkConfigCommandResponse(b *testing.B) {
	cfg := DefaultConfig()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Simulate config command: get all values with sources
		start := time.Now()
		values := GetAllConfigValuesWithSources(cfg, ".")
		duration := time.Since(start)

		if len(values) == 0 {
			b.Fatal("Expected values, got none")
		}

		// Verify performance target: <100ms
		if duration > 100*time.Millisecond {
			b.Logf("Warning: Config command took %v, target is <100ms", duration)
		}
	}
}

// BenchmarkConfigFieldCacheValidation benchmarks cache validation performance.
// This measures the overhead of cache validation logic.
func BenchmarkConfigFieldCacheValidation(b *testing.B) {
	cfg := DefaultConfig()

	// Warm up cache
	GetAllConfigFields(cfg)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// This will hit cache validation logic
		_ = globalFieldCache.getCachedFields()
	}
}

// BenchmarkStructHashComputation benchmarks struct hash computation.
// This measures the cost of cache invalidation detection.
func BenchmarkStructHashComputation(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		hash := getConfigStructHash()
		if hash == 0 {
			b.Fatal("Expected non-zero hash")
		}
	}
}

// TestPerformanceTargets validates that performance targets are met.
// This is a regular test that measures performance and fails if targets aren't met.
func TestPerformanceTargets(t *testing.T) {
	cfg := DefaultConfig()

	// Test 1: Field discovery cache miss should be <50ms
	globalFieldCache.invalidateCache()
	start := time.Now()
	fields := GetAllConfigFields(cfg)
	cacheMissDuration := time.Since(start)

	if len(fields) == 0 {
		t.Fatal("Expected fields, got none")
	}

	if cacheMissDuration > 50*time.Millisecond {
		t.Errorf("Field discovery cache miss took %v, target is <50ms", cacheMissDuration)
	}

	// Test 2: Field discovery cache hit should be <10ms
	start = time.Now()
	fields = GetAllConfigFields(cfg)
	cacheHitDuration := time.Since(start)

	if cacheHitDuration > 10*time.Millisecond {
		t.Errorf("Field discovery cache hit took %v, target is <10ms", cacheHitDuration)
	}

	// Test 3: Full config command should be <100ms
	start = time.Now()
	values := GetAllConfigValuesWithSources(cfg, ".")
	configCommandDuration := time.Since(start)

	if len(values) == 0 {
		t.Fatal("Expected values, got none")
	}

	if configCommandDuration > 100*time.Millisecond {
		t.Errorf("Config command took %v, target is <100ms", configCommandDuration)
	}

	// Test 4: Single field access should be <10ms
	start = time.Now()
	_, err := GetConfigFieldValue(cfg, "ArchiveDirPath")
	singleFieldDuration := time.Since(start)

	if err != nil {
		t.Fatal("Single field access failed:", err)
	}

	if singleFieldDuration > 10*time.Millisecond {
		t.Errorf("Single field access took %v, target is <10ms", singleFieldDuration)
	}

	t.Logf("Performance targets met:")
	t.Logf("  Cache miss: %v (<50ms)", cacheMissDuration)
	t.Logf("  Cache hit: %v (<10ms)", cacheHitDuration)
	t.Logf("  Config command: %v (<100ms)", configCommandDuration)
	t.Logf("  Single field: %v (<10ms)", singleFieldDuration)
}

// TestGetConfigValuesOverridesOnlyDebug debugs the overrides-only filtering.
// This helps understand why the benchmark is failing.
func TestGetConfigValuesOverridesOnlyDebug(t *testing.T) {
	cfg := DefaultConfig()
	defaultCfg := DefaultConfig()

	// Show original defaults
	t.Logf("Default ArchiveDirPath: %q", defaultCfg.ArchiveDirPath)
	t.Logf("Default UseCurrentDirName: %v", defaultCfg.UseCurrentDirName)
	t.Logf("Default IncludeGitInfo: %v", defaultCfg.IncludeGitInfo)

	// Create obvious overrides
	cfg.ArchiveDirPath = "/definitely/custom/override/path"
	cfg.UseCurrentDirName = false
	cfg.IncludeGitInfo = true

	// Show modified values
	t.Logf("Modified ArchiveDirPath: %q", cfg.ArchiveDirPath)
	t.Logf("Modified UseCurrentDirName: %v", cfg.UseCurrentDirName)
	t.Logf("Modified IncludeGitInfo: %v", cfg.IncludeGitInfo)

	// Get all config values first to see sources
	allValues := GetConfigValuesWithSourcesFiltered(cfg, ".", nil)
	t.Logf("Total config values found: %d", len(allValues))

	// Check specific fields
	for _, val := range allValues {
		if val.Name == "archive_dir_path" || val.Name == "use_current_dir_name" || val.Name == "include_git_info" {
			t.Logf("%s: value=%q source=%s isOverridden=%v", val.Name, val.Value, val.Source, val.IsOverridden)
		}
	}

	// Now test overrides-only filter
	filter := &ConfigFilter{
		OverridesOnly: true,
	}

	values := GetConfigValuesWithSourcesFiltered(cfg, ".", filter)
	t.Logf("Override-only values found: %d", len(values))

	for _, val := range values {
		t.Logf("Override: %s=%s (source: %s)", val.Name, val.Value, val.Source)
	}
}
