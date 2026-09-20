#!/usr/bin/env python3
"""Run automated logic-audit gates for one or all IMPL tokens.

Checks: boilerplate, sidecar structure, REQ token resolution, specctl validate/check-leads,
optional three-way lead sync (via impl_pseudocode_remediation.classify_token).

Exit 0 when all requested tokens pass. Use mark_impl_logic_audit_complete.py after manual REQ review.
"""
from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"
REQ_INDEX = REPO / "tied/requirements.yaml"
REQ_DIR = REPO / "tied/requirements"
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"

BOILERPLATE = (
    "Block implements documented behavior",
    "per resource-management atomic I/O",
)

sys.path.insert(0, str(REPO / "scripts"))
from impl_pseudocode_remediation import classify_token, load_checklist_tokens  # noqa: E402
from load_improvement_queue import load_queue  # noqa: E402
from req_criteria_loader import find_orphan_criteria  # noqa: E402


def load_req_tokens() -> set[str]:
    tokens: set[str] = set()
    if REQ_INDEX.exists():
        data = yaml.safe_load(REQ_INDEX.read_text())
        if isinstance(data, dict):
            for k in data:
                if k.startswith("REQ-"):
                    tokens.add(k)
    if REQ_DIR.exists():
        for f in REQ_DIR.glob("REQ-*.yaml"):
            tokens.add(f.stem)
    meth = REPO / "tied/methodology/requirements.yaml"
    if meth.exists():
        data = yaml.safe_load(meth.read_text())
        if isinstance(data, dict):
            for k in data:
                if k.startswith("REQ-"):
                    tokens.add(k)
    return tokens


def req_tokens_in_sidecar(text: str) -> set[str]:
    return set(re.findall(r"\[REQ-[A-Z0-9_]+\]", text))


def sidecar_issues(
    token: str,
    req_criteria_strict: bool = False,
    req_criteria_scoped: bool = False,
) -> list[str]:
    path = IMPL_DIR / f"{token}-pseudocode.md"
    if not path.exists():
        return ["missing_sidecar"]
    text = path.read_text()
    issues: list[str] = []
    for b in BOILERPLATE:
        if b in text:
            issues.append(f"boilerplate_sidecar:{b[:30]}")
    if not re.search(r"^## Summary contract", text, re.M):
        issues.append("missing_summary_contract")
    leads = [ln for ln in text.splitlines() if ln.startswith("- [IMPL-")]
    if not leads:
        issues.append("missing_block_leads")
    for lead in leads:
        if "How:" not in lead:
            issues.append("lead_missing_how")
            break
        how = lead.split("How:", 1)[-1].strip()
        if len(how) < 12:
            issues.append("generic_how")
            break
        if "Block implements documented behavior" in how:
            issues.append("boilerplate_how")
            break
    # formal DSL: SPEC-ID on runtime blocks
    blocks = re.findall(r"^## ([A-Z_][A-Z0-9_]*)", text, re.M)
    for name in blocks:
        if name == "Summary contract":
            continue
        if f"SPEC-ID: {token}::{name}" not in text and f"SPEC-ID: {token}::" not in text:
            # allow partial — at least one SPEC-ID per token
            pass
    if blocks and "SPEC-ID:" not in text:
        issues.append("missing_spec_id")
    if blocks and not re.search(r"^STEP T\d+:", text, re.M) and not re.search(r"^PROCEDURE ", text, re.M):
        issues.append("missing_step_or_procedure")
    known = load_req_tokens()
    for rt in req_tokens_in_sidecar(text):
        tok = rt.strip("[]")
        if known and tok not in known:
            issues.append(f"unknown_req:{tok}")
    if req_criteria_strict:
        issues.extend(find_orphan_criteria(text, impl_token=token, mode="strict"))
    elif req_criteria_scoped:
        issues.extend(find_orphan_criteria(text, impl_token=token, mode="scoped"))
    return issues


def run_specctl(subcmd: str, sidecar: Path) -> tuple[bool, str]:
    r = subprocess.run(
        ["go", "run", "./cmd/specctl", subcmd, str(sidecar)],
        cwd=REPO,
        capture_output=True,
        text=True,
    )
    out = (r.stdout + r.stderr).strip()
    return r.returncode == 0, out


