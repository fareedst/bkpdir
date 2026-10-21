# [IMPL-LAYERED_EXTRACTION] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Extract packages in dependency order: infrastructure first, then utilities, framework, and patterns last, with no circular references between layers.

INPUT: packages to extract, layer assignment per package
OUTPUT: four stable layers with acyclic dependencies
DATA: layer order, package dependency graph

## EXTRACT_INFRASTRUCTURE_LAYER

SPEC-ID: IMPL-LAYERED_EXTRACTION::EXTRACT_INFRASTRUCTURE_LAYER
STEP T001: extract config and errors packages first as the foundation for all higher layers

- [IMPL-LAYERED_EXTRACTION] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: extract config and errors packages first as the foundation for all higher layers.

PROCEDURE extract_infrastructure_layer():
  EXTRACT config and errors into infrastructure packages
  VALIDATE no infrastructure package depends on utilities, framework, or patterns
  VALIDATE higher layers may depend on infrastructure only through public APIs

## EXTRACT_UTILITIES_LAYER

SPEC-ID: IMPL-LAYERED_EXTRACTION::EXTRACT_UTILITIES_LAYER
STEP T001: extract formatter and git helpers after infrastructure is stable

- [IMPL-LAYERED_EXTRACTION] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: extract formatter and git helpers after infrastructure is stable.

PROCEDURE extract_utilities_layer():
  EXTRACT formatter and git utilities
  ASSERT utilities depend only on infrastructure and core runtime libraries
  VALIDATE utilities do not depend on framework or patterns packages

## EXTRACT_FRAMEWORK_LAYER

SPEC-ID: IMPL-LAYERED_EXTRACTION::EXTRACT_FRAMEWORK_LAYER
STEP T001: extract CLI framework components once utilities and infrastructure are in place

- [IMPL-LAYERED_EXTRACTION] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: extract CLI framework components once utilities and infrastructure are in place.

PROCEDURE extract_framework_layer():
  EXTRACT CLI framework package
  ASSERT framework depends on infrastructure and utilities, not patterns
  VALIDATE CLI wiring compiles with extracted dependencies

## EXTRACT_PATTERNS_LAYER

SPEC-ID: IMPL-LAYERED_EXTRACTION::EXTRACT_PATTERNS_LAYER
STEP T001: extract processing patterns last so each layer builds on stable lower layers without cycles

- [IMPL-LAYERED_EXTRACTION] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: extract processing patterns last so each layer builds on stable lower layers without cycles.

PROCEDURE extract_patterns_layer():
  EXTRACT processing patterns package
  ASSERT patterns depend on infrastructure, utilities, and framework
  VALIDATE dependency graph has no cycles across all four layers
