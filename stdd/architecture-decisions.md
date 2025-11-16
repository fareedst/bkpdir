# Architecture Decisions

**STDD Methodology Version**: 1.0.0

## Overview
This document captures the high-level architectural decisions for the bkpdir project. All decisions are cross-referenced with requirements using semantic tokens `[REQ:*]` and assigned architecture tokens `[ARCH:*]` for traceability.

## 1. Language and Runtime [ARCH:LANGUAGE_SELECTION]

### Decision: Go (Golang)
**Rationale:**
- Excellent standard library support for file operations, ZIP archives, and CLI tools
- Strong concurrency support with goroutines and channels
- Cross-platform compilation to single binaries
- Excellent tooling (go fmt, go test, go vet)
- Fast compilation and execution

**Alternatives Considered:**
- Python: Rejected due to runtime dependencies and slower performance
- Rust: Rejected due to complexity and longer development time
- C/C++: Rejected due to manual memory management complexity

**Cross-References**: [REQ:PERFORMANCE], [REQ:MAINTAINABILITY]

## 2. Project Structure [ARCH:PROJECT_STRUCTURE]

### Decision: Flat structure with internal packages
```
bkpdir/
├── main.go              # Main application entry point
├── archive.go           # Archive creation and management
├── backup.go            # File backup operations
├── config.go            # Configuration management
├── formatter.go         # Output formatting
├── errors.go            # Error handling
├── comparison.go        # Directory comparison
├── exclude.go           # File exclusion patterns
├── git.go              # Git integration
├── internal/           # Internal packages
│   ├── testutil/       # Test utilities
│   └── validation/     # Validation utilities
├── pkg/                # Extracted packages (future)
└── docs/               # Documentation
```

**Rationale:**
- Clear separation of concerns
- Standard Go project layout
- Testable components
- Supports future package extraction

**Cross-References**: [REQ:MAINTAINABILITY]

## 2.1. System Components Architecture [ARCH:SYSTEM_COMPONENTS]

### Decision: Layered architecture with clear component boundaries

**System Components Diagram:**
```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                            │
├─────────────────────────────────────────────────────────────┤
│  Command Handlers  │  Flag Processing  │  Output Formatting │
│  Context Support   │  Error Handling   │  Resource Cleanup  │
│  Dual-Mode Ops     │  Backward Compat  │  File Operations   │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                    Configuration Layer                      │
├─────────────────────────────────────────────────────────────┤
│  Config Discovery  │  Environment Vars │  Status Codes      │
│  Template Configs  │  Format Strings   │  Regex Patterns    │
│  Dual-Mode Config  │  File Backup Cfg  │  Archive Config    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                      Core Services                          │
├─────────────────────────────────────────────────────────────┤
│  Archive Service   │  Git Service      │  Verification      │
│  Resource Manager  │  Template Engine  │  Error Handler     │
│  Context Manager   │  File Operations  │  Backup Service    │
│  Comparison Svc    │  Formatter Svc    │  Config Service    │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                     Storage Layer                           │
├─────────────────────────────────────────────────────────────┤
│  File System        │  ZIP Archives    │  Checksums         │
│  File Backups       │  Metadata        │  Verification      │
│  Directory Trees    │  Atomic Ops      │  Resource Cleanup  │
└─────────────────────────────────────────────────────────────┘
```

**Rationale:**
- Clear separation of concerns across layers
- CLI layer handles user interaction and command routing
- Configuration layer manages all configuration sources and merging
- Core services provide business logic and orchestration
- Storage layer handles all file system and archive operations
- Enables testability and maintainability

**Cross-References**: [REQ:MAINTAINABILITY], [REQ:USABILITY]

## 3. Archive Format [ARCH:ARCHIVE_FORMAT] [REQ:ARCHIVE_VERIFICATION]

### Decision: ZIP format for all archive operations
**Rationale:**
- Cross-platform compatibility
- Built-in compression
- Wide tool support
- Standard library support in Go (`archive/zip`)
- Supports incremental archives

**Alternatives Considered:**
- tar.gz: Rejected due to less universal support
- tar.bz2: Rejected due to slower compression
- Custom format: Rejected due to compatibility concerns

**Implementation:**
- Uses Go's `archive/zip` package
- Supports compression levels
- Atomic file operations prevent corruption

**Cross-References**: [REQ:RELIABILITY], [REQ:ARCHIVE_VERIFICATION]

## 4. Configuration System [ARCH:CONFIG_SYSTEM] [REQ:CONFIGURATION]

### Decision: YAML configuration files with environment variable support
**Rationale:**
- Human-readable format
- Supports complex nested structures
- Good Go library support (`gopkg.in/yaml.v3`)
- Industry standard for configuration
- Environment variable support for flexibility

**Configuration Discovery:**
1. Command line flags (highest priority)
2. Environment variables (`BKPDIR_CONFIG` for search paths)
3. Configuration files in search path order:
   - `./.bkpdir.yml` (current directory)
   - `~/.bkpdir.yml` (home directory)
   - Custom paths from `BKPDIR_CONFIG`
4. Default values (lowest priority)

**Alternatives Considered:**
- JSON: Rejected due to lack of comments
- TOML: Rejected due to less common usage
- INI: Rejected due to limited nesting support

**Cross-References**: [REQ:CONFIGURATION], [REQ:USABILITY]

## 5. Layered Configuration Inheritance [ARCH:CFG_005] [REQ:CFG_005]

