# Governance and traceability hub

Normative **requirements, architecture, and implementation** are maintained in **TIED** under [`tied-project/`](../tied-project/). Do not treat this page as a second registry; use it for navigation and visuals only.

## Where to look

| Need | Location |
|------|----------|
| REQ index + detail YAML | [`tied-project/requirements.yaml`](../tied-project/requirements.yaml), [`tied-project/requirements/`](../tied-project/requirements/) |
| ARCH index + detail YAML | [`tied-project/architecture-decisions.yaml`](../tied-project/architecture-decisions.yaml), [`tied-project/architecture-decisions/`](../tied-project/architecture-decisions/) |
| IMPL index + detail YAML (incl. pseudo-code) | [`tied-project/implementation-decisions.yaml`](../tied-project/implementation-decisions.yaml), [`tied-project/implementation-decisions/`](../tied-project/implementation-decisions/) |
| Token registry | [`tied-project/semantic-tokens.yaml`](../tied-project/semantic-tokens.yaml), [`tied-bundle/docs/semantic-tokens.md`](../tied-bundle/docs/semantic-tokens.md) |
| Agent process | [`AGENTS.md`](../AGENTS.md), [`ai-principles.md`](../ai-principles.md) |
| Governance narrative | [`governance/index.md`](governance/index.md), [`governance/workflow.md`](governance/workflow.md) |
| Optional tasks | [`tied-bundle/docs/tasks.md`](../tied-bundle/docs/tasks.md) |

Project governance-related requirements that exist in TIED today include **[REQ-GOV_REGISTRY_COMPLETENESS]** and **[REQ-GOV_DISCOVERABILITY]** (see their detail files under `tied-project/requirements/`). Additional REQ tokens are listed only in `tied-project/requirements.yaml` and detail YAML—cross-links to ARCH/IMPL live in each file’s `traceability` section.

## Token-first visuals (`[REQ-STDD_VIS]`)

- Layered flow: [`docs/images/stdd-flow.svg`](images/stdd-flow.svg) — `[REQ-CFG_005]`, `[REQ-LIST_LIMIT]` chains ([ARCH-STDD_VIS_FLOW], [IMPL-STDD_VIS_DATA_PIPELINE]).
- Timeline storyboard: [`docs/images/stdd-timeline.svg`](images/stdd-timeline.svg), data in [`docs/images/stdd-timeline.json`](images/stdd-timeline.json), [`docs/images/stdd-timeline.md`](images/stdd-timeline.md).
- Token stats: [`docs/images/stdd-token-stats.svg`](images/stdd-token-stats.svg) (data: [`docs/data/stdd-token-stats.json`](data/stdd-token-stats.json)).
- Config hierarchy: [`docs/images/stdd-config-flow.svg`](images/stdd-config-flow.svg) (data: [`docs/data/stdd-config-trace.json`](data/stdd-config-trace.json)).
- Config heatmap: [`docs/images/stdd-config-heatmap.svg`](images/stdd-config-heatmap.svg) (data: [`docs/data/stdd-config-token-stats.json`](data/stdd-config-token-stats.json)).
