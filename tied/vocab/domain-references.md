# Domain vocabulary index (canonical)

**Scope:** Single-page directory for all domain vocabulary glossaries under `tied/vocab/`. Lists priority, scope, and cross-topic notes. This page is an **index only** — canonical terms live in the linked sibling files. Algorithms and step-by-step behavior stay in `tied/implementation-decisions/*-pseudocode.md`.

**Traceability:** [PROC-VOCABULARY_INDEX](../docs/processes.md) · [REQ-TIED_SETUP](../requirements/REQ-TIED_SETUP.yaml) · [REQ-GOV_DISCOVERABILITY](../requirements/REQ-GOV_DISCOVERABILITY.yaml)

**Help coverage:** N/A (no in-app Help in this repository)

**Checklist path:** [`../docs/agent-req-implementation-checklist.yaml`](../docs/agent-req-implementation-checklist.yaml) sets `VOCAB_INDEX: ./tied/vocab`. Agents **CALL** `sub-vocabulary-sync` (RESOLVE before naming/writing; RECORD when concepts are generated or artifacts change) per [`../docs/processes.md`](../docs/processes.md) § `[PROC-VOCABULARY_INDEX]`.

**Standards:** [`../../docs/vocabulary-index-analysis-and-standards.md`](../../docs/vocabulary-index-analysis-and-standards.md).

**See also:** [`tied-methodology.md`](tied-methodology.md) · [`immutable-contracts.md`](immutable-contracts.md) · [`../../docs/vocabulary-index-analysis-and-standards.md`](../../docs/vocabulary-index-analysis-and-standards.md)

---

## Canonical glossaries

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

## Alphabetical index

| Term | Section |
|------|---------|
| agent-stream | Cross-topic notes |
| agentstream | Cross-topic notes |
| bkpdir product glossaries | Cross-topic notes |
| Domain vocabulary index | Title |
| IMMUTABLE contracts | Cross-topic notes |
| IMPL grammar vocabulary | Authoring guides |
| STDD tooling glossaries | Cross-topic notes |
| sub-vocabulary-sync | Scope |
| VOCAB_INDEX | Scope |