### Decision: Configuration files support explicit inheritance with merge strategies
**Rationale:**
- Enables hierarchical configuration management
- Supports complex project structures
- Flexible merge strategies for different use cases
- Maintains backward compatibility

**Merge Strategies:**
- Standard override (no prefix): Replace parent values (for non-array fields)
- Array merge (no prefix for array fields): Default merge strategy for array fields per CFG-005
- Array merge (`+` prefix): Explicit append to parent arrays
- Array prepend (`^` prefix): Prepend to parent arrays
- Array replace (`!` prefix): Replace parent arrays completely
- Default value (`=` prefix): Use only if parent not set

**Array Field Default Behavior:**
- Array fields (like `exclude_patterns`) default to merge (accumulate) strategy when no prefix is specified
- This applies to both inheritance chains and sequential file processing
- Explicit prefixes (`!`, `+`, `^`, `=`) override the default merge behavior
- Default patterns are preserved when child configs add values

**Implementation:**
- Inheritance chain processing with circular dependency detection
- Source tracking for debugging
- Performance optimization with caching
- Thread-safe configuration loading

**Cross-References**: [REQ:CFG_005], [REQ:PERFORMANCE]

## 6. Error Handling Strategy [ARCH:ERROR_HANDLING] [REQ:ERROR_HANDLING]

### Decision: Structured error types with status codes and operation context
**Rationale:**
- Consistent error handling across operations
- Machine-readable status codes for scripting
- Enhanced debugging with operation context
- Supports template-based error formatting

**Error Types:**
- `ArchiveError`: Structured errors for archive operations
- `BackupError`: Structured errors for backup operations
- Both include: Message, StatusCode, Operation, Path, Err

**Error Detection:**
- `isDiskFullError()`: Comprehensive disk space error detection
- Enhanced error messages with operation context
- Panic recovery mechanisms

**Cross-References**: [REQ:ERROR_HANDLING], [REQ:RELIABILITY]

## 7. Resource Management [ARCH:RESOURCE_MANAGEMENT] [REQ:RESOURCE_MANAGEMENT]

### Decision: ResourceManager for automatic cleanup with panic recovery
**Rationale:**
- Prevents resource leaks on operation failure
- Handles panic scenarios gracefully
- Thread-safe for concurrent operations
- Simplifies error handling in operations

**Implementation:**
- `ResourceManager` struct with mutex protection
- Methods: `AddTempFile()`, `AddTempDir()`, `RemoveResource()`, `CleanupWithPanicRecovery()`
- Thread-safe resource tracking
- Error-resilient cleanup operations

**Cross-References**: [REQ:RESOURCE_MANAGEMENT], [REQ:RELIABILITY]

## 8. Context-Aware Operations [ARCH:CONTEXT_SUPPORT] [REQ:CONTEXT_SUPPORT]

### Decision: Support context cancellation for long-running operations
**Rationale:**
- Enables operation timeouts
- Supports graceful cancellation
- Better user experience for long operations
- Standard Go pattern for cancellable operations

**Implementation:**
- All long-running operations accept `context.Context`
- Periodic cancellation checks during operations
- Context-aware file operations
- Timeout support via `context.WithTimeout()`

**Cross-References**: [REQ:CONTEXT_SUPPORT], [REQ:USABILITY]

## 9. Output Formatting [ARCH:OUTPUT_FORMATTING] [REQ:OUTPUT_FORMATTING]

### Decision: Dual-mode formatting (printf-style and template-based)
**Rationale:**
- Printf for simple, backward-compatible formatting
- Templates for rich data extraction and advanced formatting
- Graceful fallback from template to printf
- Supports ANSI colors and structured output

**Implementation:**
- `OutputFormatter`: Printf-style formatting
- `TemplateFormatter`: Advanced template processing
- Regex extraction for rich data formatting
- Support for both Go text/template and placeholder syntax

**Cross-References**: [REQ:OUTPUT_FORMATTING], [REQ:USABILITY]

## 10. Git Integration [ARCH:GIT_INTEGRATION] [REQ:GIT_INTEGRATION]

### Decision: Git command-line interface for repository information
**Rationale:**
- Simplicity over Git library dependencies
- Relies on user's Git installation
- Consistent with user's Git configuration
- Lightweight implementation

**Implementation:**
- Uses `exec.Command("git", ...)` for Git operations
- Repository detection via `git rev-parse --is-inside-work-tree`
- Branch and hash extraction via `git rev-parse`
- Status detection via `git status --porcelain`

**Cross-References**: [REQ:GIT_INTEGRATION]

## 11. Testing Strategy [ARCH:TESTING_STRATEGY]

### Decision: Comprehensive testing with unit, integration, and infrastructure tests
**Rationale:**
- Comprehensive test coverage (>90% target)
- Fast unit tests for development workflow
- Integration tests for end-to-end scenarios
- Testing infrastructure for complex scenarios
- Aligns with validation criteria defined in requirements

**Structure:**
- Unit tests for individual functions
- Integration tests for workflows
- Test utilities for complex scenarios (corruption, disk space, permissions)
- Performance benchmarks for critical functions

**Note**: This testing strategy implements the validation criteria specified in `requirements.md`. Each requirement's validation criteria informs what types of tests are needed.

**Cross-References**: [REQ:RELIABILITY], [REQ:MAINTAINABILITY]

