# REQ → ARCH / IMPL cross-links

**Do not maintain a duplicate mapping table here.** Traceability is canonical in TIED:

- Each **requirement** detail file under [`tied/requirements/`](../../tied/requirements/) includes `traceability.architecture`, `traceability.implementation`, and related fields.
- Architecture and implementation detail files under [`tied/architecture-decisions/`](../../tied/architecture-decisions/) and [`tied/implementation-decisions/`](../../tied/implementation-decisions/) list related REQ/ARCH/IMPL tokens in `cross_references` / narrative as applicable.

**How to explore**

- Read the YAML indexes: [`tied/requirements.yaml`](../../tied/requirements.yaml), [`tied/architecture-decisions.yaml`](../../tied/architecture-decisions.yaml), [`tied/implementation-decisions.yaml`](../../tied/implementation-decisions.yaml).
- With **TIED MCP**, use tools such as `get_decisions_for_requirement` / `get_requirements_for_decision` (see [`tied/docs/ai-agent-tied-mcp-usage.md`](../../tied/docs/ai-agent-tied-mcp-usage.md)).
- Run **`tied_validate_consistency`** after edits to TIED YAML.
