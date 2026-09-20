# [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY]

## Summary contract

pkg/fileops provides secure path validation and directory traversal with exclusions, depth limits, symlink policy, and file listing helpers.

INPUT: root path, TraversalOptions, exclude patterns
OUTPUT: visitor callbacks, file path lists, validation errors
DATA: DefaultTraverser, PathValidator, PatternMatcher

## NEW_TRAVERSER

SPEC-ID: IMPL-FILE_OPERATIONS::NEW_TRAVERSER
STEP T001: RETURN DefaultTraverser WITH no preloaded exclusion patterns

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with no preloaded exclusion patterns.

PROCEDURE NEW_TRAVERSER():
  RETURN DefaultTraverser{}

## NEW_TRAVERSER_WITH_PATTERNS

SPEC-ID: IMPL-FILE_OPERATIONS::NEW_TRAVERSER_WITH_PATTERNS
STEP T001: CONSTRUCT DefaultTraverser WITH PatternMatcher FROM patterns

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: construct DefaultTraverser with PatternMatcher built from initial pattern list.

PROCEDURE NEW_TRAVERSER_WITH_PATTERNS(patterns):
  RETURN DefaultTraverser{matcher: NewPatternMatcher(patterns)}

## WALK

SPEC-ID: IMPL-FILE_OPERATIONS::WALK
STEP T001: DELEGATE to WalkWithOptions WITH default TraversalOptions

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: delegate to WalkWithOptions with default options (no exclusions, unlimited depth).

PROCEDURE WALK(root, visitor):
  OPTIONS = default TraversalOptions
  RETURN WalkWithOptions(root, OPTIONS, visitor)

## WALK_WITH_EXCLUSIONS

SPEC-ID: IMPL-FILE_OPERATIONS::WALK_WITH_EXCLUSIONS
STEP T001: WalkWithOptions WITH ExcludePatterns AND no symlink follow

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: WalkWithOptions with ExcludePatterns, no symlink follow, ignore permission errors.

PROCEDURE WALK_WITH_EXCLUSIONS(root, excludePatterns, visitor):
  OPTIONS = {ExcludePatterns, FollowSymlinks false, IgnorePermErrors true}
  RETURN traverser.WalkWithOptions(root, OPTIONS, visitor)

## WALK_WITH_OPTIONS

SPEC-ID: IMPL-FILE_OPERATIONS::WALK_WITH_OPTIONS
STEP T001: VALIDATE root path AND existence
STEP T002: INVOKE walkRecursive WITH exclusionMatcher AND visitor

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: validate root exists, build PatternMatcher from options or traverser matcher, invoke walkRecursive.

PROCEDURE WALK_WITH_OPTIONS(root, options, visitor):
  VALIDATE root path and existence
  BUILD exclusionMatcher FROM options.ExcludePatterns OR traverser.matcher
  RETURN walkRecursive(root, root, depth 0, options, exclusionMatcher, visitor)

## WALK_RECURSIVE

SPEC-ID: IMPL-FILE_OPERATIONS::WALK_RECURSIVE
STEP T001: ENFORCE depth limit AND exclusion rules
STEP T002: CALL visitor AND recurse into directories

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: enforce depth limit, apply exclusion and hidden-file rules, optionally follow symlinks, call visitor, recurse into directories.

PROCEDURE WALK_RECURSIVE(root, currentPath, depth, options, matcher, visitor):
  IF depth > MaxDepth THEN RETURN
  LSTAT currentPath; compute relPath for exclusion
  IF matcher.ShouldExclude(relPath) THEN skip file or SkipDir
  IF IgnoreHidden AND basename starts with . THEN skip
  IF symlink AND NOT FollowSymlinks THEN visitor only ELSE resolve and continue
  CALL visitor; IF directory THEN ReadDir and recurse children

## LIST_FILES

SPEC-ID: IMPL-FILE_OPERATIONS::LIST_FILES
STEP T001: COLLECT regular file paths via ListFilesWithExclusions

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: collect regular file paths under root with optional recursion via ListFilesWithExclusions.

PROCEDURE LIST_FILES(root, recursive):
  RETURN traverser.ListFiles(root, recursive)

## LIST_FILES_WITH_EXCLUSIONS

SPEC-ID: IMPL-FILE_OPERATIONS::LIST_FILES_WITH_EXCLUSIONS
STEP T001: WALK WITH exclusions AND collect non-directory paths

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: walk tree applying exclusion patterns and append file paths only.

PROCEDURE LIST_FILES_WITH_EXCLUSIONS(root, excludePatterns, recursive):
  WALK WITH exclusions; COLLECT non-directory paths into slice

## VALIDATE_PATH

SPEC-ID: IMPL-FILE_OPERATIONS::VALIDATE_PATH
STEP T001: REJECT empty paths AND run IsSecurePath
STEP T002: CLEAN path WITHOUT failing on normalization drift

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject empty paths, run IsSecurePath, clean path without failing on normalization drift.

PROCEDURE VALIDATE_PATH(path):
  IF path empty THEN error
  IF NOT IsSecurePath(path) THEN error unsafe path
  RETURN nil

## VALIDATE_EXISTENCE

SPEC-ID: IMPL-FILE_OPERATIONS::VALIDATE_EXISTENCE
STEP T001: VALIDATE_PATH THEN Stat path

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then Stat path; return not-exist or access errors.

PROCEDURE VALIDATE_EXISTENCE(path):
  VALIDATE_PATH(path)
  Stat path OR return does-not-exist error

## VALIDATE_READABLE

SPEC-ID: IMPL-FILE_OPERATIONS::VALIDATE_READABLE
STEP T001: VALIDATE_EXISTENCE THEN Open path for read

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidateExistence then Open path for read; close on success.

PROCEDURE VALIDATE_READABLE(path):
  VALIDATE_EXISTENCE(path)
  Open for read OR return not-readable error

## VALIDATE_WRITABLE

SPEC-ID: IMPL-FILE_OPERATIONS::VALIDATE_WRITABLE
STEP T001: VALIDATE_PATH THEN test write via temp file OR OpenFile WRONLY

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: ValidatePath then test write via temp file in directory or OpenFile WRONLY for files; recurse to parent when path missing.

PROCEDURE VALIDATE_WRITABLE(path):
  VALIDATE_PATH(path)
  IF not exists THEN ValidateWritable(parentDir)
  IF directory THEN create and remove temp file ELSE open WRONLY

## IS_SECURE_PATH

SPEC-ID: IMPL-FILE_OPERATIONS::IS_SECURE_PATH
STEP T001: REJECT .. traversal AND suspicious shell-expansion characters

- [IMPL-FILE_OPERATIONS] [ARCH-FILE_OPERATIONS] [REQ-RELIABILITY] — How: reject .. traversal, null bytes, newlines, and shell-expansion characters in path strings.

PROCEDURE IS_SECURE_PATH(path):
  IF CONTAINS ".." OR suspicious chars (~, $, null, CR, LF) THEN RETURN false
  RETURN true
