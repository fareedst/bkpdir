#!/usr/bin/env python3
"""Sync literal block-lead comments from IMPL sidecars into Go production and test files."""
from __future__ import annotations

import argparse
import glob
import re
import sys
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"


def parse_sidecar_leads(sidecar_path: Path) -> list[tuple[str, str]]:
    """Return list of (section_name, block_lead_line_without_dash)."""
    text = sidecar_path.read_text()
    leads: list[tuple[str, str]] = []
    section = "SUMMARY"
    for line in text.splitlines():
        if line.startswith("## "):
            section = line[3:].strip()
        m = re.match(r"^-\s+(\[IMPL-[^\n]+)$", line.strip())
        if m:
            leads.append((section, m.group(1)))
    return leads


def go_func_pattern(name: str) -> re.Pattern[str]:
    # Match func (recv) Name( or func Name(
    esc = re.escape(name)
    return re.compile(rf"^func\s+(?:\([^)]+\)\s+)?{esc}\s*\(", re.MULTILINE)


def insert_comment_before_func(content: str, func_name: str, comment_line: str) -> tuple[str, bool]:
    """Insert // comment_line immediately before func if not already present literally."""
    if comment_line in content:
        return content, False
    pat = go_func_pattern(func_name)
    m = pat.search(content)
    if not m:
        return content, False
    insert_at = m.start()
    # Skip existing comment block directly above
    prefix = content[:insert_at].rstrip("\n")
    new_block = f"// - {comment_line}\n"
    return prefix + "\n" + new_block + content[insert_at:], True


def section_to_func_names(section: str) -> list[str]:
    """Map sidecar section names to likely Go function names."""
    if section.startswith("EMBEDDED_MINITEST"):
        return []
    if section.startswith("ATOMICWRITER_"):
        method = section.split("_", 1)[1].capitalize()
        if method == "Cleanup":
            return ["cleanup"]
        return [method]
    if section.startswith("RESOURCE_"):
        rest = section.replace("RESOURCE_", "")
        parts = rest.lower().split("_")
        return ["".join(p.capitalize() for p in parts)]
    parts = section.lower().split("_")
    camel = "".join(p.capitalize() for p in parts)
    names = [camel]
    if section == "FORMAT_LIST_ARCHIVE_SIMPLE":
        names = ["formatListArchiveSimple"]
    return names


def resolve_file_for_package(pkg: str, files: list[str]) -> str | None:
    prefix = f"pkg/{pkg}/"
    for f in files:
        if f.startswith(prefix):
            return f
    for f in files:
        if f"/{pkg}/" in f or f.endswith(f"/{pkg}.go"):
            return f
    return files[0] if files else None


def load_impl_functions(yaml_path: Path) -> list[tuple[str, str]]:
    with open(yaml_path) as f:
        data = yaml.safe_load(f)
    token = list(data.keys())[0]
    rec = data[token]
    raw_files = rec.get("code_locations", {}).get("files") or []
    files: list[str] = []
    for f in raw_files:
        if isinstance(f, str):
            files.append(f)
        elif isinstance(f, dict):
            p = f.get("path") or f.get("file")
            if p:
                files.append(str(p))
    out: list[tuple[str, str]] = []
    seen: set[tuple[str, str]] = set()
    for fn in rec.get("code_locations", {}).get("functions") or []:
        if isinstance(fn, dict):
            name = fn.get("name", "")
            if "." in name:
                name = name.split(".")[-1]
            file = fn.get("file")
            if name and file:
                key = (name, file)
                if key not in seen:
                    seen.add(key)
                    out.append(key)
        elif isinstance(fn, str):
            name = fn.strip()
            m = re.match(r"^([\w]+)\.([\w]+)", name)
            if m:
                pkg, name = m.group(1), m.group(2)
                rel = resolve_file_for_package(pkg, files)
            else:
                rel = files[0] if files else None
            if name and rel:
                key = (name, rel)
                if key not in seen:
                    seen.add(key)
                    out.append(key)
    return out


