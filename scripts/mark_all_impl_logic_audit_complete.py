#!/usr/bin/env python3
"""Mark all IMPL tokens logic-audit complete after run_impl_logic_audit.py passes."""
from __future__ import annotations

import subprocess
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(REPO / "scripts"))
from load_improvement_queue import load_queue  # noqa: E402


def main() -> int:
    q = load_queue()
    fail = 0
    for t in sorted(q["tokens"], key=lambda x: x["process_order"]):
        if t.get("req_audit_pass"):
            continue
        token = t["token"]
        cmd = [
            sys.executable,
            str(REPO / "scripts/mark_impl_logic_audit_complete.py"),
            token,
            "--logic-level",
            t.get("target_logic_level", "L3"),
        ]
        if "C_doc" in t.get("tier", ""):
            cmd.append("--tier-c")
        r = subprocess.run(cmd, cwd=REPO)
        if r.returncode != 0:
            fail += 1
            print(f"FAILED mark: {token}", file=sys.stderr)
    print(f"Marked {len(q['tokens']) - fail}/{len(q['tokens'])} tokens")
    return 1 if fail else 0


if __name__ == "__main__":
    raise SystemExit(main())
