// Package config provides inheritance chain building and path resolution.
//
// This file implements the core inheritance functionality that builds dependency
// chains, resolves file paths, and detects circular dependencies in configuration
// inheritance hierarchies.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ⭐ CFG-005: Inheritance chain builder implementation - 🔧 Core inheritance functionality

// DefaultInheritanceChainBuilder implements InheritanceChainBuilder interface.
type DefaultInheritanceChainBuilder struct {
	fileOps ConfigFileOperations
}

// NewInheritanceChainBuilder creates a new inheritance chain builder.
func NewInheritanceChainBuilder(fileOps ConfigFileOperations) *DefaultInheritanceChainBuilder {
	return &DefaultInheritanceChainBuilder{
		fileOps: fileOps,
	}
}

// ⭐ CFG-005: Inheritance chain building logic - 🔍 Dependency chain construction
func (b *DefaultInheritanceChainBuilder) BuildChain(configPath string, pathResolver PathResolver) (*InheritanceChain, error) {
	startTime := time.Now()

	chain := &InheritanceChain{
		Files:        make([]string, 0),
		Visited:      make(map[string]bool),
		Sources:      make(map[string]string),
		Dependencies: make(map[string][]string),
	}

	// Build the dependency chain recursively
	err := b.buildChainRecursive(configPath, "", pathResolver, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to build inheritance chain from %s: %w", configPath, err)
	}

	// Record build time for performance tracking
	buildTime := time.Since(startTime).Nanoseconds()

	// Add metadata
	chain.Sources["build_time"] = fmt.Sprintf("%d", buildTime)
	chain.Sources["root_file"] = configPath

	return chain, nil
}

// ⭐ CFG-005: Recursive chain building - 🔍 Depth-first dependency resolution
func (b *DefaultInheritanceChainBuilder) buildChainRecursive(configPath, basePath string, pathResolver PathResolver, chain *InheritanceChain) error {
	// Resolve the full path
	resolvedPath, err := pathResolver.ResolvePath(configPath, basePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", configPath, err)
	}

	// Check for circular dependency
	if chain.Visited[resolvedPath] {
		return fmt.Errorf("circular dependency detected: %s already visited", resolvedPath)
	}

	// Validate path exists and is accessible
	if err := pathResolver.ValidatePath(resolvedPath); err != nil {
		return fmt.Errorf("invalid path %s: %w", resolvedPath, err)
	}

	// Mark as visited for circular dependency detection
	chain.Visited[resolvedPath] = true

	// Load inheritance metadata from the configuration file
	inheritance, err := b.loadInheritanceMetadata(resolvedPath)
	if err != nil {
		return fmt.Errorf("failed to load inheritance metadata from %s: %w", resolvedPath, err)
	}

	// Record dependencies
	chain.Dependencies[resolvedPath] = inheritance.Inherit

	// Process parent files first (bottom-up dependency resolution)
	for _, parentPath := range inheritance.Inherit {
		err := b.buildChainRecursive(parentPath, filepath.Dir(resolvedPath), pathResolver, chain)
		if err != nil {
			return fmt.Errorf("failed to process parent %s from %s: %w", parentPath, resolvedPath, err)
		}
	}

	// Add current file to chain (after dependencies)
	chain.Files = append(chain.Files, resolvedPath)
	chain.Sources[resolvedPath] = fmt.Sprintf("inheritance:%d", len(chain.Files))

	return nil
}

// ⭐ CFG-005: Inheritance metadata loading - 📝 YAML inheritance parsing
func (b *DefaultInheritanceChainBuilder) loadInheritanceMetadata(configPath string) (*ConfigInheritance, error) {
	// Read the configuration file
	data, err := b.fileOps.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse only the inheritance metadata
	var metadata struct {
		ConfigInheritance `yaml:",inline"`
	}

	err = yaml.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inheritance metadata from %s: %w", configPath, err)
	}

	return &metadata.ConfigInheritance, nil
}

