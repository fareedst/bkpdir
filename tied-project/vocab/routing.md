# Client vocabulary routing index

**Ownership:** Client-owned discovery handoff. TIED methodology vocabulary is refreshed under [`../tied-bundle/vocab/`](../tied-bundle/vocab/).

**Procedure:**
1. Read the [TIED methodology routing index](../tied-bundle/vocab/routing.md) for methodology terms.
2. Read the client glossary routing table below for product terms.
3. PRELOAD only the matched glossary for the task.
4. Use the [client vocabulary catalog](domain-references.md) for cross-topic links.

---

## TIED methodology vocabulary

Use [`../tied-bundle/vocab/routing.md`](../tied-bundle/vocab/routing.md) for TIED layout, process, validation, and tooling concepts. Do not copy methodology terms into client glossaries.

## Client glossary routing table

| Pri | File | Keywords / When to read |
|-----|------|------------------------|
| 1 | [bkpdir-cli.md](bkpdir-cli.md) | `bkpdir` CLI, subcommands, flags, auto-detection |
| 1b | [bkpdir-configuration.md](bkpdir-configuration.md) | Config discovery, YAML keys, merge/inheritance, `config` subcommand |
| 2 | [bkpdir-archives-backup-diff.md](bkpdir-archives-backup-diff.md) | ZIP archives, file backup, snapshots, diff |
| 2b | [bkpdir-output-formatting.md](bkpdir-output-formatting.md) | `format_*`, `template_*`, `pattern_*`, delayed output |
| 3 | [bkpdir-git-integration.md](bkpdir-git-integration.md) | Git detection, naming segments, dirty suffix |
| 3b | [bkpdir-errors-resources.md](bkpdir-errors-resources.md) | `status_*` codes, structured errors, atomic writes |
| 4 | [immutable-contracts.md](immutable-contracts.md) | `REQ-IMMUTABLE_*` contracts index |

## Ownership

Files under `tied-project/vocab/` are client-owned. Files under `tied-bundle/vocab/` are TIED-owned and are replaced during methodology refresh.

## Alphabetical index

| Term | Section |
|------|---------|
| client vocabulary routing index | Title |
| client glossary routing table | Client glossary routing table |
| TIED methodology vocabulary | TIED methodology vocabulary |
