# bkpdir output formatting (canonical)

**Scope:** Printf-style (`format_*`), template-style (`template_*`), regex (`pattern_*`) configuration keys; placeholders; dual formatting; delayed output via `OutputCollector`; and file statistics placeholders. **Vocabulary only** — formatting algorithms live in [`../../pkg/formatter/`](../../pkg/formatter/) and [`../../docs/format-strings-reference.md`](../../docs/format-strings-reference.md).

**Traceability:** [REQ-OUTPUT_FORMATTING](../requirements/REQ-OUTPUT_FORMATTING.yaml) · [REQ-TEMPLATE_FORMATTING](../requirements/REQ-TEMPLATE_FORMATTING.yaml) · [REQ-CUSTOMIZABLE_FORMAT_STRINGS](../requirements/REQ-CUSTOMIZABLE_FORMAT_STRINGS.yaml) · [REQ-OUT_002](../requirements/REQ-OUT_002.yaml) · [REQ-IMMUTABLE_OUTPUT_FORMATTING](../requirements/REQ-IMMUTABLE_OUTPUT_FORMATTING.yaml) · [REQ-IMMUTABLE_TEMPLATE_FORMATTING](../requirements/REQ-IMMUTABLE_TEMPLATE_FORMATTING.yaml) · [ARCH-OUTPUT_FORMATTING](../architecture-decisions/ARCH-OUTPUT_FORMATTING.yaml) · [ARCH-FILE_STATISTICS](../architecture-decisions/ARCH-FILE_STATISTICS.yaml) · [ARCH-CUSTOMIZABLE_FORMAT_STRINGS](../architecture-decisions/ARCH-CUSTOMIZABLE_FORMAT_STRINGS.yaml) · [IMPL-DUAL_FORMATTING](../implementation-decisions/IMPL-DUAL_FORMATTING.yaml) · [IMPL-DELAYED_OUTPUT](../implementation-decisions/IMPL-DELAYED_OUTPUT.yaml) · [IMPL-FILE_STATISTICS](../implementation-decisions/IMPL-FILE_STATISTICS.yaml) · [IMPL-FILE_STATISTICS_TEMPLATE_FIX](../implementation-decisions/IMPL-FILE_STATISTICS_TEMPLATE_FIX.yaml) · [IMPL-CUSTOMIZABLE_FORMAT_STRINGS](../implementation-decisions/IMPL-CUSTOMIZABLE_FORMAT_STRINGS.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-configuration.md`](bkpdir-configuration.md) · [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md) · [`../../docs/format-strings-reference.md`](../../docs/format-strings-reference.md) · [`../../docs/delayed-output-usage.md`](../../docs/delayed-output-usage.md)

**Split/merge rule:** User-visible message keys here; exit status integers in [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md).

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **format string** | message template (alone) | Printf-style `format_*` YAML keys |
| **template string** | format template | Named-placeholder `template_*` keys |
| **pattern string** | regex config | `pattern_*` keys for parsing filenames |
| **dual formatting** | two formatters | Printf + template paths both supported |
| **delayed output** | buffered output | `OutputCollector` defers stdout/stderr |
| **file statistics** | file stats | Size/mtime enrichment for list output |
| **placeholder** | variable | Template tokens like `#{path}` |

---

## Key families (catalog)

| Prefix | Purpose | Example keys |
|--------|---------|--------------|
| `format_*` | Printf-style messages | `format_created_archive`, `format_diff_added` |
| `template_*` | Named-placeholder messages | `template_created_archive`, `template_list_archive` |
| `pattern_*` | Regex patterns | `pattern_archive_filename`, `pattern_timestamp` |

---

## Common placeholders

| Placeholder | Used in | Meaning |
|-------------|---------|---------|
| `#{path}` | templates | File or archive path |
| `#{size_human}` | templates | Human-readable size (e.g. `1.2MB`) |
| `#{mtime}` | templates | Modification time |
| `#{branch}` | templates | Git branch name |
| `#{hash}` | templates | Git short hash |
| `#{note}` | templates | User note |
| `#{operation}` | templates | Operation label |
| `#{message}` | templates | Error or status message |
| `%s`, `%d`, `%v` | format_* | Go printf verbs |

---

## Directory / archive format keys

| Key | Purpose |
|-----|---------|
| `format_created_archive` | Archive created success |
| `format_identical_archive` | Directory unchanged vs archive |
| `format_list_archive` | Archive list line |
| `format_dry_run_archive` | Dry-run archive message |
| `format_incremental_created` | Incremental created |
| `format_incremental_skipped_no_changes` | Incremental skipped |
| `format_no_archives_found` | Empty archive list |
| `format_diff_no_changes` | Diff: no changes |
| `format_diff_changes` | Diff: header |
| `format_diff_added` | Diff: added file |
| `format_diff_modified` | Diff: modified file |
| `format_diff_deleted` | Diff: deleted file |

Each `format_*` key has a parallel `template_*` counterpart where applicable.

---

## File backup format keys

| Key | Purpose |
|-----|---------|
| `format_created_backup` | Backup created |
| `format_identical_backup` | File unchanged |
| `format_list_backup` | Backup list line |
| `format_dry_run_backup` | Dry-run backup |
| `format_no_backups_found` | Empty backup list |

---

## Pattern keys

| Key | Purpose |
|-----|---------|
| `pattern_archive_filename` | Parse archive names |
| `pattern_backup_filename` | Parse backup names |
| `pattern_config_line` | Parse config display lines |
| `pattern_timestamp` | Parse timestamp segments |

---

## Formatter types (code catalog)

| Type | Package | Role |
|------|---------|------|
| `OutputFormatter` | `pkg/formatter` | Primary formatting API |
| `OutputCollector` | `pkg/formatter` | Delayed/buffered output |
| `FormatType` | `pkg/formatter` | Format string category enum |
| `PatternType` | `pkg/formatter` | Pattern category enum |
| `FileStatInfo` | `pkg/formatter` | Gathered stats for list lines |
| `GatherFileStatInfo()` | `pkg/formatter` | Collect size/mtime for templates |

Delayed output API: `FlushAll`, `FlushStdout`, `FlushStderr`, `Clear`.

---

## Naming bridge

| Concept | YAML key | Go field (Config) | Formatter API |
|---------|----------|-------------------|---------------|
| Created archive message | `format_created_archive` | `FormatCreatedArchive` | `FormatCreatedArchive` |
| Template variant | `template_created_archive` | `TemplateCreatedArchive` | template path |
| Human size placeholder | — | — | `#{size_human}` via `FileStatInfo` |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Dual format dispatch | `DUAL_FORMATTING` | [IMPL-DUAL_FORMATTING](../implementation-decisions/IMPL-DUAL_FORMATTING.yaml) |
| Delayed output flush | `FLUSH_ALL` | [IMPL-DELAYED_OUTPUT](../implementation-decisions/IMPL-DELAYED_OUTPUT.yaml) |
| Gather file statistics | `GATHER_FILE_STAT_INFO` | [IMPL-FILE_STATISTICS](../implementation-decisions/IMPL-FILE_STATISTICS.yaml) |
| Template placeholder replace | `PLACEHOLDER_REPLACE` | [IMPL-FILE_STATISTICS_TEMPLATE_FIX](../implementation-decisions/IMPL-FILE_STATISTICS_TEMPLATE_FIX.yaml) |
| Custom format string lookup | `CUSTOMIZABLE_FORMAT_STRINGS` | [IMPL-CUSTOMIZABLE_FORMAT_STRINGS](../implementation-decisions/IMPL-CUSTOMIZABLE_FORMAT_STRINGS.yaml) |

---

## Alphabetical index

| Term | Section |
|------|---------|
| delayed output | Preferred terms |
| dual formatting | Preferred terms |
| file statistics | Preferred terms |
| format string | Preferred terms |
| format_created_archive | Directory / archive format keys |
| format_diff_added | Directory / archive format keys |
| OutputCollector | Formatter types |
| pattern string | Preferred terms |
| pattern_archive_filename | Pattern keys |
| placeholder | Preferred terms |
| template string | Preferred terms |
| #{size_human} | Common placeholders |
