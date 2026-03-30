# AI Assistant Context Documentation Index

This project uses **TIED** (traceability via `tied/` YAML + semantic tokens). Normative requirements, architecture, and implementation live in **TIED detail files**, not in duplicate prose here.

## Canonical sources

1. **[AGENTS.md](../../AGENTS.md)** — Agent operating guide (methodology, MCP-first TIED edits, checklists).
2. **[ai-principles.md](../../ai-principles.md)** — AI-first principles and validation expectations.
3. **[tied/requirements.yaml](../../tied/requirements.yaml)** — REQ index; per-token detail under [tied/requirements/](../../tied/requirements/).
4. **[tied/architecture-decisions.yaml](../../tied/architecture-decisions.yaml)** — ARCH index; detail under [tied/architecture-decisions/](../../tied/architecture-decisions/).
5. **[tied/implementation-decisions.yaml](../../tied/implementation-decisions.yaml)** — IMPL index; detail under [tied/implementation-decisions/](../../tied/implementation-decisions/) (includes `essence_pseudocode`).
6. **[tied/semantic-tokens.yaml](../../tied/semantic-tokens.yaml)** and **[tied/semantic-tokens.md](../../tied/semantic-tokens.md)** — Token registry and guide.
7. **Human-readable digests** (summaries of the above): [tied/requirements.md](../../tied/requirements.md), [tied/architecture-decisions.md](../../tied/architecture-decisions.md), [tied/implementation-decisions.md](../../tied/implementation-decisions.md).

Optional task tracking: [tied/tasks.md](../../tied/tasks.md) (optional per TIED 2.2.0).

## MCP (optional)

For reading/writing TIED YAML with validation, see [tied/docs/ai-agent-tied-mcp-usage.md](../../tied/docs/ai-agent-tied-mcp-usage.md). Use **`tied_validate_consistency`** after substantive TIED changes.

## User-facing docs

- **[../user/specification.md](../user/specification.md)** — Features and usage narrative (normative rules link to `tied/requirements/`).

## Development flow (summary)

Before code changes: read `AGENTS.md` and `ai-principles.md`; locate related `[REQ-*]`, `[ARCH-*]`, `[IMPL-*]` in `tied/`; extend pseudo-code and decisions before tests/code per [tied/docs/agent-req-implementation-checklist.md](../../tied/docs/agent-req-implementation-checklist.md).

### Code comment pattern

```go
// [REQ-FILE_BACKUP] Create backup of single file with comparison
// [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
func CreateFileBackup(cfg *Config, filePath string, note string, dryRun bool) error {
    // ...
}
```

### Quick validation

```bash
make test && make lint
```

Use `./scripts/validate_tokens.sh` or project Makefile targets for token enforcement when documented in the root README.

---

**Last Updated**: 2026-03-29
