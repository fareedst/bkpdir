#!/usr/bin/env python3
"""Migrate inline IMPL essence_pseudocode to sidecar markdown files.

Track C brownfield migration: extract inline YAML essence, de-Go-ify, write
tied/implementation-decisions/IMPL-*-pseudocode.md, optionally register via tied-cli.
"""
from __future__ import annotations

import argparse
import glob
import os
import re
import subprocess
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
TIED = REPO / "tied"
IMPL_DIR = TIED / "implementation-decisions"
TIED_CLI = REPO / ".cursor/skills/tied-yaml/scripts/tied-cli.sh"


def de_goify(text: str) -> str:
    """Replace common Go-specific idioms with language-agnostic TIED vocabulary."""
    out = text
    replacements = [
        (r"\bfmt\.Sprintf\b", "FORMAT_WITH_PLACEHOLDERS"),
        (r"\bstrings\.Contains\b", "CONTAINS"),
        (r"\bos\.Rename\b", "ATOMIC_RENAME"),
        (r"\bos\.Create\b", "CREATE_FILE"),
        (r"\bioutil\.TempFile\b", "CREATE_TEMP_FILE"),
        (r"`([^`]+)`", r"\1"),
        (r"\bprocedure\b", "PROCEDURE"),
    ]
    for pat, repl in replacements:
        out = re.sub(pat, repl, out)
    # Strip Go struct field tags
    out = re.sub(r"\s+`yaml:\"[^\"]*\"`", "", out)
    out = re.sub(r"\s+`json:\"[^\"]*\"`", "", out)
    return out


def essence_to_sidecar(token: str, essence: str, arch: list[str], req: list[str]) -> str:
    """Convert inline essence string to sidecar markdown with H2 blocks."""
    essence = de_goify(essence.strip())
    if not essence:
        return ""

    arch_tokens = " ".join(f"[{a}]" for a in arch[:3]) if arch else ""
    req_tokens = " ".join(f"[{r}]" for r in req[:3]) if req else ""
    header_parts = [f"[{token}]"]
    if arch_tokens:
        header_parts.append(arch_tokens)
    if req_tokens:
        header_parts.append(req_tokens)
    header = "# " + " ".join(header_parts)

    lines = essence.splitlines()
    body_lines: list[str] = []
    current_h2: str | None = None

    proc_re = re.compile(
        r"^(?:#+\s*)?(?:PROCEDURE|procedure)\s+([A-Za-z0-9_]+)",
        re.IGNORECASE,
    )
    block_re = re.compile(r"^#\s+\[IMPL-", re.IGNORECASE)

    for line in lines:
        proc_m = proc_re.match(line.strip())
        if proc_m:
            name = proc_m.group(1).upper()
            if current_h2 != name:
                if body_lines and body_lines[-1] != "":
                    body_lines.append("")
                body_lines.append(f"## {name}")
                body_lines.append("")
                current_h2 = name
            body_lines.append(line if line.startswith("PROCEDURE") else re.sub(r"^procedure", "PROCEDURE", line, flags=re.I))
            continue
        if block_re.match(line.strip()) and "Block implements" in line:
            # Keep as block lead under current or new section
            if not line.strip().startswith("- "):
                line = "- " + line.strip().lstrip("# ").strip()
            body_lines.append(line)
            continue
        body_lines.append(line)

    if not any(l.startswith("## ") for l in body_lines):
        # Fallback: single runtime block
        body = "\n".join(body_lines).strip()
        return f"{header}\n\n## Summary contract\n\n{body}\n"

    return f"{header}\n\n## Summary contract\n\nLanguage-agnostic implementation logic for [{token}].\n\n" + "\n".join(body_lines).strip() + "\n"


def count_blocks(sidecar_text: str) -> int:
    return len(re.findall(r"^## ", sidecar_text, re.MULTILINE))


def register_sidecar(token: str, rel_path: str, dry_run: bool) -> bool:
    if dry_run:
        return True
    env = os.environ.copy()
    env["TIED_BASE_PATH"] = str(TIED)
    payload = {
        "token": token,
        "essence_pseudocode_path": rel_path,
        "metadata_last_updated": {
            "date": "2026-06-02",
            "author": "AI agent",
            "reason": "Track C bulk sidecar migration from inline essence",
        },
    }
    import json

    r = subprocess.run(
        [str(TIED_CLI), "impl_detail_set_essence_pseudocode", json.dumps(payload)],
        env=env,
        capture_output=True,
        text=True,
    )
    if r.returncode != 0:
        print(f"ERROR registering {token}: {r.stderr or r.stdout}", file=sys.stderr)
        return False
    return True


def migrate_token(path: Path, register: bool, dry_run: bool, skip_existing: bool) -> dict:
    with open(path) as f:
        data = yaml.safe_load(f)
    token = list(data.keys())[0]
    rec = data[token]
    sidecar_name = f"{token}-pseudocode.md"
    sidecar_path = IMPL_DIR / sidecar_name
    rel_path = f"implementation-decisions/{sidecar_name}"

    if skip_existing and sidecar_path.exists() and token != "IMPL-LIST_FORMAT_SAFETY":
        text = sidecar_path.read_text()
        return {
            "token": token,
            "skipped": True,
            "blocks": count_blocks(text),
            "registered": False,
        }

    trace = rec.get("traceability") or {}
    arch = trace.get("architecture") or rec.get("cross_references") or []
    req = trace.get("requirements") or []
    # Filter to ARCH/REQ only
    arch = [t for t in arch if str(t).startswith("ARCH-")]
    req = [t for t in req if str(t).startswith("REQ-")]
    if not arch:
        arch = [t for t in (rec.get("cross_references") or []) if str(t).startswith("ARCH-")]
    if not req:
        req = [t for t in (rec.get("cross_references") or []) if str(t).startswith("REQ-")]

    essence = rec.get("essence_pseudocode") or ""
    sidecar_text = essence_to_sidecar(token, essence, arch, req)
    if not dry_run:
        sidecar_path.write_text(sidecar_text)

    ok = True
    if register and not dry_run:
        ok = register_sidecar(token, rel_path, dry_run=False)

    return {
        "token": token,
        "skipped": False,
        "blocks": count_blocks(sidecar_text),
        "registered": ok,
        "sidecar": str(sidecar_path),
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--register", action="store_true", help="Register sidecars via tied-cli")
    parser.add_argument("--dry-run", action="store_true")
    parser.add_argument("--skip-existing", action="store_true", default=True)
    parser.add_argument("--token", action="append", help="Process only these IMPL tokens")
    args = parser.parse_args()

    paths = sorted(IMPL_DIR.glob("IMPL-*.yaml"))
    if args.token:
        wanted = set(args.token)
        paths = [p for p in paths if any(p.stem == t for t in wanted)]

    results = []
    for path in paths:
        if path.name.endswith("-pseudocode.yaml"):
            continue
        r = migrate_token(path, args.register, args.dry_run, args.skip_existing)
        results.append(r)
        status = "skip" if r.get("skipped") else ("ok" if r.get("registered", True) else "fail")
        print(f"{r['token']}: {status} blocks={r.get('blocks', 0)}")

    failed = [r for r in results if not r.get("skipped") and r.get("registered") is False]
    return 1 if failed else 0


if __name__ == "__main__":
    sys.exit(main())
