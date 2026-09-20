# [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-FILE_BACKUP] [REQ-GOV_REGISTRY_COMPLETENESS]

## Summary contract

Governance tooling to extract legacy project tokens, synthesize cross-links to STDD REQ/ARCH/IMPL, author registry tables, and decommission legacy YAML in favor of tied/semantic-tokens.yaml.

INPUT: project-tokens.yaml, source trees, semantic token registry
OUTPUT: normalized JSON artifacts, semantic-tokens.md tables, pointer stubs
DATA: token groups (actions, features, semantic_tokens), migration status

## STRUCTURED_EXTRACTION

SPEC-ID: IMPL-TOKEN_SYSTEM::STRUCTURED_EXTRACTION
STEP T001: parse legacy token YAML groups and emit normalized JSON with descriptions, status, and source paths per category
STEP T002: INTEGRATE automated validation with existing tied_validate_consistency; run automated validation integration testing and navigation accuracy testing for AI assistants

- [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: parse legacy token YAML groups and emit normalized JSON with descriptions, status, and source paths per category.

PROCEDURE structured_extraction():
  PARSE project-tokens.yaml
  FOR EACH group IN actions, features, semantic_tokens:
    EMIT normalized JSON artifact
  ARCHIVE artifacts under docs/governance/

## CROSS_LINK_SYNTHESIS

SPEC-ID: IMPL-TOKEN_SYSTEM::CROSS_LINK_SYNTHESIS
STEP T001: resolve canonical REQ/ARCH/IMPL references for each legacy token and flag gaps when links are missing

- [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-GOV_REGISTRY_COMPLETENESS] [REQ-IMMUTABLE_DIRECTORY_OPERATIONS] — How: resolve canonical REQ/ARCH/IMPL references for each legacy token and flag gaps when links are missing.

PROCEDURE cross_link_synthesis():
  FOR EACH legacy_token:
    MAP identifiers to STDD tokens via prefix heuristics
    IF missing REQ/ARCH/IMPL link THEN OPEN tracking gap

## REGISTRY_AUTHORING (DOC-ONLY NOTE)

SPEC-ID: IMPL-TOKEN_SYSTEM::REGISTRY_AUTHORING
STEP T001: EXECUTE block PROCEDURE body

- [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: generate Markdown registry tables from normalized token artifacts for semantic-tokens.md.

NOTE: Markdown registry table generation is process documentation; no dedicated Go implementation lead required.

## YAML_DECOMMISSIONING (DOC-ONLY NOTE)

SPEC-ID: IMPL-TOKEN_SYSTEM::YAML_DECOMMISSIONING
STEP T001: EXECUTE block PROCEDURE body

- [IMPL-TOKEN_SYSTEM] [ARCH-RESOURCE_MANAGEMENT] [ARCH-TOKEN_SYSTEM] [REQ-DOC_016] [REQ-GOV_REGISTRY_COMPLETENESS] — How: replace legacy project token YAML with pointer stubs after migration to tied/semantic-tokens.yaml.

NOTE: Legacy YAML stub replacement is a migration process step; no dedicated Go implementation lead required.
