#!/usr/bin/env python3
"""Run rewrite + block-lead sync + checklist mark for IMPL tokens in process_order."""
from __future__ import annotations

import subprocess
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
SKIP = {
    "IMPL-LIST_FORMAT_SAFETY",
    "IMPL-CLI_FRAMEWORK",
    "IMPL-CONFIG_STRUCT",
    "IMPL-STRUCTURED_ERRORS",
}


def load_order() -> list[tuple[int, str, str]]:
    with open(CHECKLIST) as f:
        raw = f.read()
    _, _, body = raw.partition("description:")
    data = yaml.safe_load("description:" + body)
    out = []
    for e in data["impl_tokens"]:
        out.append((e.get("process_order", 0), e["token"], e.get("tier", "")))
    return sorted(out)


def run(cmd: list[str]) -> None:
    subprocess.run(cmd, cwd=REPO, check=True)


def main() -> int:
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--from-order", type=int, default=1)
    parser.add_argument("--to-order", type=int, default=72)
    parser.add_argument("--skip-sync", action="store_true")
    args = parser.parse_args()

    for order, token, tier in load_order():
        if order < args.from_order or order > args.to_order:
            continue
        if token in SKIP:
            print(f"skip pilot {token}")
            continue
        print(f"=== {order} {token} ===")
        run([sys.executable, "scripts/rewrite_bulk_sidecar_to_pilot.py", "--token", token])
        if not args.skip_sync:
            run([sys.executable, "scripts/sync_impl_block_leads.py", "--token", token])
        tier_c = tier == "C_doc_process"
        cmd = [sys.executable, "scripts/mark_impl_deep_sync_complete.py", token]
        if tier_c:
            cmd.append("--tier-c")
        run(cmd)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
