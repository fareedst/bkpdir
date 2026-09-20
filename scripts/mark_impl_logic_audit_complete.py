#!/usr/bin/env python3
"""Mark one IMPL token logic-audit complete (req_audit_pass + logic_level).

Updates improvement queue, sync checklist, and MANUAL_DEEP_SYNCED guard.
Only call after per-token exit checklist in impl-deep-sync-agent-guide.md §14.
"""
from __future__ import annotations

import argparse
import subprocess
import sys
from datetime import date
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
IMPL_DIR = REPO / "tied/implementation-decisions"

sys.path.insert(0, str(REPO / "scripts"))
from load_improvement_queue import load_queue, save_queue  # noqa: E402
from mark_impl_deep_sync_complete import add_to_manual_guard, count_block_leads  # noqa: E402
from req_criteria_loader import find_orphan_criteria  # noqa: E402


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("token", help="IMPL-TOKEN")
    parser.add_argument("--logic-level", default="", help="L1, L2, L3, or L4 (default: target_logic_level from queue)")
    parser.add_argument("--tier-c", action="store_true", help="Tier C note (no Go lead sync required)")
    parser.add_argument("--skip-gates", action="store_true", help="Skip specctl gate (not for production use)")
    parser.add_argument("--audit-notes", default="", help="Short REQ audit summary")
    parser.add_argument(
        "--req-criteria-strict",
        action="store_true",
        help="L4 pilot: block mark if orphan REQ criteria (CRIT-001)",
    )
    args = parser.parse_args()
    token = args.token if args.token.startswith("IMPL-") else f"IMPL-{args.token}"

    sidecar = IMPL_DIR / f"{token}-pseudocode.md"
    if not sidecar.exists():
        print(f"Missing sidecar: {sidecar}", file=sys.stderr)
        return 1

    if args.req_criteria_strict:
        orphans = find_orphan_criteria(sidecar.read_text(), impl_token=token, mode="strict")
        if orphans:
            for o in orphans:
                print(o, file=sys.stderr)
            print(f"CRIT-001: {len(orphans)} orphan criteria — fix or LEAP before mark complete", file=sys.stderr)
            return 1

    if not args.skip_gates:
        for cmd in (
            ["go", "run", "./cmd/specctl", "validate", str(sidecar)],
            ["go", "run", "./cmd/specctl", "check-leads", str(sidecar)],
        ):
            r = subprocess.run(cmd, cwd=REPO, capture_output=True, text=True)
            if r.returncode != 0:
                print(r.stdout + r.stderr, file=sys.stderr)
                print(f"Gate failed: {' '.join(cmd)}", file=sys.stderr)
                return 1

    blocks = count_block_leads(sidecar.read_text())
    today = date.today().isoformat()

    q = load_queue()
    q_entry = None
    for entry in q["tokens"]:
        if entry["token"] == token:
            q_entry = entry
            logic = args.logic_level or entry.get("target_logic_level", "L3")
            entry["req_audit_pass"] = True
            entry["status"] = "logic_verified"
            entry["logic_level"] = logic
            entry["issues"] = []
            note = args.audit_notes or f"Logic audit complete {today}; REQ criteria mapped to STEP/PROCEDURE"
            entry["audit_notes"] = note
            if logic == "L4":
                entry["mutation_verified"] = True
            entry.setdefault("mutation_verified", False)
            entry.setdefault("mutation_score", None)
            if args.tier_c:
                entry["notes"] = "Manual logic audit: Tier C sidecar + REQ traceability"
            else:
                entry["notes"] = f"Manual logic audit: three-way sync + REQ audit ({logic})"
            entry["deep_sync_complete"] = True
            break
    else:
        print(f"Token not in improvement queue: {token}", file=sys.stderr)
        return 1

    save_queue(q)

    raw = CHECKLIST.read_text()
    header, _, body = raw.partition("description:")
    data = yaml.safe_load("description:" + body)
    for entry in data["impl_tokens"]:
        if entry["token"] != token:
            continue
        entry["req_audit_pass"] = True
        entry["status"] = "logic_verified"
        entry["logic_level"] = q_entry["logic_level"]
        entry["deep_sync_complete"] = True
        steps = entry.setdefault("steps", {})
        steps["blocks_total"] = blocks
        steps["blocks_synced"] = blocks
        steps["layer_a_validated"] = True
        steps["layer_b_validated"] = True
        steps["tests_green"] = True
        steps["req_audit_pass"] = True
        break
    with open(CHECKLIST, "w") as f:
        if header:
            f.write(header)
        yaml.dump(data, f, default_flow_style=False, sort_keys=False, allow_unicode=True, width=120)

    add_to_manual_guard(token)
    print(f"Marked {token} logic audit complete ({blocks} blocks, {q_entry['logic_level']})")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