## 12. Build and Distribution [ARCH:BUILD_DISTRIBUTION]

### Decision: Single binary distribution with cross-platform compilation
**Rationale:**
- Easy deployment (no runtime dependencies)
- Cross-platform support (macOS, Linux)
- Fast execution
- Simple installation

**Build Targets:**
- macOS (ARM64 and AMD64)
- Ubuntu (20.04, 22.04, 24.04)
- Single binary per platform

**Implementation:**
- Makefile with build targets
- Version information injected via `-ldflags`
- Quality gates: linting and testing before build

**Cross-References**: [REQ:USABILITY]

## 13. Code Organization Principles [ARCH:CODE_ORGANIZATION]

### Decision: Clear responsibilities with interfaces for testability
**Rationale:**
- Testable components
- Clear responsibilities
- Easy to extend
- Maintainable codebase

**Principles:**
- Each module has a single, clear responsibility
- Functions are small and focused
- Interfaces where appropriate for testability
- Avoid global state where possible
- Resource management centralized

**Cross-References**: [REQ:MAINTAINABILITY]

## 14. Package Extraction Architecture [ARCH:PACKAGE_EXTRACTION] [REQ:MAINTAINABILITY]

### Decision: Layered extraction approach with zero-breaking-change principle
**Rationale:**
- Enables code reuse across projects
- Improves maintainability through modularization
- Supports gradual migration without breaking existing functionality
- Creates reusable components for CLI application development
- Follows interface-driven design for flexibility

**Extraction Strategy:**
- **Layered Approach**: Extract in dependency order
  1. Foundation layer: Infrastructure components (config, errors)
  2. Utility layer: Domain utilities (git, fileops)
  3. Framework layer: Application framework (cli)
  4. Pattern layer: Processing patterns (processing, formatter)
- **Zero-Breaking-Change**: Maintain backward compatibility through wrapper functions
- **Interface-Driven**: Define interfaces before extraction to prevent coupling
- **Circular Dependency Prevention**: Extract components with no internal dependencies first

**Extracted Packages:**
- `pkg/config`: Configuration management (1,322 lines)
- `pkg/errors`: Error handling and classification (918 lines)
- `pkg/resources`: Resource management and cleanup (486 lines)
- `pkg/formatter`: Output formatting and templates (1,056 lines)
- `pkg/git`: Git integration utilities (255 lines)
- `pkg/cli`: CLI command framework (722 lines)
- `pkg/fileops`: File operations and utilities (1,173 lines)
- `pkg/processing`: Data processing patterns (1,705 lines)

**Total Extracted**: 7,637 lines across 33 files

**Backward Compatibility:**
- Legacy type aliases maintain compatibility
- Wrapper functions preserve existing API
- Original application continues to work unchanged

**Alternatives Considered:**
- Big-bang extraction: Rejected due to high risk of breaking changes
- No extraction: Rejected due to code duplication and maintainability concerns

**Cross-References**: [REQ:MAINTAINABILITY], [REQ:CODE_QUALITY]

## 15. CLI Framework Architecture [ARCH:CLI_FRAMEWORK] [REQ:USABILITY]

### Decision: Extracted CLI framework with builder patterns and Cobra integration
**Rationale:**
- Provides reusable CLI building blocks
- Standardizes command construction patterns
- Supports rapid CLI application development
- Enables consistent user experience across applications
- Integrates with popular Cobra framework

**Framework Components:**
- **Command Builder**: Builder pattern for constructing Cobra commands
- **Flag Manager**: Standardized flag registration and binding
- **Context Manager**: Context propagation for cancellation and timeouts
- **Dry-Run Manager**: Built-in dry-run operation support
- **Version Manager**: Version display and build information handling

**Design Patterns:**
- Builder pattern for command construction
- Interface-based design for testability
- Manager pattern for cross-cutting concerns
- Template pattern for command structure

**Integration:**
- Built on Cobra CLI framework
- Supports standard Cobra features (flags, subcommands, help)
- Extends Cobra with additional patterns and utilities

**Alternatives Considered:**
- Custom CLI framework: Rejected due to maintenance overhead
- Direct Cobra usage: Rejected due to code duplication across projects

**Cross-References**: [REQ:USABILITY], [REQ:MAINTAINABILITY]

## 16. File Operations Architecture [ARCH:FILE_OPERATIONS] [REQ:RELIABILITY]

### Decision: Extracted file operations package with atomic operations and security focus
**Rationale:**
- Provides safe, reliable file operations
- Prevents data corruption through atomic operations
- Enhances security with path validation
- Supports complex directory traversal scenarios
- Enables file comparison and snapshot capabilities

**Core Features:**
- **Atomic Operations**: Safe file writing with rollback capabilities
- **Path Validation**: Security and existence checking with comprehensive validation
- **Directory Traversal**: Configurable walking with exclusion patterns
- **File Comparison**: Hash-based content verification and snapshot comparison
- **Pattern Exclusion**: Doublestar glob pattern matching for file filtering

**Security Focus:**
- Path traversal attack prevention
- Permission validation
- Secure path checking
- Input sanitization

**Performance Characteristics:**
- Atomic operations: ~2x slower than direct writes (acceptable trade-off for safety)
- Path validation: ~50μs per path
- Directory traversal: ~1000 files/ms on modern SSDs
- Hash calculation: ~100MB/s for SHA-256

