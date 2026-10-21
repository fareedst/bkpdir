# [IMPL-INTERFACE_FIRST] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Define package interfaces before concrete extraction so consumers depend on contracts, implementations stay swappable, and tests use mocks.

INPUT: consumer call sites, target packages for extraction
OUTPUT: minimal interfaces, concrete implementations, versioned evolution path
DATA: interface method sets, API documentation contracts, mock types

## DESIGN_PACKAGE_INTERFACES

SPEC-ID: IMPL-INTERFACE_FIRST::DESIGN_PACKAGE_INTERFACES
STEP T001: derive minimal interfaces from consumer call sites and document contracts before moving implementation code

- [IMPL-INTERFACE_FIRST] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: derive minimal interfaces from consumer call sites and document contracts before moving implementation code.

PROCEDURE design_package_interfaces():
  FOR each module targeted for extraction:
    ANALYZE consumer call sites to determine required methods
    DEFINE minimal interface with only methods consumers actually call
    PLACE interface definition in the consuming module
    DOCUMENT interface contract with API comments
  END FOR

## IMPLEMENT_AGAINST_INTERFACES

SPEC-ID: IMPL-INTERFACE_FIRST::IMPLEMENT_AGAINST_INTERFACES
STEP T001: implement concrete types in provider packages and wire them at the composition root only

- [IMPL-INTERFACE_FIRST] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: implement concrete types in provider packages and wire them at the composition root only.

PROCEDURE implement_against_interfaces():
  FOR each defined interface:
    CREATE concrete implementation in provider module
    VERIFY implementation satisfies interface contract through validation checks
    WIRE concrete type to interface at composition root
  END FOR

## ENABLE_INDEPENDENT_TESTING

SPEC-ID: IMPL-INTERFACE_FIRST::ENABLE_INDEPENDENT_TESTING
STEP T001: test each extracted package with mocks so unit tests avoid real downstream dependencies

- [IMPL-INTERFACE_FIRST] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: test each extracted package with mocks so unit tests avoid real downstream dependencies.

PROCEDURE enable_independent_testing():
  FOR each extracted module:
    GENERATE mock implementations of consumed interfaces
    WRITE unit tests using mocks without real external IO dependencies
    VALIDATE module builds and tests pass without other modules
  END FOR

## VERSION_INTERFACE_EVOLUTION

SPEC-ID: IMPL-INTERFACE_FIRST::VERSION_INTERFACE_EVOLUTION
STEP T001: evolve APIs via additive interfaces or v2 types with adapters instead of breaking existing consumers

- [IMPL-INTERFACE_FIRST] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: evolve APIs via additive interfaces or v2 types with adapters instead of breaking existing consumers.

PROCEDURE version_interface_evolution():
  WHEN interface change is needed:
    IF change is additive THEN
      EXTEND via embedding or new interface type
    ELSE IF change is breaking THEN
      CREATE v2 interface with new type name
      MAINTAIN adapter between old and new during transition
    END IF
  END WHEN
