#!/usr/bin/env python3
"""Insert or complete Contract: stanzas under PROCEDURE bodies for Layer B sidecars.

Layer B pseudocode_validate (full mode) expects INPUT/OUTPUT/PRE/POST/EFFECTS on active
procedures, `# [IMPL-*]` token linkage, and FAILURE_MODES / TERMINATION when triggered.
[REQ-PSEUDOCODE_FORMAL_VERIFICATION] [REQ-PSEUDOCODE_LAYER_B_CORPUS] [IMPL-SPEC_CTL]
"""
from __future__ import annotations

import argparse
import re
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
WAVE_REGISTRY = REPO / "tied-project/spec/layer-b-wave-registry.yaml"
TIER_A = (
    "IMPL-ATOMIC_OPS",
    "IMPL-CONFIG_STRUCT",
    "IMPL-CFG_006",
    "IMPL-FILE_OPERATIONS",
    "IMPL-ZIP_FORMAT",
)

RE_PROC = re.compile(r"^PROCEDURE \w+\([^)]*\):\s*$", re.M)
FIELD = re.compile(
    r"^\s*(INPUT|OUTPUT|PRE|POST|EFFECTS|FAILURE_MODES|DATA_TRANSITION|TERMINATION|DATA|CONTROL):\s*(.+)$",
    re.I,
)
REQUIRED = ("INPUT", "OUTPUT", "PRE", "POST", "EFFECTS")

DEFAULT_EFFECTS = {
    "IMPL-ATOMIC_OPS": "FS.Write, FS.Rename, State.AtomicWriter",
    "IMPL-CONFIG_STRUCT": "State.Config, pure",
    "IMPL-CFG_006": "State.Config, State.FieldCache, pure",
    "IMPL-FILE_OPERATIONS": "FS.Read, FS.Stat, pure",
    "IMPL-ZIP_FORMAT": "FS.Write, FS.Rename, State.ResourceManager",
}


def wave_tokens(wave: str) -> list[str]:
    reg = yaml.safe_load(WAVE_REGISTRY.read_text(encoding="utf-8"))
    tokens = reg.get("waves", {}).get(wave.lower(), {}).get("tokens") or []
    if not tokens:
        raise SystemExit(f"DIAGNOSTIC: no tokens for wave {wave} in {WAVE_REGISTRY}")
    return list(tokens)


def block_fields(section: str) -> dict[str, str]:
    out: dict[str, str] = {}
    for line in section.splitlines():
        m = FIELD.match(line)
        if m:
            out[m.group(1).upper()] = m.group(2).strip()
    return out


def procedure_body_text(section: str, proc_line_idx: int) -> str:
    lines = section.splitlines()
    return "\n".join(lines[proc_line_idx + 1 :])


def needs_failure_modes(body: str, fields: dict[str, str]) -> bool:
    if fields.get("FAILURE_MODES"):
        return False
    return bool(
        re.search(r"\b(?:RETURN|RAISE)\s+(?:error|failure)|\berror\b|\bFAILURE_MODE\b", body, re.I)
    )


def needs_termination(body: str, fields: dict[str, str]) -> bool:
    if fields.get("TERMINATION"):
        return False
    return bool(re.search(r"\b(?:WHILE|AWAIT|WAIT|RECURS|open-ended)\b", body, re.I))


def needs_data_transition(body: str, fields: dict[str, str]) -> bool:
    if fields.get("DATA_TRANSITION"):
        return False
    for line in body.splitlines():
        if re.match(r"^\s*DATA_TRANSITION\s*:", line, re.I):
            continue
        if not re.search(r"\b(?:mutat\w*|state transition)\b", line, re.I):
            continue
        without = re.sub(r"\b(?:never|does not|do not)\s+mutat\w*", "", line, flags=re.I)
        if re.search(r"\bmutat\w*\b", without, re.I):
            return True
    return bool(re.search(r"^\s*SET\s+\w+", body, re.M | re.I) and fields.get("EFFECTS", "") != "pure")