**Alternatives Considered:**
- Direct standard library usage: Rejected due to lack of atomic operations and security checks
- Third-party file libraries: Rejected due to dependency concerns

**Cross-References**: [REQ:RELIABILITY], [REQ:RESOURCE_MANAGEMENT]

## 17. Processing Patterns Architecture [ARCH:PROCESSING_PATTERNS] [REQ:PERFORMANCE]

### Decision: Extracted data processing patterns with pipeline and concurrent processing support
**Rationale:**
- Provides reusable processing patterns from archive/backup operations
- Enables pipeline-based data processing
- Supports concurrent processing for performance
- Standardizes naming conventions
- Integrates verification capabilities

**Core Patterns:**
- **Pipeline Pattern**: Stage-based processing with composition
- **Concurrent Processing**: Worker pool pattern for parallel execution
- **Naming Conventions**: Timestamp-based naming with Git integration
- **Verification Integration**: Checksum verification with multiple algorithms

**Pipeline Stages:**
- Collection: Gather files/directories
- Processing: Transform or process data
- Verification: Validate integrity
- Output: Generate results

**Concurrent Processing:**
- Worker pool management
- Context propagation for cancellation
- Atomic counters for progress tracking
- Resource cleanup on completion

**Naming Conventions:**
- ISO 8601 timestamp formatting (2006-01-02T150405)
- Git branch/hash integration with dirty status
- Metadata and note support
- Archive, backup, and incremental patterns

**Verification Support:**
- SHA-256, SHA-512, MD5 algorithms
- Checksum calculation and verification
- Serialization support
- Manager pattern for algorithm selection

**Alternatives Considered:**
- Inline processing logic: Rejected due to code duplication
- Third-party processing libraries: Rejected due to dependency and customization concerns

**Cross-References**: [REQ:PERFORMANCE], [REQ:MAINTAINABILITY]

## 18. Auto-Detection Architecture [ARCH:AUTO_DETECTION] [REQ:USABILITY]

### Decision: Automatic command detection based on path type (CLI-015)
**Rationale:**
- Improves user experience by reducing command verbosity
- Enables intuitive CLI usage (just provide path)
- Maintains backward compatibility with explicit commands
- Supports both file and directory operations seamlessly

**Detection Logic:**
- Path type detection via `os.Stat()` and file mode checking
- File detection: Routes to file backup operation
- Directory detection: Routes to directory archive operation
- Unsupported types: Clear error messages with guidance

**Command Routing:**
- Known commands bypass auto-detection (create, config, list, verify, backup, etc.)
- Global flags handled normally (--config, --dry-run, --list)
- Positional arguments trigger auto-detection when not a known command
- Error handling provides clear feedback for unsupported types

**Backward Compatibility:**
- All existing explicit commands continue to work
- Auto-detection only activates for unrecognized positional arguments
- No breaking changes to existing CLI interface

**Implementation:**
- `isFile()`: Detects regular files
- `isDirectory()`: Detects directories
- `validatePath()`: Validates path existence and accessibility
- `handleAutoDetectedCommand()`: Routes to appropriate handler
- `executeWithAutoDetection()`: Integrates with Cobra command execution

**Alternatives Considered:**
- Separate auto-detect command: Rejected due to reduced usability
- No auto-detection: Rejected due to user experience concerns

**Cross-References**: [REQ:USABILITY]

## 19. File Statistics Architecture [ARCH:FILE_STATISTICS] [REQ:OUTPUT_FORMATTING]

### Decision: File statistics gathering for enhanced output formatting
**Rationale:**
- Provides rich file information for output display
- Enables human-readable size formatting
- Supports file type detection
- Enhances user experience with detailed file information

**Statistics Collected:**
- File path and name
- File size (bytes and human-readable)
- Modification time (time.Time and Unix timestamp)
- File permissions/mode
- File type (regular, directory, symlink, device, pipe, socket, other)

**Human-Readable Formatting:**
- Size conversion: bytes → KB → MB → GB → TB
- Format: "1.2MB", "455KB", "1.0GB"
- Precision: One decimal place for readability

**File Type Detection:**
- Based on `os.FileMode` flags
- Supports all standard Unix file types
- Returns descriptive string identifiers

**Integration:**
- Used by output formatter for enhanced display
- Supports template-based formatting with named placeholders
- Integrates with archive and backup output

**Alternatives Considered:**
- No statistics: Rejected due to limited output information
- External library: Rejected due to dependency concerns

**Cross-References**: [REQ:OUTPUT_FORMATTING], [REQ:USABILITY]

## 20. Directory Comparison Architecture [ARCH:DIRECTORY_COMPARISON] [REQ:ARCHIVE_VERIFICATION]

### Decision: Snapshot-based directory comparison system
**Rationale:**
- Enables efficient directory comparison without full traversal
- Supports archive-to-directory comparison
- Provides hash-based content verification
- Enables incremental archive detection

**Snapshot Structure:**
- **DirectorySnapshot**: Map of file paths to FileInfo
- **FileInfo**: Path, size, modification time, hash, directory flag
- Hash-based content identification for accurate comparison

**Comparison Operations:**
- **CreateDirectorySnapshot**: Build snapshot from directory with exclusions
- **CreateArchiveSnapshot**: Build snapshot from ZIP archive
- **CompareSnapshots**: Compare two snapshots for differences
- **IsDirectoryIdenticalToArchive**: Direct directory-to-archive comparison

