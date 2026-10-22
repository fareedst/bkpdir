# AI Assistant Context Documentation Index

This project uses **TIED** (traceability via `tied-project/` YAML + semantic tokens). Normative requirements, architecture, and implementation live in **TIED detail files**, not in duplicate prose here.

## Canonical sources

1. **[AGENTS.md](../../AGENTS.md)** — Agent operating guide (methodology, MCP-first TIED edits, checklists).
2. **[ai-principles.md](../../ai-principles.md)** — AI-first principles and validation expectations.
3. **[tied-project/requirements.yaml](../../tied-project/requirements.yaml)** — REQ index; per-token detail under [tied-project/requirements/](../../tied-project/requirements/).
4. **[tied-project/architecture-decisions.yaml](../../tied-project/architecture-decisions.yaml)** — ARCH index; detail under [tied-project/architecture-decisions/](../../tied-project/architecture-decisions/).
5. **[tied-project/implementation-decisions.yaml](../../tied-project/implementation-decisions.yaml)** — IMPL index; detail under [tied-project/implementation-decisions/](../../tied-project/implementation-decisions/) (includes `essence_pseudocode`).
6. **[tied-project/semantic-tokens.yaml](../../tied-project/semantic-tokens.yaml)** and **[tied-bundle/docs/semantic-tokens.md](../../tied-bundle/docs/semantic-tokens.md)** — Token registry and guide.
7. **Human-readable digests** (summaries of the above): [tied-project/requirements.md](../../tied-project/requirements.md), [tied-project/architecture-decisions.md](../../tied-project/architecture-decisions.md), [tied-project/implementation-decisions.md](../../tied-project/implementation-decisions.md).

Optional task tracking: [tied-bundle/docs/tasks.md](../../tied-bundle/docs/tasks.md) (optional per TIED 2.2.0).

## MCP (optional)

For reading/writing TIED YAML with validation, see [tied-bundle/docs/ai-agent-tied-mcp-usage.md](../../tied-bundle/docs/ai-agent-tied-mcp-usage.md). Use **`tied_validate_consistency`** after substantive TIED changes.

## User-facing docs

- **[../user/specification.md](../user/specification.md)** — Features and usage narrative (normative rules link to `tied-project/requirements/`).

## Development flow (summary)

Before code changes: read `AGENTS.md` and `ai-principles.md`; locate related `[REQ-*]`, `[ARCH-*]`, `[IMPL-*]` in `tied-project/`; extend pseudo-code and decisions before tests/code per [tied-bundle/docs/agent-req-implementation-checklist.md](../../tied-bundle/docs/agent-req-implementation-checklist.md).

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
