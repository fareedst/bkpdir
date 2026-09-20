# [IMPL-STDD_VIS_DATA_PIPELINE] [ARCH-STDD_VIS_FLOW] [REQ-STDD_VIS] [REQ-MODULE_VALIDATION]

## Summary contract

Extract and normalize STDD documents into a graph JSON and sample CSV consumed by visualization modules without re-parsing.

INPUT: semantic-tokens.md, requirements/architecture/implementation docs, code/test anchors
OUTPUT: docs/data/stdd-trace.json, docs/data/stdd-samples.csv
DATA: normalized nodes by REQ/ARCH/IMPL/TEST/CODE, cross-reference edges, token metadata

## INGEST_SOURCES

SPEC-ID: IMPL-STDD_VIS_DATA_PIPELINE::INGEST_SOURCES
STEP T001: Load token registry and markdown indexes plus sampled code/test anchors as pipeline inputs

- [IMPL-STDD_VIS_DATA_PIPELINE] [ARCH-STDD_VIS_FLOW] [REQ-STDD_VIS] [REQ-MODULE_VALIDATION] — How: load token registry and markdown indexes plus sampled code/test anchors as pipeline inputs.

PROCEDURE INGEST_SOURCES():
  READ semantic-tokens.md and REQ/ARCH/IMPL indexes
  SAMPLE code and test anchor references

## NORMALIZE_GRAPH

SPEC-ID: IMPL-STDD_VIS_DATA_PIPELINE::NORMALIZE_GRAPH
STEP T001: Classify nodes by token type and emit directed edges for documented cross-references

- [IMPL-STDD_VIS_DATA_PIPELINE] [ARCH-STDD_VIS_FLOW] [REQ-STDD_VIS] [REQ-MODULE_VALIDATION] — How: classify nodes by token type and emit directed edges for documented cross-references.

PROCEDURE NORMALIZE_GRAPH(sources):
  CREATE nodes with type, status, sample refs
  EMIT edges for each documented cross-reference

## EXPORT_ARTIFACTS

SPEC-ID: IMPL-STDD_VIS_DATA_PIPELINE::EXPORT_ARTIFACTS
STEP T001: Write stdd-trace.json graph and stdd-samples.csv chain samples for downstream visual modules

- [IMPL-STDD_VIS_DATA_PIPELINE] [ARCH-STDD_VIS_FLOW] [REQ-STDD_VIS] [REQ-MODULE_VALIDATION] — How: write stdd-trace.json graph and stdd-samples.csv chain samples for downstream visual modules.

PROCEDURE EXPORT_ARTIFACTS(graph):
  WRITE docs/data/stdd-trace.json
  WRITE docs/data/stdd-samples.csv with sample chains

## PIPELINE_VALIDATION

SPEC-ID: IMPL-STDD_VIS_DATA_PIPELINE::PIPELINE_VALIDATION
STEP T001: Schema-check exports, compare node counts to token registry, and spot-check two sample chains

- [IMPL-STDD_VIS_DATA_PIPELINE] [ARCH-STDD_VIS_FLOW] [REQ-STDD_VIS] [REQ-MODULE_VALIDATION] — How: schema-check exports, compare node counts to token registry, and spot-check two sample chains.

PROCEDURE PIPELINE_VALIDATION(exports):
  VALIDATE JSON schema
  ASSERT node counts consistent with registry
  SPOT-CHECK two REQ→CODE chains