**Hash-Based Comparison:**
- Uses file content hashes for accurate comparison
- Handles file renames and moves correctly
- Detects content changes even with same modification time
- Supports exclusion patterns for filtering

**Use Cases:**
- Incremental archive detection (compare directory to previous archive)
- Archive verification (compare archive to source directory)
- Directory synchronization (identify differences)

**Alternatives Considered:**
- Timestamp-only comparison: Rejected due to accuracy concerns
- Full content comparison: Rejected due to performance impact

**Cross-References**: [REQ:ARCHIVE_VERIFICATION], [REQ:PERFORMANCE]

## 21. Exclusion Patterns Architecture [ARCH:EXCLUSION_PATTERNS] [REQ:CONFIGURATION]

### Decision: Pattern-based file exclusion system with glob support
**Rationale:**
- Enables flexible file filtering for archives and backups
- Supports common exclusion patterns (version control, build artifacts)
- Provides consistent exclusion behavior across operations
- Integrates with configuration system

**Pattern Matching:**
- **Doublestar Glob Patterns**: Supports `**` for recursive matching
- **Pattern Matcher**: Compiled pattern matcher for efficient matching
- **Multiple Patterns**: Supports arrays of exclusion patterns
- **Path-Based Matching**: Matches against full file paths

**Common Patterns:**
- Version control: `.git/`, `.svn/`, `.hg/`
- Build artifacts: `build/`, `dist/`, `target/`, `*.o`
- Dependencies: `node_modules/`, `vendor/`
- Temporary files: `*.tmp`, `*.temp`, `*.swp`
- OS files: `.DS_Store`, `Thumbs.db`

**Integration:**
- Used by directory traversal for filtering
- Configurable via configuration files
- Supports per-operation exclusion patterns
- Integrates with archive and backup operations

**Performance:**
- Pattern compilation for efficiency
- ~10μs per file for complex glob patterns
- Supports large pattern sets

**Alternatives Considered:**
- Regex patterns: Rejected due to complexity and performance
- Simple string matching: Rejected due to limited functionality

**Cross-References**: [REQ:CONFIGURATION], [REQ:PERFORMANCE]

## 22. Verification Architecture [ARCH:VERIFICATION] [REQ:ARCHIVE_VERIFICATION]

### Decision: Archive verification system with checksum support and integrity checking
**Rationale:**
- Ensures archive integrity and data reliability
- Detects corruption early
- Supports multiple checksum algorithms
- Provides verification status tracking

**Verification Components:**
- **Archive Integrity Checking**: Validates ZIP structure and file readability
- **Checksum Generation**: Calculates checksums for files (SHA-256, SHA-512, MD5)
- **Checksum Verification**: Validates stored checksums against file content
- **Verification Status**: Tracks verification results with timestamps and errors

**Verification Process:**
- Open archive and validate ZIP structure
- Verify each file can be read without errors
- Calculate checksums for files (if enabled)
- Compare calculated checksums with stored checksums (if available)
- Report verification status with detailed error information

**Checksum Algorithms:**
- SHA-256: Default, secure, widely supported
- SHA-512: Higher security, larger checksums
- MD5: Legacy support, faster but less secure

**Verification Status:**
- **VerificationStatus**: Tracks verification timestamp, success status, checksum presence, errors
- JSON serialization for persistence
- Error collection for detailed reporting

**Integration:**
- Optional verification after archive creation
- Standalone verification command
- Verification status in archive listings
- Checksum storage in archive metadata

**Alternatives Considered:**
- No verification: Rejected due to reliability concerns
- External verification tools: Rejected due to dependency and integration concerns

**Cross-References**: [REQ:ARCHIVE_VERIFICATION], [REQ:RELIABILITY]

## 23. Configuration Reflection Architecture [ARCH:CFG_006] [REQ:CFG_006]

### Decision: Reflection-based automatic field discovery with hierarchical source attribution
**Rationale:**
- Zero maintenance: New configuration fields automatically appear without manual updates
- Complete visibility: Shows inheritance chains, merge strategies, and source attribution
- Debugging support: Clear source attribution enables troubleshooting configuration issues
- Performance optimized: Fast inspection suitable for development workflow
- Backward compatible: Existing configuration command functionality preserved and enhanced

**Architecture Components:**

1. **Automatic Field Discovery System**
   - Go reflection-based enumeration of all Config struct fields
   - Recursive traversal of nested structures (e.g., VerificationConfig)
   - Type-aware processing for strings, bools, ints, slices, pointers, structs
   - Field categorization (archive, backup, formatting, verification, etc.)
   - Metadata extraction (field names, types, YAML tags, categories)

2. **Source Tracking Architecture**
   - Integration with CFG-005 layered inheritance system
   - Complete source chain visualization (environment → inheritance → defaults)
   - Merge strategy attribution (override, append, prepend, replace, default)
   - Override point identification in inheritance chain
   - Environment variable detection and attribution

3. **Display Formatting Architecture**
   - Multiple output formats: table (quick scanning), tree (hierarchical), JSON (programmatic)
   - Category-based field organization
   - Type-aware value formatting
   - Zero value detection and display
   - Hierarchical tree visualization with inheritance paths

4. **Performance Optimization Architecture**
   - Reflection result caching with thread-safe access (sync.RWMutex)
   - Schema hash-based cache invalidation
   - Lazy source evaluation (only resolve sources for displayed fields)
   - Incremental resolution support (single field access without full enumeration)
   - Pattern-based field filtering for efficient queries

