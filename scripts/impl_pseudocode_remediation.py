#!/usr/bin/env python3
"""Post-wave IMPL pseudocode remediation helpers.

- audit_impl_alignment: report sidecar vs Go lead mismatches and boilerplate
- seed_improvement_queue: write tied/docs/impl-pseudocode-improvement-queue.yaml
- add_test_leads_from_sidecar: prepend // - leads to *_test.go files by package mapping
"""
from __future__ import annotations

import argparse
import re
import sys
from collections import defaultdict
from pathlib import Path

import yaml

REPO = Path(__file__).resolve().parents[1]
IMPL_DIR = REPO / "tied" / "implementation-decisions"
CHECKLIST = REPO / "tied/docs/impl-pseudocode-sync-checklist.yaml"
QUEUE = REPO / "tied/docs/impl-pseudocode-improvement-queue.yaml"

BOILERPLATE = (
    "Block implements documented behavior",
    "per resource-management atomic I/O",
)

TOKEN_TEST_DIRS: dict[str, list[str]] = {
    "IMPL-ATOMIC_OPS": ["pkg/fileops"],
    "IMPL-AUTO_DETECTION": ["."],
    "IMPL-CFG_006": ["."],
    "IMPL-CFG_INHERITANCE_PATH_RESOLUTION": ["pkg/config"],
    "IMPL-CFG_MERGE_BEHAVIOR_REGISTRY": ["pkg/config"],
    "IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE": ["pkg/config"],
    "IMPL-CFG_PRECEDENCE_FIX": ["pkg/config"],
    "IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX": ["pkg/config"],
    "IMPL-CFG_MIXED_MODE_MERGE_FIX": ["pkg/config"],
    "IMPL-CFG_QUOTED_KEY_PREFIX": ["pkg/config"],
    "IMPL-CLI_FRAMEWORK": ["pkg/cli"],
    "IMPL-CONFIG_OUTPUT_GROUPING": ["."],
    "IMPL-CONFIG_SCHEMA_FLEX": ["."],
    "IMPL-CONFIG_STRUCT": ["."],
    "IMPL-CONTEXT_OPS": ["pkg/resources"],
    "IMPL-CUSTOMIZABLE_FORMAT_STRINGS": ["."],
    "IMPL-DATA_MODELS": ["."],
    "IMPL-DELAYED_OUTPUT": ["."],
    "IMPL-DIFF_COMMAND": ["."],
    "IMPL-DIRECTORY_COMPARISON": ["pkg/fileops", "."],
    "IMPL-DUAL_FORMATTING": [".", "pkg/formatter"],
    "IMPL-EXCLUDE_MERGE_FIX": ["pkg/config"],
    "IMPL-EXCLUSION_PATTERNS": ["pkg/fileops"],
    "IMPL-FILE_OPERATIONS": ["pkg/fileops"],
    "IMPL-FILE_STATISTICS": ["."],
    "IMPL-FILE_STATISTICS_TEMPLATE_FIX": ["."],
    "IMPL-GIT_CLI": ["pkg/git"],
    "IMPL-GIT_DIRTY_CONFIG": ["."],
    "IMPL-INCREMENTAL_DUPLICATE_PREVENTION": ["."],
    "IMPL-LIST_FORMAT_SAFETY": ["."],
    "IMPL-LIST_LIMIT": ["."],
    "IMPL-PACKAGE_EXTRACTION": ["."],
    "IMPL-PROCESSING_PATTERNS": ["pkg/processing"],
    "IMPL-RESOURCE_MANAGER": ["pkg/resources"],
    "IMPL-STRUCTURED_ERRORS": ["pkg/errors", "."],
    "IMPL-TESTING": ["pkg/testutil"],
    "IMPL-TOKEN_COVERAGE_AUDIT": ["."],
    "IMPL-ZIP_FORMAT": ["."],
    "IMPL-TOKEN_SYSTEM": ["cmd/token-suggester"],
    "IMPL-TRACEABILITY": ["test/metrics"],
    "IMPL-SPEC_CTL": ["test/specconformance", "internal/specparse", "internal/specmodel"],
    "IMPL-DOC_ENHANCEMENT": ["test/metrics", "test/specconformance"],
}


