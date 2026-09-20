#!/usr/bin/env python3
"""D15: remove paraphrased // [IMPL-*] banners when // - block leads exist. [REQ-PSEUDOCODE_FORMAL_VERIFICATION]"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

from lead_hygiene import REPO, run_banner_audit, strip_redundant_banners

DEFAULT_OUT = (
    REPO / "working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/evidence/d15-banner-audit.json"
)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, default=DEFAULT_OUT)
    parser.add_argument("--fail-on-redundant", action="store_true")
    parser.add_argument("--apply", action="store_true", help="Rewrite files (default dry-run)")
    args = parser.parse_args()

    before = run_banner_audit(args.output)
    total = sum(r.redundant_banner_count for r in before)
    print(f"DEBUG: D15 redundant banners: {total} lines in {len(before)} files")
    for r in before[:15]:
        print(f"  {r.path}: {r.redundant_banner_count} ({', '.join(r.tokens[:4])}{'…' if len(r.tokens)>4 else ''})")

    if not args.apply:
        if args.fail_on_redundant and total:
            return 1
        return 0

    removed_total = 0
    for path in REPO.rglob("*.go"):
        rel = str(path.relative_to(REPO)).replace("\\", "/")
        if "vendor" in rel.split("/") or rel.startswith("scripts/testdata/"):
            continue
        content = path.read_text(errors="replace")
        new_content, n = strip_redundant_banners(content)
        if n:
            path.write_text(new_content)
            removed_total += n
            print(f"DEBUG: stripped {n} banner(s) from {rel}")
    print(f"DEBUG: removed {removed_total} redundant banner lines total")
    run_banner_audit(args.output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