def audit_token(
    row: dict,
    skip_specctl: bool = False,
    req_criteria_strict: bool = False,
    req_criteria_scoped: bool = False,
) -> list[str]:
    token = row["token"]
    issues = sidecar_issues(
        token,
        req_criteria_strict=req_criteria_strict,
        req_criteria_scoped=req_criteria_scoped,
    )
    pri, align_issues = classify_token(row)
    # P0 alignment issues are blocking for Tier A/B with go_refs
    tier = row.get("tier", "")
    go_refs = row.get("go_refs", 0)
    if tier != "C_doc_process" and go_refs > 0:
        blocking = [i for i in align_issues if i.startswith(("missing_in_go", "extra_in_go", "boilerplate"))]
        issues.extend(blocking)
    sidecar = IMPL_DIR / f"{token}-pseudocode.md"
    if not skip_specctl and sidecar.exists():
        ok, out = run_specctl("validate", sidecar)
        if not ok:
            issues.append("specctl_validate_fail")
        ok, out = run_specctl("check-leads", sidecar)
        if not ok and tier != "C_doc_process" and go_refs > 0:
            issues.append("specctl_check_leads_fail")
    return issues


def filter_batch(tokens: list[dict], batch: str) -> list[dict]:
    if batch == "all":
        return tokens
    if batch == "0":
        return []
    ranges = {"1": (1, 24), "2": (25, 42), "3": (43, 58), "4": (59, 60), "5": (61, 72), "6": (73, 73)}
    lo, hi = ranges.get(batch, (1, 999))
    return [t for t in tokens if lo <= t["process_order"] <= hi]


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--token", default="", help="Single IMPL-TOKEN")
    parser.add_argument("--batch", default="all", help="Batch 1-6 or all")
    parser.add_argument("--skip-specctl", action="store_true")
    parser.add_argument("--mark-complete", action="store_true", help="Auto-mark passing tokens complete")
    parser.add_argument("--tier-c", action="store_true", help="Use tier-c when marking complete")
    parser.add_argument(
        "--req-criteria-strict",
        action="store_true",
        help="L4: fail on all REQ criteria orphans (CRIT-001 strict)",
    )
    parser.add_argument(
        "--req-criteria-scoped",
        action="store_true",
        help="L4: fail on scoped REQ criteria orphans (CRIT-001; uses scope registry)",
    )
    parser.add_argument(
        "--mark-crit001",
        action="store_true",
        help="Auto-mark passing tokens crit001_pass via mark_crit001_complete.py",
    )
    args = parser.parse_args()

    rows = load_checklist_tokens()
    if args.token:
        tok = args.token if args.token.startswith("IMPL-") else f"IMPL-{args.token}"
        rows = [r for r in rows if r["token"] == tok]
    else:
        rows = filter_batch(rows, args.batch)

    fail = 0
    for row in rows:
        token = row["token"]
        issues = audit_token(
            row,
            skip_specctl=args.skip_specctl,
            req_criteria_strict=args.req_criteria_strict,
            req_criteria_scoped=args.req_criteria_scoped,
        )
        if issues:
            fail += 1
            print(f"FAIL {token}: {issues}")
        else:
            print(f"PASS {token}")
            if args.mark_crit001:
                mode = "strict" if args.req_criteria_strict else "scoped"
                subprocess.run(
                    [
                        sys.executable,
                        str(REPO / "scripts/mark_crit001_complete.py"),
                        token,
                        "--mode",
                        mode,
                        "--skip-gates",
                    ],
                    cwd=REPO,
                    check=True,
                )
            if args.mark_complete:
                tier_c = args.tier_c or row.get("tier") == "C_doc_process"
                cmd = [
                    sys.executable,
                    str(REPO / "scripts/mark_impl_logic_audit_complete.py"),
                    token,
                    "--skip-gates",
                ]
                if tier_c:
                    cmd.append("--tier-c")
                q = load_queue()
                for qt in q.get("tokens", []):
                    if qt["token"] == token:
                        cmd.extend(["--logic-level", qt.get("target_logic_level", "L3")])
                        break
                subprocess.run(cmd, cwd=REPO, check=True)

    print(f"--- {len(rows) - fail}/{len(rows)} passed ---")
    return 1 if fail else 0


if __name__ == "__main__":
    raise SystemExit(main())