def load_checklist_tokens() -> list[dict]:
    text = CHECKLIST.read_text()
    m = re.search(r"^impl_tokens:", text, re.M)
    if not m:
        raise RuntimeError("impl_tokens not found in checklist")
    data = yaml.safe_load(text[m.start() :])
    return sorted(data["impl_tokens"], key=lambda x: x["process_order"])


def sidecar_leads(token: str) -> list[str]:
    path = IMPL_DIR / f"{token}-pseudocode.md"
    if not path.exists():
        return []
    return [line[2:] for line in path.read_text().splitlines() if line.startswith("- [IMPL-")]


def _skip_go_scan(path: Path) -> bool:
    if "vendor" in path.parts:
        return True
    # Hygiene unit fixtures use realistic IMPL tokens but are not production leads.
    return "scripts/testdata" in path.as_posix()


def go_block_leads(token: str) -> list[str]:
    leads: list[str] = []
    for gf in REPO.rglob("*.go"):
        if _skip_go_scan(gf):
            continue
        for line in gf.read_text(errors="replace").splitlines():
            m = re.match(r"^\s*// - (\[IMPL-[^\]]+\].*)$", line)
            if m and f"[{token}]" in m.group(1):
                leads.append(m.group(1))
    return leads


def test_block_leads(token: str) -> list[str]:
    leads: list[str] = []
    for gf in REPO.rglob("*_test.go"):
        if _skip_go_scan(gf):
            continue
        for line in gf.read_text(errors="replace").splitlines():
            m = re.match(r"^\s*// - (\[IMPL-[^\]]+\].*)$", line)
            if m and f"[{token}]" in m.group(1):
                leads.append(m.group(1))
    return leads


def count_boilerplate_go(token: str) -> int:
    n = 0
    for gf in REPO.rglob("*.go"):
        if _skip_go_scan(gf):
            continue
        for line in gf.read_text(errors="replace").splitlines():
            if f"[{token}]" in line and any(b in line for b in BOILERPLATE):
                n += 1
    return n


def _sidecar_leads_for_alignment(token: str) -> set[str]:
    """Sidecar block leads used for three-way Go sync (excludes INFRA policy blocks)."""
    sc = set(sidecar_leads(token))
    if token == "IMPL-TOKEN_COVERAGE_AUDIT":
        # HYGIENE_POLICY is COVER-001 waived (INFRA-policy-sidecar-block); proof is hygiene scripts/tests.
        sc = {lead for lead in sc if "audit_package_level_leads.py" not in lead}
    return sc


def classify_token(row: dict) -> tuple[str, list[str]]:
    token = row["token"]
    tier = row.get("tier", "")
    go_refs = row.get("go_refs", 0)
    sc = _sidecar_leads_for_alignment(token)
    gc = set(go_block_leads(token))
    tc = set(test_block_leads(token))
    issues: list[str] = []

    boiler = count_boilerplate_go(token)
    if boiler:
        issues.append(f"boilerplate_go:{boiler}")

    if tier != "C_doc_process" and go_refs > 0:
        missing_go = sc - gc
        extra_go = gc - sc
        # Dual-token leads tagged CONFIG_DISPLAY_FLATTENING+CFG_006 are not CFG_006 sidecar blocks
        if token == "IMPL-CFG_006":
            extra_go = {l for l in extra_go if not l.startswith("[IMPL-CONFIG_DISPLAY_FLATTENING]")}
        if missing_go:
            issues.append(f"missing_in_go:{len(missing_go)}")
        if extra_go:
            issues.append(f"extra_in_go:{len(extra_go)}")
        if sc and not tc:
            issues.append("test_leads_missing")
        elif sc and len(sc - tc) > len(sc) // 2:
            issues.append(f"test_gap:{len(sc - tc)}/{len(sc)}")

    path = IMPL_DIR / f"{token}-pseudocode.md"
    if path.exists():
        c = path.read_text()
        if re.search(r"\b(func |:= |fmt\.|strings\.|import )", c):
            if tier != "C_doc_process":
                issues.append("goisms")
        blocks = len(re.findall(r"^## [A-Z_]", c, re.M))
        proc_kw = len(re.findall(r"^PROCEDURE ", c, re.M))
        if blocks and proc_kw < blocks * 0.5:
            issues.append("no_procedure_bodies")

    if token in ("IMPL-CONFIG_STRUCT", "IMPL-CFG_006", "IMPL-CONFIG_DISPLAY_FLATTENING"):
        if path.exists() and "reflectConfigFields" in path.read_text():
            pass  # expected cross-IMPL reflection; not a blocking issue post-LEAP

    if not issues:
        return "P3", []
    # Non-blocking residual flags after remediation pass
    non_blocking = {"goisms", "test_gap", "no_procedure_bodies", "cfg_duplication"}
    blocking = [i for i in issues if not any(i.startswith(nb) for nb in non_blocking)]
    if not blocking:
        return "P3", issues
    if boiler or (sc - gc and go_refs > 0) or (gc - sc and len(gc - sc) > 2):
        return "P0", issues
    if "test_leads_missing" in issues or "no_procedure_bodies" in issues:
        return "P1", issues
    if any(i.startswith("test_gap") for i in issues):
        return "P2", issues
    return "P3", issues


