# [IMPL-EXTRACTION_PRINCIPLES] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Plan module extraction in strict dependency order so each layer validates and tests independently before the next layer moves.

INPUT: module dependency graph, extraction_layers ordered list
OUTPUT: extraction schedule with validated layer boundaries
DATA: layer deps, interface boundaries per module

## PLAN_EXTRACTION_ORDER

SPEC-ID: IMPL-EXTRACTION_PRINCIPLES::PLAN_EXTRACTION_ORDER
STEP T001: schedule modules only after their declared layer dependencies and interface boundaries validate

- [IMPL-EXTRACTION_PRINCIPLES] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: schedule modules only after their declared layer dependencies and interface boundaries validate.

PROCEDURE plan_extraction_order():
  FOR layer in extraction_layers ascending:
    FOR module in layer.modules:
      VALIDATE module has no dependencies outside declared layer deps
      VALIDATE module exposes interfaces, not concrete types, at boundaries
      SCHEDULE module for extraction after dependencies are extracted
    END FOR
  END FOR
