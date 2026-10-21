# [IMPL-DEPENDENCY_MGMT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Keeps extracted packages within declared layer import boundaries and uses Go modules for versioned package graphs.

INPUT: package source trees, dependency hierarchy manifest
OUTPUT: violation reports and tidy go.mod graphs
DATA: allowed_deps per layer, go.mod pins

## ENFORCE_DEPENDENCY_RULES

SPEC-ID: IMPL-DEPENDENCY_MGMT::ENFORCE_DEPENDENCY_RULES
STEP T001: EXECUTE PROCEDURE enforce_dependency_rules

- [IMPL-DEPENDENCY_MGMT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: for each package compare actual imports against allowed_deps for its layer and report out-of-layer imports with path context.

PROCEDURE enforce_dependency_rules():
  FOR EACH package:
    actual = parse imports from Go sources
    allowed = dependency_hierarchy[package.layer]
    IF import not in allowed THEN REPORT violation

## VERSION_WITH_GO_MODULES

SPEC-ID: IMPL-DEPENDENCY_MGMT::VERSION_WITH_GO_MODULES
STEP T001: Pin module versions in go.mod and run go mod tidy to keep the dependency graph consistent after extractions

- [IMPL-DEPENDENCY_MGMT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: pin module versions in go.mod and run go mod tidy to keep the dependency graph consistent after extractions.

PROCEDURE version_with_go_modules():
  EACH extracted package maintains go.mod (or sub-module path)
  PIN versions in require directives
  RUN go mod tidy to validate graph
