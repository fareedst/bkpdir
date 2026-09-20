#!/usr/bin/env python3
"""Re-baseline logic_level and req_audit_pass from run_impl_logic_audit gates (no mark-complete).

Updates tied/docs/impl-pseudocode-improvement-queue.yaml and impl-pseudocode-sync-checklist.yaml
from live audit results and oracle-registry targets. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
"""
from __future__ import annotations

import argparse
import sys
from datetime import date
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
ORACLE_REG = REPO / "tied/spec/oracle-registry.yaml"

sys.path.insert(0, str(REPO / "scripts"))
from impl_pseudocode_remediation import load_checklist_tokens  # noqa: E402
from load_improvement_queue import load_queue, save_queue  # noqa: E402
from reset_impl_logic_audit_queue import batch_priority, target_logic_level  # noqa: E402
from run_impl_logic_audit import audit_token  # noqa: E402


def verified_oracle_impls() -> set[str]:
    if not ORACLE_REG.exists():
        return set()
    reg = yaml.safe_load(ORACLE_REG.read_text())
    return {o["impl"] for o in reg.get("oracles", []) if o.get("status") == "verified"}


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--dry-run", action="store_true", help="Print summary only")
    args = parser.parse_args()

    in_oracle = verified_oracle_impls()
    today = date.today().isoformat()
    rows = load_checklist_tokens()
    q = load_queue()
    q_by = {t["token"]: t for t in q.get("tokens", [])}

    pass_n = fail_n = 0
    downgraded = 0

    for row in rows:
        token = row["token"]
        tier = row.get("tier", "")
        tgt = target_logic_level(tier, token, token in in_oracle)
        issues = audit_token(row, req_criteria_scoped=True)
        entry = q_by.get(token)
        if not entry:
            continue

        entry["target_logic_level"] = tgt
        entry["formal_spec"] = row.get("formal_spec", True)
        entry["tier"] = tier
        entry["process_order"] = row["process_order"]
        entry["priority"] = batch_priority(row["process_order"])
        entry["go_refs"] = row.get("go_refs", 0)
        entry["oracle_status"] = "verified" if token in in_oracle else entry.get("oracle_status", "n/a")

        prev_logic = entry.get("logic_level", tgt)
        if issues:
            fail_n += 1
            entry["req_audit_pass"] = False
            entry["status"] = "pending_audit"
            entry["issues"] = issues
            entry["logic_level"] = "L1" if "C_doc" in tier else tgt
            entry["audit_notes"] = f"Re-baseline {today}: audit failed — {issues[:3]}"
            row["req_audit_pass"] = False
            row["status"] = "pending_audit"
            row["logic_level"] = entry["logic_level"]
        else:
            pass_n += 1
            entry["req_audit_pass"] = True
            entry["status"] = "logic_verified"
            entry["issues"] = []
            entry["logic_level"] = tgt
            if prev_logic != tgt and prev_logic == "L3" and tgt in ("L1", "L2"):
                downgraded += 1
            entry["audit_notes"] = (
                f"Re-baseline {today}: automated gates pass; logic_level={tgt} (oracle={token in in_oracle})"
            )
            row["req_audit_pass"] = True
            row["status"] = "logic_verified"
            row["logic_level"] = tgt
        row["target_logic_level"] = tgt

    if args.dry_run:
        print(f"DIAGNOSTIC: dry-run pass={pass_n} fail={fail_n} downgraded_from_L3={downgraded}")
        return 1 if fail_n else 0

    save_queue(q)
    checklist = yaml.safe_load(CHECKLIST.read_text())
    checklist["impl_tokens"] = rows
    CHECKLIST.write_text(yaml.dump(checklist, sort_keys=False, allow_unicode=True))

    print(
        f"DIAGNOSTIC: rebaseline complete pass={pass_n} fail={fail_n} "
        f"logic_downgrades={downgraded} date={today}"
    )
    return 1 if fail_n else 0


if __name__ == "__main__":
    raise SystemExit(main())
