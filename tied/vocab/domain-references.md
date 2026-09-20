# Client vocabulary catalog

**Scope:** Index of client-owned domain vocabulary. TIED methodology vocabulary is cataloged separately under [`../methodology/vocab/domain-references.md`](../methodology/vocab/domain-references.md).

**Procedure:** Read [`routing.md`](routing.md) first. Use the methodology catalog for TIED concepts and this catalog for client product concepts.

---

## TIED methodology catalog

The refreshable TIED vocabulary catalog is [`../methodology/vocab/domain-references.md`](../methodology/vocab/domain-references.md).

## Client canonical glossaries

| Priority | Document | Scope |
|----------|----------|-------|
| 0 | [`domain-references.md`](domain-references.md) | This index |
| 1 | [`bkpdir-cli.md`](bkpdir-cli.md) | `bkpdir` CLI: subcommands, flags, auto-detection |
| 1b | [`bkpdir-configuration.md`](bkpdir-configuration.md) | Config discovery, YAML keys, merge/inheritance, `config` subcommand |
| 2 | [`bkpdir-archives-backup-diff.md`](bkpdir-archives-backup-diff.md) | ZIP archives, file backup, snapshots, diff |
| 2b | [`bkpdir-output-formatting.md`](bkpdir-output-formatting.md) | `format_*`, `template_*`, `pattern_*`, delayed output |
| 3 | [`bkpdir-git-integration.md`](bkpdir-git-integration.md) | Git detection, naming segments, dirty suffix |
| 3b | [`bkpdir-errors-resources.md`](bkpdir-errors-resources.md) | `status_*` codes, structured errors, atomic writes |
| 4 | [`immutable-contracts.md`](immutable-contracts.md) | Index of all `REQ-IMMUTABLE_*` contracts |
| 5 | [`tied-methodology.md`](tied-methodology.md) | TIED layout, semantic tokens, module validation, PROC-* names |
| 6 | [`tied-yaml-mcp.md`](tied-yaml-mcp.md) | TIED YAML MCP, `tied-cli`, validation/verify |
| 6b | [`feedback-to-tied.md`](feedback-to-tied.md) | Upstream feedback (`feedback.yaml`) |
| 7 | [`leap-proposal-queue.md`](leap-proposal-queue.md) | Non-canonical LEAP proposals |
| 8 | [`agentstream.md`](agentstream.md) | Go `agentstream` CLI and library |
| 8b | [`agent-stream-ruby.md`](agent-stream-ruby.md) | Ruby ATDD runner parity |
| 9 | [`pseudocode-and-citdp.md`](pseudocode-and-citdp.md) | Domain vocab vs IMPL grammar; CITDP |

---

## Authoring guides (not glossaries)

| Document | Role |
|----------|------|
| [`../../docs/vocabulary-index-analysis-and-standards.md`](../../docs/vocabulary-index-analysis-and-standards.md) | Meta-standard for glossary structure and TIED integration |
| [`../docs/pseudocode-writing-and-validation.md`](../docs/pseudocode-writing-and-validation.md) | IMPL pseudo-code lifecycle (not domain term registry) |
| [`../docs/implementation-decisions.md`](../docs/implementation-decisions.md) | IMPL grammar vocabulary (INPUT/OUTPUT/DATA) — distinct from domain vocab |

---

## Cross-topic notes

- **STDD / TIED repository layout:** canonical domain glossaries live at `tied/vocab/<topic>.md` (no `-vocabulary` filename suffix). Meta-standard: [`../../docs/vocabulary-index-analysis-and-standards.md`](../../docs/vocabulary-index-analysis-and-standards.md) § STDD convention.
- **bkpdir product vs STDD tooling:** product glossaries (`bkpdir-*`, `immutable-contracts`) cover the backup CLI; methodology/STDD glossaries cover TIED, MCP, agentstream, LEAP, CITDP.
- **agentstream** (Go) vs **agent-stream** (Ruby directory) vs **run-feature-batch** scripts — define once in [`agentstream.md`](agentstream.md) and [`agent-stream-ruby.md`](agent-stream-ruby.md).
- **Domain vocabulary** (this tree) vs **IMPL grammar vocabulary** (INPUT/OUTPUT/DATA) — define once in [`pseudocode-and-citdp.md`](pseudocode-and-citdp.md).
- **TIED base path** / **project YAML** vs **methodology YAML** — define once in [`tied-methodology.md`](tied-methodology.md).
- **Non-canonical LEAP proposals** (`leap-proposals/`) never mutate project TIED YAML — see [`leap-proposal-queue.md`](leap-proposal-queue.md).
- **Immutable contracts** index at [`immutable-contracts.md`](immutable-contracts.md) links each `REQ-IMMUTABLE_*` to its owning product glossary.

---

## Ownership

This catalog and all non-index glossaries in `tied/vocab/` are client-owned. The methodology catalog and its linked glossaries are refreshed under `tied/methodology/vocab/`.

## Alphabetical index

| Term | Section |
|------|---------|
| client canonical glossaries | Client canonical glossaries |
| client vocabulary catalog | Title |
| TIED methodology catalog | TIED methodology catalog |
| agent-stream | Cross-topic notes |
| agentstream | Cross-topic notes |
| bkpdir product glossaries | Cross-topic notes |
| Domain vocabulary index | Title |
| IMMUTABLE contracts | Cross-topic notes |
| IMPL grammar vocabulary | Authoring guides |
| STDD tooling glossaries | Cross-topic notes |
| sub-vocabulary-sync | Scope |
| VOCAB_INDEX | Scope |
