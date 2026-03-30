# Package Interdependency Mapping — EXTRACT-008

**Tokens:** [REQ-EXTRACT_008_INTERDEP_MAPPING] [ARCH-EXTRACT_008_INTERDEP] [IMPL-EXTRACT_008_DOC_MIGRATION]

> **EXTRACT-008**: Canonical **map** of extracted packages and dependencies. Extended **narrative, patterns, and examples** live in [integration-guide.md](integration-guide.md); **architecture rationale** lives in TIED (do not duplicate here).

## Normative architecture (TIED)

- [ARCH-PACKAGE_EXTRACTION](../tied/architecture-decisions/ARCH-PACKAGE_EXTRACTION.yaml)
- [ARCH-SYSTEM_COMPONENTS](../tied/architecture-decisions/ARCH-SYSTEM_COMPONENTS.yaml)
- [REQ-EXTRACT_008_INTERDEP_MAPPING](../tied/requirements/REQ-EXTRACT_008_INTERDEP_MAPPING.yaml)

## Quick reference

| Package | Purpose | External deps | Internal deps |
|---------|---------|---------------|---------------|
| **pkg/config** | Configuration | yaml.v3 | — |
| **pkg/errors** | Structured errors | — | — |
| **pkg/resources** | Resource lifecycle | — | pkg/errors |
| **pkg/formatter** | Output formatting | — | pkg/errors |
| **pkg/git** | Git helpers | — | pkg/errors, pkg/formatter |
| **pkg/cli** | CLI framework | cobra | pkg/errors, pkg/fileops, pkg/formatter, … |
| **pkg/fileops** | File operations | doublestar | pkg/errors |
| **pkg/processing** | Concurrency / pipelines | — | pkg/errors |

Line counts change over time; treat this table as a **map**, not a spec.

## Dependency sketch

```
pkg/config ──┐
             ├─→ pkg/errors ──┐
pkg/git ─────┤                ├─→ pkg/resources
             └─→ pkg/formatter ┘
                      │
pkg/cli ──────────────┼─→ pkg/fileops
                      │
pkg/processing ───────┘
```

## Diagram asset

Static SVG: [images/package-interdependency-mapping.svg](images/package-interdependency-mapping.svg).

## Where to read more

- **[integration-guide.md](integration-guide.md)** — Quick start, patterns, troubleshooting, longer examples.
- **Per-package READMEs** under `pkg/*/README.md`.
- **Examples**: [examples/](examples/).

## Deliverables (EXTRACT-008)

- This file — quick map + TIED pointers (REQ-named deliverable).
- [integration-guide.md](integration-guide.md) — integration narrative and code samples.
- `docs/images/package-interdependency-mapping.svg` — visual diagram.
- Tokens cross-referenced in `tied/semantic-tokens.yaml` / `tied/semantic-tokens.md`.

## Remaining actions

- Optional CI: import-cycle checks and diagram presence (`go list -deps` / script).
- Run project token validation scripts when changing package boundaries (see [README.md](../README.md)).
