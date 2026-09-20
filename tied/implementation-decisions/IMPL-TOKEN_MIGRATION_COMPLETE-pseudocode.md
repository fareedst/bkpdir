# [IMPL-TOKEN_MIGRATION_COMPLETE] [ARCH-TOKEN_SYSTEM] [REQ-CODE_QUALITY] [REQ-DOC_016]

## Summary contract

Complete migration of legacy tokens to canonical STDD registry with inventory scripts and traceability validation.

INPUT: legacy token artifacts, project token stubs
OUTPUT: updated semantic-tokens.md tables, archived governance docs, validation reports
DATA: token_inventory.json, render_token_tables output, validate-token-traceability results

## ARCHIVE_MIGRATION_ARTIFACTS

SPEC-ID: IMPL-TOKEN_MIGRATION_COMPLETE::ARCHIVE_MIGRATION_ARTIFACTS
STEP T001: Extract and archive legacy migration documents under docs/governance for historical reference

- [IMPL-TOKEN_MIGRATION_COMPLETE] [ARCH-TOKEN_SYSTEM] [REQ-CODE_QUALITY] [REQ-DOC_016] — How: extract and archive legacy migration documents under docs/governance for historical reference.

PROCEDURE ARCHIVE_MIGRATION_ARTIFACTS():
  MOVE legacy token docs TO docs/governance archive

## GENERATE_TOKEN_INVENTORY

SPEC-ID: IMPL-TOKEN_MIGRATION_COMPLETE::GENERATE_TOKEN_INVENTORY
STEP T001: Run token_inventory.py to produce canonical token JSON used by downstream render scripts

- [IMPL-TOKEN_MIGRATION_COMPLETE] [ARCH-TOKEN_SYSTEM] [REQ-CODE_QUALITY] [REQ-DOC_016] — How: run token_inventory.py to produce canonical token JSON used by downstream render scripts.

PROCEDURE GENERATE_TOKEN_INVENTORY():
  RUN token_inventory.py
  OUTPUT token JSON registry snapshot

## RENDER_TOKEN_TABLES

SPEC-ID: IMPL-TOKEN_MIGRATION_COMPLETE::RENDER_TOKEN_TABLES
STEP T001: Run render_token_tables.py to regenerate semantic-tokens.md tables from inventory JSON

- [IMPL-TOKEN_MIGRATION_COMPLETE] [ARCH-TOKEN_SYSTEM] [REQ-CODE_QUALITY] [REQ-DOC_016] — How: run render_token_tables.py to regenerate semantic-tokens.md tables from inventory JSON.

PROCEDURE RENDER_TOKEN_TABLES():
  RUN render_token_tables.py
  UPDATE semantic-tokens.md tables

## VALIDATE_TRACEABILITY

SPEC-ID: IMPL-TOKEN_MIGRATION_COMPLETE::VALIDATE_TRACEABILITY
STEP T001: EXECUTE PROCEDURE VALIDATE_TRACEABILITY

- [IMPL-TOKEN_MIGRATION_COMPLETE] [ARCH-TOKEN_SYSTEM] [REQ-CODE_QUALITY] [REQ-DOC_016] — How: run validate-token-traceability.sh and token-coverage-analysis.sh expecting full coverage before replacing project-tokens.yaml stub.

PROCEDURE VALIDATE_TRACEABILITY():
  RUN validate-token-traceability.sh
  RUN token-coverage-analysis.sh expecting 100% coverage
  REPLACE project-tokens.yaml WITH pointer stub
