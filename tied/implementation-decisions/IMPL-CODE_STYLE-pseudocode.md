# [IMPL-CODE_STYLE] [REQ-CODE_QUALITY]

## Summary contract

Enforce project naming, formatting, and documentation conventions so public APIs stay consistent and lint-clean across the project.

INPUT: source files, public symbols
OUTPUT: pass/fail validation results
DATA: naming rules, formatting rules, documentation rules

## NAMING_RULES

SPEC-ID: IMPL-CODE_STYLE::NAMING_RULES
STEP T001: enforce UpperCamelCase for public symbols and lowerCamelCase for internal symbols

- [IMPL-CODE_STYLE] [REQ-CODE_QUALITY] — How: enforce UpperCamelCase for public symbols and lowerCamelCase for internal symbols.

DATA NAMING_RULES:
  public: UpperCamelCase
  internal: lowerCamelCase
  acronyms: preserve case (e.g. HTTPServer, xmlParser)

PROCEDURE VALIDATE_NAMING(symbol):
  IF symbol.is_public THEN ASSERT name matches UpperCamelCase
  ELSE ASSERT name matches lowerCamelCase
  RETURN pass OR fail

## FORMATTING_RULES

SPEC-ID: IMPL-CODE_STYLE::FORMATTING_RULES
STEP T001: apply project-standard formatting and lint rules on each file
STEP T002: VERIFY all operations handle errors appropriately via lint and return-value checks

- [IMPL-CODE_STYLE] [REQ-CODE_QUALITY] — How: apply project-standard formatting and lint rules on each file.

DATA FORMATTING_RULES:
  formatter: project standard formatter
  linter: project standard linter
  standard: repository formatting profile

PROCEDURE VALIDATE_FORMATTING(file):
  RUN formatter ON file
  RUN linter ON file
  RETURN lint results

## DOCUMENTATION_RULES

SPEC-ID: IMPL-CODE_STYLE::DOCUMENTATION_RULES
STEP T001: require module-level and public-symbol doc comments; inline comments only for non-obvious logic

- [IMPL-CODE_STYLE] [REQ-CODE_QUALITY] — How: require module-level and public-symbol doc comments; inline comments only for non-obvious logic.

DATA DOCUMENTATION_RULES:
  module_level: every module has module-level comment
  public_symbols: every public type, procedure, and method has doc comment
  complex_logic: inline comments for non-obvious logic only
  examples: test suites contain usage examples

PROCEDURE VALIDATE_DOCUMENTATION(module):
  FOR EACH public symbol IN module:
    ASSERT doc comment exists
  ASSERT module-level comment exists
  RETURN pass OR fail
