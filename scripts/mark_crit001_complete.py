#!/usr/bin/env python3
"""Mark one IMPL token CRIT-001 scoped audit complete (crit001_pass)."""
from __future__ import annotations

import argparse
import subprocess
import sys
from datetime import date
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"

sys.path.insert(0, str(REPO / "scripts"))
from load_improvement_queue import load_queue, save_queue  # noqa: E402
from req_criteria_loader import find_orphan_criteria  # noqa: E402


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("token", help="IMPL-TOKEN")
    parser.add_argument("--mode", choices=("scoped", "strict"), default="scoped")
    parser.add_argument("--scope-notes", default="", help="CRIT-001 scope audit summary")
    parser.add_argument("--skip-gates", action="store_true")
    args = parser.parse_args()
    token = args.token if args.token.startswith("IMPL-") else f"IMPL-{args.token}"

    sidecar = IMPL_DIR / f"{token}-pseudocode.md"
    if not sidecar.exists():
        print(f"Missing sidecar: {sidecar}", file=sys.stderr)
        return 1

    text = sidecar.read_text()
    orphans = find_orphan_criteria(text, impl_token=token, mode=args.mode)
    if orphans and not args.skip_gates:
        for o in orphans:
            print(o, file=sys.stderr)
        print(f"CRIT-001: {len(orphans)} orphan(s) in {args.mode} mode", file=sys.stderr)
        return 1

    if not args.skip_gates:
        r = subprocess.run(
            ["go", "run", "./cmd/specctl", "validate", str(sidecar)],
            cwd=REPO,
            capture_output=True,
            text=True,
        )
        if r.returncode != 0:
            print(r.stdout + r.stderr, file=sys.stderr)
            return 1

    today = date.today().isoformat()
    q = load_queue()
    for entry in q.get("tokens", []):
        if entry["token"] != token:
            continue
        entry["crit001_pass"] = True
        entry["crit001_orphans_remaining"] = 0
        entry["crit001_scope_notes"] = (
            args.scope_notes or f"CRIT-001 {args.mode} audit complete {today}"
        )
        break
    else:
        print(f"Token not in improvement queue: {token}", file=sys.stderr)
        return 1

    save_queue(q)
    print(f"Marked {token} crit001_pass ({args.mode} mode)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