// ValidateChain validates an inheritance chain for consistency and circular dependencies.
func (b *DefaultInheritanceChainBuilder) ValidateChain(chain *InheritanceChain) error {
	if chain == nil {
		return fmt.Errorf("inheritance chain cannot be nil")
	}

	// Check for empty chain
	if len(chain.Files) == 0 {
		return fmt.Errorf("inheritance chain cannot be empty")
	}

	// Validate each file in the chain exists
	for _, filePath := range chain.Files {
		if !b.fileOps.FileExists(filePath) {
			return fmt.Errorf("file in inheritance chain does not exist: %s", filePath)
		}
	}

	// Validate dependencies are consistent
	for file, deps := range chain.Dependencies {
		for _, dep := range deps {
			found := false
			for _, chainFile := range chain.Files {
				if strings.HasSuffix(chainFile, dep) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("dependency %s of file %s not found in chain", dep, file)
			}
		}
	}

	return nil
}

// GetChainMetadata returns metadata about an inheritance chain.
func (b *DefaultInheritanceChainBuilder) GetChainMetadata(chain *InheritanceChain) *InheritanceChainMetadata {
	if chain == nil {
		return &InheritanceChainMetadata{}
	}

	// Calculate maximum depth
	maxDepth := 0
	for _, deps := range chain.Dependencies {
		if len(deps) > maxDepth {
			maxDepth = len(deps)
		}
	}

	// Extract build time if available
	var buildTime int64 = 0
	if buildTimeStr, exists := chain.Sources["build_time"]; exists {
		fmt.Sscanf(buildTimeStr, "%d", &buildTime)
	}

	return &InheritanceChainMetadata{
		ChainLength:  len(chain.Files),
		MaxDepth:     maxDepth,
		SourceFiles:  append([]string(nil), chain.Files...), // Copy slice
		Dependencies: chain.Dependencies,
		BuildTime:    buildTime,
	}
}

// ⭐ CFG-005: Path resolver implementation - 🔍 Path resolution and validation

// DefaultPathResolver implements PathResolver interface.
type DefaultPathResolver struct {
	fileOps ConfigFileOperations
}

// NewPathResolver creates a new path resolver.
func NewPathResolver(fileOps ConfigFileOperations) *DefaultPathResolver {
	return &DefaultPathResolver{
		fileOps: fileOps,
	}
}

// ⭐ CFG-005: Path resolution logic - 🔧 Path expansion and normalization
func (r *DefaultPathResolver) ResolvePath(path string, basePath string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Expand path variables
	expandedPath, err := r.ExpandPath(path)
	if err != nil {
		return "", fmt.Errorf("failed to expand path %s: %w", path, err)
	}

	// Handle absolute paths
	if filepath.IsAbs(expandedPath) {
		return filepath.Clean(expandedPath), nil
	}

	// Handle relative paths
	if basePath != "" {
		resolvedPath := filepath.Join(basePath, expandedPath)
		return filepath.Clean(resolvedPath), nil
	}

	// Use current working directory as base
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	resolvedPath := filepath.Join(cwd, expandedPath)
	return filepath.Clean(resolvedPath), nil
}

// ⭐ CFG-005: Path expansion logic - 🔧 Variable expansion support
func (r *DefaultPathResolver) ExpandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Handle home directory expansion
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}

	// Handle environment variable expansion
	expandedPath := os.ExpandEnv(path)

	return expandedPath, nil
}

// ValidatePath validates that a configuration file path is safe and accessible.
func (r *DefaultPathResolver) ValidatePath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	// Check if file exists
	if !r.fileOps.FileExists(path) {
		return fmt.Errorf("configuration file does not exist: %s", path)
	}

	// Check if file is readable
	_, err := r.fileOps.ReadFile(path)
	if err != nil {
		return fmt.Errorf("configuration file is not readable: %s: %w", path, err)
	}

	// Validate file extension (optional - could be enforced)
	ext := filepath.Ext(path)
	validExtensions := []string{".yml", ".yaml", ".json"}

	validExt := false
	for _, validExtension := range validExtensions {
		if ext == validExtension {
			validExt = true
			break
		}
	}

	if !validExt {
		return fmt.Errorf("invalid configuration file extension %s (expected: %v)", ext, validExtensions)
	}

	return nil
}

