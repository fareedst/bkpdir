# STDD Visualization Definitions & Upgrade Plan

> Tokens: [REQ-STDD_VIS], [ARCH-STDD_VIS_FLOW], [IMPL-STDD_VIS_DATA_PIPELINE], [IMPL-STDD_VIS_ASSETS]

## Layered Flow (docs/images/stdd-flow.svg)
- **Definition**: Static SVG generated from `stdd-flow.mmd`/`stdd-flow.json`; shows `[REQ-CFG_005]` and `[REQ-LIST_LIMIT]` chains traversing REQ → ARCH → IMPL → TEST → CODE lanes.
- **Data Source**: `docs/data/stdd-trace.json` (nodes/edges/styling hints) ensures reproducibility.
- **Upgrade Plan**: When adding more chains, append nodes/edges to `stdd-trace.json` and regenerate the Mermaid/SVG; keep color palette consistent for token recognition.

## Token Timeline (docs/images/stdd-timeline.svg)
- **Definition**: Static multi-step SVG summarizing the storyboard in `stdd-timeline.json`/`stdd-timeline.md`; serves as the animation substitute.
- **Upgrade Plan**:
  1. Use `stdd-timeline.json` to drive frame generation (presentation software or scripted renderer).
  2. Export MP4/GIF/animated SVG; save as `docs/images/stdd-timeline.mp4` (or similar).
  3. Update docs to embed both the static SVG (fallback) and the animated asset; rerun validation log entry in `docs/data/stdd-validation-log.md`.

## Token Statistics Visual (docs/images/stdd-token-stats.svg)
- **Definition**: Bar-chart SVG sourced from `docs/data/stdd-token-stats.json`; shows counts for `[REQ:*]`, `[ARCH:*]`, `[IMPL:*]`, `[TEST:*]` tokens as of 2025-12-31.
- **Upgrade Plan**: Refresh the JSON counts whenever `tied/semantic-tokens.yaml` / `tied/semantic-tokens.md` changes significantly, then regenerate the SVG to keep visuals current.

## Configuration Token Visuals
- **Config Hierarchy Map (SVG)**:
  - **Status**: ✅ Implemented (`docs/images/stdd-config-flow.svg`) using `docs/data/stdd-config-trace.json`.
  - **Description**: Layered diagram showing `[REQ-CONFIGURATION]`, `[REQ-CFG_005]`, `[REQ-CFG_006]`, `[ARCH-CONFIG_SYSTEM]`, `[ARCH-CFG_005]`, `[ARCH-CFG_006]`, `[IMPL-CONFIG_STRUCT]`, `[IMPL-EXCLUDE_MERGE_FIX]`, `[IMPL-CFG_006]`, `[IMPL-CONFIG_OUTPUT_GROUPING]`, and their test/code anchors.

- **Configuration State Timeline (SVG)**:
  - **Status**: 🔜 Planned.
  - **Goal**: Depict how defaults propagate and are overridden via layered configs, referencing `[REQ-CFG_005]`, `[REQ-CFG_006]`, `[IMPL-CFG_INHERITANCE_PATH_RESOLUTION]`, `[IMPL-CONFIG_OUTPUT_GROUPING]`.
  - **Format**: Timeline storyboard/animation similar to `stdd-timeline`.

- **Configuration Defaults Heatmap (SVG)**:
  - **Status**: ✅ Implemented (`docs/images/stdd-config-heatmap.svg`) using `docs/data/stdd-config-token-stats.json`.
  - **Description**: Matrix showing counts of configuration tokens per category (Defaults, Inheritance, Reflection, CLI Output) across REQ/ARCH/IMPL/TEST.

## Integration Notes
- `docs/index.md` references both SVGs so stakeholders can view the visuals directly.
- Validation results recorded in `docs/data/stdd-validation-log.md` for traceability.

