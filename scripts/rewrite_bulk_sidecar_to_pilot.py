#!/usr/bin/env python3
"""Rewrite bulk-migrated IMPL sidecars toward pilot format (in-place, preserves structure)."""
from __future__ import annotations

import argparse
import re
from pathlib import Path

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied/implementation-decisions"

SKIP_TOKENS = {
    "IMPL-LIST_FORMAT_SAFETY",
    "IMPL-CLI_FRAMEWORK",
    "IMPL-CONFIG_STRUCT",
    "IMPL-STRUCTURED_ERRORS",
    "IMPL-ATOMIC_OPS",
}

GO_PATTERNS = [
    (r"\bfmt\.\w+", "FORMAT"),
    (r"\bos\.\w+", "OS"),
    (r"\bio\.\w+", "IO"),
    (r"\breflect\.\w+", "REFLECT"),
    (r"\byaml\.\w+", "YAML"),
    (r"\bcobra\.\w+", "CLI"),
    (r":=\s*", "= "),
    (r"\*cobra\.Command", "CLI_COMMAND"),
]


def parse_h1_tokens(text: str) -> tuple[str, str, str]:
    m = re.search(r"#\s*\[(IMPL-[^\]]+)\]([^\n]*)", text)
    if not m:
        return "IMPL-UNKNOWN", "[ARCH-UNKNOWN]", "[REQ-UNKNOWN]"
    impl = m.group(1)
    rest = m.group(1) + m.group(2)
    archs = " ".join(dict.fromkeys(re.findall(r"\[ARCH-[^\]]+\]", rest)))
    reqs = " ".join(dict.fromkeys(re.findall(r"\[REQ-[^\]]+\]", rest)))
    return impl, archs or "[ARCH-UNKNOWN]", reqs or "[REQ-UNKNOWN]"


def de_goify_line(line: str) -> str:
    s = line
    for pat, repl in GO_PATTERNS:
        s = re.sub(pat, repl, s)
    return s


def how_from_boilerplate(line: str, section_hint: str = "") -> str | None:
    if "Block implements documented behavior" not in line:
        return None
    m = re.search(r"for:\s*(?:TYPE|INTERFACE|PROCEDURE|function)\s+([^:]+)", line)
    if m:
        name = m.group(1).strip().rstrip(":")
        return f"implement {name} per linked architecture and requirements"
    return f"implement {section_hint} per linked architecture and requirements" if section_hint else None


def count_meaningful_leads(text: str) -> int:
    return sum(
        1
        for line in text.splitlines()
        if line.startswith("- [IMPL-") and "— How:" in line and "Block implements" not in line
    )


def rewrite_sidecar(path: Path, dry_run: bool = False) -> bool:
    text = path.read_text()
    if "Block implements documented behavior" not in text:
        return False
    if count_meaningful_leads(text) >= 2:
        return False

    impl, archs, reqs = parse_h1_tokens(text)
    current_section = ""
    out_lines: list[str] = []
    changed = False

    for line in text.splitlines():
        if line.startswith("## "):
            current_section = line[3:].strip()

        if line.strip().startswith("- [IMPL-") and "Block implements" in line:
            how = how_from_boilerplate(line, current_section)
            if how:
                out_lines.append(f"- [{impl}] {archs} {reqs} — How: {how}.")
                changed = True
                continue

        if line.startswith("# [IMPL-") and line.count("# [IMPL-") > 0 and len(out_lines) > 2:
            # drop duplicate H1-style lines after first header
            if out_lines and out_lines[0].startswith("# ["):
                continue

        out_lines.append(de_goify_line(line))

    if not changed:
        return False

    # Ensure single top H1
    while out_lines and out_lines[0].strip() == "":
        out_lines.pop(0)
    if out_lines and not out_lines[0].startswith(f"# [{impl}]"):
        out_lines.insert(0, f"# [{impl}] {archs} {reqs}")
        out_lines.insert(1, "")

    new_text = "\n".join(out_lines).rstrip() + "\n"
    if not dry_run:
        path.write_text(new_text)
    return True


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--token", action="append")
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    tokens = args.token or [p.stem.replace("-pseudocode", "") for p in sorted(IMPL_DIR.glob("IMPL-*-pseudocode.md"))]
    changed = 0
    for token in tokens:
        if token in SKIP_TOKENS:
            continue
        path = IMPL_DIR / f"{token}-pseudocode.md"
        if rewrite_sidecar(path, dry_run=args.dry_run):
            changed += 1
            print(f"rewrote {token}")
    print(f"Rewrote {changed} sidecars")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
