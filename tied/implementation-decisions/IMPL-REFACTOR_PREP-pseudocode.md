# [IMPL-REFACTOR_PREP] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Prepare the codebase for extraction by defining interfaces first, splitting large files, and eliminating circular dependencies before packages move.

INPUT: components targeted for extraction, dependency graph, large source files
OUTPUT: interface boundaries, decomposed files, acyclic dependency layers
DATA: REFACTOR preparation markers, interface definitions, dependency clusters

## DEFINE_INTERFACES_BEFORE_EXTRACTION

SPEC-ID: IMPL-REFACTOR_PREP::DEFINE_INTERFACES_BEFORE_EXTRACTION
STEP T001: replace direct struct coupling with interfaces at every boundary before moving code into new packages

- [IMPL-REFACTOR_PREP] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: replace direct struct coupling with interfaces at every boundary before moving code into new packages.

PROCEDURE define_interfaces_before_extraction():
  FOR each component targeted for extraction:
    IDENTIFY inbound dependencies (callers of this component)
    IDENTIFY outbound dependencies (callees of this component)
    DEFINE interface for each dependency boundary
    REPLACE direct struct usage with interface usage at call sites
  END FOR
  VALIDATE no concrete type crosses a package boundary

## DECOMPOSE_LARGE_FILES

SPEC-ID: IMPL-REFACTOR_PREP::DECOMPOSE_LARGE_FILES
STEP T001: split oversized files by responsibility while preserving package-level API behavior

- [IMPL-REFACTOR_PREP] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: split oversized files by responsibility while preserving package-level API behavior.

PROCEDURE decompose_large_files():
  FOR each source_file exceeding complexity threshold:
    IDENTIFY distinct responsibilities within the file
    SPLIT into focused files with one responsibility per file
    PRESERVE package-level API with no external behavior change
    UPDATE internal imports within the package
  END FOR

## ANALYZE_AND_CLEANUP_DEPENDENCIES

SPEC-ID: IMPL-REFACTOR_PREP::ANALYZE_AND_CLEANUP_DEPENDENCIES
STEP T001: break mutual dependency clusters with interfaces or lower-layer shared code before extraction

- [IMPL-REFACTOR_PREP] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: break mutual dependency clusters with interfaces or lower-layer shared code before extraction.

PROCEDURE analyze_and_cleanup_dependencies():
  BUILD dependency_graph of current codebase
  IDENTIFY clusters of mutual dependencies
  FOR each circular dependency cluster:
    INTRODUCE interface to break the cycle
    OR RESTRUCTURE shared code into a lower dependency layer
  END FOR
