# bkpdir configuration (canonical)

**Scope:** Configuration discovery, file precedence, environment overrides, YAML keys, inheritance (`inherit`), merge prefixes (`+`, `^`, `!`, `=`), `exclude_patterns` semantics, and the `config` subcommand inspection modes. **Vocabulary only** — merge algorithms live in [`../../config.go`](../../config.go) and [`../../docs/configuration-inheritance.md`](../../docs/configuration-inheritance.md).

**Traceability:** [REQ-CONFIGURATION](../requirements/REQ-CONFIGURATION.yaml) · [REQ-CFG_001](../requirements/REQ-CFG_001.yaml) · [REQ-CFG_005](../requirements/REQ-CFG_005.yaml) · [REQ-CFG_006](../requirements/REQ-CFG_006.yaml) · [REQ-CONFIG_OUTPUT_GROUPING](../requirements/REQ-CONFIG_OUTPUT_GROUPING.yaml) · [REQ-TEST_EXCLUDE_MERGE](../requirements/REQ-TEST_EXCLUDE_MERGE.yaml) · [ARCH-CONFIG_SYSTEM](../architecture-decisions/ARCH-CONFIG_SYSTEM.yaml) · [ARCH-CFG_005](../architecture-decisions/ARCH-CFG_005.yaml) · [ARCH-CFG_006](../architecture-decisions/ARCH-CFG_006.yaml) · [ARCH-CONFIG_OUTPUT_GROUPING](../architecture-decisions/ARCH-CONFIG_OUTPUT_GROUPING.yaml) · [ARCH-EXCLUDE_MERGE_FIX](../architecture-decisions/ARCH-EXCLUDE_MERGE_FIX.yaml) · [ARCH-TEST_EXCLUDE_MERGE](../architecture-decisions/ARCH-TEST_EXCLUDE_MERGE.yaml) · [IMPL-CONFIG_STRUCT](../implementation-decisions/IMPL-CONFIG_STRUCT.yaml) · [IMPL-CFG_006](../implementation-decisions/IMPL-CFG_006.yaml) · [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) · [IMPL-CONFIG_OUTPUT_GROUPING](../implementation-decisions/IMPL-CONFIG_OUTPUT_GROUPING.yaml) · [IMPL-TEST_EXCLUDE_MERGE](../implementation-decisions/IMPL-TEST_EXCLUDE_MERGE.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-cli.md`](bkpdir-cli.md) · [`bkpdir-git-integration.md`](bkpdir-git-integration.md) · [`../../docs/configuration-inheritance.md`](../../docs/configuration-inheritance.md) · [`../../docs/configuration-examples.md`](../../docs/configuration-examples.md) · [`domain-references.md`](domain-references.md)

**Split/merge rule:** Configuration keys and merge semantics stay here; Git nested block details in [`bkpdir-git-integration.md`](bkpdir-git-integration.md); format string keys in [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md).

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **configuration discovery** | config load (alone) | Search path + ordered file processing |
| **search path** | config path list | Colon-separated list from `BKPDIR_CONFIG` or default |
| **earlier file wins** | first file priority | Sequential files: earlier overrides later |
| **inherit** | include config | Explicit inheritance chain via `inherit:` list |
| **merge prefix** | yaml merge key | `+` append, `^` prepend, `!` replace, `=` default-if-unset |
| **override** | replace (unprefixed) | Unprefixed scalar/list replaces prior value |
| **exclude_patterns** | excludes, ignore list | List-valued glob patterns; default merge for arrays |
| **single-segment directory pattern** | root-only exclude | e.g. `node_modules/` matches at any path depth |
| **source tracking** | provenance | Which file/env/default supplied each value |

---

## Naming bridge: discovery and storage

| Concept | Env / CLI | File | YAML key | Go type |
|---------|-----------|------|----------|---------|
| Config search path | `BKPDIR_CONFIG` | — | — | `GetConfigSearchPaths()` |
| Default search | (unset env) | `./.bkpdir.yml:~/.bkpdir.yml` | — | hard-coded default |
| Project config | — | `.bkpdir.yml` | all fields | `Config` |
| User config | — | `~/.bkpdir.yml` | all fields | `Config` |
| Template output | `bkpdir template -o` | `.bkpdir.default-YYYY-MM-DD.yml` | — | — |
| Field env override | `BKPDIR_<FIELD>` | — | maps to struct field | reflection-based |
| Inheritance | — | paths in `inherit:` | `inherit` | `[]string` |

---

## Core YAML keys (catalog)

| Key | Type | Default (typical) | Notes |
|-----|------|-------------------|-------|
| `archive_dir_path` | string | `../.bkpdir` | Directory archive storage |
| `use_current_dir_name` | bool | `true` | Include CWD name in archive path |
| `exclude_patterns` | list | `[]` | Glob patterns; array merge default |
| `include_git_info` | bool | legacy | Prefer `git.include_info` |
| `show_git_dirty_status` | bool | legacy | Prefer `git.show_dirty_status` |
| `skip_broken_symlinks` | bool | — | Traversal behavior |
| `inherit` | list | — | Explicit inheritance file paths |
| `git` | map | — | Nested Git block; see git glossary |
| `backup_dir_path` | string | — | File backup storage |
| `use_current_dir_name_for_files` | bool | — | File backup path layout |

Format, template, pattern, and status keys: [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) and [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md).

---

## Merge prefix catalog

| Prefix | Name | Array behavior | Scalar behavior |
|--------|------|----------------|-----------------|
| (none) | **override** | Replace entire list | Replace value |
| `+` | **merge** | Append child elements | N/A |
| `^` | **prepend** | Prepend child elements | N/A |
| `!` | **replace** | Explicit full replace | Replace value |
| `=` | **default** | Set only if unset | Set only if unset |

---

## `config` subcommand output modes

| Mode | Flag | Notes |
|------|------|-------|
| **table** | `--format table` | Default grouped table |
| **tree** | `--format tree` | Hierarchical view |
| **json** | `--format json` | Machine-readable |
| **flat** | `--flat` | Ungrouped field list |
| **overrides-only** | `--overrides-only` | Non-default values |
| **sources** | `--sources` | Per-value source attribution |

---

## Exclusion pattern matching (vocabulary)

| Rule | Example | Meaning |
|------|---------|---------|
| Trailing `/` | `node_modules/` | Directory segment |
| Single-segment at any depth | `node_modules/` | Matches as path component |
| Glob | `*.log`, `**/*.tmp` | Standard glob / doublestar |
| Unprefixed first-file list | first `.bkpdir.yml` | Replaces empty built-in default list |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Apply merge strategies | `APPLY_MERGE_STRATEGIES` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Merge operation dispatch | `APPLY_MERGE_OPERATION` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Array merge | `APPLY_MERGE` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Override merge | `APPLY_OVERRIDE` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Prepend merge | `APPLY_PREPEND` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Replace merge | `APPLY_REPLACE` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Default-if-unset | `APPLY_DEFAULT` | [IMPL-EXCLUDE_MERGE_FIX](../implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml) |
| Load with inheritance | `LOAD_CONFIG_WITH_INHERITANCE` | [IMPL-CFG_006](../implementation-decisions/IMPL-CFG_006.yaml) |
| Config output grouping | `CONFIG_OUTPUT_GROUPING` | [IMPL-CONFIG_OUTPUT_GROUPING](../implementation-decisions/IMPL-CONFIG_OUTPUT_GROUPING.yaml) |

---

## STDD proposed terms (not implemented in bkpdir)

Terms below describe a **planned** layered-config product pattern in STDD; bkpdir uses **merge-default arrays** per [REQ-CFG_005](../requirements/REQ-CFG_005.yaml), not unprefixed list replace:

| Term `(proposed)` | Notes |
|-------------------|-------|
| **project-local layer** | First project file layer in multi-product design |
| **unprefixed list replace** | STDD-only; conflicts with bkpdir merge semantics |

---

## Alphabetical index

| Term | Section |
|------|---------|
| archive_dir_path | Core YAML keys |
| BKPDIR_CONFIG | Naming bridge |
| config subcommand | config subcommand output modes |
| configuration discovery | Preferred terms |
| earlier file wins | Preferred terms |
| exclude_patterns | Preferred terms |
| flat | config subcommand output modes |
| inherit | Preferred terms |
| json | config subcommand output modes |
| merge prefix | Preferred terms |
| override | Merge prefix catalog |
| search path | Preferred terms |
| single-segment directory pattern | Preferred terms |
| source tracking | Preferred terms |
| table | config subcommand output modes |
| tree | config subcommand output modes |
| unprefixed list replace | STDD proposed terms |
