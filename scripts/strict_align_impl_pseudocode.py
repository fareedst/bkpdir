#!/usr/bin/env python3
"""Strict [PROC-IMPL_PSEUDOCODE_TOKENS] alignment for IMPL essence_pseudocode.

1) Every `#` comment line in essence_pseudocode must include at least one [IMPL-*], one
   [ARCH-*], and one [REQ-*] (filled from this file's token + primary ARCH/REQ from metadata).
2) Top-level block starters (procedure, function, DATA/type/interface, macros, VALIDATE,
   class, map, adapter, contract INPUT/OUTPUT/CONTROL lines, etc.) must be immediately
   preceded (ignoring blank lines) by a comment with the full triplet.

Indented body lines do not each require their own comment (they belong to the preceding block).
"""

from __future__ import annotations

import re
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
IMPL_DIR = ROOT / "tied" / "implementation-decisions"

BLOCK_START = re.compile(
    r"^(procedure\s+\S+|data\s+\S+|type\s+\S+|interface\s+\S+|struct\s+\S+|"
    r"function\s+\S+|class\s+\S+|map\s+\S+|adapter\s+\S+|VALIDATE\s+\S+|"
    r"FOR\s+each|FOR\s+\w|WHEN\s+|"
    r"(?:INPUT|OUTPUT|CONTROL):|"
    r"(?:Added|Extended|Implemented|Removed)\s+\S+|"
    r"(?:Eliminated|Achieved|Reduced|Increased)\s+\S+|"
    r"TEST\s+\S+|"
    r"[A-Z][A-Z0-9_]*\([^:]*\)\s*:)",
    re.IGNORECASE,
)


def pick_arch(body: dict) -> list[str]:
    crs = body.get("cross_references") or []
    arch = [x for x in crs if isinstance(x, str) and x.startswith("ARCH-")]
    tr = (body.get("traceability") or {}).get("architecture") or []
    for x in tr:
        if isinstance(x, str) and x.startswith("ARCH-") and x not in arch:
            arch.append(x)
    return arch or ["ARCH-SYSTEM_COMPONENTS"]


def pick_req(body: dict) -> list[str]:
    crs = body.get("cross_references") or []
    req = [x for x in crs if isinstance(x, str) and x.startswith("REQ-")]
    tr = (body.get("traceability") or {}).get("requirements") or []
    for x in tr:
        if isinstance(x, str) and x.startswith("REQ-") and x not in req:
            req.append(x)
    return req or ["REQ-MAINTAINABILITY"]


def triplet_prefix(impl: str, arch: list[str], req: list[str]) -> str:
    a0, r0 = arch[0], req[0]
    return f"# [{impl}] [{a0}] [{r0}]"


def align_comment_line(line: str, impl: str, arch: list[str], req: list[str]) -> str:
    """Ensure every `#` line carries IMPL + ARCH + REQ bracket tokens."""
    s = line.rstrip("\n")
    if not s.strip().startswith("#"):
        return line
    has_impl = bool(re.search(r"\[IMPL-", s))
    has_arch = bool(re.search(r"\[ARCH-", s))
    has_req = bool(re.search(r"\[REQ-", s))
    if has_impl and has_arch and has_req:
        return line
    add: list[str] = []
    if not has_impl:
        add.append(f"[{impl}]")
    if not has_arch:
        add.append(f"[{arch[0]}]")
    if not has_req:
        add.append(f"[{req[0]}]")
    return s + " " + " ".join(add) + "\n"


def leading_indent(line: str) -> int:
    return len(line) - len(line.lstrip())


def is_top_level_block_start(line: str) -> bool:
    if not line.strip():
        return False
    if line.strip().startswith("#"):
        return False
    # Indented lines are body of current block
    if leading_indent(line) > 0:
        return False
    return bool(BLOCK_START.match(line.strip()))


def comment_has_full_triplet(s: str, impl: str) -> bool:
    if not s.strip().startswith("#"):
        return False
    return bool(re.search(r"\[IMPL-", s) and re.search(r"\[ARCH-", s) and re.search(r"\[REQ-", s))


def process_epilogue(ep: str, impl: str, arch: list[str], req: list[str]) -> tuple[str, bool]:
    lines = ep.splitlines(keepends=False)
    out: list[str] = []
    for line in lines:
        if line.strip().startswith("#"):
            new_line = align_comment_line(line + "\n", impl, arch, req).rstrip("\n")
            out.append(new_line)
            continue
        if not line.strip():
            out.append(line)
            continue
        if is_top_level_block_start(line):
            p = len(out) - 1
            while p >= 0 and not out[p].strip():
                p -= 1
            ok = p >= 0 and comment_has_full_triplet(out[p], impl)
            if not ok:
                stub = (
                    f"{triplet_prefix(impl, arch, req)} "
                    f"— Block implements documented behavior for: {line.strip()[:120]}"
                )
                out.append(stub)
        out.append(line)
    result = "\n".join(out)
    if ep.endswith("\n"):
        result += "\n"
    return result, result != ep


def process_file(path: Path) -> bool:
    text = path.read_text(encoding="utf-8")
    data = yaml.safe_load(text)
    if not data or len(data) != 1:
        return False
    token, body = next(iter(data.items()))
    if not str(token).startswith("IMPL-"):
        return False
    ep = body.get("essence_pseudocode")
    if not ep or not isinstance(ep, str):
        return False
    arch = pick_arch(body)
    req = pick_req(body)
    new_ep, changed = process_epilogue(ep, token, arch, req)
    if new_ep == ep:
        return False
    body["essence_pseudocode"] = new_ep
    with open(path, "w", encoding="utf-8") as f:
        yaml.dump(data, f, default_flow_style=False, allow_unicode=True, sort_keys=False, width=1000)
    return True


def main() -> None:
    n = 0
    for path in sorted(IMPL_DIR.glob("IMPL-*.yaml")):
        if process_file(path):
            n += 1
            print("strict-aligned", path.name)
    print("files updated:", n)


if __name__ == "__main__":
    main()