5. **CLI Interface Architecture**
   - Enhanced config command with comprehensive filtering options
   - --all flag: Show all configuration fields (default behavior)
   - --overrides-only flag: Display only non-default values
   - --sources flag: Show detailed source attribution
   - --format flag: Choose output format (table, tree, JSON)
   - --filter flag: Filter fields by pattern or category

**Integration Points:**
- **CFG-005 Integration**: Leverages layered inheritance system for source tracking
- **Configuration System**: Builds on existing Config struct and loading mechanisms
- **CLI Framework**: Integrates with existing Cobra command structure
- **Testing Infrastructure**: Comprehensive test coverage with performance benchmarks

**Performance Characteristics:**
- Field discovery: <50ms for cache miss, <10ms for cache hit
- Config command: <100ms end-to-end response time
- Single field access: <10ms with incremental resolution
- Memory overhead: <5% increase for caching infrastructure
- Reflection caching: 60%+ performance improvement over uncached reflection

**Alternatives Considered:**
- **Manual maintenance**: Rejected due to maintenance overhead and risk of missing new fields
- **Code generation**: Rejected due to complexity, build-time dependencies, and reduced flexibility
- **External configuration registry**: Rejected due to dependency concerns and synchronization issues

**Cross-References**: [REQ:CFG_006], [ARCH:CONFIG_SYSTEM], [ARCH:CFG_005], [REQ:PERFORMANCE]

## 24. Security Architecture [ARCH:SECURITY] [REQ:RELIABILITY]

### Decision: Comprehensive security measures for file system and archive operations
**Rationale:**
- Prevents security vulnerabilities
- Protects against malicious input
- Ensures data integrity
- Maintains system reliability

**File System Security:**
- **Path Validation**: Prevent directory traversal with enhanced validation
- **Permission Checking**: Verify access rights with structured errors
- **Symlink Handling**: Safe symbolic link processing
- **Temporary Files**: Secure temporary file creation with automatic cleanup

**Archive Security:**
- **Checksum Verification**: SHA-256 integrity checks with context support
- **Compression Bombs**: ZIP bomb protection during creation
- **Path Sanitization**: Safe archive extraction with validation
- **Metadata Validation**: Archive metadata verification with structured errors

**Cross-References**: [REQ:RELIABILITY], [REQ:RESOURCE_MANAGEMENT]

## 25. Extensibility Architecture [ARCH:EXTENSIBILITY] [REQ:MAINTAINABILITY]

### Decision: Plugin-based extensibility with hook system
**Rationale:**
- Enables future feature extensions
- Supports custom formatters and processors
- Maintains backward compatibility
- Provides flexible integration points

**Plugin Interface:**
```go
type ArchivePlugin interface {
    Name() string
    Process(ctx context.Context, archive *Archive) error
    Cleanup() error
}

type BackupPlugin interface {
    Name() string
    Process(ctx context.Context, backup *BackupInfo) error
    Cleanup() error
}
```

**Hook System:**
- **Pre-Archive**: Before archive creation with context
- **Post-Archive**: After archive creation with verification
- **Pre-Verification**: Before verification with context
- **Post-Verification**: After verification with status
- **Pre-Backup**: Before file backup creation
- **Post-Backup**: After file backup creation

**Custom Formatters:**
```go
type OutputFormatter interface {
    Format(template string, data interface{}) (string, error)
    RegisterFunction(name string, fn interface{})
}

type TemplateFormatter interface {
    FormatWithTemplate(input, pattern, tmplStr string) (string, error)
    FormatWithPlaceholders(format string, data map[string]string) string
}
```

**Alternatives Considered:**
- No extensibility: Rejected due to limited flexibility
- Full plugin system: Deferred due to complexity

**Cross-References**: [REQ:MAINTAINABILITY]

## 26. Deployment Architecture [ARCH:DEPLOYMENT] [REQ:USABILITY]

### Decision: Single binary distribution with embedded configuration
**Rationale:**
- Easy deployment (no runtime dependencies)
- Cross-platform support
- Fast execution
- Simple installation

**Binary Distribution:**
- **Single Binary**: Self-contained executable with embedded defaults
- **Cross-Platform**: Linux, macOS, Windows support
- **Static Linking**: No external dependencies
- **Version Embedding**: Build-time version info with Git integration

**Configuration Management:**
- **Default Configuration**: Embedded defaults with immutable specifications
- **Configuration Validation**: Schema-based validation with error reporting
- **Migration Support**: Configuration version upgrades
- **Environment Integration**: Container-friendly configuration with BKPDIR_CONFIG

**Cross-References**: [REQ:USABILITY], [REQ:MAINTAINABILITY]

## 27. Performance Architecture [ARCH:PERFORMANCE] [REQ:PERFORMANCE]

### Decision: Optimized performance with streaming, caching, and concurrent processing
**Rationale:**
- Ensures responsive user experience
- Handles large directories efficiently
- Minimizes resource usage
- Supports scalable operations

**Memory Management:**
- **Streaming Processing**: Large file handling with context support
- **Buffer Pooling**: Memory reuse with efficient allocation
- **Garbage Collection**: Minimal allocation with resource cleanup
- **Resource Cleanup**: Automatic resource management with panic recovery

