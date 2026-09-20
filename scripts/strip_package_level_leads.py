#!/usr/bin/env python3
"""Strip erroneous package-level IMPL block leads from test files. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

from lead_hygiene import REPO, iter_go_files, load_impl_code_files, strip_package_leads


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--wave", choices=["a", "b", "c", "all"], default="all")
    parser.add_argument("--apply", action="store_true", help="Write changes (default dry-run)")
    parser.add_argument("paths", nargs="*", help="Optional explicit files relative to repo root")
    args = parser.parse_args()

    code_map = load_impl_code_files()
    paths: list[Path]
    if args.paths:
        paths = [REPO / p for p in args.paths]
    elif args.wave == "all":
        paths = []
        for w in ("a", "b", "c"):
            paths.extend(iter_go_files(w))
        paths = sorted(set(paths))
    else:
        paths = iter_go_files(args.wave)

    total_removed = 0
    files_changed = 0
    for path in paths:
        if not path.is_file():
            continue
        rel = str(path.relative_to(REPO)).replace("\\", "/")
        content = path.read_text(errors="replace")
        new_content, removed = strip_package_leads(content, rel, code_map)
        if removed:
            files_changed += 1
            total_removed += removed
            mode = "APPLY" if args.apply else "DRY-RUN"
            print(f"{mode}: {rel} removed {removed} package-level lead(s)")
            if args.apply:
                path.write_text(new_content)

    print(f"DEBUG: {files_changed} file(s), {total_removed} lead line(s) removed")
    return 0


if __name__ == "__main__":
    sys.exit(main())