def sync_token(token: str, dry_run: bool = False) -> dict:
    sidecar = IMPL_DIR / f"{token}-pseudocode.md"
    yaml_path = IMPL_DIR / f"{token}.yaml"
    if not sidecar.exists() or not yaml_path.exists():
        return {"token": token, "updated": 0, "skipped": True}

    leads = parse_sidecar_leads(sidecar)
    if not leads:
        return {"token": token, "updated": 0, "skipped": True}

    lead_by_section = {s: lead for s, lead in leads}
    updated = 0
    files_touched: set[str] = set()

    with open(yaml_path) as f:
        yaml_data = yaml.safe_load(f)
    token_key = list(yaml_data.keys())[0]
    code_files = yaml_data[token_key].get("code_locations", {}).get("files") or []

    func_pairs = load_impl_functions(yaml_path)
    for sec, ld in leads:
        if sec.startswith("EMBEDDED"):
            continue
        for fname in section_to_func_names(sec):
            if not fname:
                continue
            for rel in code_files:
                key = (fname, rel)
                if key not in func_pairs:
                    func_pairs.append(key)
    clean_pairs: list[tuple[str, str]] = []
    for item in func_pairs:
        if isinstance(item, tuple) and len(item) == 2 and all(isinstance(x, str) for x in item):
            clean_pairs.append(item)
    func_pairs = list(dict.fromkeys(clean_pairs))
    for func_name, rel_file in func_pairs:
        lead = None
        for sec, ld in leads:
            for fname in section_to_func_names(sec):
                if fname == func_name:
                    lead = ld
                    break
            if lead:
                break
            if sec.replace("_", "").lower() == func_name.replace("_", "").lower():
                lead = ld
                break
            if func_name.lower() in sec.lower().replace("_", ""):
                lead = ld
                break
        if not lead:
            for sec, ld in leads:
                if not sec.startswith("EMBEDDED"):
                    lead = ld
                    break
        if not lead:
            continue

        path = REPO / rel_file
        if not path.exists():
            continue
        content = path.read_text()
        new_content, changed = insert_comment_before_func(content, func_name, lead)
        if changed:
            updated += 1
            files_touched.add(str(rel_file))
            if not dry_run:
                path.write_text(new_content)

    # Sync test files: use EMBEDDED_MINITEST leads for Test* functions
    test_leads = [(s, l) for s, l in leads if s.startswith("EMBEDDED_MINITEST")]
    if test_leads:
        r = __import__("subprocess").run(
            ["rg", "-l", rf"\[{token}\]", "--glob", "*_test.go", "."],
            cwd=REPO,
            capture_output=True,
            text=True,
        )
        test_files = [x for x in r.stdout.strip().split("\n") if x]
        for tf in test_files:
            path = REPO / tf
            content = path.read_text()
            for i, (sec, lead) in enumerate(test_leads):
                # Match test func containing token
                for m in re.finditer(rf"^func (Test[^(]+)\(", content, re.MULTILINE):
                    tname = m.group(1)
                    if token.replace("IMPL-", "") not in tname and i > 0:
                        continue
                    if lead in content:
                        continue
                    insert_at = m.start()
                    prefix = content[:insert_at].rstrip("\n")
                    new_block = f"// - {lead}\n"
                    content = prefix + "\n" + new_block + content[insert_at:]
                    updated += 1
                    files_touched.add(tf)
            if not dry_run and tf in files_touched:
                path.write_text(content)

    return {"token": token, "updated": updated, "files": sorted(files_touched)}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--token", action="append")
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    tokens = args.token or [
        p.stem.replace("-pseudocode", "")
        for p in sorted(IMPL_DIR.glob("IMPL-*-pseudocode.md"))
    ]

    total = 0
    for token in tokens:
        if not token.startswith("IMPL-"):
            token = f"IMPL-{token}"
        r = sync_token(token, dry_run=args.dry_run)
        if not r.get("skipped"):
            print(f"{r['token']}: updated={r['updated']} files={r.get('files', [])}")
            total += r["updated"]
    print(f"Total comment insertions: {total}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
