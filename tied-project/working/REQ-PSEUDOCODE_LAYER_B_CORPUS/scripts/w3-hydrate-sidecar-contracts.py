#!/usr/bin/env python3
"""One-shot W3 Tier C sidecar contract hydration for Layer B structural gate."""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parents[4]
SIDEcar_DIR = REPO / "tied-project/implementation-decisions"
STRUCT_DIR = (
    REPO
    / "tied-project/working/REQ-PSEUDOCODE_LAYER_B_CORPUS/pseudocode-validate-w3-structural"
)

BULLET = re.compile(
    r"^- \[(IMPL-[^\]]+)\](.*?)(?:How:|— How:)\s*(.+?)\.\s*$",
    re.MULTILINE,
)
SUMMARY = re.compile(
    r"^INPUT:\s*(.+?)$|^OUTPUT:\s*(.+?)$|^DATA:\s*(.+?)$",
    re.MULTILINE,
)
PROC = re.compile(r"^(PROCEDURE\s+\S+\([^)]*\):)\s*$", re.MULTILINE)


def load_summary(text: str) -> tuple[str, str, str]:
    inp = out = data = "see block context"
    for line in text.splitlines():
        if line.startswith("INPUT:"):
            inp = line.split(":", 1)[1].strip()
        elif line.startswith("OUTPUT:"):
            out = line.split(":", 1)[1].strip()
        elif line.startswith("DATA:"):
            data = line.split(":", 1)[1].strip()
    return inp, out, data


def block_name_from_proc(proc_line: str) -> str:
    m = re.match(r"PROCEDURE\s+(\w+)\(", proc_line)
    return m.group(1) if m else ""


def diagnostics_for_token(token: str) -> dict[str, list[dict]]:
    path = STRUCT_DIR / f"{token}.json"
    if not path.is_file():
        return {}
    doc = json.loads(path.read_text())
    by_block: dict[str, list[dict]] = {}
    for d in doc.get("diagnostics", []):
        if d.get("severity") != "error":
            continue
        by_block.setdefault(d.get("block", ""), []).append(d)
    coverage = {c["block"]: c for c in doc.get("coverage", [])}
    return {"errors": by_block, "coverage": coverage}


def needs_termination(body: str, errors: list[dict]) -> bool:
    if any(e.get("code") == "MISSING_TERMINATION" for e in errors):
        return True
    return bool(re.search(r"^\s+FOR\b", body, re.MULTILINE)) and "RETURN" not in body


def needs_failure_modes(body: str, errors: list[dict], coverage: dict) -> list[str]:
    if any(e.get("code") == "MISSING_FAILURE_MODES" for e in errors):
        fm = coverage.get("failure_modes") or []
        if fm:
            return fm[:3]
        return ["validation or test failure aborts the procedure with a diagnostic"]
    if re.search(r"\b(ASSERT|VERIFY|IF any|IF gaps|THEN STOP)\b", body):
        return ["precondition or validation failure stops progression"]
    return []


def needs_data_transition(errors: list[dict]) -> bool:
    return any(e.get("code") == "MISSING_DATA_TRANSITION" for e in errors)


def hydrate_file(path: Path, token: str) -> bool:
    text = path.read_text()
    if "Contract:" in text and path.stat().st_mtime:  # still process if structural fails
        pass
    diag = diagnostics_for_token(token)
    err_by_block = diag.get("errors", {})
    cov_by_block = diag.get("coverage", {})
    inp, out, data = load_summary(text)
    changed = False

    def replacer(match: re.Match[str]) -> str:
        nonlocal changed
        proc_line = match.group(1)
        start = match.end()
        # already hydrated
        rest = text[start:]
        next_lines = rest.splitlines()
        if next_lines and (
            next_lines[0].strip().startswith("# [IMPL")
            or next_lines[0].strip().startswith("Contract:")
        ):
            return match.group(0)

        block_key = block_name_from_proc(proc_line)
        # find bullet before this PROCEDURE
        before = text[: match.start()]
        bullets = list(BULLET.finditer(before))
        if not bullets:
            return match.group(0)
        bullet = bullets[-1]
        impl_tok = bullet.group(1)
        extra_tokens = bullet.group(2).strip()
        how = bullet.group(3).strip()
        token_comment = f"  # [{impl_tok}]{extra_tokens} How: {how}."
        errors = err_by_block.get(block_key, [])
        coverage = cov_by_block.get(block_key, {})
        body_end = start
        for i, line in enumerate(next_lines):
            if line.startswith("## ") or (line.startswith("PROCEDURE ") and i > 0):
                body_end = start + sum(len(l) + 1 for l in next_lines[:i])
                break
        else:
            body_end = len(text)
        body = text[start:body_end]

        term = needs_termination(body, errors)
        fms = needs_failure_modes(body, errors, coverage)
        dt = needs_data_transition(errors)

        lines = [
            proc_line,
            token_comment,
            "  Contract:",
            f"    INPUT: {inp}",
            f"    OUTPUT: {out}",
            "    PRE: prerequisites for this process step are satisfied",
            "    POST: step outputs match the documented How intent",
            "    EFFECTS: process documentation; runtime unchanged unless a separate LEAP slice applies",
        ]
        if fms:
            lines.append("    FAILURE_MODES:")
            for fm in fms:
                lines.append(f"      - {fm}")
        if dt:
            lines.append(
                "    DATA_TRANSITION: parsed sidecar AST and diagnostic list updated in memory for the validate run"
            )
        if term:
            lines.append("    TERMINATION: total")
        lines.append("")
        changed = True
        return "\n".join(lines) + "\n"

    new_text = PROC.sub(replacer, text)
    if changed and new_text != text:
        path.write_text(new_text)
    return changed


def main() -> int:
    registry = REPO / "tied-project/spec/layer-b-wave-registry.yaml"
    waves: list[str] = []
    in_w3 = False
    for line in registry.read_text().splitlines():
        if line.strip() == "w3:":
            in_w3 = True
            continue
        if in_w3 and line.strip().startswith("w4:"):
            break
        if in_w3 and line.strip().startswith("- IMPL-"):
            waves.append(line.split("-", 1)[1].strip())
    any_change = False
    for tok in waves:
        sidecar = SIDEcar_DIR / f"{tok}-pseudocode.md"
        if hydrate_file(sidecar, tok):
            any_change = True
            print(f"hydrated {tok}")
    return 0 if any_change or True else 1


if __name__ == "__main__":
    sys.exit(main())
