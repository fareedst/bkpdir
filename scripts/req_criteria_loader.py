#!/usr/bin/env python3
"""Load REQ satisfaction/validation criteria from TIED detail YAML (project + methodology).

Used by run_impl_logic_audit.py for L4 CRIT-001 orphan detection.
"""
from __future__ import annotations

import re
from dataclasses import dataclass
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
REQ_DIRS = (
    REPO / "tied/requirements",
    REPO / "tied/methodology/requirements",
)
SCOPE_REG = REPO / "tied/spec/req-criteria-scope-registry.yaml"
IMPL_DECISION_DIR = REPO / "tied/implementation-decisions"

STOP_WORDS = frozenset(
    """
    a an the and or not for with from into that this each every all must may
    when then than using use uses used per via are be is was were has have had
    will can should could would being been also only other such their them they
    method coverage automated testing validation criterion criteria
    """.split()
)


@dataclass(frozen=True)
class Criterion:
    req_token: str
    crit_id: str
    text: str
    kind: str  # satisfaction | validation


def _detail_record(data: dict, token: str) -> dict | None:
    if not isinstance(data, dict):
        return None
    if token in data and isinstance(data[token], dict):
        return data[token]
    return data if "satisfaction_criteria" in data or "validation_criteria" in data else None


def load_req_detail(req_token: str) -> dict | None:
    for req_dir in REQ_DIRS:
        path = req_dir / f"{req_token}.yaml"
        if not path.exists():
            continue
        data = yaml.safe_load(path.read_text())
        return _detail_record(data, req_token)
    return None


def load_scope_registry() -> dict[tuple[str, str], dict]:
    if not SCOPE_REG.exists():
        return {}
    data = yaml.safe_load(SCOPE_REG.read_text()) or {}
    out: dict[tuple[str, str], dict] = {}
    for row in data.get("scopes") or []:
        impl = row.get("impl", "")
        req = row.get("req", "")
        if impl and req:
            out[(impl, req)] = row
    return out


def req_primary_impls(req_token: str) -> set[str]:
    owners: set[str] = set()
    detail = load_req_detail(req_token)
    if not detail:
        return owners
    impl_field = detail.get("implementation") or ""
    if isinstance(impl_field, str):
        for m in re.findall(r"IMPL-[A-Z0-9_]+", impl_field):
            owners.add(m)
    trace = detail.get("traceability") or {}
    for impl in trace.get("implementation") or []:
        if isinstance(impl, str) and impl.startswith("IMPL-"):
            owners.add(impl)
    return owners


def infer_scope_role(impl_token: str, req_token: str, registry: dict[tuple[str, str], dict]) -> tuple[str, list[str] | None]:
    """Return (role, applicable_crit_ids or None for all)."""
    key = (impl_token, req_token)
    if key in registry:
        row = registry[key]
        role = row.get("role", "auxiliary")
        ids = row.get("applicable_crit_ids")
        if ids is None:
            ids = []
        return role, list(ids)
    if impl_token in req_primary_impls(req_token):
        return "primary", None
    return "auxiliary", []


def criteria_for_pair(
    impl_token: str,
    req_token: str,
    all_criteria: dict[str, list[Criterion]],
    mode: str,
    registry: dict[tuple[str, str], dict],
) -> list[Criterion]:
    crits = all_criteria.get(req_token, [])
    if mode == "strict" or not impl_token:
        return crits
    role, applicable = infer_scope_role(impl_token, req_token, registry)
    if role == "auxiliary" and applicable == []:
        return []
    if applicable is None or role == "primary":
        return crits
    id_set = set(applicable)
    return [c for c in crits if c.crit_id in id_set]


