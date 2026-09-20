#!/usr/bin/env python3
"""Refresh oracle/logic_level fields on improvement-queue and sync-checklist.

Does NOT set req_audit_pass — use mark_impl_logic_audit_complete.py after manual REQ audit.
"""
from __future__ import annotations

import re
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
ORACLE_REG = REPO / "tied/spec/oracle-registry.yaml"
UPDATE_SCRIPT = REPO / "scripts/update_impl_sync_checklist.py"


def manual_deep_synced() -> set[str]:
    text = UPDATE_SCRIPT.read_text()
    m = re.search(r"MANUAL_DEEP_SYNCED = \{([^}]*)\}", text, re.DOTALL)
    if not m:
        return set()
    return set(re.findall(r'"(IMPL-[^"]+)"', m.group(1)))


def oracle_verified() -> set[str]:
    data = yaml.safe_load(ORACLE_REG.read_text())
    return {o["impl"] for o in data.get("oracles", []) if o.get("status") == "verified"}


def is_tier_c(tier: str) -> bool:
    return "C_doc" in tier


def main() -> int:
    verified = oracle_verified()
    deep = manual_deep_synced()
    checklist = yaml.safe_load(CHECKLIST.read_text())
    go_refs = {t["token"]: t.get("go_refs", 0) for t in checklist["impl_tokens"]}
    tier_by = {t["token"]: t.get("tier", "") for t in checklist["impl_tokens"]}

    q = yaml.safe_load(QUEUE.read_text())
    for t in q["tokens"]:
        token = t["token"]
        tier = tier_by.get(token, t.get("tier", ""))
        refs = go_refs.get(token, 0)
        if is_tier_c(tier) or refs == 0 and token not in verified:
            t.setdefault("target_logic_level", "L1")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L1"
            t["oracle_status"] = "n/a"
            t.setdefault("req_audit_pass", False)
        elif token == "IMPL-SPEC_CTL":
            t.setdefault("target_logic_level", "L3")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L2"
            t["oracle_status"] = "verified"
            t.setdefault("req_audit_pass", False)
        elif token in verified:
            t.setdefault("target_logic_level", "L3")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L3"
            t["oracle_status"] = "verified"
            t.setdefault("req_audit_pass", False)
    QUEUE.write_text(yaml.dump(q, sort_keys=False, allow_unicode=True))

    c = yaml.safe_load(CHECKLIST.read_text())
    for t in c["impl_tokens"]:
        token = t["token"]
        tier = t.get("tier", "")
        refs = t.get("go_refs", 0)
        if is_tier_c(tier) or refs == 0 and token not in verified:
            t.setdefault("target_logic_level", "L1")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L1"
            t["oracle_status"] = "n/a"
            t.setdefault("req_audit_pass", False)
        elif token == "IMPL-SPEC_CTL":
            t.setdefault("target_logic_level", "L3")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L2"
            t["oracle_status"] = "verified"
            t.setdefault("req_audit_pass", False)
        elif token in verified:
            t.setdefault("target_logic_level", "L3")
            if not t.get("req_audit_pass"):
                t["logic_level"] = "L3"
            t["oracle_status"] = "verified"
            t.setdefault("req_audit_pass", False)
    # preserve process_tokens if present
    out = yaml.dump(c, sort_keys=False, allow_unicode=True)
    CHECKLIST.write_text(out)
    print(f"updated queue ({len(q['tokens'])} tokens) and checklist ({len(c['impl_tokens'])} tokens)")
    print(f"oracle verified: {len(verified)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
