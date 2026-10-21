# [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY]

## Summary contract

Phased extraction of root packages into pkg/ with legacy type aliases and wrapper functions preserving CLI compatibility.

INPUT: phase package list, source root files
OUTPUT: pkg/<package>/ modules, root aliases and wrappers
DATA: interface definitions, import path updates, deprecation comments

## EXTRACT_PHASE

SPEC-ID: IMPL-PACKAGE_EXTRACTION::EXTRACT_PHASE
STEP T001: EXECUTE PROCEDURE EXTRACT_PHASE

- [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: for each package in a phase move code to pkg/, define interfaces first, update imports, and require go build plus all tests green.

PROCEDURE EXTRACT_PHASE(phase):
  FOR EACH package IN phase.packages:
    DEFINE interfaces for dependencies
    MOVE code TO pkg/<package>/
    CREATE root type aliases and wrapper functions
    UPDATE internal imports
    VALIDATE go build ./... AND all tests pass

## MAINTAIN_LEGACY_COMPATIBILITY

SPEC-ID: IMPL-PACKAGE_EXTRACTION::MAINTAIN_LEGACY_COMPATIBILITY
STEP T001: Expose type aliases and delegating wrappers at root with deprecation comments pointing to new pkg paths

- [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: expose type aliases and delegating wrappers at root with deprecation comments pointing to new pkg paths.

PROCEDURE MAINTAIN_LEGACY_COMPATIBILITY():
  ADD type aliases: type T = pkg.T
  ADD wrapper functions delegating to pkg
  ADD deprecation comments with new import paths
  VALIDATE CLI entry points behave identically
