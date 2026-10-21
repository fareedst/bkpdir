# [IMPL-INTERFACE_DRIVEN] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Define minimal consumer-side interfaces from real call sites, implement concretely in provider packages, and verify boundaries with mocks.

INPUT: provider package, consumer roots (main, pkg, tests)
OUTPUT: interface definitions, wiring at composition root, mock-backed unit tests
DATA: method names, argument shapes, fake implementations

## DISCOVER_CONSUMER_SURFACE

SPEC-ID: IMPL-INTERFACE_DRIVEN::DISCOVER_CONSUMER_SURFACE
STEP T001: map call sites in main, pkg, and tests to the minimal method set each consumer needs

- [IMPL-INTERFACE_DRIVEN] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: map call sites in main, pkg, and tests to the minimal method set each consumer needs.

PROCEDURE discover_consumer_surface(provider_pkg, consumer_roots):
  FOR EACH reference FROM consumer_roots TO provider_pkg:
    RECORD method name and argument shapes
  END FOR
  RETURN deduplicated method list

## DEFINE_INTERFACE_IN_CONSUMER

SPEC-ID: IMPL-INTERFACE_DRIVEN::DEFINE_INTERFACE_IN_CONSUMER
STEP T001: declare consumer-side interfaces with documentation linked to maintainability requirements

- [IMPL-INTERFACE_DRIVEN] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: declare consumer-side interfaces with documentation linked to maintainability requirements.

PROCEDURE define_interface_in_consumer(name, methods):
  DEFINE interface named name with methods
  DOCUMENT each method with API documentation linking [REQ-MAINTAINABILITY]

## IMPLEMENT_AND_REGISTER

SPEC-ID: IMPL-INTERFACE_DRIVEN::IMPLEMENT_AND_REGISTER
STEP T001: wire concrete provider types into composition root only after compile-time interface satisfaction checks

- [IMPL-INTERFACE_DRIVEN] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: wire concrete provider types into composition root only after compile-time interface satisfaction checks.

PROCEDURE implement_and_register(root_wire):
  FOR EACH interface I required by root_wire:
    INSTANTIATE concrete type T from provider package
    ASSERT T satisfies I at compile time
    INJECT T into handlers (CLI, tests)
  END FOR

## VERIFY_WITH_MOCKS

SPEC-ID: IMPL-INTERFACE_DRIVEN::VERIFY_WITH_MOCKS
STEP T001: substitute fakes at each interface boundary so unit tests avoid live filesystem and version-control dependencies

- [IMPL-INTERFACE_DRIVEN] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: substitute fakes at each interface boundary so unit tests avoid live filesystem and version-control dependencies.

PROCEDURE verify_with_mocks():
  FOR EACH module with external dependencies:
    SUBSTITUTE mock implementations in unit test files
    RUN unit test suite without live filesystem or version-control dependencies where possible
  END FOR
