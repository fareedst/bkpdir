# [IMPL-LARGE_FILE_CHALLENGE] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY]

## Summary contract

Decomposes oversized formatter.go into focused packages (templates, printf, output collector, ANSI) while preserving public interfaces used by the CLI.

INPUT: monolithic formatter.go responsibilities
OUTPUT: extracted packages with unchanged adapter surfaces
DATA: interface contracts, test coverage for extracted units

## LARGE_FILE_DECOMPOSITION

SPEC-ID: IMPL-LARGE_FILE_CHALLENGE::LARGE_FILE_DECOMPOSITION
STEP T001: extract template engine, printf formatter, output collector, and ANSI helpers into separate packages and keep Default* adapters delegating to new modules

- [IMPL-LARGE_FILE_CHALLENGE] [ARCH-CODE_ORGANIZATION] [REQ-MAINTAINABILITY] — How: extract template engine, printf formatter, output collector, and ANSI helpers into separate packages and keep Default* adapters delegating to new modules.

PROCEDURE large_file_decomposition():
  EXTRACT template engine package
  EXTRACT printf formatter package
  EXTRACT output collector package
  EXTRACT ANSI support package
  VERIFY interface compatibility via existing tests and adapters
