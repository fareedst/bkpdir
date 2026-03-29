#!/usr/bin/env python3
"""Enrich IMPL essence_pseudocode comment lines with [ARCH-*] and [REQ-*] where only [IMPL-*] appears.
Uses each detail file's cross_references (ARCH) and traceability.requirements (REQ). [PROC-IMPL_PSEUDOCODE_TOKENS]

For strict alignment (every `#` line + block starts), run `strict_align_impl_pseudocode.py` after this script.
"""

from __future__ import annotations

import re
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
IMPL_DIR = ROOT / "tied" / "implementation-decisions"


def pick_arch_refs(body: dict) -> list[str]:
    crs = body.get("cross_references") or []
    arch = [x for x in crs if isinstance(x, str) and x.startswith("ARCH-")]
    tr = (body.get("traceability") or {}).get("architecture") or []
    for x in tr:
        if isinstance(x, str) and x.startswith("ARCH-") and x not in arch:
            arch.append(x)
    return arch[:3]  # keep primary refs concise in repeated comments


def pick_req_refs(body: dict) -> list[str]:
    crs = body.get("cross_references") or []
    req = [x for x in crs if isinstance(x, str) and x.startswith("REQ-")]
    tr = (body.get("traceability") or {}).get("requirements") or []
    for x in tr:
        if isinstance(x, str) and x.startswith("REQ-") and x not in req:
            req.append(x)
    return req[:4]


def enrich_line(line: str, impl_token: str, arch: list[str], req: list[str]) -> str:
    """Ensure comment lines that reference this IMPL also name ARCH and REQ when missing."""
    if not line.strip().startswith("#"):
        return line
    if "[" not in line:
        return line
    if not re.search(r"\[" + re.escape(impl_token) + r"\]", line):
        return line
    stripped = line.rstrip("\n")
    add = []
    if not re.search(r"\[ARCH-", stripped):
        add.extend(f"[{a}]" for a in arch)
    if not re.search(r"\[REQ-", stripped):
        add.extend(f"[{r}]" for r in req)
    if not add:
        return line
    return stripped + " " + " ".join(add) + "\n"


def process_file(path: Path) -> bool:
    with open(path, encoding="utf-8") as f:
        data = yaml.safe_load(f)
    if not data or len(data) != 1:
        return False
    token, body = next(iter(data.items()))
    if not token.startswith("IMPL-"):
        return False
    ep = body.get("essence_pseudocode")
    if not ep or not isinstance(ep, str):
        return False
    arch = pick_arch_refs(body)
    req = pick_req_refs(body)
    if not arch and not req:
        return False
    lines = ep.splitlines(keepends=True)
    out = []
    changed = False
    for line in lines:
        new_line = enrich_line(line, token, arch, req)
        if new_line != line:
            changed = True
        out.append(new_line)
    if not changed:
        return False
    body["essence_pseudocode"] = "".join(out)
    with open(path, "w", encoding="utf-8") as f:
        yaml.dump(data, f, default_flow_style=False, allow_unicode=True, sort_keys=False, width=1000)
    return True


def main() -> None:
    n = 0
    for path in sorted(IMPL_DIR.glob("IMPL-*.yaml")):
        if process_file(path):
            n += 1
            print("enriched", path.name)
    print("total updated:", n)


if __name__ == "__main__":
    main()
