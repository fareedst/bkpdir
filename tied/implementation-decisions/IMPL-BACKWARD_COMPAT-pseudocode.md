# [IMPL-BACKWARD_COMPAT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Preserve the existing backup application API while extracting components into separate packages using aliases, wrappers, and interface boundaries.

INPUT: extracted package symbols, root package consumers
OUTPUT: unchanged compile-time API at root; isolated extracted modules
DATA: type aliases, wrapper functions, boundary interfaces

## PROVIDE_TYPE_ALIASES

SPEC-ID: IMPL-BACKWARD_COMPAT::PROVIDE_TYPE_ALIASES
STEP T001: re-export moved types at original locations so existing imports compile without change

- [IMPL-BACKWARD_COMPAT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: re-export moved types at original locations so existing imports compile without change.

PROCEDURE provide_type_aliases():
  FOR each type T moved to extracted package pkg:
    DECLARE legacy type alias at original location that maps T to pkg.T
    ADD deprecation comment pointing to pkg.T
  END FOR
  VALIDATE all existing consumers resolve symbols without path changes

## PROVIDE_WRAPPER_FUNCTIONS

SPEC-ID: IMPL-BACKWARD_COMPAT::PROVIDE_WRAPPER_FUNCTIONS
STEP T001: delegate moved functions through root wrappers that preserve original signatures

- [IMPL-BACKWARD_COMPAT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: delegate moved functions through root wrappers that preserve original signatures.

PROCEDURE provide_wrapper_functions():
  FOR each function F moved to extracted package pkg:
    CREATE wrapper at original location that calls pkg.F with same parameters and returns
    PRESERVE original function signature exactly
    ADD deprecation comment pointing to pkg.F
  END FOR

## ISOLATE_VIA_INTERFACES

SPEC-ID: IMPL-BACKWARD_COMPAT::ISOLATE_VIA_INTERFACES
STEP T001: ensure extracted packages depend only on declared interfaces, not root application code

- [IMPL-BACKWARD_COMPAT] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: ensure extracted packages depend only on declared interfaces, not root application code.

PROCEDURE isolate_via_interfaces():
  FOR each extracted package:
    ASSERT dependencies are limited to runtime core modules and declared packages
    ASSERT package does not depend on root-level application code
    ASSERT communication with root uses interfaces only, not concrete root types
  END FOR