**I/O Optimization:**
- **Buffered I/O**: Efficient file operations with context cancellation
- **Parallel Processing**: Concurrent file handling with thread safety
- **Compression Tuning**: Optimal compression settings with performance monitoring
- **Disk Space Monitoring**: Available space checking with enhanced error detection

**Scalability:**
- **Large Directories**: Efficient directory traversal with context support
- **Many Files**: Scalable file processing with resource management
- **Deep Hierarchies**: Stack-safe recursion with cancellation checks
- **Long-Running Operations**: Progress reporting with context cancellation

**Cross-References**: [REQ:PERFORMANCE], [REQ:RELIABILITY]

## 28. CLI Commands Architecture [ARCH:CLI_COMMANDS] [REQ:USABILITY]

### Decision: Intuitive command structure with backward compatibility
**Rationale:**
- Easy to learn and use
- Maintains backward compatibility
- Supports both directory and file operations
- Clear command hierarchy

**Command Structure:**
```
bkpdir
├── create [NOTE]                    # Create full directory archive
├── create --incremental [NOTE]     # Create incremental directory archive
├── backup FILE_PATH [NOTE]         # Create file backup
├── list                            # List directory archives
├── --list FILE_PATH                # List file backups for specific file
├── verify [ARCHIVE_NAME]           # Verify archive integrity
├── config                          # Display configuration
├── full [NOTE]                     # Backward compatibility: create full archive
├── inc [NOTE]                      # Backward compatibility: create incremental archive
└── --config                        # Backward compatibility: display configuration
```

**Command Implementation:**
- **Directory Archive Commands**: `create`, `create --incremental`, `full`, `inc`
- **File Backup Commands**: `backup FILE_PATH [NOTE]`, `--list FILE_PATH`
- **Management Commands**: `list`, `verify [ARCHIVE_NAME]`, `config`

**Command Flags:**
- **Global Flags**: `--dry-run`, `--verify`, `--checksum`, `--note`, `--config`, `--list`
- **Context Support**: All long-running operations support context cancellation
- **Output Formatting**: Configurable output with printf-style and template-based formatting

**Cross-References**: [REQ:USABILITY], [REQ:OUTPUT_FORMATTING]

## 29. CI/CD Pipeline Architecture [ARCH:CICD_PIPELINE] [REQ:CODE_QUALITY]

### Decision: AI-optimized CI/CD pipeline with automated quality gates
**Rationale:**
- Supports AI-first development workflow
- Automated quality validation
- Non-blocking execution for AI assistants
- Comprehensive validation reporting

**Pipeline Architecture:**
- **Phase 1: AI Protocol Validation**: Token compliance, documentation sync, protocol adherence
- **Phase 2: Automated Quality Gates**: Lint checking, test execution, coverage validation, build verification
- **Phase 3: AI-Optimized Reporting**: Standardized icon output, token-based traceability, AI-readable status messages

**Technical Components:**
- **Pipeline Orchestration**: Low-priority task scheduling, background processing
- **AI Protocol Validation Engine**: Token validation, documentation sync, protocol compliance
- **Quality Gate Automation**: Zero-human gates, configurable thresholds, failure recovery
- **AI-Optimized Reporting**: Structured output, icon integration, token attribution

**Alternatives Considered:**
- Traditional CI/CD: Rejected due to lack of AI-first optimization
- Manual validation: Rejected due to overhead

**Cross-References**: [REQ:CODE_QUALITY], [REQ:MAINTAINABILITY]

## 30. AI-First Documentation Architecture [ARCH:AI_DOCUMENTATION] [REQ:MAINTAINABILITY]

### Decision: AI-optimized documentation with semantic tokens and validation framework
**Rationale:**
- Enhances AI assistant comprehension
- Enables automated validation
- Maintains consistency across documentation
- Supports AI-first development workflow

**AI Documentation Framework:**
- **AI Validation Framework**: Integration with validation systems, pre-submission hooks, real-time validation
- **AI-Optimized Error Processing**: AI-readable error messages, remediation guidance, structured error responses
- **Compliance Monitoring**: AI behavior tracking, adherence monitoring, bypass audit trail
- **AI-First Documentation Strategy**: Implementation token standardization, icon usage guidelines, cross-reference integrity

**Documentation Standards:**
- **Implementation Token Manager**: Standardized token format with priority and action icons
- **Cross-Reference Validator**: Automated link validation and bidirectional reference checking
- **Content Structure Optimizer**: AI-friendly formatting and hierarchy management
- **Template Management System**: Comprehensive template library for all documentation types

**Alternatives Considered:**
- Traditional documentation: Rejected due to lack of AI optimization
- Manual documentation: Rejected due to maintenance overhead

**Cross-References**: [REQ:MAINTAINABILITY]

## 31. Token System Architecture [ARCH:TOKEN_SYSTEM] [REQ:DOC_016]

### Decision: Comprehensive semantic token system for feature traceability across all layers
**Rationale:**
- Ensures complete traceability from requirements through architecture to implementation
- Enables AI assistant navigation and comprehension
- Supports automated validation of token consistency
- Maintains long-term code maintainability

**Token System Components:**
- **Token Registry**: Central registry in `semantic-tokens.md` with cross-references to all documentation layers
- **Cross-Layer Traceability**: Requirements → Architecture → Implementation → Tests → Code
- **Automated Validation**: Scripts validate token format, consistency, and traceability
- **AI Assistant Integration**: Token-based navigation and feature discovery

