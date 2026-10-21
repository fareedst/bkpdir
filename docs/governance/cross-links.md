# REQ → ARCH / IMPL cross-links

**Do not maintain a duplicate mapping table here.** Traceability is canonical in TIED:

- Each **requirement** detail file under [`tied-project/requirements/`](../../tied-project/requirements/) includes `traceability.architecture`, `traceability.implementation`, and related fields.
- Architecture and implementation detail files under [`tied-project/architecture-decisions/`](../../tied-project/architecture-decisions/) and [`tied-project/implementation-decisions/`](../../tied-project/implementation-decisions/) list related REQ/ARCH/IMPL tokens in `cross_references` / narrative as applicable.

**How to explore**

- Read the YAML indexes: [`tied-project/requirements.yaml`](../../tied-project/requirements.yaml), [`tied-project/architecture-decisions.yaml`](../../tied-project/architecture-decisions.yaml), [`tied-project/implementation-decisions.yaml`](../../tied-project/implementation-decisions.yaml).
- With **TIED MCP**, use tools such as `get_decisions_for_requirement` / `get_requirements_for_decision` (see [`tied-bundle/docs/ai-agent-tied-mcp-usage.md`](../../tied-bundle/docs/ai-agent-tied-mcp-usage.md)).
- Run **`tied_validate_consistency`** after edits to TIED YAML.
