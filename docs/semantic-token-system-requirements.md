# Semantic token system (pointer)

Normative requirements for traceability and tokens are maintained in **TIED**, not in this file.

## Canonical sources

- **[REQ-DOC_016](../tied-project/requirements/REQ-DOC_016.yaml)** — AI-first token/traceability requirements (status and criteria in YAML).
- **[ARCH-TOKEN_SYSTEM](../tied-project/architecture-decisions/ARCH-TOKEN_SYSTEM.yaml)** — Token system architecture.
- **[IMPL-TOKEN_SYSTEM](../tied-project/implementation-decisions/IMPL-TOKEN_SYSTEM.yaml)** — Implementation approach and `essence_pseudocode`.
- **[tied-project/semantic-tokens.yaml](../tied-project/semantic-tokens.yaml)** and **[tied-bundle/docs/semantic-tokens.md](../tied-bundle/docs/semantic-tokens.md)** — Registry and human guide.
- **[AGENTS.md](../AGENTS.md)** and **[ai-principles.md](../ai-principles.md)** — Agent rules (`[PROC-IMPL_PSEUDOCODE_TOKENS]`, validation expectations).

## Legacy bridge

**[project-tokens.yaml](../project-tokens.yaml)** is a short pointer to `tied-project/` (indexes and `tied-project/semantic-tokens.yaml`); do not add token matrices there. The **canonical** registry for `[REQ-*]` / `[ARCH-*]` / `[IMPL-*]` is `tied-project/semantic-tokens.yaml`.

## Validation

Use Makefile targets and scripts described in the root **[README.md](../README.md)** (for example `make validate-token-enforcement`, `./scripts/validate_tokens.sh`, and TIED MCP **`tied_validate_consistency`** when editing TIED YAML).