def load_all_criteria() -> dict[str, list[Criterion]]:
    out: dict[str, list[Criterion]] = {}
    seen: set[str] = set()
    for req_dir in REQ_DIRS:
        if not req_dir.exists():
            continue
        for path in sorted(req_dir.glob("REQ-*.yaml")):
            token = path.stem
            if token in seen:
                continue
            seen.add(token)
            data = yaml.safe_load(path.read_text())
            rec = _detail_record(data, token)
            if not rec:
                continue
            crits: list[Criterion] = []
            idx = 1
            for item in rec.get("satisfaction_criteria") or []:
                if isinstance(item, dict):
                    text = str(item.get("criterion") or item.get("coverage") or "").strip()
                else:
                    text = str(item).strip()
                if text:
                    crits.append(
                        Criterion(
                            req_token=token,
                            crit_id=f"CRIT-{token}-{idx:03d}",
                            text=text,
                            kind="satisfaction",
                        )
                    )
                    idx += 1
            vidx = 1
            for item in rec.get("validation_criteria") or []:
                if isinstance(item, dict):
                    text = str(item.get("coverage") or item.get("method") or item.get("criterion") or "").strip()
                else:
                    text = str(item).strip()
                if text:
                    crits.append(
                        Criterion(
                            req_token=token,
                            crit_id=f"CRIT-{token}-V{vidx:03d}",
                            text=text,
                            kind="validation",
                        )
                    )
                    vidx += 1
            if crits:
                out[token] = crits
    return out


def criterion_keywords(text: str) -> list[str]:
    words = re.findall(r"[a-zA-Z][a-zA-Z0-9_]{3,}", text.lower())
    return [w for w in words if w not in STOP_WORDS]


def explicit_crit_tags(sidecar_text: str) -> set[str]:
    return set(re.findall(r"CRIT-REQ-[A-Z0-9_]+(?:-V?\d+)?", sidecar_text))


def parse_block_sections(sidecar_text: str) -> list[dict]:
    """Split sidecar into H2 blocks with coverage text."""
    lines = sidecar_text.splitlines()
    blocks: list[dict] = []
    current: dict | None = None
    proc_lines: list[str] = []
    in_proc = False

    def flush_proc() -> None:
        nonlocal in_proc, proc_lines
        if current is not None and proc_lines:
            current["procedures"].append("\n".join(proc_lines))
        proc_lines = []
        in_proc = False

    for line in lines:
        if line.startswith("## "):
            flush_proc()
            if current is not None:
                blocks.append(current)
            name = line[3:].strip()
            current = {
                "name": name,
                "lead": "",
                "steps": [],
                "pre": [],
                "post": [],
                "branches": [],
                "errors": [],
                "procedures": [],
            }
            continue
        if current is None:
            continue
        if line.startswith("- [IMPL-"):
            current["lead"] = line
            continue
        if re.match(r"^STEP T\d+:", line):
            current["steps"].append(line)
            continue
        if line.startswith("PRE:"):
            current["pre"].append(line)
            continue
        if line.startswith("POST:"):
            current["post"].append(line)
            continue
        if re.match(r"^BRANCH B\d+", line):
            current["branches"].append(line)
            continue
        if re.match(r"^ERROR E\d+", line):
            current["errors"].append(line)
            continue
        if line.startswith("PROCEDURE "):
            flush_proc()
            in_proc = True
            proc_lines = [line]
            continue
        if in_proc:
            if line.startswith("## ") or (line.startswith("- [IMPL-") and proc_lines):
                flush_proc()
                if line.startswith("## "):
                    # re-process this line — handled on next iter by restructuring
                    # simpler: push back by treating as non-proc
                    pass
            else:
                proc_lines.append(line)
                continue
    flush_proc()
    if current is not None:
        blocks.append(current)
    return blocks


def req_tokens_in_text(text: str) -> set[str]:
    return {m.strip("[]") for m in re.findall(r"\[REQ-[A-Z0-9_]+\]", text)}


def coverage_corpus_for_req(blocks: list[dict], req_token: str, h1_text: str, full_text: str) -> str:
    req_tag = f"[{req_token}]"
    if req_tag not in full_text:
        return ""
    # Sidecar-wide corpus: any REQ cited in the IMPL may satisfy criteria anywhere in the sidecar.
    return full_text.lower()


