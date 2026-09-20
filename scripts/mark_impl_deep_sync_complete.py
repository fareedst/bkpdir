#!/usr/bin/env python3
"""Mark one IMPL token as pilot-validated in checklist and MANUAL_DEEP_SYNCED guard."""
from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
IMPL_DIR = REPO / "tied/implementation-decisions"
UPDATE_SCRIPT = REPO / "scripts/update_impl_sync_checklist.py"


def count_block_leads(text: str) -> int:
    return len(re.findall(r"^- \[IMPL-", text, re.MULTILINE))


def add_to_manual_guard(token: str) -> None:
    text = UPDATE_SCRIPT.read_text()
    m = re.search(r"MANUAL_DEEP_SYNCED = \{([^}]*)\}", text, re.DOTALL)
    if not m:
        raise RuntimeError("MANUAL_DEEP_SYNCED not found")
    block = m.group(1)
    if f'"{token}"' in block:
        return
    tokens = sorted(set(re.findall(r'"IMPL-[^"]+"', block)) | {f'"{token}"'})
    new_block = ",\n    ".join(tokens) + ",\n"
    text = text[: m.start(1)] + new_block + text[m.end(1) :]
    UPDATE_SCRIPT.write_text(text)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("token", help="IMPL-TOKEN")
    parser.add_argument("--tier-c", action="store_true", help="Sidecar-only Tier C note")
    args = parser.parse_args()
    token = args.token if args.token.startswith("IMPL-") else f"IMPL-{args.token}"

    sidecar = IMPL_DIR / f"{token}-pseudocode.md"
    if not sidecar.exists():
        print(f"Missing sidecar: {sidecar}", file=sys.stderr)
        return 1
    blocks = count_block_leads(sidecar.read_text())

    with open(CHECKLIST) as f:
        raw = f.read()
    header, _, body = raw.partition("description:")
    data = yaml.safe_load("description:" + body)

    for entry in data["impl_tokens"]:
        if entry["token"] != token:
            continue
        entry["status"] = "validated"
        if args.tier_c:
            entry["notes"] = "Manual deep sync: Tier C sidecar quality (no code comment sync)"
        else:
            entry["notes"] = "Manual deep sync: pilot-quality three-way block leads"
        steps = entry.setdefault("steps", {})
        steps["discovery"] = True
        steps["sidecar_created"] = True
        steps["blocks_total"] = blocks
        steps["blocks_synced"] = blocks
        steps["metadata_updated"] = True
        steps["layer_a_validated"] = True
        steps["layer_b_validated"] = True
        steps["tests_green"] = True
        break
    else:
        print(f"Token not in checklist: {token}", file=sys.stderr)
        return 1

    with open(CHECKLIST, "w") as f:
        f.write(header)
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120)

    add_to_manual_guard(token)
    print(f"Marked {token} validated ({blocks} blocks)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
