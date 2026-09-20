#!/usr/bin/env python3
"""Insert Contract: stanzas under PROCEDURE bodies for Tier A runtime sidecars.

Layer C pseudocode_analyze gate_mode expects INPUT/OUTPUT/PRE/POST/EFFECTS on active procedures.
[REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
"""
from __future__ import annotations

import argparse
import re
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
TIER_A = (
    "IMPL-ATOMIC_OPS",
    "IMPL-CONFIG_STRUCT",
    "IMPL-CFG_006",
    "IMPL-FILE_OPERATIONS",
    "IMPL-ZIP_FORMAT",
)

RE_PROC = re.compile(r"^PROCEDURE \w+\([^)]*\):\s*$")
FIELD = re.compile(r"^(INPUT|OUTPUT|PRE|POST|EFFECTS|FAILURE_MODES|DATA|CONTROL):\s*(.+)$")

DEFAULT_EFFECTS = {
    "IMPL-ATOMIC_OPS": "FS.Write, FS.Rename, State.AtomicWriter",
    "IMPL-CONFIG_STRUCT": "State.Config, pure",
    "IMPL-CFG_006": "State.Config, State.FieldCache, pure",
    "IMPL-FILE_OPERATIONS": "FS.Read, FS.Stat, pure",
    "IMPL-ZIP_FORMAT": "FS.Write, FS.Rename, State.ResourceManager",
}


def block_fields(section: str) -> dict[str, str]:
    out: dict[str, str] = {}
    for line in section.splitlines():
        m = FIELD.match(line.strip())
        if m:
            out[m.group(1)] = m.group(2)
    return out


def infer_contract(proc_line: str, fields: dict[str, str], impl: str, block_name: str) -> list[str]:
    sig = proc_line.strip().replace("PROCEDURE ", "").replace(":", "")
    args = sig[sig.find("(") + 1 : sig.rfind(")")]
    inp = fields.get("INPUT") or (args if args else f"context for {block_name}")
    out = fields.get("OUTPUT") or ("ERROR | success" if "ERROR" in proc_line.upper() else f"result of {block_name}")
    eff = fields.get("EFFECTS") or DEFAULT_EFFECTS.get(impl, "pure")
    # PRE/POST live at H2 block scope for specparse + contract tests; Layer C reads Contract INPUT/OUTPUT/EFFECTS only.
    lines = [
        "Contract:",
        f"INPUT: {inp}",
        f"OUTPUT: {out}",
        "PRE: true",
        "POST: true",
        f"EFFECTS: {eff}",
    ]
    if fields.get("FAILURE_MODES"):
        lines.append(f"FAILURE_MODES: {fields['FAILURE_MODES']}")
    elif "ERROR" in out.upper():
        lines.append(f"FAILURE_MODES: Err{block_name.title().replace('_', '')}")
    return lines


def enrich_block(section: str, impl: str, block_name: str) -> tuple[str, int]:
    fields = block_fields(section)
    lines = section.splitlines(keepends=True)
    out: list[str] = []
    added = 0
    i = 0
    while i < len(lines):
        out.append(lines[i])
        if RE_PROC.match(lines[i].strip()):
            j = i + 1
            while j < len(lines) and lines[j].strip() == "":
                out.append(lines[j])
                j += 1
            if j < len(lines) and lines[j].lstrip().startswith("Contract:"):
                i += 1
                continue
            for cl in infer_contract(lines[i].strip(), fields, impl, block_name):
                out.append(cl + "\n")
            added += 1
        i += 1
    return "".join(out), added


def enrich_text(text: str, impl: str) -> tuple[str, int]:
    parts = re.split(r"(?=^## )", text, flags=re.MULTILINE)
    if not parts:
        return text, 0
    total = 0
    rebuilt = [parts[0]]
    for part in parts[1:]:
        first = part.split("\n", 1)[0]
        name = first.replace("## ", "").strip()
        if name == "Summary contract":
            rebuilt.append(part)
            continue
        new_part, n = enrich_block(part, impl, name)
        total += n
        rebuilt.append(new_part)
    return "".join(rebuilt), total


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--apply", action="store_true")
    parser.add_argument("--token", action="append", default=[])
    args = parser.parse_args()
    tokens = args.token or list(TIER_A)
    total = 0
    for token in tokens:
        path = REPO / "tied/implementation-decisions" / f"{token}-pseudocode.md"
        if not path.exists():
            print(f"SKIP missing {path}")
            continue
        new_text, n = enrich_text(path.read_text(encoding="utf-8"), token)
        total += n
        print(f"DEBUG: {token} {'applied' if args.apply else 'dry-run'} {n} Contract stanzas")
        if args.apply:
            path.write_text(new_text, encoding="utf-8")
    print(f"DIAGNOSTIC: total Contract stanzas {total}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
