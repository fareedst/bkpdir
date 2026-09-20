#!/usr/bin/env python3
"""Batch 0: reset logic-audit tracking for IMPL logic alignment program.

Re-prioritizes improvement queue (P0–P2 by process_order batch), sets
req_audit_pass=false and status=pending_audit on all tokens unless --preserve-audited.
"""
from __future__ import annotations

import argparse
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"

QUEUE_HEADER = """# IMPL Logic Alignment improvement queue
# Batch priorities: P0 orders 1-42, P1 orders 43-60, P2 orders 61-73
# Reset: scripts/reset_impl_logic_audit_queue.py
# Complete: scripts/mark_impl_logic_audit_complete.py IMPL-TOKEN [--logic-level L3]
# REQ audit runbook: tied/docs/impl-deep-sync-agent-guide.md §14
# Do NOT use seed_logic_verification_tracking.py to set req_audit_pass.
"""


def batch_priority(process_order: int) -> str:
    if process_order <= 42:
        return "P0"
    if process_order <= 60:
        return "P1"
    return "P2"


def target_logic_level(tier: str, token: str, in_oracle: bool) -> str:
    if token == "IMPL-SPEC_CTL":
        return "L3"
    if "C_doc" in tier:
        return "L1"
    if in_oracle:
        return "L3"
    if tier.startswith("B_test"):
        return "L2"
    if tier.startswith("B_"):
        return "L2"
    return "L3"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--preserve-audited",
        action="store_true",
        help="Keep req_audit_pass=true rows that already passed manual audit",
    )
    args = parser.parse_args()

    checklist = yaml.safe_load(CHECKLIST.read_text())
    tier_by = {t["token"]: t.get("tier", "") for t in checklist["impl_tokens"]}
    order_by = {t["token"]: t["process_order"] for t in checklist["impl_tokens"]}
    go_refs = {t["token"]: t.get("go_refs", 0) for t in checklist["impl_tokens"]}
    formal = {t["token"]: t.get("formal_spec", False) for t in checklist["impl_tokens"]}

    oracle_reg = REPO / "tied/spec/oracle-registry.yaml"
    in_oracle: set[str] = set()
    if oracle_reg.exists():
        reg = yaml.safe_load(oracle_reg.read_text())
        in_oracle = {o["impl"] for o in reg.get("oracles", []) if o.get("status") == "verified"}

    old_queue = yaml.safe_load(QUEUE.read_text()) if QUEUE.exists() else {"tokens": []}
    old_by_token = {t["token"]: t for t in old_queue.get("tokens", [])}

    tokens = []
    for row in sorted(checklist["impl_tokens"], key=lambda x: x["process_order"]):
        token = row["token"]
        order = row["process_order"]
        tier = row.get("tier", tier_by.get(token, ""))
        prev = old_by_token.get(token, {})
        preserve = args.preserve_audited and prev.get("req_audit_pass") is True
        tgt = target_logic_level(tier, token, token in in_oracle)

        entry = {
            "token": token,
            "process_order": order,
            "tier": tier,
            "priority": batch_priority(order),
            "formal_spec": formal.get(token, prev.get("formal_spec", True)),
            "target_logic_level": tgt,
            "logic_level": prev.get("logic_level", tgt) if preserve else row.get("logic_level", tgt),
            "oracle_status": "verified" if token in in_oracle else ("verified" if token == "IMPL-SPEC_CTL" else "n/a"),
            "deep_sync_complete": row.get("deep_sync_complete", True),
            "req_audit_pass": True if preserve else False,
            "status": "logic_verified" if preserve else "pending_audit",
            "issues": [] if preserve else prev.get("issues", []),
            "audit_notes": prev.get("audit_notes", "") if preserve else "",
            "go_refs": go_refs.get(token, 0),
            "notes": row.get("notes", prev.get("notes", "")),
        }
        tokens.append(entry)

    body = yaml.safe_dump(
        {"description": "Per-token logic alignment tracking (REQ audit + three-way sync).", "tokens": tokens},
        sort_keys=False,
        allow_unicode=True,
    )
    QUEUE.write_text(QUEUE_HEADER + body)

    for row in checklist["impl_tokens"]:
        token = row["token"]
        prev = old_by_token.get(token, {})
        preserve = args.preserve_audited and prev.get("req_audit_pass") is True
        if not preserve:
            row["req_audit_pass"] = False
            if row.get("status") == "logic_verified":
                row["status"] = "pending_audit"
        row["target_logic_level"] = target_logic_level(
            row.get("tier", ""), token, token in in_oracle
        )
    CHECKLIST.write_text(yaml.dump(checklist, sort_keys=False, allow_unicode=True))

    pending = sum(1 for t in tokens if not t["req_audit_pass"])
    print(f"Reset queue: {len(tokens)} tokens, {pending} pending_audit, priorities P0/P1/P2 by batch")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
