// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// - [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: simulate disk space, permissions, corruption, and error injection in internal/testutil for stress paths.
// Package testutil provides common testing utilities and patterns for CLI applications.
//
// This package extracts reusable testing patterns from the bkpdir application into
// a generic testing framework that can be used by any Go CLI application.
//
// # Key Components
//
// The package provides the following main categories of testing utilities:
//
//   - Configuration Testing: Helpers for creating test configurations and managing
//     environment variable isolation
//   - File System Testing: Utilities for creating temporary files, directories, and
//     ZIP archives for testing purposes
//   - CLI Testing: Helpers for testing command-line interfaces and cobra commands
//   - Test Assertions: Common assertion functions for various data types
//   - Test Fixtures: Reusable test data and setup patterns
//
// # Usage Examples
//
// Configuration Testing:
//
//	config, cleanup := testutil.CreateTestConfig(t, map[string]interface{}{
//		"setting1": "value1",
//		"setting2": true,
//	})
//	defer cleanup()
//
// File System Testing:
//
//	files := map[string]string{
//		"file1.txt": "content1",
//		"dir/file2.txt": "content2",
//	}
//	archivePath := testutil.CreateTestZipArchive(t, "test.zip", files)
//	defer os.Remove(archivePath)
//
// CLI Testing:
//
//	cmd := testutil.CreateTestCommand("myapp", func(cmd *cobra.Command, args []string) {
//		// command implementation
//	})
//	testutil.ExecuteCommand(t, cmd, []string{"arg1", "arg2"})
//
// Test Assertions:
//
//	testutil.AssertStringEqual(t, "field name", got, want)
//	testutil.AssertSliceEqual(t, "slice field", gotSlice, wantSlice)
//
// # Design Principles
//
// The testutil package follows these design principles:
//
//   - Generic and Reusable: All utilities are designed to work with any Go application,
//     not just the original bkpdir application
//   - Zero Dependencies: Minimal external dependencies to ensure broad compatibility
//   - Interface-Based: Uses interfaces where possible for maximum flexibility
//   - Composable: Utilities can be combined for complex testing scenarios
//   - Clean and Simple: Focuses on common patterns rather than specialized functionality
//
// # Integration with testing.T
//
// All utilities are designed to work seamlessly with Go's standard testing framework:
//
//   - Functions accept *testing.T as the first parameter
//   - Use t.Helper() to ensure proper test failure attribution
//   - Follow Go testing conventions and best practices
//   - Provide clear, actionable error messages
//
// # Relationship to internal/testutil
//
// This package complements the existing internal/testutil package:
//
//   - pkg/testutil: Common testing patterns suitable for external use
//   - internal/testutil: Specialized testing infrastructure for complex scenarios
//
// Applications can use both packages together for comprehensive testing capabilities.
//
// - [IMPL-TOKEN_COVERAGE_AUDIT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] — How: scan each module file for REQ/ARCH/IMPL comments and record missing expected tokens.
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: parse legacy token YAML groups and emit normalized JSON with descriptions, status, and source paths per category.
// - [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-GOV_REGISTRY_COMPLETENESS] [REQ-IMMUTABLE_DIRECTORY_OPERATIONS] — How: resolve canonical REQ/ARCH/IMPL references for each legacy token and flag gaps when links are missing.
// - [IMPL-TRACEABILITY] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_003] — How: hash interface signatures plus documented behavioral contracts to assign stable feature fingerprints linked to semantic tokens.
// - [IMPL-TESTING_COMPLEXITY] [ARCH-TESTING_STRATEGY] [REQ-RELIABILITY] — How: extract shared test utilities first, add package-focused tests, and keep root integration tests green after package extraction.
package testutil
