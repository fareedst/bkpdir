#!/usr/bin/env python3
"""Report CRIT-001 orphan criteria per IMPL (scoped or strict mode)."""
from __future__ import annotations

import argparse
import csv
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"

sys.path.insert(0, str(REPO / "scripts"))
from impl_pseudocode_remediation import load_checklist_tokens  # noqa: E402
from req_criteria_loader import find_orphan_criteria, orphan_details  # noqa: E402


def filter_batch(tokens: list[dict], batch: str) -> list[dict]:
    if batch == "all":
        return tokens
    ranges = {
        "1": (1, 12),
        "2": (13, 24),
        "3": (25, 36),
        "4": (37, 48),
        "5": (49, 60),
        "6": (61, 73),
    }
    lo, hi = ranges.get(batch, (1, 999))
    return [t for t in tokens if lo <= t.get("process_order", 0) <= hi]


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--token", default="")
    parser.add_argument("--batch", default="all")
    parser.add_argument("--mode", choices=("scoped", "strict"), default="scoped")
    parser.add_argument("--csv", action="store_true")
    args = parser.parse_args()

    rows = load_checklist_tokens()
    if args.token:
        tok = args.token if args.token.startswith("IMPL-") else f"IMPL-{args.token}"
        rows = [r for r in rows if r["token"] == tok]
    else:
        rows = filter_batch(rows, args.batch)

    all_details: list[dict] = []
    pass_n = 0
    for row in rows:
        token = row["token"]
        path = IMPL_DIR / f"{token}-pseudocode.md"
        if not path.exists():
            continue
        text = path.read_text()
        orphans = find_orphan_criteria(text, impl_token=token, mode=args.mode)
        details = orphan_details(text, impl_token=token, mode=args.mode)
        all_details.extend(details)
        if orphans:
            if not args.csv:
                print(f"FAIL {token}: {len(orphans)} orphan(s)")
                for o in orphans:
                    print(f"  {o}")
        else:
            pass_n += 1
            if not args.csv:
                print(f"PASS {token}")

    if args.csv:
        w = csv.DictWriter(
            sys.stdout,
            fieldnames=["impl", "req", "crit_id", "role", "text", "suggested_action"],
        )
        w.writeheader()
        for d in all_details:
            w.writerow(d)
    else:
        print(f"--- {pass_n}/{len(rows)} passed ({args.mode} mode) ---")
    return 0 if pass_n == len(rows) else 1


if __name__ == "__main__":
    raise SystemExit(main())
