# Immutable contracts index (canonical)

**Scope:** Index of all **`REQ-IMMUTABLE_*`** behavioral contracts — preferred terms, owning product glossary, and semantic token links. **Vocabulary only** — normative criteria live in each REQ detail file under [`../requirements/`](../requirements/).

**Traceability:** All `REQ-IMMUTABLE_*` tokens in [`../semantic-tokens.yaml`](../semantic-tokens.yaml) · [REQ-IMMUTABLE_FEATURE_PRESERVATION](../requirements/REQ-IMMUTABLE_FEATURE_PRESERVATION.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-cli.md`](bkpdir-cli.md) · [`bkpdir-configuration.md`](bkpdir-configuration.md) · [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) · [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) · [`bkpdir-git-integration.md`](bkpdir-git-integration.md) · [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md) · [`tied-methodology.md`](tied-methodology.md)

**Split/merge rule:** This file indexes immutables only; domain terms are defined in the linked product glossaries. Do not duplicate CLI or config catalogs here.

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **immutable contract** | frozen req, locked behavior | `REQ-IMMUTABLE_*` prefix |
| **behavioral guarantee** | backward compat (alone) | Must not change without explicit versioning |
| **owning glossary** | parent doc | Product vocab file for domain terms |

---

## Immutable REQ index

| Token | Preferred term | Owning glossary |
|-------|----------------|-----------------|
| [REQ-IMMUTABLE_ARCHIVE_NAMING](../requirements/REQ-IMMUTABLE_ARCHIVE_NAMING.yaml) | **archive naming grammar** | [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) |
| [REQ-IMMUTABLE_CLI_COMMANDS](../requirements/REQ-IMMUTABLE_CLI_COMMANDS.yaml) | **CLI command surface** | [`bkpdir-cli.md`](bkpdir-cli.md) |
| [REQ-IMMUTABLE_GLOBAL_OPTIONS](../requirements/REQ-IMMUTABLE_GLOBAL_OPTIONS.yaml) | **global CLI flags** | [`bkpdir-cli.md`](bkpdir-cli.md) |
| [REQ-IMMUTABLE_CONFIGURATION_DEFAULTS](../requirements/REQ-IMMUTABLE_CONFIGURATION_DEFAULTS.yaml) | **configuration defaults** | [`bkpdir-configuration.md`](bkpdir-configuration.md) |
| [REQ-IMMUTABLE_DIRECTORY_OPERATIONS](../requirements/REQ-IMMUTABLE_DIRECTORY_OPERATIONS.yaml) | **directory archive operations** | [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) |
| [REQ-IMMUTABLE_ERROR_HANDLING](../requirements/REQ-IMMUTABLE_ERROR_HANDLING.yaml) | **error handling contracts** | [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md) |
| [REQ-IMMUTABLE_FEATURE_PRESERVATION](../requirements/REQ-IMMUTABLE_FEATURE_PRESERVATION.yaml) | **feature preservation** | (cross-cutting; all product glossaries) |
| [REQ-IMMUTABLE_FILE_BACKUP_NAMING](../requirements/REQ-IMMUTABLE_FILE_BACKUP_NAMING.yaml) | **file backup naming** | [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) |
| [REQ-IMMUTABLE_FILE_BACKUP_OPERATIONS](../requirements/REQ-IMMUTABLE_FILE_BACKUP_OPERATIONS.yaml) | **file backup operations** | [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) |
| [REQ-IMMUTABLE_FILE_EXCLUSION](../requirements/REQ-IMMUTABLE_FILE_EXCLUSION.yaml) | **file exclusion behavior** | [`bkpdir-configuration.md`](bkpdir-configuration.md) |
| [REQ-IMMUTABLE_GIT_INTEGRATION](../requirements/REQ-IMMUTABLE_GIT_INTEGRATION.yaml) | **Git integration contracts** | [`bkpdir-git-integration.md`](bkpdir-git-integration.md) |
| [REQ-IMMUTABLE_OUTPUT_FORMATTING](../requirements/REQ-IMMUTABLE_OUTPUT_FORMATTING.yaml) | **output formatting defaults** | [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) |
| [REQ-IMMUTABLE_TEMPLATE_FORMATTING](../requirements/REQ-IMMUTABLE_TEMPLATE_FORMATTING.yaml) | **template formatting defaults** | [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) |
| [REQ-IMMUTABLE_RESOURCE_MANAGEMENT](../requirements/REQ-IMMUTABLE_RESOURCE_MANAGEMENT.yaml) | **resource management contracts** | [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md) |
| [REQ-IMMUTABLE_BUILD_SYSTEM](../requirements/REQ-IMMUTABLE_BUILD_SYSTEM.yaml) | **build system contracts** | [`tied-methodology.md`](tied-methodology.md) |
| [REQ-IMMUTABLE_CODE_QUALITY](../requirements/REQ-IMMUTABLE_CODE_QUALITY.yaml) | **code quality gates** | [`tied-methodology.md`](tied-methodology.md) |
| [REQ-IMMUTABLE_PERFORMANCE](../requirements/REQ-IMMUTABLE_PERFORMANCE.yaml) | **performance contracts** | (cross-cutting) |
| [REQ-IMMUTABLE_PLATFORM_COMPATIBILITY](../requirements/REQ-IMMUTABLE_PLATFORM_COMPATIBILITY.yaml) | **platform compatibility** | (cross-cutting) |
| [REQ-IMMUTABLE_TESTING_INFRASTRUCTURE](../requirements/REQ-IMMUTABLE_TESTING_INFRASTRUCTURE.yaml) | **testing infrastructure** | [`tied-methodology.md`](tied-methodology.md) |

---

## Naming bridge

| Concept | REQ prefix | Detail path pattern |
|---------|------------|---------------------|
| Immutable contract | `REQ-IMMUTABLE_*` | `tied/requirements/REQ-IMMUTABLE_*.yaml` |
| User spec pointer | — | [`../../docs/user/specification.md`](../../docs/user/specification.md) § Normative traceability |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| (index only; blocks live in product IMPLs) | — | See owning glossary |

---

## Alphabetical index

| Term | Section |
|------|---------|
| archive naming grammar | Immutable REQ index |
| behavioral guarantee | Preferred terms |
| CLI command surface | Immutable REQ index |
| configuration defaults | Immutable REQ index |
| error handling contracts | Immutable REQ index |
| feature preservation | Immutable REQ index |
| file backup naming | Immutable REQ index |
| file exclusion behavior | Immutable REQ index |
| global CLI flags | Immutable REQ index |
| immutable contract | Preferred terms |
| owning glossary | Preferred terms |
| template formatting defaults | Immutable REQ index |