def is_meta_validation_criterion(text: str) -> bool:
    lower = text.lower()
    if ".yaml" in lower or ".yml" in lower or lower.startswith("tied/"):
        return True
    if "benchmark" in lower and "performance validation" in lower:
        return True
    if "performance validation" in lower and ("reflection" in lower or "caching" in lower):
        return True
    return False


def criterion_covered(crit: Criterion, corpus: str, explicit_tags: set[str]) -> bool:
    if crit.kind == "validation" and is_meta_validation_criterion(crit.text):
        return True
    if crit.crit_id in explicit_tags:
        return True
    for tag in explicit_tags:
        if crit.req_token in tag and crit.crit_id.split("-")[-1] in tag:
            return True
    norm_corpus = corpus.lower()
    clean_text = re.sub(r"\*+", "", crit.text)
    norm_text = clean_text.lower()
    if len(norm_text) >= 10 and norm_text[:35] in norm_corpus:
        return True
    keywords = criterion_keywords(clean_text)
    if not keywords:
        return True
    hits = sum(1 for kw in keywords if kw in norm_corpus)
    if hits >= 2:
        return True
    if hits == 1 and len(keywords) == 1:
        return True
    if hits >= 1 and len(keywords) <= 4:
        return True
    # REQ token stem words in sidecar (e.g. configuration, maintainability)
    stem = crit.req_token.replace("REQ-", "").replace("_", " ").lower()
    for part in stem.split():
        if len(part) >= 5 and part in norm_corpus and hits >= 1:
            return True
    return False


def find_orphan_criteria(
    sidecar_text: str,
    impl_token: str = "",
    mode: str = "strict",
) -> list[str]:
    """Return CRIT-001 issue strings for uncovered REQ criteria in this sidecar."""
    all_criteria = load_all_criteria()
    registry = load_scope_registry()
    h1 = ""
    for line in sidecar_text.splitlines():
        if line.startswith("# "):
            h1 = line
            break
    blocks = parse_block_sections(sidecar_text)
    reqs = req_tokens_in_text(sidecar_text)
    if not reqs:
        return []
    explicit = explicit_crit_tags(sidecar_text)
    orphans: list[str] = []
    for req in sorted(reqs):
        crits = criteria_for_pair(impl_token, req, all_criteria, mode, registry)
        if not crits:
            continue
        corpus = coverage_corpus_for_req(blocks, req, h1, sidecar_text)
        if not corpus:
            continue
        for crit in crits:
            if not criterion_covered(crit, corpus, explicit):
                snippet = crit.text[:60].replace("\n", " ")
                orphans.append(f"CRIT-001:orphan:{crit.crit_id}:{req}:{snippet}")
    return orphans


def orphan_details(
    sidecar_text: str,
    impl_token: str = "",
    mode: str = "scoped",
) -> list[dict]:
    """Structured orphan rows for report_crit001_orphans.py."""
    all_criteria = load_all_criteria()
    registry = load_scope_registry()
    blocks = parse_block_sections(sidecar_text)
    h1 = ""
    for line in sidecar_text.splitlines():
        if line.startswith("# "):
            h1 = line
            break
    explicit = explicit_crit_tags(sidecar_text)
    rows: list[dict] = []
    for req in sorted(req_tokens_in_text(sidecar_text)):
        role, applicable = infer_scope_role(impl_token, req, registry)
        crits = criteria_for_pair(impl_token, req, all_criteria, mode, registry)
        corpus = coverage_corpus_for_req(blocks, req, h1, sidecar_text)
        for crit in crits:
            covered = criterion_covered(crit, corpus, explicit) if corpus else False
            if not covered:
                rows.append(
                    {
                        "impl": impl_token,
                        "req": req,
                        "crit_id": crit.crit_id,
                        "role": role,
                        "text": crit.text[:80],
                        "suggested_action": (
                            "expand STEP/PROCEDURE or add CRIT tag"
                            if role == "primary"
                            else "scope registry subset or sidecar map"
                        ),
                    }
                )
    return rows
