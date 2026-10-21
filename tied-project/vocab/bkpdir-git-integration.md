# bkpdir Git integration (canonical)

**Scope:** Git repository detection, branch/hash inclusion in archive names, dirty working-tree suffix (`-dirty`), legacy top-level YAML keys vs nested `git:` block, and `pkg/git` types. **Vocabulary only** — Git command execution lives in [`../../pkg/git/`](../../pkg/git/) and [`../../config.go`](../../config.go).

**Traceability:** [REQ-GIT_INTEGRATION](../requirements/REQ-GIT_INTEGRATION.yaml) · [REQ-IMMUTABLE_GIT_INTEGRATION](../requirements/REQ-IMMUTABLE_GIT_INTEGRATION.yaml) · [ARCH-GIT_INTEGRATION](../architecture-decisions/ARCH-GIT_INTEGRATION.yaml) · [IMPL-GIT_CLI](../implementation-decisions/IMPL-GIT_CLI.yaml) · [IMPL-GIT_DIRTY_CONFIG](../implementation-decisions/IMPL-GIT_DIRTY_CONFIG.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**See also:** [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) · [`bkpdir-configuration.md`](bkpdir-configuration.md) · [`../../docs/user/specification.md`](../../docs/user/specification.md)

**Split/merge rule:** Git vocabulary here; archive naming grammar in [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md); general config keys in [`bkpdir-configuration.md`](bkpdir-configuration.md).

---

## Preferred terms vs synonyms

| Preferred | Avoid | Notes |
|-----------|-------|-------|
| **Git integration** | git support (alone) | Optional repo detection and metadata |
| **nested git block** | git section | `git:` map in YAML |
| **legacy Git keys** | old git flags | Top-level `include_git_info`, `show_git_dirty_status` |
| **dirty suffix** | dirty flag | `-dirty` appended to archive name when WIP changes |
| **branch and hash segment** | git tag | `=BRANCH=HASH` in archive filename |
| **auto-detect repo** | find git | Walk for `.git` when enabled |

---

## Legacy vs nested keys

| Legacy key (top-level) | Nested key (`git:`) | Purpose |
|------------------------|---------------------|---------|
| `include_git_info` | `include_info` | Include branch/hash in names |
| `show_git_dirty_status` | `show_dirty_status` | Append `-dirty` when dirty |

Prefer **nested `git:` block** for new configuration; legacy keys retained for backward compatibility.

---

## Nested `git:` block catalog

| Key | Type | Purpose |
|-----|------|---------|
| `enabled` | bool | Master Git integration switch |
| `include_info` | bool | Include branch/hash (legacy: `include_git_info`) |
| `show_dirty_status` | bool | Dirty suffix (legacy: `show_git_dirty_status`) |
| `command` | string | Git binary path (default `git`) |
| `working_directory` | string | CWD for Git commands |
| `require_clean_repo` | bool | Fail if working tree dirty |
| `auto_detect_repo` | bool | Auto-detect repository |
| `include_submodules` | bool | Submodule metadata |
| `include_branch` | bool | Include branch name |
| `include_hash` | bool | Include short commit hash |
| `include_status` | bool | Include working directory status |
| `command_timeout` | duration | Git command timeout |
| `max_submodule_depth` | int | Submodule walk limit |

---

## Git commands (catalog)

| Command | Purpose |
|---------|---------|
| `git rev-parse --is-inside-work-tree` | Detect repo |
| `git rev-parse --abbrev-ref HEAD` | Branch name |
| `git rev-parse --short HEAD` | Short hash |
| `git status --porcelain` | Dirty detection |

---

## `pkg/git` types

| Type | Role |
|------|------|
| `Repository` | Discovered repo handle |
| `Info` | Branch, hash, dirty flag |
| `Discover()` | Locate repo from path |

---

## Naming bridge

| Concept | Archive filename | Config | Code |
|---------|------------------|--------|------|
| Branch/hash | `=main=abc123` | `git.include_branch`, `git.include_hash` | `GitConfig.IncludeBranch` |
| Dirty suffix | `-dirty` | `git.show_dirty_status` | `GitConfig.ShowDirtyStatus` |
| Disable Git | (no segment) | `git.enabled: false` | `GitConfig.Enabled` |

---

## Pseudo-code block names

| Preferred term | UPPER_SNAKE block | Owning IMPL |
|----------------|-------------------|-------------|
| Git CLI discovery | `GIT_DISCOVER` | [IMPL-GIT_CLI](../implementation-decisions/IMPL-GIT_CLI.yaml) |
| Git info extraction | `GIT_INFO_EXTRACT` | [IMPL-GIT_CLI](../implementation-decisions/IMPL-GIT_CLI.yaml) |
| Dirty status config | `GIT_DIRTY_CONFIG` | [IMPL-GIT_DIRTY_CONFIG](../implementation-decisions/IMPL-GIT_DIRTY_CONFIG.yaml) |
| Append dirty suffix | `APPEND_DIRTY_SUFFIX` | [IMPL-GIT_DIRTY_CONFIG](../implementation-decisions/IMPL-GIT_DIRTY_CONFIG.yaml) |

---

## Alphabetical index

| Term | Section |
|------|---------|
| auto-detect repo | Preferred terms |
| branch and hash segment | Preferred terms |
| dirty suffix | Preferred terms |
| Git integration | Preferred terms |
| include_git_info | Legacy vs nested keys |
| legacy Git keys | Preferred terms |
| nested git block | Preferred terms |
| show_git_dirty_status | Legacy vs nested keys |
| -dirty | Naming bridge |