**Token Types:**
- `[REQ:*]` - Requirements tokens in `requirements.md`
- `[ARCH:*]` - Architecture decision tokens in `architecture-decisions.md`
- `[IMPL:*]` - Implementation tokens in code comments
- `[TEST:*]` - Test specification tokens (referenced via `[REQ:*]` in test names)

**Token Format Standards:**
- Format: `[TYPE:IDENTIFIER]` with UPPER_SNAKE_CASE identifiers
- Cross-references: Implementation tokens must reference `[ARCH:*]` and `[REQ:*]` tokens
- Consistency: Same token must be used consistently across all layers

**Validation Requirements:**
- Token format compliance across all files
- Cross-layer token consistency validation
- Token traceability completeness checking
- Missing token detection

**Alternatives Considered:**
- Feature IDs without tokens: Rejected due to lack of traceability
- Manual documentation only: Rejected due to maintenance overhead
- External tooling: Rejected due to dependency complexity

**Cross-References**: [REQ:DOC_016], [REQ:MAINTAINABILITY]

## 32. Validation Performance Optimization Architecture [ARCH:PERF_VALIDATION] [REQ:PERFORMANCE]

### Decision: Performance-optimized validation system with caching and fast mode
**Rationale:**
- Addresses performance bottlenecks in validation
- Improves throughput and latency
- Reduces resource utilization
- Enables CI/CD integration

**Performance Optimization Components:**
- **Intelligent Caching System**: SHA-256 content hash, cache TTL (1hr), cache validation, file modification time-based keys
- **Optimized Validation Scripts**: Fast mode execution, batch processing, file size limits, token caching, memory controls
- **Memory Measurement System**: `/usr/bin/time` integration, realistic models, platform detection, actual memory stats

**Performance Results:**
- **Throughput Improvements**: 2,500% improvement for primary validation script
- **Latency Improvements**: 81ms average latency, 53ms fast mode, 4.6-14.7% cache hit improvement
- **Memory Measurement**: 4-25MB realistic memory usage, accurate measurement with `/usr/bin/time`

**Fast Mode Architecture:**
- Skip expensive operations in CI/CD environments
- Fast mode execution for quick validation
- Performance monitoring during execution
- Realistic performance targets

**Cross-References**: [REQ:PERFORMANCE], [REQ:CODE_QUALITY]

## 22. Configuration Testing Architecture [ARCH:TEST_EXCLUDE_MERGE] [REQ:TEST_EXCLUDE_MERGE] [REQ:CONFIGURATION] [REQ:CFG_006]

### Decision: Test-driven validation of exclude patterns merging and source tracking
**Rationale:**
- Validates that configuration merging works correctly for array fields
- Ensures source tracking accurately reflects merged values
- Provides regression protection for merge behavior changes
- Enables debugging of configuration issues through test scenarios

**Test Architecture:**
- **Test Scenario Setup**: Create temporary directory structure with multiple config files
- **Merge Validation**: Verify patterns from all sources are accumulated, not replaced
- **Source Tracking Validation**: Verify config command shows correct source attribution
- **Edge Case Coverage**: Empty arrays, duplicates, order preservation

**Test Structure:**
- Unit test for merge logic validation
- Integration test for config command output validation
- Test fixtures with realistic config file scenarios

**Integration Points:**
- Uses existing config loading infrastructure
- Leverages CFG-006 source tracking system
- Validates CFG-005 merge strategy behavior

**Cross-References**: [REQ:TEST_EXCLUDE_MERGE], [REQ:CONFIGURATION], [REQ:CFG_006], [ARCH:CFG_006]

## 23. Array Field Default Merge Strategy Implementation [ARCH:EXCLUDE_MERGE_FIX] [REQ:CFG_005] [REQ:CONFIGURATION]

### Decision: Array fields default to "merge" strategy to satisfy CFG-005 requirement
**Rationale:**
- Implements CFG-005 requirement that array fields default to merge (accumulate) strategy in all contexts
- Ensures default patterns (`.git/`, `vendor/`) are preserved when users add local patterns
- Fixes implementation bug where array fields were using "override" instead of "merge" by default
- Provides graceful handling of unknown config fields to prevent merge failures
- Handles YAML type conversions robustly

**Architecture Changes:**
- **Merge Strategy Detection**: `extractStrategy` function detects array field keys (like `exclude_patterns`) and defaults to "merge" strategy instead of "override" in all contexts (inheritance chains and sequential file processing)
- **State Management**: `applyMergeStrategies` uses current state (`resultMap`) instead of original state (`dstMap`) for merge operations to ensure proper accumulation
- **Pattern Preservation**: Default patterns explicitly copied into result before merge operations
- **Deduplication**: Merge operations deduplicate values to prevent duplicates
- **Unknown Field Handling**: Unknown config fields (like `inherit`, `verification`) are gracefully skipped instead of aborting merge operations
- **YAML Type Conversion**: Robust handling of `[]interface{}` types from YAML unmarshaling, converting to `[]string` as needed
- **Metadata Field Filtering**: Metadata fields used for inheritance processing (like `inherit`) are filtered out before merge operations

**Integration Points:**
- Implements CFG-005 merge strategy system requirement
- Integrates with CFG-006 source tracking
- Preserves CFG-002 merge behavior specification
- Ensures backward compatibility with existing config files

**Cross-References**: [REQ:CFG_005], [REQ:CONFIGURATION], [ARCH:CFG_005], [ARCH:CFG_006]
