#!/usr/bin/env python3
"""Audit package-level IMPL block leads. [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

from lead_hygiene import REPO, run_audit, SUSPECT_MIN_PACKAGE_LEADS

DEFAULT_OUT = (
    REPO
    / "working/REQ-PSEUDOCODE_FORMAL_VERIFICATION/evidence/archive-lead-hygiene-20260919/bulk-lead-audit.json"
)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, default=DEFAULT_OUT)
    parser.add_argument("--fail-on-suspect", action="store_true")
    parser.add_argument("--markdown", type=Path, default=None)
    args = parser.parse_args()

    results = run_audit(args.output)
    suspects = [r for r in results if r.suspect]
    print(f"DEBUG: audited {len(results)} Go files, {len(suspects)} suspect (package leads >= {SUSPECT_MIN_PACKAGE_LEADS} or ZIP fingerprint)")
    for r in suspects[:20]:
        print(f"  {r.path}: count={r.package_level_lead_count} action={r.action} pool={r.in_check_leads_pool}")
    if len(suspects) > 20:
        print(f"  ... and {len(suspects) - 20} more")

    if args.markdown:
        lines = [
            "# Bulk package-level lead audit",
            "",
            f"Suspect files: **{len(suspects)}**",
            "",
            "| path | package_leads | action | check_leads_pool |",
            "|------|---------------|--------|------------------|",
        ]
        for r in suspects:
            lines.append(
                f"| `{r.path}` | {r.package_level_lead_count} | {r.action} | {r.in_check_leads_pool} |"
            )
        args.markdown.parent.mkdir(parents=True, exist_ok=True)
        args.markdown.write_text("\n".join(lines) + "\n")

    if args.fail_on_suspect and suspects:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
