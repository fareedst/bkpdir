# bkpdir errors and resources (canonical)

**Scope:** Structured error types, configurable exit status codes (`status_*` YAML keys), atomic file operations, resource cleanup, and context-aware error handling. **Vocabulary only** — error handling algorithms live in [`../../pkg/errors/`](../../pkg/errors/), [`../../pkg/resources/`](../../pkg/resources/), and [`../../errors.go`](../../errors.go).

**Traceability:** [REQ-ERROR_HANDLING](../requirements/REQ-ERROR_HANDLING.yaml) · [REQ-RESOURCE_MANAGEMENT](../requirements/REQ-RESOURCE_MANAGEMENT.yaml) · [REQ-CONTEXT_SUPPORT](../requirements/REQ-CONTEXT_SUPPORT.yaml) · [REQ-IMMUTABLE_ERROR_HANDLING](../requirements/REQ-IMMUTABLE_ERROR_HANDLING.yaml) · [REQ-IMMUTABLE_RESOURCE_MANAGEMENT](../requirements/REQ-IMMUTABLE_RESOURCE_MANAGEMENT.yaml) · [ARCH-ERROR_HANDLING](../architecture-decisions/ARCH-ERROR_HANDLING.yaml) · [ARCH-RESOURCE_MANAGEMENT](../architecture-decisions/ARCH-RESOURCE_MANAGEMENT.yaml) · [ARCH-CONTEXT_SUPPORT](../architecture-decisions/ARCH-CONTEXT_SUPPORT.yaml) · [IMPL-STRUCTURED_ERRORS](../implementation-decisions/IMPL-STRUCTURED_ERRORS.yaml) · [IMPL-ATOMIC_OPS](../implementation-decisions/IMPL-ATOMIC_OPS.yaml) · [IMPL-RESOURCE_MANAGER](../implementation-decisions/IMPL-RESOURCE_MANAGER.yaml) · [IMPL-CONTEXT_OPS](../implementation-decisions/IMPL-CONTEXT_OPS.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) · [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) · [`immutable-contracts.md`](immutable-contracts.md)

**Split/merge rule:** Exit codes and error types here; user-facing error message format keys in [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md).

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **status code** | exit code (alone) | Configurable integer exit values via `status_*` keys |
| **structured error** | wrapped error (alone) | Typed errors with classification |
| **atomic write** | temp file rename | Write-to-temp then rename pattern |
| **resource manager** | cleanup handler | Panic recovery and deferred cleanup |
| **context support** | ctx (alone) | Context-aware cancellation in operations |

---

## Directory operation status codes

| YAML key | Typical meaning |
|----------|-----------------|
| `status_created_archive` | Archive created successfully |
| `status_failed_to_create_archive_directory` | Cannot create archive dir |
| `status_directory_is_identical_to_existing_archive` | No changes; skip create |
| `status_directory_not_found` | Source directory missing |
| `status_invalid_directory_type` | Path not a directory |
| `status_permission_denied` | Permission error |
| `status_disk_full` | Insufficient disk space |
| `status_config_error` | Configuration load/validation failure |

---

## File backup status codes

| YAML key | Typical meaning |
|----------|-----------------|
| `status_created_backup` | Backup created successfully |
| `status_failed_to_create_backup_directory` | Cannot create backup dir |
| `status_file_is_identical_to_existing_backup` | File unchanged; skip |
| `status_file_not_found` | Source file missing |
| `status_invalid_file_type` | Path not a regular file |

---

## Error format keys (cross-reference)

User-facing error messages use `format_*` keys documented in [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md): `format_disk_full_error`, `format_permission_error`, `format_directory_not_found`, `format_file_not_found`, `format_invalid_directory`, `format_invalid_file`, `format_failed_write_temp`, `format_failed_finalize_file`, etc.

---

## Error types (code catalog)

| Symbol | Package | Role |
|--------|---------|------|
| `ErrorTypeDiskFull` | `pkg/errors` | Disk space classification |
| `ErrorTypePermission` | `pkg/errors` | Permission classification |
| `HandleArchiveError()` | `main` | Archive error dispatch |
| `ResourceManagerInterface` | `pkg/resources` | Cleanup contract |

---

## Resource concepts

| Term | Meaning |
|------|---------|
| **atomic rename** | Finalize via `rename` after temp write |
| **temp file** | Intermediate path during write |
| **panic recovery** | Resource manager catches panics |
| **deferred cleanup** | Close/remove temp resources on exit |

---

## Naming bridge

| Concept | YAML key | Go field | Handler |
|---------|----------|----------|---------|
| Identical directory exit | `status_directory_is_identical_to_existing_archive` | `StatusDirectoryIsIdenticalToExistingArchive` | archive create path |
| Disk full exit | `status_disk_full` | `StatusDiskFull` | `ErrorTypeDiskFull` |
| Atomic finalize | — | — | `IMPL-ATOMIC_OPS` |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Structured error classify | `CLASSIFY_ERROR` | [IMPL-STRUCTURED_ERRORS](../implementation-decisions/IMPL-STRUCTURED_ERRORS.yaml) |
| Handle archive error | `HANDLE_ARCHIVE_ERROR` | [IMPL-STRUCTURED_ERRORS](../implementation-decisions/IMPL-STRUCTURED_ERRORS.yaml) |
| Atomic file write | `ATOMIC_WRITE` | [IMPL-ATOMIC_OPS](../implementation-decisions/IMPL-ATOMIC_OPS.yaml) |
| Resource cleanup | `RESOURCE_CLEANUP` | [IMPL-RESOURCE_MANAGER](../implementation-decisions/IMPL-RESOURCE_MANAGER.yaml) |
| Context cancel check | `CONTEXT_CANCEL_CHECK` | [IMPL-CONTEXT_OPS](../implementation-decisions/IMPL-CONTEXT_OPS.yaml) |

---

## Alphabetical index

| Term | Section |
|------|---------|
| atomic rename | Resource concepts |
| atomic write | Preferred terms |
| context support | Preferred terms |
| resource manager | Preferred terms |
| status code | Preferred terms |
| status_config_error | Directory operation status codes |
| status_created_archive | Directory operation status codes |
| status_disk_full | Directory operation status codes |
| status_file_not_found | File backup status codes |
| structured error | Preferred terms |
| temp file | Resource concepts |