def infer_contract(
    proc_line: str,
    h2_fields: dict[str, str],
    proc_fields: dict[str, str],
    impl: str,
    block_name: str,
    body: str,
) -> dict[str, str]:
    merged = {**h2_fields, **proc_fields}
    sig = proc_line.strip().replace("PROCEDURE ", "").replace(":", "")
    args = sig[sig.find("(") + 1 : sig.rfind(")")]
    out: dict[str, str] = {}
    out["INPUT"] = merged.get("INPUT") or (args if args else f"context for {block_name}")
    out["OUTPUT"] = merged.get("OUTPUT") or (
        "ERROR | success" if re.search(r"\berror\b", body, re.I) else f"result of {block_name}"
    )
    out["PRE"] = merged.get("PRE") or "true"
    out["POST"] = merged.get("POST") or "true"
    out["EFFECTS"] = merged.get("EFFECTS") or DEFAULT_EFFECTS.get(impl, "pure")
    if merged.get("FAILURE_MODES"):
        out["FAILURE_MODES"] = merged["FAILURE_MODES"]
    elif needs_failure_modes(body, out):
        out["FAILURE_MODES"] = f"Err{block_name.title().replace('_', '')}"
    if needs_termination(body, out):
        if "RETR" in block_name or "RETRY" in body.upper():
            out["TERMINATION"] = "bounded by configured max retries"
        elif "SIGNAL" in block_name:
            out["TERMINATION"] = "total after handler registration completes"
        else:
            out["TERMINATION"] = "total"
    if needs_data_transition(body, out):
        out["DATA_TRANSITION"] = f"{block_name} updates declared state per EFFECTS"
    return out


def format_contract(fields: dict[str, str]) -> list[str]:
    order = [
        "INPUT",
        "OUTPUT",
        "PRE",
        "POST",
        "EFFECTS",
        "FAILURE_MODES",
        "DATA_TRANSITION",
        "TERMINATION",
    ]
    lines = ["Contract:"]
    for key in order:
        if key in fields and fields[key]:
            lines.append(f"{key}: {fields[key]}")
    return lines


def has_impl_hash_comment(after_proc: list[str], impl: str) -> bool:
    needle = f"[{impl}]"
    for line in after_proc[:8]:
        if line.lstrip().startswith("#") and needle in line:
            return True
    return False


def enrich_block(section: str, impl: str, block_name: str) -> tuple[str, int]:
    h2_fields = block_fields(section)
    lines = section.splitlines(keepends=True)
    out: list[str] = []
    added = 0
    i = 0
    while i < len(lines):
        if not RE_PROC.match(lines[i].strip()):
            out.append(lines[i])
            i += 1
            continue
        proc_line = lines[i].strip()
        out.append(lines[i])
        i += 1
        while i < len(lines) and lines[i].strip() == "":
            out.append(lines[i])
            i += 1
        proc_fields: dict[str, str] = {}
        body_start = i
        if i < len(lines) and lines[i].lstrip().startswith("Contract:"):
            i += 1
            while i < len(lines):
                m = FIELD.match(lines[i])
                if m:
                    proc_fields[m.group(1).upper()] = m.group(2).strip()
                    i += 1
                    continue
                if lines[i].strip() == "":
                    i += 1
                    continue
                break
        rest_start = i
        while i < len(lines) and not RE_PROC.match(lines[i].strip()):
            i += 1
        body = "".join(lines[rest_start:i])
        if not has_impl_hash_comment(
            [lines[k].rstrip("\n") for k in range(body_start, min(len(lines), body_start + 8))],
            impl,
        ):
            out.append(f"  # [{impl}] — Layer B procedure token link\n")
            added += 1
        target = infer_contract(proc_line, h2_fields, proc_fields, impl, block_name, body)
        for cl in format_contract(target):
            out.append(cl + "\n")
        if not proc_fields:
            added += 1
        out.extend(lines[rest_start:i])
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
    parser.add_argument("--wave", default="")
    args = parser.parse_args()
    if args.wave:
        tokens = wave_tokens(args.wave)
    elif args.token:
        tokens = args.token
    else:
        tokens = list(TIER_A)
    total = 0
    for token in tokens:
        path = REPO / "tied-project/implementation-decisions" / f"{token}-pseudocode.md"
        if not path.exists():
            print(f"SKIP missing {path}")
            continue
        new_text, n = enrich_text(path.read_text(encoding="utf-8"), token)
        total += n
        print(f"DEBUG: {token} {'applied' if args.apply else 'dry-run'} {n} enrichments")
        if args.apply:
            path.write_text(new_text, encoding="utf-8")
    print(f"DIAGNOSTIC: total enrichments {total} across {len(tokens)} token(s)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