// ⭐ CFG-005: Circular dependency detector implementation - 🛡️ Cycle detection

// DefaultCircularDependencyDetector implements CircularDependencyDetector interface.
type DefaultCircularDependencyDetector struct {
	visited map[string]bool
	inStack map[string]bool
	path    []string
}

// NewCircularDependencyDetector creates a new circular dependency detector.
func NewCircularDependencyDetector() *DefaultCircularDependencyDetector {
	return &DefaultCircularDependencyDetector{
		visited: make(map[string]bool),
		inStack: make(map[string]bool),
		path:    make([]string, 0),
	}
}

// ⭐ CFG-005: Cycle detection algorithm - 🔍 Depth-first search implementation
func (d *DefaultCircularDependencyDetector) DetectCycle(startFile string, resolver PathResolver) error {
	// Reset detector state
	d.Reset()

	// Start cycle detection from the given file
	return d.detectCycleRecursive(startFile, "", resolver)
}

// detectCycleRecursive performs depth-first search for cycle detection.
func (d *DefaultCircularDependencyDetector) detectCycleRecursive(configPath, basePath string, resolver PathResolver) error {
	// Resolve the path
	resolvedPath, err := resolver.ResolvePath(configPath, basePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path %s: %w", configPath, err)
	}

	// Check if we're already processing this file (cycle detected)
	if d.inStack[resolvedPath] {
		// Add current file to show the cycle
		d.path = append(d.path, resolvedPath)
		return fmt.Errorf("circular dependency detected: %s", strings.Join(d.path, " -> "))
	}

	// Skip if already visited and not in current stack
	if d.visited[resolvedPath] {
		return nil
	}

	// Mark as visited and in current stack
	d.visited[resolvedPath] = true
	d.inStack[resolvedPath] = true
	d.path = append(d.path, resolvedPath)

	// Load dependencies and check recursively
	inheritance, err := d.loadInheritanceMetadata(resolvedPath, resolver)
	if err != nil {
		return fmt.Errorf("failed to load inheritance metadata from %s: %w", resolvedPath, err)
	}

	for _, parentPath := range inheritance.Inherit {
		err := d.detectCycleRecursive(parentPath, filepath.Dir(resolvedPath), resolver)
		if err != nil {
			return err
		}
	}

	// Remove from stack after processing
	d.inStack[resolvedPath] = false
	if len(d.path) > 0 {
		d.path = d.path[:len(d.path)-1]
	}

	return nil
}

// loadInheritanceMetadata loads inheritance metadata for cycle detection.
func (d *DefaultCircularDependencyDetector) loadInheritanceMetadata(configPath string, resolver PathResolver) (*ConfigInheritance, error) {
	// Validate path
	if err := resolver.ValidatePath(configPath); err != nil {
		return nil, err
	}

	// Read file content
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Parse inheritance metadata
	var metadata struct {
		ConfigInheritance `yaml:",inline"`
	}

	err = yaml.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inheritance metadata from %s: %w", configPath, err)
	}

	return &metadata.ConfigInheritance, nil
}

// GetCyclePath returns the path of detected circular dependency.
func (d *DefaultCircularDependencyDetector) GetCyclePath() []string {
	return append([]string(nil), d.path...) // Return copy
}

// Reset resets the detector state for new cycle detection.
func (d *DefaultCircularDependencyDetector) Reset() {
	d.visited = make(map[string]bool)
	d.inStack = make(map[string]bool)
	d.path = make([]string, 0)
}
