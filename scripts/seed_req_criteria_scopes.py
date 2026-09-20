#!/usr/bin/env python3
"""Seed req-criteria-scope-registry.yaml from REQ/IMPL traceability + sidecar REQ cites."""
from __future__ import annotations

import re
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
SCOPE_REG = REPO / "tied/spec/req-criteria-scope-registry.yaml"
IMPL_DIR = REPO / "tied/implementation-decisions"
REQ_DIRS = (
    REPO / "tied/requirements",
    REPO / "tied/methodology/requirements",
)
IMPL_DECISION_DIRS = (REPO / "tied/implementation-decisions",)

sys_path = REPO / "scripts"
import sys

sys.path.insert(0, str(sys_path))
from req_criteria_loader import load_all_criteria  # noqa: E402


def _detail_record(data: dict, token: str) -> dict | None:
    if not isinstance(data, dict):
        return None
    if token in data and isinstance(data[token], dict):
        return data[token]
    return data if isinstance(data, dict) else None


def req_primary_impls(req_token: str) -> set[str]:
    owners: set[str] = set()
    for req_dir in REQ_DIRS:
        path = req_dir / f"{req_token}.yaml"
        if not path.exists():
            continue
        data = yaml.safe_load(path.read_text())
        rec = _detail_record(data, req_token)
        if not rec:
            continue
        impl_field = rec.get("implementation") or ""
        if isinstance(impl_field, str):
            for m in re.findall(r"IMPL-[A-Z0-9_]+", impl_field):
                owners.add(m)
        trace = rec.get("traceability") or {}
        for impl in trace.get("implementation") or []:
            if isinstance(impl, str) and impl.startswith("IMPL-"):
                owners.add(impl)
    return owners


def impl_declared_reqs(impl_token: str) -> set[str]:
    reqs: set[str] = set()
    for d in IMPL_DECISION_DIRS:
        path = d / f"{impl_token}.yaml"
        if not path.exists():
            continue
        data = yaml.safe_load(path.read_text())
        rec = _detail_record(data, impl_token)
        if not rec:
            continue
        trace = rec.get("traceability") or {}
        for r in trace.get("requirements") or []:
            if isinstance(r, str) and r.startswith("REQ-"):
                reqs.add(r)
    return reqs


def sidecar_reqs(impl_token: str) -> set[str]:
    path = IMPL_DIR / f"{impl_token}-pseudocode.md"
    if not path.exists():
        return set()
    text = path.read_text()
    return {m.strip("[]") for m in re.findall(r"\[REQ-[A-Z0-9_]+\]", text)}


def all_formal_impls() -> list[str]:
    tokens: list[str] = []
    for path in sorted(IMPL_DIR.glob("IMPL-*-pseudocode.md")):
        tokens.append(path.stem.replace("-pseudocode", ""))
    return tokens


def classify_role(impl: str, req: str) -> str:
    owners = req_primary_impls(req)
    if impl in owners:
        return "primary"
    declared = impl_declared_reqs(impl)
    if req in declared and not owners:
        return "primary"
    if req in declared and impl in owners:
        return "primary"
    return "auxiliary"


def crit_ids_for_req(req: str, all_crits: dict) -> list[str]:
    return [c.crit_id for c in all_crits.get(req, [])]


def build_scopes(preserve_manual: bool = True) -> list[dict]:
    all_crits = load_all_criteria()
    manual: dict[tuple[str, str], dict] = {}
    if preserve_manual and SCOPE_REG.exists():
        existing = yaml.safe_load(SCOPE_REG.read_text()) or {}
        for row in existing.get("scopes") or []:
            key = (row.get("impl", ""), row.get("req", ""))
            if row.get("manual"):
                manual[key] = row

    scopes: list[dict] = []
    seen: set[tuple[str, str]] = set()
    for impl in all_formal_impls():
        reqs = sidecar_reqs(impl) | impl_declared_reqs(impl)
        for req in sorted(reqs):
            key = (impl, req)
            if key in seen:
                continue
            seen.add(key)
            if key in manual:
                scopes.append(manual[key])
                continue
            role = classify_role(impl, req)
            if role == "primary":
                applicable = crit_ids_for_req(req, all_crits)
            else:
                applicable = []
            scopes.append(
                {
                    "impl": impl,
                    "req": req,
                    "role": role,
                    "applicable_crit_ids": applicable,
                    "rationale": (
                        f"Primary owner per REQ traceability.implementation"
                        if role == "primary"
                        else "Auxiliary cross-ref; traceability-only unless crit ids listed"
                    ),
                }
            )
    scopes.sort(key=lambda r: (r["impl"], r["req"]))
    return scopes


def main() -> int:
    scopes = build_scopes()
    data = {
        "description": "Scoped REQ satisfaction criteria for formal_spec sidecars (L4 CRIT-001 pilot).",
        "scopes": scopes,
    }
    header = """# REQ criteria scope per (IMPL, REQ) pair for CRIT-001 scoped audit.
# [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
# Seed: scripts/seed_req_criteria_scopes.py
# Runbook: tied/docs/crit001-corpus-alignment-runbook.md
"""
    SCOPE_REG.write_text(header + yaml.dump(data, default_flow_style=False, sort_keys=False, allow_unicode=True))
    primary_n = sum(1 for s in scopes if s["role"] == "primary")
    aux_n = sum(1 for s in scopes if s["role"] == "auxiliary")
    print(f"DEBUG: seeded {len(scopes)} scope rows ({primary_n} primary, {aux_n} auxiliary)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