def seed_queue() -> None:
    rows = []
    for row in load_checklist_tokens():
        priority, issues = classify_token(row)
        status = "pilot_complete" if not issues or priority == "P3" else "pending"
        rows.append(
            {
                "token": row["token"],
                "process_order": row["process_order"],
                "tier": row.get("tier", ""),
                "priority": priority,
                "issues": issues,
                "status": status,
                "notes": row.get("notes", ""),
            }
        )
    header = """# IMPL Pseudocode Improvement Queue (post-wave audit)
# Seeded by scripts/impl_pseudocode_remediation.py seed-queue
# Agent guide: tied/docs/impl-deep-sync-agent-guide.md
# status: pending | in_progress | pilot_complete
"""
    body = yaml.safe_dump(
        {"description": "Per-token remediation tracking after waves 1-5 checklist validation.", "tokens": rows},
        sort_keys=False,
        allow_unicode=True,
    )
    QUEUE.write_text(header + body)


def add_test_leads(token: str, dry_run: bool = False) -> int:
    """Insert each sidecar lead once, immediately before a matching Test/Benchmark func (never at package)."""
    leads = sidecar_leads(token)
    if not leads:
        return 0
    dirs = TOKEN_TEST_DIRS.get(token, [])
    if not dirs:
        return 0
    added = 0
    token_suffix = token.replace("IMPL-", "").lower()
    func_re = re.compile(r"^func (Test|Benchmark)(\w+)\(", re.MULTILINE)
    for rel in dirs:
        base = REPO / rel
        if not base.exists():
            continue
        for tf in base.rglob("*_test.go"):
            if "vendor" in tf.parts:
                continue
            content = tf.read_text(errors="replace")
            file_added = 0
            for i, lead in enumerate(leads):
                needle = f"// - {lead}"
                if needle in content:
                    continue
                matches = list(func_re.finditer(content))
                if not matches:
                    continue
                target = None
                for m in matches:
                    tname = m.group(2).lower()
                    if token_suffix in tname or i == 0:
                        target = m
                        break
                if target is None:
                    target = matches[0]
                insert_at = target.start()
                prefix = content[:insert_at].rstrip("\n")
                content = prefix + "\n" + needle + "\n" + content[insert_at:]
                file_added += 1
            if file_added:
                added += file_added
                if not dry_run:
                    tf.write_text(content)
    return added


def add_all_test_leads(dry_run: bool = False) -> None:
    raise SystemExit(
        "add-test-leads --all is disabled: package-level paste caused bulk lead pollution. "
        "Use --token IMPL-* for func-scoped insertion, or sync_impl_block_leads.py for production."
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="cmd", required=True)
    sub.add_parser("seed-queue")
    p_test = sub.add_parser("add-test-leads")
    p_test.add_argument("--token", default="")
    p_test.add_argument("--all", action="store_true")
    p_test.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()
    if args.cmd == "seed-queue":
        seed_queue()
        print(f"Wrote {QUEUE}")
    elif args.cmd == "add-test-leads":
        if args.all:
            add_all_test_leads(dry_run=args.dry_run)
        elif args.token:
            n = add_test_leads(args.token, dry_run=args.dry_run)
            print(f"Added {n} lead lines for {args.token}")
        else:
            sys.exit("Specify --token or --all")


if __name__ == "__main__":
    main()
