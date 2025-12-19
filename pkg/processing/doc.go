// ARCH-001: See architecture.md - Core Architecture [DECISION:maintenance]
// Package processing provides generalized data processing patterns extracted from BkpDir.
//
// This package contains reusable patterns for:
//   - Timestamp-based naming conventions with metadata integration
//   - Processing pipelines with context support and atomic operations
//   - Concurrent processing with worker pools and resource management
//
// The package is designed to be used independently or in combination with other
// extracted BkpDir packages like pkg/config, pkg/errors, and pkg/formatter.
//
// Example usage:
//
//	// Create a naming provider for timestamp-based names
//	naming := processing.NewNamingProvider()
//	name, err := naming.GenerateName(processing.NamingTemplate{
//		Prefix:    "backup",
//		Timestamp: time.Now(),
//		Metadata:  map[string]string{"branch": "main", "note": "initial"},
//	})
//
//
//	// Create a processing pipeline
//	pipeline := processing.NewPipeline()
//	pipeline.AddStage(processing.CollectionStage{})
//	pipeline.AddStage(processing.ProcessingStage{})
//	result, err := pipeline.Execute(ctx, input)
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
//
// [REQ:PERFORMANCE] Processing patterns must support high-performance concurrent execution
// [ARCH:PROCESSING_PATTERNS] Processing pipeline architecture with worker pools and naming providers
// [IMPL:PROCESSING_PATTERNS] Implementation of pipeline stages, naming providers, and concurrent processors
package processing
