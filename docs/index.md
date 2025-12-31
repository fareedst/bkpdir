# Governance Discovery Hub

Overview: Central discovery point for governance tokens and their connected decisions. This hub links to the REQ token registry and to related ARCH/IMPL decisions.

## Tokens (REQ-focused)
- [REQ:GOVERNANCE_DOCUMENTATION] Governance rules must be documented with explicit anchors and cross-links to ARCH/IMPL decisions.
- [REQ:GOV_POLICY_MIGRATION] Migration of governance policies into STDD decision documents.
- [REQ:GOV_CHANGE_CONTROL] Requirements for proposing, approving, and rolling out governance changes.
- [REQ:GOV_AUDIT] Regular auditing of governance rules; ensure alignment with decisions and tests.
- [REQ:GOV_REGISTRY_COMPLETENESS] Registry must capture all governance-related tokens and their intents.
- [REQ:GOV_DISCOVERABILITY] Centralized discoverability of governance content via docs/index.md with links to tokens and decisions.
- [REQ:GOV_NAMING_CONVENTIONS] Establish and enforce token naming conventions for governance-related tokens.
- [REQ:GOV_CROSS_LINKING] Each REQ token must be cross-linked to relevant ARCH/IMPL decisions.
- [REQ:GOV_MIGRATION_STATUS] Track the progress/status of governance-content migration to ARCH/IMPL docs.
- [REQ:REGISTRY_COMPLETION_VERIFICATION] Verification rule that ensures registry remains complete over time.

## Cross-links (ARCH/IMPL)
- See cross-link mappings in the governance page and the architecture/implementation decision docs for each REQ token.

## Quick Links
- `docs/governance/index.md` (governance page)
- `docs/architecture-decisions.md`
- `docs/implementation-decisions.md`
- `stdd/tasks.md`

## Discoverability
- This hub is the primary entry point for governance token discovery. Use it to navigate to concrete decisions and migration efforts.

## STDD Visualization Gallery
- `[REQ:STDD_VIS]` token-first visuals demonstrate how requirements, architecture, implementation, tests, and code stay linked via semantic tokens.
- Layered flow diagram (static SVG): `docs/images/stdd-flow.svg` — highlights `[REQ:CFG_005]` and `[REQ:LIST_LIMIT]` chains across REQ/ARCH/IMPL/TEST/CODE lanes ([ARCH:STDD_VIS_FLOW], [IMPL:STDD_VIS_DATA_PIPELINE]).
- Token timeline storyboard (static SVG substitute for animation): `docs/images/stdd-timeline.svg` — shows sequential hand-offs for both chains, referencing `[IMPL:STDD_VIS_ASSETS]`.
- Storyboard + data sources live in `docs/images/stdd-timeline.json` and `docs/images/stdd-timeline.md` for future upgrades to animated media.
- Token statistics snapshot: `docs/images/stdd-token-stats.svg` (data in `docs/data/stdd-token-stats.json`) visualizes aggregate counts for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`, `[TEST:*]`.
- Configuration hierarchy map: `docs/images/stdd-config-flow.svg` (data in `docs/data/stdd-config-trace.json`) details how configuration tokens (`[REQ:CONFIGURATION]`, `[REQ:CFG_005]`, `[REQ:CFG_006]`, etc.) propagate through architecture, implementation, tests, and code.
- Configuration heatmap: `docs/images/stdd-config-heatmap.svg` (data in `docs/data/stdd-config-token-stats.json`) highlights token distribution across Defaults, Inheritance, Reflection, and CLI Output categories.
