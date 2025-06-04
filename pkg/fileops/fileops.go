// Package fileops provides file operations and utilities for CLI applications.
//
// This package extracts common file operation patterns including:
// - File and directory comparison with hash-based content verification
// - Pattern-based file exclusion using doublestar glob patterns
// - Path validation with security and existence checking
// - Atomic file operations for safe writing with rollback capabilities
// - Directory traversal with configurable exclusion and security options
//
// The package is designed to be reusable across different CLI applications
// while maintaining the security and reliability patterns established in
// the original BkpDir application.
//
// Key Components:
//
// Comparison: File and directory comparison utilities
//   - CreateDirectorySnapshot() - Creates snapshots for comparison
//   - CreateArchiveSnapshot() - Creates snapshots from ZIP archives
//   - CompareSnapshots() - Compares two snapshots for differences
//   - IsDirectoryIdenticalToArchive() - Checks directory-archive identity
//
// Exclusion: Pattern-based file exclusion system
//   - PatternMatcher - Handles doublestar glob pattern matching
//   - ShouldExcludeFile() - Main exclusion interface
//   - Directory, glob, and root-level pattern matching support
//
// Validation: Path validation and security checking
//   - ValidatePath() - Comprehensive path validation
//   - ValidateExistence() - Path existence checking
//   - ValidateReadable()/ValidateWritable() - Permission validation
//   - IsSecurePath() - Security validation (path traversal, etc.)
//
// Atomic Operations: Safe file writing with rollback
//   - AtomicWriter - Atomic file writing with temporary files
//   - AtomicCopy() - Atomic file copying
//   - AtomicWriteFile() - Atomic file creation
//   - Automatic cleanup on errors or rollback
//
// Traversal: Safe directory walking with exclusions
//   - Walk() - Basic directory traversal
//   - WalkWithOptions() - Configurable traversal with exclusions
//   - ListFiles() - File listing with optional recursion
//   - Support for symlinks, depth limits, hidden files
//
// Example Usage:
//
//	// File comparison
//	identical, err := fileops.IsDirectoryIdenticalToArchive(
//	    "/path/to/dir", "/path/to/archive.zip", []string{".git/", "*.tmp"})
//
//	// Safe file writing
//	err := fileops.AtomicWriteFile("/path/to/file", data, 0644)
//
//	// Directory traversal with exclusions
//	err := fileops.WalkWithExclusions("/path/to/dir", []string{"*.log"},
//	    func(path string, info os.FileInfo, err error) error {
//	        fmt.Println(path)
//	        return nil
//	    })
//
//	// Pattern-based exclusion
//	shouldExclude := fileops.ShouldExcludeFile("logs/debug.log", []string{"*.log"})
//
// All operations maintain the security and reliability patterns from the
// original implementation while providing clean, reusable interfaces.
package fileops

// ⭐ EXTRACT-006: File operations package interface - 📝

// Version information
const (
	// Version of the fileops package
	Version = "1.0.0"

	// PackageName for identification
	PackageName = "fileops"
)

// ⭐ EXTRACT-006: Package version and identification - 📝
