#!/usr/bin/env python3
"""Reset non-pilot IMPL checklist rows for manual deep-sync queue."""
from __future__ import annotations

from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"

PILOT_COMPLETE = {
    "IMPL-LIST_FORMAT_SAFETY",
    "IMPL-CLI_FRAMEWORK",
    "IMPL-CONFIG_STRUCT",
    "IMPL-STRUCTURED_ERRORS",
}

HEADER = """# IMPL Pseudocode Sync Checklist (canonical template)
#
# Process tokens: [PROC-IMPL_CODE_TEST_SYNC], [PROC-PSEUDOCODE_VALIDATION], [PROC-IMPL_PSEUDOCODE_TOKENS]
# Workflow: Track C (brownfield) — tests/code → sidecar → literal block-lead sync
# Reference: tied/docs/pseudocode-writing-and-validation.md § Track C
#
# Per-token status:
#   pending | in_progress | sidecar_migrated | manual_deep_sync | comments_synced | validated
#   validated = pilot bar only (language-agnostic sidecar + literal // - block leads in code/tests)
# Process one IMPL token at a time; update steps and status after each pass.
#
# Repeatable 12-step process (summary):
#   A1-A3 Discovery | B4-B6 Read tests/code, list blocks | C7-C10 Sidecar migrate
#   D-F11-15 Literal block-lead sync | G16-17 Metadata | H18-20 Validate + go test
"""


def main() -> None:
    with open(CHECKLIST) as f:
        raw = f.read()
    _, _, body = raw.partition("description:")
    data = yaml.safe_load("description:" + body)

    reset = 0
    for entry in data["impl_tokens"]:
        token = entry["token"]
        if token in PILOT_COMPLETE:
            continue
        entry["status"] = "manual_deep_sync"
        entry["notes"] = ""
        steps = entry.setdefault("steps", {})
        steps["layer_b_validated"] = False
        steps["blocks_synced"] = 0
        steps["discovery"] = True
        steps["sidecar_created"] = True
        steps["metadata_updated"] = False
        steps["tests_green"] = False
        steps["layer_a_validated"] = False
        reset += 1

    with open(CHECKLIST, "w") as f:
        f.write(HEADER)
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120)
    print(f"Reset {reset} tokens to manual_deep_sync; preserved {len(PILOT_COMPLETE)} pilot rows")


if __name__ == "__main__":
    main()
